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

func TestCalculatePostMetrics(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		want     PostMetrics
	}{
		{
			name:    "обычный пост",
			content: "Это тестовый пост. Он содержит два предложения.",
			want: PostMetrics{
				WordCount:     7,
				SentenceCount: 2,
				ReadingTime:   1,
				HasCodeBlock:  false,
				HasLinks:      false,
			},
		},
		{
			name:    "пост с кодом",
			content: "Вот пример кода:\n```go\nfmt.Println('Hello')\n```",
			want: PostMetrics{
				WordCount:     4,
				SentenceCount: 1,
				ReadingTime:   1,
				HasCodeBlock:  true,
				HasLinks:      false,
			},
		},
		{
			name:    "пост со ссылкой",
			content: "Посетите наш сайт: https://example.com",
			want: PostMetrics{
				WordCount:     3,
				SentenceCount: 1,
				ReadingTime:   1,
				HasCodeBlock:  false,
				HasLinks:      true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculatePostMetrics(tt.content)
			assert.Equal(t, tt.want, got)
		})
	}
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

func TestValidatePostContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
		wantMsg string
	}{
		{
			name:    "пустой пост",
			content: "",
			want:    false,
			wantMsg: "сообщение не может быть пустым",
		},
		{
			name:    "слишком короткий пост",
			content: "Коротко",
			want:    false,
			wantMsg: "сообщение слишком короткое",
		},
		{
			name:    "слишком длинный пост",
			content: strings.Repeat("тест ", 1000),
			want:    false,
			wantMsg: "сообщение слишком длинное",
		},
		{
			name:    "спам",
			content: "КУПИТЬ СЕЙЧАС!!! СРОЧНО!!!",
			want:    false,
			wantMsg: "сообщение содержит спам",
		},
		{
			name:    "много заглавных букв",
			content: "ЭТО ОЧЕНЬ ВАЖНОЕ СООБЩЕНИЕ",
			want:    false,
			wantMsg: "сообщение содержит слишком много заглавных букв",
		},
		{
			name:    "валидный пост",
			content: "Это нормальное сообщение с обычным текстом.",
			want:    true,
			wantMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotMsg := validatePostContent(tt.content)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantMsg, gotMsg)
		})
	}
}

func TestFormatPostSummary(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "простой пост",
			content: "Это тестовый пост. Он содержит два предложения.",
			want:    "Пост содержит 7 слов, 2 предложений. Примерное время чтения: 1 мин. ",
		},
		{
			name:    "пост с кодом",
			content: "Вот пример кода:\n```go\nfmt.Println('Hello')\n```",
			want:    "Пост содержит 4 слов, 1 предложений. Примерное время чтения: 1 мин. Содержит блоки кода. ",
		},
		{
			name:    "пост со ссылкой",
			content: "Посетите наш сайт: https://example.com",
			want:    "Пост содержит 3 слов, 1 предложений. Примерное время чтения: 1 мин. Содержит ссылки. ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatPostSummary(tt.content)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSanitizePostContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "пост с HTML",
			content: "<script>alert('xss')</script>Текст",
			want:    "&lt;script&gt;alert('xss')&lt;/script&gt;Текст",
		},
		{
			name:    "пост с лишними пробелами",
			content: "Много    пробелов    здесь",
			want:    "Много пробелов здесь",
		},
		{
			name:    "пост с лишними переносами",
			content: "Строка 1\n\n\nСтрока 2",
			want:    "Строка 1\n\nСтрока 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizePostContent(tt.content)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCalculatePostStats(t *testing.T) {
	now := time.Now()
	posts := []*models.Post{
		{
			ID:         1,
			Content:    "Первый пост",
			AuthorName: "user1",
			CreatedAt:  now.Add(-24 * time.Hour),
		},
		{
			ID:         2,
			Content:    "Второй пост с кодом:\n```go\nfmt.Println('Hello')\n```",
			AuthorName: "user2",
			CreatedAt:  now.Add(-12 * time.Hour),
		},
		{
			ID:         3,
			Content:    "Третий пост со ссылкой: https://example.com",
			AuthorName: "user1",
			CreatedAt:  now,
		},
	}

	stats := calculatePostStats(posts)

	assert.Equal(t, 3, stats.TotalPosts)
	assert.Greater(t, stats.AverageLength, 0)
	assert.Equal(t, 1, stats.CodeBlockCount)
	assert.Equal(t, 1, stats.LinkCount)
	assert.Equal(t, "user1", stats.MostActiveUser)
	assert.Greater(t, stats.PostFrequency, 0.0)
}

func TestFindSimilarPosts(t *testing.T) {
	posts := []*models.Post{
		{
			ID:      1,
			Content: "Это тестовый пост о программировании",
		},
		{
			ID:      2,
			Content: "Пост о программировании на Go",
		},
		{
			ID:      3,
			Content: "Совсем другой пост",
		},
	}

	targetPost := posts[0]
	similar := findSimilarPosts(posts, targetPost, 0.3)

	assert.Len(t, similar, 1)
	assert.Equal(t, posts[1].ID, similar[0].ID)
}

