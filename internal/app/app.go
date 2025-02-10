package app

import (
	"database/sql"
	"net/http"

	"github.com/DobryySoul/test-task/config"
	"github.com/DobryySoul/test-task/internal/http/routes"
	"github.com/DobryySoul/test-task/internal/repo/postgres"
	"github.com/DobryySoul/test-task/internal/service"
	"github.com/DobryySoul/test-task/pkg/logger"

	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {
	log := logger.New()
	log.Info("Initializing application")

	db, err := sql.Open("postgres", cfg.PG.URL)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Errorf("Error closing database connection: %v", err)
		}
	}()

	log.Info("Connected to database")

	repo := postgres.NewSongRepository(db, log)
	service := service.NewSongService(repo)
	h := gin.Default()

	routes.NewRouter(h, service)

	log.Infof("Starting server on port %s", cfg.Port)

	if err := http.ListenAndServe("localhost:"+cfg.Port, h); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
