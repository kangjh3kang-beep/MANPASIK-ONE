// Package postgres — ConfigMetadata/ConfigTranslation PostgreSQL 리포지토리 구현.
// DB 스키마: infrastructure/database/init/25-admin-settings-ext.sql
package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/admin-service/internal/service"
)

// ============================================================================
// ConfigMetadataRepository — PostgreSQL 구현
// ============================================================================

// ConfigMetadataRepository는 PostgreSQL 기반 설정 메타데이터 저장소입니다.
type ConfigMetadataRepository struct {
	pool *pgxpool.Pool
}

// NewConfigMetadataRepository는 PostgreSQL ConfigMetadataRepository를 생성합니다.
func NewConfigMetadataRepository(pool *pgxpool.Pool) *ConfigMetadataRepository {
	return &ConfigMetadataRepository{pool: pool}
}

// GetByKey는 키로 메타데이터를 조회합니다.
func (r *ConfigMetadataRepository) GetByKey(ctx context.Context, key string) (*service.ConfigMetadata, error) {
	const q = `SELECT
		config_key, category::text, value_type::text, security_level::text,
		is_required, COALESCE(default_value,''), allowed_values,
		COALESCE(validation_regex,''), validation_min, validation_max,
		COALESCE(depends_on,''), COALESCE(depends_value,''),
		COALESCE(env_var_name,''), COALESCE(service_name,''),
		restart_required, display_order, is_active,
		created_at, updated_at
	FROM config_metadata WHERE config_key = $1`

	var m service.ConfigMetadata
	var allowedValues []string
	err := r.pool.QueryRow(ctx, q, key).Scan(
		&m.ConfigKey, &m.Category, &m.ValueType, &m.SecurityLevel,
		&m.IsRequired, &m.DefaultValue, &allowedValues,
		&m.ValidationRegex, &m.ValidationMin, &m.ValidationMax,
		&m.DependsOn, &m.DependsValue,
		&m.EnvVarName, &m.ServiceName,
		&m.RestartRequired, &m.DisplayOrder, &m.IsActive,
		&m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	m.AllowedValues = allowedValues
	return &m, nil
}

// ListByCategory는 카테고리별 활성 메타데이터를 조회합니다.
func (r *ConfigMetadataRepository) ListByCategory(ctx context.Context, category string) ([]*service.ConfigMetadata, error) {
	const q = `SELECT
		config_key, category::text, value_type::text, security_level::text,
		is_required, COALESCE(default_value,''), allowed_values,
		COALESCE(validation_regex,''), validation_min, validation_max,
		COALESCE(depends_on,''), COALESCE(depends_value,''),
		COALESCE(env_var_name,''), COALESCE(service_name,''),
		restart_required, display_order, is_active,
		created_at, updated_at
	FROM config_metadata
	WHERE is_active = true AND ($1 = '' OR category::text = $1)
	ORDER BY display_order ASC`

	rows, err := r.pool.Query(ctx, q, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanMetadataRows(rows)
}

// ListAll은 모든 활성 메타데이터를 조회합니다.
func (r *ConfigMetadataRepository) ListAll(ctx context.Context) ([]*service.ConfigMetadata, error) {
	const q = `SELECT
		config_key, category::text, value_type::text, security_level::text,
		is_required, COALESCE(default_value,''), allowed_values,
		COALESCE(validation_regex,''), validation_min, validation_max,
		COALESCE(depends_on,''), COALESCE(depends_value,''),
		COALESCE(env_var_name,''), COALESCE(service_name,''),
		restart_required, display_order, is_active,
		created_at, updated_at
	FROM config_metadata
	WHERE is_active = true
	ORDER BY display_order ASC`

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanMetadataRows(rows)
}

// CountByCategory는 카테고리별 활성 설정 수를 반환합니다.
func (r *ConfigMetadataRepository) CountByCategory(ctx context.Context) (map[string]int32, error) {
	const q = `SELECT category::text, COUNT(*)::int
	FROM config_metadata
	WHERE is_active = true
	GROUP BY category`

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int32)
	for rows.Next() {
		var cat string
		var cnt int32
		if err := rows.Scan(&cat, &cnt); err != nil {
			return nil, err
		}
		counts[cat] = cnt
	}
	return counts, rows.Err()
}

// scanMetadataRows는 pgx.Rows에서 ConfigMetadata 슬라이스를 생성합니다.
func scanMetadataRows(rows pgx.Rows) ([]*service.ConfigMetadata, error) {
	var result []*service.ConfigMetadata
	for rows.Next() {
		var m service.ConfigMetadata
		var allowedValues []string
		if err := rows.Scan(
			&m.ConfigKey, &m.Category, &m.ValueType, &m.SecurityLevel,
			&m.IsRequired, &m.DefaultValue, &allowedValues,
			&m.ValidationRegex, &m.ValidationMin, &m.ValidationMax,
			&m.DependsOn, &m.DependsValue,
			&m.EnvVarName, &m.ServiceName,
			&m.RestartRequired, &m.DisplayOrder, &m.IsActive,
			&m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		m.AllowedValues = allowedValues
		result = append(result, &m)
	}
	return result, rows.Err()
}

// ============================================================================
// ConfigTranslationRepository — PostgreSQL 구현
// ============================================================================

// ConfigTranslationRepository는 PostgreSQL 기반 설정 번역 저장소입니다.
type ConfigTranslationRepository struct {
	pool *pgxpool.Pool
}

// NewConfigTranslationRepository는 PostgreSQL ConfigTranslationRepository를 생성합니다.
func NewConfigTranslationRepository(pool *pgxpool.Pool) *ConfigTranslationRepository {
	return &ConfigTranslationRepository{pool: pool}
}

// GetByKeyAndLang는 키와 언어 코드로 번역을 조회합니다.
func (r *ConfigTranslationRepository) GetByKeyAndLang(ctx context.Context, key, lang string) (*service.ConfigTranslation, error) {
	const q = `SELECT
		id::text, config_key, language_code,
		display_name, description,
		COALESCE(placeholder,''), COALESCE(help_text,''), COALESCE(validation_message,''),
		created_at, updated_at
	FROM config_translations
	WHERE config_key = $1 AND language_code = $2`

	var t service.ConfigTranslation
	err := r.pool.QueryRow(ctx, q, key, lang).Scan(
		&t.ID, &t.ConfigKey, &t.LanguageCode,
		&t.DisplayName, &t.Description,
		&t.Placeholder, &t.HelpText, &t.ValidationMessage,
		&t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

// ListByKey는 키로 모든 번역을 조회합니다.
func (r *ConfigTranslationRepository) ListByKey(ctx context.Context, key string) ([]*service.ConfigTranslation, error) {
	const q = `SELECT
		id::text, config_key, language_code,
		display_name, description,
		COALESCE(placeholder,''), COALESCE(help_text,''), COALESCE(validation_message,''),
		created_at, updated_at
	FROM config_translations
	WHERE config_key = $1
	ORDER BY language_code`

	rows, err := r.pool.Query(ctx, q, key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTranslationRows(rows)
}

// ListByLang는 언어 코드로 모든 번역을 조회합니다 (key→translation 맵).
func (r *ConfigTranslationRepository) ListByLang(ctx context.Context, lang string) (map[string]*service.ConfigTranslation, error) {
	const q = `SELECT
		id::text, config_key, language_code,
		display_name, description,
		COALESCE(placeholder,''), COALESCE(help_text,''), COALESCE(validation_message,''),
		created_at, updated_at
	FROM config_translations
	WHERE language_code = $1`

	rows, err := r.pool.Query(ctx, q, lang)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]*service.ConfigTranslation)
	for rows.Next() {
		var t service.ConfigTranslation
		if err := rows.Scan(
			&t.ID, &t.ConfigKey, &t.LanguageCode,
			&t.DisplayName, &t.Description,
			&t.Placeholder, &t.HelpText, &t.ValidationMessage,
			&t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result[t.ConfigKey] = &t
	}
	return result, rows.Err()
}

// scanTranslationRows는 pgx.Rows에서 ConfigTranslation 슬라이스를 생성합니다.
func scanTranslationRows(rows pgx.Rows) ([]*service.ConfigTranslation, error) {
	var result []*service.ConfigTranslation
	for rows.Next() {
		var t service.ConfigTranslation
		if err := rows.Scan(
			&t.ID, &t.ConfigKey, &t.LanguageCode,
			&t.DisplayName, &t.Description,
			&t.Placeholder, &t.HelpText, &t.ValidationMessage,
			&t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, &t)
	}
	return result, rows.Err()
}
