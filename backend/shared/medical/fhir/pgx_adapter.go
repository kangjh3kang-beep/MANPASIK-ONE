package fhir

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PgxBundleAdapter 는 pgxpool.Pool 을 fhir.SQLDB 인터페이스로 어댑팅 (Phase AW-3).
//
// pgx 의 메서드 시그니처가 database/sql 과 다르므로 어댑터로 통일.
//
// 사용 예 (gateway main.go):
//
//	pool, _ := pgxpool.New(ctx, dsn)
//	store, _ := fhir.NewPostgresBundleStore(fhir.NewPgxBundleAdapter(pool))
//	restHandler.SetHL7BundleStore(store)
type PgxBundleAdapter struct {
	pool *pgxpool.Pool
}

// NewPgxBundleAdapter 생성. pool=nil 이면 nil 반환.
func NewPgxBundleAdapter(pool *pgxpool.Pool) *PgxBundleAdapter {
	if pool == nil {
		return nil
	}
	return &PgxBundleAdapter{pool: pool}
}

// ExecContext 는 pool.Exec 위임.
func (a *PgxBundleAdapter) ExecContext(ctx context.Context, sql string, args ...interface{}) (SQLResult, error) {
	if a == nil || a.pool == nil {
		return nil, errors.New("pgx adapter nil")
	}
	tag, err := a.pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return pgxBundleResult{tag: tag}, nil
}

// QueryContext 는 pool.Query 위임.
func (a *PgxBundleAdapter) QueryContext(ctx context.Context, sql string, args ...interface{}) (SQLRows, error) {
	if a == nil || a.pool == nil {
		return nil, errors.New("pgx adapter nil")
	}
	rows, err := a.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return pgxBundleRows{rows: rows}, nil
}

// QueryRowContext 는 pool.QueryRow 위임.
func (a *PgxBundleAdapter) QueryRowContext(ctx context.Context, sql string, args ...interface{}) SQLRow {
	if a == nil || a.pool == nil {
		return pgxBundleRow{err: errors.New("pgx adapter nil")}
	}
	return pgxBundleRow{row: a.pool.QueryRow(ctx, sql, args...)}
}

// pgxBundleResult — pgx.CommandTag → SQLResult.
type pgxBundleResult struct {
	tag pgconn.CommandTag
}

func (r pgxBundleResult) RowsAffected() (int64, error) {
	return r.tag.RowsAffected(), nil
}

// pgxBundleRows — pgx.Rows → SQLRows.
type pgxBundleRows struct {
	rows pgx.Rows
}

func (r pgxBundleRows) Next() bool                   { return r.rows.Next() }
func (r pgxBundleRows) Scan(dest ...interface{}) error { return r.rows.Scan(dest...) }
func (r pgxBundleRows) Close() error                  { r.rows.Close(); return nil }
func (r pgxBundleRows) Err() error                    { return r.rows.Err() }

// pgxBundleRow — pgx.Row → SQLRow. 사전 에러도 그대로 보존.
type pgxBundleRow struct {
	row pgx.Row
	err error
}

func (r pgxBundleRow) Scan(dest ...interface{}) error {
	if r.err != nil {
		return r.err
	}
	err := r.row.Scan(dest...)
	if errors.Is(err, pgx.ErrNoRows) {
		// 공통 "no rows" 문자열로 정규화 → isNoRowsErr 가 감지하도록
		return errors.New("no rows in result set")
	}
	return err
}
