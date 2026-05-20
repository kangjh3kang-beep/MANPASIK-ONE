package compliance

import (
	"fmt"
	"time"
)

// ComplianceChecker는 시스템 상태를 기반으로 규제 준수를 자동 검사합니다.
// GetComplianceOverview()의 하드코딩을 대체합니다.
type ComplianceChecker struct {
	systemState SystemState
}

// SystemState는 자동 검사에 필요한 시스템 현황입니다.
// 각 서비스가 시작 시 또는 주기적으로 갱신합니다.
type SystemState struct {
	// 암호화
	PHIEncryptionEnabled bool   // AES-256-GCM PHI 암호화 활성화 여부
	EncryptionAlgorithm  string // 사용 중인 암호화 알고리즘

	// 감사
	AuditLogEnabled    bool // 감사 로그 기록 활성화
	AuditSigningEnabled bool // 감사 로그 전자서명 활성화
	AuditChainEnabled  bool // 해시 체인 무결성 활성화

	// 접근 제어
	RBACEnabled       bool // 역할 기반 접근 제어
	JWTAuthEnabled    bool // JWT 인증
	SessionTimeout    time.Duration // 세션 타임아웃

	// 데이터 보호
	ConsentManagementEnabled bool     // 동의 관리 시스템
	RetentionPoliciesCount   int      // 등록된 보존 정책 수
	DeletionProcessorReady   bool     // 삭제 프로세서 등록 완료
	RegisteredDeleters       []string // 등록된 DataDeleter 유형

	// 문서
	SRSVersion     string // IEC 62304 SRS 버전
	RiskPlanReady  bool   // ISO 14971 위험관리 계획
	FMEAReady      bool   // FMEA 완료
	SafetyClass    string // IEC 62304 안전 분류 (A/B/C)

	// 외부 인증
	ClinicalTrialDone bool // 임상시험 완료
	GMPCertified      bool // ISO 13485 GMP 인증
	FDAPreSubmission  bool // FDA Pre-submission 완료
}

// NewComplianceChecker는 시스템 상태로 자동 검사기를 생성합니다.
func NewComplianceChecker(state SystemState) *ComplianceChecker {
	return &ComplianceChecker{systemState: state}
}

// RunChecks는 모든 규제 프레임워크에 대해 자동 검사를 수행합니다.
func (c *ComplianceChecker) RunChecks() []ComplianceStatus {
	var results []ComplianceStatus
	results = append(results, c.checkGDPR()...)
	results = append(results, c.checkHIPAA()...)
	results = append(results, c.checkMFDS()...)
	results = append(results, c.checkFDA()...)
	results = append(results, c.checkCEIVDR()...)
	results = append(results, c.checkCFR11()...)
	return results
}

// ComputeScore는 전체 준수율을 계산합니다 (0~100).
func ComputeScore(statuses []ComplianceStatus) ComplianceScore {
	if len(statuses) == 0 {
		return ComplianceScore{}
	}

	score := ComplianceScore{TotalItems: len(statuses)}
	for _, s := range statuses {
		switch s.Status {
		case "compliant":
			score.Compliant++
		case "partial":
			score.Partial++
		case "non_compliant":
			score.NonCompliant++
		}
	}

	// compliant=100%, partial=50%, non_compliant=0%
	score.Rate = float64(score.Compliant*100+score.Partial*50) / float64(score.TotalItems)
	return score
}

// ComplianceScore는 종합 준수율입니다.
type ComplianceScore struct {
	TotalItems   int
	Compliant    int
	Partial      int
	NonCompliant int
	Rate         float64 // 0~100
}

func (c *ComplianceChecker) checkGDPR() []ComplianceStatus {
	now := time.Now()
	s := c.systemState
	return []ComplianceStatus{
		{
			Framework: "GDPR", Category: "consent_management",
			Status:      boolStatus(s.ConsentManagementEnabled),
			Description: "데이터 공유 동의 관리 (Art.6, Art.9)",
			LastCheckedAt: now,
		},
		{
			Framework: "GDPR", Category: "phi_encryption",
			Status:      boolStatus(s.PHIEncryptionEnabled && s.EncryptionAlgorithm == "AES-256-GCM"),
			Description: "AES-256-GCM 필드 레벨 암호화 (Art.32)",
			LastCheckedAt: now,
			Recommendations: conditionalRec(!s.PHIEncryptionEnabled,
				"PHI_ENCRYPTION_KEY 환경변수 설정 및 암호화 활성화 필요"),
		},
		{
			Framework: "GDPR", Category: "data_retention",
			Status:      thresholdStatus(s.RetentionPoliciesCount, 5),
			Description: fmt.Sprintf("데이터 보존 정책 %d개 등록 (Art.5)", s.RetentionPoliciesCount),
			LastCheckedAt: now,
		},
		{
			Framework: "GDPR", Category: "right_to_erasure",
			Status:      boolStatus(s.DeletionProcessorReady && len(s.RegisteredDeleters) > 0),
			Description: fmt.Sprintf("삭제권 처리기 (%d개 DataDeleter 등록) (Art.17)", len(s.RegisteredDeleters)),
			LastCheckedAt: now,
			Recommendations: conditionalRec(len(s.RegisteredDeleters) == 0,
				"서비스별 DataDeleter 구현 및 등록 필요"),
		},
		{
			Framework: "GDPR", Category: "audit_trail",
			Status:      boolStatus(s.AuditLogEnabled),
			Description: "감사 로그 기록 (Art.5(2) 책임 원칙)",
			LastCheckedAt: now,
		},
	}
}

