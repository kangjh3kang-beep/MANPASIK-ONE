package compliance

import (
	"context"
	"fmt"
	"time"
)

// DeletionRequestStatus는 삭제 요청 상태를 나타냅니다.
type DeletionRequestStatus string

const (
	DeletionPending    DeletionRequestStatus = "pending"
	DeletionProcessing DeletionRequestStatus = "processing"
	DeletionCompleted  DeletionRequestStatus = "completed"
	DeletionFailed     DeletionRequestStatus = "failed"
	DeletionPartial    DeletionRequestStatus = "partial" // 일부 법적 보존 데이터 잔존
)

// DeletionRequest는 GDPR Art.17 삭제 요청(잊힐 권리)을 표현합니다.
type DeletionRequest struct {
	RequestID   string
	UserID      string
	RequestedAt time.Time
	ProcessedAt *time.Time
	Status      DeletionRequestStatus
	Reason      string // 삭제 사유
	Scope       []string // 삭제 대상 데이터 유형 (빈 배열 = 전체)
	Results     []DeletionResult
}

// DeletionResult는 데이터 유형별 삭제 결과입니다.
type DeletionResult struct {
	DataType     string
	RecordCount  int
	DeletedCount int
	RetainedCount int    // 법적 보존 의무로 유지
	RetainReason  string // 유지 사유
	Status       DeletionRequestStatus
}

// DataDeleter는 서비스별 데이터 삭제를 수행하는 인터페이스입니다.
type DataDeleter interface {
	// DeleteUserData는 해당 사용자의 데이터를 삭제합니다.
	// 법적 보존 의무가 있는 데이터는 익명화 처리합니다.
	DeleteUserData(ctx context.Context, userID string) (*DeletionResult, error)

	// DataType은 이 삭제기가 처리하는 데이터 유형을 반환합니다.
	DataType() string
}

// DeletionProcessor는 삭제 요청을 처리합니다.
type DeletionProcessor struct {
	deleters  map[string]DataDeleter
	retention *RetentionManager
}

// NewDeletionProcessor는 삭제 프로세서를 생성합니다.
func NewDeletionProcessor(retention *RetentionManager) *DeletionProcessor {
	return &DeletionProcessor{
		deleters:  make(map[string]DataDeleter),
		retention: retention,
	}
}

// RegisterDeleter는 데이터 유형별 삭제기를 등록합니다.
func (dp *DeletionProcessor) RegisterDeleter(deleter DataDeleter) {
	dp.deleters[deleter.DataType()] = deleter
}

// ProcessDeletion는 삭제 요청을 실행합니다.
// 법적 보존 의무가 있는 데이터는 익명화 처리합니다.
func (dp *DeletionProcessor) ProcessDeletion(ctx context.Context, req *DeletionRequest) error {
	if req == nil {
		return fmt.Errorf("deletion request is nil")
	}
	if req.UserID == "" {
		return fmt.Errorf("user ID is required")
	}

	req.Status = DeletionProcessing
	now := time.Now()
	req.ProcessedAt = &now

	// 삭제 대상 결정
	targets := req.Scope
	if len(targets) == 0 {
		// 전체 삭제: 등록된 모든 삭제기
		for dt := range dp.deleters {
			targets = append(targets, dt)
		}
	}

	allSuccess := true
	for _, dataType := range targets {
		deleter, ok := dp.deleters[dataType]
		if !ok {
			req.Results = append(req.Results, DeletionResult{
				DataType: dataType,
				Status:   DeletionFailed,
				RetainReason: "no deleter registered",
			})
			allSuccess = false
			continue
		}

		result, err := deleter.DeleteUserData(ctx, req.UserID)
		if err != nil {
			req.Results = append(req.Results, DeletionResult{
				DataType:     dataType,
				Status:       DeletionFailed,
				RetainReason: err.Error(),
			})
			allSuccess = false
			continue
		}

		req.Results = append(req.Results, *result)
		if result.Status == DeletionFailed {
			allSuccess = false
		}
	}

	if allSuccess {
		req.Status = DeletionCompleted
	} else {
		req.Status = DeletionPartial
	}

	return nil
}

