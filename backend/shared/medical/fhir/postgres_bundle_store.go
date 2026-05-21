package fhir

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// ============================================================================
// PostgresBundleStore — Phase AW
// ============================================================================
//
// 인메모리 MemoryBundleStore 의 영속화 백엔드. 44-bundle-store.sql 스키마 사용.
// SQLDB 인터페이스를 통해 pgxpool.Pool 또는 *sql.DB 모두 호환.

// SQLDB 는 PostgresBundleStore 가 사용하는 최소 DB 인터페이스.
//
// 외부 의존을 피하기 위해 별도 정의 — 실제로는 pgxpool.Pool 또는 *sql.DB 가
// 이 인터페이스를 만족 (PgxBundleAdapter 참고).
type SQLDB interface {
	ExecContext(ctx context.Context, sql string, args ...interface{}) (SQLResult, error)
	QueryContext(ctx context.Context, sql string, args ...interface{}) (SQLRows, error)
	QueryRowContext(ctx context.Context, sql string, args ...interface{}) SQLRow
}

type SQLResult interface {
	RowsAffected() (int64, error)
}

type SQLRows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
	Err() error
}

type SQLRow interface {
	Scan(dest ...interface{}) error
}

// PostgresBundleStore 는 BundleStore 의 PostgreSQL 백엔드.
type PostgresBundleStore struct {
	db SQLDB
}

// NewPostgresBundleStore 생성. db=nil 이면 에러.
func NewPostgresBundleStore(db SQLDB) (*PostgresBundleStore, error) {
	if db == nil {
		return nil, errors.New("db 필수")
	}
	return &PostgresBundleStore{db: db}, nil
}

// Save 는 Bundle 저장 (UPSERT — 같은 (tenant_id, bundle_id) 이면 update).
//
// 멱등성: 같은 ID 두 번 저장해도 한 entry. stored_at 은 최초 저장 시점 보존
// (MemoryBundleStore 와 동일 의미).
func (s *PostgresBundleStore) Save(ctx context.Context, b *Bundle, tenantID, patientID string) error {
	if b == nil || b.ID == "" {
		return errors.New("Bundle.ID 가 필요")
	}
	payload, err := json.Marshal(b)
	if err != nil {
		return fmt.Errorf("Bundle 직렬화: %w", err)
	}
	// stored_at 은 ON CONFLICT 시 보존 — payload/patient_id 만 갱신
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO bundle_store (tenant_id, bundle_id, patient_id, payload, stored_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (tenant_id, bundle_id) DO UPDATE
		SET payload    = EXCLUDED.payload,
		    patient_id = EXCLUDED.patient_id
		    -- stored_at 은 의도적으로 보존 (멱등 update)
	`, tenantID, b.ID, patientID, payload)
	if err != nil {
		return fmt.Errorf("Save: %w", err)
	}
	return nil
}

// Get 은 단일 Bundle 조회 (tenant 일치 필수).
func (s *PostgresBundleStore) Get(ctx context.Context, bundleID, tenantID string) (*StoredBundle, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT bundle_id, tenant_id, patient_id, payload, stored_at
		FROM bundle_store
		WHERE tenant_id = $1 AND bundle_id = $2
	`, tenantID, bundleID)
	return scanStoredBundle(row)
}

// ListByPatient 는 tenant 내 환자별 Bundle 목록 (최신순).
func (s *PostgresBundleStore) ListByPatient(ctx context.Context, tenantID, patientID string, limit int) ([]*StoredBundle, error) {
	if patientID == "" {
		return nil, errors.New("patientID 필수")
	}
	q := `
		SELECT bundle_id, tenant_id, patient_id, payload, stored_at
		FROM bundle_store
		WHERE tenant_id = $1 AND patient_id = $2
		ORDER BY stored_at DESC
	`
	args := []interface{}{tenantID, patientID}
	if limit > 0 {
		q += " LIMIT $3"
		args = append(args, limit)
	}
	return s.queryBundles(ctx, q, args...)
}

// ListByTenant 는 tenant 의 전체 Bundle (최신순).
func (s *PostgresBundleStore) ListByTenant(ctx context.Context, tenantID string, limit int) ([]*StoredBundle, error) {
	q := `
		SELECT bundle_id, tenant_id, patient_id, payload, stored_at
		FROM bundle_store
		WHERE tenant_id = $1
		ORDER BY stored_at DESC
	`
	args := []interface{}{tenantID}
	if limit > 0 {
		q += " LIMIT $2"
		args = append(args, limit)
	}
	return s.queryBundles(ctx, q, args...)
}

// Delete 는 단일 Bundle 삭제. 미존재 시 ErrBundleNotFound.
func (s *PostgresBundleStore) Delete(ctx context.Context, bundleID, tenantID string) error {
	res, err := s.db.ExecContext(ctx,
		`DELETE FROM bundle_store WHERE tenant_id = $1 AND bundle_id = $2`,
		tenantID, bundleID)
	if err != nil {
		return fmt.Errorf("Delete: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrBundleNotFound
	}
	return nil
}

// Count 는 tenant 내 항목 수 (운영 메트릭).
func (s *PostgresBundleStore) Count(ctx context.Context, tenantID string) int {
	row := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM bundle_store WHERE tenant_id = $1`, tenantID)
	var n int
	if err := row.Scan(&n); err != nil {
		return 0
	}
	return n
}

// queryBundles 는 공통 다중 행 쿼리 헬퍼.
func (s *PostgresBundleStore) queryBundles(ctx context.Context, q string, args ...interface{}) ([]*StoredBundle, error) {
	rows, err := s.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var out []*StoredBundle
	for rows.Next() {
		entry, err := scanStoredBundle(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}
	return out, nil
}

// scanStoredBundle 은 단일 행 → StoredBundle.
//
// rowOrRows 는 SQLRow 또는 SQLRows 모두 가능 (둘 다 Scan 메서드를 가짐).
func scanStoredBundle(rowOrRows interface{ Scan(...interface{}) error }) (*StoredBundle, error) {
	var (
		bundleID  string
		tenantID  string
		patientID string
		payload   []byte
		storedAt  time.Time
	)
	if err := rowOrRows.Scan(&bundleID, &tenantID, &patientID, &payload, &storedAt); err != nil {
		// pgx.ErrNoRows / sql.ErrNoRows 도 ErrBundleNotFound 로 정규화
		if isNoRowsErr(err) {
			return nil, ErrBundleNotFound
		}
		return nil, fmt.Errorf("scan: %w", err)
	}
	bundle := &Bundle{}
	if err := json.Unmarshal(payload, bundle); err != nil {
		return nil, fmt.Errorf("Bundle 역직렬화: %w", err)
	}
	// JSON 에서 ID 가 비어있는 경우 column 값 사용 (안전망)
	if bundle.ID == "" {
		bundle.ID = bundleID
	}
	return &StoredBundle{
		Bundle:    bundle,
		TenantID:  tenantID,
		PatientID: patientID,
		StoredAt:  storedAt,
	}, nil
}

// isNoRowsErr 는 pgx / sql 의 "no rows" 에러 감지 (string 매칭으로 의존 회피).
func isNoRowsErr(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return s == "no rows in result set" ||
		s == "sql: no rows in result set"
}
