package service_test

import (
	"math"
	"testing"

	"github.com/manpasik/backend/services/digital-twin-service/internal/service"
)

func TestStateMachine_HappyPath(t *testing.T) {
	sm := service.NewStateMachine()

	transitions := [][2]service.HealthState{
		{service.HealthNominal, service.HealthWarning},
		{service.HealthWarning, service.HealthDrift},
		{service.HealthDrift, service.HealthRecalibrate},
		{service.HealthRecalibrate, service.HealthNominal},
	}

	for _, tr := range transitions {
		if err := sm.Transition(tr[0], tr[1]); err != nil {
			t.Errorf("%s → %s 실패: %v", tr[0], tr[1], err)
		}
	}
}

func TestStateMachine_InvalidTransition(t *testing.T) {
	sm := service.NewStateMachine()
	// nominal → recalibrate 직접 불가
	if err := sm.Transition(service.HealthNominal, service.HealthRecalibrate); err == nil {
		t.Error("nominal → recalibrate 직접 전이가 허용됨")
	}
}

func TestStateMachine_CanTransition(t *testing.T) {
	sm := service.NewStateMachine()
	if !sm.CanTransition(service.HealthNominal, service.HealthWarning) {
		t.Error("nominal → warning 미허용")
	}
	if sm.CanTransition(service.HealthDrift, service.HealthWarning) {
		t.Error("drift → warning 잘못 허용")
	}
}

func TestStateMachine_FailedFromAnyState(t *testing.T) {
	sm := service.NewStateMachine()
	for _, from := range []service.HealthState{
		service.HealthNominal, service.HealthWarning,
		service.HealthDrift, service.HealthRecalibrate,
	} {
		if !sm.CanTransition(from, service.HealthFailed) {
			t.Errorf("%s → failed 미허용", from)
		}
	}
}

func TestStateMachine_AllowedNext(t *testing.T) {
	sm := service.NewStateMachine()
	next := sm.AllowedNext(service.HealthNominal)
	if len(next) == 0 {
		t.Error("nominal 다음 상태 없음")
	}
}

func TestAnalyzeMultiWindow_FullWindows(t *testing.T) {
	// 50개 측정 (1~50)
	values := make([]float64, 50)
	for i := 0; i < 50; i++ {
		values[i] = float64(i + 1)
	}

	analysis := service.AnalyzeMultiWindow(values)
	if analysis.Window5 == nil || analysis.Window5.Count != 5 {
		t.Errorf("Window5 = %v", analysis.Window5)
	}
	if analysis.Window10 == nil || analysis.Window10.Count != 10 {
		t.Errorf("Window10 = %v", analysis.Window10)
	}
	if analysis.Window30 == nil || analysis.Window30.Count != 30 {
		t.Errorf("Window30 = %v", analysis.Window30)
	}
	if analysis.OverallStats.Count != 50 {
		t.Errorf("Overall = %d", analysis.OverallStats.Count)
	}

	// 1~50의 OLS 기울기는 1.0
	if math.Abs(analysis.OverallStats.TrendSlope-1.0) > 0.01 {
		t.Errorf("TrendSlope = %f, want 1.0", analysis.OverallStats.TrendSlope)
	}
}

func TestAnalyzeMultiWindow_SmallSample(t *testing.T) {
	values := []float64{1, 2, 3}
	analysis := service.AnalyzeMultiWindow(values)
	// Window5/10/30은 모두 3개로 클램프
	if analysis.Window30.Count != 3 {
		t.Errorf("Window30 (small sample) = %d", analysis.Window30.Count)
	}
}

func TestAnalyzeMultiWindow_Empty(t *testing.T) {
	analysis := service.AnalyzeMultiWindow([]float64{})
	if analysis.OverallStats != nil && analysis.OverallStats.Count != 0 {
		t.Error("빈 입력에 통계가 생성됨")
	}
}

func TestWindowAnalysis_IsDrifting_True(t *testing.T) {
	// 첫 5개는 100, 다음 5개는 130 (30% 증가)
	values := []float64{100, 100, 100, 100, 100, 130, 130, 130, 130, 130}
	analysis := service.AnalyzeMultiWindow(values)
	if !analysis.IsDrifting(0.05) {
		t.Error("30% 변화인데 미감지")
	}
}

func TestWindowAnalysis_IsDrifting_False(t *testing.T) {
	values := []float64{100, 101, 99, 100, 102, 101, 100, 99, 101, 100}
	analysis := service.AnalyzeMultiWindow(values)
	if analysis.IsDrifting(0.05) {
		t.Error("정상 범위인데 드리프트 감지")
	}
}

