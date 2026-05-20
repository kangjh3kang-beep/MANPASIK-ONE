package service

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

// fakeSubRepo는 테스트용 인메모리 구독 저장소입니다.
type fakeSubRepo struct {
	byUserID map[string]*Subscription
	byID     map[string]*Subscription
}

func newFakeSubRepo() *fakeSubRepo {
	return &fakeSubRepo{
		byUserID: make(map[string]*Subscription),
		byID:     make(map[string]*Subscription),
	}
}

func (r *fakeSubRepo) GetByUserID(_ context.Context, userID string) (*Subscription, error) {
	s, ok := r.byUserID[userID]
	if !ok {
		return nil, nil
	}
	return s, nil
}

func (r *fakeSubRepo) GetByID(_ context.Context, id string) (*Subscription, error) {
	s, ok := r.byID[id]
	if !ok {
		return nil, nil
	}
	return s, nil
}

func (r *fakeSubRepo) Create(_ context.Context, sub *Subscription) error {
	r.byUserID[sub.UserID] = sub
	r.byID[sub.ID] = sub
	return nil
}

func (r *fakeSubRepo) Update(_ context.Context, sub *Subscription) error {
	r.byUserID[sub.UserID] = sub
	r.byID[sub.ID] = sub
	return nil
}

// fakeHistoryRepo는 테스트용 이력 저장소입니다.
type fakeHistoryRepo struct {
	entries []*SubscriptionHistoryEntry
}

func newFakeHistoryRepo() *fakeHistoryRepo {
	return &fakeHistoryRepo{
		entries: make([]*SubscriptionHistoryEntry, 0),
	}
}

func (r *fakeHistoryRepo) Record(_ context.Context, entry *SubscriptionHistoryEntry) error {
	r.entries = append(r.entries, entry)
	return nil
}

func (r *fakeHistoryRepo) ListByUserID(_ context.Context, userID string, limit, offset int32) ([]*SubscriptionHistoryEntry, error) {
	var result []*SubscriptionHistoryEntry
	for _, e := range r.entries {
		if e.UserID == userID {
			result = append(result, e)
		}
	}
	return result, nil
}

func newTestService() (*SubscriptionService, *fakeSubRepo, *fakeHistoryRepo) {
	subRepo := newFakeSubRepo()
	historyRepo := newFakeHistoryRepo()
	svc := NewSubscriptionService(zap.NewNop(), subRepo, historyRepo)
	return svc, subRepo, historyRepo
}

func TestCreateSubscription_Free(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	sub, err := svc.CreateSubscription(ctx, "user-1", TierFree)
	if err != nil {
		t.Fatalf("CreateSubscription Free 실패: %v", err)
	}
	if sub.Tier != TierFree {
		t.Errorf("티어: got %d, want %d", sub.Tier, TierFree)
	}
	if sub.MaxDevices != 1 {
		t.Errorf("MaxDevices: got %d, want 1", sub.MaxDevices)
	}
	if sub.MonthlyPriceKRW != 0 {
		t.Errorf("MonthlyPriceKRW: got %d, want 0", sub.MonthlyPriceKRW)
	}
	if sub.AutoRenew != false {
		t.Error("Free 티어는 AutoRenew=false여야 합니다")
	}
}

func TestCreateSubscription_Pro(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	sub, err := svc.CreateSubscription(ctx, "user-2", TierPro)
	if err != nil {
		t.Fatalf("CreateSubscription Pro 실패: %v", err)
	}
	if sub.Tier != TierPro {
		t.Errorf("티어: got %d, want %d", sub.Tier, TierPro)
	}
	if sub.MaxDevices != 5 {
		t.Errorf("MaxDevices: got %d, want 5", sub.MaxDevices)
	}
	if sub.AICoachingEnabled != true {
		t.Error("Pro 티어는 AI 코칭이 활성화되어야 합니다")
	}
	if sub.MonthlyPriceKRW != 29900 {
		t.Errorf("MonthlyPriceKRW: got %d, want 29900", sub.MonthlyPriceKRW)
	}
}

func TestCreateSubscription_Duplicate(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	_, err := svc.CreateSubscription(ctx, "user-dup", TierFree)
	if err != nil {
		t.Fatalf("첫 번째 구독 생성 실패: %v", err)
	}

	_, err = svc.CreateSubscription(ctx, "user-dup", TierBasic)
	if err == nil {
		t.Error("중복 구독 생성이 허용되었습니다")
	}
}

func TestGetSubscription(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	_, _ = svc.CreateSubscription(ctx, "user-get", TierBasic)

	sub, err := svc.GetSubscription(ctx, "user-get")
	if err != nil {
		t.Fatalf("GetSubscription 실패: %v", err)
	}
	if sub.Tier != TierBasic {
		t.Errorf("티어: got %d, want %d", sub.Tier, TierBasic)
	}
}

