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

func setupChatRepositoryTest(t *testing.T) (*chatRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewChatRepository(db).(*chatRepository)

	cleanup := func() {
		db.Close()
	}

	return repo, mock, cleanup
}

func TestChatRepository_CreateMessage(t *testing.T) {
	repo, mock, cleanup := setupChatRepositoryTest(t)
	defer cleanup()

	authorID := 1
	content := "Test Message"

	mock.ExpectQuery("INSERT INTO chat_messages").
		WithArgs(authorID, content).
		WillReturnRows(sqlmock.NewRows([]string{"id", "author_id", "content", "created_at"}).
			AddRow(1, authorID, content, time.Now()))

	message, err := repo.CreateMessage(authorID, content)
	require.NoError(t, err)
	assert.Equal(t, 1, message.ID)
	assert.Equal(t, authorID, message.AuthorID)
	assert.Equal(t, content, message.Content)
}

func TestChatRepository_GetAllMessages(t *testing.T) {
	repo, mock, cleanup := setupChatRepositoryTest(t)
	defer cleanup()

	expectedMessages := []*models.ChatMessage{
		{
			ID:        1,
			AuthorID:  1,
			Content:   "Test Message 1",
			CreatedAt: time.Now(),
			AuthorName: "Test User 1",
		},
		{
			ID:        2,
			AuthorID:  2,
			Content:   "Test Message 2",
			CreatedAt: time.Now(),
			AuthorName: "Test User 2",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "author_id", "content", "created_at", "author_name"})
	for _, message := range expectedMessages {
		rows.AddRow(message.ID, message.AuthorID, message.Content, message.CreatedAt, message.AuthorName)
	}

	mock.ExpectQuery("SELECT cm.id, cm.author_id, cm.content, cm.created_at, u.username as author_name FROM chat_messages cm LEFT JOIN users u ON cm.author_id = u.id ORDER BY cm.created_at ASC").
		WillReturnRows(rows)

	messages, err := repo.GetAllMessages()
	require.NoError(t, err)
	assert.Equal(t, len(expectedMessages), len(messages))
	for i, message := range messages {
		assert.Equal(t, expectedMessages[i].ID, message.ID)
		assert.Equal(t, expectedMessages[i].Content, message.Content)
		assert.Equal(t, expectedMessages[i].AuthorName, message.AuthorName)
	}
}

func TestChatRepository_DeleteOldMessages(t *testing.T) {
	repo, mock, cleanup := setupChatRepositoryTest(t)
	defer cleanup()

	mock.ExpectQuery("DELETE FROM chat_messages WHERE created_at < NOW\\(\\) - INTERVAL '1 minute' RETURNING id").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2))

	err := repo.DeleteOldMessages()
	require.NoError(t, err)
}

func TestChatRepository_CleanOldMessages(t *testing.T) {
	repo, mock, cleanup := setupChatRepositoryTest(t)
	defer cleanup()

	mock.ExpectExec("DELETE FROM chat_messages WHERE created_at < NOW\\(\\) - INTERVAL '24 hours'").
		WillReturnResult(sqlmock.NewResult(0, 2))

	mock.ExpectQuery("SELECT id FROM chat_messages WHERE created_at < NOW\\(\\) - INTERVAL '24 hours'").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2))

	err := repo.CleanOldMessages()
	require.NoError(t, err)
}

func TestChatRepository_Cleanup(t *testing.T) {
	repo, mock, cleanup := setupChatRepositoryTest(t)
	defer cleanup()

	mock.ExpectExec("DELETE FROM chat_messages WHERE created_at < NOW\\(\\) - INTERVAL '24 hours'").
		WillReturnResult(sqlmock.NewResult(0, 2))

	err := repo.Cleanup()
	require.NoError(t, err)
} 