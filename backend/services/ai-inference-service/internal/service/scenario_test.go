package service_test

import (
	"testing"

	"github.com/manpasik/backend/services/ai-inference-service/internal/service"
)

// TestScenario_FullModelLifecycle은 모델 등록 → 검증 → 운영 → 폐기 흐름입니다.
func TestScenario_FullModelLifecycle(t *testing.T) {
	r := service.NewModelRegistry()

	// 1. 등록 (draft)
	v1 := &service.ModelMetadata{
		Name:      "biomarker-v1",
		Version:   "1.0.0",
		Type:      service.ModelBiomarkerClassifier,
		Framework: "tflite",
		Checksum:  "sha256:abc",
	}
	if err := r.Register(v1); err != nil {
		t.Fatalf("Register 실패: %v", err)
	}

	// 2. 메트릭 갱신 후 staging 승격
	metrics := &service.ModelMetric{
		Accuracy: 0.92, Precision: 0.91, Recall: 0.89, F1: 0.90, AUC: 0.95,
	}
	_ = r.UpdateMetrics("biomarker-v1", "1.0.0", metrics)
	if err := r.Promote("biomarker-v1", "1.0.0", service.RegistryStatusStaging); err != nil {
		t.Fatalf("staging 승격 실패: %v", err)
	}

	// 3. canary → production
	_ = r.Promote("biomarker-v1", "1.0.0", service.RegistryStatusCanary)
	_ = r.SetCanaryRatio("biomarker-v1", "1.0.0", 0.1)
	_ = r.Promote("biomarker-v1", "1.0.0", service.RegistryStatusProduction)

	// 4. v2 등록 → production 시 v1 자동 deprecated
	v2 := &service.ModelMetadata{
		Name: "biomarker-v1", Version: "2.0.0",
		Status: service.RegistryStatusStaging,
	}
	_ = r.Register(v2)
	_ = r.Promote("biomarker-v1", "2.0.0", service.RegistryStatusProduction)

	v1Updated, _ := r.Get("biomarker-v1", "1.0.0")
	if v1Updated.Status != service.RegistryStatusDeprecated {
		t.Errorf("v1 status = %q, want deprecated", v1Updated.Status)
	}

	// 5. 활성 모델 확인
	active, err := r.GetActive("biomarker-v1")
	if err != nil {
		t.Fatalf("GetActive 실패: %v", err)
	}
	if active.Version != "2.0.0" {
		t.Errorf("active version = %q, want 2.0.0", active.Version)
	}
}

// TestScenario_XAIWithRealisticFeatures는 실제 바이오마커 특성으로 XAI를 검증합니다.
func TestScenario_XAIWithRealisticFeatures(t *testing.T) {
	e := service.NewXAIExplainer()

	features := map[string]float64{
		"glucose":         180.0, // 높음 (정상 70-100)
		"hba1c":           7.5,    // 높음 (정상 <5.7)
		"systolic_bp":     145.0, // 높음 (정상 <120)
		"diastolic_bp":    92.0,  // 높음 (정상 <80)
		"weight_kg":       80.0,
	}
	weights := map[string]float64{
		"glucose":         0.4,
		"hba1c":           0.3,
		"systolic_bp":     0.15,
		"diastolic_bp":    0.10,
		"weight_kg":       0.05,
	}

	importances := e.Explain(features, weights)
	if len(importances) != 5 {
		t.Fatalf("importances = %d, want 5", len(importances))
	}

	// glucose가 가장 큰 영향
	if importances[0].Feature != "glucose" {
		t.Errorf("top feature = %q, want glucose", importances[0].Feature)
	}

	// 모두 positive 방향
	for _, imp := range importances {
		if imp.Direction != "positive" {
			t.Errorf("%s direction = %q, want positive", imp.Feature, imp.Direction)
		}
	}

	// Top 3
	top3 := e.TopK(importances, 3)
	if len(top3) != 3 {
		t.Errorf("top3 = %d, want 3", len(top3))
	}
}

