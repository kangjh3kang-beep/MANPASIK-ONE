package service_test

import (
	"testing"

	"github.com/manpasik/backend/services/ai-inference-service/internal/service"
)

// TestModelRegistry_Register는 모델 등록을 검증합니다.
func TestModelRegistry_Register(t *testing.T) {
	r := service.NewModelRegistry()

	m := &service.ModelMetadata{
		Name:      "biomarker-clf",
		Version:   "v1.0.0",
		Type:      service.ModelBiomarkerClassifier,
		Framework: "tflite",
		Checksum:  "abc123",
	}

	if err := r.Register(m); err != nil {
		t.Fatalf("Register 실패: %v", err)
	}

	got, err := r.Get("biomarker-clf", "v1.0.0")
	if err != nil {
		t.Fatalf("Get 실패: %v", err)
	}
	if got.Status != service.RegistryStatusDraft {
		t.Errorf("Status = %q, want draft", got.Status)
	}
}

// TestModelRegistry_RegisterDuplicate는 중복 등록 거부를 검증합니다.
func TestModelRegistry_RegisterDuplicate(t *testing.T) {
	r := service.NewModelRegistry()
	m := &service.ModelMetadata{Name: "x", Version: "v1"}
	_ = r.Register(m)
	if err := r.Register(m); err == nil {
		t.Error("중복 등록이 허용됨")
	}
}

// TestModelRegistry_PromoteHappyPath는 정상 승격 흐름을 검증합니다.
func TestModelRegistry_PromoteHappyPath(t *testing.T) {
	r := service.NewModelRegistry()
	_ = r.Register(&service.ModelMetadata{Name: "m", Version: "v1"})

	transitions := []string{
		service.RegistryStatusStaging,
		service.RegistryStatusCanary,
		service.RegistryStatusProduction,
	}
	for _, s := range transitions {
		if err := r.Promote("m", "v1", s); err != nil {
			t.Fatalf("Promote(%s) 실패: %v", s, err)
		}
	}

	got, _ := r.Get("m", "v1")
	if got.Status != service.RegistryStatusProduction {
		t.Errorf("Status = %q, want production", got.Status)
	}
	if got.PromotedAt == nil {
		t.Error("PromotedAt 미설정")
	}
}

// TestModelRegistry_PromoteAutoDeprecation은 새 모델 production 승격 시
// 기존 production 모델이 자동 deprecated되는지 검증합니다.
func TestModelRegistry_PromoteAutoDeprecation(t *testing.T) {
	r := service.NewModelRegistry()
	_ = r.Register(&service.ModelMetadata{Name: "m", Version: "v1", Status: service.RegistryStatusProduction})
	_ = r.Register(&service.ModelMetadata{Name: "m", Version: "v2", Status: service.RegistryStatusStaging})

	if err := r.Promote("m", "v2", service.RegistryStatusProduction); err != nil {
		t.Fatalf("Promote 실패: %v", err)
	}

	old, _ := r.Get("m", "v1")
	if old.Status != service.RegistryStatusDeprecated {
		t.Errorf("v1 Status = %q, want deprecated", old.Status)
	}

	new_, _ := r.Get("m", "v2")
	if new_.Status != service.RegistryStatusProduction {
		t.Errorf("v2 Status = %q, want production", new_.Status)
	}
}

// TestModelRegistry_PromoteInvalidTransition는 잘못된 승격을 거부합니다.
func TestModelRegistry_PromoteInvalidTransition(t *testing.T) {
	r := service.NewModelRegistry()
	_ = r.Register(&service.ModelMetadata{Name: "m", Version: "v1"})

	// draft → production 직접 승격은 불가
	err := r.Promote("m", "v1", service.RegistryStatusProduction)
	if err == nil {
		t.Error("draft → production 직접 승격이 허용됨")
	}
}

// TestModelRegistry_GetActive는 운영 중인 모델 조회를 검증합니다.
func TestModelRegistry_GetActive(t *testing.T) {
	r := service.NewModelRegistry()
	_ = r.Register(&service.ModelMetadata{Name: "m", Version: "v1", Status: service.RegistryStatusProduction})

	active, err := r.GetActive("m")
	if err != nil {
		t.Fatalf("GetActive 실패: %v", err)
	}
	if active.Version != "v1" {
		t.Errorf("Version = %q, want v1", active.Version)
	}
}