func (c *ComplianceChecker) checkHIPAA() []ComplianceStatus {
	now := time.Now()
	s := c.systemState
	return []ComplianceStatus{
		{
			Framework: "HIPAA", Category: "access_control",
			Status:      boolStatus(s.RBACEnabled && s.JWTAuthEnabled),
			Description: "JWT + RBAC 접근 제어 (§164.312(a)(1))",
			LastCheckedAt: now,
		},
		{
			Framework: "HIPAA", Category: "audit_controls",
			Status:      boolStatus(s.AuditLogEnabled && s.AuditSigningEnabled),
			Description: "감사 로그 + 전자서명 (§164.312(b))",
			LastCheckedAt: now,
			Recommendations: conditionalRec(!s.AuditSigningEnabled,
				"HMAC-SHA256 기반 감사 로그 서명 활성화 필요"),
		},
		{
			Framework: "HIPAA", Category: "integrity",
			Status:      boolStatus(s.AuditChainEnabled),
			Description: "해시 체인 기반 무결성 검증 (§164.312(c)(1))",
			LastCheckedAt: now,
			Recommendations: conditionalRec(!s.AuditChainEnabled,
				"AuditChain 해시 체인 활성화 필요"),
		},
		{
			Framework: "HIPAA", Category: "encryption",
			Status:      boolStatus(s.PHIEncryptionEnabled),
			Description: "PHI 암호화 (§164.312(a)(2)(iv))",
			LastCheckedAt: now,
		},
	}
}

func (c *ComplianceChecker) checkMFDS() []ComplianceStatus {
	now := time.Now()
	s := c.systemState
	return []ComplianceStatus{
		{
			Framework: "MFDS", Category: "software_classification",
			Status:      boolStatus(s.SafetyClass != ""),
			Description: fmt.Sprintf("IEC 62304 Class %s 분류", s.SafetyClass),
			LastCheckedAt: now,
		},
		{
			Framework: "MFDS", Category: "risk_management",
			Status:      boolStatus(s.RiskPlanReady && s.FMEAReady),
			Description: "ISO 14971 위험관리 + FMEA",
			LastCheckedAt: now,
		},
		{
			Framework: "MFDS", Category: "software_requirements",
			Status:      boolStatus(s.SRSVersion != ""),
			Description: fmt.Sprintf("IEC 62304 SRS %s", s.SRSVersion),
			LastCheckedAt: now,
		},
		{
			Framework: "MFDS", Category: "clinical_validation",
			Status:      boolStatus(s.ClinicalTrialDone),
			Description: "임상시험 (IRB 승인 + 100명 이상)",
			LastCheckedAt: now,
			Recommendations: conditionalRec(!s.ClinicalTrialDone,
				"IRB 승인 후 임상시험 실시 필요 (blocking)"),
		},
		{
			Framework: "MFDS", Category: "gmp_certification",
			Status:      boolStatus(s.GMPCertified),
			Description: "ISO 13485 GMP 인증",
			LastCheckedAt: now,
			Recommendations: conditionalRec(!s.GMPCertified,
				"ISO 13485 인증 획득 필요"),
		},
	}
}

func (c *ComplianceChecker) checkFDA() []ComplianceStatus {
	now := time.Now()
	s := c.systemState
	return []ComplianceStatus{
		{
			Framework: "FDA", Category: "predicate_device",
			Status:      "compliant", // 조사 완료
			Description: "Predicate Device 조사 완료",
			LastCheckedAt: now,
		},
		{
			Framework: "FDA", Category: "pre_submission",
			Status:      boolStatus(s.FDAPreSubmission),
			Description: "FDA Pre-submission Meeting",
			LastCheckedAt: now,
			Recommendations: conditionalRec(!s.FDAPreSubmission,
				"FDA 사전 상담 일정 확보 필요"),
		},
	}
}

func (c *ComplianceChecker) checkCEIVDR() []ComplianceStatus {
	now := time.Now()
	return []ComplianceStatus{
		{
			Framework: "CE-IVDR", Category: "technical_documentation",
			Status:      "partial",
			Description: "기술 문서 구조 정의 (Annex I)",
			LastCheckedAt: now,
			Recommendations: []string{"Notified Body 선정 및 컨설팅 필요"},
		},
	}
}

func (c *ComplianceChecker) checkCFR11() []ComplianceStatus {
	now := time.Now()
	s := c.systemState
	return []ComplianceStatus{
		{
			Framework: "21CFR11", Category: "electronic_signatures",
			Status:      boolStatus(s.AuditSigningEnabled),
			Description: "전자 서명 (§11.100)",
			LastCheckedAt: now,
		},
		{
			Framework: "21CFR11", Category: "audit_trail_protection",
			Status:      boolStatus(s.AuditChainEnabled && s.AuditSigningEnabled),
			Description: "감사 로그 보호 — 해시 체인 + 서명 (§11.10(e))",
			LastCheckedAt: now,
		},
		{
			Framework: "21CFR11", Category: "access_controls",
			Status:      boolStatus(s.RBACEnabled && s.JWTAuthEnabled && s.SessionTimeout > 0),
			Description: fmt.Sprintf("접근 제어 + 세션 타임아웃 %v (§11.10(d))", s.SessionTimeout),
			LastCheckedAt: now,
		},
	}
}

// boolStatus는 boolean을 compliance 상태 문자열로 변환합니다.
func boolStatus(ok bool) string {
	if ok {
		return "compliant"
	}
	return "non_compliant"
}

// thresholdStatus는 값이 최소 기준 이상이면 compliant, 0보다 크면 partial, 아니면 non_compliant입니다.
func thresholdStatus(value, min int) string {
	if value >= min {
		return "compliant"
	}
	if value > 0 {
		return "partial"
	}
	return "non_compliant"
}

// conditionalRec는 조건이 true일 때만 권고사항을 반환합니다.
func conditionalRec(condition bool, rec string) []string {
	if condition {
		return []string{rec}
	}
	return nil
}
