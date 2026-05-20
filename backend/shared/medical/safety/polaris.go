// Package safety는 의료 도메인 특화 안전망 (Polaris-style)을 제공합니다.
//
// FMEA F-V03 (디지털 주치의 의료법 오발화, RPN 216)의 핵심 완화책.
// 디지털 주치의 응답을 발송 전에 검사하여 위험 발화를 차단/수정합니다.
package safety

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// ============================================================================
// 안전 카테고리
// ============================================================================

// Category는 안전 위반 카테고리입니다.
type Category string

const (
	CategoryDiagnosis      Category = "diagnosis"        // 확정 진단 발화 (의료법 위반)
	CategoryPrescription   Category = "prescription"     // 처방 권유 (약사법 위반)
	CategorySurgery        Category = "surgery"          // 수술 권유
	CategorySelfHarm       Category = "self_harm"        // 자해 위험
	CategoryDosageChange   Category = "dosage_change"    // 약 복용량 변경
	CategoryEmergency      Category = "emergency"        // 응급 상황 (반드시 119 안내)
	CategoryDangerousAdvice Category = "dangerous_advice" // 위험한 의료 조언
	CategoryHallucination  Category = "hallucination"    // 사실 오류
	CategoryPIIExposure    Category = "pii_exposure"     // PII 노출
)

// Action은 위반 시 취할 조치입니다.
type Action string

const (
	ActionBlock      Action = "block"        // 응답 차단
	ActionRewrite    Action = "rewrite"      // 자동 수정
	ActionWarn       Action = "warn"         // 경고만 (로그)
	ActionEscalate   Action = "escalate"     // 인간 의사 인계
	ActionRequireACK Action = "require_ack"  // 사용자 확인 요구
)

// ============================================================================
// 룰 + 패턴
// ============================================================================

// Rule은 안전 규칙입니다.
type Rule struct {
	ID           string
	Category     Category
	Severity     int            // 1-10
	Action       Action
	Pattern      *regexp.Regexp // 감지 정규식
	RequireKeyword []string     // OR 조건 추가 키워드
	Description  string
	SuggestedFix string         // ActionRewrite 시 대체 문구
}

// 만파식 의료 안전 룰 (베이스라인)
var defaultRules = []*Rule{
	{
		ID: "R-DIAG-001", Category: CategoryDiagnosis, Severity: 9, Action: ActionRewrite,
		Pattern:     regexp.MustCompile(`(확실히|분명히|틀림없이|반드시)\s*(\S+)`),
		Description: "확정 진단 발화 (의료법 위반 가능)",
		SuggestedFix: "전문의 진료 후 정확한 진단을 받으시는 것을 권장드립니다",
	},
	{
		ID: "R-DIAG-002", Category: CategoryDiagnosis, Severity: 9, Action: ActionRewrite,
		Pattern:     regexp.MustCompile(`100%\s*(맞|확실|진단)`),
		Description: "확정성 표현",
		SuggestedFix: "가능성이 있어 보이지만, 정확한 진단을 위해 의료진 상담이 필요합니다",
	},
	{
		ID: "R-PRESC-001", Category: CategoryPrescription, Severity: 9, Action: ActionBlock,
		Pattern:     regexp.MustCompile(`(이|그|저)\s*약(을|울)?\s*(처방|복용|드세요|먹으세요)`),
		Description: "약물 처방 권유 (약사법 위반)",
	},
	{
		ID: "R-DOSAGE-001", Category: CategoryDosageChange, Severity: 8, Action: ActionBlock,
		Pattern:     regexp.MustCompile(`(복용량을?|용량을?)[^.]{0,30}(늘리|줄이|변경|바꾸|두\s*배)`),
		Description: "복용량 임의 변경 권유",
	},
	{
		ID: "R-SURG-001", Category: CategorySurgery, Severity: 8, Action: ActionRewrite,
		Pattern:     regexp.MustCompile(`수술(이|을)?\s*(필요|해야|받으세요)`),
		Description: "수술 권유",
		SuggestedFix: "전문의와 상담하여 치료 옵션을 논의하시는 것을 권장합니다",
	},
	{
		ID: "R-DANGER-001", Category: CategoryDangerousAdvice, Severity: 8, Action: ActionBlock,
		Pattern:     regexp.MustCompile(`(병원|의사)\s*(가지|안가|방문하지)\s*(말|않)`),
		Description: "병원 방문 회피 권유",
	},
	{
		ID: "R-DANGER-002", Category: CategoryDangerousAdvice, Severity: 8, Action: ActionBlock,
		Pattern:     regexp.MustCompile(`수술\s*없이\s*(완치|치료|해결)`),
		Description: "수술 없이 완치 약속",
	},
	{
		ID: "R-EMERG-001", Category: CategoryEmergency, Severity: 10, Action: ActionEscalate,
		Pattern:     regexp.MustCompile(`(의식\s*없|의식이?\s*흐려|호흡\s*곤란|호흡이?\s*힘들|숨이?\s*가쁘|심정지|쓰러|기절|출혈|119)`),
		Description: "응급 상황 키워드 감지 → 즉시 119 안내 + 의사 인계",
	},
	{
		ID: "R-SELF-001", Category: CategorySelfHarm, Severity: 10, Action: ActionEscalate,
		Pattern:     regexp.MustCompile(`(자살|자해|죽고\s*싶|살기\s*싫)`),
		Description: "자해/자살 위험 → 정신건강 위기상담 인계",
	},
	{
		ID: "R-PII-001", Category: CategoryPIIExposure, Severity: 7, Action: ActionRewrite,
		Pattern:     regexp.MustCompile(`\d{6}-\d{7}`), // 주민등록번호
		Description: "주민등록번호 노출",
		SuggestedFix: "[개인정보]",
	},
	{
		ID: "R-PII-002", Category: CategoryPIIExposure, Severity: 6, Action: ActionRewrite,
		Pattern:     regexp.MustCompile(`010-\d{4}-\d{4}|01\d{8,9}`), // 전화번호
		Description: "전화번호 노출",
		SuggestedFix: "[연락처]",
	},
}

