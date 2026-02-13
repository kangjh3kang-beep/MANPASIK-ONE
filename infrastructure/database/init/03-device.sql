-- ManPaSik device-service 데이터베이스 초기화
-- 디바이스, 디바이스 이벤트 테이블

-- 디바이스 테이블
CREATE TABLE IF NOT EXISTS devices (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id         VARCHAR(100) NOT NULL,   -- 하드웨어 고유 ID (BLE MAC / 시리얼)
    user_id           UUID NOT NULL,
    name              VARCHAR(200) DEFAULT '',
    serial_number     VARCHAR(100) NOT NULL,
    firmware_version  VARCHAR(50) NOT NULL,
    status            VARCHAR(20) NOT NULL DEFAULT 'online',
    battery_percent   INT NOT NULL DEFAULT 100,
    last_seen         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    registered_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(device_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_devices_user_id ON devices(user_id);
CREATE INDEX IF NOT EXISTS idx_devices_device_id ON devices(device_id);
CREATE INDEX IF NOT EXISTS idx_devices_status ON devices(status);

-- 디바이스 이벤트 로그 테이블
CREATE TABLE IF NOT EXISTS device_events (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id       UUID NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    event_type      VARCHAR(50) NOT NULL,
    payload         JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_device_events_device ON device_events(device_id);
CREATE INDEX IF NOT EXISTS idx_device_events_type ON device_events(event_type);
CREATE INDEX IF NOT EXISTS idx_device_events_created ON device_events(created_at DESC);

-- 펌웨어 정보 테이블
CREATE TABLE IF NOT EXISTS firmware_versions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version         VARCHAR(50) NOT NULL UNIQUE,
    release_date    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    download_url    TEXT NOT NULL,
    checksum        VARCHAR(128) NOT NULL,
    size_bytes      BIGINT NOT NULL DEFAULT 0,
    release_notes   TEXT DEFAULT '',
    mandatory       BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
