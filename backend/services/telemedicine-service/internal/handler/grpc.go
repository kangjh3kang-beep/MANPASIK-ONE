// Package handler는 telemedicine-service의 gRPC 핸들러입니다.
package handler

import (
	"context"

	"github.com/manpasik/backend/services/telemedicine-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TelemedicineHandler는 TelemedicineService gRPC 서버를 구현합니다.
type TelemedicineHandler struct {
	v1.UnimplementedTelemedicineServiceServer
	svc *service.TelemedicineService
	log *zap.Logger
}

// NewTelemedicineHandler는 TelemedicineHandler를 생성합니다.
func NewTelemedicineHandler(svc *service.TelemedicineService, log *zap.Logger) *TelemedicineHandler {
	return &TelemedicineHandler{svc: svc, log: log}
}

// CreateConsultation은 상담 생성 RPC입니다.
func (h *TelemedicineHandler) CreateConsultation(ctx context.Context, req *v1.CreateConsultationRequest) (*v1.Consultation, error) {
	if req == nil || req.PatientUserId == "" {
		return nil, status.Error(codes.InvalidArgument, "patient_user_id는 필수입니다")
	}
	if req.ChiefComplaint == "" {
		return nil, status.Error(codes.InvalidArgument, "chief_complaint는 필수입니다")
	}

	consultation, err := h.svc.CreateConsultation(
		ctx,
		req.PatientUserId,
		protoSpecialtyToService(req.Specialty),
		req.ChiefComplaint,
		req.Description,
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	return consultationToProto(consultation), nil
}

// GetConsultation은 상담 조회 RPC입니다.
func (h *TelemedicineHandler) GetConsultation(ctx context.Context, req *v1.GetConsultationRequest) (*v1.Consultation, error) {
	if req == nil || req.ConsultationId == "" {
		return nil, status.Error(codes.InvalidArgument, "consultation_id는 필수입니다")
	}

	consultation, err := h.svc.GetConsultation(ctx, req.ConsultationId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return consultationToProto(consultation), nil
}

// ListConsultations는 상담 목록 조회 RPC입니다.
func (h *TelemedicineHandler) ListConsultations(ctx context.Context, req *v1.ListConsultationsRequest) (*v1.ListConsultationsResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	consultations, total, err := h.svc.ListConsultations(
		ctx,
		req.UserId,
		protoConsultationStatusToService(req.StatusFilter),
		int(req.Limit),
		int(req.Offset),
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	var pbConsultations []*v1.Consultation
	for _, c := range consultations {
		pbConsultations = append(pbConsultations, consultationToProto(c))
	}

	return &v1.ListConsultationsResponse{
		Consultations: pbConsultations,
		TotalCount:    int32(total),
	}, nil
}

// MatchDoctor는 의사 매칭 RPC입니다.
func (h *TelemedicineHandler) MatchDoctor(ctx context.Context, req *v1.MatchDoctorRequest) (*v1.MatchDoctorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "요청이 비어 있습니다")
	}

	doctors, err := h.svc.MatchDoctor(
		ctx,
		protoSpecialtyToService(req.Specialty),
		req.Language,
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	var pbDoctors []*v1.DoctorProfile
	for _, d := range doctors {
		pbDoctors = append(pbDoctors, doctorProfileToProto(d))
	}

	return &v1.MatchDoctorResponse{
		Doctors:        pbDoctors,
		TotalAvailable: int32(len(pbDoctors)),
	}, nil
}

// StartVideoSession은 비디오 세션 시작 RPC입니다.
func (h *TelemedicineHandler) StartVideoSession(ctx context.Context, req *v1.StartVideoSessionRequest) (*v1.VideoSession, error) {
	if req == nil || req.ConsultationId == "" {
		return nil, status.Error(codes.InvalidArgument, "consultation_id는 필수입니다")
	}
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	session, err := h.svc.StartVideoSession(ctx, req.ConsultationId, req.UserId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return videoSessionToProto(session), nil
}

// EndVideoSession은 비디오 세션 종료 RPC입니다.
func (h *TelemedicineHandler) EndVideoSession(ctx context.Context, req *v1.EndVideoSessionRequest) (*v1.VideoSession, error) {
	if req == nil || req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id는 필수입니다")
	}

	session, err := h.svc.EndVideoSession(ctx, req.SessionId, req.ConsultationId, req.DoctorNotes, req.Diagnosis)
	if err != nil {
		return nil, toGRPC(err)
	}

	return videoSessionToProto(session), nil
}

// RateConsultation은 상담 평점 등록 RPC입니다.
func (h *TelemedicineHandler) RateConsultation(ctx context.Context, req *v1.RateConsultationRequest) (*v1.RateConsultationResponse, error) {
	if req == nil || req.ConsultationId == "" {
		return nil, status.Error(codes.InvalidArgument, "consultation_id는 필수입니다")
	}

	newRating, err := h.svc.RateConsultation(ctx, req.ConsultationId, req.Rating)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.RateConsultationResponse{
		Success:          true,
		NewAverageRating: newRating,
	}, nil
}

// ============================================================================
// 변환 헬퍼
// ============================================================================

func consultationToProto(c *service.Consultation) *v1.Consultation {
	pb := &v1.Consultation{
		ConsultationId:  c.ID,
		PatientUserId:   c.PatientUserID,
		DoctorId:        c.DoctorID,
		Specialty:       serviceSpecialtyToProto(c.Specialty),
		ChiefComplaint:  c.ChiefComplaint,
		Description:     c.Description,
		Status:          serviceConsultationStatusToProto(c.Status),
		Diagnosis:       c.Diagnosis,
		DoctorNotes:     c.DoctorNotes,
		PrescriptionId:  c.PrescriptionID,
		DurationMinutes: int32(c.DurationMinutes),
		Rating:          c.Rating,
		CreatedAt:       timestamppb.New(c.CreatedAt),
	}
	if !c.ScheduledAt.IsZero() {
		pb.ScheduledAt = timestamppb.New(c.ScheduledAt)
	}
	if !c.StartedAt.IsZero() {
		pb.StartedAt = timestamppb.New(c.StartedAt)
	}
	if !c.EndedAt.IsZero() {
		pb.EndedAt = timestamppb.New(c.EndedAt)
	}
	return pb
}

func doctorProfileToProto(d *service.DoctorProfile) *v1.DoctorProfile {
	return &v1.DoctorProfile{
		DoctorId:           d.ID,
		Name:               d.Name,
		Specialty:          serviceSpecialtyToProto(d.Specialty),
		Hospital:           d.Hospital,
		LicenseNumber:      d.LicenseNumber,
		ExperienceYears:    int32(d.ExperienceYears),
		Rating:             d.Rating,
		TotalConsultations: int32(d.TotalConsultations),
		IsAvailable:        d.IsAvailable,
		Languages:          d.Languages,
		ProfileImageUrl:    d.ProfileImageURL,
	}
}

func videoSessionToProto(s *service.VideoSession) *v1.VideoSession {
	pb := &v1.VideoSession{
		SessionId:       s.ID,
		ConsultationId:  s.ConsultationID,
		RoomUrl:         s.RoomURL,
		Token:           s.Token,
		Status:          serviceVideoSessionStatusToProto(s.Status),
		DurationSeconds: int32(s.DurationSeconds),
	}
	if !s.StartedAt.IsZero() {
		pb.StartedAt = timestamppb.New(s.StartedAt)
	}
	if !s.EndedAt.IsZero() {
		pb.EndedAt = timestamppb.New(s.EndedAt)
	}
	return pb
}

// --- ConsultationStatus 변환 ---

func protoConsultationStatusToService(s v1.ConsultationStatus) service.ConsultationStatus {
	switch s {
	case v1.ConsultationStatus_CONSULTATION_STATUS_REQUESTED:
		return service.StatusRequested
	case v1.ConsultationStatus_CONSULTATION_STATUS_MATCHED:
		return service.StatusMatched
	case v1.ConsultationStatus_CONSULTATION_STATUS_SCHEDULED:
		return service.StatusScheduled
	case v1.ConsultationStatus_CONSULTATION_STATUS_IN_PROGRESS:
		return service.StatusInProgress
	case v1.ConsultationStatus_CONSULTATION_STATUS_COMPLETED:
		return service.StatusCompleted
	case v1.ConsultationStatus_CONSULTATION_STATUS_CANCELLED:
		return service.StatusCancelled
	case v1.ConsultationStatus_CONSULTATION_STATUS_NO_SHOW:
		return service.StatusNoShow
	default:
		return service.StatusUnknown
	}
}

func serviceConsultationStatusToProto(s service.ConsultationStatus) v1.ConsultationStatus {
	switch s {
	case service.StatusRequested:
		return v1.ConsultationStatus_CONSULTATION_STATUS_REQUESTED
	case service.StatusMatched:
		return v1.ConsultationStatus_CONSULTATION_STATUS_MATCHED
	case service.StatusScheduled:
		return v1.ConsultationStatus_CONSULTATION_STATUS_SCHEDULED
	case service.StatusInProgress:
		return v1.ConsultationStatus_CONSULTATION_STATUS_IN_PROGRESS
	case service.StatusCompleted:
		return v1.ConsultationStatus_CONSULTATION_STATUS_COMPLETED
	case service.StatusCancelled:
		return v1.ConsultationStatus_CONSULTATION_STATUS_CANCELLED
	case service.StatusNoShow:
		return v1.ConsultationStatus_CONSULTATION_STATUS_NO_SHOW
	default:
		return v1.ConsultationStatus_CONSULTATION_STATUS_UNKNOWN
	}
}

// --- DoctorSpecialty 변환 ---

func protoSpecialtyToService(s v1.DoctorSpecialty) service.DoctorSpecialty {
	switch s {
	case v1.DoctorSpecialty_DOCTOR_SPECIALTY_GENERAL:
		return service.SpecialtyGeneral
	case v1.DoctorSpecialty_DOCTOR_SPECIALTY_INTERNAL:
		return service.SpecialtyInternal
	case v1.DoctorSpecialty_DOCTOR_SPECIALTY_CARDIOLOGY:
		return service.SpecialtyCardiology
	case v1.DoctorSpecialty_DOCTOR_SPECIALTY_ENDOCRINOLOGY:
		return service.SpecialtyEndocrinology
	case v1.DoctorSpecialty_DOCTOR_SPECIALTY_DERMATOLOGY:
		return service.SpecialtyDermatology
	case v1.DoctorSpecialty_DOCTOR_SPECIALTY_PEDIATRICS:
		return service.SpecialtyPediatrics
	case v1.DoctorSpecialty_DOCTOR_SPECIALTY_PSYCHIATRY:
		return service.SpecialtyPsychiatry
	case v1.DoctorSpecialty_DOCTOR_SPECIALTY_ORTHOPEDICS:
		return service.SpecialtyOrthopedics
	case v1.DoctorSpecialty_DOCTOR_SPECIALTY_OPHTHALMOLOGY:
		return service.SpecialtyOphthalmology
	case v1.DoctorSpecialty_DOCTOR_SPECIALTY_ENT:
		return service.SpecialtyENT
	default:
		return service.SpecialtyUnknown
	}
}

func serviceSpecialtyToProto(s service.DoctorSpecialty) v1.DoctorSpecialty {
	switch s {
	case service.SpecialtyGeneral:
		return v1.DoctorSpecialty_DOCTOR_SPECIALTY_GENERAL
	case service.SpecialtyInternal:
		return v1.DoctorSpecialty_DOCTOR_SPECIALTY_INTERNAL
	case service.SpecialtyCardiology:
		return v1.DoctorSpecialty_DOCTOR_SPECIALTY_CARDIOLOGY
	case service.SpecialtyEndocrinology:
		return v1.DoctorSpecialty_DOCTOR_SPECIALTY_ENDOCRINOLOGY
	case service.SpecialtyDermatology:
		return v1.DoctorSpecialty_DOCTOR_SPECIALTY_DERMATOLOGY
	case service.SpecialtyPediatrics:
		return v1.DoctorSpecialty_DOCTOR_SPECIALTY_PEDIATRICS
	case service.SpecialtyPsychiatry:
		return v1.DoctorSpecialty_DOCTOR_SPECIALTY_PSYCHIATRY
	case service.SpecialtyOrthopedics:
		return v1.DoctorSpecialty_DOCTOR_SPECIALTY_ORTHOPEDICS
	case service.SpecialtyOphthalmology:
		return v1.DoctorSpecialty_DOCTOR_SPECIALTY_OPHTHALMOLOGY
	case service.SpecialtyENT:
		return v1.DoctorSpecialty_DOCTOR_SPECIALTY_ENT
	default:
		return v1.DoctorSpecialty_DOCTOR_SPECIALTY_UNKNOWN
	}
}

// --- VideoSessionStatus 변환 ---

func serviceVideoSessionStatusToProto(s service.VideoSessionStatus) v1.VideoSessionStatus {
	switch s {
	case service.SessionWaiting:
		return v1.VideoSessionStatus_VIDEO_SESSION_STATUS_WAITING
	case service.SessionConnected:
		return v1.VideoSessionStatus_VIDEO_SESSION_STATUS_CONNECTED
	case service.SessionEnded:
		return v1.VideoSessionStatus_VIDEO_SESSION_STATUS_ENDED
	case service.SessionFailed:
		return v1.VideoSessionStatus_VIDEO_SESSION_STATUS_FAILED
	default:
		return v1.VideoSessionStatus_VIDEO_SESSION_STATUS_UNKNOWN
	}
}

func toGRPC(err error) error {
	if appErr, ok := err.(*apperrors.AppError); ok {
		return appErr.ToGRPC()
	}
	return status.Error(codes.Internal, err.Error())
}
