# 만파식 생태계 최종 통합 구현·구축 계획안 (Final Master Implementation Plan)

**문서번호**: MPK-FINAL-MASTER-v1.0  
**작성일**: 2026-02-12  
**목적**: 전체 기능·페이지를 상호유기적으로 연결하고, 모세혈관 수준의 세부 구현 기획·구축 계획을 수립. 갭 분석·미래 확장성·편의성·신규성·진보성 강화 반영.  
**상위 문서**: [COMPREHENSIVE-SYSTEM-BLUEPRINT-v3.0](COMPREHENSIVE-SYSTEM-BLUEPRINT-v3.0.md)

---

## I. 전체 만파식 생태계 상호연결도 및 추적성

### I.1 마스터 연결 구조

- **클라이언트**: Flutter(iOS/Android/Web/Desktop) → Kong(8090) → Keycloak(인증) → Gateway → 22개 gRPC 서비스.
- **데이터**: PostgreSQL(사용자·구독·주문·건강기록·처방·커뮤니티·관리자 등), TimescaleDB(측정 시계열), Milvus(핑거프린트), Redis(세션·캐시), Kafka(17개 토픽 이벤트), Elasticsearch(로그·검색), MinIO(파일).
- **외부**: Toss PG, FCM, HealthKit/Health Connect, TURN/STUN, 공공 API(대기질·수질·감염병), STT/TTS·음성복제 API.
- **Rust 코어**: manpasik-engine(차동·핑거프린트·BLE·NFC·암호화)·flutter-bridge → FFI 활성화 시 Flutter에서 직접 호출.

### I.2 기능 ↔ 페이지 ↔ API ↔ 이벤트 ↔ 데이터 추적성 매트릭스

| 기능 영역 | 주요 페이지/경로 | 핵심 gRPC/API | 주요 이벤트 | DB/저장소 |
| --- | --- | --- | --- | --- |
| 인증 | Splash, Login, Register | AuthService(Register,Login,Refresh,Logout,ValidateToken) | user.registered | users, refresh_tokens |
| 본인 건강 요약 | Home, DataHub(요약/타임라인/트렌드) | User, Measurement, Coaching, HealthRecord(GetHealthSummary) | measurement.completed | measurements, health_records, coaching_messages |
| 측정 | Measurement, MeasurementResult | MeasurementService(StartSession,Stream,EndSession,GetHistory), Device, Cartridge | measurement.completed, session.started/ended | measurements, measurement_sessions |
| AI 주치의·코칭 | AiCoach | CoachingService(SetGoal,GenerateCoaching,Reports,Recommendations), AiInference | measurement.completed, user.profile.updated | health_goals, coaching_messages |
| 전문가 진료 입력 | ExpertDashboard, PatientDetail | HealthRecord(CreateRecordFromClinician), Prescription | clinician_input_completed | health_records, prescriptions |
| 식단·칼로리 | AiCoach, DataHub(영양) | DietService(LogMeal,ListMeals,GetDailyNutritionSummary), AiInference(AnalyzeFoodImage) | meal.logged, nutrition.daily.summarized | diet_logs, daily_nutrition_summary |
| 마켓·결제 | Shop, Cart, Order, Payment | Shop, Payment, Subscription, Cartridge(CheckAccess) | payment.completed/failed, subscription.changed | products, orders, payments, subscriptions |
| 예약·화상·처방 | Medical | Reservation, Video, Prescription, HealthRecord | reservation.created, prescription.created | reservations, video_rooms, prescriptions |
| 커뮤니티 | Community | Community(Post,Comment,Challenge), Translation | community.post.created | posts, comments, challenges |
| 가족 | Family, 가족 건강 대시 | FamilyService(Invite,Sharing,GetSharedHealthData) | family.member.invited, family.sharing.updated | family_groups, family_members, sharing_preferences |
| 기기·카트리지 | Devices, 측정(자동감지) | DeviceService, CartridgeService, Calibration | device.registered, cartridge.verified/depleted | devices, cartridge_usage_log |
| 설정·알림 | Settings | User, Subscription, Notification(Preferences) | notification.preferences.updated, profile.updated | user_profiles, notification_preferences |
| 관리자 | Admin Portal, Admin Settings, LLM Assistant | Admin(ListConfigs,GetConfigWithMeta,BulkSet,Audit), AiInference(ConfigSession) | config.changed | system_configs, config_metadata, audit_logs |
| 목적별·컨셉 | 대시(본인/가족/위치/수질/공기/기업) | ConceptService, OrganizationService, LocationStats | location.alert.raised | concepts, concept_devices, organization_units, location_statistics |
| 익명 학습·기업 제공 | (백엔드 파이프라인) | DataProvisionService, 익명화 파이프라인 | (내부) | 익명 저장소, data_products, data_provision_contracts |
| AI 비서·주치의 | 전역(플로팅/홈/AiCoach) | AssistantService(ProcessUserInput), STT/TTS, 오케스트레이터→기존 gRPC | (세션 내) | assistant_sessions, assistant_turns |
| 측정·분석·AI 확장 | Measurement, MeasurementResult, DataHub | MeasurementService(StartSession/Stream/End), AiInferenceService, 전처리→차동→특징→핑거프린트→AI | measurement.completed | measurements, measurement_sessions, Milvus(핑거프린트), TimescaleDB(시계열) |
| 카트리지 스토어·SDK | /store, /store/item/:id, /store/my-cartridges, developer.manpasik.com | CartridgeRegistryService, DeveloperService, StoreService, ReviewService, RevenueService, CartridgeAnalyticsService | cartridge.registered, review.approved, purchase.completed | cartridge_types, developers, store_items, review_submissions, revenue_transactions, cartridge_usage_stats |

