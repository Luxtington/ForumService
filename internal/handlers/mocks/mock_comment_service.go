package mocks

import (
	"ForumService/internal/models"
)

type MockCommentService struct {
	CreateCommentFunc      func(postID int, authorID int, content string) (*models.Comment, error)
	GetCommentByIDFunc     func(id int) (*models.Comment, error)
	DeleteCommentFunc      func(id int, userID int) error
	GetCommentsByPostIDFunc func(postID int) ([]models.Comment, error)
}

func (m *MockCommentService) CreateComment(postID int, authorID int, content string) (*models.Comment, error) {
	return m.CreateCommentFunc(postID, authorID, content)
}

func (m *MockCommentService) GetCommentByID(id int) (*models.Comment, error) {
	return m.GetCommentByIDFunc(id)
}

func (m *MockCommentService) DeleteComment(id int, userID int) error {
	return m.DeleteCommentFunc(id, userID)
}

func (m *MockCommentService) GetCommentsByPostID(postID int) ([]models.Comment, error) {
	return m.GetCommentsByPostIDFunc(postID)
} 