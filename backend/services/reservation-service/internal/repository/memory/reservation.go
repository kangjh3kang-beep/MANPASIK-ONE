// Package memory는 reservation-service의 인메모리 저장소입니다.
package memory

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/manpasik/backend/services/reservation-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
)

// ============================================================================
// RegionRepository
// ============================================================================

// RegionRepository는 지역 인메모리 저장소입니다.
type RegionRepository struct {
	mu    sync.RWMutex
	store map[string]*service.Region // regionID -> Region
}

// NewRegionRepository는 시드 데이터가 포함된 인메모리 지역 저장소를 생성합니다.
func NewRegionRepository() *RegionRepository {
	regions := []*service.Region{
		// 한국 (KR)
		{ID: "kr-seoul-gangnam", CountryCode: "KR", RegionCode: "seoul", DistrictCode: "gangnam", Name: "Gangnam, Seoul", NameLocal: "서울 강남구", ParentID: "kr-seoul", Timezone: "Asia/Seoul"},
		{ID: "kr-seoul-songpa", CountryCode: "KR", RegionCode: "seoul", DistrictCode: "songpa", Name: "Songpa, Seoul", NameLocal: "서울 송파구", ParentID: "kr-seoul", Timezone: "Asia/Seoul"},
		{ID: "kr-busan-haeundae", CountryCode: "KR", RegionCode: "busan", DistrictCode: "haeundae", Name: "Haeundae, Busan", NameLocal: "부산 해운대구", ParentID: "kr-busan", Timezone: "Asia/Seoul"},
		{ID: "kr-daejeon-yuseong", CountryCode: "KR", RegionCode: "daejeon", DistrictCode: "yuseong", Name: "Yuseong, Daejeon", NameLocal: "대전 유성구", ParentID: "kr-daejeon", Timezone: "Asia/Seoul"},
		{ID: "kr-jeju-jeju", CountryCode: "KR", RegionCode: "jeju", DistrictCode: "jeju", Name: "Jeju, Jeju", NameLocal: "제주 제주시", ParentID: "kr-jeju", Timezone: "Asia/Seoul"},
		// 일본 (JP)
		{ID: "jp-tokyo-shibuya", CountryCode: "JP", RegionCode: "tokyo", DistrictCode: "shibuya", Name: "Shibuya, Tokyo", NameLocal: "東京都渋谷区", ParentID: "jp-tokyo", Timezone: "Asia/Tokyo"},
		{ID: "jp-osaka-namba", CountryCode: "JP", RegionCode: "osaka", DistrictCode: "namba", Name: "Namba, Osaka", NameLocal: "大阪市浪速区", ParentID: "jp-osaka", Timezone: "Asia/Tokyo"},
		{ID: "jp-kyoto-gion", CountryCode: "JP", RegionCode: "kyoto", DistrictCode: "gion", Name: "Gion, Kyoto", NameLocal: "京都市東山区", ParentID: "jp-kyoto", Timezone: "Asia/Tokyo"},
	}

	store := make(map[string]*service.Region)
	for _, r := range regions {
		store[r.ID] = r
	}
	return &RegionRepository{store: store}
}

// ListByCountry는 국가 코드로 지역 목록을 조회합니다.
func (r *RegionRepository) ListByCountry(_ context.Context, countryCode string) ([]*service.Region, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*service.Region
	for _, region := range r.store {
		if strings.EqualFold(region.CountryCode, countryCode) {
			result = append(result, region)
		}
	}
	return result, nil
}

// ListByRegion는 국가+지역 코드로 구/군 목록을 조회합니다.
func (r *RegionRepository) ListByRegion(_ context.Context, countryCode, regionCode string) ([]*service.Region, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*service.Region
	for _, region := range r.store {
		if strings.EqualFold(region.CountryCode, countryCode) && strings.EqualFold(region.RegionCode, regionCode) {
			result = append(result, region)
		}
	}
	return result, nil
}

