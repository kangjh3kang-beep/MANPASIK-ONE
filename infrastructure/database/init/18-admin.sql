-- =============================================================================
-- Admin Service Schema (Phase 3B)
-- =============================================================================

-- 관리자 역할
CREATE TYPE admin_role AS ENUM (
    'super_admin', 'admin', 'moderator', 'support', 'analyst'
);

-- 감사 로그 액션
CREATE TYPE audit_action AS ENUM (
    'login', 'logout', 'create', 'update', 'delete',
    'config_change', 'user_ban', 'user_unban', 'role_change'
);

-- 관리자
CREATE TABLE IF NOT EXISTS admin_users (
    admin_id        VARCHAR(36) PRIMARY KEY,
    user_id         VARCHAR(36) NOT NULL UNIQUE,
    email           VARCHAR(255) NOT NULL UNIQUE,
    display_name    VARCHAR(100) NOT NULL,
    role            admin_role NOT NULL DEFAULT 'support',
    is_active       BOOLEAN DEFAULT true,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    last_login_at   TIMESTAMPTZ,
    created_by      VARCHAR(36),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 감사 로그
CREATE TABLE IF NOT EXISTS audit_logs (
    entry_id        VARCHAR(36) PRIMARY KEY,
    admin_id        VARCHAR(36) NOT NULL REFERENCES admin_users(admin_id),
    admin_email     VARCHAR(255),
    action          audit_action NOT NULL,
    resource_type   VARCHAR(100),
    resource_id     VARCHAR(36),
    description     TEXT,
    ip_address      VARCHAR(45),
    metadata        JSONB,
    timestamp       TIMESTAMPTZ DEFAULT NOW()
);

-- 시스템 설정
CREATE TABLE IF NOT EXISTS system_configs (
    key             VARCHAR(200) PRIMARY KEY,
    value           TEXT NOT NULL,
    description     TEXT,
    updated_by      VARCHAR(36),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 인덱스
CREATE INDEX IF NOT EXISTS idx_admin_users_role ON admin_users(role);
CREATE INDEX IF NOT EXISTS idx_admin_users_active ON admin_users(is_active);
CREATE INDEX IF NOT EXISTS idx_audit_logs_admin ON audit_logs(admin_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_timestamp ON audit_logs(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_resource ON audit_logs(resource_type, resource_id);

-- 기본 시스템 설정 삽입
INSERT INTO system_configs (key, value, description) VALUES
    ('maintenance_mode', 'false', '시스템 점검 모드'),
    ('max_devices_per_user', '5', '사용자당 최대 디바이스 수'),
    ('default_language', 'ko', '기본 언어'),
    ('session_timeout_minutes', '30', '세션 타임아웃 (분)'),
    ('max_file_upload_mb', '50', '최대 파일 업로드 크기 (MB)')
ON CONFLICT (key) DO NOTHING;
