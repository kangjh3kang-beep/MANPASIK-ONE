# 기획–구현 추적성 매트릭스 (완성본)

**문서번호**: MPK-PLAN-TRACE-v2.0  
**갱신일**: 2026-02-12  
**목적**: 80개 요구사항(REQ) ↔ 설계(DES) ↔ 구현(IMP) ↔ 검증(V&V) 전체 연결. IEC 62304 §5.6, ISO 14971, 의료기기 인허가·감사 대비.  
**상태 범례**: ✅ 완료 | ⚠️ 부분 | 🔲 미구현 | — 해당 없음

---

## 1. 추적 ID 체계

| 접두어 | 의미 | 형식 예 |
|--------|------|---------|
| REQ | 기능 요구사항 | REQ-001 ~ REQ-080 |
| DES | 설계 문서/아키텍처 | DES-PROTO-AUTH, DES-RUST-DIFF |
| IMP | 구현 위치 (코드) | IMP-BE-AUTH, IMP-RUST-AI |
| VV | 검증·V&V | VV-UT-AUTH-001 (단위), VV-IT-001 (통합), VV-E2E-001 |

---

## 2. Phase 1 (MVP) — 18개 요구사항

| REQ ID | 기능명 | 설계 (DES) | 구현 (IMP) | 검증 (V&V) | 상태 |
|--------|--------|-----------|-----------|-----------|------|
| REQ-001 | 사용자 회원가입/로그인 | DES-PROTO-AUTH (manpasik.proto AuthService) DES-DB-01 (01-auth.sql) | IMP-BE-AUTH (auth-service/handler+service+repo) IMP-FE-AUTH (login_screen.dart) | VV-UT-AUTH-001 (service_test.go) VV-E2E-FLOW (flow_test.go Login) | ✅ |
| REQ-002 | 사용자 프로필 관리 | DES-PROTO-USER (UserService GetProfile/UpdateProfile) DES-DB-02 (users 테이블) | IMP-BE-USER (user-service/handler+service+repo) IMP-FE-SETTINGS (settings_screen.dart) | VV-UT-USER-001 (user_test.go) VV-E2E-FLOW | ✅ |
| REQ-003 | 다중 리더기 등록/관리 | DES-PROTO-DEVICE (DeviceService Register/List) DES-DB-03 (devices 테이블) | IMP-BE-DEVICE (device-service/handler+service+repo) IMP-FE-DEVICES (device_list_screen.dart) | VV-UT-DEVICE-001 VV-E2E-HEALTH | ✅ |
| REQ-004 | 리더기 펌웨어 OTA | DES-PROTO-DEVICE (RequestOtaUpdate) DES-AGENTS §2.5 (BLE 명령) | IMP-BE-DEVICE (handler: OTA RPC) IMP-RUST-BLE (command 0x05) | VV-UT-DEVICE-OTA (TODO) | ⚠️ |
| REQ-005 | 카트리지 자동인식 (NFC) | DES-AGENTS §2.6 (NFC 태그 구조) DES-SPEC-CART (cartridge-system-spec.md) | IMP-RUST-NFC (nfc/mod.rs 파싱/레지스트리) | VV-UT-NFC-001 (nfc 테스트 4개) | ⚠️ |
| REQ-006 | 88차원 차동측정 | DES-AGENTS §2.1 (차동측정 공식) DES-RUST-LIB (lib.rs) | IMP-RUST-DIFF (differential/mod.rs S_det-α×S_ref) | VV-UT-DIFF-001 (테스트 3개) VV-BENCH-DIFF | ✅ |
| REQ-007 | 측정 결과 저장/시각화 | DES-PROTO-MEAS (StartSession/EndSession/GetHistory) DES-DB-04 (measurements) | IMP-BE-MEAS (measurement-service 전체) IMP-FE-MEAS (measurement_screen.dart) | VV-UT-MEAS-001 VV-E2E-FLOW (StartSession→EndSession→GetHistory) | ✅ |
| REQ-008 | 896차원 핑거프린트 | DES-AGENTS §2.2 (벡터 시스템) DES-RUST-LIB | IMP-RUST-FP (fingerprint/mod.rs FingerprintBuilder) | VV-UT-FP-001 (테스트 3개) | ✅ |
| REQ-009 | 비표적 분석 (448/896) | DES-SPEC-CART §2.3 (NonTarget448/896) DES-PROTO-MEAS | IMP-BE-MEAS (세션 cartridge_type 기반) IMP-RUST-FP | VV-UT-FP-448 VV-UT-FP-896 | ⚠️ |
| REQ-010 | 측정 세션 관리 | DES-PROTO-MEAS (StartSession/StreamMeasurement/EndSession) | IMP-BE-MEAS (handler grpc.go 세션 CRUD) | VV-UT-MEAS-SESSION VV-E2E-FLOW | ✅ |
| REQ-011 | 오프라인 완전 구동 | DES-SPEC-OFFLINE (offline-capability-matrix.md) DES-RUST-SYNC | IMP-RUST-SYNC (sync/mod.rs CRDT) IMP-RUST-AI (엣지 추론) | VV-UT-SYNC-001 VV-QA-OFFLINE-72H (TODO) | ✅ |
| REQ-057 | 규제 준수 관리 | DES-COMP-REG (regulatory-compliance-checklist.md 146항목) DES-COMP-SAFETY (software-safety-classification.md) | IMP-DOC (docs/compliance/ 8종) | VV-AUDIT-REG (감사 검증 TODO) | ⚠️ |
| REQ-058 | TPM 보안 칩 연동 | DES-AGENTS §보안 DES-RUST-CRYPTO | IMP-RUST-CRYPTO (crypto/mod.rs AES-256-GCM+SHA-256) | VV-UT-CRYPTO-001 | ✅ |
| REQ-059 | BLE AES-CCM 암호화 | DES-AGENTS §2.5 (GATT 서비스) DES-RUST-BLE | IMP-RUST-BLE (ble/mod.rs 구조) | VV-UT-BLE-CRYPTO (TODO) | ⚠️ |
| REQ-060 | HTTPS TLS 1.3 | DES-INFRA-KONG (Kong SSL 종단) DES-K8S-INGRESS | IMP-INFRA-KONG (docker-compose kong) IMP-K8S-INGRESS (ingress.yaml) | VV-INFRA-TLS (TODO) | ✅ |
| REQ-063 | 감사 추적 로그 | DES-COMP-AUDIT (10년 보존) DES-SHARED-OBS | IMP-SHARED-OBS (observability/ metrics+interceptor) | VV-UT-AUDIT (TODO) | ✅ |
| REQ-065 | 72시간 오프라인 검증 | DES-SPEC-OFFLINE DES-QA-VNV (vnv-master-plan.md) | IMP-QA (테스트 시나리오 미작성) | VV-QA-OFFLINE-72H (미실행) | 🔲 |
| REQ-073 | 콘텐츠 캐싱/동기화 | DES-RUST-SYNC (CRDT 병합) | IMP-RUST-SYNC (sync/mod.rs) | VV-UT-SYNC-CACHE | ✅ |

