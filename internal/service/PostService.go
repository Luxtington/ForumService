package service

import (
	"ForumService/internal/models"
	"ForumService/internal/repository"
)

type PostService interface {
	CreatePost(post *models.Post) error
	GetPostByID(id int) (*models.Post, error)
	GetPostWithComments(postID int) (*models.Post, []models.Comment, error)
	GetPostsWithCommentsByThreadID(threadID int) ([]models.Post, map[int][]models.Comment, error)
	UpdatePost(post *models.Post, postID int, userID int) error
	DeletePost(postID int, userID int) error
	GetAllPosts() ([]*models.Post, error)
	CreateComment(comment *models.Comment) error
	GetCommentByID(id int) (*models.Comment, error)
	DeleteComment(id int) error
	GetPost(id int) (*models.Post, error)
	GetPostsByThreadID(threadID int) ([]*models.Post, error)
	GetThreadByID(id int) (*models.Thread, error)
}

type postService struct {
	repo        repository.PostRepository
	commentRepo repository.CommentRepository
	threadRepo  repository.ThreadRepository
	userRepo    repository.UserRepository
}

func NewPostService(repo repository.PostRepository, commentRepo repository.CommentRepository, threadRepo repository.ThreadRepository, userRepo repository.UserRepository) PostService {
	return &postService{
		repo:        repo,
		commentRepo: commentRepo,
		threadRepo:  threadRepo,
		userRepo:    userRepo,
	}
}

func (s *postService) CreatePost(post *models.Post) error {
	return s.repo.SavePost(post)
}

func (s *postService) GetPostByID(id int) (*models.Post, error) {
	return s.repo.GetPostByID(id)
}

func (s *postService) GetPostWithComments(postID int) (*models.Post, []models.Comment, error) {
	post, comments, err := s.repo.GetPostWithComments(postID)
	if err != nil {
		return nil, nil, err
	}

	// Добавляем информацию о возможности редактирования для комментариев
	for i := range comments {
		comments[i].CanDelete = comments[i].AuthorID == post.AuthorID
	}

	return post, comments, nil
}

func (s *postService) GetPostsWithCommentsByThreadID(threadID int) ([]models.Post, map[int][]models.Comment, error) {
	posts, commentsMap, err := s.repo.GetPostsWithCommentsByThreadID(threadID)
	if err != nil {
		return nil, nil, err
	}

	// Добавляем информацию о возможности редактирования для комментариев
	for postID, comments := range commentsMap {
		for i := range comments {
			comments[i].CanDelete = comments[i].AuthorID == posts[postID].AuthorID
		}
	}

	return posts, commentsMap, nil
}

func (s *postService) UpdatePost(post *models.Post, postID int, userID int) error {
	existingPost, err := s.repo.GetPostByID(postID)
	if err != nil {
		return err
	}

	// Получаем роль пользователя
	userRole, err := s.userRepo.GetUserRole(userID)
	if err != nil {
		return err
	}

	if existingPost.AuthorID != userID && userRole != "admin" {
		return ErrNoPermission
	}

	return s.repo.UpdatePost(post, postID)
}

func (s *postService) DeletePost(postID int, userID int) error {
	post, err := s.repo.GetPostByID(postID)
	if err != nil {
		return err
	}

	// Получаем роль пользователя
	userRole, err := s.userRepo.GetUserRole(userID)
	if err != nil {
		return err
	}

	if post.AuthorID != userID && userRole != "admin" {
		return ErrNoPermission
	}

	return s.repo.DeletePost(postID)
}

func (s *postService) GetAllPosts() ([]*models.Post, error) {
	return s.repo.GetByThreadID(0) // 0 означает все посты
}

func (s *postService) CreateComment(comment *models.Comment) error {
	return s.commentRepo.SaveComment(comment)
}

func (s *postService) GetCommentByID(id int) (*models.Comment, error) {
	return s.commentRepo.GetCommentByID(id)
}

func (s *postService) DeleteComment(id int) error {
	return s.commentRepo.DeleteComment(id)
}

func (s *postService) GetPost(id int) (*models.Post, error) {
	return s.repo.GetPostByID(id)
}

func (s *postService) GetPostsByThreadID(threadID int) ([]*models.Post, error) {
	return s.repo.GetByThreadID(threadID)
}

func (s *postService) GetThreadByID(id int) (*models.Thread, error) {
	return s.threadRepo.GetByID(id)
}
