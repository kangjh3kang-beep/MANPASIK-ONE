package service

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/manpasik/backend/services/emergency-service/internal/notifier"
)

// ============================================================================
// 도메인 상수 및 타입
// ============================================================================

// EmergencyType은 응급 상황 유형입니다.
type EmergencyType string

const (
	EmergencyTypeCardiac            EmergencyType = "cardiac"
	EmergencyTypeFall               EmergencyType = "fall"
	EmergencyTypeHypoglycemia       EmergencyType = "hypoglycemia"
	EmergencyTypeHypertensiveCrisis EmergencyType = "hypertensive_crisis"
	EmergencyTypeSevereAllergy      EmergencyType = "severe_allergy"
	EmergencyTypeOther              EmergencyType = "other"
)

// EmergencySeverity는 응급 상황 심각도입니다.
type EmergencySeverity int

const (
	SeverityLow      EmergencySeverity = 1
	SeverityModerate EmergencySeverity = 2
	SeverityHigh     EmergencySeverity = 3
	SeverityCritical EmergencySeverity = 4
)

// EmergencyStatus는 응급 상황 처리 상태입니다.
type EmergencyStatus string

const (
	StatusReported   EmergencyStatus = "reported"
	StatusDispatched EmergencyStatus = "dispatched"
	StatusResponding EmergencyStatus = "responding"
	StatusResolved   EmergencyStatus = "resolved"
	StatusCancelled  EmergencyStatus = "cancelled"
)

// severityKeywords는 설명 기반 심각도 상향 키워드입니다.
var severityKeywords = map[string]EmergencySeverity{
	"chest pain":  SeverityCritical,
	"흉통":         SeverityCritical,
	"breathing":   SeverityHigh,
	"호흡곤란":       SeverityHigh,
	"unconscious": SeverityCritical,
	"의식":         SeverityHigh,
	"bleeding":    SeverityHigh,
	"출혈":         SeverityHigh,
	"seizure":     SeverityCritical,
	"경련":         SeverityCritical,
}

// ============================================================================
// 도메인 모델
// ============================================================================

// Emergency represents an emergency report.
type Emergency struct {
	ID          string
	UserID      string
	Type        string
	Location    string
	Description string
	Status      string
	Severity    EmergencySeverity
	ContactIDs  []string
	Resolution  string
	Timestamp   time.Time
	ResolvedAt  *time.Time
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
	UserID              string   `json:"user_id"`
	AutoCall119         bool     `json:"auto_call_119"`
	EmergencyContactIDs []string `json:"emergency_contact_ids"`
	MedicalInfo         string   `json:"medical_info"`
}

// ============================================================================
// Repository 인터페이스
// ============================================================================

// EmergencyRepository defines the storage interface for emergency data (H1 모듈독립성).
type EmergencyRepository interface {
	CreateEmergency(e *Emergency) (string, error)
	GetEmergency(id string) (*Emergency, error)
	UpdateEmergency(e *Emergency) error
	ListEmergenciesByUser(userID string) ([]*Emergency, error)
	GetContactsByUser(userID string) ([]*EmergencyContact, error)
	AddContact(c *EmergencyContact) error
	RemoveContact(userID, contactID string) error
	GetSettings(userID string) (*EmergencySettings, error)
	SaveSettings(s *EmergencySettings) error
}

// ============================================================================
// 에러 정의
// ============================================================================

var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrNotFound          = errors.New("not found")
	ErrAlreadyResolved   = errors.New("emergency already resolved")
	ErrInvalidTransition = errors.New("invalid status transition")
)

// ============================================================================
// EmergencyService
// ============================================================================

// EmergencyService provides emergency-related business logic.
type EmergencyService struct {
	repo             EmergencyRepository
	emergencyNotifier notifier.EmergencyNotifier
}

// NewEmergencyService creates a new EmergencyService backed by the given repository.
func NewEmergencyService(repo EmergencyRepository) *EmergencyService {
	return &EmergencyService{repo: repo}
}

// SetEmergencyNotifier는 외부 응급 알림 제공자를 설정합니다.
func (s *EmergencyService) SetEmergencyNotifier(n notifier.EmergencyNotifier) {
	s.emergencyNotifier = n
}

// ReportEmergencyInput holds the parameters for reporting an emergency.
type ReportEmergencyInput struct {
	UserID      string
	Type        string
	Location    string
	Description string
}

// classifySeverity는 응급 유형과 설명으로 심각도를 자동 분류합니다.
func classifySeverity(emergencyType, description string) EmergencySeverity {
	var baseSeverity EmergencySeverity
	switch EmergencyType(emergencyType) {
	case EmergencyTypeCardiac, EmergencyTypeSevereAllergy:
		baseSeverity = SeverityCritical
	case EmergencyTypeFall:
		baseSeverity = SeverityHigh
	case EmergencyTypeHypoglycemia, EmergencyTypeHypertensiveCrisis:
		baseSeverity = SeverityModerate
	default:
		baseSeverity = SeverityLow
	}

	lower := strings.ToLower(description)
	for keyword, severity := range severityKeywords {
		if strings.Contains(lower, keyword) && severity > baseSeverity {
			baseSeverity = severity
		}
	}

	return baseSeverity
}

