package main

import (
	"ForumService/internal/handlers"
	"ForumService/internal/middleware"
	"ForumService/internal/repository"
	"ForumService/internal/service"
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"strconv"
)

func main() {
	// Подключение к базе данных
	dsn := "host=localhost user=postgres password=postgres dbname=forum port=5432 sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Проверка подключения
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Инициализация репозиториев
	threadRepo := repository.NewThreadRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	chatRepo := repository.NewChatMessageRepository(db)

	// Инициализация сервисов
	threadService := service.NewThreadService(threadRepo, postRepo)
	postService := service.NewPostService(postRepo, commentRepo)
	commentService := service.NewCommentService(commentRepo)
	chatService := service.NewChatService(chatRepo)

	// Инициализация обработчиков
	threadHandler := handlers.NewThreadHandler(threadService)
	postHandler := handlers.NewPostHandler(postService)
	commentHandler := handlers.NewCommentHandler(commentService)
	chatHandler := handlers.NewChatHandler(chatService)

	// Создание экземпляра Gin
	r := gin.Default()

	// Загрузка HTML шаблонов
	r.LoadHTMLGlob("templates/*")
	// Настройка статических файлов
	r.Static("/static", "./static")

	// Инициализация middleware для аутентификации
	authMiddleware := middleware.AuthServiceMiddleware("http://localhost:8081")

	// Группа защищенных маршрутов
	protected := r.Group("/api")
	protected.Use(authMiddleware)
	{
		// Маршруты для тредов
		protected.POST("/threads", threadHandler.CreateThread)
		protected.PUT("/threads/:id", threadHandler.UpdateThread)
		protected.DELETE("/threads/:id", threadHandler.DeleteThread)

		// Маршруты для постов
		protected.POST("/posts", postHandler.CreatePost)
		protected.PUT("/posts/:id", postHandler.UpdatePost)
		protected.DELETE("/posts/:id", postHandler.DeletePost)

		// Маршруты для комментариев
		protected.POST("/comments", commentHandler.CreateComment)
		protected.DELETE("/comments/:id", commentHandler.DeleteComment)

		// Маршруты для чата
		protected.POST("/chat", chatHandler.CreateMessage)
		protected.GET("/chat", chatHandler.GetMessages)
	}

	// Публичные маршруты
	// Главная страница со списком тредов
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "Главная страница",
		})
	})

	// Получение всех тредов (HTML)
	r.GET("/threads", func(c *gin.Context) {
		threads, err := threadService.GetAllThreads()
		if err != nil {
			c.HTML(500, "error.html", gin.H{
				"error": err.Error(),
			})
			return
		}
		c.HTML(200, "threads.html", gin.H{
			"threads": threads,
		})
	})

	// Получение конкретного треда с постами (HTML)
	r.GET("/threads/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.HTML(400, "bad_request.html", gin.H{
				"error": "Неверный ID треда",
			})
			return
		}
		thread, posts, err := threadService.GetThreadWithPosts(id)
		if err != nil {
			c.HTML(404, "not_found.html", gin.H{
				"error": "Тред не найден",
			})
			return
		}
		c.HTML(200, "thread.html", gin.H{
			"thread": thread,
			"posts":  posts,
		})
	})

	// Запуск сервера
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
