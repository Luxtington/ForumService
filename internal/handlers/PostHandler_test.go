package handlers

import (
	"bytes"
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

func setupPostTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			switch err.Type {
			case gin.ErrorTypeBind:
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			case gin.ErrorTypePrivate:
				if err.Err.Error() == "пользователь не аутентифицирован" {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не аутентифицирован"})
					return
				}
				if err.Err.Error() == "нет прав для создания поста в этом треде" {
					c.JSON(http.StatusForbidden, gin.H{"error": "Нет прав для создания поста в этом треде"})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	})
	return router
}

func TestPostHandler_GetAllPosts(t *testing.T) {
	// Создаем мок сервиса
	mockPostService := &mocks.MockPostService{
		GetAllPostsFunc: func() ([]*models.Post, error) {
			return []*models.Post{
				{
					ID:       1,
					Content:  "Test post 1",
					AuthorID: 1,
				},
				{
					ID:       2,
					Content:  "Test post 2",
					AuthorID: 2,
				},
			}, nil
		},
	}

	// Создаем обработчик
	handler := NewPostHandler(mockPostService)

	// Настраиваем тестовый роутер
	router := setupPostTestRouter()
	router.GET("/posts", handler.GetAllPosts)

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем содержимое ответа
	var response []*models.Post
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "Test post 1", response[0].Content)
	assert.Equal(t, "Test post 2", response[1].Content)
}

func TestPostHandler_ShowCreateForm(t *testing.T) {
	// Создаем мок сервиса с реализацией всех необходимых методов
	mockPostService := &mocks.MockPostService{
		GetAllPostsFunc: func() ([]*models.Post, error) {
			return []*models.Post{}, nil
		},
		GetPostByIDFunc: func(id int) (*models.Post, error) {
			return nil, nil
		},
		GetPostWithCommentsFunc: func(postID int) (*models.Post, []models.Comment, error) {
			return nil, nil, nil
		},
		GetPostsWithCommentsByThreadIDFunc: func(threadID int) ([]models.Post, map[int][]models.Comment, error) {
			return nil, nil, nil
		},
		UpdatePostFunc: func(post *models.Post, postID int, userID int) error {
			return nil
		},
		DeletePostFunc: func(postID int, userID int) error {
			return nil
		},
		CreatePostFunc: func(post *models.Post) error {
			return nil
		},
		CreateCommentFunc: func(comment *models.Comment) error {
			return nil
		},
		GetCommentByIDFunc: func(id int) (*models.Comment, error) {
			return nil, nil
		},
		DeleteCommentFunc: func(id int) error {
			return nil
		},
		GetPostFunc: func(id int) (*models.Post, error) {
			return nil, nil
		},
		GetPostsByThreadIDFunc: func(threadID int) ([]*models.Post, error) {
			return nil, nil
		},
		GetThreadByIDFunc: func(id int) (*models.Thread, error) {
			return nil, nil
		},
	}

	// Создаем обработчик
	handler := NewPostHandler(mockPostService)

	// Настраиваем тестовый роутер
	router := setupPostTestRouter()
	
	// Настраиваем шаблонизатор
	router.LoadHTMLGlob("../../templates/*")
	router.Static("/static", "../../static")

	router.GET("/posts/create", handler.ShowCreateForm)

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts/create", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostHandler_CreatePost(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		userID         interface{}
		mockPost       *models.Post
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "успешное создание поста",
			requestBody: map[string]interface{}{
				"title":    "Test Post",
				"content":  "Test Content",
				"thread_id": 1,
			},
			userID: uint32(1),
			mockPost: &models.Post{
				ID:        1,
				Title:     "Test Post",
				Content:   "Test Content",
				AuthorID:  1,
				ThreadID:  1,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"id":          float64(1),
				"title":       "Test Post",
				"content":     "Test Content",
				"author_id":   float64(1),
				"thread_id":   float64(1),
				"created_at":  "0001-01-01T00:00:00Z",
				"updated_at":  "0001-01-01T00:00:00Z",
				"author_name": "",
				"can_edit":    false,
			},
		},
		{
			name: "пользователь не аутентифицирован",
			requestBody: map[string]interface{}{
				"title":    "Test Post",
				"content":  "Test Content",
				"thread_id": 1,
			},
			userID:         nil,
			mockPost:       nil,
			mockError:      nil,
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Пользователь не аутентифицирован",
			},
		},
		{
			name: "нет прав для создания поста",
			requestBody: map[string]interface{}{
				"title":    "Test Post",
				"content":  "Test Content",
				"thread_id": 1,
			},
			userID:         uint32(1),
			mockPost:       nil,
			mockError:      errors.New("нет прав для создания поста в этом треде"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Ошибка при создании поста: нет прав для создания поста в этом треде",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPostService := &mocks.MockPostService{
				CreatePostFunc: func(post *models.Post) error {
					if tt.mockError != nil {
						return tt.mockError
					}
					*post = *tt.mockPost
					return nil
				},
			}

			handler := NewPostHandler(mockPostService)
			router := setupPostTestRouter()
			router.POST("/posts", func(c *gin.Context) {
				if tt.userID != nil {
					c.Set("user_id", tt.userID)
				} else {
					c.Error(&gin.Error{
						Err:  errors.New("пользователь не аутентифицирован"),
						Type: gin.ErrorTypePrivate,
					})
					return
				}
				handler.CreatePost(c)
			})

			jsonBody, _ := json.Marshal(tt.requestBody)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonBody))
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

