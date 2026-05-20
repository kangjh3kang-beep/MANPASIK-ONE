package tenancy

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

// InvitationStatus 는 초대 상태.
type InvitationStatus string

const (
	InvitationPending  InvitationStatus = "pending"
	InvitationAccepted InvitationStatus = "accepted"
	InvitationRevoked  InvitationStatus = "revoked"
	InvitationExpired  InvitationStatus = "expired"
)

// Invitation 은 조직 가입 초대.
//
// 초대장은 일회용 토큰을 가지며, 수락 시 Membership 으로 변환되고
// Status 는 accepted 로 갱신됨.
type Invitation struct {
	Token       string           // 32 hex chars (16 random bytes)
	TenantID    TenantID
	InviterID   string           // 초대를 보낸 사용자
	InviteeHint string           // 이메일/전화 등 (보낸 대상 식별; 인증된 사용자 ID 가 일치하는지는 별도 검증)
	Role        TenantRole       // 수락 후 부여할 역할
	Status      InvitationStatus
	IssuedAt    time.Time
	ExpiresAt   time.Time
	AcceptedBy  string           // 수락한 사용자 (수락 후 채워짐)
	AcceptedAt  *time.Time
}

// IsExpired 는 만료 여부.
func (i *Invitation) IsExpired() bool {
	return time.Now().After(i.ExpiresAt)
}

// IsActive 는 수락 가능한 상태인지.
func (i *Invitation) IsActive() bool {
	return i.Status == InvitationPending && !i.IsExpired()
}

// InvitationStore 는 초대 저장소.
type InvitationStore interface {
	Add(inv Invitation) error
	Get(token string) (*Invitation, error)
	Update(inv Invitation) error
	ListByTenant(tenantID TenantID) []*Invitation
}

// MemoryInvitationStore 는 sync.RWMutex 기반 인메모리 저장소.
type MemoryInvitationStore struct {
	mu          sync.RWMutex
	invitations map[string]*Invitation
}

// NewMemoryInvitationStore 생성.
func NewMemoryInvitationStore() *MemoryInvitationStore {
	return &MemoryInvitationStore{invitations: make(map[string]*Invitation)}
}

// Add 신규 초대 등록 (token 충돌 시 에러).
func (s *MemoryInvitationStore) Add(inv Invitation) error {
	if inv.Token == "" {
		return errors.New("Token 필수")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.invitations[inv.Token]; exists {
		return errors.New("Token 중복")
	}
	cp := inv
	s.invitations[inv.Token] = &cp
	return nil
}

// Get 토큰으로 조회.
func (s *MemoryInvitationStore) Get(token string) (*Invitation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	inv, ok := s.invitations[token]
	if !ok {
		return nil, ErrInvitationNotFound
	}
	cp := *inv
	return &cp, nil
}

// Update 초대 갱신 (status, AcceptedBy 등).
func (s *MemoryInvitationStore) Update(inv Invitation) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.invitations[inv.Token]; !ok {
		return ErrInvitationNotFound
	}
	cp := inv
	s.invitations[inv.Token] = &cp
	return nil
}

// ListByTenant 는 조직별 초대 목록.
func (s *MemoryInvitationStore) ListByTenant(tenantID TenantID) []*Invitation {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []*Invitation
	for _, inv := range s.invitations {
		if inv.TenantID == tenantID {
			cp := *inv
			out = append(out, &cp)
		}
	}
	return out
}

// 에러 정의.
var (
	ErrInvitationNotFound = errors.New("초대 없음")
	ErrInvitationExpired  = errors.New("초대 만료")
	ErrInvitationConsumed = errors.New("이미 처리된 초대")
	ErrInviteeMismatch    = errors.New("수락자 불일치")
)

// generateToken 은 16 바이트 랜덤 헥스 (32 chars) 생성.
func generateToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
