package handlers

import (
	"ForumService/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"fmt"
	_ "ForumService/internal/models"
)

type CommentHandler struct {
	service service.CommentService
}

type CreateCommentRequest struct {
	PostID  string `json:"post_id"`
	Content string `json:"content"`
}

func NewCommentHandler(service service.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	var request struct {
		PostID  int    `json:"post_id" binding:"required"`
		Content string `json:"content" binding:"required"`
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

	comment, err := h.service.CreateComment(request.PostID, userIDInt, request.Content)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, comment)
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID комментария"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "пользователь не аутентифицирован"})
		return
	}

	userIDInt := int(userID.(uint))

	err = h.service.DeleteComment(id, userIDInt)
	if err != nil {
		if err == service.ErrNoPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "нет прав для удаления этого комментария"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *CommentHandler) CreateChatMessage(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Ошибка при разборе JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат данных",
		})
		return
	}

	// Для чата используем post_id = 0
	comment, err := h.service.CreateComment(0, 1, req.Content) // author_id = 1 для тестирования
	if err != nil {
		fmt.Printf("Ошибка при создании сообщения чата: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при создании сообщения",
		})
		return
	}

	c.JSON(http.StatusCreated, comment)
}
