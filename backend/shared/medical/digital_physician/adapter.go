// Package digital_physician은 디지털 주치의(ACE-style 2D 아바타) 어댑터입니다.
//
// P10-09b D-04 차별화의 핵심 인터페이스. 외부 솔루션(NVIDIA ACE, OSS 대체)을
// 추상화하여 음성/표정/응답을 표준 인터페이스로 제어합니다.
//
// SSOT 무결성: 본 패키지는 SW 영역만 다루며, 카메라/마이크 하드웨어 사양은
// Phase 5 SSOT.
package digital_physician

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// ============================================================================
// 도메인 모델
// ============================================================================

// Emotion은 아바타의 감정 표현입니다.
type Emotion string

const (
	EmotionNeutral  Emotion = "neutral"
	EmotionConcern  Emotion = "concern"  // 걱정 (위험 수치 발견 시)
	EmotionWarmth   Emotion = "warmth"   // 따뜻함 (위로/지지)
	EmotionAttention Emotion = "attention" // 주의 (응급 시)
	EmotionEncouragement Emotion = "encouragement" // 격려
)

// Gesture는 아바타 제스처입니다.
type Gesture string

const (
	GestureNod       Gesture = "nod"        // 끄덕임
	GestureHeadTilt  Gesture = "head_tilt"  // 고개 기울임 (경청)
	GestureHandsUp   Gesture = "hands_up"   // 손짓 설명
	GesturePointing  Gesture = "pointing"   // 가리킴 (강조)
	GestureNone      Gesture = "none"
)

// VoiceProfile은 음성 합성 프로필입니다.
type VoiceProfile struct {
	Locale    string  // "ko-KR", "en-US", "ja-JP", "zh-CN"
	Gender    string  // "female", "male", "neutral"
	AgeRange  string  // "young", "adult", "elder"
	Pitch     float64 // 0.5 ~ 1.5 (1.0 = 표준)
	Speed     float64 // 0.5 ~ 2.0 (1.0 = 표준)
	Style     string  // "calm", "energetic", "professional"
}

// DefaultProfile은 기본 음성 프로필 (한국어 여성, 차분한 의료 전문 톤)입니다.
var DefaultProfile = VoiceProfile{
	Locale: "ko-KR", Gender: "female", AgeRange: "adult",
	Pitch: 1.0, Speed: 1.0, Style: "professional",
}

// Utterance는 단일 발화 명령입니다.
type Utterance struct {
	Text       string
	Emotion    Emotion
	Gesture    Gesture
	VoiceProfile *VoiceProfile
	Locale     string // 발화 시점 언어 오버라이드
	PlayBackgroundMusic bool   // 환영/감사 시 BG 사용 여부
	WaitForResponse     bool   // 환자 응답 대기 후 다음 액션
	MaxDurationSec      float64 // 최대 발화 시간 (안전망)
}

// AvatarState는 현재 아바타 상태입니다.
type AvatarState struct {
	SessionID    string
	Emotion      Emotion
	Gesture      Gesture
	IsListening  bool
	IsSpeaking   bool
	CurrentLocale string
	UpdatedAt    time.Time
}

// PlaybackResult는 발화 재생 결과입니다.
type PlaybackResult struct {
	UtteranceID string
	Provider    string
	Status      string // "queued", "playing", "completed", "failed"
	DurationMs  int64
	StartedAt   time.Time
	EndedAt     time.Time
	Error       string
}

// ============================================================================
// Adapter 인터페이스
// ============================================================================

// Adapter는 디지털 주치의 어댑터 인터페이스입니다.
type Adapter interface {
	StartSession(ctx context.Context, sessionID string) error
	SetProfile(ctx context.Context, sessionID string, profile *VoiceProfile) error
	Speak(ctx context.Context, sessionID string, utterance *Utterance) (*PlaybackResult, error)
	SetEmotion(ctx context.Context, sessionID string, emotion Emotion) error
	SetGesture(ctx context.Context, sessionID string, gesture Gesture) error
	GetState(ctx context.Context, sessionID string) (*AvatarState, error)
	EndSession(ctx context.Context, sessionID string) error
	Provider() string
	HealthCheck(ctx context.Context) error
}

// ============================================================================
// 검증
// ============================================================================

var supportedLocales = map[string]bool{
	"ko-KR": true, "en-US": true, "ja-JP": true, "zh-CN": true,
}

