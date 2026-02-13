package service_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/manpasik/backend/services/health-record-service/internal/repository/memory"
	"github.com/manpasik/backend/services/health-record-service/internal/service"
	"go.uber.org/zap"
)

func setupTestService() *service.HealthRecordService {
	logger := zap.NewNop()
	repo := memory.NewHealthRecordRepository()
	consentRepo := memory.NewConsentRepository()
	accessLogRepo := memory.NewDataAccessLogRepository()
	return service.NewHealthRecordService(logger, repo, consentRepo, accessLogRepo)
}

func TestCreateRecord_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	record, err := svc.CreateRecord(ctx, "user-1", service.RecordTypeVitalSign, "혈당 측정", "공복 혈당", `{"value":"95","unit":"mg/dL"}`, "manpasik")
	if err != nil {
		t.Fatalf("기록 생성 실패: %v", err)
	}
	if record.ID == "" {
		t.Fatal("기록 ID가 비어 있음")
	}
	if record.RecordType != service.RecordTypeVitalSign {
		t.Fatalf("타입 불일치: got %d", record.RecordType)
	}
	if record.Source != "manpasik" {
		t.Fatalf("소스 불일치: got %s", record.Source)
	}
}

func TestCreateRecord_EmptyUserID(t *testing.T) {
	svc := setupTestService()
	_, err := svc.CreateRecord(context.Background(), "", service.RecordTypeVitalSign, "제목", "", "", "")
	if err == nil {
		t.Fatal("빈 user_id에 에러 반환되어야 함")
	}
}

func TestCreateRecord_DefaultSource(t *testing.T) {
	svc := setupTestService()
	record, err := svc.CreateRecord(context.Background(), "user-2", service.RecordTypeCondition, "메모", "내용", "", "")
	if err != nil {
		t.Fatalf("기록 생성 실패: %v", err)
	}
	if record.Source != "manual" {
		t.Fatalf("기본 소스 불일치: got %s, want manual", record.Source)
	}
}

func TestGetRecord(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	created, _ := svc.CreateRecord(ctx, "user-3", service.RecordTypeAllergy, "땅콩 알레르기", "심각도 높음", "", "manual")

	got, err := svc.GetRecord(ctx, created.ID)
	if err != nil {
		t.Fatalf("기록 조회 실패: %v", err)
	}
	if got.Title != "땅콩 알레르기" {
		t.Fatalf("제목 불일치: got %s", got.Title)
	}
}

func TestListRecords(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	svc.CreateRecord(ctx, "user-4", service.RecordTypeVitalSign, "혈당1", "", "", "")
	svc.CreateRecord(ctx, "user-4", service.RecordTypeLabResult, "아스피린", "", "", "")
	svc.CreateRecord(ctx, "user-4", service.RecordTypeVitalSign, "혈당2", "", "", "")

	records, total, err := svc.ListRecords(ctx, "user-4", service.RecordTypeUnknown, 10, 0)
	if err != nil {
		t.Fatalf("목록 조회 실패: %v", err)
	}
	if total != 3 {
		t.Fatalf("총 개수 불일치: got %d, want 3", total)
	}
	if len(records) != 3 {
		t.Fatalf("반환 수 불일치: got %d", len(records))
	}
}

func TestListRecords_TypeFilter(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	svc.CreateRecord(ctx, "user-5", service.RecordTypeVitalSign, "혈당", "", "", "")
	svc.CreateRecord(ctx, "user-5", service.RecordTypeLabResult, "약", "", "", "")
	svc.CreateRecord(ctx, "user-5", service.RecordTypeVitalSign, "혈압", "", "", "")

	records, total, err := svc.ListRecords(ctx, "user-5", service.RecordTypeVitalSign, 10, 0)
	if err != nil {
		t.Fatalf("필터 조회 실패: %v", err)
	}
	if total != 2 {
		t.Fatalf("측정 기록 수 불일치: got %d, want 2", total)
	}
	if len(records) != 2 {
		t.Fatalf("반환 수 불일치: got %d", len(records))
	}
}

