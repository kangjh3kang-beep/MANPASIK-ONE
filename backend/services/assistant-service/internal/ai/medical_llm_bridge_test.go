package ai_test

import (
	"context"
	"strings"
	"testing"

	"github.com/manpasik/backend/services/assistant-service/internal/ai"
	"github.com/manpasik/backend/shared/medical/llm"
)

func TestMedicalLLMResponder_NoopAdapter(t *testing.T) {
	r := ai.NewMedicalLLMResponder(ai.MedicalLLMConfig{
		Adapter: llm.NewNoopAdapter(),
	})

	resp, err := r.GenerateResponse(context.Background(),
		[]ai.ChatMessage{{Role: "user", Content: "안녕하세요"}},
		"두통이 있어요")
	if err != nil {
		t.Fatalf("GenerateResponse 실패: %v", err)
	}
	if resp.Content == "" {
		t.Error("Content 비어 있음")
	}
}

func TestMedicalLLMResponder_DisclaimerEmbedded(t *testing.T) {
	r := ai.NewMedicalLLMResponder(ai.MedicalLLMConfig{
		Adapter: llm.NewNoopAdapter(),
	})

	resp, _ := r.GenerateResponse(context.Background(), nil, "약을 처방해주세요")
	if !strings.Contains(resp.Content, "[안내]") {
		t.Error("의료 면책 미포함")
	}
}

func TestMedicalLLMResponder_DefaultSystemPrompt(t *testing.T) {
	r := ai.NewMedicalLLMResponder(ai.MedicalLLMConfig{}) // 모든 옵션 기본값
	resp, err := r.GenerateResponse(context.Background(), nil, "테스트")
	if err != nil {
		t.Fatalf("기본 설정으로 실패: %v", err)
	}
	if resp == nil {
		t.Fatal("응답 nil")
	}
}

func TestMedicalLLMResponder_HistoryConversion(t *testing.T) {
	noopLLM := llm.NewNoopAdapter()
	r := ai.NewMedicalLLMResponder(ai.MedicalLLMConfig{Adapter: noopLLM})

	history := []ai.ChatMessage{
		{Role: "user", Content: "처음 메시지"},
		{Role: "assistant", Content: "첫 응답"},
		{Role: "user", Content: "두 번째 질문"},
	}

	_, err := r.GenerateResponse(context.Background(), history, "현재 질문")
	if err != nil {
		t.Fatal(err)
	}

	last := noopLLM.LastRequest()
	if last == nil {
		t.Fatal("LastRequest nil")
	}
	// system + 3 history + 1 current = 5
	if len(last.Messages) != 5 {
		t.Errorf("Messages = %d, want 5", len(last.Messages))
	}
	if last.Messages[0].Role != llm.RoleSystem {
		t.Errorf("first role = %q, want system", last.Messages[0].Role)
	}
}

func TestMedicalLLMResponder_Provider(t *testing.T) {
	r := ai.NewMedicalLLMResponder(ai.MedicalLLMConfig{Adapter: llm.NewNoopAdapter()})
	if r.Provider() != "noop" {
		t.Errorf("Provider = %q", r.Provider())
	}
}

func TestMedicalLLMResponder_HealthCheck(t *testing.T) {
	r := ai.NewMedicalLLMResponder(ai.MedicalLLMConfig{Adapter: llm.NewNoopAdapter()})
	if err := r.HealthCheck(context.Background()); err != nil {
		t.Errorf("HealthCheck = %v", err)
	}
}

func TestMedicalLLMResponder_FailoverChain(t *testing.T) {
	// 다중 어댑터 폴백 통합 검증
	primary := llm.NewNoopAdapter()
	chain := llm.NewFailoverChain(primary)
	r := ai.NewMedicalLLMResponder(ai.MedicalLLMConfig{Adapter: chain})

	resp, err := r.GenerateResponse(context.Background(), nil, "테스트")
	if err != nil {
		t.Fatalf("Failover 통합 실패: %v", err)
	}
	if resp.Content == "" {
		t.Error("Failover 응답 비어있음")
	}
}

func TestMedicalLLMResponder_CustomSystemPrompt(t *testing.T) {
	noopLLM := llm.NewNoopAdapter()
	r := ai.NewMedicalLLMResponder(ai.MedicalLLMConfig{
		Adapter:      noopLLM,
		SystemPrompt: "Custom medical prompt",
	})

	_, _ = r.GenerateResponse(context.Background(), nil, "x")

	last := noopLLM.LastRequest()
	if last.Messages[0].Content != "Custom medical prompt" {
		t.Errorf("System prompt = %q", last.Messages[0].Content)
	}
}