// GetByID는 지역 ID로 지역을 조회합니다.
func (r *RegionRepository) GetByID(_ context.Context, regionID string) (*service.Region, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	region, ok := r.store[regionID]
	if !ok {
		return nil, apperrors.New(apperrors.ErrNotFound, "지역을 찾을 수 없습니다")
	}
	return region, nil
}

// ============================================================================
// DoctorRepository
// ============================================================================

type DoctorRepository struct {
	mu    sync.RWMutex
	store map[string]*service.Doctor // doctorID -> Doctor
	byFac map[string][]string        // facilityID -> doctorIDs
}

func NewDoctorRepository() *DoctorRepository {
	tomorrow := time.Now().Truncate(24*time.Hour).Add(24*time.Hour + 9*time.Hour) // 내일 09:00
	doctors := []*service.Doctor{
		{ID: "doc-001", FacilityID: "fac-001", Name: "김내과", Specialty: "internal", LicenseNumber: "M12345", IsAvailable: true, Rating: 4.8, TotalConsultations: 1200, NextAvailableAt: tomorrow, ConsultationFee: 15000, AcceptsTelemedicine: true, AvailableRegionCodes: []string{"kr-seoul-gangnam", "kr-seoul-songpa"}},
		{ID: "doc-002", FacilityID: "fac-001", Name: "이심장", Specialty: "cardiology", LicenseNumber: "M12346", IsAvailable: true, Rating: 4.9, TotalConsultations: 800, NextAvailableAt: tomorrow, ConsultationFee: 25000, AcceptsTelemedicine: false, AvailableRegionCodes: []string{"kr-seoul-gangnam"}},
		{ID: "doc-003", FacilityID: "fac-001", Name: "박피부", Specialty: "dermatology", LicenseNumber: "M12347", IsAvailable: true, Rating: 4.7, TotalConsultations: 600, NextAvailableAt: tomorrow, ConsultationFee: 20000, AcceptsTelemedicine: true, AvailableRegionCodes: []string{"kr-seoul-gangnam", "kr-seoul-songpa"}},
		{ID: "doc-004", FacilityID: "fac-002", Name: "정소아", Specialty: "pediatrics", LicenseNumber: "M12348", IsAvailable: true, Rating: 4.6, TotalConsultations: 1500, NextAvailableAt: tomorrow, ConsultationFee: 12000, AcceptsTelemedicine: true, AvailableRegionCodes: []string{"kr-seoul-gangnam"}},
		{ID: "doc-005", FacilityID: "fac-003", Name: "최일반", Specialty: "general", LicenseNumber: "M12349", IsAvailable: true, Rating: 4.5, TotalConsultations: 2000, NextAvailableAt: tomorrow, ConsultationFee: 10000, AcceptsTelemedicine: false, AvailableRegionCodes: []string{"kr-busan-haeundae"}},
	}
	store := make(map[string]*service.Doctor)
	byFac := make(map[string][]string)
	for _, d := range doctors {
		store[d.ID] = d
		byFac[d.FacilityID] = append(byFac[d.FacilityID], d.ID)
	}
	return &DoctorRepository{store: store, byFac: byFac}
}

func (r *DoctorRepository) ListByFacility(_ context.Context, facilityID string, specialty string) ([]*service.Doctor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids, ok := r.byFac[facilityID]
	if !ok {
		return nil, nil
	}
	var result []*service.Doctor
	for _, id := range ids {
		d, ok := r.store[id]
		if !ok {
			continue
		}
		if specialty != "" && d.Specialty != specialty {
			continue
		}
		result = append(result, d)
	}
	return result, nil
}

// FindByID는 의사 ID로 의사를 조회합니다.
func (r *DoctorRepository) FindByID(_ context.Context, doctorID string) (*service.Doctor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.store[doctorID]
	if !ok {
		return nil, apperrors.New(apperrors.ErrNotFound, "의사를 찾을 수 없습니다")
	}
	return d, nil
}

