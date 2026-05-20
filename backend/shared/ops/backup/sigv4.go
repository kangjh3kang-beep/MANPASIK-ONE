package backup

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// AWS SigV4 — Authorization 헤더 기반 서명.
//
// 참고: https://docs.aws.amazon.com/general/latest/gr/sigv4-create-canonical-request.html
//
// 외부 SDK 의존 없이 표준 라이브러리(crypto/sha256/hmac)만 사용.

// SigV4Config 는 SigV4 서명 설정.
type SigV4Config struct {
	// AccessKeyID 는 IAM access key.
	AccessKeyID string
	// SecretAccessKey 는 IAM secret.
	SecretAccessKey string
	// SessionToken 는 STS 세션 토큰 (선택).
	SessionToken string
	// Region 은 AWS 리전 (예: "us-east-1", "ap-northeast-2").
	Region string
	// Service 는 AWS 서비스 (S3 의 경우 "s3").
	Service string
}

// Validate 는 필수 필드 검증.
func (c *SigV4Config) Validate() error {
	if c.AccessKeyID == "" || c.SecretAccessKey == "" {
		return errors.New("AccessKeyID/SecretAccessKey 필수")
	}
	if c.Region == "" {
		return errors.New("Region 필수")
	}
	if c.Service == "" {
		return errors.New("Service 필수 (예: s3)")
	}
	return nil
}

// SigV4Signer 는 HTTP 요청에 AWS SigV4 Authorization 헤더 부착.
type SigV4Signer struct {
	cfg SigV4Config
	// nowFn 은 테스트용 — 기본 time.Now.
	nowFn func() time.Time
}

// NewSigV4Signer 생성.
func NewSigV4Signer(cfg SigV4Config) (*SigV4Signer, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &SigV4Signer{cfg: cfg, nowFn: time.Now}, nil
}

// SetClock 은 테스트용 clock 주입.
func (s *SigV4Signer) SetClock(now func() time.Time) {
	if now != nil {
		s.nowFn = now
	}
}

// Sign 은 req 에 SigV4 헤더(Authorization, X-Amz-Date, X-Amz-Content-Sha256) 부착.
//
// payload 는 요청 본문 — req.Body 가 io.Reader 일 때 별도로 받아야 함 (Body 는
// 단방향 stream 이라 hash 후 다시 사용 불가). payload=nil 이면 빈 페이로드.
func (s *SigV4Signer) Sign(req *http.Request, payload []byte) error {
	if req == nil {
		return errors.New("req nil")
	}
	if req.URL == nil {
		return errors.New("req.URL nil")
	}

	now := s.nowFn().UTC()
	amzDate := now.Format("20060102T150405Z")
	dateStamp := now.Format("20060102")

	if payload == nil {
		payload = []byte{}
	}
	payloadHash := hex.EncodeToString(hashSHA256(payload))

	req.Header.Set("X-Amz-Date", amzDate)
	req.Header.Set("X-Amz-Content-Sha256", payloadHash)
	if s.cfg.SessionToken != "" {
		req.Header.Set("X-Amz-Security-Token", s.cfg.SessionToken)
	}
	if req.Header.Get("Host") == "" && req.URL.Host != "" {
		req.Host = req.URL.Host
	}

	canonicalReq, signedHeaders := buildCanonicalRequest(req, payloadHash)
	credentialScope := fmt.Sprintf("%s/%s/%s/aws4_request", dateStamp, s.cfg.Region, s.cfg.Service)
	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256",
		amzDate,
		credentialScope,
		hex.EncodeToString(hashSHA256([]byte(canonicalReq))),
	}, "\n")

	signingKey := deriveSigningKey(s.cfg.SecretAccessKey, dateStamp, s.cfg.Region, s.cfg.Service)
	signature := hex.EncodeToString(hmacSHA256(signingKey, []byte(stringToSign)))

	authHeader := fmt.Sprintf(
		"AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		s.cfg.AccessKeyID, credentialScope, signedHeaders, signature,
	)
	req.Header.Set("Authorization", authHeader)
	return nil
}

// buildCanonicalRequest 는 SigV4 canonical request 와 signed headers 문자열을 반환.
func buildCanonicalRequest(req *http.Request, payloadHash string) (canonical string, signedHeaders string) {
	method := req.Method
	canonicalURI := canonicalURIFromURL(req.URL)
	canonicalQuery := canonicalQueryString(req.URL)
	canonicalHeaders, signedList := buildCanonicalHeaders(req)

	canonical = strings.Join([]string{
		method,
		canonicalURI,
		canonicalQuery,
		canonicalHeaders,
		signedList,
		payloadHash,
	}, "\n")
	return canonical, signedList
}

