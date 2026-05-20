package memory

import (
	"context"
	"sync"

	"github.com/manpasik/backend/services/assistant-service/internal/service"
)

// AssistantRepository는 인메모리 비서 저장소입니다.
type AssistantRepository struct {
	mu       sync.RWMutex
	sessions map[string]*service.AssistantSession
	byUser   map[string][]string // userID -> sessionIDs
	turns    map[string][]*service.AssistantTurn // sessionID -> turns
}

func NewAssistantRepository() *AssistantRepository {
	return &AssistantRepository{
		sessions: make(map[string]*service.AssistantSession),
		byUser:   make(map[string][]string),
		turns:    make(map[string][]*service.AssistantTurn),
	}
}

func (r *AssistantRepository) CreateSession(_ context.Context, s *service.AssistantSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *s
	r.sessions[s.ID] = &cp
	r.byUser[s.UserID] = append(r.byUser[s.UserID], s.ID)
	return nil
}

func (r *AssistantRepository) GetSession(_ context.Context, id string) (*service.AssistantSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.sessions[id]
	if !ok {
		return nil, nil
	}
	cp := *s
	return &cp, nil
}

func (r *AssistantRepository) ListSessions(_ context.Context, userID string, limit, offset int32) ([]*service.AssistantSession, int32, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := r.byUser[userID]
	total := int32(len(ids))
	start := int(offset)
	if start >= len(ids) {
		return nil, total, nil
	}
	end := start + int(limit)
	if end > len(ids) {
		end = len(ids)
	}
	result := make([]*service.AssistantSession, 0, end-start)
	for _, id := range ids[start:end] {
		if s, ok := r.sessions[id]; ok {
			cp := *s
			result = append(result, &cp)
		}
	}
	return result, total, nil
}

func (r *AssistantRepository) DeleteSession(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	s, ok := r.sessions[id]
	if !ok {
		return nil
	}
	delete(r.sessions, id)
	delete(r.turns, id)
	ids := r.byUser[s.UserID]
	for i, sid := range ids {
		if sid == id {
			r.byUser[s.UserID] = append(ids[:i], ids[i+1:]...)
			break
		}
	}
	return nil
}

func (r *AssistantRepository) AddTurn(_ context.Context, t *service.AssistantTurn) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *t
	r.turns[t.SessionID] = append(r.turns[t.SessionID], &cp)
	return nil
}

func (r *AssistantRepository) ListTurns(_ context.Context, sessionID string, limit, offset int32) ([]*service.AssistantTurn, int32, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	turns := r.turns[sessionID]
	total := int32(len(turns))
	start := int(offset)
	if start >= len(turns) {
		return nil, total, nil
	}
	end := start + int(limit)
	if end > len(turns) {
		end = len(turns)
	}
	result := make([]*service.AssistantTurn, 0, end-start)
	for _, t := range turns[start:end] {
		cp := *t
		result = append(result, &cp)
	}
	return result, total, nil
}

func (r *AssistantRepository) UpdateSession(_ context.Context, s *service.AssistantSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *s
	r.sessions[s.ID] = &cp
	return nil
}
