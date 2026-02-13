// config_models.go는 설정 관리 확장에 필요한 도메인 모델과 리포지토리 인터페이스를 정의합니다.
package service

import (
	"context"
	"time"
)

// ConfigMetadata는 설정 항목의 메타데이터(스키마)입니다.
type ConfigMetadata struct {
	ConfigKey       string
	Category        string // general, payment, auth, ...
	ValueType       string // string, number, boolean, secret, url, email, json, select, multiline
	SecurityLevel   string // public, internal, confidential, secret
	IsRequired      bool
	DefaultValue    string
	AllowedValues   []string
	ValidationRegex string
	ValidationMin   *float64
	ValidationMax   *float64
	DependsOn       string
	DependsValue    string
	EnvVarName      string
	ServiceName     string
	RestartRequired bool
	DisplayOrder    int
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ConfigTranslation은 설정의 다국어 번역입니다.
type ConfigTranslation struct {
	ID                string
	ConfigKey         string
	LanguageCode      string
	DisplayName       string
	Description       string
	Placeholder       string
	HelpText          string
	ValidationMessage string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// ConfigWithMeta는 설정 값 + 메타데이터 + 번역을 합친 뷰 모델입니다.
type ConfigWithMeta struct {
	// 값
	Key      string
	Value    string // secret이면 마스킹
	RawValue string // include_secrets=true일 때만

	// 메타데이터
	Meta *ConfigMetadata

	// 번역
	Translation *ConfigTranslation

	// 변경 정보
	UpdatedBy string
	UpdatedAt time.Time
}

// ValidateResult는 설정 값 검증 결과입니다.
type ValidateResult struct {
	Valid       bool
	ErrorMsg    string
	Suggestions []string
}

// ConfigMetadataRepository는 설정 메타데이터 저장소 인터페이스입니다.
type ConfigMetadataRepository interface {
	GetByKey(ctx context.Context, key string) (*ConfigMetadata, error)
	ListByCategory(ctx context.Context, category string) ([]*ConfigMetadata, error)
	ListAll(ctx context.Context) ([]*ConfigMetadata, error)
	CountByCategory(ctx context.Context) (map[string]int32, error)
}

// ConfigTranslationRepository는 설정 다국어 번역 저장소 인터페이스입니다.
type ConfigTranslationRepository interface {
	GetByKeyAndLang(ctx context.Context, key, lang string) (*ConfigTranslation, error)
	ListByKey(ctx context.Context, key string) ([]*ConfigTranslation, error)
	ListByLang(ctx context.Context, lang string) (map[string]*ConfigTranslation, error)
}
