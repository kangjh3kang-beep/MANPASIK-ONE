-- ============================================================================
-- 카트리지 수익/분석 (Revenue/Analytics Service) — Phase 5
-- ============================================================================

-- 수익 기록
CREATE TABLE IF NOT EXISTS revenue_records (
    id VARCHAR(64) PRIMARY KEY,
    developer_id VARCHAR(64) NOT NULL,
    listing_id VARCHAR(64) NOT NULL,
    amount_krw BIGINT DEFAULT 0,
    platform_fee_krw BIGINT DEFAULT 0,
    developer_earning_krw BIGINT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_revenue_dev ON revenue_records(developer_id);

-- 정산 요청
CREATE TABLE IF NOT EXISTS payout_requests (
    id VARCHAR(64) PRIMARY KEY,
    developer_id VARCHAR(64) NOT NULL,
    amount_krw BIGINT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'pending',
    requested_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_payout_dev ON payout_requests(developer_id);

-- 수익 배분 설정
CREATE TABLE IF NOT EXISTS revenue_split_configs (
    id VARCHAR(64) PRIMARY KEY,
    developer_id VARCHAR(64) NOT NULL UNIQUE,
    developer_pct INT DEFAULT 70,
    platform_pct INT DEFAULT 30,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 카트리지 사용 분석
CREATE TABLE IF NOT EXISTS cartridge_usage_analytics (
    id VARCHAR(64) PRIMARY KEY,
    listing_id VARCHAR(64) NOT NULL,
    date DATE NOT NULL,
    usage_count INT DEFAULT 0,
    unique_users INT DEFAULT 0,
    UNIQUE(listing_id, date)
);
CREATE INDEX IF NOT EXISTS idx_usage_analytics_listing ON cartridge_usage_analytics(listing_id);
