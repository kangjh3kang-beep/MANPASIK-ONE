package fhir

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	hl7 "github.com/manpasik/backend/shared/hl7-parser"
)

// ============================================================================
// HL7 v2 ↔ FHIR 브릿지 (Phase AT)
// ============================================================================
//
// HL7 v2 메시지 (ADT, ORU 등) → FHIR Bundle 자동 변환.
// 만파식은 외부 LIS/EHR 가 HL7 v2 만 지원하는 경우에도 내부적으로 FHIR
// Observation/DiagnosticReport 형식으로 통합 저장하기 위해 본 브릿지를 사용.
//
// 주요 매핑:
//
//	ORU^R01 메시지       → Bundle (type=collection)
//	  ├─ MSH-9, MSH-10  → Bundle.identifier
//	  ├─ MSH-12         → Bundle.meta.fhirVersion (참조용)
//	  ├─ PID-3 / PID-5  → HL7Patient 메타 (Subject 참조 식별자)
//	  └─ OBX 각각        → Observation (LOINC code, value, units, range)

// ImportHL7v2 는 raw HL7 메시지 문자열 → FHIR Bundle 변환 1-step 헬퍼.
//
// `hl7.Parse(raw)` + `FromHL7v2Message(msg, opts)` 를 결합한 단축 함수로,
// 외부 시스템 통합 (gateway HTTP 핸들러 등) 에서 raw 메시지를 받아 즉시
// FHIR Bundle 로 변환할 때 사용합니다.
//
// 빈 문자열 / 유효하지 않은 헤더 / 변환 실패 시 에러 반환.
func ImportHL7v2(raw string, opts FromHL7v2Options) (*Bundle, error) {
	msg, err := hl7.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("HL7 파싱: %w", err)
	}
	return FromHL7v2Message(msg, opts)
}

// HL7Patient 는 PID segment 에서 추출된 환자 메타.
//
// 본격적인 Patient 리소스 빌드 대신 외부 시스템 식별자로 Subject 참조에 사용.
// Subject 식별자가 이미 등록된 경우 그대로 활용 (멱등성).
type HL7Patient struct {
	ID       string // PID-3 (Patient Identifier List, 첫 component)
	Name     string // PID-5 결합 "family, given"
	Birthday string // PID-7 (raw YYYYMMDD)
	Sex      string // PID-8 ("M"/"F"/"U")
}

// ExtractHL7Patient 는 메시지에서 환자 정보 추출. PID segment 없으면 nil.
func ExtractHL7Patient(msg *hl7.Message) *HL7Patient {
	pid := msg.SegmentByName("PID")
	if pid == nil {
		return nil
	}
	return &HL7Patient{
		ID:       pid.Field(3).Value(),
		Name:     msg.PatientName(),
		Birthday: pid.Field(7).Value(),
		Sex:      pid.Field(8).Value(),
	}
}

// FromHL7v2Options 는 변환 옵션.
type FromHL7v2Options struct {
	// FhirVersion 은 출력 리소스 버전 (기본 R5).
	FhirVersion Version

	// PatientResourceIDPrefix 는 Bundle 내부 Patient 참조의 prefix
	// (예: "Patient/" + PID-3 값). 비워두면 "Patient/" 사용.
	PatientResourceIDPrefix string

	// DefaultEffectiveDateTime 은 OBX-14 (Date/Time of Observation) 가 비어있을 때
	// 사용할 기본값. 비어있으면 MSH-7 의 메시지 발생 시간 사용, 그것도 없으면 time.Now().
	DefaultEffectiveDateTime time.Time
}

func (o *FromHL7v2Options) resolveVersion() Version {
	if o.FhirVersion == "" {
		return VersionR5
	}
	return o.FhirVersion
}

func (o *FromHL7v2Options) patientRef(patientID string) string {
	prefix := o.PatientResourceIDPrefix
	if prefix == "" {
		prefix = "Patient/"
	}
	return prefix + patientID
}

// FromHL7v2Message 는 HL7 v2 메시지를 FHIR Bundle 로 변환.
//
// 지원 메시지: ORU^R01 (검사 결과). 다른 타입은 OBX 가 없으면 빈 Bundle 반환.
// 메시지가 nil 이거나 MSH 가 누락된 경우 에러 반환.
func FromHL7v2Message(msg *hl7.Message, opts FromHL7v2Options) (*Bundle, error) {
	if msg == nil {
		return nil, errors.New("nil HL7 메시지")
	}
	msh := msg.SegmentByName("MSH")
	if msh == nil {
		return nil, errors.New("MSH segment 누락")
	}
	version := opts.resolveVersion()

	// MSH-7 → Bundle.timestamp / Observation effective default
	defaultTime := opts.DefaultEffectiveDateTime
	if mshTime, err := hl7.ParseHL7Time(msh.Field(7).Value()); err == nil {
		if defaultTime.IsZero() {
			defaultTime = mshTime
		}
	}
	if defaultTime.IsZero() {
		defaultTime = time.Now().UTC()
	}

	// Patient 식별자 (있으면 사용)
	patient := ExtractHL7Patient(msg)
	patientRef := ""
	if patient != nil && patient.ID != "" {
		patientRef = opts.patientRef(patient.ID)
	}

	// OBX 들 → Observation 들
	obxs := msg.ExtractOBX()
	resources := make([]Resource, 0, len(obxs))
	for i, obx := range obxs {
		ob, err := buildObservationFromOBX(obx, msh, patientRef, defaultTime, version, i+1)
		if err != nil {
			return nil, fmt.Errorf("OBX %d 변환 실패: %w", i+1, err)
		}
		resources = append(resources, ob)
	}

	bundle := NewBundle("collection", resources...)
	bundle.Timestamp = defaultTime
	// MSH-10 (Message Control ID) → Bundle.id 로 사용 (멱등성)
	if mcid := msh.Field(10).Value(); mcid != "" {
		bundle.ID = "hl7v2-" + mcid
	}
	return bundle, nil
}

