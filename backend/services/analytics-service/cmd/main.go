package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/analytics-service/internal/handler"
	"github.com/manpasik/backend/services/analytics-service/internal/repository/memory"
	"github.com/manpasik/backend/services/analytics-service/internal/repository/postgres"
	"github.com/manpasik/backend/services/analytics-service/internal/service"
	"github.com/manpasik/backend/shared/config"
)

const serviceName = "analytics-service"

func main() {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = ":8080"
	}

	cfg := config.LoadFromEnv(serviceName)
	log.Printf("[%s] Starting v%s...", serviceName, cfg.Version)

	var repo service.AnalyticsRepository

	if _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
		connCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, poolErr := pgxpool.New(connCtx, cfg.DB.DSN())
		connCancel()
		if poolErr != nil {
			log.Printf("[%s] DB connection failed, using memory: %v", serviceName, poolErr)
			repo = memory.NewAnalyticsRepository()
		} else {
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
			if pingErr := pool.Ping(pingCtx); pingErr != nil {
				pingCancel()
				pool.Close()
				log.Printf("[%s] DB ping failed, using memory: %v", serviceName, pingErr)
				repo = memory.NewAnalyticsRepository()
			} else {
				pingCancel()
				defer pool.Close()
				log.Printf("[%s] Connected to PostgreSQL", serviceName)
				repo = postgres.NewAnalyticsRepository(pool)
			}
		}
	} else {
		repo = memory.NewAnalyticsRepository()
	}

	svc := service.NewAnalyticsService(repo)

	// REST 핸들러 초기화 + 라우트 등록
	h := handler.NewAnalyticsHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"serving","service":"analytics-service"}`))
	})
	h.RegisterRoutes(mux)

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
