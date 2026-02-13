// Package errors는 표준화된 에러 응답을 제공합니다.
//
// 보안 원칙: 내부 스택 트레이스/시스템 정보를 절대 노출하지 않습니다.
package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorCode는 만파식 표준 에러 코드입니다.
type ErrorCode string

const (
	// 인증/인가 에러
	ErrUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrForbidden    ErrorCode = "FORBIDDEN"
	ErrTokenExpired ErrorCode = "TOKEN_EXPIRED"
	ErrInvalidToken ErrorCode = "INVALID_TOKEN"

	// 입력 검증 에러
	ErrInvalidInput     ErrorCode = "INVALID_INPUT"
	ErrMissingField     ErrorCode = "MISSING_FIELD"
	ErrValidationFailed ErrorCode = "VALIDATION_FAILED"

	// 리소스 에러
	ErrNotFound      ErrorCode = "NOT_FOUND"
	ErrAlreadyExists ErrorCode = "ALREADY_EXISTS"
	ErrConflict      ErrorCode = "CONFLICT"

	// 서버 에러 (내부 상세는 숨김)
	ErrInternal           ErrorCode = "INTERNAL_ERROR"
	ErrServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"

	// 비즈니스 에러
	ErrDeviceLimitExceeded  ErrorCode = "DEVICE_LIMIT_EXCEEDED"
	ErrSubscriptionRequired ErrorCode = "SUBSCRIPTION_REQUIRED"
	ErrCartridgeExpired     ErrorCode = "CARTRIDGE_EXPIRED"
	ErrMeasurementFailed    ErrorCode = "MEASUREMENT_FAILED"
)

// AppError는 만파식 애플리케이션 에러입니다.
type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// New는 새 AppError를 생성합니다.
func New(code ErrorCode, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// WithDetails는 상세 정보를 추가합니다.
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// ToGRPC는 AppError를 gRPC Status로 변환합니다.
func (e *AppError) ToGRPC() error {
	return status.Error(e.grpcCode(), e.Message)
}

// grpcCode는 ErrorCode를 gRPC 코드로 매핑합니다.
func (e *AppError) grpcCode() codes.Code {
	switch e.Code {
	case ErrUnauthorized, ErrTokenExpired, ErrInvalidToken:
		return codes.Unauthenticated
	case ErrForbidden:
		return codes.PermissionDenied
	case ErrInvalidInput, ErrMissingField, ErrValidationFailed:
		return codes.InvalidArgument
	case ErrNotFound:
		return codes.NotFound
	case ErrAlreadyExists, ErrConflict:
		return codes.AlreadyExists
	case ErrServiceUnavailable:
		return codes.Unavailable
	case ErrDeviceLimitExceeded, ErrSubscriptionRequired, ErrCartridgeExpired:
		return codes.FailedPrecondition
	default:
		return codes.Internal
	}
}

// FromGRPC는 gRPC Status를 AppError로 변환합니다.
func FromGRPC(err error) *AppError {
	st, ok := status.FromError(err)
	if !ok {
		return New(ErrInternal, "알 수 없는 에러")
	}

	code := ErrInternal
	switch st.Code() {
	case codes.Unauthenticated:
		code = ErrUnauthorized
	case codes.PermissionDenied:
		code = ErrForbidden
	case codes.InvalidArgument:
		code = ErrInvalidInput
	case codes.NotFound:
		code = ErrNotFound
	case codes.AlreadyExists:
		code = ErrAlreadyExists
	}

	return New(code, st.Message())
}