// ============================================================================
// 분석 결과
// ============================================================================

// Violation은 감지된 안전 위반 1건입니다.
type Violation struct {
	RuleID      string
	Category    Category
	Severity    int
	Action      Action
	MatchedText string
	StartIdx    int
	EndIdx      int
	Description string
}

// CheckResult는 안전 검사 결과입니다.
type CheckResult struct {
	Safe          bool
	OriginalText  string
	SanitizedText string
	Violations    []*Violation
	HighestAction Action
	HighestSeverity int
	CheckedAt     time.Time
}

// HasBlockedContent는 차단할 내용이 있는지 반환합니다.
func (r *CheckResult) HasBlockedContent() bool {
	for _, v := range r.Violations {
		if v.Action == ActionBlock {
			return true
		}
	}
	return false
}

// RequiresEscalation은 인간 의사 인계가 필요한지 반환합니다.
func (r *CheckResult) RequiresEscalation() bool {
	for _, v := range r.Violations {
		if v.Action == ActionEscalate {
			return true
		}
	}
	return false
}

// CategoriesViolated는 위반된 카테고리 목록을 반환합니다.
func (r *CheckResult) CategoriesViolated() []Category {
	seen := make(map[Category]bool)
	var cats []Category
	for _, v := range r.Violations {
		if !seen[v.Category] {
			seen[v.Category] = true
			cats = append(cats, v.Category)
		}
	}
	return cats
}

// ============================================================================
// PolarisSafetyNet
// ============================================================================

// PolarisSafetyNet은 의료 도메인 안전 검사를 수행합니다.
type PolarisSafetyNet struct {
	mu               sync.RWMutex
	rules            []*Rule
	disclaimer       string
	requireDisclaimer bool
}

// NewPolarisSafetyNet은 새 안전망을 생성합니다.
func NewPolarisSafetyNet() *PolarisSafetyNet {
	return &PolarisSafetyNet{
		rules:             append([]*Rule{}, defaultRules...),
		disclaimer:        "본 응답은 의료 진단을 대체할 수 없으며, 정확한 진단과 치료는 의료진과의 상담을 통해 받으시기 바랍니다.",
		requireDisclaimer: true,
	}
}

// SetDisclaimer는 자동 첨부 면책 조항을 변경합니다.
func (p *PolarisSafetyNet) SetDisclaimer(disclaimer string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.disclaimer = disclaimer
}

