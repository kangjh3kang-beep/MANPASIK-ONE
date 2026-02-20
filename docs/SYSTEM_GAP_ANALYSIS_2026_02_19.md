# ManPaSik 시스템 기획-계획-구현 갭 분석 보고서

**날짜**: 2026-02-19
**분석 범위**: 구현기획안(Sprint 11~15) + 검증보고서(v7.0) + 실제 코드베이스 교차 검증
**분석 방법**: 기획서 전수 대조, 코드베이스 grep/glob 실측, 빌드 검증

---

## 1. Executive Summary

```
┌──────────────────────────────────────────────────────────────────────┐
│  ManPaSik 시스템 현황 종합 (Sprint 14 기준, 2026-02-19)               │
│                                                                       │
│  ■ 기획서 대비 구현율                                                 │
│    Phase 1 (MVP)       ████████████████░░░░  80%                     │
│    Phase 2 (Core)      ████████████░░░░░░░░  60%                     │
│    Phase 3 (Advanced)  ██████████░░░░░░░░░░  50%                     │
│    Phase 4 (Ecosystem) ░░░░░░░░░░░░░░░░░░░░   0% (미착수)           │
│    Phase 5 (Future)    ░░░░░░░░░░░░░░░░░░░░   0% (미착수)           │
│                                                                       │
│  ■ 기술 계층별 달성률                                                 │
│    Go 백엔드 서비스      ████████████████████  100% (22/22 빌드)     │
│    Proto gRPC Handler    ████████████████████  100% (180 RPC)        │
│    Gateway REST 노출     ██████████████████░░   91% (164/180)        │
│    DB 스키마             ████████████████████  100% (85+ 테이블)     │
│    Flutter 화면          ████████████████████  100% (69 화면)        │
│    Flutter REST 연동     ████████████████████  100% (182 메서드)     │
│    Repository Impl       ████████████████████  100% (16/16, 에러 0)  │
│    SDK 외부 연동          ████████████████░░░░   80% (7개 래퍼)      │
│    UI 스토리보드 일치율   █████████████░░░░░░░   67% (73씬 기준)     │
│    E2E 통합 테스트        ░░░░░░░░░░░░░░░░░░░░    0% (미착수)       │
│    모니터링/HoloBody     ████████████████████  100% (v6.2)           │
│                                                                       │
│  ■ 스프린트 진행 현황                                                 │
│    Sprint 11 (Gateway REST)    ████████████████████ 완료 (164 EP)    │
│    Sprint 12 (Flutter 연동)    ████████████████████ 완료 (182 메서드) │
│    Sprint 13 (SDK 연동)        ████████████████████ 완료 (7개 서비스) │
│    Sprint 14 (UI 폴리싱)       █████████████░░░░░░░  65% (HoloBody만)│
│    Sprint 15 (E2E 테스트)      ░░░░░░░░░░░░░░░░░░░░   0% (미착수)   │
└──────────────────────────────────────────────────────────────────────┘
```

---

## 2. 기획 vs 계획 vs 구현 — 3자 대조표

### 2.1 스프린트 11: Gateway REST 100%

| 계획 항목 | 계획 목표 | 실측값 | 상태 | 비고 |
|-----------|----------|--------|------|------|
| REST 엔드포인트 총수 | 159개+ | **164개** | ✅ 초과 달성 | |
| health-record REST | 13개 추가 | 포함 (community_routes) | ✅ | 전용 파일 미분리 |
| prescription REST | 12개 추가 | 포함 (community_routes) | ✅ | 전용 파일 미분리 |
| family REST | 8개 추가 | 포함 (community_routes) | ✅ | 전용 파일 미분리 |
| community REST | 8개 추가 | community_routes 46개 | ✅ | 다른 서비스 라우트 포함 |
| video REST | 6개 추가 | 포함 | ✅ | |
| notification REST | 5개 추가 | notification_routes 8개 | ✅ | |
| translation REST | 6개 추가 | translation_routes 7개 | ✅ | |
| telemedicine REST | 5개 추가 | 포함 | ✅ | |
| 기존 잔여 REST | 11개 추가 | 포함 | ✅ | |
| Placeholder 수정 7건 | 실제 gRPC 연동 | **미확인** | ⚠️ | 일부 우회 잔존 가능 |
| 미구현 Flutter 화면 3개 | 신규 생성 | **3/3 존재** | ✅ | escalation, subscription_cancel, data_export |

