// Package fhir는 HL7 FHIR R4↔R5 Bridge와 POCT 측정값 → FHIR 리소스 변환을 제공합니다.
//
// 만파식 화상진료 P10-09b L3 표준 계층의 핵심 모듈로, POCT1-A 메시지 →
// FHIR Observation/DiagnosticReport/DocumentReference로 변환합니다.
package fhir

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/manpasik/backend/shared/medical/poct"
)

// ============================================================================
// FHIR 버전
// ============================================================================

type Version string

const (
	VersionR4 Version = "4.0.1"
	VersionR5 Version = "5.0.0"
)

// ============================================================================
// 핵심 리소스 (R4/R5 공통 + R5 확장)
// ============================================================================

// Resource는 모든 FHIR 리소스의 공통 인터페이스입니다.
type Resource interface {
	GetResourceType() string
	GetID() string
	GetVersion() Version
}

// Identifier는 리소스 식별자입니다.
type Identifier struct {
	System string `json:"system,omitempty"`
	Value  string `json:"value,omitempty"`
}

// Reference는 다른 리소스에 대한 참조입니다.
type Reference struct {
	Reference string `json:"reference"`
	Display   string `json:"display,omitempty"`
}

// Coding은 코드 시스템의 단일 코드입니다.
type Coding struct {
	System  string `json:"system"`
	Code    string `json:"code"`
	Display string `json:"display,omitempty"`
}

// CodeableConcept는 인간이 읽을 수 있는 코드 표현입니다.
type CodeableConcept struct {
	Coding []Coding `json:"coding,omitempty"`
	Text   string   `json:"text,omitempty"`
}

// Quantity는 단위가 있는 수치값입니다.
type Quantity struct {
	Value  float64 `json:"value"`
	Unit   string  `json:"unit"`
	System string  `json:"system,omitempty"` // UCUM
	Code   string  `json:"code,omitempty"`
}

// Range는 값의 범위입니다.
type Range struct {
	Low  *Quantity `json:"low,omitempty"`
	High *Quantity `json:"high,omitempty"`
}

// ReferenceRange는 정상 범위입니다.
type ReferenceRange struct {
	Low  *Quantity `json:"low,omitempty"`
	High *Quantity `json:"high,omitempty"`
	Type *CodeableConcept `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}

// Period는 시간 구간입니다.
type Period struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end,omitempty"`
}

// ============================================================================
// Observation 리소스
// ============================================================================

// Observation은 단일 측정값/관측값입니다.
type Observation struct {
	ResourceType     string             `json:"resourceType"`
	ID               string             `json:"id"`
	FhirVersion      Version            `json:"meta.fhirVersion"`
	Identifier       []Identifier       `json:"identifier,omitempty"`
	Status           string             `json:"status"` // "final", "preliminary", "amended"
	Category         []CodeableConcept  `json:"category,omitempty"`
	Code             CodeableConcept    `json:"code"`
	Subject          Reference          `json:"subject"`
	Encounter        *Reference         `json:"encounter,omitempty"`
	EffectiveDateTime time.Time         `json:"effectiveDateTime"`
	ValueQuantity    *Quantity          `json:"valueQuantity,omitempty"`
	Interpretation   []CodeableConcept  `json:"interpretation,omitempty"`
	ReferenceRange   []ReferenceRange   `json:"referenceRange,omitempty"`
	Performer        []Reference        `json:"performer,omitempty"`
	Device           *Reference         `json:"device,omitempty"`

	// 만파식 확장
	AlphaCorrection float64 `json:"_manpasik_alpha,omitempty"`
	CILow           float64 `json:"_manpasik_ci_low,omitempty"`
	CIHigh          float64 `json:"_manpasik_ci_high,omitempty"`
}

// GetResourceType은 리소스 타입을 반환합니다.
func (o *Observation) GetResourceType() string { return "Observation" }

// GetID는 ID를 반환합니다.
func (o *Observation) GetID() string { return o.ID }

// GetVersion은 FHIR 버전을 반환합니다.
func (o *Observation) GetVersion() Version { return o.FhirVersion }

// ============================================================================
// DiagnosticReport 리소스
// ============================================================================

