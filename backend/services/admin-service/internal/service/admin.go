// Package service는 admin-service의 비즈니스 로직을 구현합니다.
package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/manpasik/backend/shared/errors"
	"go.uber.org/zap"
)

// AdminRole은 관리자 역할입니다.
type AdminRole int

const (
	RoleUnknown    AdminRole = iota
	RoleSuperAdmin           // 최고 관리자
	RoleAdmin                // 관리자
	RoleModerator            // 중재자
	RoleSupport              // 지원 담당
	RoleAnalyst              // 분석 담당
)

// AuditAction은 감사 로그 액션 유형입니다.
type AuditAction int

const (
	ActionUnknown      AuditAction = iota
	ActionLogin                    // 로그인
	ActionLogout                   // 로그아웃
	ActionCreate                   // 생성
	ActionUpdate                   // 수정
	ActionDelete                   // 삭제
	ActionConfigChange             // 설정 변경
	ActionUserBan                  // 사용자 차단
	ActionUserUnban                // 사용자 차단 해제
	ActionRoleChange               // 역할 변경
)

// AdminUser는 관리자 사용자 엔티티입니다.
type AdminUser struct {
	AdminID      string
	UserID       string
	Email        string
	DisplayName  string
	Role         AdminRole
	IsActive     bool
	CreatedAt    time.Time
	LastLoginAt  time.Time
	CreatedBy    string
	CountryCode  string
	RegionCode   string
	DistrictCode string
}

// AuditLogEntry는 감사 로그 항목입니다.
type AuditLogEntry struct {
	EntryID      string
	AdminID      string
	AdminEmail   string
	Action       AuditAction
	ResourceType string
	ResourceID   string
	Description  string
	IPAddress    string
	Timestamp    time.Time
}

// AdminUserSummary는 관리자가 조회하는 사용자 요약 정보입니다.
type AdminUserSummary struct {
	UserID           string
	Email            string
	DisplayName      string
	Tier             int
	IsActive         bool
	DeviceCount      int
	MeasurementCount int
	CreatedAt        time.Time
	LastActiveAt     time.Time
}

// SystemConfig는 시스템 설정 항목입니다.
type SystemConfig struct {
	Key         string
	Value       string
	Description string
	UpdatedBy   string
	UpdatedAt   time.Time
}

// SystemStats는 시스템 통계입니다.
type SystemStats struct {
	TotalUsers         int
	ActiveUsers        int
	TotalMeasurements  int
	TotalDevices       int
	TotalConsultations int
	TotalReservations  int
	UsersByTier        map[string]int
	MeasurementsByDay  map[string]int
	SystemHealthScore  float64
	CalculatedAt       time.Time
}

// ============================================================================
// Repository 인터페이스
// ============================================================================

// AdminRepository는 관리자 데이터 저장소 인터페이스입니다.
type AdminRepository interface {
	Save(ctx context.Context, admin *AdminUser) error
	FindByID(ctx context.Context, adminID string) (*AdminUser, error)
	FindAll(ctx context.Context, roleFilter AdminRole, activeOnly bool, limit, offset int32) ([]*AdminUser, int, error)
	Update(ctx context.Context, admin *AdminUser) error
	ListByRegion(ctx context.Context, countryCode, regionCode string) ([]*AdminUser, error)
}

// AuditLogRepository는 감사 로그 저장소 인터페이스입니다.
type AuditLogRepository interface {
	Save(ctx context.Context, entry *AuditLogEntry) error
	FindAll(ctx context.Context, adminID string, actionFilter AuditAction, startDate, endDate *time.Time, limit, offset int32) ([]*AuditLogEntry, int, error)
}

// SystemConfigRepository는 시스템 설정 저장소 인터페이스입니다.
type SystemConfigRepository interface {
	Save(ctx context.Context, cfg *SystemConfig) error
	FindByKey(ctx context.Context, key string) (*SystemConfig, error)
}

// UserSummaryRepository는 사용자 요약 정보 저장소 인터페이스입니다.
type UserSummaryRepository interface {
	FindAll(ctx context.Context, keyword, roleFilter string, activeOnly bool, limit, offset int32) ([]*AdminUserSummary, int, error)
}

