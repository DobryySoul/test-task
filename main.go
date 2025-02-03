package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/DobryySoul/test-task/config"
	"github.com/DobryySoul/test-task/internal/handlers"
	"github.com/DobryySoul/test-task/internal/service"
	"github.com/DobryySoul/test-task/internal/storage/postgres"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// @title: Music info
// @version: 0.0.1
func main() {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	if err := runMigrations(cfg); err != nil {
		log.Fatal(err)
	}

	repo := postgres.NewSongRepository(db, log)
	service := service.NewSongService(repo)
	handler := handlers.NewHandler(service)

	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Получение данных библиотеки с фильтрацией по всем полям и пагинацией
	r.GET("/songs", handler.GetSongs)
	// Получение текста песни с пагинацией по куплетам
	r.GET("/song/:id/text", handler.GetSongText)
	// Удаление песни
	r.DELETE("/song/:id", handler.DeleteSong)
	// Изменение данных песни
	r.PUT("/song/:id", handler.UpdateSong)
	// Добавление новой песни в формате JSON
	r.POST("/song", handler.CreateSong)

	// r.GET("/song/:id", handler.GetSongByID) не нужный хендлер
	// swagger

	log.Infof("Server starting on port %d", cfg.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r); err != nil {
		log.Fatal(err)
	}
}

func runMigrations(cfg *config.Config) error {
	m, err := migrate.New(
		cfg.MigrationsPath,
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName))
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		if err := m.Force(1); err != nil {
			log.Fatal("Force migration failed:", err)
		}
		if err := m.Up(); err != nil {
			log.Fatal("Second migration attempt failed:", err)
		}
	}
	return nil
}
