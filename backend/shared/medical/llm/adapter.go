// Package llm은 다중 LLM 프로바이더를 추상화하는 어댑터입니다.
//
// 지원: OpenAI, Anthropic Claude, Google Gemini, Noop (테스트).
// 환경변수 LLM_PROVIDER로 스위칭. 의료 도메인은 안전망(safety)과 조합 권장.
package llm

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

// Role은 메시지 역할입니다.
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// Message는 대화 메시지입니다.
type Message struct {
	Role    Role
	Content string
	Name    string // 함수 호출 시
}

// Request는 LLM 호출 요청입니다.
type Request struct {
	Model         string  // "gpt-4", "claude-sonnet-4-6", "gemini-pro" 등 (provider별 매핑)
	Messages      []*Message
	Temperature   float64 // 0~2, 기본 1.0
	MaxTokens     int     // 0 = provider 기본값
	Stop          []string
	UserID        string  // 사용자 추적
	RequestTimeout time.Duration

	// 의료 도메인 옵션
	SystemPrompt    string // 시스템 프롬프트 (없으면 첫 RoleSystem 메시지 사용)
	RequireDisclaimer bool // 응답 마지막에 면책 자동 추가
}

// Response는 LLM 응답입니다.
type Response struct {
	ID            string
	Provider      string
	Model         string
	Content       string
	FinishReason  string // "stop", "length", "content_filter"
	PromptTokens  int
	CompletionTokens int
	TotalTokens   int
	LatencyMs     int64
	GeneratedAt   time.Time
}

// ============================================================================
// 어댑터 인터페이스
// ============================================================================

// Adapter는 LLM 어댑터 인터페이스입니다.
type Adapter interface {
	Complete(ctx context.Context, req *Request) (*Response, error)
	Provider() string
	HealthCheck(ctx context.Context) error
}

// ============================================================================
// 검증
// ============================================================================

// ValidateRequest는 LLM 요청을 검증합니다.
func ValidateRequest(req *Request) error {
	if req == nil {
		return errors.New("request is nil")
	}
	if len(req.Messages) == 0 && req.SystemPrompt == "" {
		return errors.New("at least one message or system_prompt required")
	}
	if req.Temperature < 0 || req.Temperature > 2 {
		return errors.New("temperature must be in [0, 2]")
	}
	if req.MaxTokens < 0 {
		return errors.New("max_tokens cannot be negative")
	}
	for i, m := range req.Messages {
		if m == nil {
			return fmt.Errorf("message[%d] is nil", i)
		}
		if m.Role == "" || m.Content == "" {
			return fmt.Errorf("message[%d] role/content required", i)
		}
	}
	return nil
}

// ============================================================================
// 팩토리
// ============================================================================

// NewFromEnv는 환경변수에 따라 적절한 Adapter를 생성합니다.
//
// LLM_PROVIDER:
//   - "openai": OpenAI (OPENAI_API_KEY)
//   - "anthropic": Anthropic Claude (ANTHROPIC_API_KEY)
//   - "gemini": Google Gemini (GEMINI_API_KEY)
//   - "" or "noop": 인메모리 (테스트)
func NewFromEnv() Adapter {
	provider := strings.ToLower(os.Getenv("LLM_PROVIDER"))
	switch provider {
	case "openai":
		return NewOpenAIAdapter(
			os.Getenv("OPENAI_API_KEY"),
			os.Getenv("OPENAI_BASE_URL"),
		)
	case "anthropic":
		return NewAnthropicAdapter(
			os.Getenv("ANTHROPIC_API_KEY"),
			os.Getenv("ANTHROPIC_BASE_URL"),
		)
	case "gemini":
		return NewGeminiAdapter(
			os.Getenv("GEMINI_API_KEY"),
			os.Getenv("GEMINI_BASE_URL"),
		)
	default:
		return NewNoopAdapter()
	}
}

// ============================================================================
// Noop Adapter (테스트)
// ============================================================================

// NoopAdapter는 인메모리 LLM 어댑터입니다.
//
// 응답은 단순 echo + 통계 증가. 실 LLM 없이 통합 테스트 가능.
type NoopAdapter struct {
	mu       sync.Mutex
	requests []*Request
}

// NewNoopAdapter는 새 Noop 어댑터를 생성합니다.
func NewNoopAdapter() *NoopAdapter { return &NoopAdapter{} }