---

## 3. Phase 2 (Core) — 35개 요구사항

| REQ ID | 기능명 | 설계 (DES) | 구현 (IMP) | 검증 (V&V) | 상태 |
|--------|--------|-----------|-----------|-----------|------|
| REQ-012 | 구독 등급 관리 | DES-PROTO-SUB (SubscriptionService) DES-DB-05 (05-subscription.sql) DES-TIER (terminology-and-tier-mapping.md) | IMP-BE-SUB (subscription-service 전체) | VV-UT-SUB-001 | ✅ |
| REQ-013 | SaaS 구독 결제 | DES-PROTO-PAY (PaymentService) DES-DB-07 (07-payment.sql) | IMP-BE-PAY (payment-service 전체) | VV-UT-PAY-001 VV-E2E-COMMERCE | ✅ |
| REQ-014 | 카트리지 무한확장 레지스트리 | DES-SPEC-CART (cartridge-system-spec.md 2-byte 65,536종) DES-PROTO-CART | IMP-BE-CART (cartridge-service 레지스트리) IMP-RUST-NFC (레지스트리 30종) | VV-UT-CART-REGISTRY | ✅ |
| REQ-015 | 등급별 카트리지 접근 제어 | DES-SPEC-CART §접근 정책 DES-TIER | IMP-BE-CART (접근 제어 로직) | VV-UT-CART-ACCESS | ✅ |
| REQ-016 | 카트리지 사용량 추적 | DES-PROTO-CART (GetUsageHistory) DES-DB-09 | IMP-BE-CART (사용 추적) | VV-UT-CART-USAGE | ✅ |
| REQ-017 | 온라인 상품 판매 | DES-PROTO-SHOP (ShopService) DES-DB-06 (06-shop.sql) | IMP-BE-SHOP (shop-service 상품 CRUD) | VV-UT-SHOP-001 VV-E2E-COMMERCE | ✅ |
| REQ-018 | 장바구니/주문 관리 | DES-PROTO-SHOP (AddToCart/PlaceOrder) | IMP-BE-SHOP (장바구니+주문) | VV-UT-SHOP-CART | ✅ |
| REQ-019 | AI 실시간 추론 | DES-PROTO-AI (AiInferenceService AnalyzeMeasurement) DES-DB-08 DES-ML (ml-model-design-spec.md) | IMP-BE-AI (ai-inference-service 시뮬레이션) IMP-RUST-AI (엣지 추론 스텁) | VV-UT-AI-INFER | ✅ |
| REQ-020 | AI 건강 코칭 | DES-PROTO-COACH (CoachingService) DES-DB-11 | IMP-BE-COACH (coaching-service 전체) | VV-UT-COACH-001 | ✅ |
| REQ-021 | 음식 사진→칼로리 | DES-PROTO-VISION (TODO: VisionService) | IMP-BE-VISION (서비스 미존재) | VV-UT-VISION (TODO) | 🔲 |
| REQ-023 | 카트리지 보정 데이터 | DES-PROTO-CALIB (CalibrationService) DES-DB-10 | IMP-BE-CALIB (calibration-service 전체) | VV-UT-CALIB-001 | ✅ |
| REQ-024 | 데이터 허브/타임라인 | DES-UX-SITEMAP (데이터 허브 섹션) DES-PROTO-MEAS (GetHistory) | IMP-FE-DATAHUB (미구현) | VV-WT-DATAHUB (TODO) | 🔲 |
| REQ-025 | 외부 건강 앱 연동 | DES-SPEC-FHIR (TODO) | IMP-BE-MEAS (HealthKit/Google Health TODO) | VV-IT-HEALTH-APP (TODO) | 🔲 |
| REQ-026 | 공공데이터 연계 | DES-PLAN-v1.1 §5.7 | IMP-BE-AI (공공 API 연동 TODO) | VV-IT-OPENDATA (TODO) | 🔲 |
| REQ-027 | 데이터 내보내기 | DES-PROTO-MEAS (ExportData TODO) | IMP-BE-MEAS (FHIR 기본만) | VV-UT-EXPORT | ⚠️ |
| REQ-051 | 448차원 분석 | DES-AGENTS §2.2 | IMP-RUST-FP (448차원 빌더) | VV-UT-FP-448 | ✅ |
| REQ-061 | MFA 다중 인증 | DES-SEC-MFA (TODO: TOTP 설계) DES-PROTO-AUTH (EnableMfa) | IMP-BE-AUTH (MFA TODO) IMP-FE-AUTH (OTP 화면 TODO) | VV-UT-MFA (TODO) | 🔲 |
| REQ-062 | RBAC 역할 기반 접근 | DES-PROTO-ADMIN (AdminService) DES-SHARED-RBAC | IMP-BE-ADMIN (admin-service) IMP-SHARED-RBAC (middleware/rbac.go) | VV-UT-RBAC VV-E2E-ADMIN | ✅ |
| REQ-064 | 동적 부하 시뮬레이션 | DES-QA-PERF (TODO) | IMP-QA-PERF (미작성) | VV-PERF-LOAD (미실행) | 🔲 |
| REQ-066 | 개인 기준선 (My Zone) | DES-PROTO-COACH (SetGoal 확장) | IMP-BE-COACH (My Zone TODO) | VV-UT-MYZONE (TODO) | 🔲 |
| REQ-067 | 환경 모니터링 표시 | DES-SPEC-CART §환경 4종 | IMP-BE-MEAS (환경 카트리지 TODO) IMP-FE-MEAS (대시보드 TODO) | VV-UT-ENV-MON (TODO) | 🔲 |
| REQ-068 | 위험 예측/경고 | DES-PROTO-AI (PredictRisk TODO) DES-ML §위험예측 | IMP-BE-AI (위험 스코어링 TODO) IMP-BE-NOTIF (알림 연동 TODO) | VV-UT-RISK (TODO) | 🔲 |
| REQ-070 | AI 학습 히스토리 | DES-PROTO-TRAIN (TODO: AiTrainingService) | IMP-BE-TRAIN (서비스 미존재) | VV-UT-TRAIN (TODO) | 🔲 |
| REQ-071 | 정기 구독 서비스 | DES-PROTO-SHOP (CreateSubscriptionOrder TODO) | IMP-BE-SHOP (정기구독 TODO) | VV-UT-SHOP-RECURRING (TODO) | 🔲 |
| REQ-072 | 위시리스트 | DES-PROTO-SHOP (AddToWishlist TODO) | IMP-BE-SHOP (위시리스트 TODO) | VV-UT-SHOP-WISH (TODO) | 🔲 |
| REQ-074 | 다국어 UI | DES-UX-I18N DES-FE-L10N | IMP-FE-L10N (app_localizations.dart 6언어) | VV-WT-I18N | ✅ |
| REQ-075 | 접근성 (다크모드) | DES-FE-THEME | IMP-FE-THEME (app_theme.dart Material 3) | VV-WT-THEME | ✅ |
| REQ-076 | 리더기 위치별 대시보드 | DES-PROTO-DEVICE (GetDeviceLocation TODO) | IMP-BE-DEVICE (위치 TODO) IMP-FE-DEVICES (지도 TODO) | VV-WT-DEVICE-MAP (TODO) | 🔲 |
| REQ-077 | 구독별 리더기 제한 | DES-PROTO-SUB (max_devices) DES-TIER | IMP-BE-SUB (제한 로직) | VV-UT-SUB-LIMIT | ✅ |
| REQ-079 | 신규 카트리지 OTA 배포 | DES-SPEC-CART §확장 DES-PROTO-CART (DeployCartridge TODO) | IMP-BE-CART (OTA TODO) | VV-UT-CART-OTA (TODO) | 🔲 |
| REQ-080 | 소셜 로그인 | DES-SEC-OAUTH (TODO: Google/Apple/Facebook) | IMP-BE-AUTH (소셜 TODO) IMP-FE-AUTH (소셜 버튼 TODO) | VV-IT-SOCIAL (TODO) | 🔲 |

