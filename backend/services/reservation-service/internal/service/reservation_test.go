package service_test

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/manpasik/backend/services/reservation-service/internal/repository/memory"
	"github.com/manpasik/backend/services/reservation-service/internal/service"
	"go.uber.org/zap"
)

func setupTestService() *service.ReservationService {
	logger := zap.NewNop()
	facilityRepo := memory.NewFacilityRepository()
	slotRepo := memory.NewSlotRepository()
	reservationRepo := memory.NewReservationRepository()
	doctorRepo := memory.NewDoctorRepository()
	regionRepo := memory.NewRegionRepository()
	return service.NewReservationService(logger, facilityRepo, slotRepo, reservationRepo, doctorRepo, regionRepo)
}

func TestSearchFacilities_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	facilities, total, err := svc.SearchFacilities(ctx, service.FacilityUnknown, "", service.SpecUnknown, 20, 0, "", "", "", 0, 0)
	if err != nil {
		t.Fatalf("시설 검색 실패: %v", err)
	}
	if total == 0 {
		t.Fatal("시설이 최소 1개는 있어야 함")
	}
	if len(facilities) == 0 {
		t.Fatal("반환된 시설이 없음")
	}
}

func TestSearchFacilities_ByType(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	facilities, _, err := svc.SearchFacilities(ctx, service.FacilityHospital, "", service.SpecUnknown, 20, 0, "", "", "", 0, 0)
	if err != nil {
		t.Fatalf("시설 검색 실패: %v", err)
	}
	for _, f := range facilities {
		if f.Type != service.FacilityHospital {
			t.Fatalf("시설 유형 불일치: got %d, want %d", f.Type, service.FacilityHospital)
		}
	}
}

func TestSearchFacilities_ByKeyword(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	facilities, total, err := svc.SearchFacilities(ctx, service.FacilityUnknown, "서울대", service.SpecUnknown, 20, 0, "", "", "", 0, 0)
	if err != nil {
		t.Fatalf("시설 검색 실패: %v", err)
	}
	if total == 0 {
		t.Fatal("'서울대' 키워드로 검색 시 결과가 있어야 함")
	}
	if len(facilities) == 0 {
		t.Fatal("반환된 시설이 없음")
	}
}

func TestGetFacility_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	facility, err := svc.GetFacility(ctx, "fac-001")
	if err != nil {
		t.Fatalf("시설 조회 실패: %v", err)
	}
	if facility.ID != "fac-001" {
		t.Fatalf("시설 ID 불일치: got %s, want fac-001", facility.ID)
	}
	if facility.Name != "서울대학교병원" {
		t.Fatalf("시설 이름 불일치: got %s", facility.Name)
	}
}

func TestGetFacility_NotFound(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, err := svc.GetFacility(ctx, "non-existent-id")
	if err == nil {
		t.Fatal("존재하지 않는 시설에 에러가 반환되어야 함")
	}
}

func TestGetAvailableSlots_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	slots, err := svc.GetAvailableSlots(ctx, "fac-001", time.Time{}, "", service.SpecUnknown)
	if err != nil {
		t.Fatalf("시간대 조회 실패: %v", err)
	}
	if len(slots) == 0 {
		t.Fatal("fac-001에 시간대가 있어야 함")
	}
}

func TestCreateReservation_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	res, err := svc.CreateReservation(ctx, "user-1", "fac-001", "doc-001", "slot-001-1-10", service.SpecInternal, "두통", "오전 진료 희망")
	if err != nil {
		t.Fatalf("예약 생성 실패: %v", err)
	}
	if res.ID == "" {
		t.Fatal("예약 ID가 비어 있음")
	}
	if res.UserID != "user-1" {
		t.Fatalf("UserID 불일치: got %s, want user-1", res.UserID)
	}
	if res.FacilityID != "fac-001" {
		t.Fatalf("FacilityID 불일치: got %s, want fac-001", res.FacilityID)
	}
	if res.FacilityName != "서울대학교병원" {
		t.Fatalf("FacilityName 불일치: got %s", res.FacilityName)
	}
	if res.Status != service.ResPending {
		t.Fatalf("Status 불일치: got %d, want %d", res.Status, service.ResPending)
	}
}

func TestCreateReservation_MissingUserID(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, err := svc.CreateReservation(ctx, "", "fac-001", "doc-001", "slot-1", service.SpecGeneral, "이유", "")
	if err == nil {
		t.Fatal("빈 user_id에 에러가 반환되어야 함")
	}
}

func TestCreateReservation_MissingFacilityID(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, err := svc.CreateReservation(ctx, "user-1", "", "doc-001", "slot-1", service.SpecGeneral, "이유", "")
	if err == nil {
		t.Fatal("빈 facility_id에 에러가 반환되어야 함")
	}
}

