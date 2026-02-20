// Package handler는 admin-service의 gRPC 핸들러입니다.
package handler

import (
	"context"

	"github.com/manpasik/backend/services/admin-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AdminHandler는 AdminService gRPC 서버를 구현합니다.
type AdminHandler struct {
	v1.UnimplementedAdminServiceServer
	svc       *service.AdminService
	cfgMgr    *service.ConfigManager
	log       *zap.Logger
}

// NewAdminHandler는 AdminHandler를 생성합니다.
func NewAdminHandler(svc *service.AdminService, log *zap.Logger) *AdminHandler {
	return &AdminHandler{svc: svc, log: log}
}

// SetConfigManager는 ConfigManager를 설정합니다.
func (h *AdminHandler) SetConfigManager(cfgMgr *service.ConfigManager) {
	h.cfgMgr = cfgMgr
}

// CreateAdmin은 관리자 생성 RPC입니다.
func (h *AdminHandler) CreateAdmin(ctx context.Context, req *v1.CreateAdminRequest) (*v1.AdminUser, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "요청이 비어 있습니다")
	}

	role := protoAdminRoleToService(req.Role)
	admin, err := h.svc.CreateAdmin(ctx, "", req.Email, req.DisplayName, role, "")
	if err != nil {
		return nil, toGRPC(err)
	}

	return adminUserToProto(admin), nil
}