// ============================================================================
// AdminService
// ============================================================================

// AuditLogStore는 확장된 감사 로그 저장소 인터페이스입니다.
// OldValue/NewValue 추적이 가능한 감사 로그를 지원합니다.
type AuditLogStore interface {
	Create(ctx context.Context, log *AuditLogDetail) error
	ListByAdmin(ctx context.Context, adminID string, limit, offset int) ([]*AuditLogDetail, error)
	ListByAction(ctx context.Context, action string, limit, offset int) ([]*AuditLogDetail, error)
	ListAll(ctx context.Context, limit, offset int) ([]*AuditLogDetail, error)
	Count(ctx context.Context) (int, error)
}

// AuditLogDetail은 확장된 감사 로그 엔티티입니다.
// OldValue/NewValue를 포함하여 변경 추적이 가능합니다.
type AuditLogDetail struct {
	ID        string
	AdminID   string
	Action    string // "config_update", "config_create", "config_delete", "user_ban", "system_restart" 등
	Resource  string // 대상 리소스 (예: "config:security.jwt_ttl_hours")
	OldValue  string
	NewValue  string
	IPAddress string
	UserAgent string
	CreatedAt time.Time
}

// AdminService는 관리자 비즈니스 로직입니다.
type AdminService struct {
	logger        *zap.Logger
	adminRepo     AdminRepository
	auditRepo     AuditLogRepository
	configRepo    SystemConfigRepository
	userRepo      UserSummaryRepository
	auditLogStore AuditLogStore // 확장 감사 로그 저장소 (선택)
}

// NewAdminService는 새 AdminService를 생성합니다.
func NewAdminService(
	logger *zap.Logger,
	adminRepo AdminRepository,
	auditRepo AuditLogRepository,
	configRepo SystemConfigRepository,
	userRepo UserSummaryRepository,
) *AdminService {
	return &AdminService{
		logger:     logger,
		adminRepo:  adminRepo,
		auditRepo:  auditRepo,
		configRepo: configRepo,
		userRepo:   userRepo,
	}
}

// SetAuditLogStore는 확장 감사 로그 저장소를 설정합니다.
func (s *AdminService) SetAuditLogStore(store AuditLogStore) {
	s.auditLogStore = store
}

// recordAuditDetail은 확장 감사 로그를 기록합니다.
// auditLogStore가 설정되지 않은 경우 무시합니다.
func (s *AdminService) recordAuditDetail(ctx context.Context, adminID, action, resource, oldValue, newValue string) {
	if s.auditLogStore == nil {
		return
	}
	now := time.Now().UTC()
	_ = s.auditLogStore.Create(ctx, &AuditLogDetail{
		ID:        uuid.New().String(),
		AdminID:   adminID,
		Action:    action,
		Resource:  resource,
		OldValue:  oldValue,
		NewValue:  newValue,
		CreatedAt: now,
	})
}

// GetAuditLogDetails는 확장 감사 로그를 조회합니다.
func (s *AdminService) GetAuditLogDetails(ctx context.Context, limit, offset int) ([]*AuditLogDetail, int, error) {
	if s.auditLogStore == nil {
		return nil, 0, nil
	}
	logs, err := s.auditLogStore.ListAll(ctx, limit, offset)
	if err != nil {
		s.logger.Error("확장 감사 로그 조회 실패", zap.Error(err))
		return nil, 0, err
	}
	count, _ := s.auditLogStore.Count(ctx)
	return logs, count, nil
}

// GetAuditLogDetailsByAdmin은 특정 관리자의 확장 감사 로그를 조회합니다.
func (s *AdminService) GetAuditLogDetailsByAdmin(ctx context.Context, adminID string, limit, offset int) ([]*AuditLogDetail, error) {
	if s.auditLogStore == nil {
		return nil, nil
	}
	return s.auditLogStore.ListByAdmin(ctx, adminID, limit, offset)
}

