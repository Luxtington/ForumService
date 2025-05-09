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

	// Инициализация сервисов
	threadService := service.NewThreadService(threadRepo, postRepo)
	postService := service.NewPostService(postRepo, commentRepo)
	commentService := service.NewCommentService(commentRepo)

	// Инициализация обработчиков
	threadHandler := handlers.NewThreadHandler(threadService)
	postHandler := handlers.NewPostHandler(postService)
	commentHandler := handlers.NewCommentHandler(commentService)

	// Создание экземпляра Gin
	r := gin.Default()

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
	}

	// Публичные маршруты
	// Получение всех тредов
	r.GET("/threads", threadHandler.GetAllThreads)
	// Получение конкретного треда с постами
	r.GET("/threads/:id", threadHandler.GetThreadWithPosts)

	// Запуск сервера
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
