-- ============================================================================
-- 카트리지 스토어 (Developer/Store/Review Service) — Phase 3
-- ============================================================================

-- 개발자 프로필
CREATE TABLE IF NOT EXISTS developers (
    id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL UNIQUE,
    company_name VARCHAR(200) DEFAULT '',
    tier VARCHAR(30) DEFAULT 'explorer',
    status VARCHAR(20) DEFAULT 'active',
    cartridge_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_developers_user ON developers(user_id);

-- 개발자 API 키
CREATE TABLE IF NOT EXISTS developer_api_keys (
    id VARCHAR(64) PRIMARY KEY,
    developer_id VARCHAR(64) NOT NULL REFERENCES developers(id) ON DELETE CASCADE,
    name VARCHAR(100) DEFAULT '',
    key_hash VARCHAR(256) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_api_keys_dev ON developer_api_keys(developer_id);

-- 스토어 리스팅
CREATE TABLE IF NOT EXISTS store_listings (
    id VARCHAR(64) PRIMARY KEY,
    developer_id VARCHAR(64) NOT NULL REFERENCES developers(id),
    cartridge_name VARCHAR(200) NOT NULL,
    category VARCHAR(50) DEFAULT '',
    description TEXT DEFAULT '',
    price_krw INT DEFAULT 0,
    rating DOUBLE PRECISION DEFAULT 0,
    review_count INT DEFAULT 0,
    download_count INT DEFAULT 0,
    icon_url VARCHAR(500) DEFAULT '',
    package_url VARCHAR(500) DEFAULT '',
    status VARCHAR(20) DEFAULT 'draft',
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_store_listings_dev ON store_listings(developer_id);
CREATE INDEX IF NOT EXISTS idx_store_listings_category ON store_listings(category);
CREATE INDEX IF NOT EXISTS idx_store_listings_status ON store_listings(status);

-- 구매 기록
CREATE TABLE IF NOT EXISTS store_purchases (
    id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL,
    listing_id VARCHAR(64) NOT NULL REFERENCES store_listings(id),
    price_krw INT DEFAULT 0,
    purchased_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_store_purchases_user ON store_purchases(user_id);

-- 카트리지 제출/심사
CREATE TABLE IF NOT EXISTS cartridge_submissions (
    id VARCHAR(64) PRIMARY KEY,
    developer_id VARCHAR(64) NOT NULL REFERENCES developers(id),
    cartridge_name VARCHAR(200) NOT NULL,
    category VARCHAR(50) DEFAULT '',
    description TEXT DEFAULT '',
    price_krw INT DEFAULT 0,
    package_url VARCHAR(500) DEFAULT '',
    status VARCHAR(20) DEFAULT 'submitted',
    reviewer_id VARCHAR(64) DEFAULT '',
    reviewer_comment TEXT DEFAULT '',
    submitted_at TIMESTAMPTZ DEFAULT NOW(),
    reviewed_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_submissions_dev ON cartridge_submissions(developer_id);
CREATE INDEX IF NOT EXISTS idx_submissions_status ON cartridge_submissions(status);

-- 사용자 리뷰/평점
CREATE TABLE IF NOT EXISTS user_cartridge_ratings (
    id VARCHAR(64) PRIMARY KEY,
    listing_id VARCHAR(64) NOT NULL REFERENCES store_listings(id),
    user_id VARCHAR(64) NOT NULL,
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(listing_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_ratings_listing ON user_cartridge_ratings(listing_id);
