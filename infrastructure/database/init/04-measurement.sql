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

    -- 차동측정 원시 데이터 (SSOT: S_diff = S_det - alpha * S_ref)
    raw_channels      DOUBLE PRECISION[] DEFAULT '{}',
    s_det             DOUBLE PRECISION NOT NULL DEFAULT 0,
    s_ref             DOUBLE PRECISION NOT NULL DEFAULT 0,
    alpha             DOUBLE PRECISION NOT NULL DEFAULT 0.98,
    s_corrected       DOUBLE PRECISION NOT NULL DEFAULT 0,

    -- 측정 결과
    primary_value     DOUBLE PRECISION NOT NULL DEFAULT 0,
    unit              VARCHAR(20) NOT NULL DEFAULT '',
    confidence        DOUBLE PRECISION NOT NULL DEFAULT 0,
    evidence_status   VARCHAR(32) NOT NULL DEFAULT 'unknown',
    diagnostic_ready  BOOLEAN NOT NULL DEFAULT FALSE,
    evidence_gaps     TEXT[] NOT NULL DEFAULT '{}',

    -- 핑거프린트 벡터 (88→896→1792 Stage별 확장)
    fingerprint_dim   INT NOT NULL DEFAULT 0,
    fingerprint_vector FLOAT4[],

    -- AI 추론 결과
    anomaly_score     DOUBLE PRECISION,
    classification    JSONB,
    concentrations    JSONB,
    model_id          VARCHAR(64),
    model_hash        VARCHAR(64),

    -- 디지털 트윈 동기화
    twin_predicted    DOUBLE PRECISION[],
    twin_residual     DOUBLE PRECISION[],
    drift_score       DOUBLE PRECISION,

    -- 데이터 무결성 (해시 체인)
    hash_chain        VARCHAR(64),
    data_integrity    VARCHAR(64),

    -- FHIR 호환 (HL7 Observation 리소스)
    fhir_observation_id VARCHAR(64),

    -- 환경 메타데이터
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

-- 핑거프린트 벡터 인덱스 (GIN for array containment queries)
CREATE INDEX IF NOT EXISTS idx_mdata_model ON measurement_data(model_id) WHERE model_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_mdata_anomaly ON measurement_data(anomaly_score) WHERE anomaly_score IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_mdata_fhir ON measurement_data(fhir_observation_id) WHERE fhir_observation_id IS NOT NULL;

-- 측정 결과 요약 뷰 (자주 사용되는 조회 최적화)
CREATE OR REPLACE VIEW measurement_summary AS
SELECT
    ms.id AS session_id,
    md.cartridge_type,
    md.primary_value,
    md.unit,
    md.confidence,
    md.evidence_status,
    md.diagnostic_ready,
    md.evidence_gaps,
    md.anomaly_score,
    md.classification,
    md.drift_score,
    md.time AS measured_at,
    ms.user_id
FROM measurement_sessions ms
JOIN measurement_data md ON ms.id = md.session_id
WHERE ms.status = 'completed';
