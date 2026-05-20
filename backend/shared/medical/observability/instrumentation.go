// Package observability는 medical 모듈의 자동 계측을 제공합니다.
//
// Phase N tracing + sla 모니터링과 medical 도메인 모듈을 결합해
// 의료 워크플로우의 분산 추적과 SLA 준수를 자동으로 측정합니다.
package observability

import (
	"context"
	"fmt"
	"time"

	"github.com/manpasik/backend/shared/medical/sla"
	"github.com/manpasik/backend/shared/ops/tracing"
)

// ============================================================================
// MedicalInstrumentor
// ============================================================================

// MedicalInstrumentor는 medical 작업에 자동 계측을 적용합니다.
type MedicalInstrumentor struct {
	tracer  tracing.Tracer
	sla     *sla.Monitor
	service string
}

// NewMedicalInstrumentor는 새 계측기를 생성합니다.
func NewMedicalInstrumentor(tracer tracing.Tracer, slaMonitor *sla.Monitor, serviceName string) *MedicalInstrumentor {
	if tracer == nil {
		tracer = tracing.NewNoopTracer()
	}
	if slaMonitor == nil {
		slaMonitor = sla.NewMonitor(nil)
	}
	return &MedicalInstrumentor{
		tracer:  tracer,
		sla:     slaMonitor,
		service: serviceName,
	}
}

// ============================================================================
// 계측 컨텍스트
// ============================================================================

// InstrumentationContext는 단일 작업의 계측 핸들입니다.
type InstrumentationContext struct {
	span        *tracing.Span
	slaTimer    func() *sla.Measurement
	tracer      tracing.Tracer
	ctx         context.Context
	operationName string
	startedAt   time.Time
}

// StartOperation은 작업 시작 시 분산 추적 + SLA 타이머를 동시에 시작합니다.
//
// 사용:
//
//	inst := observability.NewMedicalInstrumentor(tracer, slaMon, "telemedicine")
//	op := inst.StartOperation(ctx, "ProcessMeasurement", sla.TierSTAT, "session-1", "2345-7")
//	defer op.End(nil)
//	// ... 작업 ...
//	op.SetAttribute("patient_id", "P-789")
//	op.End(err) // 자동으로 span/SLA 종료
func (m *MedicalInstrumentor) StartOperation(
	ctx context.Context,
	name string,
	tier sla.Tier,
	sessionID string,
	testCode string,
) *InstrumentationContext {
	spanCtx, span := m.tracer.StartSpan(ctx, name, tracing.SpanKindInternal)
	span.SetAttribute("medical.session_id", sessionID)
	span.SetAttribute("medical.test_code", testCode)
	span.SetAttribute("medical.tier", string(tier))

	timer := m.sla.StartTimer(sessionID, tier, m.service, testCode)

	return &InstrumentationContext{
		span:          span,
		slaTimer:      timer,
		tracer:        m.tracer,
		ctx:           spanCtx,
		operationName: name,
		startedAt:     time.Now().UTC(),
	}
}

// SetAttribute는 span에 속성을 추가합니다.
func (i *InstrumentationContext) SetAttribute(key, value string) {
	if i.span != nil {
		i.span.SetAttribute(key, value)
	}
}

// AddEvent는 span에 이벤트를 추가합니다.
func (i *InstrumentationContext) AddEvent(name string, attrs map[string]string) {
	if i.span != nil {
		i.span.AddEvent(name, attrs)
	}
}

// Context는 자식 작업이 사용할 추적 컨텍스트를 반환합니다.
func (i *InstrumentationContext) Context() context.Context {
	return i.ctx
}

// End는 작업 종료를 기록합니다. err이 nil이 아니면 span을 error 상태로 표시.
func (i *InstrumentationContext) End(err error) *sla.Measurement {
	if i.span != nil {
		if err != nil {
			i.span.SetStatus(tracing.StatusError, err.Error())
			i.span.SetAttribute("error", "true")
			i.span.SetAttribute("error.message", err.Error())
		} else {
			i.span.SetStatus(tracing.StatusOK, "")
		}
		i.span.End(i.tracer)
	}

	var measurement *sla.Measurement
	if i.slaTimer != nil {
		measurement = i.slaTimer()
	}

	return measurement
}

// ============================================================================
// Workflow 단계 계측
// ============================================================================

// InstrumentedWorkflow는 5계층 의료 워크플로우 자동 계측 래퍼입니다.
type InstrumentedWorkflow struct {
	instrumentor *MedicalInstrumentor
	sessionID    string
	patientID    string
}

// NewInstrumentedWorkflow는 새 워크플로우 추적기를 생성합니다.
func NewInstrumentedWorkflow(inst *MedicalInstrumentor, sessionID, patientID string) *InstrumentedWorkflow {
	return &InstrumentedWorkflow{
		instrumentor: inst,
		sessionID:    sessionID,
		patientID:    patientID,
	}
}

// TrackPOCT는 POCT 파싱 단계를 계측합니다.
func (w *InstrumentedWorkflow) TrackPOCT(ctx context.Context, fn func(context.Context) error) error {
	op := w.instrumentor.StartOperation(ctx, "medical.poct.parse", sla.TierSTAT, w.sessionID, "")
	op.SetAttribute("patient_id", w.patientID)
	err := fn(op.Context())
	op.End(err)
	return err
}

