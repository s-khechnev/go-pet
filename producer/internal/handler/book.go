package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

var validate = validator.New()

type BookService interface {
	Post(request PostBookRequest) (string, error)
}

type BookHandler struct {
	service BookService
}

func New(service BookService) *BookHandler {
	return &BookHandler{service: service}
}

type PostBookRequest struct {
	Title   string   `json:"title" validate:"required,gte=1,lte=255"`
	Authors []string `json:"authors" validate:"gte=1,lte=255"`
	Text    string   `json:"text" validate:"gte=1,lte=10000"`
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (h *BookHandler) Post(c *gin.Context) {
	var req PostBookRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}

	id, err := h.service.Post(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}