// ValidateUtterance는 발화 명령을 검증합니다.
func ValidateUtterance(u *Utterance) error {
	if u == nil {
		return errors.New("utterance is nil")
	}
	if strings.TrimSpace(u.Text) == "" {
		return errors.New("text is empty")
	}
	if len([]rune(u.Text)) > 2000 {
		return errors.New("text exceeds 2000 characters")
	}
	if u.MaxDurationSec < 0 {
		return errors.New("max_duration_sec cannot be negative")
	}
	if u.Locale != "" && !supportedLocales[u.Locale] {
		return fmt.Errorf("unsupported locale: %s", u.Locale)
	}
	return nil
}

// ValidateProfile은 음성 프로필을 검증합니다.
func ValidateProfile(p *VoiceProfile) error {
	if p == nil {
		return errors.New("profile is nil")
	}
	if !supportedLocales[p.Locale] {
		return fmt.Errorf("unsupported locale: %s", p.Locale)
	}
	if p.Pitch < 0.5 || p.Pitch > 1.5 {
		return errors.New("pitch must be in [0.5, 1.5]")
	}
	if p.Speed < 0.5 || p.Speed > 2.0 {
		return errors.New("speed must be in [0.5, 2.0]")
	}
	return nil
}

// ============================================================================
// 팩토리
// ============================================================================

// NewFromEnv는 환경변수에 따라 Adapter를 생성합니다.
//
// DIGITAL_PHYSICIAN_PROVIDER:
//   - "ace": NVIDIA ACE
//   - "" or "noop": 인메모리 (테스트/개발)
func NewFromEnv() Adapter {
	provider := strings.ToLower(os.Getenv("DIGITAL_PHYSICIAN_PROVIDER"))
	switch provider {
	case "ace":
		return NewACEAdapter(
			os.Getenv("NVIDIA_ACE_ENDPOINT"),
			os.Getenv("NVIDIA_ACE_API_KEY"),
		)
	default:
		return NewNoopAdapter()
	}
}

// ============================================================================
// Noop Adapter (테스트/개발)
// ============================================================================

// NoopAdapter는 발화/제스처를 메모리에만 기록하는 어댑터입니다.
type NoopAdapter struct {
	mu       sync.RWMutex
	sessions map[string]*AvatarState
	history  map[string][]*PlaybackResult
}

// NewNoopAdapter는 새 Noop 어댑터를 생성합니다.
func NewNoopAdapter() *NoopAdapter {
	return &NoopAdapter{
		sessions: make(map[string]*AvatarState),
		history:  make(map[string][]*PlaybackResult),
	}
}

// StartSession은 새 세션을 시작합니다.
func (a *NoopAdapter) StartSession(_ context.Context, sessionID string) error {
	if sessionID == "" {
		return errors.New("session_id required")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, exists := a.sessions[sessionID]; exists {
		return fmt.Errorf("session %s already exists", sessionID)
	}
	a.sessions[sessionID] = &AvatarState{
		SessionID:     sessionID,
		Emotion:       EmotionNeutral,
		Gesture:       GestureNone,
		CurrentLocale: "ko-KR",
		UpdatedAt:     time.Now().UTC(),
	}
	return nil
}

// SetProfile은 세션의 음성 프로필을 변경합니다.
func (a *NoopAdapter) SetProfile(_ context.Context, sessionID string, profile *VoiceProfile) error {
	if err := ValidateProfile(profile); err != nil {
		return err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	state, ok := a.sessions[sessionID]
	if !ok {
		return fmt.Errorf("session %s not found", sessionID)
	}
	state.CurrentLocale = profile.Locale
	state.UpdatedAt = time.Now().UTC()
	return nil
}

// Speak는 발화를 시뮬레이션합니다.
//
// ctx 취소를 존중하며, 락은 단일 임계영역으로 처리하여 EndSession과의
// race condition으로 인한 고아 history를 방지합니다.
func (a *NoopAdapter) Speak(ctx context.Context, sessionID string, utterance *Utterance) (*PlaybackResult, error) {
	if err := ValidateUtterance(utterance); err != nil {
		return nil, err
	}
	if ctx != nil && ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// 시뮬레이션 길이 사전 계산
	durationMs := int64(float64(len([]rune(utterance.Text))) / 200.0 * 1000.0)
	if durationMs < 500 {
		durationMs = 500
	}

	now := time.Now().UTC()
	result := &PlaybackResult{
		UtteranceID: fmt.Sprintf("utt-%d", now.UnixNano()),
		Provider:    "noop",
		Status:      "completed",
		DurationMs:  durationMs,
		StartedAt:   now,
		EndedAt:     now.Add(time.Duration(durationMs) * time.Millisecond),
	}

	// 단일 락으로 세션 검증 + 상태 갱신 + history 추가 (atomic)
	a.mu.Lock()
	defer a.mu.Unlock()

	state, ok := a.sessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("session %s not found", sessionID)
	}
	if utterance.Emotion != "" {
		state.Emotion = utterance.Emotion
	}
	if utterance.Gesture != "" {
		state.Gesture = utterance.Gesture
	}
	state.IsSpeaking = false // 시뮬레이션이므로 즉시 완료
	state.UpdatedAt = time.Now().UTC()
	a.history[sessionID] = append(a.history[sessionID], result)
	return result, nil
}

// SetEmotion은 감정을 변경합니다.
func (a *NoopAdapter) SetEmotion(_ context.Context, sessionID string, emotion Emotion) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	state, ok := a.sessions[sessionID]
	if !ok {
		return fmt.Errorf("session %s not found", sessionID)
	}
	state.Emotion = emotion
	state.UpdatedAt = time.Now().UTC()
	return nil
}

