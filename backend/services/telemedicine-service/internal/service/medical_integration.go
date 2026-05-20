// Package service: 화상진료 ↔ medical 도메인 통합 어댑터
//
// telemedicine-service가 backend/shared/medical/ 5개 모듈을 유기적으로 호출하는
// 통합 지점입니다. 모세혈관 워크플로우의 핵심 보강.
//
//   Flow:
//     L1 디바이스 → L2 POCT 파서 → L3 FHIR 변환 → L4 안전망/추론 → L5 SOAP
//
// SSOT 무결성: 본 파일은 SW 영역만을 다루며, 디바이스 BoM 직접 참조 없음.
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/manpasik/backend/shared/medical/fhir"
	"github.com/manpasik/backend/shared/medical/observability"
	"github.com/manpasik/backend/shared/medical/poct"
	"github.com/manpasik/backend/shared/medical/reasoning"
	"github.com/manpasik/backend/shared/medical/safety"
	"github.com/manpasik/backend/shared/medical/scribe"
	"github.com/manpasik/backend/shared/medical/sla"
	"github.com/manpasik/backend/shared/ops/tracing"
)

// ============================================================================
// Medical Integration — 5계층 통합 코디네이터
// ============================================================================

// CalibrationProvider는 차분식 보정에 필요한 기준 신호와 표준오차를 제공합니다.
//
// Phase 5 SSOT (캘리브레이션 데이터)에서 주입됩니다.
// 본 SW 영역은 인터페이스만 정의하고 실제 캘리브레이션 값은 외부에서 공급받습니다.
type CalibrationProvider interface {
	// RefSignal은 측정 채널의 기준 신호 R_n을 반환합니다.
	RefSignal(loincCode string, rawSignal float64) float64
	// StandardError는 측정값의 95% CI 계산용 표준오차를 반환합니다.
	StandardError(loincCode string, value float64) float64
}

// fallbackCalibrationProvider는 운영 캘리브레이션이 미주입된 환경에서
// 안전한 기본값(보정 미적용)을 반환합니다.
//
// 운영 환경에서는 반드시 SetCalibrationProvider를 호출하여
// Phase 5 SSOT 기반 실 캘리브레이션을 주입해야 합니다.
type fallbackCalibrationProvider struct{}

func (fallbackCalibrationProvider) RefSignal(_ string, _ float64) float64 { return 0 }
func (fallbackCalibrationProvider) StandardError(_ string, _ float64) float64 {
	return 0
}

// MedicalIntegrator는 화상진료 세션에서 의료 도메인 모듈들을 조율합니다.
type MedicalIntegrator struct {
	parser       *poct.Parser
	diffEngine   *poct.DifferentialEngine
	fhirBuilder  *fhir.Builder
	safetyNet    *safety.PolarisSafetyNet
	auditLog     *safety.AuditLog
	reasoner     *reasoning.IntegratedReasoner
	scribe       *scribe.Scribe
	calibration  CalibrationProvider
	instrumentor *observability.MedicalInstrumentor
}

// NewMedicalIntegrator는 새 통합 코디네이터를 생성합니다.
//
// 기본 캘리브레이션 Provider는 보정 0(원시값 보존)을 반환합니다.
// 기본 계측은 인메모리 tracer + 표준 SLA 임계값.
// 운영에서는 SetCalibrationProvider, SetInstrumentor로 실 의존성을 주입하세요.
func NewMedicalIntegrator() *MedicalIntegrator {
	return &MedicalIntegrator{
		parser:      poct.NewParser(),
		diffEngine:  poct.NewDifferentialEngine(),
		fhirBuilder: fhir.NewBuilder(fhir.VersionR5),
		safetyNet:   safety.NewPolarisSafetyNet(),
		auditLog:    safety.NewAuditLog(10000),
		reasoner:    reasoning.NewIntegratedReasoner(),
		scribe:      scribe.NewScribe("ko"),
		calibration: fallbackCalibrationProvider{},
		instrumentor: observability.NewMedicalInstrumentor(
			tracing.NewMemoryTracer("telemedicine"),
			sla.NewMonitor(nil),
			"telemedicine",
		),
	}
}

