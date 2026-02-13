// Package service는 prescription-service의 비즈니스 로직을 구현합니다.
package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/manpasik/backend/shared/errors"
	"go.uber.org/zap"
)

// PrescriptionStatus는 처방전 상태입니다.
type PrescriptionStatus int

const (
	StatusUnknown   PrescriptionStatus = iota
	StatusDraft
	StatusActive
	StatusDispensed
	StatusCompleted
	StatusCancelled
	StatusExpired
)

// InteractionSeverity는 약물 상호작용 심각도입니다.
type InteractionSeverity int

const (
	SeverityUnknown         InteractionSeverity = iota
	SeverityNone
	SeverityMinor
	SeverityModerate
	SeverityMajor
	SeverityContraindicated
)

// FulfillmentType represents how the patient receives medication
type FulfillmentType string

const (
	FulfillmentPickup   FulfillmentType = "PICKUP"   // Patient picks up at pharmacy
	FulfillmentCourier  FulfillmentType = "COURIER"   // Courier delivery
	FulfillmentDelivery FulfillmentType = "DELIVERY"  // Standard mail delivery
)

// DispensaryStatus tracks pharmacy preparation progress
type DispensaryStatus string

const (
	DispensaryPending    DispensaryStatus = "pending"
	DispensaryPreparing  DispensaryStatus = "preparing"
	DispensaryReady      DispensaryStatus = "ready"
	DispensaryDispensed  DispensaryStatus = "dispensed"
)

// Prescription은 처방전 도메인 객체입니다.
type Prescription struct {
	ID              string
	PatientUserID   string
	DoctorID        string
	ConsultationID  string
	Status          PrescriptionStatus
	Medications     []*Medication
	Diagnosis       string
	Notes           string
	PharmacyID      string
	PharmacyName    string
	FulfillmentType FulfillmentType
	ShippingAddress string
	FulfillmentToken string
	DispensaryStatus DispensaryStatus
	PrescribedAt    time.Time
	ExpiresAt       time.Time
	SentToPharmacyAt time.Time
	DispensedAt     time.Time
	CreatedAt       time.Time
}

// Medication은 약물 도메인 객체입니다.
type Medication struct {
	ID               string
	DrugName         string
	DrugCode         string
	Dosage           string
	Frequency        string
	DurationDays     int
	Route            string
	Instructions     string
	Quantity         int
	RefillsRemaining int
	IsGenericAllowed bool
}

// DrugInteraction은 약물 상호작용 정보입니다.
type DrugInteraction struct {
	DrugA          string
	DrugB          string
	Severity       InteractionSeverity
	Description    string
	Recommendation string
}

// MedicationReminder는 복약 알림입니다.
type MedicationReminder struct {
	PrescriptionID string
	MedicationID   string
	DrugName       string
	Dosage         string
	TimeOfDay      string
	Instructions   string
	IsTaken        bool
}

// PrescriptionRepository는 처방전 저장소 인터페이스입니다.
type PrescriptionRepository interface {
	Save(ctx context.Context, p *Prescription) error
	FindByID(ctx context.Context, id string) (*Prescription, error)
	FindByUserID(ctx context.Context, userID string, statusFilter PrescriptionStatus, limit, offset int) ([]*Prescription, int, error)
	Update(ctx context.Context, p *Prescription) error
}

// DrugInteractionRepository는 약물 상호작용 저장소 인터페이스입니다.
type DrugInteractionRepository interface {
	CheckInteractions(ctx context.Context, drugCodes []string) ([]*DrugInteraction, error)
}

// FulfillmentToken은 처방전을 약국에 전달할 때 사용하는 토큰입니다.
type FulfillmentToken struct {
	Token          string
	PrescriptionID string
	PharmacyID     string
	CreatedAt      time.Time
	ExpiresAt      time.Time
	IsUsed         bool
	UsedAt         time.Time
}

// TokenRepository는 조제 토큰 저장소 인터페이스입니다.
type TokenRepository interface {
	Create(ctx context.Context, token *FulfillmentToken) error
	GetByToken(ctx context.Context, token string) (*FulfillmentToken, error)
	MarkUsed(ctx context.Context, token string) error
}

