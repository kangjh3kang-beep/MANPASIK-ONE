package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// OpenAIResponder는 OpenAI Chat Completions API를 사용하는 응답 생성기입니다.
type OpenAIResponder struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

// OpenAIConfig는 OpenAI 연결 설정입니다.
type OpenAIConfig struct {
	APIKey  string
	Model   string // 기본값: "gpt-4o"
	BaseURL string // 기본값: "https://api.openai.com/v1"
}

// NewOpenAIResponder는 OpenAI 기반 응답 생성기를 생성합니다.
func NewOpenAIResponder(cfg OpenAIConfig) *OpenAIResponder {
	if cfg.Model == "" {
		cfg.Model = "gpt-4o"
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.openai.com/v1"
	}
	return &OpenAIResponder{
		apiKey:  cfg.APIKey,
		model:   cfg.Model,
		baseURL: cfg.BaseURL,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

const systemPrompt = `당신은 만파식(ManPaSik) 건강 관리 앱의 AI 비서입니다.
사용자의 건강 관련 질문에 친절하고 전문적으로 답변합니다.

역할:
- 건강 측정 결과 해석 (혈당, 혈압, 콜레스테롤 등)
- 건강 코칭 및 생활습관 조언
- 의료 예약, 처방전, 알림 안내
- 식사/운동/수면 관리 도움

규칙:
- 한국어로 응답
- 의학적 진단이나 처방은 하지 않음
- 긴급 상황(흉통, 호흡곤란 등) 감지 시 즉시 119 안내
- 간결하고 실용적인 답변 (200자 이내 권장)`

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// GenerateResponse는 OpenAI API로 대화 응답을 생성합니다.
func (r *OpenAIResponder) GenerateResponse(ctx context.Context, history []ChatMessage, userMessage string) (*AIResponse, error) {
	if r.apiKey == "" {
		return nil, fmt.Errorf("OpenAI API 키가 설정되지 않았습니다")
	}

	// 메시지 조합: 시스템 + 이력 + 사용자
	messages := []chatMessage{
		{Role: "system", Content: systemPrompt},
	}
	for _, h := range history {
		messages = append(messages, chatMessage{Role: h.Role, Content: h.Content})
	}
	messages = append(messages, chatMessage{Role: "user", Content: userMessage})

	reqBody := chatRequest{
		Model:    r.model,
		Messages: messages,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("요청 직렬화 실패: %w", err)
	}

	url := r.baseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("HTTP 요청 생성 실패: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.apiKey)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("OpenAI API 호출 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API 응답 오류: status=%d", resp.StatusCode)
	}

	var chatResp chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, fmt.Errorf("OpenAI 응답 파싱 실패: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("OpenAI 응답에 선택지가 없습니다")
	}

	content := chatResp.Choices[0].Message.Content
	if content == "" {
		return nil, fmt.Errorf("OpenAI 응답이 비어 있습니다")
	}

	return &AIResponse{
		Content: content,
	}, nil
}
