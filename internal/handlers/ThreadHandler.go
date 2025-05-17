package handlers

import (
	"ForumService/internal/models"
	"ForumService/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ThreadHandler struct {
	service service.ThreadService
}

func NewThreadHandler(service service.ThreadService) *ThreadHandler {
	return &ThreadHandler{service: service}
}

type CreateThreadRequest struct {
	Title string `json:"title" binding:"required"`
}

type UpdateThreadRequest struct {
	Title string `json:"title" binding:"required"`
}

// func (h *ThreadHandler) RegisterRoutes(r *gin.RouterGroup) {
// 	threads := r.Group("/threads")
// 	{
// 		threads.POST("", h.CreateThread)
// 		threads.GET("/:id", h.GetThreadWithPosts)
// 		threads.DELETE("/:id", h.DeleteThread)
// 	}
// }

func (h *ThreadHandler) CreateThread(c *gin.Context) {
	var request struct {
		Title string `json:"title" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "неверный формат данных"})
		return
	}

	// Получаем ID пользователя из контекста
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "пользователь не аутентифицирован"})
		return
	}

	// Преобразуем uint в int
	userIDInt := int(userID.(uint))

	thread, err := h.service.CreateThread(request.Title, userIDInt)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, thread)
}

func (h *ThreadHandler) GetThreadWithPosts(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid thread ID"})
		return
	}

	thread, posts, err := h.service.GetThreadWithPosts(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "thread not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"thread": thread,
		"posts":  posts,
	})
}

func (h *ThreadHandler) DeleteThread(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid thread ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userIDInt := int(userID.(uint))

	if err := h.service.DeleteThread(id, userIDInt); err != nil {
		if err == service.ErrNoPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "no permission to delete this thread"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ThreadHandler) UpdateThread(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid thread ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userIDInt := int(userID.(uint))

	var req UpdateThreadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	thread := &models.Thread{
		ID:    id,
		Title: req.Title,
	}

	if err := h.service.UpdateThread(thread, userIDInt); err != nil {
		if err == service.ErrNoPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "no permission to update this thread"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, thread)
}

func (h *ThreadHandler) GetAllThreads(c *gin.Context) {
	threads, err := h.service.GetAllThreads()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Добавляем информацию об авторе для каждого треда
	for _, thread := range threads {
		user, err := h.service.GetUserByID(thread.AuthorID)
		if err == nil && user != nil {
			thread.AuthorName = user.Username
		}
	}

	c.JSON(http.StatusOK, threads)
}

func (h *ThreadHandler) GetThreadPosts(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid thread ID"})
		return
	}

	posts, err := h.service.GetPostsByThreadID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "thread not found"})
		return
	}

	c.JSON(http.StatusOK, posts)
}
