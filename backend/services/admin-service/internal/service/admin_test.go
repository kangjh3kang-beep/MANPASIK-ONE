package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/manpasik/backend/services/admin-service/internal/service"
	"go.uber.org/zap"
)

// ============================================================================
// Fake Repositories
// ============================================================================

// fakeAdminRepo는 테스트용 인메모리 관리자 저장소입니다.
type fakeAdminRepo struct {
	admins map[string]*service.AdminUser
}

func newFakeAdminRepo() *fakeAdminRepo {
	return &fakeAdminRepo{admins: make(map[string]*service.AdminUser)}
}

func (r *fakeAdminRepo) Save(_ context.Context, admin *service.AdminUser) error {
	r.admins[admin.AdminID] = admin
	return nil
}

func (r *fakeAdminRepo) FindByID(_ context.Context, adminID string) (*service.AdminUser, error) {
	a, ok := r.admins[adminID]
	if !ok {
		return nil, nil
	}
	return a, nil
}

func (r *fakeAdminRepo) FindAll(_ context.Context, roleFilter service.AdminRole, activeOnly bool, limit, offset int32) ([]*service.AdminUser, int, error) {
	var result []*service.AdminUser
	for _, a := range r.admins {
		if roleFilter != service.RoleUnknown && a.Role != roleFilter {
			continue
		}
		if activeOnly && !a.IsActive {
			continue
		}
		result = append(result, a)
	}
	total := len(result)
	start := int(offset)
	if start >= total {
		return nil, total, nil
	}
	end := start + int(limit)
	if end > total {
		end = total
	}
	return result[start:end], total, nil
}

func (r *fakeAdminRepo) Update(_ context.Context, admin *service.AdminUser) error {
	r.admins[admin.AdminID] = admin
	return nil
}

func (r *fakeAdminRepo) ListByRegion(_ context.Context, countryCode, regionCode string) ([]*service.AdminUser, error) {
	var result []*service.AdminUser
	for _, a := range r.admins {
		if a.CountryCode != countryCode {
			continue
		}
		if regionCode != "" && a.RegionCode != regionCode {
			continue
		}
		result = append(result, a)
	}
	return result, nil
}

// fakeAuditRepo는 테스트용 인메모리 감사 로그 저장소입니다.
type fakeAuditRepo struct {
	entries []*service.AuditLogEntry
}

func newFakeAuditRepo() *fakeAuditRepo {
	return &fakeAuditRepo{entries: make([]*service.AuditLogEntry, 0)}
}

func (r *fakeAuditRepo) Save(_ context.Context, entry *service.AuditLogEntry) error {
	r.entries = append(r.entries, entry)
	return nil
}

func (r *fakeAuditRepo) FindAll(_ context.Context, adminID string, actionFilter service.AuditAction, startDate, endDate *time.Time, limit, offset int32) ([]*service.AuditLogEntry, int, error) {
	var result []*service.AuditLogEntry
	for _, e := range r.entries {
		if adminID != "" && e.AdminID != adminID {
			continue
		}
		if actionFilter != service.ActionUnknown && e.Action != actionFilter {
			continue
		}
		result = append(result, e)
	}
	total := len(result)
	start := int(offset)
	if start >= total {
		return nil, total, nil
	}
	end := start + int(limit)
	if end > total {
		end = total
	}
	return result[start:end], total, nil
}

// fakeConfigRepo는 테스트용 인메모리 시스템 설정 저장소입니다.
type fakeConfigRepo struct {
	configs map[string]*service.SystemConfig
}

func newFakeConfigRepo() *fakeConfigRepo {
	return &fakeConfigRepo{configs: make(map[string]*service.SystemConfig)}
}

func (r *fakeConfigRepo) Save(_ context.Context, cfg *service.SystemConfig) error {
	r.configs[cfg.Key] = cfg
	return nil
}

func (r *fakeConfigRepo) FindByKey(_ context.Context, key string) (*service.SystemConfig, error) {
	c, ok := r.configs[key]
	if !ok {
		return nil, nil
	}
	return c, nil
}

