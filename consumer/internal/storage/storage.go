package storage

import (
	"consumer/internal/entity"
	"github.com/google/uuid"
)

type BookRow struct {
	Id      uuid.UUID
	Title   string
	Authors []string
	Text    string
}

func FromModel(e entity.Book) BookRow {
	return BookRow{
		Id:      uuid.MustParse(e.Id),
		Title:   e.Title,
		Authors: e.Authors,
		Text:    e.Text,
	}
}

func ToModel(e BookRow) entity.Book {
	return entity.Book{
		Id:      e.Id.String(),
		Title:   e.Title,
		Authors: e.Authors,
		Text:    e.Text,
	}
}
