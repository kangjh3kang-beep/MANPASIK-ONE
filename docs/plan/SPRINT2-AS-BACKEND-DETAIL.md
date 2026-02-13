# Sprint 2 — 관리자 설정 백엔드 세부 구현 기획서

> 범위: AS-1 ~ AS-5 | 예상 기간: 2~3일

---

## 1. 구현 순서

```
AS-1 (DB 스키마) → AS-2 (RPC 확장) → AS-3 (암호화) → AS-4 (Kafka 이벤트) → AS-5 (DB config 로드)
```

---

## 2. AS-1: DB 스키마 확장

### 파일: `infrastructure/database/init/25-admin-settings-ext.sql`

```sql
-- Custom ENUM types
CREATE TYPE config_category AS ENUM (
    'general','payment','auth','storage','messaging','database','ai','notification','security','integration'
);
CREATE TYPE config_value_type AS ENUM (
    'string','number','boolean','secret','url','email','json','select','multiline'
);
CREATE TYPE config_security_level AS ENUM (
    'public','internal','confidential','secret'
);

-- config_metadata: 설정 항목의 스키마 정의
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

-- config_translations: 다국어 설명
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

-- llm_config_sessions: LLM 대화 세션
CREATE TABLE IF NOT EXISTS llm_config_sessions (
    session_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    admin_id         VARCHAR(36) NOT NULL,
    language_code    VARCHAR(5) NOT NULL DEFAULT 'ko',
    status           VARCHAR(20) NOT NULL DEFAULT 'active',
    context_category config_category,
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    updated_at       TIMESTAMPTZ DEFAULT NOW()
);

-- llm_config_messages: 대화 메시지
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

-- config_change_queue: 설정 변경 대기열
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
```

### 시드 데이터
- system_configs INSERT: 25+ 항목 (toss.*, jwt.*, keycloak.*, s3.*, kafka.*, fcm.*, security.*, ai.*, llm.*, maintenance_mode, max_devices_per_user, default_language, session_timeout_minutes, max_file_upload_mb)
- config_metadata INSERT: 25+ 항목 매핑 (category, value_type, security_level, env_var_name, service_name, restart_required, display_order)
- config_translations INSERT: 주요 10개 키 × 3언어(ko/en/ja) = 30행

---

## 3. AS-2: admin-service RPC 확장

### 3.1 Proto 확장 (`backend/shared/proto/manpasik.proto`)

```protobuf
// --- AdminService 확장 ---
message ListSystemConfigsRequest {
  string language_code = 1;
  string category = 2;
  bool include_secrets = 3;
}
message ListSystemConfigsResponse {
  repeated ConfigWithMeta configs = 1;
  map<string, int32> category_counts = 2;
}
message ConfigWithMeta {
  string key = 1;
  string value = 2;
  string raw_value = 3;
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
  string display_name = 20;
  string description = 21;
  string placeholder = 22;
  string help_text = 23;
  string validation_message = 24;
  string updated_by = 30;
  google.protobuf.Timestamp updated_at = 31;
}
message GetConfigWithMetaRequest {
  string key = 1;
  string language_code = 2;
}
message ValidateConfigValueRequest {
  string key = 1;
  string value = 2;
}
message ValidateConfigValueResponse {
  bool valid = 1;
  string error_message = 2;
  repeated string suggestions = 3;
}
message BulkSetConfigsRequest {
  repeated SetSystemConfigRequest configs = 1;
  string reason = 2;
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

### 3.2 리포지토리 인터페이스

#### `backend/services/admin-service/internal/repository/config_meta.go`

```go
type ConfigMetadataRepository interface {
    GetByKey(ctx context.Context, key string) (*models.ConfigMetadata, error)
    ListByCategory(ctx context.Context, category string) ([]*models.ConfigMetadata, error)
    ListAll(ctx context.Context) ([]*models.ConfigMetadata, error)
    CountByCategory(ctx context.Context) (map[string]int32, error)
}
```

#### `backend/services/admin-service/internal/repository/config_translation.go`

```go
type ConfigTranslationRepository interface {
    GetByKeyAndLang(ctx context.Context, key, lang string) (*models.ConfigTranslation, error)
    ListByKey(ctx context.Context, key string) ([]*models.ConfigTranslation, error)
    ListByLang(ctx context.Context, lang string) ([]*models.ConfigTranslation, error)
}
```

#### PostgreSQL 구현: `internal/repository/postgres/config_meta.go`, `postgres/config_translation.go`
#### Memory 구현: `internal/repository/memory/config_meta.go`, `memory/config_translation.go`

### 3.3 서비스 레이어

#### `backend/services/admin-service/internal/service/config_manager.go`

```go
type ConfigManager struct {
    configRepo      repository.SystemConfigRepository
    metaRepo        repository.ConfigMetadataRepository
    translationRepo repository.ConfigTranslationRepository
    encryptor       *crypto.AESEncryptor   // AS-3
    eventPublisher  events.EventPublisher  // AS-4
}

