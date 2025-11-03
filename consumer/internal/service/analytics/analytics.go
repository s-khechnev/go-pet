package analytics

import (
	"context"
	"log/slog"
)

type BookAnalyticsRepository interface {
	GetCountBooks(ctx context.Context) (int64, error)
	GetCountTextSymbols(ctx context.Context) (int64, error)
	GetCountAuthors(ctx context.Context) (int64, error)
}

type BookAnalyticsService struct {
	bookRepository BookAnalyticsRepository
}

func NewBookAnalyticsService(repo BookAnalyticsRepository) *BookAnalyticsService {
	return &BookAnalyticsService{
		bookRepository: repo,
	}
}

type Stats struct {
	CountTextSymbols int64
	CountBooks       int64
	CountAuthors     int64
}

func (s *BookAnalyticsService) GetStatistics(ctx context.Context) (Stats, error) {
	countBooks, err := s.bookRepository.GetCountBooks(ctx)
	if err != nil {
		slog.Error("failed to count books", slog.String("error", err.Error()))
		return Stats{}, err
	}

	countTextSymbols, err := s.bookRepository.GetCountTextSymbols(ctx)
	if err != nil {
		slog.Error("failed to count books", slog.String("error", err.Error()))
		return Stats{}, err
	}

	countAuthors, err := s.bookRepository.GetCountAuthors(ctx)
	if err != nil {
		slog.Error("failed to count authors", slog.String("error", err.Error()))
		return Stats{}, err
	}

	return Stats{
		CountBooks:       countBooks,
		CountTextSymbols: countTextSymbols,
		CountAuthors:     countAuthors,
	}, nil
}