func TestGetReservation_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	created, _ := svc.CreateReservation(ctx, "user-2", "fac-002", "doc-004", "slot-002-1-10", service.SpecPediatrics, "소아 진료", "")

	got, err := svc.GetReservation(ctx, created.ID)
	if err != nil {
		t.Fatalf("예약 조회 실패: %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("ID 불일치: got %s, want %s", got.ID, created.ID)
	}
	if got.UserID != "user-2" {
		t.Fatalf("UserID 불일치: got %s", got.UserID)
	}
}

func TestListReservations_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		_, err := svc.CreateReservation(ctx, "user-list", "fac-001", "doc-001", "slot-001-1-10", service.SpecInternal, "이유", "")
		if err != nil {
			t.Fatalf("예약 생성 실패: %v", err)
		}
	}

	reservations, total, err := svc.ListReservations(ctx, "user-list", service.ResUnknown, 10, 0)
	if err != nil {
		t.Fatalf("목록 조회 실패: %v", err)
	}
	if total != 3 {
		t.Fatalf("총 개수 불일치: got %d, want 3", total)
	}
	if len(reservations) != 3 {
		t.Fatalf("반환 개수 불일치: got %d, want 3", len(reservations))
	}
}

func TestCancelReservation_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	created, _ := svc.CreateReservation(ctx, "user-cancel", "fac-001", "doc-001", "slot-001-1-10", service.SpecInternal, "진료", "")

	err := svc.CancelReservation(ctx, created.ID, "user-cancel", "일정 변경")
	if err != nil {
		t.Fatalf("예약 취소 실패: %v", err)
	}

	// 취소 후 상태 확인
	cancelled, _ := svc.GetReservation(ctx, created.ID)
	if cancelled.Status != service.ResCancelled {
		t.Fatalf("Status 불일치: got %d, want %d", cancelled.Status, service.ResCancelled)
	}
}

func TestEndToEnd_ReservationFlow(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	// 1. 시설 검색
	facilities, total, err := svc.SearchFacilities(ctx, service.FacilityHospital, "", service.SpecInternal, 20, 0, "", "", "", 0, 0)
	if err != nil {
		t.Fatalf("시설 검색 실패: %v", err)
	}
	if total == 0 {
		t.Fatal("내과 전문 병원이 최소 1개는 있어야 함")
	}

	targetFacility := facilities[0]

	// 2. 시설 상세 조회
	facility, err := svc.GetFacility(ctx, targetFacility.ID)
	if err != nil {
		t.Fatalf("시설 조회 실패: %v", err)
	}
	if facility.Name == "" {
		t.Fatal("시설 이름이 비어 있음")
	}

	// 3. 예약 가능 시간대 조회
	slots, err := svc.GetAvailableSlots(ctx, targetFacility.ID, time.Time{}, "", service.SpecUnknown)
	if err != nil {
		t.Fatalf("시간대 조회 실패: %v", err)
	}
	// 슬롯이 있을 수도 있고 없을 수도 있음

	_ = slots // 사용

	// 4. 예약 생성
	reservation, err := svc.CreateReservation(ctx, "user-e2e", targetFacility.ID, "doc-001", "slot-001-1-10", service.SpecInternal, "건강검진", "오전 예약 희망")
	if err != nil {
		t.Fatalf("예약 생성 실패: %v", err)
	}
	if reservation.Status != service.ResPending {
		t.Fatalf("초기 상태 불일치: got %d, want %d", reservation.Status, service.ResPending)
	}

	// 5. 예약 조회
	got, err := svc.GetReservation(ctx, reservation.ID)
	if err != nil {
		t.Fatalf("예약 조회 실패: %v", err)
	}
	if got.FacilityName != targetFacility.Name {
		t.Fatalf("FacilityName 불일치: got %s, want %s", got.FacilityName, targetFacility.Name)
	}

	// 6. 목록 조회
	list, listTotal, err := svc.ListReservations(ctx, "user-e2e", service.ResUnknown, 10, 0)
	if err != nil {
		t.Fatalf("목록 조회 실패: %v", err)
	}
	if listTotal != 1 {
		t.Fatalf("총 개수 불일치: got %d, want 1", listTotal)
	}
	if len(list) != 1 {
		t.Fatalf("반환 개수 불일치: got %d, want 1", len(list))
	}

	// 7. 예약 취소
	err = svc.CancelReservation(ctx, reservation.ID, "user-e2e", "개인 사정")
	if err != nil {
		t.Fatalf("예약 취소 실패: %v", err)
	}

	// 8. 취소 후 상태 확인
	cancelled, _ := svc.GetReservation(ctx, reservation.ID)
	if cancelled.Status != service.ResCancelled {
		t.Fatalf("취소 상태 불일치: got %d, want %d", cancelled.Status, service.ResCancelled)
	}
}

// ============================================================================
// 지역 검색 및 의사 프로필 확장 테스트
// ============================================================================

