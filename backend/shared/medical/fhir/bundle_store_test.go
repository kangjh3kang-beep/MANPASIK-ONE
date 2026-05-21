package fhir_test

import (
	"context"
	"errors"
	"testing"

	"github.com/manpasik/backend/shared/medical/fhir"
)

func bundleWith(id, patient string) *fhir.Bundle {
	b := fhir.NewBundle("collection")
	b.ID = id
	return b
}

func TestMemoryBundleStore_SaveAndGet(t *testing.T) {
	store := fhir.NewMemoryBundleStore(0)
	ctx := context.Background()
	b := bundleWith("bundle-1", "PAT1")

	if err := store.Save(ctx, b, "tenantA", "PAT1"); err != nil {
		t.Fatal(err)
	}
	got, err := store.Get(ctx, "bundle-1", "tenantA")
	if err != nil {
		t.Fatal(err)
	}
	if got.Bundle.ID != "bundle-1" || got.PatientID != "PAT1" {
		t.Errorf("got = %+v", got)
	}
	if got.StoredAt.IsZero() {
		t.Error("StoredAt 미설정")
	}
}

func TestMemoryBundleStore_Idempotent(t *testing.T) {
	store := fhir.NewMemoryBundleStore(0)
	ctx := context.Background()
	b1 := bundleWith("X", "P1")
	if err := store.Save(ctx, b1, "T", "P1"); err != nil {
		t.Fatal(err)
	}
	firstStoredAt, _ := store.Get(ctx, "X", "T")

	// 같은 ID 로 두 번째 저장 — 새 entry 가 만들어지지 않고 update.
	b2 := bundleWith("X", "P1")
	if err := store.Save(ctx, b2, "T", "P1-updated"); err != nil {
		t.Fatal(err)
	}
	if store.Count(ctx, "T") != 1 {
		t.Errorf("count = %d, 1 기대 (멱등)", store.Count(ctx, "T"))
	}
	second, _ := store.Get(ctx, "X", "T")
	if second.PatientID != "P1-updated" {
		t.Errorf("PatientID = %s, P1-updated 기대", second.PatientID)
	}
	if !second.StoredAt.Equal(firstStoredAt.StoredAt) {
		t.Error("StoredAt 이 최초 저장 시점에서 변경됨 — 멱등 update 시 보존되어야 함")
	}
}

func TestMemoryBundleStore_TenantIsolation(t *testing.T) {
	store := fhir.NewMemoryBundleStore(0)
	ctx := context.Background()

	_ = store.Save(ctx, bundleWith("shared-id", "P1"), "tenantA", "P1")
	_ = store.Save(ctx, bundleWith("shared-id", "P2"), "tenantB", "P2")

	// tenantA 가 tenantB 의 Bundle 을 조회할 수 없음 (서로 다른 namespace)
	a, _ := store.Get(ctx, "shared-id", "tenantA")
	if a.PatientID != "P1" {
		t.Errorf("tenant 격리 위반: A.PatientID = %s", a.PatientID)
	}
	b, _ := store.Get(ctx, "shared-id", "tenantB")
	if b.PatientID != "P2" {
		t.Errorf("tenant 격리 위반: B.PatientID = %s", b.PatientID)
	}
	if a == b {
		t.Error("tenant 별 격리 실패: 같은 객체 참조")
	}
}

func TestMemoryBundleStore_GetNotFound(t *testing.T) {
	store := fhir.NewMemoryBundleStore(0)
	_, err := store.Get(context.Background(), "missing", "tenantA")
	if !errors.Is(err, fhir.ErrBundleNotFound) {
		t.Errorf("err = %v, ErrBundleNotFound 기대", err)
	}
}

func TestMemoryBundleStore_ListByPatient(t *testing.T) {
	store := fhir.NewMemoryBundleStore(0)
	ctx := context.Background()

	_ = store.Save(ctx, bundleWith("b1", "P1"), "T", "P1")
	_ = store.Save(ctx, bundleWith("b2", "P1"), "T", "P1")
	_ = store.Save(ctx, bundleWith("b3", "P2"), "T", "P2")

	list, err := store.ListByPatient(ctx, "T", "P1", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Errorf("P1 list = %d, 2 기대", len(list))
	}
	for _, e := range list {
		if e.PatientID != "P1" {
			t.Errorf("P1 필터 위반: %s", e.PatientID)
		}
	}
}