### I.3 크로스 참조 — 문서 내 상세 위치

- **인증·스플래시·로그인·회원가입**: 본 문서 II.1.
- **홈·데이터허브·측정·측정결과**: Blueprint 6.5, 3.2, 10.2.
- **AI 코치·식단·칼로리**: Blueprint 10.2, 10.5, 10.6.
- **전문가 입력·피검자 AI**: Blueprint 10.11.
- **마켓·결제·구독**: Blueprint 3.5, Part 5.1.
- **예약·화상·처방·의료**: Blueprint 3.7, 10.4, 6.6.
- **커뮤니티·번역**: Blueprint 10.4, 3.6.
- **가족**: Blueprint 3.9, 10.13.
- **기기·카트리지·측정 편의**: Blueprint 10.12, 3.8.
- **설정·알림**: 본 문서 II.10.
- **관리자·설정·LLM**: Blueprint 3.11, SPRINT2-AS-BACKEND-DETAIL.
- **목적별·레고형·지역 통계**: Blueprint 10.12, 10.13.
- **익명 학습·기업 제공**: Blueprint 10.9, 10.10.
- **AI 비서·주치의(텍스트·음성 명령)**: Blueprint 10.14, [AI-ASSISTANT-MASTER-SPEC.md](AI-ASSISTANT-MASTER-SPEC.md).
- **측정·분석·AI 확장(88~1792차원 파이프라인)**: Blueprint 10.15, [MEASUREMENT-ANALYSIS-AI-SPEC.md](MEASUREMENT-ANALYSIS-AI-SPEC.md).
- **카트리지 스토어·개발자 SDK(오픈 생태계)**: Blueprint 10.16, [CARTRIDGE-STORE-SDK-SPEC.md](CARTRIDGE-STORE-SDK-SPEC.md).

---

## II. 기능·페이지별 세부 기획 보강 (전체 동일 수준)

아래는 Blueprint에서 상대적으로 간략했던 영역을 6.5·10.11 수준으로 보강한 요약. 구현 시 Blueprint 해당 Part와 본 절을 함께 참조.

### II.1 인증·스플래시·로그인·회원가입

- **설계 원칙**: 최소 입력(이메일·비밀번호·표시명), 단계별 검증(실시간 형식 검사), 보안(비밀번호 강도·저장 불가), 접근성(라벨·에러 메시지 명확).
- **스플래시**: 앱 로고·브랜드 1~2초 → 토큰 유효성 검사(ValidateToken 또는 로컬 만료) → 유효 시 /home, 만료/없음 시 /login. 오프라인 시 마지막 캐시 프로필로 진입 가능 옵션.
- **로그인**: 이메일·비밀번호 필드, “로그인 유지” 체크(Refresh 토큰 연장 정책). 실패 시 “이메일 또는 비밀번호를 확인해주세요” 통합 메시지(보안). 2FA 있으면 다음 단계 입력 화면.
- **회원가입**: 이메일(중복 검사)·비밀번호(강도 표시)·표시명. 약관·개인정보·건강데이터 AI 활용(선택) 동의 체크. 가입 완료 시 자동 로그인 → /home 또는 온보딩 1회.
- **데이터**: AuthService.Register → users, refresh_tokens. manpasik.user.registered 발행 → NotificationService(환영 메일/푸시).

