# V&V 마스터 플랜 (Verification & Validation Master Plan)

## ManPaSik (만파식) 차동측정 기반 헬스케어 AI 생태계

**문서 번호**: MPS-VMP-001  
**버전**: 1.0  
**분류**: 기밀 (Confidential)  
**작성일**: 2026-02-10  
**작성자**: Claude (Security & Compliance Agent)  
**근거 표준**: IEC 62304:2015, ISO 13485:2016, FDA Software Validation Guidance  

---

## 1. 목적

본 문서는 ManPaSik 소프트웨어 시스템의 검증(Verification)과 확인(Validation) 전략을 정의합니다.

- **Verification (검증)**: "올바르게 만들었는가?" — 설계 명세 대비 구현 확인
- **Validation (확인)**: "올바른 것을 만들었는가?" — 사용자 요구사항 충족 확인

---

## 2. 적용 범위

| 구성요소 | IEC 62304 등급 | V&V 수준 |
|----------|-------------|---------|
| Rust 코어 엔진 (10모듈: ble, nfc, dsp, ai, crypto, sync, storage, ffi, config, error) | Class B | 전체 V&V |
| Go 마이크로서비스 (22서비스: auth, user, measurement, device, cartridge 외 17개) | Class B | 전체 V&V |
| Flutter 모바일 앱 (12개 Feature: auth, home, measurement, devices, chat, settings, admin, ai_coach, data_hub, community, family, market, medical) | Class B | 전체 V&V |
| Next.js 웹 대시보드 | Class A | 기본 V&V |
| AI/ML 모델 (5종: 바이오마커 분류/회귀, 이상탐지, 건강코치, 카트리지 QC) | Class B | 강화 V&V (FDA AI/ML) |
| 인프라 (Docker Compose, K8s, PostgreSQL 21 스키마) | Class A | 기본 V&V |

---

## 3. V&V 전략

### 3.1 검증 (Verification) 전략

```
Level 1: 단위 검증 (Unit Verification)
  ├── 코드 리뷰 (모든 PR 필수)
  ├── 정적 분석 (Rust clippy, Go golangci-lint, Dart flutter_lints)
  ├── 단위 테스트 (커버리지 ≥ 80%)
  └── 속성 기반 테스트 (Rust proptest — 안전 관련 모듈)

Level 2: 통합 검증 (Integration Verification)
  ├── 모듈 간 인터페이스 테스트
  ├── API 계약 테스트 (gRPC Proto 준수)
  ├── 데이터베이스 통합 테스트 (Testcontainers)
  └── BLE/NFC 프로토콜 테스트

Level 3: 시스템 검증 (System Verification)
  ├── 엔드투엔드 테스트 (측정 파이프라인 전체)
  ├── 성능 테스트 (부하, 스트레스, 지구력)
  ├── 보안 테스트 (침투 테스트, OWASP 검증)
  └── 호환성 테스트 (OS, 디바이스, 브라우저)
```

### 3.2 확인 (Validation) 전략

```
Level 4: 사용적합성 확인 (Usability Validation)
  ├── 사용적합성 공학 프로세스 (IEC 62366-1)
  ├── 사용자 인터뷰 및 관찰
  ├── 형성적 평가 (Formative Evaluation)
  └── 종합적 평가 (Summative Evaluation)

Level 5: 임상 확인 (Clinical Validation)
  ├── 분석 성능 평가 (정확도, 정밀도, LOD/LOQ)
  ├── 임상 성능 평가 (민감도, 특이도, 양/음성 예측도)
  ├── 참고치 상관성 연구 (방법 비교)
  └── 안정성 연구 (카트리지 유효기간)
```

---

## 4. 테스트 유형별 상세

### 4.1 단위 테스트

| 구성요소 | 프레임워크 | 커버리지 목표 | 명명 규칙 |
|----------|----------|-----------|----------|
| Rust 코어 | `#[cfg(test)]` + criterion | ≥ 85% | `test_기능_시나리오_기대결과` |
| Go 백엔드 | `testing` + testify | ≥ 80% | `Test기능_시나리오` |
| Flutter | flutter_test + mocktail | ≥ 80% | `test기능 시나리오` |
| TypeScript | Jest + RTL | ≥ 75% | `it('should 동작')` |

### 4.2 통합 테스트

| 테스트 | 대상 | 도구 | 환경 |
|--------|------|------|------|
| Rust ↔ Flutter | FFI 브리지 10개 API | flutter_rust_bridge 테스트 | 로컬 |
| Go ↔ PostgreSQL | CRUD 전 테이블 | testcontainers-go | Docker |
| Go ↔ TimescaleDB | 시계열 쿼리 | testcontainers-go | Docker |
| Go ↔ Redis | 캐시/세션 | testcontainers-go | Docker |
| Go ↔ Milvus | 벡터 검색 | testcontainers-go | Docker |
| gRPC 서비스 간 | 4개 서비스 호출 | grpc-testing | Docker Compose |
| BLE 통신 | 리더기 ↔ 앱 | 시뮬레이터 | 로컬 |

