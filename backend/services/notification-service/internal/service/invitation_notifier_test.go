package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/tenancy"
)

func TestInvitationNotifierAdapter_NoDispatcher(t *testing.T) {
	a := NewInvitationNotifierAdapter(nil, nil)
	err := a.NotifyInvitation(context.Background(), tenancy.Invitation{
		InviteeHint: "x@y.com",
	})
	if err != nil {
		t.Errorf("nil dispatcher 에러 발생: %v", err)
	}
}

func TestInvitationNotifierAdapter_NoHint(t *testing.T) {
	disp := NewMultiChannelDispatcher(nil, nil, nil, nil, nil)
	a := NewInvitationNotifierAdapter(disp, nil)
	err := a.NotifyInvitation(context.Background(), tenancy.Invitation{
		Token: "tok-1", TenantID: "t", Role: tenancy.TenantRoleMember,
	})
	if err != nil {
		t.Errorf("hint 비어있는데 에러: %v", err)
	}
}

func TestInvitationNotifierAdapter_DefaultEmailFallback(t *testing.T) {
	// disp 가 noop 이어도 lookup → email 매핑이 동작함을 검증
	a := NewInvitationNotifierAdapter(nil, nil)
	contact, err := a.contactLookup(context.Background(), "user@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if contact.Email != "user@example.com" {
		t.Errorf("Email = %q", contact.Email)
	}
}

func TestInvitationNotifierAdapter_NonEmailHint(t *testing.T) {
	a := NewInvitationNotifierAdapter(nil, nil)
	contact, _ := a.contactLookup(context.Background(), "010-1234-5678")
	if contact.Email != "" {
		t.Errorf("이메일이 아닌데 Email 설정됨: %q", contact.Email)
	}
}

func TestInvitationNotifierAdapter_LookupError(t *testing.T) {
	disp := NewMultiChannelDispatcher(nil, nil, nil, nil, nil)
	failingLookup := func(_ context.Context, _ string) (InviteeContact, error) {
		return InviteeContact{}, errors.New("user-service down")
	}
	a := NewInvitationNotifierAdapter(disp, failingLookup)
	err := a.NotifyInvitation(context.Background(), tenancy.Invitation{
		Token: "tok-2", InviteeHint: "x@y.com",
	})
	if err == nil {
		t.Error("lookup 실패에 에러 없음")
	}
	if !strings.Contains(err.Error(), "contact lookup") {
		t.Errorf("err = %v", err)
	}
}

func TestInvitationNotifierAdapter_EmptyContactSkipsDispatch(t *testing.T) {
	disp := NewMultiChannelDispatcher(nil, nil, nil, nil, nil)
	emptyLookup := func(_ context.Context, _ string) (InviteeContact, error) {
		return InviteeContact{}, nil
	}
	a := NewInvitationNotifierAdapter(disp, emptyLookup)
	err := a.NotifyInvitation(context.Background(), tenancy.Invitation{
		Token: "tok-3", InviteeHint: "ghost",
	})
	if err != nil {
		t.Errorf("빈 contact 에 에러: %v", err)
	}
}

func TestInvitationNotifierAdapter_BodyContainsToken(t *testing.T) {
	a := NewInvitationNotifierAdapter(nil, nil)
	body := buildBody(a, tenancy.Invitation{
		Token:     "abc123",
		TenantID:  "hospA",
		Role:      tenancy.TenantRoleMedicalStaff,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	})
	if !strings.Contains(body, "abc123") {
		t.Errorf("body 에 token 누락: %q", body)
	}
	if !strings.Contains(body, "hospA") {
		t.Errorf("body 에 tenant 누락: %q", body)
	}
	if !strings.Contains(body, "medical_staff") {
		t.Errorf("body 에 role 누락: %q", body)
	}
	if !strings.Contains(body, "manpasik://invite/abc123") {
		t.Errorf("body 에 딥링크 누락: %q", body)
	}
}

func TestInvitationNotifierAdapter_DefaultLocale(t *testing.T) {
	a := NewInvitationNotifierAdapter(nil, nil)
	if a.Locale() != "ko" {
		t.Errorf("default locale = %q, want ko", a.Locale())
	}
}

func TestInvitationNotifierAdapter_SetLocale_Supported(t *testing.T) {
	a := NewInvitationNotifierAdapter(nil, nil)
	for _, loc := range []string{"ko", "en", "ja", "zh"} {
		a.SetLocale(loc)
		if a.Locale() != loc {
			t.Errorf("SetLocale(%q) → %q", loc, a.Locale())
		}
	}
}

func TestInvitationNotifierAdapter_SetLocale_Unsupported(t *testing.T) {
	a := NewInvitationNotifierAdapter(nil, nil)
	a.SetLocale("zz") // 미지원
	if a.Locale() != "ko" {
		t.Errorf("미지원 locale → %q (want ko fallback)", a.Locale())
	}
}

