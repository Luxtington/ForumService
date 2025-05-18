package main

import (
	"ForumService/internal/client"
	"ForumService/internal/handlers"
	"ForumService/internal/middleware"
	"ForumService/internal/repository"
	"ForumService/internal/service"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"time"
	"github.com/gorilla/websocket"
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

	// Инициализация gRPC клиента для аутентификации
	authClient, err := client.NewAuthClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to create auth client: %v", err)
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
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8082"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Cookie", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
	r.GET("/ws", func(c *gin.Context) {
		log.Printf("Получен запрос на WebSocket соединение")
		
		// Проверяем токен в куках
		token, err := c.Cookie("auth_token")
		if err != nil {
			log.Printf("Ошибка получения токена из куки: %v", err)
			return
		}

		if token == "" {
			log.Printf("Пустой токен в куки")
			return
		}

		log.Printf("Токен получен успешно")
		
		// Добавляем токен в заголовок
		c.Request.Header.Set("Authorization", "Bearer "+token)
		
		// Вызываем middleware для установки контекста
		authMiddleware(c)
		
		if c.IsAborted() {
			log.Printf("Middleware прервал запрос")
			return
		}

		// Проверяем, что пользователь аутентифицирован
		userID, exists := c.Get("user_id")
		if !exists || userID == nil {
			log.Printf("Пользователь не аутентифицирован")
			return
		}

		username, _ := c.Get("username")
		log.Printf("WebSocket соединение устанавливается для пользователя ID: %v, username: %v", userID, username)

		// Обновляем соединение до WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Ошибка при обновлении соединения до WebSocket: %v", err)
			return
		}

		// Создаем нового клиента
		client := &handlers.Client{
			Conn:     conn,
			Send:     make(chan []byte, 256),
			Username: username.(string),
			UserID:   int(userID.(uint32)),
		}

		hub.Register <- client

		go hub.WritePump(client)
		go hub.ReadPump(client)
	})

	// Главная страница со списком тредов
	r.GET("/", func(c *gin.Context) {
		user, _ := c.Get("user")
		userID, _ := c.Get("user_id")
		userRole, _ := c.Get("user_role")
		username, _ := c.Get("username")
		if userRole == nil {
			userRole = "user"
		}
		fmt.Printf("Debug - User Role in /: %v (type: %T)\n", userRole, userRole)
		fmt.Printf("Debug - Raw user role in /: %q\n", userRole)

		// Преобразуем userID в int для корректного сравнения
		var userIDInt int
		if userID != nil {
			userIDInt = int(userID.(uint))
		}

		fmt.Printf("Debug - User ID in /: %d\n", userIDInt)
		fmt.Printf("Debug - User Role in /: %v\n", userRole)

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
		fmt.Printf("Debug - User Role in /threads: %v (type: %T)\n", userRole, userRole)
		fmt.Printf("Debug - Raw user role in /threads: %q\n", userRole)

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
			userIDInt = int(userID.(uint))
		}

		fmt.Printf("Debug - User ID in /threads: %d\n", userIDInt)
		fmt.Printf("Debug - User Role in /threads: %v\n", userRole)

		c.HTML(200, "threads.html", gin.H{
			"threads":   threads,
			"user":      user,
			"user_id":   userIDInt,
			"user_role": userRole,
			"username": username,
		})
	})

	// Получение конкретного треда с постами (HTML)
	r.GET("/threads/:id", func(c *gin.Context) {
		user, _ := c.Get("user")
		userID, _ := c.Get("user_id")
		userRole, _ := c.Get("user_role")
		username, _ := c.Get("username")
		if userRole == nil {
			userRole = "user"
		}
		fmt.Printf("Debug - User Role in /threads/:id: %v (type: %T)\n", userRole, userRole)
		fmt.Printf("Debug - Raw user role in /threads/:id: %q\n", userRole)

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
			userIDInt = int(userID.(uint))
		}

		fmt.Printf("Debug - Thread Author ID: %d\n", thread.AuthorID)
		fmt.Printf("Debug - User ID: %d\n", userIDInt)
		fmt.Printf("Debug - User Role: %v\n", userRole)

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
	r.GET("/posts/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.HTML(400, "error.html", gin.H{
				"error": "Неверный ID поста",
			})
			return
		}

		post, comments, err := postService.GetPostWithComments(id)
		if err != nil {
			log.Printf("Ошибка при получении поста с комментариями: %v", err)
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

		log.Printf("Debug - User info: user=%+v, userID=%+v", user, userID)

		// Проверяем, может ли пользователь редактировать пост
		if userID != nil {
			userIDInt := int(userID.(uint))
			if post.AuthorID == userIDInt {
				post.CanEdit = true
			}
		}

		// Добавляем флаг CanDelete для комментариев
		for i := range comments {
			if userID != nil {
				userIDInt := int(userID.(uint))
				comments[i].CanDelete = comments[i].AuthorID == userIDInt
			}
		}

		log.Printf("Отправка данных в шаблон: post=%+v, comments=%+v, user=%+v, userID=%+v", post, comments, user, userID)

		c.HTML(200, "post.html", gin.H{
			"post":     post,
			"comments": comments,
			"user":     user,
			"user_id":  userID,
			"user_role": userRole,
			"username": username,
			"CanEdit":  post.CanEdit,
		})
	})

	// Запуск сервера
	port := 8081
	log.Printf("Server is running on port %d", port)
	if err := r.Run(":" + strconv.Itoa(port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