### II.2 측정 플로우 (진입→세션→스트리밍→결과) 모세혈관

- **진입**: MeasurementScreen. 컨셉(10.13) 선택 시 “이 컨셉에서 사용할 카트리지” 목록 또는 자동감지(10.12.1) 대기.
- **카트리지**: NFC 터치 또는 BLE 리더기 삽입 → ReadCartridge/ValidateCartridge → 타입·잔여 사용 횟수 표시. “이 카트리지로 측정” 확인 → StartSession(session_type, cartridge_uid, device_id, concept_id 선택).
- **세션**: StreamMeasurement 양방향 스트림. 클라이언트는 실시간 패킷 수신 → 미니 스파크라인·진행률 표시. 타임아웃·연결 끊김 시 자동 EndSession 호출·부분 결과 저장 옵션.
- **종료**: EndSession → measurement.completed 발행(payload: session_id, sample_type, biomarker_summary, location_id 선택) → AiInference·Coaching·Notification 구독.
- **결과 화면**: MeasurementResultScreen. 상단 등급·한 줄 요약, 항목 카드(이름·수치·단위·스파크라인·뱃지), “다음에 할 일”(재측정/병원). Blueprint 6.5.2 3) 참조.

### II.3 기기 관리 (Devices) 세부

- **목록**: DeviceService.ListDevices → 카드(이름·모델·마지막 연결·배터리·소속 컨셉). “추가” → BLE 스캔(또는 수동 시리얼 입력).
- **등록**: RegisterDevice(device_id, name, concept_id 선택) → device.registered. OTA: RequestOtaUpdate → 디바이스 상태 스트림으로 진행률.
- **설정**: 기기별 “이 기기를 ○○ 컨셉에 사용” 할당(10.13). “기본 측정 리더기” 지정(자동감지 시 우선).

### II.4 마켓·장바구니·주문·결제 세부

- **마켓**: ShopService.ListProducts(필터: 카테고리·구독 포함 여부). CartridgeService.ListAccessibleCartridges·CheckCartridgeAccess로 “내 구독으로 쓸 수 있는 카트리지” 표시. 상품 카드 → GetProduct → “장바구니”/“바로 구매”.
- **장바구니**: AddToCart, GetCart. 수량 변경·삭제. “주문하기” → CreateOrder(배송지·요청 사항) → Order 생성.
- **결제**: PaymentService.CreatePayment(order_id, method) → PG 리다이렉트 또는 앱 내 결제. 복귀 후 ConfirmPayment → payment.completed → 구독 갱신·알림.
- **페이지**: 상단 “내 구독·이번 달 혜택” 요약. 상품 그리드·필터·검색. 장바구니 FAB·배지. 주문 완료 시 “주문 상세”·“배송 추적” 링크.

### II.5 커뮤니티 세부

- **포럼**: ListPosts(정렬·필터)·CreatePost·GetPost. 댓글(ListComments, CreateComment), 좋아요(LikePost). 실시간 번역(10.4) 옵션.
- **챌린지**: ListChallenges(타입·상태)·JoinChallenge·진행률. “N일 연속 측정” 등 게이미피케이션과 연동(10.7).
- **페이지**: 탭(피드 / 챌린지 / 내 글). 카드(제목·미리보기·번역 토글)·무한 스크롤. 작성 시 마크다운·이미지 첨부.

### II.6 의료(예약·화상·처방) 세부

- **예약**: ReservationService.SearchFacilities(지역·과)·GetAvailableSlots·CreateReservation. “다가오는 예약” 카드·취소(CancelReservation).
- **화상**: VideoService.CreateRoom/JoinRoom. 10.4 실시간 번역·10.8 같은 음성 번역 옵션. 종료 시 EndRoom.
- **처방**: PrescriptionService.ListPrescriptions·GetPrescription. SelectPharmacyAndFulfillment·SendPrescriptionToPharmacy. 복약 리마인더(GetMedicationReminders)와 10.11 피검자 AI 연동.
- **페이지**: MedicalScreen 탭(예약 / 화상 / 처방 / 건강데이터 공유). 의료진에게 “요약 패킷” 공유(동의) 버튼.