// EventPublisher is an optional interface for publishing domain events.
type EventPublisher interface {
	Publish(ctx context.Context, event interface{}) error
}

// PrescriptionService는 처방전 서비스 핵심 로직입니다.
type PrescriptionService struct {
	log              *zap.Logger
	prescriptionRepo PrescriptionRepository
	interactionRepo  DrugInteractionRepository
	tokenRepo        TokenRepository
	eventPub         EventPublisher
}

// NewPrescriptionService는 PrescriptionService를 생성합니다.
// eventPub is optional and may be nil.
func NewPrescriptionService(log *zap.Logger, prescriptionRepo PrescriptionRepository, interactionRepo DrugInteractionRepository, tokenRepo TokenRepository, eventPub ...EventPublisher) *PrescriptionService {
	svc := &PrescriptionService{
		log:              log,
		prescriptionRepo: prescriptionRepo,
		interactionRepo:  interactionRepo,
		tokenRepo:        tokenRepo,
	}
	if len(eventPub) > 0 && eventPub[0] != nil {
		svc.eventPub = eventPub[0]
	}
	return svc
}

// CreatePrescription은 새 처방전을 생성합니다.
func (s *PrescriptionService) CreatePrescription(ctx context.Context, patientUserID, doctorID, consultationID, diagnosis, notes string, medications []*Medication) (*Prescription, error) {
	if patientUserID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "patient_user_id는 필수입니다")
	}
	if doctorID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "doctor_id는 필수입니다")
	}

	now := time.Now()
	for _, m := range medications {
		if m.ID == "" {
			m.ID = uuid.New().String()
		}
	}

	prescription := &Prescription{
		ID:             uuid.New().String(),
		PatientUserID:  patientUserID,
		DoctorID:       doctorID,
		ConsultationID: consultationID,
		Status:         StatusActive,
		Medications:    medications,
		Diagnosis:      diagnosis,
		Notes:          notes,
		PrescribedAt:   now,
		ExpiresAt:      now.Add(30 * 24 * time.Hour), // 기본 30일
		CreatedAt:      now,
	}

	if err := s.prescriptionRepo.Save(ctx, prescription); err != nil {
		return nil, fmt.Errorf("처방전 저장 실패: %w", err)
	}

	s.log.Info("처방전 생성 완료",
		zap.String("prescription_id", prescription.ID),
		zap.String("patient_user_id", patientUserID),
		zap.Int("medication_count", len(medications)),
	)

	// Publish prescription.created event
	if s.eventPub != nil {
		_ = s.eventPub.Publish(ctx, map[string]interface{}{
			"type":            "prescription.created",
			"prescription_id": prescription.ID,
			"user_id":         prescription.PatientUserID,
			"doctor_id":       prescription.DoctorID,
			"diagnosis":       prescription.Diagnosis,
		})
	}

	return prescription, nil
}

// GetPrescription은 처방전을 조회합니다.
func (s *PrescriptionService) GetPrescription(ctx context.Context, prescriptionID string) (*Prescription, error) {
	if prescriptionID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "prescription_id는 필수입니다")
	}
	return s.prescriptionRepo.FindByID(ctx, prescriptionID)
}

// ListPrescriptions는 환자의 처방전 목록을 조회합니다.
func (s *PrescriptionService) ListPrescriptions(ctx context.Context, patientUserID string, statusFilter PrescriptionStatus, limit, offset int) ([]*Prescription, int, error) {
	if patientUserID == "" {
		return nil, 0, apperrors.New(apperrors.ErrInvalidInput, "patient_user_id는 필수입니다")
	}
	if limit <= 0 {
		limit = 20
	}
	return s.prescriptionRepo.FindByUserID(ctx, patientUserID, statusFilter, limit, offset)
}

