// Package service는 reservation-service의 비즈니스 로직을 구현합니다.
package service

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/manpasik/backend/shared/errors"
	"go.uber.org/zap"
)

// Haversine calculates distance between two lat/lon points in km.
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371.0 // Earth radius km
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

// Region represents a geographic region hierarchy.
type Region struct {
	ID           string
	CountryCode  string // "KR", "US", "JP"
	RegionCode   string // "seoul", "tokyo", "new-york"
	DistrictCode string // "gangnam", "shibuya", "manhattan"
	Name         string
	NameLocal    string // Local name (한국어, 日本語)
	ParentID     string
	Timezone     string
}

// RegionRepository interface
type RegionRepository interface {
	ListByCountry(ctx context.Context, countryCode string) ([]*Region, error)
	ListByRegion(ctx context.Context, countryCode, regionCode string) ([]*Region, error)
	GetByID(ctx context.Context, regionID string) (*Region, error)
}

// FacilityType은 의료 시설 유형입니다.
type FacilityType int

const (
	FacilityUnknown  FacilityType = iota
	FacilityHospital              // 병원
	FacilityClinic                // 의원
	FacilityPharmacy              // 약국
	FacilityDental                // 치과
	FacilityOriental              // 한의원
)

// ReservationStatus는 예약 상태입니다.
type ReservationStatus int

const (
	ResUnknown   ReservationStatus = iota
	ResPending                     // 대기
	ResConfirmed                   // 확인
	ResCompleted                   // 완료
	ResCancelled                   // 취소
	ResNoShow                      // 노쇼
)

// Specialty는 의사 전문 분야입니다.
type Specialty int

const (
	SpecUnknown       Specialty = iota
	SpecGeneral                 // 일반의
	SpecInternal                // 내과
	SpecCardiology              // 심장내과
	SpecEndocrinology           // 내분비내과
	SpecDermatology             // 피부과
	SpecPediatrics              // 소아과
	SpecPsychiatry              // 정신과
	SpecOrthopedics             // 정형외과
	SpecOphthalmology           // 안과
	SpecENT                     // 이비인후과
)

// Facility는 의료 시설 도메인 객체입니다.
type Facility struct {
	ID                 string
	Name               string
	Type               FacilityType
	Address            string
	Phone              string
	Latitude           float64
	Longitude          float64
	DistanceKM         float64
	Rating             float64
	ReviewCount        int
	Specialties        []Specialty
	OperatingHours     string
	IsOpenNow          bool
	AcceptsReservation bool
	ImageURL           string
	CountryCode        string // "KR", "US", "JP"
	RegionCode         string // "seoul", "tokyo"
	DistrictCode       string // "gangnam", "shibuya"
	Timezone           string // "Asia/Seoul", "Asia/Tokyo"
	HasTelemedicine    bool
}

// Doctor는 의사 도메인 객체입니다.
type Doctor struct {
	ID                   string
	FacilityID           string
	UserID               string
	Name                 string
	Specialty            string // e.g. "internal", "cardiology"
	LicenseNumber        string
	Languages            []string
	IsAvailable          bool
	Rating               float64
	TotalConsultations   int
	NextAvailableAt      time.Time
	ConsultationFee      int32
	AcceptsTelemedicine  bool
	AvailableRegionCodes []string
}

// TimeSlot은 예약 가능 시간대입니다.
type TimeSlot struct {
	ID          string
	StartTime   time.Time
	EndTime     time.Time
	IsAvailable bool
	DoctorID    string
	DoctorName  string
}

