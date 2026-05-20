-- emergency-service 초기화 스크립트
-- Phase A: 응급 데이터 영속화

-- 응급 상태 enum
DO $$ BEGIN
    CREATE TYPE emergency_status AS ENUM ('reported', 'dispatched', 'resolved', 'cancelled');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- 응급 신고 테이블
CREATE TABLE IF NOT EXISTS emergencies (
    id              VARCHAR(64) PRIMARY KEY,
    user_id         VARCHAR(64) NOT NULL,
    type            VARCHAR(50) NOT NULL,
    location        VARCHAR(200),
    description     TEXT,
    status          emergency_status NOT NULL DEFAULT 'reported',
    contact_ids     TEXT[] DEFAULT '{}',
    timestamp       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_emergencies_user ON emergencies(user_id);
CREATE INDEX IF NOT EXISTS idx_emergencies_status ON emergencies(status);
CREATE INDEX IF NOT EXISTS idx_emergencies_timestamp ON emergencies(timestamp DESC);

-- 비상 연락처 테이블
CREATE TABLE IF NOT EXISTS emergency_contacts (
    id              VARCHAR(64) PRIMARY KEY,
    user_id         VARCHAR(64) NOT NULL,
    name            VARCHAR(100) NOT NULL,
    phone           VARCHAR(20) NOT NULL,
    relationship    VARCHAR(50)
);

CREATE INDEX IF NOT EXISTS idx_emergency_contacts_user ON emergency_contacts(user_id);

-- 응급 설정 테이블
CREATE TABLE IF NOT EXISTS emergency_settings (
    user_id                 VARCHAR(64) PRIMARY KEY,
    auto_call_119           BOOLEAN NOT NULL DEFAULT FALSE,
    emergency_contact_ids   TEXT[] DEFAULT '{}',
    medical_info            TEXT DEFAULT ''
);
