package memory

import (
	"context"
	"sync"
	"github.com/google/uuid"
	"github.com/manpasik/backend/services/data-platform-service/internal/service"
)

type PlatformRepository struct {
	mu       sync.RWMutex
	datasets map[string]*service.AnonymizedDataset
	consents map[string]*service.ConsentInfo
	insights []*service.AggregatedInsight
}

func NewPlatformRepository() *PlatformRepository {
	return &PlatformRepository{
		datasets: make(map[string]*service.AnonymizedDataset),
		consents: make(map[string]*service.ConsentInfo),
	}
}

func (r *PlatformRepository) GetRegionStats(_ context.Context, _, _, _ string) (*service.RegionStats, error) { return nil, nil }
func (r *PlatformRepository) ListHotspots(_ context.Context, _, _ string, _ int32) ([]*service.RegionStats, error) { return nil, nil }
func (r *PlatformRepository) GetTrend(_ context.Context, _, _ string) (*service.RegionTrend, error) { return nil, nil }

func (r *PlatformRepository) ListDatasets(_ context.Context, dsType string, limit, offset int32) ([]*service.AnonymizedDataset, int32, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	out := make([]*service.AnonymizedDataset, 0)
	for _, ds := range r.datasets {
		if dsType != "" && ds.DatasetType != dsType { continue }
		out = append(out, ds)
	}
	total := int32(len(out))
	start := int(offset); if start > len(out) { start = len(out) }
	end := start + int(limit); if end > len(out) { end = len(out) }
	return out[start:end], total, nil
}
func (r *PlatformRepository) GetDataset(_ context.Context, id string) (*service.AnonymizedDataset, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	if v, ok := r.datasets[id]; ok { return v, nil }
	return nil, service.ErrNotFound
}
func (r *PlatformRepository) SaveDataset(_ context.Context, ds *service.AnonymizedDataset) error {
	r.mu.Lock(); defer r.mu.Unlock()
	r.datasets[ds.ID] = ds; return nil
}
func (r *PlatformRepository) CreateAccessRequest(_ context.Context, requesterID, datasetID, purpose string) (*service.DataAccessResp, error) {
	return &service.DataAccessResp{RequestID: uuid.New().String(), Status: "pending"}, nil
}
func (r *PlatformRepository) GetConsent(_ context.Context, userID string) (*service.ConsentInfo, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	if v, ok := r.consents[userID]; ok { return v, nil }
	return nil, nil
}
func (r *PlatformRepository) UpdateConsent(_ context.Context, c *service.ConsentInfo) error {
	r.mu.Lock(); defer r.mu.Unlock()
	r.consents[c.UserID] = c; return nil
}
func (r *PlatformRepository) ListInsights(_ context.Context, biomarker string, limit int32) ([]*service.AggregatedInsight, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	out := make([]*service.AggregatedInsight, 0)
	for _, i := range r.insights {
		if biomarker != "" && i.Biomarker != biomarker { continue }
		out = append(out, i)
		if int32(len(out)) >= limit { break }
	}
	return out, nil
}
func (r *PlatformRepository) SaveInsight(_ context.Context, i *service.AggregatedInsight) error {
	r.mu.Lock(); defer r.mu.Unlock()
	r.insights = append(r.insights, i); return nil
}
