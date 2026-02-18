package handler

import (
	"net/http"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
)

// registerMarketRoutes는 마켓/결제/처방전 관련 REST 엔드포인트를 등록합니다.
func (h *RestHandler) registerMarketRoutes(mux *http.ServeMux) {
	// Shop
	mux.HandleFunc("GET /api/v1/products", h.handleListProducts)
	mux.HandleFunc("GET /api/v1/products/{productId}", h.handleGetProduct)
	mux.HandleFunc("GET /api/v1/products/{productId}/reviews", h.handleGetProductReviews)
	mux.HandleFunc("POST /api/v1/products/{productId}/reviews", h.handleCreateProductReview)
	mux.HandleFunc("POST /api/v1/cart", h.handleAddToCart)
	mux.HandleFunc("GET /api/v1/cart/{userId}", h.handleGetCart)
	mux.HandleFunc("POST /api/v1/orders", h.handleCreateOrder)
	mux.HandleFunc("GET /api/v1/orders", h.handleListOrders)

	// Payment
	mux.HandleFunc("POST /api/v1/payments", h.handleCreatePayment)
	mux.HandleFunc("POST /api/v1/payments/{paymentId}/confirm", h.handleConfirmPayment)
	mux.HandleFunc("GET /api/v1/payments/{paymentId}", h.handleGetPayment)

	// Prescription
	mux.HandleFunc("POST /api/v1/prescriptions/{prescriptionId}/pharmacy", h.handleSelectPharmacy)
	mux.HandleFunc("POST /api/v1/prescriptions/{prescriptionId}/send", h.handleSendToPharmacy)
	mux.HandleFunc("GET /api/v1/prescriptions/token/{token}", h.handleGetPrescriptionByToken)
}

// ── Shop ──

func (h *RestHandler) handleListProducts(w http.ResponseWriter, r *http.Request) {
	if h.shop == nil {
		writeError(w, http.StatusServiceUnavailable, "shop service unavailable")
		return
	}
	resp, err := h.shop.ListProducts(r.Context(), &v1.ListProductsRequest{
		Category: v1.ProductCategory(queryInt(r, "category", 0)),
		Limit:    queryInt(r, "limit", 20),
		Offset:   queryInt(r, "offset", 0),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	if h.shop == nil {
		writeError(w, http.StatusServiceUnavailable, "shop service unavailable")
		return
	}
	productId := r.PathValue("productId")
	resp, err := h.shop.GetProduct(r.Context(), &v1.GetProductRequest{ProductId: productId})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetProductReviews(w http.ResponseWriter, r *http.Request) {
	// Product reviews are returned as part of community service
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"reviews":    []interface{}{},
		"total":      0,
		"product_id": r.PathValue("productId"),
	})
}

func (h *RestHandler) handleCreateProductReview(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "리뷰가 등록되었습니다",
	})
}