// GetAdmin은 관리자 조회 RPC입니다.
func (h *AdminHandler) GetAdmin(ctx context.Context, req *v1.GetAdminRequest) (*v1.AdminUser, error) {
	if req == nil || req.AdminId == "" {
		return nil, status.Error(codes.InvalidArgument, "admin_id는 필수입니다")
	}

	admin, err := h.svc.GetAdmin(ctx, req.AdminId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return adminUserToProto(admin), nil
}

// ListAdmins는 관리자 목록 조회 RPC입니다.
func (h *AdminHandler) ListAdmins(ctx context.Context, req *v1.ListAdminsRequest) (*v1.ListAdminsResponse, error) {
	if req == nil {
		req = &v1.ListAdminsRequest{}
	}

	roleFilter := protoAdminRoleToService(req.RoleFilter)
	admins, total, err := h.svc.ListAdmins(ctx, roleFilter, false, req.Limit, req.Offset)
	if err != nil {
		return nil, toGRPC(err)
	}

	protoAdmins := make([]*v1.AdminUser, 0, len(admins))
	for _, a := range admins {
		protoAdmins = append(protoAdmins, adminUserToProto(a))
	}

	return &v1.ListAdminsResponse{
		Admins:     protoAdmins,
		TotalCount: int32(total),
	}, nil
}

// UpdateAdminRole은 관리자 역할 변경 RPC입니다.
func (h *AdminHandler) UpdateAdminRole(ctx context.Context, req *v1.UpdateAdminRoleRequest) (*v1.AdminUser, error) {
	if req == nil || req.AdminId == "" {
		return nil, status.Error(codes.InvalidArgument, "admin_id는 필수입니다")
	}

	newRole := protoAdminRoleToService(req.NewRole)
	admin, err := h.svc.UpdateAdminRole(ctx, req.AdminId, newRole, "")
	if err != nil {
		return nil, toGRPC(err)
	}

	return adminUserToProto(admin), nil
}

// DeactivateAdmin은 관리자 비활성화 RPC입니다.
func (h *AdminHandler) DeactivateAdmin(ctx context.Context, req *v1.DeactivateAdminRequest) (*v1.AdminUser, error) {
	if req == nil || req.AdminId == "" {
		return nil, status.Error(codes.InvalidArgument, "admin_id는 필수입니다")
	}

	_, _, err := h.svc.DeactivateAdmin(ctx, req.AdminId, "", "")
	if err != nil {
		return nil, toGRPC(err)
	}

	admin, err := h.svc.GetAdmin(ctx, req.AdminId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return adminUserToProto(admin), nil
}

// ListUsers는 관리자용 사용자 목록 조회 RPC입니다.
func (h *AdminHandler) ListUsers(ctx context.Context, req *v1.AdminListUsersRequest) (*v1.AdminListUsersResponse, error) {
	if req == nil {
		req = &v1.AdminListUsersRequest{}
	}

	users, total, err := h.svc.ListUsers(ctx, req.Query, "", false, req.Limit, req.Offset)
	if err != nil {
		return nil, toGRPC(err)
	}

	protoUsers := make([]*v1.AdminUserSummary, 0, len(users))
	for _, u := range users {
		protoUsers = append(protoUsers, userSummaryToProto(u))
	}

	return &v1.AdminListUsersResponse{
		Users:      protoUsers,
		TotalCount: int32(total),
	}, nil
}

// GetSystemStats는 시스템 통계 조회 RPC입니다.
func (h *AdminHandler) GetSystemStats(ctx context.Context, req *v1.GetSystemStatsRequest) (*v1.GetSystemStatsResponse, error) {
	if req == nil {
		req = &v1.GetSystemStatsRequest{}
	}

	stats, err := h.svc.GetSystemStats(ctx, 0)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.GetSystemStatsResponse{
		TotalUsers:        int32(stats.TotalUsers),
		ActiveUsers:       int32(stats.ActiveUsers),
		TotalDevices:      int32(stats.TotalDevices),
		TotalMeasurements: int32(stats.TotalMeasurements),
	}, nil
}

// GetAuditLog는 감사 로그 조회 RPC입니다.
func (h *AdminHandler) GetAuditLog(ctx context.Context, req *v1.GetAuditLogRequest) (*v1.GetAuditLogResponse, error) {
	if req == nil {
		req = &v1.GetAuditLogRequest{}
	}

	actionFilter := protoAuditActionToService(req.ActionFilter)

	entries, total, err := h.svc.GetAuditLog(ctx, req.AdminId, actionFilter, nil, nil, req.Limit, req.Offset)
	if err != nil {
		return nil, toGRPC(err)
	}

	protoEntries := make([]*v1.AuditLogEntry, 0, len(entries))
	for _, e := range entries {
		protoEntries = append(protoEntries, auditLogEntryToProto(e))
	}

	return &v1.GetAuditLogResponse{
		Entries:    protoEntries,
		TotalCount: int32(total),
	}, nil
}

// SetSystemConfig는 시스템 설정 저장 RPC입니다.
func (h *AdminHandler) SetSystemConfig(ctx context.Context, req *v1.SetSystemConfigRequest) (*v1.SystemConfig, error) {
	if req == nil || req.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "key는 필수입니다")
	}

	cfg, err := h.svc.SetSystemConfig(ctx, req.Key, req.Value, req.Description, "")
	if err != nil {
		return nil, toGRPC(err)
	}

	return systemConfigToProto(cfg), nil
}

// GetSystemConfig는 시스템 설정 조회 RPC입니다.
func (h *AdminHandler) GetSystemConfig(ctx context.Context, req *v1.GetSystemConfigRequest) (*v1.SystemConfig, error) {
	if req == nil || req.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "key는 필수입니다")
	}

	cfg, err := h.svc.GetSystemConfig(ctx, req.Key)
	if err != nil {
		return nil, toGRPC(err)
	}

	return systemConfigToProto(cfg), nil
}

// ListAdminsByRegion은 지역별 관리자 목록 조회 RPC입니다.
func (h *AdminHandler) ListAdminsByRegion(ctx context.Context, req *v1.ListAdminsByRegionRequest) (*v1.ListAdminsResponse, error) {
	if req == nil || req.CountryCode == "" {
		return nil, status.Error(codes.InvalidArgument, "country_code는 필수입니다")
	}

	admins, err := h.svc.ListAdminsByRegion(ctx, req.CountryCode, req.RegionCode)
	if err != nil {
		return nil, toGRPC(err)
	}

	var pbAdmins []*v1.AdminUser
	for _, a := range admins {
		pbAdmins = append(pbAdmins, adminUserToProto(a))
	}

	return &v1.ListAdminsResponse{Admins: pbAdmins}, nil
}

