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
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Для тестирования используем фиксированный ID пользователя
    message, err := h.service.CreateMessage(1, request.Content)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, message)
}

func (h *ChatHandler) GetMessages(c *gin.Context) {
    messages, err := h.service.GetAllMessages()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, messages)
} 