// GetAuditLogDetailsByAction은 특정 액션 유형의 확장 감사 로그를 조회합니다.
func (s *AdminService) GetAuditLogDetailsByAction(ctx context.Context, action string, limit, offset int) ([]*AuditLogDetail, error) {
	if s.auditLogStore == nil {
		return nil, nil
	}
	return s.auditLogStore.ListByAction(ctx, action, limit, offset)
}

// CreateAdmin은 새 관리자를 생성합니다.
func (s *AdminService) CreateAdmin(ctx context.Context, userID, email, displayName string, role AdminRole, createdBy string) (*AdminUser, error) {
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}
	if email == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "email은 필수입니다")
	}

	now := time.Now().UTC()
	admin := &AdminUser{
		AdminID:     uuid.New().String(),
		UserID:      userID,
		Email:       email,
		DisplayName: displayName,
		Role:        role,
		IsActive:    true,
		CreatedAt:   now,
		LastLoginAt: now,
		CreatedBy:   createdBy,
	}

	if err := s.adminRepo.Save(ctx, admin); err != nil {
		s.logger.Error("관리자 생성 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "관리자 생성에 실패했습니다")
	}

	// 감사 로그 기록
	_ = s.auditRepo.Save(ctx, &AuditLogEntry{
		EntryID:      uuid.New().String(),
		AdminID:      createdBy,
		AdminEmail:   "",
		Action:       ActionCreate,
		ResourceType: "admin",
		ResourceID:   admin.AdminID,
		Description:  "관리자 생성: " + email,
		Timestamp:    now,
	})

	// 확장 감사 로그 기록
	s.recordAuditDetail(ctx, createdBy, "admin_create", "admin:"+admin.AdminID, "", email)

	s.logger.Info("관리자 생성 완료", zap.String("admin_id", admin.AdminID), zap.String("email", email))
	return admin, nil
}

// GetAdmin은 관리자 정보를 조회합니다.
func (s *AdminService) GetAdmin(ctx context.Context, adminID string) (*AdminUser, error) {
	if adminID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "admin_id는 필수입니다")
	}

	admin, err := s.adminRepo.FindByID(ctx, adminID)
	if err != nil {
		s.logger.Error("관리자 조회 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "관리자 조회에 실패했습니다")
	}
	if admin == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "관리자를 찾을 수 없습니다")
	}
	return admin, nil
}

// ListAdmins는 관리자 목록을 조회합니다.
func (s *AdminService) ListAdmins(ctx context.Context, roleFilter AdminRole, activeOnly bool, limit, offset int32) ([]*AdminUser, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	admins, total, err := s.adminRepo.FindAll(ctx, roleFilter, activeOnly, limit, offset)
	if err != nil {
		s.logger.Error("관리자 목록 조회 실패", zap.Error(err))
		return nil, 0, apperrors.New(apperrors.ErrInternal, "관리자 목록 조회에 실패했습니다")
	}
	return admins, total, nil
}

// ListAdminsByRegion은 국가 코드와 지역 코드로 관리자 목록을 조회합니다.
func (s *AdminService) ListAdminsByRegion(ctx context.Context, countryCode, regionCode string) ([]*AdminUser, error) {
	if countryCode == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "country_code는 필수입니다")
	}

	admins, err := s.adminRepo.ListByRegion(ctx, countryCode, regionCode)
	if err != nil {
		s.logger.Error("지역별 관리자 조회 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "지역별 관리자 조회에 실패했습니다")
	}
	return admins, nil
}

