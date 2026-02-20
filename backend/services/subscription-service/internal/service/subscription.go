// Package service는 subscription-service의 비즈니스 로직을 구현합니다.
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/manpasik/backend/shared/errors"
	"go.uber.org/zap"
)

// SubscriptionTier는 구독 티어입니다.
type SubscriptionTier int32

const (
	TierFree     SubscriptionTier = 0
	TierBasic    SubscriptionTier = 1
	TierPro      SubscriptionTier = 2
	TierClinical SubscriptionTier = 3
)

// SubscriptionStatus는 구독 상태입니다.
type SubscriptionStatus int32

const (
	StatusUnknown   SubscriptionStatus = 0
	StatusActive    SubscriptionStatus = 1
	StatusCancelled SubscriptionStatus = 2
	StatusExpired   SubscriptionStatus = 3
	StatusSuspended SubscriptionStatus = 4
	StatusTrial     SubscriptionStatus = 5
)

// Subscription은 구독 엔티티입니다.
type Subscription struct {
	ID                  string
	UserID              string
	Tier                SubscriptionTier
	Status              SubscriptionStatus
	StartedAt           time.Time
	ExpiresAt           time.Time
	CancelledAt         *time.Time
	MaxDevices          int32
	MaxFamilyMembers    int32
	AICoachingEnabled   bool
	TelemedicineEnabled bool
	MonthlyPriceKRW     int32
	AutoRenew           bool
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// SubscriptionPlan은 구독 플랜 정의입니다.
type SubscriptionPlan struct {
	Tier                SubscriptionTier
	Name                string
	Description         string
	MonthlyPriceKRW     int32
	MaxDevices          int32
	MaxFamilyMembers    int32
	AICoachingEnabled   bool
	TelemedicineEnabled bool
	Features            []string
}

// SubscriptionRepository는 구독 데이터 저장소 인터페이스입니다.
type SubscriptionRepository interface {
	GetByUserID(ctx context.Context, userID string) (*Subscription, error)
	GetByID(ctx context.Context, id string) (*Subscription, error)
	Create(ctx context.Context, sub *Subscription) error
	Update(ctx context.Context, sub *Subscription) error
}

// SubscriptionHistoryRepository는 구독 변경 이력 저장소입니다.
type SubscriptionHistoryRepository interface {
	Record(ctx context.Context, entry *SubscriptionHistoryEntry) error
	ListByUserID(ctx context.Context, userID string, limit, offset int32) ([]*SubscriptionHistoryEntry, error)
}

// SubscriptionHistoryEntry는 구독 변경 이력 항목입니다.
type SubscriptionHistoryEntry struct {
	ID        string
	UserID    string
	OldTier   SubscriptionTier
	NewTier   SubscriptionTier
	Action    string // "create", "upgrade", "downgrade", "cancel", "renew"
	Reason    string
	CreatedAt time.Time
}

// ============================================================================
// 이벤트 타입 (Kafka 발행용)
// ============================================================================

// SubscriptionChangedEvent는 구독 변경 이벤트입니다.
type SubscriptionChangedEvent struct {
	SubscriptionID      string
	UserID              string
	PreviousTier        string
	NewTier             string
	ChangeType          string // "create", "upgrade", "downgrade", "cancel"
	EffectiveAt         time.Time
	MaxDevices          int32
	MaxFamilyMembers    int32
	AICoachingEnabled   bool
	TelemedicineEnabled bool
}

// EventPublisher는 이벤트 발행 인터페이스입니다 (Kafka).
type EventPublisher interface {
	PublishSubscriptionChanged(ctx context.Context, event *SubscriptionChangedEvent) error
}

// SubscriptionService는 구독 비즈니스 로직입니다.
type SubscriptionService struct {
	logger         *zap.Logger
	subRepo        SubscriptionRepository
	historyRepo    SubscriptionHistoryRepository
	eventPublisher EventPublisher
}

// NewSubscriptionService는 새 SubscriptionService를 생성합니다.
func NewSubscriptionService(
	logger *zap.Logger,
	subRepo SubscriptionRepository,
	historyRepo SubscriptionHistoryRepository,
) *SubscriptionService {
	return &SubscriptionService{
		logger:      logger,
		subRepo:     subRepo,
		historyRepo: historyRepo,
	}
}

// SetEventPublisher는 이벤트 발행기를 설정합니다 (optional).
func (s *SubscriptionService) SetEventPublisher(ep EventPublisher) {
	s.eventPublisher = ep
}

// 플랜 정의 (불변 참조)
var plans = map[SubscriptionTier]*SubscriptionPlan{
	TierFree: {
		Tier:                TierFree,
		Name:                "Free",
		Description:         "기본 무료 플랜",
		MonthlyPriceKRW:     0,
		MaxDevices:          1,
		MaxFamilyMembers:    0,
		AICoachingEnabled:   false,
		TelemedicineEnabled: false,
		Features:            []string{"기본 측정", "측정 이력 조회", "단일 리더기"},
	},
	TierBasic: {
		Tier:                TierBasic,
		Name:                "Basic Safety",
		Description:         "기본 안전 관리 플랜 (₩9,900/월)",
		MonthlyPriceKRW:     9900,
		MaxDevices:          3,
		MaxFamilyMembers:    2,
		AICoachingEnabled:   false,
		TelemedicineEnabled: false,
		Features:            []string{"기본 측정", "측정 이력 조회", "리더기 3대", "가족 2명", "데이터 내보내기"},
	},
	TierPro: {
		Tier:                TierPro,
		Name:                "Bio-Optimization",
		Description:         "AI 코칭 포함 프로 플랜 (₩29,900/월)",
		MonthlyPriceKRW:     29900,
		MaxDevices:          5,
		MaxFamilyMembers:    5,
		AICoachingEnabled:   true,
		TelemedicineEnabled: false,
		Features:            []string{"기본 측정", "측정 이력 조회", "리더기 5대", "가족 5명", "데이터 내보내기", "AI 건강 코칭", "트렌드 분석", "건강 점수"},
	},
	TierClinical: {
		Tier:                TierClinical,
		Name:                "Clinical Guard",
		Description:         "화상진료 포함 클리니컬 플랜 (₩59,900/월)",
		MonthlyPriceKRW:     59900,
		MaxDevices:          10,
		MaxFamilyMembers:    10,
		AICoachingEnabled:   true,
		TelemedicineEnabled: true,
		Features:            []string{"기본 측정", "측정 이력 조회", "리더기 10대", "가족 10명", "데이터 내보내기", "AI 건강 코칭", "트렌드 분석", "건강 점수", "화상진료", "의료진 매칭", "FHIR 연동"},
	},
}

// CreateSubscription은 새 구독을 생성합니다.
func (s *SubscriptionService) CreateSubscription(ctx context.Context, userID string, tier SubscriptionTier) (*Subscription, error) {
	// 기존 구독 확인
	existing, _ := s.subRepo.GetByUserID(ctx, userID)
	if existing != nil && existing.Status == StatusActive {
		return nil, apperrors.New(apperrors.ErrAlreadyExists, "이미 활성 구독이 있습니다")
	}

	plan, ok := plans[tier]
	if !ok {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "유효하지 않은 구독 티어입니다")
	}

	now := time.Now().UTC()
	sub := &Subscription{
		ID:                  uuid.New().String(),
		UserID:              userID,
		Tier:                tier,
		Status:              StatusActive,
		StartedAt:           now,
		ExpiresAt:           now.AddDate(0, 1, 0), // 1개월
		MaxDevices:          plan.MaxDevices,
		MaxFamilyMembers:    plan.MaxFamilyMembers,
		AICoachingEnabled:   plan.AICoachingEnabled,
		TelemedicineEnabled: plan.TelemedicineEnabled,
		MonthlyPriceKRW:     plan.MonthlyPriceKRW,
		AutoRenew:           tier != TierFree,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	if err := s.subRepo.Create(ctx, sub); err != nil {
		s.logger.Error("구독 생성 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "구독 생성에 실패했습니다")
	}

	// 이력 기록
	_ = s.historyRepo.Record(ctx, &SubscriptionHistoryEntry{
		ID:        uuid.New().String(),
		UserID:    userID,
		OldTier:   TierFree,
		NewTier:   tier,
		Action:    "create",
		CreatedAt: now,
	})

	s.logger.Info("구독 생성 완료", zap.String("user_id", userID), zap.Int32("tier", int32(tier)))

	// 구독 생성 이벤트 발행 (Kafka, optional)
	if s.eventPublisher != nil {
		evt := &SubscriptionChangedEvent{
			SubscriptionID:      sub.ID,
			UserID:              userID,
			PreviousTier:        "free",
			NewTier:             tierToString(tier),
			ChangeType:          "create",
			EffectiveAt:         now,
			MaxDevices:          plan.MaxDevices,
			MaxFamilyMembers:    plan.MaxFamilyMembers,
			AICoachingEnabled:   plan.AICoachingEnabled,
			TelemedicineEnabled: plan.TelemedicineEnabled,
		}
		if err := s.eventPublisher.PublishSubscriptionChanged(ctx, evt); err != nil {
			s.logger.Warn("구독 생성 이벤트 발행 실패 (비치명적)", zap.Error(err))
		}
	}

	return sub, nil
}

// GetSubscription은 사용자의 구독 정보를 조회합니다.
func (s *SubscriptionService) GetSubscription(ctx context.Context, userID string) (*Subscription, error) {
	sub, err := s.subRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("구독 조회 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "구독 조회에 실패했습니다")
	}
	if sub == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "구독 정보가 없습니다")
	}
	return sub, nil
}

