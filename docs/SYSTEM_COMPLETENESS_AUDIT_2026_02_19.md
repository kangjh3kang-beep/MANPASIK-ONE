# 만파식(ManPaSik) 시스템 종합 완성도 전수조사 보고서

**조사일**: 2026-02-19
**조사범위**: 기획 문서 50+건, 백엔드 28개 서비스, 프론트엔드 16개 모듈, 인프라 전체
**현재 스프린트**: Sprint 14 (HoloBody v6.2 완료)

---

## 1. 전체 시스템 현황 요약

| 영역 | 항목수 | 완성 | 부분 | 미구현 | 완성률 |
|------|--------|------|------|--------|--------|
| **백엔드 Go 서비스** | 28 | 22 | 3 | 3 | 78.6% |
| **Proto RPC 구현** | 164 | 151 | - | 13 | 92.1% |
| **Gateway REST** | 175 | 175 | - | - | 100% |
| **Flutter 화면** | 73 | 73 | - | - | 100% |
| **Flutter REST 연동** | 182 | 110 | - | 72 | 60.4% |
| **Flutter 테스트** | 목표200 | 86 | - | 114 | 43.0% |
| **Go 테스트** | 목표300 | 340 | - | - | 100%+ |
| **E2E 테스트** | 목표30 | 0 | - | 30 | 0% |
| **접근성(Semantics)** | 목표80% | 2건 | - | - | ~3% |
| **국제화(i18n)** | 목표200키 | 60키 | - | 140키 | 30% |
| **DB 스키마** | 81 | 81 | - | - | 100% |
| **Dockerfile** | 28 | 22 | - | 6 | 78.6% |
| **K8s 매니페스트** | 24 | 24 | - | - | 100% |
| **CI/CD 파이프라인** | 2 | 2 | - | - | 100% |
| **모니터링** | 3 | 3 | - | - | 100% |

**종합 가중 완성도: 73.2%**

---

## 2. 미구현 항목 전수 목록

### 2.1 백엔드 — Proto 미구현 RPC 13건

| # | 서비스 | RPC | 유형 | 우선순위 |
|---|--------|-----|------|----------|
| B-1 | AuthService | `SocialLogin` | Unary | A |
| B-2 | AuthService | `ResetPassword` | Unary | A |
| B-3 | MeasurementService | `StreamMeasurement` | BiDi Stream | C |
| B-4 | DeviceService | `StreamDeviceStatus` | BiDi Stream | C |
| B-5 | DeviceService | `RequestOtaUpdate` | Unary | B |
| B-6 | SubscriptionService | `CheckCartridgeAccess` | Unary | A |
| B-7 | SubscriptionService | `ListAccessibleCartridges` | Unary | A |
| B-8 | AiInferenceService | `StreamChat` | Server Stream | C |
| B-9 | AdminService | `GetRevenueStats` | Unary | B |
| B-10 | AdminService | `GetInventoryStats` | Unary | B |
| B-11 | CommunityService | `GetChallengeLeaderboard` | Unary | A |
| B-12 | CommunityService | `UpdateChallengeProgress` | Unary | A |
| B-13 | TranslationService | `TranslateRealtime` | Unary | C |

### 2.2 백엔드 — Proto 미정의 서비스 6건

| # | 서비스 | 현재 상태 | 코드줄 | main.go | Dockerfile | 테스트 | 우선순위 |
|---|--------|-----------|--------|---------|------------|--------|----------|
| S-1 | emergency-service | HTTP 핸들러만 | 542 | ❌ | ❌ | 8개 | A |
| S-2 | marketplace-service | gRPC 리스너만 | 622 | 부분 | ❌ | 0 | B |
| S-3 | vision-service | gRPC 서버+시뮬 | 508 | ✅ | ❌ | 0 | B |
| S-4 | analytics-service | 빈 디렉토리 | 0 | ❌ | ❌ | 0 | C |
| S-5 | iot-gateway-service | 빈 디렉토리 | 0 | ❌ | ❌ | 0 | C |
| S-6 | nlp-service | 빈 디렉토리 | 0 | ❌ | ❌ | 0 | C |

