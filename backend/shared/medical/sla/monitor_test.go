package sla_test

import (
	"testing"
	"time"

	"github.com/manpasik/backend/shared/medical/sla"
)

func TestMonitor_StatRecord_WithinSLA(t *testing.T) {
	m := sla.NewMonitor(nil)

	measurement := &sla.Measurement{
		SessionID: "s-1",
		Tier:      sla.TierSTAT,
		StartedAt: time.Now().UTC(),
		ResultAt:  time.Now().UTC().Add(15 * time.Second), // 15초 < 30초
		Service:   "ai-inference",
	}
	if err := m.Record(measurement); err != nil {
		t.Fatalf("Record 실패: %v", err)
	}
	if !measurement.WithinSLA {
		t.Error("15초가 STAT 30초 SLA 위반으로 잘못 판정")
	}
	if measurement.Breach {
		t.Error("정상 응답이 breach로 판정")
	}
}

func TestMonitor_StatRecord_Breach(t *testing.T) {
	m := sla.NewMonitor(nil)

	now := time.Now().UTC()
	measurement := &sla.Measurement{
		SessionID: "s-2",
		Tier:      sla.TierSTAT,
		StartedAt: now,
		ResultAt:  now.Add(45 * time.Second), // 45초 > 30초 + 15초(50%) → critical
		Service:   "ai-inference",
	}
	_ = m.Record(measurement)

	if !measurement.Breach {
		t.Error("breach 미감지")
	}
	if measurement.BreachSeverity != "critical" {
		t.Errorf("Severity = %q, want critical (45s > 30+15s)", measurement.BreachSeverity)
	}

	if m.AlertCount() != 1 {
		t.Errorf("알람 = %d, want 1", m.AlertCount())
	}
}

func TestMonitor_WarningSeverity(t *testing.T) {
	m := sla.NewMonitor(nil)

	now := time.Now().UTC()
	measurement := &sla.Measurement{
		SessionID: "s-3",
		Tier:      sla.TierSTAT,
		StartedAt: now,
		ResultAt:  now.Add(35 * time.Second), // 35초 (30초+5초) → warning
	}
	_ = m.Record(measurement)
	if measurement.BreachSeverity != "warning" {
		t.Errorf("Severity = %q, want warning", measurement.BreachSeverity)
	}
}

func TestMonitor_StartTimer(t *testing.T) {
	m := sla.NewMonitor(nil)

	stop := m.StartTimer("session-timer", sla.TierSTAT, "service-a", "2345-7")
	time.Sleep(20 * time.Millisecond)
	measurement := stop()

	if measurement.DurationMs < 15 || measurement.DurationMs > 100 {
		t.Errorf("DurationMs = %d, want ~20ms", measurement.DurationMs)
	}
	if !measurement.WithinSLA {
		t.Error("20ms가 SLA 위반")
	}
}

func TestMonitor_TierUrgent(t *testing.T) {
	m := sla.NewMonitor(nil)

	now := time.Now().UTC()
	measurement := &sla.Measurement{
		SessionID: "s-urgent",
		Tier:      sla.TierUrgent,
		StartedAt: now,
		ResultAt:  now.Add(50 * time.Second), // 50초 < 60초 (urgent SLA)
	}
	_ = m.Record(measurement)
	if !measurement.WithinSLA {
		t.Error("urgent 50초가 SLA 위반으로 잘못 판정")
	}
}

func TestMonitor_Stats_Percentiles(t *testing.T) {
	m := sla.NewMonitor(nil)

	// 100개 샘플 (1~100ms)
	now := time.Now().UTC()
	for i := 1; i <= 100; i++ {
		measurement := &sla.Measurement{
			SessionID: "s",
			Tier:      sla.TierSTAT,
			StartedAt: now,
			ResultAt:  now.Add(time.Duration(i) * time.Millisecond),
		}
		_ = m.Record(measurement)
	}

	stats := m.Stats(sla.TierSTAT, time.Time{})
	if stats.TotalSamples != 100 {
		t.Errorf("Samples = %d, want 100", stats.TotalSamples)
	}
	if stats.P95Ms < 90 || stats.P95Ms > 100 {
		t.Errorf("P95 = %f, want 90-100", stats.P95Ms)
	}
	if stats.BreachCount != 0 {
		t.Errorf("BreachCount = %d, want 0 (모두 30초 이내)", stats.BreachCount)
	}
}

