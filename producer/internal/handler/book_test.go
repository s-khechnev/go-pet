package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	ctrl := gomock.NewController(t)
	bookService := NewMockBookService(ctrl)

	handler := New(bookService)

	router := gin.Default()
	router.POST("/books", handler.Post)

	tests := []struct {
		name       string
		body       string
		expectCode int
	}{
		{
			name:       "cant parse body",
			body:       "invalid json",
			expectCode: http.StatusBadRequest,
		},
		{
			name:       "invalid json",
			body:       `{"title":"", "text":"123"}`,
			expectCode: http.StatusUnprocessableEntity,
		},
		{
			name:       "success",
			body:       `{"title":"War and Peace", "authors":["Lev Tolstoi"], "text":"123"}`,
			expectCode: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectCode >= 200 && tt.expectCode <= 299 {
				bookService.EXPECT().Post(gomock.Any()).Return(uuid.New().String(), nil).Times(1)
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/books", strings.NewReader(tt.body))
			router.ServeHTTP(w, req)

			if w.Code != tt.expectCode {
				t.Errorf("expect code %d, but got %d", tt.expectCode, w.Code)
			}
		})
	}
}
