package mocks

import (
	"ForumService/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockThreadRepo struct{ mock.Mock }
func (m *MockThreadRepo) Create(thread *models.Thread) error { args := m.Called(thread); return args.Error(0) }
func (m *MockThreadRepo) GetByID(id int) (*models.Thread, error) { args := m.Called(id); return args.Get(0).(*models.Thread), args.Error(1) }
func (m *MockThreadRepo) Update(thread *models.Thread) error { args := m.Called(thread); return args.Error(0) }
func (m *MockThreadRepo) Delete(id int) error { args := m.Called(id); return args.Error(0) }
func (m *MockThreadRepo) GetAllThreads() ([]*models.Thread, error) { args := m.Called(); return args.Get(0).([]*models.Thread), args.Error(1) }
func (m *MockThreadRepo) GetThreadWithPosts(threadID int) (*models.Thread, []models.Post, map[int][]models.Comment, error) { args := m.Called(threadID); return args.Get(0).(*models.Thread), args.Get(1).([]models.Post), args.Get(2).(map[int][]models.Comment), args.Error(3) }

type MockPostRepo struct{ mock.Mock }
func (m *MockPostRepo) SavePost(post *models.Post) error { args := m.Called(post); return args.Error(0) }
func (m *MockPostRepo) GetPostByID(id int) (*models.Post, error) { args := m.Called(id); return args.Get(0).(*models.Post), args.Error(1) }
func (m *MockPostRepo) GetPostWithComments(postID int) (*models.Post, []models.Comment, error) { args := m.Called(postID); return args.Get(0).(*models.Post), args.Get(1).([]models.Comment), args.Error(2) }
func (m *MockPostRepo) GetPostsWithCommentsByThreadID(threadID int) ([]models.Post, map[int][]models.Comment, error) { args := m.Called(threadID); return args.Get(0).([]models.Post), args.Get(1).(map[int][]models.Comment), args.Error(2) }
func (m *MockPostRepo) UpdatePost(post *models.Post, postID int) error { args := m.Called(post, postID); return args.Error(0) }
func (m *MockPostRepo) DeletePost(postID int) error { args := m.Called(postID); return args.Error(0) }
func (m *MockPostRepo) GetByThreadID(threadID int) ([]*models.Post, error) { args := m.Called(threadID); return args.Get(0).([]*models.Post), args.Error(1) }

type MockCommentRepo struct{ mock.Mock }
func (m *MockCommentRepo) SaveComment(comment *models.Comment) error { args := m.Called(comment); return args.Error(0) }
func (m *MockCommentRepo) GetCommentByID(id int) (*models.Comment, error) { args := m.Called(id); return args.Get(0).(*models.Comment), args.Error(1) }
func (m *MockCommentRepo) DeleteComment(id int) error { args := m.Called(id); return args.Error(0) }
func (m *MockCommentRepo) GetCommentsByPostID(postID int) ([]models.Comment, error) { args := m.Called(postID); return args.Get(0).([]models.Comment), args.Error(1) }

type MockUserRepo struct{ mock.Mock }
func (m *MockUserRepo) GetUserByID(id int) (*models.User, error) { args := m.Called(id); return args.Get(0).(*models.User), args.Error(1) }
func (m *MockUserRepo) GetUserRole(userID int) (string, error) { args := m.Called(userID); return args.String(0), args.Error(1) }
func (m *MockUserRepo) GetUserPosts(userID int) ([]*models.Post, error) { args := m.Called(userID); return args.Get(0).([]*models.Post), args.Error(1) }
func (m *MockUserRepo) GetUserCommentCount(userID int) (int, error) { args := m.Called(userID); return args.Int(0), args.Error(1) }
func (m *MockUserRepo) SaveUser(user *models.User) error { args := m.Called(user); return args.Error(0) }

type MockChatRepo struct{ mock.Mock }
func (m *MockChatRepo) CreateMessage(authorID int, content string) (*models.ChatMessage, error) { args := m.Called(authorID, content); return args.Get(0).(*models.ChatMessage), args.Error(1) }
func (m *MockChatRepo) GetAllMessages() ([]*models.ChatMessage, error) { args := m.Called(); return args.Get(0).([]*models.ChatMessage), args.Error(1) }
func (m *MockChatRepo) CleanOldMessages() error { args := m.Called(); return args.Error(0) }
func (m *MockChatRepo) Cleanup() error { args := m.Called(); return args.Error(0) }
func (m *MockChatRepo) DeleteOldMessages() error { args := m.Called(); return args.Error(0) } 