// ============================================================================
// FacilityRepository
// ============================================================================

// FacilityRepository는 시설 인메모리 저장소입니다.
type FacilityRepository struct {
	mu    sync.RWMutex
	store map[string]*service.Facility
}

// NewFacilityRepository는 샘플 시설 데이터가 포함된 인메모리 저장소를 생성합니다.
func NewFacilityRepository() *FacilityRepository {
	repo := &FacilityRepository{
		store: make(map[string]*service.Facility),
	}

	// 샘플 시설 5곳 시드
	facilities := []*service.Facility{
		{
			ID:                 "fac-001",
			Name:               "서울대학교병원",
			Type:               service.FacilityHospital,
			Address:            "서울특별시 종로구 대학로 101",
			Phone:              "02-2072-2114",
			Latitude:           37.5796,
			Longitude:          127.0040,
			Rating:             4.8,
			ReviewCount:        3200,
			Specialties:        []service.Specialty{service.SpecInternal, service.SpecCardiology, service.SpecEndocrinology, service.SpecOrthopedics},
			OperatingHours:     "평일 09:00-17:30, 토 09:00-12:30",
			IsOpenNow:          true,
			AcceptsReservation: true,
			ImageURL:           "https://img.manpasik.com/facilities/fac-001.jpg",
			CountryCode:        "KR",
			RegionCode:         "seoul",
			DistrictCode:       "gangnam",
			Timezone:           "Asia/Seoul",
			HasTelemedicine:    true,
		},
		{
			ID:                 "fac-002",
			Name:               "강남세브란스병원",
			Type:               service.FacilityHospital,
			Address:            "서울특별시 강남구 언주로 211",
			Phone:              "02-2019-3114",
			Latitude:           37.4976,
			Longitude:          127.0474,
			Rating:             4.7,
			ReviewCount:        2800,
			Specialties:        []service.Specialty{service.SpecInternal, service.SpecDermatology, service.SpecPediatrics},
			OperatingHours:     "평일 08:30-17:30, 토 08:30-12:30",
			IsOpenNow:          true,
			AcceptsReservation: true,
			ImageURL:           "https://img.manpasik.com/facilities/fac-002.jpg",
			CountryCode:        "KR",
			RegionCode:         "seoul",
			DistrictCode:       "gangnam",
			Timezone:           "Asia/Seoul",
			HasTelemedicine:    true,
		},
		{
			ID:                 "fac-003",
			Name:               "미소내과의원",
			Type:               service.FacilityClinic,
			Address:            "서울특별시 서초구 서초대로 256",
			Phone:              "02-555-1234",
			Latitude:           37.4917,
			Longitude:          127.0079,
			Rating:             4.5,
			ReviewCount:        450,
			Specialties:        []service.Specialty{service.SpecInternal, service.SpecGeneral},
			OperatingHours:     "평일 09:00-18:00, 토 09:00-13:00",
			IsOpenNow:          true,
			AcceptsReservation: true,
			ImageURL:           "https://img.manpasik.com/facilities/fac-003.jpg",
			CountryCode:        "KR",
			RegionCode:         "busan",
			DistrictCode:       "haeundae",
			Timezone:           "Asia/Seoul",
			HasTelemedicine:    false,
		},
		{
			ID:                 "fac-004",
			Name:               "건강약국",
			Type:               service.FacilityPharmacy,
			Address:            "서울특별시 강남구 역삼로 123",
			Phone:              "02-555-5678",
			Latitude:           37.5003,
			Longitude:          127.0367,
			Rating:             4.3,
			ReviewCount:        120,
			Specialties:        []service.Specialty{},
			OperatingHours:     "평일 09:00-21:00, 토 09:00-18:00",
			IsOpenNow:          true,
			AcceptsReservation: false,
			ImageURL:           "https://img.manpasik.com/facilities/fac-004.jpg",
			CountryCode:        "KR",
			RegionCode:         "seoul",
			DistrictCode:       "songpa",
			Timezone:           "Asia/Seoul",
			HasTelemedicine:    false,
		},
		{
			ID:                 "fac-005",
			Name:               "밝은눈안과의원",
			Type:               service.FacilityClinic,
			Address:            "서울특별시 마포구 월드컵북로 396",
			Phone:              "02-333-9999",
			Latitude:           37.5563,
			Longitude:          126.9062,
			Rating:             4.6,
			ReviewCount:        680,
			Specialties:        []service.Specialty{service.SpecOphthalmology},
			OperatingHours:     "평일 09:00-18:00, 토 09:00-13:00",
			IsOpenNow:          false,
			AcceptsReservation: true,
			ImageURL:           "https://img.manpasik.com/facilities/fac-005.jpg",
			CountryCode:        "JP",
			RegionCode:         "tokyo",
			DistrictCode:       "shibuya",
			Timezone:           "Asia/Tokyo",
			HasTelemedicine:    true,
		},
	}

	for _, f := range facilities {
		repo.store[f.ID] = f
	}

	return repo
}

