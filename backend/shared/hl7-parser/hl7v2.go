// Package hl7parser 는 HL7 v2.x 메시지 파서를 제공합니다 (Phase AS-2).
//
// HL7 v2 는 ASCII 텍스트 기반 의료 메시지 표준으로, 전 세계 EHR/LIS 시스템
// 간 통신에 가장 널리 사용됩니다. 메시지는 segment(라인) 단위로 구성되며
// 첫 segment 는 항상 MSH (Message Header) 입니다.
//
// 만파식은 외부 병원 LIS 와의 검사 결과 송수신을 위해 HL7 v2 를 보조 채널로
// 지원합니다 (1차는 FHIR R4/R5 — backend/shared/medical/fhir/).
//
// 구분자 (MSH 첫 5바이트로 동적 결정):
//   - field    : `|` (기본)
//   - component: `^`
//   - repetition: `~`
//   - escape   : `\`
//   - subcomp  : `&`
package hl7parser

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Delimiters 는 메시지에서 사용된 5개 구분자 (MSH-1, MSH-2).
type Delimiters struct {
	Field        byte // 보통 '|'
	Component    byte // 보통 '^'
	Repetition   byte // 보통 '~'
	Escape       byte // 보통 '\\'
	Subcomponent byte // 보통 '&'
}

// DefaultDelimiters 는 HL7 v2 표준 기본값.
func DefaultDelimiters() Delimiters {
	return Delimiters{
		Field:        '|',
		Component:    '^',
		Repetition:   '~',
		Escape:       '\\',
		Subcomponent: '&',
	}
}

// Message 는 파싱된 HL7 v2 메시지.
type Message struct {
	Delimiters Delimiters
	Segments   []Segment
}

// Segment 는 한 줄의 segment (예: MSH, PID, OBX).
//
// Fields[0] 은 segment 이름 (예: "MSH"). MSH segment 의 경우
// HL7 표준에 따라 Fields[1] 은 encoding chars ("^~\&") 이며 Fields[2]
// 부터 sending application 등으로 매핑됩니다.
type Segment struct {
	Name   string
	Fields []Field
}

// Field 는 하나의 field. 반복(repetition)이 있으면 Repetitions 가 2개 이상.
//
// 단순 문자열만 필요하면 `f.Value()` 로 첫 repetition 의 첫 component 의
// 첫 subcomponent 를 가져올 수 있습니다.
type Field struct {
	Repetitions []Repetition
}

// Repetition 은 한 field 내 반복.
type Repetition struct {
	Components []Component
}

// Component 는 ^ 로 구분되는 단위. 더 깊은 subcomponent 가 있으면 Subcomponents 가 2개 이상.
type Component struct {
	Subcomponents []string
}

// ErrInvalidHeader 는 MSH 가 아닌 segment 로 시작하거나 길이 부족 시 반환.
var ErrInvalidHeader = errors.New("MSH segment 헤더가 유효하지 않음")

// Parse 는 HL7 v2 메시지 문자열을 파싱합니다.
//
// segment 구분자는 CR(`\r`), LF(`\n`), CRLF 모두 허용 (실 시스템 호환).
// 비어있는 끝줄은 무시.
func Parse(raw string) (*Message, error) {
	if len(raw) < 8 {
		return nil, ErrInvalidHeader
	}
	if !strings.HasPrefix(raw, "MSH") {
		return nil, ErrInvalidHeader
	}
	// MSH-1: field separator (4번째 문자), MSH-2: encoding chars (5~8번째)
	dl := Delimiters{
		Field:        raw[3],
		Component:    raw[4],
		Repetition:   raw[5],
		Escape:       raw[6],
		Subcomponent: raw[7],
	}
	msg := &Message{Delimiters: dl}

	// segment 분리 (CR/LF/CRLF 모두 허용)
	normalized := strings.ReplaceAll(raw, "\r\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r", "\n")
	for _, line := range strings.Split(normalized, "\n") {
		if line == "" {
			continue
		}
		seg, err := parseSegment(line, dl)
		if err != nil {
			return nil, fmt.Errorf("segment %q: %w", line[:minInt(len(line), 20)], err)
		}
		msg.Segments = append(msg.Segments, seg)
	}
	if len(msg.Segments) == 0 || msg.Segments[0].Name != "MSH" {
		return nil, ErrInvalidHeader
	}
	return msg, nil
}

