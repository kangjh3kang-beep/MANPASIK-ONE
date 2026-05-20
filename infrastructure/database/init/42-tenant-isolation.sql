-- 42-tenant-isolation.sql
-- Phase AG-1: 멀티테넌트 데이터 격리 적용
--
-- 기존 테이블에 tenant_id 컬럼 추가하여 SELECT 쿼리가 tenant 단위로 격리되도록.
-- legacy 데이터 (NULL) 는 호환을 위해 허용하되, 신규 데이터는 NOT NULL 권장.

-- =========================================================
-- measurement_data: tenant 격리 컬럼 + 인덱스
-- =========================================================

ALTER TABLE measurement_data
    ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(64);

-- 사용자별 + 테넌트별 빠른 조회를 위한 복합 인덱스
CREATE INDEX IF NOT EXISTS idx_measurement_data_user_tenant
    ON measurement_data(user_id, tenant_id, time DESC);

-- 테넌트별 조회용
CREATE INDEX IF NOT EXISTS idx_measurement_data_tenant_time
    ON measurement_data(tenant_id, time DESC) WHERE tenant_id IS NOT NULL;

COMMENT ON COLUMN measurement_data.tenant_id IS
    '멀티테넌트 격리 (Phase AG-1). NULL=legacy/personal, 값 있으면 해당 조직 데이터.';

-- =========================================================
-- health_records: 동일 적용
-- =========================================================

ALTER TABLE IF EXISTS health_records
    ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(64);

CREATE INDEX IF NOT EXISTS idx_health_records_user_tenant
    ON health_records(user_id, tenant_id);

-- =========================================================
-- prescriptions: 동일 적용 (의료 데이터 격리)
-- =========================================================

ALTER TABLE IF EXISTS prescriptions
    ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(64);

CREATE INDEX IF NOT EXISTS idx_prescriptions_tenant
    ON prescriptions(tenant_id) WHERE tenant_id IS NOT NULL;
