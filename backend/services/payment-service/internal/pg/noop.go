// Package pg는 PG사(Toss 등) 결제 연동 구현을 제공합니다.
package pg

import "context"

// NoopGateway는 PG 연동이 없을 때 사용하는 no-op 구현입니다.
// Confirm은 호출만 하고 pgTransactionID는 빈 문자열, err nil 반환.
// Cancel은 호출만 하고 nil 반환. B-3 Toss 구현 전까지 기본값으로 사용합니다.
type NoopGateway struct{}

// NewNoopGateway는 NoopGateway를 생성합니다.
func NewNoopGateway() *NoopGateway {
	return &NoopGateway{}
}

// Confirm은 no-op: 실제 PG 호출 없이 성공으로 처리합니다.
func (n *NoopGateway) Confirm(ctx context.Context, paymentKey, orderId string, amountKRW int32) (string, error) {
	return "", nil
}

// Cancel은 no-op: 실제 PG 호출 없이 성공으로 처리합니다.
func (n *NoopGateway) Cancel(ctx context.Context, paymentKey, reason string) error {
	return nil
}