func parseSegment(line string, dl Delimiters) (Segment, error) {
	parts := strings.Split(line, string(dl.Field))
	if len(parts) == 0 {
		return Segment{}, errors.New("빈 segment")
	}
	name := parts[0]
	if len(name) != 3 {
		return Segment{}, fmt.Errorf("segment 이름 길이 = %d (3 기대)", len(name))
	}
	seg := Segment{Name: name}
	// MSH 의 경우 first field separator 가 field 1 의 역할 — encoding chars 는 field 2
	if name == "MSH" {
		// MSH-1 = field separator 자체 (단일 문자)
		seg.Fields = append(seg.Fields, Field{
			Repetitions: []Repetition{{Components: []Component{{Subcomponents: []string{string(dl.Field)}}}}},
		})
		// MSH-2 부터는 parts[1] 부터
		for _, f := range parts[1:] {
			seg.Fields = append(seg.Fields, parseField(f, dl))
		}
	} else {
		for _, f := range parts[1:] {
			seg.Fields = append(seg.Fields, parseField(f, dl))
		}
	}
	return seg, nil
}

func parseField(raw string, dl Delimiters) Field {
	reps := strings.Split(raw, string(dl.Repetition))
	f := Field{Repetitions: make([]Repetition, 0, len(reps))}
	for _, r := range reps {
		comps := strings.Split(r, string(dl.Component))
		rep := Repetition{Components: make([]Component, 0, len(comps))}
		for _, c := range comps {
			subs := strings.Split(c, string(dl.Subcomponent))
			rep.Components = append(rep.Components, Component{Subcomponents: subs})
		}
		f.Repetitions = append(f.Repetitions, rep)
	}
	return f
}

// Value 는 가장 단순한 단일 문자열 추출: rep[0].comp[0].sub[0].
func (f Field) Value() string {
	if len(f.Repetitions) == 0 || len(f.Repetitions[0].Components) == 0 ||
		len(f.Repetitions[0].Components[0].Subcomponents) == 0 {
		return ""
	}
	return f.Repetitions[0].Components[0].Subcomponents[0]
}

// Components 는 첫 repetition 의 component 들을 단순 문자열 슬라이스로 반환.
func (f Field) Components() []string {
	if len(f.Repetitions) == 0 {
		return nil
	}
	out := make([]string, 0, len(f.Repetitions[0].Components))
	for _, c := range f.Repetitions[0].Components {
		if len(c.Subcomponents) > 0 {
			out = append(out, c.Subcomponents[0])
		} else {
			out = append(out, "")
		}
	}
	return out
}

// Segment 단위 접근 헬퍼

// SegmentByName 은 첫 번째로 매칭되는 segment 반환. 없으면 nil.
func (m *Message) SegmentByName(name string) *Segment {
	for i := range m.Segments {
		if m.Segments[i].Name == name {
			return &m.Segments[i]
		}
	}
	return nil
}

// SegmentsByName 은 매칭되는 모든 segment 반환 (예: 다중 OBX).
func (m *Message) SegmentsByName(name string) []*Segment {
	var out []*Segment
	for i := range m.Segments {
		if m.Segments[i].Name == name {
			out = append(out, &m.Segments[i])
		}
	}
	return out
}

// Field 는 segment 의 1-based field index (HL7 표준 표기) 로 field 접근.
//
// 예: PID-5 (환자 이름) → segment.Field(5).
// MSH segment 의 경우 MSH-1 (field separator) 부터 시작.
func (s *Segment) Field(index int) Field {
	if index < 1 || index > len(s.Fields) {
		return Field{}
	}
	return s.Fields[index-1]
}

// MSH 헬퍼 (메시지 헤더 표준 필드)

// MessageType 은 MSH-9 의 component 0 (예: "ADT", "ORU").
func (m *Message) MessageType() string {
	msh := m.SegmentByName("MSH")
	if msh == nil {
		return ""
	}
	return msh.Field(9).Value()
}

