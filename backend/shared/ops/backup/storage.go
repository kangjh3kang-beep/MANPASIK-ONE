package backup

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// ObjectStorage 는 외부 스토리지(S3/MinIO/로컬) 추상화.
//
// AWS SDK 의존을 피해 io.Reader/Writer 만 사용. S3 어댑터는 별도 패키지에서
// 이 인터페이스를 구현하여 주입한다.
type ObjectStorage interface {
	Put(ctx context.Context, key string, data io.Reader) error
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	List(ctx context.Context, prefix string) ([]string, error)
	Provider() string
}

// LocalFileStorage 는 로컬 디렉토리를 스토리지로 사용 (테스트/개발용).
type LocalFileStorage struct {
	root string
	mu   sync.Mutex
}

// NewLocalFileStorage 는 root 를 베이스 디렉토리로 사용.
func NewLocalFileStorage(root string) (*LocalFileStorage, error) {
	if root == "" {
		return nil, errors.New("root 필수")
	}
	if err := os.MkdirAll(root, 0o750); err != nil {
		return nil, err
	}
	return &LocalFileStorage{root: root}, nil
}

func (s *LocalFileStorage) resolve(key string) string {
	clean := filepath.Clean(strings.TrimPrefix(key, "/"))
	return filepath.Join(s.root, clean)
}

// Put 은 key 위치에 데이터 기록.
func (s *LocalFileStorage) Put(ctx context.Context, key string, data io.Reader) error {
	if key == "" {
		return errors.New("key 필수")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	path := s.resolve(key)
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// ctx 캔슬 시 즉시 종료
	done := make(chan error, 1)
	go func() {
		_, copyErr := io.Copy(f, data)
		done <- copyErr
	}()
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		_ = os.Remove(path)
		return ctx.Err()
	}
}

// Get 은 key 위치의 데이터 반환.
func (s *LocalFileStorage) Get(_ context.Context, key string) (io.ReadCloser, error) {
	path := s.resolve(key)
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("키 없음: %s", key)
		}
		return nil, err
	}
	return f, nil
}

// Delete 는 key 제거 (없어도 에러 아님).
func (s *LocalFileStorage) Delete(_ context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	path := s.resolve(key)
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// List 는 prefix 로 시작하는 모든 키 반환 (root 기준 상대 경로).
func (s *LocalFileStorage) List(_ context.Context, prefix string) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	target := s.resolve(prefix)
	var out []string
	err := filepath.Walk(s.root, func(path string, info os.FileInfo, _ error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		if strings.HasPrefix(path, target) {
			rel, err := filepath.Rel(s.root, path)
			if err == nil {
				out = append(out, filepath.ToSlash(rel))
			}
		}
		return nil
	})
	sort.Strings(out)
	return out, err
}

// Provider 이름.
func (s *LocalFileStorage) Provider() string { return "local" }

// NoopStorage 는 모든 호출이 즉시 성공하지만 실제로는 아무것도 하지 않음.
//
// PgDumpProvider 의 destination="" 이거나 단순히 로컬에만 백업하고 싶을 때 사용.
type NoopStorage struct{}

// NewNoopStorage 생성.
func NewNoopStorage() *NoopStorage { return &NoopStorage{} }

func (NoopStorage) Put(_ context.Context, _ string, data io.Reader) error {
	// data 를 모두 읽고 버림 (드레인 안 하면 호출자 측 reader 가 막힐 수 있음)
	if data != nil {
		_, _ = io.Copy(io.Discard, data)
	}
	return nil
}

func (NoopStorage) Get(_ context.Context, _ string) (io.ReadCloser, error) {
	return nil, errors.New("noop storage: get 미지원")
}

func (NoopStorage) Delete(_ context.Context, _ string) error  { return nil }
func (NoopStorage) List(_ context.Context, _ string) ([]string, error) { return nil, nil }
func (NoopStorage) Provider() string                          { return "noop" }
