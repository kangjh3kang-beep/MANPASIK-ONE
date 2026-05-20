package handler

import (
	"context"
	"github.com/manpasik/backend/services/concept-service/internal/service"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ConceptHandler struct {
	v1.UnimplementedConceptServiceServer
	v1.UnimplementedOrganizationServiceServer
	svc *service.ConceptService
	log *zap.Logger
}

func NewConceptHandler(svc *service.ConceptService, log *zap.Logger) *ConceptHandler {
	return &ConceptHandler{svc: svc, log: log}
}

func (h *ConceptHandler) ListConcepts(ctx context.Context, req *v1.ListConceptsRequest) (*v1.ListConceptsResponse, error) {
	list, err := h.svc.ListConcepts(ctx)
	if err != nil { return nil, status.Error(codes.Internal, err.Error()) }
	out := make([]*v1.Concept, 0, len(list))
	for _, c := range list { out = append(out, conceptToProto(c)) }
	return &v1.ListConceptsResponse{Concepts: out}, nil
}
func (h *ConceptHandler) GetConcept(ctx context.Context, req *v1.GetConceptRequest) (*v1.Concept, error) {
	c, err := h.svc.GetConcept(ctx, req.GetConceptId())
	if err != nil { return nil, toGRPC(err) }
	return conceptToProto(c), nil
}
func (h *ConceptHandler) CreateConcept(ctx context.Context, req *v1.CreateConceptRequest) (*v1.Concept, error) {
	c, err := h.svc.CreateConcept(ctx, req.GetName(), req.GetDescription(), req.GetCategory(), "", req.GetOwnerId())
	if err != nil { return nil, toGRPC(err) }
	return conceptToProto(c), nil
}
func (h *ConceptHandler) AssignDeviceToConcept(ctx context.Context, req *v1.AssignDeviceToConceptRequest) (*v1.AssignDeviceToConceptResponse, error) {
	if err := h.svc.AssignDevice(ctx, req.GetConceptId(), req.GetDeviceId()); err != nil { return nil, toGRPC(err) }
	return &v1.AssignDeviceToConceptResponse{Success: true}, nil
}
func (h *ConceptHandler) GetConceptStats(ctx context.Context, req *v1.GetConceptStatsRequest) (*v1.ConceptStats, error) {
	devices, _, _, err := h.svc.GetConceptStats(ctx, req.GetConceptId())
	if err != nil { return nil, toGRPC(err) }
	return &v1.ConceptStats{ConceptId: req.GetConceptId(), DeviceCount: devices, MeasurementCount: 0, AvgHealthScore: 0.0}, nil
}
func (h *ConceptHandler) GetConceptDashboard(ctx context.Context, req *v1.GetConceptDashboardRequest) (*v1.ConceptDashboard, error) {
	devices, _, _, _, err := h.svc.GetConceptDashboard(ctx, req.GetConceptId())
	if err != nil { return nil, toGRPC(err) }
	stats := &v1.ConceptStats{ConceptId: req.GetConceptId(), DeviceCount: devices}
	return &v1.ConceptDashboard{ConceptId: req.GetConceptId(), Stats: stats, RecentAlerts: []string{}}, nil
}

func (h *ConceptHandler) CreateOrganization(ctx context.Context, req *v1.CreateOrganizationRequest) (*v1.Organization, error) {
	o, err := h.svc.CreateOrganization(ctx, req.GetName(), req.GetDescription())
	if err != nil { return nil, toGRPC(err) }
	return orgToProto(o), nil
}
func (h *ConceptHandler) GetOrganization(ctx context.Context, req *v1.GetOrganizationRequest) (*v1.Organization, error) {
	o, err := h.svc.GetOrganization(ctx, req.GetOrganizationId())
	if err != nil { return nil, toGRPC(err) }
	return orgToProto(o), nil
}
func (h *ConceptHandler) ListOrganizations(ctx context.Context, req *v1.ListOrganizationsRequest) (*v1.ListOrganizationsResponse, error) {
	list, err := h.svc.ListOrganizations(ctx)
	if err != nil { return nil, status.Error(codes.Internal, err.Error()) }
	out := make([]*v1.Organization, 0, len(list))
	for _, o := range list { out = append(out, orgToProto(o)) }
	return &v1.ListOrganizationsResponse{Organizations: out}, nil
}
func (h *ConceptHandler) AddMember(ctx context.Context, req *v1.AddOrgMemberRequest) (*v1.AddOrgMemberResponse, error) {
	if err := h.svc.AddMember(ctx, req.GetOrganizationId(), req.GetUserId()); err != nil { return nil, toGRPC(err) }
	return &v1.AddOrgMemberResponse{Success: true}, nil
}
func (h *ConceptHandler) RemoveMember(ctx context.Context, req *v1.RemoveOrgMemberRequest) (*v1.RemoveOrgMemberResponse, error) {
	if err := h.svc.RemoveMember(ctx, req.GetOrganizationId(), req.GetUserId()); err != nil { return nil, toGRPC(err) }
	return &v1.RemoveOrgMemberResponse{Success: true}, nil
}

func conceptToProto(c *service.Concept) *v1.Concept {
	return &v1.Concept{Id: c.ID, Name: c.Name, Description: c.Description, Category: c.Category,
		IconUrl: c.IconURL, OwnerId: c.OwnerID, DeviceIds: c.DeviceIDs, CreatedAt: timestamppb.New(c.CreatedAt)}
}
func orgToProto(o *service.Organization) *v1.Organization {
	return &v1.Organization{Id: o.ID, Name: o.Name, Description: o.Description,
		OwnerId: o.OwnerID, MemberCount: int32(len(o.MemberIDs)), CreatedAt: timestamppb.New(o.CreatedAt)}
}
func toGRPC(err error) error {
	switch err {
	case service.ErrNotFound: return status.Error(codes.NotFound, err.Error())
	case service.ErrInvalidInput: return status.Error(codes.InvalidArgument, err.Error())
	case service.ErrAlreadyExists: return status.Error(codes.AlreadyExists, err.Error())
	default: return status.Error(codes.Internal, err.Error())
	}
}