### 2.3 프론트엔드 — 미사용 REST Client 메서드 72건

**인증 (1건)**: resetPassword
**측정 (2건)**: exportSingleMeasurement, exportToFhirObservations
**디바이스 (2건)**: requestOtaUpdate, updateDeviceStatus
**구독 (3건)**: checkFeatureAccess, checkCartridgeAccess, listAccessibleCartridges
**AI/코칭 (5건)**: listAiModels, analyzeFoodImage, analyzeExerciseVideo, listCoachingMessages, getWeeklyReport
**커뮤니티 (7건)**: getChallenges(중복), getQnaQuestions, createPostWithImage, getChallengeLeaderboard, updateChallengeProgress, createChallenge, getChallenge
**가족 (3건)**: updateMemberRole, listFamilyMembers, validateSharingAccess
**비디오/WebRTC (7건)**: createVideoRoom, getVideoRoom, endVideoRoom, sendVideoSignal, getVideoSignals, listVideoParticipants, getVideoRoomStats
**번역 (5건)**: detectLanguage, listSupportedLanguages, translateBatch, getTranslationHistory, getTranslationUsage
**원격진료 (3건)**: getConsultation, startVideoSession, endVideoSession
**처방전 (9건)**: createPrescription~updateDispensaryStatus
**관리자 (11건)**: getSystemMetrics~bulkSetConfigs
**결제 (1건)**: refundPayment
**상품리뷰 (2건)**: getProductReviews, createProductReview
**기타 (5건)**: uploadAvatar, streamChat, createInquiry, getPrescriptionByToken, getAlertDetail
**건강기록 (3건)**: getFhirRecord, shareFhirRecord, getAccessLogs
**예약 (3건)**: searchFacilities, getAvailableSlots, createReservation

### 2.4 프론트엔드 — 접근성 미적용

- **현재**: Semantics 위젯 2건 (monitoring_dashboard_screen.dart에서만)
- **목표**: 모든 인터랙티브 위젯에 Semantics label 추가
- **대상 화면**: 68개 전체

### 2.5 프론트엔드 — 국제화 부족

- **현재**: 60개 ARB 키 (6개 언어)
- **목표**: 200+ 키
- **하드코딩 한국어 문자열**: 8개 파일, 33건
- **미번역 화면**: 대부분

### 2.6 프론트엔드 — 테스트 부족

- **현재**: 86개 테스트 케이스 (15개 파일)
- **목표**: 200+ 케이스
- **미테스트 영역**: Widget 테스트, Integration 테스트, Repository 테스트

### 2.7 시뮬레이션/더미 데이터 잔존

| 화면/파일 | 더미 데이터 | 대체 필요 |
|-----------|------------|-----------|
| MarketRepositoryRest | 일반상품 8개 시뮬 | Shop API 연동 |
| MedicalRepositoryRest | 의사 5명 하드코딩 | Reservation API |
| MarketScreen | 카트리지 4종 하드코딩 | Cartridge API |
| HomeScreen | Timer 기반 랜덤 데이터 | Health Summary API |

### 2.8 인프라 미구현

| # | 항목 | 현재 상태 | 우선순위 |
|---|------|-----------|----------|
| I-1 | Terraform (IaC) | 미존재 | C |
| I-2 | ELK/Loki 로깅 스택 | stdout만 | C |
| I-3 | Istio 서비스 메시 | 미적용 | C |
| I-4 | Vault/KMS 시크릿 | K8s Secret만 | C |
| I-5 | Rust 코어 FFI 빌드 | 스텁만 존재 | C |

---

## 3. 100% 완성도 달성 구현계획

### Phase 1: 코어 완성 (Sprint 15-A, 3일)

#### 1-1. Proto 미구현 RPC 7건 구현 (A급)
**작업량**: ~700줄 Go 코드