func TestMemoryBundleStore_ListByTenant_SortedNewestFirst(t *testing.T) {
	store := fhir.NewMemoryBundleStore(0)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		_ = store.Save(ctx, bundleWith("bundle-"+rune2str(i), "P"), "T", "P")
	}
	list, _ := store.ListByTenant(ctx, "T", 0)
	if len(list) != 3 {
		t.Fatalf("list = %d", len(list))
	}
	// 최신순 (storedAt 내림차순) — 정확한 시간 차이가 작아도 깨지지 않게 순서가 일관되는지 본다.
	for i := 0; i < len(list)-1; i++ {
		if list[i].StoredAt.Before(list[i+1].StoredAt) {
			t.Errorf("정렬 위반: list[%d].StoredAt < list[%d].StoredAt", i, i+1)
		}
	}
}

func TestMemoryBundleStore_LimitClamping(t *testing.T) {
	store := fhir.NewMemoryBundleStore(0)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		_ = store.Save(ctx, bundleWith("b-"+rune2str(i), "P"), "T", "P")
	}
	list, _ := store.ListByTenant(ctx, "T", 2)
	if len(list) != 2 {
		t.Errorf("limit=2 적용 안 됨: %d", len(list))
	}
	listAll, _ := store.ListByTenant(ctx, "T", 0)
	if len(listAll) != 5 {
		t.Errorf("limit=0 (전체) = %d, 5 기대", len(listAll))
	}
}

func TestMemoryBundleStore_MaxSizeFIFO(t *testing.T) {
	store := fhir.NewMemoryBundleStore(2) // tenant 당 최대 2개
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		_ = store.Save(ctx, bundleWith("b-"+rune2str(i), "P"), "T", "P")
	}
	if store.Count(ctx, "T") != 2 {
		t.Errorf("max=2 인데 count = %d", store.Count(ctx, "T"))
	}
	// 가장 최근 두 개 (b-3, b-4) 만 남아야 함
	if _, err := store.Get(ctx, "b-0", "T"); !errors.Is(err, fhir.ErrBundleNotFound) {
		t.Error("FIFO 위반: 오래된 b-0 가 남아있음")
	}
	if _, err := store.Get(ctx, "b-4", "T"); err != nil {
		t.Error("최신 b-4 가 없어짐")
	}
}

func TestMemoryBundleStore_Delete(t *testing.T) {
	store := fhir.NewMemoryBundleStore(0)
	ctx := context.Background()

	_ = store.Save(ctx, bundleWith("b1", "P"), "T", "P")
	if err := store.Delete(ctx, "b1", "T"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Get(ctx, "b1", "T"); !errors.Is(err, fhir.ErrBundleNotFound) {
		t.Error("Delete 후에도 조회됨")
	}
	if err := store.Delete(ctx, "missing", "T"); !errors.Is(err, fhir.ErrBundleNotFound) {
		t.Error("미존재 항목 Delete 통과")
	}
}

func TestMemoryBundleStore_PatientFilter_EmptyError(t *testing.T) {
	store := fhir.NewMemoryBundleStore(0)
	if _, err := store.ListByPatient(context.Background(), "T", "", 0); err == nil {
		t.Error("빈 patientID 통과")
	}
}

func TestMemoryBundleStore_SaveNilBundle(t *testing.T) {
	store := fhir.NewMemoryBundleStore(0)
	if err := store.Save(context.Background(), nil, "T", "P"); err == nil {
		t.Error("nil Bundle 통과")
	}
	b := fhir.NewBundle("collection")
	b.ID = "" // empty
	if err := store.Save(context.Background(), b, "T", "P"); err == nil {
		t.Error("빈 Bundle.ID 통과")
	}
}

// 작은 헬퍼 — rune → ascii 단일 문자열 (테스트용 id 생성).
func rune2str(i int) string {
	return string(rune('0' + i))
}
