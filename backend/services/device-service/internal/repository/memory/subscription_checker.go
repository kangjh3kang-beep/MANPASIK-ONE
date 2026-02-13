package memory

import (
	"context"
)

// SubscriptionChecker는 인메모리 구독 확인기입니다 (개발용 - 항상 무제한 허용).
type SubscriptionChecker struct{}

// NewSubscriptionChecker는 인메모리 SubscriptionChecker를 생성합니다.
func NewSubscriptionChecker() *SubscriptionChecker {
	return &SubscriptionChecker{}
}

// GetMaxDevices는 최대 디바이스 수를 반환합니다 (개발용: 무제한).
func (c *SubscriptionChecker) GetMaxDevices(_ context.Context, _ string) (int, error) {
	return 999, nil
}