// SetCalibrationProvider는 Phase 5 SSOT 기반 실 캘리브레이션을 주입합니다.
func (m *MedicalIntegrator) SetCalibrationProvider(p CalibrationProvider) {
	if p != nil {
		m.calibration = p
	}
}

// SetInstrumentor는 운영 환경의 분산 추적 + SLA 모니터를 주입합니다.
//
// 운영에서는 OTLP tracer + 공유 SLA Monitor를 주입하여 모든 telemedicine
// 인스턴스의 측정값을 단일 백엔드로 집계할 수 있습니다.
func (m *MedicalIntegrator) SetInstrumentor(inst *observability.MedicalInstrumentor) {
	if inst != nil {
		m.instrumentor = inst
	}
}

// InstrumentationReport는 단일 세션의 추적/SLA 보고서를 반환합니다.
func (m *MedicalIntegrator) InstrumentationReport(consultationID, patientID string) *observability.WorkflowReport {
	return m.instrumentor.GenerateReport(consultationID, patientID)
}

// ConsultationContext는 단일 화상진료 세션의 의료 도메인 컨텍스트입니다.
type ConsultationContext struct {
	ConsultationID string
	PatientID      string
	DoctorID       string
	StartedAt      time.Time
}

// ============================================================================
// Step 1: POCT 측정값 수신 + 차분식 보정
// ============================================================================

// ProcessMeasurement는 POCT1-A 메시지를 파싱·보정·FHIR 변환합니다.
//
// 분산 추적: medical.poct.parse + medical.fhir.convert 두 span 자동 생성.
// SLA: STAT 30초 등급 자동 측정.
func (m *MedicalIntegrator) ProcessMeasurement(
	ctx context.Context,
	consultation *ConsultationContext,
	rawPOCTMessage string,
	deviceRef string, // 예: "Device/dev-001" (메타데이터 reference만, 부품 정보 미포함)
) (*MeasurementResult, error) {
	if consultation == nil {
		return nil, errors.New("consultation context required")
	}

	op := m.instrumentor.StartOperation(ctx, "medical.poct.parse",
		sla.TierSTAT, consultation.ConsultationID, "")
	op.SetAttribute("patient_id", consultation.PatientID)

	// 1) POCT1-A 파싱
	msg, err := m.parser.Parse(rawPOCTMessage)
	if err != nil {
		op.End(err)
		return nil, fmt.Errorf("POCT 파싱 실패: %w", err)
	}
	op.SetAttribute("observation.count", fmt.Sprintf("%d", len(msg.Observations)))

	// 2) 차분식 보정 + 95% CI 적용
	//
	// 기준 신호 R_n과 표준오차 모두 CalibrationProvider(Phase 5 SSOT)에서 주입.
	// fallback Provider는 0을 반환하여 보정/CI를 비활성화하므로 운영 시
	// SetCalibrationProvider로 실 데이터를 반드시 주입해야 합니다.
	for _, obs := range msg.Observations {
		obs.RawSignal = obs.Value // 원시 신호 보존
		if obs.RefSignal == 0 {
			obs.RefSignal = m.calibration.RefSignal(obs.LOINCCode, obs.Value)
		}
		m.diffEngine.CorrectObservation(obs, obs.LOINCCode)
		stdErr := m.calibration.StandardError(obs.LOINCCode, obs.Value)
		if stdErr > 0 {
			poct.ApplyCI95(obs, stdErr)
		}
	}

	op.End(nil) // POCT 파싱 + 보정 단계 종료

	// 3) FHIR R5 리소스 생성 (별도 span)
	fhirOp := m.instrumentor.StartOperation(op.Context(), "medical.fhir.convert",
		sla.TierUrgent, consultation.ConsultationID, "")
	fhirOp.SetAttribute("patient_id", consultation.PatientID)

	patientRef := fmt.Sprintf("Patient/%s", consultation.PatientID)
	report, observations, err := m.fhirBuilder.FromPOCTMessage(msg, patientRef, deviceRef)
	if err != nil {
		fhirOp.End(err)
		return nil, fmt.Errorf("FHIR 변환 실패: %w", err)
	}
	fhirOp.End(nil)

	return &MeasurementResult{
		POCTMessage:  msg,
		Report:       report,
		Observations: observations,
	}, nil
}

