package llm_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/manpasik/backend/shared/medical/llm"
)

func newReq() *llm.Request {
	return &llm.Request{
		Model:       "test-model",
		Temperature: 0.7,
		Messages: []*llm.Message{
			{Role: llm.RoleSystem, Content: "You are a medical assistant."},
			{Role: llm.RoleUser, Content: "어디가 불편하신가요?"},
		},
	}
}

func TestNoopAdapter_Complete(t *testing.T) {
	a := llm.NewNoopAdapter()
	resp, err := a.Complete(context.Background(), newReq())
	if err != nil {
		t.Fatalf("Complete 실패: %v", err)
	}
	if resp.Provider != "noop" {
		t.Errorf("Provider = %q", resp.Provider)
	}
	if resp.FinishReason != "stop" {
		t.Errorf("FinishReason = %q", resp.FinishReason)
	}
	if !strings.Contains(resp.Content, "어디가 불편") {
		t.Errorf("Content echo 실패: %q", resp.Content)
	}
	if a.RequestCount() != 1 {
		t.Errorf("RequestCount = %d", a.RequestCount())
	}
}

func TestNoopAdapter_RequireDisclaimer(t *testing.T) {
	a := llm.NewNoopAdapter()
	req := newReq()
	req.RequireDisclaimer = true

	resp, _ := a.Complete(context.Background(), req)
	if !strings.Contains(resp.Content, "[안내]") {
		t.Error("면책 조항 미추가")
	}
}

func TestValidateRequest(t *testing.T) {
	cases := []struct {
		name    string
		req     *llm.Request
		wantErr bool
	}{
		{"nil", nil, true},
		{"empty messages no system", &llm.Request{Temperature: 0.5}, true},
		{"valid messages", newReq(), false},
		{"system prompt only", &llm.Request{SystemPrompt: "x", Temperature: 0.5}, false},
		{"temperature too high", &llm.Request{
			Messages: []*llm.Message{{Role: llm.RoleUser, Content: "x"}},
			Temperature: 3.0,
		}, true},
		{"negative max_tokens", &llm.Request{
			Messages: []*llm.Message{{Role: llm.RoleUser, Content: "x"}},
			MaxTokens: -1,
		}, true},
		{"empty role", &llm.Request{
			Messages: []*llm.Message{{Content: "x"}},
		}, true},
		{"nil message", &llm.Request{
			Messages: []*llm.Message{nil},
		}, true},
	}
	for _, c := range cases {
		err := llm.ValidateRequest(c.req)
		if (err != nil) != c.wantErr {
			t.Errorf("%s: err=%v, wantErr=%v", c.name, err, c.wantErr)
		}
	}
}

func TestNoopAdapter_ContextCancelled(t *testing.T) {
	a := llm.NewNoopAdapter()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := a.Complete(ctx, newReq())
	if err == nil {
		t.Error("취소된 ctx에서 통과")
	}
}

func TestNoopAdapter_Clear(t *testing.T) {
	a := llm.NewNoopAdapter()
	for i := 0; i < 3; i++ {
		_, _ = a.Complete(context.Background(), newReq())
	}
	a.Clear()
	if a.RequestCount() != 0 {
		t.Errorf("Clear 후 Count = %d", a.RequestCount())
	}
}

func TestNoopAdapter_LastRequest(t *testing.T) {
	a := llm.NewNoopAdapter()
	if a.LastRequest() != nil {
		t.Error("초기 LastRequest != nil")
	}

	req := newReq()
	_, _ = a.Complete(context.Background(), req)

	last := a.LastRequest()
	if last == nil {
		t.Fatal("LastRequest nil")
	}
}

func TestOpenAIAdapter_HealthCheck(t *testing.T) {
	a := llm.NewOpenAIAdapter("api-key", "")
	if err := a.HealthCheck(context.Background()); err != nil {
		t.Errorf("HealthCheck 실패: %v", err)
	}

	a2 := llm.NewOpenAIAdapter("", "")
	if err := a2.HealthCheck(context.Background()); err == nil {
		t.Error("API 키 없이 통과")
	}
}

func TestOpenAIAdapter_CompleteRequiresAPIKey(t *testing.T) {
	a := llm.NewOpenAIAdapter("", "")
	_, err := a.Complete(context.Background(), newReq())
	if err == nil {
		t.Error("API 키 없이 통과")
	}
}

func TestOpenAIAdapter_CompleteSuccess(t *testing.T) {
	a := llm.NewOpenAIAdapter("test-key", "")
	resp, err := a.Complete(context.Background(), newReq())
	if err != nil {
		t.Fatalf("Complete 실패: %v", err)
	}
	if resp.Provider != "openai" {
		t.Errorf("Provider = %q", resp.Provider)
	}
	if resp.Model != "test-model" {
		t.Errorf("Model = %q", resp.Model)
	}
}

func TestAnthropicAdapter_DefaultModel(t *testing.T) {
	a := llm.NewAnthropicAdapter("key", "")
	req := newReq()
	req.Model = "" // 모델 미지정

	resp, err := a.Complete(context.Background(), req)
	if err != nil {
		t.Fatalf("Complete 실패: %v", err)
	}
	if resp.Model != "claude-sonnet-4-6" {
		t.Errorf("기본 Model = %q, want claude-sonnet-4-6", resp.Model)
	}
}

