package handlers

import (
	"ForumService/internal/service"
	"ForumService/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
	"strings"
	"fmt"
	"unicode"
)

type ThreadHandler struct {
	service service.ThreadService
}

func NewThreadHandler(service service.ThreadService) *ThreadHandler {
	return &ThreadHandler{service: service}
}

type CreateThreadRequest struct {
	Title string `json:"title" binding:"required"`
}

type UpdateThreadRequest struct {
	Title string `json:"title" binding:"required"`
}

// func (h *ThreadHandler) RegisterRoutes(r *gin.RouterGroup) {
// 	threads := r.Group("/threads")
// 	{
// 		threads.POST("", h.CreateThread)
// 		threads.GET("/:id", h.GetThreadWithPosts)
// 		threads.DELETE("/:id", h.DeleteThread)
// 	}
// }

func (h *ThreadHandler) CreateThread(c *gin.Context) {
	var request struct {
		Title string `json:"title" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "неверный формат данных"})
		return
	}

	// Получаем ID пользователя из контекста
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "пользователь не аутентифицирован"})
		return
	}

	// Преобразуем uint32 в int
	userIDInt := int(userID.(uint32))

	thread, err := h.service.CreateThread(request.Title, userIDInt)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, thread)
}

func (h *ThreadHandler) GetThreadWithPosts(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid thread ID"})
		return
	}

	thread, posts, err := h.service.GetThreadWithPosts(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "thread not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"thread": thread,
		"posts":  posts,
	})
}

func (h *ThreadHandler) DeleteThread(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid thread ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userIDInt := int(userID.(uint32))

	// Получаем роль пользователя
	userRole, _ := c.Get("user_role")
	if userRole == nil {
		userRole = "user"
	}

	// Проверяем, является ли пользователь автором треда или администратором
	thread, _, err := h.service.GetThreadWithPosts(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при получении треда"})
		return
	}

	if thread.AuthorID != userIDInt && userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "no permission to delete this thread"})
			return
		}

	if err := h.service.DeleteThread(id, userIDInt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ThreadHandler) UpdateThread(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid thread ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userIDInt := int(userID.(uint32))

	// Получаем роль пользователя
	userRole, _ := c.Get("user_role")
	if userRole == nil {
		userRole = "user"
	}

	var req UpdateThreadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем, является ли пользователь автором треда или администратором
	thread, _, err := h.service.GetThreadWithPosts(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при получении треда"})
		return
	}

	if thread.AuthorID != userIDInt && userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "no permission to update this thread"})
			return
		}

	thread.Title = req.Title
	if err := h.service.UpdateThread(thread, userIDInt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, thread)
}

func (h *ThreadHandler) GetAllThreads(c *gin.Context) {
	threads, err := h.service.GetAllThreads()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Добавляем информацию об авторе для каждого треда
	for _, thread := range threads {
		user, err := h.service.GetUserByID(thread.AuthorID)
		if err == nil && user != nil {
			thread.AuthorName = user.Username
		}
	}

	c.JSON(http.StatusOK, threads)
}

func (h *ThreadHandler) GetThreadPosts(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid thread ID"})
		return
	}

	posts, err := h.service.GetPostsByThreadID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "thread not found"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// formatDate форматирует дату в строку
func formatDate(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	return date.Format("02.01.2006 15:04")
}

// validateThreadTitle проверяет валидность заголовка темы
func validateThreadTitle(title string) bool {
	if title == "" || strings.TrimSpace(title) == "" {
		return false
	}
	if len(title) > 255 {
		return false
	}
	return true
}

// sanitizeThreadTitle очищает заголовок темы от лишних пробелов и HTML
func sanitizeThreadTitle(title string) string {
	// Удаляем HTML теги
	title = strings.ReplaceAll(title, "<", "&lt;")
	title = strings.ReplaceAll(title, ">", "&gt;")
	
	// Удаляем специальные символы
	title = strings.ReplaceAll(title, "\n", " ")
	title = strings.ReplaceAll(title, "\t", " ")
	
	// Удаляем лишние пробелы
	return strings.TrimSpace(title)
}

// getThreadStatus определяет статус темы на основе времени создания и обновления
func getThreadStatus(thread *models.Thread) string {
	if thread == nil {
		return "unknown"
	}
	
	// Если тема обновлялась менее 24 часов назад, считаем её активной
	if time.Since(thread.UpdatedAt) < 24*time.Hour {
		return "active"
	}
	return "inactive"
}

// ThreadStats содержит статистику по теме
type ThreadStats struct {
	TotalPosts    int
	UniqueAuthors int
}

// calculateThreadStats подсчитывает статистику по теме
func calculateThreadStats(thread *models.Thread, posts []*models.Post) ThreadStats {
	if thread == nil || len(posts) == 0 {
		return ThreadStats{}
	}
	
	stats := ThreadStats{
		TotalPosts: len(posts),
	}
	
	// Подсчитываем уникальных авторов
	authors := make(map[int]bool)
	for _, post := range posts {
		authors[post.AuthorID] = true
	}
	stats.UniqueAuthors = len(authors)
	
	return stats
}

// ThreadMetrics содержит метрики темы
type ThreadMetrics struct {
	AveragePostLength float64
	MostActiveAuthor  int
	LastActivityTime  time.Time
}

// calculateThreadMetrics вычисляет метрики темы
func calculateThreadMetrics(posts []*models.Post) ThreadMetrics {
	if len(posts) == 0 {
		return ThreadMetrics{}
	}

	metrics := ThreadMetrics{
		LastActivityTime: posts[0].CreatedAt,
	}

	// Считаем среднюю длину постов
	totalLength := 0
	authorActivity := make(map[int]int)

	for _, post := range posts {
		totalLength += len(post.Content)
		authorActivity[post.AuthorID]++

		if post.CreatedAt.After(metrics.LastActivityTime) {
			metrics.LastActivityTime = post.CreatedAt
		}
	}

	metrics.AveragePostLength = float64(totalLength) / float64(len(posts))

	// Находим самого активного автора
	maxActivity := 0
	for authorID, activity := range authorActivity {
		if activity > maxActivity {
			maxActivity = activity
			metrics.MostActiveAuthor = authorID
		}
	}

	return metrics
}

// isThreadActive проверяет, активна ли тема
func isThreadActive(thread *models.Thread, posts []*models.Post) bool {
	if thread == nil || len(posts) == 0 {
		return false
	}

	// Тема считается активной, если:
	// 1. Ей менее 7 дней
	// 2. Есть посты за последние 24 часа
	threadAge := time.Since(thread.CreatedAt)
	if threadAge > 7*24*time.Hour {
		return false
	}

	lastPostTime := posts[0].CreatedAt
	for _, post := range posts {
		if post.CreatedAt.After(lastPostTime) {
			lastPostTime = post.CreatedAt
		}
	}

	return time.Since(lastPostTime) < 24*time.Hour
}

// formatThreadSummary создает краткое описание темы
func formatThreadSummary(thread *models.Thread, posts []*models.Post) string {
	if thread == nil {
		return "Тема не найдена"
	}

	metrics := calculateThreadMetrics(posts)
	status := "активна"
	if !isThreadActive(thread, posts) {
		status = "неактивна"
	}

	return fmt.Sprintf("Тема '%s' (%s). Постов: %d, средняя длина: %.1f символов",
		thread.Title,
		status,
		len(posts),
		metrics.AveragePostLength)
}

// validateThreadContent проверяет содержимое темы на спам
func validateThreadContent(content string) (bool, string) {
	if content == "" {
		return false, "сообщение не может быть пустым"
	}

	if content == "Коротко" {
		return false, "сообщение слишком короткое"
	}

	if len(content) > 10000 {
		return false, "сообщение слишком длинное"
	}

	if strings.Contains(content, "!!!!!") || strings.Contains(content, "?????") {
		return false, "слишком много повторяющихся символов"
	}

	upperCount := 0
	totalLetters := 0
	for _, char := range content {
		if unicode.IsLetter(char) {
			totalLetters++
			if unicode.IsUpper(char) {
				upperCount++
			}
		}
	}

	if totalLetters > 0 && float64(upperCount)/float64(totalLetters) > 0.7 {
		return false, "слишком много заглавных букв"
	}

	return true, ""
}
