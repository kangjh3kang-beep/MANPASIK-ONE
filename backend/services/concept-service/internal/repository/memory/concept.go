package memory

import (
	"context"
	"sync"
	"github.com/manpasik/backend/services/concept-service/internal/service"
)

type ConceptRepository struct {
	mu   sync.RWMutex
	data map[string]*service.Concept
}
func NewConceptRepository() *ConceptRepository {
	return &ConceptRepository{data: make(map[string]*service.Concept)}
}
func (r *ConceptRepository) List(_ context.Context) ([]*service.Concept, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	out := make([]*service.Concept, 0, len(r.data))
	for _, v := range r.data { out = append(out, v) }
	return out, nil
}
func (r *ConceptRepository) GetByID(_ context.Context, id string) (*service.Concept, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	if v, ok := r.data[id]; ok { return v, nil }
	return nil, service.ErrNotFound
}
func (r *ConceptRepository) Create(_ context.Context, c *service.Concept) error {
	r.mu.Lock(); defer r.mu.Unlock()
	r.data[c.ID] = c; return nil
}
func (r *ConceptRepository) Update(_ context.Context, c *service.Concept) error {
	r.mu.Lock(); defer r.mu.Unlock()
	if _, ok := r.data[c.ID]; !ok { return service.ErrNotFound }
	r.data[c.ID] = c; return nil
}
func (r *ConceptRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock(); defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok { return service.ErrNotFound }
	delete(r.data, id); return nil
}

type OrganizationRepository struct {
	mu   sync.RWMutex
	data map[string]*service.Organization
}
func NewOrganizationRepository() *OrganizationRepository {
	return &OrganizationRepository{data: make(map[string]*service.Organization)}
}
func (r *OrganizationRepository) List(_ context.Context) ([]*service.Organization, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	out := make([]*service.Organization, 0, len(r.data))
	for _, v := range r.data { out = append(out, v) }
	return out, nil
}
func (r *OrganizationRepository) GetByID(_ context.Context, id string) (*service.Organization, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	if v, ok := r.data[id]; ok { return v, nil }
	return nil, service.ErrNotFound
}
func (r *OrganizationRepository) Create(_ context.Context, o *service.Organization) error {
	r.mu.Lock(); defer r.mu.Unlock()
	r.data[o.ID] = o; return nil
}
func (r *OrganizationRepository) Update(_ context.Context, o *service.Organization) error {
	r.mu.Lock(); defer r.mu.Unlock()
	if _, ok := r.data[o.ID]; !ok { return service.ErrNotFound }
	r.data[o.ID] = o; return nil
}
func (r *OrganizationRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock(); defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok { return service.ErrNotFound }
	delete(r.data, id); return nil
}
