package repository

import (
	"ForumService/internal/models"
	"database/sql"
	"fmt"
)

type CommentRepositoryImpl struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &CommentRepositoryImpl{db: db}
}

func (r *CommentRepositoryImpl) SaveComment(comment *models.Comment) error {
	const query = `INSERT INTO comments (post_id, author_id, content, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id`
	return r.db.QueryRow(query, comment.PostID, comment.AuthorID, comment.Content).Scan(&comment.ID)
}

func (r *CommentRepositoryImpl) GetCommentByID(id int) (*models.Comment, error) {
	const query = `SELECT id, post_id, author_id, content, created_at FROM comments WHERE id = $1`
	comment := &models.Comment{}
	err := r.db.QueryRow(query, id).Scan(&comment.ID, &comment.PostID, &comment.AuthorID, &comment.Content, &comment.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("комментарий не найден")
		}
		return nil, fmt.Errorf("ошибка при получении комментария: %v", err)
	}
	return comment, nil
}

//func (r *CommentRepositoryImpl) List() ([]models.Comment, error) {
//	const query = `SELECT * FROM comments ORDER BY created_at DESC RETURNING id, post_id, author_id, content, created_at`
//	rows, err := r.db.Query(query)
//	defer rows.Close()
//	if err != nil {
//		return nil, err
//	}
//	allComments := []models.Comment{}
//	for rows.Next() {
//		comment := models.Comment{}
//		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.AuthorID, &comment.Content, &comment.CreatedAt); err != nil {
//			return nil, err
//		}
//		allComments = append(allComments, comment)
//	}
//	return allComments, nil
//}
//func (r *CommentRepositoryImpl) UpdateUser(comment *models.Comment, id int) error {
//	query := `UPDATE comments SET content = $1, updated_at = NOW() WHERE id = $2 RETURNING updated_at`
//	return r.db.QueryRow(query, comment.Content, id).Scan(&comment.UpdatedAt)
//}

func (r *CommentRepositoryImpl) DeleteComment(id int) error {
	const query = `DELETE FROM comments WHERE id = $1 RETURNING id`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *CommentRepositoryImpl) GetCommentsByPostID(postID int) ([]models.Comment, error) {
	const query = `
        SELECT id, post_id, author_id, content, created_at
        FROM comments
        WHERE post_id = $1
        ORDER BY created_at ASC`

	rows, err := r.db.Query(query, postID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении комментариев: %v", err)
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.AuthorID,
			&comment.Content,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании комментария: %v", err)
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка после сканирования комментариев: %v", err)
	}

	return comments, nil
}