// UpdateAdminRole은 관리자 역할을 변경합니다.
func (s *AdminService) UpdateAdminRole(ctx context.Context, adminID string, newRole AdminRole, updatedBy string) (*AdminUser, error) {
	if adminID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "admin_id는 필수입니다")
	}

	admin, err := s.adminRepo.FindByID(ctx, adminID)
	if err != nil {
		s.logger.Error("관리자 조회 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "관리자 조회에 실패했습니다")
	}
	if admin == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "관리자를 찾을 수 없습니다")
	}

	oldRole := admin.Role
	admin.Role = newRole

	if err := s.adminRepo.Update(ctx, admin); err != nil {
		s.logger.Error("관리자 역할 변경 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "관리자 역할 변경에 실패했습니다")
	}

	// 감사 로그 기록
	now := time.Now().UTC()
	_ = s.auditRepo.Save(ctx, &AuditLogEntry{
		EntryID:      uuid.New().String(),
		AdminID:      updatedBy,
		Action:       ActionRoleChange,
		ResourceType: "admin",
		ResourceID:   adminID,
		Description:  "역할 변경: " + roleToString(oldRole) + " → " + roleToString(newRole),
		Timestamp:    now,
	})

	// 확장 감사 로그 기록
	s.recordAuditDetail(ctx, updatedBy, "role_change", "admin:"+adminID, roleToString(oldRole), roleToString(newRole))

	s.logger.Info("관리자 역할 변경 완료",
		zap.String("admin_id", adminID),
		zap.String("old_role", roleToString(oldRole)),
		zap.String("new_role", roleToString(newRole)),
	)
	return admin, nil
}

// DeactivateAdmin은 관리자를 비활성화합니다.
func (s *AdminService) DeactivateAdmin(ctx context.Context, adminID, deactivatedBy, reason string) (bool, string, error) {
	if adminID == "" {
		return false, "", apperrors.New(apperrors.ErrInvalidInput, "admin_id는 필수입니다")
	}

	admin, err := s.adminRepo.FindByID(ctx, adminID)
	if err != nil {
		s.logger.Error("관리자 조회 실패", zap.Error(err))
		return false, "", apperrors.New(apperrors.ErrInternal, "관리자 조회에 실패했습니다")
	}
	if admin == nil {
		return false, "", apperrors.New(apperrors.ErrNotFound, "관리자를 찾을 수 없습니다")
	}

	admin.IsActive = false
	if err := s.adminRepo.Update(ctx, admin); err != nil {
		s.logger.Error("관리자 비활성화 실패", zap.Error(err))
		return false, "", apperrors.New(apperrors.ErrInternal, "관리자 비활성화에 실패했습니다")
	}

	// 감사 로그 기록
	now := time.Now().UTC()
	_ = s.auditRepo.Save(ctx, &AuditLogEntry{
		EntryID:      uuid.New().String(),
		AdminID:      deactivatedBy,
		Action:       ActionDelete,
		ResourceType: "admin",
		ResourceID:   adminID,
		Description:  "관리자 비활성화: " + reason,
		Timestamp:    now,
	})

	// 확장 감사 로그 기록
	s.recordAuditDetail(ctx, deactivatedBy, "admin_deactivate", "admin:"+adminID, "active", "inactive")

	s.logger.Info("관리자 비활성화 완료", zap.String("admin_id", adminID))
	return true, "관리자가 비활성화되었습니다", nil
}

// ListUsers는 관리자용 사용자 목록을 조회합니다.
func (s *AdminService) ListUsers(ctx context.Context, keyword, roleFilter string, activeOnly bool, limit, offset int32) ([]*AdminUserSummary, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	users, total, err := s.userRepo.FindAll(ctx, keyword, roleFilter, activeOnly, limit, offset)
	if err != nil {
		s.logger.Error("사용자 목록 조회 실패", zap.Error(err))
		return nil, 0, apperrors.New(apperrors.ErrInternal, "사용자 목록 조회에 실패했습니다")
	}
	return users, total, nil
}

// GetSystemStats는 시스템 통계를 반환합니다.
func (s *AdminService) GetSystemStats(ctx context.Context, days int32) (*SystemStats, error) {
	// 인메모리 구현: 더미 통계 반환 (실 DB 연동 시 교체)
	users, totalUsers, _ := s.userRepo.FindAll(ctx, "", "", false, 1000, 0)

	activeUsers := 0
	usersByTier := map[string]int{
		"free":     0,
		"basic":    0,
		"pro":      0,
		"clinical": 0,
	}
	totalDevices := 0
	totalMeasurements := 0

	for _, u := range users {
		if u.IsActive {
			activeUsers++
		}
		switch u.Tier {
		case 0:
			usersByTier["free"]++
		case 1:
			usersByTier["basic"]++
		case 2:
			usersByTier["pro"]++
		case 3:
			usersByTier["clinical"]++
		}
		totalDevices += u.DeviceCount
		totalMeasurements += u.MeasurementCount
	}

	now := time.Now().UTC()
	measurementsByDay := map[string]int{
		now.AddDate(0, 0, -1).Format("2006-01-02"): 120,
		now.Format("2006-01-02"):                    85,
	}

	return &SystemStats{
		TotalUsers:         totalUsers,
		ActiveUsers:        activeUsers,
		TotalMeasurements:  totalMeasurements,
		TotalDevices:       totalDevices,
		TotalConsultations: 15,
		TotalReservations:  8,
		UsersByTier:        usersByTier,
		MeasurementsByDay:  measurementsByDay,
		SystemHealthScore:  98.5,
		CalculatedAt:       now,
	}, nil
}

