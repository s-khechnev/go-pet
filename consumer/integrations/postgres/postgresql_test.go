package postgres

import (
	"consumer/internal/entity"
	"consumer/internal/service/processor"
	storagePkg "consumer/internal/storage"
	"consumer/internal/storage/postgresql"
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"log"
	"slices"
	"strings"
	"testing"

	_ "github.com/lib/pq"
)

func Test(t *testing.T) {
	ctx := context.Background()

	dbName := "bookdb"
	dbUser := "postgres"
	dbPassword := "12345"
	postgresContainer, err := postgres.Run(ctx,
		"postgres:latest",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)
	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get conn string: %s", err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to create connect to postgres for migrations: %s", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set postgres dialect: %s", err)
	}

	if err := goose.Up(db, "../../migrations"); err != nil {
		log.Fatalf("failed to up migrations: %s", err)
	}
	db.Close()

	storage, err := postgresql.NewStorage(ctx, connStr)
	if err != nil {
		log.Fatalf("failed to create storage: %s", err)
	}

	processor := processor.NewBookProcessorService(storage)

	books := []entity.Book{
		{
			Id:      uuid.New().String(),
			Title:   "Some title1",
			Authors: []string{"Author1", "Author2"},
			Text:    "some text",
		},
		{
			Id:      uuid.New().String(),
			Title:   "super title 2",
			Authors: []string{"single author"},
			Text:    "super text 2",
		},
	}
	cmpBook := func(a, b entity.Book) int {
		return strings.Compare(a.Id, b.Id)
	}
	slices.SortFunc(books, cmpBook)

	for _, book := range books {
		if err := processor.Process(ctx, book); err != nil {
			log.Fatalf("failed to process book: %s", err)
		}
	}

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("failed to connect to postgres for test processors: %s", err)
	}

	rows, err := conn.Query(ctx, "SELECT * FROM books")
	if err != nil {
		log.Fatalf("failed to query books: %s", err)
	}
	defer rows.Close()
	actualBooks := make([]entity.Book, 0)
	for rows.Next() {
		var actualBook storagePkg.BookRow
		if err := rows.Scan(&actualBook.Id, &actualBook.Title, &actualBook.Text); err != nil {
			log.Fatalf("failed to scan actualBook: %s", err)
		}
		actualBooks = append(actualBooks, storagePkg.ToModel(actualBook))
	}
	slices.SortFunc(actualBooks, cmpBook)

	if len(actualBooks) != len(books) {
		log.Fatalf("expected %d books, got %d", len(books), len(actualBooks))
	}

	for i := range actualBooks {
		if actualBooks[i].Id != books[i].Id {
			log.Fatalf("expected %s got %s", actualBooks[i].Id, books[i].Id)
		}

		if actualBooks[i].Title != books[i].Title {
			log.Fatalf("expected %s got %s", actualBooks[i].Title, books[i].Title)
		}

		if actualBooks[i].Text != strings.ToUpper(books[i].Text) {
			log.Fatalf("wrong processing: expected %s got %s", actualBooks[i].Text, books[i].Text)
		}
	}
}
