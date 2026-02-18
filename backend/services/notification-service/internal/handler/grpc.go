// Package handler는 notification-service의 gRPC 핸들러입니다.
package handler

import (
	"context"
	"encoding/json"

	"github.com/manpasik/backend/services/notification-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NotificationHandler는 NotificationService gRPC 서버를 구현합니다.
type NotificationHandler struct {
	v1.UnimplementedNotificationServiceServer
	svc *service.NotificationService
	log *zap.Logger
}

// NewNotificationHandler는 NotificationHandler를 생성합니다.
func NewNotificationHandler(svc *service.NotificationService, log *zap.Logger) *NotificationHandler {
	return &NotificationHandler{svc: svc, log: log}
}

// SendNotification은 알림 발송 RPC입니다.
func (h *NotificationHandler) SendNotification(ctx context.Context, req *v1.SendNotificationRequest) (*v1.Notification, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}
	if req.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "title은 필수입니다")
	}

	// req.Data is map[string]string in proto; service expects a JSON string.
	var dataStr string
	if len(req.Data) > 0 {
		dataBytes, err := json.Marshal(req.Data)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "data 직렬화 실패: %v", err)
		}
		dataStr = string(dataBytes)
	}

	noti, err := h.svc.SendNotification(
		ctx,
		req.UserId,
		protoNotificationTypeToService(req.Type),
		protoChannelToService(req.Channel),
		protoPriorityToService(req.Priority),
		req.Title,
		req.Body,
		dataStr,
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	return notificationToProto(noti), nil
}

