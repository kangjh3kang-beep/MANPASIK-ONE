# 만파식 생태계 상세 구현 계획서

**문서번호**: MPK-IMPL-PLAN-v1.0  
**작성일**: 2026-02-12  
**목적**: 최종 기획안(Blueprint v3.0 + FINAL-MASTER + 3개 상세 명세서)에 따라, Phase 0~5 전체를 **주차(Week)별·일(Day)별·태스크별**로 구체화한 상세 구현 계획을 수립한다.  
**상위 문서**: [FINAL-MASTER-IMPLEMENTATION-PLAN](FINAL-MASTER-IMPLEMENTATION-PLAN.md), [COMPREHENSIVE-SYSTEM-BLUEPRINT-v3.0](COMPREHENSIVE-SYSTEM-BLUEPRINT-v3.0.md)  
**참조 명세**: [AI-ASSISTANT-MASTER-SPEC](AI-ASSISTANT-MASTER-SPEC.md), [MEASUREMENT-ANALYSIS-AI-SPEC](MEASUREMENT-ANALYSIS-AI-SPEC.md), [CARTRIDGE-STORE-SDK-SPEC](CARTRIDGE-STORE-SDK-SPEC.md)

---

## 범례

- **BE**: 백엔드 (Go gRPC 서비스)
- **FE**: 프론트엔드 (Flutter)
- **RC**: Rust 코어 (manpasik-engine)
- **IF**: 인프라 (Docker/K8s/CI)
- **AI**: AI/ML 모델·파이프라인
- **QA**: 테스트·검증
- **DOC**: 문서·규제

---

## Phase 0 — 기반 안정화 (Week 1~2)

### Week 1: DB 마이그레이션 + Proto 정식화

| Day | 태스크 ID | 담당 | 태스크 | 산출물 | 완료 기준 |
| --- | --- | --- | --- | --- | --- |
| D1 | P0-01 | BE | config_metadata 테이블 생성 마이그레이션 SQL 작성 | `infrastructure/database/init/22-config-metadata.sql` | psql로 적용 후 테이블 생성 확인 |
| D1 | P0-02 | BE | config_translations 테이블 생성 마이그레이션 SQL 작성 | `infrastructure/database/init/23-config-translations.sql` | psql로 적용 후 테이블 생성 확인 |
| D2 | P0-03 | BE | admin-service: 인메모리 → PostgreSQL 리포지토리 전환 | `admin-service/internal/repository/postgres/` | 단위 테스트 3개 이상 통과 |
| D2 | P0-04 | BE | ConfigMetadata CRUD 구현 (GetConfigWithMeta, UpdateConfigMeta) | admin-service handler 확장 | E2E: 설정 메타 조회·수정 통과 |
| D3 | P0-05 | BE | Proto 정식 재생성 스크립트 작성 | `backend/scripts/gen_proto.sh` | `make proto` 1회 실행으로 모든 .pb.go 재생성 |
| D3 | P0-06 | BE | admin_config_ext.go 수동 코드 제거, 자동 생성 코드로 대체 | 파일 삭제 + go build 성공 | `go build ./...` 에러 없음 |
| D4 | P0-07 | IF | Gateway·E2E 포트 일치 재검증 | docker-compose.dev.yml 수정 | `make e2e` 전체 통과 |
| D4 | P0-08 | IF | Docker Compose 볼륨·시드 스크립트 갱신 (22, 23번 SQL 포함) | docker-compose.dev.yml | `docker compose up` → 테이블 자동 생성 |
| D5 | P0-09 | QA | Phase 0 게이트: 전체 E2E + 빌드 검증 | 통과 리포트 | health + flow + admin 테스트 모두 PASS |

### Week 2: Rust FFI 복구 + Flutter 테스트 강화

