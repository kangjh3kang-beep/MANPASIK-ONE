-- =============================================================================
-- ManPaSik Phase 2: coaching-service 테이블 초기화
-- =============================================================================

-- 목표 카테고리 ENUM
DO $$ BEGIN
  CREATE TYPE goal_category AS ENUM (
    'BLOOD_GLUCOSE', 'BLOOD_PRESSURE', 'CHOLESTEROL', 'WEIGHT',
    'EXERCISE', 'NUTRITION', 'SLEEP', 'STRESS', 'CUSTOM'
  );
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- 목표 상태 ENUM
DO $$ BEGIN
  CREATE TYPE goal_status AS ENUM ('ACTIVE', 'ACHIEVED', 'PAUSED', 'CANCELLED');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- 코칭 유형 ENUM
DO $$ BEGIN
  CREATE TYPE coaching_type AS ENUM (
    'MEASUREMENT_FEEDBACK', 'DAILY_TIP', 'GOAL_PROGRESS',
    'ALERT', 'MOTIVATION', 'RECOMMENDATION'
  );
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- 추천 유형 ENUM
DO $$ BEGIN
  CREATE TYPE recommendation_type AS ENUM (
    'FOOD', 'EXERCISE', 'SUPPLEMENT', 'LIFESTYLE', 'CHECKUP'
  );
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- ============================================================================
-- 건강 목표 테이블
-- ============================================================================
CREATE TABLE IF NOT EXISTS health_goals (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          VARCHAR(64) NOT NULL,
    category         goal_category NOT NULL,
    metric_name      VARCHAR(100) NOT NULL,
    target_value     DOUBLE PRECISION NOT NULL,
    current_value    DOUBLE PRECISION DEFAULT 0.0,
    unit             VARCHAR(30) DEFAULT '',
    progress_pct     DOUBLE PRECISION DEFAULT 0.0,
    status           goal_status NOT NULL DEFAULT 'ACTIVE',
    description      TEXT DEFAULT '',
    target_date      TIMESTAMPTZ,
    achieved_at      TIMESTAMPTZ,
    created_at       TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_health_goals_user ON health_goals (user_id);
CREATE INDEX IF NOT EXISTS idx_health_goals_status ON health_goals (user_id, status);
CREATE INDEX IF NOT EXISTS idx_health_goals_category ON health_goals (user_id, category);

-- ============================================================================
-- 코칭 메시지 테이블
-- ============================================================================
CREATE TABLE IF NOT EXISTS coaching_messages (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          VARCHAR(64) NOT NULL,
    coaching_type    coaching_type NOT NULL,
    title            VARCHAR(500) NOT NULL,
    body             TEXT NOT NULL,
    risk_level       INTEGER DEFAULT 0,    -- 0: unspecified, 1: low, 2: moderate, 3: high, 4: critical
    action_items     JSONB DEFAULT '[]'::jsonb,
    related_metric   VARCHAR(100) DEFAULT '',
    related_value    DOUBLE PRECISION DEFAULT 0.0,
    is_read          BOOLEAN DEFAULT FALSE,
    created_at       TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_coaching_messages_user ON coaching_messages (user_id);
CREATE INDEX IF NOT EXISTS idx_coaching_messages_type ON coaching_messages (user_id, coaching_type);
CREATE INDEX IF NOT EXISTS idx_coaching_messages_created ON coaching_messages (user_id, created_at DESC);

-- ============================================================================
-- 일일 건강 리포트 테이블
-- ============================================================================
CREATE TABLE IF NOT EXISTS daily_health_reports (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             VARCHAR(64) NOT NULL,
    report_date         DATE NOT NULL,
    overall_score       DOUBLE PRECISION DEFAULT 0.0,
    measurements_count  INTEGER DEFAULT 0,
    summary             TEXT DEFAULT '',
    recommendations     JSONB DEFAULT '[]'::jsonb,
    created_at          TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, report_date)
);

CREATE INDEX IF NOT EXISTS idx_daily_reports_user ON daily_health_reports (user_id);
CREATE INDEX IF NOT EXISTS idx_daily_reports_date ON daily_health_reports (user_id, report_date DESC);

-- ============================================================================
-- 주간 건강 리포트 테이블
-- ============================================================================
CREATE TABLE IF NOT EXISTS weekly_health_reports (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             VARCHAR(64) NOT NULL,
    week_start          DATE NOT NULL,
    week_end            DATE NOT NULL,
    average_score       DOUBLE PRECISION DEFAULT 0.0,
    score_trend         VARCHAR(20) DEFAULT 'stable',  -- improving, stable, declining
    total_measurements  INTEGER DEFAULT 0,
    goals_achieved      INTEGER DEFAULT 0,
    goals_active        INTEGER DEFAULT 0,
    weekly_summary      TEXT DEFAULT '',
    key_insights        JSONB DEFAULT '[]'::jsonb,
    created_at          TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, week_start)
);

CREATE INDEX IF NOT EXISTS idx_weekly_reports_user ON weekly_health_reports (user_id);
CREATE INDEX IF NOT EXISTS idx_weekly_reports_week ON weekly_health_reports (user_id, week_start DESC);

-- ============================================================================
-- 개인화 추천 테이블
-- ============================================================================
CREATE TABLE IF NOT EXISTS recommendations (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             VARCHAR(64) NOT NULL,
    recommendation_type recommendation_type NOT NULL,
    title               VARCHAR(500) NOT NULL,
    description         TEXT NOT NULL,
    reason              TEXT DEFAULT '',
    priority            INTEGER DEFAULT 1,   -- 1: low, 2: moderate, 3: high, 4: critical
    action_steps        JSONB DEFAULT '[]'::jsonb,
    related_metric      VARCHAR(100) DEFAULT '',
    is_dismissed        BOOLEAN DEFAULT FALSE,
    created_at          TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    expires_at          TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_recommendations_user ON recommendations (user_id);
CREATE INDEX IF NOT EXISTS idx_recommendations_type ON recommendations (user_id, recommendation_type);
CREATE INDEX IF NOT EXISTS idx_recommendations_active ON recommendations (user_id, is_dismissed, created_at DESC);
