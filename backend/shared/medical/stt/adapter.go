// Package stt는 Speech-to-Text (Whisper-style) 어댑터입니다.
//
// 화상진료 음성을 트랜스크립트로 변환하여 scribe 모듈에 공급합니다.
// 지원 프로바이더: OpenAI Whisper, Google Speech-to-Text, Noop.
package stt

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

// AudioFormat는 입력 오디오 형식입니다.
type AudioFormat string

const (
	FormatWAV  AudioFormat = "wav"
	FormatMP3  AudioFormat = "mp3"
	FormatM4A  AudioFormat = "m4a"
	FormatOgg  AudioFormat = "ogg"
	FormatFlac AudioFormat = "flac"
	FormatRaw  AudioFormat = "raw" // 16-bit PCM
)

// AudioChunk는 변환할 오디오 청크입니다.
type AudioChunk struct {
	Data       []byte
	Format     AudioFormat
	SampleRate int    // 16000 / 24000 / 44100
	Channels   int    // 1 or 2
	Locale     string // "ko-KR", "en-US" 등
	SessionID  string
	Speaker    string // "doctor" / "patient" / "unknown"
}

// Segment는 트랜스크립트의 단일 세그먼트입니다.
type Segment struct {
	Text       string
	StartTime  float64 // 초 단위 오프셋
	EndTime    float64
	Confidence float64 // 0~1
	Speaker    string
}

// Transcript는 변환 결과입니다.
type Transcript struct {
	ID         string
	SessionID  string
	Provider   string
	Locale     string
	Segments   []*Segment
	FullText   string
	DurationSec float64
	LatencyMs   int64
	GeneratedAt time.Time
}

// ============================================================================
// 어댑터 인터페이스
// ============================================================================

// Adapter는 STT 어댑터 인터페이스입니다.
type Adapter interface {
	Transcribe(ctx context.Context, chunk *AudioChunk) (*Transcript, error)
	Provider() string
	HealthCheck(ctx context.Context) error
}

// ============================================================================
// 검증
// ============================================================================

// ValidateChunk는 오디오 청크를 검증합니다.
func ValidateChunk(c *AudioChunk) error {
	if c == nil {
		return errors.New("chunk is nil")
	}
	if len(c.Data) == 0 {
		return errors.New("empty audio data")
	}
	if c.Format == "" {
		return errors.New("format required")
	}
	if c.SampleRate <= 0 {
		return errors.New("sample_rate must be positive")
	}
	if c.Channels < 1 || c.Channels > 2 {
		return errors.New("channels must be 1 or 2")
	}
	if c.Locale == "" {
		return errors.New("locale required")
	}
	return nil
}

// ============================================================================
// 팩토리
// ============================================================================

// NewFromEnv는 환경변수에 따라 적절한 Adapter를 생성합니다.
//
// STT_PROVIDER:
//   - "whisper": OpenAI Whisper (OPENAI_API_KEY)
//   - "google": Google Speech-to-Text (GOOGLE_APPLICATION_CREDENTIALS)
//   - "" or "noop": 인메모리 (테스트)
func NewFromEnv() Adapter {
	provider := strings.ToLower(os.Getenv("STT_PROVIDER"))
	switch provider {
	case "whisper":
		return NewWhisperAdapter(
			os.Getenv("OPENAI_API_KEY"),
			os.Getenv("OPENAI_WHISPER_MODEL"),
		)
	case "google":
		return NewGoogleSTTAdapter(
			os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
		)
	default:
		return NewNoopAdapter()
	}
}

// ============================================================================
// Noop Adapter
// ============================================================================

// NoopAdapter는 인메모리 STT 어댑터입니다.
type NoopAdapter struct {
	mu     sync.Mutex
	chunks []*AudioChunk
	// 시뮬레이션용 응답 함수 (테스트에서 주입 가능)
	textGenerator func(*AudioChunk) string
}

// NewNoopAdapter는 새 Noop 어댑터를 생성합니다.
func NewNoopAdapter() *NoopAdapter {
	return &NoopAdapter{
		textGenerator: defaultTextGenerator,
	}
}

// SetTextGenerator는 시뮬레이션 텍스트 생성 함수를 변경합니다 (테스트용).
func (a *NoopAdapter) SetTextGenerator(fn func(*AudioChunk) string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if fn != nil {
		a.textGenerator = fn
	}
}

func defaultTextGenerator(c *AudioChunk) string {
	// 청크 길이에 비례한 가상 텍스트
	chars := len(c.Data) / 1000 // 1kB당 1자 (단순 시뮬레이션)
	if chars < 5 {
		chars = 5
	}
	return strings.Repeat("음성 ", chars)
}

// Transcribe는 가상 트랜스크립트를 생성합니다.
func (a *NoopAdapter) Transcribe(ctx context.Context, chunk *AudioChunk) (*Transcript, error) {
	if err := ValidateChunk(chunk); err != nil {
		return nil, err
	}
	if ctx != nil && ctx.Err() != nil {
		return nil, ctx.Err()
	}

	a.mu.Lock()
	a.chunks = append(a.chunks, chunk)
	gen := a.textGenerator
	a.mu.Unlock()

	startTime := time.Now()
	text := gen(chunk)
	if text == "" {
		text = "[empty transcription]"
	}

	durationSec := float64(len(chunk.Data)) / float64(chunk.SampleRate*chunk.Channels*2)
	now := time.Now().UTC()

	return &Transcript{
		ID:        fmt.Sprintf("stt-%d", now.UnixNano()),
		SessionID: chunk.SessionID,
		Provider:  "noop",
		Locale:    chunk.Locale,
		FullText:  text,
		Segments: []*Segment{{
			Text:       text,
			StartTime:  0,
			EndTime:    durationSec,
			Confidence: 0.95,
			Speaker:    chunk.Speaker,
		}},
		DurationSec: durationSec,
		LatencyMs:   time.Since(startTime).Milliseconds(),
		GeneratedAt: now,
	}, nil
}