| Day | 태스크 ID | 담당 | 태스크 | 산출물 | 완료 기준 |
| --- | --- | --- | --- | --- | --- |
| D1 | P0-10 | RC | RustBridge.init() 활성화, FFI 바인딩 검증 | `flutter-bridge/src/lib.rs` 수정 | Flutter에서 Rust 함수 호출 성공 |
| D1 | P0-11 | RC | BLE/NFC 스텁 함수 → 시뮬레이션 모드 전환 | ble/mod.rs, nfc/mod.rs | `cargo test` 전체 통과 |
| D2 | P0-12 | RC | DSP 기본 함수 단위 테스트 보강 (filter, fft, normalize) | `src/dsp/mod.rs` 테스트 추가 | 테스트 15개 이상, 커버리지 60%+ |
| D2 | P0-13 | RC | differential/fingerprint 단위 테스트 보강 | 각 mod.rs 테스트 | 테스트 10개 이상 |
| D3 | P0-14 | FE | Flutter 단위 테스트 60개 목표 (현재→60) | `test/` 폴더 | `flutter test` 60개 이상 PASS |
| D3 | P0-15 | FE | CI 파이프라인 Flutter 테스트 필수 추가 | `.github/workflows/ci.yml` 수정 | CI에서 Flutter 테스트 실패 시 빌드 실패 |
| D4 | P0-16 | IF | CI 파이프라인 Rust 테스트 추가 | `.github/workflows/ci.yml` 수정 | CI에서 `cargo test` 실행·통과 |
| D4 | P0-17 | QA | Rust FFI 통합 스모크 테스트 (Flutter→Rust→결과) | 스모크 테스트 스크립트 | 측정 시뮬레이션 → 결과 반환 성공 |
| D5 | P0-18 | QA | **Phase 0 완료 게이트** | 체크리스트 | P0-01~P0-17 전체 PASS, III.1 갭 3건 해결(Config PG, Proto, Rust FFI) |

---

## Phase 1 — 핵심 사용자 경험 (Week 3~6)

### Week 3: 인증·온보딩 + 측정 세션 기반

| Day | 태스크 ID | 담당 | 태스크 | 산출물 | 완료 기준 |
| --- | --- | --- | --- | --- | --- |
| D1 | P1-01 | FE | SplashScreen: 토큰 유효성 검사 → /home 또는 /login 분기 | splash_screen.dart 수정 | 앱 시작 시 자동 분기 동작 |
| D1 | P1-02 | FE | LoginScreen: 이메일·비밀번호 유효성 검사, 에러 메시지 통합 | login_screen.dart 수정 | 로그인 성공/실패 시나리오 통과 |
| D2 | P1-03 | FE | RegisterScreen: 이메일 중복 검사, 비밀번호 강도 표시, 약관 동의 | register_screen.dart 신규 | 가입 → 자동 로그인 → /home |
| D2 | P1-04 | BE | AuthService: 이메일 중복 검사 RPC 추가, 약관 동의 필드 | auth-service handler 확장 | 단위 테스트 통과 |
| D3 | P1-05 | RC | NFC 카트리지 인식 플로우: ReadCartridge → CartridgeInfo 반환 | nfc/mod.rs | 시뮬레이션 모드에서 카트리지 정보 반환 |
| D3 | P1-06 | BE | CartridgeService.ValidateCartridge: 타입·잔여 횟수·보정 계수 | cartridge-service 확장 | RPC 호출 → 보정 계수 반환 |
| D4 | P1-07 | BE | MeasurementService.StartSession: session_type, cartridge_uid, device_id, concept_id | measurement-service 확장 | 세션 생성 → session_id 반환 |
| D4 | P1-08 | BE | MeasurementService.StreamMeasurement: 양방향 gRPC 스트림 기반 | measurement-service 확장 | 클라이언트→서버 패킷 송수신 |
| D5 | P1-09 | BE | MeasurementService.EndSession: 세션 종료 → measurement.completed 발행 | measurement-service + Kafka | 이벤트 발행 확인 |
| D5 | P1-10 | QA | 인증 + 측정 세션 E2E 통합 테스트 | e2e/auth_measurement_test.go | 가입→로그인→세션 시작→종료 성공 |

### Week 4: 측정 파이프라인(88차원) + 결과 화면

| Day | 태스크 ID | 담당 | 태스크 | 산출물 | 완료 기준 |
| --- | --- | --- | --- | --- | --- |
| D1 | P1-11 | RC | 88차원 DSP 전처리 파이프라인 (이상치→필터→정규화→FFT) | dsp/mod.rs 확장 | 합성 데이터 입력 → 전처리 결과 검증 |
| D1 | P1-12 | RC | 차동 보정 (S_det − α×S_ref + 오프셋/게인/온도) | differential/mod.rs | 단위 테스트: 보정 전후 SNR 향상 |
| D2 | P1-13 | RC | 핑거프린트 생성 (Basic 88-dim, L2 정규화) | fingerprint/mod.rs | 88-dim 벡터 생성, 코사인 유사도 계산 |
| D2 | P1-14 | RC | 온디바이스 AI: Calibration 모델 (88→88 MLP) TFLite 로드 | ai/mod.rs | TFLite 모델 로드 → 추론 <50ms |
| D3 | P1-15 | RC | 온디바이스 AI: BasicClassifier (88→5클래스) + AnomalyDetector (88→스코어) | ai/mod.rs | 분류 결과 + 이상 스코어 반환 |
| D3 | P1-16 | RC | QualityAssessment (88→3클래스: 좋음/보통/나쁨) | ai/mod.rs | 품질 게이트 판정 |
| D4 | P1-17 | FE | MeasurementScreen: 카트리지 NFC 인식 → 자동 설정 → 실시간 스파크라인 | measurement_screen.dart 대폭 수정 | 카트리지 인식 → 측정 → 결과 전환 |
| D4 | P1-18 | FE | MeasurementResultScreen: 등급·항목 카드·AI 요약·"다음에 할 일" | measurement_result_screen.dart 수정 | 결과 화면 렌더링 완료 |
| D5 | P1-19 | QA | 측정 E2E: NFC 시뮬→세션→스트림→보정→AI→결과 | 통합 테스트 | 전체 파이프라인 1회 완주 성공 |

