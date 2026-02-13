package memory_test

import (
	"context"
	"testing"
	"time"

	"github.com/manpasik/backend/services/admin-service/internal/repository/memory"
)

// ============================================================================
// AuditLogStore 단위 테스트
// ============================================================================

// TestAuditLogStore_Create는 감사 로그 생성을 테스트합니다.
func TestAuditLogStore_Create(t *testing.T) {
	store := memory.NewAuditLogStore()
	ctx := context.Background()

	log := &memory.AuditLog{
		AdminID:   "admin-001",
		Action:    "config_update",
		Resource:  "config:security.jwt_ttl_hours",
		OldValue:  "24",
		NewValue:  "48",
		IPAddress: "192.168.1.100",
		UserAgent: "Mozilla/5.0",
		CreatedAt: time.Now().UTC(),
	}

	err := store.Create(ctx, log)
	if err != nil {
		t.Fatalf("Create 실패: %v", err)
	}

	// 생성 확인
	count, err := store.Count(ctx)
	if err != nil {
		t.Fatalf("Count 실패: %v", err)
	}
	if count != 1 {
		t.Errorf("Count: got %d, want 1", count)
	}

	// ID 자동 생성 확인
	all, err := store.ListAll(ctx, 10, 0)
	if err != nil {
		t.Fatalf("ListAll 실패: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("ListAll 결과 수: got %d, want 1", len(all))
	}
	if all[0].ID == "" {
		t.Error("ID가 자동 생성되어야 합니다")
	}
	if all[0].AdminID != "admin-001" {
		t.Errorf("AdminID: got %s, want admin-001", all[0].AdminID)
	}
	if all[0].Action != "config_update" {
		t.Errorf("Action: got %s, want config_update", all[0].Action)
	}
	if all[0].OldValue != "24" {
		t.Errorf("OldValue: got %s, want 24", all[0].OldValue)
	}
	if all[0].NewValue != "48" {
		t.Errorf("NewValue: got %s, want 48", all[0].NewValue)
	}
}

// TestAuditLogStore_CreateWithID는 ID가 지정된 경우 해당 ID를 사용하는지 테스트합니다.
func TestAuditLogStore_CreateWithID(t *testing.T) {
	store := memory.NewAuditLogStore()
	ctx := context.Background()

	log := &memory.AuditLog{
		ID:        "custom-id-001",
		AdminID:   "admin-002",
		Action:    "user_ban",
		Resource:  "user:user-123",
		OldValue:  "active",
		NewValue:  "banned",
		CreatedAt: time.Now().UTC(),
	}

	err := store.Create(ctx, log)
	if err != nil {
		t.Fatalf("Create 실패: %v", err)
	}

	all, _ := store.ListAll(ctx, 10, 0)
	if len(all) != 1 {
		t.Fatalf("결과 수: got %d, want 1", len(all))
	}
	if all[0].ID != "custom-id-001" {
		t.Errorf("ID: got %s, want custom-id-001", all[0].ID)
	}
}

// TestAuditLogStore_CreateNil은 nil 로그 생성 시 에러가 없는지 테스트합니다.
func TestAuditLogStore_CreateNil(t *testing.T) {
	store := memory.NewAuditLogStore()
	ctx := context.Background()

	err := store.Create(ctx, nil)
	if err != nil {
		t.Fatalf("nil 로그 Create에서 에러 발생: %v", err)
	}

	count, _ := store.Count(ctx)
	if count != 0 {
		t.Errorf("nil 로그는 저장되지 않아야 합니다: count=%d", count)
	}
}

// TestAuditLogStore_ListByAdmin은 관리자별 조회를 테스트합니다.
func TestAuditLogStore_ListByAdmin(t *testing.T) {
	store := memory.NewAuditLogStore()
	ctx := context.Background()

	// 여러 관리자의 로그 생성
	logs := []*memory.AuditLog{
		{AdminID: "admin-A", Action: "config_update", Resource: "config:key1", CreatedAt: time.Now().UTC().Add(-3 * time.Minute)},
		{AdminID: "admin-B", Action: "user_ban", Resource: "user:user-1", CreatedAt: time.Now().UTC().Add(-2 * time.Minute)},
		{AdminID: "admin-A", Action: "config_create", Resource: "config:key2", CreatedAt: time.Now().UTC().Add(-1 * time.Minute)},
		{AdminID: "admin-A", Action: "system_restart", Resource: "system:main", CreatedAt: time.Now().UTC()},
		{AdminID: "admin-C", Action: "config_delete", Resource: "config:key3", CreatedAt: time.Now().UTC()},
	}
	for _, l := range logs {
		_ = store.Create(ctx, l)
	}

	// admin-A의 로그만 조회 (3개)
	result, err := store.ListByAdmin(ctx, "admin-A", 10, 0)
	if err != nil {
		t.Fatalf("ListByAdmin 실패: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("admin-A 로그 수: got %d, want 3", len(result))
	}

	// 최신순 정렬 확인
	if len(result) >= 2 {
		if result[0].Action != "system_restart" {
			t.Errorf("최신 로그 Action: got %s, want system_restart", result[0].Action)
		}
	}

	// admin-B의 로그 (1개)
	result, _ = store.ListByAdmin(ctx, "admin-B", 10, 0)
	if len(result) != 1 {
		t.Errorf("admin-B 로그 수: got %d, want 1", len(result))
	}

	// 존재하지 않는 관리자
	result, _ = store.ListByAdmin(ctx, "admin-Z", 10, 0)
	if len(result) != 0 {
		t.Errorf("admin-Z 로그 수: got %d, want 0", len(result))
	}
}

// TestAuditLogStore_ListByAction은 액션별 조회를 테스트합니다.
func TestAuditLogStore_ListByAction(t *testing.T) {
	store := memory.NewAuditLogStore()
	ctx := context.Background()

	logs := []*memory.AuditLog{
		{AdminID: "admin-1", Action: "config_update", Resource: "config:key1"},
		{AdminID: "admin-2", Action: "config_update", Resource: "config:key2"},
		{AdminID: "admin-1", Action: "user_ban", Resource: "user:user-1"},
		{AdminID: "admin-3", Action: "config_update", Resource: "config:key3"},
		{AdminID: "admin-1", Action: "config_delete", Resource: "config:key4"},
	}
	for _, l := range logs {
		_ = store.Create(ctx, l)
	}

	// config_update 조회 (3개)
	result, err := store.ListByAction(ctx, "config_update", 10, 0)
	if err != nil {
		t.Fatalf("ListByAction 실패: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("config_update 로그 수: got %d, want 3", len(result))
	}

	// user_ban 조회 (1개)
	result, _ = store.ListByAction(ctx, "user_ban", 10, 0)
	if len(result) != 1 {
		t.Errorf("user_ban 로그 수: got %d, want 1", len(result))
	}

	// 존재하지 않는 액션
	result, _ = store.ListByAction(ctx, "nonexistent_action", 10, 0)
	if len(result) != 0 {
		t.Errorf("nonexistent_action 로그 수: got %d, want 0", len(result))
	}
}

// TestAuditLogStore_ListAll_Pagination은 전체 조회 및 페이지네이션을 테스트합니다.
func TestAuditLogStore_ListAll_Pagination(t *testing.T) {
	store := memory.NewAuditLogStore()
	ctx := context.Background()

	// 5개 로그 생성
	for i := 0; i < 5; i++ {
		_ = store.Create(ctx, &memory.AuditLog{
			AdminID:   "admin-pager",
			Action:    "config_update",
			Resource:  "config:key",
			CreatedAt: time.Now().UTC(),
		})
	}

	// 전체 조회
	all, err := store.ListAll(ctx, 10, 0)
	if err != nil {
		t.Fatalf("ListAll 실패: %v", err)
	}
	if len(all) != 5 {
		t.Errorf("전체 로그 수: got %d, want 5", len(all))
	}

	// 페이지네이션: limit=2, offset=0 → 2개
	page1, _ := store.ListAll(ctx, 2, 0)
	if len(page1) != 2 {
		t.Errorf("page1 수: got %d, want 2", len(page1))
	}

	// 페이지네이션: limit=2, offset=2 → 2개
	page2, _ := store.ListAll(ctx, 2, 2)
	if len(page2) != 2 {
		t.Errorf("page2 수: got %d, want 2", len(page2))
	}

	// 페이지네이션: limit=2, offset=4 → 1개
	page3, _ := store.ListAll(ctx, 2, 4)
	if len(page3) != 1 {
		t.Errorf("page3 수: got %d, want 1", len(page3))
	}

	// 페이지네이션: offset 초과 → 0개
	page4, _ := store.ListAll(ctx, 2, 10)
	if len(page4) != 0 {
		t.Errorf("offset 초과 시 결과: got %d, want 0", len(page4))
	}
}

// TestAuditLogStore_Count는 전체 카운트를 테스트합니다.
func TestAuditLogStore_Count(t *testing.T) {
	store := memory.NewAuditLogStore()
	ctx := context.Background()

	// 초기 카운트 = 0
	count, err := store.Count(ctx)
	if err != nil {
		t.Fatalf("Count 실패: %v", err)
	}
	if count != 0 {
		t.Errorf("초기 Count: got %d, want 0", count)
	}

	// 3개 추가
	for i := 0; i < 3; i++ {
		_ = store.Create(ctx, &memory.AuditLog{
			AdminID: "admin-counter",
			Action:  "config_update",
		})
	}

	count, _ = store.Count(ctx)
	if count != 3 {
		t.Errorf("3개 추가 후 Count: got %d, want 3", count)
	}
}

// TestAuditLogStore_DeepCopy는 저장된 로그가 원본과 분리되는지 테스트합니다.
func TestAuditLogStore_DeepCopy(t *testing.T) {
	store := memory.NewAuditLogStore()
	ctx := context.Background()

	log := &memory.AuditLog{
		AdminID:  "admin-copy",
		Action:   "config_update",
		Resource: "config:original",
		OldValue: "v1",
		NewValue: "v2",
	}

	_ = store.Create(ctx, log)

	// 원본 수정
	log.Resource = "config:modified"
	log.NewValue = "v3"

	// 저장된 로그는 원본 변경의 영향을 받지 않아야 함
	all, _ := store.ListAll(ctx, 10, 0)
	if len(all) != 1 {
		t.Fatalf("결과 수: got %d, want 1", len(all))
	}
	if all[0].Resource != "config:original" {
		t.Errorf("Resource가 변경됨: got %s, want config:original", all[0].Resource)
	}
	if all[0].NewValue != "v2" {
		t.Errorf("NewValue가 변경됨: got %s, want v2", all[0].NewValue)
	}
}

// TestAuditLogStore_AutoCreatedAt은 CreatedAt가 없을 때 자동 설정되는지 테스트합니다.
func TestAuditLogStore_AutoCreatedAt(t *testing.T) {
	store := memory.NewAuditLogStore()
	ctx := context.Background()

	log := &memory.AuditLog{
		AdminID: "admin-time",
		Action:  "config_update",
		// CreatedAt 미설정
	}

	before := time.Now().UTC()
	_ = store.Create(ctx, log)
	after := time.Now().UTC()

	all, _ := store.ListAll(ctx, 10, 0)
	if len(all) != 1 {
		t.Fatalf("결과 수: got %d, want 1", len(all))
	}

	createdAt := all[0].CreatedAt
	if createdAt.Before(before) || createdAt.After(after) {
		t.Errorf("CreatedAt가 자동 설정 범위를 벗어남: %v (expected between %v and %v)", createdAt, before, after)
	}
}
