package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"ForumService/internal/models"
	"github.com/gorilla/websocket"
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
			log.Printf("Клиент зарегистрирован: %s (ID: %d)", client.Username, client.UserID)
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				log.Printf("Клиент отрегистрирован: %s (ID: %d)", client.Username, client.UserID)
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
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Ошибка при обновлении соединения до WebSocket: %v", err)
		return
	}

	// Получаем данные из контекста
	username, exists := r.Context().Value("username").(string)
	if !exists {
		log.Printf("Имя пользователя не найдено в контексте")
		conn.Close()
		return
	}

	userID, exists := r.Context().Value("user_id").(uint)
	if !exists {
		log.Printf("ID пользователя не найден в контексте")
		conn.Close()
		return
	}

	log.Printf("Создание нового WebSocket клиента: username=%s, userID=%d", username, userID)

	client := &Client{
		Conn:     conn,
		Send:     make(chan []byte, 256),
		Username: username,
		UserID:   int(userID),
	}

	h.Register <- client

	go h.WritePump(client)
	go h.ReadPump(client)
}

func (h *Hub) WritePump(c *Client) {
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
				log.Printf("Ошибка при получении writer: %v", err)
				return
			}
			w.Write(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				log.Printf("Ошибка при закрытии writer: %v", err)
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Ошибка при отправке ping: %v", err)
				return
			}
		}
	}
}

func (h *Hub) ReadPump(c *Client) {
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
				log.Printf("Ошибка чтения сообщения: %v", err)
			}
			break
		}

		log.Printf("Получено сообщение от пользователя %s (ID: %d): %s", c.Username, c.UserID, string(message))

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Ошибка при разборе сообщения: %v", err)
			continue
		}

		if msg.Content == "" {
			log.Printf("Получено пустое сообщение от пользователя %s", c.Username)
			continue
		}

		log.Printf("Попытка сохранения сообщения в БД: authorID=%d, content=%s", c.UserID, msg.Content)

		chatMessage, err := h.ChatRepo.CreateMessage(c.UserID, msg.Content)
		if err != nil {
			log.Printf("Ошибка при сохранении сообщения в БД: %v", err)
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

		log.Printf("Сообщение успешно сохранено в БД: ID=%d", chatMessage.ID)

		msg.AuthorName = c.Username
		msg.AuthorID = c.UserID
		msg.CreatedAt = chatMessage.CreatedAt.Format(time.RFC3339)

		messageBytes, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Ошибка при сериализации сообщения: %v", err)
			continue
		}

		log.Printf("Отправка сообщения всем клиентам: %s", string(messageBytes))
		h.Broadcast <- messageBytes
	}
} 