// UpdatePrescriptionStatus는 처방전 상태를 변경합니다.
func (s *PrescriptionService) UpdatePrescriptionStatus(ctx context.Context, prescriptionID string, newStatus PrescriptionStatus, pharmacyID string) (*Prescription, error) {
	if prescriptionID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "prescription_id는 필수입니다")
	}

	prescription, err := s.prescriptionRepo.FindByID(ctx, prescriptionID)
	if err != nil {
		return nil, err
	}

	prescription.Status = newStatus
	if pharmacyID != "" {
		prescription.PharmacyID = pharmacyID
	}
	if newStatus == StatusDispensed {
		prescription.DispensedAt = time.Now()
	}

	if err := s.prescriptionRepo.Update(ctx, prescription); err != nil {
		return nil, fmt.Errorf("처방전 상태 업데이트 실패: %w", err)
	}

	s.log.Info("처방전 상태 변경",
		zap.String("prescription_id", prescriptionID),
		zap.Int("new_status", int(newStatus)),
	)

	return prescription, nil
}

// AddMedication은 처방전에 약물을 추가합니다.
func (s *PrescriptionService) AddMedication(ctx context.Context, prescriptionID string, medication *Medication) (*Prescription, error) {
	if prescriptionID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "prescription_id는 필수입니다")
	}
	if medication == nil || medication.DrugName == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "약물 정보는 필수입니다")
	}

	prescription, err := s.prescriptionRepo.FindByID(ctx, prescriptionID)
	if err != nil {
		return nil, err
	}

	if medication.ID == "" {
		medication.ID = uuid.New().String()
	}
	prescription.Medications = append(prescription.Medications, medication)

	if err := s.prescriptionRepo.Update(ctx, prescription); err != nil {
		return nil, fmt.Errorf("약물 추가 실패: %w", err)
	}

	s.log.Info("약물 추가",
		zap.String("prescription_id", prescriptionID),
		zap.String("drug_name", medication.DrugName),
	)

	return prescription, nil
}

// RemoveMedication은 처방전에서 약물을 제거합니다.
func (s *PrescriptionService) RemoveMedication(ctx context.Context, prescriptionID, medicationID string) (*Prescription, error) {
	if prescriptionID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "prescription_id는 필수입니다")
	}
	if medicationID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "medication_id는 필수입니다")
	}

	prescription, err := s.prescriptionRepo.FindByID(ctx, prescriptionID)
	if err != nil {
		return nil, err
	}

	found := false
	var updated []*Medication
	for _, m := range prescription.Medications {
		if m.ID == medicationID {
			found = true
			continue
		}
		updated = append(updated, m)
	}

	if !found {
		return nil, apperrors.New(apperrors.ErrNotFound, "해당 약물을 찾을 수 없습니다")
	}

	prescription.Medications = updated
	if err := s.prescriptionRepo.Update(ctx, prescription); err != nil {
		return nil, fmt.Errorf("약물 제거 실패: %w", err)
	}

	s.log.Info("약물 제거",
		zap.String("prescription_id", prescriptionID),
		zap.String("medication_id", medicationID),
	)

	return prescription, nil
}

// CheckDrugInteraction은 약물 간 상호작용을 검사합니다.
func (s *PrescriptionService) CheckDrugInteraction(ctx context.Context, drugCodes []string) ([]*DrugInteraction, error) {
	if len(drugCodes) < 2 {
		return nil, nil // 2개 미만이면 상호작용 없음
	}
	return s.interactionRepo.CheckInteractions(ctx, drugCodes)
}

// GetMedicationReminders는 특정 날짜의 복약 알림 목록을 반환합니다.
func (s *PrescriptionService) GetMedicationReminders(ctx context.Context, patientUserID string) ([]*MedicationReminder, error) {
	if patientUserID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "patient_user_id는 필수입니다")
	}

	prescriptions, _, err := s.prescriptionRepo.FindByUserID(ctx, patientUserID, StatusActive, 100, 0)
	if err != nil {
		return nil, err
	}

	var reminders []*MedicationReminder
	for _, p := range prescriptions {
		for _, m := range p.Medications {
			times := parseFrequencyToTimes(m.Frequency)
			for _, t := range times {
				reminders = append(reminders, &MedicationReminder{
					PrescriptionID: p.ID,
					MedicationID:   m.ID,
					DrugName:       m.DrugName,
					Dosage:         m.Dosage,
					TimeOfDay:      t,
					Instructions:   m.Instructions,
					IsTaken:        false,
				})
			}
		}
	}

	return reminders, nil
}