// MeasurementResult는 측정 처리 결과입니다.
type MeasurementResult struct {
	POCTMessage  *poct.Message
	Report       *fhir.DiagnosticReport
	Observations []*fhir.Observation
}

// ============================================================================
// Step 2: 발화 안전망 검사 + 감사 기록
// ============================================================================

// CheckUtterance는 디지털 주치의/AI 발화를 안전망으로 검증하고 감사합니다.
//
// 분산 추적: medical.safety.check span 자동 생성. 위반 시 span 에러로 표시.
func (m *MedicalIntegrator) CheckUtterance(
	consultation *ConsultationContext,
	utterance string,
) *safety.CheckResult {
	if consultation == nil {
		return &safety.CheckResult{Safe: false}
	}

	op := m.instrumentor.StartOperation(context.Background(), "medical.safety.check",
		sla.TierUrgent, consultation.ConsultationID, "")
	op.SetAttribute("patient_id", consultation.PatientID)

	result := m.safetyNet.Check(utterance)
	m.auditLog.Record(consultation.PatientID, consultation.ConsultationID, result)

	if !result.Safe {
		op.SetAttribute("safety.violations", fmt.Sprintf("%d", len(result.Violations)))
		op.SetAttribute("safety.action", string(result.HighestAction))
		op.End(fmt.Errorf("safety violations: %d", len(result.Violations)))
	} else {
		op.End(nil)
	}
	return result
}

// SafetyAuditCount는 세션의 누적 감사 로그 건수를 반환합니다.
func (m *MedicalIntegrator) SafetyAuditCount(consultationID string) int {
	return len(m.auditLog.FindBySession(consultationID))
}

// ============================================================================
// Step 3: CoT + SHAP 통합 추론
// ============================================================================

// AnalyzeReasoning은 측정값과 진단 후보로 임상 추론을 수행합니다.
//
// 분산 추적: medical.cot.reason span 자동 생성. SLA STAT 30s 등급 측정.
func (m *MedicalIntegrator) AnalyzeReasoning(
	consultation *ConsultationContext,
	question string,
	observations map[string]float64, // LOINC code → value
	weights map[string]float64,
	baseline float64,
	candidates []reasoning.DiagnosisCandidate,
) (*reasoning.ReasoningResult, error) {
	if consultation == nil {
		return nil, errors.New("consultation context required")
	}

	op := m.instrumentor.StartOperation(context.Background(), "medical.cot.reason",
		sla.TierSTAT, consultation.ConsultationID, "")
	op.SetAttribute("patient_id", consultation.PatientID)
	op.SetAttribute("candidates.count", fmt.Sprintf("%d", len(candidates)))

	result, err := m.reasoner.Reason(
		consultation.PatientID,
		question,
		observations, weights, baseline,
		candidates,
	)
	op.End(err)
	return result, err
}

// ============================================================================
// Step 4: SOAP 노트 자동 생성
// ============================================================================

