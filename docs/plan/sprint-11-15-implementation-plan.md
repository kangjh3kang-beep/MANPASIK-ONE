# ManPaSik 100% 달성 상세 구현 계획서

> **문서 ID**: MPK-IMPL-PLAN-v1.0
> **작성일**: 2026-02-18
> **기준 문서**: system-completeness-verification-v7.0.md
> **목표**: 모든 지표 100% 달성 (REST 49%→100%, UI 67%→100%, SDK 연동 완료, E2E 테스트 완비)
> **예상 기간**: Sprint 11~15 (5 스프린트)

---

## 목차

1. [전체 로드맵](#1-전체-로드맵)
2. [Sprint 11: Gateway REST 100% + 미구현 화면](#2-sprint-11)
3. [Sprint 12: Flutter 프론트엔드 연동 100%](#3-sprint-12)
4. [Sprint 13: SDK 외부 연동 (PG/FCM/OAuth/WebRTC)](#4-sprint-13)
5. [Sprint 14: 스토리보드 UI 100% + 폴리싱](#5-sprint-14)
6. [Sprint 15: E2E 테스트 + 빌드 검증 + 출시 준비](#6-sprint-15)
7. [작업 의존성 그래프](#7-작업-의존성-그래프)

---

## 1. 전체 로드맵

```
Sprint 11  ▸ 백엔드 REST Gateway 76개 엔드포인트 추가 (49%→100%)
           ▸ Placeholder 7건 실제 gRPC 연동
           ▸ 미구현 Flutter 화면 3개 신규 생성

Sprint 12  ▸ Flutter REST 클라이언트 76개 메서드 추가
           ▸ Repository 구현체 placeholder 해소
           ▸ Flutter ↔ Gateway 완전 연동

Sprint 13  ▸ Toss Payments SDK 연동
           ▸ Firebase FCM 설정
           ▸ flutter_webrtc 활성화
           ▸ OAuth2 (Google/Kakao/Apple) 연동
           ▸ PASS 본인인증 연동

Sprint 14  ▸ 스토리보드 UI 갭 해소 (67%→100%)
           ▸ Lottie 애니메이션 교체
           ▸ 119 에스컬레이션 UI 구현
           ▸ 데이터 허브 내보내기/외부연동 UI

Sprint 15  ▸ E2E 통합 테스트 30개 시나리오
           ▸ Go 빌드 + Flutter analyze 재검증
           ▸ 성능 프로파일링 + 보안 감사
           ▸ 출시 체크리스트 완료
```

---

## 2. Sprint 11: Gateway REST 100% + 미구현 화면

> **목표**: REST 노출률 49% → 100% (83개 → 159개 엔드포인트)
> **작업량**: 76개 REST 핸들러 + Placeholder 수정 7건 + Flutter 화면 3개

### 2.1 Day 1: Health Record + Prescription REST (25개)

#### Task 11-01: Health Record REST 엔드포인트 (13개)

**파일**: `backend/services/gateway/internal/handler/health_record_routes.go` (신규)

```
구현할 엔드포인트:
POST   /api/v1/health-records                              → CreateRecord
GET    /api/v1/health-records/{recordId}                    → GetRecord
GET    /api/v1/health-records                               → ListRecords
PUT    /api/v1/health-records/{recordId}                    → UpdateRecord
DELETE /api/v1/health-records/{recordId}                    → DeleteRecord
POST   /api/v1/health-records/export-fhir                   → ExportToFHIR
POST   /api/v1/health-records/import-fhir                   → ImportFromFHIR
GET    /api/v1/health-records/summary                       → GetHealthSummary
POST   /api/v1/health-records/consents                      → CreateDataSharingConsent
DELETE /api/v1/health-records/consents/{consentId}          → RevokeDataSharingConsent
GET    /api/v1/health-records/consents                      → ListDataSharingConsents
POST   /api/v1/health-records/share-provider                → ShareWithProvider
GET    /api/v1/health-records/access-log                    → GetDataAccessLog
```

**구현 패턴** (모든 핸들러 동일):
```go
func (h *RestHandler) registerHealthRecordRoutes(mux *http.ServeMux) {
    mux.HandleFunc("POST /api/v1/health-records", h.requireAuth(h.handleCreateHealthRecord))
    mux.HandleFunc("GET /api/v1/health-records/{recordId}", h.requireAuth(h.handleGetHealthRecord))
    // ... 13개 등록
}

func (h *RestHandler) handleCreateHealthRecord(w http.ResponseWriter, r *http.Request) {
    if h.healthRecord == nil {
        writeError(w, http.StatusServiceUnavailable, "health-record-service unavailable")
        return
    }
    var body struct {
        UserId     string `json:"user_id"`
        RecordType string `json:"record_type"`
        Data       string `json:"data"`
    }
    if err := readJSON(r, &body); err != nil {
        writeError(w, http.StatusBadRequest, "invalid request body")
        return
    }
    resp, err := h.healthRecord.CreateRecord(r.Context(), &v1.CreateHealthRecordRequest{
        UserId:     body.UserId,
        RecordType: body.RecordType,
        Data:       body.Data,
    })
    if err != nil {
        writeError(w, http.StatusInternalServerError, err.Error())
        return
    }
    writeProtoJSON(w, http.StatusCreated, resp)
}
```

#### Task 11-02: Prescription REST 엔드포인트 (12개)

**파일**: `backend/services/gateway/internal/handler/prescription_routes.go` (신규)

```
구현할 엔드포인트:
POST   /api/v1/prescriptions                                        → CreatePrescription
GET    /api/v1/prescriptions/{prescriptionId}                        → GetPrescription
GET    /api/v1/prescriptions                                         → ListPrescriptions
PATCH  /api/v1/prescriptions/{prescriptionId}/status                 → UpdatePrescriptionStatus
POST   /api/v1/prescriptions/{prescriptionId}/medications            → AddMedication
DELETE /api/v1/prescriptions/{prescriptionId}/medications/{medId}    → RemoveMedication
POST   /api/v1/prescriptions/check-drug-interaction                  → CheckDrugInteraction
GET    /api/v1/prescriptions/reminders                               → GetMedicationReminders
POST   /api/v1/prescriptions/{prescriptionId}/select-pharmacy        → SelectPharmacyAndFulfillment
POST   /api/v1/prescriptions/{prescriptionId}/send-pharmacy          → SendPrescriptionToPharmacy
GET    /api/v1/prescriptions/by-token/{token}                        → GetPrescriptionByToken
PATCH  /api/v1/prescriptions/{prescriptionId}/dispensary-status      → UpdateDispensaryStatus
```

**검증**: `GOWORK=off go build ./...` 통과

---

### 2.2 Day 2: Family + Community + Video REST (22개)

#### Task 11-03: Family REST 엔드포인트 (8개 추가)

**파일**: `backend/services/gateway/internal/handler/rest_handler.go` (기존 수정)

기존 `handleListFamilyGroups` placeholder를 실제 gRPC 연동으로 교체하고 8개 추가:

```
수정:
GET    /api/v1/family/groups                                 → 빈 배열 → ListFamilyMembers 연동

추가:
POST   /api/v1/family/groups                                 → CreateFamilyGroup
GET    /api/v1/family/groups/{groupId}                       → GetFamilyGroup
POST   /api/v1/family/groups/{groupId}/invite                → InviteMember
POST   /api/v1/family/invitations/{invitationId}/respond     → RespondToInvitation
DELETE /api/v1/family/groups/{groupId}/members/{userId}      → RemoveMember
PUT    /api/v1/family/groups/{groupId}/members/{userId}/role → UpdateMemberRole
GET    /api/v1/family/groups/{groupId}/members               → ListFamilyMembers
PUT    /api/v1/family/groups/{groupId}/sharing-prefs         → SetSharingPreferences
```

#### Task 11-04: Community REST 엔드포인트 (8개 추가)

**파일**: `backend/services/gateway/internal/handler/rest_handler.go` (기존 수정)

```
추가:
POST   /api/v1/posts/{postId}/comments                      → CreateComment
GET    /api/v1/posts/{postId}/comments                       → ListComments
POST   /api/v1/challenges                                    → CreateChallenge
GET    /api/v1/challenges/{challengeId}                      → GetChallenge
POST   /api/v1/challenges/{challengeId}/join                 → JoinChallenge
GET    /api/v1/challenges                                    → ListChallenges
GET    /api/v1/challenges/{challengeId}/leaderboard          → GetChallengeLeaderboard
PUT    /api/v1/challenges/{challengeId}/progress             → UpdateChallengeProgress
```

#### Task 11-05: Video REST 엔드포인트 (6개 추가)

**파일**: `backend/services/gateway/internal/handler/video_routes.go` (신규)

```
추가:
POST   /api/v1/video/rooms                                  → CreateRoom
GET    /api/v1/video/rooms/{roomId}                          → GetRoom
POST   /api/v1/video/rooms/{roomId}/end                      → EndRoom
POST   /api/v1/video/rooms/{roomId}/signal                   → SendSignal
GET    /api/v1/video/rooms/{roomId}/participants              → ListParticipants
GET    /api/v1/video/rooms/{roomId}/stats                     → GetRoomStats
```

---

### 2.3 Day 3: Notification + Translation + Telemedicine + 나머지 REST (27개)

#### Task 11-06: Notification REST (5개 추가)

```
POST   /api/v1/notifications                                → SendNotification
POST   /api/v1/notifications/mark-all-read                   → MarkAllAsRead
PUT    /api/v1/notifications/preferences                     → UpdateNotificationPreferences
GET    /api/v1/notifications/preferences                     → GetNotificationPreferences
POST   /api/v1/notifications/send-template                   → SendFromTemplate
```

#### Task 11-07: Translation REST (6개 추가)

```
POST   /api/v1/translations/detect-language                  → DetectLanguage
GET    /api/v1/translations/languages                        → ListSupportedLanguages
POST   /api/v1/translations/batch                            → TranslateBatch
GET    /api/v1/translations/history                          → GetTranslationHistory
GET    /api/v1/translations/usage                            → GetTranslationUsage
POST   /api/v1/translations/realtime                         → TranslateRealtime
```

#### Task 11-08: Telemedicine REST (5개 추가)

```
GET    /api/v1/telemedicine/consultations/{consultationId}   → GetConsultation
GET    /api/v1/telemedicine/consultations                    → ListConsultations
POST   /api/v1/telemedicine/consultations/{id}/start-video   → StartVideoSession
POST   /api/v1/telemedicine/consultations/{id}/end-video     → EndVideoSession
POST   /api/v1/telemedicine/consultations/{id}/rate          → RateConsultation
```

#### Task 11-09: 기존 서비스 잔여 REST (11개 추가)

```
Measurement:
GET    /api/v1/measurements/latest                           → GetLatestMeasurement
GET    /api/v1/measurements/statistics                       → GetStatistics

Device:
GET    /api/v1/devices/{deviceId}                            → GetDevice
PUT    /api/v1/devices/{deviceId}                            → UpdateDevice
DELETE /api/v1/devices/{deviceId}                            → DeleteDevice

Subscription:
PUT    /api/v1/subscriptions/{subscriptionId}/upgrade        → UpgradeSubscription
GET    /api/v1/subscriptions/{userId}/usage                  → GetUsageStats

Payment:
POST   /api/v1/payments/{paymentId}/refund                   → RefundPayment
GET    /api/v1/payments                                      → ListPayments

Coaching:
PUT    /api/v1/coaching/goals/{goalId}/progress              → UpdateGoalProgress
DELETE /api/v1/coaching/goals/{goalId}                       → DeleteGoal
```

---

### 2.4 Day 4: Placeholder 수정 + REST 연동 교체 (7건)

#### Task 11-10: Placeholder → 실제 gRPC 연동

| # | 엔드포인트 | 현재 | 수정 내용 |
|---|-----------|------|-----------|
| 1 | `POST /users/{userId}/avatar` | 하드코딩 URL | UserService.UpdateProfile + 파일 업로드 스트림 |
| 2 | `PUT /users/{userId}/emergency-settings` | UpdateProfile 우회 | 전용 필드 세팅으로 변경 (proto 내 emergency 필드) |
| 3 | `GET /family/groups` | 빈 배열 반환 | Task 11-03에서 해결 |
| 4 | `GET /products/{id}/reviews` | CommunityService 우회 | ShopService에 리뷰 전용 proto 추가 또는 주석 명시 |
| 5 | `POST /products/{id}/reviews` | CommunityService 우회 | 상동 |
| 6 | `POST /ai/food-analyze` | AnalyzeMeasurement 우회 | AI proto에 AnalyzeFood RPC 추가 또는 주석 명시 |
| 7 | `POST /ai/exercise-analyze` | AnalyzeMeasurement 우회 | AI proto에 AnalyzeExercise RPC 추가 또는 주석 명시 |

**판단**: #4~#7은 현재 우회 방식이 의도적인 경우 주석으로 명시하고 문서화. 향후 전용 proto 분리 시 교체.

---

### 2.5 Day 5: 미구현 Flutter 화면 3개 신규 생성

#### Task 11-11: 119 에스컬레이션 진행 화면

**파일**: `frontend/flutter-app/lib/features/settings/presentation/escalation_progress_screen.dart` (신규)

```dart
/// 119 에스컬레이션 4단계 진행 화면
/// 스토리보드: storyboard-emergency-response.md 장면 3
///
/// UI 요소:
/// - 4단계 타이머 프로그레스 (1→본인확인 2→보호자 3→AI음성 4→119)
/// - 각 단계별 남은 시간 카운트다운
/// - GPS 위치 + 건강 데이터 전송 상태
/// - 신고 취소 버튼 (단계 3 이전)
/// - 119 통화 연결 상태 (단계 4)
class EscalationProgressScreen extends ConsumerStatefulWidget { ... }
```

**라우트 추가**: `app_router.dart`에 `/settings/escalation` 추가

#### Task 11-12: 구독 해지/다운그레이드 화면

**파일**: `frontend/flutter-app/lib/features/market/presentation/subscription_cancel_screen.dart` (신규)

```dart
/// 구독 해지/다운그레이드 화면
/// 스토리보드: storyboard-subscription-upgrade.md 장면 4
///
/// UI 요소:
/// - 현재 플랜 정보
/// - 다운그레이드 시 잃는 기능 목록
/// - 해지 사유 선택 (라디오 버튼 5개)
/// - 잔여 기간 안내
/// - 환불 정책 안내
/// - 최종 확인 다이얼로그 (2단계)
class SubscriptionCancelScreen extends ConsumerWidget { ... }
```

**라우트 추가**: `app_router.dart`에 `/market/subscription/cancel` 추가

#### Task 11-13: 데이터 내보내기 화면

**파일**: `frontend/flutter-app/lib/features/data_hub/presentation/data_export_screen.dart` (신규)

```dart
/// 데이터 내보내기 화면
/// 스토리보드: storyboard-data-hub.md 장면 2
///
/// UI 요소:
/// - 내보내기 형식 선택 (PDF/CSV/FHIR)
/// - 기간 선택 (DateRangePicker)
/// - 바이오마커 선택 (체크박스 리스트)
/// - 내보내기 진행 프로그레스
/// - 파일 공유 (share_plus)
class DataExportScreen extends ConsumerWidget { ... }
```

**라우트 추가**: `app_router.dart`에 `/data/export` 추가

---

### Sprint 11 검증 기준

```bash
# 1. Go 빌드 (Gateway 포함 전체)
for svc in gateway auth-service user-service measurement-service device-service \
  subscription-service shop-service payment-service ai-inference-service \
  cartridge-service calibration-service coaching-service notification-service \
  family-service health-record-service telemedicine-service reservation-service \
  community-service admin-service prescription-service translation-service video-service; do
  wsl -d Ubuntu -- bash -c "export PATH=/usr/local/go/bin:/usr/bin:/bin && \
    cd /home/kangjh3kang/Manpasik/backend/services/$svc && GOWORK=off go build ./..."
done

# 2. Flutter analyze
wsl -d Ubuntu -- bash --norc --noprofile -c \
  'export PATH=/home/kangjh3kang/flutter/bin:/usr/local/go/bin:/usr/local/bin:/usr/bin:/bin && \
   cd /home/kangjh3kang/Manpasik/frontend/flutter-app && flutter analyze'

# 3. REST 엔드포인트 카운트 (159개 이상)
grep -c 'mux.HandleFunc' backend/services/gateway/internal/handler/*.go
```

**통과 조건**: Go 빌드 ALL PASS, Flutter analyze 에러 0, REST 159개+

---

## 3. Sprint 12: Flutter 프론트엔드 연동 100%

> **목표**: Flutter Repository 구현체 placeholder 0건, REST 클라이언트 완전 연동
> **작업량**: rest_client.dart 76개 메서드 추가 + Repository 구현체 20개 수정

### 3.1 Day 1: REST 클라이언트 메서드 추가 (76개)

#### Task 12-01: rest_client.dart 확장

**파일**: `frontend/flutter-app/lib/core/services/rest_client.dart`

서비스별 메서드 추가:

```dart
// ─── Health Record (13개) ───
Future<Map<String, dynamic>> createHealthRecord(String userId, String type, String data);
Future<Map<String, dynamic>> getHealthRecord(String recordId);
Future<List<Map<String, dynamic>>> listHealthRecords({String? userId, String? type, int limit = 20, int offset = 0});
Future<Map<String, dynamic>> updateHealthRecord(String recordId, String data);
Future<void> deleteHealthRecord(String recordId);
Future<Map<String, dynamic>> exportToFhir(String userId, String format);
Future<Map<String, dynamic>> importFromFhir(String fhirBundle);
Future<Map<String, dynamic>> getHealthSummary(String userId);
Future<Map<String, dynamic>> createDataSharingConsent(String userId, String providerId, List<String> dataTypes);
Future<void> revokeDataSharingConsent(String consentId);
Future<List<Map<String, dynamic>>> listDataSharingConsents(String userId);
Future<Map<String, dynamic>> shareWithProvider(String recordId, String providerId);
Future<List<Map<String, dynamic>>> getDataAccessLog(String recordId);

// ─── Prescription (12개) ───
Future<Map<String, dynamic>> createPrescription(String consultationId, List<Map<String, dynamic>> medications);
Future<Map<String, dynamic>> getPrescription(String prescriptionId);
Future<List<Map<String, dynamic>>> listPrescriptions({String? userId, String? status});
Future<Map<String, dynamic>> updatePrescriptionStatus(String prescriptionId, String status);
Future<Map<String, dynamic>> addMedication(String prescriptionId, Map<String, dynamic> medication);
Future<void> removeMedication(String prescriptionId, String medicationId);
Future<Map<String, dynamic>> checkDrugInteraction(List<String> medicationIds);
Future<List<Map<String, dynamic>>> getMedicationReminders(String userId, {String? date});
Future<Map<String, dynamic>> selectPharmacy(String prescriptionId, String pharmacyId);
Future<Map<String, dynamic>> sendPrescriptionToPharmacy(String prescriptionId);
Future<Map<String, dynamic>> getPrescriptionByToken(String token);
Future<Map<String, dynamic>> updateDispensaryStatus(String prescriptionId, String status);

// ─── Family (8개) ───
Future<Map<String, dynamic>> createFamilyGroup(String name, String creatorId);
Future<Map<String, dynamic>> getFamilyGroup(String groupId);
Future<Map<String, dynamic>> inviteFamilyMember(String groupId, String phone, String role);
Future<Map<String, dynamic>> respondToInvitation(String invitationId, bool accept);
Future<void> removeFamilyMember(String groupId, String userId);
Future<Map<String, dynamic>> updateMemberRole(String groupId, String userId, String role);
Future<List<Map<String, dynamic>>> listFamilyMembers(String groupId);
Future<void> setSharingPreferences(String groupId, Map<String, dynamic> prefs);

// ─── Community (8개) ───
Future<Map<String, dynamic>> createComment(String postId, String content);
Future<List<Map<String, dynamic>>> listComments(String postId, {int limit = 20, int offset = 0});
Future<Map<String, dynamic>> createChallenge(String title, String description, DateTime startDate, DateTime endDate);
Future<Map<String, dynamic>> getChallenge(String challengeId);
Future<void> joinChallenge(String challengeId);
Future<List<Map<String, dynamic>>> listChallenges({int limit = 20, int offset = 0});
Future<List<Map<String, dynamic>>> getChallengeLeaderboard(String challengeId);
Future<void> updateChallengeProgress(String challengeId, double progress);

// ─── Video (6개) ───
Future<Map<String, dynamic>> createVideoRoom(String consultationId);
Future<Map<String, dynamic>> getVideoRoom(String roomId);
Future<Map<String, dynamic>> endVideoRoom(String roomId);
Future<void> sendVideoSignal(String roomId, String type, String payload);
Future<List<Map<String, dynamic>>> listVideoParticipants(String roomId);
Future<Map<String, dynamic>> getVideoRoomStats(String roomId);

// ─── Notification (5개) ───
Future<Map<String, dynamic>> sendNotification(String userId, String title, String body, String type);
Future<void> markAllNotificationsAsRead(String userId);
Future<void> updateNotificationPreferences(String userId, Map<String, dynamic> prefs);
Future<Map<String, dynamic>> getNotificationPreferences(String userId);
Future<Map<String, dynamic>> sendNotificationFromTemplate(String templateId, String userId, Map<String, String> params);

// ─── Translation (6개) ───
Future<Map<String, dynamic>> detectLanguage(String text);
Future<List<Map<String, dynamic>>> listSupportedLanguages();
Future<List<Map<String, dynamic>>> translateBatch(List<String> texts, String targetLang);
Future<List<Map<String, dynamic>>> getTranslationHistory(String userId);
Future<Map<String, dynamic>> getTranslationUsage(String userId);
Future<Map<String, dynamic>> translateRealtime(String text, String sourceLang, String targetLang);

// ─── Telemedicine (5개) ───
Future<Map<String, dynamic>> getConsultation(String consultationId);
Future<List<Map<String, dynamic>>> listConsultations({String? userId, String? status});
Future<Map<String, dynamic>> startVideoSession(String consultationId);
Future<Map<String, dynamic>> endVideoSession(String consultationId);
Future<void> rateConsultation(String consultationId, int rating, String comment);

// ─── 기존 서비스 잔여 (11개) ───
Future<Map<String, dynamic>> getLatestMeasurement(String userId);
Future<Map<String, dynamic>> getMeasurementStatistics(String userId, {String? period});
Future<Map<String, dynamic>> getDevice(String deviceId);
Future<Map<String, dynamic>> updateDevice(String deviceId, Map<String, dynamic> data);
Future<void> deleteDevice(String deviceId);
Future<Map<String, dynamic>> upgradeSubscription(String subscriptionId, String planId);
Future<Map<String, dynamic>> getSubscriptionUsageStats(String userId);
Future<Map<String, dynamic>> refundPayment(String paymentId, {String? reason});
Future<List<Map<String, dynamic>>> listPayments({String? userId, int limit = 20});
Future<void> updateGoalProgress(String goalId, double progress);
Future<void> deleteGoal(String goalId);
```

---

### 3.2 Day 2-3: Repository 구현체 수정 (20개 파일)

#### Task 12-02: Community Repository - Placeholder 해소

**파일**: `frontend/flutter-app/lib/features/community/data/community_repository_rest.dart`

```dart
// 수정 전: throw UnimplementedError('Comment creation not available via REST yet')
// 수정 후:
@override
Future<Comment> createComment(String postId, String content) async {
  final data = await _rest.createComment(postId, content);
  return Comment.fromJson(data);
}

@override
Future<List<Comment>> getComments(String postId, {int limit = 20, int offset = 0}) async {
  final data = await _rest.listComments(postId, limit: limit, offset: offset);
  return data.map((e) => Comment.fromJson(e)).toList();
}

// getChallenges, joinChallenge 등 6개 메서드 동일 패턴
```

#### Task 12-03: Family Repository - Placeholder 해소

**파일**: `frontend/flutter-app/lib/features/family/data/family_repository_rest.dart`

```dart
// 수정 전: 로컬 하드코딩 데이터
// 수정 후:
@override
Future<FamilyGroup> createGroup(String name) async {
  final data = await _rest.createFamilyGroup(name, _userId);
  return FamilyGroup.fromJson(data);
}

@override
Future<void> acceptInvitation(String invitationCode) async {
  await _rest.respondToInvitation(invitationCode, true);
}

// updateMemberPermission, removeMember, sendMeasurementReminder 등 5개
```

#### Task 12-04: Medical Repository - Placeholder 해소

**파일**: `frontend/flutter-app/lib/features/medical/data/medical_repository_rest.dart`

```dart
// cancelReservation: placeholder → REST 연동
@override
Future<void> cancelReservation(String reservationId) async {
  await _rest.cancelReservation(reservationId);
}

// sendEmergencyAlert: placeholder → REST 연동
@override
Future<void> sendEmergencyAlert(EmergencyAlertRequest request) async {
  await _rest.sendEmergencyAlert(request.userId, request.type, request.location);
}
```

#### Task 12-05: Market Repository - 구독 관리 메서드 추가

**파일**: `frontend/flutter-app/lib/features/market/domain/market_repository.dart`

```dart
// 인터페이스에 추가:
Future<Subscription> getCurrentSubscription(String userId);
Future<void> upgradeSubscription(String planId);
Future<void> downgradeSubscription(String planId);
Future<void> cancelSubscription(String reason);
```

**파일**: `frontend/flutter-app/lib/features/market/data/market_repository_rest.dart`

```dart
// 구현체 추가:
@override
Future<Subscription> getCurrentSubscription(String userId) async {
  final data = await _rest.getSubscription(userId);
  return Subscription.fromJson(data);
}

@override
Future<void> cancelSubscription(String reason) async {
  await _rest.cancelSubscription(_subscriptionId, reason);
}
```

#### Task 12-06: Notification Repository - 확장 메서드 추가

**파일**: `frontend/flutter-app/lib/features/notification/domain/notification_repository.dart`

```dart
// 추가:
Future<void> markAllAsRead(String userId);
Future<NotificationPreferences> getPreferences(String userId);
Future<void> updatePreferences(String userId, NotificationPreferences prefs);
```

#### Task 12-07: Medical/Prescription - 약국 전송 기능 추가

**파일**: `frontend/flutter-app/lib/features/medical/domain/medical_repository.dart`

```dart
// 추가:
Future<void> sendPrescriptionToPharmacy(String prescriptionId, String pharmacyId);
Future<PrescriptionFulfillment> getPrescriptionFulfillment(String prescriptionId);
Future<List<MedicationReminder>> getMedicationReminders(String userId);
```

---

### 3.3 Day 4-5: Provider 연동 + 화면 바인딩

#### Task 12-08: gRPC Provider 확장

**파일**: `frontend/flutter-app/lib/core/providers/grpc_provider.dart`

```dart
// 기존 패턴 유지하며 새로운 Provider 추가:

// Health Record Providers
final healthRecordSummaryProvider = FutureProvider.autoDispose<HealthSummary>((ref) async {
  final repo = ref.watch(healthRecordRepositoryProvider);
  return repo.getHealthSummary();
});

final dataSharingConsentsProvider = FutureProvider.autoDispose<List<DataSharingConsent>>((ref) async {
  final repo = ref.watch(healthRecordRepositoryProvider);
  return repo.listConsents();
});

// Challenge Providers
final challengeListProvider = FutureProvider.autoDispose<List<HealthChallenge>>((ref) async {
  final repo = ref.watch(communityRepositoryProvider);
  return repo.getChallenges();
});

final challengeLeaderboardProvider = FutureProvider.autoDispose.family<List<LeaderboardEntry>, String>((ref, challengeId) async {
  final repo = ref.watch(communityRepositoryProvider);
  return repo.getLeaderboard(challengeId);
});

// Medication Reminders
final medicationRemindersProvider = FutureProvider.autoDispose<List<MedicationReminder>>((ref) async {
  final repo = ref.watch(medicalRepositoryProvider);
  return repo.getMedicationReminders();
});
```

#### Task 12-09: 화면 연동 수정

| 화면 파일 | 수정 내용 |
|-----------|-----------|
| `community_screen.dart` | `_ChallengeTab`에서 `challengeListProvider` 사용 |
| `challenge_screen.dart` | `challengeLeaderboardProvider` 연동 |
| `escalation_progress_screen.dart` | REST 연동 (Sprint 11에서 생성) |
| `subscription_cancel_screen.dart` | `cancelSubscription` 연동 |
| `data_export_screen.dart` | `exportToFhir` 연동 |
| `notification_screen.dart` | `markAllAsRead`, `getPreferences` 연동 |
| `consent_management_screen.dart` | `dataSharingConsentsProvider` 연동 |

---

### Sprint 12 검증 기준

```bash
# 1. Flutter analyze 에러 0
# 2. 모든 Repository 구현체에서 UnimplementedError 0건
grep -rn "UnimplementedError\|throw.*not.*implemented\|placeholder\|빈 배열 반환" \
  frontend/flutter-app/lib/features/*/data/
# → 결과 0건이어야 함

# 3. rest_client.dart 메서드 수 159개+
grep -c "Future<" frontend/flutter-app/lib/core/services/rest_client.dart
```

---

## 4. Sprint 13: SDK 외부 연동

> **목표**: 시뮬레이션 서비스 0건, 실제 SDK 연동 100%
> **전제**: 각 SDK의 계약/키 발급 필요

### 4.1 Day 1-2: Toss Payments SDK

#### Task 13-01: Toss Payments 웹뷰 연동

**파일**: `frontend/flutter-app/lib/core/services/payment_service.dart`

```dart
/// SimulatedPaymentService → TossPaymentService 교체
class TossPaymentService implements PaymentService {
  static const _clientKey = String.fromEnvironment('TOSS_CLIENT_KEY');

  @override
  Future<PaymentResult> requestPayment({
    required String orderName,
    required int amountKrw,
    required String orderId,
  }) async {
    // 1. 토스 결제창 웹뷰 호출
    final result = await Navigator.push(context, MaterialPageRoute(
      builder: (_) => TossPaymentWebView(
        clientKey: _clientKey,
        orderId: orderId,
        orderName: orderName,
        amount: amountKrw,
        successUrl: 'manpasik://payment/success',
        failUrl: 'manpasik://payment/fail',
      ),
    ));

    // 2. 성공 시 서버에 확인 요청
    if (result?.paymentKey != null) {
      await _rest.confirmPayment(result.paymentKey, orderId, amountKrw);
      return PaymentResult(success: true, paymentKey: result.paymentKey);
    }
    return PaymentResult(success: false, errorMessage: result?.errorMessage);
  }
}
```

**신규 파일**: `frontend/flutter-app/lib/shared/widgets/toss_payment_webview.dart`

**pubspec.yaml 추가**: `webview_flutter: ^4.x`

---

### 4.2 Day 2: Firebase FCM

#### Task 13-02: FCM 푸시 알림 연동

**파일**: `frontend/flutter-app/lib/core/services/push_notification_service.dart`

```dart
/// PollingNotificationService → FcmNotificationService 교체
class FcmNotificationService implements PushNotificationService {
  @override
  Future<void> initialize() async {
    await Firebase.initializeApp();
    final messaging = FirebaseMessaging.instance;

    // iOS 권한 요청
    await messaging.requestPermission(alert: true, badge: true, sound: true);

    // 토큰 획득 및 서버 등록
    final token = await messaging.getToken();
    if (token != null) {
      await _rest.registerPushToken(token);
    }

    // 포그라운드 알림 처리
    FirebaseMessaging.onMessage.listen((message) {
      _controller.add(PushNotification.fromFcm(message));
    });

    // 백그라운드 알림 탭 처리
    FirebaseMessaging.onMessageOpenedApp.listen((message) {
      _handleNotificationTap(message);
    });
  }
}
```

**필요 파일**:
- `android/app/google-services.json`
- `ios/Runner/GoogleService-Info.plist`
- `firebase_options.dart` (flutterfire configure)

**pubspec.yaml 추가**: `firebase_core`, `firebase_messaging`

---

### 4.3 Day 3: flutter_webrtc

#### Task 13-03: WebRTC P2P 연결 활성화

**파일**: `frontend/flutter-app/lib/features/medical/presentation/video_call_screen.dart`

```dart
// 주석 해제 및 구현:
Future<void> _initWebRtc() async {
  final config = {
    'iceServers': [
      {'urls': 'stun:stun.l.google.com:19302'},
      {
        'urls': 'turn:turn.manpasik.com:3478',
        'username': _turnUsername,
        'credential': _turnCredential,
      },
    ],
  };

  _peerConnection = await createPeerConnection(config);

  // 로컬 스트림
  _localStream = await navigator.mediaDevices.getUserMedia({
    'audio': true,
    'video': {'facingMode': 'user', 'width': 640, 'height': 480},
  });
  _localStream!.getTracks().forEach((track) {
    _peerConnection!.addTrack(track, _localStream!);
  });

  // 리모트 스트림
  _peerConnection!.onTrack = (event) {
    setState(() => _remoteStream = event.streams[0]);
  };

  // ICE candidate
  _peerConnection!.onIceCandidate = (candidate) {
    _rest.sendVideoSignal(widget.roomId, 'ice-candidate', jsonEncode(candidate.toMap()));
  };

  // Signaling (REST polling 또는 WebSocket)
  _startSignalingPoll();
}
```

**pubspec.yaml 추가**: `flutter_webrtc: ^0.11.x`

---

### 4.4 Day 4: OAuth2 (Google/Kakao/Apple)

#### Task 13-04: 소셜 로그인 연동

**파일**: `frontend/flutter-app/lib/features/auth/data/auth_repository_impl.dart`

```dart
// Google Sign-In
Future<AuthResult> loginWithGoogle() async {
  final googleUser = await GoogleSignIn(scopes: ['email', 'profile']).signIn();
  if (googleUser == null) throw AuthException('Google sign-in cancelled');

  final googleAuth = await googleUser.authentication;
  final result = await _rest.socialLogin('google', googleAuth.idToken!);
  return AuthResult.fromJson(result);
}

// Kakao Login
Future<AuthResult> loginWithKakao() async {
  final OAuthToken token;
  if (await isKakaoTalkInstalled()) {
    token = await UserApi.instance.loginWithKakaoTalk();
  } else {
    token = await UserApi.instance.loginWithKakaoAccount();
  }
  final result = await _rest.socialLogin('kakao', token.accessToken);
  return AuthResult.fromJson(result);
}

// Apple Sign-In
Future<AuthResult> loginWithApple() async {
  final credential = await SignInWithApple.getAppleIDCredential(
    scopes: [AppleIDAuthorizationScopes.email, AppleIDAuthorizationScopes.fullName],
  );
  final result = await _rest.socialLogin('apple', credential.identityToken!);
  return AuthResult.fromJson(result);
}
```

**pubspec.yaml 추가**: `google_sign_in`, `kakao_flutter_sdk_user`, `sign_in_with_apple`

---

### 4.5 Day 5: PASS 본인인증 + HealthKit

#### Task 13-05: PASS 본인인증 연동

**파일**: `frontend/flutter-app/lib/core/services/identity_verification_service.dart`

```dart
/// SimulatedIdentityService → PassIdentityService
class PassIdentityService implements IdentityVerificationService {
  @override
  Future<VerificationResult> verify() async {
    // PASS 인증 웹뷰 호출
    final result = await Navigator.push(context, MaterialPageRoute(
      builder: (_) => PassVerificationWebView(
        merchantId: const String.fromEnvironment('PASS_MERCHANT_ID'),
        callbackUrl: 'manpasik://identity/callback',
      ),
    ));
    return VerificationResult(
      success: result.success,
      name: result.name,
      phone: result.phone,
      birthDate: result.birthDate,
    );
  }
}
```

#### Task 13-06: HealthKit / Health Connect 연동

**파일**: `frontend/flutter-app/lib/core/services/health_connect_service.dart`

```dart
/// SimulatedHealthConnect → RealHealthConnectService
class RealHealthConnectService implements HealthConnectService {
  final Health _health = Health();

  @override
  Future<bool> requestPermission() async {
    final types = [
      HealthDataType.BLOOD_GLUCOSE,
      HealthDataType.HEART_RATE,
      HealthDataType.WEIGHT,
      HealthDataType.STEPS,
    ];
    return _health.requestAuthorization(types);
  }

  @override
  Future<List<HealthDataPoint>> readData(DateTime start, DateTime end) async {
    return _health.getHealthDataFromTypes(start, end, [
      HealthDataType.BLOOD_GLUCOSE,
      HealthDataType.HEART_RATE,
    ]);
  }
}
```

**pubspec.yaml 추가**: `health: ^10.x`

---

### Sprint 13 검증 기준

```
✅ Toss Payments 웹뷰 로드 확인 (sandbox 모드)
✅ FCM 토큰 발급 확인 (Firebase console)
✅ WebRTC localStream 획득 확인 (카메라 프리뷰)
✅ Google/Kakao 로그인 토큰 획득 확인
✅ HealthKit 권한 다이얼로그 표시 확인
✅ SimulatedXxxService 참조 0건 (프로덕션 빌드에서)
```

---

## 5. Sprint 14: 스토리보드 UI 100% + 폴리싱

> **목표**: 스토리보드 평균 UI 일치율 67% → 95%+
> **작업량**: UI 갭 해소 40건 + Lottie 교체 5건 + 애니메이션 추가 8건

### 5.1 Day 1: Lottie 애니메이션 교체 (5건)

#### Task 14-01: Lottie 에셋 교체

| 위치 | 현재 | 교체 |
|------|------|------|
| `splash_screen.dart` 로고 | 플레이스홀더 텍스트 | `assets/lottie/logo_intro.json` |
| `onboarding_screen.dart` 웰컴 | 텍스트만 | `assets/lottie/confetti.json` |
| `order_complete_screen.dart` | 텍스트 체크마크 | `assets/lottie/check_success.json` |
| `measurement_screen.dart` 대기 | 정적 아이콘 | `assets/lottie/blood_drop.json` |
| `ble_scan_dialog.dart` 검색 | CircularProgress | `assets/lottie/ble_scanning.json` |

**pubspec.yaml 확인**: `lottie: ^3.x` (이미 포함된 경우 에셋만 추가)

---

### 5.2 Day 2: 스토리보드 갭 해소 — 홈/알림/측정

#### Task 14-02: 홈 대시보드 UI 보강

**파일**: `features/home/presentation/home_screen.dart`

```
추가 요소:
- 미니 스파크라인 차트 (최근 7일 바이오마커 추세)
- 무한 스크롤 페이지네이션 (최근 측정 목록)
- 날짜별 그룹핑 헤더
```

#### Task 14-03: 알림 센터 필터 탭

**파일**: `features/notification/presentation/notification_screen.dart`

```
추가 요소:
- 카테고리 필터 탭 (전체/측정/가족/의료/시스템)
- 스와이프 삭제 (Dismissible)
- 딥링크 이동 (알림 탭 시 해당 화면으로)
```

---

### 5.3 Day 3: 스토리보드 갭 해소 — 커뮤니티/가족

#### Task 14-04: 커뮤니티 UI 보강

**파일**: `features/community/presentation/community_screen.dart`

```
추가 요소:
- 댓글 목록 (PostCard 확장)
- 챌린지 리더보드 위젯
- 챌린지 오늘 인증하기 버튼
- 게시글 작성 시 익명/실명 선택 토글
- 게시글 작성 시 측정 데이터 첨부 체크박스
```

#### Task 14-05: 가족 관리 UI 보강

**파일**: `features/family/presentation/family_create_screen.dart`

```
추가 요소:
- QR 코드 초대 (qr_flutter 패키지)
- 링크 공유 (share_plus 패키지)
- 보호자 대시보드 주간 리포트 위젯
```

---

### 5.4 Day 4: 스토리보드 갭 해소 — 데이터허브/구독/의료

#### Task 14-06: 데이터 허브 UI 보강

**파일**: `features/data_hub/presentation/data_hub_screen.dart`

```
추가 요소:
- My Zone 오버레이 (개인 기준선 범위 표시)
- 추세 화살표 (↑↓→) 각 바이오마커 옆
- 핀치 줌 제스처
- 외부 연동 설정 섹션 (HealthKit/Health Connect 토글)
```

#### Task 14-07: 구독 UI 보강

**파일**: `features/market/presentation/plan_comparison_screen.dart`

```
추가 요소:
- 연간/월간 토글 (연간 20% 할인)
- 쿠폰 코드 입력 필드
- 결제 수단 변경 기능
```

#### Task 14-08: 화상진료 UI 보강

**파일**: `features/medical/presentation/video_call_screen.dart`

```
추가 요소:
- 종료 확인 다이얼로그 (2단계)
- 실시간 바이오데이터 패널 (최근 측정 표시)
- 자동 재연결 로직 (ICE restart)
- 영상 품질 자동 조절
```

---

### 5.5 Day 5: 나머지 UI 갭 + 접근성

#### Task 14-09: 설정/고객지원 UI 보강

```
- 약관 변경 이력 화면
- 답변 알림 설정 토글
- FAQ 봇 문의 (AI 챗봇 연결)
```

#### Task 14-10: 접근성 보강

```
- Semantics 라벨 전수 추가 (모든 버튼/아이콘)
- 최소 터치 영역 48x48dp 검증
- 고대비 모드 색상 대비 4.5:1 검증
- 시니어 모드 폰트 1.5배 검증
```

---

### Sprint 14 검증 기준

```
✅ 스토리보드 18개 × 73씬 UI 일치율 재측정: 95%+
✅ Lottie 애니메이션 5개 정상 재생
✅ 접근성 검사 (flutter test --accessibility)
✅ Flutter analyze 에러 0
```

---

## 6. Sprint 15: E2E 테스트 + 빌드 검증 + 출시 준비

> **목표**: E2E 30개 시나리오 PASS, 보안 감사 PASS, 출시 체크리스트 100%

### 6.1 Day 1-2: E2E 통합 테스트 작성 (30개 시나리오)

#### Task 15-01: Flutter Integration Test

**파일**: `frontend/flutter-app/integration_test/` (신규 디렉토리)

```dart
// e2e_auth_test.dart
void main() {
  testWidgets('E2E-01: 회원가입 → 로그인 → 홈', (tester) async { ... });
  testWidgets('E2E-02: 소셜 로그인 (Google)', (tester) async { ... });
  testWidgets('E2E-03: 비밀번호 찾기 흐름', (tester) async { ... });
}

// e2e_measurement_test.dart
void main() {
  testWidgets('E2E-04: 측정 시작 → 결과 → 저장', (tester) async { ... });
  testWidgets('E2E-05: 오프라인 측정 → 동기화', (tester) async { ... });
  testWidgets('E2E-06: 핑거프린트 분석 표시', (tester) async { ... });
}

// e2e_market_test.dart
void main() {
  testWidgets('E2E-07: 마켓 → 상품 → 장바구니 → 결제', (tester) async { ... });
  testWidgets('E2E-08: 구독 업그레이드 흐름', (tester) async { ... });
  testWidgets('E2E-09: 구독 해지 흐름', (tester) async { ... });
}

// e2e_medical_test.dart
void main() {
  testWidgets('E2E-10: 병원 검색 → 예약 → 화상진료', (tester) async { ... });
  testWidgets('E2E-11: 처방전 조회 → 약국 전송', (tester) async { ... });
  testWidgets('E2E-12: 예약 취소', (tester) async { ... });
}

// e2e_family_test.dart
void main() {
  testWidgets('E2E-13: 가족 그룹 생성 → 초대', (tester) async { ... });
  testWidgets('E2E-14: 보호자 대시보드 확인', (tester) async { ... });
  testWidgets('E2E-15: 가족 건강 리포트', (tester) async { ... });
}

// e2e_community_test.dart
void main() {
  testWidgets('E2E-16: 게시글 작성 → 댓글 → 좋아요', (tester) async { ... });
  testWidgets('E2E-17: 챌린지 참여 → 진행률', (tester) async { ... });
  testWidgets('E2E-18: QnA 질문 → 답변 확인', (tester) async { ... });
}

// e2e_ai_test.dart
void main() {
  testWidgets('E2E-19: AI 코칭 인사이트 확인', (tester) async { ... });
  testWidgets('E2E-20: 음식 사진 → 칼로리 분석', (tester) async { ... });
  testWidgets('E2E-21: AI 챗봇 대화', (tester) async { ... });
}

// e2e_data_test.dart
void main() {
  testWidgets('E2E-22: 데이터 허브 트렌드 차트', (tester) async { ... });
  testWidgets('E2E-23: FHIR 내보내기', (tester) async { ... });
  testWidgets('E2E-24: 충돌 해결 화면', (tester) async { ... });
}

// e2e_admin_test.dart
void main() {
  testWidgets('E2E-25: 관리자 대시보드 통계', (tester) async { ... });
  testWidgets('E2E-26: 사용자 검색 → 정지', (tester) async { ... });
  testWidgets('E2E-27: 감사 로그 필터링', (tester) async { ... });
}

// e2e_settings_test.dart
void main() {
  testWidgets('E2E-28: 설정 → 테마 변경 → 언어 변경', (tester) async { ... });
  testWidgets('E2E-29: 긴급 연락처 추가/삭제', (tester) async { ... });
  testWidgets('E2E-30: 동의 관리 → 철회', (tester) async { ... });
}
```

#### Task 15-02: Go 백엔드 통합 테스트

**파일**: `backend/services/gateway/internal/handler/integration_test.go` (신규)

```go
// 30개 REST 엔드포인트 통합 테스트
func TestE2E_AuthFlow(t *testing.T) { ... }
func TestE2E_MeasurementFlow(t *testing.T) { ... }
func TestE2E_MarketPurchaseFlow(t *testing.T) { ... }
func TestE2E_FamilyManagementFlow(t *testing.T) { ... }
func TestE2E_TelemedicineFlow(t *testing.T) { ... }
```

---

### 6.2 Day 3: 전체 빌드 + 정적 분석 재검증

#### Task 15-03: 빌드 검증

```bash
# Go 전체 빌드
for svc in $(ls backend/services/); do
  wsl -d Ubuntu -- bash -c "export PATH=/usr/local/go/bin:/usr/bin:/bin && \
    cd /home/kangjh3kang/Manpasik/backend/services/$svc && \
    GOWORK=off go build ./..."
done

# Go 전체 테스트
for svc in $(ls backend/services/); do
  wsl -d Ubuntu -- bash -c "export PATH=/usr/local/go/bin:/usr/bin:/bin && \
    cd /home/kangjh3kang/Manpasik/backend/services/$svc && \
    GOWORK=off go test ./..."
done

# Flutter analyze
wsl -d Ubuntu -- bash --norc --noprofile -c \
  'export PATH=/home/kangjh3kang/flutter/bin:$PATH && \
   cd /home/kangjh3kang/Manpasik/frontend/flutter-app && \
   flutter analyze && flutter test'
```

#### Task 15-04: 정량 지표 최종 검증

```bash
# REST 엔드포인트 수 (목표: 159+)
grep -c 'mux.HandleFunc' backend/services/gateway/internal/handler/*.go

# Proto RPC 수 vs REST 수 비교
grep -c 'rpc ' backend/proto/*.proto
# → 두 수가 일치해야 함

# Flutter placeholder 잔여 0건
grep -rn "UnimplementedError\|SimulatedPayment\|SimulatedIdentity\|SimulatedPharmacy\|SimulatedPublicData\|SimulatedHealthConnect" \
  frontend/flutter-app/lib/

# 테스트 함수 총 수 (목표: 400+)
grep -c "func Test\|testWidgets" backend/services/*/internal/service/*_test.go \
  frontend/flutter-app/integration_test/*.dart frontend/flutter-app/test/**/*_test.dart
```

---

### 6.3 Day 4: 보안 감사 + 성능

#### Task 15-05: 보안 체크리스트

| # | 항목 | 검증 방법 | 통과 기준 |
|---|------|-----------|-----------|
| 1 | SQL Injection | Parameterized query 전수 검사 | 문자열 연결 0건 |
| 2 | XSS | HTML sanitize 검사 | 사용자 입력 직접 렌더 0건 |
| 3 | CSRF | Token 검증 | 모든 POST/PUT/DELETE에 인증 |
| 4 | SSL Pinning | 프로덕션 빌드 확인 | kReleaseMode에서 활성 |
| 5 | 민감 데이터 로깅 | 로그 검사 | PII 노출 0건 |
| 6 | API 인증 | Gateway middleware | 인증 없는 POST 0건 (auth 제외) |
| 7 | 의존성 취약점 | `flutter pub audit` / `go vuln check` | Critical 0건 |

#### Task 15-06: 성능 프로파일링

```
- Flutter DevTools 프로파일링: 60fps 유지 확인
- 메모리 누수 검사: 화면 전환 100회 후 메모리 증가 < 10MB
- API 응답 시간: P95 < 500ms
- 오프라인 동기화: 100건 큐 처리 < 10초
```

---

### 6.4 Day 5: 출시 체크리스트

#### Task 15-07: 최종 출시 체크리스트

```
[ ] Go 빌드 22/22 PASS (gateway + 21 서비스)
[ ] Go 테스트 400+ PASS
[ ] Flutter analyze 에러 0
[ ] Flutter E2E 30/30 PASS
[ ] REST 엔드포인트 159/169 노출 (내부전용 10개 제외)
[ ] Repository placeholder 0건
[ ] SimulatedXxxService 0건 (프로덕션)
[ ] Lottie 5개 정상 재생
[ ] 스토리보드 UI 일치율 95%+
[ ] 보안 감사 Critical 0건
[ ] 60fps 프로파일링 PASS
[ ] Android APK 빌드 성공
[ ] iOS IPA 빌드 성공 (Apple Developer 계정 필요)
[ ] Firebase 프로젝트 설정 완료
[ ] Toss 가맹점 심사 완료
[ ] PASS 본인인증 계약 완료
[ ] 개인정보 처리방침 최종 검토
[ ] 의료 면책 고지문 법률 검토
```

---

## 7. 작업 의존성 그래프

```
Sprint 11 (백엔드 REST)
  ├─ Task 11-01~09: REST 엔드포인트 76개
  ├─ Task 11-10: Placeholder 수정 7건
  └─ Task 11-11~13: Flutter 화면 3개 ──────┐
                                            │
Sprint 12 (Flutter 연동)  ◀── depends on ───┘
  ├─ Task 12-01: rest_client 76개 메서드
  ├─ Task 12-02~07: Repository 구현체 20파일
  └─ Task 12-08~09: Provider + 화면 바인딩 ──┐
                                              │
Sprint 13 (SDK 연동)  ◀── depends on ─────────┘
  ├─ Task 13-01: Toss Payments (독립)
  ├─ Task 13-02: FCM (독립)
  ├─ Task 13-03: WebRTC (Video REST 필요 → S11 의존)
  ├─ Task 13-04: OAuth2 (독립)
  └─ Task 13-05~06: PASS/HealthKit (독립)
                                              │
Sprint 14 (UI 폴리싱)  ◀── depends on ────────┘
  ├─ Task 14-01: Lottie (독립)
  ├─ Task 14-02~09: 스토리보드 갭 (S12 연동 필요)
  └─ Task 14-10: 접근성 (독립)
                                              │
Sprint 15 (E2E + 출시)  ◀── depends on ───────┘
  ├─ Task 15-01~02: E2E 테스트 (전체 의존)
  ├─ Task 15-03~04: 빌드 검증
  ├─ Task 15-05~06: 보안/성능
  └─ Task 15-07: 출시 체크리스트
```

---

## 부록: 작업량 요약

| Sprint | 백엔드 (Go) | 프론트엔드 (Flutter) | 신규 파일 | 수정 파일 |
|--------|------------|---------------------|----------|----------|
| 11 | REST 핸들러 76개 + placeholder 수정 | 화면 3개 신규 | ~8 | ~5 |
| 12 | — | rest_client 76메서드 + repo 20파일 | ~0 | ~25 |
| 13 | — | SDK 연동 6건 + 웹뷰 2개 | ~4 | ~8 |
| 14 | — | UI 보강 40건 + Lottie 5개 | ~5(에셋) | ~15 |
| 15 | 통합 테스트 | E2E 30개 + 빌드 검증 | ~12 | ~5 |
| **합계** | **~84 핸들러** | **~180 메서드/위젯** | **~29** | **~58** |

### 예상 완료 후 지표

| 지표 | 현재 | 목표 | 달성 Sprint |
|------|------|------|------------|
| REST 노출률 | 49% | **100%** | Sprint 11 |
| Flutter Placeholder | 20건 | **0건** | Sprint 12 |
| 시뮬레이션 서비스 | 7건 | **0건** | Sprint 13 |
| 스토리보드 UI 일치율 | 67% | **95%+** | Sprint 14 |
| E2E 테스트 | 0개 | **30+** | Sprint 15 |
| Go 빌드 | 11/11 | **22/22** | Sprint 15 |
| Go 테스트 | 319개 | **400+** | Sprint 15 |
| Flutter analyze | 에러 0 | **에러 0** | 매 Sprint |
