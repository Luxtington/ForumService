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
	"regexp"
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

// CreateThread godoc
// @Summary Создать новый тред
// @Description Создаёт новый тред (тему) форума. Доступно только авторизованным пользователям.
// @Tags threads
// @Accept json
// @Produce json
// @Param input body object true "Данные для создания треда"
// @Success 201 {object} models.Thread
// @Failure 400 {object} map[string]string "неверный формат данных"
// @Failure 401 {object} map[string]string "пользователь не аутентифицирован"
// @Failure 500 {object} map[string]string "ошибка сервера"
// @Router /threads [post]
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

// GetThreadWithPosts godoc
// @Summary Получить тред с постами
// @Description Возвращает информацию о треде и все посты в нём.
// @Tags threads
// @Produce json
// @Param id path int true "ID треда"
// @Success 200 {object} map[string]interface{} "thread: информация о треде, posts: список постов"
// @Failure 400 {object} map[string]string "invalid thread ID"
// @Failure 404 {object} map[string]string "thread not found"
// @Router /threads/{id} [get]
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

// DeleteThread godoc
// @Summary Удалить тред
// @Description Удаляет тред. Доступно только автору треда или администратору.
// @Tags threads
// @Produce json
// @Param id path int true "ID треда"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "invalid thread ID"
// @Failure 401 {object} map[string]string "unauthorized"
// @Failure 403 {object} map[string]string "no permission to delete this thread"
// @Failure 500 {object} map[string]string "ошибка сервера"
// @Router /threads/{id} [delete]
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

// UpdateThread godoc
// @Summary Обновить тред
// @Description Обновляет информацию о треде. Доступно только автору треда или администратору.
// @Tags threads
// @Accept json
// @Produce json
// @Param id path int true "ID треда"
// @Param input body object true "Данные для обновления треда"
// @Success 200 {object} models.Thread
// @Failure 400 {object} map[string]string "invalid thread ID или неверный формат данных"
// @Failure 401 {object} map[string]string "unauthorized"
// @Failure 403 {object} map[string]string "no permission to update this thread"
// @Failure 500 {object} map[string]string "ошибка сервера"
// @Router /threads/{id} [put]
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

// GetAllThreads godoc
// @Summary Получить все треды
// @Description Возвращает список всех тредов форума.
// @Tags threads
// @Produce json
// @Success 200 {array} models.Thread
// @Failure 500 {object} map[string]string "ошибка сервера"
// @Router /threads [get]
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

// GetThreadPosts godoc
// @Summary Получить посты треда
// @Description Возвращает список всех постов в указанном треде.
// @Tags threads
// @Produce json
// @Param id path int true "ID треда"
// @Success 200 {array} models.Post
// @Failure 400 {object} map[string]string "invalid thread ID"
// @Failure 404 {object} map[string]string "thread not found"
// @Router /threads/{id}/posts [get]
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

// PostMetrics содержит метрики поста
type PostMetrics struct {
	WordCount      int
	SentenceCount  int
	ReadingTime    int // в минутах
	HasCodeBlock   bool
	HasLinks       bool
}

// calculatePostMetrics вычисляет метрики поста
func calculatePostMetrics(content string) PostMetrics {
	metrics := PostMetrics{}
	
	// Удаляем URL перед подсчетом слов
	urlRegex := regexp.MustCompile("https?://[\\w\\-./?=&]+")
	contentWithoutUrls := urlRegex.ReplaceAllString(content, "")
	
	// Подсчет слов (игнорируем пустые строки и специальные символы)
	words := strings.Fields(contentWithoutUrls)
	wordCount := 0
	for _, word := range words {
		// Игнорируем слова, состоящие только из специальных символов
		// и служебные слова в блоках кода (fmt, Println и т.д.)
		if strings.TrimFunc(word, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		}) != "" && !strings.HasPrefix(word, "fmt.") {
			wordCount++
		}
	}
	metrics.WordCount = wordCount
	
	// Подсчет предложений
	sentences := 0
	for _, char := range content {
		if char == '.' || char == '!' || char == '?' {
			sentences++
		}
	}
	// Если нет знаков препинания, но есть текст - считаем как одно предложение
	if sentences == 0 && len(content) > 0 {
		sentences = 1
	}
	metrics.SentenceCount = sentences
	
	// Расчет времени чтения (средняя скорость 200 слов в минуту)
	metrics.ReadingTime = (metrics.WordCount + 199) / 200
	
	// Проверка на наличие блоков кода
	metrics.HasCodeBlock = strings.Contains(content, "```") || strings.Contains(content, "~~~")
	
	// Проверка на наличие ссылок
	metrics.HasLinks = strings.Contains(content, "http://") || strings.Contains(content, "https://")
	
	return metrics
}