### 4.3 시스템 테스트

| 테스트 유형 | 시나리오 | 통과 기준 |
|-----------|---------|----------|
| 측정 파이프라인 E2E | BLE 연결 → 측정 → 결과 표시 | 정상 완료, 데이터 무결성 |
| 오프라인 동기화 | 오프라인 측정 → 재연결 → 동기화 | 데이터 손실 0건, CRDT 병합 정상 |
| 다중 리더기 | 3개 리더기 동시 연결 | 독립 측정 정상 |
| 장시간 연속 측정 | 8시간 연속 사용 | 메모리 누수 없음, 성능 저하 없음 |
| 대량 동시 접속 | 1,000 동시 사용자 | 응답 시간 < 500ms (p99) |
| 페일오버 | DB/서비스 장애 복구 | RTO < 5분, 데이터 손실 0 |

### 4.4 보안 테스트

| 테스트 | 도구 | 기준 |
|--------|------|------|
| SAST (정적 분석) | cargo-audit, gosec, semgrep | 심각(Critical) 0건 |
| DAST (동적 분석) | OWASP ZAP | 높음(High) 0건 |
| 침투 테스트 | 전문 업체 위탁 | OWASP Top 10 전 항목 |
| 의존성 취약점 | cargo-deny, go-vulncheck | CVE 심각 0건 |
| 암호화 검증 | 자체 + 외부 감사 | FIPS 140-2 Level 1 |

### 4.5 AI/ML 모델 검증 (FDA AI/ML SaMD 가이드라인)

| 검증 항목 | 방법 | 기준 |
|----------|------|------|
| 학습 데이터 품질 | 데이터 감사, 라벨 품질 평가 | 라벨 오류율 < 1% |
| 모델 정확도 | 5-Fold 교차 검증 | AUC > 0.90 (분류), R² > 0.95 (회귀) |
| 모델 편향 | 인구통계별 성능 비교 | 그룹 간 성능 차이 < 5% |
| 로버스트니스 | 적대적 입력, 노이즈 주입 | 정확도 감소 < 10% |
| 설명가능성 | SHAP/LIME | 주요 특성 해석 가능 |
| 모델 드리프트 | 시간별 성능 모니터링 | 성능 저하 시 자동 알림 |

---

## 5. 테스트 환경

### 5.1 환경 구성

| 환경 | 용도 | 인프라 |
|------|------|--------|
| **개발 (Dev)** | 단위/통합 테스트 | Docker Compose (15+ 서비스) |
| **CI** | 자동화 테스트 | GitHub Actions / GitLab CI |
| **스테이징 (Staging)** | 시스템/성능 테스트 | K8s 클러스터 (프로덕션 미러) |
| **프로덕션 (Prod)** | 임상 확인 | K8s 클러스터 (GMP 준수) |

### 5.2 테스트 데이터 관리

- **개인정보**: 합성(Synthetic) 데이터만 사용 (개발/CI/스테이징)
- **임상 데이터**: IRB 승인 후 프로덕션 환경에서만 처리
- **데이터 분리**: 각 환경 간 데이터 격리 필수
- **시드 데이터**: 각 카트리지 타입별 대표 측정 데이터셋

---

## 6. CI/CD 파이프라인 통합

```
[코드 커밋] → [Lint/Format] → [단위 테스트] → [정적 분석]
                                                   ↓
                                         [통합 테스트] → [보안 스캔]
                                                            ↓
                                              [Docker 빌드] → [E2E 테스트]
                                                                   ↓
                                              [스테이징 배포] → [성능 테스트]
                                                                   ↓
                                                         [수동 승인] → [프로덕션]
```

### 게이트 기준

| 게이트 | 조건 | 실패 시 |
|--------|------|---------|
| PR 병합 | 모든 테스트 통과 + 코드 리뷰 승인 | 병합 차단 |
| 스테이징 배포 | 통합 + 보안 테스트 통과 | 배포 차단 |
| 프로덕션 배포 | E2E + 성능 + 보안 전 통과 + 수동 승인 | 배포 차단 |

---

## 7. 추적성 매트릭스 (Traceability)

```
사용자 요구사항 (URS)
  ↕ 양방향 추적
소프트웨어 요구사항 (SRS)
  ↕ 양방향 추적
소프트웨어 설계 (SDS)
  ↕ 양방향 추적
소프트웨어 구현 (소스 코드)
  ↕ 양방향 추적
테스트 케이스 (V&V)
  ↕ 양방향 추적
위험 통제 조치 (ISO 14971)
```

- 모든 요구사항은 최소 하나의 테스트 케이스로 검증
- 모든 위험 통제 조치는 테스트로 효과 검증
- 변경 시 영향받는 테스트 자동 식별 (CI 트리거)

---

## 8. 결함 관리

### 8.1 결함 심각도 분류