// fakeUserRepo는 테스트용 인메모리 사용자 요약 저장소입니다.
type fakeUserRepo struct {
	users []*service.AdminUserSummary
}

func newFakeUserRepo() *fakeUserRepo {
	now := time.Now().UTC()
	return &fakeUserRepo{
		users: []*service.AdminUserSummary{
			{UserID: "user-001", Email: "hong@example.com", DisplayName: "홍길동", Tier: 1, IsActive: true, DeviceCount: 2, MeasurementCount: 45, CreatedAt: now, LastActiveAt: now},
			{UserID: "user-002", Email: "kim@example.com", DisplayName: "김철수", Tier: 2, IsActive: true, DeviceCount: 3, MeasurementCount: 120, CreatedAt: now, LastActiveAt: now},
			{UserID: "user-003", Email: "lee@example.com", DisplayName: "이영희", Tier: 0, IsActive: false, DeviceCount: 1, MeasurementCount: 5, CreatedAt: now, LastActiveAt: now},
		},
	}
}

func (r *fakeUserRepo) FindAll(_ context.Context, keyword, roleFilter string, activeOnly bool, limit, offset int32) ([]*service.AdminUserSummary, int, error) {
	var result []*service.AdminUserSummary
	for _, u := range r.users {
		if activeOnly && !u.IsActive {
			continue
		}
		result = append(result, u)
	}
	total := len(result)
	start := int(offset)
	if start >= total {
		return nil, total, nil
	}
	end := start + int(limit)
	if end > total {
		end = total
	}
	return result[start:end], total, nil
}

// ============================================================================
// 테스트 헬퍼
// ============================================================================

func newTestService() (*service.AdminService, *fakeAdminRepo, *fakeAuditRepo, *fakeConfigRepo, *fakeUserRepo) {
	adminRepo := newFakeAdminRepo()
	auditRepo := newFakeAuditRepo()
	configRepo := newFakeConfigRepo()
	userRepo := newFakeUserRepo()
	svc := service.NewAdminService(zap.NewNop(), adminRepo, auditRepo, configRepo, userRepo)
	return svc, adminRepo, auditRepo, configRepo, userRepo
}

// ============================================================================
// 테스트
// ============================================================================

func TestCreateAdmin_Success(t *testing.T) {
	svc, _, _, _, _ := newTestService()
	ctx := context.Background()

	admin, err := svc.CreateAdmin(ctx, "user-100", "admin@test.com", "테스트 관리자", service.RoleAdmin, "creator-1")
	if err != nil {
		t.Fatalf("CreateAdmin 실패: %v", err)
	}
	if admin.AdminID == "" {
		t.Error("AdminID가 비어 있습니다")
	}
	if admin.UserID != "user-100" {
		t.Errorf("UserID: got %s, want user-100", admin.UserID)
	}
	if admin.Email != "admin@test.com" {
		t.Errorf("Email: got %s, want admin@test.com", admin.Email)
	}
	if admin.Role != service.RoleAdmin {
		t.Errorf("Role: got %d, want %d", admin.Role, service.RoleAdmin)
	}
	if !admin.IsActive {
		t.Error("새 관리자는 IsActive=true여야 합니다")
	}
}

func TestCreateAdmin_MissingUserID(t *testing.T) {
	svc, _, _, _, _ := newTestService()
	ctx := context.Background()

	_, err := svc.CreateAdmin(ctx, "", "admin@test.com", "관리자", service.RoleAdmin, "creator-1")
	if err == nil {
		t.Error("user_id 없이 관리자 생성이 허용되었습니다")
	}
}

func TestCreateAdmin_MissingEmail(t *testing.T) {
	svc, _, _, _, _ := newTestService()
	ctx := context.Background()

	_, err := svc.CreateAdmin(ctx, "user-101", "", "관리자", service.RoleAdmin, "creator-1")
	if err == nil {
		t.Error("email 없이 관리자 생성이 허용되었습니다")
	}
}