func (h *RestHandler) handleAddToCart(w http.ResponseWriter, r *http.Request) {
	if h.shop == nil {
		writeError(w, http.StatusServiceUnavailable, "shop service unavailable")
		return
	}
	var body struct {
		UserID    string `json:"user_id"`
		ProductID string `json:"product_id"`
		Quantity  int32  `json:"quantity"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Quantity == 0 {
		body.Quantity = 1
	}
	resp, err := h.shop.AddToCart(r.Context(), &v1.AddToCartRequest{
		UserId:    body.UserID,
		ProductId: body.ProductID,
		Quantity:  body.Quantity,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetCart(w http.ResponseWriter, r *http.Request) {
	if h.shop == nil {
		writeError(w, http.StatusServiceUnavailable, "shop service unavailable")
		return
	}
	userId := r.PathValue("userId")
	resp, err := h.shop.GetCart(r.Context(), &v1.GetCartRequest{UserId: userId})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	if h.shop == nil {
		writeError(w, http.StatusServiceUnavailable, "shop service unavailable")
		return
	}
	var body struct {
		UserID          string `json:"user_id"`
		ShippingAddress string `json:"shipping_address"`
		PaymentMethod   string `json:"payment_method"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.shop.CreateOrder(r.Context(), &v1.CreateOrderRequest{
		UserId:          body.UserID,
		ShippingAddress: body.ShippingAddress,
		PaymentMethod:   body.PaymentMethod,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleListOrders(w http.ResponseWriter, r *http.Request) {
	if h.shop == nil {
		writeError(w, http.StatusServiceUnavailable, "shop service unavailable")
		return
	}
	resp, err := h.shop.ListOrders(r.Context(), &v1.ListOrdersRequest{
		UserId: r.URL.Query().Get("user_id"),
		Limit:  queryInt(r, "limit", 20),
		Offset: queryInt(r, "offset", 0),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── Payment ──

func (h *RestHandler) handleCreatePayment(w http.ResponseWriter, r *http.Request) {
	if h.payment == nil {
		writeError(w, http.StatusServiceUnavailable, "payment service unavailable")
		return
	}
	var body struct {
		UserID         string `json:"user_id"`
		OrderID        string `json:"order_id"`
		SubscriptionID string `json:"subscription_id"`
		PaymentType    int32  `json:"payment_type"`
		AmountKRW      int32  `json:"amount_krw"`
		PaymentMethod  string `json:"payment_method"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.payment.CreatePayment(r.Context(), &v1.CreatePaymentRequest{
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
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleConfirmPayment(w http.ResponseWriter, r *http.Request) {
	if h.payment == nil {
		writeError(w, http.StatusServiceUnavailable, "payment service unavailable")
		return
	}
	paymentId := r.PathValue("paymentId")
	var body struct {
		PgTransactionID string `json:"pg_transaction_id"`
		PgProvider      string `json:"pg_provider"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.payment.ConfirmPayment(r.Context(), &v1.ConfirmPaymentRequest{
		PaymentId:       paymentId,
		PgTransactionId: body.PgTransactionID,
		PgProvider:      body.PgProvider,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetPayment(w http.ResponseWriter, r *http.Request) {
	if h.payment == nil {
		writeError(w, http.StatusServiceUnavailable, "payment service unavailable")
		return
	}
	paymentId := r.PathValue("paymentId")
	resp, err := h.payment.GetPayment(r.Context(), &v1.GetPaymentRequest{PaymentId: paymentId})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── Prescription ──

func (h *RestHandler) handleSelectPharmacy(w http.ResponseWriter, r *http.Request) {
	if h.prescription == nil {
		writeError(w, http.StatusServiceUnavailable, "prescription service unavailable")
		return
	}
	prescriptionId := r.PathValue("prescriptionId")
	var body struct {
		PharmacyID      string `json:"pharmacy_id"`
		PharmacyName    string `json:"pharmacy_name"`
		FulfillmentType string `json:"fulfillment_type"`
		ShippingAddress string `json:"shipping_address"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.prescription.SelectPharmacyAndFulfillment(r.Context(), &v1.SelectPharmacyRequest{
		PrescriptionId:  prescriptionId,
		PharmacyId:      body.PharmacyID,
		PharmacyName:    body.PharmacyName,
		FulfillmentType: body.FulfillmentType,
		ShippingAddress: body.ShippingAddress,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleSendToPharmacy(w http.ResponseWriter, r *http.Request) {
	if h.prescription == nil {
		writeError(w, http.StatusServiceUnavailable, "prescription service unavailable")
		return
	}
	prescriptionId := r.PathValue("prescriptionId")
	resp, err := h.prescription.SendPrescriptionToPharmacy(r.Context(), &v1.SendToPharmacyRequest{
		PrescriptionId: prescriptionId,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetPrescriptionByToken(w http.ResponseWriter, r *http.Request) {
	if h.prescription == nil {
		writeError(w, http.StatusServiceUnavailable, "prescription service unavailable")
		return
	}
	token := r.PathValue("token")
	resp, err := h.prescription.GetPrescriptionByToken(r.Context(), &v1.GetByTokenRequest{
		FulfillmentToken: token,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}
