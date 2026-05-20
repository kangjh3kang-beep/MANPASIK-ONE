package tenancy

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PgxAdapter 는 pgxpool.Pool 을 SQLDB 인터페이스로 어댑팅.
//
// pgx 의 메서드 시그니처가 database/sql 과 다르므로 (예: ExecContext 대신
// Exec(ctx,...)) 이 어댑터로 통일.
//
// 사용 예 (서비스 main.go):
//
//	pool, _ := pgxpool.New(ctx, dsn)
//	memStore, _ := tenancy.NewPostgresMembershipStore(tenancy.NewPgxAdapter(pool))
//	invStore, _ := tenancy.NewPostgresInvitationStore(tenancy.NewPgxAdapter(pool))
type PgxAdapter struct {
	pool *pgxpool.Pool
}

// NewPgxAdapter 생성. pool=nil 이면 nil 반환 (Store 가 db=nil 검사로 거부).
func NewPgxAdapter(pool *pgxpool.Pool) *PgxAdapter {
	if pool == nil {
		return nil
	}
	return &PgxAdapter{pool: pool}
}

// ExecContext 는 pool.Exec 위임.
func (a *PgxAdapter) ExecContext(ctx context.Context, sql string, args ...interface{}) (SQLResult, error) {
	if a == nil || a.pool == nil {
		return nil, errors.New("pgx adapter nil")
	}
	tag, err := a.pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return pgxResult{tag: tag}, nil
}

// QueryContext 는 pool.Query 위임.
func (a *PgxAdapter) QueryContext(ctx context.Context, sql string, args ...interface{}) (SQLRows, error) {
	if a == nil || a.pool == nil {
		return nil, errors.New("pgx adapter nil")
	}
	rows, err := a.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &pgxRows{rows: rows}, nil
}

// QueryRowContext 는 pool.QueryRow 위임.
func (a *PgxAdapter) QueryRowContext(ctx context.Context, sql string, args ...interface{}) SQLRow {
	if a == nil || a.pool == nil {
		return errRow{err: errors.New("pgx adapter nil")}
	}
	return pgxRow{row: a.pool.QueryRow(ctx, sql, args...)}
}

// pgxResult 는 pgconn.CommandTag → SQLResult 어댑터.
type pgxResult struct {
	tag pgconn.CommandTag
}

func (r pgxResult) RowsAffected() (int64, error) {
	return r.tag.RowsAffected(), nil
}

// pgxRows 는 pgx.Rows → SQLRows 어댑터.
type pgxRows struct {
	rows pgx.Rows
}

func (r *pgxRows) Next() bool                   { return r.rows.Next() }
func (r *pgxRows) Scan(dest ...interface{}) error { return r.rows.Scan(dest...) }
func (r *pgxRows) Close() error                  { r.rows.Close(); return nil }
func (r *pgxRows) Err() error                    { return r.rows.Err() }

// pgxRow 는 pgx.Row → SQLRow 어댑터.
//
// pgx.Row.Scan 은 행이 없으면 pgx.ErrNoRows 반환. 이를 일관된 메시지
// "no rows in result set" 으로 변환하여 scanInvitation/scanMembership 의 패턴
// 매칭이 동작하도록 함.
type pgxRow struct {
	row pgx.Row
}

func (r pgxRow) Scan(dest ...interface{}) error {
	err := r.row.Scan(dest...)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return errors.New("no rows in result set")
	}
	return err
}

// errRow 는 어댑터 자체 에러를 Scan 으로 전달.
type errRow struct{ err error }

func (r errRow) Scan(_ ...interface{}) error { return r.err }
