package mocks

import (
	"ForumService/internal/models"
)

type MockChatService struct {
	CreateMessageFunc     func(authorID int, content string) (*models.ChatMessage, error)
	GetAllMessagesFunc    func() ([]*models.ChatMessage, error)
}

func (m *MockChatService) CreateMessage(authorID int, content string) (*models.ChatMessage, error) {
	return m.CreateMessageFunc(authorID, content)
}

func (m *MockChatService) GetAllMessages() ([]*models.ChatMessage, error) {
	return m.GetAllMessagesFunc()
} 