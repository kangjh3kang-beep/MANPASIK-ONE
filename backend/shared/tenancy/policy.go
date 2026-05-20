package tenancy

// PolicyEngine 은 (사용자, 리소스, 작업) → Decision 평가.
type PolicyEngine struct {
	store MembershipStore
	// rolePerms 는 역할별 허용 작업.
	rolePerms map[TenantRole]map[Action]bool
	// allowOwnerOverride 가 true 면 OwnerID == userID 시 항상 허용.
	allowOwnerOverride bool
}

// NewPolicyEngine 은 기본 역할-권한 매트릭스로 초기화.
func NewPolicyEngine(store MembershipStore) *PolicyEngine {
	return &PolicyEngine{
		store:              store,
		rolePerms:          defaultRolePerms(),
		allowOwnerOverride: true,
	}
}

func defaultRolePerms() map[TenantRole]map[Action]bool {
	allActions := func() map[Action]bool {
		return map[Action]bool{
			ActionRead: true, ActionWrite: true, ActionDelete: true,
			ActionShare: true, ActionAdmin: true,
		}
	}
	rwShare := func() map[Action]bool {
		return map[Action]bool{
			ActionRead: true, ActionWrite: true, ActionShare: true,
		}
	}
	rw := func() map[Action]bool {
		return map[Action]bool{
			ActionRead: true, ActionWrite: true,
		}
	}
	readOnly := func() map[Action]bool {
		return map[Action]bool{ActionRead: true}
	}

	return map[TenantRole]map[Action]bool{
		TenantRoleOwner:        allActions(),
		TenantRoleAdmin:        allActions(),
		TenantRoleMedicalStaff: rwShare(),
		TenantRoleMember:       rw(),
		TenantRoleViewer:       readOnly(),
	}
}

// SetRolePermission 은 특정 역할-작업 조합 허용/거부 설정.
func (p *PolicyEngine) SetRolePermission(role TenantRole, action Action, allowed bool) {
	if p.rolePerms[role] == nil {
		p.rolePerms[role] = make(map[Action]bool)
	}
	p.rolePerms[role][action] = allowed
}

// SetAllowOwnerOverride 는 소유자 자동 허용 여부 설정.
func (p *PolicyEngine) SetAllowOwnerOverride(b bool) {
	p.allowOwnerOverride = b
}

// Evaluate 는 (사용자가 자신의 컨텍스트 테넌트 currentTenant 에서) 리소스에
// action 을 수행할 수 있는지 평가.
//
// 평가 순서:
//  1. currentTenant 와 resource.TenantID 가 다르면 즉시 거부 (교차 조직 차단)
//  2. resource.OwnerID == userID 이고 owner override 면 허용
//  3. resource.IsSharedWith(userID) 이고 action == read 면 허용
//  4. 멤버십 조회 → 비활성이면 거부
//  5. 역할의 권한 매트릭스 확인
func (p *PolicyEngine) Evaluate(userID string, currentTenant TenantID, resource *Resource, action Action) Decision {
	if userID == "" {
		return Decision{Reason: "사용자 ID 없음"}
	}
	if currentTenant.IsZero() {
		return Decision{Reason: ErrNoTenantContext.Error()}
	}
	if resource == nil {
		return Decision{Reason: "리소스 없음"}
	}

	// 1. 교차 조직 차단
	if !resource.TenantID.IsZero() && resource.TenantID != currentTenant {
		return Decision{Reason: ErrCrossTenant.Error()}
	}

	// 2. 소유자 우선 허용 (자기 데이터)
	if p.allowOwnerOverride && resource.OwnerID != "" && resource.OwnerID == userID {
		return Decision{Allowed: true, Reason: "owner_override"}
	}

	// 3. 명시 공유 (read 한정)
	if action == ActionRead && resource.IsSharedWith(userID) {
		return Decision{Allowed: true, Reason: "shared_resource"}
	}

	// 4. 멤버십 확인
	m, err := p.store.Get(userID, currentTenant)
	if err != nil {
		return Decision{Reason: ErrNoMembership.Error()}
	}
	if !m.Active {
		return Decision{Reason: ErrInactiveMember.Error()}
	}

	// 5. 역할 권한 매트릭스
	perms, ok := p.rolePerms[m.Role]
	if !ok {
		return Decision{Reason: "역할 권한 미정의: " + string(m.Role)}
	}
	if perms[action] {
		return Decision{Allowed: true, Reason: "role_grant", MatchedRole: m.Role}
	}

	return Decision{
		Reason:      "역할 '" + string(m.Role) + "' 에 작업 '" + string(action) + "' 권한 없음",
		MatchedRole: m.Role,
	}
}

// CheckMembership 은 단순히 사용자가 조직의 활성 멤버인지 확인.
func (p *PolicyEngine) CheckMembership(userID string, tenantID TenantID) (*Membership, error) {
	m, err := p.store.Get(userID, tenantID)
	if err != nil {
		return nil, err
	}
	if !m.Active {
		return nil, ErrInactiveMember
	}
	return m, nil
}

// FilterResources 는 리소스 목록 중 사용자가 read 가능한 것만 반환.
func (p *PolicyEngine) FilterResources(userID string, currentTenant TenantID, resources []*Resource) []*Resource {
	out := make([]*Resource, 0, len(resources))
	for _, r := range resources {
		if d := p.Evaluate(userID, currentTenant, r, ActionRead); d.Allowed {
			out = append(out, r)
		}
	}
	return out
}
