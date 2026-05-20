#!/bin/bash
set -e
BASE="/home/kangjh3kang/Manpasik/backend/services"

# ============================================================================
# 1. concept-service (포트 50078)
# ============================================================================
SVC="$BASE/concept-service"
mkdir -p "$SVC/cmd" "$SVC/internal/handler" "$SVC/internal/service" "$SVC/internal/repository/memory"

cat > "$SVC/cmd/main.go" << 'GOEOF'
package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"net/http"

	"github.com/manpasik/backend/services/concept-service/internal/handler"
	"github.com/manpasik/backend/services/concept-service/internal/repository/memory"
	"github.com/manpasik/backend/services/concept-service/internal/service"
	"github.com/manpasik/backend/shared/config"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"github.com/manpasik/backend/shared/middleware"
	"github.com/manpasik/backend/shared/observability"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const serviceName = "concept-service"

func main() {
	cfg := config.LoadFromEnv(serviceName)
	logger, err := zap.NewProduction()
	if err != nil { logger = zap.NewNop() }
	defer logger.Sync()

	metrics := observability.NewMetrics()
	healthCheck := observability.NewHealthCheck(serviceName, cfg.Version)
	log.Printf("[%s] Starting v%s...", serviceName, cfg.Version)

	conceptRepo := memory.NewConceptRepository()
	orgRepo := memory.NewOrganizationRepository()
	svc := service.NewConceptService(logger, conceptRepo, orgRepo)
	h := handler.NewConceptHandler(svc, logger)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDInterceptor(),
			observability.UnaryServerInterceptor(metrics),
		),
	)
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_SERVING)
	v1.RegisterConceptServiceServer(grpcServer, h)
	v1.RegisterOrganizationServiceServer(grpcServer, h)
	reflection.Register(grpcServer)

	grpcPort := cfg.GRPCPort
	if grpcPort == "" { grpcPort = ":50078" }
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil { log.Fatalf("[%s] Failed to listen: %v", serviceName, err) }

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		log.Printf("[%s] Received signal %v, shutting down...", serviceName, sig)
		healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_NOT_SERVING)
		go func() { time.Sleep(cfg.ShutdownTimeout); os.Exit(1) }()
		grpcServer.GracefulStop(); cancel()
	}()
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/metrics", metrics.PrometheusHandler())
		mux.HandleFunc("/health", healthCheck.Handler())
		logger.Info("Metrics server starting", zap.String("addr", ":9100"))
		http.ListenAndServe(":9100", mux)
	}()
	log.Printf("[%s] gRPC server listening on %s", serviceName, grpcPort)
	if err := grpcServer.Serve(lis); err != nil { log.Fatalf("[%s] Failed to serve: %v", serviceName, err) }
	<-ctx.Done()
}
GOEOF

cat > "$SVC/internal/service/concept.go" << 'GOEOF'
package service

import (
	"context"
	"errors"
	"sync"
	"time"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalidInput  = errors.New("invalid input")
)

type Concept struct {
	ID, Name, Description, Category, IconURL, OwnerID string
	DeviceIDs                                         []string
	CreatedAt                                         time.Time
}
type Organization struct {
	ID, Name, Description string
	MemberIDs             []string
	CreatedAt             time.Time
}

type ConceptRepository interface {
	List(ctx context.Context) ([]*Concept, error)
	GetByID(ctx context.Context, id string) (*Concept, error)
	Create(ctx context.Context, c *Concept) error
	Update(ctx context.Context, c *Concept) error
}
type OrganizationRepository interface {
	List(ctx context.Context) ([]*Organization, error)
	GetByID(ctx context.Context, id string) (*Organization, error)
	Create(ctx context.Context, o *Organization) error
	Update(ctx context.Context, o *Organization) error
	Delete(ctx context.Context, id string) error
}

type ConceptService struct {
	mu          sync.RWMutex
	log         *zap.Logger
	concepts    ConceptRepository
	orgs        OrganizationRepository
}

func NewConceptService(log *zap.Logger, cr ConceptRepository, or OrganizationRepository) *ConceptService {
	return &ConceptService{log: log, concepts: cr, orgs: or}
}

