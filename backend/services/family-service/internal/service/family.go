// Package service는 family-service의 비즈니스 로직을 구현합니다.
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/manpasik/backend/shared/errors"
	"go.uber.org/zap"
)

// FamilyRole은 가족 멤버 역할입니다.
type FamilyRole int

const (
	RoleUnknown  FamilyRole = iota
	RoleOwner               // 그룹 생성자 (관리자)
	RoleGuardian            // 보호자 (전체 열람)
	RoleMember              // 일반 멤버
	RoleChild               // 자녀 (제한적)
	RoleElderly             // 어르신 (보호 대상)
)

// String 은 FamilyRole 의 사람 친화 표현 (tenancy 동기화 등에 사용).
func (r FamilyRole) String() string {
	switch r {
	case RoleOwner:
		return "owner"
	case RoleGuardian:
		return "guardian"
	case RoleMember:
		return "member"
	case RoleChild:
		return "child"
	case RoleElderly:
		return "elderly"
	default:
		return "unknown"
	}
}

// InvitationStatus는 초대 상태입니다.
type InvitationStatus int

const (
	InviteUnknown  InvitationStatus = iota
	InvitePending                   // 대기 중
	InviteAccepted                  // 수락됨
	InviteDeclined                  // 거절됨
	InviteExpired                   // 만료됨
)

// FamilyGroup은 가족 그룹 도메인 객체입니다.
type FamilyGroup struct {
	ID          string
	OwnerUserID string
	GroupName   string
	Description string
	MaxMembers  int
	CreatedAt   time.Time
}

// FamilyMember는 가족 멤버 도메인 객체입니다.
type FamilyMember struct {
	UserID         string
	GroupID        string
	DisplayName    string
	Email          string
	Role           FamilyRole
	JoinedAt       time.Time
	SharingEnabled bool
}

// FamilyInvitation은 초대 도메인 객체입니다.
type FamilyInvitation struct {
	ID            string
	GroupID       string
	GroupName     string
	InviterUserID string
	InviteeEmail  string
	Role          FamilyRole
	Message       string
	Status        InvitationStatus
	CreatedAt     time.Time
	ExpiresAt     time.Time
}

// SharingPreferences는 건강 데이터 공유 설정입니다.
type SharingPreferences struct {
	UserID               string
	GroupID              string
	ShareMeasurements    bool
	ShareHealthScore     bool
	ShareGoals           bool
	ShareCoaching        bool
	ShareAlerts          bool
	AllowedViewerIDs     []string
	MeasurementDaysLimit int      // 최근 N일만 공유 (0 = 무제한)
	AllowedBiomarkers    []string // 비어있으면 전체, 아니면 특정 바이오마커만
	RequireApproval      bool     // true이면 공유 요청마다 승인 필요
}

// SharedHealthSummary는 공유된 건강 데이터 요약입니다.
type SharedHealthSummary struct {
	UserID              string
	DisplayName         string
	HealthScore         float64
	MeasurementsCount   int
	ScoreTrend          string
	LatestAlert         string
	LastMeasurementAt   *time.Time
}

// FamilyGroupRepository는 가족 그룹 저장소 인터페이스입니다.
type FamilyGroupRepository interface {
	Save(ctx context.Context, g *FamilyGroup) error
	FindByID(ctx context.Context, id string) (*FamilyGroup, error)
}

// FamilyMemberRepository는 가족 멤버 저장소 인터페이스입니다.
type FamilyMemberRepository interface {
	Save(ctx context.Context, m *FamilyMember) error
	FindByGroupID(ctx context.Context, groupID string) ([]*FamilyMember, error)
	FindByUserIDAndGroupID(ctx context.Context, userID, groupID string) (*FamilyMember, error)
	Remove(ctx context.Context, userID, groupID string) error
	CountByGroupID(ctx context.Context, groupID string) (int, error)
}

// InvitationRepository는 초대 저장소 인터페이스입니다.
type InvitationRepository interface {
	Save(ctx context.Context, inv *FamilyInvitation) error
	FindByID(ctx context.Context, id string) (*FamilyInvitation, error)
	Update(ctx context.Context, inv *FamilyInvitation) error
}