### II.7 가족 세부

- **초대**: FamilyService.CreateFamilyGroup·InviteMember(이메일/연락처). RespondToInvitation(수락/거절). ListFamilyMembers·SetSharingPreferences(공유 범위).
- **가족 건강 대시**: 10.13. 멤버 선택 → 해당 멤버의 DataHub·측정 이력·코칭 요약(공유 설정 내). 보호자 알림(ai.risk.detected 시 선택).

### II.8 설정 세부

- **계정**: UserService.GetProfile/UpdateProfile. 프로필 사진·표시명·연락처.
- **구독**: SubscriptionService.GetSubscription. 플랜 변경·해지(CancelSubscription). CheckFeatureAccess로 “이 기능 사용 가능” 표시.
- **알림**: NotificationService.UpdateNotificationPreferences(채널·시간대)·GetNotificationPreferences.
- **접근성**: 테마(라이트/다크/고대비), 글자 크기, 스크린 리더 안내. Part 4.3.
- **긴급 대응**: 6.5.4·Part 4.2. 긴급 연락망·119 자동 신고·안전 모드(야간/독거) 설정.
- **페이지**: SettingsScreen 그룹(계정 / 구독 / 알림 / 접근성 / 긴급 / 약관·개인정보). 각 그룹 내 1~2탭 깊이.

### II.9 관리자 포탈 세부

- **설정 관리**: AdminService.ListSystemConfigs(카테고리)·GetConfigWithMeta·ValidateConfigValue·BulkSetConfigs. Blueprint SPRINT2-AS-BACKEND-DETAIL. 카테고리 탭·ConfigCard·ConfigEditDialog·재시작 필요 배너.
- **LLM 어시스턴트**: AiInferenceService.StartConfigSession·SendConfigMessage·ApplyConfigSuggestion·EndConfigSession. 채팅 UI(3:2 분할)·제안 적용/거부.
- **감사·통계**: GetAuditLog(필터)·GetSystemStats. 사용자/측정/결제 요약 대시.
- **계층형(미구현)**: 10.13 조직 단위 확장. 국가/지역/지점/판매점 노드·권한 상속.

### II.10 알림·푸시 세부

- **수신**: NotificationService.ListNotifications(페이징)·MarkAsRead·GetUnreadCount. FCM 연동(ConfigWatcher로 fcm.server_key 핫리로드).
- **표시**: 홈 상단 “알림 센터” 드로어 또는 아이콘+배지. 리스트(제목·본문·시각·읽음). 탭 시 해당 화면(측정 결과/코칭/주문 등) 딥링크.

### II.11 AI 비서·주치의 세부

- **진입점**: 앱 전역 플로팅 버튼·홈 상단·AiCoach 내 "말로/글로 요청하기". 텍스트 입력 또는 음성(STT)으로 동일 파이프라인.
- **범위**: 사용자 직접 수행 이벤트 전부 — 측정 시작/조회, 건강·코칭·식단, 기기 등록/설정, 마켓·장바구니·결제, 예약·처방, 커뮤니티·가족, 설정·알림·컨셉, 관리자(역할 시). 의도 분류·도구(RPC/API) 매핑·확인 정책(결제·해지 등) 상세는 [AI-ASSISTANT-MASTER-SPEC.md](AI-ASSISTANT-MASTER-SPEC.md).
- **아키텍처**: ProcessUserInput(text/audio) → NLU(의도·엔티티) → 오케스트레이터 → 기존 gRPC/REST 호출 → 완전 문장 응답 → TTS(선택). CoachingService·AiInferenceService·10.2 주치의와 동일 페르소나·인증·RBAC·이벤트 연동.
- **학습·성장**: assistant_sessions·assistant_turns 저장; 개인 맥락·선호 암호화; 익명 집계(의도·성공률)로 10.9 생물형 AI 연동; 피드백 반영·주기적 재학습.

### II.12 측정·분석·AI 확장 (88~1792차원) 세부

