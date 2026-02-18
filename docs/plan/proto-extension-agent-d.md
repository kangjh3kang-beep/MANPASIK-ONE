# Proto Extension Proposal — Agent D (기반 서비스 강화)

> Sprint 1 Phase 2 | Field Range: **350–399** | Services: admin-service, notification-service, family-service

## 1. 개요

Agent D는 admin-service, notification-service, family-service의 Phase 1 서비스 로직에서 도입된 도메인 확장을 proto에 반영합니다.

## 2. 신규 Enum 정의

```protobuf
// Field range 350–359: Enums

// Admin
enum AdminRole {
  ADMIN_ROLE_UNSPECIFIED = 0;
  ADMIN_ROLE_SUPER_ADMIN = 1;
  ADMIN_ROLE_ADMIN       = 2;
  ADMIN_ROLE_MODERATOR   = 3;
  ADMIN_ROLE_SUPPORT     = 4;
  ADMIN_ROLE_ANALYST     = 5;
}

// Family
enum FamilyRole {
  FAMILY_ROLE_UNSPECIFIED = 0;
  FAMILY_ROLE_OWNER       = 1;
  FAMILY_ROLE_GUARDIAN     = 2;
  FAMILY_ROLE_MEMBER       = 3;
  FAMILY_ROLE_CHILD        = 4;
  FAMILY_ROLE_ELDERLY      = 5;
}

enum InvitationStatus {
  INVITATION_STATUS_UNSPECIFIED = 0;
  INVITATION_STATUS_PENDING     = 1;
  INVITATION_STATUS_ACCEPTED    = 2;
  INVITATION_STATUS_DECLINED    = 3;
  INVITATION_STATUS_EXPIRED     = 4;
}

// Notification
enum NotificationChannel {
  NOTIFICATION_CHANNEL_UNSPECIFIED = 0;
  NOTIFICATION_CHANNEL_PUSH        = 1;
  NOTIFICATION_CHANNEL_EMAIL       = 2;
  NOTIFICATION_CHANNEL_SMS         = 3;
  NOTIFICATION_CHANNEL_IN_APP      = 4;
}

enum NotificationPriority {
  NOTIFICATION_PRIORITY_UNSPECIFIED = 0;
  NOTIFICATION_PRIORITY_LOW         = 1;
  NOTIFICATION_PRIORITY_NORMAL      = 2;
  NOTIFICATION_PRIORITY_HIGH        = 3;
  NOTIFICATION_PRIORITY_URGENT      = 4;
}
```

## 3. 신규 메시지 정의

### 3.1 Admin Service

#### AuditLogDetail (확장 감사 로그)

```protobuf
message AuditLogDetail {
  string id         = 350;
  string admin_id   = 351;
  string action     = 352;   // "config_update", "role_change", "admin_deactivate"
  string resource   = 353;   // "config:security.jwt_ttl_hours", "admin:{id}"
  string old_value  = 354;
  string new_value  = 355;
  string ip_address = 356;
  string user_agent = 357;
  google.protobuf.Timestamp created_at = 358;
}
```

#### AdminUser 확장 필드

```protobuf
message AdminUser {
  // 기존 필드 유지...
  string country_code  = 350;
  string region_code   = 351;
  string district_code = 352;
}
```

### 3.2 Notification Service

#### NotificationTemplate (알림 템플릿)

```protobuf
message NotificationTemplate {
  string key       = 350;
  string title     = 351;
  string body_fmt  = 352;   // fmt.Sprintf 형식
  string type      = 353;   // "prescription", "appointment", "health_alert"
  string priority  = 354;   // "low", "normal", "high", "urgent"
  string channel   = 355;   // "push", "email", "sms", "in_app"
}
```

#### NotificationPreferences 확장 필드

```protobuf
message NotificationPreferences {
  // 기존 필드 유지...
  bool   health_alert_enabled = 350;
  bool   coaching_enabled     = 351;
  bool   promotion_enabled    = 352;
  string quiet_hours_start    = 353;   // "HH:MM"
  string quiet_hours_end      = 354;
  string language             = 355;
}
```

### 3.3 Family Service

#### SharingPreferences (공유 설정)

```protobuf
message SharingPreferences {
  string user_id                = 350;
  string group_id               = 351;
  bool   share_measurements     = 352;
  bool   share_health_score     = 353;
  bool   share_goals            = 354;
  bool   share_coaching         = 355;
  bool   share_alerts           = 356;
  repeated string allowed_viewer_ids = 357;
  int32  measurement_days_limit = 358;   // 0 = 무제한
  repeated string allowed_biomarkers = 359;
  bool   require_approval       = 360;
}
```

#### SharedHealthSummary (건강 요약)

```protobuf
message SharedHealthSummary {
  string user_id              = 350;
  string display_name         = 351;
  double health_score         = 352;
  int32  measurements_count   = 353;
  string score_trend          = 354;   // "improving", "stable", "declining"
  string latest_alert         = 355;
  google.protobuf.Timestamp last_measurement_at = 356;
}
```

