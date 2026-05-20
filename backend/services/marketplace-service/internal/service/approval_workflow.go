// Package service: 마켓플레이스 6단계 파트너 승인 워크플로우
//
// 신청 → 서류심사 → 실사 → 승인/거부 → 보류 → 활성화
// 각 단계는 검토자, 검토 시각, 검토 결과를 기록합니다.
package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// ============================================================================
// 워크플로우 상태 (6단계)
// ============================================================================

const (
	WorkflowSubmitted     = "submitted"      // 1. 신청 제출
	WorkflowReviewing     = "reviewing"      // 2. 서류 심사 중
	WorkflowInspecting    = "inspecting"     // 3. 실사 진행 중
	WorkflowApproved      = "approved"       // 4. 승인 완료
	WorkflowRejected      = "rejected"       // 5. 거부
	WorkflowOnHold        = "on_hold"        // 6. 보류 (서류 보완 필요)
)

// ============================================================================
// 우선순위
// ============================================================================

const (
	PriorityLow    = 1
	PriorityNormal = 2
	PriorityHigh   = 3
	PriorityUrgent = 4
)

// ============================================================================
// 도메인 모델
// ============================================================================

// ApprovalRecord는 승인 워크플로우 기록입니다.
type ApprovalRecord struct {
	ID         string
	PartnerID  string
	Stage      string
	Reviewer   string
	Decision   string // "pass", "fail", "needs_revision"
	Notes      string
	ReviewedAt time.Time
}

// PartnerApplication은 파트너 신청서입니다.
type PartnerApplication struct {
	ID                string
	PartnerID         string
	Stage             string
	Priority          int
	BusinessNumber    string
	BusinessLicense   string
	BankAccount       string
	TaxID             string
	Documents         []string // 서류 URL 목록
	History           []*ApprovalRecord
	SubmittedAt       time.Time
	LastUpdatedAt     time.Time
	EstimatedDecision time.Time // 예상 결정일
}

// ============================================================================
// 검토 큐
// ============================================================================

// ReviewQueue는 우선순위 기반 검토 대기열입니다.
type ReviewQueue struct {
	mu           sync.RWMutex
	applications map[string]*PartnerApplication
}

// NewReviewQueue는 새 검토 대기열을 생성합니다.
func NewReviewQueue() *ReviewQueue {
	return &ReviewQueue{
		applications: make(map[string]*PartnerApplication),
	}
}

// Enqueue는 신청서를 큐에 추가합니다.
func (q *ReviewQueue) Enqueue(app *PartnerApplication) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.applications[app.ID] = app
}

// Dequeue는 우선순위 + 제출 시각 기준으로 다음 검토 대상을 반환합니다.
func (q *ReviewQueue) Dequeue() *PartnerApplication {
	q.mu.Lock()
	defer q.mu.Unlock()

	var best *PartnerApplication
	for _, app := range q.applications {
		if app.Stage != WorkflowSubmitted && app.Stage != WorkflowOnHold {
			continue
		}
		if best == nil {
			best = app
			continue
		}
		// 우선순위가 더 높거나, 같으면 더 일찍 제출된 것 우선
		if app.Priority > best.Priority ||
			(app.Priority == best.Priority && app.SubmittedAt.Before(best.SubmittedAt)) {
			best = app
		}
	}

	if best != nil {
		delete(q.applications, best.ID)
	}
	return best
}

// Size는 큐의 크기를 반환합니다.
func (q *ReviewQueue) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.applications)
}

// PendingByPriority는 우선순위별 대기 건수를 반환합니다.
func (q *ReviewQueue) PendingByPriority() map[int]int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	result := make(map[int]int)
	for _, app := range q.applications {
		result[app.Priority]++
	}
	return result
}

// ============================================================================
// 승인 워크플로우 엔진
// ============================================================================

// ApprovalWorkflowEngine은 6단계 파트너 승인 워크플로우를 관리합니다.
type ApprovalWorkflowEngine struct {
	mu      sync.RWMutex
	apps    map[string]*PartnerApplication
	queue   *ReviewQueue
	transitions map[string][]string // 허용된 전이
}

// NewApprovalWorkflowEngine은 새 워크플로우 엔진을 생성합니다.
func NewApprovalWorkflowEngine() *ApprovalWorkflowEngine {
	return &ApprovalWorkflowEngine{
		apps:  make(map[string]*PartnerApplication),
		queue: NewReviewQueue(),
		transitions: map[string][]string{
			WorkflowSubmitted:  {WorkflowReviewing, WorkflowOnHold, WorkflowRejected},
			WorkflowReviewing:  {WorkflowInspecting, WorkflowOnHold, WorkflowRejected},
			WorkflowInspecting: {WorkflowApproved, WorkflowOnHold, WorkflowRejected},
			WorkflowOnHold:     {WorkflowReviewing, WorkflowRejected},
			WorkflowApproved:   {}, // 종착
			WorkflowRejected:   {}, // 종착
		},
	}
}

