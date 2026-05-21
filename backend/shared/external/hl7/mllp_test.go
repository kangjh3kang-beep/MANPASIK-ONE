package hl7mllp_test

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	hl7 "github.com/manpasik/backend/shared/hl7-parser"
	mllp "github.com/manpasik/backend/shared/external/hl7"
)

const testORU = "MSH|^~\\&|LIS|HospA|EMR|HospA|20260521103000||ORU^R01|MSGCLI|P|2.5\r" +
	"PID|1||PAT1||김^영수||19880101|M\r" +
	"OBR|1||LAB-1|GLU^Glucose^L\r" +
	"OBX|1|NM|GLU^Glucose^LN||95|mg/dL|70-110|N|||F\r"

// =============================================================================
// Framing 테스트
// =============================================================================

func TestEncodeFrame_StartEndBytes(t *testing.T) {
	frame := mllp.EncodeFrame([]byte("hello"))
	if frame[0] != 0x0B {
		t.Errorf("start byte = 0x%02X, 0x0B 기대", frame[0])
	}
	if frame[len(frame)-2] != 0x1C || frame[len(frame)-1] != 0x0D {
		t.Errorf("end bytes = 0x%02X 0x%02X, 0x1C 0x0D 기대",
			frame[len(frame)-2], frame[len(frame)-1])
	}
}

func TestDecodeFrame_RoundTrip(t *testing.T) {
	payload := []byte("MSH|^~\\&|...")
	frame := mllp.EncodeFrame(payload)
	decoded, err := mllp.DecodeFrame(frame)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(decoded, payload) {
		t.Errorf("decoded != original: %q vs %q", decoded, payload)
	}
}

func TestDecodeFrame_InvalidFrames(t *testing.T) {
	cases := [][]byte{
		nil,
		{0x0B},
		{0x0B, 'A', 0x0D}, // missing FS
		{'A', 0x1C, 0x0D}, // missing VT
	}
	for i, c := range cases {
		if _, err := mllp.DecodeFrame(c); err == nil {
			t.Errorf("case %d: invalid frame 통과: %v", i, c)
		}
	}
}

