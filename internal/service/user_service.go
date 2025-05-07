package service

import (
	"ForumService/internal/models"
	"ForumService/internal/repository"
)

type UserService interface {
	GetUserByID(id int) (*models.User, error)
	GetUserPosts(userID int) ([]*models.Post, error)
	GetUserCommentCount(userID int) (int, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUserByID(id int) (*models.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *userService) GetUserPosts(userID int) ([]*models.Post, error) {
	return s.repo.GetUserPosts(userID)
}

func (s *userService) GetUserCommentCount(userID int) (int, error) {
	return s.repo.GetUserCommentCount(userID)
} 