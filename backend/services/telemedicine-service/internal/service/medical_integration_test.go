package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/manpasik/backend/services/telemedicine-service/internal/service"
	"github.com/manpasik/backend/shared/medical/reasoning"
	"github.com/manpasik/backend/shared/medical/scribe"
)

const samplePOCT = `MSH|^~\&|D|F|R|F||20260430120000|ORU^R01|MSG-INT-1
PID|||PATIENT-INT-1
OBR|1|ORD-INT||5811-5^Glucose Panel
OBX|1|NM|2345-7^Glucose|1|180.5|mg/dL|70-100|H||F|20260430120030`

func newTestConsultation() *service.ConsultationContext {
	return &service.ConsultationContext{
		ConsultationID: "consult-int-001",
		PatientID:      "PATIENT-INT-1",
		DoctorID:       "doctor-001",
		StartedAt:      time.Now().UTC(),
	}
}

// testCalibrationProvider는 테스트용 캘리브레이션 stub.
type testCalibrationProvider struct{}

func (testCalibrationProvider) RefSignal(_ string, value float64) float64 {
	return value * 0.05
}
func (testCalibrationProvider) StandardError(_ string, value float64) float64 {
	return value * 0.02
}

// newTestIntegrator는 캘리브레이션이 주입된 통합기를 반환합니다.
func newTestIntegrator() *service.MedicalIntegrator {
	mi := service.NewMedicalIntegrator()
	mi.SetCalibrationProvider(testCalibrationProvider{})
	return mi
}

// TestMedicalIntegrator_Capillary_FullFlow는 모세혈관 단위 5계층 통합 흐름을 검증합니다.
func TestMedicalIntegrator_Capillary_FullFlow(t *testing.T) {
	mi := newTestIntegrator()
	consultation := newTestConsultation()
	ctx := context.Background()

	// L1→L2: POCT 측정 처리
	measurement, err := mi.ProcessMeasurement(ctx, consultation, samplePOCT, "Device/dev-1")
	if err != nil {
		t.Fatalf("ProcessMeasurement 실패: %v", err)
	}
	if len(measurement.Observations) != 1 {
		t.Errorf("Observations = %d, want 1", len(measurement.Observations))
	}
	if measurement.Observations[0].AlphaCorrection == 0 {
		t.Error("차분식 alpha 미적용")
	}
	if measurement.Observations[0].CIHigh == 0 {
		t.Error("95% CI 미계산")
	}

	// L4: 발화 안전 검사
	check := mi.CheckUtterance(consultation, "혈당이 다소 높습니다. 내분비내과 상담을 권장드립니다.")
	if !check.Safe {
		t.Errorf("정상 발화가 unsafe로 분류: %v", check.Violations)
	}

	dangerCheck := mi.CheckUtterance(consultation, "이 약을 처방해드릴게요")
	if !dangerCheck.HasBlockedContent() {
		t.Error("위험 발화 미차단")
	}
	if mi.SafetyAuditCount("consult-int-001") < 1 {
		t.Error("감사 로그 미기록")
	}

	// L4: CoT 추론
	candidates := []reasoning.DiagnosisCandidate{
		{Name: "고혈당", ICD10: "R73.9", Probability: 0.7,
			Reasoning: "공복혈당 ≥126 mg/dL"},
	}
	reason, err := mi.AnalyzeReasoning(consultation,
		"혈당 이상 평가",
		map[string]float64{"2345-7": 180.5},
		map[string]float64{"2345-7": 0.5},
		50.0,
		candidates,
	)
	if err != nil {
		t.Fatalf("AnalyzeReasoning 실패: %v", err)
	}
	if reason.Chain == nil || len(reason.Chain.Steps) == 0 {
		t.Error("추론 단계 미생성")
	}

	// L5: SOAP 생성
	transcript := &scribe.Transcript{
		SessionID: consultation.ConsultationID,
		Segments: []*scribe.TranscriptSegment{
			{Speaker: "doctor", Text: "어떻게 오셨나요?", Confidence: 0.95},
			{Speaker: "patient", Text: "최근 갈증이 심하고 머리가 아파요", Confidence: 0.92},
			{Speaker: "doctor", Text: "내분비내과 진료 권장합니다. 다음 주에 다시 오세요.", Confidence: 0.94},
		},
		StartedAt: time.Now().UTC().Add(-15 * time.Minute),
		EndedAt:   time.Now().UTC(),
	}

	soap, docRef, err := mi.GenerateSOAP(consultation, transcript, measurement, reason,
		"고혈당", "R73.9", []string{"전당뇨"}, "moderate")
	if err != nil {
		t.Fatalf("GenerateSOAP 실패: %v", err)
	}
	if soap.Assessment.PrimaryDiagnosis != "고혈당" {
		t.Errorf("진단 = %q", soap.Assessment.PrimaryDiagnosis)
	}
	if docRef == nil || len(docRef.Content) == 0 {
		t.Error("DocumentReference 누락")
	}

	// L5: Bundle 묶음
	bundle := mi.BuildTransactionBundle(measurement, docRef)
	if bundle == nil || len(bundle.Entry) == 0 {
		t.Fatal("Bundle 생성 실패")
	}
	if bundle.Type != "transaction" {
		t.Errorf("Bundle Type = %q", bundle.Type)
	}

	// 통합 요약
	summary := mi.Summarize(consultation, measurement, reason, bundle)
	if summary == nil {
		t.Fatal("Summary 생성 실패")
	}
	if summary.MeasurementCount != 1 {
		t.Errorf("MeasurementCount = %d", summary.MeasurementCount)
	}
	if summary.AbnormalCount != 1 {
		t.Errorf("AbnormalCount = %d, want 1 (혈당 H)", summary.AbnormalCount)
	}
	if summary.SafetyViolations < 1 {
		t.Errorf("SafetyViolations = %d, want >= 1", summary.SafetyViolations)
	}
	if summary.FHIRBundleID == "" {
		t.Error("FHIR Bundle ID 누락")
	}
}

