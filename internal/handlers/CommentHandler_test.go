package handlers

import (
	"bytes"
	"encoding/json"
	"ForumService/internal/handlers/mocks"
	"ForumService/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateComment(t *testing.T) {
	tests := []struct {
		name           string
		postID         string
		requestBody    map[string]interface{}
		mockResponse   *models.Comment
		mockError      error
		expectedStatus int
		userID         uint
	}{
		{
			name:   "успешное создание комментария",
			postID: "1",
			requestBody: map[string]interface{}{
				"post_id": 1,
				"content": "Test Comment Content",
			},
			mockResponse: &models.Comment{
				ID:       1,
				Content:  "Test Comment Content",
				PostID:   1,
				AuthorID: 1,
			},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			userID:         1,
		},
		{
			name:   "неверный формат данных",
			postID: "1",
			requestBody: map[string]interface{}{
				"invalid": "data",
			},
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			userID:         1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.MockCommentService{
				CreateCommentFunc: func(postID int, authorID int, content string) (*models.Comment, error) {
					return tt.mockResponse, tt.mockError
				},
			}

			handler := NewCommentHandler(mockService)
			router := setupTestRouter()
			router.POST("/posts/:postID/comments", func(c *gin.Context) {
				c.Set("user_id", tt.userID)
				c.Params = []gin.Param{{Key: "postID", Value: tt.postID}}
				handler.CreateComment(c)
			})

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/posts/"+tt.postID+"/comments", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestDeleteComment(t *testing.T) {
	tests := []struct {
		name           string
		commentID      string
		userID         uint
		userRole       string
		mockComment    *models.Comment
		mockError      error
		expectedStatus int
	}{
		{
			name:      "успешное удаление комментария автором",
			commentID: "1",
			userID:    1,
			userRole:  "user",
			mockComment: &models.Comment{
				ID:       1,
				Content:  "Test Comment",
				PostID:   1,
				AuthorID: 1,
			},
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:      "успешное удаление комментария админом",
			commentID: "1",
			userID:    2,
			userRole:  "admin",
			mockComment: &models.Comment{
				ID:       1,
				Content:  "Test Comment",
				PostID:   1,
				AuthorID: 1,
			},
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:      "отказ в доступе",
			commentID: "1",
			userID:    2,
			userRole:  "user",
			mockComment: &models.Comment{
				ID:       1,
				Content:  "Test Comment",
				PostID:   1,
				AuthorID: 1,
			},
			mockError:      nil,
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.MockCommentService{
				GetCommentByIDFunc: func(id int) (*models.Comment, error) {
					return tt.mockComment, nil
				},
				DeleteCommentFunc: func(id int, userID int) error {
					return tt.mockError
				},
			}

			handler := NewCommentHandler(mockService)
			router := setupTestRouter()
			router.DELETE("/comments/:id", func(c *gin.Context) {
				c.Set("user_id", tt.userID)
				c.Set("user_role", tt.userRole)
				c.Params = []gin.Param{{Key: "id", Value: tt.commentID}}
				handler.DeleteComment(c)
			})

			req := httptest.NewRequest("DELETE", "/comments/"+tt.commentID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
} 