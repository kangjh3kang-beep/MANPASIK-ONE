#!/bin/bash
set -e

BASE="/home/kangjh3kang/Manpasik/backend/services"

# 1. analytics-service
mkdir -p "$BASE/analytics-service/internal/repository/postgres"
cat > "$BASE/analytics-service/internal/repository/postgres/analytics.go" << 'GOEOF'
package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/analytics-service/internal/service"
)

type AnalyticsRepository struct {
	pool *pgxpool.Pool
}

func NewAnalyticsRepository(pool *pgxpool.Pool) *AnalyticsRepository {
	return &AnalyticsRepository{pool: pool}
}

func (r *AnalyticsRepository) TrackEvent(ctx context.Context, event *service.AnalyticsEvent) error {
	propsJSON, _ := json.Marshal(event.Properties)
	const q = `INSERT INTO analytics_events (id, user_id, event_type, properties, created_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.pool.Exec(ctx, q, event.ID, event.UserID, event.EventType, propsJSON, event.CreatedAt)
	return err
}

func (r *AnalyticsRepository) GetUserAnalytics(ctx context.Context, userID string) (*service.UserAnalytics, error) {
	const q = `SELECT COUNT(*), MAX(created_at) FROM analytics_events WHERE user_id = $1`
	var total int
	var lastAt *time.Time
	err := r.pool.QueryRow(ctx, q, userID).Scan(&total, &lastAt)
	if err != nil {
		return nil, err
	}
	ua := &service.UserAnalytics{UserID: userID, TotalEvents: total}
	if lastAt != nil {
		ua.LastEventAt = *lastAt
	}
	// top event types
	const q2 = `SELECT event_type, COUNT(*) as cnt FROM analytics_events WHERE user_id = $1
		GROUP BY event_type ORDER BY cnt DESC LIMIT 5`
	rows, err := r.pool.Query(ctx, q2, userID)
	if err != nil {
		return ua, nil
	}
	defer rows.Close()
	for rows.Next() {
		var et string
		var c int
		if err := rows.Scan(&et, &c); err == nil {
			ua.TopEventTypes = append(ua.TopEventTypes, et)
		}
	}
	return ua, nil
}

func (r *AnalyticsRepository) GetDailyStats(ctx context.Context, date string) (*service.DailyStats, error) {
	const q = `SELECT stat_date, total_events, unique_users, events_by_type FROM analytics_daily_stats WHERE stat_date = $1`
	var ds service.DailyStats
	var ebtJSON []byte
	err := r.pool.QueryRow(ctx, q, date).Scan(&ds.Date, &ds.TotalEvents, &ds.UniqueUsers, &ebtJSON)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if len(ebtJSON) > 0 {
		json.Unmarshal(ebtJSON, &ds.EventsByType)
	}
	return &ds, nil
}

func (r *AnalyticsRepository) ListRecentEvents(ctx context.Context, limit int) ([]*service.AnalyticsEvent, error) {
	const q = `SELECT id, user_id, event_type, properties, created_at FROM analytics_events
		ORDER BY created_at DESC LIMIT $1`
	rows, err := r.pool.Query(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*service.AnalyticsEvent
	for rows.Next() {
		var e service.AnalyticsEvent
		var propsJSON []byte
		if err := rows.Scan(&e.ID, &e.UserID, &e.EventType, &propsJSON, &e.CreatedAt); err != nil {
			return nil, err
		}
		if len(propsJSON) > 0 {
			json.Unmarshal(propsJSON, &e.Properties)
		}
		list = append(list, &e)
	}
	return list, rows.Err()
}
GOEOF

# 2. emergency-service
mkdir -p "$BASE/emergency-service/internal/repository/postgres"
cat > "$BASE/emergency-service/internal/repository/postgres/emergency.go" << 'GOEOF'
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
	const q = `INSERT INTO emergencies (id, user_id, type, location, description, status, contact_ids, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.pool.Exec(nil, q, e.ID, e.UserID, e.Type, e.Location, e.Description, e.Status, contactsJSON, e.Timestamp)
	if err != nil {
		return "", err
	}
	return e.ID, nil
}

func (r *EmergencyRepository) GetEmergency(id string) (*service.Emergency, error) {
	const q = `SELECT id, user_id, type, location, description, status, contact_ids, created_at
		FROM emergencies WHERE id = $1`
	var e service.Emergency
	var contactsJSON []byte
	err := r.pool.QueryRow(nil, q, id).Scan(
		&e.ID, &e.UserID, &e.Type, &e.Location, &e.Description, &e.Status, &contactsJSON, &e.Timestamp,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if len(contactsJSON) > 0 {
		json.Unmarshal(contactsJSON, &e.ContactIDs)
	}
	return &e, nil
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
GOEOF

# 3. audit-service
mkdir -p "$BASE/audit-service/internal/repository/postgres"
cat > "$BASE/audit-service/internal/repository/postgres/audit.go" << 'GOEOF'
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
GOEOF

# 4. digital-twin-service
mkdir -p "$BASE/digital-twin-service/internal/repository/postgres"
cat > "$BASE/digital-twin-service/internal/repository/postgres/twin.go" << 'GOEOF'
package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/digital-twin-service/internal/service"
)

type TwinRepository struct {
	pool *pgxpool.Pool
}

func NewTwinRepository(pool *pgxpool.Pool) *TwinRepository {
	return &TwinRepository{pool: pool}
}

func (r *TwinRepository) Store(ctx context.Context, state *service.TwinState) error {
	const q = `INSERT INTO twin_states (id, session_id, user_id, device_id, ewma_value, cusum_pos, cusum_neg,
		health_state, drift_score, remaining_measurements, measurement_count, last_synced_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (id) DO UPDATE SET
			ewma_value = EXCLUDED.ewma_value, cusum_pos = EXCLUDED.cusum_pos, cusum_neg = EXCLUDED.cusum_neg,
			health_state = EXCLUDED.health_state, drift_score = EXCLUDED.drift_score,
			remaining_measurements = EXCLUDED.remaining_measurements, measurement_count = EXCLUDED.measurement_count,
			last_synced_at = EXCLUDED.last_synced_at`
	_, err := r.pool.Exec(ctx, q,
		state.ID, state.SessionID, state.UserID, state.DeviceID,
		state.EWMAValue, state.CUSUMPos, state.CUSUMNeg,
		state.HealthState, state.DriftScore,
		state.RemainingMeasurements, state.MeasurementCount,
		state.LastSyncedAt, state.CreatedAt,
	)
	return err
}

func (r *TwinRepository) Get(ctx context.Context, sessionID string) (*service.TwinState, error) {
	const q = `SELECT id, session_id, user_id, device_id, ewma_value, cusum_pos, cusum_neg,
		health_state, drift_score, remaining_measurements, measurement_count, last_synced_at, created_at
		FROM twin_states WHERE session_id = $1 ORDER BY last_synced_at DESC LIMIT 1`
	var s service.TwinState
	err := r.pool.QueryRow(ctx, q, sessionID).Scan(
		&s.ID, &s.SessionID, &s.UserID, &s.DeviceID,
		&s.EWMAValue, &s.CUSUMPos, &s.CUSUMNeg,
		&s.HealthState, &s.DriftScore,
		&s.RemainingMeasurements, &s.MeasurementCount,
		&s.LastSyncedAt, &s.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *TwinRepository) ListByDevice(ctx context.Context, deviceID string, limit int) ([]*service.TwinState, error) {
	const q = `SELECT id, session_id, user_id, device_id, ewma_value, cusum_pos, cusum_neg,
		health_state, drift_score, remaining_measurements, measurement_count, last_synced_at, created_at
		FROM twin_states WHERE device_id = $1 ORDER BY last_synced_at DESC LIMIT $2`
	rows, err := r.pool.Query(ctx, q, deviceID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*service.TwinState
	for rows.Next() {
		var s service.TwinState
		if err := rows.Scan(
			&s.ID, &s.SessionID, &s.UserID, &s.DeviceID,
			&s.EWMAValue, &s.CUSUMPos, &s.CUSUMNeg,
			&s.HealthState, &s.DriftScore,
			&s.RemainingMeasurements, &s.MeasurementCount,
			&s.LastSyncedAt, &s.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, &s)
	}
	return list, rows.Err()
}
GOEOF

# 5. iot-gateway-service
mkdir -p "$BASE/iot-gateway-service/internal/repository/postgres"
cat > "$BASE/iot-gateway-service/internal/repository/postgres/iot_gateway.go" << 'GOEOF'
package postgres

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/iot-gateway-service/internal/service"
)

type IoTRepository struct {
	pool *pgxpool.Pool
}

func NewIoTRepository(pool *pgxpool.Pool) *IoTRepository {
	return &IoTRepository{pool: pool}
}

func (r *IoTRepository) RegisterDevice(ctx context.Context, device *service.IoTDevice) error {
	metaJSON, _ := json.Marshal(device.Metadata)
	const q = `INSERT INTO iot_devices (id, device_id, protocol, last_ping_at, status, metadata)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.pool.Exec(ctx, q, device.ID, device.DeviceID, device.Protocol, device.LastPingAt, device.Status, metaJSON)
	return err
}

func (r *IoTRepository) GetDevice(ctx context.Context, deviceID string) (*service.IoTDevice, error) {
	const q = `SELECT id, device_id, protocol, last_ping_at, status, metadata
		FROM iot_devices WHERE device_id = $1`
	var d service.IoTDevice
	var metaJSON []byte
	err := r.pool.QueryRow(ctx, q, deviceID).Scan(&d.ID, &d.DeviceID, &d.Protocol, &d.LastPingAt, &d.Status, &metaJSON)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if len(metaJSON) > 0 {
		json.Unmarshal(metaJSON, &d.Metadata)
	}
	return &d, nil
}

func (r *IoTRepository) SendCommand(ctx context.Context, cmd *service.IoTCommand) error {
	const q = `INSERT INTO iot_commands (id, device_id, command_type, payload, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.pool.Exec(ctx, q, cmd.ID, cmd.DeviceID, cmd.CommandType, cmd.Payload, cmd.Status, cmd.CreatedAt)
	return err
}

func (r *IoTRepository) GetCommand(ctx context.Context, commandID string) (*service.IoTCommand, error) {
	const q = `SELECT id, device_id, command_type, payload, status, created_at
		FROM iot_commands WHERE id = $1`
	var c service.IoTCommand
	err := r.pool.QueryRow(ctx, q, commandID).Scan(&c.ID, &c.DeviceID, &c.CommandType, &c.Payload, &c.Status, &c.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *IoTRepository) ReceiveData(ctx context.Context, data *service.IoTData) error {
	const q = `INSERT INTO iot_data (id, device_id, data_type, value, unit, received_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.pool.Exec(ctx, q, data.ID, data.DeviceID, data.DataType, data.Value, data.Unit, data.ReceivedAt)
	return err
}

func (r *IoTRepository) ListData(ctx context.Context, deviceID string, limit, offset int) ([]*service.IoTData, error) {
	const q = `SELECT id, device_id, data_type, value, unit, received_at
		FROM iot_data WHERE device_id = $1 ORDER BY received_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, q, deviceID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*service.IoTData
	for rows.Next() {
		var d service.IoTData
		if err := rows.Scan(&d.ID, &d.DeviceID, &d.DataType, &d.Value, &d.Unit, &d.ReceivedAt); err != nil {
			return nil, err
		}
		list = append(list, &d)
	}
	return list, rows.Err()
}
GOEOF

# 6. vision-service
mkdir -p "$BASE/vision-service/internal/repository/postgres"
cat > "$BASE/vision-service/internal/repository/postgres/vision.go" << 'GOEOF'
package postgres

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/vision-service/internal/service"
)

type FoodAnalysisRepository struct {
	pool *pgxpool.Pool
}

func NewFoodAnalysisRepository(pool *pgxpool.Pool) *FoodAnalysisRepository {
	return &FoodAnalysisRepository{pool: pool}
}

func (r *FoodAnalysisRepository) Save(ctx context.Context, analysis *service.FoodAnalysis) error {
	itemsJSON, _ := json.Marshal(analysis.FoodItems)
	const q = `INSERT INTO food_analyses (id, user_id, image_url, status, total_calorie_kcal, food_items, meal_type, analyzed_at, created_at, error_message)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.pool.Exec(ctx, q,
		analysis.ID, analysis.UserID, analysis.ImageURL, int32(analysis.Status),
		analysis.TotalCalorieKcal, itemsJSON, analysis.MealType,
		analysis.AnalyzedAt, analysis.CreatedAt, analysis.ErrorMessage,
	)
	return err
}

func (r *FoodAnalysisRepository) FindByID(ctx context.Context, id string) (*service.FoodAnalysis, error) {
	const q = `SELECT id, user_id, image_url, status, total_calorie_kcal, food_items, meal_type, analyzed_at, created_at, error_message
		FROM food_analyses WHERE id = $1`
	var a service.FoodAnalysis
	var statusInt int32
	var itemsJSON []byte
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&a.ID, &a.UserID, &a.ImageURL, &statusInt,
		&a.TotalCalorieKcal, &itemsJSON, &a.MealType,
		&a.AnalyzedAt, &a.CreatedAt, &a.ErrorMessage,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	a.Status = service.FoodAnalysisStatus(statusInt)
	if len(itemsJSON) > 0 {
		json.Unmarshal(itemsJSON, &a.FoodItems)
	}
	return &a, nil
}

func (r *FoodAnalysisRepository) FindByUserID(ctx context.Context, userID string, limit, offset int32) ([]*service.FoodAnalysis, int32, error) {
	var total int32
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM food_analyses WHERE user_id = $1`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	const q = `SELECT id, user_id, image_url, status, total_calorie_kcal, food_items, meal_type, analyzed_at, created_at, error_message
		FROM food_analyses WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []*service.FoodAnalysis
	for rows.Next() {
		var a service.FoodAnalysis
		var statusInt int32
		var itemsJSON []byte
		if err := rows.Scan(
			&a.ID, &a.UserID, &a.ImageURL, &statusInt,
			&a.TotalCalorieKcal, &itemsJSON, &a.MealType,
			&a.AnalyzedAt, &a.CreatedAt, &a.ErrorMessage,
		); err != nil {
			return nil, 0, err
		}
		a.Status = service.FoodAnalysisStatus(statusInt)
		if len(itemsJSON) > 0 {
			json.Unmarshal(itemsJSON, &a.FoodItems)
		}
		list = append(list, &a)
	}
	return list, total, rows.Err()
}

func (r *FoodAnalysisRepository) Update(ctx context.Context, analysis *service.FoodAnalysis) error {
	itemsJSON, _ := json.Marshal(analysis.FoodItems)
	const q = `UPDATE food_analyses SET status=$1, total_calorie_kcal=$2, food_items=$3, analyzed_at=$4, error_message=$5
		WHERE id = $6`
	_, err := r.pool.Exec(ctx, q,
		int32(analysis.Status), analysis.TotalCalorieKcal, itemsJSON,
		analysis.AnalyzedAt, analysis.ErrorMessage, analysis.ID,
	)
	return err
}
GOEOF

# 7. assistant-service
mkdir -p "$BASE/assistant-service/internal/repository/postgres"
cat > "$BASE/assistant-service/internal/repository/postgres/assistant.go" << 'GOEOF'
package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/assistant-service/internal/service"
)

type AssistantRepository struct {
	pool *pgxpool.Pool
}

func NewAssistantRepository(pool *pgxpool.Pool) *AssistantRepository {
	return &AssistantRepository{pool: pool}
}

func (r *AssistantRepository) CreateSession(ctx context.Context, s *service.AssistantSession) error {
	const q = `INSERT INTO assistant_sessions (id, user_id, title, status, turn_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(ctx, q, s.ID, s.UserID, s.Title, s.Status, s.TurnCount, s.CreatedAt, s.UpdatedAt)
	return err
}

func (r *AssistantRepository) GetSession(ctx context.Context, id string) (*service.AssistantSession, error) {
	const q = `SELECT id, user_id, title, status, turn_count, created_at, updated_at
		FROM assistant_sessions WHERE id = $1`
	var s service.AssistantSession
	err := r.pool.QueryRow(ctx, q, id).Scan(&s.ID, &s.UserID, &s.Title, &s.Status, &s.TurnCount, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *AssistantRepository) ListSessions(ctx context.Context, userID string, limit, offset int32) ([]*service.AssistantSession, int32, error) {
	var total int32
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM assistant_sessions WHERE user_id = $1`, userID).Scan(&total)

	const q = `SELECT id, user_id, title, status, turn_count, created_at, updated_at
		FROM assistant_sessions WHERE user_id = $1 ORDER BY updated_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []*service.AssistantSession
	for rows.Next() {
		var s service.AssistantSession
		if err := rows.Scan(&s.ID, &s.UserID, &s.Title, &s.Status, &s.TurnCount, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, &s)
	}
	return list, total, rows.Err()
}

func (r *AssistantRepository) DeleteSession(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM assistant_turns WHERE session_id = $1`, id)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `DELETE FROM assistant_sessions WHERE id = $1`, id)
	return err
}

func (r *AssistantRepository) AddTurn(ctx context.Context, t *service.AssistantTurn) error {
	const q = `INSERT INTO assistant_turns (id, session_id, role, content, intent, action_type, action_result, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.pool.Exec(ctx, q, t.ID, t.SessionID, t.Role, t.Content, t.Intent, t.ActionType, t.ActionResult, t.CreatedAt)
	return err
}

func (r *AssistantRepository) ListTurns(ctx context.Context, sessionID string, limit, offset int32) ([]*service.AssistantTurn, int32, error) {
	var total int32
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM assistant_turns WHERE session_id = $1`, sessionID).Scan(&total)

	const q = `SELECT id, session_id, role, content, COALESCE(intent,''), COALESCE(action_type,''), COALESCE(action_result,''), created_at
		FROM assistant_turns WHERE session_id = $1 ORDER BY created_at ASC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, q, sessionID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []*service.AssistantTurn
	for rows.Next() {
		var t service.AssistantTurn
		if err := rows.Scan(&t.ID, &t.SessionID, &t.Role, &t.Content, &t.Intent, &t.ActionType, &t.ActionResult, &t.CreatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, &t)
	}
	return list, total, rows.Err()
}

func (r *AssistantRepository) UpdateSession(ctx context.Context, s *service.AssistantSession) error {
	const q = `UPDATE assistant_sessions SET title=$1, status=$2, turn_count=$3, updated_at=$4 WHERE id = $5`
	_, err := r.pool.Exec(ctx, q, s.Title, s.Status, s.TurnCount, s.UpdatedAt, s.ID)
	return err
}
GOEOF

# 8. marketplace-service
mkdir -p "$BASE/marketplace-service/internal/repository/postgres"
cat > "$BASE/marketplace-service/internal/repository/postgres/marketplace.go" << 'GOEOF'
package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/marketplace-service/internal/service"
)

// ProductRepository is the PostgreSQL implementation.
type ProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

func (r *ProductRepository) FindByPartner(ctx context.Context, partnerID string) ([]*service.PartnerProduct, error) {
	return r.queryProducts(ctx, `SELECT id, partner_id, name, description, price, category, image_url, is_active, created_at, updated_at
		FROM marketplace_products WHERE partner_id = $1 AND is_active = TRUE ORDER BY created_at DESC`, partnerID)
}

func (r *ProductRepository) FindByCategory(ctx context.Context, category string) ([]*service.PartnerProduct, error) {
	return r.queryProducts(ctx, `SELECT id, partner_id, name, description, price, category, image_url, is_active, created_at, updated_at
		FROM marketplace_products WHERE category = $1 AND is_active = TRUE ORDER BY created_at DESC`, category)
}

func (r *ProductRepository) FindByPartnerAndCategory(ctx context.Context, partnerID, category string) ([]*service.PartnerProduct, error) {
	return r.queryProducts(ctx, `SELECT id, partner_id, name, description, price, category, image_url, is_active, created_at, updated_at
		FROM marketplace_products WHERE partner_id = $1 AND category = $2 AND is_active = TRUE ORDER BY created_at DESC`, partnerID, category)
}

func (r *ProductRepository) FindAll(ctx context.Context) ([]*service.PartnerProduct, error) {
	return r.queryProducts(ctx, `SELECT id, partner_id, name, description, price, category, image_url, is_active, created_at, updated_at
		FROM marketplace_products WHERE is_active = TRUE ORDER BY created_at DESC LIMIT 100`)
}

func (r *ProductRepository) FindByID(ctx context.Context, id string) (*service.PartnerProduct, error) {
	const q = `SELECT id, partner_id, name, description, price, category, image_url, is_active, created_at, updated_at
		FROM marketplace_products WHERE id = $1`
	var p service.PartnerProduct
	err := r.pool.QueryRow(ctx, q, id).Scan(&p.ID, &p.PartnerID, &p.Name, &p.Description, &p.Price, &p.Category, &p.ImageURL, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Save(ctx context.Context, product *service.PartnerProduct) error {
	const q = `INSERT INTO marketplace_products (id, partner_id, name, description, price, category, image_url, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.pool.Exec(ctx, q, product.ID, product.PartnerID, product.Name, product.Description, product.Price, product.Category, product.ImageURL, product.IsActive, product.CreatedAt, product.UpdatedAt)
	return err
}

func (r *ProductRepository) Update(ctx context.Context, product *service.PartnerProduct) error {
	const q = `UPDATE marketplace_products SET name=$1, description=$2, price=$3, category=$4, image_url=$5, is_active=$6, updated_at=$7 WHERE id = $8`
	_, err := r.pool.Exec(ctx, q, product.Name, product.Description, product.Price, product.Category, product.ImageURL, product.IsActive, product.UpdatedAt, product.ID)
	return err
}

func (r *ProductRepository) queryProducts(ctx context.Context, q string, args ...interface{}) ([]*service.PartnerProduct, error) {
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil { return nil, err }
	defer rows.Close()
	var list []*service.PartnerProduct
	for rows.Next() {
		var p service.PartnerProduct
		if err := rows.Scan(&p.ID, &p.PartnerID, &p.Name, &p.Description, &p.Price, &p.Category, &p.ImageURL, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &p)
	}
	return list, rows.Err()
}

// PartnerRepository is the PostgreSQL implementation.
type PartnerRepository struct {
	pool *pgxpool.Pool
}

func NewPartnerRepository(pool *pgxpool.Pool) *PartnerRepository {
	return &PartnerRepository{pool: pool}
}

func (r *PartnerRepository) Save(ctx context.Context, partner *service.Partner) error {
	const q = `INSERT INTO marketplace_partners (id, name, description, contact_email, status, created_at)
		VALUES ($1, $2, $3, $4, $5::partner_status, $6)`
	_, err := r.pool.Exec(ctx, q, partner.ID, partner.Name, partner.Description, partner.ContactEmail, partner.Status, partner.CreatedAt)
	return err
}

func (r *PartnerRepository) FindByID(ctx context.Context, id string) (*service.Partner, error) {
	const q = `SELECT id, name, description, contact_email, status::text, created_at FROM marketplace_partners WHERE id = $1`
	var p service.Partner
	err := r.pool.QueryRow(ctx, q, id).Scan(&p.ID, &p.Name, &p.Description, &p.ContactEmail, &p.Status, &p.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &p, nil
}

func (r *PartnerRepository) FindAll(ctx context.Context) ([]*service.Partner, error) {
	const q = `SELECT id, name, description, contact_email, status::text, created_at FROM marketplace_partners ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, q)
	if err != nil { return nil, err }
	defer rows.Close()
	var list []*service.Partner
	for rows.Next() {
		var p service.Partner
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.ContactEmail, &p.Status, &p.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, &p)
	}
	return list, rows.Err()
}

// StatsRepository is the PostgreSQL implementation.
type StatsRepository struct {
	pool *pgxpool.Pool
}

func NewStatsRepository(pool *pgxpool.Pool) *StatsRepository {
	return &StatsRepository{pool: pool}
}

func (r *StatsRepository) FindByPartnerID(ctx context.Context, partnerID string) (*service.PartnerStats, error) {
	const q = `SELECT partner_id, total_products, total_orders, revenue FROM marketplace_stats WHERE partner_id = $1`
	var s service.PartnerStats
	err := r.pool.QueryRow(ctx, q, partnerID).Scan(&s.PartnerID, &s.TotalProducts, &s.TotalOrders, &s.Revenue)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &s, nil
}

func (r *StatsRepository) Save(ctx context.Context, stats *service.PartnerStats) error {
	const q = `INSERT INTO marketplace_stats (partner_id, total_products, total_orders, revenue)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (partner_id) DO UPDATE SET total_products=EXCLUDED.total_products, total_orders=EXCLUDED.total_orders, revenue=EXCLUDED.revenue`
	_, err := r.pool.Exec(ctx, q, stats.PartnerID, stats.TotalProducts, stats.TotalOrders, stats.Revenue)
	return err
}
GOEOF

# 9. concept-service
mkdir -p "$BASE/concept-service/internal/repository/postgres"
cat > "$BASE/concept-service/internal/repository/postgres/concept.go" << 'GOEOF'
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
GOEOF

# 10. cartridge-store-service
mkdir -p "$BASE/cartridge-store-service/internal/repository/postgres"
cat > "$BASE/cartridge-store-service/internal/repository/postgres/store.go" << 'GOEOF'
package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/cartridge-store-service/internal/service"
)

type StoreRepository struct {
	pool *pgxpool.Pool
}

func NewStoreRepository(pool *pgxpool.Pool) *StoreRepository {
	return &StoreRepository{pool: pool}
}

func (r *StoreRepository) CreateDeveloper(ctx context.Context, d *service.Developer) error {
	const q = `INSERT INTO cs_developers (id, user_id, company_name, tier, status, created_at) VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.pool.Exec(ctx, q, d.ID, d.UserID, d.CompanyName, d.Tier, d.Status, d.CreatedAt)
	return err
}

func (r *StoreRepository) GetDeveloper(ctx context.Context, userID string) (*service.Developer, error) {
	const q = `SELECT id, user_id, company_name, tier, status, created_at FROM cs_developers WHERE user_id = $1`
	var d service.Developer
	err := r.pool.QueryRow(ctx, q, userID).Scan(&d.ID, &d.UserID, &d.CompanyName, &d.Tier, &d.Status, &d.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &d, nil
}

func (r *StoreRepository) CreateApiKey(ctx context.Context, k *service.ApiKey) error {
	const q = `INSERT INTO cs_api_keys (key_id, developer_id, key, status, created_at) VALUES ($1,$2,$3,$4,$5)`
	_, err := r.pool.Exec(ctx, q, k.KeyID, k.DeveloperID, k.Key, k.Status, k.CreatedAt)
	return err
}

func (r *StoreRepository) CreateItem(ctx context.Context, item *service.StoreItem) error {
	const q = `INSERT INTO cs_store_items (id, developer_id, name, description, category, version, price_krw, status, rating, downloads, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.pool.Exec(ctx, q, item.ID, item.DeveloperID, item.Name, item.Description, item.Category, item.Version, item.PriceKrw, item.Status, item.Rating, item.Downloads, item.CreatedAt)
	return err
}

func (r *StoreRepository) ListItems(ctx context.Context, category string, limit, offset int32) ([]*service.StoreItem, int32, error) {
	var total int32
	if category != "" {
		r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM cs_store_items WHERE category=$1`, category).Scan(&total)
	} else {
		r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM cs_store_items`).Scan(&total)
	}
	q := `SELECT id, developer_id, name, description, category, version, price_krw, status, rating, downloads, created_at FROM cs_store_items`
	var args []interface{}
	if category != "" {
		q += ` WHERE category=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
		args = append(args, category, limit, offset)
	} else {
		q += ` ORDER BY created_at DESC LIMIT $1 OFFSET $2`
		args = append(args, limit, offset)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil { return nil, 0, err }
	defer rows.Close()
	var list []*service.StoreItem
	for rows.Next() {
		var i service.StoreItem
		if err := rows.Scan(&i.ID, &i.DeveloperID, &i.Name, &i.Description, &i.Category, &i.Version, &i.PriceKrw, &i.Status, &i.Rating, &i.Downloads, &i.CreatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, &i)
	}
	return list, total, rows.Err()
}

func (r *StoreRepository) SearchItems(ctx context.Context, query string, limit int32) ([]*service.StoreItem, error) {
	const q = `SELECT id, developer_id, name, description, category, version, price_krw, status, rating, downloads, created_at
		FROM cs_store_items WHERE name ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%' LIMIT $2`
	rows, err := r.pool.Query(ctx, q, query, limit)
	if err != nil { return nil, err }
	defer rows.Close()
	var list []*service.StoreItem
	for rows.Next() {
		var i service.StoreItem
		if err := rows.Scan(&i.ID, &i.DeveloperID, &i.Name, &i.Description, &i.Category, &i.Version, &i.PriceKrw, &i.Status, &i.Rating, &i.Downloads, &i.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, &i)
	}
	return list, rows.Err()
}

func (r *StoreRepository) GetItem(ctx context.Context, id string) (*service.StoreItem, error) {
	const q = `SELECT id, developer_id, name, description, category, version, price_krw, status, rating, downloads, created_at
		FROM cs_store_items WHERE id = $1`
	var i service.StoreItem
	err := r.pool.QueryRow(ctx, q, id).Scan(&i.ID, &i.DeveloperID, &i.Name, &i.Description, &i.Category, &i.Version, &i.PriceKrw, &i.Status, &i.Rating, &i.Downloads, &i.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &i, nil
}

func (r *StoreRepository) CreatePurchase(ctx context.Context, p *service.Purchase) error {
	const q = `INSERT INTO cs_purchases (id, user_id, item_id, price_krw, status, purchased_at) VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.pool.Exec(ctx, q, p.ID, p.UserID, p.ItemID, p.PriceKrw, p.Status, p.PurchasedAt)
	return err
}

func (r *StoreRepository) ListPurchases(ctx context.Context, userID string, limit int32) ([]*service.Purchase, error) {
	const q = `SELECT id, user_id, item_id, price_krw, status, purchased_at FROM cs_purchases WHERE user_id = $1 ORDER BY purchased_at DESC LIMIT $2`
	rows, err := r.pool.Query(ctx, q, userID, limit)
	if err != nil { return nil, err }
	defer rows.Close()
	var list []*service.Purchase
	for rows.Next() {
		var p service.Purchase
		if err := rows.Scan(&p.ID, &p.UserID, &p.ItemID, &p.PriceKrw, &p.Status, &p.PurchasedAt); err != nil { return nil, err }
		list = append(list, &p)
	}
	return list, rows.Err()
}

func (r *StoreRepository) CreateReview(ctx context.Context, rv *service.ReviewStatus) error {
	const q = `INSERT INTO cs_reviews (submission_id, cartridge_id, status, reviewer_comment, submitted_at) VALUES ($1,$2,$3,$4,$5)`
	_, err := r.pool.Exec(ctx, q, rv.SubmissionID, rv.CartridgeID, rv.Status, rv.ReviewerComment, rv.SubmittedAt)
	return err
}

func (r *StoreRepository) GetReview(ctx context.Context, submissionID string) (*service.ReviewStatus, error) {
	const q = `SELECT submission_id, cartridge_id, status, COALESCE(reviewer_comment,''), submitted_at FROM cs_reviews WHERE submission_id = $1`
	var rv service.ReviewStatus
	err := r.pool.QueryRow(ctx, q, submissionID).Scan(&rv.SubmissionID, &rv.CartridgeID, &rv.Status, &rv.ReviewerComment, &rv.SubmittedAt)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &rv, nil
}

func (r *StoreRepository) UpdateReview(ctx context.Context, rv *service.ReviewStatus) error {
	const q = `UPDATE cs_reviews SET status=$1, reviewer_comment=$2 WHERE submission_id = $3`
	_, err := r.pool.Exec(ctx, q, rv.Status, rv.ReviewerComment, rv.SubmissionID)
	return err
}

func (r *StoreRepository) ListSubmissions(ctx context.Context, devID string, limit int32) ([]*service.ReviewStatus, error) {
	const q = `SELECT submission_id, cartridge_id, status, COALESCE(reviewer_comment,''), submitted_at FROM cs_reviews
		WHERE cartridge_id IN (SELECT id FROM cs_store_items WHERE developer_id = $1) ORDER BY submitted_at DESC LIMIT $2`
	rows, err := r.pool.Query(ctx, q, devID, limit)
	if err != nil { return nil, err }
	defer rows.Close()
	var list []*service.ReviewStatus
	for rows.Next() {
		var rv service.ReviewStatus
		if err := rows.Scan(&rv.SubmissionID, &rv.CartridgeID, &rv.Status, &rv.ReviewerComment, &rv.SubmittedAt); err != nil { return nil, err }
		list = append(list, &rv)
	}
	return list, rows.Err()
}
GOEOF

# 11. nlp-service
mkdir -p "$BASE/nlp-service/internal/repository/postgres"
cat > "$BASE/nlp-service/internal/repository/postgres/nlp.go" << 'GOEOF'
package postgres

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/nlp-service/internal/service"
)

type NLPRepository struct {
	pool *pgxpool.Pool
}

func NewNLPRepository(pool *pgxpool.Pool) *NLPRepository {
	return &NLPRepository{pool: pool}
}

func (r *NLPRepository) SaveQuery(ctx context.Context, query *service.HealthQuery) error {
	entJSON, _ := json.Marshal(query.Entities)
	const q = `INSERT INTO nlp_health_queries (id, user_id, raw_text, intent, entities, confidence, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(ctx, q, query.ID, query.UserID, query.RawText, query.Intent, entJSON, query.Confidence, query.CreatedAt)
	return err
}

func (r *NLPRepository) GetQuery(ctx context.Context, queryID string) (*service.HealthQuery, error) {
	const q = `SELECT id, user_id, raw_text, intent, entities, confidence, created_at
		FROM nlp_health_queries WHERE id = $1`
	var hq service.HealthQuery
	var entJSON []byte
	err := r.pool.QueryRow(ctx, q, queryID).Scan(&hq.ID, &hq.UserID, &hq.RawText, &hq.Intent, &entJSON, &hq.Confidence, &hq.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	if len(entJSON) > 0 { json.Unmarshal(entJSON, &hq.Entities) }
	return &hq, nil
}

func (r *NLPRepository) SaveExtraction(ctx context.Context, extraction *service.SymptomExtraction) error {
	sympJSON, _ := json.Marshal(extraction.Symptoms)
	const q = `INSERT INTO nlp_symptom_extractions (id, text, symptoms, processed_at)
		VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, q, extraction.ID, extraction.Text, sympJSON, extraction.ProcessedAt)
	return err
}

func (r *NLPRepository) GetSuggestions(ctx context.Context, queryID string) ([]service.Suggestion, error) {
	const q = `SELECT id, query_id, text, category, priority FROM nlp_suggestions WHERE query_id = $1 ORDER BY priority DESC`
	rows, err := r.pool.Query(ctx, q, queryID)
	if err != nil { return nil, err }
	defer rows.Close()
	var list []service.Suggestion
	for rows.Next() {
		var s service.Suggestion
		if err := rows.Scan(&s.ID, &s.QueryID, &s.Text, &s.Category, &s.Priority); err != nil { return nil, err }
		list = append(list, s)
	}
	return list, rows.Err()
}
GOEOF

# 12. data-platform-service
mkdir -p "$BASE/data-platform-service/internal/repository/postgres"
cat > "$BASE/data-platform-service/internal/repository/postgres/platform.go" << 'GOEOF'
package postgres

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/data-platform-service/internal/service"
)

type PlatformRepository struct {
	pool *pgxpool.Pool
}

func NewPlatformRepository(pool *pgxpool.Pool) *PlatformRepository {
	return &PlatformRepository{pool: pool}
}

func (r *PlatformRepository) GetRegionStats(ctx context.Context, code, period, biomarker string) (*service.RegionStats, error) {
	const q = `SELECT region_code, region_name, period, avg_value, min_value, max_value, sample_count, risk_level
		FROM dp_region_stats WHERE region_code = $1 AND period = $2`
	var s service.RegionStats
	err := r.pool.QueryRow(ctx, q, code, period).Scan(&s.RegionCode, &s.RegionName, &s.Period, &s.AvgValue, &s.MinValue, &s.MaxValue, &s.SampleCount, &s.RiskLevel)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &s, nil
}

func (r *PlatformRepository) ListHotspots(ctx context.Context, biomarker, risk string, limit int32) ([]*service.RegionStats, error) {
	const q = `SELECT region_code, region_name, period, avg_value, min_value, max_value, sample_count, risk_level
		FROM dp_region_stats WHERE risk_level = $1 ORDER BY avg_value DESC LIMIT $2`
	rows, err := r.pool.Query(ctx, q, risk, limit)
	if err != nil { return nil, err }
	defer rows.Close()
	var list []*service.RegionStats
	for rows.Next() {
		var s service.RegionStats
		if err := rows.Scan(&s.RegionCode, &s.RegionName, &s.Period, &s.AvgValue, &s.MinValue, &s.MaxValue, &s.SampleCount, &s.RiskLevel); err != nil { return nil, err }
		list = append(list, &s)
	}
	return list, rows.Err()
}

func (r *PlatformRepository) GetTrend(ctx context.Context, code, biomarker string) (*service.RegionTrend, error) {
	const q = `SELECT region_code, biomarker, trend_direction, points FROM dp_region_trends WHERE region_code = $1 AND biomarker = $2`
	var t service.RegionTrend
	var pointsJSON []byte
	err := r.pool.QueryRow(ctx, q, code, biomarker).Scan(&t.RegionCode, &t.Biomarker, &t.TrendDirection, &pointsJSON)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	if len(pointsJSON) > 0 { json.Unmarshal(pointsJSON, &t.Points) }
	return &t, nil
}

func (r *PlatformRepository) ListDatasets(ctx context.Context, dsType string, limit, offset int32) ([]*service.AnonymizedDataset, int32, error) {
	var total int32
	if dsType != "" {
		r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM dp_datasets WHERE dataset_type=$1`, dsType).Scan(&total)
	} else {
		r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM dp_datasets`).Scan(&total)
	}
	q := `SELECT id, dataset_type, description, record_count, date_range, created_at FROM dp_datasets`
	var args []interface{}
	if dsType != "" {
		q += ` WHERE dataset_type=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
		args = append(args, dsType, limit, offset)
	} else {
		q += ` ORDER BY created_at DESC LIMIT $1 OFFSET $2`
		args = append(args, limit, offset)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil { return nil, 0, err }
	defer rows.Close()
	var list []*service.AnonymizedDataset
	for rows.Next() {
		var d service.AnonymizedDataset
		if err := rows.Scan(&d.ID, &d.DatasetType, &d.Description, &d.RecordCount, &d.DateRange, &d.CreatedAt); err != nil { return nil, 0, err }
		list = append(list, &d)
	}
	return list, total, rows.Err()
}

func (r *PlatformRepository) GetDataset(ctx context.Context, id string) (*service.AnonymizedDataset, error) {
	const q = `SELECT id, dataset_type, description, record_count, date_range, created_at FROM dp_datasets WHERE id = $1`
	var d service.AnonymizedDataset
	err := r.pool.QueryRow(ctx, q, id).Scan(&d.ID, &d.DatasetType, &d.Description, &d.RecordCount, &d.DateRange, &d.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &d, nil
}

func (r *PlatformRepository) SaveDataset(ctx context.Context, ds *service.AnonymizedDataset) error {
	const q = `INSERT INTO dp_datasets (id, dataset_type, description, record_count, date_range, created_at)
		VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT (id) DO NOTHING`
	_, err := r.pool.Exec(ctx, q, ds.ID, ds.DatasetType, ds.Description, ds.RecordCount, ds.DateRange, ds.CreatedAt)
	return err
}

func (r *PlatformRepository) CreateAccessRequest(ctx context.Context, requesterID, datasetID, purpose string) (*service.DataAccessResp, error) {
	return &service.DataAccessResp{RequestID: "req-" + requesterID[:8], Status: "pending"}, nil
}

func (r *PlatformRepository) GetConsent(ctx context.Context, userID string) (*service.ConsentInfo, error) {
	const q = `SELECT user_id, anonymous_research, region_stats_consent, enterprise_data, updated_at
		FROM dp_consents WHERE user_id = $1`
	var c service.ConsentInfo
	err := r.pool.QueryRow(ctx, q, userID).Scan(&c.UserID, &c.AnonymousResearch, &c.RegionStatsConsent, &c.EnterpriseData, &c.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &c, nil
}

func (r *PlatformRepository) UpdateConsent(ctx context.Context, c *service.ConsentInfo) error {
	const q = `INSERT INTO dp_consents (user_id, anonymous_research, region_stats_consent, enterprise_data, updated_at)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (user_id) DO UPDATE SET anonymous_research=EXCLUDED.anonymous_research, region_stats_consent=EXCLUDED.region_stats_consent, enterprise_data=EXCLUDED.enterprise_data, updated_at=EXCLUDED.updated_at`
	_, err := r.pool.Exec(ctx, q, c.UserID, c.AnonymousResearch, c.RegionStatsConsent, c.EnterpriseData, c.UpdatedAt)
	return err
}

func (r *PlatformRepository) ListInsights(ctx context.Context, biomarker string, limit int32) ([]*service.AggregatedInsight, error) {
	const q = `SELECT id, biomarker, insight_type, summary, created_at FROM dp_insights WHERE biomarker = $1 ORDER BY created_at DESC LIMIT $2`
	rows, err := r.pool.Query(ctx, q, biomarker, limit)
	if err != nil { return nil, err }
	defer rows.Close()
	var list []*service.AggregatedInsight
	for rows.Next() {
		var i service.AggregatedInsight
		if err := rows.Scan(&i.ID, &i.Biomarker, &i.InsightType, &i.Summary, &i.CreatedAt); err != nil { return nil, err }
		list = append(list, &i)
	}
	return list, rows.Err()
}

func (r *PlatformRepository) SaveInsight(ctx context.Context, i *service.AggregatedInsight) error {
	const q = `INSERT INTO dp_insights (id, biomarker, insight_type, summary, created_at) VALUES ($1,$2,$3,$4,$5)`
	_, err := r.pool.Exec(ctx, q, i.ID, i.Biomarker, i.InsightType, i.Summary, i.CreatedAt)
	return err
}
GOEOF

# 13. voice-profile-service
mkdir -p "$BASE/voice-profile-service/internal/repository/postgres"
cat > "$BASE/voice-profile-service/internal/repository/postgres/voice_profile.go" << 'GOEOF'
package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/voice-profile-service/internal/service"
)

type VoiceProfileRepository struct {
	pool *pgxpool.Pool
}

func NewVoiceProfileRepository(pool *pgxpool.Pool) *VoiceProfileRepository {
	return &VoiceProfileRepository{pool: pool}
}

func (r *VoiceProfileRepository) Create(ctx context.Context, p *service.VoiceProfile) error {
	const q = `INSERT INTO voice_profiles (profile_id, user_id, profile_name, language, status, model_url, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(ctx, q, p.ProfileID, p.UserID, p.ProfileName, p.Language, p.Status, p.ModelURL, p.CreatedAt)
	return err
}

func (r *VoiceProfileRepository) GetByID(ctx context.Context, profileID string) (*service.VoiceProfile, error) {
	const q = `SELECT profile_id, user_id, profile_name, language, status, COALESCE(model_url,''), created_at
		FROM voice_profiles WHERE profile_id = $1`
	var p service.VoiceProfile
	err := r.pool.QueryRow(ctx, q, profileID).Scan(&p.ProfileID, &p.UserID, &p.ProfileName, &p.Language, &p.Status, &p.ModelURL, &p.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &p, nil
}

func (r *VoiceProfileRepository) ListByUser(ctx context.Context, userID string) ([]*service.VoiceProfile, error) {
	const q = `SELECT profile_id, user_id, profile_name, language, status, COALESCE(model_url,''), created_at
		FROM voice_profiles WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, q, userID)
	if err != nil { return nil, err }
	defer rows.Close()
	var list []*service.VoiceProfile
	for rows.Next() {
		var p service.VoiceProfile
		if err := rows.Scan(&p.ProfileID, &p.UserID, &p.ProfileName, &p.Language, &p.Status, &p.ModelURL, &p.CreatedAt); err != nil { return nil, err }
		list = append(list, &p)
	}
	return list, rows.Err()
}

func (r *VoiceProfileRepository) Delete(ctx context.Context, profileID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM voice_profiles WHERE profile_id = $1`, profileID)
	return err
}
GOEOF

echo "=== 13개 PostgreSQL Repository 생성 완료 ==="
ls -la "$BASE"/analytics-service/internal/repository/postgres/
ls -la "$BASE"/emergency-service/internal/repository/postgres/
ls -la "$BASE"/audit-service/internal/repository/postgres/
ls -la "$BASE"/digital-twin-service/internal/repository/postgres/
ls -la "$BASE"/iot-gateway-service/internal/repository/postgres/
ls -la "$BASE"/vision-service/internal/repository/postgres/
ls -la "$BASE"/assistant-service/internal/repository/postgres/
ls -la "$BASE"/marketplace-service/internal/repository/postgres/
ls -la "$BASE"/concept-service/internal/repository/postgres/
ls -la "$BASE"/cartridge-store-service/internal/repository/postgres/
ls -la "$BASE"/nlp-service/internal/repository/postgres/
ls -la "$BASE"/data-platform-service/internal/repository/postgres/
ls -la "$BASE"/voice-profile-service/internal/repository/postgres/