### 2.2 스프린트 12: Flutter 프론트엔드 연동 100%

| 계획 항목 | 계획 목표 | 실측값 | 상태 | 비고 |
|-----------|----------|--------|------|------|
| rest_client.dart 메서드 | 159개+ | **182개** | ✅ 초과 달성 | |
| Repository 구현체 | 20개 수정 | **16/16 완료** | ✅ | UnimplementedError 0건 |
| Placeholder 잔여 | 0건 | **0건** | ✅ | grep 확인 |
| Provider 연동 | 7개 신규 | 구현됨 | ✅ | |
| 화면 바인딩 | 7개 수정 | 구현됨 | ✅ | |

### 2.3 스프린트 13: SDK 외부 연동

| 계획 항목 | 계획 목표 | 실측값 | 상태 | 비고 |
|-----------|----------|--------|------|------|
| Toss Payments | 프로덕션 구현체 | **TossPaymentService** 존재 | ✅ | 환경변수 없으면 Simulated 폴백 |
| FCM 푸시 | 프로덕션 구현체 | **FcmNotificationService** 존재 | ✅ | 환경변수 없으면 Polling 폴백 |
| PASS 본인인증 | 프로덕션 구현체 | **PassIdentityService** 존재 | ✅ | 환경변수 없으면 Simulated 폴백 |
| WebRTC | flutter_webrtc | **RestSignalingWebRtcService** | ⚠️ | REST 시그널링, 실제 P2P 미검증 |
| OAuth2 (Google/Kakao/Apple) | 3개 플랫폼 | **SocialAuthService** | ⚠️ | 브라우저 OAuth, 네이티브 SDK 미연동 |
| HealthKit/Health Connect | health 패키지 | **HealthConnectService** | ⚠️ | PlatformChannel, 실제 연동 미검증 |
| 약국조회 | 공공데이터 API | **RestPharmacyService** | ✅ | 환경변수 없으면 Simulated 폴백 |
| 공공데이터 서비스 | 프로덕션 구현체 | **SimulatedPublicDataService만** | ❌ | 프로덕션 구현체 없음 |
| SimulatedXxxService 0건 | 프로덕션 빌드에서 | **5개 Simulated 클래스 잔존** | ⚠️ | 팩토리 패턴으로 폴백 |

### 2.4 스프린트 14: 스토리보드 UI 100% + 폴리싱

| 계획 항목 | 계획 목표 | 실측값 | 상태 | 비고 |
|-----------|----------|--------|------|------|
| **Lottie 애니메이션 교체 5건** | 실제 에셋 | **미구현** | ❌ | 플레이스홀더 텍스트 유지 |
| **홈 대시보드 미니차트** | 스파크라인 추가 | **미구현** | ❌ | |
| **알림 카테고리 필터 탭** | 5탭 필터 | **미구현** | ❌ | |
| **커뮤니티 댓글/챌린지 UI** | 리더보드, 인증 | **미구현** | ❌ | |
| **가족 QR 초대** | qr_flutter | **미구현** | ❌ | |
| **데이터허브 My Zone 오버레이** | 개인 기준선 | **미구현** | ❌ | |
| **구독 연간/월간 토글** | 쿠폰 입력 | **미구현** | ❌ | |
| **화상진료 재연결 로직** | ICE restart | **미구현** | ❌ | |
| **설정 약관변경 이력** | 화면 | **미구현** | ❌ | |
| **접근성 보강 전수** | Semantics 48x48dp | **미구현** | ❌ | |
| HoloBody v6.2 | 품질 고도화 | **완료** | ✅ | 5000파티클, HSL, Catmull-Rom |
| 모니터링 대시보드 | 크래시 해결 | **완료** | ✅ | SafeHitTestWrapper + SliverList.builder |

### 2.5 스프린트 15: E2E 테스트 + 빌드 검증 + 출시 준비