- **파이프라인**: ① 카트리지 NFC 인식 → ② 세션 시작(채널 수·주파수·시간 자동 설정) → ③ BLE 원시 데이터 스트리밍 → ④ 전처리(이상치 Z-score>3σ 보간, Band-pass/Notch 필터, 정규화, FFT) → ⑤ 차동 보정(S_det − α×S_ref + 오프셋/게인/온도 + ML 보정) → ⑥ 특징 추출(주파수·시간·통계·임피던스·GAF·자동) → ⑦ 핑거프린트(88/448/896/1792-dim 벡터 L2 정규화) → ⑧ AI 추론(온디바이스 TFLite + 서버 GPU 앙상블) → ⑨ 후처리(TimescaleDB·Milvus·Kafka·코칭·알림).
- **차원 성장**: 88(Phase 1) → 448(Phase 2: 전자코/전자혀) → 896(Phase 3: 완전 융합) → 1792(Phase 5: 시간축 확장). 각 단계별 Rust 모듈·AI 모델·백엔드·프론트엔드 로드맵 명시.
- **AI 확장 기능 10종**: 실시간 이상탐지, 바이오마커 정량, 다중 패널, 핑거프린트 유사도 검색, 시계열 예측, 식품 이미지 분석, 품질 게이트, GAF 시각화, 개인화 보정, 교차 검증.
- **품질 게이트**: SNR≥20dB, 패킷유실<5%, 온도 15~40°C, 자가진단 OK, 배터리>10%.
- **상세**: [MEASUREMENT-ANALYSIS-AI-SPEC.md](MEASUREMENT-ANALYSIS-AI-SPEC.md).

### II.13 카트리지 스토어 & 개발자 SDK 세부

- **스토어 구조**: 앱 내 카트리지 스토어(홈·검색·카테고리·상세·내 카트리지). 카테고리: 건강·호르몬·감염·전자코·전자혀·식품·환경·산업·연구·기타. 구매 → 물리 배송 또는 구독 포함 활성화.
- **오픈 SDK (CDK)**: ManPaSik Cartridge Development Kit — docs(하드웨어·전극·NFC·보정·AI·심사 가이드), tools(cdk-cli·시뮬레이터·HDK), sdk(Rust/Python/Web-API), templates(basic/e-nose/e-tongue/multi-panel/custom), examples.
- **cartridge.toml**: 카트리지·하드웨어·NFC·보정·AI·표시·가격·규제 전 필드 정의.
- **개발 워크플로**: 등록 → 프로젝트 생성 → 설계 → 보정 → AI 모델 → 시뮬레이션 → 패키지(.mpk) → 심사 제출 → 게시 → 판매·정산.
- **심사**: 자동 검증(스키마·NFC 코드·보안) → 기술(전극·DSP·보정 R²≥0.95·AI F1≥0.85) → 안전(검체·안내·경고) → 규제(의료용 시 KFDA/CE/FDA).
- **수익 배분**: 개발자 70% / 만파식 30% (기본). 연 매출 1억 미만 소규모: 85/15. 구독 포함 시 사용 횟수 비례 배분.
- **개발자 등급**: Explorer(무료·SDK만) → Developer($99/년) → Professional($299/년) → Enterprise(맞춤) → Research(무료·교육기관).
- **보안**: 루트 CA → 개발자 HMAC 서명 → NFC 프로비저닝 키. 위조·복제·악성모델 방지.
- **인프라**: 6개 신규 서비스(CartridgeRegistry·Developer·Store·Review·Revenue·Analytics). DB 7개 테이블. 프론트 라우트 6개.
- **상세**: [CARTRIDGE-STORE-SDK-SPEC.md](CARTRIDGE-STORE-SDK-SPEC.md).

---

## III. 갭 분석·보강·미래 강화

### III.1 갭 분석 및 보강 항목

| 갭 | 현재 상태 | 보강 방향 | 우선순위 |
| --- | --- | --- | --- |
| ConfigMetadata/Translation PG | 인메모리만 | config_metadata·config_translations PostgreSQL 구현, 마이그레이션 | P1 |
| Rust FFI | 비활성 | RustBridge.init() 복구, BLE/NFC 스텁→실제 연동 단계적 전환 | P0 |
| Proto 수동 코드 | admin_config_ext.go | protoc 재생성 정식화, 수동 파일 제거 | P1 |
| Flutter market/medical/community/family | 미구현 | Feature 모듈·라우트·II.4~II.7 페이지 구현 | P1 |
| 오프라인 CRDT/동기화 | 문서만 | 오프라인 큐·충돌 해결·연결 시 일괄 전송 설계·검증 | P2 |
| 관리자 계층형 | 미반영 | 10.13 organization_units·역할 상속·관리자 포탈 확장 | P2 |
| 화상 TURN/WebRTC | 정의만 | TURN 서버 설정·WebRTC 연결 플로우·실패 시 폴백 | P2 |
| Diet/식단 로그 API | 언급만 | LogMeal·ListMeals·GetDailyNutritionSummary RPC·테이블 정의 | P1 |
| VoiceProfile·TranslateWithVoice | 제안만 | VoiceProfileService·TranslateWithVoice 스펙·API 확정 | P2 |
| LocationStats·Alert 수집 | 제안만 | 공공 API 연동·location_statistics·alert 규칙 엔진 | P2 |

