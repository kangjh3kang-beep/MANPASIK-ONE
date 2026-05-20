package tenancy_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/manpasik/backend/shared/tenancy"
)

func TestRequireTenant_Missing(t *testing.T) {
	if _, err := tenancy.RequireTenant(context.Background()); err != tenancy.ErrTenantContextMissing {
		t.Errorf("err = %v, want ErrTenantContextMissing", err)
	}
}

func TestRequireTenant_Present(t *testing.T) {
	ctx := tenancy.WithTenant(context.Background(), "hospA")
	tid, err := tenancy.RequireTenant(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if tid != "hospA" {
		t.Errorf("tid = %q", tid)
	}
}

func TestRequireTenant_Empty(t *testing.T) {
	ctx := tenancy.WithTenant(context.Background(), "")
	if _, err := tenancy.RequireTenant(ctx); err != tenancy.ErrTenantContextMissing {
		t.Errorf("빈 tenant 통과: %v", err)
	}
}

func TestAssertOwnership_Match(t *testing.T) {
	ctx := tenancy.WithTenant(context.Background(), "hospA")
	if err := tenancy.AssertOwnership(ctx, "hospA"); err != nil {
		t.Errorf("일치인데 에러: %v", err)
	}
}

func TestAssertOwnership_Mismatch(t *testing.T) {
	ctx := tenancy.WithTenant(context.Background(), "hospA")
	err := tenancy.AssertOwnership(ctx, "hospB")
	if !errors.Is(err, tenancy.ErrTenantMismatch) {
		t.Errorf("err = %v", err)
	}
}

func TestAssertOwnership_LegacyEmpty(t *testing.T) {
	// resource_tenant_id 가 빈 값이면 legacy → 통과
	ctx := tenancy.WithTenant(context.Background(), "hospA")
	if err := tenancy.AssertOwnership(ctx, ""); err != nil {
		t.Errorf("legacy 거부: %v", err)
	}
}

func TestAssertOwnership_NoCtx(t *testing.T) {
	if err := tenancy.AssertOwnership(context.Background(), "hospA"); err == nil {
		t.Error("ctx 없는데 통과")
	}
}

func TestSQLFilter_Basic(t *testing.T) {
	clause, args := tenancy.SQLFilter("hospA", "tenant_id", []interface{}{"u-1"})
	wantClause := "AND tenant_id = $2"
	if clause != wantClause {
		t.Errorf("clause = %q, want %q", clause, wantClause)
	}
	if len(args) != 2 || args[1] != "hospA" {
		t.Errorf("args = %v", args)
	}
}

func TestSQLFilter_DefaultColumn(t *testing.T) {
	clause, _ := tenancy.SQLFilter("hospA", "", []interface{}{})
	if !strings.Contains(clause, "tenant_id") {
		t.Errorf("기본 컬럼명 미사용: %q", clause)
	}
}

func TestSQLFilter_InjectionDefense(t *testing.T) {
	// 컬럼명에 SQL injection 시도 → 빈 clause 반환
	clause, args := tenancy.SQLFilter("hospA", "tenant_id; DROP TABLE x", []interface{}{})
	if clause != "" {
		t.Errorf("injection 통과: %q", clause)
	}
	if len(args) != 0 {
		t.Errorf("args 변경됨: %v", args)
	}
}

func TestSQLFilter_EmptyTenant(t *testing.T) {
	clause, _ := tenancy.SQLFilter("", "tenant_id", []interface{}{})
	if clause != "" {
		t.Error("빈 tenant 에서 clause 생성")
	}
}

func TestFilterByTenant_KeepsMatching(t *testing.T) {
	type res struct {
		ID       string
		TenantID string
	}
	resources := []res{
		{"a", "hospA"},
		{"b", "hospB"},
		{"c", "hospA"},
	}
	ctx := tenancy.WithTenant(context.Background(), "hospA")
	got := tenancy.FilterByTenant(ctx, resources, func(r res) string { return r.TenantID })
	if len(got) != 2 {
		t.Errorf("len = %d, want 2", len(got))
	}
}

func TestFilterByTenant_NoCtx_EmptyResult(t *testing.T) {
	type res struct{ TenantID string }
	resources := []res{{"hospA"}, {"hospB"}}
	got := tenancy.FilterByTenant(context.Background(), resources, func(r res) string { return r.TenantID })
	if got != nil {
		t.Errorf("ctx 없는데 결과 = %v (want nil)", got)
	}
}

func TestSanitizeTenantQuery_Empty(t *testing.T) {
	if err := tenancy.SanitizeTenantQuery(""); err != nil {
		t.Errorf("빈 입력 거부: %v", err)
	}
}

func TestSanitizeTenantQuery_Valid(t *testing.T) {
	for _, q := range []string{"created_at", "id ASC", "user_id, created_at"} {
		if err := tenancy.SanitizeTenantQuery(q); err != nil {
			t.Errorf("%q 거부: %v", q, err)
		}
	}
}

func TestSanitizeTenantQuery_BlocksInjection(t *testing.T) {
	dangerous := []string{
		"id; DROP TABLE measurements",
		"created_at -- comment",
		"id; DELETE FROM x",
		"id /* something */",
		"id) UNION SELECT *",
	}
	for _, q := range dangerous {
		if err := tenancy.SanitizeTenantQuery(q); err == nil {
			t.Errorf("위험 입력 통과: %q", q)
		}
	}
}

func TestSanitizeTenantQuery_BlocksSpecialChars(t *testing.T) {
	for _, q := range []string{"id\nDROP", "id$x"} {
		if err := tenancy.SanitizeTenantQuery(q); err == nil {
			t.Errorf("특수 문자 통과: %q", q)
		}
	}
}

func TestBuildTenantClause_NoTenant(t *testing.T) {
	clause, args := tenancy.BuildTenantClause(context.Background(), "md", 2)
	if !strings.Contains(clause, "md.tenant_id IS NULL") {
		t.Errorf("clause = %q", clause)
	}
	if len(args) != 0 {
		t.Errorf("args = %v", args)
	}
}

func TestBuildTenantClause_WithTenant(t *testing.T) {
	ctx := tenancy.WithTenant(context.Background(), "hospA")
	clause, args := tenancy.BuildTenantClause(ctx, "hr", 3)
	if !strings.Contains(clause, "hr.tenant_id = $3") {
		t.Errorf("clause = %q", clause)
	}
	if len(args) != 1 || args[0] != "hospA" {
		t.Errorf("args = %v", args)
	}
}

func TestBuildTenantClause_NoAlias(t *testing.T) {
	ctx := tenancy.WithTenant(context.Background(), "hospA")
	clause, _ := tenancy.BuildTenantClause(ctx, "", 2)
	if strings.Contains(clause, ".tenant_id") {
		t.Errorf("alias 없는데 . 포함: %q", clause)
	}
	if !strings.Contains(clause, "tenant_id = $2") {
		t.Errorf("clause = %q", clause)
	}
}

func TestBuildTenantClause_LargeIndex(t *testing.T) {
	ctx := tenancy.WithTenant(context.Background(), "hospA")
	clause, _ := tenancy.BuildTenantClause(ctx, "x", 15)
	if !strings.Contains(clause, "$15") {
		t.Errorf("2자리 인덱스 처리 실패: %q", clause)
	}
}
