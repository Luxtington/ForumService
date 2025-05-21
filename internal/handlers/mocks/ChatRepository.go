package mocks

import (
	"ForumService/internal/models"
	"sync/atomic"
)

type MockChatRepository struct {
	CreateMessageFunc     func(authorID int, content string) (*models.ChatMessage, error)
	createMessageCallCount int64
}

func (m *MockChatRepository) CreateMessage(authorID int, content string) (*models.ChatMessage, error) {
	atomic.AddInt64(&m.createMessageCallCount, 1)
	return m.CreateMessageFunc(authorID, content)
}

func (m *MockChatRepository) CreateMessageCallCount() int {
	return int(atomic.LoadInt64(&m.createMessageCallCount))
} 