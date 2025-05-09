package main

import (
	"ForumService/internal/config"
	"ForumService/internal/handlers"
	"ForumService/internal/repository"
	"ForumService/internal/service"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "html/template"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig("../config/database.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Подключение к базе данных
	db, err := sql.Open("postgres", cfg.Database.GetDSN())
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Настройка пула соединений
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		log.Fatalf("Ошибка проверки подключения к базе данных: %v", err)
	}

	// Инициализация репозиториев
	threadRepo := repository.NewThreadRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	chatRepo := repository.NewChatMessageRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Инициализация сервисов
	threadService := service.NewThreadService(threadRepo, postRepo)
	postService := service.NewPostService(postRepo, commentRepo)
	commentService := service.NewCommentService(commentRepo)
	chatService := service.NewChatService(chatRepo)
	userService := service.NewUserService(userRepo)

	// Создание структуры сервисов для роутера
	services := &handlers.Services{
		ThreadService:  threadService,
		PostService:    postService,
		CommentService: commentService,
		UserService:    userService,
		ChatService:    chatService,
	}

	// Инициализация роутера
	router := gin.Default()

	// Загрузка шаблонов
	router.LoadHTMLGlob("../templates/*")
	router.Static("/static", "../static")

	// Регистрация маршрутов
	handlers.RegisterRoutes(router, services)

	// Запуск сервера
	fmt.Println("Сервер запущен на http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
