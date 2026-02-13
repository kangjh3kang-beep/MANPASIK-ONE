// Package memory는 telemedicine-service의 인메모리 저장소입니다.
package memory

import (
	"context"
	"sync"

	"github.com/manpasik/backend/services/telemedicine-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
)

// ConsultationRepository는 상담 인메모리 저장소입니다.
type ConsultationRepository struct {
	mu    sync.RWMutex
	store map[string]*service.Consultation
}

// NewConsultationRepository는 새 인메모리 상담 저장소를 생성합니다.
func NewConsultationRepository() *ConsultationRepository {
	return &ConsultationRepository{
		store: make(map[string]*service.Consultation),
	}
}

// Save는 상담을 저장합니다.
func (r *ConsultationRepository) Save(_ context.Context, c *service.Consultation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[c.ID] = c
	return nil
}

// FindByID는 ID로 상담을 조회합니다.
func (r *ConsultationRepository) FindByID(_ context.Context, id string) (*service.Consultation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.store[id]
	if !ok {
		return nil, apperrors.New(apperrors.ErrNotFound, "상담을 찾을 수 없습니다")
	}
	return c, nil
}

// FindByUserID는 사용자 상담 목록을 조회합니다.
func (r *ConsultationRepository) FindByUserID(_ context.Context, userID string, statusFilter service.ConsultationStatus, limit, offset int) ([]*service.Consultation, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*service.Consultation
	for _, c := range r.store {
		if c.PatientUserID != userID {
			continue
		}
		if statusFilter != service.StatusUnknown && c.Status != statusFilter {
			continue
		}
		filtered = append(filtered, c)
	}

	// 시간 역순 정렬 (최신 먼저)
	for i := 0; i < len(filtered); i++ {
		for j := i + 1; j < len(filtered); j++ {
			if filtered[j].CreatedAt.After(filtered[i].CreatedAt) {
				filtered[i], filtered[j] = filtered[j], filtered[i]
			}
		}
	}

	total := len(filtered)
	if offset >= total {
		return nil, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}

	return filtered[offset:end], total, nil
}

// Update는 상담을 업데이트합니다.
func (r *ConsultationRepository) Update(_ context.Context, c *service.Consultation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[c.ID]; !ok {
		return apperrors.New(apperrors.ErrNotFound, "상담을 찾을 수 없습니다")
	}
	r.store[c.ID] = c
	return nil
}

// DoctorRepository는 의사 인메모리 저장소입니다.
type DoctorRepository struct {
	mu    sync.RWMutex
	store map[string]*service.DoctorProfile
}

