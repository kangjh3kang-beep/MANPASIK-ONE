# ManPaSik (만파식) AI Agent Team - Master Context & Orchestration

> 이 파일은 Cursor IDE 에이전트 팀의 최상위 오케스트레이션 문서입니다.
> 모든 에이전트가 이 파일을 참조하여 프로젝트 전체 컨텍스트를 이해합니다.

---

## 1. 프로젝트 아이


















































덴티티

- **프로젝트명**: ManPaSik (만파식, 萬波息) AI Healthcare Ecosystem
- **철학**: 홍익인간(弘益人間) - 세상의 모든 파동을 분석해 인간을 이롭게 한다
- **제품 유형**: 체외진단의료기기(IVD) + 헬스케어 SaaS 플랫폼
- **의료기기 등급**: Class II (중등도 위험)
- **대상 시장**: 한국(MFDS), 미국(FDA 510(k)), EU(CE-IVDR), 중국(NMPA), 일본(PMDA)
- **라이선스**: Proprietary - All Rights Reserved
- **버전**: 0.1.0 (MVP 개발 중)
- **로드맵**: 24개월 5단계 (MVP → Core → Advanced → Ecosystem → Evolution)
- **예상 규모**: 피크 인력 32명, 예산 ~67억원

---

## 2. 핵심 기술 깊이 이해

### 2.1 차동측정법 (Differential Measurement)
```
핵심 공식: S_corrected = S_det - α × S_ref

• S_det     : 검출 전극(Detection electrode) 신호
• S_ref     : 기준 전극(Reference electrode) 신호
• α (alpha) : 보정 계수 (기본값 = 0.95, crate::DEFAULT_ALPHA)

효과: 99% 매트릭스 노이즈 제거
```

**채널별 보정 확장:**
```
S_corrected[i] = (S_det[i] - α × S_ref[i] - offset[i]) × gain[i]
```

**파라미터 구조 (CorrectionParams):**
- `alpha: f64` — 보정 계수 (0.0~1.0)
- `channel_offsets: Vec<f64>` — 채널별 DC 오프셋
- `channel_gains: Vec<f64>` — 채널별 게인
- `temp_coefficient: f64` — 온도 보정 계수

### 2.2 핑거프린트 벡터 시스템
```
차원 확장 경로:
  88차원 (Basic)  → 단일 센서 기본 측정
  448차원 (Enhanced) → 전자코(8ch) + 전자혀(8ch) 융합
  896차원 (Full)   → 완전 융합 (MAX_CHANNELS = 896)
```

**벡터 연산:**
- L2 정규화 (`normalize()`)
- 코사인 유사도 (`cosine_similarity()`)
- 유클리드 거리 (`euclidean_distance()`)
- Milvus 벡터DB 저장 (`to_milvus_vector()`)

**FingerprintBuilder 패턴:**
```
FingerprintBuilder::new(base_88ch)
    .with_e_nose(e_nose_channels)
    .with_e_tongue(e_tongue_channels)
    .build() → FingerprintVector (88 | 448 | 896)
```

### 2.3 카트리지 시스템 (29종)
| 코드 | 카테고리 | 타입 | 채널 | 측정시간 |
|------|---------|------|------|---------|
| 0x01-0x0E | 건강 바이오마커 | Glucose, LipidPanel, HbA1c, UricAcid, Creatinine, VitaminD, VitaminB12, Ferritin, Tsh, Cortisol, Testosterone, Estrogen, Crp, Insulin | 88 | 15초 |
| 0x20-0x23 | 환경 모니터링 | WaterQuality, IndoorAirQuality, Radon, Radiation | 88 | 15초 |
| 0x30-0x33 | 식품 안전 | PesticideResidue, FoodFreshness, Allergen, DateDrug | 88 | 15초 |
| 0x40-0x42 | 전자코/전자혀 | ENose, ETongue, EhdGas | 8 | 30초 |
| 0x50-0x52 | 고급 분석 | NonTarget448, NonTarget896, MultiBiomarker | 448/896/88 | 45-90초 |
| 0xFF | 연구용 | CustomResearch | 가변 | 가변 |

### 2.4 AI 추론 엔진
| 모델 타입 | 입력 크기 | 출력 크기 | 용도 |
|----------|----------|----------|------|
| Calibration | 88 | 88 | 채널별 보정값 |
| FingerprintClassifier | 896 | 29 | 카트리지 타입 분류 |
| AnomalyDetection | 88 | 1 | 이상 스코어 |
| ValuePredictor | 88 | 1 | 단일 값 예측 |
| QualityAssessment | 88 | 3 | 품질 등급 (좋음/보통/나쁨) |

