package main

import (
	"log"

	"github.com/rafaelcoelhox/labbend/internal/app"
)

func main() {
	// Load configuration
	config := app.LoadConfig()

	// Create and start application
	application, err := app.New(config)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	// Start server (blocks until shutdown)
	if err := application.Start(); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}