// ============================================================================
// 타입 변환 헬퍼 함수
// ============================================================================

func adminUserToProto(a *service.AdminUser) *v1.AdminUser {
	return &v1.AdminUser{
		AdminId:     a.AdminID,
		Email:       a.Email,
		DisplayName: a.DisplayName,
		Role:        serviceAdminRoleToProto(a.Role),
		IsActive:    a.IsActive,
		CreatedAt:   timestamppb.New(a.CreatedAt),
		LastLoginAt: timestamppb.New(a.LastLoginAt),
	}
}

func userSummaryToProto(u *service.AdminUserSummary) *v1.AdminUserSummary {
	return &v1.AdminUserSummary{
		UserId:      u.UserID,
		DisplayName: u.DisplayName,
		Email:       u.Email,
		IsActive:    u.IsActive,
		CreatedAt:   timestamppb.New(u.CreatedAt),
	}
}

func auditLogEntryToProto(e *service.AuditLogEntry) *v1.AuditLogEntry {
	return &v1.AuditLogEntry{
		EntryId:      e.EntryID,
		AdminId:      e.AdminID,
		Action:       serviceAuditActionToProto(e.Action),
		ResourceType: e.ResourceType,
		ResourceId:   e.ResourceID,
		Details:      e.Description,
		IpAddress:    e.IPAddress,
		Timestamp:    timestamppb.New(e.Timestamp),
	}
}

func systemConfigToProto(c *service.SystemConfig) *v1.SystemConfig {
	return &v1.SystemConfig{
		Key:         c.Key,
		Value:       c.Value,
		Description: c.Description,
		UpdatedBy:   c.UpdatedBy,
		UpdatedAt:   timestamppb.New(c.UpdatedAt),
	}
}

// --- AdminRole 변환 ---

func protoAdminRoleToService(role v1.AdminRole) service.AdminRole {
	switch role {
	case v1.AdminRole_ADMIN_ROLE_SUPER:
		return service.RoleSuperAdmin
	case v1.AdminRole_ADMIN_ROLE_NATIONAL:
		return service.RoleAdmin
	case v1.AdminRole_ADMIN_ROLE_REGIONAL:
		return service.RoleModerator
	case v1.AdminRole_ADMIN_ROLE_BRANCH:
		return service.RoleSupport
	case v1.AdminRole_ADMIN_ROLE_STORE:
		return service.RoleAnalyst
	default:
		return service.RoleUnknown
	}
}

func serviceAdminRoleToProto(role service.AdminRole) v1.AdminRole {
	switch role {
	case service.RoleSuperAdmin:
		return v1.AdminRole_ADMIN_ROLE_SUPER
	case service.RoleAdmin:
		return v1.AdminRole_ADMIN_ROLE_NATIONAL
	case service.RoleModerator:
		return v1.AdminRole_ADMIN_ROLE_REGIONAL
	case service.RoleSupport:
		return v1.AdminRole_ADMIN_ROLE_BRANCH
	case service.RoleAnalyst:
		return v1.AdminRole_ADMIN_ROLE_STORE
	default:
		return v1.AdminRole_ADMIN_ROLE_UNKNOWN
	}
}

// --- AuditAction 변환 ---

func protoAuditActionToService(action v1.AuditAction) service.AuditAction {
	switch action {
	case v1.AuditAction_AUDIT_ACTION_LOGIN:
		return service.ActionLogin
	case v1.AuditAction_AUDIT_ACTION_LOGOUT:
		return service.ActionLogout
	case v1.AuditAction_AUDIT_ACTION_CREATE:
		return service.ActionCreate
	case v1.AuditAction_AUDIT_ACTION_UPDATE:
		return service.ActionUpdate
	case v1.AuditAction_AUDIT_ACTION_DELETE:
		return service.ActionDelete
	case v1.AuditAction_AUDIT_ACTION_CONFIG_CHANGE:
		return service.ActionConfigChange
	case v1.AuditAction_AUDIT_ACTION_ROLE_CHANGE:
		return service.ActionRoleChange
	default:
		return service.ActionUnknown
	}
}

