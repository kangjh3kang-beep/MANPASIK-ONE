# Proto Extension Proposal — Agent C (데이터 공유/FHIR)

> Sprint 1 Phase 2 | Field Range: **300–349** | Services: health-record-service

## 1. 개요

Agent C는 health-record-service Phase 1에서 구현된 데이터 공유 동의(Consent), FHIR R4 Export/Import, 데이터 접근 로그 도메인을 proto에 반영합니다.

## 2. 신규 Enum 정의

```protobuf
// Field range 300–309: Enums
enum ConsentType {
  CONSENT_TYPE_UNSPECIFIED      = 0;
  CONSENT_TYPE_MEASUREMENT_SHARE = 1;
  CONSENT_TYPE_RECORD_SHARE      = 2;
  CONSENT_TYPE_FULL_ACCESS       = 3;
}

enum ConsentStatus {
  CONSENT_STATUS_UNSPECIFIED = 0;
  CONSENT_STATUS_ACTIVE      = 1;
  CONSENT_STATUS_REVOKED     = 2;
  CONSENT_STATUS_EXPIRED     = 3;
}

enum FHIRResourceType {
  FHIR_RESOURCE_TYPE_UNSPECIFIED          = 0;
  FHIR_RESOURCE_TYPE_OBSERVATION          = 1;
  FHIR_RESOURCE_TYPE_CONDITION            = 2;
  FHIR_RESOURCE_TYPE_MEDICATION_STATEMENT = 3;
  FHIR_RESOURCE_TYPE_ALLERGY_INTOLERANCE  = 4;
  FHIR_RESOURCE_TYPE_IMMUNIZATION         = 5;
  FHIR_RESOURCE_TYPE_PROCEDURE            = 6;
  FHIR_RESOURCE_TYPE_DIAGNOSTIC_REPORT    = 7;
  FHIR_RESOURCE_TYPE_PATIENT              = 8;
}

enum HealthRecordType {
  HEALTH_RECORD_TYPE_UNSPECIFIED  = 0;
  HEALTH_RECORD_TYPE_LAB_RESULT   = 1;
  HEALTH_RECORD_TYPE_IMAGING      = 2;
  HEALTH_RECORD_TYPE_VITAL_SIGN   = 3;
  HEALTH_RECORD_TYPE_ALLERGY      = 4;
  HEALTH_RECORD_TYPE_CONDITION    = 5;
  HEALTH_RECORD_TYPE_IMMUNIZATION = 6;
  HEALTH_RECORD_TYPE_PROCEDURE    = 7;
}
```

## 3. 신규 메시지 정의

### 3.1 DataSharingConsent (데이터 공유 동의)

```protobuf
message DataSharingConsent {
  string id            = 300;
  string user_id       = 301;
  string provider_id   = 302;   // facility_id
  string provider_name = 303;
  ConsentType consent_type = 304;
  repeated string scope    = 305;   // ["blood_glucose", "blood_pressure", ...]
  string purpose       = 306;   // "treatment", "research", "emergency"
  ConsentStatus status = 307;
  google.protobuf.Timestamp granted_at = 308;
  google.protobuf.Timestamp expires_at = 309;
  google.protobuf.Timestamp revoked_at = 310;
  string revoke_reason = 311;
}
```

### 3.2 DataAccessLog (데이터 접근 로그)

```protobuf
message DataAccessLog {
  string id            = 300;
  string consent_id    = 301;
  string user_id       = 302;
  string provider_id   = 303;
  string action        = 304;   // "view", "export", "share"
  string resource_type = 305;   // "measurement", "health_record"
  repeated string resource_ids = 306;
  google.protobuf.Timestamp accessed_at = 307;
  string ip_address    = 308;
}
```

### 3.3 SharedDataBundle (공유 데이터 번들)

```protobuf
message SharedDataBundle {
  string consent_id       = 300;
  string provider_id      = 301;
  string fhir_bundle_json = 302;
  int32  resource_count   = 303;
  google.protobuf.Timestamp shared_at = 304;
}
```

### 3.4 HealthRecord 확장 필드

기존 `HealthRecord` 메시지에 추가:
```protobuf
message HealthRecord {
  // 기존 필드 유지...
  HealthRecordType record_type  = 300;
  string source                 = 301;   // "manpasik", "manual", "fhir_import"
  string fhir_resource_id       = 302;
  FHIRResourceType fhir_type    = 303;
}
```

## 4. 신규 RPC 정의

### HealthRecordService 확장

