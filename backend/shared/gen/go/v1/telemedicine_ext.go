// telemedicine_ext.go — telemedicine-service Proto 스텁 (protoc 재생성 후 삭제 예정)
// TODO: protoc 재생성 후 이 파일을 제거하세요.
package v1

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ============================================================================
// 원격진료 enum 타입
// ============================================================================

// ConsultationStatus은 상담 상태입니다.
type ConsultationStatus int32

const (
	ConsultationStatus_CONSULTATION_STATUS_UNKNOWN     ConsultationStatus = 0
	ConsultationStatus_CONSULTATION_STATUS_REQUESTED   ConsultationStatus = 1
	ConsultationStatus_CONSULTATION_STATUS_MATCHED     ConsultationStatus = 2
	ConsultationStatus_CONSULTATION_STATUS_SCHEDULED   ConsultationStatus = 3
	ConsultationStatus_CONSULTATION_STATUS_IN_PROGRESS ConsultationStatus = 4
	ConsultationStatus_CONSULTATION_STATUS_COMPLETED   ConsultationStatus = 5
	ConsultationStatus_CONSULTATION_STATUS_CANCELLED   ConsultationStatus = 6
	ConsultationStatus_CONSULTATION_STATUS_NO_SHOW     ConsultationStatus = 7
)

// DoctorSpecialty — manpasik.pb.go에 이미 정의됨 (중복 제거됨)

// VideoSessionStatus는 비디오 세션 상태입니다.
type VideoSessionStatus int32

const (
	VideoSessionStatus_VIDEO_SESSION_STATUS_UNKNOWN   VideoSessionStatus = 0
	VideoSessionStatus_VIDEO_SESSION_STATUS_WAITING   VideoSessionStatus = 1
	VideoSessionStatus_VIDEO_SESSION_STATUS_CONNECTED VideoSessionStatus = 2
	VideoSessionStatus_VIDEO_SESSION_STATUS_ENDED     VideoSessionStatus = 3
	VideoSessionStatus_VIDEO_SESSION_STATUS_FAILED    VideoSessionStatus = 4
)

// ============================================================================
// 원격진료 메시지 타입
// ============================================================================

// Consultation은 원격진료 상담 메시지입니다.
type Consultation struct {
	ConsultationId  string                 `json:"consultation_id,omitempty"`
	PatientUserId   string                 `json:"patient_user_id,omitempty"`
	DoctorId        string                 `json:"doctor_id,omitempty"`
	Specialty       DoctorSpecialty        `json:"specialty,omitempty"`
	ChiefComplaint  string                 `json:"chief_complaint,omitempty"`
	Description     string                 `json:"description,omitempty"`
	Status          ConsultationStatus     `json:"status,omitempty"`
	Diagnosis       string                 `json:"diagnosis,omitempty"`
	DoctorNotes     string                 `json:"doctor_notes,omitempty"`
	PrescriptionId  string                 `json:"prescription_id,omitempty"`
	DurationMinutes int32                  `json:"duration_minutes,omitempty"`
	Rating          float64                `json:"rating,omitempty"`
	CreatedAt       *timestamppb.Timestamp `json:"created_at,omitempty"`
	ScheduledAt     *timestamppb.Timestamp `json:"scheduled_at,omitempty"`
	StartedAt       *timestamppb.Timestamp `json:"started_at,omitempty"`
	EndedAt         *timestamppb.Timestamp `json:"ended_at,omitempty"`
}

// CreateConsultationRequest은 상담 생성 요청입니다.
type CreateConsultationRequest struct {
	PatientUserId  string          `json:"patient_user_id,omitempty"`
	Specialty      DoctorSpecialty `json:"specialty,omitempty"`
	ChiefComplaint string          `json:"chief_complaint,omitempty"`
	Description    string          `json:"description,omitempty"`
}

// GetConsultationRequest은 상담 조회 요청입니다.
type GetConsultationRequest struct {
	ConsultationId string `json:"consultation_id,omitempty"`
}

