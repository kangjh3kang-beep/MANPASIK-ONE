# IEC 62304:2015 소프트웨어 요구사항 명세서 (SRS)

> **문서 ID**: DOC-SRS-001  
> **버전**: v3.0  
> **안전 등급**: Class B (IEC 62304)  
> **작성일**: 2026-02-13  
> **상태**: 초안 (검토 대기)  
> **적용 표준**: IEC 62304:2006+AMD1:2015 §5.2 (소프트웨어 요구사항 분석)  
> **선행 문서**: DOC-SDP-001 v1.1 (소프트웨어 개발 계획서)

---

## 문서 제어 (IEC 62304 §5.2)

### 승인 이력

| 역할 | 이름 | 서명 | 일자 |
|------|------|------|------|
| 작성자 | — | — | 2026-02-13 |
| 기술 리더 검토 | — | — | — |
| 품질 관리자 검토 | — | — | — |
| 최종 승인 | — | — | — |

### 변경 이력

| 버전 | 일자 | 작성자 | 변경 내용 |
|------|------|--------|----------|
| v1.0 | 2026-02-12 | Claude Agent | 초안 작성 (핵심 요구사항 20개, NFR 14개) |
| v2.0 | 2026-02-13 | Claude Agent | 전면 개정 — 80개 REQ 전체 수록, 안전·규정 요구사항 확충 |
| v3.0 | 2026-02-13 | Claude Agent | REQ-FUNC/NFR/IF/DATA/SAFE/REG ID 체계 통일, 추적성 매트릭스 강화 |

---

## 목차

