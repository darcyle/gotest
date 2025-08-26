package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// TODO: Parse flags for different modes
	var (
		addr   = flag.String("addr", ":8080", "HTTP server address")
		dbURL  = flag.String("db", getEnvOrDefault("DATABASE_URL", ""), "Database connection string")
		enrich = flag.Bool("enrich", false, "Run enrichment worker instead of server")
	)
	flag.Parse()

	if *dbURL == "" {
		log.Fatal("DATABASE_URL environment variable or -db flag is required")
	}

	// TODO: Connect to database with proper settings
	db, err := sql.Open("postgres", *dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// TODO: Configure connection pool
	// db.SetMaxOpenConns(25)
	// db.SetMaxIdleConns(5)
	// db.SetConnMaxLifetime(time.Hour)

	// TODO: Ping database to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// TODO: Initialize store
	store := &pgStore{db: db}

	if *enrich {
		// TODO: Run enrichment worker pool
		log.Println("Starting enrichment worker...")
		// wp := WorkerPool{Workers: 3, Batch: 10, Store: store}
		// return wp.Run(context.Background())
		return
	}

	// TODO: Create HTTP server
	server := &http.Server{
		Addr:         *addr,
		Handler:      NewServer(store),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// TODO: Implement graceful shutdown
	go func() {
		log.Printf("Server starting on %s", *addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed:", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")
	
	// TODO: Implement graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
	
	log.Println("Server stopped")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}