package scribe_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/medical/fhir"
	"github.com/manpasik/backend/shared/medical/scribe"
	"github.com/manpasik/backend/shared/medical/stt"
)

func newTestTranscript() *scribe.Transcript {
	return &scribe.Transcript{
		SessionID: "session-001",
		Segments: []*scribe.TranscriptSegment{
			{Speaker: "doctor", Text: "안녕하세요, 어디가 불편하신가요?", Confidence: 0.95},
			{Speaker: "patient", Text: "어제부터 머리가 심하게 아파요. 혈압약을 복용하고 있습니다.", Confidence: 0.92},
			{Speaker: "doctor", Text: "혈압이 좀 높아요. 진통제를 처방해드릴게요. 다음 주에 다시 오세요.", Confidence: 0.94},
			{Speaker: "patient", Text: "알레르기는 없습니다.", Confidence: 0.93},
		},
		StartedAt: time.Now().UTC().Add(-15 * time.Minute),
		EndedAt:   time.Now().UTC(),
	}
}

func TestScribe_GenerateBasicSOAP(t *testing.T) {
	s := scribe.NewScribe("ko")
	tr := newTestTranscript()

	note, err := s.Generate(tr, &scribe.GenerateOptions{
		PatientID:     "patient-001",
		EncounterID:   "enc-001",
		Provider:      "Dr. Kim",
		Diagnosis:     "고혈압성 두통",
		ICD10:         "I10",
		Differentials: []string{"긴장성 두통", "편두통"},
		RiskLevel:     "moderate",
	})
	if err != nil {
		t.Fatalf("Generate 실패: %v", err)
	}

	if note.PatientID != "patient-001" {
		t.Errorf("PatientID = %q", note.PatientID)
	}
	if note.Subjective == nil || note.Objective == nil || note.Assessment == nil || note.Plan == nil {
		t.Error("SOAP 4섹션 중 일부 누락")
	}
	if note.Assessment.PrimaryDiagnosis != "고혈압성 두통" {
		t.Errorf("진단 = %q", note.Assessment.PrimaryDiagnosis)
	}
}

func TestScribe_RequiresPatientID(t *testing.T) {
	s := scribe.NewScribe("ko")
	tr := newTestTranscript()
	if _, err := s.Generate(tr, &scribe.GenerateOptions{}); err == nil {
		t.Error("PatientID 없이 통과")
	}
}

func TestScribe_RequiresTranscript(t *testing.T) {
	s := scribe.NewScribe("ko")
	if _, err := s.Generate(nil, &scribe.GenerateOptions{PatientID: "p"}); err == nil {
		t.Error("transcript 없이 통과")
	}
}

func TestScribe_ChiefComplaintExtraction(t *testing.T) {
	s := scribe.NewScribe("ko")
	tr := newTestTranscript()

	note, _ := s.Generate(tr, &scribe.GenerateOptions{PatientID: "p"})
	if note.Subjective.ChiefComplaint == "" {
		t.Error("주증상 미추출")
	}
	if !strings.Contains(note.Subjective.ChiefComplaint, "머리") &&
		!strings.Contains(note.Subjective.ChiefComplaint, "아파") {
		t.Errorf("주증상에 키워드 없음: %q", note.Subjective.ChiefComplaint)
	}
}

func TestScribe_MedicationDetection(t *testing.T) {
	s := scribe.NewScribe("ko")
	tr := newTestTranscript()

	note, _ := s.Generate(tr, &scribe.GenerateOptions{PatientID: "p"})
	if len(note.Subjective.Medications) == 0 {
		t.Error("환자 발화에서 약물 언급 미감지")
	}
}

func TestScribe_AllergyDetection(t *testing.T) {
	s := scribe.NewScribe("ko")
	tr := newTestTranscript()

	note, _ := s.Generate(tr, &scribe.GenerateOptions{PatientID: "p"})
	if len(note.Subjective.Allergies) == 0 {
		t.Error("환자 알레르기 발화 미감지")
	}
}

func TestScribe_PrescriptionInPlan(t *testing.T) {
	s := scribe.NewScribe("ko")
	tr := newTestTranscript()

	note, _ := s.Generate(tr, &scribe.GenerateOptions{PatientID: "p"})
	if len(note.Plan.Medications) == 0 {
		t.Error("처방 발화 미감지")
	}
}

func TestScribe_FollowUpInPlan(t *testing.T) {
	s := scribe.NewScribe("ko")
	tr := newTestTranscript()

	note, _ := s.Generate(tr, &scribe.GenerateOptions{PatientID: "p"})
	if note.Plan.FollowUp == "" {
		t.Error("추적 진료 발화 미감지")
	}
}

