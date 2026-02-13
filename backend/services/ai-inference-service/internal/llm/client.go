// Package llm은 LLM(대규모 언어 모델) 클라이언트를 제공합니다.
//
// OpenAI Chat Completion API를 기반으로 하며, 환경변수 또는 DB 설정으로
// 프로바이더·모델·API 키를 구성할 수 있습니다.
// 보안 원칙: API 키를 코드에 하드코딩하지 않습니다.
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	apperrors "github.com/manpasik/backend/shared/errors"
)

// ============================================================================
// 인터페이스 및 도메인 타입
// ============================================================================

// LLMClient는 LLM과의 채팅 인터페이스입니다.
type LLMClient interface {
	// Chat는 시스템 프롬프트와 메시지 목록을 전달하고 응답을 반환합니다.
	// ctx의 deadline/cancel을 존중합니다.
	Chat(ctx context.Context, systemPrompt string, messages []ChatMessage) (*ChatResponse, error)
}

// ChatMessage는 대화 메시지 하나를 나타냅니다.
type ChatMessage struct {
	Role    string `json:"role"`    // "user", "assistant", "system"
	Content string `json:"content"` // 메시지 본문
}

// ChatResponse는 LLM 응답을 나타냅니다.
type ChatResponse struct {
	Content      string // 모델이 생성한 텍스트
	FinishReason string // "stop", "length" 등
	TokensUsed   int    // 총 토큰 사용량
}

// ============================================================================
// OpenAI 구현
// ============================================================================

const (
	// DefaultModel은 기본 OpenAI 모델입니다.
	DefaultModel = "gpt-4o"

	// DefaultBaseURL은 OpenAI API 기본 엔드포인트입니다.
	DefaultBaseURL = "https://api.openai.com/v1"

	// DefaultTimeout은 HTTP 요청 기본 타임아웃입니다.
	DefaultTimeout = 30 * time.Second
)

// OpenAIClient는 OpenAI Chat Completion API 클라이언트입니다.
type OpenAIClient struct {
	apiKey     string
	model      string
	baseURL    string
	httpClient *http.Client
}

// OpenAIOption은 OpenAIClient 생성 시 옵션을 설정하는 함수 타입입니다.
type OpenAIOption func(*OpenAIClient)

// WithBaseURL은 커스텀 API 기본 URL을 설정합니다.
func WithBaseURL(url string) OpenAIOption {
	return func(c *OpenAIClient) {
		if url != "" {
			c.baseURL = url
		}
	}
}

// WithHTTPClient는 커스텀 HTTP 클라이언트를 설정합니다.
func WithHTTPClient(hc *http.Client) OpenAIOption {
	return func(c *OpenAIClient) {
		if hc != nil {
			c.httpClient = hc
		}
	}
}

