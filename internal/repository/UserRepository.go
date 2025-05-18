package repository

import (
	"ForumService/internal/models"
	"database/sql"
	"fmt"
	"log"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) SaveUser(user *models.User) error {
	const query = `INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRow(query, user.Username, user.Email).Scan(&user.ID)
	if err != nil {
		log.Fatal("ERROR IN USER REPO 1")
		return err
	}
	return nil
}

func (r *userRepository) GetUserByID(id int) (*models.User, error) {
	const query = `
        SELECT id, username, email
        FROM users
        WHERE id = $1
    `
	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("пользователь с id %d не найден", id)
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении пользователя: %w", err)
	}
	return user, nil
}

func (r *userRepository) GetUserByUsername(username string) (*models.User, error) {
	const query = `
		SELECT id, username, email 
		FROM users 
		WHERE username = $1
	`

	var user models.User
	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email)

	if err != nil {
		log.Println("ERROR IN USER REPO 2")
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserPosts(userID int) ([]*models.Post, error) {
	query := `SELECT id, thread_id, author_id, content, created_at, updated_at FROM posts WHERE author_id = $1`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении постов пользователя: %w", err)
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(
			&post.ID,
			&post.ThreadID,
			&post.AuthorID,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании поста: %w", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *userRepository) GetUserCommentCount(userID int) (int, error) {
	query := `SELECT COUNT(*) FROM comments WHERE author_id = $1`
	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("ошибка при подсчете комментариев пользователя: %w", err)
	}
	return count, nil
}

func (r *userRepository) GetUserRole(userID int) (string, error) {
	var role string
	err := r.db.QueryRow("SELECT role FROM users WHERE id = $1", userID).Scan(&role)
	if err != nil {
		fmt.Printf("Debug - UserRepository.GetUserRole - Error: %v\n", err)
		return "", fmt.Errorf("couldn't get user role: %w", err)
	}
	fmt.Printf("Debug - UserRepository.GetUserRole - User ID: %d, Role: %s\n", userID, role)
	return role, nil
}

/*
CREATE TYPE user_role AS ENUM ('admin', 'user');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'user',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для быстрого поиска
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
*/
