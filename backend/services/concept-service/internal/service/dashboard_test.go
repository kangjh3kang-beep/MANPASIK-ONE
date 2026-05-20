package service_test

import (
	"testing"
	"time"

	"github.com/manpasik/backend/services/concept-service/internal/service"
)

// TestStatsEngine_RecordMeasurement는 측정 기록을 검증합니다.
func TestStatsEngine_RecordMeasurement(t *testing.T) {
	e := service.NewStatsEngine()
	e.RecordMeasurement("c1", false, time.Now().UTC())
	e.RecordMeasurement("c1", true, time.Now().UTC())

	c := &service.Concept{ID: "c1", Name: "혈당관리", Category: "medical"}
	stats := e.GetConceptStats(c)

	if stats.MeasurementCount != 2 {
		t.Errorf("MeasurementCount = %d, want 2", stats.MeasurementCount)
	}
	if stats.AbnormalCount != 1 {
		t.Errorf("AbnormalCount = %d, want 1", stats.AbnormalCount)
	}
	if stats.LastActivityAt == nil {
		t.Error("LastActivityAt 미설정")
	}
}

// TestStatsEngine_GetTrends는 일별 트렌드 조회를 검증합니다.
func TestStatsEngine_GetTrends(t *testing.T) {
	e := service.NewStatsEngine()
	now := time.Now().UTC()
	e.RecordMeasurement("c1", false, now)
	e.RecordMeasurement("c1", true, now.AddDate(0, 0, -1))
	e.RecordMeasurement("c1", false, now.AddDate(0, 0, -2))

	trends := e.GetTrends("c1", 5)
	if len(trends) != 3 {
		t.Errorf("trends = %d, want 3", len(trends))
	}

	// 내림차순 정렬 (가장 최근이 첫 번째)
	if trends[0].Date.Before(trends[1].Date) {
		t.Error("트렌드가 내림차순으로 정렬되지 않음")
	}
}

// TestStatsEngine_GetOrgKPI는 조직 KPI 계산을 검증합니다.
func TestStatsEngine_GetOrgKPI(t *testing.T) {
	e := service.NewStatsEngine()
	now := time.Now().UTC()
	e.RecordMeasurement("c1", false, now)
	e.RecordMeasurement("c1", true, now)
	e.RecordMeasurement("c2", false, now)

	org := &service.Organization{
		ID:        "org1",
		Name:      "건강병원",
		MemberIDs: []string{"u1", "u2", "u3"},
	}

	kpi := e.GetOrgKPI(org, 5)
	if kpi.ActiveUsers != 3 {
		t.Errorf("ActiveUsers = %d, want 3", kpi.ActiveUsers)
	}
	if kpi.TotalConcepts != 5 {
		t.Errorf("TotalConcepts = %d, want 5", kpi.TotalConcepts)
	}
	if kpi.MeasurementsToday != 3 {
		t.Errorf("MeasurementsToday = %d, want 3", kpi.MeasurementsToday)
	}
	expectedRate := 1.0 / 3.0
	if kpi.AbnormalRate < expectedRate-0.01 || kpi.AbnormalRate > expectedRate+0.01 {
		t.Errorf("AbnormalRate = %f, want %f", kpi.AbnormalRate, expectedRate)
	}
}

// TestStatsEngine_AvgUsagePerDay는 일평균 사용량을 검증합니다.
func TestStatsEngine_AvgUsagePerDay(t *testing.T) {
	e := service.NewStatsEngine()
	now := time.Now().UTC()
	e.RecordMeasurement("c1", false, now)
	e.RecordMeasurement("c1", false, now.AddDate(0, 0, -1))
	e.RecordMeasurement("c1", false, now.AddDate(0, 0, -2))

	c := &service.Concept{ID: "c1"}
	stats := e.GetConceptStats(c)
	if stats.AvgUsagePerDay != 1.0 {
		t.Errorf("AvgUsagePerDay = %f, want 1.0", stats.AvgUsagePerDay)
	}
}