func (s *ConceptService) ListConcepts(ctx context.Context) ([]*Concept, error) {
	return s.concepts.List(ctx)
}
func (s *ConceptService) GetConcept(ctx context.Context, id string) (*Concept, error) {
	if id == "" { return nil, ErrInvalidInput }
	return s.concepts.GetByID(ctx, id)
}
func (s *ConceptService) CreateConcept(ctx context.Context, name, description, category, iconURL, ownerID string) (*Concept, error) {
	if name == "" { return nil, ErrInvalidInput }
	c := &Concept{ID: uuid.New().String(), Name: name, Description: description, Category: category,
		IconURL: iconURL, OwnerID: ownerID, CreatedAt: time.Now().UTC()}
	return c, s.concepts.Create(ctx, c)
}
func (s *ConceptService) AssignDevice(ctx context.Context, conceptID, deviceID string) error {
	if conceptID == "" || deviceID == "" { return ErrInvalidInput }
	c, err := s.concepts.GetByID(ctx, conceptID)
	if err != nil { return err }
	for _, d := range c.DeviceIDs { if d == deviceID { return ErrAlreadyExists } }
	c.DeviceIDs = append(c.DeviceIDs, deviceID)
	return s.concepts.Update(ctx, c)
}
func (s *ConceptService) GetConceptStats(ctx context.Context, id string) (int32, int32, float64, error) {
	if id == "" { return 0, 0, 0, ErrInvalidInput }
	c, err := s.concepts.GetByID(ctx, id)
	if err != nil { return 0, 0, 0, err }
	return int32(len(c.DeviceIDs)), 0, 0.0, nil
}
func (s *ConceptService) GetConceptDashboard(ctx context.Context, id string) (int32, int32, float64, []string, error) {
	if id == "" { return 0, 0, 0, nil, ErrInvalidInput }
	c, err := s.concepts.GetByID(ctx, id)
	if err != nil { return 0, 0, 0, nil, err }
	return int32(len(c.DeviceIDs)), 0, 0.0, []string{}, nil
}

func (s *ConceptService) CreateOrganization(ctx context.Context, name, description string) (*Organization, error) {
	if name == "" { return nil, ErrInvalidInput }
	o := &Organization{ID: uuid.New().String(), Name: name, Description: description, CreatedAt: time.Now().UTC()}
	return o, s.orgs.Create(ctx, o)
}
func (s *ConceptService) GetOrganization(ctx context.Context, id string) (*Organization, error) {
	if id == "" { return nil, ErrInvalidInput }
	return s.orgs.GetByID(ctx, id)
}
func (s *ConceptService) ListOrganizations(ctx context.Context) ([]*Organization, error) {
	return s.orgs.List(ctx)
}
func (s *ConceptService) AddMember(ctx context.Context, orgID, memberID string) error {
	if orgID == "" || memberID == "" { return ErrInvalidInput }
	o, err := s.orgs.GetByID(ctx, orgID)
	if err != nil { return err }
	for _, m := range o.MemberIDs { if m == memberID { return ErrAlreadyExists } }
	o.MemberIDs = append(o.MemberIDs, memberID)
	return s.orgs.Update(ctx, o)
}
func (s *ConceptService) RemoveMember(ctx context.Context, orgID, memberID string) error {
	if orgID == "" || memberID == "" { return ErrInvalidInput }
	o, err := s.orgs.GetByID(ctx, orgID)
	if err != nil { return err }
	newM := make([]string, 0)
	found := false
	for _, m := range o.MemberIDs {
		if m == memberID { found = true; continue }
		newM = append(newM, m)
	}
	if !found { return ErrNotFound }
	o.MemberIDs = newM
	return s.orgs.Update(ctx, o)
}
GOEOF

cat > "$SVC/internal/repository/memory/concept.go" << 'GOEOF'
package memory

import (
	"context"
	"sync"
	"github.com/manpasik/backend/services/concept-service/internal/service"
)

