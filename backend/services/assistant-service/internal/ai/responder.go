// Package ai는 AI 비서 응답 생성을 위한 외부 API 연동 인터페이스를 제공합니다.
package ai

import "context"

// ChatMessage는 대화 이력의 단일 메시지입니다.
type ChatMessage struct {
	Role    string // "user" 또는 "assistant"
	Content string
}

// AIResponse는 AI 응답 결과입니다.
type AIResponse struct {
	Content      string // 응답 텍스트
	Intent       string // 감지된 인텐트 (선택)
	ActionType   string // "query", "action", "" (선택)
	ActionResult string // JSON 액션 데이터 (선택)
}

// AIResponder는 AI 기반 대화 응답 생성 인터페이스입니다.
type AIResponder interface {
	// GenerateResponse는 대화 이력과 사용자 메시지를 기반으로 응답을 생성합니다.
	GenerateResponse(ctx context.Context, history []ChatMessage, userMessage string) (*AIResponse, error)
}

// NoopResponder는 에러를 반환하여 기존 키워드 기반 응답으로 폴백하도록 유도합니다.
type NoopResponder struct{}

// NewNoopResponder는 NoopResponder를 생성합니다.
func NewNoopResponder() *NoopResponder {
	return &NoopResponder{}
}

// GenerateResponse는 항상 nil을 반환하여 폴백을 유도합니다.
func (n *NoopResponder) GenerateResponse(_ context.Context, _ []ChatMessage, _ string) (*AIResponse, error) {
	return nil, nil
}
