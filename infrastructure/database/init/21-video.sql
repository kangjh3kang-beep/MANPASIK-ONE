-- video-service 데이터베이스 스키마
-- Phase 3C: WebRTC 시그널링·미디어

-- ENUM 타입
DO $$ BEGIN
    CREATE TYPE room_status AS ENUM (
        'waiting', 'active', 'ended', 'failed'
    );
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    CREATE TYPE room_type AS ENUM (
        'one_to_one', 'group', 'webinar', 'consultation'
    );
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    CREATE TYPE signal_type AS ENUM (
        'offer', 'answer', 'ice_candidate', 'renegotiate', 'mute', 'unmute'
    );
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- 회의실 테이블
CREATE TABLE IF NOT EXISTS video_rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    room_type room_type NOT NULL DEFAULT 'one_to_one',
    status room_status NOT NULL DEFAULT 'waiting',
    created_by UUID NOT NULL,
    max_participants INT DEFAULT 2,
    recording_url TEXT,
    duration_seconds INT DEFAULT 0,
    total_bytes_transferred BIGINT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_rooms_created_by ON video_rooms(created_by);
CREATE INDEX IF NOT EXISTS idx_rooms_status ON video_rooms(status);

-- 참가자 테이블
CREATE TABLE IF NOT EXISTS video_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES video_rooms(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    display_name VARCHAR(100),
    role VARCHAR(50) DEFAULT 'participant',
    is_audio_enabled BOOLEAN DEFAULT TRUE,
    is_video_enabled BOOLEAN DEFAULT TRUE,
    is_screen_sharing BOOLEAN DEFAULT FALSE,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    left_at TIMESTAMPTZ,
    UNIQUE(room_id, user_id, joined_at)
);

CREATE INDEX IF NOT EXISTS idx_participants_room ON video_participants(room_id);
CREATE INDEX IF NOT EXISTS idx_participants_user ON video_participants(user_id);

-- WebRTC 시그널 테이블
CREATE TABLE IF NOT EXISTS video_signals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES video_rooms(id) ON DELETE CASCADE,
    from_user_id UUID NOT NULL,
    to_user_id UUID,
    signal_type signal_type NOT NULL,
    payload JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_signals_room ON video_signals(room_id);
CREATE INDEX IF NOT EXISTS idx_signals_created ON video_signals(created_at);

-- 회의실 통계 테이블
CREATE TABLE IF NOT EXISTS video_room_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES video_rooms(id) ON DELETE CASCADE UNIQUE,
    total_participants_joined INT DEFAULT 0,
    peak_participants INT DEFAULT 0,
    total_signals INT DEFAULT 0,
    average_latency_ms REAL DEFAULT 0,
    total_bytes_transferred BIGINT DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_room_stats_room ON video_room_stats(room_id);
