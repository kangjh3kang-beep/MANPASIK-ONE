// Package postgres는 telemedicine-service의 PostgreSQL 저장소 구현입니다.
//
// DB 스키마: infrastructure/database/init/15-telemedicine.sql
// 테이블: doctors, consultations, video_sessions
package postgres

import (
	"context"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/telemedicine-service/internal/service"
)

// ============================================================================
// 열거형 변환 유틸
// ============================================================================

func specialtyToString(s service.DoctorSpecialty) string {
	switch s {
	case service.SpecialtyGeneral:
		return "general"
	case service.SpecialtyInternal:
		return "internal"
	case service.SpecialtyCardiology:
		return "cardiology"
	case service.SpecialtyEndocrinology:
		return "endocrinology"
	case service.SpecialtyDermatology:
		return "dermatology"
	case service.SpecialtyPediatrics:
		return "pediatrics"
	case service.SpecialtyPsychiatry:
		return "psychiatry"
	case service.SpecialtyOrthopedics:
		return "orthopedics"
	case service.SpecialtyOphthalmology:
		return "ophthalmology"
	case service.SpecialtyENT:
		return "ent"
	default:
		return "general"
	}
}

func stringToSpecialty(s string) service.DoctorSpecialty {
	switch s {
	case "general":
		return service.SpecialtyGeneral
	case "internal":
		return service.SpecialtyInternal
	case "cardiology":
		return service.SpecialtyCardiology
	case "endocrinology":
		return service.SpecialtyEndocrinology
	case "dermatology":
		return service.SpecialtyDermatology
	case "pediatrics":
		return service.SpecialtyPediatrics
	case "psychiatry":
		return service.SpecialtyPsychiatry
	case "orthopedics":
		return service.SpecialtyOrthopedics
	case "ophthalmology":
		return service.SpecialtyOphthalmology
	case "ent":
		return service.SpecialtyENT
	default:
		return service.SpecialtyUnknown
	}
}

func consultStatusToString(s service.ConsultationStatus) string {
	switch s {
	case service.StatusRequested:
		return "requested"
	case service.StatusMatched:
		return "matched"
	case service.StatusScheduled:
		return "scheduled"
	case service.StatusInProgress:
		return "in_progress"
	case service.StatusCompleted:
		return "completed"
	case service.StatusCancelled:
		return "cancelled"
	case service.StatusNoShow:
		return "no_show"
	default:
		return "requested"
	}
}

func stringToConsultStatus(s string) service.ConsultationStatus {
	switch s {
	case "requested":
		return service.StatusRequested
	case "matched":
		return service.StatusMatched
	case "scheduled":
		return service.StatusScheduled
	case "in_progress":
		return service.StatusInProgress
	case "completed":
		return service.StatusCompleted
	case "cancelled":
		return service.StatusCancelled
	case "no_show":
		return service.StatusNoShow
	default:
		return service.StatusUnknown
	}
}

func sessionStatusToString(s service.VideoSessionStatus) string {
	switch s {
	case service.SessionWaiting:
		return "waiting"
	case service.SessionConnected:
		return "connected"
	case service.SessionEnded:
		return "ended"
	case service.SessionFailed:
		return "failed"
	default:
		return "waiting"
	}
}

func stringToSessionStatus(s string) service.VideoSessionStatus {
	switch s {
	case "waiting":
		return service.SessionWaiting
	case "connected":
		return service.SessionConnected
	case "ended":
		return service.SessionEnded
	case "failed":
		return service.SessionFailed
	default:
		return service.SessionUnknown
	}
}

// ============================================================================
// ConsultationRepository
// ============================================================================

// ConsultationRepository는 PostgreSQL 기반 상담 저장소입니다.
type ConsultationRepository struct {
	pool *pgxpool.Pool
}

// NewConsultationRepository는 ConsultationRepository를 생성합니다.
func NewConsultationRepository(pool *pgxpool.Pool) *ConsultationRepository {
	return &ConsultationRepository{pool: pool}
}

