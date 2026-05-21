// Package hl7mllp 는 HL7 v2 MLLP (Minimum Lower Layer Protocol) 어댑터를
// 제공합니다 (Phase AX).
//
// MLLP 는 HL7 v2 메시지의 표준 TCP 전송 프로토콜로, 의료기관 LIS/EHR 간
// 통신에서 가장 널리 사용됩니다. 만파식은 외부 시스템이 HTTP 가 아닌 TCP
// 전용일 때 이 어댑터로 통합합니다 (Phase AU 의 HTTP 와 병행).
//
// 프레임 형식:
//
//	<VT> message <FS><CR>
//	0x0B  payload  0x1C 0x0D
//
// 한 TCP 연결로 다중 메시지 송수신 가능. 송신자는 각 메시지마다 수신자의
// ACK 응답을 기다림.
package hl7mllp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
)

// MLLP 프레임 구분자 (HL7 표준 — 변경 불가).
const (
	startBlock = 0x0B // <VT>
	endBlock1  = 0x1C // <FS>
	endBlock2  = 0x0D // <CR>
)

// maxMessageBytes 는 단일 MLLP 메시지의 최대 크기 (안전 마진).
const maxMessageBytes = 1 << 20 // 1 MB

// 에러 sentinel.
var (
	// ErrInvalidFrame 은 MLLP 프레임이 정해진 구분자로 시작/끝나지 않음.
	ErrInvalidFrame = errors.New("MLLP 프레임 형식 오류")

	// ErrMessageTooLarge 는 1MB 상한 초과.
	ErrMessageTooLarge = errors.New("MLLP 메시지 크기 초과")
)

// EncodeFrame 은 payload 를 MLLP 프레임으로 감쌉니다.
//
//	output = <VT> payload <FS><CR>
func EncodeFrame(payload []byte) []byte {
	out := make([]byte, 0, len(payload)+3)
	out = append(out, startBlock)
	out = append(out, payload...)
	out = append(out, endBlock1, endBlock2)
	return out
}

// DecodeFrame 은 MLLP 프레임에서 payload 를 추출합니다.
//
// 프레임 구분자가 누락된 경우 ErrInvalidFrame.
func DecodeFrame(frame []byte) ([]byte, error) {
	if len(frame) < 3 {
		return nil, ErrInvalidFrame
	}
	if frame[0] != startBlock {
		return nil, ErrInvalidFrame
	}
	// 끝 구분자 검증 (마지막 2바이트)
	last := len(frame) - 1
	if frame[last] != endBlock2 || frame[last-1] != endBlock1 {
		return nil, ErrInvalidFrame
	}
	return frame[1:last-1], nil
}

// ReadFrame 은 reader 에서 한 프레임을 읽어 payload 반환.
//
// 한 TCP 연결로 여러 메시지를 받을 때 반복 호출. EOF / 연결 종료 시 io.EOF
// 또는 다른 에러 반환.
func ReadFrame(r *bufio.Reader) ([]byte, error) {
	// startBlock 까지 무시 (HL7 표준 — 프레임 사이 잡음 허용).
	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		if b == startBlock {
			break
		}
	}

	// payload 누적 (endBlock1 만나면 종료)
	var buf bytes.Buffer
	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		if b == endBlock1 {
			// 다음 바이트는 반드시 endBlock2
			next, err := r.ReadByte()
			if err != nil {
				return nil, err
			}
			if next != endBlock2 {
				return nil, ErrInvalidFrame
			}
			return buf.Bytes(), nil
		}
		if buf.Len() >= maxMessageBytes {
			return nil, ErrMessageTooLarge
		}
		buf.WriteByte(b)
	}
}

// WriteFrame 은 payload 를 framing 한 후 writer 에 송신.
//
// 자동으로 flush 까지 수행 — 호출 측이 bufio.Writer 를 쓰더라도 즉시 전송.
func WriteFrame(w io.Writer, payload []byte) error {
	frame := EncodeFrame(payload)
	if _, err := w.Write(frame); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	if bw, ok := w.(*bufio.Writer); ok {
		return bw.Flush()
	}
	return nil
}
