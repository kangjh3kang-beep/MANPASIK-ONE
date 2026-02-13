package service_test

import (
	"context"
	"testing"

	"github.com/manpasik/backend/services/admin-service/internal/crypto"
	"github.com/manpasik/backend/services/admin-service/internal/repository/memory"
	"github.com/manpasik/backend/services/admin-service/internal/service"
	"github.com/manpasik/backend/shared/events"
	"go.uber.org/zap"
)

func newTestConfigManager() *service.ConfigManager {
	logger := zap.NewNop()
	configRepo := memory.NewSystemConfigRepository()
	metaRepo := memory.NewConfigMetadataRepository()
	transRepo := memory.NewConfigTranslationRepository()
	auditRepo := memory.NewAuditLogRepository()
	eventBus := events.NewEventBus()

	return service.NewConfigManager(logger, configRepo, metaRepo, transRepo, auditRepo, nil, eventBus)
}

// newTestConfigManagerWithBus는 이벤트 버스를 외부에서 참조할 수 있도록 반환합니다.
func newTestConfigManagerWithBus() (*service.ConfigManager, *events.EventBus) {
	logger := zap.NewNop()
	configRepo := memory.NewSystemConfigRepository()
	metaRepo := memory.NewConfigMetadataRepository()
	transRepo := memory.NewConfigTranslationRepository()
	auditRepo := memory.NewAuditLogRepository()
	eventBus := events.NewEventBus()

	cm := service.NewConfigManager(logger, configRepo, metaRepo, transRepo, auditRepo, nil, eventBus)
	return cm, eventBus
}

// newTestConfigManagerWithEncryption은 암호화 기능이 활성화된 ConfigManager를 생성합니다.
func newTestConfigManagerWithEncryption(t *testing.T) (*service.ConfigManager, *memory.SystemConfigRepository) {
	t.Helper()
	logger := zap.NewNop()
	configRepo := memory.NewSystemConfigRepository()
	metaRepo := memory.NewConfigMetadataRepository()
	transRepo := memory.NewConfigTranslationRepository()
	auditRepo := memory.NewAuditLogRepository()
	eventBus := events.NewEventBus()

	// 테스트용 암호화 키 (32 bytes = 64 hex chars)
	enc, err := crypto.NewAESEncryptor("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("암호화기 생성 실패: %v", err)
	}

	cm := service.NewConfigManager(logger, configRepo, metaRepo, transRepo, auditRepo, enc, eventBus)
	return cm, configRepo
}

func TestListSystemConfigs_All(t *testing.T) {
	cm := newTestConfigManager()
	ctx := context.Background()

	configs, counts, err := cm.ListSystemConfigs(ctx, "ko", "", false)
	if err != nil {
		t.Fatalf("ListSystemConfigs 실패: %v", err)
	}
	if len(configs) == 0 {
		t.Error("설정이 0개입니다")
	}
	if len(counts) == 0 {
		t.Error("카테고리 수가 0개입니다")
	}
}

func TestListSystemConfigs_ByCategory(t *testing.T) {
	cm := newTestConfigManager()
	ctx := context.Background()

	configs, _, err := cm.ListSystemConfigs(ctx, "ko", "general", false)
	if err != nil {
		t.Fatalf("카테고리별 조회 실패: %v", err)
	}
	for _, c := range configs {
		if c.Meta != nil && c.Meta.Category != "general" {
			t.Errorf("general 카테고리만 있어야 하는데 %s 발견", c.Meta.Category)
		}
	}
}

func TestListSystemConfigs_SecretMasking(t *testing.T) {
	cm := newTestConfigManager()
	ctx := context.Background()

	// toss.secret_key에 값 설정
	_, _ = cm.SetConfigWithMeta(ctx, "toss.secret_key", "test_secret_value", "admin")

	configs, _, _ := cm.ListSystemConfigs(ctx, "ko", "payment", false)
	for _, c := range configs {
		if c.Key == "toss.secret_key" {
			if c.Value != "****" {
				t.Errorf("시크릿은 마스킹되어야 합니다: %s", c.Value)
			}
			if c.RawValue != "" {
				t.Error("includeSecrets=false일 때 RawValue는 비어있어야 합니다")
			}
		}
	}
}

func TestGetConfigWithMeta(t *testing.T) {
	cm := newTestConfigManager()
	ctx := context.Background()

	cfg, err := cm.GetConfigWithMeta(ctx, "maintenance_mode", "ko")
	if err != nil {
		t.Fatalf("GetConfigWithMeta 실패: %v", err)
	}
	if cfg == nil {
		t.Fatal("cfg가 nil입니다")
	}
	if cfg.Key != "maintenance_mode" {
		t.Errorf("key가 다릅니다: %s", cfg.Key)
	}
	if cfg.Translation != nil && cfg.Translation.DisplayName != "유지보수 모드" {
		t.Errorf("한국어 번역이 다릅니다: %s", cfg.Translation.DisplayName)
	}
}

func TestGetConfigWithMeta_English(t *testing.T) {
	cm := newTestConfigManager()
	ctx := context.Background()

	cfg, err := cm.GetConfigWithMeta(ctx, "maintenance_mode", "en")
	if err != nil {
		t.Fatalf("영어 조회 실패: %v", err)
	}
	if cfg.Translation != nil && cfg.Translation.DisplayName != "Maintenance Mode" {
		t.Errorf("영어 번역이 다릅니다: %s", cfg.Translation.DisplayName)
	}
}

