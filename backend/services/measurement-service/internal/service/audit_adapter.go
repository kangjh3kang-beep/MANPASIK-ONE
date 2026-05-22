// Package service provides the audit client adapter for measurement-service.
//
// H3 사전인증: 측정 작업을 감사 로그로 추적합니다.
// audit-service의 내부 서비스 계층을 gRPC로 호출하거나,
// audit-service가 사용 불가능할 때는 로컬 로깅으로 폴백합니다 (H6 실패 격리).
package service

import (
	"context"

	"go.uber.org/zap"
)

// LoggingAuditClient는 감사 로그를 로컬 로거로 기록하는 폴백 구현입니다.
// audit-service gRPC 연결이 불가능할 때 사용됩니다.
type LoggingAuditClient struct {
	logger *zap.Logger
}

// NewLoggingAuditClient는 로깅 기반 감사 클라이언트를 생성합니다.
func NewLoggingAuditClient(logger *zap.Logger) *LoggingAuditClient {
	return &LoggingAuditClient{logger: logger}
}

// RecordAction은 감사 액션을 로컬 로그로 기록합니다.
func (c *LoggingAuditClient) RecordAction(
	ctx context.Context,
	adminID, action, resourceType, resourceID, description string,
) error {
	c.logger.Info("audit: action recorded (local)",
		zap.String("admin_id", adminID),
		zap.String("action", action),
		zap.String("resource_type", resourceType),
		zap.String("resource_id", resourceID),
		zap.String("description", description),
	)
	return nil
}

// GRPCAuditClient는 audit-service gRPC를 통해 감사 로그를 기록합니다.
// GetAuditLogDetails RPC가 조회 전용이므로, 현재는 AdminServiceClient를
// 통한 로깅만 지원합니다. 향후 RecordAction RPC 추가 시 교체 예정.
type GRPCAuditClient struct {
	logger  *zap.Logger
	fallback *LoggingAuditClient
}

// NewGRPCAuditClient는 gRPC 기반 감사 클라이언트를 생성합니다.
func NewGRPCAuditClient(logger *zap.Logger) *GRPCAuditClient {
	return &GRPCAuditClient{
		logger:  logger,
		fallback: NewLoggingAuditClient(logger),
	}
}

// RecordAction은 감사 액션을 audit-service로 전송합니다.
// 현재 proto에 RecordAction RPC가 없으므로 로컬 로깅으로 폴백합니다.
// TODO: audit-service proto에 RecordAction RPC 추가 후 gRPC 호출로 교체
func (c *GRPCAuditClient) RecordAction(
	ctx context.Context,
	adminID, action, resourceType, resourceID, description string,
) error {
	// 현재 audit-service proto에 RecordAction RPC가 미정의이므로
	// 구조화된 로컬 로깅으로 기록 (H6: 감사 서비스 불가 시 측정 흐름 미차단)
	return c.fallback.RecordAction(ctx, adminID, action, resourceType, resourceID, description)
}
