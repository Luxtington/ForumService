package models

import "time"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Thread struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	AuthorID   int       `json:"author_id"`
	AuthorName string    `json:"author_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Post struct {
	ID        int       `json:"id"`
	ThreadID  int       `json:"thread_id"`
	AuthorID  int       `json:"author_id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	AuthorName string   `json:"author_name"`
	Comments  []Comment `json:"comments,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CanEdit   bool      `json:"can_edit"`
}

type Comment struct {
	ID         int       `json:"id"`
	PostID     int       `json:"post_id"`
	AuthorID   int       `json:"author_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	CanDelete  bool      `json:"can_delete"`
	AuthorName string    `json:"author_name"`
}

//	type CreateThreadRequest struct {
//		Content  string `json:"content"`
//		AuthorID int    `json:"author_id"`
//	}
//type CreatePostRequest struct {
//	ThreadID int    `json:"thread_id"`
//	Content  string `json:"content"`
//	AuthorID int    `json:"author_id"`
//}

//
//type CreateCommentRequest struct {
//	PostID   int    `json:"post_id"`
//	Content  string `json:"content"`
//	AuthorID int    `json:"author_id"`
//}
//
//type UpdatePostRequest struct {
//	PostID    int    `json:"post_id"` // ID поста, который нужно обновить
//	Content   string `json:"content"` // Новое содержимое поста
//	UpdaterID int    `json:"updater_id"`
//}
//
//type DeletePostRequest struct {
//	PostID    int `json:"post_id"`
//	DeleterID int `json:"author_id"`
//}