// TestHasPermission_Owner는 owner의 모든 권한을 검증합니다.
func TestHasPermission_Owner(t *testing.T) {
	actions := []string{"read", "write", "delete", "manage", "invite"}
	for _, a := range actions {
		if !service.HasPermission(service.RoleOwner, a) {
			t.Errorf("Owner는 %q 권한 보유해야 함", a)
		}
	}
}

// TestHasPermission_Guest는 guest의 read 권한만 보유함을 검증합니다.
func TestHasPermission_Guest(t *testing.T) {
	if !service.HasPermission(service.RoleGuest, "read") {
		t.Error("Guest는 read 권한 보유해야 함")
	}
	if service.HasPermission(service.RoleGuest, "write") {
		t.Error("Guest는 write 권한 없어야 함")
	}
	if service.HasPermission(service.RoleGuest, "delete") {
		t.Error("Guest는 delete 권한 없어야 함")
	}
}

// TestIsHigherOrEqualRole은 역할 순위 비교를 검증합니다.
func TestIsHigherOrEqualRole(t *testing.T) {
	cases := []struct {
		role1, role2 string
		expected     bool
	}{
		{service.RoleOwner, service.RoleAdmin, true},
		{service.RoleAdmin, service.RoleOwner, false},
		{service.RoleAdmin, service.RoleAdmin, true},
		{service.RoleMember, service.RoleGuest, true},
		{service.RoleGuest, service.RoleMember, false},
	}

	for _, c := range cases {
		got := service.IsHigherOrEqualRole(c.role1, c.role2)
		if got != c.expected {
			t.Errorf("IsHigherOrEqualRole(%s, %s) = %v, want %v",
				c.role1, c.role2, got, c.expected)
		}
	}
}

// TestInvitationManager_SendAccept는 초대 발송과 수락을 검증합니다.
func TestInvitationManager_SendAccept(t *testing.T) {
	m := service.NewInvitationManager()

	inv := &service.Invitation{
		ID:             "inv-1",
		OrganizationID: "org-1",
		InviterID:      "user-inviter",
		InviteeEmail:   "test@example.com",
		Role:           service.RoleMember,
	}

	if err := m.Send(inv); err != nil {
		t.Fatalf("Send 실패: %v", err)
	}

	got, _ := m.Get("inv-1")
	if got.Status != service.InviteStatusPending {
		t.Errorf("Status = %q, want pending", got.Status)
	}
	if got.ExpiresAt.IsZero() {
		t.Error("ExpiresAt 미설정")
	}

	if err := m.Accept("inv-1", "user-invitee"); err != nil {
		t.Fatalf("Accept 실패: %v", err)
	}

	got, _ = m.Get("inv-1")
	if got.Status != service.InviteStatusAccepted {
		t.Errorf("Status = %q, want accepted", got.Status)
	}
	if got.InviteeUserID != "user-invitee" {
		t.Errorf("InviteeUserID = %q, want user-invitee", got.InviteeUserID)
	}
}

// TestInvitationManager_Decline는 초대 거절을 검증합니다.
func TestInvitationManager_Decline(t *testing.T) {
	m := service.NewInvitationManager()
	inv := &service.Invitation{ID: "i", InviterID: "u", InviteeEmail: "t@a.com"}
	_ = m.Send(inv)

	if err := m.Decline("i"); err != nil {
		t.Fatalf("Decline 실패: %v", err)
	}
	got, _ := m.Get("i")
	if got.Status != service.InviteStatusDeclined {
		t.Errorf("Status = %q, want declined", got.Status)
	}
}

// TestInvitationManager_Revoke는 초대 회수를 검증합니다.
func TestInvitationManager_Revoke(t *testing.T) {
	m := service.NewInvitationManager()
	inv := &service.Invitation{ID: "i", InviterID: "u-inv", InviteeEmail: "t@a.com"}
	_ = m.Send(inv)

	if err := m.Revoke("i", "u-inv"); err != nil {
		t.Fatalf("Revoke 실패: %v", err)
	}
	got, _ := m.Get("i")
	if got.Status != service.InviteStatusRevoked {
		t.Errorf("Status = %q, want revoked", got.Status)
	}

	// 다른 사람은 revoke 불가
	inv2 := &service.Invitation{ID: "i2", InviterID: "u-inv", InviteeEmail: "t@a.com"}
	_ = m.Send(inv2)
	if err := m.Revoke("i2", "other-user"); err == nil {
		t.Error("발신자가 아닌 사람이 revoke함")
	}
}

