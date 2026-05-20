package postgres

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/concept-service/internal/service"
)

type ConceptRepository struct {
	pool *pgxpool.Pool
}

func NewConceptRepository(pool *pgxpool.Pool) *ConceptRepository {
	return &ConceptRepository{pool: pool}
}

func (r *ConceptRepository) List(ctx context.Context) ([]*service.Concept, error) {
	const q = `SELECT id, name, description, category, icon_url, owner_id, device_ids, created_at FROM concepts ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, q)
	if err != nil { return nil, err }
	defer rows.Close()
	var list []*service.Concept
	for rows.Next() {
		c, err := scanConcept(rows)
		if err != nil { return nil, err }
		list = append(list, c)
	}
	return list, rows.Err()
}

func (r *ConceptRepository) GetByID(ctx context.Context, id string) (*service.Concept, error) {
	const q = `SELECT id, name, description, category, icon_url, owner_id, device_ids, created_at FROM concepts WHERE id = $1`
	var c service.Concept
	var deviceJSON []byte
	err := r.pool.QueryRow(ctx, q, id).Scan(&c.ID, &c.Name, &c.Description, &c.Category, &c.IconURL, &c.OwnerID, &deviceJSON, &c.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	if len(deviceJSON) > 0 { json.Unmarshal(deviceJSON, &c.DeviceIDs) }
	return &c, nil
}

func (r *ConceptRepository) Create(ctx context.Context, c *service.Concept) error {
	deviceJSON, _ := json.Marshal(c.DeviceIDs)
	const q = `INSERT INTO concepts (id, name, description, category, icon_url, owner_id, device_ids, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.pool.Exec(ctx, q, c.ID, c.Name, c.Description, c.Category, c.IconURL, c.OwnerID, deviceJSON, c.CreatedAt)
	return err
}

func (r *ConceptRepository) Update(ctx context.Context, c *service.Concept) error {
	deviceJSON, _ := json.Marshal(c.DeviceIDs)
	const q = `UPDATE concepts SET name=$1, description=$2, category=$3, icon_url=$4, device_ids=$5 WHERE id = $6`
	_, err := r.pool.Exec(ctx, q, c.Name, c.Description, c.Category, c.IconURL, deviceJSON, c.ID)
	return err
}

func (r *ConceptRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM concepts WHERE id = $1`, id)
	return err
}

func scanConcept(rows pgx.Rows) (*service.Concept, error) {
	var c service.Concept
	var deviceJSON []byte
	if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.Category, &c.IconURL, &c.OwnerID, &deviceJSON, &c.CreatedAt); err != nil {
		return nil, err
	}
	if len(deviceJSON) > 0 { json.Unmarshal(deviceJSON, &c.DeviceIDs) }
	return &c, nil
}

type OrganizationRepository struct {
	pool *pgxpool.Pool
}

func NewOrganizationRepository(pool *pgxpool.Pool) *OrganizationRepository {
	return &OrganizationRepository{pool: pool}
}

func (r *OrganizationRepository) List(ctx context.Context) ([]*service.Organization, error) {
	const q = `SELECT id, name, description, owner_id, member_ids, created_at FROM organizations ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, q)
	if err != nil { return nil, err }
	defer rows.Close()
	var list []*service.Organization
	for rows.Next() {
		var o service.Organization
		var memberJSON []byte
		if err := rows.Scan(&o.ID, &o.Name, &o.Description, &o.OwnerID, &memberJSON, &o.CreatedAt); err != nil { return nil, err }
		if len(memberJSON) > 0 { json.Unmarshal(memberJSON, &o.MemberIDs) }
		list = append(list, &o)
	}
	return list, rows.Err()
}

func (r *OrganizationRepository) GetByID(ctx context.Context, id string) (*service.Organization, error) {
	const q = `SELECT id, name, description, owner_id, member_ids, created_at FROM organizations WHERE id = $1`
	var o service.Organization
	var memberJSON []byte
	err := r.pool.QueryRow(ctx, q, id).Scan(&o.ID, &o.Name, &o.Description, &o.OwnerID, &memberJSON, &o.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	if len(memberJSON) > 0 { json.Unmarshal(memberJSON, &o.MemberIDs) }
	return &o, nil
}

func (r *OrganizationRepository) Create(ctx context.Context, o *service.Organization) error {
	memberJSON, _ := json.Marshal(o.MemberIDs)
	const q = `INSERT INTO organizations (id, name, description, owner_id, member_ids, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.pool.Exec(ctx, q, o.ID, o.Name, o.Description, o.OwnerID, memberJSON, o.CreatedAt)
	return err
}

func (r *OrganizationRepository) Update(ctx context.Context, o *service.Organization) error {
	memberJSON, _ := json.Marshal(o.MemberIDs)
	const q = `UPDATE organizations SET name=$1, description=$2, member_ids=$3 WHERE id = $4`
	_, err := r.pool.Exec(ctx, q, o.Name, o.Description, memberJSON, o.ID)
	return err
}

func (r *OrganizationRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM organizations WHERE id = $1`, id)
	return err
}
