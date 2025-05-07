package service

import (
	"ForumService/internal/models"
	"ForumService/internal/repository"
	"errors"
	"fmt"
)

type ThreadService interface {
	GetThreadWithPosts(threadID int) (*models.Thread, []*models.Post, error)
	CreateThread(thread *models.Thread) error
	UpdateThread(thread *models.Thread) error
	DeleteThread(threadID int) error
	GetAllThreads() ([]*models.Thread, error)
}

type threadService struct {
	threadRepo repository.ThreadRepository
	postRepo   repository.PostRepository
}

func NewThreadService(threadRepo repository.ThreadRepository, postRepo repository.PostRepository) ThreadService {
	return &threadService{
		threadRepo: threadRepo,
		postRepo:   postRepo,
	}
}

func (s *threadService) GetThreadWithPosts(threadID int) (*models.Thread, []*models.Post, error) {
	fmt.Printf("Получение треда с ID: %d\n", threadID)
	thread, err := s.threadRepo.GetByID(threadID)
	if err != nil {
		fmt.Printf("Ошибка при получении треда из репозитория: %v\n", err)
		return nil, nil, err
	}
	if thread == nil {
		fmt.Printf("Тред не найден\n")
		return nil, nil, errors.New("thread not found")
	}

	fmt.Printf("Тред найден: %+v\n", thread)
	posts, err := s.postRepo.GetByThreadID(threadID)
	if err != nil {
		fmt.Printf("Ошибка при получении постов: %v\n", err)
		return nil, nil, err
	}

	fmt.Printf("Получено постов: %d\n", len(posts))
	return thread, posts, nil
}

func (s *threadService) CreateThread(thread *models.Thread) error {
	return s.threadRepo.Create(thread)
}

func (s *threadService) UpdateThread(thread *models.Thread) error {
	return s.threadRepo.Update(thread)
}

func (s *threadService) DeleteThread(threadID int) error {
	return s.threadRepo.Delete(threadID)
}

func (s *threadService) GetAllThreads() ([]*models.Thread, error) {
	threads, err := s.threadRepo.GetAllThreads()
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении тредов: %v", err)
	}
	return threads, nil
}