// ReportEmergency creates a new emergency report and returns its ID.
func (s *EmergencyService) ReportEmergency(in ReportEmergencyInput) (string, error) {
	if in.UserID == "" {
		return "", errors.New("user_id is required")
	}
	if in.Type == "" {
		return "", errors.New("type is required")
	}

	severity := classifySeverity(in.Type, in.Description)

	e := &Emergency{
		UserID:      in.UserID,
		Type:        in.Type,
		Location:    in.Location,
		Description: in.Description,
		Status:      string(StatusReported),
		Severity:    severity,
	}

	// Attach user's emergency contacts if available.
	contacts, _ := s.repo.GetContactsByUser(in.UserID)
	ids := make([]string, 0, len(contacts))
	for _, c := range contacts {
		ids = append(ids, c.ID)
	}
	e.ContactIDs = ids

	// AutoCall119 설정 + 심각도 High 이상이면 자동 dispatched
	settings, _ := s.repo.GetSettings(in.UserID)
	if settings != nil && settings.AutoCall119 && severity >= SeverityHigh {
		e.Status = string(StatusDispatched)
	}

	emergencyID, err := s.repo.CreateEmergency(e)
	if err != nil {
		return "", err
	}

	// 외부 응급 알림 제공자를 통해 119 신고 + 연락처 알림
	if s.emergencyNotifier != nil && severity >= SeverityHigh {
		s.dispatchExternalAlert(emergencyID, e, settings)
	}

	return emergencyID, nil
}

// dispatchExternalAlert는 외부 응급 알림 서비스에 신고를 전달합니다.
func (s *EmergencyService) dispatchExternalAlert(emergencyID string, e *Emergency, settings *EmergencySettings) {
	ctx := context.Background()

	alert := &notifier.EmergencyAlert{
		EmergencyID: emergencyID,
		UserID:      e.UserID,
		Type:        e.Type,
		Severity:    int(e.Severity),
		Location:    e.Location,
		Description: e.Description,
	}
	if settings != nil {
		alert.MedicalInfo = settings.MedicalInfo
	}

	// 119 자동 신고
	if settings != nil && settings.AutoCall119 {
		result, err := s.emergencyNotifier.Dispatch119(ctx, alert)
		if err != nil {
			log.Printf("[Emergency] 119 신고 실패 (%s): %v", s.emergencyNotifier.ProviderName(), err)
		} else if result != nil {
			log.Printf("[Emergency] 119 신고 완료: dispatch_id=%s, eta=%d분", result.DispatchID, result.EstimatedETA)
		}
	}

	// 비상 연락처 알림
	contacts, _ := s.repo.GetContactsByUser(e.UserID)
	if len(contacts) > 0 {
		phones := make([]string, 0, len(contacts))
		for _, c := range contacts {
			phones = append(phones, c.Phone)
		}
		if err := s.emergencyNotifier.NotifyContacts(ctx, alert, phones); err != nil {
			log.Printf("[Emergency] 비상 연락처 알림 실패 (%s): %v", s.emergencyNotifier.ProviderName(), err)
		}
	}
}

// GetEmergencyContacts returns the emergency contacts for a user.
func (s *EmergencyService) GetEmergencyContacts(userID string) ([]*EmergencyContact, error) {
	if userID == "" {
		return nil, errors.New("user_id is required")
	}
	return s.repo.GetContactsByUser(userID)
}

// AddContact는 사용자에게 비상 연락처를 추가합니다.
func (s *EmergencyService) AddContact(userID string, contact *EmergencyContact) error {
	if userID == "" {
		return errors.New("user_id is required")
	}
	if contact == nil {
		return errors.New("contact must not be nil")
	}
	if contact.Name == "" {
		return errors.New("contact name is required")
	}
	if contact.Phone == "" {
		return errors.New("contact phone is required")
	}
	contact.UserID = userID
	return s.repo.AddContact(contact)
}

// RemoveContact는 사용자의 비상 연락처를 제거합니다.
func (s *EmergencyService) RemoveContact(userID, contactID string) error {
	if userID == "" {
		return errors.New("user_id is required")
	}
	if contactID == "" {
		return errors.New("contact_id is required")
	}
	return s.repo.RemoveContact(userID, contactID)
}

// ResolveEmergency는 응급 상황을 해결 처리합니다.
func (s *EmergencyService) ResolveEmergency(emergencyID, resolution string) error {
	if emergencyID == "" {
		return errors.New("emergency_id is required")
	}

	e, err := s.repo.GetEmergency(emergencyID)
	if err != nil {
		return err
	}
	if e == nil {
		return ErrNotFound
	}

	status := EmergencyStatus(e.Status)
	if status == StatusResolved {
		return ErrAlreadyResolved
	}
	if status == StatusCancelled {
		return ErrInvalidTransition
	}

	now := time.Now()
	e.Status = string(StatusResolved)
	e.Resolution = resolution
	e.ResolvedAt = &now
	return s.repo.UpdateEmergency(e)
}

// GetEmergencyHistory는 사용자의 응급 이력을 조회합니다.
func (s *EmergencyService) GetEmergencyHistory(userID string) ([]*Emergency, error) {
	if userID == "" {
		return nil, errors.New("user_id is required")
	}
	return s.repo.ListEmergenciesByUser(userID)
}

// UpdateEmergencySettings saves emergency settings for a user.
func (s *EmergencyService) UpdateEmergencySettings(settings *EmergencySettings) error {
	if settings == nil {
		return errors.New("settings must not be nil")
	}
	if settings.UserID == "" {
		return errors.New("user_id is required")
	}
	return s.repo.SaveSettings(settings)
}

// GetEmergencySettings retrieves emergency settings for a user.
func (s *EmergencyService) GetEmergencySettings(userID string) (*EmergencySettings, error) {
	if userID == "" {
		return nil, errors.New("user_id is required")
	}
	return s.repo.GetSettings(userID)
}
