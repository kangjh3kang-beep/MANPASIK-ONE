package digital_physician_test

import (
	"context"
	"strings"
	"testing"

	dp "github.com/manpasik/backend/shared/medical/digital_physician"
)

func TestNoopAdapter_StartSession(t *testing.T) {
	a := dp.NewNoopAdapter()
	ctx := context.Background()

	if err := a.StartSession(ctx, "session-1"); err != nil {
		t.Fatalf("StartSession 실패: %v", err)
	}
	if a.SessionCount() != 1 {
		t.Errorf("SessionCount = %d, want 1", a.SessionCount())
	}
}

func TestNoopAdapter_StartSession_Duplicate(t *testing.T) {
	a := dp.NewNoopAdapter()
	_ = a.StartSession(context.Background(), "s")
	if err := a.StartSession(context.Background(), "s"); err == nil {
		t.Error("중복 세션 통과")
	}
}

func TestNoopAdapter_RequiresSessionID(t *testing.T) {
	a := dp.NewNoopAdapter()
	if err := a.StartSession(context.Background(), ""); err == nil {
		t.Error("빈 session_id 통과")
	}
}

func TestNoopAdapter_Speak(t *testing.T) {
	a := dp.NewNoopAdapter()
	ctx := context.Background()
	_ = a.StartSession(ctx, "s1")

	result, err := a.Speak(ctx, "s1", &dp.Utterance{
		Text:    "안녕하세요. 어디가 불편하신가요?",
		Emotion: dp.EmotionWarmth,
		Gesture: dp.GestureNod,
	})
	if err != nil {
		t.Fatalf("Speak 실패: %v", err)
	}
	if result.Status != "completed" {
		t.Errorf("Status = %q", result.Status)
	}
	if result.DurationMs <= 0 {
		t.Error("DurationMs <= 0")
	}

	// 발화 후 상태 확인
	state, _ := a.GetState(ctx, "s1")
	if state.Emotion != dp.EmotionWarmth {
		t.Errorf("Emotion = %q, want warmth", state.Emotion)
	}
	if state.IsSpeaking {
		t.Error("발화 완료 후에도 IsSpeaking=true")
	}
}

func TestNoopAdapter_Speak_NoSession(t *testing.T) {
	a := dp.NewNoopAdapter()
	_, err := a.Speak(context.Background(), "missing", &dp.Utterance{Text: "x"})
	if err == nil {
		t.Error("미존재 세션에 발화 통과")
	}
}

func TestNoopAdapter_SetEmotion(t *testing.T) {
	a := dp.NewNoopAdapter()
	ctx := context.Background()
	_ = a.StartSession(ctx, "s")

	if err := a.SetEmotion(ctx, "s", dp.EmotionConcern); err != nil {
		t.Fatalf("SetEmotion 실패: %v", err)
	}
	state, _ := a.GetState(ctx, "s")
	if state.Emotion != dp.EmotionConcern {
		t.Errorf("Emotion = %q", state.Emotion)
	}
}

func TestNoopAdapter_SetGesture(t *testing.T) {
	a := dp.NewNoopAdapter()
	ctx := context.Background()
	_ = a.StartSession(ctx, "s")

	if err := a.SetGesture(ctx, "s", dp.GesturePointing); err != nil {
		t.Fatalf("SetGesture 실패: %v", err)
	}
	state, _ := a.GetState(ctx, "s")
	if state.Gesture != dp.GesturePointing {
		t.Errorf("Gesture = %q", state.Gesture)
	}
}

func TestNoopAdapter_SetProfile(t *testing.T) {
	a := dp.NewNoopAdapter()
	ctx := context.Background()
	_ = a.StartSession(ctx, "s")

	profile := &dp.VoiceProfile{
		Locale: "en-US", Gender: "male", AgeRange: "adult",
		Pitch: 0.9, Speed: 1.1, Style: "professional",
	}
	if err := a.SetProfile(ctx, "s", profile); err != nil {
		t.Fatalf("SetProfile 실패: %v", err)
	}
	state, _ := a.GetState(ctx, "s")
	if state.CurrentLocale != "en-US" {
		t.Errorf("Locale = %q", state.CurrentLocale)
	}
}

func TestNoopAdapter_EndSession(t *testing.T) {
	a := dp.NewNoopAdapter()
	ctx := context.Background()
	_ = a.StartSession(ctx, "s")
	_ = a.EndSession(ctx, "s")

	if a.SessionCount() != 0 {
		t.Errorf("SessionCount = %d, want 0", a.SessionCount())
	}
}

func TestNoopAdapter_History(t *testing.T) {
	a := dp.NewNoopAdapter()
	ctx := context.Background()
	_ = a.StartSession(ctx, "s")

	for i := 0; i < 3; i++ {
		_, _ = a.Speak(ctx, "s", &dp.Utterance{Text: "발화 " + string(rune('a'+i))})
	}

	history := a.History("s")
	if len(history) != 3 {
		t.Errorf("History = %d, want 3", len(history))
	}
}

func TestValidateUtterance(t *testing.T) {
	cases := []struct {
		name    string
		utt     *dp.Utterance
		wantErr bool
	}{
		{"nil", nil, true},
		{"empty text", &dp.Utterance{Text: ""}, true},
		{"valid", &dp.Utterance{Text: "안녕하세요"}, false},
		{"too long", &dp.Utterance{Text: strings.Repeat("a", 2001)}, true},
		{"unsupported locale", &dp.Utterance{Text: "x", Locale: "fr-FR"}, true},
		{"valid locale", &dp.Utterance{Text: "x", Locale: "ja-JP"}, false},
		{"negative duration", &dp.Utterance{Text: "x", MaxDurationSec: -1}, true},
	}
	for _, c := range cases {
		err := dp.ValidateUtterance(c.utt)
		if (err != nil) != c.wantErr {
			t.Errorf("%s: err=%v, wantErr=%v", c.name, err, c.wantErr)
		}
	}
}