func TestGetSubscription_NotFound(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	_, err := svc.GetSubscription(ctx, "non-existent")
	if err == nil {
		t.Error("존재하지 않는 구독 조회가 성공했습니다")
	}
}

func TestUpdateSubscription_Upgrade(t *testing.T) {
	svc, _, historyRepo := newTestService()
	ctx := context.Background()

	_, _ = svc.CreateSubscription(ctx, "user-up", TierBasic)

	sub, err := svc.UpdateSubscription(ctx, "user-up", TierPro)
	if err != nil {
		t.Fatalf("UpdateSubscription 실패: %v", err)
	}
	if sub.Tier != TierPro {
		t.Errorf("업그레이드 후 티어: got %d, want %d", sub.Tier, TierPro)
	}
	if sub.MaxDevices != 5 {
		t.Errorf("업그레이드 후 MaxDevices: got %d, want 5", sub.MaxDevices)
	}

	// 이력 확인 (create + upgrade = 2개)
	if len(historyRepo.entries) != 2 {
		t.Errorf("이력 수: got %d, want 2", len(historyRepo.entries))
	}
	if historyRepo.entries[1].Action != "upgrade" {
		t.Errorf("이력 액션: got %s, want upgrade", historyRepo.entries[1].Action)
	}
}

func TestUpdateSubscription_Downgrade(t *testing.T) {
	svc, _, historyRepo := newTestService()
	ctx := context.Background()

	_, _ = svc.CreateSubscription(ctx, "user-down", TierClinical)

	sub, err := svc.UpdateSubscription(ctx, "user-down", TierBasic)
	if err != nil {
		t.Fatalf("UpdateSubscription 다운그레이드 실패: %v", err)
	}
	if sub.Tier != TierBasic {
		t.Errorf("다운그레이드 후 티어: got %d, want %d", sub.Tier, TierBasic)
	}
	if sub.TelemedicineEnabled != false {
		t.Error("Basic 티어는 화상진료가 비활성화되어야 합니다")
	}

	if historyRepo.entries[1].Action != "downgrade" {
		t.Errorf("이력 액션: got %s, want downgrade", historyRepo.entries[1].Action)
	}
}

func TestUpdateSubscription_SameTier(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	_, _ = svc.CreateSubscription(ctx, "user-same", TierBasic)

	_, err := svc.UpdateSubscription(ctx, "user-same", TierBasic)
	if err == nil {
		t.Error("동일 티어 변경이 허용되었습니다")
	}
}

func TestCancelSubscription(t *testing.T) {
	svc, _, historyRepo := newTestService()
	ctx := context.Background()

	_, _ = svc.CreateSubscription(ctx, "user-cancel", TierPro)

	sub, err := svc.CancelSubscription(ctx, "user-cancel", "더 이상 필요 없음")
	if err != nil {
		t.Fatalf("CancelSubscription 실패: %v", err)
	}
	if sub.Status != StatusCancelled {
		t.Errorf("해지 후 상태: got %d, want %d", sub.Status, StatusCancelled)
	}
	if sub.CancelledAt == nil {
		t.Error("CancelledAt이 설정되지 않았습니다")
	}
	if sub.AutoRenew != false {
		t.Error("해지 후 AutoRenew=false여야 합니다")
	}

	// cancel 이력 확인
	found := false
	for _, e := range historyRepo.entries {
		if e.Action == "cancel" && e.Reason == "더 이상 필요 없음" {
			found = true
		}
	}
	if !found {
		t.Error("해지 이력이 기록되지 않았습니다")
	}
}

func TestCancelSubscription_AlreadyCancelled(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	_, _ = svc.CreateSubscription(ctx, "user-cc", TierBasic)
	_, _ = svc.CancelSubscription(ctx, "user-cc", "test")

	_, err := svc.CancelSubscription(ctx, "user-cc", "다시 해지")
	if err == nil {
		t.Error("이미 해지된 구독의 재해지가 허용되었습니다")
	}
}

func TestCheckFeatureAccess_Allowed(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	_, _ = svc.CreateSubscription(ctx, "user-feat", TierPro)

	allowed, _, _, _ := svc.CheckFeatureAccess(ctx, "user-feat", "ai_coaching")
	if !allowed {
		t.Error("Pro 티어는 ai_coaching에 접근 가능해야 합니다")
	}
}

func TestCheckFeatureAccess_Denied(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	_, _ = svc.CreateSubscription(ctx, "user-deny", TierBasic)

	allowed, requiredTier, currentTier, msg := svc.CheckFeatureAccess(ctx, "user-deny", "ai_coaching")
	if allowed {
		t.Error("Basic 티어는 ai_coaching에 접근 불가해야 합니다")
	}
	if requiredTier != TierPro {
		t.Errorf("필요 티어: got %d, want %d", requiredTier, TierPro)
	}
	if currentTier != TierBasic {
		t.Errorf("현재 티어: got %d, want %d", currentTier, TierBasic)
	}
	if msg == "" {
		t.Error("거부 메시지가 비어 있습니다")
	}
}