// Reservation은 예약 도메인 객체입니다.
type Reservation struct {
	ID           string
	UserID       string
	FacilityID   string
	FacilityName string
	DoctorID     string
	DoctorName   string
	Specialty    Specialty
	Status       ReservationStatus
	Reason       string
	Notes        string
	ScheduledAt  time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// FacilityRepository는 의료 시설 저장소 인터페이스입니다.
type FacilityRepository interface {
	Search(ctx context.Context, facilityType FacilityType, keyword string, specialty Specialty, limit, offset int) ([]*Facility, int, error)
	FindByID(ctx context.Context, id string) (*Facility, error)
}

// SlotRepository는 예약 가능 시간대 저장소 인터페이스입니다.
type SlotRepository interface {
	FindByFacilityAndDate(ctx context.Context, facilityID string, date time.Time, doctorID string, specialty Specialty) ([]*TimeSlot, error)
}

// DoctorRepository는 의사 저장소 인터페이스입니다.
type DoctorRepository interface {
	ListByFacility(ctx context.Context, facilityID string, specialty string) ([]*Doctor, error)
	FindByID(ctx context.Context, doctorID string) (*Doctor, error)
}

// ReservationRepository는 예약 저장소 인터페이스입니다.
type ReservationRepository interface {
	Save(ctx context.Context, r *Reservation) error
	FindByID(ctx context.Context, id string) (*Reservation, error)
	FindByUserID(ctx context.Context, userID string, statusFilter ReservationStatus, limit, offset int) ([]*Reservation, int, error)
	Update(ctx context.Context, r *Reservation) error
}

// EventPublisher is an optional interface for publishing domain events.
type EventPublisher interface {
	Publish(ctx context.Context, event interface{}) error
}

// ReservationService는 예약 서비스 핵심 로직입니다.
type ReservationService struct {
	log             *zap.Logger
	facilityRepo    FacilityRepository
	slotRepo        SlotRepository
	reservationRepo ReservationRepository
	doctorRepo      DoctorRepository
	regionRepo      RegionRepository
	eventPub        EventPublisher
}

// NewReservationService는 ReservationService를 생성합니다.
// eventPub is optional and may be nil.
func NewReservationService(log *zap.Logger, facilityRepo FacilityRepository, slotRepo SlotRepository, reservationRepo ReservationRepository, doctorRepo DoctorRepository, regionRepo RegionRepository, eventPub ...EventPublisher) *ReservationService {
	svc := &ReservationService{
		log:             log,
		facilityRepo:    facilityRepo,
		slotRepo:        slotRepo,
		reservationRepo: reservationRepo,
		doctorRepo:      doctorRepo,
		regionRepo:      regionRepo,
	}
	if len(eventPub) > 0 && eventPub[0] != nil {
		svc.eventPub = eventPub[0]
	}
	return svc
}

// SearchFacilities는 의료 시설을 검색합니다.
// countryCode, regionCode, districtCode로 지역 필터링을 수행하고,
// userLat, userLon이 0이 아닌 경우 Haversine 거리를 계산합니다.
func (s *ReservationService) SearchFacilities(ctx context.Context, facilityType FacilityType, keyword string, specialty Specialty, limit, offset int, countryCode, regionCode, districtCode string, userLat, userLon float64) ([]*Facility, int, error) {
	if limit <= 0 {
		limit = 20
	}
	keyword = strings.TrimSpace(keyword)

	facilities, total, err := s.facilityRepo.Search(ctx, facilityType, keyword, specialty, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("시설 검색 실패: %w", err)
	}

	// 지역 필터링
	if countryCode != "" || regionCode != "" || districtCode != "" {
		var filtered []*Facility
		for _, f := range facilities {
			if countryCode != "" && !strings.EqualFold(f.CountryCode, countryCode) {
				continue
			}
			if regionCode != "" && !strings.EqualFold(f.RegionCode, regionCode) {
				continue
			}
			if districtCode != "" && !strings.EqualFold(f.DistrictCode, districtCode) {
				continue
			}
			filtered = append(filtered, f)
		}
		facilities = filtered
		total = len(filtered)
	}

	// 사용자 좌표가 주어진 경우 Haversine 거리 계산
	if userLat != 0 || userLon != 0 {
		for _, f := range facilities {
			f.DistanceKM = Haversine(userLat, userLon, f.Latitude, f.Longitude)
		}
	}

	s.log.Info("시설 검색 완료",
		zap.Int("total", total),
		zap.Int("returned", len(facilities)),
	)

	return facilities, total, nil
}

// GetFacility는 시설을 조회합니다.
func (s *ReservationService) GetFacility(ctx context.Context, facilityID string) (*Facility, error) {
	if facilityID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "facility_id는 필수입니다")
	}
	return s.facilityRepo.FindByID(ctx, facilityID)
}

// ListDoctorsByFacility는 시설별 의사 목록을 조회합니다.
func (s *ReservationService) ListDoctorsByFacility(ctx context.Context, facilityID string, specialty string) ([]*Doctor, error) {
	if facilityID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "facility_id는 필수입니다")
	}
	return s.doctorRepo.ListByFacility(ctx, facilityID, specialty)
}

// GetAvailableSlots는 예약 가능 시간대를 조회합니다.
func (s *ReservationService) GetAvailableSlots(ctx context.Context, facilityID string, date time.Time, doctorID string, specialty Specialty) ([]*TimeSlot, error) {
	if facilityID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "facility_id는 필수입니다")
	}

	slots, err := s.slotRepo.FindByFacilityAndDate(ctx, facilityID, date, doctorID, specialty)
	if err != nil {
		return nil, fmt.Errorf("시간대 조회 실패: %w", err)
	}

	return slots, nil
}

