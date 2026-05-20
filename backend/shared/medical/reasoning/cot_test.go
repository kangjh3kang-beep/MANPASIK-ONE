package reasoning_test

import (
	"strings"
	"testing"

	"github.com/manpasik/backend/shared/medical/reasoning"
)

func TestChainBuilder_FullFlow(t *testing.T) {
	chain := reasoning.NewChainBuilder("patient-001", "환자의 주증상은 무엇인가?").
		Observe("측정값", "혈당 180 mg/dL, HbA1c 7.5%", []string{"glucose=180", "hba1c=7.5"}).
		Differential([]reasoning.DiagnosisCandidate{
			{Name: "제2형 당뇨병", ICD10: "E11", Probability: 0.85},
			{Name: "전당뇨", ICD10: "R73.03", Probability: 0.10},
		}).
		Evidence("진단 근거", []string{"공복혈당 ≥126", "HbA1c ≥6.5"}, nil).
		Hypothesize("제2형 당뇨병", 0.85, "공복혈당과 HbA1c 모두 진단 기준 충족").
		Recommend([]string{"내분비내과 진료", "식이 조절 시작"}, "routine").
		Build("당뇨병 가능성 높음", 0.85)

	if len(chain.Steps) != 5 {
		t.Errorf("Steps = %d, want 5", len(chain.Steps))
	}
	if chain.OverallConfidence != 0.85 {
		t.Errorf("Confidence = %f", chain.OverallConfidence)
	}
	if chain.Duration <= 0 {
		t.Error("Duration이 0 이하")
	}
}

func TestChain_AddStep_AutoOrder(t *testing.T) {
	chain := &reasoning.Chain{}
	for i := 0; i < 3; i++ {
		chain.AddStep(&reasoning.Step{Type: reasoning.StepObservation, Title: "x"})
	}
	for i, s := range chain.Steps {
		if s.Order != i+1 {
			t.Errorf("step %d Order = %d, want %d", i, s.Order, i+1)
		}
	}
}

func TestChain_FilterByType(t *testing.T) {
	chain := reasoning.NewChainBuilder("p", "q").
		Observe("o1", "x", nil).
		Observe("o2", "y", nil).
		Recommend([]string{"a"}, "low").
		Build("done", 0.9)

	obs := chain.FilterByType(reasoning.StepObservation)
	if len(obs) != 2 {
		t.Errorf("observation = %d, want 2", len(obs))
	}

	rec := chain.FilterByType(reasoning.StepRecommendation)
	if len(rec) != 1 {
		t.Errorf("recommendation = %d, want 1", len(rec))
	}
}

func TestChain_HighConfidenceSteps(t *testing.T) {
	chain := &reasoning.Chain{}
	chain.AddStep(&reasoning.Step{Title: "high", Confidence: 0.95})
	chain.AddStep(&reasoning.Step{Title: "low", Confidence: 0.5})
	chain.AddStep(&reasoning.Step{Title: "mid", Confidence: 0.75})

	high := chain.HighConfidenceSteps(0.7)
	if len(high) != 2 {
		t.Errorf("high = %d, want 2", len(high))
	}
}

func TestSHAPExplainer_BasicExplanation(t *testing.T) {
	e := reasoning.NewSHAPExplainer()

	features := map[string]float64{
		"glucose":   180.0,
		"bp_sys":    140.0,
		"weight_kg": 75.0,
	}
	weights := map[string]float64{
		"glucose":   0.3,
		"bp_sys":    0.2,
		"weight_kg": 0.05,
	}

	exp, err := e.Explain(features, weights, 50.0)
	if err != nil {
		t.Fatalf("Explain 실패: %v", err)
	}

	if exp.BaseValue != 50.0 {
		t.Errorf("BaseValue = %f", exp.BaseValue)
	}
	expected := 50.0 + 180*0.3 + 140*0.2 + 75*0.05
	if exp.PredictedValue < expected-0.01 || exp.PredictedValue > expected+0.01 {
		t.Errorf("PredictedValue = %f, want %f", exp.PredictedValue, expected)
	}

	// 기여도 정렬 검증
	for i := 1; i < len(exp.Contributions); i++ {
		if exp.Contributions[i].AbsoluteImpact > exp.Contributions[i-1].AbsoluteImpact {
			t.Errorf("정렬 위반 idx %d", i)
		}
	}
}

func TestSHAPExplainer_DirectionAssignment(t *testing.T) {
	e := reasoning.NewSHAPExplainer()
	exp, _ := e.Explain(
		map[string]float64{"f1": 1.0, "f2": -2.0, "f3": 0.0},
		map[string]float64{"f1": 0.5, "f2": 0.5, "f3": 0.5},
		0.0,
	)

	for _, c := range exp.Contributions {
		if c.Feature == "f1" && c.Direction != "increases_risk" {
			t.Errorf("f1 direction = %q", c.Direction)
		}
		if c.Feature == "f2" && c.Direction != "decreases_risk" {
			t.Errorf("f2 direction = %q", c.Direction)
		}
	}
}