func TestInvitationNotifierAdapter_LocaleSpecificBody(t *testing.T) {
	cases := map[string]string{
		"ko": "조직에",
		"en": "invited",
		"ja": "招待",
		"zh": "邀请",
	}
	for loc, expectedSub := range cases {
		a := NewInvitationNotifierAdapter(nil, nil)
		a.SetLocale(loc)
		body := buildBody(a, tenancy.Invitation{
			Token: "tok-1", TenantID: "hospA", Role: tenancy.TenantRoleMember,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		})
		if !strings.Contains(body, expectedSub) {
			t.Errorf("locale=%q body 에 %q 누락:\n%s", loc, expectedSub, body)
		}
	}
}

func TestInvitationNotifierAdapter_LocaleSpecificTitle(t *testing.T) {
	cases := map[string]string{
		"ko": "만파식",
		"en": "ManPaSik",
		"ja": "万波息",
		"zh": "万波息",
	}
	for loc, expectedSub := range cases {
		a := NewInvitationNotifierAdapter(nil, nil)
		a.SetLocale(loc)
		titleFmt, _ := a.resolveTemplates()
		if !strings.Contains(titleFmt, expectedSub) {
			t.Errorf("locale=%q title 에 %q 누락: %s", loc, expectedSub, titleFmt)
		}
	}
}

func TestInvitationNotifierAdapter_OverrideTakesPriority(t *testing.T) {
	a := NewInvitationNotifierAdapter(nil, nil)
	a.SetLocale("en")
	a.SetTitleFormat("[Custom] %s")
	titleFmt, _ := a.resolveTemplates()
	if titleFmt != "[Custom] %s" {
		t.Errorf("override 무시: %q", titleFmt)
	}
}

func TestInvitationNotifierAdapter_AllLocalesContainDeepLink(t *testing.T) {
	for _, loc := range []string{"ko", "en", "ja", "zh"} {
		a := NewInvitationNotifierAdapter(nil, nil)
		a.SetLocale(loc)
		body := buildBody(a, tenancy.Invitation{
			Token: "abc", TenantID: "t", Role: tenancy.TenantRoleMember,
			ExpiresAt: time.Now(),
		})
		if !strings.Contains(body, "manpasik://invite/abc") {
			t.Errorf("locale=%q 에 딥링크 누락:\n%s", loc, body)
		}
	}
}

func TestLocaleResolver_Override(t *testing.T) {
	a := NewInvitationNotifierAdapter(nil, nil)
	a.SetLocale("ko") // default
	a.SetLocaleResolver(func(_ context.Context, hint string) string {
		if strings.Contains(hint, "@en.") {
			return "en"
		}
		return ""
	})

	titleFmt, _ := a.resolveTemplatesFor(context.Background(), "user@en.test")
	if !strings.Contains(titleFmt, "ManPaSik") {
		t.Errorf("resolver=en 인데 ko 사용: %q", titleFmt)
	}

	// 미해결 hint 는 default ko
	titleFmt2, _ := a.resolveTemplatesFor(context.Background(), "user@kr.test")
	if !strings.Contains(titleFmt2, "만파식") {
		t.Errorf("미해결 hint 에서 fallback 안됨: %q", titleFmt2)
	}
}

func TestLocaleResolver_UnsupportedFallback(t *testing.T) {
	a := NewInvitationNotifierAdapter(nil, nil)
	a.SetLocale("en")
	a.SetLocaleResolver(func(_ context.Context, _ string) string {
		return "unsupported-locale"
	})

	titleFmt, _ := a.resolveTemplatesFor(context.Background(), "x")
	// resolver 의 미지원 결과는 무시 → default (en) 사용
	if !strings.Contains(titleFmt, "ManPaSik") {
		t.Errorf("미지원 resolver 결과 무시 안됨: %q", titleFmt)
	}
}

func TestLocaleResolver_EmptyResultUsesDefault(t *testing.T) {
	a := NewInvitationNotifierAdapter(nil, nil)
	a.SetLocale("ja")
	a.SetLocaleResolver(func(_ context.Context, _ string) string { return "" })

	titleFmt, _ := a.resolveTemplatesFor(context.Background(), "x")
	if !strings.Contains(titleFmt, "万波息") {
		t.Errorf("빈 resolver 결과 → default 미사용: %q", titleFmt)
	}
}

func TestLocaleResolver_NotSet_UsesDefault(t *testing.T) {
	a := NewInvitationNotifierAdapter(nil, nil)
	a.SetLocale("zh")
	titleFmt, _ := a.resolveTemplatesFor(context.Background(), "x")
	if !strings.Contains(titleFmt, "万波息") {
		t.Errorf("resolver 미설정 시 default 미사용: %q", titleFmt)
	}
}

// buildBody 는 어댑터 내부 포맷을 테스트용으로 호출 (invitation_notifier.go 와 동일).
func buildBody(a *InvitationNotifierAdapter, inv tenancy.Invitation) string {
	_, bodyFmt := a.resolveTemplates()
	return fmt.Sprintf(bodyFmt,
		string(inv.TenantID),
		string(inv.Role),
		inv.Token,
		inv.ExpiresAt.Format(time.RFC3339),
		inv.Token,
	)
}
