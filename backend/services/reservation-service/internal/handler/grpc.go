// Package handler는 reservation-service의 gRPC 핸들러입니다.
package handler

import (
	"context"
	"time"

	"github.com/manpasik/backend/services/reservation-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ReservationHandler는 ReservationService gRPC 서버를 구현합니다.
type ReservationHandler struct {
	v1.UnimplementedReservationServiceServer
	svc *service.ReservationService
	log *zap.Logger
}

// NewReservationHandler는 ReservationHandler를 생성합니다.
func NewReservationHandler(svc *service.ReservationService, log *zap.Logger) *ReservationHandler {
	return &ReservationHandler{svc: svc, log: log}
}

// SearchFacilities는 시설 검색 RPC입니다.
func (h *ReservationHandler) SearchFacilities(ctx context.Context, req *v1.SearchFacilitiesRequest) (*v1.SearchFacilitiesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "요청이 비어 있습니다")
	}

	facilities, total, err := h.svc.SearchFacilities(
		ctx,
		protoFacilityTypeToService(req.Type),
		req.Query,
		0,
		int(req.Limit),
		0,
		"", "", "", // countryCode, regionCode, districtCode
		0, 0, // userLat, userLon
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	var pbFacilities []*v1.Facility
	for _, f := range facilities {
		pbFacilities = append(pbFacilities, facilityToProto(f))
	}

	return &v1.SearchFacilitiesResponse{
		Facilities: pbFacilities,
		TotalCount: int32(total),
	}, nil
}

// GetFacility는 시설 조회 RPC입니다.
func (h *ReservationHandler) GetFacility(ctx context.Context, req *v1.GetFacilityRequest) (*v1.Facility, error) {
	if req == nil || req.FacilityId == "" {
		return nil, status.Error(codes.InvalidArgument, "facility_id는 필수입니다")
	}

	facility, err := h.svc.GetFacility(ctx, req.FacilityId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return facilityToProto(facility), nil
}

// ListDoctorsByFacility는 시설별 의사 목록 조회 RPC입니다.
func (h *ReservationHandler) ListDoctorsByFacility(ctx context.Context, req *v1.ListDoctorsByFacilityRequest) (*v1.ListDoctorsByFacilityResponse, error) {
	if req == nil || req.FacilityId == "" {
		return nil, status.Error(codes.InvalidArgument, "facility_id는 필수입니다")
	}

	doctors, err := h.svc.ListDoctorsByFacility(ctx, req.FacilityId, "")
	if err != nil {
		return nil, toGRPC(err)
	}

	var pbDoctors []*v1.Doctor
	for _, d := range doctors {
		pbDoctors = append(pbDoctors, doctorToProto(d))
	}

	return &v1.ListDoctorsByFacilityResponse{Doctors: pbDoctors}, nil
}

// GetDoctorAvailability는 의사별 예약 가능 시간대 조회 RPC입니다.
func (h *ReservationHandler) GetDoctorAvailability(ctx context.Context, req *v1.GetDoctorAvailabilityRequest) (*v1.GetDoctorAvailabilityResponse, error) {
	if req == nil || req.DoctorId == "" {
		return nil, status.Error(codes.InvalidArgument, "doctor_id는 필수입니다")
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "날짜 형식 오류: %v", err)
	}

	slots, err := h.svc.GetDoctorAvailability(ctx, req.DoctorId, date)
	if err != nil {
		return nil, toGRPC(err)
	}

	var pbSlots []*v1.TimeSlotDetail
	for _, s := range slots {
		pbSlots = append(pbSlots, &v1.TimeSlotDetail{
			StartTime:   s.StartTime.Format(time.RFC3339),
			EndTime:     s.EndTime.Format(time.RFC3339),
			IsAvailable: s.IsAvailable,
		})
	}

	return &v1.GetDoctorAvailabilityResponse{Slots: pbSlots}, nil
}

