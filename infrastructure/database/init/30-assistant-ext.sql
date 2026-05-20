-- ============================================================================
-- AI 비서 확장 (Assistant Service Phase 5 고도화)
-- ============================================================================

-- 사용자 선호도 (다중 턴 메모리)
CREATE TABLE IF NOT EXISTS assistant_user_preferences (
    id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL,
    key VARCHAR(100) NOT NULL,
    value TEXT NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, key)
);
CREATE INDEX IF NOT EXISTS idx_asst_prefs_user ON assistant_user_preferences(user_id);
