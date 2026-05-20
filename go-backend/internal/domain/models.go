package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// User represents the users table
type User struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Email     string     `json:"email" db:"email"`
	Nickname  string     `json:"nickname" db:"nickname"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"` // Pointer for nullable column
}

// Device represents the STM32 hardware readers in the devices table
type Device struct {
	MACAddress      string    `json:"mac_address" db:"mac_address"`
	UserID          uuid.UUID `json:"user_id" db:"user_id"`
	HardwareVersion string    `json:"hardware_version" db:"hardware_version"`
	FirmwareVersion string    `json:"firmware_version" db:"firmware_version"`
	Status          string    `json:"status" db:"status"`
	LastSyncAt      time.Time `json:"last_sync_at" db:"last_sync_at"`
}

// Measurement represents the complex layered analytical data
type Measurement struct {
	ID            uuid.UUID       `json:"id" db:"id"`
	UserID        uuid.UUID       `json:"user_id" db:"user_id"`
	DeviceMAC     string          `json:"device_mac" db:"device_mac"`
	MeasuredAt    time.Time       `json:"measured_at" db:"measured_at"`
	DiffSignal    json.RawMessage `json:"diff_signal" db:"diff_signal"` // JSON stored array of raw float signals
	Fingerprint   json.RawMessage `json:"fingerprint" db:"fingerprint"` // JSON stored 896-dimension array
	HealthScore   int             `json:"health_score" db:"health_score"`
	RiskLabel     string          `json:"risk_label" db:"risk_label"`
	ClientLocalID string          `json:"client_local_id" db:"client_local_id"`
	SyncedAt      time.Time       `json:"synced_at" db:"synced_at"`
}

// DiagnosticLog represents hardware/firmware anomalies reported by Rust Layer
type DiagnosticLog struct {
	ID         uuid.UUID       `json:"id" db:"id"`
	DeviceMAC  string          `json:"device_mac" db:"device_mac"`
	EventType  string          `json:"event_type" db:"event_type"`
	Severity   string          `json:"severity" db:"severity"`
	Payload    json.RawMessage `json:"payload" db:"payload"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
}
