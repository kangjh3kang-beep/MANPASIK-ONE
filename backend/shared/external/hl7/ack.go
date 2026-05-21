package hl7mllp

import (
	"fmt"
	"strings"
	"time"

	hl7 "github.com/manpasik/backend/shared/hl7-parser"
)

// ============================================================================
// HL7 v2 ACK 메시지 생성 (Phase AX-4)
// ============================================================================
//
// ACK 표준 형식 (HL7 v2.5):
//
//	MSH|^~\&|<sending>|<facility>|<receiving>|<facility>|<datetime>||ACK^<event>|<control_id>|P|2.5
//	MSA|<ack_code>|<orig_control_id>|<text>
//
// ack_code:
//   AA = Application Accept (정상 수신 + 처리 성공)
//   AE = Application Error  (수신했으나 처리 실패)
//   AR = Application Reject (메시지 자체 거부 — 파싱 실패 / 권한 / 형식)
//
// 정상 흐름 — 외부 시스템이 ACK 가 없으면 메시지를 재전송할 수 있음. 따라서
// 처리 결과와 ACK 송신은 강하게 결합되어야 함.

// BuildACK 는 수신 메시지에 대한 ACK 본문 생성 (MLLP framing 은 미포함).
//
// originalMsg 가 nil 이면 (파싱 불가) 최소 정보로 ACK 작성. 이 경우 ack_code 는
// 보통 "AR".
func BuildACK(originalMsg *hl7.Message, ackCode, textMsg string) string {
	if ackCode == "" {
		ackCode = "AA"
	}
	// 원본 메시지에서 가능한 정보 추출
	var (
		sendingApp, receivingApp string
		eventCode                string
		origControlID            string
		version                  = "2.5"
	)
	if originalMsg != nil {
		// 원본의 수신측이 ACK 의 송신측, 그 반대도 마찬가지 (역방향).
		sendingApp = safeMSH(originalMsg, 5)  // MSH-5: receiving app → ACK sender
		receivingApp = safeMSH(originalMsg, 3) // MSH-3: sending app  → ACK receiver
		origControlID = originalMsg.MessageControlID()
		event := originalMsg.MessageTypeEvent()
		if event != "" {
			eventCode = event
		}
		if v := originalMsg.Version(); v != "" {
			version = v
		}
	}
	if sendingApp == "" {
		sendingApp = "MANPASIK"
	}
	if receivingApp == "" {
		receivingApp = "EXTERNAL"
	}
	if eventCode == "" {
		eventCode = "R01"
	}

	now := time.Now().UTC().Format("20060102150405")
	ackControlID := fmt.Sprintf("ACK-%s", now)

	// MSH
	msh := strings.Join([]string{
		"MSH",
		"^~\\&",
		sendingApp,
		"MANPASIK_FAC",
		receivingApp,
		"EXTERNAL_FAC",
		now,
		"",
		"ACK^" + eventCode,
		ackControlID,
		"P",
		version,
	}, "|")

	// MSA
	msa := strings.Join([]string{
		"MSA",
		ackCode,
		origControlID,
		sanitizeText(textMsg),
	}, "|")

	return msh + "\r" + msa + "\r"
}

// safeMSH 는 MSH-N field 안전 접근 (nil 또는 부족 필드 시 빈 문자열).
func safeMSH(msg *hl7.Message, index int) string {
	if msg == nil {
		return ""
	}
	seg := msg.SegmentByName("MSH")
	if seg == nil {
		return ""
	}
	return seg.Field(index).Value()
}

// sanitizeText 는 HL7 구분자 (| ^ ~ \ &) 를 제거해 ACK 본문이 깨지지 않게 합니다.
func sanitizeText(s string) string {
	if s == "" {
		return ""
	}
	r := strings.NewReplacer(
		"|", " ", "^", " ", "~", " ", "\\", " ", "&", " ",
		"\r", " ", "\n", " ",
	)
	return r.Replace(s)
}

// ParseACK 는 받은 ACK 메시지에서 ackCode + textMsg 추출.
//
// 형식 위배 시 ackCode = "" 반환. ackCode 가 "" 이면 호출 측이 송신 실패로 처리.
func ParseACK(raw string) (ackCode, textMsg string) {
	msg, err := hl7.Parse(raw)
	if err != nil {
		return "", ""
	}
	msa := msg.SegmentByName("MSA")
	if msa == nil {
		return "", ""
	}
	ackCode = msa.Field(1).Value()
	textMsg = msa.Field(3).Value()
	return ackCode, textMsg
}