func TestHaversine(t *testing.T) {
	// 서울 (37.5665, 126.9780) → 부산 (35.1796, 129.0756) ≈ 325km
	seoulLat, seoulLon := 37.5665, 126.9780
	busanLat, busanLon := 35.1796, 129.0756

	dist := service.Haversine(seoulLat, seoulLon, busanLat, busanLon)

	// 허용 오차 ±10km
	if math.Abs(dist-325) > 10 {
		t.Fatalf("서울→부산 거리 오류: got %.1f km, want ≈325 km", dist)
	}

	// 같은 지점 거리 = 0
	zero := service.Haversine(seoulLat, seoulLon, seoulLat, seoulLon)
	if zero != 0 {
		t.Fatalf("같은 지점 거리는 0이어야 함: got %.4f", zero)
	}
}

func TestSearchFacilitiesByRegion(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	// KR + seoul 으로 필터링: fac-001, fac-002, fac-004가 매칭되어야 함
	facilities, total, err := svc.SearchFacilities(ctx, service.FacilityUnknown, "", service.SpecUnknown, 20, 0, "KR", "seoul", "", 0, 0)
	if err != nil {
		t.Fatalf("지역 검색 실패: %v", err)
	}
	if total == 0 {
		t.Fatal("KR+seoul로 검색 시 결과가 있어야 함")
	}
	for _, f := range facilities {
		if f.CountryCode != "KR" {
			t.Fatalf("국가 코드 불일치: got %s, want KR", f.CountryCode)
		}
		if f.RegionCode != "seoul" {
			t.Fatalf("지역 코드 불일치: got %s, want seoul", f.RegionCode)
		}
	}

	// JP 필터링: fac-005만 매칭
	jpFacilities, jpTotal, err := svc.SearchFacilities(ctx, service.FacilityUnknown, "", service.SpecUnknown, 20, 0, "JP", "", "", 0, 0)
	if err != nil {
		t.Fatalf("JP 검색 실패: %v", err)
	}
	if jpTotal != 1 {
		t.Fatalf("JP 시설 개수 불일치: got %d, want 1", jpTotal)
	}
	if jpFacilities[0].ID != "fac-005" {
		t.Fatalf("JP 시설 ID 불일치: got %s, want fac-005", jpFacilities[0].ID)
	}
}

func TestSearchFacilitiesByDistrict(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	// KR + seoul + gangnam: fac-001, fac-002
	facilities, total, err := svc.SearchFacilities(ctx, service.FacilityUnknown, "", service.SpecUnknown, 20, 0, "KR", "seoul", "gangnam", 0, 0)
	if err != nil {
		t.Fatalf("구/군 검색 실패: %v", err)
	}
	if total == 0 {
		t.Fatal("KR+seoul+gangnam으로 검색 시 결과가 있어야 함")
	}
	for _, f := range facilities {
		if f.DistrictCode != "gangnam" {
			t.Fatalf("구/군 코드 불일치: got %s, want gangnam", f.DistrictCode)
		}
	}

	// 존재하지 않는 구/군
	_, noTotal, err := svc.SearchFacilities(ctx, service.FacilityUnknown, "", service.SpecUnknown, 20, 0, "KR", "seoul", "nonexistent", 0, 0)
	if err != nil {
		t.Fatalf("빈 결과 검색 실패: %v", err)
	}
	if noTotal != 0 {
		t.Fatalf("존재하지 않는 구/군 결과 개수 불일치: got %d, want 0", noTotal)
	}
}

func TestGetDoctorAvailability(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	slots, err := svc.GetDoctorAvailability(ctx, "doc-001", time.Time{})
	if err != nil {
		t.Fatalf("의사 가용성 조회 실패: %v", err)
	}
	if len(slots) == 0 {
		t.Fatal("doc-001에 가용 시간대가 있어야 함")
	}
	for _, slot := range slots {
		if !slot.IsAvailable {
			t.Fatal("반환된 슬롯은 모두 가용 상태여야 함")
		}
		if slot.DoctorID != "doc-001" {
			t.Fatalf("의사 ID 불일치: got %s, want doc-001", slot.DoctorID)
		}
	}
}

func TestSelectDoctor(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	// 정상 선택
	doctor, err := svc.SelectDoctor(ctx, "fac-001", "doc-001", "user-1")
	if err != nil {
		t.Fatalf("의사 선택 실패: %v", err)
	}
	if doctor.ID != "doc-001" {
		t.Fatalf("의사 ID 불일치: got %s, want doc-001", doctor.ID)
	}
	if doctor.FacilityID != "fac-001" {
		t.Fatalf("시설 ID 불일치: got %s, want fac-001", doctor.FacilityID)
	}

	// 다른 시설의 의사 선택 시도 → 에러
	_, err = svc.SelectDoctor(ctx, "fac-002", "doc-001", "user-1")
	if err == nil {
		t.Fatal("다른 시설의 의사 선택 시 에러가 반환되어야 함")
	}

	// 빈 파라미터 검증
	_, err = svc.SelectDoctor(ctx, "", "doc-001", "user-1")
	if err == nil {
		t.Fatal("빈 facility_id에 에러가 반환되어야 함")
	}
	_, err = svc.SelectDoctor(ctx, "fac-001", "", "user-1")
	if err == nil {
		t.Fatal("빈 doctor_id에 에러가 반환되어야 함")
	}
}
