// Package service: 다채널 알림 디스패처 (notification ↔ external 어댑터 통합)
//
// notification-service가 backend/shared/external/{sms,email,fcm,kakao}/ 5개 어댑터를
// 유기적으로 호출하는 모세혈관 통합 지점입니다.
//
//   채널 우선순위 정책:
//     critical → FCM Push + SMS + Email + Kakao (4채널 동시)
//     high     → FCM Push + SMS (2채널 폴백)
//     normal   → FCM Push (단일 채널)
//     low      → InApp 또는 Email
package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/manpasik/backend/shared/external/email"
	"github.com/manpasik/backend/shared/external/fcm"
	"github.com/manpasik/backend/shared/external/kakao"
	"github.com/manpasik/backend/shared/external/sms"
	"github.com/manpasik/backend/shared/medical/poise"
)

// ============================================================================
// MultiChannelDispatcher
// ============================================================================

// MultiChannelDispatcher는 우선순위에 따라 적절한 외부 채널로 알림을 분배합니다.
//
// PoISE Loop 통합: 발송 실패/지연 시 자동으로 poise.CollectFeedback 호출
// → 메트릭 누적 → 임계값 초과 시 개선 제안 자동 생성.
type MultiChannelDispatcher struct {
	smsAdapter   sms.Adapter
	emailAdapter email.Adapter
	fcmAdapter   fcm.Adapter
	kakaoAdapter kakao.Adapter

	defaultEmailFrom *email.Address
	poiseLoop        *poise.Loop // optional, nil이면 피드백 수집 비활성
}

// NewMultiChannelDispatcher는 새 디스패처를 생성합니다.
//
// 어댑터를 nil로 전달하면 해당 채널은 비활성화됩니다.
func NewMultiChannelDispatcher(
	smsA sms.Adapter,
	emailA email.Adapter,
	fcmA fcm.Adapter,
	kakaoA kakao.Adapter,
	emailFrom *email.Address,
) *MultiChannelDispatcher {
	if emailFrom == nil {
		emailFrom = &email.Address{Email: "noreply@manpasik.com", Name: "만파식"}
	}
	return &MultiChannelDispatcher{
		smsAdapter:       smsA,
		emailAdapter:     emailA,
		fcmAdapter:       fcmA,
		kakaoAdapter:     kakaoA,
		defaultEmailFrom: emailFrom,
	}
}

// NewMultiChannelDispatcherFromEnv는 환경변수 기반으로 모든 어댑터를 초기화합니다.
func NewMultiChannelDispatcherFromEnv() *MultiChannelDispatcher {
	return NewMultiChannelDispatcher(
		sms.NewFromEnv(),
		email.NewFromEnv(),
		fcm.NewFromEnv(),
		kakao.NewFromEnv(),
		nil,
	)
}

// SetPoISELoop는 PoISE 자동 개선 루프와 dispatcher를 연결합니다.
//
// 주입 후 모든 발송 결과(실패율 / 응답시간)가 PoISE 메트릭으로 자동 수집됩니다.
func (d *MultiChannelDispatcher) SetPoISELoop(loop *poise.Loop) {
	if loop != nil {
		d.poiseLoop = loop
	}
}

// ============================================================================
// 알림 요청 모델
// ============================================================================

// MultiChannelRequest는 다채널 알림 요청입니다.
type MultiChannelRequest struct {
	UserID       string
	Title        string
	Body         string
	Priority     string // "critical" | "high" | "normal" | "low"
	Data         map[string]string

	// 채널별 수신자 (옵션 — 없으면 해당 채널 미사용)
	FCMToken    string
	PhoneNumber string // E.164
	EmailAddr   string
	KakaoPhone  string
	KakaoTemplateCode string
	KakaoVariables    map[string]string
}

// MultiChannelResult는 다채널 발송 결과입니다.
type MultiChannelResult struct {
	UserID       string
	Priority     string
	AttemptedChannels []string
	SuccessChannels   []string
	Errors            map[string]error
	StartedAt         time.Time
	CompletedAt       time.Time
	Duration          time.Duration
}

// HasAnySuccess는 적어도 하나의 채널이 성공했는지 반환합니다.
func (r *MultiChannelResult) HasAnySuccess() bool {
	return len(r.SuccessChannels) > 0
}

