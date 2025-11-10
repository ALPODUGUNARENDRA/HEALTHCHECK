package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	
	"loadbalancer/backend"
	"loadbalancer/handler"
	"loadbalancer/config"
	"loadbalancer/utils"
)

func main() {
	// Initialize loggers
	utils.InitLoggers()
	
	// Load configuration
	cfg := config.LoadConfig()
	
	// Create backend pool
	backendPool := backend.NewPool()
	
	// Add backend servers from config
	for _, server := range cfg.BackendServers {
		err := backendPool.AddBackend(server.URL, server.Weight)
		if err != nil {
			log.Fatalf("Failed to add backend %s: %v", server.URL, err)
		}
	}
	
	// Create load balancer handler
	lbHandler := handler.NewLoadBalancerHandler(backendPool, cfg.Algorithm)
	
	// Start health checks
	healthChecker := handler.NewHealthChecker(backendPool, cfg.HealthCheckInterval)
	healthChecker.Start()
	
	// Setup HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: lbHandler,
	}
	
	// Start server in goroutine
	go func() {
		log.Printf("Load balancer started on port %d", cfg.Port)
		log.Printf("Using algorithm: %s", cfg.Algorithm)
		log.Printf("Backend servers: %v", cfg.BackendServers)
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	
	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down load balancer...")
	healthChecker.Stop()
	
	if err := server.Close(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	
	log.Println("Load balancer exited")
}