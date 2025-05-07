package handlers

import (
	"ForumService/internal/models"
	"ForumService/internal/service"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"fmt"
)

type ViewsHandler struct {
	threadService service.ThreadService
	commentService service.CommentService
}

func NewViewsHandler(threadService service.ThreadService, commentService service.CommentService) *ViewsHandler {
	return &ViewsHandler{
		threadService: threadService,
		commentService: commentService,
	}
}

func (h *ViewsHandler) Index(c *gin.Context) {
	threads, err := h.threadService.GetAllThreads()
	if err != nil {
		fmt.Printf("Error getting threads: %v\n", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Title":   "Ошибка",
			"Message": "Не удалось загрузить треды",
		})
		return
	}

	fmt.Printf("Rendering index with %d threads\n", len(threads))

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Threads": threads,
	})
}

func (h *ViewsHandler) GetThreadWithPosts(w http.ResponseWriter, r *http.Request) {
	// Получаем ID треда из URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	threadIDStr := parts[2]
	threadID, err := strconv.Atoi(threadIDStr)
	if err != nil {
		http.Error(w, "Invalid thread ID", http.StatusBadRequest)
		return
	}

	thread, posts, err := h.threadService.GetThreadWithPosts(threadID)
	if err != nil {
		http.Error(w, "Thread not found", http.StatusNotFound)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/thread.html"))
	data := struct {
		Thread *models.Thread
		Posts  []*models.Post
	}{
		Thread: thread,
		Posts:  posts,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h *ViewsHandler) ShowThread(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Title":   "Ошибка",
			"Message": "Неверный ID треда",
		})
		return
	}

	thread, posts, err := h.threadService.GetThreadWithPosts(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"Title":   "Ошибка",
			"Message": "Тред не найден",
		})
		return
	}

	userID := 0
	if user, exists := c.Get("user"); exists {
		if u, ok := user.(*models.User); ok {
			userID = u.ID
		}
	}

	// Получаем комментарии для каждого поста
	for i := range posts {
		comments, err := h.commentService.GetCommentsByPostID(posts[i].ID)
		if err != nil {
			comments = make([]models.Comment, 0)
		}
		posts[i].Comments = comments
		
		// Проверяем права на редактирование
		posts[i].CanEdit = posts[i].AuthorID == userID
	}

	c.HTML(http.StatusOK, "thread.html", gin.H{
		"Thread": thread,
		"Posts":  posts,
		"User":   c.MustGet("user"),
	})
}