// canonicalURIFromURL 은 URI 의 path 부분 (이중 인코딩 안 함, S3 의 경우).
func canonicalURIFromURL(u *url.URL) string {
	if u.Path == "" {
		return "/"
	}
	// S3: 경로 컴포넌트 별로 RFC 3986 unreserved 만 인코딩 (/ 유지).
	parts := strings.Split(u.Path, "/")
	for i, p := range parts {
		parts[i] = uriEscapeS3(p)
	}
	return strings.Join(parts, "/")
}

// uriEscapeS3 는 S3 SigV4 의 경로 인코딩 (slash 미인코딩, unreserved 만 raw).
func uriEscapeS3(s string) string {
	var b strings.Builder
	for _, r := range s {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') ||
			(r >= '0' && r <= '9') || r == '-' || r == '_' || r == '.' || r == '~' {
			b.WriteRune(r)
		} else {
			b.WriteString(fmt.Sprintf("%%%02X", r))
		}
	}
	return b.String()
}

// canonicalQueryString 은 쿼리 파라미터를 (이름,값) 알파벳 정렬 + 인코딩.
func canonicalQueryString(u *url.URL) string {
	if u.RawQuery == "" {
		return ""
	}
	q := u.Query()
	keys := make([]string, 0, len(q))
	for k := range q {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		vs := q[k]
		sort.Strings(vs)
		for _, v := range vs {
			parts = append(parts, url.QueryEscape(k)+"="+url.QueryEscape(v))
		}
	}
	return strings.Join(parts, "&")
}

// buildCanonicalHeaders 는 정렬된 (lowercase header : trimmed value) + signed list.
func buildCanonicalHeaders(req *http.Request) (canonicalHeaders, signedHeaders string) {
	headers := make(map[string][]string)
	for k, v := range req.Header {
		lk := strings.ToLower(k)
		headers[lk] = v
	}
	// host 는 명시 추가 (req.Header 에 없을 수 있음)
	if _, ok := headers["host"]; !ok {
		host := req.URL.Host
		if req.Host != "" {
			host = req.Host
		}
		headers["host"] = []string{host}
	}

	keys := make([]string, 0, len(headers))
	for k := range headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var lines []string
	for _, k := range keys {
		var vals []string
		for _, v := range headers[k] {
			vals = append(vals, strings.TrimSpace(collapseWhitespace(v)))
		}
		lines = append(lines, k+":"+strings.Join(vals, ","))
	}
	canonicalHeaders = strings.Join(lines, "\n") + "\n"
	signedHeaders = strings.Join(keys, ";")
	return canonicalHeaders, signedHeaders
}

// collapseWhitespace 는 연속 공백을 하나로 압축 (SigV4 요건).
func collapseWhitespace(s string) string {
	var b strings.Builder
	prevSpace := false
	for _, r := range s {
		if r == ' ' || r == '\t' {
			if !prevSpace {
				b.WriteByte(' ')
			}
			prevSpace = true
		} else {
			b.WriteRune(r)
			prevSpace = false
		}
	}
	return b.String()
}

// deriveSigningKey 는 4단계 HMAC-SHA256 으로 서명 키 도출.
func deriveSigningKey(secret, date, region, service string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+secret), []byte(date))
	kRegion := hmacSHA256(kDate, []byte(region))
	kService := hmacSHA256(kRegion, []byte(service))
	kSigning := hmacSHA256(kService, []byte("aws4_request"))
	return kSigning
}

func hmacSHA256(key, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}

func hashSHA256(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}

// ============================================================================
// SigV4-aware ObjectStorage
// ============================================================================

// S3StorageConfig 는 SigV4 인증 + S3 호환 저장소 설정.
//
// 동작 모드:
//   - virtual-hosted-style: BaseURL = "https://bucket-name.s3.region.amazonaws.com"
//   - path-style: BaseURL = "https://s3.region.amazonaws.com/bucket-name"
//
// MinIO 등 호환 서비스는 path-style + 자체 endpoint 사용.
type S3StorageConfig struct {
	BaseURL        string
	SigV4          SigV4Config
	RequestTimeout time.Duration
}

// S3Storage 는 SigV4 서명 후 HTTP 호출하는 ObjectStorage 구현.
type S3Storage struct {
	cfg     S3StorageConfig
	signer  *SigV4Signer
	doer    HTTPDoer
}