func TestMonitor_Stats_HighBreachRate(t *testing.T) {
	m := sla.NewMonitor(nil)

	now := time.Now().UTC()
	// 50% breach
	for i := 0; i < 100; i++ {
		dur := 15 * time.Second
		if i%2 == 0 {
			dur = 45 * time.Second
		}
		_ = m.Record(&sla.Measurement{
			SessionID: "s", Tier: sla.TierSTAT,
			StartedAt: now, ResultAt: now.Add(dur),
		})
	}

	stats := m.Stats(sla.TierSTAT, time.Time{})
	if stats.BreachCount != 50 {
		t.Errorf("BreachCount = %d, want 50", stats.BreachCount)
	}
	if stats.BreachRate < 0.49 || stats.BreachRate > 0.51 {
		t.Errorf("BreachRate = %f, want 0.5", stats.BreachRate)
	}
}

func TestMonitor_EvaluateBurnRate(t *testing.T) {
	m := sla.NewMonitor(nil)

	now := time.Now().UTC()
	// 30% breach (default 임계값 10% 초과)
	for i := 0; i < 100; i++ {
		dur := 15 * time.Second
		if i < 30 {
			dur = 45 * time.Second
		}
		_ = m.Record(&sla.Measurement{
			SessionID: "s", Tier: sla.TierSTAT,
			StartedAt: now, ResultAt: now.Add(dur),
		})
	}

	alert := m.EvaluateBurnRate(sla.TierSTAT, time.Time{})
	if alert == nil {
		t.Fatal("burn rate 임계 초과인데 알람 없음")
	}
	if alert.Type != "burn_rate" {
		t.Errorf("Type = %q", alert.Type)
	}
	if alert.BurnRate < 0.29 || alert.BurnRate > 0.31 {
		t.Errorf("BurnRate = %f, want 0.3", alert.BurnRate)
	}
}

func TestMonitor_BurnRate_NoAlert(t *testing.T) {
	m := sla.NewMonitor(nil)

	now := time.Now().UTC()
	// 5% breach만 (10% 미만 → 알람 없음)
	for i := 0; i < 100; i++ {
		dur := 15 * time.Second
		if i < 5 {
			dur = 45 * time.Second
		}
		_ = m.Record(&sla.Measurement{
			SessionID: "s", Tier: sla.TierSTAT,
			StartedAt: now, ResultAt: now.Add(dur),
		})
	}

	alert := m.EvaluateBurnRate(sla.TierSTAT, time.Time{})
	if alert != nil {
		t.Errorf("5%% breach인데 burn rate 알람: %+v", alert)
	}
}

func TestMonitor_AlertsBySeverity(t *testing.T) {
	m := sla.NewMonitor(nil)

	now := time.Now().UTC()
	// 1 critical + 2 warning
	_ = m.Record(&sla.Measurement{SessionID: "s", Tier: sla.TierSTAT,
		StartedAt: now, ResultAt: now.Add(50 * time.Second)}) // critical
	_ = m.Record(&sla.Measurement{SessionID: "s", Tier: sla.TierSTAT,
		StartedAt: now, ResultAt: now.Add(35 * time.Second)}) // warning
	_ = m.Record(&sla.Measurement{SessionID: "s", Tier: sla.TierSTAT,
		StartedAt: now, ResultAt: now.Add(36 * time.Second)}) // warning

	if len(m.AlertsBySeverity("critical")) != 1 {
		t.Errorf("critical = %d, want 1", len(m.AlertsBySeverity("critical")))
	}
	if len(m.AlertsBySeverity("warning")) != 2 {
		t.Errorf("warning = %d, want 2", len(m.AlertsBySeverity("warning")))
	}
}