func TestGetAdmin_Success(t *testing.T) {
	svc, _, _, _, _ := newTestService()
	ctx := context.Background()

	created, _ := svc.CreateAdmin(ctx, "user-200", "get@test.com", "조회 테스트", service.RoleModerator, "creator-1")

	admin, err := svc.GetAdmin(ctx, created.AdminID)
	if err != nil {
		t.Fatalf("GetAdmin 실패: %v", err)
	}
	if admin.Email != "get@test.com" {
		t.Errorf("Email: got %s, want get@test.com", admin.Email)
	}
	if admin.Role != service.RoleModerator {
		t.Errorf("Role: got %d, want %d", admin.Role, service.RoleModerator)
	}
}

func TestGetAdmin_NotFound(t *testing.T) {
	svc, _, _, _, _ := newTestService()
	ctx := context.Background()

	_, err := svc.GetAdmin(ctx, "non-existent-id")
	if err == nil {
		t.Error("존재하지 않는 관리자 조회가 성공했습니다")
	}
}

func TestListAdmins_Success(t *testing.T) {
	svc, _, _, _, _ := newTestService()
	ctx := context.Background()

	_, _ = svc.CreateAdmin(ctx, "user-301", "list1@test.com", "관리자1", service.RoleAdmin, "creator-1")
	_, _ = svc.CreateAdmin(ctx, "user-302", "list2@test.com", "관리자2", service.RoleSupport, "creator-1")
	_, _ = svc.CreateAdmin(ctx, "user-303", "list3@test.com", "관리자3", service.RoleAdmin, "creator-1")

	admins, total, err := svc.ListAdmins(ctx, service.RoleUnknown, false, 10, 0)
	if err != nil {
		t.Fatalf("ListAdmins 실패: %v", err)
	}
	if total != 3 {
		t.Errorf("총 관리자 수: got %d, want 3", total)
	}
	if len(admins) != 3 {
		t.Errorf("반환된 관리자 수: got %d, want 3", len(admins))
	}
}

func TestUpdateAdminRole_Success(t *testing.T) {
	svc, _, auditRepo, _, _ := newTestService()
	ctx := context.Background()

	created, _ := svc.CreateAdmin(ctx, "user-400", "role@test.com", "역할 변경 테스트", service.RoleSupport, "creator-1")

	updated, err := svc.UpdateAdminRole(ctx, created.AdminID, service.RoleAdmin, "updater-1")
	if err != nil {
		t.Fatalf("UpdateAdminRole 실패: %v", err)
	}
	if updated.Role != service.RoleAdmin {
		t.Errorf("역할: got %d, want %d", updated.Role, service.RoleAdmin)
	}

	// 감사 로그 확인 (create + role_change = 2개)
	if len(auditRepo.entries) < 2 {
		t.Errorf("감사 로그 수: got %d, want >= 2", len(auditRepo.entries))
	}
	lastEntry := auditRepo.entries[len(auditRepo.entries)-1]
	if lastEntry.Action != service.ActionRoleChange {
		t.Errorf("감사 로그 액션: got %d, want %d", lastEntry.Action, service.ActionRoleChange)
	}
}

func TestDeactivateAdmin_Success(t *testing.T) {
	svc, _, _, _, _ := newTestService()
	ctx := context.Background()

	created, _ := svc.CreateAdmin(ctx, "user-500", "deactivate@test.com", "비활성화 테스트", service.RoleAnalyst, "creator-1")

	success, msg, err := svc.DeactivateAdmin(ctx, created.AdminID, "deactivator-1", "퇴사")
	if err != nil {
		t.Fatalf("DeactivateAdmin 실패: %v", err)
	}
	if !success {
		t.Error("비활성화 결과가 false입니다")
	}
	if msg == "" {
		t.Error("비활성화 메시지가 비어 있습니다")
	}

	// 비활성화 확인
	admin, _ := svc.GetAdmin(ctx, created.AdminID)
	if admin.IsActive {
		t.Error("비활성화 후에도 IsActive=true입니다")
	}
}

