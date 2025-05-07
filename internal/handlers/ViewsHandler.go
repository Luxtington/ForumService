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
	threadService  service.ThreadService
	postService    service.PostService
	commentService service.CommentService
	chatService    service.ChatService
}

func NewViewsHandler(
	threadService service.ThreadService,
	postService service.PostService,
	commentService service.CommentService,
	chatService service.ChatService,
) *ViewsHandler {
	return &ViewsHandler{
		threadService:  threadService,
		postService:    postService,
		commentService: commentService,
		chatService:    chatService,
	}
}

func (h *ViewsHandler) Index(c *gin.Context) {
	threads, err := h.threadService.GetAllThreads()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка при получении списка тредов",
		})
		return
	}

	chatMessages, err := h.chatService.GetAllMessages()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка при получении сообщений чата",
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Threads":      threads,
		"ChatMessages": chatMessages,
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

func (h *ViewsHandler) ShowPost(c *gin.Context) {
	fmt.Printf("Начало обработки запроса ShowPost\n")
	
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Printf("Ошибка при парсинге ID поста: %v\n", err)
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Неверный ID поста",
		})
		return
	}
	fmt.Printf("Получение поста с ID: %d\n", postID)

	post, err := h.postService.GetPost(postID)
	if err != nil {
		fmt.Printf("Ошибка при получении поста: %v\n", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": fmt.Sprintf("Ошибка при получении поста: %v", err),
		})
		return
	}
	fmt.Printf("Пост найден: %+v\n", post)

	comments, err := h.commentService.GetCommentsByPostID(postID)
	if err != nil {
		fmt.Printf("Ошибка при получении комментариев: %v\n", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": fmt.Sprintf("Ошибка при получении комментариев: %v", err),
		})
		return
	}
	fmt.Printf("Получено комментариев: %d\n", len(comments))

	c.HTML(http.StatusOK, "post.html", gin.H{
		"Post":     post,
		"Comments": comments,
	})
}
