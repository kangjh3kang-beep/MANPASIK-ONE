package backup_test

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/manpasik/backend/shared/ops/backup"
)

// fakeRunner 는 stdoutPath 에 정해진 내용을 기록하고 성공/실패 반환.
type fakeRunner struct {
	content   string
	failWith  error
	gotCmd    string
	gotArgs   []string
	gotEnv    map[string]string
}

func (f *fakeRunner) Run(_ context.Context, cmd string, args []string, env map[string]string, stdoutPath string) error {
	f.gotCmd = cmd
	f.gotArgs = args
	f.gotEnv = env
	if f.failWith != nil {
		return f.failWith
	}
	return os.WriteFile(stdoutPath, []byte(f.content), 0o600)
}

func TestParsePgDumpSource_Valid(t *testing.T) {
	src, err := backup.ParsePgDumpSource("host=localhost port=5433 user=u dbname=db password=secret")
	if err != nil {
		t.Fatal(err)
	}
	if src.Host != "localhost" || src.Port != "5433" || src.User != "u" ||
		src.DBName != "db" || src.Password != "secret" {
		t.Errorf("파싱 실패: %+v", src)
	}
}

func TestParsePgDumpSource_DefaultPort(t *testing.T) {
	src, err := backup.ParsePgDumpSource("host=h user=u dbname=d")
	if err != nil {
		t.Fatal(err)
	}
	if src.Port != "5432" {
		t.Errorf("기본 포트 = %q", src.Port)
	}
}

func TestParsePgDumpSource_MissingFields(t *testing.T) {
	if _, err := backup.ParsePgDumpSource("host=h user=u"); err == nil {
		t.Error("dbname 누락에 에러 없음")
	}
}

func TestPgDumpProvider_Backup_Success(t *testing.T) {
	dir := t.TempDir()
	runner := &fakeRunner{content: "-- dump\nCREATE TABLE foo(id int);\n"}
	prov, err := backup.NewPgDumpProvider(backup.PgDumpConfig{OutputDir: dir}, runner, nil)
	if err != nil {
		t.Fatal(err)
	}

	rec, err := prov.Backup(context.Background(),
		backup.KindPostgres,
		"host=h user=u dbname=manpasik password=p",
		"")
	if err != nil {
		t.Fatalf("Backup err = %v", err)
	}
	if rec.Status != backup.StatusCompleted {
		t.Errorf("Status = %s", rec.Status)
	}
	if rec.SizeBytes == 0 {
		t.Error("SizeBytes = 0")
	}
	if !strings.HasPrefix(rec.Checksum, "sha256:") {
		t.Errorf("Checksum = %s", rec.Checksum)
	}
	if rec.CompletedAt == nil {
		t.Error("CompletedAt nil")
	}

	// runner 에 PGPASSWORD 가 전달되어야 함
	if runner.gotEnv["PGPASSWORD"] != "p" {
		t.Errorf("PGPASSWORD = %q", runner.gotEnv["PGPASSWORD"])
	}
	// args 에 --no-owner --no-privileges 기본 포함
	args := strings.Join(runner.gotArgs, " ")
	if !strings.Contains(args, "--no-owner") || !strings.Contains(args, "--no-privileges") {
		t.Errorf("args = %s", args)
	}
}

func TestPgDumpProvider_Backup_RunnerFails(t *testing.T) {
	dir := t.TempDir()
	runner := &fakeRunner{failWith: errors.New("pg_dump: connection refused")}
	prov, _ := backup.NewPgDumpProvider(backup.PgDumpConfig{OutputDir: dir}, runner, nil)

	rec, err := prov.Backup(context.Background(),
		backup.KindPostgres, "host=h user=u dbname=d", "")
	if err == nil {
		t.Fatal("실패에 에러 없음")
	}
	if rec.Status != backup.StatusFailed {
		t.Errorf("Status = %s", rec.Status)
	}
	if rec.Error == "" {
		t.Error("Error 미기록")
	}
}

func TestPgDumpProvider_WrongKind(t *testing.T) {
	dir := t.TempDir()
	prov, _ := backup.NewPgDumpProvider(backup.PgDumpConfig{OutputDir: dir}, &fakeRunner{}, nil)
	if _, err := prov.Backup(context.Background(), backup.KindMilvus, "x", ""); err == nil {
		t.Error("Milvus kind 통과")
	}
}

func TestPgDumpProvider_BackupWithStorage(t *testing.T) {
	dir := t.TempDir()
	storageDir := t.TempDir()
	storage, err := backup.NewLocalFileStorage(storageDir)
	if err != nil {
		t.Fatal(err)
	}
	runner := &fakeRunner{content: "-- dump"}
	prov, _ := backup.NewPgDumpProvider(backup.PgDumpConfig{OutputDir: dir}, runner, storage)

	rec, err := prov.Backup(context.Background(),
		backup.KindPostgres, "host=h user=u dbname=mydb",
		"daily/mydb-1.dump")
	if err != nil {
		t.Fatal(err)
	}
	if rec.Status != backup.StatusCompleted {
		t.Errorf("Status = %s", rec.Status)
	}

	// 스토리지에 파일이 존재해야 함
	if _, statErr := os.Stat(filepath.Join(storageDir, "daily/mydb-1.dump")); statErr != nil {
		t.Errorf("스토리지 업로드 실패: %v", statErr)
	}
}

