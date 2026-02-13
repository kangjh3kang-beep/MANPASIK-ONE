-- =============================================================================
-- 25-admin-settings-ext.sql
-- 관리자 설정 관리 확장: 메타데이터, 다국어, LLM 세션, 변경 대기열
-- =============================================================================

-- ENUM: 설정 카테고리
DO $$ BEGIN
    CREATE TYPE config_category AS ENUM (
        'general','payment','auth','storage','messaging','database','ai','notification','security','integration'
    );
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- ENUM: 설정 값 유형
DO $$ BEGIN
    CREATE TYPE config_value_type AS ENUM (
        'string','number','boolean','secret','url','email','json','select','multiline'
    );
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- ENUM: 설정 보안 등급
DO $$ BEGIN
    CREATE TYPE config_security_level AS ENUM (
        'public','internal','confidential','secret'
    );
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- 설정 메타데이터 (설정 항목의 스키마 정의)
CREATE TABLE IF NOT EXISTS config_metadata (
    config_key       VARCHAR(200) PRIMARY KEY REFERENCES system_configs(key) ON DELETE CASCADE,
    category         config_category NOT NULL DEFAULT 'general',
    value_type       config_value_type NOT NULL DEFAULT 'string',
    security_level   config_security_level NOT NULL DEFAULT 'public',
    is_required      BOOLEAN DEFAULT false,
    default_value    TEXT,
    allowed_values   TEXT[],
    validation_regex TEXT,
    validation_min   NUMERIC,
    validation_max   NUMERIC,
    depends_on       VARCHAR(200),
    depends_value    TEXT,
    env_var_name     TEXT,
    service_name     TEXT,
    restart_required BOOLEAN DEFAULT false,
    display_order    INTEGER DEFAULT 0,
    is_active        BOOLEAN DEFAULT true,
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    updated_at       TIMESTAMPTZ DEFAULT NOW()
);

-- 설정 다국어 설명
CREATE TABLE IF NOT EXISTS config_translations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    config_key      VARCHAR(200) NOT NULL REFERENCES system_configs(key) ON DELETE CASCADE,
    language_code   VARCHAR(5) NOT NULL,
    display_name    VARCHAR(200) NOT NULL,
    description     TEXT NOT NULL,
    placeholder     TEXT,
    help_text       TEXT,
    validation_message TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(config_key, language_code)
);
CREATE INDEX IF NOT EXISTS idx_config_translations_key ON config_translations(config_key);
CREATE INDEX IF NOT EXISTS idx_config_translations_lang ON config_translations(language_code);

-- LLM 설정 어시스턴트 대화 세션
CREATE TABLE IF NOT EXISTS llm_config_sessions (
    session_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    admin_id         VARCHAR(36) NOT NULL,
    language_code    VARCHAR(5) NOT NULL DEFAULT 'ko',
    status           VARCHAR(20) NOT NULL DEFAULT 'active',
    context_category VARCHAR(50),
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    updated_at       TIMESTAMPTZ DEFAULT NOW()
);

