package tests

import (
	"database/sql"
	_"ForumService/internal/models"
	"ForumService/internal/repository"
	"testing"
	_"time"

	_"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/lib/pq"
)

func setupTest(t *testing.T) (*sql.DB, repository.UserRepository, repository.PostRepository, repository.CommentRepository, func()) {
	db, err := setupTestDB()
	require.NoError(t, err)

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)

	// Функция очистки после тестов
	cleanup := func() {
		db.Close()
	}

	return db, userRepo, postRepo, commentRepo, cleanup
}

/*
func TestUserCreation(t *testing.T) {
	_, userRepo, _, _, cleanup := setupTest(t)
	defer cleanup()

	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
	}

	err := userRepo.SaveUser(user)
	require.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestPostCreation(t *testing.T) {
	db, userRepo, postRepo, _, cleanup := setupTest(t)
	defer cleanup()

	// Сначала создаем пользователя
	user := &models.User{
		Username: "postauthor",
		Email:    "post@example.com",
	}
	err := userRepo.SaveUser(user)
	require.NoError(t, err)

	// Создаем тему
	thread := &models.Thread{
		Title:    "Test Thread",
		AuthorID: user.ID,
	}
	threadRepo := repository.NewThreadRepository(db)
	err = threadRepo.Create(thread)
	require.NoError(t, err)

	// Теперь создаем пост
	post := &models.Post{
		Content:    "Test post content",
		AuthorID:   user.ID,
		ThreadID:   thread.ID,
		AuthorName: user.Username,
	}

	err = postRepo.SavePost(post)
	require.NoError(t, err)
	assert.NotZero(t, post.ID)
}

func TestCommentCreation(t *testing.T) {
	db, userRepo, postRepo, commentRepo, cleanup := setupTest(t)
	defer cleanup()

	// Создаем пользователя
	user := &models.User{
		Username: "commentauthor",
		Email:    "comment@example.com",
	}
	err := userRepo.SaveUser(user)
	require.NoError(t, err)

	// Создаем тему
	thread := &models.Thread{
		Title:    "Test Thread for Comment",
		AuthorID: user.ID,
	}
	threadRepo := repository.NewThreadRepository(db)
	err = threadRepo.Create(thread)
	require.NoError(t, err)

	// Создаем пост
	post := &models.Post{
		Content:    "Test post for comment",
		AuthorID:   user.ID,
		ThreadID:   thread.ID,
		AuthorName: user.Username,
	}
	err = postRepo.SavePost(post)
	require.NoError(t, err)

	// Создаем комментарий
	comment := &models.Comment{
		Content:  "Test comment content",
		PostID:   post.ID,
		AuthorID: user.ID,
	}

	err = commentRepo.SaveComment(comment)
	require.NoError(t, err)
	assert.NotZero(t, comment.ID)
}

func TestPostRetrieval(t *testing.T) {
	db, userRepo, postRepo, _, cleanup := setupTest(t)
	defer cleanup()

	// Создаем пользователя
	user := &models.User{
		Username: "postretrieval",
		Email:    "retrieval@example.com",
	}
	err := userRepo.SaveUser(user)
	require.NoError(t, err)

	// Создаем тему
	thread := &models.Thread{
		Title:    "Test Thread for Retrieval",
		AuthorID: user.ID,
	}
	threadRepo := repository.NewThreadRepository(db)
	err = threadRepo.Create(thread)
	require.NoError(t, err)

	// Создаем пост
	post := &models.Post{
		Content:    "Test post for retrieval",
		AuthorID:   user.ID,
		ThreadID:   thread.ID,
		AuthorName: user.Username,
	}
	err = postRepo.SavePost(post)
	require.NoError(t, err)

	// Получаем пост
	retrievedPost, err := postRepo.GetPostByID(post.ID)
	require.NoError(t, err)
	assert.NotNil(t, retrievedPost)
	assert.Equal(t, post.Content, retrievedPost.Content)
	assert.Equal(t, post.AuthorID, retrievedPost.AuthorID)
	assert.Equal(t, post.ThreadID, retrievedPost.ThreadID)
}

func TestCommentRetrieval(t *testing.T) {
	db, userRepo, postRepo, commentRepo, cleanup := setupTest(t)
	defer cleanup()

	// Создаем пользователя
	user := &models.User{
		Username: "commentretrieval",
		Email:    "commentretrieval@example.com",
	}
	err := userRepo.SaveUser(user)
	require.NoError(t, err)

	// Создаем тему
	thread := &models.Thread{
		Title:    "Test Thread for Comment Retrieval",
		AuthorID: user.ID,
	}
	threadRepo := repository.NewThreadRepository(db)
	err = threadRepo.Create(thread)
	require.NoError(t, err)

	// Создаем пост
	post := &models.Post{
		Content:    "Test post for comment retrieval",
		AuthorID:   user.ID,
		ThreadID:   thread.ID,
		AuthorName: user.Username,
	}
	err = postRepo.SavePost(post)
	require.NoError(t, err)

	// Создаем комментарий
	comment := &models.Comment{
		Content:  "Test comment for retrieval",
		PostID:   post.ID,
		AuthorID: user.ID,
	}
	err = commentRepo.SaveComment(comment)
	require.NoError(t, err)

	// Получаем комментарий
	retrievedComment, err := commentRepo.GetCommentByID(comment.ID)
	require.NoError(t, err)
	assert.NotNil(t, retrievedComment)
	assert.Equal(t, comment.Content, retrievedComment.Content)
	assert.Equal(t, comment.PostID, retrievedComment.PostID)
	assert.Equal(t, comment.AuthorID, retrievedComment.AuthorID)
}

func TestUserRetrieval(t *testing.T) {
	_, userRepo, _, _, cleanup := setupTest(t)
	defer cleanup()

	user := &models.User{
		Username: "testuser2",
		Email:    "test2@example.com",
	}

	err := userRepo.SaveUser(user)
	require.NoError(t, err)

	retrievedUser, err := userRepo.GetUserByID(user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.Username, retrievedUser.Username)
	assert.Equal(t, user.Email, retrievedUser.Email)
}

func TestFullForumWorkflow(t *testing.T) {
	db, userRepo, postRepo, commentRepo, cleanup := setupTest(t)
	defer cleanup()

	threadRepo := repository.NewThreadRepository(db)

	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
	}
	err := userRepo.SaveUser(user)
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
		ThreadID:   thread.ID,
		AuthorID:   user.ID,
		Content:    "Test post content",
		AuthorName: user.Username,
	}
	err = postRepo.SavePost(post)
	require.NoError(t, err)
	assert.NotZero(t, post.ID)

	comment := &models.Comment{
		PostID:   post.ID,
		AuthorID: user.ID,
		Content:  "Test comment",
	}
	err = commentRepo.SaveComment(comment)
	require.NoError(t, err)
	assert.NotZero(t, comment.ID)

	retrievedThread, err := threadRepo.GetByID(thread.ID)
	require.NoError(t, err)
	assert.Equal(t, thread.Title, retrievedThread.Title)
	assert.Equal(t, user.ID, retrievedThread.AuthorID)

	retrievedPost, err := postRepo.GetPostByID(post.ID)
	require.NoError(t, err)
	assert.Equal(t, post.Content, retrievedPost.Content)
	assert.Equal(t, thread.ID, retrievedPost.ThreadID)
	assert.Equal(t, user.ID, retrievedPost.AuthorID)

	retrievedComment, err := commentRepo.GetCommentByID(comment.ID)
	require.NoError(t, err)
	assert.Equal(t, comment.Content, retrievedComment.Content)
	assert.Equal(t, post.ID, retrievedComment.PostID)
	assert.Equal(t, user.ID, retrievedComment.AuthorID)
}

func TestChatWorkflow(t *testing.T) {
	db, userRepo, _, _, cleanup := setupTest(t)
	defer cleanup()

	chatRepo := repository.NewChatRepository(db)

	user := &models.User{
		Username: "chatuser",
		Email:    "chat@example.com",
	}
	err := userRepo.SaveUser(user)
	require.NoError(t, err)
	assert.NotZero(t, user.ID)

	messageContent := "Test chat message"
	message, err := chatRepo.CreateMessage(user.ID, messageContent)
	require.NoError(t, err)
	assert.NotNil(t, message)

	messages, err := chatRepo.GetAllMessages()
	require.NoError(t, err)
	assert.NotEmpty(t, messages)

	found := false
	for _, msg := range messages {
		if msg.Content == messageContent && msg.AuthorID == user.ID {
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
		if msg.Content == messageContent && msg.AuthorID == user.ID {
			found = true
			break
		}
	}
	assert.False(t, found, "Old message was not deleted")
}
*/

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