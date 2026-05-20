package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/audit-service/internal/service"
)

type AuditRepository struct {
	pool *pgxpool.Pool
}

func NewAuditRepository(pool *pgxpool.Pool) *AuditRepository {
	return &AuditRepository{pool: pool}
}

func (r *AuditRepository) Store(ctx context.Context, entry *service.AuditEntry) error {
	metaJSON, _ := json.Marshal(entry.Metadata)
	const q = `INSERT INTO audit_entries (id, admin_id, action, resource_type, resource_id, description, ip_address, user_agent, metadata, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.pool.Exec(ctx, q,
		entry.ID, entry.AdminID, entry.Action, entry.ResourceType, entry.ResourceID,
		entry.Description, entry.IPAddress, entry.UserAgent, metaJSON, entry.Timestamp,
	)
	return err
}

func (r *AuditRepository) List(ctx context.Context, filter service.AuditFilter) ([]*service.AuditEntry, int, error) {
	where := "WHERE 1=1"
	args := []interface{}{}
	idx := 1

	if filter.AdminID != "" {
		where += fmt.Sprintf(" AND admin_id = $%d", idx)
		args = append(args, filter.AdminID)
		idx++
	}
	if filter.Action != "" {
		where += fmt.Sprintf(" AND action = $%d", idx)
		args = append(args, filter.Action)
		idx++
	}
	if filter.ResourceType != "" {
		where += fmt.Sprintf(" AND resource_type = $%d", idx)
		args = append(args, filter.ResourceType)
		idx++
	}
	if !filter.StartTime.IsZero() {
		where += fmt.Sprintf(" AND timestamp >= $%d", idx)
		args = append(args, filter.StartTime)
		idx++
	}
	if !filter.EndTime.IsZero() {
		where += fmt.Sprintf(" AND timestamp <= $%d", idx)
		args = append(args, filter.EndTime)
		idx++
	}

	// count
	countQ := "SELECT COUNT(*) FROM audit_entries " + where
	var total int
	if err := r.pool.QueryRow(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// select
	selQ := fmt.Sprintf(`SELECT id, admin_id, action, resource_type, resource_id, description, ip_address, user_agent, metadata, timestamp
		FROM audit_entries %s ORDER BY timestamp DESC LIMIT $%d OFFSET $%d`, where, idx, idx+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.pool.Query(ctx, selQ, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*service.AuditEntry
	for rows.Next() {
		var e service.AuditEntry
		var metaJSON []byte
		if err := rows.Scan(&e.ID, &e.AdminID, &e.Action, &e.ResourceType, &e.ResourceID,
			&e.Description, &e.IPAddress, &e.UserAgent, &metaJSON, &e.Timestamp); err != nil {
			return nil, 0, err
		}
		if len(metaJSON) > 0 {
			json.Unmarshal(metaJSON, &e.Metadata)
		}
		list = append(list, &e)
	}
	return list, total, rows.Err()
}

func (r *AuditRepository) Get(ctx context.Context, entryID string) (*service.AuditEntry, error) {
	const q = `SELECT id, admin_id, action, resource_type, resource_id, description, ip_address, user_agent, metadata, timestamp
		FROM audit_entries WHERE id = $1`
	var e service.AuditEntry
	var metaJSON []byte
	err := r.pool.QueryRow(ctx, q, entryID).Scan(
		&e.ID, &e.AdminID, &e.Action, &e.ResourceType, &e.ResourceID,
		&e.Description, &e.IPAddress, &e.UserAgent, &metaJSON, &e.Timestamp,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if len(metaJSON) > 0 {
		json.Unmarshal(metaJSON, &e.Metadata)
	}
	return &e, nil
}
