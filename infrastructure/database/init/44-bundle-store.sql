-- Phase AW: HL7/FHIR Bundle 영속화 저장소
--
-- 외부 LIS/EHR 가 송신한 HL7 v2 메시지를 FHIR Bundle 로 변환한 후 보관.
-- 의료 표준 정보 모델 단위 (Bundle 트랜잭션) 유지로 감사/추적성 확보.
--
-- 인덱스:
--   - (tenant_id, bundle_id) PK 복합 키 — tenant 격리 + 단일 조회
--   - (tenant_id, stored_at DESC) — tenant 별 최신순 목록
--   - (tenant_id, patient_id, stored_at DESC) — 환자별 검사 결과 조회
--
-- 멀티테넌트: tenant_id 빈 문자열 ("") 은 단일 테넌트 모드 (TENANCY_ENFORCED=false).
-- payload 는 FHIR Bundle JSON 전체.

CREATE TABLE IF NOT EXISTS bundle_store (
    tenant_id   TEXT        NOT NULL DEFAULT '',
    bundle_id   TEXT        NOT NULL,
    patient_id  TEXT        NOT NULL DEFAULT '',
    payload     JSONB       NOT NULL,
    stored_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (tenant_id, bundle_id)
);

CREATE INDEX IF NOT EXISTS idx_bundle_store_tenant_latest
    ON bundle_store (tenant_id, stored_at DESC);

CREATE INDEX IF NOT EXISTS idx_bundle_store_tenant_patient_latest
    ON bundle_store (tenant_id, patient_id, stored_at DESC)
    WHERE patient_id <> '';

COMMENT ON TABLE bundle_store IS 'HL7 v2 외부 검사 결과 FHIR Bundle 보관 (Phase AW)';
COMMENT ON COLUMN bundle_store.tenant_id IS 'tenant 격리 키 (빈 문자열 = 단일 테넌트)';
COMMENT ON COLUMN bundle_store.bundle_id IS '"hl7v2-{MSH-10}" 형식';
COMMENT ON COLUMN bundle_store.patient_id IS 'PID-3 (없으면 빈 문자열)';
COMMENT ON COLUMN bundle_store.payload IS 'FHIR Bundle 전체 JSONB';
