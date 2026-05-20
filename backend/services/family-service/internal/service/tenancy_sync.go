package service

import (
	"context"

	"github.com/manpasik/backend/shared/tenancy"
)

// TenancySyncAdapter 는 family-service.TenancySync 인터페이스의 표준 구현.
//
// tenancy.MembershipStore 와 직접 통신하여 family 도메인의 그룹/멤버 변경을
// tenancy 측 영속화. main.go 에서 PostgresMembershipStore (또는 인메모리) 를
// 주입.
//
// Phase AK-1: 옵션 webhook 발송 — 그룹/멤버 변경 시 외부 시스템에 자동 알림.
//
// 역할 매핑:
//
//	family.RoleOwner       → tenancy.TenantRoleOwner
//	family.RoleGuardian    → tenancy.TenantRoleAdmin (보호자 = 그룹 운영)
//	family.RoleMember      → tenancy.TenantRoleMember
//	family.RoleChild       → tenancy.TenantRoleMember (자녀)
//	family.RoleElderly     → tenancy.TenantRoleMember (노인)
//	기타                   → tenancy.TenantRoleViewer
type TenancySyncAdapter struct {
	memStore tenancy.MembershipStore
	webhook  *tenancy.WebhookDispatcher // 옵션 — 멤버 변경 webhook 자동 발송
}

// NewTenancySyncAdapter 생성. memStore 는 필수.
func NewTenancySyncAdapter(memStore tenancy.MembershipStore) *TenancySyncAdapter {
	return &TenancySyncAdapter{memStore: memStore}
}

// SetWebhookDispatcher 는 멤버 변경 시 webhook 자동 발송 (Phase AK-1).
//
// 미설정 시 webhook 미발송. Invite/Accept 흐름은 InvitationService 가
// 별도로 webhook 발송하므로, 여기서는 family-service 가 직접 처리하는
// 그룹 생성/멤버 추가/제거/역할 변경 4가지 이벤트만 송신.
func (a *TenancySyncAdapter) SetWebhookDispatcher(d *tenancy.WebhookDispatcher) {
	a.webhook = d
}

// OnGroupCreated 는 그룹 생성 시 owner 를 tenancy 의 owner 멤버십으로 등록.
func (a *TenancySyncAdapter) OnGroupCreated(_ context.Context, groupID, ownerUserID string) error {
	if a.memStore == nil {
		return nil
	}
	if err := a.memStore.Add(tenancy.Membership{
		UserID:   ownerUserID,
		TenantID: tenancy.TenantID(groupID),
		Role:     tenancy.TenantRoleOwner,
	}); err != nil {
		return err
	}
	a.dispatchEvent(tenancy.EventMembershipCreated, groupID, ownerUserID, ownerUserID,
		map[string]string{"role": "owner", "source": "family.group_created"})
	return nil
}

// OnMemberAdded 는 가족 멤버 추가 시 tenancy 멤버십 등록.
func (a *TenancySyncAdapter) OnMemberAdded(_ context.Context, groupID, userID, familyRole string) error {
	if a.memStore == nil {
		return nil
	}
	role := mapFamilyRoleToTenant(familyRole)
	if err := a.memStore.Add(tenancy.Membership{
		UserID:   userID,
		TenantID: tenancy.TenantID(groupID),
		Role:     role,
	}); err != nil {
		return err
	}
	a.dispatchEvent(tenancy.EventMembershipCreated, groupID, userID, userID,
		map[string]string{"role": string(role), "family_role": familyRole})
	return nil
}

// OnMemberRemoved 는 가족 멤버 제거 시 tenancy 측도 제거.
func (a *TenancySyncAdapter) OnMemberRemoved(_ context.Context, groupID, userID string) error {
	if a.memStore == nil {
		return nil
	}
	if err := a.memStore.Remove(userID, tenancy.TenantID(groupID)); err != nil {
		return err
	}
	a.dispatchEvent(tenancy.EventMembershipRemoved, groupID, userID, "",
		map[string]string{"source": "family.member_removed"})
	return nil
}

// OnMemberRoleChanged 는 역할 변경 시 동기화.
func (a *TenancySyncAdapter) OnMemberRoleChanged(_ context.Context, groupID, userID, newFamilyRole string) error {
	if a.memStore == nil {
		return nil
	}
	role := mapFamilyRoleToTenant(newFamilyRole)
	if err := a.memStore.UpdateRole(userID, tenancy.TenantID(groupID), role); err != nil {
		return err
	}
	a.dispatchEvent(tenancy.EventMembershipRoleChanged, groupID, userID, "",
		map[string]string{"new_role": string(role), "family_role": newFamilyRole})
	return nil
}

// dispatchEvent 는 webhook 이 설정된 경우 비동기 이벤트 발송.
func (a *TenancySyncAdapter) dispatchEvent(t tenancy.EventType, groupID, userID, actorID string,
	payload map[string]string) {
	if a.webhook == nil {
		return
	}
	a.webhook.DispatchAsync(tenancy.Event{
		Type:     t,
		TenantID: groupID,
		UserID:   userID,
		ActorID:  actorID,
		Payload:  payload,
	})
}

// mapFamilyRoleToTenant 는 family.FamilyRole → tenancy.TenantRole 매핑.
//
// FamilyRole 은 int 이지만 family.go 의 String() 메서드로 안정적인 이름 표현
// 을 보장. 호출자가 String() 결과를 전달.
func mapFamilyRoleToTenant(familyRole string) tenancy.TenantRole {
	switch familyRole {
	case "owner":
		return tenancy.TenantRoleOwner
	case "guardian":
		return tenancy.TenantRoleAdmin
	case "member", "child", "elderly":
		return tenancy.TenantRoleMember
	default:
		return tenancy.TenantRoleViewer
	}
}
