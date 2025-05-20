package service

import (
	"errors"
	"testing"
	"ForumService/internal/models"
	"ForumService/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreatePost_Success(t *testing.T) {
	repo := new(mocks.MockPostRepo)
	commentRepo := new(mocks.MockCommentRepo)
	threadRepo := new(mocks.MockThreadRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewPostService(repo, commentRepo, threadRepo, userRepo)

	post := &models.Post{ID: 1, Content: "test"}
	repo.On("SavePost", post).Return(nil)
	err := service.CreatePost(post)
	assert.NoError(t, err)
}

func TestGetPostByID_Success(t *testing.T) {
	repo := new(mocks.MockPostRepo)
	commentRepo := new(mocks.MockCommentRepo)
	threadRepo := new(mocks.MockThreadRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewPostService(repo, commentRepo, threadRepo, userRepo)

	post := &models.Post{ID: 1}
	repo.On("GetPostByID", 1).Return(post, nil)
	res, err := service.GetPostByID(1)
	assert.NoError(t, err)
	assert.Equal(t, post, res)
}

func TestUpdatePost_NoPermission(t *testing.T) {
	repo := new(mocks.MockPostRepo)
	commentRepo := new(mocks.MockCommentRepo)
	threadRepo := new(mocks.MockThreadRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewPostService(repo, commentRepo, threadRepo, userRepo)

	post := &models.Post{ID: 1, AuthorID: 2}
	repo.On("GetPostByID", 1).Return(post, nil)
	userRepo.On("GetUserRole", 3).Return("user", nil)

	err := service.UpdatePost(&models.Post{ID: 1}, 1, 3)
	assert.ErrorIs(t, err, ErrNoPermission)
}

func TestUpdatePost_Admin(t *testing.T) {
	repo := new(mocks.MockPostRepo)
	commentRepo := new(mocks.MockCommentRepo)
	threadRepo := new(mocks.MockThreadRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewPostService(repo, commentRepo, threadRepo, userRepo)

	post := &models.Post{ID: 1, AuthorID: 2}
	repo.On("GetPostByID", 1).Return(post, nil)
	userRepo.On("GetUserRole", 4).Return("admin", nil)
	repo.On("UpdatePost", mock.AnythingOfType("*models.Post"), 1).Return(nil)

	err := service.UpdatePost(&models.Post{ID: 1}, 1, 4)
	assert.NoError(t, err)
}

func TestDeletePost_NoPermission(t *testing.T) {
	repo := new(mocks.MockPostRepo)
	commentRepo := new(mocks.MockCommentRepo)
	threadRepo := new(mocks.MockThreadRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewPostService(repo, commentRepo, threadRepo, userRepo)

	post := &models.Post{ID: 1, AuthorID: 2}
	repo.On("GetPostByID", 1).Return(post, nil)
	userRepo.On("GetUserRole", 3).Return("user", nil)

	err := service.DeletePost(1, 3)
	assert.ErrorIs(t, err, ErrNoPermission)
}

func TestDeletePost_Admin(t *testing.T) {
	repo := new(mocks.MockPostRepo)
	commentRepo := new(mocks.MockCommentRepo)
	threadRepo := new(mocks.MockThreadRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewPostService(repo, commentRepo, threadRepo, userRepo)

	post := &models.Post{ID: 1, AuthorID: 2}
	repo.On("GetPostByID", 1).Return(post, nil)
	userRepo.On("GetUserRole", 4).Return("admin", nil)
	repo.On("DeletePost", 1).Return(nil)

	err := service.DeletePost(1, 4)
	assert.NoError(t, err)
}

func TestCreatePost_Error(t *testing.T) {
	repo := new(mocks.MockPostRepo)
	commentRepo := new(mocks.MockCommentRepo)
	threadRepo := new(mocks.MockThreadRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewPostService(repo, commentRepo, threadRepo, userRepo)

	post := &models.Post{ID: 1, Content: "test"}
	repo.On("SavePost", post).Return(errors.New("db error"))
	err := service.CreatePost(post)
	assert.Error(t, err)
}

func TestGetPostByID_Error(t *testing.T) {
	repo := new(mocks.MockPostRepo)
	commentRepo := new(mocks.MockCommentRepo)
	threadRepo := new(mocks.MockThreadRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewPostService(repo, commentRepo, threadRepo, userRepo)

	repo.On("GetPostByID", 1).Return((*models.Post)(nil), errors.New("db error"))
	res, err := service.GetPostByID(1)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestUpdatePost_Error(t *testing.T) {
	repo := new(mocks.MockPostRepo)
	commentRepo := new(mocks.MockCommentRepo)
	threadRepo := new(mocks.MockThreadRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewPostService(repo, commentRepo, threadRepo, userRepo)

	post := &models.Post{ID: 1, AuthorID: 1}
	repo.On("GetPostByID", 1).Return(post, nil)
	userRepo.On("GetUserRole", 1).Return("user", nil)
	repo.On("UpdatePost", mock.AnythingOfType("*models.Post"), 1).Return(errors.New("db error"))

	err := service.UpdatePost(&models.Post{ID: 1}, 1, 1)
	assert.Error(t, err)
}

func TestDeletePost_Error(t *testing.T) {
	repo := new(mocks.MockPostRepo)
	commentRepo := new(mocks.MockCommentRepo)
	threadRepo := new(mocks.MockThreadRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewPostService(repo, commentRepo, threadRepo, userRepo)

	post := &models.Post{ID: 1, AuthorID: 1}
	repo.On("GetPostByID", 1).Return(post, nil)
	userRepo.On("GetUserRole", 1).Return("user", nil)
	repo.On("DeletePost", 1).Return(errors.New("db error"))

	err := service.DeletePost(1, 1)
	assert.Error(t, err)
}

func TestGetPostWithComments_Success(t *testing.T) {
	repo := new(mocks.MockPostRepo)
	commentRepo := new(mocks.MockCommentRepo)
	threadRepo := new(mocks.MockThreadRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewPostService(repo, commentRepo, threadRepo, userRepo)

	post := &models.Post{ID: 1, ThreadID: 1, AuthorID: 1, Content: "post"}
	comments := []models.Comment{{ID: 1, PostID: 1, AuthorID: 1, Content: "comment"}}
	repo.On("GetPostWithComments", 1).Return(post, comments, nil)

	resPost, resComments, err := service.GetPostWithComments(1)
	assert.NoError(t, err)
	assert.Equal(t, post, resPost)
	assert.Equal(t, comments, resComments)
}

func TestGetPostWithComments_Error(t *testing.T) {
	repo := new(mocks.MockPostRepo)
	commentRepo := new(mocks.MockCommentRepo)
	threadRepo := new(mocks.MockThreadRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewPostService(repo, commentRepo, threadRepo, userRepo)

	repo.On("GetPostWithComments", 1).Return((*models.Post)(nil), ([]models.Comment)(nil), errors.New("db error"))
	resPost, resComments, err := service.GetPostWithComments(1)
	assert.Error(t, err)
	assert.Nil(t, resPost)
	assert.Nil(t, resComments)
}

func TestGetPostsWithCommentsByThreadID_Success(t *testing.T) {
	repo := new(mocks.MockPostRepo)
	commentRepo := new(mocks.MockCommentRepo)
	threadRepo := new(mocks.MockThreadRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewPostService(repo, commentRepo, threadRepo, userRepo)

	posts := []models.Post{{ID: 1, ThreadID: 1, AuthorID: 1, Content: "post"}}
	comments := map[int][]models.Comment{
		0: {{ID: 1, PostID: 1, AuthorID: 1, Content: "comment"}},
	}
	repo.On("GetPostsWithCommentsByThreadID", 1).Return(posts, comments, nil)

	resPosts, resComments, err := service.GetPostsWithCommentsByThreadID(1)
	assert.NoError(t, err)
	assert.Equal(t, posts, resPosts)
	assert.Equal(t, comments, resComments)
}

func TestGetPostsWithCommentsByThreadID_Error(t *testing.T) {
	repo := new(mocks.MockPostRepo)
	commentRepo := new(mocks.MockCommentRepo)
	threadRepo := new(mocks.MockThreadRepo)
	userRepo := new(mocks.MockUserRepo)
	service := NewPostService(repo, commentRepo, threadRepo, userRepo)

	repo.On("GetPostsWithCommentsByThreadID", 1).Return(([]models.Post)(nil), (map[int][]models.Comment)(nil), errors.New("db error"))
	resPosts, resComments, err := service.GetPostsWithCommentsByThreadID(1)
	assert.Error(t, err)
	assert.Nil(t, resPosts)
	assert.Nil(t, resComments)
} 