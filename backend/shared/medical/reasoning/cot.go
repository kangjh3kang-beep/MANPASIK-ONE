// Package reasoning은 의료 임상 의사결정을 위한 CoT(Chain-of-Thought) +
// SHAP-style 설명을 제공합니다.
//
// P10-09b D-03 차별화의 핵심: 추론사슬 + 95% CI + 특성 기여도 가시화.
package reasoning

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

// ============================================================================
// CoT 추론 단계
// ============================================================================

// StepType은 추론 단계 종류입니다.
type StepType string

const (
	StepObservation     StepType = "observation"      // 관측: 측정값/증상
	StepDifferential    StepType = "differential"     // 감별진단 후보
	StepEvidence        StepType = "evidence"         // 근거 제시
	StepRuleOut         StepType = "rule_out"         // 배제진단
	StepHypothesis      StepType = "hypothesis"       // 가설 채택
	StepRecommendation  StepType = "recommendation"   // 권고
	StepCaution         StepType = "caution"          // 주의 사항
)

// Step은 추론 단계 1개입니다.
type Step struct {
	ID           string
	Type         StepType
	Order        int
	Title        string
	Description  string
	Confidence   float64 // 0~1
	Evidence     []string
	References   []string // PMID, LOINC 코드 등
	Citations    []*Citation
	Children     []*Step
	Timestamp    time.Time
}

// Citation은 인용 정보입니다.
type Citation struct {
	Source      string // "PMID", "LOINC", "ICD-10", "internal"
	ID          string
	Title       string
	URL         string
	AccessedAt  time.Time
}

// ============================================================================
// 추론 사슬
// ============================================================================

// Chain은 완전한 추론 사슬입니다.
type Chain struct {
	ID            string
	PatientID     string
	EncounterID   string
	Question      string
	Answer        string
	Steps         []*Step
	Conclusion    string
	OverallConfidence float64
	StartedAt     time.Time
	CompletedAt   time.Time
	Duration      time.Duration
}

// AddStep은 추론 단계를 추가합니다.
func (c *Chain) AddStep(step *Step) {
	if step.Order == 0 {
		step.Order = len(c.Steps) + 1
	}
	if step.Timestamp.IsZero() {
		step.Timestamp = time.Now().UTC()
	}
	c.Steps = append(c.Steps, step)
}

// Complete는 추론 사슬을 완료 처리합니다.
func (c *Chain) Complete(conclusion string, overallConfidence float64) {
	c.Conclusion = conclusion
	c.OverallConfidence = overallConfidence
	c.CompletedAt = time.Now().UTC()
	if !c.StartedAt.IsZero() {
		c.Duration = c.CompletedAt.Sub(c.StartedAt)
	}
}

// FilterByType은 특정 타입의 단계만 반환합니다.
func (c *Chain) FilterByType(stepType StepType) []*Step {
	var result []*Step
	for _, s := range c.Steps {
		if s.Type == stepType {
			result = append(result, s)
		}
	}
	return result
}

// HighConfidenceSteps는 신뢰도 임계값 이상의 단계만 반환합니다.
func (c *Chain) HighConfidenceSteps(threshold float64) []*Step {
	var result []*Step
	for _, s := range c.Steps {
		if s.Confidence >= threshold {
			result = append(result, s)
		}
	}
	return result
}

// ============================================================================
// CoT 빌더
// ============================================================================

// ChainBuilder는 추론 사슬을 단계적으로 구축합니다.
type ChainBuilder struct {
	chain *Chain
}

// NewChainBuilder는 새 빌더를 생성합니다.
func NewChainBuilder(patientID, question string) *ChainBuilder {
	return &ChainBuilder{
		chain: &Chain{
			ID:        fmt.Sprintf("cot-%d", time.Now().UnixNano()),
			PatientID: patientID,
			Question:  question,
			StartedAt: time.Now().UTC(),
		},
	}
}

// Observe는 관측 단계를 추가합니다.
func (b *ChainBuilder) Observe(title, description string, evidence []string) *ChainBuilder {
	b.chain.AddStep(&Step{
		ID:          fmt.Sprintf("step-obs-%d", len(b.chain.Steps)+1),
		Type:        StepObservation,
		Title:       title,
		Description: description,
		Confidence:  1.0, // 관측은 항상 확정
		Evidence:    evidence,
	})
	return b
}

// Differential은 감별진단 단계를 추가합니다.
func (b *ChainBuilder) Differential(diagnoses []DiagnosisCandidate) *ChainBuilder {
	desc := "감별진단 후보:\n"
	var refs []string
	for _, d := range diagnoses {
		desc += fmt.Sprintf("  - %s (%.0f%%)\n", d.Name, d.Probability*100)
		if d.ICD10 != "" {
			refs = append(refs, "ICD-10:"+d.ICD10)
		}
	}
	b.chain.AddStep(&Step{
		ID:          fmt.Sprintf("step-dx-%d", len(b.chain.Steps)+1),
		Type:        StepDifferential,
		Title:       "감별진단",
		Description: desc,
		Confidence:  0.8,
		References:  refs,
	})
	return b
}