func TestListUsers_Success(t *testing.T) {
	svc, _, _, _, _ := newTestService()
	ctx := context.Background()

	users, total, err := svc.ListUsers(ctx, "", "", false, 10, 0)
	if err != nil {
		t.Fatalf("ListUsers 실패: %v", err)
	}
	if total < 1 {
		t.Errorf("총 사용자 수: got %d, want >= 1", total)
	}
	if len(users) < 1 {
		t.Errorf("반환된 사용자 수: got %d, want >= 1", len(users))
	}
}

func TestGetSystemStats_Success(t *testing.T) {
	svc, _, _, _, _ := newTestService()
	ctx := context.Background()

	stats, err := svc.GetSystemStats(ctx, 7)
	if err != nil {
		t.Fatalf("GetSystemStats 실패: %v", err)
	}
	if stats.TotalUsers < 1 {
		t.Errorf("TotalUsers: got %d, want >= 1", stats.TotalUsers)
	}
	if stats.SystemHealthScore <= 0 {
		t.Errorf("SystemHealthScore: got %f, want > 0", stats.SystemHealthScore)
	}
	if stats.CalculatedAt.IsZero() {
		t.Error("CalculatedAt가 zero입니다")
	}
}

func TestGetAuditLog_Success(t *testing.T) {
	svc, _, _, _, _ := newTestService()
	ctx := context.Background()

	// 관리자 생성을 통해 감사 로그 발생
	_, _ = svc.CreateAdmin(ctx, "user-700", "audit@test.com", "감사 로그 테스트", service.RoleAdmin, "creator-1")

	entries, total, err := svc.GetAuditLog(ctx, "", service.ActionUnknown, nil, nil, 50, 0)
	if err != nil {
		t.Fatalf("GetAuditLog 실패: %v", err)
	}
	if total < 1 {
		t.Errorf("감사 로그 수: got %d, want >= 1", total)
	}
	if len(entries) < 1 {
		t.Errorf("반환된 로그 수: got %d, want >= 1", len(entries))
	}
}

func TestSetSystemConfig_Success(t *testing.T) {
	svc, _, _, _, _ := newTestService()
	ctx := context.Background()

	cfg, err := svc.SetSystemConfig(ctx, "test_key", "test_value", "테스트 설정", "admin-1")
	if err != nil {
		t.Fatalf("SetSystemConfig 실패: %v", err)
	}
	if cfg.Key != "test_key" {
		t.Errorf("Key: got %s, want test_key", cfg.Key)
	}
	if cfg.Value != "test_value" {
		t.Errorf("Value: got %s, want test_value", cfg.Value)
	}
	if cfg.UpdatedBy != "admin-1" {
		t.Errorf("UpdatedBy: got %s, want admin-1", cfg.UpdatedBy)
	}
}

func TestGetSystemConfig_Success(t *testing.T) {
	svc, _, _, _, _ := newTestService()
	ctx := context.Background()

	_, _ = svc.SetSystemConfig(ctx, "get_key", "get_value", "조회 테스트", "admin-1")

	cfg, err := svc.GetSystemConfig(ctx, "get_key")
	if err != nil {
		t.Fatalf("GetSystemConfig 실패: %v", err)
	}
	if cfg.Value != "get_value" {
		t.Errorf("Value: got %s, want get_value", cfg.Value)
	}
}

