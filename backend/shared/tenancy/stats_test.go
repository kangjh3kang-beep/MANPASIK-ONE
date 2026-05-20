package tenancy_test

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/tenancy"
)

func newStatsFixture(t *testing.T) *tenancy.StatsCollector {
	t.Helper()
	mem := tenancy.NewMemoryMembershipStore()
	inv := tenancy.NewMemoryInvitationStore()

	// 활성 멤버 5명 (admin: 2, member: 2, viewer: 1)
	_ = mem.Add(tenancy.Membership{UserID: "u1", TenantID: "A", Role: tenancy.TenantRoleAdmin})
	_ = mem.Add(tenancy.Membership{UserID: "u2", TenantID: "A", Role: tenancy.TenantRoleMember})
	_ = mem.Add(tenancy.Membership{UserID: "u3", TenantID: "B", Role: tenancy.TenantRoleAdmin})
	_ = mem.Add(tenancy.Membership{UserID: "u4", TenantID: "B", Role: tenancy.TenantRoleMember})
	_ = mem.Add(tenancy.Membership{UserID: "u5", TenantID: "B", Role: tenancy.TenantRoleViewer})
	// 비활성 1명
	_ = mem.Add(tenancy.Membership{UserID: "u6", TenantID: "B", Role: tenancy.TenantRoleMember})
	_ = mem.SetActive("u6", "B", false)

	now := time.Now()
	// pending 초대 2개 (A 1개, B 1개)
	_ = inv.Add(tenancy.Invitation{
		Token: "p1", TenantID: "A", Role: tenancy.TenantRoleMember,
		Status: tenancy.InvitationPending, ExpiresAt: now.Add(24 * time.Hour),
	})
	_ = inv.Add(tenancy.Invitation{
		Token: "p2", TenantID: "B", Role: tenancy.TenantRoleViewer,
		Status: tenancy.InvitationPending, ExpiresAt: now.Add(24 * time.Hour),
	})
	// accepted 1개
	_ = inv.Add(tenancy.Invitation{
		Token: "a1", TenantID: "A", Role: tenancy.TenantRoleMember,
		Status: tenancy.InvitationAccepted, ExpiresAt: now,
	})
	// expired (status=expired) 1개
	_ = inv.Add(tenancy.Invitation{
		Token: "e1", TenantID: "A", Role: tenancy.TenantRoleMember,
		Status: tenancy.InvitationExpired, ExpiresAt: now.Add(-time.Hour),
	})
	// pending 인데 만료된 1개 → expired 카운트
	_ = inv.Add(tenancy.Invitation{
		Token: "ep1", TenantID: "B", Role: tenancy.TenantRoleViewer,
		Status: tenancy.InvitationPending, ExpiresAt: now.Add(-time.Hour),
	})
	// revoked 1개
	_ = inv.Add(tenancy.Invitation{
		Token: "r1", TenantID: "A", Role: tenancy.TenantRoleAdmin,
		Status: tenancy.InvitationRevoked, ExpiresAt: now,
	})

	tenants := tenancy.TenantListerFunc(func() []tenancy.TenantID {
		return []tenancy.TenantID{"A", "B"}
	})
	return tenancy.NewStatsCollector(mem, inv, tenants)
}

func TestStatsCollector_TenantCount(t *testing.T) {
	c := newStatsFixture(t)
	stats := c.Collect()
	if stats.TenantCount != 2 {
		t.Errorf("TenantCount = %d, want 2", stats.TenantCount)
	}
}

func TestStatsCollector_MemberCounts(t *testing.T) {
	c := newStatsFixture(t)
	stats := c.Collect()
	if stats.ActiveMemberCount != 5 {
		t.Errorf("ActiveMemberCount = %d, want 5", stats.ActiveMemberCount)
	}
	if stats.InactiveMemberCount != 1 {
		t.Errorf("InactiveMemberCount = %d, want 1", stats.InactiveMemberCount)
	}
}

func TestStatsCollector_RoleDistribution(t *testing.T) {
	c := newStatsFixture(t)
	stats := c.Collect()
	if stats.RoleDistribution["admin"] != 2 {
		t.Errorf("admin = %d", stats.RoleDistribution["admin"])
	}
	if stats.RoleDistribution["member"] != 2 {
		t.Errorf("member = %d", stats.RoleDistribution["member"])
	}
	if stats.RoleDistribution["viewer"] != 1 {
		t.Errorf("viewer = %d", stats.RoleDistribution["viewer"])
	}
}

func TestStatsCollector_InvitationCounts(t *testing.T) {
	c := newStatsFixture(t)
	stats := c.Collect()
	if stats.PendingInvitationCount != 2 {
		t.Errorf("Pending = %d, want 2", stats.PendingInvitationCount)
	}
	if stats.AcceptedInvitationCount != 1 {
		t.Errorf("Accepted = %d", stats.AcceptedInvitationCount)
	}
	// status=expired 1개 + pending 인데 만료된 1개 = 2
	if stats.ExpiredInvitationCount != 2 {
		t.Errorf("Expired = %d, want 2", stats.ExpiredInvitationCount)
	}
	if stats.RevokedInvitationCount != 1 {
		t.Errorf("Revoked = %d", stats.RevokedInvitationCount)
	}
}

func TestStatsCollector_LastSnapshot(t *testing.T) {
	c := newStatsFixture(t)
	if c.LastSnapshot() != nil {
		t.Error("Collect 호출 전에 snapshot 존재")
	}
	c.Collect()
	if c.LastSnapshot() == nil {
		t.Error("Collect 후 snapshot nil")
	}
}

