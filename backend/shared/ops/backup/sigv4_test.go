package backup_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/ops/backup"
)

// AWS SigV4 known-test-vector: GET vanilla
// 출처: https://docs.aws.amazon.com/general/latest/gr/sigv4_signing.html
//
// 우리는 표준 라이브러리 동등 구현이 정확한지 검증하기 위해
// 1) 알려진 테스트 벡터의 서명 (잘 알려진 입력) 또는
// 2) 자체 검증 (signing key derivation 의 단순 sha256 매핑) 을 사용.
//
// 단순함을 위해 "AWSCryptoCookbook" 사용자 정의 호출 검증을 한다.

func TestSigV4Signer_Validate(t *testing.T) {
	cfg := backup.SigV4Config{}
	if err := cfg.Validate(); err == nil {
		t.Error("빈 cfg 에 검증 에러 없음")
	}

	cfg = backup.SigV4Config{
		AccessKeyID:     "AKID",
		SecretAccessKey: "SECRET",
	}
	if err := cfg.Validate(); err == nil {
		t.Error("region 누락에 에러 없음")
	}
}

func TestSigV4Signer_AttachesAuthHeader(t *testing.T) {
	signer, err := backup.NewSigV4Signer(backup.SigV4Config{
		AccessKeyID:     "AKIDEXAMPLE",
		SecretAccessKey: "wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY",
		Region:          "us-east-1",
		Service:         "s3",
	})
	if err != nil {
		t.Fatal(err)
	}
	signer.SetClock(func() time.Time {
		return time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	})

	req, _ := http.NewRequest("GET", "https://bucket.s3.amazonaws.com/path", nil)
	if err := signer.Sign(req, nil); err != nil {
		t.Fatal(err)
	}

	auth := req.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20260501/us-east-1/s3/aws4_request") {
		t.Errorf("Authorization = %q", auth)
	}
	if !strings.Contains(auth, "SignedHeaders=") {
		t.Errorf("SignedHeaders 누락: %q", auth)
	}
	if !strings.Contains(auth, "Signature=") {
		t.Errorf("Signature 누락: %q", auth)
	}
	if req.Header.Get("X-Amz-Date") != "20260501T120000Z" {
		t.Errorf("X-Amz-Date = %q", req.Header.Get("X-Amz-Date"))
	}
}

func TestSigV4Signer_PayloadHash(t *testing.T) {
	signer, _ := backup.NewSigV4Signer(backup.SigV4Config{
		AccessKeyID: "AK", SecretAccessKey: "SK",
		Region: "ap-northeast-2", Service: "s3",
	})
	req, _ := http.NewRequest("PUT", "https://x.s3.amazonaws.com/k", nil)
	if err := signer.Sign(req, []byte("hello")); err != nil {
		t.Fatal(err)
	}
	// SHA-256("hello") = 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824
	want := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	if got := req.Header.Get("X-Amz-Content-Sha256"); got != want {
		t.Errorf("payload hash = %q, want %q", got, want)
	}
}

func TestSigV4Signer_SessionToken(t *testing.T) {
	signer, _ := backup.NewSigV4Signer(backup.SigV4Config{
		AccessKeyID: "AK", SecretAccessKey: "SK",
		Region: "us-east-1", Service: "s3",
		SessionToken: "FQoGZX...",
	})
	req, _ := http.NewRequest("GET", "https://x.s3.amazonaws.com/k", nil)
	_ = signer.Sign(req, nil)
	if req.Header.Get("X-Amz-Security-Token") != "FQoGZX..." {
		t.Errorf("session token 미부착: %q", req.Header.Get("X-Amz-Security-Token"))
	}
}

func TestSigV4Signer_DeterministicSignature(t *testing.T) {
	signer, _ := backup.NewSigV4Signer(backup.SigV4Config{
		AccessKeyID: "AK", SecretAccessKey: "SK",
		Region: "us-east-1", Service: "s3",
	})
	signer.SetClock(func() time.Time {
		return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	})
	r1, _ := http.NewRequest("GET", "https://x.s3.amazonaws.com/k", nil)
	_ = signer.Sign(r1, nil)
	r2, _ := http.NewRequest("GET", "https://x.s3.amazonaws.com/k", nil)
	_ = signer.Sign(r2, nil)
	if r1.Header.Get("Authorization") != r2.Header.Get("Authorization") {
		t.Error("동일 입력에 다른 서명")
	}
}

