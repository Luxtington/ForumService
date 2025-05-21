package handlers

import (
	"bytes"
	"encoding/json"
	"ForumService/internal/handlers/mocks"
	"ForumService/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

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
			userID: uint(1),
			mockMessage: &models.ChatMessage{
				ID:        1,
				Content:   "Test message",
				AuthorID:  1,
				AuthorName: "",
				CreatedAt: time.Time{},
			},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"id":          float64(1),
				"content":     "Test message",
				"author_id":   float64(1),
				"author_name": "",
				"created_at":  "0001-01-01T00:00:00Z",
			},
		},
		{
			name: "неверный формат данных",
			requestBody: map[string]interface{}{
				"invalid": "data",
			},
			userID:         uint(1),
			mockMessage:    nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "неверный формат данных",
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
				"error": "пользователь не аутентифицирован",
			},
		},
		{
			name: "ошибка сервиса",
			requestBody: map[string]interface{}{
				"content": "Test message",
			},
			userID:         uint(1),
			mockMessage:    nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": assert.AnError.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockChatService := &mocks.MockChatService{
				CreateMessageFunc: func(authorID int, content string) (*models.ChatMessage, error) {
					return tt.mockMessage, tt.mockError
				},
			}

			handler := NewChatHandler(mockChatService)
			router := setupViewsTestRouter()
			router.POST("/messages", func(c *gin.Context) {
				if tt.userID != nil {
					c.Set("user_id", tt.userID)
				}
				handler.CreateMessage(c)
			})

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/messages", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

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
					ID:         1,
					Content:    "Test message 1",
					AuthorID:   1,
					AuthorName: "",
					CreatedAt:  time.Time{},
				},
				{
					ID:         2,
					Content:    "Test message 2",
					AuthorID:   2,
					AuthorName: "",
					CreatedAt:  time.Time{},
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: []interface{}{
				map[string]interface{}{
					"id":          float64(1),
					"content":     "Test message 1",
					"author_id":   float64(1),
					"author_name": "",
					"created_at":  "0001-01-01T00:00:00Z",
				},
				map[string]interface{}{
					"id":          float64(2),
					"content":     "Test message 2",
					"author_id":   float64(2),
					"author_name": "",
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
				"error": assert.AnError.Error(),
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
			router := setupViewsTestRouter()
			router.GET("/messages", handler.GetMessages)

			req := httptest.NewRequest("GET", "/messages", nil)
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