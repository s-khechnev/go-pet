package processor

import (
	"consumer/internal/entity"
	"log/slog"
	"strings"
)

type BookRepository interface {
	Save(book *entity.Book) error
}

type BookProcessorService struct {
	bookRepository BookRepository
}

func NewBookProcessorService(repo BookRepository) *BookProcessorService {
	return &BookProcessorService{
		bookRepository: repo,
	}
}

func (s *BookProcessorService) Process(book entity.Book) error {
	// some super complicated processing
	book.Text = strings.ToUpper(book.Text)

	if err := s.bookRepository.Save(&book); err != nil {
		slog.Error("failed to save book", slog.String("id", book.Id), slog.String("error", err.Error()))
		return err
	}
	return nil
}
