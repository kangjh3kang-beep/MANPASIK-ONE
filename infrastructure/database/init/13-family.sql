-- family-service 초기화 스크립트
-- Phase 3A: 가족 그룹 서비스 DB

-- 가족 역할 enum
DO $$ BEGIN
    CREATE TYPE family_role AS ENUM ('owner', 'guardian', 'member', 'child', 'elderly');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- 초대 상태 enum
DO $$ BEGIN
    CREATE TYPE invitation_status AS ENUM ('pending', 'accepted', 'declined', 'expired');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- 가족 그룹 테이블
CREATE TABLE IF NOT EXISTS family_groups (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_user_id   UUID NOT NULL,
    group_name      VARCHAR(100) NOT NULL,
    description     TEXT,
    max_members     INT NOT NULL DEFAULT 10,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_family_groups_owner ON family_groups(owner_user_id);

-- 가족 멤버 테이블
CREATE TABLE IF NOT EXISTS family_members (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    group_id        UUID NOT NULL REFERENCES family_groups(id) ON DELETE CASCADE,
    display_name    VARCHAR(100),
    email           VARCHAR(255),
    role            family_role NOT NULL DEFAULT 'member',
    sharing_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    joined_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, group_id)
);

CREATE INDEX IF NOT EXISTS idx_family_members_group ON family_members(group_id);
CREATE INDEX IF NOT EXISTS idx_family_members_user ON family_members(user_id);

-- 초대 테이블
CREATE TABLE IF NOT EXISTS family_invitations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id        UUID NOT NULL REFERENCES family_groups(id) ON DELETE CASCADE,
    inviter_user_id UUID NOT NULL,
    invitee_email   VARCHAR(255) NOT NULL,
    role            family_role NOT NULL DEFAULT 'member',
    message         TEXT,
    status          invitation_status NOT NULL DEFAULT 'pending',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at      TIMESTAMPTZ NOT NULL DEFAULT (NOW() + INTERVAL '7 days')
);

CREATE INDEX IF NOT EXISTS idx_family_invitations_group ON family_invitations(group_id);
CREATE INDEX IF NOT EXISTS idx_family_invitations_invitee ON family_invitations(invitee_email);

-- 건강 데이터 공유 설정 테이블
CREATE TABLE IF NOT EXISTS sharing_preferences (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL,
    group_id            UUID NOT NULL REFERENCES family_groups(id) ON DELETE CASCADE,
    share_measurements  BOOLEAN NOT NULL DEFAULT FALSE,
    share_health_score  BOOLEAN NOT NULL DEFAULT FALSE,
    share_goals         BOOLEAN NOT NULL DEFAULT FALSE,
    share_coaching      BOOLEAN NOT NULL DEFAULT FALSE,
    share_alerts        BOOLEAN NOT NULL DEFAULT FALSE,
    allowed_viewer_ids  UUID[] DEFAULT '{}',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, group_id)
);

CREATE INDEX IF NOT EXISTS idx_sharing_prefs_user_group ON sharing_preferences(user_id, group_id);
