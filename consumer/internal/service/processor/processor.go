package processor

import (
	"consumer/internal/entity"
	"context"
	"log/slog"
	"strings"
	"time"
)

type BookRepository interface {
	SaveBook(ctx context.Context, b entity.Book) error
}

type BookProcessorService struct {
	bookRepository BookRepository
}

func NewBookProcessorService(repo BookRepository) *BookProcessorService {
	return &BookProcessorService{
		bookRepository: repo,
	}
}

func (s *BookProcessorService) Process(ctx context.Context, book entity.Book) error {
	// some super complicated processing
	book.Text = strings.ToUpper(book.Text)

	timeoutToSave := time.Second * 3
	ctx, cancel := context.WithTimeout(ctx, timeoutToSave)
	defer cancel()

	if err := s.bookRepository.SaveBook(ctx, book); err != nil {
		slog.Error("failed to save book", slog.String("id", book.Id), slog.String("error", err.Error()))
		return err
	}
	slog.Info("book is saved", slog.String("id", book.Id), slog.String("text", book.Text))

	return nil
}
