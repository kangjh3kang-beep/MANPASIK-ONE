// Package postgres는 reservation-service의 PostgreSQL 저장소 구현입니다.
package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/reservation-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
)

// facilityTypeToString converts FacilityType enum to DB string.
func facilityTypeToString(ft service.FacilityType) string {
	switch ft {
	case service.FacilityHospital:
		return "hospital"
	case service.FacilityClinic:
		return "clinic"
	case service.FacilityPharmacy:
		return "pharmacy"
	case service.FacilityDental:
		return "dental"
	case service.FacilityOriental:
		return "oriental"
	default:
		return ""
	}
}

// stringToFacilityType converts DB string to FacilityType enum.
func stringToFacilityType(s string) service.FacilityType {
	switch s {
	case "hospital":
		return service.FacilityHospital
	case "clinic":
		return service.FacilityClinic
	case "pharmacy":
		return service.FacilityPharmacy
	case "dental":
		return service.FacilityDental
	case "oriental":
		return service.FacilityOriental
	default:
		return service.FacilityUnknown
	}
}

// reservationStatusToString converts ReservationStatus to DB enum string.
func reservationStatusToString(s service.ReservationStatus) string {
	switch s {
	case service.ResPending:
		return "pending"
	case service.ResConfirmed:
		return "confirmed"
	case service.ResCompleted:
		return "completed"
	case service.ResCancelled:
		return "cancelled"
	case service.ResNoShow:
		return "no_show"
	default:
		return ""
	}
}

// stringToReservationStatus converts DB enum string to ReservationStatus.
func stringToReservationStatus(s string) service.ReservationStatus {
	switch s {
	case "pending":
		return service.ResPending
	case "confirmed":
		return service.ResConfirmed
	case "completed":
		return service.ResCompleted
	case "cancelled":
		return service.ResCancelled
	case "no_show":
		return service.ResNoShow
	default:
		return service.ResUnknown
	}
}

// specialtyToString converts Specialty enum to DB string.
func specialtyToString(sp service.Specialty) string {
	switch sp {
	case service.SpecGeneral:
		return "general"
	case service.SpecInternal:
		return "internal"
	case service.SpecCardiology:
		return "cardiology"
	case service.SpecEndocrinology:
		return "endocrinology"
	case service.SpecDermatology:
		return "dermatology"
	case service.SpecPediatrics:
		return "pediatrics"
	case service.SpecPsychiatry:
		return "psychiatry"
	case service.SpecOrthopedics:
		return "orthopedics"
	case service.SpecOphthalmology:
		return "ophthalmology"
	case service.SpecENT:
		return "ent"
	default:
		return ""
	}
}

// stringToSpecialty converts DB string to Specialty enum.
func stringToSpecialty(s string) service.Specialty {
	switch strings.ToLower(s) {
	case "general":
		return service.SpecGeneral
	case "internal":
		return service.SpecInternal
	case "cardiology":
		return service.SpecCardiology
	case "endocrinology":
		return service.SpecEndocrinology
	case "dermatology":
		return service.SpecDermatology
	case "pediatrics":
		return service.SpecPediatrics
	case "psychiatry":
		return service.SpecPsychiatry
	case "orthopedics":
		return service.SpecOrthopedics
	case "ophthalmology":
		return service.SpecOphthalmology
	case "ent":
		return service.SpecENT
	default:
		return service.SpecUnknown
	}
}

// ============================================================================
// RegionRepository
// ============================================================================

// RegionRepository는 PostgreSQL 기반 RegionRepository 구현입니다.
type RegionRepository struct {
	pool *pgxpool.Pool
}

// NewRegionRepository는 RegionRepository를 생성합니다.
func NewRegionRepository(pool *pgxpool.Pool) *RegionRepository {
	return &RegionRepository{pool: pool}
}

