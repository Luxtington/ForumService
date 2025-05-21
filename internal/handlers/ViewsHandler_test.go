package handlers

import (
	"errors"
	"ForumService/internal/handlers/mocks"
	"ForumService/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupViewsTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.LoadHTMLGlob("../../templates/*")
	return router
}

func TestViewsHandler_Index_Error(t *testing.T) {
	// Создаем моки сервисов
	mockThreadService := &mocks.MockThreadService{
		GetAllThreadsFunc: func() ([]*models.Thread, error) {
			return nil, errors.New("ошибка получения тредов")
		},
	}
	mockChatService := &mocks.MockChatService{
		GetAllMessagesFunc: func() ([]*models.ChatMessage, error) {
			return nil, errors.New("ошибка получения сообщений")
		},
	}

	// Создаем обработчик
	handler := NewViewsHandler(mockThreadService, nil, nil, mockChatService)

	// Настраиваем тестовый роутер
	router := setupViewsTestRouter()
	router.GET("/", handler.Index)

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestViewsHandler_ShowThread_InvalidID(t *testing.T) {
	// Создаем моки сервисов
	mockThreadService := &mocks.MockThreadService{}
	mockPostService := &mocks.MockPostService{}
	mockCommentService := &mocks.MockCommentService{}
	mockChatService := &mocks.MockChatService{}

	// Создаем обработчик
	handler := NewViewsHandler(mockThreadService, mockPostService, mockCommentService, mockChatService)

	// Настраиваем тестовый роутер
	router := setupViewsTestRouter()
	router.GET("/thread/:id", handler.ShowThread)

	// Создаем тестовый запрос с неверным ID
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/thread/invalid", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestViewsHandler_ShowThread_NotFound(t *testing.T) {
	// Создаем моки сервисов
	mockThreadService := &mocks.MockThreadService{
		GetThreadWithPostsFunc: func(id int) (*models.Thread, []*models.Post, error) {
			return nil, nil, errors.New("тред не найден")
		},
	}
	mockPostService := &mocks.MockPostService{}
	mockCommentService := &mocks.MockCommentService{}
	mockChatService := &mocks.MockChatService{}

	// Создаем обработчик
	handler := NewViewsHandler(mockThreadService, mockPostService, mockCommentService, mockChatService)

	// Настраиваем тестовый роутер
	router := setupViewsTestRouter()
	router.GET("/thread/:id", handler.ShowThread)

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/thread/999", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestViewsHandler_ShowPost_InvalidID(t *testing.T) {
	// Создаем моки сервисов
	mockThreadService := &mocks.MockThreadService{}
	mockPostService := &mocks.MockPostService{}
	mockCommentService := &mocks.MockCommentService{}
	mockChatService := &mocks.MockChatService{}

	// Создаем обработчик
	handler := NewViewsHandler(mockThreadService, mockPostService, mockCommentService, mockChatService)

	// Настраиваем тестовый роутер
	router := setupViewsTestRouter()
	router.GET("/post/:id", handler.ShowPost)

	// Создаем тестовый запрос с неверным ID
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/post/invalid", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestViewsHandler_ShowPost_Error(t *testing.T) {
	// Создаем моки сервисов
	mockThreadService := &mocks.MockThreadService{}
	mockPostService := &mocks.MockPostService{
		GetPostWithCommentsFunc: func(postID int) (*models.Post, []models.Comment, error) {
			return nil, nil, errors.New("ошибка получения поста")
		},
	}
	mockCommentService := &mocks.MockCommentService{}
	mockChatService := &mocks.MockChatService{}

	// Создаем обработчик
	handler := NewViewsHandler(mockThreadService, mockPostService, mockCommentService, mockChatService)

	// Настраиваем тестовый роутер
	router := setupViewsTestRouter()
	router.GET("/post/:id", handler.ShowPost)

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/post/999", nil)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestViewsHandler_GetThreadWithPosts_InvalidURL(t *testing.T) {
	// Создаем моки сервисов
	mockThreadService := &mocks.MockThreadService{}
	mockPostService := &mocks.MockPostService{}
	mockCommentService := &mocks.MockCommentService{}
	mockChatService := &mocks.MockChatService{}

	// Создаем обработчик
	handler := NewViewsHandler(mockThreadService, mockPostService, mockCommentService, mockChatService)

	// Создаем тестовый запрос с неверным URL
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/thread", nil)

	// Выполняем запрос
	handler.GetThreadWithPosts(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestViewsHandler_GetThreadWithPosts_InvalidID(t *testing.T) {
	// Создаем моки сервисов
	mockThreadService := &mocks.MockThreadService{}
	mockPostService := &mocks.MockPostService{}
	mockCommentService := &mocks.MockCommentService{}
	mockChatService := &mocks.MockChatService{}

	// Создаем обработчик
	handler := NewViewsHandler(mockThreadService, mockPostService, mockCommentService, mockChatService)

	// Создаем тестовый запрос с неверным ID
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/thread/invalid", nil)

	// Выполняем запрос
	handler.GetThreadWithPosts(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestViewsHandler_GetThreadWithPosts_NotFound(t *testing.T) {
	// Создаем моки сервисов
	mockThreadService := &mocks.MockThreadService{
		GetThreadWithPostsFunc: func(id int) (*models.Thread, []*models.Post, error) {
			return nil, nil, errors.New("тред не найден")
		},
	}
	mockPostService := &mocks.MockPostService{}
	mockCommentService := &mocks.MockCommentService{}
	mockChatService := &mocks.MockChatService{}

	// Создаем обработчик
	handler := NewViewsHandler(mockThreadService, mockPostService, mockCommentService, mockChatService)

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/thread/999", nil)

	// Выполняем запрос
	handler.GetThreadWithPosts(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusNotFound, w.Code)
} 