func TestPgDumpProvider_FormatFlag(t *testing.T) {
	dir := t.TempDir()
	for _, fmt := range []string{"plain", "custom", "tar"} {
		runner := &fakeRunner{content: "x"}
		prov, _ := backup.NewPgDumpProvider(backup.PgDumpConfig{
			OutputDir: dir, Format: fmt,
		}, runner, nil)
		_, _ = prov.Backup(context.Background(),
			backup.KindPostgres, "host=h user=u dbname=d", "")

		args := strings.Join(runner.gotArgs, " ")
		if !strings.Contains(args, "-F") {
			t.Errorf("[%s] -F 없음", fmt)
		}
	}
}

func TestPgDumpProvider_Restore_DryRun(t *testing.T) {
	dir := t.TempDir()
	runner := &fakeRunner{content: "x"}
	prov, _ := backup.NewPgDumpProvider(backup.PgDumpConfig{OutputDir: dir}, runner, nil)

	rec, _ := prov.Backup(context.Background(),
		backup.KindPostgres, "host=h user=u dbname=d", "")
	if err := prov.Restore(context.Background(), &backup.RestoreRequest{
		BackupID: rec.ID, DryRun: true,
	}); err != nil {
		t.Errorf("dry-run err = %v", err)
	}
}

func TestPgDumpProvider_Restore_NonDryRunRefused(t *testing.T) {
	dir := t.TempDir()
	runner := &fakeRunner{content: "x"}
	prov, _ := backup.NewPgDumpProvider(backup.PgDumpConfig{OutputDir: dir}, runner, nil)
	rec, _ := prov.Backup(context.Background(),
		backup.KindPostgres, "host=h user=u dbname=d", "")
	if err := prov.Restore(context.Background(), &backup.RestoreRequest{
		BackupID: rec.ID, DryRun: false,
	}); err == nil {
		t.Error("실 복구가 거부되지 않음")
	}
}

func TestPgDumpProvider_VerifyAndList(t *testing.T) {
	dir := t.TempDir()
	prov, _ := backup.NewPgDumpProvider(backup.PgDumpConfig{OutputDir: dir},
		&fakeRunner{content: "x"}, nil)
	rec, _ := prov.Backup(context.Background(),
		backup.KindPostgres, "host=h user=u dbname=d", "")
	if err := prov.Verify(context.Background(), rec.ID); err != nil {
		t.Errorf("Verify err = %v", err)
	}
	list, _ := prov.List(context.Background(), backup.KindPostgres)
	if len(list) != 1 || list[0].Status != backup.StatusVerified {
		t.Errorf("List = %v", list)
	}
}

func TestPgDumpProvider_Delete(t *testing.T) {
	dir := t.TempDir()
	prov, _ := backup.NewPgDumpProvider(backup.PgDumpConfig{OutputDir: dir},
		&fakeRunner{content: "x"}, nil)
	rec, _ := prov.Backup(context.Background(),
		backup.KindPostgres, "host=h user=u dbname=d", "")
	if err := prov.Delete(context.Background(), rec.ID); err != nil {
		t.Errorf("Delete err = %v", err)
	}
	list, _ := prov.List(context.Background(), backup.KindPostgres)
	if len(list) != 0 {
		t.Errorf("Delete 후 목록 = %v", list)
	}
}

func TestPgDumpProvider_NoOutputDir(t *testing.T) {
	if _, err := backup.NewPgDumpProvider(backup.PgDumpConfig{}, nil, nil); err == nil {
		t.Error("OutputDir 누락에 에러 없음")
	}
}

func TestPgDumpProvider_ProviderName(t *testing.T) {
	dir := t.TempDir()
	prov, _ := backup.NewPgDumpProvider(backup.PgDumpConfig{OutputDir: dir}, nil, nil)
	if prov.Provider() != "pg_dump" {
		t.Errorf("Provider = %q", prov.Provider())
	}
}

func TestLocalFileStorage_PutGetDelete(t *testing.T) {
	root := t.TempDir()
	s, err := backup.NewLocalFileStorage(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Put(context.Background(), "a/b.txt", strings.NewReader("hello")); err != nil {
		t.Fatal(err)
	}
	rc, err := s.Get(context.Background(), "a/b.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer rc.Close()
	body, _ := io.ReadAll(rc)
	if string(body) != "hello" {
		t.Errorf("body = %q", body)
	}
	if err := s.Delete(context.Background(), "a/b.txt"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Get(context.Background(), "a/b.txt"); err == nil {
		t.Error("삭제 후에도 Get 성공")
	}
}

func TestLocalFileStorage_List(t *testing.T) {
	root := t.TempDir()
	s, _ := backup.NewLocalFileStorage(root)
	for _, k := range []string{"daily/x.dump", "daily/y.dump", "weekly/z.dump"} {
		_ = s.Put(context.Background(), k, strings.NewReader(k))
	}
	keys, err := s.List(context.Background(), "daily/")
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 2 {
		t.Errorf("daily 목록 = %v", keys)
	}
}

func TestNoopStorage(t *testing.T) {
	s := backup.NewNoopStorage()
	if err := s.Put(context.Background(), "k", strings.NewReader("data")); err != nil {
		t.Errorf("Put err = %v", err)
	}
	if _, err := s.Get(context.Background(), "k"); err == nil {
		t.Error("Noop Get 이 성공함")
	}
	if err := s.Delete(context.Background(), "k"); err != nil {
		t.Error(err)
	}
	keys, _ := s.List(context.Background(), "")
	if len(keys) != 0 {
		t.Errorf("List = %v", keys)
	}
	if s.Provider() != "noop" {
		t.Errorf("Provider = %q", s.Provider())
	}
}

func TestExecRunner_NoCmd(t *testing.T) {
	if err := (backup.ExecRunner{}).Run(context.Background(), "", nil, nil, "/tmp/x"); err == nil {
		t.Error("빈 cmd 통과")
	}
}
