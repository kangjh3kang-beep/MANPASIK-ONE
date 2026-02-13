// config_manager.go는 설정 관리 확장 비즈니스 로직입니다.
package service

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/manpasik/backend/services/admin-service/internal/crypto"
	"github.com/manpasik/backend/shared/events"
	"go.uber.org/zap"
)

// ConfigManager는 확장된 설정 관리 비즈니스 로직입니다.
type ConfigManager struct {
	logger          *zap.Logger
	configRepo      SystemConfigRepository
	metaRepo        ConfigMetadataRepository
	translationRepo ConfigTranslationRepository
	auditRepo       AuditLogRepository
	encryptor       *crypto.AESEncryptor
	eventPublisher  events.EventPublisher
}

// NewConfigManager는 ConfigManager를 생성합니다.
func NewConfigManager(
	logger *zap.Logger,
	configRepo SystemConfigRepository,
	metaRepo ConfigMetadataRepository,
	translationRepo ConfigTranslationRepository,
	auditRepo AuditLogRepository,
	encryptor *crypto.AESEncryptor,
	eventPublisher events.EventPublisher,
) *ConfigManager {
	return &ConfigManager{
		logger:          logger,
		configRepo:      configRepo,
		metaRepo:        metaRepo,
		translationRepo: translationRepo,
		auditRepo:       auditRepo,
		encryptor:       encryptor,
		eventPublisher:  eventPublisher,
	}
}

// ListSystemConfigs는 시스템 설정 목록을 메타데이터·번역과 함께 반환합니다.
func (cm *ConfigManager) ListSystemConfigs(ctx context.Context, lang, category string, includeSecrets bool) ([]*ConfigWithMeta, map[string]int32, error) {
	if lang == "" {
		lang = "ko"
	}

	// 메타데이터 조회
	var metas []*ConfigMetadata
	var err error
	if category != "" {
		metas, err = cm.metaRepo.ListByCategory(ctx, category)
	} else {
		metas, err = cm.metaRepo.ListAll(ctx)
	}
	if err != nil {
		cm.logger.Error("메타데이터 조회 실패", zap.Error(err))
		return nil, nil, fmt.Errorf("메타데이터 조회 실패: %w", err)
	}

	// 번역 맵 조회
	translations, err := cm.translationRepo.ListByLang(ctx, lang)
	if err != nil {
		cm.logger.Error("번역 조회 실패", zap.Error(err))
		translations = make(map[string]*ConfigTranslation)
	}

	// 카테고리별 수
	categoryCounts, _ := cm.metaRepo.CountByCategory(ctx)

	// 결과 구성
	results := make([]*ConfigWithMeta, 0, len(metas))
	for _, meta := range metas {
		cfg, _ := cm.configRepo.FindByKey(ctx, meta.ConfigKey)
		if cfg == nil {
			continue
		}

		item := &ConfigWithMeta{
			Key:       cfg.Key,
			Meta:      meta,
			UpdatedBy: cfg.UpdatedBy,
			UpdatedAt: cfg.UpdatedAt,
		}

		// 시크릿 처리
		if meta.SecurityLevel == "secret" {
			item.Value = "****"
			if includeSecrets {
				// 복호화
				if cm.encryptor != nil {
					decrypted, err := cm.encryptor.Decrypt(cfg.Value)
					if err == nil {
						item.RawValue = decrypted
					} else {
						item.RawValue = cfg.Value // 암호화 안 된 경우 그대로
					}
				} else {
					item.RawValue = cfg.Value
				}
			}
		} else {
			item.Value = cfg.Value
		}

		// 번역 첨부
		if trans, ok := translations[meta.ConfigKey]; ok {
			item.Translation = trans
		}

		results = append(results, item)
	}

	return results, categoryCounts, nil
}