func TestCheckFeatureAccess_UnknownFeature(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	_, _ = svc.CreateSubscription(ctx, "user-unk", TierClinical)

	allowed, _, _, _ := svc.CheckFeatureAccess(ctx, "user-unk", "unknown_feature")
	if allowed {
		t.Error("알 수 없는 기능에 대한 접근이 허용되었습니다")
	}
}

func TestListSubscriptionPlans(t *testing.T) {
	svc, _, _ := newTestService()

	planList := svc.ListSubscriptionPlans()
	if len(planList) != 4 {
		t.Fatalf("플랜 수: got %d, want 4", len(planList))
	}

	// 순서 확인 (Free → Basic → Pro → Clinical)
	expectedTiers := []SubscriptionTier{TierFree, TierBasic, TierPro, TierClinical}
	for i, plan := range planList {
		if plan.Tier != expectedTiers[i] {
			t.Errorf("플랜[%d] 티어: got %d, want %d", i, plan.Tier, expectedTiers[i])
		}
	}

	// Clinical 가격 확인
	if planList[3].MonthlyPriceKRW != 59900 {
		t.Errorf("Clinical 가격: got %d, want 59900", planList[3].MonthlyPriceKRW)
	}
}

// =============================================================================
// Phase F 테스트 보강: 카트리지 접근 제어 + 엣지 케이스
// =============================================================================

func TestCheckCartridgeAccess_FreeGlucose(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()
	svc.CreateSubscription(ctx, "user-free", TierFree)

	// Free: HealthBiomarker Glucose (cat=1, type=1) = LIMITED (일 3회)
	result := svc.CheckCartridgeAccess(ctx, "user-free", 1, 1)
	if !result.Allowed {
		t.Fatal("Free 티어에서 Glucose는 Limited로 접근 가능해야 함")
	}
	if result.AccessLevel != AccessLimited {
		t.Errorf("AccessLevel: got %s, want %s", result.AccessLevel, AccessLimited)
	}
	if result.RemainingDaily != 3 {
		t.Errorf("RemainingDaily: got %d, want 3", result.RemainingDaily)
	}
}

func TestCheckCartridgeAccess_FreeRestricted(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()
	svc.CreateSubscription(ctx, "user-free2", TierFree)

	// Free: Environmental (cat=2) = RESTRICTED (글로벌 기본값 적용)
	result := svc.CheckCartridgeAccess(ctx, "user-free2", 2, 1)
	if result.Allowed {
		t.Fatal("Free 티어에서 Environmental은 접근 불가해야 함")
	}
	if result.AccessLevel != AccessRestricted {
		t.Errorf("AccessLevel: got %s, want %s", result.AccessLevel, AccessRestricted)
	}
}

func TestCheckCartridgeAccess_BasicIncluded(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()
	svc.CreateSubscription(ctx, "user-basic", TierBasic)

	// Basic: HealthBiomarker 전체 (cat=1) = INCLUDED
	result := svc.CheckCartridgeAccess(ctx, "user-basic", 1, 5)
	if !result.Allowed {
		t.Fatal("Basic 티어에서 HealthBiomarker는 접근 가능해야 함")
	}
	if result.AccessLevel != AccessIncluded {
		t.Errorf("AccessLevel: got %s, want %s", result.AccessLevel, AccessIncluded)
	}
	if result.RemainingDaily != -1 {
		t.Errorf("RemainingDaily: got %d, want -1 (무제한)", result.RemainingDaily)
	}
}

func TestCheckCartridgeAccess_BasicAddOn(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()
	svc.CreateSubscription(ctx, "user-basic2", TierBasic)

	// Basic: Environmental (cat=2) = ADD_ON
	result := svc.CheckCartridgeAccess(ctx, "user-basic2", 2, 1)
	if result.Allowed {
		t.Fatal("Basic 티어에서 Environmental은 ADD_ON (구매 필요)")
	}
	if result.AccessLevel != AccessAddOn {
		t.Errorf("AccessLevel: got %s, want %s", result.AccessLevel, AccessAddOn)
	}
}

func TestCheckCartridgeAccess_ProIncluded(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()
	svc.CreateSubscription(ctx, "user-pro", TierPro)

	// Pro: Environmental (cat=2) = INCLUDED
	result := svc.CheckCartridgeAccess(ctx, "user-pro", 2, 1)
	if !result.Allowed {
		t.Fatal("Pro 티어에서 Environmental은 접근 가능해야 함")
	}
	if result.AccessLevel != AccessIncluded {
		t.Errorf("AccessLevel: got %s, want %s", result.AccessLevel, AccessIncluded)
	}
}

