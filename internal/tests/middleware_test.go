package tests

import (
	"ForumService/internal/middleware"
	"ForumService/internal/models"
	"ForumService/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockUserService - мок для UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserByID(id int) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUserPosts(userID int) ([]*models.Post, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Post), args.Error(1)
}

func (m *MockUserService) GetUserCommentCount(userID int) (int, error) {
	args := m.Called(userID)
	return args.Int(0), args.Error(1)
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		cookieValue    string
		mockUser       *models.User
		mockError      error
		expectedStatus int
		expectedUser   bool
	}{
		{
			name:           "No cookie",
			cookieValue:    "",
			expectedStatus: http.StatusOK,
			expectedUser:   false,
		},
		{
			name:           "Invalid cookie format",
			cookieValue:    "invalid",
			expectedStatus: http.StatusOK,
			expectedUser:   false,
		},
		{
			name:           "User not found",
			cookieValue:    "123",
			mockError:      service.ErrUserNotFound,
			expectedStatus: http.StatusOK,
			expectedUser:   false,
		},
		{
			name:        "Valid user",
			cookieValue: "123",
			mockUser: &models.User{
				ID:       123,
				Username: "testuser",
				Email:    "test@example.com",
			},
			expectedStatus: http.StatusOK,
			expectedUser:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			if tt.cookieValue == "123" {
				mockService.On("GetUserByID", 123).Return(tt.mockUser, tt.mockError)
			}

			router := gin.New()
			router.Use(middleware.AuthMiddleware(mockService))
			router.GET("/test", func(c *gin.Context) {
				user, exists := c.Get("user")
				if tt.expectedUser {
					assert.True(t, exists)
					assert.NotNil(t, user)
					u := user.(*models.User)
					assert.Equal(t, tt.mockUser.ID, u.ID)
					assert.Equal(t, tt.mockUser.Username, u.Username)
					assert.Equal(t, tt.mockUser.Email, u.Email)
				} else {
					assert.False(t, exists)
				}
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			if tt.cookieValue != "" {
				req.AddCookie(&http.Cookie{
					Name:  "session",
					Value: tt.cookieValue,
				})
			}
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestRequireAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setUserID      bool
		expectedStatus int
		expectedRedirect bool
	}{
		{
			name:           "No user ID",
			setUserID:      false,
			expectedStatus: http.StatusFound,
			expectedRedirect: true,
		},
		{
			name:           "Has user ID",
			setUserID:      true,
			expectedStatus: http.StatusOK,
			expectedRedirect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(func(c *gin.Context) {
				if tt.setUserID {
					c.Set("userID", 123)
				}
				c.Next()
			})
			router.Use(middleware.RequireAuth())
			router.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedRedirect {
				assert.Equal(t, "/login", w.Header().Get("Location"))
			}
		})
	}
} 