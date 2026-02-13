// Package postgres는 family-service의 PostgreSQL 저장소 구현입니다.
package postgres

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/family-service/internal/service"
)

var roleToDB = map[service.FamilyRole]string{
	service.RoleOwner:    "owner",
	service.RoleGuardian: "guardian",
	service.RoleMember:   "member",
	service.RoleChild:    "child",
	service.RoleElderly:   "elderly",
	service.RoleUnknown:  "member",
}

var dbToRole = map[string]service.FamilyRole{
	"owner":    service.RoleOwner,
	"guardian": service.RoleGuardian,
	"member":   service.RoleMember,
	"child":    service.RoleChild,
	"elderly":  service.RoleElderly,
}

var inviteStatusToDB = map[service.InvitationStatus]string{
	service.InvitePending:  "pending",
	service.InviteAccepted: "accepted",
	service.InviteDeclined: "declined",
	service.InviteExpired:   "expired",
	service.InviteUnknown:  "pending",
}

var dbToInviteStatus = map[string]service.InvitationStatus{
	"pending":  service.InvitePending,
	"accepted": service.InviteAccepted,
	"declined": service.InviteDeclined,
	"expired":  service.InviteExpired,
}

// ============================================================================
// GroupRepository (FamilyGroupRepository)
// ============================================================================

// GroupRepository는 PostgreSQL 기반 FamilyGroupRepository 구현입니다.
type GroupRepository struct {
	pool *pgxpool.Pool
}

// NewGroupRepository는 GroupRepository를 생성합니다.
func NewGroupRepository(pool *pgxpool.Pool) *GroupRepository {
	return &GroupRepository{pool: pool}
}

// Save는 가족 그룹을 저장합니다.
func (r *GroupRepository) Save(ctx context.Context, g *service.FamilyGroup) error {
	const q = `INSERT INTO family_groups (id, owner_user_id, group_name, description, max_members, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $6)`
	_, err := r.pool.Exec(ctx, q,
		g.ID, g.OwnerUserID, g.GroupName, nullIfEmpty(g.Description), g.MaxMembers, g.CreatedAt,
	)
	return err
}

// FindByID는 ID로 가족 그룹을 조회합니다.
func (r *GroupRepository) FindByID(ctx context.Context, id string) (*service.FamilyGroup, error) {
	const q = `SELECT id, owner_user_id, group_name, COALESCE(description,''), max_members, created_at
		FROM family_groups WHERE id = $1`
	var g service.FamilyGroup
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&g.ID, &g.OwnerUserID, &g.GroupName, &g.Description, &g.MaxMembers, &g.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &g, nil
}

// ============================================================================
// MemberRepository (FamilyMemberRepository)
// ============================================================================

// MemberRepository는 PostgreSQL 기반 FamilyMemberRepository 구현입니다.
type MemberRepository struct {
	pool *pgxpool.Pool
}

// NewMemberRepository는 MemberRepository를 생성합니다.
func NewMemberRepository(pool *pgxpool.Pool) *MemberRepository {
	return &MemberRepository{pool: pool}
}

// Save는 가족 멤버를 저장합니다.
func (r *MemberRepository) Save(ctx context.Context, m *service.FamilyMember) error {
	roleStr, _ := roleToDB[m.Role]
	if roleStr == "" {
		roleStr = "member"
	}
	const q = `INSERT INTO family_members (user_id, group_id, display_name, email, role, sharing_enabled, joined_at)
		VALUES ($1, $2, $3, $4, $5::family_role, $6, $7)
		ON CONFLICT (user_id, group_id) DO UPDATE SET display_name = EXCLUDED.display_name, email = EXCLUDED.email, role = EXCLUDED.role, sharing_enabled = EXCLUDED.sharing_enabled`
	_, err := r.pool.Exec(ctx, q,
		m.UserID, m.GroupID, nullIfEmpty(m.DisplayName), nullIfEmpty(m.Email), roleStr, m.SharingEnabled, m.JoinedAt,
	)
	return err
}

