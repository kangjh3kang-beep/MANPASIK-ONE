# ManPaSik 미구현 사항 및 구현 계획서

> **문서번호**: MPK-IMPL-PLAN-2026-02-11
> **작성자**: Claude Opus 4.5
> **목적**: 기획안 대비 미구현 사항 식별 및 Phase별 세부 구현 계획 수립
> **공유 대상**: 모든 IDE·AI 에이전트 (Cursor, VS Code, Claude, ChatGPT 등)

---

## 목차

1. [분석 요약](#1-분석-요약)
2. [기능 요구사항 총괄 (80개)](#2-기능-요구사항-총괄-80개)
3. [구현 현황 분석](#3-구현-현황-분석)
4. [미구현 사항 목록](#4-미구현-사항-목록)
5. [Phase별 세부 구현 계획](#5-phase별-세부-구현-계획)
6. [우선순위 매트릭스](#6-우선순위-매트릭스)
7. [의존성 맵](#7-의존성-맵)
8. [일정 추정](#8-일정-추정)

---

## 1. 분석 요약

### 1.1 전체 현황

| 영역 | 총 항목 | 구현 완료 | 미구현 | 진행률 |
|------|---------|----------|--------|--------|
| **기능 요구사항** | 80개 | 35개 | 45개 | 44% |
| **백엔드 서비스** | 29개 | 23개 | 6개 | 79% |
| **Rust 코어 모듈** | 9개 | 6개 | 3개 (부분) | 85% |
| **Flutter Feature** | 12개 | 6개 | 6개 | 50% |
| **규정 문서** | 15개 | 8개 | 7개 | 53% |

### 1.2 Phase별 구현 현황

| Phase | 기능 수 | 구현 | 미구현 | 상태 |
|-------|---------|------|--------|------|
| Phase 1 (MVP) | 18개 | 16개 | 2개 | ✅ 89% |
| Phase 2 (Core) | 35개 | 15개 | 20개 | 🔄 43% |
| Phase 3 (Advanced) | 16개 | 4개 | 12개 | 🔲 25% |
| Phase 4 (Ecosystem) | 13개 | 0개 | 13개 | 🔲 0% |
| Phase 5 (Future) | 2개 | 0개 | 2개 | 🔲 0% |

---

## 2. 기능 요구사항 총괄 (80개)

### 2.1 Phase 1 (MVP) - 18개 기능

| ID | 기능명 | 서비스 | 상태 | 비고 |
|----|--------|--------|------|------|
| REQ-001 | 사용자 회원가입/로그인 | auth-service | ✅ | JWT, Keycloak |
| REQ-002 | 사용자 프로필 관리 | user-service | ✅ | |
| REQ-003 | 다중 리더기 등록/관리 | device-service | ✅ | |
| REQ-004 | 리더기 펌웨어 OTA | device-service | ⚠️ | 인터페이스만 |
| REQ-005 | 카트리지 자동인식 | Rust nfc | ⚠️ | 스텁만 구현 |
| REQ-006 | 88차원 차동측정 | Rust differential | ✅ | 100% 완료 |
| REQ-007 | 측정 결과 저장/시각화 | measurement-service | ✅ | |
| REQ-008 | 896차원 핑거프린트 | Rust fingerprint | ✅ | 구조 완성 |
| REQ-009 | 비표적 분석 | measurement-service | ⚠️ | 기본만 |
| REQ-010 | 측정 세션 관리 | measurement-service | ✅ | |
| REQ-011 | 오프라인 완전 구동 | Rust sync | ✅ | CRDT 구현 |
| REQ-057 | 규제 준수 관리 | 문서 | ⚠️ | 초안 |
| REQ-058 | TPM 보안 칩 | Rust crypto | ✅ | |
| REQ-059 | BLE AES-CCM 암호화 | Rust ble | ⚠️ | 구조만 |
| REQ-060 | HTTPS TLS 1.3 | 인프라 | ✅ | |
| REQ-063 | 감사 추적 로그 | 전역 | ✅ | |
| REQ-065 | 72시간 오프라인 검증 | QA | 🔲 | 미실행 |
| REQ-073 | 콘텐츠 캐싱/동기화 | Rust sync | ✅ | |

### 2.2 Phase 2 (Core) - 35개 기능

| ID | 기능명 | 서비스 | 상태 | 비고 |
|----|--------|--------|------|------|
| REQ-012 | 구독 등급 관리 | subscription-service | ✅ | 4티어 |
| REQ-013 | SaaS 구독 결제 | payment-service | ✅ | |
| REQ-014 | 카트리지 무한확장 레지스트리 | cartridge-service | ✅ | 65,536종 |
| REQ-015 | 등급별 카트리지 접근 제어 | cartridge-service | ✅ | |
| REQ-016 | 카트리지 사용량 추적 | cartridge-service | ✅ | |
| REQ-017 | 온라인 상품 판매 | shop-service | ✅ | |
| REQ-018 | 장바구니/주문 관리 | shop-service | ✅ | |
| REQ-019 | AI 실시간 추론 | ai-inference-service | ✅ | 시뮬레이션 |
| REQ-020 | AI 건강 코칭 | coaching-service | ✅ | |
| REQ-021 | 음식 사진→칼로리 | vision-service | 🔲 | **미구현** |
| REQ-023 | 카트리지 보정 데이터 | calibration-service | ✅ | |
| REQ-024 | 데이터 허브/타임라인 | Flutter data_hub | 🔲 | **미구현** |
| REQ-025 | 외부 건강 앱 연동 | measurement-service | 🔲 | **미구현** |
| REQ-026 | 공공데이터 연계 | ai-inference-service | 🔲 | **미구현** |
| REQ-027 | 데이터 내보내기 | measurement-service | ⚠️ | FHIR만 |
| REQ-051 | 448차원 분석 | Rust fingerprint | ✅ | |
| REQ-061 | MFA 다중 인증 | auth-service | 🔲 | **미구현** |
| REQ-062 | RBAC 역할 기반 접근 | admin-service | ✅ | |
| REQ-064 | 동적 부하 시뮬레이션 | 인프라 | 🔲 | **미실행** |
| REQ-066 | 개인 기준선(My Zone) | coaching-service | 🔲 | **미구현** |
| REQ-067 | 환경 모니터링 표시 | measurement-service | 🔲 | **미구현** |
| REQ-068 | 위험 예측/경고 | ai-inference-service | 🔲 | **미구현** |
| REQ-070 | AI 학습 히스토리 | ai-training-service | 🔲 | **미구현** |
| REQ-071 | 정기 구독 서비스 | shop-service | 🔲 | **미구현** |
| REQ-072 | 위시리스트 | shop-service | 🔲 | **미구현** |
| REQ-074 | 다국어 UI | Flutter l10n | ✅ | 6언어 |
| REQ-075 | 접근성(다크모드) | Flutter theme | ✅ | |
| REQ-076 | 리더기 위치별 대시보드 | device-service | 🔲 | **미구현** |
| REQ-077 | 구독별 리더기 제한 | subscription-service | ✅ | |
| REQ-079 | 신규 카트리지 OTA 배포 | cartridge-service | 🔲 | **미구현** |
| REQ-080 | 소셜 로그인 | auth-service | 🔲 | **미구현** |

### 2.3 Phase 3 (Advanced) - 16개 기능

| ID | 기능명 | 서비스 | 상태 | 비고 |
|----|--------|--------|------|------|
| REQ-022 | 운동 영상 분석 | vision-service | 🔲 | **미구현** |
| REQ-028 | 화상진료 | telemedicine-service | ✅ | 구조 완료 |
| REQ-029 | 병원/약국 검색/예약 | reservation-service | ✅ | |
| REQ-030 | 처방전 관리 | prescription-service | ✅ | |
| REQ-031 | 가족 그룹 관리 | family-service | ✅ | |
| REQ-032 | 보호자 모니터링 | family-service | 🔲 | **미구현** |
| REQ-033 | 가족 건강 리포트 | health-record-service | 🔲 | **미구현** |
| REQ-034 | 건강 기록 (FHIR) | health-record-service | ✅ | R4 호환 |
| REQ-035 | 커뮤니티 포럼 | community-service | ✅ | |
| REQ-036 | 전문가 Q&A | community-service | 🔲 | **미구현** |
| REQ-037 | 글로벌 챌린지 | community-service | 🔲 | **미구현** |
| REQ-038 | 실시간 번역 | translation-service | ✅ | |
| REQ-039 | 푸시/이메일/SMS 알림 | notification-service | ✅ | |
| REQ-040 | 계층형 관리자 포탈 | admin-service | ✅ | |
| REQ-069 | 긴급 연락망 설정 | notification-service | 🔲 | **미구현** |
| REQ-078 | 정기 결과 내보내기 | measurement-service | 🔲 | **미구현** |

### 2.4 Phase 4 (Ecosystem) - 13개 기능

| ID | 기능명 | 서비스 | 상태 | 비고 |
|----|--------|--------|------|------|
| REQ-041 | 재고 관리 | inventory-service | 🔲 | **서비스 미구현** |
| REQ-042 | 배송 추적 | logistics-service | 🔲 | **서비스 미구현** |
| REQ-043 | 비즈니스 인텔리전스 | analytics-service | 🔲 | **서비스 미구현** |
| REQ-044 | SDK 마켓 | marketplace-service | 🔲 | **서비스 미구현** |
| REQ-045 | 수익 분배 시스템 | marketplace-service | 🔲 | **서비스 미구현** |
| REQ-046 | AI 모델 학습 | ai-training-service | 🔲 | **서비스 미구현** |
| REQ-047 | NLP 자연어 처리 | nlp-service | 🔲 | **서비스 미구현** |
| REQ-048 | IoT 게이트웨이 | iot-gateway-service | 🔲 | **서비스 미구현** |
| REQ-049 | 리더기 위치 추적 | location-service | 🔲 | **서비스 미구현** |
| REQ-050 | 긴급 대응 시스템 | emergency-service | 🔲 | **서비스 미구현** |
| REQ-052 | 전자코/전자혀 통합 | Rust measurement | 🔲 | **미구현** |
| REQ-054 | 음성 명령 인터페이스 | nlp-service | 🔲 | **서비스 미구현** |

### 2.5 Phase 5 (Future) - 2개 기능

| ID | 기능명 | 서비스 | 상태 | 비고 |
|----|--------|--------|------|------|
| REQ-053 | 1792차원 궁극 분석 | Rust fingerprint | 🔲 | 구조만 준비 |
| REQ-055 | 웨어러블 디바이스 | iot-gateway-service | 🔲 | **미구현** |
| REQ-056 | 스마트홈 통합 | iot-gateway-service | 🔲 | **미구현** |

---

## 3. 구현 현황 분석

### 3.1 백엔드 서비스 (29개)

#### ✅ 구현 완료 (23개)
```
auth-service (50051)        user-service (50052)
device-service (50053)      measurement-service (50054)
subscription-service (50055) shop-service (50056)
payment-service (50057)     ai-inference-service (50058)
cartridge-service (50059)   calibration-service (50060)
coaching-service (50061)    admin-service
community-service           family-service
health-record-service       notification-service
prescription-service        reservation-service
telemedicine-service        translation-service
video-service
```

#### 🔲 미구현 (6개)
| 서비스 | Phase | 우선순위 | 예상 기능 |
|--------|-------|----------|----------|
| **analytics-service** | 4 | P1 | BI 대시보드, 매출 분석 |
| **emergency-service** | 4 | P0 | 119 연동, 위치 공유 |
| **iot-gateway-service** | 4 | P1 | MQTT, 다중 디바이스 |
| **marketplace-service** | 4 | P0 | SDK 마켓, 수익 분배 |
| **nlp-service** | 4 | P1 | TTS/STT, 음성 명령 |
| **vision-service** | 2 | P1 | 음식 인식, 운동 분석 |

### 3.2 Rust 코어 (9개 모듈)

#### ✅ 완성 (6개)
- differential (100%) - 차동측정 알고리즘
- fingerprint (95%) - 88→1792차원 성장
- crypto (95%) - AES-256, 해시체인
- sync (85%) - CRDT 오프라인
- dsp (90%) - FFT, 필터링
- flutter-bridge (95%) - FFI 래퍼

#### ⚠️ 부분 구현 (3개)
| 모듈 | 구현율 | 미완성 항목 | 우선순위 |
|------|--------|-------------|----------|
| **ai** | 50% | TFLite 실제 추론 | P0 |
| **ble** | 50% | btleplug GATT 통신 | P0 |
| **nfc** | 60% | ISO 14443A 실제 읽기 | P0 |

### 3.3 Flutter 앱 (12개 Feature)

#### ✅ 구현 완료 (6개)
- auth (로그인/회원가입)
- home (대시보드)
- measurement (측정)
- devices (기기 관리)
- settings (설정)
- user (프로필)

#### 🔲 미구현 (6개)
| Feature | Phase | 우선순위 | 예상 화면 |
|---------|-------|----------|----------|
| **data_hub** | 2 | P0 | 타임라인, 트렌드 차트 |
| **ai_coach** | 2 | P0 | AI 대화, 코칭 카드 |
| **market** | 2 | P1 | 상품 목록, 장바구니 |
| **community** | 3 | P1 | 포럼, Q&A |
| **medical** | 3 | P0 | 화상진료, 예약 |
| **family** | 3 | P1 | 가족 관리 |

---

## 4. 미구현 사항 목록

### 4.1 P0 (Critical) - 즉시 구현 필요

| # | 영역 | 항목 | Phase | 예상 공수 |
|---|------|------|-------|----------|
| 1 | Rust | AI 모듈 TFLite 실제 추론 | 1 | 3일 |
| 2 | Rust | BLE 모듈 btleplug 통신 | 1 | 5일 |
| 3 | Rust | NFC 모듈 실제 읽기 | 1 | 3일 |
| 4 | Flutter | data_hub Feature | 2 | 5일 |
| 5 | Flutter | ai_coach Feature | 2 | 5일 |
| 6 | Backend | vision-service | 2 | 7일 |
| 7 | Backend | emergency-service | 4 | 5일 |
| 8 | Backend | marketplace-service | 4 | 7일 |
| 9 | 규정 | IEC 62304 정식 문서 3종 | 1 | 10일 |
| 10 | 규정 | ISO 14971 위험관리 5종 | 1 | 10일 |

### 4.2 P1 (High) - 다음 Sprint

| # | 영역 | 항목 | Phase | 예상 공수 |
|---|------|------|-------|----------|
| 1 | Flutter | market Feature | 2 | 5일 |
| 2 | Flutter | medical Feature | 3 | 7일 |
| 3 | Backend | analytics-service | 4 | 5일 |
| 4 | Backend | iot-gateway-service | 4 | 7일 |
| 5 | Backend | nlp-service | 4 | 7일 |
| 6 | 기능 | MFA 다중 인증 | 2 | 3일 |
| 7 | 기능 | 소셜 로그인 | 2 | 3일 |
| 8 | 기능 | 위험 예측/경고 | 2 | 5일 |
| 9 | 기능 | 개인 기준선(My Zone) | 2 | 3일 |

### 4.3 P2 (Medium) - 후순위

| # | 영역 | 항목 | Phase | 예상 공수 |
|---|------|------|-------|----------|
| 1 | Flutter | community Feature | 3 | 5일 |
| 2 | Flutter | family Feature | 3 | 5일 |
| 3 | 기능 | 운동 영상 분석 | 3 | 7일 |
| 4 | 기능 | 글로벌 챌린지 | 3 | 5일 |
| 5 | 기능 | 전문가 Q&A | 3 | 5일 |
| 6 | 기능 | 공공데이터 연계 | 2 | 5일 |
| 7 | 기능 | HealthKit/Google Health | 2 | 5일 |

---

## 5. Phase별 세부 구현 계획

### 5.1 Phase 1 잔여 작업 (2주)

#### Week 1: Rust 코어 완성
```
Day 1-2: AI 모듈 TFLite 실제 구현
  - tflitec::Interpreter 초기화
  - 모델 파일 로드 (tflite-models/)
  - 입력 텐서 바인딩, 추론 실행
  - 테스트: predict() 실제 동작

Day 3-4: BLE 모듈 btleplug 구현
  - btleplug Manager 초기화
  - 디바이스 스캔 (scan_devices)
  - 연결 및 GATT 서비스 조회
  - 특성 읽기/쓰기/Notification

Day 5: NFC 모듈 실제 구현
  - 플랫폼별 NFC 라이브러리 연동
  - ISO 14443A 프로토콜 구현
  - 카트리지 UID/데이터 읽기
```

#### Week 2: 규정 문서 정식화
```
Day 1-3: IEC 62304 문서
  - Software Development Plan (SDP)
  - Software Requirements Specification (SRS)
  - Software Architecture Document (SAD)

Day 4-5: ISO 14971 문서
  - FMEA 분석 보고서
  - 위험 추정/평가 보고서
```

### 5.2 Phase 2 구현 계획 (8주)

#### Sprint 1 (Week 1-2): 핵심 Flutter Feature
```
data_hub Feature:
  - 타임라인 화면 (측정 기록 시간순)
  - 트렌드 차트 (fl_chart)
  - 데이터 필터링 (기간/카테고리)
  - gRPC 연동 (measurement-service)

ai_coach Feature:
  - AI 대화 화면 (채팅 UI)
  - 코칭 카드 (일일/주간)
  - 목표 설정 화면
  - gRPC 연동 (coaching-service)
```

#### Sprint 2 (Week 3-4): vision-service 구현
```
서비스 구조:
  - cmd/main.go
  - internal/handler/grpc.go
  - internal/service/vision.go
  - internal/repository/memory/

gRPC 메서드:
  - AnalyzeFood(image) → calories, nutrients
  - AnalyzeExercise(video) → calories_burned
  - DetectObjects(image) → objects
  - GetModelInfo() → model_version
```

#### Sprint 3 (Week 5-6): 부가 기능
```
MFA 다중 인증:
  - auth-service TOTP 구현
  - Flutter OTP 입력 화면

소셜 로그인:
  - Google OAuth 연동
  - Apple Sign In
  - Facebook Login

위험 예측/경고:
  - ai-inference-service 확장
  - 푸시 알림 연동
```

#### Sprint 4 (Week 7-8): Market Feature
```
Flutter market Feature:
  - 상품 목록 화면
  - 상품 상세 화면
  - 장바구니 화면
  - 주문/결제 화면
  - gRPC 연동 (shop-service, payment-service)
```

### 5.3 Phase 3 구현 계획 (12주)

#### Sprint 1-2: Medical Feature
```
Flutter medical Feature:
  - 화상진료 예약 화면
  - WebRTC 화상 통화 화면
  - 병원/약국 검색 화면
  - 예약 관리 화면
  - gRPC 연동 (telemedicine, reservation, prescription)
```

#### Sprint 3-4: Family/Community Feature
```
Flutter family Feature:
  - 가족 그룹 관리
  - 구성원 초대
  - 보호자 모니터링 대시보드

Flutter community Feature:
  - 포럼 목록/상세
  - 게시글 작성
  - 댓글
  - 챌린지 참여
```

#### Sprint 5-6: 고급 기능
```
전문가 Q&A:
  - community-service 확장
  - 전문가 인증 시스템

글로벌 챌린지:
  - 챌린지 생성/참여
  - 순위표
  - 보상 시스템
```

### 5.4 Phase 4 구현 계획 (24주)

#### Sprint 1-4: 핵심 서비스
```
emergency-service:
  - 위험 감지 로직
  - 119 API 연동
  - 위치 공유
  - 음성 통화

marketplace-service:
  - SDK 등록/승인
  - 카트리지 마켓
  - 수익 분배 (30:70)
```

#### Sprint 5-8: AI/IoT 서비스
```
ai-training-service:
  - Flower 연합학습
  - 모델 버전 관리
  - 배포 파이프라인

iot-gateway-service:
  - MQTT 브로커
  - 다중 디바이스
  - 프로토콜 변환

nlp-service:
  - TTS/STT
  - 음성 명령 처리
  - 자연어 이해
```

#### Sprint 9-12: 분석/물류
```
analytics-service:
  - BI 대시보드
  - 매출 분석
  - 사용자 행동 분석

inventory-service / logistics-service:
  - 재고 관리
  - 발주 자동화
  - 배송 추적
```

---

## 6. 우선순위 매트릭스

```
                    높은 영향
                       ▲
                       │
     P0 Critical       │       P1 High
     ┌─────────────────┼─────────────────┐
     │ • Rust AI/BLE/NFC│ • MFA 인증      │
     │ • IEC 62304 문서 │ • 소셜 로그인   │
     │ • data_hub      │ • vision-service│
     │ • ai_coach      │ • market Feature│
     │ • emergency-svc │ • analytics-svc │
높은 ├─────────────────┼─────────────────┤ 낮은
긴급 │ P2 Medium       │ P3 Low          │ 긴급
     │ • community     │ • 웨어러블      │
     │ • family        │ • 스마트홈      │
     │ • 전문가 Q&A    │ • 1792차원      │
     │ • 글로벌 챌린지 │                 │
     └─────────────────┼─────────────────┘
                       │
                       ▼
                    낮은 영향
```

---

## 7. 의존성 맵

```
Phase 1 (완료)
    ├── auth-service ─────────────────────────────┐
    ├── measurement-service ──────────────────────┤
    ├── device-service ───────────────────────────┤
    └── Rust Core (differential, fingerprint) ────┤
                                                  │
Phase 2 (진행중)                                   ▼
    ├── subscription-service ◄── payment-service
    ├── shop-service ◄── payment-service
    ├── ai-inference-service ◄── Rust AI 완성 필요
    ├── vision-service (신규) ◄── AI 모델 필요
    ├── data_hub Feature ◄── measurement-service
    ├── ai_coach Feature ◄── coaching-service
    └── market Feature ◄── shop-service, payment-service
                                                  │
Phase 3                                           ▼
    ├── telemedicine-service ◄── video-service
    ├── reservation-service
    ├── family-service ◄── user-service
    ├── medical Feature ◄── telemedicine, reservation
    └── community Feature ◄── community-service
                                                  │
Phase 4                                           ▼
    ├── emergency-service ◄── notification, location
    ├── marketplace-service ◄── payment-service
    ├── ai-training-service ◄── ai-inference-service
    ├── iot-gateway-service
    └── nlp-service
```

---

## 8. 일정 추정

### 8.1 전체 로드맵

| Phase | 기간 | 주요 목표 |
|-------|------|----------|
| **Phase 1 잔여** | 2주 | Rust 코어 완성, 규정 문서 |
| **Phase 2** | 8주 | Core 기능, Flutter Feature 4개 |
| **Phase 3** | 12주 | Advanced 기능, 의료/커뮤니티 |
| **Phase 4** | 24주 | Ecosystem, AI/IoT |
| **Phase 5** | 24주 | Future, 1792차원 |

### 8.2 마일스톤

| 마일스톤 | 목표일 | 완료 조건 |
|----------|--------|----------|
| **M1: Rust 완성** | +2주 | AI/BLE/NFC 실제 구현 |
| **M2: Phase 2 MVP** | +10주 | data_hub, ai_coach, vision |
| **M3: MFDS 인허가** | +14주 | 규제 문서 제출 |
| **M4: Phase 3 완료** | +26주 | 의료/커뮤니티 |
| **M5: Phase 4 완료** | +50주 | Ecosystem |
| **M6: 글로벌 출시** | +70주 | FDA/CE 인증 |

---

## 9. 참조 문서

| 문서 | 경로 |
|------|------|
| 기획안 v1.1 | docs/plan/MPK-ECO-PLAN-v1.1-COMPLETE.md |
| MSA 로드맵 | docs/plan/msa-expansion-roadmap.md |
| 카트리지 시스템 | docs/specs/cartridge-system-spec.md |
| 품질 게이트 | QUALITY_GATES.md |
| 구현 현황 보고서 | docs/reports/implementation-status-2026-02-11.md |

---

**문서 종료**

*본 계획서는 모든 IDE 및 AI 에이전트가 참조할 수 있도록 작성되었습니다. 작업 진행 시 CHANGELOG.md에 기록하고, 완료 시 본 문서의 상태를 갱신해주세요.*
