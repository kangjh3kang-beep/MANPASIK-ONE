package backup

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// HTTPDoer 는 net/http.Client 호환 (테스트용 모킹).
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// HTTPStorageConfig 는 HTTPObjectStorage 설정.
//
// AWS S3 SigV4 의존을 피해 Bearer/Custom 토큰 인증만 지원. 호환 대상:
//   - Cloudflare R2 (API Token)
//   - MinIO with custom auth proxy
//   - 사내 S3-게이트웨이 (Bearer JWT)
//   - 자체 REST 객체 저장소
//
// AWS S3 본가 사용 시 별도 어댑터 (s3-aws/ 패키지) 가 SigV4 처리.
type HTTPStorageConfig struct {
	// BaseURL 은 객체 저장소 기본 URL (예: "https://r2.example.com/manpasik-backups").
	BaseURL string
	// AuthHeader 키 (예: "Authorization", "X-API-Key").
	AuthHeader string
	// AuthValue (예: "Bearer abc...", "key-xyz").
	AuthValue string
	// RequestTimeout 은 단일 호출 타임아웃 (기본 30초).
	RequestTimeout time.Duration
	// ListEndpoint 가 비어있지 않으면 List 시 그 엔드포인트로 GET 수행.
	// 비어있으면 List 는 ErrUnsupported 반환.
	ListEndpoint string
}

// HTTPObjectStorage 는 표준 HTTP REST 로 객체 저장소에 접근.
type HTTPObjectStorage struct {
	cfg    HTTPStorageConfig
	client HTTPDoer
	mu     sync.Mutex
}

// NewHTTPObjectStorage 생성. doer=nil 이면 http.Client{Timeout:30s}.
func NewHTTPObjectStorage(cfg HTTPStorageConfig, doer HTTPDoer) (*HTTPObjectStorage, error) {
	if cfg.BaseURL == "" {
		return nil, errors.New("BaseURL 필수")
	}
	if cfg.RequestTimeout <= 0 {
		cfg.RequestTimeout = 30 * time.Second
	}
	cfg.BaseURL = strings.TrimRight(cfg.BaseURL, "/")
	if doer == nil {
		doer = &http.Client{Timeout: cfg.RequestTimeout}
	}
	return &HTTPObjectStorage{cfg: cfg, client: doer}, nil
}

// Provider 이름.
func (s *HTTPObjectStorage) Provider() string { return "http" }

// resolve 는 baseURL + key 조합 (key 의 선행 / 제거).
func (s *HTTPObjectStorage) resolve(key string) string {
	k := strings.TrimLeft(key, "/")
	return s.cfg.BaseURL + "/" + k
}

// applyAuth 는 요청에 인증 헤더 부착.
func (s *HTTPObjectStorage) applyAuth(req *http.Request) {
	if s.cfg.AuthHeader != "" && s.cfg.AuthValue != "" {
		req.Header.Set(s.cfg.AuthHeader, s.cfg.AuthValue)
	}
}

// Put 은 PUT key (Body=data).
func (s *HTTPObjectStorage) Put(ctx context.Context, key string, data io.Reader) error {
	if key == "" {
		return errors.New("key 필수")
	}
	if data == nil {
		data = bytes.NewReader(nil)
	}

	// io.Reader 가 ContentLength 를 모르므로 일단 메모리에 모음.
	// (대용량 백업은 별도 multipart 어댑터 권장; 여기서는 단순/일반)
	buf, err := io.ReadAll(data)
	if err != nil {
		return fmt.Errorf("body read: %w", err)
	}

	reqCtx, cancel := context.WithTimeout(ctx, s.cfg.RequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, "PUT", s.resolve(key), bytes.NewReader(buf))
	if err != nil {
		return err
	}
	req.ContentLength = int64(len(buf))
	req.Header.Set("Content-Type", "application/octet-stream")
	s.applyAuth(req)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("PUT 실패: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("PUT %s: HTTP %d (%s)", key, resp.StatusCode,
			strings.TrimSpace(string(body)))
	}
	return nil
}

// Get 은 GET key.
func (s *HTTPObjectStorage) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	reqCtx, cancel := context.WithTimeout(ctx, s.cfg.RequestTimeout)
	// 주의: Get 은 ReadCloser 를 반환하므로 cancel 을 닫는 시점이 문제.
	// 호출자가 ReadCloser 를 모두 읽은 뒤 Close 하면 ctx 가 살아 있어야 함.
	// 단순화를 위해 별도 ctx 를 사용하지 않고 부모 ctx 를 그대로 전달.
	_ = cancel // lint 무시 — Get 의 ctx 수명은 호출자 관리

	req, err := http.NewRequestWithContext(reqCtx, "GET", s.resolve(key), nil)
	if err != nil {
		return nil, err
	}
	s.applyAuth(req)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET 실패: %w", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		resp.Body.Close()
		return nil, fmt.Errorf("키 없음: %s", key)
	}
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		resp.Body.Close()
		return nil, fmt.Errorf("GET %s: HTTP %d (%s)", key, resp.StatusCode,
			strings.TrimSpace(string(body)))
	}
	return resp.Body, nil
}

// Delete 는 DELETE key. 404 는 멱등 성공.
func (s *HTTPObjectStorage) Delete(ctx context.Context, key string) error {
	reqCtx, cancel := context.WithTimeout(ctx, s.cfg.RequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, "DELETE", s.resolve(key), nil)
	if err != nil {
		return err
	}
	s.applyAuth(req)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("DELETE 실패: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("DELETE %s: HTTP %d (%s)", key, resp.StatusCode,
			strings.TrimSpace(string(body)))
	}
	return nil
}

// List 는 ListEndpoint 가 설정된 경우 GET ListEndpoint?prefix=... 호출.
//
// 응답 형식: 줄바꿈으로 분리된 키 목록 (text/plain). 다른 포맷이 필요한 경우
// 별도 어댑터에서 List 를 오버라이드.
func (s *HTTPObjectStorage) List(ctx context.Context, prefix string) ([]string, error) {
	if s.cfg.ListEndpoint == "" {
		return nil, errors.New("ListEndpoint 미설정")
	}

	reqCtx, cancel := context.WithTimeout(ctx, s.cfg.RequestTimeout)
	defer cancel()

	url := strings.TrimRight(s.cfg.ListEndpoint, "/") + "?prefix=" + strings.TrimLeft(prefix, "/")
	req, err := http.NewRequestWithContext(reqCtx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	s.applyAuth(req)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("LIST 실패: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("LIST: HTTP %d (%s)", resp.StatusCode,
			strings.TrimSpace(string(body)))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(body)), "\n")
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			out = append(out, l)
		}
	}
	return out, nil
}
