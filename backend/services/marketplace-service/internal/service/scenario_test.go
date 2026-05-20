package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/manpasik/backend/services/marketplace-service/internal/service"
)

// TestScenario_PartnerOnboardingToSettlement은 파트너 가입 → 승인 → 거래 → 정산 전체 흐름입니다.
func TestScenario_PartnerOnboardingToSettlement(t *testing.T) {
	engine := service.NewApprovalWorkflowEngine()
	settlement := service.NewSettlementEngine()
	ctx := context.Background()

	// 1. 신청
	app := &service.PartnerApplication{
		ID:              "p-onboard-001",
		PartnerID:       "partner-A",
		Priority:        service.PriorityHigh,
		BusinessNumber:  "111-22-33333",
		BusinessLicense: "license.pdf",
	}
	if err := engine.SubmitApplication(ctx, app); err != nil {
		t.Fatalf("Submit 실패: %v", err)
	}

	// 2. 검토 큐 처리
	next := engine.GetNextForReview()
	if next == nil || next.ID != "p-onboard-001" {
		t.Fatalf("큐 dequeue 실패")
	}

	// 3. 단계별 승인
	stages := []string{
		service.WorkflowReviewing,
		service.WorkflowInspecting,
		service.WorkflowApproved,
	}
	for _, s := range stages {
		if err := engine.AdvanceStage(ctx, "p-onboard-001", s, "admin", "검토 통과"); err != nil {
			t.Fatalf("%s 전이 실패: %v", s, err)
		}
	}

	final, _ := engine.GetApplication("p-onboard-001")
	if final.Stage != service.WorkflowApproved {
		t.Errorf("Stage = %q", final.Stage)
	}
	if len(final.History) != 3 {
		t.Errorf("History = %d, want 3", len(final.History))
	}

	// 4. 정산 (월간)
	items := []*service.SettlementItem{
		{ProductID: "p1", OrderID: "o1", GrossPrice: 50000, Commission: 7500, OrderedAt: time.Now()},
		{ProductID: "p2", OrderID: "o2", GrossPrice: 100000, Commission: 15000, OrderedAt: time.Now()},
		{ProductID: "p3", OrderID: "o3", GrossPrice: 30000, Commission: 4500, OrderedAt: time.Now()},
	}
	record := settlement.GenerateMonthlySettlement("partner-A", 2026, 4, items)

	if record.GrossRevenue != 180000 {
		t.Errorf("GrossRevenue = %f", record.GrossRevenue)
	}
	if record.Commission != 27000 {
		t.Errorf("Commission = %f", record.Commission)
	}
	if record.NetPayout != 153000 {
		t.Errorf("NetPayout = %f", record.NetPayout)
	}

	// 5. 정산 완료
	if err := settlement.MarkSettled(record); err != nil {
		t.Fatalf("MarkSettled 실패: %v", err)
	}
	if record.Status != "paid" {
		t.Errorf("Status = %q", record.Status)
	}
}

// TestScenario_PartnerRejectionFlow는 거부 흐름을 검증합니다.
func TestScenario_PartnerRejectionFlow(t *testing.T) {
	engine := service.NewApprovalWorkflowEngine()
	ctx := context.Background()

	app := &service.PartnerApplication{
		ID: "p-reject", PartnerID: "p-r", BusinessNumber: "x", BusinessLicense: "y",
	}
	_ = engine.SubmitApplication(ctx, app)

	// 서류 심사 → 거부
	_ = engine.AdvanceStage(ctx, "p-reject", service.WorkflowReviewing, "admin", "")
	if err := engine.AdvanceStage(ctx, "p-reject", service.WorkflowRejected, "admin", "사업자번호 위조"); err != nil {
		t.Fatalf("거부 전이 실패: %v", err)
	}

	got, _ := engine.GetApplication("p-reject")
	if got.Stage != service.WorkflowRejected {
		t.Errorf("Stage = %q", got.Stage)
	}

	// 거부된 신청은 더 이상 전이 불가
	if err := engine.AdvanceStage(ctx, "p-reject", service.WorkflowReviewing, "admin", ""); err == nil {
		t.Error("거부 후 재심사 전이가 허용됨")
	}
}

// TestScenario_OnHoldRecoveryFlow는 보류 → 재심사 회복 흐름입니다.
func TestScenario_OnHoldRecoveryFlow(t *testing.T) {
	engine := service.NewApprovalWorkflowEngine()
	ctx := context.Background()

	app := &service.PartnerApplication{
		ID: "p-hold", PartnerID: "p-h", BusinessNumber: "x", BusinessLicense: "y",
	}
	_ = engine.SubmitApplication(ctx, app)
	_ = engine.AdvanceStage(ctx, "p-hold", service.WorkflowOnHold, "admin", "서류 보완 요청")

	// 보완 후 재검토
	if err := engine.AdvanceStage(ctx, "p-hold", service.WorkflowReviewing, "admin", "보완 완료"); err != nil {
		t.Fatalf("재검토 실패: %v", err)
	}
	_ = engine.AdvanceStage(ctx, "p-hold", service.WorkflowInspecting, "admin", "")
	if err := engine.AdvanceStage(ctx, "p-hold", service.WorkflowApproved, "admin", "통과"); err != nil {
		t.Fatalf("최종 승인 실패: %v", err)
	}
}