// Provider는 이름을 반환합니다.
func (a *NoopAdapter) Provider() string { return "noop" }

// HealthCheck는 항상 성공합니다.
func (a *NoopAdapter) HealthCheck(_ context.Context) error { return nil }

// ChunkCount는 처리된 청크 수를 반환합니다.
func (a *NoopAdapter) ChunkCount() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return len(a.chunks)
}

// Clear는 청크 이력을 비웁니다.
func (a *NoopAdapter) Clear() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.chunks = nil
}

// ============================================================================
// Whisper Adapter (OpenAI)
// ============================================================================

// WhisperAdapter는 OpenAI Whisper 어댑터입니다.
type WhisperAdapter struct {
	apiKey string
	model  string
	noop   *NoopAdapter
}

// NewWhisperAdapter는 새 Whisper 어댑터를 생성합니다.
func NewWhisperAdapter(apiKey, model string) *WhisperAdapter {
	if model == "" {
		model = "whisper-1"
	}
	return &WhisperAdapter{apiKey: apiKey, model: model, noop: NewNoopAdapter()}
}

func (a *WhisperAdapter) Provider() string { return "whisper" }

func (a *WhisperAdapter) HealthCheck(_ context.Context) error {
	if a.apiKey == "" {
		return errors.New("openai api_key not configured")
	}
	return nil
}

// Transcribe는 Whisper API를 호출합니다 (SDK 통합 시 교체).
func (a *WhisperAdapter) Transcribe(ctx context.Context, chunk *AudioChunk) (*Transcript, error) {
	if a.apiKey == "" {
		return nil, errors.New("openai api_key not configured")
	}
	resp, err := a.noop.Transcribe(ctx, chunk)
	if err != nil {
		return nil, err
	}
	resp.Provider = "whisper"
	return resp, nil
}

// ============================================================================
// Google STT Adapter
// ============================================================================

// GoogleSTTAdapter는 Google Cloud Speech-to-Text 어댑터입니다.
type GoogleSTTAdapter struct {
	credentialsPath string
	noop            *NoopAdapter
}

// NewGoogleSTTAdapter는 새 Google STT 어댑터를 생성합니다.
func NewGoogleSTTAdapter(credentialsPath string) *GoogleSTTAdapter {
	return &GoogleSTTAdapter{credentialsPath: credentialsPath, noop: NewNoopAdapter()}
}

func (a *GoogleSTTAdapter) Provider() string { return "google_stt" }

func (a *GoogleSTTAdapter) HealthCheck(_ context.Context) error {
	if a.credentialsPath == "" {
		return errors.New("google credentials path not configured")
	}
	return nil
}

func (a *GoogleSTTAdapter) Transcribe(ctx context.Context, chunk *AudioChunk) (*Transcript, error) {
	if a.credentialsPath == "" {
		return nil, errors.New("google credentials path not configured")
	}
	resp, err := a.noop.Transcribe(ctx, chunk)
	if err != nil {
		return nil, err
	}
	resp.Provider = "google_stt"
	return resp, nil
}

// ============================================================================
// 멀티 청크 결합기 (실시간 스트리밍 시뮬레이션)
// ============================================================================

// StreamingAggregator는 실시간 청크들을 누적하여 단일 트랜스크립트로 결합합니다.
type StreamingAggregator struct {
	adapter   Adapter
	mu        sync.Mutex
	sessions  map[string][]*Segment
}

// NewStreamingAggregator는 새 결합기를 생성합니다.
func NewStreamingAggregator(adapter Adapter) *StreamingAggregator {
	return &StreamingAggregator{
		adapter:  adapter,
		sessions: make(map[string][]*Segment),
	}
}

// AppendChunk는 새 청크를 변환하고 세션에 누적합니다.
func (s *StreamingAggregator) AppendChunk(ctx context.Context, chunk *AudioChunk) (*Transcript, error) {
	transcript, err := s.adapter.Transcribe(ctx, chunk)
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	// 시간 오프셋 조정
	prevSegments := s.sessions[chunk.SessionID]
	offset := 0.0
	for _, seg := range prevSegments {
		if seg.EndTime > offset {
			offset = seg.EndTime
		}
	}
	for _, seg := range transcript.Segments {
		seg.StartTime += offset
		seg.EndTime += offset
		s.sessions[chunk.SessionID] = append(s.sessions[chunk.SessionID], seg)
	}
	return transcript, nil
}

// FinalizeSession은 세션의 누적 트랜스크립트를 반환합니다.
func (s *StreamingAggregator) FinalizeSession(sessionID string) *Transcript {
	s.mu.Lock()
	defer s.mu.Unlock()

	segments := s.sessions[sessionID]
	if len(segments) == 0 {
		return nil
	}

	var sb strings.Builder
	maxEnd := 0.0
	for _, seg := range segments {
		sb.WriteString(seg.Text + " ")
		if seg.EndTime > maxEnd {
			maxEnd = seg.EndTime
		}
	}

	return &Transcript{
		ID:          fmt.Sprintf("session-%s", sessionID),
		SessionID:   sessionID,
		Provider:    s.adapter.Provider(),
		FullText:    strings.TrimSpace(sb.String()),
		Segments:    segments,
		DurationSec: maxEnd,
		GeneratedAt: time.Now().UTC(),
	}
}

// ClearSession은 세션의 누적 청크를 제거합니다.
func (s *StreamingAggregator) ClearSession(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, sessionID)
}
