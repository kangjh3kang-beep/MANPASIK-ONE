package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/manpasik/backend/services/marketplace-service/internal/service"
)

func newTestApp(id string, priority int) *service.PartnerApplication {
	return &service.PartnerApplication{
		ID:              id,
		PartnerID:       "p-" + id,
		Priority:        priority,
		BusinessNumber:  "123-45-67890",
		BusinessLicense: "license.pdf",
		BankAccount:     "신한 110-123-456",
	}
}

// TestSubmitApplication는 신청서 제출을 검증합니다.
func TestSubmitApplication(t *testing.T) {
	engine := service.NewApprovalWorkflowEngine()
	app := newTestApp("app-001", service.PriorityNormal)

	if err := engine.SubmitApplication(context.Background(), app); err != nil {
		t.Fatalf("SubmitApplication 실패: %v", err)
	}

	got, err := engine.GetApplication("app-001")
	if err != nil {
		t.Fatalf("GetApplication 실패: %v", err)
	}
	if got.Stage != service.WorkflowSubmitted {
		t.Errorf("Stage = %q, want %q", got.Stage, service.WorkflowSubmitted)
	}
	if got.EstimatedDecision.IsZero() {
		t.Error("EstimatedDecision이 설정되지 않음")
	}
}

// TestSubmitApplication_MissingBusinessNumber는 사업자번호 누락 검증.
func TestSubmitApplication_MissingBusinessNumber(t *testing.T) {
	engine := service.NewApprovalWorkflowEngine()
	app := &service.PartnerApplication{ID: "x", PartnerID: "p-x", BusinessLicense: "l.pdf"}

	err := engine.SubmitApplication(context.Background(), app)
	if err == nil {
		t.Error("사업자번호 누락 시 에러가 발생해야 함")
	}
}

// TestAdvanceStage_HappyPath는 정상 워크플로우 진행을 검증합니다.
func TestAdvanceStage_HappyPath(t *testing.T) {
	engine := service.NewApprovalWorkflowEngine()
	app := newTestApp("app-002", service.PriorityNormal)
	_ = engine.SubmitApplication(context.Background(), app)

	stages := []string{
		service.WorkflowReviewing,
		service.WorkflowInspecting,
		service.WorkflowApproved,
	}

	for _, s := range stages {
		if err := engine.AdvanceStage(context.Background(), "app-002", s, "admin", "ok"); err != nil {
			t.Fatalf("%s 전이 실패: %v", s, err)
		}
	}

	got, _ := engine.GetApplication("app-002")
	if got.Stage != service.WorkflowApproved {
		t.Errorf("최종 Stage = %q, want %q", got.Stage, service.WorkflowApproved)
	}
	if len(got.History) != 3 {
		t.Errorf("History 길이 = %d, want 3", len(got.History))
	}
}

// TestAdvanceStage_InvalidTransition는 잘못된 전이를 거부합니다.
func TestAdvanceStage_InvalidTransition(t *testing.T) {
	engine := service.NewApprovalWorkflowEngine()
	app := newTestApp("app-003", service.PriorityNormal)
	_ = engine.SubmitApplication(context.Background(), app)

	// submitted → approved 직접 전이는 불가
	err := engine.AdvanceStage(context.Background(), "app-003", service.WorkflowApproved, "admin", "")
	if err == nil {
		t.Error("불가능한 전이가 허용됨")
	}
}

// TestAdvanceStage_OnHoldRecovery는 보류→재심사 회복을 검증합니다.
func TestAdvanceStage_OnHoldRecovery(t *testing.T) {
	engine := service.NewApprovalWorkflowEngine()
	app := newTestApp("app-004", service.PriorityNormal)
	_ = engine.SubmitApplication(context.Background(), app)
	_ = engine.AdvanceStage(context.Background(), "app-004", service.WorkflowOnHold, "admin", "서류 보완 필요")
	if err := engine.AdvanceStage(context.Background(), "app-004", service.WorkflowReviewing, "admin", "재심사"); err != nil {
		t.Fatalf("on_hold → reviewing 실패: %v", err)
	}
}

// TestReviewQueue_Priority는 우선순위 큐 정렬을 검증합니다.
func TestReviewQueue_Priority(t *testing.T) {
	engine := service.NewApprovalWorkflowEngine()
	low := newTestApp("low", service.PriorityLow)
	urgent := newTestApp("urgent", service.PriorityUrgent)
	normal := newTestApp("normal", service.PriorityNormal)

	_ = engine.SubmitApplication(context.Background(), low)
	_ = engine.SubmitApplication(context.Background(), urgent)
	_ = engine.SubmitApplication(context.Background(), normal)

	first := engine.GetNextForReview()
	if first == nil || first.ID != "urgent" {
		t.Errorf("첫 번째 검토 = %v, want urgent", first)
	}

	second := engine.GetNextForReview()
	if second == nil || second.ID != "normal" {
		t.Errorf("두 번째 검토 = %v, want normal", second)
	}
}

