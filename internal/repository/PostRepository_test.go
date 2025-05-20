package repository

import (
	_"database/sql"
	"ForumService/internal/models"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

func setupPostRepositoryTest(t *testing.T) (*postRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewPostRepository(db).(*postRepository)

	cleanup := func() {
		db.Close()
	}

	return repo, mock, cleanup
}

func TestPostRepository_SavePost(t *testing.T) {
	repo, mock, cleanup := setupPostRepositoryTest(t)
	defer cleanup()

	post := &models.Post{
		ThreadID: 1,
		AuthorID: 1,
		Content:  "Test Post",
	}

	mock.ExpectQuery("INSERT INTO posts").
		WithArgs(post.ThreadID, post.AuthorID, post.Content).
		WillReturnRows(sqlmock.NewRows([]string{"id", "thread_id", "author_id", "content", "created_at"}).
			AddRow(1, post.ThreadID, post.AuthorID, post.Content, time.Now()))

	err := repo.SavePost(post)
	require.NoError(t, err)
	assert.Equal(t, 1, post.ID)
	assert.Equal(t, 1, post.ThreadID)
	assert.Equal(t, 1, post.AuthorID)
	assert.Equal(t, "Test Post", post.Content)
}

func TestPostRepository_GetPostByID(t *testing.T) {
	repo, mock, cleanup := setupPostRepositoryTest(t)
	defer cleanup()

	expectedPost := &models.Post{
		ID:        1,
		ThreadID:  1,
		AuthorID:  1,
		Content:   "Test Post",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		AuthorName: "Test User",
	}

	mock.ExpectQuery("SELECT p.id, p.thread_id, p.author_id, p.content, p.created_at, p.updated_at, u.username as author_name FROM posts p LEFT JOIN users u ON p.author_id = u.id WHERE p.id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "thread_id", "author_id", "content", "created_at", "updated_at", "author_name"}).
			AddRow(expectedPost.ID, expectedPost.ThreadID, expectedPost.AuthorID, expectedPost.Content, expectedPost.CreatedAt, expectedPost.UpdatedAt, expectedPost.AuthorName))

	post, err := repo.GetPostByID(1)
	require.NoError(t, err)
	assert.Equal(t, expectedPost.ID, post.ID)
	assert.Equal(t, expectedPost.ThreadID, post.ThreadID)
	assert.Equal(t, expectedPost.AuthorID, post.AuthorID)
	assert.Equal(t, expectedPost.Content, post.Content)
	assert.Equal(t, expectedPost.AuthorName, post.AuthorName)
}

func TestPostRepository_GetPostWithComments(t *testing.T) {
	repo, mock, cleanup := setupPostRepositoryTest(t)
	defer cleanup()

	// Мок для получения поста
	expectedPost := &models.Post{
		ID:        1,
		ThreadID:  1,
		AuthorID:  1,
		Content:   "Test Post",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		AuthorName: "Test User",
	}

	mock.ExpectQuery("SELECT p.id, p.thread_id, p.author_id, p.content, p.created_at, p.updated_at, u.username as author_name FROM posts p LEFT JOIN users u ON p.author_id = u.id WHERE p.id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "thread_id", "author_id", "content", "created_at", "updated_at", "author_name"}).
			AddRow(expectedPost.ID, expectedPost.ThreadID, expectedPost.AuthorID, expectedPost.Content, expectedPost.CreatedAt, expectedPost.UpdatedAt, expectedPost.AuthorName))

	// Мок для получения комментариев
	expectedComments := []models.Comment{
		{
			ID:        1,
			PostID:    1,
			AuthorID:  1,
			Content:   "Test Comment",
			CreatedAt: time.Now(),
			AuthorName: "Test User",
		},
	}

	mock.ExpectQuery("SELECT c.id, c.post_id, c.author_id, c.content, c.created_at, u.username as author_name FROM comments c LEFT JOIN users u ON c.author_id = u.id WHERE c.post_id = \\$1 ORDER BY c.created_at ASC").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "author_id", "content", "created_at", "author_name"}).
			AddRow(expectedComments[0].ID, expectedComments[0].PostID, expectedComments[0].AuthorID, expectedComments[0].Content, expectedComments[0].CreatedAt, expectedComments[0].AuthorName))

	post, comments, err := repo.GetPostWithComments(1)
	require.NoError(t, err)
	assert.Equal(t, expectedPost.ID, post.ID)
	assert.Equal(t, expectedPost.Content, post.Content)
	assert.Equal(t, len(expectedComments), len(comments))
	assert.Equal(t, expectedComments[0].Content, comments[0].Content)
}