### Week 5: 홈·데이터허브 + 기기 관리

| Day | 태스크 ID | 담당 | 태스크 | 산출물 | 완료 기준 |
| --- | --- | --- | --- | --- | --- |
| D1 | P1-20 | FE | HomeScreen 리디자인: 건강 점수 원형 게이지, 측정·환경·AI 요약 카드 | home_screen.dart | 홈 화면 Blueprint 6.5 부합 |
| D1 | P1-21 | BE | HealthRecordService.GetHealthSummary: 최근 측정·코칭·환경 통합 | health-record-service 신규/확장 | 요약 데이터 JSON 반환 |
| D2 | P1-22 | FE | DataHubScreen 신규: 요약 탭 (분류별 카드 그리드) | data_hub_screen.dart 신규 | 요약 탭 렌더링 |
| D2 | P1-23 | FE | DataHubScreen: 타임라인 탭 (날짜별 리스트, 측정 N회, 상태 뱃지) | data_hub_screen.dart | 타임라인 스크롤 |
| D3 | P1-24 | FE | DataHubScreen: 트렌드 탭 (기간 선택, 라인/영역 차트) | data_hub_screen.dart + chart 패키지 | 차트 렌더링 |
| D3 | P1-25 | FE | DataHubScreen: 내 기준선 탭 ("My Zone" 막대) | data_hub_screen.dart | 개인 기준선 시각화 |
| D4 | P1-26 | FE | DeviceListScreen 개선: BLE 스캔, RegisterDevice, OTA, 컨셉 할당 | device_list_screen.dart | 기기 추가·설정·OTA 흐름 |
| D4 | P1-27 | FE | app_router.dart: /data, /coach(스텁), /store(스텁) 라우트 추가 | app_router.dart | 탭 전환 동작 |
| D5 | P1-28 | QA | Week 5 통합 검증: 홈→데이터허브→측정→결과→기기 | E2E 시나리오 | 전체 흐름 성공 |

### Week 6: 설정 + Phase 1 마무리

| Day | 태스크 ID | 담당 | 태스크 | 산출물 | 완료 기준 |
| --- | --- | --- | --- | --- | --- |
| D1 | P1-29 | FE | SettingsScreen 개선: 계정·구독·알림·접근성·긴급·약관 6그룹 | settings_screen.dart | 설정 전체 UI 완료 |
| D1 | P1-30 | BE | NotificationService.UpdateNotificationPreferences | notification-service | 알림 설정 저장·조회 |
| D2 | P1-31 | FE | 알림 센터 UI: 홈 상단 아이콘 + 배지 + 드로어 리스트 | notification_drawer.dart 신규 | 알림 목록 표시·읽음 처리 |
| D2 | P1-32 | BE | NotificationService.ListNotifications, MarkAsRead, GetUnreadCount | notification-service | 알림 CRUD 완료 |
| D3 | P1-33 | FE | 테마 라이트/다크/고대비 전환, 글자 크기 설정 | app_theme.dart 확장 | 테마 즉시 전환 |
| D3 | P1-34 | FE | 긴급 대응 설정: 긴급 연락망·119 자동 신고·야간 모드 | emergency_settings.dart 신규 | 긴급 설정 저장 |
| D4 | P1-35 | QA | Phase 1 전체 회귀 테스트 | E2E 전체 실행 | 모든 시나리오 PASS |
| D5 | P1-36 | QA | **Phase 1 완료 게이트** | 체크리스트·릴리스 노트 | 인증+측정+홈+데이터허브+기기+설정 완료 |

---

## Phase 2 — AI·상거래·관리자 (Week 7~10)

### Week 7: AI 주치의·코칭·식단

