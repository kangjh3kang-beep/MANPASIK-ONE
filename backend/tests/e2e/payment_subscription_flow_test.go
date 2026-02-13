package e2e

import (
	"context"
	"testing"
	"time"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TestPaymentFlow E2E: CreatePayment → GetPayment
func TestPaymentFlow(t *testing.T) {
	dialCtx, dialCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer dialCancel()

	payConn, err := grpc.DialContext(dialCtx, PaymentAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Skipf("payment 서비스 연결 불가 (기동 후 재실행): %v", err)
	}
	defer payConn.Close()

	client := v1.NewPaymentServiceClient(payConn)
	rpcCtx, rpcCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer rpcCancel()

	// 1) CreatePayment — Proto: user_id, amount_krw, payment_type, payment_method
	createResp, err := client.CreatePayment(rpcCtx, &v1.CreatePaymentRequest{
		UserId:        "e2e-user-pay-001",
		AmountKrw:     9900,
		PaymentType:   v1.PaymentType_PAYMENT_TYPE_ONE_TIME,
		PaymentMethod: "card",
	})
	if err != nil {
		t.Fatalf("CreatePayment 실패: %v", err)
	}
	if createResp.PaymentId == "" {
		t.Fatal("payment_id가 비어 있습니다")
	}
	t.Logf("CreatePayment 성공: id=%s, status=%v", createResp.PaymentId, createResp.Status)

	// 2) GetPayment
	getResp, err := client.GetPayment(rpcCtx, &v1.GetPaymentRequest{
		PaymentId: createResp.PaymentId,
	})
	if err != nil {
		t.Fatalf("GetPayment 실패: %v", err)
	}
	if getResp.PaymentId != createResp.PaymentId {
		t.Errorf("ID 불일치: %s != %s", getResp.PaymentId, createResp.PaymentId)
	}
	t.Logf("GetPayment 성공: %s, amount=%d, status=%v", getResp.PaymentId, getResp.AmountKrw, getResp.Status)

	// 3) ConfirmPayment — Proto: payment_id, pg_transaction_id, pg_provider, payment_key
	confirmResp, err := client.ConfirmPayment(rpcCtx, &v1.ConfirmPaymentRequest{
		PaymentId:  createResp.PaymentId,
		PaymentKey: "test_pay_key_001",
		PgProvider: "toss",
	})
	if err != nil {
		t.Logf("ConfirmPayment 실패 (PG 미연동 시 예상): %v", err)
	} else {
		t.Logf("ConfirmPayment 성공: status=%v", confirmResp.Status)
	}
}

// TestSubscriptionFlow E2E: CreateSubscription → GetSubscription → UpgradeSubscription
func TestSubscriptionFlow(t *testing.T) {
	dialCtx, dialCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer dialCancel()

	subConn, err := grpc.DialContext(dialCtx, SubscriptionAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Skipf("subscription 서비스 연결 불가 (기동 후 재실행): %v", err)
	}
	defer subConn.Close()

	client := v1.NewSubscriptionServiceClient(subConn)
	rpcCtx, rpcCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer rpcCancel()

	// 1) CreateSubscription
	createResp, err := client.CreateSubscription(rpcCtx, &v1.CreateSubscriptionRequest{
		UserId: "e2e-user-sub-001",
		Tier:   v1.SubscriptionTier_SUBSCRIPTION_TIER_BASIC,
	})
	if err != nil {
		t.Fatalf("CreateSubscription 실패: %v", err)
	}
	if createResp.SubscriptionId == "" {
		t.Fatal("subscription_id가 비어 있습니다")
	}
	t.Logf("CreateSubscription 성공: id=%s, tier=%v", createResp.SubscriptionId, createResp.Tier)

	// 2) GetSubscription — Proto: GetSubscriptionDetailRequest
	getResp, err := client.GetSubscription(rpcCtx, &v1.GetSubscriptionDetailRequest{
		UserId: "e2e-user-sub-001",
	})
	if err != nil {
		t.Fatalf("GetSubscription 실패: %v", err)
	}
	t.Logf("GetSubscription 성공: tier=%v, status=%v", getResp.Tier, getResp.Status)

	// 3) UpdateSubscription (업그레이드) — Proto: UpdateSubscriptionRequest
	upgradeResp, err := client.UpdateSubscription(rpcCtx, &v1.UpdateSubscriptionRequest{
		UserId:  "e2e-user-sub-001",
		NewTier: v1.SubscriptionTier_SUBSCRIPTION_TIER_PRO,
	})
	if err != nil {
		t.Logf("UpdateSubscription 실패 (예상 가능): %v", err)
	} else {
		t.Logf("UpdateSubscription 성공: new_tier=%v", upgradeResp.Tier)
	}
}
