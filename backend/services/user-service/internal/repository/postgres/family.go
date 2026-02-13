package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/user-service/internal/service"
)

// FamilyRepository는 PostgreSQL 기반 FamilyRepository 구현입니다.
type FamilyRepository struct {
	pool *pgxpool.Pool
}

// NewFamilyRepository는 FamilyRepository를 생성합니다.
func NewFamilyRepository(pool *pgxpool.Pool) *FamilyRepository {
	return &FamilyRepository{pool: pool}
}

// GetGroup은 그룹 ID로 가족 그룹을 조회합니다 (멤버 포함).
func (r *FamilyRepository) GetGroup(ctx context.Context, groupID string) (*service.FamilyGroup, error) {
	const groupQ = `SELECT id, name, owner_id FROM family_groups WHERE id = $1`
	var g service.FamilyGroup
	err := r.pool.QueryRow(ctx, groupQ, groupID).Scan(&g.ID, &g.Name, &g.OwnerID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	members, err := r.getMembers(ctx, groupID)
	if err != nil {
		return nil, err
	}
	g.Members = members
	return &g, nil
}

// GetUserGroups는 사용자가 속한 모든 가족 그룹을 조회합니다.
func (r *FamilyRepository) GetUserGroups(ctx context.Context, userID string) ([]*service.FamilyGroup, error) {
	const q = `SELECT DISTINCT fg.id, fg.name, fg.owner_id
		FROM family_groups fg
		LEFT JOIN family_members fm ON fg.id = fm.group_id
		WHERE fg.owner_id = $1 OR fm.user_id = $1`
	rows, err := r.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*service.FamilyGroup
	for rows.Next() {
		var g service.FamilyGroup
		if err := rows.Scan(&g.ID, &g.Name, &g.OwnerID); err != nil {
			return nil, err
		}
		members, err := r.getMembers(ctx, g.ID)
		if err != nil {
			return nil, err
		}
		g.Members = members
		groups = append(groups, &g)
	}
	return groups, rows.Err()
}

// CreateGroup은 새 가족 그룹을 생성합니다.
func (r *FamilyRepository) CreateGroup(ctx context.Context, group *service.FamilyGroup) error {
	const q = `INSERT INTO family_groups (id, name, owner_id) VALUES ($1, $2, $3)`
	_, err := r.pool.Exec(ctx, q, group.ID, group.Name, group.OwnerID)
	return err
}

// AddMember는 가족 그룹에 구성원을 추가합니다.
func (r *FamilyRepository) AddMember(ctx context.Context, groupID, userID, role string) error {
	// 그룹 존재 확인
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM family_groups WHERE id = $1)`, groupID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("group %s not found", groupID)
	}

	const q = `INSERT INTO family_members (group_id, user_id, role) VALUES ($1, $2, $3)
		ON CONFLICT (group_id, user_id) DO UPDATE SET role = EXCLUDED.role`
	_, err = r.pool.Exec(ctx, q, groupID, userID, role)
	return err
}

// RemoveMember는 가족 그룹에서 구성원을 제거합니다.
func (r *FamilyRepository) RemoveMember(ctx context.Context, groupID, userID string) error {
	const q = `DELETE FROM family_members WHERE group_id = $1 AND user_id = $2`
	_, err := r.pool.Exec(ctx, q, groupID, userID)
	return err
}

// getMembers는 그룹의 멤버 목록을 조회합니다.
func (r *FamilyRepository) getMembers(ctx context.Context, groupID string) ([]*service.FamilyMember, error) {
	const q = `SELECT user_id, role FROM family_members WHERE group_id = $1`
	rows, err := r.pool.Query(ctx, q, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*service.FamilyMember
	for rows.Next() {
		var m service.FamilyMember
		if err := rows.Scan(&m.UserID, &m.Role); err != nil {
			return nil, err
		}
		members = append(members, &m)
	}
	return members, rows.Err()
}
