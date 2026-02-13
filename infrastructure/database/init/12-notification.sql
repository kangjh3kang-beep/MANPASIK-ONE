-- notification-service 초기화 스크립트
-- Phase 3A: 알림 서비스 DB

-- 알림 채널 enum
DO $$ BEGIN
    CREATE TYPE notification_channel AS ENUM ('push', 'email', 'sms', 'in_app');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- 알림 타입 enum
DO $$ BEGIN
    CREATE TYPE notification_type AS ENUM (
        'measurement', 'health_alert', 'coaching', 'subscription',
        'order', 'family', 'system', 'reminder', 'promotion'
    );
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- 알림 우선순위 enum
DO $$ BEGIN
    CREATE TYPE notification_priority AS ENUM ('low', 'normal', 'high', 'urgent');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- 알림 테이블
CREATE TABLE IF NOT EXISTS notifications (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    type            notification_type NOT NULL DEFAULT 'system',
    channel         notification_channel NOT NULL DEFAULT 'in_app',
    priority        notification_priority NOT NULL DEFAULT 'normal',
    title           VARCHAR(255) NOT NULL,
    body            TEXT,
    data            JSONB DEFAULT '{}',
    is_read         BOOLEAN NOT NULL DEFAULT FALSE,
    silent          BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    read_at         TIMESTAMPTZ
);

-- 인덱스
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_user_unread ON notifications(user_id) WHERE is_read = FALSE;
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(user_id, type);
CREATE INDEX IF NOT EXISTS idx_notifications_created ON notifications(user_id, created_at DESC);

-- 알림 설정 테이블
CREATE TABLE IF NOT EXISTS notification_preferences (
    user_id                 UUID PRIMARY KEY,
    push_enabled            BOOLEAN NOT NULL DEFAULT TRUE,
    email_enabled           BOOLEAN NOT NULL DEFAULT TRUE,
    sms_enabled             BOOLEAN NOT NULL DEFAULT FALSE,
    in_app_enabled          BOOLEAN NOT NULL DEFAULT TRUE,
    health_alert_enabled    BOOLEAN NOT NULL DEFAULT TRUE,
    coaching_enabled        BOOLEAN NOT NULL DEFAULT TRUE,
    promotion_enabled       BOOLEAN NOT NULL DEFAULT FALSE,
    quiet_hours_start       VARCHAR(5),  -- "HH:MM"
    quiet_hours_end         VARCHAR(5),
    language                VARCHAR(5) NOT NULL DEFAULT 'ko',
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 디바이스 토큰 테이블 (푸시 알림용)
CREATE TABLE IF NOT EXISTS device_tokens (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    token           TEXT NOT NULL,
    platform        VARCHAR(20) NOT NULL,  -- 'ios', 'android', 'web'
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_device_tokens_user ON device_tokens(user_id);
