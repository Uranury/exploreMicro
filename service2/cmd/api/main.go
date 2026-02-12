package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Uranury/exploreMicro/service2/internal/handlers"
	"github.com/Uranury/exploreMicro/service2/internal/http_pack"
	"github.com/Uranury/exploreMicro/service2/internal/service"
	"github.com/Uranury/exploreMicro/service2/internal/storage"
)

func main() {
	store := storage.NewStore()

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	userServiceURL := getEnv("USER_SERVICE_URL", "http://localhost:8080")
	userClient := http_pack.NewHTTPUserClient(userServiceURL, httpClient)

	svc := service.NewService(store, userClient)

	handler := handlers.NewOrderHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", handler.CreateOrder)
	mux.HandleFunc("GET /orders/{id}", handler.GetOrder)
	mux.HandleFunc("GET /orders", handler.ListOrders)
	mux.HandleFunc("DELETE /orders/{id}", handler.CancelOrder)

	port := getEnv("PORT", "8081")
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("Starting order service on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