// TestInvitationManager_ExpireOld는 만료 처리를 검증합니다.
func TestInvitationManager_ExpireOld(t *testing.T) {
	m := service.NewInvitationManager()

	inv := &service.Invitation{
		ID:           "exp",
		InviterID:    "u",
		InviteeEmail: "t@a.com",
		ExpiresAt:    time.Now().UTC().Add(-1 * time.Hour), // 이미 만료
	}
	_ = m.Send(inv)

	// Send에서 ExpiresAt이 zero가 아니므로 그대로 유지되어야 하지만
	// Send 내부에서 7일 만료로 덮어씌우는지 확인
	got, _ := m.Get("exp")
	// 실제로는 Send 내부에서 ExpiresAt이 IsZero일 때만 덮어쓰므로 과거 시각 유지됨
	if got.ExpiresAt.After(time.Now().UTC()) {
		t.Errorf("ExpiresAt이 미래로 변경됨: %v", got.ExpiresAt)
	}

	expired := m.ExpireOld()
	if expired != 1 {
		t.Errorf("expired = %d, want 1", expired)
	}

	got, _ = m.Get("exp")
	if got.Status != service.InviteStatusExpired {
		t.Errorf("Status = %q, want expired", got.Status)
	}
}

// TestInvitationManager_AcceptExpired는 만료된 초대 수락 차단을 검증합니다.
func TestInvitationManager_AcceptExpired(t *testing.T) {
	m := service.NewInvitationManager()
	inv := &service.Invitation{
		ID: "exp", InviterID: "u", InviteeEmail: "t@a.com",
		ExpiresAt: time.Now().UTC().Add(-1 * time.Hour),
	}
	_ = m.Send(inv)

	err := m.Accept("exp", "user")
	if err == nil {
		t.Error("만료된 초대 수락이 허용됨")
	}
}

// TestInvitationManager_InvalidRole는 잘못된 역할 거부를 검증합니다.
func TestInvitationManager_InvalidRole(t *testing.T) {
	m := service.NewInvitationManager()
	inv := &service.Invitation{
		ID: "i", InviterID: "u", InviteeEmail: "t@a.com", Role: "invalid_role",
	}
	if err := m.Send(inv); err == nil {
		t.Error("잘못된 역할 초대가 허용됨")
	}
}

// TestInvitationManager_ListByOrganization는 조직별 초대 목록 조회를 검증합니다.
func TestInvitationManager_ListByOrganization(t *testing.T) {
	m := service.NewInvitationManager()
	_ = m.Send(&service.Invitation{ID: "1", OrganizationID: "org-A", InviterID: "u", InviteeEmail: "a@x.com"})
	_ = m.Send(&service.Invitation{ID: "2", OrganizationID: "org-A", InviterID: "u", InviteeEmail: "b@x.com"})
	_ = m.Send(&service.Invitation{ID: "3", OrganizationID: "org-B", InviterID: "u", InviteeEmail: "c@x.com"})

	orgA := m.ListByOrganization("org-A")
	if len(orgA) != 2 {
		t.Errorf("org-A 초대 수 = %d, want 2", len(orgA))
	}
}

// TestInvitationManager_CountByStatus는 상태별 통계를 검증합니다.
func TestInvitationManager_CountByStatus(t *testing.T) {
	m := service.NewInvitationManager()
	_ = m.Send(&service.Invitation{ID: "1", InviterID: "u", InviteeEmail: "a@x.com"})
	_ = m.Send(&service.Invitation{ID: "2", InviterID: "u", InviteeEmail: "b@x.com"})
	_ = m.Decline("2")

	counts := m.CountByStatus()
	if counts[service.InviteStatusPending] != 1 {
		t.Errorf("pending = %d, want 1", counts[service.InviteStatusPending])
	}
	if counts[service.InviteStatusDeclined] != 1 {
		t.Errorf("declined = %d, want 1", counts[service.InviteStatusDeclined])
	}
}
