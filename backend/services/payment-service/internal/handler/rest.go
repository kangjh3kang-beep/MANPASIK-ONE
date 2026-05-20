// Package handler는 payment-service의 REST 핸들러를 제공합니다.
package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/manpasik/backend/services/payment-service/internal/service"
)

// RESTHandler는 결제 REST API 핸들러입니다.
type RESTHandler struct {
	svc           *service.PaymentService
	webhookSecret string // Toss 웹훅 시그니처 검증 키
}

// NewRESTHandler는 새 RESTHandler를 생성합니다.
func NewRESTHandler(svc *service.PaymentService, webhookSecret string) *RESTHandler {
	return &RESTHandler{svc: svc, webhookSecret: webhookSecret}
}

// RegisterRoutes는 결제 REST 라우트를 등록합니다.
func (h *RESTHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/payments", h.handleCreatePayment)
	mux.HandleFunc("POST /api/v1/payments/confirm", h.handleConfirmPayment)
	mux.HandleFunc("GET /api/v1/payments/{paymentID}", h.handleGetPayment)
	mux.HandleFunc("GET /api/v1/payments", h.handleListPayments)
	mux.HandleFunc("POST /api/v1/payments/{paymentID}/refund", h.handleRefundPayment)
	mux.HandleFunc("POST /webhooks/payments/toss", h.handleTossWebhook)
}

// ----------- Create Payment -----------

type createPaymentRequest struct {
	UserID         string `json:"user_id"`
	OrderID        string `json:"order_id"`
	SubscriptionID string `json:"subscription_id"`
	PaymentType    int32  `json:"payment_type"`
	AmountKRW      int32  `json:"amount_krw"`
	PaymentMethod  string `json:"payment_method"`
}

func (h *RESTHandler) handleCreatePayment(w http.ResponseWriter, r *http.Request) {
	var req createPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.UserID == "" {
		writeError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	p, err := h.svc.CreatePayment(r.Context(), req.UserID, req.OrderID, req.SubscriptionID,
		service.PaymentType(req.PaymentType), req.AmountKRW, req.PaymentMethod)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, p)
}

// ----------- Confirm Payment -----------

type confirmPaymentRequest struct {
	PaymentID     string `json:"payment_id"`
	PgTxID        string `json:"pg_transaction_id"`
	PgProvider    string `json:"pg_provider"`
	PaymentKey    string `json:"payment_key"`
}

func (h *RESTHandler) handleConfirmPayment(w http.ResponseWriter, r *http.Request) {
	var req confirmPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.PaymentID == "" {
		writeError(w, http.StatusBadRequest, "payment_id is required")
		return
	}

	p, err := h.svc.ConfirmPayment(r.Context(), req.PaymentID, req.PgTxID, req.PgProvider, req.PaymentKey)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, p)
}

// ----------- Get Payment -----------

func (h *RESTHandler) handleGetPayment(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentID")
	if paymentID == "" {
		writeError(w, http.StatusBadRequest, "payment_id is required")
		return
	}

	p, err := h.svc.GetPayment(r.Context(), paymentID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, p)
}

// ----------- List Payments -----------

func (h *RESTHandler) handleListPayments(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	payments, total, err := h.svc.ListPayments(r.Context(), userID, 50, 0)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"payments":    payments,
		"total_count": total,
	})
}

// ----------- Refund Payment -----------

type refundRequest struct {
	RefundAmountKRW int32  `json:"refund_amount_krw"`
	Reason          string `json:"reason"`
}

func (h *RESTHandler) handleRefundPayment(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentID")
	if paymentID == "" {
		writeError(w, http.StatusBadRequest, "payment_id is required")
		return
	}

	var req refundRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	refund, payment, err := h.svc.RefundPayment(r.Context(), paymentID, req.RefundAmountKRW, req.Reason)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"refund_id":        refund.ID,
		"payment_id":       refund.PaymentID,
		"refund_amount_krw": refund.RefundAmountKRW,
		"payment_status":   payment.Status,
		"refunded_at":      refund.CreatedAt,
	})
}

// ----------- Toss Webhook -----------

type tossWebhookPayload struct {
	EventType  string `json:"eventType"`
	PaymentKey string `json:"paymentKey"`
	OrderID    string `json:"orderId"`
	Status     string `json:"status"`
	Amount     int32  `json:"totalAmount"`
}

func (h *RESTHandler) handleTossWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "cannot read body")
		return
	}

	// HMAC-SHA256 시그니처 검증 (Toss 웹훅 보안)
	if h.webhookSecret != "" {
		signature := r.Header.Get("X-Toss-Signature")
		if !h.verifySignature(body, signature) {
			log.Printf("[payment-webhook] signature verification failed")
			writeError(w, http.StatusUnauthorized, "invalid signature")
			return
		}
	}

	var payload tossWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid webhook payload")
		return
	}

	log.Printf("[payment-webhook] event=%s paymentKey=%s orderID=%s status=%s",
		payload.EventType, payload.PaymentKey, payload.OrderID, payload.Status)

	switch payload.EventType {
	case "PAYMENT_STATUS_CHANGED":
		if payload.Status == "DONE" {
			// 결제 완료: paymentKey로 결제 확인
			_, confirmErr := h.svc.ConfirmPaymentByWebhook(r.Context(), payload.PaymentKey, payload.OrderID, payload.Amount)
			if confirmErr != nil {
				log.Printf("[payment-webhook] confirm error: %v", confirmErr)
				// 웹훅은 200 반환해야 재전송 방지
			}
		}
	case "PAYMENT_CANCEL":
		log.Printf("[payment-webhook] payment cancelled: %s", payload.PaymentKey)
	}

	// Toss 웹훅은 항상 200 반환 (실패해도 재전송 방지)
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *RESTHandler) verifySignature(body []byte, signature string) bool {
	if signature == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(h.webhookSecret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

// ----------- Helpers -----------

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeJSON error: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
