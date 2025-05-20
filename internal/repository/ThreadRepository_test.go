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

func setupThreadRepositoryTest(t *testing.T) (*threadRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewThreadRepository(db).(*threadRepository)

	cleanup := func() {
		db.Close()
	}

	return repo, mock, cleanup
}

func TestThreadRepository_Create(t *testing.T) {
	repo, mock, cleanup := setupThreadRepositoryTest(t)
	defer cleanup()

	thread := &models.Thread{
		Title:    "Test Thread",
		AuthorID: 1,
	}

	mock.ExpectQuery("INSERT INTO threads").
		WithArgs(thread.Title, thread.AuthorID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author_id", "created_at", "updated_at"}).
			AddRow(1, thread.Title, thread.AuthorID, time.Now(), time.Now()))

	err := repo.Create(thread)
	require.NoError(t, err)
	assert.Equal(t, 1, thread.ID)
	assert.Equal(t, "Test Thread", thread.Title)
	assert.Equal(t, 1, thread.AuthorID)
}

func TestThreadRepository_GetByID(t *testing.T) {
	repo, mock, cleanup := setupThreadRepositoryTest(t)
	defer cleanup()

	expectedThread := &models.Thread{
		ID:        1,
		Title:     "Test Thread",
		AuthorID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectQuery("SELECT id, title, author_id, created_at, updated_at FROM threads WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author_id", "created_at", "updated_at"}).
			AddRow(expectedThread.ID, expectedThread.Title, expectedThread.AuthorID, expectedThread.CreatedAt, expectedThread.UpdatedAt))

	thread, err := repo.GetByID(1)
	require.NoError(t, err)
	assert.Equal(t, expectedThread.ID, thread.ID)
	assert.Equal(t, expectedThread.Title, thread.Title)
	assert.Equal(t, expectedThread.AuthorID, thread.AuthorID)
}

func TestThreadRepository_Update(t *testing.T) {
	repo, mock, cleanup := setupThreadRepositoryTest(t)
	defer cleanup()

	thread := &models.Thread{
		ID:    1,
		Title: "Updated Thread",
	}

	mock.ExpectExec("UPDATE threads SET title = \\$1, updated_at = CURRENT_TIMESTAMP WHERE id = \\$2").
		WithArgs(thread.Title, thread.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(thread)
	require.NoError(t, err)
}

func TestThreadRepository_Delete(t *testing.T) {
	repo, mock, cleanup := setupThreadRepositoryTest(t)
	defer cleanup()

	mock.ExpectExec("DELETE FROM threads WHERE id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Delete(1)
	require.NoError(t, err)
}

func TestThreadRepository_GetAllThreads(t *testing.T) {
	repo, mock, cleanup := setupThreadRepositoryTest(t)
	defer cleanup()

	expectedThreads := []*models.Thread{
		{
			ID:        1,
			Title:     "Thread 1",
			AuthorID:  1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			AuthorName: "Test User 1",
		},
		{
			ID:        2,
			Title:     "Thread 2",
			AuthorID:  2,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			AuthorName: "Test User 2",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "author_id", "created_at", "updated_at", "author_name"})
	for _, thread := range expectedThreads {
		rows.AddRow(thread.ID, thread.Title, thread.AuthorID, thread.CreatedAt, thread.UpdatedAt, thread.AuthorName)
	}

	mock.ExpectQuery("SELECT t.id, t.title, t.author_id, t.created_at, t.updated_at, u.username as author_name FROM threads t LEFT JOIN users u ON t.author_id = u.id ORDER BY t.created_at DESC").
		WillReturnRows(rows)

	threads, err := repo.GetAllThreads()
	require.NoError(t, err)
	assert.Equal(t, len(expectedThreads), len(threads))
	for i, thread := range threads {
		assert.Equal(t, expectedThreads[i].ID, thread.ID)
		assert.Equal(t, expectedThreads[i].Title, thread.Title)
		assert.Equal(t, expectedThreads[i].AuthorID, thread.AuthorID)
		assert.Equal(t, expectedThreads[i].AuthorName, thread.AuthorName)
	}
}

func TestThreadRepository_GetThreadWithPosts(t *testing.T) {
	repo, mock, cleanup := setupThreadRepositoryTest(t)
	defer cleanup()

	// Настройка моков для транзакции
	mock.ExpectBegin()

	// Мок для получения треда
	expectedThread := &models.Thread{
		ID:        1,
		Title:     "Test Thread",
		AuthorID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectQuery("SELECT id, title, author_id, created_at, updated_at FROM threads WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author_id", "created_at", "updated_at"}).
			AddRow(expectedThread.ID, expectedThread.Title, expectedThread.AuthorID, expectedThread.CreatedAt, expectedThread.UpdatedAt))

	// Мок для получения постов
	expectedPosts := []models.Post{
		{
			ID:        1,
			ThreadID:  1,
			AuthorID:  1,
			Content:   "Post 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	postRows := sqlmock.NewRows([]string{"id", "thread_id", "user_id", "content", "created_at", "updated_at"})
	for _, post := range expectedPosts {
		postRows.AddRow(post.ID, post.ThreadID, post.AuthorID, post.Content, post.CreatedAt, post.UpdatedAt)
	}

	mock.ExpectQuery("SELECT id, thread_id, user_id, content, created_at, updated_at FROM posts WHERE thread_id = \\$1 ORDER BY created_at DESC").
		WithArgs(1).
		WillReturnRows(postRows)

	// Мок для получения комментариев
	expectedComments := map[int][]models.Comment{
		1: {
			{
				ID:        1,
				PostID:    1,
				AuthorID:  1,
				Content:   "Comment 1",
				CreatedAt: time.Now(),
			},
		},
	}

	commentRows := sqlmock.NewRows([]string{"id", "post_id", "user_id", "content", "created_at"})
	for _, comments := range expectedComments {
		for _, comment := range comments {
			commentRows.AddRow(comment.ID, comment.PostID, comment.AuthorID, comment.Content, comment.CreatedAt)
		}
	}

	mock.ExpectQuery("SELECT id, post_id, user_id, content, created_at FROM comments WHERE post_id = ANY\\(\\$1\\) ORDER BY created_at ASC").
		WithArgs(pq.Array([]int{1})).
		WillReturnRows(commentRows)

	mock.ExpectCommit()

	thread, posts, comments, err := repo.GetThreadWithPosts(1)
	require.NoError(t, err)
	assert.Equal(t, expectedThread.ID, thread.ID)
	assert.Equal(t, expectedThread.Title, thread.Title)
	assert.Equal(t, len(expectedPosts), len(posts))
	assert.Equal(t, len(expectedComments), len(comments))
} 