func TestListAdminsByRegion(t *testing.T) {
	svc, adminRepo, _, _, _ := newTestService()
	ctx := context.Background()

	// 시드 데이터: 지역 정보가 있는 관리자
	adminRepo.admins["region-admin-1"] = &service.AdminUser{
		AdminID:      "region-admin-1",
		UserID:       "user-r1",
		Email:        "kr-sel@test.com",
		DisplayName:  "서울 관리자",
		Role:         service.RoleAdmin,
		IsActive:     true,
		CreatedAt:    time.Now().UTC(),
		CountryCode:  "KR",
		RegionCode:   "SEL",
		DistrictCode: "GN",
	}
	adminRepo.admins["region-admin-2"] = &service.AdminUser{
		AdminID:      "region-admin-2",
		UserID:       "user-r2",
		Email:        "kr-bus@test.com",
		DisplayName:  "부산 관리자",
		Role:         service.RoleAdmin,
		IsActive:     true,
		CreatedAt:    time.Now().UTC(),
		CountryCode:  "KR",
		RegionCode:   "BUS",
		DistrictCode: "HU",
	}
	adminRepo.admins["region-admin-3"] = &service.AdminUser{
		AdminID:      "region-admin-3",
		UserID:       "user-r3",
		Email:        "us-nyc@test.com",
		DisplayName:  "NY 관리자",
		Role:         service.RoleAdmin,
		IsActive:     true,
		CreatedAt:    time.Now().UTC(),
		CountryCode:  "US",
		RegionCode:   "NYC",
		DistrictCode: "MN",
	}

	// 한국 전체 조회
	admins, err := svc.ListAdminsByRegion(ctx, "KR", "")
	if err != nil {
		t.Fatalf("ListAdminsByRegion(KR, '') 실패: %v", err)
	}
	if len(admins) != 2 {
		t.Errorf("KR 관리자 수: got %d, want 2", len(admins))
	}

	// 한국 서울 조회
	admins, err = svc.ListAdminsByRegion(ctx, "KR", "SEL")
	if err != nil {
		t.Fatalf("ListAdminsByRegion(KR, SEL) 실패: %v", err)
	}
	if len(admins) != 1 {
		t.Errorf("KR-SEL 관리자 수: got %d, want 1", len(admins))
	}
	if len(admins) > 0 && admins[0].Email != "kr-sel@test.com" {
		t.Errorf("KR-SEL 관리자 이메일: got %s, want kr-sel@test.com", admins[0].Email)
	}

	// 미국 조회
	admins, err = svc.ListAdminsByRegion(ctx, "US", "")
	if err != nil {
		t.Fatalf("ListAdminsByRegion(US, '') 실패: %v", err)
	}
	if len(admins) != 1 {
		t.Errorf("US 관리자 수: got %d, want 1", len(admins))
	}

	// 빈 국가 코드 → 에러
	_, err = svc.ListAdminsByRegion(ctx, "", "SEL")
	if err == nil {
		t.Error("빈 country_code에 에러가 반환되어야 합니다")
	}

	// 존재하지 않는 국가 코드 → 빈 결과
	admins, err = svc.ListAdminsByRegion(ctx, "JP", "")
	if err != nil {
		t.Fatalf("ListAdminsByRegion(JP, '') 실패: %v", err)
	}
	if len(admins) != 0 {
		t.Errorf("JP 관리자 수: got %d, want 0", len(admins))
	}
}

// ============================================================================
// 확장 감사 로그 (AuditLogStore) 통합 테스트
// ============================================================================

// fakeAuditLogStore는 테스트용 확장 감사 로그 저장소입니다.
type fakeAuditLogStore struct {
	logs []*service.AuditLogDetail
}

func newFakeAuditLogStore() *fakeAuditLogStore {
	return &fakeAuditLogStore{logs: make([]*service.AuditLogDetail, 0)}
}

func (s *fakeAuditLogStore) Create(_ context.Context, log *service.AuditLogDetail) error {
	if log != nil {
		s.logs = append(s.logs, log)
	}
	return nil
}

func (s *fakeAuditLogStore) ListByAdmin(_ context.Context, adminID string, limit, offset int) ([]*service.AuditLogDetail, error) {
	var result []*service.AuditLogDetail
	for _, l := range s.logs {
		if l.AdminID == adminID {
			result = append(result, l)
		}
	}
	return result, nil
}

func (s *fakeAuditLogStore) ListByAction(_ context.Context, action string, limit, offset int) ([]*service.AuditLogDetail, error) {
	var result []*service.AuditLogDetail
	for _, l := range s.logs {
		if l.Action == action {
			result = append(result, l)
		}
	}
	return result, nil
}

func (s *fakeAuditLogStore) ListAll(_ context.Context, limit, offset int) ([]*service.AuditLogDetail, error) {
	return s.logs, nil
}

