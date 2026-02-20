// emergency-service: 응급 상황 관리 마이크로서비스
//
// 포트: HTTP :8080 (Health + REST API)
// 의존: 인메모리 저장소
//
// 기능:
// - 응급 상황 신고 (ReportEmergency)
// - 비상 연락처 조회 (GetEmergencyContacts)
// - 응급 설정 관리 (Get/UpdateEmergencySettings)
package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/manpasik/backend/services/emergency-service/internal/handler"
	"github.com/manpasik/backend/services/emergency-service/internal/repository/memory"
	"github.com/manpasik/backend/services/emergency-service/internal/service"
)

const serviceName = "emergency-service"

func main() {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = ":8080"
	}

	log.Printf("[%s] Starting...", serviceName)

	// 인메모리 저장소 초기화
	repo := memory.NewEmergencyRepository()

	// 서비스 레이어 초기화
	svc := service.NewEmergencyService(repo)

	// HTTP 핸들러 초기화
	h := handler.NewEmergencyHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"serving","service":"emergency-service"}`))
	})
	h.RegisterRoutes(mux)

	go func() {
		log.Printf("[%s] HTTP server on %s", serviceName, httpPort)
		if err := http.ListenAndServe(httpPort, mux); err != nil {
			log.Fatalf("[%s] HTTP server error: %v", serviceName, err)
		}
	}()

	// 시그널 대기
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	log.Printf("[%s] Received signal %v, shutting down...", serviceName, sig)
	log.Printf("[%s] Shutdown complete", serviceName)
}
