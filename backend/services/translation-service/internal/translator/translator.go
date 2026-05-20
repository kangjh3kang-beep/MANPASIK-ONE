// Package translator는 외부 번역 API 클라이언트를 구현합니다.
package translator

import (
	"context"
	"fmt"
)

// Translator는 외부 번역 API 인터페이스입니다.
type Translator interface {
	Translate(ctx context.Context, text, sourceLang, targetLang string) (string, error)
}

// NoopTranslator는 외부 번역 API 없이 에러를 반환합니다.
// 에러 반환 시 TranslationService가 기존 의료사전+시뮬레이션 폴백을 사용합니다.
type NoopTranslator struct{}

// NewNoopTranslator는 NoopTranslator를 생성합니다.
func NewNoopTranslator() *NoopTranslator {
	return &NoopTranslator{}
}

// Translate는 에러를 반환합니다.
func (n *NoopTranslator) Translate(_ context.Context, _, _, _ string) (string, error) {
	return "", fmt.Errorf("외부 번역 API 미설정")
}
