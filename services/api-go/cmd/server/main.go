package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"eda-demo/api-go/internal/config"
	"eda-demo/api-go/internal/consumer"
	"eda-demo/api-go/internal/handler"
	"eda-demo/api-go/internal/service"
)

// @title Go Inventory & Billing Service API
// @version 1.0.0
// @description High-performance microservice in Go for asynchronous inventory reservation and fiscal invoicing in an Event-Driven Architecture.
// @contact.name Arquitectura EDA Demo
// @BasePath /

func main() {
	log.Println("[api-go] Starting High-Performance Inventory & Billing Service...")

	cfg := config.LoadConfig()
	billingSvc := service.NewInventoryBillingService()
	h := handler.NewHandler(billingSvc)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize and run RabbitMQ consumer in a goroutine
	orderConsumer := consumer.NewOrderConsumer(cfg, billingSvc)
	go func() {
		if err := orderConsumer.Start(ctx); err != nil {
			log.Printf("[api-go] Consumer stopped with error: %v", err)
		}
	}()
	defer orderConsumer.Close()

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("[api-go] HTTP server listening on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[api-go] HTTP server error: %v", err)
		}
	}()

	// Listen for termination signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("[api-go] Shutting down gracefully...")
	cancel() // Stop consumer

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("[api-go] HTTP server shutdown error: %v", err)
	}

	log.Println("[api-go] Service exited cleanly.")
}
