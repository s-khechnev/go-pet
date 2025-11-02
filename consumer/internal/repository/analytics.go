package repository

import (
	"consumer/internal/entity"
	"log/slog"
)

type AnalyticsPostgresqlRepo struct {
}

func NewAnalyticsPostgresqlRepo() *AnalyticsPostgresqlRepo {
	return &AnalyticsPostgresqlRepo{}
}

func (r *AnalyticsPostgresqlRepo) SaveCountBook(book entity.Book) error {
	slog.Info("trying save book", slog.String("id", book.Id))
	return nil
}
