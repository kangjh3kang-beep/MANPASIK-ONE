# ManPaSik 비기능 요구사항 명세서 (Non-Functional Requirements)

**문서번호**: MPK-NFR-v1.0  
**갱신일**: 2026-02-12  
**목적**: 성능, 가용성, 확장성, 보안, 스토리지 등 비기능 요구사항의 정량적 목표 정의  
**적용 범위**: 전체 시스템 (Rust 코어, Go 백엔드, Flutter 앱, 인프라)

---

## 1. 성능 요구사항 (Performance)

### 1.1 API 응답시간

| 카테고리 | P50 | P95 | P99 | 측정 기준 |
|---------|-----|-----|-----|----------|
| 인증 API (Login/Register) | < 100ms | < 300ms | < 500ms | 서버 사이드, Kong 이후 |
| 조회 API (GetProfile, ListDevices) | < 50ms | < 150ms | < 300ms | 서버 사이드 |
| 측정 API (StartSession/EndSession) | < 100ms | < 200ms | < 500ms | 서버 사이드 |
| 스트리밍 API (StreamMeasurement) | < 30ms/메시지 | < 50ms | < 100ms | gRPC 단방향 지연 |
| 벡터 검색 (Milvus 유사도) | < 50ms | < 200ms | < 500ms | 896차원, 100만 벡터 기준 |
| AI 추론 (엣지, TFLite) | < 100ms | < 200ms | < 300ms | 모바일 디바이스 |
| AI 추론 (서버, 시뮬레이션) | < 200ms | < 500ms | < 1s | ai-inference-service |
| 파일 업로드 (MinIO) | < 500ms | < 1s | < 2s | 5MB 이미지 기준 |
| 검색 (Elasticsearch) | < 100ms | < 300ms | < 500ms | 단순 텍스트 검색 |

### 1.2 처리량 (Throughput)

| Phase | 동시 접속 (CCU) | 초당 요청 (RPS) | 초당 측정 | 비고 |
|-------|---------------|----------------|---------|------|
| Phase 1 (MVP) | 100 | 500 | 50 | 개발/테스트 |
| Phase 2 (Core) | 1,000 | 5,000 | 500 | 베타 출시 |
| Phase 3 (Advanced) | 10,000 | 50,000 | 2,000 | 정식 출시 (한국) |
| Phase 4 (Ecosystem) | 100,000 | 200,000 | 10,000 | 글로벌 확장 |
| Phase 5 (Future) | 500,000 | 1,000,000 | 50,000 | 최대 부하 |

### 1.3 Rust 코어 성능

| 연산 | 목표 | 측정 기준 |
|------|------|----------|
| 차동측정 (88ch 단일) | < 1μs | criterion 벤치마크 |
| 차동측정 (896ch 배치) | < 10μs | criterion 벤치마크 |
| 핑거프린트 생성 (896차원) | < 100μs | FingerprintBuilder::build() |
| 코사인 유사도 (896차원 쌍) | < 5μs | — |
| AES-256-GCM 암호화 (1KB) | < 50μs | ring 기반 |
| SHA-256 해시 (1KB) | < 10μs | — |
| FFT (1024-point) | < 100μs | rustfft 기반 |
| BLE 패킷 파싱 | < 10μs | MeasurementPacket::parse() |
| NFC 태그 파싱 (v2.0) | < 50μs | CartridgeTag::parse() |

---

## 2. 가용성 요구사항 (Availability)

### 2.1 서비스 가용성 SLA

| Phase | 가용성 목표 | 허용 다운타임/월 | 비고 |
|-------|-----------|----------------|------|
| Phase 1 | 99.0% | 7.3시간 | 개발 환경, 계획 중단 허용 |
| Phase 2 | 99.5% | 3.7시간 | 베타, 야간 유지보수 허용 |
| Phase 3 | 99.9% | 43.2분 | 프로덕션, 무중단 배포 |
| Phase 4 | 99.95% | 21.6분 | 글로벌, 다중 리전 |
| Phase 5 | 99.99% | 4.3분 | 의료기기, 자동 복구 |

### 2.2 복구 목표

| 항목 | 목표 (Phase 3+) | 비고 |
|------|----------------|------|
| **RTO** (Recovery Time Objective) | < 5분 | 서비스 장애 복구 시간 |
| **RPO** (Recovery Point Objective) | < 1분 | 데이터 유실 허용 시점 |
| **MTTR** (Mean Time To Recovery) | < 15분 | 평균 복구 시간 |
| **MTBF** (Mean Time Between Failures) | > 720시간 (30일) | 평균 장애 간격 |

### 2.3 오프라인 가용성 (Rust 코어)

