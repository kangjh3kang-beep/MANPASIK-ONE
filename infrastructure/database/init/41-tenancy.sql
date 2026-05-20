-- 41-tenancy.sql
-- 멀티테넌트 RBAC 영속화 스키마
--
-- 사용자가 여러 조직(병원/연구소/가족 그룹)에 속할 수 있고, 조직마다 다른 역할을
-- 가질 수 있다. 시스템 전역 역할(roles 테이블)과는 별도로 운영된다.
--
-- backend/shared/tenancy/ 의 MembershipStore 인터페이스를 영속 구현.

CREATE TABLE IF NOT EXISTS tenant_memberships (
    user_id     VARCHAR(64)   NOT NULL,
    tenant_id   VARCHAR(64)   NOT NULL,
    role        VARCHAR(32)   NOT NULL,
    -- TenantRole 값: owner | admin | medical_staff | member | viewer
    active      BOOLEAN       NOT NULL DEFAULT TRUE,
    joined_at   BIGINT        NOT NULL,
    -- 가입 unix timestamp (초 단위)
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, tenant_id)
);

-- 인덱스: 사용자별 조직 조회
CREATE INDEX IF NOT EXISTS idx_tenant_memberships_user
    ON tenant_memberships(user_id) WHERE active = TRUE;

-- 인덱스: 조직별 멤버 조회
CREATE INDEX IF NOT EXISTS idx_tenant_memberships_tenant
    ON tenant_memberships(tenant_id) WHERE active = TRUE;

-- 역할 검증 (CHECK 제약)
ALTER TABLE tenant_memberships
    DROP CONSTRAINT IF EXISTS tenant_memberships_role_check;
ALTER TABLE tenant_memberships
    ADD CONSTRAINT tenant_memberships_role_check
    CHECK (role IN ('owner', 'admin', 'medical_staff', 'member', 'viewer'));

-- updated_at 자동 갱신 트리거
CREATE OR REPLACE FUNCTION update_tenant_memberships_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_tenant_memberships_updated_at ON tenant_memberships;
CREATE TRIGGER trg_tenant_memberships_updated_at
    BEFORE UPDATE ON tenant_memberships
    FOR EACH ROW EXECUTE FUNCTION update_tenant_memberships_updated_at();

COMMENT ON TABLE tenant_memberships IS
    '멀티테넌트 멤버십 (Phase Y-2). MembershipStore 영속 구현 대상.';

-- =========================================================
-- 초대장 (Phase AA-1)
-- 토큰 기반 초대 발급 + 수락 + 만료/취소.
-- =========================================================

CREATE TABLE IF NOT EXISTS tenant_invitations (
    token         VARCHAR(64)  PRIMARY KEY,
    tenant_id     VARCHAR(64)  NOT NULL,
    inviter_id    VARCHAR(64)  NOT NULL,
    invitee_hint  VARCHAR(255),
    role          VARCHAR(32)  NOT NULL,
    status        VARCHAR(16)  NOT NULL DEFAULT 'pending',
    issued_at     TIMESTAMPTZ  NOT NULL,
    expires_at    TIMESTAMPTZ  NOT NULL,
    accepted_by   VARCHAR(64),
    accepted_at   TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_tenant_invitations_tenant
    ON tenant_invitations(tenant_id) WHERE status = 'pending';

CREATE INDEX IF NOT EXISTS idx_tenant_invitations_inviter
    ON tenant_invitations(inviter_id) WHERE status = 'pending';

CREATE INDEX IF NOT EXISTS idx_tenant_invitations_expires
    ON tenant_invitations(expires_at) WHERE status = 'pending';

ALTER TABLE tenant_invitations
    DROP CONSTRAINT IF EXISTS tenant_invitations_status_check;
ALTER TABLE tenant_invitations
    ADD CONSTRAINT tenant_invitations_status_check
    CHECK (status IN ('pending', 'accepted', 'revoked', 'expired'));

ALTER TABLE tenant_invitations
    DROP CONSTRAINT IF EXISTS tenant_invitations_role_check;
ALTER TABLE tenant_invitations
    ADD CONSTRAINT tenant_invitations_role_check
    CHECK (role IN ('owner', 'admin', 'medical_staff', 'member', 'viewer'));

COMMENT ON TABLE tenant_invitations IS
    '테넌시 초대장 (Phase AA-1). InvitationStore 영속 구현 대상.';
