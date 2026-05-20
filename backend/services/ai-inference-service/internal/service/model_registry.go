// Package service: AI 모델 레지스트리 + XAI + Drift Detection
//
// MLOps에서 모델 버전 관리, 카나리 배포, 입력 분포 변화 감지(KS-test 단순화),
// SHAP 유사 특성 중요도를 제공합니다.
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
// 모델 메타데이터
// ============================================================================

// ModelStatus는 모델 상태입니다.
const (
	RegistryStatusDraft      = "draft"      // 개발 중
	RegistryStatusStaging    = "staging"    // 스테이징
	RegistryStatusCanary     = "canary"     // 카나리 배포
	RegistryStatusProduction = "production" // 운영 활성
	RegistryStatusDeprecated = "deprecated" // 폐기 예정
	RegistryStatusArchived   = "archived"   // 보관됨
)

// ModelMetric은 모델 평가 지표입니다.
type ModelMetric struct {
	Accuracy   float64
	Precision  float64
	Recall     float64
	F1         float64
	AUC        float64
	LossValue  float64
	EvaluatedAt time.Time
}

// ModelMetadata는 모델의 메타데이터입니다.
type ModelMetadata struct {
	Name          string
	Version       string
	Type          AiModelType
	Status        string
	Checksum      string // SHA-256
	Framework     string // "tflite", "onnx", "pytorch", "keras"
	SizeMB        float64
	InputSchema   map[string]string // field -> dtype
	OutputSchema  map[string]string
	Metrics       *ModelMetric
	CanaryRatio   float64 // 0.0~1.0, canary 트래픽 비율
	CreatedAt     time.Time
	PromotedAt    *time.Time
	DeprecatedAt  *time.Time
	Tags          []string
	Description   string
}

// ModelKey는 모델 식별 키입니다.
func ModelKey(name, version string) string {
	return fmt.Sprintf("%s:%s", name, version)
}

// ============================================================================
// 모델 레지스트리
// ============================================================================

// ModelRegistry는 AI 모델 메타데이터 저장소입니다.
type ModelRegistry struct {
	mu     sync.RWMutex
	models map[string]*ModelMetadata
}

// NewModelRegistry는 새 모델 레지스트리를 생성합니다.
func NewModelRegistry() *ModelRegistry {
	return &ModelRegistry{models: make(map[string]*ModelMetadata)}
}

// Register는 새 모델을 등록합니다.
func (r *ModelRegistry) Register(m *ModelMetadata) error {
	if m == nil {
		return errors.New("model is nil")
	}
	if m.Name == "" || m.Version == "" {
		return errors.New("name and version required")
	}
	if m.Status == "" {
		m.Status = RegistryStatusDraft
	}
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now().UTC()
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	key := ModelKey(m.Name, m.Version)
	if _, exists := r.models[key]; exists {
		return fmt.Errorf("model %s already registered", key)
	}
	r.models[key] = m
	return nil
}

// Get은 모델을 조회합니다.
func (r *ModelRegistry) Get(name, version string) (*ModelMetadata, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.models[ModelKey(name, version)]
	if !ok {
		return nil, fmt.Errorf("model %s:%s not found", name, version)
	}
	return m, nil
}

// ListByName은 모델명으로 모든 버전을 조회합니다.
func (r *ModelRegistry) ListByName(name string) []*ModelMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var results []*ModelMetadata
	for _, m := range r.models {
		if m.Name == name {
			results = append(results, m)
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Version > results[j].Version
	})
	return results
}

// GetActive는 운영 또는 카나리 상태의 모델을 반환합니다.
func (r *ModelRegistry) GetActive(name string) (*ModelMetadata, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, m := range r.models {
		if m.Name == name && m.Status == RegistryStatusProduction {
			return m, nil
		}
	}
	return nil, fmt.Errorf("no active model for %q", name)
}

// Promote는 모델을 다음 단계로 승격합니다.
func (r *ModelRegistry) Promote(name, version, targetStatus string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	m, ok := r.models[ModelKey(name, version)]
	if !ok {
		return fmt.Errorf("model %s:%s not found", name, version)
	}

	if !isValidStatusTransition(m.Status, targetStatus) {
		return fmt.Errorf("invalid transition: %s → %s", m.Status, targetStatus)
	}

	// production으로 승격 시 기존 production 모델은 deprecated로 전환
	if targetStatus == RegistryStatusProduction {
		for _, existing := range r.models {
			if existing.Name == name && existing.Status == RegistryStatusProduction {
				existing.Status = RegistryStatusDeprecated
				now := time.Now().UTC()
				existing.DeprecatedAt = &now
			}
		}
	}

	m.Status = targetStatus
	now := time.Now().UTC()
	m.PromotedAt = &now
	return nil
}

