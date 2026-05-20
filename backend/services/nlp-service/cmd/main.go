// nlp-service: 자연어 처리 마이크로서비스
//
// 포트: HTTP :8080
// 의존: 없음 — 인메모리 저장소 사용
//
// 기능:
// - 건강 질의 파싱 (의도/엔티티 추출)
// - 증상 키워드 추출
// - 건강 제안 조회
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
	"github.com/manpasik/backend/services/nlp-service/internal/classifier"
	"github.com/manpasik/backend/services/nlp-service/internal/repository/memory"
	"github.com/manpasik/backend/services/nlp-service/internal/repository/postgres"
	"github.com/manpasik/backend/services/nlp-service/internal/service"
	"github.com/manpasik/backend/shared/config"
)

const serviceName = "nlp-service"

func main() {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = ":8080"
	}

	cfg := config.LoadFromEnv(serviceName)
	log.Printf("[%s] Starting v%s...", serviceName, cfg.Version)

	var repo service.NLPRepository

	if _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
		connCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, poolErr := pgxpool.New(connCtx, cfg.DB.DSN())
		connCancel()
		if poolErr != nil {
			log.Printf("[%s] DB connection failed, using memory: %v", serviceName, poolErr)
			repo = memory.NewNLPRepository()
		} else {
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
			if pingErr := pool.Ping(pingCtx); pingErr != nil {
				pingCancel()
				pool.Close()
				log.Printf("[%s] DB ping failed, using memory: %v", serviceName, pingErr)
				repo = memory.NewNLPRepository()
			} else {
				pingCancel()
				defer pool.Close()
				log.Printf("[%s] Connected to PostgreSQL", serviceName)
				repo = postgres.NewNLPRepository(pool)
			}
		}
	} else {
		repo = memory.NewNLPRepository()
	}

	svc := service.NewNLPService(repo)

	// OpenAI 인텐트 분류기: NLP_API_KEY 환경변수 기반
	if apiKey := os.Getenv("NLP_API_KEY"); apiKey != "" {
		model := os.Getenv("NLP_MODEL")
		if model == "" {
			model = "gpt-4o-mini"
		}
		baseURL := os.Getenv("NLP_BASE_URL")
		cls := classifier.NewOpenAIClassifier(classifier.OpenAIConfig{
			APIKey:  apiKey,
			Model:   model,
			BaseURL: baseURL,
		})
		svc.SetIntentClassifier(cls)
		log.Printf("[%s] OpenAI 인텐트 분류기 활성화 (model=%s)", serviceName, model)
	} else {
		log.Printf("[%s] NLP_API_KEY 미설정, 키워드 기반 분류 사용", serviceName)
	}

	_ = svc

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"serving","service":"nlp-service"}`))
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