| 등급 | 설명 | 수정 기한 |
|------|------|----------|
| **Blocker** | 시스템 사용 불가 | 즉시 |
| **Critical** | 핵심 기능 실패, 데이터 손실 위험 | 24시간 |
| **Major** | 주요 기능 장애, 우회 방법 존재 | 1주 |
| **Minor** | UI/UX 문제, 성능 저하 | 다음 릴리스 |
| **Trivial** | 오타, 미용적 이슈 | 백로그 |

### 8.2 안전 관련 결함

- 위험 통제 관련 결함은 **자동으로 Critical** 이상 분류
- 별도의 CAPA (시정/예방 조치) 프로세스 적용
- 잔여 위험 재평가 필수

---

## 9. V&V 일정

| Phase | 활동 | 기간 | 산출물 |
|-------|------|------|--------|
| Phase 1 (MVP) | 단위 테스트 프레임워크 구축, CI 파이프라인 | Month 1-4 | 테스트 인프라 |
| Phase 2 (Core) | 통합 테스트, 보안 테스트 | Month 5-8 | 통합 테스트 보고서 |
| Phase 3 (Advanced) | 시스템 테스트, 성능 테스트 | Month 9-12 | 시스템 테스트 보고서 |
| Phase 4 (Validation) | 사용적합성 평가, 임상 확인 | Month 13-16 | V&V 최종 보고서 |

---

## 10. 구현된 테스트 산출물 현황

### 10.1 E2E 테스트 (`tests/e2e/`)

| 파일 | 시나리오 수 | 커버리지 |
|------|----------|---------|
| auth_flow_test.go | 5 | 회원가입, 로그인, 중복, 토큰갱신, 로그아웃 |
| measurement_flow_test.go | 5 | 전체 플로우, 차동측정 정확도, 멀티카트리지, 히스토리, 미인증 차단 |
| device_management_test.go | 4 | 디바이스 CRUD, 멀티디바이스, 펌웨어, 접근제어 |
| medical_service_test.go | 5 | 비대면진료, 처방전, 건강리포트, 응급알림, 프라이버시 |
| community_family_test.go | 5 | 게시글 CRUD, 좋아요/댓글, 가족그룹, 데이터공유, 권한 |
| admin_test.go | 5 | 사용자관리, 카트리지관리, 시스템설정, 감사로그, 번역관리 |
| **합계** | **29** | 6개 도메인 전체 커버 |

### 10.2 보안 테스트 (`tests/security/`)

| 파일 | 테스트 수 | OWASP 커버리지 |
|------|----------|--------------|
| auth_security_test.go | 8 | A01 (권한), A04 (Rate Limit), A07 (인증) |
| api_security_test.go | 9 | A03 (인젝션), A05 (설정), A06 (취약점), A08 (무결성), A09 (로깅) |
| data_security_test.go | 11 | A02 (암호화), A10 (SSRF), GDPR/PIPA/HIPAA |
| dependency_scan.sh | - | Go/Dart/Rust/Docker 의존성 스캔 |
| **합계** | **28** | **OWASP Top 10 전 항목 (10/10)** |

### 10.3 부하 테스트 (`tests/load/`, k6 기반)

| 파일 | 시나리오 | 목표 지표 |
|------|---------|----------|
| auth_load_test.js | 인증 부하 | 500 VU, p95 < 200ms |
| measurement_load_test.js | 측정 부하 | 200 VU, p95 < 500ms |
| stress_test.js | 스트레스 | 2000 VU 단계적 증가 |
| spike_test.js | 스파이크 | 0→1000 VU 급증 |
| concurrent_users_test.js | 동시접속 | 1000 동시 사용자 |
| api_gateway_load_test.js | 게이트웨이 | 혼합 트래픽 500 VU |

### 10.4 위험 통제 검증 (ISO 14971 ↔ V&V 연계)

| FMEA ID | 위험 | 통제 조치 | 검증 테스트 |
|---------|------|---------|-----------|
| FM-AI-002 | AI 위음성 | SafetyValidator 이중검증 | ai/mod.rs 단위테스트 3개 |
| FM-BLE-001 | BLE 연결 끊김 | BleStateMachine + ReconnectionStrategy | ble/mod.rs 단위테스트 5개 |
| FM-BLE-003 | 데이터 손실 | ChunkReassembler + CRC32 | ble/mod.rs 단위테스트 1개 |
| FM-SYNC-001 | 오프라인 동기화 실패 | CRDT (GCounter, LWWRegister, ORSet) | sync 모듈 테스트 |
| FM-CRYPTO-001 | 암호화 실패 | AES-256-GCM + 키체인 | security 테스트 11개 |

---

## 11. 문서 이력

| 버전 | 날짜 | 변경 내용 | 작성자 |
|------|------|----------|--------|
| 1.0 | 2026-02-10 | 초안 작성 | Claude |
| 1.1 | 2026-02-14 | 서비스 수 정정 (4→22), 구현된 테스트 산출물 현황 추가, ISO 14971 연계 추가 | Claude |

---

*본 문서는 IEC 62304:2015 5.5절(소프트웨어 검증), 5.7절(소프트웨어 확인) 및*
*FDA "General Principles of Software Validation" 가이던스를 준수합니다.*
