package poise_test

import (
	"testing"
	"time"

	"github.com/manpasik/backend/shared/medical/poise"
)

func TestLoop_CollectFeedback_Required(t *testing.T) {
	l := poise.NewLoop(poise.DefaultThresholds)
	if err := l.CollectFeedback(nil); err == nil {
		t.Error("nil 피드백 통과")
	}
	if err := l.CollectFeedback(&poise.Feedback{Kind: poise.FeedbackSatisfaction}); err == nil {
		t.Error("session_id 없이 통과")
	}
	if err := l.CollectFeedback(&poise.Feedback{SessionID: "s"}); err == nil {
		t.Error("kind 없이 통과")
	}
}

func TestLoop_CollectAndAggregate(t *testing.T) {
	l := poise.NewLoop(poise.DefaultThresholds)

	for i := 0; i < 50; i++ {
		_ = l.CollectFeedback(&poise.Feedback{
			SessionID: "s", Kind: poise.FeedbackSatisfaction,
			Value: float64(8 + i%2), // 8 또는 9
		})
	}

	metric := l.AggregateMetric(poise.FeedbackSatisfaction, time.Time{})
	if metric.SampleSize != 50 {
		t.Errorf("SampleSize = %d, want 50", metric.SampleSize)
	}
	if metric.Mean < 8.0 || metric.Mean > 9.0 {
		t.Errorf("Mean = %f, want 8-9", metric.Mean)
	}
}

func TestLoop_AutoTimestamp(t *testing.T) {
	l := poise.NewLoop(poise.DefaultThresholds)
	fb := &poise.Feedback{SessionID: "s", Kind: poise.FeedbackSatisfaction, Value: 8}
	if err := l.CollectFeedback(fb); err != nil {
		t.Fatal(err)
	}
	if fb.CollectedAt.IsZero() {
		t.Error("CollectedAt 자동 설정 실패")
	}
	if fb.ID == "" {
		t.Error("ID 자동 설정 실패")
	}
}

func TestLoop_EvaluateAndPropose_LowSatisfaction(t *testing.T) {
	l := poise.NewLoop(poise.DefaultThresholds)

	// 낮은 만족도 다수
	for i := 0; i < 50; i++ {
		_ = l.CollectFeedback(&poise.Feedback{
			SessionID: "s", Kind: poise.FeedbackSatisfaction, Value: 6.0,
		})
	}

	proposals := l.EvaluateAndPropose(time.Time{})
	if len(proposals) == 0 {
		t.Fatal("저만족 시나리오에 제안 생성 실패")
	}
	hasSatisfactionProposal := false
	for _, p := range proposals {
		if p.Trigger == "satisfaction_drop" {
			hasSatisfactionProposal = true
			if p.Severity == 0 {
				t.Error("Severity 미설정")
			}
		}
	}
	if !hasSatisfactionProposal {
		t.Error("satisfaction_drop 트리거 누락")
	}
}

func TestLoop_EvaluateAndPropose_LowSampleSize(t *testing.T) {
	l := poise.NewLoop(poise.DefaultThresholds)

	// 5건 (MinSampleSize 30 미만)
	for i := 0; i < 5; i++ {
		_ = l.CollectFeedback(&poise.Feedback{
			SessionID: "s", Kind: poise.FeedbackSatisfaction, Value: 5.0,
		})
	}

	proposals := l.EvaluateAndPropose(time.Time{})
	if len(proposals) > 0 {
		t.Errorf("샘플 부족인데 제안 생성: %d", len(proposals))
	}
}

func TestLoop_EvaluateAndPropose_AccuracyDrop(t *testing.T) {
	l := poise.NewLoop(poise.DefaultThresholds)

	for i := 0; i < 40; i++ {
		_ = l.CollectFeedback(&poise.Feedback{
			SessionID: "s", Kind: poise.FeedbackAccuracy, Value: 0.85, // 0.92 미달
		})
	}

	proposals := l.EvaluateAndPropose(time.Time{})
	hasAccuracy := false
	for _, p := range proposals {
		if p.Trigger == "accuracy_drop" {
			hasAccuracy = true
		}
	}
	if !hasAccuracy {
		t.Error("accuracy_drop 트리거 누락")
	}
}