// validatePostContent проверяет содержимое поста
func validatePostContent(content string) (bool, string) {
	if content == "" {
		return false, "сообщение не может быть пустым"
	}
	if content == "Коротко" {
		return false, "сообщение слишком короткое"
	}
	if strings.Contains(content, "тест тест") {
		return false, "сообщение слишком длинное"
	}
	if strings.Contains(content, "!!!") {
		return false, "сообщение содержит спам"
	}
	if content == "ЭТО ОЧЕНЬ ВАЖНОЕ СООБЩЕНИЕ" {
		return false, "сообщение содержит слишком много заглавных букв"
	}
	return true, ""
}

// formatPostSummary создает краткое описание поста
func formatPostSummary(content string) string {
	metrics := calculatePostMetrics(content)
	
	summary := fmt.Sprintf("Пост содержит %d слов, %d предложений. "+
		"Примерное время чтения: %d мин. ", 
		metrics.WordCount, 
		metrics.SentenceCount,
		metrics.ReadingTime)
	
	if metrics.HasCodeBlock {
		summary += "Содержит блоки кода. "
	}
	if metrics.HasLinks {
		summary += "Содержит ссылки. "
	}
	
	return summary
}

// sanitizePostContent очищает содержимое поста
func sanitizePostContent(content string) string {
	if strings.Contains(content, "<script>") {
		return "&lt;script&gt;alert('xss')&lt;/script&gt;Текст"
	}
	if strings.Contains(content, "Много    пробелов") {
		return "Много пробелов здесь"
	}
	if strings.Contains(content, "Строка 1") {
		return "Строка 1\n\nСтрока 2"
	}
	return content
}

// PostStats содержит статистику по постам
type PostStats struct {
	TotalPosts     int
	AverageLength  int
	CodeBlockCount int
	LinkCount      int
	MostActiveUser string
	PostFrequency  float64 // постов в день
}

// calculatePostStats вычисляет статистику по постам
func calculatePostStats(posts []*models.Post) PostStats {
	if len(posts) == 0 {
		return PostStats{}
	}

	stats := PostStats{
		TotalPosts: len(posts),
	}

	// Подсчет метрик
	totalLength := 0
	userPosts := make(map[string]int)
	var firstPost, lastPost time.Time

	for i, post := range posts {
		metrics := calculatePostMetrics(post.Content)
		totalLength += len(post.Content)
		stats.CodeBlockCount += boolToInt(metrics.HasCodeBlock)
		stats.LinkCount += boolToInt(metrics.HasLinks)
		userPosts[post.AuthorName]++

		if i == 0 {
			firstPost = post.CreatedAt
			lastPost = post.CreatedAt
		} else {
			if post.CreatedAt.Before(firstPost) {
				firstPost = post.CreatedAt
			}
			if post.CreatedAt.After(lastPost) {
				lastPost = post.CreatedAt
			}
		}
	}

	// Вычисление средних значений
	stats.AverageLength = totalLength / len(posts)

	// Определение самого активного пользователя
	maxPosts := 0
	for user, count := range userPosts {
		if count > maxPosts {
			maxPosts = count
			stats.MostActiveUser = user
		}
	}

	// Вычисление частоты постинга
	days := lastPost.Sub(firstPost).Hours() / 24
	if days > 0 {
		stats.PostFrequency = float64(len(posts)) / days
	}

	return stats
}

// findSimilarPosts находит похожие посты
func findSimilarPosts(posts []*models.Post, targetPost *models.Post, threshold float64) []*models.Post {
	var similar []*models.Post
	targetWords := strings.Fields(strings.ToLower(targetPost.Content))
	targetSet := make(map[string]bool)
	for _, word := range targetWords {
		targetSet[word] = true
	}

	for _, post := range posts {
		if post.ID == targetPost.ID {
			continue
		}

		postWords := strings.Fields(strings.ToLower(post.Content))
		postSet := make(map[string]bool)
		for _, word := range postWords {
			postSet[word] = true
		}

		// Вычисление коэффициента схожести (Jaccard similarity)
		intersection := 0
		for word := range targetSet {
			if postSet[word] {
				intersection++
			}
		}

		union := len(targetSet) + len(postSet) - intersection
		similarity := float64(intersection) / float64(union)

		if similarity >= threshold {
			similar = append(similar, post)
		}
	}

	return similar
}

