package handler

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/manpasik/backend/shared/medical/fhir"
)

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

	// 응답: FHIR Bundle JSON + 메타 카운터 + Bundle ID
	resp := map[string]interface{}{
		"bundle":            bundle,
		"observation_count": countObservations(bundle),
		"bundle_id":         bundle.ID,
		"fhir_version":      string(responseVersion),
	}
	// Bundle.Entry 수가 0 인 메시지는 의미 있는 데이터가 없으므로 알림 헤더 추가
	if len(bundle.Entry) == 0 {
		w.Header().Set("X-Manpasik-Warning", "no observations extracted")
	}
	w.Header().Set("X-Manpasik-Body-Bytes", strconv.Itoa(len(raw)))
	writeJSON(w, http.StatusOK, resp)
}

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