// CreateReservation은 새 예약을 생성합니다.
func (s *ReservationService) CreateReservation(ctx context.Context, userID, facilityID, doctorID, slotID string, specialty Specialty, reason, notes string) (*Reservation, error) {
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}
	if facilityID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "facility_id는 필수입니다")
	}

	// 시설 조회로 이름 가져오기
	facility, err := s.facilityRepo.FindByID(ctx, facilityID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	reservation := &Reservation{
		ID:           uuid.New().String(),
		UserID:       userID,
		FacilityID:   facilityID,
		FacilityName: facility.Name,
		DoctorID:     doctorID,
		DoctorName:   "", // 슬롯에서 가져올 수도 있으나 간략화
		Specialty:    specialty,
		Status:       ResPending,
		Reason:       reason,
		Notes:        notes,
		ScheduledAt:  now.Add(24 * time.Hour), // 기본: 내일
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.reservationRepo.Save(ctx, reservation); err != nil {
		return nil, fmt.Errorf("예약 저장 실패: %w", err)
	}

	s.log.Info("예약 생성 완료",
		zap.String("reservation_id", reservation.ID),
		zap.String("user_id", userID),
		zap.String("facility_id", facilityID),
	)

	// Publish reservation.created event
	if s.eventPub != nil {
		_ = s.eventPub.Publish(ctx, map[string]interface{}{
			"type":           "reservation.created",
			"reservation_id": reservation.ID,
			"user_id":        reservation.UserID,
			"facility_name":  reservation.FacilityName,
			"date":           reservation.ScheduledAt.Format("2006-01-02"),
		})
	}

	return reservation, nil
}

// GetReservation은 예약을 조회합니다.
func (s *ReservationService) GetReservation(ctx context.Context, reservationID string) (*Reservation, error) {
	if reservationID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "reservation_id는 필수입니다")
	}
	return s.reservationRepo.FindByID(ctx, reservationID)
}

// ListReservations는 사용자의 예약 목록을 조회합니다.
func (s *ReservationService) ListReservations(ctx context.Context, userID string, statusFilter ReservationStatus, limit, offset int) ([]*Reservation, int, error) {
	if userID == "" {
		return nil, 0, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}
	if limit <= 0 {
		limit = 20
	}
	return s.reservationRepo.FindByUserID(ctx, userID, statusFilter, limit, offset)
}

// GetDoctorAvailability는 특정 의사의 특정 날짜 예약 가능 시간대를 조회합니다.
func (s *ReservationService) GetDoctorAvailability(ctx context.Context, doctorID string, date time.Time) ([]*TimeSlot, error) {
	if doctorID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "doctor_id는 필수입니다")
	}

	// 의사 정보 조회하여 시설 ID 가져오기
	doctor, err := s.doctorRepo.FindByID(ctx, doctorID)
	if err != nil {
		return nil, err
	}

	slots, err := s.slotRepo.FindByFacilityAndDate(ctx, doctor.FacilityID, date, doctorID, SpecUnknown)
	if err != nil {
		return nil, fmt.Errorf("의사 시간대 조회 실패: %w", err)
	}

	// 예약 가능한 슬롯만 필터링
	var available []*TimeSlot
	for _, slot := range slots {
		if slot.IsAvailable {
			available = append(available, slot)
		}
	}

	return available, nil
}

// SelectDoctor는 시설 내 특정 의사를 선택합니다.
func (s *ReservationService) SelectDoctor(ctx context.Context, facilityID, doctorID, userID string) (*Doctor, error) {
	if facilityID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "facility_id는 필수입니다")
	}
	if doctorID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "doctor_id는 필수입니다")
	}
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}

	doctor, err := s.doctorRepo.FindByID(ctx, doctorID)
	if err != nil {
		return nil, err
	}

	// 의사가 해당 시설에 속하는지 확인
	if doctor.FacilityID != facilityID {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "해당 시설에 소속된 의사가 아닙니다")
	}

	if !doctor.IsAvailable {
		return nil, apperrors.New(apperrors.ErrConflict, "현재 예약 불가능한 의사입니다")
	}

	s.log.Info("의사 선택 완료",
		zap.String("doctor_id", doctorID),
		zap.String("facility_id", facilityID),
		zap.String("user_id", userID),
	)

	return doctor, nil
}

// CancelReservation은 예약을 취소합니다.
func (s *ReservationService) CancelReservation(ctx context.Context, reservationID, userID, reason string) error {
	if reservationID == "" {
		return apperrors.New(apperrors.ErrInvalidInput, "reservation_id는 필수입니다")
	}
	if userID == "" {
		return apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}

	reservation, err := s.reservationRepo.FindByID(ctx, reservationID)
	if err != nil {
		return err
	}

	if reservation.UserID != userID {
		return apperrors.New(apperrors.ErrInvalidInput, "본인의 예약만 취소할 수 있습니다")
	}

	if reservation.Status == ResCancelled {
		return apperrors.New(apperrors.ErrConflict, "이미 취소된 예약입니다")
	}
	if reservation.Status == ResCompleted {
		return apperrors.New(apperrors.ErrConflict, "완료된 예약은 취소할 수 없습니다")
	}

	reservation.Status = ResCancelled
	reservation.Notes = reason
	reservation.UpdatedAt = time.Now()

	if err := s.reservationRepo.Update(ctx, reservation); err != nil {
		return fmt.Errorf("예약 업데이트 실패: %w", err)
	}

	s.log.Info("예약 취소 완료",
		zap.String("reservation_id", reservationID),
		zap.String("user_id", userID),
	)

	// Publish reservation.cancelled event
	if s.eventPub != nil {
		_ = s.eventPub.Publish(ctx, map[string]interface{}{
			"type":           "reservation.cancelled",
			"reservation_id": reservation.ID,
			"user_id":        reservation.UserID,
			"facility_name":  reservation.FacilityName,
			"reason":         reason,
		})
	}

	return nil
}