// ListByCountry는 국가 코드로 지역 목록을 조회합니다.
func (r *RegionRepository) ListByCountry(ctx context.Context, countryCode string) ([]*service.Region, error) {
	const q = `SELECT id, country_code, COALESCE(region_code, ''), COALESCE(district_code, ''),
		COALESCE(name_en, ''), COALESCE(name_local, ''), COALESCE(parent_id::text, '')
		FROM regions WHERE country_code = $1 ORDER BY name_en`
	rows, err := r.pool.Query(ctx, q, countryCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var regions []*service.Region
	for rows.Next() {
		var reg service.Region
		if err := rows.Scan(
			&reg.ID, &reg.CountryCode, &reg.RegionCode, &reg.DistrictCode,
			&reg.Name, &reg.NameLocal, &reg.ParentID,
		); err != nil {
			return nil, err
		}
		regions = append(regions, &reg)
	}
	return regions, rows.Err()
}

// ListByRegion는 국가+지역 코드로 구/군 목록을 조회합니다.
func (r *RegionRepository) ListByRegion(ctx context.Context, countryCode, regionCode string) ([]*service.Region, error) {
	const q = `SELECT id, country_code, COALESCE(region_code, ''), COALESCE(district_code, ''),
		COALESCE(name_en, ''), COALESCE(name_local, ''), COALESCE(parent_id::text, '')
		FROM regions WHERE country_code = $1 AND region_code = $2 ORDER BY name_en`
	rows, err := r.pool.Query(ctx, q, countryCode, regionCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var regions []*service.Region
	for rows.Next() {
		var reg service.Region
		if err := rows.Scan(
			&reg.ID, &reg.CountryCode, &reg.RegionCode, &reg.DistrictCode,
			&reg.Name, &reg.NameLocal, &reg.ParentID,
		); err != nil {
			return nil, err
		}
		regions = append(regions, &reg)
	}
	return regions, rows.Err()
}

// GetByID는 지역 ID로 지역을 조회합니다.
func (r *RegionRepository) GetByID(ctx context.Context, regionID string) (*service.Region, error) {
	const q = `SELECT id, country_code, COALESCE(region_code, ''), COALESCE(district_code, ''),
		COALESCE(name_en, ''), COALESCE(name_local, ''), COALESCE(parent_id::text, '')
		FROM regions WHERE id = $1`
	var reg service.Region
	err := r.pool.QueryRow(ctx, q, regionID).Scan(
		&reg.ID, &reg.CountryCode, &reg.RegionCode, &reg.DistrictCode,
		&reg.Name, &reg.NameLocal, &reg.ParentID,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.New(apperrors.ErrNotFound, "지역을 찾을 수 없습니다")
		}
		return nil, err
	}
	return &reg, nil
}

// ============================================================================
// FacilityRepository
// ============================================================================

// FacilityRepository는 PostgreSQL 기반 FacilityRepository 구현입니다.
type FacilityRepository struct {
	pool *pgxpool.Pool
}

// NewFacilityRepository는 FacilityRepository를 생성합니다.
func NewFacilityRepository(pool *pgxpool.Pool) *FacilityRepository {
	return &FacilityRepository{pool: pool}
}

// Search는 조건에 맞는 시설을 검색합니다.
func (r *FacilityRepository) Search(ctx context.Context, facilityType service.FacilityType, keyword string, specialty service.Specialty, limit, offset int) ([]*service.Facility, int, error) {
	// 동적 WHERE 절 구성
	var conditions []string
	var args []interface{}
	argIdx := 1

	if facilityType != service.FacilityUnknown {
		conditions = append(conditions, fmt.Sprintf("f.type = $%d::facility_type", argIdx))
		args = append(args, facilityTypeToString(facilityType))
		argIdx++
	}
	if keyword != "" {
		conditions = append(conditions, fmt.Sprintf("(LOWER(f.name) LIKE $%d OR LOWER(f.address) LIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+strings.ToLower(keyword)+"%")
		argIdx++
	}
	if specialty != service.SpecUnknown {
		conditions = append(conditions, fmt.Sprintf("$%d = ANY(f.specialties)", argIdx))
		args = append(args, specialtyToString(specialty))
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// count 쿼리
	countQ := fmt.Sprintf(`SELECT COUNT(*) FROM facilities f %s`, whereClause)
	var total int
	if err := r.pool.QueryRow(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// 데이터 쿼리
	dataQ := fmt.Sprintf(`SELECT f.facility_id, f.name, f.type, f.address, COALESCE(f.phone, ''),
		COALESCE(f.latitude, 0), COALESCE(f.longitude, 0), COALESCE(f.rating, 0), COALESCE(f.review_count, 0),
		COALESCE(f.specialties, '{}'), COALESCE(f.operating_hours, ''), COALESCE(f.is_open_now, false),
		COALESCE(f.accepts_reservation, true), COALESCE(f.image_url, ''),
		COALESCE(f.country_code, ''), COALESCE(f.region_code, ''), COALESCE(f.district_code, ''),
		COALESCE(f.timezone, ''), COALESCE(f.has_telemedicine, false)
		FROM facilities f %s ORDER BY f.rating DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, dataQ, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var facilities []*service.Facility
	for rows.Next() {
		var f service.Facility
		var ftStr string
		var specialtiesArr []string
		if err := rows.Scan(
			&f.ID, &f.Name, &ftStr, &f.Address, &f.Phone,
			&f.Latitude, &f.Longitude, &f.Rating, &f.ReviewCount,
			&specialtiesArr, &f.OperatingHours, &f.IsOpenNow,
			&f.AcceptsReservation, &f.ImageURL,
			&f.CountryCode, &f.RegionCode, &f.DistrictCode,
			&f.Timezone, &f.HasTelemedicine,
		); err != nil {
			return nil, 0, err
		}
		f.Type = stringToFacilityType(ftStr)
		for _, s := range specialtiesArr {
			sp := stringToSpecialty(s)
			if sp != service.SpecUnknown {
				f.Specialties = append(f.Specialties, sp)
			}
		}
		facilities = append(facilities, &f)
	}
	return facilities, total, rows.Err()
}

// FindByID는 ID로 시설을 조회합니다.
func (r *FacilityRepository) FindByID(ctx context.Context, id string) (*service.Facility, error) {
	const q = `SELECT f.facility_id, f.name, f.type, f.address, COALESCE(f.phone, ''),
		COALESCE(f.latitude, 0), COALESCE(f.longitude, 0), COALESCE(f.rating, 0), COALESCE(f.review_count, 0),
		COALESCE(f.specialties, '{}'), COALESCE(f.operating_hours, ''), COALESCE(f.is_open_now, false),
		COALESCE(f.accepts_reservation, true), COALESCE(f.image_url, ''),
		COALESCE(f.country_code, ''), COALESCE(f.region_code, ''), COALESCE(f.district_code, ''),
		COALESCE(f.timezone, ''), COALESCE(f.has_telemedicine, false)
		FROM facilities f WHERE f.facility_id = $1`
	var f service.Facility
	var ftStr string
	var specialtiesArr []string
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&f.ID, &f.Name, &ftStr, &f.Address, &f.Phone,
		&f.Latitude, &f.Longitude, &f.Rating, &f.ReviewCount,
		&specialtiesArr, &f.OperatingHours, &f.IsOpenNow,
		&f.AcceptsReservation, &f.ImageURL,
		&f.CountryCode, &f.RegionCode, &f.DistrictCode,
		&f.Timezone, &f.HasTelemedicine,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.New(apperrors.ErrNotFound, "시설을 찾을 수 없습니다")
		}
		return nil, err
	}
	f.Type = stringToFacilityType(ftStr)
	for _, s := range specialtiesArr {
		sp := stringToSpecialty(s)
		if sp != service.SpecUnknown {
			f.Specialties = append(f.Specialties, sp)
		}
	}
	return &f, nil
}

// ============================================================================
// DoctorRepository
// ============================================================================

// DoctorRepository는 PostgreSQL 기반 DoctorRepository 구현입니다.
type DoctorRepository struct {
	pool *pgxpool.Pool
}

// NewDoctorRepository는 DoctorRepository를 생성합니다.
func NewDoctorRepository(pool *pgxpool.Pool) *DoctorRepository {
	return &DoctorRepository{pool: pool}
}

// ListByFacility는 시설별 의사 목록을 조회합니다.
func (r *DoctorRepository) ListByFacility(ctx context.Context, facilityID string, specialty string) ([]*service.Doctor, error) {
	var q string
	var args []interface{}
	if specialty != "" {
		q = `SELECT id, facility_id, COALESCE(user_id::text, ''), name, COALESCE(specialty, ''),
			COALESCE(license_number, ''), COALESCE(languages, '{}'), COALESCE(is_available, true),
			COALESCE(rating, 0), COALESCE(total_consultations, 0)
			FROM doctors WHERE facility_id = $1 AND specialty = $2 ORDER BY rating DESC`
		args = []interface{}{facilityID, specialty}
	} else {
		q = `SELECT id, facility_id, COALESCE(user_id::text, ''), name, COALESCE(specialty, ''),
			COALESCE(license_number, ''), COALESCE(languages, '{}'), COALESCE(is_available, true),
			COALESCE(rating, 0), COALESCE(total_consultations, 0)
			FROM doctors WHERE facility_id = $1 ORDER BY rating DESC`
		args = []interface{}{facilityID}
	}

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doctors []*service.Doctor
	for rows.Next() {
		var d service.Doctor
		var languages []string
		if err := rows.Scan(
			&d.ID, &d.FacilityID, &d.UserID, &d.Name, &d.Specialty,
			&d.LicenseNumber, &languages, &d.IsAvailable,
			&d.Rating, &d.TotalConsultations,
		); err != nil {
			return nil, err
		}
		d.Languages = languages
		doctors = append(doctors, &d)
	}
	return doctors, rows.Err()
}

// FindByID는 의사 ID로 의사를 조회합니다.
func (r *DoctorRepository) FindByID(ctx context.Context, doctorID string) (*service.Doctor, error) {
	const q = `SELECT id, facility_id, COALESCE(user_id::text, ''), name, COALESCE(specialty, ''),
		COALESCE(license_number, ''), COALESCE(languages, '{}'), COALESCE(is_available, true),
		COALESCE(rating, 0), COALESCE(total_consultations, 0)
		FROM doctors WHERE id = $1`
	var d service.Doctor
	var languages []string
	err := r.pool.QueryRow(ctx, q, doctorID).Scan(
		&d.ID, &d.FacilityID, &d.UserID, &d.Name, &d.Specialty,
		&d.LicenseNumber, &languages, &d.IsAvailable,
		&d.Rating, &d.TotalConsultations,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.New(apperrors.ErrNotFound, "의사를 찾을 수 없습니다")
		}
		return nil, err
	}
	d.Languages = languages
	return &d, nil
}

// ============================================================================
// SlotRepository
// ============================================================================

// SlotRepository는 PostgreSQL 기반 SlotRepository 구현입니다.
type SlotRepository struct {
	pool *pgxpool.Pool
}

// NewSlotRepository는 SlotRepository를 생성합니다.
func NewSlotRepository(pool *pgxpool.Pool) *SlotRepository {
	return &SlotRepository{pool: pool}
}

// FindByFacilityAndDate는 시설 ID와 날짜로 예약 가능 시간대를 조회합니다.
func (r *SlotRepository) FindByFacilityAndDate(ctx context.Context, facilityID string, date time.Time, doctorID string, specialty service.Specialty) ([]*service.TimeSlot, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	conditions = append(conditions, fmt.Sprintf("ts.facility_id = $%d", argIdx))
	args = append(args, facilityID)
	argIdx++

	conditions = append(conditions, fmt.Sprintf("ts.slot_date = $%d", argIdx))
	args = append(args, date.Format("2006-01-02"))
	argIdx++

	conditions = append(conditions, "ts.is_available = true")

	if doctorID != "" {
		conditions = append(conditions, fmt.Sprintf("ts.doctor_id = $%d", argIdx))
		args = append(args, doctorID)
		argIdx++
	}

	whereClause := strings.Join(conditions, " AND ")
	q := fmt.Sprintf(`SELECT ts.slot_id, ts.start_time, ts.end_time, ts.is_available,
		COALESCE(ts.doctor_id, ''), COALESCE(ts.doctor_name, '')
		FROM time_slots ts WHERE %s ORDER BY ts.start_time`, whereClause)

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []*service.TimeSlot
	for rows.Next() {
		var s service.TimeSlot
		if err := rows.Scan(
			&s.ID, &s.StartTime, &s.EndTime, &s.IsAvailable,
			&s.DoctorID, &s.DoctorName,
		); err != nil {
			return nil, err
		}
		slots = append(slots, &s)
	}
	return slots, rows.Err()
}

// ============================================================================
// ReservationRepository
// ============================================================================

// ReservationRepository는 PostgreSQL 기반 ReservationRepository 구현입니다.
type ReservationRepository struct {
	pool *pgxpool.Pool
}

// NewReservationRepository는 ReservationRepository를 생성합니다.
func NewReservationRepository(pool *pgxpool.Pool) *ReservationRepository {
	return &ReservationRepository{pool: pool}
}

// Save는 예약을 저장합니다.
func (r *ReservationRepository) Save(ctx context.Context, res *service.Reservation) error {
	const q = `INSERT INTO reservations (reservation_id, user_id, facility_id, facility_name, doctor_id, doctor_name, specialty, status, reason, notes, scheduled_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8::reservation_status, $9, $10, $11, $12, $13)`
	_, err := r.pool.Exec(ctx, q,
		res.ID,
		res.UserID,
		res.FacilityID,
		res.FacilityName,
		res.DoctorID,
		res.DoctorName,
		specialtyToString(res.Specialty),
		reservationStatusToString(res.Status),
		res.Reason,
		res.Notes,
		res.ScheduledAt,
		res.CreatedAt,
		res.UpdatedAt,
	)
	return err
}

// FindByID는 ID로 예약을 조회합니다.
func (r *ReservationRepository) FindByID(ctx context.Context, id string) (*service.Reservation, error) {
	const q = `SELECT reservation_id, user_id, facility_id, COALESCE(facility_name, ''),
		COALESCE(doctor_id, ''), COALESCE(doctor_name, ''), COALESCE(specialty, ''),
		status, COALESCE(reason, ''), COALESCE(notes, ''), scheduled_at, created_at, updated_at
		FROM reservations WHERE reservation_id = $1`
	var res service.Reservation
	var statusStr, specialtyStr string
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&res.ID, &res.UserID, &res.FacilityID, &res.FacilityName,
		&res.DoctorID, &res.DoctorName, &specialtyStr,
		&statusStr, &res.Reason, &res.Notes, &res.ScheduledAt, &res.CreatedAt, &res.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.New(apperrors.ErrNotFound, "예약을 찾을 수 없습니다")
		}
		return nil, err
	}
	res.Status = stringToReservationStatus(statusStr)
	res.Specialty = stringToSpecialty(specialtyStr)
	return &res, nil
}

// FindByUserID는 사용자 예약 목록을 조회합니다.
func (r *ReservationRepository) FindByUserID(ctx context.Context, userID string, statusFilter service.ReservationStatus, limit, offset int) ([]*service.Reservation, int, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIdx))
	args = append(args, userID)
	argIdx++

	if statusFilter != service.ResUnknown {
		conditions = append(conditions, fmt.Sprintf("status = $%d::reservation_status", argIdx))
		args = append(args, reservationStatusToString(statusFilter))
		argIdx++
	}

	whereClause := strings.Join(conditions, " AND ")

	// count 쿼리
	countQ := fmt.Sprintf(`SELECT COUNT(*) FROM reservations WHERE %s`, whereClause)
	var total int
	if err := r.pool.QueryRow(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// 데이터 쿼리
	dataQ := fmt.Sprintf(`SELECT reservation_id, user_id, facility_id, COALESCE(facility_name, ''),
		COALESCE(doctor_id, ''), COALESCE(doctor_name, ''), COALESCE(specialty, ''),
		status, COALESCE(reason, ''), COALESCE(notes, ''), scheduled_at, created_at, updated_at
		FROM reservations WHERE %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		whereClause, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, dataQ, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var reservations []*service.Reservation
	for rows.Next() {
		var res service.Reservation
		var statusStr, specialtyStr string
		if err := rows.Scan(
			&res.ID, &res.UserID, &res.FacilityID, &res.FacilityName,
			&res.DoctorID, &res.DoctorName, &specialtyStr,
			&statusStr, &res.Reason, &res.Notes, &res.ScheduledAt, &res.CreatedAt, &res.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		res.Status = stringToReservationStatus(statusStr)
		res.Specialty = stringToSpecialty(specialtyStr)
		reservations = append(reservations, &res)
	}
	return reservations, total, rows.Err()
}

// Update는 예약을 업데이트합니다.
func (r *ReservationRepository) Update(ctx context.Context, res *service.Reservation) error {
	const q = `UPDATE reservations SET status = $1::reservation_status, reason = $2, notes = $3, updated_at = $4
		WHERE reservation_id = $5`
	tag, err := r.pool.Exec(ctx, q,
		reservationStatusToString(res.Status),
		res.Reason,
		res.Notes,
		res.UpdatedAt,
		res.ID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apperrors.New(apperrors.ErrNotFound, "예약을 찾을 수 없습니다")
	}
	return nil
}
