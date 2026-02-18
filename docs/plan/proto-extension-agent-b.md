# Proto Extension Proposal — Agent B (처방/약국)

> Sprint 1 Phase 2 | Field Range: **250–299** | Services: prescription-service

## 1. 개요

Agent B는 prescription-service Phase 1에서 구현된 약국 선택, 조제 토큰, 상태 머신, 복약 알림 도메인을 proto에 반영합니다.

## 2. 신규 Enum 정의

```protobuf
// Field range 250–259: Enums
enum FulfillmentType {
  FULFILLMENT_TYPE_UNSPECIFIED = 0;
  FULFILLMENT_TYPE_PICKUP      = 1;   // 약국 직접 수령
  FULFILLMENT_TYPE_COURIER     = 2;   // 퀵/택배 배송
  FULFILLMENT_TYPE_DELIVERY    = 3;   // 일반 배송
}

enum DispensaryStatus {
  DISPENSARY_STATUS_UNSPECIFIED = 0;
  DISPENSARY_STATUS_PENDING     = 1;   // 대기
  DISPENSARY_STATUS_PREPARING   = 2;   // 조제 중
  DISPENSARY_STATUS_READY       = 3;   // 조제 완료
  DISPENSARY_STATUS_DISPENSED   = 4;   // 수령 완료
}

enum InteractionSeverity {
  INTERACTION_SEVERITY_UNSPECIFIED     = 0;
  INTERACTION_SEVERITY_NONE            = 1;
  INTERACTION_SEVERITY_MINOR           = 2;
  INTERACTION_SEVERITY_MODERATE        = 3;
  INTERACTION_SEVERITY_MAJOR           = 4;
  INTERACTION_SEVERITY_CONTRAINDICATED = 5;
}
```

## 3. 신규 메시지 정의

### 3.1 Medication (약물) 확장

```protobuf
message Medication {
  string id                = 1;
  string drug_name         = 2;
  string drug_code         = 3;
  string dosage            = 4;
  string frequency         = 5;
  int32  duration_days     = 250;
  string route             = 251;   // "oral", "injection", "topical"
  string instructions      = 252;
  int32  quantity          = 253;
  int32  refills_remaining = 254;
  bool   is_generic_allowed = 255;
}
```

### 3.2 DrugInteraction (약물 상호작용)

```protobuf
message DrugInteraction {
  string drug_a          = 250;
  string drug_b          = 251;
  InteractionSeverity severity = 252;
  string description     = 253;
  string recommendation  = 254;
}
```

### 3.3 MedicationReminder (복약 알림)

```protobuf
message MedicationReminder {
  string prescription_id = 250;
  string medication_id   = 251;
  string drug_name       = 252;
  string dosage          = 253;
  string time_of_day     = 254;  // "08:00", "13:00", "20:00"
  string instructions    = 255;
  bool   is_taken        = 256;
}
```

### 3.4 FulfillmentToken (조제 토큰)

```protobuf
message FulfillmentToken {
  string token           = 250;
  string prescription_id = 251;
  string pharmacy_id     = 252;
  google.protobuf.Timestamp created_at = 253;
  google.protobuf.Timestamp expires_at = 254;
  bool   is_used         = 255;
  google.protobuf.Timestamp used_at    = 256;
}
```

### 3.5 Prescription 확장 필드

기존 `Prescription` 메시지에 추가:
```protobuf
message Prescription {
  // 기존 필드 유지...
  string pharmacy_id        = 250;
  string pharmacy_name      = 251;
  FulfillmentType fulfillment_type = 252;
  string shipping_address   = 253;
  string fulfillment_token  = 254;
  DispensaryStatus dispensary_status = 255;
  google.protobuf.Timestamp sent_to_pharmacy_at = 256;
  google.protobuf.Timestamp dispensed_at        = 257;
}
```

## 4. 신규 RPC 정의

### PrescriptionService 확장