// ListConsultationsRequest은 상담 목록 조회 요청입니다.
type ListConsultationsRequest struct {
	UserId       string             `json:"user_id,omitempty"`
	StatusFilter ConsultationStatus `json:"status_filter,omitempty"`
	Limit        int32              `json:"limit,omitempty"`
	Offset       int32              `json:"offset,omitempty"`
}

// ListConsultationsResponse은 상담 목록 조회 응답입니다.
type ListConsultationsResponse struct {
	Consultations []*Consultation `json:"consultations,omitempty"`
	TotalCount    int32           `json:"total_count,omitempty"`
}

// MatchDoctorRequest은 의사 매칭 요청입니다.
type MatchDoctorRequest struct {
	Specialty DoctorSpecialty `json:"specialty,omitempty"`
	Language  string          `json:"language,omitempty"`
}

// MatchDoctorResponse은 의사 매칭 응답입니다.
type MatchDoctorResponse struct {
	Doctors        []*DoctorProfile `json:"doctors,omitempty"`
	TotalAvailable int32            `json:"total_available,omitempty"`
}

// DoctorProfile은 의사 프로필입니다.
type DoctorProfile struct {
	DoctorId           string          `json:"doctor_id,omitempty"`
	Name               string          `json:"name,omitempty"`
	Specialty          DoctorSpecialty `json:"specialty,omitempty"`
	Hospital           string          `json:"hospital,omitempty"`
	LicenseNumber      string          `json:"license_number,omitempty"`
	ExperienceYears    int32           `json:"experience_years,omitempty"`
	Rating             float64         `json:"rating,omitempty"`
	TotalConsultations int32           `json:"total_consultations,omitempty"`
	IsAvailable        bool            `json:"is_available,omitempty"`
	Languages          []string        `json:"languages,omitempty"`
	ProfileImageUrl    string          `json:"profile_image_url,omitempty"`
}

// StartVideoSessionRequest은 비디오 세션 시작 요청입니다.
type StartVideoSessionRequest struct {
	ConsultationId string `json:"consultation_id,omitempty"`
	UserId         string `json:"user_id,omitempty"`
}

// VideoSession은 비디오 세션 메시지입니다.
type VideoSession struct {
	SessionId       string                 `json:"session_id,omitempty"`
	ConsultationId  string                 `json:"consultation_id,omitempty"`
	RoomUrl         string                 `json:"room_url,omitempty"`
	Token           string                 `json:"token,omitempty"`
	Status          VideoSessionStatus     `json:"status,omitempty"`
	DurationSeconds int32                  `json:"duration_seconds,omitempty"`
	StartedAt       *timestamppb.Timestamp `json:"started_at,omitempty"`
	EndedAt         *timestamppb.Timestamp `json:"ended_at,omitempty"`
}

// EndVideoSessionRequest은 비디오 세션 종료 요청입니다.
type EndVideoSessionRequest struct {
	SessionId      string `json:"session_id,omitempty"`
	ConsultationId string `json:"consultation_id,omitempty"`
	DoctorNotes    string `json:"doctor_notes,omitempty"`
	Diagnosis      string `json:"diagnosis,omitempty"`
}

// RateConsultationRequest은 상담 평점 요청입니다.
type RateConsultationRequest struct {
	ConsultationId string  `json:"consultation_id,omitempty"`
	Rating         float64 `json:"rating,omitempty"`
}

// RateConsultationResponse은 상담 평점 응답입니다.
type RateConsultationResponse struct {
	Success          bool    `json:"success,omitempty"`
	NewAverageRating float64 `json:"new_average_rating,omitempty"`
}

// ============================================================================
// TelemedicineService gRPC 인터페이스
// ============================================================================

