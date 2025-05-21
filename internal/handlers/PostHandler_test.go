package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"ForumService/internal/handlers/mocks"
	"ForumService/internal/models"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupPostTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.LoadHTMLGlob("../../templates/*")
	return router
}

func TestPostHandler_GetAllPosts(t *testing.T) {
	// Создаем мок сервиса
	mockPostService := &mocks.MockPostService{}

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
}

func TestPostHandler_ShowCreateForm(t *testing.T) {
	// Создаем мок сервиса
	mockPostService := &mocks.MockPostService{}

	// Создаем обработчик
	handler := NewPostHandler(mockPostService)

	// Настраиваем тестовый роутер
	router := setupPostTestRouter()
	router.GET("/posts/create", handler.ShowCreateForm)

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts/create", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostHandler_CreatePost_Success(t *testing.T) {
	// Создаем мок сервиса
	mockPostService := &mocks.MockPostService{
		GetThreadByIDFunc: func(id int) (*models.Thread, error) {
			return &models.Thread{
				ID:       1,
				AuthorID: 1,
			}, nil
		},
		CreatePostFunc: func(post *models.Post) error {
			return nil
		},
	}

	// Создаем обработчик
	handler := NewPostHandler(mockPostService)

	// Настраиваем тестовый роутер
	router := setupPostTestRouter()
	router.POST("/posts", func(c *gin.Context) {
		c.Set("user_id", uint32(1))
		c.Set("user_role", "user")
		handler.CreatePost(c)
	})

	// Создаем тестовый запрос
	requestBody := map[string]interface{}{
		"thread_id": 1,
		"content":   "Test post content",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusCreated, w.Code)
	
	// Проверяем содержимое ответа
	var response models.Post
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Logf("Response body: %s", w.Body.String())
		t.Fatal(err)
	}
	
	// Проверяем поля созданного поста
	assert.Equal(t, 1, response.ThreadID)
	assert.Equal(t, "Test post content", response.Content)
	assert.Equal(t, 1, response.AuthorID)
	assert.False(t, response.CanEdit)
}

