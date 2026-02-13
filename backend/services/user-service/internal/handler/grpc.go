// Package handler는 user-service의 gRPC 핸들러입니다.
package handler

import (
	"context"

	"github.com/manpasik/backend/services/user-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserHandler는 UserService gRPC 서버를 구현합니다.
type UserHandler struct {
	v1.UnimplementedUserServiceServer
	svc *service.UserService
	log *zap.Logger
}

// NewUserHandler는 UserHandler를 생성합니다.
func NewUserHandler(svc *service.UserService, log *zap.Logger) *UserHandler {
	return &UserHandler{svc: svc, log: log}
}

// GetProfile은 사용자 프로필 조회 RPC입니다.
func (h *UserHandler) GetProfile(ctx context.Context, req *v1.GetProfileRequest) (*v1.UserProfile, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	profile, err := h.svc.GetProfile(ctx, req.UserId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.UserProfile{
		UserId:      profile.UserID,
		Email:       profile.Email,
		DisplayName: profile.DisplayName,
		AvatarUrl:   profile.AvatarURL,
		Language:    profile.Language,
		Timezone:    profile.Timezone,
		CreatedAt:   timestamppb.New(profile.CreatedAt),
	}, nil
}

// UpdateProfile은 사용자 프로필 업데이트 RPC입니다.
func (h *UserHandler) UpdateProfile(ctx context.Context, req *v1.UpdateProfileRequest) (*v1.UserProfile, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	profile, err := h.svc.UpdateProfile(ctx, req.UserId, req.DisplayName, req.AvatarUrl, req.Language, req.Timezone)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.UserProfile{
		UserId:      profile.UserID,
		Email:       profile.Email,
		DisplayName: profile.DisplayName,
		AvatarUrl:   profile.AvatarURL,
		Language:    profile.Language,
		Timezone:    profile.Timezone,
		CreatedAt:   timestamppb.New(profile.CreatedAt),
	}, nil
}

// GetSubscription은 구독 정보 조회 RPC입니다.
func (h *UserHandler) GetSubscription(ctx context.Context, req *v1.GetSubscriptionRequest) (*v1.SubscriptionInfo, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	sub, err := h.svc.GetSubscription(ctx, req.UserId)
	if err != nil {
		return nil, toGRPC(err)
	}

	tierMap := map[service.SubscriptionTier]v1.SubscriptionTier{
		service.TierFree:     v1.SubscriptionTier_SUBSCRIPTION_TIER_FREE,
		service.TierBasic:    v1.SubscriptionTier_SUBSCRIPTION_TIER_BASIC,
		service.TierPro:      v1.SubscriptionTier_SUBSCRIPTION_TIER_PRO,
		service.TierClinical: v1.SubscriptionTier_SUBSCRIPTION_TIER_CLINICAL,
	}

	info := &v1.SubscriptionInfo{
		UserId:              sub.UserID,
		Tier:                tierMap[sub.Tier],
		MaxDevices:          int32(sub.MaxDevices),
		MaxFamilyMembers:    int32(sub.MaxFamilyMembers),
		AiCoachingEnabled:   sub.AICoachingEnabled,
		TelemedicineEnabled: sub.TelemedicineEnabled,
	}

	if !sub.StartedAt.IsZero() {
		info.StartedAt = timestamppb.New(sub.StartedAt)
	}
	if !sub.ExpiresAt.IsZero() {
		info.ExpiresAt = timestamppb.New(sub.ExpiresAt)
	}

	return info, nil
}

// toGRPC는 AppError를 gRPC status로 변환합니다.
func toGRPC(err error) error {
	if err == nil {
		return nil
	}
	if ae, ok := err.(*apperrors.AppError); ok {
		return ae.ToGRPC()
	}
	if s, ok := status.FromError(err); ok {
		return s.Err()
	}
	return status.Error(codes.Internal, "내부 오류가 발생했습니다")
}