// TrackFHIR는 FHIR 변환 단계를 계측합니다.
func (w *InstrumentedWorkflow) TrackFHIR(ctx context.Context, fn func(context.Context) error) error {
	op := w.instrumentor.StartOperation(ctx, "medical.fhir.convert", sla.TierUrgent, w.sessionID, "")
	op.SetAttribute("patient_id", w.patientID)
	err := fn(op.Context())
	op.End(err)
	return err
}

// TrackReasoning는 CoT 추론 단계를 계측합니다.
func (w *InstrumentedWorkflow) TrackReasoning(ctx context.Context, fn func(context.Context) error) error {
	op := w.instrumentor.StartOperation(ctx, "medical.cot.reason", sla.TierSTAT, w.sessionID, "")
	op.SetAttribute("patient_id", w.patientID)
	err := fn(op.Context())
	op.End(err)
	return err
}

// TrackSafety는 Polaris 안전망 단계를 계측합니다.
func (w *InstrumentedWorkflow) TrackSafety(ctx context.Context, fn func(context.Context) error) error {
	op := w.instrumentor.StartOperation(ctx, "medical.safety.check", sla.TierUrgent, w.sessionID, "")
	op.SetAttribute("patient_id", w.patientID)
	err := fn(op.Context())
	op.End(err)
	return err
}

// TrackScribe는 SOAP 생성 단계를 계측합니다.
func (w *InstrumentedWorkflow) TrackScribe(ctx context.Context, fn func(context.Context) error) error {
	op := w.instrumentor.StartOperation(ctx, "medical.scribe.generate", sla.TierRoutine, w.sessionID, "")
	op.SetAttribute("patient_id", w.patientID)
	err := fn(op.Context())
	op.End(err)
	return err
}

// ============================================================================
// 종합 보고서
// ============================================================================

// WorkflowReport는 워크플로우 추적 종합 보고서입니다.
type WorkflowReport struct {
	SessionID         string
	PatientID         string
	TotalSpans        int
	ErrorSpans        int
	TotalDurationMs   int64
	StatBreaches      int
	GeneratedAt       time.Time
}

// GenerateReport는 인메모리 tracer + SLA monitor에서 보고서를 생성합니다.
//
// 운영에서는 OTel collector 또는 Jaeger에서 직접 조회 권장.
func (m *MedicalInstrumentor) GenerateReport(sessionID, patientID string) *WorkflowReport {
	report := &WorkflowReport{
		SessionID:   sessionID,
		PatientID:   patientID,
		GeneratedAt: time.Now().UTC(),
	}

	// MemoryTracer에서 trace ID 매칭 스팬 조회
	if memTracer, ok := m.tracer.(*tracing.MemoryTracer); ok {
		spans := memTracer.Spans()
		for _, span := range spans {
			if span.Attributes["medical.session_id"] != sessionID {
				continue
			}
			report.TotalSpans++
			report.TotalDurationMs += span.Duration().Milliseconds()
			if span.Status == tracing.StatusError {
				report.ErrorSpans++
			}
		}
	}

	// SLA breach 카운트는 알람에서 추출
	for _, alert := range m.sla.AlertsBySeverity("critical") {
		if alert.Service == m.service && alert.Tier == sla.TierSTAT {
			report.StatBreaches++
		}
	}
	for _, alert := range m.sla.AlertsBySeverity("warning") {
		if alert.Service == m.service && alert.Tier == sla.TierSTAT {
			report.StatBreaches++
		}
	}

	return report
}

// ============================================================================
// 프로파일러 (단순 latency 분포)
// ============================================================================

// LatencyProfile은 작업별 지연 분포 분석입니다.
type LatencyProfile struct {
	Operations map[string]*OperationStats
}

// OperationStats는 특정 작업의 통계입니다.
type OperationStats struct {
	Count       int
	TotalMs     int64
	MeanMs      float64
	MaxMs       int64
}

// ProfileLatency는 메모리 tracer에서 작업별 latency 분포를 추출합니다.
func (m *MedicalInstrumentor) ProfileLatency() *LatencyProfile {
	profile := &LatencyProfile{Operations: make(map[string]*OperationStats)}

	memTracer, ok := m.tracer.(*tracing.MemoryTracer)
	if !ok {
		return profile
	}

	for _, span := range memTracer.Spans() {
		stats, exists := profile.Operations[span.Name]
		if !exists {
			stats = &OperationStats{}
			profile.Operations[span.Name] = stats
		}
		ms := span.Duration().Milliseconds()
		stats.Count++
		stats.TotalMs += ms
		if ms > stats.MaxMs {
			stats.MaxMs = ms
		}
		if stats.Count > 0 {
			stats.MeanMs = float64(stats.TotalMs) / float64(stats.Count)
		}
	}
	return profile
}

// String returns a summary string.
func (p *LatencyProfile) String() string {
	if len(p.Operations) == 0 {
		return "[empty profile]"
	}
	out := "Latency Profile:\n"
	for name, stats := range p.Operations {
		out += fmt.Sprintf("  %s: count=%d mean=%.1fms max=%dms\n",
			name, stats.Count, stats.MeanMs, stats.MaxMs)
	}
	return out
}
