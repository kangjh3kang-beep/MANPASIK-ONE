package e2e

import (
	"fmt"
	"net/http"
	"testing"
)

// ─── E2E-MKT-001: 상품 → 장바구니 → 결제 → 완료 전체 플로우 ───

func TestMarketPurchaseFullFlow(t *testing.T) {
	email := uniqueEmail("mkt-full")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Market User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 1. 상품 목록 조회
	listResp, listResult := apiRequest(t, "GET", "/api/v1/market/products?category=cartridge&limit=10", nil, token)
	defer listResp.Body.Close()
	t.Logf("1. 상품 목록: status=%d, result=%v", listResp.StatusCode, listResult)

	// 2. 상품 상세 조회
	detailResp, detailResult := apiRequest(t, "GET", "/api/v1/market/products/PROD-GLU-001", nil, token)
	defer detailResp.Body.Close()
	t.Logf("2. 상품 상세: status=%d, result=%v", detailResp.StatusCode, detailResult)

	// 3. 장바구니 추가
	cartBody := map[string]interface{}{
		"product_id": "PROD-GLU-001",
		"quantity":   2,
	}
	addResp, addResult := apiRequest(t, "POST", "/api/v1/market/cart", cartBody, token)
	defer addResp.Body.Close()
	t.Logf("3. 장바구니 추가: status=%d, result=%v", addResp.StatusCode, addResult)

	// 4. 장바구니 조회
	cartResp, cartResult := apiRequest(t, "GET", "/api/v1/market/cart", nil, token)
	defer cartResp.Body.Close()
	t.Logf("4. 장바구니 조회: status=%d, result=%v", cartResp.StatusCode, cartResult)

	// 5. 결제 시작
	checkoutBody := map[string]interface{}{
		"payment_method": "card",
		"shipping_address": map[string]string{
			"name":    "홍길동",
			"phone":   "010-1234-5678",
			"address": "서울시 강남구 역삼로 123",
			"zip":     "06234",
		},
	}
	checkoutResp, checkoutResult := apiRequest(t, "POST", "/api/v1/market/checkout", checkoutBody, token)
	defer checkoutResp.Body.Close()
	t.Logf("5. 결제: status=%d, result=%v", checkoutResp.StatusCode, checkoutResult)

	// 6. 주문 이력 조회
	orderResp, orderResult := apiRequest(t, "GET", "/api/v1/market/orders?limit=5", nil, token)
	defer orderResp.Body.Close()
	t.Logf("6. 주문 이력: status=%d, result=%v", orderResp.StatusCode, orderResult)
}

// ─── E2E-MKT-002: 장바구니 수량 변경/삭제 ───

func TestMarketCartManipulation(t *testing.T) {
	email := uniqueEmail("mkt-cart")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Cart User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 상품 2종 추가
	for i := 1; i <= 2; i++ {
		body := map[string]interface{}{
			"product_id": fmt.Sprintf("PROD-TEST-%03d", i),
			"quantity":   i,
		}
		resp, _ := apiRequest(t, "POST", "/api/v1/market/cart", body, token)
		resp.Body.Close()
	}

	// 수량 변경
	updateBody := map[string]interface{}{
		"quantity": 5,
	}
	updateResp, _ := apiRequest(t, "PATCH", "/api/v1/market/cart/PROD-TEST-001", updateBody, token)
	defer updateResp.Body.Close()
	t.Logf("수량 변경: status=%d", updateResp.StatusCode)

	// 상품 삭제
	delResp, _ := apiRequest(t, "DELETE", "/api/v1/market/cart/PROD-TEST-002", nil, token)
	defer delResp.Body.Close()
	t.Logf("상품 삭제: status=%d", delResp.StatusCode)

	// 장바구니 확인
	cartResp, cartResult := apiRequest(t, "GET", "/api/v1/market/cart", nil, token)
	defer cartResp.Body.Close()
	t.Logf("장바구니 확인: status=%d, result=%v", cartResp.StatusCode, cartResult)
}

// ─── E2E-MKT-003: 구독 플랜 변경 ───

func TestMarketSubscriptionChange(t *testing.T) {
	email := uniqueEmail("mkt-sub")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Subscription User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 1. 구독 플랜 목록 조회
	plansResp, plansResult := apiRequest(t, "GET", "/api/v1/market/subscription/plans", nil, token)
	defer plansResp.Body.Close()
	t.Logf("1. 플랜 목록: status=%d, result=%v", plansResp.StatusCode, plansResult)

	// 2. 현재 구독 조회
	currentResp, currentResult := apiRequest(t, "GET", "/api/v1/market/subscription", nil, token)
	defer currentResp.Body.Close()
	t.Logf("2. 현재 구독: status=%d, result=%v", currentResp.StatusCode, currentResult)

	// 3. 플랜 업그레이드
	upgradeBody := map[string]interface{}{
		"plan_id":        "premium",
		"payment_method": "card",
	}
	upgradeResp, upgradeResult := apiRequest(t, "POST", "/api/v1/market/subscription/upgrade", upgradeBody, token)
	defer upgradeResp.Body.Close()
	t.Logf("3. 업그레이드: status=%d, result=%v", upgradeResp.StatusCode, upgradeResult)
}

// ─── E2E-MKT-004: 주문 상세 및 배송 추적 ───

func TestMarketOrderTracking(t *testing.T) {
	email := uniqueEmail("mkt-track")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Track User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 주문 상세 조회
	detailResp, detailResult := apiRequest(t, "GET", "/api/v1/market/orders/ORD-001", nil, token)
	defer detailResp.Body.Close()
	t.Logf("주문 상세: status=%d, result=%v", detailResp.StatusCode, detailResult)

	// 배송 추적 조회
	trackResp, trackResult := apiRequest(t, "GET", "/api/v1/market/orders/ORD-001/tracking", nil, token)
	defer trackResp.Body.Close()
	t.Logf("배송 추적: status=%d, result=%v", trackResp.StatusCode, trackResult)
}

// ─── E2E-MKT-005: 인증 없는 결제 차단 ───

func TestMarketCheckoutWithoutAuth(t *testing.T) {
	checkoutBody := map[string]interface{}{
		"payment_method": "card",
	}
	resp, _ := apiRequest(t, "POST", "/api/v1/market/checkout", checkoutBody, "")
	defer resp.Body.Close()
	assertNotStatus(t, resp.StatusCode, http.StatusOK, "인증 없는 결제")
	assertNotStatus(t, resp.StatusCode, http.StatusCreated, "인증 없는 결제")
	t.Logf("인증 없는 결제 차단 확인: status=%d", resp.StatusCode)
}
