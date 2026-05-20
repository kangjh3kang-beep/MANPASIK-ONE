package tenancy

import (
	"sync"
	"time"
)

// MembershipStore 는 사용자-조직-역할 매핑 저장소.
type MembershipStore interface {
	Add(m Membership) error
	Remove(userID string, tenantID TenantID) error
	Get(userID string, tenantID TenantID) (*Membership, error)
	ListUserTenants(userID string) []*Membership
	ListTenantMembers(tenantID TenantID) []*Membership
	UpdateRole(userID string, tenantID TenantID, newRole TenantRole) error
	SetActive(userID string, tenantID TenantID, active bool) error
}

// MemoryMembershipStore 는 sync.RWMutex 기반 인메모리 저장소.
//
// 키 형식: "{userID}::{tenantID}" (충돌 방지용 더블콜론).
type MemoryMembershipStore struct {
	mu      sync.RWMutex
	members map[string]*Membership
}

// NewMemoryMembershipStore 생성.
func NewMemoryMembershipStore() *MemoryMembershipStore {
	return &MemoryMembershipStore{
		members: make(map[string]*Membership),
	}
}

func memberKey(userID string, tenantID TenantID) string {
	return userID + "::" + string(tenantID)
}

// Add 는 새 멤버십 등록 (중복 시 덮어쓰기).
func (s *MemoryMembershipStore) Add(m Membership) error {
	if err := m.Validate(); err != nil {
		return err
	}
	if m.JoinedAt == 0 {
		m.JoinedAt = time.Now().Unix()
	}
	if !m.Active && m.JoinedAt == time.Now().Unix() {
		// 새 멤버는 기본 활성
		m.Active = true
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := m
	s.members[memberKey(m.UserID, m.TenantID)] = &cp
	return nil
}

// Remove 멤버십 삭제.
func (s *MemoryMembershipStore) Remove(userID string, tenantID TenantID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.members, memberKey(userID, tenantID))
	return nil
}

// Get 멤버십 조회 (없으면 (nil, ErrNoMembership)).
func (s *MemoryMembershipStore) Get(userID string, tenantID TenantID) (*Membership, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.members[memberKey(userID, tenantID)]
	if !ok {
		return nil, ErrNoMembership
	}
	cp := *m
	return &cp, nil
}

// ListUserTenants 는 사용자가 속한 모든 조직 멤버십.
func (s *MemoryMembershipStore) ListUserTenants(userID string) []*Membership {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []*Membership{}
	for _, m := range s.members {
		if m.UserID == userID {
			cp := *m
			out = append(out, &cp)
		}
	}
	return out
}

// ListTenantMembers 는 조직의 모든 멤버.
func (s *MemoryMembershipStore) ListTenantMembers(tenantID TenantID) []*Membership {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []*Membership{}
	for _, m := range s.members {
		if m.TenantID == tenantID {
			cp := *m
			out = append(out, &cp)
		}
	}
	return out
}

// UpdateRole 멤버의 역할 변경.
func (s *MemoryMembershipStore) UpdateRole(userID string, tenantID TenantID, newRole TenantRole) error {
	if !newRole.IsKnown() {
		return ErrCrossTenant // 부정확하지만 임시; 실제로는 별도 에러
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	m, ok := s.members[memberKey(userID, tenantID)]
	if !ok {
		return ErrNoMembership
	}
	m.Role = newRole
	return nil
}

// SetActive 활성/비활성 토글.
func (s *MemoryMembershipStore) SetActive(userID string, tenantID TenantID, active bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	m, ok := s.members[memberKey(userID, tenantID)]
	if !ok {
		return ErrNoMembership
	}
	m.Active = active
	return nil
}
