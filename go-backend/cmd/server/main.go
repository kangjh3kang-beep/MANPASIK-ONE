package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/kangjh3kang/manpasik/go-backend/internal/grpc/handler"
	"github.com/kangjh3kang/manpasik/go-backend/internal/repository"
	pb "github.com/kangjh3kang/manpasik/go-backend/proto/api_v1"
)

const (
	port = ":50051"
)

func main() {
	log.Println("[ManPaSik] Go Backend (v2.4.3) 부팅 중...")

	// 1. PostgreSQL/TimescaleDB 커넥션 풀 연결 (가상 URL, 실제 운영 시 환경변수 치환 권장)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// (주의: 모의 스키마. Postgres 환경이 준비되지 않았다면 에러가 날 수 있으므로 에러 핸들링 유연하게 처리)
	dbURL := "postgres://postgres:password@localhost:5432/manpasik?sslmode=disable"
	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Printf("[경고] 데이터베이스 연결 실패 (운영 환경 시 필수 구성 요소): %v\n", err)
		log.Println("[INFO] DB가 없지만 gRPC 서버 포트는 스텁 모드로 개방합니다 (개발용).")
	} else {
		defer dbpool.Close()
		log.Println("[SUCCESS] PostgreSQL 연결 성공 (CRDT 동기화 채널 확보)")
	}

	// 2. 레퍼지토리 및 핸들러(컨트롤러) 인스턴스화
	measureRepo := repository.NewMeasurementRepository(dbpool)
	measureHandler := handler.NewMeasurementHandler(measureRepo)

	// 3. gRPC 서버 구동
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("포트 바인딩 실패: %v", err)
	}

	grpcServer := grpc.NewServer()
	
	// MeasurementService API 인터페이스 등록
	pb.RegisterMeasurementServiceServer(grpcServer, measureHandler)

	// Postman이나 gRPC 호환 클라이언트 도구에서 메서드 조회가 가능하도록 리플렉션 허용
	reflection.Register(grpcServer)

	// 4. 안전한 서버 종료(Graceful Shutdown) 처리를 위한 고루틴
	go func() {
		log.Printf("[SUCCESS] gRPC 서버가 포트 %s에서 수신을 시작했습니다\n", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC 서버 크래시: %v", err)
		}
	}()

	// 5. [신규] HTTP/REST 브릿지 서버 가동 (Next.js 연동용)
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write([]byte(`{"status":"ok", "version":"2.4.3", "service":"ManPaSik Core"}`))
		})

		// 임상 데이터 REST 엔드포인트 목업 (실제 Repo 연결 가능)
		mux.HandleFunc("/api/v1/patients/vitals", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write([]byte(`{"data":[{"timestamp":"08:00","heartRate":72,"bloodPressureSys":120,"spo2":98}]}`))
		})

		log.Println("[SUCCESS] HTTP REST 게이트웨이가 포트 :8080에서 가동 중입니다.")
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Printf("[ERROR] HTTP 서버 오류: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[INFO] 서버 종료 시그널 수신, 그레이스풀 셧다운 트리거...")
	grpcServer.GracefulStop()
	log.Println("[SUCCESS] 서버가 안전하게 종료되었습니다.")
}