-- LLM 대화 메시지
CREATE TABLE IF NOT EXISTS llm_config_messages (
    message_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id       UUID NOT NULL REFERENCES llm_config_sessions(session_id) ON DELETE CASCADE,
    role             VARCHAR(20) NOT NULL,
    content          TEXT NOT NULL,
    suggested_configs JSONB,
    applied          BOOLEAN DEFAULT false,
    applied_at       TIMESTAMPTZ,
    created_at       TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_llm_messages_session ON llm_config_messages(session_id, created_at);

-- 설정 변경 대기열 (LLM 제안 → 관리자 승인 → 적용)
CREATE TABLE IF NOT EXISTS config_change_queue (
    change_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id   UUID REFERENCES llm_config_sessions(session_id),
    config_key   VARCHAR(200) NOT NULL,
    old_value    TEXT,
    new_value    TEXT NOT NULL,
    reason       TEXT,
    suggested_by VARCHAR(20) NOT NULL DEFAULT 'admin',
    status       VARCHAR(20) NOT NULL DEFAULT 'pending',
    approved_by  VARCHAR(36),
    approved_at  TIMESTAMPTZ,
    applied_at   TIMESTAMPTZ,
    created_at   TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_change_queue_status ON config_change_queue(status);
CREATE INDEX IF NOT EXISTS idx_change_queue_key ON config_change_queue(config_key);

-- =============================================================================
-- 시드 데이터: system_configs 기본 설정
-- =============================================================================
INSERT INTO system_configs (key, value, description) VALUES
    ('toss.secret_key',          '', 'Toss Payments 시크릿 키'),
    ('toss.api_url',             'https://api.tosspayments.com', 'Toss API URL'),
    ('toss.sandbox_mode',        'true', 'Toss 샌드박스 모드'),
    ('jwt.secret',               '', 'JWT 서명 시크릿'),
    ('jwt.access_ttl_minutes',   '15', 'JWT Access Token TTL (분)'),
    ('jwt.refresh_ttl_days',     '7', 'JWT Refresh Token TTL (일)'),
    ('keycloak.url',             'http://keycloak:9090', 'Keycloak URL'),
    ('keycloak.realm',           'manpasik', 'Keycloak Realm'),
    ('keycloak.client_id',       'manpasik-api', 'Keycloak Client ID'),
    ('keycloak.client_secret',   '', 'Keycloak Client Secret'),
    ('s3.endpoint',              'minio:9000', 'S3 엔드포인트'),
    ('s3.access_key',            '', 'S3 Access Key'),
    ('s3.secret_key',            '', 'S3 Secret Key'),
    ('s3.bucket',                'manpasik', 'S3 버킷 이름'),
    ('kafka.brokers',            'redpanda:19092', 'Kafka 브로커 주소'),
    ('fcm.server_key',           '', 'FCM 서버 키'),
    ('fcm.project_id',           '', 'Firebase 프로젝트 ID'),
    ('security.cors_origins',    '*', '허용 CORS 오리진'),
    ('security.rate_limit_rpm',  '60', 'API 레이트리밋 (요청/분)'),
    ('ai.default_model',         'biomarker_classifier', '기본 AI 모델'),
    ('ai.confidence_threshold',  '0.85', 'AI 분석 신뢰도 임계값'),
    ('llm.provider',             'openai', 'LLM 제공자'),
    ('llm.api_key',              '', 'LLM API 키'),
    ('llm.model',                'gpt-4o', 'LLM 모델 이름'),
    ('llm.max_tokens',           '2048', 'LLM 최대 토큰 수'),
    ('llm.temperature',          '0.3', 'LLM 온도'),
    ('session_timeout_minutes',  '30', '세션 타임아웃 (분)'),
    ('max_file_upload_mb',       '50', '최대 파일 업로드 크기 (MB)')
ON CONFLICT (key) DO NOTHING;

-- =============================================================================
-- 시드 데이터: config_metadata
-- =============================================================================
INSERT INTO config_metadata (config_key, category, value_type, security_level, is_required, default_value, env_var_name, service_name, restart_required, display_order) VALUES
    ('toss.secret_key',         'payment',      'secret',  'secret',       true,  '', 'TOSS_SECRET_KEY',        'payment-service',       false, 1),
    ('toss.api_url',            'payment',      'url',     'internal',     false, 'https://api.tosspayments.com', 'TOSS_API_URL',      'payment-service',       false, 2),
    ('toss.sandbox_mode',       'payment',      'boolean', 'internal',     false, 'true', NULL,                  'payment-service',       false, 3),
    ('jwt.secret',              'auth',         'secret',  'secret',       true,  '', 'JWT_SECRET',              'auth-service',          true,  1),
    ('jwt.access_ttl_minutes',  'auth',         'number',  'internal',     false, '15', 'JWT_ACCESS_TTL_MINUTES','auth-service',          false, 2),
    ('jwt.refresh_ttl_days',    'auth',         'number',  'internal',     false, '7',  'JWT_REFRESH_TTL_DAYS',  'auth-service',          false, 3),
    ('keycloak.url',            'auth',         'url',     'internal',     false, 'http://keycloak:9090', 'KEYCLOAK_URL',          'auth-service',          true,  4),
    ('keycloak.realm',          'auth',         'string',  'internal',     false, 'manpasik', 'KEYCLOAK_REALM',  'auth-service',          true,  5),
    ('keycloak.client_id',      'auth',         'string',  'internal',     false, 'manpasik-api', 'KEYCLOAK_CLIENT_ID','auth-service',    true,  6),
    ('keycloak.client_secret',  'auth',         'secret',  'secret',       true,  '', 'KEYCLOAK_CLIENT_SECRET', 'auth-service',          true,  7),
    ('s3.endpoint',             'storage',      'url',     'internal',     false, 'minio:9000', 'S3_ENDPOINT',   'gateway',               false, 1),
    ('s3.access_key',           'storage',      'secret',  'secret',       true,  '', 'S3_ACCESS_KEY',           'gateway',               false, 2),
    ('s3.secret_key',           'storage',      'secret',  'secret',       true,  '', 'S3_SECRET_KEY',           'gateway',               false, 3),
    ('s3.bucket',               'storage',      'string',  'internal',     false, 'manpasik', 'S3_BUCKET',       'gateway',               false, 4),
    ('kafka.brokers',           'messaging',    'string',  'internal',     false, 'redpanda:19092', 'KAFKA_BROKERS','*',                   true,  1),
    ('fcm.server_key',          'notification', 'secret',  'secret',       false, '', 'FCM_SERVER_KEY',          'notification-service',  false, 1),
    ('fcm.project_id',          'notification', 'string',  'internal',     false, '', 'FCM_PROJECT_ID',          'notification-service',  false, 2),
    ('security.cors_origins',   'security',     'string',  'internal',     false, '*', NULL,                     'gateway',               false, 1),
    ('security.rate_limit_rpm', 'security',     'number',  'public',       false, '60', NULL,                    'gateway',               false, 2),
    ('ai.default_model',        'ai',           'select',  'public',       false, 'biomarker_classifier', NULL,  'ai-inference-service',  false, 1),
    ('ai.confidence_threshold', 'ai',           'number',  'public',       false, '0.85', NULL,                  'ai-inference-service',  false, 2),
    ('llm.provider',            'ai',           'select',  'internal',     false, 'openai', NULL,                'ai-inference-service',  false, 10),
    ('llm.api_key',             'ai',           'secret',  'secret',       false, '', NULL,                      'ai-inference-service',  false, 11),
    ('llm.model',               'ai',           'string',  'internal',     false, 'gpt-4o', NULL,               'ai-inference-service',  false, 12),
    ('llm.max_tokens',          'ai',           'number',  'internal',     false, '2048', NULL,                  'ai-inference-service',  false, 13),
    ('llm.temperature',         'ai',           'number',  'internal',     false, '0.3', NULL,                   'ai-inference-service',  false, 14),
    ('maintenance_mode',        'general',      'boolean', 'public',       false, 'false', NULL,                 '*',                     false, 1),
    ('max_devices_per_user',    'general',      'number',  'public',       false, '5', NULL,                     '*',                     false, 2),
    ('default_language',        'general',      'select',  'public',       false, 'ko', NULL,                    '*',                     false, 3),
    ('session_timeout_minutes', 'general',      'number',  'public',       false, '30', NULL,                    '*',                     false, 4),
    ('max_file_upload_mb',      'general',      'number',  'public',       false, '50', NULL,                    '*',                     false, 5)
ON CONFLICT (config_key) DO NOTHING;

-- allowed_values 설정
UPDATE config_metadata SET allowed_values = ARRAY['biomarker_classifier','anomaly_detector','trend_predictor','health_scorer','food_calorie_estimator'] WHERE config_key = 'ai.default_model';
UPDATE config_metadata SET allowed_values = ARRAY['openai','anthropic','local'] WHERE config_key = 'llm.provider';
UPDATE config_metadata SET allowed_values = ARRAY['ko','en','ja','zh','fr','hi'] WHERE config_key = 'default_language';
UPDATE config_metadata SET validation_min = 0.0, validation_max = 1.0 WHERE config_key = 'ai.confidence_threshold';
UPDATE config_metadata SET validation_min = 1, validation_max = 120 WHERE config_key = 'jwt.access_ttl_minutes';
UPDATE config_metadata SET validation_min = 1, validation_max = 90 WHERE config_key = 'jwt.refresh_ttl_days';
UPDATE config_metadata SET validation_min = 1, validation_max = 10000 WHERE config_key = 'security.rate_limit_rpm';
UPDATE config_metadata SET validation_min = 5, validation_max = 120 WHERE config_key = 'session_timeout_minutes';
UPDATE config_metadata SET validation_min = 1, validation_max = 500 WHERE config_key = 'max_file_upload_mb';

-- =============================================================================
-- 시드 데이터: config_translations (ko/en/ja 주요 항목)
-- =============================================================================
INSERT INTO config_translations (config_key, language_code, display_name, description, placeholder, help_text) VALUES
    ('toss.secret_key', 'ko', 'Toss 시크릿 키', 'Toss Payments 결제 승인/취소 API 호출 시 인증에 사용됩니다.', 'test_sk_... 또는 live_sk_...', '## 발급 방법\n1. [Toss 개발자센터](https://developers.tosspayments.com/) 로그인\n2. **내 개발정보** → **API 키** 선택\n3. **시크릿 키** 복사\n\n> 테스트 키(`test_sk_`)와 라이브 키(`live_sk_`)를 구분하세요.'),
    ('toss.secret_key', 'en', 'Toss Secret Key', 'Used for authentication when calling Toss Payments confirmation/cancellation APIs.', 'test_sk_... or live_sk_...', '## How to obtain\n1. Log in to [Toss Developer Center](https://developers.tosspayments.com/)\n2. Go to **My Dev Info** → **API Keys**\n3. Copy the **Secret Key**'),
    ('toss.secret_key', 'ja', 'Toss シークレットキー', 'Toss Payments の決済承認・取消API呼び出し時の認証に使用されます。', 'test_sk_... または live_sk_...', '## 取得方法\n1. [Toss開発者センター](https://developers.tosspayments.com/)にログイン\n2. **開発情報** → **APIキー** 選択\n3. **シークレットキー**をコピー'),
    ('jwt.secret', 'ko', 'JWT 시크릿', 'JWT 토큰 서명에 사용되는 비밀 키입니다.', '32자 이상의 랜덤 문자열', '강력한 랜덤 문자열을 사용하세요. `openssl rand -hex 32`로 생성 가능합니다.'),
    ('jwt.secret', 'en', 'JWT Secret', 'Secret key used for signing JWT tokens.', 'Random string (32+ chars)', 'Use a strong random string. Generate with `openssl rand -hex 32`.'),
    ('maintenance_mode', 'ko', '유지보수 모드', '활성화하면 일반 사용자의 접근이 차단됩니다.', '', '> **주의**: 활성화 시 모든 일반 사용자가 서비스를 이용할 수 없습니다.'),
    ('maintenance_mode', 'en', 'Maintenance Mode', 'When enabled, regular users cannot access the service.', '', '> **Warning**: All regular users will be blocked when enabled.'),
    ('ai.default_model', 'ko', '기본 AI 모델', '측정 분석에 사용할 기본 AI 모델입니다.', '', '사용 가능한 모델: biomarker_classifier, anomaly_detector, trend_predictor, health_scorer, food_calorie_estimator'),
    ('ai.default_model', 'en', 'Default AI Model', 'Default AI model for measurement analysis.', '', 'Available models: biomarker_classifier, anomaly_detector, trend_predictor, health_scorer, food_calorie_estimator'),
    ('fcm.server_key', 'ko', 'FCM 서버 키', 'Firebase Cloud Messaging 푸시 알림 전송에 사용됩니다.', '', '## 발급 방법\n1. [Firebase Console](https://console.firebase.google.com/) 접속\n2. 프로젝트 설정 → Cloud Messaging 탭\n3. **서버 키** 복사'),
    ('fcm.server_key', 'en', 'FCM Server Key', 'Used for sending push notifications via Firebase Cloud Messaging.', '', '## How to obtain\n1. Go to [Firebase Console](https://console.firebase.google.com/)\n2. Project Settings → Cloud Messaging\n3. Copy **Server Key**')
ON CONFLICT (config_key, language_code) DO NOTHING;