// SharingPreferencesRepository는 공유 설정 저장소 인터페이스입니다.
type SharingPreferencesRepository interface {
	Save(ctx context.Context, pref *SharingPreferences) error
	FindByUserIDAndGroupID(ctx context.Context, userID, groupID string) (*SharingPreferences, error)
}

// TenancySync 는 family-service 의 그룹/멤버 변경을 tenancy 영속 store 에
// 자동 반영하는 어댑터.
//
// family-service 는 도메인별 (가족 공유 데이터, 알림 등) 로직을 담당하고,
// tenancy 는 권한/멤버십 격리를 담당. 두 개념이 분리되어 있으므로 동기화
// 어댑터로 교차 도메인 일관성 유지.
type TenancySync interface {
	// OnGroupCreated: 그룹 생성 시 owner 를 tenancy.MembershipOwner 로 등록.
	OnGroupCreated(ctx context.Context, groupID, ownerUserID string) error
	// OnMemberAdded: 멤버 가입 시 tenancy.Member 등록 (역할은 family role 에서 변환).
	OnMemberAdded(ctx context.Context, groupID, userID, familyRole string) error
	// OnMemberRemoved: 멤버 제거 시 tenancy 멤버십도 제거.
	OnMemberRemoved(ctx context.Context, groupID, userID string) error
	// OnMemberRoleChanged: 역할 변경 시 tenancy 측 역할도 동기화.
	OnMemberRoleChanged(ctx context.Context, groupID, userID, newFamilyRole string) error
}

// FamilyService는 가족 서비스 핵심 로직입니다.
type FamilyService struct {
	log         *zap.Logger
	groupRepo   FamilyGroupRepository
	memberRepo  FamilyMemberRepository
	inviteRepo  InvitationRepository
	sharingRepo SharingPreferencesRepository
	tenancySync TenancySync // 옵션 — 미설정 시 동기화 안 함
}

// NewFamilyService는 FamilyService를 생성합니다.
func NewFamilyService(log *zap.Logger, groupRepo FamilyGroupRepository, memberRepo FamilyMemberRepository, inviteRepo InvitationRepository, sharingRepo SharingPreferencesRepository) *FamilyService {
	return &FamilyService{
		log:         log,
		groupRepo:   groupRepo,
		memberRepo:  memberRepo,
		inviteRepo:  inviteRepo,
		sharingRepo: sharingRepo,
	}
}

// SetTenancySync 는 tenancy 동기화 어댑터 등록.
//
// 운영 환경에서는 main.go 에서 호출. 미설정 시 family-service 는 기존대로
// 동작 (tenancy 미연동).
func (s *FamilyService) SetTenancySync(sync TenancySync) {
	s.tenancySync = sync
}

// CreateFamilyGroup은 가족 그룹을 생성합니다.
func (s *FamilyService) CreateFamilyGroup(ctx context.Context, ownerUserID, groupName, description string) (*FamilyGroup, error) {
	if ownerUserID == "" || groupName == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}

	group := &FamilyGroup{
		ID:          uuid.New().String(),
		OwnerUserID: ownerUserID,
		GroupName:   groupName,
		Description: description,
		MaxMembers:  10, // 기본값, 구독 티어에 따라 조정
		CreatedAt:   time.Now(),
	}

	if err := s.groupRepo.Save(ctx, group); err != nil {
		return nil, fmt.Errorf("가족 그룹 저장 실패: %w", err)
	}

	// Owner를 멤버로 자동 추가
	owner := &FamilyMember{
		UserID:         ownerUserID,
		GroupID:        group.ID,
		DisplayName:    "그룹장",
		Role:           RoleOwner,
		JoinedAt:       time.Now(),
		SharingEnabled: true,
	}
	if err := s.memberRepo.Save(ctx, owner); err != nil {
		return nil, fmt.Errorf("Owner 멤버 등록 실패: %w", err)
	}

	// Tenancy 동기화 (옵션): 그룹 ID 를 tenant ID 로, owner 를 owner 역할로 등록.
	// 실패는 family-service 핵심 로직과 격리 — 로그만 남기고 진행.
	if s.tenancySync != nil {
		if err := s.tenancySync.OnGroupCreated(ctx, group.ID, ownerUserID); err != nil {
			s.log.Warn("tenancy 동기화 실패 (그룹 생성)",
				zap.String("group_id", group.ID),
				zap.Error(err))
		}
	}

	s.log.Info("가족 그룹 생성",
		zap.String("group_id", group.ID),
		zap.String("owner", ownerUserID),
	)

	return group, nil
}