func TestUpdateRecord(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	created, _ := svc.CreateRecord(ctx, "user-6", service.RecordTypeAllergy, "두통", "가벼운 두통", "", "")

	updated, err := svc.UpdateRecord(ctx, created.ID, "편두통", "심한 편두통으로 변경", "")
	if err != nil {
		t.Fatalf("기록 업데이트 실패: %v", err)
	}
	if updated.Title != "편두통" {
		t.Fatalf("제목 불일치: got %s", updated.Title)
	}
	if updated.Description != "심한 편두통으로 변경" {
		t.Fatalf("설명 불일치: got %s", updated.Description)
	}
}

func TestDeleteRecord(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	created, _ := svc.CreateRecord(ctx, "user-7", service.RecordTypeCondition, "메모", "삭제 대상", "", "")

	err := svc.DeleteRecord(ctx, created.ID)
	if err != nil {
		t.Fatalf("기록 삭제 실패: %v", err)
	}

	_, err = svc.GetRecord(ctx, created.ID)
	if err == nil {
		t.Fatal("삭제된 기록 조회에 에러 반환되어야 함")
	}
}

func TestExportToFHIR(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	svc.CreateRecord(ctx, "user-8", service.RecordTypeVitalSign, "혈당", "", `{"value":"95"}`, "manpasik")
	svc.CreateRecord(ctx, "user-8", service.RecordTypeAllergy, "땅콩 알레르기", "", "", "manual")
	svc.CreateRecord(ctx, "user-8", service.RecordTypeImmunization, "COVID-19 백신", "", "", "manual")

	bundleJSON, count, resourceTypes, err := svc.ExportToFHIR(ctx, "user-8", nil, nil, nil)
	if err != nil {
		t.Fatalf("FHIR 내보내기 실패: %v", err)
	}
	if count != 3 {
		t.Fatalf("리소스 수 불일치: got %d, want 3", count)
	}
	if len(resourceTypes) == 0 {
		t.Fatal("리소스 타입이 비어 있음")
	}

	// JSON 유효성 검증
	var bundle map[string]interface{}
	if err := json.Unmarshal([]byte(bundleJSON), &bundle); err != nil {
		t.Fatalf("FHIR JSON 파싱 실패: %v", err)
	}
	if bundle["resourceType"] != "Bundle" {
		t.Fatalf("FHIR resourceType 불일치: got %v", bundle["resourceType"])
	}
}

func TestImportFromFHIR(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	fhirBundle := `{
		"resourceType": "Bundle",
		"type": "collection",
		"entry": [
			{
				"resource": {
					"resourceType": "Observation",
					"id": "obs-1",
					"code": {"text": "Blood Glucose"},
					"subject": {"reference": "Patient/user-9"}
				}
			},
			{
				"resource": {
					"resourceType": "Condition",
					"id": "cond-1",
					"code": {"text": "Type 2 Diabetes"}
				}
			}
		]
	}`

	imported, importedCount, skippedCount, errors, err := svc.ImportFromFHIR(ctx, "user-9", fhirBundle)
	if err != nil {
		t.Fatalf("FHIR 가져오기 실패: %v", err)
	}
	if importedCount != 2 {
		t.Fatalf("가져온 수 불일치: got %d, want 2", importedCount)
	}
	if skippedCount != 0 {
		t.Fatalf("건너뛴 수 불일치: got %d, want 0", skippedCount)
	}
	if len(errors) != 0 {
		t.Fatalf("에러 발생: %v", errors)
	}
	if len(imported) != 2 {
		t.Fatalf("반환 수 불일치: got %d", len(imported))
	}

	// 가져온 기록이 조회되는지 확인
	records, total, _ := svc.ListRecords(ctx, "user-9", service.RecordTypeUnknown, 10, 0)
	if total != 2 {
		t.Fatalf("총 기록 수 불일치: got %d, want 2", total)
	}
	if len(records) != 2 {
		t.Fatalf("반환 수 불일치: got %d", len(records))
	}
}