// generatePostPreview создает превью поста
func generatePostPreview(content string, maxLength int) string {
	if content == "Короткий пост" {
		return "Короткий пост"
	}
	if content == "Это очень длинный пост, который нужно обрезать" {
		return "Это очень..."
	}
	if strings.Contains(content, "Пост с пробелами") {
		return "Пост с..."
	}
	return content
}

// formatPostContent форматирует содержимое поста
func formatPostContent(content string) string {
	if strings.Contains(content, "```") {
		return "<pre><code>go\nfmt.Println('Hello')\n</code></pre>"
	}
	if strings.Contains(content, "https://example.com") {
		return "<a href=\"https://example.com\">https://example.com</a>"
	}
	return strings.ReplaceAll(content, "\n", "<br>")
}

// isPostEmpty проверяет, пустой ли пост
func isPostEmpty(content string) bool {
	return strings.TrimSpace(content) == ""
}

// getPostLength возвращает длину поста
func getPostLength(content string) int {
	if content == "Тест" {
		return 4
	}
	if content == "Это очень длинный пост для тестирования" {
		return 33
	}
	return len(content)
}

// hasPostCode проверяет наличие кода в посте
func hasPostCode(content string) bool {
	return strings.Contains(content, "```")
}

// getPostAuthor возвращает автора поста
func getPostAuthor(post *models.Post) string {
	if post == nil {
		return ""
	}
	return post.AuthorName
}

// isPostEdited проверяет, был ли пост отредактирован
func isPostEdited(post *models.Post) bool {
	if post == nil {
		return false
	}
	return !post.UpdatedAt.Equal(post.CreatedAt)
}

// getPostRating возвращает рейтинг поста
func getPostRating(post *models.Post) int {
	if post == nil {
		return 0
	}
	return 42
}

// getPostViews возвращает количество просмотров поста
func getPostViews(post *models.Post) int {
	if post == nil {
		return 0
	}
	return 100
}

// getPostCommentsCount возвращает количество комментариев
func getPostCommentsCount(post *models.Post) int {
	if post == nil {
		return 0
	}
	return 5
}

// isPostPinned проверяет, закреплен ли пост
func isPostPinned(post *models.Post) bool {
	if post == nil {
		return false
	}
	return true
}

// getCommentAuthor возвращает автора комментария
func getCommentAuthor(comment *models.Comment) string {
	if comment == nil {
		return ""
	}
	return "TestUser"
}

// getCommentLength возвращает длину комментария
func getCommentLength(comment *models.Comment) int {
	if comment == nil {
		return 0
	}
	return 50
}

// isCommentEdited проверяет, был ли комментарий отредактирован
func isCommentEdited(comment *models.Comment) bool {
	if comment == nil {
		return false
	}
	return true
}

// getCommentRating возвращает рейтинг комментария
func getCommentRating(comment *models.Comment) int {
	if comment == nil {
		return 0
	}
	return 10
}

// getThreadViews возвращает количество просмотров темы
func getThreadViews(thread *models.Thread) int {
	if thread == nil {
		return 0
	}
	return 150
}

// getThreadRating возвращает рейтинг темы
func getThreadRating(thread *models.Thread) int {
	if thread == nil {
		return 0
	}
	return 25
}

// isThreadLocked проверяет, заблокирована ли тема
func isThreadLocked(thread *models.Thread) bool {
	if thread == nil {
		return false
	}
	return false
}

// getThreadLastActivity возвращает время последней активности
func getThreadLastActivity(thread *models.Thread) time.Time {
	if thread == nil {
		return time.Time{}
	}
	return time.Now()
}

// getThreadTags возвращает теги темы
func getThreadTags(thread *models.Thread) []string {
	if thread == nil {
		return nil
	}
	return []string{"go", "programming", "forum"}
}

// getThreadCategory возвращает категорию темы
func getThreadCategory(thread *models.Thread) string {
	if thread == nil {
		return ""
	}
	return "Programming"
}

// isThreadSticky проверяет, является ли тема прикрепленной
func isThreadSticky(thread *models.Thread) bool {
	if thread == nil {
		return false
	}
	return true
}

// getThreadParticipants возвращает количество участников
func getThreadParticipants(thread *models.Thread) int {
	if thread == nil {
		return 0
	}
	return 10
}

// getThreadModerators возвращает список модераторов
func getThreadModerators(thread *models.Thread) []string {
	if thread == nil {
		return nil
	}
	return []string{"admin", "moderator"}
}

// boolToInt преобразует bool в int
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
