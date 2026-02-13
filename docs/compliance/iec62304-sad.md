# IEC 62304:2015 소프트웨어 아키텍처 설계서 (SAD)

> **문서 ID**: DOC-SAD-001  
> **표준 근거**: IEC 62304:2006+AMD1:2015 §5.3 (소프트웨어 아키텍처 설계)  
> **버전**: v2.0  
> **안전 등급**: Class B (IEC 62304)  
> **작성일**: 2026-02-13  
> **상태**: 초안 (검토 대기)  
> **승인**: — (승인 대기)

---

## 변경 이력

| 버전 | 일자 | 작성자 | 변경 내용 |
|------|------|--------|----------|
| v1.0 | 2026-02-12 | Claude Agent | 초안 작성 (23개 서비스, 기본 아키텍처) |
| v2.0 | 2026-02-13 | Claude Agent | 전면 개정 — IEC 62304 §5.3 필수 항목 완결, 3-Tier·데이터·보안·배포 상세화, SOUP 확장 |

---

## 참조 문서

| 구분 | 문서 | 경로 |
|------|------|------|
| 상위 | 소프트웨어 개발 계획서 (SDP) | `docs/compliance/iec62304-sdp.md` |
| 요구사항 | 소프트웨어 요구사항 명세서 (SRS) | `docs/compliance/iec62304-srs.md` |
| 안전 판정 | 소프트웨어 안전 등급 판정서 | `docs/compliance/software-safety-classification.md` |
| 위험관리 | ISO 14971 위험관리 계획 | `docs/compliance/iso14971-risk-management-plan.md` |
| 검증 | V&V 마스터 플랜 | `docs/compliance/vnv-master-plan.md` |
| 이벤트 | Kafka 이벤트 스키마 명세 | `docs/specs/event-schema-specification.md` |
| Proto | gRPC API 정의 | `backend/shared/proto/manpasik.proto` |

---

## 목차

