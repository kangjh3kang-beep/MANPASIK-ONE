// Package ai: medical/llm.Adapter ↔ AIResponder 브릿지
//
// 기존 자체 OpenAIResponder를 medical/llm 어댑터로 점진 마이그레이션합니다.
// FailoverChain을 활용하여 LLM 프로바이더 장애 시 자동 폴백 가능.
package ai

import (
	"context"
	"errors"

	"github.com/manpasik/backend/shared/medical/llm"
)

// MedicalLLMResponder는 medical/llm.Adapter를 AIResponder 인터페이스로 래핑합니다.
//
// 장점:
//   - LLM 프로바이더 교체 (OpenAI/Claude/Gemini) 가능
//   - FailoverChain으로 장애 자동 폴백
//   - safety/Polaris 안전망과 일관 통합
type MedicalLLMResponder struct {
	adapter      llm.Adapter
	defaultModel string
	systemPrompt string
}

// MedicalLLMConfig는 브릿지 설정입니다.
type MedicalLLMConfig struct {
	Adapter      llm.Adapter // medical/llm.NewFromEnv() 결과 주입
	DefaultModel string      // "gpt-4", "claude-sonnet-4-6" 등 (provider 의존)
	SystemPrompt string      // 의료 도메인 시스템 프롬프트
}

// NewMedicalLLMResponder는 브릿지를 생성합니다.
func NewMedicalLLMResponder(cfg MedicalLLMConfig) *MedicalLLMResponder {
	if cfg.Adapter == nil {
		cfg.Adapter = llm.NewNoopAdapter()
	}
	if cfg.SystemPrompt == "" {
		cfg.SystemPrompt = "당신은 만파식 건강 비서입니다. 의료 진단을 내리지 않으며, 사용자에게 친절하고 정확한 일반 건강 정보를 제공합니다."
	}
	return &MedicalLLMResponder{
		adapter:      cfg.Adapter,
		defaultModel: cfg.DefaultModel,
		systemPrompt: cfg.SystemPrompt,
	}
}

// GenerateResponse는 대화 이력과 사용자 메시지로 응답을 생성합니다.
func (r *MedicalLLMResponder) GenerateResponse(
	ctx context.Context,
	history []ChatMessage,
	userMessage string,
) (*AIResponse, error) {
	if r.adapter == nil {
		return nil, errors.New("llm adapter not configured")
	}

	// AIResponder.history → llm.Message 변환
	messages := make([]*llm.Message, 0, len(history)+2)
	messages = append(messages, &llm.Message{
		Role:    llm.RoleSystem,
		Content: r.systemPrompt,
	})
	for _, h := range history {
		role := llm.RoleAssistant
		if h.Role == "user" {
			role = llm.RoleUser
		}
		messages = append(messages, &llm.Message{
			Role:    role,
			Content: h.Content,
		})
	}
	messages = append(messages, &llm.Message{
		Role:    llm.RoleUser,
		Content: userMessage,
	})

	req := &llm.Request{
		Model:             r.defaultModel,
		Temperature:       0.5, // 의료 도메인은 보수적
		MaxTokens:         800,
		Messages:          messages,
		RequireDisclaimer: true,
	}

	resp, err := r.adapter.Complete(ctx, req)
	if err != nil {
		return nil, err
	}

	return &AIResponse{
		Content: resp.Content,
		// Intent/ActionType은 nlp-service에서 별도 분석
		Intent: "general_query",
	}, nil
}

// Provider는 현재 LLM 프로바이더 이름을 반환합니다.
func (r *MedicalLLMResponder) Provider() string {
	if r.adapter == nil {
		return "unknown"
	}
	return r.adapter.Provider()
}

// HealthCheck는 LLM 어댑터의 헬스체크를 위임합니다.
func (r *MedicalLLMResponder) HealthCheck(ctx context.Context) error {
	if r.adapter == nil {
		return errors.New("adapter not configured")
	}
	return r.adapter.HealthCheck(ctx)
}