func (s *fakeAuditLogStore) Count(_ context.Context) (int, error) {
	return len(s.logs), nil
}

func TestAuditLogStore_ConfigChangeTracking(t *testing.T) {
	svc, _, _, _, _ := newTestService()
	store := newFakeAuditLogStore()
	svc.SetAuditLogStore(store)
	ctx := context.Background()

	// 1. 새 설정 생성 → config_create
	_, err := svc.SetSystemConfig(ctx, "test.key", "initial_value", "테스트", "admin-1")
	if err != nil {
		t.Fatalf("SetSystemConfig 실패: %v", err)
	}

	if len(store.logs) < 1 {
		t.Fatalf("확장 감사 로그가 기록되어야 합니다: got %d", len(store.logs))
	}

	firstLog := store.logs[0]
	if firstLog.Action != "config_create" {
		t.Errorf("새 설정은 config_create여야 합니다: got %s", firstLog.Action)
	}
	if firstLog.Resource != "config:test.key" {
		t.Errorf("Resource: got %s, want config:test.key", firstLog.Resource)
	}
	if firstLog.OldValue != "" {
		t.Errorf("새 설정의 OldValue는 비어 있어야 합니다: got %s", firstLog.OldValue)
	}
	if firstLog.NewValue != "initial_value" {
		t.Errorf("NewValue: got %s, want initial_value", firstLog.NewValue)
	}

	// 2. 설정 업데이트 → config_update
	_, err = svc.SetSystemConfig(ctx, "test.key", "updated_value", "업데이트", "admin-1")
	if err != nil {
		t.Fatalf("SetSystemConfig 업데이트 실패: %v", err)
	}

	if len(store.logs) < 2 {
		t.Fatalf("두 번째 감사 로그가 필요합니다: got %d", len(store.logs))
	}

	secondLog := store.logs[1]
	if secondLog.Action != "config_update" {
		t.Errorf("기존 설정 변경은 config_update여야 합니다: got %s", secondLog.Action)
	}
	if secondLog.OldValue != "initial_value" {
		t.Errorf("OldValue: got %s, want initial_value", secondLog.OldValue)
	}
	if secondLog.NewValue != "updated_value" {
		t.Errorf("NewValue: got %s, want updated_value", secondLog.NewValue)
	}
}

func TestAuditLogStore_AdminCreateLogging(t *testing.T) {
	svc, _, _, _, _ := newTestService()
	store := newFakeAuditLogStore()
	svc.SetAuditLogStore(store)
	ctx := context.Background()

	_, err := svc.CreateAdmin(ctx, "user-audit-1", "audit@test.com", "감사 테스트", service.RoleAdmin, "creator-1")
	if err != nil {
		t.Fatalf("CreateAdmin 실패: %v", err)
	}

	if len(store.logs) != 1 {
		t.Fatalf("감사 로그 수: got %d, want 1", len(store.logs))
	}

	log := store.logs[0]
	if log.Action != "admin_create" {
		t.Errorf("Action: got %s, want admin_create", log.Action)
	}
	if log.AdminID != "creator-1" {
		t.Errorf("AdminID: got %s, want creator-1", log.AdminID)
	}
	if log.NewValue != "audit@test.com" {
		t.Errorf("NewValue: got %s, want audit@test.com", log.NewValue)
	}
}

func TestAuditLogStore_RoleChangeLogging(t *testing.T) {
	svc, _, _, _, _ := newTestService()
	store := newFakeAuditLogStore()
	svc.SetAuditLogStore(store)
	ctx := context.Background()

	created, _ := svc.CreateAdmin(ctx, "user-role-1", "role@test.com", "역할 변경 테스트", service.RoleSupport, "system")
	store.logs = nil // 생성 로그 초기화

	_, err := svc.UpdateAdminRole(ctx, created.AdminID, service.RoleAdmin, "updater-1")
	if err != nil {
		t.Fatalf("UpdateAdminRole 실패: %v", err)
	}

	if len(store.logs) != 1 {
		t.Fatalf("감사 로그 수: got %d, want 1", len(store.logs))
	}

	log := store.logs[0]
	if log.Action != "role_change" {
		t.Errorf("Action: got %s, want role_change", log.Action)
	}
	if log.OldValue != "support" {
		t.Errorf("OldValue: got %s, want support", log.OldValue)
	}
	if log.NewValue != "admin" {
		t.Errorf("NewValue: got %s, want admin", log.NewValue)
	}
}