// UpdateSubscription은 구독 티어를 변경합니다.
func (s *SubscriptionService) UpdateSubscription(ctx context.Context, userID string, newTier SubscriptionTier) (*Subscription, error) {
	sub, err := s.subRepo.GetByUserID(ctx, userID)
	if err != nil || sub == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "구독 정보가 없습니다")
	}

	if sub.Status != StatusActive {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "활성 상태의 구독만 변경할 수 있습니다")
	}

	if sub.Tier == newTier {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "동일한 티어로는 변경할 수 없습니다")
	}

	plan, ok := plans[newTier]
	if !ok {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "유효하지 않은 구독 티어입니다")
	}

	oldTier := sub.Tier
	now := time.Now().UTC()

	sub.Tier = newTier
	sub.MaxDevices = plan.MaxDevices
	sub.MaxFamilyMembers = plan.MaxFamilyMembers
	sub.AICoachingEnabled = plan.AICoachingEnabled
	sub.TelemedicineEnabled = plan.TelemedicineEnabled
	sub.MonthlyPriceKRW = plan.MonthlyPriceKRW
	sub.AutoRenew = newTier != TierFree
	sub.UpdatedAt = now

	if err := s.subRepo.Update(ctx, sub); err != nil {
		s.logger.Error("구독 업데이트 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "구독 변경에 실패했습니다")
	}

	// 이력 기록
	action := "upgrade"
	if newTier < oldTier {
		action = "downgrade"
	}
	_ = s.historyRepo.Record(ctx, &SubscriptionHistoryEntry{
		ID:        uuid.New().String(),
		UserID:    userID,
		OldTier:   oldTier,
		NewTier:   newTier,
		Action:    action,
		CreatedAt: now,
	})

	s.logger.Info("구독 변경 완료",
		zap.String("user_id", userID),
		zap.Int32("old_tier", int32(oldTier)),
		zap.Int32("new_tier", int32(newTier)),
	)

	// 구독 변경 이벤트 발행 (Kafka, optional)
	if s.eventPublisher != nil {
		evt := &SubscriptionChangedEvent{
			SubscriptionID:      sub.ID,
			UserID:              userID,
			PreviousTier:        tierToString(oldTier),
			NewTier:             tierToString(newTier),
			ChangeType:          action,
			EffectiveAt:         now,
			MaxDevices:          plan.MaxDevices,
			MaxFamilyMembers:    plan.MaxFamilyMembers,
			AICoachingEnabled:   plan.AICoachingEnabled,
			TelemedicineEnabled: plan.TelemedicineEnabled,
		}
		if err := s.eventPublisher.PublishSubscriptionChanged(ctx, evt); err != nil {
			s.logger.Warn("구독 변경 이벤트 발행 실패 (비치명적)", zap.Error(err))
		}
	}

	return sub, nil
}