func (cm *ConfigManager) ListSystemConfigs(ctx context.Context, lang, category string, includeSecrets bool) (*pb.ListSystemConfigsResponse, error)
func (cm *ConfigManager) GetConfigWithMeta(ctx context.Context, key, lang string) (*pb.ConfigWithMeta, error)
func (cm *ConfigManager) ValidateConfigValue(ctx context.Context, key, value string) (*pb.ValidateConfigValueResponse, error)
func (cm *ConfigManager) BulkSetConfigs(ctx context.Context, configs []*pb.SetSystemConfigRequest, reason string) (*pb.BulkSetConfigsResponse, error)
```

**ValidateConfigValue 의사코드:**
```
1. metaRepo.GetByKey(key)
2. value_type == "number" → ParseFloat → check min/max
3. value_type == "boolean" → check "true"/"false"
4. value_type == "url" → url.Parse
5. value_type == "email" → regexp 검증
6. value_type == "select" → check allowed_values 포함
7. validation_regex != "" → regexp.MatchString
8. return valid/error_message/suggestions
```

### 3.4 gRPC 핸들러

#### `backend/services/admin-service/internal/handler/grpc.go` 확장

```go
func (h *AdminHandler) ListSystemConfigs(ctx context.Context, req *pb.ListSystemConfigsRequest) (*pb.ListSystemConfigsResponse, error)
func (h *AdminHandler) GetConfigWithMeta(ctx context.Context, req *pb.GetConfigWithMetaRequest) (*pb.ConfigWithMeta, error)
func (h *AdminHandler) ValidateConfigValue(ctx context.Context, req *pb.ValidateConfigValueRequest) (*pb.ValidateConfigValueResponse, error)
func (h *AdminHandler) BulkSetConfigs(ctx context.Context, req *pb.BulkSetConfigsRequest) (*pb.BulkSetConfigsResponse, error)
```

---

## 4. AS-3: AES-256-GCM 암호화

### 파일: `backend/services/admin-service/internal/crypto/aes.go`

```go
type AESEncryptor struct {
    key []byte // 32 bytes (AES-256)
}

func NewAESEncryptor(keyHex string) (*AESEncryptor, error)
func (e *AESEncryptor) Encrypt(plaintext string) (string, error)   // → base64(nonce+ciphertext)
func (e *AESEncryptor) Decrypt(ciphertext string) (string, error)  // base64 decode → AES-GCM open
```

**Encrypt 의사코드:**
```
1. aes.NewCipher(key)
2. cipher.NewGCM(block)
3. nonce = make([]byte, gcm.NonceSize()); io.ReadFull(rand.Reader, nonce)
4. sealed = gcm.Seal(nonce, nonce, plaintext, nil)
5. return base64.StdEncoding.EncodeToString(sealed)
```

### 통합: `config_manager.go`

```go
func (cm *ConfigManager) SetSystemConfig(ctx, key, value, adminID string) error {
    meta := cm.metaRepo.GetByKey(ctx, key)
    if meta.SecurityLevel == "secret" && cm.encryptor != nil {
        value = cm.encryptor.Encrypt(value)
    }
    cm.configRepo.Set(ctx, key, value, adminID)
    cm.eventPublisher.Publish("manpasik.config.changed", ConfigChangedEvent{Key: key, ...})
}
```

### 환경변수
- `CONFIG_ENCRYPTION_KEY`: 64자리 hex (32 bytes) — admin-service만 보유

---

## 5. AS-4: Kafka config.changed 이벤트 + ConfigWatcher

### 5.1 이벤트 발행 (admin-service)

#### 이벤트 구조

```go
type ConfigChangedEvent struct {
    Key          string `json:"key"`
    NewValue     string `json:"new_value"`     // secret이면 복호화된 값 (TLS 내 전송)
    OldValue     string `json:"old_value"`
    ChangedBy    string `json:"changed_by"`
    ServiceName  string `json:"service_name"`  // 해당 설정의 대상 서비스
    Timestamp    string `json:"timestamp"`
}
```

- 토픽: `manpasik.config.changed`
- SetSystemConfig, BulkSetConfigs 완료 후 발행

### 5.2 ConfigWatcher 인터페이스

#### `backend/shared/events/config_watcher.go`

```go
type ConfigChangeHandler func(key, newValue string) error

type ConfigWatcher interface {
    Watch(ctx context.Context, serviceFilter string, handler ConfigChangeHandler) error
    Close() error
}

// KafkaConfigWatcher: Kafka consumer 구현
type KafkaConfigWatcher struct { ... }
func NewKafkaConfigWatcher(brokers []string) (*KafkaConfigWatcher, error)
func (w *KafkaConfigWatcher) Watch(ctx, serviceFilter, handler) error
func (w *KafkaConfigWatcher) Close() error

