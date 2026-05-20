package handler

import (
	"context"
	"github.com/manpasik/backend/services/voice-profile-service/internal/service"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type VoiceProfileHandler struct {
	v1.UnimplementedVoiceProfileServiceServer
	svc *service.VoiceProfileService
	log *zap.Logger
}

func NewVoiceProfileHandler(svc *service.VoiceProfileService, log *zap.Logger) *VoiceProfileHandler {
	return &VoiceProfileHandler{svc: svc, log: log}
}

func (h *VoiceProfileHandler) CreateVoiceProfile(ctx context.Context, req *v1.CreateVoiceProfileRequest) (*v1.VoiceProfile, error) {
	// VoiceSample is []byte in proto; service accepts string URL — pass empty for now
	p, err := h.svc.Create(ctx, req.GetUserId(), req.GetProfileName(), req.GetLanguage(), "")
	if err != nil {
		return nil, toGRPC(err)
	}
	return profileToProto(p), nil
}
func (h *VoiceProfileHandler) GetVoiceProfile(ctx context.Context, req *v1.GetVoiceProfileRequest) (*v1.VoiceProfile, error) {
	p, err := h.svc.Get(ctx, req.GetProfileId())
	if err != nil {
		return nil, toGRPC(err)
	}
	return profileToProto(p), nil
}
func (h *VoiceProfileHandler) ListVoiceProfiles(ctx context.Context, req *v1.ListVoiceProfilesRequest) (*v1.ListVoiceProfilesResponse, error) {
	list, err := h.svc.List(ctx, req.GetUserId())
	if err != nil {
		return nil, toGRPC(err)
	}
	out := make([]*v1.VoiceProfile, 0, len(list))
	for _, p := range list {
		out = append(out, profileToProto(p))
	}
	return &v1.ListVoiceProfilesResponse{Profiles: out}, nil
}
func (h *VoiceProfileHandler) SynthesizeTranslation(ctx context.Context, req *v1.SynthesizeTranslationRequest) (*v1.SynthesizeTranslationResponse, error) {
	audioURL, translatedText, err := h.svc.SynthesizeTranslation(ctx, req.GetProfileId(), req.GetText(), req.GetTargetLanguage())
	if err != nil {
		return nil, toGRPC(err)
	}
	// 실제 TTS 모델 연동 전까지 audioURL만 반환, 바이트 데이터는 빈 상태
	_ = audioURL
	return &v1.SynthesizeTranslationResponse{
		AudioData:       nil,
		TranslatedText:  translatedText,
		Language:        req.GetTargetLanguage(),
		SimilarityScore: 0.0,
	}, nil
}
func (h *VoiceProfileHandler) DeleteVoiceProfile(ctx context.Context, req *v1.DeleteVoiceProfileRequest) (*v1.DeleteVoiceProfileResponse, error) {
	if err := h.svc.Delete(ctx, req.GetProfileId()); err != nil {
		return nil, toGRPC(err)
	}
	return &v1.DeleteVoiceProfileResponse{Success: true}, nil
}

func profileToProto(p *service.VoiceProfile) *v1.VoiceProfile {
	return &v1.VoiceProfile{ProfileId: p.ProfileID, UserId: p.UserID, ProfileName: p.ProfileName,
		Language: p.Language, Status: p.Status, ModelUrl: p.ModelURL, CreatedAt: timestamppb.New(p.CreatedAt)}
}
func toGRPC(err error) error {
	switch err {
	case service.ErrNotFound:
		return status.Error(codes.NotFound, err.Error())
	case service.ErrInvalidInput:
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
