-- marketplace-service 초기화 스크립트
-- Phase A: 마켓플레이스 영속화

-- 파트너 상태 enum
DO $$ BEGIN
    CREATE TYPE partner_status AS ENUM ('pending', 'active', 'suspended');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- 파트너 테이블
CREATE TABLE IF NOT EXISTS marketplace_partners (
    id              VARCHAR(64) PRIMARY KEY,
    name            VARCHAR(200) NOT NULL,
    description     TEXT,
    contact_email   VARCHAR(200),
    status          partner_status NOT NULL DEFAULT 'pending',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 파트너 제품 테이블
CREATE TABLE IF NOT EXISTS marketplace_products (
    id              VARCHAR(64) PRIMARY KEY,
    partner_id      VARCHAR(64) NOT NULL REFERENCES marketplace_partners(id),
    name            VARCHAR(200) NOT NULL,
    description     TEXT,
    price           DOUBLE PRECISION NOT NULL DEFAULT 0,
    category        VARCHAR(50),
    image_url       VARCHAR(500),
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_mp_products_partner ON marketplace_products(partner_id);
CREATE INDEX IF NOT EXISTS idx_mp_products_category ON marketplace_products(category);
CREATE INDEX IF NOT EXISTS idx_mp_products_active ON marketplace_products(is_active) WHERE is_active = TRUE;

-- 파트너 통계 테이블
CREATE TABLE IF NOT EXISTS marketplace_stats (
    partner_id      VARCHAR(64) PRIMARY KEY REFERENCES marketplace_partners(id),
    total_products  INT NOT NULL DEFAULT 0,
    total_orders    INT NOT NULL DEFAULT 0,
    revenue         DOUBLE PRECISION NOT NULL DEFAULT 0
);
