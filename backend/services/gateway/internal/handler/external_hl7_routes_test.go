package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/manpasik/backend/shared/medical/fhir"
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

	// PATCH 는 등록되지 않은 메서드 (POST=import, GET=list, DELETE=delete)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/external/hl7", nil)
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

// ============================================================================
// Phase AV — BundleStore 통합 테스트
// ============================================================================

func newHandlerWithStore() *RestHandler {
	h := newNilHandler()
	h.SetHL7BundleStore(fhir.NewMemoryBundleStore(0))
	return h
}

func TestExternalHL7Import_StoreTrue_SavesBundle(t *testing.T) {
	h := newHandlerWithStore()
	mux := h.SetupRoutes()

	req := newExternalHL7Request(t, testORU, "text/plain", "store=true")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["stored"] != true {
		t.Errorf("stored = %v, true 기대", resp["stored"])
	}

	// 저장 확인: GET 으로 다시 조회 가능
	getReq := httptest.NewRequest(http.MethodGet,
		"/api/v1/external/hl7/hl7v2-MSGTEST", nil)
	getW := httptest.NewRecorder()
	mux.ServeHTTP(getW, getReq)
	if getW.Code != http.StatusOK {
		t.Errorf("GET 조회 실패: status=%d body=%s", getW.Code, getW.Body.String())
	}
	var entry map[string]interface{}
	_ = json.Unmarshal(getW.Body.Bytes(), &entry)
	if entry["patient_id"] != "PAT-ext" {
		t.Errorf("patient_id = %v, PAT-ext 기대", entry["patient_id"])
	}
}

func TestExternalHL7Import_StoreFalse_DoesNotSave(t *testing.T) {
	h := newHandlerWithStore()
	mux := h.SetupRoutes()

	req := newExternalHL7Request(t, testORU, "text/plain", "")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["stored"] != false {
		t.Errorf("stored = %v, false 기대", resp["stored"])
	}

	// store 안 했으니 GET 조회 시 404
	getReq := httptest.NewRequest(http.MethodGet,
		"/api/v1/external/hl7/hl7v2-MSGTEST", nil)
	getW := httptest.NewRecorder()
	mux.ServeHTTP(getW, getReq)
	if getW.Code != http.StatusNotFound {
		t.Errorf("미저장 항목 GET status = %d, 404 기대", getW.Code)
	}
}

func TestExternalHL7Import_StoreWithoutStoreConfigured(t *testing.T) {
	h := newNilHandler() // store 미설정
	mux := h.SetupRoutes()

	req := newExternalHL7Request(t, testORU, "text/plain", "store=true")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("status = %d, 503 기대", w.Code)
	}
}

func TestExternalHL7Get_NotFound(t *testing.T) {
	h := newHandlerWithStore()
	mux := h.SetupRoutes()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/external/hl7/missing-id", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, 404 기대", w.Code)
	}
}

func TestExternalHL7List_Empty(t *testing.T) {
	h := newHandlerWithStore()
	mux := h.SetupRoutes()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/external/hl7", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if cnt, _ := resp["count"].(float64); cnt != 0 {
		t.Errorf("count = %v, 0 기대", resp["count"])
	}
}

func TestExternalHL7List_PatientFilter(t *testing.T) {
	h := newHandlerWithStore()
	mux := h.SetupRoutes()

	// 두 번 저장 (멱등 — 같은 MSG-ID 이므로 한 entry)
	for i := 0; i < 2; i++ {
		req := newExternalHL7Request(t, testORU, "text/plain", "store=true")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("save status = %d", w.Code)
		}
	}

	// PAT-ext 로 조회
	req := httptest.NewRequest(http.MethodGet,
		"/api/v1/external/hl7?patient=PAT-ext", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if cnt, _ := resp["count"].(float64); cnt != 1 {
		t.Errorf("count = %v, 1 기대 (멱등)", resp["count"])
	}

	// 미존재 환자
	missReq := httptest.NewRequest(http.MethodGet,
		"/api/v1/external/hl7?patient=NOTFOUND", nil)
	missW := httptest.NewRecorder()
	mux.ServeHTTP(missW, missReq)
	if missW.Code != http.StatusOK {
		t.Errorf("status = %d", missW.Code)
	}
	var missResp map[string]interface{}
	_ = json.Unmarshal(missW.Body.Bytes(), &missResp)
	if cnt, _ := missResp["count"].(float64); cnt != 0 {
		t.Errorf("미존재 patient count = %v, 0 기대", missResp["count"])
	}
}

func TestExternalHL7Delete_RoundTrip(t *testing.T) {
	h := newHandlerWithStore()
	mux := h.SetupRoutes()

	// 저장
	req := newExternalHL7Request(t, testORU, "text/plain", "store=true")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("save status = %d", w.Code)
	}

	// 삭제
	delReq := httptest.NewRequest(http.MethodDelete,
		"/api/v1/external/hl7/hl7v2-MSGTEST", nil)
	delW := httptest.NewRecorder()
	mux.ServeHTTP(delW, delReq)
	if delW.Code != http.StatusOK {
		t.Errorf("delete status = %d", delW.Code)
	}

	// 다시 GET 시 404
	getReq := httptest.NewRequest(http.MethodGet,
		"/api/v1/external/hl7/hl7v2-MSGTEST", nil)
	getW := httptest.NewRecorder()
	mux.ServeHTTP(getW, getReq)
	if getW.Code != http.StatusNotFound {
		t.Errorf("삭제 후 GET status = %d, 404 기대", getW.Code)
	}

	// 두 번째 삭제 시 404
	del2W := httptest.NewRecorder()
	mux.ServeHTTP(del2W, delReq)
	if del2W.Code != http.StatusNotFound {
		t.Errorf("이미 삭제된 항목 DELETE = %d, 404 기대", del2W.Code)
	}
}

func TestExternalHL7Delete_StoreNotConfigured(t *testing.T) {
	h := newNilHandler()
	mux := h.SetupRoutes()

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/external/hl7/x", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("status = %d, 503 기대", w.Code)
	}
}

func TestExternalHL7List_WithLimit(t *testing.T) {
	h := newHandlerWithStore()
	mux := h.SetupRoutes()

	// 서로 다른 MSG-ID 로 3개 저장
	for i := 0; i < 3; i++ {
		raw := strings.Replace(testORU, "MSGTEST", "MSG-"+string(rune('A'+i)), 1)
		req := newExternalHL7Request(t, raw, "text/plain", "store=true")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("save status = %d body=%s", w.Code, w.Body.String())
		}
	}

	req := httptest.NewRequest(http.MethodGet,
		"/api/v1/external/hl7?limit=2", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if cnt, _ := resp["count"].(float64); cnt != 2 {
		t.Errorf("count = %v, 2 기대 (limit=2)", resp["count"])
	}
}
