// Package service는 notification-service의 비즈니스 로직을 구현합니다.
package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/manpasik/backend/shared/errors"
	"go.uber.org/zap"
)

// NotificationChannel은 알림 채널 유형입니다.
type NotificationChannel int

const (
	ChannelUnknown NotificationChannel = iota
	ChannelPush
	ChannelEmail
	ChannelSMS
	ChannelInApp
)

// NotificationType은 알림 종류입니다.
type NotificationType int

const (
	TypeUnknown      NotificationType = iota
	TypeMeasurement                    // 측정
	TypeHealthAlert                    // 건강 알림
	TypeAppointment                    // 예약
	TypePrescription                   // 처방
	TypeCommunity                      // 커뮤니티
	TypeSystem                         // 시스템
	TypePromotion                      // 프로모션
)

// NotificationPriority는 알림 우선순위입니다.
type NotificationPriority int

const (
	PriorityUnknown NotificationPriority = iota
	PriorityLow
	PriorityNormal
	PriorityHigh
	PriorityUrgent
)

// Notification은 알림 도메인 객체입니다.
type Notification struct {
	ID        string
	UserID    string
	Type      NotificationType
	Channel   NotificationChannel
	Priority  NotificationPriority
	Title     string
	Body      string
	Data      string
	IsRead    bool
	CreatedAt time.Time
	ReadAt    *time.Time
}

// NotificationPreferences는 사용자 알림 설정입니다.
type NotificationPreferences struct {
	UserID             string
	PushEnabled        bool
	EmailEnabled       bool
	SMSEnabled         bool
	InAppEnabled       bool
	HealthAlertEnabled bool
	CoachingEnabled    bool
	PromotionEnabled   bool
	QuietHoursStart    string // "HH:MM"
	QuietHoursEnd      string
	Language           string
}

// NotificationRepository는 알림 저장소 인터페이스입니다.
type NotificationRepository interface {
	Save(ctx context.Context, n *Notification) error
	FindByID(ctx context.Context, id string) (*Notification, error)
	FindByUserID(ctx context.Context, userID string, typeFilter NotificationType, unreadOnly bool, limit, offset int) ([]*Notification, int, error)
	MarkAsRead(ctx context.Context, id string) error
	MarkAllAsRead(ctx context.Context, userID string) (int, error)
	GetUnreadCount(ctx context.Context, userID string) (int, map[string]int, error)
}

// PreferencesRepository는 알림 설정 저장소 인터페이스입니다.
type PreferencesRepository interface {
	Save(ctx context.Context, pref *NotificationPreferences) error
	FindByUserID(ctx context.Context, userID string) (*NotificationPreferences, error)
}

// PushSender는 실제 푸시 알림 전송 인터페이스입니다 (FCM 등).
type PushSender interface {
	SendPush(ctx context.Context, userID, title, body, data string) error
}

// EmailSender는 실제 이메일 전송 인터페이스입니다 (SMTP 등).
type EmailSender interface {
	SendEmail(ctx context.Context, userID, subject, body string) error
}

// NotificationService는 알림 서비스 핵심 로직입니다.
type NotificationService struct {
	log         *zap.Logger
	notiRepo    NotificationRepository
	prefRepo    PreferencesRepository
	pushSender  PushSender  // optional: nil이면 푸시 미발송
	emailSender EmailSender // optional: nil이면 이메일 미발송
	mu          sync.RWMutex
}

// NewNotificationService는 NotificationService를 생성합니다.
func NewNotificationService(log *zap.Logger, notiRepo NotificationRepository, prefRepo PreferencesRepository) *NotificationService {
	return &NotificationService{
		log:      log,
		notiRepo: notiRepo,
		prefRepo: prefRepo,
	}
}

// SetPushSender는 푸시 알림 전송기를 설정합니다 (optional).
func (s *NotificationService) SetPushSender(ps PushSender) {
	s.pushSender = ps
}

// SetEmailSender는 이메일 전송기를 설정합니다 (optional).
func (s *NotificationService) SetEmailSender(es EmailSender) {
	s.emailSender = es
}