### 2.5 BLE 프로토콜
**GATT 서비스:**
| UUID | 서비스 |
|------|--------|
| `0000fff0-...-34fb` | 메인 측정 서비스 |
| `0000180a-...-34fb` | 디바이스 정보 서비스 |
| `0000180f-...-34fb` | 배터리 서비스 |

**GATT 특성:**
| UUID | 특성 | 모드 |
|------|------|------|
| `0000fff1-...` | 측정 명령 | Write |
| `0000fff2-...` | 측정 데이터 | Notify |
| `0000fff3-...` | 측정 상태 | Read/Notify |
| `0000fff4-...` | 보정 데이터 | Read/Write |
| `00002a26-...` | 펌웨어 버전 | Read |
| `00002a19-...` | 배터리 레벨 | Read |

**BLE 명령 코드:**
| 코드 | 명령 |
|------|------|
| 0x01 | StartMeasurement |
| 0x02 | StopMeasurement |
| 0x03 | GetStatus |
| 0x04 | StartCalibration |
| 0x05 | SetParameters |
| 0xFF | Reset |

**측정 데이터 패킷 (바이너리):**
```
[0-1]  : 시퀀스 번호 (u16 LE)
[2-3]  : 채널 수 (u16 LE)
[4-7]  : 온도 (f32 LE, 섭씨)
[8-11] : 습도 (f32 LE, %)
[12]   : 배터리 (u8, %)
[13-16]: 타임스탬프 하위 4바이트 (u32 LE, ms)
[17+]  : 채널 데이터 (f32 LE × 채널수)
```

### 2.6 NFC 카트리지 태그 구조
```
[0-7]  : 카트리지 ID (8바이트 UID)
[8]    : 카트리지 타입 코드 (CartridgeType::to_code)
[9-16] : 로트 ID (8바이트 ASCII)
[17-24]: 유효 기간 (YYYYMMDD)
[25-26]: 잔여 사용 횟수 (u16 LE)
[27-28]: 최대 사용 횟수 (u16 LE)
[29-36]: α 계수 (f64 LE)
[37-44]: 온도 보정 계수 (f64 LE)
[45-52]: 습도 보정 계수 (f64 LE)
[53+]  : 추가 보정 데이터
```

---

## 3. 시스템 아키텍처 (확정)

```
┌──────────────────────────────────────────────────────────────────┐
│                        클라이언트 계층                             │
│  ┌─────────────────────┐    ┌──────────────────────────┐        │
│  │ Flutter 앱           │    │ Next.js 14 웹 (PWA)      │        │
│  │ iOS/Android/Desktop │    │ 관리자/의료진 대시보드     │        │
│  │ Riverpod 2.x        │    │ TypeScript + Tailwind    │        │
│  └──────────┬──────────┘    └────────────┬─────────────┘        │
│             │ FFI (flutter_rust_bridge)   │ HTTPS/WebSocket      │
├─────────────┼────────────────────────────┼──────────────────────┤
│             ▼                            │                       │
│  ┌──────────────────────┐                │                       │
│  │ Rust Core Engine     │  코어 계층      │                       │
│  │ (manpasik-engine)    │                │                       │
│  │ ┌─────────────────┐  │                │                       │
│  │ │ differential     │  │  S_det-α×S_ref│                       │
│  │ │ fingerprint      │  │  88→448→896   │                       │
│  │ │ ai (TFLite)      │  │  엣지 추론    │                       │
│  │ │ ble (btleplug)   │  │  GATT 통신    │                       │
│  │ │ nfc              │  │  카트리지      │                       │
│  │ │ dsp (rustfft)    │  │  FFT/필터     │                       │
│  │ │ crypto (ring)    │  │  AES-256+SHA  │                       │
│  │ │ sync (CRDT)      │  │  오프라인      │                       │
│  │ └─────────────────┘  │                │                       │
│  └──────────────────────┘                │                       │
├──────────────────────────────────────────┼──────────────────────┤
│                   API Gateway 계층        │                       │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │ Kong 3.7 (인증, 레이트리밋, 라우팅, SSL 종단)            │    │
│  │ → Keycloak 25.0 (OIDC/OAuth2.0, MFA, RBAC)             │    │
│  └─────────────────────────┬───────────────────────────────┘    │
│                             │ gRPC + mTLS                        │
├─────────────────────────────┼────────────────────────────────────┤
│                  마이크로서비스 계층 (Go 1.22+)                    │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌─────────────────┐    │
│  │ auth     │ │ user     │ │ device   │ │ measurement     │    │
│  │ :50051   │ │ :50052   │ │ :50053   │ │ :50054          │    │
│  └────┬─────┘ └────┬─────┘ └────┬─────┘ └──────┬──────────┘    │
│       │            │            │               │                │
│  ┌────┼────────────┼────────────┼───────────────┼───────────┐   │
│  │    ▼            ▼            ▼               ▼           │   │
│  │              데이터 계층                                   │   │
│  │  ┌───────────┐ ┌────────────┐ ┌──────┐ ┌──────┐         │   │
│  │  │PostgreSQL │ │TimescaleDB │ │Milvus│ │Redis │         │   │
│  │  │16 :5432   │ │:5433       │ │:19530│ │:6379 │         │   │
│  │  │메인 DB    │ │시계열      │ │벡터  │ │캐시  │         │   │
│  │  └───────────┘ └────────────┘ └──────┘ └──────┘         │   │
│  │  ┌──────────┐ ┌────────────┐ ┌───────┐ ┌───────┐       │   │
│  │  │Kafka     │ │Elasticsearch│ │MinIO │ │MQTT   │       │   │
│  │  │:19092    │ │:9200       │ │:9010  │ │:1883  │       │   │
│  │  │이벤트    │ │검색/로그   │ │파일   │ │IoT    │       │   │
│  │  └──────────┘ └────────────┘ └───────┘ └───────┘       │   │
│  └──────────────────────────────────────────────────────────┘   │
├──────────────────────────────────────────────────────────────────┤
│                    관측성 계층                                     │
│  Prometheus :9090 → Grafana :3000 → Alerting                     │
│  OpenTelemetry → Jaeger (분산 트레이싱)                           │
│  ELK Stack (구조화된 로그)                                        │
└──────────────────────────────────────────────────────────────────┘
```

