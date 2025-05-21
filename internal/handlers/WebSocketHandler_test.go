package handlers

import (
	"ForumService/internal/handlers/mocks"
	"ForumService/internal/models"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestHub_HandleWebSocket(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		userID         int
		expectedStatus int
	}{
		{
			name:           "успешное подключение",
			username:       "testuser",
			userID:         1,
			expectedStatus: http.StatusSwitchingProtocols,
		},
		{
			name:           "отсутствует имя пользователя",
			username:       "",
			userID:         1,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "отсутствует ID пользователя",
			username:       "testuser",
			userID:         0,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockChatRepo := &mocks.MockChatRepository{
				CreateMessageFunc: func(authorID int, content string) (*models.ChatMessage, error) {
					return &models.ChatMessage{
						ID:       1,
						Content:  content,
						AuthorID: authorID,
					}, nil
				},
			}

			hub := NewHub(mockChatRepo)
			go hub.Run()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), "username", tt.username)
				ctx = context.WithValue(ctx, "user_id", tt.userID)
				hub.HandleWebSocket(w, r.WithContext(ctx))
			}))
			defer server.Close()

			wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
			ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
			if err != nil {
				if tt.expectedStatus == http.StatusSwitchingProtocols {
					t.Errorf("Ожидалось успешное подключение, но получили ошибку: %v", err)
				}
				return
			}
			defer ws.Close()

			// Проверяем, что клиент был зарегистрирован
			time.Sleep(100 * time.Millisecond)
			assert.Equal(t, 1, len(hub.Clients))
		})
	}
}

func TestHub_MessageHandling(t *testing.T) {
	mockChatRepo := &mocks.MockChatRepository{
		CreateMessageFunc: func(authorID int, content string) (*models.ChatMessage, error) {
			return &models.ChatMessage{
				ID:       1,
				Content:  content,
				AuthorID: authorID,
			}, nil
		},
	}

	hub := NewHub(mockChatRepo)
	go hub.Run()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "username", "testuser")
		ctx = context.WithValue(ctx, "user_id", 1)
		hub.HandleWebSocket(w, r.WithContext(ctx))
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Не удалось подключиться к WebSocket: %v", err)
	}
	defer ws.Close()

	// Отправляем тестовое сообщение
	message := Message{
		Type:       "message",
		Content:    "Test message",
		AuthorID:   1,
		AuthorName: "testuser",
		CreatedAt:  time.Now().Format(time.RFC3339),
	}
	messageBytes, _ := json.Marshal(message)
	err = ws.WriteMessage(websocket.TextMessage, messageBytes)
	assert.NoError(t, err)

	// Ждем немного, чтобы сообщение было обработано
	time.Sleep(100 * time.Millisecond)

	// Проверяем, что сообщение было создано в репозитории
	assert.Equal(t, 1, mockChatRepo.CreateMessageCallCount())
}

func TestHub_ClientDisconnection(t *testing.T) {
	mockChatRepo := &mocks.MockChatRepository{
		CreateMessageFunc: func(authorID int, content string) (*models.ChatMessage, error) {
			return &models.ChatMessage{
				ID:       1,
				Content:  content,
				AuthorID: authorID,
			}, nil
		},
	}

	hub := NewHub(mockChatRepo)
	go hub.Run()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "username", "testuser")
		ctx = context.WithValue(ctx, "user_id", 1)
		hub.HandleWebSocket(w, r.WithContext(ctx))
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Не удалось подключиться к WebSocket: %v", err)
	}

	// Проверяем, что клиент был зарегистрирован
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 1, len(hub.Clients))

	// Закрываем соединение
	ws.Close()

	// Ждем немного, чтобы клиент был отрегистрирован
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 0, len(hub.Clients))
}

