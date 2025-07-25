package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/leonvanderhaeghen/stockplatform/services/storeSvc/internal/config"
	"github.com/leonvanderhaeghen/stockplatform/services/storeSvc/internal/database"
	"github.com/leonvanderhaeghen/stockplatform/services/storeSvc/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize and start server
	srv := server.New(cfg, db)
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Store service started on port %s", cfg.Server.Port)

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down store service...")
	srv.Stop()
}
