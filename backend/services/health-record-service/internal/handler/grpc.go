// Package handler는 health-record-service의 gRPC 핸들러입니다.
package handler

import (
	"context"
	"time"

	"github.com/manpasik/backend/services/health-record-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// HealthRecordHandler는 HealthRecordService gRPC 서버를 구현합니다.
type HealthRecordHandler struct {
	v1.UnimplementedHealthRecordServiceServer
	svc *service.HealthRecordService
	log *zap.Logger
}

// NewHealthRecordHandler는 HealthRecordHandler를 생성합니다.
func NewHealthRecordHandler(svc *service.HealthRecordService, log *zap.Logger) *HealthRecordHandler {
	return &HealthRecordHandler{svc: svc, log: log}
}

// CreateRecord는 건강 기록 생성 RPC입니다.
func (h *HealthRecordHandler) CreateRecord(ctx context.Context, req *v1.CreateHealthRecordRequest) (*v1.HealthRecord, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}
	if req.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "title은 필수입니다")
	}

	// Extract data_json and source from metadata map (proto doesn't have DataJson/Source fields)
	dataJson := ""
	source := ""
	if req.Metadata != nil {
		dataJson = req.Metadata["data_json"]
		source = req.Metadata["source"]
	}

	record, err := h.svc.CreateRecord(
		ctx,
		req.UserId,
		protoRecordTypeToService(req.RecordType),
		req.Title,
		req.Description,
		dataJson,
		source,
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	return recordToProto(record), nil
}

// GetRecord는 건강 기록 조회 RPC입니다.
func (h *HealthRecordHandler) GetRecord(ctx context.Context, req *v1.GetHealthRecordRequest) (*v1.HealthRecord, error) {
	if req == nil || req.RecordId == "" {
		return nil, status.Error(codes.InvalidArgument, "record_id는 필수입니다")
	}

	record, err := h.svc.GetRecord(ctx, req.RecordId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return recordToProto(record), nil
}

// ListRecords는 건강 기록 목록 조회 RPC입니다.
func (h *HealthRecordHandler) ListRecords(ctx context.Context, req *v1.ListHealthRecordsRequest) (*v1.ListHealthRecordsResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	records, total, err := h.svc.ListRecords(
		ctx,
		req.UserId,
		protoRecordTypeToService(req.TypeFilter),
		int(req.Limit),
		int(req.Offset),
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	var pbRecords []*v1.HealthRecord
	for _, r := range records {
		pbRecords = append(pbRecords, recordToProto(r))
	}

	return &v1.ListHealthRecordsResponse{
		Records:    pbRecords,
		TotalCount: int32(total),
	}, nil
}

// UpdateRecord는 건강 기록 업데이트 RPC입니다.
func (h *HealthRecordHandler) UpdateRecord(ctx context.Context, req *v1.UpdateHealthRecordRequest) (*v1.HealthRecord, error) {
	if req == nil || req.RecordId == "" {
		return nil, status.Error(codes.InvalidArgument, "record_id는 필수입니다")
	}

	// Extract data_json from metadata map (proto doesn't have DataJson field)
	dataJson := ""
	if req.Metadata != nil {
		dataJson = req.Metadata["data_json"]
	}

	record, err := h.svc.UpdateRecord(ctx, req.RecordId, req.Title, req.Description, dataJson)
	if err != nil {
		return nil, toGRPC(err)
	}

	return recordToProto(record), nil
}

// DeleteRecord는 건강 기록 삭제 RPC입니다.
func (h *HealthRecordHandler) DeleteRecord(ctx context.Context, req *v1.DeleteHealthRecordRequest) (*v1.DeleteHealthRecordResponse, error) {
	if req == nil || req.RecordId == "" {
		return nil, status.Error(codes.InvalidArgument, "record_id는 필수입니다")
	}

	err := h.svc.DeleteRecord(ctx, req.RecordId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.DeleteHealthRecordResponse{Success: true}, nil
}

// ExportToFHIR는 FHIR R4 내보내기 RPC입니다.
func (h *HealthRecordHandler) ExportToFHIR(ctx context.Context, req *v1.ExportToFHIRRequest) (*v1.ExportToFHIRResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	bundleJSON, count, _, err := h.svc.ExportToFHIR(ctx, req.UserId, nil, nil, nil)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.ExportToFHIRResponse{
		FhirBundleJson: bundleJSON,
		ResourceCount:  int32(count),
	}, nil
}

// ImportFromFHIR는 FHIR R4 가져오기 RPC입니다.
func (h *HealthRecordHandler) ImportFromFHIR(ctx context.Context, req *v1.ImportFromFHIRRequest) (*v1.ImportFromFHIRResponse, error) {
	if req == nil || req.UserId == "" || req.FhirBundleJson == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id, fhir_bundle_json은 필수입니다")
	}

	_, importedCount, skippedCount, errors, err := h.svc.ImportFromFHIR(ctx, req.UserId, req.FhirBundleJson)
	if err != nil && importedCount == 0 {
		return nil, toGRPC(err)
	}

	return &v1.ImportFromFHIRResponse{
		ImportedCount: int32(importedCount),
		SkippedCount:  int32(skippedCount),
		Errors:        errors,
	}, nil
}

// GetHealthSummary는 건강 요약 조회 RPC입니다.
func (h *HealthRecordHandler) GetHealthSummary(ctx context.Context, req *v1.GetHealthSummaryRequest) (*v1.GetHealthSummaryResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	total, recordsByType, recentRecords, _, err := h.svc.GetHealthSummary(ctx, req.UserId, 0)
	if err != nil {
		return nil, toGRPC(err)
	}

	// Convert map[string]int to map[string]int32
	pbRecordsByType := make(map[string]int32, len(recordsByType))
	for k, v := range recordsByType {
		pbRecordsByType[k] = int32(v)
	}

	// Convert service records to proto records
	var pbRecentRecords []*v1.HealthRecord
	for _, r := range recentRecords {
		pbRecentRecords = append(pbRecentRecords, recordToProto(r))
	}

	return &v1.GetHealthSummaryResponse{
		UserId:        req.UserId,
		TotalRecords:  int32(total),
		RecordsByType: pbRecordsByType,
		RecentRecords: pbRecentRecords,
	}, nil
}

// CreateDataSharingConsent는 데이터 공유 동의 생성 RPC입니다.
func (h *HealthRecordHandler) CreateDataSharingConsent(ctx context.Context, req *v1.CreateConsentRequest) (*v1.DataSharingConsent, error) {
	consent := &service.DataSharingConsent{
		UserID:       req.UserId,
		ProviderID:   req.ProviderId,
		ProviderName: req.ProviderName,
		ConsentType:  service.ConsentType(req.ConsentType),
		Scope:        req.Scope,
		Purpose:      req.Purpose,
	}
	result, err := h.svc.CreateDataSharingConsent(ctx, consent)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "동의 생성 실패: %v", err)
	}
	return consentToProto(result), nil
}

