package service_test

import (
	"strings"
	"testing"

	"github.com/manpasik/backend/services/nlp-service/internal/repository/memory"
	"github.com/manpasik/backend/services/nlp-service/internal/service"
)

// TestKoreanTokenizer_BasicTokens는 기본 토큰화를 검증합니다.
func TestKoreanTokenizer_BasicTokens(t *testing.T) {
	kt := service.NewKoreanTokenizer()
	tokens := kt.Tokenize("오늘 두통이 심해요")

	if len(tokens) == 0 {
		t.Fatal("토큰이 추출되지 않았습니다")
	}

	hasTime := false
	hasWord := false
	for _, tk := range tokens {
		if tk.Type == service.TokenTime {
			hasTime = true
		}
		if tk.Type == service.TokenWord {
			hasWord = true
		}
	}
	if !hasTime {
		t.Error("시간 토큰('오늘')이 추출되지 않음")
	}
	if !hasWord {
		t.Error("단어 토큰이 추출되지 않음")
	}
}

// TestKoreanTokenizer_NumberAndUnit은 숫자/단위 추출을 검증합니다.
func TestKoreanTokenizer_NumberAndUnit(t *testing.T) {
	kt := service.NewKoreanTokenizer()
	tokens := kt.Tokenize("혈당 120 mg/dL 입니다")

	hasNumber := false
	hasUnit := false
	for _, tk := range tokens {
		if tk.Type == service.TokenNumber && tk.Text == "120" {
			hasNumber = true
		}
		if tk.Type == service.TokenUnit && tk.Text == "mg/dL" {
			hasUnit = true
		}
	}
	if !hasNumber {
		t.Error("숫자 토큰('120')이 추출되지 않음")
	}
	if !hasUnit {
		t.Error("단위 토큰('mg/dL')이 추출되지 않음")
	}
}

// TestKoreanTokenizer_TimeExpressions는 다양한 시간 표현 추출을 검증합니다.
func TestKoreanTokenizer_TimeExpressions(t *testing.T) {
	kt := service.NewKoreanTokenizer()
	cases := []string{"어제", "오전", "3시간 전", "내일 아침"}

	for _, c := range cases {
		tokens := kt.Tokenize(c)
		hasTime := false
		for _, tk := range tokens {
			if tk.Type == service.TokenTime {
				hasTime = true
				break
			}
		}
		if !hasTime {
			t.Errorf("시간 표현 %q에서 시간 토큰이 추출되지 않음", c)
		}
	}
}

// TestIntentRouter_Emergency는 응급 의도 분류를 검증합니다.
func TestIntentRouter_Emergency(t *testing.T) {
	router := service.NewIntentRouter()
	cases := []string{
		"갑자기 의식없이 쓰러졌어요",
		"호흡곤란이 심해요 119 불러주세요",
		"심정지 같아요",
	}

	for _, c := range cases {
		result := router.Classify(c)
		if result.Intent != service.IntentEmergency {
			t.Errorf("텍스트 %q: Intent = %q, want %q", c, result.Intent, service.IntentEmergency)
		}
		if result.Severity < 70 {
			t.Errorf("텍스트 %q: Severity = %d, want >= 70", c, result.Severity)
		}
	}
}

// TestIntentRouter_Symptom은 증상 신고 의도 분류를 검증합니다.
func TestIntentRouter_Symptom(t *testing.T) {
	router := service.NewIntentRouter()
	cases := []string{
		"머리가 아파요 두통이 심해요",
		"기침이 계속 나오고 인후통이 있어요",
		"복통과 설사가 있어요",
	}

	for _, c := range cases {
		result := router.Classify(c)
		if result.Intent != service.IntentSymptom {
			t.Errorf("텍스트 %q: Intent = %q, want %q", c, result.Intent, service.IntentSymptom)
		}
	}
}

// TestIntentRouter_Reservation은 예약 의도 분류를 검증합니다.
func TestIntentRouter_Reservation(t *testing.T) {
	router := service.NewIntentRouter()
	cases := []string{
		"내일 병원 예약 잡아주세요",
		"진료 일정 변경하고 싶어요",
		"방문 예약 취소해주세요",
	}

	for _, c := range cases {
		result := router.Classify(c)
		if result.Intent != service.IntentReservation {
			t.Errorf("텍스트 %q: Intent = %q, want %q", c, result.Intent, service.IntentReservation)
		}
	}
}

// TestIntentRouter_Question은 질문 의도 분류를 검증합니다.
func TestIntentRouter_Question(t *testing.T) {
	router := service.NewIntentRouter()
	cases := []string{
		"혈당 정상 범위가 무엇인가요?",
		"운동을 어떻게 해야 하나요?",
		"왜 혈압이 올라갔나요?",
	}

	for _, c := range cases {
		result := router.Classify(c)
		if result.Intent != service.IntentQuestion {
			t.Errorf("텍스트 %q: Intent = %q, want %q", c, result.Intent, service.IntentQuestion)
		}
	}
}

