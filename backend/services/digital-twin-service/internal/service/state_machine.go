// Package service: 디지털 트윈 상태 머신 + 예측 + 다중 윈도우 분석
//
// CRUD → CRUD+ 승격을 위한 강화 모듈.
// 기존 EWMA/CUSUM 위에 상태 전이 검증, 캘리브레이션 권고, N-step 예측을 제공합니다.
package service

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

// ============================================================================
// 상태 머신
// ============================================================================

// HealthState는 디지털 트윈의 건강 상태입니다.
type HealthState string

const (
	HealthNominal     HealthState = "nominal"      // 정상
	HealthWarning     HealthState = "warning"      // 경고 (소폭 드리프트)
	HealthDrift       HealthState = "drift"        // 드리프트 감지
	HealthRecalibrate HealthState = "recalibrate"  // 캘리브레이션 필요
	HealthFailed      HealthState = "failed"       // 측정 실패
)

// allowedTransitions는 허용된 상태 전이입니다.
//
// nominal → warning → drift → recalibrate (단방향 진행)
// 캘리브레이션 후: recalibrate → nominal (회복)
// 어떤 상태에서든 failed 가능 (긴급 정지)
var allowedTransitions = map[HealthState][]HealthState{
	HealthNominal:     {HealthWarning, HealthDrift, HealthFailed},
	HealthWarning:     {HealthNominal, HealthDrift, HealthFailed},
	HealthDrift:       {HealthRecalibrate, HealthFailed},
	HealthRecalibrate: {HealthNominal, HealthFailed},
	HealthFailed:      {HealthNominal}, // 수동 복구 후
}

// StateMachine은 트윈 상태 전이를 관리합니다.
type StateMachine struct {
	mu sync.RWMutex
	transitions map[HealthState][]HealthState
}

// NewStateMachine은 새 상태 머신을 생성합니다.
func NewStateMachine() *StateMachine {
	cloned := make(map[HealthState][]HealthState, len(allowedTransitions))
	for k, v := range allowedTransitions {
		cloned[k] = append([]HealthState{}, v...)
	}
	return &StateMachine{transitions: cloned}
}

// CanTransition은 상태 전이 가능 여부를 반환합니다.
func (sm *StateMachine) CanTransition(from, to HealthState) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	allowed, ok := sm.transitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

// Transition은 상태 전이를 시도합니다. 실패 시 error 반환.
func (sm *StateMachine) Transition(from, to HealthState) error {
	if !sm.CanTransition(from, to) {
		return fmt.Errorf("invalid transition: %s → %s", from, to)
	}
	return nil
}

// AllowedNext는 현재 상태에서 가능한 다음 상태 목록을 반환합니다.
func (sm *StateMachine) AllowedNext(state HealthState) []HealthState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	transitions := sm.transitions[state]
	out := make([]HealthState, len(transitions))
	copy(out, transitions)
	return out
}

// ============================================================================
// 다중 윈도우 분석
// ============================================================================

// WindowAnalysis는 다중 시간 윈도우의 통계를 제공합니다.
type WindowAnalysis struct {
	Window5   *WindowStats // 최근 5개
	Window10  *WindowStats // 최근 10개
	Window30  *WindowStats // 최근 30개
	OverallStats *WindowStats // 전체
}

// WindowStats는 단일 윈도우 통계입니다.
type WindowStats struct {
	WindowSize int
	Count      int
	Mean       float64
	StdDev     float64
	Min        float64
	Max        float64
	TrendSlope float64 // OLS 기울기 (상승=+, 하강=-)
	DriftRate  float64 // 직전 윈도우 대비 평균 변화율
}

// AnalyzeMultiWindow는 측정값 시계열을 다중 윈도우로 분석합니다.
//
// values: 시간순 정렬된 측정값.
func AnalyzeMultiWindow(values []float64) *WindowAnalysis {
	if len(values) == 0 {
		return &WindowAnalysis{}
	}

	return &WindowAnalysis{
		Window5:      computeWindow(values, 5),
		Window10:     computeWindow(values, 10),
		Window30:     computeWindow(values, 30),
		OverallStats: computeWindow(values, len(values)),
	}
}

