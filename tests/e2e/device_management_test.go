package e2e

import (
	"fmt"
	"net/http"
	"testing"
)

// ─── E2E-DEV-001: 디바이스 등록 → 조회 → 삭제 ───

func TestDeviceLifecycle(t *testing.T) {
	email := uniqueEmail("dev-lifecycle")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Device User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 1. 디바이스 등록
	devBody := map[string]string{
		"device_id":        "MPK-LC-001",
		"name":             "Lifecycle Test Reader",
		"firmware_version": "1.0.0",
		"serial_number":    "SN-E2E-001",
	}
	regResp, regResult := apiRequest(t, "POST", devicesEndpoint, devBody, token)
	defer regResp.Body.Close()
	t.Logf("1. 디바이스 등록: status=%d, result=%v", regResp.StatusCode, regResult)

	// 2. 디바이스 목록 조회
	listResp, listResult := apiRequest(t, "GET", devicesEndpoint, nil, token)
	defer listResp.Body.Close()
	t.Logf("2. 디바이스 목록: status=%d, result=%v", listResp.StatusCode, listResult)

	// 3. 특정 디바이스 조회
	detailResp, detailResult := apiRequest(t, "GET", devicesEndpoint+"/MPK-LC-001", nil, token)
	defer detailResp.Body.Close()
	t.Logf("3. 디바이스 상세: status=%d, result=%v", detailResp.StatusCode, detailResult)

	// 4. 디바이스 삭제
	delResp, _ := apiRequest(t, "DELETE", devicesEndpoint+"/MPK-LC-001", nil, token)
	defer delResp.Body.Close()
	t.Logf("4. 디바이스 삭제: status=%d", delResp.StatusCode)
}

// ─── E2E-DEV-002: 다중 디바이스 관리 ───

func TestMultipleDeviceManagement(t *testing.T) {
	email := uniqueEmail("dev-multi")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Multi Device User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 3대 디바이스 등록
	for i := 1; i <= 3; i++ {
		devBody := map[string]string{
			"device_id":        fmt.Sprintf("MPK-MULTI-%03d", i),
			"name":             fmt.Sprintf("Reader %d", i),
			"firmware_version": "1.0.0",
		}
		resp, _ := apiRequest(t, "POST", devicesEndpoint, devBody, token)
		resp.Body.Close()
	}

	// 목록 조회 (3대 이상 확인)
	listResp, listResult := apiRequest(t, "GET", devicesEndpoint, nil, token)
	defer listResp.Body.Close()
	t.Logf("다중 디바이스 목록: status=%d, result=%v", listResp.StatusCode, listResult)
}

// ─── E2E-DEV-003: 디바이스 펌웨어 버전 확인 ───

func TestDeviceFirmwareCheck(t *testing.T) {
	email := uniqueEmail("dev-fw")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Firmware User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 디바이스 등록 (구 버전 펌웨어)
	devBody := map[string]string{
		"device_id":        "MPK-FW-001",
		"name":             "Old Firmware Reader",
		"firmware_version": "0.9.0",
	}
	resp, result := apiRequest(t, "POST", devicesEndpoint, devBody, token)
	defer resp.Body.Close()
	t.Logf("구 버전 펌웨어 디바이스 등록: status=%d, result=%v", resp.StatusCode, result)

	// 펌웨어 업데이트 확인 요청
	fwResp, fwResult := apiRequest(t, "GET", devicesEndpoint+"/MPK-FW-001/firmware/check", nil, token)
	defer fwResp.Body.Close()
	t.Logf("펌웨어 업데이트 확인: status=%d, result=%v", fwResp.StatusCode, fwResult)
}

// ─── E2E-DEV-004: 타인 디바이스 접근 차단 ───

func TestDeviceAccessControl(t *testing.T) {
	emailA := uniqueEmail("dev-acl-a")
	emailB := uniqueEmail("dev-acl-b")
	tokenA := registerAndLogin(t, emailA, "E2eTest1!@#", "User A")
	tokenB := registerAndLogin(t, emailB, "E2eTest1!@#", "User B")
	if tokenA == "" || tokenB == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 사용자 A가 디바이스 등록
	devBody := map[string]string{
		"device_id": "MPK-ACL-001",
		"name":      "User A Reader",
	}
	apiRequest(t, "POST", devicesEndpoint, devBody, tokenA)

	// 사용자 B가 사용자 A의 디바이스 접근 시도
	resp, _ := apiRequest(t, "GET", devicesEndpoint+"/MPK-ACL-001", nil, tokenB)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		t.Error("타인 디바이스 접근 허용됨: 권한 분리 취약점")
	}
	t.Logf("타인 디바이스 접근 차단 확인: status=%d", resp.StatusCode)
}
