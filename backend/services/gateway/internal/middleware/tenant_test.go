package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/grpc/metadata"

	gw "github.com/manpasik/backend/services/gateway/internal/middleware"
	"github.com/manpasik/backend/shared/tenancy"
)

func TestTenantPropagation_HeadersInjected(t *testing.T) {
	captured := struct {
		tenant string
		user   string
		ok     bool
	}{}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if tid, ok := tenancy.TenantFromContext(ctx); ok {
			captured.tenant = string(tid)
		}
		if uid, ok := tenancy.UserFromContext(ctx); ok {
			captured.user = uid
		}
		captured.ok = true
		w.WriteHeader(200)
	})

	h := gw.TenantPropagation(next)
	r := httptest.NewRequest("GET", "/x", nil)
	r.Header.Set("X-Tenant-ID", "hospital-A")
	r.Header.Set("X-User-ID", "u-001")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if !captured.ok {
		t.Fatal("next 미호출")
	}
	if captured.tenant != "hospital-A" {
		t.Errorf("tenant = %q", captured.tenant)
	}
	if captured.user != "u-001" {
		t.Errorf("user = %q", captured.user)
	}
}

func TestTenantPropagation_MissingHeaders(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		// 미주입 시 ctx 에 테넌트/유저 없어야
		if _, ok := tenancy.TenantFromContext(r.Context()); ok {
			t.Error("미설정 시 tenant 가 포함됨")
		}
		w.WriteHeader(200)
	})
	h := gw.TenantPropagation(next)
	r := httptest.NewRequest("GET", "/x", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if !called {
		t.Fatal("next 미호출")
	}
}

func TestWithTenantMetadata(t *testing.T) {
	ctx := tenancy.WithUser(tenancy.WithTenant(httptest.NewRequest("GET", "/", nil).Context(),
		"t-1"), "u-1")
	ctx = gw.WithTenantMetadata(ctx)
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		t.Fatal("outgoing metadata 없음")
	}
	if vals := md.Get(tenancy.MetadataTenantKey); len(vals) == 0 || vals[0] != "t-1" {
		t.Errorf("tenant md = %v", vals)
	}
	if vals := md.Get(tenancy.MetadataUserKey); len(vals) == 0 || vals[0] != "u-1" {
		t.Errorf("user md = %v", vals)
	}
}

func TestWithTenantMetadata_NoContext(t *testing.T) {
	// ctx 에 테넌트/유저 없으면 outgoing metadata 추가 안 됨
	ctx := httptest.NewRequest("GET", "/", nil).Context()
	ctx = gw.WithTenantMetadata(ctx)
	md, ok := metadata.FromOutgoingContext(ctx)
	// 메타데이터 없거나 비어 있어야
	if ok && (len(md.Get(tenancy.MetadataTenantKey)) > 0 || len(md.Get(tenancy.MetadataUserKey)) > 0) {
		t.Errorf("빈 ctx 에서 outgoing md 추가됨: %v", md)
	}
}

func TestTenantPropagation_OnlyTenant(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := tenancy.TenantFromContext(r.Context()); !ok {
			t.Error("tenant 누락")
		}
		if _, ok := tenancy.UserFromContext(r.Context()); ok {
			t.Error("user 헤더 미설정인데 ctx 에 존재")
		}
	})
	h := gw.TenantPropagation(next)
	r := httptest.NewRequest("GET", "/x", nil)
	r.Header.Set("X-Tenant-ID", "t-only")
	h.ServeHTTP(httptest.NewRecorder(), r)
}