// GetConfigWithMeta는 단일 설정을 메타데이터·번역과 함께 반환합니다.
func (cm *ConfigManager) GetConfigWithMeta(ctx context.Context, key, lang string) (*ConfigWithMeta, error) {
	if key == "" {
		return nil, fmt.Errorf("설정 key는 필수입니다")
	}
	if lang == "" {
		lang = "ko"
	}

	cfg, err := cm.configRepo.FindByKey(ctx, key)
	if err != nil || cfg == nil {
		return nil, fmt.Errorf("설정을 찾을 수 없습니다: %s", key)
	}

	meta, _ := cm.metaRepo.GetByKey(ctx, key)
	trans, _ := cm.translationRepo.GetByKeyAndLang(ctx, key, lang)

	item := &ConfigWithMeta{
		Key:         cfg.Key,
		Meta:        meta,
		Translation: trans,
		UpdatedBy:   cfg.UpdatedBy,
		UpdatedAt:   cfg.UpdatedAt,
	}

	if meta != nil && meta.SecurityLevel == "secret" {
		item.Value = "****"
		if cm.encryptor != nil {
			decrypted, err := cm.encryptor.Decrypt(cfg.Value)
			if err == nil {
				item.RawValue = decrypted
			} else {
				item.RawValue = cfg.Value
			}
		} else {
			item.RawValue = cfg.Value
		}
	} else {
		item.Value = cfg.Value
	}

	return item, nil
}

// ValidateConfigValue는 설정 값의 유효성을 검증합니다.
func (cm *ConfigManager) ValidateConfigValue(ctx context.Context, key, value string) *ValidateResult {
	meta, err := cm.metaRepo.GetByKey(ctx, key)
	if err != nil || meta == nil {
		// 메타데이터 없으면 검증 생략 (유효 처리)
		return &ValidateResult{Valid: true}
	}

	switch meta.ValueType {
	case "number":
		num, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return &ValidateResult{Valid: false, ErrorMsg: "숫자 형식이 아닙니다"}
		}
		if meta.ValidationMin != nil && num < *meta.ValidationMin {
			return &ValidateResult{Valid: false, ErrorMsg: fmt.Sprintf("최소값은 %v입니다", *meta.ValidationMin)}
		}
		if meta.ValidationMax != nil && num > *meta.ValidationMax {
			return &ValidateResult{Valid: false, ErrorMsg: fmt.Sprintf("최대값은 %v입니다", *meta.ValidationMax)}
		}

	case "boolean":
		if value != "true" && value != "false" {
			return &ValidateResult{Valid: false, ErrorMsg: "true 또는 false만 가능합니다", Suggestions: []string{"true", "false"}}
		}

	case "url":
		if _, err := url.ParseRequestURI(value); err != nil {
			return &ValidateResult{Valid: false, ErrorMsg: "올바른 URL 형식이 아닙니다"}
		}

	case "email":
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(value) {
			return &ValidateResult{Valid: false, ErrorMsg: "올바른 이메일 형식이 아닙니다"}
		}

	case "select":
		if len(meta.AllowedValues) > 0 {
			found := false
			for _, av := range meta.AllowedValues {
				if av == value {
					found = true
					break
				}
			}
			if !found {
				return &ValidateResult{
					Valid:       false,
					ErrorMsg:    fmt.Sprintf("허용되지 않는 값입니다: %s", value),
					Suggestions: meta.AllowedValues,
				}
			}
		}
	}

	// 정규식 검증
	if meta.ValidationRegex != "" {
		re, err := regexp.Compile(meta.ValidationRegex)
		if err == nil && !re.MatchString(value) {
			return &ValidateResult{Valid: false, ErrorMsg: "유효성 검증 패턴에 맞지 않습니다"}
		}
	}

	return &ValidateResult{Valid: true}
}

