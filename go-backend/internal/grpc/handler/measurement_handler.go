package handler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kangjh3kang/manpasik/go-backend/internal/domain"
	"github.com/kangjh3kang/manpasik/go-backend/internal/repository"
	pb "github.com/kangjh3kang/manpasik/go-backend/proto/api_v1"
)

// MeasurementHandler serves gRPC requests for MeasurementService.
type MeasurementHandler struct {
	pb.UnimplementedMeasurementServiceServer
	repo repository.MeasurementRepository
}

// NewMeasurementHandler creates a new gRPC handler.
func NewMeasurementHandler(repo repository.MeasurementRepository) *MeasurementHandler {
	return &MeasurementHandler{repo: repo}
}

// SyncMeasurements: 모바일 기기(Flutter) 내 Isar DB에서 오프라인 모드로 모아둔 데이터를
// 네트워크 복구 시 한방에(Bulk) 빨아들이는 CRDT 핵심 동기화 파이프라인.
func (h *MeasurementHandler) SyncMeasurements(ctx context.Context, req *pb.SyncMeasurementsRequest) (*pb.SyncMeasurementsResponse, error) {
	log.Printf("[gRPC] SyncMeasurements 요청 수신: 총 %d건 대기중", len(req.Records))

	var syncedCount int32
	var failedLocalIds []string

	for _, rec := range req.Records {
		// 16채널 Float 배열(diff_signal)과 896차원 배열(fingerprint)을 JSON DB 스키마에 맞게 압축
		diffBytes, err := json.Marshal(rec.DiffSignal)
		if err != nil {
			log.Printf("diffSignal 파싱 실패 (ID: %s): %v", rec.Id, err)
			failedLocalIds = append(failedLocalIds, rec.Id)
			continue
		}

		fpBytes, err := json.Marshal(rec.Fingerprint)
		if err != nil {
			log.Printf("fingerprint 파싱 실패 (ID: %s): %v", rec.Id, err)
			failedLocalIds = append(failedLocalIds, rec.Id)
			continue
		}

		// 임의의 유저 UUID 매핑 (실제로는 gRPC Interceptor에서 Context Auth 토큰 파싱 후 주입)
		mockUserID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

		// 도메인 엔티티 매핑
		m := &domain.Measurement{
			ID:            uuid.New(),                 // 서버 자체 고유 적재 키 발급
			UserID:        mockUserID,
			DeviceMAC:     rec.DeviceMac,
			MeasuredAt:    time.UnixMilli(rec.Timestamp),
			DiffSignal:    diffBytes,
			Fingerprint:   fpBytes,
			HealthScore:   int(rec.HealthScore),
			RiskLabel:     rec.RiskLabel,
			ClientLocalID: rec.Id,                     // 클라이언트 동기화 무결성 보장을 위한 식별자 기록
		}

		// Repository를 통한 타임스케일DB/포스트그레스 삽입
		err = h.repo.InsertMeasurement(ctx, m)
		if err != nil {
			log.Printf("DB 적재(CRDT 동기화) 실패 (LocalID: %s): %v", rec.Id, err)
			failedLocalIds = append(failedLocalIds, rec.Id)
			continue
		}

		syncedCount++
	}

	log.Printf("[gRPC] SyncMeasurements 완료: 성공 %d건, 실패 %d건", syncedCount, len(failedLocalIds))

	// 모바일에서 Sync 응답을 받으면 Isar 테이블의 isSynced = true 로 플래그 업데이트
	return &pb.SyncMeasurementsResponse{
		SyncedCount: syncedCount,
		FailedIds:   failedLocalIds,
	}, nil
}

// SubmitCheckupSession: V2.1 이후 도입될 프리미엄 종합 검진 서비스
func (h *MeasurementHandler) SubmitCheckupSession(ctx context.Context, req *pb.SubmitCheckupRequest) (*pb.SubmitCheckupResponse, error) {
	log.Printf("[gRPC] SubmitCheckupSession 기능은 준비 중입니다.")
	return &pb.SubmitCheckupResponse{
		Success:        false,
		CompositeScore: 0.0,
	}, nil
}
