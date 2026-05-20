package security_test

import (
	"testing"
	"time"

	"github.com/manpasik/backend/shared/ops/security"
)

func newTestFinding(id string, sev security.Severity) *security.Finding {
	return &security.Finding{
		ID:           id,
		Type:         security.ScanCodeQL,
		Severity:     sev,
		Title:        "Test Finding " + id,
		DiscoveredAt: time.Now().UTC(),
	}
}

func TestSeverityScore(t *testing.T) {
	cases := []struct {
		sev      security.Severity
		minScore float64
	}{
		{security.SeverityCritical, 9.0},
		{security.SeverityHigh, 7.0},
		{security.SeverityMedium, 4.0},
		{security.SeverityLow, 2.0},
	}
	for _, c := range cases {
		got := security.SeverityScore(c.sev)
		if got < c.minScore {
			t.Errorf("%s = %f, want >= %f", c.sev, got, c.minScore)
		}
	}
}

func TestSeverityScore_Ordering(t *testing.T) {
	if !(security.SeverityScore(security.SeverityCritical) > security.SeverityScore(security.SeverityHigh)) {
		t.Error("Critical >= High 위반")
	}
	if !(security.SeverityScore(security.SeverityHigh) > security.SeverityScore(security.SeverityMedium)) {
		t.Error("High > Medium 위반")
	}
}

func TestPolicyEvaluator_BlockOnCritical(t *testing.T) {
	e := security.NewPolicyEvaluator(security.DefaultPolicy)
	report := &security.ScanReport{
		ID: "x", Type: security.ScanCodeQL,
		Findings: []*security.Finding{
			newTestFinding("c1", security.SeverityCritical),
		},
	}
	d := e.Evaluate(report)
	if d.Allowed {
		t.Error("critical 발견에도 통과")
	}
	if len(d.BlockedBy) != 1 {
		t.Errorf("BlockedBy = %d, want 1", len(d.BlockedBy))
	}
}

func TestPolicyEvaluator_AllowsLow(t *testing.T) {
	e := security.NewPolicyEvaluator(security.DefaultPolicy)
	report := &security.ScanReport{
		ID: "x", Type: security.ScanCodeQL,
		Findings: []*security.Finding{
			newTestFinding("l1", security.SeverityLow),
			newTestFinding("l2", security.SeverityInfo),
		},
	}
	d := e.Evaluate(report)
	if !d.Allowed {
		t.Errorf("Low/Info만 있는 보고서가 차단됨: %s", d.Reason)
	}
}

func TestPolicyEvaluator_MaxMediumExceeded(t *testing.T) {
	policy := security.DefaultPolicy
	policy.MaxMediumPerService = 2
	e := security.NewPolicyEvaluator(policy)

	findings := []*security.Finding{
		newTestFinding("m1", security.SeverityMedium),
		newTestFinding("m2", security.SeverityMedium),
		newTestFinding("m3", security.SeverityMedium),
	}
	report := &security.ScanReport{ID: "x", Findings: findings}
	d := e.Evaluate(report)
	if d.Allowed {
		t.Error("medium 초과인데 통과")
	}
}

func TestPolicyEvaluator_AcknowledgedSkipped(t *testing.T) {
	e := security.NewPolicyEvaluator(security.DefaultPolicy)
	f := newTestFinding("ack", security.SeverityCritical)
	f.Acknowledged = true

	report := &security.ScanReport{ID: "x", Findings: []*security.Finding{f}}
	d := e.Evaluate(report)
	if !d.Allowed {
		t.Errorf("ack된 finding이 차단함: %s", d.Reason)
	}
}

func TestPolicyEvaluator_ExemptCVE(t *testing.T) {
	policy := security.DefaultPolicy
	policy.ExemptCVEs = []string{"CVE-2024-9999"}
	e := security.NewPolicyEvaluator(policy)

	f := newTestFinding("cve", security.SeverityCritical)
	f.CVE = "CVE-2024-9999"

	report := &security.ScanReport{ID: "x", Findings: []*security.Finding{f}}
	d := e.Evaluate(report)
	if !d.Allowed {
		t.Error("예외 CVE가 차단됨")
	}
}

func TestPolicyEvaluator_OldFindingBlocked(t *testing.T) {
	e := security.NewPolicyEvaluator(security.DefaultPolicy)
	old := newTestFinding("old", security.SeverityLow)
	old.DiscoveredAt = time.Now().UTC().AddDate(0, 0, -60) // 60일 전

	report := &security.ScanReport{ID: "x", Findings: []*security.Finding{old}}
	d := e.Evaluate(report)
	if d.Allowed {
		t.Error("30일 초과 finding이 통과")
	}
}