func TestPredictNSteps_LinearExtrapolation(t *testing.T) {
	// 1~10 → 11, 12, 13 예측 기대
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	result := service.PredictNSteps(values, 3, service.HealthNominal)

	if len(result.Predictions) != 3 {
		t.Fatalf("Predictions = %d, want 3", len(result.Predictions))
	}
	if math.Abs(result.Predictions[0]-11.0) > 0.5 {
		t.Errorf("[0] = %f, want ~11", result.Predictions[0])
	}
	if math.Abs(result.Predictions[2]-13.0) > 0.5 {
		t.Errorf("[2] = %f, want ~13", result.Predictions[2])
	}

	// CI 폭은 잔차에 비례 (완벽한 선형이면 0)
	// 폭이 0이어도 high >= low여야 함
	for i := 0; i < 3; i++ {
		if result.ConfidenceHigh[i] < result.ConfidenceLow[i] {
			t.Errorf("CI[%d] 역전: high=%f low=%f", i, result.ConfidenceHigh[i], result.ConfidenceLow[i])
		}
	}
}

func TestPredictNSteps_NonLinearData_HasCIWidth(t *testing.T) {
	// 잔차가 있는 데이터 → CI 폭 양수
	values := []float64{1, 3, 2, 5, 4, 7, 6, 9, 8, 11}
	result := service.PredictNSteps(values, 3, service.HealthNominal)

	for i := 0; i < 3; i++ {
		width := result.ConfidenceHigh[i] - result.ConfidenceLow[i]
		if width <= 0 {
			t.Errorf("CI[%d] 폭 = %f, 잔차 있는 데이터인데 0", i, width)
		}
	}
}

func TestPredictNSteps_PredictedState_Stable(t *testing.T) {
	values := []float64{100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	result := service.PredictNSteps(values, 3, service.HealthNominal)
	if result.PredictedState != service.HealthNominal {
		t.Errorf("안정 데이터 PredictedState = %q, want nominal", result.PredictedState)
	}
}

func TestPredictNSteps_Drift_HighChange(t *testing.T) {
	// 급격히 증가 → 변화율 > 20% → recalibrate 예상
	values := []float64{100, 110, 130, 160, 200}
	result := service.PredictNSteps(values, 3, service.HealthWarning)
	if result.PredictedState != service.HealthRecalibrate &&
		result.PredictedState != service.HealthDrift {
		t.Errorf("PredictedState = %q, want drift/recalibrate", result.PredictedState)
	}
}

func TestPredictNSteps_TooFewSamples(t *testing.T) {
	result := service.PredictNSteps([]float64{5}, 3, service.HealthNominal)
	if len(result.Predictions) > 0 {
		t.Error("샘플 부족인데 예측 생성")
	}
}

func TestRecommendCalibration_RecalibrateState(t *testing.T) {
	state := &service.TwinState{HealthState: string(service.HealthRecalibrate)}
	rec := service.RecommendCalibration(state, nil)
	if !rec.Required {
		t.Error("recalibrate 상태인데 Required=false")
	}
	if rec.Urgency != "immediate" {
		t.Errorf("Urgency = %q", rec.Urgency)
	}
}

func TestRecommendCalibration_DriftScheduled(t *testing.T) {
	state := &service.TwinState{HealthState: string(service.HealthDrift)}
	rec := service.RecommendCalibration(state, nil)
	if !rec.Required {
		t.Error("drift 상태인데 Required=false")
	}
	if rec.Urgency != "scheduled" {
		t.Errorf("Urgency = %q", rec.Urgency)
	}
}

func TestRecommendCalibration_NominalNoAction(t *testing.T) {
	state := &service.TwinState{HealthState: string(service.HealthNominal)}
	rec := service.RecommendCalibration(state, nil)
	if rec.Required {
		t.Error("nominal 상태인데 Required=true")
	}
}

func TestRecommendCalibration_NilState(t *testing.T) {
	if rec := service.RecommendCalibration(nil, nil); rec != nil {
		t.Error("nil state에 권고 반환")
	}
}

func TestRecommendCalibration_DriftRateUpgradesUrgency(t *testing.T) {
	state := &service.TwinState{HealthState: string(service.HealthNominal)}

	// 큰 드리프트율 시뮬레이션 (Window5 = 20% 증가)
	values := []float64{100, 100, 100, 100, 100, 120, 120, 120, 120, 120}
	analysis := service.AnalyzeMultiWindow(values)

	rec := service.RecommendCalibration(state, analysis)
	if rec.Urgency != "scheduled" {
		t.Errorf("Urgency = %q, want scheduled (drift_rate 15%% 초과)", rec.Urgency)
	}
}

func TestSortedCopy(t *testing.T) {
	in := []float64{3, 1, 2}
	out := service.SortedCopy(in)
	if out[0] != 1 || out[1] != 2 || out[2] != 3 {
		t.Errorf("정렬 실패: %v", out)
	}
	// 원본 미변경
	if in[0] != 3 {
		t.Error("원본 변경됨")
	}
}

func TestStateMachine_FailedRecovery(t *testing.T) {
	sm := service.NewStateMachine()
	if !sm.CanTransition(service.HealthFailed, service.HealthNominal) {
		t.Error("failed → nominal 복구 미허용")
	}
}

func TestPredictNSteps_NegativeN(t *testing.T) {
	result := service.PredictNSteps([]float64{1, 2, 3}, -1, service.HealthNominal)
	if len(result.Predictions) > 0 {
		t.Error("음수 n에 예측 생성")
	}
}