// SendNotification은 알림을 발송합니다.
func (s *NotificationService) SendNotification(ctx context.Context, userID string, nType NotificationType, channel NotificationChannel, priority NotificationPriority, title, body string, data string) (*Notification, error) {
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}
	if title == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}

	// 채널이 0이면 사용자 설정에 따라 자동 선택
	if channel == ChannelUnknown {
		pref, err := s.prefRepo.FindByUserID(ctx, userID)
		if err != nil || pref == nil {
			channel = ChannelInApp // 기본값
		} else {
			channel = s.selectBestChannel(pref, priority)
		}
	}

	// 우선순위 기본값
	if priority == PriorityUnknown {
		priority = PriorityNormal
	}

	noti := &Notification{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      nType,
		Channel:   channel,
		Priority:  priority,
		Title:     title,
		Body:      body,
		Data:      data,
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	if err := s.notiRepo.Save(ctx, noti); err != nil {
		return nil, fmt.Errorf("알림 저장 실패: %w", err)
	}

	// 실제 채널별 전송 (optional, 실패해도 알림 저장은 유지)
	switch channel {
	case ChannelPush:
		if s.pushSender != nil {
			if err := s.pushSender.SendPush(ctx, userID, title, body, data); err != nil {
				s.log.Warn("FCM 푸시 전송 실패 (알림은 저장됨)", zap.Error(err), zap.String("notification_id", noti.ID))
			}
		}
	case ChannelEmail:
		if s.emailSender != nil {
			if err := s.emailSender.SendEmail(ctx, userID, title, body); err != nil {
				s.log.Warn("이메일 전송 실패 (알림은 저장됨)", zap.Error(err), zap.String("notification_id", noti.ID))
			}
		}
	}

	s.log.Info("알림 발송 완료",
		zap.String("notification_id", noti.ID),
		zap.String("user_id", userID),
		zap.Int("channel", int(channel)),
		zap.Int("type", int(nType)),
	)

	return noti, nil
}

// selectBestChannel은 사용자 설정과 우선순위에 따라 최적의 채널을 선택합니다.
func (s *NotificationService) selectBestChannel(pref *NotificationPreferences, priority NotificationPriority) NotificationChannel {
	if priority >= PriorityHigh && pref.PushEnabled {
		return ChannelPush
	}
	if pref.InAppEnabled {
		return ChannelInApp
	}
	if pref.PushEnabled {
		return ChannelPush
	}
	if pref.EmailEnabled {
		return ChannelEmail
	}
	return ChannelInApp
}

// ListNotifications는 알림 목록을 조회합니다.
func (s *NotificationService) ListNotifications(ctx context.Context, userID string, typeFilter NotificationType, unreadOnly bool, limit, offset int) ([]*Notification, int, int, error) {
	if userID == "" {
		return nil, 0, 0, apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}
	if limit <= 0 {
		limit = 20
	}

	notis, total, err := s.notiRepo.FindByUserID(ctx, userID, typeFilter, unreadOnly, limit, offset)
	if err != nil {
		return nil, 0, 0, err
	}

	// 읽지 않은 수 조회
	unreadCount, _, err := s.notiRepo.GetUnreadCount(ctx, userID)
	if err != nil {
		return nil, 0, 0, err
	}

	return notis, total, unreadCount, nil
}

// MarkAsRead는 특정 알림을 읽음 처리합니다.
func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID string) error {
	if notificationID == "" {
		return apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}
	return s.notiRepo.MarkAsRead(ctx, notificationID)
}

// MarkAllAsRead는 사용자의 모든 알림을 읽음 처리합니다.
func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID string) (int, error) {
	if userID == "" {
		return 0, apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}
	return s.notiRepo.MarkAllAsRead(ctx, userID)
}

// GetUnreadCount는 읽지 않은 알림 수를 반환합니다.
func (s *NotificationService) GetUnreadCount(ctx context.Context, userID string) (int, map[string]int, error) {
	if userID == "" {
		return 0, nil, apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}
	return s.notiRepo.GetUnreadCount(ctx, userID)
}

// UpdatePreferences는 알림 설정을 업데이트합니다.
func (s *NotificationService) UpdatePreferences(ctx context.Context, pref *NotificationPreferences) (*NotificationPreferences, error) {
	if pref == nil || pref.UserID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}
	if err := s.prefRepo.Save(ctx, pref); err != nil {
		return nil, fmt.Errorf("알림 설정 저장 실패: %w", err)
	}
	return pref, nil
}

// GetPreferences는 알림 설정을 조회합니다.
func (s *NotificationService) GetPreferences(ctx context.Context, userID string) (*NotificationPreferences, error) {
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}

	pref, err := s.prefRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 설정이 없으면 기본값 반환
	if pref == nil {
		pref = &NotificationPreferences{
			UserID:             userID,
			PushEnabled:        true,
			EmailEnabled:       true,
			SMSEnabled:         false,
			InAppEnabled:       true,
			HealthAlertEnabled: true,
			CoachingEnabled:    true,
			PromotionEnabled:   false,
			Language:           "ko",
		}
	}

	return pref, nil
}

// NotificationTemplate은 재사용 가능한 알림 템플릿을 정의합니다.
type NotificationTemplate struct {
	Key      string
	Title    string
	BodyFmt  string // fmt.Sprintf 포맷 문자열
	Type     string // NotificationType 이름 ("prescription", "appointment", ...)
	Priority string // "low", "normal", "high", "urgent"
	Channel  string // "push", "email", "sms", "in_app"
}