// TestModelRegistry_SetCanaryRatio는 카나리 비율 설정을 검증합니다.
func TestModelRegistry_SetCanaryRatio(t *testing.T) {
	r := service.NewModelRegistry()
	_ = r.Register(&service.ModelMetadata{Name: "m", Version: "v1", Status: service.RegistryStatusCanary})

	if err := r.SetCanaryRatio("m", "v1", 0.25); err != nil {
		t.Fatalf("SetCanaryRatio 실패: %v", err)
	}

	got, _ := r.Get("m", "v1")
	if got.CanaryRatio != 0.25 {
		t.Errorf("CanaryRatio = %f, want 0.25", got.CanaryRatio)
	}

	// 비율 범위 검증
	if err := r.SetCanaryRatio("m", "v1", 1.5); err == nil {
		t.Error("1.0 초과 비율이 허용됨")
	}
}

// TestModelRegistry_UpdateMetrics는 평가 지표 갱신을 검증합니다.
func TestModelRegistry_UpdateMetrics(t *testing.T) {
	r := service.NewModelRegistry()
	_ = r.Register(&service.ModelMetadata{Name: "m", Version: "v1"})

	metrics := &service.ModelMetric{
		Accuracy:  0.92,
		Precision: 0.91,
		Recall:    0.89,
		F1:        0.90,
		AUC:       0.95,
	}

	if err := r.UpdateMetrics("m", "v1", metrics); err != nil {
		t.Fatalf("UpdateMetrics 실패: %v", err)
	}

	got, _ := r.Get("m", "v1")
	if got.Metrics.Accuracy != 0.92 {
		t.Errorf("Accuracy = %f, want 0.92", got.Metrics.Accuracy)
	}
	if got.Metrics.EvaluatedAt.IsZero() {
		t.Error("EvaluatedAt 미설정")
	}
}

// TestConfidenceGate_Pass는 신뢰도 게이팅을 검증합니다.
func TestConfidenceGate_Pass(t *testing.T) {
	gate := service.NewConfidenceGate(0.8)

	if !gate.Pass(0.9) {
		t.Error("0.9가 0.8 임계값을 통과하지 못함")
	}
	if gate.Pass(0.5) {
		t.Error("0.5가 0.8 임계값을 통과함")
	}

	pass, reason := gate.PassWithReason(0.5)
	if pass {
		t.Error("PassWithReason: 0.5가 통과로 판정됨")
	}
	if reason == "" {
		t.Error("실패 이유가 비어 있음")
	}
}

// TestXAIExplainer_Explain은 특성 중요도 계산을 검증합니다.
func TestXAIExplainer_Explain(t *testing.T) {
	e := service.NewXAIExplainer()

	features := map[string]float64{
		"glucose":  120.0,
		"bp_sys":   130.0,
		"bp_dia":   85.0,
		"weight":   70.0,
	}
	weights := map[string]float64{
		"glucose": 0.5,
		"bp_sys":  0.3,
		"bp_dia":  0.1,
		"weight":  -0.05,
	}

	importances := e.Explain(features, weights)
	if len(importances) != 4 {
		t.Fatalf("importance 수 = %d, want 4", len(importances))
	}

	// 중요도가 내림차순 정렬되어야 함
	for i := 1; i < len(importances); i++ {
		if importances[i].Importance > importances[i-1].Importance {
			t.Errorf("정렬 오류: [%d]=%f > [%d]=%f",
				i, importances[i].Importance, i-1, importances[i-1].Importance)
		}
	}

	// 합이 1.0이어야 함 (정규화)
	sum := 0.0
	for _, imp := range importances {
		sum += imp.Importance
	}
	if sum < 0.99 || sum > 1.01 {
		t.Errorf("정규화 합 = %f, want 1.0", sum)
	}
}

// TestXAIExplainer_TopK는 상위 K 추출을 검증합니다.
func TestXAIExplainer_TopK(t *testing.T) {
	e := service.NewXAIExplainer()
	importances := []service.FeatureImportance{
		{Feature: "a", Importance: 0.5},
		{Feature: "b", Importance: 0.3},
		{Feature: "c", Importance: 0.2},
	}

	top2 := e.TopK(importances, 2)
	if len(top2) != 2 {
		t.Errorf("top2 수 = %d, want 2", len(top2))
	}
}

