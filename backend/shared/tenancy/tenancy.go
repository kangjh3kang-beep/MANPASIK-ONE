// Package tenancy 는 만파식의 멀티테넌트(조직 단위) 권한 격리를 제공합니다.
//
// 단일 역할 RBAC(middleware/rbac.go)와는 별도 레이어로 동작:
//   - middleware/rbac.go: gRPC 메서드 단위의 역할 게이트
//   - tenancy/*: 리소스 단위(레코드 수준) 조직 격리 + per-org 역할
//
// 한 사용자는 여러 조직에 속할 수 있고, 조직마다 다른 역할을 가질 수 있습니다.
// 예: 김의사는 A병원에서 admin, B병원에서 viewer.
package tenancy

import (
	"errors"
	"strings"
)

// TenantID 는 조직(병원/연구소/가족 그룹) 식별자.
type TenantID string

func (t TenantID) String() string { return string(t) }
func (t TenantID) IsZero() bool   { return t == "" }

// TenantRole 은 조직 내에서 사용자가 가지는 역할.
//
// 시스템 전역 역할(middleware.RoleAdmin)과는 별도로 운영되며,
// 사용자는 시스템 전역 user 이지만 특정 조직에서 owner 일 수 있음.
type TenantRole string

const (
	// TenantRoleOwner 는 조직의 최종 의사결정권자 (소유자).
	TenantRoleOwner TenantRole = "owner"
	// TenantRoleAdmin 은 조직의 운영자 (멤버 추가/제거, 설정).
	TenantRoleAdmin TenantRole = "admin"
	// TenantRoleMedicalStaff 는 의료 전문가 (의사/간호사/약사).
	TenantRoleMedicalStaff TenantRole = "medical_staff"
	// TenantRoleMember 는 일반 멤버 (자신의 데이터 + 공유받은 데이터 접근).
	TenantRoleMember TenantRole = "member"
	// TenantRoleViewer 는 읽기 전용.
	TenantRoleViewer TenantRole = "viewer"
)

// Action 은 리소스에 대한 작업 종류.
type Action string

const (
	ActionRead   Action = "read"
	ActionWrite  Action = "write"
	ActionDelete Action = "delete"
	ActionShare  Action = "share"
	ActionAdmin  Action = "admin"
)

// Membership 은 사용자-조직-역할 연결.
type Membership struct {
	UserID   string
	TenantID TenantID
	Role     TenantRole
	// JoinedAt unix timestamp (생성 시각).
	JoinedAt int64
	// Active 가 false 면 비활성 멤버 (탈퇴, 정직 등).
	Active bool
}

// Validate 는 Membership 의 필수 필드를 검증.
func (m *Membership) Validate() error {
	if m.UserID == "" {
		return errors.New("UserID 필수")
	}
	if m.TenantID.IsZero() {
		return errors.New("TenantID 필수")
	}
	if m.Role == "" {
		return errors.New("Role 필수")
	}
	if !m.Role.IsKnown() {
		return errors.New("알 수 없는 역할: " + string(m.Role))
	}
	return nil
}

// IsKnown 은 정의된 역할인지 확인.
func (r TenantRole) IsKnown() bool {
	switch r {
	case TenantRoleOwner, TenantRoleAdmin, TenantRoleMedicalStaff,
		TenantRoleMember, TenantRoleViewer:
		return true
	}
	return false
}

// Resource 는 테넌시 검사 대상.
type Resource struct {
	// TenantID 는 리소스가 속한 조직 (필수).
	TenantID TenantID
	// OwnerID 는 리소스를 생성한 사용자 (선택).
	OwnerID string
	// Type 은 리소스 종류 (예: "measurement", "prescription").
	Type string
	// SharedWith 는 명시적으로 공유받은 사용자 ID 목록.
	SharedWith []string
}

// IsSharedWith 는 사용자에게 명시 공유되었는지 확인.
func (r *Resource) IsSharedWith(userID string) bool {
	for _, uid := range r.SharedWith {
		if uid == userID {
			return true
		}
	}
	return false
}

// Decision 은 권한 평가 결과.
type Decision struct {
	Allowed bool
	Reason  string
	// MatchedRole 은 Allowed=true 일 때 어떤 역할이 매칭되었는지.
	MatchedRole TenantRole
}

func (d Decision) String() string {
	var sb strings.Builder
	if d.Allowed {
		sb.WriteString("ALLOW")
	} else {
		sb.WriteString("DENY")
	}
	if d.Reason != "" {
		sb.WriteString(": ")
		sb.WriteString(d.Reason)
	}
	return sb.String()
}

// 에러 정의
var (
	ErrNoTenantContext = errors.New("테넌트 컨텍스트 없음")
	ErrNoMembership    = errors.New("조직 멤버 아님")
	ErrInactiveMember  = errors.New("비활성 멤버")
	ErrCrossTenant     = errors.New("교차 조직 접근 차단")
)
