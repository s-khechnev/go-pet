package postgresql

import (
	"consumer/internal/entity"
	"consumer/internal/storage"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

type BookStorage struct {
	pool *pgxpool.Pool
}

func NewStorage(ctx context.Context, connStr string) (*BookStorage, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	return &BookStorage{
		pool: pool,
	}, nil
}

func (s *BookStorage) SaveBook(ctx context.Context, b entity.Book) error {
	book := storage.FromModel(b)

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		"INSERT INTO books (id, title, text) VALUES ($1, $2, $3)",
		book.Id, book.Title, book.Text)
	if err != nil {
		return fmt.Errorf("failed to insert book: %w", err)
	}

	for _, authorName := range book.Authors {
		authorName = strings.TrimSpace(authorName)
		if authorName == "" {
			continue
		}

		var authorID int64

		err = tx.QueryRow(ctx,
			"SELECT id FROM authors WHERE name = $1", authorName).Scan(&authorID)
		if errors.Is(err, pgx.ErrNoRows) {
			err = tx.QueryRow(ctx,
				"INSERT INTO authors (name) VALUES ($1) RETURNING id",
				authorName).Scan(&authorID)
			if err != nil {
				return fmt.Errorf("failed to insert author %s: %w", authorName, err)
			}
		} else if err != nil {
			return fmt.Errorf("failed to query author %s: %w", authorName, err)
		}

		_, err = tx.Exec(ctx,
			"INSERT INTO book_authors (book_id, author_id) VALUES ($1, $2)",
			book.Id, authorID)
		if err != nil {
			return fmt.Errorf("failed to link book with author %s: %w", authorName, err)
		}
	}

	return tx.Commit(ctx)
}

func (s *BookStorage) GetCountBooks(ctx context.Context) (int64, error) {
	var count int64
	err := s.pool.QueryRow(ctx, "SELECT COUNT(id) FROM books").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to query count books: %w", err)
	}
	return count, nil
}

func (s *BookStorage) GetCountTextSymbols(ctx context.Context) (int64, error) {
	var count int64
	err := s.pool.QueryRow(ctx, "SELECT SUM(length(text)) FROM books").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to query count text symbols: %w", err)
	}
	return count, nil
}

func (s *BookStorage) GetCountAuthors(ctx context.Context) (int64, error) {
	var count int64
	err := s.pool.QueryRow(ctx, "SELECT COUNT(id) FROM authors").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to query count authors: %w", err)
	}
	return count, nil
}