| Day | 태스크 ID | 담당 | 태스크 | 산출물 | 완료 기준 |
| --- | --- | --- | --- | --- | --- |
| D1 | P2-01 | FE | AiCoachScreen 신규: 대화형 UI, 오늘의 코칭 카드, 이번 주 요약 | ai_coach_screen.dart 신규 | 대화 UI 렌더링 |
| D1 | P2-02 | BE | CoachingService.GenerateCoaching: LLM 기반 코칭 메시지 생성 | coaching-service 확장 | 코칭 메시지 반환 |
| D2 | P2-03 | BE | CoachingService.GetIntegratedSummary: 측정·환경·복약·영양 통합 요약 | coaching-service | JSON 통합 요약 |
| D2 | P2-04 | BE | DietService 신규: LogMeal, ListMeals, GetDailyNutritionSummary, DeleteMeal | diet-service 신규 | 식단 CRUD 완료 |
| D3 | P2-05 | BE | diet_logs, daily_nutrition_summary 마이그레이션 SQL | `infrastructure/database/init/24-diet.sql` | 테이블 생성 |
| D3 | P2-06 | AI | AiInferenceService.AnalyzeFoodImage: 이미지→영양소 벡터 | ai-inference-service 확장 | 음식 사진 → 칼로리 결과 |
| D4 | P2-07 | FE | AiCoachScreen: 음식 촬영/갤러리 → 칼로리 분석 → 식단 로그 연동 | ai_coach_screen.dart 확장 | 사진→분석→로그 저장 |
| D5 | P2-08 | BE | CoachingService: SetHealthGoal, CheckGoalAchievement, 보상(10.7 기본) | coaching-service 확장 | 목표 설정·달성 검증 |

### Week 8: 마켓·결제 + 관리자

| Day | 태스크 ID | 담당 | 태스크 | 산출물 | 완료 기준 |
| --- | --- | --- | --- | --- | --- |
| D1 | P2-09 | FE | ShopScreen 신규: 상품 목록·필터·검색, 구독 혜택 요약 | shop_screen.dart 신규 | 상품 그리드 표시 |
| D1 | P2-10 | FE | CartScreen, OrderScreen, PaymentScreen 신규 | 3개 화면 신규 | 장바구니→주문→결제 흐름 |
| D2 | P2-11 | BE | ShopService: ListProducts, GetProduct, AddToCart, GetCart | shop-service 구현 | 상품·장바구니 API |
| D2 | P2-12 | BE | PaymentService: CreatePayment, ConfirmPayment (Toss PG 연동) | payment-service 확장 | 결제 승인·확인 |
| D3 | P2-13 | BE | SubscriptionService: GetSubscription, CheckFeatureAccess, 플랜 변경 | subscription-service 확장 | 구독 관리 완료 |
| D3 | P2-14 | FE | AdminPortal: 설정 카테고리 탭, ConfigCard, ConfigEditDialog | admin 모듈 신규 | 관리자 설정 UI |
| D4 | P2-15 | FE | AdminPortal: LLM 어시스턴트 채팅 UI (3:2 분할) | admin 모듈 확장 | LLM 채팅 동작 |
| D5 | P2-16 | QA | 마켓·결제 E2E: 상품 조회→장바구니→결제→구독 갱신 | E2E 테스트 | 결제 성공 시나리오 통과 |

### Week 9: AI 비서(텍스트) + 측정 448차원 확장

| Day | 태스크 ID | 담당 | 태스크 | 산출물 | 완료 기준 |
| --- | --- | --- | --- | --- | --- |
| D1 | P2-17 | BE | AssistantService 신규: ProcessUserInput(text), 의도 분류(NLU) | assistant-service 신규 | 텍스트 입력 → 의도 반환 |
| D1 | P2-18 | BE | 오케스트레이터: 의도 → gRPC 도구 호출 → 응답 생성 | assistant-service | "혈당 측정 시작해줘" → StartSession 호출 |
| D2 | P2-19 | BE | assistant_sessions, assistant_turns 마이그레이션 SQL | `infrastructure/database/init/25-assistant.sql` | 테이블 생성 |
| D2 | P2-20 | FE | AI 비서 플로팅 버튼 + 텍스트 입력 UI (전역) | assistant_fab.dart 신규 | 앱 전역 플로팅 진입 |
| D3 | P2-21 | RC | 448차원 핑거프린트 (Enhanced): 전자코/전자혀 채널 포함 | fingerprint/mod.rs 확장 | 448-dim 벡터 생성 |
| D3 | P2-22 | AI | EnhancedClassifier (448→15클래스) CNN-1D 모델 훈련·TFLite 변환 | `models/enhanced_classifier_v1.tflite` | 표준 데이터셋 F1≥0.85 |
| D4 | P2-23 | AI | AiInferenceService GPU 연동: 448/896 모델 서빙 | ai-inference-service 확장 | 서버 추론 <500ms |
| D5 | P2-24 | QA | AI 비서 + 448차원 통합 테스트 | E2E 시나리오 | 텍스트 명령으로 측정·조회 성공 |

