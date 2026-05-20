package service_test

import (
	"context"
	"strings"
	"testing"

	"github.com/manpasik/backend/services/nlp-service/internal/service"
)

// TestScenario_FullEmergencyFlow는 응급 메시지 분석 + 엔티티 추출 흐름입니다.
func TestScenario_FullEmergencyFlow(t *testing.T) {
	router := service.NewIntentRouter()

	text := "환자가 갑자기 가슴이 아프고 호흡곤란이 심해요 의식이 흐려져요"
	result := router.Classify(text)

	if result.Intent != service.IntentEmergency {
		t.Errorf("Intent = %q, want emergency", result.Intent)
	}
	if result.Severity < 80 {
		t.Errorf("Severity = %d, want >= 80", result.Severity)
	}

	entities := router.ExtractEntities(text)
	hasChest := false
	for _, bp := range entities.BodyParts {
		if bp == "가슴" {
			hasChest = true
		}
	}
	if !hasChest {
		t.Errorf("가슴 미추출: %v", entities.BodyParts)
	}
}

// TestScenario_MultiTurnContext는 다중 턴 분석 일관성을 검증합니다.
func TestScenario_MultiTurnContext(t *testing.T) {
	router := service.NewIntentRouter()

	turns := []string{
		"안녕하세요",
		"머리가 아파서 진료 예약 잡고 싶어요",
		"감사합니다",
	}

	expected := []string{
		service.IntentSmalltalk,
		service.IntentReservation,
		service.IntentGratitude,
	}

	for i, t_text := range turns {
		result := router.Classify(t_text)
		if result.Intent != expected[i] {
			t.Errorf("turn %d (%q): Intent = %q, want %q", i, t_text, result.Intent, expected[i])
		}
	}
}

// TestScenario_SymptomSeverityProgression는 증상 심각도 단계 검증입니다.
func TestScenario_SymptomSeverityProgression(t *testing.T) {
	router := service.NewIntentRouter()

	cases := []struct {
		text   string
		minSev int
	}{
		{"가벼운 두통이 있어요", 0},
		{"고열이 나고 피로감이 심해요", 30},
		{"심정지 의심됩니다", 90},
	}

	for _, c := range cases {
		result := router.Classify(c.text)
		if result.Severity < c.minSev {
			t.Errorf("%q: Severity = %d, want >= %d", c.text, result.Severity, c.minSev)
		}
	}
}

// TestScenario_KoreanTokenizerWithMedicalText는 의료 도메인 텍스트 토큰화를 검증합니다.
func TestScenario_KoreanTokenizerWithMedicalText(t *testing.T) {
	kt := service.NewKoreanTokenizer()

	text := "오늘 아침 8시에 측정한 혈압은 140/90 mmHg 이고 심박수는 75 bpm 입니다"
	tokens := kt.Tokenize(text)

	hasTime := false
	hasUnit := false
	for _, tk := range tokens {
		if tk.Type == service.TokenTime {
			hasTime = true
		}
		if tk.Type == service.TokenUnit {
			hasUnit = true
		}
	}
	if !hasTime {
		t.Error("시간 토큰 없음")
	}
	if !hasUnit {
		t.Error("단위 토큰 없음")
	}
}

// TestScenario_AnalyzeKoreanIntegration은 NLPService.AnalyzeKorean 종합 분석입니다.
func TestScenario_AnalyzeKoreanIntegration(t *testing.T) {
	repo := newScenarioMockRepo()
	svc := service.NewNLPService(repo)

	cases := []struct {
		text   string
		intent string
	}{
		{"두통이 심해요", service.IntentSymptom},
		{"내일 진료 예약 잡아주세요", service.IntentReservation},
		{"심정지 의심됩니다", service.IntentEmergency},
	}

	for _, c := range cases {
		analysis := svc.AnalyzeKorean(c.text)
		if analysis.Intent != c.intent {
			t.Errorf("%q: Intent = %q, want %q", c.text, analysis.Intent, c.intent)
		}
		if len(analysis.Tokens) == 0 {
			t.Errorf("%q: 토큰이 비어 있음", c.text)
		}
		if analysis.Entities == nil {
			t.Errorf("%q: 엔티티가 nil", c.text)
		}
	}
}

