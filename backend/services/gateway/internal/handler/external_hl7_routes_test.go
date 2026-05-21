package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testORU = "MSH|^~\\&|LIS|HospA|EMR|HospA|20260521103000||ORU^R01|MSGTEST|P|2.5\r" +
	"PID|1||PAT-ext^^^MRN||김^철수||19800101|M\r" +
	"OBR|1||LAB-1|GLU^Glucose^L\r" +
	"OBX|1|NM|GLU^Glucose^LN||110|mg/dL|70-110|N|||F\r" +
	"OBX|2|NM|HBA1C^HbA1c^LN||5.4|%|<5.7|N|||F\r"

func newExternalHL7Request(t *testing.T, body, contentType, query string) *http.Request {
	t.Helper()
	url := "/api/v1/external/hl7"
	if query != "" {
		url += "?" + query
	}
	req := httptest.NewRequest(http.MethodPost, url, strings.NewReader(body))
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return req
}

func TestExternalHL7Import_Success(t *testing.T) {
	h := newNilHandler()
	mux := h.SetupRoutes()

	req := newExternalHL7Request(t, testORU, "text/plain", "")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if got, _ := resp["observation_count"].(float64); got != 2 {
		t.Errorf("observation_count = %v, 2 기대", resp["observation_count"])
	}
	if resp["bundle_id"] != "hl7v2-MSGTEST" {
		t.Errorf("bundle_id = %v", resp["bundle_id"])
	}
	if resp["fhir_version"] != "5.0.0" {
		t.Errorf("fhir_version = %v, 5.0.0 기대", resp["fhir_version"])
	}
	if w.Header().Get("X-Manpasik-Body-Bytes") == "" {
		t.Error("X-Manpasik-Body-Bytes 헤더 누락")
	}
}

func TestExternalHL7Import_R4VersionQuery(t *testing.T) {
	h := newNilHandler()
	mux := h.SetupRoutes()

	req := newExternalHL7Request(t, testORU, "text/plain", "version=R4")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["fhir_version"] != "4.0.1" {
		t.Errorf("fhir_version = %v, 4.0.1 기대", resp["fhir_version"])
	}
}

func TestExternalHL7Import_CustomPatientPrefix(t *testing.T) {
	h := newNilHandler()
	mux := h.SetupRoutes()

	req := newExternalHL7Request(t, testORU, "application/hl7-v2",
		"patient_prefix=ManpasikUser%2F")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "ManpasikUser/PAT-ext") {
		t.Errorf("patient_prefix 적용 안 됨:\n%s", body[:min(len(body), 500)])
	}
}

func TestExternalHL7Import_EmptyBody(t *testing.T) {
	h := newNilHandler()
	mux := h.SetupRoutes()

	req := newExternalHL7Request(t, "", "text/plain", "")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, 400 기대", w.Code)
	}
}

func TestExternalHL7Import_InvalidHL7(t *testing.T) {
	h := newNilHandler()
	mux := h.SetupRoutes()

	req := newExternalHL7Request(t, "PID|nothing", "text/plain", "")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, 422 기대", w.Code)
	}
	assertErrorContains(t, w.Body.Bytes(), "HL7")
}

func TestExternalHL7Import_UnsupportedContentType(t *testing.T) {
	h := newNilHandler()
	mux := h.SetupRoutes()

	req := newExternalHL7Request(t, testORU, "application/json", "")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusUnsupportedMediaType {
		t.Errorf("status = %d, 415 기대", w.Code)
	}
}

func TestExternalHL7Import_AllowedContentTypes(t *testing.T) {
	cases := []string{
		"text/plain",
		"text/plain; charset=utf-8",
		"application/hl7-v2",
		"application/x-hl7",
		"application/octet-stream",
		"", // 미지정도 허용
	}
	for _, ct := range cases {
		t.Run(ct, func(t *testing.T) {
			h := newNilHandler()
			mux := h.SetupRoutes()
			req := newExternalHL7Request(t, testORU, ct, "")
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			if w.Code != http.StatusOK {
				t.Errorf("ct=%q status=%d body=%s", ct, w.Code, w.Body.String())
			}
		})
	}
}

func TestExternalHL7Import_WrongMethod(t *testing.T) {
	h := newNilHandler()
	mux := h.SetupRoutes()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/external/hl7", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Go 1.22+ ServeMux 는 명시되지 않은 메서드에 대해 405 반환
	if w.Code != http.StatusMethodNotAllowed && w.Code != http.StatusNotFound {
		t.Errorf("status = %d, 405/404 기대", w.Code)
	}
}

func TestExternalHL7Import_NoObservationWarning(t *testing.T) {
	// OBX 가 없는 ADT 메시지 — Bundle 변환 자체는 성공하나 entry=0
	raw := "MSH|^~\\&|H|H|L|L|20260521||ADT^A01|MSG00099|P|2.5\r" +
		"PID|1||PATX\r"
	h := newNilHandler()
	mux := h.SetupRoutes()

	req := newExternalHL7Request(t, raw, "text/plain", "")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	if w.Header().Get("X-Manpasik-Warning") == "" {
		t.Error("X-Manpasik-Warning 헤더 누락")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