// CancelSubscription은 구독을 해지합니다.
func (s *SubscriptionService) CancelSubscription(ctx context.Context, userID, reason string) (*Subscription, error) {
	sub, err := s.subRepo.GetByUserID(ctx, userID)
	if err != nil || sub == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "구독 정보가 없습니다")
	}

	if sub.Status != StatusActive {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "활성 상태의 구독만 해지할 수 있습니다")
	}

	now := time.Now().UTC()
	sub.Status = StatusCancelled
	sub.CancelledAt = &now
	sub.AutoRenew = false
	sub.UpdatedAt = now

	if err := s.subRepo.Update(ctx, sub); err != nil {
		s.logger.Error("구독 해지 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "구독 해지에 실패했습니다")
	}

	// 이력 기록
	_ = s.historyRepo.Record(ctx, &SubscriptionHistoryEntry{
		ID:        uuid.New().String(),
		UserID:    userID,
		OldTier:   sub.Tier,
		NewTier:   TierFree,
		Action:    "cancel",
		Reason:    reason,
		CreatedAt: now,
	})

	s.logger.Info("구독 해지 완료", zap.String("user_id", userID))

	// 구독 해지 이벤트 발행 (Kafka, optional)
	if s.eventPublisher != nil {
		evt := &SubscriptionChangedEvent{
			SubscriptionID:      sub.ID,
			UserID:              userID,
			PreviousTier:        tierToString(sub.Tier),
			NewTier:             "free",
			ChangeType:          "cancel",
			EffectiveAt:         now,
			MaxDevices:          0,
			MaxFamilyMembers:    0,
			AICoachingEnabled:   false,
			TelemedicineEnabled: false,
		}
		if err := s.eventPublisher.PublishSubscriptionChanged(ctx, evt); err != nil {
			s.logger.Warn("구독 해지 이벤트 발행 실패 (비치명적)", zap.Error(err))
		}
	}

	return sub, nil
}

