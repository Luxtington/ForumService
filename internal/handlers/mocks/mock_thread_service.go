package mocks

import (
	"ForumService/internal/models"
)

type MockThreadService struct {
	CreateThreadFunc      func(title string, authorID int) (*models.Thread, error)
	GetThreadByIDFunc     func(id int) (*models.Thread, error)
	GetThreadWithPostsFunc func(id int) (*models.Thread, []*models.Post, error)
	DeleteThreadFunc      func(id int, userID int) error
	UpdateThreadFunc      func(thread *models.Thread, userID int) error
	GetAllThreadsFunc     func() ([]*models.Thread, error)
	GetPostsByThreadIDFunc func(id int) ([]*models.Post, error)
	GetUserByIDFunc       func(id int) (*models.User, error)
}

func (m *MockThreadService) CreateThread(title string, authorID int) (*models.Thread, error) {
	return m.CreateThreadFunc(title, authorID)
}

func (m *MockThreadService) GetThreadByID(id int) (*models.Thread, error) {
	return m.GetThreadByIDFunc(id)
}

func (m *MockThreadService) GetThreadWithPosts(id int) (*models.Thread, []*models.Post, error) {
	return m.GetThreadWithPostsFunc(id)
}

func (m *MockThreadService) DeleteThread(id int, userID int) error {
	return m.DeleteThreadFunc(id, userID)
}

func (m *MockThreadService) UpdateThread(thread *models.Thread, userID int) error {
	return m.UpdateThreadFunc(thread, userID)
}

func (m *MockThreadService) GetAllThreads() ([]*models.Thread, error) {
	return m.GetAllThreadsFunc()
}

func (m *MockThreadService) GetPostsByThreadID(id int) ([]*models.Post, error) {
	return m.GetPostsByThreadIDFunc(id)
}

func (m *MockThreadService) GetUserByID(id int) (*models.User, error) {
	return m.GetUserByIDFunc(id)
} 