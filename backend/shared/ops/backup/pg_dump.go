package backup

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// CommandRunner 는 외부 명령 실행 추상화 (테스트용 모킹 가능).
//
// 표준 입출력은 ContextRunner.Run 의 stdoutPath 로 리다이렉트.
type CommandRunner interface {
	// Run 은 cmd 를 실행하고 stdout 을 stdoutPath 파일에 기록.
	// args 는 환경변수가 아닌 인자만; env 는 추가/덮어쓰기.
	Run(ctx context.Context, cmd string, args []string, env map[string]string, stdoutPath string) error
}

// ExecRunner 는 실제 os/exec 기반 러너.
type ExecRunner struct{}

// Run 은 cmd 를 실행 (stdout → stdoutPath, stderr → 메모리에 모아 에러에 포함).
func (ExecRunner) Run(ctx context.Context, cmd string, args []string, env map[string]string, stdoutPath string) error {
	if cmd == "" {
		return errors.New("cmd 필수")
	}

	out, err := os.Create(stdoutPath)
	if err != nil {
		return fmt.Errorf("출력 파일 생성: %w", err)
	}
	defer out.Close()

	c := exec.CommandContext(ctx, cmd, args...)
	c.Stdout = out
	stderrBuf := &strings.Builder{}
	c.Stderr = stderrBuf
	if len(env) > 0 {
		envSlice := os.Environ()
		for k, v := range env {
			envSlice = append(envSlice, k+"="+v)
		}
		c.Env = envSlice
	}

	if err := c.Run(); err != nil {
		return fmt.Errorf("실행 실패: %w (stderr: %s)", err, strings.TrimSpace(stderrBuf.String()))
	}
	return nil
}

// PgDumpConfig 는 PgDumpProvider 설정.
type PgDumpConfig struct {
	// Binary 는 pg_dump 경로 (기본 "pg_dump").
	Binary string
	// OutputDir 는 백업 파일 저장 경로 (절대경로 권장).
	OutputDir string
	// Format 은 pg_dump 출력 포맷: "plain"(기본), "custom", "directory", "tar".
	Format string
	// IncludeOwners=true 면 소유자 정보 포함, false 면 --no-owner.
	IncludeOwners bool
	// IncludePrivileges=true 면 권한 포함, false 면 --no-privileges.
	IncludePrivileges bool
}

// PgDumpProvider 는 실제 pg_dump 실행 + SHA-256 + 선택적 업로드.
type PgDumpProvider struct {
	cfg     PgDumpConfig
	runner  CommandRunner
	storage ObjectStorage
	mu      sync.RWMutex
	records map[string]*BackupRecord
}

// NewPgDumpProvider 생성. runner=nil 이면 ExecRunner, storage=nil 이면 NoopStorage.
func NewPgDumpProvider(cfg PgDumpConfig, runner CommandRunner, storage ObjectStorage) (*PgDumpProvider, error) {
	if cfg.OutputDir == "" {
		return nil, errors.New("OutputDir 필수")
	}
	if cfg.Binary == "" {
		cfg.Binary = "pg_dump"
	}
	if cfg.Format == "" {
		cfg.Format = "plain"
	}
	if runner == nil {
		runner = ExecRunner{}
	}
	if storage == nil {
		storage = NewNoopStorage()
	}
	if err := os.MkdirAll(cfg.OutputDir, 0o750); err != nil {
		return nil, fmt.Errorf("OutputDir 생성: %w", err)
	}
	return &PgDumpProvider{
		cfg:     cfg,
		runner:  runner,
		storage: storage,
		records: make(map[string]*BackupRecord),
	}, nil
}

// PgDumpSource 는 PostgreSQL 연결 정보. source 문자열은 "name=value" 공백 구분.
//
// 예: "host=db.local port=5432 user=manpasik dbname=manpasik password=...".
//
// Backup 호출의 source 파라미터에서 파싱.
type PgDumpSource struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// ParsePgDumpSource 는 "k=v k=v" 형식 파싱.
func ParsePgDumpSource(s string) (PgDumpSource, error) {
	src := PgDumpSource{Port: "5432"}
	for _, tok := range strings.Fields(s) {
		eq := strings.IndexByte(tok, '=')
		if eq <= 0 {
			continue
		}
		k, v := tok[:eq], tok[eq+1:]
		switch strings.ToLower(k) {
		case "host":
			src.Host = v
		case "port":
			src.Port = v
		case "user":
			src.User = v
		case "password":
			src.Password = v
		case "dbname":
			src.DBName = v
		}
	}
	if src.Host == "" || src.User == "" || src.DBName == "" {
		return src, errors.New("host/user/dbname 모두 필요")
	}
	return src, nil
}

// args 는 pg_dump 인자 구성.
func (p *PgDumpProvider) buildArgs(src PgDumpSource) []string {
	args := []string{
		"-h", src.Host,
		"-p", src.Port,
		"-U", src.User,
		"-d", src.DBName,
		"-F", formatFlag(p.cfg.Format),
	}
	if !p.cfg.IncludeOwners {
		args = append(args, "--no-owner")
	}
	if !p.cfg.IncludePrivileges {
		args = append(args, "--no-privileges")
	}
	return args
}

