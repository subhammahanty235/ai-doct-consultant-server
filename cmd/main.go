package main

import (
	"log"

	"github.com/subhammahanty235/medai/internal/config"
	"github.com/subhammahanty235/medai/internal/db"
	"github.com/subhammahanty235/medai/internal/shared"
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
	router := shared.SetupRouter(database, cfg)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
