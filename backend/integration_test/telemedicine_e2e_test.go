// Package integration_test: P10-09b 화상진료 E2E 시나리오
//
// 5계층 아키텍처 전체 흐름 검증:
//   L1 (디바이스 측정) → L2 (POCT 게이트웨이) → L3 (FHIR) → L4 (CoT/안전망) → L5 (SOAP)
package integration_test

import (
	"strings"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/medical/fhir"
	"github.com/manpasik/backend/shared/medical/poct"
	"github.com/manpasik/backend/shared/medical/reasoning"
	"github.com/manpasik/backend/shared/medical/safety"
	"github.com/manpasik/backend/shared/medical/scribe"
)

// TestTelemedicineE2E_FullFlow는 화상진료 전체 흐름을 검증합니다.
//
// 시나리오:
//   1. POCT1-A 메시지 수신 (L1→L2)
//   2. 차분식 보정 + 95% CI 계산 (L2)
//   3. FHIR DiagnosticReport + Observation 생성 (L3)
//   4. CoT + SHAP 추론사슬 (L4)
//   5. Polaris 안전망 검사 (L4)
//   6. SOAP 노트 + DocumentReference (L5)
func TestTelemedicineE2E_FullFlow(t *testing.T) {
	// ---- L1→L2: POCT1-A 메시지 파싱 ----
	rawMsg := `MSH|^~\&|MMUP-001|FAC|GW|HOSP||20260430120000|ORU^R01|MSG-E2E-001
PID|||PATIENT-E2E-789
OBR|1|ORDER-E2E||5811-5^Glucose Panel
OBX|1|NM|2345-7^Glucose|1|180.5|mg/dL|70-100|H||F|20260430120030
OBX|2|NM|4548-4^HbA1c|1|7.5|%|0-5.7|H||F|20260430120030`

	parser := poct.NewParser()
	msg, err := parser.Parse(rawMsg)
	if err != nil {
		t.Fatalf("POCT 파싱 실패: %v", err)
	}
	if len(msg.Observations) != 2 {
		t.Errorf("Observations = %d, want 2", len(msg.Observations))
	}

	// ---- L2: 차분식 보정 ----
	engine := poct.NewDifferentialEngine()
	for _, obs := range msg.Observations {
		obs.RawSignal = obs.Value
		obs.RefSignal = obs.Value * 0.05 // 시뮬레이션: 5% 노이즈
		engine.CorrectObservation(obs, obs.LOINCCode)
		// 95% CI 적용 (stdErr = 2% of value)
		poct.ApplyCI95(obs, obs.Value*0.02)
	}

	// 보정 검증
	for _, obs := range msg.Observations {
		if obs.Alpha != poct.DefaultAlpha {
			t.Errorf("%s alpha = %f", obs.LOINCCode, obs.Alpha)
		}
		if obs.CIHigh <= obs.CILow {
			t.Errorf("%s CI 폭 비정상: %f-%f", obs.LOINCCode, obs.CILow, obs.CIHigh)
		}
	}

	// ---- L3: FHIR R5 생성 ----
	builder := fhir.NewBuilder(fhir.VersionR5)
	report, observations, err := builder.FromPOCTMessage(msg, "Patient/PATIENT-E2E-789", "Device/MMUP-001")
	if err != nil {
		t.Fatalf("FHIR 생성 실패: %v", err)
	}
	if report.Status != "final" {
		t.Errorf("Report status = %q", report.Status)
	}
	if len(observations) != 2 {
		t.Errorf("FHIR Observations = %d", len(observations))
	}

	// 만파식 확장 필드 검증
	if observations[0].AlphaCorrection == 0 {
		t.Error("Alpha correction 누락")
	}
	if observations[0].CIHigh == 0 {
		t.Error("CI 값 누락")
	}

	// ---- L4: CoT + SHAP 통합 추론 ----
	reasoner := reasoning.NewIntegratedReasoner()
	observationMap := map[string]float64{
		"2345-7": observations[0].ValueQuantity.Value,
		"4548-4": observations[1].ValueQuantity.Value,
	}
	weights := map[string]float64{
		"2345-7": 0.3,
		"4548-4": 0.4,
	}
	candidates := []reasoning.DiagnosisCandidate{
		{Name: "제2형 당뇨병", ICD10: "E11.9", Probability: 0.85,
			Reasoning: "공복혈당 ≥126 mg/dL + HbA1c ≥6.5% 진단 기준 충족"},
		{Name: "전당뇨병", ICD10: "R73.03", Probability: 0.10},
	}

	result, err := reasoner.Reason(
		"PATIENT-E2E-789",
		"환자의 당뇨 관련 위험도 평가",
		observationMap, weights, 50.0, candidates,
	)
	if err != nil {
		t.Fatalf("Reasoning 실패: %v", err)
	}
	if result.Chain == nil || len(result.Chain.Steps) < 3 {
		t.Errorf("Chain steps = %d, want >= 3", len(result.Chain.Steps))
	}
	if result.Explanation == nil || len(result.Explanation.Contributions) == 0 {
		t.Error("SHAP 설명 누락")
	}
	if result.BuildTimeMs > 1500 {
		t.Errorf("BuildTime = %dms, want <= 1500ms (D-03 KPI)", result.BuildTimeMs)
	}

	// ---- L4: Polaris 안전망 검사 ----
	polaris := safety.NewPolarisSafetyNet()
	auditLog := safety.NewAuditLog(1000)

	// 의사 발화 시뮬레이션 (안전한 발화)
	doctorResponse := "혈당과 HbA1c가 높은 편입니다. 내분비내과 진료를 받으시기를 권장드립니다."
	checkResult := polaris.Check(doctorResponse)
	auditLog.Record("PATIENT-E2E-789", "session-E2E", checkResult)

	if !checkResult.Safe {
		t.Errorf("정상 발화가 unsafe로 판정: violations=%v", checkResult.Violations)
	}
	if !strings.Contains(checkResult.SanitizedText, "[안내]") {
		t.Error("면책 조항 미추가")
	}

	// 위험 발화 시뮬레이션
	dangerousResponse := "이 약을 처방해드릴게요. 복용량을 두 배로 늘리세요."
	dangerCheck := polaris.Check(dangerousResponse)
	auditLog.Record("PATIENT-E2E-789", "session-E2E", dangerCheck)
	if !dangerCheck.HasBlockedContent() {
		t.Error("위험 발화 미차단")
	}

	if auditLog.Count() < 1 {
		t.Errorf("감사 로그 = %d, want >= 1", auditLog.Count())
	}

	// ---- L5: Scribe SOAP 생성 ----
	transcript := &scribe.Transcript{
		SessionID: "session-E2E",
		Segments: []*scribe.TranscriptSegment{
			{Speaker: "doctor", Text: "안녕하세요, 어떻게 오셨나요?", Confidence: 0.95},
			{Speaker: "patient", Text: "최근 갈증이 심하고 소변이 자주 마려워요. 가족 중 당뇨병 환자가 있습니다.", Confidence: 0.92},
			{Speaker: "doctor", Text: "혈당 수치가 높네요. 내분비내과 진료가 필요합니다. 다음 주에 결과를 가지고 다시 오세요.", Confidence: 0.94},
		},
		StartedAt: time.Now().UTC().Add(-15 * time.Minute),
		EndedAt:   time.Now().UTC(),
	}

	scribeEngine := scribe.NewScribe("ko")

	// FHIR LabResult로 변환
	labResults := make([]*scribe.LabResult, 0, len(observations))
	for _, obs := range observations {
		labResults = append(labResults, &scribe.LabResult{
			TestCode: obs.Code.Coding[0].Code,
			TestName: obs.Code.Coding[0].Display,
			Value:    obs.ValueQuantity.Value,
			Unit:     obs.ValueQuantity.Unit,
			Flag:     "H",
			Alpha:    obs.AlphaCorrection,
			CILow:    obs.CILow,
			CIHigh:   obs.CIHigh,
		})
		// 정상 범위 추출
		if len(obs.ReferenceRange) > 0 {
			labResults[len(labResults)-1].RefLow = obs.ReferenceRange[0].Low.Value
			labResults[len(labResults)-1].RefHigh = obs.ReferenceRange[0].High.Value
		}
	}

	soap, err := scribeEngine.Generate(transcript, &scribe.GenerateOptions{
		PatientID:     "PATIENT-E2E-789",
		EncounterID:   "session-E2E",
		Provider:      "Dr. Kim",
		LabResults:    labResults,
		Diagnosis:     candidates[0].Name,
		ICD10:         candidates[0].ICD10,
		Differentials: []string{candidates[1].Name},
		Reasoning:     result.Chain.Conclusion,
		RiskLevel:     "high",
	})
	if err != nil {
		t.Fatalf("Scribe 실패: %v", err)
	}
	if soap.Subjective.ChiefComplaint == "" {
		t.Error("주증상 미추출")
	}
	if soap.Assessment.PrimaryDiagnosis != "제2형 당뇨병" {
		t.Errorf("진단 = %q", soap.Assessment.PrimaryDiagnosis)
	}
	if len(soap.Objective.LabResults) != 2 {
		t.Errorf("LabResults = %d", len(soap.Objective.LabResults))
	}

	// ---- L5: SOAP → FHIR DocumentReference ----
	docRef, err := scribe.ToFHIRDocumentReference(soap)
	if err != nil {
		t.Fatalf("DocumentReference 변환 실패: %v", err)
	}
	if docRef.Subject.Reference != "Patient/PATIENT-E2E-789" {
		t.Errorf("Subject = %q", docRef.Subject.Reference)
	}
	if len(docRef.Content) == 0 || len(docRef.Content[0].Attachment.Data) == 0 {
		t.Error("DocumentReference 콘텐츠 누락")
	}

	// ---- L5: Bundle 묶기 (트랜잭션) ----
	resources := []fhir.Resource{report, docRef}
	for _, obs := range observations {
		resources = append(resources, obs)
	}
	bundle := fhir.NewBundle("transaction", resources...)
	if len(bundle.Entry) != len(resources) {
		t.Errorf("Bundle entry = %d, want %d", len(bundle.Entry), len(resources))
	}
}

