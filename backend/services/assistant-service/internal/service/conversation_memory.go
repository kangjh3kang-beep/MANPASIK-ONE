// Package service: 어시스턴트 대화 메모리 + 의료 프롬프트 + 안전 필터
//
// 슬라이딩 윈도우 메모리, 7종 의료 프롬프트 템플릿, 안전 필터, 후속 질문 생성을 제공합니다.
package service

import (
	"strings"
	"sync"
	"time"
)

// ============================================================================
// 슬라이딩 윈도우 대화 메모리
// ============================================================================

// MemoryEntry는 메모리의 단일 항목입니다.
type MemoryEntry struct {
	Role      string // "user", "assistant"
	Content   string
	Intent    string
	Timestamp time.Time
}

// ConversationMemory는 슬라이딩 윈도우 기반 대화 메모리입니다.
type ConversationMemory struct {
	mu        sync.RWMutex
	sessions  map[string][]MemoryEntry
	windowSize int
}

// NewConversationMemory는 새 대화 메모리를 생성합니다.
//
// windowSize: 세션당 보관할 최대 턴 수 (기본 10)
func NewConversationMemory(windowSize int) *ConversationMemory {
	if windowSize <= 0 {
		windowSize = 10
	}
	return &ConversationMemory{
		sessions:   make(map[string][]MemoryEntry),
		windowSize: windowSize,
	}
}

// Append는 세션에 새 메모리 항목을 추가합니다.
// 윈도우 크기 초과 시 가장 오래된 항목부터 제거됩니다.
func (m *ConversationMemory) Append(sessionID string, entry MemoryEntry) {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	entries := m.sessions[sessionID]
	entries = append(entries, entry)

	// 윈도우 초과 시 오래된 항목 제거
	if len(entries) > m.windowSize {
		entries = entries[len(entries)-m.windowSize:]
	}
	m.sessions[sessionID] = entries
}

// GetRecent는 최근 N개의 메모리 항목을 반환합니다.
func (m *ConversationMemory) GetRecent(sessionID string, n int) []MemoryEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entries := m.sessions[sessionID]
	if n <= 0 || n >= len(entries) {
		return append([]MemoryEntry{}, entries...)
	}
	return append([]MemoryEntry{}, entries[len(entries)-n:]...)
}

// GetAll은 세션의 모든 메모리를 반환합니다.
func (m *ConversationMemory) GetAll(sessionID string) []MemoryEntry {
	return m.GetRecent(sessionID, 0)
}

// Clear는 세션의 메모리를 삭제합니다.
func (m *ConversationMemory) Clear(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, sessionID)
}

// SessionCount는 메모리에 저장된 세션 수를 반환합니다.
func (m *ConversationMemory) SessionCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.sessions)
}

// BuildContext는 LLM에 전달할 컨텍스트 문자열을 빌드합니다.
func (m *ConversationMemory) BuildContext(sessionID string, maxTurns int) string {
	entries := m.GetRecent(sessionID, maxTurns)
	if len(entries) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, e := range entries {
		sb.WriteString(e.Role)
		sb.WriteString(": ")
		sb.WriteString(e.Content)
		sb.WriteString("\n")
	}
	return sb.String()
}

// ============================================================================
// 의료 프롬프트 템플릿 (7종)
// ============================================================================

// PromptTemplate는 의료 도메인 프롬프트 템플릿입니다.
type PromptTemplate struct {
	Name         string
	IntentMatch  []string // 매칭 의도
	SystemPrompt string
	UserExample  string
	Disclaimer   string // 의료 면책 조항
}

