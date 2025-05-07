package service

import (
	"ForumService/internal/models"
	"ForumService/internal/repository"
	"fmt"
)

type CommentService interface {
	CreateComment(postID, authorID int, content string) (*models.Comment, error)
	GetCommentByID(id int) (*models.Comment, error)
	GetCommentsByPostID(postID int) ([]models.Comment, error)
	DeleteComment(id int, userID int, isAdmin bool) error
}

type commentService struct {
	repo repository.CommentRepository
}

func NewCommentService(repo repository.CommentRepository) CommentService {
	return &commentService{repo: repo}
}

func (s *commentService) CreateComment(postID, authorID int, content string) (*models.Comment, error) {
	comment := &models.Comment{
		PostID:   postID,
		AuthorID: authorID,
		Content:  content,
	}

	if err := s.repo.SaveComment(comment); err != nil {
		return nil, fmt.Errorf("couldn't create comment: %w", err)
	}

	return comment, nil
}

func (s *commentService) GetCommentByID(id int) (*models.Comment, error) {
	comment, err := s.repo.GetCommentByID(id)
	if err != nil {
		return nil, fmt.Errorf("couldn't get comment: %w", err)
	}
	return comment, nil
}

func (s *commentService) GetCommentsByPostID(postID int) ([]models.Comment, error) {
	comments, err := s.repo.GetCommentsByPostID(postID)
	if err != nil {
		return nil, fmt.Errorf("couldn't get comments by post: %w", err)
	}
	return comments, nil
}

func (s *commentService) DeleteComment(id int, userID int, isAdmin bool) error {
	// Проверяем, существует ли комментарий
	comment, err := s.repo.GetCommentByID(id)
	if err != nil {
		return fmt.Errorf("comment not found: %w", err)
	}

	if !isAdmin && comment.AuthorID != userID {
		return fmt.Errorf("there are no rights for comment delete")
	}

	if err := s.repo.DeleteComment(id); err != nil {
		return fmt.Errorf("couldn't delete comment: %w", err)
	}

	return nil
}
