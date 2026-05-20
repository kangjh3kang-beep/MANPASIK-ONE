package postgres

import (
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/emergency-service/internal/service"
)

type EmergencyRepository struct {
	pool *pgxpool.Pool
}

func NewEmergencyRepository(pool *pgxpool.Pool) *EmergencyRepository {
	return &EmergencyRepository{pool: pool}
}

func (r *EmergencyRepository) CreateEmergency(e *service.Emergency) (string, error) {
	contactsJSON, _ := json.Marshal(e.ContactIDs)
	const q = `INSERT INTO emergencies (id, user_id, type, location, description, status, severity, contact_ids, resolution, created_at, resolved_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.pool.Exec(nil, q, e.ID, e.UserID, e.Type, e.Location, e.Description, e.Status,
		int(e.Severity), contactsJSON, e.Resolution, e.Timestamp, e.ResolvedAt)
	if err != nil {
		return "", err
	}
	return e.ID, nil
}

func (r *EmergencyRepository) GetEmergency(id string) (*service.Emergency, error) {
	const q = `SELECT id, user_id, type, location, description, status, COALESCE(severity, 0), contact_ids,
		COALESCE(resolution, ''), created_at, resolved_at
		FROM emergencies WHERE id = $1`
	var e service.Emergency
	var contactsJSON []byte
	var severity int
	err := r.pool.QueryRow(nil, q, id).Scan(
		&e.ID, &e.UserID, &e.Type, &e.Location, &e.Description, &e.Status,
		&severity, &contactsJSON, &e.Resolution, &e.Timestamp, &e.ResolvedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	e.Severity = service.EmergencySeverity(severity)
	if len(contactsJSON) > 0 {
		json.Unmarshal(contactsJSON, &e.ContactIDs)
	}
	return &e, nil
}

func (r *EmergencyRepository) UpdateEmergency(e *service.Emergency) error {
	contactsJSON, _ := json.Marshal(e.ContactIDs)
	const q = `UPDATE emergencies SET status = $2, severity = $3, resolution = $4, resolved_at = $5, contact_ids = $6
		WHERE id = $1`
	_, err := r.pool.Exec(nil, q, e.ID, e.Status, int(e.Severity), e.Resolution, e.ResolvedAt, contactsJSON)
	return err
}

func (r *EmergencyRepository) ListEmergenciesByUser(userID string) ([]*service.Emergency, error) {
	const q = `SELECT id, user_id, type, location, description, status, COALESCE(severity, 0), contact_ids,
		COALESCE(resolution, ''), created_at, resolved_at
		FROM emergencies WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(nil, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*service.Emergency
	for rows.Next() {
		var e service.Emergency
		var contactsJSON []byte
		var severity int
		if err := rows.Scan(&e.ID, &e.UserID, &e.Type, &e.Location, &e.Description, &e.Status,
			&severity, &contactsJSON, &e.Resolution, &e.Timestamp, &e.ResolvedAt); err != nil {
			return nil, err
		}
		e.Severity = service.EmergencySeverity(severity)
		if len(contactsJSON) > 0 {
			json.Unmarshal(contactsJSON, &e.ContactIDs)
		}
		list = append(list, &e)
	}
	return list, rows.Err()
}

func (r *EmergencyRepository) GetContactsByUser(userID string) ([]*service.EmergencyContact, error) {
	const q = `SELECT id, user_id, name, phone, relationship FROM emergency_contacts WHERE user_id = $1`
	rows, err := r.pool.Query(nil, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*service.EmergencyContact
	for rows.Next() {
		var c service.EmergencyContact
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Phone, &c.Relationship); err != nil {
			return nil, err
		}
		list = append(list, &c)
	}
	return list, rows.Err()
}

func (r *EmergencyRepository) AddContact(c *service.EmergencyContact) error {
	const q = `INSERT INTO emergency_contacts (id, user_id, name, phone, relationship)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.pool.Exec(nil, q, c.ID, c.UserID, c.Name, c.Phone, c.Relationship)
	return err
}

func (r *EmergencyRepository) RemoveContact(userID, contactID string) error {
	const q = `DELETE FROM emergency_contacts WHERE user_id = $1 AND id = $2`
	tag, err := r.pool.Exec(nil, q, userID, contactID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return service.ErrNotFound
	}
	return nil
}

func (r *EmergencyRepository) GetSettings(userID string) (*service.EmergencySettings, error) {
	const q = `SELECT user_id, auto_call_119, emergency_contact_ids, medical_info
		FROM emergency_settings WHERE user_id = $1`
	var s service.EmergencySettings
	var contactIDsJSON []byte
	err := r.pool.QueryRow(nil, q, userID).Scan(&s.UserID, &s.AutoCall119, &contactIDsJSON, &s.MedicalInfo)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if len(contactIDsJSON) > 0 {
		json.Unmarshal(contactIDsJSON, &s.EmergencyContactIDs)
	}
	return &s, nil
}

func (r *EmergencyRepository) SaveSettings(s *service.EmergencySettings) error {
	contactIDsJSON, _ := json.Marshal(s.EmergencyContactIDs)
	const q = `INSERT INTO emergency_settings (user_id, auto_call_119, emergency_contact_ids, medical_info)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE SET
			auto_call_119 = EXCLUDED.auto_call_119,
			emergency_contact_ids = EXCLUDED.emergency_contact_ids,
			medical_info = EXCLUDED.medical_info`
	_, err := r.pool.Exec(nil, q, s.UserID, s.AutoCall119, contactIDsJSON, s.MedicalInfo)
	return err
}
