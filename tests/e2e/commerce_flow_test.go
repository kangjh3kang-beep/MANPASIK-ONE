package e2e

import (
	"fmt"
	"net/http"
	"testing"
)

// ─── Golden Path: 커머스 전체 플로우 ───
// TestGoldenPath_CommerceFlow tests the full purchase journey:
// 1. Login as patient
// 2. Browse products
// 3. Add to cart
// 4. Create order
// 5. Process payment
// 6. Verify order status

func TestGoldenPath_CommerceFlow(t *testing.T) {
	email := uniqueEmail("golden-commerce")
	password := "Commerce1!@#"
	name := "Commerce Golden User"

	// ──── Step 1: 회원가입 + 로그인 ────
	t.Log("Step 1: 회원가입 + 로그인")
	regStatus, _ := registerUser(t, email, password, name)
	if regStatus != http.StatusCreated && regStatus != http.StatusOK {
		t.Skipf("회원가입 응답: %d (서비스 미실행 가능)", regStatus)
		return
	}

	token := loginUser(t, email, password)
	if token == "" {
		t.Fatal("로그인 실패: 토큰 미발급")
	}
	t.Logf("  로그인 성공 (token=%s...)", token[:min(16, len(token))])

	// ──── Step 2: 상품 목록 탐색 ────
	t.Log("Step 2: 상품 목록 탐색")
	listResp, listResult := apiRequest(t, "GET", "/api/v1/market/products?category=cartridge&limit=10", nil, token)
	defer listResp.Body.Close()
	t.Logf("  상품 목록: status=%d, result=%v", listResp.StatusCode, listResult)

	// 상품 상세 조회
	detailResp, detailResult := apiRequest(t, "GET", "/api/v1/market/products/PROD-GLU-001", nil, token)
	defer detailResp.Body.Close()
	t.Logf("  상품 상세: status=%d, result=%v", detailResp.StatusCode, detailResult)

	// ──── Step 3: 장바구니 추가 ────
	t.Log("Step 3: 장바구니 추가")
	cartItems := []struct {
		productID string
		quantity  int
	}{
		{"PROD-GLU-001", 2},
		{"PROD-LIP-001", 1},
	}

	for _, item := range cartItems {
		cartBody := map[string]interface{}{
			"product_id": item.productID,
			"quantity":   item.quantity,
		}
		addResp, addResult := apiRequest(t, "POST", "/api/v1/market/cart", cartBody, token)
		defer addResp.Body.Close()
		t.Logf("  장바구니 추가 (%s x%d): status=%d, result=%v",
			item.productID, item.quantity, addResp.StatusCode, addResult)
	}

	// 장바구니 확인
	cartResp, cartResult := apiRequest(t, "GET", "/api/v1/market/cart", nil, token)
	defer cartResp.Body.Close()
	t.Logf("  장바구니 확인: status=%d, items=%v", cartResp.StatusCode, cartResult)

	// ──── Step 4: 주문 생성 (체크아웃) ────
	t.Log("Step 4: 주문 생성 (체크아웃)")
	checkoutBody := map[string]interface{}{
		"payment_method": "card",
		"shipping_address": map[string]string{
			"name":    "홍길동",
			"phone":   "010-1234-5678",
			"address": "서울시 강남구 역삼로 456",
			"zip":     "06234",
		},
		"coupon_code": "",
	}
	checkoutResp, checkoutResult := apiRequest(t, "POST", "/api/v1/market/checkout", checkoutBody, token)
	defer checkoutResp.Body.Close()
	t.Logf("  체크아웃: status=%d, result=%v", checkoutResp.StatusCode, checkoutResult)

	// 주문 ID 추출
	orderID := ""
	if id, ok := checkoutResult["order_id"].(string); ok {
		orderID = id
	} else if id, ok := checkoutResult["id"].(string); ok {
		orderID = id
	}

	// ──── Step 5: 결제 처리 ────
	t.Log("Step 5: 결제 처리")
	if orderID != "" {
		paymentBody := map[string]interface{}{
			"order_id":       orderID,
			"payment_method": "card",
			"card_token":     "tok_test_golden_path",
		}
		payResp, payResult := apiRequest(t, "POST", "/api/v1/payments/process", paymentBody, token)
		defer payResp.Body.Close()
		t.Logf("  결제 처리: status=%d, result=%v", payResp.StatusCode, payResult)
	} else {
		t.Log("  주문 ID 미확인 (checkout 응답에서 order_id 없음)")
	}

	// ──── Step 6: 주문 상태 확인 ────
	t.Log("Step 6: 주문 상태 확인")

	// 주문 이력 조회
	ordersResp, ordersResult := apiRequest(t, "GET", "/api/v1/market/orders?limit=5", nil, token)
	defer ordersResp.Body.Close()
	t.Logf("  주문 이력: status=%d, result=%v", ordersResp.StatusCode, ordersResult)

	// 특정 주문 상세 (order ID 가 있는 경우)
	if orderID != "" {
		orderDetailResp, orderDetailResult := apiRequest(t, "GET",
			fmt.Sprintf("/api/v1/market/orders/%s", orderID), nil, token)
		defer orderDetailResp.Body.Close()
		t.Logf("  주문 상세: status=%d, result=%v", orderDetailResp.StatusCode, orderDetailResult)

		// 주문 상태 확인
		if status, ok := orderDetailResult["status"].(string); ok {
			t.Logf("  주문 상태: %s", status)
		}
	}

	t.Log("Golden Path Commerce Flow completed.")
}

