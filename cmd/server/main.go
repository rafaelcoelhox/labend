package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rafaelcoelhox/labbend/internal/app"
)

func main() {
	config := app.LoadConfig()

	application, err := app.NewApp(config)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	// Setup context para graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start application in goroutine
	go func() {
		if err := application.Start(ctx); err != nil {
			log.Printf("Application failed to start: %v", err)
			cancel()
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutdown signal received, starting graceful shutdown...")

	// Cancel context to signal shutdown
	cancel()

	// Give time for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Stop application
	done := make(chan struct{})
	go func() {
		if err := application.Stop(); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
		close(done)
	}()

	// Wait for shutdown completion or timeout
	select {
	case <-done:
		log.Println("Application shutdown completed gracefully")
	case <-shutdownCtx.Done():
		log.Println("Shutdown timeout reached, forcing exit")
	}
}