// GetFamilyGroup은 가족 그룹을 조회합니다.
func (s *FamilyService) GetFamilyGroup(ctx context.Context, groupID string) (*FamilyGroup, int, error) {
	if groupID == "" {
		return nil, 0, apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}

	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.memberRepo.CountByGroupID(ctx, groupID)
	if err != nil {
		return nil, 0, err
	}

	return group, count, nil
}

// InviteMember는 가족 멤버를 초대합니다.
func (s *FamilyService) InviteMember(ctx context.Context, groupID, inviterUserID, inviteeEmail string, role FamilyRole, message string) (*FamilyInvitation, error) {
	if groupID == "" || inviterUserID == "" || inviteeEmail == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}

	// 그룹 존재 확인
	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	// 초대자 권한 확인 (Owner/Guardian만 초대 가능)
	inviter, err := s.memberRepo.FindByUserIDAndGroupID(ctx, inviterUserID, groupID)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrUnauthorized, "권한이 없습니다")
	}
	if inviter.Role != RoleOwner && inviter.Role != RoleGuardian {
		return nil, apperrors.New(apperrors.ErrUnauthorized, "권한이 없습니다")
	}

	// 최대 멤버 수 확인
	count, err := s.memberRepo.CountByGroupID(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if count >= group.MaxMembers {
		return nil, fmt.Errorf("최대 멤버 수(%d)에 도달했습니다", group.MaxMembers)
	}

	if role == RoleUnknown {
		role = RoleMember
	}

	invitation := &FamilyInvitation{
		ID:            uuid.New().String(),
		GroupID:       groupID,
		GroupName:     group.GroupName,
		InviterUserID: inviterUserID,
		InviteeEmail:  inviteeEmail,
		Role:          role,
		Message:       message,
		Status:        InvitePending,
		CreatedAt:     time.Now(),
		ExpiresAt:     time.Now().Add(7 * 24 * time.Hour), // 7일 유효
	}

	if err := s.inviteRepo.Save(ctx, invitation); err != nil {
		return nil, fmt.Errorf("초대 저장 실패: %w", err)
	}

	s.log.Info("가족 멤버 초대",
		zap.String("invitation_id", invitation.ID),
		zap.String("group_id", groupID),
		zap.String("invitee", inviteeEmail),
	)

	return invitation, nil
}

// RespondToInvitation은 초대를 수락하거나 거절합니다.
func (s *FamilyService) RespondToInvitation(ctx context.Context, invitationID, userID string, accept bool) (string, string, error) {
	if invitationID == "" || userID == "" {
		return "", "", apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}

	invitation, err := s.inviteRepo.FindByID(ctx, invitationID)
	if err != nil {
		return "", "", err
	}

	if invitation.Status != InvitePending {
		return "", "", fmt.Errorf("이미 처리된 초대입니다")
	}

	if time.Now().After(invitation.ExpiresAt) {
		invitation.Status = InviteExpired
		s.inviteRepo.Update(ctx, invitation)
		return "", "", fmt.Errorf("만료된 초대입니다")
	}

	if accept {
		invitation.Status = InviteAccepted
		s.inviteRepo.Update(ctx, invitation)

		member := &FamilyMember{
			UserID:         userID,
			GroupID:        invitation.GroupID,
			DisplayName:    invitation.InviteeEmail,
			Email:          invitation.InviteeEmail,
			Role:           invitation.Role,
			JoinedAt:       time.Now(),
			SharingEnabled: false,
		}
		if err := s.memberRepo.Save(ctx, member); err != nil {
			return "", "", fmt.Errorf("멤버 등록 실패: %w", err)
		}

		// Tenancy 동기화 (옵션): family role → tenancy role 변환 후 등록.
		if s.tenancySync != nil {
			if err := s.tenancySync.OnMemberAdded(ctx, invitation.GroupID, userID, invitation.Role.String()); err != nil {
				s.log.Warn("tenancy 동기화 실패 (멤버 추가)",
					zap.String("group_id", invitation.GroupID),
					zap.String("user_id", userID),
					zap.Error(err))
			}
		}

		return invitation.GroupID, "초대를 수락했습니다", nil
	}

	invitation.Status = InviteDeclined
	s.inviteRepo.Update(ctx, invitation)
	return "", "초대를 거절했습니다", nil
}

