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
	UpdatePost(post *models.Post, postID int) error
	DeletePost(postID int) error
	GetAllPosts() ([]*models.Post, error)
	CreateComment(comment *models.Comment) error
	GetCommentByID(id int) (*models.Comment, error)
	DeleteComment(id int) error
	GetPost(id int) (*models.Post, error)
	GetPostsByThreadID(threadID int) ([]*models.Post, error)
}

type postService struct {
	repo        repository.PostRepository
	commentRepo repository.CommentRepository
}

func NewPostService(repo repository.PostRepository, commentRepo repository.CommentRepository) PostService {
	return &postService{
		repo:        repo,
		commentRepo: commentRepo,
	}
}

func (s *postService) CreatePost(post *models.Post) error {
	return s.repo.SavePost(post)
}

func (s *postService) GetPostByID(id int) (*models.Post, error) {
	return s.repo.GetPostByID(id)
}

func (s *postService) GetPostWithComments(postID int) (*models.Post, []models.Comment, error) {
	return s.repo.GetPostWithComments(postID)
}

func (s *postService) GetPostsWithCommentsByThreadID(threadID int) ([]models.Post, map[int][]models.Comment, error) {
	return s.repo.GetPostsWithCommentsByThreadID(threadID)
}

func (s *postService) UpdatePost(post *models.Post, postID int) error {
	return s.repo.UpdatePost(post, postID)
}

func (s *postService) DeletePost(postID int) error {
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
