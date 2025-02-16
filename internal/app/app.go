package app

import (
	"database/sql"
	"net/http"

	"github.com/DobryySoul/test-task/config"
	"github.com/DobryySoul/test-task/internal/http/routes/handlers"
	"github.com/DobryySoul/test-task/internal/http/routes/router"
	"github.com/DobryySoul/test-task/internal/repo/postgres"
	"github.com/DobryySoul/test-task/internal/service"
	"github.com/DobryySoul/test-task/pkg/logger"
	"github.com/go-playground/validator/v10"
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

	repo := postgres.NewRepository(db)
	service := service.NewSongService(repo)
	handler := handlers.NewHandler(service, *validator.New())
	r := router.NewRouter(handler)

	log.Infof("Starting server on port %s", cfg.Port)

	if err := http.ListenAndServe("localhost:"+cfg.Port, r.Router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
