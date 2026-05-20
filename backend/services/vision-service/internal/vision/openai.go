// Package vision은 AI 비전 분석 외부 클라이언트를 구현합니다.
package vision

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// FoodItemResult는 OpenAI Vision API 응답의 음식 항목입니다.
type FoodItemResult struct {
	Name        string  `json:"name"`
	Confidence  float64 `json:"confidence"`
	CalorieKcal float64 `json:"calorie_kcal"`
	PortionG    float64 `json:"portion_g"`
	Nutrients   []struct {
		Name   string  `json:"name"`
		Amount float64 `json:"amount"`
		Unit   string  `json:"unit"`
		DV     float64 `json:"dv"`
	} `json:"nutrients"`
}

// OpenAIVisionAnalyzer는 OpenAI Vision API를 사용하는 음식 분석기입니다.
type OpenAIVisionAnalyzer struct {
	apiKey     string
	model      string
	baseURL    string
	httpClient *http.Client
}

// NewOpenAIVisionAnalyzer는 OpenAIVisionAnalyzer를 생성합니다.
func NewOpenAIVisionAnalyzer(apiKey, model, baseURL string) *OpenAIVisionAnalyzer {
	if model == "" {
		model = "gpt-4o"
	}
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	return &OpenAIVisionAnalyzer{
		apiKey:     apiKey,
		model:      model,
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

const visionSystemPrompt = `당신은 한국 음식 영양 분석 전문가입니다.
이미지에 있는 음식을 분석하고 다음 JSON 형식으로만 반환하세요. 설명 없이 JSON 배열만:
[{"name":"음식이름","confidence":0.0-1.0,"calorie_kcal":숫자,"portion_g":숫자,"nutrients":[{"name":"영양소","amount":숫자,"unit":"g","dv":0.0-1.0}]}]`

type visionRequest struct {
	Model    string           `json:"model"`
	Messages []visionMessage  `json:"messages"`
}

type visionMessage struct {
	Role    string        `json:"role"`
	Content interface{}   `json:"content"`
}

type contentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *imageURL `json:"image_url,omitempty"`
}

type imageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail"`
}

type visionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// AnalyzeFood는 이미지 URL을 OpenAI Vision API로 분석합니다.
func (a *OpenAIVisionAnalyzer) AnalyzeFood(ctx context.Context, imgURL string) ([]FoodItemResult, error) {
	if a.apiKey == "" {
		return nil, fmt.Errorf("Vision API 키가 설정되지 않았습니다")
	}

	reqBody := visionRequest{
		Model: a.model,
		Messages: []visionMessage{
			{Role: "system", Content: visionSystemPrompt},
			{Role: "user", Content: []contentPart{
				{Type: "text", Text: "이 이미지의 음식을 분석해 주세요."},
				{Type: "image_url", ImageURL: &imageURL{URL: imgURL, Detail: "auto"}},
			}},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("요청 직렬화 실패: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/chat/completions", a.baseURL), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("요청 생성 실패: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Vision API 호출 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Vision API 오류 (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var vResp visionResponse
	if err := json.NewDecoder(resp.Body).Decode(&vResp); err != nil {
		return nil, fmt.Errorf("Vision API 응답 파싱 실패: %w", err)
	}

	if len(vResp.Choices) == 0 {
		return nil, fmt.Errorf("Vision API 응답에 결과가 없습니다")
	}

	content := vResp.Choices[0].Message.Content

	// JSON 배열 파싱
	var items []FoodItemResult
	if err := json.Unmarshal([]byte(content), &items); err != nil {
		return nil, fmt.Errorf("음식 목록 파싱 실패: %w", err)
	}

	return items, nil
}