type ConceptRepository struct {
	mu   sync.RWMutex
	data map[string]*service.Concept
}
func NewConceptRepository() *ConceptRepository {
	return &ConceptRepository{data: make(map[string]*service.Concept)}
}
func (r *ConceptRepository) List(_ context.Context) ([]*service.Concept, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	out := make([]*service.Concept, 0, len(r.data))
	for _, v := range r.data { out = append(out, v) }
	return out, nil
}
func (r *ConceptRepository) GetByID(_ context.Context, id string) (*service.Concept, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	if v, ok := r.data[id]; ok { return v, nil }
	return nil, service.ErrNotFound
}
func (r *ConceptRepository) Create(_ context.Context, c *service.Concept) error {
	r.mu.Lock(); defer r.mu.Unlock()
	r.data[c.ID] = c; return nil
}
func (r *ConceptRepository) Update(_ context.Context, c *service.Concept) error {
	r.mu.Lock(); defer r.mu.Unlock()
	if _, ok := r.data[c.ID]; !ok { return service.ErrNotFound }
	r.data[c.ID] = c; return nil
}

type OrganizationRepository struct {
	mu   sync.RWMutex
	data map[string]*service.Organization
}
func NewOrganizationRepository() *OrganizationRepository {
	return &OrganizationRepository{data: make(map[string]*service.Organization)}
}
func (r *OrganizationRepository) List(_ context.Context) ([]*service.Organization, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	out := make([]*service.Organization, 0, len(r.data))
	for _, v := range r.data { out = append(out, v) }
	return out, nil
}
func (r *OrganizationRepository) GetByID(_ context.Context, id string) (*service.Organization, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	if v, ok := r.data[id]; ok { return v, nil }
	return nil, service.ErrNotFound
}
func (r *OrganizationRepository) Create(_ context.Context, o *service.Organization) error {
	r.mu.Lock(); defer r.mu.Unlock()
	r.data[o.ID] = o; return nil
}
func (r *OrganizationRepository) Update(_ context.Context, o *service.Organization) error {
	r.mu.Lock(); defer r.mu.Unlock()
	if _, ok := r.data[o.ID]; !ok { return service.ErrNotFound }
	r.data[o.ID] = o; return nil
}
func (r *OrganizationRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock(); defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok { return service.ErrNotFound }
	delete(r.data, id); return nil
}
GOEOF

cat > "$SVC/internal/handler/grpc.go" << 'GOEOF'
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
	return &v1.ListConceptsResponse{Concepts: out, Total: int32(len(out))}, nil
}
func (h *ConceptHandler) GetConcept(ctx context.Context, req *v1.GetConceptRequest) (*v1.Concept, error) {
	c, err := h.svc.GetConcept(ctx, req.GetConceptId())
	if err != nil { return nil, toGRPC(err) }
	return conceptToProto(c), nil
}
func (h *ConceptHandler) CreateConcept(ctx context.Context, req *v1.CreateConceptRequest) (*v1.Concept, error) {
	c, err := h.svc.CreateConcept(ctx, req.GetName(), req.GetDescription(), req.GetCategory(), req.GetIconUrl(), req.GetOwnerId())
	if err != nil { return nil, toGRPC(err) }
	return conceptToProto(c), nil
}
func (h *ConceptHandler) AssignDeviceToConcept(ctx context.Context, req *v1.AssignDeviceToConceptRequest) (*v1.AssignDeviceToConceptResponse, error) {
	if err := h.svc.AssignDevice(ctx, req.GetConceptId(), req.GetDeviceId()); err != nil { return nil, toGRPC(err) }
	return &v1.AssignDeviceToConceptResponse{Success: true}, nil
}
func (h *ConceptHandler) GetConceptStats(ctx context.Context, req *v1.GetConceptStatsRequest) (*v1.ConceptStats, error) {
	devices, users, completion, err := h.svc.GetConceptStats(ctx, req.GetConceptId())
	if err != nil { return nil, toGRPC(err) }
	return &v1.ConceptStats{ConceptId: req.GetConceptId(), TotalDevices: devices, ActiveUsers: users, CompletionRate: completion}, nil
}
func (h *ConceptHandler) GetConceptDashboard(ctx context.Context, req *v1.GetConceptDashboardRequest) (*v1.ConceptDashboard, error) {
	devices, users, completion, events, err := h.svc.GetConceptDashboard(ctx, req.GetConceptId())
	if err != nil { return nil, toGRPC(err) }
	return &v1.ConceptDashboard{ConceptId: req.GetConceptId(), TotalDevices: devices, ActiveUsers: users, CompletionRate: completion, RecentEvents: events}, nil
}