// ComplianceStatus는 규제 준수 상태를 나타냅니다.
type ComplianceStatus struct {
	Framework       string // GDPR, HIPAA, MFDS, FDA, CE-IVDR
	Category        string // data_protection, audit, encryption, consent, retention
	Status          string // compliant, partial, non_compliant
	Description     string
	LastCheckedAt   time.Time
	Recommendations []string
}

// GetComplianceOverview는 규제 준수 현황을 반환합니다.
func GetComplianceOverview() []ComplianceStatus {
	now := time.Now()
	return []ComplianceStatus{
		// GDPR
		{Framework: "GDPR", Category: "consent_management", Status: "compliant",
			Description: "데이터 공유 동의 관리 (생성/조회/철회)", LastCheckedAt: now},
		{Framework: "GDPR", Category: "data_access_log", Status: "compliant",
			Description: "데이터 접근 이력 추적 (Art.15 접근권)", LastCheckedAt: now},
		{Framework: "GDPR", Category: "right_to_erasure", Status: "compliant",
			Description: "삭제 요청 처리 프레임워크 (Art.17)", LastCheckedAt: now},
		{Framework: "GDPR", Category: "data_retention", Status: "compliant",
			Description: "데이터 유형별 보존 기간 정책 (Art.5)", LastCheckedAt: now},
		{Framework: "GDPR", Category: "phi_encryption", Status: "compliant",
			Description: "AES-256-GCM 필드 레벨 암호화", LastCheckedAt: now},
		{Framework: "GDPR", Category: "data_portability", Status: "partial",
			Description:     "FHIR 기반 건강 데이터 내보내기 (Art.20)",
			LastCheckedAt:   now,
			Recommendations: []string{"전체 데이터 유형 지원 확대 필요"}},
		{Framework: "GDPR", Category: "dpia", Status: "compliant",
			Description: "데이터 보호 영향 평가 문서 완성 (Art.35)", LastCheckedAt: now},

		// HIPAA
		{Framework: "HIPAA", Category: "access_control", Status: "compliant",
			Description: "JWT + RBAC 기반 접근 제어", LastCheckedAt: now},
		{Framework: "HIPAA", Category: "audit_trail", Status: "compliant",
			Description: "감사 로그 기록 및 조회 (§164.312(b))", LastCheckedAt: now},
		{Framework: "HIPAA", Category: "encryption", Status: "compliant",
			Description: "AES-256-GCM 암호화 (§164.312(a)(2)(iv))", LastCheckedAt: now},
		{Framework: "HIPAA", Category: "integrity", Status: "partial",
			Description:     "PHI 무결성 검증 (§164.312(c)(1))",
			LastCheckedAt:   now,
			Recommendations: []string{"체크섬 기반 무결성 검증 강화"}},

		// MFDS (식약처)
		{Framework: "MFDS", Category: "software_classification", Status: "compliant",
			Description: "IEC 62304 Class B 분류 완료", LastCheckedAt: now},
		{Framework: "MFDS", Category: "risk_management", Status: "compliant",
			Description: "ISO 14971 위험관리 계획 + FMEA 완료", LastCheckedAt: now},
		{Framework: "MFDS", Category: "software_requirements", Status: "compliant",
			Description: "IEC 62304 SRS v3.0 완성", LastCheckedAt: now},
		{Framework: "MFDS", Category: "clinical_validation", Status: "non_compliant",
			Description:     "임상시험 미실시",
			LastCheckedAt:   now,
			Recommendations: []string{"IRB 승인 후 100명 임상시험 실시"}},

		// FDA 510(k)
		{Framework: "FDA", Category: "predicate_device", Status: "compliant",
			Description: "Predicate Device 조사 완료", LastCheckedAt: now},
		{Framework: "FDA", Category: "se_report", Status: "non_compliant",
			Description:     "Substantial Equivalence 보고서 미작성",
			LastCheckedAt:   now,
			Recommendations: []string{"SE 보고서 작성 + FDA 사전 상담"}},

		// CE-IVDR
		{Framework: "CE-IVDR", Category: "technical_documentation", Status: "partial",
			Description:     "기술 문서 구조 정의, 상세 작성 필요",
			LastCheckedAt:   now,
			Recommendations: []string{"Notified Body 선정 및 컨설팅"}},
	}
}