// DiagnosticReport은 진단 보고서 (관련 Observation들의 그룹)입니다.
type DiagnosticReport struct {
	ResourceType  string            `json:"resourceType"`
	ID            string            `json:"id"`
	FhirVersion   Version           `json:"meta.fhirVersion"`
	Status        string            `json:"status"` // "final", "partial", "preliminary"
	Category      []CodeableConcept `json:"category,omitempty"`
	Code          CodeableConcept   `json:"code"`
	Subject       Reference         `json:"subject"`
	Encounter     *Reference        `json:"encounter,omitempty"`
	Effective     Period            `json:"effectivePeriod,omitempty"`
	Issued        time.Time         `json:"issued"`
	Performer     []Reference       `json:"performer,omitempty"`
	Result        []Reference       `json:"result,omitempty"` // → Observation
	Conclusion    string            `json:"conclusion,omitempty"`
	ConclusionCode []CodeableConcept `json:"conclusionCode,omitempty"`
}

// GetResourceType은 리소스 타입을 반환합니다.
func (r *DiagnosticReport) GetResourceType() string { return "DiagnosticReport" }

// GetID는 ID를 반환합니다.
func (r *DiagnosticReport) GetID() string { return r.ID }

// GetVersion은 FHIR 버전을 반환합니다.
func (r *DiagnosticReport) GetVersion() Version { return r.FhirVersion }

// ============================================================================
// DocumentReference (Scribe SOAP 출력)
// ============================================================================

// DocumentReference는 외부 문서(SOAP 노트, PDF 등) 참조입니다.
type DocumentReference struct {
	ResourceType string             `json:"resourceType"`
	ID           string             `json:"id"`
	FhirVersion  Version            `json:"meta.fhirVersion"`
	Status       string             `json:"status"` // "current", "superseded"
	DocStatus    string             `json:"docStatus,omitempty"` // "final", "preliminary"
	Type         CodeableConcept    `json:"type"`
	Category     []CodeableConcept  `json:"category,omitempty"`
	Subject      Reference          `json:"subject"`
	Date         time.Time          `json:"date"`
	Author       []Reference        `json:"author,omitempty"`
	Description  string             `json:"description,omitempty"`
	Content      []DocContent       `json:"content"`
	Context      *DocContext        `json:"context,omitempty"`
}

// DocContent는 문서 콘텐츠입니다.
type DocContent struct {
	Attachment Attachment `json:"attachment"`
	Format     *Coding    `json:"format,omitempty"`
}

// Attachment는 첨부 파일입니다.
type Attachment struct {
	ContentType string `json:"contentType"` // "text/plain", "application/pdf"
	Data        []byte `json:"data,omitempty"` // base64
	URL         string `json:"url,omitempty"`
	Size        int    `json:"size,omitempty"`
	Title       string `json:"title,omitempty"`
}

// DocContext는 문서 컨텍스트입니다.
type DocContext struct {
	Encounter []Reference `json:"encounter,omitempty"`
	Period    *Period     `json:"period,omitempty"`
}

// GetResourceType은 리소스 타입을 반환합니다.
func (d *DocumentReference) GetResourceType() string { return "DocumentReference" }

// GetID는 ID를 반환합니다.
func (d *DocumentReference) GetID() string { return d.ID }

// GetVersion은 FHIR 버전을 반환합니다.
func (d *DocumentReference) GetVersion() Version { return d.FhirVersion }

// ============================================================================
// 변환 (POCT → FHIR)
// ============================================================================

// Builder는 POCT 메시지에서 FHIR 리소스를 빌드합니다.
type Builder struct {
	defaultVersion Version
}

// NewBuilder는 새 빌더를 생성합니다.
func NewBuilder(version Version) *Builder {
	if version == "" {
		version = VersionR5
	}
	return &Builder{defaultVersion: version}
}

