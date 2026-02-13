-- =============================================================================
-- Data sharing consents (GDPR/개인정보 동의)
-- 23: 데이터 공유·처리 동의 저장
-- =============================================================================

DO $$ BEGIN
    CREATE TYPE consent_scope AS ENUM (
        'health_record', 'prescription', 'measurement', 'telemedicine',
        'marketing', 'analytics', 'third_party_share', 'research'
    );
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    CREATE TYPE consent_status AS ENUM ('granted', 'revoked', 'expired');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS data_sharing_consents (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    scope           consent_scope NOT NULL,
    status          consent_status NOT NULL DEFAULT 'granted',
    version         INTEGER NOT NULL DEFAULT 1,
    granted_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at      TIMESTAMPTZ,
    expires_at       TIMESTAMPTZ,
    ip_address      INET,
    user_agent      TEXT,
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_consents_user ON data_sharing_consents(user_id);
CREATE INDEX IF NOT EXISTS idx_consents_user_scope ON data_sharing_consents(user_id, scope);
CREATE INDEX IF NOT EXISTS idx_consents_status ON data_sharing_consents(status);
CREATE INDEX IF NOT EXISTS idx_consents_granted_at ON data_sharing_consents(granted_at);