// RevokeDataSharingConsent는 데이터 공유 동의 철회 RPC입니다.
func (h *HealthRecordHandler) RevokeDataSharingConsent(ctx context.Context, req *v1.RevokeConsentRequest) (*v1.RevokeConsentResponse, error) {
	err := h.svc.RevokeDataSharingConsent(ctx, req.ConsentId, req.Reason)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "동의 철회 실패: %v", err)
	}
	return &v1.RevokeConsentResponse{Success: true}, nil
}

// ListDataSharingConsents는 데이터 공유 동의 목록 조회 RPC입니다.
func (h *HealthRecordHandler) ListDataSharingConsents(ctx context.Context, req *v1.ListConsentsRequest) (*v1.ListConsentsResponse, error) {
	consents, err := h.svc.ListDataSharingConsents(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "동의 목록 조회 실패: %v", err)
	}
	var pbConsents []*v1.DataSharingConsent
	for _, c := range consents {
		pbConsents = append(pbConsents, consentToProto(c))
	}
	return &v1.ListConsentsResponse{Consents: pbConsents}, nil
}

// ShareWithProvider는 의료 데이터를 제공자와 공유하는 RPC입니다.
func (h *HealthRecordHandler) ShareWithProvider(ctx context.Context, req *v1.ShareWithProviderRequest) (*v1.ShareWithProviderResponse, error) {
	bundle, err := h.svc.ShareWithProvider(ctx, req.ConsentId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "데이터 공유 실패: %v", err)
	}
	return &v1.ShareWithProviderResponse{
		FhirBundleJson: bundle.FHIRBundleJSON,
		ResourceCount:  int32(bundle.ResourceCount),
		SharedAt:       bundle.SharedAt.Format(time.RFC3339),
	}, nil
}

// GetDataAccessLog는 데이터 접근 로그 조회 RPC입니다.
func (h *HealthRecordHandler) GetDataAccessLog(ctx context.Context, req *v1.GetDataAccessLogRequest) (*v1.GetDataAccessLogResponse, error) {
	entries, total, err := h.svc.GetDataAccessLog(ctx, req.UserId, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "접근 로그 조회 실패: %v", err)
	}
	var pbEntries []*v1.DataAccessLogEntry
	for _, e := range entries {
		pbEntries = append(pbEntries, &v1.DataAccessLogEntry{
			Id:           e.ID,
			ConsentId:    e.ConsentID,
			UserId:       e.UserID,
			ProviderId:   e.ProviderID,
			Action:       e.Action,
			ResourceType: e.ResourceType,
			ResourceIds:  e.ResourceIDs,
			AccessedAt:   e.AccessedAt.Format(time.RFC3339),
		})
	}
	return &v1.GetDataAccessLogResponse{Entries: pbEntries, Total: int32(total)}, nil
}