func TestPostHandler_GetPost(t *testing.T) {
	tests := []struct {
		name           string
		postID         string
		mockPost       *models.Post
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "успешное получение поста",
			postID: "1",
			mockPost: &models.Post{
				ID:        1,
				Content:   "Test post",
				AuthorID:  1,
				ThreadID:  1,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":          float64(1),
				"content":     "Test post",
				"author_id":   float64(1),
				"thread_id":   float64(1),
				"created_at":  "0001-01-01T00:00:00Z",
				"updated_at":  "0001-01-01T00:00:00Z",
				"author_name": "",
				"can_edit":    false,
				"title":       "",
			},
		},
		{
			name:           "неверный ID поста",
			postID:         "invalid",
			mockPost:       nil,
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Неверный ID поста: strconv.Atoi: parsing \"invalid\": invalid syntax",
			},
		},
		{
			name:           "пост не найден",
			postID:         "999",
			mockPost:       nil,
			mockError:      errors.New("post not found"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Пост не найден: post not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPostService := &mocks.MockPostService{
				GetPostFunc: func(id int) (*models.Post, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return tt.mockPost, nil
				},
			}

			handler := NewPostHandler(mockPostService)
			router := setupPostTestRouter()
			router.GET("/posts/:id", func(c *gin.Context) {
				c.Params = []gin.Param{{Key: "id", Value: tt.postID}}
				handler.GetPost(c)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/posts/"+tt.postID, nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}

func TestPostHandler_UpdatePost(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		userID         interface{}
		postID         string
		mockPost       *models.Post
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "успешное обновление поста",
			requestBody: map[string]interface{}{
				"content": "Updated post",
			},
			userID:   uint32(1),
			postID:   "1",
			mockPost: &models.Post{
				ID:        1,
				Content:   "Updated post",
				AuthorID:  1,
				ThreadID:  1,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":          float64(1),
				"content":     "Updated post",
				"author_id":   float64(1),
				"thread_id":   float64(1),
				"created_at":  "0001-01-01T00:00:00Z",
				"updated_at":  "0001-01-01T00:00:00Z",
				"author_name": "",
				"can_edit":    false,
				"title":       "",
			},
		},
		{
			name: "пользователь не аутентифицирован",
			requestBody: map[string]interface{}{
				"content": "Updated post",
			},
			userID:         nil,
			postID:         "1",
			mockPost:       nil,
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Пользователь не аутентифицирован",
			},
		},
		{
			name: "нет прав для редактирования",
			requestBody: map[string]interface{}{
				"content": "Updated post",
			},
			userID:         uint32(2),
			postID:         "1",
			mockPost: &models.Post{
				ID:        1,
				Content:   "Test post",
				AuthorID:  1,
				ThreadID:  1,
				CreatedAt: time.Time{},
			},
			mockError:      errors.New("нет прав для редактирования поста"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Пост не найден: нет прав для редактирования поста",
			},
		},
		{
			name: "ошибка сервиса",
			requestBody: map[string]interface{}{
				"content": "Updated post",
			},
			userID:         uint32(1),
			postID:         "1",
			mockPost:       nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Пост не найден: assert.AnError general error for testing",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPostService := &mocks.MockPostService{
				GetPostFunc: func(id int) (*models.Post, error) {
					return tt.mockPost, tt.mockError
				},
				UpdatePostFunc: func(post *models.Post, postID int, userID int) error {
					return tt.mockError
				},
			}

			handler := NewPostHandler(mockPostService)
			router := setupPostTestRouter()
			router.PUT("/posts/:id", func(c *gin.Context) {
				if tt.userID != nil {
					c.Set("user_id", tt.userID)
				} else {
					c.Error(&gin.Error{
						Err:  errors.New("Пользователь не аутентифицирован"),
						Type: gin.ErrorTypePrivate,
					})
					return
				}
				c.Params = []gin.Param{{Key: "id", Value: tt.postID}}
				handler.UpdatePost(c)
			})

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PUT", "/posts/1", bytes.NewBuffer(body))
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

func TestPostHandler_DeletePost(t *testing.T) {
	tests := []struct {
		name           string
		userID         interface{}
		postID         string
		mockPost       *models.Post
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "успешное удаление поста",
			userID: uint32(1),
			postID: "1",
			mockPost: &models.Post{
				ID:        1,
				Content:   "Test post",
				AuthorID:  1,
				ThreadID:  1,
				CreatedAt: time.Time{},
			},
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
			expectedBody:   nil,
		},
		{
			name:           "пользователь не аутентифицирован",
			userID:         nil,
			postID:         "1",
			mockPost:       nil,
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Пользователь не аутентифицирован",
			},
		},
		{
			name:   "нет прав для удаления",
			userID: uint32(2),
			postID: "1",
			mockPost: &models.Post{
				ID:        1,
				Content:   "Test post",
				AuthorID:  1,
				ThreadID:  1,
				CreatedAt: time.Time{},
			},
			mockError:      errors.New("нет прав для удаления поста"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Пост не найден: нет прав для удаления поста",
			},
		},
		{
			name:           "ошибка сервиса",
			userID:         uint32(1),
			postID:         "1",
			mockPost:       nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Пост не найден: assert.AnError general error for testing",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPostService := &mocks.MockPostService{
				GetPostFunc: func(id int) (*models.Post, error) {
					return tt.mockPost, tt.mockError
				},
				DeletePostFunc: func(postID int, userID int) error {
					return tt.mockError
				},
			}

			handler := NewPostHandler(mockPostService)
			router := setupPostTestRouter()
			router.DELETE("/posts/:id", func(c *gin.Context) {
				if tt.userID != nil {
					c.Set("user_id", tt.userID)
				} else {
					c.Error(&gin.Error{
						Err:  errors.New("Пользователь не аутентифицирован"),
						Type: gin.ErrorTypePrivate,
					})
					return
				}
				c.Params = []gin.Param{{Key: "id", Value: tt.postID}}
				handler.DeletePost(c)
			})

			req := httptest.NewRequest("DELETE", "/posts/1", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestPostHandler_GetPostComments(t *testing.T) {
	tests := []struct {
		name           string
		postID         string
		mockComments   []models.Comment
		mockPost       *models.Post
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "успешное получение комментариев",
			postID: "1",
			mockComments: []models.Comment{
				{
					ID:        1,
					Content:   "Test comment 1",
					AuthorID:  1,
					PostID:    1,
					CreatedAt: time.Time{},
				},
				{
					ID:        2,
					Content:   "Test comment 2",
					AuthorID:  2,
					PostID:    1,
					CreatedAt: time.Time{},
				},
			},
			mockPost: &models.Post{
				ID:        1,
				Content:   "Test post",
				AuthorID:  1,
				ThreadID:  1,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"comments": []interface{}{
					map[string]interface{}{
						"id":          float64(1),
						"content":     "Test comment 1",
						"author_id":   float64(1),
						"post_id":     float64(1),
						"created_at":  "0001-01-01T00:00:00Z",
						"author_name": "",
						"can_delete":  false,
					},
					map[string]interface{}{
						"id":          float64(2),
						"content":     "Test comment 2",
						"author_id":   float64(2),
						"post_id":     float64(1),
						"created_at":  "0001-01-01T00:00:00Z",
						"author_name": "",
						"can_delete":  false,
					},
				},
				"post": map[string]interface{}{
					"id":          float64(1),
					"content":     "Test post",
					"author_id":   float64(1),
					"thread_id":   float64(1),
					"created_at":  "0001-01-01T00:00:00Z",
					"updated_at":  "0001-01-01T00:00:00Z",
					"author_name": "",
					"can_edit":    false,
					"title":       "",
				},
			},
		},
		{
			name:           "неверный ID поста",
			postID:         "invalid",
			mockComments:   nil,
			mockPost:       nil,
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Неверный ID поста: strconv.Atoi: parsing \"invalid\": invalid syntax",
			},
		},
		{
			name:           "пост не найден",
			postID:         "1",
			mockComments:   nil,
			mockPost:       nil,
			mockError:      errors.New("post not found"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Пост не найден: post not found",
			},
		},
		{
			name:           "ошибка сервиса",
			postID:         "1",
			mockComments:   nil,
			mockPost:       nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Пост не найден: assert.AnError general error for testing",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPostService := &mocks.MockPostService{
				GetPostWithCommentsFunc: func(postID int) (*models.Post, []models.Comment, error) {
					return tt.mockPost, tt.mockComments, tt.mockError
				},
			}

			handler := NewPostHandler(mockPostService)
			router := setupPostTestRouter()
			router.GET("/posts/:id/comments", func(c *gin.Context) {
				c.Params = []gin.Param{{Key: "id", Value: tt.postID}}
				handler.GetPostComments(c)
			})

			req := httptest.NewRequest("GET", "/posts/1/comments", nil)
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