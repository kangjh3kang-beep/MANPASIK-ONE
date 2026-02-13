-- Measurement Service Database Initialization (TimescaleDB)
-- 테이블은 POSTGRES_DB(기본 manpasik_ts)에 생성됩니다.

-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TABLE IF NOT EXISTS measurements (
    time TIMESTAMP WITH TIME ZONE NOT NULL,
    device_id VARCHAR(100) NOT NULL,
    user_id UUID NOT NULL,
    cartridge_type VARCHAR(50) NOT NULL,
    
    -- 차동측정 원시 데이터
    s_det DOUBLE PRECISION,
    s_ref DOUBLE PRECISION,
    alpha DOUBLE PRECISION DEFAULT 0.95,
    s_corrected DOUBLE PRECISION, -- s_det - alpha * s_ref
    
    -- 결과
    primary_value DOUBLE PRECISION,
    unit VARCHAR(20),
    status VARCHAR(20), -- success, fail, error
    
    -- 환경
    temperature DOUBLE PRECISION,
    humidity DOUBLE PRECISION,
    battery_level INTEGER
);

-- Convert to Hypertable
SELECT create_hypertable('measurements', 'time');

CREATE INDEX idx_measurements_user_time ON measurements (user_id, time DESC);
CREATE INDEX idx_measurements_device_time ON measurements (device_id, time DESC);
