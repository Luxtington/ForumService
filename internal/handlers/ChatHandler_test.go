package handlers

import (
	"bytes"
	"encoding/json"
	"ForumService/internal/handlers/mocks"
	"ForumService/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateMessage(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockResponse   *models.ChatMessage
		mockError      error
		expectedStatus int
		userID         uint
	}{
		{
			name: "успешное создание сообщения",
			requestBody: map[string]interface{}{
				"content": "Test Message Content",
			},
			mockResponse: &models.ChatMessage{
				ID:       1,
				Content:  "Test Message Content",
				AuthorID: 1,
			},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			userID:         1,
		},
		{
			name: "неверный формат данных",
			requestBody: map[string]interface{}{
				"invalid": "data",
			},
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			userID:         1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.MockChatService{
				CreateMessageFunc: func(authorID int, content string) (*models.ChatMessage, error) {
					return tt.mockResponse, tt.mockError
				},
			}

			handler := NewChatHandler(mockService)
			router := setupTestRouter()
			router.POST("/chat/messages", func(c *gin.Context) {
				c.Set("user_id", tt.userID)
				handler.CreateMessage(c)
			})

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/chat/messages", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestGetMessages(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   []*models.ChatMessage
		mockError      error
		expectedStatus int
	}{
		{
			name: "успешное получение сообщений",
			mockResponse: []*models.ChatMessage{
				{
					ID:       1,
					Content:  "Test Message 1",
					AuthorID: 1,
				},
				{
					ID:       2,
					Content:  "Test Message 2",
					AuthorID: 2,
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "ошибка при получении сообщений",
			mockResponse:   nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.MockChatService{
				GetAllMessagesFunc: func() ([]*models.ChatMessage, error) {
					return tt.mockResponse, tt.mockError
				},
			}

			handler := NewChatHandler(mockService)
			router := setupTestRouter()
			router.GET("/chat/messages", handler.GetMessages)

			req := httptest.NewRequest("GET", "/chat/messages", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
} 