func (h *ConceptHandler) CreateOrganization(ctx context.Context, req *v1.CreateOrganizationRequest) (*v1.Organization, error) {
	o, err := h.svc.CreateOrganization(ctx, req.GetName(), req.GetDescription())
	if err != nil { return nil, toGRPC(err) }
	return orgToProto(o), nil
}
func (h *ConceptHandler) GetOrganization(ctx context.Context, req *v1.GetOrganizationRequest) (*v1.Organization, error) {
	o, err := h.svc.GetOrganization(ctx, req.GetOrgId())
	if err != nil { return nil, toGRPC(err) }
	return orgToProto(o), nil
}
func (h *ConceptHandler) ListOrganizations(ctx context.Context, req *v1.ListOrganizationsRequest) (*v1.ListOrganizationsResponse, error) {
	list, err := h.svc.ListOrganizations(ctx)
	if err != nil { return nil, status.Error(codes.Internal, err.Error()) }
	out := make([]*v1.Organization, 0, len(list))
	for _, o := range list { out = append(out, orgToProto(o)) }
	return &v1.ListOrganizationsResponse{Organizations: out, Total: int32(len(out))}, nil
}
func (h *ConceptHandler) AddMember(ctx context.Context, req *v1.AddOrgMemberRequest) (*v1.AddOrgMemberResponse, error) {
	if err := h.svc.AddMember(ctx, req.GetOrgId(), req.GetMemberId()); err != nil { return nil, toGRPC(err) }
	return &v1.AddOrgMemberResponse{Success: true}, nil
}
func (h *ConceptHandler) RemoveMember(ctx context.Context, req *v1.RemoveOrgMemberRequest) (*v1.RemoveOrgMemberResponse, error) {
	if err := h.svc.RemoveMember(ctx, req.GetOrgId(), req.GetMemberId()); err != nil { return nil, toGRPC(err) }
	return &v1.RemoveOrgMemberResponse{Success: true}, nil
}

func conceptToProto(c *service.Concept) *v1.Concept {
	return &v1.Concept{Id: c.ID, Name: c.Name, Description: c.Description, Category: c.Category,
		IconUrl: c.IconURL, OwnerId: c.OwnerID, DeviceIds: c.DeviceIDs, CreatedAt: timestamppb.New(c.CreatedAt)}
}
func orgToProto(o *service.Organization) *v1.Organization {
	return &v1.Organization{Id: o.ID, Name: o.Name, Description: o.Description,
		MemberIds: o.MemberIDs, CreatedAt: timestamppb.New(o.CreatedAt)}
}
func toGRPC(err error) error {
	switch err {
	case service.ErrNotFound: return status.Error(codes.NotFound, err.Error())
	case service.ErrInvalidInput: return status.Error(codes.InvalidArgument, err.Error())
	case service.ErrAlreadyExists: return status.Error(codes.AlreadyExists, err.Error())
	default: return status.Error(codes.Internal, err.Error())
	}
}
GOEOF

cat > "$SVC/internal/service/concept_test.go" << 'GOEOF'
package service_test

import (
	"context"
	"testing"
	"github.com/manpasik/backend/services/concept-service/internal/repository/memory"
	"github.com/manpasik/backend/services/concept-service/internal/service"
	"go.uber.org/zap"
)

func newSvc() *service.ConceptService {
	return service.NewConceptService(zap.NewNop(), memory.NewConceptRepository(), memory.NewOrganizationRepository())
}

func TestCreateConcept(t *testing.T) {
	s := newSvc()
	c, err := s.CreateConcept(context.Background(), "테스트", "설명", "personal", "", "u1")
	if err != nil { t.Fatal(err) }
	if c.ID == "" || c.Name != "테스트" { t.Error("unexpected concept") }
}
func TestListConcepts(t *testing.T) {
	s := newSvc(); ctx := context.Background()
	s.CreateConcept(ctx, "A", "", "", "", "u1")
	s.CreateConcept(ctx, "B", "", "", "", "u1")
	list, _ := s.ListConcepts(ctx)
	if len(list) != 2 { t.Errorf("expected 2, got %d", len(list)) }
}
func TestAssignDevice(t *testing.T) {
	s := newSvc(); ctx := context.Background()
	c, _ := s.CreateConcept(ctx, "X", "", "", "", "u1")
	if err := s.AssignDevice(ctx, c.ID, "dev1"); err != nil { t.Fatal(err) }
	got, _ := s.GetConcept(ctx, c.ID)
	if len(got.DeviceIDs) != 1 { t.Error("device not assigned") }
}
func TestCreateOrg(t *testing.T) {
	s := newSvc()
	o, err := s.CreateOrganization(context.Background(), "병원", "테스트")
	if err != nil { t.Fatal(err) }
	if o.ID == "" { t.Error("empty id") }
}
func TestAddRemoveMember(t *testing.T) {
	s := newSvc(); ctx := context.Background()
	o, _ := s.CreateOrganization(ctx, "병원", "")
	s.AddMember(ctx, o.ID, "m1")
	got, _ := s.GetOrganization(ctx, o.ID)
	if len(got.MemberIDs) != 1 { t.Error("member not added") }
	s.RemoveMember(ctx, o.ID, "m1")
	got2, _ := s.GetOrganization(ctx, o.ID)
	if len(got2.MemberIDs) != 0 { t.Error("member not removed") }
}
GOEOF

