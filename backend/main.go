package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FoxyHunter7/geocachingbrughia-backend/internal/config"
	"github.com/FoxyHunter7/geocachingbrughia-backend/internal/database"
	"github.com/FoxyHunter7/geocachingbrughia-backend/internal/router"
	"github.com/FoxyHunter7/geocachingbrughia-backend/internal/services/email"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists (development)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Validate configuration (fails in production with insecure settings)
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Initialize database
	db, err := database.New(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Seed default data (languages and static content)
	if err := db.SeedDefaults(); err != nil {
		log.Printf("Warning: Failed to seed defaults: %v", err)
	}

	// Initialize email service
	emailService := email.New(cfg.SMTP)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start reminder scheduler (checks every hour)
	go emailService.StartReminderScheduler(db, cfg.ReminderDays)

	// Set up router
	r := router.New(db, cfg, emailService)

	// Create HTTP server
	addr := ":" + cfg.Port
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Channel for shutdown signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on %s", addr)
		log.Printf("Database: %s", cfg.DatabasePath)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-shutdown
	log.Println("Shutdown signal received, gracefully stopping...")

	// Cancel context to stop background services
	cancel()

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 30*time.Second)
	defer shutdownCancel()

	// Gracefully shutdown server
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped gracefully")
}
