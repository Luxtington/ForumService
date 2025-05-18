package repository

import (
	"database/sql"

	"ForumService/internal/models"
	"fmt"
)

type ChatRepository interface {
	CreateMessage(authorID int, content string) (*models.ChatMessage, error)
	GetAllMessages() ([]*models.ChatMessage, error)
}

type chatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) ChatRepository {
	return &chatRepository{db: db}
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
		ORDER BY cm.created_at DESC`
	
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