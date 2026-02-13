// Package postgres는 device-service의 PostgreSQL 저장소 구현입니다.
package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/device-service/internal/service"
)

// ============================================================================
// DeviceRepository
// ============================================================================

// DeviceRepository는 PostgreSQL 기반 DeviceRepository 구현입니다.
type DeviceRepository struct {
	pool *pgxpool.Pool
}

// NewDeviceRepository는 DeviceRepository를 생성합니다.
func NewDeviceRepository(pool *pgxpool.Pool) *DeviceRepository {
	return &DeviceRepository{pool: pool}
}

// Create는 디바이스를 생성합니다.
func (r *DeviceRepository) Create(ctx context.Context, device *service.Device) error {
	const q = `INSERT INTO devices (id, device_id, user_id, name, serial_number, firmware_version, status, battery_percent, last_seen, registered_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.pool.Exec(ctx, q,
		device.ID,
		device.DeviceID,
		device.UserID,
		device.Name,
		device.SerialNumber,
		device.FirmwareVersion,
		string(device.Status),
		device.BatteryPercent,
		device.LastSeen,
		device.RegisteredAt,
	)
	return err
}

// GetByID는 ID로 디바이스를 조회합니다.
func (r *DeviceRepository) GetByID(ctx context.Context, deviceID string) (*service.Device, error) {
	const q = `SELECT id, device_id, user_id, COALESCE(name, ''), serial_number, firmware_version, status, battery_percent, last_seen, registered_at
		FROM devices WHERE id = $1`
	var d service.Device
	var status string
	err := r.pool.QueryRow(ctx, q, deviceID).Scan(
		&d.ID, &d.DeviceID, &d.UserID, &d.Name, &d.SerialNumber,
		&d.FirmwareVersion, &status, &d.BatteryPercent, &d.LastSeen, &d.RegisteredAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	d.Status = service.DeviceStatus(status)
	return &d, nil
}

// ListByUser는 사용자의 디바이스 목록을 조회합니다.
func (r *DeviceRepository) ListByUser(ctx context.Context, userID string) ([]*service.Device, error) {
	const q = `SELECT id, device_id, user_id, COALESCE(name, ''), serial_number, firmware_version, status, battery_percent, last_seen, registered_at
		FROM devices WHERE user_id = $1 ORDER BY registered_at DESC`
	rows, err := r.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []*service.Device
	for rows.Next() {
		var d service.Device
		var status string
		if err := rows.Scan(
			&d.ID, &d.DeviceID, &d.UserID, &d.Name, &d.SerialNumber,
			&d.FirmwareVersion, &status, &d.BatteryPercent, &d.LastSeen, &d.RegisteredAt,
		); err != nil {
			return nil, err
		}
		d.Status = service.DeviceStatus(status)
		devices = append(devices, &d)
	}
	return devices, rows.Err()
}

// UpdateStatus는 디바이스 상태를 업데이트합니다.
func (r *DeviceRepository) UpdateStatus(ctx context.Context, deviceID string, status service.DeviceStatus, battery int, lastSeen time.Time) error {
	const q = `UPDATE devices SET status = $1, battery_percent = $2, last_seen = $3 WHERE id = $4`
	_, err := r.pool.Exec(ctx, q, string(status), battery, lastSeen, deviceID)
	return err
}

// CountByUser는 사용자의 디바이스 수를 반환합니다.
func (r *DeviceRepository) CountByUser(ctx context.Context, userID string) (int, error) {
	const q = `SELECT COUNT(*) FROM devices WHERE user_id = $1`
	var count int
	err := r.pool.QueryRow(ctx, q, userID).Scan(&count)
	return count, err
}

// Delete는 디바이스를 삭제합니다.
func (r *DeviceRepository) Delete(ctx context.Context, deviceID string) error {
	const q = `DELETE FROM devices WHERE id = $1`
	_, err := r.pool.Exec(ctx, q, deviceID)
	return err
}

// ============================================================================
// DeviceEventRepository
// ============================================================================

// DeviceEventRepository는 PostgreSQL 기반 DeviceEventRepository 구현입니다.
type DeviceEventRepository struct {
	pool *pgxpool.Pool
}

// NewDeviceEventRepository는 DeviceEventRepository를 생성합니다.
func NewDeviceEventRepository(pool *pgxpool.Pool) *DeviceEventRepository {
	return &DeviceEventRepository{pool: pool}
}

// LogEvent는 디바이스 이벤트를 기록합니다.
func (r *DeviceEventRepository) LogEvent(ctx context.Context, event *service.DeviceEvent) error {
	payloadJSON, err := json.Marshal(event.Payload)
	if err != nil {
		payloadJSON = []byte("{}")
	}
	const q = `INSERT INTO device_events (id, device_id, event_type, payload, created_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err = r.pool.Exec(ctx, q,
		event.ID,
		event.DeviceID,
		event.EventType,
		payloadJSON,
		event.CreatedAt,
	)
	return err
}
