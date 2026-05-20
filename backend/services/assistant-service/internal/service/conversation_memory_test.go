package service_test

import (
	"strings"
	"testing"
	"time"

	"github.com/manpasik/backend/services/assistant-service/internal/service"
)

// TestConversationMemory_AppendAndGet은 메모리 저장과 조회를 검증합니다.
func TestConversationMemory_AppendAndGet(t *testing.T) {
	mem := service.NewConversationMemory(5)

	mem.Append("s1", service.MemoryEntry{Role: "user", Content: "안녕하세요"})
	mem.Append("s1", service.MemoryEntry{Role: "assistant", Content: "안녕하세요. 무엇을 도와드릴까요?"})

	entries := mem.GetAll("s1")
	if len(entries) != 2 {
		t.Errorf("entries = %d, want 2", len(entries))
	}
	if entries[0].Role != "user" {
		t.Errorf("첫 번째 role = %q, want user", entries[0].Role)
	}
}

// TestConversationMemory_SlidingWindow는 윈도우 크기 초과 시 오래된 항목 제거를 검증합니다.
func TestConversationMemory_SlidingWindow(t *testing.T) {
	mem := service.NewConversationMemory(3)

	for i := 0; i < 5; i++ {
		mem.Append("s1", service.MemoryEntry{
			Role:    "user",
			Content: string(rune('a' + i)),
		})
	}

	entries := mem.GetAll("s1")
	if len(entries) != 3 {
		t.Errorf("entries = %d, want 3", len(entries))
	}
	// 오래된 'a', 'b'는 제거, 'c', 'd', 'e'만 남아야 함
	if entries[0].Content != "c" {
		t.Errorf("첫 항목 = %q, want c", entries[0].Content)
	}
	if entries[2].Content != "e" {
		t.Errorf("마지막 항목 = %q, want e", entries[2].Content)
	}
}

// TestConversationMemory_GetRecent는 최근 N개 조회를 검증합니다.
func TestConversationMemory_GetRecent(t *testing.T) {
	mem := service.NewConversationMemory(10)
	for i := 0; i < 5; i++ {
		mem.Append("s", service.MemoryEntry{Content: string(rune('1' + i))})
	}

	recent := mem.GetRecent("s", 2)
	if len(recent) != 2 {
		t.Errorf("recent = %d, want 2", len(recent))
	}
	if recent[1].Content != "5" {
		t.Errorf("마지막 = %q, want 5", recent[1].Content)
	}
}

// TestConversationMemory_BuildContext는 컨텍스트 문자열 빌드를 검증합니다.
func TestConversationMemory_BuildContext(t *testing.T) {
	mem := service.NewConversationMemory(5)
	mem.Append("s", service.MemoryEntry{Role: "user", Content: "두통"})
	mem.Append("s", service.MemoryEntry{Role: "assistant", Content: "언제부터?"})

	ctx := mem.BuildContext("s", 5)
	if !strings.Contains(ctx, "user: 두통") {
		t.Errorf("context에 user 메시지 없음: %q", ctx)
	}
	if !strings.Contains(ctx, "assistant: 언제부터?") {
		t.Errorf("context에 assistant 메시지 없음: %q", ctx)
	}
}

// TestConversationMemory_Clear는 세션 메모리 삭제를 검증합니다.
func TestConversationMemory_Clear(t *testing.T) {
	mem := service.NewConversationMemory(10)
	mem.Append("s", service.MemoryEntry{Content: "x"})

	mem.Clear("s")

	if len(mem.GetAll("s")) != 0 {
		t.Error("Clear 후에도 메모리가 남아있음")
	}
	if mem.SessionCount() != 0 {
		t.Error("SessionCount가 0이 아님")
	}
}

// TestConversationMemory_AutoTimestamp는 자동 타임스탬프 설정을 검증합니다.
func TestConversationMemory_AutoTimestamp(t *testing.T) {
	mem := service.NewConversationMemory(5)
	before := time.Now().UTC()
	mem.Append("s", service.MemoryEntry{Content: "x"})
	after := time.Now().UTC()

	entry := mem.GetAll("s")[0]
	if entry.Timestamp.Before(before) || entry.Timestamp.After(after) {
		t.Error("자동 타임스탬프 설정 오류")
	}
}

// TestMedicalPromptEngine_SelectByIntent_Symptom은 증상 의도 매칭을 검증합니다.
func TestMedicalPromptEngine_SelectByIntent_Symptom(t *testing.T) {
	e := service.NewMedicalPromptEngine()
	tmpl := e.SelectByIntent("symptom")

	if tmpl.Name != "증상상담" {
		t.Errorf("Name = %q, want 증상상담", tmpl.Name)
	}
	if !strings.Contains(tmpl.SystemPrompt, "증상") {
		t.Error("SystemPrompt에 '증상' 키워드 없음")
	}
}

// TestMedicalPromptEngine_SelectByIntent_Emergency는 응급 의도 매칭을 검증합니다.
func TestMedicalPromptEngine_SelectByIntent_Emergency(t *testing.T) {
	e := service.NewMedicalPromptEngine()
	tmpl := e.SelectByIntent("emergency")

	if tmpl.Name != "응급상황" {
		t.Errorf("Name = %q, want 응급상황", tmpl.Name)
	}
	if !strings.Contains(tmpl.Disclaimer, "119") {
		t.Errorf("Disclaimer에 '119' 없음: %q", tmpl.Disclaimer)
	}
}

