package memory

import (
	"context"
	"sync"

	"github.com/manpasik/backend/services/ai-inference-service/internal/service"
)

// ============================================================================
// In-Memory Analysis Repository
// ============================================================================

type AnalysisRepository struct {
	mu   sync.RWMutex
	data map[string]*service.AnalysisResult
}

func NewAnalysisRepository() *AnalysisRepository {
	return &AnalysisRepository{data: make(map[string]*service.AnalysisResult)}
}

func (r *AnalysisRepository) Save(_ context.Context, result *service.AnalysisResult) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[result.AnalysisID] = result
	return nil
}

func (r *AnalysisRepository) FindByID(_ context.Context, id string) (*service.AnalysisResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if v, ok := r.data[id]; ok {
		return v, nil
	}
	return nil, nil
}

func (r *AnalysisRepository) FindByUserID(_ context.Context, userID string, limit int) ([]*service.AnalysisResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var results []*service.AnalysisResult
	for _, v := range r.data {
		if v.UserID == userID {
			results = append(results, v)
			if limit > 0 && len(results) >= limit {
				break
			}
		}
	}
	return results, nil
}

// ============================================================================
// In-Memory HealthScore Repository
// ============================================================================

type HealthScoreRepository struct {
	mu   sync.RWMutex
	data map[string]*service.HealthScore // keyed by userID (latest only)
}

func NewHealthScoreRepository() *HealthScoreRepository {
	return &HealthScoreRepository{data: make(map[string]*service.HealthScore)}
}

func (r *HealthScoreRepository) Save(_ context.Context, score *service.HealthScore) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[score.UserID] = score
	return nil
}

func (r *HealthScoreRepository) FindLatestByUserID(_ context.Context, userID string) (*service.HealthScore, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if v, ok := r.data[userID]; ok {
		return v, nil
	}
	return nil, nil
}