| RPC | 서비스 | 예상 줄수 | 설명 |
|-----|--------|-----------|------|
| SocialLogin | auth-service | 80 | OAuth 토큰 검증 → JWT 발급 |
| ResetPassword | auth-service | 60 | 이메일 토큰 → 비밀번호 변경 |
| CheckCartridgeAccess | subscription-service | 50 | 구독 등급별 카트리지 접근 확인 |
| ListAccessibleCartridges | subscription-service | 60 | 사용자 구독에서 접근 가능 카트리지 목록 |
| GetChallengeLeaderboard | community-service | 70 | 챌린지 참여자 순위 조회 |
| UpdateChallengeProgress | community-service | 60 | 챌린지 진행률 업데이트 |
| RequestOtaUpdate | device-service | 80 | OTA 펌웨어 업데이트 요청 |

#### 1-2. 부분구현 서비스 3건 완성
**작업량**: ~600줄 Go 코드 + Dockerfile 3개

| 서비스 | 필요 작업 |
|--------|----------|
| emergency-service | cmd/main.go 생성 + Dockerfile 추가 |
| marketplace-service | gRPC 서버 시작 활성화 + Dockerfile + 테스트 10개 |
| vision-service | Dockerfile + 테스트 10개 |

#### 1-3. 미구현 서비스 3건 스켈레톤 생성
**작업량**: ~1200줄 Go 코드 + Dockerfile 3개

| 서비스 | 파일 구성 |
|--------|----------|
| analytics-service | cmd/main.go + handler + service(5메서드) + repo + test + Dockerfile |
| iot-gateway-service | cmd/main.go + handler + service(5메서드) + repo + test + Dockerfile |
| nlp-service | cmd/main.go + handler + service(5메서드) + repo + test + Dockerfile |

#### 1-4. Proto RPC B급 2건 구현
| RPC | 서비스 | 예상 줄수 |
|-----|--------|-----------|
| GetRevenueStats | admin-service | 80 |
| GetInventoryStats | admin-service | 80 |

### Phase 2: 프론트엔드 연동 완성 (Sprint 15-B, 3일)

#### 2-1. 미사용 REST 메서드 화면 연동 (핵심 32건)

**인증 그룹 (3건)**:
- resetPassword → forgot_password_screen 연동
- uploadAvatar → profile_edit_screen 연동
- socialLogin → login_screen 개선

**처방전 그룹 (9건)**:
- 처방전 CRUD → prescription_detail_screen 완전 연동
- checkDrugInteraction → 약물 상호작용 경고
- getMedicationReminders → 복약 알림

**화상진료 그룹 (7건)**:
- WebRTC 시그널링 → video_call_screen 완전 연동
- createVideoRoom/endVideoRoom → 실제 세션 관리

**구독 접근제어 (3건)**:
- checkCartridgeAccess → market_screen 잠금 UI
- checkFeatureAccess → 기능 접근 게이트

**커뮤니티 확장 (5건)**:
- 챌린지 리더보드/진행률 → challenge_screen 연동
- Q&A 질문 목록 → qna_screen 연동

**건강기록 (3건)**:
- FHIR 내보내기 → data_hub_screen 연동
- 접근 로그 → settings 연동

**결제 환불 (1건)**:
- refundPayment → order_detail_screen 연동

**상품 리뷰 (1건)**:
- 리뷰 조회/작성 → product_detail_screen 연동

#### 2-2. 시뮬레이션 데이터 → 실제 API 전환 (4건)

| 화면 | 현재 | 변경 |
|------|------|------|
| MarketRepositoryRest | 시뮬 상품 8개 | shop API 연동 |
| MedicalRepositoryRest | 의사 5명 하드코딩 | reservation API |
| MarketScreen | 카트리지 하드코딩 | cartridge API |
| HomeScreen | Timer 랜덤 | healthSummary API |

#### 2-3. TODO/FIXME 주석 해소 (33건)
- 각 파일의 TODO/FIXME를 실제 구현으로 전환

### Phase 3: 테스트 (Sprint 15-C, 4일)

#### 3-1. Flutter 위젯 테스트 추가 (+114건 → 200개 목표)

