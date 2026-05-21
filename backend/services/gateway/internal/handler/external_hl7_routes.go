package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	hl7 "github.com/manpasik/backend/shared/hl7-parser"
	"github.com/manpasik/backend/shared/medical/fhir"
	"github.com/manpasik/backend/shared/tenancy"
)

// fhirBundleStore 는 외부 HL7 핸들러가 사용하는 BundleStore 최소 인터페이스.
// fhir.BundleStore 와 시그니처 동일 — gateway 가 직접 fhir 의 내부 구현에
// 의존하지 않고 인터페이스로만 통신.
type fhirBundleStore interface {
	Save(ctx context.Context, b *fhir.Bundle, tenantID, patientID string) error
	Get(ctx context.Context, bundleID, tenantID string) (*fhir.StoredBundle, error)
	ListByPatient(ctx context.Context, tenantID, patientID string, limit int) ([]*fhir.StoredBundle, error)
	ListByTenant(ctx context.Context, tenantID string, limit int) ([]*fhir.StoredBundle, error)
	Delete(ctx context.Context, bundleID, tenantID string) error
	Count(ctx context.Context, tenantID string) int
}

// SetHL7BundleStore 는 외부 HL7 통합용 BundleStore 주입 (선택).
// 미설정 시 ?store=true 요청 시 503 / 조회 엔드포인트도 503.
func (h *RestHandler) SetHL7BundleStore(s fhirBundleStore) {
	h.hl7BundleStore = s
}

// registerExternalHL7Routes 는 외부 LIS/EHR 시스템이 HL7 v2 메시지를 만파식에
// HTTP 로 송신할 수 있는 엔드포인트를 등록합니다 (Phase AU-2).
//
// 외부 시스템이 MLLP/TCP 가 아닌 HTTP 만 지원해도 통합 가능 — 의료기관 통합
// 게이트웨이 우회용 경로.
//
// 엔드포인트:
//
//	POST /api/v1/external/hl7
//	    Content-Type: text/plain (또는 application/hl7-v2)
//	    Body: raw HL7 v2.x 메시지 (CR/LF/CRLF segment 구분)
//	    Query: ?version=R4|R5 (기본 R5)
//	          ?patient_prefix=ManpasikUser/ (선택, Subject 참조 prefix)
//	    Response 200: FHIR Bundle JSON (type=collection)
//	    Response 4xx: 파싱/변환 실패 사유 JSON
//
// 인증: gateway 의 기본 인증 미들웨어로 보호 (mux 외부에서 적용). 본 핸들러
// 자체는 raw 본문 처리에만 집중합니다.
func (h *RestHandler) registerExternalHL7Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/external/hl7", h.handleExternalHL7Import)
	// Phase AV — 조회 엔드포인트
	mux.HandleFunc("GET /api/v1/external/hl7", h.handleExternalHL7List)
	mux.HandleFunc("GET /api/v1/external/hl7/{bundleId}", h.handleExternalHL7Get)
	mux.HandleFunc("DELETE /api/v1/external/hl7/{bundleId}", h.handleExternalHL7Delete)
}

// maxHL7BodyBytes 는 단일 메시지 본문 상한. 일반적인 HL7 v2 메시지는 수 KB
// 이지만 대량 검사 결과는 수십 KB 가능. 안전 마진 1 MB.
const maxHL7BodyBytes = 1 << 20

func (h *RestHandler) handleExternalHL7Import(w http.ResponseWriter, r *http.Request) {
	// Content-Type 검증 (text/plain | application/hl7-v2 | application/x-hl7)
	ct := strings.ToLower(strings.TrimSpace(r.Header.Get("Content-Type")))
	// "; charset=" 등 파라미터 제거
	if i := strings.Index(ct, ";"); i >= 0 {
		ct = strings.TrimSpace(ct[:i])
	}
	if ct != "" && !isHL7ContentType(ct) {
		writeError(w, http.StatusUnsupportedMediaType,
			"Content-Type 은 text/plain 또는 application/hl7-v2 권장")
		return
	}

	// 본문 읽기 (크기 제한 적용)
	r.Body = http.MaxBytesReader(w, r.Body, maxHL7BodyBytes)
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "본문 읽기 실패: "+err.Error())
		return
	}
	if len(raw) == 0 {
		writeError(w, http.StatusBadRequest, "본문 비어있음")
		return
	}

	version := parseHL7FhirVersion(r.URL.Query().Get("version"))
	// 응답에 명시할 FHIR 버전 — 미지정 시 기본 R5 (FromHL7v2Options 동일 정책)
	responseVersion := version
	if responseVersion == "" {
		responseVersion = fhir.VersionR5
	}

	opts := fhir.FromHL7v2Options{
		FhirVersion:             version,
		PatientResourceIDPrefix: r.URL.Query().Get("patient_prefix"),
	}

	bundle, err := fhir.ImportHL7v2(string(raw), opts)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, "HL7 변환 실패: "+err.Error())
		return
	}

	// Phase AV — ?store=true 시 Bundle 자동 저장
	stored := false
	storeRequested := r.URL.Query().Get("store") == "true"
	if storeRequested {
		if h.hl7BundleStore == nil {
			writeError(w, http.StatusServiceUnavailable, "BundleStore 미설정")
			return
		}
		// 환자 ID 추출 — 이미 변환 성공한 raw 이므로 재파싱 안전
		patientID := ""
		if msg, perr := hl7.Parse(string(raw)); perr == nil {
			if p := fhir.ExtractHL7Patient(msg); p != nil {
				patientID = p.ID
			}
		}
		tenantID := tenantFromHTTPRequest(r)
		if err := h.hl7BundleStore.Save(r.Context(), bundle, tenantID, patientID); err != nil {
			writeError(w, http.StatusInternalServerError, "Bundle 저장 실패: "+err.Error())
			return
		}
		stored = true
	}

	// 응답: FHIR Bundle JSON + 메타 카운터 + Bundle ID
	resp := map[string]interface{}{
		"bundle":            bundle,
		"observation_count": countObservations(bundle),
		"bundle_id":         bundle.ID,
		"fhir_version":      string(responseVersion),
		"stored":            stored,
	}
	// Bundle.Entry 수가 0 인 메시지는 의미 있는 데이터가 없으므로 알림 헤더 추가
	if len(bundle.Entry) == 0 {
		w.Header().Set("X-Manpasik-Warning", "no observations extracted")
	}
	w.Header().Set("X-Manpasik-Body-Bytes", strconv.Itoa(len(raw)))
	writeJSON(w, http.StatusOK, resp)
}