// SubmitApplication은 신청서를 제출하여 워크플로우를 시작합니다.
func (e *ApprovalWorkflowEngine) SubmitApplication(ctx context.Context, app *PartnerApplication) error {
	if app == nil || app.PartnerID == "" {
		return errors.New("partner_id required")
	}
	if app.BusinessNumber == "" {
		return errors.New("business_number required")
	}
	if app.BusinessLicense == "" {
		return errors.New("business_license required")
	}

	now := time.Now().UTC()
	app.Stage = WorkflowSubmitted
	app.SubmittedAt = now
	app.LastUpdatedAt = now
	if app.Priority == 0 {
		app.Priority = PriorityNormal
	}
	// 5영업일 후 결정 예상
	app.EstimatedDecision = now.AddDate(0, 0, 5)

	e.mu.Lock()
	e.apps[app.ID] = app
	e.mu.Unlock()

	e.queue.Enqueue(app)
	return nil
}

// AdvanceStage는 워크플로우를 다음 단계로 진행시킵니다.
func (e *ApprovalWorkflowEngine) AdvanceStage(ctx context.Context, appID, nextStage, reviewer, notes string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	app, ok := e.apps[appID]
	if !ok {
		return fmt.Errorf("application %s not found", appID)
	}

	allowed, ok := e.transitions[app.Stage]
	if !ok {
		return fmt.Errorf("no transitions from %s (terminal stage)", app.Stage)
	}

	valid := false
	for _, s := range allowed {
		if s == nextStage {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid transition: %s → %s", app.Stage, nextStage)
	}

	now := time.Now().UTC()
	record := &ApprovalRecord{
		ID:         fmt.Sprintf("ar-%s-%d", appID, now.UnixNano()),
		PartnerID:  app.PartnerID,
		Stage:      nextStage,
		Reviewer:   reviewer,
		Notes:      notes,
		ReviewedAt: now,
		Decision:   decisionForStage(nextStage),
	}

	app.History = append(app.History, record)
	app.Stage = nextStage
	app.LastUpdatedAt = now

	// on_hold 또는 reviewing 상태면 큐에 다시 넣기
	if nextStage == WorkflowOnHold || nextStage == WorkflowReviewing {
		e.queue.Enqueue(app)
	}

	return nil
}

func decisionForStage(stage string) string {
	switch stage {
	case WorkflowApproved:
		return "pass"
	case WorkflowRejected:
		return "fail"
	case WorkflowOnHold:
		return "needs_revision"
	default:
		return "in_progress"
	}
}

// GetApplication은 신청서를 조회합니다.
func (e *ApprovalWorkflowEngine) GetApplication(appID string) (*PartnerApplication, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	app, ok := e.apps[appID]
	if !ok {
		return nil, fmt.Errorf("application %s not found", appID)
	}
	return app, nil
}

// GetNextForReview는 검토 큐에서 다음 검토 대상을 가져옵니다.
func (e *ApprovalWorkflowEngine) GetNextForReview() *PartnerApplication {
	return e.queue.Dequeue()
}

// CountByStage는 단계별 신청서 건수를 반환합니다.
func (e *ApprovalWorkflowEngine) CountByStage() map[string]int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	counts := make(map[string]int)
	for _, app := range e.apps {
		counts[app.Stage]++
	}
	return counts
}

// ============================================================================
// 정산 엔진
// ============================================================================

// SettlementRecord는 정산 명세서입니다.
type SettlementRecord struct {
	ID            string
	PartnerID     string
	PeriodStart   time.Time
	PeriodEnd     time.Time
	GrossRevenue  float64
	Commission    float64
	NetPayout     float64
	Currency      string
	Status        string // "pending", "paid", "disputed"
	OrderCount    int
	Items         []*SettlementItem
	GeneratedAt   time.Time
	SettledAt     *time.Time
}

// SettlementItem은 정산 명세 항목입니다.
type SettlementItem struct {
	ProductID  string
	OrderID    string
	GrossPrice float64
	Commission float64
	Net        float64
	OrderedAt  time.Time
}

// SettlementEngine은 월간 정산을 처리합니다.
type SettlementEngine struct {
	defaultCurrency string
}

// NewSettlementEngine은 새 정산 엔진을 생성합니다.
func NewSettlementEngine() *SettlementEngine {
	return &SettlementEngine{defaultCurrency: "KRW"}
}

// GenerateMonthlySettlement은 월간 정산명세서를 생성합니다.
func (s *SettlementEngine) GenerateMonthlySettlement(
	partnerID string,
	year int,
	month int,
	items []*SettlementItem,
) *SettlementRecord {
	periodStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	periodEnd := periodStart.AddDate(0, 1, 0).Add(-time.Second)

	var gross, commission float64
	var orderIDs = make(map[string]bool)

	for _, item := range items {
		gross += item.GrossPrice
		commission += item.Commission
		orderIDs[item.OrderID] = true
	}

	record := &SettlementRecord{
		ID:           fmt.Sprintf("st-%s-%d-%02d", partnerID, year, month),
		PartnerID:    partnerID,
		PeriodStart:  periodStart,
		PeriodEnd:    periodEnd,
		GrossRevenue: gross,
		Commission:   commission,
		NetPayout:    gross - commission,
		Currency:     s.defaultCurrency,
		Status:       "pending",
		OrderCount:   len(orderIDs),
		Items:        items,
		GeneratedAt:  time.Now().UTC(),
	}

	return record
}

