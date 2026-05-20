package handler

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/manpasik/backend/services/measurement-service/internal/service"
	"github.com/manpasik/backend/shared/assay"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type streamSessionRepo struct {
	sessions map[string]*service.MeasurementSession
}

func (r *streamSessionRepo) CreateSession(_ context.Context, session *service.MeasurementSession) error {
	r.sessions[session.ID] = session
	return nil
}

func (r *streamSessionRepo) GetSession(_ context.Context, sessionID string) (*service.MeasurementSession, error) {
	return r.sessions[sessionID], nil
}

func (r *streamSessionRepo) EndSession(_ context.Context, sessionID string, totalMeasurements int, endedAt time.Time) error {
	session := r.sessions[sessionID]
	if session == nil {
		return nil
	}
	session.TotalMeasurements = totalMeasurements
	session.EndedAt = &endedAt
	session.Status = "completed"
	return nil
}

type streamMeasureRepo struct {
	stored []*service.MeasurementData
	history []*service.MeasurementSummary
}

func (r *streamMeasureRepo) Store(_ context.Context, data *service.MeasurementData) error {
	copyData := *data
	copyData.RawChannels = append([]float64(nil), data.RawChannels...)
	copyData.FingerprintVector = append([]float32(nil), data.FingerprintVector...)
	r.stored = append(r.stored, &copyData)
	return nil
}

func (r *streamMeasureRepo) GetHistory(_ context.Context, _ string, _, _ time.Time, _, _ int) ([]*service.MeasurementSummary, int, error) {
	return r.history, len(r.history), nil
}

type streamVectorRepo struct {
	vectors map[string][]float32
}

func (r *streamVectorRepo) StoreFingerprint(_ context.Context, sessionID string, vector []float32) error {
	r.vectors[sessionID] = append([]float32(nil), vector...)
	return nil
}

func (r *streamVectorRepo) SearchSimilar(_ context.Context, _ []float32, _ int) ([]service.SimilarResult, error) {
	return nil, nil
}

type streamEventPublisher struct{}

func (p *streamEventPublisher) PublishMeasurementCompleted(_ context.Context, _ *service.MeasurementCompletedEvent) error {
	return nil
}

type measurementStreamFake struct {
	ctx       context.Context
	requests  []*v1.MeasurementData
	responses []*v1.MeasurementResult
	recvIndex int
}

func (s *measurementStreamFake) Recv() (*v1.MeasurementData, error) {
	if s.recvIndex >= len(s.requests) {
		return nil, io.EOF
	}
	request := s.requests[s.recvIndex]
	s.recvIndex++
	return request, nil
}

func (s *measurementStreamFake) Send(response *v1.MeasurementResult) error {
	s.responses = append(s.responses, response)
	return nil
}

func (s *measurementStreamFake) SetHeader(metadata.MD) error {
	return nil
}

func (s *measurementStreamFake) SendHeader(metadata.MD) error {
	return nil
}

func (s *measurementStreamFake) SetTrailer(metadata.MD) {}

func (s *measurementStreamFake) Context() context.Context {
	if s.ctx != nil {
		return s.ctx
	}
	return context.Background()
}

func (s *measurementStreamFake) SendMsg(any) error {
	return nil
}

func (s *measurementStreamFake) RecvMsg(any) error {
	return nil
}

func newStreamMeasurementHandler(session *service.MeasurementSession) (*MeasurementHandler, *streamMeasureRepo, *streamVectorRepo) {
	sessionRepo := &streamSessionRepo{
		sessions: map[string]*service.MeasurementSession{
			session.ID: session,
		},
	}
	measureRepo := &streamMeasureRepo{}
	vectorRepo := &streamVectorRepo{vectors: make(map[string][]float32)}
	svc := service.NewMeasurementService(zap.NewNop(), sessionRepo, measureRepo, vectorRepo, &streamEventPublisher{})
	return NewMeasurementHandler(svc, zap.NewNop()), measureRepo, vectorRepo
}

