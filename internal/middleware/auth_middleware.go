package middleware

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

// AuthServiceMiddleware проверяет JWT токен через AuthService
func AuthServiceMiddleware(authServiceURL string) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "требуется аутентификация"})
            c.Abort()
            return
        }

        // Проверяем токен через AuthService
        req, err := http.NewRequest("GET", authServiceURL+"/api/auth/validate", nil)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при проверке аутентификации"})
            c.Abort()
            return
        }
        req.Header.Set("Authorization", token)

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

        c.Next()
    }
} 