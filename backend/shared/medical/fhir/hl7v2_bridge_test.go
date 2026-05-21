package fhir_test

import (
	"strings"
	"testing"

	hl7 "github.com/manpasik/backend/shared/hl7-parser"
	"github.com/manpasik/backend/shared/medical/fhir"
)

// 실제 HL7 v2.5 ORU^R01 — 혈당 + HbA1c 측정 결과.
const sampleORU = "MSH|^~\\&|LIS-Lab1|HospA|EMR|HospA|20260521103000||ORU^R01|MSG54321|P|2.5\r" +
	"PID|1||PAT789^^^MRN||김^영수||19880101|M\r" +
	"OBR|1||LAB-456|GLU^Glucose Panel^L\r" +
	"OBX|1|NM|GLU^Glucose^LN||145|mg/dL|70-110|H|||F\r" +
	"OBX|2|NM|HBA1C^HbA1c^LN||6.8|%|<5.7|H|||F\r" +
	"OBX|3|NM|CHOL^Cholesterol^LN||220|mg/dL|<200|H|||F\r"

func mustParseHL7(t *testing.T, raw string) *hl7.Message {
	t.Helper()
	msg, err := hl7.Parse(raw)
	if err != nil {
		t.Fatalf("HL7 파싱 실패: %v", err)
	}
	return msg
}

func TestFromHL7v2Message_BundleStructure(t *testing.T) {
	msg := mustParseHL7(t, sampleORU)
	bundle, err := fhir.FromHL7v2Message(msg, fhir.FromHL7v2Options{})
	if err != nil {
		t.Fatal(err)
	}
	if bundle.Type != "collection" {
		t.Errorf("Bundle.type = %s, collection 기대", bundle.Type)
	}
	if bundle.ID != "hl7v2-MSG54321" {
		t.Errorf("Bundle.id = %s", bundle.ID)
	}
	if len(bundle.Entry) != 3 {
		t.Fatalf("Bundle entries = %d, 3 기대", len(bundle.Entry))
	}
	if bundle.Timestamp.Year() != 2026 || bundle.Timestamp.Month() != 5 {
		t.Errorf("Bundle.timestamp = %v (MSH-7 2026-05 기대)", bundle.Timestamp)
	}
}

func TestFromHL7v2Message_GlucoseObservation(t *testing.T) {
	msg := mustParseHL7(t, sampleORU)
	bundle, _ := fhir.FromHL7v2Message(msg, fhir.FromHL7v2Options{})
	if len(bundle.Entry) == 0 {
		t.Fatal("Entry 없음")
	}
	// 첫 번째 entry = GLU
	res := bundle.Entry[0].Resource
	ob, ok := res.(*fhir.Observation)
	if !ok {
		t.Fatalf("리소스 타입 = %T", res)
	}
	if ob.Status != "final" {
		t.Errorf("status = %s, final 기대", ob.Status)
	}
	if ob.ValueQuantity == nil || ob.ValueQuantity.Value != 145 {
		t.Errorf("ValueQuantity = %+v", ob.ValueQuantity)
	}
	if ob.ValueQuantity.Unit != "mg/dL" {
		t.Errorf("Unit = %s", ob.ValueQuantity.Unit)
	}
	if ob.ValueQuantity.System != "http://unitsofmeasure.org" {
		t.Errorf("UCUM system 누락: %s", ob.ValueQuantity.System)
	}
	// LOINC coding
	if len(ob.Code.Coding) == 0 || ob.Code.Coding[0].System != "http://loinc.org" {
		t.Errorf("LOINC coding 누락: %+v", ob.Code.Coding)
	}
	if ob.Code.Coding[0].Code != "GLU" {
		t.Errorf("Code = %s", ob.Code.Coding[0].Code)
	}
}

