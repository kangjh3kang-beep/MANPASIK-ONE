// Package handler는 video-service의 gRPC 핸들러입니다.
package handler

import (
	"context"

	"github.com/manpasik/backend/services/video-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// VideoHandler는 VideoService gRPC 서버를 구현합니다.
type VideoHandler struct {
	v1.UnimplementedVideoServiceServer
	svc *service.VideoService
	log *zap.Logger
}

// NewVideoHandler는 VideoHandler를 생성합니다.
func NewVideoHandler(svc *service.VideoService, log *zap.Logger) *VideoHandler {
	return &VideoHandler{svc: svc, log: log}
}

// CreateRoom은 회의실 생성 RPC입니다.
func (h *VideoHandler) CreateRoom(ctx context.Context, req *v1.CreateRoomRequest) (*v1.Room, error) {
	if req == nil || req.HostUserId == "" {
		return nil, status.Error(codes.InvalidArgument, "host_user_id는 필수입니다")
	}

	room, err := h.svc.CreateRoom(
		ctx,
		req.Title,
		protoRoomTypeToService(req.RoomType),
		req.HostUserId,
		int(req.MaxParticipants),
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	return roomToProto(room), nil
}

// GetRoom은 회의실 조회 RPC입니다.
func (h *VideoHandler) GetRoom(ctx context.Context, req *v1.GetRoomRequest) (*v1.Room, error) {
	if req == nil || req.RoomId == "" {
		return nil, status.Error(codes.InvalidArgument, "room_id는 필수입니다")
	}

	room, err := h.svc.GetRoom(ctx, req.RoomId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return roomToProto(room), nil
}

// JoinRoom은 참가자 입장 RPC입니다.
func (h *VideoHandler) JoinRoom(ctx context.Context, req *v1.JoinRoomRequest) (*v1.JoinRoomResponse, error) {
	if req == nil || req.RoomId == "" {
		return nil, status.Error(codes.InvalidArgument, "room_id는 필수입니다")
	}
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	room, token, _, err := h.svc.JoinRoom(ctx, req.RoomId, req.UserId, "", "")
	if err != nil {
		return nil, toGRPC(err)
	}

	var participants []*v1.Participant
	for _, p := range room.Participants {
		if p.LeftAt.IsZero() {
			participants = append(participants, participantToProto(p))
		}
	}

	return &v1.JoinRoomResponse{
		Success:      true,
		Token:        token,
		Room:         roomToProto(room),
		Participants: participants,
	}, nil
}

// LeaveRoom은 참가자 퇴장 RPC입니다.
func (h *VideoHandler) LeaveRoom(ctx context.Context, req *v1.LeaveRoomRequest) (*v1.LeaveRoomResponse, error) {
	if req == nil || req.RoomId == "" {
		return nil, status.Error(codes.InvalidArgument, "room_id는 필수입니다")
	}

	_, err := h.svc.LeaveRoom(ctx, req.RoomId, req.UserId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.LeaveRoomResponse{
		Success: true,
		Message: "참가자가 퇴장했습니다",
	}, nil
}

// EndRoom은 회의실 종료 RPC입니다.
func (h *VideoHandler) EndRoom(ctx context.Context, req *v1.EndRoomRequest) (*v1.Room, error) {
	if req == nil || req.RoomId == "" {
		return nil, status.Error(codes.InvalidArgument, "room_id는 필수입니다")
	}

	room, err := h.svc.EndRoom(ctx, req.RoomId, req.UserId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return roomToProto(room), nil
}

// SendSignal은 WebRTC 시그널 전송 RPC입니다.
func (h *VideoHandler) SendSignal(ctx context.Context, req *v1.SendSignalRequest) (*v1.SendSignalResponse, error) {
	if req == nil || req.RoomId == "" {
		return nil, status.Error(codes.InvalidArgument, "room_id는 필수입니다")
	}
	if req.FromUserId == "" {
		return nil, status.Error(codes.InvalidArgument, "from_user_id는 필수입니다")
	}

	_, err := h.svc.SendSignal(
		ctx,
		req.RoomId,
		req.FromUserId,
		req.ToUserId,
		protoSignalTypeToService(req.SignalType),
		req.Payload,
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.SendSignalResponse{
		Success: true,
	}, nil
}

// ListParticipants는 참가자 목록 조회 RPC입니다.
func (h *VideoHandler) ListParticipants(ctx context.Context, req *v1.ListParticipantsRequest) (*v1.ListParticipantsResponse, error) {
	if req == nil || req.RoomId == "" {
		return nil, status.Error(codes.InvalidArgument, "room_id는 필수입니다")
	}

	participants, err := h.svc.ListParticipants(ctx, req.RoomId)
	if err != nil {
		return nil, toGRPC(err)
	}

	var pb []*v1.Participant
	for _, p := range participants {
		pb = append(pb, participantToProto(p))
	}

	return &v1.ListParticipantsResponse{
		Participants: pb,
	}, nil
}

// GetRoomStats는 회의실 통계 조회 RPC입니다.
func (h *VideoHandler) GetRoomStats(ctx context.Context, req *v1.GetRoomStatsRequest) (*v1.GetRoomStatsResponse, error) {
	if req == nil || req.RoomId == "" {
		return nil, status.Error(codes.InvalidArgument, "room_id는 필수입니다")
	}

	room, err := h.svc.GetRoomStats(ctx, req.RoomId)
	if err != nil {
		return nil, toGRPC(err)
	}

	currentActive := int32(0)
	for _, p := range room.Participants {
		if p.LeftAt.IsZero() {
			currentActive++
		}
	}

	resp := &v1.GetRoomStatsResponse{
		RoomId:              room.ID,
		TotalParticipants:   int32(len(room.Participants)),
		CurrentParticipants: currentActive,
		DurationSeconds:     int64(room.DurationSeconds),
		SignalCount:         int32(room.SignalCount),
	}
	if !room.StartedAt.IsZero() {
		resp.StartedAt = timestamppb.New(room.StartedAt)
	}

	return resp, nil
}

// ============================================================================
// 변환 헬퍼
// ============================================================================

func roomToProto(r *service.Room) *v1.Room {
	currentActive := int32(0)
	for _, p := range r.Participants {
		if p.LeftAt.IsZero() {
			currentActive++
		}
	}

	pb := &v1.Room{
		RoomId:           r.ID,
		HostUserId:       r.CreatedBy,
		RoomType:         serviceRoomTypeToProto(r.RoomType),
		Status:           serviceRoomStatusToProto(r.Status),
		Title:            r.Name,
		ParticipantCount: currentActive,
		MaxParticipants:  int32(r.MaxParticipants),
		CreatedAt:        timestamppb.New(r.CreatedAt),
	}
	if !r.StartedAt.IsZero() {
		pb.StartedAt = timestamppb.New(r.StartedAt)
	}
	if !r.EndedAt.IsZero() {
		pb.EndedAt = timestamppb.New(r.EndedAt)
	}
	return pb
}

func participantToProto(p *service.Participant) *v1.Participant {
	return &v1.Participant{
		UserId:      p.UserID,
		DisplayName: p.DisplayName,
		IsHost:      p.Role == "host",
		IsMuted:     !p.IsAudioEnabled,
		IsVideoOn:   p.IsVideoEnabled,
		JoinedAt:    timestamppb.New(p.JoinedAt),
	}
}

func protoRoomTypeToService(t v1.RoomType) service.RoomType {
	switch t {
	case v1.RoomType_ROOM_TYPE_ONE_TO_ONE:
		return service.RoomTypeOneToOne
	case v1.RoomType_ROOM_TYPE_GROUP:
		return service.RoomTypeGroup
	case v1.RoomType_ROOM_TYPE_WEBINAR:
		return service.RoomTypeWebinar
	case v1.RoomType_ROOM_TYPE_CONSULTATION:
		return service.RoomTypeConsultation
	default:
		return service.RoomTypeUnknown
	}
}

func serviceRoomTypeToProto(t service.RoomType) v1.RoomType {
	switch t {
	case service.RoomTypeOneToOne:
		return v1.RoomType_ROOM_TYPE_ONE_TO_ONE
	case service.RoomTypeGroup:
		return v1.RoomType_ROOM_TYPE_GROUP
	case service.RoomTypeWebinar:
		return v1.RoomType_ROOM_TYPE_WEBINAR
	case service.RoomTypeConsultation:
		return v1.RoomType_ROOM_TYPE_CONSULTATION
	default:
		return v1.RoomType_ROOM_TYPE_UNKNOWN
	}
}

func serviceRoomStatusToProto(s service.RoomStatus) v1.RoomStatus {
	switch s {
	case service.RoomStatusWaiting:
		return v1.RoomStatus_ROOM_STATUS_WAITING
	case service.RoomStatusActive:
		return v1.RoomStatus_ROOM_STATUS_ACTIVE
	case service.RoomStatusEnded:
		return v1.RoomStatus_ROOM_STATUS_ENDED
	case service.RoomStatusFailed:
		return v1.RoomStatus_ROOM_STATUS_FAILED
	default:
		return v1.RoomStatus_ROOM_STATUS_UNKNOWN
	}
}

func protoSignalTypeToService(t v1.SignalType) service.SignalType {
	switch t {
	case v1.SignalType_SIGNAL_TYPE_OFFER:
		return service.SignalTypeOffer
	case v1.SignalType_SIGNAL_TYPE_ANSWER:
		return service.SignalTypeAnswer
	case v1.SignalType_SIGNAL_TYPE_ICE_CANDIDATE:
		return service.SignalTypeICECandidate
	case v1.SignalType_SIGNAL_TYPE_RENEGOTIATE:
		return service.SignalTypeRenegotiate
	case v1.SignalType_SIGNAL_TYPE_MUTE:
		return service.SignalTypeMute
	case v1.SignalType_SIGNAL_TYPE_UNMUTE:
		return service.SignalTypeUnmute
	default:
		return service.SignalTypeUnknown
	}
}

func toGRPC(err error) error {
	if appErr, ok := err.(*apperrors.AppError); ok {
		return appErr.ToGRPC()
	}
	return status.Error(codes.Internal, err.Error())
}
