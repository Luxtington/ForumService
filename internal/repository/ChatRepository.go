package repository

import (
	"database/sql"
	"log"
	"time"

	"ForumService/internal/models"
	"fmt"
)

type ChatRepository interface {
	CreateMessage(authorID int, content string) (*models.ChatMessage, error)
	GetAllMessages() ([]*models.ChatMessage, error)
	DeleteOldMessages() error
}

type chatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) ChatRepository {
	repo := &chatRepository{db: db}
	go repo.startMessageCleanup()
	return repo
}

func (r *chatRepository) CreateMessage(authorID int, content string) (*models.ChatMessage, error) {
	message := &models.ChatMessage{
		AuthorID: authorID,
		Content:  content,
	}

	query := `INSERT INTO chat_messages (author_id, content, created_at) 
			  VALUES ($1, $2, CURRENT_TIMESTAMP) 
			  RETURNING id, author_id, content, created_at`
	
	fmt.Printf("Создание сообщения чата: author_id=%d, content=%s\n", authorID, content)
	
	err := r.db.QueryRow(query, message.AuthorID, message.Content).Scan(
		&message.ID,
		&message.AuthorID,
		&message.Content,
		&message.CreatedAt,
	)
	if err != nil {
		fmt.Printf("Ошибка при создании сообщения чата: %v\n", err)
		return nil, err
	}

	fmt.Printf("Сообщение чата создано: id=%d, author_id=%d, content=%s, created_at=%v\n",
		message.ID, message.AuthorID, message.Content, message.CreatedAt)

	return message, nil
}

func (r *chatRepository) GetAllMessages() ([]*models.ChatMessage, error) {
	query := `
		SELECT cm.id, cm.author_id, cm.content, cm.created_at, u.username as author_name 
		FROM chat_messages cm
		LEFT JOIN users u ON cm.author_id = u.id 
		ORDER BY cm.created_at ASC`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.ChatMessage
	for rows.Next() {
		message := &models.ChatMessage{}
		err := rows.Scan(
			&message.ID,
			&message.AuthorID,
			&message.Content,
			&message.CreatedAt,
			&message.AuthorName,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *chatRepository) DeleteOldMessages() error {
	query := `DELETE FROM chat_messages WHERE created_at < NOW() - INTERVAL '1 minute' RETURNING id`
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("Ошибка при удалении старых сообщений: %v", err)
		return err
	}
	defer rows.Close()

	var deletedIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Printf("Ошибка при сканировании ID удаленного сообщения: %v", err)
			continue
		}
		deletedIDs = append(deletedIDs, id)
	}

	if len(deletedIDs) > 0 {
		log.Printf("Удалены сообщения с ID: %v", deletedIDs)
	}

	return nil
}

func (r *chatRepository) startMessageCleanup() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := r.DeleteOldMessages(); err != nil {
			log.Printf("Ошибка при очистке старых сообщений: %v", err)
		}
	}
} 