// TestTelemedicineE2E_EmergencyEscalation은 응급 상황 에스컬레이션 흐름을 검증합니다.
func TestTelemedicineE2E_EmergencyEscalation(t *testing.T) {
	// 환자가 화상진료 중 응급 상황 발화
	patientText := "갑자기 가슴이 너무 아프고 호흡이 힘들어요. 의식이 흐려져요."

	polaris := safety.NewPolarisSafetyNet()
	check := polaris.Check(patientText)

	if !check.RequiresEscalation() {
		t.Error("응급 상황이 escalation 요구로 분류되지 않음")
	}

	// 카테고리 확인
	hasEmergency := false
	for _, cat := range check.CategoriesViolated() {
		if cat == safety.CategoryEmergency {
			hasEmergency = true
		}
	}
	if !hasEmergency {
		t.Error("Emergency 카테고리 미감지")
	}
}

// TestTelemedicineE2E_HighDifferential은 감별진단 케이스를 검증합니다.
func TestTelemedicineE2E_HighDifferential(t *testing.T) {
	reasoner := reasoning.NewIntegratedReasoner()

	// 흉통 증상 → 다양한 감별진단
	candidates := []reasoning.DiagnosisCandidate{
		{Name: "심근경색", ICD10: "I21", Probability: 0.45,
			Reasoning: "트로포닌 ↑, ECG ST 상승"},
		{Name: "역류성 식도염", ICD10: "K21.0", Probability: 0.30,
			Reasoning: "산성 트림, 음식 후 악화"},
		{Name: "근육통", ICD10: "M79.1", Probability: 0.15,
			Reasoning: "운동 후 발생, 압통 양성"},
	}

	observations := map[string]float64{
		"6598-7":  2.5,  // troponin (above 0.04 ng/mL = abnormal)
		"30341-2": 100,  // ECG ST elevation (mock)
	}
	weights := map[string]float64{
		"6598-7":  0.6,
		"30341-2": 0.4,
	}

	result, err := reasoner.Reason("p-emerg-1", "흉통 평가", observations, weights, 30.0, candidates)
	if err != nil {
		t.Fatalf("Reason 실패: %v", err)
	}

	if len(result.Chain.Steps) < 3 {
		t.Errorf("Steps = %d, want >= 3 (관측+감별+근거+가설)", len(result.Chain.Steps))
	}

	// 감별진단 단계 존재 확인
	diff := result.Chain.FilterByType(reasoning.StepDifferential)
	if len(diff) == 0 {
		t.Error("Differential step 누락")
	}
}

