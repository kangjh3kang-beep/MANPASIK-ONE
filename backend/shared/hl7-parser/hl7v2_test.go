package hl7parser_test

import (
	"strings"
	"testing"

	hl7 "github.com/manpasik/backend/shared/hl7-parser"
)

// 실제 HL7 v2.5 ADT^A01 (환자 입원) 예시 메시지.
const sampleADT_A01 = "MSH|^~\\&|HIS|HospA|LIS|LabA|20260521103000||ADT^A01|MSG00001|P|2.5\r" +
	"EVN|A01|20260521103000\r" +
	"PID|1||PATID1234^^^MRN||김^철수^^^Mr||19800505|M|||서울시 강남구\r" +
	"PV1|1|I|2000^2012^01||||004777^홍^길동^A.|||SUR||||||||S|400123\r"

// 실제 HL7 v2.5 ORU^R01 (검사 결과) 예시 메시지 — 혈당 측정.
const sampleORU_R01 = "MSH|^~\\&|LIS|HospA|EMR|HospA|20260521103000||ORU^R01|MSG00002|P|2.5\r" +
	"PID|1||PATID5678^^^MRN||이^영희||19900101|F\r" +
	"OBR|1||LAB123|GLU^Glucose^L\r" +
	"OBX|1|NM|GLU^Glucose^LN||95|mg/dL|70-110|N|||F\r" +
	"OBX|2|NM|HBA1C^HbA1c^LN||5.6|%|<5.7|N|||F\r"

func TestParse_BasicADT(t *testing.T) {
	msg, err := hl7.Parse(sampleADT_A01)
	if err != nil {
		t.Fatal(err)
	}
	if len(msg.Segments) != 4 {
		t.Errorf("segments = %d, 4 기대", len(msg.Segments))
	}
	if msg.Segments[0].Name != "MSH" {
		t.Errorf("first segment = %s", msg.Segments[0].Name)
	}
	if msg.MessageType() != "ADT" {
		t.Errorf("MessageType = %s, ADT 기대", msg.MessageType())
	}
	if msg.MessageTypeEvent() != "A01" {
		t.Errorf("MessageTypeEvent = %s, A01 기대", msg.MessageTypeEvent())
	}
	if msg.Version() != "2.5" {
		t.Errorf("Version = %s", msg.Version())
	}
	if msg.MessageControlID() != "MSG00001" {
		t.Errorf("MessageControlID = %s", msg.MessageControlID())
	}
}

func TestParse_PatientFields(t *testing.T) {
	msg, _ := hl7.Parse(sampleADT_A01)
	if msg.PatientID() != "PATID1234" {
		t.Errorf("PatientID = %s", msg.PatientID())
	}
	// PID-5: 김^철수^^^Mr → family=김, given=철수
	if name := msg.PatientName(); name != "김, 철수" {
		t.Errorf("PatientName = %s", name)
	}
}

func TestParse_ORU_OBXExtraction(t *testing.T) {
	msg, err := hl7.Parse(sampleORU_R01)
	if err != nil {
		t.Fatal(err)
	}
	obxs := msg.ExtractOBX()
	if len(obxs) != 2 {
		t.Fatalf("OBX 개수 = %d, 2 기대", len(obxs))
	}
	if obxs[0].ObservationID != "GLU" || obxs[0].Value != "95" || obxs[0].Units != "mg/dL" {
		t.Errorf("OBX[0] = %+v", obxs[0])
	}
	if obxs[0].ObservationCS != "LN" {
		t.Errorf("LOINC coding system 누락: %+v", obxs[0])
	}
	if obxs[0].ResultStatus != "F" {
		t.Errorf("ResultStatus = %s, F 기대", obxs[0].ResultStatus)
	}
	if obxs[1].ObservationID != "HBA1C" || obxs[1].Value != "5.6" || obxs[1].Units != "%" {
		t.Errorf("OBX[1] = %+v", obxs[1])
	}
}

func TestParse_NonStandardDelimiters(t *testing.T) {
	// HL7 표준은 동적 구분자 허용 — MSH-1/2 로 정의. field='#', component='$'.
	// 모든 field separator 도 '#' 로 통일 필요.
	// MSH 필드 카운트: MSH-1=#, MSH-2=$~\&, MSH-3=HIS, ..., MSH-7=20260521,
	// MSH-8=(empty), MSH-9=ADT$A01 → 20260521 다음 빈 필드 1개.
	raw := "MSH#$~\\&#HIS#HospA#LIS#LabA#20260521##ADT$A01#MID#P#2.5\r" +
		"PID#1##PATID9999"
	msg, err := hl7.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if msg.Delimiters.Field != '#' || msg.Delimiters.Component != '$' {
		t.Errorf("delimiters = %+v", msg.Delimiters)
	}
	if msg.MessageType() != "ADT" || msg.MessageTypeEvent() != "A01" {
		t.Errorf("동적 구분자 메시지 타입 파싱 실패: type=%s event=%s",
			msg.MessageType(), msg.MessageTypeEvent())
	}
	if msg.PatientID() != "PATID9999" {
		t.Errorf("PatientID = %s", msg.PatientID())
	}
}

