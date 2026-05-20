# Phase F: 테스트/품질 — Go + Flutter 테스트 강화

## 개요
Phase F는 Go 백엔드와 Flutter 프론트엔드의 테스트 커버리지 GAP을 보강하는 단계입니다.
정밀 실사 결과, Go 35개 서비스 중 4개가 취약(10건 이하), Flutter는 auth·data_hub 영역이 부분 테스트 상태였습니다.

## 사전 실사 결과

### Go 테스트 현황 (Phase F 이전)
- 전체: 35 서비스, 약 997개 테스트 함수
- 취약 서비스 (10건 이하):
  - digital-twin: 5건
  - device: 10건
  - user: 10건
  - voice-profile: 10건

### Flutter 테스트 현황 (Phase F 이전)
- 전체: 약 51개 파일, 약 380건
- 취약 영역:
  - auth: 기본 모델 테스트만 (AuthResult 팩토리 + AuthState 생성)
  - data_hub: 도메인 모델 테스트만 (건강 점수 계산/모니터링 로직 미검증)
  - coach: 테스트 0건 (UI와 로직 강결합)

## 구현 내역

### F-1: Go 취약 서비스 테스트 강화

#### digital-twin-service (5→14, +9건)
**파일**: `backend/services/digital-twin-service/internal/service/twin_test.go`

| 테스트 | 설명 |
|--------|------|
| TestGetTwinState_Exists | 기존 세션 조회 |
| TestGetTwinState_NotFound | 미존재 세션 → nil |
| TestListDeviceTwins | 디바이스별 트윈 목록 |
| TestListDeviceTwins_DefaultLimit | 빈 결과 기본 처리 |
| TestSyncTwin_ServerOverridesClient | 서버 건강 상태 오버라이드 |
| TestSyncTwin_NegativeCUSUM | 음수 CUSUM → recalibrate |
| TestDetermineAction_Drift | drift → recalibrate 액션 |
| TestDetermineAction_Warning | warning → continue 액션 |
| TestSeverityOrder | 심각도 순서 검증 |

#### device-service (10→18, +8건)
**파일**: `backend/services/device-service/internal/service/device_test.go`

| 테스트 | 설명 |
|--------|------|
| TestRegisterDevice_기본_이름_생성 | 시리얼 번호 기반 자동 이름 |
| TestRequestOtaUpdate_체크섬_결정적 | 동일 입력 → 동일 SHA-256 |
| TestListDevices_다른_유저_격리 | 사용자별 데이터 격리 |
| TestRegisterDevice_Kafka이벤트_발행 | Kafka 등록 이벤트 발행 |
| TestUpdateDeviceStatus_Kafka이벤트_발행 | Kafka 상태 변경 이벤트 |
| TestRequestOtaUpdate_이벤트_기록 | OTA 이벤트 로깅 |
| TestUpdateDeviceStatus_연속_변경 | 연속 상태 변경 + 이벤트 누적 |
| mockKafkaPublisher | 테스트용 Kafka mock |

#### user-service (10→18, +8건)
**파일**: `backend/services/user-service/internal/service/user_test.go`

| 테스트 | 설명 |
|--------|------|
| TestUpdateProfile_부분_업데이트_이름만 | 이름만 변경, 나머지 유지 |
| TestUpdateProfile_아바타URL_업데이트 | 아바타 URL 업데이트 |
| TestUpdateProfile_빈_유저ID | 빈 유저 ID 검증 |
| TestGetSubscription_모든_티어_설정 | Free/Basic/Pro/Clinical 4티어 |
| TestGetMaxDevices_기본값_폴백 | 미구독자 Free 기본값 |
| TestUpdateProfile_타임존_업데이트 | 타임존 변경 |
| TestUpdateProfile_모든_언어 | ko/en/zh/ja 4개 언어 |
| TestTierConfig_검증 | 티어별 디바이스 수 정합성 |

#### voice-profile-service (10→21, +11건)
**파일**: `backend/services/voice-profile-service/internal/service/voice_profile_test.go`

