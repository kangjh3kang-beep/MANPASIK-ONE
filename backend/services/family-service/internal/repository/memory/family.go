// Package memory는 family-service의 인메모리 저장소입니다.
package memory

import (
	"context"
	"sync"

	"github.com/manpasik/backend/services/family-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
)

// FamilyGroupRepository는 가족 그룹 인메모리 저장소입니다.
type FamilyGroupRepository struct {
	mu    sync.RWMutex
	store map[string]*service.FamilyGroup
}

func NewFamilyGroupRepository() *FamilyGroupRepository {
	return &FamilyGroupRepository{store: make(map[string]*service.FamilyGroup)}
}

func (r *FamilyGroupRepository) Save(_ context.Context, g *service.FamilyGroup) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[g.ID] = g
	return nil
}

func (r *FamilyGroupRepository) FindByID(_ context.Context, id string) (*service.FamilyGroup, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, ok := r.store[id]
	if !ok {
		return nil, apperrors.New(apperrors.ErrNotFound, "리소스를 찾을 수 없습니다")
	}
	return g, nil
}

// FamilyMemberRepository는 가족 멤버 인메모리 저장소입니다.
type FamilyMemberRepository struct {
	mu    sync.RWMutex
	store map[string]*service.FamilyMember // key: "userID:groupID"
}

func NewFamilyMemberRepository() *FamilyMemberRepository {
	return &FamilyMemberRepository{store: make(map[string]*service.FamilyMember)}
}

func memberKey(userID, groupID string) string {
	return userID + ":" + groupID
}

func (r *FamilyMemberRepository) Save(_ context.Context, m *service.FamilyMember) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[memberKey(m.UserID, m.GroupID)] = m
	return nil
}

func (r *FamilyMemberRepository) FindByGroupID(_ context.Context, groupID string) ([]*service.FamilyMember, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var members []*service.FamilyMember
	for _, m := range r.store {
		if m.GroupID == groupID {
			members = append(members, m)
		}
	}
	return members, nil
}

func (r *FamilyMemberRepository) FindByUserIDAndGroupID(_ context.Context, userID, groupID string) (*service.FamilyMember, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.store[memberKey(userID, groupID)]
	if !ok {
		return nil, apperrors.New(apperrors.ErrNotFound, "리소스를 찾을 수 없습니다")
	}
	return m, nil
}

func (r *FamilyMemberRepository) Remove(_ context.Context, userID, groupID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := memberKey(userID, groupID)
	if _, ok := r.store[key]; !ok {
		return apperrors.New(apperrors.ErrNotFound, "리소스를 찾을 수 없습니다")
	}
	delete(r.store, key)
	return nil
}

func (r *FamilyMemberRepository) CountByGroupID(_ context.Context, groupID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, m := range r.store {
		if m.GroupID == groupID {
			count++
		}
	}
	return count, nil
}

// InvitationRepository는 초대 인메모리 저장소입니다.
type InvitationRepository struct {
	mu    sync.RWMutex
	store map[string]*service.FamilyInvitation
}

func NewInvitationRepository() *InvitationRepository {
	return &InvitationRepository{store: make(map[string]*service.FamilyInvitation)}
}

func (r *InvitationRepository) Save(_ context.Context, inv *service.FamilyInvitation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[inv.ID] = inv
	return nil
}

func (r *InvitationRepository) FindByID(_ context.Context, id string) (*service.FamilyInvitation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	inv, ok := r.store[id]
	if !ok {
		return nil, apperrors.New(apperrors.ErrNotFound, "리소스를 찾을 수 없습니다")
	}
	return inv, nil
}

func (r *InvitationRepository) Update(_ context.Context, inv *service.FamilyInvitation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[inv.ID] = inv
	return nil
}

// SharingPreferencesRepository는 공유 설정 인메모리 저장소입니다.
type SharingPreferencesRepository struct {
	mu    sync.RWMutex
	store map[string]*service.SharingPreferences // key: "userID:groupID"
}

func NewSharingPreferencesRepository() *SharingPreferencesRepository {
	return &SharingPreferencesRepository{store: make(map[string]*service.SharingPreferences)}
}

func (r *SharingPreferencesRepository) Save(_ context.Context, pref *service.SharingPreferences) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[pref.UserID+":"+pref.GroupID] = pref
	return nil
}

func (r *SharingPreferencesRepository) FindByUserIDAndGroupID(_ context.Context, userID, groupID string) (*service.SharingPreferences, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	pref, ok := r.store[userID+":"+groupID]
	if !ok {
		return nil, nil
	}
	return pref, nil
}
