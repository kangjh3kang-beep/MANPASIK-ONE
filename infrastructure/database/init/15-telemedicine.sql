-- =============================================================================
-- Telemedicine Service Schema (Phase 3B)
-- =============================================================================

-- 의사 전문분야
CREATE TYPE doctor_specialty AS ENUM (
    'general', 'internal', 'cardiology', 'endocrinology',
    'dermatology', 'pediatrics', 'psychiatry', 'orthopedics',
    'ophthalmology', 'ent'
);

-- 진료 상태
CREATE TYPE consultation_status AS ENUM (
    'requested', 'matched', 'scheduled', 'in_progress',
    'completed', 'cancelled', 'no_show'
);

-- 화상 세션 상태
CREATE TYPE video_session_status AS ENUM (
    'waiting', 'connected', 'ended', 'failed'
);

-- 의사 프로필
CREATE TABLE IF NOT EXISTS doctors (
    doctor_id       VARCHAR(36) PRIMARY KEY,
    user_id         VARCHAR(36),
    name            VARCHAR(100) NOT NULL,
    specialty       doctor_specialty NOT NULL DEFAULT 'general',
    hospital        VARCHAR(200),
    license_number  VARCHAR(50) UNIQUE NOT NULL,
    experience_years INTEGER DEFAULT 0,
    rating          DECIMAL(3,2) DEFAULT 0.00,
    total_consultations INTEGER DEFAULT 0,
    is_available    BOOLEAN DEFAULT true,
    languages       TEXT[] DEFAULT ARRAY['ko'],
    profile_image_url TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 진료 기록
CREATE TABLE IF NOT EXISTS consultations (
    consultation_id VARCHAR(36) PRIMARY KEY,
    patient_user_id VARCHAR(36) NOT NULL,
    doctor_id       VARCHAR(36) REFERENCES doctors(doctor_id),
    specialty       doctor_specialty NOT NULL DEFAULT 'general',
    chief_complaint TEXT NOT NULL,
    description     TEXT,
    status          consultation_status NOT NULL DEFAULT 'requested',
    diagnosis       TEXT,
    doctor_notes    TEXT,
    prescription_id VARCHAR(36),
    duration_minutes INTEGER DEFAULT 0,
    rating          DECIMAL(3,2) DEFAULT 0.00,
    is_urgent       BOOLEAN DEFAULT false,
    scheduled_at    TIMESTAMPTZ,
    started_at      TIMESTAMPTZ,
    ended_at        TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 화상 세션
CREATE TABLE IF NOT EXISTS video_sessions (
    session_id      VARCHAR(36) PRIMARY KEY,
    consultation_id VARCHAR(36) NOT NULL REFERENCES consultations(consultation_id),
    room_url        TEXT NOT NULL,
    token           TEXT NOT NULL,
    status          video_session_status NOT NULL DEFAULT 'waiting',
    started_at      TIMESTAMPTZ,
    ended_at        TIMESTAMPTZ,
    duration_seconds INTEGER DEFAULT 0,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 인덱스
CREATE INDEX IF NOT EXISTS idx_consultations_patient ON consultations(patient_user_id);
CREATE INDEX IF NOT EXISTS idx_consultations_doctor ON consultations(doctor_id);
CREATE INDEX IF NOT EXISTS idx_consultations_status ON consultations(status);
CREATE INDEX IF NOT EXISTS idx_doctors_specialty ON doctors(specialty);
CREATE INDEX IF NOT EXISTS idx_doctors_available ON doctors(is_available);
CREATE INDEX IF NOT EXISTS idx_video_sessions_consultation ON video_sessions(consultation_id);
CREATE INDEX IF NOT EXISTS idx_consultations_created ON consultations(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_consultations_patient_status ON consultations(patient_user_id, status);