cat > "$SVC/Dockerfile" << 'GOEOF'
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /concept-service ./services/concept-service/cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /concept-service /concept-service
EXPOSE 50078
ENTRYPOINT ["/concept-service"]
GOEOF

echo "[OK] concept-service created"

# ============================================================================
# 2. voice-profile-service (포트 50081)
# ============================================================================
SVC="$BASE/voice-profile-service"
mkdir -p "$SVC/cmd" "$SVC/internal/handler" "$SVC/internal/service" "$SVC/internal/repository/memory"

cat > "$SVC/cmd/main.go" << 'GOEOF'
package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"net/http"

	"github.com/manpasik/backend/services/voice-profile-service/internal/handler"
	"github.com/manpasik/backend/services/voice-profile-service/internal/repository/memory"
	"github.com/manpasik/backend/services/voice-profile-service/internal/service"
	"github.com/manpasik/backend/shared/config"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"github.com/manpasik/backend/shared/middleware"
	"github.com/manpasik/backend/shared/observability"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const serviceName = "voice-profile-service"

func main() {
	cfg := config.LoadFromEnv(serviceName)
	logger, err := zap.NewProduction()
	if err != nil { logger = zap.NewNop() }
	defer logger.Sync()

	metrics := observability.NewMetrics()
	healthCheck := observability.NewHealthCheck(serviceName, cfg.Version)
	log.Printf("[%s] Starting v%s...", serviceName, cfg.Version)

	repo := memory.NewVoiceProfileRepository()
	svc := service.NewVoiceProfileService(logger, repo)
	h := handler.NewVoiceProfileHandler(svc, logger)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDInterceptor(),
			observability.UnaryServerInterceptor(metrics),
		),
	)
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_SERVING)
	v1.RegisterVoiceProfileServiceServer(grpcServer, h)
	reflection.Register(grpcServer)

	grpcPort := cfg.GRPCPort
	if grpcPort == "" { grpcPort = ":50081" }
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil { log.Fatalf("[%s] Failed to listen: %v", serviceName, err) }

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		log.Printf("[%s] Received signal %v, shutting down...", serviceName, sig)
		healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_NOT_SERVING)
		go func() { time.Sleep(cfg.ShutdownTimeout); os.Exit(1) }()
		grpcServer.GracefulStop(); cancel()
	}()
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/metrics", metrics.PrometheusHandler())
		mux.HandleFunc("/health", healthCheck.Handler())
		http.ListenAndServe(":9100", mux)
	}()
	log.Printf("[%s] gRPC server listening on %s", serviceName, grpcPort)
	if err := grpcServer.Serve(lis); err != nil { log.Fatalf("[%s] Failed to serve: %v", serviceName, err) }
	<-ctx.Done()
}
GOEOF

cat > "$SVC/internal/service/voice_profile.go" << 'GOEOF'
package service

import (
	"context"
	"errors"
	"time"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
)

type VoiceProfile struct {
	ProfileID, UserID, ProfileName, Language, Status, ModelURL string
	CreatedAt                                                  time.Time
}

type VoiceProfileRepository interface {
	Create(ctx context.Context, p *VoiceProfile) error
	GetByID(ctx context.Context, profileID string) (*VoiceProfile, error)
	ListByUser(ctx context.Context, userID string) ([]*VoiceProfile, error)
	Delete(ctx context.Context, profileID string) error
}