func TestCheckCartridgeAccess_ClinicalAll(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()
	svc.CreateSubscription(ctx, "user-clin", TierClinical)

	// Clinical: 글로벌 기본값 = INCLUDED
	result := svc.CheckCartridgeAccess(ctx, "user-clin", 5, 1)
	if !result.Allowed {
		t.Fatal("Clinical 티어에서 모든 카트리지 접근 가능해야 함")
	}
	if result.AccessLevel != AccessIncluded {
		t.Errorf("AccessLevel: got %s, want %s", result.AccessLevel, AccessIncluded)
	}
}

func TestCheckCartridgeAccess_NoSubscription(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	result := svc.CheckCartridgeAccess(ctx, "no-user", 1, 1)
	if result.Allowed {
		t.Fatal("구독 없는 사용자는 접근 불가해야 함")
	}
}

func TestListAccessibleCartridges_Free(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()
	svc.CreateSubscription(ctx, "user-list-free", TierFree)

	entries := svc.ListAccessibleCartridges(ctx, "user-list-free")
	if len(entries) == 0 {
		t.Fatal("카트리지 목록이 비어있으면 안 됨")
	}

	// Free: Glucose는 Limited
	var glucoseEntry *CartridgeAccessEntry
	for _, e := range entries {
		if e.CategoryCode == 1 && e.TypeIndex == 1 {
			glucoseEntry = e
			break
		}
	}
	if glucoseEntry == nil {
		t.Fatal("Glucose 카트리지가 목록에 없음")
	}
	if glucoseEntry.AccessLevel != AccessLimited {
		t.Errorf("Free Glucose AccessLevel: got %s, want %s", glucoseEntry.AccessLevel, AccessLimited)
	}
}

func TestListAccessibleCartridges_Clinical(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()
	svc.CreateSubscription(ctx, "user-list-clin", TierClinical)

	entries := svc.ListAccessibleCartridges(ctx, "user-list-clin")
	includedCount := 0
	for _, e := range entries {
		if e.AccessLevel == AccessIncluded {
			includedCount++
		}
	}
	// Clinical은 대부분 INCLUDED
	if includedCount < 10 {
		t.Errorf("Clinical 포함 카트리지가 10개 이상이어야 함: got %d", includedCount)
	}
}

func TestGetSubscription_EmptyUserID(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()
	_, err := svc.GetSubscription(ctx, "")
	if err == nil {
		t.Fatal("빈 user_id에 에러가 반환되어야 함")
	}
}

func TestUpdateSubscription_NotFound(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()
	_, err := svc.UpdateSubscription(ctx, "non-existent", TierPro)
	if err == nil {
		t.Fatal("존재하지 않는 구독 업데이트가 성공했습니다")
	}
}

func TestCancelSubscription_NotFound(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()
	_, err := svc.CancelSubscription(ctx, "non-existent", "이유")
	if err == nil {
		t.Fatal("존재하지 않는 구독 해지가 성공했습니다")
	}
}

func TestCreateSubscription_Clinical(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	sub, err := svc.CreateSubscription(ctx, "user-clin", TierClinical)
	if err != nil {
		t.Fatalf("CreateSubscription Clinical 실패: %v", err)
	}
	if sub.Tier != TierClinical {
		t.Errorf("티어: got %d, want %d", sub.Tier, TierClinical)
	}
	if !sub.TelemedicineEnabled {
		t.Error("Clinical 티어는 화상진료가 활성화되어야 합니다")
	}
	if sub.MaxDevices != 10 {
		t.Errorf("MaxDevices: got %d, want 10", sub.MaxDevices)
	}
	if sub.MonthlyPriceKRW != 59900 {
		t.Errorf("MonthlyPriceKRW: got %d, want 59900", sub.MonthlyPriceKRW)
	}
}

func TestUpdateSubscription_CancelledSubscription(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	svc.CreateSubscription(ctx, "user-canup", TierBasic)
	svc.CancelSubscription(ctx, "user-canup", "test")

	_, err := svc.UpdateSubscription(ctx, "user-canup", TierPro)
	if err == nil {
		t.Fatal("해지된 구독의 업그레이드가 허용되었습니다")
	}
}

func TestTierToString(t *testing.T) {
	tests := []struct {
		tier SubscriptionTier
		want string
	}{
		{TierFree, "free"},
		{TierBasic, "basic"},
		{TierPro, "pro"},
		{TierClinical, "clinical"},
		{SubscriptionTier(99), "unknown"},
	}
	for _, tt := range tests {
		got := tierToString(tt.tier)
		if got != tt.want {
			t.Errorf("tierToString(%d) = %s, want %s", tt.tier, got, tt.want)
		}
	}
}