// RemoveMember는 가족 멤버를 제거합니다.
func (s *FamilyService) RemoveMember(ctx context.Context, groupID, userID string) error {
	if groupID == "" || userID == "" {
		return apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}

	// 멤버 존재 확인
	member, err := s.memberRepo.FindByUserIDAndGroupID(ctx, userID, groupID)
	if err != nil {
		return apperrors.New(apperrors.ErrNotFound, "멤버를 찾을 수 없습니다")
	}

	// Owner는 자기 자신을 제거할 수 없음
	if member.Role == RoleOwner {
		return fmt.Errorf("그룹 Owner는 탈퇴할 수 없습니다. 그룹을 삭제하세요")
	}

	if err := s.memberRepo.Remove(ctx, userID, groupID); err != nil {
		return err
	}

	// Tenancy 동기화 (옵션)
	if s.tenancySync != nil {
		if err := s.tenancySync.OnMemberRemoved(ctx, groupID, userID); err != nil {
			s.log.Warn("tenancy 동기화 실패 (멤버 제거)",
				zap.String("group_id", groupID),
				zap.String("user_id", userID),
				zap.Error(err))
		}
	}
	return nil
}

// UpdateMemberRole은 멤버 역할을 변경합니다.
func (s *FamilyService) UpdateMemberRole(ctx context.Context, groupID, userID string, newRole FamilyRole) (*FamilyMember, error) {
	if groupID == "" || userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}

	member, err := s.memberRepo.FindByUserIDAndGroupID(ctx, userID, groupID)
	if err != nil {
		return nil, err
	}

	member.Role = newRole
	if err := s.memberRepo.Save(ctx, member); err != nil {
		return nil, err
	}

	if s.tenancySync != nil {
		if err := s.tenancySync.OnMemberRoleChanged(ctx, groupID, userID, newRole.String()); err != nil {
			s.log.Warn("tenancy 동기화 실패 (역할 변경)",
				zap.String("group_id", groupID),
				zap.String("user_id", userID),
				zap.Error(err))
		}
	}

	return member, nil
}

// ListFamilyMembers는 가족 멤버 목록을 조회합니다.
func (s *FamilyService) ListFamilyMembers(ctx context.Context, groupID string) ([]*FamilyMember, error) {
	if groupID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}
	return s.memberRepo.FindByGroupID(ctx, groupID)
}

// SetSharingPreferences는 건강 데이터 공유 설정을 변경합니다.
func (s *FamilyService) SetSharingPreferences(ctx context.Context, pref *SharingPreferences) (*SharingPreferences, error) {
	if pref == nil || pref.UserID == "" || pref.GroupID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}

	// 멤버 확인
	_, err := s.memberRepo.FindByUserIDAndGroupID(ctx, pref.UserID, pref.GroupID)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrUnauthorized, "권한이 없습니다")
	}

	if err := s.sharingRepo.Save(ctx, pref); err != nil {
		return nil, fmt.Errorf("공유 설정 저장 실패: %w", err)
	}

	return pref, nil
}

