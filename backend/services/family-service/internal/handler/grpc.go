// Package handler는 family-service의 gRPC 핸들러입니다.
package handler

import (
	"context"

	"github.com/manpasik/backend/services/family-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// FamilyHandler는 FamilyService gRPC 서버를 구현합니다.
type FamilyHandler struct {
	v1.UnimplementedFamilyServiceServer
	svc *service.FamilyService
	log *zap.Logger
}

// NewFamilyHandler는 FamilyHandler를 생성합니다.
func NewFamilyHandler(svc *service.FamilyService, log *zap.Logger) *FamilyHandler {
	return &FamilyHandler{svc: svc, log: log}
}

// CreateFamilyGroup은 가족 그룹 생성 RPC입니다.
func (h *FamilyHandler) CreateFamilyGroup(ctx context.Context, req *v1.CreateFamilyGroupRequest) (*v1.FamilyGroup, error) {
	if req == nil || req.OwnerUserId == "" {
		return nil, status.Error(codes.InvalidArgument, "owner_user_id는 필수입니다")
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name은 필수입니다")
	}

	group, err := h.svc.CreateFamilyGroup(ctx, req.OwnerUserId, req.Name, "")
	if err != nil {
		return nil, toGRPC(err)
	}

	return familyGroupToProto(group, 1), nil
}

// GetFamilyGroup은 가족 그룹 조회 RPC입니다.
func (h *FamilyHandler) GetFamilyGroup(ctx context.Context, req *v1.GetFamilyGroupRequest) (*v1.FamilyGroup, error) {
	if req == nil || req.GroupId == "" {
		return nil, status.Error(codes.InvalidArgument, "group_id는 필수입니다")
	}

	group, memberCount, err := h.svc.GetFamilyGroup(ctx, req.GroupId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return familyGroupToProto(group, memberCount), nil
}

// InviteMember는 멤버 초대 RPC입니다.
func (h *FamilyHandler) InviteMember(ctx context.Context, req *v1.InviteMemberRequest) (*v1.FamilyInvitation, error) {
	if req == nil || req.GroupId == "" || req.InviterUserId == "" || req.InviteeEmail == "" {
		return nil, status.Error(codes.InvalidArgument, "group_id, inviter_user_id, invitee_email은 필수입니다")
	}

	inv, err := h.svc.InviteMember(ctx, req.GroupId, req.InviterUserId, req.InviteeEmail, protoRoleToService(req.Role), "")
	if err != nil {
		return nil, toGRPC(err)
	}

	return invitationToProto(inv), nil
}

// RespondToInvitation은 초대 응답 RPC입니다.
func (h *FamilyHandler) RespondToInvitation(ctx context.Context, req *v1.RespondToInvitationRequest) (*v1.RespondToInvitationResponse, error) {
	if req == nil || req.InvitationId == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "invitation_id, user_id는 필수입니다")
	}

	groupID, _, err := h.svc.RespondToInvitation(ctx, req.InvitationId, req.UserId, req.Accept)
	if err != nil {
		return nil, toGRPC(err)
	}

	// Fetch the updated group to include in the response.
	group, memberCount, _ := h.svc.GetFamilyGroup(ctx, groupID)
	var pbGroup *v1.FamilyGroup
	if group != nil {
		pbGroup = familyGroupToProto(group, memberCount)
	}

	return &v1.RespondToInvitationResponse{
		Success: true,
		Message: "초대 응답이 처리되었습니다",
		Group:   pbGroup,
	}, nil
}