| 테스트 | 설명 |
|--------|------|
| TestCreate_EmptyUserID | 빈 유저 ID 검증 |
| TestCreate_EmptyName | 빈 프로필 이름 검증 |
| TestCreate_DefaultLanguage | 기본 언어 ko |
| TestCreate_NameTooLong | 50자 초과 이름 |
| TestGet_EmptyID | 빈 프로필 ID |
| TestList_EmptyUserID | 빈 유저 ID |
| TestDelete_EmptyID | 빈 프로필 ID |
| TestSynthesize_EmptyText | 빈 텍스트 합성 |
| TestCreate_SupportedLanguages | 7개 지원 언어 |
| TestUpdateName_TooLong | 긴 이름 업데이트 |
| TestUpdateStatus_EmptyID | 빈 프로필 ID |

### F-2: Flutter 취약 영역 테스트 강화

#### auth (기존 14→27, +13건)
**파일**: `test/shared/providers/auth_notifier_test.dart` (+6건)

| 테스트 | 설명 |
|--------|------|
| 데모 로그인 isDemo | loginAsDemo → isDemo=true |
| 게스트 isDemo=false | loginAsGuest → isDemo=false |
| socialLogin restClient 없음 | restClient null → false |
| role 기본값 user | 로그인 후 role='user' |
| displayName 이메일 앞부분 | displayName null → 이메일 파싱 |
| 데모→로그아웃→일반 전환 | 상태 전환 시나리오 |

**파일**: `test/shared/providers/auth_state_test.dart` (+7건)

| 테스트 | 설명 |
|--------|------|
| isAdmin admin | role='admin' → isAdmin=true |
| isAdmin super_admin | role='super_admin' → true |
| isAdmin user | role='user' → false |
| isDemo demo-user-id | userId='demo-user-id' → true |
| isDemo 일반 사용자 | userId='user-123' → false |
| role 기본값 | AuthState() → role='user' |
| copyWith role→isAdmin | role 변경 후 isAdmin 반영 |

#### data_hub (기존 8→30, +22건)
**파일**: `test/features/data_hub/domain/data_hub_health_score_test.dart` (신규)

- **건강 점수 계산** (9건): 빈 목록, 정중앙, 경계, 초과, 평균, null 제외, range=0, 50%, 데모
- **모니터링 요약** (5건): 빈 목록, 혼합, 전원 연결, 배터리 20/19% 경계
- **기기 타입 필터링** (4건): All/Bio/Gas/Env 탭
- **경고 기기 필터링** (4건): disconnected, 배터리 부족, 경고 없음

## 검증 결과

| 구분 | 결과 |
|------|------|
| Rust manpasik-engine | 91 PASS |
| Go 35 서비스 | ALL PASS (0 FAIL) |
| Flutter 전체 | 432+ PASS |

## Phase F 전후 비교

| 영역 | Before | After | 증가 |
|------|--------|-------|------|
| Go digital-twin | 5건 | 14건 | +9 |
| Go device | 10건 | 18건 | +8 |
| Go user | 10건 | 18건 | +8 |
| Go voice-profile | 10건 | 21건 | +11 |
| Flutter auth | 14건 | 27건 | +13 |
| Flutter data_hub | 8건 | 30건 | +22 |
| **합계** | 57건 | 128건 | **+71** |

## 변경 파일 목록
1. `backend/services/digital-twin-service/internal/service/twin_test.go` — +9 테스트
2. `backend/services/device-service/internal/service/device_test.go` — +8 테스트 + mockKafkaPublisher
3. `backend/services/user-service/internal/service/user_test.go` — +8 테스트
4. `backend/services/voice-profile-service/internal/service/voice_profile_test.go` — +11 테스트
5. `frontend/flutter-app/test/shared/providers/auth_notifier_test.dart` — +6 테스트
6. `frontend/flutter-app/test/shared/providers/auth_state_test.dart` — +7 테스트
7. `frontend/flutter-app/test/features/data_hub/domain/data_hub_health_score_test.dart` — 신규 22 테스트