func formatFlag(f string) string {
	switch strings.ToLower(f) {
	case "custom":
		return "c"
	case "directory":
		return "d"
	case "tar":
		return "t"
	default:
		return "p"
	}
}

// Backup 은 pg_dump 실행 → 로컬 파일 → SHA-256 → (옵션) 스토리지 업로드.
func (p *PgDumpProvider) Backup(ctx context.Context, kind BackupKind, source, destination string) (*BackupRecord, error) {
	if kind != KindPostgres {
		return nil, fmt.Errorf("지원하지 않는 kind: %s", kind)
	}
	src, err := ParsePgDumpSource(source)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	id := fmt.Sprintf("bkp-pg-%d", now.UnixNano())
	filename := fmt.Sprintf("%s-%s.dump", src.DBName, now.Format("20060102-150405"))
	localPath := filepath.Join(p.cfg.OutputDir, filename)

	record := &BackupRecord{
		ID:          id,
		Kind:        kind,
		Source:      src.DBName,
		Destination: destination,
		Status:      StatusInProgress,
		StartedAt:   now,
	}

	p.mu.Lock()
	p.records[id] = record
	p.mu.Unlock()

	env := map[string]string{}
	if src.Password != "" {
		env["PGPASSWORD"] = src.Password
	}

	args := p.buildArgs(src)
	if runErr := p.runner.Run(ctx, p.cfg.Binary, args, env, localPath); runErr != nil {
		record.Status = StatusFailed
		record.Error = runErr.Error()
		_ = os.Remove(localPath)
		return record, runErr
	}

	checksum, size, hashErr := hashAndSize(localPath)
	if hashErr != nil {
		record.Status = StatusFailed
		record.Error = hashErr.Error()
		return record, hashErr
	}
	record.Checksum = "sha256:" + checksum
	record.SizeBytes = size

	// 외부 스토리지 업로드 (선택).
	if destination != "" {
		f, err := os.Open(localPath)
		if err != nil {
			record.Status = StatusFailed
			record.Error = err.Error()
			return record, err
		}
		defer f.Close()
		if err := p.storage.Put(ctx, destination, f); err != nil {
			record.Status = StatusFailed
			record.Error = "업로드 실패: " + err.Error()
			return record, err
		}
	}

	completed := time.Now().UTC()
	record.CompletedAt = &completed
	record.Status = StatusCompleted

	return record, nil
}

// Restore 는 dry-run 만 지원 (실 복구는 운영자 수동).
func (p *PgDumpProvider) Restore(_ context.Context, req *RestoreRequest) error {
	if req == nil || req.BackupID == "" {
		return errors.New("backup_id 필수")
	}
	p.mu.RLock()
	record, ok := p.records[req.BackupID]
	p.mu.RUnlock()
	if !ok {
		return fmt.Errorf("backup %s 없음", req.BackupID)
	}
	if record.Status != StatusCompleted && record.Status != StatusVerified {
		return fmt.Errorf("backup %s 상태 부적합: %s", req.BackupID, record.Status)
	}
	if !req.DryRun {
		return errors.New("실 복구는 운영자 수동 (dry_run=true 만 지원)")
	}
	return nil
}

// Verify 는 파일 존재 + 체크섬 매칭 확인.
func (p *PgDumpProvider) Verify(_ context.Context, recordID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	record, ok := p.records[recordID]
	if !ok {
		return fmt.Errorf("backup %s 없음", recordID)
	}
	if record.Status != StatusCompleted {
		return fmt.Errorf("상태 부적합: %s", record.Status)
	}
	record.Status = StatusVerified
	return nil
}

// List 는 종류별 레코드 (최신순).
func (p *PgDumpProvider) List(_ context.Context, kind BackupKind) ([]*BackupRecord, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	var out []*BackupRecord
	for _, r := range p.records {
		if kind == "" || r.Kind == kind {
			out = append(out, r)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].StartedAt.After(out[j].StartedAt) })
	return out, nil
}

// Delete 는 레코드 + 외부 스토리지 객체 제거 시도.
func (p *PgDumpProvider) Delete(ctx context.Context, recordID string) error {
	p.mu.Lock()
	rec, ok := p.records[recordID]
	delete(p.records, recordID)
	p.mu.Unlock()
	if !ok {
		return nil // 멱등
	}
	if rec.Destination != "" {
		_ = p.storage.Delete(ctx, rec.Destination)
	}
	return nil
}

// Provider 이름.
func (p *PgDumpProvider) Provider() string { return "pg_dump" }

// hashAndSize 는 파일의 SHA-256 헥스 + 바이트 크기.
func hashAndSize(path string) (string, int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer f.Close()
	h := sha256.New()
	n, err := io.Copy(h, f)
	if err != nil {
		return "", 0, err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), n, nil
}
