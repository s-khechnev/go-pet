package service

import (
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"producer/internal/entity"
	"producer/internal/handler"
)

type BookProducer interface {
	Produce(book entity.Book) error
}

type BookService struct {
	bookProducer BookProducer
}

func New(producer BookProducer) *BookService {
	return &BookService{
		bookProducer: producer,
	}
}

var ProducerError = errors.New("produce book error")

func (s *BookService) Post(req handler.PostBookRequest) (string, error) {
	id := uuid.New().String()
	book := entity.Book{
		Id:      id,
		Title:   req.Title,
		Authors: req.Authors,
		Text:    req.Text,
	}

	if err := s.bookProducer.Produce(book); err != nil {
		slog.Error("failed to push book", slog.String("error", err.Error()))
		return uuid.Nil.String(), ProducerError
	}

	return id, nil
}
