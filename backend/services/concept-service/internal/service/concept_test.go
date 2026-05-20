package service_test

import (
	"context"
	"testing"

	"github.com/manpasik/backend/services/concept-service/internal/repository/memory"
	"github.com/manpasik/backend/services/concept-service/internal/service"
	"go.uber.org/zap"
)

func newSvc() *service.ConceptService {
	return service.NewConceptService(zap.NewNop(), memory.NewConceptRepository(), memory.NewOrganizationRepository())
}

func TestCreateConcept(t *testing.T) {
	s := newSvc()
	c, err := s.CreateConcept(context.Background(), "테스트", "설명", "personal", "", "u1")
	if err != nil {
		t.Fatal(err)
	}
	if c.ID == "" || c.Name != "테스트" {
		t.Error("unexpected concept")
	}
}
func TestListConcepts(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	s.CreateConcept(ctx, "AA", "", "", "", "u1")
	s.CreateConcept(ctx, "BB", "", "", "", "u1")
	list, _ := s.ListConcepts(ctx)
	if len(list) != 2 {
		t.Errorf("expected 2, got %d", len(list))
	}
}
func TestAssignDevice(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	c, _ := s.CreateConcept(ctx, "XX", "", "", "", "u1")
	if err := s.AssignDevice(ctx, c.ID, "dev1"); err != nil {
		t.Fatal(err)
	}
	got, _ := s.GetConcept(ctx, c.ID)
	if len(got.DeviceIDs) != 1 {
		t.Error("device not assigned")
	}
}
func TestCreateOrg(t *testing.T) {
	s := newSvc()
	o, err := s.CreateOrganization(context.Background(), "병원", "테스트")
	if err != nil {
		t.Fatal(err)
	}
	if o.ID == "" {
		t.Error("empty id")
	}
}
func TestAddRemoveMember(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	o, _ := s.CreateOrganization(ctx, "병원", "")
	s.AddMember(ctx, o.ID, "m1")
	got, _ := s.GetOrganization(ctx, o.ID)
	if len(got.MemberIDs) != 1 {
		t.Error("member not added")
	}
	s.RemoveMember(ctx, o.ID, "m1")
	got2, _ := s.GetOrganization(ctx, o.ID)
	if len(got2.MemberIDs) != 0 {
		t.Error("member not removed")
	}
}

// ============================================================================
// Phase B 신규 테스트
// ============================================================================

func TestCreateConcept_InvalidCategory(t *testing.T) {
	s := newSvc()
	_, err := s.CreateConcept(context.Background(), "테스트", "", "invalid_cat", "", "u1")
	if err == nil {
		t.Error("지원하지 않는 카테고리에 에러가 반환되어야 합니다")
	}
}

func TestCreateConcept_DefaultCategory(t *testing.T) {
	s := newSvc()
	c, err := s.CreateConcept(context.Background(), "기본카테고리", "", "", "", "u1")
	if err != nil {
		t.Fatal(err)
	}
	if c.Category != "personal" {
		t.Errorf("기본 카테고리 = %s, want personal", c.Category)
	}
}

func TestCreateConcept_NameLength(t *testing.T) {
	s := newSvc()
	// 1글자 이름 (too short)
	_, err := s.CreateConcept(context.Background(), "A", "", "personal", "", "u1")
	if err == nil {
		t.Error("너무 짧은 이름에 에러가 반환되어야 합니다")
	}
}

func TestAssignDevice_Limit(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	c, _ := s.CreateConcept(ctx, "리밋테스트", "", "personal", "", "u1")

	// 20개까지 등록
	for i := 0; i < service.MaxDevicesPerConcept; i++ {
		err := s.AssignDevice(ctx, c.ID, "dev-"+string(rune('A'+i)))
		if err != nil {
			t.Fatalf("디바이스 %d 등록 실패: %v", i, err)
		}
	}

	// 21번째 → 에러
	err := s.AssignDevice(ctx, c.ID, "dev-overflow")
	if err == nil {
		t.Error("디바이스 제한 초과 시 에러가 반환되어야 합니다")
	}
}

func TestUnassignDevice(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	c, _ := s.CreateConcept(ctx, "언어사인", "", "personal", "", "u1")
	s.AssignDevice(ctx, c.ID, "dev1")
	s.AssignDevice(ctx, c.ID, "dev2")

	err := s.UnassignDevice(ctx, c.ID, "dev1")
	if err != nil {
		t.Fatalf("UnassignDevice 실패: %v", err)
	}

	got, _ := s.GetConcept(ctx, c.ID)
	if len(got.DeviceIDs) != 1 {
		t.Errorf("디바이스 수 = %d, want 1", len(got.DeviceIDs))
	}

	// 존재하지 않는 디바이스 제거
	err = s.UnassignDevice(ctx, c.ID, "dev-nonexistent")
	if err == nil {
		t.Error("존재하지 않는 디바이스 제거에 에러가 반환되어야 합니다")
	}
}

func TestDeleteConcept(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	c, _ := s.CreateConcept(ctx, "삭제대상", "", "personal", "", "u1")

	err := s.DeleteConcept(ctx, c.ID)
	if err != nil {
		t.Fatalf("DeleteConcept 실패: %v", err)
	}

	_, err = s.GetConcept(ctx, c.ID)
	if err == nil {
		t.Error("삭제된 컨셉 조회에 에러가 반환되어야 합니다")
	}
}

