package middleware

import (
	"ForumService/internal/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ErrorHandler middleware для обработки ошибок
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Проверяем, есть ли ошибки
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			
			// Проверяем, является ли ошибка кастомной ошибкой форума
			if forumErr, ok := err.(*errors.ForumError); ok {
				c.JSON(forumErr.Code, gin.H{
					"error": forumErr.Message,
				})
				return
			}

			// Если это не кастомная ошибка, возвращаем 500
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Внутренняя ошибка сервера",
			})
		}
	}
} 