// Search는 조건에 맞는 시설을 검색합니다.
func (r *FacilityRepository) Search(_ context.Context, facilityType service.FacilityType, keyword string, specialty service.Specialty, limit, offset int) ([]*service.Facility, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*service.Facility
	for _, f := range r.store {
		// 시설 유형 필터
		if facilityType != service.FacilityUnknown && f.Type != facilityType {
			continue
		}
		// 키워드 필터 (이름 또는 주소)
		if keyword != "" {
			lowerKeyword := strings.ToLower(keyword)
			if !strings.Contains(strings.ToLower(f.Name), lowerKeyword) &&
				!strings.Contains(strings.ToLower(f.Address), lowerKeyword) {
				continue
			}
		}
		// 전문 분야 필터
		if specialty != service.SpecUnknown {
			found := false
			for _, sp := range f.Specialties {
				if sp == specialty {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		filtered = append(filtered, f)
	}

	// 평점 내림차순 정렬
	for i := 0; i < len(filtered); i++ {
		for j := i + 1; j < len(filtered); j++ {
			if filtered[j].Rating > filtered[i].Rating {
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

// FindByID는 ID로 시설을 조회합니다.
func (r *FacilityRepository) FindByID(_ context.Context, id string) (*service.Facility, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	f, ok := r.store[id]
	if !ok {
		return nil, apperrors.New(apperrors.ErrNotFound, "시설을 찾을 수 없습니다")
	}
	return f, nil
}

// ============================================================================
// SlotRepository
// ============================================================================

// SlotRepository는 시간대 인메모리 저장소입니다.
type SlotRepository struct {
	mu    sync.RWMutex
	store map[string][]*service.TimeSlot // facilityID -> slots
}

// NewSlotRepository는 샘플 시간대 데이터가 포함된 인메모리 저장소를 생성합니다.
func NewSlotRepository() *SlotRepository {
	repo := &SlotRepository{
		store: make(map[string][]*service.TimeSlot),
	}

	// 시설별 시간대 생성
	doctors := []struct {
		id   string
		name string
	}{
		{"doc-001", "김내과"},
		{"doc-002", "이심장"},
		{"doc-003", "박피부"},
	}

	// fac-001 시간대 생성
	baseDate := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour) // 내일
	for i, doc := range doctors {
		for hour := 9; hour < 17; hour++ {
			slotStart := baseDate.Add(time.Duration(hour) * time.Hour)
			slotEnd := slotStart.Add(30 * time.Minute)
			slot := &service.TimeSlot{
				ID:          fmt.Sprintf("slot-001-%d-%02d", i+1, hour),
				StartTime:   slotStart,
				EndTime:     slotEnd,
				IsAvailable: hour%3 != 0, // 매 3시간째는 이미 예약됨
				DoctorID:    doc.id,
				DoctorName:  doc.name,
			}
			repo.store["fac-001"] = append(repo.store["fac-001"], slot)
		}
	}

	// fac-002 시간대 생성
	for hour := 9; hour < 17; hour++ {
		slotStart := baseDate.Add(time.Duration(hour) * time.Hour)
		slotEnd := slotStart.Add(30 * time.Minute)
		slot := &service.TimeSlot{
			ID:          fmt.Sprintf("slot-002-1-%02d", hour),
			StartTime:   slotStart,
			EndTime:     slotEnd,
			IsAvailable: true,
			DoctorID:    "doc-004",
			DoctorName:  "정소아",
		}
		repo.store["fac-002"] = append(repo.store["fac-002"], slot)
	}

	// fac-003 시간대 생성
	for hour := 9; hour < 18; hour++ {
		slotStart := baseDate.Add(time.Duration(hour) * time.Hour)
		slotEnd := slotStart.Add(30 * time.Minute)
		slot := &service.TimeSlot{
			ID:          fmt.Sprintf("slot-003-1-%02d", hour),
			StartTime:   slotStart,
			EndTime:     slotEnd,
			IsAvailable: true,
			DoctorID:    "doc-005",
			DoctorName:  "최일반",
		}
		repo.store["fac-003"] = append(repo.store["fac-003"], slot)
	}

	return repo
}

// FindByFacilityAndDate는 시설 ID와 날짜로 시간대를 조회합니다.
func (r *SlotRepository) FindByFacilityAndDate(_ context.Context, facilityID string, date time.Time, doctorID string, specialty service.Specialty) ([]*service.TimeSlot, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	slots, ok := r.store[facilityID]
	if !ok {
		return nil, nil // 슬롯이 없으면 빈 리스트 반환
	}

	var result []*service.TimeSlot
	for _, slot := range slots {
		// 의사 ID 필터
		if doctorID != "" && slot.DoctorID != doctorID {
			continue
		}
		result = append(result, slot)
	}

	return result, nil
}

// ============================================================================
// ReservationRepository
// ============================================================================

// ReservationRepository는 예약 인메모리 저장소입니다.
type ReservationRepository struct {
	mu    sync.RWMutex
	store map[string]*service.Reservation
}

// NewReservationRepository는 새 인메모리 예약 저장소를 생성합니다.
func NewReservationRepository() *ReservationRepository {
	return &ReservationRepository{
		store: make(map[string]*service.Reservation),
	}
}

// Save는 예약을 저장합니다.
func (r *ReservationRepository) Save(_ context.Context, res *service.Reservation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[res.ID] = res
	return nil
}

// FindByID는 ID로 예약을 조회합니다.
func (r *ReservationRepository) FindByID(_ context.Context, id string) (*service.Reservation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res, ok := r.store[id]
	if !ok {
		return nil, apperrors.New(apperrors.ErrNotFound, "예약을 찾을 수 없습니다")
	}
	return res, nil
}

// FindByUserID는 사용자 예약 목록을 조회합니다.
func (r *ReservationRepository) FindByUserID(_ context.Context, userID string, statusFilter service.ReservationStatus, limit, offset int) ([]*service.Reservation, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*service.Reservation
	for _, res := range r.store {
		if res.UserID != userID {
			continue
		}
		if statusFilter != service.ResUnknown && res.Status != statusFilter {
			continue
		}
		filtered = append(filtered, res)
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

// Update는 예약을 업데이트합니다.
func (r *ReservationRepository) Update(_ context.Context, res *service.Reservation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[res.ID]; !ok {
		return apperrors.New(apperrors.ErrNotFound, "예약을 찾을 수 없습니다")
	}
	r.store[res.ID] = res
	return nil
}
