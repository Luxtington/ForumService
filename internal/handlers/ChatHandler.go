package handlers

import (
	"ForumService/internal/service"
	"ForumService/internal/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ChatHandler struct {
	service service.ChatService
}

func NewChatHandler(service service.ChatService) *ChatHandler {
	return &ChatHandler{service: service}
}

type CreateMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

// CreateMessage godoc
// @Summary Создать сообщение в чате
// @Description Создаёт новое сообщение в общем чате. Доступно только авторизованным пользователям.
// @Tags chat
// @Accept json
// @Produce json
// @Param input body object true "Данные для создания сообщения"
// @Success 201 {object} object "Returns created message"
// @Failure 400 {object} map[string]string "неверный формат данных"
// @Failure 401 {object} map[string]string "пользователь не аутентифицирован"
// @Failure 500 {object} map[string]string "ошибка сервера"
// @Router /chat [post]
func (h *ChatHandler) CreateMessage(c *gin.Context) {
	var request CreateMessageRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(errors.NewValidationError("Неверный формат данных", err))
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.Error(errors.NewUnauthorizedError("Пользователь не аутентифицирован", nil))
		return
	}

	userIDInt := int(userID.(uint32))
	message, err := h.service.CreateMessage(userIDInt, request.Content)
	if err != nil {
		c.Error(errors.NewInternalServerError("Ошибка при создании сообщения", err))
		return
	}

	c.JSON(http.StatusCreated, message)
}

// GetMessages godoc
// @Summary Получить все сообщения чата
// @Description Возвращает список всех сообщений в общем чате.
// @Tags chat
// @Produce json
// @Success 200 {array} object "Returns list of messages"
// @Failure 500 {object} map[string]string "ошибка сервера"
// @Router /chat [get]
func (h *ChatHandler) GetMessages(c *gin.Context) {
	messages, err := h.service.GetAllMessages()
	if err != nil {
		c.Error(errors.NewInternalServerError("Ошибка при получении сообщений", err))
		return
	}

	c.JSON(http.StatusOK, messages)
} 