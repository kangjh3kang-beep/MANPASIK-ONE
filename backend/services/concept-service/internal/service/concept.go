package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ============================================================================
// 에러 및 상수
// ============================================================================

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalidInput  = errors.New("invalid input")
	ErrDeviceLimit   = errors.New("device limit exceeded")
)

// 허용되는 컨셉 카테고리
var validCategories = map[string]bool{
	"personal":  true,
	"medical":   true,
	"fitness":   true,
	"nutrition": true,
	"research":  true,
}

// MaxDevicesPerConcept은 컨셉 당 최대 디바이스 수입니다.
const MaxDevicesPerConcept = 20

// ============================================================================
// 도메인 모델
// ============================================================================

type Concept struct {
	ID, Name, Description, Category, IconURL, OwnerID string
	DeviceIDs                                         []string
	CreatedAt                                         time.Time
}

type Organization struct {
	ID, Name, Description, OwnerID string
	MemberIDs                      []string
	CreatedAt                      time.Time
}

// ============================================================================
// Repository 인터페이스
// ============================================================================

type ConceptRepository interface {
	List(ctx context.Context) ([]*Concept, error)
	GetByID(ctx context.Context, id string) (*Concept, error)
	Create(ctx context.Context, c *Concept) error
	Update(ctx context.Context, c *Concept) error
	Delete(ctx context.Context, id string) error
}

type OrganizationRepository interface {
	List(ctx context.Context) ([]*Organization, error)
	GetByID(ctx context.Context, id string) (*Organization, error)
	Create(ctx context.Context, o *Organization) error
	Update(ctx context.Context, o *Organization) error
	Delete(ctx context.Context, id string) error
}

// ============================================================================
// ConceptService
// ============================================================================

type ConceptService struct {
	mu       sync.RWMutex
	log      *zap.Logger
	concepts ConceptRepository
	orgs     OrganizationRepository
}

func NewConceptService(log *zap.Logger, cr ConceptRepository, or OrganizationRepository) *ConceptService {
	return &ConceptService{log: log, concepts: cr, orgs: or}
}

func (s *ConceptService) ListConcepts(ctx context.Context) ([]*Concept, error) {
	return s.concepts.List(ctx)
}

func (s *ConceptService) GetConcept(ctx context.Context, id string) (*Concept, error) {
	if id == "" {
		return nil, ErrInvalidInput
	}
	return s.concepts.GetByID(ctx, id)
}

func (s *ConceptService) CreateConcept(ctx context.Context, name, description, category, iconURL, ownerID string) (*Concept, error) {
	if name == "" {
		return nil, ErrInvalidInput
	}
	if len(name) < 2 || len(name) > 100 {
		return nil, fmt.Errorf("%w: 이름은 2~100자여야 합니다", ErrInvalidInput)
	}
	if category != "" && !validCategories[category] {
		return nil, fmt.Errorf("%w: 지원하지 않는 카테고리: %s", ErrInvalidInput, category)
	}
	if category == "" {
		category = "personal"
	}

	c := &Concept{
		ID: uuid.New().String(), Name: name, Description: description, Category: category,
		IconURL: iconURL, OwnerID: ownerID, CreatedAt: time.Now().UTC(),
	}
	return c, s.concepts.Create(ctx, c)
}

// DeleteConcept은 컨셉을 삭제합니다.
func (s *ConceptService) DeleteConcept(ctx context.Context, id string) error {
	if id == "" {
		return ErrInvalidInput
	}
	return s.concepts.Delete(ctx, id)
}

func (s *ConceptService) AssignDevice(ctx context.Context, conceptID, deviceID string) error {
	if conceptID == "" || deviceID == "" {
		return ErrInvalidInput
	}
	c, err := s.concepts.GetByID(ctx, conceptID)
	if err != nil {
		return err
	}
	for _, d := range c.DeviceIDs {
		if d == deviceID {
			return ErrAlreadyExists
		}
	}
	if len(c.DeviceIDs) >= MaxDevicesPerConcept {
		return fmt.Errorf("%w: 최대 %d개 디바이스", ErrDeviceLimit, MaxDevicesPerConcept)
	}
	c.DeviceIDs = append(c.DeviceIDs, deviceID)
	return s.concepts.Update(ctx, c)
}

