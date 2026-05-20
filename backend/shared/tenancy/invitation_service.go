package tenancy

import (
	"context"
	"errors"
	"time"
)

// InvitationServiceConfig 는 초대 서비스 동작 설정.
type InvitationServiceConfig struct {
	// DefaultTTL 은 초대 유효 기간 (기본 7일).
	DefaultTTL time.Duration
	// TokenGenerator 는 토큰 생성 함수 (테스트 시 결정적 토큰 주입).
	// nil 이면 crypto/rand 기반 generateToken.
	TokenGenerator func() (string, error)
	// Now 는 시간 소스 (테스트용). nil 이면 time.Now.
	Now func() time.Time
}

// InvitationNotifier 는 초대 발급 시 외부 알림(이메일/SMS/카카오톡 등) 발송 담당.
//
// 인터페이스로 추상화하여 notification-service 또는 외부 어댑터가 구현.
// SetNotifier 미호출 시 알림 발송 안 함 (옵셔널).
type InvitationNotifier interface {
	NotifyInvitation(ctx context.Context, inv Invitation) error
}

// InvitationService 는 초대 발급/수락/취소 비즈니스 로직.
//
// 의존: InvitationStore + MembershipStore + PolicyEngine (선택, 권한 체크용).
type InvitationService struct {
	invStore  InvitationStore
	memStore  MembershipStore
	cfg       InvitationServiceConfig
	policy    *PolicyEngine // nil 가능: 권한 체크 없이 어떤 사용자든 invite 가능
	notifier  InvitationNotifier
	webhook   *WebhookDispatcher // 옵션 — invite 이벤트 webhook 발송
}

// NewInvitationService 생성.
func NewInvitationService(invStore InvitationStore, memStore MembershipStore,
	cfg InvitationServiceConfig) (*InvitationService, error) {
	if invStore == nil || memStore == nil {
		return nil, errors.New("invStore/memStore 필수")
	}
	if cfg.DefaultTTL <= 0 {
		cfg.DefaultTTL = 7 * 24 * time.Hour
	}
	if cfg.TokenGenerator == nil {
		cfg.TokenGenerator = generateToken
	}
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	return &InvitationService{
		invStore: invStore,
		memStore: memStore,
		cfg:      cfg,
	}, nil
}

// SetPolicyEngine 은 invite 발급 시 권한 검사용 PolicyEngine 등록.
//
// 미설정 시 모든 사용자가 invite 발급 가능 (개발용). 운영 환경에서는 반드시
// 등록하여 ActionAdmin 또는 ActionShare 권한 검사.
func (s *InvitationService) SetPolicyEngine(p *PolicyEngine) {
	s.policy = p
}

// SetNotifier 는 초대 발급 시 외부 알림 발송 담당자 등록 (선택).
//
// 발송은 비동기 goroutine 으로 호출되어 Invite() 응답을 막지 않음. 발송 실패는
// 로깅 차원이므로 InvitationService 자체는 에러를 무시.
func (s *InvitationService) SetNotifier(n InvitationNotifier) {
	s.notifier = n
}

// SetWebhookDispatcher 는 invite 이벤트 (Created/Accepted/Revoked) webhook
// 발송기 등록 (Phase AJ-2). 미설정 시 webhook 미발송.
//
// 발송은 dispatcher.DispatchAsync 로 비동기 — service 응답 시간에 영향 없음.
func (s *InvitationService) SetWebhookDispatcher(d *WebhookDispatcher) {
	s.webhook = d
}

// InviteRequest 는 초대 발급 요청.
type InviteRequest struct {
	InviterID   string
	TenantID    TenantID
	InviteeHint string
	Role        TenantRole
	TTL         time.Duration // 0 이면 cfg.DefaultTTL
}

// Invite 는 초대 발급. 권한 검사 후 토큰 생성.
//
// 권한 체크: PolicyEngine 등록 시 inviter 가 tenant 의 admin/owner 또는 share 권한 보유 확인.
func (s *InvitationService) Invite(req InviteRequest) (*Invitation, error) {
	if req.InviterID == "" || req.TenantID.IsZero() {
		return nil, errors.New("InviterID/TenantID 필수")
	}
	if !req.Role.IsKnown() {
		return nil, errors.New("알 수 없는 역할: " + string(req.Role))
	}

	// 권한 체크 (PolicyEngine 있을 때만)
	if s.policy != nil {
		// inviter 가 tenant 의 admin 권한이 있는지 확인
		// invite 는 ActionAdmin 권한 필요 (멤버 추가는 admin 작업)
		res := &Resource{TenantID: req.TenantID}
		d := s.policy.Evaluate(req.InviterID, req.TenantID, res, ActionAdmin)
		if !d.Allowed {
			return nil, errors.New("초대 권한 없음: " + d.Reason)
		}
	}

	token, err := s.cfg.TokenGenerator()
	if err != nil {
		return nil, err
	}

	ttl := req.TTL
	if ttl <= 0 {
		ttl = s.cfg.DefaultTTL
	}
	now := s.cfg.Now()
	inv := Invitation{
		Token:       token,
		TenantID:    req.TenantID,
		InviterID:   req.InviterID,
		InviteeHint: req.InviteeHint,
		Role:        req.Role,
		Status:      InvitationPending,
		IssuedAt:    now,
		ExpiresAt:   now.Add(ttl),
	}
	if err := s.invStore.Add(inv); err != nil {
		return nil, err
	}

	// 알림 발송 (비동기). 실패는 로깅 차원이며 invite 발급 자체는 성공.
	if s.notifier != nil {
		notifierCopy := s.notifier
		invCopy := inv
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_ = notifierCopy.NotifyInvitation(ctx, invCopy)
		}()
	}

	// Webhook 발송 (Phase AJ-2)
	if s.webhook != nil {
		s.webhook.DispatchAsync(Event{
			Type:     EventInvitationCreated,
			TenantID: string(inv.TenantID),
			ActorID:  inv.InviterID,
			Payload: map[string]string{
				"role":         string(inv.Role),
				"invitee_hint": inv.InviteeHint,
			},
		})
	}
	return &inv, nil
}