---

## 4. gRPC API 상세 (Proto 정의)

### 4.1 MeasurementService (manpasik.v1)
| RPC | 타입 | 요청 | 응답 |
|-----|------|------|------|
| StartSession | Unary | device_id, cartridge_id, user_id | session_id, started_at |
| StreamMeasurement | Bidi-Stream | session_id, raw_channels[], differential{s_det,s_ref,alpha,s_corrected}, env_meta{temp_c,humidity_pct,pressure_kpa} | session_id, primary_value, unit, confidence, fingerprint_vector[] |
| EndSession | Unary | session_id | session_id, total_measurements, ended_at |
| GetMeasurementHistory | Unary | user_id, start_time, end_time, limit, offset | measurements[], total_count |

### 4.2 DeviceService (manpasik.v1)
| RPC | 타입 | 요청 | 응답 |
|-----|------|------|------|
| RegisterDevice | Unary | device_id, serial_number, firmware_version, user_id | device_id, registration_token |
| ListDevices | Unary | user_id | devices[] |
| StreamDeviceStatus | Bidi-Stream | device_id, status, battery_percent, signal_strength | command_id, command_type, payload |
| RequestOtaUpdate | Unary | device_id, target_version | update_id, download_url, checksum |

### 4.3 UserService (manpasik.v1)
| RPC | 타입 | 요청 | 응답 |
|-----|------|------|------|
| GetProfile | Unary | user_id | email, display_name, avatar_url, language, timezone, subscription_tier |
| UpdateProfile | Unary | user_id, display_name, ... | UserProfile |
| GetSubscription | Unary | user_id | tier, started_at, expires_at, max_devices, max_family_members, ai_coaching_enabled, telemedicine_enabled |

### 4.4 Health (manpasik.health.v1)
| RPC | 타입 | 용도 |
|-----|------|------|
| Check | Unary | 서비스 상태 확인 |
| Watch | Server-Stream | 상태 감시 |
| GetStatus | Unary | 상세 상태 (메타데이터, 의존성, 메트릭) |
| CheckDependencies | Unary | DB/캐시/큐 상태 |
| Ready | Unary | K8s readiness probe |
| Live | Unary | K8s liveness probe |

### 4.5 구독 티어
| Enum | 이름 | 가격 |
|------|------|------|
| SUBSCRIPTION_TIER_FREE (0) | Free | 무료 |
| SUBSCRIPTION_TIER_BASIC (1) | Basic Safety | ₩9,900/월 |
| SUBSCRIPTION_TIER_PRO (2) | Bio-Optimization | ₩29,900/월 |
| SUBSCRIPTION_TIER_CLINICAL (3) | Clinical Guard | ₩59,900/월 |