// NewOpenAIClient는 새 OpenAI 클라이언트를 생성합니다.
//
// apiKey가 비어있으면 Chat 호출 시 에러를 반환합니다.
// model이 비어있으면 기본값 "gpt-4o"가 사용됩니다.
func NewOpenAIClient(apiKey, model string, opts ...OpenAIOption) *OpenAIClient {
	if model == "" {
		model = DefaultModel
	}

	c := &OpenAIClient{
		apiKey:  apiKey,
		model:   model,
		baseURL: DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// NewOpenAIClientFromEnv는 환경변수에서 설정을 읽어 OpenAI 클라이언트를 생성합니다.
//
// 환경변수:
//   - LLM_API_KEY: API 키 (필수)
//   - LLM_MODEL: 모델 이름 (기본값: gpt-4o)
//   - LLM_BASE_URL: API 기본 URL (기본값: https://api.openai.com/v1)
func NewOpenAIClientFromEnv() *OpenAIClient {
	apiKey := os.Getenv("LLM_API_KEY")
	model := os.Getenv("LLM_MODEL")
	baseURL := os.Getenv("LLM_BASE_URL")

	return NewOpenAIClient(apiKey, model, WithBaseURL(baseURL))
}

// Model은 현재 설정된 모델 이름을 반환합니다.
func (c *OpenAIClient) Model() string {
	return c.model
}

// BaseURL은 현재 설정된 API 기본 URL을 반환합니다.
func (c *OpenAIClient) BaseURL() string {
	return c.baseURL
}

// Chat는 OpenAI Chat Completion API를 호출합니다.
//
// ctx의 deadline/cancel을 존중하며, API 에러·네트워크 타임아웃·응답 파싱 에러를 처리합니다.
func (c *OpenAIClient) Chat(ctx context.Context, systemPrompt string, messages []ChatMessage) (*ChatResponse, error) {
	// API 키 검증
	if c.apiKey == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "LLM API 키가 설정되지 않았습니다")
	}

	// 요청 메시지 구성: system 프롬프트 + 사용자 메시지
	reqMessages := make([]openaiMessage, 0, len(messages)+1)
	if systemPrompt != "" {
		reqMessages = append(reqMessages, openaiMessage{
			Role:    "system",
			Content: systemPrompt,
		})
	}
	for _, m := range messages {
		reqMessages = append(reqMessages, openaiMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}

	reqBody := openaiChatRequest{
		Model:    c.model,
		Messages: reqMessages,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrInternal, "LLM 요청 직렬화에 실패했습니다").
			WithDetails(fmt.Sprintf("json marshal: %v", err))
	}

	// HTTP 요청 생성 (context 전파)
	endpoint := fmt.Sprintf("%s/chat/completions", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, apperrors.New(apperrors.ErrInternal, "LLM HTTP 요청 생성에 실패했습니다").
			WithDetails(fmt.Sprintf("new request: %v", err))
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	// HTTP 호출
	resp, err := c.httpClient.Do(req)
	if err != nil {
		// context 취소/타임아웃 확인
		if ctx.Err() != nil {
			return nil, apperrors.New(apperrors.ErrServiceUnavailable, "LLM 요청이 취소되었거나 타임아웃되었습니다").
				WithDetails(fmt.Sprintf("context: %v", ctx.Err()))
		}
		return nil, apperrors.New(apperrors.ErrServiceUnavailable, "LLM API 호출에 실패했습니다").
			WithDetails(fmt.Sprintf("http do: %v", err))
	}
	defer resp.Body.Close()

	// 응답 바디 읽기
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrInternal, "LLM 응답 읽기에 실패했습니다").
			WithDetails(fmt.Sprintf("read body: %v", err))
	}

	// HTTP 상태 코드 확인
	if resp.StatusCode != http.StatusOK {
		return nil, c.parseAPIError(resp.StatusCode, respBody)
	}

	// 응답 파싱
	var chatResp openaiChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, apperrors.New(apperrors.ErrInternal, "LLM 응답 파싱에 실패했습니다").
			WithDetails(fmt.Sprintf("json unmarshal: %v", err))
	}

	if len(chatResp.Choices) == 0 {
		return nil, apperrors.New(apperrors.ErrInternal, "LLM 응답에 선택지가 없습니다")
	}

	return &ChatResponse{
		Content:      chatResp.Choices[0].Message.Content,
		FinishReason: chatResp.Choices[0].FinishReason,
		TokensUsed:   chatResp.Usage.TotalTokens,
	}, nil
}

// parseAPIError는 OpenAI API 에러 응답을 파싱합니다.
func (c *OpenAIClient) parseAPIError(statusCode int, body []byte) error {
	var apiErr openaiErrorResponse
	if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Error.Message != "" {
		switch statusCode {
		case http.StatusUnauthorized:
			return apperrors.New(apperrors.ErrUnauthorized, "LLM API 인증에 실패했습니다").
				WithDetails(apiErr.Error.Message)
		case http.StatusTooManyRequests:
			return apperrors.New(apperrors.ErrServiceUnavailable, "LLM API 요청 한도를 초과했습니다").
				WithDetails(apiErr.Error.Message)
		default:
			return apperrors.New(apperrors.ErrInternal, "LLM API 오류가 발생했습니다").
				WithDetails(fmt.Sprintf("status=%d: %s", statusCode, apiErr.Error.Message))
		}
	}

	// 파싱 불가 시 일반 에러
	return apperrors.New(apperrors.ErrInternal, "LLM API 오류가 발생했습니다").
		WithDetails(fmt.Sprintf("status=%d, body=%s", statusCode, string(body)))
}

// ============================================================================
// OpenAI API 요청/응답 구조체 (내부용)
// ============================================================================

// openaiMessage는 OpenAI Chat API 메시지입니다.
type openaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// openaiChatRequest는 OpenAI Chat Completion 요청입니다.
type openaiChatRequest struct {
	Model    string           `json:"model"`
	Messages []openaiMessage  `json:"messages"`
}

// openaiChatResponse는 OpenAI Chat Completion 응답입니다.
type openaiChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int           `json:"index"`
		Message      openaiMessage `json:"message"`
		FinishReason string        `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// openaiErrorResponse는 OpenAI API 에러 응답입니다.
type openaiErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

// 컴파일 타임 인터페이스 구현 확인
var _ LLMClient = (*OpenAIClient)(nil)