// AllSuccess는 시도한 모든 채널이 성공했는지 반환합니다.
func (r *MultiChannelResult) AllSuccess() bool {
	return len(r.AttemptedChannels) > 0 &&
		len(r.AttemptedChannels) == len(r.SuccessChannels)
}

// ============================================================================
// 디스패치 로직
// ============================================================================

// Dispatch는 우선순위에 맞는 채널로 알림을 발송합니다.
//
// critical/high 우선순위는 채널을 동시(병렬) 시도하여 SLA 보장.
func (d *MultiChannelDispatcher) Dispatch(ctx context.Context, req *MultiChannelRequest) (*MultiChannelResult, error) {
	if req == nil {
		return nil, errors.New("request required")
	}
	if req.UserID == "" {
		return nil, errors.New("user_id required")
	}

	result := &MultiChannelResult{
		UserID:    req.UserID,
		Priority:  req.Priority,
		Errors:    make(map[string]error),
		StartedAt: time.Now().UTC(),
	}

	channels := d.selectChannels(req)
	if len(channels) == 0 {
		return result, errors.New("no eligible channels (need fcm_token / phone / email / kakao_phone)")
	}

	// 병렬 발송
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, ch := range channels {
		wg.Add(1)
		go func(channel string) {
			defer wg.Done()
			err := d.sendOne(ctx, channel, req)

			mu.Lock()
			result.AttemptedChannels = append(result.AttemptedChannels, channel)
			if err != nil {
				result.Errors[channel] = err
			} else {
				result.SuccessChannels = append(result.SuccessChannels, channel)
			}
			mu.Unlock()
		}(ch)
	}
	wg.Wait()

	result.CompletedAt = time.Now().UTC()
	result.Duration = result.CompletedAt.Sub(result.StartedAt)

	// PoISE 자동 피드백: 응답 시간 + 실패율을 메트릭으로 수집
	if d.poiseLoop != nil {
		_ = d.poiseLoop.CollectFeedback(&poise.Feedback{
			SessionID: req.UserID,
			Kind:      poise.FeedbackResponseTime,
			Value:     float64(result.Duration.Milliseconds()),
			Comment:   fmt.Sprintf("priority=%s channels_attempted=%d", req.Priority, len(result.AttemptedChannels)),
		})
		// 실패가 있었다면 안전 알람으로도 기록 (운영 가시성 강화)
		if len(result.Errors) > 0 {
			_ = d.poiseLoop.CollectFeedback(&poise.Feedback{
				SessionID: req.UserID,
				Kind:      poise.FeedbackSafetyAlert,
				Value:     float64(len(result.Errors)), // 채널 실패 수
				Comment:   fmt.Sprintf("notification_dispatch_failures: %d/%d", len(result.Errors), len(result.AttemptedChannels)),
			})
		}
	}
	return result, nil
}

// selectChannels는 우선순위와 수신자 정보로 사용할 채널 목록을 결정합니다.
func (d *MultiChannelDispatcher) selectChannels(req *MultiChannelRequest) []string {
	var channels []string

	switch req.Priority {
	case "critical":
		// 가능한 모든 채널 동원
		if req.FCMToken != "" && d.fcmAdapter != nil {
			channels = append(channels, "fcm")
		}
		if req.PhoneNumber != "" && d.smsAdapter != nil {
			channels = append(channels, "sms")
		}
		if req.EmailAddr != "" && d.emailAdapter != nil {
			channels = append(channels, "email")
		}
		if req.KakaoPhone != "" && req.KakaoTemplateCode != "" && d.kakaoAdapter != nil {
			channels = append(channels, "kakao")
		}
	case "high":
		// FCM + SMS 폴백
		if req.FCMToken != "" && d.fcmAdapter != nil {
			channels = append(channels, "fcm")
		}
		if req.PhoneNumber != "" && d.smsAdapter != nil {
			channels = append(channels, "sms")
		}
	case "normal":
		// FCM 단일
		if req.FCMToken != "" && d.fcmAdapter != nil {
			channels = append(channels, "fcm")
		} else if req.EmailAddr != "" && d.emailAdapter != nil {
			channels = append(channels, "email")
		}
	default: // low / 기본
		if req.EmailAddr != "" && d.emailAdapter != nil {
			channels = append(channels, "email")
		}
	}
	return channels
}

