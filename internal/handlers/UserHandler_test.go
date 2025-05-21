package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupUserTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.LoadHTMLGlob("../../templates/*")
	return router
}

func TestUserHandler_ShowLoginForm(t *testing.T) {
	handler := NewUserHandler()
	router := setupUserTestRouter()
	router.GET("/login", handler.ShowLoginForm)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/", w.Header().Get("Location"))
}

func TestUserHandler_ShowRegisterForm(t *testing.T) {
	handler := NewUserHandler()
	router := setupUserTestRouter()
	router.GET("/register", handler.ShowRegisterForm)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/register", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/", w.Header().Get("Location"))
}

func TestUserHandler_Register(t *testing.T) {
	handler := NewUserHandler()
	router := setupUserTestRouter()
	router.POST("/register", handler.Register)

	requestBody := map[string]interface{}{
		"username": "testuser",
		"password": "testpass",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/", w.Header().Get("Location"))
}

func TestUserHandler_Login(t *testing.T) {
	handler := NewUserHandler()
	router := setupUserTestRouter()
	router.POST("/login", handler.Login)

	requestBody := map[string]interface{}{
		"username": "testuser",
		"password": "testpass",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/", w.Header().Get("Location"))
}

func TestUserHandler_Logout(t *testing.T) {
	handler := NewUserHandler()
	router := setupUserTestRouter()
	router.POST("/logout", handler.Logout)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/logout", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/", w.Header().Get("Location"))
}

func TestUserHandler_ShowProfile(t *testing.T) {
	handler := NewUserHandler()
	router := setupUserTestRouter()
	router.GET("/profile", handler.ShowProfile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/profile", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/", w.Header().Get("Location"))
} 

func TestUserHandler_Register_WithoutJSONContentType(t *testing.T) {
    handler := NewUserHandler()
    router := setupUserTestRouter()
    router.POST("/register", handler.Register)

    requestBody := map[string]interface{}{
        "username": "testuser",
        "password": "testpass",
    }
    jsonBody, _ := json.Marshal(requestBody)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
    // Не устанавливаем Content-Type

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusFound, w.Code)
    assert.Equal(t, "/", w.Header().Get("Location"))
}