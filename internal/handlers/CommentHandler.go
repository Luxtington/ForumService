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
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Ошибка при разборе JSON: %v\n", err)
		fmt.Printf("Полученные данные: %+v\n", req)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при разборе JSON: " + err.Error()})
		return
	}

	fmt.Printf("Получен запрос на создание комментария: %+v\n", req)

	// Получаем ID пользователя из контекста
	userID, exists := c.Get("user_id")
	if !exists {
		fmt.Printf("Пользователь не авторизован\n")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	fmt.Printf("ID пользователя: %v\n", userID)

	// Преобразуем post_id из строки в int
	postID, err := strconv.Atoi(req.PostID)
	if err != nil {
		fmt.Printf("Ошибка при преобразовании post_id: %v\n", err)
		fmt.Printf("Полученный post_id: %v\n", req.PostID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат post_id"})
		return
	}

	fmt.Printf("Преобразованный post_id: %v\n", postID)

	comment, err := h.service.CreateComment(postID, userID.(int), req.Content)
	if err != nil {
		fmt.Printf("Ошибка при создании комментария: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Комментарий успешно создан: %+v\n", comment)
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
