package config

import (
	"os"
	"testing"
)

func TestLoadConfigFromDB_NilPool(t *testing.T) {
	result := LoadConfigFromDB(nil, "any_key")
	if result != "" {
		t.Errorf("nil pool에서 빈 문자열이어야 합니다: %s", result)
	}
}

func TestLoadConfigFromDB_EmptyKey(t *testing.T) {
	result := LoadConfigFromDB(nil, "")
	if result != "" {
		t.Errorf("빈 키에서 빈 문자열이어야 합니다: %s", result)
	}
}

func TestLoadConfigWithFallback_EnvVar(t *testing.T) {
	// DB 없이 환경변수 fallback 테스트
	testKey := "TEST_CONFIG_FALLBACK_" + t.Name()
	testValue := "fallback_value_123"
	os.Setenv(testKey, testValue)
	defer os.Unsetenv(testKey)

	result := LoadConfigWithFallback(nil, "nonexistent.db.key", testKey)
	if result != testValue {
		t.Errorf("환경변수 fallback 실패: %s != %s", result, testValue)
	}
}

func TestLoadConfigWithFallback_NoEnv(t *testing.T) {
	result := LoadConfigWithFallback(nil, "nonexistent.key", "NONEXISTENT_ENV_VAR_12345")
	if result != "" {
		t.Errorf("DB도 env도 없으면 빈 문자열이어야 합니다: %s", result)
	}
}
