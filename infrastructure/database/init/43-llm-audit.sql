-- 43-llm-audit.sql
-- Phase AI-1: LLM 호출 감사 로그 영속화
--
-- 모든 LLM 호출에 대한 tenant 별 추적/비용 분리/규정 준수 감사.
-- HIPAA/GDPR 7년 보존 정책 (tenancy/audit-service 와 함께 운영).

CREATE TABLE IF NOT EXISTS llm_audit (
    id                BIGSERIAL    PRIMARY KEY,
    tenant_id         VARCHAR(64)  NOT NULL,        -- "personal" 포함 (NULL 대신)
    user_id           VARCHAR(64),
    provider          VARCHAR(32),
    model             VARCHAR(64),
    prompt_tokens     INTEGER      NOT NULL DEFAULT 0,
    completion_tokens INTEGER      NOT NULL DEFAULT 0,
    total_tokens      INTEGER      NOT NULL DEFAULT 0,
    latency_ms        BIGINT       NOT NULL DEFAULT 0,
    success           BOOLEAN      NOT NULL DEFAULT FALSE,
    error_message     TEXT,
    recorded_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- tenant 별 시간순 조회 (운영 통계)
CREATE INDEX IF NOT EXISTS idx_llm_audit_tenant_time
    ON llm_audit(tenant_id, recorded_at DESC);

-- 실패 호출 빠른 조회 (트러블슈팅)
CREATE INDEX IF NOT EXISTS idx_llm_audit_failures
    ON llm_audit(tenant_id, recorded_at DESC) WHERE success = FALSE;

-- 사용자별 조회
CREATE INDEX IF NOT EXISTS idx_llm_audit_user
    ON llm_audit(user_id, recorded_at DESC) WHERE user_id IS NOT NULL;

-- =========================================================
-- llm_quota: tenant 별 동적 LLM 한도 (Phase AJ-3)
-- =========================================================

CREATE TABLE IF NOT EXISTS llm_quota (
    tenant_id            VARCHAR(64)   PRIMARY KEY,
    daily_token_limit    INTEGER       NOT NULL DEFAULT 0,
    monthly_token_limit  INTEGER       NOT NULL DEFAULT 0,
    daily_request_limit  INTEGER       NOT NULL DEFAULT 0,
    updated_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    created_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE llm_quota IS
    'Tenant 별 LLM 사용 한도 (Phase AJ-3). DynamicQuota 의 영속 backing store.';

COMMENT ON TABLE llm_audit IS
    'LLM 호출 감사 로그 (Phase AI-1). tenant 별 비용/사용량 추적.';
