// Package handler는 cartridge-service의 gRPC 핸들러입니다.
package handler

import (
	"context"

	"github.com/manpasik/backend/services/cartridge-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CartridgeHandler는 CartridgeService gRPC 서버를 구현합니다.
type CartridgeHandler struct {
	v1.UnimplementedCartridgeServiceServer
	svc *service.CartridgeService
	log *zap.Logger
}

// NewCartridgeHandler는 CartridgeHandler를 생성합니다.
func NewCartridgeHandler(svc *service.CartridgeService, log *zap.Logger) *CartridgeHandler {
	return &CartridgeHandler{svc: svc, log: log}
}

// ReadCartridge는 NFC 태그 읽기 RPC입니다.
func (h *CartridgeHandler) ReadCartridge(ctx context.Context, req *v1.ReadCartridgeRequest) (*v1.CartridgeDetail, error) {
	if req == nil || len(req.NfcTagData) == 0 {
		return nil, status.Error(codes.InvalidArgument, "nfc_tag_data는 필수입니다")
	}

	detail, err := h.svc.ReadCartridge(ctx, req.NfcTagData, req.TagVersion)
	if err != nil {
		return nil, toGRPC(err)
	}

	// ReadCartridge 후 상태 저장소에 초기화
	_ = h.svc.InitCartridgeState(ctx, detail)

	return cartridgeDetailToProto(detail), nil
}

// RecordUsage는 카트리지 사용 기록 RPC입니다.
func (h *CartridgeHandler) RecordUsage(ctx context.Context, req *v1.RecordUsageRequest) (*v1.RecordUsageResponse, error) {
	if req == nil || req.UserId == "" || req.CartridgeUid == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id와 cartridge_uid는 필수입니다")
	}

	remaining, err := h.svc.RecordUsage(ctx, req.UserId, req.SessionId, req.CartridgeUid, req.CategoryCode, req.TypeIndex)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.RecordUsageResponse{
		Success:        true,
		RemainingUses:  remaining,
		RemainingDaily: -1, // 무제한 (구독 연동 시 변경)
		RemainingMonthly: -1,
	}, nil
}

// GetUsageHistory는 사용 이력 조회 RPC입니다.
func (h *CartridgeHandler) GetUsageHistory(ctx context.Context, req *v1.GetUsageHistoryRequest) (*v1.GetUsageHistoryResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	records, totalCount, err := h.svc.GetUsageHistory(ctx, req.UserId, req.Limit, req.Offset)
	if err != nil {
		return nil, toGRPC(err)
	}

	protoRecords := make([]*v1.CartridgeUsageRecord, 0, len(records))
	for _, r := range records {
		protoRecords = append(protoRecords, usageRecordToProto(r))
	}

	return &v1.GetUsageHistoryResponse{
		Records:    protoRecords,
		TotalCount: totalCount,
	}, nil
}

// GetCartridgeType은 카트리지 타입 정보 조회 RPC입니다.
func (h *CartridgeHandler) GetCartridgeType(ctx context.Context, req *v1.GetCartridgeTypeRequest) (*v1.CartridgeTypeInfo, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "요청이 비어 있습니다")
	}

	typeInfo, err := h.svc.GetCartridgeType(ctx, req.CategoryCode, req.TypeIndex)
	if err != nil {
		return nil, toGRPC(err)
	}

	return typeInfoToProto(typeInfo), nil
}

// ListCategories는 카테고리 목록 조회 RPC입니다.
func (h *CartridgeHandler) ListCategories(_ context.Context, _ *v1.ListCategoriesRequest) (*v1.ListCategoriesResponse, error) {
	categories := h.svc.ListCategories()

	protoCategories := make([]*v1.CartridgeCategoryInfo, 0, len(categories))
	for _, cat := range categories {
		protoCategories = append(protoCategories, &v1.CartridgeCategoryInfo{
			Code:        cat.Code,
			NameEn:      cat.NameEN,
			NameKo:      cat.NameKO,
			Description: cat.Description,
			TypeCount:   cat.TypeCount,
			IsActive:    cat.IsActive,
		})
	}

	return &v1.ListCategoriesResponse{
		Categories: protoCategories,
	}, nil
}

// ListTypesByCategory는 카테고리별 타입 목록 조회 RPC입니다.
func (h *CartridgeHandler) ListTypesByCategory(ctx context.Context, req *v1.ListTypesByCategoryRequest) (*v1.ListTypesByCategoryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "요청이 비어 있습니다")
	}

	types, err := h.svc.ListTypesByCategory(ctx, req.CategoryCode)
	if err != nil {
		return nil, toGRPC(err)
	}

	protoTypes := make([]*v1.CartridgeTypeInfo, 0, len(types))
	for _, t := range types {
		protoTypes = append(protoTypes, typeInfoToProto(t))
	}

	return &v1.ListTypesByCategoryResponse{
		Types: protoTypes,
	}, nil
}

