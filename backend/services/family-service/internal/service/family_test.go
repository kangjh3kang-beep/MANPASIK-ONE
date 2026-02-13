package service_test

import (
	"context"
	"testing"

	"github.com/manpasik/backend/services/family-service/internal/repository/memory"
	"github.com/manpasik/backend/services/family-service/internal/service"
	"go.uber.org/zap"
)

func setupTestService() *service.FamilyService {
	logger := zap.NewNop()
	return service.NewFamilyService(
		logger,
		memory.NewFamilyGroupRepository(),
		memory.NewFamilyMemberRepository(),
		memory.NewInvitationRepository(),
		memory.NewSharingPreferencesRepository(),
	)
}

func TestCreateFamilyGroup_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	group, err := svc.CreateFamilyGroup(ctx, "owner-1", "김씨 가족", "우리 가족 그룹")
	if err != nil {
		t.Fatalf("가족 그룹 생성 실패: %v", err)
	}
	if group.ID == "" {
		t.Fatal("그룹 ID가 비어 있음")
	}
	if group.OwnerUserID != "owner-1" {
		t.Fatalf("Owner 불일치: got %s", group.OwnerUserID)
	}
	if group.GroupName != "김씨 가족" {
		t.Fatalf("그룹명 불일치: got %s", group.GroupName)
	}
}

func TestCreateFamilyGroup_EmptyOwner(t *testing.T) {
	svc := setupTestService()
	_, err := svc.CreateFamilyGroup(context.Background(), "", "가족", "")
	if err == nil {
		t.Fatal("빈 owner에 에러가 반환되어야 함")
	}
}

func TestGetFamilyGroup(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	group, _ := svc.CreateFamilyGroup(ctx, "owner-2", "이씨 가족", "")

	got, count, err := svc.GetFamilyGroup(ctx, group.ID)
	if err != nil {
		t.Fatalf("그룹 조회 실패: %v", err)
	}
	if got.GroupName != "이씨 가족" {
		t.Fatalf("그룹명 불일치: got %s", got.GroupName)
	}
	if count != 1 { // Owner 자동 포함
		t.Fatalf("멤버 수 불일치: got %d, want 1", count)
	}
}

func TestInviteAndAccept(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	group, _ := svc.CreateFamilyGroup(ctx, "owner-3", "박씨 가족", "")

	// 초대
	inv, err := svc.InviteMember(ctx, group.ID, "owner-3", "member@test.com", service.RoleMember, "가족 초대합니다")
	if err != nil {
		t.Fatalf("멤버 초대 실패: %v", err)
	}
	if inv.Status != service.InvitePending {
		t.Fatalf("초대 상태 불일치: got %d, want %d", inv.Status, service.InvitePending)
	}

	// 수락
	groupID, msg, err := svc.RespondToInvitation(ctx, inv.ID, "user-member", true)
	if err != nil {
		t.Fatalf("초대 수락 실패: %v", err)
	}
	if groupID != group.ID {
		t.Fatalf("수락 후 그룹 ID 불일치: got %s", groupID)
	}
	if msg == "" {
		t.Fatal("메시지가 비어 있음")
	}

	// 멤버 수 확인
	_, count, _ := svc.GetFamilyGroup(ctx, group.ID)
	if count != 2 {
		t.Fatalf("멤버 수 불일치: got %d, want 2", count)
	}
}

func TestInviteAndDecline(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	group, _ := svc.CreateFamilyGroup(ctx, "owner-4", "최씨 가족", "")
	inv, _ := svc.InviteMember(ctx, group.ID, "owner-4", "decline@test.com", service.RoleMember, "")

	_, msg, err := svc.RespondToInvitation(ctx, inv.ID, "user-decline", false)
	if err != nil {
		t.Fatalf("초대 거절 실패: %v", err)
	}
	if msg == "" {
		t.Fatal("거절 메시지가 비어 있음")
	}

	// 멤버 수 변동 없음 (Owner 1명)
	_, count, _ := svc.GetFamilyGroup(ctx, group.ID)
	if count != 1 {
		t.Fatalf("멤버 수 불일치: got %d, want 1", count)
	}
}

func TestInvite_UnauthorizedRole(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	group, _ := svc.CreateFamilyGroup(ctx, "owner-5", "정씨 가족", "")

	// 일반 멤버 추가
	inv, _ := svc.InviteMember(ctx, group.ID, "owner-5", "normal@test.com", service.RoleMember, "")
	svc.RespondToInvitation(ctx, inv.ID, "normal-user", true)

	// 일반 멤버가 초대 시도 → 권한 없음
	_, err := svc.InviteMember(ctx, group.ID, "normal-user", "another@test.com", service.RoleMember, "")
	if err == nil {
		t.Fatal("일반 멤버 초대에 에러가 반환되어야 함")
	}
}