// SelectPharmacyAndFulfillment은 처방전에 약국과 수령 방식을 설정합니다.
func (s *PrescriptionService) SelectPharmacyAndFulfillment(ctx context.Context, prescriptionID, pharmacyID, pharmacyName string, fulfillment FulfillmentType, shippingAddr string) error {
	if prescriptionID == "" {
		return apperrors.New(apperrors.ErrInvalidInput, "prescription_id는 필수입니다")
	}
	if pharmacyID == "" {
		return apperrors.New(apperrors.ErrInvalidInput, "pharmacy_id는 필수입니다")
	}
	if pharmacyName == "" {
		return apperrors.New(apperrors.ErrInvalidInput, "pharmacy_name은 필수입니다")
	}
	if fulfillment == "" {
		return apperrors.New(apperrors.ErrInvalidInput, "fulfillment_type은 필수입니다")
	}
	if (fulfillment == FulfillmentCourier || fulfillment == FulfillmentDelivery) && shippingAddr == "" {
		return apperrors.New(apperrors.ErrInvalidInput, "배송/택배 수령 시 배송 주소는 필수입니다")
	}

	prescription, err := s.prescriptionRepo.FindByID(ctx, prescriptionID)
	if err != nil {
		return err
	}

	if prescription.Status != StatusActive {
		return apperrors.New(apperrors.ErrInvalidInput, "활성 상태의 처방전만 약국 설정이 가능합니다")
	}

	prescription.PharmacyID = pharmacyID
	prescription.PharmacyName = pharmacyName
	prescription.FulfillmentType = fulfillment
	prescription.ShippingAddress = shippingAddr

	if err := s.prescriptionRepo.Update(ctx, prescription); err != nil {
		return fmt.Errorf("약국 설정 실패: %w", err)
	}

	s.log.Info("약국 및 수령 방식 설정",
		zap.String("prescription_id", prescriptionID),
		zap.String("pharmacy_id", pharmacyID),
		zap.String("fulfillment_type", string(fulfillment)),
	)

	return nil
}

// SendPrescriptionToPharmacy는 처방전을 약국에 전송하고 조제 토큰을 발급합니다.
func (s *PrescriptionService) SendPrescriptionToPharmacy(ctx context.Context, prescriptionID string) (*FulfillmentToken, error) {
	if prescriptionID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "prescription_id는 필수입니다")
	}

	prescription, err := s.prescriptionRepo.FindByID(ctx, prescriptionID)
	if err != nil {
		return nil, err
	}

	if prescription.PharmacyID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "약국이 설정되지 않은 처방전입니다")
	}

	now := time.Now()
	tokenStr := generateFulfillmentToken()

	token := &FulfillmentToken{
		Token:          tokenStr,
		PrescriptionID: prescriptionID,
		PharmacyID:     prescription.PharmacyID,
		CreatedAt:      now,
		ExpiresAt:      now.Add(24 * time.Hour),
		IsUsed:         false,
	}

	if err := s.tokenRepo.Create(ctx, token); err != nil {
		return nil, fmt.Errorf("토큰 생성 실패: %w", err)
	}

	prescription.FulfillmentToken = tokenStr
	prescription.SentToPharmacyAt = now
	prescription.DispensaryStatus = DispensaryPending

	if err := s.prescriptionRepo.Update(ctx, prescription); err != nil {
		return nil, fmt.Errorf("처방전 상태 업데이트 실패: %w", err)
	}

	s.log.Info("처방전 약국 전송 완료",
		zap.String("prescription_id", prescriptionID),
		zap.String("pharmacy_id", prescription.PharmacyID),
		zap.String("token", tokenStr),
	)

	// Publish prescription.sent_to_pharmacy event
	if s.eventPub != nil {
		_ = s.eventPub.Publish(ctx, map[string]interface{}{
			"type":            "prescription.sent_to_pharmacy",
			"prescription_id": prescriptionID,
			"pharmacy_id":     prescription.PharmacyID,
			"pharmacy_name":   prescription.PharmacyName,
			"token":           tokenStr,
		})
	}

	return token, nil
}

