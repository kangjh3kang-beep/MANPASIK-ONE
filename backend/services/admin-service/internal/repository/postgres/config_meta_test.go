// Package postgres는 ConfigMetadata/ConfigTranslation PostgreSQL 리포지토리 테스트입니다.
// DB 연결이 없으면 자동 스킵됩니다. 통합 테스트 실행 시 DB_HOST 환경변수를 설정하세요.
package postgres

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/admin-service/internal/service"
)

// testPool은 테스트용 DB 연결 풀입니다.
func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	host := os.Getenv("DB_HOST")
	if host == "" {
		t.Skip("DB_HOST 미설정: PostgreSQL 통합 테스트 건너뜀")
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "manpasik"
	}
	pass := os.Getenv("DB_PASSWORD")
	if pass == "" {
		pass = "manpasik_dev"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "manpasik_dev"
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, dbname)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Skipf("DB 연결 실패, 건너뜀: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Skipf("DB ping 실패, 건너뜀: %v", err)
	}

	t.Cleanup(func() { pool.Close() })
	return pool
}

// ============================================================================
// ConfigMetadataRepository 테스트
// ============================================================================

func TestConfigMetadataRepo_GetByKey_존재하는키(t *testing.T) {
	pool := testPool(t)
	repo := NewConfigMetadataRepository(pool)
	ctx := context.Background()

	// DB에 시드 데이터(25-admin-settings-ext.sql)의 maintenance_mode가 있어야 함
	meta, err := repo.GetByKey(ctx, "maintenance_mode")
	if err != nil {
		t.Fatalf("GetByKey 에러: %v", err)
	}
	if meta == nil {
		t.Fatal("maintenance_mode 메타데이터가 nil — 시드 데이터 확인 필요")
	}
	if meta.ConfigKey != "maintenance_mode" {
		t.Errorf("ConfigKey = %q, want %q", meta.ConfigKey, "maintenance_mode")
	}
	if meta.Category != "general" {
		t.Errorf("Category = %q, want %q", meta.Category, "general")
	}
	if meta.ValueType != "boolean" {
		t.Errorf("ValueType = %q, want %q", meta.ValueType, "boolean")
	}
}

func TestConfigMetadataRepo_GetByKey_존재하지않는키(t *testing.T) {
	pool := testPool(t)
	repo := NewConfigMetadataRepository(pool)
	ctx := context.Background()

	meta, err := repo.GetByKey(ctx, "nonexistent_key_xxx")
	if err != nil {
		t.Fatalf("GetByKey 에러: %v", err)
	}
	if meta != nil {
		t.Errorf("존재하지 않는 키에 대해 nil 반환 기대, got: %+v", meta)
	}
}

func TestConfigMetadataRepo_ListByCategory(t *testing.T) {
	pool := testPool(t)
	repo := NewConfigMetadataRepository(pool)
	ctx := context.Background()

	list, err := repo.ListByCategory(ctx, "general")
	if err != nil {
		t.Fatalf("ListByCategory 에러: %v", err)
	}
	if len(list) == 0 {
		t.Error("general 카테고리 메타데이터가 비어 있음 — 시드 데이터 확인 필요")
	}
	for _, m := range list {
		if m.Category != "general" {
			t.Errorf("카테고리 필터링 실패: got %q, want %q", m.Category, "general")
		}
		if !m.IsActive {
			t.Error("비활성 항목이 포함됨")
		}
	}
}

func TestConfigMetadataRepo_ListAll(t *testing.T) {
	pool := testPool(t)
	repo := NewConfigMetadataRepository(pool)
	ctx := context.Background()

	list, err := repo.ListAll(ctx)
	if err != nil {
		t.Fatalf("ListAll 에러: %v", err)
	}
	if len(list) == 0 {
		t.Error("메타데이터가 비어 있음 — 시드 데이터 확인 필요")
	}
	for _, m := range list {
		if !m.IsActive {
			t.Error("비활성 항목이 포함됨")
		}
	}
}