// FromPOCTMessage는 POCT1-A 메시지에서 FHIR 리소스 묶음을 생성합니다.
//
// 반환:
//   - DiagnosticReport (전체 묶음)
//   - 각 Observation (각 OBX → 1 Observation)
func (b *Builder) FromPOCTMessage(msg *poct.Message, patientRef, deviceRef string) (*DiagnosticReport, []*Observation, error) {
	if msg == nil {
		return nil, nil, errors.New("message is nil")
	}
	if patientRef == "" {
		return nil, nil, errors.New("patient reference required")
	}

	observations := make([]*Observation, 0, len(msg.Observations))
	resultRefs := make([]Reference, 0, len(msg.Observations))

	for _, obx := range msg.Observations {
		obs := b.buildObservation(obx, msg, patientRef, deviceRef)
		observations = append(observations, obs)
		resultRefs = append(resultRefs, Reference{
			Reference: fmt.Sprintf("Observation/%s", obs.ID),
		})
	}

	report := &DiagnosticReport{
		ResourceType: "DiagnosticReport",
		ID:           fmt.Sprintf("dr-%s", msg.MessageID),
		FhirVersion:  b.defaultVersion,
		Status:       mapStatus(msg.Observations),
		Category: []CodeableConcept{{
			Coding: []Coding{{
				System: "http://terminology.hl7.org/CodeSystem/v2-0074",
				Code:   "POCT",
				Display: "Point of Care Testing",
			}},
		}},
		Code: CodeableConcept{
			Coding: []Coding{{
				System:  "http://loinc.org",
				Code:    msg.TestCode,
				Display: msg.TestName,
			}},
		},
		Subject:   Reference{Reference: patientRef},
		Effective: Period{Start: msg.Timestamp},
		Issued:    time.Now().UTC(),
		Result:    resultRefs,
	}

	if deviceRef != "" {
		report.Performer = []Reference{{Reference: deviceRef, Display: msg.Sender}}
	}

	return report, observations, nil
}

func (b *Builder) buildObservation(obx *poct.Observation, msg *poct.Message, patientRef, deviceRef string) *Observation {
	obs := &Observation{
		ResourceType: "Observation",
		ID:           fmt.Sprintf("obs-%s-%d", msg.MessageID, obx.SequenceNum),
		FhirVersion:  b.defaultVersion,
		Status:       mapObsStatus(obx.Status),
		Category: []CodeableConcept{{
			Coding: []Coding{{
				System:  "http://terminology.hl7.org/CodeSystem/observation-category",
				Code:    "laboratory",
				Display: "Laboratory",
			}},
		}},
		Code: CodeableConcept{
			Coding: []Coding{{
				System:  "http://loinc.org",
				Code:    obx.LOINCCode,
				Display: obx.Description,
			}},
		},
		Subject:           Reference{Reference: patientRef},
		EffectiveDateTime: obx.Timestamp,
		ValueQuantity: &Quantity{
			Value:  obx.Value,
			Unit:   obx.Unit,
			System: "http://unitsofmeasure.org",
			Code:   obx.Unit,
		},
	}

	// 만파식 확장: alpha, CI
	if obx.Alpha > 0 {
		obs.AlphaCorrection = obx.Alpha
	}
	if obx.CILow != 0 || obx.CIHigh != 0 {
		obs.CILow = obx.CILow
		obs.CIHigh = obx.CIHigh
	}

	// 정상 범위
	if obx.RefLow > 0 || obx.RefHigh > 0 {
		obs.ReferenceRange = []ReferenceRange{{
			Low:  &Quantity{Value: obx.RefLow, Unit: obx.Unit},
			High: &Quantity{Value: obx.RefHigh, Unit: obx.Unit},
		}}
	}

	// 해석 (high/low/normal)
	if interp := mapInterpretation(obx.Flag); interp != nil {
		obs.Interpretation = []CodeableConcept{*interp}
	}

	// 디바이스
	if deviceRef != "" {
		obs.Device = &Reference{Reference: deviceRef}
	}

	return obs
}

func mapStatus(obs []*poct.Observation) string {
	for _, o := range obs {
		if o.Status == "P" {
			return "preliminary"
		}
	}
	return "final"
}

func mapObsStatus(s string) string {
	switch s {
	case "F":
		return "final"
	case "P":
		return "preliminary"
	case "C":
		return "amended"
	default:
		return "final"
	}
}

func mapInterpretation(flag string) *CodeableConcept {
	codes := map[string]Coding{
		"H": {System: "http://terminology.hl7.org/CodeSystem/v3-ObservationInterpretation", Code: "H", Display: "High"},
		"L": {System: "http://terminology.hl7.org/CodeSystem/v3-ObservationInterpretation", Code: "L", Display: "Low"},
		"A": {System: "http://terminology.hl7.org/CodeSystem/v3-ObservationInterpretation", Code: "A", Display: "Abnormal"},
		"N": {System: "http://terminology.hl7.org/CodeSystem/v3-ObservationInterpretation", Code: "N", Display: "Normal"},
	}
	if c, ok := codes[flag]; ok {
		return &CodeableConcept{Coding: []Coding{c}}
	}
	return nil
}

// ============================================================================
// R4 ↔ R5 변환
// ============================================================================

