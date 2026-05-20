package stt_test

import (
	"context"
	"strings"
	"testing"

	"github.com/manpasik/backend/shared/medical/stt"
)

func newTestChunk(sessionID string) *stt.AudioChunk {
	return &stt.AudioChunk{
		Data:       make([]byte, 32_000), // ~1초 16kHz 16bit mono
		Format:     stt.FormatWAV,
		SampleRate: 16_000,
		Channels:   1,
		Locale:     "ko-KR",
		SessionID:  sessionID,
		Speaker:    "patient",
	}
}

func TestNoopAdapter_Transcribe(t *testing.T) {
	a := stt.NewNoopAdapter()
	chunk := newTestChunk("s-1")

	transcript, err := a.Transcribe(context.Background(), chunk)
	if err != nil {
		t.Fatalf("Transcribe 실패: %v", err)
	}
	if transcript.Provider != "noop" {
		t.Errorf("Provider = %q", transcript.Provider)
	}
	if transcript.SessionID != "s-1" {
		t.Errorf("SessionID = %q", transcript.SessionID)
	}
	if transcript.FullText == "" {
		t.Error("FullText 비어 있음")
	}
	if len(transcript.Segments) == 0 {
		t.Error("Segments 비어 있음")
	}
	if a.ChunkCount() != 1 {
		t.Errorf("ChunkCount = %d", a.ChunkCount())
	}
}

func TestValidateChunk(t *testing.T) {
	cases := []struct {
		name    string
		c       *stt.AudioChunk
		wantErr bool
	}{
		{"nil", nil, true},
		{"empty data", &stt.AudioChunk{Format: stt.FormatWAV, SampleRate: 16000, Channels: 1, Locale: "ko"}, true},
		{"no format", &stt.AudioChunk{Data: []byte{1}, SampleRate: 16000, Channels: 1, Locale: "ko"}, true},
		{"zero sample rate", &stt.AudioChunk{Data: []byte{1}, Format: "wav", Channels: 1, Locale: "ko"}, true},
		{"3 channels", &stt.AudioChunk{Data: []byte{1}, Format: "wav", SampleRate: 16000, Channels: 3, Locale: "ko"}, true},
		{"no locale", &stt.AudioChunk{Data: []byte{1}, Format: "wav", SampleRate: 16000, Channels: 1}, true},
		{"valid", newTestChunk("s"), false},
	}
	for _, c := range cases {
		err := stt.ValidateChunk(c.c)
		if (err != nil) != c.wantErr {
			t.Errorf("%s: err=%v, wantErr=%v", c.name, err, c.wantErr)
		}
	}
}

func TestNoopAdapter_CustomTextGenerator(t *testing.T) {
	a := stt.NewNoopAdapter()
	a.SetTextGenerator(func(c *stt.AudioChunk) string {
		return "사용자 정의 응답: " + c.Speaker
	})

	chunk := newTestChunk("s")
	transcript, _ := a.Transcribe(context.Background(), chunk)
	if !strings.Contains(transcript.FullText, "사용자 정의") {
		t.Errorf("FullText = %q", transcript.FullText)
	}
}

func TestNoopAdapter_DurationCalculation(t *testing.T) {
	a := stt.NewNoopAdapter()
	chunk := newTestChunk("s")
	// 16kHz mono 16-bit → 1초 = 32000 bytes
	chunk.Data = make([]byte, 32_000)

	transcript, _ := a.Transcribe(context.Background(), chunk)
	if transcript.DurationSec < 0.9 || transcript.DurationSec > 1.1 {
		t.Errorf("DurationSec = %f, want ~1.0", transcript.DurationSec)
	}
}

func TestNoopAdapter_ContextCancelled(t *testing.T) {
	a := stt.NewNoopAdapter()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := a.Transcribe(ctx, newTestChunk("s"))
	if err == nil {
		t.Error("취소 ctx 통과")
	}
}

func TestWhisperAdapter_HealthCheck(t *testing.T) {
	a := stt.NewWhisperAdapter("api-key", "")
	if err := a.HealthCheck(context.Background()); err != nil {
		t.Errorf("HealthCheck 실패: %v", err)
	}

	a2 := stt.NewWhisperAdapter("", "")
	if err := a2.HealthCheck(context.Background()); err == nil {
		t.Error("API 키 없이 통과")
	}
}