### 4.6 디바이스 상태 Enum
```
DEVICE_STATUS_UNKNOWN(0), ONLINE(1), OFFLINE(2), MEASURING(3), UPDATING(4), ERROR(5)
```

### 4.7 명령 타입 Enum
```
COMMAND_TYPE_UNKNOWN(0), START_MEASUREMENT(1), STOP_MEASUREMENT(2),
CALIBRATE(3), REBOOT(4), OTA_UPDATE(5)
```

---

## 5. 데이터베이스 스키마

### 5.1 PostgreSQL (메인 DB, :5432)
```sql
-- 사용자
CREATE TABLE users (
    id UUID PRIMARY KEY, email VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(100), avatar_url TEXT,
    language VARCHAR(10) DEFAULT 'ko', timezone VARCHAR(50) DEFAULT 'Asia/Seoul',
    created_at/updated_at TIMESTAMPTZ
);

-- 구독
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY, user_id UUID FK→users,
    tier VARCHAR(20) ['free','basic','pro','clinical'],
    started_at/expires_at TIMESTAMPTZ, max_devices INT, max_family_members INT
);

-- 가족 그룹
CREATE TABLE family_groups (id UUID PK, name VARCHAR(100), owner_id UUID FK→users);
CREATE TABLE family_members (id UUID PK, family_id FK, user_id FK, role VARCHAR(20), UNIQUE(family_id,user_id));

-- 디바이스
CREATE TABLE devices (
    id UUID PK, device_id VARCHAR(100) UNIQUE, user_id FK→users,
    name, firmware_version, last_seen, status, battery_percent
);
CREATE TABLE device_events (id UUID PK, device_id FK→devices, event_type, payload JSONB);
```

### 5.2 TimescaleDB (시계열, :5433)
```sql
CREATE TABLE measurements (
    time TIMESTAMPTZ NOT NULL, session_id UUID, device_id UUID, user_id UUID,
    cartridge_type VARCHAR(50),
    raw_channels DOUBLE PRECISION[], s_det/s_ref/alpha/s_corrected DOUBLE PRECISION,
    primary_value DOUBLE PRECISION, unit VARCHAR(20), confidence FLOAT,
    fingerprint_vector REAL[], -- 88/448/896 차원
    temp_c FLOAT, humidity_pct FLOAT, battery_pct INT
);
SELECT create_hypertable('measurements', 'time');
CREATE INDEX idx_measurements_session ON measurements (session_id);
CREATE INDEX idx_measurements_user ON measurements (user_id, time DESC);
```

### 5.3 Milvus (벡터 DB, :19530)
- 핑거프린트 유사도 검색 (코사인/유클리드)
- 차원: 88, 448, 896

### 5.4 Redis (캐시, :6379)
- JWT 세션, 임시 측정 데이터, 디바이스 상태 캐시

---

## 6. 에이전트 팀 구성 및 역할

### 6.1 에이전트 목록
| 에이전트 | 규칙 파일 | 담당 디렉토리 | 전문 영역 |
|---------|----------|-------------|-----------|
| **Global** | `manpasik-project.mdc` | `**/*` | 프로젝트 공통 원칙, 보안, TDD |
| **Rust Core** | `rust-core.mdc` | `rust-core/**` | 차동측정, 핑거프린트, BLE, NFC, AI, DSP, 암호화, CRDT |
| **Go Backend** | `go-backend.mdc` | `backend/**` | gRPC 서비스, DB, 인증, 이벤트 |
| **Frontend** | `frontend.mdc` | `frontend/**` | Flutter UI, Next.js 웹, 디자인 시스템 |
| **Security** | `security-compliance.mdc` | 전체 | HIPAA/GDPR/MFDS/FDA, OWASP, STRIDE |
| **Infra** | `infrastructure.mdc` | `infrastructure/**` | Docker, K8s, CI/CD, 모니터링 |

### 6.2 작업 유형별 에이전트 라우팅
```
[차동측정 공식/핑거프린트/BLE/NFC] → rust-core-agent
[gRPC API/DB CRUD/인증/구독]       → go-backend-agent
[Flutter UI/디자인/라우팅/상태관리]  → frontend-agent
[보안 검토/규정 분석/위협 모델링]    → security-agent
[Docker/K8s/CI/CD/모니터링]         → infra-agent
[풀스택 기능/통합]                   → orchestrator → 병렬 위임
```

