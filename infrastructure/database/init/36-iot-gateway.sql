-- iot-gateway-service 초기화 스크립트
-- Phase A: IoT 디바이스/데이터 영속화

-- IoT 디바이스 테이블
CREATE TABLE IF NOT EXISTS iot_devices (
    id              VARCHAR(64) PRIMARY KEY,
    device_id       VARCHAR(64) NOT NULL UNIQUE,
    protocol        VARCHAR(30) NOT NULL DEFAULT 'ble',
    last_ping_at    TIMESTAMPTZ,
    status          VARCHAR(20) NOT NULL DEFAULT 'offline',
    metadata        JSONB DEFAULT '{}'
);

CREATE INDEX IF NOT EXISTS idx_iot_devices_status ON iot_devices(status);

-- IoT 명령 테이블
CREATE TABLE IF NOT EXISTS iot_commands (
    id              VARCHAR(64) PRIMARY KEY,
    device_id       VARCHAR(64) NOT NULL,
    command_type    VARCHAR(50) NOT NULL,
    payload         TEXT,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_iot_commands_device ON iot_commands(device_id);
CREATE INDEX IF NOT EXISTS idx_iot_commands_status ON iot_commands(status);

-- IoT 데이터 테이블
CREATE TABLE IF NOT EXISTS iot_data (
    id              VARCHAR(64) PRIMARY KEY,
    device_id       VARCHAR(64) NOT NULL,
    data_type       VARCHAR(50) NOT NULL,
    value           DOUBLE PRECISION NOT NULL,
    unit            VARCHAR(20),
    received_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_iot_data_device ON iot_data(device_id);
CREATE INDEX IF NOT EXISTS idx_iot_data_received ON iot_data(device_id, received_at DESC);