// Evidence는 근거 단계를 추가합니다.
func (b *ChainBuilder) Evidence(title string, evidenceList []string, citations []*Citation) *ChainBuilder {
	b.chain.AddStep(&Step{
		ID:          fmt.Sprintf("step-ev-%d", len(b.chain.Steps)+1),
		Type:        StepEvidence,
		Title:       title,
		Description: strings.Join(evidenceList, "\n"),
		Confidence:  0.9,
		Evidence:    evidenceList,
		Citations:   citations,
	})
	return b
}

// RuleOut은 배제진단 단계를 추가합니다.
func (b *ChainBuilder) RuleOut(diagnosis string, reason string) *ChainBuilder {
	b.chain.AddStep(&Step{
		ID:          fmt.Sprintf("step-ro-%d", len(b.chain.Steps)+1),
		Type:        StepRuleOut,
		Title:       fmt.Sprintf("배제: %s", diagnosis),
		Description: reason,
		Confidence:  0.85,
	})
	return b
}

// Hypothesize는 가설 채택 단계를 추가합니다.
func (b *ChainBuilder) Hypothesize(diagnosis string, confidence float64, reasoning string) *ChainBuilder {
	b.chain.AddStep(&Step{
		ID:          fmt.Sprintf("step-hyp-%d", len(b.chain.Steps)+1),
		Type:        StepHypothesis,
		Title:       fmt.Sprintf("주 진단 가설: %s", diagnosis),
		Description: reasoning,
		Confidence:  confidence,
	})
	return b
}

// Recommend는 권고 단계를 추가합니다.
func (b *ChainBuilder) Recommend(actions []string, urgency string) *ChainBuilder {
	desc := fmt.Sprintf("긴급도: %s\n", urgency)
	for i, a := range actions {
		desc += fmt.Sprintf("  %d. %s\n", i+1, a)
	}
	b.chain.AddStep(&Step{
		ID:          fmt.Sprintf("step-rec-%d", len(b.chain.Steps)+1),
		Type:        StepRecommendation,
		Title:       "권고 조치",
		Description: desc,
		Confidence:  0.9,
	})
	return b
}

// Caution은 주의 사항을 추가합니다.
func (b *ChainBuilder) Caution(warning string) *ChainBuilder {
	b.chain.AddStep(&Step{
		ID:          fmt.Sprintf("step-caut-%d", len(b.chain.Steps)+1),
		Type:        StepCaution,
		Title:       "주의 사항",
		Description: warning,
		Confidence:  1.0,
	})
	return b
}

// Build는 최종 추론 사슬을 생성합니다.
func (b *ChainBuilder) Build(conclusion string, overallConfidence float64) *Chain {
	b.chain.Answer = conclusion
	b.chain.Complete(conclusion, overallConfidence)
	return b.chain
}

// DiagnosisCandidate는 감별진단 후보입니다.
type DiagnosisCandidate struct {
	Name        string
	ICD10       string
	Probability float64
	Reasoning   string
}

// ============================================================================
// SHAP-style 특성 기여도
// ============================================================================

// FeatureContribution은 단일 특성의 기여도입니다.
type FeatureContribution struct {
	Feature      string
	Value        float64
	Contribution float64 // 양수=양성 방향, 음수=음성 방향
	AbsoluteImpact float64
	Direction    string // "increases_risk", "decreases_risk", "neutral"
	Explanation  string
}

// SHAPExplanation은 SHAP-style 설명입니다.
type SHAPExplanation struct {
	BaseValue       float64               // 평균 예측값
	PredictedValue  float64               // 실제 예측값
	Contributions   []*FeatureContribution
	TopPositive     []*FeatureContribution // 상위 양성 기여
	TopNegative     []*FeatureContribution // 상위 음성 기여
	Summary         string
}

// SHAPExplainer는 SHAP-유사 설명을 생성합니다.
//
// 입력: 특성값, 가중치, 베이스라인 → 각 특성의 contribution = weight * (value - baseline)
type SHAPExplainer struct{}

// NewSHAPExplainer는 새 설명자를 생성합니다.
func NewSHAPExplainer() *SHAPExplainer { return &SHAPExplainer{} }

