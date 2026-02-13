-- =============================================================================
-- Reservation Service Schema (Phase 3B)
-- =============================================================================

-- 시설 유형
CREATE TYPE facility_type AS ENUM (
    'hospital', 'clinic', 'pharmacy', 'lab', 'dental', 'oriental'
);

-- 예약 상태
CREATE TYPE reservation_status AS ENUM (
    'pending', 'confirmed', 'checked_in', 'completed', 'cancelled', 'no_show'
);

-- 의료 시설
CREATE TABLE IF NOT EXISTS facilities (
    facility_id     VARCHAR(36) PRIMARY KEY,
    name            VARCHAR(200) NOT NULL,
    type            facility_type NOT NULL DEFAULT 'hospital',
    address         TEXT NOT NULL,
    phone           VARCHAR(20),
    latitude        DECIMAL(10,7),
    longitude       DECIMAL(10,7),
    rating          DECIMAL(3,2) DEFAULT 0.00,
    review_count    INTEGER DEFAULT 0,
    specialties     TEXT[],
    operating_hours TEXT,
    is_open_now     BOOLEAN DEFAULT false,
    accepts_reservation BOOLEAN DEFAULT true,
    image_url       TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 예약 가능 시간대
CREATE TABLE IF NOT EXISTS time_slots (
    slot_id         VARCHAR(36) PRIMARY KEY,
    facility_id     VARCHAR(36) NOT NULL REFERENCES facilities(facility_id),
    doctor_id       VARCHAR(36),
    doctor_name     VARCHAR(100),
    slot_date       DATE NOT NULL,
    start_time      TIMESTAMPTZ NOT NULL,
    end_time        TIMESTAMPTZ NOT NULL,
    is_available    BOOLEAN DEFAULT true,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 예약
CREATE TABLE IF NOT EXISTS reservations (
    reservation_id  VARCHAR(36) PRIMARY KEY,
    user_id         VARCHAR(36) NOT NULL,
    facility_id     VARCHAR(36) NOT NULL REFERENCES facilities(facility_id),
    facility_name   VARCHAR(200),
    doctor_id       VARCHAR(36),
    doctor_name     VARCHAR(100),
    specialty       VARCHAR(50),
    status          reservation_status NOT NULL DEFAULT 'pending',
    reason          TEXT,
    notes           TEXT,
    scheduled_at    TIMESTAMPTZ NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 인덱스
CREATE INDEX IF NOT EXISTS idx_facilities_type ON facilities(type);
CREATE INDEX IF NOT EXISTS idx_facilities_location ON facilities(latitude, longitude);
CREATE INDEX IF NOT EXISTS idx_time_slots_facility_date ON time_slots(facility_id, slot_date);
CREATE INDEX IF NOT EXISTS idx_time_slots_available ON time_slots(is_available);
CREATE INDEX IF NOT EXISTS idx_reservations_user ON reservations(user_id);
CREATE INDEX IF NOT EXISTS idx_reservations_facility ON reservations(facility_id);
CREATE INDEX IF NOT EXISTS idx_reservations_status ON reservations(status);
CREATE INDEX IF NOT EXISTS idx_reservations_scheduled ON reservations(scheduled_at);