| 테스트 그룹 | 파일 수 | 테스트 수 |
|------------|--------|-----------|
| Auth 화면 테스트 | 4 | 16 |
| Home 화면 테스트 | 1 | 8 |
| Market 화면 테스트 | 3 | 12 |
| Medical 화면 테스트 | 3 | 12 |
| Community 화면 테스트 | 2 | 10 |
| Settings 화면 테스트 | 2 | 8 |
| DataHub 화면 테스트 | 2 | 10 |
| Repository 테스트 | 8 | 24 |
| Provider 테스트 | 4 | 14 |
| **합계** | **29** | **114** |

#### 3-2. E2E 테스트 시나리오 30건

| # | 시나리오 | 서비스 범위 |
|---|---------|------------|
| 1 | 회원가입→로그인→프로필 | auth, user |
| 2 | 디바이스 등록→목록 | device |
| 3 | 측정 시작→종료→이력 | measurement |
| 4 | 전체 플로우: Login→Measure→History | auth, device, measurement |
| 5 | 서비스 헬스체크 (전체) | all |
| 6 | 구독 생성→업그레이드→카트리지 접근 | subscription, cartridge |
| 7 | 상품 조회→장바구니→주문→결제 | shop, payment |
| 8 | 측정→AI 추론→코칭 | measurement, ai, coaching |
| 9 | 카트리지 인증→보정→측정 | cartridge, calibration, measurement |
| 10 | 위험 감지→알림 발송 | measurement, notification |
| 11 | AI 코칭 세션→리포트 | coaching |
| 12 | 음식 분석→칼로리 추적 | vision |
| 13 | 데이터 허브→차트 | data-hub |
| 14 | 커뮤니티 게시→댓글 | community |
| 15 | 챌린지 참여→진행률→리더보드 | community |
| 16 | 가족 그룹→초대→모니터링 | family |
| 17 | 화상진료 예약→상담→처방 | telemedicine, prescription |
| 18 | 게시글→번역→댓글 | community, translation |
| 19 | 건강기록 FHIR 내보내기 | health-record |
| 20 | 처방전→약국 전송→조제 추적 | prescription |
| 21 | 관리자 KPI 대시보드 | admin |
| 22 | 관리자 사용자 관리 | admin, user |
| 23 | 비상 연락망→응급 신고 | emergency |
| 24 | 마켓플레이스 파트너 상품 | marketplace |
| 25 | 알림 설정→푸시 수신 | notification |
| 26 | 데이터 내보내기 (JSON/CSV) | data-hub |
| 27 | 다국어 전환 | translation |
| 28 | 구독 취소→다운그레이드 | subscription |
| 29 | 오프라인 측정→동기화 | measurement |
| 30 | 보안 설정→2FA 활성화 | auth |

#### 3-3. Go 통합 테스트 (+60건)

| 테스트 그룹 | 테스트 수 |
|------------|-----------|
| Gateway→Auth 연동 | 8 |
| Gateway→Measurement 연동 | 8 |
| Gateway→Shop→Payment 연동 | 10 |
| Gateway→Subscription→Cartridge 연동 | 8 |
| Gateway→Community 연동 | 8 |
| Gateway→Medical 연동 | 10 |
| Kafka 이벤트 흐름 | 8 |
| **합계** | **60** |

### Phase 4: 품질 보강 (Sprint 15-D, 3일)

#### 4-1. 접근성(Semantics) 전수 추가

**대상**: 68개 화면, 32개 공유 위젯
**작업**: 모든 인터랙티브 위젯에 `Semantics(label:)` 추가
**예상 변경점**: ~300개 위젯

| 화면 그룹 | 위젯 수 (추정) |
|-----------|--------------|
| Auth (4 화면) | 20 |
| Home (1 화면) | 15 |
| Market (12 화면) | 45 |
| Medical (6 화면) | 30 |
| Community (6 화면) | 25 |
| Settings (13 화면) | 40 |
| DataHub (2+9위젯) | 25 |
| Admin (8 화면) | 35 |
| Family (6 화면) | 20 |
| Devices (4 화면) | 15 |
| 공유 위젯 (32개) | 30 |
| **합계** | **~300** |

