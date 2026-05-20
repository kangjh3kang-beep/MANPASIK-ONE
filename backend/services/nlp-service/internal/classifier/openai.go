package classifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// OpenAIClassifier는 OpenAI Chat Completions API를 사용하는 인텐트 분류기입니다.
type OpenAIClassifier struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

// OpenAIConfig는 OpenAI 연결 설정입니다.
type OpenAIConfig struct {
	APIKey  string
	Model   string // 기본값: "gpt-4o-mini"
	BaseURL string // 기본값: "https://api.openai.com/v1"
}

// NewOpenAIClassifier는 OpenAI 기반 인텐트 분류기를 생성합니다.
func NewOpenAIClassifier(cfg OpenAIConfig) *OpenAIClassifier {
	if cfg.Model == "" {
		cfg.Model = "gpt-4o-mini"
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.openai.com/v1"
	}
	return &OpenAIClassifier{
		apiKey:  cfg.APIKey,
		model:   cfg.Model,
		baseURL: cfg.BaseURL,
		client:  &http.Client{Timeout: 15 * time.Second},
	}
}

const classifySystemPrompt = `당신은 건강 관리 앱의 NLP 인텐트 분류기입니다.
사용자 텍스트에서 인텐트, 엔티티, 긴급도를 추출하여 JSON으로 반환합니다.

가능한 인텐트:
- health_inquiry.blood_sugar, health_inquiry.blood_pressure, health_inquiry.cholesterol
- health_inquiry.heart_rate, health_inquiry.body_temperature, health_inquiry.oxygen_saturation
- health_inquiry.weight, health_inquiry.sleep, health_inquiry.exercise, health_inquiry.nutrition
- health_inquiry.measurement, health_inquiry (일반)

가능한 긴급도: routine, urgent, emergency
- emergency: 흉통, 호흡곤란, 의식불명 등 생명 위협 증상
- urgent: 2개 이상 중등도 증상 동시 발현
- routine: 기타

반드시 다음 JSON 형식으로만 응답:
{"intent":"...","entities":["..."],"confidence":0.95,"urgency":"routine"}`

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

// Classify는 OpenAI API로 텍스트의 인텐트를 분류합니다.
func (c *OpenAIClassifier) Classify(ctx context.Context, text string) (*ClassificationResult, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("OpenAI API 키가 설정되지 않았습니다")
	}

	reqBody := chatRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "system", Content: classifySystemPrompt},
			{Role: "user", Content: text},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("요청 직렬화 실패: %w", err)
	}

	url := c.baseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("HTTP 요청 생성 실패: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
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

	var result ClassificationResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("분류 결과 JSON 파싱 실패: %w", err)
	}

	return &result, nil
}