func TestStreamMeasurementStoresMeasurementAndFingerprint(t *testing.T) {
	handler, measureRepo, vectorRepo := newStreamMeasurementHandler(&service.MeasurementSession{
		ID:          "session-1",
		DeviceID:    "device-1",
		CartridgeID: "glucose",
		UserID:      "user-1",
		StartedAt:   time.Now().UTC(),
		Status:      "active",
	})

	stream := &measurementStreamFake{
		ctx: context.Background(),
		requests: []*v1.MeasurementData{
			{
				SessionId:   "session-1",
				RawChannels: []float64{1, 2, 3},
				Differential: &v1.DifferentialCorrection{
					SDet:       10,
					SRef:       1,
					Alpha:      0.95,
					SCorrected: 9,
				},
				EnvMeta: &v1.EnvironmentMeta{
					TempC:       24.5,
					HumidityPct: 45,
				},
			},
		},
	}

	if err := handler.StreamMeasurement(stream); err != nil {
		t.Fatalf("StreamMeasurement failed: %v", err)
	}

	if len(measureRepo.stored) != 1 {
		t.Fatalf("stored measurement count = %d, want 1", len(measureRepo.stored))
	}
	stored := measureRepo.stored[0]
	if stored.SessionID != "session-1" {
		t.Fatalf("stored SessionID = %q, want session-1", stored.SessionID)
	}
	if stored.DeviceID != "device-1" || stored.UserID != "user-1" || stored.CartridgeType != "glucose" {
		t.Fatalf("session metadata not backfilled: device=%q user=%q cartridge=%q", stored.DeviceID, stored.UserID, stored.CartridgeType)
	}
	if stored.PrimaryValue != 9 || stored.SCorrected != 9 || stored.Unit != "mg/dL" {
		t.Fatalf("stored result fields mismatch: primary=%f corrected=%f unit=%q confidence=%f", stored.PrimaryValue, stored.SCorrected, stored.Unit, stored.Confidence)
	}
	if stored.Confidence <= 0.90 || stored.Confidence >= 0.95 {
		t.Fatalf("stored confidence = %f, want assay-derived low completeness confidence", stored.Confidence)
	}
	if stored.TempC != 24.5 || stored.HumidityPct != 45 {
		t.Fatalf("environment metadata mismatch: temp=%f humidity=%f", stored.TempC, stored.HumidityPct)
	}
	assertFloat64Slice(t, stored.RawChannels, []float64{1, 2, 3})
	assertFloat32Slice(t, stored.FingerprintVector, []float32{1, 2, 3})
	assertFloat32Slice(t, vectorRepo.vectors["session-1"], []float32{1, 2, 3})

	if len(stream.responses) != 1 {
		t.Fatalf("stream response count = %d, want 1", len(stream.responses))
	}
	response := stream.responses[0]
	if response.SessionId != "session-1" || response.PrimaryValue != 9 || response.Unit != "mg/dL" {
		t.Fatalf("response mismatch: session=%q primary=%f unit=%q confidence=%f", response.SessionId, response.PrimaryValue, response.Unit, response.Confidence)
	}
	if response.Confidence != stored.Confidence {
		t.Fatalf("response confidence = %f, want stored confidence %f", response.Confidence, stored.Confidence)
	}
	if response.EvidenceStatus != "research_only" {
		t.Fatalf("response EvidenceStatus = %q, want research_only", response.EvidenceStatus)
	}
	if response.DiagnosticReady {
		t.Fatal("research-only response must not be diagnostic ready")
	}
	assertContainsString(t, response.EvidenceGaps, "clinical_lock_required")
	assertFloat32Slice(t, response.FingerprintVector, []float32{1, 2, 3})
	if response.ProcessedAt == nil {
		t.Fatal("response ProcessedAt is nil")
	}
}