func TestScribe_LabResultsIntegration(t *testing.T) {
	s := scribe.NewScribe("ko")
	tr := newTestTranscript()

	labResults := []*scribe.LabResult{
		{TestCode: "2345-7", TestName: "Glucose", Value: 130, Unit: "mg/dL",
			RefLow: 70, RefHigh: 100, Flag: "H", Alpha: 0.98, CILow: 128, CIHigh: 132},
	}

	note, _ := s.Generate(tr, &scribe.GenerateOptions{
		PatientID:  "p",
		LabResults: labResults,
	})
	if len(note.Objective.LabResults) != 1 {
		t.Errorf("LabResults = %d, want 1", len(note.Objective.LabResults))
	}
	if note.Objective.LabResults[0].CILow != 128 {
		t.Errorf("CI 데이터 누락")
	}
}

func TestScribe_VitalSignsIntegration(t *testing.T) {
	s := scribe.NewScribe("ko")
	tr := newTestTranscript()

	vitals := []*scribe.VitalSign{
		{Type: "blood_pressure", Value: "140/90", Unit: "mmHg", Timestamp: time.Now(), Abnormal: true},
		{Type: "heart_rate", Value: "85", Unit: "bpm", Timestamp: time.Now()},
	}

	note, _ := s.Generate(tr, &scribe.GenerateOptions{
		PatientID:  "p",
		VitalSigns: vitals,
	})
	if len(note.Objective.VitalSigns) != 2 {
		t.Errorf("VitalSigns = %d, want 2", len(note.Objective.VitalSigns))
	}
}

func TestScribe_ConfidenceComputation(t *testing.T) {
	s := scribe.NewScribe("ko")
	tr := newTestTranscript()

	note, _ := s.Generate(tr, &scribe.GenerateOptions{PatientID: "p"})
	if note.Confidence < 0.9 || note.Confidence > 1.0 {
		t.Errorf("Confidence = %f, want 0.9-1.0", note.Confidence)
	}
}

func TestTranscript_FilterBySpeaker(t *testing.T) {
	tr := newTestTranscript()
	doctor := tr.FilterBySpeaker("doctor")
	patient := tr.FilterBySpeaker("patient")

	if len(doctor) != 2 {
		t.Errorf("doctor = %d", len(doctor))
	}
	if len(patient) != 2 {
		t.Errorf("patient = %d", len(patient))
	}
}

func TestTranscript_PatientSpeech(t *testing.T) {
	tr := newTestTranscript()
	speech := tr.PatientSpeech()
	if !strings.Contains(speech, "머리") {
		t.Errorf("환자 발화에 '머리' 누락: %q", speech)
	}
}

func TestFormatSOAP_Structure(t *testing.T) {
	note := &scribe.SOAPNote{
		ID:        "soap-1",
		PatientID: "p-001",
		Date:      time.Now().UTC(),
		Subjective: &scribe.Subjective{
			ChiefComplaint:          "두통",
			HistoryOfPresentIllness: "어제부터 시작",
		},
		Objective: &scribe.Objective{
			LabResults: []*scribe.LabResult{
				{TestCode: "2345-7", TestName: "Glucose", Value: 110, Unit: "mg/dL", RefLow: 70, RefHigh: 100, Flag: "H"},
			},
		},
		Assessment: &scribe.Assessment{
			PrimaryDiagnosis: "긴장성 두통",
			ICD10:           "G44.2",
		},
		Plan: &scribe.Plan{
			FollowUp: "1주 후 재내원",
		},
		Confidence: 0.92,
		GeneratedBy: "scribe-v1.0",
	}

	formatted := scribe.FormatSOAP(note)
	if !strings.Contains(formatted, "## S (Subjective)") {
		t.Error("S 섹션 헤더 누락")
	}
	if !strings.Contains(formatted, "## O (Objective)") {
		t.Error("O 섹션 헤더 누락")
	}
	if !strings.Contains(formatted, "## A (Assessment)") {
		t.Error("A 섹션 헤더 누락")
	}
	if !strings.Contains(formatted, "## P (Plan)") {
		t.Error("P 섹션 헤더 누락")
	}
	if !strings.Contains(formatted, "긴장성 두통") {
		t.Error("진단 누락")
	}
}

func TestToFHIRDocumentReference(t *testing.T) {
	note := &scribe.SOAPNote{
		ID:         "n-1",
		PatientID:  "p-001",
		Date:       time.Now().UTC(),
		Confidence: 0.9,
		Subjective: &scribe.Subjective{ChiefComplaint: "두통"},
		Objective:  &scribe.Objective{},
		Assessment: &scribe.Assessment{PrimaryDiagnosis: "두통"},
		Plan:       &scribe.Plan{},
	}

	doc, err := scribe.ToFHIRDocumentReference(note)
	if err != nil {
		t.Fatalf("FHIR 변환 실패: %v", err)
	}
	if doc.ResourceType != "DocumentReference" {
		t.Errorf("ResourceType = %q", doc.ResourceType)
	}
	if doc.Subject.Reference != "Patient/p-001" {
		t.Errorf("Subject = %q", doc.Subject.Reference)
	}
	if len(doc.Content) == 0 {
		t.Error("Content 누락")
	}
	if doc.DocStatus != "preliminary" {
		t.Errorf("DocStatus = %q, want preliminary (검토 미완료)", doc.DocStatus)
	}
}