### 6.3 협업 파이프라인

#### 새 기능 (Full-Stack)
```
security-agent  →  위협 모델 & 보안 요구사항
orchestrator    →  작업 분할
[parallel]
  rust-core-agent → 코어 로직
  go-backend-agent → API
  frontend-agent → UI
security-agent  →  코드 보안 리뷰
infra-agent     →  배포
```

#### 측정 기능 (엔드투엔드)
```
rust-core-agent → BLE GATT 통신 + 차동측정 + 핑거프린트 생성
go-backend-agent → gRPC StreamMeasurement + TimescaleDB + Milvus 저장
frontend-agent → Flutter 측정 화면 + 결과 표시
```

---

## 7. 프로젝트 현황 (실시간)

### 7.1 진행률
| 영역 | 상태 | 진행률 | 담당 |
|------|------|--------|------|
| 프로젝트 구조 | ✅ 완료 | 100% | Antigravity |
| Rust 코어 엔진 | ✅ 완료 (8모듈, 스텁 3개) | 90% | Antigravity |
| Docker 인프라 | ✅ 완료 (15+ 서비스) | 100% | Antigravity |
| gRPC Proto 정의 | ✅ 완료 | 100% | Antigravity |
| Flutter 앱 | 🔲 대기 | 0% | ChatGPT |
| Go 백엔드 서비스 | 🔲 구조만 생성 | 5% | ChatGPT |
| 규정/보안 분석 | 🔲 스펙 전달 대기 | 0% | Claude |
| 통합 테스트 | 🔲 대기 | 0% | Antigravity |
| AI/ML 파이프라인 | 🔲 대기 | 0% | TBD |
| 에이전트 팀 설정 | ✅ 완료 | 100% | Claude |

### 7.2 Rust 모듈 구현 상태
| 모듈 | 상태 | 코드량 | 테스트 |
|------|------|--------|--------|
| differential | ✅ 구현 완료 | 213줄 | ✅ 3개 |
| fingerprint | ✅ 구현 완료 | 273줄 | ✅ 3개 |
| ble | ✅ 구현 완료 | 396줄 | ✅ 3개 |
| nfc | ✅ 구현 완료 | 478줄 | ✅ 4개 |
| ai | ✅ 구현 완료 | 368줄 | ✅ 4개 |
| crypto | 🔸 스텁 | 34줄 | ❌ |
| dsp | 🔸 스텁 | 36줄 | ❌ |
| sync | 🔸 스텁 | 82줄 | ❌ |
| flutter-bridge | ✅ 구현 완료 | 261줄 | ✅ 3개 |
| lib.rs (메인) | ✅ 구현 완료 | 147줄 | ✅ 2개 |

### 7.3 현재 Phase
```
Phase 1: MVP (Month 1-4)
  ✅ Week 1-4: 기반 인프라, Rust 코어
  ✅ Week 5-8: BLE/NFC/AI 모듈
  🔲 Week 9-12: Flutter 앱 기본 (대기)
  🔲 Week 13-16: Go 서비스 구현 (대기)
```

---

## 8. 핵심 결정 사항 (변경 불가)

1. **차동측정 공식**: `S_det - α × S_ref` (α 기본값 = 0.95)
2. **핑거프린트 차원**: 88 → 448 → 896 확장 경로
3. **리더기 관리**: 무제한 확장 (티어별 기본값만 다름)
4. **오프라인**: 100% 완전 동작 (CRDT 동기화)
5. **데이터 패킷**: 패밀리C 표준 준수
6. **AI 모델**: TFLite 엣지 + 클라우드 하이브리드
7. **보안**: AES-256-GCM + SHA-256 해시체인 + E2E 암호화
8. **인증**: Keycloak OIDC + JWT (Access + Refresh)
9. **API Gateway**: Kong 3.7

---

## 9. 크로스커팅 규칙 (모든 에이전트 필수 준수)

### 보안 (Security First - OWASP Top 10)
- 모든 입력값: 엄격한 검증 (Rust: validator, Go: binding 태그, Dart: form validation)
- SQL: ORM만 사용 (직접 쿼리 절대 금지)
- 인증: JWT + RBAC 필수
- 민감 정보: 코드 하드코딩 절대 금지
- 의료 데이터(PHI): AES-256-GCM 암호화 + 접근 감사 로그