// 7종 의료 프롬프트 템플릿
var medicalPromptTemplates = map[string]*PromptTemplate{
	"symptom": {
		Name:        "증상상담",
		IntentMatch: []string{"symptom", "symptom_report"},
		SystemPrompt: `당신은 의료 도메인 어시스턴트입니다. 사용자가 보고한 증상을 듣고:
1. 가능성 있는 원인을 일반적인 수준에서 설명합니다.
2. 자가관리 방법을 안내합니다.
3. 의료진 진료를 받아야 할 상황을 명확히 알립니다.
진단을 직접 내리지 않으며, 전문 의료진과의 상담을 권유하세요.`,
		Disclaimer: "본 응답은 진단이 아니며 의료적 판단을 대체할 수 없습니다. 증상이 지속되면 의료진을 방문하세요.",
	},
	"medication": {
		Name:        "약물복용",
		IntentMatch: []string{"medication", "prescription"},
		SystemPrompt: `당신은 약물 복용 안내 어시스턴트입니다.
1. 약물의 일반적 복용 방법, 효과, 부작용을 설명합니다.
2. 복용 시간/식사 관계를 명확히 안내합니다.
3. 다른 약물과의 상호작용 위험을 경고합니다.
약사/의사 상담 없이 복용 변경을 권하지 마세요.`,
		Disclaimer: "약 복용 변경 전 반드시 처방 의사 또는 약사에게 상담하세요.",
	},
	"diet": {
		Name:        "식이코칭",
		IntentMatch: []string{"diet", "nutrition", "food"},
		SystemPrompt: `당신은 영양 식이 코치입니다. 사용자의 건강 상태에 따라:
1. 균형잡힌 식단을 제안합니다.
2. 피해야 할 음식을 안내합니다.
3. 칼로리/영양소 정보를 제공합니다.
의학적 식이 처방은 영양사/의료진의 영역임을 명심하세요.`,
		Disclaimer: "특수 식이가 필요한 경우 임상영양사 상담을 권합니다.",
	},
	"exercise": {
		Name:        "운동코칭",
		IntentMatch: []string{"exercise", "workout"},
		SystemPrompt: `당신은 운동 코칭 어시스턴트입니다.
1. 사용자의 건강 상태와 목표에 맞는 운동을 추천합니다.
2. 강도/빈도/시간을 명확히 제시합니다.
3. 부상 위험을 경고합니다.
심혈관/근골격계 질환자는 의료진 상담 후 시작하도록 권유하세요.`,
		Disclaimer: "심장 질환, 관절 질환자는 의사와 상담 후 운동을 시작하세요.",
	},
	"emergency": {
		Name:        "응급상황",
		IntentMatch: []string{"emergency", "urgent"},
		SystemPrompt: `당신은 응급 상황 안내 어시스턴트입니다. 사용자의 메시지가 응급 상황을 시사할 경우:
1. 즉시 119 또는 가까운 응급실 방문을 강하게 권합니다.
2. 응급처치 방법을 간결히 안내합니다 (CPR, AED, 출혈 지혈 등).
3. 환자 상태 모니터링 항목을 알립니다.
지체 없이 응급 의료서비스를 호출하도록 유도하세요.`,
		Disclaimer: "이 메시지는 진료가 아닙니다. 즉시 119에 전화하거나 응급실로 이동하세요.",
	},
	"reservation": {
		Name:        "예약안내",
		IntentMatch: []string{"reservation", "schedule"},
		SystemPrompt: `당신은 진료 예약 어시스턴트입니다.
1. 사용자의 증상/요구에 맞는 진료과를 제안합니다.
2. 가까운 의료기관 검색 방법을 안내합니다.
3. 예약 시 준비물 (보험증, 처방전 등)을 알립니다.`,
		Disclaimer: "예약 정보는 변동될 수 있으니 의료기관에 직접 확인하세요.",
	},
	"general": {
		Name:        "일반안내",
		IntentMatch: []string{"question", "general", "smalltalk"},
		SystemPrompt: `당신은 친절하고 정중한 건강 어시스턴트입니다.
1. 사용자의 질문에 명확하고 이해하기 쉽게 답합니다.
2. 의학적 사실 기반으로 응답합니다.
3. 답을 모를 경우 솔직히 인정하고 적절한 정보원을 안내합니다.`,
		Disclaimer: "본 정보는 일반적인 안내이며 의료적 진단을 대체할 수 없습니다.",
	},
}

// MedicalPromptEngine은 의도에 따른 프롬프트를 선택합니다.
type MedicalPromptEngine struct{}

// NewMedicalPromptEngine은 새 프롬프트 엔진을 생성합니다.
func NewMedicalPromptEngine() *MedicalPromptEngine {
	return &MedicalPromptEngine{}
}

// SelectByIntent는 의도에 매칭되는 프롬프트 템플릿을 선택합니다.
func (e *MedicalPromptEngine) SelectByIntent(intent string) *PromptTemplate {
	if intent == "" {
		return medicalPromptTemplates["general"]
	}

	for _, tmpl := range medicalPromptTemplates {
		for _, m := range tmpl.IntentMatch {
			if m == intent {
				return tmpl
			}
		}
	}
	return medicalPromptTemplates["general"]
}

// BuildPrompt는 메모리 컨텍스트와 사용자 입력으로 최종 프롬프트를 생성합니다.
func (e *MedicalPromptEngine) BuildPrompt(intent, contextStr, userInput string) string {
	tmpl := e.SelectByIntent(intent)

	var sb strings.Builder
	sb.WriteString("[SYSTEM]\n")
	sb.WriteString(tmpl.SystemPrompt)
	sb.WriteString("\n\n")

	if contextStr != "" {
		sb.WriteString("[CONVERSATION HISTORY]\n")
		sb.WriteString(contextStr)
		sb.WriteString("\n")
	}

	sb.WriteString("[USER]\n")
	sb.WriteString(userInput)
	sb.WriteString("\n\n")

	sb.WriteString("[ASSISTANT INSTRUCTION]\n")
	sb.WriteString("위 사용자 입력에 대해 응답하고, 마지막에 다음 면책 조항을 추가하세요:\n")
	sb.WriteString("DISCLAIMER: ")
	sb.WriteString(tmpl.Disclaimer)

	return sb.String()
}