// TestScenario_TimeExpressionDiversity는 다양한 시간 표현을 검증합니다.
func TestScenario_TimeExpressionDiversity(t *testing.T) {
	kt := service.NewKoreanTokenizer()

	expressions := []string{
		"오늘", "어제", "내일",
		"오전 9시", "오후 3시",
		"3시간 전", "2일 전", "1주 전",
		"지금", "방금",
	}

	for _, e := range expressions {
		tokens := kt.Tokenize(e + " 측정함")
		hasTime := false
		for _, tk := range tokens {
			if tk.Type == service.TokenTime {
				hasTime = true
				break
			}
		}
		if !hasTime {
			t.Errorf("%q: 시간 토큰 미추출", e)
		}
	}
}

// TestScenario_EntityCombination은 다양한 엔티티 조합 추출을 검증합니다.
func TestScenario_EntityCombination(t *testing.T) {
	router := service.NewIntentRouter()

	text := "어제 아침 8시에 머리 통증이 있었고 가슴도 답답했어요. 혈압 140/90 mmHg"
	entities := router.ExtractEntities(text)

	if len(entities.BodyParts) < 2 {
		t.Errorf("BodyParts = %d, want >= 2", len(entities.BodyParts))
	}
	if len(entities.Numbers) < 2 {
		t.Errorf("Numbers = %d, want >= 2", len(entities.Numbers))
	}
	if len(entities.Times) < 1 {
		t.Errorf("Times = %d, want >= 1", len(entities.Times))
	}

	hasMmHg := false
	for _, u := range entities.Units {
		if strings.Contains(u, "mmHg") {
			hasMmHg = true
		}
	}
	if !hasMmHg {
		t.Errorf("mmHg 미추출: %v", entities.Units)
	}
}

// TestScenario_SeverityToUrgency는 긴급도 점수 변환 임계값을 검증합니다.
func TestScenario_SeverityToUrgency(t *testing.T) {
	cases := []struct {
		score int
		want  string
	}{
		{0, service.UrgencyRoutine},
		{34, service.UrgencyRoutine},
		{35, service.UrgencyUrgent},
		{69, service.UrgencyUrgent},
		{70, service.UrgencyEmergency},
		{100, service.UrgencyEmergency},
	}

	for _, c := range cases {
		got := service.SeverityToUrgency(c.score)
		if got != c.want {
			t.Errorf("SeverityToUrgency(%d) = %q, want %q", c.score, got, c.want)
		}
	}
}

// TestScenario_QuestionDetection은 질문형 분류를 검증합니다.
func TestScenario_QuestionDetection(t *testing.T) {
	router := service.NewIntentRouter()

	questions := []string{
		"혈당 정상 수치가 어떻게 되나요?",
		"운동을 얼마나 해야 하나요?",
		"왜 혈압이 갑자기 올라갔나요?",
		"혈당 검사는 언제 하면 좋나요?",
	}

	for _, q := range questions {
		result := router.Classify(q)
		if result.Intent != service.IntentQuestion {
			t.Errorf("%q: Intent = %q, want question", q, result.Intent)
		}
	}
}

// ============================================================================
// 헬퍼
// ============================================================================

type scenarioMockRepo struct{}

func newScenarioMockRepo() *scenarioMockRepo { return &scenarioMockRepo{} }

func (m *scenarioMockRepo) SaveQuery(_ context.Context, _ *service.HealthQuery) error      { return nil }
func (m *scenarioMockRepo) GetQuery(_ context.Context, _ string) (*service.HealthQuery, error) {
	return nil, nil
}
func (m *scenarioMockRepo) SaveExtraction(_ context.Context, _ *service.SymptomExtraction) error {
	return nil
}
func (m *scenarioMockRepo) GetSuggestions(_ context.Context, _ string) ([]service.Suggestion, error) {
	return nil, nil
}
func (m *scenarioMockRepo) SaveSuggestions(_ context.Context, _ string, _ []service.Suggestion) error {
	return nil
}