func TestToFHIRDocumentReference_HumanReviewed(t *testing.T) {
	now := time.Now().UTC()
	note := &scribe.SOAPNote{
		ID: "n", PatientID: "p", Date: now,
		HumanReviewed: true,
		ReviewerID: "Dr. Kim", ReviewedAt: &now,
		Subjective: &scribe.Subjective{}, Objective: &scribe.Objective{},
		Assessment: &scribe.Assessment{}, Plan: &scribe.Plan{},
	}

	doc, _ := scribe.ToFHIRDocumentReference(note)
	if doc.DocStatus != "final" {
		t.Errorf("DocStatus = %q, want final", doc.DocStatus)
	}
}

func TestSortVitalSignsByTime(t *testing.T) {
	t1 := time.Now().UTC()
	t2 := t1.Add(1 * time.Minute)
	t3 := t1.Add(2 * time.Minute)

	vitals := []*scribe.VitalSign{
		{Type: "c", Timestamp: t3},
		{Type: "a", Timestamp: t1},
		{Type: "b", Timestamp: t2},
	}

	scribe.SortVitalSignsByTime(vitals)
	if vitals[0].Type != "a" || vitals[1].Type != "b" || vitals[2].Type != "c" {
		t.Errorf("정렬 실패: %s,%s,%s", vitals[0].Type, vitals[1].Type, vitals[2].Type)
	}
}

func TestScribe_LocaleDefault(t *testing.T) {
	s := scribe.NewScribe("")
	tr := newTestTranscript()
	note, _ := s.Generate(tr, &scribe.GenerateOptions{PatientID: "p"})
	if note.Locale != "ko" {
		t.Errorf("Locale = %q, want ko (default)", note.Locale)
	}
}

// FHIR import 사용 확인
func TestFHIRPackageImport(t *testing.T) {
	// fhir.VersionR5 사용으로 import 검증
	if fhir.VersionR5 == "" {
		t.Error("FHIR VersionR5 미정의")
	}
}

// TestScribe_GenerateFromAudio_RequiresAdapter는 STT 어댑터 미주입 시 거부를 검증합니다.
func TestScribe_GenerateFromAudio_RequiresAdapter(t *testing.T) {
	s := scribe.NewScribe("ko")
	_, err := s.GenerateFromAudio(context.Background(),
		&scribe.AudioInput{Chunks: []*stt.AudioChunk{{Data: []byte{1}}}},
		&scribe.GenerateOptions{PatientID: "p"})
	if err == nil {
		t.Error("STT 미주입인데 통과")
	}
}

// TestScribe_GenerateFromAudio_HappyPath는 오디오 입력 → SOAP 변환을 검증합니다.
func TestScribe_GenerateFromAudio_HappyPath(t *testing.T) {
	s := scribe.NewScribe("ko")
	s.SetSTTAdapter(stt.NewNoopAdapter())

	chunks := []*stt.AudioChunk{
		{
			Data: make([]byte, 32000), Format: stt.FormatWAV,
			SampleRate: 16000, Channels: 1, Locale: "ko-KR",
			SessionID: "audio-test", Speaker: "patient",
		},
		{
			Data: make([]byte, 32000), Format: stt.FormatWAV,
			SampleRate: 16000, Channels: 1, Locale: "ko-KR",
			SessionID: "audio-test", Speaker: "doctor",
		},
	}

	soap, err := s.GenerateFromAudio(context.Background(),
		&scribe.AudioInput{Chunks: chunks},
		&scribe.GenerateOptions{
			PatientID: "P-001", Diagnosis: "두통", ICD10: "G44.2",
		})
	if err != nil {
		t.Fatalf("GenerateFromAudio 실패: %v", err)
	}
	if soap == nil {
		t.Fatal("SOAP nil")
	}
	if soap.Assessment.PrimaryDiagnosis != "두통" {
		t.Errorf("진단 = %q", soap.Assessment.PrimaryDiagnosis)
	}
}

// TestScribe_GenerateFromAudio_RequiresChunks는 빈 오디오 거부를 검증합니다.
func TestScribe_GenerateFromAudio_RequiresChunks(t *testing.T) {
	s := scribe.NewScribe("ko")
	s.SetSTTAdapter(stt.NewNoopAdapter())

	_, err := s.GenerateFromAudio(context.Background(),
		&scribe.AudioInput{},
		&scribe.GenerateOptions{PatientID: "p"})
	if err == nil {
		t.Error("빈 청크 통과")
	}
}