// NewDoctorRepository는 샘플 의사 데이터가 포함된 인메모리 저장소를 생성합니다.
func NewDoctorRepository() *DoctorRepository {
	repo := &DoctorRepository{
		store: make(map[string]*service.DoctorProfile),
	}

	// 샘플 의사 5명 시드
	doctors := []*service.DoctorProfile{
		{
			ID:                 "doc-001",
			Name:               "김내과",
			Specialty:          service.SpecialtyInternal,
			Hospital:           "서울대학교병원",
			LicenseNumber:      "MED-2020-001",
			ExperienceYears:    15,
			Rating:             4.8,
			TotalConsultations: 1200,
			IsAvailable:        true,
			Languages:          []string{"ko", "en"},
			ProfileImageURL:    "https://img.manpasik.com/doctors/doc-001.jpg",
		},
		{
			ID:                 "doc-002",
			Name:               "이심장",
			Specialty:          service.SpecialtyCardiology,
			Hospital:           "세브란스병원",
			LicenseNumber:      "MED-2018-042",
			ExperienceYears:    20,
			Rating:             4.9,
			TotalConsultations: 2500,
			IsAvailable:        true,
			Languages:          []string{"ko", "en", "ja"},
			ProfileImageURL:    "https://img.manpasik.com/doctors/doc-002.jpg",
		},
		{
			ID:                 "doc-003",
			Name:               "박피부",
			Specialty:          service.SpecialtyDermatology,
			Hospital:           "삼성서울병원",
			LicenseNumber:      "MED-2019-108",
			ExperienceYears:    10,
			Rating:             4.7,
			TotalConsultations: 800,
			IsAvailable:        true,
			Languages:          []string{"ko"},
			ProfileImageURL:    "https://img.manpasik.com/doctors/doc-003.jpg",
		},
		{
			ID:                 "doc-004",
			Name:               "정소아",
			Specialty:          service.SpecialtyPediatrics,
			Hospital:           "서울아산병원",
			LicenseNumber:      "MED-2017-055",
			ExperienceYears:    12,
			Rating:             4.6,
			TotalConsultations: 950,
			IsAvailable:        false,
			Languages:          []string{"ko", "en"},
			ProfileImageURL:    "https://img.manpasik.com/doctors/doc-004.jpg",
		},
		{
			ID:                 "doc-005",
			Name:               "최일반",
			Specialty:          service.SpecialtyGeneral,
			Hospital:           "고려대학교 안암병원",
			LicenseNumber:      "MED-2021-200",
			ExperienceYears:    8,
			Rating:             4.5,
			TotalConsultations: 600,
			IsAvailable:        true,
			Languages:          []string{"ko", "en"},
			ProfileImageURL:    "https://img.manpasik.com/doctors/doc-005.jpg",
		},
	}

	for _, d := range doctors {
		repo.store[d.ID] = d
	}

	return repo
}

// FindBySpecialty는 전문 분야별 의사 목록을 조회합니다.
func (r *DoctorRepository) FindBySpecialty(_ context.Context, specialty service.DoctorSpecialty, language string) ([]*service.DoctorProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*service.DoctorProfile
	for _, d := range r.store {
		if specialty != service.SpecialtyUnknown && d.Specialty != specialty {
			continue
		}
		if language != "" {
			found := false
			for _, lang := range d.Languages {
				if lang == language {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		result = append(result, d)
	}

	return result, nil
}

// FindByID는 ID로 의사를 조회합니다.
func (r *DoctorRepository) FindByID(_ context.Context, id string) (*service.DoctorProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.store[id]
	if !ok {
		return nil, apperrors.New(apperrors.ErrNotFound, "의사를 찾을 수 없습니다")
	}
	return d, nil
}

// VideoSessionRepository는 비디오 세션 인메모리 저장소입니다.
type VideoSessionRepository struct {
	mu    sync.RWMutex
	store map[string]*service.VideoSession
}

// NewVideoSessionRepository는 새 인메모리 비디오 세션 저장소를 생성합니다.
func NewVideoSessionRepository() *VideoSessionRepository {
	return &VideoSessionRepository{
		store: make(map[string]*service.VideoSession),
	}
}

// Save는 비디오 세션을 저장합니다.
func (r *VideoSessionRepository) Save(_ context.Context, s *service.VideoSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[s.ID] = s
	return nil
}

// FindByID는 ID로 비디오 세션을 조회합니다.
func (r *VideoSessionRepository) FindByID(_ context.Context, id string) (*service.VideoSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.store[id]
	if !ok {
		return nil, apperrors.New(apperrors.ErrNotFound, "비디오 세션을 찾을 수 없습니다")
	}
	return s, nil
}

// FindByConsultationID는 상담 ID로 비디오 세션을 조회합니다.
func (r *VideoSessionRepository) FindByConsultationID(_ context.Context, consultationID string) (*service.VideoSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, s := range r.store {
		if s.ConsultationID == consultationID {
			return s, nil
		}
	}
	return nil, apperrors.New(apperrors.ErrNotFound, "비디오 세션을 찾을 수 없습니다")
}

// Update는 비디오 세션을 업데이트합니다.
func (r *VideoSessionRepository) Update(_ context.Context, s *service.VideoSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[s.ID]; !ok {
		return apperrors.New(apperrors.ErrNotFound, "비디오 세션을 찾을 수 없습니다")
	}
	r.store[s.ID] = s
	return nil
}