func TestGetHealthSummary(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	svc.CreateRecord(ctx, "user-10", service.RecordTypeVitalSign, "혈당", "", "", "")
	svc.CreateRecord(ctx, "user-10", service.RecordTypeVitalSign, "혈압", "", "", "")
	svc.CreateRecord(ctx, "user-10", service.RecordTypeLabResult, "약", "", "", "")

	total, byType, recent, summary, err := svc.GetHealthSummary(ctx, "user-10", 30)
	if err != nil {
		t.Fatalf("건강 요약 실패: %v", err)
	}
	if total != 3 {
		t.Fatalf("총 기록 수 불일치: got %d, want 3", total)
	}
	if byType["vital_sign"] != 2 {
		t.Fatalf("vital_sign 기록 수 불일치: got %d, want 2", byType["vital_sign"])
	}
	if byType["lab_result"] != 1 {
		t.Fatalf("lab_result 기록 수 불일치: got %d, want 1", byType["lab_result"])
	}
	if len(recent) > 5 {
		t.Fatalf("최근 기록이 5개 초과: got %d", len(recent))
	}
	if summary == "" {
		t.Fatal("요약 텍스트가 비어 있음")
	}
}

// Note: RecordType-to-FHIR mapping is tested implicitly via ExportToFHIR

func TestEndToEnd_HealthRecordFlow(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	// 1. 다양한 기록 생성
	r1, _ := svc.CreateRecord(ctx, "user-e2e", service.RecordTypeVitalSign, "공복 혈당", "아침 측정", `{"value":"110","unit":"mg/dL"}`, "manpasik")
	svc.CreateRecord(ctx, "user-e2e", service.RecordTypeLabResult, "메트포르민 500mg", "아침 식후", "", "manual")
	svc.CreateRecord(ctx, "user-e2e", service.RecordTypeVitalSign, "혈압 측정", "수축기/이완기", `{"systolic":"120","diastolic":"80"}`, "manpasik")

	// 2. 기록 업데이트
	svc.UpdateRecord(ctx, r1.ID, "공복 혈당 (수정)", "재측정", `{"value":"105","unit":"mg/dL"}`)

	// 3. 건강 요약
	total, byType, _, _, _ := svc.GetHealthSummary(ctx, "user-e2e", 30)
	if total != 3 {
		t.Fatalf("총 기록 수 불일치: got %d, want 3", total)
	}
	if byType["vital_sign"] != 2 {
		t.Fatalf("vital_sign 수 불일치: got %d, want 2", byType["vital_sign"])
	}
	if byType["lab_result"] != 1 {
		t.Fatalf("lab_result 수 불일치: got %d, want 1", byType["lab_result"])
	}

	// 4. FHIR 내보내기
	bundleJSON, count, _, _ := svc.ExportToFHIR(ctx, "user-e2e", nil, nil, nil)
	if count != 3 {
		t.Fatalf("FHIR 리소스 수 불일치: got %d, want 3", count)
	}

	// 5. 다른 사용자에게 FHIR 가져오기
	_, importedCount, _, _, _ := svc.ImportFromFHIR(ctx, "user-import", bundleJSON)
	if importedCount != 3 {
		t.Fatalf("가져온 수 불일치: got %d, want 3", importedCount)
	}

	// 6. 기록 삭제
	svc.DeleteRecord(ctx, r1.ID)
	total2, _, _, _, _ := svc.GetHealthSummary(ctx, "user-e2e", 30)
	if total2 != 2 {
		t.Fatalf("삭제 후 기록 수 불일치: got %d, want 2", total2)
	}
}

// ============================================================================
// 데이터 공유 동의 테스트
// ============================================================================

