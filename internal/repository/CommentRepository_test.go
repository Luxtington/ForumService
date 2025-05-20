package repository

import (
	_"database/sql"
	"ForumService/internal/models"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/DATA-DOG/go-sqlmock"
)

func setupCommentRepositoryTest(t *testing.T) (*CommentRepositoryImpl, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewCommentRepository(db).(*CommentRepositoryImpl)

	cleanup := func() {
		db.Close()
	}

	return repo, mock, cleanup
}

func TestCommentRepository_SaveComment(t *testing.T) {
	repo, mock, cleanup := setupCommentRepositoryTest(t)
	defer cleanup()

	comment := &models.Comment{
		PostID:   1,
		AuthorID: 1,
		Content:  "Test Comment",
	}

	mock.ExpectQuery("INSERT INTO comments").
		WithArgs(comment.PostID, comment.AuthorID, comment.Content).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err := repo.SaveComment(comment)
	require.NoError(t, err)
	assert.Equal(t, 1, comment.ID)
}

func TestCommentRepository_GetCommentByID(t *testing.T) {
	repo, mock, cleanup := setupCommentRepositoryTest(t)
	defer cleanup()

	expectedComment := &models.Comment{
		ID:        1,
		PostID:    1,
		AuthorID:  1,
		Content:   "Test Comment",
		CreatedAt: time.Now(),
	}

	mock.ExpectQuery("SELECT id, post_id, author_id, content, created_at FROM comments WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "author_id", "content", "created_at"}).
			AddRow(expectedComment.ID, expectedComment.PostID, expectedComment.AuthorID, expectedComment.Content, expectedComment.CreatedAt))

	comment, err := repo.GetCommentByID(1)
	require.NoError(t, err)
	assert.Equal(t, expectedComment.ID, comment.ID)
	assert.Equal(t, expectedComment.PostID, comment.PostID)
	assert.Equal(t, expectedComment.AuthorID, comment.AuthorID)
	assert.Equal(t, expectedComment.Content, comment.Content)
}

func TestCommentRepository_DeleteComment(t *testing.T) {
	repo, mock, cleanup := setupCommentRepositoryTest(t)
	defer cleanup()

	mock.ExpectExec("DELETE FROM comments WHERE id = \\$1 RETURNING id").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteComment(1)
	require.NoError(t, err)
}

func TestCommentRepository_GetCommentsByPostID(t *testing.T) {
	repo, mock, cleanup := setupCommentRepositoryTest(t)
	defer cleanup()

	expectedComments := []models.Comment{
		{
			ID:        1,
			PostID:    1,
			AuthorID:  1,
			Content:   "Test Comment 1",
			CreatedAt: time.Now(),
		},
		{
			ID:        2,
			PostID:    1,
			AuthorID:  2,
			Content:   "Test Comment 2",
			CreatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "post_id", "author_id", "content", "created_at"})
	for _, comment := range expectedComments {
		rows.AddRow(comment.ID, comment.PostID, comment.AuthorID, comment.Content, comment.CreatedAt)
	}

	mock.ExpectQuery("SELECT id, post_id, author_id, content, created_at FROM comments WHERE post_id = \\$1 ORDER BY created_at ASC").
		WithArgs(1).
		WillReturnRows(rows)

	comments, err := repo.GetCommentsByPostID(1)
	require.NoError(t, err)
	assert.Equal(t, len(expectedComments), len(comments))
	for i, comment := range comments {
		assert.Equal(t, expectedComments[i].ID, comment.ID)
		assert.Equal(t, expectedComments[i].Content, comment.Content)
	}
} 