func serviceAuditActionToProto(action service.AuditAction) v1.AuditAction {
	switch action {
	case service.ActionLogin:
		return v1.AuditAction_AUDIT_ACTION_LOGIN
	case service.ActionLogout:
		return v1.AuditAction_AUDIT_ACTION_LOGOUT
	case service.ActionCreate:
		return v1.AuditAction_AUDIT_ACTION_CREATE
	case service.ActionUpdate:
		return v1.AuditAction_AUDIT_ACTION_UPDATE
	case service.ActionDelete:
		return v1.AuditAction_AUDIT_ACTION_DELETE
	case service.ActionConfigChange:
		return v1.AuditAction_AUDIT_ACTION_CONFIG_CHANGE
	case service.ActionRoleChange:
		return v1.AuditAction_AUDIT_ACTION_ROLE_CHANGE
	default:
		return v1.AuditAction_AUDIT_ACTION_UNKNOWN
	}
}

// --- SubscriptionTier 변환 ---

func serviceSubscriptionTierToProto(tier int) v1.SubscriptionTier {
	switch tier {
	case 1:
		return v1.SubscriptionTier_SUBSCRIPTION_TIER_BASIC
	case 2:
		return v1.SubscriptionTier_SUBSCRIPTION_TIER_PRO
	case 3:
		return v1.SubscriptionTier_SUBSCRIPTION_TIER_CLINICAL
	default:
		return v1.SubscriptionTier_SUBSCRIPTION_TIER_FREE
	}
}

// ============================================================================
// 설정 관리 확장 RPC
// ============================================================================

// ListSystemConfigs는 설정 목록 조회 RPC입니다 (메타데이터·번역 포함).
func (h *AdminHandler) ListSystemConfigs(ctx context.Context, req *v1.ListSystemConfigsRequest) (*v1.ListSystemConfigsResponse, error) {
	if h.cfgMgr == nil {
		return nil, status.Error(codes.Unimplemented, "ConfigManager가 설정되지 않았습니다")
	}
	if req == nil {
		req = &v1.ListSystemConfigsRequest{}
	}

	configs, categoryCounts, err := h.cfgMgr.ListSystemConfigs(ctx, req.LanguageCode, req.Category, req.IncludeSecrets)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "설정 목록 조회 실패: %v", err)
	}

	pbConfigs := make([]*v1.ConfigWithMeta, 0, len(configs))
	for _, c := range configs {
		pbConfigs = append(pbConfigs, configWithMetaToProto(c))
	}

	return &v1.ListSystemConfigsResponse{
		Configs:        pbConfigs,
		CategoryCounts: categoryCounts,
	}, nil
}

// GetConfigWithMeta는 단일 설정 조회 RPC입니다 (메타데이터·번역 포함).
func (h *AdminHandler) GetConfigWithMeta(ctx context.Context, req *v1.GetConfigWithMetaRequest) (*v1.ConfigWithMeta, error) {
	if h.cfgMgr == nil {
		return nil, status.Error(codes.Unimplemented, "ConfigManager가 설정되지 않았습니다")
	}
	if req == nil || req.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "key는 필수입니다")
	}

	cfg, err := h.cfgMgr.GetConfigWithMeta(ctx, req.Key, req.LanguageCode)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "설정 조회 실패: %v", err)
	}

	return configWithMetaToProto(cfg), nil
}

// ValidateConfigValue는 설정 값 유효성 검증 RPC입니다.
func (h *AdminHandler) ValidateConfigValue(ctx context.Context, req *v1.ValidateConfigValueRequest) (*v1.ValidateConfigValueResponse, error) {
	if h.cfgMgr == nil {
		return nil, status.Error(codes.Unimplemented, "ConfigManager가 설정되지 않았습니다")
	}
	if req == nil || req.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "key는 필수입니다")
	}

	result := h.cfgMgr.ValidateConfigValue(ctx, req.Key, req.Value)
	return &v1.ValidateConfigValueResponse{
		Valid:        result.Valid,
		ErrorMessage: result.ErrorMsg,
		Suggestions:  result.Suggestions,
	}, nil
}

