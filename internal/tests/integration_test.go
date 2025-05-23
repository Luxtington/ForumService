package tests

import (
	"ForumService/internal/repository"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func setupTest(t *testing.T) (*sql.DB, repository.UserRepository, repository.PostRepository, repository.CommentRepository, func()) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=forum sslmode=disable")
	require.NoError(t, err)

	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)

	cleanup := func() {
		db.Close()
	}

	return db, userRepo, postRepo, commentRepo, cleanup
}

func TestGetNonExistingUser(t *testing.T) {
	_, userRepo, _, _, cleanup := setupTest(t)
	defer cleanup()

	user, err := userRepo.GetUserByID(-1)
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestGetNonExistingPost(t *testing.T) {
	_, _, postRepo, _, cleanup := setupTest(t)
	defer cleanup()

	post, err := postRepo.GetPostByID(-1)
	assert.Error(t, err)
	assert.Nil(t, post)
}

func TestGetNonExistingComment(t *testing.T) {
	_, _, _, commentRepo, cleanup := setupTest(t)
	defer cleanup()

	comment, err := commentRepo.GetCommentByID(-1)
	assert.Error(t, err)
	assert.Nil(t, comment)
}

func TestGetUserPostsEmpty(t *testing.T) {
	_, userRepo, _, _, cleanup := setupTest(t)
	defer cleanup()

	posts, err := userRepo.GetUserPosts(-1)
	assert.NoError(t, err)
	assert.True(t, posts == nil || len(posts) == 0)
}

func TestGetUserCommentCountZero(t *testing.T) {
	_, userRepo, _, _, cleanup := setupTest(t)
	defer cleanup()

	count, err := userRepo.GetUserCommentCount(-1)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestGetNonExistingThread(t *testing.T) {
	db, _, _, _, cleanup := setupTest(t)
	defer cleanup()

	threadRepo := repository.NewThreadRepository(db)
	thread, err := threadRepo.GetByID(-1)
	assert.NoError(t, err)
	assert.Nil(t, thread)
}

func TestGetPostsByNonExistingThreadID(t *testing.T) {
	_, _, postRepo, _, cleanup := setupTest(t)
	defer cleanup()

	posts, err := postRepo.GetByThreadID(-1)
	assert.NoError(t, err)
	assert.True(t, posts == nil || len(posts) == 0)
}

func TestGetCommentsByNonExistingPostID(t *testing.T) {
	_, _, _, commentRepo, cleanup := setupTest(t)
	defer cleanup()

	comments, err := commentRepo.GetCommentsByPostID(-1)
	assert.NoError(t, err)
	assert.True(t, comments == nil || len(comments) == 0)
} 