func TestSHAPExplainer_TopPositiveNegative(t *testing.T) {
	e := reasoning.NewSHAPExplainer()
	features := map[string]float64{
		"a": 10, "b": 20, "c": -5, "d": 15, "e": -10,
	}
	weights := map[string]float64{
		"a": 1, "b": 1, "c": 1, "d": 1, "e": 1,
	}

	exp, _ := e.Explain(features, weights, 0)

	if len(exp.TopPositive) == 0 {
		t.Error("TopPositive 비어있음")
	}
	if len(exp.TopNegative) == 0 {
		t.Error("TopNegative 비어있음")
	}
}

func TestSHAPExplainer_Summary(t *testing.T) {
	e := reasoning.NewSHAPExplainer()
	exp, _ := e.Explain(
		map[string]float64{"glucose": 180},
		map[string]float64{"glucose": 0.3},
		50,
	)

	if !strings.Contains(exp.Summary, "위험도") {
		t.Errorf("Summary = %q", exp.Summary)
	}
}

func TestSHAPExplainer_EmptyInput(t *testing.T) {
	e := reasoning.NewSHAPExplainer()
	_, err := e.Explain(map[string]float64{}, map[string]float64{"x": 1}, 0)
	if err == nil {
		t.Error("빈 features 통과")
	}
}

func TestIntegratedReasoner_FullPipeline(t *testing.T) {
	r := reasoning.NewIntegratedReasoner()

	observations := map[string]float64{
		"2345-7":  180.0, // glucose
		"4548-4":  7.5,   // hba1c
	}
	weights := map[string]float64{
		"2345-7": 0.3,
		"4548-4": 0.4,
	}
	candidates := []reasoning.DiagnosisCandidate{
		{Name: "제2형 당뇨병", ICD10: "E11", Probability: 0.85, Reasoning: "공복혈당 + HbA1c 모두 충족"},
	}

	result, err := r.Reason("patient-001", "당뇨 가능성 평가", observations, weights, 50.0, candidates)
	if err != nil {
		t.Fatalf("Reason 실패: %v", err)
	}

	if result.Chain == nil {
		t.Error("Chain nil")
	}
	if result.Explanation == nil {
		t.Error("Explanation nil")
	}
	if len(result.Chain.Steps) < 2 {
		t.Errorf("Steps = %d, want >= 2", len(result.Chain.Steps))
	}
}

func TestIntegratedReasoner_NoPatient(t *testing.T) {
	r := reasoning.NewIntegratedReasoner()
	_, err := r.Reason("", "q", map[string]float64{"x": 1}, map[string]float64{"x": 1}, 0, nil)
	if err == nil {
		t.Error("환자 ID 없이 통과")
	}
}

func TestChainBuilder_Caution(t *testing.T) {
	chain := reasoning.NewChainBuilder("p", "q").
		Caution("이 결과는 진단이 아닙니다").
		Build("done", 0.9)

	cautions := chain.FilterByType(reasoning.StepCaution)
	if len(cautions) != 1 {
		t.Errorf("cautions = %d, want 1", len(cautions))
	}
}

func TestChainBuilder_RuleOut(t *testing.T) {
	chain := reasoning.NewChainBuilder("p", "q").
		RuleOut("심근경색", "트로포닌 정상").
		Build("done", 0.9)

	ruleOuts := chain.FilterByType(reasoning.StepRuleOut)
	if len(ruleOuts) != 1 {
		t.Errorf("rule_outs = %d, want 1", len(ruleOuts))
	}
}

func TestSHAPExplainer_BaselinePreserved(t *testing.T) {
	e := reasoning.NewSHAPExplainer()
	exp, _ := e.Explain(
		map[string]float64{"x": 0}, // 기여도 0
		map[string]float64{"x": 1.0},
		100.0,
	)
	if exp.PredictedValue != 100.0 {
		t.Errorf("0 기여인데 predicted = %f, want 100", exp.PredictedValue)
	}
}

func TestChainBuilder_Citations(t *testing.T) {
	citations := []*reasoning.Citation{
		{Source: "PMID", ID: "31341544", Title: "ARTESH study"},
		{Source: "PMID", ID: "41619752", Title: "TRICORDER RCT"},
	}

	chain := reasoning.NewChainBuilder("p", "q").
		Evidence("연구 근거", []string{"AR/햅틱 효과", "TRICORDER RCT 결과"}, citations).
		Build("done", 0.9)

	ev := chain.FilterByType(reasoning.StepEvidence)
	if len(ev) == 0 {
		t.Fatal("Evidence step 없음")
	}
	if len(ev[0].Citations) != 2 {
		t.Errorf("Citations = %d, want 2", len(ev[0].Citations))
	}
}

func TestIntegratedReasoner_BuildTimeReasonable(t *testing.T) {
	r := reasoning.NewIntegratedReasoner()
	result, _ := r.Reason(
		"p", "q",
		map[string]float64{"a": 1},
		map[string]float64{"a": 1},
		0,
		nil,
	)
	if result.BuildTimeMs > 1000 {
		t.Errorf("BuildTimeMs = %d, want <= 1000ms", result.BuildTimeMs)
	}
}