// Save는 상담을 저장합니다.
func (r *ConsultationRepository) Save(ctx context.Context, c *service.Consultation) error {
	const q = `INSERT INTO consultations (consultation_id, patient_user_id, doctor_id, specialty, chief_complaint,
		description, status, diagnosis, doctor_notes, prescription_id, duration_minutes, rating,
		scheduled_at, started_at, ended_at, created_at)
		VALUES ($1, $2, $3, $4::doctor_specialty, $5, $6, $7::consultation_status, $8, $9, $10, $11, $12, $13, $14, $15, $16)`

	var doctorID, diagnosis, notes, prescriptionID *string
	var scheduledAt, startedAt, endedAt *time.Time

	if c.DoctorID != "" {
		doctorID = &c.DoctorID
	}
	if c.Diagnosis != "" {
		diagnosis = &c.Diagnosis
	}
	if c.DoctorNotes != "" {
		notes = &c.DoctorNotes
	}
	if c.PrescriptionID != "" {
		prescriptionID = &c.PrescriptionID
	}
	if !c.ScheduledAt.IsZero() {
		scheduledAt = &c.ScheduledAt
	}
	if !c.StartedAt.IsZero() {
		startedAt = &c.StartedAt
	}
	if !c.EndedAt.IsZero() {
		endedAt = &c.EndedAt
	}

	_, err := r.pool.Exec(ctx, q,
		c.ID, c.PatientUserID, doctorID, specialtyToString(c.Specialty),
		c.ChiefComplaint, c.Description, consultStatusToString(c.Status),
		diagnosis, notes, prescriptionID, c.DurationMinutes, c.Rating,
		scheduledAt, startedAt, endedAt, c.CreatedAt,
	)
	return err
}

// FindByID는 상담을 ID로 조회합니다.
func (r *ConsultationRepository) FindByID(ctx context.Context, id string) (*service.Consultation, error) {
	const q = `SELECT consultation_id, patient_user_id, COALESCE(doctor_id,''), specialty,
		chief_complaint, COALESCE(description,''), status, COALESCE(diagnosis,''),
		COALESCE(doctor_notes,''), COALESCE(prescription_id,''), duration_minutes, COALESCE(rating,0),
		COALESCE(scheduled_at, '0001-01-01'), COALESCE(started_at, '0001-01-01'),
		COALESCE(ended_at, '0001-01-01'), created_at
		FROM consultations WHERE consultation_id = $1`
	var c service.Consultation
	var specStr, statusStr string
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&c.ID, &c.PatientUserID, &c.DoctorID, &specStr,
		&c.ChiefComplaint, &c.Description, &statusStr, &c.Diagnosis,
		&c.DoctorNotes, &c.PrescriptionID, &c.DurationMinutes, &c.Rating,
		&c.ScheduledAt, &c.StartedAt, &c.EndedAt, &c.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	c.Specialty = stringToSpecialty(specStr)
	c.Status = stringToConsultStatus(statusStr)
	return &c, nil
}

