package service

import (
	"errors"
	"testing"

	"ForumService/internal/models"
	"ForumService/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetThreadWithPosts_Success(t *testing.T) {
	threadRepo := new(mocks.MockThreadRepo)
	postRepo := new(mocks.MockPostRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewThreadService(threadRepo, postRepo, userRepo)

	thread := &models.Thread{ID: 1, Title: "Test", AuthorID: 1}
	posts := []*models.Post{{ID: 1, ThreadID: 1, AuthorID: 1, Content: "post"}}
	threadRepo.On("GetByID", 1).Return(thread, nil)
	postRepo.On("GetByThreadID", 1).Return(posts, nil)

	resThread, resPosts, err := service.GetThreadWithPosts(1)
	assert.NoError(t, err)
	assert.Equal(t, thread, resThread)
	assert.Equal(t, posts, resPosts)
}

func TestGetThreadWithPosts_ThreadNotFound(t *testing.T) {
	threadRepo := new(mocks.MockThreadRepo)
	postRepo := new(mocks.MockPostRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewThreadService(threadRepo, postRepo, userRepo)

	threadRepo.On("GetByID", 2).Return((*models.Thread)(nil), errors.New("not found"))

	resThread, resPosts, err := service.GetThreadWithPosts(2)
	assert.Error(t, err)
	assert.Nil(t, resThread)
	assert.Nil(t, resPosts)
}

func TestCreateThread_Success(t *testing.T) {
	threadRepo := new(mocks.MockThreadRepo)
	postRepo := new(mocks.MockPostRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewThreadService(threadRepo, postRepo, userRepo)

	threadRepo.On("Create", mock.AnythingOfType("*models.Thread")).Return(nil)
	thread, err := service.CreateThread("title", 1)
	assert.NoError(t, err)
	assert.Equal(t, "title", thread.Title)
	assert.Equal(t, 1, thread.AuthorID)
}

func TestUpdateThread_NoPermission(t *testing.T) {
	threadRepo := new(mocks.MockThreadRepo)
	postRepo := new(mocks.MockPostRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewThreadService(threadRepo, postRepo, userRepo)

	thread := &models.Thread{ID: 1, AuthorID: 2}
	threadRepo.On("GetByID", 1).Return(thread, nil)
	userRepo.On("GetUserRole", 3).Return("user", nil)

	err := service.UpdateThread(&models.Thread{ID: 1}, 3)
	assert.ErrorIs(t, err, ErrNoPermission)
}

func TestUpdateThread_Admin(t *testing.T) {
	threadRepo := new(mocks.MockThreadRepo)
	postRepo := new(mocks.MockPostRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewThreadService(threadRepo, postRepo, userRepo)

	thread := &models.Thread{ID: 1, AuthorID: 2}
	threadRepo.On("GetByID", 1).Return(thread, nil)
	userRepo.On("GetUserRole", 4).Return("admin", nil)
	threadRepo.On("Update", mock.AnythingOfType("*models.Thread")).Return(nil)

	err := service.UpdateThread(&models.Thread{ID: 1}, 4)
	assert.NoError(t, err)
}

func TestDeleteThread_NoPermission(t *testing.T) {
	threadRepo := new(mocks.MockThreadRepo)
	postRepo := new(mocks.MockPostRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewThreadService(threadRepo, postRepo, userRepo)

	thread := &models.Thread{ID: 1, AuthorID: 2}
	threadRepo.On("GetByID", 1).Return(thread, nil)
	userRepo.On("GetUserRole", 3).Return("user", nil)

	err := service.DeleteThread(1, 3)
	assert.ErrorIs(t, err, ErrNoPermission)
}

func TestDeleteThread_Admin(t *testing.T) {
	threadRepo := new(mocks.MockThreadRepo)
	postRepo := new(mocks.MockPostRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewThreadService(threadRepo, postRepo, userRepo)

	thread := &models.Thread{ID: 1, AuthorID: 2}
	threadRepo.On("GetByID", 1).Return(thread, nil)
	userRepo.On("GetUserRole", 4).Return("admin", nil)
	threadRepo.On("Delete", 1).Return(nil)

	err := service.DeleteThread(1, 4)
	assert.NoError(t, err)
}

func TestGetAllThreads_Error(t *testing.T) {
	threadRepo := new(mocks.MockThreadRepo)
	postRepo := new(mocks.MockPostRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewThreadService(threadRepo, postRepo, userRepo)

	threadRepo.On("GetAllThreads").Return(([]*models.Thread)(nil), errors.New("fail"))
	threads, err := service.GetAllThreads()
	assert.Error(t, err)
	assert.Nil(t, threads)
}

func TestGetAllThreads_Success(t *testing.T) {
	threadRepo := new(mocks.MockThreadRepo)
	postRepo := new(mocks.MockPostRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewThreadService(threadRepo, postRepo, userRepo)

	threads := []*models.Thread{
		{ID: 1, Title: "Thread 1", AuthorID: 1},
		{ID: 2, Title: "Thread 2", AuthorID: 2},
	}
	threadRepo.On("GetAllThreads").Return(threads, nil)
	res, err := service.GetAllThreads()
	assert.NoError(t, err)
	assert.Equal(t, threads, res)
}

func TestGetPostsByThreadID(t *testing.T) {
	threadRepo := new(mocks.MockThreadRepo)
	postRepo := new(mocks.MockPostRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewThreadService(threadRepo, postRepo, userRepo)

	posts := []*models.Post{{ID: 1, ThreadID: 1, AuthorID: 1, Content: "post"}}
	postRepo.On("GetByThreadID", 1).Return(posts, nil)
	res, err := service.GetPostsByThreadID(1)
	assert.NoError(t, err)
	assert.Equal(t, posts, res)
}

func TestGetUserByID(t *testing.T) {
	threadRepo := new(mocks.MockThreadRepo)
	postRepo := new(mocks.MockPostRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewThreadService(threadRepo, postRepo, userRepo)

	user := &models.User{ID: 1, Username: "test"}
	userRepo.On("GetUserByID", 1).Return(user, nil)
	res, err := service.GetUserByID(1)
	assert.NoError(t, err)
	assert.Equal(t, user, res)
}

func TestGetThreadWithPosts_WithComments(t *testing.T) {
	threadRepo := new(mocks.MockThreadRepo)
	postRepo := new(mocks.MockPostRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewThreadService(threadRepo, postRepo, userRepo)

	thread := &models.Thread{ID: 1, Title: "thread"}
	posts := []*models.Post{{ID: 1, ThreadID: 1, AuthorID: 1, Content: "post"}}
	threadRepo.On("GetByID", 1).Return(thread, nil)
	postRepo.On("GetByThreadID", 1).Return(posts, nil)

	resThread, resPosts, err := service.GetThreadWithPosts(1)
	assert.NoError(t, err)
	assert.Equal(t, thread, resThread)
	assert.Equal(t, posts, resPosts)
}

func TestGetThreadWithPosts_Error(t *testing.T) {
	threadRepo := new(mocks.MockThreadRepo)
	postRepo := new(mocks.MockPostRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewThreadService(threadRepo, postRepo, userRepo)

	threadRepo.On("GetByID", 1).Return((*models.Thread)(nil), errors.New("db error"))

	resThread, resPosts, err := service.GetThreadWithPosts(1)
	assert.Error(t, err)
	assert.Nil(t, resThread)
	assert.Nil(t, resPosts)
} 