-- Cartridge System Database Initialization (무한확장 체계)
-- 카트리지 카테고리·타입 레지스트리 + 등급별 접근 정책

-- ============================================================================
-- 카트리지 카테고리 (무한 확장)
-- ============================================================================
CREATE TABLE IF NOT EXISTS cartridge_categories (
    code              SMALLINT PRIMARY KEY,           -- 0x01~0xFF
    name_en           VARCHAR(100) NOT NULL,
    name_ko           VARCHAR(100) NOT NULL,
    description       TEXT DEFAULT '',
    icon_url          VARCHAR(500) DEFAULT '',
    sort_order        INTEGER DEFAULT 0,
    is_active         BOOLEAN DEFAULT TRUE,
    phase             INTEGER DEFAULT 1,
    created_at        TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 기본 카테고리 초기 데이터
INSERT INTO cartridge_categories (code, name_en, name_ko, description, sort_order, phase) VALUES
    (1,   'HealthBiomarker',   '건강 바이오마커',    '혈액/타액/체액 기반 건강 바이오마커', 1, 1),
    (2,   'Environmental',     '환경 모니터링',      '수질/공기/방사능 환경 분석',         2, 1),
    (3,   'FoodSafety',        '식품 안전',         '농약/신선도/알레르겐 식품 분석',     3, 1),
    (4,   'ElectronicSensor',  '전자코/전자혀',      '전자코/전자혀/EHD 센서 분석',       4, 1),
    (5,   'AdvancedAnalysis',  '고급 분석',         '비표적/다중패널 고급 분석',          5, 2),
    (6,   'Industrial',        '산업용',            '화학물질/중금속/유해가스 산업 분석',  6, 3),
    (7,   'Veterinary',        '수의학',            '동물 혈액/바이오마커 분석',          7, 3),
    (8,   'Pharmaceutical',    '제약',              '약물 성분/농도 분석',              8, 3),
    (9,   'Agricultural',      '농업',              '토양/비료/작물 분석',              9, 4),
    (10,  'Cosmetic',          '화장품',            '성분/피부 타입 분석',              10, 4),
    (11,  'Forensic',          '법의학',            '체액/약물/독물 분석',              11, 4),
    (12,  'Marine',            '해양',              '해수/양식장/선박 연료 분석',        12, 4),
    (254, 'Beta',              '베타/실험용',        '베타 테스트 및 실험용 카트리지',     99, 2),
    (255, 'CustomResearch',    '맞춤형 연구',        '연구 목적 맞춤형 카트리지',         100, 1)
ON CONFLICT (code) DO NOTHING;

-- ============================================================================
-- 카트리지 타입 (무한 확장, 카테고리당 최대 256종)
-- ============================================================================
CREATE TABLE IF NOT EXISTS cartridge_types (
    category_code     SMALLINT NOT NULL REFERENCES cartridge_categories(code),
    type_index        SMALLINT NOT NULL,              -- 0x01~0xFF
    legacy_code       SMALLINT DEFAULT 0,             -- v1.0 호환 코드 (0이면 없음)
    name_en           VARCHAR(100) NOT NULL,
    name_ko           VARCHAR(100) NOT NULL,
    description       TEXT DEFAULT '',
    required_channels INTEGER NOT NULL DEFAULT 88,
    measurement_secs  INTEGER NOT NULL DEFAULT 15,
    unit              VARCHAR(30) DEFAULT '',
    reference_range   VARCHAR(100) DEFAULT '',
    is_active         BOOLEAN DEFAULT TRUE,
    is_beta           BOOLEAN DEFAULT FALSE,
    phase             INTEGER DEFAULT 1,
    manufacturer      VARCHAR(200) DEFAULT 'ManPaSik',
    sdk_vendor_id     VARCHAR(100) DEFAULT '',
    created_at        TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (category_code, type_index)
);

-- 기본 29종 카트리지 타입 초기 데이터
-- HealthBiomarker (14종)
INSERT INTO cartridge_types (category_code, type_index, legacy_code, name_en, name_ko, required_channels, measurement_secs, unit, reference_range) VALUES
    (1, 1,  1,  'Glucose',       '혈당',          88, 15, 'mg/dL',   '70-100'),
    (1, 2,  2,  'LipidPanel',    '지질 패널',      88, 15, 'mg/dL',   ''),
    (1, 3,  3,  'HbA1c',         '당화혈색소',     88, 15, '%',       '4.0-5.6'),
    (1, 4,  4,  'UricAcid',      '요산',          88, 15, 'mg/dL',   '3.5-7.2'),
    (1, 5,  5,  'Creatinine',    '크레아티닌',     88, 15, 'mg/dL',   '0.7-1.3'),
    (1, 6,  6,  'VitaminD',      '비타민 D',      88, 15, 'ng/mL',   '30-100'),
    (1, 7,  7,  'VitaminB12',    '비타민 B12',    88, 15, 'pg/mL',   '200-900'),
    (1, 8,  8,  'Ferritin',      '철분(페리틴)',   88, 15, 'ng/mL',   '12-300'),
    (1, 9,  9,  'Tsh',           '갑상선(TSH)',   88, 15, 'mIU/L',   '0.4-4.0'),
    (1, 10, 10, 'Cortisol',      '코르티솔',      88, 15, 'μg/dL',   '6-23'),
    (1, 11, 11, 'Testosterone',  '테스토스테론',   88, 15, 'ng/dL',   '300-1000'),
    (1, 12, 12, 'Estrogen',      '에스트로겐',    88, 15, 'pg/mL',   '15-350'),
    (1, 13, 13, 'Crp',           'C-반응성단백',  88, 15, 'mg/L',    '0-3'),
    (1, 14, 14, 'Insulin',       '인슐린',        88, 15, 'μIU/mL',  '2.6-24.9'),
    -- Environmental (4종)
    (2, 1, 32, 'WaterQuality',     '수질 검사',     88, 15, '', ''),
    (2, 2, 33, 'IndoorAirQuality', '실내 공기질',   88, 15, '', ''),
    (2, 3, 34, 'Radon',            '라돈',          88, 15, 'Bq/m³', ''),
    (2, 4, 35, 'Radiation',        '방사능',        88, 15, 'μSv/h', ''),
    -- FoodSafety (4종)
    (3, 1, 48, 'PesticideResidue', '농약 잔류',     88, 15, '', ''),
    (3, 2, 49, 'FoodFreshness',    '식품 신선도',   88, 15, '', ''),
    (3, 3, 50, 'Allergen',         '알레르겐',      88, 15, '', ''),
    (3, 4, 51, 'DateDrug',         '데이트약물',    88, 15, '', ''),
    -- ElectronicSensor (3종)
    (4, 1, 64, 'ENose',   '전자코',     8,  30, '', ''),
    (4, 2, 65, 'ETongue', '전자혀',     8,  30, '', ''),
    (4, 3, 66, 'EhdGas',  'EHD 기체',   8,  30, '', ''),
    -- AdvancedAnalysis (4종, NonTarget1792 = Phase 5 궁극 확장)
    (5, 1, 80, 'NonTarget448',    '비표적 448차원',          448,  60,  '', ''),
    (5, 2, 81, 'NonTarget896',    '비표적 896차원',          896,  90,  '', ''),
    (5, 3, 82, 'NonTarget1792',   '비표적 1792차원(궁극)',   1792, 180, '', ''),
    (5, 4, 83, 'MultiBiomarker',  '다중 바이오마커',          88,   45,  '', ''),
    -- CustomResearch (1종)
    (255, 1, 255, 'CustomResearch', '맞춤형 연구용', 896, 90, '', '')
ON CONFLICT (category_code, type_index) DO NOTHING;

CREATE INDEX IF NOT EXISTS idx_cartridge_types_legacy ON cartridge_types (legacy_code) WHERE legacy_code > 0;
CREATE INDEX IF NOT EXISTS idx_cartridge_types_active ON cartridge_types (is_active);

-- ============================================================================
-- 등급별 카트리지 접근 정책 (동적 관리)
-- ============================================================================
CREATE TABLE IF NOT EXISTS cartridge_tier_access (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tier              INTEGER NOT NULL,                -- 0: Free, 1: Basic, 2: Pro, 3: Clinical
    category_code     SMALLINT NOT NULL DEFAULT 0,     -- 0 = 전체 카테고리
    type_index        SMALLINT NOT NULL DEFAULT 0,     -- 0 = 카테고리 내 전체
    access_level      VARCHAR(20) NOT NULL DEFAULT 'restricted',
                      -- included, limited, add_on, restricted, beta
    daily_limit       INTEGER DEFAULT 0,               -- 0 = 무제한
    monthly_limit     INTEGER DEFAULT 0,               -- 0 = 무제한
    addon_price_krw   INTEGER DEFAULT 0,               -- add_on 시 건당 가격
    priority          INTEGER DEFAULT 0,               -- 높을수록 우선
    is_active         BOOLEAN DEFAULT TRUE,
    effective_from    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    effective_until   TIMESTAMPTZ,                     -- NULL = 무기한
    created_at        TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (tier, category_code, type_index)
);

-- 기본 정책 초기 데이터
-- Free 등급: 기본 건강 3종 LIMITED(일3회), 나머지 전체 RESTRICTED
INSERT INTO cartridge_tier_access (tier, category_code, type_index, access_level, priority) VALUES
    -- Free: 글로벌 기본값 = RESTRICTED
    (0, 0, 0, 'restricted', 0),
    -- Free: HealthBiomarker 기본 3종 = LIMITED (일 3회)
    (0, 1, 1, 'limited', 10),   -- Glucose
    (0, 1, 2, 'limited', 10),   -- LipidPanel
    (0, 1, 3, 'limited', 10),   -- HbA1c
    -- Basic: 글로벌 기본값 = RESTRICTED
    (1, 0, 0, 'restricted', 0),
    -- Basic: HealthBiomarker 전체 = INCLUDED
    (1, 1, 0, 'included', 5),
    -- Basic: Environmental = ADD_ON
    (1, 2, 0, 'add_on', 5),
    -- Basic: FoodSafety = ADD_ON
    (1, 3, 0, 'add_on', 5),
    -- Pro: 글로벌 기본값 = RESTRICTED
    (2, 0, 0, 'restricted', 0),
    -- Pro: HealthBiomarker = INCLUDED
    (2, 1, 0, 'included', 5),
    -- Pro: Environmental = INCLUDED
    (2, 2, 0, 'included', 5),
    -- Pro: FoodSafety = INCLUDED
    (2, 3, 0, 'included', 5),
    -- Pro: ElectronicSensor = INCLUDED
    (2, 4, 0, 'included', 5),
    -- Pro: AdvancedAnalysis = ADD_ON
    (2, 5, 0, 'add_on', 5),
    -- Pro: Veterinary = ADD_ON
    (2, 7, 0, 'add_on', 5),
    -- Pro: Agricultural = ADD_ON
    (2, 9, 0, 'add_on', 5),
    -- Pro: Cosmetic = INCLUDED
    (2, 10, 0, 'included', 5),
    -- Pro: Marine = ADD_ON
    (2, 12, 0, 'add_on', 5),
    -- Pro: ThirdParty = ADD_ON (0xF0~0xFD → code 240~253)
    -- Clinical: 글로벌 기본값 = INCLUDED (전체 무제한)
    (3, 0, 0, 'included', 0),
    -- Clinical: Beta = BETA
    (3, 254, 0, 'beta', 10)
ON CONFLICT (tier, category_code, type_index) DO NOTHING;

-- Free 등급 일일 제한 설정
UPDATE cartridge_tier_access SET daily_limit = 3
WHERE tier = 0 AND access_level = 'limited';

CREATE INDEX IF NOT EXISTS idx_cta_tier ON cartridge_tier_access (tier);
CREATE INDEX IF NOT EXISTS idx_cta_category ON cartridge_tier_access (category_code);
CREATE INDEX IF NOT EXISTS idx_cta_active ON cartridge_tier_access (is_active, effective_from, effective_until);

-- ============================================================================
-- 애드온 구매 내역
-- ============================================================================
CREATE TABLE IF NOT EXISTS cartridge_addon_purchases (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id           UUID NOT NULL,
    category_code     SMALLINT NOT NULL,
    type_index        SMALLINT NOT NULL,
    remaining_uses    INTEGER NOT NULL DEFAULT 0,
    total_purchased   INTEGER NOT NULL DEFAULT 0,
    price_krw         INTEGER NOT NULL DEFAULT 0,
    purchased_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    expires_at        TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_addon_user ON cartridge_addon_purchases (user_id);
CREATE INDEX IF NOT EXISTS idx_addon_cartridge ON cartridge_addon_purchases (category_code, type_index);

-- ============================================================================
-- 카트리지 사용 로그 (감사 추적 + 사용량 추적)
-- ============================================================================
CREATE TABLE IF NOT EXISTS cartridge_usage_log (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id           UUID NOT NULL,
    session_id        UUID NOT NULL,
    cartridge_uid     VARCHAR(100) NOT NULL,
    category_code     SMALLINT NOT NULL,
    type_index        SMALLINT NOT NULL,
    tier_at_usage     INTEGER NOT NULL,
    access_level      VARCHAR(20) NOT NULL,
    used_at           TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_usage_user ON cartridge_usage_log (user_id);
CREATE INDEX IF NOT EXISTS idx_usage_session ON cartridge_usage_log (session_id);
CREATE INDEX IF NOT EXISTS idx_usage_date ON cartridge_usage_log (used_at);
CREATE INDEX IF NOT EXISTS idx_usage_user_date ON cartridge_usage_log (user_id, used_at);
