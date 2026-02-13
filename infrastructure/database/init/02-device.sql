-- Device Service Database Initialization
-- 테이블은 POSTGRES_DB(기본 manpasik)에 생성됩니다.

CREATE TABLE IF NOT EXISTS devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id VARCHAR(100) UNIQUE NOT NULL, -- Serial or MAC
    user_id UUID NOT NULL, -- Logical reference to user-service
    name VARCHAR(100),
    firmware_version VARCHAR(50),
    status VARCHAR(20) DEFAULT 'offline',
    last_seen_at TIMESTAMP WITH TIME ZONE,
    registered_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    metadata JSONB
);

CREATE TABLE IF NOT EXISTS firmware_updates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version VARCHAR(50) NOT NULL,
    file_url TEXT NOT NULL,
    checksum VARCHAR(64) NOT NULL,
    release_notes TEXT,
    is_mandatory BOOLEAN DEFAULT FALSE,
    released_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_devices_user_id ON devices(user_id);
CREATE INDEX idx_devices_device_id ON devices(device_id);
