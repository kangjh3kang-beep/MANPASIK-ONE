package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = ":50080"
	}

	log.Printf("[Location] 부팅 중... (gRPC %s)", port)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("포트 바인딩 실패: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	// Health check HTTP server (metrics port 9100)
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status":"ok","service":"location-service"}`))
		})
		log.Println("[Location] Health endpoint on :9100")
		http.ListenAndServe(":9100", mux)
	}()

	go func() {
		log.Printf("[Location] gRPC 서버 시작 (%s)", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC 서버 오류: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[Location] 그레이스풀 셧다운...")
	grpcServer.GracefulStop()
	log.Println("[Location] 종료 완료")
}
