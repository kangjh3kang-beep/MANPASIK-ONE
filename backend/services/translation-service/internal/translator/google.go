package translator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// GoogleTranslator는 Google Cloud Translation API v2 클라이언트입니다.
type GoogleTranslator struct {
	apiKey     string
	httpClient *http.Client
	apiURL     string // 테스트에서 교체 가능
}

// NewGoogleTranslator는 GoogleTranslator를 생성합니다.
func NewGoogleTranslator(apiKey string) *GoogleTranslator {
	return &GoogleTranslator{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 15 * time.Second},
		apiURL:     "https://translation.googleapis.com/language/translate/v2",
	}
}

// SetAPIURL은 API URL을 설정합니다 (테스트용).
func (g *GoogleTranslator) SetAPIURL(url string) {
	g.apiURL = url
}

type googleTranslateResponse struct {
	Data struct {
		Translations []struct {
			TranslatedText         string `json:"translatedText"`
			DetectedSourceLanguage string `json:"detectedSourceLanguage"`
		} `json:"translations"`
	} `json:"data"`
}

// Translate는 Google Translate API로 텍스트를 번역합니다.
func (g *GoogleTranslator) Translate(ctx context.Context, text, sourceLang, targetLang string) (string, error) {
	if text == "" {
		return "", fmt.Errorf("번역할 텍스트가 비어있습니다")
	}
	if targetLang == "" {
		return "", fmt.Errorf("대상 언어가 비어있습니다")
	}

	params := url.Values{}
	params.Set("q", text)
	params.Set("target", targetLang)
	params.Set("format", "text")
	params.Set("key", g.apiKey)
	if sourceLang != "" {
		params.Set("source", sourceLang)
	}

	reqURL := fmt.Sprintf("%s?%s", g.apiURL, params.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("Google Translate 요청 생성 실패: %w", err)
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("Google Translate API 호출 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Google Translate API 오류 (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var result googleTranslateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("Google Translate 응답 파싱 실패: %w", err)
	}

	if len(result.Data.Translations) == 0 {
		return "", fmt.Errorf("Google Translate 응답에 번역 결과가 없습니다")
	}

	return result.Data.Translations[0].TranslatedText, nil
}
