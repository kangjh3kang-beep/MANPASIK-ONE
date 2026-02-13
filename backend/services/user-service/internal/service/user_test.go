package service

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"
)

// =============================================================================
// 목(Mock) 저장소
// =============================================================================

type mockProfileRepo struct {
	profiles map[string]*UserProfile
}

func newMockProfileRepo() *mockProfileRepo {
	return &mockProfileRepo{profiles: make(map[string]*UserProfile)}
}

func (r *mockProfileRepo) GetByID(_ context.Context, userID string) (*UserProfile, error) {
	p, ok := r.profiles[userID]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (r *mockProfileRepo) Update(_ context.Context, profile *UserProfile) error {
	r.profiles[profile.UserID] = profile
	return nil
}

type mockSubRepo struct {
	subs map[string]*Subscription
}

func newMockSubRepo() *mockSubRepo {
	return &mockSubRepo{subs: make(map[string]*Subscription)}
}

func (r *mockSubRepo) GetByUserID(_ context.Context, userID string) (*Subscription, error) {
	s, ok := r.subs[userID]
	if !ok {
		return nil, nil
	}
	return s, nil
}

func (r *mockSubRepo) Create(_ context.Context, sub *Subscription) error {
	r.subs[sub.UserID] = sub
	return nil
}

func (r *mockSubRepo) Update(_ context.Context, sub *Subscription) error {
	r.subs[sub.UserID] = sub
	return nil
}

type mockFamilyRepo struct{}

func newMockFamilyRepo() *mockFamilyRepo { return &mockFamilyRepo{} }

func (r *mockFamilyRepo) GetGroup(_ context.Context, _ string) (*FamilyGroup, error) {
	return nil, nil
}
func (r *mockFamilyRepo) GetUserGroups(_ context.Context, _ string) ([]*FamilyGroup, error) {
	return nil, nil
}
func (r *mockFamilyRepo) CreateGroup(_ context.Context, _ *FamilyGroup) error { return nil }
func (r *mockFamilyRepo) AddMember(_ context.Context, _, _, _ string) error   { return nil }
func (r *mockFamilyRepo) RemoveMember(_ context.Context, _, _ string) error   { return nil }

// =============================================================================
// 헬퍼
// =============================================================================

func newTestUserService() (*UserService, *mockProfileRepo, *mockSubRepo) {
	logger, _ := zap.NewDevelopment()
	profileRepo := newMockProfileRepo()
	subRepo := newMockSubRepo()
	familyRepo := newMockFamilyRepo()
	svc := NewUserService(logger, profileRepo, subRepo, familyRepo)
	return svc, profileRepo, subRepo
}

// =============================================================================
// GetProfile 테스트
// =============================================================================

func TestGetProfile_성공(t *testing.T) {
	svc, profileRepo, _ := newTestUserService()
	ctx := context.Background()

	// 시드 데이터
	profileRepo.profiles["user-1"] = &UserProfile{
		UserID:      "user-1",
		Email:       "test@manpasik.com",
		DisplayName: "테스트 사용자",
		Language:    "ko",
		Timezone:    "Asia/Seoul",
		CreatedAt:   time.Now().UTC(),
	}

	profile, err := svc.GetProfile(ctx, "user-1")
	if err != nil {
		t.Fatalf("프로필 조회 실패: %v", err)
	}
	if profile.Email != "test@manpasik.com" {
		t.Errorf("이메일 불일치: got %s, want test@manpasik.com", profile.Email)
	}
	if profile.DisplayName != "테스트 사용자" {
		t.Errorf("이름 불일치: got %s, want 테스트 사용자", profile.DisplayName)
	}
}

func TestGetProfile_존재하지_않는_사용자(t *testing.T) {
	svc, _, _ := newTestUserService()
	ctx := context.Background()

	_, err := svc.GetProfile(ctx, "nonexistent")
	if err == nil {
		t.Fatal("존재하지 않는 사용자에 대해 에러가 발생해야 합니다")
	}
}

func TestGetProfile_빈_유저ID(t *testing.T) {
	svc, _, _ := newTestUserService()
	ctx := context.Background()

	_, err := svc.GetProfile(ctx, "")
	if err == nil {
		t.Fatal("빈 user_id에 대해 에러가 발생해야 합니다")
	}
}

// =============================================================================
// UpdateProfile 테스트
// =============================================================================

func TestUpdateProfile_성공(t *testing.T) {
	svc, profileRepo, _ := newTestUserService()
	ctx := context.Background()

	profileRepo.profiles["user-1"] = &UserProfile{
		UserID:      "user-1",
		Email:       "test@manpasik.com",
		DisplayName: "원래 이름",
		Language:    "ko",
	}

	updated, err := svc.UpdateProfile(ctx, "user-1", "새 이름", "", "en", "")
	if err != nil {
		t.Fatalf("프로필 업데이트 실패: %v", err)
	}
	if updated.DisplayName != "새 이름" {
		t.Errorf("이름 불일치: got %s, want 새 이름", updated.DisplayName)
	}
	if updated.Language != "en" {
		t.Errorf("언어 불일치: got %s, want en", updated.Language)
	}
}

func TestUpdateProfile_잘못된_언어(t *testing.T) {
	svc, profileRepo, _ := newTestUserService()
	ctx := context.Background()

	profileRepo.profiles["user-1"] = &UserProfile{
		UserID: "user-1",
		Email:  "test@manpasik.com",
	}

	_, err := svc.UpdateProfile(ctx, "user-1", "", "", "fr", "")
	if err == nil {
		t.Fatal("지원하지 않는 언어에 대해 에러가 발생해야 합니다")
	}
}

func TestUpdateProfile_존재하지_않는_사용자(t *testing.T) {
	svc, _, _ := newTestUserService()
	ctx := context.Background()

	_, err := svc.UpdateProfile(ctx, "nonexistent", "이름", "", "", "")
	if err == nil {
		t.Fatal("존재하지 않는 사용자에 대해 에러가 발생해야 합니다")
	}
}

// =============================================================================
// GetSubscription 테스트
// =============================================================================

func TestGetSubscription_존재하는_구독(t *testing.T) {
	svc, _, subRepo := newTestUserService()
	ctx := context.Background()

	subRepo.subs["user-1"] = &Subscription{
		ID:                "sub-1",
		UserID:            "user-1",
		Tier:              TierPro,
		MaxDevices:        10,
		MaxFamilyMembers:  5,
		AICoachingEnabled: true,
	}

	sub, err := svc.GetSubscription(ctx, "user-1")
	if err != nil {
		t.Fatalf("구독 조회 실패: %v", err)
	}
	if sub.Tier != TierPro {
		t.Errorf("티어 불일치: got %s, want pro", sub.Tier)
	}
	if sub.MaxDevices != 10 {
		t.Errorf("최대 디바이스 불일치: got %d, want 10", sub.MaxDevices)
	}
}

func TestGetSubscription_기본_Free_티어(t *testing.T) {
	svc, _, _ := newTestUserService()
	ctx := context.Background()

	sub, err := svc.GetSubscription(ctx, "user-no-sub")
	if err != nil {
		t.Fatalf("구독 조회 실패: %v", err)
	}
	if sub.Tier != TierFree {
		t.Errorf("기본 티어가 Free여야 합니다: got %s", sub.Tier)
	}
	if sub.MaxDevices != 1 {
		t.Errorf("Free 티어 최대 디바이스는 1이어야 합니다: got %d", sub.MaxDevices)
	}
}

func TestGetMaxDevices_성공(t *testing.T) {
	svc, _, subRepo := newTestUserService()
	ctx := context.Background()

	subRepo.subs["user-1"] = &Subscription{
		UserID:     "user-1",
		Tier:       TierClinical,
		MaxDevices: 999,
	}

	max, err := svc.GetMaxDevices(ctx, "user-1")
	if err != nil {
		t.Fatalf("GetMaxDevices 실패: %v", err)
	}
	if max != 999 {
		t.Errorf("Clinical 티어 최대 디바이스는 999여야 합니다: got %d", max)
	}
}

func TestGetSubscription_빈_유저ID(t *testing.T) {
	svc, _, _ := newTestUserService()
	ctx := context.Background()

	_, err := svc.GetSubscription(ctx, "")
	if err == nil {
		t.Fatal("빈 user_id에 대해 에러가 발생해야 합니다")
	}
}