// TestMedicalPromptEngine_SelectByIntent_Unknown은 미지의 의도 시 일반 템플릿 폴백을 검증합니다.
func TestMedicalPromptEngine_SelectByIntent_Unknown(t *testing.T) {
	e := service.NewMedicalPromptEngine()
	tmpl := e.SelectByIntent("xyz_unknown")

	if tmpl.Name != "일반안내" {
		t.Errorf("Name = %q, want 일반안내", tmpl.Name)
	}
}

// TestMedicalPromptEngine_BuildPrompt는 최종 프롬프트 빌드를 검증합니다.
func TestMedicalPromptEngine_BuildPrompt(t *testing.T) {
	e := service.NewMedicalPromptEngine()
	prompt := e.BuildPrompt("symptom", "user: 두통\nassistant: 언제?", "오늘부터 머리가 아파요")

	if !strings.Contains(prompt, "[SYSTEM]") {
		t.Error("프롬프트에 [SYSTEM] 섹션 없음")
	}
	if !strings.Contains(prompt, "[USER]") {
		t.Error("프롬프트에 [USER] 섹션 없음")
	}
	if !strings.Contains(prompt, "DISCLAIMER:") {
		t.Error("프롬프트에 DISCLAIMER 없음")
	}
	if !strings.Contains(prompt, "오늘부터 머리가 아파요") {
		t.Error("사용자 입력이 프롬프트에 포함되지 않음")
	}
}

// TestMedicalPromptEngine_AvailableTemplates는 7개 템플릿 존재를 검증합니다.
func TestMedicalPromptEngine_AvailableTemplates(t *testing.T) {
	e := service.NewMedicalPromptEngine()
	names := e.AvailableTemplates()
	if len(names) != 7 {
		t.Errorf("templates = %d, want 7", len(names))
	}
}

// TestSafetyFilter_DetectDangerous는 위험 키워드 감지를 검증합니다.
func TestSafetyFilter_DetectDangerous(t *testing.T) {
	filter := service.NewSafetyFilter()
	check := filter.Check("이 약을 처방해드릴게요")

	if check.Safe {
		t.Error("위험 키워드 감지 실패")
	}
	if len(check.BlockedKeywords) == 0 {
		t.Error("BlockedKeywords가 비어 있음")
	}
}

// TestSafetyFilter_EmergencySignal은 응급 신호 감지를 검증합니다.
func TestSafetyFilter_EmergencySignal(t *testing.T) {
	filter := service.NewSafetyFilter()
	check := filter.Check("환자가 갑자기 쓰러졌어요")

	if !check.HasEmergencySignal {
		t.Error("응급 신호 감지 실패")
	}
}

// TestSafetyFilter_SafeText는 안전한 텍스트의 통과를 검증합니다.
func TestSafetyFilter_SafeText(t *testing.T) {
	filter := service.NewSafetyFilter()
	check := filter.Check("두통이 있어요. 어떻게 해야 하나요?")

	if !check.Safe {
		t.Error("안전한 텍스트가 차단됨")
	}
}

// TestSafetyFilter_Sanitize는 위험 키워드 마스킹을 검증합니다.
func TestSafetyFilter_Sanitize(t *testing.T) {
	filter := service.NewSafetyFilter()
	sanitized := filter.SanitizeResponse("의사 없이 약을 권합니다")

	if strings.Contains(sanitized, "약을 권합니다") {
		t.Error("위험 키워드가 마스킹되지 않음")
	}
	if !strings.Contains(sanitized, "[필터됨]") {
		t.Error("[필터됨] 마커가 없음")
	}
}

// TestSafetyFilter_AppendDisclaimer는 면책 조항 자동 추가를 검증합니다.
func TestSafetyFilter_AppendDisclaimer(t *testing.T) {
	filter := service.NewSafetyFilter()

	// 새로 추가
	r1 := filter.AppendDisclaimer("응답입니다", "의료 면책")
	if !strings.Contains(r1, "[안내]") {
		t.Error("면책 조항이 추가되지 않음")
	}

	// 이미 있으면 추가하지 않음
	r2 := filter.AppendDisclaimer("응답입니다 DISCLAIMER: 기존", "새로운")
	count := strings.Count(r2, "DISCLAIMER")
	if count > 1 {
		t.Error("DISCLAIMER가 중복 추가됨")
	}
}

// TestFollowUpGenerator_Symptom은 증상 의도 후속 질문을 검증합니다.
func TestFollowUpGenerator_Symptom(t *testing.T) {
	gen := service.NewFollowUpGenerator()
	questions := gen.Generate("symptom")

	if len(questions) != 3 {
		t.Errorf("questions = %d, want 3", len(questions))
	}
	if !strings.Contains(questions[0], "언제") {
		t.Errorf("첫 질문에 '언제' 없음: %q", questions[0])
	}
}

// TestFollowUpGenerator_Emergency는 응급 의도 후속 질문을 검증합니다.
func TestFollowUpGenerator_Emergency(t *testing.T) {
	gen := service.NewFollowUpGenerator()
	questions := gen.Generate("emergency")

	hasEmergencyQuestion := false
	for _, q := range questions {
		if strings.Contains(q, "119") || strings.Contains(q, "의식") {
			hasEmergencyQuestion = true
			break
		}
	}
	if !hasEmergencyQuestion {
		t.Errorf("응급 관련 질문 없음: %v", questions)
	}
}

// TestFollowUpGenerator_FallbackGeneral은 미지 의도 시 일반 질문 폴백을 검증합니다.
func TestFollowUpGenerator_FallbackGeneral(t *testing.T) {
	gen := service.NewFollowUpGenerator()
	questions := gen.Generate("xyz_unknown")

	if len(questions) == 0 {
		t.Error("폴백 질문이 비어 있음")
	}
}
