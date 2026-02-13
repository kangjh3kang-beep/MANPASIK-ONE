package e2e

import (
	"context"
	"os"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// 서비스 엔드포인트 (환경변수로 오버라이드 가능)
var (
	gatewayAddr            = getEnvOrDefault("GATEWAY_ADDR", "localhost:8080")
	authServiceAddr        = getEnvOrDefault("AUTH_SERVICE_ADDR", "localhost:50051")
	userServiceAddr        = getEnvOrDefault("USER_SERVICE_ADDR", "localhost:50052")
	deviceServiceAddr      = getEnvOrDefault("DEVICE_SERVICE_ADDR", "localhost:50053")
	measurementServiceAddr = getEnvOrDefault("MEASUREMENT_SERVICE_ADDR", "localhost:50054")
)

// TestServiceHealth 모든 서비스 헬스체크 테스트
func TestServiceHealth(t *testing.T) {
	services := []struct {
		name string
		addr string
	}{
		{"auth-service", authServiceAddr},
		{"user-service", userServiceAddr},
		{"device-service", deviceServiceAddr},
		{"measurement-service", measurementServiceAddr},
	}

	for _, svc := range services {
		t.Run(svc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			conn, err := grpc.DialContext(ctx, svc.addr,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithBlock(),
			)
			if err != nil {
				t.Fatalf("연결 실패 %s: %v", svc.name, err)
			}
			defer conn.Close()

			client := healthpb.NewHealthClient(conn)
			resp, err := client.Check(ctx, &healthpb.HealthCheckRequest{
				Service: svc.name,
			})
			if err != nil {
				t.Fatalf("헬스체크 실패 %s: %v", svc.name, err)
			}

			if resp.Status != healthpb.HealthCheckResponse_SERVING {
				t.Errorf("%s 상태: %v, 기대: SERVING", svc.name, resp.Status)
			}

			t.Logf("✅ %s 정상 동작", svc.name)
		})
	}
}

// TestMeasurementFlow 측정 플로우 E2E 테스트
func TestMeasurementFlow(t *testing.T) {
	t.Skip("TODO: gRPC 클라이언트 구현 후 활성화")

	// 1. 디바이스 등록
	// 2. 카트리지 스캔
	// 3. 측정 시작
	// 4. 측정 데이터 수신
	// 5. 결과 확인
}

// TestDifferentialMeasurement 차동측정 계산 테스트
func TestDifferentialMeasurement(t *testing.T) {
	tests := []struct {
		name       string
		sDet       float64
		sRef       float64
		alpha      float64
		expected   float64
		tolerance  float64
	}{
		{
			name:      "기본 차동측정",
			sDet:      1.234,
			sRef:      0.012,
			alpha:     0.95,
			expected:  1.2226,  // 1.234 - 0.95 * 0.012
			tolerance: 0.0001,
		},
		{
			name:      "높은 노이즈",
			sDet:      2.5,
			sRef:      0.5,
			alpha:     0.95,
			expected:  2.025,   // 2.5 - 0.95 * 0.5
			tolerance: 0.0001,
		},
		{
			name:      "알파 1.0",
			sDet:      1.0,
			sRef:      0.1,
			alpha:     1.0,
			expected:  0.9,     // 1.0 - 1.0 * 0.1
			tolerance: 0.0001,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := differentialCorrection(tc.sDet, tc.sRef, tc.alpha)
			diff := abs(result - tc.expected)
			if diff > tc.tolerance {
				t.Errorf("차동측정 결과 = %v, 기대 = %v (오차: %v)", result, tc.expected, diff)
			}
			t.Logf("✅ S_det=%.4f, S_ref=%.4f, α=%.2f → S_corrected=%.4f", 
				tc.sDet, tc.sRef, tc.alpha, result)
		})
	}
}

// 차동측정 보정 함수
func differentialCorrection(sDet, sRef, alpha float64) float64 {
	return sDet - alpha*sRef
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func getEnvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
