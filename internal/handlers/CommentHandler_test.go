package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"ForumService/internal/handlers/mocks"
	"ForumService/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupCommentTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.LoadHTMLGlob("../../templates/*")
	return router
}

func TestCommentHandler_CreateComment_Success(t *testing.T) {
	mockCommentService := &mocks.MockCommentService{
		CreateCommentFunc: func(postID int, authorID int, content string) (*models.Comment, error) {
			return &models.Comment{
				ID:       1,
				PostID:   postID,
				AuthorID: authorID,
				Content:  content,
			}, nil
		},
	}

	handler := NewCommentHandler(mockCommentService)

	router := setupCommentTestRouter()
	router.POST("/comments", func(c *gin.Context) {
		c.Set("user_id", uint32(1))
		handler.CreateComment(c)
	})

	requestBody := map[string]interface{}{
		"post_id": 1,
		"content": "Test comment",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/comments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCommentHandler_CreateComment_Unauthorized(t *testing.T) {
	mockCommentService := &mocks.MockCommentService{}
	handler := NewCommentHandler(mockCommentService)

	router := setupCommentTestRouter()
	router.POST("/comments", handler.CreateComment)

	requestBody := map[string]interface{}{
		"post_id": 1,
		"content": "Test comment",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/comments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCommentHandler_CreateComment_InvalidData(t *testing.T) {
	mockCommentService := &mocks.MockCommentService{}
	handler := NewCommentHandler(mockCommentService)

	router := setupCommentTestRouter()
	router.POST("/comments", func(c *gin.Context) {
		c.Set("user_id", uint32(1))
		handler.CreateComment(c)
	})

	requestBody := map[string]interface{}{
		"post_id": "invalid",
		"content": "",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/comments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCommentHandler_CreateComment_ServiceError(t *testing.T) {
	mockCommentService := &mocks.MockCommentService{
		CreateCommentFunc: func(postID int, authorID int, content string) (*models.Comment, error) {
			return nil, errors.New("service error")
		},
	}

	handler := NewCommentHandler(mockCommentService)

	router := setupCommentTestRouter()
	router.POST("/comments", func(c *gin.Context) {
		c.Set("user_id", uint32(1))
		handler.CreateComment(c)
	})

	requestBody := map[string]interface{}{
		"post_id": 1,
		"content": "Test comment",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/comments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCommentHandler_DeleteComment_Success(t *testing.T) {
	mockCommentService := &mocks.MockCommentService{
		GetCommentByIDFunc: func(id int) (*models.Comment, error) {
			return &models.Comment{
				ID:       1,
				PostID:   1,
				AuthorID: 1,
				Content:  "Test comment",
			}, nil
		},
		DeleteCommentFunc: func(id int, userID int) error {
			return nil
		},
	}

	handler := NewCommentHandler(mockCommentService)

	router := setupCommentTestRouter()
	router.DELETE("/comments/:id", func(c *gin.Context) {
		c.Set("user_id", uint32(1))
		c.Set("user_role", "user")
		handler.DeleteComment(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/comments/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestCommentHandler_DeleteComment_Unauthorized(t *testing.T) {
	mockCommentService := &mocks.MockCommentService{}
	handler := NewCommentHandler(mockCommentService)

	router := setupCommentTestRouter()
	router.DELETE("/comments/:id", handler.DeleteComment)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/comments/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCommentHandler_DeleteComment_InvalidID(t *testing.T) {
	mockCommentService := &mocks.MockCommentService{}
	handler := NewCommentHandler(mockCommentService)

	router := setupCommentTestRouter()
	router.DELETE("/comments/:id", func(c *gin.Context) {
		c.Set("user_id", uint32(1))
		c.Set("user_role", "user")
		handler.DeleteComment(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/comments/invalid", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCommentHandler_DeleteComment_NoPermission(t *testing.T) {
	mockCommentService := &mocks.MockCommentService{
		GetCommentByIDFunc: func(id int) (*models.Comment, error) {
			return &models.Comment{
				ID:       1,
				PostID:   1,
				AuthorID: 2,
				Content:  "Test comment",
			}, nil
		},
	}

	handler := NewCommentHandler(mockCommentService)

	router := setupCommentTestRouter()
	router.DELETE("/comments/:id", func(c *gin.Context) {
		c.Set("user_id", uint32(1))
		c.Set("user_role", "user")
		handler.DeleteComment(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/comments/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCommentHandler_DeleteComment_AdminSuccess(t *testing.T) {
	mockCommentService := &mocks.MockCommentService{
		GetCommentByIDFunc: func(id int) (*models.Comment, error) {
			return &models.Comment{
				ID:       1,
				PostID:   1,
				AuthorID: 2,
				Content:  "Test comment",
			}, nil
		},
		DeleteCommentFunc: func(id int, userID int) error {
			return nil
		},
	}

	handler := NewCommentHandler(mockCommentService)

	router := setupCommentTestRouter()
	router.DELETE("/comments/:id", func(c *gin.Context) {
		c.Set("user_id", uint32(1))
		c.Set("user_role", "admin")
		handler.DeleteComment(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/comments/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestCommentHandler_DeleteComment_GetCommentError(t *testing.T) {
	mockCommentService := &mocks.MockCommentService{
		GetCommentByIDFunc: func(id int) (*models.Comment, error) {
			return nil, errors.New("error getting comment")
		},
	}

	handler := NewCommentHandler(mockCommentService)

	router := setupCommentTestRouter()
	router.DELETE("/comments/:id", func(c *gin.Context) {
		c.Set("user_id", uint32(1))
		c.Set("user_role", "user")
		handler.DeleteComment(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/comments/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCommentHandler_DeleteComment_DeleteError(t *testing.T) {
	mockCommentService := &mocks.MockCommentService{
		GetCommentByIDFunc: func(id int) (*models.Comment, error) {
			return &models.Comment{
				ID:       1,
				PostID:   1,
				AuthorID: 1,
				Content:  "Test comment",
			}, nil
		},
		DeleteCommentFunc: func(id int, userID int) error {
			return errors.New("error deleting comment")
		},
	}

	handler := NewCommentHandler(mockCommentService)

	router := setupCommentTestRouter()
	router.DELETE("/comments/:id", func(c *gin.Context) {
		c.Set("user_id", uint32(1))
		c.Set("user_role", "user")
		handler.DeleteComment(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/comments/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCommentHandler_CreateChatMessage_Success(t *testing.T) {
	mockCommentService := &mocks.MockCommentService{
		CreateCommentFunc: func(postID int, authorID int, content string) (*models.Comment, error) {
			return &models.Comment{
				ID:       1,
				PostID:   0,
				AuthorID: 1,
				Content:  content,
			}, nil
		},
	}

	handler := NewCommentHandler(mockCommentService)

	router := setupCommentTestRouter()
	router.POST("/chat/messages", handler.CreateChatMessage)

	requestBody := map[string]interface{}{
		"content": "Test chat message",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/chat/messages", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCommentHandler_CreateChatMessage_InvalidData(t *testing.T) {
	mockCommentService := &mocks.MockCommentService{}
	handler := NewCommentHandler(mockCommentService)

	router := setupCommentTestRouter()
	router.POST("/chat/messages", handler.CreateChatMessage)

	requestBody := map[string]interface{}{
		"content": "",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/chat/messages", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCommentHandler_CreateChatMessage_ServiceError(t *testing.T) {
	mockCommentService := &mocks.MockCommentService{
		CreateCommentFunc: func(postID int, authorID int, content string) (*models.Comment, error) {
			return nil, errors.New("service error")
		},
	}

	handler := NewCommentHandler(mockCommentService)

	router := setupCommentTestRouter()
	router.POST("/chat/messages", handler.CreateChatMessage)

	requestBody := map[string]interface{}{
		"content": "Test chat message",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/chat/messages", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}