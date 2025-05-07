package handlers

import (
	"ForumService/internal/models"
	"ForumService/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"fmt"
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

// func (h *ThreadHandler) RegisterRoutes(r *gin.RouterGroup) {
// 	threads := r.Group("/threads")
// 	{
// 		threads.POST("", h.CreateThread)
// 		threads.GET("/:id", h.GetThreadWithPosts)
// 		threads.DELETE("/:id", h.DeleteThread)
// 	}
// }

func (h *ThreadHandler) CreateThread(c *gin.Context) {
	var req CreateThreadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Ошибка при разборе JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Для тестирования используем фиксированный ID пользователя
	thread := &models.Thread{
		Title: req.Title,
		AuthorID: 1, // Фиксированный ID для тестирования
	}

	if err := h.service.CreateThread(thread); err != nil {
		fmt.Printf("Ошибка при создании треда: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, thread)
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

	if err := h.service.DeleteThread(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
