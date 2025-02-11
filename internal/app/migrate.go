package app

import (
	"errors"
	"log"
	"os"

	// migrate tools
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func init() {
	databaseURL := os.Getenv("PG_URL")
	if databaseURL == "" {
		log.Fatalf("failed to migrate database: PG_URL not found")
	}

	m, err := migrate.New("file://migrations", databaseURL)
	if err != nil {
		log.Fatalf("Migration initialization failed: %v", err)
	}

	defer m.Close()

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("Migrate: no change")

		return
	}

	log.Println("Migrate: up success")
}
