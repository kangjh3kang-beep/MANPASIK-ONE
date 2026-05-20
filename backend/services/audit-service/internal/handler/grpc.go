package handler

import (
	"context"

	"github.com/manpasik/backend/services/audit-service/internal/service"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AuditHandler는 감사 로그 gRPC 핸들러입니다.
// AdminService의 GetAuditLog/GetAuditLogDetails RPC를 위임받아 처리합니다.
type AuditHandler struct {
	v1.UnimplementedAdminServiceServer
	svc *service.AuditService
	log *zap.Logger
}

// NewAuditHandler는 새 AuditHandler를 생성합니다.
func NewAuditHandler(svc *service.AuditService, log *zap.Logger) *AuditHandler {
	return &AuditHandler{svc: svc, log: log}
}

// GetAuditLog는 감사 로그 조회 RPC입니다.
func (h *AuditHandler) GetAuditLog(ctx context.Context, req *v1.GetAuditLogRequest) (*v1.GetAuditLogResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "요청이 비어있습니다")
	}

	filter := service.AuditFilter{
		AdminID: req.AdminId,
		Limit:   int(req.Limit),
		Offset:  int(req.Offset),
	}
	if req.StartTime != nil {
		filter.StartTime = req.StartTime.AsTime()
	}
	if req.EndTime != nil {
		filter.EndTime = req.EndTime.AsTime()
	}

	entries, total, err := h.svc.ListEntries(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbEntries := make([]*v1.AuditLogEntry, 0, len(entries))
	for _, e := range entries {
		pbEntries = append(pbEntries, &v1.AuditLogEntry{
			EntryId:      e.ID,
			AdminId:      e.AdminID,
			ResourceType: e.ResourceType,
			ResourceId:   e.ResourceID,
			Details:      e.Description,
			IpAddress:    e.IPAddress,
			Timestamp:    timestamppb.New(e.Timestamp),
		})
	}

	return &v1.GetAuditLogResponse{
		Entries:    pbEntries,
		TotalCount: int32(total),
	}, nil
}

// GetAuditLogDetails는 상세 감사 로그 조회 RPC입니다.
func (h *AuditHandler) GetAuditLogDetails(ctx context.Context, req *v1.GetAuditLogDetailsRequest) (*v1.GetAuditLogDetailsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "요청이 비어있습니다")
	}

	filter := service.AuditFilter{
		AdminID: req.AdminId,
		Action:  req.Action,
		Limit:   int(req.Limit),
		Offset:  int(req.Offset),
	}

	entries, total, err := h.svc.ListEntries(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	details := make([]*v1.AuditLogDetail, 0, len(entries))
	for _, e := range entries {
		details = append(details, &v1.AuditLogDetail{
			Id:           e.ID,
			AdminId:      e.AdminID,
			Action:       e.Action,
			ResourceType: e.ResourceType,
			ResourceId:   e.ResourceID,
			Description:  e.Description,
			IpAddress:    e.IPAddress,
			CreatedAt:    timestamppb.New(e.Timestamp),
		})
	}

	return &v1.GetAuditLogDetailsResponse{
		Details: details,
		Total:   int32(total),
	}, nil
}