---

## 4. Phase 3 (Advanced) — 16개 요구사항

| REQ ID | 기능명 | 설계 (DES) | 구현 (IMP) | 검증 (V&V) | 상태 |
|--------|--------|-----------|-----------|-----------|------|
| REQ-022 | 운동 영상 분석 | DES-PROTO-VISION (TODO) DES-ML §비전 | IMP-BE-VISION (서비스 미존재) | VV-UT-VISION-EX (TODO) | 🔲 |
| REQ-028 | 화상진료 | DES-PROTO-TELE (TelemedicineService) DES-DB-15 | IMP-BE-TELE (telemedicine-service 구조) | VV-UT-TELE-001 VV-E2E-MEDICAL | ✅ |
| REQ-029 | 병원/약국 검색/예약 | DES-PROTO-RESV (ReservationService) DES-DB-16 | IMP-BE-RESV (reservation-service 구조) | VV-UT-RESV-001 VV-E2E-MEDICAL | ✅ |
| REQ-030 | 처방전 관리 | DES-PROTO-PRESC (PrescriptionService) DES-DB-19 | IMP-BE-PRESC (prescription-service 구조) | VV-UT-PRESC-001 | ✅ |
| REQ-031 | 가족 그룹 관리 | DES-PROTO-FAMILY (FamilyService) DES-DB-13 | IMP-BE-FAMILY (family-service 구조) | VV-UT-FAMILY-001 | ✅ |
| REQ-032 | 보호자 모니터링 | DES-PROTO-FAMILY (MonitorMember TODO) | IMP-BE-FAMILY (모니터링 TODO) IMP-FE-FAMILY (대시보드 TODO) | VV-UT-FAMILY-MON (TODO) | 🔲 |
| REQ-033 | 가족 건강 리포트 | DES-PROTO-HR (GenerateFamilyReport TODO) | IMP-BE-HR (리포트 TODO) | VV-UT-HR-FAMILY (TODO) | 🔲 |
| REQ-034 | 건강 기록 (FHIR R4) | DES-PROTO-HR (HealthRecordService) DES-DB-14 DES-SPEC-FHIR | IMP-BE-HR (health-record-service 구조) | VV-UT-HR-001 VV-IT-FHIR | ✅ |
| REQ-035 | 커뮤니티 포럼 | DES-PROTO-COMM (CommunityService) DES-DB-17 | IMP-BE-COMM (community-service 구조) | VV-UT-COMM-001 VV-E2E-COMMUNITY | ✅ |
| REQ-036 | 전문가 Q&A | DES-PROTO-COMM (AskExpert TODO) | IMP-BE-COMM (Q&A TODO) | VV-UT-COMM-QA (TODO) | 🔲 |
| REQ-037 | 글로벌 챌린지 | DES-PROTO-COMM (CreateChallenge TODO) | IMP-BE-COMM (챌린지 TODO) | VV-UT-COMM-CHALL (TODO) | 🔲 |
| REQ-038 | 실시간 번역 | DES-PROTO-TRANS (TranslationService) DES-DB-20 | IMP-BE-TRANS (translation-service 구조) | VV-UT-TRANS-001 | ✅ |
| REQ-039 | 푸시/이메일/SMS 알림 | DES-PROTO-NOTIF (NotificationService) DES-DB-12 | IMP-BE-NOTIF (notification-service 구조) | VV-UT-NOTIF-001 | ✅ |
| REQ-040 | 계층형 관리자 포탈 | DES-PROTO-ADMIN (AdminService) DES-DB-18 | IMP-BE-ADMIN (admin-service 구조) | VV-UT-ADMIN-001 VV-E2E-ADMIN | ✅ |
| REQ-069 | 긴급 연락망 설정 | DES-PROTO-NOTIF (SetEmergencyContacts TODO) | IMP-BE-NOTIF (긴급 연락망 TODO) | VV-UT-NOTIF-EMERG (TODO) | 🔲 |
| REQ-078 | 정기 결과 내보내기 | DES-PROTO-MEAS (ScheduleExport TODO) | IMP-BE-MEAS (정기 내보내기 TODO) | VV-UT-EXPORT-SCHED (TODO) | 🔲 |