func TestRemoveMember(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	group, _ := svc.CreateFamilyGroup(ctx, "owner-6", "한씨 가족", "")
	inv, _ := svc.InviteMember(ctx, group.ID, "owner-6", "remove@test.com", service.RoleMember, "")
	svc.RespondToInvitation(ctx, inv.ID, "user-remove", true)

	// Owner가 멤버 제거
	err := svc.RemoveMember(ctx, group.ID, "user-remove")
	if err != nil {
		t.Fatalf("멤버 제거 실패: %v", err)
	}

	_, count, _ := svc.GetFamilyGroup(ctx, group.ID)
	if count != 1 {
		t.Fatalf("멤버 수 불일치: got %d, want 1", count)
	}
}

func TestRemoveMember_OwnerCannotLeave(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	group, _ := svc.CreateFamilyGroup(ctx, "owner-7", "오씨 가족", "")

	err := svc.RemoveMember(ctx, group.ID, "owner-7")
	if err == nil {
		t.Fatal("Owner 탈퇴에 에러가 반환되어야 함")
	}
}

func TestUpdateMemberRole(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	group, _ := svc.CreateFamilyGroup(ctx, "owner-8", "장씨 가족", "")
	inv, _ := svc.InviteMember(ctx, group.ID, "owner-8", "role@test.com", service.RoleMember, "")
	svc.RespondToInvitation(ctx, inv.ID, "user-role", true)

	member, err := svc.UpdateMemberRole(ctx, group.ID, "user-role", service.RoleGuardian)
	if err != nil {
		t.Fatalf("역할 변경 실패: %v", err)
	}
	if member.Role != service.RoleGuardian {
		t.Fatalf("역할 불일치: got %d, want %d", member.Role, service.RoleGuardian)
	}
}

func TestListFamilyMembers(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	group, _ := svc.CreateFamilyGroup(ctx, "owner-9", "송씨 가족", "")
	inv, _ := svc.InviteMember(ctx, group.ID, "owner-9", "m1@test.com", service.RoleMember, "")
	svc.RespondToInvitation(ctx, inv.ID, "m1", true)
	inv2, _ := svc.InviteMember(ctx, group.ID, "owner-9", "m2@test.com", service.RoleChild, "")
	svc.RespondToInvitation(ctx, inv2.ID, "m2", true)

	members, err := svc.ListFamilyMembers(ctx, group.ID)
	if err != nil {
		t.Fatalf("멤버 목록 조회 실패: %v", err)
	}
	if len(members) != 3 { // Owner + 2명
		t.Fatalf("멤버 수 불일치: got %d, want 3", len(members))
	}
}

func TestSetSharingPreferences(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	group, _ := svc.CreateFamilyGroup(ctx, "owner-10", "유씨 가족", "")

	pref := &service.SharingPreferences{
		UserID:            "owner-10",
		GroupID:           group.ID,
		ShareMeasurements: true,
		ShareHealthScore:  true,
		ShareGoals:        false,
		ShareCoaching:     true,
		ShareAlerts:       true,
	}

	result, err := svc.SetSharingPreferences(ctx, pref)
	if err != nil {
		t.Fatalf("공유 설정 실패: %v", err)
	}
	if !result.ShareMeasurements {
		t.Fatal("ShareMeasurements가 true여야 함")
	}
}

func TestGetSharedHealthData_OwnerCanView(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	group, _ := svc.CreateFamilyGroup(ctx, "owner-11", "조씨 가족", "")
	inv, _ := svc.InviteMember(ctx, group.ID, "owner-11", "data@test.com", service.RoleMember, "")
	svc.RespondToInvitation(ctx, inv.ID, "data-user", true)

	summaries, err := svc.GetSharedHealthData(ctx, "owner-11", "", group.ID, 7)
	if err != nil {
		t.Fatalf("공유 데이터 조회 실패: %v", err)
	}
	if len(summaries) != 1 {
		t.Fatalf("요약 수 불일치: got %d, want 1 (본인 제외 1명)", len(summaries))
	}
	if summaries[0].UserID != "data-user" {
		t.Fatalf("사용자 ID 불일치: got %s", summaries[0].UserID)
	}
}