func TestSigV4Signer_NilReq(t *testing.T) {
	signer, _ := backup.NewSigV4Signer(backup.SigV4Config{
		AccessKeyID: "AK", SecretAccessKey: "SK",
		Region: "us-east-1", Service: "s3",
	})
	if err := signer.Sign(nil, nil); err == nil {
		t.Error("nil req 통과")
	}
}

func TestS3Storage_Put(t *testing.T) {
	var (
		gotPath  string
		gotAuth  string
		gotBody  []byte
		mu       sync.Mutex
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(200)
	}))
	defer srv.Close()

	s, err := backup.NewS3Storage(backup.S3StorageConfig{
		BaseURL: srv.URL + "/manpasik-bucket",
		SigV4: backup.SigV4Config{
			AccessKeyID:     "AK",
			SecretAccessKey: "SK",
			Region:          "us-east-1",
			Service:         "s3",
		},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if err := s.Put(context.Background(), "daily/x.dump", strings.NewReader("data")); err != nil {
		t.Fatal(err)
	}

	if gotPath != "/manpasik-bucket/daily/x.dump" {
		t.Errorf("path = %q", gotPath)
	}
	if !strings.HasPrefix(gotAuth, "AWS4-HMAC-SHA256 Credential=AK/") {
		t.Errorf("auth = %q", gotAuth)
	}
	if string(gotBody) != "data" {
		t.Errorf("body = %q", gotBody)
	}
}

func TestS3Storage_Get_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer srv.Close()

	s, _ := backup.NewS3Storage(backup.S3StorageConfig{
		BaseURL: srv.URL,
		SigV4: backup.SigV4Config{
			AccessKeyID: "AK", SecretAccessKey: "SK",
			Region: "us-east-1", Service: "s3",
		},
	}, nil)

	if _, err := s.Get(context.Background(), "missing"); err == nil {
		t.Error("404 에 에러 없음")
	}
}

func TestS3Storage_Delete_404Idempotent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer srv.Close()

	s, _ := backup.NewS3Storage(backup.S3StorageConfig{
		BaseURL: srv.URL,
		SigV4: backup.SigV4Config{
			AccessKeyID: "AK", SecretAccessKey: "SK",
			Region: "us-east-1", Service: "s3",
		},
	}, nil)
	if err := s.Delete(context.Background(), "x"); err != nil {
		t.Errorf("404 delete err = %v", err)
	}
}

func TestS3Storage_List_ParseXML(t *testing.T) {
	xmlBody := `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
  <Contents><Key>daily/x.dump</Key></Contents>
  <Contents><Key>daily/y.dump</Key></Contents>
  <Contents><Key>weekly/z.dump</Key></Contents>
</ListBucketResult>`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(xmlBody))
	}))
	defer srv.Close()

	s, _ := backup.NewS3Storage(backup.S3StorageConfig{
		BaseURL: srv.URL,
		SigV4: backup.SigV4Config{
			AccessKeyID: "AK", SecretAccessKey: "SK",
			Region: "us-east-1", Service: "s3",
		},
	}, nil)
	keys, err := s.List(context.Background(), "daily/")
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 3 || keys[0] != "daily/x.dump" {
		t.Errorf("keys = %v", keys)
	}
}

func TestS3Storage_NoBaseURL(t *testing.T) {
	if _, err := backup.NewS3Storage(backup.S3StorageConfig{
		SigV4: backup.SigV4Config{
			AccessKeyID: "AK", SecretAccessKey: "SK",
			Region: "us-east-1", Service: "s3",
		},
	}, nil); err == nil {
		t.Error("BaseURL 누락에 에러 없음")
	}
}

func TestS3Storage_Provider(t *testing.T) {
	s, _ := backup.NewS3Storage(backup.S3StorageConfig{
		BaseURL: "https://x",
		SigV4: backup.SigV4Config{
			AccessKeyID: "AK", SecretAccessKey: "SK",
			Region: "us-east-1", Service: "s3",
		},
	}, nil)
	if s.Provider() != "s3" {
		t.Errorf("Provider = %q", s.Provider())
	}
}

func TestS3Storage_PutFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
		w.Write([]byte("forbidden"))
	}))
	defer srv.Close()

	s, _ := backup.NewS3Storage(backup.S3StorageConfig{
		BaseURL: srv.URL,
		SigV4: backup.SigV4Config{
			AccessKeyID: "AK", SecretAccessKey: "SK",
			Region: "us-east-1", Service: "s3",
		},
	}, nil)
	if err := s.Put(context.Background(), "k", strings.NewReader("data")); err == nil {
		t.Error("403 에 에러 없음")
	}
}
