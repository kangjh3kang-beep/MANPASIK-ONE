// admin_config_ext.go는 AdminService 설정 관리 확장 메시지 수동 정의입니다.
// NOTE: protoc 재생성 시 manpasik.pb.go에 포함되면 이 파일을 제거합니다.
package v1

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ListSystemConfigsRequest는 설정 목록 조회 요청입니다.
type ListSystemConfigsRequest struct {
	LanguageCode   string `protobuf:"bytes,1,opt,name=language_code,json=languageCode,proto3" json:"language_code,omitempty"`
	Category       string `protobuf:"bytes,2,opt,name=category,proto3" json:"category,omitempty"`
	IncludeSecrets bool   `protobuf:"varint,3,opt,name=include_secrets,json=includeSecrets,proto3" json:"include_secrets,omitempty"`
}

func (x *ListSystemConfigsRequest) GetLanguageCode() string {
	if x != nil {
		return x.LanguageCode
	}
	return ""
}
func (x *ListSystemConfigsRequest) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}
func (x *ListSystemConfigsRequest) GetIncludeSecrets() bool {
	if x != nil {
		return x.IncludeSecrets
	}
	return false
}

// ListSystemConfigsResponse는 설정 목록 조회 응답입니다.
type ListSystemConfigsResponse struct {
	Configs        []*ConfigWithMeta `protobuf:"bytes,1,rep,name=configs,proto3" json:"configs,omitempty"`
	CategoryCounts map[string]int32  `protobuf:"bytes,2,rep,name=category_counts,json=categoryCounts,proto3" json:"category_counts,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *ListSystemConfigsResponse) GetConfigs() []*ConfigWithMeta {
	if x != nil {
		return x.Configs
	}
	return nil
}
func (x *ListSystemConfigsResponse) GetCategoryCounts() map[string]int32 {
	if x != nil {
		return x.CategoryCounts
	}
	return nil
}

// ConfigWithMeta는 설정 값 + 메타데이터 + 번역입니다.
type ConfigWithMeta struct {
	Key              string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value            string                 `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	RawValue         string                 `protobuf:"bytes,3,opt,name=raw_value,json=rawValue,proto3" json:"raw_value,omitempty"`
	Category         string                 `protobuf:"bytes,4,opt,name=category,proto3" json:"category,omitempty"`
	ValueType        string                 `protobuf:"bytes,5,opt,name=value_type,json=valueType,proto3" json:"value_type,omitempty"`
	SecurityLevel    string                 `protobuf:"bytes,6,opt,name=security_level,json=securityLevel,proto3" json:"security_level,omitempty"`
	IsRequired       bool                   `protobuf:"varint,7,opt,name=is_required,json=isRequired,proto3" json:"is_required,omitempty"`
	DefaultValue     string                 `protobuf:"bytes,8,opt,name=default_value,json=defaultValue,proto3" json:"default_value,omitempty"`
	AllowedValues    []string               `protobuf:"bytes,9,rep,name=allowed_values,json=allowedValues,proto3" json:"allowed_values,omitempty"`
	ValidationRegex  string                 `protobuf:"bytes,10,opt,name=validation_regex,json=validationRegex,proto3" json:"validation_regex,omitempty"`
	ValidationMin    float64                `protobuf:"fixed64,11,opt,name=validation_min,json=validationMin,proto3" json:"validation_min,omitempty"`
	ValidationMax    float64                `protobuf:"fixed64,12,opt,name=validation_max,json=validationMax,proto3" json:"validation_max,omitempty"`
	DependsOn        string                 `protobuf:"bytes,13,opt,name=depends_on,json=dependsOn,proto3" json:"depends_on,omitempty"`
	DependsValue     string                 `protobuf:"bytes,14,opt,name=depends_value,json=dependsValue,proto3" json:"depends_value,omitempty"`
	EnvVarName       string                 `protobuf:"bytes,15,opt,name=env_var_name,json=envVarName,proto3" json:"env_var_name,omitempty"`
	ServiceName      string                 `protobuf:"bytes,16,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
	RestartRequired  bool                   `protobuf:"varint,17,opt,name=restart_required,json=restartRequired,proto3" json:"restart_required,omitempty"`
	DisplayName      string                 `protobuf:"bytes,20,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	Description      string                 `protobuf:"bytes,21,opt,name=description,proto3" json:"description,omitempty"`
	Placeholder      string                 `protobuf:"bytes,22,opt,name=placeholder,proto3" json:"placeholder,omitempty"`
	HelpText         string                 `protobuf:"bytes,23,opt,name=help_text,json=helpText,proto3" json:"help_text,omitempty"`
	ValidationMessage string                `protobuf:"bytes,24,opt,name=validation_message,json=validationMessage,proto3" json:"validation_message,omitempty"`
	UpdatedBy        string                 `protobuf:"bytes,30,opt,name=updated_by,json=updatedBy,proto3" json:"updated_by,omitempty"`
	UpdatedAt        *timestamppb.Timestamp `protobuf:"bytes,31,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

// GetConfigWithMetaRequest는 단일 설정 조회 요청입니다.
type GetConfigWithMetaRequest struct {
	Key          string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	LanguageCode string `protobuf:"bytes,2,opt,name=language_code,json=languageCode,proto3" json:"language_code,omitempty"`
}

func (x *GetConfigWithMetaRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}
func (x *GetConfigWithMetaRequest) GetLanguageCode() string {
	if x != nil {
		return x.LanguageCode
	}
	return ""
}

// ValidateConfigValueRequest는 설정 값 유효성 검증 요청입니다.
type ValidateConfigValueRequest struct {
	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *ValidateConfigValueRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}
func (x *ValidateConfigValueRequest) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

// ValidateConfigValueResponse는 유효성 검증 응답입니다.
type ValidateConfigValueResponse struct {
	Valid        bool     `protobuf:"varint,1,opt,name=valid,proto3" json:"valid,omitempty"`
	ErrorMessage string   `protobuf:"bytes,2,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
	Suggestions  []string `protobuf:"bytes,3,rep,name=suggestions,proto3" json:"suggestions,omitempty"`
}

// BulkSetConfigsRequest는 일괄 설정 변경 요청입니다.
type BulkSetConfigsRequest struct {
	Configs []*SetSystemConfigRequest `protobuf:"bytes,1,rep,name=configs,proto3" json:"configs,omitempty"`
	Reason  string                    `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
}

func (x *BulkSetConfigsRequest) GetConfigs() []*SetSystemConfigRequest {
	if x != nil {
		return x.Configs
	}
	return nil
}
func (x *BulkSetConfigsRequest) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

// BulkSetConfigsResponse는 일괄 설정 변경 응답입니다.
type BulkSetConfigsResponse struct {
	Results      []*ConfigChangeResult `protobuf:"bytes,1,rep,name=results,proto3" json:"results,omitempty"`
	SuccessCount int32                 `protobuf:"varint,2,opt,name=success_count,json=successCount,proto3" json:"success_count,omitempty"`
	FailureCount int32                 `protobuf:"varint,3,opt,name=failure_count,json=failureCount,proto3" json:"failure_count,omitempty"`
}

// ConfigChangeResult는 개별 설정 변경 결과입니다.
type ConfigChangeResult struct {
	Key          string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Success      bool   `protobuf:"varint,2,opt,name=success,proto3" json:"success,omitempty"`
	ErrorMessage string `protobuf:"bytes,3,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
}
