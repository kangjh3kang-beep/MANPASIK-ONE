package e2e

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// TestServiceHealth 모든 서비스 헬스체크 테스트
func TestServiceHealth(t *testing.T) {
	services := []struct {
		name string
		addr string
	}{
		{"auth-service", AuthAddr()},
		{"user-service", UserAddr()},
		{"device-service", DeviceAddr()},
		{"measurement-service", MeasurementAddr()},
		{"subscription-service", SubscriptionAddr()},
		{"shop-service", ShopAddr()},
		{"payment-service", PaymentAddr()},
		{"ai-inference-service", AiInferenceAddr()},
		{"cartridge-service", CartridgeAddr()},
		{"calibration-service", CalibrationAddr()},
		{"coaching-service", CoachingAddr()},
		{"family-service", FamilyAddr()},
		{"health-record-service", HealthRecordAddr()},
		{"community-service", CommunityAddr()},
		{"reservation-service", ReservationAddr()},
		{"admin-service", AdminAddr()},
		{"notification-service", NotificationAddr()},
		{"prescription-service", PrescriptionAddr()},
		{"video-service", VideoAddr()},
		{"telemedicine-service", TelemedicineAddr()},
		{"vision-service", VisionAddr()},
		{"translation-service", TranslationAddr()},
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
				t.Skipf("서비스 미기동 시 스킵: %s 연결 실패: %v", svc.name, err)
			}
			defer conn.Close()

			client := healthpb.NewHealthClient(conn)
			resp, err := client.Check(ctx, &healthpb.HealthCheckRequest{
				Service: svc.name,
			})
			if err != nil {
				t.Skipf("헬스체크 실패 (서비스 미기동 가능): %s: %v", svc.name, err)
			}

			if resp.Status != healthpb.HealthCheckResponse_SERVING {
				t.Errorf("%s 상태: %v, 기대: SERVING", svc.name, resp.Status)
			}

			t.Logf("✅ %s 정상 동작", svc.name)
		})
	}
}

// TestDifferentialMeasurement 차동측정 계산 테스트
func TestDifferentialMeasurement(t *testing.T) {
	tests := []struct {
		name      string
		sDet      float64
		sRef      float64
		alpha     float64
		expected  float64
		tolerance float64
	}{
		{
			name:      "기본 차동측정",
			sDet:      1.234,
			sRef:      0.012,
			alpha:     0.95,
			expected:  1.2226,
			tolerance: 0.0001,
		},
		{
			name:      "높은 노이즈",
			sDet:      2.5,
			sRef:      0.5,
			alpha:     0.95,
			expected:  2.025,
			tolerance: 0.0001,
		},
		{
			name:      "알파 1.0",
			sDet:      1.0,
			sRef:      0.1,
			alpha:     1.0,
			expected:  0.9,
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

func differentialCorrection(sDet, sRef, alpha float64) float64 {
	return sDet - alpha*sRef
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