```protobuf
service PrescriptionService {
  // 기존 RPC 유지...

  // --- Phase 1 추가 ---
  rpc SelectPharmacyAndFulfillment(SelectPharmacyRequest) returns (SelectPharmacyResponse);
  rpc SendToPharmacy(SendToPharmacyRequest) returns (SendToPharmacyResponse);
  rpc GetByToken(GetByTokenRequest) returns (GetByTokenResponse);
  rpc UpdateDispensaryStatus(UpdateDispensaryStatusRequest) returns (UpdateDispensaryStatusResponse);
  rpc CheckDrugInteraction(CheckDrugInteractionRequest) returns (CheckDrugInteractionResponse);
  rpc GetMedicationReminders(GetMedicationRemindersRequest) returns (GetMedicationRemindersResponse);
}

// --- Request/Response ---
message SelectPharmacyRequest {
  string prescription_id  = 1;
  string pharmacy_id      = 2;
  string pharmacy_name    = 3;
  FulfillmentType fulfillment_type = 4;
  string shipping_address = 5;
}
message SelectPharmacyResponse {
  bool success = 1;
}

message SendToPharmacyRequest {
  string prescription_id = 1;
}
message SendToPharmacyResponse {
  FulfillmentToken token = 1;
}

message GetByTokenRequest {
  string token = 1;
}
message GetByTokenResponse {
  Prescription prescription = 1;
}

message UpdateDispensaryStatusRequest {
  string prescription_id       = 1;
  DispensaryStatus new_status  = 2;
}
message UpdateDispensaryStatusResponse {
  bool success = 1;
}

message CheckDrugInteractionRequest {
  repeated string drug_codes = 1;
}
message CheckDrugInteractionResponse {
  repeated DrugInteraction interactions = 1;
}

message GetMedicationRemindersRequest {
  string patient_user_id = 1;
}
message GetMedicationRemindersResponse {
  repeated MedicationReminder reminders = 1;
}
```

## 5. 상태 전이 규칙 (State Machine)

```
DispensaryStatus 전이:
  pending → preparing → ready → dispensed

PrescriptionStatus 전이:
  draft → active → dispensed → completed
                 → cancelled
                 → expired (시간 초과)
```

유효하지 않은 전이 시 `INVALID_ARGUMENT` gRPC 에러를 반환합니다.

## 6. Kafka 이벤트

| 이벤트 | Topic | Payload |
|--------|-------|---------|
| `prescription.created` | `prescription-events` | prescription_id, user_id, doctor_id, diagnosis |
| `prescription.sent_to_pharmacy` | `prescription-events` | prescription_id, pharmacy_id, pharmacy_name, token |
| `prescription.dispensed` | `prescription-events` | prescription_id, user_id, pharmacy_id |

## 7. 영향 분석

| 항목 | 내용 |
|------|------|
| 영향받는 서비스 | prescription-service |
| 연동 서비스 | notification-service (조제 완료 알림), reservation-service (진료→처방 연결) |
| 하위 호환성 | 신규 필드/RPC 추가만 — 기존 필드 변경 없음 |
| 필드 번호 충돌 | 없음 (250–299 범위 전용) |
| DB 마이그레이션 | prescriptions 테이블에 pharmacy/fulfillment 컬럼 추가, fulfillment_tokens 테이블 신규 생성 |

## 8. Phase 3 핸들러 구현 계획

Proto 병합 후 `handler/grpc.go`에서:
1. `SelectPharmacyAndFulfillment` — 약국 선택 + 수령 방식 설정
2. `SendToPharmacy` — 조제 토큰 발급 + 약국 전송
3. `GetByToken` — 토큰으로 처방전 조회
4. `UpdateDispensaryStatus` — 조제 상태 전이 (state machine)
5. `CheckDrugInteraction` — 약물 상호작용 검사
6. `GetMedicationReminders` — 복약 알림 스케줄 조회
