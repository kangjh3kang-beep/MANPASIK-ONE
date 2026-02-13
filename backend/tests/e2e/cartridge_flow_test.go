package e2e

import (
	"context"
	"testing"
	"time"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TestCartridgeFlow E2E: ValidateCartridge → RecordUsage → GetUsageHistory → GetRemainingUses
func TestCartridgeFlow(t *testing.T) {
	dialCtx, dialCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer dialCancel()

	conn, err := grpc.DialContext(dialCtx, CartridgeAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Skipf("cartridge 서비스 연결 불가: %v", err)
	}
	defer conn.Close()

	client := v1.NewCartridgeServiceClient(conn)
	rpcCtx, rpcCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer rpcCancel()

	cartridgeUID := "e2e-cart-uid-001"

	// 1) ValidateCartridge — 카트리지 유효성 검증
	validateResp, err := client.ValidateCartridge(rpcCtx, &v1.ValidateCartridgeRequest{
		CartridgeUid: cartridgeUID,
		CategoryCode: 0x01,
		TypeIndex:    1,
		UserId:       "e2e-user-cart-001",
	})
	if err != nil {
		t.Fatalf("ValidateCartridge 실패: %v", err)
	}
	t.Logf("ValidateCartridge 성공: valid=%v, reason=%s, remaining=%d, access=%v",
		validateResp.IsValid, validateResp.Reason, validateResp.RemainingUses, validateResp.AccessLevel)

	// 2) RecordUsage — 카트리지 사용 기록
	usageResp, err := client.RecordUsage(rpcCtx, &v1.RecordUsageRequest{
		UserId:       "e2e-user-cart-001",
		SessionId:    "e2e-session-cart-001",
		CartridgeUid: cartridgeUID,
		CategoryCode: 0x01,
		TypeIndex:    1,
	})
	if err != nil {
		t.Logf("RecordUsage 실패 (카트리지 미등록 시 예상): %v", err)
	} else {
		t.Logf("RecordUsage 성공: success=%v, remaining_uses=%d, remaining_daily=%d",
			usageResp.Success, usageResp.RemainingUses, usageResp.RemainingDaily)
	}

	// 3) GetUsageHistory — 사용 이력 조회
	historyResp, err := client.GetUsageHistory(rpcCtx, &v1.GetUsageHistoryRequest{
		UserId: "e2e-user-cart-001",
		Limit:  10,
	})
	if err != nil {
		t.Fatalf("GetUsageHistory 실패: %v", err)
	}
	t.Logf("GetUsageHistory 성공: total=%d", historyResp.TotalCount)

	// 4) GetRemainingUses — 잔여 사용 횟수 조회
	remainResp, err := client.GetRemainingUses(rpcCtx, &v1.GetRemainingUsesRequest{
		CartridgeUid: cartridgeUID,
	})
	if err != nil {
		t.Logf("GetRemainingUses 실패 (카트리지 미등록 시 예상): %v", err)
	} else {
		t.Logf("GetRemainingUses 성공: uid=%s, remaining=%d/%d, expired=%v",
			remainResp.CartridgeUid, remainResp.RemainingUses, remainResp.MaxUses, remainResp.IsExpired)
	}

	// 5) ListCategories — 카테고리 목록 조회
	catResp, err := client.ListCategories(rpcCtx, &v1.ListCategoriesRequest{})
	if err != nil {
		t.Logf("ListCategories 실패 (레지스트리 미초기화 시 예상): %v", err)
	} else {
		t.Logf("ListCategories 성공: %d개 카테고리", len(catResp.Categories))
	}

	t.Logf("✅ 카트리지 플로우 완료: uid=%s", cartridgeUID)
}