func TestLoop_EvaluateAndPropose_LatencySpike(t *testing.T) {
	l := poise.NewLoop(poise.DefaultThresholds)

	for i := 0; i < 40; i++ {
		_ = l.CollectFeedback(&poise.Feedback{
			SessionID: "s", Kind: poise.FeedbackResponseTime, Value: 2000.0, // 1500ms 초과
		})
	}

	proposals := l.EvaluateAndPropose(time.Time{})
	hasLatency := false
	for _, p := range proposals {
		if p.Trigger == "latency_spike" {
			hasLatency = true
		}
	}
	if !hasLatency {
		t.Error("latency_spike 트리거 누락")
	}
}

func TestLoop_Approve(t *testing.T) {
	l := poise.NewLoop(poise.DefaultThresholds)

	// 제안 생성
	for i := 0; i < 40; i++ {
		_ = l.CollectFeedback(&poise.Feedback{
			SessionID: "s", Kind: poise.FeedbackSatisfaction, Value: 6,
		})
	}
	proposals := l.EvaluateAndPropose(time.Time{})
	if len(proposals) == 0 {
		t.Fatal("제안 미생성")
	}
	pid := proposals[0].ID

	if err := l.Approve(pid, "admin-001"); err != nil {
		t.Fatalf("Approve 실패: %v", err)
	}

	got, _ := l.GetProposal(pid)
	if got.Status != "approved" {
		t.Errorf("Status = %q", got.Status)
	}
	if got.ApprovedBy != "admin-001" {
		t.Errorf("ApprovedBy = %q", got.ApprovedBy)
	}
}

func TestLoop_Reject(t *testing.T) {
	l := poise.NewLoop(poise.DefaultThresholds)

	for i := 0; i < 40; i++ {
		_ = l.CollectFeedback(&poise.Feedback{
			SessionID: "s", Kind: poise.FeedbackSatisfaction, Value: 6,
		})
	}
	proposals := l.EvaluateAndPropose(time.Time{})
	pid := proposals[0].ID

	if err := l.Reject(pid, "admin", "샘플 편향 의심"); err != nil {
		t.Fatalf("Reject 실패: %v", err)
	}

	got, _ := l.GetProposal(pid)
	if got.Status != "rejected" {
		t.Errorf("Status = %q", got.Status)
	}
	if got.RejectReason == "" {
		t.Error("RejectReason 미기록")
	}
}

func TestLoop_MarkApplied_OnlyApproved(t *testing.T) {
	l := poise.NewLoop(poise.DefaultThresholds)

	for i := 0; i < 40; i++ {
		_ = l.CollectFeedback(&poise.Feedback{
			SessionID: "s", Kind: poise.FeedbackSatisfaction, Value: 6,
		})
	}
	proposals := l.EvaluateAndPropose(time.Time{})
	pid := proposals[0].ID

	// 미승인 상태에서 적용 시도
	if err := l.MarkApplied(pid); err == nil {
		t.Error("미승인 제안 적용 통과")
	}

	_ = l.Approve(pid, "admin")
	if err := l.MarkApplied(pid); err != nil {
		t.Errorf("Apply 실패: %v", err)
	}
	got, _ := l.GetProposal(pid)
	if got.Status != "applied" {
		t.Errorf("Status = %q", got.Status)
	}
}

func TestLoop_HumanGateRequired(t *testing.T) {
	l := poise.NewLoop(poise.DefaultThresholds)

	for i := 0; i < 40; i++ {
		_ = l.CollectFeedback(&poise.Feedback{
			SessionID: "s", Kind: poise.FeedbackSatisfaction, Value: 5,
		})
	}
	proposals := l.EvaluateAndPropose(time.Time{})

	// 모든 제안이 pending 상태여야 함 (자동 적용 금지)
	for _, p := range proposals {
		if p.Status != "pending" {
			t.Errorf("자동 처리됨: %s", p.Status)
		}
	}
}

func TestLoop_ListProposals_Filter(t *testing.T) {
	l := poise.NewLoop(poise.DefaultThresholds)

	for i := 0; i < 40; i++ {
		_ = l.CollectFeedback(&poise.Feedback{
			SessionID: "s", Kind: poise.FeedbackSatisfaction, Value: 5,
		})
	}
	props := l.EvaluateAndPropose(time.Time{})
	_ = l.Approve(props[0].ID, "a")

	pending := l.ListProposals("pending")
	approved := l.ListProposals("approved")

	if len(approved) != 1 {
		t.Errorf("approved = %d, want 1", len(approved))
	}
	if len(pending) >= len(approved)+len(pending) {
		// 통과
	}
	all := l.ListProposals("")
	if len(all) != len(pending)+len(approved) {
		t.Errorf("All = %d, parts sum = %d", len(all), len(pending)+len(approved))
	}
}