// RemoveMember는 멤버 제거 RPC입니다.
func (h *FamilyHandler) RemoveMember(ctx context.Context, req *v1.RemoveMemberRequest) (*v1.RemoveMemberResponse, error) {
	if req == nil || req.GroupId == "" || req.RequesterUserId == "" || req.TargetUserId == "" {
		return nil, status.Error(codes.InvalidArgument, "group_id, requester_user_id, target_user_id는 필수입니다")
	}

	err := h.svc.RemoveMember(ctx, req.GroupId, req.TargetUserId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.RemoveMemberResponse{Success: true}, nil
}

// UpdateMemberRole은 멤버 역할 변경 RPC입니다.
func (h *FamilyHandler) UpdateMemberRole(ctx context.Context, req *v1.UpdateMemberRoleRequest) (*v1.FamilyMember, error) {
	if req == nil || req.GroupId == "" || req.RequesterUserId == "" || req.TargetUserId == "" {
		return nil, status.Error(codes.InvalidArgument, "group_id, requester_user_id, target_user_id는 필수입니다")
	}

	member, err := h.svc.UpdateMemberRole(ctx, req.GroupId, req.TargetUserId, protoRoleToService(req.NewRole))
	if err != nil {
		return nil, toGRPC(err)
	}

	return memberToProto(member), nil
}

// ListFamilyMembers는 멤버 목록 조회 RPC입니다.
func (h *FamilyHandler) ListFamilyMembers(ctx context.Context, req *v1.ListFamilyMembersRequest) (*v1.ListFamilyMembersResponse, error) {
	if req == nil || req.GroupId == "" {
		return nil, status.Error(codes.InvalidArgument, "group_id는 필수입니다")
	}

	members, err := h.svc.ListFamilyMembers(ctx, req.GroupId)
	if err != nil {
		return nil, toGRPC(err)
	}

	var pbMembers []*v1.FamilyMember
	for _, m := range members {
		pbMembers = append(pbMembers, memberToProto(m))
	}

	return &v1.ListFamilyMembersResponse{Members: pbMembers}, nil
}

// SetSharingPreferences는 공유 설정 변경 RPC입니다.
func (h *FamilyHandler) SetSharingPreferences(ctx context.Context, req *v1.SetSharingPreferencesRequest) (*v1.SharingPreferences, error) {
	if req == nil || req.UserId == "" || req.GroupId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id, group_id는 필수입니다")
	}

	pref := &service.SharingPreferences{
		UserID:               req.UserId,
		GroupID:              req.GroupId,
		ShareMeasurements:    req.ShareMeasurements,
		ShareHealthScore:     req.ShareHealthScore,
		ShareGoals:           req.ShareGoals,
		ShareCoaching:        req.ShareCoaching,
		ShareAlerts:          req.ShareAlerts,
		AllowedViewerIDs:     req.SharedWithUserIds,
		MeasurementDaysLimit: int(req.MeasurementDaysLimit),
		AllowedBiomarkers:    req.AllowedBiomarkers,
		RequireApproval:      req.RequireApproval,
	}

	result, err := h.svc.SetSharingPreferences(ctx, pref)
	if err != nil {
		return nil, toGRPC(err)
	}

	return sharingToProto(result), nil
}

// GetSharedHealthData는 공유 건강 데이터 조회 RPC입니다.
func (h *FamilyHandler) GetSharedHealthData(ctx context.Context, req *v1.GetSharedHealthDataRequest) (*v1.GetSharedHealthDataResponse, error) {
	if req == nil || req.RequesterUserId == "" || req.GroupId == "" {
		return nil, status.Error(codes.InvalidArgument, "requester_user_id, group_id는 필수입니다")
	}

	summaries, err := h.svc.GetSharedHealthData(ctx, req.RequesterUserId, req.TargetUserId, req.GroupId, 0)
	if err != nil {
		return nil, toGRPC(err)
	}

	resp := &v1.GetSharedHealthDataResponse{}
	if len(summaries) > 0 {
		resp.TargetUserId = summaries[0].UserID
		resp.TargetDisplayName = summaries[0].DisplayName
		resp.HealthScore = summaries[0].HealthScore
		if summaries[0].LastMeasurementAt != nil {
			resp.LastActive = summaries[0].LastMeasurementAt.Format("2006-01-02T15:04:05Z")
		}
	}

	return resp, nil
}

// ValidateSharingAccess는 공유 접근 검증 RPC입니다.
func (h *FamilyHandler) ValidateSharingAccess(ctx context.Context, req *v1.ValidateSharingAccessRequest) (*v1.ValidateSharingAccessResponse, error) {
	if req == nil || req.GroupId == "" || req.RequesterUserId == "" || req.TargetUserId == "" {
		return nil, status.Error(codes.InvalidArgument, "group_id, requester_user_id, target_user_id는 필수입니다")
	}

	allowed, reason, err := h.svc.ValidateSharingAccess(ctx, req.GroupId, req.RequesterUserId, req.TargetUserId, req.Biomarker)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.ValidateSharingAccessResponse{
		Allowed:    allowed,
		DenyReason: reason,
	}, nil
}