// tenantFromHTTPRequest 는 TenantPropagation 미들웨어가 ctx 에 주입한 tenant 추출.
// 미설정 시 "" 반환 (단일 테넌트 모드).
func tenantFromHTTPRequest(r *http.Request) string {
	if tid, ok := tenancy.TenantFromContext(r.Context()); ok {
		return string(tid)
	}
	return ""
}

// ============================================================================
// Phase AV — 저장된 Bundle 조회 / 삭제 엔드포인트
// ============================================================================

// handleExternalHL7Get 은 단일 Bundle 조회.
// GET /api/v1/external/hl7/{bundleId}
func (h *RestHandler) handleExternalHL7Get(w http.ResponseWriter, r *http.Request) {
	if h.hl7BundleStore == nil {
		writeError(w, http.StatusServiceUnavailable, "BundleStore 미설정")
		return
	}
	bundleID := r.PathValue("bundleId")
	if bundleID == "" {
		writeError(w, http.StatusBadRequest, "bundleId 필수")
		return
	}
	tenantID := tenantFromHTTPRequest(r)
	entry, err := h.hl7BundleStore.Get(r.Context(), bundleID, tenantID)
	if err != nil {
		if errors.Is(err, fhir.ErrBundleNotFound) {
			writeError(w, http.StatusNotFound, "Bundle 미존재")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, entry)
}

// handleExternalHL7List 는 tenant 또는 환자별 Bundle 목록.
// GET /api/v1/external/hl7
// GET /api/v1/external/hl7?patient=PAT123
// GET /api/v1/external/hl7?limit=20
func (h *RestHandler) handleExternalHL7List(w http.ResponseWriter, r *http.Request) {
	if h.hl7BundleStore == nil {
		writeError(w, http.StatusServiceUnavailable, "BundleStore 미설정")
		return
	}
	tenantID := tenantFromHTTPRequest(r)
	limit := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	patient := r.URL.Query().Get("patient")
	var (
		list []*fhir.StoredBundle
		err  error
	)
	if patient != "" {
		list, err = h.hl7BundleStore.ListByPatient(r.Context(), tenantID, patient, limit)
	} else {
		list, err = h.hl7BundleStore.ListByTenant(r.Context(), tenantID, limit)
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"count":  len(list),
		"items":  list,
	})
}

// handleExternalHL7Delete 는 단일 Bundle 제거.
// DELETE /api/v1/external/hl7/{bundleId}
func (h *RestHandler) handleExternalHL7Delete(w http.ResponseWriter, r *http.Request) {
	if h.hl7BundleStore == nil {
		writeError(w, http.StatusServiceUnavailable, "BundleStore 미설정")
		return
	}
	bundleID := r.PathValue("bundleId")
	if bundleID == "" {
		writeError(w, http.StatusBadRequest, "bundleId 필수")
		return
	}
	tenantID := tenantFromHTTPRequest(r)
	if err := h.hl7BundleStore.Delete(r.Context(), bundleID, tenantID); err != nil {
		if errors.Is(err, fhir.ErrBundleNotFound) {
			writeError(w, http.StatusNotFound, "Bundle 미존재")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"deleted": bundleID})
}

// _ 는 사용되지 않는 import 차단 방지 (context 가 hl7 패키지 사용 위해 필요).
var _ context.Context

// isHL7ContentType 는 HL7 v2 raw 메시지로 허용할 Content-Type 인지 검사.
func isHL7ContentType(ct string) bool {
	switch ct {
	case "text/plain",
		"application/hl7-v2",
		"application/x-hl7",
		"application/octet-stream":
		return true
	}
	return false
}

// parseHL7FhirVersion 은 query "version" 값을 fhir.Version 으로 변환.
//
// "R4"/"r4"/"4.0.1" → R4, "R5"/"r5"/"5.0.0" → R5, 기타 → 빈값 (기본 R5 사용됨).
func parseHL7FhirVersion(v string) fhir.Version {
	switch strings.ToUpper(strings.TrimSpace(v)) {
	case "R4", "4.0", "4.0.1", "4":
		return fhir.VersionR4
	case "R5", "5.0", "5.0.0", "5":
		return fhir.VersionR5
	}
	return ""
}

// countObservations 는 Bundle 에 포함된 Observation 리소스 개수.
func countObservations(b *fhir.Bundle) int {
	if b == nil {
		return 0
	}
	n := 0
	for _, e := range b.Entry {
		if e.Resource != nil && e.Resource.GetResourceType() == "Observation" {
			n++
		}
	}
	return n
}
