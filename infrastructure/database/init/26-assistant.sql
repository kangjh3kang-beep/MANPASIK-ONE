-- ============================================================================
-- AI 비서 (Assistant Service) — Phase 2
-- ============================================================================

-- 비서 세션
CREATE TABLE IF NOT EXISTS assistant_sessions (
    id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL,
    title VARCHAR(200) DEFAULT '',
    status VARCHAR(20) DEFAULT 'active',
    turn_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_assistant_sessions_user ON assistant_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_assistant_sessions_status ON assistant_sessions(status);

-- 비서 턴(대화)
CREATE TABLE IF NOT EXISTS assistant_turns (
    id VARCHAR(64) PRIMARY KEY,
    session_id VARCHAR(64) NOT NULL REFERENCES assistant_sessions(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL CHECK (role IN ('user', 'assistant', 'system')),
    content TEXT NOT NULL,
    intent VARCHAR(100) DEFAULT '',
    action_type VARCHAR(100) DEFAULT '',
    action_result TEXT DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_assistant_turns_session ON assistant_turns(session_id);
CREATE INDEX IF NOT EXISTS idx_assistant_turns_intent ON assistant_turns(intent);