// ============================================================================
// 변환 헬퍼
// ============================================================================

func familyGroupToProto(g *service.FamilyGroup, memberCount int) *v1.FamilyGroup {
	return &v1.FamilyGroup{
		GroupId:     g.ID,
		Name:        g.GroupName,
		OwnerUserId: g.OwnerUserID,
		MaxMembers:  int32(g.MaxMembers),
		CreatedAt:   timestamppb.New(g.CreatedAt),
	}
}

func invitationToProto(inv *service.FamilyInvitation) *v1.FamilyInvitation {
	return &v1.FamilyInvitation{
		InvitationId:  inv.ID,
		GroupId:       inv.GroupID,
		InviterUserId: inv.InviterUserID,
		InviteeEmail:  inv.InviteeEmail,
		Role:          serviceRoleToProto(inv.Role),
		Status:        serviceInviteStatusToProto(inv.Status),
		CreatedAt:     timestamppb.New(inv.CreatedAt),
		ExpiresAt:     timestamppb.New(inv.ExpiresAt),
	}
}

func memberToProto(m *service.FamilyMember) *v1.FamilyMember {
	return &v1.FamilyMember{
		UserId:      m.UserID,
		DisplayName: m.DisplayName,
		Role:        serviceRoleToProto(m.Role),
		JoinedAt:    timestamppb.New(m.JoinedAt),
	}
}

func sharingToProto(p *service.SharingPreferences) *v1.SharingPreferences {
	return &v1.SharingPreferences{
		UserId:               p.UserID,
		GroupId:              p.GroupID,
		ShareMeasurements:    p.ShareMeasurements,
		ShareHealthScore:     p.ShareHealthScore,
		ShareGoals:           p.ShareGoals,
		ShareCoaching:        p.ShareCoaching,
		ShareAlerts:          p.ShareAlerts,
		SharedWithUserIds:    p.AllowedViewerIDs,
		MeasurementDaysLimit: int32(p.MeasurementDaysLimit),
		AllowedBiomarkers:    p.AllowedBiomarkers,
		RequireApproval:      p.RequireApproval,
	}
}

func protoRoleToService(r v1.FamilyRole) service.FamilyRole {
	switch r {
	case v1.FamilyRole_FAMILY_ROLE_OWNER:
		return service.RoleOwner
	case v1.FamilyRole_FAMILY_ROLE_GUARDIAN:
		return service.RoleGuardian
	case v1.FamilyRole_FAMILY_ROLE_MEMBER:
		return service.RoleMember
	case v1.FamilyRole_FAMILY_ROLE_CHILD:
		return service.RoleChild
	case v1.FamilyRole_FAMILY_ROLE_ELDERLY:
		return service.RoleElderly
	default:
		return service.RoleUnknown
	}
}

func serviceRoleToProto(r service.FamilyRole) v1.FamilyRole {
	switch r {
	case service.RoleOwner:
		return v1.FamilyRole_FAMILY_ROLE_OWNER
	case service.RoleGuardian:
		return v1.FamilyRole_FAMILY_ROLE_GUARDIAN
	case service.RoleMember:
		return v1.FamilyRole_FAMILY_ROLE_MEMBER
	case service.RoleChild:
		return v1.FamilyRole_FAMILY_ROLE_CHILD
	case service.RoleElderly:
		return v1.FamilyRole_FAMILY_ROLE_ELDERLY
	default:
		return v1.FamilyRole_FAMILY_ROLE_UNKNOWN
	}
}

func serviceInviteStatusToProto(s service.InvitationStatus) v1.InvitationStatus {
	switch s {
	case service.InvitePending:
		return v1.InvitationStatus_INVITATION_STATUS_PENDING
	case service.InviteAccepted:
		return v1.InvitationStatus_INVITATION_STATUS_ACCEPTED
	case service.InviteDeclined:
		return v1.InvitationStatus_INVITATION_STATUS_REJECTED
	case service.InviteExpired:
		return v1.InvitationStatus_INVITATION_STATUS_EXPIRED
	default:
		return v1.InvitationStatus_INVITATION_STATUS_UNKNOWN
	}
}

func toGRPC(err error) error {
	if appErr, ok := err.(*apperrors.AppError); ok {
		return appErr.ToGRPC()
	}
	return status.Error(codes.Internal, err.Error())
}
