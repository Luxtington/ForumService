package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "ForumService/internal/service"
)

type ChatHandler struct {
    service service.ChatService
}

func NewChatHandler(service service.ChatService) *ChatHandler {
    return &ChatHandler{service: service}
}

func (h *ChatHandler) CreateMessage(c *gin.Context) {
    var request struct {
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

    message, err := h.service.CreateMessage(userIDInt, request.Content)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(201, message)
}

func (h *ChatHandler) GetMessages(c *gin.Context) {
    messages, err := h.service.GetAllMessages()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, messages)
} 