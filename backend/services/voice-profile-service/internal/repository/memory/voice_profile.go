package memory

import (
	"context"
	"sync"
	"github.com/manpasik/backend/services/voice-profile-service/internal/service"
)

type VoiceProfileRepository struct {
	mu   sync.RWMutex
	data map[string]*service.VoiceProfile
}

func NewVoiceProfileRepository() *VoiceProfileRepository {
	return &VoiceProfileRepository{data: make(map[string]*service.VoiceProfile)}
}
func (r *VoiceProfileRepository) Create(_ context.Context, p *service.VoiceProfile) error {
	r.mu.Lock(); defer r.mu.Unlock()
	r.data[p.ProfileID] = p; return nil
}
func (r *VoiceProfileRepository) GetByID(_ context.Context, id string) (*service.VoiceProfile, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	if v, ok := r.data[id]; ok { return v, nil }
	return nil, service.ErrNotFound
}
func (r *VoiceProfileRepository) ListByUser(_ context.Context, userID string) ([]*service.VoiceProfile, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	out := make([]*service.VoiceProfile, 0)
	for _, v := range r.data { if v.UserID == userID { out = append(out, v) } }
	return out, nil
}
func (r *VoiceProfileRepository) Update(_ context.Context, p *service.VoiceProfile) error {
	r.mu.Lock(); defer r.mu.Unlock()
	if _, ok := r.data[p.ProfileID]; !ok { return service.ErrNotFound }
	r.data[p.ProfileID] = p; return nil
}
func (r *VoiceProfileRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock(); defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok { return service.ErrNotFound }
	delete(r.data, id); return nil
}
