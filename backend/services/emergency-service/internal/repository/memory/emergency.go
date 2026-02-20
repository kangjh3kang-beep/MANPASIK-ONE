package memory

import (
	"errors"
	"sync"
	"time"
)

// Emergency represents an emergency report.
type Emergency struct {
	ID          string
	UserID      string
	Type        string
	Location    string
	Description string
	Status      string
	ContactIDs  []string
	Timestamp   time.Time
}

// EmergencyContact represents a user's emergency contact.
type EmergencyContact struct {
	ID           string
	UserID       string
	Name         string
	Phone        string
	Relationship string
}

// EmergencySettings represents a user's emergency settings.
type EmergencySettings struct {
	UserID              string
	AutoCall119         bool
	EmergencyContactIDs []string
	MedicalInfo         string
}

// EmergencyRepository provides in-memory storage for emergency data.
type EmergencyRepository struct {
	mu          sync.RWMutex
	emergencies map[string]*Emergency
	contacts    map[string][]*EmergencyContact
	settings    map[string]*EmergencySettings
	nextID      int
}

// NewEmergencyRepository creates a new in-memory emergency repository.
func NewEmergencyRepository() *EmergencyRepository {
	return &EmergencyRepository{
		emergencies: make(map[string]*Emergency),
		contacts:    make(map[string][]*EmergencyContact),
		settings:    make(map[string]*EmergencySettings),
	}
}

// CreateEmergency stores a new emergency and returns its ID.
func (r *EmergencyRepository) CreateEmergency(e *Emergency) (string, error) {
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
func (r *EmergencyRepository) GetEmergency(id string) (*Emergency, error) {
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

// GetContactsByUser returns all emergency contacts for a user.
func (r *EmergencyRepository) GetContactsByUser(userID string) ([]*EmergencyContact, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	src := r.contacts[userID]
	out := make([]*EmergencyContact, len(src))
	for i, c := range src {
		cp := *c
		out[i] = &cp
	}
	return out, nil
}

// AddContact adds an emergency contact for a user.
func (r *EmergencyRepository) AddContact(c *EmergencyContact) error {
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

// GetSettings retrieves emergency settings for a user.
func (r *EmergencyRepository) GetSettings(userID string) (*EmergencySettings, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s, ok := r.settings[userID]
	if !ok {
		return &EmergencySettings{
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
func (r *EmergencyRepository) SaveSettings(s *EmergencySettings) error {
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