func TestAuditLogStore_DeactivateLogging(t *testing.T) {
	svc, _, _, _, _ := newTestService()
	store := newFakeAuditLogStore()
	svc.SetAuditLogStore(store)
	ctx := context.Background()

	created, _ := svc.CreateAdmin(ctx, "user-deact-1", "deact@test.com", "비활성화 테스트", service.RoleAnalyst, "system")
	store.logs = nil

	_, _, err := svc.DeactivateAdmin(ctx, created.AdminID, "deactivator-1", "퇴사")
	if err != nil {
		t.Fatalf("DeactivateAdmin 실패: %v", err)
	}

	if len(store.logs) != 1 {
		t.Fatalf("감사 로그 수: got %d, want 1", len(store.logs))
	}

	log := store.logs[0]
	if log.Action != "admin_deactivate" {
		t.Errorf("Action: got %s, want admin_deactivate", log.Action)
	}
	if log.OldValue != "active" {
		t.Errorf("OldValue: got %s, want active", log.OldValue)
	}
	if log.NewValue != "inactive" {
		t.Errorf("NewValue: got %s, want inactive", log.NewValue)
	}
}

func TestAuditLogStore_QueryMethods(t *testing.T) {
	svc, _, _, _, _ := newTestService()
	store := newFakeAuditLogStore()
	svc.SetAuditLogStore(store)
	ctx := context.Background()

	// 여러 액션 수행
	_, _ = svc.CreateAdmin(ctx, "user-q1", "q1@test.com", "쿼리 테스트 1", service.RoleAdmin, "admin-A")
	_, _ = svc.CreateAdmin(ctx, "user-q2", "q2@test.com", "쿼리 테스트 2", service.RoleSupport, "admin-B")
	_, _ = svc.SetSystemConfig(ctx, "query.test.key", "value1", "쿼리 설정", "admin-A")

	// 전체 조회
	details, total, err := svc.GetAuditLogDetails(ctx, 100, 0)
	if err != nil {
		t.Fatalf("GetAuditLogDetails 실패: %v", err)
	}
	if total < 3 {
		t.Errorf("전체 감사 로그 수: got %d, want >= 3", total)
	}
	if len(details) < 3 {
		t.Errorf("반환된 감사 로그 수: got %d, want >= 3", len(details))
	}

	// 관리자별 조회
	adminALogs, err := svc.GetAuditLogDetailsByAdmin(ctx, "admin-A", 100, 0)
	if err != nil {
		t.Fatalf("GetAuditLogDetailsByAdmin 실패: %v", err)
	}
	if len(adminALogs) != 2 {
		t.Errorf("admin-A 로그 수: got %d, want 2", len(adminALogs))
	}

	// 액션별 조회
	createLogs, err := svc.GetAuditLogDetailsByAction(ctx, "admin_create", 100, 0)
	if err != nil {
		t.Fatalf("GetAuditLogDetailsByAction 실패: %v", err)
	}
	if len(createLogs) != 2 {
		t.Errorf("admin_create 로그 수: got %d, want 2", len(createLogs))
	}
}

// ============================================================================
// 기존 E2E 테스트
// ============================================================================

