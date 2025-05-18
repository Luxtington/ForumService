package repository

import "ForumService/internal/models"

type CommentRepository interface {
	SaveComment(comment *models.Comment) error
	GetCommentByID(id int) (*models.Comment, error)
	DeleteComment(id int) error
	GetCommentsByPostID(postID int) ([]models.Comment, error)
}

type ThreadRepository interface {
	Create(thread *models.Thread) error
	GetByID(id int) (*models.Thread, error)
	Update(thread *models.Thread) error
	Delete(id int) error
	GetAllThreads() ([]*models.Thread, error)
	GetThreadWithPosts(threadID int) (*models.Thread, []models.Post, map[int][]models.Comment, error)
}

type PostRepository interface {
	SavePost(post *models.Post) error
	GetPostByID(id int) (*models.Post, error)
	GetPostWithComments(postID int) (*models.Post, []models.Comment, error)
	GetPostsWithCommentsByThreadID(threadID int) ([]models.Post, map[int][]models.Comment, error)
	UpdatePost(post *models.Post, postID int) error
	DeletePost(postID int) error
	GetByThreadID(threadID int) ([]*models.Post, error)
}

type UserRepository interface {
	SaveUser(user *models.User) error
	GetUserByID(id int) (*models.User, error)
	GetUserPosts(userID int) ([]*models.Post, error)
	GetUserCommentCount(userID int) (int, error)
	GetUserRole(userID int) (string, error)
} 