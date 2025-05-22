package handlers

import (
	"ForumService/internal/models"
	"ForumService/internal/service"
	"ForumService/internal/errors"
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

type CreatePostRequest struct {
	ThreadID int    `json:"thread_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type UpdatePostRequest struct {
	Content string `json:"content" binding:"required"`
}

// GetAllPosts godoc
// @Summary Получить все посты
// @Description Возвращает список всех постов форума.
// @Tags posts
// @Produce json
// @Success 200 {array} models.Post
// @Failure 500 {object} map[string]string "ошибка сервера"
// @Router /posts [get]
func (h *PostHandler) GetAllPosts(c *gin.Context) {
	posts, err := h.service.GetAllPosts()
	if err != nil {
		c.Error(errors.NewInternalServerError("Ошибка при получении списка постов", err))
		return
	}

	c.JSON(http.StatusOK, posts)
}

// ShowCreateForm godoc
// @Summary Показать форму создания поста
// @Description Отображает HTML-страницу с формой для создания нового поста.
// @Tags posts
// @Produce html
// @Success 200 {string} string "HTML страница"
// @Router /posts/create [get]
func (h *PostHandler) ShowCreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "create_post.html", gin.H{
		"title": "Создать пост",
	})
}

// CreatePost godoc
// @Summary Создать новый пост
// @Description Создаёт новый пост в указанном треде. Доступно только авторизованным пользователям.
// @Tags posts
// @Accept json
// @Produce json
// @Param input body object true "Данные для создания поста"
// @Success 201 {object} models.Post
// @Failure 400 {object} map[string]string "неверный формат данных"
// @Failure 401 {object} map[string]string "пользователь не аутентифицирован"
// @Failure 403 {object} map[string]string "нет прав для создания поста в этом треде"
// @Failure 500 {object} map[string]string "ошибка сервера"
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	var request CreatePostRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(errors.NewValidationError("Неверный формат данных", err))
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.Error(errors.NewUnauthorizedError("Пользователь не аутентифицирован", nil))
		return
	}

	userIDInt := int(userID.(uint32))
	post := &models.Post{
		ThreadID: request.ThreadID,
		AuthorID: userIDInt,
		Content:  request.Content,
	}

	if err := h.service.CreatePost(post); err != nil {
		c.Error(errors.NewInternalServerError("Ошибка при создании поста", err))
		return
	}

	c.JSON(http.StatusCreated, post)
}

// GetPost godoc
// @Summary Получить пост по ID
// @Description Возвращает информацию о посте по его ID.
// @Tags posts
// @Produce json
// @Param id path int true "ID поста"
// @Success 200 {object} models.Post
// @Failure 400 {object} map[string]string "invalid post ID"
// @Failure 404 {object} map[string]string "post not found"
// @Router /posts/{id} [get]
func (h *PostHandler) GetPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.NewBadRequestError("Неверный ID поста", err))
		return
	}

	post, err := h.service.GetPost(id)
	if err != nil {
		c.Error(errors.NewNotFoundError("Пост не найден", err))
		return
	}

	c.JSON(http.StatusOK, post)
}

// GetPostWithComments godoc
// @Summary Получить пост с комментариями
// @Description Возвращает информацию о посте и все комментарии к нему.
// @Tags posts
// @Produce json
// @Param id path int true "ID поста"
// @Success 200 {object} map[string]interface{} "post: информация о посте, comments: список комментариев"
// @Failure 400 {object} map[string]string "invalid post ID"
// @Failure 404 {object} map[string]string "post not found"
// @Router /posts/{id}/comments [get]
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

// ShowEditForm godoc
// @Summary Показать форму редактирования поста
// @Description Отображает HTML-страницу с формой для редактирования поста.
// @Tags posts
// @Produce html
// @Param id path int true "ID поста"
// @Success 200 {string} string "HTML страница"
// @Failure 400 {object} map[string]string "Неверный ID поста"
// @Failure 404 {object} map[string]string "Пост не найден"
// @Router /posts/{id}/edit [get]
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

// UpdatePost godoc
// @Summary Обновить пост
// @Description Обновляет информацию о посте. Доступно только автору поста или администратору.
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "ID поста"
// @Param input body object true "Данные для обновления поста"
// @Success 200 {object} models.Post
// @Failure 400 {object} map[string]string "invalid post ID или неверный формат данных"
// @Failure 401 {object} map[string]string "unauthorized"
// @Failure 403 {object} map[string]string "no permission to update this post"
// @Failure 500 {object} map[string]string "ошибка сервера"
// @Router /posts/{id} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.NewBadRequestError("Неверный ID поста", err))
		return
	}

	var request UpdatePostRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(errors.NewValidationError("Неверный формат данных", err))
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.Error(errors.NewUnauthorizedError("Пользователь не аутентифицирован", nil))
		return
	}

	userIDInt := int(userID.(uint32))
	userRole, _ := c.Get("user_role")

	post, err := h.service.GetPost(id)
	if err != nil {
		c.Error(errors.NewNotFoundError("Пост не найден", err))
		return
	}

	if post.AuthorID != userIDInt && userRole != "admin" {
		c.Error(errors.NewPermissionDeniedError("Нет прав для редактирования поста", nil))
		return
	}

	post.Content = request.Content
	if err := h.service.UpdatePost(post, id, userIDInt); err != nil {
		c.Error(errors.NewInternalServerError("Ошибка при обновлении поста", err))
		return
	}

	c.JSON(http.StatusOK, post)
}

// DeletePost godoc
// @Summary Удалить пост
// @Description Удаляет пост. Доступно только автору поста или администратору.
// @Tags posts
// @Produce json
// @Param id path int true "ID поста"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "invalid post ID"
// @Failure 401 {object} map[string]string "unauthorized"
// @Failure 403 {object} map[string]string "no permission to delete this post"
// @Failure 500 {object} map[string]string "ошибка сервера"
// @Router /posts/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.NewBadRequestError("Неверный ID поста", err))
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.Error(errors.NewUnauthorizedError("Пользователь не аутентифицирован", nil))
		return
	}

	userIDInt := int(userID.(uint32))
	userRole, _ := c.Get("user_role")

	post, err := h.service.GetPost(id)
	if err != nil {
		c.Error(errors.NewNotFoundError("Пост не найден", err))
		return
	}

	if post.AuthorID != userIDInt && userRole != "admin" {
		c.Error(errors.NewPermissionDeniedError("Нет прав для удаления поста", nil))
		return
	}

	if err := h.service.DeletePost(id, userIDInt); err != nil {
		c.Error(errors.NewInternalServerError("Ошибка при удалении поста", err))
		return
	}

	c.Status(http.StatusNoContent)
}

// ListPosts godoc
// @Summary Список постов
// @Description Отображает HTML-страницу со списком всех постов.
// @Tags posts
// @Produce html
// @Success 200 {string} string "HTML страница"
// @Failure 500 {object} map[string]string "ошибка сервера"
// @Router /posts [get]
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

// ShowPost godoc
// @Summary Показать пост
// @Description Отображает HTML-страницу с информацией о посте.
// @Tags posts
// @Produce html
// @Param id path int true "ID поста"
// @Success 200 {string} string "HTML страница"
// @Failure 400 {object} map[string]string "Неверный ID поста"
// @Failure 404 {object} map[string]string "Пост не найден"
// @Router /posts/{id} [get]
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
		c.Error(errors.NewBadRequestError("Неверный ID поста", err))
		return
	}

	post, comments, err := h.service.GetPostWithComments(id)
	if err != nil {
		c.Error(errors.NewNotFoundError("Пост не найден", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
		"comments": comments,
	})
}
