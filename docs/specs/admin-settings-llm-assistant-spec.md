# 관리자 설정 관리 시스템 + LLM 대리 설정 기능 세부 기획서

> **문서 ID**: SPEC-ADMIN-SETTINGS-LLM-001  
> **버전**: v1.0  
> **작성일**: 2026-02-12  
> **상태**: 기획 확정 (구현 대기)  
> **대상**: admin-service (백엔드), Flutter admin UI (프론트엔드), ai-inference-service 확장

---

## 1. 개요

### 1.1 배경

현재 ManPaSik 시스템은 Toss PG 시크릿 키, Keycloak 비밀, JWT 시크릿, S3 키 등 **민감 설정값을 환경변수(`.env`)로만 관리**한다. 이 방식은:
- 서버 재시작 없이 설정 변경 불가
- 개발자/운영자가 SSH로 직접 서버에 접근해야 함
- 설정값의 의미나 올바른 입력 방법을 설명할 수 없음
- 다국어 환경에서 비한국어 관리자가 설정을 이해하기 어려움

### 1.2 목표

1. **관리자 대시보드에서 모든 시스템 설정을 UI로 관리** — 키·값·설명·유효성·보안 등급
2. **관리자 지정 언어로 설정 설명을 자동 번역** — 6개 언어(ko/en/ja/zh/fr/hi) + 10개(translation-service)
3. **LLM 기반 설정 어시스턴트** — 설정의 의미 설명, 올바른 값 추천, 유효성 검증, 관리자 승인 후 대리 설정 수행

### 1.3 핵심 원칙

- **보안 우선**: 민감 키는 AES-256 암호화 저장, UI에서 마스킹, 감사 로그 기록
- **무중단 반영**: 설정 변경 시 서비스 재시작 없이 즉시 또는 점진적 적용
- **역할 기반 접근**: Super Admin만 민감 설정 변경 가능, Moderator는 읽기만
- **LLM 안전성**: LLM은 설정을 "제안"만 하고, 관리자 "확인" 후 적용. 민감 키값은 LLM에 전달하지 않음

---

## 2. 시스템 아키텍처

```
┌─────────────────────────────────────────────────────────────┐
│  Flutter Admin Dashboard                                     │
│  ┌──────────┐ ┌──────────────┐ ┌──────────────────────────┐ │
│  │설정 목록  │ │설정 편집기    │ │LLM 어시스턴트 채팅 패널  │ │
│  │(카테고리별)│ │(유형별 입력)  │ │(설정 질문/제안/실행)     │ │
│  └──────────┘ └──────────────┘ └──────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│  API Gateway (REST→gRPC)                                     │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐ ┌──────────────┐ ┌──────────────────────┐ │
│  │admin-service │ │translation-  │ │ai-inference-service  │ │
│  │              │ │service       │ │(LLM 어시스턴트 확장) │ │
│  │SetConfig     │ │TranslateText │ │ConfigAssistantChat   │ │
│  │GetConfig     │ │              │ │SuggestConfigValue    │ │
│  │ListConfigs   │ │              │ │ValidateConfig        │ │
│  │GetConfigMeta │ │              │ │ApplyConfigSuggestion │ │
│  └─────┬───────┘ └──────────────┘ └──────────┬───────────┘ │
│        │                                       │             │
│  ┌─────▼───────────────────────────────────────▼───────────┐ │
│  │  PostgreSQL                                              │ │
│  │  system_configs + config_metadata + config_translations  │ │
│  │  + llm_config_sessions + config_change_queue             │ │
│  └──────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

---

## 3. 데이터 모델

### 3.1 DB 스키마 확장 (`25-admin-settings-ext.sql`)

```sql
-- =============================================================================
-- Admin Settings Extension: 설정 메타데이터, 다국어, LLM 세션
-- =============================================================================

-- 설정 카테고리
CREATE TYPE config_category AS ENUM (
    'general',         -- 일반 (유지보수 모드, 기본 언어 등)
    'payment',         -- 결제 (Toss PG, 환불 정책 등)
    'auth',            -- 인증 (Keycloak, JWT 등)
    'storage',         -- 스토리지 (S3/MinIO 등)
    'messaging',       -- 메시징 (Kafka, MQTT 등)
    'database',        -- 데이터베이스 (PG, Redis, Milvus, ES 등)
    'ai',              -- AI/ML (모델 설정, 추론 설정 등)
    'notification',    -- 알림 (FCM, 이메일 등)
    'security',        -- 보안 (암호화, CORS, 레이트리밋 등)
    'integration'      -- 외부 연동 (API 키, 웹훅 등)
);

-- 설정 값 유형
CREATE TYPE config_value_type AS ENUM (
    'string',          -- 일반 문자열
    'number',          -- 숫자 (정수·실수)
    'boolean',         -- true/false
    'secret',          -- 암호화 저장, UI 마스킹
    'url',             -- URL 형식 검증
    'email',           -- 이메일 형식 검증
    'json',            -- JSON 구조
    'select',          -- 선택지 (allowed_values)
    'multiline'        -- 여러 줄 텍스트
);