// tierToString은 SubscriptionTier를 문자열로 변환합니다.
func tierToString(tier SubscriptionTier) string {
	switch tier {
	case TierFree:
		return "free"
	case TierBasic:
		return "basic"
	case TierPro:
		return "pro"
	case TierClinical:
		return "clinical"
	default:
		return "unknown"
	}
}

// CheckFeatureAccess는 기능 접근 권한을 확인합니다.
func (s *SubscriptionService) CheckFeatureAccess(ctx context.Context, userID, featureName string) (bool, SubscriptionTier, SubscriptionTier, string) {
	sub, err := s.subRepo.GetByUserID(ctx, userID)
	if err != nil || sub == nil {
		return false, TierFree, TierFree, "구독 정보가 없습니다"
	}

	currentTier := sub.Tier

	// 기능별 필요 티어 매핑
	featureTierMap := map[string]SubscriptionTier{
		"basic_measurement":   TierFree,
		"measurement_history": TierFree,
		"data_export":         TierBasic,
		"multi_device":        TierBasic,
		"family_sharing":      TierBasic,
		"ai_coaching":         TierPro,
		"trend_analysis":      TierPro,
		"health_score":        TierPro,
		"telemedicine":        TierClinical,
		"medical_matching":    TierClinical,
		"fhir_integration":    TierClinical,
	}

	requiredTier, ok := featureTierMap[featureName]
	if !ok {
		return false, TierFree, currentTier, "알 수 없는 기능입니다"
	}

	if currentTier >= requiredTier {
		return true, requiredTier, currentTier, ""
	}

	return false, requiredTier, currentTier, "상위 플랜이 필요합니다"
}

// ListSubscriptionPlans는 모든 구독 플랜 목록을 반환합니다.
func (s *SubscriptionService) ListSubscriptionPlans() []*SubscriptionPlan {
	result := make([]*SubscriptionPlan, 0, len(plans))
	for _, tier := range []SubscriptionTier{TierFree, TierBasic, TierPro, TierClinical} {
		if plan, ok := plans[tier]; ok {
			result = append(result, plan)
		}
	}
	return result
}

// ============================================================================
// 카트리지 접근 제어 (Cartridge Access Control)
// ============================================================================

// CartridgeAccessLevel은 카트리지 접근 레벨입니다.
type CartridgeAccessLevel string

const (
	AccessIncluded   CartridgeAccessLevel = "included"
	AccessLimited    CartridgeAccessLevel = "limited"
	AccessAddOn      CartridgeAccessLevel = "add_on"
	AccessRestricted CartridgeAccessLevel = "restricted"
	AccessBeta       CartridgeAccessLevel = "beta"
)

// CartridgeAccessResult는 카트리지 접근 검증 결과입니다.
type CartridgeAccessResult struct {
	Allowed          bool
	AccessLevel      CartridgeAccessLevel
	RemainingDaily   int32 // -1 = 무제한
	RemainingMonthly int32 // -1 = 무제한
	RequiredTier     SubscriptionTier
	CurrentTier      SubscriptionTier
	Message          string
	AddonPriceKRW    int32
}

