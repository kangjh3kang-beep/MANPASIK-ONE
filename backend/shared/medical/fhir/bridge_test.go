package fhir_test

import (
	"testing"
	"time"

	"github.com/manpasik/backend/shared/medical/fhir"
	"github.com/manpasik/backend/shared/medical/poct"
)

func newTestPOCTMessage() *poct.Message {
	return &poct.Message{
		MessageID:   "MSG-001",
		MessageType: "ORU^R01",
		Sender:      "DEVICE-001",
		Timestamp:   time.Date(2026, 4, 30, 12, 0, 0, 0, time.UTC),
		PatientID:   "PATIENT-789",
		OrderID:     "ORDER-100",
		TestCode:    "5811-5",
		TestName:    "Glucose Panel",
		Observations: []*poct.Observation{
			{
				SequenceNum: 1, ValueType: "NM",
				LOINCCode: "2345-7", Description: "Glucose",
				Value: 110.5, Unit: "mg/dL",
				RefLow: 70, RefHigh: 100, Flag: "H", Status: "F",
				Timestamp: time.Date(2026, 4, 30, 12, 0, 30, 0, time.UTC),
				Alpha: 0.98, CILow: 108.5, CIHigh: 112.5,
			},
			{
				SequenceNum: 2, ValueType: "NM",
				LOINCCode: "2093-3", Description: "Cholesterol",
				Value: 190, Unit: "mg/dL",
				RefLow: 0, RefHigh: 200, Flag: "N", Status: "F",
				Timestamp: time.Date(2026, 4, 30, 12, 0, 30, 0, time.UTC),
			},
		},
	}
}

func TestBuilder_FromPOCTMessage_BasicResources(t *testing.T) {
	b := fhir.NewBuilder(fhir.VersionR5)
	msg := newTestPOCTMessage()

	report, observations, err := b.FromPOCTMessage(msg, "Patient/p-789", "Device/dev-001")
	if err != nil {
		t.Fatalf("FromPOCTMessage 실패: %v", err)
	}

	if report.ResourceType != "DiagnosticReport" {
		t.Errorf("ResourceType = %q", report.ResourceType)
	}
	if report.Status != "final" {
		t.Errorf("Status = %q, want final", report.Status)
	}
	if len(observations) != 2 {
		t.Errorf("Observations = %d, want 2", len(observations))
	}
	if len(report.Result) != 2 {
		t.Errorf("Result refs = %d, want 2", len(report.Result))
	}
}

func TestBuilder_PreliminaryStatus(t *testing.T) {
	b := fhir.NewBuilder(fhir.VersionR5)
	msg := newTestPOCTMessage()
	msg.Observations[0].Status = "P" // preliminary

	report, _, _ := b.FromPOCTMessage(msg, "Patient/p", "Device/d")
	if report.Status != "preliminary" {
		t.Errorf("Status = %q, want preliminary", report.Status)
	}
}

func TestBuilder_ObservationFields(t *testing.T) {
	b := fhir.NewBuilder(fhir.VersionR5)
	msg := newTestPOCTMessage()

	_, observations, _ := b.FromPOCTMessage(msg, "Patient/p-789", "Device/dev-001")

	glucose := observations[0]
	if glucose.Code.Coding[0].Code != "2345-7" {
		t.Errorf("LOINC = %q", glucose.Code.Coding[0].Code)
	}
	if glucose.ValueQuantity.Value != 110.5 {
		t.Errorf("Value = %f", glucose.ValueQuantity.Value)
	}
	if glucose.AlphaCorrection != 0.98 {
		t.Errorf("Alpha = %f", glucose.AlphaCorrection)
	}
	if glucose.CILow != 108.5 || glucose.CIHigh != 112.5 {
		t.Errorf("CI = %f-%f", glucose.CILow, glucose.CIHigh)
	}
}