// Converter는 FHIR 버전 간 변환을 수행합니다.
//
// R5 → R4 다운그레이드 시:
//   - DocumentReference.context.encounter (R5)는 그대로 유지
//   - Observation 만파식 확장 필드는 R4에서도 supported (extension 사용)
type Converter struct{}

// NewConverter는 새 변환기를 생성합니다.
func NewConverter() *Converter { return &Converter{} }

// ToR4는 R5 리소스를 R4로 다운그레이드합니다.
func (c *Converter) ToR4(resource Resource) (Resource, error) {
	switch r := resource.(type) {
	case *Observation:
		clone := *r
		clone.FhirVersion = VersionR4
		return &clone, nil
	case *DiagnosticReport:
		clone := *r
		clone.FhirVersion = VersionR4
		return &clone, nil
	case *DocumentReference:
		clone := *r
		clone.FhirVersion = VersionR4
		return &clone, nil
	}
	return nil, fmt.Errorf("unsupported resource type: %T", resource)
}

// ToR5는 R4 리소스를 R5로 업그레이드합니다.
func (c *Converter) ToR5(resource Resource) (Resource, error) {
	switch r := resource.(type) {
	case *Observation:
		clone := *r
		clone.FhirVersion = VersionR5
		return &clone, nil
	case *DiagnosticReport:
		clone := *r
		clone.FhirVersion = VersionR5
		return &clone, nil
	case *DocumentReference:
		clone := *r
		clone.FhirVersion = VersionR5
		return &clone, nil
	}
	return nil, fmt.Errorf("unsupported resource type: %T", resource)
}

// ============================================================================
// Bundle 빌더 (전송용)
// ============================================================================

// Bundle은 여러 리소스를 묶은 트랜잭션입니다.
type Bundle struct {
	ResourceType string         `json:"resourceType"`
	ID           string         `json:"id"`
	Type         string         `json:"type"` // "transaction", "collection", "document"
	Timestamp    time.Time      `json:"timestamp"`
	Entry        []*BundleEntry `json:"entry"`
}

// BundleEntry는 묶음 항목입니다.
type BundleEntry struct {
	FullURL  string   `json:"fullUrl"`
	Resource Resource `json:"resource"`
	Request  *BundleRequest `json:"request,omitempty"`
}

// BundleRequest는 트랜잭션 요청 정보입니다.
type BundleRequest struct {
	Method string `json:"method"` // "POST", "PUT", "DELETE"
	URL    string `json:"url"`
}

// NewBundle은 새 트랜잭션 묶음을 생성합니다.
func NewBundle(bundleType string, resources ...Resource) *Bundle {
	now := time.Now().UTC()
	bundle := &Bundle{
		ResourceType: "Bundle",
		ID:           fmt.Sprintf("bundle-%d", now.UnixNano()),
		Type:         bundleType,
		Timestamp:    now,
		Entry:        make([]*BundleEntry, 0, len(resources)),
	}
	for _, r := range resources {
		entry := &BundleEntry{
			FullURL:  fmt.Sprintf("urn:uuid:%s", r.GetID()),
			Resource: r,
		}
		if bundleType == "transaction" {
			entry.Request = &BundleRequest{
				Method: "POST",
				URL:    r.GetResourceType(),
			}
		}
		bundle.Entry = append(bundle.Entry, entry)
	}
	return bundle
}

// SortByEffectiveDate는 Observation들을 시간순으로 정렬합니다.
func SortByEffectiveDate(observations []*Observation) {
	sort.Slice(observations, func(i, j int) bool {
		return observations[i].EffectiveDateTime.Before(observations[j].EffectiveDateTime)
	})
}

// ============================================================================
// LOINC 코드 검증
// ============================================================================

// IsValidLOINC는 LOINC 코드 형식을 검증합니다.
//
// 형식: NNNN-N (1~5자리 숫자 + 하이픈 + 1자리 숫자)
func IsValidLOINC(code string) bool {
	if !strings.Contains(code, "-") {
		return false
	}
	parts := strings.Split(code, "-")
	if len(parts) != 2 {
		return false
	}
	if len(parts[0]) < 1 || len(parts[0]) > 5 {
		return false
	}
	if len(parts[1]) != 1 {
		return false
	}
	for _, c := range parts[0] {
		if c < '0' || c > '9' {
			return false
		}
	}
	if parts[1][0] < '0' || parts[1][0] > '9' {
		return false
	}
	return true
}