// GetAuditLog는 감사 로그를 조회합니다.
func (s *AdminService) GetAuditLog(ctx context.Context, adminID string, actionFilter AuditAction, startDate, endDate *time.Time, limit, offset int32) ([]*AuditLogEntry, int, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	entries, total, err := s.auditRepo.FindAll(ctx, adminID, actionFilter, startDate, endDate, limit, offset)
	if err != nil {
		s.logger.Error("감사 로그 조회 실패", zap.Error(err))
		return nil, 0, apperrors.New(apperrors.ErrInternal, "감사 로그 조회에 실패했습니다")
	}
	return entries, total, nil
}

// SetSystemConfig는 시스템 설정을 저장합니다.
func (s *AdminService) SetSystemConfig(ctx context.Context, key, value, description, updatedBy string) (*SystemConfig, error) {
	if key == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "설정 key는 필수입니다")
	}
	if strings.TrimSpace(value) == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "설정 value는 필수입니다")
	}

	// 이전 값 조회 (확장 감사 로그용)
	oldValue := ""
	if oldCfg, err := s.configRepo.FindByKey(ctx, key); err == nil && oldCfg != nil {
		oldValue = oldCfg.Value
	}

	now := time.Now().UTC()
	cfg := &SystemConfig{
		Key:         key,
		Value:       value,
		Description: description,
		UpdatedBy:   updatedBy,
		UpdatedAt:   now,
	}

	if err := s.configRepo.Save(ctx, cfg); err != nil {
		s.logger.Error("시스템 설정 저장 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "시스템 설정 저장에 실패했습니다")
	}

	// 감사 로그 기록
	_ = s.auditRepo.Save(ctx, &AuditLogEntry{
		EntryID:      uuid.New().String(),
		AdminID:      updatedBy,
		Action:       ActionConfigChange,
		ResourceType: "config",
		ResourceID:   key,
		Description:  "시스템 설정 변경: " + key + " = " + value,
		Timestamp:    now,
	})

	// 확장 감사 로그: OldValue/NewValue 추적
	action := "config_update"
	if oldValue == "" {
		action = "config_create"
	}
	s.recordAuditDetail(ctx, updatedBy, action, "config:"+key, oldValue, value)

	s.logger.Info("시스템 설정 저장 완료", zap.String("key", key))
	return cfg, nil
}

// GetSystemConfig는 시스템 설정을 조회합니다.
func (s *AdminService) GetSystemConfig(ctx context.Context, key string) (*SystemConfig, error) {
	if key == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "설정 key는 필수입니다")
	}

	cfg, err := s.configRepo.FindByKey(ctx, key)
	if err != nil {
		s.logger.Error("시스템 설정 조회 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "시스템 설정 조회에 실패했습니다")
	}
	if cfg == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "설정을 찾을 수 없습니다: "+key)
	}
	return cfg, nil
}

// ============================================================================
// 매출/재고 통계
// ============================================================================

// RevenueStats는 매출 통계입니다.
type RevenueStats struct {
	TotalRevenueKRW        int64
	SubscriptionRevenueKRW int64
	ProductRevenueKRW      int64
	TotalTransactions      int32
	Periods                []RevenuePeriod
	RevenueByTier          map[string]int64
}