// ValidateSharingAccess는 특정 바이오마커에 대한 공유 접근 권한을 검증합니다.
// allowed=true이면 접근 허용, reason은 거부 사유 또는 허용 사유입니다.
func (s *FamilyService) ValidateSharingAccess(ctx context.Context, groupID, requesterID, targetUserID, biomarker string) (bool, string, error) {
	if groupID == "" || requesterID == "" || targetUserID == "" {
		return false, "", apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}

	// 1. 요청자가 그룹 멤버인지 확인
	_, err := s.memberRepo.FindByUserIDAndGroupID(ctx, requesterID, groupID)
	if err != nil {
		return false, "요청자가 그룹 멤버가 아닙니다", nil
	}

	// 2. 대상 사용자가 그룹 멤버인지 확인
	_, err = s.memberRepo.FindByUserIDAndGroupID(ctx, targetUserID, groupID)
	if err != nil {
		return false, "대상 사용자가 그룹 멤버가 아닙니다", nil
	}

	// 3. 대상 사용자의 공유 설정 확인
	pref, _ := s.sharingRepo.FindByUserIDAndGroupID(ctx, targetUserID, groupID)
	if pref == nil {
		return false, "대상 사용자의 공유 설정이 없습니다", nil
	}

	// 공유가 활성화되어 있는지 확인
	if !pref.ShareMeasurements {
		return false, "대상 사용자가 측정 데이터 공유를 비활성화했습니다", nil
	}

	// 4. 승인 필요 여부 확인
	if pref.RequireApproval {
		return false, "대상 사용자가 공유 승인을 요구합니다", nil
	}

	// 5. 허용된 바이오마커 확인
	if biomarker != "" && len(pref.AllowedBiomarkers) > 0 {
		found := false
		for _, allowed := range pref.AllowedBiomarkers {
			if allowed == biomarker {
				found = true
				break
			}
		}
		if !found {
			return false, fmt.Sprintf("바이오마커 '%s'에 대한 접근이 허용되지 않았습니다", biomarker), nil
		}
	}

	// 6. MeasurementDaysLimit 정보 포함
	reason := "접근이 허용되었습니다"
	if pref.MeasurementDaysLimit > 0 {
		reason = fmt.Sprintf("접근이 허용되었습니다 (최근 %d일 데이터만)", pref.MeasurementDaysLimit)
	}

	return true, reason, nil
}

// GetSharedHealthData는 공유된 건강 데이터를 조회합니다.
func (s *FamilyService) GetSharedHealthData(ctx context.Context, requesterUserID, targetUserID, groupID string, days int) ([]*SharedHealthSummary, error) {
	if requesterUserID == "" || groupID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}

	// 요청자가 그룹 멤버인지 확인
	requester, err := s.memberRepo.FindByUserIDAndGroupID(ctx, requesterUserID, groupID)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrUnauthorized, "권한이 없습니다")
	}

	if days <= 0 {
		days = 7
	}

	// targetUserID가 지정되면 해당 멤버의 데이터만, 아니면 전체 멤버
	members, err := s.memberRepo.FindByGroupID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	var summaries []*SharedHealthSummary
	for _, m := range members {
		if targetUserID != "" && m.UserID != targetUserID {
			continue
		}
		if m.UserID == requesterUserID {
			continue // 본인 데이터는 제외
		}

		// 공유 설정 확인
		pref, _ := s.sharingRepo.FindByUserIDAndGroupID(ctx, m.UserID, groupID)

		// Guardian/Owner는 공유 설정 무관하게 열람 가능
		canView := false
		if requester.Role == RoleOwner || requester.Role == RoleGuardian {
			canView = true
		} else if pref != nil && pref.ShareHealthScore {
			canView = true
		}

		if canView {
			summary := &SharedHealthSummary{
				UserID:            m.UserID,
				DisplayName:       m.DisplayName,
				HealthScore:       75.0 + float64(len(m.UserID)%25), // 시뮬레이션 데이터
				MeasurementsCount: days * 2,
				ScoreTrend:        "stable",
			}
			summaries = append(summaries, summary)
		}
	}

	return summaries, nil
}