// UnassignDevice는 컨셉에서 디바이스를 제거합니다.
func (s *ConceptService) UnassignDevice(ctx context.Context, conceptID, deviceID string) error {
	if conceptID == "" || deviceID == "" {
		return ErrInvalidInput
	}
	c, err := s.concepts.GetByID(ctx, conceptID)
	if err != nil {
		return err
	}
	newIDs := make([]string, 0)
	found := false
	for _, d := range c.DeviceIDs {
		if d == deviceID {
			found = true
			continue
		}
		newIDs = append(newIDs, d)
	}
	if !found {
		return ErrNotFound
	}
	c.DeviceIDs = newIDs
	return s.concepts.Update(ctx, c)
}

func (s *ConceptService) GetConceptStats(ctx context.Context, id string) (int32, int32, float64, error) {
	if id == "" {
		return 0, 0, 0, ErrInvalidInput
	}
	c, err := s.concepts.GetByID(ctx, id)
	if err != nil {
		return 0, 0, 0, err
	}
	deviceCount := int32(len(c.DeviceIDs))
	activeCount := deviceCount
	var rate float64
	if deviceCount > 0 {
		rate = float64(activeCount) / float64(MaxDevicesPerConcept) * 100.0
	}
	return deviceCount, activeCount, rate, nil
}

func (s *ConceptService) GetConceptDashboard(ctx context.Context, id string) (int32, int32, float64, []string, error) {
	if id == "" {
		return 0, 0, 0, nil, ErrInvalidInput
	}
	c, err := s.concepts.GetByID(ctx, id)
	if err != nil {
		return 0, 0, 0, nil, err
	}
	deviceCount := int32(len(c.DeviceIDs))
	activeCount := deviceCount
	var rate float64
	if deviceCount > 0 {
		rate = float64(activeCount) / float64(MaxDevicesPerConcept) * 100.0
	}
	return deviceCount, activeCount, rate, c.DeviceIDs, nil
}

// ============================================================================
// Organization 관련
// ============================================================================

func (s *ConceptService) CreateOrganization(ctx context.Context, name, description string) (*Organization, error) {
	if name == "" {
		return nil, ErrInvalidInput
	}
	if len(name) < 2 || len(name) > 100 {
		return nil, fmt.Errorf("%w: 이름은 2~100자여야 합니다", ErrInvalidInput)
	}
	o := &Organization{ID: uuid.New().String(), Name: name, Description: description, CreatedAt: time.Now().UTC()}
	return o, s.orgs.Create(ctx, o)
}

func (s *ConceptService) GetOrganization(ctx context.Context, id string) (*Organization, error) {
	if id == "" {
		return nil, ErrInvalidInput
	}
	return s.orgs.GetByID(ctx, id)
}

func (s *ConceptService) ListOrganizations(ctx context.Context) ([]*Organization, error) {
	return s.orgs.List(ctx)
}

// DeleteOrganization은 조직을 삭제합니다.
func (s *ConceptService) DeleteOrganization(ctx context.Context, id string) error {
	if id == "" {
		return ErrInvalidInput
	}
	return s.orgs.Delete(ctx, id)
}

func (s *ConceptService) AddMember(ctx context.Context, orgID, memberID string) error {
	if orgID == "" || memberID == "" {
		return ErrInvalidInput
	}
	o, err := s.orgs.GetByID(ctx, orgID)
	if err != nil {
		return err
	}
	for _, m := range o.MemberIDs {
		if m == memberID {
			return ErrAlreadyExists
		}
	}
	o.MemberIDs = append(o.MemberIDs, memberID)
	return s.orgs.Update(ctx, o)
}

func (s *ConceptService) RemoveMember(ctx context.Context, orgID, memberID string) error {
	if orgID == "" || memberID == "" {
		return ErrInvalidInput
	}
	o, err := s.orgs.GetByID(ctx, orgID)
	if err != nil {
		return err
	}
	newM := make([]string, 0)
	found := false
	for _, m := range o.MemberIDs {
		if m == memberID {
			found = true
			continue
		}
		newM = append(newM, m)
	}
	if !found {
		return ErrNotFound
	}
	o.MemberIDs = newM
	return s.orgs.Update(ctx, o)
}