// NewS3Storage 생성. doer=nil 이면 http.Client{Timeout:cfg.RequestTimeout}.
func NewS3Storage(cfg S3StorageConfig, doer HTTPDoer) (*S3Storage, error) {
	if cfg.BaseURL == "" {
		return nil, errors.New("BaseURL 필수")
	}
	if cfg.RequestTimeout <= 0 {
		cfg.RequestTimeout = 30 * time.Second
	}
	signer, err := NewSigV4Signer(cfg.SigV4)
	if err != nil {
		return nil, err
	}
	if doer == nil {
		doer = &http.Client{Timeout: cfg.RequestTimeout}
	}
	cfg.BaseURL = strings.TrimRight(cfg.BaseURL, "/")
	return &S3Storage{cfg: cfg, signer: signer, doer: doer}, nil
}

// SetClock 은 테스트용 — signer 의 clock 주입.
func (s *S3Storage) SetClock(now func() time.Time) {
	s.signer.SetClock(now)
}

// Provider 이름.
func (s *S3Storage) Provider() string { return "s3" }

func (s *S3Storage) buildURL(key string) string {
	return s.cfg.BaseURL + "/" + strings.TrimLeft(key, "/")
}

// Put 은 PUT key, payload 를 SigV4 서명하여 송신.
func (s *S3Storage) Put(ctx context.Context, key string, data io.Reader) error {
	if key == "" {
		return errors.New("key 필수")
	}
	if data == nil {
		data = strings.NewReader("")
	}
	body, err := io.ReadAll(data)
	if err != nil {
		return fmt.Errorf("body read: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", s.buildURL(key), strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	req.ContentLength = int64(len(body))
	req.Header.Set("Content-Type", "application/octet-stream")
	if err := s.signer.Sign(req, body); err != nil {
		return err
	}
	resp, err := s.doer.Do(req)
	if err != nil {
		return fmt.Errorf("PUT 실패: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		buf, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("PUT %s: HTTP %d (%s)", key, resp.StatusCode, strings.TrimSpace(string(buf)))
	}
	return nil
}

// Get 은 GET key (SigV4 서명, payload 빈 페이로드 해시).
func (s *S3Storage) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.buildURL(key), nil)
	if err != nil {
		return nil, err
	}
	if err := s.signer.Sign(req, nil); err != nil {
		return nil, err
	}
	resp, err := s.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET 실패: %w", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		resp.Body.Close()
		return nil, fmt.Errorf("키 없음: %s", key)
	}
	if resp.StatusCode >= 300 {
		buf, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		resp.Body.Close()
		return nil, fmt.Errorf("GET %s: HTTP %d (%s)", key, resp.StatusCode, strings.TrimSpace(string(buf)))
	}
	return resp.Body, nil
}

// Delete 는 DELETE key (404 멱등).
func (s *S3Storage) Delete(ctx context.Context, key string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", s.buildURL(key), nil)
	if err != nil {
		return err
	}
	if err := s.signer.Sign(req, nil); err != nil {
		return err
	}
	resp, err := s.doer.Do(req)
	if err != nil {
		return fmt.Errorf("DELETE 실패: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}
	if resp.StatusCode >= 300 {
		buf, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("DELETE %s: HTTP %d (%s)", key, resp.StatusCode, strings.TrimSpace(string(buf)))
	}
	return nil
}

// List 는 S3 ListObjectsV2 (?list-type=2&prefix=...) — XML 응답 파싱은 단순화.
//
// XML 파싱은 별도 어댑터에서 강화. 기본 구현은 응답 본문을 그대로 [1]개 키로 반환.
func (s *S3Storage) List(ctx context.Context, prefix string) ([]string, error) {
	u := s.cfg.BaseURL + "/?list-type=2&prefix=" + url.QueryEscape(prefix)
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}
	if err := s.signer.Sign(req, nil); err != nil {
		return nil, err
	}
	resp, err := s.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("LIST 실패: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		buf, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("LIST: HTTP %d (%s)", resp.StatusCode, strings.TrimSpace(string(buf)))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return parseListBucketResultKeys(body), nil
}

// parseListBucketResultKeys 는 ListBucketResult XML 에서 <Key>...</Key> 만 추출.
//
// 정식 XML 파싱 대신 정규식적 단순 추출 (의존성 최소화).
func parseListBucketResultKeys(body []byte) []string {
	s := string(body)
	var keys []string
	for {
		i := strings.Index(s, "<Key>")
		if i < 0 {
			break
		}
		j := strings.Index(s[i+5:], "</Key>")
		if j < 0 {
			break
		}
		keys = append(keys, s[i+5:i+5+j])
		s = s[i+5+j+6:]
	}
	return keys
}