1. [개요](#1-개요-clause-531)
2. [시스템 아키텍처 개요](#2-시스템-아키텍처-개요-clause-532)
3. [소프트웨어 항목 분해표](#3-소프트웨어-항목-분해표-clause-533)
4. [인터페이스 설계](#4-인터페이스-설계-clause-534)
5. [데이터 아키텍처](#5-데이터-아키텍처-clause-535)
6. [보안 아키텍처](#6-보안-아키텍처-clause-536)
7. [배포 아키텍처](#7-배포-아키텍처-clause-537)
8. [안전 등급 매핑](#8-안전-등급-매핑-clause-538)
9. [외부 의존성 (SOUP)](#9-외부-의존성-soup-clause-539)
10. [부록](#10-부록)

---

## 1. 개요 (Clause 5.3.1)

### 1.1 목적

본 문서는 **ManPaSik (만파식)** 의료기기 소프트웨어의 아키텍처를 정의합니다. IEC 62304:2015 §5.3에 따라 소프트웨어 시스템의 구조, 소프트웨어 항목 분해, 인터페이스, 데이터 흐름, 보안, 배포를 명시합니다.

### 1.2 범위

| 항목 | 내용 |
|------|------|
| **제품명** | ManPaSik (만파식, 萬波息) — 차동측정 기반 범용 분석 헬스케어 AI 생태계 |
| **의료기기 등급** | Class II (체외진단의료기기, IVD) |
| **소프트웨어 안전 등급** | IEC 62304 Class B |
| **기술 스택** | Go gRPC MSA 30+서비스, Rust 코어엔진(BLE/NFC/DSP/AI), Flutter 모바일앱 |
| **인프라** | PostgreSQL, TimescaleDB, Milvus, Redis, Kafka, MinIO, Elasticsearch |
| **배포** | K8s 배포, Docker Compose 개발환경 |

### 1.3 아키텍처 원칙

1. **마이크로서비스 분리**: 각 서비스는 독립 배포·확장 가능
2. **오프라인 우선**: 코어 엔진은 네트워크 없이 100% 동작
3. **보안 계층화**: 심층 방어 (Defense in Depth)
4. **관찰 가능성**: 모든 서비스에 메트릭·로깅·트레이싱 적용

---

## 2. 시스템 아키텍처 개요 (Clause 5.3.2)

### 2.1 3-Tier 아키텍처

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│  Tier 1: 클라이언트 계층                                                          │
│  ┌─────────────────────────────────────────────────────────────────────────────┐ │
│  │ Flutter 모바일 앱 (iOS/Android)                                              │ │
│  │ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────────────┐ │ │
│  │ │ UI (Material) │ │ Rust Core    │ │ gRPC Client   │ │ REST Client (Dio)    │ │ │
│  │ │ 7+ 화면       │ │ (FFI Bridge) │ │ (Protobuf)    │ │ (파일 업로드 등)      │ │ │
│  │ └──────────────┘ └──────────────┘ └───────┬───────┘ └──────────┬───────────┘ │ │
│  └──────────────────────────────────────────┼────────────────────┼─────────────┘ │
├──────────────────────────────────────────────┼────────────────────┼───────────────┤
│  Tier 2: API Gateway 계층                    │                    │               │
│  ┌──────────────────────────────────────────┼────────────────────┼─────────────┐ │
│  │ Kong 3.7 (REST/gRPC 라우팅) + Keycloak 25.0 (OIDC 인증)                    │ │
│  │ JWT 검증 · RBAC · Rate Limiting · SSL 종단 · CORS                            │ │
│  └──────────────────────────────────────────┼────────────────────┼─────────────┘ │
├──────────────────────────────────────────────┼────────────────────┼───────────────┤
│  Tier 3: 마이크로서비스 계층 (Go 1.22+, gRPC)│                    │               │
│  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐             │
│  │ auth   │ │ user   │ │ device │ │measure-│ │payment │ │shop    │             │
│  │service │ │service │ │service │ │ment-svc│ │service │ │service │  … 30+      │
│  └───┬────┘ └───┬────┘ └───┬────┘ └───┬────┘ └───┬────┘ └───┬────┘             │
│      └──────────┴──────────┴─────────┴──────────┴──────────┴───────────────────│
│                                    │                                          │
│  ┌─────────────────────────────────────────────────────────────────────────┐  │
│  │ 메시지 버스: Kafka (Redpanda v24.2) — 이벤트 기반 비동기 통신             │  │
│  │ 18개 토픽: measurement.completed, payment.completed, subscription.changed 등 │
│  └─────────────────────────────────────────────────────────────────────────┘  │
├────────────────────────────────────────────────────────────────────────────────┤
│  데이터 계층                                                                   │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌────────┐ │
│  │PostgreSQL│ │ Redis    │ │ Milvus   │ │Elastic-  │ │ MinIO    │ │ Kafka  │ │
│  │ 16       │ │ 7        │ │ 2.4      │ │search 8  │ │          │ │(Redpanda)│
│  │TimescaleDB│ │(캐시)    │ │(벡터 DB) │ │(검색)    │ │(객체)    │ │(스트림)│
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘ └────────┘ │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 메시지 버스 (Kafka) 역할

| 역할 | 설명 |
|------|------|
| **이벤트 발행** | 서비스가 도메인 이벤트를 토픽에 발행 (예: 측정 완료, 결제 완료) |
| **비동기 처리** | 발행-구독 패턴으로 서비스 간 느슨한 결합 유지 |
| **데이터 파이프라인** | measurement → ai-inference, coaching, notification 등 다중 소비 |
| **DLQ** | 실패 메시지는 `manpasik.dlq` 토픽으로 전달 후 수동 재처리 |

---

## 3. 소프트웨어 항목 분해표 (Clause 5.3.3)

### 3.1 Gateway

| 항목 | 내용 |
|------|------|
| **소프트웨어 항목** | Gateway |
| **기술** | Kong 3.7 + Keycloak 25.0 (또는 Go 커스텀 Gateway) |
| **역할** | REST/gRPC 라우팅, JWT 인증 검증, RBAC, Rate Limiting, 파일 업로드 프록시 |
| **포트** | 50050 (gRPC), 8080 (REST) |
| **안전 등급** | Class A |

### 3.2 핵심 서비스 (Auth, User, Device, Measurement)

| # | 서비스 | 역할 | 포트 | 안전 등급 |
|---|--------|------|------|----------|
| 1 | auth-service | 사용자 인증, JWT 발급, OIDC 연동 | 50051 | Class B |
| 2 | user-service | 사용자 프로필, 구독 상태 관리 | 50052 | Class B |
| 3 | device-service | 리더기 등록/관리, OTA, 상태 모니터링 | 50053 | Class B |
| 4 | measurement-service | 측정 세션, 결과 저장, 벡터 검색, Milvus 연동 | 50054 | Class B |

### 3.3 상거래 서비스 (Payment, Subscription, Shop)

| # | 서비스 | 역할 | 포트 | 안전 등급 |
|---|--------|------|------|----------|
| 5 | payment-service | PG 결제 (Toss), 결제 완료/취소 | 50057 | Class A |
| 6 | subscription-service | 구독 관리 (Free/Basic/Clinical Guard 등 4티어) | 50055 | Class A |
| 7 | shop-service | 상품, 장바구니, 주문 관리 | 50056 | Class A |

### 3.4 AI·코칭 서비스

| # | 서비스 | 역할 | 포트 | 안전 등급 |
|---|--------|------|------|----------|
| 8 | ai-inference-service | AI 추론, 모델 서빙, 클라우드 분석 | 50058 | Class B |
| 9 | coaching-service | AI 건강 코칭, 추천 리포트 | 50061 | Class A |

### 3.5 카트리지·보정 서비스

| # | 서비스 | 역할 | 포트 | 안전 등급 |
|---|--------|------|------|----------|
| 10 | cartridge-service | 카트리지 인증(NFC), 추적, 29종 지원 | 50059 | Class B |
| 11 | calibration-service | 디바이스 보정, 기준값 관리 | 50060 | Class B |

### 3.6 원격진료·미디어·번역 서비스

| # | 서비스 | 역할 | 포트 | 안전 등급 |
|---|--------|------|------|----------|
| 12 | telemedicine-service | 화상진료, 예약 연동 | 50066 | Class B |
| 13 | video-service | 비디오 관리, 스트리밍 | 50064 | Class A |
| 14 | translation-service | 다국어 번역 | 50065 | Class A |

### 3.7 기타 서비스

| # | 서비스 | 역할 | 포트 | 안전 등급 |
|---|--------|------|------|----------|
| 15 | notification-service | 푸시 알림 (FCM), 이메일, SMS | 50062 | Class A |
| 16 | reservation-service | 예약 관리 (병원/의원) | 50067 | Class A |
| 17 | prescription-service | 처방 관리, FHIR R4 | 50069 | Class B |
| 18 | community-service | 커뮤니티 게시판 | 50063 | Class A |
| 19 | family-service | 가족 구성원 관리 | 50070 | Class A |
| 20 | health-record-service | 건강 기록, FHIR 매핑 | 50071 | Class B |
| 21 | vision-service | 음식 사진 분석 | 50072 | Class B |
| 22 | admin-service | 시스템 관리, 설정, 감사 로그 | 50068 | Class A |

### 3.8 Rust 코어엔진 (비서비스 컴포넌트)

| 모듈 | 기능 | 안전 등급 |
|------|------|----------|
| **differential** | 차동측정 연산 \( S_{corrected} = S_{det} - \alpha \times S_{ref} \) | Class B |
| **fingerprint** | 핑거프린트 벡터 생성 (88/448/896차원) | Class B |
| **dsp** | FFT, 디지털 필터 (신호 전처리) | Class B |
| **ble** | BLE 5.0 GATT 통신 (리더기 ↔ 앱) | Class B |
| **nfc** | NFC 카트리지 인식/검증 | Class B |
| **crypto** | AES-256-GCM, SHA-256 해시체인 | Class B |
| **ai** | TFLite 엣지 추론 (5종 모델) | Class B |
| **sync** | CRDT 오프라인 동기화 | Class A |
| **flutter-bridge** | FFI 브리지 (Rust ↔ Flutter) | Class B |

### 3.9 소프트웨어 항목 총괄표

| 구분 | 항목 수 | 비고 |
|------|---------|------|
| Go 마이크로서비스 | 22개 | auth, user, device, measurement 등 |
| Gateway | 1 | Kong + Keycloak |
| Rust 코어 모듈 | 9개 | differential, fingerprint, ble, nfc, ai 등 |
| Flutter 앱 | 1 | 모바일 클라이언트 |
| **합계** | **33+** | 30+ 서비스 생태계 |

---

## 4. 인터페이스 설계 (Clause 5.3.4)

### 4.1 서비스 간 gRPC 인터페이스

| 호출자 | 피호출자 | 프로토콜 | 용도 |
|--------|----------|---------|------|
| Flutter 앱 | Gateway | gRPC-Web / gRPC | 모든 API 요청 |
| Gateway | auth-service | gRPC | JWT 검증, 사용자 인증 |
| Gateway | 각 마이크로서비스 | gRPC | API 라우팅 (요청 위임) |
| measurement-service | ai-inference-service | gRPC | 클라우드 AI 분석 요청 |
| payment-service | subscription-service | gRPC | 구독 활성화 |

**Proto 정의**: `backend/shared/proto/manpasik.proto`  
**생성 코드**: Go `backend/shared/gen/go/v1/`, Dart `frontend/flutter-app/lib/generated/`

### 4.2 REST 외부 API

| 서비스 | 외부 시스템 | 프로토콜 | 용도 |
|--------|-------------|---------|------|
| payment-service | Toss Payments API | HTTPS REST | 결제 승인/취소 |
| notification-service | FCM (Firebase Cloud Messaging) | HTTPS | 푸시 알림 |
| auth-service | Keycloak (OIDC) | OAuth2/OIDC | 인증 위임 |
| translation-service | 번역 API (선택) | HTTPS | 다국어 |
| vision-service | AI Vision API (선택) | HTTPS | 이미지 분석 |

### 4.3 하드웨어 인터페이스 (BLE/NFC)

| 인터페이스 | 프로토콜 | 역할 | 구현 |
|------------|---------|------|------|
| **BLE 5.0** | GATT | 리더기 ↔ Flutter 앱 간 측정 데이터 전송 | Rust `ble` 모듈 (btleplug) |
| **NFC** | ISO 14443 | 카트리지 인증, UID/펌웨어 읽기 | Rust `nfc` 모듈 |

**데이터 흐름**: 리더기 → [BLE GATT] → Flutter → Rust FFI → 차동측정 → 핑거프린트 → gRPC → Backend

---

## 5. 데이터 아키텍처 (Clause 5.3.5)

### 5.1 저장소별 역할

| 저장소 | 버전 | 역할 | 사용 서비스 |
|--------|------|------|------------|
| **PostgreSQL 16** | 16 | 주 관계형 데이터, 트랜잭션 | 전체 (22+ 서비스) |
| **TimescaleDB** | pg16 확장 | 시계열 측정 데이터 (hypertable) | measurement-service |
| **Redis 7** | 7 | 캐시, 세션, Rate Limit 카운터 | auth, device, subscription, gateway |
| **Milvus 2.4** | 2.4 | 벡터 유사도 검색 (896차원 핑거프린트) | measurement-service |
| **Elasticsearch 8.14** | 8.14 | 전문 검색, 로그 집계 | measurement, community |
| **MinIO** | S3 호환 | 오브젝트 스토리지 (파일 업로드) | gateway |
| **Kafka (Redpanda)** | v24.2 | 이벤트 스트림, 메시지 버스 | 전체 (이벤트 발행/구독) |

### 5.2 스키마 개요 (PostgreSQL)

| 스키마/파일 | 용도 |
|-------------|------|
| `01-auth.sql` | users, refresh_tokens |
| `02-device.sql`, `03-device.sql` | devices, device_events |
| `02-user.sql` | profiles, subscriptions |
| `03-measurement.sql`, `04-measurement.sql` | measurements, measurement_sessions |
| `05-subscription.sql` | subscription_plans, user_subscriptions |
| `06-shop.sql` | products, cart_items, orders |
| `07-payment.sql` | payments, refunds |
| `08-ai-inference.sql` | inference_logs, models |
| `09-cartridge.sql` | cartridges, cartridge_types |
| `10-calibration.sql` | calibrations |
| `11-coaching.sql` | coaching_sessions, recommendations |
| `12-notification.sql` | notification_logs |
| `13-family.sql` | families, family_members |
| `14-health-record.sql` | health_records (FHIR 호환) |
| `15-telemedicine.sql` | telemedicine_sessions |
| `16-reservation.sql` | reservations |
| `17-community.sql` | posts, comments |
| `18-admin.sql` | admin_settings, audit_logs |
| `19-prescription.sql` | prescriptions |
| `20-translation.sql` | translations |
| `21-video.sql` | videos |
| `22-regions-facilities-doctors.sql` | 지역/시설/의사 마스터 |
| `23-data-sharing-consents.sql` | 동의 관리 |
| `24-prescription-fulfillment.sql` | 처방 이행 |
| `25-admin-settings-ext.sql` | LLM 어시스턴트 설정 등 |

### 5.3 데이터 흐름도 (핵심 시나리오)

```
[측정 플로우]
리더기 → [BLE GATT] → Flutter 앱
  → Rust FFI: 차동측정 (S_det - α × S_ref)
  → Rust FFI: 핑거프린트 생성 (896차원)
  → Rust FFI: 엣지 AI 추론 (TFLite)
  → 로컬 저장 (SQLite/Hive)
  → [gRPC] → gateway → measurement-service
    → PostgreSQL (측정 데이터 저장)
    → Milvus (벡터 저장)
    → Kafka (manpasik.measurement.completed)
      → ai-inference-service (클라우드 분석)
      → coaching-service (코칭 리포트)
      → notification-service (결과 알림)
  → Flutter 앱: 결과 표시

[결제 플로우]
Flutter 앱 → gateway → payment-service
  → CreatePayment (DB 저장)
  → [Toss REST API] → ConfirmPayment
  → DB 상태 업데이트 (confirmed)
  → Kafka (manpasik.payment.completed)
    → subscription-service (구독 활성화)
    → notification-service (결제 완료 알림)
```

---

## 6. 보안 아키텍처 (Clause 5.3.6)

### 6.1 인증·인가

| 항목 | 구현 |
|------|------|
| **인증** | Keycloak 25.0 (OIDC) → JWT 발급 |
| **인가** | RBAC (5개 역할: SuperAdmin, Admin, Moderator, Support, Analyst) |
| **API Gateway** | JWT 검증, 역할 확인, Rate Limiting |
| **서비스 간** | 내부 gRPC 호출 시 서비스 계정 또는 JWT 전달 |

### 6.2 데이터 암호화

| 대상 | 방식 |
|------|------|
| **전송 중** | TLS 1.3 (gRPC, REST, BLE 암호화) |
| **저장 시** | AES-256-GCM (PHI, PII, 시크릿 설정) |
| **암호화 키** | 환경변수 (`CONFIG_ENCRYPTION_KEY`, `JWT_SECRET`) — 코드 하드코딩 금지 |

### 6.3 감사 로그 (Audit Trail)

| 대상 | 보존 기간 | 용도 |
|------|----------|------|
| admin_settings 변경 | 최소 10년 | 규제 대응 |
| 측정 데이터 접근 | 최소 10년 | PHI 접근 추적 |
| 결제/구독 변경 | 7년 | 상업 거래 증빙 |
| 로그인 실패 | 1년 | 보안 사고 분석 |

### 6.4 네트워크 보안

| 항목 | 구현 |
|------|------|
| **서비스 간** | K8s 네트워크 정책 (내부 전용) |
| **외부 접근** | API Gateway만 외부 노출 |
| **CORS** | 화이트리스트 기반 origin 제한 |
| **Rate Limiting** | 60 req/min (기본), API별 상이 가능 |

---

## 7. 배포 아키텍처 (Clause 5.3.7)

### 7.1 Docker 컨테이너화

| 항목 | 내용 |
|------|------|
| **빌드** | 멀티스테이지 Dockerfile (`backend/services/*/Dockerfile`) |
| **이미지** | `manpasik/{service-name}:{version}` 또는 `ghcr.io/{repo}/{service}` |
| **개발 환경** | `infrastructure/docker/docker-compose.dev.yml` |

### 7.2 Kubernetes 배포

| 환경 | 경로 | 비고 |
|------|------|------|
| **base** | `infrastructure/kubernetes/base/` | 공통 Deployment, Service, ConfigMap |
| **dev** | `infrastructure/kubernetes/overlays/dev/` | 개발용 리소스·레플리카 패치 |
| **staging** | `infrastructure/kubernetes/overlays/staging/` | 스테이징 검증 |
| **production** | `infrastructure/kubernetes/overlays/production/` | HPA, PDB, 프로덕션 리소스 |

**오케스트레이션**: Deployment + Service + HPA (서비스별), Kustomize 환경별 오버레이

### 7.3 CI/CD (GitHub Actions)

| 워크플로우 | 파일 | 역할 |
|------------|------|------|
| **CI** | `.github/workflows/ci.yml` | push/PR 시 빌드·테스트·린트 (Rust, Go, Flutter) |
| **CD** | `.github/workflows/cd.yml` | main 브랜치 푸시 시 Docker 빌드·K8s 배포 (Rolling Update) |

**CI 단계**: Rust (clippy, test) → Go (golangci-lint, test) → Flutter (analyze, test) → Docker 이미지 빌드

---

## 8. 안전 등급 매핑 (Clause 5.3.8)

### 8.1 서비스별 IEC 62304 Class 분류

| 안전 등급 | 서비스/모듈 | 근거 |
|----------|-------------|------|
| **Class B** | auth-service, user-service, device-service, measurement-service | 인증·측정 데이터 직접 처리, PHI 접근 |
| **Class B** | ai-inference-service, cartridge-service, calibration-service | 측정 정확성·카트리지 무결성 |
| **Class B** | telemedicine-service, prescription-service, health-record-service, vision-service | 의료·건강 데이터 처리 |
| **Class B** | Rust: differential, fingerprint, dsp, ble, nfc, crypto, ai, flutter-bridge | 측정·암호화·하드웨어 연동 |
| **Class A** | payment-service, subscription-service, shop-service | 건강 데이터에 직접 영향 없음 |
| **Class A** | coaching-service, notification-service, community-service, video-service, translation-service | 부가 서비스 |
| **Class A** | reservation-service, family-service, admin-service | 관리·예약 |
| **Class A** | Gateway, Rust sync | 라우팅·동기화 |

### 8.2 Class B 개발 요구사항

Class B 항목은 다음을 준수해야 합니다:
- 소프트웨어 아키텍처 설계 (§5.3) — 본 문서
- 소프트웨어 통합 및 통합 시험 (§5.6)
- 소프트웨어 시스템 시험 (§5.7)
- SOUP 위험 평가 (§5.3.9)

---

## 9. 외부 의존성 (SOUP) (Clause 5.3.9)

### 9.1 Go 주요 의존성

| 라이브러리 | 버전 | 용도 | 라이선스 | 안전 영향 |
|-----------|------|------|---------|----------|
| google.golang.org/grpc | 1.78+ | gRPC 서버/클라이언트 | Apache-2.0 | B |
| github.com/jackc/pgx/v5 | 5.x | PostgreSQL 드라이버 | MIT | B |
| github.com/twmb/franz-go | 1.x | Kafka 클라이언트 | BSD-3 | A |
| github.com/redis/go-redis/v9 | 9.x | Redis 클라이언트 | BSD-2 | A |
| github.com/golang-jwt/jwt/v5 | 5.x | JWT 토큰 | MIT | B |
| go.uber.org/zap | 1.x | 구조화 로깅 | MIT | A |
| github.com/google/uuid | 1.x | UUID 생성 | BSD-3 | A |
| google.golang.org/protobuf | 1.x | Protobuf | BSD-3 | B |
| github.com/milvus-io/milvus-sdk-go/v2 | 2.4 | Milvus 벡터 DB | Apache-2.0 | B |
| github.com/minio/minio-go/v7 | 7.x | MinIO 객체 스토리지 | Apache-2.0 | A |

### 9.2 Rust 주요 의존성

| 라이브러리 | 버전 | 용도 | 라이선스 |
|-----------|------|------|---------|
| flutter_rust_bridge | 2.x | FFI 브리지 | MIT |
| ring | 0.17+ | 암호화 (AES, SHA) | ISC |
| aes-gcm | 1.x | AES-256-GCM | MIT/Apache-2.0 |
| sha2 | 0.10 | SHA-256 해시체인 | MIT |
| tokio | 1.x | 비동기 런타임 | MIT |
| serde | 1.x | 직렬화 | MIT/Apache-2.0 |
| btleplug | 0.11 | BLE 통신 | MIT |
| tflitec | 0.5 | TFLite 엣지 추론 | Apache-2.0 |
| rustfft, dasp | - | DSP, FFT | MIT/Apache-2.0 |
| ndarray | - | 선형 대수 (핑거프린트) | Apache-2.0 |

### 9.3 Flutter 주요 의존성

| 라이브러리 | 버전 | 용도 | 라이선스 |
|-----------|------|------|---------|
| grpc | 4.x | gRPC 클라이언트 | BSD-3 |
| protobuf | 3.x | Protobuf 직렬화 | BSD-3 |
| flutter_riverpod | 2.x | 상태 관리 | MIT |
| go_router | 13.x | 라우팅 | BSD-3 |
| fl_chart | 0.69 | 차트 | MIT |
| dio | 5.x | HTTP 클라이언트 | MIT |
| sqflite / hive_flutter | - | 로컬 저장소 | BSD-3 |
| flutter_secure_storage | (선택) | 보안 저장 | BSD-3 |

### 9.4 인프라 SOUP

| 구성요소 | 버전 | 용도 |
|----------|------|------|
| PostgreSQL | 16 | 관계형 DB |
| TimescaleDB | pg16 | 시계열 확장 |
| Milvus | 2.4 | 벡터 DB |
| Redis | 7 | 캐시·세션 |
| Kafka (Redpanda) | v24.2 | 이벤트 스트림 |
| Elasticsearch | 8.14 | 검색 |
| MinIO | S3 호환 | 객체 스토리지 |
| Kong | 3.7 | API Gateway |
| Keycloak | 25.0 | 인증 (OIDC) |
| Kubernetes | - | 오케스트레이션 |
| Docker | - | 컨테이너화 |

---

## 10. 부록

### 부록 A: 서비스 통신 매트릭스

| 호출자 → 피호출자 | 프로토콜 | 용도 |
|------------------|---------|------|
| Flutter → Gateway | gRPC-Web | 모든 API 요청 |
| Gateway → auth-service | gRPC | JWT 검증 |
| Gateway → 각 서비스 | gRPC | API 라우팅 |
| payment → Toss API | HTTPS | 결제 승인/취소 |
| notification → FCM | HTTPS | 푸시 알림 |
| admin → Kafka | Kafka | config.changed 발행 |
| payment ← Kafka | Kafka | config.changed 구독 |

### 부록 B: Kafka 토픽 목록 (18개)

| 토픽 | Producer | Consumer(s) |
|------|----------|-------------|
| manpasik.measurement.completed | measurement-service | ai-inference, coaching, notification |
| manpasik.measurement.session.started | measurement-service | device-service |
| manpasik.measurement.session.ended | measurement-service | coaching, health-record |
| manpasik.payment.completed | payment-service | subscription, shop, notification |
| manpasik.payment.failed | payment-service | notification, admin |
| manpasik.subscription.changed | subscription-service | cartridge, notification, user |
| manpasik.cartridge.verified | cartridge-service | measurement, notification |
| manpasik.cartridge.depleted | cartridge-service | notification, shop |
| manpasik.notification.send | 여러 서비스 | notification-service |
| manpasik.user.registered | auth-service | user, notification, coaching |
| manpasik.user.profile.updated | user-service | coaching, health-record |
| manpasik.device.registered | device-service | notification, admin |
| manpasik.device.status.changed | device-service | notification |
| manpasik.ai.risk.detected | ai-inference-service | notification, coaching |
| manpasik.reservation.created | reservation-service | notification, telemedicine |
| manpasik.prescription.created | prescription-service | notification, health-record |
| manpasik.community.post.created | community-service | notification, translation |
| manpasik.dlq | 전체 (실패 시) | admin (수동 재처리) |

### 부록 C: 서비스 내부 구조 (표준 패턴)

모든 Go 백엔드 서비스는 다음 계층 구조를 따릅니다:

```
service-name/
├── cmd/main.go                    # 엔트리 포인트, DI
└── internal/
    ├── handler/grpc.go            # gRPC 핸들러 (프레젠테이션)
    ├── service/*.go               # 비즈니스 로직 + 도메인 모델
    └── repository/
        ├── postgres/*.go          # PostgreSQL 구현
        └── memory/*.go            # 인메모리 (테스트/개발)
```

**패턴**: handler → service → repository (의존성 역전)

---

**마지막 업데이트**: 2026-02-13 (v2.0)  
**승인**: — (품질관리 책임자 검토 대기)
