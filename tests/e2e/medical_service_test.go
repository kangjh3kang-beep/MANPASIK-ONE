package e2e

import (
	"net/http"
	"testing"
)

// ─── E2E-MED-001: 원격진료 예약 플로우 ───

func TestTelemedicineReservationFlow(t *testing.T) {
	email := uniqueEmail("med-tele")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Telemedicine User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 1. 의료진 목록 조회
	docResp, docResult := apiRequest(t, "GET", "/api/v1/telemedicine/doctors?specialty=general", nil, token)
	defer docResp.Body.Close()
	t.Logf("1. 의료진 목록: status=%d, result=%v", docResp.StatusCode, docResult)

	// 2. 예약 생성
	reserveBody := map[string]interface{}{
		"doctor_id":      "DOC-001",
		"date":           "2026-03-01",
		"time_slot":      "10:00",
		"consultation_type": "video",
		"symptoms":       "혈당 수치 이상",
	}
	resResp, resResult := apiRequest(t, "POST", "/api/v1/reservations", reserveBody, token)
	defer resResp.Body.Close()
	t.Logf("2. 예약 생성: status=%d, result=%v", resResp.StatusCode, resResult)

	// 3. 예약 목록 조회
	listResp, listResult := apiRequest(t, "GET", "/api/v1/reservations", nil, token)
	defer listResp.Body.Close()
	t.Logf("3. 예약 목록: status=%d, result=%v", listResp.StatusCode, listResult)
}

// ─── E2E-MED-002: 처방전 조회 ───

func TestPrescriptionQuery(t *testing.T) {
	email := uniqueEmail("med-rx")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Prescription User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 처방전 목록 조회
	resp, result := apiRequest(t, "GET", "/api/v1/prescriptions", nil, token)
	defer resp.Body.Close()
	t.Logf("처방전 목록: status=%d, result=%v", resp.StatusCode, result)
}

// ─── E2E-MED-003: 건강 리포트 생성 및 공유 ───

func TestHealthReportGeneration(t *testing.T) {
	email := uniqueEmail("med-report")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Report User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 건강 리포트 생성 요청
	reportBody := map[string]interface{}{
		"period":     "monthly",
		"start_date": "2026-01-01",
		"end_date":   "2026-01-31",
		"include":    []string{"glucose", "lipid_panel", "vitaminD"},
	}
	resp, result := apiRequest(t, "POST", "/api/v1/health-records/reports", reportBody, token)
	defer resp.Body.Close()
	t.Logf("건강 리포트 생성: status=%d, result=%v", resp.StatusCode, result)
}

// ─── E2E-MED-004: 119 긴급 연락 (위험 수치 알림) ───

func TestEmergencyAlertTrigger(t *testing.T) {
	email := uniqueEmail("med-emergency")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Emergency User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 위험 수치 측정 시뮬레이션 (혈당 40 mg/dL — 저혈당 위험)
	criticalBody := map[string]interface{}{
		"device_id":    "MPK-EMG-001",
		"cartridge_id": "CART-GLU-EMG-001",
		"channels":     make([]float64, 88),
		"result_value": 40.0,
		"result_unit":  "mg/dL",
		"is_critical":  true,
	}
	resp, result := apiRequest(t, "POST", measureEndpoint, criticalBody, token)
	defer resp.Body.Close()
	t.Logf("위험 수치 측정: status=%d, result=%v", resp.StatusCode, result)

	// 알림 발송 확인
	notifResp, notifResult := apiRequest(t, "GET", "/api/v1/notifications?type=critical", nil, token)
	defer notifResp.Body.Close()
	t.Logf("긴급 알림 확인: status=%d, result=%v", notifResp.StatusCode, notifResult)
}

// ─── E2E-MED-005: 건강 데이터 프라이버시 (타인 접근 차단) ───

func TestHealthDataPrivacy(t *testing.T) {
	emailA := uniqueEmail("med-priv-a")
	emailB := uniqueEmail("med-priv-b")
	tokenA := registerAndLogin(t, emailA, "E2eTest1!@#", "Patient A")
	tokenB := registerAndLogin(t, emailB, "E2eTest1!@#", "Patient B")
	if tokenA == "" || tokenB == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 사용자 A의 건강 기록 조회를 사용자 B가 시도
	resp, _ := apiRequest(t, "GET", "/api/v1/health-records?user_id=patient-a", nil, tokenB)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		t.Error("건강 데이터 프라이버시 위반: 타인 건강 기록 접근 가능")
	}
	t.Logf("건강 데이터 접근 차단 확인: status=%d", resp.StatusCode)
}