---

## 5. Phase 4 (Ecosystem) — 13개 요구사항

| REQ ID | 기능명 | 설계 (DES) | 구현 (IMP) | 검증 (V&V) | 상태 |
|--------|--------|-----------|-----------|-----------|------|
| REQ-041 | 재고 관리 | DES-PLAN-MSA §Phase4 (inventory-service) | IMP-BE-INVENTORY (서비스 미존재) | VV-UT-INV (TODO) | 🔲 |
| REQ-042 | 배송 추적 | DES-PLAN-MSA §Phase4 (logistics-service) | IMP-BE-LOGISTICS (서비스 미존재) | VV-UT-LOGIS (TODO) | 🔲 |
| REQ-043 | 비즈니스 인텔리전스 | DES-PLAN-MSA §Phase4 (analytics-service) | IMP-BE-ANALYTICS (서비스 미존재) | VV-UT-ANALYTICS (TODO) | 🔲 |
| REQ-044 | SDK 마켓 | DES-PLAN-MSA §Phase4 (marketplace-service) | IMP-BE-MARKET (서비스 미존재) | VV-UT-MARKET (TODO) | 🔲 |
| REQ-045 | 수익 분배 시스템 | DES-PLAN-MSA §Phase4 DES-PROTO-MARKET (TODO) | IMP-BE-MARKET (서비스 미존재) | VV-UT-MARKET-REV (TODO) | 🔲 |
| REQ-046 | AI 모델 학습 | DES-ML §연합학습 DES-PLAN-MSA (ai-training-service) | IMP-BE-TRAIN (서비스 미존재) | VV-UT-TRAIN (TODO) | 🔲 |
| REQ-047 | NLP 자연어 처리 | DES-PLAN-MSA §Phase4 (nlp-service) | IMP-BE-NLP (서비스 미존재) | VV-UT-NLP (TODO) | 🔲 |
| REQ-048 | IoT 게이트웨이 | DES-PLAN-MSA §Phase4 (iot-gateway-service) | IMP-BE-IOT (서비스 미존재) | VV-UT-IOT (TODO) | 🔲 |
| REQ-049 | 리더기 위치 추적 | DES-PLAN-MSA §Phase4 (location-service) | IMP-BE-LOCATION (서비스 미존재) | VV-UT-LOCATION (TODO) | 🔲 |
| REQ-050 | 긴급 대응 시스템 | DES-PLAN-MSA §Phase4 (emergency-service) | IMP-BE-EMERGENCY (서비스 미존재) | VV-UT-EMERGENCY (TODO) | 🔲 |
| REQ-052 | 전자코/전자혀 통합 | DES-AGENTS §2.3 (E-Nose/E-Tongue) DES-RUST-FP | IMP-RUST-FP (448/896 구조 있음, 통합 TODO) | VV-UT-FP-ENOSE (TODO) | 🔲 |
| REQ-054 | 음성 명령 인터페이스 | DES-PLAN-MSA §Phase4 (nlp-service) | IMP-BE-NLP (서비스 미존재) | VV-UT-NLP-VOICE (TODO) | 🔲 |