// MessageTypeEvent 는 MSH-9 의 component 1 (예: "A01", "R01").
func (m *Message) MessageTypeEvent() string {
	msh := m.SegmentByName("MSH")
	if msh == nil {
		return ""
	}
	comps := msh.Field(9).Components()
	if len(comps) < 2 {
		return ""
	}
	return comps[1]
}

// SendingApplication 은 MSH-3.
func (m *Message) SendingApplication() string {
	msh := m.SegmentByName("MSH")
	if msh == nil {
		return ""
	}
	return msh.Field(3).Value()
}

// MessageControlID 는 MSH-10 (메시지 고유 ID).
func (m *Message) MessageControlID() string {
	msh := m.SegmentByName("MSH")
	if msh == nil {
		return ""
	}
	return msh.Field(10).Value()
}

// Version 은 MSH-12 (예: "2.5.1").
func (m *Message) Version() string {
	msh := m.SegmentByName("MSH")
	if msh == nil {
		return ""
	}
	return msh.Field(12).Value()
}

// PatientID 는 PID-3 의 첫 component (환자 식별자).
func (m *Message) PatientID() string {
	pid := m.SegmentByName("PID")
	if pid == nil {
		return ""
	}
	return pid.Field(3).Value()
}

// PatientName 은 PID-5 의 family^given 결합 문자열.
func (m *Message) PatientName() string {
	pid := m.SegmentByName("PID")
	if pid == nil {
		return ""
	}
	comps := pid.Field(5).Components()
	switch len(comps) {
	case 0:
		return ""
	case 1:
		return comps[0]
	default:
		// family^given (HL7 표준 — family name 먼저)
		given := strings.TrimSpace(comps[1])
		family := strings.TrimSpace(comps[0])
		if given == "" {
			return family
		}
		return family + ", " + given
	}
}

// ParseHL7Time 은 HL7 v2 datetime (YYYYMMDDHHMMSS[.SSSS][+/-ZZZZ]) → time.Time.
//
// 형식은 truncate 허용: "20260521", "202605211030" 등도 수용.
// 빈 문자열 또는 파싱 실패 시 time.Time{} + 에러.
func ParseHL7Time(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, errors.New("빈 시간")
	}
	// 타임존 분리
	tz := ""
	if i := strings.IndexAny(s, "+-"); i > 4 {
		tz = s[i:]
		s = s[:i]
	}
	layouts := []string{
		"20060102150405.0000",
		"20060102150405",
		"200601021504",
		"2006010215",
		"20060102",
		"200601",
		"2006",
	}
	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			if tz != "" && len(tz) == 5 {
				if loc, lerr := time.Parse("-0700", tz); lerr == nil {
					return t.In(loc.Location()), nil
				}
			}
			return t.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("HL7 시간 형식 인식 실패: %q", s)
}

// OBXResult 는 OBX (Observation/Result) segment 의 핵심 필드.
type OBXResult struct {
	SetID         string // OBX-1
	ValueType     string // OBX-2 (예: "NM" 숫자, "ST" 문자)
	ObservationID string // OBX-3 component 0 (예: "GLU")
	ObservationCS string // OBX-3 component 2 (예: "LN" for LOINC)
	Value         string // OBX-5
	Units         string // OBX-6
	ReferenceRange string // OBX-7
	AbnormalFlags string // OBX-8
	ResultStatus  string // OBX-11 (예: "F" final, "P" preliminary)
}

// ExtractOBX 는 메시지에서 OBX segment 들을 OBXResult 슬라이스로 추출.
func (m *Message) ExtractOBX() []OBXResult {
	segs := m.SegmentsByName("OBX")
	out := make([]OBXResult, 0, len(segs))
	for _, s := range segs {
		obsID := s.Field(3).Components()
		first, third := "", ""
		if len(obsID) > 0 {
			first = obsID[0]
		}
		if len(obsID) > 2 {
			third = obsID[2]
		}
		out = append(out, OBXResult{
			SetID:          s.Field(1).Value(),
			ValueType:      s.Field(2).Value(),
			ObservationID:  first,
			ObservationCS:  third,
			Value:          s.Field(5).Value(),
			Units:          s.Field(6).Value(),
			ReferenceRange: s.Field(7).Value(),
			AbnormalFlags:  s.Field(8).Value(),
			ResultStatus:   s.Field(11).Value(),
		})
	}
	return out
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
