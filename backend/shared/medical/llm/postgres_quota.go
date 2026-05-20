package llm

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// QuotaQuerier 는 PostgresQuotaStore 의 SELECT 인터페이스.
//
// QueryContext 는 List 용 (다중 행).
type QuotaQuerier interface {
	QueryRowContext(ctx context.Context, sql string, args ...interface{}) AuditRow
	QueryContext(ctx context.Context, sql string, args ...interface{}) (QuotaRows, error)
}

// QuotaRows 는 다중 행 SELECT 결과.
type QuotaRows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
	Err() error
}

// PostgresQuotaStore 는 QuotaStore 의 PostgreSQL 영속 구현.
type PostgresQuotaStore struct {
	exec    AuditExecer
	querier QuotaQuerier
}

// NewPostgresQuotaStore 생성. exec/querier 둘 다 필요 (List 위해).
func NewPostgresQuotaStore(exec AuditExecer, querier QuotaQuerier) (*PostgresQuotaStore, error) {
	if exec == nil || querier == nil {
		return nil, errors.New("exec/querier 필수")
	}
	return &PostgresQuotaStore{exec: exec, querier: querier}, nil
}

// Get 은 tenant 의 quota config 조회.
func (s *PostgresQuotaStore) Get(ctx context.Context, tenantID string) (*QuotaConfig, error) {
	row := s.querier.QueryRowContext(ctx, `
		SELECT tenant_id, daily_token_limit, monthly_token_limit,
		       daily_request_limit, updated_at
		FROM llm_quota WHERE tenant_id = $1
	`, tenantID)

	var cfg QuotaConfig
	err := row.Scan(&cfg.TenantID, &cfg.DailyTokenLimit, &cfg.MonthlyTokenLimit,
		&cfg.DailyRequestLimit, &cfg.UpdatedAt)
	if err != nil {
		// pgx.ErrNoRows 또는 sql.ErrNoRows → ErrQuotaConfigNotFound
		msg := err.Error()
		if msg == "no rows in result set" || msg == "sql: no rows in result set" {
			return nil, ErrQuotaConfigNotFound
		}
		return nil, fmt.Errorf("Get llm_quota: %w", err)
	}
	return &cfg, nil
}

// Set 은 quota config UPSERT.
func (s *PostgresQuotaStore) Set(ctx context.Context, cfg QuotaConfig) error {
	if cfg.TenantID == "" {
		return errors.New("TenantID 필수")
	}
	if cfg.UpdatedAt.IsZero() {
		cfg.UpdatedAt = time.Now()
	}
	_, err := s.exec.ExecContext(ctx, `
		INSERT INTO llm_quota
			(tenant_id, daily_token_limit, monthly_token_limit, daily_request_limit, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (tenant_id) DO UPDATE
		SET daily_token_limit = EXCLUDED.daily_token_limit,
		    monthly_token_limit = EXCLUDED.monthly_token_limit,
		    daily_request_limit = EXCLUDED.daily_request_limit,
		    updated_at = EXCLUDED.updated_at
	`, cfg.TenantID, cfg.DailyTokenLimit, cfg.MonthlyTokenLimit,
		cfg.DailyRequestLimit, cfg.UpdatedAt)
	if err != nil {
		return fmt.Errorf("UPSERT llm_quota: %w", err)
	}
	return nil
}

// Delete 는 quota 제거.
func (s *PostgresQuotaStore) Delete(ctx context.Context, tenantID string) error {
	_, err := s.exec.ExecContext(ctx, `
		DELETE FROM llm_quota WHERE tenant_id = $1
	`, tenantID)
	return err
}

// List 는 모든 quota config 반환.
func (s *PostgresQuotaStore) List(ctx context.Context) ([]*QuotaConfig, error) {
	rows, err := s.querier.QueryContext(ctx, `
		SELECT tenant_id, daily_token_limit, monthly_token_limit,
		       daily_request_limit, updated_at
		FROM llm_quota ORDER BY tenant_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*QuotaConfig
	for rows.Next() {
		var cfg QuotaConfig
		if err := rows.Scan(&cfg.TenantID, &cfg.DailyTokenLimit,
			&cfg.MonthlyTokenLimit, &cfg.DailyRequestLimit, &cfg.UpdatedAt); err != nil {
			continue
		}
		c := cfg
		out = append(out, &c)
	}
	return out, rows.Err()
}
