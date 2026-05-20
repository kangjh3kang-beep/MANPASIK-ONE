package backup_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/manpasik/backend/shared/ops/backup"
)

func newRecordingServer(t *testing.T, fn http.HandlerFunc) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(fn)
	t.Cleanup(srv.Close)
	return srv
}

func TestHTTPObjectStorage_Put(t *testing.T) {
	var (
		gotPath   string
		gotMethod string
		gotAuth   string
		gotBody   []byte
		mu        sync.Mutex
	)
	srv := newRecordingServer(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		gotPath = r.URL.Path
		gotMethod = r.Method
		gotAuth = r.Header.Get("Authorization")
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(200)
	})

	s, err := backup.NewHTTPObjectStorage(backup.HTTPStorageConfig{
		BaseURL:    srv.URL + "/bucket",
		AuthHeader: "Authorization",
		AuthValue:  "Bearer abc",
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if err := s.Put(context.Background(), "daily/x.dump", strings.NewReader("payload")); err != nil {
		t.Fatal(err)
	}

	if gotMethod != "PUT" {
		t.Errorf("method = %q", gotMethod)
	}
	if gotPath != "/bucket/daily/x.dump" {
		t.Errorf("path = %q", gotPath)
	}
	if gotAuth != "Bearer abc" {
		t.Errorf("auth = %q", gotAuth)
	}
	if string(gotBody) != "payload" {
		t.Errorf("body = %q", gotBody)
	}
}

func TestHTTPObjectStorage_Get(t *testing.T) {
	srv := newRecordingServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("method = %q", r.Method)
		}
		w.Write([]byte("hello"))
	})
	s, _ := backup.NewHTTPObjectStorage(backup.HTTPStorageConfig{
		BaseURL: srv.URL + "/b",
	}, nil)

	rc, err := s.Get(context.Background(), "k.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer rc.Close()
	body, _ := io.ReadAll(rc)
	if string(body) != "hello" {
		t.Errorf("body = %q", body)
	}
}

func TestHTTPObjectStorage_Get_NotFound(t *testing.T) {
	srv := newRecordingServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	s, _ := backup.NewHTTPObjectStorage(backup.HTTPStorageConfig{BaseURL: srv.URL}, nil)
	if _, err := s.Get(context.Background(), "missing"); err == nil {
		t.Error("404 에 에러 없음")
	}
}

func TestHTTPObjectStorage_Get_Error(t *testing.T) {
	srv := newRecordingServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("internal err"))
	})
	s, _ := backup.NewHTTPObjectStorage(backup.HTTPStorageConfig{BaseURL: srv.URL}, nil)
	if _, err := s.Get(context.Background(), "x"); err == nil {
		t.Error("500 에 에러 없음")
	}
}

func TestHTTPObjectStorage_Delete(t *testing.T) {
	called := false
	srv := newRecordingServer(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		if r.Method != "DELETE" {
			t.Errorf("method = %q", r.Method)
		}
		w.WriteHeader(204)
	})
	s, _ := backup.NewHTTPObjectStorage(backup.HTTPStorageConfig{BaseURL: srv.URL}, nil)
	if err := s.Delete(context.Background(), "x"); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("서버 미호출")
	}
}

func TestHTTPObjectStorage_Delete_404Idempotent(t *testing.T) {
	srv := newRecordingServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})
	s, _ := backup.NewHTTPObjectStorage(backup.HTTPStorageConfig{BaseURL: srv.URL}, nil)
	if err := s.Delete(context.Background(), "x"); err != nil {
		t.Errorf("404 에 에러: %v", err)
	}
}

func TestHTTPObjectStorage_List(t *testing.T) {
	var gotURL string
	srv := newRecordingServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotURL = r.URL.String()
		w.Write([]byte("a/x.dump\na/y.dump\nb/z.dump\n"))
	})
	s, _ := backup.NewHTTPObjectStorage(backup.HTTPStorageConfig{
		BaseURL:      srv.URL,
		ListEndpoint: srv.URL + "/list",
	}, nil)

	keys, err := s.List(context.Background(), "a/")
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 3 {
		t.Errorf("keys = %v", keys)
	}
	if !strings.Contains(gotURL, "prefix=a/") {
		t.Errorf("URL = %q", gotURL)
	}
}

func TestHTTPObjectStorage_List_Unsupported(t *testing.T) {
	s, _ := backup.NewHTTPObjectStorage(backup.HTTPStorageConfig{
		BaseURL: "http://x",
	}, nil)
	if _, err := s.List(context.Background(), ""); err == nil {
		t.Error("ListEndpoint 미설정인데 에러 없음")
	}
}

func TestHTTPObjectStorage_NoBaseURL(t *testing.T) {
	if _, err := backup.NewHTTPObjectStorage(backup.HTTPStorageConfig{}, nil); err == nil {
		t.Error("빈 BaseURL 통과")
	}
}

func TestHTTPObjectStorage_PutFailure(t *testing.T) {
	srv := newRecordingServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
		w.Write([]byte("forbidden"))
	})
	s, _ := backup.NewHTTPObjectStorage(backup.HTTPStorageConfig{BaseURL: srv.URL}, nil)
	if err := s.Put(context.Background(), "x", strings.NewReader("y")); err == nil {
		t.Error("403 에 에러 없음")
	}
}

// 가짜 doer 로 timeout 흉내
type slowDoer struct{}

func (slowDoer) Do(req *http.Request) (*http.Response, error) {
	<-req.Context().Done()
	return nil, req.Context().Err()
}

func TestHTTPObjectStorage_Timeout(t *testing.T) {
	s, _ := backup.NewHTTPObjectStorage(backup.HTTPStorageConfig{
		BaseURL:        "http://x",
		RequestTimeout: 50_000_000, // 50ms
	}, slowDoer{})
	if err := s.Put(context.Background(), "k", strings.NewReader("d")); err == nil {
		t.Error("타임아웃에 에러 없음")
	}
}

func TestHTTPObjectStorage_Provider(t *testing.T) {
	s, _ := backup.NewHTTPObjectStorage(backup.HTTPStorageConfig{BaseURL: "http://x"}, nil)
	if s.Provider() != "http" {
		t.Errorf("Provider = %q", s.Provider())
	}
}

func TestHTTPObjectStorage_PutNilData(t *testing.T) {
	srv := newRecordingServer(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if len(body) != 0 {
			t.Errorf("body = %q (want empty)", body)
		}
		w.WriteHeader(200)
	})
	s, _ := backup.NewHTTPObjectStorage(backup.HTTPStorageConfig{BaseURL: srv.URL}, nil)
	if err := s.Put(context.Background(), "k", nil); err != nil {
		t.Fatal(err)
	}
}