func TestConfigMetadataRepo_CountByCategory(t *testing.T) {
	pool := testPool(t)
	repo := NewConfigMetadataRepository(pool)
	ctx := context.Background()

	counts, err := repo.CountByCategory(ctx)
	if err != nil {
		t.Fatalf("CountByCategory 에러: %v", err)
	}
	if len(counts) == 0 {
		t.Error("카테고리 카운트가 비어 있음")
	}
	if counts["general"] == 0 {
		t.Error("general 카테고리 카운트가 0")
	}
}

// ============================================================================
// ConfigTranslationRepository 테스트
// ============================================================================

func TestConfigTranslationRepo_GetByKeyAndLang_존재(t *testing.T) {
	pool := testPool(t)
	repo := NewConfigTranslationRepository(pool)
	ctx := context.Background()

	trans, err := repo.GetByKeyAndLang(ctx, "maintenance_mode", "ko")
	if err != nil {
		t.Fatalf("GetByKeyAndLang 에러: %v", err)
	}
	if trans == nil {
		t.Fatal("maintenance_mode/ko 번역이 nil — 시드 데이터 확인 필요")
	}
	if trans.ConfigKey != "maintenance_mode" {
		t.Errorf("ConfigKey = %q, want %q", trans.ConfigKey, "maintenance_mode")
	}
	if trans.LanguageCode != "ko" {
		t.Errorf("LanguageCode = %q, want %q", trans.LanguageCode, "ko")
	}
	if trans.DisplayName == "" {
		t.Error("DisplayName이 비어 있음")
	}
}

func TestConfigTranslationRepo_GetByKeyAndLang_미존재(t *testing.T) {
	pool := testPool(t)
	repo := NewConfigTranslationRepository(pool)
	ctx := context.Background()

	trans, err := repo.GetByKeyAndLang(ctx, "nonexistent_key", "ko")
	if err != nil {
		t.Fatalf("GetByKeyAndLang 에러: %v", err)
	}
	if trans != nil {
		t.Errorf("존재하지 않는 키에 대해 nil 기대, got: %+v", trans)
	}
}

func TestConfigTranslationRepo_ListByKey(t *testing.T) {
	pool := testPool(t)
	repo := NewConfigTranslationRepository(pool)
	ctx := context.Background()

	list, err := repo.ListByKey(ctx, "maintenance_mode")
	if err != nil {
		t.Fatalf("ListByKey 에러: %v", err)
	}
	if len(list) == 0 {
		t.Error("maintenance_mode 번역이 비어 있음 — 시드 데이터 확인 필요")
	}
	for _, tr := range list {
		if tr.ConfigKey != "maintenance_mode" {
			t.Errorf("키 필터링 실패: got %q, want %q", tr.ConfigKey, "maintenance_mode")
		}
	}
}

func TestConfigTranslationRepo_ListByLang(t *testing.T) {
	pool := testPool(t)
	repo := NewConfigTranslationRepository(pool)
	ctx := context.Background()

	m, err := repo.ListByLang(ctx, "ko")
	if err != nil {
		t.Fatalf("ListByLang 에러: %v", err)
	}
	if len(m) == 0 {
		t.Error("ko 언어 번역이 비어 있음 — 시드 데이터 확인 필요")
	}

	// 반환된 모든 항목이 ko 언어여야 함
	for key, tr := range m {
		if tr.ConfigKey != key {
			t.Errorf("맵 키(%q)와 ConfigKey(%q) 불일치", key, tr.ConfigKey)
		}
		if tr.LanguageCode != "ko" {
			t.Errorf("언어 필터링 실패: got %q, want %q", tr.LanguageCode, "ko")
		}
	}
}

// ============================================================================
// 인터페이스 준수 컴파일 타임 검증
// ============================================================================

var _ service.ConfigMetadataRepository = (*ConfigMetadataRepository)(nil)
var _ service.ConfigTranslationRepository = (*ConfigTranslationRepository)(nil)