// TestSettlementEngine_GenerateMonthly는 월간 정산 생성을 검증합니다.
func TestSettlementEngine_GenerateMonthly(t *testing.T) {
	engine := service.NewSettlementEngine()

	items := []*service.SettlementItem{
		{ProductID: "p1", OrderID: "o1", GrossPrice: 10000, Commission: 1500, Net: 8500, OrderedAt: time.Now()},
		{ProductID: "p2", OrderID: "o2", GrossPrice: 20000, Commission: 3000, Net: 17000, OrderedAt: time.Now()},
		{ProductID: "p3", OrderID: "o3", GrossPrice: 5000, Commission: 750, Net: 4250, OrderedAt: time.Now()},
	}

	settlement := engine.GenerateMonthlySettlement("partner-001", 2026, 4, items)

	if settlement.GrossRevenue != 35000 {
		t.Errorf("GrossRevenue = %f, want 35000", settlement.GrossRevenue)
	}
	if settlement.Commission != 5250 {
		t.Errorf("Commission = %f, want 5250", settlement.Commission)
	}
	if settlement.NetPayout != 29750 {
		t.Errorf("NetPayout = %f, want 29750", settlement.NetPayout)
	}
	if settlement.OrderCount != 3 {
		t.Errorf("OrderCount = %d, want 3", settlement.OrderCount)
	}
	if settlement.Status != "pending" {
		t.Errorf("Status = %q, want pending", settlement.Status)
	}
}

// TestSettlementEngine_MarkSettled는 정산 완료 처리를 검증합니다.
func TestSettlementEngine_MarkSettled(t *testing.T) {
	engine := service.NewSettlementEngine()
	settlement := engine.GenerateMonthlySettlement("partner-001", 2026, 4, nil)

	if err := engine.MarkSettled(settlement); err != nil {
		t.Fatalf("MarkSettled 실패: %v", err)
	}
	if settlement.Status != "paid" {
		t.Errorf("Status = %q, want paid", settlement.Status)
	}
	if settlement.SettledAt == nil {
		t.Error("SettledAt이 설정되지 않음")
	}

	// 두 번째 호출은 실패해야 함
	if err := engine.MarkSettled(settlement); err == nil {
		t.Error("이미 정산 완료된 항목에 다시 정산 처리됨")
	}
}

// TestCommissionForCategory_PartnerTier는 파트너 등급별 수수료 할인을 검증합니다.
func TestCommissionForCategory_PartnerTier(t *testing.T) {
	engine := service.NewSettlementEngine()

	// 카트리지 기본 15%
	plain := engine.CommissionForCategory(10000, "cartridge", "")
	platinum := engine.CommissionForCategory(10000, "cartridge", "platinum")

	if platinum >= plain {
		t.Errorf("플래티넘 수수료(%f) >= 기본 수수료(%f), 할인이 적용되지 않음", platinum, plain)
	}

	// 최소 5% 보장
	tiny := engine.CommissionForCategory(100, "cartridge", "platinum")
	if tiny < 100*0.05 {
		t.Errorf("최소 수수료 5%% 미만: %f", tiny)
	}
}

// TestDisputeResolver_FileAndResolve는 분쟁 등록/해결 흐름을 검증합니다.
func TestDisputeResolver_FileAndResolve(t *testing.T) {
	resolver := service.NewDisputeResolver()

	d := &service.Dispute{
		ID:        "dp-001",
		PartnerID: "p-001",
		OrderID:   "o-001",
		UserID:    "u-001",
		Reason:    "상품 불량",
	}

	if err := resolver.FileDispute(d); err != nil {
		t.Fatalf("FileDispute 실패: %v", err)
	}

	got, _ := resolver.GetDispute("dp-001")
	if got.Status != service.DisputeOpen {
		t.Errorf("Status = %q, want %q", got.Status, service.DisputeOpen)
	}

	if err := resolver.AssignMediator("dp-001", "med-jin"); err != nil {
		t.Fatalf("AssignMediator 실패: %v", err)
	}

	got, _ = resolver.GetDispute("dp-001")
	if got.Status != service.DisputeMediating {
		t.Errorf("Status = %q, want %q", got.Status, service.DisputeMediating)
	}

	if err := resolver.ResolveDispute("dp-001", "환불 처리"); err != nil {
		t.Fatalf("ResolveDispute 실패: %v", err)
	}

	got, _ = resolver.GetDispute("dp-001")
	if got.Status != service.DisputeResolved {
		t.Errorf("Status = %q, want %q", got.Status, service.DisputeResolved)
	}
	if got.ResolvedAt == nil {
		t.Error("ResolvedAt 미설정")
	}
}

// TestDisputeResolver_Escalate는 에스컬레이션을 검증합니다.
func TestDisputeResolver_Escalate(t *testing.T) {
	resolver := service.NewDisputeResolver()
	d := &service.Dispute{ID: "dp-002", PartnerID: "p-2", OrderID: "o-2", UserID: "u-2"}
	_ = resolver.FileDispute(d)

	if err := resolver.EscalateDispute("dp-002", "법무팀 검토 필요"); err != nil {
		t.Fatalf("EscalateDispute 실패: %v", err)
	}
	got, _ := resolver.GetDispute("dp-002")
	if got.Status != service.DisputeEscalated {
		t.Errorf("Status = %q, want %q", got.Status, service.DisputeEscalated)
	}
}

