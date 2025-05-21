package mocks

import (
	"ForumService/internal/models"
)

type MockPostService struct {
	CreatePostFunc func(post *models.Post) error
	GetPostByIDFunc func(id int) (*models.Post, error)
	GetPostWithCommentsFunc func(postID int) (*models.Post, []models.Comment, error)
	GetPostsWithCommentsByThreadIDFunc func(threadID int) ([]models.Post, map[int][]models.Comment, error)
	UpdatePostFunc func(post *models.Post, postID int, userID int) error
	DeletePostFunc func(postID int, userID int) error
	GetAllPostsFunc func() ([]*models.Post, error)
	CreateCommentFunc func(comment *models.Comment) error
	GetCommentByIDFunc func(id int) (*models.Comment, error)
	DeleteCommentFunc func(id int) error
	GetPostFunc func(id int) (*models.Post, error)
	GetPostsByThreadIDFunc func(threadID int) ([]*models.Post, error)
	GetThreadByIDFunc func(id int) (*models.Thread, error)
}

func (m *MockPostService) CreatePost(post *models.Post) error {
	return m.CreatePostFunc(post)
}

func (m *MockPostService) GetPostByID(id int) (*models.Post, error) {
	return m.GetPostByIDFunc(id)
}

func (m *MockPostService) GetPostWithComments(postID int) (*models.Post, []models.Comment, error) {
	return m.GetPostWithCommentsFunc(postID)
}

func (m *MockPostService) GetPostsWithCommentsByThreadID(threadID int) ([]models.Post, map[int][]models.Comment, error) {
	return m.GetPostsWithCommentsByThreadIDFunc(threadID)
}

func (m *MockPostService) UpdatePost(post *models.Post, postID int, userID int) error {
	return m.UpdatePostFunc(post, postID, userID)
}

func (m *MockPostService) DeletePost(postID int, userID int) error {
	return m.DeletePostFunc(postID, userID)
}

func (m *MockPostService) GetAllPosts() ([]*models.Post, error) {
	return m.GetAllPostsFunc()
}

func (m *MockPostService) CreateComment(comment *models.Comment) error {
	return m.CreateCommentFunc(comment)
}

func (m *MockPostService) GetCommentByID(id int) (*models.Comment, error) {
	return m.GetCommentByIDFunc(id)
}

func (m *MockPostService) DeleteComment(id int) error {
	return m.DeleteCommentFunc(id)
}

func (m *MockPostService) GetPost(id int) (*models.Post, error) {
	return m.GetPostFunc(id)
}

func (m *MockPostService) GetPostsByThreadID(threadID int) ([]*models.Post, error) {
	return m.GetPostsByThreadIDFunc(threadID)
}

func (m *MockPostService) GetThreadByID(id int) (*models.Thread, error) {
	return m.GetThreadByIDFunc(id)
} 