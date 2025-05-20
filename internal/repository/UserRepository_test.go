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

func setupUserRepositoryTest(t *testing.T) (*userRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewUserRepository(db).(*userRepository)

	cleanup := func() {
		db.Close()
	}

	return repo, mock, cleanup
}

func TestUserRepository_SaveUser(t *testing.T) {
	repo, mock, cleanup := setupUserRepositoryTest(t)
	defer cleanup()

	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
	}

	mock.ExpectQuery("INSERT INTO users").
		WithArgs(user.Username, user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err := repo.SaveUser(user)
	require.NoError(t, err)
	assert.Equal(t, 1, user.ID)
}

func TestUserRepository_GetUserByID(t *testing.T) {
	repo, mock, cleanup := setupUserRepositoryTest(t)
	defer cleanup()

	expectedUser := &models.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}

	mock.ExpectQuery("SELECT id, username, email FROM users WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email"}).
			AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Email))

	user, err := repo.GetUserByID(1)
	require.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Username, user.Username)
	assert.Equal(t, expectedUser.Email, user.Email)
}

func TestUserRepository_GetUserByUsername(t *testing.T) {
	repo, mock, cleanup := setupUserRepositoryTest(t)
	defer cleanup()

	expectedUser := &models.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}

	mock.ExpectQuery("SELECT id, username, email FROM users WHERE username = \\$1").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email"}).
			AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Email))

	user, err := repo.GetUserByUsername("testuser")
	require.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Username, user.Username)
	assert.Equal(t, expectedUser.Email, user.Email)
}

func TestUserRepository_GetUserPosts(t *testing.T) {
	repo, mock, cleanup := setupUserRepositoryTest(t)
	defer cleanup()

	expectedPosts := []*models.Post{
		{
			ID:        1,
			ThreadID:  1,
			AuthorID:  1,
			Content:   "Test Post 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			ThreadID:  1,
			AuthorID:  1,
			Content:   "Test Post 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "thread_id", "author_id", "content", "created_at", "updated_at"})
	for _, post := range expectedPosts {
		rows.AddRow(post.ID, post.ThreadID, post.AuthorID, post.Content, post.CreatedAt, post.UpdatedAt)
	}

	mock.ExpectQuery("SELECT id, thread_id, author_id, content, created_at, updated_at FROM posts WHERE author_id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)

	posts, err := repo.GetUserPosts(1)
	require.NoError(t, err)
	assert.Equal(t, len(expectedPosts), len(posts))
	for i, post := range posts {
		assert.Equal(t, expectedPosts[i].ID, post.ID)
		assert.Equal(t, expectedPosts[i].Content, post.Content)
	}
}

func TestUserRepository_GetUserCommentCount(t *testing.T) {
	repo, mock, cleanup := setupUserRepositoryTest(t)
	defer cleanup()

	expectedCount := 5

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM comments WHERE author_id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

	count, err := repo.GetUserCommentCount(1)
	require.NoError(t, err)
	assert.Equal(t, expectedCount, count)
}

func TestUserRepository_GetUserRole(t *testing.T) {
	repo, mock, cleanup := setupUserRepositoryTest(t)
	defer cleanup()

	expectedRole := "admin"

	mock.ExpectQuery("SELECT role FROM users WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow(expectedRole))

	role, err := repo.GetUserRole(1)
	require.NoError(t, err)
	assert.Equal(t, expectedRole, role)
} 