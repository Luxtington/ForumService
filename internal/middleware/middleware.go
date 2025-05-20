package middleware

import (
	"github.com/gin-gonic/gin"
)

func DummyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
} 