// Explain은 모델 예측에 대한 설명을 생성합니다.
//
// features: 특성값
// weights: 모델 가중치 (linear 모델 가정)
// baseline: 평균 예측값
func (e *SHAPExplainer) Explain(
	features map[string]float64,
	weights map[string]float64,
	baseline float64,
) (*SHAPExplanation, error) {
	if len(features) == 0 || len(weights) == 0 {
		return nil, errors.New("features and weights required")
	}

	predicted := baseline
	contributions := make([]*FeatureContribution, 0, len(features))

	for f, v := range features {
		w, ok := weights[f]
		if !ok {
			continue
		}
		contribution := w * v
		predicted += contribution

		direction := "neutral"
		if contribution > 0.001 {
			direction = "increases_risk"
		} else if contribution < -0.001 {
			direction = "decreases_risk"
		}

		contributions = append(contributions, &FeatureContribution{
			Feature:        f,
			Value:          v,
			Contribution:   contribution,
			AbsoluteImpact: math.Abs(contribution),
			Direction:      direction,
		})
	}

	// 절대값 내림차순 정렬
	sort.Slice(contributions, func(i, j int) bool {
		return contributions[i].AbsoluteImpact > contributions[j].AbsoluteImpact
	})

	// 양성/음성 분리
	var positive, negative []*FeatureContribution
	for _, c := range contributions {
		if c.Contribution > 0 {
			positive = append(positive, c)
		} else if c.Contribution < 0 {
			negative = append(negative, c)
		}
	}

	return &SHAPExplanation{
		BaseValue:      baseline,
		PredictedValue: predicted,
		Contributions:  contributions,
		TopPositive:    topK(positive, 5),
		TopNegative:    topK(negative, 5),
		Summary:        generateSummary(contributions, baseline, predicted),
	}, nil
}

func topK(in []*FeatureContribution, k int) []*FeatureContribution {
	if k <= 0 || k >= len(in) {
		return in
	}
	return in[:k]
}

func generateSummary(contributions []*FeatureContribution, baseline, predicted float64) string {
	if len(contributions) == 0 {
		return "기여 특성 없음"
	}

	delta := predicted - baseline
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("기준값 %.2f에서 %.2f로 ", baseline, predicted))
	if delta > 0 {
		sb.WriteString(fmt.Sprintf("위험도 증가 (+%.2f).", delta))
	} else if delta < 0 {
		sb.WriteString(fmt.Sprintf("위험도 감소 (%.2f).", delta))
	} else {
		sb.WriteString("변화 없음.")
	}

	sb.WriteString(" 가장 큰 영향: ")
	if len(contributions) > 0 {
		c := contributions[0]
		sb.WriteString(fmt.Sprintf("%s (값=%.2f, 기여=%.2f)", c.Feature, c.Value, c.Contribution))
	}
	return sb.String()
}

// ============================================================================
// 통합 추론 엔진 (CoT + SHAP)
// ============================================================================

// IntegratedReasoner는 CoT 추론과 SHAP 설명을 결합합니다.
type IntegratedReasoner struct {
	shap *SHAPExplainer
}

// NewIntegratedReasoner는 새 통합 추론기를 생성합니다.
func NewIntegratedReasoner() *IntegratedReasoner {
	return &IntegratedReasoner{shap: NewSHAPExplainer()}
}

// ReasoningResult는 통합 추론 결과입니다.
type ReasoningResult struct {
	Chain       *Chain
	Explanation *SHAPExplanation
	BuildTimeMs int64
}

// Reason은 측정값으로부터 통합 추론을 수행합니다.
//
// patientID: 환자 ID
// question: 임상 질문
// observations: LOINC code → 측정값 (e.g., "2345-7" → 110.5)
// weights: 모델 가중치
// baseline: 평균 예측값
// diagnosisCandidates: 감별진단 후보
func (r *IntegratedReasoner) Reason(
	patientID string,
	question string,
	observations map[string]float64,
	weights map[string]float64,
	baseline float64,
	candidates []DiagnosisCandidate,
) (*ReasoningResult, error) {
	if patientID == "" {
		return nil, errors.New("patient_id required")
	}
	startTime := time.Now()

	// 1. SHAP 설명 생성
	explanation, err := r.shap.Explain(observations, weights, baseline)
	if err != nil {
		return nil, fmt.Errorf("SHAP 실패: %w", err)
	}

	// 2. CoT 사슬 빌드
	builder := NewChainBuilder(patientID, question)

	// 관측 단계
	var evidence []string
	for code, value := range observations {
		evidence = append(evidence, fmt.Sprintf("%s: %.2f", code, value))
	}
	builder.Observe("측정값 수집", "POCT 디바이스로부터 다음 값들을 측정했습니다", evidence)

	// 감별진단
	if len(candidates) > 0 {
		builder.Differential(candidates)
	}

	// SHAP 근거
	if len(explanation.TopPositive) > 0 {
		var ev []string
		for _, c := range explanation.TopPositive[:min(3, len(explanation.TopPositive))] {
			ev = append(ev, fmt.Sprintf("%s 기여 +%.2f", c.Feature, c.Contribution))
		}
		builder.Evidence("위험 증가 요인", ev, nil)
	}

	// 주 진단 가설
	if len(candidates) > 0 {
		top := candidates[0]
		builder.Hypothesize(top.Name, top.Probability, top.Reasoning)
	}

	// 권고
	if explanation.PredictedValue > baseline*1.5 {
		builder.Recommend([]string{"추가 검사 권장", "전문의 진료 예약"}, "urgent")
	}

	chain := builder.Build("통합 분석 완료", explanation.PredictedValue)

	return &ReasoningResult{
		Chain:       chain,
		Explanation: explanation,
		BuildTimeMs: time.Since(startTime).Milliseconds(),
	}, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
