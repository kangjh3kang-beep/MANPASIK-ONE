package memory

import (
	"errors"
	"sync"
	"time"

	"github.com/manpasik/backend/services/emergency-service/internal/service"
)

// EmergencyRepository provides in-memory storage for emergency data.
type EmergencyRepository struct {
	mu          sync.RWMutex
	emergencies map[string]*service.Emergency
	contacts    map[string][]*service.EmergencyContact
	settings    map[string]*service.EmergencySettings
	nextID      int
}

// NewEmergencyRepository creates a new in-memory emergency repository.
func NewEmergencyRepository() *EmergencyRepository {
	return &EmergencyRepository{
		emergencies: make(map[string]*service.Emergency),
		contacts:    make(map[string][]*service.EmergencyContact),
		settings:    make(map[string]*service.EmergencySettings),
	}
}

// CreateEmergency stores a new emergency and returns its ID.
func (r *EmergencyRepository) CreateEmergency(e *service.Emergency) (string, error) {
	if e == nil {
		return "", errors.New("emergency must not be nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nextID++
	e.ID = idFromInt(r.nextID)
	e.Timestamp = time.Now()
	if e.Status == "" {
		e.Status = "reported"
	}

	stored := *e
	stored.ContactIDs = copyStrings(e.ContactIDs)
	r.emergencies[e.ID] = &stored
	return e.ID, nil
}

// GetEmergency retrieves an emergency by ID.
func (r *EmergencyRepository) GetEmergency(id string) (*service.Emergency, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	e, ok := r.emergencies[id]
	if !ok {
		return nil, errors.New("emergency not found")
	}
	out := *e
	out.ContactIDs = copyStrings(e.ContactIDs)
	return &out, nil
}

// UpdateEmergency updates an existing emergency record.
func (r *EmergencyRepository) UpdateEmergency(e *service.Emergency) error {
	if e == nil {
		return errors.New("emergency must not be nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.emergencies[e.ID]; !ok {
		return errors.New("emergency not found")
	}
	stored := *e
	stored.ContactIDs = copyStrings(e.ContactIDs)
	r.emergencies[e.ID] = &stored
	return nil
}

// ListEmergenciesByUser returns all emergencies for a user.
func (r *EmergencyRepository) ListEmergenciesByUser(userID string) ([]*service.Emergency, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*service.Emergency
	for _, e := range r.emergencies {
		if e.UserID == userID {
			out := *e
			out.ContactIDs = copyStrings(e.ContactIDs)
			result = append(result, &out)
		}
	}
	return result, nil
}

// GetContactsByUser returns all emergency contacts for a user.
func (r *EmergencyRepository) GetContactsByUser(userID string) ([]*service.EmergencyContact, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	src := r.contacts[userID]
	out := make([]*service.EmergencyContact, len(src))
	for i, c := range src {
		cp := *c
		out[i] = &cp
	}
	return out, nil
}

// AddContact adds an emergency contact for a user.
func (r *EmergencyRepository) AddContact(c *service.EmergencyContact) error {
	if c == nil {
		return errors.New("contact must not be nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nextID++
	c.ID = idFromInt(r.nextID)
	stored := *c
	r.contacts[c.UserID] = append(r.contacts[c.UserID], &stored)
	return nil
}

// RemoveContact removes an emergency contact by ID.
func (r *EmergencyRepository) RemoveContact(userID, contactID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	contacts := r.contacts[userID]
	found := false
	newList := make([]*service.EmergencyContact, 0, len(contacts))
	for _, c := range contacts {
		if c.ID == contactID {
			found = true
			continue
		}
		newList = append(newList, c)
	}
	if !found {
		return errors.New("contact not found")
	}
	r.contacts[userID] = newList
	return nil
}

// GetSettings retrieves emergency settings for a user.
func (r *EmergencyRepository) GetSettings(userID string) (*service.EmergencySettings, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s, ok := r.settings[userID]
	if !ok {
		return &service.EmergencySettings{
			UserID:              userID,
			AutoCall119:         false,
			EmergencyContactIDs: []string{},
			MedicalInfo:         "",
		}, nil
	}
	out := *s
	out.EmergencyContactIDs = copyStrings(s.EmergencyContactIDs)
	return &out, nil
}

// SaveSettings creates or updates emergency settings for a user.
func (r *EmergencyRepository) SaveSettings(s *service.EmergencySettings) error {
	if s == nil {
		return errors.New("settings must not be nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	stored := *s
	stored.EmergencyContactIDs = copyStrings(s.EmergencyContactIDs)
	r.settings[s.UserID] = &stored
	return nil
}

func idFromInt(n int) string {
	s := ""
	if n == 0 {
		return "emg-0"
	}
	for n > 0 {
		s = string(rune(48+n%10)) + s
		n /= 10
	}
	return "emg-" + s
}

func copyStrings(src []string) []string {
	if src == nil {
		return nil
	}
	out := make([]string, len(src))
	copy(out, src)
	return out
}