#### 4-2. 국제화(i18n) 확장

**목표**: 60키 → 200+키
**작업**:
1. 하드코딩 한국어 33건 → ARB로 이동
2. 신규 키 140개 추가 (모든 화면 라벨/버튼/메시지)
3. 6개 언어 ARB 동기화

#### 4-3. 보안 점검 (OWASP Top 10)

| # | 항목 | 현재 | 조치 |
|---|------|------|------|
| 1 | SQL Injection | PGX 파라미터 쿼리 ✅ | 검증 완료 |
| 2 | XSS | Flutter 자동 이스케이프 ✅ | 검증 완료 |
| 3 | CSRF | JWT Bearer ✅ | 검증 완료 |
| 4 | 인증 결함 | JWT RS256 ✅ | MFA 추가 검토 |
| 5 | 접근 제어 | RBAC 미들웨어 ✅ | 세분화 검토 |
| 6 | 보안 구성 | TLS, CORS ✅ | 헤더 강화 |
| 7 | 암호화 | AES-256-GCM ✅ | - |

#### 4-4. 성능 기준선 측정

| 시나리오 | 목표 |
|---------|------|
| 로그인 API P95 | < 300ms |
| 측정 조회 P95 | < 150ms |
| 대시보드 폴링 P95 | < 200ms |
| 벡터 검색 P95 | < 200ms |

### Phase 5: 최종 검증 (Sprint 15-E, 1일)

#### 5-1. 전체 빌드 검증
- Go 28개 서비스 전수 빌드 (`GOWORK=off go build ./...`)
- Flutter 정적 분석 (`flutter analyze` 오류 0건)
- Flutter 웹 빌드 (`flutter build web`)
- Docker 이미지 22+6=28개 빌드

#### 5-2. 완성도 재측정
- Proto RPC 커버리지 → 100%
- REST 메서드 사용률 → 80%+
- 테스트 케이스 → 200+ (Flutter), 400+ (Go)
- 접근성 → 80%+ 위젯 Semantics
- 국제화 → 200+ ARB 키

---

## 4. 스프린트 일정 요약

```
Sprint 15-A (3일): 코어 완성
├── Day 1: Proto RPC 7건 + 부분구현 서비스 3건 완성
├── Day 2: 미구현 서비스 3건 스켈레톤 + Admin RPC 2건
└── Day 3: 전체 Go 빌드 검증 + Dockerfile 6개

Sprint 15-B (3일): 프론트엔드 연동
├── Day 4: 처방전/화상진료/인증 연동 (19건)
├── Day 5: 구독/커뮤니티/건강기록/결제 연동 (13건)
└── Day 6: 시뮬데이터 전환 + TODO 해소 + Flutter 검증

Sprint 15-C (4일): 테스트
├── Day 7-8: Flutter 위젯 테스트 114건
├── Day 9: E2E 시나리오 30건
└── Day 10: Go 통합 테스트 60건

Sprint 15-D (3일): 품질 보강
├── Day 11: 접근성 Semantics 300건
├── Day 12: i18n 140키 추가 + ARB 동기화
└── Day 13: 보안 점검 + 성능 기준선

Sprint 15-E (1일): 최종 검증
└── Day 14: 전체 빌드 + 완성도 재측정 + 보고서
```

**총 소요: 14일 (Sprint 15 = 2주)**

---

## 5. 항목별 체크리스트

### 백엔드 체크리스트 (33건)

