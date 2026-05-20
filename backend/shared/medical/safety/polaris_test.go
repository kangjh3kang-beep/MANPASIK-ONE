package safety_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/manpasik/backend/shared/medical/safety"
)

func TestPolaris_DiagnosisRewrite(t *testing.T) {
	p := safety.NewPolarisSafetyNet()
	result := p.Check("당뇨병이 확실히 맞습니다")

	if result.Safe {
		// rewrite는 safe로 분류 가능 (block이 아니므로)
		t.Logf("Safe = %v (rewrite는 safe 가능)", result.Safe)
	}
	if len(result.Violations) == 0 {
		t.Error("확정 진단 발화 미감지")
	}
	if !strings.Contains(result.SanitizedText, "전문의") {
		t.Errorf("Sanitize 실패: %q", result.SanitizedText)
	}
}

func TestPolaris_PrescriptionBlocked(t *testing.T) {
	p := safety.NewPolarisSafetyNet()
	result := p.Check("이 약을 처방해드립니다")

	if !result.HasBlockedContent() {
		t.Error("처방 권유가 차단되지 않음")
	}
	if result.Safe {
		t.Error("처방 발화가 safe로 분류")
	}
}

func TestPolaris_DosageChange(t *testing.T) {
	p := safety.NewPolarisSafetyNet()
	result := p.Check("복용량을 두 배로 늘리세요")

	if !result.HasBlockedContent() {
		t.Error("복용량 변경 권유 미차단")
	}
}

func TestPolaris_EmergencyEscalation(t *testing.T) {
	p := safety.NewPolarisSafetyNet()
	result := p.Check("환자가 의식없이 쓰러졌습니다")

	if !result.RequiresEscalation() {
		t.Error("응급 신호가 escalation을 요구하지 않음")
	}
	if !contains(result.CategoriesViolated(), safety.CategoryEmergency) {
		t.Error("Emergency 카테고리 미감지")
	}
}

func TestPolaris_SelfHarmEscalation(t *testing.T) {
	p := safety.NewPolarisSafetyNet()
	result := p.Check("저는 살기 싫어요")

	if !result.RequiresEscalation() {
		t.Error("자해 위험 escalation 안됨")
	}
}

func TestPolaris_PIIRedaction(t *testing.T) {
	p := safety.NewPolarisSafetyNet()
	result := p.Check("주민번호 901020-1234567 입니다")

	if !strings.Contains(result.SanitizedText, "[개인정보]") {
		t.Errorf("주민등록번호 미가림: %q", result.SanitizedText)
	}
	if strings.Contains(result.SanitizedText, "901020-1234567") {
		t.Error("주민등록번호가 노출됨")
	}
}

func TestPolaris_PhoneRedaction(t *testing.T) {
	p := safety.NewPolarisSafetyNet()
	result := p.Check("연락처는 010-1234-5678 입니다")

	if !strings.Contains(result.SanitizedText, "[연락처]") {
		t.Errorf("전화번호 미가림: %q", result.SanitizedText)
	}
}

func TestPolaris_SafeText(t *testing.T) {
	p := safety.NewPolarisSafetyNet()
	result := p.Check("두통이 있을 때는 충분한 휴식을 권합니다")

	if !result.Safe {
		t.Errorf("정상 텍스트가 unsafe로 분류: violations=%v", result.Violations)
	}
	if !strings.Contains(result.SanitizedText, "[안내]") {
		t.Error("정상 텍스트에 면책 조항 미추가")
	}
}

func TestPolaris_DangerousAdvice(t *testing.T) {
	p := safety.NewPolarisSafetyNet()
	result := p.Check("병원 가지 말고 집에서 쉬세요")

	if !result.HasBlockedContent() {
		t.Error("위험 조언 미차단")
	}
}

func TestPolaris_SurgeryRewrite(t *testing.T) {
	p := safety.NewPolarisSafetyNet()
	result := p.Check("이 경우 수술이 필요합니다")

	hasSurgery := false
	for _, v := range result.Violations {
		if v.Category == safety.CategorySurgery {
			hasSurgery = true
			break
		}
	}
	if !hasSurgery {
		t.Error("수술 권유 미감지")
	}
}

