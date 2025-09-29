package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(&storagePath, "storage_path", "", "path to storage")
	flag.StringVar(&migrationsPath, "migrations_path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations_table", "migrations", "name of migrations table")
	flag.Parse()

	if storagePath == "" {
		panic("storage path is not set")
	}

	if migrationsPath == "" {
		panic("migrations path is not set")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("migrations are up to date")
			return
		}

		panic(err)
	}

	logger.Info("migrations are applied")
}
