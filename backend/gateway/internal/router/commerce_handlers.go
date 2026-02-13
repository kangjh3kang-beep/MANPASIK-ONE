package router

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
)

// ============================================================================
// Subscription Service Handlers
// ============================================================================

// GET /api/v1/subscriptions/{userId}
func (r *Router) handleGetSubscription(w http.ResponseWriter, req *http.Request) {
	userID := req.PathValue("userId")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user_id가 필요합니다")
		return
	}

	conn, err := dialGRPC(r.subscriptionAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "구독 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewSubscriptionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetSubscription(ctx, &v1.GetSubscriptionDetailRequest{
		UserId: userID,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, subscriptionDetailToMap(resp))
}

// POST /api/v1/subscriptions
func (r *Router) handleCreateSubscription(w http.ResponseWriter, req *http.Request) {
	var body struct {
		UserID string `json:"user_id"`
		Tier   int32  `json:"tier"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := dialGRPC(r.subscriptionAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "구독 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewSubscriptionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.CreateSubscription(ctx, &v1.CreateSubscriptionRequest{
		UserId: body.UserID,
		Tier:   v1.SubscriptionTier(body.Tier),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, subscriptionDetailToMap(resp))
}

// DELETE /api/v1/subscriptions/{subscriptionId}
func (r *Router) handleCancelSubscription(w http.ResponseWriter, req *http.Request) {
	subscriptionID := req.PathValue("subscriptionId")
	if subscriptionID == "" {
		writeError(w, http.StatusBadRequest, "subscription_id가 필요합니다")
		return
	}

	var body struct {
		UserID string `json:"user_id"`
		Reason string `json:"reason"`
	}
	// Body is optional for DELETE
	json.NewDecoder(req.Body).Decode(&body)

	conn, err := dialGRPC(r.subscriptionAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "구독 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewSubscriptionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.CancelSubscription(ctx, &v1.CancelSubscriptionRequest{
		UserId: body.UserID,
		Reason: body.Reason,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result := map[string]interface{}{
		"success": resp.GetSuccess(),
	}
	if resp.GetCancelledAt() != nil {
		result["cancelled_at"] = resp.GetCancelledAt().AsTime().Format(time.RFC3339)
	}
	if resp.GetEffectiveUntil() != nil {
		result["effective_until"] = resp.GetEffectiveUntil().AsTime().Format(time.RFC3339)
	}

	writeJSON(w, http.StatusOK, result)
}

// GET /api/v1/subscriptions/plans
func (r *Router) handleListSubscriptionPlans(w http.ResponseWriter, req *http.Request) {
	conn, err := dialGRPC(r.subscriptionAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "구독 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewSubscriptionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.ListSubscriptionPlans(ctx, &v1.ListSubscriptionPlansRequest{})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	plans := make([]map[string]interface{}, 0, len(resp.GetPlans()))
	for _, p := range resp.GetPlans() {
		plans = append(plans, map[string]interface{}{
			"tier":                 p.GetTier().String(),
			"name":                 p.GetName(),
			"description":          p.GetDescription(),
			"monthly_price_krw":    p.GetMonthlyPriceKrw(),
			"max_devices":          p.GetMaxDevices(),
			"max_family_members":   p.GetMaxFamilyMembers(),
			"ai_coaching_enabled":  p.GetAiCoachingEnabled(),
			"telemedicine_enabled": p.GetTelemedicineEnabled(),
			"features":             p.GetFeatures(),
		})
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"plans": plans,
	})
}

func subscriptionDetailToMap(s *v1.SubscriptionDetail) map[string]interface{} {
	m := map[string]interface{}{
		"subscription_id":      s.GetSubscriptionId(),
		"user_id":              s.GetUserId(),
		"tier":                 s.GetTier().String(),
		"status":               s.GetStatus().String(),
		"max_devices":          s.GetMaxDevices(),
		"max_family_members":   s.GetMaxFamilyMembers(),
		"ai_coaching_enabled":  s.GetAiCoachingEnabled(),
		"telemedicine_enabled": s.GetTelemedicineEnabled(),
		"monthly_price_krw":    s.GetMonthlyPriceKrw(),
		"auto_renew":           s.GetAutoRenew(),
	}
	if s.GetStartedAt() != nil {
		m["started_at"] = s.GetStartedAt().AsTime().Format(time.RFC3339)
	}
	if s.GetExpiresAt() != nil {
		m["expires_at"] = s.GetExpiresAt().AsTime().Format(time.RFC3339)
	}
	if s.GetCancelledAt() != nil {
		m["cancelled_at"] = s.GetCancelledAt().AsTime().Format(time.RFC3339)
	}
	return m
}

// ============================================================================
// Shop Service Handlers
// ============================================================================

// GET /api/v1/products?category=...&limit=...&offset=...
func (r *Router) handleListProducts(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	category, _ := strconv.Atoi(q.Get("category"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	if limit <= 0 {
		limit = 20
	}

	conn, err := dialGRPC(r.shopAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "쇼핑 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewShopServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.ListProducts(ctx, &v1.ListProductsRequest{
		Category: v1.ProductCategory(category),
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	products := make([]map[string]interface{}, 0, len(resp.GetProducts()))
	for _, p := range resp.GetProducts() {
		products = append(products, productToMap(p))
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"products":    products,
		"total_count": resp.GetTotalCount(),
	})
}

// GET /api/v1/products/{productId}
func (r *Router) handleGetProduct(w http.ResponseWriter, req *http.Request) {
	productID := req.PathValue("productId")
	if productID == "" {
		writeError(w, http.StatusBadRequest, "product_id가 필요합니다")
		return
	}

	conn, err := dialGRPC(r.shopAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "쇼핑 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewShopServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetProduct(ctx, &v1.GetProductRequest{
		ProductId: productID,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, productToMap(resp))
}

// POST /api/v1/cart
func (r *Router) handleAddToCart(w http.ResponseWriter, req *http.Request) {
	var body struct {
		UserID    string `json:"user_id"`
		ProductID string `json:"product_id"`
		Quantity  int32  `json:"quantity"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := dialGRPC(r.shopAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "쇼핑 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewShopServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.AddToCart(ctx, &v1.AddToCartRequest{
		UserId:    body.UserID,
		ProductId: body.ProductID,
		Quantity:  body.Quantity,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, cartToMap(resp))
}

// GET /api/v1/cart/{userId}
func (r *Router) handleGetCart(w http.ResponseWriter, req *http.Request) {
	userID := req.PathValue("userId")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user_id가 필요합니다")
		return
	}

	conn, err := dialGRPC(r.shopAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "쇼핑 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewShopServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetCart(ctx, &v1.GetCartRequest{
		UserId: userID,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, cartToMap(resp))
}

