package service

import (
	"errors"

	"github.com/manpasik/backend/services/emergency-service/internal/repository/memory"
)

// EmergencyService provides emergency-related business logic.
type EmergencyService struct {
	repo *memory.EmergencyRepository
}

// NewEmergencyService creates a new EmergencyService backed by the given repository.
func NewEmergencyService(repo *memory.EmergencyRepository) *EmergencyService {
	return &EmergencyService{repo: repo}
}

// ReportEmergencyInput holds the parameters for reporting an emergency.
type ReportEmergencyInput struct {
	UserID      string
	Type        string
	Location    string
	Description string
}

// ReportEmergency creates a new emergency report and returns its ID.
func (s *EmergencyService) ReportEmergency(in ReportEmergencyInput) (string, error) {
	if in.UserID == "" {
		return "", errors.New("user_id is required")
	}
	if in.Type == "" {
		return "", errors.New("type is required")
	}

	e := &memory.Emergency{
		UserID:      in.UserID,
		Type:        in.Type,
		Location:    in.Location,
		Description: in.Description,
		Status:      "reported",
	}

	// Attach user's emergency contacts if available.
	contacts, _ := s.repo.GetContactsByUser(in.UserID)
	ids := make([]string, 0, len(contacts))
	for _, c := range contacts {
		ids = append(ids, c.ID)
	}
	e.ContactIDs = ids

	return s.repo.CreateEmergency(e)
}

// GetEmergencyContacts returns the emergency contacts for a user.
func (s *EmergencyService) GetEmergencyContacts(userID string) ([]*memory.EmergencyContact, error) {
	if userID == "" {
		return nil, errors.New("user_id is required")
	}
	return s.repo.GetContactsByUser(userID)
}

// UpdateEmergencySettings saves emergency settings for a user.
func (s *EmergencyService) UpdateEmergencySettings(settings *memory.EmergencySettings) error {
	if settings == nil {
		return errors.New("settings must not be nil")
	}
	if settings.UserID == "" {
		return errors.New("user_id is required")
	}
	return s.repo.SaveSettings(settings)
}

// GetEmergencySettings retrieves emergency settings for a user.
func (s *EmergencyService) GetEmergencySettings(userID string) (*memory.EmergencySettings, error) {
	if userID == "" {
		return nil, errors.New("user_id is required")
	}
	return s.repo.GetSettings(userID)
}