// BulkSetConfigs는 일괄 설정 변경 RPC입니다.
func (h *AdminHandler) BulkSetConfigs(ctx context.Context, req *v1.BulkSetConfigsRequest) (*v1.BulkSetConfigsResponse, error) {
	if h.cfgMgr == nil {
		return nil, status.Error(codes.Unimplemented, "ConfigManager가 설정되지 않았습니다")
	}
	if req == nil || len(req.Configs) == 0 {
		return nil, status.Error(codes.InvalidArgument, "configs는 필수입니다")
	}

	configs := make([]struct{ Key, Value string }, 0, len(req.Configs))
	for _, c := range req.Configs {
		configs = append(configs, struct{ Key, Value string }{c.Key, c.Value})
	}

	successes, failures, errs := h.cfgMgr.BulkSetConfigs(ctx, configs, "", req.Reason)

	results := make([]*v1.ConfigChangeResult, 0, len(configs))
	errIdx := 0
	for _, c := range configs {
		cr := &v1.ConfigChangeResult{Key: c.Key, Success: true}
		if errIdx < len(errs) {
			for _, e := range errs {
				if len(e) > len(c.Key) && e[:len(c.Key)] == c.Key {
					cr.Success = false
					cr.ErrorMessage = e
					break
				}
			}
		}
		results = append(results, cr)
	}

	return &v1.BulkSetConfigsResponse{
		Results:      results,
		SuccessCount: int32(successes),
		FailureCount: int32(failures),
	}, nil
}

// configWithMetaToProto는 ConfigWithMeta를 protobuf로 변환합니다.
func configWithMetaToProto(c *service.ConfigWithMeta) *v1.ConfigWithMeta {
	pb := &v1.ConfigWithMeta{
		Key:       c.Key,
		Value:     c.Value,
		RawValue:  c.RawValue,
		UpdatedBy: c.UpdatedBy,
		UpdatedAt: timestamppb.New(c.UpdatedAt),
	}

	if c.Meta != nil {
		pb.Category = c.Meta.Category
		pb.ValueType = c.Meta.ValueType
		pb.SecurityLevel = c.Meta.SecurityLevel
		pb.IsRequired = c.Meta.IsRequired
		pb.DefaultValue = c.Meta.DefaultValue
		pb.AllowedValues = c.Meta.AllowedValues
		pb.ValidationRegex = c.Meta.ValidationRegex
		if c.Meta.ValidationMin != nil {
			pb.ValidationMin = *c.Meta.ValidationMin
		}
		if c.Meta.ValidationMax != nil {
			pb.ValidationMax = *c.Meta.ValidationMax
		}
		pb.DependsOn = c.Meta.DependsOn
		pb.DependsValue = c.Meta.DependsValue
		pb.EnvVarName = c.Meta.EnvVarName
		pb.ServiceName = c.Meta.ServiceName
		pb.RestartRequired = c.Meta.RestartRequired
	}

	if c.Translation != nil {
		pb.DisplayName = c.Translation.DisplayName
		pb.Description = c.Translation.Description
		pb.Placeholder = c.Translation.Placeholder
		pb.HelpText = c.Translation.HelpText
		pb.ValidationMessage = c.Translation.ValidationMessage
	}

	return pb
}