// SetConfigWithMeta는 검증·암호화·감사로그·이벤트 발행을 포함한 설정 저장입니다.
func (cm *ConfigManager) SetConfigWithMeta(ctx context.Context, key, value, adminID string) (*SystemConfig, error) {
	if key == "" {
		return nil, fmt.Errorf("설정 key는 필수입니다")
	}

	// 유효성 검증
	result := cm.ValidateConfigValue(ctx, key, value)
	if !result.Valid {
		return nil, fmt.Errorf("유효성 검증 실패: %s", result.ErrorMsg)
	}

	// 이전 값 조회
	oldCfg, _ := cm.configRepo.FindByKey(ctx, key)
	oldValue := ""
	if oldCfg != nil {
		oldValue = oldCfg.Value
	}

	// 암호화
	meta, _ := cm.metaRepo.GetByKey(ctx, key)
	storeValue := value
	if meta != nil && meta.SecurityLevel == "secret" && cm.encryptor != nil {
		encrypted, err := cm.encryptor.Encrypt(value)
		if err != nil {
			cm.logger.Error("암호화 실패", zap.String("key", key), zap.Error(err))
			return nil, fmt.Errorf("암호화 실패: %w", err)
		}
		storeValue = encrypted
	}

	// 저장
	now := time.Now().UTC()
	cfg := &SystemConfig{
		Key:         key,
		Value:       storeValue,
		Description: "",
		UpdatedBy:   adminID,
		UpdatedAt:   now,
	}
	if oldCfg != nil {
		cfg.Description = oldCfg.Description
	}

	if err := cm.configRepo.Save(ctx, cfg); err != nil {
		cm.logger.Error("설정 저장 실패", zap.Error(err))
		return nil, fmt.Errorf("설정 저장 실패: %w", err)
	}

	// 감사 로그
	if cm.auditRepo != nil {
		_ = cm.auditRepo.Save(ctx, &AuditLogEntry{
			EntryID:      uuid.New().String(),
			AdminID:      adminID,
			Action:       ActionConfigChange,
			ResourceType: "system_config",
			ResourceID:   key,
			Description:  fmt.Sprintf("설정 변경: %s", key),
			Timestamp:    now,
		})
	}

	// 이벤트 발행
	if cm.eventPublisher != nil {
		serviceName := ""
		if meta != nil {
			serviceName = meta.ServiceName
		}
		_ = events.PublishConfigChanged(ctx, cm.eventPublisher, key, value, oldValue, adminID, serviceName)
	}

	cm.logger.Info("설정 변경 완료", zap.String("key", key), zap.String("admin_id", adminID))
	return cfg, nil
}

// BulkSetConfigs는 여러 설정을 일괄 변경합니다.
func (cm *ConfigManager) BulkSetConfigs(ctx context.Context, configs []struct{ Key, Value string }, adminID, reason string) (successes, failures int, errors []string) {
	for _, c := range configs {
		_, err := cm.SetConfigWithMeta(ctx, c.Key, c.Value, adminID)
		if err != nil {
			failures++
			errors = append(errors, fmt.Sprintf("%s: %v", c.Key, err))
		} else {
			successes++
		}
	}

	// 일괄 변경 감사 로그
	if cm.auditRepo != nil && reason != "" {
		_ = cm.auditRepo.Save(ctx, &AuditLogEntry{
			EntryID:      uuid.New().String(),
			AdminID:      adminID,
			Action:       ActionConfigChange,
			ResourceType: "system_config",
			ResourceID:   "bulk",
			Description:  fmt.Sprintf("일괄 설정 변경 (성공: %d, 실패: %d): %s", successes, failures, reason),
			Timestamp:    time.Now().UTC(),
		})
	}

	return successes, failures, errors
}

// maskValue는 시크릿 값을 마스킹합니다.
func maskValue(val string) string {
	if val == "" {
		return ""
	}
	if len(val) <= 4 {
		return strings.Repeat("*", len(val))
	}
	return val[:2] + strings.Repeat("*", len(val)-4) + val[len(val)-2:]
}
