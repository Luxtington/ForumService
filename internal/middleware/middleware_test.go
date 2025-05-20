package middleware

import (
	"ForumService/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"ForumService/internal/client"
	"context"
	"AuthService/proto"
	"google.golang.org/grpc"
)

// Мок для UserService
type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) GetUserByID(id int) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockUserService) GetUserPosts(userID int) ([]*models.Post, error) {
	return nil, nil
}

func (m *mockUserService) GetUserCommentCount(userID int) (int, error) {
	return 0, nil
}

// Мок для AuthClient
type mockAuthClient struct {
	mock.Mock
}

func (m *mockAuthClient) ValidateToken(token string) (uint32, string, string, error) {
	args := m.Called(token)
	return args.Get(0).(uint32), args.String(1), args.String(2), args.Error(3)
}

func (m *mockAuthClient) Register(username, password string) (uint32, string, string, error) {
	args := m.Called(username, password)
	return args.Get(0).(uint32), args.String(1), args.String(2), args.Error(3)
}

func (m *mockAuthClient) Login(username, password string) (uint32, string, string, error) {
	args := m.Called(username, password)
	return args.Get(0).(uint32), args.String(1), args.String(2), args.Error(3)
}

// Мок для proto.AuthServiceClient
type mockProtoAuthClient struct {
	proto.UnimplementedAuthServiceServer
}

func (m *mockProtoAuthClient) ValidateToken(ctx context.Context, req *proto.ValidateTokenRequest, opts ...grpc.CallOption) (*proto.ValidateTokenResponse, error) {
	if req.Token == "valid_token" {
		return &proto.ValidateTokenResponse{UserId: 1, Username: "test_user", Role: "user"}, nil
	}
	return nil, assert.AnError
}

func (m *mockProtoAuthClient) Login(ctx context.Context, req *proto.LoginRequest, opts ...grpc.CallOption) (*proto.LoginResponse, error) {
	return nil, nil
}

func (m *mockProtoAuthClient) Register(ctx context.Context, req *proto.RegisterRequest, opts ...grpc.CallOption) (*proto.RegisterResponse, error) {
	return nil, nil
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userService := new(mockUserService)
	user := &models.User{ID: 1, Username: "test_user"}
	userService.On("GetUserByID", 1).Return(user, nil)

	router := gin.New()
	router.Use(AuthMiddleware(userService))
	router.GET("/test", func(c *gin.Context) {
		u, _ := c.Get("user")
		c.JSON(200, gin.H{"user": u})
	})

	// Тест с валидной кукой
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "1"})
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// Тест без куки
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// Тест с невалидной кукой
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "invalid"})
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestRequireAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequireAuth())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Тест без пользователя в контексте
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 302, w.Code) // Редирект на /login

	// Тест с пользователем в контексте
	router = gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.Next()
	})
	router.Use(RequireAuth())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestDummyMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(DummyMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestAuthServiceMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ac := &client.AuthClient{Client: &mockProtoAuthClient{}}

	router := gin.New()
	router.Use(AuthServiceMiddleware(ac))
	router.GET("/test", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")
		role, _ := c.Get("user_role")
		c.JSON(200, gin.H{"user_id": userID, "username": username, "role": role})
	})

	// Тест с валидным токеном в заголовке
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer valid_token")
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// Тест с валидным токеном в куке
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: "valid_token"})
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// Тест с невалидным токеном
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)

	// Тест без токена
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
} 