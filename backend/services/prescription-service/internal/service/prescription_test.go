package service_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/manpasik/backend/services/prescription-service/internal/repository/memory"
	"github.com/manpasik/backend/services/prescription-service/internal/service"
	"go.uber.org/zap"
)

func setupPrescriptionService() *service.PrescriptionService {
	logger := zap.NewNop()
	prescriptionRepo := memory.NewPrescriptionRepository()
	interactionRepo := memory.NewDrugInteractionRepository()
	tokenRepo := memory.NewTokenRepository()
	return service.NewPrescriptionService(logger, prescriptionRepo, interactionRepo, tokenRepo)
}

func TestCreatePrescription(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()

	meds := []*service.Medication{
		{DrugName: "아모시실린", DrugCode: "AMOX001", Dosage: "500mg", Frequency: "1일 3회", DurationDays: 7, Route: "경구", Instructions: "식후 30분", Quantity: 21},
	}

	p, err := svc.CreatePrescription(ctx, "user-1", "doctor-1", "consult-1", "급성 상기도 감염", "항생제 처방", meds)
	if err != nil {
		t.Fatalf("CreatePrescription 실패: %v", err)
	}
	if p.ID == "" {
		t.Fatal("처방전 ID가 비어 있습니다")
	}
	if p.Status != service.StatusActive {
		t.Fatalf("예상 상태 Active, 실제: %d", p.Status)
	}
	if len(p.Medications) != 1 {
		t.Fatalf("약물 수 예상 1, 실제: %d", len(p.Medications))
	}
}

func TestCreatePrescription_EmptyPatient(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()

	_, err := svc.CreatePrescription(ctx, "", "doctor-1", "", "", "", nil)
	if err == nil {
		t.Fatal("빈 patient_user_id에 에러가 발생해야 합니다")
	}
}

func TestCreatePrescription_EmptyDoctor(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()

	_, err := svc.CreatePrescription(ctx, "user-1", "", "", "", "", nil)
	if err == nil {
		t.Fatal("빈 doctor_id에 에러가 발생해야 합니다")
	}
}

func TestGetPrescription(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()

	meds := []*service.Medication{{DrugName: "타이레놀", DrugCode: "ACET001", Dosage: "500mg", Frequency: "1일 3회"}}
	created, _ := svc.CreatePrescription(ctx, "user-1", "doctor-1", "", "두통", "", meds)

	got, err := svc.GetPrescription(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetPrescription 실패: %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("ID 불일치: %s != %s", got.ID, created.ID)
	}
}

func TestGetPrescription_NotFound(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()

	_, err := svc.GetPrescription(ctx, "nonexistent")
	if err == nil {
		t.Fatal("존재하지 않는 처방전에 에러가 발생해야 합니다")
	}
}

func TestListPrescriptions(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()

	svc.CreatePrescription(ctx, "user-1", "doctor-1", "", "감기", "", []*service.Medication{{DrugName: "약1", DrugCode: "D001"}})
	svc.CreatePrescription(ctx, "user-1", "doctor-2", "", "두통", "", []*service.Medication{{DrugName: "약2", DrugCode: "D002"}})
	svc.CreatePrescription(ctx, "user-2", "doctor-1", "", "복통", "", []*service.Medication{{DrugName: "약3", DrugCode: "D003"}})

	list, total, err := svc.ListPrescriptions(ctx, "user-1", service.StatusUnknown, 10, 0)
	if err != nil {
		t.Fatalf("ListPrescriptions 실패: %v", err)
	}
	if total != 2 {
		t.Fatalf("user-1 처방전 수 예상 2, 실제: %d", total)
	}
	if len(list) != 2 {
		t.Fatalf("반환 수 예상 2, 실제: %d", len(list))
	}
}

func TestUpdatePrescriptionStatus(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()

	meds := []*service.Medication{{DrugName: "약", DrugCode: "D001"}}
	created, _ := svc.CreatePrescription(ctx, "user-1", "doctor-1", "", "진단", "", meds)

	updated, err := svc.UpdatePrescriptionStatus(ctx, created.ID, service.StatusDispensed, "pharmacy-1")
	if err != nil {
		t.Fatalf("UpdatePrescriptionStatus 실패: %v", err)
	}
	if updated.Status != service.StatusDispensed {
		t.Fatalf("예상 상태 Dispensed, 실제: %d", updated.Status)
	}
	if updated.PharmacyID != "pharmacy-1" {
		t.Fatalf("PharmacyID 불일치: %s", updated.PharmacyID)
	}
}

