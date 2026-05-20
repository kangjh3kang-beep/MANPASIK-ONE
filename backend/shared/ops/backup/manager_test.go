package backup_test

import (
	"context"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/ops/backup"
)

func TestMemoryProvider_BackupSuccess(t *testing.T) {
	p := backup.NewMemoryProvider()
	rec, err := p.Backup(context.Background(), backup.KindPostgres, "manpasik_db", "/backup/postgres/")
	if err != nil {
		t.Fatalf("Backup 실패: %v", err)
	}
	if rec.Status != backup.StatusCompleted {
		t.Errorf("Status = %q", rec.Status)
	}
	if rec.Checksum == "" {
		t.Error("Checksum 미생성")
	}
	if rec.SizeBytes <= 0 {
		t.Error("SizeBytes <= 0")
	}
}

func TestMemoryProvider_BackupRequiresSource(t *testing.T) {
	p := backup.NewMemoryProvider()
	_, err := p.Backup(context.Background(), backup.KindPostgres, "", "/dest/")
	if err == nil {
		t.Error("source 없이 통과됨")
	}
}

func TestMemoryProvider_RestoreHappyPath(t *testing.T) {
	p := backup.NewMemoryProvider()
	rec, _ := p.Backup(context.Background(), backup.KindPostgres, "db", "/dest/")

	err := p.Restore(context.Background(), &backup.RestoreRequest{
		BackupID: rec.ID,
		TargetDB: "manpasik_db_restored",
	})
	if err != nil {
		t.Errorf("Restore 실패: %v", err)
	}
}

func TestMemoryProvider_RestoreDryRun(t *testing.T) {
	p := backup.NewMemoryProvider()
	rec, _ := p.Backup(context.Background(), backup.KindPostgres, "db", "/dest/")

	err := p.Restore(context.Background(), &backup.RestoreRequest{
		BackupID: rec.ID,
		TargetDB: "test",
		DryRun:   true,
	})
	if err != nil {
		t.Errorf("DryRun 실패: %v", err)
	}
}

func TestMemoryProvider_RestoreNotFound(t *testing.T) {
	p := backup.NewMemoryProvider()
	err := p.Restore(context.Background(), &backup.RestoreRequest{BackupID: "missing"})
	if err == nil {
		t.Error("미존재 백업 복구가 통과됨")
	}
}

func TestMemoryProvider_Verify(t *testing.T) {
	p := backup.NewMemoryProvider()
	rec, _ := p.Backup(context.Background(), backup.KindPostgres, "db", "/dest/")

	if err := p.Verify(context.Background(), rec.ID); err != nil {
		t.Errorf("Verify 실패: %v", err)
	}

	list, _ := p.List(context.Background(), backup.KindPostgres)
	if len(list) == 0 || list[0].Status != backup.StatusVerified {
		t.Errorf("Status not verified")
	}
}

func TestMemoryProvider_List(t *testing.T) {
	p := backup.NewMemoryProvider()
	for i := 0; i < 3; i++ {
		_, _ = p.Backup(context.Background(), backup.KindPostgres, "db", "/x/")
	}
	for i := 0; i < 2; i++ {
		_, _ = p.Backup(context.Background(), backup.KindMilvus, "vectors", "/y/")
	}

	pgList, _ := p.List(context.Background(), backup.KindPostgres)
	if len(pgList) != 3 {
		t.Errorf("postgres = %d, want 3", len(pgList))
	}
	all, _ := p.List(context.Background(), "")
	if len(all) != 5 {
		t.Errorf("all = %d, want 5", len(all))
	}
}

func TestMemoryProvider_Delete(t *testing.T) {
	p := backup.NewMemoryProvider()
	rec, _ := p.Backup(context.Background(), backup.KindPostgres, "db", "/x/")

	if err := p.Delete(context.Background(), rec.ID); err != nil {
		t.Errorf("Delete 실패: %v", err)
	}
	if p.Count() != 0 {
		t.Errorf("Count = %d, want 0", p.Count())
	}
}

func TestRetentionEnforcer_DeletesOldBackups(t *testing.T) {
	p := backup.NewMemoryProvider()
	enforcer := backup.NewRetentionEnforcer(p)

	// 신규 + 오래된 백업 혼합
	for i := 0; i < 3; i++ {
		_, _ = p.Backup(context.Background(), backup.KindPostgres, "db", "/x/")
	}

	// 보관 0일 = 모두 삭제
	deleted, _ := enforcer.Enforce(context.Background(), backup.KindPostgres, 0)
	if deleted != 0 {
		t.Errorf("0일 보관은 삭제 안함, got %d", deleted)
	}

	// 보관 -1일 처리: 0 이하는 무시
	if d, _ := enforcer.Enforce(context.Background(), backup.KindPostgres, -1); d != 0 {
		t.Errorf("-1일도 삭제 안함, got %d", d)
	}
}

func TestScheduleManager_AddGet(t *testing.T) {
	m := backup.NewScheduleManager()
	s := &backup.Schedule{
		Name:      "daily-pg",
		Kind:      backup.KindPostgres,
		Source:    "manpasik_db",
		Cron:      "0 2 * * *",
		Retention: 30,
		Enabled:   true,
	}
	if err := m.Add(s); err != nil {
		t.Fatalf("Add 실패: %v", err)
	}
	got, err := m.Get("daily-pg")
	if err != nil {
		t.Fatalf("Get 실패: %v", err)
	}
	if got.Cron != "0 2 * * *" {
		t.Errorf("Cron = %q", got.Cron)
	}
}

func TestScheduleManager_AddRequiresName(t *testing.T) {
	m := backup.NewScheduleManager()
	if err := m.Add(&backup.Schedule{Cron: "* * * * *"}); err == nil {
		t.Error("이름 없이 통과됨")
	}
}

func TestScheduleManager_Remove(t *testing.T) {
	m := backup.NewScheduleManager()
	_ = m.Add(&backup.Schedule{Name: "x", Cron: "* * * * *"})
	m.Remove("x")
	if _, err := m.Get("x"); err == nil {
		t.Error("Remove 후에도 조회됨")
	}
}

func TestScheduleManager_EnabledSchedules(t *testing.T) {
	m := backup.NewScheduleManager()
	_ = m.Add(&backup.Schedule{Name: "a", Cron: "* * * * *", Enabled: true})
	_ = m.Add(&backup.Schedule{Name: "b", Cron: "* * * * *", Enabled: false})
	_ = m.Add(&backup.Schedule{Name: "c", Cron: "* * * * *", Enabled: true})

	enabled := m.EnabledSchedules()
	if len(enabled) != 2 {
		t.Errorf("enabled = %d, want 2", len(enabled))
	}
}

func TestScheduleManager_SetEnabled(t *testing.T) {
	m := backup.NewScheduleManager()
	_ = m.Add(&backup.Schedule{Name: "x", Cron: "* * * * *", Enabled: false})

	if err := m.SetEnabled("x", true); err != nil {
		t.Fatalf("SetEnabled 실패: %v", err)
	}
	got, _ := m.Get("x")
	if !got.Enabled {
		t.Error("Enabled가 변경되지 않음")
	}

	if err := m.SetEnabled("missing", true); err == nil {
		t.Error("미존재 스케줄 통과됨")
	}
}

func TestBackupRecord_Status(t *testing.T) {
	now := time.Now().UTC()
	r := &backup.BackupRecord{
		ID:          "x",
		Kind:        backup.KindPostgres,
		Status:      backup.StatusCompleted,
		StartedAt:   now,
		CompletedAt: &now,
	}
	if r.Status != backup.StatusCompleted {
		t.Error("Status mismatch")
	}
}