// FindByGroupID는 그룹의 멤버 목록을 조회합니다.
func (r *MemberRepository) FindByGroupID(ctx context.Context, groupID string) ([]*service.FamilyMember, error) {
	const q = `SELECT user_id, group_id, COALESCE(display_name,''), COALESCE(email,''), role::text, sharing_enabled, joined_at
		FROM family_members WHERE group_id = $1 ORDER BY joined_at ASC`
	rows, err := r.pool.Query(ctx, q, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*service.FamilyMember
	for rows.Next() {
		var m service.FamilyMember
		var roleStr string
		if err := rows.Scan(&m.UserID, &m.GroupID, &m.DisplayName, &m.Email, &roleStr, &m.SharingEnabled, &m.JoinedAt); err != nil {
			return nil, err
		}
		m.Role = dbToRole[roleStr]
		if m.Role == 0 {
			m.Role = service.RoleMember
		}
		list = append(list, &m)
	}
	return list, rows.Err()
}

// FindByUserIDAndGroupID는 사용자와 그룹으로 멤버를 조회합니다.
func (r *MemberRepository) FindByUserIDAndGroupID(ctx context.Context, userID, groupID string) (*service.FamilyMember, error) {
	const q = `SELECT user_id, group_id, COALESCE(display_name,''), COALESCE(email,''), role::text, sharing_enabled, joined_at
		FROM family_members WHERE user_id = $1 AND group_id = $2`
	var m service.FamilyMember
	var roleStr string
	err := r.pool.QueryRow(ctx, q, userID, groupID).Scan(
		&m.UserID, &m.GroupID, &m.DisplayName, &m.Email, &roleStr, &m.SharingEnabled, &m.JoinedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	m.Role = dbToRole[roleStr]
	if m.Role == 0 {
		m.Role = service.RoleMember
	}
	return &m, nil
}

// Remove는 가족 멤버를 제거합니다.
func (r *MemberRepository) Remove(ctx context.Context, userID, groupID string) error {
	const q = `DELETE FROM family_members WHERE user_id = $1 AND group_id = $2`
	_, err := r.pool.Exec(ctx, q, userID, groupID)
	return err
}

// CountByGroupID는 그룹의 멤버 수를 반환합니다.
func (r *MemberRepository) CountByGroupID(ctx context.Context, groupID string) (int, error) {
	const q = `SELECT COUNT(*) FROM family_members WHERE group_id = $1`
	var count int
	err := r.pool.QueryRow(ctx, q, groupID).Scan(&count)
	return count, err
}

// ============================================================================
// InvitationRepository
// ============================================================================

// InvitationRepository는 PostgreSQL 기반 InvitationRepository 구현입니다.
type InvitationRepository struct {
	pool *pgxpool.Pool
}

// NewInvitationRepository는 InvitationRepository를 생성합니다.
func NewInvitationRepository(pool *pgxpool.Pool) *InvitationRepository {
	return &InvitationRepository{pool: pool}
}

// Save는 초대를 저장합니다.
func (r *InvitationRepository) Save(ctx context.Context, inv *service.FamilyInvitation) error {
	roleStr, _ := roleToDB[inv.Role]
	if roleStr == "" {
		roleStr = "member"
	}
	statusStr, _ := inviteStatusToDB[inv.Status]
	if statusStr == "" {
		statusStr = "pending"
	}
	const q = `INSERT INTO family_invitations (id, group_id, inviter_user_id, invitee_email, role, message, status, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5::family_role, $6, $7::invitation_status, $8, $9)`
	_, err := r.pool.Exec(ctx, q,
		inv.ID, inv.GroupID, inv.InviterUserID, inv.InviteeEmail, roleStr, nullIfEmpty(inv.Message), statusStr, inv.CreatedAt, inv.ExpiresAt,
	)
	return err
}

// FindByID는 ID로 초대를 조회합니다.
func (r *InvitationRepository) FindByID(ctx context.Context, id string) (*service.FamilyInvitation, error) {
	const q = `SELECT inv.id, inv.group_id, fg.group_name, inv.inviter_user_id, inv.invitee_email, inv.role::text, COALESCE(inv.message,''), inv.status::text, inv.created_at, inv.expires_at
		FROM family_invitations inv
		JOIN family_groups fg ON inv.group_id = fg.id
		WHERE inv.id = $1`
	var inv service.FamilyInvitation
	var roleStr, statusStr string
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&inv.ID, &inv.GroupID, &inv.GroupName, &inv.InviterUserID, &inv.InviteeEmail,
		&roleStr, &inv.Message, &statusStr, &inv.CreatedAt, &inv.ExpiresAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	inv.Role = dbToRole[roleStr]
	if inv.Role == 0 {
		inv.Role = service.RoleMember
	}
	inv.Status = dbToInviteStatus[statusStr]
	if inv.Status == 0 {
		inv.Status = service.InvitePending
	}
	return &inv, nil
}

// Update는 초대를 업데이트합니다.
func (r *InvitationRepository) Update(ctx context.Context, inv *service.FamilyInvitation) error {
	statusStr, _ := inviteStatusToDB[inv.Status]
	if statusStr == "" {
		statusStr = "pending"
	}
	const q = `UPDATE family_invitations SET status = $1::invitation_status WHERE id = $2`
	_, err := r.pool.Exec(ctx, q, statusStr, inv.ID)
	return err
}

// ============================================================================
// SharingPreferencesRepository
// ============================================================================

// SharingPreferencesRepository는 PostgreSQL 기반 SharingPreferencesRepository 구현입니다.
type SharingPreferencesRepository struct {
	pool *pgxpool.Pool
}

// NewSharingPreferencesRepository는 SharingPreferencesRepository를 생성합니다.
func NewSharingPreferencesRepository(pool *pgxpool.Pool) *SharingPreferencesRepository {
	return &SharingPreferencesRepository{pool: pool}
}

// Save는 공유 설정을 저장합니다 (UPSERT).
func (r *SharingPreferencesRepository) Save(ctx context.Context, pref *service.SharingPreferences) error {
	// Base columns; extended columns (measurement_days_limit, allowed_biomarkers, require_approval) require 13a-family-sharing-extended.sql
	const q = `INSERT INTO sharing_preferences (user_id, group_id, share_measurements, share_health_score, share_goals, share_coaching, share_alerts, allowed_viewer_ids, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8::text[]::uuid[], NOW(), NOW())
		ON CONFLICT (user_id, group_id) DO UPDATE SET
			share_measurements = EXCLUDED.share_measurements,
			share_health_score = EXCLUDED.share_health_score,
			share_goals = EXCLUDED.share_goals,
			share_coaching = EXCLUDED.share_coaching,
			share_alerts = EXCLUDED.share_alerts,
			allowed_viewer_ids = EXCLUDED.allowed_viewer_ids,
			updated_at = NOW()`
	_, err := r.pool.Exec(ctx, q,
		pref.UserID, pref.GroupID,
		pref.ShareMeasurements, pref.ShareHealthScore, pref.ShareGoals, pref.ShareCoaching, pref.ShareAlerts,
		pref.AllowedViewerIDs,
	)
	return err
}

// FindByUserIDAndGroupID는 사용자와 그룹으로 공유 설정을 조회합니다.
func (r *SharingPreferencesRepository) FindByUserIDAndGroupID(ctx context.Context, userID, groupID string) (*service.SharingPreferences, error) {
	// Base columns; allowed_viewer_ids scanned via array_to_string
	const q = `SELECT user_id, group_id, share_measurements, share_health_score, share_goals, share_coaching, share_alerts,
		COALESCE(array_to_string(allowed_viewer_ids, ','), ''), 0, '{}'::text[], FALSE
		FROM sharing_preferences WHERE user_id = $1 AND group_id = $2`
	var pref service.SharingPreferences
	var allowedViewersStr string
	var measurementDaysLimit int
	var allowedBiomarkers []string
	var requireApproval bool
	err := r.pool.QueryRow(ctx, q, userID, groupID).Scan(
		&pref.UserID, &pref.GroupID,
		&pref.ShareMeasurements, &pref.ShareHealthScore, &pref.ShareGoals, &pref.ShareCoaching, &pref.ShareAlerts,
		&allowedViewersStr, &measurementDaysLimit, &allowedBiomarkers, &requireApproval,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	// Parse "uuid1,uuid2" to []string
	if allowedViewersStr != "" {
		pref.AllowedViewerIDs = splitAndTrim(allowedViewersStr, ",")
	}
	pref.MeasurementDaysLimit = measurementDaysLimit
	pref.AllowedBiomarkers = allowedBiomarkers
	pref.RequireApproval = requireApproval
	return &pref, nil
}

func splitAndTrim(s, sep string) []string {
	if s == "" {
		return nil
	}
	var result []string
	for _, p := range strings.Split(s, sep) {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func nullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