func TestAggregator_AddGet(t *testing.T) {
	a := security.NewAggregator()
	r := &security.ScanReport{
		ID:   "r1",
		Type: security.ScanCodeQL,
		Findings: []*security.Finding{
			newTestFinding("f1", security.SeverityHigh),
		},
	}
	if err := a.Add(r); err != nil {
		t.Fatalf("Add 실패: %v", err)
	}
	got, _ := a.Get("r1")
	if len(got.Findings) != 1 {
		t.Errorf("Findings = %d", len(got.Findings))
	}
}

func TestAggregator_AllFindingsSorted(t *testing.T) {
	a := security.NewAggregator()
	_ = a.Add(&security.ScanReport{
		ID: "r1",
		Findings: []*security.Finding{
			newTestFinding("low", security.SeverityLow),
			newTestFinding("crit", security.SeverityCritical),
			newTestFinding("med", security.SeverityMedium),
		},
	})

	all := a.AllFindings()
	if len(all) != 3 {
		t.Errorf("findings = %d, want 3", len(all))
	}
	if all[0].Severity != security.SeverityCritical {
		t.Errorf("first = %s, want critical", all[0].Severity)
	}
	if all[2].Severity != security.SeverityLow {
		t.Errorf("last = %s, want low", all[2].Severity)
	}
}

func TestAggregator_FindingsByType(t *testing.T) {
	a := security.NewAggregator()
	r1 := &security.ScanReport{ID: "1", Type: security.ScanCodeQL, Findings: []*security.Finding{newTestFinding("a", security.SeverityHigh)}}
	r2 := &security.ScanReport{ID: "2", Type: security.ScanTrivy, Findings: []*security.Finding{newTestFinding("b", security.SeverityLow)}}
	_ = a.Add(r1)
	_ = a.Add(r2)

	codeql := a.FindingsByType(security.ScanCodeQL)
	if len(codeql) != 1 {
		t.Errorf("codeql = %d, want 1", len(codeql))
	}
	trivy := a.FindingsByType(security.ScanTrivy)
	if len(trivy) != 1 {
		t.Errorf("trivy = %d, want 1", len(trivy))
	}
}

func TestAggregator_CriticalCount(t *testing.T) {
	a := security.NewAggregator()
	_ = a.Add(&security.ScanReport{
		ID: "x",
		Findings: []*security.Finding{
			newTestFinding("c1", security.SeverityCritical),
			newTestFinding("c2", security.SeverityCritical),
			newTestFinding("h1", security.SeverityHigh),
		},
	})

	if a.CriticalCount() != 2 {
		t.Errorf("CriticalCount = %d, want 2", a.CriticalCount())
	}
}

func TestAggregator_AcknowledgeFinding(t *testing.T) {
	a := security.NewAggregator()
	r := &security.ScanReport{
		ID: "r",
		Findings: []*security.Finding{newTestFinding("f1", security.SeverityCritical)},
	}
	_ = a.Add(r)

	if err := a.AcknowledgeFinding("r", "f1"); err != nil {
		t.Fatalf("AcknowledgeFinding 실패: %v", err)
	}
	if a.CriticalCount() != 0 {
		t.Errorf("ack 후에도 critical count = %d", a.CriticalCount())
	}
}

func TestAggregator_SummaryByScanType(t *testing.T) {
	a := security.NewAggregator()
	_ = a.Add(&security.ScanReport{
		ID: "1", Type: security.ScanCodeQL,
		Findings: []*security.Finding{newTestFinding("a", security.SeverityLow), newTestFinding("b", security.SeverityHigh)},
	})
	_ = a.Add(&security.ScanReport{
		ID: "2", Type: security.ScanTrivy,
		Findings: []*security.Finding{newTestFinding("c", security.SeverityMedium)},
	})

	summary := a.SummaryByScanType()
	if summary[security.ScanCodeQL] != 2 {
		t.Errorf("codeql = %d, want 2", summary[security.ScanCodeQL])
	}
	if summary[security.ScanTrivy] != 1 {
		t.Errorf("trivy = %d, want 1", summary[security.ScanTrivy])
	}
}

func TestScanReport_CountBySeverity(t *testing.T) {
	r := &security.ScanReport{
		Findings: []*security.Finding{
			newTestFinding("c", security.SeverityCritical),
			newTestFinding("h1", security.SeverityHigh),
			newTestFinding("h2", security.SeverityHigh),
		},
	}
	counts := r.CountBySeverity()
	if counts[security.SeverityCritical] != 1 {
		t.Errorf("critical = %d", counts[security.SeverityCritical])
	}
	if counts[security.SeverityHigh] != 2 {
		t.Errorf("high = %d", counts[security.SeverityHigh])
	}
}