-- 설정 보안 등급
CREATE TYPE config_security_level AS ENUM (
    'public',          -- 모든 관리자 조회 가능
    'internal',        -- Admin 이상만 조회
    'confidential',    -- Super Admin만 조회/변경
    'secret'           -- Super Admin만 변경, 값은 암호화 저장, LLM에 전달 금지
);

-- 설정 메타데이터 (설정 항목의 스키마 정의)
CREATE TABLE IF NOT EXISTS config_metadata (
    config_key          VARCHAR(200) PRIMARY KEY REFERENCES system_configs(key) ON DELETE CASCADE,
    category            config_category NOT NULL DEFAULT 'general',
    value_type          config_value_type NOT NULL DEFAULT 'string',
    security_level      config_security_level NOT NULL DEFAULT 'public',
    is_required         BOOLEAN DEFAULT false,
    default_value       TEXT,
    allowed_values      TEXT[],                  -- select 유형일 때 허용 값 목록
    validation_regex    TEXT,                    -- 정규식 유효성 검증
    validation_min      NUMERIC,                -- 숫자 최소값
    validation_max      NUMERIC,                -- 숫자 최대값
    depends_on          VARCHAR(200),            -- 의존 설정 키 (다른 설정이 특정 값일 때만 활성)
    depends_value       TEXT,                    -- 의존 설정의 활성 조건 값
    env_var_name        TEXT,                    -- 대응하는 환경변수 이름 (참조용)
    service_name        TEXT,                    -- 관련 서비스 이름
    restart_required    BOOLEAN DEFAULT false,   -- 변경 시 서비스 재시작 필요 여부
    display_order       INTEGER DEFAULT 0,       -- 카테고리 내 표시 순서
    is_active           BOOLEAN DEFAULT true,    -- 비활성화 시 UI에서 숨김
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);

-- 설정 다국어 설명
CREATE TABLE IF NOT EXISTS config_translations (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    config_key          VARCHAR(200) NOT NULL REFERENCES system_configs(key) ON DELETE CASCADE,
    language_code       VARCHAR(5) NOT NULL,     -- 'ko', 'en', 'ja', 'zh', 'fr', 'hi' 등
    display_name        VARCHAR(200) NOT NULL,   -- 설정 표시 이름
    description         TEXT NOT NULL,            -- 상세 설명
    placeholder         TEXT,                    -- 입력 필드 플레이스홀더
    help_text           TEXT,                    -- 추가 도움말 (마크다운 지원)
    validation_message  TEXT,                    -- 유효성 검증 실패 시 메시지
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(config_key, language_code)
);

CREATE INDEX IF NOT EXISTS idx_config_translations_key ON config_translations(config_key);
CREATE INDEX IF NOT EXISTS idx_config_translations_lang ON config_translations(language_code);

-- LLM 설정 어시스턴트 대화 세션
CREATE TABLE IF NOT EXISTS llm_config_sessions (
    session_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    admin_id            VARCHAR(36) NOT NULL,
    language_code       VARCHAR(5) NOT NULL DEFAULT 'ko',
    status              VARCHAR(20) NOT NULL DEFAULT 'active', -- 'active', 'completed', 'cancelled'
    context_category    config_category,         -- 특정 카테고리 컨텍스트 (null=전체)
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);

