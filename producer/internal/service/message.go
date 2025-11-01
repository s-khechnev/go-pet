package service

import (
	"errors"
	"log/slog"
	"producer/internal/entity"
	"producer/internal/handler"
)

type MessageProducer interface {
	Produce(message entity.Message) error
}

type MessageService struct {
	messageProducer MessageProducer
}

func New(producer MessageProducer) *MessageService {
	return &MessageService{
		messageProducer: producer,
	}
}

var ProducerError = errors.New("produce message error")

func (s *MessageService) Put(req handler.PutMessageRequest) error {
	if err := s.messageProducer.Produce(entity.Message{IP: req.IP, Payload: req.Message}); err != nil {
		slog.Error("failed to push message", slog.String("error", err.Error()))
		return ProducerError
	}

	return nil
}