// buildObservationFromOBX 는 OBX → Observation 단일 변환.
func buildObservationFromOBX(
	obx hl7.OBXResult,
	msh *hl7.Segment,
	patientRef string,
	defaultTime time.Time,
	version Version,
	idx int,
) (*Observation, error) {
	id := obx.SetID
	if id == "" {
		id = strconv.Itoa(idx)
	}
	ob := &Observation{
		ResourceType: "Observation",
		ID:           "hl7v2-obx-" + id,
		FhirVersion:  version,
		Status:       mapHL7ResultStatus(obx.ResultStatus),
		Code: CodeableConcept{
			Coding: []Coding{{
				System:  resolveLOINCSystem(obx.ObservationCS),
				Code:    obx.ObservationID,
				Display: obx.ObservationID,
			}},
			Text: obx.ObservationID,
		},
		EffectiveDateTime: defaultTime,
	}
	if patientRef != "" {
		ob.Subject = Reference{Reference: patientRef}
	}
	if obx.Value != "" && (obx.ValueType == "NM" || obx.ValueType == "SN") {
		if v, err := strconv.ParseFloat(obx.Value, 64); err == nil {
			ob.ValueQuantity = &Quantity{
				Value:  v,
				Unit:   obx.Units,
				System: "http://unitsofmeasure.org", // UCUM
				Code:   obx.Units,
			}
		}
	}
	if obx.AbnormalFlags != "" {
		if interp := mapInterpretation(obx.AbnormalFlags); interp != nil {
			ob.Interpretation = []CodeableConcept{*interp}
		}
	}
	if rr := parseHL7ReferenceRange(obx.ReferenceRange, obx.Units); rr != nil {
		ob.ReferenceRange = []ReferenceRange{*rr}
	}
	// MSH-3 (Sending Application) → performer 참조 (식별 용도)
	if sendingApp := msh.Field(3).Value(); sendingApp != "" {
		ob.Performer = []Reference{{Reference: "Organization/" + sendingApp}}
	}
	return ob, nil
}

// mapHL7ResultStatus 는 OBX-11 (HL7 result status) → FHIR Observation.status.
//
// HL7: F=final, P=preliminary, C=corrected, X=cancelled, S=partial.
// FHIR: registered | preliminary | final | amended | corrected | cancelled.
func mapHL7ResultStatus(s string) string {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "F":
		return "final"
	case "P":
		return "preliminary"
	case "C":
		return "corrected"
	case "X":
		return "cancelled"
	case "S":
		return "preliminary"
	default:
		return "final"
	}
}

// resolveLOINCSystem 은 OBX-3 component 2 (coding system) 에 따라 URI 반환.
//
// LOINC 가 가장 일반적. SNOMED-CT, IC-D-10 도 지원.
func resolveLOINCSystem(cs string) string {
	switch strings.ToUpper(strings.TrimSpace(cs)) {
	case "LN", "LOINC":
		return "http://loinc.org"
	case "SCT", "SNOMED":
		return "http://snomed.info/sct"
	case "I10", "ICD10", "ICD-10":
		return "http://hl7.org/fhir/sid/icd-10"
	case "I11", "ICD11", "ICD-11":
		return "http://id.who.int/icd/release/11/mms"
	case "":
		// coding system 미지정 — 로컬 코드로 처리
		return "urn:manpasik:local-code"
	default:
		return "urn:hl7:cs:" + strings.ToLower(cs)
	}
}

// parseHL7ReferenceRange 는 OBX-7 의 정상 범위 문자열을 파싱.
//
// 형식: "low-high" (예: "70-110"), "<high" (예: "<5.7"), ">low" (예: ">3.0").
// 단위는 OBX-6 에서 별도로 받아서 Quantity 에 부착.
func parseHL7ReferenceRange(raw, units string) *ReferenceRange {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	mk := func(v float64) *Quantity {
		return &Quantity{Value: v, Unit: units, System: "http://unitsofmeasure.org", Code: units}
	}
	if strings.HasPrefix(raw, "<") {
		if v, err := strconv.ParseFloat(strings.TrimPrefix(raw, "<"), 64); err == nil {
			return &ReferenceRange{High: mk(v)}
		}
		return nil
	}
	if strings.HasPrefix(raw, ">") {
		if v, err := strconv.ParseFloat(strings.TrimPrefix(raw, ">"), 64); err == nil {
			return &ReferenceRange{Low: mk(v)}
		}
		return nil
	}
	// "low-high" 형식
	parts := strings.SplitN(raw, "-", 2)
	if len(parts) != 2 {
		return nil
	}
	low, errL := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	high, errH := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if errL != nil || errH != nil {
		return nil
	}
	return &ReferenceRange{Low: mk(low), High: mk(high)}
}
