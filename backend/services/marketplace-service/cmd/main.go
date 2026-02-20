// marketplace-service: 파트너 마켓플레이스 마이크로서비스
//
// 포트: HTTP :8080 (Health + REST API)
// 의존: 인메모리 저장소
//
// 기능:
// - 파트너 상품 조회 (카테고리/파트너 필터)
// - 파트너 등록
// - 파트너 통계 조회
// - 상품 업데이트
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/manpasik/backend/services/marketplace-service/internal/repository/memory"
	"github.com/manpasik/backend/services/marketplace-service/internal/service"
)

const serviceName = "marketplace-service"

func main() {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = ":8080"
	}

	log.Printf("[%s] Starting...", serviceName)

	// 인메모리 저장소 초기화
	productRepo := memory.NewProductRepository()
	partnerRepo := memory.NewPartnerRepository()
	statsRepo := memory.NewStatsRepository()

	// 서비스 레이어 초기화
	svc := service.NewMarketplaceService(productRepo, partnerRepo, statsRepo)

	// HTTP 라우터
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"serving","service":"marketplace-service"}`))
	})

	// REST 엔드포인트
	mux.HandleFunc("/api/v1/marketplace/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		partnerID := r.URL.Query().Get("partner_id")
		category := r.URL.Query().Get("category")
		products, err := svc.ListPartnerProducts(r.Context(), partnerID, category)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	})

	mux.HandleFunc("/api/v1/marketplace/partners", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var partner service.Partner
		if err := json.NewDecoder(r.Body).Decode(&partner); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		id, err := svc.RegisterPartner(r.Context(), &partner)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"partner_id": id})
	})

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