| 계획 항목 | 계획 목표 | 실측값 | 상태 | 비고 |
|-----------|----------|--------|------|------|
| Flutter E2E 30개 시나리오 | 30개 | **0개** | ❌ | integration_test 디렉토리 없음 |
| Go 통합 테스트 | 추가 | 342개 (단위만) | ⚠️ | 통합 테스트 별도 없음 |
| Go 빌드 재검증 | 22/22 | **22/22 PASS** | ✅ | |
| Flutter analyze | 에러 0 | **에러 0 (737 info)** | ✅ | |
| 보안 감사 | 7항목 | **미실시** | ❌ | |
| 성능 프로파일링 | 60fps | **미실시** | ❌ | |
| 출시 체크리스트 | 17항목 | **미실시** | ❌ | |

---

## 3. 백엔드 상세 검증 결과

### 3.1 Go 서비스 28개 상태

| # | 서비스 | main.go | Dockerfile | 테스트 | 비고 |
|---|--------|---------|------------|--------|------|
| 1 | admin-service | ✅ | ✅ | 33 | config_manager_test 포함 |
| 2 | ai-inference-service | ✅ | ✅ | 23 | |
| 3 | **analytics-service** | ❌ | ❌ | 0 | **빈 서비스 (Phase 4)** |
| 4 | auth-service | ✅ | ✅ | 8 | |
| 5 | calibration-service | ✅ | ✅ | 12 | |
| 6 | cartridge-service | ✅ | ✅ | 27 | |
| 7 | coaching-service | ✅ | ✅ | 11 | |
| 8 | community-service | ✅ | ✅ | 15 | |
| 9 | device-service | ✅ | ✅ | 10 | |
| 10 | **emergency-service** | ❌ | ❌ | 0 | **빈 서비스** |
| 11 | family-service | ✅ | ✅ | 19 | |
| 12 | gateway | ✅ | ✅ | 0 | 라우터, 테스트 별도 |
| 13 | health-record-service | ✅ | ✅ | 19 | |
| 14 | **iot-gateway-service** | ❌ | ❌ | 0 | **빈 서비스 (Phase 4)** |
| 15 | **marketplace-service** | ❌ | ❌ | 0 | **빈 서비스 (Phase 4)** |
| 16 | measurement-service | ✅ | ✅ | 14 | |
| 17 | **nlp-service** | ❌ | ❌ | 0 | **빈 서비스 (Phase 4)** |
| 18 | notification-service | ✅ | ✅ | 18 | |
| 19 | payment-service | ✅ | ✅ | 11 | |
| 20 | prescription-service | ✅ | ✅ | 29 | |
| 21 | reservation-service | ✅ | ✅ | 18 | |
| 22 | shop-service | ✅ | ✅ | 12 | |
| 23 | subscription-service | ✅ | ✅ | 14 | |
| 24 | telemedicine-service | ✅ | ✅ | 12 | |
| 25 | translation-service | ✅ | ✅ | 13 | |
| 26 | user-service | ✅ | ✅ | 10 | |
| 27 | video-service | ✅ | ✅ | 14 | |
| 28 | **vision-service** | ✅ | ❌ | 0 | main.go만, Dockerfile/테스트 없음 |
| | **합계** | **23/28** | **22/28** | **342** | |

### 3.2 Gateway REST 라우트 파일별 분포

| 라우트 파일 | 엔드포인트 수 | 포함 서비스 |
|------------|-------------|-----------|
| community_routes.go | 46 | community + health-record + prescription + video + telemedicine + family |
| measurement_routes.go | 31 | measurement + device + calibration + cartridge |
| market_routes.go | 23 | shop + payment + subscription(일부) |
| admin_routes.go | 20 | admin |
| notification_routes.go | 8 | notification |
| subscription_routes.go | 8 | subscription |
| coaching_routes.go | 7 | coaching |
| translation_routes.go | 7 | translation |
| user_routes.go | 7 | user |
| auth_routes.go | 6 | auth |
| rest_handler.go | 1 | 메인 라우트 등록 |
| **합계** | **164** | |

### 3.3 Proto vs REST 노출률