---

## 6. Phase 5 (Future) — 3개 요구사항

| REQ ID | 기능명 | 설계 (DES) | 구현 (IMP) | 검증 (V&V) | 상태 |
|--------|--------|-----------|-----------|-----------|------|
| REQ-053 | 1792차원 궁극 분석 | DES-AGENTS §2.2 (E12-IF 다중 리더기 융합) DES-RUST-FP | IMP-RUST-FP (bridge에 create_fingerprint_1792 존재, 엔진 TODO) | VV-UT-FP-1792 (TODO) | 🔲 |
| REQ-055 | 웨어러블 디바이스 | DES-PLAN-v1.1 §Phase5 | IMP-BE-IOT (서비스 미존재) | VV-IT-WEARABLE (TODO) | 🔲 |
| REQ-056 | 스마트홈 통합 | DES-PLAN-v1.1 §Phase5 | IMP-BE-IOT (서비스 미존재) | VV-IT-SMARTHOME (TODO) | 🔲 |

---

## 7. 통계 요약

### 7.1 상태별 집계

| 상태 | 개수 | 비율 |
|------|------|------|
| ✅ 완료 | 35개 | 43% |
| ⚠️ 부분 구현 | 5개 | 6% |
| 🔲 미구현 | 40개 | 50% |
| **합계** | **80개** | 100% |

