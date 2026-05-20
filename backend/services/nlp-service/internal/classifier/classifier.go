// Package classifier는 NLP 인텐트 분류를 위한 외부 API 연동 인터페이스를 제공합니다.
package classifier

import "context"

// ClassificationResult는 인텐트 분류 결과입니다.
type ClassificationResult struct {
	Intent     string   // 분류된 인텐트 (예: "health_inquiry.blood_sugar")
	Entities   []string // 추출된 엔티티
	Confidence float64  // 분류 신뢰도 (0.0~1.0)
	Urgency    string   // "routine", "urgent", "emergency"
}

// IntentClassifier는 텍스트에서 인텐트를 분류하는 인터페이스입니다.
type IntentClassifier interface {
	// Classify는 텍스트에서 인텐트, 엔티티, 긴급도를 추출합니다.
	Classify(ctx context.Context, text string) (*ClassificationResult, error)
}

// NoopClassifier는 nil을 반환하여 기존 키워드 기반 분류로 폴백하도록 유도합니다.
type NoopClassifier struct{}

// NewNoopClassifier는 NoopClassifier를 생성합니다.
func NewNoopClassifier() *NoopClassifier {
	return &NoopClassifier{}
}

// Classify는 항상 nil을 반환하여 폴백을 유도합니다.
func (n *NoopClassifier) Classify(_ context.Context, _ string) (*ClassificationResult, error) {
	return nil, nil
}