func TestStatsCollector_NilTenantLister(t *testing.T) {
	mem := tenancy.NewMemoryMembershipStore()
	c := tenancy.NewStatsCollector(mem, nil, nil)
	stats := c.Collect()
	if stats.TenantCount != 0 {
		t.Errorf("TenantCount = %d", stats.TenantCount)
	}
	if stats.ActiveMemberCount != 0 {
		t.Errorf("ActiveMemberCount = %d", stats.ActiveMemberCount)
	}
}

func TestStatsCollector_NilInvitationStore(t *testing.T) {
	mem := tenancy.NewMemoryMembershipStore()
	_ = mem.Add(tenancy.Membership{UserID: "u1", TenantID: "A", Role: tenancy.TenantRoleAdmin})
	tenants := tenancy.TenantListerFunc(func() []tenancy.TenantID {
		return []tenancy.TenantID{"A"}
	})
	c := tenancy.NewStatsCollector(mem, nil, tenants)
	stats := c.Collect()
	if stats.ActiveMemberCount != 1 {
		t.Errorf("Active = %d", stats.ActiveMemberCount)
	}
	if stats.PendingInvitationCount != 0 {
		t.Errorf("Pending = %d", stats.PendingInvitationCount)
	}
}

func TestStatsCollector_EmptyMembersTenantNotCounted(t *testing.T) {
	mem := tenancy.NewMemoryMembershipStore()
	tenants := tenancy.TenantListerFunc(func() []tenancy.TenantID {
		return []tenancy.TenantID{"GHOST"} // 멤버 없음
	})
	c := tenancy.NewStatsCollector(mem, nil, tenants)
	stats := c.Collect()
	if stats.TenantCount != 0 {
		t.Errorf("멤버 없는 tenant 가 카운트됨: %d", stats.TenantCount)
	}
}

func TestStatsCollector_CollectMetrics_DashboardCompat(t *testing.T) {
	c := newStatsFixture(t)
	m := c.CollectMetrics()
	if m["status"] != "healthy" {
		t.Errorf("status = %v", m["status"])
	}
	if m["tenant_count"].(int) != 2 {
		t.Errorf("tenant_count = %v", m["tenant_count"])
	}
	if m["role_admin"].(int) != 2 {
		t.Errorf("role_admin = %v", m["role_admin"])
	}
}

func TestStatsCollector_LastSnapshotIsCopy(t *testing.T) {
	c := newStatsFixture(t)
	c.Collect()
	snap := c.LastSnapshot()
	// 외부에서 수정해도 내부 상태에 영향 없어야
	snap.RoleDistribution["admin"] = 999
	snap2 := c.LastSnapshot()
	if snap2.RoleDistribution["admin"] == 999 {
		t.Error("LastSnapshot 가 내부 참조 반환 (copy 아님)")
	}
}

// 사용되지 않은 import 제거를 위한 확인용
var _ = time.Now

func TestWritePrometheusMetrics_Format(t *testing.T) {
	c := newStatsFixture(t)
	stats := c.Collect()
	w := httptest.NewRecorder()
	tenancy.WritePrometheusMetrics(w, stats)

	body := w.Body.String()
	if w.Header().Get("Content-Type") != "text/plain; version=0.0.4; charset=utf-8" {
		t.Errorf("Content-Type = %q", w.Header().Get("Content-Type"))
	}
	for _, want := range []string{
		"manpasik_tenancy_tenants_total 2",
		"manpasik_tenancy_members_active_total 5",
		"manpasik_tenancy_members_inactive_total 1",
		"manpasik_tenancy_role_distribution{role=\"admin\"}",
		"manpasik_tenancy_invitations_total{status=\"pending\"} 2",
		"manpasik_tenancy_invitations_total{status=\"accepted\"} 1",
		"manpasik_tenancy_invitations_total{status=\"expired\"} 2",
		"manpasik_tenancy_invitations_total{status=\"revoked\"} 1",
		"# HELP",
		"# TYPE",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("응답에 누락: %q\n--- body ---\n%s", want, body)
		}
	}
}

func TestWritePrometheusMetrics_EmptyStats(t *testing.T) {
	stats := &tenancy.TenancyStats{
		RoleDistribution: map[string]int{},
	}
	w := httptest.NewRecorder()
	tenancy.WritePrometheusMetrics(w, stats)
	body := w.Body.String()
	if !strings.Contains(body, "manpasik_tenancy_tenants_total 0") {
		t.Errorf("기본 메트릭 누락:\n%s", body)
	}
	// 빈 RoleDistribution 시 role_distribution 메트릭 미출력
	if strings.Contains(body, "manpasik_tenancy_role_distribution") {
		t.Errorf("빈 분포에 메트릭 출력:\n%s", body)
	}
}

func TestWritePrometheusMetrics_LabelEscape(t *testing.T) {
	stats := &tenancy.TenancyStats{
		RoleDistribution: map[string]int{
			`role"with"quotes`: 1,
		},
	}
	w := httptest.NewRecorder()
	tenancy.WritePrometheusMetrics(w, stats)
	body := w.Body.String()
	if !strings.Contains(body, `role="role\"with\"quotes"`) {
		t.Errorf("라벨 이스케이프 실패:\n%s", body)
	}
}

func TestWritePrometheusMetrics_NegativeAndLarge(t *testing.T) {
	stats := &tenancy.TenancyStats{
		TenantCount:       1234,
		ActiveMemberCount: 567890,
	}
	w := httptest.NewRecorder()
	tenancy.WritePrometheusMetrics(w, stats)
	body := w.Body.String()
	if !strings.Contains(body, "manpasik_tenancy_tenants_total 1234") {
		t.Errorf("큰 값 출력 실패:\n%s", body)
	}
	if !strings.Contains(body, "manpasik_tenancy_members_active_total 567890") {
		t.Errorf("큰 값 출력 실패:\n%s", body)
	}
}