func TestPostRepository_GetPostsWithCommentsByThreadID(t *testing.T) {
	repo, mock, cleanup := setupPostRepositoryTest(t)
	defer cleanup()

	// Мок для получения постов
	expectedPosts := []models.Post{
		{
			ID:        1,
			ThreadID:  1,
			Content:   "Test Post 1",
			CreatedAt: time.Now(),
			AuthorID:  1,
			AuthorName: "Test User 1",
		},
		{
			ID:        2,
			ThreadID:  1,
			Content:   "Test Post 2",
			CreatedAt: time.Now(),
			AuthorID:  2,
			AuthorName: "Test User 2",
		},
	}

	postRows := sqlmock.NewRows([]string{"id", "thread_id", "content", "created_at", "author_id", "author_username"})
	postRows.AddRow(1, 1, "Test Post 1", expectedPosts[0].CreatedAt, 1, "Test User 1")
	postRows.AddRow(2, 1, "Test Post 2", expectedPosts[1].CreatedAt, 2, "Test User 2")

	mock.ExpectQuery("SELECT p.id, p.thread_id, p.content, p.created_at, u.id as author_id, u.username as author_username FROM posts p JOIN users u ON p.user_id = u.id WHERE p.thread_id = \\$1 ORDER BY p.created_at DESC LIMIT \\$2 OFFSET \\$3").
		WithArgs(1, 20, 0).
		WillReturnRows(postRows)

	// Мок для получения комментариев
	commentRows := sqlmock.NewRows([]string{"id", "post_id", "author_id", "content", "created_at"})
	commentRows.AddRow(1, 1, 1, "Test Comment 1", time.Now())
	commentRows.AddRow(2, 2, 2, "Test Comment 2", time.Now())

	mock.ExpectQuery("SELECT c.id, c.post_id, c.author_id, c.content, c.created_at FROM comments c WHERE c.post_id = ANY\\(\\$1\\) ORDER BY c.created_at ASC").
		WithArgs(pq.Array([]int{1, 2})).
		WillReturnRows(commentRows)

	posts, comments, err := repo.GetPostsWithCommentsByThreadID(1)
	require.NoError(t, err)
	assert.Equal(t, 2, len(posts))
	assert.Equal(t, 2, len(comments))
	
	// Проверяем, что все поля правильно заполнены
	for i, post := range posts {
		assert.Equal(t, expectedPosts[i].ID, post.ID)
		assert.Equal(t, expectedPosts[i].ThreadID, post.ThreadID)
		assert.Equal(t, expectedPosts[i].Content, post.Content)
		assert.Equal(t, expectedPosts[i].AuthorID, post.AuthorID)
		assert.Equal(t, expectedPosts[i].AuthorName, post.AuthorName)
	}
}

func TestPostRepository_UpdatePost(t *testing.T) {
	repo, mock, cleanup := setupPostRepositoryTest(t)
	defer cleanup()

	post := &models.Post{
		ID:      1,
		Content: "Updated Post",
	}

	mock.ExpectExec("UPDATE posts SET content = \\$1, updated_at = CURRENT_TIMESTAMP WHERE id = \\$2").
		WithArgs(post.Content, post.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdatePost(post, post.ID)
	require.NoError(t, err)
}

func TestPostRepository_DeletePost(t *testing.T) {
	repo, mock, cleanup := setupPostRepositoryTest(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM comments WHERE post_id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("DELETE FROM posts WHERE id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.DeletePost(1)
	require.NoError(t, err)
}

func TestPostRepository_GetByThreadID(t *testing.T) {
	repo, mock, cleanup := setupPostRepositoryTest(t)
	defer cleanup()

	expectedPosts := []*models.Post{
		{
			ID:        1,
			ThreadID:  1,
			AuthorID:  1,
			Content:   "Test Post 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			AuthorName: "Test User 1",
		},
		{
			ID:        2,
			ThreadID:  1,
			AuthorID:  2,
			Content:   "Test Post 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			AuthorName: "Test User 2",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "thread_id", "author_id", "content", "created_at", "updated_at", "author_name"})
	for _, post := range expectedPosts {
		rows.AddRow(post.ID, post.ThreadID, post.AuthorID, post.Content, post.CreatedAt, post.UpdatedAt, post.AuthorName)
	}

	mock.ExpectQuery("SELECT p.id, p.thread_id, p.author_id, p.content, p.created_at, p.updated_at, u.username as author_name FROM posts p LEFT JOIN users u ON p.author_id = u.id WHERE p.thread_id = \\$1 ORDER BY p.created_at ASC").
		WithArgs(1).
		WillReturnRows(rows)

	posts, err := repo.GetByThreadID(1)
	require.NoError(t, err)
	assert.Equal(t, len(expectedPosts), len(posts))
	for i, post := range posts {
		assert.Equal(t, expectedPosts[i].ID, post.ID)
		assert.Equal(t, expectedPosts[i].Content, post.Content)
		assert.Equal(t, expectedPosts[i].AuthorName, post.AuthorName)
	}
} 