// GetRemainingUses는 잔여 사용 횟수 조회 RPC입니다.
func (h *CartridgeHandler) GetRemainingUses(ctx context.Context, req *v1.GetRemainingUsesRequest) (*v1.GetRemainingUsesResponse, error) {
	if req == nil || req.CartridgeUid == "" {
		return nil, status.Error(codes.InvalidArgument, "cartridge_uid는 필수입니다")
	}

	info, err := h.svc.GetRemainingUses(ctx, req.CartridgeUid)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.GetRemainingUsesResponse{
		CartridgeUid:  info.CartridgeUID,
		RemainingUses: info.RemainingUses,
		MaxUses:       info.MaxUses,
		ExpiryDate:    info.ExpiryDate,
		IsExpired:     info.IsExpired,
	}, nil
}

// ValidateCartridge는 카트리지 유효성 검증 RPC입니다.
func (h *CartridgeHandler) ValidateCartridge(ctx context.Context, req *v1.ValidateCartridgeRequest) (*v1.ValidateCartridgeResponse, error) {
	if req == nil || req.CartridgeUid == "" {
		return nil, status.Error(codes.InvalidArgument, "cartridge_uid는 필수입니다")
	}

	result, err := h.svc.ValidateCartridge(ctx, req.CartridgeUid, req.CategoryCode, req.TypeIndex, req.UserId)
	if err != nil {
		return nil, toGRPC(err)
	}

	resp := &v1.ValidateCartridgeResponse{
		IsValid:       result.IsValid,
		Reason:        result.Reason,
		RemainingUses: result.RemainingUses,
		AccessLevel:   serviceAccessLevelToProto(result.AccessLevel),
	}

	if result.Detail != nil {
		resp.Detail = cartridgeDetailToProto(result.Detail)
	}

	return resp, nil
}

// ============================================================================
// 헬퍼 함수
// ============================================================================

func cartridgeDetailToProto(d *service.CartridgeDetail) *v1.CartridgeDetail {
	return &v1.CartridgeDetail{
		CartridgeUid:        d.CartridgeUID,
		CategoryCode:        d.CategoryCode,
		TypeIndex:           d.TypeIndex,
		LegacyCode:          d.LegacyCode,
		NameKo:              d.NameKO,
		NameEn:              d.NameEN,
		LotId:               d.LotID,
		ExpiryDate:          d.ExpiryDate,
		RemainingUses:       d.RemainingUses,
		MaxUses:             d.MaxUses,
		AlphaCoefficient:    d.AlphaCoefficient,
		TempCoefficient:     d.TempCoefficient,
		HumidityCoefficient: d.HumidityCoefficient,
		RequiredChannels:    d.RequiredChannels,
		MeasurementSecs:     d.MeasurementSecs,
		Unit:                d.Unit,
		ReferenceRange:      d.ReferenceRange,
		IsValid:             d.IsValid,
		ValidationMessage:   d.ValidationMessage,
	}
}

func usageRecordToProto(r *service.CartridgeUsageRecord) *v1.CartridgeUsageRecord {
	return &v1.CartridgeUsageRecord{
		RecordId:     r.RecordID,
		UserId:       r.UserID,
		SessionId:    r.SessionID,
		CartridgeUid: r.CartridgeUID,
		CategoryCode: r.CategoryCode,
		TypeIndex:    r.TypeIndex,
		TypeNameKo:   r.TypeNameKO,
		UsedAt:       timestamppb.New(r.UsedAt),
	}
}

func typeInfoToProto(t *service.CartridgeTypeInfo) *v1.CartridgeTypeInfo {
	return &v1.CartridgeTypeInfo{
		CategoryCode:     t.CategoryCode,
		TypeIndex:        t.TypeIndex,
		LegacyCode:       t.LegacyCode,
		NameEn:           t.NameEN,
		NameKo:           t.NameKO,
		Description:      t.Description,
		RequiredChannels: t.RequiredChannels,
		MeasurementSecs:  t.MeasurementSecs,
		Unit:             t.Unit,
		ReferenceRange:   t.ReferenceRange,
		IsActive:         t.IsActive,
		IsBeta:           t.IsBeta,
		Manufacturer:     t.Manufacturer,
	}
}

func serviceAccessLevelToProto(level service.CartridgeAccessLevel) v1.CartridgeAccessLevel {
	switch level {
	case service.AccessIncluded:
		return v1.CartridgeAccessLevel_CARTRIDGE_ACCESS_INCLUDED
	case service.AccessLimited:
		return v1.CartridgeAccessLevel_CARTRIDGE_ACCESS_LIMITED
	case service.AccessAddOn:
		return v1.CartridgeAccessLevel_CARTRIDGE_ACCESS_ADD_ON
	case service.AccessRestricted:
		return v1.CartridgeAccessLevel_CARTRIDGE_ACCESS_RESTRICTED
	case service.AccessBeta:
		return v1.CartridgeAccessLevel_CARTRIDGE_ACCESS_BETA
	default:
		return v1.CartridgeAccessLevel_CARTRIDGE_ACCESS_UNKNOWN
	}
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