// RevenuePeriod는 기간별 매출입니다.
type RevenuePeriod struct {
	Label            string
	RevenueKRW       int64
	TransactionCount int32
}

// InventoryStats는 재고 통계입니다.
type InventoryStats struct {
	Items           []InventoryItem
	TotalProducts   int32
	LowStockCount   int32
	OutOfStockCount int32
}

// InventoryItem은 재고 항목입니다.
type InventoryItem struct {
	ProductID         string
	ProductName       string
	Category          int32
	CurrentStock      int32
	MinStockThreshold int32
	MonthlySales      int32
	PriceKRW          int32
	Status            string
}

// GetRevenueStats는 매출 통계를 반환합니다.
func (s *AdminService) GetRevenueStats(ctx context.Context, period, startDate, endDate string) (*RevenueStats, error) {
	now := time.Now().UTC()

	periods := []RevenuePeriod{
		{Label: now.AddDate(0, -2, 0).Format("2006-01"), RevenueKRW: 12500000, TransactionCount: 45},
		{Label: now.AddDate(0, -1, 0).Format("2006-01"), RevenueKRW: 15800000, TransactionCount: 52},
		{Label: now.Format("2006-01"), RevenueKRW: 8900000, TransactionCount: 31},
	}

	return &RevenueStats{
		TotalRevenueKRW:        37200000,
		SubscriptionRevenueKRW: 28500000,
		ProductRevenueKRW:      8700000,
		TotalTransactions:      128,
		Periods:                periods,
		RevenueByTier: map[string]int64{
			"free": 0, "basic": 9500000, "pro": 14000000, "clinical": 5000000,
		},
	}, nil
}

// GetInventoryStats는 재고 통계를 반환합니다.
func (s *AdminService) GetInventoryStats(ctx context.Context, categoryFilter int32) (*InventoryStats, error) {
	items := []InventoryItem{
		{ProductID: "CTR-BG-001", ProductName: "혈당 카트리지 (50회)", Category: 1, CurrentStock: 1200, MinStockThreshold: 200, MonthlySales: 350, PriceKRW: 35000, Status: "in_stock"},
		{ProductID: "CTR-CL-001", ProductName: "콜레스테롤 카트리지 (25회)", Category: 1, CurrentStock: 85, MinStockThreshold: 100, MonthlySales: 120, PriceKRW: 45000, Status: "low_stock"},
		{ProductID: "DEV-MPS-001", ProductName: "만파식적 리더기 v3", Category: 2, CurrentStock: 320, MinStockThreshold: 50, MonthlySales: 40, PriceKRW: 299000, Status: "in_stock"},
		{ProductID: "CTR-UA-001", ProductName: "요산 카트리지 (25회)", Category: 1, CurrentStock: 0, MinStockThreshold: 100, MonthlySales: 80, PriceKRW: 42000, Status: "out_of_stock"},
		{ProductID: "CTR-HB-001", ProductName: "헤모글로빈 카트리지 (25회)", Category: 1, CurrentStock: 450, MinStockThreshold: 150, MonthlySales: 200, PriceKRW: 38000, Status: "in_stock"},
	}

	if categoryFilter > 0 {
		filtered := make([]InventoryItem, 0)
		for _, item := range items {
			if item.Category == categoryFilter {
				filtered = append(filtered, item)
			}
		}
		items = filtered
	}

	lowStock := int32(0)
	outOfStock := int32(0)
	for _, item := range items {
		if item.Status == "low_stock" {
			lowStock++
		}
		if item.Status == "out_of_stock" {
			outOfStock++
		}
	}

	return &InventoryStats{
		Items:           items,
		TotalProducts:   int32(len(items)),
		LowStockCount:   lowStock,
		OutOfStockCount: outOfStock,
	}, nil
}

// ============================================================================
// 헬퍼 함수
// ============================================================================

func roleToString(r AdminRole) string {
	switch r {
	case RoleSuperAdmin:
		return "super_admin"
	case RoleAdmin:
		return "admin"
	case RoleModerator:
		return "moderator"
	case RoleSupport:
		return "support"
	case RoleAnalyst:
		return "analyst"
	default:
		return "unknown"
	}
}
