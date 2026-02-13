// Package config는 모든 마이크로서비스의 공통 설정을 제공합니다.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// ServiceConfig는 마이크로서비스 공통 설정입니다.
type ServiceConfig struct {
	// 서비스 정보
	ServiceName string
	Version     string
	GRPCPort    string
	HTTPPort    string

	// 데이터베이스
	DB DatabaseConfig

	// Redis 캐시
	Redis RedisConfig

	// Kafka 메시징
	Kafka KafkaConfig

	// Keycloak 인증
	Keycloak KeycloakConfig

	// JWT 설정
	JWT JWTConfig

	// Milvus 벡터 DB
	Milvus MilvusConfig

	// Elasticsearch 검색
	Elasticsearch ElasticsearchConfig

	// S3 호환 스토리지
	S3 S3Config

	// Toss Payments PG (payment-service)
	Toss TossConfig

	// 로그 레벨
	LogLevel string

	// 서버 설정
	ShutdownTimeout time.Duration
}

// DatabaseConfig는 PostgreSQL 접속 설정입니다.
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int
}

// DSN은 PostgreSQL 접속 문자열을 반환합니다.
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.SSLMode,
	)
}

// RedisConfig는 Redis 접속 설정입니다.
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// Addr은 Redis 접속 주소를 반환합니다.
func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

// KafkaConfig는 Kafka 접속 설정입니다.
type KafkaConfig struct {
	Brokers []string
	GroupID string
}

// KeycloakConfig는 Keycloak OIDC 설정입니다.
type KeycloakConfig struct {
	BaseURL  string
	Realm    string
	ClientID string
	Secret   string
}

// JWTConfig는 JWT 토큰 설정입니다.
type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Issuer          string
}

// MilvusConfig는 Milvus 벡터 DB 설정입니다.
type MilvusConfig struct {
	Host           string
	Port           int
	CollectionName string
}

// Addr은 Milvus 접속 주소를 반환합니다.
func (m *MilvusConfig) Addr() string {
	return fmt.Sprintf("%s:%d", m.Host, m.Port)
}

// ElasticsearchConfig는 Elasticsearch 설정입니다.
type ElasticsearchConfig struct {
	URL      string
	Username string
	Password string
}

// S3Config는 S3 호환 스토리지 설정입니다.
type S3Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
	UseSSL    bool
}

// TossConfig는 Toss Payments PG 설정입니다.
type TossConfig struct {
	SecretKey string // TOSS_SECRET_KEY (시크릿 키, 코드에 미기재)
	APIURL    string // TOSS_API_URL (기본: https://api.tosspayments.com)
}

// LoadFromEnv는 환경변수에서 설정을 로드합니다.
func LoadFromEnv(serviceName string) *ServiceConfig {
	return &ServiceConfig{
		ServiceName: serviceName,
		Version:     getEnv("VERSION", "dev"),
		GRPCPort:    getEnv("GRPC_PORT", ":50051"),
		HTTPPort:    getEnv("HTTP_PORT", ":8080"),

		DB: DatabaseConfig{
			Host:     getEnv("DB_HOST", "postgres"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "manpasik"),
			Password: getEnv("DB_PASSWORD", "manpasik_dev"),
			DBName:   getEnv("DB_NAME", "manpasik"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			MaxConns: getEnvInt("DB_MAX_CONNS", 20),
		},

		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "redis"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},

		Kafka: KafkaConfig{
			Brokers: []string{getEnv("KAFKA_BROKERS", "redpanda:19092")},
			GroupID: getEnv("KAFKA_GROUP_ID", serviceName),
		},

		Keycloak: KeycloakConfig{
			BaseURL:  getEnv("KEYCLOAK_URL", "http://keycloak:9090"),
			Realm:    getEnv("KEYCLOAK_REALM", "manpasik"),
			ClientID: getEnv("KEYCLOAK_CLIENT_ID", "manpasik-api"),
			Secret:   getEnv("KEYCLOAK_CLIENT_SECRET", ""),
		},

		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "dev-secret-change-in-production"),
			AccessTokenTTL:  time.Duration(getEnvInt("JWT_ACCESS_TTL_MINUTES", 15)) * time.Minute,
			RefreshTokenTTL: time.Duration(getEnvInt("JWT_REFRESH_TTL_DAYS", 7)) * 24 * time.Hour,
			Issuer:          getEnv("JWT_ISSUER", "manpasik-auth"),
		},

		Milvus: MilvusConfig{
			Host:           getEnv("MILVUS_HOST", "milvus"),
			Port:           getEnvInt("MILVUS_PORT", 19530),
			CollectionName: getEnv("MILVUS_COLLECTION", "manpasik_fingerprints"),
		},

		Elasticsearch: ElasticsearchConfig{
			URL:      getEnv("ELASTICSEARCH_URL", "http://elasticsearch:9200"),
			Username: getEnv("ELASTICSEARCH_USERNAME", ""),
			Password: getEnv("ELASTICSEARCH_PASSWORD", ""),
		},

		S3: S3Config{
			Endpoint:  getEnv("S3_ENDPOINT", "minio:9000"),
			AccessKey: getEnv("S3_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("S3_SECRET_KEY", "minioadmin"),
			Bucket:    getEnv("S3_BUCKET", "manpasik"),
			Region:    getEnv("S3_REGION", "us-east-1"),
			UseSSL:    getEnv("S3_USE_SSL", "false") == "true",
		},

		Toss: TossConfig{
			SecretKey: getEnv("TOSS_SECRET_KEY", ""),
			APIURL:    getEnv("TOSS_API_URL", "https://api.tosspayments.com"),
		},

		LogLevel:        getEnv("LOG_LEVEL", "info"),
		ShutdownTimeout: time.Duration(getEnvInt("SHUTDOWN_TIMEOUT_SECONDS", 5)) * time.Second,
	}
}

// getEnv는 환경변수를 조회하고 기본값을 반환합니다.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvInt는 정수형 환경변수를 조회합니다.
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
