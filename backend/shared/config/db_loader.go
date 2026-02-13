// db_loader.go는 DB system_configs에서 설정을 로드하는 유틸리티입니다.
package config

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// LoadConfigFromDB는 system_configs 테이블에서 설정 값을 조회합니다.
// 키가 없거나 DB 오류 시 빈 문자열을 반환합니다.
func LoadConfigFromDB(pool *pgxpool.Pool, key string) string {
	if pool == nil || key == "" {
		return ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var value string
	err := pool.QueryRow(ctx, "SELECT value FROM system_configs WHERE key=$1", key).Scan(&value)
	if err != nil {
		return ""
	}
	return value
}

// LoadConfigWithFallback는 DB에서 설정을 먼저 조회하고,
// 없으면 환경변수에서 폴백합니다.
func LoadConfigWithFallback(pool *pgxpool.Pool, dbKey, envVar string) string {
	if val := LoadConfigFromDB(pool, dbKey); val != "" {
		return val
	}
	return os.Getenv(envVar)
}