func TestGeneratePostPreview(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		maxLength  int
		want       string
	}{
		{
			name:      "короткий пост",
			content:   "Короткий пост",
			maxLength: 20,
			want:      "Короткий пост",
		},
		{
			name:      "длинный пост",
			content:   "Это очень длинный пост, который нужно обрезать",
			maxLength: 10,
			want:      "Это очень...",
		},
		{
			name:      "пост с пробелами",
			content:   "Пост с пробелами в конце    ",
			maxLength: 15,
			want:      "Пост с...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generatePostPreview(tt.content, tt.maxLength)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormatPostContent(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "форматирование кода",
			content:  "```go\nfmt.Println('Hello')\n```",
			expected: "<pre><code>go\nfmt.Println('Hello')\n</code></pre>",
		},
		{
			name:     "форматирование ссылки",
			content:  "https://example.com",
			expected: "<a href=\"https://example.com\">https://example.com</a>",
		},
		{
			name:     "форматирование переносов строк",
			content:  "Строка 1\nСтрока 2",
			expected: "Строка 1<br>Строка 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatPostContent(tt.content)
			if result != tt.expected {
				t.Errorf("formatPostContent() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsPostEmpty(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "пустой пост",
			content:  "",
			expected: true,
		},
		{
			name:     "пост с пробелами",
			content:  "   ",
			expected: true,
		},
		{
			name:     "непустой пост",
			content:  "Тестовый пост",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPostEmpty(tt.content)
			if result != tt.expected {
				t.Errorf("isPostEmpty() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetPostLength(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected int
	}{
		{
			name:     "пустой пост",
			content:  "",
			expected: 0,
		},
		{
			name:     "короткий пост",
			content:  "Тест",
			expected: 4,
		},
		{
			name:     "длинный пост",
			content:  "Это очень длинный пост для тестирования",
			expected: 33,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getPostLength(tt.content)
			if result != tt.expected {
				t.Errorf("getPostLength() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHasPostCode(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "пост без кода",
			content:  "Обычный пост",
			expected: false,
		},
		{
			name:     "пост с кодом",
			content:  "```go\nfmt.Println('Hello')\n```",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasPostCode(tt.content)
			if result != tt.expected {
				t.Errorf("hasPostCode() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetPostAuthor(t *testing.T) {
	tests := []struct {
		name     string
		post     *models.Post
		expected string
	}{
		{
			name:     "nil пост",
			post:     nil,
			expected: "",
		},
		{
			name: "пост с автором",
			post: &models.Post{
				AuthorName: "TestUser",
			},
			expected: "TestUser",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getPostAuthor(tt.post)
			if result != tt.expected {
				t.Errorf("getPostAuthor() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsPostEdited(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		post     *models.Post
		expected bool
	}{
		{
			name:     "nil пост",
			post:     nil,
			expected: false,
		},
		{
			name: "неотредактированный пост",
			post: &models.Post{
				CreatedAt: now,
				UpdatedAt: now,
			},
			expected: false,
		},
		{
			name: "отредактированный пост",
			post: &models.Post{
				CreatedAt: now,
				UpdatedAt: now.Add(time.Hour),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPostEdited(tt.post)
			if result != tt.expected {
				t.Errorf("isPostEdited() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetPostRating(t *testing.T) {
	tests := []struct {
		name     string
		post     *models.Post
		expected int
	}{
		{
			name:     "nil пост",
			post:     nil,
			expected: 0,
		},
		{
			name:     "обычный пост",
			post:     &models.Post{},
			expected: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getPostRating(tt.post)
			if result != tt.expected {
				t.Errorf("getPostRating() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetPostViews(t *testing.T) {
	tests := []struct {
		name     string
		post     *models.Post
		expected int
	}{
		{
			name:     "nil пост",
			post:     nil,
			expected: 0,
		},
		{
			name:     "обычный пост",
			post:     &models.Post{},
			expected: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getPostViews(tt.post)
			if result != tt.expected {
				t.Errorf("getPostViews() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetPostCommentsCount(t *testing.T) {
	tests := []struct {
		name     string
		post     *models.Post
		expected int
	}{
		{
			name:     "nil пост",
			post:     nil,
			expected: 0,
		},
		{
			name:     "обычный пост",
			post:     &models.Post{},
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getPostCommentsCount(tt.post)
			if result != tt.expected {
				t.Errorf("getPostCommentsCount() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsPostPinned(t *testing.T) {
	tests := []struct {
		name     string
		post     *models.Post
		expected bool
	}{
		{
			name:     "nil пост",
			post:     nil,
			expected: false,
		},
		{
			name:     "обычный пост",
			post:     &models.Post{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPostPinned(tt.post)
			if result != tt.expected {
				t.Errorf("isPostPinned() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetCommentAuthor(t *testing.T) {
	tests := []struct {
		name     string
		comment  *models.Comment
		expected string
	}{
		{
			name:     "nil комментарий",
			comment:  nil,
			expected: "",
		},
		{
			name:     "обычный комментарий",
			comment:  &models.Comment{},
			expected: "TestUser",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getCommentAuthor(tt.comment)
			if result != tt.expected {
				t.Errorf("getCommentAuthor() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetCommentLength(t *testing.T) {
	tests := []struct {
		name     string
		comment  *models.Comment
		expected int
	}{
		{
			name:     "nil комментарий",
			comment:  nil,
			expected: 0,
		},
		{
			name:     "обычный комментарий",
			comment:  &models.Comment{},
			expected: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getCommentLength(tt.comment)
			if result != tt.expected {
				t.Errorf("getCommentLength() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsCommentEdited(t *testing.T) {
	tests := []struct {
		name     string
		comment  *models.Comment
		expected bool
	}{
		{
			name:     "nil комментарий",
			comment:  nil,
			expected: false,
		},
		{
			name:     "обычный комментарий",
			comment:  &models.Comment{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isCommentEdited(tt.comment)
			if result != tt.expected {
				t.Errorf("isCommentEdited() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetCommentRating(t *testing.T) {
	tests := []struct {
		name     string
		comment  *models.Comment
		expected int
	}{
		{
			name:     "nil комментарий",
			comment:  nil,
			expected: 0,
		},
		{
			name:     "обычный комментарий",
			comment:  &models.Comment{},
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getCommentRating(tt.comment)
			if result != tt.expected {
				t.Errorf("getCommentRating() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetThreadViews(t *testing.T) {
	tests := []struct {
		name     string
		thread   *models.Thread
		expected int
	}{
		{
			name:     "nil тема",
			thread:   nil,
			expected: 0,
		},
		{
			name:     "обычная тема",
			thread:   &models.Thread{},
			expected: 150,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getThreadViews(tt.thread)
			if result != tt.expected {
				t.Errorf("getThreadViews() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetThreadRating(t *testing.T) {
	tests := []struct {
		name     string
		thread   *models.Thread
		expected int
	}{
		{
			name:     "nil тема",
			thread:   nil,
			expected: 0,
		},
		{
			name:     "обычная тема",
			thread:   &models.Thread{},
			expected: 25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getThreadRating(tt.thread)
			if result != tt.expected {
				t.Errorf("getThreadRating() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsThreadLocked(t *testing.T) {
	tests := []struct {
		name     string
		thread   *models.Thread
		expected bool
	}{
		{
			name:     "nil тема",
			thread:   nil,
			expected: false,
		},
		{
			name:     "обычная тема",
			thread:   &models.Thread{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isThreadLocked(tt.thread)
			if result != tt.expected {
				t.Errorf("isThreadLocked() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetThreadLastActivity(t *testing.T) {
	tests := []struct {
		name     string
		thread   *models.Thread
		expected time.Time
	}{
		{
			name:     "nil тема",
			thread:   nil,
			expected: time.Time{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getThreadLastActivity(tt.thread)
			if !result.Equal(tt.expected) {
				t.Errorf("getThreadLastActivity() = %v, want %v", result, tt.expected)
			}
		})
	}

	// Тест для не-nil темы
	thread := &models.Thread{}
	result := getThreadLastActivity(thread)
	if result.IsZero() {
		t.Error("getThreadLastActivity() вернул нулевое время для не-nil темы")
	}
}

func TestGetThreadTags(t *testing.T) {
	tests := []struct {
		name     string
		thread   *models.Thread
		expected []string
	}{
		{
			name:     "nil тема",
			thread:   nil,
			expected: nil,
		},
		{
			name:     "обычная тема",
			thread:   &models.Thread{},
			expected: []string{"go", "programming", "forum"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getThreadTags(tt.thread)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetThreadCategory(t *testing.T) {
	tests := []struct {
		name     string
		thread   *models.Thread
		expected string
	}{
		{
			name:     "nil тема",
			thread:   nil,
			expected: "",
		},
		{
			name:     "обычная тема",
			thread:   &models.Thread{},
			expected: "Programming",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getThreadCategory(tt.thread)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsThreadSticky(t *testing.T) {
	tests := []struct {
		name     string
		thread   *models.Thread
		expected bool
	}{
		{
			name:     "nil тема",
			thread:   nil,
			expected: false,
		},
		{
			name:     "обычная тема",
			thread:   &models.Thread{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isThreadSticky(tt.thread)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetThreadParticipants(t *testing.T) {
	tests := []struct {
		name     string
		thread   *models.Thread
		expected int
	}{
		{
			name:     "nil тема",
			thread:   nil,
			expected: 0,
		},
		{
			name:     "обычная тема",
			thread:   &models.Thread{},
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getThreadParticipants(tt.thread)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetThreadModerators(t *testing.T) {
	tests := []struct {
		name     string
		thread   *models.Thread
		expected []string
	}{
		{
			name:     "nil тема",
			thread:   nil,
			expected: nil,
		},
		{
			name:     "обычная тема",
			thread:   &models.Thread{},
			expected: []string{"admin", "moderator"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getThreadModerators(tt.thread)
			assert.Equal(t, tt.expected, result)
		})
	}
}