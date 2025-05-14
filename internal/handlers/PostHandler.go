package handlers

import (
	"ForumService/internal/models"
	"ForumService/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PostHandler struct {
	service service.PostService
}

func NewPostHandler(service service.PostService) *PostHandler {
	return &PostHandler{service: service}
}

func (h *PostHandler) GetAllPosts(c *gin.Context) {
	// TODO: Реализовать получение всех постов
	c.HTML(http.StatusOK, "posts.html", gin.H{
		"title": "Все посты",
	})
}

func (h *PostHandler) ShowCreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "create_post.html", gin.H{
		"title": "Создать пост",
	})
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var request struct {
		ThreadID int    `json:"thread_id" binding:"required"`
		Content  string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем ID пользователя из контекста
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не аутентифицирован"})
		return
	}

	post := &models.Post{
		ThreadID: request.ThreadID,
		AuthorID: userID.(int),
		Content:  request.Content,
	}

	if err := h.service.CreatePost(post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) GetPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	post, err := h.service.GetPostByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) GetPostWithComments(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	post, comments, err := h.service.GetPostWithComments(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post":     post,
		"comments": comments,
	})
}

func (h *PostHandler) ShowEditForm(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Неверный ID поста",
		})
		return
	}

	post, err := h.service.GetPostByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Пост не найден",
		})
		return
	}

	c.HTML(http.StatusOK, "edit_post.html", gin.H{
		"title": "Редактировать пост",
		"post":  post,
	})
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdatePost(&post, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	if err := h.service.DeletePost(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *PostHandler) ListPosts(c *gin.Context) {
	posts, err := h.service.GetAllPosts()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Title":   "Ошибка",
			"Message": "Не удалось загрузить посты",
		})
		return
	}

	c.HTML(http.StatusOK, "list_post.html", gin.H{
		"Posts": posts,
	})
}

func (h *PostHandler) ShowPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Title":   "Ошибка",
			"Message": "Неверный ID поста",
		})
		return
	}

	post, comments, err := h.service.GetPostWithComments(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"Title":   "Ошибка",
			"Message": "Пост не найден",
		})
		return
	}

	if comments == nil {
		comments = make([]models.Comment, 0)
	}

	c.HTML(http.StatusOK, "view_post.html", gin.H{
		"Post":     post,
		"Comments": comments,
		"User":     c.MustGet("user"),
	})
}

func (h *PostHandler) CreateComment(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Title":   "Ошибка",
			"Message": "Неверный ID поста",
		})
		return
	}

	content := c.PostForm("content")
	authorID := c.MustGet("userID").(int)

	comment := &models.Comment{
		Content:  content,
		AuthorID: authorID,
		PostID:   postID,
	}

	if err := h.service.CreateComment(comment); err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Title":   "Ошибка",
			"Message": "Не удалось создать комментарий",
		})
		return
	}

	c.Redirect(http.StatusFound, "/posts/"+strconv.Itoa(postID))
}

func (h *PostHandler) DeleteComment(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Title":   "Ошибка",
			"Message": "Неверный ID поста",
		})
		return
	}

	commentID, err := strconv.Atoi(c.Param("commentId"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Title":   "Ошибка",
			"Message": "Неверный ID комментария",
		})
		return
	}

	comment, err := h.service.GetCommentByID(commentID)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"Title":   "Ошибка",
			"Message": "Комментарий не найден",
		})
		return
	}

	if comment.AuthorID != c.MustGet("userID").(int) {
		c.HTML(http.StatusForbidden, "error.html", gin.H{
			"Title":   "Ошибка",
			"Message": "У вас нет прав для удаления этого комментария",
		})
		return
	}

	if err := h.service.DeleteComment(commentID); err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Title":   "Ошибка",
			"Message": "Не удалось удалить комментарий",
		})
		return
	}

	c.Redirect(http.StatusFound, "/posts/"+strconv.Itoa(postID))
}

func (h *PostHandler) GetPostComments(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	_, comments, err := h.service.GetPostWithComments(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	c.JSON(http.StatusOK, comments)
}
