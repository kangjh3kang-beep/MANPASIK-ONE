// Package service는 user-service의 비즈니스 로직을 구현합니다.
//
// 기능: 사용자 프로필 CRUD, 구독 관리, 가족 그룹
package service

import (
	"context"
	"time"

	apperrors "github.com/manpasik/backend/shared/errors"
	"go.uber.org/zap"
)

// UserService는 사용자 관리 서비스입니다.
type UserService struct {
	logger      *zap.Logger
	profileRepo ProfileRepository
	subRepo     SubscriptionRepository
	familyRepo  FamilyRepository
}

// ProfileRepository는 사용자 프로필 저장소 인터페이스입니다.
type ProfileRepository interface {
	GetByID(ctx context.Context, userID string) (*UserProfile, error)
	Update(ctx context.Context, profile *UserProfile) error
}

// SubscriptionRepository는 구독 저장소 인터페이스입니다.
type SubscriptionRepository interface {
	GetByUserID(ctx context.Context, userID string) (*Subscription, error)
	Create(ctx context.Context, sub *Subscription) error
	Update(ctx context.Context, sub *Subscription) error
}

// FamilyRepository는 가족 그룹 저장소 인터페이스입니다.
type FamilyRepository interface {
	GetGroup(ctx context.Context, groupID string) (*FamilyGroup, error)
	GetUserGroups(ctx context.Context, userID string) ([]*FamilyGroup, error)
	CreateGroup(ctx context.Context, group *FamilyGroup) error
	AddMember(ctx context.Context, groupID, userID, role string) error
	RemoveMember(ctx context.Context, groupID, userID string) error
}

// UserProfile은 사용자 프로필 엔티티입니다.
type UserProfile struct {
	UserID           string
	Email            string
	DisplayName      string
	AvatarURL        string
	Language         string // "ko", "en", "zh", "ja"
	Timezone         string
	SubscriptionTier string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// Subscription은 구독 엔티티입니다.
type Subscription struct {
	ID                  string
	UserID              string
	Tier                SubscriptionTier
	StartedAt           time.Time
	ExpiresAt           time.Time
	MaxDevices          int
	MaxFamilyMembers    int
	AICoachingEnabled   bool
	TelemedicineEnabled bool
}

// SubscriptionTier는 구독 티어입니다.
type SubscriptionTier string

const (
	TierFree     SubscriptionTier = "free"
	TierBasic    SubscriptionTier = "basic"    // ₩9,900/월
	TierPro      SubscriptionTier = "pro"      // ₩29,900/월
	TierClinical SubscriptionTier = "clinical" // ₩59,900/월
)

// TierConfig는 티어별 설정입니다.
var TierConfig = map[SubscriptionTier]struct {
	MaxDevices          int
	MaxFamilyMembers    int
	AICoachingEnabled   bool
	TelemedicineEnabled bool
}{
	TierFree:     {MaxDevices: 1, MaxFamilyMembers: 0, AICoachingEnabled: false, TelemedicineEnabled: false},
	TierBasic:    {MaxDevices: 3, MaxFamilyMembers: 2, AICoachingEnabled: false, TelemedicineEnabled: false},
	TierPro:      {MaxDevices: 10, MaxFamilyMembers: 5, AICoachingEnabled: true, TelemedicineEnabled: false},
	TierClinical: {MaxDevices: 999, MaxFamilyMembers: 10, AICoachingEnabled: true, TelemedicineEnabled: true},
}

// FamilyGroup은 가족 그룹 엔티티입니다.
type FamilyGroup struct {
	ID      string
	Name    string
	OwnerID string
	Members []*FamilyMember
}

// FamilyMember는 가족 구성원입니다.
type FamilyMember struct {
	UserID string
	Role   string // "owner", "adult", "child", "elderly"
}

// NewUserService는 새 UserService를 생성합니다.
func NewUserService(
	logger *zap.Logger,
	profileRepo ProfileRepository,
	subRepo SubscriptionRepository,
	familyRepo FamilyRepository,
) *UserService {
	return &UserService{
		logger:      logger,
		profileRepo: profileRepo,
		subRepo:     subRepo,
		familyRepo:  familyRepo,
	}
}

// GetProfile은 사용자 프로필을 조회합니다.
func (s *UserService) GetProfile(ctx context.Context, userID string) (*UserProfile, error) {
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}

	profile, err := s.profileRepo.GetByID(ctx, userID)
	if err != nil || profile == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "사용자를 찾을 수 없습니다")
	}

	return profile, nil
}

// UpdateProfile은 사용자 프로필을 업데이트합니다.
func (s *UserService) UpdateProfile(
	ctx context.Context,
	userID, displayName, avatarURL, language, timezone string,
) (*UserProfile, error) {
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}

	profile, err := s.profileRepo.GetByID(ctx, userID)
	if err != nil || profile == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "사용자를 찾을 수 없습니다")
	}

	// 변경 사항 적용
	if displayName != "" {
		profile.DisplayName = displayName
	}
	if avatarURL != "" {
		profile.AvatarURL = avatarURL
	}
	if language != "" {
		// 지원 언어 검증
		validLangs := map[string]bool{"ko": true, "en": true, "zh": true, "ja": true}
		if !validLangs[language] {
			return nil, apperrors.New(apperrors.ErrInvalidInput, "지원하지 않는 언어입니다 (ko, en, zh, ja)")
		}
		profile.Language = language
	}
	if timezone != "" {
		profile.Timezone = timezone
	}
	profile.UpdatedAt = time.Now().UTC()

	if err := s.profileRepo.Update(ctx, profile); err != nil {
		return nil, apperrors.New(apperrors.ErrInternal, "프로필 업데이트에 실패했습니다")
	}

	s.logger.Info("프로필 업데이트", zap.String("user_id", userID))
	return profile, nil
}

// GetSubscription은 구독 정보를 조회합니다.
func (s *UserService) GetSubscription(ctx context.Context, userID string) (*Subscription, error) {
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}

	sub, err := s.subRepo.GetByUserID(ctx, userID)
	if err != nil || sub == nil {
		// Free 티어 기본 반환
		config := TierConfig[TierFree]
		return &Subscription{
			UserID:              userID,
			Tier:                TierFree,
			MaxDevices:          config.MaxDevices,
			MaxFamilyMembers:    config.MaxFamilyMembers,
			AICoachingEnabled:   config.AICoachingEnabled,
			TelemedicineEnabled: config.TelemedicineEnabled,
		}, nil
	}

	return sub, nil
}

// GetMaxDevices는 사용자의 최대 디바이스 수를 반환합니다 (SubscriptionChecker 구현).
func (s *UserService) GetMaxDevices(ctx context.Context, userID string) (int, error) {
	sub, err := s.GetSubscription(ctx, userID)
	if err != nil {
		return 1, nil // 기본값
	}
	return sub.MaxDevices, nil
}
