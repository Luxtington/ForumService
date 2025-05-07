package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// Временные функции для тестирования
func (h *UserHandler) ShowLoginForm(c *gin.Context) {
	c.Redirect(http.StatusFound, "/")
}

func (h *UserHandler) ShowRegisterForm(c *gin.Context) {
	c.Redirect(http.StatusFound, "/")
}

func (h *UserHandler) Register(c *gin.Context) {
	c.Redirect(http.StatusFound, "/")
}

func (h *UserHandler) Login(c *gin.Context) {
	c.Redirect(http.StatusFound, "/")
}

func (h *UserHandler) Logout(c *gin.Context) {
	c.Redirect(http.StatusFound, "/")
}

func (h *UserHandler) ShowProfile(c *gin.Context) {
	c.Redirect(http.StatusFound, "/")
} 