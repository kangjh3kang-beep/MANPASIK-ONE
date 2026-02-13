// Package memory는 인메모리 설정 메타데이터/번역 저장소입니다 (개발/테스트용).
package memory

import (
	"context"
	"sync"

	"github.com/manpasik/backend/services/admin-service/internal/service"
)

// ConfigMetadataRepository는 인메모리 설정 메타데이터 저장소입니다.
type ConfigMetadataRepository struct {
	mu    sync.RWMutex
	items map[string]*service.ConfigMetadata
}

// NewConfigMetadataRepository는 인메모리 ConfigMetadataRepository를 생성합니다.
func NewConfigMetadataRepository() *ConfigMetadataRepository {
	repo := &ConfigMetadataRepository{
		items: make(map[string]*service.ConfigMetadata),
	}

	// 시드 데이터: 주요 설정 메타
	seedItems := []*service.ConfigMetadata{
		{ConfigKey: "maintenance_mode", Category: "general", ValueType: "boolean", SecurityLevel: "public", DefaultValue: "false", ServiceName: "*", DisplayOrder: 1, IsActive: true},
		{ConfigKey: "max_devices_per_user", Category: "general", ValueType: "number", SecurityLevel: "public", DefaultValue: "5", ServiceName: "*", DisplayOrder: 2, IsActive: true},
		{ConfigKey: "default_language", Category: "general", ValueType: "select", SecurityLevel: "public", DefaultValue: "ko", AllowedValues: []string{"ko", "en", "ja", "zh", "fr", "hi"}, ServiceName: "*", DisplayOrder: 3, IsActive: true},
		{ConfigKey: "toss.secret_key", Category: "payment", ValueType: "secret", SecurityLevel: "secret", IsRequired: true, EnvVarName: "TOSS_SECRET_KEY", ServiceName: "payment-service", DisplayOrder: 1, IsActive: true},
		{ConfigKey: "toss.api_url", Category: "payment", ValueType: "url", SecurityLevel: "internal", DefaultValue: "https://api.tosspayments.com", EnvVarName: "TOSS_API_URL", ServiceName: "payment-service", DisplayOrder: 2, IsActive: true},
		{ConfigKey: "jwt.secret", Category: "auth", ValueType: "secret", SecurityLevel: "secret", IsRequired: true, EnvVarName: "JWT_SECRET", ServiceName: "auth-service", RestartRequired: true, DisplayOrder: 1, IsActive: true},
		{ConfigKey: "fcm.server_key", Category: "notification", ValueType: "secret", SecurityLevel: "secret", EnvVarName: "FCM_SERVER_KEY", ServiceName: "notification-service", DisplayOrder: 1, IsActive: true},
		{ConfigKey: "ai.default_model", Category: "ai", ValueType: "select", SecurityLevel: "public", DefaultValue: "biomarker_classifier", AllowedValues: []string{"biomarker_classifier", "anomaly_detector", "trend_predictor", "health_scorer", "food_calorie_estimator"}, ServiceName: "ai-inference-service", DisplayOrder: 1, IsActive: true},
	}

	for _, item := range seedItems {
		repo.items[item.ConfigKey] = item
	}

	return repo
}

// GetByKey는 키로 메타데이터를 조회합니다.
func (r *ConfigMetadataRepository) GetByKey(_ context.Context, key string) (*service.ConfigMetadata, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.items[key]
	if !ok {
		return nil, nil
	}
	cp := *item
	return &cp, nil
}

// ListByCategory는 카테고리별 메타데이터를 조회합니다.
func (r *ConfigMetadataRepository) ListByCategory(_ context.Context, category string) ([]*service.ConfigMetadata, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*service.ConfigMetadata
	for _, item := range r.items {
		if category != "" && item.Category != category {
			continue
		}
		if !item.IsActive {
			continue
		}
		cp := *item
		result = append(result, &cp)
	}
	return result, nil
}

// ListAll은 모든 메타데이터를 조회합니다.
func (r *ConfigMetadataRepository) ListAll(_ context.Context) ([]*service.ConfigMetadata, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*service.ConfigMetadata, 0, len(r.items))
	for _, item := range r.items {
		if !item.IsActive {
			continue
		}
		cp := *item
		result = append(result, &cp)
	}
	return result, nil
}

// CountByCategory는 카테고리별 설정 수를 반환합니다.
func (r *ConfigMetadataRepository) CountByCategory(_ context.Context) (map[string]int32, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	counts := make(map[string]int32)
	for _, item := range r.items {
		if item.IsActive {
			counts[item.Category]++
		}
	}
	return counts, nil
}

// ConfigTranslationRepository는 인메모리 설정 번역 저장소입니다.
type ConfigTranslationRepository struct {
	mu    sync.RWMutex
	items []*service.ConfigTranslation
}

// NewConfigTranslationRepository는 인메모리 ConfigTranslationRepository를 생성합니다.
func NewConfigTranslationRepository() *ConfigTranslationRepository {
	return &ConfigTranslationRepository{
		items: []*service.ConfigTranslation{
			{ConfigKey: "maintenance_mode", LanguageCode: "ko", DisplayName: "유지보수 모드", Description: "활성화하면 일반 사용자 접근이 차단됩니다."},
			{ConfigKey: "maintenance_mode", LanguageCode: "en", DisplayName: "Maintenance Mode", Description: "When enabled, regular users cannot access the service."},
			{ConfigKey: "toss.secret_key", LanguageCode: "ko", DisplayName: "Toss 시크릿 키", Description: "Toss Payments 결제 승인/취소 API 인증에 사용됩니다.", HelpText: "Toss 개발자센터에서 발급받으세요."},
			{ConfigKey: "toss.secret_key", LanguageCode: "en", DisplayName: "Toss Secret Key", Description: "Used for Toss Payments API authentication.", HelpText: "Obtain from Toss Developer Center."},
			{ConfigKey: "ai.default_model", LanguageCode: "ko", DisplayName: "기본 AI 모델", Description: "측정 분석에 사용할 기본 AI 모델입니다."},
			{ConfigKey: "ai.default_model", LanguageCode: "en", DisplayName: "Default AI Model", Description: "Default AI model for measurement analysis."},
		},
	}
}

// GetByKeyAndLang는 키와 언어 코드로 번역을 조회합니다.
func (r *ConfigTranslationRepository) GetByKeyAndLang(_ context.Context, key, lang string) (*service.ConfigTranslation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, t := range r.items {
		if t.ConfigKey == key && t.LanguageCode == lang {
			cp := *t
			return &cp, nil
		}
	}
	return nil, nil
}

// ListByKey는 키로 모든 번역을 조회합니다.
func (r *ConfigTranslationRepository) ListByKey(_ context.Context, key string) ([]*service.ConfigTranslation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*service.ConfigTranslation
	for _, t := range r.items {
		if t.ConfigKey == key {
			cp := *t
			result = append(result, &cp)
		}
	}
	return result, nil
}

// ListByLang는 언어 코드로 모든 번역을 조회합니다 (key→translation 맵).
func (r *ConfigTranslationRepository) ListByLang(_ context.Context, lang string) (map[string]*service.ConfigTranslation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]*service.ConfigTranslation)
	for _, t := range r.items {
		if t.LanguageCode == lang {
			cp := *t
			result[t.ConfigKey] = &cp
		}
	}
	return result, nil
}
