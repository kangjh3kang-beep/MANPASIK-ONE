-- health-record-service 초기화 스크립트
-- Phase 3A: 건강 기록 서비스 DB

-- 건강 기록 유형 enum
DO $$ BEGIN
    CREATE TYPE health_record_type AS ENUM (
        'measurement', 'medication', 'symptom', 'vital_sign', 'lab_result',
        'allergy', 'condition', 'immunization', 'procedure', 'note'
    );
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- FHIR 리소스 유형 enum
DO $$ BEGIN
    CREATE TYPE fhir_resource_type AS ENUM (
        'observation', 'condition', 'medication_statement', 'allergy_intolerance',
        'immunization', 'procedure', 'diagnostic_report', 'patient'
    );
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- 건강 기록 테이블
CREATE TABLE IF NOT EXISTS health_records (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL,
    record_type         health_record_type NOT NULL,
    title               VARCHAR(255) NOT NULL,
    description         TEXT,
    data                JSONB DEFAULT '{}',
    source              VARCHAR(50) NOT NULL DEFAULT 'manual',  -- 'manpasik', 'manual', 'fhir_import'
    fhir_resource_id    VARCHAR(100),
    fhir_type           fhir_resource_type,
    recorded_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ   -- 소프트 삭제
);

-- 인덱스
CREATE INDEX IF NOT EXISTS idx_health_records_user ON health_records(user_id);
CREATE INDEX IF NOT EXISTS idx_health_records_user_type ON health_records(user_id, record_type);
CREATE INDEX IF NOT EXISTS idx_health_records_recorded ON health_records(user_id, recorded_at DESC);
CREATE INDEX IF NOT EXISTS idx_health_records_fhir ON health_records(fhir_resource_id) WHERE fhir_resource_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_health_records_source ON health_records(source);

-- 건강 기록 첨부파일 테이블
CREATE TABLE IF NOT EXISTS health_record_attachments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    record_id       UUID NOT NULL REFERENCES health_records(id) ON DELETE CASCADE,
    file_name       VARCHAR(255) NOT NULL,
    file_type       VARCHAR(50) NOT NULL,  -- 'image/png', 'application/pdf' 등
    file_size       BIGINT NOT NULL,
    storage_path    TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_hr_attachments_record ON health_record_attachments(record_id);

-- 건강 기록 공유 이력 테이블 (FHIR 전송 등)
CREATE TABLE IF NOT EXISTS health_record_shares (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    record_id       UUID NOT NULL REFERENCES health_records(id) ON DELETE CASCADE,
    shared_with     VARCHAR(255) NOT NULL,  -- 의료기관, 사용자 등
    share_type      VARCHAR(50) NOT NULL,   -- 'fhir_export', 'family_share', 'provider_share'
    fhir_bundle_id  VARCHAR(100),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_hr_shares_record ON health_record_shares(record_id);