// NoopConfigWatcher: 개발용
type NoopConfigWatcher struct{}
```

### 5.3 서비스별 적용

#### payment-service

```go
// cmd/main.go
watcher := events.NewKafkaConfigWatcher(cfg.Kafka.Brokers)
watcher.Watch(ctx, "payment-service", func(key, newValue string) error {
    switch key {
    case "toss.secret_key":
        newClient := pg.NewTossClient(newValue, currentAPIURL)
        svc.SetPaymentGateway(newClient)
    case "toss.api_url":
        // ...
    }
    return nil
})
```

#### notification-service

```go
watcher.Watch(ctx, "notification-service", func(key, newValue string) error {
    switch key {
    case "fcm.server_key":
        newFCM := fcm.NewClient(newValue, projectID)
        svc.SetFCMClient(newFCM)
    }
    return nil
})
```

---

## 6. AS-5: DB Config 우선 로드

### 공통 패턴: `backend/shared/config/db_loader.go`

```go
func LoadConfigFromDB(pool *pgxpool.Pool, key string) (string, error) {
    var value string
    err := pool.QueryRow(ctx, "SELECT value FROM system_configs WHERE key=$1", key).Scan(&value)
    if err != nil { return "", err }
    return value, nil
}

func LoadConfigWithFallback(pool *pgxpool.Pool, dbKey, envVar string) string {
    if val, err := LoadConfigFromDB(pool, dbKey); err == nil && val != "" {
        return val
    }
    return os.Getenv(envVar)
}
```

### payment-service 적용

```go
// cmd/main.go
tossSecret := config.LoadConfigWithFallback(pool, "toss.secret_key", "TOSS_SECRET_KEY")
tossAPIURL := config.LoadConfigWithFallback(pool, "toss.api_url", "TOSS_API_URL")
if tossSecret != "" {
    svc.SetPaymentGateway(pg.NewTossClient(tossSecret, tossAPIURL))
}
```

---

## 7. 테스트 전략

| 파일 | 테스트 항목 |
|------|------------|
| `crypto/aes_test.go` | Encrypt/Decrypt 왕복, 잘못된 키, 빈 문자열, 유니코드 |
| `service/config_manager_test.go` | ListConfigs(카테고리 필터), GetConfigWithMeta(번역 포함), ValidateValue(각 유형별), BulkSet(성공/부분실패), 암호화 저장 확인 |
| `repository/postgres/config_meta_test.go` | CRUD, 카테고리별 조회, CountByCategory |
| `repository/postgres/config_translation_test.go` | 키+언어 조회, 없는 언어 폴백 |
| `events/config_watcher_test.go` | Mock Kafka consumer, serviceFilter 동작, 핸들러 호출 확인 |
| `config/db_loader_test.go` | DB 값 존재/미존재, env fallback |

**예상 테스트 수:** 30+

---

## 8. 검증 기준

```bash
# 빌드
cd backend && go build ./...

# 정적 분석
cd backend && go vet ./...

# 테스트
cd backend && go test ./services/admin-service/... -v -count=1
cd backend && go test ./shared/events/... -v -count=1
cd backend && go test ./shared/config/... -v -count=1

# 전체
cd backend && go test ./... -count=1
```

---

## 9. 수정·생성 파일 목록

| 상태 | 파일 |
|------|------|
| 신규 | `infrastructure/database/init/25-admin-settings-ext.sql` |
| 신규 | `backend/services/admin-service/internal/crypto/aes.go` |
| 신규 | `backend/services/admin-service/internal/crypto/aes_test.go` |
| 신규 | `backend/services/admin-service/internal/repository/config_meta.go` |
| 신규 | `backend/services/admin-service/internal/repository/config_translation.go` |
| 신규 | `backend/services/admin-service/internal/repository/postgres/config_meta.go` |
| 신규 | `backend/services/admin-service/internal/repository/postgres/config_translation.go` |
| 신규 | `backend/services/admin-service/internal/repository/memory/config_meta.go` |
| 신규 | `backend/services/admin-service/internal/repository/memory/config_translation.go` |
| 신규 | `backend/services/admin-service/internal/service/config_manager.go` |
| 신규 | `backend/services/admin-service/internal/service/config_manager_test.go` |
| 신규 | `backend/shared/events/config_watcher.go` |
| 신규 | `backend/shared/events/config_watcher_test.go` |
| 신규 | `backend/shared/config/db_loader.go` |
| 신규 | `backend/shared/config/db_loader_test.go` |
| 수정 | `backend/shared/proto/manpasik.proto` |
| 수정 | `backend/services/admin-service/internal/handler/grpc.go` |
| 수정 | `backend/services/admin-service/cmd/main.go` |
| 수정 | `backend/services/payment-service/cmd/main.go` |
| 수정 | `backend/services/notification-service/cmd/main.go` |
| 수정 | `backend/shared/config/config.go` |

**마지막 업데이트**: 2026-02-12