func TestWhisperAdapter_Transcribe(t *testing.T) {
	a := stt.NewWhisperAdapter("test-key", "")
	transcript, err := a.Transcribe(context.Background(), newTestChunk("s"))
	if err != nil {
		t.Fatalf("Transcribe 실패: %v", err)
	}
	if transcript.Provider != "whisper" {
		t.Errorf("Provider = %q", transcript.Provider)
	}
}

func TestGoogleSTTAdapter_HealthCheck(t *testing.T) {
	a := stt.NewGoogleSTTAdapter("/path/to/creds.json")
	if err := a.HealthCheck(context.Background()); err != nil {
		t.Errorf("HealthCheck 실패: %v", err)
	}

	a2 := stt.NewGoogleSTTAdapter("")
	if err := a2.HealthCheck(context.Background()); err == nil {
		t.Error("credentials 없이 통과")
	}
}

func TestNewFromEnv_Default(t *testing.T) {
	t.Setenv("STT_PROVIDER", "")
	a := stt.NewFromEnv()
	if a.Provider() != "noop" {
		t.Errorf("Provider = %q", a.Provider())
	}
}

func TestNewFromEnv_Whisper(t *testing.T) {
	t.Setenv("STT_PROVIDER", "whisper")
	a := stt.NewFromEnv()
	if a.Provider() != "whisper" {
		t.Errorf("Provider = %q", a.Provider())
	}
}

func TestNewFromEnv_Google(t *testing.T) {
	t.Setenv("STT_PROVIDER", "google")
	a := stt.NewFromEnv()
	if a.Provider() != "google_stt" {
		t.Errorf("Provider = %q", a.Provider())
	}
}

func TestStreamingAggregator_MultiChunk(t *testing.T) {
	adapter := stt.NewNoopAdapter()
	agg := stt.NewStreamingAggregator(adapter)

	for i := 0; i < 3; i++ {
		chunk := newTestChunk("session-stream")
		_, err := agg.AppendChunk(context.Background(), chunk)
		if err != nil {
			t.Fatalf("AppendChunk %d 실패: %v", i, err)
		}
	}

	final := agg.FinalizeSession("session-stream")
	if final == nil {
		t.Fatal("FinalizeSession nil")
	}
	if len(final.Segments) != 3 {
		t.Errorf("Segments = %d, want 3", len(final.Segments))
	}
	// 시간 오프셋 누적 확인
	if final.Segments[2].StartTime <= final.Segments[1].StartTime {
		t.Error("시간 오프셋 누적 안됨")
	}
}

func TestStreamingAggregator_FinalizeEmpty(t *testing.T) {
	adapter := stt.NewNoopAdapter()
	agg := stt.NewStreamingAggregator(adapter)

	if final := agg.FinalizeSession("missing"); final != nil {
		t.Error("빈 세션이 nil이 아님")
	}
}

func TestStreamingAggregator_ClearSession(t *testing.T) {
	adapter := stt.NewNoopAdapter()
	agg := stt.NewStreamingAggregator(adapter)

	_, _ = agg.AppendChunk(context.Background(), newTestChunk("s"))
	agg.ClearSession("s")

	if agg.FinalizeSession("s") != nil {
		t.Error("Clear 후에도 세션 존재")
	}
}

func TestNoopAdapter_Clear(t *testing.T) {
	a := stt.NewNoopAdapter()
	for i := 0; i < 3; i++ {
		_, _ = a.Transcribe(context.Background(), newTestChunk("s"))
	}
	a.Clear()
	if a.ChunkCount() != 0 {
		t.Errorf("Clear 후 ChunkCount = %d", a.ChunkCount())
	}
}

func TestSegment_Confidence(t *testing.T) {
	a := stt.NewNoopAdapter()
	transcript, _ := a.Transcribe(context.Background(), newTestChunk("s"))
	if transcript.Segments[0].Confidence < 0.5 {
		t.Errorf("기본 Confidence = %f", transcript.Segments[0].Confidence)
	}
}

func TestSpeakerPropagation(t *testing.T) {
	a := stt.NewNoopAdapter()
	chunk := newTestChunk("s")
	chunk.Speaker = "doctor"

	transcript, _ := a.Transcribe(context.Background(), chunk)
	if transcript.Segments[0].Speaker != "doctor" {
		t.Errorf("Speaker = %q, want doctor", transcript.Segments[0].Speaker)
	}
}
