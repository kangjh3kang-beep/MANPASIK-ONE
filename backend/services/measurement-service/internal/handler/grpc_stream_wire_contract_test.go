package handler

import (
	"context"
	"encoding/hex"
	"math"
	"testing"
	"time"

	"github.com/manpasik/backend/services/measurement-service/internal/service"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"google.golang.org/protobuf/proto"
)

const dartMeasurementDataGoldenHex = "" +
	"0a0973657373696f6e2d31" +
	"1218000000000000f03f00000000000000400000000000000840" +
	"1a24" +
	"090000000000005940" +
	"110000000000001440" +
	"19666666666666ee3f" +
	"210000000000d05740" +
	"220f" +
	"0d0000c441" +
	"1500003442" +
	"1d9a99ca42"

func TestStreamMeasurementAcceptsDartWireGoldenFrame(t *testing.T) {
	payload, err := hex.DecodeString(dartMeasurementDataGoldenHex)
	if err != nil {
		t.Fatalf("decode golden hex: %v", err)
	}

	var frame v1.MeasurementData
	if err := proto.Unmarshal(payload, &frame); err != nil {
		t.Fatalf("unmarshal Dart golden MeasurementData: %v", err)
	}

	if frame.SessionId != "session-1" {
		t.Fatalf("SessionId = %q, want session-1", frame.SessionId)
	}
	assertFloat64Slice(t, frame.RawChannels, []float64{1, 2, 3})
	if frame.GetDifferential().GetSDet() != 100 ||
		frame.GetDifferential().GetSRef() != 5 ||
		frame.GetDifferential().GetAlpha() != 0.95 ||
		frame.GetDifferential().GetSCorrected() != 95.25 {
		t.Fatalf("differential mismatch: %+v", frame.GetDifferential())
	}
	if frame.GetEnvMeta().GetTempC() != 24.5 ||
		frame.GetEnvMeta().GetHumidityPct() != 45 ||
		math.Abs(float64(frame.GetEnvMeta().GetPressureKpa()-101.3)) > 0.001 {
		t.Fatalf("env_meta mismatch: %+v", frame.GetEnvMeta())
	}

	handler, measureRepo, vectorRepo := newStreamMeasurementHandler(&service.MeasurementSession{
		ID:          "session-1",
		DeviceID:    "device-1",
		CartridgeID: "glucose",
		UserID:      "user-1",
		StartedAt:   time.Now().UTC(),
		Status:      "active",
	})
	stream := &measurementStreamFake{
		ctx:      context.Background(),
		requests: []*v1.MeasurementData{&frame},
	}

	if err := handler.StreamMeasurement(stream); err != nil {
		t.Fatalf("StreamMeasurement with Dart golden frame failed: %v", err)
	}
	if len(measureRepo.stored) != 1 {
		t.Fatalf("stored measurement count = %d, want 1", len(measureRepo.stored))
	}
	stored := measureRepo.stored[0]
	if stored.SessionID != "session-1" || stored.DeviceID != "device-1" ||
		stored.UserID != "user-1" || stored.CartridgeType != "glucose" {
		t.Fatalf("stored metadata mismatch: %+v", stored)
	}
	if stored.PrimaryValue != 95.25 || stored.SCorrected != 95.25 {
		t.Fatalf("stored corrected values mismatch: primary=%f corrected=%f", stored.PrimaryValue, stored.SCorrected)
	}
	assertFloat32Slice(t, stored.FingerprintVector, []float32{1, 2, 3})
	assertFloat32Slice(t, vectorRepo.vectors["session-1"], []float32{1, 2, 3})
	if len(stream.responses) != 1 {
		t.Fatalf("stream response count = %d, want 1", len(stream.responses))
	}
	if stream.responses[0].PrimaryValue != 95.25 || stream.responses[0].SessionId != "session-1" {
		t.Fatalf("response mismatch: %+v", stream.responses[0])
	}
}