// sendOne은 단일 채널로 발송합니다. 모든 어댑터에 nil guard 적용.
func (d *MultiChannelDispatcher) sendOne(ctx context.Context, channel string, req *MultiChannelRequest) error {
	switch channel {
	case "fcm":
		if d.fcmAdapter == nil {
			return errors.New("fcm adapter not configured")
		}
		_, err := d.fcmAdapter.Send(ctx, &fcm.Message{
			Token: req.FCMToken,
			Notification: &fcm.Notification{Title: req.Title, Body: req.Body},
			Data:         req.Data,
		})
		return err

	case "sms":
		if d.smsAdapter == nil {
			return errors.New("sms adapter not configured")
		}
		body := fmt.Sprintf("[%s] %s", req.Title, req.Body)
		_, err := d.smsAdapter.Send(ctx, &sms.Message{
			To:   req.PhoneNumber,
			Body: body,
		})
		return err

	case "email":
		if d.emailAdapter == nil {
			return errors.New("email adapter not configured")
		}
		_, err := d.emailAdapter.Send(ctx, &email.Message{
			From:     d.defaultEmailFrom,
			To:       []*email.Address{{Email: req.EmailAddr}},
			Subject:  req.Title,
			TextBody: req.Body,
		})
		return err

	case "kakao":
		if d.kakaoAdapter == nil {
			return errors.New("kakao adapter not configured")
		}
		_, err := d.kakaoAdapter.Send(ctx, &kakao.Message{
			TemplateCode: req.KakaoTemplateCode,
			To:           req.KakaoPhone,
			Variables:    req.KakaoVariables,
		})
		return err
	}
	return fmt.Errorf("unknown channel: %s", channel)
}

// ============================================================================
// 헬스체크
// ============================================================================

// ChannelHealth는 채널별 헬스체크 결과입니다.
type ChannelHealth struct {
	Channel  string
	Healthy  bool
	Provider string
	Error    string
}

// HealthCheck는 모든 채널의 헬스체크를 수행합니다.
func (d *MultiChannelDispatcher) HealthCheck(ctx context.Context) []*ChannelHealth {
	results := make([]*ChannelHealth, 0, 4)

	if d.fcmAdapter != nil {
		err := d.fcmAdapter.HealthCheck(ctx)
		results = append(results, &ChannelHealth{
			Channel: "fcm", Healthy: err == nil, Provider: d.fcmAdapter.Provider(),
			Error: errMsg(err),
		})
	}
	if d.smsAdapter != nil {
		err := d.smsAdapter.HealthCheck(ctx)
		results = append(results, &ChannelHealth{
			Channel: "sms", Healthy: err == nil, Provider: d.smsAdapter.Provider(),
			Error: errMsg(err),
		})
	}
	if d.emailAdapter != nil {
		err := d.emailAdapter.HealthCheck(ctx)
		results = append(results, &ChannelHealth{
			Channel: "email", Healthy: err == nil, Provider: d.emailAdapter.Provider(),
			Error: errMsg(err),
		})
	}
	if d.kakaoAdapter != nil {
		err := d.kakaoAdapter.HealthCheck(ctx)
		results = append(results, &ChannelHealth{
			Channel: "kakao", Healthy: err == nil, Provider: d.kakaoAdapter.Provider(),
			Error: errMsg(err),
		})
	}
	return results
}

func errMsg(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// ProviderSummary는 활성 채널 ↔ 프로바이더 매핑을 반환합니다.
func (d *MultiChannelDispatcher) ProviderSummary() map[string]string {
	summary := make(map[string]string)
	if d.fcmAdapter != nil {
		summary["fcm"] = d.fcmAdapter.Provider()
	}
	if d.smsAdapter != nil {
		summary["sms"] = d.smsAdapter.Provider()
	}
	if d.emailAdapter != nil {
		summary["email"] = d.emailAdapter.Provider()
	}
	if d.kakaoAdapter != nil {
		summary["kakao"] = d.kakaoAdapter.Provider()
	}
	return summary
}
