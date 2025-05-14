package middleware

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "encoding/json"
    "strings"
)

// AuthServiceMiddleware проверяет JWT токен через AuthService
func AuthServiceMiddleware(authServiceURL string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Сначала проверяем заголовок Authorization
        authHeader := c.GetHeader("Authorization")
        var token string
        if authHeader != "" {
            // Убираем префикс "Bearer " если он есть
            token = strings.TrimPrefix(authHeader, "Bearer ")
        } else {
            // Если нет в заголовке, проверяем куки
            var err error
            token, err = c.Cookie("auth_token")
            if err != nil {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "требуется аутентификация"})
                c.Abort()
                return
            }
        }

        // Проверяем токен через AuthService
        req, err := http.NewRequest("GET", authServiceURL+"/api/auth/validate", nil)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при проверке аутентификации"})
            c.Abort()
            return
        }

        // Добавляем куки в запрос
        req.AddCookie(&http.Cookie{
            Name:  "auth_token",
            Value: token,
        })

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при проверке аутентификации"})
            c.Abort()
            return
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "недействительный токен"})
            c.Abort()
            return
        }

        // Получаем данные пользователя из ответа
        var user struct {
            ID int `json:"id"`
        }
        if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при получении данных пользователя"})
            c.Abort()
            return
        }

        // Добавляем ID пользователя в контекст
        c.Set("user_id", user.ID)
        c.Next()
    }
} 