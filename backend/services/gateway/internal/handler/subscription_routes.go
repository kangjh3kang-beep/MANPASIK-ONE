package handler

import (
	"net/http"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
)

// registerSubscriptionRoutes는 구독 관련 REST 엔드포인트를 등록합니다.
func (h *RestHandler) registerSubscriptionRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/subscriptions", h.handleCreateSubscription)
	mux.HandleFunc("GET /api/v1/subscriptions", h.handleGetSubscription)
	mux.HandleFunc("PUT /api/v1/subscriptions", h.handleUpdateSubscription)
	mux.HandleFunc("POST /api/v1/subscriptions/cancel", h.handleCancelSubscription)
	mux.HandleFunc("GET /api/v1/subscriptions/feature-access", h.handleCheckFeatureAccess)
	mux.HandleFunc("GET /api/v1/subscriptions/plans", h.handleListSubscriptionPlans)
	mux.HandleFunc("GET /api/v1/subscriptions/cartridge-access", h.handleCheckCartridgeAccess)
	mux.HandleFunc("GET /api/v1/subscriptions/accessible-cartridges", h.handleListAccessibleCartridges)
}

func (h *RestHandler) handleCreateSubscription(w http.ResponseWriter, r *http.Request) {
	if h.subscription == nil {
		writeError(w, http.StatusServiceUnavailable, "subscription service unavailable")
		return
	}
	var body struct {
		UserID string `json:"user_id"`
		Tier   int32  `json:"tier"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.subscription.CreateSubscription(r.Context(), &v1.CreateSubscriptionRequest{
		UserId: body.UserID,
		Tier:   v1.SubscriptionTier(body.Tier),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleGetSubscription(w http.ResponseWriter, r *http.Request) {
	if h.subscription == nil {
		writeError(w, http.StatusServiceUnavailable, "subscription service unavailable")
		return
	}
	resp, err := h.subscription.GetSubscription(r.Context(), &v1.GetSubscriptionDetailRequest{
		UserId: r.URL.Query().Get("user_id"),
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleUpdateSubscription(w http.ResponseWriter, r *http.Request) {
	if h.subscription == nil {
		writeError(w, http.StatusServiceUnavailable, "subscription service unavailable")
		return
	}
	var body struct {
		UserID    string `json:"user_id"`
		NewTier   int32  `json:"new_tier"`
		PaymentID string `json:"payment_id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.subscription.UpdateSubscription(r.Context(), &v1.UpdateSubscriptionRequest{
		UserId:    body.UserID,
		NewTier:   v1.SubscriptionTier(body.NewTier),
		PaymentId: body.PaymentID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleCancelSubscription(w http.ResponseWriter, r *http.Request) {
	if h.subscription == nil {
		writeError(w, http.StatusServiceUnavailable, "subscription service unavailable")
		return
	}
	var body struct {
		UserID string `json:"user_id"`
		Reason string `json:"reason"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.subscription.CancelSubscription(r.Context(), &v1.CancelSubscriptionRequest{
		UserId: body.UserID,
		Reason: body.Reason,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleCheckFeatureAccess(w http.ResponseWriter, r *http.Request) {
	if h.subscription == nil {
		writeError(w, http.StatusServiceUnavailable, "subscription service unavailable")
		return
	}
	resp, err := h.subscription.CheckFeatureAccess(r.Context(), &v1.CheckFeatureAccessRequest{
		UserId:      r.URL.Query().Get("user_id"),
		FeatureName: r.URL.Query().Get("feature"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleListSubscriptionPlans(w http.ResponseWriter, r *http.Request) {
	if h.subscription == nil {
		writeError(w, http.StatusServiceUnavailable, "subscription service unavailable")
		return
	}
	resp, err := h.subscription.ListSubscriptionPlans(r.Context(), &v1.ListSubscriptionPlansRequest{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleCheckCartridgeAccess(w http.ResponseWriter, r *http.Request) {
	if h.subscription == nil {
		writeError(w, http.StatusServiceUnavailable, "subscription service unavailable")
		return
	}
	resp, err := h.subscription.CheckCartridgeAccess(r.Context(), &v1.CheckCartridgeAccessRequest{
		UserId:       r.URL.Query().Get("user_id"),
		CategoryCode: queryInt(r, "category_code", 0),
		TypeIndex:    queryInt(r, "type_index", 0),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleListAccessibleCartridges(w http.ResponseWriter, r *http.Request) {
	if h.subscription == nil {
		writeError(w, http.StatusServiceUnavailable, "subscription service unavailable")
		return
	}
	resp, err := h.subscription.ListAccessibleCartridges(r.Context(), &v1.ListAccessibleCartridgesRequest{
		UserId: r.URL.Query().Get("user_id"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}
