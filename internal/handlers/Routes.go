package handlers

import (
	"ForumService/internal/service"
	"github.com/gin-gonic/gin"
)

type Services struct {
	UserService    service.UserService
	PostService    service.PostService
	CommentService service.CommentService
	ThreadService  service.ThreadService
}

func RegisterRoutes(router *gin.Engine, services *Services) {
	// Инициализация обработчиков
	threadHandler := NewThreadHandler(services.ThreadService)
	viewsHandler := NewViewsHandler(services.ThreadService, services.CommentService)

	// Главная страница
	router.GET("/", viewsHandler.Index)

	// API маршруты
	api := router.Group("/api")
	{
		threads := api.Group("/threads")
		{
			threads.POST("", threadHandler.CreateThread)
			threads.GET("/:id", threadHandler.GetThreadWithPosts)
			threads.DELETE("/:id", threadHandler.DeleteThread)
		}
	}

	router.GET("/threads/:id", viewsHandler.ShowThread)

	// Обработчики ошибок
	router.NoRoute(func(c *gin.Context) {
		c.HTML(404, "error.html", gin.H{
			"Title":   "Страница не найдена",
			"Message": "Запрашиваемая страница не существует",
		})
	})
}