// Accept 는 초대 수락. 토큰 검증 → 멤버십 생성 → 초대 상태 갱신.
//
// 동일 사용자가 동일 초대를 여러 번 수락 시 첫 번째만 성공 (멱등 X — 명시적 에러).
// 운영 정책: 이미 멤버인 경우 초대 자체를 무효 처리하고 에러 반환 (중복 가입 방지).
func (s *InvitationService) Accept(token, accepterID string) (*Membership, error) {
	if token == "" || accepterID == "" {
		return nil, errors.New("token/accepterID 필수")
	}

	inv, err := s.invStore.Get(token)
	if err != nil {
		return nil, err
	}
	if inv.Status != InvitationPending {
		return nil, ErrInvitationConsumed
	}
	if inv.IsExpired() {
		// 만료 표시 갱신 (실패해도 무시)
		inv.Status = InvitationExpired
		_ = s.invStore.Update(*inv)
		return nil, ErrInvitationExpired
	}

	// 이미 활성 멤버인지 확인 → 중복 추가 차단
	if existing, err := s.memStore.Get(accepterID, inv.TenantID); err == nil && existing.Active {
		return nil, errors.New("이미 활성 멤버")
	}

	m := Membership{
		UserID:   accepterID,
		TenantID: inv.TenantID,
		Role:     inv.Role,
		Active:   true,
		JoinedAt: s.cfg.Now().Unix(),
	}
	if err := s.memStore.Add(m); err != nil {
		return nil, err
	}

	// 초대 상태 갱신
	now := s.cfg.Now()
	inv.Status = InvitationAccepted
	inv.AcceptedBy = accepterID
	inv.AcceptedAt = &now
	if err := s.invStore.Update(*inv); err != nil {
		return nil, err
	}

	// Webhook 발송 (Phase AJ-2)
	if s.webhook != nil {
		s.webhook.DispatchAsync(Event{
			Type:     EventInvitationAccepted,
			TenantID: string(inv.TenantID),
			UserID:   accepterID,
			ActorID:  accepterID,
			Payload: map[string]string{
				"role":       string(inv.Role),
				"inviter_id": inv.InviterID,
			},
		})
		s.webhook.DispatchAsync(Event{
			Type:     EventMembershipCreated,
			TenantID: string(inv.TenantID),
			UserID:   accepterID,
			ActorID:  accepterID,
			Payload: map[string]string{
				"role": string(m.Role),
			},
		})
	}
	return &m, nil
}

// Revoke 는 초대 취소. 발급자 또는 admin 만 가능.
func (s *InvitationService) Revoke(token, revokerID string) error {
	inv, err := s.invStore.Get(token)
	if err != nil {
		return err
	}
	if inv.Status != InvitationPending {
		return ErrInvitationConsumed
	}
	// 발급자 본인 또는 PolicyEngine 검사 통과
	if inv.InviterID != revokerID && s.policy != nil {
		res := &Resource{TenantID: inv.TenantID}
		if d := s.policy.Evaluate(revokerID, inv.TenantID, res, ActionAdmin); !d.Allowed {
			return errors.New("취소 권한 없음")
		}
	}
	inv.Status = InvitationRevoked
	if err := s.invStore.Update(*inv); err != nil {
		return err
	}

	// Webhook 발송 (Phase AJ-2)
	if s.webhook != nil {
		s.webhook.DispatchAsync(Event{
			Type:     EventInvitationRevoked,
			TenantID: string(inv.TenantID),
			ActorID:  revokerID,
			Payload: map[string]string{
				"original_inviter": inv.InviterID,
			},
		})
	}
	return nil
}

// ListPendingByTenant 는 활성 (pending + 만료 안 됨) 초대만 반환.
func (s *InvitationService) ListPendingByTenant(tenantID TenantID) []*Invitation {
	all := s.invStore.ListByTenant(tenantID)
	out := make([]*Invitation, 0, len(all))
	for _, inv := range all {
		if inv.IsActive() {
			out = append(out, inv)
		}
	}
	return out
}