// POST /api/v1/orders
func (r *Router) handleCreateOrder(w http.ResponseWriter, req *http.Request) {
	var body struct {
		UserID          string `json:"user_id"`
		ShippingAddress string `json:"shipping_address"`
		PaymentMethod   string `json:"payment_method"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := dialGRPC(r.shopAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "쇼핑 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewShopServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.CreateOrder(ctx, &v1.CreateOrderRequest{
		UserId:          body.UserID,
		ShippingAddress: body.ShippingAddress,
		PaymentMethod:   body.PaymentMethod,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, orderToMap(resp))
}

// GET /api/v1/orders?user_id=...&limit=...&offset=...
func (r *Router) handleListOrders(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	userID := q.Get("user_id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user_id 쿼리 파라미터가 필요합니다")
		return
	}

	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	if limit <= 0 {
		limit = 20
	}

	conn, err := dialGRPC(r.shopAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "쇼핑 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewShopServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.ListOrders(ctx, &v1.ListOrdersRequest{
		UserId: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	orders := make([]map[string]interface{}, 0, len(resp.GetOrders()))
	for _, o := range resp.GetOrders() {
		orders = append(orders, orderToMap(o))
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"orders":      orders,
		"total_count": resp.GetTotalCount(),
	})
}

func productToMap(p *v1.Product) map[string]interface{} {
	m := map[string]interface{}{
		"product_id":  p.GetProductId(),
		"name":        p.GetName(),
		"description": p.GetDescription(),
		"category":    p.GetCategory().String(),
		"price_krw":   p.GetPriceKrw(),
		"stock":       p.GetStock(),
		"image_url":   p.GetImageUrl(),
		"is_active":   p.GetIsActive(),
	}
	if p.GetCreatedAt() != nil {
		m["created_at"] = p.GetCreatedAt().AsTime().Format(time.RFC3339)
	}
	return m
}

func cartToMap(c *v1.Cart) map[string]interface{} {
	items := make([]map[string]interface{}, 0, len(c.GetItems()))
	for _, item := range c.GetItems() {
		items = append(items, map[string]interface{}{
			"cart_item_id":   item.GetCartItemId(),
			"product_id":    item.GetProductId(),
			"product_name":  item.GetProductName(),
			"quantity":      item.GetQuantity(),
			"unit_price_krw": item.GetUnitPriceKrw(),
			"total_price_krw": item.GetTotalPriceKrw(),
		})
	}
	return map[string]interface{}{
		"user_id":         c.GetUserId(),
		"items":           items,
		"total_price_krw": c.GetTotalPriceKrw(),
	}
}

func orderToMap(o *v1.Order) map[string]interface{} {
	items := make([]map[string]interface{}, 0, len(o.GetItems()))
	for _, item := range o.GetItems() {
		items = append(items, map[string]interface{}{
			"product_id":     item.GetProductId(),
			"product_name":   item.GetProductName(),
			"quantity":       item.GetQuantity(),
			"unit_price_krw": item.GetUnitPriceKrw(),
			"total_price_krw": item.GetTotalPriceKrw(),
		})
	}
	m := map[string]interface{}{
		"order_id":         o.GetOrderId(),
		"user_id":          o.GetUserId(),
		"items":            items,
		"total_price_krw":  o.GetTotalPriceKrw(),
		"status":           o.GetStatus().String(),
		"shipping_address": o.GetShippingAddress(),
		"payment_id":       o.GetPaymentId(),
	}
	if o.GetCreatedAt() != nil {
		m["created_at"] = o.GetCreatedAt().AsTime().Format(time.RFC3339)
	}
	if o.GetUpdatedAt() != nil {
		m["updated_at"] = o.GetUpdatedAt().AsTime().Format(time.RFC3339)
	}
	return m
}

// ============================================================================
// Payment Service Handlers
// ============================================================================

// POST /api/v1/payments
func (r *Router) handleCreatePayment(w http.ResponseWriter, req *http.Request) {
	var body struct {
		UserID         string `json:"user_id"`
		OrderID        string `json:"order_id"`
		SubscriptionID string `json:"subscription_id"`
		PaymentType    int32  `json:"payment_type"`
		AmountKRW      int32  `json:"amount_krw"`
		PaymentMethod  string `json:"payment_method"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := dialGRPC(r.paymentAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "결제 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewPaymentServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.CreatePayment(ctx, &v1.CreatePaymentRequest{
		UserId:         body.UserID,
		OrderId:        body.OrderID,
		SubscriptionId: body.SubscriptionID,
		PaymentType:    v1.PaymentType(body.PaymentType),
		AmountKrw:      body.AmountKRW,
		PaymentMethod:  body.PaymentMethod,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, paymentDetailToMap(resp))
}

// POST /api/v1/payments/{paymentId}/confirm
func (r *Router) handleConfirmPayment(w http.ResponseWriter, req *http.Request) {
	paymentID := req.PathValue("paymentId")
	if paymentID == "" {
		writeError(w, http.StatusBadRequest, "payment_id가 필요합니다")
		return
	}

	var body struct {
		PgTransactionID string `json:"pg_transaction_id"`
		PgProvider      string `json:"pg_provider"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := dialGRPC(r.paymentAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "결제 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewPaymentServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.ConfirmPayment(ctx, &v1.ConfirmPaymentRequest{
		PaymentId:       paymentID,
		PgTransactionId: body.PgTransactionID,
		PgProvider:      body.PgProvider,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, paymentDetailToMap(resp))
}

// GET /api/v1/payments/{paymentId}
func (r *Router) handleGetPayment(w http.ResponseWriter, req *http.Request) {
	paymentID := req.PathValue("paymentId")
	if paymentID == "" {
		writeError(w, http.StatusBadRequest, "payment_id가 필요합니다")
		return
	}

	conn, err := dialGRPC(r.paymentAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "결제 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewPaymentServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetPayment(ctx, &v1.GetPaymentRequest{
		PaymentId: paymentID,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, paymentDetailToMap(resp))
}

func paymentDetailToMap(p *v1.PaymentDetail) map[string]interface{} {
	m := map[string]interface{}{
		"payment_id":        p.GetPaymentId(),
		"user_id":           p.GetUserId(),
		"order_id":          p.GetOrderId(),
		"subscription_id":   p.GetSubscriptionId(),
		"payment_type":      p.GetPaymentType().String(),
		"amount_krw":        p.GetAmountKrw(),
		"status":            p.GetStatus().String(),
		"payment_method":    p.GetPaymentMethod(),
		"pg_transaction_id": p.GetPgTransactionId(),
		"pg_provider":       p.GetPgProvider(),
	}
	if p.GetCreatedAt() != nil {
		m["created_at"] = p.GetCreatedAt().AsTime().Format(time.RFC3339)
	}
	if p.GetCompletedAt() != nil {
		m["completed_at"] = p.GetCompletedAt().AsTime().Format(time.RFC3339)
	}
	return m
}
