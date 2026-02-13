// Package memory는 인메모리 관리자 저장소입니다 (개발/테스트용).
package memory

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/manpasik/backend/services/admin-service/internal/service"
)

// ============================================================================
// AdminRepository
// ============================================================================

// AdminRepository는 인메모리 관리자 저장소입니다.
type AdminRepository struct {
	mu     sync.RWMutex
	admins map[string]*service.AdminUser // key: AdminID
}

// NewAdminRepository는 인메모리 AdminRepository를 생성합니다.
// 기본 SuperAdmin 시드 데이터를 포함합니다.
func NewAdminRepository() *AdminRepository {
	repo := &AdminRepository{
		admins: make(map[string]*service.AdminUser),
	}

	// 시드 데이터: 기본 SuperAdmin
	now := time.Now().UTC()
	repo.admins["seed-admin-001"] = &service.AdminUser{
		AdminID:      "seed-admin-001",
		UserID:       "seed-user-001",
		Email:        "superadmin@manpasik.com",
		DisplayName:  "시스템 관리자",
		Role:         service.RoleSuperAdmin,
		IsActive:     true,
		CreatedAt:    now,
		LastLoginAt:  now,
		CreatedBy:    "system",
		CountryCode:  "KR",
		RegionCode:   "SEL",
		DistrictCode: "GN",
	}

	return repo
}

// Save는 관리자를 저장합니다.
func (r *AdminRepository) Save(_ context.Context, admin *service.AdminUser) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *admin
	r.admins[admin.AdminID] = &cp
	return nil
}

// FindByID는 관리자 ID로 조회합니다.
func (r *AdminRepository) FindByID(_ context.Context, adminID string) (*service.AdminUser, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.admins[adminID]
	if !ok {
		return nil, nil
	}
	cp := *a
	return &cp, nil
}

// FindAll은 관리자 목록을 조회합니다.
func (r *AdminRepository) FindAll(_ context.Context, roleFilter service.AdminRole, activeOnly bool, limit, offset int32) ([]*service.AdminUser, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*service.AdminUser
	for _, a := range r.admins {
		if roleFilter != service.RoleUnknown && a.Role != roleFilter {
			continue
		}
		if activeOnly && !a.IsActive {
			continue
		}
		cp := *a
		filtered = append(filtered, &cp)
	}

	total := len(filtered)

	// 페이지네이션
	start := int(offset)
	if start >= total {
		return nil, total, nil
	}
	end := start + int(limit)
	if end > total {
		end = total
	}

	return filtered[start:end], total, nil
}

// Update는 관리자 정보를 업데이트합니다.
func (r *AdminRepository) Update(_ context.Context, admin *service.AdminUser) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *admin
	r.admins[admin.AdminID] = &cp
	return nil
}

// ListByRegion은 국가 코드와 지역 코드로 관리자를 검색합니다.
func (r *AdminRepository) ListByRegion(_ context.Context, countryCode, regionCode string) ([]*service.AdminUser, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*service.AdminUser
	for _, a := range r.admins {
		if a.CountryCode != countryCode {
			continue
		}
		if regionCode != "" && a.RegionCode != regionCode {
			continue
		}
		cp := *a
		result = append(result, &cp)
	}
	return result, nil
}

// ============================================================================
// AuditLogRepository
// ============================================================================

// AuditLogRepository는 인메모리 감사 로그 저장소입니다.
type AuditLogRepository struct {
	mu      sync.RWMutex
	entries []*service.AuditLogEntry
}

// NewAuditLogRepository는 인메모리 AuditLogRepository를 생성합니다.
func NewAuditLogRepository() *AuditLogRepository {
	return &AuditLogRepository{
		entries: make([]*service.AuditLogEntry, 0),
	}
}

// Save는 감사 로그를 저장합니다.
func (r *AuditLogRepository) Save(_ context.Context, entry *service.AuditLogEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *entry
	r.entries = append(r.entries, &cp)
	return nil
}

// FindAll은 감사 로그를 조회합니다.
func (r *AuditLogRepository) FindAll(_ context.Context, adminID string, actionFilter service.AuditAction, startDate, endDate *time.Time, limit, offset int32) ([]*service.AuditLogEntry, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*service.AuditLogEntry
	for _, e := range r.entries {
		if adminID != "" && e.AdminID != adminID {
			continue
		}
		if actionFilter != service.ActionUnknown && e.Action != actionFilter {
			continue
		}
		if startDate != nil && e.Timestamp.Before(*startDate) {
			continue
		}
		if endDate != nil && e.Timestamp.After(*endDate) {
			continue
		}
		cp := *e
		filtered = append(filtered, &cp)
	}

	total := len(filtered)

	start := int(offset)
	if start >= total {
		return nil, total, nil
	}
	end := start + int(limit)
	if end > total {
		end = total
	}

	return filtered[start:end], total, nil
}

// ============================================================================
// SystemConfigRepository
// ============================================================================

