package compliance

import (
	"fmt"
	"time"
)

// RetentionPolicy는 데이터 유형별 보존 기간 정책을 정의합니다 (GDPR Art.5).
type RetentionPolicy struct {
	DataType        string        // 데이터 유형 (measurement, health_record, audit_log 등)
	RetentionPeriod time.Duration // 보존 기간
	LegalBasis      string        // 법적 근거
	Description     string        // 정책 설명
}

// RetentionCheckResult는 보존 기간 검증 결과입니다.
type RetentionCheckResult struct {
	DataType        string    // 데이터 유형
	RecordID        string    // 레코드 ID
	CreatedAt       time.Time // 생성일
	ExpiresAt       time.Time // 만료일
	IsExpired       bool      // 만료 여부
	DaysRemaining   int       // 남은 일수 (만료 시 음수)
	RetentionPolicy string    // 적용된 정책명
}

// DefaultRetentionPolicies는 만파식 기본 데이터 보존 정책을 반환합니다.
func DefaultRetentionPolicies() []RetentionPolicy {
	return []RetentionPolicy{
		{
			DataType:        "measurement",
			RetentionPeriod: 10 * 365 * 24 * time.Hour, // 10년
			LegalBasis:      "의료법 제22조 (진료기록 보존의무 10년)",
			Description:     "차동측정 원시 데이터 및 결과",
		},
		{
			DataType:        "health_record",
			RetentionPeriod: 10 * 365 * 24 * time.Hour, // 10년
			LegalBasis:      "의료법 제22조, HIPAA §164.530(j)",
			Description:     "건강 기록 (검사결과, 진단, 처방)",
		},
		{
			DataType:        "audit_log",
			RetentionPeriod: 6 * 365 * 24 * time.Hour, // 6년
			LegalBasis:      "GDPR Art.5(1)(e), HIPAA §164.530(j)",
			Description:     "시스템 감사 로그",
		},
		{
			DataType:        "consent",
			RetentionPeriod: 5 * 365 * 24 * time.Hour, // 5년
			LegalBasis:      "GDPR Art.7(1), 개인정보보호법 제39조의6",
			Description:     "데이터 공유 동의 기록 (철회 후 5년)",
		},
		{
			DataType:        "user_profile",
			RetentionPeriod: 3 * 365 * 24 * time.Hour, // 3년
			LegalBasis:      "개인정보보호법 제21조 (비활성 계정 파기)",
			Description:     "사용자 프로필 (비활성 기준)",
		},
		{
			DataType:        "payment",
			RetentionPeriod: 5 * 365 * 24 * time.Hour, // 5년
			LegalBasis:      "전자상거래법 제6조, 국세기본법 제26조의2",
			Description:     "결제 및 거래 기록",
		},
		{
			DataType:        "communication",
			RetentionPeriod: 3 * 365 * 24 * time.Hour, // 3년
			LegalBasis:      "전자상거래법 제6조 (소비자 불만 처리)",
			Description:     "상담, 문의, 알림 기록",
		},
	}
}

// RetentionManager는 데이터 보존 정책을 관리합니다.
type RetentionManager struct {
	policies map[string]RetentionPolicy
}

// NewRetentionManager는 기본 정책으로 관리자를 생성합니다.
func NewRetentionManager() *RetentionManager {
	rm := &RetentionManager{
		policies: make(map[string]RetentionPolicy),
	}
	for _, p := range DefaultRetentionPolicies() {
		rm.policies[p.DataType] = p
	}
	return rm
}

// GetPolicy는 데이터 유형에 해당하는 보존 정책을 반환합니다.
func (rm *RetentionManager) GetPolicy(dataType string) (*RetentionPolicy, error) {
	p, ok := rm.policies[dataType]
	if !ok {
		return nil, fmt.Errorf("no retention policy for data type: %s", dataType)
	}
	return &p, nil
}

// CheckRetention는 레코드의 보존 기간 만료 여부를 검증합니다.
func (rm *RetentionManager) CheckRetention(dataType, recordID string, createdAt time.Time) (*RetentionCheckResult, error) {
	policy, err := rm.GetPolicy(dataType)
	if err != nil {
		return nil, err
	}

	expiresAt := createdAt.Add(policy.RetentionPeriod)
	now := time.Now()
	daysRemaining := int(expiresAt.Sub(now).Hours() / 24)

	return &RetentionCheckResult{
		DataType:        dataType,
		RecordID:        recordID,
		CreatedAt:       createdAt,
		ExpiresAt:       expiresAt,
		IsExpired:       now.After(expiresAt),
		DaysRemaining:   daysRemaining,
		RetentionPolicy: policy.LegalBasis,
	}, nil
}

// ListPolicies는 모든 보존 정책을 반환합니다.
func (rm *RetentionManager) ListPolicies() []RetentionPolicy {
	result := make([]RetentionPolicy, 0, len(rm.policies))
	for _, p := range rm.policies {
		result = append(result, p)
	}
	return result
}