- [ ] B-1: AuthService.SocialLogin 구현
- [ ] B-2: AuthService.ResetPassword 구현
- [ ] B-3: MeasurementService.StreamMeasurement (스트리밍, Phase 후순위)
- [ ] B-4: DeviceService.StreamDeviceStatus (스트리밍, Phase 후순위)
- [ ] B-5: DeviceService.RequestOtaUpdate 구현
- [ ] B-6: SubscriptionService.CheckCartridgeAccess 구현
- [ ] B-7: SubscriptionService.ListAccessibleCartridges 구현
- [ ] B-8: AiInferenceService.StreamChat (스트리밍, Phase 후순위)
- [ ] B-9: AdminService.GetRevenueStats 구현
- [ ] B-10: AdminService.GetInventoryStats 구현
- [ ] B-11: CommunityService.GetChallengeLeaderboard 구현
- [ ] B-12: CommunityService.UpdateChallengeProgress 구현
- [ ] B-13: TranslationService.TranslateRealtime (Phase 후순위)
- [ ] S-1: emergency-service main.go + Dockerfile
- [ ] S-2: marketplace-service gRPC 활성화 + Dockerfile + 테스트
- [ ] S-3: vision-service Dockerfile + 테스트
- [ ] S-4: analytics-service 전체 스켈레톤
- [ ] S-5: iot-gateway-service 전체 스켈레톤
- [ ] S-6: nlp-service 전체 스켈레톤
- [ ] 빌드 검증: 28개 서비스 전수

### 프론트엔드 체크리스트 (7건)

- [ ] F-1: 핵심 REST 메서드 32건 화면 연동
- [ ] F-2: 시뮬레이션 데이터 4건 실API 전환
- [ ] F-3: TODO/FIXME 33건 해소
- [ ] F-4: 접근성 Semantics ~300건 추가
- [ ] F-5: i18n ARB 140키 추가 + 하드코딩 33건 제거
- [ ] F-6: 위젯 테스트 114건 추가
- [ ] F-7: Flutter 분석 오류 0건 + 웹빌드 성공

### 테스트 체크리스트 (3건)

- [ ] T-1: E2E 시나리오 30건 작성
- [ ] T-2: Go 통합 테스트 60건
- [ ] T-3: 성능 기준선 측정

### 인프라 체크리스트 (Phase 후순위, 5건)

- [ ] I-1: Terraform IaC 구성
- [ ] I-2: ELK/Loki 로깅 스택
- [ ] I-3: Istio 서비스 메시
- [ ] I-4: Vault/KMS 시크릿 관리
- [ ] I-5: Rust 코어 FFI 통합 빌드

---

## 6. 우선순위 분류

### A급 — 즉시 (기능 완전성)
1. Proto RPC 미구현 7건 (Auth, Subscription, Community, Device)
2. 부분구현 서비스 3건 완성 (emergency, marketplace, vision)
3. 핵심 REST 메서드 32건 화면 연동

### B급 — 중요 (품질)
4. 미구현 서비스 3건 스켈레톤 (analytics, iot-gateway, nlp)
5. Admin RPC 2건 (GetRevenueStats, GetInventoryStats)
6. 시뮬레이션 데이터 실API 전환
7. TODO/FIXME 해소

### C급 — 보강 (비기능)
8. 접근성 Semantics 300건
9. i18n 140키 확장
10. 위젯 테스트 114건
11. E2E 시나리오 30건
12. 성능 기준선

### D급 — 장기 (인프라)
13. 스트리밍 RPC 3건 (Stream*)
14. TranslateRealtime
15. Terraform/Vault/Istio/ELK

---

## 7. 리스크 및 의존성

| 리스크 | 영향 | 완화 방안 |
|--------|------|----------|
| Proto 변경 시 Gen 재생성 | 빌드 실패 | protoc 자동화 스크립트 활용 |
| 스트리밍 RPC는 Gateway REST 미지원 | WebSocket 필요 | Phase 후순위로 분리 |
| Rust FFI 스텁 상태 | 실제 측정 불가 | 시뮬레이션 모드 유지 |
| 3개 빈 서비스 기획 부재 | 스켈레톤만 가능 | 최소 CRUD 구현 후 확장 |

---

**보고서 작성**: Claude
**최종 목표**: Sprint 15 (2주) 완료 후 종합 완성도 100%
**다음 행동**: Phase 1 즉시 시작 (Proto RPC 7건 + 서비스 완성)
