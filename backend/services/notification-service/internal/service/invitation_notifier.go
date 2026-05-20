package service

import (
	"context"
	"fmt"
	"time"

	"github.com/manpasik/backend/shared/tenancy"
)

// InvitationNotifierAdapter 는 tenancy.InvitationNotifier 인터페이스를
// MultiChannelDispatcher 로 구현하는 어댑터.
//
// 사용 예 (gateway/admin main.go):
//
//	dispatcher := service.NewMultiChannelDispatcherFromEnv()
//	notifier := service.NewInvitationNotifierAdapter(dispatcher, lookupContact)
//	invSvc.SetNotifier(notifier)
//
// lookupContact 함수는 invitee_hint (이메일 등) 를 받아 채널별 수신자 정보로
// 변환. 누락된 채널은 자동으로 발송 대상에서 제외.
// LocaleResolver 는 invitee_hint 를 받아 적절한 locale 코드 반환.
//
// 운영 예: user-service.GetByEmail → user.LangCode → 매핑 (ko/en/ja/zh).
// 빈 문자열 또는 미지원 locale 반환 시 default locale 사용.
type LocaleResolver func(ctx context.Context, hint string) string

type InvitationNotifierAdapter struct {
	dispatcher     *MultiChannelDispatcher
	contactLookup  ContactLookupFunc
	localeResolver LocaleResolver // 옵션 — 설정 시 매 알림마다 호출
	locale         string         // 기본 "ko" — SetLocale 으로 변경
	// titleFormat/bodyFormat 은 SetTitleFormat 호출 시에만 설정 (직접 오버라이드용).
	// 미설정 시 locale 기반 템플릿 사용.
	titleFormat string
	bodyFormat  string
}

// invitationLocaleTemplates 는 4개 언어별 (title/body) 템플릿.
//
// title 의 %s = tenant_id, body 의 %s = (tenant, role, token, expires, token).
var invitationLocaleTemplates = map[string]struct {
	Title string
	Body  string
}{
	"ko": {
		Title: "[만파식] %s 조직 초대",
		Body: "%s 조직에 %s 역할로 초대되었습니다.\n" +
			"수락 토큰: %s\n" +
			"만료: %s\n" +
			"앱에서 [설정 > 조직 초대 수락] 으로 입력하거나, " +
			"딥링크 manpasik://invite/%s 를 클릭하세요.",
	},
	"en": {
		Title: "[ManPaSik] You're invited to %s",
		Body: "You've been invited to %s as %s.\n" +
			"Acceptance token: %s\n" +
			"Expires: %s\n" +
			"Open the app and go to [Settings > Accept Invitation], " +
			"or tap the deep link manpasik://invite/%s.",
	},
	"ja": {
		Title: "[万波息] %s 組織への招待",
		Body: "%s 組織に %s 役割で招待されました。\n" +
			"受諾トークン: %s\n" +
			"有効期限: %s\n" +
			"アプリで [設定 > 招待を受諾] から入力するか、" +
			"ディープリンク manpasik://invite/%s をタップしてください。",
	},
	"zh": {
		Title: "[万波息] 邀请您加入 %s",
		Body: "您已被邀请加入 %s,角色为 %s。\n" +
			"接受令牌: %s\n" +
			"过期时间: %s\n" +
			"在应用中前往 [设置 > 接受邀请] 输入,或点击深链接 manpasik://invite/%s。",
	},
}

// ContactLookupFunc 는 invitee_hint → 채널 수신자 변환 함수.
//
// invitee_hint 는 이메일/전화/사용자ID 등 발급자가 임의로 입력한 식별자.
// notification-service 는 이를 user-service/profile 조회로 풀어서
// 적절한 채널 (FCM 토큰 / 이메일 / 카카오 폰) 을 결정.
type ContactLookupFunc func(ctx context.Context, hint string) (InviteeContact, error)

// InviteeContact 는 알림 발송용 수신자 정보.
type InviteeContact struct {
	UserID    string
	FCMToken  string
	Email     string
	Phone     string
	KakaoPhone string
}

// NewInvitationNotifierAdapter 생성. lookup=nil 이면 invitee_hint 를
// EmailAddr 으로 직접 사용 (간단한 fallback). 기본 locale=ko.
func NewInvitationNotifierAdapter(dispatcher *MultiChannelDispatcher,
	lookup ContactLookupFunc) *InvitationNotifierAdapter {
	if lookup == nil {
		lookup = defaultEmailFallbackLookup
	}
	return &InvitationNotifierAdapter{
		dispatcher:    dispatcher,
		contactLookup: lookup,
		locale:        "ko",
	}
}

// SetLocale 은 알림 기본 언어 변경 (Phase AF-3). 지원: ko, en, ja, zh.
//
// 미지원 코드 또는 빈 값은 ko 로 fallback.
func (a *InvitationNotifierAdapter) SetLocale(locale string) {
	if _, ok := invitationLocaleTemplates[locale]; ok {
		a.locale = locale
		return
	}
	a.locale = "ko"
}

// SetLocaleResolver 는 invitee 별 locale 자동 감지 함수 등록 (Phase AG-3).
//
// 미설정 시 default locale 사용. 설정 시 매 알림 발송 직전 호출되어 결과로
// title/body 템플릿 결정.
func (a *InvitationNotifierAdapter) SetLocaleResolver(resolver LocaleResolver) {
	a.localeResolver = resolver
}

