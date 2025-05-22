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

// CreateComment godoc
// @Summary Создать новый комментарий
// @Description Создаёт новый комментарий к посту. Доступно только авторизованным пользователям.
// @Tags comments
// @Accept json
// @Produce json
// @Param input body object true "Данные для создания комментария"
// @Success 201 {object} models.Comment
// @Failure 400 {object} map[string]string "неверный формат данных"
// @Failure 401 {object} map[string]string "пользователь не аутентифицирован"
// @Failure 500 {object} map[string]string "ошибка сервера"
// @Router /comments [post]
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

	// Преобразуем uint32 в int
	userIDInt := int(userID.(uint32))

	comment, err := h.service.CreateComment(request.PostID, userIDInt, request.Content)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, comment)
}

// DeleteComment godoc
// @Summary Удалить комментарий
// @Description Удаляет комментарий. Доступно только автору комментария или администратору.
// @Tags comments
// @Produce json
// @Param id path int true "ID комментария"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "Неверный ID комментария"
// @Failure 401 {object} map[string]string "пользователь не аутентифицирован"
// @Failure 403 {object} map[string]string "нет прав для удаления этого комментария"
// @Failure 500 {object} map[string]string "ошибка сервера"
// @Router /comments/{id} [delete]
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

	// Преобразуем uint32 в int
	userIDInt := int(userID.(uint32))

	// Получаем роль пользователя
	userRole, _ := c.Get("user_role")
	if userRole == nil {
		userRole = "user"
	}

	// Отладочная информация
	fmt.Printf("Debug - DeleteComment - User ID: %d, Role: %v\n", userIDInt, userRole)
	fmt.Printf("Debug - DeleteComment - User Role type: %T\n", userRole)

	// Проверяем, является ли пользователь автором комментария или администратором
	comment, err := h.service.GetCommentByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при получении комментария"})
		return
	}

	if comment.AuthorID != userIDInt && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "нет прав для удаления этого комментария"})
		return
	}

	err = h.service.DeleteComment(id, userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// CreateChatMessage godoc
// @Summary Создать сообщение в чате
// @Description Создаёт новое сообщение в общем чате. Доступно только авторизованным пользователям.
// @Tags comments
// @Accept json
// @Produce json
// @Param input body object true "Данные для создания сообщения"
// @Success 201 {object} models.Comment
// @Failure 400 {object} map[string]string "Неверный формат данных"
// @Failure 500 {object} map[string]string "Ошибка при создании сообщения"
// @Router /comments/chat [post]
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