func TestValidateConfigValue_Number(t *testing.T) {
	cm := newTestConfigManager()
	ctx := context.Background()

	// number 유형은 없는 메타 → 유효
	result := cm.ValidateConfigValue(ctx, "nonexistent", "42")
	if !result.Valid {
		t.Errorf("메타 없으면 유효해야 합니다")
	}
}

func TestValidateConfigValue_Boolean(t *testing.T) {
	cm := newTestConfigManager()
	ctx := context.Background()

	// maintenance_mode는 boolean
	result := cm.ValidateConfigValue(ctx, "maintenance_mode", "true")
	if !result.Valid {
		t.Error("true는 유효해야 합니다")
	}

	result = cm.ValidateConfigValue(ctx, "maintenance_mode", "yes")
	if result.Valid {
		t.Error("yes는 유효하지 않아야 합니다")
	}
}

func TestValidateConfigValue_Select(t *testing.T) {
	cm := newTestConfigManager()
	ctx := context.Background()

	result := cm.ValidateConfigValue(ctx, "default_language", "ko")
	if !result.Valid {
		t.Error("ko는 유효해야 합니다")
	}

	result = cm.ValidateConfigValue(ctx, "default_language", "xx")
	if result.Valid {
		t.Error("xx는 유효하지 않아야 합니다")
	}
	if len(result.Suggestions) == 0 {
		t.Error("허용값 제안이 있어야 합니다")
	}
}

func TestValidateConfigValue_URL(t *testing.T) {
	cm := newTestConfigManager()
	ctx := context.Background()

	result := cm.ValidateConfigValue(ctx, "toss.api_url", "https://api.tosspayments.com")
	if !result.Valid {
		t.Error("유효한 URL이어야 합니다")
	}

	result = cm.ValidateConfigValue(ctx, "toss.api_url", "not-a-url")
	if result.Valid {
		t.Error("잘못된 URL은 유효하지 않아야 합니다")
	}
}

func TestSetConfigWithMeta_WithAuditAndEvent(t *testing.T) {
	cm, eventBus := newTestConfigManagerWithBus()
	ctx := context.Background()

	eventReceived := false
	eventBus.Subscribe(events.EventConfigChanged, func(_ context.Context, event events.Event) error {
		eventReceived = true
		if event.Payload["key"] != "maintenance_mode" {
			t.Errorf("이벤트 key가 다릅니다: %v", event.Payload["key"])
		}
		return nil
	})

	cfg, err := cm.SetConfigWithMeta(ctx, "maintenance_mode", "true", "test-admin")
	if err != nil {
		t.Fatalf("설정 저장 실패: %v", err)
	}
	if cfg == nil {
		t.Fatal("cfg가 nil입니다")
	}

	if !eventReceived {
		t.Error("config.changed 이벤트가 발생해야 합니다")
	}
}

func TestSetConfigWithMeta_ValidationFailure(t *testing.T) {
	cm := newTestConfigManager()
	ctx := context.Background()

	_, err := cm.SetConfigWithMeta(ctx, "maintenance_mode", "invalid", "admin")
	if err == nil {
		t.Error("유효성 검증 실패 시 에러가 발생해야 합니다")
	}
}

func TestSetConfigWithMeta_Encryption(t *testing.T) {
	cm, configRepo := newTestConfigManagerWithEncryption(t)
	ctx := context.Background()

	// toss.secret_key는 security_level=secret
	_, err := cm.SetConfigWithMeta(ctx, "toss.secret_key", "live_sk_test123", "admin")
	if err != nil {
		t.Fatalf("시크릿 설정 저장 실패: %v", err)
	}

	// DB에 저장된 값은 암호화되어야 함
	raw, _ := configRepo.FindByKey(ctx, "toss.secret_key")
	if raw.Value == "live_sk_test123" {
		t.Error("시크릿 값이 평문으로 저장되었습니다")
	}

	// GetConfigWithMeta로 복호화 확인
	cfgMeta, _ := cm.GetConfigWithMeta(ctx, "toss.secret_key", "ko")
	if cfgMeta.Value != "****" {
		t.Errorf("마스킹이 안 되었습니다: %s", cfgMeta.Value)
	}
	if cfgMeta.RawValue != "live_sk_test123" {
		t.Errorf("복호화된 값이 다릅니다: %s", cfgMeta.RawValue)
	}
}

func TestBulkSetConfigs(t *testing.T) {
	cm := newTestConfigManager()
	ctx := context.Background()

	configs := []struct{ Key, Value string }{
		{"maintenance_mode", "true"},
		{"default_language", "en"},
		{"default_language", "invalid"}, // 실패 케이스
	}

	successes, failures, errs := cm.BulkSetConfigs(ctx, configs, "admin", "일괄 변경 테스트")
	if successes != 2 {
		t.Errorf("성공 수가 2여야 합니다: %d", successes)
	}
	if failures != 1 {
		t.Errorf("실패 수가 1이어야 합니다: %d", failures)
	}
	if len(errs) != 1 {
		t.Errorf("에러 메시지가 1개여야 합니다: %d", len(errs))
	}
}