func TestFromHL7v2Message_ReferenceRangeRange(t *testing.T) {
	msg := mustParseHL7(t, sampleORU)
	bundle, _ := fhir.FromHL7v2Message(msg, fhir.FromHL7v2Options{})
	gluObs := bundle.Entry[0].Resource.(*fhir.Observation)
	if len(gluObs.ReferenceRange) != 1 {
		t.Fatalf("ReferenceRange 개수 = %d, 1 기대", len(gluObs.ReferenceRange))
	}
	rr := gluObs.ReferenceRange[0]
	if rr.Low == nil || rr.Low.Value != 70 {
		t.Errorf("Low = %+v", rr.Low)
	}
	if rr.High == nil || rr.High.Value != 110 {
		t.Errorf("High = %+v", rr.High)
	}
}

func TestFromHL7v2Message_ReferenceRangeLessThan(t *testing.T) {
	msg := mustParseHL7(t, sampleORU)
	bundle, _ := fhir.FromHL7v2Message(msg, fhir.FromHL7v2Options{})
	// HBA1C: <5.7
	hba1c := bundle.Entry[1].Resource.(*fhir.Observation)
	if len(hba1c.ReferenceRange) != 1 {
		t.Fatalf("HBA1C ReferenceRange = %d", len(hba1c.ReferenceRange))
	}
	rr := hba1c.ReferenceRange[0]
	if rr.Low != nil {
		t.Errorf("Low 가 nil 이어야 함 (<5.7): %+v", rr.Low)
	}
	if rr.High == nil || rr.High.Value != 5.7 {
		t.Errorf("High = %+v", rr.High)
	}
}

func TestFromHL7v2Message_AbnormalFlag(t *testing.T) {
	msg := mustParseHL7(t, sampleORU)
	bundle, _ := fhir.FromHL7v2Message(msg, fhir.FromHL7v2Options{})
	for i, e := range bundle.Entry {
		ob := e.Resource.(*fhir.Observation)
		if len(ob.Interpretation) == 0 {
			t.Errorf("Observation[%d] interpretation 누락 (H 기대)", i)
			continue
		}
		if ob.Interpretation[0].Coding[0].Code != "H" {
			t.Errorf("Observation[%d] interpretation = %s, H 기대",
				i, ob.Interpretation[0].Coding[0].Code)
		}
	}
}

func TestFromHL7v2Message_SubjectReference(t *testing.T) {
	msg := mustParseHL7(t, sampleORU)
	bundle, _ := fhir.FromHL7v2Message(msg, fhir.FromHL7v2Options{})
	ob := bundle.Entry[0].Resource.(*fhir.Observation)
	if !strings.HasSuffix(ob.Subject.Reference, "PAT789") {
		t.Errorf("Subject ref = %s, PAT789 포함 기대", ob.Subject.Reference)
	}
	if !strings.HasPrefix(ob.Subject.Reference, "Patient/") {
		t.Errorf("Subject ref prefix = %s, Patient/ 기대", ob.Subject.Reference)
	}
}

func TestFromHL7v2Message_CustomPrefix(t *testing.T) {
	msg := mustParseHL7(t, sampleORU)
	bundle, _ := fhir.FromHL7v2Message(msg, fhir.FromHL7v2Options{
		PatientResourceIDPrefix: "ManpasikUser/",
	})
	ob := bundle.Entry[0].Resource.(*fhir.Observation)
	if ob.Subject.Reference != "ManpasikUser/PAT789" {
		t.Errorf("custom prefix Subject = %s", ob.Subject.Reference)
	}
}

func TestFromHL7v2Message_R4Version(t *testing.T) {
	msg := mustParseHL7(t, sampleORU)
	bundle, _ := fhir.FromHL7v2Message(msg, fhir.FromHL7v2Options{
		FhirVersion: fhir.VersionR4,
	})
	ob := bundle.Entry[0].Resource.(*fhir.Observation)
	if ob.FhirVersion != fhir.VersionR4 {
		t.Errorf("FhirVersion = %s, R4 기대", ob.FhirVersion)
	}
}