// TestIntentRouter_Gratitude는 감사 인사 분류를 검증합니다.
func TestIntentRouter_Gratitude(t *testing.T) {
	router := service.NewIntentRouter()
	result := router.Classify("정말 감사합니다 고마워요")

	if result.Intent != service.IntentGratitude {
		t.Errorf("Intent = %q, want %q", result.Intent, service.IntentGratitude)
	}
}

// TestIntentRouter_Severity는 긴급도 점수를 검증합니다.
func TestIntentRouter_Severity(t *testing.T) {
	router := service.NewIntentRouter()

	cases := []struct {
		text    string
		minSev  int
	}{
		{"심정지 의심", 95},
		{"의식없는 상태", 90},
		{"호흡곤란이 심해요", 80},
	}

	for _, c := range cases {
		result := router.Classify(c.text)
		if result.Severity < c.minSev {
			t.Errorf("텍스트 %q: Severity = %d, want >= %d", c.text, result.Severity, c.minSev)
		}
	}
}

// TestExtractEntities_BodyParts는 신체 부위 추출을 검증합니다.
func TestExtractEntities_BodyParts(t *testing.T) {
	router := service.NewIntentRouter()
	ext := router.ExtractEntities("머리가 아프고 가슴 통증도 있어요")

	if len(ext.BodyParts) < 2 {
		t.Errorf("BodyParts 수 = %d, 최소 2개 예상 (머리, 가슴)", len(ext.BodyParts))
	}

	hasHead := false
	hasChest := false
	for _, bp := range ext.BodyParts {
		if bp == "머리" {
			hasHead = true
		}
		if bp == "가슴" {
			hasChest = true
		}
	}
	if !hasHead || !hasChest {
		t.Errorf("머리/가슴 미추출. BodyParts=%v", ext.BodyParts)
	}
}

// TestExtractEntities_Symptoms는 증상 추출을 검증합니다.
func TestExtractEntities_Symptoms(t *testing.T) {
	router := service.NewIntentRouter()
	ext := router.ExtractEntities("두통과 발열, 그리고 기침이 있어요")

	if len(ext.Symptoms) < 3 {
		t.Errorf("Symptoms 수 = %d, 최소 3개 예상", len(ext.Symptoms))
	}
}

// TestExtractEntities_NumbersAndUnits는 숫자+단위 추출을 검증합니다.
func TestExtractEntities_NumbersAndUnits(t *testing.T) {
	router := service.NewIntentRouter()
	ext := router.ExtractEntities("혈압이 140/90 mmHg 이고 심박수 85 bpm")

	if len(ext.Numbers) < 2 {
		t.Errorf("Numbers 수 = %d, 최소 2개 예상", len(ext.Numbers))
	}
	if len(ext.Units) < 2 {
		t.Errorf("Units 수 = %d, 최소 2개 예상 (mmHg, bpm)", len(ext.Units))
	}
}

// TestSeverityToUrgency는 긴급도 점수→urgency 변환을 검증합니다.
func TestSeverityToUrgency(t *testing.T) {
	cases := []struct {
		score  int
		expect string
	}{
		{0, service.UrgencyRoutine},
		{30, service.UrgencyRoutine},
		{50, service.UrgencyUrgent},
		{75, service.UrgencyEmergency},
		{100, service.UrgencyEmergency},
	}

	for _, c := range cases {
		got := service.SeverityToUrgency(c.score)
		if got != c.expect {
			t.Errorf("SeverityToUrgency(%d) = %q, want %q", c.score, got, c.expect)
		}
	}
}

// TestAnalyzeKorean_Integration은 종합 분석을 검증합니다.
func TestAnalyzeKorean_Integration(t *testing.T) {
	repo := memory.NewNLPRepository()
	svc := service.NewNLPService(repo)

	analysis := svc.AnalyzeKorean("오늘 아침부터 두통이 심하고 어지러워요")

	if analysis.Intent != service.IntentSymptom {
		t.Errorf("Intent = %q, want %q", analysis.Intent, service.IntentSymptom)
	}
	if analysis.Confidence <= 0.5 {
		t.Errorf("Confidence = %f, want > 0.5", analysis.Confidence)
	}
	if len(analysis.Tokens) == 0 {
		t.Error("Tokens가 비어 있음")
	}
	if analysis.Entities == nil {
		t.Fatal("Entities가 nil")
	}
	if len(analysis.Entities.Symptoms) == 0 {
		t.Error("Symptoms가 비어 있음")
	}

	// 시간 표현 확인
	hasTime := false
	for _, tk := range analysis.Tokens {
		if tk.Type == service.TokenTime && strings.Contains(tk.Text, "오늘") {
			hasTime = true
			break
		}
	}
	if !hasTime {
		t.Error("시간 표현 '오늘'이 추출되지 않음")
	}
}