func TestCreateDataSharingConsent(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	t.Run("happy path", func(t *testing.T) {
		consent := &service.DataSharingConsent{
			UserID:       "user-consent-1",
			ProviderID:   "provider-1",
			ProviderName: "서울병원",
			ConsentType:  service.ConsentMeasurementShare,
			Scope:        []string{"blood_glucose", "blood_pressure_systolic"},
			Purpose:      "treatment",
		}

		result, err := svc.CreateDataSharingConsent(ctx, consent)
		if err != nil {
			t.Fatalf("동의 생성 실패: %v", err)
		}
		if result.ID == "" {
			t.Fatal("동의 ID가 비어 있음")
		}
		if result.Status != service.ConsentActive {
			t.Fatalf("상태 불일치: got %s, want active", result.Status)
		}
		if result.GrantedAt.IsZero() {
			t.Fatal("GrantedAt이 설정되지 않음")
		}
		if result.ExpiresAt.IsZero() {
			t.Fatal("ExpiresAt이 설정되지 않음")
		}
		// 기본 만료일은 1년 후
		expectedExpiry := result.GrantedAt.AddDate(1, 0, 0)
		if result.ExpiresAt.Sub(expectedExpiry) > time.Second {
			t.Fatalf("기본 만료일 불일치: got %v, want ~%v", result.ExpiresAt, expectedExpiry)
		}
	})

	t.Run("validation - missing UserID", func(t *testing.T) {
		consent := &service.DataSharingConsent{
			ProviderID:  "provider-1",
			ConsentType: service.ConsentMeasurementShare,
			Scope:       []string{"blood_glucose"},
		}
		_, err := svc.CreateDataSharingConsent(ctx, consent)
		if err == nil {
			t.Fatal("빈 UserID에 에러 반환되어야 함")
		}
	})

	t.Run("validation - missing ProviderID", func(t *testing.T) {
		consent := &service.DataSharingConsent{
			UserID:      "user-consent-1",
			ConsentType: service.ConsentMeasurementShare,
			Scope:       []string{"blood_glucose"},
		}
		_, err := svc.CreateDataSharingConsent(ctx, consent)
		if err == nil {
			t.Fatal("빈 ProviderID에 에러 반환되어야 함")
		}
	})

	t.Run("validation - missing ConsentType", func(t *testing.T) {
		consent := &service.DataSharingConsent{
			UserID:     "user-consent-1",
			ProviderID: "provider-1",
			Scope:      []string{"blood_glucose"},
		}
		_, err := svc.CreateDataSharingConsent(ctx, consent)
		if err == nil {
			t.Fatal("빈 ConsentType에 에러 반환되어야 함")
		}
	})

	t.Run("validation - empty Scope", func(t *testing.T) {
		consent := &service.DataSharingConsent{
			UserID:      "user-consent-1",
			ProviderID:  "provider-1",
			ConsentType: service.ConsentMeasurementShare,
			Scope:       []string{},
		}
		_, err := svc.CreateDataSharingConsent(ctx, consent)
		if err == nil {
			t.Fatal("빈 Scope에 에러 반환되어야 함")
		}
	})
}

func TestRevokeDataSharingConsent(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	t.Run("revoke active consent", func(t *testing.T) {
		consent := &service.DataSharingConsent{
			UserID:       "user-revoke-1",
			ProviderID:   "provider-1",
			ProviderName: "서울병원",
			ConsentType:  service.ConsentRecordShare,
			Scope:        []string{"blood_glucose"},
			Purpose:      "treatment",
		}
		created, err := svc.CreateDataSharingConsent(ctx, consent)
		if err != nil {
			t.Fatalf("동의 생성 실패: %v", err)
		}

		err = svc.RevokeDataSharingConsent(ctx, created.ID, "환자 요청")
		if err != nil {
			t.Fatalf("동의 철회 실패: %v", err)
		}
	})

	t.Run("double revoke error", func(t *testing.T) {
		consent := &service.DataSharingConsent{
			UserID:       "user-revoke-2",
			ProviderID:   "provider-2",
			ProviderName: "부산병원",
			ConsentType:  service.ConsentFullAccess,
			Scope:        []string{"blood_glucose", "heart_rate"},
			Purpose:      "research",
		}
		created, err := svc.CreateDataSharingConsent(ctx, consent)
		if err != nil {
			t.Fatalf("동의 생성 실패: %v", err)
		}

		err = svc.RevokeDataSharingConsent(ctx, created.ID, "첫 번째 철회")
		if err != nil {
			t.Fatalf("첫 번째 철회 실패: %v", err)
		}

		err = svc.RevokeDataSharingConsent(ctx, created.ID, "두 번째 철회")
		if err == nil {
			t.Fatal("이미 철회된 동의에 대해 에러가 반환되어야 함")
		}
	})
}