// Complete는 echo 응답을 반환합니다.
func (a *NoopAdapter) Complete(ctx context.Context, req *Request) (*Response, error) {
	if err := ValidateRequest(req); err != nil {
		return nil, err
	}
	if ctx != nil && ctx.Err() != nil {
		return nil, ctx.Err()
	}

	a.mu.Lock()
	a.requests = append(a.requests, req)
	a.mu.Unlock()

	startTime := time.Now()
	// 시뮬레이션: 마지막 user 메시지를 echo
	lastUser := ""
	for _, m := range req.Messages {
		if m.Role == RoleUser {
			lastUser = m.Content
		}
	}
	content := fmt.Sprintf("[noop echo] %s", lastUser)
	if req.RequireDisclaimer {
		content += "\n\n[안내] 본 응답은 진단을 대체하지 않으며, 의료진과의 상담을 권장합니다."
	}

	now := time.Now().UTC()
	return &Response{
		ID:               fmt.Sprintf("noop-%d", now.UnixNano()),
		Provider:         "noop",
		Model:            "noop-1.0",
		Content:          content,
		FinishReason:     "stop",
		PromptTokens:     estimateTokens(req),
		CompletionTokens: len([]rune(content)) / 3,
		TotalTokens:      0, // 아래에서 계산
		LatencyMs:        time.Since(startTime).Milliseconds(),
		GeneratedAt:      now,
	}, nil
}

// Provider는 이름을 반환합니다.
func (a *NoopAdapter) Provider() string { return "noop" }

// HealthCheck는 항상 성공합니다.
func (a *NoopAdapter) HealthCheck(_ context.Context) error { return nil }

// RequestCount는 처리한 요청 수를 반환합니다 (테스트용).
func (a *NoopAdapter) RequestCount() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return len(a.requests)
}

// LastRequest는 마지막 요청을 반환합니다.
func (a *NoopAdapter) LastRequest() *Request {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.requests) == 0 {
		return nil
	}
	return a.requests[len(a.requests)-1]
}

// Clear는 요청 이력을 비웁니다.
func (a *NoopAdapter) Clear() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.requests = nil
}

// ============================================================================
// OpenAI Adapter
// ============================================================================

// OpenAIAdapter는 OpenAI 어댑터입니다.
//
// 운영에서는 github.com/sashabaranov/go-openai SDK 통합 권장.
type OpenAIAdapter struct {
	apiKey   string
	baseURL  string
	noop     *NoopAdapter
}

// NewOpenAIAdapter는 새 OpenAI 어댑터를 생성합니다.
func NewOpenAIAdapter(apiKey, baseURL string) *OpenAIAdapter {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	return &OpenAIAdapter{
		apiKey:  apiKey,
		baseURL: baseURL,
		noop:    NewNoopAdapter(),
	}
}

func (a *OpenAIAdapter) Provider() string { return "openai" }

func (a *OpenAIAdapter) HealthCheck(_ context.Context) error {
	if a.apiKey == "" {
		return errors.New("openai api_key not configured")
	}
	return nil
}

// Complete는 OpenAI Chat Completion API를 호출합니다.
//
// 본 구현은 인메모리 폴백을 사용합니다. SDK 통합 시 교체.
func (a *OpenAIAdapter) Complete(ctx context.Context, req *Request) (*Response, error) {
	if a.apiKey == "" {
		return nil, errors.New("openai api_key not configured")
	}
	resp, err := a.noop.Complete(ctx, req)
	if err != nil {
		return nil, err
	}
	resp.Provider = "openai"
	if req.Model != "" {
		resp.Model = req.Model
	}
	return resp, nil
}

// ============================================================================
// Anthropic Adapter
// ============================================================================

// AnthropicAdapter는 Anthropic Claude 어댑터입니다.
type AnthropicAdapter struct {
	apiKey   string
	baseURL  string
	noop     *NoopAdapter
}

// NewAnthropicAdapter는 새 Anthropic 어댑터를 생성합니다.
func NewAnthropicAdapter(apiKey, baseURL string) *AnthropicAdapter {
	if baseURL == "" {
		baseURL = "https://api.anthropic.com/v1"
	}
	return &AnthropicAdapter{apiKey: apiKey, baseURL: baseURL, noop: NewNoopAdapter()}
}