func TestEndToEnd_AdminFlow(t *testing.T) {
	svc, _, auditRepo, _, _ := newTestService()
	ctx := context.Background()

	// 1. 관리자 생성
	admin, err := svc.CreateAdmin(ctx, "user-e2e", "e2e@test.com", "E2E 테스트", service.RoleSupport, "system")
	if err != nil {
		t.Fatalf("E2E: 관리자 생성 실패: %v", err)
	}
	if admin.Role != service.RoleSupport {
		t.Errorf("E2E: 초기 역할: got %d, want %d", admin.Role, service.RoleSupport)
	}

	// 2. 관리자 조회
	fetched, err := svc.GetAdmin(ctx, admin.AdminID)
	if err != nil {
		t.Fatalf("E2E: 관리자 조회 실패: %v", err)
	}
	if fetched.Email != "e2e@test.com" {
		t.Errorf("E2E: Email: got %s, want e2e@test.com", fetched.Email)
	}

	// 3. 역할 변경
	updated, err := svc.UpdateAdminRole(ctx, admin.AdminID, service.RoleAdmin, "system")
	if err != nil {
		t.Fatalf("E2E: 역할 변경 실패: %v", err)
	}
	if updated.Role != service.RoleAdmin {
		t.Errorf("E2E: 변경 후 역할: got %d, want %d", updated.Role, service.RoleAdmin)
	}

	// 4. 관리자 비활성화
	success, _, err := svc.DeactivateAdmin(ctx, admin.AdminID, "system", "E2E 테스트")
	if err != nil {
		t.Fatalf("E2E: 비활성화 실패: %v", err)
	}
	if !success {
		t.Error("E2E: 비활성화 실패")
	}

	// 5. 비활성화 확인
	deactivated, err := svc.GetAdmin(ctx, admin.AdminID)
	if err != nil {
		t.Fatalf("E2E: 비활성화 후 조회 실패: %v", err)
	}
	if deactivated.IsActive {
		t.Error("E2E: 비활성화 후 IsActive=true")
	}

	// 6. 설정 변경 + 조회
	_, err = svc.SetSystemConfig(ctx, "e2e_config", "e2e_value", "E2E 테스트 설정", "system")
	if err != nil {
		t.Fatalf("E2E: 설정 저장 실패: %v", err)
	}
	cfg, err := svc.GetSystemConfig(ctx, "e2e_config")
	if err != nil {
		t.Fatalf("E2E: 설정 조회 실패: %v", err)
	}
	if cfg.Value != "e2e_value" {
		t.Errorf("E2E: 설정 값: got %s, want e2e_value", cfg.Value)
	}

	// 7. 감사 로그 확인 (여러 액션이 기록되어야 함)
	entries, total, err := svc.GetAuditLog(ctx, "", service.ActionUnknown, nil, nil, 100, 0)
	if err != nil {
		t.Fatalf("E2E: 감사 로그 조회 실패: %v", err)
	}
	if total < 3 {
		t.Errorf("E2E: 감사 로그 수: got %d, want >= 3", total)
	}

	// 감사 로그에 다양한 액션이 있는지 확인
	actionSet := make(map[service.AuditAction]bool)
	for _, e := range entries {
		actionSet[e.Action] = true
	}
	if !actionSet[service.ActionCreate] {
		t.Error("E2E: 감사 로그에 ActionCreate가 없습니다")
	}
	if !actionSet[service.ActionRoleChange] {
		t.Error("E2E: 감사 로그에 ActionRoleChange가 없습니다")
	}

	// auditRepo 직접 확인
	if len(auditRepo.entries) < 3 {
		t.Errorf("E2E: auditRepo entries: got %d, want >= 3", len(auditRepo.entries))
	}

	// 8. 시스템 통계 확인
	stats, err := svc.GetSystemStats(ctx, 7)
	if err != nil {
		t.Fatalf("E2E: 시스템 통계 조회 실패: %v", err)
	}
	if stats.TotalUsers < 1 {
		t.Errorf("E2E: TotalUsers: got %d, want >= 1", stats.TotalUsers)
	}

	// 9. 사용자 목록 확인
	users, userTotal, err := svc.ListUsers(ctx, "", "", false, 10, 0)
	if err != nil {
		t.Fatalf("E2E: 사용자 목록 조회 실패: %v", err)
	}
	if userTotal < 1 {
		t.Errorf("E2E: 사용자 수: got %d, want >= 1", userTotal)
	}
	if len(users) < 1 {
		t.Errorf("E2E: 반환된 사용자 수: got %d, want >= 1", len(users))
	}
}
