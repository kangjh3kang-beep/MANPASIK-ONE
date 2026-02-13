-- ManPaSik measurement-service 데이터베이스 초기화
-- 측정 세션, 시계열 데이터, 핑거프린트 벡터 테이블
-- TimescaleDB 확장 사용 (가용 시)

-- TimescaleDB 확장 활성화 (설치된 경우만)
DO $$
BEGIN
    CREATE EXTENSION IF NOT EXISTS timescaledb;
EXCEPTION WHEN OTHERS THEN
    RAISE NOTICE 'TimescaleDB 확장 없음 - 일반 PostgreSQL 테이블 사용';
END
$$;

-- 측정 세션 테이블
CREATE TABLE IF NOT EXISTS measurement_sessions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id           VARCHAR(100) NOT NULL,
    cartridge_id        VARCHAR(100) NOT NULL,
    user_id             UUID NOT NULL,
    status              VARCHAR(20) NOT NULL DEFAULT 'active',
    total_measurements  INT NOT NULL DEFAULT 0,
    started_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ended_at            TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_sessions_user ON measurement_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_device ON measurement_sessions(device_id);
CREATE INDEX IF NOT EXISTS idx_sessions_status ON measurement_sessions(status);
CREATE INDEX IF NOT EXISTS idx_sessions_started ON measurement_sessions(started_at DESC);

-- 측정 데이터 시계열 테이블
CREATE TABLE IF NOT EXISTS measurement_data (
    time              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    session_id        UUID NOT NULL REFERENCES measurement_sessions(id) ON DELETE CASCADE,
    device_id         VARCHAR(100) NOT NULL,
    user_id           UUID NOT NULL,
    cartridge_type    VARCHAR(50) NOT NULL,
    raw_channels      DOUBLE PRECISION[] DEFAULT '{}',
    s_det             DOUBLE PRECISION NOT NULL DEFAULT 0,
    s_ref             DOUBLE PRECISION NOT NULL DEFAULT 0,
    alpha             DOUBLE PRECISION NOT NULL DEFAULT 0.95,
    s_corrected       DOUBLE PRECISION NOT NULL DEFAULT 0,
    primary_value     DOUBLE PRECISION NOT NULL DEFAULT 0,
    unit              VARCHAR(20) NOT NULL DEFAULT '',
    confidence        DOUBLE PRECISION NOT NULL DEFAULT 0,
    fingerprint_dim   INT NOT NULL DEFAULT 0,
    temp_c            REAL DEFAULT 0,
    humidity_pct      REAL DEFAULT 0,
    battery_pct       INT DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_mdata_session ON measurement_data(session_id);
CREATE INDEX IF NOT EXISTS idx_mdata_user ON measurement_data(user_id);
CREATE INDEX IF NOT EXISTS idx_mdata_time ON measurement_data(time DESC);

-- TimescaleDB 하이퍼테이블 변환 (TimescaleDB 가용 시)
DO $$
BEGIN
    PERFORM create_hypertable('measurement_data', 'time', if_not_exists => TRUE);
    RAISE NOTICE 'measurement_data를 TimescaleDB 하이퍼테이블로 변환 완료';
EXCEPTION WHEN OTHERS THEN
    RAISE NOTICE 'TimescaleDB 하이퍼테이블 변환 건너뜀 (일반 테이블 유지)';
END
$$;

-- 측정 결과 요약 뷰 (자주 사용되는 조회 최적화)
CREATE OR REPLACE VIEW measurement_summary AS
SELECT
    ms.id AS session_id,
    md.cartridge_type,
    md.primary_value,
    md.unit,
    md.time AS measured_at,
    ms.user_id
FROM measurement_sessions ms
JOIN measurement_data md ON ms.id = md.session_id
WHERE ms.status = 'completed';