| 항목 | 수치 |
|------|------|
| Proto 파일 | 2개 (manpasik.proto + health.proto) |
| Proto service 정의 | 23개 (21 + 2) |
| Proto RPC 메서드 | **180개** (174 + 6) |
| Gateway REST 엔드포인트 | **164개** |
| **REST 노출률** | **91.1%** (164/180) |
| 미노출 RPC | **~16개** |

### 3.4 빈 서비스 (main.go 없음, 디렉토리만 존재)

| 서비스 | 용도 | Phase | 상태 |
|--------|------|-------|------|
| analytics-service | 사용자 행동 분석 | Phase 4 | 디렉토리만 |
| emergency-service | 긴급 대응 전용 | Phase 3 | 디렉토리만 |
| iot-gateway-service | IoT 디바이스 게이트웨이 | Phase 4 | 디렉토리만 |
| marketplace-service | SDK 마켓플레이스 | Phase 4 | 디렉토리만 |
| nlp-service | 자연어 처리 | Phase 4 | 디렉토리만 |
| vision-service | 컴퓨터 비전 (음식 분석) | Phase 3 | main.go만, 테스트/Docker 없음 |

---

## 4. 프론트엔드 상세 검증 결과

### 4.1 정량 지표

| 항목 | v7.0 보고서 (Sprint 10) | 실측 (Sprint 14) | 변화 |
|------|------------------------|------------------|------|
| Flutter 화면 | 72개 | **69개** | -3 (통합/리팩터링) |
| GoRouter 라우트 | 72개 | **84개** | +12 |
| REST 클라이언트 메서드 | (계획 159) | **182개** | 계획 초과 |
| Repository REST | (계획 20) | **16개** | 구조 최적화 |
| UnimplementedError | (계획 0) | **0건** | ✅ |
| Flutter 단위 테스트 파일 | ? | **14개** | |
| E2E 통합 테스트 | 0개 | **0개** | ❌ 미착수 |

### 4.2 SDK 서비스 현황 (환경변수별 폴백)

| 서비스 | 프로덕션 구현체 | 시뮬레이션 폴백 | 환경변수 | 상태 |
|--------|---------------|----------------|---------|------|
| 결제 | TossPaymentService | SimulatedPaymentService | `TOSS_CLIENT_KEY` | ⚠️ 키 없으면 시뮬 |
| 푸시알림 | FcmNotificationService | PollingNotificationService | `FCM_ENABLED` | ⚠️ FCM 없으면 폴링 |
| 본인인증 | PassIdentityService | SimulatedIdentityService | `PASS_MERCHANT_ID` | ⚠️ 키 없으면 시뮬 |
| 화상통화 | RestSignalingWebRtcService | (시뮬레이션 모드) | — | ⚠️ P2P 미검증 |
| 소셜로그인 | SocialAuthService | (시뮬레이션) | `GOOGLE_CLIENT_ID` 등 | ⚠️ 네이티브 SDK 미연동 |
| 건강연동 | HealthConnectService | (시뮬레이션) | — | ⚠️ 실제 연동 미검증 |
| 약국조회 | RestPharmacyService | SimulatedPharmacyService | `PHARMACY_API_KEY` | ⚠️ 키 없으면 시뮬 |
| **공공데이터** | **없음** | **SimulatedPublicDataService** | — | ❌ 프로덕션 구현체 없음 |

### 4.3 네이티브 바이너리 스텁

| 항목 | 파일 | 상태 |
|------|------|------|
| BLE 스캔/연결 | rust_ffi_stub.dart | 시뮬레이션만 (Rust FFI 미빌드) |
| NFC 카트리지 인식 | rust_ffi_stub.dart | 시뮬레이션만 |
| WebRTC P2P | webrtc_service.dart | REST 시그널링만 (flutter_webrtc 패키지 미사용) |

---

## 5. 인프라/DB 검증 결과

| 항목 | 수치 |
|------|------|
| DB 스키마 파일 | 27개 (.sql) |
| DB 테이블 | 85+ (CREATE TABLE 합계) |
| Docker Compose | `infrastructure/docker-compose.dev.yml` 존재 |
| Kubernetes | **없음** |
| Proto 파일 | 2개 (manpasik.proto, health.proto) |
| Kafka 어댑터 | `backend/shared/events/kafka_adapter.go` 존재 |
| Config 관리 | 각 서비스별 `infrastructure/config/` |

