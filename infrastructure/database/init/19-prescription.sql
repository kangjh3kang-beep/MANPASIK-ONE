-- prescription-service 데이터베이스 스키마
-- Phase 3C: 처방전 관리

-- ENUM 타입
DO $$ BEGIN
    CREATE TYPE prescription_status AS ENUM (
        'draft', 'active', 'dispensed', 'completed', 'cancelled', 'expired'
    );
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    CREATE TYPE drug_interaction_severity AS ENUM (
        'none', 'minor', 'moderate', 'major', 'contraindicated'
    );
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- 처방전 테이블
CREATE TABLE IF NOT EXISTS prescriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    patient_user_id UUID NOT NULL,
    doctor_id UUID NOT NULL,
    consultation_id UUID,
    status prescription_status NOT NULL DEFAULT 'active',
    diagnosis TEXT,
    notes TEXT,
    pharmacy_id UUID,
    prescribed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    dispensed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_prescriptions_patient ON prescriptions(patient_user_id);
CREATE INDEX IF NOT EXISTS idx_prescriptions_doctor ON prescriptions(doctor_id);
CREATE INDEX IF NOT EXISTS idx_prescriptions_status ON prescriptions(status);
CREATE INDEX IF NOT EXISTS idx_prescriptions_consultation ON prescriptions(consultation_id);

-- 처방 약물 테이블
CREATE TABLE IF NOT EXISTS prescription_medications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    prescription_id UUID NOT NULL REFERENCES prescriptions(id) ON DELETE CASCADE,
    drug_name VARCHAR(200) NOT NULL,
    drug_code VARCHAR(50),
    dosage VARCHAR(100),
    frequency VARCHAR(100),
    duration_days INT DEFAULT 0,
    route VARCHAR(50) DEFAULT '경구',
    instructions TEXT,
    quantity INT DEFAULT 0,
    refills_remaining INT DEFAULT 0,
    is_generic_allowed BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_medications_prescription ON prescription_medications(prescription_id);

-- 약물 상호작용 테이블
CREATE TABLE IF NOT EXISTS drug_interactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    drug_a_code VARCHAR(50) NOT NULL,
    drug_b_code VARCHAR(50) NOT NULL,
    severity drug_interaction_severity NOT NULL DEFAULT 'none',
    description TEXT,
    recommendation TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(drug_a_code, drug_b_code)
);

CREATE INDEX IF NOT EXISTS idx_interactions_drugs ON drug_interactions(drug_a_code, drug_b_code);

-- 복약 기록 테이블
CREATE TABLE IF NOT EXISTS medication_taken_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    prescription_id UUID NOT NULL REFERENCES prescriptions(id),
    medication_id UUID NOT NULL REFERENCES prescription_medications(id),
    patient_user_id UUID NOT NULL,
    taken_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    skipped BOOLEAN DEFAULT FALSE,
    notes TEXT
);

CREATE INDEX IF NOT EXISTS idx_taken_logs_patient ON medication_taken_logs(patient_user_id, taken_at);

-- 시드 데이터: 주요 약물 상호작용
INSERT INTO drug_interactions (drug_a_code, drug_b_code, severity, description, recommendation) VALUES
    ('WARF001', 'ASPR001', 'major', '와파린과 아스피린의 병용은 출혈 위험을 크게 증가시킵니다', 'INR 모니터링 강화, 출혈 징후 관찰'),
    ('WARF001', 'IBUP001', 'major', '와파린과 이부프로펜은 출혈 위험을 증가시킵니다', 'NSAIDs 대신 아세트아미노펜 사용 고려'),
    ('SSRI001', 'MAOI001', 'contraindicated', 'SSRI와 MAOI 병용은 세로토닌 증후군을 유발할 수 있습니다', '절대 병용 금지'),
    ('SIMV001', 'CLAR001', 'major', '심바스타틴과 클래리스로마이신 병용 시 횡문근융해증 위험', '클래리스로마이신 사용 중 심바스타틴 일시 중단'),
    ('DIGO001', 'AMIO001', 'major', '디곡신과 아미오다론 병용 시 디곡신 혈중 농도 상승', '디곡신 용량 50% 감량')
ON CONFLICT (drug_a_code, drug_b_code) DO NOTHING;
