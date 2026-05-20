package handler

import (
	"context"

	"github.com/manpasik/backend/services/assistant-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AssistantHandler struct {
	v1.UnimplementedAssistantServiceServer
	svc *service.AssistantService
	log *zap.Logger
}

func NewAssistantHandler(svc *service.AssistantService, log *zap.Logger) *AssistantHandler {
	return &AssistantHandler{svc: svc, log: log}
}

func (h *AssistantHandler) SendCommand(ctx context.Context, req *v1.AssistantCommandRequest) (*v1.AssistantCommandResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}
	userTurn, asstTurn, session, err := h.svc.SendCommand(ctx, req.UserId, req.SessionId, req.Text)
	if err != nil {
		return nil, toGRPC(err)
	}
	_ = userTurn // 사용자 턴은 별도 저장됨
	return &v1.AssistantCommandResponse{
		SessionId:    session.ID,
		TurnId:       asstTurn.ID,
		ResponseText: asstTurn.Content,
		Intent:       asstTurn.Intent,
		ActionType:   asstTurn.ActionType,
		ActionResult: asstTurn.ActionResult,
	}, nil
}

func (h *AssistantHandler) GetSession(ctx context.Context, req *v1.GetAssistantSessionRequest) (*v1.AssistantSession, error) {
	if req == nil || req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id는 필수입니다")
	}
	sess, err := h.svc.GetSession(ctx, req.SessionId)
	if err != nil {
		return nil, toGRPC(err)
	}
	return sessionToProto(sess), nil
}

func (h *AssistantHandler) ListSessions(ctx context.Context, req *v1.ListAssistantSessionsRequest) (*v1.ListAssistantSessionsResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}
	sessions, total, err := h.svc.ListSessions(ctx, req.UserId, req.Limit, req.Offset)
	if err != nil {
		return nil, toGRPC(err)
	}
	proto := make([]*v1.AssistantSession, 0, len(sessions))
	for _, s := range sessions {
		proto = append(proto, sessionToProto(s))
	}
	return &v1.ListAssistantSessionsResponse{Sessions: proto, Total: total}, nil
}

func (h *AssistantHandler) ListTurns(ctx context.Context, req *v1.ListAssistantTurnsRequest) (*v1.ListAssistantTurnsResponse, error) {
	if req == nil || req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id는 필수입니다")
	}
	turns, total, err := h.svc.ListTurns(ctx, req.SessionId, req.Limit, req.Offset)
	if err != nil {
		return nil, toGRPC(err)
	}
	proto := make([]*v1.AssistantTurn, 0, len(turns))
	for _, t := range turns {
		proto = append(proto, turnToProto(t))
	}
	return &v1.ListAssistantTurnsResponse{Turns: proto, Total: total}, nil
}

func (h *AssistantHandler) DeleteSession(ctx context.Context, req *v1.DeleteAssistantSessionRequest) (*v1.DeleteAssistantSessionResponse, error) {
	if req == nil || req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id는 필수입니다")
	}
	if err := h.svc.DeleteSession(ctx, req.SessionId); err != nil {
		return nil, toGRPC(err)
	}
	return &v1.DeleteAssistantSessionResponse{Success: true}, nil
}

func sessionToProto(s *service.AssistantSession) *v1.AssistantSession {
	return &v1.AssistantSession{
		Id:     s.ID,
		UserId: s.UserID,
		Title:     s.Title,
		Status:    s.Status,
		TurnCount: s.TurnCount,
		CreatedAt: timestamppb.New(s.CreatedAt),
		UpdatedAt: timestamppb.New(s.UpdatedAt),
	}
}

func turnToProto(t *service.AssistantTurn) *v1.AssistantTurn {
	return &v1.AssistantTurn{
		Id:           t.ID,
		SessionId:    t.SessionID,
		Role:         t.Role,
		Content:      t.Content,
		Intent:       t.Intent,
		ActionType:   t.ActionType,
		ActionResult: t.ActionResult,
		CreatedAt:    timestamppb.New(t.CreatedAt),
	}
}

func toGRPC(err error) error {
	if err == nil {
		return nil
	}
	if ae, ok := err.(*apperrors.AppError); ok {
		return ae.ToGRPC()
	}
	return status.Error(codes.Internal, "내부 오류가 발생했습니다")
}