func isValidStatusTransition(from, to string) bool {
	allowed := map[string][]string{
		RegistryStatusDraft:      {RegistryStatusStaging},
		RegistryStatusStaging:    {RegistryStatusCanary, RegistryStatusProduction, RegistryStatusDeprecated},
		RegistryStatusCanary:     {RegistryStatusProduction, RegistryStatusDeprecated},
		RegistryStatusProduction: {RegistryStatusDeprecated},
		RegistryStatusDeprecated: {RegistryStatusArchived},
		RegistryStatusArchived:   {},
	}
	for _, t := range allowed[from] {
		if t == to {
			return true
		}
	}
	return false
}

// SetCanaryRatio는 카나리 트래픽 비율을 설정합니다.
func (r *ModelRegistry) SetCanaryRatio(name, version string, ratio float64) error {
	if ratio < 0 || ratio > 1 {
		return errors.New("canary ratio must be between 0 and 1")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	m, ok := r.models[ModelKey(name, version)]
	if !ok {
		return fmt.Errorf("model %s:%s not found", name, version)
	}
	if m.Status != RegistryStatusCanary {
		return fmt.Errorf("model is not in canary state (current: %s)", m.Status)
	}
	m.CanaryRatio = ratio
	return nil
}

// UpdateMetrics는 모델 평가 지표를 갱신합니다.
func (r *ModelRegistry) UpdateMetrics(name, version string, metrics *ModelMetric) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	m, ok := r.models[ModelKey(name, version)]
	if !ok {
		return fmt.Errorf("model %s:%s not found", name, version)
	}
	if metrics.EvaluatedAt.IsZero() {
		metrics.EvaluatedAt = time.Now().UTC()
	}
	m.Metrics = metrics
	return nil
}

// CountByStatus는 상태별 모델 수를 반환합니다.
func (r *ModelRegistry) CountByStatus() map[string]int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	counts := make(map[string]int)
	for _, m := range r.models {
		counts[m.Status]++
	}
	return counts
}

// ============================================================================
// 추론 신뢰도 + 임계값 게이팅
// ============================================================================

// ConfidenceGate는 신뢰도 임계값 검증입니다.
type ConfidenceGate struct {
	MinConfidence float64
}

// NewConfidenceGate는 새 신뢰도 게이트를 생성합니다.
func NewConfidenceGate(min float64) *ConfidenceGate {
	if min < 0 {
		min = 0
	}
	if min > 1 {
		min = 1
	}
	return &ConfidenceGate{MinConfidence: min}
}

// Pass는 신뢰도가 임계값을 통과하는지 확인합니다.
func (g *ConfidenceGate) Pass(confidence float64) bool {
	return confidence >= g.MinConfidence
}

// PassWithReason은 통과 여부와 이유를 반환합니다.
func (g *ConfidenceGate) PassWithReason(confidence float64) (bool, string) {
	if confidence >= g.MinConfidence {
		return true, ""
	}
	return false, fmt.Sprintf("confidence %.3f below threshold %.3f", confidence, g.MinConfidence)
}

// ============================================================================
// XAI - 특성 중요도 (SHAP-유사)
// ============================================================================

// FeatureImportance는 특성 중요도입니다.
type FeatureImportance struct {
	Feature    string
	Importance float64 // 0~1
	Direction  string  // "positive", "negative", "neutral"
	Value      float64
}

// XAIExplainer는 모델 예측에 대한 설명을 생성합니다.
type XAIExplainer struct{}

// NewXAIExplainer는 새 XAI 설명자를 생성합니다.
func NewXAIExplainer() *XAIExplainer {
	return &XAIExplainer{}
}