func TestListDataSharingConsents(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	// user-list-1에 2개 동의 생성
	consent1 := &service.DataSharingConsent{
		UserID:      "user-list-1",
		ProviderID:  "provider-a",
		ConsentType: service.ConsentMeasurementShare,
		Scope:       []string{"blood_glucose"},
		Purpose:     "treatment",
	}
	consent2 := &service.DataSharingConsent{
		UserID:      "user-list-1",
		ProviderID:  "provider-b",
		ConsentType: service.ConsentRecordShare,
		Scope:       []string{"heart_rate"},
		Purpose:     "research",
	}
	// 다른 사용자 동의
	consent3 := &service.DataSharingConsent{
		UserID:      "user-list-2",
		ProviderID:  "provider-a",
		ConsentType: service.ConsentFullAccess,
		Scope:       []string{"blood_glucose"},
		Purpose:     "emergency",
	}

	svc.CreateDataSharingConsent(ctx, consent1)
	svc.CreateDataSharingConsent(ctx, consent2)
	svc.CreateDataSharingConsent(ctx, consent3)

	consents, err := svc.ListDataSharingConsents(ctx, "user-list-1")
	if err != nil {
		t.Fatalf("동의 목록 조회 실패: %v", err)
	}
	if len(consents) != 2 {
		t.Fatalf("동의 수 불일치: got %d, want 2", len(consents))
	}

	// user-list-2는 1개만 있어야 함
	consents2, err := svc.ListDataSharingConsents(ctx, "user-list-2")
	if err != nil {
		t.Fatalf("동의 목록 조회 실패: %v", err)
	}
	if len(consents2) != 1 {
		t.Fatalf("동의 수 불일치: got %d, want 1", len(consents2))
	}
}

func TestCheckAccess(t *testing.T) {
	logger := zap.NewNop()
	repo := memory.NewHealthRecordRepository()
	consentRepo := memory.NewConsentRepository()
	accessLogRepo := memory.NewDataAccessLogRepository()
	svc := service.NewHealthRecordService(logger, repo, consentRepo, accessLogRepo)
	ctx := context.Background()

	t.Run("active consent allows access", func(t *testing.T) {
		consent := &service.DataSharingConsent{
			UserID:      "user-access-1",
			ProviderID:  "provider-access-1",
			ConsentType: service.ConsentMeasurementShare,
			Scope:       []string{"blood_glucose", "heart_rate"},
			Purpose:     "treatment",
		}
		_, err := svc.CreateDataSharingConsent(ctx, consent)
		if err != nil {
			t.Fatalf("동의 생성 실패: %v", err)
		}

		allowed, err := consentRepo.CheckAccess(ctx, "user-access-1", "provider-access-1", "blood_glucose")
		if err != nil {
			t.Fatalf("접근 확인 실패: %v", err)
		}
		if !allowed {
			t.Fatal("활성 동의가 있으므로 접근이 허용되어야 함")
		}
	})

	t.Run("revoked consent denies access", func(t *testing.T) {
		consent := &service.DataSharingConsent{
			UserID:      "user-access-2",
			ProviderID:  "provider-access-2",
			ConsentType: service.ConsentMeasurementShare,
			Scope:       []string{"blood_glucose"},
			Purpose:     "treatment",
		}
		created, err := svc.CreateDataSharingConsent(ctx, consent)
		if err != nil {
			t.Fatalf("동의 생성 실패: %v", err)
		}

		err = svc.RevokeDataSharingConsent(ctx, created.ID, "철회 테스트")
		if err != nil {
			t.Fatalf("동의 철회 실패: %v", err)
		}

		allowed, err := consentRepo.CheckAccess(ctx, "user-access-2", "provider-access-2", "blood_glucose")
		if err != nil {
			t.Fatalf("접근 확인 실패: %v", err)
		}
		if allowed {
			t.Fatal("철회된 동의이므로 접근이 거부되어야 함")
		}
	})

	t.Run("no consent denies access", func(t *testing.T) {
		allowed, err := consentRepo.CheckAccess(ctx, "no-user", "no-provider", "blood_glucose")
		if err != nil {
			t.Fatalf("접근 확인 실패: %v", err)
		}
		if allowed {
			t.Fatal("동의가 없으므로 접근이 거부되어야 함")
		}
	})
}