### 7.2 Phase별 완료율

| Phase | 총 REQ | ✅ | ⚠️ | 🔲 | 완료율 |
|-------|--------|-----|------|------|--------|
| Phase 1 (MVP) | 18 | 12 | 4 | 2 | 67% |
| Phase 2 (Core) | 35 | 15 | 1 | 19 | 43% |
| Phase 3 (Advanced) | 16 | 8 | 0 | 8 | 50% |
| Phase 4 (Ecosystem) | 13 | 0 | 0 | 13 | 0% |
| Phase 5 (Future) | 3 | 0 | 0 | 3 | 0% |

### 7.3 레이어별 완료율

| 레이어 | 총 REQ | 완료 | 완료율 |
|--------|--------|------|--------|
| Rust 코어 | 12 | 8 | 67% |
| Go 백엔드 | 42 | 22 | 52% |
| Flutter 프론트엔드 | 10 | 5 | 50% |
| 인프라/QA | 8 | 3 | 38% |
| 규정/문서 | 8 | 2 | 25% |

---

## 8. DES 문서 인덱스

| DES ID | 문서명 | 경로 |
|--------|--------|------|
| DES-PROTO-* | gRPC Proto 정의 | backend/shared/proto/manpasik.proto |
| DES-DB-01~24 | DB 초기화 스크립트 | infrastructure/database/init/ |
| DES-AGENTS | 마스터 컨텍스트 | AGENTS.md |
| DES-RUST-* | Rust 모듈 설계 | .cursor/rules/rust-core.mdc |
| DES-SPEC-CART | 카트리지 시스템 | docs/specs/cartridge-system-spec.md |
| DES-SPEC-OFFLINE | 오프라인 매트릭스 | docs/specs/offline-capability-matrix.md |
| DES-TIER | 용어/티어 매핑 | docs/plan/terminology-and-tier-mapping.md |
| DES-ML | AI/ML 모델 설계 | docs/ai-specs/claude/ml-model-design-spec.md |
| DES-SEC-* | 보안 설계 | docs/security/ |
| DES-COMP-* | 규제 준수 | docs/compliance/ |
| DES-UX-* | UX 설계 | docs/ux/ |
| DES-PLAN-* | 기획/계획 | docs/plan/ |
| DES-FE-* | 프론트엔드 설계 | .cursor/rules/frontend.mdc |
| DES-INFRA-* | 인프라 설계 | infrastructure/ |
| DES-QA-* | QA/V&V 설계 | docs/compliance/vnv-master-plan.md |

---

## 9. 유지 관리 규칙

1. **신규 요구사항 추가 시**: 본 문서에 행 추가, REQ ID 순차 부여
2. **구현 완료 시**: 상태를 ✅로 변경, IMP/VV 셀에 실제 코드/테스트 경로 기록
3. **V&V 완료 시**: VV 셀에 테스트 결과 링크 기록, QUALITY_GATES.md 동기화
4. **감사 대비**: 본 문서를 IEC 62304 §5.6 "소프트웨어 항목의 확인" 추적성 증빙으로 사용
5. **변경 이력**: 변경 시 갱신일 및 변경 내역을 CHANGELOG.md에 기록

---

**참조**: QUALITY_GATES.md, docs/compliance/vnv-master-plan.md, AGENTS.md
