package observability_test

import (
	"context"
	"errors"
	"testing"

	"github.com/manpasik/backend/shared/medical/observability"
	"github.com/manpasik/backend/shared/medical/sla"
	"github.com/manpasik/backend/shared/ops/tracing"
)

func newTestInstrumentor() *observability.MedicalInstrumentor {
	tracer := tracing.NewMemoryTracer("test")
	slaMon := sla.NewMonitor(nil)
	return observability.NewMedicalInstrumentor(tracer, slaMon, "telemedicine-test")
}

func TestStartOperation_HappyPath(t *testing.T) {
	inst := newTestInstrumentor()
	ctx := context.Background()

	op := inst.StartOperation(ctx, "test-op", sla.TierSTAT, "session-1", "2345-7")
	if op == nil {
		t.Fatal("StartOperation returned nil")
	}

	op.SetAttribute("custom_key", "custom_value")
	op.AddEvent("checkpoint", map[string]string{"step": "1"})
	measurement := op.End(nil)

	if measurement == nil {
		t.Error("End는 SLA Measurement 반환해야 함")
	}
	if measurement.SessionID != "session-1" {
		t.Errorf("SessionID = %q", measurement.SessionID)
	}
}

func TestStartOperation_ErrorMarksSpan(t *testing.T) {
	inst := newTestInstrumentor()
	op := inst.StartOperation(context.Background(), "fail-op", sla.TierSTAT, "s", "x")

	op.End(errors.New("simulated failure"))

	report := inst.GenerateReport("s", "p")
	if report.ErrorSpans < 1 {
		t.Errorf("ErrorSpans = %d, want >= 1", report.ErrorSpans)
	}
}

func TestStartOperation_Context_ChildSpans(t *testing.T) {
	inst := newTestInstrumentor()
	parent := inst.StartOperation(context.Background(), "parent", sla.TierSTAT, "s", "x")

	// 부모 컨텍스트로 자식 작업 생성
	childCtx := parent.Context()
	parentSpan := tracing.SpanFromContext(childCtx)
	if parentSpan == nil {
		t.Error("Context에서 부모 span 추출 실패")
	}

	parent.End(nil)
}

func TestInstrumentedWorkflow_5Stages(t *testing.T) {
	inst := newTestInstrumentor()
	wf := observability.NewInstrumentedWorkflow(inst, "session-wf", "P-001")

	ctx := context.Background()
	stages := []func(context.Context) error{
		func(c context.Context) error {
			return wf.TrackPOCT(c, func(_ context.Context) error { return nil })
		},
		func(c context.Context) error {
			return wf.TrackFHIR(c, func(_ context.Context) error { return nil })
		},
		func(c context.Context) error {
			return wf.TrackReasoning(c, func(_ context.Context) error { return nil })
		},
		func(c context.Context) error {
			return wf.TrackSafety(c, func(_ context.Context) error { return nil })
		},
		func(c context.Context) error {
			return wf.TrackScribe(c, func(_ context.Context) error { return nil })
		},
	}

	for i, stage := range stages {
		if err := stage(ctx); err != nil {
			t.Errorf("stage %d 실패: %v", i, err)
		}
	}

	report := inst.GenerateReport("session-wf", "P-001")
	if report.TotalSpans != 5 {
		t.Errorf("TotalSpans = %d, want 5", report.TotalSpans)
	}
	if report.ErrorSpans != 0 {
		t.Errorf("ErrorSpans = %d, want 0", report.ErrorSpans)
	}
}

func TestInstrumentedWorkflow_ErrorPropagation(t *testing.T) {
	inst := newTestInstrumentor()
	wf := observability.NewInstrumentedWorkflow(inst, "session-err", "P")

	expectedErr := errors.New("FHIR conversion failed")
	got := wf.TrackFHIR(context.Background(), func(_ context.Context) error {
		return expectedErr
	})

	if got != expectedErr {
		t.Errorf("err = %v, want %v", got, expectedErr)
	}
}

func TestGenerateReport_FilterBySession(t *testing.T) {
	inst := newTestInstrumentor()
	wf1 := observability.NewInstrumentedWorkflow(inst, "session-A", "P")
	wf2 := observability.NewInstrumentedWorkflow(inst, "session-B", "P")

	_ = wf1.TrackPOCT(context.Background(), func(_ context.Context) error { return nil })
	_ = wf2.TrackPOCT(context.Background(), func(_ context.Context) error { return nil })

	reportA := inst.GenerateReport("session-A", "P")
	if reportA.TotalSpans != 1 {
		t.Errorf("session-A spans = %d, want 1", reportA.TotalSpans)
	}
}

func TestProfileLatency(t *testing.T) {
	inst := newTestInstrumentor()
	wf := observability.NewInstrumentedWorkflow(inst, "session-prof", "P")

	for i := 0; i < 3; i++ {
		_ = wf.TrackPOCT(context.Background(), func(_ context.Context) error { return nil })
	}

	profile := inst.ProfileLatency()
	stats, ok := profile.Operations["medical.poct.parse"]
	if !ok {
		t.Fatal("medical.poct.parse 작업 누락")
	}
	if stats.Count != 3 {
		t.Errorf("Count = %d, want 3", stats.Count)
	}
	if stats.MeanMs < 0 {
		t.Error("MeanMs 음수")
	}
}

func TestProfileLatency_String(t *testing.T) {
	inst := newTestInstrumentor()
	_ = observability.NewInstrumentedWorkflow(inst, "s", "p").
		TrackPOCT(context.Background(), func(_ context.Context) error { return nil })

	profile := inst.ProfileLatency()
	str := profile.String()
	if str == "" {
		t.Error("String empty")
	}
}

func TestProfileLatency_Empty(t *testing.T) {
	inst := newTestInstrumentor()
	profile := inst.ProfileLatency()
	str := profile.String()
	if str != "[empty profile]" {
		t.Errorf("String = %q", str)
	}
}

func TestNewInstrumentor_NilDefaults(t *testing.T) {
	// 두 의존성 모두 nil이어도 동작해야 함
	inst := observability.NewMedicalInstrumentor(nil, nil, "default-svc")
	if inst == nil {
		t.Fatal("nil 의존성 시 nil 반환")
	}

	op := inst.StartOperation(context.Background(), "x", sla.TierSTAT, "s", "")
	op.End(nil)
}

func TestStartOperation_AttributesAreRecorded(t *testing.T) {
	tracer := tracing.NewMemoryTracer("test")
	slaMon := sla.NewMonitor(nil)
	inst := observability.NewMedicalInstrumentor(tracer, slaMon, "svc")

	op := inst.StartOperation(context.Background(), "named-op", sla.TierUrgent, "session-attr", "5811-5")
	op.End(nil)

	spans := tracer.Spans()
	found := false
	for _, s := range spans {
		if s.Name == "named-op" {
			found = true
			if s.Attributes["medical.session_id"] != "session-attr" {
				t.Errorf("session_id 미기록: %v", s.Attributes)
			}
			if s.Attributes["medical.test_code"] != "5811-5" {
				t.Errorf("test_code 미기록")
			}
			if s.Attributes["medical.tier"] != "urgent" {
				t.Errorf("tier = %q", s.Attributes["medical.tier"])
			}
		}
	}
	if !found {
		t.Error("named-op span 미발견")
	}
}