// 기본 카트리지 접근 정책 (인메모리, DB 연동 시 교체 예정)
// key: "tier:category:type" (type=0이면 카테고리 전체, category=0이면 글로벌 기본값)
type cartridgeAccessPolicy struct {
	AccessLevel CartridgeAccessLevel
	DailyLimit  int32 // 0 = 무제한
	Priority    int32
}

var defaultCartridgePolicy = map[string]*cartridgeAccessPolicy{
	// Free: 글로벌 기본값 = RESTRICTED
	"0:0:0": {AccessLevel: AccessRestricted, Priority: 0},
	// Free: HealthBiomarker Glucose/LipidPanel/HbA1c = LIMITED (일 3회)
	"0:1:1": {AccessLevel: AccessLimited, DailyLimit: 3, Priority: 10},
	"0:1:2": {AccessLevel: AccessLimited, DailyLimit: 3, Priority: 10},
	"0:1:3": {AccessLevel: AccessLimited, DailyLimit: 3, Priority: 10},

	// Basic: 글로벌 기본값 = RESTRICTED
	"1:0:0": {AccessLevel: AccessRestricted, Priority: 0},
	// Basic: HealthBiomarker 전체 = INCLUDED
	"1:1:0": {AccessLevel: AccessIncluded, Priority: 5},
	// Basic: Environmental/FoodSafety = ADD_ON
	"1:2:0": {AccessLevel: AccessAddOn, Priority: 5},
	"1:3:0": {AccessLevel: AccessAddOn, Priority: 5},

	// Pro: 글로벌 기본값 = RESTRICTED
	"2:0:0": {AccessLevel: AccessRestricted, Priority: 0},
	// Pro: Health/Env/Food/Sensor/Cosmetic = INCLUDED
	"2:1:0":  {AccessLevel: AccessIncluded, Priority: 5},
	"2:2:0":  {AccessLevel: AccessIncluded, Priority: 5},
	"2:3:0":  {AccessLevel: AccessIncluded, Priority: 5},
	"2:4:0":  {AccessLevel: AccessIncluded, Priority: 5},
	"2:10:0": {AccessLevel: AccessIncluded, Priority: 5},
	// Pro: Advanced/Veterinary/Agricultural/Marine = ADD_ON
	"2:5:0":  {AccessLevel: AccessAddOn, Priority: 5},
	"2:7:0":  {AccessLevel: AccessAddOn, Priority: 5},
	"2:9:0":  {AccessLevel: AccessAddOn, Priority: 5},
	"2:12:0": {AccessLevel: AccessAddOn, Priority: 5},

	// Clinical: 글로벌 기본값 = INCLUDED (전체 무제한)
	"3:0:0": {AccessLevel: AccessIncluded, Priority: 0},
	// Clinical: Beta = BETA
	"3:254:0": {AccessLevel: AccessBeta, Priority: 10},
}

// CheckCartridgeAccess는 사용자의 구독 등급에 따라 카트리지 접근 가능 여부를 검증합니다.
// 정책 적용 우선순위: 타입별 오버라이드 > 카테고리별 > 글로벌 기본값
func (s *SubscriptionService) CheckCartridgeAccess(ctx context.Context, userID string, categoryCode, typeIndex int32) *CartridgeAccessResult {
	sub, err := s.subRepo.GetByUserID(ctx, userID)
	if err != nil || sub == nil {
		return &CartridgeAccessResult{
			Allowed:     false,
			AccessLevel: AccessRestricted,
			CurrentTier: TierFree,
			Message:     "구독 정보가 없습니다. 가입 후 이용해 주세요.",
		}
	}

	currentTier := sub.Tier
	tierStr := fmt.Sprintf("%d", int32(currentTier))

	// 1단계: 타입별 오버라이드 확인
	typeKey := fmt.Sprintf("%s:%d:%d", tierStr, categoryCode, typeIndex)
	if policy, ok := defaultCartridgePolicy[typeKey]; ok {
		return s.buildAccessResult(policy, currentTier)
	}

	// 2단계: 카테고리별 정책 확인
	catKey := fmt.Sprintf("%s:%d:0", tierStr, categoryCode)
	if policy, ok := defaultCartridgePolicy[catKey]; ok {
		return s.buildAccessResult(policy, currentTier)
	}

	// 3단계: 글로벌 기본값
	globalKey := fmt.Sprintf("%s:0:0", tierStr)
	if policy, ok := defaultCartridgePolicy[globalKey]; ok {
		return s.buildAccessResult(policy, currentTier)
	}

	// fallback: RESTRICTED
	return &CartridgeAccessResult{
		Allowed:     false,
		AccessLevel: AccessRestricted,
		CurrentTier: currentTier,
		Message:     "해당 카트리지는 현재 등급에서 사용할 수 없습니다.",
	}
}