func TestHub_MessageErrorHandling(t *testing.T) {
	tests := []struct {
		name           string
		message        interface{}
		mockError      error
		expectedStatus int
	}{
		{
			name: "ошибка при создании сообщения",
			message: Message{
				Type:       "message",
				Content:    "Test message",
				AuthorID:   1,
				AuthorName: "testuser",
				CreatedAt:  time.Now().Format(time.RFC3339),
			},
			mockError:      errors.New("ошибка создания сообщения"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "некорректный формат сообщения",
			message:        "invalid message",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "пустое сообщение",
			message: Message{
				Type:       "message",
				Content:    "",
				AuthorID:   1,
				AuthorName: "testuser",
				CreatedAt:  time.Now().Format(time.RFC3339),
			},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockChatRepo := &mocks.MockChatRepository{
				CreateMessageFunc: func(authorID int, content string) (*models.ChatMessage, error) {
					return nil, tt.mockError
				},
			}

			hub := NewHub(mockChatRepo)
			go hub.Run()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), "username", "testuser")
				ctx = context.WithValue(ctx, "user_id", 1)
				hub.HandleWebSocket(w, r.WithContext(ctx))
			}))
			defer server.Close()

			wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
			ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
			if err != nil {
				t.Fatalf("Не удалось подключиться к WebSocket: %v", err)
			}
			defer ws.Close()

			// Отправляем тестовое сообщение
			var messageBytes []byte
			if str, ok := tt.message.(string); ok {
				messageBytes = []byte(str)
			} else {
				messageBytes, _ = json.Marshal(tt.message)
			}

			err = ws.WriteMessage(websocket.TextMessage, messageBytes)
			assert.NoError(t, err)

			// Ждем немного, чтобы сообщение было обработано
			time.Sleep(100 * time.Millisecond)

			// Проверяем, что сообщение было обработано
			if tt.mockError != nil {
				assert.Equal(t, 1, mockChatRepo.CreateMessageCallCount())
			}
		})
	}
}

func TestHub_ConnectionErrorHandling(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		userID         int
		expectedStatus int
	}{
		{
			name:           "неверный формат ID пользователя",
			username:       "testuser",
			userID:         -1,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "пустое имя пользователя",
			username:       "",
			userID:         1,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "отсутствует ID пользователя",
			username:       "testuser",
			userID:         0,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockChatRepo := &mocks.MockChatRepository{
				CreateMessageFunc: func(authorID int, content string) (*models.ChatMessage, error) {
					return &models.ChatMessage{
						ID:       1,
						Content:  content,
						AuthorID: authorID,
					}, nil
				},
			}

			hub := NewHub(mockChatRepo)
			go hub.Run()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), "username", tt.username)
				ctx = context.WithValue(ctx, "user_id", tt.userID)
				hub.HandleWebSocket(w, r.WithContext(ctx))
			}))
			defer server.Close()

			wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
			ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
			if err != nil {
				if tt.expectedStatus == http.StatusSwitchingProtocols {
					t.Errorf("Ожидалось успешное подключение, но получили ошибку: %v", err)
				}
				return
			}
			defer ws.Close()

			// Проверяем, что клиент не был зарегистрирован
			time.Sleep(100 * time.Millisecond)
			assert.Equal(t, 0, len(hub.Clients))
		})
	}
}

func TestHub_ReadWriteErrorHandling(t *testing.T) {
	mockChatRepo := &mocks.MockChatRepository{
		CreateMessageFunc: func(authorID int, content string) (*models.ChatMessage, error) {
			return &models.ChatMessage{
				ID:       1,
				Content:  content,
				AuthorID: authorID,
			}, nil
		},
	}

	hub := NewHub(mockChatRepo)
	go hub.Run()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "username", "testuser")
		ctx = context.WithValue(ctx, "user_id", 1)
		hub.HandleWebSocket(w, r.WithContext(ctx))
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Не удалось подключиться к WebSocket: %v", err)
	}

	// Проверяем, что клиент был зарегистрирован
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 1, len(hub.Clients))

	// Закрываем соединение на стороне сервера
	ws.Close()

	// Ждем немного, чтобы клиент был отрегистрирован
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 0, len(hub.Clients))

	// Пытаемся отправить сообщение в закрытое соединение
	message := Message{
		Type:       "message",
		Content:    "Test message",
		AuthorID:   1,
		AuthorName: "testuser",
		CreatedAt:  time.Now().Format(time.RFC3339),
	}
	messageBytes, _ := json.Marshal(message)
	err = ws.WriteMessage(websocket.TextMessage, messageBytes)
	assert.Error(t, err)
} 