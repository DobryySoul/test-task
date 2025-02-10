package main

import (
	"log"

	"github.com/DobryySoul/test-task/config"
	_ "github.com/DobryySoul/test-task/docs"
	"github.com/DobryySoul/test-task/internal/app"
)

// @title Music info
// @version 0.0.1
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("can't init config: %s", err)
	}

	app.Run(cfg)
}