### Week 10: 전문가 입력 + Phase 2 마무리

| Day | 태스크 ID | 담당 | 태스크 | 산출물 | 완료 기준 |
| --- | --- | --- | --- | --- | --- |
| D1 | P2-25 | FE | ExpertDashboard 신규: 환자 선택, 실시간 변화추이, 위험 신호 | expert_dashboard.dart 신규 | 전문가 대시 UI |
| D1 | P2-26 | BE | HealthRecordService.CreateRecordFromClinician: 전문가 진료 입력 | health-record-service 확장 | 전문가 입력 → 건강기록 저장 |
| D2 | P2-27 | FE | PatientDetail 신규: 종합/측정 이력/트렌드/복약·진료/내보내기 탭 | patient_detail.dart 신규 | 피검자 상세 표시 |
| D2 | P2-28 | BE | 10.11 연동: 전문가 입력 → clinician_input_completed → AI 코칭 반영 | Kafka + coaching-service | 이벤트 체인 확인 |
| D3 | P2-29 | FE | DataHub 식단 카드 추가, 영양 요약 차트 | data_hub_screen.dart 확장 | 식단·영양 탭/카드 |
| D4 | P2-30 | QA | Phase 2 전체 회귀 테스트 | E2E 전체 | 모든 시나리오 PASS |
| D5 | P2-31 | QA | **Phase 2 완료 게이트** | 체크리스트 | AI 코치+마켓+관리자+비서+전문가 완료 |

---

## Phase 3 — 의료·커뮤니티·가족·목적별 (Week 11~14)

### Week 11: 의료(예약·화상·처방) + 실시간 번역

| Day | 태스크 ID | 담당 | 태스크 | 산출물 | 완료 기준 |
| --- | --- | --- | --- | --- | --- |
| D1 | P3-01 | FE | MedicalScreen 신규: 예약/화상/처방/건강공유 4탭 | medical_screen.dart 신규 | 탭 전환 UI |
| D1 | P3-02 | BE | ReservationService: SearchFacilities, GetAvailableSlots, CreateReservation | reservation-service 구현 | 예약 CRUD |
| D2 | P3-03 | BE | VideoService: CreateRoom, JoinRoom, TURN/WebRTC 기본 설정 | video-service 구현 | 화상방 생성·참여 |
| D2 | P3-04 | BE | TranslationService: TranslateText, TranslateStream (실시간 번역) | translation-service 신규 | 텍스트/스트림 번역 |
| D3 | P3-05 | BE | PrescriptionService: ListPrescriptions, GetPrescription, 복약 리마인더 | prescription-service 확장 | 처방 조회·리마인더 |
| D3 | P3-06 | FE | 화상진료 UI: WebRTC 연결, 실시간 번역 자막 오버레이 | video_call_screen.dart 신규 | 화상 연결·번역 표시 |
| D4 | P3-07 | FE | 예약 UI: 시설 검색, 슬롯 선택, 예약 확인, 취소 | reservation_screen.dart 신규 | 예약 생성·취소 |
| D5 | P3-08 | QA | 의료 E2E: 예약→화상→처방→복약 리마인더 | E2E 테스트 | 전체 흐름 통과 |

### Week 12: 커뮤니티 + 가족

| Day | 태스크 ID | 담당 | 태스크 | 산출물 | 완료 기준 |
| --- | --- | --- | --- | --- | --- |
| D1 | P3-09 | FE | CommunityScreen 신규: 피드/챌린지/내 글 3탭, 무한 스크롤 | community_screen.dart 신규 | 포럼 UI |
| D1 | P3-10 | BE | CommunityService: CreatePost, ListPosts, Comment, LikePost | community-service 구현 | 게시판 CRUD |
| D2 | P3-11 | BE | CommunityService: ListChallenges, JoinChallenge, 진행률 | community-service 확장 | 챌린지 참여 |
| D2 | P3-12 | FE | 챌린지 UI: 진행률 바, "N일 연속 측정" 카드, 참가 버튼 | challenge_card.dart 신규 | 챌린지 UI |
| D3 | P3-13 | FE | FamilyScreen 신규: 가족 구성원 목록, 초대, 공유 설정 | family_screen.dart 신규 | 가족 관리 UI |
| D3 | P3-14 | BE | FamilyService: CreateFamilyGroup, InviteMember, SetSharingPreferences | family-service 구현 | 가족 CRUD |
| D4 | P3-15 | FE | 가족 건강 대시: 멤버 선택 → DataHub 공유 뷰 | family_health_dash.dart 신규 | 가족 건강 조회 |
| D5 | P3-16 | QA | 커뮤니티+가족 E2E | E2E 테스트 | 게시·챌린지·가족초대 통과 |

