package service

import (
	"errors"
	"testing"
	"ForumService/internal/models"
	"ForumService/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByID_Success(t *testing.T) {
	repo := new(mocks.MockUserRepo)
	service := NewUserService(repo)
	user := &models.User{ID: 1, Username: "test"}
	repo.On("GetUserByID", 1).Return(user, nil)
	res, err := service.GetUserByID(1)
	assert.NoError(t, err)
	assert.Equal(t, user, res)
}

func TestGetUserPosts_Success(t *testing.T) {
	repo := new(mocks.MockUserRepo)
	service := NewUserService(repo)
	posts := []*models.Post{{ID: 1, AuthorID: 1, Content: "post"}}
	repo.On("GetUserPosts", 1).Return(posts, nil)
	res, err := service.GetUserPosts(1)
	assert.NoError(t, err)
	assert.Equal(t, posts, res)
}

func TestGetUserCommentCount_Success(t *testing.T) {
	repo := new(mocks.MockUserRepo)
	service := NewUserService(repo)
	repo.On("GetUserCommentCount", 1).Return(5, nil)
	count, err := service.GetUserCommentCount(1)
	assert.NoError(t, err)
	assert.Equal(t, 5, count)
}

func TestGetUserByID_Error(t *testing.T) {
	repo := new(mocks.MockUserRepo)
	service := NewUserService(repo)

	repo.On("GetUserByID", 1).Return((*models.User)(nil), errors.New("db error"))
	res, err := service.GetUserByID(1)
	assert.Error(t, err)
	assert.Nil(t, res)
} 