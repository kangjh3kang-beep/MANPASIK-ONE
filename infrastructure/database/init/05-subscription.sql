-- Subscription Service Database Initialization
-- 테이블은 POSTGRES_DB(기본 manpasik)에 생성됩니다.

-- 구독 테이블
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    tier INTEGER NOT NULL DEFAULT 0,         -- 0: Free, 1: Basic, 2: Pro, 3: Clinical
    status INTEGER NOT NULL DEFAULT 1,       -- 1: Active, 2: Cancelled, 3: Expired, 4: Suspended, 5: Trial
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE,
    cancelled_at TIMESTAMP WITH TIME ZONE,
    max_devices INTEGER NOT NULL DEFAULT 1,
    max_family_members INTEGER NOT NULL DEFAULT 0,
    ai_coaching_enabled BOOLEAN DEFAULT FALSE,
    telemedicine_enabled BOOLEAN DEFAULT FALSE,
    monthly_price_krw INTEGER NOT NULL DEFAULT 0,
    auto_renew BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 사용자별 활성 구독은 하나만 허용 (unique index)
CREATE UNIQUE INDEX IF NOT EXISTS idx_subscriptions_user_id_active
    ON subscriptions (user_id) WHERE status = 1;

CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions (user_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_status ON subscriptions (status);
CREATE INDEX IF NOT EXISTS idx_subscriptions_expires_at ON subscriptions (expires_at);

-- 구독 플랜 테이블 (참조 데이터, 불변)
CREATE TABLE IF NOT EXISTS subscription_plans (
    tier INTEGER PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT DEFAULT '',
    monthly_price_krw INTEGER NOT NULL DEFAULT 0,
    max_devices INTEGER NOT NULL DEFAULT 1,
    max_family_members INTEGER NOT NULL DEFAULT 0,
    ai_coaching_enabled BOOLEAN DEFAULT FALSE,
    telemedicine_enabled BOOLEAN DEFAULT FALSE,
    features JSONB DEFAULT '[]'::jsonb,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 플랜 초기 데이터
INSERT INTO subscription_plans (tier, name, description, monthly_price_krw, max_devices, max_family_members, ai_coaching_enabled, telemedicine_enabled, features)
VALUES
    (0, 'Free', '기본 무료 플랜', 0, 1, 0, FALSE, FALSE, '["기본 측정", "측정 이력 조회", "단일 리더기"]'),
    (1, 'Basic Safety', '기본 안전 관리 플랜 (₩9,900/월)', 9900, 3, 2, FALSE, FALSE, '["기본 측정", "측정 이력 조회", "리더기 3대", "가족 2명", "데이터 내보내기"]'),
    (2, 'Bio-Optimization', 'AI 코칭 포함 프로 플랜 (₩29,900/월)', 29900, 5, 5, TRUE, FALSE, '["기본 측정", "측정 이력 조회", "리더기 5대", "가족 5명", "데이터 내보내기", "AI 건강 코칭", "트렌드 분석", "건강 점수"]'),
    (3, 'Clinical Guard', '화상진료 포함 클리니컬 플랜 (₩59,900/월)', 59900, 10, 10, TRUE, TRUE, '["기본 측정", "측정 이력 조회", "리더기 10대", "가족 10명", "데이터 내보내기", "AI 건강 코칭", "트렌드 분석", "건강 점수", "화상진료", "의료진 매칭", "FHIR 연동"]')
ON CONFLICT (tier) DO NOTHING;

-- 구독 변경 이력 테이블
CREATE TABLE IF NOT EXISTS subscription_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    old_tier INTEGER,
    new_tier INTEGER,
    action VARCHAR(50) NOT NULL,   -- create, upgrade, downgrade, cancel, renew
    reason TEXT DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_subscription_history_user_id ON subscription_history (user_id);
CREATE INDEX IF NOT EXISTS idx_subscription_history_created_at ON subscription_history (created_at);