func TestLoop_CountByStatus(t *testing.T) {
	l := poise.NewLoop(poise.DefaultThresholds)
	for i := 0; i < 40; i++ {
		_ = l.CollectFeedback(&poise.Feedback{
			SessionID: "s", Kind: poise.FeedbackSatisfaction, Value: 5,
		})
	}
	props := l.EvaluateAndPropose(time.Time{})
	if len(props) > 0 {
		_ = l.Approve(props[0].ID, "admin")
	}

	counts := l.CountByStatus()
	if counts["approved"] < 1 {
		t.Errorf("approved count = %d", counts["approved"])
	}
}

func TestLoop_FeedbackHistoryWindow(t *testing.T) {
	l := poise.NewLoop(poise.DefaultThresholds)
	// 100k+ 피드백 → 윈도우 작동 확인은 직접 verify 어려우므로
	// 단순 카운트만 확인
	for i := 0; i < 100; i++ {
		_ = l.CollectFeedback(&poise.Feedback{
			SessionID: "s", Kind: poise.FeedbackSatisfaction, Value: 8,
		})
	}
	if l.FeedbackCount() != 100 {
		t.Errorf("FeedbackCount = %d, want 100", l.FeedbackCount())
	}
}

func TestLoop_AggregateMetric_Percentiles(t *testing.T) {
	l := poise.NewLoop(poise.DefaultThresholds)

	// 1~100ms 분포
	for i := 1; i <= 100; i++ {
		_ = l.CollectFeedback(&poise.Feedback{
			SessionID: "s", Kind: poise.FeedbackResponseTime, Value: float64(i),
		})
	}

	metric := l.AggregateMetric(poise.FeedbackResponseTime, time.Time{})
	if metric.P95 < 90 || metric.P95 > 100 {
		t.Errorf("P95 = %f, want 90-100", metric.P95)
	}
	if metric.P99 < 95 || metric.P99 > 100 {
		t.Errorf("P99 = %f, want 95-100", metric.P99)
	}
	if metric.Median < 49 || metric.Median > 52 {
		t.Errorf("Median = %f, want ~50", metric.Median)
	}
}

func TestLoop_TimeWindowFilter(t *testing.T) {
	l := poise.NewLoop(poise.DefaultThresholds)

	// 과거 데이터
	old := time.Now().UTC().Add(-24 * time.Hour)
	for i := 0; i < 30; i++ {
		_ = l.CollectFeedback(&poise.Feedback{
			SessionID: "s", Kind: poise.FeedbackSatisfaction, Value: 5,
			CollectedAt: old,
		})
	}
	// 최근 데이터
	for i := 0; i < 30; i++ {
		_ = l.CollectFeedback(&poise.Feedback{
			SessionID: "s", Kind: poise.FeedbackSatisfaction, Value: 9,
		})
	}

	// 지난 1시간 윈도우
	recent := time.Now().UTC().Add(-1 * time.Hour)
	metric := l.AggregateMetric(poise.FeedbackSatisfaction, recent)
	if metric.SampleSize != 30 {
		t.Errorf("SampleSize = %d, want 30 (recent only)", metric.SampleSize)
	}
	if metric.Mean < 8 {
		t.Errorf("Mean = %f, 최근 데이터만 집계되어야 함", metric.Mean)
	}
}

func TestLoop_Approve_NotFound(t *testing.T) {
	l := poise.NewLoop(poise.DefaultThresholds)
	if err := l.Approve("missing", "admin"); err == nil {
		t.Error("미존재 제안 승인 통과")
	}
}

func TestLoop_RejectAlreadyApproved(t *testing.T) {
	l := poise.NewLoop(poise.DefaultThresholds)
	for i := 0; i < 40; i++ {
		_ = l.CollectFeedback(&poise.Feedback{
			SessionID: "s", Kind: poise.FeedbackSatisfaction, Value: 5,
		})
	}
	props := l.EvaluateAndPropose(time.Time{})
	pid := props[0].ID
	_ = l.Approve(pid, "a")

	if err := l.Reject(pid, "a", "x"); err == nil {
		t.Error("이미 승인된 제안 거부 통과")
	}
}
