# Phase A' — PostgreSQL 영속화 상세 구현 계획서

**문서 ID**: MPS-PLAN-PHASE-A'-v1.0
**작성일**: 2026-04-25
**상태**: 검증 완료 — 구조적 영속화 100% 달성

---

## 1. 현황 요약 (2026-04-25 실사)

### 1.1 Phase A' 달성률: **100% (구조적 완료)**

| 항목 | 기대치 | 실측 | 상태 |
|------|--------|------|------|
| PostgreSQL 저장소 파일 | 34/35 서비스 | 39파일/34서비스 | ✅ |
| SQL 쿼리 구현 | 모든 인터페이스 메서드 | 100% 매칭 | ✅ |
| main.go DB_HOST 스위칭 | 34/34 서비스 | 34/34 서���스 | ✅ |
| 공통 config (DSN) | config.go | ✅ DatabaseConfig.DSN() | ✅ |
| DB config loader | db_loader.go | ✅ system_configs 테이블 | ✅ |
| 전체 빌드 | 35/35 PASS | 35/35 PASS | ✅ |
| 전체 테스트 | 35/35 PASS | 35/35 PASS | ✅ |

### 1.2 초기 "GAP" 재분류

초기 분석에서 7개 서비스가 "postgres 메서드 부족"으로 분류되었으나, 상세 분석 결과 **실제 PostgreSQL 갭은 0건**:

| 서비스 | 누락으로 보인 항목 | 실제 소속 | DB 갭 여부 |
|--------|-------------------|----------|-----------|
| admin-service | InMemoryAuditLogStore | 내부 레거시 스토어 (서비스 인터페이스 아님) | ❌ |
| audit-service | Seed() | 개발용 시드 데이터 (인터페이스 아님) | ��� |
| auth-service | TokenRepository | Redis 전용 (의도적 설계) | ❌ |
| device-service | EventPublisher, SubscriptionChecker | Kafka/gRPC 서비스 호출 | ❌ |
| emergency-service | idFromInt, copyStrings | 내부 헬퍼 함수 | ❌ |
| measurement-service | SearchIndexer, VectorRepo, EventPublisher | Elasticsearch/Milvus/Kafka | ❌ |
| subscription-service | EventPublisher | Kafka 이벤트 발행 | ❌ |

---

## 2. 아키텍처 패턴 (확정)

### 2.1 저장소 스위칭 패턴 (DB_HOST 기반)
```go
// cmd/main.go — 모든 34서비스 공통 패턴
var repo service.XxxRepository

if _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
    connCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
    pool, poolErr := pgxpool.New(connCtx, cfg.DB.DSN())
    connCancel()
    if poolErr != nil {
        log.Printf("[%s] DB connection failed, using memory: %v", serviceName, poolErr)
        repo = memory.NewXxxRepository()
    } else {
        pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
        if pingErr := pool.Ping(pingCtx); pingErr != nil {
            pingCancel()
            pool.Close()
            repo = memory.NewXxxRepository()
        } else {
            pingCancel()
            defer pool.Close()
            repo = postgres.NewXxxRepository(pool)
        }
    }
} else {
    repo = memory.NewXxxRepository()
}
```

### 2.2 공통 설정 (`backend/shared/config/`)
- `config.go`: `DatabaseConfig.DSN()` → `postgres://user:pass@host:port/dbname?sslmode=disable`
- `db_loader.go`: `LoadConfigFromDB(pool, key)` → system_configs 테이블 폴백
- 환경변수: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`

### 2.3 PostgreSQL 저장소 패턴
```go
// internal/repository/postgres/xxx.go
type XxxRepository struct {
    pool *pgxpool.Pool
}

func NewXxxRepository(pool *pgxpool.Pool) *XxxRepository {
    return &XxxRepository{pool: pool}
}