func TestValidateSharingAccess_Allowed(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	group, _ := svc.CreateFamilyGroup(ctx, "owner-va", "접근검증 가족", "")
	inv, _ := svc.InviteMember(ctx, group.ID, "owner-va", "target@test.com", service.RoleMember, "")
	svc.RespondToInvitation(ctx, inv.ID, "target-user", true)

	// 대상 사용자의 공유 설정
	svc.SetSharingPreferences(ctx, &service.SharingPreferences{
		UserID:            "target-user",
		GroupID:           group.ID,
		ShareMeasurements: true,
		ShareHealthScore:  true,
	})

	allowed, reason, err := svc.ValidateSharingAccess(ctx, group.ID, "owner-va", "target-user", "blood_glucose")
	if err != nil {
		t.Fatalf("ValidateSharingAccess 실패: %v", err)
	}
	if !allowed {
		t.Fatalf("접근이 허용되어야 합니다: reason=%s", reason)
	}
}

func TestValidateSharingAccess_NotMember(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	group, _ := svc.CreateFamilyGroup(ctx, "owner-nm", "멤버아님 가족", "")

	allowed, reason, err := svc.ValidateSharingAccess(ctx, group.ID, "nonmember", "owner-nm", "")
	if err != nil {
		t.Fatalf("에러 발생: %v", err)
	}
	if allowed {
		t.Fatal("비멤버 접근이 거부되어야 합니다")
	}
	if reason == "" {
		t.Fatal("거부 사유가 있어야 합니다")
	}
}

func TestValidateSharingAccess_SharingDisabled(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	group, _ := svc.CreateFamilyGroup(ctx, "owner-sd", "공유비활성 가족", "")
	inv, _ := svc.InviteMember(ctx, group.ID, "owner-sd", "sd@test.com", service.RoleMember, "")
	svc.RespondToInvitation(ctx, inv.ID, "user-sd", true)

	// 공유 비활성화 설정
	svc.SetSharingPreferences(ctx, &service.SharingPreferences{
		UserID:            "user-sd",
		GroupID:           group.ID,
		ShareMeasurements: false,
	})

	allowed, _, err := svc.ValidateSharingAccess(ctx, group.ID, "owner-sd", "user-sd", "")
	if err != nil {
		t.Fatalf("에러 발생: %v", err)
	}
	if allowed {
		t.Fatal("공유 비활성화 시 접근이 거부되어야 합니다")
	}
}

func TestValidateSharingAccess_RequireApproval(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	group, _ := svc.CreateFamilyGroup(ctx, "owner-ra", "승인필요 가족", "")
	inv, _ := svc.InviteMember(ctx, group.ID, "owner-ra", "ra@test.com", service.RoleMember, "")
	svc.RespondToInvitation(ctx, inv.ID, "user-ra", true)

	svc.SetSharingPreferences(ctx, &service.SharingPreferences{
		UserID:            "user-ra",
		GroupID:           group.ID,
		ShareMeasurements: true,
		RequireApproval:   true,
	})

	allowed, _, err := svc.ValidateSharingAccess(ctx, group.ID, "owner-ra", "user-ra", "")
	if err != nil {
		t.Fatalf("에러 발생: %v", err)
	}
	if allowed {
		t.Fatal("승인 필요 시 접근이 거부되어야 합니다")
	}
}

func TestValidateSharingAccess_BiomarkerRestriction(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	group, _ := svc.CreateFamilyGroup(ctx, "owner-br", "바이오마커 제한 가족", "")
	inv, _ := svc.InviteMember(ctx, group.ID, "owner-br", "br@test.com", service.RoleMember, "")
	svc.RespondToInvitation(ctx, inv.ID, "user-br", true)

	svc.SetSharingPreferences(ctx, &service.SharingPreferences{
		UserID:            "user-br",
		GroupID:           group.ID,
		ShareMeasurements: true,
		AllowedBiomarkers: []string{"blood_glucose", "heart_rate"},
	})

	// 허용된 바이오마커
	allowed, _, err := svc.ValidateSharingAccess(ctx, group.ID, "owner-br", "user-br", "blood_glucose")
	if err != nil {
		t.Fatalf("에러 발생: %v", err)
	}
	if !allowed {
		t.Fatal("허용된 바이오마커 접근이 허용되어야 합니다")
	}

	// 허용되지 않은 바이오마커
	allowed, _, err = svc.ValidateSharingAccess(ctx, group.ID, "owner-br", "user-br", "cholesterol_total")
	if err != nil {
		t.Fatalf("에러 발생: %v", err)
	}
	if allowed {
		t.Fatal("허용되지 않은 바이오마커 접근이 거부되어야 합니다")
	}
}