// consentToProto는 service.DataSharingConsent를 v1.DataSharingConsent로 변환합니다.
func consentToProto(c *service.DataSharingConsent) *v1.DataSharingConsent {
	pb := &v1.DataSharingConsent{
		Id:           c.ID,
		UserId:       c.UserID,
		ProviderId:   c.ProviderID,
		ProviderName: c.ProviderName,
		ConsentType:  string(c.ConsentType),
		Scope:        c.Scope,
		Purpose:      c.Purpose,
		Status:       string(c.Status),
		GrantedAt:    c.GrantedAt.Format(time.RFC3339),
		ExpiresAt:    c.ExpiresAt.Format(time.RFC3339),
	}
	if !c.RevokedAt.IsZero() {
		pb.RevokedAt = c.RevokedAt.Format(time.RFC3339)
		pb.RevokeReason = c.RevokeReason
	}
	return pb
}

// ============================================================================
// 변환 헬퍼
// ============================================================================

func recordToProto(r *service.HealthRecord) *v1.HealthRecord {
	// Store Data and Source in the metadata map (proto uses metadata instead of dedicated fields)
	metadata := map[string]string{}
	if r.Data != "" {
		metadata["data_json"] = r.Data
	}
	if r.Source != "" {
		metadata["source"] = r.Source
	}

	return &v1.HealthRecord{
		RecordId:    r.ID,
		UserId:      r.UserID,
		RecordType:  serviceRecordTypeToProto(r.RecordType),
		Title:       r.Title,
		Description: r.Description,
		Metadata:    metadata,
		CreatedAt:   timestamppb.New(r.CreatedAt),
		UpdatedAt:   timestamppb.New(r.UpdatedAt),
	}
}

func protoRecordTypeToService(t v1.HealthRecordType) service.HealthRecordType {
	switch t {
	case v1.HealthRecordType_HEALTH_RECORD_TYPE_LAB_RESULT:
		return service.RecordTypeLabResult
	case v1.HealthRecordType_HEALTH_RECORD_TYPE_IMAGING:
		return service.RecordTypeImaging
	case v1.HealthRecordType_HEALTH_RECORD_TYPE_VITAL_SIGN:
		return service.RecordTypeVitalSign
	case v1.HealthRecordType_HEALTH_RECORD_TYPE_ALLERGY:
		return service.RecordTypeAllergy
	case v1.HealthRecordType_HEALTH_RECORD_TYPE_CONDITION:
		return service.RecordTypeCondition
	case v1.HealthRecordType_HEALTH_RECORD_TYPE_IMMUNIZATION:
		return service.RecordTypeImmunization
	case v1.HealthRecordType_HEALTH_RECORD_TYPE_PROCEDURE:
		return service.RecordTypeProcedure
	default:
		return service.RecordTypeUnknown
	}
}

func serviceRecordTypeToProto(t service.HealthRecordType) v1.HealthRecordType {
	switch t {
	case service.RecordTypeLabResult:
		return v1.HealthRecordType_HEALTH_RECORD_TYPE_LAB_RESULT
	case service.RecordTypeImaging:
		return v1.HealthRecordType_HEALTH_RECORD_TYPE_IMAGING
	case service.RecordTypeVitalSign:
		return v1.HealthRecordType_HEALTH_RECORD_TYPE_VITAL_SIGN
	case service.RecordTypeAllergy:
		return v1.HealthRecordType_HEALTH_RECORD_TYPE_ALLERGY
	case service.RecordTypeCondition:
		return v1.HealthRecordType_HEALTH_RECORD_TYPE_CONDITION
	case service.RecordTypeImmunization:
		return v1.HealthRecordType_HEALTH_RECORD_TYPE_IMMUNIZATION
	case service.RecordTypeProcedure:
		return v1.HealthRecordType_HEALTH_RECORD_TYPE_PROCEDURE
	default:
		return v1.HealthRecordType_HEALTH_RECORD_TYPE_UNKNOWN
	}
}

func serviceFHIRTypeToProto(t service.FHIRResourceType) v1.FHIRResourceType {
	switch t {
	case service.FHIRObservation:
		return v1.FHIRResourceType_FHIR_RESOURCE_TYPE_OBSERVATION
	case service.FHIRCondition:
		return v1.FHIRResourceType_FHIR_RESOURCE_TYPE_CONDITION
	case service.FHIRMedicationStatement:
		return v1.FHIRResourceType_FHIR_RESOURCE_TYPE_OBSERVATION // medication mapped to observation
	case service.FHIRAllergyIntolerance:
		return v1.FHIRResourceType_FHIR_RESOURCE_TYPE_ALLERGY_INTOLERANCE
	case service.FHIRImmunization:
		return v1.FHIRResourceType_FHIR_RESOURCE_TYPE_IMMUNIZATION
	case service.FHIRProcedure:
		return v1.FHIRResourceType_FHIR_RESOURCE_TYPE_PROCEDURE
	case service.FHIRDiagnosticReport:
		return v1.FHIRResourceType_FHIR_RESOURCE_TYPE_DIAGNOSTIC_REPORT
	case service.FHIRPatient:
		return v1.FHIRResourceType_FHIR_RESOURCE_TYPE_BUNDLE
	default:
		return v1.FHIRResourceType_FHIR_RESOURCE_TYPE_UNKNOWN
	}
}

func toGRPC(err error) error {
	if appErr, ok := err.(*apperrors.AppError); ok {
		return appErr.ToGRPC()
	}
	return status.Error(codes.Internal, err.Error())
}