// AddRule은 추가 룰을 등록합니다.
func (p *PolarisSafetyNet) AddRule(rule *Rule) error {
	if rule == nil || rule.ID == "" {
		return errors.New("rule and ID required")
	}
	if rule.Pattern == nil {
		return errors.New("pattern required")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.rules = append(p.rules, rule)
	return nil
}

// Check는 텍스트를 모든 룰에 대해 검사합니다.
func (p *PolarisSafetyNet) Check(text string) *CheckResult {
	if text == "" {
		return &CheckResult{
			Safe: true, OriginalText: "", SanitizedText: "",
			CheckedAt: time.Now().UTC(),
		}
	}

	p.mu.RLock()
	rules := append([]*Rule{}, p.rules...)
	p.mu.RUnlock()

	result := &CheckResult{
		OriginalText:  text,
		SanitizedText: text,
		CheckedAt:     time.Now().UTC(),
	}

	for _, rule := range rules {
		matches := rule.Pattern.FindAllStringIndex(text, -1)
		for _, m := range matches {
			matchedText := text[m[0]:m[1]]

			// 추가 키워드 조건 검사
			if len(rule.RequireKeyword) > 0 {
				found := false
				for _, kw := range rule.RequireKeyword {
					if strings.Contains(text, kw) {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}

			result.Violations = append(result.Violations, &Violation{
				RuleID:      rule.ID,
				Category:    rule.Category,
				Severity:    rule.Severity,
				Action:      rule.Action,
				MatchedText: matchedText,
				StartIdx:    m[0],
				EndIdx:      m[1],
				Description: rule.Description,
			})

			if rule.Severity > result.HighestSeverity {
				result.HighestSeverity = rule.Severity
				result.HighestAction = rule.Action
			}
		}
	}

	// 정렬: 심각도 내림차순
	sort.SliceStable(result.Violations, func(i, j int) bool {
		return result.Violations[i].Severity > result.Violations[j].Severity
	})

	result.Safe = !result.HasBlockedContent() && !result.RequiresEscalation()

	// Sanitize: rewrite 룰 적용
	result.SanitizedText = p.sanitize(text, result.Violations, rules)

	// 면책 조항 추가
	if p.requireDisclaimer && result.Safe {
		result.SanitizedText = p.appendDisclaimer(result.SanitizedText)
	}

	return result
}

// sanitize는 ActionRewrite 위반을 자동 대체합니다.
func (p *PolarisSafetyNet) sanitize(text string, violations []*Violation, rules []*Rule) string {
	ruleByID := make(map[string]*Rule)
	for _, r := range rules {
		ruleByID[r.ID] = r
	}

	// 뒤에서부터 치환 (인덱스 보존)
	sorted := append([]*Violation{}, violations...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].StartIdx > sorted[j].StartIdx
	})

	result := text
	for _, v := range sorted {
		if v.Action != ActionRewrite {
			continue
		}
		rule := ruleByID[v.RuleID]
		if rule == nil || rule.SuggestedFix == "" {
			continue
		}
		result = result[:v.StartIdx] + rule.SuggestedFix + result[v.EndIdx:]
	}
	return result
}

func (p *PolarisSafetyNet) appendDisclaimer(text string) string {
	if strings.Contains(text, "[안내]") || strings.Contains(text, p.disclaimer) {
		return text
	}
	return text + "\n\n[안내] " + p.disclaimer
}

// ============================================================================
// 감사 로그 (모든 차단/수정 기록)
// ============================================================================

// AuditEntry는 검사 감사 로그 1건입니다.
type AuditEntry struct {
	ID            string
	UserID        string
	SessionID     string
	OriginalText  string
	SanitizedText string
	Violations    []*Violation
	Action        Action
	CheckedAt     time.Time
}

// AuditLog는 안전 검사 감사 로그입니다.
type AuditLog struct {
	mu      sync.Mutex
	entries []*AuditEntry
	maxSize int
}

// NewAuditLog는 새 감사 로그를 생성합니다.
func NewAuditLog(maxSize int) *AuditLog {
	if maxSize <= 0 {
		maxSize = 10000
	}
	return &AuditLog{maxSize: maxSize}
}

// Record는 검사 결과를 감사 로그에 기록합니다.
func (l *AuditLog) Record(userID, sessionID string, result *CheckResult) {
	if result == nil || len(result.Violations) == 0 {
		return // 위반 없으면 기록하지 않음
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	entry := &AuditEntry{
		ID:            fmt.Sprintf("audit-%d", time.Now().UnixNano()),
		UserID:        userID,
		SessionID:     sessionID,
		OriginalText:  result.OriginalText,
		SanitizedText: result.SanitizedText,
		Violations:    result.Violations,
		Action:        result.HighestAction,
		CheckedAt:     result.CheckedAt,
	}
	l.entries = append(l.entries, entry)

	// 윈도우 초과 시 오래된 항목 제거
	if len(l.entries) > l.maxSize {
		l.entries = l.entries[len(l.entries)-l.maxSize:]
	}
}

// FindBySession은 세션의 감사 로그를 반환합니다.
func (l *AuditLog) FindBySession(sessionID string) []*AuditEntry {
	l.mu.Lock()
	defer l.mu.Unlock()
	var result []*AuditEntry
	for _, e := range l.entries {
		if e.SessionID == sessionID {
			result = append(result, e)
		}
	}
	return result
}

// CountByCategory는 카테고리별 위반 건수를 반환합니다.
func (l *AuditLog) CountByCategory() map[Category]int {
	l.mu.Lock()
	defer l.mu.Unlock()
	counts := make(map[Category]int)
	for _, e := range l.entries {
		for _, v := range e.Violations {
			counts[v.Category]++
		}
	}
	return counts
}

// Count는 총 감사 로그 수를 반환합니다.
func (l *AuditLog) Count() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.entries)
}
