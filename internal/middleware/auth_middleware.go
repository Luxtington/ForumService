package middleware

import (
	"ForumService/internal/client"
	"github.com/gin-gonic/gin"
	"strings"
	"fmt"
)

// AuthServiceMiddleware проверяет JWT токен через AuthService
func AuthServiceMiddleware(authClient *client.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка Authorization или из куки
		authHeader := c.GetHeader("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Если токен не найден в заголовке, пробуем получить из куки
		if token == "" {
			if cookieToken, err := c.Cookie("auth_token"); err == nil {
				token = cookieToken
			}
		}

		if token == "" {
			c.JSON(401, gin.H{"error": "токен не предоставлен"})
			c.Abort()
			return
		}

		// Проверяем токен через gRPC клиент
		userID, username, role, err := authClient.ValidateToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "недействительный токен"})
			c.Abort()
			return
		}

		// Сохраняем информацию о пользователе в контексте
		c.Set("user_id", userID)
		c.Set("username", username)
		c.Set("user_role", role)
		
		// Отладочный вывод
		fmt.Printf("Debug - User ID: %d, Username: %s, Role: %s\n", userID, username, role)
		fmt.Printf("Debug - User Role type: %T\n", role)
		fmt.Printf("Debug - Raw user role: %q\n", role)
		
		// Проверяем, что роль установлена в контексте
		if role, exists := c.Get("user_role"); exists {
			fmt.Printf("Debug - Role in context: %v (type: %T)\n", role, role)
			fmt.Printf("Debug - Raw role in context: %q\n", role)
			
			// Проверяем, что роль является строкой
			if roleStr, ok := role.(string); ok {
				fmt.Printf("Debug - Role is string: %q\n", roleStr)
				// Устанавливаем роль заново, чтобы убедиться, что это строка
				c.Set("user_role", roleStr)
			} else {
				fmt.Printf("Debug - Role is not string: %T\n", role)
				// Если роль не строка, устанавливаем значение по умолчанию
				c.Set("user_role", "user")
			}
		} else {
			fmt.Printf("Debug - Role not found in context\n")
			// Если роль не найдена, устанавливаем значение по умолчанию
			c.Set("user_role", "user")
		}

		// Проверяем, что ID установлен в контексте
		if id, exists := c.Get("user_id"); exists {
			fmt.Printf("Debug - User ID in context: %v (type: %T)\n", id, id)
		} else {
			fmt.Printf("Debug - User ID not found in context\n")
		}

		c.Next()
	}
} 