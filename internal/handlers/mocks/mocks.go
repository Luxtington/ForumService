package mocks

import (
	"ForumService/internal/models"
)

type MockUserService struct {
	RegisterFunc    func(username, password string) (*models.User, error)
	LoginFunc       func(username, password string) (*models.User, error)
	GetUserByIDFunc func(id int) (*models.User, error)
}

func (m *MockUserService) Register(username, password string) (*models.User, error) {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(username, password)
	}
	return nil, nil
}

func (m *MockUserService) Login(username, password string) (*models.User, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(username, password)
	}
	return nil, nil
}

func (m *MockUserService) GetUserByID(id int) (*models.User, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(id)
	}
	return nil, nil
} 