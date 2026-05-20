-- ============================================================================
-- 식단/칼로리 (Vision/Diet Service) — Phase 2
-- ============================================================================

-- 음식 분석 결과
CREATE TABLE IF NOT EXISTS food_analyses (
    id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL,
    meal_type VARCHAR(20) DEFAULT 'snack',
    image_url VARCHAR(500) DEFAULT '',
    total_calories DOUBLE PRECISION DEFAULT 0,
    total_carbs_g DOUBLE PRECISION DEFAULT 0,
    total_protein_g DOUBLE PRECISION DEFAULT 0,
    total_fat_g DOUBLE PRECISION DEFAULT 0,
    analyzed_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_food_analyses_user ON food_analyses(user_id);

-- 음식 항목
CREATE TABLE IF NOT EXISTS food_items (
    id VARCHAR(64) PRIMARY KEY,
    analysis_id VARCHAR(64) NOT NULL REFERENCES food_analyses(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    calories DOUBLE PRECISION DEFAULT 0,
    portion_g DOUBLE PRECISION DEFAULT 0,
    carbs_g DOUBLE PRECISION DEFAULT 0,
    protein_g DOUBLE PRECISION DEFAULT 0,
    fat_g DOUBLE PRECISION DEFAULT 0,
    confidence DOUBLE PRECISION DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_food_items_analysis ON food_items(analysis_id);

-- 식사 로그
CREATE TABLE IF NOT EXISTS meal_logs (
    id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL,
    meal_type VARCHAR(20) NOT NULL,
    description TEXT DEFAULT '',
    total_calories DOUBLE PRECISION DEFAULT 0,
    analysis_id VARCHAR(64) DEFAULT '',
    logged_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_meal_logs_user ON meal_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_meal_logs_date ON meal_logs(logged_at);

-- 일일 영양 요약
CREATE TABLE IF NOT EXISTS daily_nutrition_summaries (
    id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL,
    date DATE NOT NULL,
    total_calories DOUBLE PRECISION DEFAULT 0,
    total_carbs_g DOUBLE PRECISION DEFAULT 0,
    total_protein_g DOUBLE PRECISION DEFAULT 0,
    total_fat_g DOUBLE PRECISION DEFAULT 0,
    meal_count INT DEFAULT 0,
    UNIQUE(user_id, date)
);
CREATE INDEX IF NOT EXISTS idx_daily_nutrition_user_date ON daily_nutrition_summaries(user_id, date);
