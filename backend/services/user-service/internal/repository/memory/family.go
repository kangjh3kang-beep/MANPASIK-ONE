package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/manpasik/backend/services/user-service/internal/service"
)

// FamilyRepository는 인메모리 가족 그룹 저장소입니다.
type FamilyRepository struct {
	mu     sync.RWMutex
	groups map[string]*service.FamilyGroup // key: groupID
}

// NewFamilyRepository는 인메모리 FamilyRepository를 생성합니다.
func NewFamilyRepository() *FamilyRepository {
	return &FamilyRepository{
		groups: make(map[string]*service.FamilyGroup),
	}
}

func (r *FamilyRepository) GetGroup(_ context.Context, groupID string) (*service.FamilyGroup, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, ok := r.groups[groupID]
	if !ok {
		return nil, nil
	}
	return g, nil
}

func (r *FamilyRepository) GetUserGroups(_ context.Context, userID string) ([]*service.FamilyGroup, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*service.FamilyGroup
	for _, g := range r.groups {
		if g.OwnerID == userID {
			result = append(result, g)
			continue
		}
		for _, m := range g.Members {
			if m.UserID == userID {
				result = append(result, g)
				break
			}
		}
	}
	return result, nil
}

func (r *FamilyRepository) CreateGroup(_ context.Context, group *service.FamilyGroup) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.groups[group.ID] = group
	return nil
}

func (r *FamilyRepository) AddMember(_ context.Context, groupID, userID, role string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	g, ok := r.groups[groupID]
	if !ok {
		return fmt.Errorf("group %s not found", groupID)
	}
	g.Members = append(g.Members, &service.FamilyMember{UserID: userID, Role: role})
	return nil
}

func (r *FamilyRepository) RemoveMember(_ context.Context, groupID, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	g, ok := r.groups[groupID]
	if !ok {
		return fmt.Errorf("group %s not found", groupID)
	}
	for i, m := range g.Members {
		if m.UserID == userID {
			g.Members = append(g.Members[:i], g.Members[i+1:]...)
			return nil
		}
	}
	return nil
}