// TestScenario_DriftMonitoringPipeline은 운영 모델 드리프트 모니터링을 검증합니다.
func TestScenario_DriftMonitoringPipeline(t *testing.T) {
	detector := service.NewDriftDetector(0.15)

	// 정상 운영 분포 (혈당 + 혈압)
	baselineGlucose := []float64{95, 98, 102, 100, 99, 105, 97, 101}
	baselineBP := []float64{120, 122, 118, 121, 119, 123, 120}

	// 가설 시나리오 1: 정상 (드리프트 없음)
	currentGlucoseNormal := []float64{96, 99, 101, 100, 98, 104, 99, 102}
	currentBPNormal := []float64{121, 120, 119, 122, 120, 121, 120}

	r1 := detector.Compare("glucose", baselineGlucose, currentGlucoseNormal)
	r2 := detector.Compare("bp", baselineBP, currentBPNormal)

	if r1.Status == service.DriftSignificant {
		t.Errorf("정상 분포에 significant drift 감지: ks=%f", r1.KSStatistic)
	}
	if r2.Status == service.DriftSignificant {
		t.Errorf("정상 분포에 significant drift 감지: ks=%f", r2.KSStatistic)
	}

	// 가설 시나리오 2: 큰 드리프트 (분석 모델 갱신 필요)
	currentGlucoseShifted := []float64{180, 185, 175, 190, 178, 195, 188}
	r3 := detector.Compare("glucose", baselineGlucose, currentGlucoseShifted)
	if r3.Status != service.DriftSignificant {
		t.Errorf("드리프트 미감지: ks=%f", r3.KSStatistic)
	}
}

// TestScenario_ConfidenceGateGating은 신뢰도 게이팅 운영 시나리오입니다.
func TestScenario_ConfidenceGateGating(t *testing.T) {
	// 신뢰도 0.85 미만은 차단 (의료 도메인은 높은 신뢰도 필요)
	gate := service.NewConfidenceGate(0.85)

	// 정상 케이스: 0.92 신뢰도
	if !gate.Pass(0.92) {
		t.Error("0.92 신뢰도가 통과되지 않음")
	}

	// 차단 케이스: 0.7 신뢰도
	pass, reason := gate.PassWithReason(0.7)
	if pass {
		t.Error("0.7 신뢰도가 통과됨")
	}
	if reason == "" {
		t.Error("실패 이유 미기록")
	}
}

// TestScenario_BatchDriftAnalysis는 다중 특성 일괄 드리프트 분석입니다.
func TestScenario_BatchDriftAnalysis(t *testing.T) {
	detector := service.NewDriftDetector(0.15)

	baseline := map[string][]float64{
		"glucose":  {95, 98, 102, 100, 99},
		"bp_sys":   {120, 122, 118, 121, 119},
		"hr":       {72, 75, 70, 74, 71},
	}
	current := map[string][]float64{
		"glucose":  {180, 185, 175, 190, 178}, // 큰 드리프트
		"bp_sys":   {121, 120, 119, 122, 120}, // 정상
		"hr":       {73, 76, 71, 75, 72},      // 정상
	}

	reports := detector.CompareBatch(baseline, current)
	if len(reports) != 3 {
		t.Fatalf("reports = %d, want 3", len(reports))
	}

	if !service.HasSignificantDrift(reports) {
		t.Error("significant drift 미감지")
	}

	// 첫 번째가 가장 큰 드리프트 (glucose)
	if reports[0].Feature != "glucose" {
		t.Errorf("first feature = %q, want glucose", reports[0].Feature)
	}
}

// TestScenario_ModelVersionConcurrency는 동일 모델명의 여러 버전 관리를 검증합니다.
func TestScenario_ModelVersionConcurrency(t *testing.T) {
	r := service.NewModelRegistry()

	versions := []string{"1.0.0", "1.0.1", "1.1.0", "2.0.0"}
	for _, v := range versions {
		_ = r.Register(&service.ModelMetadata{
			Name: "shared-model", Version: v,
			Status: service.RegistryStatusStaging,
		})
	}

	all := r.ListByName("shared-model")
	if len(all) != 4 {
		t.Errorf("ListByName = %d, want 4", len(all))
	}

	// 내림차순 정렬 확인
	if all[0].Version != "2.0.0" {
		t.Errorf("first = %q, want 2.0.0", all[0].Version)
	}
}