func TestAddMedication(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()

	meds := []*service.Medication{{DrugName: "약1", DrugCode: "D001"}}
	created, _ := svc.CreatePrescription(ctx, "user-1", "doctor-1", "", "진단", "", meds)

	newMed := &service.Medication{DrugName: "약2", DrugCode: "D002", Dosage: "10mg"}
	updated, err := svc.AddMedication(ctx, created.ID, newMed)
	if err != nil {
		t.Fatalf("AddMedication 실패: %v", err)
	}
	if len(updated.Medications) != 2 {
		t.Fatalf("약물 수 예상 2, 실제: %d", len(updated.Medications))
	}
}

func TestRemoveMedication(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()

	med1 := &service.Medication{ID: "med-1", DrugName: "약1", DrugCode: "D001"}
	med2 := &service.Medication{ID: "med-2", DrugName: "약2", DrugCode: "D002"}
	created, _ := svc.CreatePrescription(ctx, "user-1", "doctor-1", "", "진단", "", []*service.Medication{med1, med2})

	updated, err := svc.RemoveMedication(ctx, created.ID, "med-1")
	if err != nil {
		t.Fatalf("RemoveMedication 실패: %v", err)
	}
	if len(updated.Medications) != 1 {
		t.Fatalf("약물 수 예상 1, 실제: %d", len(updated.Medications))
	}
	if updated.Medications[0].ID != "med-2" {
		t.Fatalf("남은 약물 ID 예상 med-2, 실제: %s", updated.Medications[0].ID)
	}
}

func TestRemoveMedication_NotFound(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()

	created, _ := svc.CreatePrescription(ctx, "user-1", "doctor-1", "", "진단", "", []*service.Medication{{DrugName: "약1"}})

	_, err := svc.RemoveMedication(ctx, created.ID, "nonexistent")
	if err == nil {
		t.Fatal("존재하지 않는 약물 제거에 에러가 발생해야 합니다")
	}
}

func TestCheckDrugInteraction(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()

	interactions, err := svc.CheckDrugInteraction(ctx, []string{"WARF001", "ASPR001"})
	if err != nil {
		t.Fatalf("CheckDrugInteraction 실패: %v", err)
	}
	if len(interactions) != 1 {
		t.Fatalf("상호작용 수 예상 1, 실제: %d", len(interactions))
	}
	if interactions[0].Severity != service.SeverityMajor {
		t.Fatalf("심각도 예상 Major, 실제: %d", interactions[0].Severity)
	}
}

func TestCheckDrugInteraction_NoInteraction(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()

	interactions, err := svc.CheckDrugInteraction(ctx, []string{"AMOX001", "ACET001"})
	if err != nil {
		t.Fatalf("CheckDrugInteraction 실패: %v", err)
	}
	if len(interactions) != 0 {
		t.Fatalf("상호작용 수 예상 0, 실제: %d", len(interactions))
	}
}

func TestCheckDrugInteraction_Contraindicated(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()

	interactions, err := svc.CheckDrugInteraction(ctx, []string{"SSRI001", "MAOI001"})
	if err != nil {
		t.Fatalf("CheckDrugInteraction 실패: %v", err)
	}
	if len(interactions) != 1 {
		t.Fatalf("상호작용 수 예상 1, 실제: %d", len(interactions))
	}
	if interactions[0].Severity != service.SeverityContraindicated {
		t.Fatalf("심각도 예상 Contraindicated, 실제: %d", interactions[0].Severity)
	}
}

func TestGetMedicationReminders(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()

	meds := []*service.Medication{
		{DrugName: "약1", DrugCode: "D001", Dosage: "500mg", Frequency: "1일 3회", Instructions: "식후"},
		{DrugName: "약2", DrugCode: "D002", Dosage: "10mg", Frequency: "1일 2회", Instructions: "식전"},
	}
	svc.CreatePrescription(ctx, "user-1", "doctor-1", "", "진단", "", meds)

	reminders, err := svc.GetMedicationReminders(ctx, "user-1")
	if err != nil {
		t.Fatalf("GetMedicationReminders 실패: %v", err)
	}
	// 1일 3회(3) + 1일 2회(2) = 5 알림
	if len(reminders) != 5 {
		t.Fatalf("알림 수 예상 5, 실제: %d", len(reminders))
	}
}

