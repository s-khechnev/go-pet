package service

import "producer/internal/handler"

type MessageService struct {
}

func New() *MessageService {
	return &MessageService{}
}

func (MessageService) Put(req handler.PutMessageRequest) error {
	return nil
}
