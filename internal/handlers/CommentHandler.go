package handlers

import (
	"ForumService/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"fmt"
)

type CommentHandler struct {
	service service.CommentService
}

func NewCommentHandler(service service.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req struct {
		PostID  int    `json:"post_id" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Ошибка при разборе JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат данных",
		})
		return
	}

	// Для тестирования устанавливаем author_id = 1
	comment, err := h.service.CreateComment(req.PostID, 1, req.Content)
	if err != nil {
		fmt.Printf("Ошибка при создании комментария: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при создании комментария",
		})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID комментария"})
		return
	}

	// Получаем userID и isAdmin из контекста с проверкой на nil
	userID, _ := c.Get("user_id")
	isAdmin, _ := c.Get("is_admin")

	// Устанавливаем значения по умолчанию, если они nil
	var userIDInt int
	if userID == nil {
		userIDInt = 1 // Для тестирования используем ID 1
	} else {
		userIDInt = userID.(int)
	}

	var isAdminBool bool
	if isAdmin == nil {
		isAdminBool = false
	} else {
		isAdminBool = isAdmin.(bool)
	}

	err = h.service.DeleteComment(id, userIDInt, isAdminBool)
	if err != nil {
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