// TestScenario_PriorityQueue_UrgentBeforeNormal은 긴급 우선 처리를 검증합니다.
func TestScenario_PriorityQueue_UrgentBeforeNormal(t *testing.T) {
	engine := service.NewApprovalWorkflowEngine()
	ctx := context.Background()

	low := &service.PartnerApplication{ID: "low", PartnerID: "p-l", Priority: service.PriorityLow, BusinessNumber: "x", BusinessLicense: "y"}
	urgent := &service.PartnerApplication{ID: "urg", PartnerID: "p-u", Priority: service.PriorityUrgent, BusinessNumber: "x", BusinessLicense: "y"}

	_ = engine.SubmitApplication(ctx, low)
	_ = engine.SubmitApplication(ctx, urgent)

	first := engine.GetNextForReview()
	if first.ID != "urg" {
		t.Errorf("first = %q, want urg", first.ID)
	}
}

// TestScenario_DisputeFullLifecycle은 분쟁 신고 → 중재 → 해결 흐름입니다.
func TestScenario_DisputeFullLifecycle(t *testing.T) {
	resolver := service.NewDisputeResolver()

	d := &service.Dispute{
		ID:          "dp-life-001",
		PartnerID:   "p-001",
		OrderID:     "o-001",
		UserID:      "u-001",
		Reason:      "상품 불량",
		Description: "수령 시 외관 손상",
	}

	// 1. 신고
	if err := resolver.FileDispute(d); err != nil {
		t.Fatalf("FileDispute 실패: %v", err)
	}

	// 2. 중재자 배정
	if err := resolver.AssignMediator("dp-life-001", "mediator-001"); err != nil {
		t.Fatalf("AssignMediator 실패: %v", err)
	}

	got, _ := resolver.GetDispute("dp-life-001")
	if got.Status != service.DisputeMediating {
		t.Errorf("Status = %q", got.Status)
	}
	if got.Mediator != "mediator-001" {
		t.Errorf("Mediator = %q", got.Mediator)
	}

	// 3. 해결
	if err := resolver.ResolveDispute("dp-life-001", "전액 환불 후 상품 회수"); err != nil {
		t.Fatalf("ResolveDispute 실패: %v", err)
	}

	final, _ := resolver.GetDispute("dp-life-001")
	if final.Status != service.DisputeResolved {
		t.Errorf("Status = %q", final.Status)
	}
	if final.ResolvedAt == nil {
		t.Error("ResolvedAt 미설정")
	}
}

// TestScenario_TieredCommission는 파트너 등급별 차등 수수료를 검증합니다.
func TestScenario_TieredCommission(t *testing.T) {
	settlement := service.NewSettlementEngine()
	price := 10000.0

	// 카트리지 카테고리: 기본 15%
	plain := settlement.CommissionForCategory(price, "cartridge", "")
	silver := settlement.CommissionForCategory(price, "cartridge", "silver")  // 14%
	gold := settlement.CommissionForCategory(price, "cartridge", "gold")      // 12%
	platinum := settlement.CommissionForCategory(price, "cartridge", "platinum") // 10%

	if plain != 1500 {
		t.Errorf("plain = %f, want 1500", plain)
	}
	if silver >= plain {
		t.Errorf("silver(%f) >= plain(%f)", silver, plain)
	}
	if gold >= silver {
		t.Errorf("gold(%f) >= silver(%f)", gold, silver)
	}
	if platinum >= gold {
		t.Errorf("platinum(%f) >= gold(%f)", platinum, gold)
	}
}

// TestScenario_DisputeStatistics는 분쟁 상태별 통계를 검증합니다.
func TestScenario_DisputeStatistics(t *testing.T) {
	resolver := service.NewDisputeResolver()

	for i := 0; i < 5; i++ {
		_ = resolver.FileDispute(&service.Dispute{
			ID:        "stat-" + string(rune('a'+i)),
			PartnerID: "p", OrderID: "o", UserID: "u",
		})
	}

	// 1개 중재 중, 2개 해결, 1개 에스컬레이션
	_ = resolver.AssignMediator("stat-a", "med")
	_ = resolver.ResolveDispute("stat-b", "환불")
	_ = resolver.ResolveDispute("stat-c", "교환")
	_ = resolver.EscalateDispute("stat-d", "법무 검토")

	counts := resolver.CountByStatus()
	if counts[service.DisputeOpen] != 1 {
		t.Errorf("open = %d, want 1", counts[service.DisputeOpen])
	}
	if counts[service.DisputeMediating] != 1 {
		t.Errorf("mediating = %d, want 1", counts[service.DisputeMediating])
	}
	if counts[service.DisputeResolved] != 2 {
		t.Errorf("resolved = %d, want 2", counts[service.DisputeResolved])
	}
	if counts[service.DisputeEscalated] != 1 {
		t.Errorf("escalated = %d, want 1", counts[service.DisputeEscalated])
	}
}
