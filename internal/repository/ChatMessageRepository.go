package repository

import (
    "database/sql"
    "ForumService/internal/models"
)

type ChatMessageRepository struct {
    db *sql.DB
}

func NewChatMessageRepository(db *sql.DB) *ChatMessageRepository {
    return &ChatMessageRepository{db: db}
}

func (r *ChatMessageRepository) CreateMessage(authorID int, content string) (*models.ChatMessage, error) {
    var message models.ChatMessage
    err := r.db.QueryRow(
        "INSERT INTO chat_messages (author_id, content) VALUES ($1, $2) RETURNING id, author_id, content, created_at",
        authorID, content,
    ).Scan(&message.ID, &message.AuthorID, &message.Content, &message.CreatedAt)
    
    if err != nil {
        return nil, err
    }
    
    return &message, nil
}

func (r *ChatMessageRepository) GetAllMessages() ([]*models.ChatMessage, error) {
    rows, err := r.db.Query(
        "SELECT id, author_id, content, created_at FROM chat_messages ORDER BY created_at ASC",
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var messages []*models.ChatMessage
    for rows.Next() {
        var message models.ChatMessage
        if err := rows.Scan(&message.ID, &message.AuthorID, &message.Content, &message.CreatedAt); err != nil {
            return nil, err
        }
        messages = append(messages, &message)
    }

    return messages, nil
}

func (r *ChatMessageRepository) DeleteMessage(id int) error {
    result, err := r.db.Exec("DELETE FROM chat_messages WHERE id = $1", id)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return sql.ErrNoRows
    }

    return nil
} 