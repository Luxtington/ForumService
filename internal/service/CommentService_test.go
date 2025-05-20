package service

import (
	"errors"
	"testing"
	"ForumService/internal/models"
	"ForumService/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateComment_Success(t *testing.T) {
	repo := new(mocks.MockCommentRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewCommentService(repo, userRepo)

	comment := &models.Comment{PostID: 1, AuthorID: 2, Content: "content"}
	repo.On("SaveComment", comment).Return(nil)
	userRepo.On("GetUserRole", 2).Return("user", nil)
	res, err := service.CreateComment(1, 2, "content")
	assert.NoError(t, err)
	assert.Equal(t, 1, res.PostID)
	assert.Equal(t, 2, res.AuthorID)
	assert.Equal(t, "content", res.Content)
}

func TestGetCommentByID_Success(t *testing.T) {
	repo := new(mocks.MockCommentRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewCommentService(repo, userRepo)

	comment := &models.Comment{ID: 1}
	repo.On("GetCommentByID", 1).Return(comment, nil)
	res, err := service.GetCommentByID(1)
	assert.NoError(t, err)
	assert.Equal(t, comment, res)
}

func TestGetCommentsByPostID_Success(t *testing.T) {
	repo := new(mocks.MockCommentRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewCommentService(repo, userRepo)

	comments := []models.Comment{{ID: 1, PostID: 1}}
	repo.On("GetCommentsByPostID", 1).Return(comments, nil)
	res, err := service.GetCommentsByPostID(1)
	assert.NoError(t, err)
	assert.Equal(t, comments, res)
}

func TestDeleteComment_NoPermission(t *testing.T) {
	repo := new(mocks.MockCommentRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewCommentService(repo, userRepo)

	comment := &models.Comment{ID: 1, AuthorID: 2}
	repo.On("GetCommentByID", 1).Return(comment, nil)
	userRepo.On("GetUserRole", 3).Return("user", nil)

	err := service.DeleteComment(1, 3)
	assert.ErrorIs(t, err, ErrNoPermission)
}

func TestDeleteComment_Admin(t *testing.T) {
	repo := new(mocks.MockCommentRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewCommentService(repo, userRepo)

	comment := &models.Comment{ID: 1, AuthorID: 2}
	repo.On("GetCommentByID", 1).Return(comment, nil)
	userRepo.On("GetUserRole", 4).Return("admin", nil)
	repo.On("DeleteComment", 1).Return(nil)

	err := service.DeleteComment(1, 4)
	assert.NoError(t, err)
}

func TestCreateComment_Error(t *testing.T) {
	repo := new(mocks.MockCommentRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewCommentService(repo, userRepo)

	comment := &models.Comment{PostID: 1, AuthorID: 2, Content: "fail"}
	repo.On("SaveComment", comment).Return(errors.New("db error"))
	userRepo.On("GetUserRole", 2).Return("user", nil)
	res, err := service.CreateComment(1, 2, "fail")
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestGetCommentByID_Error(t *testing.T) {
	repo := new(mocks.MockCommentRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewCommentService(repo, userRepo)

	repo.On("GetCommentByID", 1).Return((*models.Comment)(nil), errors.New("db error"))
	res, err := service.GetCommentByID(1)
	assert.Error(t, err)
	assert.Nil(t, res)
} 