### III.2 미래 확장성 강화

- **카트리지**: 4-byte 타입 확장·서드파티 마켓(0xF0~0xFD) 온보딩·검증 프로세스.
- **다국가**: tenant_id·locale·통화·규제(CE/FDA/PMDA)별 기능 플래그·동의 템플릿.
- **신규 서비스**: Proto에 새 서비스 추가 시 Gateway·라우트·RBAC 한 곳에서 등록하는 컨벤션.
- **연합학습**: Flower·Secure Aggregation 표준화, “참여 동의” 플로우·보상(포인트) 옵션.

### III.3 편의성 강화

- **단축 플로우**: “측정만 하기”(컨셉·카트리지 자동)·“오늘 코칭만 보기”·“한 번에 결제”(장바구니 없이 1개 상품).
- **오프라인 표시**: 네트워크 끊김 시 상단 배너 “오프라인 — 측정은 저장되고 나중에 동기화돼요”.
- **검색 통합**: 앱 내 전역 검색(설정·코칭·측정 이력·상품·커뮤니티) 한 입력창.
- **위젯/바로가기**: 홈 화면 위젯(오늘 점수·다음 측정 리마인더)·앱 바로가기(측정·코치).

### III.4 신규성·진보성 강화

- **예측 UX**: morning_prediction·measurement_anticipation·anomaly_response(Part 9). “내일 이 시간 추천 측정”·“이상 패턴 감지 시 차분한 색·안내” 적용.
- **음성·손없이**: 핵심 기능 음성 명령·TTS 요약. 10.8 같은 음성 번역으로 글로벌 진료·커뮤니티.
- **감정·맥락 인식**: “스트레스 높은 날” 감지 시 코칭 톤·UI 색 조정(Part 9).
- **블록체인·무결성(선택)**: 측정·진료 결과 해시 저장·검증 가능성(규제·연구용).

---

## IV. 최종 구현·구축 계획안 (Phase별·모세혈관 수준)

### IV.1 Phase 0 — 기반 안정화 (2주)

- ConfigMetadata/ConfigTranslation PostgreSQL 마이그레이션·시드. admin-service 전환.
- Proto 정식 재생성·admin_config_ext.go 제거. Gateway·E2E 포트 일치 재검증.
- Rust FFI 복구: RustBridge.init() 활성화, 스모크 테스트.
- Flutter 테스트 60개 목표 달성·CI 필수.

### IV.2 Phase 1 — 핵심 사용자 경험 (4주)

- **인증·온보딩**: 스플래시·로그인·회원가입(II.1) 세부 적용. 2FA 옵션.
- **측정 E2E**: 카트리지 자동감지(10.12.1)·StartSession·Stream·EndSession·결과 화면(6.5.2). measurement.completed 연동.
- **측정·분석 기본(II.12)**: 88차원 DSP 전처리·차동 보정·핑거프린트(Basic)·온디바이스 AI(Calibration+BasicClassifier+Anomaly+Quality). MEASUREMENT-ANALYSIS-AI-SPEC Phase 1 범위.
- **홈·데이터허브**: 요약·타임라인·트렌드·내 기준선(6.5.2). GetHealthSummary·Coaching 연동.
- **기기**: ListDevices·RegisterDevice·OTA·컨셉 할당(10.13) 기본.
- **설정**: 프로필·구독·알림·접근성·긴급(II.8).

### IV.3 Phase 2 — AI·상거래·관리자 (4주)

