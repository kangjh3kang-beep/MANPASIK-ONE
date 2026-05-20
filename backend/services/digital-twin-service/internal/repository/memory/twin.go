package memory

import (
	"context"
	"sort"
	"sync"

	"github.com/manpasik/backend/services/digital-twin-service/internal/service"
)

// TwinRepository는 디지털 트윈 상태의 인메모리 저장소입니다.
type TwinRepository struct {
	mu     sync.RWMutex
	states map[string]*service.TwinState // key: sessionID
}

// NewTwinRepository는 새 인메모리 트윈 저장소를 생성합니다.
func NewTwinRepository() *TwinRepository {
	return &TwinRepository{
		states: make(map[string]*service.TwinState),
	}
}

func (r *TwinRepository) Store(ctx context.Context, state *service.TwinState) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.states[state.SessionID] = state
	return nil
}

func (r *TwinRepository) Get(ctx context.Context, sessionID string) (*service.TwinState, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.states[sessionID], nil
}

func (r *TwinRepository) ListByDevice(ctx context.Context, deviceID string, limit int) ([]*service.TwinState, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*service.TwinState
	for _, s := range r.states {
		if s.DeviceID == deviceID {
			result = append(result, s)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].LastSyncedAt.After(result[j].LastSyncedAt)
	})

	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}
	return result, nil
}