// TelemedicineServiceServer는 원격진료 서비스 gRPC 서버 인터페이스입니다.
type TelemedicineServiceServer interface {
	CreateConsultation(context.Context, *CreateConsultationRequest) (*Consultation, error)
	GetConsultation(context.Context, *GetConsultationRequest) (*Consultation, error)
	ListConsultations(context.Context, *ListConsultationsRequest) (*ListConsultationsResponse, error)
	MatchDoctor(context.Context, *MatchDoctorRequest) (*MatchDoctorResponse, error)
	StartVideoSession(context.Context, *StartVideoSessionRequest) (*VideoSession, error)
	EndVideoSession(context.Context, *EndVideoSessionRequest) (*VideoSession, error)
	RateConsultation(context.Context, *RateConsultationRequest) (*RateConsultationResponse, error)
}

// UnimplementedTelemedicineServiceServer는 미구현 서버입니다.
type UnimplementedTelemedicineServiceServer struct{}

func (UnimplementedTelemedicineServiceServer) CreateConsultation(context.Context, *CreateConsultationRequest) (*Consultation, error) {
	return nil, nil
}
func (UnimplementedTelemedicineServiceServer) GetConsultation(context.Context, *GetConsultationRequest) (*Consultation, error) {
	return nil, nil
}
func (UnimplementedTelemedicineServiceServer) ListConsultations(context.Context, *ListConsultationsRequest) (*ListConsultationsResponse, error) {
	return nil, nil
}
func (UnimplementedTelemedicineServiceServer) MatchDoctor(context.Context, *MatchDoctorRequest) (*MatchDoctorResponse, error) {
	return nil, nil
}
func (UnimplementedTelemedicineServiceServer) StartVideoSession(context.Context, *StartVideoSessionRequest) (*VideoSession, error) {
	return nil, nil
}
func (UnimplementedTelemedicineServiceServer) EndVideoSession(context.Context, *EndVideoSessionRequest) (*VideoSession, error) {
	return nil, nil
}
func (UnimplementedTelemedicineServiceServer) RateConsultation(context.Context, *RateConsultationRequest) (*RateConsultationResponse, error) {
	return nil, nil
}

// RegisterTelemedicineServiceServer는 gRPC 서버에 원격진료 서비스를 등록합니다.
func RegisterTelemedicineServiceServer(s *grpc.Server, srv TelemedicineServiceServer) {
	s.RegisterService(&_TelemedicineService_serviceDesc, srv)
}

var _TelemedicineService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "manpasik.v1.TelemedicineService",
	HandlerType: (*TelemedicineServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "CreateConsultation", Handler: _TelemedicineService_CreateConsultation_Handler},
		{MethodName: "GetConsultation", Handler: _TelemedicineService_GetConsultation_Handler},
		{MethodName: "ListConsultations", Handler: _TelemedicineService_ListConsultations_Handler},
		{MethodName: "MatchDoctor", Handler: _TelemedicineService_MatchDoctor_Handler},
		{MethodName: "StartVideoSession", Handler: _TelemedicineService_StartVideoSession_Handler},
		{MethodName: "EndVideoSession", Handler: _TelemedicineService_EndVideoSession_Handler},
		{MethodName: "RateConsultation", Handler: _TelemedicineService_RateConsultation_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "manpasik.proto",
}

func _TelemedicineService_CreateConsultation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateConsultationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	return srv.(TelemedicineServiceServer).CreateConsultation(ctx, in)
}

func _TelemedicineService_GetConsultation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetConsultationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	return srv.(TelemedicineServiceServer).GetConsultation(ctx, in)
}

func _TelemedicineService_ListConsultations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListConsultationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	return srv.(TelemedicineServiceServer).ListConsultations(ctx, in)
}

func _TelemedicineService_MatchDoctor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MatchDoctorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	return srv.(TelemedicineServiceServer).MatchDoctor(ctx, in)
}

func _TelemedicineService_StartVideoSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartVideoSessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	return srv.(TelemedicineServiceServer).StartVideoSession(ctx, in)
}

func _TelemedicineService_EndVideoSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EndVideoSessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	return srv.(TelemedicineServiceServer).EndVideoSession(ctx, in)
}

func _TelemedicineService_RateConsultation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RateConsultationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	return srv.(TelemedicineServiceServer).RateConsultation(ctx, in)
}