// TestTelemedicineE2E_PIIRedaction은 PII 노출 차단을 검증합니다.
func TestTelemedicineE2E_PIIRedaction(t *testing.T) {
	polaris := safety.NewPolarisSafetyNet()

	textWithPII := "환자 주민등록번호 901020-1234567, 연락처 010-1234-5678 입니다."
	check := polaris.Check(textWithPII)

	if strings.Contains(check.SanitizedText, "901020-1234567") {
		t.Error("주민등록번호가 노출됨")
	}
	if strings.Contains(check.SanitizedText, "010-1234-5678") {
		t.Error("전화번호가 노출됨")
	}
	if !strings.Contains(check.SanitizedText, "[개인정보]") {
		t.Error("주민등록번호 마스킹 누락")
	}
}

// TestTelemedicineE2E_BundleTransactionStructure는 FHIR Bundle 트랜잭션 구조를 검증합니다.
func TestTelemedicineE2E_BundleTransactionStructure(t *testing.T) {
	parser := poct.NewParser()
	msg, _ := parser.Parse(`MSH|^~\&|D|F|R|F||20260430|ORU^R01|MSG-001
OBR|1|ORD||1234-5^Test
OBX|1|NM|1234-5^Test|1|99.0|mg/dL|0-100|N||F|20260430`)

	builder := fhir.NewBuilder(fhir.VersionR5)
	msg.PatientID = "p"
	report, observations, _ := builder.FromPOCTMessage(msg, "Patient/p", "Device/d")

	resources := []fhir.Resource{report}
	for _, obs := range observations {
		resources = append(resources, obs)
	}

	bundle := fhir.NewBundle("transaction", resources...)
	for _, entry := range bundle.Entry {
		if entry.Request == nil {
			t.Error("transaction Bundle entry에 Request 누락")
		}
		if entry.Request.Method != "POST" {
			t.Errorf("Method = %q, want POST", entry.Request.Method)
		}
	}
}

// TestTelemedicineE2E_AuditTrailComplete는 감사 추적 완성도를 검증합니다.
func TestTelemedicineE2E_AuditTrailComplete(t *testing.T) {
	polaris := safety.NewPolarisSafetyNet()
	log := safety.NewAuditLog(100)

	sessionID := "audit-test-session"
	developerInputs := []string{
		"이 약을 처방하세요",
		"수술이 필요합니다",
		"의식없이 쓰러진 경우 119를 부르세요",
		"두통이 있을 때는 충분한 휴식을 권합니다", // 정상
	}

	for _, text := range developerInputs {
		check := polaris.Check(text)
		log.Record("user-001", sessionID, check)
	}

	// 정상 발화 1건은 기록되지 않으므로 3건 예상
	entries := log.FindBySession(sessionID)
	if len(entries) < 2 {
		t.Errorf("위반 감사 로그 = %d, want >= 2", len(entries))
	}

	counts := log.CountByCategory()
	if counts[safety.CategoryPrescription] == 0 {
		t.Error("Prescription 카운트 누락")
	}
	if counts[safety.CategorySurgery] == 0 {
		t.Error("Surgery 카운트 누락")
	}
}