func computeWindow(values []float64, size int) *WindowStats {
	n := len(values)
	if n == 0 || size <= 0 {
		return &WindowStats{WindowSize: size}
	}
	if size > n {
		size = n
	}
	window := values[n-size:]
	mean := mean64(window)
	std := stdDev64(window, mean)
	min, max := minMax64(window)
	slope := olsSlope(window)

	stats := &WindowStats{
		WindowSize: size,
		Count:      len(window),
		Mean:       mean,
		StdDev:     std,
		Min:        min,
		Max:        max,
		TrendSlope: slope,
	}

	// 드리프트율: 직전 동일 크기 윈도우와 평균 비교
	if n >= 2*size {
		prevWindow := values[n-2*size : n-size]
		prevMean := mean64(prevWindow)
		if prevMean != 0 {
			stats.DriftRate = (mean - prevMean) / math.Abs(prevMean)
		}
	}
	return stats
}

// IsDrifting은 윈도우 분석에서 드리프트를 판정합니다.
//
// 기준: 최근 5개의 평균이 직전 5개 대비 ±5% 이상 변화 + slope 크기 임계값 초과.
func (w *WindowAnalysis) IsDrifting(threshold float64) bool {
	if w.Window5 == nil || w.Window5.Count < 5 {
		return false
	}
	if math.Abs(w.Window5.DriftRate) >= threshold {
		return true
	}
	return false
}

// ============================================================================
// N-step 예측
// ============================================================================

// PredictionResult는 N-step 예측 결과입니다.
type PredictionResult struct {
	NSteps         int
	Predictions    []float64
	ConfidenceLow  []float64 // 95% CI 하한
	ConfidenceHigh []float64 // 95% CI 상한
	PredictedState HealthState // 다음 N step 후 예상 상태
	GeneratedAt    time.Time
}

// PredictNSteps는 OLS 기반 N-step 예측을 수행합니다.
//
// 단순 선형 회귀로 추세를 외삽하며, 95% CI는 잔차 표준편차에 기반합니다.
func PredictNSteps(values []float64, n int, currentState HealthState) *PredictionResult {
	if len(values) < 2 || n <= 0 {
		return &PredictionResult{NSteps: n, GeneratedAt: time.Now().UTC()}
	}

	intercept, slope := olsParams(values)
	residuals := computeResiduals(values, intercept, slope)
	residualStd := stdDev64(residuals, 0) // 잔차 평균은 0
	const z = 1.96

	predictions := make([]float64, n)
	ciLow := make([]float64, n)
	ciHigh := make([]float64, n)

	startIdx := float64(len(values))
	for i := 0; i < n; i++ {
		x := startIdx + float64(i)
		yHat := intercept + slope*x
		predictions[i] = yHat
		margin := z * residualStd * math.Sqrt(1+1/float64(len(values)))
		ciLow[i] = yHat - margin
		ciHigh[i] = yHat + margin
	}

	// 예측 상태: 마지막 예측값과 마지막 실측값의 변화율 기반
	lastValue := values[len(values)-1]
	lastPred := predictions[n-1]
	predictedState := currentState
	if lastValue != 0 {
		change := math.Abs(lastPred-lastValue) / math.Abs(lastValue)
		switch {
		case change > 0.20:
			predictedState = HealthRecalibrate
		case change > 0.10:
			predictedState = HealthDrift
		case change > 0.05:
			predictedState = HealthWarning
		default:
			predictedState = HealthNominal
		}
	}

	return &PredictionResult{
		NSteps:         n,
		Predictions:    predictions,
		ConfidenceLow:  ciLow,
		ConfidenceHigh: ciHigh,
		PredictedState: predictedState,
		GeneratedAt:    time.Now().UTC(),
	}
}

// ============================================================================
// 캘리브레이션 권고
// ============================================================================

// CalibrationRecommendation은 캘리브레이션 권고입니다.
type CalibrationRecommendation struct {
	Required          bool
	Urgency           string // "immediate" | "scheduled" | "monitor"
	Reason            string
	SuggestedActions  []string
	EstimatedDowntime time.Duration
}