// =============================================================================
// Phase F 테스트 보강
// =============================================================================

func TestGetConcept_EmptyID(t *testing.T) {
	s := newSvc()
	_, err := s.GetConcept(context.Background(), "")
	if err == nil {
		t.Error("빈 ID에 에러가 반환되어야 합니다")
	}
}

func TestGetConcept_NotFound(t *testing.T) {
	s := newSvc()
	_, err := s.GetConcept(context.Background(), "nonexistent")
	if err == nil {
		t.Error("존재하지 않는 컨셉에 에러가 반환되어야 합니다")
	}
}

func TestGetConceptStats(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	c, _ := s.CreateConcept(ctx, "통계테스트", "", "personal", "", "u1")
	s.AssignDevice(ctx, c.ID, "d1")
	s.AssignDevice(ctx, c.ID, "d2")

	total, active, rate, err := s.GetConceptStats(ctx, c.ID)
	if err != nil {
		t.Fatal(err)
	}
	if total != 2 {
		t.Errorf("total: got %d, want 2", total)
	}
	if active != 2 {
		t.Errorf("active: got %d, want 2", active)
	}
	if rate <= 0 {
		t.Errorf("rate should be > 0: got %f", rate)
	}
}

func TestGetConceptStats_EmptyID(t *testing.T) {
	s := newSvc()
	_, _, _, err := s.GetConceptStats(context.Background(), "")
	if err == nil {
		t.Error("빈 ID에 에러가 반환되어야 합니다")
	}
}

func TestGetConceptDashboard(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	c, _ := s.CreateConcept(ctx, "대시보드", "", "fitness", "", "u1")
	s.AssignDevice(ctx, c.ID, "d1")

	total, active, rate, devices, err := s.GetConceptDashboard(ctx, c.ID)
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || active != 1 {
		t.Errorf("total=%d active=%d, want 1,1", total, active)
	}
	if rate <= 0 {
		t.Errorf("rate should be > 0: got %f", rate)
	}
	if len(devices) != 1 || devices[0] != "d1" {
		t.Errorf("devices: got %v, want [d1]", devices)
	}
}

func TestGetConceptDashboard_EmptyID(t *testing.T) {
	s := newSvc()
	_, _, _, _, err := s.GetConceptDashboard(context.Background(), "")
	if err == nil {
		t.Error("빈 ID에 에러가 반환되어야 합니다")
	}
}

func TestListOrganizations(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	s.CreateOrganization(ctx, "조직1", "")
	s.CreateOrganization(ctx, "조직2", "")

	list, err := s.ListOrganizations(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2, got %d", len(list))
	}
}

func TestDeleteOrganization(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	o, _ := s.CreateOrganization(ctx, "삭제대상", "")

	err := s.DeleteOrganization(ctx, o.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteOrganization_EmptyID(t *testing.T) {
	s := newSvc()
	err := s.DeleteOrganization(context.Background(), "")
	if err == nil {
		t.Error("빈 ID에 에러가 반환되어야 합니다")
	}
}

func TestAssignDevice_EmptyParams(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	err := s.AssignDevice(ctx, "", "d1")
	if err == nil {
		t.Error("빈 conceptID에 에러 반환 필요")
	}
	c, _ := s.CreateConcept(ctx, "테스트2", "", "personal", "", "u1")
	err = s.AssignDevice(ctx, c.ID, "")
	if err == nil {
		t.Error("빈 deviceID에 에러 반환 필요")
	}
}

func TestAssignDevice_Duplicate(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	c, _ := s.CreateConcept(ctx, "중복테스트", "", "personal", "", "u1")
	s.AssignDevice(ctx, c.ID, "d1")
	err := s.AssignDevice(ctx, c.ID, "d1")
	if err == nil {
		t.Error("중복 디바이스에 에러가 반환되어야 합니다")
	}
}

func TestAddMember_Duplicate(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	o, _ := s.CreateOrganization(ctx, "중복멤버", "")
	s.AddMember(ctx, o.ID, "m1")
	err := s.AddMember(ctx, o.ID, "m1")
	if err == nil {
		t.Error("중복 멤버에 에러가 반환되어야 합니다")
	}
}

func TestRemoveMember_NotFound(t *testing.T) {
	s := newSvc()
	ctx := context.Background()
	o, _ := s.CreateOrganization(ctx, "없는멤버", "")
	err := s.RemoveMember(ctx, o.ID, "nonexistent")
	if err == nil {
		t.Error("존재하지 않는 멤버 제거에 에러가 반환되어야 합니다")
	}
}

func TestCreateOrganization_ShortName(t *testing.T) {
	s := newSvc()
	_, err := s.CreateOrganization(context.Background(), "A", "")
	if err == nil {
		t.Error("너무 짧은 이름에 에러가 반환되어야 합니다")
	}
}

func TestGetOrganization_EmptyID(t *testing.T) {
	s := newSvc()
	_, err := s.GetOrganization(context.Background(), "")
	if err == nil {
		t.Error("빈 ID에 에러가 반환되어야 합니다")
	}
}

func TestDeleteConcept_EmptyID(t *testing.T) {
	s := newSvc()
	err := s.DeleteConcept(context.Background(), "")
	if err == nil {
		t.Error("빈 ID에 에러가 반환되어야 합니다")
	}
}
