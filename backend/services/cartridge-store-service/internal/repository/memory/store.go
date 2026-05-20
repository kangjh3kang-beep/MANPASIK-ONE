package memory

import (
	"context"
	"strings"
	"sync"
	"github.com/manpasik/backend/services/cartridge-store-service/internal/service"
)

type StoreRepository struct {
	mu          sync.RWMutex
	developers  map[string]*service.Developer
	devByUser   map[string]string
	apiKeys     map[string]*service.ApiKey
	items       map[string]*service.StoreItem
	purchases   []*service.Purchase
	reviews     map[string]*service.ReviewStatus
}

func NewStoreRepository() *StoreRepository {
	return &StoreRepository{
		developers: make(map[string]*service.Developer),
		devByUser:  make(map[string]string),
		apiKeys:    make(map[string]*service.ApiKey),
		items:      make(map[string]*service.StoreItem),
		reviews:    make(map[string]*service.ReviewStatus),
	}
}
func (r *StoreRepository) CreateDeveloper(_ context.Context, d *service.Developer) error {
	r.mu.Lock(); defer r.mu.Unlock()
	r.developers[d.ID] = d; r.devByUser[d.UserID] = d.ID; return nil
}
func (r *StoreRepository) GetDeveloper(_ context.Context, userID string) (*service.Developer, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	if id, ok := r.devByUser[userID]; ok { return r.developers[id], nil }
	return nil, service.ErrNotFound
}
func (r *StoreRepository) CreateApiKey(_ context.Context, k *service.ApiKey) error {
	r.mu.Lock(); defer r.mu.Unlock()
	r.apiKeys[k.KeyID] = k; return nil
}
func (r *StoreRepository) CreateItem(_ context.Context, item *service.StoreItem) error {
	r.mu.Lock(); defer r.mu.Unlock()
	r.items[item.ID] = item; return nil
}
func (r *StoreRepository) ListItems(_ context.Context, category string, limit, offset int32) ([]*service.StoreItem, int32, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	out := make([]*service.StoreItem, 0)
	for _, item := range r.items {
		if category != "" && item.Category != category { continue }
		out = append(out, item)
	}
	total := int32(len(out))
	s := int(offset); if s > len(out) { s = len(out) }
	e := s + int(limit); if e > len(out) { e = len(out) }
	return out[s:e], total, nil
}
func (r *StoreRepository) SearchItems(_ context.Context, query string, limit int32) ([]*service.StoreItem, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	out := make([]*service.StoreItem, 0)
	q := strings.ToLower(query)
	for _, item := range r.items {
		if strings.Contains(strings.ToLower(item.Name), q) || strings.Contains(strings.ToLower(item.Description), q) {
			out = append(out, item)
			if int32(len(out)) >= limit { break }
		}
	}
	return out, nil
}
func (r *StoreRepository) GetItem(_ context.Context, id string) (*service.StoreItem, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	if v, ok := r.items[id]; ok { return v, nil }
	return nil, service.ErrNotFound
}
func (r *StoreRepository) CreatePurchase(_ context.Context, p *service.Purchase) error {
	r.mu.Lock(); defer r.mu.Unlock()
	r.purchases = append(r.purchases, p); return nil
}
func (r *StoreRepository) ListPurchases(_ context.Context, userID string, limit int32) ([]*service.Purchase, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	out := make([]*service.Purchase, 0)
	for _, p := range r.purchases {
		if p.UserID == userID { out = append(out, p) }
		if int32(len(out)) >= limit { break }
	}
	return out, nil
}
func (r *StoreRepository) CreateReview(_ context.Context, rv *service.ReviewStatus) error {
	r.mu.Lock(); defer r.mu.Unlock()
	r.reviews[rv.SubmissionID] = rv; return nil
}
func (r *StoreRepository) GetReview(_ context.Context, id string) (*service.ReviewStatus, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	if v, ok := r.reviews[id]; ok { return v, nil }
	return nil, service.ErrNotFound
}
func (r *StoreRepository) UpdateReview(_ context.Context, rv *service.ReviewStatus) error {
	r.mu.Lock(); defer r.mu.Unlock()
	r.reviews[rv.SubmissionID] = rv; return nil
}
func (r *StoreRepository) ListSubmissions(_ context.Context, devID string, limit int32) ([]*service.ReviewStatus, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	out := make([]*service.ReviewStatus, 0)
	for _, rv := range r.reviews { out = append(out, rv) }
	return out, nil
}
