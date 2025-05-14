package service

import (
    "ForumService/internal/models"
    "ForumService/internal/repository"
)

type ChatService interface {
    CreateMessage(authorID int, content string) (*models.ChatMessage, error)
    GetAllMessages() ([]*models.ChatMessage, error)
}

type chatService struct {
    repo repository.ChatRepository
}

func NewChatService(repo repository.ChatRepository) ChatService {
    return &chatService{repo: repo}
}

func (s *chatService) CreateMessage(authorID int, content string) (*models.ChatMessage, error) {
    return s.repo.CreateMessage(authorID, content)
}

func (s *chatService) GetAllMessages() ([]*models.ChatMessage, error) {
    return s.repo.GetAllMessages()
} 