---

## 6. 미구현/미완성 현황 종합

### 6.1 심각도 CRITICAL (서비스 불가 — 4건)

| # | 항목 | 상세 | 필요 작업 |
|---|------|------|----------|
| C-1 | **Rust FFI 네이티브 빌드** | BLE/NFC 실제 디바이스 연결 불가 | flutter_rust_bridge 크로스 컴파일 |
| C-2 | **WebRTC P2P 실제 연결** | REST 시그널링만, 실제 P2P 스트림 미동작 | flutter_webrtc 패키지 + TURN 서버 |
| C-3 | **결제 실제 연동** | Toss 키 없으면 시뮬레이션 | Toss 가맹점 계약 + 키 발급 |
| C-4 | **E2E 통합 테스트 전무** | integration_test 디렉토리조차 없음 | 30개 시나리오 작성 |

### 6.2 심각도 HIGH (주요 기능 제한 — 8건)

| # | 항목 | 상세 | 필요 작업 |
|---|------|------|----------|
| H-1 | **Sprint 14 UI 보강 미완** | 계획된 40건 중 ~35건 미구현 | Lottie, 미니차트, 필터탭 등 |
| H-2 | **FCM 실제 구성** | Firebase 프로젝트 미설정 | google-services.json + flutterfire |
| H-3 | **OAuth 네이티브 SDK** | 브라우저 OAuth만, 네이티브 미연동 | google_sign_in, kakao_sdk 패키지 |
| H-4 | **공공데이터 프로덕션** | SimulatedPublicDataService만 존재 | 공공데이터포털 API 키 + 구현체 |
| H-5 | **빈 Go 서비스 5개** | analytics, emergency, iot-gateway, marketplace, nlp | Phase 4 구현 시 필요 |
| H-6 | **vision-service 미완성** | main.go만, Dockerfile/테스트 없음 | Docker + 테스트 추가 |
| H-7 | **보안 감사 미실시** | SQL Injection, XSS 등 7항목 미검증 | Sprint 15 보안 체크리스트 |
| H-8 | **접근성 보강 미실시** | Semantics, 48x48dp, 고대비 미검증 | 전수 검사 필요 |

### 6.3 심각도 MEDIUM (UX 저하 — 12건)

| # | 항목 | 상세 |
|---|------|------|
| M-1 | Lottie 애니메이션 플레이스홀더 5건 | splash, onboarding, order_complete, measurement, ble_scan |
| M-2 | 홈 대시보드 미니 스파크라인 차트 | 최근 7일 추세 없음 |
| M-3 | 알림 센터 카테고리 필터 탭 | 전체/측정/가족/의료/시스템 분류 없음 |
| M-4 | 커뮤니티 챌린지 리더보드 | UI 없음 |
| M-5 | 가족 QR 코드 초대 | qr_flutter 미사용 |
| M-6 | 데이터허브 My Zone 오버레이 | 개인 기준선 범위 없음 |
| M-7 | 구독 연간/월간 토글 + 쿠폰 | 미구현 |
| M-8 | 화상진료 자동 재연결 | ICE restart 미구현 |
| M-9 | 설정 약관 변경 이력 화면 | 미구현 |
| M-10 | 예약 취소 확인 다이얼로그 강화 | 미구현 |
| M-11 | Gateway Placeholder 7건 | avatar, emergency, food/exercise analyze 우회 잔존 가능 |
| M-12 | Gateway REST ~16개 미노출 | Proto 180 - REST 164 = 16 RPC 미노출 |

### 6.4 심각도 LOW (미래 확장 — 5건)

| # | 항목 | Phase |
|---|------|-------|
| L-1 | SDK 마켓플레이스 (marketplace-service) | Phase 4 |
| L-2 | 연합학습 AI (analytics-service) | Phase 4 |
| L-3 | IoT 게이트웨이 (iot-gateway-service) | Phase 4 |
| L-4 | NLP 서비스 (nlp-service) | Phase 4 |
| L-5 | 음성 명령 / 웨어러블 / 1792차원 | Phase 5 |

---

## 7. 스프린트별 완료 현황 요약

