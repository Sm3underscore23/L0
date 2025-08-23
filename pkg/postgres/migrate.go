package postgres

import (
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// migrationsUp -.
func (pg *Postgres) migrationsUp(migrationsDir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	sqlDB := stdlib.OpenDBFromPool(pg.Pool)
	defer sqlDB.Close()

	log.Println(migrationsDir)

	if err := goose.Up(sqlDB, migrationsDir); err != nil {
		return fmt.Errorf("goose up failed: %w", err)
	}

	return nil
}