// Locale 반환 (default).
func (a *InvitationNotifierAdapter) Locale() string { return a.locale }

// resolveLocale 은 hint 와 ctx 로부터 locale 결정.
//
// 우선순위:
//  1. localeResolver 가 설정되었고 결과가 지원 locale 이면 그 locale
//  2. 그렇지 않으면 a.locale (default)
func (a *InvitationNotifierAdapter) resolveLocale(ctx context.Context, hint string) string {
	if a.localeResolver != nil {
		if loc := a.localeResolver(ctx, hint); loc != "" {
			if _, ok := invitationLocaleTemplates[loc]; ok {
				return loc
			}
		}
	}
	return a.locale
}

// SetTitleFormat / SetBodyFormat — 운영자 커스터마이징용 (locale 템플릿 무시).
//
// 미설정 시 locale 기반 자동 템플릿 사용.
func (a *InvitationNotifierAdapter) SetTitleFormat(format string) {
	if format != "" {
		a.titleFormat = format
	}
}

func (a *InvitationNotifierAdapter) SetBodyFormat(format string) {
	if format != "" {
		a.bodyFormat = format
	}
}

// resolveTemplates 는 default locale 기반 (title, body) 반환 (테스트 호환용).
//
// 운영에서는 resolveTemplatesFor 가 invitee 별 자동 감지 사용.
func (a *InvitationNotifierAdapter) resolveTemplates() (title, body string) {
	return a.resolveTemplatesFor(context.Background(), "")
}

// resolveTemplatesFor 는 ctx/hint 로 locale 자동 감지 후 (title, body) 반환.
//
// localeResolver 가 설정되어 있으면 그 결과 우선; 명시적 SetTitleFormat/
// SetBodyFormat 가 호출되어 있으면 locale 무시하고 override 사용.
func (a *InvitationNotifierAdapter) resolveTemplatesFor(ctx context.Context, hint string) (title, body string) {
	loc := a.resolveLocale(ctx, hint)
	t, ok := invitationLocaleTemplates[loc]
	if !ok {
		t = invitationLocaleTemplates["ko"]
	}
	if a.titleFormat != "" {
		title = a.titleFormat
	} else {
		title = t.Title
	}
	if a.bodyFormat != "" {
		body = a.bodyFormat
	} else {
		body = t.Body
	}
	return title, body
}

// NotifyInvitation 은 tenancy.InvitationNotifier 구현.
//
// invitee_hint 가 비어있으면 알림 미발송 (수신자 미상). dispatcher 가 nil 이면
// 즉시 nil 반환 (운영자가 명시적으로 알림 비활성화한 경우).
func (a *InvitationNotifierAdapter) NotifyInvitation(ctx context.Context,
	inv tenancy.Invitation) error {
	if a.dispatcher == nil {
		return nil
	}
	if inv.InviteeHint == "" {
		return nil // 수신자 미상이면 알림 발송 안 함
	}

	contact, err := a.contactLookup(ctx, inv.InviteeHint)
	if err != nil {
		return fmt.Errorf("contact lookup failed: %w", err)
	}
	if contact.UserID == "" && contact.FCMToken == "" &&
		contact.Email == "" && contact.Phone == "" && contact.KakaoPhone == "" {
		// 수신 가능한 채널이 하나도 없음 → 발송 생략
		return nil
	}

	titleFmt, bodyFmt := a.resolveTemplatesFor(ctx, inv.InviteeHint)
	title := fmt.Sprintf(titleFmt, string(inv.TenantID))
	body := fmt.Sprintf(bodyFmt,
		string(inv.TenantID),
		string(inv.Role),
		inv.Token,
		inv.ExpiresAt.Format(time.RFC3339),
		inv.Token,
	)

	req := &MultiChannelRequest{
		UserID:      contact.UserID,
		Title:       title,
		Body:        body,
		Priority:    "high",
		FCMToken:    contact.FCMToken,
		EmailAddr:   contact.Email,
		PhoneNumber: contact.Phone,
		KakaoPhone:  contact.KakaoPhone,
		Data: map[string]string{
			"type":        "tenant_invitation",
			"tenant_id":   string(inv.TenantID),
			"invite_token": inv.Token,
		},
	}
	// UserID 가 비어있어도 dispatcher 가 거부하지 않도록 hint 자체를 ID 로 사용
	if req.UserID == "" {
		req.UserID = inv.InviteeHint
	}

	if _, err := a.dispatcher.Dispatch(ctx, req); err != nil {
		return fmt.Errorf("dispatch failed: %w", err)
	}
	return nil
}

// defaultEmailFallbackLookup 은 hint 를 이메일로 간주하고 그 외 채널은 비움.
//
// 운영 환경에서는 user-service.GetByEmail 등으로 교체 권장.
func defaultEmailFallbackLookup(_ context.Context, hint string) (InviteeContact, error) {
	c := InviteeContact{}
	if isLikelyEmail(hint) {
		c.Email = hint
	}
	return c, nil
}

// isLikelyEmail 은 단순 휴리스틱 — '@' 가 포함되면 이메일로 간주.
func isLikelyEmail(s string) bool {
	for _, r := range s {
		if r == '@' {
			return true
		}
	}
	return false
}