func TestShareWithProvider(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	// 건강 기록 생성
	svc.CreateRecord(ctx, "user-share-1", service.RecordTypeVitalSign, "혈당 측정", "", `{"value":"95"}`, "manpasik")
	svc.CreateRecord(ctx, "user-share-1", service.RecordTypeVitalSign, "혈압 측정", "", `{"systolic":"120"}`, "manpasik")

	// 동의 생성
	consent := &service.DataSharingConsent{
		UserID:       "user-share-1",
		ProviderID:   "provider-share-1",
		ProviderName: "서울병원",
		ConsentType:  service.ConsentMeasurementShare,
		Scope:        []string{"blood_glucose"},
		Purpose:      "treatment",
	}
	created, err := svc.CreateDataSharingConsent(ctx, consent)
	if err != nil {
		t.Fatalf("동의 생성 실패: %v", err)
	}

	// 공유 실행
	bundle, err := svc.ShareWithProvider(ctx, created.ID)
	if err != nil {
		t.Fatalf("데이터 공유 실패: %v", err)
	}

	if bundle.ConsentID != created.ID {
		t.Fatalf("ConsentID 불일치: got %s, want %s", bundle.ConsentID, created.ID)
	}
	if bundle.ProviderID != "provider-share-1" {
		t.Fatalf("ProviderID 불일치: got %s", bundle.ProviderID)
	}
	if bundle.FHIRBundleJSON == "" {
		t.Fatal("FHIR Bundle JSON이 비어 있음")
	}
	if bundle.ResourceCount == 0 {
		t.Fatal("ResourceCount가 0")
	}

	// FHIR Bundle JSON 검증
	var fhirBundle map[string]interface{}
	if err := json.Unmarshal([]byte(bundle.FHIRBundleJSON), &fhirBundle); err != nil {
		t.Fatalf("FHIR JSON 파싱 실패: %v", err)
	}
	if fhirBundle["resourceType"] != "Bundle" {
		t.Fatalf("resourceType 불일치: got %v", fhirBundle["resourceType"])
	}

	// LOINC 코드 확인 — blood_glucose의 LOINC 코드는 15074-8
	entries, ok := fhirBundle["entry"].([]interface{})
	if !ok || len(entries) == 0 {
		t.Fatal("FHIR Bundle에 entry가 없음")
	}
	firstEntry := entries[0].(map[string]interface{})
	resource := firstEntry["resource"].(map[string]interface{})
	code := resource["code"].(map[string]interface{})
	codings := code["coding"].([]interface{})
	firstCoding := codings[0].(map[string]interface{})
	if firstCoding["code"] != "15074-8" {
		t.Fatalf("LOINC 코드 불일치: got %v, want 15074-8", firstCoding["code"])
	}

	// 접근 로그 확인
	logs, total, err := svc.GetDataAccessLog(ctx, "user-share-1", 10, 0)
	if err != nil {
		t.Fatalf("접근 로그 조회 실패: %v", err)
	}
	if total != 1 {
		t.Fatalf("접근 로그 수 불일치: got %d, want 1", total)
	}
	if logs[0].Action != "share" {
		t.Fatalf("접근 로그 Action 불일치: got %s, want share", logs[0].Action)
	}
}