// Explain은 입력 특성과 가중치로 중요도를 계산합니다.
//
// features: 입력 특성 (key=feature name, value=raw input value)
// weights: 모델 가중치 (key=feature name, value=가중치)
//
// 결과는 |weight × value|를 정규화하여 importance로 반환합니다.
func (e *XAIExplainer) Explain(features map[string]float64, weights map[string]float64) []FeatureImportance {
	if len(features) == 0 || len(weights) == 0 {
		return nil
	}

	contributions := make(map[string]float64)
	totalAbs := 0.0
	for f, v := range features {
		w, ok := weights[f]
		if !ok {
			continue
		}
		c := w * v
		contributions[f] = c
		totalAbs += math.Abs(c)
	}

	if totalAbs == 0 {
		return nil
	}

	importances := make([]FeatureImportance, 0, len(contributions))
	for f, c := range contributions {
		direction := "neutral"
		if c > 0.001 {
			direction = "positive"
		} else if c < -0.001 {
			direction = "negative"
		}
		importances = append(importances, FeatureImportance{
			Feature:    f,
			Importance: math.Abs(c) / totalAbs,
			Direction:  direction,
			Value:      features[f],
		})
	}

	// 중요도 내림차순 정렬
	sort.Slice(importances, func(i, j int) bool {
		return importances[i].Importance > importances[j].Importance
	})

	return importances
}

// TopK는 상위 K개 중요 특성을 반환합니다.
func (e *XAIExplainer) TopK(importances []FeatureImportance, k int) []FeatureImportance {
	if k <= 0 || k >= len(importances) {
		return importances
	}
	return importances[:k]
}

// ============================================================================
// 드리프트 감지 (단순 KS-test)
// ============================================================================

// DriftStatus는 드리프트 상태입니다.
const (
	DriftNone     = "none"
	DriftMinor    = "minor"
	DriftSignificant = "significant"
)

// DriftReport는 드리프트 분석 결과입니다.
type DriftReport struct {
	Feature       string
	BaselineMean  float64
	CurrentMean   float64
	BaselineStd   float64
	CurrentStd    float64
	KSStatistic   float64 // 0~1, 분포 차이 척도
	Status        string
	DetectedAt    time.Time
}

// DriftDetector는 입력 분포 변화를 감지합니다.
type DriftDetector struct {
	threshold float64 // 0.1 = minor, 0.2+ = significant
}

// NewDriftDetector는 새 드리프트 감지기를 생성합니다.
func NewDriftDetector(threshold float64) *DriftDetector {
	if threshold <= 0 {
		threshold = 0.1
	}
	return &DriftDetector{threshold: threshold}
}

// Compare는 baseline과 current 샘플 분포를 비교합니다.
//
// 단순 KS-유사 통계: 평균/표준편차 차이를 정규화한 값.
func (d *DriftDetector) Compare(feature string, baseline, current []float64) *DriftReport {
	if len(baseline) == 0 || len(current) == 0 {
		return nil
	}

	bMean, bStd := meanStd(baseline)
	cMean, cStd := meanStd(current)

	// 표준화된 평균 차이 (effect size)
	pooledStd := (bStd + cStd) / 2
	if pooledStd == 0 {
		pooledStd = 1e-9
	}
	ks := math.Abs(bMean-cMean) / pooledStd
	if ks > 1.0 {
		ks = 1.0
	}

	report := &DriftReport{
		Feature:      feature,
		BaselineMean: bMean,
		CurrentMean:  cMean,
		BaselineStd:  bStd,
		CurrentStd:   cStd,
		KSStatistic:  ks,
		DetectedAt:   time.Now().UTC(),
	}

	switch {
	case ks >= d.threshold*2:
		report.Status = DriftSignificant
	case ks >= d.threshold:
		report.Status = DriftMinor
	default:
		report.Status = DriftNone
	}

	return report
}

// CompareBatch는 여러 특성의 드리프트를 일괄 분석합니다.
func (d *DriftDetector) CompareBatch(baseline, current map[string][]float64) []*DriftReport {
	var reports []*DriftReport
	for feature, base := range baseline {
		cur, ok := current[feature]
		if !ok {
			continue
		}
		if r := d.Compare(feature, base, cur); r != nil {
			reports = append(reports, r)
		}
	}
	sort.Slice(reports, func(i, j int) bool {
		return reports[i].KSStatistic > reports[j].KSStatistic
	})
	return reports
}

// HasSignificantDrift는 보고서 중 significant drift가 있는지 확인합니다.
func HasSignificantDrift(reports []*DriftReport) bool {
	for _, r := range reports {
		if r.Status == DriftSignificant {
			return true
		}
	}
	return false
}

// ============================================================================
// 헬퍼
// ============================================================================

func meanStd(values []float64) (mean, std float64) {
	if len(values) == 0 {
		return 0, 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	mean = sum / float64(len(values))

	if len(values) < 2 {
		return mean, 0
	}
	sumSq := 0.0
	for _, v := range values {
		sumSq += (v - mean) * (v - mean)
	}
	std = math.Sqrt(sumSq / float64(len(values)-1))
	return mean, std
}
