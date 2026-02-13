-- translation-service 데이터베이스 스키마
-- Phase 3C: 실시간 번역

-- 번역 이력 테이블
CREATE TABLE IF NOT EXISTS translation_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID,
    source_text TEXT NOT NULL,
    translated_text TEXT NOT NULL,
    source_language VARCHAR(10) NOT NULL,
    target_language VARCHAR(10) NOT NULL,
    confidence REAL DEFAULT 0.0,
    context_hint VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_translations_user ON translation_records(user_id);
CREATE INDEX IF NOT EXISTS idx_translations_created ON translation_records(created_at);

-- 번역 사용량 테이블
CREATE TABLE IF NOT EXISTS translation_usage (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE,
    total_characters BIGINT DEFAULT 0,
    monthly_characters BIGINT DEFAULT 0,
    monthly_limit BIGINT DEFAULT 100000,
    total_requests INT DEFAULT 0,
    monthly_requests INT DEFAULT 0,
    month_start DATE NOT NULL DEFAULT DATE_TRUNC('month', CURRENT_DATE),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_usage_user ON translation_usage(user_id);

-- 의료 용어 사전 테이블
CREATE TABLE IF NOT EXISTS medical_terms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_language VARCHAR(10) NOT NULL,
    target_language VARCHAR(10) NOT NULL,
    source_term VARCHAR(200) NOT NULL,
    target_term VARCHAR(200) NOT NULL,
    category VARCHAR(50) DEFAULT 'general',
    verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(source_language, target_language, source_term)
);

CREATE INDEX IF NOT EXISTS idx_medical_terms_lang ON medical_terms(source_language, target_language);

-- 지원 언어 테이블
CREATE TABLE IF NOT EXISTS supported_languages (
    language_code VARCHAR(10) PRIMARY KEY,
    language_name VARCHAR(100) NOT NULL,
    native_name VARCHAR(100) NOT NULL,
    supports_medical BOOLEAN DEFAULT FALSE,
    supports_realtime BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 시드 데이터: 지원 언어
INSERT INTO supported_languages (language_code, language_name, native_name, supports_medical, supports_realtime) VALUES
    ('ko', 'Korean', '한국어', TRUE, TRUE),
    ('en', 'English', 'English', TRUE, TRUE),
    ('ja', 'Japanese', '日本語', TRUE, TRUE),
    ('zh', 'Chinese', '中文', TRUE, TRUE),
    ('fr', 'French', 'Français', TRUE, TRUE),
    ('de', 'German', 'Deutsch', TRUE, FALSE),
    ('es', 'Spanish', 'Español', TRUE, TRUE),
    ('hi', 'Hindi', 'हिन्दी', FALSE, TRUE),
    ('pt', 'Portuguese', 'Português', FALSE, FALSE),
    ('ar', 'Arabic', 'العربية', FALSE, FALSE)
ON CONFLICT (language_code) DO NOTHING;

-- 시드 데이터: 의료 용어 (한→영)
INSERT INTO medical_terms (source_language, target_language, source_term, target_term, category, verified) VALUES
    ('ko', 'en', '혈압', 'blood pressure', 'vital_signs', TRUE),
    ('ko', 'en', '혈당', 'blood glucose', 'vital_signs', TRUE),
    ('ko', 'en', '체온', 'body temperature', 'vital_signs', TRUE),
    ('ko', 'en', '처방전', 'prescription', 'pharmacy', TRUE),
    ('ko', 'en', '복용량', 'dosage', 'pharmacy', TRUE),
    ('ko', 'en', '부작용', 'side effect', 'pharmacy', TRUE),
    ('ko', 'en', '진단', 'diagnosis', 'clinical', TRUE),
    ('ko', 'en', '증상', 'symptom', 'clinical', TRUE),
    ('ko', 'en', '항생제', 'antibiotic', 'pharmacy', TRUE),
    ('ko', 'en', '진통제', 'analgesic', 'pharmacy', TRUE)
ON CONFLICT (source_language, target_language, source_term) DO NOTHING;
