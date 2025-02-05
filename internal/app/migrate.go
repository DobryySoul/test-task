package app

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	// migrate tools
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}
	DBUser := os.Getenv("DB_USER")
	DBPassword := os.Getenv("DB_PASSWORD")
	DBPORT := os.Getenv("DB_PORT")
	DBHost := os.Getenv("DB_HOST")
	DBName := os.Getenv("DB_NAME")
	fmt.Println(DBUser, DBPassword, DBPORT, DBHost, DBName)
	connURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		DBUser, DBPassword, DBHost, DBPORT, DBName)

	m, err := migrate.New("file://../migrations", connURL)
	if err != nil {
		log.Fatalf("Migration initialization failed: %v", err)
	}

	err = m.Up()

	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Migrate: no change")

		return
	}

	log.Printf("Migrate: up success")
}