func TestPolaris_HighestSeverityTracking(t *testing.T) {
	p := safety.NewPolarisSafetyNet()
	// emergency(10) > diagnosis(9)
	result := p.Check("당뇨병이 확실히 맞으니 의식없이 쓰러진 것입니다")

	if result.HighestSeverity != 10 {
		t.Errorf("HighestSeverity = %d, want 10 (emergency)", result.HighestSeverity)
	}
	if result.HighestAction != safety.ActionEscalate {
		t.Errorf("HighestAction = %q, want escalate", result.HighestAction)
	}
}

func TestPolaris_AddCustomRule(t *testing.T) {
	p := safety.NewPolarisSafetyNet()

	rule := &safety.Rule{
		ID:          "CUSTOM-001",
		Category:    safety.CategoryDangerousAdvice,
		Severity:    7,
		Action:      safety.ActionWarn,
		Pattern:     regexp.MustCompile(`(검증되지\s*않은|미승인)`),
		Description: "미승인 약물 언급",
	}
	if err := p.AddRule(rule); err != nil {
		t.Fatalf("AddRule 실패: %v", err)
	}

	result := p.Check("미승인 약물입니다")
	if len(result.Violations) == 0 {
		t.Error("커스텀 룰이 적용되지 않음")
	}
}

func TestPolaris_AddRule_Validation(t *testing.T) {
	p := safety.NewPolarisSafetyNet()

	if err := p.AddRule(nil); err == nil {
		t.Error("nil rule 통과")
	}
	if err := p.AddRule(&safety.Rule{ID: "x"}); err == nil {
		t.Error("Pattern 없이 통과")
	}
}

func TestPolaris_EmptyText(t *testing.T) {
	p := safety.NewPolarisSafetyNet()
	result := p.Check("")
	if !result.Safe {
		t.Error("빈 텍스트가 unsafe")
	}
}

func TestPolaris_NoDoubleDisclaimer(t *testing.T) {
	p := safety.NewPolarisSafetyNet()
	result := p.Check("정상적인 안내. [안내] 이미 있음")
	count := strings.Count(result.SanitizedText, "[안내]")
	if count > 1 {
		t.Errorf("[안내] 중복 추가: %d번", count)
	}
}

func TestPolaris_CustomDisclaimer(t *testing.T) {
	p := safety.NewPolarisSafetyNet()
	p.SetDisclaimer("커스텀 면책 조항입니다")

	result := p.Check("정상적인 안내")
	if !strings.Contains(result.SanitizedText, "커스텀 면책") {
		t.Error("커스텀 면책 미적용")
	}
}

func TestAuditLog_RecordAndFind(t *testing.T) {
	log := safety.NewAuditLog(100)
	p := safety.NewPolarisSafetyNet()

	result := p.Check("이 약을 처방하세요")
	log.Record("user-001", "session-A", result)

	entries := log.FindBySession("session-A")
	if len(entries) != 1 {
		t.Errorf("entries = %d, want 1", len(entries))
	}
}

func TestAuditLog_NoViolationNotRecorded(t *testing.T) {
	log := safety.NewAuditLog(100)
	p := safety.NewPolarisSafetyNet()

	result := p.Check("두통이 있을 때 휴식을 권합니다")
	log.Record("user", "session-X", result)

	if log.Count() != 0 {
		t.Error("정상 텍스트가 감사 로그에 기록됨")
	}
}

func TestAuditLog_CountByCategory(t *testing.T) {
	log := safety.NewAuditLog(100)
	p := safety.NewPolarisSafetyNet()

	log.Record("u", "s", p.Check("이 약을 처방"))         // prescription
	log.Record("u", "s", p.Check("의식없이 쓰러"))         // emergency
	log.Record("u", "s", p.Check("복용량을 늘리세요"))    // dosage_change

	counts := log.CountByCategory()
	if counts[safety.CategoryPrescription] == 0 {
		t.Error("Prescription 카운트 누락")
	}
	if counts[safety.CategoryEmergency] == 0 {
		t.Error("Emergency 카운트 누락")
	}
}

func TestAuditLog_MaxSizeWindow(t *testing.T) {
	log := safety.NewAuditLog(2)
	p := safety.NewPolarisSafetyNet()

	for i := 0; i < 5; i++ {
		log.Record("u", "s", p.Check("이 약을 처방"))
	}

	if log.Count() != 2 {
		t.Errorf("Count = %d, want 2 (max)", log.Count())
	}
}

func contains(s []safety.Category, target safety.Category) bool {
	for _, c := range s {
		if c == target {
			return true
		}
	}
	return false
}