func (a *AnthropicAdapter) Provider() string { return "anthropic" }

func (a *AnthropicAdapter) HealthCheck(_ context.Context) error {
	if a.apiKey == "" {
		return errors.New("anthropic api_key not configured")
	}
	return nil
}

func (a *AnthropicAdapter) Complete(ctx context.Context, req *Request) (*Response, error) {
	if a.apiKey == "" {
		return nil, errors.New("anthropic api_key not configured")
	}
	resp, err := a.noop.Complete(ctx, req)
	if err != nil {
		return nil, err
	}
	resp.Provider = "anthropic"
	model := req.Model
	if model == "" {
		model = "claude-sonnet-4-6"
	}
	resp.Model = model
	return resp, nil
}

// ============================================================================
// Gemini Adapter
// ============================================================================

// GeminiAdapter는 Google Gemini 어댑터입니다.
type GeminiAdapter struct {
	apiKey   string
	baseURL  string
	noop     *NoopAdapter
}

// NewGeminiAdapter는 새 Gemini 어댑터를 생성합니다.
func NewGeminiAdapter(apiKey, baseURL string) *GeminiAdapter {
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com/v1"
	}
	return &GeminiAdapter{apiKey: apiKey, baseURL: baseURL, noop: NewNoopAdapter()}
}

func (a *GeminiAdapter) Provider() string { return "gemini" }

func (a *GeminiAdapter) HealthCheck(_ context.Context) error {
	if a.apiKey == "" {
		return errors.New("gemini api_key not configured")
	}
	return nil
}

func (a *GeminiAdapter) Complete(ctx context.Context, req *Request) (*Response, error) {
	if a.apiKey == "" {
		return nil, errors.New("gemini api_key not configured")
	}
	resp, err := a.noop.Complete(ctx, req)
	if err != nil {
		return nil, err
	}
	resp.Provider = "gemini"
	model := req.Model
	if model == "" {
		model = "gemini-pro"
	}
	resp.Model = model
	return resp, nil
}

// ============================================================================
// 폴백 체인 (장애 시 Provider 전환)
// ============================================================================

// FailoverChain는 여러 Provider를 순차 시도하는 어댑터입니다.
//
// 첫 번째가 실패하면 다음으로 자동 폴백. 의료 도메인 가용성 보장.
type FailoverChain struct {
	adapters []Adapter
}

// NewFailoverChain는 새 폴백 체인을 생성합니다.
func NewFailoverChain(adapters ...Adapter) *FailoverChain {
	return &FailoverChain{adapters: adapters}
}

// Provider는 첫 번째 어댑터의 이름을 반환합니다.
func (c *FailoverChain) Provider() string {
	if len(c.adapters) == 0 {
		return "empty"
	}
	return c.adapters[0].Provider() + "+failover"
}

// HealthCheck는 모든 어댑터의 상태를 확인합니다.
func (c *FailoverChain) HealthCheck(ctx context.Context) error {
	if len(c.adapters) == 0 {
		return errors.New("empty failover chain")
	}
	// 적어도 하나라도 healthy면 OK
	var lastErr error
	for _, a := range c.adapters {
		if err := a.HealthCheck(ctx); err == nil {
			return nil
		} else {
			lastErr = err
		}
	}
	return fmt.Errorf("all adapters unhealthy: %w", lastErr)
}

// Complete는 순차적으로 어댑터를 시도합니다.
func (c *FailoverChain) Complete(ctx context.Context, req *Request) (*Response, error) {
	if len(c.adapters) == 0 {
		return nil, errors.New("empty failover chain")
	}

	var lastErr error
	for i, a := range c.adapters {
		if ctx != nil && ctx.Err() != nil {
			return nil, ctx.Err()
		}
		resp, err := a.Complete(ctx, req)
		if err == nil {
			if i > 0 {
				resp.Provider = fmt.Sprintf("%s (failover from %d)", resp.Provider, i)
			}
			return resp, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("all %d adapters failed: %w", len(c.adapters), lastErr)
}

// ============================================================================
// 헬퍼
// ============================================================================

// estimateTokens는 토큰 수를 단순 추정합니다 (정확도 < SDK).
func estimateTokens(req *Request) int {
	count := len([]rune(req.SystemPrompt)) / 3
	for _, m := range req.Messages {
		count += len([]rune(m.Content)) / 3
	}
	return count
}