- **AI 주치의**: 10.2 전반(다체액·환경·복약·패턴·영양·진료)·10.11 전문가 입력→피검자 AI. GetIntegratedSummary·CreateRecordFromClinician.
- **코칭·보상**: GenerateCoaching·일/주간 리포트·10.7 보상(CheckGoalAchievement·GrantReward).
- **식단·칼로리**: 10.5·10.6. LogMeal·AnalyzeFoodImage(또는 확장). 식단 카드 DataHub.
- **마켓·결제**: Shop·Cart·Order·Payment(II.4). payment.completed·subscription.changed 연동.
- **관리자**: ListSystemConfigs·BulkSet·ConfigWatcher·LLM 어시스턴트 UI. Audit·Stats.
- **AI 비서 기본(II.11)**: AssistantService(ProcessUserInput 텍스트), NLU 의도 분류, 기존 gRPC 도구 호출 오케스트레이터, assistant_sessions·assistant_turns DB. 텍스트 비서 MVP.
- **측정·분석 확장(II.12)**: 448차원 EnhancedClassifier, 전자코/전자혀 UI, AiInferenceService GPU 연동. MEASUREMENT-ANALYSIS-AI-SPEC Phase 2 범위.

### IV.4 Phase 3 — 의료·커뮤니티·가족·목적별 (4주)

- **예약·화상·처방**: Reservation·Video(TURN)·Prescription(II.6). 실시간 번역(10.4)·같은 음성(10.8) 옵션.
- **커뮤니티**: Post·Comment·Challenge(II.5). 실시간 번역 채팅.
- **가족**: Invite·Sharing·가족 건강 대시(10.13).
- **목적별·레고형**: concepts·concept_devices·organization_units. 대시(본인/가족/위치/수질/공기/기업)·라우트·GetStatsForConcept.
- **측정·분석 확장(II.12)**: 896차원 FullFusionClassifier, Milvus ANN 검색, GAF 시각화. MEASUREMENT-ANALYSIS-AI-SPEC Phase 3 범위.
- **카트리지 스토어 기본(II.13)**: 1st-party 카트리지 등록, CartridgeRegistryService·StoreService 기본, 스토어 홈·상세 UI. CARTRIDGE-STORE-SDK-SPEC Phase 1 범위.

### IV.5 Phase 4 — 지역 통계·익명 학습·기업 제공 (3주)

- **지역 통계·경보**: LocationStatsService·GetStatsByLocation·GetAlertsForLocation·GetPrediction(10.12). 공공 API 연동·location.alert.raised.
- **익명 학습**: 10.9 파이프라인 완성·학습 전용 저장소·모델 배포·감사.
- **기업 제공**: 10.10 제공 분기·DataProvisionService·DTA·카탈로그.
- **카트리지 스토어 확장(II.13)**: 개발자 등록·SDK(CDK) 공개, DeveloperService·ReviewService, cdk-cli·시뮬레이터, 서드파티 카트리지 심사. CARTRIDGE-STORE-SDK-SPEC Phase 2 범위.

### IV.6 Phase 5 — 음성·품질·확장 (2주)

- **음성 복제 번역**: 10.8 VoiceProfileService·TranslateWithVoice(선택 구간). Flutter “내 음성 등록”·“같은 음성 번역” UI.
- **AI 비서 고도화(II.11)**: 음성 STT/TTS 파이프라인 통합, 음성 명령 전 기능 수행, 학습·성장(피드백·익명 집계·재학습). AI-ASSISTANT-MASTER-SPEC 음성 고도화 범위.
- **측정·분석 궁극(II.12)**: 1792차원 UltimateDiagnostic·TemporalPredictor, Dual-Encoder Transformer, GPU 클러스터. MEASUREMENT-ANALYSIS-AI-SPEC Phase 5 범위.
- **카트리지 스토어 전면(II.13)**: RevenueService·CartridgeAnalyticsService, 수익 배분·정산, 개발자 콘솔 전체, 전면 오픈·AI 모델 마켓플레이스. CARTRIDGE-STORE-SDK-SPEC Phase 3~5 범위.
- **품질·성능**: NFR P95·RPS 검증. 부하 테스트·롤백 리허설.
- **문서·규제**: SDP/SRS/SAD 최종안·추적성 매트릭스. 갭(III.1) P0·P1 정리.

### IV.7 반복 검증 사이클

- **스프린트 종료 시**: 해당 Phase 기능 목록 체크·E2E 시나리오 통과·갭 항목 재평가.
- **분기 검토**: 확장성·편의성·신규성(III.2~III.4) 중 미반영 아이디어 1~2건 선택 반영.
- **릴리스 전**: 보안·접근성·규제 체크리스트·최종 승인.

### IV.8 모세혈관 수준 태스크 예시 (측정 플로우)