func TestBuilder_Interpretation(t *testing.T) {
	b := fhir.NewBuilder(fhir.VersionR5)
	msg := newTestPOCTMessage()

	_, observations, _ := b.FromPOCTMessage(msg, "Patient/p", "Device/d")

	if len(observations[0].Interpretation) == 0 {
		t.Fatal("H flag 인터프리테이션 누락")
	}
	if observations[0].Interpretation[0].Coding[0].Code != "H" {
		t.Errorf("Interp = %q", observations[0].Interpretation[0].Coding[0].Code)
	}
}

func TestBuilder_ReferenceRange(t *testing.T) {
	b := fhir.NewBuilder(fhir.VersionR5)
	msg := newTestPOCTMessage()

	_, observations, _ := b.FromPOCTMessage(msg, "Patient/p", "Device/d")

	if len(observations[0].ReferenceRange) == 0 {
		t.Fatal("ReferenceRange 누락")
	}
	rng := observations[0].ReferenceRange[0]
	if rng.Low.Value != 70 || rng.High.Value != 100 {
		t.Errorf("Range = %f-%f", rng.Low.Value, rng.High.Value)
	}
}

func TestBuilder_Validation_NilMessage(t *testing.T) {
	b := fhir.NewBuilder(fhir.VersionR5)
	if _, _, err := b.FromPOCTMessage(nil, "Patient/p", "Device/d"); err == nil {
		t.Error("nil 메시지 통과")
	}
}

func TestBuilder_Validation_NoPatient(t *testing.T) {
	b := fhir.NewBuilder(fhir.VersionR5)
	msg := newTestPOCTMessage()
	if _, _, err := b.FromPOCTMessage(msg, "", "Device/d"); err == nil {
		t.Error("환자 ref 없이 통과")
	}
}

func TestConverter_R5ToR4(t *testing.T) {
	c := fhir.NewConverter()
	r5 := &fhir.Observation{ID: "x", FhirVersion: fhir.VersionR5}

	r4, err := c.ToR4(r5)
	if err != nil {
		t.Fatalf("ToR4 실패: %v", err)
	}
	if r4.GetVersion() != fhir.VersionR4 {
		t.Errorf("Version = %q, want %q", r4.GetVersion(), fhir.VersionR4)
	}
}

func TestConverter_R4ToR5(t *testing.T) {
	c := fhir.NewConverter()
	r4 := &fhir.Observation{ID: "x", FhirVersion: fhir.VersionR4}

	r5, err := c.ToR5(r4)
	if err != nil {
		t.Fatalf("ToR5 실패: %v", err)
	}
	if r5.GetVersion() != fhir.VersionR5 {
		t.Errorf("Version = %q", r5.GetVersion())
	}
}

func TestConverter_DiagnosticReport(t *testing.T) {
	c := fhir.NewConverter()
	dr := &fhir.DiagnosticReport{ID: "dr", FhirVersion: fhir.VersionR5}

	converted, err := c.ToR4(dr)
	if err != nil {
		t.Fatalf("ToR4 실패: %v", err)
	}
	if converted.GetVersion() != fhir.VersionR4 {
		t.Errorf("Version = %q", converted.GetVersion())
	}
	if converted.GetResourceType() != "DiagnosticReport" {
		t.Errorf("ResourceType = %q", converted.GetResourceType())
	}
}

func TestNewBundle_Transaction(t *testing.T) {
	obs := &fhir.Observation{ID: "obs-1", ResourceType: "Observation", FhirVersion: fhir.VersionR5}
	dr := &fhir.DiagnosticReport{ID: "dr-1", ResourceType: "DiagnosticReport", FhirVersion: fhir.VersionR5}

	bundle := fhir.NewBundle("transaction", obs, dr)

	if bundle.Type != "transaction" {
		t.Errorf("Type = %q", bundle.Type)
	}
	if len(bundle.Entry) != 2 {
		t.Errorf("Entry = %d, want 2", len(bundle.Entry))
	}
	if bundle.Entry[0].Request == nil {
		t.Error("transaction Request 누락")
	}
	if bundle.Entry[0].Request.Method != "POST" {
		t.Errorf("Method = %q", bundle.Entry[0].Request.Method)
	}
}

