-- digital-twin-service 초기화 스크립트
-- Phase A: 디지털 트윈 상태 영속화

-- 트윈 상태 enum
DO $$ BEGIN
    CREATE TYPE twin_health_state AS ENUM ('nominal', 'warning', 'drift', 'recalibrate');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- 디지털 트윈 상태 테이블
CREATE TABLE IF NOT EXISTS twin_states (
    id                      VARCHAR(64) PRIMARY KEY,
    session_id              VARCHAR(64) NOT NULL,
    user_id                 VARCHAR(64),
    device_id               VARCHAR(64),
    ewma_value              DOUBLE PRECISION NOT NULL DEFAULT 0,
    cusum_pos               DOUBLE PRECISION NOT NULL DEFAULT 0,
    cusum_neg               DOUBLE PRECISION NOT NULL DEFAULT 0,
    health_state            VARCHAR(20) NOT NULL DEFAULT 'nominal',
    drift_score             DOUBLE PRECISION NOT NULL DEFAULT 0,
    remaining_measurements  INT NOT NULL DEFAULT 0,
    measurement_count       INT NOT NULL DEFAULT 0,
    last_synced_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_twin_states_session ON twin_states(session_id);
CREATE INDEX IF NOT EXISTS idx_twin_states_device ON twin_states(device_id);
CREATE INDEX IF NOT EXISTS idx_twin_states_user ON twin_states(user_id);
