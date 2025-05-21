package handlers

import (
	"bytes"
	"encoding/json"
	_"errors"
	"ForumService/internal/handlers/mocks"
	"ForumService/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestCreateThread(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockResponse   *models.Thread
		mockError      error
		expectedStatus int
		userID         uint
	}{
		{
			name: "успешное создание треда",
			requestBody: map[string]interface{}{
				"title": "Test Thread",
			},
			mockResponse: &models.Thread{
				ID:       1,
				Title:    "Test Thread",
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
			mockService := &mocks.MockThreadService{
				CreateThreadFunc: func(title string, authorID int) (*models.Thread, error) {
					return tt.mockResponse, tt.mockError
				},
			}

			handler := NewThreadHandler(mockService)
			router := setupTestRouter()
			router.POST("/threads", func(c *gin.Context) {
				c.Set("user_id", tt.userID)
				handler.CreateThread(c)
			})

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/threads", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestGetThreadWithPosts(t *testing.T) {
	tests := []struct {
		name           string
		threadID       string
		mockThread     *models.Thread
		mockPosts      []*models.Post
		mockError      error
		expectedStatus int
	}{
		{
			name:     "успешное получение треда с постами",
			threadID: "1",
			mockThread: &models.Thread{
				ID:       1,
				Title:    "Test Thread",
				AuthorID: 1,
			},
			mockPosts: []*models.Post{
				{
					ID:       1,
					Content:  "Test Post",
					ThreadID: 1,
					AuthorID: 1,
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "неверный ID треда",
			threadID:       "invalid",
			mockThread:     nil,
			mockPosts:      nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.MockThreadService{
				GetThreadWithPostsFunc: func(id int) (*models.Thread, []*models.Post, error) {
					return tt.mockThread, tt.mockPosts, tt.mockError
				},
			}

			handler := NewThreadHandler(mockService)
			router := setupTestRouter()
			router.GET("/threads/:id", handler.GetThreadWithPosts)

			req := httptest.NewRequest("GET", "/threads/"+tt.threadID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestUpdateThread(t *testing.T) {
	tests := []struct {
		name           string
		threadID       string
		userID         uint
		userRole       string
		requestBody    map[string]interface{}
		mockThread     *models.Thread
		mockError      error
		expectedStatus int
	}{
		{
			name:     "успешное обновление треда автором",
			threadID: "1",
			userID:   1,
			userRole: "user",
			requestBody: map[string]interface{}{
				"title": "Updated Thread Title",
			},
			mockThread: &models.Thread{
				ID:       1,
				Title:    "Original Title",
				AuthorID: 1,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:     "отказ в доступе",
			threadID: "1",
			userID:   2,
			userRole: "user",
			requestBody: map[string]interface{}{
				"title": "Updated Thread Title",
			},
			mockThread: &models.Thread{
				ID:       1,
				Title:    "Original Title",
				AuthorID: 1,
			},
			mockError:      nil,
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.MockThreadService{
				GetThreadWithPostsFunc: func(id int) (*models.Thread, []*models.Post, error) {
					return tt.mockThread, nil, nil
				},
				UpdateThreadFunc: func(thread *models.Thread, userID int) error {
					return tt.mockError
				},
			}

			handler := NewThreadHandler(mockService)
			router := setupTestRouter()
			router.PUT("/threads/:id", func(c *gin.Context) {
				c.Set("user_id", tt.userID)
				c.Set("user_role", tt.userRole)
				handler.UpdateThread(c)
			})

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PUT", "/threads/"+tt.threadID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestDeleteThread(t *testing.T) {
	tests := []struct {
		name           string
		threadID       string
		userID         uint
		userRole       string
		mockThread     *models.Thread
		mockError      error
		expectedStatus int
	}{
		{
			name:     "успешное удаление треда автором",
			threadID: "1",
			userID:   1,
			userRole: "user",
			mockThread: &models.Thread{
				ID:       1,
				Title:    "Test Thread",
				AuthorID: 1,
			},
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:     "успешное удаление треда админом",
			threadID: "1",
			userID:   2,
			userRole: "admin",
			mockThread: &models.Thread{
				ID:       1,
				Title:    "Test Thread",
				AuthorID: 1,
			},
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:     "отказ в доступе",
			threadID: "1",
			userID:   2,
			userRole: "user",
			mockThread: &models.Thread{
				ID:       1,
				Title:    "Test Thread",
				AuthorID: 1,
			},
			mockError:      nil,
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.MockThreadService{
				GetThreadWithPostsFunc: func(id int) (*models.Thread, []*models.Post, error) {
					return tt.mockThread, nil, nil
				},
				DeleteThreadFunc: func(id int, userID int) error {
					return tt.mockError
				},
			}

			handler := NewThreadHandler(mockService)
			router := setupTestRouter()
			router.DELETE("/threads/:id", func(c *gin.Context) {
				c.Set("user_id", tt.userID)
				c.Set("user_role", tt.userRole)
				handler.DeleteThread(c)
			})

			req := httptest.NewRequest("DELETE", "/threads/"+tt.threadID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
} 