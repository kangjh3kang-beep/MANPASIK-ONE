-- ============================================================================
-- 데이터 플랫폼 (LocationStats/DataProvision Service) — Phase 4
-- ============================================================================

-- 지역 통계 (배치 집계)
CREATE TABLE IF NOT EXISTS region_stats (
    id VARCHAR(64) PRIMARY KEY,
    region_code VARCHAR(20) NOT NULL,
    region_name VARCHAR(100) DEFAULT '',
    period VARCHAR(20) NOT NULL,
    biomarker VARCHAR(50) NOT NULL,
    avg_value DOUBLE PRECISION DEFAULT 0,
    min_value DOUBLE PRECISION DEFAULT 0,
    max_value DOUBLE PRECISION DEFAULT 0,
    sample_count INT DEFAULT 0,
    risk_level VARCHAR(20) DEFAULT 'low',
    computed_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_region_stats_code ON region_stats(region_code);
CREATE INDEX IF NOT EXISTS idx_region_stats_biomarker ON region_stats(biomarker);

-- 익명 집계
CREATE TABLE IF NOT EXISTS anonymous_aggregations (
    id VARCHAR(64) PRIMARY KEY,
    batch_date DATE NOT NULL,
    biomarker VARCHAR(50) NOT NULL,
    cohort_key VARCHAR(100) DEFAULT '',
    record_count INT DEFAULT 0,
    avg_value DOUBLE PRECISION DEFAULT 0,
    status VARCHAR(20) DEFAULT 'completed',
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_anon_agg_date ON anonymous_aggregations(batch_date);

-- 집계 인사이트
CREATE TABLE IF NOT EXISTS aggregated_insights (
    id VARCHAR(64) PRIMARY KEY,
    aggregation_id VARCHAR(64) REFERENCES anonymous_aggregations(id),
    biomarker VARCHAR(50) NOT NULL,
    insight_type VARCHAR(50) NOT NULL,
    summary TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 데이터 접근 요청
CREATE TABLE IF NOT EXISTS data_access_requests (
    id VARCHAR(64) PRIMARY KEY,
    requester_id VARCHAR(64) NOT NULL,
    dataset_id VARCHAR(64) NOT NULL,
    purpose TEXT DEFAULT '',
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_data_access_requester ON data_access_requests(requester_id);

-- 익명화 데이터셋
CREATE TABLE IF NOT EXISTS anonymized_datasets (
    id VARCHAR(64) PRIMARY KEY,
    dataset_type VARCHAR(50) NOT NULL,
    description TEXT DEFAULT '',
    record_count INT DEFAULT 0,
    date_range VARCHAR(50) DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW()
);