type VoiceProfileService struct {
	log  *zap.Logger
	repo VoiceProfileRepository
}

func NewVoiceProfileService(log *zap.Logger, repo VoiceProfileRepository) *VoiceProfileService {
	return &VoiceProfileService{log: log, repo: repo}
}

func (s *VoiceProfileService) Create(ctx context.Context, userID, profileName, language, voiceSampleURL string) (*VoiceProfile, error) {
	if userID == "" { return nil, ErrInvalidInput }
	if language == "" { language = "ko" }
	p := &VoiceProfile{
		ProfileID: uuid.New().String(), UserID: userID, ProfileName: profileName,
		Language: language, Status: "processing", CreatedAt: time.Now().UTC(),
	}
	return p, s.repo.Create(ctx, p)
}
func (s *VoiceProfileService) Get(ctx context.Context, profileID string) (*VoiceProfile, error) {
	if profileID == "" { return nil, ErrInvalidInput }
	return s.repo.GetByID(ctx, profileID)
}
func (s *VoiceProfileService) List(ctx context.Context, userID string) ([]*VoiceProfile, error) {
	if userID == "" { return nil, ErrInvalidInput }
	return s.repo.ListByUser(ctx, userID)
}
func (s *VoiceProfileService) Delete(ctx context.Context, profileID string) error {
	if profileID == "" { return ErrInvalidInput }
	return s.repo.Delete(ctx, profileID)
}
func (s *VoiceProfileService) SynthesizeTranslation(ctx context.Context, profileID, text, targetLang string) (string, string, error) {
	if profileID == "" || text == "" { return "", "", ErrInvalidInput }
	p, err := s.repo.GetByID(ctx, profileID)
	if err != nil { return "", "", err }
	audioURL := "https://tts.manpasik.com/audio/" + uuid.New().String() + ".wav"
	translatedText := "[" + targetLang + "] " + text
	_ = p
	return audioURL, translatedText, nil
}
GOEOF

cat > "$SVC/internal/repository/memory/voice_profile.go" << 'GOEOF'
package memory

import (
	"context"
	"sync"
	"github.com/manpasik/backend/services/voice-profile-service/internal/service"
)

type VoiceProfileRepository struct {
	mu   sync.RWMutex
	data map[string]*service.VoiceProfile
}

func NewVoiceProfileRepository() *VoiceProfileRepository {
	return &VoiceProfileRepository{data: make(map[string]*service.VoiceProfile)}
}
func (r *VoiceProfileRepository) Create(_ context.Context, p *service.VoiceProfile) error {
	r.mu.Lock(); defer r.mu.Unlock()
	r.data[p.ProfileID] = p; return nil
}
func (r *VoiceProfileRepository) GetByID(_ context.Context, id string) (*service.VoiceProfile, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	if v, ok := r.data[id]; ok { return v, nil }
	return nil, service.ErrNotFound
}
func (r *VoiceProfileRepository) ListByUser(_ context.Context, userID string) ([]*service.VoiceProfile, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	out := make([]*service.VoiceProfile, 0)
	for _, v := range r.data { if v.UserID == userID { out = append(out, v) } }
	return out, nil
}
func (r *VoiceProfileRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock(); defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok { return service.ErrNotFound }
	delete(r.data, id); return nil
}
GOEOF