// AvailableTemplates는 사용 가능한 모든 템플릿 이름을 반환합니다.
func (e *MedicalPromptEngine) AvailableTemplates() []string {
	names := make([]string, 0, len(medicalPromptTemplates))
	for k := range medicalPromptTemplates {
		names = append(names, k)
	}
	return names
}

// ============================================================================
// 안전 필터 (위험 키워드 차단 + 면책 자동 추가)
// ============================================================================

// 위험 키워드 (응답 차단 또는 강한 경고 필요)
var dangerousKeywords = []string{
	// 자가 진단 관련
	"확실히 진단", "100% 맞다", "절대 ~다",
	// 처방 권유
	"이 약을 처방", "약을 권합니다", "복용량 늘리세요",
	// 위험 행위 권유
	"수술 없이", "병원 안 가도",
}

// 응급 신호 키워드 (응급 응답 강제)
var emergencySignals = []string{
	"자살", "죽고 싶", "쓰러", "의식 없", "호흡 멈", "심정지", "출혈 멈추지 않",
}

// SafetyFilter는 응답 안전성 필터입니다.
type SafetyFilter struct{}

// NewSafetyFilter는 새 안전 필터를 생성합니다.
func NewSafetyFilter() *SafetyFilter {
	return &SafetyFilter{}
}

// SafetyCheck는 안전 검증 결과입니다.
type SafetyCheck struct {
	Safe              bool
	HasEmergencySignal bool
	BlockedKeywords    []string
	Reason             string
}

// Check는 입력 텍스트의 안전성을 검증합니다.
func (s *SafetyFilter) Check(text string) *SafetyCheck {
	check := &SafetyCheck{Safe: true}

	for _, kw := range emergencySignals {
		if strings.Contains(text, kw) {
			check.HasEmergencySignal = true
			check.BlockedKeywords = append(check.BlockedKeywords, kw)
		}
	}

	for _, kw := range dangerousKeywords {
		if strings.Contains(text, kw) {
			check.Safe = false
			check.BlockedKeywords = append(check.BlockedKeywords, kw)
			check.Reason = "위험한 의료 권유 키워드가 감지되었습니다"
		}
	}

	return check
}

// SanitizeResponse는 응답에서 위험 키워드를 마스킹합니다.
func (s *SafetyFilter) SanitizeResponse(text string) string {
	result := text
	for _, kw := range dangerousKeywords {
		result = strings.ReplaceAll(result, kw, "[필터됨]")
	}
	return result
}

// AppendDisclaimer는 응답에 면책 조항을 추가합니다 (이미 있으면 스킵).
func (s *SafetyFilter) AppendDisclaimer(response, disclaimer string) string {
	if strings.Contains(response, "DISCLAIMER") || strings.Contains(response, "면책") {
		return response
	}
	return response + "\n\n[안내] " + disclaimer
}

// ============================================================================
// 후속 질문 생성기
// ============================================================================

// FollowUpGenerator는 의도/문맥에 기반한 후속 질문을 생성합니다.
type FollowUpGenerator struct{}

// NewFollowUpGenerator는 새 후속 질문 생성기를 생성합니다.
func NewFollowUpGenerator() *FollowUpGenerator {
	return &FollowUpGenerator{}
}

// 의도별 후속 질문 사전
var followUpsByIntent = map[string][]string{
	"symptom": {
		"증상이 언제부터 시작되었나요?",
		"통증의 강도를 1-10으로 표현하면 어느 정도인가요?",
		"비슷한 증상이 과거에도 있었나요?",
	},
	"medication": {
		"현재 다른 약을 함께 복용 중이신가요?",
		"이 약을 처음 처방받으신 건가요?",
		"복용 후 어떤 불편함이 있으신가요?",
	},
	"diet": {
		"평소 식사 습관은 어떻게 되시나요?",
		"알레르기나 식이 제한이 있으신가요?",
		"체중 목표가 있으신가요? (감량/유지/증량)",
	},
	"exercise": {
		"평소 운동 빈도는 어떻게 되시나요?",
		"무리 없이 할 수 있는 강도를 알려주세요.",
		"부상이나 관절 통증이 있으신가요?",
	},
	"emergency": {
		"환자의 의식 상태는 어떻습니까?",
		"호흡과 맥박은 정상인가요?",
		"119에 전화하셨나요?",
	},
	"reservation": {
		"어떤 진료과를 원하시나요?",
		"방문 가능한 시간대를 알려주세요.",
		"가까운 지역을 알려주세요.",
	},
	"general": {
		"더 자세히 알고 싶은 부분이 있으신가요?",
		"다른 궁금한 점이 있으신가요?",
		"도움이 되셨나요?",
	},
}

// Generate는 의도에 맞는 후속 질문 3개를 반환합니다.
func (g *FollowUpGenerator) Generate(intent string) []string {
	questions, ok := followUpsByIntent[intent]
	if !ok {
		questions = followUpsByIntent["general"]
	}
	if len(questions) > 3 {
		return questions[:3]
	}
	return questions
}
