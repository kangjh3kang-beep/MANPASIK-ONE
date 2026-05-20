-- audit-service 초기화 스크립트
-- Phase A: 감사 로그 영속화

-- 감사 이벤트 테이블
CREATE TABLE IF NOT EXISTS audit_entries (
    id              VARCHAR(64) PRIMARY KEY,
    admin_id        VARCHAR(64) NOT NULL,
    action          VARCHAR(100) NOT NULL,
    resource_type   VARCHAR(50) NOT NULL,
    resource_id     VARCHAR(64),
    description     TEXT,
    ip_address      VARCHAR(45),
    user_agent      VARCHAR(500),
    metadata        JSONB DEFAULT '{}',
    timestamp       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_entries_admin ON audit_entries(admin_id);
CREATE INDEX IF NOT EXISTS idx_audit_entries_action ON audit_entries(action);
CREATE INDEX IF NOT EXISTS idx_audit_entries_resource ON audit_entries(resource_type, resource_id);
CREATE INDEX IF NOT EXISTS idx_audit_entries_timestamp ON audit_entries(timestamp DESC);
