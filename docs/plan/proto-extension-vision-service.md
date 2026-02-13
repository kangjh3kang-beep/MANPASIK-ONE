# Vision Service Proto 확장 제안

> **목적**: vision-service의 gRPC API를 Proto로 정의하여 Gateway·Flutter 클라이언트와 계약을 맞춘다.  
> **상태**: 제안(미반영). `manpasik.proto` 또는 별도 `vision/v1/vision.proto`에 반영 시 이 문서를 참고한다.

## 1. 배경

- vision-service는 이미 **서비스·리포지토리·시뮬레이션 분석**이 구현되어 있으나, **Proto에 VisionService가 정의되어 있지 않다**.
- 핸들러(`internal/handler/grpc.go`)와 `cmd/main.go`에서 `RegisterVisionServiceServer`가 주석 처리된 상태이다.
- 음식 이미지 분석(칼로리·영양소)·이력 조회·일일 요약 API를 Proto로 정의하면 Gateway 라우팅 및 클라이언트 코드 생성이 가능하다.

## 2. 제안 RPC

| RPC | 요청 | 응답 | 비고 |
|-----|------|------|------|
| **AnalyzeFood** | user_id, image_url, meal_type | FoodAnalysis | 이미지 URL 기반 분석(업로드는 별도 스토리지 연동) |
| **GetAnalysis** | analysis_id | FoodAnalysis | 단건 조회 |
| **ListAnalyses** | user_id, limit, offset | ListAnalysesResponse (items, total) | 사용자별 이력 |
| **GetDailySummary** | user_id | GetDailySummaryResponse (total_kcal, meal_breakdown) | 당일 칼로리·끼니별 요약 |

- **인증**: 모든 RPC는 `user_id` 또는 JWT에서 추출한 사용자로 검증 권장.
- **이미지 업로드**: 현재 서비스는 `image_url`(S3/MinIO 등 저장 후 URL)을 받는다. 업로드 스트림 RPC는 필요 시 별도 `UploadFoodImage` 등으로 확장.

## 3. 제안 메시지 정의 (초안)

```protobuf
// ============================================================================
// Vision Service (음식 이미지 분석)
// ============================================================================

service VisionService {
  rpc AnalyzeFood(AnalyzeFoodRequest) returns (FoodAnalysis);
  rpc GetAnalysis(GetAnalysisRequest) returns (FoodAnalysis);
  rpc ListAnalyses(ListAnalysesRequest) returns (ListAnalysesResponse);
  rpc GetDailySummary(GetDailySummaryRequest) returns (GetDailySummaryResponse);
}

message AnalyzeFoodRequest {
  string user_id = 1;
  string image_url = 2;
  string meal_type = 3;  // "breakfast" | "lunch" | "dinner" | "snack"
}

message GetAnalysisRequest {
  string analysis_id = 1;
}

message ListAnalysesRequest {
  string user_id = 1;
  int32 limit = 2;
  int32 offset = 3;
}

message ListAnalysesResponse {
  repeated FoodAnalysis items = 1;
  int32 total = 2;
}

message GetDailySummaryRequest {
  string user_id = 1;
}

message GetDailySummaryResponse {
  double total_kcal = 1;
  map<string, double> meal_breakdown = 2;  // meal_type -> kcal
}

enum FoodAnalysisStatus {
  FOOD_ANALYSIS_STATUS_UNKNOWN = 0;
  FOOD_ANALYSIS_STATUS_PENDING = 1;
  FOOD_ANALYSIS_STATUS_PROCESSING = 2;
  FOOD_ANALYSIS_STATUS_COMPLETED = 3;
  FOOD_ANALYSIS_STATUS_FAILED = 4;
}

message NutrientInfo {
  string name = 1;
  double amount = 2;
  string unit = 3;
  double dv = 4;  // 일일 권장량 대비 비율 0.0~1.0
}

message FoodItem {
  string name = 1;
  double confidence = 2;
  double calorie_kcal = 3;
  double portion_g = 4;
  repeated NutrientInfo nutrients = 5;
}

message FoodAnalysis {
  string id = 1;
  string user_id = 2;
  string image_url = 3;
  FoodAnalysisStatus status = 4;
  double total_calorie_kcal = 5;
  repeated FoodItem food_items = 6;
  string meal_type = 7;
  google.protobuf.Timestamp analyzed_at = 8;
  google.protobuf.Timestamp created_at = 9;
  string error_message = 10;
}
```

## 4. 반영 시 할 일

1. **Proto 파일**: `backend/shared/proto/manpasik.proto` 끝에 위 서비스·메시지 추가 또는 `shared/proto/vision/v1/vision.proto` 신규 후 import.
2. **코드 생성**: `protoc`로 Go/Flutter 스텁 재생성.
3. **vision-service**: `internal/handler/grpc.go`에서 `v1.RegisterVisionServiceServer(grpcServer, visionHandler)` 활성화, `cmd/main.go`에서 동일.
4. **Gateway**: VisionService 메서드 라우팅 및 인증 미들웨어 적용.
5. **Flutter**: 생성된 클라이언트로 Vision 화면 연동.

## 5. 참고

- **기존 서비스 코드**: `backend/services/vision-service/internal/service/vision.go` (AnalyzeFood, GetAnalysis, ListAnalyses, GetDailySummary).
- **핸들러**: `internal/handler/grpc.go` — Proto 메시지와 도메인 구조체 변환만 추가하면 됨.

---

**마지막 업데이트**: 2026-02-12
