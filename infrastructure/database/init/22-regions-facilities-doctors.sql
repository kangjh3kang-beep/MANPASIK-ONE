-- =============================================================================
-- Regions, Facilities (extensions), Doctors — reservation-service
-- 22: 지역·시설·의사 테이블 (16-reservation.sql 시설 보강)
-- =============================================================================

-- 지역(국가/광역/구군) 계층
CREATE TABLE IF NOT EXISTS regions (
    id              VARCHAR(64) PRIMARY KEY,
    country_code    VARCHAR(2) NOT NULL,
    region_code     VARCHAR(32) NOT NULL DEFAULT '',
    district_code   VARCHAR(32) NOT NULL DEFAULT '',
    name_en         VARCHAR(200) NOT NULL DEFAULT '',
    name_local      VARCHAR(200) NOT NULL DEFAULT '',
    parent_id       VARCHAR(64),
    timezone        VARCHAR(64) NOT NULL DEFAULT 'UTC',
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_regions_country ON regions(country_code);
CREATE INDEX IF NOT EXISTS idx_regions_country_region ON regions(country_code, region_code);
CREATE INDEX IF NOT EXISTS idx_regions_parent ON regions(parent_id);

-- facilities 테이블에 지역·화상진료 컬럼 추가 (16에서 이미 facilities 생성된 경우)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'facilities' AND column_name = 'country_code') THEN
        ALTER TABLE facilities ADD COLUMN country_code VARCHAR(2) DEFAULT '';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'facilities' AND column_name = 'region_code') THEN
        ALTER TABLE facilities ADD COLUMN region_code VARCHAR(32) DEFAULT '';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'facilities' AND column_name = 'district_code') THEN
        ALTER TABLE facilities ADD COLUMN district_code VARCHAR(32) DEFAULT '';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'facilities' AND column_name = 'timezone') THEN
        ALTER TABLE facilities ADD COLUMN timezone VARCHAR(64) DEFAULT 'UTC';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'facilities' AND column_name = 'has_telemedicine') THEN
        ALTER TABLE facilities ADD COLUMN has_telemedicine BOOLEAN DEFAULT false;
    END IF;
EXCEPTION
    WHEN undefined_table THEN NULL;  -- facilities가 아직 없으면 무시 (16이 나중에 로드되는 경우)
END $$;

-- 의사 (16-reservation.sql 적용 후 실행)
CREATE TABLE IF NOT EXISTS doctors (
    id                   VARCHAR(36) PRIMARY KEY,
    facility_id           VARCHAR(36) NOT NULL REFERENCES facilities(facility_id) ON DELETE CASCADE,
    user_id              UUID,
    name                 VARCHAR(200) NOT NULL,
    specialty            VARCHAR(50) NOT NULL DEFAULT '',
    license_number       VARCHAR(50) DEFAULT '',
    languages            TEXT[] DEFAULT '{}',
    is_available         BOOLEAN DEFAULT true,
    rating               DECIMAL(3,2) DEFAULT 0.00,
    total_consultations   INTEGER DEFAULT 0,
    created_at           TIMESTAMPTZ DEFAULT NOW(),
    updated_at           TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_doctors_facility ON doctors(facility_id);
CREATE INDEX IF NOT EXISTS idx_doctors_specialty ON doctors(facility_id, specialty);
CREATE INDEX IF NOT EXISTS idx_doctors_available ON doctors(is_available);
