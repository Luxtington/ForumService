package main

import (
	"ForumService/internal/client"
	"ForumService/internal/handlers"
	"ForumService/internal/middleware"
	"ForumService/internal/repository"
	"ForumService/internal/service"
	_"context"
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/Luxtington/Shared/logger"
	"net/http"
	"strconv"
	_"time"
	"github.com/gorilla/websocket"
	_"github.com/golang/protobuf/proto"
	_"github.com/golang/protobuf/ptypes/empty"
	_"google.golang.org/grpc"
	"go.uber.org/zap"
	"context"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Subprotocols: []string{"chat"},
}

func main() {
	logger.InitLogger()
	log := logger.GetLogger()

	// Подключение к базе данных
	dsn := "host=localhost user=postgres password=postgres dbname=forum port=5432 sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Проверка подключения
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database", zap.Error(err))
	}

	// Инициализация gRPC клиента для аутентификации
	authClient, err := client.NewAuthClient("localhost:50051")
	if err != nil {
		log.Fatal("Failed to create auth client", zap.Error(err))
	}

	// Инициализация репозиториев
	threadRepo := repository.NewThreadRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	chatRepo := repository.NewChatRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Инициализация сервисов
	postService := service.NewPostService(postRepo, commentRepo, threadRepo, userRepo)
	commentService := service.NewCommentService(commentRepo, userRepo)
	threadService := service.NewThreadService(threadRepo, postRepo, userRepo)
	chatService := service.NewChatService(chatRepo)

	// Инициализация обработчиков
	threadHandler := handlers.NewThreadHandler(threadService)
	postHandler := handlers.NewPostHandler(postService)
	commentHandler := handlers.NewCommentHandler(commentService)
	chatHandler := handlers.NewChatHandler(chatService)

	// Создание экземпляра Gin
	r := gin.Default()

	// Настройка CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8081"}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Загрузка HTML шаблонов
	r.LoadHTMLGlob("templates/*")
	// Настройка статических файлов
	r.Static("/static", "./static")

	// Инициализация middleware для аутентификации
	authMiddleware := middleware.AuthServiceMiddleware(authClient)

	// Инициализация Hub для веб-сокетов
	hub := handlers.NewHub(chatRepo)
	go hub.Run()

	// Группа защищенных маршрутов
	protected := r.Group("/api")
	protected.Use(authMiddleware)

	// Маршруты для тредов
	protected.GET("/threads", threadHandler.GetAllThreads)
	protected.GET("/threads/:id", threadHandler.GetThreadWithPosts)
	protected.GET("/threads/:id/posts", threadHandler.GetThreadPosts)
	protected.POST("/threads", threadHandler.CreateThread)
	protected.PUT("/threads/:id", threadHandler.UpdateThread)
	protected.DELETE("/threads/:id", threadHandler.DeleteThread)

	// Маршруты для постов
	protected.GET("/posts", postHandler.GetAllPosts)
	protected.GET("/posts/:id", postHandler.GetPost)
	protected.GET("/posts/:id/comments", postHandler.GetPostComments)
	protected.POST("/posts", postHandler.CreatePost)
	protected.PUT("/posts/:id", postHandler.UpdatePost)
	protected.DELETE("/posts/:id", postHandler.DeletePost)

	// Маршруты для комментариев
	protected.POST("/comments", commentHandler.CreateComment)
	protected.DELETE("/comments/:id", commentHandler.DeleteComment)

	// Маршруты для чата
	protected.POST("/chat", chatHandler.CreateMessage)
	protected.GET("/chat", chatHandler.GetMessages)

	// WebSocket маршрут (публичный)
	r.GET("/ws", authMiddleware, func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			log.Error("ID пользователя не найден в контексте")
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}

		username, exists := c.Get("username")
		if !exists {
			log.Error("Имя пользователя не найдено в контексте")
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}

		// Преобразуем userID в int
		userIDInt, ok := userID.(uint32)
		if !ok {
			log.Error("Неверный тип ID пользователя")
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}

		log.Info("WebSocket подключение", 
			zap.Int("user_id", int(userIDInt)),
			zap.String("username", username.(string)))

		// Создаем новый контекст с данными пользователя
		ctx := context.WithValue(c.Request.Context(), "user_id", int(userIDInt))
		ctx = context.WithValue(ctx, "username", username.(string))
		
		// Создаем новый запрос с обновленным контекстом
		newReq := c.Request.WithContext(ctx)
		
		hub.HandleWebSocket(c.Writer, newReq)
	})

	// Главная страница со списком тредов
	r.GET("/", authMiddleware, func(c *gin.Context) {
		user, _ := c.Get("user")
		userID, _ := c.Get("user_id")
		userRole, _ := c.Get("user_role")
		username, _ := c.Get("username")
		if userRole == nil {
			userRole = "user"
		}
		log.Info("Debug - User Role in /", zap.Any("role", userRole))
		log.Info("Debug - Raw user role in /", zap.String("role", fmt.Sprintf("%v", userRole)))

		// Преобразуем userID в int для корректного сравнения
		var userIDInt int
		if userID != nil {
			userIDInt = int(userID.(uint32))
		}

		log.Info("Debug - User ID in /", zap.Int("user_id", userIDInt))
		log.Info("Debug - User Role in /", zap.Any("role", userRole))

		threads, err := threadService.GetAllThreads()
		if err != nil {
			c.HTML(500, "error.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		c.HTML(200, "index.html", gin.H{
			"title":     "Главная страница",
			"user":      user,
			"user_id":   userIDInt,
			"user_role": userRole,
			"username": username,
			"Threads":   threads,
		})
	})

	// Получение всех тредов (HTML)
	r.GET("/threads", func(c *gin.Context) {
		user, _ := c.Get("user")
		userID, _ := c.Get("user_id")
		userRole, _ := c.Get("user_role")
		username, _ := c.Get("username")
		if userRole == nil {
			userRole = "user"
		}
		log.Info("Debug - User Role in /threads", zap.Any("role", userRole))
		log.Info("Debug - Raw user role in /threads", zap.String("role", fmt.Sprintf("%q", userRole)))

		threads, err := threadService.GetAllThreads()
		if err != nil {
			c.HTML(500, "error.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		// Преобразуем userID в int для корректного сравнения
		var userIDInt int
		if userID != nil {
			userIDInt = int(userID.(uint32))
		}

		log.Info("Debug - User ID in /threads", zap.Int("user_id", userIDInt))
		log.Info("Debug - User Role in /threads", zap.Any("role", userRole))

		c.HTML(200, "threads.html", gin.H{
			"threads":   threads,
			"user":      user,
			"user_id":   userIDInt,
			"user_role": userRole,
			"username": username,
		})
	})

	// Получение конкретного треда с постами (HTML)
	r.GET("/threads/:id", authMiddleware, func(c *gin.Context) {
		user, _ := c.Get("user")
		userID, _ := c.Get("user_id")
		userRole, _ := c.Get("user_role")
		username, _ := c.Get("username")
		if userRole == nil {
			userRole = "user"
		}
		log.Info("Debug - User Role in /threads/:id", zap.Any("role", userRole))
		log.Info("Debug - Raw user role in /threads/:id", zap.String("role", fmt.Sprintf("%q", userRole)))

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

		// Преобразуем userID в int для корректного сравнения
		var userIDInt int
		if userID != nil {
			userIDInt = int(userID.(uint32))
		}

		log.Info("Debug - Thread Author ID", zap.Int("author_id", thread.AuthorID))
		log.Info("Debug - User ID", zap.Int("user_id", userIDInt))
		log.Info("Debug - User Role", zap.Any("role", userRole))

		c.HTML(200, "thread.html", gin.H{
			"Thread":    thread,
			"posts":     posts,
			"user":      user,
			"user_id":   userIDInt,
			"user_role": userRole,
			"username": username,
		})
	})

	// Получение конкретного поста (HTML)
	r.GET("/posts/:id", authMiddleware, func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.HTML(400, "error.html", gin.H{
				"error": "Неверный ID поста",
			})
			return
		}

		post, comments, err := postService.GetPostWithComments(id)
		if err != nil {
			log.Error("Ошибка при получении поста с комментариями", zap.Error(err))
			c.HTML(404, "error.html", gin.H{
				"error": "Пост не найден",
			})
			return
		}

		// Получаем информацию о пользователе из контекста
		user, _ := c.Get("user")
		userID, _ := c.Get("user_id")
		userRole, _ := c.Get("user_role")
		username, _ := c.Get("username")

		// Преобразуем userID в int для корректного сравнения
		var userIDInt int
		if userID != nil {
			userIDInt = int(userID.(uint32))
		}

		// Проверяем, может ли пользователь редактировать пост
		post.CanEdit = userIDInt == post.AuthorID || userRole == "admin"

		// Добавляем флаг CanDelete для комментариев
		for i := range comments {
			comments[i].CanDelete = comments[i].AuthorID == userIDInt || userRole == "admin"
		}

		c.HTML(200, "post.html", gin.H{
			"post":      post,
			"comments":  comments,
			"user":      user,
			"user_id":   userIDInt,
			"user_role": userRole,
			"username":  username,
		})
	})

	// Запуск сервера
	port := 8081
	log.Info("Server is running", zap.Int("port", port))
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