// ListNotifications는 알림 목록 조회 RPC입니다.
func (h *NotificationHandler) ListNotifications(ctx context.Context, req *v1.ListNotificationsRequest) (*v1.ListNotificationsResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	notis, total, _, err := h.svc.ListNotifications(
		ctx,
		req.UserId,
		service.TypeUnknown,
		req.UnreadOnly,
		int(req.Limit),
		int(req.Offset),
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	var pbNotis []*v1.Notification
	for _, n := range notis {
		pbNotis = append(pbNotis, notificationToProto(n))
	}

	return &v1.ListNotificationsResponse{
		Notifications: pbNotis,
		TotalCount:    int32(total),
	}, nil
}

// MarkAsRead는 알림 읽음 처리 RPC입니다.
func (h *NotificationHandler) MarkAsRead(ctx context.Context, req *v1.MarkAsReadRequest) (*v1.MarkAsReadResponse, error) {
	if req == nil || req.NotificationId == "" {
		return nil, status.Error(codes.InvalidArgument, "notification_id는 필수입니다")
	}

	err := h.svc.MarkAsRead(ctx, req.NotificationId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.MarkAsReadResponse{Success: true}, nil
}

// MarkAllAsRead는 모든 알림 읽음 처리 RPC입니다.
func (h *NotificationHandler) MarkAllAsRead(ctx context.Context, req *v1.MarkAllAsReadRequest) (*v1.MarkAllAsReadResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	count, err := h.svc.MarkAllAsRead(ctx, req.UserId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.MarkAllAsReadResponse{MarkedCount: int32(count)}, nil
}

// GetUnreadCount는 읽지 않은 알림 수 조회 RPC입니다.
func (h *NotificationHandler) GetUnreadCount(ctx context.Context, req *v1.GetUnreadCountRequest) (*v1.GetUnreadCountResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	count, _, err := h.svc.GetUnreadCount(ctx, req.UserId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.GetUnreadCountResponse{
		Count: int32(count),
	}, nil
}

// UpdateNotificationPreferences는 알림 설정 업데이트 RPC입니다.
func (h *NotificationHandler) UpdateNotificationPreferences(ctx context.Context, req *v1.UpdateNotificationPreferencesRequest) (*v1.NotificationPreferences, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}
	pref := &service.NotificationPreferences{
		UserID:             req.UserId,
		PushEnabled:        req.PushEnabled,
		EmailEnabled:       req.EmailEnabled,
		SMSEnabled:         req.SmsEnabled,
		HealthAlertEnabled: req.HealthAlerts,
		CoachingEnabled:    req.AppointmentReminders,
		PromotionEnabled:   req.Promotions,
		QuietHoursStart:    req.QuietHoursStart,
		QuietHoursEnd:      req.QuietHoursEnd,
	}

	result, err := h.svc.UpdatePreferences(ctx, pref)
	if err != nil {
		return nil, toGRPC(err)
	}

	return preferencesToProto(result), nil
}

// GetNotificationPreferences는 알림 설정 조회 RPC입니다.
func (h *NotificationHandler) GetNotificationPreferences(ctx context.Context, req *v1.GetNotificationPreferencesRequest) (*v1.NotificationPreferences, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	pref, err := h.svc.GetPreferences(ctx, req.UserId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return preferencesToProto(pref), nil
}

// SendFromTemplate은 사전 정의된 템플릿으로 알림을 발송하는 RPC입니다.
func (h *NotificationHandler) SendFromTemplate(ctx context.Context, req *v1.SendFromTemplateRequest) (*v1.Notification, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}
	if req.TemplateKey == "" {
		return nil, status.Error(codes.InvalidArgument, "template_key는 필수입니다")
	}

	// data map의 값들을 args로 변환
	var args []interface{}
	for _, v := range req.Data {
		args = append(args, v)
	}

	err := h.svc.SendFromTemplate(ctx, req.UserId, req.TemplateKey, args...)
	if err != nil {
		return nil, toGRPC(err)
	}

	// SendFromTemplate은 내부에서 SendNotification을 호출하므로
	// 가장 최근 알림을 조회하여 반환
	notis, _, _, err := h.svc.ListNotifications(ctx, req.UserId, service.TypeUnknown, false, 1, 0)
	if err != nil || len(notis) == 0 {
		return &v1.Notification{UserId: req.UserId}, nil
	}

	return notificationToProto(notis[0]), nil
}

// ============================================================================
// 변환 헬퍼
// ============================================================================

func notificationToProto(n *service.Notification) *v1.Notification {
	// n.Data is a JSON string; proto expects map[string]string.
	dataMap := make(map[string]string)
	if n.Data != "" {
		if err := json.Unmarshal([]byte(n.Data), &dataMap); err != nil {
			// Fallback: store the raw string under a "raw" key.
			dataMap["raw"] = n.Data
		}
	}

	pb := &v1.Notification{
		NotificationId: n.ID,
		UserId:         n.UserID,
		Type:           serviceNotificationTypeToProto(n.Type),
		Channel:        serviceChannelToProto(n.Channel),
		Priority:       servicePriorityToProto(n.Priority),
		Title:          n.Title,
		Body:           n.Body,
		Data:           dataMap,
		IsRead:         n.IsRead,
		CreatedAt:      timestamppb.New(n.CreatedAt),
	}
	if n.ReadAt != nil {
		pb.ReadAt = timestamppb.New(*n.ReadAt)
	}
	return pb
}

func preferencesToProto(p *service.NotificationPreferences) *v1.NotificationPreferences {
	return &v1.NotificationPreferences{
		UserId:               p.UserID,
		PushEnabled:          p.PushEnabled,
		EmailEnabled:         p.EmailEnabled,
		SmsEnabled:           p.SMSEnabled,
		HealthAlerts:         p.HealthAlertEnabled,
		AppointmentReminders: p.CoachingEnabled,
		Promotions:           p.PromotionEnabled,
		QuietHoursStart:      p.QuietHoursStart,
		QuietHoursEnd:        p.QuietHoursEnd,
	}
}

func protoNotificationTypeToService(t v1.NotificationType) service.NotificationType {
	switch t {
	case v1.NotificationType_NOTIFICATION_TYPE_MEASUREMENT:
		return service.TypeMeasurement
	case v1.NotificationType_NOTIFICATION_TYPE_HEALTH_ALERT:
		return service.TypeHealthAlert
	case v1.NotificationType_NOTIFICATION_TYPE_APPOINTMENT:
		return service.TypeAppointment
	case v1.NotificationType_NOTIFICATION_TYPE_PRESCRIPTION:
		return service.TypePrescription
	case v1.NotificationType_NOTIFICATION_TYPE_COMMUNITY:
		return service.TypeCommunity
	case v1.NotificationType_NOTIFICATION_TYPE_SYSTEM:
		return service.TypeSystem
	case v1.NotificationType_NOTIFICATION_TYPE_PROMOTION:
		return service.TypePromotion
	default:
		return service.TypeUnknown
	}
}

func serviceNotificationTypeToProto(t service.NotificationType) v1.NotificationType {
	switch t {
	case service.TypeMeasurement:
		return v1.NotificationType_NOTIFICATION_TYPE_MEASUREMENT
	case service.TypeHealthAlert:
		return v1.NotificationType_NOTIFICATION_TYPE_HEALTH_ALERT
	case service.TypeAppointment:
		return v1.NotificationType_NOTIFICATION_TYPE_APPOINTMENT
	case service.TypePrescription:
		return v1.NotificationType_NOTIFICATION_TYPE_PRESCRIPTION
	case service.TypeCommunity:
		return v1.NotificationType_NOTIFICATION_TYPE_COMMUNITY
	case service.TypeSystem:
		return v1.NotificationType_NOTIFICATION_TYPE_SYSTEM
	case service.TypePromotion:
		return v1.NotificationType_NOTIFICATION_TYPE_PROMOTION
	default:
		return v1.NotificationType_NOTIFICATION_TYPE_UNKNOWN
	}
}

func protoChannelToService(c v1.NotificationChannel) service.NotificationChannel {
	switch c {
	case v1.NotificationChannel_NOTIFICATION_CHANNEL_PUSH:
		return service.ChannelPush
	case v1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL:
		return service.ChannelEmail
	case v1.NotificationChannel_NOTIFICATION_CHANNEL_SMS:
		return service.ChannelSMS
	case v1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP:
		return service.ChannelInApp
	default:
		return service.ChannelUnknown
	}
}

func serviceChannelToProto(c service.NotificationChannel) v1.NotificationChannel {
	switch c {
	case service.ChannelPush:
		return v1.NotificationChannel_NOTIFICATION_CHANNEL_PUSH
	case service.ChannelEmail:
		return v1.NotificationChannel_NOTIFICATION_CHANNEL_EMAIL
	case service.ChannelSMS:
		return v1.NotificationChannel_NOTIFICATION_CHANNEL_SMS
	case service.ChannelInApp:
		return v1.NotificationChannel_NOTIFICATION_CHANNEL_IN_APP
	default:
		return v1.NotificationChannel_NOTIFICATION_CHANNEL_UNKNOWN
	}
}

func protoPriorityToService(p v1.NotificationPriority) service.NotificationPriority {
	switch p {
	case v1.NotificationPriority_NOTIFICATION_PRIORITY_LOW:
		return service.PriorityLow
	case v1.NotificationPriority_NOTIFICATION_PRIORITY_NORMAL:
		return service.PriorityNormal
	case v1.NotificationPriority_NOTIFICATION_PRIORITY_HIGH:
		return service.PriorityHigh
	case v1.NotificationPriority_NOTIFICATION_PRIORITY_URGENT:
		return service.PriorityUrgent
	default:
		return service.PriorityUnknown
	}
}

func servicePriorityToProto(p service.NotificationPriority) v1.NotificationPriority {
	switch p {
	case service.PriorityLow:
		return v1.NotificationPriority_NOTIFICATION_PRIORITY_LOW
	case service.PriorityNormal:
		return v1.NotificationPriority_NOTIFICATION_PRIORITY_NORMAL
	case service.PriorityHigh:
		return v1.NotificationPriority_NOTIFICATION_PRIORITY_HIGH
	case service.PriorityUrgent:
		return v1.NotificationPriority_NOTIFICATION_PRIORITY_URGENT
	default:
		return v1.NotificationPriority_NOTIFICATION_PRIORITY_UNKNOWN
	}
}

func toGRPC(err error) error {
	if appErr, ok := err.(*apperrors.AppError); ok {
		return appErr.ToGRPC()
	}
	return status.Error(codes.Internal, err.Error())
}