// GetAuditLogDetails는 확장 감사 로그(OldValue/NewValue 포함)를 조회하는 RPC입니다.
func (h *AdminHandler) GetAuditLogDetails(ctx context.Context, req *v1.GetAuditLogDetailsRequest) (*v1.GetAuditLogDetailsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "요청이 비어 있습니다")
	}

	limit := int(req.Limit)
	offset := int(req.Offset)
	if limit <= 0 {
		limit = 20
	}

	var details []*service.AuditLogDetail
	var total int
	var err error

	switch {
	case req.AdminId != "":
		details, err = h.svc.GetAuditLogDetailsByAdmin(ctx, req.AdminId, limit, offset)
	case req.Action != "":
		details, err = h.svc.GetAuditLogDetailsByAction(ctx, req.Action, limit, offset)
	default:
		details, total, err = h.svc.GetAuditLogDetails(ctx, limit, offset)
	}
	if err != nil {
		return nil, toGRPC(err)
	}

	var pbDetails []*v1.AuditLogDetail
	for _, d := range details {
		pbDetails = append(pbDetails, &v1.AuditLogDetail{
			Id:           d.ID,
			AdminId:      d.AdminID,
			Action:       d.Action,
			ResourceType: d.Resource,
			OldValue:     d.OldValue,
			NewValue:     d.NewValue,
			IpAddress:    d.IPAddress,
			Description:  d.UserAgent,
			CreatedAt:    timestamppb.New(d.CreatedAt),
		})
	}

	return &v1.GetAuditLogDetailsResponse{
		Details: pbDetails,
		Total:   int32(total),
	}, nil
}

// ============================================================================
// GetRevenueStats — 매출 통계 조회
// ============================================================================

func (h *AdminHandler) GetRevenueStats(ctx context.Context, req *v1.GetRevenueStatsRequest) (*v1.GetRevenueStatsResponse, error) {
	if req == nil {
		req = &v1.GetRevenueStatsRequest{}
	}

	stats, err := h.svc.GetRevenueStats(ctx, req.Period, req.StartDate, req.EndDate)
	if err != nil {
		return nil, toGRPC(err)
	}

	protoPeriods := make([]*v1.RevenuePeriod, 0, len(stats.Periods))
	for _, p := range stats.Periods {
		protoPeriods = append(protoPeriods, &v1.RevenuePeriod{
			Label:            p.Label,
			RevenueKrw:       p.RevenueKRW,
			TransactionCount: p.TransactionCount,
		})
	}

	return &v1.GetRevenueStatsResponse{
		TotalRevenueKrw:        stats.TotalRevenueKRW,
		SubscriptionRevenueKrw: stats.SubscriptionRevenueKRW,
		ProductRevenueKrw:      stats.ProductRevenueKRW,
		TotalTransactions:      stats.TotalTransactions,
		Periods:                protoPeriods,
		RevenueByTier:          stats.RevenueByTier,
	}, nil
}

// ============================================================================
// GetInventoryStats — 재고 통계 조회
// ============================================================================

func (h *AdminHandler) GetInventoryStats(ctx context.Context, req *v1.GetInventoryStatsRequest) (*v1.GetInventoryStatsResponse, error) {
	if req == nil {
		req = &v1.GetInventoryStatsRequest{}
	}

	stats, err := h.svc.GetInventoryStats(ctx, req.CategoryFilter)
	if err != nil {
		return nil, toGRPC(err)
	}

	protoItems := make([]*v1.InventoryItem, 0, len(stats.Items))
	for _, item := range stats.Items {
		protoItems = append(protoItems, &v1.InventoryItem{
			ProductId:         item.ProductID,
			ProductName:       item.ProductName,
			Category:          item.Category,
			CurrentStock:      item.CurrentStock,
			MinStockThreshold: item.MinStockThreshold,
			MonthlySales:      item.MonthlySales,
			PriceKrw:          item.PriceKRW,
			Status:            item.Status,
		})
	}

	return &v1.GetInventoryStatsResponse{
		Items:           protoItems,
		TotalProducts:   stats.TotalProducts,
		LowStockCount:   stats.LowStockCount,
		OutOfStockCount: stats.OutOfStockCount,
	}, nil
}

// --- toGRPC 에러 변환 ---

func toGRPC(err error) error {
	if err == nil {
		return nil
	}
	if appErr, ok := err.(*apperrors.AppError); ok {
		return appErr.ToGRPC()
	}
	if s, ok := status.FromError(err); ok {
		return s.Err()
	}
	return status.Error(codes.Internal, "내부 오류가 발생했습니다")
}