### Week 13: 목적별·레고형 + 카트리지 스토어 기본 + 896차원

| Day | 태스크 ID | 담당 | 태스크 | 산출물 | 완료 기준 |
| --- | --- | --- | --- | --- | --- |
| D1 | P3-17 | BE | ConceptService 신규: CreateConcept, ListConcepts, AssignDevice | concept-service 신규 | 컨셉 CRUD |
| D1 | P3-18 | BE | concepts, concept_devices, organization_units 마이그레이션 SQL | `infrastructure/database/init/26-concepts.sql` | 테이블 생성 |
| D2 | P3-19 | FE | 목적별 대시 라우트: /dashboard, /family/health, /location/:id/env, /water, /air, /org | 6개 라우트 + 스텁 화면 | 라우트 전환 동작 |
| D2 | P3-20 | FE | CartridgeStoreScreen 신규: 스토어 홈 (배너·카테고리·추천) | cartridge_store_screen.dart 신규 | 스토어 홈 UI |
| D3 | P3-21 | BE | CartridgeRegistryService.RegisterCartridgeType, GetCartridgeSpec | cartridge-registry-service 신규 | 1st-party 등록 |
| D3 | P3-22 | BE | StoreService.ListStoreItems, GetStoreItem | store-service 신규 | 스토어 리스팅 |
| D4 | P3-23 | RC | 896차원 핑거프린트 (Full): 교차전극 448채널 추가 | fingerprint/mod.rs 확장 | 896-dim 벡터 |
| D4 | P3-24 | AI | FullFusionClassifier (896→30클래스) CGMA-Net 훈련 | `models/full_fusion_v1.onnx` | F1≥0.85 |
| D5 | P3-25 | QA | Phase 3 중간 검증 | E2E | 의료+커뮤니티+가족+스토어+896 |

### Week 14: Milvus + GAF + Phase 3 마무리

| Day | 태스크 ID | 담당 | 태스크 | 산출물 | 완료 기준 |
| --- | --- | --- | --- | --- | --- |
| D1 | P3-26 | IF | Milvus 벡터DB Docker Compose 추가, 컬렉션 생성 | docker-compose.dev.yml | Milvus 기동·접근 |
| D1 | P3-27 | BE | 핑거프린트 벡터 Milvus 인덱싱·ANN 유사도 검색 API | measurement-service 확장 | Top-K 유사 패턴 반환 |
| D2 | P3-28 | RC | GAF 이미지 변환 (896-dim → 이미지) | dsp/mod.rs 또는 feature/mod.rs | GAF 이미지 생성 |
| D2 | P3-29 | FE | DataHub: GAF 시각화 카드, 유사 패턴 비교 뷰 | data_hub_screen.dart 확장 | GAF 이미지 표시 |
| D3 | P3-30 | FE | CartridgeStore: 상세 페이지, 내 카트리지, 구매 | store_item_detail.dart 등 | 스토어 전체 UI |
| D4 | P3-31 | QA | Phase 3 전체 회귀 테스트 | E2E 전체 | 모든 시나리오 PASS |
| D5 | P3-32 | QA | **Phase 3 완료 게이트** | 체크리스트 | 의료+커뮤니티+가족+목적별+스토어+896 완료 |

---

## Phase 4 — 지역 통계·익명 학습·기업 (Week 15~17)

### Week 15: 지역 통계 + 익명 학습 파이프라인

| Day | 태스크 ID | 담당 | 태스크 | 산출물 | 완료 기준 |
| --- | --- | --- | --- | --- | --- |
| D1 | P4-01 | BE | LocationStatsService 신규: GetStatsByLocation, GetAlertsForLocation | location-stats-service 신규 | 위치별 통계 |
| D1 | P4-02 | BE | location_statistics 마이그레이션 SQL + 공공 API 수집 스케줄러 | init/27-location-stats.sql + 크론 | 통계 수집 |
| D2 | P4-03 | BE | 경보 규칙 엔진: 임계치 기반 위험/주의/정보 등급 산출 | location-stats-service | alert 생성·발행 |
| D2 | P4-04 | FE | 홈/DataHub: "내 위치 기반 경보" 블록, 지도 오버레이(선택) | home_screen + data_hub 확장 | 경보 표시 |
| D3 | P4-05 | BE | 익명화 파이프라인: 식별자 제거·준식별자 구간화·K-익명성 | anonymization-service 신규 | 익명 데이터 생성 |
| D3 | P4-06 | IF | 학습 전용 저장소 (별도 스키마/버킷) 설정 | Docker + MinIO | 격리된 저장소 |
| D4 | P4-07 | AI | Federated Learning 파이프라인 기초: 로컬 미세조정 → 가중치 집계 | 파이프라인 스크립트 | 시뮬레이션 1회 실행 |
| D5 | P4-08 | QA | 위치 통계 + 익명화 E2E | 테스트 | 통계 수집→경보→익명화 |