func TestStreamMeasurementUsesAssayUnitFromSessionCartridge(t *testing.T) {
	handler, measureRepo, _ := newStreamMeasurementHandler(&service.MeasurementSession{
		ID:          "session-crp",
		DeviceID:    "device-1",
		CartridgeID: "crp",
		UserID:      "user-1",
		StartedAt:   time.Now().UTC(),
		Status:      "active",
	})

	stream := &measurementStreamFake{
		ctx: context.Background(),
		requests: []*v1.MeasurementData{
			{
				SessionId:   "session-crp",
				RawChannels: make([]float64, 88),
				Differential: &v1.DifferentialCorrection{
					SDet:       5,
					SRef:       1,
					Alpha:      0.98,
					SCorrected: 4.2,
				},
				EnvMeta: &v1.EnvironmentMeta{},
			},
		},
	}

	if err := handler.StreamMeasurement(stream); err != nil {
		t.Fatalf("StreamMeasurement failed: %v", err)
	}
	if len(measureRepo.stored) != 1 {
		t.Fatalf("stored measurement count = %d, want 1", len(measureRepo.stored))
	}
	if measureRepo.stored[0].Unit != "mg/L" {
		t.Fatalf("stored CRP unit = %q, want mg/L", measureRepo.stored[0].Unit)
	}
	if len(stream.responses) != 1 || stream.responses[0].Unit != "mg/L" {
		t.Fatalf("response CRP unit mismatch: responses=%v", stream.responses)
	}
}

func TestStreamMeasurementRequiresSessionID(t *testing.T) {
	handler, measureRepo, vectorRepo := newStreamMeasurementHandler(&service.MeasurementSession{
		ID:        "session-1",
		StartedAt: time.Now().UTC(),
		Status:    "active",
	})
	stream := &measurementStreamFake{
		ctx:      context.Background(),
		requests: []*v1.MeasurementData{{RawChannels: []float64{1}}},
	}

	err := handler.StreamMeasurement(stream)
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("status code = %v, want %v; err=%v", status.Code(err), codes.InvalidArgument, err)
	}
	if len(measureRepo.stored) != 0 {
		t.Fatalf("stored measurement count = %d, want 0", len(measureRepo.stored))
	}
	if len(vectorRepo.vectors) != 0 {
		t.Fatalf("stored vector count = %d, want 0", len(vectorRepo.vectors))
	}
}

func TestGetMeasurementHistoryIncludesEvidenceFields(t *testing.T) {
	handler, measureRepo, _ := newStreamMeasurementHandler(&service.MeasurementSession{
		ID:        "session-history",
		StartedAt: time.Now().UTC(),
		Status:    "completed",
	})
	measureRepo.history = []*service.MeasurementSummary{
		{
			SessionID:       "session-history",
			CartridgeType:   "glucose",
			PrimaryValue:    99.5,
			Unit:            "mg/dL",
			EvidenceStatus:  assay.EvidenceStatusResearchOnly,
			DiagnosticReady: false,
			EvidenceGaps:    []string{"clinical_lock_required"},
			MeasuredAt:      time.Now().UTC(),
		},
	}

	response, err := handler.GetMeasurementHistory(context.Background(), &v1.GetHistoryRequest{
		UserId: "user-1",
		Limit:  10,
	})
	if err != nil {
		t.Fatalf("GetMeasurementHistory failed: %v", err)
	}
	if len(response.Measurements) != 1 {
		t.Fatalf("history count = %d, want 1", len(response.Measurements))
	}
	summary := response.Measurements[0]
	if summary.EvidenceStatus != "research_only" {
		t.Fatalf("EvidenceStatus = %q, want research_only", summary.EvidenceStatus)
	}
	if summary.DiagnosticReady {
		t.Fatal("research_only history summary must not be diagnostic ready")
	}
	assertContainsString(t, summary.EvidenceGaps, "clinical_lock_required")
}

func assertFloat64Slice(t *testing.T, got, want []float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("float64 slice length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("float64 slice[%d] = %f, want %f", i, got[i], want[i])
		}
	}
}

func assertFloat32Slice(t *testing.T, got, want []float32) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("float32 slice length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("float32 slice[%d] = %f, want %f", i, got[i], want[i])
		}
	}
}

func assertContainsString(t *testing.T, values []string, want string) {
	t.Helper()
	for _, value := range values {
		if value == want {
			return
		}
	}
	t.Fatalf("values = %v, want %q", values, want)
}
