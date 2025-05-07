package middleware

import (
	"ForumService/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func AuthMiddleware(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := c.Cookie("session")
		if err != nil {
			c.Next()
			return
		}

		userID, err := strconv.Atoi(session)
		if err != nil {
			c.Next()
			return
		}

		user, err := userService.GetUserByID(userID)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user", user)
		c.Set("userID", userID)
		c.Next()
	}
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := c.Get("userID"); !exists {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}
