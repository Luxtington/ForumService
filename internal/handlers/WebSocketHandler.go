package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"ForumService/internal/models"
	"github.com/gorilla/websocket"
	"github.com/Luxtington/Shared/logger"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	Conn     *websocket.Conn
	Send     chan []byte
	Username string
	UserID   int
}

type Message struct {
	Type       string `json:"type"`
	Content    string `json:"content"`
	AuthorID   int    `json:"author_id"`
	AuthorName string `json:"author_name"`
	CreatedAt  string `json:"created_at"`
}

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	ChatRepo   ChatRepository
}

type ChatRepository interface {
	CreateMessage(authorID int, content string) (*models.ChatMessage, error)
}

func NewHub(chatRepo ChatRepository) *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		ChatRepo:   chatRepo,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			log := logger.GetLogger()
			log.Info("Клиент зарегистрирован", 
				zap.String("username", client.Username),
				zap.Int("user_id", client.UserID))
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				log := logger.GetLogger()
				log.Info("Клиент отрегистрирован",
					zap.String("username", client.Username),
					zap.Int("user_id", client.UserID))
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}

func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	username, exists := r.Context().Value("username").(string)
	if !exists || username == "" {
		log.Error("Имя пользователя не найдено в контексте или пустое")
		http.Error(w, "Имя пользователя не найдено", http.StatusBadRequest)
		return
	}

	userID, exists := r.Context().Value("user_id").(int)
	if !exists || userID <= 0 {
		log.Error("ID пользователя не найден в контексте или некорректный")
		http.Error(w, "ID пользователя не найден", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Ошибка при обновлении соединения до WebSocket", zap.Error(err))
		return
	}

	log.Info("Создание нового WebSocket клиента",
		zap.String("username", username),
		zap.Int("user_id", userID))

	client := &Client{
		Conn:     conn,
		Send:     make(chan []byte, 256),
		Username: username,
		UserID:   userID,
	}

	h.Register <- client

	go h.WritePump(client)
	go h.ReadPump(client)
}

func (h *Hub) WritePump(c *Client) {
	log := logger.GetLogger()
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Error("Ошибка при получении writer", zap.Error(err))
				return
			}
			w.Write(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				log.Error("Ошибка при закрытии writer", zap.Error(err))
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Error("Ошибка при отправке ping", zap.Error(err))
				return
			}
		}
	}
}

func (h *Hub) ReadPump(c *Client) {
	log := logger.GetLogger()
	defer func() {
		h.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512 * 1024) // 512KB
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error("Ошибка чтения сообщения", zap.Error(err))
			}
			break
		}

		log.Info("Получено сообщение от пользователя",
			zap.String("username", c.Username),
			zap.Int("user_id", c.UserID),
			zap.String("message", string(message)))

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Error("Ошибка при разборе сообщения", zap.Error(err))
			continue
		}

		if msg.Content == "" {
			log.Info("Получено пустое сообщение", zap.String("username", c.Username))
			continue
		}

		log.Info("Попытка сохранения сообщения в БД",
			zap.Int("author_id", c.UserID),
			zap.String("content", msg.Content))

		chatMessage, err := h.ChatRepo.CreateMessage(c.UserID, msg.Content)
		if err != nil {
			log.Error("Ошибка при сохранении сообщения в БД", zap.Error(err))
			errorMsg := Message{
				Type:       "error",
				Content:    "Не удалось сохранить сообщение",
				AuthorID:   c.UserID,
				AuthorName: c.Username,
				CreatedAt:  time.Now().Format(time.RFC3339),
			}
			errorBytes, _ := json.Marshal(errorMsg)
			c.Send <- errorBytes
			continue
		}

		log.Info("Сообщение успешно сохранено в БД", zap.Int("message_id", chatMessage.ID))

		msg.AuthorName = c.Username
		msg.AuthorID = c.UserID
		msg.CreatedAt = chatMessage.CreatedAt.Format(time.RFC3339)
		msg.Type = "message"

		messageBytes, err := json.Marshal(msg)
		if err != nil {
			log.Error("Ошибка при сериализации сообщения", zap.Error(err))
			continue
		}

		log.Info("Отправка сообщения всем клиентам", 
			zap.String("message", string(messageBytes)),
			zap.Int("clients_count", len(h.Clients)))

		// Отправляем сообщение всем клиентам
		for client := range h.Clients {
			select {
			case client.Send <- messageBytes:
				log.Info("Сообщение отправлено клиенту",
					zap.String("username", client.Username),
					zap.Int("user_id", client.UserID))
			default:
				log.Error("Не удалось отправить сообщение клиенту",
					zap.String("username", client.Username),
					zap.Int("user_id", client.UserID))
				close(client.Send)
				delete(h.Clients, client)
			}
		}
	}
} 