// GetPrescriptionByToken은 조제 토큰으로 처방전을 조회합니다.
func (s *PrescriptionService) GetPrescriptionByToken(ctx context.Context, token string) (*Prescription, error) {
	if token == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "token은 필수입니다")
	}

	ft, err := s.tokenRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if time.Now().After(ft.ExpiresAt) {
		return nil, apperrors.New(apperrors.ErrTokenExpired, "만료된 토큰입니다")
	}

	prescription, err := s.prescriptionRepo.FindByID(ctx, ft.PrescriptionID)
	if err != nil {
		return nil, err
	}

	return prescription, nil
}

// UpdateDispensaryStatus는 약국 조제 상태를 변경합니다.
func (s *PrescriptionService) UpdateDispensaryStatus(ctx context.Context, prescriptionID string, status DispensaryStatus) error {
	if prescriptionID == "" {
		return apperrors.New(apperrors.ErrInvalidInput, "prescription_id는 필수입니다")
	}

	prescription, err := s.prescriptionRepo.FindByID(ctx, prescriptionID)
	if err != nil {
		return err
	}

	if !isValidDispensaryTransition(prescription.DispensaryStatus, status) {
		return apperrors.New(apperrors.ErrInvalidInput, fmt.Sprintf("잘못된 상태 전환입니다: %s → %s", prescription.DispensaryStatus, status))
	}

	prescription.DispensaryStatus = status

	if status == DispensaryDispensed {
		prescription.Status = StatusDispensed
		prescription.DispensedAt = time.Now()
	}

	if err := s.prescriptionRepo.Update(ctx, prescription); err != nil {
		return fmt.Errorf("조제 상태 업데이트 실패: %w", err)
	}

	s.log.Info("조제 상태 변경",
		zap.String("prescription_id", prescriptionID),
		zap.String("dispensary_status", string(status)),
	)

	// Publish prescription.dispensed event when dispensed
	if status == DispensaryDispensed && s.eventPub != nil {
		_ = s.eventPub.Publish(ctx, map[string]interface{}{
			"type":            "prescription.dispensed",
			"prescription_id": prescriptionID,
			"user_id":         prescription.PatientUserID,
			"pharmacy_id":     prescription.PharmacyID,
		})
	}

	return nil
}

// isValidDispensaryTransition은 조제 상태 전환이 유효한지 검사합니다.
func isValidDispensaryTransition(from, to DispensaryStatus) bool {
	switch from {
	case DispensaryPending:
		return to == DispensaryPreparing
	case DispensaryPreparing:
		return to == DispensaryReady
	case DispensaryReady:
		return to == DispensaryDispensed
	default:
		return false
	}
}

// generateFulfillmentToken은 6자리 영숫자 조제 토큰을 생성합니다.
func generateFulfillmentToken() string {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // Exclude ambiguous chars
	b := make([]byte, 6)
	rand.Read(b)
	for i := range b {
		b[i] = chars[int(b[i])%len(chars)]
	}
	return string(b)
}

// parseFrequencyToTimes는 복용 빈도를 시간대로 변환합니다.
func parseFrequencyToTimes(frequency string) []string {
	switch frequency {
	case "1일 1회", "QD":
		return []string{"08:00"}
	case "1일 2회", "BID":
		return []string{"08:00", "20:00"}
	case "1일 3회", "TID":
		return []string{"08:00", "13:00", "20:00"}
	case "1일 4회", "QID":
		return []string{"08:00", "12:00", "18:00", "22:00"}
	case "취침전", "HS":
		return []string{"22:00"}
	default:
		return []string{"08:00"}
	}
}
