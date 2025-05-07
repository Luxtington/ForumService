package handlers

import (
	"ForumService/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CommentHandler struct {
	service service.CommentService
}

func NewCommentHandler(service service.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req struct {
		PostID  int    `form:"post_id" binding:"required"`
		Content string `form:"content" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Неверные данные",
		})
		return
	}

	userID, _ := c.Get("userID")
	_, err := h.service.CreateComment(req.PostID, userID.(int), req.Content)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка при создании комментария",
		})
		return
	}

	c.Redirect(http.StatusFound, "/posts/"+strconv.Itoa(req.PostID))
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Неверный ID комментария",
		})
		return
	}

	userID, _ := c.Get("userID")
	isAdmin, _ := c.Get("isAdmin")
	err = h.service.DeleteComment(id, userID.(int), isAdmin.(bool))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка при удалении комментария",
		})
		return
	}

	c.Redirect(http.StatusFound, "/posts")
}
