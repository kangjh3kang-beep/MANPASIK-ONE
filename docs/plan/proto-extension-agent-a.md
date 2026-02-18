# Proto Extension Proposal — Agent A (의료/예약)

> Sprint 1 Phase 2 | Field Range: **200–249** | Services: reservation-service, telemedicine-service

## 1. 개요

Agent A는 reservation-service와 telemedicine-service의 Phase 1 서비스 로직에서 도입된 신규 도메인 객체를 proto에 반영합니다.
현재 `manpasik.proto`의 `ReservationService`와 `TelemedicineService`에 없는 메시지/RPC를 추가합니다.

## 2. 신규 Enum 정의

```protobuf
// Field range 200–209: Enums
enum FacilityType {
  FACILITY_TYPE_UNSPECIFIED = 0;
  FACILITY_TYPE_HOSPITAL    = 1;   // 병원
  FACILITY_TYPE_CLINIC      = 2;   // 의원
  FACILITY_TYPE_PHARMACY    = 3;   // 약국
  FACILITY_TYPE_DENTAL      = 4;   // 치과
  FACILITY_TYPE_ORIENTAL    = 5;   // 한의원
}

enum Specialty {
  SPECIALTY_UNSPECIFIED   = 0;
  SPECIALTY_GENERAL       = 1;   // 일반의
  SPECIALTY_INTERNAL      = 2;   // 내과
  SPECIALTY_CARDIOLOGY    = 3;   // 심장내과
  SPECIALTY_ENDOCRINOLOGY = 4;   // 내분비내과
  SPECIALTY_DERMATOLOGY   = 5;   // 피부과
  SPECIALTY_PEDIATRICS    = 6;   // 소아과
  SPECIALTY_PSYCHIATRY    = 7;   // 정신과
  SPECIALTY_ORTHOPEDICS   = 8;   // 정형외과
  SPECIALTY_OPHTHALMOLOGY = 9;   // 안과
  SPECIALTY_ENT           = 10;  // 이비인후과
}

enum ReservationStatus {
  RESERVATION_STATUS_UNSPECIFIED = 0;
  RESERVATION_STATUS_PENDING     = 1;
  RESERVATION_STATUS_CONFIRMED   = 2;
  RESERVATION_STATUS_COMPLETED   = 3;
  RESERVATION_STATUS_CANCELLED   = 4;
  RESERVATION_STATUS_NO_SHOW     = 5;
}
```

## 3. 신규 메시지 정의

### 3.1 Region (지역 계층)

```protobuf
message Region {
  string id            = 200;
  string country_code  = 201;  // "KR", "US", "JP"
  string region_code   = 202;  // "seoul", "tokyo"
  string district_code = 203;  // "gangnam", "shibuya"
  string name          = 204;
  string name_local    = 205;  // 현지어 이름
  string parent_id     = 206;
  string timezone      = 207;  // "Asia/Seoul"
}
```

### 3.2 Facility (의료 시설) 확장

```protobuf
message Facility {
  string id                  = 1;
  string name                = 2;
  FacilityType type          = 200;
  string address             = 3;
  string phone               = 4;
  double latitude            = 5;
  double longitude           = 6;
  double distance_km         = 201;   // 계산된 Haversine 거리
  double rating              = 202;
  int32  review_count        = 203;
  repeated Specialty specialties = 204;
  string operating_hours     = 205;
  bool   is_open_now         = 206;
  bool   accepts_reservation = 207;
  string image_url           = 208;
  string country_code        = 209;
  string region_code         = 210;
  string district_code       = 211;
  string timezone            = 212;
  bool   has_telemedicine    = 213;
}
```

### 3.3 Doctor (의사)

```protobuf
message Doctor {
  string id                     = 1;
  string facility_id            = 2;
  string user_id                = 200;
  string name                   = 3;
  string specialty              = 4;
  string license_number         = 201;
  repeated string languages     = 202;
  bool   is_available           = 203;
  double rating                 = 204;
  int32  total_consultations    = 205;
  google.protobuf.Timestamp next_available_at = 206;
  int32  consultation_fee       = 207;
  bool   accepts_telemedicine   = 208;
  repeated string available_region_codes = 209;
}
```

### 3.4 TimeSlot (예약 가능 시간대)

```protobuf
message TimeSlot {
  string id          = 200;
  google.protobuf.Timestamp start_time = 201;
  google.protobuf.Timestamp end_time   = 202;
  bool   is_available = 203;
  string doctor_id    = 204;
  string doctor_name  = 205;
}
```

## 4. 신규 RPC 정의

### ReservationService 확장

```protobuf
service ReservationService {
  // 기존 RPC 유지...

  // --- Phase 1 추가 ---
  rpc ListDoctorsByFacility(ListDoctorsRequest) returns (ListDoctorsResponse);
  rpc GetDoctorAvailability(GetDoctorAvailabilityRequest) returns (GetDoctorAvailabilityResponse);
  rpc SelectDoctor(SelectDoctorRequest) returns (SelectDoctorResponse);
  rpc ListRegions(ListRegionsRequest) returns (ListRegionsResponse);
}

// --- Request/Response ---
message ListDoctorsRequest {
  string facility_id = 1;
  string specialty   = 2;
}
message ListDoctorsResponse {
  repeated Doctor doctors = 1;
}

message GetDoctorAvailabilityRequest {
  string doctor_id = 1;
  google.protobuf.Timestamp date = 2;
}
message GetDoctorAvailabilityResponse {
  repeated TimeSlot slots = 1;
}

message SelectDoctorRequest {
  string facility_id = 1;
  string doctor_id   = 2;
  string user_id     = 3;
}
message SelectDoctorResponse {
  Doctor doctor = 1;
}

message ListRegionsRequest {
  string country_code = 1;
  string region_code  = 2;  // optional: 하위 district만 조회
}
message ListRegionsResponse {
  repeated Region regions = 1;
}
```

### SearchFacilities 확장 필드

기존 `SearchFacilitiesRequest`에 추가:
```protobuf
message SearchFacilitiesRequest {
  // 기존 필드 유지...
  string country_code  = 200;
  string region_code   = 201;
  string district_code = 202;
  double user_lat      = 203;
  double user_lon      = 204;
}
```

## 5. 영향 분석

| 항목 | 내용 |
|------|------|
| 영향받는 서비스 | reservation-service, telemedicine-service |
| 하위 호환성 | 신규 필드/RPC 추가만 — 기존 필드 변경 없음 |
| 필드 번호 충돌 | 없음 (200–249 범위 전용) |
| Kafka 이벤트 | `reservation.created`, `reservation.cancelled` 이벤트 스키마 변경 없음 |
| 마이그레이션 | DB 스키마 변경 필요 (facilities 테이블에 region 컬럼 추가) |

## 6. Phase 3 핸들러 구현 계획

Proto 확장 병합 후 `handler/grpc.go`에서:
1. `SearchFacilities` — 신규 region 필터 + Haversine 거리 파라미터 매핑
2. `ListDoctorsByFacility` — 신규 RPC 핸들러
3. `GetDoctorAvailability` — 신규 RPC 핸들러
4. `SelectDoctor` — 신규 RPC 핸들러
5. `ListRegions` — 신규 RPC 핸들러