// CartridgeAccessEntry는 접근 가능 카트리지 항목입니다.
type CartridgeAccessEntry struct {
	CategoryCode     int32
	TypeIndex        int32
	Name             string
	AccessLevel      CartridgeAccessLevel
	RemainingDaily   int32
	RemainingMonthly int32
}

// ListAccessibleCartridges는 사용자가 접근 가능한 카트리지 목록을 반환합니다.
func (s *SubscriptionService) ListAccessibleCartridges(ctx context.Context, userID string) []*CartridgeAccessEntry {
	sub, err := s.subRepo.GetByUserID(ctx, userID)
	tier := TierFree
	if err == nil && sub != nil {
		tier = sub.Tier
	}

	// 카트리지 타입 목록 (주요 14종 + 환경 4종 + 식품 4종)
	cartridgeTypes := []struct {
		Cat  int32
		Type int32
		Name string
	}{
		{1, 1, "혈당 (Glucose)"}, {1, 2, "지질 패널 (Lipid Panel)"}, {1, 3, "당화혈색소 (HbA1c)"},
		{1, 4, "요산 (Uric Acid)"}, {1, 5, "크레아티닌 (Creatinine)"}, {1, 6, "빌리루빈 (Bilirubin)"},
		{1, 7, "ALT"}, {1, 8, "AST"}, {1, 9, "TSH"}, {1, 10, "CRP"},
		{1, 11, "비타민 D"}, {1, 12, "철분 (Ferritin)"}, {1, 13, "전해질"}, {1, 14, "BNP"},
		{2, 1, "미세먼지 (PM2.5)"}, {2, 2, "VOC"}, {2, 3, "이산화탄소 (CO2)"}, {2, 4, "라돈"},
		{3, 1, "잔류농약"}, {3, 2, "중금속"}, {3, 3, "세균"}, {3, 4, "항생제"},
	}

	entries := make([]*CartridgeAccessEntry, 0, len(cartridgeTypes))
	for _, ct := range cartridgeTypes {
		result := s.CheckCartridgeAccess(ctx, userID, ct.Cat, ct.Type)
		_ = tier // already used via CheckCartridgeAccess
		entries = append(entries, &CartridgeAccessEntry{
			CategoryCode:     ct.Cat,
			TypeIndex:        ct.Type,
			Name:             ct.Name,
			AccessLevel:      result.AccessLevel,
			RemainingDaily:   result.RemainingDaily,
			RemainingMonthly: result.RemainingMonthly,
		})
	}

	return entries
}

func (s *SubscriptionService) buildAccessResult(policy *cartridgeAccessPolicy, currentTier SubscriptionTier) *CartridgeAccessResult {
	result := &CartridgeAccessResult{
		AccessLevel: policy.AccessLevel,
		CurrentTier: currentTier,
	}

	switch policy.AccessLevel {
	case AccessIncluded:
		result.Allowed = true
		result.RemainingDaily = -1
		result.RemainingMonthly = -1
	case AccessLimited:
		result.Allowed = true
		result.RemainingDaily = policy.DailyLimit
		result.RemainingMonthly = -1
		result.Message = fmt.Sprintf("일일 %d회 제한", policy.DailyLimit)
	case AccessAddOn:
		result.Allowed = false // 구매 여부 별도 확인 필요
		result.Message = "별도 구매가 필요한 카트리지입니다."
	case AccessRestricted:
		result.Allowed = false
		result.Message = "상위 구독 등급이 필요합니다."
	case AccessBeta:
		if currentTier >= TierClinical {
			result.Allowed = true
			result.Message = "베타 카트리지 (Clinical 전용)"
		} else {
			result.Allowed = false
			result.Message = "Clinical Guard 등급 전용 베타 카트리지입니다."
			result.RequiredTier = TierClinical
		}
	}

	return result
}