// SetGesture는 제스처를 변경합니다.
func (a *NoopAdapter) SetGesture(_ context.Context, sessionID string, gesture Gesture) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	state, ok := a.sessions[sessionID]
	if !ok {
		return fmt.Errorf("session %s not found", sessionID)
	}
	state.Gesture = gesture
	state.UpdatedAt = time.Now().UTC()
	return nil
}

// GetState는 현재 상태를 반환합니다.
func (a *NoopAdapter) GetState(_ context.Context, sessionID string) (*AvatarState, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	state, ok := a.sessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("session %s not found", sessionID)
	}
	clone := *state
	return &clone, nil
}

// EndSession은 세션을 종료합니다.
func (a *NoopAdapter) EndSession(_ context.Context, sessionID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.sessions, sessionID)
	return nil
}

// Provider는 이름을 반환합니다.
func (a *NoopAdapter) Provider() string { return "noop" }

// HealthCheck는 항상 성공합니다.
func (a *NoopAdapter) HealthCheck(_ context.Context) error { return nil }

// History는 세션의 발화 이력을 반환합니다 (테스트용).
func (a *NoopAdapter) History(sessionID string) []*PlaybackResult {
	a.mu.RLock()
	defer a.mu.RUnlock()
	src := a.history[sessionID]
	out := make([]*PlaybackResult, len(src))
	copy(out, src)
	return out
}

// SessionCount는 활성 세션 수를 반환합니다.
func (a *NoopAdapter) SessionCount() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.sessions)
}

// ============================================================================
// ACE Adapter (스텁 — 실 SDK 통합 시 구현)
// ============================================================================

// ACEAdapter는 NVIDIA ACE 어댑터입니다.
//
// 운영에서는 https://developer.nvidia.com/ace SDK 통합 권장.
type ACEAdapter struct {
	endpoint string
	apiKey   string
	noop     *NoopAdapter // 실 SDK 통합 전 폴백
}

// NewACEAdapter는 새 ACE 어댑터를 생성합니다.
func NewACEAdapter(endpoint, apiKey string) *ACEAdapter {
	return &ACEAdapter{
		endpoint: endpoint, apiKey: apiKey,
		noop: NewNoopAdapter(),
	}
}

// Provider는 이름을 반환합니다.
func (a *ACEAdapter) Provider() string { return "nvidia_ace" }

// HealthCheck는 자격증명을 확인합니다.
func (a *ACEAdapter) HealthCheck(_ context.Context) error {
	if a.endpoint == "" || a.apiKey == "" {
		return errors.New("ACE endpoint/api_key not configured")
	}
	return nil
}

// 본 구현은 인메모리 폴백을 활용. 실 SDK 통합 시 noop 호출을 SDK 호출로 교체.
func (a *ACEAdapter) StartSession(ctx context.Context, sessionID string) error {
	return a.noop.StartSession(ctx, sessionID)
}

func (a *ACEAdapter) SetProfile(ctx context.Context, sessionID string, profile *VoiceProfile) error {
	return a.noop.SetProfile(ctx, sessionID, profile)
}

func (a *ACEAdapter) Speak(ctx context.Context, sessionID string, utterance *Utterance) (*PlaybackResult, error) {
	return a.noop.Speak(ctx, sessionID, utterance)
}

func (a *ACEAdapter) SetEmotion(ctx context.Context, sessionID string, emotion Emotion) error {
	return a.noop.SetEmotion(ctx, sessionID, emotion)
}

func (a *ACEAdapter) SetGesture(ctx context.Context, sessionID string, gesture Gesture) error {
	return a.noop.SetGesture(ctx, sessionID, gesture)
}

func (a *ACEAdapter) GetState(ctx context.Context, sessionID string) (*AvatarState, error) {
	return a.noop.GetState(ctx, sessionID)
}

func (a *ACEAdapter) EndSession(ctx context.Context, sessionID string) error {
	return a.noop.EndSession(ctx, sessionID)
}