1. [개요](#1-개요-clause-521)
2. [참조 문서](#2-참조-문서)
3. [시스템 개요](#3-시스템-개요)
4. [기능 요구사항](#4-기능-요구사항-clause-522)
5. [비기능 요구사항](#5-비기능-요구사항-clause-523)
6. [인터페이스 요구사항](#6-인터페이스-요구사항-clause-524)
7. [데이터 요구사항 (REQ-DATA)](#7-데이터-요구사항-req-data-xxx)
8. [안전 요구사항](#8-안전-요구사항-clause-525)
9. [규정 요구사항](#9-규정-요구사항-clause-526)
10. [사용 환경](#10-사용-환경)
11. [추적성](#11-추적성-clause-527)
12. [부록](#부록)

---

## 1. 개요 (Clause 5.2.1)

### 1.1 목적

본 문서는 **ManPaSik (만파식, 萬波息)** 의료기기 소프트웨어의 기능적·비기능적·안전·규정 요구사항을 체계적으로 정의합니다. IEC 62304:2015 §5.2 요구사항을 충족하며, 모든 요구사항은 고유 ID를 부여받아 설계(SAD)·구현(코드)·검증(V&V)과 양방향 추적이 가능합니다.

### 1.2 범위

| 항목 | 내용 |
|------|------|
| **제품명** | ManPaSik (만파식) — 차동측정 기반 범용 분석 헬스케어 AI 생태계 |
| **의료기기 등급** | Class II (체외진단의료기기, IVD) |
| **소프트웨어 안전 등급** | IEC 62304 Class B |
| **대상 시장** | KR(MFDS), US(FDA 510(k)), EU(CE-IVDR), CN(NMPA), JP(PMDA) |
| **핵심 기술** | 차동측정(`S_det − α × S_ref`), 1792차원 핑거프린트, 엣지 AI(TFLite), 연합학습 |
| **기술 스택** | Go gRPC MSA 30+서비스, Rust 코어엔진, Flutter 모바일, PostgreSQL/TimescaleDB/Milvus/Redis |

### 1.3 용어 정의

| 용어 | 정의 |
|------|------|
| SRS | 소프트웨어 요구사항 명세서 (Software Requirements Specification) |
| REQ-FUNC | 기능 요구사항 식별자 접두어 |
| REQ-NFR | 비기능 요구사항 (Non-Functional Requirement) |
| REQ-SAFE | 안전 요구사항 (Safety Requirement) |
| REQ-REG | 규정 요구사항 (Regulatory Requirement) |
| REQ-IF | 인터페이스 요구사항 (Interface Requirement) |
| PHI | 보호 대상 건강 정보 (Protected Health Information) |
| PII | 개인 식별 정보 (Personally Identifiable Information) |
| SOUP | 알려진 출처의 소프트웨어 (Software of Unknown Provenance) |
| V&V | 검증 및 유효성 확인 (Verification and Validation) |
| CRDT | 충돌 없는 복제 데이터 타입 (Conflict-free Replicated Data Type) |
| RBAC | 역할 기반 접근 제어 (Role-Based Access Control) |

### 1.4 요구사항 ID 체계 (IEC 62304 §5.2)

| 접두어 | 유형 | 형식 예 | 용도 |
|--------|------|---------|------|
| **REQ-FUNC-xxx** | 기능 요구사항 | REQ-FUNC-AUTH-001, REQ-FUNC-MEAS-003 | 측정, AI, 인증, 구독, 카트리지 등 |
| **REQ-NFR-xxx** | 비기능 요구사항 | REQ-NFR-PERF-001, REQ-NFR-SEC-002 | 성능, 보안, 가용성, 확장성 |
| **REQ-IF-xxx** | 인터페이스 요구사항 | REQ-IF-001, REQ-IF-002 | gRPC, BLE, NFC, REST |
| **REQ-DATA-xxx** | 데이터 요구사항 | REQ-DATA-001, REQ-DATA-002 | PostgreSQL, TimescaleDB, Milvus, Redis |
| **REQ-SAFE-xxx** | 안전 요구사항 | REQ-SAFE-001, REQ-SAFE-002 | 측정 정확도, 데이터 무결성, 경고 |
| **REQ-REG-xxx** | 규정 요구사항 | REQ-REG-001, REQ-REG-002 | FDA 510(k), CE-IVDR, MFDS |

---

## 2. 참조 문서

### 2.1 외부 규정·표준

| 문서 | 버전 | 적용 |
|------|------|------|
| **IEC 62304** | 2006+AMD1:2015 | 의료기기 소프트웨어 — 소프트웨어 수명주기 프로세스 |
| **ISO 14971** | 2019 | 의료기기 — 위험관리 |
| **ISO 13485** | 2016 | 의료기기 — 품질경영시스템 |
| **IEC 62366-1** | 2015+Amd.1 | 의료기기 — 사용적합성 공학 |
| **FDA SW Guidance** | 2002 | General Principles of Software Validation |
| **FDA Cybersecurity** | 2023 | Cybersecurity in Medical Devices |
| **EU IVDR** | 2017/746 | 체외진단 의료기기 규정 |
| **GDPR** | 2016/679 | 일반 데이터 보호 규정 |
| **개인정보보호법 (PIPA)** | 2023 개정 | 한국 개인정보 보호 |
| **의료기기법** | 2024 개정 | 한국 의료기기 규제 |

### 2.2 프로젝트 내부 문서

| 문서 | 경로 | 관계 |
|------|------|------|
| SDP (개발 계획) | `docs/compliance/iec62304-sdp.md` | 본 문서의 상위 프로세스 |
| SAD (아키텍처 설계) | `docs/compliance/iec62304-sad.md` | 본 문서 REQ → SAD DES 추적 |
| 안전 등급 판정 | `docs/compliance/software-safety-classification.md` | Class B 판정 근거 |
| 위험관리 계획 | `docs/compliance/iso14971-risk-management-plan.md` | 위해 시나리오 → SAF 요구사항 도출 |
| STRIDE 위협 모델 | `docs/security/stride-threat-model.md` | 31개 위협 → 보안 통제 조치 |
| V&V 마스터 플랜 | `docs/compliance/vnv-master-plan.md` | REQ → 테스트 케이스 매핑 |
| 추적성 매트릭스 | `docs/plan/plan-traceability-matrix.md` | 80개 REQ ↔ DES ↔ IMP ↔ V&V |
| 비기능 요구사항 정량화 | `docs/specs/non-functional-requirements.md` | NFR 정량적 목표 상세 |
| 카트리지 시스템 명세 | `docs/specs/cartridge-system-spec.md` | 무한확장 체계 |
| 이벤트 스키마 명세 | `docs/specs/event-schema-specification.md` | 18개 Kafka 토픽 |
| Proto 정의 | `backend/shared/proto/manpasik.proto` | gRPC API 인터페이스 |

---

## 3. 시스템 개요

### 3.1 시스템 구성

ManPaSik 소프트웨어 시스템은 다음 4개 계층으로 구성됩니다:

```
┌──────────────────────────────────────────────────────────────┐
│ 1. 클라이언트 계층: Flutter 앱 (iOS/Android) + Next.js 웹    │
├──────────────────────────────────────────────────────────────┤
│ 2. 코어 계층: Rust 엔진 (차동측정, 핑거프린트, 엣지 AI)      │
├──────────────────────────────────────────────────────────────┤
│ 3. 서비스 계층: API Gateway (Kong) + Go gRPC MSA 30+서비스   │
├──────────────────────────────────────────────────────────────┤
│ 4. 데이터 계층: PostgreSQL + TimescaleDB + Milvus + Redis    │
│                 Kafka + Elasticsearch + MinIO                │
└──────────────────────────────────────────────────────────────┘
```

### 3.2 소프트웨어 항목 (Software Items) 요약

| 소프트웨어 항목 | 기술 스택 | 안전 등급 | 주요 기능 |
|---------------|----------|----------|----------|
| Flutter 모바일 앱 | Flutter 3.x + Dart | Class B | 사용자 UI, BLE/NFC 연동, 오프라인 측정 |
| Rust 코어 엔진 | Rust Ed.2021 | Class B | 차동측정, 핑거프린트, 엣지 AI, 암호화, CRDT |
| Go 백엔드 서비스 (30+) | Go 1.22+ gRPC | Class A~B | 인증, 측정, AI, 결제, 원격진료 등 |
| API Gateway | Kong 3.7 + Keycloak 25.0 | Class A | 인증, 라우팅, RBAC, SSL 종단 |
| 데이터 인프라 | PostgreSQL 16, Milvus 2.4, Redis 7 등 | — | 관계형·벡터·캐시·이벤트 저장 |

### 3.3 의도된 사용

| 항목 | 내용 |
|------|------|
| **사용 목적** | 차동측정 기반 건강 바이오마커 분석 및 건강 관리 정보 제공 |
| **사용자** | 일반 소비자 (환자가 아닌 건강 관리 목적), 의료 전문가 (원격진료) |
| **사용 환경** | 가정, 사무실, 의료기관 (실내 환경) |
| **사용 제한** | 진단 기기가 아님; 의료 전문가 확인 없이 자가 관리 참고 정보 제공 |

---

## 4. 기능 요구사항 (Clause 5.2.2)

### 4.1 인증 및 사용자 관리

| REQ ID | 요구사항 | 안전 등급 | 서비스 | 우선순위 |
|--------|---------|----------|--------|---------|
| REQ-FUNC-AUTH-001 | 시스템은 이메일/비밀번호 기반 사용자 회원가입을 지원해야 한다. 비밀번호는 bcrypt(cost≥12)로 해시 저장한다. | B | auth-service | P0 |
| REQ-FUNC-AUTH-002 | 시스템은 JWT 기반 인증을 제공해야 한다. Access Token 만료 15분, Refresh Token 만료 7일. | B | auth-service | P0 |
| REQ-FUNC-AUTH-003 | 시스템은 Refresh Token을 이용한 Access Token 자동 갱신을 지원해야 한다. | B | auth-service | P0 |
| REQ-FUNC-AUTH-004 | 시스템은 비밀번호 재설정 기능을 이메일 인증을 통해 제공해야 한다. | A | auth-service | P1 |
| REQ-FUNC-AUTH-005 | 시스템은 소셜 로그인(Google, Apple, Facebook)을 Keycloak OIDC를 통해 지원해야 한다. | A | auth-service | P2 |
| REQ-FUNC-AUTH-006 | 시스템은 TOTP 기반 다중 인증(MFA)을 제공해야 한다. | B | auth-service | P2 |
| REQ-FUNC-AUTH-007 | 시스템은 토큰 검증(ValidateToken) RPC를 통해 서비스 간 인증 확인을 지원해야 한다. | B | auth-service | P0 |
| REQ-FUNC-USER-001 | 시스템은 사용자 프로필(이름, 아바타, 언어, 시간대) 조회 및 수정을 지원해야 한다. | A | user-service | P0 |
| REQ-FUNC-USER-002 | 시스템은 사용자 구독 정보(티어, 만료일, 최대 디바이스 수) 조회를 지원해야 한다. | A | user-service | P0 |
| REQ-FUNC-USER-003 | 시스템은 사용자 삭제(탈퇴) 시 개인정보를 30일 후 완전 삭제해야 한다. | B | user-service | P1 |

### 4.2 디바이스 관리

| REQ ID | 요구사항 | 안전 등급 | 서비스 | 우선순위 |
|--------|---------|----------|--------|---------|
| REQ-FUNC-DEV-001 | 시스템은 BLE 5.0을 통한 만파식 리더기 검색 및 페어링을 지원해야 한다. | B | device-service, Rust Core | P0 |
| REQ-FUNC-DEV-002 | 시스템은 리더기 등록(시리얼 번호, 펌웨어 버전)을 지원해야 한다. 구독 티어별 최대 디바이스 수를 제한한다. | B | device-service | P0 |
| REQ-FUNC-DEV-003 | 시스템은 등록된 디바이스 목록 조회 및 상태(온라인/오프라인/측정중/업데이트중/에러)를 실시간 제공해야 한다. | B | device-service | P0 |
| REQ-FUNC-DEV-004 | 시스템은 BLE 양방향 스트림을 통해 디바이스 상태(배터리, 신호강도)를 실시간 모니터링해야 한다. | B | device-service | P1 |
| REQ-FUNC-DEV-005 | 시스템은 리더기 펌웨어 OTA 업데이트를 요청하고 진행 상태를 추적해야 한다. | B | device-service | P2 |
| REQ-FUNC-DEV-006 | 시스템은 디바이스 이벤트(연결, 해제, 에러, 보정)를 기록하고 이력을 조회할 수 있어야 한다. | A | device-service | P1 |

### 4.3 카트리지 관리

| REQ ID | 요구사항 | 안전 등급 | 서비스 | 우선순위 |
|--------|---------|----------|--------|---------|
| REQ-FUNC-CART-001 | 시스템은 NFC ISO 14443A를 통해 카트리지 태그를 자동으로 인식하고 파싱(UID, 타입, 로트, 유효기간, 잔여횟수, 보정계수)해야 한다. | B | cartridge-service, Rust Core | P0 |
| REQ-FUNC-CART-002 | 시스템은 무한 확장 카트리지 레지스트리(2-byte 코드, 최대 65,536종)를 관리해야 한다. 기본 29종이 사전 등록된다. | B | cartridge-service | P0 |
| REQ-FUNC-CART-003 | 시스템은 구독 티어별 카트리지 접근 제어(INCLUDED/LIMITED/ADD_ON/RESTRICTED/BETA)를 적용해야 한다. | A | cartridge-service | P1 |
| REQ-FUNC-CART-004 | 시스템은 카트리지 사용 횟수를 추적하고, 잔여 횟수가 0인 카트리지 사용을 차단해야 한다. | B | cartridge-service | P0 |
| REQ-FUNC-CART-005 | 시스템은 카트리지 유효기간을 검증하고, 만료된 카트리지 사용 시 경고를 표시해야 한다. | B | cartridge-service | P0 |
| REQ-FUNC-CART-006 | 시스템은 새 카트리지 타입의 OTA 배포를 지원해야 한다. | A | cartridge-service | P3 |

### 4.4 측정 및 분석

| REQ ID | 요구사항 | 안전 등급 | 서비스 | 우선순위 |
|--------|---------|----------|--------|---------|
| REQ-FUNC-MEAS-001 | 시스템은 측정 세션을 생성(StartSession)하고 세션 ID를 발급해야 한다. 세션에는 디바이스 ID, 카트리지 ID, 사용자 ID가 기록된다. | B | measurement-service | P0 |
| REQ-FUNC-MEAS-002 | 시스템은 양방향 gRPC 스트림(StreamMeasurement)을 통해 실시간 측정 데이터를 수집해야 한다. 각 메시지에는 원시 채널 데이터, 차동측정 값, 환경 메타데이터(온도, 습도, 기압)가 포함된다. | B | measurement-service | P0 |
| REQ-FUNC-MEAS-003 | Rust 코어 엔진은 차동측정 공식 `S_corrected = S_det − α × S_ref`를 구현하고, 채널별 오프셋·게인·온도 보정을 적용해야 한다. α 기본값은 0.95이며 카트리지별 조정이 가능하다. | B | Rust Core | P0 |
| REQ-FUNC-MEAS-004 | Rust 코어 엔진은 88/448/896차원 핑거프린트 벡터를 생성해야 한다. 벡터는 L2 정규화되며 코사인 유사도 검색을 지원한다. | B | Rust Core | P0 |
| REQ-FUNC-MEAS-005 | 시스템은 측정 세션을 종료(EndSession)하고 총 측정 횟수 및 종료 시각을 반환해야 한다. | B | measurement-service | P0 |
| REQ-FUNC-MEAS-006 | 시스템은 사용자별 측정 이력을 날짜 범위, 페이지네이션으로 조회(GetMeasurementHistory)할 수 있어야 한다. | B | measurement-service | P0 |
| REQ-FUNC-MEAS-007 | 시스템은 측정 데이터를 PostgreSQL(메타데이터), TimescaleDB(시계열), Milvus(핑거프린트 벡터)에 분산 저장해야 한다. | B | measurement-service | P0 |
| REQ-FUNC-MEAS-008 | 시스템은 측정 완료 이벤트(measurement.completed)를 Kafka로 발행하여 AI 추론·코칭·알림 서비스에 전달해야 한다. | B | measurement-service | P1 |
| REQ-FUNC-MEAS-009 | 시스템은 FHIR R4 Observation 형식으로 측정 데이터 내보내기를 지원해야 한다. | A | measurement-service | P3 |

### 4.5 AI 추론

| REQ ID | 요구사항 | 안전 등급 | 서비스 | 우선순위 |
|--------|---------|----------|--------|---------|
| REQ-FUNC-AI-001 | 시스템은 서버 측 AI 추론 서비스를 제공해야 한다. 측정 데이터를 입력받아 바이오마커 분석 결과(예측값, 단위, 신뢰구간, 위험도)를 반환한다. | B | ai-inference-service | P1 |
| REQ-FUNC-AI-002 | Rust 코어 엔진은 TFLite 기반 엣지 AI 추론을 지원해야 한다. 5종 모델(Calibration, FingerprintClassifier, AnomalyDetection, ValuePredictor, QualityAssessment)을 네트워크 없이 실행한다. | B | Rust Core | P0 |
| REQ-FUNC-AI-003 | AI 추론 서비스는 모델 버전을 관리하고, 활성 모델 목록을 조회할 수 있어야 한다. | A | ai-inference-service | P2 |
| REQ-FUNC-AI-004 | 시스템은 AI 분석 완료 이벤트(ai.analysis_completed)를 Kafka로 발행해야 한다. | A | ai-inference-service | P2 |
| REQ-FUNC-AI-005 | 시스템은 위험 예측 및 이상 감지 결과에 따라 사용자에게 경고 알림을 전송해야 한다. | B | ai-inference-service, notification-service | P2 |

### 4.6 보정 (Calibration)

| REQ ID | 요구사항 | 안전 등급 | 서비스 | 우선순위 |
|--------|---------|----------|--------|---------|
| REQ-FUNC-CAL-001 | 시스템은 팩토리 보정(제조 시 초기 보정)을 수행하고 보정 인증서를 기록해야 한다. | B | calibration-service | P1 |
| REQ-FUNC-CAL-002 | 시스템은 현장 보정(사용자 환경 보정)을 수행하고 보정 결과(오프셋, 게인, R² 점수)를 저장해야 한다. | B | calibration-service | P1 |
| REQ-FUNC-CAL-003 | 시스템은 22종 보정 모델을 관리하고, 카트리지 타입별 적합한 보정 모델을 자동 선택해야 한다. | B | calibration-service | P1 |
| REQ-FUNC-CAL-004 | 시스템은 보정 이력을 추적하고, 보정 만료 시 재보정을 알림으로 안내해야 한다. | B | calibration-service | P2 |

### 4.7 구독 관리

| REQ ID | 요구사항 | 안전 등급 | 서비스 | 우선순위 |
|--------|---------|----------|--------|---------|
| REQ-FUNC-SUB-001 | 시스템은 4단계 구독 티어(Free/Basic ₩9,900/Pro ₩29,900/Clinical ₩59,900)를 관리해야 한다. | A | subscription-service | P1 |
| REQ-FUNC-SUB-002 | 시스템은 구독 생성, 업그레이드, 다운그레이드, 취소를 지원해야 한다. | A | subscription-service | P1 |
| REQ-FUNC-SUB-003 | 시스템은 구독 만료 시 자동으로 Free 티어로 변경해야 한다. | A | subscription-service | P1 |
| REQ-FUNC-SUB-004 | 시스템은 구독 티어에 따라 최대 디바이스 수, 가족 구성원 수, AI 코칭 활성화 여부, 원격진료 활성화 여부를 제한해야 한다. | A | subscription-service | P1 |

### 4.8 결제

| REQ ID | 요구사항 | 안전 등급 | 서비스 | 우선순위 |
|--------|---------|----------|--------|---------|
| REQ-FUNC-PAY-001 | 시스템은 Toss Payments PG를 통한 결제 생성 및 승인을 지원해야 한다. | A | payment-service | P1 |
| REQ-FUNC-PAY-002 | 시스템은 결제 취소(전체/부분 취소)를 지원해야 한다. | A | payment-service | P1 |
| REQ-FUNC-PAY-003 | 시스템은 결제 완료 이벤트(payment.completed)를 Kafka로 발행하여 구독 활성화, 알림 전송을 트리거해야 한다. | A | payment-service | P1 |
| REQ-FUNC-PAY-004 | 시스템은 결제 이력을 조회할 수 있어야 한다. | A | payment-service | P1 |

### 4.9 상점 (Commerce)

| REQ ID | 요구사항 | 안전 등급 | 서비스 | 우선순위 |
|--------|---------|----------|--------|---------|
| REQ-FUNC-SHOP-001 | 시스템은 상품 등록, 조회, 카테고리 분류를 지원해야 한다. | A | shop-service | P1 |
| REQ-FUNC-SHOP-002 | 시스템은 장바구니 관리(추가, 수정, 삭제)를 지원해야 한다. | A | shop-service | P1 |
| REQ-FUNC-SHOP-003 | 시스템은 주문 생성 및 상태 추적(대기→확인→배송→완료)을 지원해야 한다. | A | shop-service | P1 |

### 4.10 코칭

| REQ ID | 요구사항 | 안전 등급 | 서비스 | 우선순위 |
|--------|---------|----------|--------|---------|
| REQ-FUNC-COACH-001 | 시스템은 AI 기반 건강 코칭 프로그램 시작 및 종료를 지원해야 한다. | A | coaching-service | P2 |
| REQ-FUNC-COACH-002 | 시스템은 일일/주간 건강 리포트를 자동 생성해야 한다. | A | coaching-service | P2 |
| REQ-FUNC-COACH-003 | 시스템은 개인 건강 목표 설정 및 진행률 추적을 지원해야 한다. | A | coaching-service | P2 |
| REQ-FUNC-COACH-004 | 시스템은 개인화된 건강 추천(운동, 식단, 생활습관)을 제공해야 한다. | A | coaching-service | P2 |

### 4.11 원격진료

| REQ ID | 요구사항 | 안전 등급 | 서비스 | 우선순위 |
|--------|---------|----------|--------|---------|
| REQ-FUNC-TELE-001 | 시스템은 화상진료 세션 생성 및 종료를 지원해야 한다. | B | telemedicine-service | P3 |
| REQ-FUNC-TELE-002 | 시스템은 진료 세션에 측정 데이터를 첨부할 수 있어야 한다. | B | telemedicine-service | P3 |
| REQ-FUNC-TELE-003 | 시스템은 진료 기록(메모, 진단, 녹화 URL)을 저장하고 조회할 수 있어야 한다. | B | telemedicine-service | P3 |
| REQ-FUNC-TELE-004 | 시스템은 의사 목록을 전문분야, 가능 시간대로 검색할 수 있어야 한다. | A | telemedicine-service | P3 |

### 4.12 처방

| REQ ID | 요구사항 | 안전 등급 | 서비스 | 우선순위 |
|--------|---------|----------|--------|---------|
| REQ-FUNC-RX-001 | 시스템은 의료진이 처방전을 생성할 수 있어야 한다. 처방전에는 약품 목록, 용법·용량이 포함된다. | B | prescription-service | P3 |
| REQ-FUNC-RX-002 | 시스템은 처방전을 약국에 전자 전송할 수 있어야 한다. | B | prescription-service | P3 |
| REQ-FUNC-RX-003 | 시스템은 처방 이행 상태(대기→발송→조제→수령)를 추적해야 한다. | B | prescription-service | P3 |

### 4.13 건강 기록

| REQ ID | 요구사항 | 안전 등급 | 서비스 | 우선순위 |
|--------|---------|----------|--------|---------|
| REQ-FUNC-HR-001 | 시스템은 사용자 건강 기록을 FHIR R4 형식으로 생성·저장해야 한다. | B | health-record-service | P3 |
| REQ-FUNC-HR-002 | 시스템은 가족 건강 리포트를 생성할 수 있어야 한다. | A | health-record-service | P3 |
| REQ-FUNC-HR-003 | 시스템은 건강 기록을 외부 시스템(HealthKit, Google Health)과 연동할 수 있어야 한다. | A | health-record-service | P4 |

### 4.14 가족 관리

| REQ ID | 요구사항 | 안전 등급 | 서비스 | 우선순위 |
|--------|---------|----------|--------|---------|
| REQ-FUNC-FAM-001 | 시스템은 가족 그룹 생성 및 구성원 초대/삭제를 지원해야 한다. | A | family-service | P3 |
| REQ-FUNC-FAM-002 | 시스템은 가족 그룹 내 보호자-피보호자 역할을 관리해야 한다. | A | family-service | P3 |
| REQ-FUNC-FAM-003 | 시스템은 보호자가 피보호자의 측정 데이터를 모니터링할 수 있어야 한다. | B | family-service | P3 |

### 4.15 알림

| REQ ID | 요구사항 | 안전 등급 | 서비스 | 우선순위 |
|--------|---------|----------|--------|---------|
| REQ-FUNC-NOTI-001 | 시스템은 Firebase FCM을 통한 푸시 알림 전송을 지원해야 한다. | A | notification-service | P1 |
| REQ-FUNC-NOTI-002 | 시스템은 알림 이력 조회 및 읽음 처리를 지원해야 한다. | A | notification-service | P2 |
| REQ-FUNC-NOTI-003 | 시스템은 사용자별 알림 설정(수신 여부, 유형별 필터)을 관리해야 한다. | A | notification-service | P2 |

### 4.16 커뮤니티

| REQ ID | 요구사항 | 안전 등급 | 서비스 | 우선순위 |
|--------|---------|----------|--------|---------|
| REQ-FUNC-COMM-001 | 시스템은 커뮤니티 게시글/댓글 CRUD를 지원해야 한다. | A | community-service | P3 |
| REQ-FUNC-COMM-002 | 시스템은 게시글 전문 검색(Elasticsearch)을 지원해야 한다. | A | community-service | P3 |
| REQ-FUNC-COMM-003 | 시스템은 게시글 좋아요, 신고 기능을 지원해야 한다. | A | community-service | P3 |

### 4.17 예약

| REQ ID | 요구사항 | 안전 등급 | 서비스 | 우선순위 |
|--------|---------|----------|--------|---------|
| REQ-FUNC-RESV-001 | 시스템은 병원/약국 검색(지역, 전문분야, 가용 시간)을 지원해야 한다. | A | reservation-service | P3 |
| REQ-FUNC-RESV-002 | 시스템은 예약 생성, 수정, 취소를 지원해야 한다. | A | reservation-service | P3 |

### 4.18 관리자 기능

| REQ ID | 요구사항 | 안전 등급 | 서비스 | 우선순위 |
|--------|---------|----------|--------|---------|
| REQ-FUNC-ADMIN-001 | 시스템은 관리자가 시스템 설정을 카테고리별(결제, 알림, 보안, AI 등 8개 카테고리)로 관리할 수 있어야 한다. 민감 설정은 AES-256-GCM으로 암호화 저장한다. | A | admin-service | P1 |
| REQ-FUNC-ADMIN-002 | 시스템은 모든 주요 관리 작업에 대해 감사 로그를 기록해야 한다. 감사 로그에는 작업자, 시각, 변경 전/후 값이 포함된다. | B | admin-service | P0 |
| REQ-FUNC-ADMIN-003 | 시스템은 RBAC 5개 역할(SuperAdmin, Admin, Moderator, Support, Analyst)을 지원해야 한다. | B | admin-service | P1 |
| REQ-FUNC-ADMIN-004 | 시스템은 설정 변경 시 config.changed 이벤트를 Kafka로 발행하여 관련 서비스(payment, notification 등)에 실시간 반영해야 한다. | A | admin-service | P1 |

### 4.19 기타 서비스

| REQ ID | 요구사항 | 안전 등급 | 서비스 | 우선순위 |
|--------|---------|----------|--------|---------|
| REQ-FUNC-TRANS-001 | 시스템은 6개 언어(ko/en/ja/zh/fr/hi) 실시간 번역을 지원해야 한다. | A | translation-service | P3 |
| REQ-FUNC-VIDEO-001 | 시스템은 비디오 업로드, 트랜스코딩, 스트리밍을 지원해야 한다. | A | video-service | P3 |
| REQ-FUNC-VISION-001 | 시스템은 음식 사진을 분석하여 칼로리/영양소를 추정해야 한다. | B | vision-service | P3 |
| REQ-FUNC-GW-001 | API Gateway는 모든 외부 요청을 인증·라우팅·레이트리밋 처리해야 한다. | A | gateway | P0 |
| REQ-FUNC-GW-002 | API Gateway는 MinIO를 통한 파일 업로드/다운로드/삭제를 지원해야 한다. | A | gateway | P1 |

### 4.20 오프라인 기능

| REQ ID | 요구사항 | 안전 등급 | 서비스 | 우선순위 |
|--------|---------|----------|--------|---------|
| REQ-FUNC-OFF-001 | 시스템은 네트워크 연결 없이 차동측정, 핑거프린트 생성, 엣지 AI 추론을 100% 수행할 수 있어야 한다. | B | Rust Core | P0 |
| REQ-FUNC-OFF-002 | 시스템은 오프라인 상태에서 최소 72시간 분량의 측정 데이터를 로컬(SQLite CRDT)에 보존해야 한다. | B | Rust Core | P0 |
| REQ-FUNC-OFF-003 | 시스템은 온라인 복귀 시 CRDT 기반 자동 동기화를 수행해야 한다. 100건 동기화는 30초 이내 완료되어야 한다. | B | Rust Core | P1 |

---

## 5. 비기능 요구사항 (Clause 5.2.3)

### 5.1 성능 요구사항

| NFR ID | 유형 | 요구사항 | 목표값 | 측정 기준 |
|--------|------|---------|-------|----------|
| REQ-NFR-PERF-001 | 응답시간 | 인증 API P95 응답시간 | < 300ms | 서버 사이드, Kong 이후 |
| REQ-NFR-PERF-002 | 응답시간 | 조회 API P95 응답시간 | < 150ms | 서버 사이드 |
| REQ-NFR-PERF-003 | 응답시간 | 측정 API P95 응답시간 | < 200ms | 서버 사이드 |
| REQ-NFR-PERF-004 | 응답시간 | 스트리밍 메시지 지연 P95 | < 50ms | gRPC 단방향 |
| REQ-NFR-PERF-005 | 응답시간 | 벡터 검색(Milvus) P95 | < 200ms | 896차원, 100만 벡터 |
| REQ-NFR-PERF-006 | 응답시간 | AI 엣지 추론(TFLite) P95 | < 200ms | 모바일 디바이스 |
| REQ-NFR-PERF-007 | 처리량 | Phase 3 동시 접속 | 10,000 CCU | — |
| REQ-NFR-PERF-008 | 처리량 | Phase 3 초당 요청 | 50,000 RPS | — |
| REQ-NFR-PERF-009 | 처리량 | Phase 3 초당 측정 | 2,000 TPS | — |
| REQ-NFR-PERF-010 | Rust 코어 | 차동측정 88ch 단일 | < 1μs | criterion 벤치마크 |
| REQ-NFR-PERF-011 | Rust 코어 | 핑거프린트 896차원 생성 | < 100μs | FingerprintBuilder |
| REQ-NFR-PERF-012 | Rust 코어 | AES-256-GCM 암호화 1KB | < 50μs | ring 기반 |

### 5.2 가용성 요구사항

| NFR ID | 유형 | 요구사항 | 목표값 |
|--------|------|---------|-------|
| REQ-NFR-AVAIL-001 | 가용성 | Phase 3 서비스 가용률 | 99.9% (43.2분/월 허용) |
| REQ-NFR-AVAIL-002 | 복구시간 | RTO (복구 시간 목표) | < 5분 |
| REQ-NFR-AVAIL-003 | 복구시점 | RPO (복구 시점 목표) | < 1분 |
| REQ-NFR-AVAIL-004 | 평균복구 | MTTR (평균 복구 시간) | < 15분 |
| REQ-NFR-AVAIL-005 | 장애간격 | MTBF (평균 장애 간격) | > 720시간 |
| REQ-NFR-AVAIL-006 | 오프라인 | 오프라인 측정 가용률 | 100% |
| REQ-NFR-AVAIL-007 | 오프라인 | 오프라인 데이터 보존 | ≥ 72시간 분량 |

### 5.3 확장성 요구사항

| NFR ID | 유형 | 요구사항 | 목표값 |
|--------|------|---------|-------|
| REQ-NFR-SCALE-001 | 수평 확장 | 서비스별 자동 확장 | K8s HPA (CPU 70%, 메모리 80%) |
| REQ-NFR-SCALE-002 | DB 확장 | TimescaleDB 시계열 파티셔닝 | 자동 하이퍼테이블 |
| REQ-NFR-SCALE-003 | 벡터 확장 | Milvus 벡터 인덱스 | IVF_FLAT, 100만→1억 벡터 |
| REQ-NFR-SCALE-004 | 메시지 확장 | Kafka 파티션 | 토픽당 최소 3 파티션 |

### 5.4 보안 요구사항

| NFR ID | 유형 | 요구사항 | 목표값 |
|--------|------|---------|-------|
| REQ-NFR-SEC-001 | 저장 암호화 | PHI/PII 저장 데이터 | AES-256-GCM |
| REQ-NFR-SEC-002 | 전송 암호화 | 모든 네트워크 통신 | TLS 1.3 |
| REQ-NFR-SEC-003 | 인증 | 사용자 인증 방식 | JWT + OAuth2 (Keycloak OIDC) |
| REQ-NFR-SEC-004 | 인가 | 역할 기반 접근 제어 | RBAC 5개 역할 |
| REQ-NFR-SEC-005 | 비밀번호 | 비밀번호 해시 저장 | bcrypt cost ≥ 12 |
| REQ-NFR-SEC-006 | 토큰 | Access Token 만료 | 15분 |
| REQ-NFR-SEC-007 | 토큰 | Refresh Token 만료 | 7일 |
| REQ-NFR-SEC-008 | 레이트리밋 | API Rate Limiting | 60 req/min (기본) |
| REQ-NFR-SEC-009 | 해시체인 | 측정 데이터 무결성 | SHA-256 해시 |
| REQ-NFR-SEC-010 | 입력검증 | 모든 사용자 입력 | 서버 측 엄격 검증 |
| REQ-NFR-SEC-011 | 감사 | 주요 작업 감사 로그 | 10년 보존 |

### 5.5 규정 준수 요구사항

| NFR ID | 유형 | 요구사항 | 목표값 |
|--------|------|---------|-------|
| REQ-NFR-REG-001 | 동의관리 | 개인정보 동의 관리 | GDPR/PIPA 준수 |
| REQ-NFR-REG-002 | 데이터보존 | 의료 데이터 보존 기간 | 최소 10년 |
| REQ-NFR-REG-003 | 감사추적 | 감사 추적 기록 | 모든 주요 작업 |
| REQ-NFR-REG-004 | 데이터삭제 | 사용자 탈퇴 시 데이터 삭제 | 30일 내 완전 삭제 |
| REQ-NFR-REG-005 | 데이터이전 | 데이터 포터빌리티 | FHIR R4 형식 내보내기 |

### 5.6 국제화 요구사항

| NFR ID | 유형 | 요구사항 | 목표값 |
|--------|------|---------|-------|
| REQ-NFR-I18N-001 | 다국어 | UI 다국어 지원 | 6개 언어 (ko/en/ja/zh/fr/hi) |
| REQ-NFR-I18N-002 | 시간대 | 사용자별 시간대 지원 | 전체 IANA 시간대 |
| REQ-NFR-I18N-003 | 통화 | 다중 통화 지원 | KRW, USD, EUR, JPY, CNY |

### 5.7 사용성 요구사항

| NFR ID | 유형 | 요구사항 | 목표값 |
|--------|------|---------|-------|
| REQ-NFR-UX-001 | 접근성 | 다크 모드 지원 | Material Design 3 테마 |
| REQ-NFR-UX-002 | 접근성 | 화면 크기 대응 | 반응형 레이아웃 |
| REQ-NFR-UX-003 | 사용성 | 첫 측정까지 걸리는 시간 | < 5분 (온보딩 포함) |

---

## 6. 인터페이스 요구사항 (Clause 5.2.4)

### 6.1 외부 하드웨어 인터페이스

| REQ-IF ID | 인터페이스 | 프로토콜 | 방향 | 상세 |
|--------|-----------|---------|------|------|
| REQ-IF-001 | 만파식 리더기 ↔ 모바일 앱 | BLE 5.0 (GATT) | 양방향 | GATT 서비스 UUID `0000fff0-...`, 측정 데이터 Notify, 명령 Write |
| REQ-IF-002 | 카트리지 ↔ 모바일 앱 | NFC ISO 14443A | 단방향 (읽기) | 53+ 바이트 태그: UID(8B) + 타입(1B) + 로트(8B) + 유효기간(8B) + 잔여횟수(2B) + 보정계수(24B) |
| REQ-IF-003 | 리더기 배터리 서비스 | BLE (Battery Service) | 읽기 | UUID `0000180f-...`, 배터리 레벨(u8, %) |

### 6.2 내부 서비스 인터페이스

| REQ-IF ID | 인터페이스 | 프로토콜 | 상세 |
|--------|-----------|---------|------|
| REQ-IF-004 | 서비스 간 통신 | gRPC (HTTP/2 + Protobuf) | 30+ 서비스 상호 호출, Proto: `manpasik.proto` |
| REQ-IF-005 | 비동기 이벤트 | Kafka (Redpanda) | 21개 토픽, JSON 스키마, DLQ 정책 |
| REQ-IF-006 | 벡터 검색 | gRPC (Milvus SDK) | 896차원 핑거프린트 유사도 검색 |
| REQ-IF-007 | 전문 검색 | HTTP REST (Elasticsearch) | 측정·커뮤니티 콘텐츠 검색 |
| REQ-IF-008 | 캐시 | Redis Protocol (TCP 6379) | 세션, 디바이스 상태, 구독 캐시 |
| REQ-IF-009 | 오브젝트 스토리지 | S3 API (MinIO) | 파일 업로드/다운로드/삭제 |

### 6.3 외부 서비스 인터페이스

| REQ-IF ID | 인터페이스 | 프로토콜 | 방향 | 상세 |
|--------|-----------|---------|------|------|
| REQ-IF-010 | 모바일 앱 ↔ API Gateway | gRPC-Web / REST | 외부 | TLS 1.3, JWT Bearer 인증 |
| REQ-IF-011 | payment-service ↔ Toss Payments | HTTPS REST | 외부 | 결제 승인/취소 API |
| REQ-IF-012 | notification-service ↔ Firebase FCM | HTTPS REST | 외부 | 푸시 알림 전송 |
| REQ-IF-013 | auth-service ↔ Keycloak | OIDC/OAuth2 | 내부/외부 | 인증 위임, MFA |
| REQ-IF-014 | Rust 코어 ↔ Flutter 앱 | FFI (flutter_rust_bridge) | 내부 | 10 API: 차동측정, 핑거프린트, AI, BLE, NFC 등 |

### 6.4 Kafka 이벤트 토픽 (주요)

| REQ-IF ID | 토픽명 | 발행 서비스 | 소비 서비스 | 페이로드 요약 |
|--------|--------|-----------|-----------|-------------|
| REQ-IF-EVT-001 | manpasik.measurement.completed | measurement | ai-inference, coaching, notification | session_id, user_id, fingerprint, primary_value |
| REQ-IF-EVT-002 | manpasik.payment.completed | payment | subscription, notification | payment_id, user_id, amount, status |
| REQ-IF-EVT-003 | manpasik.subscription.created | subscription | notification | subscription_id, user_id, tier |
| REQ-IF-EVT-004 | manpasik.config.changed | admin | payment, notification, 전체 | config_category, key, updated_by |
| REQ-IF-EVT-005 | manpasik.calibration.completed | calibration | measurement | calibration_id, device_id, model, r_squared |
| REQ-IF-EVT-006 | manpasik.ai.analysis_completed | ai-inference | coaching, notification | analysis_id, user_id, risk_level |
| REQ-IF-EVT-007 | manpasik.device.status_changed | device | notification | device_id, old_status, new_status |

---

## 7. 데이터 요구사항 (REQ-DATA-xxx)

### 7.0 데이터 요구사항 명세

| REQ-DATA ID | 요구사항 | 저장소 | 관련 기능 |
|-------------|---------|--------|----------|
| REQ-DATA-001 | 사용자·인증 데이터는 PostgreSQL에 관계형 저장 | PostgreSQL 16 | REQ-FUNC-AUTH-001~007, REQ-FUNC-USER-001 |
| REQ-DATA-002 | 측정 시계열 데이터는 TimescaleDB 하이퍼테이블에 저장 | TimescaleDB | REQ-FUNC-MEAS-002, REQ-FUNC-MEAS-007 |
| REQ-DATA-003 | 핑거프린트 벡터는 Milvus에 저장하여 유사도 검색 지원 | Milvus 2.4 | REQ-FUNC-MEAS-004, REQ-FUNC-MEAS-007 |
| REQ-DATA-004 | 세션·디바이스·구독 캐시는 Redis에 TTL 적용 저장 | Redis 7 | REQ-FUNC-GW-001, REQ-FUNC-DEV-003, REQ-FUNC-SUB-001 |
| REQ-DATA-005 | 비동기 이벤트는 Kafka 토픽으로 발행·소비 | Kafka (Redpanda) | REQ-FUNC-MEAS-008, REQ-FUNC-PAY-003, REQ-FUNC-ADMIN-004 |
| REQ-DATA-006 | 전문 검색·로그는 Elasticsearch에 인덱싱 | Elasticsearch 8.14 | REQ-FUNC-MEAS-006, REQ-FUNC-COMM-002 |
| REQ-DATA-007 | 파일·이미지는 MinIO S3 API로 오브젝트 저장 | MinIO | REQ-FUNC-GW-002, REQ-FUNC-VISION-001 |
| REQ-DATA-008 | PHI/PII 저장 시 AES-256-GCM 암호화 적용 | 전체 DB | REQ-NFR-SEC-001, REQ-REG-015 |
| REQ-DATA-009 | 의료 데이터 최소 10년 보존 | PostgreSQL, TimescaleDB | REQ-NFR-REG-002, REQ-REG-019 |
| REQ-DATA-010 | 측정 데이터 무결성 SHA-256 해시체인 보장 | 측정 스키마 | REQ-SAFE-008, REQ-NFR-SEC-009 |

### 7.1 데이터 저장소 매핑

| 데이터 유형 | 저장소 | 보존 기간 | 암호화 | 백업 주기 |
|-----------|--------|----------|--------|---------|
| 사용자 정보 (PII) | PostgreSQL | 탈퇴 후 30일 | AES-256-GCM | 일일 |
| 인증 자격증명 | PostgreSQL | 계정 수명 | bcrypt 해시 | 일일 |
| 측정 메타데이터 | PostgreSQL | 10년 | AES-256-GCM | 일일 |
| 측정 시계열 | TimescaleDB | 10년 | AES-256-GCM | 일일 |
| 핑거프린트 벡터 | Milvus | 10년 | N/A (내부) | 일일 |
| 건강 기록 (PHI) | PostgreSQL | 10년 | AES-256-GCM | 일일 |
| 처방 기록 | PostgreSQL | 10년 | AES-256-GCM | 일일 |
| 감사 로그 | PostgreSQL + Elasticsearch | 10년 | 무결성 해시 (SHA-256) | 일일 |
| 결제 기록 | PostgreSQL | 5년 | AES-256-GCM | 일일 |
| 디바이스 상태 캐시 | Redis | TTL 기반 | N/A | — |
| JWT 세션 | Redis | 토큰 만료 시 | N/A | — |
| 시스템 로그 | Elasticsearch | 90일 | N/A | 주간 |
| 파일/이미지 | MinIO | 사용자 수명 | 서버 측 암호화 | 일일 |
| 이벤트 스트림 | Kafka | 7일 (기본) | TLS in-transit | — |

### 7.2 데이터베이스 스키마 개요

| DB 초기화 스크립트 | 도메인 | 주요 테이블 |
|------------------|--------|-----------|
| `01-auth.sql` | 인증 | users, user_credentials, refresh_tokens |
| `02-device.sql` | 디바이스 | devices, device_events, device_firmware |
| `03-measurement.sql` | 측정 | measurements (hypertable), measurement_sessions |
| `05-subscription.sql` | 구독 | subscriptions, subscription_history |
| `06-shop.sql` | 상점 | products, carts, cart_items, orders, order_items |
| `07-payment.sql` | 결제 | payments, payment_refunds |
| `08-ai-inference.sql` | AI | ai_models, ai_analyses, ai_predictions |
| `09-cartridge.sql` | 카트리지 | cartridge_types, cartridge_registry, cartridge_usage |
| `10-calibration.sql` | 보정 | calibrations, calibration_models, calibration_certs |
| `11-coaching.sql` | 코칭 | coaching_programs, goals, daily_reports |
| `12-notification.sql` | 알림 | notifications, notification_settings |
| `13-family.sql` | 가족 | family_groups, family_members |
| `14-health-record.sql` | 건강기록 | health_records, fhir_observations |
| `15-telemedicine.sql` | 원격진료 | telemedicine_sessions, session_notes |
| `16-reservation.sql` | 예약 | facilities, doctors, reservations |
| `17-community.sql` | 커뮤니티 | posts, comments, likes, reports |
| `18-admin.sql` | 관리 | admin_settings, audit_logs |
| `19-prescription.sql` | 처방 | prescriptions, prescription_items, fulfillments |
| `20-translation.sql` | 번역 | translation_keys, translations |
| `21-video.sql` | 비디오 | videos, video_metadata |

### 7.3 데이터 무결성

| 항목 | 요구사항 |
|------|---------|
| 측정 데이터 | SHA-256 해시로 변조 방지, 해시체인으로 순서 보장 |
| 감사 로그 | 삭제 불가(append-only), 무결성 해시 적용 |
| 결제 기록 | 트랜잭션 원자성 보장, 이중 기록 방지 (멱등성 키) |
| 참조 무결성 | PostgreSQL 외래키 제약조건 적용 |
| 동시성 제어 | 낙관적 잠금 (updated_at 기반) |

---

## 8. 안전 요구사항 (Clause 5.2.5)

### 8.1 위험관리 기반 안전 요구사항

ISO 14971 위험관리 계획서에서 식별된 9개 위해 시나리오 및 STRIDE 위협 모델링 31개 위협으로부터 도출된 소프트웨어 안전 통제 조치입니다.

| REQ-SAFE ID | 위해 시나리오 | 안전 통제 조치 | 관련 REQ | 검증 방법 |
|--------|-------------|-------------|---------|----------|
| REQ-SAFE-001 | 측정 데이터 부정확 → 잘못된 건강 판단 | 차동측정 보정 계수 검증, 채널별 범위 점검, 이상값 탐지 | REQ-FUNC-MEAS-003, REQ-FUNC-CAL-001 | 단위 테스트 + 벤치마크, 참조 물질 테스트 |
| REQ-SAFE-002 | 핑거프린트 벡터 오류 → 잘못된 유사도 분석 | L2 정규화 검증, 차원 유효성 검사 (88/448/896), NaN/Inf 방어 | REQ-FUNC-MEAS-004 | 단위 테스트 (코사인 유사도 범위 0~1) |
| REQ-SAFE-003 | AI 추론 오류 → 잘못된 바이오마커 분석 | 신뢰구간 표시 필수, 임계값 이하 시 "결과 불확실" 경고, 의료 전문가 확인 권고 문구 | REQ-FUNC-AI-001, REQ-FUNC-AI-002 | 모델 검증 (테스트셋 AUC ≥ 0.85) |
| REQ-SAFE-004 | 오프라인 데이터 손실 → 측정 기록 유실 | CRDT 로컬 저장 72시간 보존, 온라인 복귀 시 자동 동기화, 충돌 해결 로그 | REQ-FUNC-OFF-001, REQ-FUNC-OFF-002, REQ-FUNC-OFF-003 | 72시간 연속 오프라인 테스트 |
| REQ-SAFE-005 | 보정 오류 → 체계적 측정 편향 | 보정 R² ≥ 0.95 요구, 보정 만료 시 경고, 팩토리 보정 인증서 필수 | REQ-FUNC-CAL-001, REQ-FUNC-CAL-002 | 보정 정확도 테스트 |
| REQ-SAFE-006 | 카트리지 인식 오류 → 잘못된 측정 파라미터 | NFC 태그 CRC 검증, 카트리지 타입-채널 수 교차 검증, 유효기간/잔여횟수 검증 | REQ-FUNC-CART-001, REQ-FUNC-CART-004, REQ-FUNC-CART-005 | NFC 파싱 단위 테스트 |
| REQ-SAFE-007 | 인증 우회 → 무단 의료 데이터 접근 | JWT 서명 검증(RS256), 토큰 만료 엄격 적용, RBAC 적용, 민감 작업 재인증 | REQ-FUNC-AUTH-002, REQ-FUNC-ADMIN-003 | 보안 침투 테스트 |
| REQ-SAFE-008 | 데이터 변조 → 측정 이력 위변조 | SHA-256 해시체인, 감사 로그 append-only, TLS 전송 암호화 | REQ-FUNC-ADMIN-002, REQ-NFR-SEC-009 | 무결성 검증 테스트 |
| REQ-SAFE-009 | 처방 오류 → 잘못된 약품/용량 전달 | 처방전 의사 전자서명 필수, 약품 상호작용 경고, 이중 확인 프로세스 | REQ-FUNC-RX-001, REQ-FUNC-RX-002 | 처방 워크플로우 E2E 테스트 |
| REQ-SAFE-010 | 서비스 장애 → 긴급 측정 불가 | 오프라인 모드 100% 동작, 서비스 자동 복구(K8s), 헬스체크 모니터링 | REQ-FUNC-OFF-001, REQ-NFR-AVAIL-001 | 장애 주입 테스트 (Chaos) |

### 8.2 STRIDE 위협 통제 매핑

| 위협 ID | STRIDE 유형 | 위협 설명 | 소프트웨어 통제 | 관련 NFR |
|---------|------------|---------|---------------|---------|
| THR-001 | Spoofing | 사용자 신원 위장 | 강력한 비밀번호 정책 + MFA | REQ-NFR-SEC-003, REQ-NFR-SEC-005 |
| THR-002 | Spoofing | JWT 토큰 탈취 | 토큰 만료(15분) + Refresh 회전 | REQ-NFR-SEC-006, REQ-NFR-SEC-007 |
| THR-005 | Tampering | 측정 데이터 변조 | SHA-256 해시체인 + TLS 1.3 | REQ-NFR-SEC-002, REQ-NFR-SEC-009 |
| THR-010 | Repudiation | 관리 작업 부인 | 감사 로그 (append-only, 10년) | REQ-NFR-SEC-011 |
| THR-015 | Info Disclosure | PHI 유출 | AES-256-GCM 저장 암호화 | REQ-NFR-SEC-001 |
| THR-020 | DoS | 서비스 거부 공격 | Rate Limiting 60 req/min + K8s HPA | REQ-NFR-SEC-008, REQ-NFR-SCALE-001 |
| THR-025 | Elevation | 권한 상승 | RBAC 역할 검증 + 최소 권한 원칙 | REQ-NFR-SEC-004 |
| THR-030 | Tampering | 카트리지 위조 | NFC 태그 서명 검증 + 서버 인증 | REQ-FUNC-CART-001 |
| THR-031 | Info Disclosure | BLE 통신 도청 | BLE AES-CCM 암호화 + 페어링 인증 | REQ-IF-001 |

> **전체 31개 위협 상세**: `docs/security/stride-threat-model.md` 참조

### 8.3 소프트웨어 안전 관련 경고/알림

| 경고 조건 | 심각도 | 사용자 표시 | 관련 REQ |
|----------|--------|-----------|---------|
| AI 신뢰구간 < 70% | Warning | "결과 불확실 — 재측정 권장" | REQ-FUNC-AI-001, REQ-SAFE-003 |
| 측정값이 정상 범위 밖 | Critical | "비정상 수치 감지 — 의료 전문가 상담 권장" | REQ-SAFE-001 |
| 카트리지 잔여 횟수 ≤ 5 | Info | "카트리지 교체 예정" | REQ-FUNC-CART-004 |
| 카트리지 유효기간 만료 | Error | "만료된 카트리지 — 사용 불가" | REQ-FUNC-CART-005, REQ-SAFE-006 |
| 보정 만료 (90일 초과) | Warning | "재보정 필요 — 정확도 보장 불가" | REQ-FUNC-CAL-004, REQ-SAFE-005 |
| 디바이스 배터리 < 10% | Warning | "배터리 부족 — 충전 필요" | REQ-FUNC-DEV-003 |
| 오프라인 72시간 초과 | Warning | "데이터 동기화 권장" | REQ-FUNC-OFF-002, REQ-SAFE-004 |

---

## 9. 규정 요구사항 (Clause 5.2.6)

### 9.1 의료기기 규정

| REG ID | 규정 | 관할 | 요구사항 | 적용 상태 |
|--------|------|------|---------|----------|
| REQ-REG-001 | IEC 62304:2015 | 국제 | 소프트웨어 수명주기 프로세스 (Class B) — SDP, SRS, SAD, V&V 필수 | 구현 중 |
| REQ-REG-002 | ISO 14971:2019 | 국제 | 위험관리 프로세스 — 9개 위해 시나리오 식별, 5×5 매트릭스 | 완료 |
| REQ-REG-003 | ISO 13485:2016 | 국제 | 품질경영시스템 — 설계 관리, 문서 관리, CAPA | 계획 중 |
| REQ-REG-004 | IEC 62366-1:2015 | 국제 | 사용적합성 공학 — 사용 오류 위험 분석 | 계획 중 |
| REQ-REG-005 | 의료기기법 (MFDS) | 한국 | KGMP 인증, 기술문서 심사, 허가 신청 | 계획 중 |
| REQ-REG-006 | FDA 510(k) | 미국 | 실질적 동등성 입증, SW Validation Guidance 준수 | 계획 중 |
| REQ-REG-007 | EU IVDR 2017/746 | EU | 체외진단 의료기기 규정, 기술문서 (Annex II/III) | 계획 중 |
| REQ-REG-008 | NMPA | 중국 | 의료기기 등록, 임상 평가 | 계획 중 |
| REQ-REG-009 | PMDA | 일본 | 의료기기 승인, STED 문서 | 계획 중 |

### 9.2 데이터 보호 규정

| REG ID | 규정 | 관할 | 요구사항 | 적용 상태 |
|--------|------|------|---------|----------|
| REQ-REG-010 | GDPR 2016/679 | EU | 동의 관리, 데이터 이전권, 잊혀질 권리, DPO 임명 | 설계 중 |
| REQ-REG-011 | 개인정보보호법 (PIPA) | 한국 | 동의 수집, 목적 외 이용 금지, 파기 의무 | 설계 중 |
| REQ-REG-012 | HIPAA | 미국 | PHI 보호, BAA 체결, 기술적·관리적·물리적 보호 | 설계 중 |
| REQ-REG-013 | PIPL | 중국 | 개인정보 국외 이전 제한, 보호 영향 평가 | 계획 중 |
| REQ-REG-014 | APPI | 일본 | 개인정보 취급 사업자 의무, 국외 이전 규칙 | 계획 중 |

### 9.3 규정 준수 소프트웨어 요구사항

| REQ-REG ID | 요구사항 | 관련 NFR/REQ |
|--------|---------|-------------|
| REQ-REG-015 | 모든 PHI 저장 시 AES-256-GCM 암호화 적용 | REQ-NFR-SEC-001 |
| REQ-REG-016 | 모든 네트워크 통신에 TLS 1.3 적용 | REQ-NFR-SEC-002 |
| REQ-REG-017 | 사용자 동의 없이 개인정보 수집/처리 금지 | REQ-NFR-REG-001 |
| REQ-REG-018 | 데이터 주체의 접근/정정/삭제 요청 처리 기능 | REQ-NFR-REG-004 |
| REQ-REG-019 | 의료 데이터 최소 10년 보존 | REQ-NFR-REG-002 |
| REQ-REG-020 | 감사 추적 기록 10년 보존, 삭제 불가 | REQ-NFR-SEC-011, REQ-NFR-REG-003 |
| REQ-REG-021 | 소프트웨어 변경 시 위험 영향 분석 수행 | REQ-SAFE-001~010 |
| REQ-REG-022 | 모든 요구사항에 대한 양방향 추적성 유지 | §11 추적성 |

> **상세 146항목 체크리스트**: `docs/compliance/regulatory-compliance-checklist.md` 참조

---

## 10. 사용 환경

### 10.1 클라이언트 환경

| 항목 | 최소 요구사항 | 권장 |
|------|-------------|------|
| Android 버전 | 10 (API 29) | 13+ |
| iOS 버전 | 15.0 | 17.0+ |
| Bluetooth | BLE 5.0 | BLE 5.2 |
| NFC | ISO 14443A 지원 | — |
| 네트워크 | Wi-Fi 또는 LTE (오프라인 지원) | 5G |
| 저장 공간 | 200MB (앱) + 500MB (오프라인 데이터) | 2GB |

### 10.2 서버 환경

| 항목 | 요구사항 |
|------|---------|
| 서버 OS | Linux (Ubuntu 22.04 LTS+) |
| 컨테이너 | Docker 24.x + Kubernetes 1.28+ |
| 데이터베이스 | PostgreSQL 16, Redis 7, Milvus 2.4, Elasticsearch 8.14 |
| 메시징 | Kafka (Redpanda 24.2) |
| 스토리지 | MinIO (S3 호환) |
| 모니터링 | Prometheus + Grafana 11 + OpenTelemetry |

### 10.3 네트워크 환경

| 항목 | 요구사항 |
|------|---------|
| 프로토콜 | gRPC (HTTP/2), REST (HTTPS), WebSocket, BLE, NFC |
| 최소 대역폭 | 1 Mbps (측정 데이터 전송) |
| TLS | 1.3 필수 |
| 포트 | gRPC 50050-50072, PostgreSQL 5432, Redis 6379, Kafka 19092 |

---

## 11. 추적성 (Clause 5.2.7)

### 11.1 추적성 매트릭스 구조

모든 요구사항은 다음 4단계 양방향 추적을 유지합니다:

```
요구사항 (REQ/NFR/SAF/REG)
    ↕
설계 (DES) — SAD 소프트웨어 항목, Proto 정의, DB 스키마
    ↕
구현 (IMP) — 소스 코드 경로, 함수/메서드
    ↕
검증 (V&V) — 단위 테스트, 통합 테스트, E2E 테스트
```

### 11.2 추적성 ID 체계

| 접두어 | 의미 | 형식 예 |
|--------|------|---------|
| REQ | 기능 요구사항 | REQ-FUNC-MEAS-001 |
| NFR | 비기능 요구사항 | REQ-NFR-PERF-001 |
| SAF | 안전 요구사항 | REQ-SAFE-001 |
| REG | 규정 요구사항 | REQ-REG-001 |
| ITF | 인터페이스 요구사항 | REQ-IF-001 |
| DES | 설계 문서 | DES-PROTO-MEAS, DES-DB-03 |
| IMP | 구현 코드 | IMP-BE-MEAS, IMP-RUST-DIFF |
| VV | 검증 테스트 | VV-UT-MEAS-001, VV-E2E-FLOW |

### 11.3 요구사항 → 아키텍처 매핑 (요약)

| 요구사항 그룹 | SAD 소프트웨어 항목 | Proto 서비스 | DB 스키마 |
|-------------|-------------------|------------|----------|
| REQ-FUNC-AUTH-* | auth-service | AuthService (5 RPC) | 01-auth.sql |
| REQ-FUNC-USER-* | user-service | UserService (5 RPC) | 01-auth.sql (users) |
| REQ-FUNC-DEV-* | device-service, Rust Core BLE | DeviceService (5 RPC) | 02-device.sql |
| REQ-FUNC-CART-* | cartridge-service, Rust Core NFC | CartridgeService (5 RPC) | 09-cartridge.sql |
| REQ-FUNC-MEAS-* | measurement-service, Rust Core | MeasurementService (6 RPC) | 03-measurement.sql |
| REQ-FUNC-AI-* | ai-inference-service, Rust Core AI | AiInferenceService (5 RPC) | 08-ai-inference.sql |
| REQ-FUNC-CAL-* | calibration-service | CalibrationService (5 RPC) | 10-calibration.sql |
| REQ-FUNC-SUB-* | subscription-service | SubscriptionService (5 RPC) | 05-subscription.sql |
| REQ-FUNC-PAY-* | payment-service | PaymentService (5 RPC) | 07-payment.sql |
| REQ-FUNC-SHOP-* | shop-service | ShopService (5 RPC) | 06-shop.sql |
| REQ-FUNC-COACH-* | coaching-service | CoachingService (5 RPC) | 11-coaching.sql |
| REQ-FUNC-TELE-* | telemedicine-service | TelemedicineService (7 RPC) | 15-telemedicine.sql |
| REQ-FUNC-RX-* | prescription-service | PrescriptionService (5 RPC) | 19-prescription.sql |
| REQ-FUNC-HR-* | health-record-service | HealthRecordService (5 RPC) | 14-health-record.sql |
| REQ-FUNC-FAM-* | family-service | FamilyService (5 RPC) | 13-family.sql |
| REQ-FUNC-NOTI-* | notification-service | NotificationService (5 RPC) | 12-notification.sql |
| REQ-FUNC-COMM-* | community-service | CommunityService (5 RPC) | 17-community.sql |
| REQ-FUNC-RESV-* | reservation-service | ReservationService (5 RPC) | 16-reservation.sql |
| REQ-FUNC-ADMIN-* | admin-service | AdminService (15 RPC) | 18-admin.sql |
| REQ-FUNC-OFF-* | Rust Core (sync, ai) | — | SQLite CRDT (로컬) |

### 11.4 추적성 매트릭스: REQ → 아키텍처 항목 → 테스트 케이스

| REQ ID (대표) | 아키텍처 항목 (SAD) | 구현 경로 (IMP) | 테스트 케이스 (V&V) |
|--------------|---------------------|-----------------|---------------------|
| REQ-FUNC-AUTH-001~007 | auth-service, AuthService Proto | backend/services/auth-service/ | VV-UT-AUTH-001 (service_test.go), VV-E2E-FLOW (flow_test.go) |
| REQ-FUNC-USER-001~003 | user-service, UserService Proto | backend/services/user-service/ | VV-UT-USER-001, VV-E2E-FLOW |
| REQ-FUNC-DEV-001~006 | device-service, Rust ble/ | backend/services/device-service/, rust-core/manpasik-engine/ | VV-UT-DEVICE-001, VV-E2E-HEALTH |
| REQ-FUNC-CART-001~006 | cartridge-service, Rust nfc/ | backend/services/cartridge-service/, rust-core/ | VV-UT-CART-001, VV-E2E-CARTRIDGE |
| REQ-FUNC-MEAS-001~009 | measurement-service, Rust differential/, fingerprint/ | backend/services/measurement-service/, rust-core/ | VV-UT-MEAS-001, VV-UT-DIFF-001, VV-E2E-FLOW |
| REQ-FUNC-AI-001~005 | ai-inference-service, Rust ai/ | backend/services/ai-inference-service/, rust-core/ | VV-UT-AI-INFER, VV-E2E-AI_HARDWARE |
| REQ-FUNC-CAL-001~004 | calibration-service | backend/services/calibration-service/ | VV-UT-CALIB-001 |
| REQ-FUNC-SUB-001~004 | subscription-service | backend/services/subscription-service/ | VV-UT-SUB-001, VV-E2E-COMMERCE |
| REQ-FUNC-PAY-001~004 | payment-service | backend/services/payment-service/ | VV-UT-PAY-001, VV-E2E-COMMERCE |
| REQ-FUNC-SHOP-001~003 | shop-service | backend/services/shop-service/ | VV-UT-SHOP-001 |
| REQ-FUNC-COACH-001~004 | coaching-service | backend/services/coaching-service/ | VV-UT-COACH-001, VV-E2E-COACHING |
| REQ-FUNC-TELE-001~004 | telemedicine-service | backend/services/telemedicine-service/ | VV-UT-TELE-001, VV-E2E-MEDICAL |
| REQ-FUNC-ADMIN-001~004 | admin-service | backend/services/admin-service/ | VV-UT-ADMIN-001, VV-E2E-ADMIN |
| REQ-FUNC-OFF-001~003 | Rust sync/, ai/ | rust-core/manpasik-engine/ | VV-UT-SYNC-001, VV-QA-OFFLINE-72H |
| REQ-NFR-PERF-001~012 | 전체 서비스 | — | VV-PERF-LOAD, criterion 벤치마크 |
| REQ-NFR-SEC-001~011 | gateway, auth, shared/middleware | infrastructure/, backend/shared/ | VV-INFRA-TLS, gosec, Trivy |
| REQ-SAFE-001~010 | 위험관리 기반 | — | VV-UT-* (도메인별), VV-E2E-*, 보안 침투 테스트 |
| REQ-DATA-001~010 | PostgreSQL, TimescaleDB, Milvus, Redis, Kafka, ES, MinIO | infrastructure/database/init/ | VV-IT-DB (DB 마이그레이션 테스트) |
| REQ-IF-001~014 | BLE, NFC, gRPC, Kafka, REST | Proto, Rust ble/nfc, Kong | VV-UT-BLE, VV-UT-NFC, VV-IT-GRPC |

**전체 80개 REQ 상세 매핑**: `docs/plan/plan-traceability-matrix.md` (v2.0)

---

## 부록

### 부록 A: 요구사항 요약 통계

| 카테고리 | 요구사항 수 | Class B | Class A |
|---------|-----------|---------|---------|
| 기능 요구사항 (REQ-FUNC) | 80+ | 42 | 38+ |
| 비기능 요구사항 (REQ-NFR) | 35 | — | — |
| 안전 요구사항 (REQ-SAFE) | 10 | 10 | — |
| 규정 요구사항 (REQ-REG) | 22 | — | — |
| 인터페이스 요구사항 (REQ-IF) | 14+ | — | — |
| 데이터 요구사항 (REQ-DATA) | 10 | — | — |
| **합계** | **171+** | — | — |

### 부록 B: 우선순위 정의

| 우선순위 | 의미 | Phase |
|---------|------|-------|
| P0 | 필수 — MVP 출시 조건 | Phase 1 |
| P1 | 중요 — Core 기능 | Phase 2 |
| P2 | 높음 — Advanced 기능 | Phase 2~3 |
| P3 | 보통 — 확장 기능 | Phase 3 |
| P4 | 낮음 — 미래 기능 | Phase 4~5 |

### 부록 C: IEC 62304 §5.2 체크리스트

| §5.2 요구사항 | 대응 섹션 | 상태 |
|--------------|----------|------|
| 5.2.1 요구사항 활동의 범위 정의 | §1 개요 | ✅ |
| 5.2.2 기능 요구사항 정의 | §4 기능 요구사항 (80+ REQ) | ✅ |
| 5.2.3 비기능 요구사항 정의 | §5 비기능 요구사항 (35 REQ-NFR) | ✅ |
| 5.2.4 인터페이스 요구사항 정의 | §6 인터페이스 요구사항 (14+ REQ-IF) | ✅ |
| 5.2.5 안전 관련 요구사항 정의 | §8 안전 요구사항 (10 REQ-SAFE) | ✅ |
| 5.2.6 규정 요구사항 포함 | §9 규정 요구사항 (22 REQ-REG) | ✅ |
| 5.2.7 추적성 확립 | §11 추적성 | ✅ |
| 위험 통제 조치 식별 | §8.1, §8.2 (STRIDE 31개 위협) | ✅ |
| 요구사항 검토 및 업데이트 | V&V 마스터 플랜 참조 | ✅ |

---

**마지막 업데이트**: 2026-02-13 (v3.0)  
**다음 검토**: Phase 완료 시 또는 분기 1회  
**승인 대기**: 기술 리더 검토 → 품질 관리자 검토 → 승인