func TestReadFrame_IgnoresLeadingNoise(t *testing.T) {
	// VT 이전의 잡음은 무시되어야 함
	buf := bytes.NewBuffer(append([]byte("junk-junk"), mllp.EncodeFrame([]byte("MSH..."))...))
	r := bufio.NewReader(buf)
	payload, err := mllp.ReadFrame(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(payload) != "MSH..." {
		t.Errorf("payload = %q", payload)
	}
}

func TestReadFrame_EOF(t *testing.T) {
	r := bufio.NewReader(bytes.NewBuffer(nil))
	_, err := mllp.ReadFrame(r)
	if !errors.Is(err, io.EOF) {
		t.Errorf("err = %v, EOF 기대", err)
	}
}

// =============================================================================
// ACK 빌더 테스트
// =============================================================================

func TestBuildACK_FromValidORU(t *testing.T) {
	msg, _ := hl7.Parse(testORU)
	ack := mllp.BuildACK(msg, "AA", "ok")
	if !strings.HasPrefix(ack, "MSH|") {
		t.Errorf("ACK MSH 누락:\n%s", ack)
	}
	if !strings.Contains(ack, "MSA|AA|MSGCLI|") {
		t.Errorf("MSA 필드 누락:\n%s", ack)
	}
	if !strings.Contains(ack, "ACK^R01") {
		t.Errorf("ACK event code (R01 상속) 누락:\n%s", ack)
	}
}

func TestBuildACK_NilMessage(t *testing.T) {
	ack := mllp.BuildACK(nil, "AR", "parse failed")
	if !strings.Contains(ack, "MSA|AR||") {
		t.Errorf("nil 메시지 ACK 형식:\n%s", ack)
	}
	if !strings.Contains(ack, "MANPASIK") {
		t.Errorf("기본 송신자 MANPASIK 누락")
	}
}

func TestBuildACK_SanitizesText(t *testing.T) {
	// 위험한 HL7 구분자가 포함된 텍스트
	ack := mllp.BuildACK(nil, "AE", "err: |broken^msg~field|")
	if strings.Contains(ack, "MSA|AE||err: |broken") {
		t.Errorf("sanitize 실패 — | 가 남아있음:\n%s", ack)
	}
}

func TestParseACK_RoundTrip(t *testing.T) {
	msg, _ := hl7.Parse(testORU)
	ack := mllp.BuildACK(msg, "AA", "okay")
	code, text := mllp.ParseACK(ack)
	if code != "AA" {
		t.Errorf("code = %s", code)
	}
	if text != "okay" {
		t.Errorf("text = %s", text)
	}
}

// =============================================================================
// Server/Client 통합 테스트
// =============================================================================

func TestMLLPServerClient_HappyPath(t *testing.T) {
	var receivedMessages int32
	server, err := mllp.NewMLLPServer(mllp.ServerConfig{
		Addr: "127.0.0.1:0",
		Handler: func(_ context.Context, raw []byte, msg *hl7.Message) (string, string, error) {
			atomic.AddInt32(&receivedMessages, 1)
			if msg == nil || msg.MessageType() != "ORU" {
				return "AR", "expected ORU", nil
			}
			return "AA", "stored", nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := server.Start(); err != nil {
		t.Fatal(err)
	}
	defer server.Stop()

	// 클라이언트로 송신
	client, err := mllp.NewClient(mllp.ClientConfig{
		Addr:         server.Addr(),
		DialTimeout:  1 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 1 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	code, text, err := client.SendString(testORU)
	if err != nil {
		t.Fatal(err)
	}
	if code != "AA" {
		t.Errorf("code = %s, AA 기대", code)
	}
	if text != "stored" {
		t.Errorf("text = %s", text)
	}
	if atomic.LoadInt32(&receivedMessages) != 1 {
		t.Errorf("server 수신 = %d, 1 기대", receivedMessages)
	}

	// 통계 확인
	stats := server.Stats()
	if stats.MessagesAccepted != 1 {
		t.Errorf("Accepted = %d", stats.MessagesAccepted)
	}
}

func TestMLLPServer_InvalidMessageRejected(t *testing.T) {
	server, _ := mllp.NewMLLPServer(mllp.ServerConfig{
		Addr:    "127.0.0.1:0",
		Handler: func(_ context.Context, _ []byte, _ *hl7.Message) (string, string, error) { return "AA", "", nil },
	})
	_ = server.Start()
	defer server.Stop()

	// raw TCP — invalid HL7 송신
	conn, err := net.Dial("tcp", server.Addr())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	if err := mllp.WriteFrame(conn, []byte("NOT_A_VALID_HL7")); err != nil {
		t.Fatal(err)
	}
	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	reader := bufio.NewReader(conn)
	ackRaw, err := mllp.ReadFrame(reader)
	if err != nil {
		t.Fatal(err)
	}
	code, _ := mllp.ParseACK(string(ackRaw))
	if code != "AR" {
		t.Errorf("ACK code = %s, AR 기대 (파싱 실패)", code)
	}
}

func TestMLLPServer_MultipleMessagesSameConnection(t *testing.T) {
	var processed int32
	server, _ := mllp.NewMLLPServer(mllp.ServerConfig{
		Addr: "127.0.0.1:0",
		Handler: func(_ context.Context, _ []byte, _ *hl7.Message) (string, string, error) {
			atomic.AddInt32(&processed, 1)
			return "AA", "", nil
		},
	})
	_ = server.Start()
	defer server.Stop()

	client, _ := mllp.NewClient(mllp.ClientConfig{Addr: server.Addr()})
	defer client.Close()

	for i := 0; i < 3; i++ {
		code, _, err := client.SendString(testORU)
		if err != nil {
			t.Fatalf("msg %d: %v", i, err)
		}
		if code != "AA" {
			t.Errorf("msg %d code = %s", i, code)
		}
	}
	if atomic.LoadInt32(&processed) != 3 {
		t.Errorf("processed = %d, 3 기대", processed)
	}
}

func TestMLLPServer_HandlerErrorCode(t *testing.T) {
	server, _ := mllp.NewMLLPServer(mllp.ServerConfig{
		Addr: "127.0.0.1:0",
		Handler: func(_ context.Context, _ []byte, _ *hl7.Message) (string, string, error) {
			return "AE", "downstream failed", errors.New("DB unavailable")
		},
	})
	_ = server.Start()
	defer server.Stop()

	client, _ := mllp.NewClient(mllp.ClientConfig{Addr: server.Addr()})
	defer client.Close()

	code, text, err := client.SendString(testORU)
	if err != nil {
		t.Fatal(err)
	}
	if code != "AE" {
		t.Errorf("code = %s, AE 기대", code)
	}
	if text != "downstream failed" {
		t.Errorf("text = %s", text)
	}
	if server.Stats().MessagesError != 1 {
		t.Errorf("Error count = %d", server.Stats().MessagesError)
	}
}

func TestMLLPClient_NoServer(t *testing.T) {
	client, _ := mllp.NewClient(mllp.ClientConfig{
		Addr:        "127.0.0.1:1", // 닫힌 포트
		DialTimeout: 200 * time.Millisecond,
	})
	_, _, err := client.SendString(testORU)
	if err == nil {
		t.Error("연결 실패가 통과")
	}
}

func TestMLLPServer_StartStopIdempotent(t *testing.T) {
	server, _ := mllp.NewMLLPServer(mllp.ServerConfig{
		Addr:    "127.0.0.1:0",
		Handler: func(_ context.Context, _ []byte, _ *hl7.Message) (string, string, error) { return "AA", "", nil },
	})
	_ = server.Start()
	// 두 번째 Start 는 에러
	if err := server.Start(); err == nil {
		t.Error("이중 Start 통과")
	}
	_ = server.Stop()
	// Stop 후 다시 Stop 은 noop
	_ = server.Stop()
}

func TestMLLPClient_AddrRequired(t *testing.T) {
	if _, err := mllp.NewClient(mllp.ClientConfig{}); err == nil {
		t.Error("Addr 없이 통과")
	}
}

func TestMLLPServer_HandlerRequired(t *testing.T) {
	if _, err := mllp.NewMLLPServer(mllp.ServerConfig{Addr: "x"}); err == nil {
		t.Error("Handler 없이 통과")
	}
}