func TestNewBundle_Collection(t *testing.T) {
	obs := &fhir.Observation{ID: "obs-1", ResourceType: "Observation"}
	bundle := fhir.NewBundle("collection", obs)
	if bundle.Entry[0].Request != nil {
		t.Error("collection에 Request가 있음")
	}
}

func TestSortByEffectiveDate(t *testing.T) {
	t1 := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 4, 2, 0, 0, 0, 0, time.UTC)
	t3 := time.Date(2026, 4, 3, 0, 0, 0, 0, time.UTC)

	obs := []*fhir.Observation{
		{ID: "c", EffectiveDateTime: t3},
		{ID: "a", EffectiveDateTime: t1},
		{ID: "b", EffectiveDateTime: t2},
	}

	fhir.SortByEffectiveDate(obs)

	if obs[0].ID != "a" || obs[1].ID != "b" || obs[2].ID != "c" {
		t.Errorf("정렬 순서 = %s,%s,%s", obs[0].ID, obs[1].ID, obs[2].ID)
	}
}

func TestIsValidLOINC(t *testing.T) {
	valid := []string{"2345-7", "5811-5", "12345-6", "1-0"}
	invalid := []string{"", "abc-1", "2345", "2345-", "-7", "12345-67", "abcdef-1"}

	for _, v := range valid {
		if !fhir.IsValidLOINC(v) {
			t.Errorf("정상 LOINC %q 거부됨", v)
		}
	}
	for _, v := range invalid {
		if fhir.IsValidLOINC(v) {
			t.Errorf("잘못된 LOINC %q 통과됨", v)
		}
	}
}

func TestObservation_GetMethods(t *testing.T) {
	obs := &fhir.Observation{
		ResourceType: "Observation",
		ID:           "obs-x",
		FhirVersion:  fhir.VersionR5,
	}
	if obs.GetResourceType() != "Observation" {
		t.Error("GetResourceType")
	}
	if obs.GetID() != "obs-x" {
		t.Error("GetID")
	}
	if obs.GetVersion() != fhir.VersionR5 {
		t.Error("GetVersion")
	}
}

func TestDocumentReference_BuildAndSerialize(t *testing.T) {
	doc := &fhir.DocumentReference{
		ResourceType: "DocumentReference",
		ID:           "doc-soap-1",
		FhirVersion:  fhir.VersionR5,
		Status:       "current",
		DocStatus:    "final",
		Type: fhir.CodeableConcept{
			Coding: []fhir.Coding{{
				System: "http://loinc.org", Code: "11488-4", Display: "Consultation note",
			}},
		},
		Subject: fhir.Reference{Reference: "Patient/p-789"},
		Date:    time.Now().UTC(),
		Content: []fhir.DocContent{{
			Attachment: fhir.Attachment{
				ContentType: "text/plain",
				Title:       "SOAP Note",
				Data:        []byte("S: 환자 두통 호소\nO: BP 130/80\nA: 두통\nP: 진통제"),
			},
		}},
	}

	if doc.GetResourceType() != "DocumentReference" {
		t.Error("ResourceType mismatch")
	}
	if len(doc.Content) != 1 {
		t.Errorf("Content = %d", len(doc.Content))
	}
}

func TestBuilder_DefaultVersion(t *testing.T) {
	b := fhir.NewBuilder("")
	msg := newTestPOCTMessage()
	report, _, _ := b.FromPOCTMessage(msg, "Patient/p", "Device/d")
	if report.FhirVersion != fhir.VersionR5 {
		t.Errorf("기본 Version = %q, want R5", report.FhirVersion)
	}
}

func TestBuilder_DeviceReference(t *testing.T) {
	b := fhir.NewBuilder(fhir.VersionR5)
	msg := newTestPOCTMessage()

	_, observations, _ := b.FromPOCTMessage(msg, "Patient/p", "Device/dev-001")
	if observations[0].Device == nil {
		t.Error("Device 참조 누락")
	}
	if observations[0].Device.Reference != "Device/dev-001" {
		t.Errorf("Device = %q", observations[0].Device.Reference)
	}
}