func TestGeminiAdapter_DefaultModel(t *testing.T) {
	a := llm.NewGeminiAdapter("key", "")
	req := newReq()
	req.Model = ""

	resp, _ := a.Complete(context.Background(), req)
	if resp.Model != "gemini-pro" {
		t.Errorf("기본 Model = %q", resp.Model)
	}
}

func TestNewFromEnv_Default(t *testing.T) {
	t.Setenv("LLM_PROVIDER", "")
	a := llm.NewFromEnv()
	if a.Provider() != "noop" {
		t.Errorf("Provider = %q, want noop", a.Provider())
	}
}

func TestNewFromEnv_OpenAI(t *testing.T) {
	t.Setenv("LLM_PROVIDER", "openai")
	t.Setenv("OPENAI_API_KEY", "k")
	a := llm.NewFromEnv()
	if a.Provider() != "openai" {
		t.Errorf("Provider = %q", a.Provider())
	}
}

func TestNewFromEnv_Anthropic(t *testing.T) {
	t.Setenv("LLM_PROVIDER", "anthropic")
	a := llm.NewFromEnv()
	if a.Provider() != "anthropic" {
		t.Errorf("Provider = %q", a.Provider())
	}
}

func TestNewFromEnv_Gemini(t *testing.T) {
	t.Setenv("LLM_PROVIDER", "gemini")
	a := llm.NewFromEnv()
	if a.Provider() != "gemini" {
		t.Errorf("Provider = %q", a.Provider())
	}
}

func TestFailoverChain_FirstSucceeds(t *testing.T) {
	primary := llm.NewNoopAdapter()
	secondary := llm.NewNoopAdapter()
	chain := llm.NewFailoverChain(primary, secondary)

	resp, err := chain.Complete(context.Background(), newReq())
	if err != nil {
		t.Fatalf("Complete 실패: %v", err)
	}
	if primary.RequestCount() != 1 {
		t.Errorf("primary requests = %d, want 1", primary.RequestCount())
	}
	if secondary.RequestCount() != 0 {
		t.Errorf("secondary requests = %d, want 0 (primary 성공)", secondary.RequestCount())
	}
	if resp == nil {
		t.Fatal("resp nil")
	}
}

type failingAdapter struct{}

func (failingAdapter) Provider() string { return "failing" }
func (failingAdapter) HealthCheck(_ context.Context) error { return errors.New("down") }
func (failingAdapter) Complete(_ context.Context, _ *llm.Request) (*llm.Response, error) {
	return nil, errors.New("simulated failure")
}

func TestFailoverChain_PrimaryFails_FallsBack(t *testing.T) {
	failing := failingAdapter{}
	healthy := llm.NewNoopAdapter()
	chain := llm.NewFailoverChain(failing, healthy)

	resp, err := chain.Complete(context.Background(), newReq())
	if err != nil {
		t.Fatalf("Complete 실패: %v", err)
	}
	if healthy.RequestCount() != 1 {
		t.Errorf("healthy = %d, want 1 (폴백)", healthy.RequestCount())
	}
	if !strings.Contains(resp.Provider, "failover") {
		t.Errorf("Provider = %q (failover 표시 누락)", resp.Provider)
	}
}

func TestFailoverChain_AllFail(t *testing.T) {
	chain := llm.NewFailoverChain(failingAdapter{}, failingAdapter{})

	_, err := chain.Complete(context.Background(), newReq())
	if err == nil {
		t.Error("모두 실패인데 통과")
	}
}

func TestFailoverChain_Empty(t *testing.T) {
	chain := llm.NewFailoverChain()

	_, err := chain.Complete(context.Background(), newReq())
	if err == nil {
		t.Error("빈 체인이 통과")
	}
	if err := chain.HealthCheck(context.Background()); err == nil {
		t.Error("빈 체인 HealthCheck 통과")
	}
}

func TestFailoverChain_HealthCheck_AnyHealthy(t *testing.T) {
	chain := llm.NewFailoverChain(failingAdapter{}, llm.NewNoopAdapter())
	if err := chain.HealthCheck(context.Background()); err != nil {
		t.Errorf("적어도 하나 healthy인데 실패: %v", err)
	}
}

func TestFailoverChain_Provider(t *testing.T) {
	chain := llm.NewFailoverChain(llm.NewNoopAdapter(), llm.NewOpenAIAdapter("k", ""))
	if !strings.Contains(chain.Provider(), "noop") {
		t.Errorf("Provider = %q", chain.Provider())
	}
}

func TestNoopAdapter_TokenEstimation(t *testing.T) {
	a := llm.NewNoopAdapter()
	resp, _ := a.Complete(context.Background(), newReq())
	if resp.PromptTokens <= 0 {
		t.Errorf("PromptTokens = %d, want > 0", resp.PromptTokens)
	}
	if resp.CompletionTokens < 0 {
		t.Errorf("CompletionTokens = %d", resp.CompletionTokens)
	}
}

func TestRoleConstants(t *testing.T) {
	if llm.RoleSystem != "system" {
		t.Error("RoleSystem")
	}
	if llm.RoleUser != "user" {
		t.Error("RoleUser")
	}
	if llm.RoleAssistant != "assistant" {
		t.Error("RoleAssistant")
	}
}
