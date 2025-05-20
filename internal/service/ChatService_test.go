package service

import (
	"testing"
	"ForumService/internal/models"
	"ForumService/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"errors"
)

func TestCreateMessage_Success(t *testing.T) {
	repo := new(mocks.MockChatRepo)
	service := NewChatService(repo)
	msg := &models.ChatMessage{ID: 1, AuthorID: 2, Content: "hi"}
	repo.On("CreateMessage", 2, "hi").Return(msg, nil)
	res, err := service.CreateMessage(2, "hi")
	assert.NoError(t, err)
	assert.Equal(t, msg, res)
}

func TestGetAllMessages_Success(t *testing.T) {
	repo := new(mocks.MockChatRepo)
	service := NewChatService(repo)
	msgs := []*models.ChatMessage{{ID: 1, AuthorID: 2, Content: "hi"}}
	repo.On("GetAllMessages").Return(msgs, nil)
	res, err := service.GetAllMessages()
	assert.NoError(t, err)
	assert.Equal(t, msgs, res)
}

func TestGetAllMessages_Error(t *testing.T) {
	repo := new(mocks.MockChatRepo)
	service := NewChatService(repo)
	repo.On("GetAllMessages").Return(([]*models.ChatMessage)(nil), errors.New("db error"))
	res, err := service.GetAllMessages()
	assert.Error(t, err)
	assert.Nil(t, res)
} 