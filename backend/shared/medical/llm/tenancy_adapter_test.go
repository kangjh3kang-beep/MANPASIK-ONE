package llm

import (
	"context"
	"errors"
	"testing"

	"github.com/manpasik/backend/shared/tenancy"
)

func TestTenancyAdapter_RecordsTenantInAudit(t *testing.T) {
	inner := NewNoopAdapter()
	audit := NewMemoryAuditLog()
	adapter := NewTenancyAdapter(inner, audit)

	ctx := tenancy.WithTenant(context.Background(), "hospA")
	_, err := adapter.Complete(ctx, &Request{
		UserID: "doc1",
		Messages: []*Message{
			{Role: RoleUser, Content: "안녕"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	entries := audit.EntriesByTenant("hospA")
	if len(entries) != 1 {
		t.Fatalf("entries = %d", len(entries))
	}
	if entries[0].UserID != "doc1" {
		t.Errorf("UserID = %q", entries[0].UserID)
	}
	if !entries[0].Success {
		t.Error("성공인데 Success=false")
	}
}

func TestTenancyAdapter_NoTenantUsesPersonal(t *testing.T) {
	inner := NewNoopAdapter()
	audit := NewMemoryAuditLog()
	adapter := NewTenancyAdapter(inner, audit)

	_, _ = adapter.Complete(context.Background(), &Request{
		UserID: "u1",
		Messages: []*Message{
			{Role: RoleUser, Content: "test"},
		},
	})

	personal := audit.EntriesByTenant("personal")
	if len(personal) != 1 {
		t.Errorf("personal entries = %d", len(personal))
	}
}

func TestTenancyAdapter_RecordsErrorInAudit(t *testing.T) {
	inner := &failingAdapter{err: errors.New("rate limited")}
	audit := NewMemoryAuditLog()
	adapter := NewTenancyAdapter(inner, audit)

	ctx := tenancy.WithTenant(context.Background(), "hospB")
	_, err := adapter.Complete(ctx, &Request{
		UserID: "u",
		Messages: []*Message{{Role: RoleUser, Content: "안녕하세요 이것은 토큰 카운트를 위한 충분히 긴 의료 상담 메시지입니다 도움이 필요합니다"}},
	})
	if err == nil {
		t.Fatal("에러가 전파되지 않음")
	}

	entries := audit.EntriesByTenant("hospB")
	if len(entries) != 1 {
		t.Fatalf("entries = %d", len(entries))
	}
	if entries[0].Success {
		t.Error("실패인데 Success=true")
	}
	if entries[0].ErrorMessage == "" {
		t.Error("ErrorMessage 누락")
	}
}

func TestTenancyAdapter_QuotaBlocks(t *testing.T) {
	inner := NewNoopAdapter()
	audit := NewMemoryAuditLog()
	adapter := NewTenancyAdapter(inner, audit)
	quota := NewMemoryQuota(100)
	adapter.SetQuota(quota)

	// hospA 가 이미 100토큰 사용
	quota.RecordUsage("hospA", 100)

	ctx := tenancy.WithTenant(context.Background(), "hospA")
	_, err := adapter.Complete(ctx, &Request{
		UserID:   "u",
		Messages: []*Message{{Role: RoleUser, Content: "안녕하세요 이것은 토큰 카운트를 위한 충분히 긴 의료 상담 메시지입니다 도움이 필요합니다"}},
	})
	if err == nil {
		t.Error("quota 초과인데 통과")
	}

	// 감사 로그에 quota 거부 기록
	entries := audit.EntriesByTenant("hospA")
	if len(entries) != 1 || entries[0].ErrorMessage != "quota exceeded" {
		t.Errorf("audit = %v", entries)
	}
}

func TestTenancyAdapter_QuotaAllowsAndRecordsUsage(t *testing.T) {
	inner := NewNoopAdapter()
	adapter := NewTenancyAdapter(inner, nil)
	quota := NewMemoryQuota(1000)
	adapter.SetQuota(quota)

	ctx := tenancy.WithTenant(context.Background(), "hospA")
	_, err := adapter.Complete(ctx, &Request{
		UserID:   "u",
		Messages: []*Message{{Role: RoleUser, Content: "안녕하세요 이것은 토큰 카운트를 위한 충분히 긴 의료 상담 메시지입니다 도움이 필요합니다"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if quota.Usage("hospA") <= 0 {
		t.Errorf("토큰 사용량 미기록: %d", quota.Usage("hospA"))
	}
}

func TestTenancyAdapter_TotalTokensByTenant(t *testing.T) {
	inner := NewNoopAdapter()
	audit := NewMemoryAuditLog()
	adapter := NewTenancyAdapter(inner, audit)

	ctxA := tenancy.WithTenant(context.Background(), "A")
	ctxB := tenancy.WithTenant(context.Background(), "B")
	for i := 0; i < 3; i++ {
		_, _ = adapter.Complete(ctxA, &Request{Messages: []*Message{{Role: RoleUser, Content: "안녕하세요 이것은 토큰 카운트를 위한 충분히 긴 의료 상담 메시지입니다 도움이 필요합니다"}}})
	}
	for i := 0; i < 2; i++ {
		_, _ = adapter.Complete(ctxB, &Request{Messages: []*Message{{Role: RoleUser, Content: "안녕하세요 이것은 토큰 카운트를 위한 충분히 긴 의료 상담 메시지입니다 도움이 필요합니다"}}})
	}

	if audit.TotalTokensByTenant("A") <= 0 {
		t.Error("A 토큰 = 0")
	}
	if audit.TotalTokensByTenant("B") <= 0 {
		t.Error("B 토큰 = 0")
	}
	// A 와 B 가 분리되어 카운트됨 (격리)
	totalA := audit.TotalTokensByTenant("A")
	totalB := audit.TotalTokensByTenant("B")
	if totalA == totalB {
		t.Error("A 와 B 토큰이 동일 (격리 미적용 가능성)")
	}
}

func TestTenancyAdapter_NilInner(t *testing.T) {
	adapter := NewTenancyAdapter(nil, nil)
	if _, err := adapter.Complete(context.Background(), &Request{}); err == nil {
		t.Error("nil inner 통과")
	}
	if err := adapter.HealthCheck(context.Background()); err == nil {
		t.Error("HealthCheck 통과")
	}
}

func TestTenancyAdapter_NilRequest(t *testing.T) {
	adapter := NewTenancyAdapter(NewNoopAdapter(), nil)
	if _, err := adapter.Complete(context.Background(), nil); err == nil {
		t.Error("nil request 통과")
	}
}

func TestTenancyAdapter_ProviderName(t *testing.T) {
	adapter := NewTenancyAdapter(NewNoopAdapter(), nil)
	if adapter.Provider() != "noop+tenancy" {
		t.Errorf("Provider = %q", adapter.Provider())
	}
}

// failingAdapter 는 항상 에러 반환하는 테스트용 어댑터.
type failingAdapter struct{ err error }

func (f *failingAdapter) Complete(_ context.Context, _ *Request) (*Response, error) {
	return nil, f.err
}
func (f *failingAdapter) Provider() string                   { return "failing" }
func (f *failingAdapter) HealthCheck(_ context.Context) error { return nil }