func TestCheckDrugInteraction_LessThanTwo(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()

	interactions, err := svc.CheckDrugInteraction(ctx, []string{"WARF001"})
	if err != nil {
		t.Fatalf("CheckDrugInteraction 실패: %v", err)
	}
	if interactions != nil {
		t.Fatalf("단일 약물에 상호작용이 없어야 합니다")
	}
}

// === Pharmacy Flow Tests ===

func createActivePrescription(t *testing.T, svc *service.PrescriptionService, ctx context.Context) *service.Prescription {
	t.Helper()
	meds := []*service.Medication{{DrugName: "약", DrugCode: "D001", Dosage: "500mg", Frequency: "1일 3회"}}
	p, err := svc.CreatePrescription(ctx, "user-1", "doctor-1", "consult-1", "진단", "메모", meds)
	if err != nil {
		t.Fatalf("처방전 생성 실패: %v", err)
	}
	return p
}

func TestSelectPharmacyAndFulfillment_Success(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()
	p := createActivePrescription(t, svc, ctx)

	err := svc.SelectPharmacyAndFulfillment(ctx, p.ID, "pharm-1", "만파약국", service.FulfillmentPickup, "")
	if err != nil {
		t.Fatalf("SelectPharmacyAndFulfillment 실패: %v", err)
	}

	got, _ := svc.GetPrescription(ctx, p.ID)
	if got.PharmacyID != "pharm-1" {
		t.Fatalf("PharmacyID 예상 pharm-1, 실제: %s", got.PharmacyID)
	}
	if got.PharmacyName != "만파약국" {
		t.Fatalf("PharmacyName 예상 만파약국, 실제: %s", got.PharmacyName)
	}
	if got.FulfillmentType != service.FulfillmentPickup {
		t.Fatalf("FulfillmentType 예상 PICKUP, 실제: %s", got.FulfillmentType)
	}
}

func TestSelectPharmacyAndFulfillment_CourierRequiresAddress(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()
	p := createActivePrescription(t, svc, ctx)

	err := svc.SelectPharmacyAndFulfillment(ctx, p.ID, "pharm-1", "만파약국", service.FulfillmentCourier, "")
	if err == nil {
		t.Fatal("택배 수령 시 빈 주소에 에러가 발생해야 합니다")
	}

	// With address should succeed
	err = svc.SelectPharmacyAndFulfillment(ctx, p.ID, "pharm-1", "만파약국", service.FulfillmentCourier, "서울시 강남구")
	if err != nil {
		t.Fatalf("주소 포함 시 성공해야 합니다: %v", err)
	}
}

func TestSelectPharmacyAndFulfillment_DeliveryRequiresAddress(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()
	p := createActivePrescription(t, svc, ctx)

	err := svc.SelectPharmacyAndFulfillment(ctx, p.ID, "pharm-1", "만파약국", service.FulfillmentDelivery, "")
	if err == nil {
		t.Fatal("배송 수령 시 빈 주소에 에러가 발생해야 합니다")
	}
}

func TestSelectPharmacyAndFulfillment_EmptyPharmacyID(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()
	p := createActivePrescription(t, svc, ctx)

	err := svc.SelectPharmacyAndFulfillment(ctx, p.ID, "", "만파약국", service.FulfillmentPickup, "")
	if err == nil {
		t.Fatal("빈 pharmacy_id에 에러가 발생해야 합니다")
	}
}

func TestSelectPharmacyAndFulfillment_NonActivePrescription(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()
	p := createActivePrescription(t, svc, ctx)

	// Change to dispensed status
	svc.UpdatePrescriptionStatus(ctx, p.ID, service.StatusDispensed, "")

	err := svc.SelectPharmacyAndFulfillment(ctx, p.ID, "pharm-1", "만파약국", service.FulfillmentPickup, "")
	if err == nil {
		t.Fatal("비활성 처방전에 에러가 발생해야 합니다")
	}
}

