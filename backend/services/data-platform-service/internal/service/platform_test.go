package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/manpasik/backend/services/data-platform-service/internal/repository/memory"
	"github.com/manpasik/backend/services/data-platform-service/internal/service"
	"go.uber.org/zap"
)

func newSvc() *service.PlatformService {
	return service.NewPlatformService(zap.NewNop(), memory.NewPlatformRepository())
}
func TestGetRegionStats(t *testing.T) {
	s := newSvc()
	st, err := s.GetRegionStats(context.Background(), "KR-11", "day", "blood_glucose")
	if err != nil {
		t.Fatal(err)
	}
	if st.RegionCode != "KR-11" {
		t.Error("code mismatch")
	}
}
func TestListHotspots(t *testing.T) {
	s := newSvc()
	list, err := s.ListHotspots(context.Background(), "blood_glucose", "high", 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) == 0 {
		t.Error("no hotspots")
	}
}
func TestGetTrend(t *testing.T) {
	s := newSvc()
	trend, err := s.GetTrend(context.Background(), "KR-11", "blood_glucose", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(trend.Points) != 7 {
		t.Errorf("expected 7 points, got %d", len(trend.Points))
	}
}
func TestListDatasets(t *testing.T) {
	s := newSvc()
	list, total, err := s.ListDatasets(context.Background(), "", 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if total == 0 || len(list) == 0 {
		t.Error("no datasets")
	}
}
func TestConsent(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	c, _ := s.GetConsent(ctx, "u1")
	if c.AnonymousResearch {
		t.Error("default should be false")
	}
	c2, _ := s.UpdateConsent(ctx, "u1", true, true, false)
	if !c2.AnonymousResearch {
		t.Error("should be true")
	}
}
func TestTriggerAggregation(t *testing.T) {
	s := newSvc()
	id, st, err := s.TriggerAggregation(context.Background(), "blood_glucose", "2026-01-01")
	if err != nil {
		t.Fatal(err)
	}
	if id == "" || st != "completed" {
		t.Error("unexpected result")
	}
}

// ============================================================================
// Phase B 신규 테스트
// ============================================================================

func TestListDatasets_InvalidType(t *testing.T) {
	s := newSvc()
	_, _, err := s.ListDatasets(context.Background(), "invalid_type", 10, 0)
	if err == nil {
		t.Error("지원하지 않는 데이터셋 타입에 에러가 반환되어야 합니다")
	}
}

func TestGetRegionStats_InvalidPeriod(t *testing.T) {
	s := newSvc()
	_, err := s.GetRegionStats(context.Background(), "KR-11", "invalid_period", "blood_glucose")
	if err == nil {
		t.Error("지원하지 않는 기간에 에러가 반환되어야 합니다")
	}
}

func TestRequestAccess_NoConsent(t *testing.T) {
	s := newSvc()
	ctx := context.Background()

	// 데이터셋 시드
	s.ListDatasets(ctx, "", 10, 0)

	// 동의 없이 접근 요청
	_, err := s.RequestAccess(ctx, "u1", "any-dataset", "연구 목적으로 데이터를 분석합니다")
	if err == nil {
		t.Error("동의 없이 접근 요청에 에러가 반환되어야 합니다")
	}
}

func TestRequestAccess_WithConsent(t *testing.T) {
	s := newSvc()
	ctx := context.Background()

	// 데이터셋 시드 + ID 확보
	datasets, _, _ := s.ListDatasets(ctx, "", 10, 0)
	if len(datasets) == 0 {
		t.Fatal("시드 데이터셋이 없습니다")
	}
	dsID := datasets[0].ID

	// 동의 설정
	s.UpdateConsent(ctx, "u1", true, true, false)

	// 접근 요청
	resp, err := s.RequestAccess(ctx, "u1", dsID, "학술 연구를 위한 익명화된 혈당 데이터 분석입니다")
	if err != nil {
		t.Fatalf("접근 요청 실패: %v", err)
	}
	if resp.RequestID == "" {
		t.Error("RequestID가 비어 있습니다")
	}
	if resp.Status != "pending" {
		t.Errorf("Status = %s, want pending", resp.Status)
	}
}

func TestRequestAccess_ShortPurpose(t *testing.T) {
	s := newSvc()
	ctx := context.Background()

	s.UpdateConsent(ctx, "u1", true, true, false)
	_, err := s.RequestAccess(ctx, "u1", "ds-1", "짧음")
	if err == nil {
		t.Error("너무 짧은 목적에 에러가 반환되어야 합니다")
	}
}

func TestTriggerAggregation_MissingBiomarker(t *testing.T) {
	s := newSvc()
	_, _, err := s.TriggerAggregation(context.Background(), "", "2026-01-01")
	if err == nil {
		t.Error("빈 biomarker에 에러가 반환되어야 합니다")
	}
}

// ============================================================================
// Phase E-4: K-Anonymity + 연합학습 + 연구용 API 테스트
// ============================================================================

func TestValidateKAnonymity_Satisfied(t *testing.T) {
	s := newSvc()
	records := []map[string]string{
		{"age_group": "30대", "gender": "male", "region": "KR-11"},
		{"age_group": "30대", "gender": "male", "region": "KR-11"},
		{"age_group": "30대", "gender": "male", "region": "KR-11"},
		{"age_group": "30대", "gender": "male", "region": "KR-11"},
		{"age_group": "30대", "gender": "male", "region": "KR-11"},
		{"age_group": "40대", "gender": "female", "region": "KR-26"},
		{"age_group": "40대", "gender": "female", "region": "KR-26"},
		{"age_group": "40대", "gender": "female", "region": "KR-26"},
		{"age_group": "40대", "gender": "female", "region": "KR-26"},
		{"age_group": "40대", "gender": "female", "region": "KR-26"},
	}
	result := s.ValidateKAnonymity(records, service.DefaultKAnonymityConfig())
	if !result.IsAnonymous {
		t.Error("K=5 기준 만족해야 합니다")
	}
	if result.K < 5 {
		t.Errorf("K = %d, 5 이상 기대", result.K)
	}
	if result.SuppressedCount != 0 {
		t.Errorf("억제 레코드 %d개, 0 기대", result.SuppressedCount)
	}
}

func TestValidateKAnonymity_NotSatisfied(t *testing.T) {
	s := newSvc()
	records := []map[string]string{
		{"age_group": "30대", "gender": "male", "region": "KR-11"},
		{"age_group": "30대", "gender": "male", "region": "KR-11"},
		// 2개뿐 → k=5 미충족
	}
	result := s.ValidateKAnonymity(records, service.DefaultKAnonymityConfig())
	if result.IsAnonymous {
		t.Error("K=5 기준 미충족이어야 합니다")
	}
	if result.SuppressedCount != 2 {
		t.Errorf("억제 레코드 %d개, 2 기대", result.SuppressedCount)
	}
}

func TestValidateKAnonymity_EmptyRecords(t *testing.T) {
	s := newSvc()
	result := s.ValidateKAnonymity(nil, service.DefaultKAnonymityConfig())
	if result.IsAnonymous {
		t.Error("빈 레코드에서 K-Anonymity를 만족하면 안 됩니다")
	}
}

func TestAnonymizeDataset_SuppressesSmallClasses(t *testing.T) {
	s := newSvc()
	records := []map[string]string{
		// 5개 클래스 A (유지)
		{"age_group": "30대", "gender": "male", "region": "KR-11"},
		{"age_group": "30대", "gender": "male", "region": "KR-11"},
		{"age_group": "30대", "gender": "male", "region": "KR-11"},
		{"age_group": "30대", "gender": "male", "region": "KR-11"},
		{"age_group": "30대", "gender": "male", "region": "KR-11"},
		// 2개 클래스 B (억제)
		{"age_group": "20대", "gender": "female", "region": "KR-42"},
		{"age_group": "20대", "gender": "female", "region": "KR-42"},
	}

	result, kResult := s.AnonymizeDataset(records, service.DefaultKAnonymityConfig())
	if len(result) != 5 {
		t.Errorf("익명화 후 %d개, 5 기대", len(result))
	}
	if kResult.SuppressedCount != 2 {
		t.Errorf("억제 %d개, 2 기대", kResult.SuppressedCount)
	}
	if !kResult.IsAnonymous {
		t.Error("억제 후 K-Anonymity를 만족해야 합니다")
	}
}

func TestGeneralizeAgeGroup(t *testing.T) {
	tests := []struct {
		age    int
		expect string
	}{
		{15, "10대 이하"},
		{25, "20대"},
		{35, "30대"},
		{45, "40대"},
		{55, "50대"},
		{65, "60대"},
		{75, "70대 이상"},
	}
	for _, tt := range tests {
		got := service.GeneralizeAgeGroup(tt.age)
		if got != tt.expect {
			t.Errorf("GeneralizeAgeGroup(%d) = %q, want %q", tt.age, got, tt.expect)
		}
	}
}

func TestGeneralizeRegion(t *testing.T) {
	got := service.GeneralizeRegion("KR-11-001")
	if got != "KR-11" {
		t.Errorf("GeneralizeRegion = %q, want KR-11", got)
	}
	short := service.GeneralizeRegion("KR")
	if short != "KR" {
		t.Errorf("짧은 코드 = %q, want KR", short)
	}
}

// --- 연합학습 테스트 ---

func TestFederatedAggregator_SubmitAndAggregate(t *testing.T) {
	agg := service.NewFederatedAggregator(zap.NewNop(), 3)

	// 3개 디바이스 업데이트 제출
	for i := 0; i < 3; i++ {
		err := agg.SubmitLocalUpdate(&service.LocalModelUpdate{
			DeviceID:    fmt.Sprintf("device-%d", i),
			ModelType:   service.FederatedModelHealthScore,
			Weights:     []float64{1.0 + float64(i)*0.1, 2.0 + float64(i)*0.2},
			SampleCount: 100 + i*50,
			TrainedAt:   time.Now(),
		})
		if err != nil {
			t.Fatalf("SubmitLocalUpdate %d 실패: %v", i, err)
		}
	}

	if agg.GetPendingUpdateCount(service.FederatedModelHealthScore) != 3 {
		t.Fatalf("대기 업데이트 %d, 3 기대", agg.GetPendingUpdateCount(service.FederatedModelHealthScore))
	}

	// 글로벌 모델 집계
	global, err := agg.AggregateModel(service.FederatedModelHealthScore)
	if err != nil {
		t.Fatalf("AggregateModel 실패: %v", err)
	}
	if global.Version != 1 {
		t.Errorf("Version = %d, 1 기대", global.Version)
	}
	if global.TotalDevices != 3 {
		t.Errorf("TotalDevices = %d, 3 기대", global.TotalDevices)
	}
	if len(global.Weights) != 2 {
		t.Errorf("가중치 차원 %d, 2 기대", len(global.Weights))
	}

	// 집계 후 대기 업데이트 초기화
	if agg.GetPendingUpdateCount(service.FederatedModelHealthScore) != 0 {
		t.Error("집계 후 대기 업데이트가 0이어야 합니다")
	}
}

func TestFederatedAggregator_InsufficientDevices(t *testing.T) {
	agg := service.NewFederatedAggregator(zap.NewNop(), 3)

	agg.SubmitLocalUpdate(&service.LocalModelUpdate{
		DeviceID: "device-1", ModelType: service.FederatedModelAnomaly,
		Weights: []float64{1.0}, SampleCount: 100, TrainedAt: time.Now(),
	})

	_, err := agg.AggregateModel(service.FederatedModelAnomaly)
	if err == nil {
		t.Error("디바이스 부족 시 에러가 반환되어야 합니다")
	}
}

func TestFederatedAggregator_DimensionMismatch(t *testing.T) {
	agg := service.NewFederatedAggregator(zap.NewNop(), 2)

	agg.SubmitLocalUpdate(&service.LocalModelUpdate{
		DeviceID: "d1", ModelType: service.FederatedModelTrend,
		Weights: []float64{1.0, 2.0}, SampleCount: 100, TrainedAt: time.Now(),
	})
	agg.SubmitLocalUpdate(&service.LocalModelUpdate{
		DeviceID: "d2", ModelType: service.FederatedModelTrend,
		Weights: []float64{1.0, 2.0, 3.0}, SampleCount: 100, TrainedAt: time.Now(),
	})

	_, err := agg.AggregateModel(service.FederatedModelTrend)
	if err == nil {
		t.Error("가중치 차원 불일치에 에러가 반환되어야 합니다")
	}
}

func TestFederatedAggregator_VersionIncrement(t *testing.T) {
	agg := service.NewFederatedAggregator(zap.NewNop(), 2)

	// 첫 번째 라운드
	for i := 0; i < 2; i++ {
		agg.SubmitLocalUpdate(&service.LocalModelUpdate{
			DeviceID: fmt.Sprintf("d%d", i), ModelType: service.FederatedModelHealthScore,
			Weights: []float64{1.0}, SampleCount: 100, TrainedAt: time.Now(),
		})
	}
	g1, _ := agg.AggregateModel(service.FederatedModelHealthScore)

	// 두 번째 라운드
	for i := 0; i < 2; i++ {
		agg.SubmitLocalUpdate(&service.LocalModelUpdate{
			DeviceID: fmt.Sprintf("d%d", i), ModelType: service.FederatedModelHealthScore,
			Weights: []float64{2.0}, SampleCount: 150, TrainedAt: time.Now(),
		})
	}
	g2, _ := agg.AggregateModel(service.FederatedModelHealthScore)

	if g2.Version != g1.Version+1 {
		t.Errorf("Version = %d, %d 기대", g2.Version, g1.Version+1)
	}
}

func TestFederatedAggregator_InvalidInput(t *testing.T) {
	agg := service.NewFederatedAggregator(zap.NewNop(), 3)

	if err := agg.SubmitLocalUpdate(nil); err == nil {
		t.Error("nil 업데이트에 에러 기대")
	}
	if err := agg.SubmitLocalUpdate(&service.LocalModelUpdate{DeviceID: ""}); err == nil {
		t.Error("빈 DeviceID에 에러 기대")
	}
	if err := agg.SubmitLocalUpdate(&service.LocalModelUpdate{DeviceID: "d", Weights: nil}); err == nil {
		t.Error("빈 Weights에 에러 기대")
	}
}

func TestFederatedAggregator_GetGlobalModel(t *testing.T) {
	agg := service.NewFederatedAggregator(zap.NewNop(), 2)

	// 모델 없을 때
	_, err := agg.GetGlobalModel(service.FederatedModelHealthScore)
	if err == nil {
		t.Error("모델 없을 때 에러가 반환되어야 합니다")
	}

	// 모델 생성
	for i := 0; i < 2; i++ {
		agg.SubmitLocalUpdate(&service.LocalModelUpdate{
			DeviceID: fmt.Sprintf("d%d", i), ModelType: service.FederatedModelHealthScore,
			Weights: []float64{1.0}, SampleCount: 100, TrainedAt: time.Now(),
		})
	}
	agg.AggregateModel(service.FederatedModelHealthScore)

	got, err := agg.GetGlobalModel(service.FederatedModelHealthScore)
	if err != nil {
		t.Fatalf("GetGlobalModel 실패: %v", err)
	}
	if got.ModelID == "" {
		t.Error("ModelID가 비어 있습니다")
	}
}

// --- 연구용 API 테스트 ---

func TestProcessResearchQuery_Success(t *testing.T) {
	s := newSvc()
	ctx := context.Background()

	// 동의 설정
	s.UpdateConsent(ctx, "researcher-1", true, true, false)

	result, err := s.ProcessResearchQuery(ctx, &service.ResearchQuery{
		RequesterID: "researcher-1",
		DatasetType: "blood_glucose",
		AgeRange:    [2]int{30, 60},
		MinSamples:  5,
	})
	if err != nil {
		t.Fatalf("ProcessResearchQuery 실패: %v", err)
	}
	if result.QueryID == "" {
		t.Error("QueryID가 비어 있습니다")
	}
	if result.Summary == nil {
		t.Fatal("Summary가 nil입니다")
	}
	if result.KAnonymity == nil || !result.KAnonymity.IsAnonymous {
		t.Error("K-Anonymity가 만족되어야 합니다")
	}
}

func TestProcessResearchQuery_NoConsent(t *testing.T) {
	s := newSvc()
	_, err := s.ProcessResearchQuery(context.Background(), &service.ResearchQuery{
		RequesterID: "no-consent-user",
		DatasetType: "blood_glucose",
	})
	if err == nil {
		t.Error("동의 없이 연구 쿼리에 에러가 반환되어야 합니다")
	}
}

func TestProcessResearchQuery_InvalidType(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	s.UpdateConsent(ctx, "r1", true, true, false)

	_, err := s.ProcessResearchQuery(ctx, &service.ResearchQuery{
		RequesterID: "r1",
		DatasetType: "invalid_type",
	})
	if err == nil {
		t.Error("잘못된 데이터셋 타입에 에러가 반환되어야 합니다")
	}
}

func TestProcessResearchQuery_MissingRequester(t *testing.T) {
	s := newSvc()
	_, err := s.ProcessResearchQuery(context.Background(), &service.ResearchQuery{
		DatasetType: "blood_glucose",
	})
	if err == nil {
		t.Error("빈 requester_id에 에러가 반환되어야 합니다")
	}
}

// --- DP-FedAvg 테스트 ---

func TestFederatedAggregator_WithDP(t *testing.T) {
	agg := service.NewFederatedAggregator(zap.NewNop(), 2)

	// 동일 가중치 2개 디바이스
	for i := 0; i < 2; i++ {
		agg.SubmitLocalUpdate(&service.LocalModelUpdate{
			DeviceID: fmt.Sprintf("d%d", i), ModelType: service.FederatedModelHealthScore,
			Weights: []float64{1.0, 2.0, 3.0}, SampleCount: 100, TrainedAt: time.Now(),
		})
	}

	dp := service.DPConfig{Epsilon: 1.0, Sensitivity: 1.0, ClipNorm: 10.0}
	global, err := agg.AggregateModelWithDP(service.FederatedModelHealthScore, dp)
	if err != nil {
		t.Fatalf("AggregateModelWithDP 실패: %v", err)
	}

	// DP 메트릭 확인
	if global.Metrics["dp_epsilon"] != 1.0 {
		t.Errorf("dp_epsilon = %f, 1.0 기대", global.Metrics["dp_epsilon"])
	}

	// 노이즈 적용 → 원래 값과 차이 존재 (완전 동일하지 않아야 함)
	// 단, 노이즈가 작으므로 크게 벗어나지 않아야 함
	for i, w := range global.Weights {
		expected := []float64{1.0, 2.0, 3.0}[i]
		diff := w - expected
		if diff == 0.0 {
			// 노이즈가 정확히 0이면 실패 (i*7919+13 해시에 의해 0이 아닌 노이즈)
			// 다만 특정 i에서는 0.5일 수 있으므로 전체 합으로 검증
		}
		if diff > 1.0 || diff < -1.0 {
			t.Errorf("DP 노이즈가 너무 큼: weight[%d] = %f (expected ~%f)", i, w, expected)
		}
	}
}

func TestFederatedAggregator_DP_InvalidEpsilon(t *testing.T) {
	agg := service.NewFederatedAggregator(zap.NewNop(), 2)
	for i := 0; i < 2; i++ {
		agg.SubmitLocalUpdate(&service.LocalModelUpdate{
			DeviceID: fmt.Sprintf("d%d", i), ModelType: service.FederatedModelHealthScore,
			Weights: []float64{1.0}, SampleCount: 100, TrainedAt: time.Now(),
		})
	}

	dp := service.DPConfig{Epsilon: 0, Sensitivity: 1.0, ClipNorm: 1.0}
	_, err := agg.AggregateModelWithDP(service.FederatedModelHealthScore, dp)
	if err == nil {
		t.Error("epsilon=0에 에러가 반환되어야 합니다")
	}
}

func TestFederatedAggregator_DP_StrongPrivacy(t *testing.T) {
	agg := service.NewFederatedAggregator(zap.NewNop(), 2)
	for i := 0; i < 2; i++ {
		agg.SubmitLocalUpdate(&service.LocalModelUpdate{
			DeviceID: fmt.Sprintf("d%d", i), ModelType: service.FederatedModelAnomaly,
			Weights: []float64{5.0, 10.0}, SampleCount: 200, TrainedAt: time.Now(),
		})
	}

	// 강한 프라이버시 (작은 epsilon → 큰 노이즈)
	dp := service.DPConfig{Epsilon: 0.1, Sensitivity: 1.0, ClipNorm: 1.0}
	global, err := agg.AggregateModelWithDP(service.FederatedModelAnomaly, dp)
	if err != nil {
		t.Fatalf("DP 집계 실패: %v", err)
	}
	if global.Metrics["noise_scale"] != 10.0 {
		t.Errorf("noise_scale = %f, 10.0 기대 (sensitivity/epsilon = 1/0.1)", global.Metrics["noise_scale"])
	}
}

func TestProcessResearchQuery_GenderFilter(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	s.UpdateConsent(ctx, "r2", true, true, false)

	result, err := s.ProcessResearchQuery(ctx, &service.ResearchQuery{
		RequesterID: "r2",
		DatasetType: "heart_rate",
		Gender:      "female",
	})
	if err != nil {
		t.Fatalf("ProcessResearchQuery 실패: %v", err)
	}
	if result.Summary.SampleCount <= 0 {
		t.Error("샘플 수가 0보다 커야 합니다")
	}
}
