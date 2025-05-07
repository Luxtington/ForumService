package models

import "time"

type ChatMessage struct {
    ID        int       `json:"id"`
    AuthorID  int       `json:"author_id"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
} 