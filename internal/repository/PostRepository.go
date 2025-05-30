package repository

import (
	"database/sql"
	"ForumService/internal/models"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"github.com/Luxtington/Shared/logger"
	"go.uber.org/zap"
)

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) GetByThreadID(threadID int) ([]*models.Post, error) {
	query := `
		SELECT p.id, p.thread_id, p.author_id, p.content, p.created_at, p.updated_at, u.username as author_name 
		FROM posts p
		LEFT JOIN users u ON p.author_id = u.id
		WHERE p.thread_id = $1 
		ORDER BY p.created_at ASC`
	rows, err := r.db.Query(query, threadID)
	if err != nil {
		return nil, err
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
			&post.AuthorName,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *postRepository) Create(post *models.Post) error {
	return r.SavePost(post)
}

func (r *postRepository) Update(post *models.Post) error {
	query := `UPDATE posts SET content = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.Exec(query, post.Content, post.ID)
	return err
}

func (r *postRepository) Delete(id int) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *postRepository) SavePost(post *models.Post) error {
	const query = `
		INSERT INTO posts (thread_id, author_id, content) 
		VALUES ($1, $2, $3)
		RETURNING id, thread_id, author_id, content, created_at
	`

	var newPost models.Post
	err := r.db.QueryRow(query, post.ThreadID, post.AuthorID, post.Content).Scan(
		&newPost.ID,
		&newPost.ThreadID,
		&newPost.AuthorID,
		&newPost.Content,
		&newPost.CreatedAt,
	)
	if err != nil {
		log := logger.GetLogger()
		log.Error("Ошибка при создании поста", zap.Error(err))
		return err
	}

	*post = newPost
	return nil
}

func (r *postRepository) GetPostByID(id int) (*models.Post, error) {
	query := `
		SELECT p.id, p.thread_id, p.author_id, p.content, p.created_at, p.updated_at, u.username as author_name
		FROM posts p
		LEFT JOIN users u ON p.author_id = u.id
		WHERE p.id = $1`

	post := &models.Post{}
	err := r.db.QueryRow(query, id).Scan(
		&post.ID,
		&post.ThreadID,
		&post.AuthorID,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.AuthorName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("пост не найден")
		}
		return nil, fmt.Errorf("ошибка при получении поста: %v", err)
	}

	return post, nil
}

func (r *postRepository) GetPostWithComments(postID int) (*models.Post, []models.Comment, error) {
	post, err := r.GetPostByID(postID)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка при получении поста: %w", err)
	}

	const query = `SELECT c.id, c.post_id, c.author_id, c.content, c.created_at, u.username as author_name
                   FROM comments c
                   LEFT JOIN users u ON c.author_id = u.id
                   WHERE c.post_id = $1
                   ORDER BY c.created_at ASC`

	rows, err := r.db.Query(query, postID)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка при получении комментариев: %w", err)
	}
	defer rows.Close()

	comments := make([]models.Comment, 0)
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.AuthorID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.AuthorName,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("ошибка при сканировании комментария: %w", err)
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("ошибка при итерации комментариев: %w", err)
	}

	return post, comments, nil
}

func (r *postRepository) GetPostsWithCommentsByThreadID(threadID int) ([]models.Post, map[int][]models.Comment, error) {
	// Получаем посты
	const postsQuery = `
        SELECT 
            p.id, p.thread_id, p.content, p.created_at,
            u.id as author_id, u.username as author_username
        FROM posts p
        JOIN users u ON p.user_id = u.id
        WHERE p.thread_id = $1
        ORDER BY p.created_at DESC
        LIMIT $2 OFFSET $3
    `

	limit := 20
	offset := 0

	postRows, err := r.db.Query(postsQuery, threadID, limit, offset)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка при получении постов: %w", err)
	}
	defer postRows.Close()

	var posts []models.Post
	var postIDs []int

	// Сканируем посты
	for postRows.Next() {
		var post models.Post
		err := postRows.Scan(
			&post.ID,
			&post.ThreadID,
			&post.Content,
			&post.CreatedAt,
			&post.AuthorID,
			&post.AuthorName,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("ошибка при сканировании поста: %w", err)
		}
		posts = append(posts, post)
		postIDs = append(postIDs, post.ID)
	}

	if err = postRows.Err(); err != nil {
		return nil, nil, fmt.Errorf("ошибка после сканирования постов: %w", err)
	}

	// Если постов нет, возвращаем пустые значения
	if len(postIDs) == 0 {
		return posts, make(map[int][]models.Comment), nil
	}

	// Получаем комментарии для всех постов
	const commentsQuery = `
        SELECT 
            c.id, c.post_id, c.author_id, c.content, c.created_at
        FROM comments c
        WHERE c.post_id = ANY($1)
        ORDER BY c.created_at ASC
    `

	commentRows, err := r.db.Query(commentsQuery, pq.Array(postIDs))
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка при получении комментариев: %w", err)
	}
	defer commentRows.Close()

	commentsByPostID := make(map[int][]models.Comment)
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
			return nil, nil, fmt.Errorf("ошибка при сканировании комментария: %w", err)
		}
		commentsByPostID[comment.PostID] = append(commentsByPostID[comment.PostID], comment)
	}

	if err = commentRows.Err(); err != nil {
		return nil, nil, fmt.Errorf("ошибка после сканирования комментариев: %w", err)
	}

	return posts, commentsByPostID, nil
}

func (r *postRepository) UpdatePost(post *models.Post, postID int) error {
	query := `UPDATE posts SET content = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.Exec(query, post.Content, postID)
	return err
}

func (r *postRepository) DeletePost(postID int) error {
	log := logger.GetLogger()
	
	tx, err := r.db.Begin()
	if err != nil {
		log.Error("error while tran begin in POST REPO 7", zap.Error(err))
		return err
	}
	defer tx.Rollback()

	const deleteCommentsQuery = `DELETE FROM comments WHERE post_id = $1`
	_, err = tx.Exec(deleteCommentsQuery, postID)
	if err != nil {
		log.Error("error while deleting comments in POST REPO 7.1", zap.Error(err))
		return err
	}

	const deletePostQuery = `DELETE FROM posts WHERE id = $1`
	result, err := tx.Exec(deletePostQuery, postID)
	if err != nil {
		log.Error("error while deleting post in POST REPO 7.2", zap.Error(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("error while checking post deleting in POST REPO 7.3", zap.Error(err))
		return err
	}

	if rowsAffected == 0 {
		log.Error("ERROR IN DELETE POST REPO 7.4")
		return errors.New("post not found in POST REPO")
	}

	if err = tx.Commit(); err != nil {
		log.Error("error while tran commit", zap.Error(err))
		return err
	}
	return nil
}

//func (r *PostRepository) GetPostsWithCommentsByThreadID(threadID int) ([]models.Post, error) {
//	const query = `
//		SELECT
//			p.id, p.thread_id, p.content, p.created_at,
//			u.id as author_id, u.username as author_username
//		FROM posts p
//		JOIN users u ON p.user_id = u.id
//		WHERE p.thread_id = $1
//		ORDER BY p.created_at
//		LIMIT $2 OFFSET $3
//	`
//
//	limit := 20 //  20 постов на странице
//	offset := 0
//
//	rows, err := r.db.Query(query, threadID, limit, offset)
//	if err != nil {
//		log.Println("ERROR IN POST REPO 3")
//		return nil, err
//	}
//	defer rows.Close()
//
//	var posts []models.Post
//	for rows.Next() {
//		var post models.Post
//		if err := rows.Scan(
//			&post.ID,
//			&post.ThreadID,
//			&post.Content,
//			&post.CreatedAt,
//			&post.AuthorID,
//		); err != nil {
//			return nil, err
//		}
//		posts = append(posts, post)
//	}
//
//	return posts, nil
//}
