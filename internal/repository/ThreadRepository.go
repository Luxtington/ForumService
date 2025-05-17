package repository

import (
	"database/sql"
	"ForumService/internal/models"
	"fmt"
	"github.com/lib/pq"
)

type threadRepository struct {
	db *sql.DB
}

func NewThreadRepository(db *sql.DB) ThreadRepository {
	return &threadRepository{db: db}
}

func (r *threadRepository) GetByID(id int) (*models.Thread, error) {
	query := `SELECT id, title, author_id, created_at, updated_at FROM threads WHERE id = $1`
	thread := &models.Thread{}
	err := r.db.QueryRow(query, id).Scan(
		&thread.ID,
		&thread.Title,
		&thread.AuthorID,
		&thread.CreatedAt,
		&thread.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return thread, nil
}

func (r *threadRepository) Create(thread *models.Thread) error {
	query := `INSERT INTO threads (title, author_id, created_at, updated_at) 
			  VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) 
			  RETURNING id, title, author_id, created_at, updated_at`
	
	fmt.Printf("Создание треда: title=%s, author_id=%d\n", thread.Title, thread.AuthorID)
	
	err := r.db.QueryRow(query, thread.Title, thread.AuthorID).Scan(
		&thread.ID,
		&thread.Title,
		&thread.AuthorID,
		&thread.CreatedAt,
		&thread.UpdatedAt,
	)
	
	if err != nil {
		fmt.Printf("Ошибка при создании треда: %v\n", err)
		return err
	}
	
	fmt.Printf("Тред создан: id=%d, title=%s, author_id=%d, created_at=%v\n", 
		thread.ID, thread.Title, thread.AuthorID, thread.CreatedAt)
	
	return nil
}

func (r *threadRepository) Update(thread *models.Thread) error {
	query := `UPDATE threads SET title = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.Exec(query, thread.Title, thread.ID)
	return err
}

func (r *threadRepository) Delete(id int) error {
	query := `DELETE FROM threads WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// CreateThread создает новый тред
func (r *threadRepository) CreateThread(thread *models.Thread) error {
	const query = `
        INSERT INTO threads (title, author_id, created_at, updated_at)
        VALUES ($1, $2, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `

	err := r.db.QueryRow(
		query,
		thread.Title,
		thread.AuthorID,
	).Scan(
		&thread.ID,
		&thread.CreatedAt,
		&thread.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create thread: %w", err)
	}

	return nil
}

// GetThreadWithPosts получает тред по ID вместе со всеми постами и их комментариями
func (r *threadRepository) GetThreadWithPosts(threadID int) (*models.Thread, []models.Post, map[int][]models.Comment, error) {
	// Начинаем транзакцию для обеспечения консистентности данных
	tx, err := r.db.Begin()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Получаем информацию о треде
	const threadQuery = `
        SELECT id, title, author_id, created_at, updated_at
        FROM threads
        WHERE id = $1
    `

	thread := &models.Thread{}
	err = tx.QueryRow(threadQuery, threadID).Scan(
		&thread.ID,
		&thread.Title,
		&thread.AuthorID,
		&thread.CreatedAt,
		&thread.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil, fmt.Errorf("thread not found")
		}
		return nil, nil, nil, fmt.Errorf("failed to get thread: %w", err)
	}

	// 2. Получаем все посты треда
	const postsQuery = `
        SELECT id, thread_id, user_id, content, created_at, updated_at
        FROM posts
        WHERE thread_id = $1
        ORDER BY created_at DESC
    `

	postRows, err := tx.Query(postsQuery, threadID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get posts: %w", err)
	}
	defer postRows.Close()

	var posts []models.Post
	var postIDs []int

	for postRows.Next() {
		var post models.Post
		err := postRows.Scan(
			&post.ID,
			&post.ThreadID,
			&post.AuthorID,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
		postIDs = append(postIDs, post.ID)
	}

	if err = postRows.Err(); err != nil {
		return nil, nil, nil, fmt.Errorf("error iterating posts: %w", err)
	}

	// 3. Если есть посты, получаем все комментарии для них
	commentsByPostID := make(map[int][]models.Comment)

	if len(postIDs) > 0 {
		const commentsQuery = `
            SELECT id, post_id, user_id, content, created_at
            FROM comments
            WHERE post_id = ANY($1)
            ORDER BY created_at ASC
        `

		commentRows, err := tx.Query(commentsQuery, pq.Array(postIDs))
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to get comments: %w", err)
		}
		defer commentRows.Close()

		for commentRows.Next() {
			var comment models.Comment
			err := commentRows.Scan(
				&comment.ID,
				&comment.PostID,
				&comment.AuthorID,
				&comment.Content,
				&comment.CreatedAt,
			)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("failed to scan comment: %w", err)
			}
			commentsByPostID[comment.PostID] = append(commentsByPostID[comment.PostID], comment)
		}

		if err = commentRows.Err(); err != nil {
			return nil, nil, nil, fmt.Errorf("error iterating comments: %w", err)
		}
	}

	// Фиксируем транзакцию
	if err = tx.Commit(); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return thread, posts, commentsByPostID, nil
}

// DeleteThread удаляет тред (каскадное удаление автоматически удалит посты и комментарии)
func (r *threadRepository) DeleteThread(threadID int) error {
	const query = `
        DELETE FROM threads
        WHERE id = $1
    `

	result, err := r.db.Exec(query, threadID)
	if err != nil {
		return fmt.Errorf("failed to delete thread: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("thread not found")
	}

	return nil
}

//CREATE TABLE threads (
//id SERIAL PRIMARY KEY,
//title VARCHAR(255) NOT NULL,
//author_id INTEGER NOT NULL REFERENCES users(id),
//created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
//updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
//);
//
//-- Таблица постов с каскадным удалением
//CREATE TABLE posts (
//id SERIAL PRIMARY KEY,
//thread_id INTEGER NOT NULL REFERENCES threads(id) ON DELETE CASCADE,
//user_id INTEGER NOT NULL REFERENCES users(id),
//content TEXT NOT NULL,
//created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
//updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
//);
//
//-- Таблица комментариев с каскадным удалением
//CREATE TABLE comments (
//id SERIAL PRIMARY KEY,
//post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
//user_id INTEGER NOT NULL REFERENCES users(id),
//content TEXT NOT NULL,
//created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
//);
//
//-- Индексы для оптимизации
//CREATE INDEX idx_threads_author_id ON threads(author_id);
//CREATE INDEX idx_threads_created_at ON threads(created_at);
//CREATE INDEX idx_posts_thread_id ON posts(thread_id);
//CREATE INDEX idx_posts_created_at ON posts(created_at);
//CREATE INDEX idx_comments_post_id ON comments(post_id);
//CREATE INDEX idx_comments_created_at ON comments(created_at);

func (r *threadRepository) GetAllThreads() ([]*models.Thread, error) {
	query := `
		SELECT t.id, t.title, t.author_id, t.created_at, t.updated_at, u.username as author_name
		FROM threads t
		LEFT JOIN users u ON t.author_id = u.id
		ORDER BY t.created_at DESC
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении тредов: %v", err)
	}
	defer rows.Close()

	var threads []*models.Thread
	for rows.Next() {
		thread := &models.Thread{}
		err := rows.Scan(
			&thread.ID,
			&thread.Title,
			&thread.AuthorID,
			&thread.CreatedAt,
			&thread.UpdatedAt,
			&thread.AuthorName,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании треда: %v", err)
		}
		threads = append(threads, thread)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по тредам: %v", err)
	}

	return threads, nil
}