-- LLM 대화 메시지
CREATE TABLE IF NOT EXISTS llm_config_messages (
    message_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id          UUID NOT NULL REFERENCES llm_config_sessions(session_id) ON DELETE CASCADE,
    role                VARCHAR(20) NOT NULL,    -- 'user', 'assistant', 'system'
    content             TEXT NOT NULL,
    suggested_configs   JSONB,                   -- 어시스턴트가 제안한 설정 변경 [{key, value, reason}]
    applied             BOOLEAN DEFAULT false,   -- 제안이 실제 적용되었는지
    applied_at          TIMESTAMPTZ,
    created_at          TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_llm_messages_session ON llm_config_messages(session_id, created_at);

-- 설정 변경 대기열 (LLM 제안 → 관리자 승인 → 적용)
CREATE TABLE IF NOT EXISTS config_change_queue (
    change_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id          UUID REFERENCES llm_config_sessions(session_id),
    config_key          VARCHAR(200) NOT NULL,
    old_value           TEXT,
    new_value           TEXT NOT NULL,
    reason              TEXT,                    -- 변경 사유
    suggested_by        VARCHAR(20) NOT NULL DEFAULT 'admin', -- 'admin', 'llm'
    status              VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending', 'approved', 'rejected', 'applied'
    approved_by         VARCHAR(36),
    approved_at         TIMESTAMPTZ,
    applied_at          TIMESTAMPTZ,
    created_at          TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_change_queue_status ON config_change_queue(status);
CREATE INDEX IF NOT EXISTS idx_change_queue_key ON config_change_queue(config_key);
```

### 3.2 초기 설정 메타데이터 시드 (주요 항목)

```sql
-- 기존 system_configs에 키가 없으면 추가
INSERT INTO system_configs (key, value, description) VALUES
    -- 결제
    ('toss.secret_key', '', 'Toss Payments 시크릿 키'),
    ('toss.api_url', 'https://api.tosspayments.com', 'Toss API URL'),
    ('toss.sandbox_mode', 'true', 'Toss 샌드박스 모드'),
    -- 인증
    ('jwt.secret', '', 'JWT 서명 시크릿'),
    ('jwt.access_ttl_minutes', '15', 'JWT Access Token TTL (분)'),
    ('jwt.refresh_ttl_days', '7', 'JWT Refresh Token TTL (일)'),
    ('keycloak.url', 'http://keycloak:9090', 'Keycloak URL'),
    ('keycloak.realm', 'manpasik', 'Keycloak Realm'),
    ('keycloak.client_id', 'manpasik-api', 'Keycloak Client ID'),
    ('keycloak.client_secret', '', 'Keycloak Client Secret'),
    -- 스토리지
    ('s3.endpoint', 'minio:9000', 'S3 엔드포인트'),
    ('s3.access_key', '', 'S3 Access Key'),
    ('s3.secret_key', '', 'S3 Secret Key'),
    ('s3.bucket', 'manpasik', 'S3 버킷 이름'),
    -- 메시징
    ('kafka.brokers', 'redpanda:19092', 'Kafka 브로커 주소'),
    -- 알림
    ('fcm.server_key', '', 'FCM 서버 키'),
    ('fcm.project_id', '', 'Firebase 프로젝트 ID'),
    -- 보안
    ('security.cors_origins', '*', '허용 CORS 오리진'),
    ('security.rate_limit_rpm', '60', 'API 레이트리밋 (요청/분)'),
    -- AI
    ('ai.default_model', 'biomarker_classifier', '기본 AI 모델'),
    ('ai.confidence_threshold', '0.85', 'AI 분석 신뢰도 임계값'),
    -- LLM
    ('llm.provider', 'openai', 'LLM 제공자 (openai, anthropic, local)'),
    ('llm.api_key', '', 'LLM API 키'),
    ('llm.model', 'gpt-4o', 'LLM 모델 이름'),
    ('llm.max_tokens', '2048', 'LLM 최대 토큰 수'),
    ('llm.temperature', '0.3', 'LLM 온도 (낮을수록 결정적)')
ON CONFLICT (key) DO NOTHING;

-- 메타데이터 시드
INSERT INTO config_metadata (
    config_key, category, value_type, security_level,
    is_required, default_value, env_var_name, service_name,
    restart_required, display_order
) VALUES
    -- 결제
    ('toss.secret_key', 'payment', 'secret', 'secret', true, '', 'TOSS_SECRET_KEY', 'payment-service', false, 1),
    ('toss.api_url', 'payment', 'url', 'internal', false, 'https://api.tosspayments.com', 'TOSS_API_URL', 'payment-service', false, 2),
    ('toss.sandbox_mode', 'payment', 'boolean', 'internal', false, 'true', NULL, 'payment-service', false, 3),
    -- 인증
    ('jwt.secret', 'auth', 'secret', 'secret', true, '', 'JWT_SECRET', 'auth-service', true, 1),
    ('jwt.access_ttl_minutes', 'auth', 'number', 'internal', false, '15', 'JWT_ACCESS_TTL_MINUTES', 'auth-service', false, 2),
    ('jwt.refresh_ttl_days', 'auth', 'number', 'internal', false, '7', 'JWT_REFRESH_TTL_DAYS', 'auth-service', false, 3),
    ('keycloak.url', 'auth', 'url', 'internal', false, 'http://keycloak:9090', 'KEYCLOAK_URL', 'auth-service', true, 4),
    ('keycloak.realm', 'auth', 'string', 'internal', false, 'manpasik', 'KEYCLOAK_REALM', 'auth-service', true, 5),
    ('keycloak.client_id', 'auth', 'string', 'internal', false, 'manpasik-api', 'KEYCLOAK_CLIENT_ID', 'auth-service', true, 6),
    ('keycloak.client_secret', 'auth', 'secret', 'secret', true, '', 'KEYCLOAK_CLIENT_SECRET', 'auth-service', true, 7),
    -- 스토리지
    ('s3.endpoint', 'storage', 'url', 'internal', false, 'minio:9000', 'S3_ENDPOINT', 'gateway', false, 1),
    ('s3.access_key', 'storage', 'secret', 'secret', true, '', 'S3_ACCESS_KEY', 'gateway', false, 2),
    ('s3.secret_key', 'storage', 'secret', 'secret', true, '', 'S3_SECRET_KEY', 'gateway', false, 3),
    ('s3.bucket', 'storage', 'string', 'internal', false, 'manpasik', 'S3_BUCKET', 'gateway', false, 4),
    -- 메시징
    ('kafka.brokers', 'messaging', 'string', 'internal', false, 'redpanda:19092', 'KAFKA_BROKERS', '*', true, 1),
    -- 알림
    ('fcm.server_key', 'notification', 'secret', 'secret', false, '', 'FCM_SERVER_KEY', 'notification-service', false, 1),
    ('fcm.project_id', 'notification', 'string', 'internal', false, '', 'FCM_PROJECT_ID', 'notification-service', false, 2),
    -- 보안
    ('security.cors_origins', 'security', 'string', 'internal', false, '*', NULL, 'gateway', false, 1),
    ('security.rate_limit_rpm', 'security', 'number', 'public', false, '60', NULL, 'gateway', false, 2),
    -- AI
    ('ai.default_model', 'ai', 'select', 'public', false, 'biomarker_classifier', NULL, 'ai-inference-service', false, 1),
    ('ai.confidence_threshold', 'ai', 'number', 'public', false, '0.85', NULL, 'ai-inference-service', false, 2),
    -- LLM
    ('llm.provider', 'ai', 'select', 'internal', false, 'openai', NULL, 'ai-inference-service', false, 10),
    ('llm.api_key', 'ai', 'secret', 'secret', false, '', NULL, 'ai-inference-service', false, 11),
    ('llm.model', 'ai', 'string', 'internal', false, 'gpt-4o', NULL, 'ai-inference-service', false, 12),
    ('llm.max_tokens', 'ai', 'number', 'internal', false, '2048', NULL, 'ai-inference-service', false, 13),
    ('llm.temperature', 'ai', 'number', 'internal', false, '0.3', NULL, 'ai-inference-service', false, 14),
    -- 일반
    ('maintenance_mode', 'general', 'boolean', 'public', false, 'false', NULL, '*', false, 1),
    ('max_devices_per_user', 'general', 'number', 'public', false, '5', NULL, '*', false, 2),
    ('default_language', 'general', 'select', 'public', false, 'ko', NULL, '*', false, 3),
    ('session_timeout_minutes', 'general', 'number', 'public', false, '30', NULL, '*', false, 4),
    ('max_file_upload_mb', 'general', 'number', 'public', false, '50', NULL, '*', false, 5)
ON CONFLICT (config_key) DO NOTHING;
```

### 3.3 다국어 설명 시드 예시 (ko/en)

```sql
-- Toss 시크릿 키 (ko)
INSERT INTO config_translations (config_key, language_code, display_name, description, placeholder, help_text) VALUES
('toss.secret_key', 'ko', 'Toss 시크릿 키',
 'Toss Payments에서 발급받은 시크릿 키입니다. 결제 승인 및 취소 API 호출 시 인증에 사용됩니다.',
 'test_sk_... 또는 live_sk_...',
 '## 발급 방법\n1. [Toss 개발자센터](https://developers.tosspayments.com/) 로그인\n2. **내 개발정보** → **API 키** 메뉴 선택\n3. **시크릿 키** 복사\n\n> **주의**: 테스트 키(`test_sk_`)와 라이브 키(`live_sk_`)를 구분하세요.\n> 샌드박스 모드가 켜져 있으면 테스트 키를, 꺼져 있으면 라이브 키를 사용합니다.',
 '올바른 Toss 시크릿 키 형식이 아닙니다'),
-- Toss 시크릿 키 (en)
('toss.secret_key', 'en', 'Toss Secret Key',
 'The secret key issued by Toss Payments. Used for authentication when calling payment confirmation and cancellation APIs.',
 'test_sk_... or live_sk_...',
 '## How to obtain\n1. Log in to [Toss Developer Center](https://developers.tosspayments.com/)\n2. Go to **My Dev Info** → **API Keys**\n3. Copy the **Secret Key**\n\n> **Note**: Distinguish between test keys (`test_sk_`) and live keys (`live_sk_`).\n> Use test key when sandbox mode is on, live key when off.',
 'Invalid Toss secret key format'),
-- Toss 시크릿 키 (ja)
('toss.secret_key', 'ja', 'Toss シークレットキー',
 'Toss Paymentsが発行したシークレットキーです。決済承認・取消API呼び出し時の認証に使用されます。',
 'test_sk_... または live_sk_...',
 '## 取得方法\n1. [Toss開発者センター](https://developers.tosspayments.com/)にログイン\n2. **開発情報** → **APIキー**メニュー選択\n3. **シークレットキー**をコピー\n\n> **注意**: テストキー(`test_sk_`)とライブキー(`live_sk_`)を区別してください。',
 '正しいTossシークレットキー形式ではありません')
ON CONFLICT (config_key, language_code) DO NOTHING;
```

---

## 4. gRPC API 확장

### 4.1 AdminService 확장 (admin-service)

```protobuf
// 설정 관리 확장 RPC
service AdminService {
  // --- 기존 ---
  rpc SetSystemConfig(SetSystemConfigRequest) returns (SystemConfig);
  rpc GetSystemConfig(GetSystemConfigRequest) returns (SystemConfig);

  // --- 신규: 설정 목록·메타·번역 ---
  rpc ListSystemConfigs(ListSystemConfigsRequest) returns (ListSystemConfigsResponse);
  rpc GetConfigWithMeta(GetConfigWithMetaRequest) returns (ConfigWithMeta);
  rpc ListConfigsByCategory(ListConfigsByCategoryRequest) returns (ListSystemConfigsResponse);
  rpc ValidateConfigValue(ValidateConfigValueRequest) returns (ValidateConfigValueResponse);
  rpc BulkSetConfigs(BulkSetConfigsRequest) returns (BulkSetConfigsResponse);

  // --- 신규: 설정 변경 대기열 ---
  rpc ListPendingChanges(ListPendingChangesRequest) returns (ListPendingChangesResponse);
  rpc ApproveChange(ApproveChangeRequest) returns (ConfigChangeResult);
  rpc RejectChange(RejectChangeRequest) returns (ConfigChangeResult);
}

// --- 메시지 정의 ---

message ListSystemConfigsRequest {
  string language_code = 1;       // 번역 언어 (빈 문자열이면 ko)
  string category = 2;            // 카테고리 필터 (빈 문자열이면 전체)
  bool include_secrets = 3;       // secret 값 포함 여부 (Super Admin만)
}

message ListSystemConfigsResponse {
  repeated ConfigWithMeta configs = 1;
  map<string, int32> category_counts = 2;  // 카테고리별 설정 수
}

message ConfigWithMeta {
  // 값
  string key = 1;
  string value = 2;               // secret이면 마스킹 (****) 반환
  string raw_value = 3;           // include_secrets=true일 때만 실제 값

  // 메타
  string category = 4;
  string value_type = 5;
  string security_level = 6;
  bool is_required = 7;
  string default_value = 8;
  repeated string allowed_values = 9;
  string validation_regex = 10;
  double validation_min = 11;
  double validation_max = 12;
  string depends_on = 13;
  string depends_value = 14;
  string env_var_name = 15;
  string service_name = 16;
  bool restart_required = 17;

  // 다국어 번역
  string display_name = 20;       // 요청 언어의 표시 이름
  string description = 21;        // 요청 언어의 설명
  string placeholder = 22;
  string help_text = 23;          // 마크다운 도움말
  string validation_message = 24;

  // 변경 정보
  string updated_by = 30;
  google.protobuf.Timestamp updated_at = 31;
}

message ValidateConfigValueRequest {
  string key = 1;
  string value = 2;
}

message ValidateConfigValueResponse {
  bool valid = 1;
  string error_message = 2;       // 유효하지 않을 때 에러 메시지 (요청 언어)
  repeated string suggestions = 3; // 올바른 값 제안
}

message BulkSetConfigsRequest {
  repeated SetSystemConfigRequest configs = 1;
  string reason = 2;              // 일괄 변경 사유
}

message BulkSetConfigsResponse {
  repeated ConfigChangeResult results = 1;
  int32 success_count = 2;
  int32 failure_count = 3;
}

message ConfigChangeResult {
  string key = 1;
  bool success = 2;
  string error_message = 3;
}
```

### 4.2 LLM 어시스턴트 RPC (ai-inference-service 확장)

```protobuf
// LLM 설정 어시스턴트 — ai-inference-service에 추가
service AiInferenceService {
  // --- 기존 ---
  rpc AnalyzeMeasurement(...) returns (...);

  // --- 신규: 설정 어시스턴트 ---
  rpc StartConfigSession(StartConfigSessionRequest) returns (ConfigSessionResponse);
  rpc SendConfigMessage(SendConfigMessageRequest) returns (ConfigAssistantResponse);
  rpc ApplyConfigSuggestion(ApplyConfigSuggestionRequest) returns (ApplyConfigSuggestionResponse);
  rpc EndConfigSession(EndConfigSessionRequest) returns (ConfigSessionResponse);
}

message StartConfigSessionRequest {
  string admin_id = 1;
  string language_code = 2;       // 어시스턴트 응답 언어
  string category = 3;            // 특정 카테고리에 대해 대화 (빈 문자열=전체)
}

message ConfigSessionResponse {
  string session_id = 1;
  string status = 2;
  string welcome_message = 3;     // 어시스턴트 환영 메시지 (요청 언어)
}

message SendConfigMessageRequest {
  string session_id = 1;
  string message = 2;             // 관리자 입력 메시지
}

message ConfigAssistantResponse {
  string session_id = 1;
  string message = 2;             // 어시스턴트 응답 (요청 언어)
  repeated ConfigSuggestion suggestions = 3; // 설정 변경 제안 (있을 때만)
  bool requires_confirmation = 4; // 제안 실행에 관리자 확인 필요
}

message ConfigSuggestion {
  string suggestion_id = 1;
  string config_key = 2;
  string current_value = 3;       // 현재 값 (secret이면 마스킹)
  string suggested_value = 4;     // 제안 값 (secret이면 클라이언트가 직접 입력해야 함)
  string reason = 5;              // 제안 이유 (요청 언어)
  bool is_secret = 6;             // true이면 클라이언트가 값을 직접 입력해야 함
}

message ApplyConfigSuggestionRequest {
  string session_id = 1;
  string suggestion_id = 2;
  string admin_id = 3;
  string override_value = 4;      // secret이면 클라이언트가 직접 입력한 값
}

message ApplyConfigSuggestionResponse {
  bool success = 1;
  string message = 2;
  ConfigChangeResult result = 3;
}
```

---

## 5. LLM 어시스턴트 동작 규격

### 5.1 시스템 프롬프트 (서비스 내부)

```
You are the ManPaSik system configuration assistant. Your role:
1. Explain system settings in the admin's chosen language ({language_code}).
2. Help admins understand what each setting does, how to obtain values (e.g., API keys), and what valid values look like.
3. Suggest configuration changes based on admin requests.
4. NEVER output or ask for the actual value of 'secret' level settings. Instead, instruct the admin to enter the value directly in the input field.
5. When suggesting changes, always provide a clear reason.
6. If asked about settings you don't know, say so honestly.

Available settings categories: {categories_json}
Current config state (non-secret only): {current_configs_json}

Respond in {language_code}. Use markdown for formatting.
```

### 5.2 대화 흐름 예시

**시나리오: 관리자가 Toss PG 설정을 요청 (한국어)**

```
[관리자] Toss 결제 연동하려면 어떻게 해야 하나요?

[어시스턴트] ## Toss Payments 연동 설정

Toss 결제를 연동하려면 다음 3개 설정이 필요합니다:

1. **Toss 시크릿 키** (`toss.secret_key`)
   - Toss 개발자센터에서 발급받은 시크릿 키
   - [발급 방법 보기 →]

2. **Toss API URL** (`toss.api_url`)
   - 현재 값: `https://api.tosspayments.com` ✅ (기본값 사용 중)

3. **샌드박스 모드** (`toss.sandbox_mode`)
   - 현재 값: `true`
   - 테스트 환경에서는 `true`, 운영 환경에서는 `false`로 설정하세요

시크릿 키를 설정하시겠습니까? 키 입력 필드를 열어드릴게요.

[관리자] 네, 설정해주세요. 그리고 샌드박스도 꺼주세요.

[어시스턴트] 다음과 같이 변경을 제안합니다:

| 설정 | 현재 | 변경 후 | 사유 |
|------|------|---------|------|
| `toss.secret_key` | (미설정) | 🔒 직접 입력 필요 | 결제 승인에 필수 |
| `toss.sandbox_mode` | true | **false** | 운영 환경 전환 |

> ⚠️ **주의**: 시크릿 키는 보안상 제가 직접 설정할 수 없습니다. 아래 입력 필드에 직접 입력해 주세요.
> 샌드박스 모드를 끄면 실제 결제가 발생합니다. 확인하시겠습니까?

[적용] [취소]
```

**시나리오: 영어 관리자가 AI 설정을 질문**

```
[Admin] What AI models are available and how do I configure them?

[Assistant] ## AI Configuration

ManPaSik currently supports 5 AI models:

| Model | Version | Purpose | Accuracy |
|-------|---------|---------|----------|
| Biomarker Classifier | v1.0.0 | Blood/urine analysis | 94.2% |
| Anomaly Detector | v1.0.0 | Time-series anomaly detection | 91.8% |
| Trend Predictor | v1.0.0 | Health indicator forecasting | 87.6% |
| Health Scorer | v1.0.0 | Comprehensive health scoring | 90.5% |
| Food Calorie Estimator | v0.9.0-beta | Food calorie estimation | 82.3% |

**Current settings:**
- Default model: `biomarker_classifier` ✅
- Confidence threshold: `0.85` (85%)

Would you like to change the default model or adjust the confidence threshold?
```

### 5.3 LLM 안전 규칙

| 규칙 | 설명 |
|------|------|
| **S1: 시크릿 미전달** | security_level="secret"인 설정의 실제 값을 LLM 프롬프트에 포함하지 않음 |
| **S2: 제안만** | LLM은 설정 변경을 "제안"만 하고, 관리자가 명시적으로 "적용" 버튼을 눌러야 반영 |
| **S3: 감사 추적** | LLM을 통한 모든 설정 변경은 audit_logs에 `suggested_by=llm` 기록 |
| **S4: 역할 검증** | LLM이 제안하더라도 해당 관리자의 역할이 설정 보안 등급에 미달하면 거부 |
| **S5: 값 검증** | LLM 제안 값도 동일한 validation_regex/min/max 검증 통과 필요 |
| **S6: 재시작 경고** | restart_required=true인 설정 변경 시 관리자에게 경고 표시 |

---

## 6. Flutter Admin UI 설계

### 6.1 화면 구성

```
📱 Admin Dashboard (관리자 대시보드)
├── /admin/settings                    ← 설정 관리 메인
│   ├── 카테고리 탭 (general, payment, auth, ...)
│   ├── 설정 카드 목록 (카테고리별)
│   │   ├── 카드: display_name + description + 현재 값
│   │   ├── 편집 버튼 → 설정 편집 다이얼로그
│   │   └── 도움말 아이콘 → help_text 팝업 (마크다운 렌더링)
│   └── 검색·필터 (키워드, 서비스별, 변경 필요 항목)
│
├── /admin/settings/assistant          ← LLM 어시스턴트
│   ├── 채팅 인터페이스 (좌: 대화, 우: 제안된 설정 변경 패널)
│   ├── 입력: 텍스트 + 카테고리 선택
│   └── 제안 카드: [적용] [거부] [수정 후 적용] 버튼
│
├── /admin/settings/pending            ← 대기 중인 변경
│   ├── LLM 제안 대기 목록
│   └── [일괄 승인] [일괄 거부]
│
└── /admin/settings/history            ← 변경 이력
    └── audit_logs (resource_type=system_config) 필터된 목록
```

### 6.2 설정 편집 다이얼로그 UX

```
┌──────────────────────────────────────────────┐
│  🔧 Toss 시크릿 키 (Toss Secret Key)         │  ← display_name (관리자 언어)
│                                                │
│  결제 승인 및 취소 API 호출 시 인증에 사용됩니다.│  ← description
│                                                │
│  ┌────────────────────────────────────────┐   │
│  │ 🔒 ••••••••••••••••••••               │   │  ← secret → 마스킹. 👁️ 토글
│  └────────────────────────────────────────┘   │
│                                                │
│  카테고리: 결제 (payment)                      │
│  서비스: payment-service                       │
│  환경변수: TOSS_SECRET_KEY                     │
│  보안 등급: 🔴 Secret (Super Admin만 변경)      │
│  재시작 필요: 아니오 ✅                         │
│                                                │
│  ─── 상세 안내 ───                              │
│  ## 발급 방법                                   │  ← help_text (마크다운)
│  1. Toss 개발자센터 로그인                      │
│  2. 내 개발정보 → API 키 메뉴 선택              │
│  3. 시크릿 키 복사                              │
│                                                │
│  > ⚠️ 테스트 키(test_sk_)와 라이브 키           │
│  > (live_sk_)를 구분하세요.                     │
│                                                │
│  ─── 변경 사유 ───                              │
│  ┌────────────────────────────────────────┐   │
│  │ (선택) 변경 사유를 입력하세요            │   │
│  └────────────────────────────────────────┘   │
│                                                │
│         [취소]        [저장]                    │
│         [🤖 AI에게 도움 요청]                   │  ← LLM 어시스턴트로 이동
└──────────────────────────────────────────────┘
```

### 6.3 LLM 채팅 인터페이스 UX

```
┌─────────────────────────────────────┬──────────────────────┐
│  🤖 설정 어시스턴트                  │  📋 제안된 변경       │
│                                     │                      │
│  [카테고리: 전체 ▼]  [언어: 한국어 ▼]│  (제안 없음)          │
│                                     │                      │
│  ┌─ 🤖 ─────────────────────────┐  │                      │
│  │ 안녕하세요! ManPaSik 설정     │  │                      │
│  │ 어시스턴트입니다.              │  │                      │
│  │ 어떤 설정을 도와드릴까요?     │  │                      │
│  └──────────────────────────────┘  │                      │
│                                     │                      │
│  ┌─ 👤 ─────────────────────────┐  │                      │
│  │ Toss 결제 설정해주세요        │  │                      │
│  └──────────────────────────────┘  │                      │
│                                     │                      │
│  ┌─ 🤖 ─────────────────────────┐  │  ┌────────────────┐ │
│  │ Toss 연동에 3개 설정이        │  │  │toss.secret_key │ │
│  │ 필요합니다: ...               │  │  │🔒 직접 입력    │ │
│  │                               │  │  │[입력] [건너뛰기]│ │
│  │ 적용하시겠습니까?             │  │  ├────────────────┤ │
│  └──────────────────────────────┘  │  │toss.sandbox    │ │
│                                     │  │true → false     │ │
│  ┌──────────────────────────────┐  │  │[적용] [거부]    │ │
│  │ 메시지 입력...          [전송]│  │  └────────────────┘ │
│  └──────────────────────────────┘  │                      │
└─────────────────────────────────────┴──────────────────────┘
```

---

## 7. 설정 동적 반영 메커니즘

### 7.1 즉시 반영 (Hot Reload)

대부분의 설정은 서비스 재시작 없이 즉시 반영 가능합니다:

```
                    admin-service
                    SetSystemConfig()
                          │
                ┌─────────▼──────────┐
                │  system_configs DB  │
                │  (값 저장)          │
                └─────────┬──────────┘
                          │
            ┌─────────────▼─────────────┐
            │  Kafka: manpasik.config.   │
            │  changed (key, new_value)  │
            └─────┬─────────────────┬───┘
                  │                 │
         ┌────────▼───┐   ┌────────▼────────┐
         │payment-svc │   │notification-svc │
         │ConfigWatcher│   │ConfigWatcher    │
         │reload()     │   │reload()         │
         └─────────────┘   └─────────────────┘
```

- **패턴**: 각 서비스에 `ConfigWatcher`가 Kafka `manpasik.config.changed` 토픽을 구독.
- **동작**: 자기 서비스에 해당하는 키가 변경되면 메모리 내 설정을 갱신.
- **예**: `toss.secret_key` 변경 → payment-service가 `pgGateway` 재초기화.

### 7.2 재시작 필요 설정

`restart_required=true`인 설정(예: Kafka brokers, Keycloak URL)은:
1. DB에 새 값 저장
2. 관리자에게 "이 설정은 서비스 재시작 후 적용됩니다" 경고
3. 재시작 전까지 이전 값으로 동작

---

## 8. 보안 설계

### 8.1 암호화 저장

- `security_level=secret`인 설정 값은 **AES-256-GCM으로 암호화**하여 `system_configs.value`에 저장
- 암호화 키: `CONFIG_ENCRYPTION_KEY` 환경변수 (32바이트, admin-service만 보유)
- 복호화: admin-service 내부에서만 수행. 다른 서비스는 Kafka로 복호화된 값을 (TLS 채널 내) 수신

### 8.2 접근 제어 매트릭스

| 보안 등급 | 조회 | 변경 | LLM 전달 |
|-----------|------|------|----------|
| `public` | 모든 관리자 | Admin 이상 | 값 포함 |
| `internal` | Admin 이상 | Admin 이상 | 값 포함 |
| `confidential` | Super Admin | Super Admin | 키 이름만 |
| `secret` | Super Admin (마스킹) | Super Admin | 키 이름만, 값 미포함 |

### 8.3 감사 로그

모든 설정 변경은 `audit_logs`에 기록:
```json
{
  "action": "config_change",
  "resource_type": "system_config",
  "resource_id": "toss.secret_key",
  "description": "설정 변경: toss.secret_key (LLM 제안, 관리자 승인)",
  "metadata": {
    "old_value": "****",
    "new_value": "****",
    "suggested_by": "llm",
    "session_id": "uuid",
    "reason": "Toss 결제 연동 설정"
  }
}
```

---

## 9. 구현 단계

### Phase 1: 기반 (Sprint 2, 1~2일)
1. DB 스키마 확장 (`25-admin-settings-ext.sql`)
2. `config_metadata`, `config_translations` 시드 데이터
3. admin-service: `ListSystemConfigs`, `GetConfigWithMeta`, `ValidateConfigValue` RPC 구현
4. admin-service: 암호화 저장/복호화 로직

### Phase 2: 동적 반영 (Sprint 2, 1일)
5. Kafka `manpasik.config.changed` 이벤트 발행 (admin-service)
6. 주요 서비스(payment, notification)에 `ConfigWatcher` 구현
7. payment-service: DB system_config에서 Toss 키 로드 (env fallback)

### Phase 3: Flutter Admin UI (Sprint 2~3, 2~3일)
8. 설정 목록 화면 (카테고리별 탭, 검색)
9. 설정 편집 다이얼로그 (유형별 입력, 마스킹, 도움말)
10. 설정 변경 이력 화면

### Phase 4: LLM 어시스턴트 (Sprint 3, 2~3일)
11. ai-inference-service: LLM 클라이언트 (OpenAI/Anthropic) 통합
12. `StartConfigSession`, `SendConfigMessage`, `ApplyConfigSuggestion` RPC 구현
13. 시스템 프롬프트 + 안전 규칙 적용
14. Flutter: LLM 채팅 UI + 제안 패널

### Phase 5: 고도화 (Sprint 4)
15. 다국어 자동 번역 (translation-service 연동: 새 설정 추가 시 미번역 언어 자동 번역)
16. 설정 변경 대기열 + 일괄 승인/거부
17. 설정 내보내기/가져오기 (JSON)
18. 설정 스냅샷 + 롤백

---

## 10. 기술 결정 요약

| 항목 | 결정 |
|------|------|
| 설정 저장소 | PostgreSQL `system_configs` + `config_metadata` (기존 테이블 확장) |
| 암호화 | AES-256-GCM, 키는 환경변수 (`CONFIG_ENCRYPTION_KEY`) |
| 동적 반영 | Kafka 이벤트 + 서비스별 ConfigWatcher |
| LLM 제공자 | OpenAI (기본), Anthropic (대안), 로컬 모델 (향후) |
| LLM 안전 | 시크릿 값 미전달, 제안→승인 2단계, 감사 로그 |
| 다국어 | DB 기반 번역 (config_translations) + translation-service 자동 번역 |
| 프론트엔드 | Flutter 관리자 화면 (Material 3, 마크다운 렌더링) |

---

## 11. 파일 위치 요약

| 파일 | 내용 |
|------|------|
| `infrastructure/database/init/25-admin-settings-ext.sql` | DB 스키마 확장 |
| `backend/services/admin-service/internal/service/config_manager.go` | 설정 관리 비즈니스 로직 |
| `backend/services/admin-service/internal/crypto/aes.go` | AES-256-GCM 암호화/복호화 |
| `backend/services/admin-service/internal/repository/postgres/config.go` | 메타·번역 리포지토리 |
| `backend/services/ai-inference-service/internal/service/config_assistant.go` | LLM 어시스턴트 로직 |
| `backend/services/ai-inference-service/internal/llm/client.go` | OpenAI/Anthropic HTTP 클라이언트 |
| `backend/shared/events/config_watcher.go` | Kafka 기반 설정 변경 감시자 |
| `frontend/flutter-app/lib/features/admin/presentation/settings_screen.dart` | 설정 관리 UI |
| `frontend/flutter-app/lib/features/admin/presentation/config_assistant_screen.dart` | LLM 채팅 UI |
| `docs/specs/admin-settings-llm-assistant-spec.md` | 이 기획서 |

---

**마지막 업데이트**: 2026-02-12 (v1.0 초안)
