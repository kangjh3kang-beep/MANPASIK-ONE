package llm

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PgxAdapter 는 pgxpool.Pool 을 AuditExecer + AuditQuerier + QuotaQuerier 로
// 어댑팅 (Phase AK-2).
//
// pgx 의 메서드 시그니처가 database/sql 과 다르므로 (Exec(ctx,...) vs
// ExecContext(ctx,...)) 어댑터로 통일.
//
// 사용 예 (admin-service main.go):
//
//	pool, _ := pgxpool.New(ctx, dsn)
//	pgxAdapter := llm.NewPgxAdapter(pool)
//	auditLog, _ := llm.NewPostgresAuditLog(pgxAdapter, pgxAdapter)
//	quotaStore, _ := llm.NewPostgresQuotaStore(pgxAdapter, pgxAdapter)
type PgxAdapter struct {
	pool *pgxpool.Pool
}

// NewPgxAdapter 생성. pool=nil 이면 nil 반환 (이후 호출이 에러 반환).
func NewPgxAdapter(pool *pgxpool.Pool) *PgxAdapter {
	if pool == nil {
		return nil
	}
	return &PgxAdapter{pool: pool}
}

// ExecContext 는 pool.Exec 위임.
func (a *PgxAdapter) ExecContext(ctx context.Context, sql string, args ...interface{}) (AuditExecResult, error) {
	if a == nil || a.pool == nil {
		return nil, errors.New("pgx adapter nil")
	}
	tag, err := a.pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return pgxExecResult{tag: tag}, nil
}

type pgxExecResult struct {
	tag interface{ RowsAffected() int64 }
}

func (r pgxExecResult) RowsAffected() (int64, error) { return r.tag.RowsAffected(), nil }

// QueryRowContext 는 pool.QueryRow 위임.
func (a *PgxAdapter) QueryRowContext(ctx context.Context, sql string, args ...interface{}) AuditRow {
	if a == nil || a.pool == nil {
		return errAuditRow{err: errors.New("pgx adapter nil")}
	}
	return pgxAuditRow{row: a.pool.QueryRow(ctx, sql, args...)}
}

// QueryContext 는 pool.Query 위임 (QuotaQuerier 인터페이스).
func (a *PgxAdapter) QueryContext(ctx context.Context, sql string, args ...interface{}) (QuotaRows, error) {
	if a == nil || a.pool == nil {
		return nil, errors.New("pgx adapter nil")
	}
	rows, err := a.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &pgxQuotaRows{rows: rows}, nil
}

type pgxAuditRow struct {
	row pgx.Row
}

func (r pgxAuditRow) Scan(dest ...interface{}) error {
	err := r.row.Scan(dest...)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return errors.New("no rows in result set")
	}
	return err
}

type errAuditRow struct{ err error }

func (r errAuditRow) Scan(_ ...interface{}) error { return r.err }

type pgxQuotaRows struct {
	rows pgx.Rows
}

func (r *pgxQuotaRows) Next() bool                   { return r.rows.Next() }
func (r *pgxQuotaRows) Scan(dest ...interface{}) error { return r.rows.Scan(dest...) }
func (r *pgxQuotaRows) Close() error                  { r.rows.Close(); return nil }
func (r *pgxQuotaRows) Err() error                    { return r.rows.Err() }
