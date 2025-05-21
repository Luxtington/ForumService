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

func TestCreatePost(t *testing.T) {
	tests := []struct {
		name           string
		threadID       string
		requestBody    map[string]interface{}
		mockResponse   *models.Post
		mockError      error
		expectedStatus int
		userID         uint
		userRole       string
	}{
		{
			name:     "успешное создание поста",
			threadID: "1",
			requestBody: map[string]interface{}{
				"thread_id": 1,
				"content":   "Test Post Content",
			},
			mockResponse: &models.Post{
				ID:       1,
				Content:  "Test Post Content",
				ThreadID: 1,
				AuthorID: 1,
			},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			userID:         1,
			userRole:       "admin",
		},
		{
			name:     "неверный формат данных",
			threadID: "1",
			requestBody: map[string]interface{}{
				"invalid": "data",
			},
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			userID:         1,
			userRole:       "admin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.MockPostService{
				CreatePostFunc: func(post *models.Post) error {
					return tt.mockError
				},
				GetThreadByIDFunc: func(id int) (*models.Thread, error) {
					return &models.Thread{ID: id}, nil
				},
			}

			handler := NewPostHandler(mockService)
			router := setupTestRouter()
			router.POST("/threads/:threadID/posts", func(c *gin.Context) {
				c.Set("user_id", tt.userID)
				c.Set("user_role", tt.userRole)
				c.Params = []gin.Param{{Key: "threadID", Value: tt.threadID}}
				handler.CreatePost(c)
			})

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/threads/"+tt.threadID+"/posts", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestGetPost(t *testing.T) {
	tests := []struct {
		name           string
		postID         string
		mockResponse   *models.Post
		mockError      error
		expectedStatus int
	}{
		{
			name:   "успешное получение поста",
			postID: "1",
			mockResponse: &models.Post{
				ID:       1,
				Title:    "Test Post",
				Content:  "Test Content",
				ThreadID: 1,
				AuthorID: 1,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "пост не найден",
			postID:         "999",
			mockResponse:   nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.MockPostService{
				GetPostByIDFunc: func(id int) (*models.Post, error) {
					return tt.mockResponse, tt.mockError
				},
			}

			handler := NewPostHandler(mockService)
			router := setupTestRouter()
			router.GET("/posts/:id", func(c *gin.Context) {
				c.Params = []gin.Param{{Key: "id", Value: tt.postID}}
				handler.GetPost(c)
			})

			req := httptest.NewRequest("GET", "/posts/"+tt.postID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestUpdatePost(t *testing.T) {
	tests := []struct {
		name           string
		postID         string
		requestBody    map[string]interface{}
		mockResponse   *models.Post
		mockError      error
		expectedStatus int
		userID         uint
		userRole       string
	}{
		{
			name:   "успешное обновление поста автором",
			postID: "1",
			requestBody: map[string]interface{}{
				"title":   "Updated Post Title",
				"content": "Updated Post Content",
			},
			mockResponse: &models.Post{
				ID:       1,
				Title:    "Updated Post Title",
				Content:  "Updated Post Content",
				ThreadID: 1,
				AuthorID: 1,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			userID:         1,
			userRole:       "user",
		},
		{
			name:   "успешное обновление поста админом",
			postID: "1",
			requestBody: map[string]interface{}{
				"title":   "Updated Post Title",
				"content": "Updated Post Content",
			},
			mockResponse: &models.Post{
				ID:       1,
				Title:    "Updated Post Title",
				Content:  "Updated Post Content",
				ThreadID: 1,
				AuthorID: 2,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			userID:         1,
			userRole:       "admin",
		},
		{
			name:   "отказ в доступе",
			postID: "1",
			requestBody: map[string]interface{}{
				"title":   "Updated Post Title",
				"content": "Updated Post Content",
			},
			mockResponse: &models.Post{
				ID:       1,
				Title:    "Original Post Title",
				Content:  "Original Post Content",
				ThreadID: 1,
				AuthorID: 2,
			},
			mockError:      nil,
			expectedStatus: http.StatusForbidden,
			userID:         1,
			userRole:       "user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.MockPostService{
				GetPostByIDFunc: func(id int) (*models.Post, error) {
					return tt.mockResponse, nil
				},
				UpdatePostFunc: func(post *models.Post, postID int, userID int) error {
					return tt.mockError
				},
			}

			handler := NewPostHandler(mockService)
			router := setupTestRouter()
			router.PUT("/posts/:id", func(c *gin.Context) {
				c.Set("user_id", tt.userID)
				c.Set("user_role", tt.userRole)
				c.Params = []gin.Param{{Key: "id", Value: tt.postID}}
				handler.UpdatePost(c)
			})

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PUT", "/posts/"+tt.postID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestDeletePost(t *testing.T) {
	tests := []struct {
		name           string
		postID         string
		mockResponse   *models.Post
		mockError      error
		expectedStatus int
		userID         uint
		userRole       string
	}{
		{
			name:   "успешное удаление поста автором",
			postID: "1",
			mockResponse: &models.Post{
				ID:       1,
				Title:    "Test Post",
				Content:  "Test Content",
				ThreadID: 1,
				AuthorID: 1,
			},
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
			userID:         1,
			userRole:       "user",
		},
		{
			name:   "успешное удаление поста админом",
			postID: "1",
			mockResponse: &models.Post{
				ID:       1,
				Title:    "Test Post",
				Content:  "Test Content",
				ThreadID: 1,
				AuthorID: 2,
			},
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
			userID:         1,
			userRole:       "admin",
		},
		{
			name:   "отказ в доступе",
			postID: "1",
			mockResponse: &models.Post{
				ID:       1,
				Title:    "Test Post",
				Content:  "Test Content",
				ThreadID: 1,
				AuthorID: 2,
			},
			mockError:      nil,
			expectedStatus: http.StatusForbidden,
			userID:         1,
			userRole:       "user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.MockPostService{
				GetPostByIDFunc: func(id int) (*models.Post, error) {
					return tt.mockResponse, nil
				},
				DeletePostFunc: func(postID int, userID int) error {
					return tt.mockError
				},
			}

			handler := NewPostHandler(mockService)
			router := setupTestRouter()
			router.DELETE("/posts/:id", func(c *gin.Context) {
				c.Set("user_id", tt.userID)
				c.Set("user_role", tt.userRole)
				c.Params = []gin.Param{{Key: "id", Value: tt.postID}}
				handler.DeletePost(c)
			})

			req := httptest.NewRequest("DELETE", "/posts/"+tt.postID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
} 