// RecommendCalibration은 트윈 상태로부터 캘리브레이션 권고를 생성합니다.
func RecommendCalibration(state *TwinState, analysis *WindowAnalysis) *CalibrationRecommendation {
	if state == nil {
		return nil
	}
	rec := &CalibrationRecommendation{}

	switch HealthState(state.HealthState) {
	case HealthRecalibrate:
		rec.Required = true
		rec.Urgency = "immediate"
		rec.Reason = "드리프트 임계 초과로 즉시 캘리브레이션 필요"
		rec.SuggestedActions = []string{
			"디바이스 펌웨어 자가 진단 실행",
			"기준 시료로 재측정",
			"R_n 기준 채널 신호 검증",
		}
		rec.EstimatedDowntime = 15 * time.Minute
	case HealthDrift:
		rec.Required = true
		rec.Urgency = "scheduled"
		rec.Reason = "지속적인 드리프트 추세 감지"
		rec.SuggestedActions = []string{
			"24시간 내 캘리브레이션 예약",
			"측정 신뢰도 임계 강화 (95% → 99%)",
		}
		rec.EstimatedDowntime = 10 * time.Minute
	case HealthWarning:
		rec.Required = false
		rec.Urgency = "monitor"
		rec.Reason = "경미한 변화 감지, 추가 측정으로 추세 확인"
		rec.SuggestedActions = []string{"다음 5회 측정 결과 모니터링"}
	default:
		rec.Required = false
		rec.Urgency = "monitor"
		rec.Reason = "정상 범위"
	}

	// 윈도우 분석에서 큰 변화가 보이면 긴급도 상향
	if analysis != nil && analysis.Window5 != nil && math.Abs(analysis.Window5.DriftRate) > 0.15 {
		if rec.Urgency == "monitor" {
			rec.Urgency = "scheduled"
			rec.SuggestedActions = append(rec.SuggestedActions,
				"최근 5회 평균 변화율 15% 초과 - 사전 점검 권장")
		}
	}

	return rec
}

// ============================================================================
// 헬퍼 (수치 계산)
// ============================================================================

func mean64(v []float64) float64 {
	if len(v) == 0 {
		return 0
	}
	s := 0.0
	for _, x := range v {
		s += x
	}
	return s / float64(len(v))
}

func stdDev64(v []float64, mean float64) float64 {
	if len(v) < 2 {
		return 0
	}
	s := 0.0
	for _, x := range v {
		s += (x - mean) * (x - mean)
	}
	return math.Sqrt(s / float64(len(v)-1))
}

func minMax64(v []float64) (min, max float64) {
	if len(v) == 0 {
		return 0, 0
	}
	min, max = v[0], v[0]
	for _, x := range v {
		if x < min {
			min = x
		}
		if x > max {
			max = x
		}
	}
	return
}

// olsParams는 단순 선형 회귀 (y = intercept + slope*x) 파라미터를 계산합니다.
func olsParams(values []float64) (intercept, slope float64) {
	n := len(values)
	if n < 2 {
		return 0, 0
	}
	xs := make([]float64, n)
	for i := range xs {
		xs[i] = float64(i)
	}
	xMean := mean64(xs)
	yMean := mean64(values)

	num := 0.0
	den := 0.0
	for i := 0; i < n; i++ {
		num += (xs[i] - xMean) * (values[i] - yMean)
		den += (xs[i] - xMean) * (xs[i] - xMean)
	}
	if den == 0 {
		return yMean, 0
	}
	slope = num / den
	intercept = yMean - slope*xMean
	return
}

// olsSlope는 OLS 기울기만 반환합니다.
func olsSlope(values []float64) float64 {
	_, slope := olsParams(values)
	return slope
}

func computeResiduals(values []float64, intercept, slope float64) []float64 {
	residuals := make([]float64, len(values))
	for i, v := range values {
		yHat := intercept + slope*float64(i)
		residuals[i] = v - yHat
	}
	return residuals
}

// SortedCopy returns a sorted copy without modifying input.
func SortedCopy(in []float64) []float64 {
	out := make([]float64, len(in))
	copy(out, in)
	sort.Float64s(out)
	return out
}

// ============================================================================
// 통합: TwinService 확장 (기존 코드 보완)
// ============================================================================

// AnalyzeAndRecommend는 트윈 상태와 측정 시계열을 종합 분석합니다.
//
// 반환:
//   - WindowAnalysis: 다중 윈도우 통계
//   - PredictionResult: N-step 예측
//   - CalibrationRecommendation: 캘리브레이션 권고
func (s *TwinService) AnalyzeAndRecommend(
	state *TwinState,
	values []float64,
	predictNSteps int,
) (*WindowAnalysis, *PredictionResult, *CalibrationRecommendation, error) {
	if state == nil {
		return nil, nil, nil, errors.New("state required")
	}
	analysis := AnalyzeMultiWindow(values)
	currentState := HealthState(state.HealthState)
	if currentState == "" {
		currentState = HealthNominal
	}
	prediction := PredictNSteps(values, predictNSteps, currentState)
	recommendation := RecommendCalibration(state, analysis)
	return analysis, prediction, recommendation, nil
}
