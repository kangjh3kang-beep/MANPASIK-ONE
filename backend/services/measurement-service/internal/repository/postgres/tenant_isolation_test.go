package postgres

import (
	"context"
	"strings"
	"testing"

	"github.com/manpasik/backend/shared/tenancy"
)

// buildTenantClause 헬퍼만 단독 테스트 (DB 없이).

func TestBuildTenantClause_NoTenant_DefaultsToNullFilter(t *testing.T) {
	clause, args := buildTenantClause(context.Background(), "md")
	if !strings.Contains(clause, "tenant_id IS NULL") {
		t.Errorf("ctx 미설정 시 NULL 필터 누락: %q", clause)
	}
	if len(args) != 0 {
		t.Errorf("args = %v", args)
	}
}

func TestBuildTenantClause_WithTenant_AddsEqualsFilter(t *testing.T) {
	ctx := tenancy.WithTenant(context.Background(), "hospA")
	clause, args := buildTenantClause(ctx, "md")
	if !strings.Contains(clause, "md.tenant_id = $2") {
		t.Errorf("clause = %q", clause)
	}
	if len(args) != 1 || args[0] != "hospA" {
		t.Errorf("args = %v", args)
	}
}

func TestBuildTenantClause_NoTableAlias(t *testing.T) {
	clause, _ := buildTenantClause(context.Background(), "")
	if !strings.Contains(clause, "tenant_id IS NULL") {
		t.Errorf("alias 미설정 시 컬럼명 fallback: %q", clause)
	}
	if strings.Contains(clause, ".tenant_id") {
		t.Errorf("alias 없는데 . 포함: %q", clause)
	}
}

func TestBuildTenantClause_EmptyTenant_FallbackToNull(t *testing.T) {
	ctx := tenancy.WithTenant(context.Background(), "")
	clause, args := buildTenantClause(ctx, "md")
	if !strings.Contains(clause, "IS NULL") {
		t.Errorf("빈 tenant → NULL 필터 미적용: %q", clause)
	}
	if len(args) != 0 {
		t.Errorf("빈 tenant 에서 args = %v", args)
	}
}

func TestBuildTenantClause_DifferentTenantsProduceDifferentArgs(t *testing.T) {
	ctxA := tenancy.WithTenant(context.Background(), "hospA")
	ctxB := tenancy.WithTenant(context.Background(), "hospB")

	_, argsA := buildTenantClause(ctxA, "md")
	_, argsB := buildTenantClause(ctxB, "md")

	if argsA[0] == argsB[0] {
		t.Error("tenant 가 다른데 args 동일")
	}
}

// 격리 위반 시뮬레이션:
//   - A 조직 ctx 로 GetHistory 호출 → 쿼리에 tenant_id = 'A' 포함
//   - B 조직 데이터는 자연스럽게 0행 반환 (DB 레벨 격리)
//   - 이는 SELECT 결과가 tenant 별로 분리됨을 보장
func TestTenantIsolation_QueryContainsTenantFilter(t *testing.T) {
	cases := []struct {
		name   string
		ctx    context.Context
		expect string
	}{
		{"hospA", tenancy.WithTenant(context.Background(), "hospA"), "= $2"},
		{"hospB", tenancy.WithTenant(context.Background(), "hospB"), "= $2"},
		{"none", context.Background(), "IS NULL"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			clause, _ := buildTenantClause(tc.ctx, "md")
			if !strings.Contains(clause, tc.expect) {
				t.Errorf("clause = %q, want to contain %q", clause, tc.expect)
			}
		})
	}
}
