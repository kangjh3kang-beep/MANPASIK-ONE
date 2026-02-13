// Package service는 telemedicine-service의 비즈니스 로직을 구현합니다.
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/manpasik/backend/shared/errors"
	"go.uber.org/zap"
)

// ConsultationStatus는 상담 상태입니다.
type ConsultationStatus int

const (
	StatusUnknown    ConsultationStatus = iota
	StatusRequested
	StatusMatched
	StatusScheduled
	StatusInProgress
	StatusCompleted
	StatusCancelled
	StatusNoShow
)

// DoctorSpecialty는 의사 전문 분야입니다.
type DoctorSpecialty int

const (
	SpecialtyUnknown       DoctorSpecialty = iota
	SpecialtyGeneral
	SpecialtyInternal
	SpecialtyCardiology
	SpecialtyEndocrinology
	SpecialtyDermatology
	SpecialtyPediatrics
	SpecialtyPsychiatry
	SpecialtyOrthopedics
	SpecialtyOphthalmology
	SpecialtyENT
)

// VideoSessionStatus는 비디오 세션 상태입니다.
type VideoSessionStatus int

const (
	SessionUnknown   VideoSessionStatus = iota
	SessionWaiting
	SessionConnected
	SessionEnded
	SessionFailed
)

// Consultation은 원격진료 상담 도메인 객체입니다.
type Consultation struct {
	ID              string
	PatientUserID   string
	DoctorID        string
	Specialty       DoctorSpecialty
	ChiefComplaint  string
	Description     string
	Status          ConsultationStatus
	Diagnosis       string
	DoctorNotes     string
	PrescriptionID  string
	DurationMinutes int
	Rating          float64
	ScheduledAt     time.Time
	StartedAt       time.Time
	EndedAt         time.Time
	CreatedAt       time.Time
}

// DoctorProfile은 의사 프로필 도메인 객체입니다.
type DoctorProfile struct {
	ID                 string
	Name               string
	Specialty          DoctorSpecialty
	Hospital           string
	LicenseNumber      string
	ExperienceYears    int
	Rating             float64
	TotalConsultations int
	IsAvailable        bool
	Languages          []string
	ProfileImageURL    string
}

// VideoSession은 비디오 세션 도메인 객체입니다.
type VideoSession struct {
	ID              string
	ConsultationID  string
	RoomURL         string
	Token           string
	Status          VideoSessionStatus
	StartedAt       time.Time
	EndedAt         time.Time
	DurationSeconds int
}

// ConsultationRepository는 상담 저장소 인터페이스입니다.
type ConsultationRepository interface {
	Save(ctx context.Context, c *Consultation) error
	FindByID(ctx context.Context, id string) (*Consultation, error)
	FindByUserID(ctx context.Context, userID string, statusFilter ConsultationStatus, limit, offset int) ([]*Consultation, int, error)
	Update(ctx context.Context, c *Consultation) error
}

// DoctorRepository는 의사 저장소 인터페이스입니다.
type DoctorRepository interface {
	FindBySpecialty(ctx context.Context, specialty DoctorSpecialty, language string) ([]*DoctorProfile, error)
	FindByID(ctx context.Context, id string) (*DoctorProfile, error)
}

// VideoSessionRepository는 비디오 세션 저장소 인터페이스입니다.
type VideoSessionRepository interface {
	Save(ctx context.Context, s *VideoSession) error
	FindByID(ctx context.Context, id string) (*VideoSession, error)
	FindByConsultationID(ctx context.Context, consultationID string) (*VideoSession, error)
	Update(ctx context.Context, s *VideoSession) error
}

// TelemedicineService는 원격진료 서비스 핵심 로직입니다.
type TelemedicineService struct {
	log         *zap.Logger
	consultRepo ConsultationRepository
	doctorRepo  DoctorRepository
	sessionRepo VideoSessionRepository
}

// NewTelemedicineService는 TelemedicineService를 생성합니다.
func NewTelemedicineService(log *zap.Logger, consultRepo ConsultationRepository, doctorRepo DoctorRepository, sessionRepo VideoSessionRepository) *TelemedicineService {
	return &TelemedicineService{
		log:         log,
		consultRepo: consultRepo,
		doctorRepo:  doctorRepo,
		sessionRepo: sessionRepo,
	}
}

// CreateConsultation은 새 상담을 생성합니다.
func (s *TelemedicineService) CreateConsultation(ctx context.Context, patientUserID string, specialty DoctorSpecialty, chiefComplaint, description string) (*Consultation, error) {
	if patientUserID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "patient_user_id는 필수입니다")
	}
	if chiefComplaint == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "chief_complaint는 필수입니다")
	}

	consultation := &Consultation{
		ID:             uuid.New().String(),
		PatientUserID:  patientUserID,
		Specialty:      specialty,
		ChiefComplaint: chiefComplaint,
		Description:    description,
		Status:         StatusRequested,
		CreatedAt:      time.Now(),
	}

	if err := s.consultRepo.Save(ctx, consultation); err != nil {
		return nil, fmt.Errorf("상담 저장 실패: %w", err)
	}

	s.log.Info("상담 생성 완료",
		zap.String("consultation_id", consultation.ID),
		zap.String("patient_user_id", patientUserID),
	)

	return consultation, nil
}

