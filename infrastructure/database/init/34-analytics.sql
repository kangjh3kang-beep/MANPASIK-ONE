-- analytics-service 초기화 스크립트
-- Phase A: 분석 이벤트 영속화

-- 분석 이벤트 테이블
CREATE TABLE IF NOT EXISTS analytics_events (
    id              VARCHAR(64) PRIMARY KEY,
    user_id         VARCHAR(64) NOT NULL,
    event_type      VARCHAR(100) NOT NULL,
    properties      JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_analytics_events_user ON analytics_events(user_id);
CREATE INDEX IF NOT EXISTS idx_analytics_events_type ON analytics_events(event_type);
CREATE INDEX IF NOT EXISTS idx_analytics_events_created ON analytics_events(created_at DESC);

-- 일별 통계 캐시 테이블
CREATE TABLE IF NOT EXISTS analytics_daily_stats (
    date            DATE PRIMARY KEY,
    total_events    INT NOT NULL DEFAULT 0,
    unique_users    INT NOT NULL DEFAULT 0,
    events_by_type  JSONB DEFAULT '{}'
);
