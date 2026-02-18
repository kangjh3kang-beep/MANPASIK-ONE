// Package handler는 prescription-service의 gRPC 핸들러입니다.
package handler

import (
	"context"
	"time"

	"github.com/manpasik/backend/services/prescription-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// PrescriptionHandler는 PrescriptionService gRPC 서버를 구현합니다.
type PrescriptionHandler struct {
	v1.UnimplementedPrescriptionServiceServer
	svc *service.PrescriptionService
	log *zap.Logger
}

// NewPrescriptionHandler는 PrescriptionHandler를 생성합니다.
func NewPrescriptionHandler(svc *service.PrescriptionService, log *zap.Logger) *PrescriptionHandler {
	return &PrescriptionHandler{svc: svc, log: log}
}

// CreatePrescription은 처방전 생성 RPC입니다.
func (h *PrescriptionHandler) CreatePrescription(ctx context.Context, req *v1.CreatePrescriptionRequest) (*v1.Prescription, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}
	if req.DoctorId == "" {
		return nil, status.Error(codes.InvalidArgument, "doctor_id는 필수입니다")
	}

	var meds []*service.Medication
	for _, m := range req.Medications {
		meds = append(meds, protoMedicationToService(m))
	}

	prescription, err := h.svc.CreatePrescription(
		ctx,
		req.UserId,
		req.DoctorId,
		"",
		req.Diagnosis,
		req.Notes,
		meds,
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	return prescriptionToProto(prescription), nil
}

// GetPrescription은 처방전 조회 RPC입니다.
func (h *PrescriptionHandler) GetPrescription(ctx context.Context, req *v1.GetPrescriptionRequest) (*v1.Prescription, error) {
	if req == nil || req.PrescriptionId == "" {
		return nil, status.Error(codes.InvalidArgument, "prescription_id는 필수입니다")
	}

	prescription, err := h.svc.GetPrescription(ctx, req.PrescriptionId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return prescriptionToProto(prescription), nil
}

// ListPrescriptions는 처방전 목록 조회 RPC입니다.
func (h *PrescriptionHandler) ListPrescriptions(ctx context.Context, req *v1.ListPrescriptionsRequest) (*v1.ListPrescriptionsResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	prescriptions, total, err := h.svc.ListPrescriptions(
		ctx,
		req.UserId,
		protoPrescriptionStatusToService(req.StatusFilter),
		int(req.Limit),
		int(req.Offset),
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	var pb []*v1.Prescription
	for _, p := range prescriptions {
		pb = append(pb, prescriptionToProto(p))
	}

	return &v1.ListPrescriptionsResponse{
		Prescriptions: pb,
		TotalCount:    int32(total),
	}, nil
}

// UpdatePrescriptionStatus는 처방전 상태 변경 RPC입니다.
func (h *PrescriptionHandler) UpdatePrescriptionStatus(ctx context.Context, req *v1.UpdatePrescriptionStatusRequest) (*v1.Prescription, error) {
	if req == nil || req.PrescriptionId == "" {
		return nil, status.Error(codes.InvalidArgument, "prescription_id는 필수입니다")
	}

	prescription, err := h.svc.UpdatePrescriptionStatus(
		ctx,
		req.PrescriptionId,
		protoPrescriptionStatusToService(req.NewStatus),
		"",
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	return prescriptionToProto(prescription), nil
}

// AddMedication은 약물 추가 RPC입니다.
func (h *PrescriptionHandler) AddMedication(ctx context.Context, req *v1.AddMedicationRequest) (*v1.Prescription, error) {
	if req == nil || req.PrescriptionId == "" {
		return nil, status.Error(codes.InvalidArgument, "prescription_id는 필수입니다")
	}
	if req.Medication == nil {
		return nil, status.Error(codes.InvalidArgument, "medication은 필수입니다")
	}

	prescription, err := h.svc.AddMedication(
		ctx,
		req.PrescriptionId,
		protoMedicationToService(req.Medication),
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	return prescriptionToProto(prescription), nil
}

// RemoveMedication은 약물 제거 RPC입니다.
func (h *PrescriptionHandler) RemoveMedication(ctx context.Context, req *v1.RemoveMedicationRequest) (*v1.Prescription, error) {
	if req == nil || req.PrescriptionId == "" {
		return nil, status.Error(codes.InvalidArgument, "prescription_id는 필수입니다")
	}

	prescription, err := h.svc.RemoveMedication(ctx, req.PrescriptionId, req.MedicationId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return prescriptionToProto(prescription), nil
}

// CheckDrugInteraction은 약물 상호작용 검사 RPC입니다.
func (h *PrescriptionHandler) CheckDrugInteraction(ctx context.Context, req *v1.CheckDrugInteractionRequest) (*v1.CheckDrugInteractionResponse, error) {
	if req == nil || len(req.MedicationNames) < 2 {
		return &v1.CheckDrugInteractionResponse{}, nil
	}

	interactions, err := h.svc.CheckDrugInteraction(ctx, req.MedicationNames)
	if err != nil {
		return nil, toGRPC(err)
	}

	var pbInteractions []*v1.DrugInteraction
	hasCritical := false
	for _, i := range interactions {
		pbInteractions = append(pbInteractions, &v1.DrugInteraction{
			DrugA:          i.DrugA,
			DrugB:          i.DrugB,
			Severity:       serviceSeverityToProto(i.Severity),
			Description:    i.Description,
			Recommendation: i.Recommendation,
		})
		if i.Severity >= service.SeverityMajor {
			hasCritical = true
		}
	}

	return &v1.CheckDrugInteractionResponse{
		Interactions: pbInteractions,
		HasCritical:  hasCritical,
	}, nil
}

// GetMedicationReminders는 복약 알림 조회 RPC입니다.
func (h *PrescriptionHandler) GetMedicationReminders(ctx context.Context, req *v1.GetMedicationRemindersRequest) (*v1.GetMedicationRemindersResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	reminders, err := h.svc.GetMedicationReminders(ctx, req.UserId)
	if err != nil {
		return nil, toGRPC(err)
	}

	var pb []*v1.MedicationReminder
	for _, r := range reminders {
		pb = append(pb, &v1.MedicationReminder{
			PrescriptionId: r.PrescriptionID,
			MedicationName: r.DrugName,
			Dosage:         r.Dosage,
			ScheduledTime:  r.TimeOfDay,
			Instructions:   r.Instructions,
			IsTaken:        r.IsTaken,
		})
	}

	return &v1.GetMedicationRemindersResponse{
		Reminders: pb,
	}, nil
}

// SelectPharmacyAndFulfillment은 약국 선택 및 수령 방식 설정 RPC입니다.
func (h *PrescriptionHandler) SelectPharmacyAndFulfillment(ctx context.Context, req *v1.SelectPharmacyRequest) (*v1.SelectPharmacyResponse, error) {
	ft := service.FulfillmentType(req.FulfillmentType)
	err := h.svc.SelectPharmacyAndFulfillment(ctx, req.PrescriptionId, req.PharmacyId, req.PharmacyName, ft, req.ShippingAddress)
	if err != nil {
		return nil, toGRPC(err)
	}
	return &v1.SelectPharmacyResponse{Success: true, Message: "약국이 선택되었습니다"}, nil
}

// SendPrescriptionToPharmacy는 처방전을 약국에 전송하는 RPC입니다.
func (h *PrescriptionHandler) SendPrescriptionToPharmacy(ctx context.Context, req *v1.SendToPharmacyRequest) (*v1.SendToPharmacyResponse, error) {
	token, err := h.svc.SendPrescriptionToPharmacy(ctx, req.PrescriptionId)
	if err != nil {
		return nil, toGRPC(err)
	}
	return &v1.SendToPharmacyResponse{
		FulfillmentToken: token.Token,
		ExpiresAt:        token.ExpiresAt.Format(time.RFC3339),
		Success:          true,
	}, nil
}

// GetPrescriptionByToken은 조제 토큰으로 처방전을 조회하는 RPC입니다.
func (h *PrescriptionHandler) GetPrescriptionByToken(ctx context.Context, req *v1.GetByTokenRequest) (*v1.Prescription, error) {
	p, err := h.svc.GetPrescriptionByToken(ctx, req.FulfillmentToken)
	if err != nil {
		return nil, toGRPC(err)
	}
	return prescriptionToProto(p), nil
}

// UpdateDispensaryStatus는 조제 상태를 업데이트하는 RPC입니다.
func (h *PrescriptionHandler) UpdateDispensaryStatus(ctx context.Context, req *v1.UpdateDispensaryStatusRequest) (*v1.Prescription, error) {
	err := h.svc.UpdateDispensaryStatus(ctx, req.PrescriptionId, service.DispensaryStatus(req.Status))
	if err != nil {
		return nil, toGRPC(err)
	}
	p, err := h.svc.GetPrescription(ctx, req.PrescriptionId)
	if err != nil {
		return nil, toGRPC(err)
	}
	return prescriptionToProto(p), nil
}

// ============================================================================
// 변환 헬퍼
// ============================================================================

func prescriptionToProto(p *service.Prescription) *v1.Prescription {
	var meds []*v1.Medication
	for _, m := range p.Medications {
		meds = append(meds, serviceMedicationToProto(m))
	}

	pb := &v1.Prescription{
		PrescriptionId: p.ID,
		UserId:         p.PatientUserID,
		DoctorId:       p.DoctorID,
		Status:         servicePrescriptionStatusToProto(p.Status),
		Medications:    meds,
		Diagnosis:      p.Diagnosis,
		Notes:          p.Notes,
	}
	if !p.PrescribedAt.IsZero() {
		pb.PrescribedAt = timestamppb.New(p.PrescribedAt)
	}
	if !p.ExpiresAt.IsZero() {
		pb.ExpiresAt = timestamppb.New(p.ExpiresAt)
	}
	if !p.CreatedAt.IsZero() {
		pb.UpdatedAt = timestamppb.New(p.CreatedAt)
	}
	return pb
}

func serviceMedicationToProto(m *service.Medication) *v1.Medication {
	return &v1.Medication{
		MedicationId: m.ID,
		Name:         m.DrugName,
		Dosage:       m.Dosage,
		Frequency:    m.Frequency,
		DurationDays: int32(m.DurationDays),
		Route:        m.Route,
		Instructions: m.Instructions,
	}
}

func protoMedicationToService(m *v1.Medication) *service.Medication {
	return &service.Medication{
		ID:           m.MedicationId,
		DrugName:     m.Name,
		Dosage:       m.Dosage,
		Frequency:    m.Frequency,
		DurationDays: int(m.DurationDays),
		Route:        m.Route,
		Instructions: m.Instructions,
	}
}

func protoPrescriptionStatusToService(s v1.PrescriptionStatus) service.PrescriptionStatus {
	switch s {
	case v1.PrescriptionStatus_PRESCRIPTION_STATUS_DRAFT:
		return service.StatusDraft
	case v1.PrescriptionStatus_PRESCRIPTION_STATUS_ACTIVE:
		return service.StatusActive
	case v1.PrescriptionStatus_PRESCRIPTION_STATUS_DISPENSED:
		return service.StatusDispensed
	case v1.PrescriptionStatus_PRESCRIPTION_STATUS_COMPLETED:
		return service.StatusCompleted
	case v1.PrescriptionStatus_PRESCRIPTION_STATUS_CANCELLED:
		return service.StatusCancelled
	case v1.PrescriptionStatus_PRESCRIPTION_STATUS_EXPIRED:
		return service.StatusExpired
	default:
		return service.StatusUnknown
	}
}

func servicePrescriptionStatusToProto(s service.PrescriptionStatus) v1.PrescriptionStatus {
	switch s {
	case service.StatusDraft:
		return v1.PrescriptionStatus_PRESCRIPTION_STATUS_DRAFT
	case service.StatusActive:
		return v1.PrescriptionStatus_PRESCRIPTION_STATUS_ACTIVE
	case service.StatusDispensed:
		return v1.PrescriptionStatus_PRESCRIPTION_STATUS_DISPENSED
	case service.StatusCompleted:
		return v1.PrescriptionStatus_PRESCRIPTION_STATUS_COMPLETED
	case service.StatusCancelled:
		return v1.PrescriptionStatus_PRESCRIPTION_STATUS_CANCELLED
	case service.StatusExpired:
		return v1.PrescriptionStatus_PRESCRIPTION_STATUS_EXPIRED
	default:
		return v1.PrescriptionStatus_PRESCRIPTION_STATUS_UNKNOWN
	}
}

func serviceSeverityToProto(s service.InteractionSeverity) v1.DrugInteractionSeverity {
	switch s {
	case service.SeverityNone:
		return v1.DrugInteractionSeverity_DRUG_INTERACTION_SEVERITY_NONE
	case service.SeverityMinor:
		return v1.DrugInteractionSeverity_DRUG_INTERACTION_SEVERITY_MINOR
	case service.SeverityModerate:
		return v1.DrugInteractionSeverity_DRUG_INTERACTION_SEVERITY_MODERATE
	case service.SeverityMajor:
		return v1.DrugInteractionSeverity_DRUG_INTERACTION_SEVERITY_MAJOR
	case service.SeverityContraindicated:
		return v1.DrugInteractionSeverity_DRUG_INTERACTION_SEVERITY_CONTRAINDICATED
	default:
		return v1.DrugInteractionSeverity_DRUG_INTERACTION_SEVERITY_UNKNOWN
	}
}

func toGRPC(err error) error {
	if appErr, ok := err.(*apperrors.AppError); ok {
		return appErr.ToGRPC()
	}
	return status.Error(codes.Internal, err.Error())
}