- MeasurementScreen 진입 → 컨셉 선택 드롭다운 로드(ConceptService.ListConcepts) → “기본: 본인” 선택 시 concept_id=default.
- “측정 시작” 탭 → BLE 어댑터 상태 확인 → 스캔 시작(또는 NFC 리스너 활성화).
- NFC 태그 감지 시 → CartridgeService.ReadCartridge(uid) → 타입·잔여 횟수 표시 → “이 카트리지로 측정” 버튼 활성화.
- BLE 디바이스 연결 → GATT 서비스·캐릭터리스틱 구독 → MeasurementService.StartSession(device_id, cartridge_uid, concept_id, sample_type from cartridge) 호출.
- 스트림 수신 루프: 패킷 파싱 → 로컬 스파크라인 데이터 추가 → 진행률(0~100%) 갱신 → 90% 도달 시 “거의 완료” 문구.
- “측정 완료” 또는 타임아웃 → EndSession(session_id) → 응답 대기 → measurement.completed 수신 대기(또는 폴링 GetMeasurementHistory 최신 1건).
- MeasurementResultScreen으로 전환(session_id 전달) → GetHistory(session_id) 또는 결과 페이로드로 항목별 카드 렌더 → “다음에 할 일” 버튼(재측정/병원) 표시.
- (선택) “AI 코치에게 물어보기” → AiCoachScreen으로 딥링크(session_id) → GenerateCoaching(context_session_id) 호출.

위와 같은 수준으로 **모든 주요 사용자 시나리오**에 대해 “진입 → API 호출 순서 → UI 상태 변화 → 이벤트/다음 화면”을 나열하면 모세혈관 수준 구현 명세가 된다. 다른 플로우(로그인·주문·예약·가족 초대 등)도 동일한 방식으로 II절과 Blueprint를 기반으로 태스크 리스트화 가능.

### IV.9 미진·보강 반복 분석 및 검증

- **매 Phase 종료**: III.1 갭 표에서 해당 Phase와 연관된 갭의 “현재 상태”를 갱신. 해결 시 P0/P1에서 제거.
- **품질 게이트**: QUALITY_GATES.md S1~S6·Phase 통과 시에만 다음 Phase 진입. 실패 시 원인(갭·리소스·일정) 분석 후 보강 계획 수립.
- **사용자 시나리오 검증**: “첫 측정부터 결과 확인까지”, “가입부터 첫 결제까지”, “전문가 입력부터 피검자 코칭 반영까지” 등 E2E 시나리오를 테스트 케이스로 등록·회귀 실행.
- **보강 우선순위**: P0 → P1 → P2. 신규 아이디어(III.3·III.4)는 Phase 5 또는 이후 “확장 스프린트”에 배치.

---

## V. 문서 간 상호 참조 및 완결성

- **본 문서**: 최종 통합 구현·구축 계획, 추적성, 갭·미래 강화, Phase별 모세혈관 태스크.
- **DETAILED-IMPLEMENTATION-PLAN.md**: 본 문서 IV절을 **주차(Week)별·일(Day)별·태스크별 154개**로 구체화한 상세 구현 계획.
- **MASTER-DOCUMENT-INDEX.md**: 전체 42개 기획·설계·명세·검증 문서의 계층 구조·상호 참조 인덱스.
- **Blueprint Part 1~11 (10.1~10.16)**: 아키텍처·플로우·상세 기능·API·이벤트·DB.
- **기능별 상세 명세**: AI-ASSISTANT-MASTER-SPEC, MEASUREMENT-ANALYSIS-AI-SPEC, CARTRIDGE-STORE-SDK-SPEC.
- **Blueprint Part 6.5~6.7**: 일반인/전문가 정보제공·페이지 구성·데이터 배치.
- **SPRINT2-***: 실행 순서·Day별. QUALITY_GATES.md: Stage·Phase 통과 기준.
- **docs/specs/**: 이벤트 스키마·NFR·배포·테스트·카트리지.
- **docs/ux/**: 사이트맵·스토리보드·DESIGN_SYSTEM.

위 문서들을 함께 참조하면 **전체 만파식 생태계 구축을 위한 시스템 구현·구축 계획**이 상호유기적으로 연결된 **최종 기획안 및 상세 기획안**으로 완성되며, 갭 보강·미래 확장성·편의성·신규성·진보성 반영과 모세혈관 수준 구현 계획까지 포함된다.

---

**문서 끝.**
