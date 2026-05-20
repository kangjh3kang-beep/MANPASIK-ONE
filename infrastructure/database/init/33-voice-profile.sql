-- ============================================================================
-- 음성 프로필 (Voice Profile Service) — Phase 5
-- ============================================================================

CREATE TABLE IF NOT EXISTS voice_profiles (
    id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL,
    profile_name VARCHAR(200) DEFAULT '',
    language VARCHAR(10) DEFAULT 'ko',
    status VARCHAR(20) DEFAULT 'processing',
    model_url VARCHAR(500) DEFAULT '',
    voice_sample_url VARCHAR(500) DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_voice_profiles_user ON voice_profiles(user_id);
