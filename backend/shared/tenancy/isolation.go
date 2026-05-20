package tenancy

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// 데이터 격리 위반 에러 (Phase AF-1).
//
// Repository 계층에서 tenant 격리를 강제할 때 사용:
//   - 컨텍스트에 tenant 가 없을 때
//   - 리소스의 tenant_id 가 컨텍스트와 다를 때
//   - 멀티 tenant 쿼리가 격리를 우회할 때
var (
	ErrTenantContextMissing = errors.New("tenant_id 가 ctx 에 없음")
	ErrTenantMismatch       = errors.New("리소스 tenant 와 ctx tenant 불일치")
	ErrCrossTenantQuery     = errors.New("교차 tenant 쿼리 차단")
)

// RequireTenant 는 ctx 에서 tenant_id 를 강제 추출.
//
// 누락 시 ErrTenantContextMissing 반환. Repository 의 모든 SELECT/UPDATE/DELETE
// 쿼리에서 호출 → tenant 미설정 시 즉시 에러로 데이터 누출 차단.
//
// 사용 예:
//
//	tid, err := tenancy.RequireTenant(ctx)
//	if err != nil { return nil, err }
//	rows := db.Query("SELECT * FROM measurements WHERE user_id = $1 AND tenant_id = $2",
//	                  userID, tid)
func RequireTenant(ctx context.Context) (TenantID, error) {
	tid, ok := TenantFromContext(ctx)
	if !ok || tid.IsZero() {
		return "", ErrTenantContextMissing
	}
	return tid, nil
}

// AssertOwnership 은 리소스의 tenant_id 가 ctx 의 tenant 와 일치하는지 검증.
//
// 일반적인 사용:
//   - SELECT 후 결과 행마다 검증 (예방)
//   - INSERT/UPDATE 직전 (방어)
//
// 리소스의 tenant_id 가 빈 값 (legacy 데이터) 이면 ctx tenant 를 자동 부여.
func AssertOwnership(ctx context.Context, resourceTenantID string) error {
	tid, err := RequireTenant(ctx)
	if err != nil {
		return err
	}
	if resourceTenantID == "" {
		return nil // legacy: 신규로 채워줌
	}
	if string(tid) != resourceTenantID {
		return fmt.Errorf("%w (resource=%s, ctx=%s)",
			ErrTenantMismatch, resourceTenantID, string(tid))
	}
	return nil
}

// SQLFilter 는 SQL WHERE 절에 tenant_id 필터를 자동 추가.
//
// 사용 예:
//
//	tid, _ := tenancy.RequireTenant(ctx)
//	clause, args := tenancy.SQLFilter(tid, "tenant_id", existingArgs)
//	q := "SELECT * FROM measurements WHERE user_id = $1 " + clause
//	rows := db.Query(q, args...)
//
// 이 헬퍼는 SQL injection 안전 — placeholder 만 추가하고 값은 args 로 전달.
func SQLFilter(tid TenantID, columnName string, existingArgs []interface{}) (string, []interface{}) {
	if tid.IsZero() {
		// 호출자가 RequireTenant 검증 후 호출하는 게 일반적이지만 방어.
		return "", existingArgs
	}
	if columnName == "" {
		columnName = "tenant_id"
	}
	// SQL injection 방어: 컬럼명에 영문/숫자/언더스코어만 허용
	for _, c := range columnName {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') || c == '_') {
			return "", existingArgs
		}
	}

	placeholder := fmt.Sprintf("$%d", len(existingArgs)+1)
	clause := fmt.Sprintf("AND %s = %s", columnName, placeholder)
	return clause, append(existingArgs, string(tid))
}

// FilterByTenant 는 리소스 슬라이스 중 ctx 의 tenant 와 일치하는 것만 반환.
//
// MemoryStore 의 List 결과 후처리에 사용. ctx 에 tenant 가 없으면 빈 슬라이스
// 반환 (안전 default — 누출보다 거부 우선).
//
// extractor 함수가 각 리소스의 tenant_id 추출.
func FilterByTenant[T any](ctx context.Context, resources []T,
	extractor func(T) string) []T {
	tid, err := RequireTenant(ctx)
	if err != nil {
		return nil // 안전 default
	}
	out := make([]T, 0, len(resources))
	for _, r := range resources {
		if extractor(r) == string(tid) {
			out = append(out, r)
		}
	}
	return out
}

// SanitizeTenantQuery 는 사용자 입력 SQL 일부 (예: ORDER BY 컬럼) 에서
// tenant 격리를 우회하려는 시도를 차단.
//
// 예: 사용자가 "id; DROP TABLE measurements" 같은 SQL injection 시도 시
// 영문/숫자/언더스코어 외 문자 발견 시 ErrCrossTenantQuery 반환.
//
// 일반적인 ORDER BY/LIMIT 등 사용자 입력 받는 곳에서만 호출.
func SanitizeTenantQuery(input string) error {
	if input == "" {
		return nil
	}
	for _, c := range input {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') || c == '_' || c == ' ' || c == ',' ||
			c == '.' || c == '\'' || c == '"') {
			return fmt.Errorf("%w: invalid char %q in %s",
				ErrCrossTenantQuery, c, abbrev(input, 30))
		}
		// 위험 키워드 차단 (SQL injection 휴리스틱)
	}
	low := strings.ToLower(input)
	for _, banned := range []string{"drop", "delete", "truncate", "update", "insert", "--", ";", "/*"} {
		if strings.Contains(low, banned) {
			return fmt.Errorf("%w: dangerous keyword %q", ErrCrossTenantQuery, banned)
		}
	}
	return nil
}

func abbrev(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// BuildTenantClause 는 PostgreSQL/pgx 쿼리에 tenant 격리 WHERE 절을 추가.
//
// 동작 (보안 우선):
//   - ctx 에 tenant 있음    → " AND {col} = ${nextIdx}", args 에 값 추가
//   - ctx 에 tenant 없음    → " AND {col} IS NULL" (legacy/personal 만)
//
// tableAlias 가 비어있지 않으면 "{alias}.tenant_id" 형태로 사용.
// nextIdx 는 호출자가 이미 사용한 placeholder 다음 번호 (예: $1=userID 면 2 전달).
//
// 사용 예:
//
//	args := []interface{}{userID}
//	clause, tArgs := tenancy.BuildTenantClause(ctx, "md", len(args)+1)
//	q := "SELECT * FROM measurement_data md WHERE md.user_id = $1" + clause
//	args = append(args, tArgs...)
//	rows, _ := db.Query(q, args...)
func BuildTenantClause(ctx context.Context, tableAlias string, nextIdx int) (string, []interface{}) {
	col := "tenant_id"
	if tableAlias != "" {
		col = tableAlias + ".tenant_id"
	}
	tid, ok := TenantFromContext(ctx)
	if !ok || tid.IsZero() {
		return " AND " + col + " IS NULL", nil
	}
	return " AND " + col + " = $" + intToStr(nextIdx), []interface{}{string(tid)}
}

// intToStr 는 nextIdx 를 placeholder 형식으로 변환 (1자리/2자리 빠른 변환).
func intToStr(n int) string {
	if n < 10 {
		return string(rune('0' + n))
	}
	if n < 100 {
		return string(rune('0'+n/10)) + string(rune('0'+n%10))
	}
	return fmt.Sprintf("%d", n)
}