// TestCountByStage는 단계별 통계를 검증합니다.
func TestCountByStage(t *testing.T) {
	engine := service.NewApprovalWorkflowEngine()

	for i := 0; i < 3; i++ {
		_ = engine.SubmitApplication(context.Background(), newTestApp("a"+string(rune('0'+i)), service.PriorityNormal))
	}

	counts := engine.CountByStage()
	if counts[service.WorkflowSubmitted] != 3 {
		t.Errorf("submitted 건수 = %d, want 3", counts[service.WorkflowSubmitted])
	}
}

// TestPendingByPriority는 우선순위별 대기 건수를 검증합니다.
func TestPendingByPriority(t *testing.T) {
	q := service.NewReviewQueue()
	q.Enqueue(newTestApp("hi", service.PriorityHigh))
	q.Enqueue(newTestApp("u", service.PriorityUrgent))
	q.Enqueue(newTestApp("u2", service.PriorityUrgent))

	counts := q.PendingByPriority()
	if counts[service.PriorityUrgent] != 2 {
		t.Errorf("urgent 건수 = %d, want 2", counts[service.PriorityUrgent])
	}
}

// TestEnqueueSize는 큐 크기를 검증합니다.
func TestEnqueueSize(t *testing.T) {
	q := service.NewReviewQueue()
	q.Enqueue(newTestApp("a", 1))
	q.Enqueue(newTestApp("b", 2))
	if q.Size() != 2 {
		t.Errorf("Size = %d, want 2", q.Size())
	}
}

// TestSubmittedAtPriority는 같은 우선순위에서 더 일찍 제출된 것이 우선인지 검증합니다.
func TestSubmittedAtPriority(t *testing.T) {
	q := service.NewReviewQueue()
	old := newTestApp("old", service.PriorityNormal)
	old.SubmittedAt = time.Now().Add(-1 * time.Hour)
	old.Stage = service.WorkflowSubmitted

	new := newTestApp("new", service.PriorityNormal)
	new.SubmittedAt = time.Now()
	new.Stage = service.WorkflowSubmitted

	q.Enqueue(new)
	q.Enqueue(old)

	first := q.Dequeue()
	if first == nil || first.ID != "old" {
		t.Errorf("first = %v, want old", first)
	}
}

// TestApprovalRecordHistory는 승인 이력 기록을 검증합니다.
func TestApprovalRecordHistory(t *testing.T) {
	engine := service.NewApprovalWorkflowEngine()
	app := newTestApp("hist", service.PriorityNormal)
	_ = engine.SubmitApplication(context.Background(), app)
	_ = engine.AdvanceStage(context.Background(), "hist", service.WorkflowReviewing, "admin1", "서류 검토 완료")
	_ = engine.AdvanceStage(context.Background(), "hist", service.WorkflowInspecting, "admin2", "실사 중")
	_ = engine.AdvanceStage(context.Background(), "hist", service.WorkflowRejected, "admin1", "사업자등록증 위조 의심")

	got, _ := engine.GetApplication("hist")
	if len(got.History) != 3 {
		t.Errorf("History 길이 = %d, want 3", len(got.History))
	}
	if got.Stage != service.WorkflowRejected {
		t.Errorf("최종 Stage = %q, want rejected", got.Stage)
	}

	// 거부 단계의 결정값
	last := got.History[len(got.History)-1]
	if last.Decision != "fail" {
		t.Errorf("Decision = %q, want fail", last.Decision)
	}
}

// TestDisputeResolver_CountByStatus는 분쟁 상태별 통계를 검증합니다.
func TestDisputeResolver_CountByStatus(t *testing.T) {
	resolver := service.NewDisputeResolver()
	_ = resolver.FileDispute(&service.Dispute{ID: "1", PartnerID: "p", OrderID: "o", UserID: "u"})
	_ = resolver.FileDispute(&service.Dispute{ID: "2", PartnerID: "p", OrderID: "o2", UserID: "u"})
	_ = resolver.FileDispute(&service.Dispute{ID: "3", PartnerID: "p", OrderID: "o3", UserID: "u"})
	_ = resolver.AssignMediator("2", "m")
	_ = resolver.ResolveDispute("3", "환불")

	counts := resolver.CountByStatus()
	if counts[service.DisputeOpen] != 1 {
		t.Errorf("open = %d, want 1", counts[service.DisputeOpen])
	}
	if counts[service.DisputeMediating] != 1 {
		t.Errorf("mediating = %d, want 1", counts[service.DisputeMediating])
	}
	if counts[service.DisputeResolved] != 1 {
		t.Errorf("resolved = %d, want 1", counts[service.DisputeResolved])
	}
}
