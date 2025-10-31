package handler

import "github.com/gin-gonic/gin"

type MessageService interface {
	Put()
}

type MessageHandler struct {
	service MessageService
}

func New(service MessageService) *MessageHandler {
	return &MessageHandler{service: service}
}

func (h *MessageHandler) Put(c *gin.Context) {

}