// GenerateSOAP는 트랜스크립트 + 측정값 + 추론을 통합해 SOAP 노트를 생성합니다.
func (m *MedicalIntegrator) GenerateSOAP(
	consultation *ConsultationContext,
	transcript *scribe.Transcript,
	measurement *MeasurementResult,
	reasoning *reasoning.ReasoningResult,
	primaryDiagnosis, icd10 string,
	differentials []string,
	riskLevel string,
) (*scribe.SOAPNote, *fhir.DocumentReference, error) {
	if consultation == nil {
		return nil, nil, errors.New("consultation context required")
	}

	// FHIR Observation → scribe.LabResult 변환
	var labResults []*scribe.LabResult
	if measurement != nil {
		labResults = make([]*scribe.LabResult, 0, len(measurement.Observations))
		for _, obs := range measurement.Observations {
			lr := &scribe.LabResult{
				TestCode: obs.Code.Coding[0].Code,
				TestName: obs.Code.Coding[0].Display,
				Value:    obs.ValueQuantity.Value,
				Unit:     obs.ValueQuantity.Unit,
				Alpha:    obs.AlphaCorrection,
				CILow:    obs.CILow,
				CIHigh:   obs.CIHigh,
			}
			if len(obs.Interpretation) > 0 && len(obs.Interpretation[0].Coding) > 0 {
				lr.Flag = obs.Interpretation[0].Coding[0].Code
			}
			if len(obs.ReferenceRange) > 0 {
				if obs.ReferenceRange[0].Low != nil {
					lr.RefLow = obs.ReferenceRange[0].Low.Value
				}
				if obs.ReferenceRange[0].High != nil {
					lr.RefHigh = obs.ReferenceRange[0].High.Value
				}
			}
			labResults = append(labResults, lr)
		}
	}

	reasoningSummary := ""
	if reasoning != nil && reasoning.Chain != nil {
		reasoningSummary = reasoning.Chain.Conclusion
	}

	op := m.instrumentor.StartOperation(context.Background(), "medical.scribe.generate",
		sla.TierRoutine, consultation.ConsultationID, "")
	op.SetAttribute("patient_id", consultation.PatientID)
	defer func() { op.End(nil) }()

	soap, err := m.scribe.Generate(transcript, &scribe.GenerateOptions{
		PatientID:     consultation.PatientID,
		EncounterID:   consultation.ConsultationID,
		Provider:      consultation.DoctorID,
		LabResults:    labResults,
		Diagnosis:     primaryDiagnosis,
		ICD10:         icd10,
		Differentials: differentials,
		Reasoning:     reasoningSummary,
		RiskLevel:     riskLevel,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("SOAP 생성 실패: %w", err)
	}

	docRef, err := scribe.ToFHIRDocumentReference(soap)
	if err != nil {
		return soap, nil, fmt.Errorf("DocumentReference 변환 실패: %w", err)
	}

	return soap, docRef, nil
}

// ============================================================================
// Step 5: Bundle 묶음 (전송용)
// ============================================================================

// BuildTransactionBundle은 화상진료 1세션의 모든 FHIR 리소스를 하나의 트랜잭션으로 묶습니다.
//
// 외부 EHR 푸시 시 단일 트랜잭션으로 전송 가능.
func (m *MedicalIntegrator) BuildTransactionBundle(
	measurement *MeasurementResult,
	docRef *fhir.DocumentReference,
) *fhir.Bundle {
	resources := make([]fhir.Resource, 0)
	if measurement != nil {
		if measurement.Report != nil {
			resources = append(resources, measurement.Report)
		}
		for _, obs := range measurement.Observations {
			resources = append(resources, obs)
		}
	}
	if docRef != nil {
		resources = append(resources, docRef)
	}
	return fhir.NewBundle("transaction", resources...)
}

// ============================================================================
// 통합 점검 헬퍼
// ============================================================================

// ConsultationSummary는 통합 흐름 요약 보고서입니다.
type ConsultationSummary struct {
	ConsultationID    string
	PatientID         string
	MeasurementCount  int
	AbnormalCount     int
	SafetyViolations  int
	DiagnosisHypothesis string
	OverallConfidence float64
	FHIRBundleID      string
	GeneratedAt       time.Time
}

// Summarize는 화상진료 세션 흐름의 결과를 요약합니다.
func (m *MedicalIntegrator) Summarize(
	consultation *ConsultationContext,
	measurement *MeasurementResult,
	reasoning *reasoning.ReasoningResult,
	bundle *fhir.Bundle,
) *ConsultationSummary {
	if consultation == nil {
		return nil
	}

	summary := &ConsultationSummary{
		ConsultationID: consultation.ConsultationID,
		PatientID:      consultation.PatientID,
		GeneratedAt:    time.Now().UTC(),
	}

	if measurement != nil {
		summary.MeasurementCount = len(measurement.Observations)
		for _, obs := range measurement.POCTMessage.Observations {
			if obs.IsAbnormal() {
				summary.AbnormalCount++
			}
		}
	}

	summary.SafetyViolations = m.SafetyAuditCount(consultation.ConsultationID)

	if reasoning != nil && reasoning.Chain != nil {
		// 첫 번째 가설 단계 추출
		hypotheses := reasoning.Chain.Steps
		for _, s := range hypotheses {
			if s.Type == "hypothesis" {
				summary.DiagnosisHypothesis = s.Title
				summary.OverallConfidence = s.Confidence
				break
			}
		}
		if summary.DiagnosisHypothesis == "" {
			summary.OverallConfidence = reasoning.Chain.OverallConfidence
		}
	}

	if bundle != nil {
		summary.FHIRBundleID = bundle.ID
	}
	return summary
}