cat > "$SVC/internal/handler/grpc.go" << 'GOEOF'
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
	p, err := h.svc.Create(ctx, req.GetUserId(), req.GetProfileName(), req.GetLanguage(), req.GetVoiceSampleUrl())
	if err != nil { return nil, toGRPC(err) }
	return profileToProto(p), nil
}
func (h *VoiceProfileHandler) GetVoiceProfile(ctx context.Context, req *v1.GetVoiceProfileRequest) (*v1.VoiceProfile, error) {
	p, err := h.svc.Get(ctx, req.GetProfileId())
	if err != nil { return nil, toGRPC(err) }
	return profileToProto(p), nil
}
func (h *VoiceProfileHandler) ListVoiceProfiles(ctx context.Context, req *v1.ListVoiceProfilesRequest) (*v1.ListVoiceProfilesResponse, error) {
	list, err := h.svc.List(ctx, req.GetUserId())
	if err != nil { return nil, toGRPC(err) }
	out := make([]*v1.VoiceProfile, 0, len(list))
	for _, p := range list { out = append(out, profileToProto(p)) }
	return &v1.ListVoiceProfilesResponse{Profiles: out, Total: int32(len(out))}, nil
}
func (h *VoiceProfileHandler) SynthesizeTranslation(ctx context.Context, req *v1.SynthesizeTranslationRequest) (*v1.SynthesizeTranslationResponse, error) {
	audioURL, translatedText, err := h.svc.SynthesizeTranslation(ctx, req.GetProfileId(), req.GetText(), req.GetTargetLanguage())
	if err != nil { return nil, toGRPC(err) }
	return &v1.SynthesizeTranslationResponse{AudioUrl: audioURL, TranslatedText: translatedText}, nil
}
func (h *VoiceProfileHandler) DeleteVoiceProfile(ctx context.Context, req *v1.DeleteVoiceProfileRequest) (*v1.DeleteVoiceProfileResponse, error) {
	if err := h.svc.Delete(ctx, req.GetProfileId()); err != nil { return nil, toGRPC(err) }
	return &v1.DeleteVoiceProfileResponse{Success: true}, nil
}

func profileToProto(p *service.VoiceProfile) *v1.VoiceProfile {
	return &v1.VoiceProfile{ProfileId: p.ProfileID, UserId: p.UserID, ProfileName: p.ProfileName,
		Language: p.Language, Status: p.Status, ModelUrl: p.ModelURL, CreatedAt: timestamppb.New(p.CreatedAt)}
}
func toGRPC(err error) error {
	switch err {
	case service.ErrNotFound: return status.Error(codes.NotFound, err.Error())
	case service.ErrInvalidInput: return status.Error(codes.InvalidArgument, err.Error())
	default: return status.Error(codes.Internal, err.Error())
	}
}
GOEOF

cat > "$SVC/internal/service/voice_profile_test.go" << 'GOEOF'
package service_test

import (
	"context"
	"testing"
	"github.com/manpasik/backend/services/voice-profile-service/internal/repository/memory"
	"github.com/manpasik/backend/services/voice-profile-service/internal/service"
	"go.uber.org/zap"
)

func newSvc() *service.VoiceProfileService {
	return service.NewVoiceProfileService(zap.NewNop(), memory.NewVoiceProfileRepository())
}
func TestCreateAndGet(t *testing.T) {
	s := newSvc(); ctx := context.Background()
	p, err := s.Create(ctx, "u1", "내 목소리", "ko", "")
	if err != nil { t.Fatal(err) }
	got, err := s.Get(ctx, p.ProfileID)
	if err != nil { t.Fatal(err) }
	if got.ProfileName != "내 목소리" { t.Error("name mismatch") }
}
func TestListProfiles(t *testing.T) {
	s := newSvc(); ctx := context.Background()
	s.Create(ctx, "u1", "A", "ko", "")
	s.Create(ctx, "u1", "B", "en", "")
	s.Create(ctx, "u2", "C", "ko", "")
	list, _ := s.List(ctx, "u1")
	if len(list) != 2 { t.Errorf("expected 2, got %d", len(list)) }
}
func TestDeleteProfile(t *testing.T) {
	s := newSvc(); ctx := context.Background()
	p, _ := s.Create(ctx, "u1", "X", "ko", "")
	if err := s.Delete(ctx, p.ProfileID); err != nil { t.Fatal(err) }
	if _, err := s.Get(ctx, p.ProfileID); err == nil { t.Error("expected error") }
}
func TestSynthesize(t *testing.T) {
	s := newSvc(); ctx := context.Background()
	p, _ := s.Create(ctx, "u1", "X", "ko", "")
	url, text, err := s.SynthesizeTranslation(ctx, p.ProfileID, "안녕하세요", "en")
	if err != nil { t.Fatal(err) }
	if url == "" || text == "" { t.Error("empty result") }
}
GOEOF

cat > "$SVC/Dockerfile" << 'GOEOF'
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /voice-profile-service ./services/voice-profile-service/cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /voice-profile-service /voice-profile-service
EXPOSE 50081
ENTRYPOINT ["/voice-profile-service"]
GOEOF

echo "[OK] voice-profile-service created"

echo "=== Part 1 done (concept + voice-profile) ==="
