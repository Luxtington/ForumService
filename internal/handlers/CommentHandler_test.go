package handlers

import (
	"bytes"
	_"context"
	"encoding/json"
	"errors"
	"ForumService/internal/handlers/mocks"
	"ForumService/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupCommentTestRouter() *gin.Engine {
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
				if err.Err.Error() == "нет прав для удаления комментария" {
					c.JSON(http.StatusForbidden, gin.H{"error": "Нет прав для удаления комментария"})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	})
	return router
}

func TestCommentHandler_CreateComment(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		userID         interface{}
		mockComment    *models.Comment
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "успешное создание комментария",
			requestBody: map[string]interface{}{
				"content": "Test comment",
				"post_id": 1,
			},
			userID: uint32(1),
			mockComment: &models.Comment{
				ID:        1,
				Content:   "Test comment",
				AuthorID:  1,
				PostID:    1,
				CreatedAt: time.Time{},
			},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"id":          float64(1),
				"content":     "Test comment",
				"author_id":   float64(1),
				"post_id":     float64(1),
				"created_at":  "0001-01-01T00:00:00Z",
				"author_name": "",
				"can_delete":  false,
			},
		},
		{
			name: "пользователь не аутентифицирован",
			requestBody: map[string]interface{}{
				"content": "Test comment",
				"post_id": 1,
			},
			userID:         nil,
			mockComment:    nil,
			mockError:      nil,
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Пользователь не аутентифицирован",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCommentService := &mocks.MockCommentService{
				CreateCommentFunc: func(postID int, authorID int, content string) (*models.Comment, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return tt.mockComment, nil
				},
			}

			handler := NewCommentHandler(mockCommentService)
			router := setupCommentTestRouter()
			router.POST("/comments", func(c *gin.Context) {
				if tt.userID != nil {
					c.Set("user_id", tt.userID)
				} else {
					c.Error(&gin.Error{
						Err:  errors.New("пользователь не аутентифицирован"),
						Type: gin.ErrorTypePrivate,
					})
					return
				}
				handler.CreateComment(c)
			})

			jsonBody, _ := json.Marshal(tt.requestBody)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/comments", bytes.NewBuffer(jsonBody))
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

func TestCommentHandler_DeleteComment(t *testing.T) {
	tests := []struct {
		name           string
		commentID      string
		userID         interface{}
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "успешное удаление комментария",
			commentID:      "1",
			userID:         uint32(1),
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
			expectedBody:   nil,
		},
		{
			name:           "пользователь не аутентифицирован",
			commentID:      "1",
			userID:         nil,
			mockError:      nil,
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Пользователь не аутентифицирован",
			},
		},
		{
			name:           "нет прав для удаления",
			commentID:      "1",
			userID:         uint32(2),
			mockError:      errors.New("нет прав для удаления комментария"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Комментарий не найден: нет прав для удаления комментария",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCommentService := &mocks.MockCommentService{
				GetCommentByIDFunc: func(id int) (*models.Comment, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return &models.Comment{ID: id, AuthorID: 1}, nil
				},
				DeleteCommentFunc: func(id int, userID int) error {
					return tt.mockError
				},
			}

			handler := NewCommentHandler(mockCommentService)
			router := setupCommentTestRouter()
			router.DELETE("/comments/:id", func(c *gin.Context) {
				if tt.userID != nil {
					c.Set("user_id", tt.userID)
				} else {
					c.Error(&gin.Error{
						Err:  errors.New("пользователь не аутентифицирован"),
						Type: gin.ErrorTypePrivate,
					})
					return
				}
				handler.DeleteComment(c)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/comments/"+tt.commentID, nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			} else {
				assert.Empty(t, w.Body.String())
			}
		})
	}
}

func TestCommentHandler_CreateChatMessage_Success(t *testing.T) {
	mockCommentService := &mocks.MockCommentService{
		CreateCommentFunc: func(postID int, authorID int, content string) (*models.Comment, error) {
			return &models.Comment{
				ID:       1,
				PostID:   0,
				AuthorID: 1,
				Content:  content,
			}, nil
		},
	}

	handler := NewCommentHandler(mockCommentService)

	router := setupCommentTestRouter()
	router.POST("/chat/messages", handler.CreateChatMessage)

	requestBody := map[string]interface{}{
		"content": "Test chat message",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/chat/messages", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCommentHandler_CreateChatMessage_InvalidData(t *testing.T) {
	mockCommentService := &mocks.MockCommentService{}
	handler := NewCommentHandler(mockCommentService)

	router := setupCommentTestRouter()
	router.POST("/chat/messages", handler.CreateChatMessage)

	requestBody := map[string]interface{}{
		"content": "",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/chat/messages", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCommentHandler_CreateChatMessage_ServiceError(t *testing.T) {
	mockCommentService := &mocks.MockCommentService{
		CreateCommentFunc: func(postID int, authorID int, content string) (*models.Comment, error) {
			return nil, errors.New("service error")
		},
	}

	handler := NewCommentHandler(mockCommentService)

	router := setupCommentTestRouter()
	router.POST("/chat/messages", handler.CreateChatMessage)

	requestBody := map[string]interface{}{
		"content": "Test chat message",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/chat/messages", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}