// TestDriftDetector_NoDrift는 드리프트 없는 분포를 검증합니다.
func TestDriftDetector_NoDrift(t *testing.T) {
	d := service.NewDriftDetector(0.3)

	// 동일한 분포 (같은 평균/표준편차)
	baseline := []float64{100, 102, 99, 101, 100, 103, 98}
	current := []float64{100, 102, 99, 101, 100, 103, 98}

	report := d.Compare("glucose", baseline, current)
	if report == nil {
		t.Fatal("report nil")
	}
	if report.Status != service.DriftNone {
		t.Errorf("Status = %q, want none (KS=%f)", report.Status, report.KSStatistic)
	}
}

// TestDriftDetector_Significant는 큰 드리프트를 감지합니다.
func TestDriftDetector_Significant(t *testing.T) {
	d := service.NewDriftDetector(0.1)

	baseline := []float64{100, 102, 99, 101, 100, 103, 98}
	current := []float64{200, 205, 198, 210, 195, 205, 200}

	report := d.Compare("glucose", baseline, current)
	if report.Status != service.DriftSignificant {
		t.Errorf("Status = %q, want significant", report.Status)
	}
	if report.KSStatistic <= 0.2 {
		t.Errorf("KSStatistic = %f, want > 0.2", report.KSStatistic)
	}
}

// TestDriftDetector_Batch는 일괄 드리프트 분석을 검증합니다.
func TestDriftDetector_Batch(t *testing.T) {
	d := service.NewDriftDetector(0.1)

	baseline := map[string][]float64{
		"glucose": {100, 102, 99, 101},
		"bp":      {120, 122, 119, 121},
	}
	current := map[string][]float64{
		"glucose": {200, 205, 198, 210},
		"bp":      {121, 120, 119, 122},
	}

	reports := d.CompareBatch(baseline, current)
	if len(reports) != 2 {
		t.Fatalf("reports = %d, want 2", len(reports))
	}

	if !service.HasSignificantDrift(reports) {
		t.Error("significant drift 미감지")
	}

	// 첫 번째가 더 큰 드리프트 (정렬 검증)
	if reports[0].KSStatistic < reports[1].KSStatistic {
		t.Error("드리프트 정렬이 내림차순이 아님")
	}
}

// TestModelRegistry_CountByStatus는 상태별 통계를 검증합니다.
func TestModelRegistry_CountByStatus(t *testing.T) {
	r := service.NewModelRegistry()
	_ = r.Register(&service.ModelMetadata{Name: "a", Version: "v1"})
	_ = r.Register(&service.ModelMetadata{Name: "b", Version: "v1", Status: service.RegistryStatusProduction})
	_ = r.Register(&service.ModelMetadata{Name: "c", Version: "v1", Status: service.RegistryStatusProduction})

	counts := r.CountByStatus()
	if counts[service.RegistryStatusDraft] != 1 {
		t.Errorf("draft = %d, want 1", counts[service.RegistryStatusDraft])
	}
	if counts[service.RegistryStatusProduction] != 2 {
		t.Errorf("production = %d, want 2", counts[service.RegistryStatusProduction])
	}
}

// TestModelRegistry_ListByName은 모델명별 버전 목록 조회를 검증합니다.
func TestModelRegistry_ListByName(t *testing.T) {
	r := service.NewModelRegistry()
	_ = r.Register(&service.ModelMetadata{Name: "m", Version: "v1"})
	_ = r.Register(&service.ModelMetadata{Name: "m", Version: "v2"})
	_ = r.Register(&service.ModelMetadata{Name: "m", Version: "v3"})
	_ = r.Register(&service.ModelMetadata{Name: "other", Version: "v1"})

	versions := r.ListByName("m")
	if len(versions) != 3 {
		t.Errorf("versions = %d, want 3", len(versions))
	}
	// 내림차순 정렬 확인
	if versions[0].Version != "v3" {
		t.Errorf("first version = %q, want v3", versions[0].Version)
	}
}
