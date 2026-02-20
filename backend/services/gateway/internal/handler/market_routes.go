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
	mux.HandleFunc("POST /api/v1/prescriptions", h.handleCreatePrescription)
	mux.HandleFunc("GET /api/v1/prescriptions", h.handleListPrescriptions)
	mux.HandleFunc("GET /api/v1/prescriptions/{prescriptionId}", h.handleGetPrescription)
	mux.HandleFunc("PUT /api/v1/prescriptions/{prescriptionId}/status", h.handleUpdatePrescriptionStatus)
	mux.HandleFunc("POST /api/v1/prescriptions/{prescriptionId}/medications", h.handleAddMedication)
	mux.HandleFunc("DELETE /api/v1/prescriptions/{prescriptionId}/medications/{medicationId}", h.handleRemoveMedication)
	mux.HandleFunc("POST /api/v1/prescriptions/drug-interactions", h.handleCheckDrugInteraction)
	mux.HandleFunc("GET /api/v1/prescriptions/reminders", h.handleGetMedicationReminders)
	mux.HandleFunc("POST /api/v1/prescriptions/{prescriptionId}/pharmacy", h.handleSelectPharmacy)
	mux.HandleFunc("POST /api/v1/prescriptions/{prescriptionId}/send", h.handleSendToPharmacy)
	mux.HandleFunc("GET /api/v1/prescriptions/token/{token}", h.handleGetPrescriptionByToken)
	mux.HandleFunc("PUT /api/v1/prescriptions/{prescriptionId}/dispensary", h.handleUpdateDispensaryStatus)
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
	productId := r.PathValue("productId")
	if h.community == nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"reviews": []interface{}{}, "total": 0, "product_id": productId,
		})
		return
	}
	resp, err := h.community.ListPosts(r.Context(), &v1.ListPostsRequest{
		Category: v1.PostCategory_POST_CATEGORY_EXPERIENCE,
		Query:    productId,
		Limit:    queryInt(r, "limit", 20),
		Offset:   queryInt(r, "offset", 0),
	})
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"reviews": []interface{}{}, "total": 0, "product_id": productId,
		})
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleCreateProductReview(w http.ResponseWriter, r *http.Request) {
	productId := r.PathValue("productId")
	if h.community == nil {
		writeError(w, http.StatusServiceUnavailable, "community service unavailable")
		return
	}
	var body struct {
		UserID  string `json:"user_id"`
		Content string `json:"content"`
		Rating  int32  `json:"rating"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.community.CreatePost(r.Context(), &v1.CreatePostRequest{
		AuthorId: body.UserID,
		Title:    "상품 리뷰: " + productId,
		Content:  body.Content,
		Category: v1.PostCategory_POST_CATEGORY_EXPERIENCE,
		Tags:     []string{"product:" + productId, "review"},
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
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

func (h *RestHandler) handleCreatePrescription(w http.ResponseWriter, r *http.Request) {
	if h.prescription == nil {
		writeError(w, http.StatusServiceUnavailable, "prescription service unavailable")
		return
	}
	var body struct {
		UserID      string `json:"user_id"`
		DoctorID    string `json:"doctor_id"`
		DoctorName  string `json:"doctor_name"`
		FacilityID  string `json:"facility_id"`
		Diagnosis   string `json:"diagnosis"`
		Notes       string `json:"notes"`
		Medications []struct {
			MedicationID string `json:"medication_id"`
			Name         string `json:"name"`
			Dosage       string `json:"dosage"`
			Frequency    string `json:"frequency"`
			Route        string `json:"route"`
			DurationDays int32  `json:"duration_days"`
			Instructions string `json:"instructions"`
			IsCritical   bool   `json:"is_critical"`
		} `json:"medications"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	meds := make([]*v1.Medication, len(body.Medications))
	for i, m := range body.Medications {
		meds[i] = &v1.Medication{
			MedicationId: m.MedicationID,
			Name:         m.Name,
			Dosage:       m.Dosage,
			Frequency:    m.Frequency,
			Route:        m.Route,
			DurationDays: m.DurationDays,
			Instructions: m.Instructions,
			IsCritical:   m.IsCritical,
		}
	}
	resp, err := h.prescription.CreatePrescription(r.Context(), &v1.CreatePrescriptionRequest{
		UserId:      body.UserID,
		DoctorId:    body.DoctorID,
		DoctorName:  body.DoctorName,
		FacilityId:  body.FacilityID,
		Diagnosis:   body.Diagnosis,
		Notes:       body.Notes,
		Medications: meds,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleGetPrescription(w http.ResponseWriter, r *http.Request) {
	if h.prescription == nil {
		writeError(w, http.StatusServiceUnavailable, "prescription service unavailable")
		return
	}
	prescriptionId := r.PathValue("prescriptionId")
	resp, err := h.prescription.GetPrescription(r.Context(), &v1.GetPrescriptionRequest{
		PrescriptionId: prescriptionId,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleListPrescriptions(w http.ResponseWriter, r *http.Request) {
	if h.prescription == nil {
		writeError(w, http.StatusServiceUnavailable, "prescription service unavailable")
		return
	}
	resp, err := h.prescription.ListPrescriptions(r.Context(), &v1.ListPrescriptionsRequest{
		UserId:       r.URL.Query().Get("user_id"),
		StatusFilter: v1.PrescriptionStatus(queryInt(r, "status", 0)),
		Limit:        queryInt(r, "limit", 20),
		Offset:       queryInt(r, "offset", 0),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleUpdatePrescriptionStatus(w http.ResponseWriter, r *http.Request) {
	if h.prescription == nil {
		writeError(w, http.StatusServiceUnavailable, "prescription service unavailable")
		return
	}
	prescriptionId := r.PathValue("prescriptionId")
	var body struct {
		NewStatus int32 `json:"new_status"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.prescription.UpdatePrescriptionStatus(r.Context(), &v1.UpdatePrescriptionStatusRequest{
		PrescriptionId: prescriptionId,
		NewStatus:      v1.PrescriptionStatus(body.NewStatus),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleAddMedication(w http.ResponseWriter, r *http.Request) {
	if h.prescription == nil {
		writeError(w, http.StatusServiceUnavailable, "prescription service unavailable")
		return
	}
	prescriptionId := r.PathValue("prescriptionId")
	var body struct {
		MedicationID string `json:"medication_id"`
		Name         string `json:"name"`
		Dosage       string `json:"dosage"`
		Frequency    string `json:"frequency"`
		Route        string `json:"route"`
		DurationDays int32  `json:"duration_days"`
		Instructions string `json:"instructions"`
		IsCritical   bool   `json:"is_critical"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.prescription.AddMedication(r.Context(), &v1.AddMedicationRequest{
		PrescriptionId: prescriptionId,
		Medication: &v1.Medication{
			MedicationId: body.MedicationID,
			Name:         body.Name,
			Dosage:       body.Dosage,
			Frequency:    body.Frequency,
			Route:        body.Route,
			DurationDays: body.DurationDays,
			Instructions: body.Instructions,
			IsCritical:   body.IsCritical,
		},
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleRemoveMedication(w http.ResponseWriter, r *http.Request) {
	if h.prescription == nil {
		writeError(w, http.StatusServiceUnavailable, "prescription service unavailable")
		return
	}
	prescriptionId := r.PathValue("prescriptionId")
	medicationId := r.PathValue("medicationId")
	resp, err := h.prescription.RemoveMedication(r.Context(), &v1.RemoveMedicationRequest{
		PrescriptionId: prescriptionId,
		MedicationId:   medicationId,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleCheckDrugInteraction(w http.ResponseWriter, r *http.Request) {
	if h.prescription == nil {
		writeError(w, http.StatusServiceUnavailable, "prescription service unavailable")
		return
	}
	var body struct {
		MedicationNames []string `json:"medication_names"`
		UserID          string   `json:"user_id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.prescription.CheckDrugInteraction(r.Context(), &v1.CheckDrugInteractionRequest{
		MedicationNames: body.MedicationNames,
		UserId:          body.UserID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetMedicationReminders(w http.ResponseWriter, r *http.Request) {
	if h.prescription == nil {
		writeError(w, http.StatusServiceUnavailable, "prescription service unavailable")
		return
	}
	resp, err := h.prescription.GetMedicationReminders(r.Context(), &v1.GetMedicationRemindersRequest{
		UserId: r.URL.Query().Get("user_id"),
		Date:   r.URL.Query().Get("date"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleUpdateDispensaryStatus(w http.ResponseWriter, r *http.Request) {
	if h.prescription == nil {
		writeError(w, http.StatusServiceUnavailable, "prescription service unavailable")
		return
	}
	prescriptionId := r.PathValue("prescriptionId")
	var body struct {
		Status string `json:"status"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.prescription.UpdateDispensaryStatus(r.Context(), &v1.UpdateDispensaryStatusRequest{
		PrescriptionId: prescriptionId,
		Status:         body.Status,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

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