func TestSendPrescriptionToPharmacy_Success(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()
	p := createActivePrescription(t, svc, ctx)

	err := svc.SelectPharmacyAndFulfillment(ctx, p.ID, "pharm-1", "만파약국", service.FulfillmentPickup, "")
	if err != nil {
		t.Fatalf("약국 설정 실패: %v", err)
	}

	token, err := svc.SendPrescriptionToPharmacy(ctx, p.ID)
	if err != nil {
		t.Fatalf("SendPrescriptionToPharmacy 실패: %v", err)
	}
	if token == nil {
		t.Fatal("토큰이 nil입니다")
	}
	if len(token.Token) != 6 {
		t.Fatalf("토큰 길이 예상 6, 실제: %d", len(token.Token))
	}
	if token.PrescriptionID != p.ID {
		t.Fatalf("PrescriptionID 불일치: %s != %s", token.PrescriptionID, p.ID)
	}
	if token.PharmacyID != "pharm-1" {
		t.Fatalf("PharmacyID 불일치: %s", token.PharmacyID)
	}
	if token.ExpiresAt.Before(time.Now()) {
		t.Fatal("토큰 만료 시간이 과거입니다")
	}

	// Prescription should be updated
	got, _ := svc.GetPrescription(ctx, p.ID)
	if got.FulfillmentToken != token.Token {
		t.Fatalf("처방전 토큰 불일치")
	}
	if got.DispensaryStatus != service.DispensaryPending {
		t.Fatalf("DispensaryStatus 예상 pending, 실제: %s", got.DispensaryStatus)
	}
	if got.SentToPharmacyAt.IsZero() {
		t.Fatal("SentToPharmacyAt이 설정되지 않았습니다")
	}
}

func TestSendPrescriptionToPharmacy_NoPharmacy(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()
	p := createActivePrescription(t, svc, ctx)

	_, err := svc.SendPrescriptionToPharmacy(ctx, p.ID)
	if err == nil {
		t.Fatal("약국 미설정 시 에러가 발생해야 합니다")
	}
}

func TestGetPrescriptionByToken_Valid(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()
	p := createActivePrescription(t, svc, ctx)

	svc.SelectPharmacyAndFulfillment(ctx, p.ID, "pharm-1", "만파약국", service.FulfillmentPickup, "")
	token, _ := svc.SendPrescriptionToPharmacy(ctx, p.ID)

	got, err := svc.GetPrescriptionByToken(ctx, token.Token)
	if err != nil {
		t.Fatalf("GetPrescriptionByToken 실패: %v", err)
	}
	if got.ID != p.ID {
		t.Fatalf("ID 불일치: %s != %s", got.ID, p.ID)
	}
}

func TestGetPrescriptionByToken_InvalidToken(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()

	_, err := svc.GetPrescriptionByToken(ctx, "INVALID")
	if err == nil {
		t.Fatal("존재하지 않는 토큰에 에러가 발생해야 합니다")
	}
}

func TestGetPrescriptionByToken_ExpiredToken(t *testing.T) {
	// We need to test with an expired token by directly creating one in the repo
	logger := zap.NewNop()
	prescriptionRepo := memory.NewPrescriptionRepository()
	interactionRepo := memory.NewDrugInteractionRepository()
	tokenRepo := memory.NewTokenRepository()
	svc := service.NewPrescriptionService(logger, prescriptionRepo, interactionRepo, tokenRepo)
	ctx := context.Background()

	p := createActivePrescription(t, svc, ctx)
	svc.SelectPharmacyAndFulfillment(ctx, p.ID, "pharm-1", "만파약국", service.FulfillmentPickup, "")

	// Create an expired token directly in the repo
	expiredToken := &service.FulfillmentToken{
		Token:          "EXPTKN",
		PrescriptionID: p.ID,
		PharmacyID:     "pharm-1",
		CreatedAt:      time.Now().Add(-48 * time.Hour),
		ExpiresAt:      time.Now().Add(-24 * time.Hour), // Expired 24h ago
		IsUsed:         false,
	}
	tokenRepo.Create(ctx, expiredToken)

	_, err := svc.GetPrescriptionByToken(ctx, "EXPTKN")
	if err == nil {
		t.Fatal("만료된 토큰에 에러가 발생해야 합니다")
	}
}

func TestGetPrescriptionByToken_EmptyToken(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()

	_, err := svc.GetPrescriptionByToken(ctx, "")
	if err == nil {
		t.Fatal("빈 토큰에 에러가 발생해야 합니다")
	}
}