| Sprint | 계획 | 실제 달성 | 완료율 | 미완 핵심 항목 |
|--------|------|----------|--------|--------------|
| **11** | REST 76개 + 화면 3개 | REST 81개 + 화면 3개 | **100%** | Placeholder 7건 일부 잔존 |
| **12** | REST 클라이언트 76 + Repo 20 | 182메서드 + 16 Repo | **100%** | — |
| **13** | SDK 7개 프로덕션 | 7개 래퍼 (팩토리 패턴) | **80%** | 실제 키 없이 Simulated 폴백 |
| **14** | UI 40건 + Lottie 5 + 접근성 | HoloBody v6.2 + 크래시 2건 | **~35%** | UI 보강 35건, Lottie, 접근성 미구현 |
| **15** | E2E 30 + 보안 + 성능 | 미착수 | **0%** | 전체 미착수 |

---

## 8. 정량 대비표 (v7.0 → Sprint 14 실측)

| 지표 | v7.0 (Sprint 10) | Sprint 14 실측 | 계획 목표 | 달성률 |
|------|------------------|---------------|----------|--------|
| Go 서비스 (빌드 가능) | 11 | **22** | 22 | 100% |
| Go 서비스 (총 디렉토리) | 21 | **28** | 22+ | — |
| Proto RPC | 169 | **180** | 169+ | 100%+ |
| Gateway REST | 83 | **164** | 159+ | **103%** |
| REST 노출률 | 49% | **91%** | 100% | 91% |
| DB 테이블 | 84+ | **85+** | 84+ | 100% |
| Go 테스트 | 319 | **342** | 400+ | 86% |
| Flutter 화면 | 72 | **69** | 72+ | 96% |
| Flutter 라우트 | 72 | **84** | 72+ | 117% |
| REST 메서드 | ? | **182** | 159+ | 115% |
| Repository REST | ? | **16** | 16+ | 100% |
| UnimplementedError | ? | **0** | 0 | 100% |
| SDK 시뮬레이션 잔여 | 7 | **5** | 0 | 29% |
| 스토리보드 UI 일치 | 67% | **~70%** | 95%+ | 74% |
| E2E 테스트 | 0 | **0** | 30+ | 0% |
| Flutter analyze 에러 | 0 | **0** | 0 | 100% |
| Flutter 웹 빌드 | ✅ | **✅ (68.1s)** | ✅ | 100% |

---

## 9. 결론 및 권장 다음 단계

### 9.1 현재 프로젝트 건강도

| 영역 | Sprint 10 | Sprint 14 | 변화 |
|------|-----------|-----------|------|
| 아키텍처 설계 | 95/100 | **95/100** | = |
| 백엔드 구현 | 90/100 | **97/100** | +7 |
| 프론트엔드 구조 | 85/100 | **92/100** | +7 |
| Gateway 정합성 | 49/100 | **91/100** | **+42** |
| 프론트엔드 완성도 | 65/100 | **75/100** | +10 |
| 외부 연동 | 30/100 | **70/100** | **+40** |
| 모니터링/홀로그램 | — | **98/100** | 신규 |
| 통합 테스트 | 20/100 | **20/100** | = |
| **종합** | **68%** | **82%** | **+14%p** |

### 9.2 MVP 출시까지 남은 핵심 작업

| 우선순위 | 작업 | 예상 규모 |
|---------|------|----------|
| **P0** | E2E 통합 테스트 30개 시나리오 | 3일 |
| **P0** | SDK 실제 키 연동 검증 (Toss/FCM/OAuth) | 2일 (키 발급 의존) |
| **P1** | Sprint 14 잔여 UI 보강 35건 | 5일 |
| **P1** | Lottie 애니메이션 에셋 교체 5건 | 1일 |
| **P1** | Gateway REST 미노출 ~16개 추가 | 1일 |
| **P2** | 보안 감사 7항목 | 2일 |
| **P2** | 접근성 보강 전수 | 2일 |
| **P2** | 성능 프로파일링 (60fps) | 1일 |
| **P3** | Rust FFI 네이티브 빌드 (BLE/NFC) | 5일 |
| **P3** | WebRTC 실제 P2P 연결 | 3일 |
| **Phase 4** | 빈 서비스 5개 구현 | 향후 |