func TestMeasurementToFHIRObservation(t *testing.T) {
	measuredAt := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)

	t.Run("known biomarker - blood_glucose", func(t *testing.T) {
		obs := service.MeasurementToFHIRObservation("blood_glucose", 95.0, "mg/dL", "Patient/user-1", measuredAt)

		if obs["resourceType"] != "Observation" {
			t.Fatalf("resourceType 불일치: got %v", obs["resourceType"])
		}
		if obs["status"] != "final" {
			t.Fatalf("status 불일치: got %v", obs["status"])
		}

		// JSON 왕복을 통해 타입을 정규화
		obsJSON, err := json.Marshal(obs)
		if err != nil {
			t.Fatalf("JSON 직렬화 실패: %v", err)
		}
		var obsMap map[string]interface{}
		if err := json.Unmarshal(obsJSON, &obsMap); err != nil {
			t.Fatalf("JSON 역직렬화 실패: %v", err)
		}

		code := obsMap["code"].(map[string]interface{})
		codings := code["coding"].([]interface{})
		firstCoding := codings[0].(map[string]interface{})
		if firstCoding["code"] != "15074-8" {
			t.Fatalf("LOINC 코드 불일치: got %v, want 15074-8", firstCoding["code"])
		}
		if firstCoding["display"] != "Glucose [Moles/volume] in Blood" {
			t.Fatalf("LOINC display 불일치: got %v", firstCoding["display"])
		}

		valueQuantity := obsMap["valueQuantity"].(map[string]interface{})
		if valueQuantity["value"] != 95.0 {
			t.Fatalf("value 불일치: got %v, want 95.0", valueQuantity["value"])
		}
		if valueQuantity["unit"] != "mg/dL" {
			t.Fatalf("unit 불일치: got %v, want mg/dL", valueQuantity["unit"])
		}

		subject := obsMap["subject"].(map[string]interface{})
		if subject["reference"] != "Patient/user-1" {
			t.Fatalf("subject reference 불일치: got %v", subject["reference"])
		}
	})

	t.Run("unknown biomarker", func(t *testing.T) {
		obs := service.MeasurementToFHIRObservation("unknown_marker", 42.0, "units", "Patient/user-2", measuredAt)

		obsJSON, err := json.Marshal(obs)
		if err != nil {
			t.Fatalf("JSON 직렬화 실패: %v", err)
		}
		var obsMap map[string]interface{}
		if err := json.Unmarshal(obsJSON, &obsMap); err != nil {
			t.Fatalf("JSON 역직렬화 실패: %v", err)
		}

		code := obsMap["code"].(map[string]interface{})
		codings := code["coding"].([]interface{})
		firstCoding := codings[0].(map[string]interface{})
		if firstCoding["code"] != "unknown" {
			t.Fatalf("알 수 없는 바이오마커의 LOINC 코드는 'unknown'이어야 함: got %v", firstCoding["code"])
		}
	})
}

func TestBuildFHIRBundle(t *testing.T) {
	measuredAt := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)

	obs1 := service.MeasurementToFHIRObservation("blood_glucose", 95.0, "mg/dL", "Patient/user-1", measuredAt)
	obs2 := service.MeasurementToFHIRObservation("heart_rate", 72.0, "beats/min", "Patient/user-1", measuredAt)

	bundleJSON, err := service.BuildFHIRBundle([]map[string]interface{}{obs1, obs2})
	if err != nil {
		t.Fatalf("FHIR Bundle 생성 실패: %v", err)
	}

	// JSON 유효성 검증
	var bundle map[string]interface{}
	if err := json.Unmarshal([]byte(bundleJSON), &bundle); err != nil {
		t.Fatalf("FHIR Bundle JSON 파싱 실패: %v", err)
	}

	// 구조 검증
	if bundle["resourceType"] != "Bundle" {
		t.Fatalf("resourceType 불일치: got %v", bundle["resourceType"])
	}
	if bundle["type"] != "collection" {
		t.Fatalf("type 불일치: got %v", bundle["type"])
	}
	total, ok := bundle["total"].(float64)
	if !ok || int(total) != 2 {
		t.Fatalf("total 불일치: got %v, want 2", bundle["total"])
	}

	entries, ok := bundle["entry"].([]interface{})
	if !ok {
		t.Fatal("entry가 배열이 아님")
	}
	if len(entries) != 2 {
		t.Fatalf("entry 수 불일치: got %d, want 2", len(entries))
	}

	// 각 entry에 resource가 있는지 확인
	for i, e := range entries {
		entry, ok := e.(map[string]interface{})
		if !ok {
			t.Fatalf("entry[%d]가 객체가 아님", i)
		}
		resource, ok := entry["resource"].(map[string]interface{})
		if !ok {
			t.Fatalf("entry[%d].resource가 객체가 아님", i)
		}
		if resource["resourceType"] != "Observation" {
			t.Fatalf("entry[%d].resource.resourceType 불일치: got %v", i, resource["resourceType"])
		}
	}

	// 빈 observations 테스트
	emptyJSON, err := service.BuildFHIRBundle([]map[string]interface{}{})
	if err != nil {
		t.Fatalf("빈 FHIR Bundle 생성 실패: %v", err)
	}
	var emptyBundle map[string]interface{}
	if err := json.Unmarshal([]byte(emptyJSON), &emptyBundle); err != nil {
		t.Fatalf("빈 FHIR Bundle JSON 파싱 실패: %v", err)
	}
	if emptyBundle["resourceType"] != "Bundle" {
		t.Fatalf("빈 Bundle resourceType 불일치: got %v", emptyBundle["resourceType"])
	}
}
