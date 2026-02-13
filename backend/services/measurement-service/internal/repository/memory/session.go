// Package memory는 measurement-service의 인메모리 저장소 구현입니다.
package memory

import (
	"context"
	"sync"
	"time"

	"github.com/manpasik/backend/services/measurement-service/internal/service"
)

// SessionRepository는 인메모리 세션 저장소입니다.
type SessionRepository struct {
	mu       sync.RWMutex
	sessions map[string]*service.MeasurementSession // key: sessionID
}

// NewSessionRepository는 인메모리 SessionRepository를 생성합니다.
func NewSessionRepository() *SessionRepository {
	return &SessionRepository{
		sessions: make(map[string]*service.MeasurementSession),
	}
}

func (r *SessionRepository) CreateSession(_ context.Context, session *service.MeasurementSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[session.ID] = session
	return nil
}

func (r *SessionRepository) GetSession(_ context.Context, sessionID string) (*service.MeasurementSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.sessions[sessionID]
	if !ok {
		return nil, nil
	}
	cp := *s
	return &cp, nil
}

func (r *SessionRepository) EndSession(_ context.Context, sessionID string, totalMeasurements int, endedAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	s, ok := r.sessions[sessionID]
	if !ok {
		return nil
	}
	s.TotalMeasurements = totalMeasurements
	s.EndedAt = &endedAt
	s.Status = "completed"
	return nil
}
