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
	DeleteComment(id int, userID int) error
}

type commentService struct {
	repo     repository.CommentRepository
	userRepo repository.UserRepository
}

func NewCommentService(repo repository.CommentRepository, userRepo repository.UserRepository) CommentService {
	return &commentService{
		repo:     repo,
		userRepo: userRepo,
	}
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

func (s *commentService) DeleteComment(commentID int, userID int) error {
	// Проверяем существование комментария
	comment, err := s.repo.GetCommentByID(commentID)
	if err != nil {
		return err
	}

	// Получаем роль пользователя
	userRole, err := s.userRepo.GetUserRole(userID)
	if err != nil {
		return err
	}

	// Отладочная информация
	fmt.Printf("Debug - CommentService.DeleteComment - User ID: %d, Role: %s\n", userID, userRole)

	// Проверяем права доступа
	if comment.AuthorID != userID && userRole != "admin" {
		return ErrNoPermission
	}

	return s.repo.DeleteComment(commentID)
}