func TestPostHandler_CreatePost_InvalidData(t *testing.T) {
	// Создаем мок сервиса
	mockPostService := &mocks.MockPostService{}

	// Создаем обработчик
	handler := NewPostHandler(mockPostService)

	// Настраиваем тестовый роутер
	router := setupPostTestRouter()
	router.POST("/posts", handler.CreatePost)

	// Создаем тестовый запрос с неверными данными
	requestBody := map[string]interface{}{
		"thread_id": "invalid",
		"content":   "",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostHandler_CreatePost_Unauthorized(t *testing.T) {
	// Создаем мок сервиса
	mockPostService := &mocks.MockPostService{}

	// Создаем обработчик
	handler := NewPostHandler(mockPostService)

	// Настраиваем тестовый роутер
	router := setupPostTestRouter()
	router.POST("/posts", handler.CreatePost)

	// Создаем тестовый запрос
	requestBody := map[string]interface{}{
		"thread_id": 1,
		"content":   "Test post content",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPostHandler_CreatePost_NoPermission(t *testing.T) {
	// Создаем мок сервиса
	mockPostService := &mocks.MockPostService{
		GetThreadByIDFunc: func(id int) (*models.Thread, error) {
			return &models.Thread{
				ID:       1,
				AuthorID: 2, // Другой автор
			}, nil
		},
	}

	// Создаем обработчик
	handler := NewPostHandler(mockPostService)

	// Настраиваем тестовый роутер
	router := setupPostTestRouter()
	router.POST("/posts", func(c *gin.Context) {
		c.Set("user_id", uint32(1))
		c.Set("user_role", "user")
		handler.CreatePost(c)
	})

	// Создаем тестовый запрос
	requestBody := map[string]interface{}{
		"thread_id": 1,
		"content":   "Test post content",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusForbidden, w.Code)
	
	// Проверяем содержимое ответа
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "нет прав для создания поста в этом треде", response["error"])
}

func TestPostHandler_GetPost_Success(t *testing.T) {
	// Создаем мок сервиса
	mockPostService := &mocks.MockPostService{
		GetPostByIDFunc: func(id int) (*models.Post, error) {
			return &models.Post{
				ID:       1,
				Content:  "Test post",
				AuthorID: 1,
			}, nil
		},
	}

	// Создаем обработчик
	handler := NewPostHandler(mockPostService)

	// Настраиваем тестовый роутер
	router := setupPostTestRouter()
	router.GET("/posts/:id", handler.GetPost)

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts/1", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Проверяем содержимое ответа
	var response models.Post
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, response.ID)
	assert.Equal(t, "Test post", response.Content)
	assert.Equal(t, 1, response.AuthorID)
}

func TestPostHandler_GetPost_InvalidID(t *testing.T) {
	// Создаем мок сервиса
	mockPostService := &mocks.MockPostService{}

	// Создаем обработчик
	handler := NewPostHandler(mockPostService)

	// Настраиваем тестовый роутер
	router := setupPostTestRouter()
	router.GET("/posts/:id", handler.GetPost)

	// Создаем тестовый запрос с неверным ID
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts/invalid", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostHandler_GetPost_NotFound(t *testing.T) {
	// Создаем мок сервиса
	mockPostService := &mocks.MockPostService{
		GetPostByIDFunc: func(id int) (*models.Post, error) {
			return nil, errors.New("post not found")
		},
	}

	// Создаем обработчик
	handler := NewPostHandler(mockPostService)

	// Настраиваем тестовый роутер
	router := setupPostTestRouter()
	router.GET("/posts/:id", handler.GetPost)

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts/999", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPostHandler_GetPostWithComments_Success(t *testing.T) {
	// Создаем мок сервиса
	mockPostService := &mocks.MockPostService{
		GetPostWithCommentsFunc: func(id int) (*models.Post, []models.Comment, error) {
			return &models.Post{
				ID:       1,
				Content:  "Test post",
				AuthorID: 1,
			}, []models.Comment{
				{
				ID:       1,
					Content:  "Test comment",
					PostID:   1,
					AuthorID: 1,
				},
			}, nil
		},
	}

	// Создаем обработчик
	handler := NewPostHandler(mockPostService)

	// Настраиваем тестовый роутер
	router := setupPostTestRouter()
	router.GET("/posts/:id/comments", handler.GetPostWithComments)

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts/1/comments", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostHandler_GetPostWithComments_InvalidID(t *testing.T) {
	// Создаем мок сервиса
	mockPostService := &mocks.MockPostService{}

	// Создаем обработчик
	handler := NewPostHandler(mockPostService)

	// Настраиваем тестовый роутер
	router := setupPostTestRouter()
	router.GET("/posts/:id/comments", handler.GetPostWithComments)

	// Создаем тестовый запрос с неверным ID
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts/invalid/comments", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostHandler_GetPostWithComments_NotFound(t *testing.T) {
	// Создаем мок сервиса
	mockPostService := &mocks.MockPostService{
		GetPostWithCommentsFunc: func(id int) (*models.Post, []models.Comment, error) {
			return nil, nil, errors.New("post not found")
		},
	}

	// Создаем обработчик
	handler := NewPostHandler(mockPostService)

	// Настраиваем тестовый роутер
	router := setupPostTestRouter()
	router.GET("/posts/:id/comments", handler.GetPostWithComments)

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts/999/comments", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPostHandler_ShowEditForm(t *testing.T) {
	// Создаем мок сервиса
	mockPostService := &mocks.MockPostService{
		GetPostByIDFunc: func(id int) (*models.Post, error) {
			return &models.Post{
				ID:       1,
				Content:  "Test post",
				AuthorID: 1,
			}, nil
		},
	}

	// Создаем обработчик
	handler := NewPostHandler(mockPostService)

	// Настраиваем тестовый роутер
	router := setupPostTestRouter()
	router.GET("/posts/:id/edit", func(c *gin.Context) {
		c.Set("user_id", uint32(1))
		c.Set("user_role", "user")
		handler.ShowEditForm(c)
	})

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts/1/edit", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostHandler_ShowEditForm_Unauthorized(t *testing.T) {
	// Создаем мок сервиса
	mockPostService := &mocks.MockPostService{
		GetPostByIDFunc: func(id int) (*models.Post, error) {
			return nil, errors.New("post not found")
		},
	}

	// Создаем обработчик
	handler := NewPostHandler(mockPostService)

	// Настраиваем тестовый роутер
	router := setupPostTestRouter()
	router.GET("/posts/:id/edit", handler.ShowEditForm)

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts/1/edit", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPostHandler_UpdatePost_Success(t *testing.T) {
	// Создаем мок сервиса
	mockPostService := &mocks.MockPostService{
				GetPostByIDFunc: func(id int) (*models.Post, error) {
			return &models.Post{
				ID:       1,
				Content:  "Old content",
				AuthorID: 1,
			}, nil
				},
				UpdatePostFunc: func(post *models.Post, postID int, userID int) error {
			return nil
		},
	}

	// Создаем обработчик
	handler := NewPostHandler(mockPostService)

	// Настраиваем тестовый роутер
	router := setupPostTestRouter()
			router.PUT("/posts/:id", func(c *gin.Context) {
		c.Set("user_id", uint32(1))
		c.Set("user_role", "user")
				handler.UpdatePost(c)
			})

	// Создаем тестовый запрос
	requestBody := map[string]interface{}{
		"content": "Updated content",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/posts/1", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostHandler_UpdatePost_Unauthorized(t *testing.T) {
	// Создаем мок сервиса
	mockPostService := &mocks.MockPostService{}

	// Создаем обработчик
	handler := NewPostHandler(mockPostService)

	// Настраиваем тестовый роутер
	router := setupPostTestRouter()
	router.PUT("/posts/:id", handler.UpdatePost)

	// Создаем тестовый запрос
	requestBody := map[string]interface{}{
		"content": "Updated content",
	}
	jsonBody, _ := json.Marshal(requestBody)

			w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/posts/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
			router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPostHandler_UpdatePost_NoPermission(t *testing.T) {
	// Создаем мок сервиса
	mockPostService := &mocks.MockPostService{
		GetPostByIDFunc: func(id int) (*models.Post, error) {
			return &models.Post{
				ID:       1,
				Content:  "Old content",
				AuthorID: 2, // Другой автор
			}, nil
		},
	}

	// Создаем обработчик
	handler := NewPostHandler(mockPostService)

	// Настраиваем тестовый роутер
	router := setupPostTestRouter()
	router.PUT("/posts/:id", func(c *gin.Context) {
		c.Set("user_id", uint32(1))
		c.Set("user_role", "user")
		handler.UpdatePost(c)
	})

	// Создаем тестовый запрос
	requestBody := map[string]interface{}{
		"content": "Updated content",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/posts/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPostHandler_DeletePost_Success(t *testing.T) {
	// Создаем мок сервиса
	mockPostService := &mocks.MockPostService{
		GetPostByIDFunc: func(id int) (*models.Post, error) {
			return &models.Post{
				ID:       1,
				Content:  "Test post",
				AuthorID: 1,
			}, nil
		},
		DeletePostFunc: func(postID int, userID int) error {
			return nil
		},
	}

	// Создаем обработчик
	handler := NewPostHandler(mockPostService)

	// Настраиваем тестовый роутер
	router := setupPostTestRouter()
	router.DELETE("/posts/:id", func(c *gin.Context) {
		c.Set("user_id", uint32(1))
		c.Set("user_role", "user")
		handler.DeletePost(c)
	})

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/posts/1", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestPostHandler_DeletePost_Unauthorized(t *testing.T) {
	// Создаем мок сервиса
	mockPostService := &mocks.MockPostService{}

	// Создаем обработчик
	handler := NewPostHandler(mockPostService)

	// Настраиваем тестовый роутер
	router := setupPostTestRouter()
	router.DELETE("/posts/:id", handler.DeletePost)

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/posts/1", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPostHandler_DeletePost_NoPermission(t *testing.T) {
	// Создаем мок сервиса
	mockPostService := &mocks.MockPostService{
		GetPostByIDFunc: func(id int) (*models.Post, error) {
			return &models.Post{
				ID:       1,
				Content:  "Test post",
				AuthorID: 2, // Другой автор
			}, nil
		},
	}

	// Создаем обработчик
	handler := NewPostHandler(mockPostService)

	// Настраиваем тестовый роутер
	router := setupPostTestRouter()
	router.DELETE("/posts/:id", func(c *gin.Context) {
		c.Set("user_id", uint32(1))
		c.Set("user_role", "user")
		handler.DeletePost(c)
	})

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/posts/1", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPostHandler_DeletePost_AdminSuccess(t *testing.T) {
	// Создаем мок сервиса
	mockPostService := &mocks.MockPostService{
				GetPostByIDFunc: func(id int) (*models.Post, error) {
			return &models.Post{
				ID:       1,
				Content:  "Test post",
				AuthorID: 2, // Другой автор
			}, nil
				},
				DeletePostFunc: func(postID int, userID int) error {
			return nil
		},
	}

	// Создаем обработчик
	handler := NewPostHandler(mockPostService)

	// Настраиваем тестовый роутер
	router := setupPostTestRouter()
			router.DELETE("/posts/:id", func(c *gin.Context) {
		c.Set("user_id", uint32(1))
		c.Set("user_role", "admin")
				handler.DeletePost(c)
			})

	// Создаем тестовый запрос
			w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/posts/1", nil)

	// Выполняем запрос
			router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestPostHandler_ListPosts_Success(t *testing.T) {
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
    router.GET("/posts/list", handler.ListPosts)

    // Создаем тестовый запрос
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/posts/list", nil)

    // Выполняем запрос
    router.ServeHTTP(w, req)

    // Проверяем результат
    assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostHandler_ListPosts_Error(t *testing.T) {
    // Создаем мок сервиса
    mockPostService := &mocks.MockPostService{
        GetAllPostsFunc: func() ([]*models.Post, error) {
            return nil, errors.New("database error")
        },
    }

    // Создаем обработчик
    handler := NewPostHandler(mockPostService)

    // Настраиваем тестовый роутер
    router := setupPostTestRouter()
    router.GET("/posts/list", handler.ListPosts)

    // Создаем тестовый запрос
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/posts/list", nil)

    // Выполняем запрос
    router.ServeHTTP(w, req)

    // Проверяем результат
    assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPostHandler_ShowPost_Success(t *testing.T) {
    // Создаем мок сервиса
    mockPostService := &mocks.MockPostService{
        GetPostWithCommentsFunc: func(id int) (*models.Post, []models.Comment, error) {
            return &models.Post{
                    ID:       1,
                    Content:  "Test post",
                    AuthorID: 1,
                }, []models.Comment{
                    {
                        ID:       1,
                        Content:  "Test comment",
                        AuthorID: 1,
                        PostID:   1,
                    },
                }, nil
        },
    }

    // Создаем обработчик
    handler := NewPostHandler(mockPostService)

    // Настраиваем тестовый роутер
    router := setupPostTestRouter()
    router.GET("/posts/view/:id", func(c *gin.Context) {
        c.Set("user", &models.User{ID: 1})
        handler.ShowPost(c)
    })

    // Создаем тестовый запрос
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/posts/view/1", nil)

    // Выполняем запрос
    router.ServeHTTP(w, req)

    // Проверяем результат
    assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostHandler_ShowPost_InvalidID(t *testing.T) {
    // Создаем мок сервиса
    mockPostService := &mocks.MockPostService{}

    // Создаем обработчик
    handler := NewPostHandler(mockPostService)

    // Настраиваем тестовый роутер
    router := setupPostTestRouter()
    router.GET("/posts/view/:id", handler.ShowPost)

    // Создаем тестовый запрос с неверным ID
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/posts/view/invalid", nil)

    // Выполняем запрос
    router.ServeHTTP(w, req)

    // Проверяем результат
    assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostHandler_ShowPost_NotFound(t *testing.T) {
    // Создаем мок сервиса
    mockPostService := &mocks.MockPostService{
        GetPostWithCommentsFunc: func(id int) (*models.Post, []models.Comment, error) {
            return nil, nil, errors.New("post not found")
        },
    }

    // Создаем обработчик
    handler := NewPostHandler(mockPostService)

    // Настраиваем тестовый роутер
    router := setupPostTestRouter()
    router.GET("/posts/view/:id", handler.ShowPost)

    // Создаем тестовый запрос
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/posts/view/999", nil)

    // Выполняем запрос
    router.ServeHTTP(w, req)

    // Проверяем результат
    assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPostHandler_CreateComment_Success(t *testing.T) {
    // Создаем мок сервиса
    mockPostService := &mocks.MockPostService{
        CreateCommentFunc: func(comment *models.Comment) error {
            return nil
        },
    }

    // Создаем обработчик
    handler := NewPostHandler(mockPostService)

    // Настраиваем тестовый роутер
    router := setupPostTestRouter()
    router.POST("/posts/:id/comments", func(c *gin.Context) {
        c.Set("userID", 1)
        handler.CreateComment(c)
    })

    // Создаем тестовый запрос
    form := url.Values{}
    form.Add("content", "Test comment")
    
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/posts/1/comments", strings.NewReader(form.Encode()))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    // Выполняем запрос
    router.ServeHTTP(w, req)

    // Проверяем результат
    assert.Equal(t, http.StatusFound, w.Code)
    assert.Equal(t, "/posts/1", w.Header().Get("Location"))
}

func TestPostHandler_CreateComment_InvalidPostID(t *testing.T) {
    // Создаем мок сервиса
    mockPostService := &mocks.MockPostService{}

    // Создаем обработчик
    handler := NewPostHandler(mockPostService)

    // Настраиваем тестовый роутер
    router := setupPostTestRouter()
    router.POST("/posts/:id/comments", func(c *gin.Context) {
        c.Set("userID", 1)
        handler.CreateComment(c)
    })

    // Создаем тестовый запрос с неверным ID поста
    form := url.Values{}
    form.Add("content", "Test comment")
    
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/posts/invalid/comments", strings.NewReader(form.Encode()))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    // Выполняем запрос
    router.ServeHTTP(w, req)

    // Проверяем результат
    assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostHandler_CreateComment_EmptyContent(t *testing.T) {
    // Создаем мок сервиса
    mockPostService := &mocks.MockPostService{
        CreateCommentFunc: func(comment *models.Comment) error {
            return errors.New("содержимое комментария не может быть пустым")
        },
    }

    // Создаем обработчик
    handler := NewPostHandler(mockPostService)

    // Настраиваем тестовый роутер
    router := setupPostTestRouter()
    router.POST("/posts/:id/comments", func(c *gin.Context) {
        c.Set("userID", 1)
        handler.CreateComment(c)
    })

    // Создаем тестовый запрос с пустым содержимым
    form := url.Values{}
    form.Add("content", "")
    
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/posts/1/comments", strings.NewReader(form.Encode()))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    // Выполняем запрос
    router.ServeHTTP(w, req)

    // Проверяем результат
    assert.Equal(t, http.StatusInternalServerError, w.Code)
    
    // Проверяем содержимое ответа
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    if err == nil {
        assert.NotNil(t, response["error"])
    }
}

func TestPostHandler_DeleteComment_Success(t *testing.T) {
    // Создаем мок сервиса
    mockPostService := &mocks.MockPostService{
        GetCommentByIDFunc: func(id int) (*models.Comment, error) {
            return &models.Comment{
                ID:       1,
                Content:  "Test comment",
                AuthorID: 1,
                PostID:   1,
            }, nil
        },
        DeleteCommentFunc: func(id int) error {
            return nil
        },
    }

    // Создаем обработчик
    handler := NewPostHandler(mockPostService)

    // Настраиваем тестовый роутер
    router := setupPostTestRouter()
    router.POST("/posts/:postId/comments/:commentId/delete", func(c *gin.Context) {
        c.Set("userID", 1)
        handler.DeleteComment(c)
    })

    // Создаем тестовый запрос
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/posts/1/comments/1/delete", nil)

    // Выполняем запрос
    router.ServeHTTP(w, req)

    // Проверяем результат
    assert.Equal(t, http.StatusFound, w.Code)
    assert.Equal(t, "/posts/1", w.Header().Get("Location"))
}

func TestPostHandler_DeleteComment_NoPermission(t *testing.T) {
    // Создаем мок сервиса
    mockPostService := &mocks.MockPostService{
        GetCommentByIDFunc: func(id int) (*models.Comment, error) {
            return &models.Comment{
                ID:       1,
                Content:  "Test comment",
                AuthorID: 2, // Другой автор
                PostID:   1,
            }, nil
        },
    }

    // Создаем обработчик
    handler := NewPostHandler(mockPostService)

    // Настраиваем тестовый роутер
    router := setupPostTestRouter()
    router.POST("/posts/:postId/comments/:commentId/delete", func(c *gin.Context) {
        c.Set("userID", 1)
        handler.DeleteComment(c)
    })

    // Создаем тестовый запрос
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/posts/1/comments/1/delete", nil)

    // Выполняем запрос
    router.ServeHTTP(w, req)

    // Проверяем результат
    assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPostHandler_GetPostComments_Success(t *testing.T) {
    // Создаем мок сервиса
    mockPostService := &mocks.MockPostService{
        GetPostWithCommentsFunc: func(id int) (*models.Post, []models.Comment, error) {
            return nil, []models.Comment{
                {
                    ID:       1,
                    Content:  "Test comment 1",
                    AuthorID: 1,
                    PostID:   1,
                },
                {
                    ID:       2,
                    Content:  "Test comment 2",
                    AuthorID: 2,
                    PostID:   1,
                },
            }, nil
        },
    }

    // Создаем обработчик
    handler := NewPostHandler(mockPostService)

    // Настраиваем тестовый роутер
    router := setupPostTestRouter()
    router.GET("/posts/:id/comments/list", handler.GetPostComments)

    // Создаем тестовый запрос
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/posts/1/comments/list", nil)

    // Выполняем запрос
    router.ServeHTTP(w, req)

    // Проверяем результат
    assert.Equal(t, http.StatusOK, w.Code)
    
    var comments []models.Comment
    err := json.Unmarshal(w.Body.Bytes(), &comments)
    assert.NoError(t, err)
    assert.Equal(t, 2, len(comments))
}

func TestPostHandler_GetPostComments_InvalidID(t *testing.T) {
    // Создаем мок сервиса
    mockPostService := &mocks.MockPostService{}

    // Создаем обработчик
    handler := NewPostHandler(mockPostService)

    // Настраиваем тестовый роутер
    router := setupPostTestRouter()
    router.GET("/posts/:id/comments/list", handler.GetPostComments)

    // Создаем тестовый запрос с неверным ID
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/posts/invalid/comments/list", nil)

    // Выполняем запрос
    router.ServeHTTP(w, req)

    // Проверяем результат
    assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostHandler_GetPostComments_NotFound(t *testing.T) {
    // Создаем мок сервиса
    mockPostService := &mocks.MockPostService{
        GetPostWithCommentsFunc: func(id int) (*models.Post, []models.Comment, error) {
            return nil, nil, errors.New("post not found")
        },
    }

    // Создаем обработчик
    handler := NewPostHandler(mockPostService)

    // Настраиваем тестовый роутер
    router := setupPostTestRouter()
    router.GET("/posts/:id/comments/list", handler.GetPostComments)

    // Создаем тестовый запрос
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/posts/999/comments/list", nil)

    // Выполняем запрос
    router.ServeHTTP(w, req)

    // Проверяем результат
    assert.Equal(t, http.StatusNotFound, w.Code)
}