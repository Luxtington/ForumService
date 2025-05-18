package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// AuthServiceMiddleware проверяет JWT токен через AuthService
func AuthServiceMiddleware(authServiceURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "токен не предоставлен"})
			c.Abort()
			return
		}

		// Убираем префикс "Bearer " если он есть
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Создаем запрос к AuthService
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/auth/validate", authServiceURL), nil)
		if err != nil {
			c.JSON(500, gin.H{"error": "ошибка при создании запроса"})
			c.Abort()
			return
		}

		// Добавляем токен в заголовок
		req.Header.Set("Authorization", "Bearer "+token)

		// Выполняем запрос
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(500, gin.H{"error": "ошибка при проверке токена"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(401, gin.H{"error": "недействительный токен"})
			c.Abort()
			return
		}

		// Декодируем ответ
		var user struct {
			ID       uint   `json:"id"`
			Username string `json:"username"`
			Role     string `json:"role"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			c.JSON(500, gin.H{"error": "ошибка при обработке ответа"})
			c.Abort()
			return
		}

		// Сохраняем информацию о пользователе в контексте
		c.Set("user_id", user.ID)
		c.Set("username", user.Username)
		c.Set("user_role", user.Role)
		
		// Отладочный вывод
		fmt.Printf("Debug - User ID: %d, Username: %s, Role: %s\n", user.ID, user.Username, user.Role)
		fmt.Printf("Debug - User Role type: %T\n", user.Role)
		fmt.Printf("Debug - Raw user role: %q\n", user.Role)
		
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