// MarkSettled는 정산 완료로 표시합니다.
func (s *SettlementEngine) MarkSettled(record *SettlementRecord) error {
	if record == nil {
		return errors.New("record is nil")
	}
	if record.Status == "paid" {
		return errors.New("already settled")
	}
	now := time.Now().UTC()
	record.Status = "paid"
	record.SettledAt = &now
	return nil
}

// CommissionForCategory는 카테고리에 따른 수수료를 계산합니다 (할인 적용).
func (s *SettlementEngine) CommissionForCategory(price float64, category string, partnerTier string) float64 {
	rate, ok := defaultCommissionRates[category]
	if !ok {
		rate = 0.20
	}

	// 파트너 티어별 할인
	discount := 0.0
	switch partnerTier {
	case "platinum":
		discount = 0.05 // 5%p 할인
	case "gold":
		discount = 0.03
	case "silver":
		discount = 0.01
	}

	finalRate := rate - discount
	if finalRate < 0.05 {
		finalRate = 0.05 // 최소 5%
	}
	return price * finalRate
}

// ============================================================================
// 분쟁 해결
// ============================================================================

// DisputeStatus는 분쟁 상태입니다.
const (
	DisputeOpen       = "open"
	DisputeMediating  = "mediating"
	DisputeResolved   = "resolved"
	DisputeEscalated  = "escalated"
)

// Dispute는 분쟁 신고입니다.
type Dispute struct {
	ID            string
	PartnerID     string
	OrderID       string
	UserID        string
	Reason        string
	Description   string
	Status        string
	Resolution    string
	OpenedAt      time.Time
	ResolvedAt    *time.Time
	Mediator      string
}

// DisputeResolver는 분쟁 해결을 처리합니다.
type DisputeResolver struct {
	mu       sync.RWMutex
	disputes map[string]*Dispute
}

// NewDisputeResolver는 새 분쟁 해결자를 생성합니다.
func NewDisputeResolver() *DisputeResolver {
	return &DisputeResolver{disputes: make(map[string]*Dispute)}
}

// FileDispute는 새 분쟁을 등록합니다.
func (dr *DisputeResolver) FileDispute(d *Dispute) error {
	if d == nil {
		return errors.New("dispute is nil")
	}
	if d.ID == "" {
		d.ID = fmt.Sprintf("dp-%d", time.Now().UnixNano())
	}
	d.Status = DisputeOpen
	d.OpenedAt = time.Now().UTC()

	dr.mu.Lock()
	defer dr.mu.Unlock()
	dr.disputes[d.ID] = d
	return nil
}

// AssignMediator는 중재자를 배정하여 분쟁을 진행시킵니다.
func (dr *DisputeResolver) AssignMediator(disputeID, mediator string) error {
	dr.mu.Lock()
	defer dr.mu.Unlock()

	d, ok := dr.disputes[disputeID]
	if !ok {
		return fmt.Errorf("dispute %s not found", disputeID)
	}
	if d.Status != DisputeOpen {
		return fmt.Errorf("can only assign to open disputes (current: %s)", d.Status)
	}
	d.Status = DisputeMediating
	d.Mediator = mediator
	return nil
}

// ResolveDispute는 분쟁을 해결합니다.
func (dr *DisputeResolver) ResolveDispute(disputeID, resolution string) error {
	dr.mu.Lock()
	defer dr.mu.Unlock()

	d, ok := dr.disputes[disputeID]
	if !ok {
		return fmt.Errorf("dispute %s not found", disputeID)
	}
	if d.Status == DisputeResolved {
		return errors.New("already resolved")
	}
	now := time.Now().UTC()
	d.Status = DisputeResolved
	d.Resolution = resolution
	d.ResolvedAt = &now
	return nil
}

// EscalateDispute는 분쟁을 상위로 에스컬레이션합니다.
func (dr *DisputeResolver) EscalateDispute(disputeID, reason string) error {
	dr.mu.Lock()
	defer dr.mu.Unlock()

	d, ok := dr.disputes[disputeID]
	if !ok {
		return fmt.Errorf("dispute %s not found", disputeID)
	}
	d.Status = DisputeEscalated
	if d.Description == "" {
		d.Description = reason
	} else {
		d.Description += "\n[ESCALATED] " + reason
	}
	return nil
}

// GetDispute는 분쟁을 조회합니다.
func (dr *DisputeResolver) GetDispute(id string) (*Dispute, error) {
	dr.mu.RLock()
	defer dr.mu.RUnlock()
	d, ok := dr.disputes[id]
	if !ok {
		return nil, fmt.Errorf("dispute %s not found", id)
	}
	return d, nil
}

// CountByStatus는 상태별 분쟁 건수를 반환합니다.
func (dr *DisputeResolver) CountByStatus() map[string]int {
	dr.mu.RLock()
	defer dr.mu.RUnlock()
	counts := make(map[string]int)
	for _, d := range dr.disputes {
		counts[d.Status]++
	}
	return counts
}
