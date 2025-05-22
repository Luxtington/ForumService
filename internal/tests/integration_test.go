package tests

import (
	"database/sql"
	"testing"
	"time"

	"ForumService/internal/models"
	"ForumService/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFullForumWorkflow(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	threadRepo := repository.NewThreadRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)

	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "user",
	}
	err = userRepo.Create(user)
	require.NoError(t, err)
	assert.NotZero(t, user.ID)

	thread := &models.Thread{
		Title:    "Test Thread",
		AuthorID: user.ID,
	}
	err = threadRepo.Create(thread)
	require.NoError(t, err)
	assert.NotZero(t, thread.ID)

	post := &models.Post{
		ThreadID: thread.ID,
		AuthorID: user.ID,
		Content:  "Test post content",
	}
	err = postRepo.Create(post)
	require.NoError(t, err)
	assert.NotZero(t, post.ID)

	comment := &models.Comment{
		PostID:   post.ID,
		AuthorID: user.ID,
		Content:  "Test comment",
	}
	err = commentRepo.Create(comment)
	require.NoError(t, err)
	assert.NotZero(t, comment.ID)

	retrievedThread, err := threadRepo.GetByID(thread.ID)
	require.NoError(t, err)
	assert.Equal(t, thread.Title, retrievedThread.Title)
	assert.Equal(t, user.ID, retrievedThread.AuthorID)

	retrievedPost, err := postRepo.GetByID(post.ID)
	require.NoError(t, err)
	assert.Equal(t, post.Content, retrievedPost.Content)
	assert.Equal(t, thread.ID, retrievedPost.ThreadID)
	assert.Equal(t, user.ID, retrievedPost.AuthorID)

	retrievedComment, err := commentRepo.GetByID(comment.ID)
	require.NoError(t, err)
	assert.Equal(t, comment.Content, retrievedComment.Content)
	assert.Equal(t, post.ID, retrievedComment.PostID)
	assert.Equal(t, user.ID, retrievedComment.AuthorID)
}

func TestChatWorkflow(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	chatRepo := repository.NewChatRepository(db)

	user := &models.User{
		Username: "chatuser",
		Email:    "chat@example.com",
		Password: "password123",
		Role:     "user",
	}
	err = userRepo.Create(user)
	require.NoError(t, err)
	assert.NotZero(t, user.ID)

	message := &models.ChatMessage{
		AuthorID: user.ID,
		Content:  "Test chat message",
	}
	err = chatRepo.CreateMessage(user.ID, message.Content)
	require.NoError(t, err)

	messages, err := chatRepo.GetAllMessages()
	require.NoError(t, err)
	assert.NotEmpty(t, messages)

	found := false
	for _, msg := range messages {
		if msg.Content == message.Content && msg.AuthorID == user.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "Created message not found in chat")

	time.Sleep(65 * time.Second)

	// Проверяем, что сообщение было удалено
	messages, err = chatRepo.GetAllMessages()
	require.NoError(t, err)
	found = false
	for _, msg := range messages {
		if msg.Content == message.Content && msg.AuthorID == user.ID {
			found = true
			break
		}
	}
	assert.False(t, found, "Old message was not deleted")
}

func setupTestDB() (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=forum_test sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Проверяем подключение
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
} 