```protobuf
service HealthRecordService {
  // 기존 RPC 유지...

  // --- Phase 1 추가: FHIR ---
  rpc ExportToFHIR(ExportToFHIRRequest) returns (ExportToFHIRResponse);
  rpc ImportFromFHIR(ImportFromFHIRRequest) returns (ImportFromFHIRResponse);
  rpc GetHealthSummary(GetHealthSummaryRequest) returns (GetHealthSummaryResponse);

  // --- Phase 1 추가: 데이터 공유 동의 ---
  rpc CreateDataSharingConsent(CreateConsentRequest) returns (CreateConsentResponse);
  rpc RevokeDataSharingConsent(RevokeConsentRequest) returns (RevokeConsentResponse);
  rpc ListDataSharingConsents(ListConsentsRequest) returns (ListConsentsResponse);
  rpc ShareWithProvider(ShareWithProviderRequest) returns (ShareWithProviderResponse);
  rpc GetDataAccessLog(GetAccessLogRequest) returns (GetAccessLogResponse);
}

// --- FHIR Request/Response ---
message ExportToFHIRRequest {
  string user_id = 1;
  repeated HealthRecordType record_types = 2;
  google.protobuf.Timestamp start_date   = 3;
  google.protobuf.Timestamp end_date     = 4;
}
message ExportToFHIRResponse {
  string fhir_bundle_json              = 1;
  int32  resource_count                = 2;
  repeated FHIRResourceType resource_types = 3;
}

message ImportFromFHIRRequest {
  string user_id        = 1;
  string bundle_json    = 2;
}
message ImportFromFHIRResponse {
  repeated HealthRecord imported_records = 1;
  int32  imported_count = 2;
  int32  skipped_count  = 3;
  repeated string errors = 4;
}

message GetHealthSummaryRequest {
  string user_id = 1;
  int32  days    = 2;   // 최근 N일 (기본 30)
}
message GetHealthSummaryResponse {
  int32  total_records              = 1;
  map<string, int32> records_by_type = 2;
  repeated HealthRecord recent_records = 3;
  string summary_text               = 4;
}

// --- Consent Request/Response ---
message CreateConsentRequest {
  string user_id       = 1;
  string provider_id   = 2;
  string provider_name = 3;
  ConsentType consent_type = 4;
  repeated string scope    = 5;
  string purpose       = 6;
  google.protobuf.Timestamp expires_at = 7;
}
message CreateConsentResponse {
  DataSharingConsent consent = 1;
}

message RevokeConsentRequest {
  string consent_id = 1;
  string reason     = 2;
}
message RevokeConsentResponse {
  bool success = 1;
}

message ListConsentsRequest {
  string user_id = 1;
}
message ListConsentsResponse {
  repeated DataSharingConsent consents = 1;
}

message ShareWithProviderRequest {
  string consent_id = 1;
}
message ShareWithProviderResponse {
  SharedDataBundle bundle = 1;
}

message GetAccessLogRequest {
  string user_id = 1;
  int32  limit   = 2;
  int32  offset  = 3;
}
message GetAccessLogResponse {
  repeated DataAccessLog logs = 1;
  int32 total = 2;
}
```

## 5. FHIR R4 매핑 테이블

| HealthRecordType | FHIR ResourceType | LOINC Code (예시) |
|-----------------|-------------------|-------------------|
| VITAL_SIGN | Observation | 85354-9 (Blood Pressure) |
| LAB_RESULT | DiagnosticReport | 2345-7 (Glucose) |
| IMAGING | DiagnosticReport | — |
| ALLERGY | AllergyIntolerance | — |
| CONDITION | Condition | — |
| IMMUNIZATION | Immunization | — |
| PROCEDURE | Procedure | — |

Manpasik 바이오마커 → LOINC 매핑 (15+ biomarkers):
| 바이오마커 | LOINC Code | 단위 |
|-----------|------------|------|
| blood_glucose | 2345-7 | mg/dL |
| blood_pressure_systolic | 8480-6 | mmHg |
| blood_pressure_diastolic | 8462-4 | mmHg |
| heart_rate | 8867-4 | bpm |
| body_temperature | 8310-5 | °C |
| spo2 | 59408-5 | % |
| bmi | 39156-5 | kg/m² |
| total_cholesterol | 2093-3 | mg/dL |
| hdl_cholesterol | 2085-9 | mg/dL |
| ldl_cholesterol | 2089-1 | mg/dL |
| triglycerides | 2571-8 | mg/dL |
| hba1c | 4548-4 | % |
| creatinine | 2160-0 | mg/dL |
| uric_acid | 3084-1 | mg/dL |
| hemoglobin | 718-7 | g/dL |

## 6. 규정 준수

| 규정 | 요구사항 | 구현 |
|------|----------|------|
| PIPA (개인정보보호법) | 동의 기반 데이터 공유 | ConsentType + Scope |
| HIPAA | 최소 필요 원칙 | Scope 필드로 세분화된 접근 |
| GDPR | 동의 철회권 | RevokeDataSharingConsent RPC |
| 의료법 | 접근 로그 보관 | DataAccessLog 5년 보관 |

## 7. 영향 분석

| 항목 | 내용 |
|------|------|
| 영향받는 서비스 | health-record-service |
| 연동 서비스 | family-service (가족 데이터 공유), admin-service (감사 로그) |
| 하위 호환성 | 신규 필드/RPC 추가만 |
| 필드 번호 충돌 | 없음 (300–349 범위 전용) |
| DB 마이그레이션 | data_sharing_consents, data_access_logs 테이블 신규 생성 |

## 8. Phase 3 핸들러 구현 계획

Proto 병합 후 `handler/grpc.go`에서:
1. `ExportToFHIR` — FHIR R4 Bundle JSON 생성
2. `ImportFromFHIR` — FHIR R4 Bundle 파싱 + 레코드 생성
3. `GetHealthSummary` — 기간별 건강 요약
4. `CreateDataSharingConsent` — 동의 생성 (scope 검증)
5. `RevokeDataSharingConsent` — 동의 철회
6. `ShareWithProvider` — 동의 범위 내 FHIR 데이터 공유
7. `GetDataAccessLog` — 접근 이력 조회