// FindByUserID는 사용자의 상담 목록을 조회합니다.
func (r *ConsultationRepository) FindByUserID(ctx context.Context, userID string, statusFilter service.ConsultationStatus, limit, offset int) ([]*service.Consultation, int, error) {
	countQ := "SELECT COUNT(*) FROM consultations WHERE patient_user_id = $1"
	listQ := `SELECT consultation_id, patient_user_id, COALESCE(doctor_id,''), specialty,
		chief_complaint, COALESCE(description,''), status, COALESCE(diagnosis,''),
		COALESCE(doctor_notes,''), COALESCE(prescription_id,''), duration_minutes, COALESCE(rating,0),
		COALESCE(scheduled_at, '0001-01-01'), COALESCE(started_at, '0001-01-01'),
		COALESCE(ended_at, '0001-01-01'), created_at
		FROM consultations WHERE patient_user_id = $1`

	args := []interface{}{userID}
	idx := 2

	if statusFilter != service.StatusUnknown {
		filter := " AND status = $2::consultation_status"
		countQ += filter
		listQ += filter
		args = append(args, consultStatusToString(statusFilter))
		idx++
	}

	var total int
	if err := r.pool.QueryRow(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	listQ += " ORDER BY created_at DESC LIMIT $" + itoa(idx) + " OFFSET $" + itoa(idx+1)
	queryArgs := append(args, limit, offset)

	rows, err := r.pool.Query(ctx, listQ, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var consultations []*service.Consultation
	for rows.Next() {
		var c service.Consultation
		var specStr, statusStr string
		if err := rows.Scan(
			&c.ID, &c.PatientUserID, &c.DoctorID, &specStr,
			&c.ChiefComplaint, &c.Description, &statusStr, &c.Diagnosis,
			&c.DoctorNotes, &c.PrescriptionID, &c.DurationMinutes, &c.Rating,
			&c.ScheduledAt, &c.StartedAt, &c.EndedAt, &c.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		c.Specialty = stringToSpecialty(specStr)
		c.Status = stringToConsultStatus(statusStr)
		consultations = append(consultations, &c)
	}
	return consultations, total, rows.Err()
}

// Update는 상담을 업데이트합니다.
func (r *ConsultationRepository) Update(ctx context.Context, c *service.Consultation) error {
	var doctorID, diagnosis, notes, prescriptionID *string
	var scheduledAt, startedAt, endedAt *time.Time

	if c.DoctorID != "" {
		doctorID = &c.DoctorID
	}
	if c.Diagnosis != "" {
		diagnosis = &c.Diagnosis
	}
	if c.DoctorNotes != "" {
		notes = &c.DoctorNotes
	}
	if c.PrescriptionID != "" {
		prescriptionID = &c.PrescriptionID
	}
	if !c.ScheduledAt.IsZero() {
		scheduledAt = &c.ScheduledAt
	}
	if !c.StartedAt.IsZero() {
		startedAt = &c.StartedAt
	}
	if !c.EndedAt.IsZero() {
		endedAt = &c.EndedAt
	}

	const q = `UPDATE consultations SET doctor_id = $1, status = $2::consultation_status, diagnosis = $3,
		doctor_notes = $4, prescription_id = $5, duration_minutes = $6, rating = $7,
		scheduled_at = $8, started_at = $9, ended_at = $10, updated_at = NOW()
		WHERE consultation_id = $11`
	_, err := r.pool.Exec(ctx, q,
		doctorID, consultStatusToString(c.Status), diagnosis, notes, prescriptionID,
		c.DurationMinutes, c.Rating, scheduledAt, startedAt, endedAt, c.ID,
	)
	return err
}

// ============================================================================
// DoctorRepository
// ============================================================================

// DoctorRepository는 PostgreSQL 기반 의사 저장소입니다.
type DoctorRepository struct {
	pool *pgxpool.Pool
}

// NewDoctorRepository는 DoctorRepository를 생성합니다.
func NewDoctorRepository(pool *pgxpool.Pool) *DoctorRepository {
	return &DoctorRepository{pool: pool}
}

// FindBySpecialty는 전문 분야별 의사를 조회합니다.
func (r *DoctorRepository) FindBySpecialty(ctx context.Context, specialty service.DoctorSpecialty, language string) ([]*service.DoctorProfile, error) {
	q := `SELECT doctor_id, name, specialty, COALESCE(hospital,''), COALESCE(license_number,''),
		experience_years, rating, total_consultations, is_available, COALESCE(languages, ARRAY['ko']),
		COALESCE(profile_image_url,'')
		FROM doctors WHERE 1=1`

	var args []interface{}
	idx := 1

	if specialty != service.SpecialtyUnknown {
		q += " AND specialty = $" + itoa(idx) + "::doctor_specialty"
		args = append(args, specialtyToString(specialty))
		idx++
	}

	if language != "" {
		q += " AND $" + itoa(idx) + " = ANY(languages)"
		args = append(args, language)
		idx++
	}

	q += " ORDER BY rating DESC, total_consultations DESC"

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doctors []*service.DoctorProfile
	for rows.Next() {
		var d service.DoctorProfile
		var specStr string
		if err := rows.Scan(
			&d.ID, &d.Name, &specStr, &d.Hospital, &d.LicenseNumber,
			&d.ExperienceYears, &d.Rating, &d.TotalConsultations, &d.IsAvailable,
			&d.Languages, &d.ProfileImageURL,
		); err != nil {
			return nil, err
		}
		d.Specialty = stringToSpecialty(specStr)
		doctors = append(doctors, &d)
	}
	return doctors, rows.Err()
}

// FindByID는 의사를 ID로 조회합니다.
func (r *DoctorRepository) FindByID(ctx context.Context, id string) (*service.DoctorProfile, error) {
	const q = `SELECT doctor_id, name, specialty, COALESCE(hospital,''), COALESCE(license_number,''),
		experience_years, rating, total_consultations, is_available, COALESCE(languages, ARRAY['ko']),
		COALESCE(profile_image_url,'')
		FROM doctors WHERE doctor_id = $1`
	var d service.DoctorProfile
	var specStr string
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&d.ID, &d.Name, &specStr, &d.Hospital, &d.LicenseNumber,
		&d.ExperienceYears, &d.Rating, &d.TotalConsultations, &d.IsAvailable,
		&d.Languages, &d.ProfileImageURL,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	d.Specialty = stringToSpecialty(specStr)
	return &d, nil
}

// ============================================================================
// VideoSessionRepository
// ============================================================================

// VideoSessionRepository는 PostgreSQL 기반 비디오 세션 저장소입니다.
type VideoSessionRepository struct {
	pool *pgxpool.Pool
}

// NewVideoSessionRepository는 VideoSessionRepository를 생성합니다.
func NewVideoSessionRepository(pool *pgxpool.Pool) *VideoSessionRepository {
	return &VideoSessionRepository{pool: pool}
}

// Save는 비디오 세션을 저장합니다.
func (r *VideoSessionRepository) Save(ctx context.Context, s *service.VideoSession) error {
	var startedAt, endedAt *time.Time
	if !s.StartedAt.IsZero() {
		startedAt = &s.StartedAt
	}
	if !s.EndedAt.IsZero() {
		endedAt = &s.EndedAt
	}

	const q = `INSERT INTO video_sessions (session_id, consultation_id, room_url, token, status, started_at, ended_at, duration_seconds)
		VALUES ($1, $2, $3, $4, $5::video_session_status, $6, $7, $8)`
	_, err := r.pool.Exec(ctx, q,
		s.ID, s.ConsultationID, s.RoomURL, s.Token,
		sessionStatusToString(s.Status), startedAt, endedAt, s.DurationSeconds,
	)
	return err
}

// FindByID는 세션을 ID로 조회합니다.
func (r *VideoSessionRepository) FindByID(ctx context.Context, id string) (*service.VideoSession, error) {
	const q = `SELECT session_id, consultation_id, room_url, token, status,
		COALESCE(started_at, '0001-01-01'), COALESCE(ended_at, '0001-01-01'), duration_seconds
		FROM video_sessions WHERE session_id = $1`
	var s service.VideoSession
	var statusStr string
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&s.ID, &s.ConsultationID, &s.RoomURL, &s.Token, &statusStr,
		&s.StartedAt, &s.EndedAt, &s.DurationSeconds,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	s.Status = stringToSessionStatus(statusStr)
	return &s, nil
}

// FindByConsultationID는 상담 ID로 세션을 조회합니다.
func (r *VideoSessionRepository) FindByConsultationID(ctx context.Context, consultationID string) (*service.VideoSession, error) {
	const q = `SELECT session_id, consultation_id, room_url, token, status,
		COALESCE(started_at, '0001-01-01'), COALESCE(ended_at, '0001-01-01'), duration_seconds
		FROM video_sessions WHERE consultation_id = $1 ORDER BY created_at DESC LIMIT 1`
	var s service.VideoSession
	var statusStr string
	err := r.pool.QueryRow(ctx, q, consultationID).Scan(
		&s.ID, &s.ConsultationID, &s.RoomURL, &s.Token, &statusStr,
		&s.StartedAt, &s.EndedAt, &s.DurationSeconds,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	s.Status = stringToSessionStatus(statusStr)
	return &s, nil
}

// Update는 세션을 업데이트합니다.
func (r *VideoSessionRepository) Update(ctx context.Context, s *service.VideoSession) error {
	var startedAt, endedAt *time.Time
	if !s.StartedAt.IsZero() {
		startedAt = &s.StartedAt
	}
	if !s.EndedAt.IsZero() {
		endedAt = &s.EndedAt
	}

	const q = `UPDATE video_sessions SET status = $1::video_session_status, started_at = $2,
		ended_at = $3, duration_seconds = $4 WHERE session_id = $5`
	_, err := r.pool.Exec(ctx, q,
		sessionStatusToString(s.Status), startedAt, endedAt, s.DurationSeconds, s.ID,
	)
	return err
}

// ============================================================================
// 유틸
// ============================================================================

func itoa(n int) string {
	return strconv.Itoa(n)
}