// SystemConfigRepository는 인메모리 시스템 설정 저장소입니다.
type SystemConfigRepository struct {
	mu      sync.RWMutex
	configs map[string]*service.SystemConfig
}

// NewSystemConfigRepository는 인메모리 SystemConfigRepository를 생성합니다.
// 기본 설정 시드 데이터를 포함합니다.
func NewSystemConfigRepository() *SystemConfigRepository {
	repo := &SystemConfigRepository{
		configs: make(map[string]*service.SystemConfig),
	}

	now := time.Now().UTC()
	// 시드 데이터: 기본 시스템 설정
	repo.configs["maintenance_mode"] = &service.SystemConfig{
		Key:         "maintenance_mode",
		Value:       "false",
		Description: "시스템 유지보수 모드 활성화 여부",
		UpdatedBy:   "system",
		UpdatedAt:   now,
	}
	repo.configs["max_devices_per_user"] = &service.SystemConfig{
		Key:         "max_devices_per_user",
		Value:       "10",
		Description: "사용자당 최대 디바이스 수",
		UpdatedBy:   "system",
		UpdatedAt:   now,
	}
	repo.configs["default_language"] = &service.SystemConfig{
		Key:         "default_language",
		Value:       "ko",
		Description: "기본 언어 설정",
		UpdatedBy:   "system",
		UpdatedAt:   now,
	}

	return repo
}

// Save는 시스템 설정을 저장합니다.
func (r *SystemConfigRepository) Save(_ context.Context, cfg *service.SystemConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *cfg
	r.configs[cfg.Key] = &cp
	return nil
}

// FindByKey는 키로 시스템 설정을 조회합니다.
func (r *SystemConfigRepository) FindByKey(_ context.Context, key string) (*service.SystemConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.configs[key]
	if !ok {
		return nil, nil
	}
	cp := *c
	return &cp, nil
}

// ============================================================================
// UserSummaryRepository
// ============================================================================

// UserSummaryRepository는 인메모리 사용자 요약 정보 저장소입니다.
type UserSummaryRepository struct {
	mu    sync.RWMutex
	users []*service.AdminUserSummary
}

// NewUserSummaryRepository는 인메모리 UserSummaryRepository를 생성합니다.
// 샘플 사용자 데이터를 포함합니다.
func NewUserSummaryRepository() *UserSummaryRepository {
	now := time.Now().UTC()
	repo := &UserSummaryRepository{
		users: []*service.AdminUserSummary{
			{
				UserID:           "user-001",
				Email:            "hong@example.com",
				DisplayName:      "홍길동",
				Tier:             1,
				IsActive:         true,
				DeviceCount:      2,
				MeasurementCount: 45,
				CreatedAt:        now.AddDate(0, -3, 0),
				LastActiveAt:     now.AddDate(0, 0, -1),
			},
			{
				UserID:           "user-002",
				Email:            "kim@example.com",
				DisplayName:      "김철수",
				Tier:             2,
				IsActive:         true,
				DeviceCount:      3,
				MeasurementCount: 120,
				CreatedAt:        now.AddDate(0, -6, 0),
				LastActiveAt:     now,
			},
			{
				UserID:           "user-003",
				Email:            "lee@example.com",
				DisplayName:      "이영희",
				Tier:             0,
				IsActive:         false,
				DeviceCount:      1,
				MeasurementCount: 5,
				CreatedAt:        now.AddDate(-1, 0, 0),
				LastActiveAt:     now.AddDate(0, -2, 0),
			},
			{
				UserID:           "user-004",
				Email:            "park@example.com",
				DisplayName:      "박민수",
				Tier:             3,
				IsActive:         true,
				DeviceCount:      5,
				MeasurementCount: 230,
				CreatedAt:        now.AddDate(0, -1, 0),
				LastActiveAt:     now,
			},
			{
				UserID:           "user-005",
				Email:            "choi@example.com",
				DisplayName:      "최지은",
				Tier:             1,
				IsActive:         true,
				DeviceCount:      1,
				MeasurementCount: 30,
				CreatedAt:        now.AddDate(0, -2, 0),
				LastActiveAt:     now.AddDate(0, 0, -3),
			},
		},
	}
	return repo
}

// FindAll은 사용자 요약 목록을 조회합니다.
func (r *UserSummaryRepository) FindAll(_ context.Context, keyword, roleFilter string, activeOnly bool, limit, offset int32) ([]*service.AdminUserSummary, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*service.AdminUserSummary
	for _, u := range r.users {
		if activeOnly && !u.IsActive {
			continue
		}
		if keyword != "" {
			kw := strings.ToLower(keyword)
			if !strings.Contains(strings.ToLower(u.Email), kw) &&
				!strings.Contains(strings.ToLower(u.DisplayName), kw) {
				continue
			}
		}
		cp := *u
		filtered = append(filtered, &cp)
	}

	total := len(filtered)

	start := int(offset)
	if start >= total {
		return nil, total, nil
	}
	end := start + int(limit)
	if end > total {
		end = total
	}

	return filtered[start:end], total, nil
}