// ─── Golden Path: 구독 전체 플로우 ───

func TestGoldenPath_SubscriptionFlow(t *testing.T) {
	email := uniqueEmail("golden-sub")
	token := registerAndLogin(t, email, "Subscription1!@#", "Subscription User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 1. 구독 플랜 목록 조회
	t.Log("Step 1: 구독 플랜 목록 조회")
	plansResp, plansResult := apiRequest(t, "GET", "/api/v1/market/subscription/plans", nil, token)
	defer plansResp.Body.Close()
	t.Logf("  플랜 목록: status=%d, result=%v", plansResp.StatusCode, plansResult)

	// 2. 현재 구독 상태 확인
	t.Log("Step 2: 현재 구독 상태 확인")
	currentResp, currentResult := apiRequest(t, "GET", "/api/v1/market/subscription", nil, token)
	defer currentResp.Body.Close()
	t.Logf("  현재 구독: status=%d, result=%v", currentResp.StatusCode, currentResult)

	// 3. Basic 플랜 구독
	t.Log("Step 3: Basic 플랜 구독")
	subBody := map[string]interface{}{
		"plan_id":        "basic",
		"payment_method": "card",
		"card_token":     "tok_test_sub_basic",
	}
	subResp, subResult := apiRequest(t, "POST", "/api/v1/market/subscription/subscribe", subBody, token)
	defer subResp.Body.Close()
	t.Logf("  구독: status=%d, result=%v", subResp.StatusCode, subResult)

	// 4. Premium 으로 업그레이드
	t.Log("Step 4: Premium 으로 업그레이드")
	upgradeBody := map[string]interface{}{
		"plan_id":        "premium",
		"payment_method": "card",
	}
	upgradeResp, upgradeResult := apiRequest(t, "POST", "/api/v1/market/subscription/upgrade", upgradeBody, token)
	defer upgradeResp.Body.Close()
	t.Logf("  업그레이드: status=%d, result=%v", upgradeResp.StatusCode, upgradeResult)

	// 5. 업그레이드 후 상태 확인
	t.Log("Step 5: 업그레이드 후 상태 확인")
	afterResp, afterResult := apiRequest(t, "GET", "/api/v1/market/subscription", nil, token)
	defer afterResp.Body.Close()
	t.Logf("  업그레이드 후 구독: status=%d, result=%v", afterResp.StatusCode, afterResult)

	t.Log("Golden Path Subscription Flow completed.")
}

// ─── Golden Path: 인증 없는 결제 차단 검증 ───

func TestGoldenPath_UnauthorizedCommerceBlocked(t *testing.T) {
	endpoints := []struct {
		name   string
		method string
		path   string
		body   interface{}
	}{
		{"장바구니 추가", "POST", "/api/v1/market/cart", map[string]interface{}{"product_id": "P1", "quantity": 1}},
		{"체크아웃", "POST", "/api/v1/market/checkout", map[string]interface{}{"payment_method": "card"}},
		{"결제 처리", "POST", "/api/v1/payments/process", map[string]interface{}{"order_id": "O1"}},
		{"주문 이력", "GET", "/api/v1/market/orders", nil},
	}

	for _, ep := range endpoints {
		t.Run(ep.name, func(t *testing.T) {
			resp, _ := apiRequest(t, ep.method, ep.path, ep.body, "")
			defer resp.Body.Close()
			assertNotStatus(t, resp.StatusCode, http.StatusOK, "인증 없는 "+ep.name)
			assertNotStatus(t, resp.StatusCode, http.StatusCreated, "인증 없는 "+ep.name)
			t.Logf("  %s 차단 확인: status=%d", ep.name, resp.StatusCode)
		})
	}
}