### Week 16: 기업 제공 + 카트리지 스토어 SDK 공개

| Day | 태스크 ID | 담당 | 태스크 | 산출물 | 완료 기준 |
| --- | --- | --- | --- | --- | --- |
| D1 | P4-09 | BE | DataProvisionService 신규: ListDataProducts, RequestDataAccess | data-provision-service | B2B 데이터 제공 |
| D1 | P4-10 | BE | data_products, data_provision_contracts 마이그레이션 SQL | init/28-data-provision.sql | 테이블 생성 |
| D2 | P4-11 | BE | DeveloperService 신규: RegisterDeveloper, CreateApiKey | developer-service | 개발자 등록 |
| D2 | P4-12 | BE | ReviewService 신규: SubmitForReview, 자동 검증 파이프라인 | review-service | 심사 제출·자동 검증 |
| D3 | P4-13 | RC | cdk-cli 기본: init, validate, simulate 명령 | `tools/cdk-cli/` | cdk init → 프로젝트 생성 |
| D3 | P4-14 | RC | 카트리지 시뮬레이터: 가상 리더기 + 합성 신호 | `tools/simulator/` | 시뮬레이션 측정 1회 |
| D4 | P4-15 | DOC | CDK 문서: getting-started, hardware-spec, calibration-guide | `manpasik-cdk/docs/` | 3개 가이드 완성 |
| D5 | P4-16 | QA | **Phase 4 완료 게이트** | 체크리스트 | 지역통계+익명학습+기업+SDK 완료 |

### Week 17: 버퍼 + Phase 4 보강

| Day | 태스크 ID | 담당 | 태스크 | 산출물 | 완료 기준 |
| --- | --- | --- | --- | --- | --- |
| D1-D3 | P4-17 | ALL | 지연 태스크 소화, 버그 수정, III.1 갭 P1 항목 해결 | 패치 | 갭 테이블 갱신 |
| D4 | P4-18 | QA | Phase 4 전체 회귀 | E2E 전체 | 전체 PASS |
| D5 | P4-19 | DOC | Phase 4 릴리스 노트, 갭 갱신, III.1 상태 업데이트 | 문서 갱신 | 갭 3건+ 해결 |

---

## Phase 5 — 음성·품질·확장 (Week 18~19)

### Week 18: 음성 AI 비서 + 1792차원 + 카트리지 스토어 전면

| Day | 태스크 ID | 담당 | 태스크 | 산출물 | 완료 기준 |
| --- | --- | --- | --- | --- | --- |
| D1 | P5-01 | BE | AssistantService: 음성 STT 파이프라인 통합 (Whisper 또는 외부 STT) | assistant-service 확장 | 음성→텍스트→의도 |
| D1 | P5-02 | BE | AssistantService: TTS 응답 (음성 출력) | assistant-service 확장 | 텍스트→음성 반환 |
| D2 | P5-03 | FE | AI 비서 마이크 입력, 음성 응답 재생 UI | assistant_fab.dart 확장 | 음성 대화 동작 |
| D2 | P5-04 | BE | VoiceProfileService: RegisterVoiceProfile, TranslateWithVoice | voice-profile-service 신규 | 음성 프로필 등록 |
| D3 | P5-05 | RC | 1792차원 핑거프린트 (Ultimate): 896(t₀) + 896(t₋₁) | fingerprint/mod.rs 확장 | 1792-dim 벡터 |
| D3 | P5-06 | AI | UltimateDiagnostic (1792→30+α) Dual-Encoder Transformer 훈련 | 서버 모델 | F1≥0.85 |
| D4 | P5-07 | AI | TemporalPredictor (1792 시계열) Bi-GRU+Attention | 서버 모델 | 바이오마커 예측 |
| D4 | P5-08 | BE | RevenueService: GetSalesReport, GetPayoutHistory, 정산 | revenue-service 신규 | 수익 배분·정산 |
| D5 | P5-09 | BE | CartridgeAnalyticsService: GetUsageStats, GetRatings | analytics-service 신규 | 사용 통계 |

### Week 19: 품질·성능·문서 + 최종 게이트