func TestValidateProfile(t *testing.T) {
	cases := []struct {
		name    string
		p       *dp.VoiceProfile
		wantErr bool
	}{
		{"nil", nil, true},
		{"valid", &dp.VoiceProfile{Locale: "ko-KR", Pitch: 1.0, Speed: 1.0}, false},
		{"unsupported locale", &dp.VoiceProfile{Locale: "fr", Pitch: 1.0, Speed: 1.0}, true},
		{"pitch too low", &dp.VoiceProfile{Locale: "ko-KR", Pitch: 0.3, Speed: 1.0}, true},
		{"pitch too high", &dp.VoiceProfile{Locale: "ko-KR", Pitch: 2.0, Speed: 1.0}, true},
		{"speed too high", &dp.VoiceProfile{Locale: "ko-KR", Pitch: 1.0, Speed: 3.0}, true},
	}
	for _, c := range cases {
		err := dp.ValidateProfile(c.p)
		if (err != nil) != c.wantErr {
			t.Errorf("%s: err=%v, wantErr=%v", c.name, err, c.wantErr)
		}
	}
}

func TestACEAdapter_HealthCheck(t *testing.T) {
	a := dp.NewACEAdapter("https://ace.nvidia.com", "key")
	if err := a.HealthCheck(context.Background()); err != nil {
		t.Errorf("HealthCheck 실패: %v", err)
	}

	a2 := dp.NewACEAdapter("", "")
	if err := a2.HealthCheck(context.Background()); err == nil {
		t.Error("자격증명 없이 통과")
	}
}

func TestACEAdapter_DelegatesToNoop(t *testing.T) {
	a := dp.NewACEAdapter("https://x", "k")
	ctx := context.Background()

	if err := a.StartSession(ctx, "s"); err != nil {
		t.Fatalf("StartSession 실패: %v", err)
	}
	_, err := a.Speak(ctx, "s", &dp.Utterance{Text: "안녕"})
	if err != nil {
		t.Fatalf("Speak 실패: %v", err)
	}
	if a.Provider() != "nvidia_ace" {
		t.Errorf("Provider = %q", a.Provider())
	}
}

func TestNewFromEnv_Default(t *testing.T) {
	t.Setenv("DIGITAL_PHYSICIAN_PROVIDER", "")
	a := dp.NewFromEnv()
	if a.Provider() != "noop" {
		t.Errorf("Provider = %q, want noop", a.Provider())
	}
}

func TestNewFromEnv_ACE(t *testing.T) {
	t.Setenv("DIGITAL_PHYSICIAN_PROVIDER", "ace")
	t.Setenv("NVIDIA_ACE_ENDPOINT", "https://x")
	t.Setenv("NVIDIA_ACE_API_KEY", "k")

	a := dp.NewFromEnv()
	if a.Provider() != "nvidia_ace" {
		t.Errorf("Provider = %q", a.Provider())
	}
}

func TestEmotionTransitions(t *testing.T) {
	a := dp.NewNoopAdapter()
	ctx := context.Background()
	_ = a.StartSession(ctx, "s")

	emotions := []dp.Emotion{
		dp.EmotionNeutral, dp.EmotionWarmth,
		dp.EmotionConcern, dp.EmotionAttention,
		dp.EmotionEncouragement,
	}
	for _, e := range emotions {
		if err := a.SetEmotion(ctx, "s", e); err != nil {
			t.Errorf("SetEmotion(%s) 실패: %v", e, err)
		}
	}
	state, _ := a.GetState(ctx, "s")
	if state.Emotion != dp.EmotionEncouragement {
		t.Errorf("최종 Emotion = %q", state.Emotion)
	}
}

func TestGestureTransitions(t *testing.T) {
	a := dp.NewNoopAdapter()
	ctx := context.Background()
	_ = a.StartSession(ctx, "s")

	gestures := []dp.Gesture{
		dp.GestureNod, dp.GestureHeadTilt,
		dp.GestureHandsUp, dp.GesturePointing, dp.GestureNone,
	}
	for _, g := range gestures {
		if err := a.SetGesture(ctx, "s", g); err != nil {
			t.Errorf("SetGesture(%s) 실패: %v", g, err)
		}
	}
}

func TestSpeakDurationScalesWithText(t *testing.T) {
	a := dp.NewNoopAdapter()
	ctx := context.Background()
	_ = a.StartSession(ctx, "s")

	short, _ := a.Speak(ctx, "s", &dp.Utterance{Text: "안녕"})
	long, _ := a.Speak(ctx, "s", &dp.Utterance{
		Text: strings.Repeat("긴 발화입니다. ", 50),
	})

	if long.DurationMs <= short.DurationMs {
		t.Errorf("긴 발화(%dms) <= 짧은 발화(%dms)", long.DurationMs, short.DurationMs)
	}
}

func TestDefaultProfile(t *testing.T) {
	if dp.DefaultProfile.Locale != "ko-KR" {
		t.Errorf("기본 Locale = %q, want ko-KR", dp.DefaultProfile.Locale)
	}
	if err := dp.ValidateProfile(&dp.DefaultProfile); err != nil {
		t.Errorf("DefaultProfile 검증 실패: %v", err)
	}
}
