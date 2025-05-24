package main

import (
	"log"

	"github.com/subhammahanty235/medai/config"
	"github.com/subhammahanty235/medai/internal/api"
	"github.com/subhammahanty235/medai/internal/db"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	database, err := db.Initialize(cfg.MongoURI, cfg.DatabaseName)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize router
	router := api.SetupRouter(database, cfg)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
