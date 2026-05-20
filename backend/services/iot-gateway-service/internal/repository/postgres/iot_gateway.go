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

func (r *IoTRepository) UpdateDevice(ctx context.Context, device *service.IoTDevice) error {
	metaJSON, _ := json.Marshal(device.Metadata)
	const q = `UPDATE iot_devices SET protocol=$1, last_ping_at=$2, status=$3, metadata=$4 WHERE device_id=$5`
	tag, err := r.pool.Exec(ctx, q, device.Protocol, device.LastPingAt, device.Status, metaJSON, device.DeviceID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *IoTRepository) UpdateCommand(ctx context.Context, cmd *service.IoTCommand) error {
	const q = `UPDATE iot_commands SET status=$1 WHERE id=$2`
	tag, err := r.pool.Exec(ctx, q, cmd.Status, cmd.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *IoTRepository) ListCommandsByDevice(ctx context.Context, deviceID string, limit int) ([]*service.IoTCommand, error) {
	const q = `SELECT id, device_id, command_type, payload, status, created_at
		FROM iot_commands WHERE device_id = $1 ORDER BY created_at DESC LIMIT $2`
	rows, err := r.pool.Query(ctx, q, deviceID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*service.IoTCommand
	for rows.Next() {
		var c service.IoTCommand
		if err := rows.Scan(&c.ID, &c.DeviceID, &c.CommandType, &c.Payload, &c.Status, &c.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, &c)
	}
	return list, rows.Err()
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