func (r *XxxRepository) GetByID(ctx context.Context, id string) (*service.Xxx, error) {
    const q = `SELECT ... FROM xxx WHERE id = $1`
    var x service.Xxx
    err := r.pool.QueryRow(ctx, q, id).Scan(&x.Field1, &x.Field2, ...)
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, nil  // Not found → nil, nil (memory 호환)
        }
        return nil, err
    }
    return &x, nil
}
```

---

## 3. 서비스별 PostgreSQL 구현 현황

### 3.1 완전 구현 (34/34 서비스)

| # | 서비스 | Postgres 파일 | 메서드 수 | DB 스키마 |
|---|--------|--------------|----------|----------|
| 1 | admin-service | admin.go, config_meta.go | 26 | 18-admin.sql |
| 2 | ai-inference-service | inference.go | 11 | - |
| 3 | analytics-service | analytics.go | 6 | 34-analytics.sql |
| 4 | assistant-service | assistant.go | 8 | 26-assistant.sql |
| 5 | audit-service | audit.go | 4 | 39-audit.sql |
| 6 | auth-service | user.go | 5 | 01-auth.sql |
| 7 | calibration-service | calibration.go | 14 | 08-calibration.sql |
| 8 | cartridge-service | cartridge.go | 7 | 09-cartridge.sql |
| 9 | cartridge-store-service | store.go | 14 | 29-cartridge-store.sql |
| 10 | coaching-service | coaching.go | 20 | 10-coaching.sql |
| 11 | community-service | community.go | 22 | 17-community.sql |
| 12 | concept-service | concept.go | 13 | 28-concept-org.sql |
| 13 | data-platform-service | platform.go | 12 | 31-data-platform.sql |
| 14 | device-service | device.go | 9 | 05-device.sql |
| 15 | digital-twin-service | twin.go | 4 | 40-digital-twin.sql |
| 16 | emergency-service | emergency.go | 10 | 35-emergency.sql |
| 17 | family-service | family.go | 18 | 13-family.sql |
| 18 | health-record-service | healthrecord.go | 23 | 14-health-record.sql |
| 19 | iot-gateway-service | iot_gateway.go | 10 | 36-iot-gateway.sql |
| 20 | marketplace-service | marketplace.go | 16 | 37-marketplace.sql |
| 21 | measurement-service | measurement.go, session.go | 8 | 04-measurement.sql |
| 22 | nlp-service | nlp.go | 6 | 38-nlp.sql |
| 23 | notification-service | notification.go | 14 | 12-notification.sql |
| 24 | payment-service | payment.go | 13 | 06-payment.sql |
| 25 | prescription-service | prescription.go | 20 | 19-prescription.sql |
| 26 | reservation-service | reservation.go | 23 | 11-reservation.sql |
| 27 | shop-service | shop.go | 17 | 07-shop.sql |
| 28 | subscription-service | subscription.go | 8 | 03-subscription.sql |
| 29 | telemedicine-service | telemedicine.go | 20 | 22-telemedicine.sql |
| 30 | translation-service | translation.go | 6 | 20-translation.sql |
| 31 | user-service | family.go, profile.go, subscription.go | 14 | 02-user.sql |
| 32 | video-service | video.go | 12 | 21-video.sql |
| 33 | vision-service | vision.go | 6 | 27-vision-diet.sql |
| 34 | voice-profile-service | voice_profile.go | 6 | 33-voice-profile.sql |

### 3.2 비-DB 인프라 (Phase C 영역)

다음은 PostgreSQL이 아닌 외부 인프라로, Phase C에서 구현:

| 구성요소 | 대상 서비스 | 외부 시스템 | Phase |
|---------|-----------|-----------|-------|
| TokenRepository (Redis) | auth-service | Redis | C-1 |
| EventPublisher (Kafka) | device, measurement, subscription, payment | Kafka/Redpanda | C-2 |
| SearchIndexer | measurement-service | Elasticsearch | C-3 |
| VectorRepository | measurement-service | Milvus | C-3 |
| SubscriptionChecker | device-service | gRPC 호출 | B-1 |

---

## 4. 남은 작업 (Phase A' → Phase A'' 최적화)

Phase A'의 구조적 영속화는 완료되었으나, 운영 품질을 위한 추가 작업:

### 4.1 공통 DB 커넥션 헬퍼 (선택적 개선)

현재 34개 main.go에 DB 연결 보일러플레이트가 중복됩니다. 공통 헬퍼 추출 가능:

```go
// backend/shared/database/pool.go (신규)
package database

func NewPool(cfg *config.DatabaseConfig) (*pgxpool.Pool, error) {
    if _, set := os.LookupEnv("DB_HOST"); !set || cfg.Host == "" {
        return nil, ErrNoDBConfig
    }
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    pool, err := pgxpool.New(ctx, cfg.DSN())
    if err != nil {
        return nil, err
    }
    pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer pingCancel()
    if err := pool.Ping(pingCtx); err != nil {
        pool.Close()
        return nil, err
    }
    return pool, nil
}
```

**우선순위**: 낮음 (현재 코드 정상 작동)

### 4.2 DB 마이그레이션 도구 (Phase G 영역)

- `golang-migrate` 또는 `goose` 도입
- 42개 SQL init 스크립트 → 버전화된 마이그레이션으로 전환
- **우선순위**: Phase G (운영/배포)

### 4.3 통합 테스트 (Phase F 영역)

- testcontainers-go를 활용한 PostgreSQL 통합 테스트
- 각 서비스의 postgres repo에 `*_integration_test.go` 추가
- **우선순위**: Phase F (테스트/품질)

---

## 5. 결론 및 다음 단계

### Phase A' 최종 판정: ✅ 구조적 완료 (100%)

- **PostgreSQL 저장소**: 34/34 서비스 구현 (39 파일, 모든 인터페이스 메서드)
- **DB 스위칭**: 34/34 main.go에 DB_HOST 환경변수 기반 폴백
- **공통 Config**: DSN, db_loader 완비
- **빌드/테스트**: 35/35 ALL PASS

### 다음 Phase 우선순위

1. **Phase B** (SKELETON→REAL): analytics, emergency, assistant, concept, data-platform, gateway
2. **Phase C** (외부 연동): Redis TokenRepo, Kafka EventPublisher, Milvus VectorRepo
3. **Phase F** (테스트): PostgreSQL 통합 테스트 추가

---

*자동 생성: 2026-04-25 | 만파식(萬波息) Phase A' 영속화 상세 구현 계획서*
