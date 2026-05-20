CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Users Table (사용자 기본 정보)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    nickname VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ -- 오프라인 CRDT 동기화를 위한 Soft Delete 처리
);

-- 2. Devices Table (기기/리더기 관리 - STM32 하드웨어 연동 대상)
CREATE TABLE devices (
    mac_address VARCHAR(17) PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    hardware_version VARCHAR(50) NOT NULL, -- ex) CSI v1.0
    firmware_version VARCHAR(50) NOT NULL,
    status VARCHAR(20) DEFAULT 'ACTIVE', -- ACTIVE, INACTIVE, REPAIR_MODE
    last_sync_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 3. Measurements Table (핵심 측정 데이터 보관)
-- ⚠️ 향후 분당/초당 대량 데이터 인입 시 TimescaleDB Hypertable로 승격 예정
CREATE TABLE measurements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    device_mac VARCHAR(17) REFERENCES devices(mac_address) ON DELETE SET NULL,
    measured_at TIMESTAMPTZ NOT NULL,
    
    -- Layer 3: Rust Core Engine 계산 결과물 (차분 신호 및 특징 벡터)
    diff_signal JSONB NOT NULL,     -- S_n - α_n * R_n (차동 측정값 배열)
    fingerprint JSONB NOT NULL,     -- 896차원 융합 특징 벡터
    
    -- Layer 6: AI 추론 결과
    health_score INTEGER NOT NULL,
    risk_label VARCHAR(50) NOT NULL,
    
    -- Sync Metadata (모바일 앱 연동 및 오프라인-퍼스트 지원)
    client_local_id VARCHAR(100) UNIQUE, -- 앱 내부 데이터베이스(Isar)의 로컬 식별자
    synced_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 인덱스 추가 (빠른 조회를 위함)
CREATE INDEX idx_measurements_user_id ON measurements(user_id);
CREATE INDEX idx_measurements_measured_at ON measurements(measured_at);
CREATE INDEX idx_measurements_client_local_id ON measurements(client_local_id);

-- 4. Diagnostics Log Table (자가치유 및 진단 로그 - 펌웨어 에러 추적)
CREATE TABLE diagnostic_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_mac VARCHAR(17) REFERENCES devices(mac_address) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL, -- ex) 'HARDWARE_FAULT', 'FIRMWARE_ROLLBACK'
    severity VARCHAR(20) NOT NULL,   -- 'INFO', 'WARNING', 'CRITICAL'
    payload JSONB,                   -- 에러 상세 로그 또는 복구 이벤트 내역
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_diagnostic_logs_device ON diagnostic_logs(device_mac);
