package handlers

import (
	"ForumService/internal/middleware"
	"ForumService/internal/service"
	"github.com/gin-gonic/gin"
)

type Services struct {
	UserService    service.UserService
	PostService    service.PostService
	CommentService service.CommentService
	ThreadService  service.ThreadService
	ChatService    service.ChatService
}

func RegisterRoutes(router *gin.Engine, services *Services) {
	// Добавляем middleware для обработки ошибок
	router.Use(middleware.ErrorHandler())

	// Инициализация обработчиков
	viewsHandler := NewViewsHandler(services.ThreadService, services.PostService, services.CommentService, services.ChatService)
	threadHandler := NewThreadHandler(services.ThreadService)
	postHandler := NewPostHandler(services.PostService)
	commentHandler := NewCommentHandler(services.CommentService)
	chatHandler := NewChatHandler(services.ChatService)

	// Главная страница
	router.GET("/", viewsHandler.Index)

	// Маршруты для отображения страниц
	router.GET("/threads/:id", viewsHandler.ShowThread)
	router.GET("/posts/:id", viewsHandler.ShowPost)

	// API маршруты
	api := router.Group("/api")
	{
		// Маршруты для тредов
		threads := api.Group("/threads")
		{
			threads.GET("", threadHandler.GetAllThreads)
			threads.GET("/:id", threadHandler.GetThreadWithPosts)
			threads.GET("/:id/posts", threadHandler.GetThreadPosts)
			threads.POST("", threadHandler.CreateThread)
			threads.PUT("/:id", threadHandler.UpdateThread)
			threads.DELETE("/:id", threadHandler.DeleteThread)
		}

		// Маршруты для постов
		posts := api.Group("/posts")
		{
			posts.GET("", postHandler.GetAllPosts)
			posts.GET("/:id", postHandler.GetPost)
			posts.GET("/:id/comments", postHandler.GetPostComments)
			posts.POST("", postHandler.CreatePost)
			posts.PUT("/:id", postHandler.UpdatePost)
			posts.DELETE("/:id", postHandler.DeletePost)
		}

		// Маршруты для комментариев
		comments := api.Group("/comments")
		{
			comments.POST("", commentHandler.CreateComment)
			comments.DELETE("/:id", commentHandler.DeleteComment)
		}

		// Чат
		api.GET("/chat", chatHandler.GetMessages)
		api.POST("/chat", chatHandler.CreateMessage)
	}

	// Обработчики ошибок
	router.NoRoute(func(c *gin.Context) {
		c.HTML(404, "error.html", gin.H{
			"Title":   "Страница не найдена",
			"Message": "Запрашиваемая страница не существует",
		})
	})
}