func TestMonitor_AcknowledgeAlert(t *testing.T) {
	m := sla.NewMonitor(nil)

	now := time.Now().UTC()
	_ = m.Record(&sla.Measurement{SessionID: "s", Tier: sla.TierSTAT,
		StartedAt: now, ResultAt: now.Add(50 * time.Second)})

	criticals := m.AlertsBySeverity("critical")
	if len(criticals) == 0 {
		t.Fatal("알람 없음")
	}
	id := criticals[0].ID

	if err := m.AcknowledgeAlert(id); err != nil {
		t.Fatalf("Ack 실패: %v", err)
	}
	if !criticals[0].Acknowledged {
		t.Error("Acknowledged = false")
	}
}

func TestMonitor_RecordValidation(t *testing.T) {
	m := sla.NewMonitor(nil)
	if err := m.Record(nil); err == nil {
		t.Error("nil 통과")
	}
	if err := m.Record(&sla.Measurement{Tier: sla.TierSTAT}); err == nil {
		t.Error("session_id 없이 통과")
	}
}

func TestMonitor_DefaultTierFromEmpty(t *testing.T) {
	m := sla.NewMonitor(nil)
	now := time.Now().UTC()

	measurement := &sla.Measurement{
		SessionID: "s", // tier 미지정 → STAT 기본값
		StartedAt: now, ResultAt: now.Add(15 * time.Second),
	}
	if err := m.Record(measurement); err != nil {
		t.Fatalf("Record 실패: %v", err)
	}
	if measurement.Tier != sla.TierSTAT {
		t.Errorf("Tier = %q, want stat (default)", measurement.Tier)
	}
}

func TestMonitor_MeasurementCount(t *testing.T) {
	m := sla.NewMonitor(nil)
	now := time.Now().UTC()

	for i := 0; i < 5; i++ {
		_ = m.Record(&sla.Measurement{
			SessionID: "s", Tier: sla.TierSTAT,
			StartedAt: now, ResultAt: now.Add(10 * time.Second),
		})
	}
	if m.MeasurementCount() != 5 {
		t.Errorf("Count = %d, want 5", m.MeasurementCount())
	}
}

func TestMonitor_BudgetUtilization(t *testing.T) {
	m := sla.NewMonitor(nil)

	now := time.Now().UTC()
	for i := 0; i < 50; i++ {
		_ = m.Record(&sla.Measurement{
			SessionID: "s", Tier: sla.TierSTAT,
			StartedAt: now, ResultAt: now.Add(15 * time.Second), // 50% utilization
		})
	}

	stats := m.Stats(sla.TierSTAT, time.Time{})
	if stats.BudgetUtilization < 0.4 || stats.BudgetUtilization > 0.6 {
		t.Errorf("Utilization = %f, want ~0.5", stats.BudgetUtilization)
	}
}

func TestMonitor_UnknownTier(t *testing.T) {
	m := sla.NewMonitor(nil)
	if err := m.Record(&sla.Measurement{
		SessionID: "s", Tier: "unknown_tier",
		StartedAt: time.Now(), ResultAt: time.Now().Add(time.Second),
	}); err == nil {
		t.Error("미지의 tier 통과")
	}
}

func TestMonitor_CustomTargets(t *testing.T) {
	custom := map[sla.Tier]sla.Target{
		"strict": {Tier: "strict", BudgetMs: 1000, ToleranceMs: 500, BurnRateAlert: 0.05},
	}
	m := sla.NewMonitor(custom)

	now := time.Now().UTC()
	_ = m.Record(&sla.Measurement{
		SessionID: "s", Tier: "strict",
		StartedAt: now, ResultAt: now.Add(800 * time.Millisecond),
	})

	stats := m.Stats("strict", time.Time{})
	if stats.BudgetMs != 1000 {
		t.Errorf("Budget = %d, want 1000", stats.BudgetMs)
	}
	if stats.BreachCount != 0 {
		t.Errorf("Breach = %d, want 0", stats.BreachCount)
	}
}