// GetConsultation은 상담을 조회합니다.
func (s *TelemedicineService) GetConsultation(ctx context.Context, consultationID string) (*Consultation, error) {
	if consultationID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "consultation_id는 필수입니다")
	}
	return s.consultRepo.FindByID(ctx, consultationID)
}

// ListConsultations는 사용자의 상담 목록을 조회합니다.
func (s *TelemedicineService) ListConsultations(ctx context.Context, userID string, statusFilter ConsultationStatus, limit, offset int) ([]*Consultation, int, error) {
	if userID == "" {
		return nil, 0, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}
	if limit <= 0 {
		limit = 20
	}
	return s.consultRepo.FindByUserID(ctx, userID, statusFilter, limit, offset)
}

// MatchDoctor는 조건에 맞는 의사 목록을 반환합니다.
func (s *TelemedicineService) MatchDoctor(ctx context.Context, specialty DoctorSpecialty, language string) ([]*DoctorProfile, error) {
	doctors, err := s.doctorRepo.FindBySpecialty(ctx, specialty, language)
	if err != nil {
		return nil, fmt.Errorf("의사 조회 실패: %w", err)
	}

	// 가용 의사만 필터링
	var available []*DoctorProfile
	for _, d := range doctors {
		if d.IsAvailable {
			available = append(available, d)
		}
	}

	return available, nil
}

// StartVideoSession은 비디오 세션을 생성합니다.
func (s *TelemedicineService) StartVideoSession(ctx context.Context, consultationID, userID string) (*VideoSession, error) {
	if consultationID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "consultation_id는 필수입니다")
	}
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}

	// 상담 존재 확인
	consultation, err := s.consultRepo.FindByID(ctx, consultationID)
	if err != nil {
		return nil, err
	}

	roomID := uuid.New().String()
	session := &VideoSession{
		ID:             uuid.New().String(),
		ConsultationID: consultationID,
		RoomURL:        fmt.Sprintf("https://meet.manpasik.com/%s", roomID),
		Token:          uuid.New().String(),
		Status:         SessionConnected,
		StartedAt:      time.Now(),
	}

	if err := s.sessionRepo.Save(ctx, session); err != nil {
		return nil, fmt.Errorf("비디오 세션 저장 실패: %w", err)
	}

	// 상담 상태를 InProgress로 변경
	consultation.Status = StatusInProgress
	consultation.StartedAt = session.StartedAt
	if err := s.consultRepo.Update(ctx, consultation); err != nil {
		return nil, fmt.Errorf("상담 상태 업데이트 실패: %w", err)
	}

	s.log.Info("비디오 세션 시작",
		zap.String("session_id", session.ID),
		zap.String("consultation_id", consultationID),
	)

	return session, nil
}

// EndVideoSession은 비디오 세션을 종료합니다.
func (s *TelemedicineService) EndVideoSession(ctx context.Context, sessionID, consultationID, doctorNotes, diagnosis string) (*VideoSession, error) {
	if sessionID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "session_id는 필수입니다")
	}

	session, err := s.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	session.Status = SessionEnded
	session.EndedAt = now
	if !session.StartedAt.IsZero() {
		session.DurationSeconds = int(now.Sub(session.StartedAt).Seconds())
	}

	if err := s.sessionRepo.Update(ctx, session); err != nil {
		return nil, fmt.Errorf("비디오 세션 업데이트 실패: %w", err)
	}

	// 상담 업데이트
	if consultationID != "" {
		consultation, err := s.consultRepo.FindByID(ctx, consultationID)
		if err == nil {
			consultation.Status = StatusCompleted
			consultation.DoctorNotes = doctorNotes
			consultation.Diagnosis = diagnosis
			consultation.EndedAt = now
			if !consultation.StartedAt.IsZero() {
				consultation.DurationMinutes = int(now.Sub(consultation.StartedAt).Minutes())
			}
			s.consultRepo.Update(ctx, consultation)
		}
	}

	s.log.Info("비디오 세션 종료",
		zap.String("session_id", sessionID),
	)

	return session, nil
}

// RateConsultation은 상담에 평점을 부여합니다.
func (s *TelemedicineService) RateConsultation(ctx context.Context, consultationID string, rating float64) (float64, error) {
	if consultationID == "" {
		return 0, apperrors.New(apperrors.ErrInvalidInput, "consultation_id는 필수입니다")
	}
	if rating < 0 || rating > 5 {
		return 0, apperrors.New(apperrors.ErrInvalidInput, "rating은 0~5 사이여야 합니다")
	}

	consultation, err := s.consultRepo.FindByID(ctx, consultationID)
	if err != nil {
		return 0, err
	}

	consultation.Rating = rating
	if err := s.consultRepo.Update(ctx, consultation); err != nil {
		return 0, fmt.Errorf("상담 평점 업데이트 실패: %w", err)
	}

	s.log.Info("상담 평점 등록",
		zap.String("consultation_id", consultationID),
		zap.Float64("rating", rating),
	)

	return rating, nil
}
