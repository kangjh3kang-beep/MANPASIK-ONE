# Proto Extension Merge Plan

> Sprint 1 → Phase 3 전환을 위한 Proto 병합 계획 | 2026-02-14

## 1. 현재 상태

- `manpasik.proto`: ~3163줄, 21개 gRPC 서비스 정의
- Phase 1 완료: 서비스 로직 + 테스트 (proto 변경 없이 Go 도메인 객체로 구현)
- Phase 2 완료: 4개 Proto Extension Proposal 문서 작성

## 2. 병합 대상 요약

### Agent A (200–249): 의료/예약
| 유형 | 항목 |
|------|------|
| Enum | FacilityType, Specialty, ReservationStatus |
| Message | Region, Facility 확장, Doctor, TimeSlot |
| RPC | ListDoctorsByFacility, GetDoctorAvailability, SelectDoctor, ListRegions |
| Request/Response | 8개 메시지 |

### Agent B (250–299): 처방/약국
| 유형 | 항목 |
|------|------|
| Enum | FulfillmentType, DispensaryStatus, InteractionSeverity |
| Message | Medication 확장, DrugInteraction, MedicationReminder, FulfillmentToken, Prescription 확장 |
| RPC | SelectPharmacyAndFulfillment, SendToPharmacy, GetByToken, UpdateDispensaryStatus, CheckDrugInteraction, GetMedicationReminders |
| Request/Response | 12개 메시지 |

### Agent C (300–349): 데이터 공유/FHIR
| 유형 | 항목 |
|------|------|
| Enum | ConsentType, ConsentStatus, FHIRResourceType, HealthRecordType |
| Message | DataSharingConsent, DataAccessLog, SharedDataBundle, HealthRecord 확장 |
| RPC | ExportToFHIR, ImportFromFHIR, GetHealthSummary, CreateDataSharingConsent, RevokeDataSharingConsent, ListDataSharingConsents, ShareWithProvider, GetDataAccessLog |
| Request/Response | 16개 메시지 |

### Agent D (350–399): 기반 서비스
| 유형 | 항목 |
|------|------|
| Enum | AdminRole, FamilyRole, InvitationStatus, NotificationChannel, NotificationPriority |
| Message | AuditLogDetail, AdminUser 확장, NotificationTemplate, NotificationPreferences 확장, SharingPreferences, SharedHealthSummary |
| RPC | ListAdminsByRegion, GetAuditLogDetails, SendFromTemplate, SetSharingPreferences, ValidateSharingAccess, GetSharedHealthData |
| Request/Response | 12개 메시지 |

## 3. 병합 절차

### Step 1: 준비 (Pre-merge)

```bash
# 1. 브랜치 생성
git checkout -b sprint1/proto-extension

# 2. 현재 proto 백업
cp backend/shared/proto/manpasik.proto backend/shared/proto/manpasik.proto.bak

# 3. 현재 빌드 상태 확인
cd backend && GOWORK=off go build ./services/...
```

### Step 2: Proto 편집

편집 순서 (manpasik.proto 내):

1. **Enum 블록** (파일 상단, 기존 enum 뒤에 추가)
   - Agent A Enums → Agent B Enums → Agent C Enums → Agent D Enums

2. **Message 블록** (관련 서비스 섹션 내에 추가)
   - 기존 메시지 확장: 기존 메시지에 신규 필드 추가 (번호 범위 준수)
   - 신규 메시지: 관련 서비스 섹션 끝에 추가

3. **RPC 블록** (각 service 정의 내에 추가)
   - 각 Agent의 신규 RPC를 해당 서비스에 추가

4. **Request/Response** (RPC 정의 아래에 추가)
   - 각 RPC의 Request/Response 메시지

### Step 3: protoc 컴파일

```bash
# protoc 설치 확인
protoc --version  # >= 3.21

# Go 플러그인 설치 확인
which protoc-gen-go
which protoc-gen-go-grpc

# 컴파일
cd backend/shared/proto
protoc \
  --go_out=../gen/go/v1 --go_opt=paths=source_relative \
  --go-grpc_out=../gen/go/v1 --go-grpc_opt=paths=source_relative \
  -I. \
  -I/usr/local/include \
  manpasik.proto
```

### Step 4: 빌드 & 테스트 검증

```bash
cd backend

# 빌드 검증
GOWORK=off go build ./services/...
# 예상: 21/21 PASS (핸들러에서 신규 RPC는 UnimplementedXxxServer로 자동 처리)

# 테스트 검증
GOWORK=off go test ./services/.../service/... -count=1
# 예상: 21/21 PASS (서비스 로직은 proto와 독립적)
```

### Step 5: Phase 3 핸들러 구현

각 서비스의 `handler/grpc.go`에서:

```go
// 예시: reservation-service handler
func (h *Handler) ListDoctorsByFacility(ctx context.Context, req *pb.ListDoctorsRequest) (*pb.ListDoctorsResponse, error) {
    doctors, err := h.svc.ListDoctorsByFacility(ctx, req.FacilityId, req.Specialty)
    if err != nil {
        return nil, mapError(err)
    }
    return &pb.ListDoctorsResponse{
        Doctors: mapDoctorsToProto(doctors),
    }, nil
}
```

## 4. 롤백 계획

Proto 병합 실패 시:

```bash
# 1. proto 백업 복원
cp backend/shared/proto/manpasik.proto.bak backend/shared/proto/manpasik.proto

# 2. protoc 재컴파일 (원본 기준)
cd backend/shared/proto && protoc ...

# 3. 빌드 확인
cd backend && GOWORK=off go build ./services/...
```

## 5. 예상 변경량

| 항목 | 변경량 |
|------|--------|
| Enum 추가 | 15개 |
| Message 추가/확장 | 20개 |
| RPC 추가 | 24개 |
| Request/Response 추가 | 48개 |
| Proto 증가 라인 | ~800–1000줄 |
| 최종 proto 크기 | ~4000줄 |

## 6. 의존성 매트릭스

Proto 병합 후 핸들러 구현 순서 (의존성 기준):

```
Phase 3-A: 독립 핸들러 (병렬 가능)
  ├── Agent A: reservation-service (ListDoctors, SelectDoctor, ListRegions)
  ├── Agent B: prescription-service (SelectPharmacy, SendToPharmacy, GetByToken, etc.)
  ├── Agent C: health-record-service (ExportFHIR, ImportFHIR, Consent CRUD)
  └── Agent D: admin-service, notification-service, family-service

Phase 3-B: 연동 핸들러 (Phase 3-A 완료 후)
  ├── telemedicine → prescription 연동
  ├── prescription → notification 연동
  ├── health-record → family 연동
  └── measurement → health-record 연동
```
