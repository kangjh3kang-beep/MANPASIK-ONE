-- =============================================================================
-- ManPaSik Phase 2: calibration-service 테이블 초기화
-- =============================================================================

-- 보정 유형 ENUM
DO $$ BEGIN
  CREATE TYPE calibration_type AS ENUM ('FACTORY', 'FIELD', 'AUTO');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- 보정 상태 ENUM
DO $$ BEGIN
  CREATE TYPE calibration_status AS ENUM ('VALID', 'EXPIRING', 'EXPIRED', 'NEEDED');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- ============================================================================
-- 보정 기록 테이블
-- ============================================================================
CREATE TABLE IF NOT EXISTS calibration_records (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id           VARCHAR(100) NOT NULL,
    cartridge_category  SMALLINT NOT NULL,
    cartridge_type_index SMALLINT NOT NULL,
    calibration_type    calibration_type NOT NULL DEFAULT 'FACTORY',
    alpha               DOUBLE PRECISION NOT NULL DEFAULT 0.95,
    channel_offsets     DOUBLE PRECISION[] DEFAULT '{}',
    channel_gains       DOUBLE PRECISION[] DEFAULT '{}',
    temp_coefficient    DOUBLE PRECISION DEFAULT 0.0,
    humidity_coefficient DOUBLE PRECISION DEFAULT 0.0,
    accuracy_score      DOUBLE PRECISION DEFAULT 0.0,   -- 0.0 ~ 1.0
    reference_standard  VARCHAR(200) DEFAULT '',
    calibrated_by       VARCHAR(200) DEFAULT '',
    calibrated_at       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at          TIMESTAMPTZ NOT NULL,
    status              calibration_status NOT NULL DEFAULT 'VALID',
    created_at          TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 디바이스 + 카트리지 타입별 최신 보정 빠른 조회
CREATE INDEX IF NOT EXISTS idx_calibration_device_type
    ON calibration_records (device_id, cartridge_category, cartridge_type_index, calibrated_at DESC);

CREATE INDEX IF NOT EXISTS idx_calibration_device
    ON calibration_records (device_id);

CREATE INDEX IF NOT EXISTS idx_calibration_expires
    ON calibration_records (expires_at);

CREATE INDEX IF NOT EXISTS idx_calibration_status
    ON calibration_records (status);

-- ============================================================================
-- 보정 모델 테이블 (카트리지 타입별 기본 보정 파라미터)
-- ============================================================================
CREATE TABLE IF NOT EXISTS calibration_models (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cartridge_category  SMALLINT NOT NULL,
    cartridge_type_index SMALLINT NOT NULL,
    name                VARCHAR(200) NOT NULL,
    version             VARCHAR(50) NOT NULL DEFAULT '1.0.0',
    default_alpha       DOUBLE PRECISION NOT NULL DEFAULT 0.95,
    validity_days       INTEGER NOT NULL DEFAULT 90,
    description         TEXT DEFAULT '',
    is_active           BOOLEAN DEFAULT TRUE,
    created_at          TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (cartridge_category, cartridge_type_index)
);

-- 기본 보정 모델 초기 데이터
INSERT INTO calibration_models (cartridge_category, cartridge_type_index, name, version, default_alpha, validity_days, description) VALUES
    -- HealthBiomarker (카테고리 1): alpha=0.95, 90일
    (1, 1,  'Glucose Calibration',       '1.0.0', 0.95, 90, '혈당 측정 보정 모델'),
    (1, 2,  'LipidPanel Calibration',    '1.0.0', 0.95, 90, '지질 패널 보정 모델'),
    (1, 3,  'HbA1c Calibration',         '1.0.0', 0.95, 90, '당화혈색소 보정 모델'),
    (1, 4,  'UricAcid Calibration',      '1.0.0', 0.95, 90, '요산 보정 모델'),
    (1, 5,  'Creatinine Calibration',    '1.0.0', 0.95, 90, '크레아티닌 보정 모델'),
    (1, 6,  'VitaminD Calibration',      '1.0.0', 0.95, 90, '비타민D 보정 모델'),
    (1, 7,  'VitaminB12 Calibration',    '1.0.0', 0.95, 90, '비타민B12 보정 모델'),
    (1, 8,  'Ferritin Calibration',      '1.0.0', 0.95, 90, '페리틴 보정 모델'),
    (1, 9,  'TSH Calibration',           '1.0.0', 0.95, 90, '갑상선자극호르몬 보정 모델'),
    (1, 10, 'Cortisol Calibration',      '1.0.0', 0.95, 90, '코르티솔 보정 모델'),
    (1, 11, 'Testosterone Calibration',  '1.0.0', 0.95, 90, '테스토스테론 보정 모델'),
    (1, 12, 'Estrogen Calibration',      '1.0.0', 0.95, 90, '에스트로겐 보정 모델'),
    (1, 13, 'CRP Calibration',           '1.0.0', 0.95, 90, 'C반응성단백질 보정 모델'),
    (1, 14, 'Insulin Calibration',       '1.0.0', 0.95, 90, '인슐린 보정 모델'),
    -- ElectronicSensor (카테고리 4): alpha=0.92, 60일
    (4, 1,  'ENose Calibration',         '1.0.0', 0.92, 60, '전자코 보정 모델'),
    (4, 2,  'ETongue Calibration',       '1.0.0', 0.92, 60, '전자혀 보정 모델'),
    (4, 3,  'EHD Gas Calibration',       '1.0.0', 0.92, 60, 'EHD 가스 보정 모델'),
    -- AdvancedAnalysis (카테고리 5): alpha=0.97, 120일
    (5, 1,  'NonTarget448 Calibration',  '1.0.0', 0.97, 120, '비표적 448차원 보정 모델'),
    (5, 2,  'NonTarget896 Calibration',  '1.0.0', 0.97, 120, '비표적 896차원 보정 모델'),
    (5, 3,  'NonTarget1792 Calibration', '1.0.0', 0.97, 120, '비표적 1792차원 궁극 확장 보정 모델'),
    (5, 4,  'MultiBiomarker Calibration','1.0.0', 0.97, 120, '다중 바이오마커 보정 모델'),
    -- CustomResearch (카테고리 255): alpha=0.95, 30일
    (255, 1, 'CustomResearch Calibration', '1.0.0', 0.95, 30, '맞춤형 연구용 보정 모델')
ON CONFLICT (cartridge_category, cartridge_type_index) DO NOTHING;