## 4. 신규 RPC 정의

### 4.1 AdminService 확장

```protobuf
service AdminService {
  // 기존 RPC 유지...

  // --- Phase 1 추가 ---
  rpc ListAdminsByRegion(ListAdminsByRegionRequest) returns (ListAdminsByRegionResponse);
  rpc GetAuditLogDetails(GetAuditLogDetailsRequest) returns (GetAuditLogDetailsResponse);
  rpc SendFromTemplate(SendFromTemplateRequest) returns (SendFromTemplateResponse);
}

message ListAdminsByRegionRequest {
  string country_code = 1;
  string region_code  = 2;
}
message ListAdminsByRegionResponse {
  repeated AdminUser admins = 1;
}

message GetAuditLogDetailsRequest {
  string admin_id = 1;   // optional: 특정 관리자 필터
  string action   = 2;   // optional: 액션 유형 필터
  int32  limit    = 3;
  int32  offset   = 4;
}
message GetAuditLogDetailsResponse {
  repeated AuditLogDetail logs = 1;
  int32 total = 2;
}
```

### 4.2 NotificationService 확장

```protobuf
service NotificationService {
  // 기존 RPC 유지...

  // --- Phase 1 추가 ---
  rpc SendFromTemplate(SendFromTemplateRequest) returns (SendFromTemplateResponse);
}

message SendFromTemplateRequest {
  string user_id      = 1;
  string template_key = 2;
  repeated string args = 3;
}
message SendFromTemplateResponse {
  bool success = 1;
}
```

### 4.3 FamilyService 확장

```protobuf
service FamilyService {
  // 기존 RPC 유지...

  // --- Phase 1 추가 ---
  rpc SetSharingPreferences(SetSharingPreferencesRequest) returns (SetSharingPreferencesResponse);
  rpc ValidateSharingAccess(ValidateSharingAccessRequest) returns (ValidateSharingAccessResponse);
  rpc GetSharedHealthData(GetSharedHealthDataRequest) returns (GetSharedHealthDataResponse);
}

message SetSharingPreferencesRequest {
  SharingPreferences preferences = 1;
}
message SetSharingPreferencesResponse {
  SharingPreferences preferences = 1;
}

message ValidateSharingAccessRequest {
  string group_id        = 1;
  string requester_id    = 2;
  string target_user_id  = 3;
  string biomarker       = 4;
}
message ValidateSharingAccessResponse {
  bool   allowed = 1;
  string reason  = 2;
}

message GetSharedHealthDataRequest {
  string requester_user_id = 1;
  string target_user_id    = 2;   // optional: 비어있으면 전체 멤버
  string group_id          = 3;
  int32  days              = 4;
}
message GetSharedHealthDataResponse {
  repeated SharedHealthSummary summaries = 1;
}
```

## 5. 사전 정의 템플릿 목록 (12개)

| Key | Title | Type | Priority | Channel |
|-----|-------|------|----------|---------|
| `prescription_created` | 새 처방전 발행 | prescription | high | push |
| `prescription_sent` | 처방전 전송 완료 | prescription | normal | push |
| `prescription_ready` | 약 조제 완료 | prescription | high | push |
| `prescription_dispensed` | 약 수령 완료 | prescription | normal | in_app |
| `delivery_started` | 배송 출발 | prescription | normal | push |
| `delivery_arrived` | 배송 완료 | prescription | high | push |
| `appointment_reminder` | 진료 예약 알림 | appointment | high | push |
| `appointment_cancelled` | 예약 취소 | appointment | normal | push |
| `health_alert_critical` | 건강 이상 감지 | health_alert | urgent | push |
| `health_alert_warning` | 건강 주의 | health_alert | high | push |
| `measurement_complete` | 측정 완료 | measurement | normal | in_app |
| `family_data_shared` | 가족 데이터 공유 | system | normal | in_app |

## 6. 영향 분석

| 항목 | 내용 |
|------|------|
| 영향받는 서비스 | admin-service, notification-service, family-service |
| 연동 서비스 | health-record-service (데이터 공유), prescription-service (알림 트리거) |
| 하위 호환성 | 신규 필드/RPC 추가만 |
| 필드 번호 충돌 | 없음 (350–399 범위 전용) |
| DB 마이그레이션 | sharing_preferences 테이블 신규, audit_log_details 테이블 신규, notification_templates 시드 데이터 |

## 7. Phase 3 핸들러 구현 계획

Admin:
1. `ListAdminsByRegion` — 지역별 관리자 조회
2. `GetAuditLogDetails` — OldValue/NewValue 포함 감사 로그

Notification:
3. `SendFromTemplate` — 템플릿 기반 알림 발송

Family:
4. `SetSharingPreferences` — 바이오마커별 세분화 공유 설정
5. `ValidateSharingAccess` — 공유 접근 권한 검증 (RequireApproval, AllowedBiomarkers)
6. `GetSharedHealthData` — Guardian/Owner 우선 접근 + 일반 멤버 설정 기반 접근
