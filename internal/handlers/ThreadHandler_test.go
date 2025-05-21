package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"ForumService/internal/handlers/mocks"
	"ForumService/internal/models"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupThreadTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.LoadHTMLGlob("../../templates/*")
	return router
}

func TestThreadHandler_CreateThread_Success(t *testing.T) {
	// Создаем мок сервиса
	mockThreadService := &mocks.MockThreadService{
		CreateThreadFunc: func(title string, authorID int) (*models.Thread, error) {
			return &models.Thread{
				ID:       1,
				Title:    title,
				AuthorID: authorID,
			}, nil
		},
	}

	// Создаем обработчик
	handler := NewThreadHandler(mockThreadService)

	// Настраиваем тестовый роутер
	router := setupThreadTestRouter()
	router.POST("/threads", func(c *gin.Context) {
		c.Set("user_id", uint32(1))
		c.Set("user_role", "user")
		handler.CreateThread(c)
	})

	// Создаем тестовый запрос
	requestBody := map[string]interface{}{
		"title": "Test Thread",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/threads", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestThreadHandler_CreateThread_Unauthorized(t *testing.T) {
	// Создаем мок сервиса
	mockThreadService := &mocks.MockThreadService{}

	// Создаем обработчик
	handler := NewThreadHandler(mockThreadService)

	// Настраиваем тестовый роутер
	router := setupThreadTestRouter()
	router.POST("/threads", handler.CreateThread)

	// Создаем тестовый запрос
	requestBody := map[string]interface{}{
		"title": "Test Thread",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/threads", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestThreadHandler_CreateThread_InvalidData(t *testing.T) {
	// Создаем мок сервиса
	mockThreadService := &mocks.MockThreadService{}

	// Создаем обработчик
	handler := NewThreadHandler(mockThreadService)

	// Настраиваем тестовый роутер
	router := setupThreadTestRouter()
	router.POST("/threads", func(c *gin.Context) {
		c.Set("user_id", uint32(1))
		c.Set("user_role", "user")
		handler.CreateThread(c)
	})

	// Создаем тестовый запрос с неверными данными
	requestBody := map[string]interface{}{
		"title": "",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/threads", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestThreadHandler_GetThreadWithPosts_Success(t *testing.T) {
	// Создаем мок сервиса
	mockThreadService := &mocks.MockThreadService{
		GetThreadWithPostsFunc: func(id int) (*models.Thread, []*models.Post, error) {
			return &models.Thread{
				ID:       1,
				Title:    "Test Thread",
				AuthorID: 1,
			}, []*models.Post{
				{
					ID:       1,
					ThreadID: 1,
					Content:  "Test Post",
				},
			}, nil
		},
	}

	// Создаем обработчик
	handler := NewThreadHandler(mockThreadService)

	// Настраиваем тестовый роутер
	router := setupThreadTestRouter()
	router.GET("/threads/:id", handler.GetThreadWithPosts)

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/threads/1", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestThreadHandler_GetThreadWithPosts_InvalidID(t *testing.T) {
	// Создаем мок сервиса
	mockThreadService := &mocks.MockThreadService{}

	// Создаем обработчик
	handler := NewThreadHandler(mockThreadService)

	// Настраиваем тестовый роутер
	router := setupThreadTestRouter()
	router.GET("/threads/:id", handler.GetThreadWithPosts)

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/threads/invalid", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestThreadHandler_GetThreadWithPosts_NotFound(t *testing.T) {
	// Создаем мок сервиса
	mockThreadService := &mocks.MockThreadService{
		GetThreadWithPostsFunc: func(id int) (*models.Thread, []*models.Post, error) {
			return nil, nil, errors.New("thread not found")
		},
	}

	// Создаем обработчик
	handler := NewThreadHandler(mockThreadService)

	// Настраиваем тестовый роутер
	router := setupThreadTestRouter()
	router.GET("/threads/:id", handler.GetThreadWithPosts)

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/threads/1", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestThreadHandler_DeleteThread_Success(t *testing.T) {
	// Создаем мок сервиса
	mockThreadService := &mocks.MockThreadService{
		GetThreadWithPostsFunc: func(id int) (*models.Thread, []*models.Post, error) {
			return &models.Thread{
				ID:       1,
				Title:    "Test Thread",
				AuthorID: 1,
			}, nil, nil
		},
		DeleteThreadFunc: func(id int, userID int) error {
			return nil
		},
	}

	// Создаем обработчик
	handler := NewThreadHandler(mockThreadService)

	// Настраиваем тестовый роутер
	router := setupThreadTestRouter()
	router.DELETE("/threads/:id", func(c *gin.Context) {
		c.Set("user_id", uint32(1))
		c.Set("user_role", "user")
		handler.DeleteThread(c)
	})

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/threads/1", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestThreadHandler_DeleteThread_Unauthorized(t *testing.T) {
	// Создаем мок сервиса
	mockThreadService := &mocks.MockThreadService{}

	// Создаем обработчик
	handler := NewThreadHandler(mockThreadService)

	// Настраиваем тестовый роутер
	router := setupThreadTestRouter()
	router.DELETE("/threads/:id", handler.DeleteThread)

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/threads/1", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestThreadHandler_DeleteThread_NoPermission(t *testing.T) {
	// Создаем мок сервиса
	mockThreadService := &mocks.MockThreadService{
		GetThreadWithPostsFunc: func(id int) (*models.Thread, []*models.Post, error) {
			return &models.Thread{
				ID:       1,
				Title:    "Test Thread",
				AuthorID: 2, // Другой автор
			}, nil, nil
		},
	}

	// Создаем обработчик
	handler := NewThreadHandler(mockThreadService)

	// Настраиваем тестовый роутер
	router := setupThreadTestRouter()
	router.DELETE("/threads/:id", func(c *gin.Context) {
		c.Set("user_id", uint32(1))
		c.Set("user_role", "user")
		handler.DeleteThread(c)
	})

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/threads/1", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestThreadHandler_DeleteThread_AdminSuccess(t *testing.T) {
	// Создаем мок сервиса
	mockThreadService := &mocks.MockThreadService{
		GetThreadWithPostsFunc: func(id int) (*models.Thread, []*models.Post, error) {
			return &models.Thread{
				ID:       1,
				Title:    "Test Thread",
				AuthorID: 2, // Другой автор
			}, nil, nil
		},
		DeleteThreadFunc: func(id int, userID int) error {
			return nil
		},
	}

	// Создаем обработчик
	handler := NewThreadHandler(mockThreadService)

	// Настраиваем тестовый роутер
	router := setupThreadTestRouter()
	router.DELETE("/threads/:id", func(c *gin.Context) {
		c.Set("user_id", uint32(1))
		c.Set("user_role", "admin")
		handler.DeleteThread(c)
	})

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/threads/1", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestThreadHandler_UpdateThread_Success(t *testing.T) {
	// Создаем мок сервиса
	mockThreadService := &mocks.MockThreadService{
		GetThreadWithPostsFunc: func(id int) (*models.Thread, []*models.Post, error) {
			return &models.Thread{
				ID:       1,
				Title:    "Old Title",
				AuthorID: 1,
			}, nil, nil
		},
		UpdateThreadFunc: func(thread *models.Thread, userID int) error {
			return nil
		},
	}

	// Создаем обработчик
	handler := NewThreadHandler(mockThreadService)

	// Настраиваем тестовый роутер
	router := setupThreadTestRouter()
	router.PUT("/threads/:id", func(c *gin.Context) {
		c.Set("user_id", uint32(1))
		c.Set("user_role", "user")
		handler.UpdateThread(c)
	})

	// Создаем тестовый запрос
	requestBody := map[string]interface{}{
		"title": "New Title",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/threads/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestThreadHandler_UpdateThread_Unauthorized(t *testing.T) {
	// Создаем мок сервиса
	mockThreadService := &mocks.MockThreadService{}

	// Создаем обработчик
	handler := NewThreadHandler(mockThreadService)

	// Настраиваем тестовый роутер
	router := setupThreadTestRouter()
	router.PUT("/threads/:id", handler.UpdateThread)

	// Создаем тестовый запрос
	requestBody := map[string]interface{}{
		"title": "New Title",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/threads/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestThreadHandler_UpdateThread_NoPermission(t *testing.T) {
	// Создаем мок сервиса
	mockThreadService := &mocks.MockThreadService{
		GetThreadWithPostsFunc: func(id int) (*models.Thread, []*models.Post, error) {
			return &models.Thread{
				ID:       1,
				Title:    "Old Title",
				AuthorID: 2, // Другой автор
			}, nil, nil
		},
	}

	// Создаем обработчик
	handler := NewThreadHandler(mockThreadService)

	// Настраиваем тестовый роутер
	router := setupThreadTestRouter()
	router.PUT("/threads/:id", func(c *gin.Context) {
		c.Set("user_id", uint32(1))
		c.Set("user_role", "user")
		handler.UpdateThread(c)
	})

	// Создаем тестовый запрос
	requestBody := map[string]interface{}{
		"title": "New Title",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/threads/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestThreadHandler_GetAllThreads_Success(t *testing.T) {
	// Создаем мок сервиса
	mockThreadService := &mocks.MockThreadService{
		GetAllThreadsFunc: func() ([]*models.Thread, error) {
			return []*models.Thread{
				{
					ID:       1,
					Title:    "Thread 1",
					AuthorID: 1,
				},
				{
					ID:       2,
					Title:    "Thread 2",
					AuthorID: 2,
				},
			}, nil
		},
		GetUserByIDFunc: func(id int) (*models.User, error) {
			return &models.User{
				ID:       id,
				Username: "user" + strconv.Itoa(id),
			}, nil
		},
	}

	// Создаем обработчик
	handler := NewThreadHandler(mockThreadService)

	// Настраиваем тестовый роутер
	router := setupThreadTestRouter()
	router.GET("/threads", handler.GetAllThreads)

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/threads", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestThreadHandler_GetAllThreads_Error(t *testing.T) {
	// Создаем мок сервиса
	mockThreadService := &mocks.MockThreadService{
		GetAllThreadsFunc: func() ([]*models.Thread, error) {
			return nil, errors.New("database error")
		},
	}

	// Создаем обработчик
	handler := NewThreadHandler(mockThreadService)

	// Настраиваем тестовый роутер
	router := setupThreadTestRouter()
	router.GET("/threads", handler.GetAllThreads)

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/threads", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestThreadHandler_GetThreadPosts_Success(t *testing.T) {
	// Создаем мок сервиса
	mockThreadService := &mocks.MockThreadService{
		GetPostsByThreadIDFunc: func(threadID int) ([]*models.Post, error) {
			return []*models.Post{
				{
					ID:       1,
					ThreadID: threadID,
					Content:  "Post 1",
				},
				{
					ID:       2,
					ThreadID: threadID,
					Content:  "Post 2",
				},
			}, nil
		},
	}

	// Создаем обработчик
	handler := NewThreadHandler(mockThreadService)

	// Настраиваем тестовый роутер
	router := setupThreadTestRouter()
	router.GET("/threads/:id/posts", handler.GetThreadPosts)

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/threads/1/posts", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestThreadHandler_GetThreadPosts_InvalidID(t *testing.T) {
	// Создаем мок сервиса
	mockThreadService := &mocks.MockThreadService{}

	// Создаем обработчик
	handler := NewThreadHandler(mockThreadService)

	// Настраиваем тестовый роутер
	router := setupThreadTestRouter()
	router.GET("/threads/:id/posts", handler.GetThreadPosts)

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/threads/invalid/posts", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestThreadHandler_GetThreadPosts_NotFound(t *testing.T) {
	// Создаем мок сервиса
	mockThreadService := &mocks.MockThreadService{
		GetPostsByThreadIDFunc: func(threadID int) ([]*models.Post, error) {
			return nil, errors.New("thread not found")
		},
	}

	// Создаем обработчик
	handler := NewThreadHandler(mockThreadService)

	// Настраиваем тестовый роутер
	router := setupThreadTestRouter()
	router.GET("/threads/:id/posts", handler.GetThreadPosts)

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/threads/1/posts", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestThreadHandler_UpdateThread_InvalidID(t *testing.T) {
	// Создаем мок сервиса
	mockThreadService := &mocks.MockThreadService{}
	// Создаем обработчик
	handler := NewThreadHandler(mockThreadService)
	// Настраиваем тестовый роутер
	router := setupThreadTestRouter()
	router.PUT("/threads/:id", handler.UpdateThread)
	// Создаем тестовый запрос
	requestBody := map[string]interface{}{
		"title": "New Title",
	}
	jsonBody, _ := json.Marshal(requestBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/threads/invalid", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	// Выполняем запрос
	router.ServeHTTP(w, req)
	// Проверяем результат
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestThreadHandler_GetThreadWithPosts_ServiceError(t *testing.T) {
	// Создаем мок сервиса
	mockThreadService := &mocks.MockThreadService{
		GetThreadWithPostsFunc: func(id int) (*models.Thread, []*models.Post, error) {
			return nil, nil, errors.New("thread not found")
		},
	}

	// Создаем обработчик
	handler := NewThreadHandler(mockThreadService)

	// Настраиваем тестовый роутер
	router := setupThreadTestRouter()
	router.GET("/threads/:id", handler.GetThreadWithPosts)

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/threads/1", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusNotFound, w.Code)
	
	// Проверяем содержимое ответа
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Logf("Response body: %s", w.Body.String())
		t.Fatal(err)
	}
	assert.Contains(t, response, "error")
	assert.Equal(t, "thread not found", response["error"])
}

func TestThreadHandler_GetAllThreads_EmptyList(t *testing.T) {
	// Создаем мок сервиса
	mockThreadService := &mocks.MockThreadService{
		GetAllThreadsFunc: func() ([]*models.Thread, error) {
			return []*models.Thread{}, nil
		},
	}
	// Создаем обработчик
	handler := NewThreadHandler(mockThreadService)
	// Настраиваем тестовый роутер
	router := setupThreadTestRouter()
	router.GET("/threads", handler.GetAllThreads)
	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/threads", nil)
	// Выполняем запрос
	router.ServeHTTP(w, req)
	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, "[]", w.Body.String())
}

func TestThreadHandler_GetAllThreads_UserNotFound(t *testing.T) {
	// Создаем мок сервиса
	mockThreadService := &mocks.MockThreadService{
		GetAllThreadsFunc: func() ([]*models.Thread, error) {
			return []*models.Thread{
				{
					ID:       1,
					Title:    "Thread 1",
					AuthorID: 1,
				},
			}, nil
		},
		GetUserByIDFunc: func(id int) (*models.User, error) {
			return nil, errors.New("user not found")
		},
	}
	// Создаем обработчик
	handler := NewThreadHandler(mockThreadService)
	// Настраиваем тестовый роутер
	router := setupThreadTestRouter()
	router.GET("/threads", handler.GetAllThreads)
	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/threads", nil)
	// Выполняем запрос
	router.ServeHTTP(w, req)
	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)
	var threads []*models.Thread
	err := json.Unmarshal(w.Body.Bytes(), &threads)
	assert.NoError(t, err)
	assert.Len(t, threads, 1)
	assert.Empty(t, threads[0].AuthorName) // AuthorName должен быть пустым
}

func TestThreadHandler_GetThreadPosts_EmptyList(t *testing.T) {
	// Создаем мок сервиса
	mockThreadService := &mocks.MockThreadService{
		GetPostsByThreadIDFunc: func(threadID int) ([]*models.Post, error) {
			return []*models.Post{}, nil
		},
	}
	// Создаем обработчик
	handler := NewThreadHandler(mockThreadService)
	// Настраиваем тестовый роутер
	router := setupThreadTestRouter()
	router.GET("/threads/:id/posts", handler.GetThreadPosts)
	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/threads/1/posts", nil)
	// Выполняем запрос
	router.ServeHTTP(w, req)
	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, "[]", w.Body.String())
}

func TestThreadHandler_DeleteThread_NotFoundAfterGet(t *testing.T) {
	// Создаем мок сервиса
	mockThreadService := &mocks.MockThreadService{
		GetThreadWithPostsFunc: func(id int) (*models.Thread, []*models.Post, error) {
			return nil, nil, errors.New("database error")
		},
	}
	// Создаем обработчик
	handler := NewThreadHandler(mockThreadService)
	// Настраиваем тестовый роутер
	router := setupThreadTestRouter()
	router.DELETE("/threads/:id", func(c *gin.Context) {
		c.Set("user_id", uint32(1))
		c.Set("user_role", "user")
		handler.DeleteThread(c)
	})
	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/threads/1", nil)
	// Выполняем запрос
	router.ServeHTTP(w, req)
	// Проверяем результат
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestFormatDate(t *testing.T) {
	// Тест форматирования даты
	date := time.Date(2024, 3, 15, 14, 30, 0, 0, time.UTC)
	formatted := formatDate(date)
	assert.Equal(t, "15.03.2024 14:30", formatted)

	// Тест с нулевой датой
	zeroDate := time.Time{}
	formatted = formatDate(zeroDate)
	assert.Equal(t, "", formatted)
}

func TestValidateThreadTitle(t *testing.T) {
	// Тест валидного заголовка
	assert.True(t, validateThreadTitle("Valid Title"))
	assert.True(t, validateThreadTitle("Title with numbers 123"))
	assert.True(t, validateThreadTitle("Title with special chars: !@#$%"))

	// Тест невалидного заголовка
	assert.False(t, validateThreadTitle(""))
	assert.False(t, validateThreadTitle("   "))
	assert.False(t, validateThreadTitle(strings.Repeat("a", 256))) // слишком длинный
}

func TestSanitizeThreadTitle(t *testing.T) {
	// Тест очистки заголовка
	assert.Equal(t, "Clean Title", sanitizeThreadTitle("  Clean Title  "))
	assert.Equal(t, "&lt;script&gt;alert('xss')&lt;/script&gt;No HTML", sanitizeThreadTitle("<script>alert('xss')</script>No HTML"))
	assert.Equal(t, "Special  Chars", sanitizeThreadTitle("Special\n\tChars"))
}

func TestGetThreadStatus(t *testing.T) {
	// Тест статуса активной темы
	activeThread := &models.Thread{
		CreatedAt: time.Now().Add(-time.Hour),
		UpdatedAt: time.Now(),
	}
	assert.Equal(t, "active", getThreadStatus(activeThread))

	// Тест статуса неактивной темы
	inactiveThread := &models.Thread{
		CreatedAt: time.Now().Add(-25 * time.Hour),
		UpdatedAt: time.Now().Add(-25 * time.Hour),
	}
	assert.Equal(t, "inactive", getThreadStatus(inactiveThread))
}

func TestCalculateThreadStats(t *testing.T) {
	// Тест подсчета статистики
	thread := &models.Thread{
		ID: 1,
	}
	posts := []*models.Post{
		{ID: 1, AuthorID: 1},
		{ID: 2, AuthorID: 2},
		{ID: 3, AuthorID: 1},
	}
	stats := calculateThreadStats(thread, posts)
	assert.Equal(t, 3, stats.TotalPosts)
	assert.Equal(t, 2, stats.UniqueAuthors)
}

func TestCalculateThreadMetrics(t *testing.T) {
	now := time.Now()
	posts := []*models.Post{
		{
			ID:        1,
			Content:   "Короткий пост",
			AuthorID:  1,
			CreatedAt: now.Add(-2 * time.Hour),
		},
		{
			ID:        2,
			Content:   "Этот пост длиннее предыдущего",
			AuthorID:  2,
			CreatedAt: now.Add(-1 * time.Hour),
		},
		{
			ID:        3,
			Content:   "И этот пост тоже довольно длинный",
			AuthorID:  1,
			CreatedAt: now,
		},
	}

	metrics := calculateThreadMetrics(posts)
	assert.Equal(t, float64(47), metrics.AveragePostLength) // (12 + 28 + 17) / 3
	assert.Equal(t, 1, metrics.MostActiveAuthor) // Автор 1 написал 2 поста
	assert.Equal(t, now, metrics.LastActivityTime)
}

func TestIsThreadActive(t *testing.T) {
	now := time.Now()
	thread := &models.Thread{
		ID:        1,
		Title:     "Test Thread",
		CreatedAt: now.Add(-3 * 24 * time.Hour), // 3 дня назад
	}

	// Активная тема (есть посты за последние 24 часа)
	activePosts := []*models.Post{
		{CreatedAt: now.Add(-12 * time.Hour)},
		{CreatedAt: now.Add(-1 * time.Hour)},
	}
	assert.True(t, isThreadActive(thread, activePosts))

	// Неактивная тема (нет постов за последние 24 часа)
	inactivePosts := []*models.Post{
		{CreatedAt: now.Add(-48 * time.Hour)},
		{CreatedAt: now.Add(-36 * time.Hour)},
	}
	assert.False(t, isThreadActive(thread, inactivePosts))

	// Старая тема (более 7 дней)
	oldThread := &models.Thread{
		ID:        2,
		Title:     "Old Thread",
		CreatedAt: now.Add(-8 * 24 * time.Hour),
	}
	assert.False(t, isThreadActive(oldThread, activePosts))
}

func TestFormatThreadSummary(t *testing.T) {
	now := time.Now()
	thread := &models.Thread{
		ID:        1,
		Title:     "Test Thread",
		CreatedAt: now.Add(-1 * time.Hour),
	}
	posts := []*models.Post{
		{Content: "Первый пост", CreatedAt: now.Add(-30 * time.Minute)},
		{Content: "Второй пост", CreatedAt: now},
	}

	summary := formatThreadSummary(thread, posts)
	assert.Contains(t, summary, "Test Thread")
	assert.Contains(t, summary, "активна")
	assert.Contains(t, summary, "2")
	assert.Contains(t, summary, "21.0") // средняя длина постов
}

func TestValidateThreadContent(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantValid bool
		wantMsg   string
	}{
		{
			name:      "пустое сообщение",
			content:   "",
			wantValid: false,
			wantMsg:   "сообщение не может быть пустым",
		},
		{
			name:      "валидное сообщение",
			content:   "Это нормальное сообщение средней длины",
			wantValid: true,
			wantMsg:   "",
		},
		{
			name:      "короткое сообщение",
			content:   "Коротко",
			wantValid: false,
			wantMsg:   "сообщение слишком короткое",
		},
		{
			name:      "длинное сообщение",
			content:   strings.Repeat("a", 10001),
			wantValid: false,
			wantMsg:   "сообщение слишком длинное",
		},
		{
			name:      "много восклицательных знаков",
			content:   "Важно!!!!! Срочно!!!!!",
			wantValid: false,
			wantMsg:   "слишком много повторяющихся символов",
		},
		{
			name:      "много вопросительных знаков",
			content:   "Что????? Как?????",
			wantValid: false,
			wantMsg:   "слишком много повторяющихся символов",
		},
		{
			name:      "много заглавных букв",
			content:   "ВСЕ БУКВЫ ЗАГЛАВНЫЕ И ЭТО НЕ ХОРОШО ВСЕ БУКВЫ ЗАГЛАВНЫЕ И ЭТО НЕ ХОРОШО",
			wantValid: false,
			wantMsg:   "слишком много заглавных букв",
		},
		{
			name:      "нормальное количество заглавных букв",
			content:   "Нормальный текст с Заглавными буквами",
			wantValid: true,
			wantMsg:   "",
		},
		{
			name:      "ровно 10 символов",
			content:   "1234567890",
			wantValid: true,
			wantMsg:   "",
		},
		{
			name:      "ровно 10000 символов",
			content:   strings.Repeat("a", 10000),
			wantValid: true,
			wantMsg:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, msg := validateThreadContent(tt.content)
			assert.Equal(t, tt.wantValid, valid, "валидность не совпадает")
			assert.Equal(t, tt.wantMsg, msg, "сообщение об ошибке не совпадает")
		})
	}
}