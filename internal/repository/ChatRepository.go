package repository

import (
	"database/sql"
	"github.com/Luxtington/Shared/logger"
	"go.uber.org/zap"
	"time"

	"ForumService/internal/models"
	"fmt"
)

type ChatRepository interface {
	CreateMessage(authorID int, content string) (*models.ChatMessage, error)
	GetAllMessages() ([]*models.ChatMessage, error)
	DeleteOldMessages() error
	CleanOldMessages() error
	Cleanup() error
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
	log := logger.GetLogger()
	query := `DELETE FROM chat_messages WHERE created_at < NOW() - INTERVAL '1 minute' RETURNING id`
	rows, err := r.db.Query(query)
	if err != nil {
		log.Error("Ошибка при удалении старых сообщений", zap.Error(err))
		return err
	}
	defer rows.Close()

	var deletedIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Error("Ошибка при сканировании ID удаленного сообщения", zap.Error(err))
			continue
		}
		deletedIDs = append(deletedIDs, id)
	}

	if len(deletedIDs) > 0 {
		log.Info("Удалены сообщения", zap.Ints("message_ids", deletedIDs))
	}

	return nil
}

func (r *chatRepository) startMessageCleanup() {
	log := logger.GetLogger()
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := r.DeleteOldMessages(); err != nil {
			log.Error("Ошибка при очистке старых сообщений", zap.Error(err))
		}
	}
}

func (r *chatRepository) CleanOldMessages() error {
	log := logger.GetLogger()
	
	// Удаляем старые сообщения
	_, err := r.db.Exec("DELETE FROM chat_messages WHERE created_at < NOW() - INTERVAL '24 hours'")
	if err != nil {
		log.Error("Ошибка при удалении старых сообщений", zap.Error(err))
		return err
	}

	// Получаем ID удаленных сообщений
	var deletedIDs []int
	rows, err := r.db.Query("SELECT id FROM chat_messages WHERE created_at < NOW() - INTERVAL '24 hours'")
	if err != nil {
		log.Error("Ошибка при сканировании ID удаленного сообщения", zap.Error(err))
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Error("Ошибка при сканировании ID удаленного сообщения", zap.Error(err))
			return err
		}
		deletedIDs = append(deletedIDs, id)
	}

	log.Info("Удалены сообщения", zap.Ints("message_ids", deletedIDs))
	return nil
}

func (r *chatRepository) Cleanup() error {
	log := logger.GetLogger()
	
	_, err := r.db.Exec("DELETE FROM chat_messages WHERE created_at < NOW() - INTERVAL '24 hours'")
	if err != nil {
		log.Error("Ошибка при очистке старых сообщений", zap.Error(err))
		return err
	}
	return nil
} 