// PredefinedTemplates는 모든 사전 정의된 알림 템플릿을 포함합니다.
var PredefinedTemplates = map[string]NotificationTemplate{
	// Prescription
	"prescription_created":   {Key: "prescription_created", Title: "새 처방전 발행", BodyFmt: "%s 의사가 처방전을 발행했습니다", Type: "prescription", Priority: "high", Channel: "push"},
	"prescription_sent":      {Key: "prescription_sent", Title: "처방전 전송 완료", BodyFmt: "%s 약국에 처방전이 전송되었습니다", Type: "prescription", Priority: "normal", Channel: "push"},
	"prescription_ready":     {Key: "prescription_ready", Title: "약 조제 완료", BodyFmt: "%s 약국에서 약 조제가 완료되었습니다. 수령해 주세요", Type: "prescription", Priority: "high", Channel: "push"},
	"prescription_dispensed": {Key: "prescription_dispensed", Title: "약 수령 완료", BodyFmt: "처방전 수령이 확인되었습니다", Type: "prescription", Priority: "normal", Channel: "in_app"},
	// Delivery
	"delivery_started": {Key: "delivery_started", Title: "배송 출발", BodyFmt: "처방약이 배송을 시작했습니다. 예상 도착: %s", Type: "prescription", Priority: "normal", Channel: "push"},
	"delivery_arrived": {Key: "delivery_arrived", Title: "배송 완료", BodyFmt: "처방약이 도착했습니다", Type: "prescription", Priority: "high", Channel: "push"},
	// Appointment
	"appointment_reminder":  {Key: "appointment_reminder", Title: "진료 예약 알림", BodyFmt: "%s에 %s 예약이 있습니다", Type: "appointment", Priority: "high", Channel: "push"},
	"appointment_cancelled": {Key: "appointment_cancelled", Title: "예약 취소", BodyFmt: "%s 예약이 취소되었습니다", Type: "appointment", Priority: "normal", Channel: "push"},
	// Health alert
	"health_alert_critical": {Key: "health_alert_critical", Title: "건강 이상 감지", BodyFmt: "%s 수치가 위험 범위입니다: %s", Type: "health_alert", Priority: "urgent", Channel: "push"},
	"health_alert_warning":  {Key: "health_alert_warning", Title: "건강 주의", BodyFmt: "%s 수치가 주의 범위입니다: %s", Type: "health_alert", Priority: "high", Channel: "push"},
	// Measurement
	"measurement_complete": {Key: "measurement_complete", Title: "측정 완료", BodyFmt: "%s 측정이 완료되었습니다. 결과를 확인하세요", Type: "measurement", Priority: "normal", Channel: "in_app"},
	// Family
	"family_data_shared": {Key: "family_data_shared", Title: "가족 데이터 공유", BodyFmt: "%s님이 건강 데이터를 공유했습니다", Type: "system", Priority: "normal", Channel: "in_app"},
}

// parseNotificationType은 문자열에서 NotificationType으로 변환합니다.
func parseNotificationType(s string) NotificationType {
	switch s {
	case "measurement":
		return TypeMeasurement
	case "health_alert":
		return TypeHealthAlert
	case "appointment":
		return TypeAppointment
	case "prescription":
		return TypePrescription
	case "community":
		return TypeCommunity
	case "system":
		return TypeSystem
	case "promotion":
		return TypePromotion
	default:
		return TypeUnknown
	}
}

// parseNotificationPriority는 문자열에서 NotificationPriority로 변환합니다.
func parseNotificationPriority(s string) NotificationPriority {
	switch s {
	case "low":
		return PriorityLow
	case "normal":
		return PriorityNormal
	case "high":
		return PriorityHigh
	case "urgent":
		return PriorityUrgent
	default:
		return PriorityNormal
	}
}

// parseNotificationChannel은 문자열에서 NotificationChannel로 변환합니다.
func parseNotificationChannel(s string) NotificationChannel {
	switch s {
	case "push":
		return ChannelPush
	case "email":
		return ChannelEmail
	case "sms":
		return ChannelSMS
	case "in_app":
		return ChannelInApp
	default:
		return ChannelUnknown
	}
}

// SendFromTemplate은 사전 정의된 템플릿을 사용하여 알림을 발송합니다.
func (s *NotificationService) SendFromTemplate(ctx context.Context, userID, templateKey string, args ...interface{}) error {
	tmpl, exists := PredefinedTemplates[templateKey]
	if !exists {
		return fmt.Errorf("알림 템플릿 없음: %s", templateKey)
	}
	body := fmt.Sprintf(tmpl.BodyFmt, args...)
	nType := parseNotificationType(tmpl.Type)
	priority := parseNotificationPriority(tmpl.Priority)
	channel := parseNotificationChannel(tmpl.Channel)
	_, err := s.SendNotification(ctx, userID, nType, channel, priority, tmpl.Title, body, "")
	return err
}

// NotificationTypeToString는 알림 타입을 문자열로 변환합니다.
func NotificationTypeToString(t NotificationType) string {
	switch t {
	case TypeMeasurement:
		return "measurement"
	case TypeHealthAlert:
		return "health_alert"
	case TypeAppointment:
		return "appointment"
	case TypePrescription:
		return "prescription"
	case TypeCommunity:
		return "community"
	case TypeSystem:
		return "system"
	case TypePromotion:
		return "promotion"
	default:
		return "unknown"
	}
}