| Day | 태스크 ID | 담당 | 태스크 | 산출물 | 완료 기준 |
| --- | --- | --- | --- | --- | --- |
| D1 | P5-10 | QA | 부하 테스트: NFR P95 <200ms, 100 RPS | k6/Locust 스크립트 | NFR 통과 |
| D1 | P5-11 | IF | 롤백 리허설: 카나리 배포 → 롤백 → 정상 복구 | 리허설 로그 | 5분 이내 롤백 |
| D2 | P5-12 | QA | 보안 점검: OWASP Top 10 체크리스트, 의존성 스캔 | 보안 리포트 | Critical 0건 |
| D2 | P5-13 | QA | 접근성 점검: 스크린 리더·키보드·고대비 | 접근성 리포트 | WCAG 2.1 AA |
| D3 | P5-14 | DOC | SDP(소프트웨어 개발 계획서) 최종안 | SDP 문서 | IEC 62304 §5 |
| D3 | P5-15 | DOC | SRS(소프트웨어 요구사항 명세서) 최종안 | SRS 문서 | 80 REQ 매핑 |
| D4 | P5-16 | DOC | SAD(소프트웨어 아키텍처 문서) 최종안 | SAD 문서 | 서비스·DB·이벤트 |
| D4 | P5-17 | DOC | 추적성 매트릭스 최종 갱신 | plan-traceability-matrix.md | 100% 연결 |
| D5 | P5-18 | QA | **Phase 5 완료 게이트 = 전체 시스템 릴리스 게이트** | 최종 체크리스트 | 전 Phase 기능 PASS, NFR PASS, 보안·접근성·규제 PASS |

---

## 태스크 총계

| Phase | 기간 | 태스크 수 | 핵심 산출물 |
| --- | --- | --- | --- |
| Phase 0 | 2주 (W1~2) | 18 | DB 마이그레이션, Proto 정식화, Rust FFI 복구 |
| Phase 1 | 4주 (W3~6) | 36 | 인증, 88차원 측정 E2E, 홈·데이터허브, 기기, 설정 |
| Phase 2 | 4주 (W7~10) | 31 | AI 코치, 식단, 마켓·결제, 관리자, AI 비서(텍스트), 448차원, 전문가 |
| Phase 3 | 4주 (W11~14) | 32 | 의료, 커뮤니티, 가족, 목적별, 카트리지 스토어, 896차원, Milvus, GAF |
| Phase 4 | 3주 (W15~17) | 19 | 지역 통계, 익명 학습, 기업 제공, SDK(CDK) 공개 |
| Phase 5 | 2주 (W18~19) | 18 | 음성 비서, 1792차원, 카트리지 전면, 품질·보안·규제 |
| **합계** | **19주** | **154** | |

---

## 게이트 체크리스트 (각 Phase 종료 시)

| 항목 | Phase 0 | Phase 1 | Phase 2 | Phase 3 | Phase 4 | Phase 5 |
| --- | --- | --- | --- | --- | --- | --- |
| E2E 테스트 전체 PASS | O | O | O | O | O | O |
| 신규 단위 테스트 추가 | 30+ | 50+ | 40+ | 40+ | 20+ | 10+ |
| III.1 갭 해결 건수 | 3 | 1 | 2 | 1 | 2 | 1 |
| Flutter 테스트 누적 | 60+ | 100+ | 150+ | 200+ | 220+ | 240+ |
| 빌드 성공 (BE+FE+RC) | O | O | O | O | O | O |
| 보안 Critical 0건 | O | O | O | O | O | O |
| 릴리스 노트 작성 | O | O | O | O | O | O |
| NFR 부하 테스트 | - | - | - | - | - | O |
| 규제 문서 (SDP/SRS/SAD) | - | - | - | - | - | O |

---

## 참조

- [FINAL-MASTER-IMPLEMENTATION-PLAN](FINAL-MASTER-IMPLEMENTATION-PLAN.md): Phase 0~5 개요, III 갭 분석, IV 로드맵
- [COMPREHENSIVE-SYSTEM-BLUEPRINT-v3.0](COMPREHENSIVE-SYSTEM-BLUEPRINT-v3.0.md): Part 1~11 전체 기획
- [MEASUREMENT-ANALYSIS-AI-SPEC](MEASUREMENT-ANALYSIS-AI-SPEC.md): 측정 파이프라인 상세
- [CARTRIDGE-STORE-SDK-SPEC](CARTRIDGE-STORE-SDK-SPEC.md): 카트리지 스토어 상세
- [AI-ASSISTANT-MASTER-SPEC](AI-ASSISTANT-MASTER-SPEC.md): AI 비서 상세
- [MASTER-DOCUMENT-INDEX](MASTER-DOCUMENT-INDEX.md): 전체 42개 문서 인덱스