| 항목 | 목표 | 검증 방법 |
|------|------|----------|
| 오프라인 측정 | 100% 기능 동작 | REQ-065 72시간 연속 테스트 |
| 오프라인 AI 추론 | 100% (TFLite 엣지) | 5종 모델 전부 로컬 |
| 오프라인 데이터 보존 | 최소 72시간 분량 | CRDT 로컬 저장소 |
| 오프라인→온라인 동기화 | < 30초 (100건) | CRDT 병합 + gRPC Stream |

---

## 3. 확장성 요구사항 (Scalability)

### 3.1 수평 확장

| 컴포넌트 | 최소 인스턴스 | 최대 인스턴스 | 자동확장 조건 |
|---------|------------|------------|-------------|
| Go gRPC 서비스 (각) | 2 | 20 | CPU > 70% 또는 RPS > 1,000/인스턴스 |
| Kong API Gateway | 2 | 10 | 요청 큐 > 100 |
| Keycloak | 2 | 5 | 활성 세션 > 10,000 |
| PostgreSQL | 1 Primary + 2 Read Replica | 1P + 5RR | 읽기 부하 > 5,000 QPS |
| Redis | 1 Master + 2 Replica | 3M + 6R (Cluster) | 메모리 > 80% |
| Milvus | 1 Standalone | 3 Cluster | 벡터 수 > 1,000만 |
| Kafka (Redpanda) | 3 Broker | 9 Broker | 파티션 수, 처리량 |
| Elasticsearch | 3 Node | 9 Node | 인덱스 크기 > 100GB |

### 3.2 데이터 증가 예측

| 데이터 유형 | 단건 크기 | Phase 2 (1년차) | Phase 4 (3년차) | 5년 총량 |
|-----------|---------|-----------------|-----------------|---------|
| 측정 레코드 (TimescaleDB) | ~2KB | 50만 건/일 → 180GB/년 | 500만 건/일 → 1.8TB/년 | ~5TB |
| 핑거프린트 벡터 (Milvus) | 3.5KB (896×f32) | 50만 벡터/일 → 630GB/년 | 500만/일 → 6.3TB/년 | ~18TB |
| 사용자 데이터 (PostgreSQL) | ~5KB/유저 | 10만 유저 → 500MB | 100만 유저 → 5GB | ~10GB |
| 이벤트 로그 (Kafka→ES) | ~500B | 200만 이벤트/일 → 360GB/년 | 2,000만/일 → 3.6TB/년 | ~10TB |
| 파일 (MinIO) | ~2MB 평균 | 1만 파일/일 → 7.3TB/년 | 10만/일 → 73TB/년 | ~200TB |

### 3.3 데이터 보존 정책

| 데이터 유형 | 보존 기간 | 근거 | 보존 후 처리 |
|-----------|---------|------|-------------|
| 측정 원시 데이터 | 10년 | IEC 62304, 의료기기 규정 | 콜드 스토리지 이전 (3년 후) |
| 감사 추적 로그 | 10년 | ISO 13485, FDA 21 CFR Part 11 | 압축 아카이브 |
| 개인정보 (PHI) | 회원 탈퇴 후 90일 | GDPR Art.17, 개인정보보호법 | 완전 삭제 (물리적) |
| 이벤트 로그 | 1년 | 운영 분석용 | 삭제 또는 집계 보존 |
| 벡터 데이터 | 5년 | AI 모델 품질 유지 | 다운샘플링 후 아카이브 |
| 백업 | 90일 (일간), 1년 (주간) | DR 대비 | 순환 삭제 |

---

## 4. 보안 요구사항 (Security)

### 4.1 정량적 보안 목표

| 항목 | 목표 | 측정 방법 |
|------|------|----------|
| 취약점 해결 시간 (Critical) | < 24시간 | JIRA SLA |
| 취약점 해결 시간 (High) | < 7일 | JIRA SLA |
| 보안 스캔 빈도 | 매 PR + 주 1회 전체 | CI/CD + Cron |
| 비밀번호 최소 복잡도 | 8자 이상, 대/소/숫/특 포함 | Validator |
| JWT Access Token TTL | 15분 | Keycloak 설정 |
| JWT Refresh Token TTL | 7일 | Redis TTL |
| 세션 타임아웃 (비활성) | 30분 | 미들웨어 |
| Rate Limiting (인증) | 5회/분/IP | middleware/rate_limit.go |
| Rate Limiting (일반 API) | 100회/분/유저 | Kong + 미들웨어 |
| 암호화 강도 | AES-256-GCM, RSA-2048+ | crypto 모듈 |
| TLS 버전 | TLS 1.3 필수 | Kong/Ingress 설정 |
| OWASP Top 10 준수 | 10/10 항목 대응 | 보안 감사 |