func TestSharingPreferencesWithLimits(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	group, _ := svc.CreateFamilyGroup(ctx, "owner-lim", "제한 가족", "")

	pref := &service.SharingPreferences{
		UserID:               "owner-lim",
		GroupID:              group.ID,
		ShareMeasurements:    true,
		ShareHealthScore:     true,
		MeasurementDaysLimit: 30,
		AllowedBiomarkers:    []string{"blood_glucose", "blood_pressure"},
		RequireApproval:      false,
	}

	result, err := svc.SetSharingPreferences(ctx, pref)
	if err != nil {
		t.Fatalf("공유 설정 실패: %v", err)
	}
	if result.MeasurementDaysLimit != 30 {
		t.Fatalf("MeasurementDaysLimit 불일치: got %d, want 30", result.MeasurementDaysLimit)
	}
	if len(result.AllowedBiomarkers) != 2 {
		t.Fatalf("AllowedBiomarkers 수 불일치: got %d, want 2", len(result.AllowedBiomarkers))
	}
	if result.RequireApproval {
		t.Fatal("RequireApproval이 false여야 합니다")
	}

	// MeasurementDaysLimit가 ValidateSharingAccess에 반영되는지 확인
	inv, _ := svc.InviteMember(ctx, group.ID, "owner-lim", "viewer@test.com", service.RoleMember, "")
	svc.RespondToInvitation(ctx, inv.ID, "viewer-lim", true)

	allowed, reason, err := svc.ValidateSharingAccess(ctx, group.ID, "viewer-lim", "owner-lim", "blood_glucose")
	if err != nil {
		t.Fatalf("ValidateSharingAccess 실패: %v", err)
	}
	if !allowed {
		t.Fatalf("접근이 허용되어야 합니다: reason=%s", reason)
	}
	// reason에 일수 제한 정보가 포함되어야 함
	if reason == "" {
		t.Fatal("reason이 비어 있으면 안 됩니다")
	}
}

func TestEndToEnd_FamilyFlow(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	// 1. 가족 그룹 생성
	group, _ := svc.CreateFamilyGroup(ctx, "dad", "우리집", "화목한 가족")

	// 2. 엄마 초대 (Guardian)
	inv1, _ := svc.InviteMember(ctx, group.ID, "dad", "mom@family.com", service.RoleGuardian, "가족으로 초대합니다")
	svc.RespondToInvitation(ctx, inv1.ID, "mom", true)

	// 3. 자녀 초대 (Child)
	inv2, _ := svc.InviteMember(ctx, group.ID, "dad", "child@family.com", service.RoleChild, "")
	svc.RespondToInvitation(ctx, inv2.ID, "child", true)

	// 4. 어르신 초대 (Elderly)
	inv3, _ := svc.InviteMember(ctx, group.ID, "dad", "grandma@family.com", service.RoleElderly, "")
	svc.RespondToInvitation(ctx, inv3.ID, "grandma", true)

	// 5. 멤버 수 확인
	_, count, _ := svc.GetFamilyGroup(ctx, group.ID)
	if count != 4 {
		t.Fatalf("전체 멤버 수 불일치: got %d, want 4", count)
	}

	// 6. 공유 설정
	svc.SetSharingPreferences(ctx, &service.SharingPreferences{
		UserID:           "grandma",
		GroupID:          group.ID,
		ShareHealthScore: true,
		ShareAlerts:      true,
	})

	// 7. 엄마(Guardian)가 전체 데이터 조회
	summaries, err := svc.GetSharedHealthData(ctx, "mom", "", group.ID, 30)
	if err != nil {
		t.Fatalf("건강 데이터 조회 실패: %v", err)
	}
	if len(summaries) != 3 { // dad, child, grandma (본인 제외)
		t.Fatalf("요약 수 불일치: got %d, want 3", len(summaries))
	}

	// 8. 역할 변경
	svc.UpdateMemberRole(ctx, group.ID, "child", service.RoleMember)
	members, _ := svc.ListFamilyMembers(ctx, group.ID)
	for _, m := range members {
		if m.UserID == "child" && m.Role != service.RoleMember {
			t.Fatalf("역할 변경 실패: got %d, want %d", m.Role, service.RoleMember)
		}
	}
}