func TestFromHL7v2Message_PerformerFromMSH3(t *testing.T) {
	msg := mustParseHL7(t, sampleORU)
	bundle, _ := fhir.FromHL7v2Message(msg, fhir.FromHL7v2Options{})
	ob := bundle.Entry[0].Resource.(*fhir.Observation)
	if len(ob.Performer) == 0 {
		t.Fatal("Performer 누락")
	}
	if ob.Performer[0].Reference != "Organization/LIS-Lab1" {
		t.Errorf("Performer = %s", ob.Performer[0].Reference)
	}
}

func TestFromHL7v2Message_NilMessage(t *testing.T) {
	if _, err := fhir.FromHL7v2Message(nil, fhir.FromHL7v2Options{}); err == nil {
		t.Error("nil 메시지가 통과")
	}
}

func TestExtractHL7Patient(t *testing.T) {
	msg := mustParseHL7(t, sampleORU)
	p := fhir.ExtractHL7Patient(msg)
	if p == nil {
		t.Fatal("Patient 추출 실패")
	}
	if p.ID != "PAT789" {
		t.Errorf("ID = %s", p.ID)
	}
	if p.Name != "김, 영수" {
		t.Errorf("Name = %s", p.Name)
	}
	if p.Birthday != "19880101" {
		t.Errorf("Birthday = %s", p.Birthday)
	}
	if p.Sex != "M" {
		t.Errorf("Sex = %s", p.Sex)
	}
}

func TestFromHL7v2Message_EmptyOBX(t *testing.T) {
	// PID 만 있고 OBX 가 없는 메시지 — 빈 Bundle 반환.
	// MSH-8 (Security) 만 한 칸 비우고 MSH-9=ORU^R01, MSH-10=MSGempty 가 되도록.
	raw := "MSH|^~\\&|H|H|L|L|20260521||ORU^R01|MSGempty|P|2.5\r" +
		"PID|1||PATX\r"
	msg := mustParseHL7(t, raw)
	bundle, err := fhir.FromHL7v2Message(msg, fhir.FromHL7v2Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(bundle.Entry) != 0 {
		t.Errorf("OBX 없는데 entry = %d", len(bundle.Entry))
	}
	if bundle.ID != "hl7v2-MSGempty" {
		t.Errorf("Bundle.ID = %s", bundle.ID)
	}
}

func TestImportHL7v2_OneStep(t *testing.T) {
	bundle, err := fhir.ImportHL7v2(sampleORU, fhir.FromHL7v2Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(bundle.Entry) != 3 {
		t.Errorf("Entry = %d, 3 기대", len(bundle.Entry))
	}
}

func TestImportHL7v2_InvalidRaw(t *testing.T) {
	if _, err := fhir.ImportHL7v2("not an hl7 message", fhir.FromHL7v2Options{}); err == nil {
		t.Error("invalid raw 통과")
	}
	if _, err := fhir.ImportHL7v2("", fhir.FromHL7v2Options{}); err == nil {
		t.Error("빈 raw 통과")
	}
}

func TestParseHL7ReferenceRange_Edge(t *testing.T) {
	// 직접 호출 테스트는 internal helper 라 lite 검증.
	// ">3.0" 형식이 정상 파싱되는지 ORU 변환을 통해 확인.
	raw := "MSH|^~\\&|H|H|L|L|20260521|||ORU^R01|MID|P|2.5\r" +
		"PID|1||PATY\r" +
		"OBX|1|NM|HDL^HDL^LN||45|mg/dL|>40|N|||F\r"
	msg := mustParseHL7(t, raw)
	bundle, _ := fhir.FromHL7v2Message(msg, fhir.FromHL7v2Options{})
	ob := bundle.Entry[0].Resource.(*fhir.Observation)
	if len(ob.ReferenceRange) != 1 {
		t.Fatal("ReferenceRange 누락")
	}
	rr := ob.ReferenceRange[0]
	if rr.Low == nil || rr.Low.Value != 40 {
		t.Errorf("Low = %+v (>40 기대)", rr.Low)
	}
	if rr.High != nil {
		t.Errorf("High 는 nil 이어야 함: %+v", rr.High)
	}
}
