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
	fmt.Printf("Начало обработки запроса ShowThread\n")
	
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Printf("Ошибка при парсинге ID треда: %v\n", err)
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Неверный ID треда",
		})
		return
	}
	fmt.Printf("Получение треда с ID: %d\n", id)

	thread, posts, err := h.threadService.GetThreadWithPosts(id)
	if err != nil {
		fmt.Printf("Ошибка при получении треда: %v\n", err)
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Тред не найден",
		})
		return
	}
	fmt.Printf("Тред найден: %+v\n", thread)
	fmt.Printf("Количество постов: %d\n", len(posts))

	c.HTML(http.StatusOK, "thread.html", gin.H{
		"Thread": thread,
		"Posts":  posts,
	})
}