// TestMedicalIntegrator_RequiresContext는 ConsultationContext nil 거부를 검증합니다.
func TestMedicalIntegrator_RequiresContext(t *testing.T) {
	mi := newTestIntegrator()
	ctx := context.Background()

	_, err := mi.ProcessMeasurement(ctx, nil, samplePOCT, "Device/x")
	if err == nil {
		t.Error("nil context로 측정 처리 통과")
	}

	if check := mi.CheckUtterance(nil, "x"); check.Safe {
		t.Error("nil context로 발화 검사 통과")
	}

	_, err = mi.AnalyzeReasoning(nil, "q",
		map[string]float64{"x": 1}, map[string]float64{"x": 1}, 0, nil)
	if err == nil {
		t.Error("nil context로 추론 통과")
	}
}

// TestMedicalIntegrator_NoBoMReference는 medical 통합에 디바이스 BoM 직접 참조가 없는지 검증합니다.
//
// SSOT 무결성: 본 통합 어댑터는 SW 영역 한정이며, 부품명·핀 수·물리 사양 직접 참조 금지.
func TestMedicalIntegrator_NoBoMReference(t *testing.T) {
	mi := newTestIntegrator()
	consultation := newTestConsultation()
	ctx := context.Background()

	// 디바이스 reference는 메타데이터 식별자만 (부품 정보 없음)
	measurement, err := mi.ProcessMeasurement(ctx, consultation, samplePOCT, "Device/manpasik-001")
	if err != nil {
		t.Fatalf("ProcessMeasurement 실패: %v", err)
	}

	// FHIR Observation의 Device reference는 ID만 포함, 부품 메타데이터 없음
	if measurement.Observations[0].Device != nil {
		ref := measurement.Observations[0].Device.Reference
		if ref != "Device/manpasik-001" {
			t.Errorf("Device ref = %q, BoM 정보 노출 의심", ref)
		}
	}
}

// TestMedicalIntegrator_BundleMatchesObservations는 Bundle 구조 일관성을 검증합니다.
func TestMedicalIntegrator_BundleMatchesObservations(t *testing.T) {
	mi := newTestIntegrator()
	consultation := newTestConsultation()

	measurement, _ := mi.ProcessMeasurement(context.Background(), consultation, samplePOCT, "Device/d")
	transcript := &scribe.Transcript{
		Segments: []*scribe.TranscriptSegment{
			{Speaker: "patient", Text: "test", Confidence: 0.9},
		},
	}
	soap, docRef, _ := mi.GenerateSOAP(consultation, transcript, measurement, nil, "X", "X00", nil, "low")
	if soap == nil {
		t.Fatal("SOAP 생성 실패")
	}

	bundle := mi.BuildTransactionBundle(measurement, docRef)
	// Report 1 + Observations N + DocRef 1 = N+2
	expected := 1 + len(measurement.Observations) + 1
	if len(bundle.Entry) != expected {
		t.Errorf("Bundle entries = %d, want %d", len(bundle.Entry), expected)
	}
}

