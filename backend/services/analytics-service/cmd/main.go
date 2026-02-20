package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/manpasik/backend/services/analytics-service/internal/repository/memory"
	"github.com/manpasik/backend/services/analytics-service/internal/service"
)

const serviceName = "analytics-service"

func main() {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = ":8080"
	}

	log.Printf("[%s] Starting...", serviceName)

	repo := memory.NewAnalyticsRepository()
	svc := service.NewAnalyticsService(repo)
	_ = svc

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"serving","service":"analytics-service"}`))
	})

	go func() {
		log.Printf("[%s] HTTP server on %s", serviceName, httpPort)
		if err := http.ListenAndServe(httpPort, mux); err != nil {
			log.Fatalf("[%s] HTTP server error: %v", serviceName, err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	log.Printf("[%s] Received signal %v, shutting down...", serviceName, sig)
}