### 품질 (TDD)
- "No Code Without Tests" - 실패하는 테스트 먼저
- 커버리지 80% 이상
- Rust: `#[cfg(test)]` + criterion 벤치마크
- Go: `_test.go` + testcontainers
- Dart: widget test + integration test

### 의료기기 규정
- IEC 62304 소프트웨어 수명주기 프로세스
- ISO 14971 위험관리
- 해시체인 무결성 검증 (SHA-256)
- 감사 추적 기록 필수 (최소 10년 보존)

### 코드 스타일
- Rust: `clippy::all`, `rust_2018_idioms`
- Go: `golangci-lint`, effective Go
- Dart: `flutter_lints` 3.x
- TypeScript: strict mode + ESLint

### 응답 언어
- 모든 응답, 커밋 메시지, 문서: **한국어**
- 코드 내 변수명/함수명: **영어**

---

## 10. 참조 문서 맵

| 문서 | 경로 | 용도 |
|------|------|------|
| 프로젝트 README | `README.md` | 개요, 기술 스택 |
| AI 협업 분담 | `docs/AI_COLLABORATION.md` | 3-AI 역할 분담 |
| 작업 로그 | `CHANGELOG.md` | 실시간 공유 작업 기록 |
| 공유 컨텍스트 | `CONTEXT.md` | AI 간 빠른 컨텍스트 공유 |
| 의료 규정 스펙 | `docs/ai-specs/claude/medical-compliance-spec.md` | HIPAA/GDPR/MFDS/FDA |
| 보안 아키텍처 스펙 | `docs/ai-specs/claude/security-architecture-spec.md` | 보안 설계 검토 |
| Flutter UI 스펙 | `docs/ai-specs/chatgpt/flutter-ui-spec.md` | Flutter 구현 가이드 |
| Go 서비스 스펙 | `docs/ai-specs/chatgpt/go-services-spec.md` | Go MSA 구현 가이드 |
| gRPC Proto | `backend/shared/proto/manpasik.proto` | API 정의 |
| 헬스체크 Proto | `backend/shared/proto/health.proto` | 헬스체크 프로토콜 |
| Docker Compose | `infrastructure/docker/docker-compose.dev.yml` | 개발 환경 |

---

## 11. 실시간 작업 기록 프로토콜 (필수)

> **모든 에이전트는 코드 변경이 수반되는 작업 완료 시 반드시 이 프로토콜을 따릅니다.**
> 만파식 프로젝트는 3개 AI(Antigravity, Claude, ChatGPT)가 병렬 작업하므로,
> 실시간 기록 없이는 작업 충돌과 컨텍스트 손실이 발생합니다.

### 세션 시작 시 (필수)
```
1. CONTEXT.md 읽기 → 현재 프로젝트 상태 파악
2. CHANGELOG.md 최근 항목 읽기 → 다른 AI의 최근 작업 확인
3. 충돌 가능성 확인 → 같은 파일 동시 수정 방지
```

### 작업 완료 시 (필수)
```
1. CHANGELOG.md 업데이트 → "## 🔄 최근 작업 로그" 바로 아래에 새 항목 추가
   - 형식: ## [날짜] [AI명] - [작업 제목]
   - 포함: 상태, 변경 사항(모든 파일), 결정 사항(이유), 다음 단계(담당 AI)
   
2. CONTEXT.md 업데이트 → 변경된 상태 반영
   - 진행률, 기술 스택, 결정 사항, AI 역할 등

3. CHANGELOG.md 현황 테이블 갱신 → 프로젝트 현황 요약 진행률 업데이트
```

### 충돌 방지 규칙
- CHANGELOG.md: **최상단에만 추가** (기존 항목 수정 금지)
- CONTEXT.md: **자신의 작업 관련 섹션만** 업데이트
- 다른 AI의 기록을 삭제/수정하지 않는다

### 공유 문서 체계
| 문서 | 용도 | 주기 |
|------|------|------|
| `CHANGELOG.md` | 시간순 작업 로그 (최신 상단) | 매 작업 완료 |
| `CONTEXT.md` | 프로젝트 현재 상태 스냅샷 | 상태 변경 시 |
| `AGENTS.md` | 마스터 컨텍스트 (아키텍처/API) | 아키텍처 변경 시 |
| `.cursor/rules/work-logging.mdc` | 이 프로토콜 상세 규칙 | 필요 시 |

---

**Document Version**: 2.1.0
**Last Updated**: 2026-02-09
**Author**: Claude (Security & Architecture Agent)
