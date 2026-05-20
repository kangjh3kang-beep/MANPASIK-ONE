package vision

import (
	"context"
	"fmt"
	"log"
)

// LogVisionAnalyzer는 Vision API 미설정 시 로그만 기록하는 분석기입니다.
// 에러를 반환하여 VisionService가 fallbackAnalysis를 사용하도록 합니다.
type LogVisionAnalyzer struct{}

// NewLogVisionAnalyzer는 LogVisionAnalyzer를 생성합니다.
func NewLogVisionAnalyzer() *LogVisionAnalyzer {
	return &LogVisionAnalyzer{}
}

// AnalyzeFood는 로그만 기록하고 에러를 반환합니다.
func (l *LogVisionAnalyzer) AnalyzeFood(_ context.Context, imageURL string) ([]FoodItemResult, error) {
	log.Printf("[VISION-LOG] imageURL=%s (Vision API 미설정, fallback 사용)", imageURL)
	return nil, fmt.Errorf("Vision API 미설정")
}
