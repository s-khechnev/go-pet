package main

import (
	"consumer/internal/config"
	"fmt"
	"github.com/pressly/goose"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.GetConfig()
	dbstring := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Database)

	db, err := goose.OpenDBWithDriver("postgres", dbstring)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}
}