func TestParse_CRLFAccepted(t *testing.T) {
	raw := strings.ReplaceAll(sampleADT_A01, "\r", "\r\n")
	msg, err := hl7.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(msg.Segments) != 4 {
		t.Errorf("CRLF 정규화 실패: segments = %d", len(msg.Segments))
	}
}

func TestParse_InvalidHeader(t *testing.T) {
	_, err := hl7.Parse("PID|1||X")
	if err == nil {
		t.Error("MSH 미시작 메시지 통과")
	}
	_, err = hl7.Parse("")
	if err == nil {
		t.Error("빈 메시지 통과")
	}
	_, err = hl7.Parse("MSH")
	if err == nil {
		t.Error("길이 부족 메시지 통과")
	}
}

func TestSegmentByName_Missing(t *testing.T) {
	msg, _ := hl7.Parse(sampleADT_A01)
	if msg.SegmentByName("OBX") != nil {
		t.Error("OBX 가 없는데 반환됨")
	}
}

func TestField_Value_EmptyField(t *testing.T) {
	// 빈 field 가 있는 PID
	raw := "MSH|^~\\&|HIS|H|L|L|20260521|||ADT^A01|M|P|2.5\r" +
		"PID|1|||emptyfamily^||\r"
	msg, _ := hl7.Parse(raw)
	pid := msg.SegmentByName("PID")
	if pid == nil {
		t.Fatal("PID 없음")
	}
	// PID-2 는 빈 field
	if v := pid.Field(2).Value(); v != "" {
		t.Errorf("빈 field Value = %q (\"\" 기대)", v)
	}
	// PID-3 도 빈
	if v := pid.Field(3).Value(); v != "" {
		t.Errorf("빈 field Value = %q", v)
	}
}

func TestParseHL7Time(t *testing.T) {
	cases := []struct {
		input string
		ok    bool
		year  int
	}{
		{"20260521103000", true, 2026},
		{"20260521", true, 2026},
		{"202605211030", true, 2026},
		{"20260521103000.5000", true, 2026},
		{"20260521103000+0900", true, 2026},
		{"", false, 0},
		{"notatime", false, 0},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			tm, err := hl7.ParseHL7Time(c.input)
			if c.ok {
				if err != nil {
					t.Errorf("err = %v", err)
				}
				if tm.Year() != c.year {
					t.Errorf("year = %d, %d 기대", tm.Year(), c.year)
				}
			} else if err == nil {
				t.Errorf("실패해야 하는데 통과: %q", c.input)
			}
		})
	}
}

func TestParse_RepetitionField(t *testing.T) {
	// PID-13 (전화번호 — Home) 에 ~ 로 두 번호 반복.
	// parts[0]=PID, parts[N]=field N. PID-13 = parts[13]. → 13 pipes total.
	// 5 pipes 가 김^철수 까지 + 8 pipes 가 그 이후 = 13. (Fields[12] = phone)
	raw := "MSH|^~\\&|H|H|L|L|20260521|||ADT^A01|M|P|2.5\r" +
		"PID|1||X||김^철수||||||||010-1111-2222~010-3333-4444\r"
	msg, _ := hl7.Parse(raw)
	pid := msg.SegmentByName("PID")
	if pid == nil {
		t.Fatal("PID 없음")
	}
	phone := pid.Field(13)
	if len(phone.Repetitions) != 2 {
		t.Errorf("Repetitions = %d, 2 기대", len(phone.Repetitions))
	}
	if phone.Repetitions[0].Components[0].Subcomponents[0] != "010-1111-2222" {
		t.Errorf("첫 반복 = %s", phone.Repetitions[0].Components[0].Subcomponents[0])
	}
	if phone.Repetitions[1].Components[0].Subcomponents[0] != "010-3333-4444" {
		t.Errorf("둘째 반복 = %s", phone.Repetitions[1].Components[0].Subcomponents[0])
	}
}

func TestSegmentsByName_MultipleOBX(t *testing.T) {
	msg, _ := hl7.Parse(sampleORU_R01)
	obx := msg.SegmentsByName("OBX")
	if len(obx) != 2 {
		t.Errorf("OBX 슬라이스 = %d", len(obx))
	}
}
