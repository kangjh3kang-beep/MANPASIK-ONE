-- ManPaSik user-service 데이터베이스 초기화
-- 사용자 프로필, 구독, 가족 그룹 테이블

-- 프로필 테이블
CREATE TABLE IF NOT EXISTS user_profiles (
    user_id         UUID PRIMARY KEY,
    email           VARCHAR(320) NOT NULL,
    display_name    VARCHAR(100) DEFAULT '',
    avatar_url      TEXT DEFAULT '',
    language        VARCHAR(5) DEFAULT 'ko',
    timezone        VARCHAR(50) DEFAULT 'Asia/Seoul',
    subscription_tier VARCHAR(20) DEFAULT 'free',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_profiles_email ON user_profiles(email);

-- 구독 테이블
CREATE TABLE IF NOT EXISTS subscriptions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES user_profiles(user_id) ON DELETE CASCADE,
    tier            VARCHAR(20) NOT NULL DEFAULT 'free',
    started_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at      TIMESTAMPTZ,
    max_devices     INT NOT NULL DEFAULT 1,
    max_family_members INT NOT NULL DEFAULT 0,
    ai_coaching_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    telemedicine_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_tier ON subscriptions(tier);

-- 가족 그룹 테이블
CREATE TABLE IF NOT EXISTS family_groups (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100) NOT NULL,
    owner_id        UUID NOT NULL REFERENCES user_profiles(user_id) ON DELETE CASCADE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_family_groups_owner ON family_groups(owner_id);

-- 가족 구성원 테이블
CREATE TABLE IF NOT EXISTS family_members (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id        UUID NOT NULL REFERENCES family_groups(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES user_profiles(user_id) ON DELETE CASCADE,
    role            VARCHAR(20) NOT NULL DEFAULT 'adult',
    joined_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(group_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_family_members_group ON family_members(group_id);
CREATE INDEX IF NOT EXISTS idx_family_members_user ON family_members(user_id);