// TestMedicalIntegrator_Summary_NoMeasurement는 측정값 없이 요약 가능한지 확인합니다.
func TestMedicalIntegrator_Summary_NoMeasurement(t *testing.T) {
	mi := newTestIntegrator()
	consultation := newTestConsultation()

	summary := mi.Summarize(consultation, nil, nil, nil)
	if summary == nil {
		t.Fatal("Summary nil")
	}
	if summary.MeasurementCount != 0 {
		t.Errorf("MeasurementCount = %d, want 0", summary.MeasurementCount)
	}
}

// TestMedicalIntegrator_Reasoning_HypothesisExtraction은 추론 단계에서 가설 추출을 검증합니다.
func TestMedicalIntegrator_Reasoning_HypothesisExtraction(t *testing.T) {
	mi := newTestIntegrator()
	consultation := newTestConsultation()

	candidates := []reasoning.DiagnosisCandidate{
		{Name: "당뇨병", ICD10: "E11", Probability: 0.85, Reasoning: "공복혈당 + HbA1c 모두 진단 기준 초과"},
	}
	result, _ := mi.AnalyzeReasoning(consultation, "Q",
		map[string]float64{"2345-7": 200},
		map[string]float64{"2345-7": 0.5},
		50.0, candidates)

	measurement, _ := mi.ProcessMeasurement(context.Background(), consultation, samplePOCT, "Device/d")
	bundle := mi.BuildTransactionBundle(measurement, nil)
	summary := mi.Summarize(consultation, measurement, result, bundle)

	if summary.OverallConfidence < 0.5 {
		t.Errorf("OverallConfidence = %f, want >= 0.5", summary.OverallConfidence)
	}
}

// TestMedicalIntegrator_InstrumentationReportCaptured는 observability 주입 후
// 워크플로우 추적이 자동 수집되는지 검증합니다.
func TestMedicalIntegrator_InstrumentationReportCaptured(t *testing.T) {
	mi := newTestIntegrator()
	consultation := newTestConsultation()
	ctx := context.Background()

	// L1→L3 측정 처리 (POCT + FHIR 두 span 생성)
	_, err := mi.ProcessMeasurement(ctx, consultation, samplePOCT, "Device/d-1")
	if err != nil {
		t.Fatalf("ProcessMeasurement 실패: %v", err)
	}

	// 안전 검사 (1 span)
	mi.CheckUtterance(consultation, "혈당이 높습니다. 의사 상담을 권장합니다.")

	// 추론 (1 span)
	candidates := []reasoning.DiagnosisCandidate{
		{Name: "고혈당", ICD10: "R73.9", Probability: 0.7, Reasoning: "공복혈당"},
	}
	_, err = mi.AnalyzeReasoning(consultation, "Q",
		map[string]float64{"2345-7": 180},
		map[string]float64{"2345-7": 0.5},
		50.0, candidates)
	if err != nil {
		t.Fatalf("AnalyzeReasoning 실패: %v", err)
	}

	// 보고서: 최소 4 span (poct.parse + fhir.convert + safety.check + cot.reason)
	report := mi.InstrumentationReport(consultation.ConsultationID, consultation.PatientID)
	if report == nil {
		t.Fatal("Report nil")
	}
	if report.TotalSpans < 4 {
		t.Errorf("TotalSpans = %d, want >= 4 (관측 누락)", report.TotalSpans)
	}
	if report.ErrorSpans != 0 {
		t.Errorf("ErrorSpans = %d, want 0 (정상 흐름)", report.ErrorSpans)
	}
}

// TestMedicalIntegrator_UnsafeUtteranceMarksErrorSpan는 안전 위반 시
// span이 error로 마킹되는지 검증합니다.
func TestMedicalIntegrator_UnsafeUtteranceMarksErrorSpan(t *testing.T) {
	mi := newTestIntegrator()
	consultation := newTestConsultation()

	mi.CheckUtterance(consultation, "이 약을 처방하세요") // 위반

	report := mi.InstrumentationReport(consultation.ConsultationID, consultation.PatientID)
	if report.ErrorSpans < 1 {
		t.Errorf("ErrorSpans = %d, want >= 1 (안전 위반은 error span)", report.ErrorSpans)
	}
}

// TestMedicalIntegrator_SafetyViolationCount는 누적 감사 카운트를 검증합니다.
func TestMedicalIntegrator_SafetyViolationCount(t *testing.T) {
	mi := newTestIntegrator()
	consultation := newTestConsultation()

	for _, text := range []string{
		"이 약을 처방하세요",
		"수술이 필요합니다",
		"의식이 없으면 119를 부르세요",
	} {
		mi.CheckUtterance(consultation, text)
	}

	count := mi.SafetyAuditCount(consultation.ConsultationID)
	if count != 3 {
		t.Errorf("SafetyAuditCount = %d, want 3", count)
	}
}