### 4.2 PHI(개인건강정보) 처리

| 항목 | 요구사항 | 구현 위치 |
|------|---------|----------|
| 저장 시 암호화 | AES-256-GCM | Rust crypto + DB TDE |
| 전송 시 암호화 | TLS 1.3 (E2E) | Kong + gRPC mTLS |
| 접근 로그 | 전 PHI 접근 기록 (최소 10년) | shared/observability |
| 최소 권한 원칙 | RBAC + 데이터 필터링 | middleware/rbac.go |
| 동의 기반 접근 | 데이터 공유 동의서 필수 | 23-data-sharing-consents.sql |
| 익명화/가명화 | 연구용 데이터 k-anonymity ≥ 5 | TODO (Phase 3) |

---

## 5. 안정성 요구사항 (Reliability)

### 5.1 오류율 목표

| 항목 | 목표 | 측정 |
|------|------|------|
| API 오류율 (5xx) | < 0.1% | Prometheus error_rate |
| gRPC 오류율 | < 0.05% | gRPC interceptor |
| 메시지 유실률 (Kafka) | 0% (at-least-once) | Consumer lag 모니터링 |
| 측정 데이터 무결성 | 100% (SHA-256 해시체인) | 해시 검증 |
| DB 트랜잭션 실패율 | < 0.01% | PostgreSQL metrics |

### 5.2 장애 복구 시나리오

| 시나리오 | RTO | RPO | 자동화 |
|---------|-----|-----|--------|
| 단일 서비스 크래시 | < 30초 | 0 | K8s Pod 재시작 |
| DB Primary 장애 | < 60초 | < 10초 | Patroni 자동 failover |
| Redis Master 장애 | < 30초 | < 5초 | Sentinel 자동 failover |
| Kafka Broker 장애 | < 60초 | 0 | Replication factor ≥ 3 |
| 전체 리전 장애 | < 5분 | < 1분 | Multi-region failover (Phase 4+) |
| 네트워크 파티션 | — | 0 | 오프라인 모드 자동 전환 (CRDT) |

---

## 6. 사용성 요구사항 (Usability)

| 항목 | 목표 | 측정 방법 |
|------|------|----------|
| 앱 시작 시간 (Cold) | < 3초 | 모바일 프로파일링 |
| 화면 전환 시간 | < 300ms | Flutter DevTools |
| 측정 시작까지 탭 수 | ≤ 3회 | UX 테스트 |
| 다국어 지원 | 6개 언어 (ko/en/ja/zh/fr/hi) | i18n 커버리지 |
| 접근성 (WCAG) | AA 등급 | axe-core 스캔 |
| 다크모드 지원 | 100% 화면 | Material 3 theme |
| 에러 메시지 이해도 | 비전문가 이해 가능 | UX 사용자 테스트 |

---

## 7. 운영 요구사항 (Operability)

| 항목 | 목표 | Phase |
|------|------|-------|
| 로그 중앙화 | 100% 서비스 → ELK | Phase 2 |
| 분산 트레이싱 | 100% gRPC → OpenTelemetry → Jaeger | Phase 2 |
| 메트릭 수집 | 100% 서비스 → Prometheus | Phase 1 (완료) |
| 알림 (Alert) | P1: 5분 내 통보, P0: 1분 내 | Phase 3 |
| 대시보드 | Grafana — 서비스별 + 비즈니스 | Phase 2 |
| 배포 빈도 | ≥ 1회/주 (Phase 3+) | CI/CD |
| 롤백 시간 | < 5분 | K8s 롤백 |
| 배포 성공률 | > 99% | CI/CD 메트릭 |

---

## 8. 검증 계획

| NFR 카테고리 | 검증 도구 | 검증 시점 | 담당 |
|-------------|---------|---------|------|
| API 응답시간 | k6, wrk, ghz (gRPC) | Phase 2 Gate, 매 릴리스 | QA |
| 처리량 | k6 부하 테스트 | Phase Gate | QA |
| Rust 성능 | criterion 벤치마크 | 매 PR | CI |
| 가용성 | Chaos Monkey, Litmus | Phase 3 Gate | Infra |
| 보안 | Trivy, gosec, cargo-audit | 매 PR + 주 1회 | CI + Security |
| 스토리지 | Prometheus disk metrics | 월 1회 검토 | Infra |
| 사용성 | Flutter DevTools, Lighthouse | Phase Gate | Frontend |

---

**참조**: QUALITY_GATES.md, docs/compliance/vnv-master-plan.md, docs/security/stride-threat-model.md