// SelectDoctor는 시설 내 의사 선택 RPC입니다.
func (h *ReservationHandler) SelectDoctor(ctx context.Context, req *v1.SelectDoctorRequest) (*v1.SelectDoctorResponse, error) {
	if req == nil || req.FacilityId == "" || req.DoctorId == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "facility_id, doctor_id, user_id는 필수입니다")
	}

	doctor, err := h.svc.SelectDoctor(ctx, req.FacilityId, req.DoctorId, req.UserId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.SelectDoctorResponse{
		Doctor:  doctorToProto(doctor),
		Success: true,
	}, nil
}

// GetAvailableSlots는 예약 가능 시간대 조회 RPC입니다.
func (h *ReservationHandler) GetAvailableSlots(ctx context.Context, req *v1.GetAvailableSlotsRequest) (*v1.GetAvailableSlotsResponse, error) {
	if req == nil || req.FacilityId == "" {
		return nil, status.Error(codes.InvalidArgument, "facility_id는 필수입니다")
	}

	var date time.Time
	if req.Date != "" {
		if parsed, err := time.Parse("2006-01-02", req.Date); err == nil {
			date = parsed
		}
	}

	slots, err := h.svc.GetAvailableSlots(
		ctx,
		req.FacilityId,
		date,
		req.DoctorId,
		protoSpecialtyToService(req.Specialty),
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	var pbSlots []*v1.TimeSlot
	for _, slot := range slots {
		pbSlots = append(pbSlots, timeSlotToProto(slot))
	}

	return &v1.GetAvailableSlotsResponse{
		Slots: pbSlots,
	}, nil
}

// CreateReservation은 예약 생성 RPC입니다.
func (h *ReservationHandler) CreateReservation(ctx context.Context, req *v1.CreateReservationRequest) (*v1.Reservation, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}
	if req.FacilityId == "" {
		return nil, status.Error(codes.InvalidArgument, "facility_id는 필수입니다")
	}

	reservation, err := h.svc.CreateReservation(
		ctx,
		req.UserId,
		req.FacilityId,
		"",
		req.SlotId,
		0,
		req.Reason,
		req.Notes,
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	return reservationToProto(reservation), nil
}

// GetReservation은 예약 조회 RPC입니다.
func (h *ReservationHandler) GetReservation(ctx context.Context, req *v1.GetReservationRequest) (*v1.Reservation, error) {
	if req == nil || req.ReservationId == "" {
		return nil, status.Error(codes.InvalidArgument, "reservation_id는 필수입니다")
	}

	reservation, err := h.svc.GetReservation(ctx, req.ReservationId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return reservationToProto(reservation), nil
}

// ListReservations는 예약 목록 조회 RPC입니다.
func (h *ReservationHandler) ListReservations(ctx context.Context, req *v1.ListReservationsRequest) (*v1.ListReservationsResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	reservations, total, err := h.svc.ListReservations(
		ctx,
		req.UserId,
		protoReservationStatusToService(req.Status),
		int(req.Limit),
		int(req.Offset),
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	var pbReservations []*v1.Reservation
	for _, r := range reservations {
		pbReservations = append(pbReservations, reservationToProto(r))
	}

	return &v1.ListReservationsResponse{
		Reservations: pbReservations,
		TotalCount:   int32(total),
	}, nil
}

// CancelReservation은 예약 취소 RPC입니다.
func (h *ReservationHandler) CancelReservation(ctx context.Context, req *v1.CancelReservationRequest) (*v1.CancelReservationResponse, error) {
	if req == nil || req.ReservationId == "" {
		return nil, status.Error(codes.InvalidArgument, "reservation_id는 필수입니다")
	}

	err := h.svc.CancelReservation(ctx, req.ReservationId, "", req.Reason)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.CancelReservationResponse{Success: true, Message: "예약이 취소되었습니다"}, nil
}

// ============================================================================
// 변환 헬퍼
// ============================================================================

func facilityToProto(f *service.Facility) *v1.Facility {
	var specialties []v1.DoctorSpecialty
	for _, sp := range f.Specialties {
		specialties = append(specialties, serviceSpecialtyToProto(sp))
	}

	return &v1.Facility{
		FacilityId:         f.ID,
		Name:               f.Name,
		Type:               serviceFacilityTypeToProto(f.Type),
		Address:            f.Address,
		Phone:              f.Phone,
		Latitude:           f.Latitude,
		Longitude:          f.Longitude,
		Rating:             float32(f.Rating),
		Specialties:        specialties,
		OperatingHours:     f.OperatingHours,
		IsOpenNow:          f.IsOpenNow,
		AcceptsReservation: f.AcceptsReservation,
	}
}

func timeSlotToProto(slot *service.TimeSlot) *v1.TimeSlot {
	pb := &v1.TimeSlot{
		SlotId:      slot.ID,
		IsAvailable: slot.IsAvailable,
		DoctorName:  slot.DoctorName,
	}
	if !slot.StartTime.IsZero() {
		pb.StartTime = timestamppb.New(slot.StartTime)
	}
	if !slot.EndTime.IsZero() {
		pb.EndTime = timestamppb.New(slot.EndTime)
	}
	return pb
}

func reservationToProto(r *service.Reservation) *v1.Reservation {
	pb := &v1.Reservation{
		ReservationId: r.ID,
		UserId:        r.UserID,
		FacilityId:    r.FacilityID,
		FacilityName:  r.FacilityName,
		DoctorName:    r.DoctorName,
		Status:        serviceReservationStatusToProto(r.Status),
		Reason:        r.Reason,
		Notes:         r.Notes,
		CreatedAt:     timestamppb.New(r.CreatedAt),
	}
	if !r.ScheduledAt.IsZero() {
		pb.AppointmentTime = timestamppb.New(r.ScheduledAt)
	}
	return pb
}

func doctorToProto(d *service.Doctor) *v1.Doctor {
	return &v1.Doctor{
		DoctorId:             d.ID,
		Name:                 d.Name,
		Specialty:            d.Specialty,
		FacilityId:           d.FacilityID,
		Rating:               float32(d.Rating),
		ConsultationFee:      d.ConsultationFee,
		AcceptsTelemedicine:  d.AcceptsTelemedicine,
		AvailableRegionCodes: d.AvailableRegionCodes,
		NextAvailableAt:      d.NextAvailableAt.Format(time.RFC3339),
	}
}

// --- FacilityType 변환 ---

func protoFacilityTypeToService(t v1.FacilityType) service.FacilityType {
	switch t {
	case v1.FacilityType_FACILITY_TYPE_HOSPITAL:
		return service.FacilityHospital
	case v1.FacilityType_FACILITY_TYPE_CLINIC:
		return service.FacilityClinic
	case v1.FacilityType_FACILITY_TYPE_PHARMACY:
		return service.FacilityPharmacy
	case v1.FacilityType_FACILITY_TYPE_DENTAL:
		return service.FacilityDental
	case v1.FacilityType_FACILITY_TYPE_ORIENTAL:
		return service.FacilityOriental
	default:
		return service.FacilityUnknown
	}
}

func serviceFacilityTypeToProto(t service.FacilityType) v1.FacilityType {
	switch t {
	case service.FacilityHospital:
		return v1.FacilityType_FACILITY_TYPE_HOSPITAL
	case service.FacilityClinic:
		return v1.FacilityType_FACILITY_TYPE_CLINIC
	case service.FacilityPharmacy:
		return v1.FacilityType_FACILITY_TYPE_PHARMACY
	case service.FacilityDental:
		return v1.FacilityType_FACILITY_TYPE_DENTAL
	case service.FacilityOriental:
		return v1.FacilityType_FACILITY_TYPE_ORIENTAL
	default:
		return v1.FacilityType_FACILITY_TYPE_UNKNOWN
	}
}

// --- ReservationStatus 변환 ---

func protoReservationStatusToService(s v1.ReservationStatus) service.ReservationStatus {
	switch s {
	case v1.ReservationStatus_RESERVATION_STATUS_PENDING:
		return service.ResPending
	case v1.ReservationStatus_RESERVATION_STATUS_CONFIRMED:
		return service.ResConfirmed
	case v1.ReservationStatus_RESERVATION_STATUS_COMPLETED:
		return service.ResCompleted
	case v1.ReservationStatus_RESERVATION_STATUS_CANCELLED:
		return service.ResCancelled
	case v1.ReservationStatus_RESERVATION_STATUS_NO_SHOW:
		return service.ResNoShow
	default:
		return service.ResUnknown
	}
}

func serviceReservationStatusToProto(s service.ReservationStatus) v1.ReservationStatus {
	switch s {
	case service.ResPending:
		return v1.ReservationStatus_RESERVATION_STATUS_PENDING
	case service.ResConfirmed:
		return v1.ReservationStatus_RESERVATION_STATUS_CONFIRMED
	case service.ResCompleted:
		return v1.ReservationStatus_RESERVATION_STATUS_COMPLETED
	case service.ResCancelled:
		return v1.ReservationStatus_RESERVATION_STATUS_CANCELLED
	case service.ResNoShow:
		return v1.ReservationStatus_RESERVATION_STATUS_NO_SHOW
	default:
		return v1.ReservationStatus_RESERVATION_STATUS_UNKNOWN
	}
}

// --- DoctorSpecialty 변환 ---

func protoSpecialtyToService(s v1.DoctorSpecialty) service.Specialty {
	switch s {
	case v1.DoctorSpecialty_DOCTOR_SPECIALTY_GENERAL:
		return service.SpecGeneral
	case v1.DoctorSpecialty_DOCTOR_SPECIALTY_INTERNAL:
		return service.SpecInternal
	case v1.DoctorSpecialty_DOCTOR_SPECIALTY_CARDIOLOGY:
		return service.SpecCardiology
	case v1.DoctorSpecialty_DOCTOR_SPECIALTY_ENDOCRINOLOGY:
		return service.SpecEndocrinology
	case v1.DoctorSpecialty_DOCTOR_SPECIALTY_DERMATOLOGY:
		return service.SpecDermatology
	case v1.DoctorSpecialty_DOCTOR_SPECIALTY_PEDIATRICS:
		return service.SpecPediatrics
	case v1.DoctorSpecialty_DOCTOR_SPECIALTY_PSYCHIATRY:
		return service.SpecPsychiatry
	case v1.DoctorSpecialty_DOCTOR_SPECIALTY_ORTHOPEDICS:
		return service.SpecOrthopedics
	case v1.DoctorSpecialty_DOCTOR_SPECIALTY_OPHTHALMOLOGY:
		return service.SpecOphthalmology
	case v1.DoctorSpecialty_DOCTOR_SPECIALTY_ENT:
		return service.SpecENT
	default:
		return service.SpecUnknown
	}
}

func serviceSpecialtyToProto(s service.Specialty) v1.DoctorSpecialty {
	switch s {
	case service.SpecGeneral:
		return v1.DoctorSpecialty_DOCTOR_SPECIALTY_GENERAL
	case service.SpecInternal:
		return v1.DoctorSpecialty_DOCTOR_SPECIALTY_INTERNAL
	case service.SpecCardiology:
		return v1.DoctorSpecialty_DOCTOR_SPECIALTY_CARDIOLOGY
	case service.SpecEndocrinology:
		return v1.DoctorSpecialty_DOCTOR_SPECIALTY_ENDOCRINOLOGY
	case service.SpecDermatology:
		return v1.DoctorSpecialty_DOCTOR_SPECIALTY_DERMATOLOGY
	case service.SpecPediatrics:
		return v1.DoctorSpecialty_DOCTOR_SPECIALTY_PEDIATRICS
	case service.SpecPsychiatry:
		return v1.DoctorSpecialty_DOCTOR_SPECIALTY_PSYCHIATRY
	case service.SpecOrthopedics:
		return v1.DoctorSpecialty_DOCTOR_SPECIALTY_ORTHOPEDICS
	case service.SpecOphthalmology:
		return v1.DoctorSpecialty_DOCTOR_SPECIALTY_OPHTHALMOLOGY
	case service.SpecENT:
		return v1.DoctorSpecialty_DOCTOR_SPECIALTY_ENT
	default:
		return v1.DoctorSpecialty_DOCTOR_SPECIALTY_UNKNOWN
	}
}

func toGRPC(err error) error {
	if appErr, ok := err.(*apperrors.AppError); ok {
		return appErr.ToGRPC()
	}
	return status.Error(codes.Internal, err.Error())
}