func TestUpdateDispensaryStatus_ValidTransitions(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()
	p := createActivePrescription(t, svc, ctx)

	svc.SelectPharmacyAndFulfillment(ctx, p.ID, "pharm-1", "만파약국", service.FulfillmentPickup, "")
	svc.SendPrescriptionToPharmacy(ctx, p.ID)

	// pending → preparing
	err := svc.UpdateDispensaryStatus(ctx, p.ID, service.DispensaryPreparing)
	if err != nil {
		t.Fatalf("pending→preparing 실패: %v", err)
	}
	got, _ := svc.GetPrescription(ctx, p.ID)
	if got.DispensaryStatus != service.DispensaryPreparing {
		t.Fatalf("예상 preparing, 실제: %s", got.DispensaryStatus)
	}

	// preparing → ready
	err = svc.UpdateDispensaryStatus(ctx, p.ID, service.DispensaryReady)
	if err != nil {
		t.Fatalf("preparing→ready 실패: %v", err)
	}
	got, _ = svc.GetPrescription(ctx, p.ID)
	if got.DispensaryStatus != service.DispensaryReady {
		t.Fatalf("예상 ready, 실제: %s", got.DispensaryStatus)
	}

	// ready → dispensed
	err = svc.UpdateDispensaryStatus(ctx, p.ID, service.DispensaryDispensed)
	if err != nil {
		t.Fatalf("ready→dispensed 실패: %v", err)
	}
	got, _ = svc.GetPrescription(ctx, p.ID)
	if got.DispensaryStatus != service.DispensaryDispensed {
		t.Fatalf("예상 dispensed, 실제: %s", got.DispensaryStatus)
	}
	if got.Status != service.StatusDispensed {
		t.Fatalf("dispensed 시 처방전 상태도 Dispensed여야 합니다, 실제: %d", got.Status)
	}
	if got.DispensedAt.IsZero() {
		t.Fatal("DispensedAt이 설정되지 않았습니다")
	}
}

func TestUpdateDispensaryStatus_InvalidTransitions(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()
	p := createActivePrescription(t, svc, ctx)

	svc.SelectPharmacyAndFulfillment(ctx, p.ID, "pharm-1", "만파약국", service.FulfillmentPickup, "")
	svc.SendPrescriptionToPharmacy(ctx, p.ID)

	// pending → ready (skip preparing)
	err := svc.UpdateDispensaryStatus(ctx, p.ID, service.DispensaryReady)
	if err == nil {
		t.Fatal("pending→ready 전환에 에러가 발생해야 합니다")
	}

	// pending → dispensed (skip preparing, ready)
	err = svc.UpdateDispensaryStatus(ctx, p.ID, service.DispensaryDispensed)
	if err == nil {
		t.Fatal("pending→dispensed 전환에 에러가 발생해야 합니다")
	}

	// Move to preparing, then try invalid transition
	svc.UpdateDispensaryStatus(ctx, p.ID, service.DispensaryPreparing)

	// preparing → dispensed (skip ready)
	err = svc.UpdateDispensaryStatus(ctx, p.ID, service.DispensaryDispensed)
	if err == nil {
		t.Fatal("preparing→dispensed 전환에 에러가 발생해야 합니다")
	}

	// preparing → pending (backward)
	err = svc.UpdateDispensaryStatus(ctx, p.ID, service.DispensaryPending)
	if err == nil {
		t.Fatal("preparing→pending 역전환에 에러가 발생해야 합니다")
	}
}

func TestFulfillmentTokenGeneration(t *testing.T) {
	svc := setupPrescriptionService()
	ctx := context.Background()

	// Generate multiple tokens and check format
	ambiguousChars := "OI10"
	tokens := make(map[string]bool)

	for i := 0; i < 20; i++ {
		p := createActivePrescription(t, svc, ctx)
		svc.SelectPharmacyAndFulfillment(ctx, p.ID, "pharm-1", "만파약국", service.FulfillmentPickup, "")
		token, err := svc.SendPrescriptionToPharmacy(ctx, p.ID)
		if err != nil {
			t.Fatalf("토큰 생성 실패: %v", err)
		}

		// Check length is 6
		if len(token.Token) != 6 {
			t.Fatalf("토큰 길이 예상 6, 실제: %d", len(token.Token))
		}

		// Check no ambiguous characters
		for _, c := range token.Token {
			if strings.ContainsRune(ambiguousChars, c) {
				t.Fatalf("토큰에 모호한 문자가 포함되어 있습니다: %c (토큰: %s)", c, token.Token)
			}
		}

		// Check all characters are valid
		const validChars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
		for _, c := range token.Token {
			if !strings.ContainsRune(validChars, c) {
				t.Fatalf("토큰에 유효하지 않은 문자가 포함되어 있습니다: %c (토큰: %s)", c, token.Token)
			}
		}

		tokens[token.Token] = true
	}

	// Check that we get some variety (not all the same)
	if len(tokens) < 2 {
		t.Fatal("생성된 토큰이 모두 동일합니다 — 랜덤성 문제")
	}
}
