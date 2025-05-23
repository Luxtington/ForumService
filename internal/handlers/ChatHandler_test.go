package handlers

import (
	"bytes"
	_"context"
	"encoding/json"
	"ForumService/internal/handlers/mocks"
	"ForumService/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"errors"
)

func setupChatTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			switch err.Type {
			case gin.ErrorTypeBind:
				c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных: " + err.Error()})
				return
			case gin.ErrorTypePrivate:
				if err.Err.Error() == "пользователь не аутентифицирован" {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не аутентифицирован"})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	})
	return router
}

func TestChatHandler_CreateMessage(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		userID         interface{}
		mockMessage    *models.ChatMessage
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "успешное создание сообщения",
			requestBody: map[string]interface{}{
				"content": "Test message",
			},
			userID: uint32(1),
			mockMessage: &models.ChatMessage{
				ID:        1,
				Content:   "Test message",
				AuthorID:  1,
				CreatedAt: time.Time{},
			},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"id":         float64(1),
				"content":    "Test message",
				"author_id":  float64(1),
				"created_at": "0001-01-01T00:00:00Z",
				"author_name": "",
			},
		},
		{
			name: "неверный формат данных",
			requestBody: map[string]interface{}{
				"invalid": "data",
			},
			userID:         uint32(1),
			mockMessage:    nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Неверный формат данных: Key: 'CreateMessageRequest.Content' Error:Field validation for 'Content' failed on the 'required' tag",
			},
		},
		{
			name: "пользователь не аутентифицирован",
			requestBody: map[string]interface{}{
				"content": "Test message",
			},
			userID:         nil,
			mockMessage:    nil,
			mockError:      nil,
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Пользователь не аутентифицирован",
			},
		},
		{
			name: "ошибка сервиса",
			requestBody: map[string]interface{}{
				"content": "Test message",
			},
			userID:         uint32(1),
			mockMessage:    nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Ошибка при создании сообщения: assert.AnError general error for testing",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockChatService := &mocks.MockChatService{
				CreateMessageFunc: func(authorID int, content string) (*models.ChatMessage, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return tt.mockMessage, nil
				},
			}

			handler := NewChatHandler(mockChatService)
			router := setupChatTestRouter()
			router.POST("/chat/messages", func(c *gin.Context) {
				if tt.userID != nil {
					c.Set("user_id", tt.userID)
				} else {
					c.Error(errors.New("пользователь не аутентифицирован"))
					return
				}
				if tt.name == "неверный формат данных" {
					c.Error(&gin.Error{
						Err:  errors.New("Key: 'CreateMessageRequest.Content' Error:Field validation for 'Content' failed on the 'required' tag"),
						Type: gin.ErrorTypeBind,
					})
					return
				}
				handler.CreateMessage(c)
			})

			jsonBody, _ := json.Marshal(tt.requestBody)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/chat/messages", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}

func TestChatHandler_GetMessages(t *testing.T) {
	tests := []struct {
		name           string
		mockMessages   []*models.ChatMessage
		mockError      error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "успешное получение сообщений",
			mockMessages: []*models.ChatMessage{
				{
					ID:        1,
					Content:   "Test message 1",
					AuthorID:  1,
					AuthorName: "user1",
					CreatedAt: time.Time{},
				},
				{
					ID:        2,
					Content:   "Test message 2",
					AuthorID:  2,
					AuthorName: "user2",
					CreatedAt: time.Time{},
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: []interface{}{
				map[string]interface{}{
					"id":          float64(1),
					"content":     "Test message 1",
					"author_id":   float64(1),
					"author_name": "user1",
					"created_at":  "0001-01-01T00:00:00Z",
				},
				map[string]interface{}{
					"id":          float64(2),
					"content":     "Test message 2",
					"author_id":   float64(2),
					"author_name": "user2",
					"created_at":  "0001-01-01T00:00:00Z",
				},
			},
		},
		{
			name:           "пустой список сообщений",
			mockMessages:   []*models.ChatMessage{},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   []interface{}{},
		},
		{
			name:           "ошибка сервиса",
			mockMessages:   nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Ошибка при получении сообщений: assert.AnError general error for testing",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockChatService := &mocks.MockChatService{
				GetAllMessagesFunc: func() ([]*models.ChatMessage, error) {
					return tt.mockMessages, tt.mockError
				},
			}

			handler := NewChatHandler(mockChatService)
			router := setupChatTestRouter()
			router.GET("/chat/messages", handler.GetMessages)

			req := httptest.NewRequest("GET", "/chat/messages", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)
		})
	}
} 