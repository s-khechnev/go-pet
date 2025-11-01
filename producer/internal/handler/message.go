package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

var validate = validator.New()

type MessageService interface {
	Put(request PutMessageRequest) error
}

type MessageHandler struct {
	service MessageService
}

func New(service MessageService) *MessageHandler {
	return &MessageHandler{service: service}
}

type PutMessageRequest struct {
	IP      string
	Message string `json:"message" validate:"required,gte=1,lte=1023"`
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (h *MessageHandler) Put(c *gin.Context) {
	var req PutMessageRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}

	req.IP = c.RemoteIP()

	if err := h.service.Put(req); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ip": req.IP, "message": req.Message})
}
