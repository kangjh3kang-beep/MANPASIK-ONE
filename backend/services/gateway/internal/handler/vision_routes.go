package handler

import (
	"net/http"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
)

func (h *RestHandler) registerVisionRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/vision/analyze-food", h.handleAnalyzeFood)
	mux.HandleFunc("GET /api/v1/vision/food-analysis/{analysisId}", h.handleGetFoodAnalysis)
	mux.HandleFunc("GET /api/v1/vision/food-analyses", h.handleListFoodAnalyses)
	mux.HandleFunc("GET /api/v1/vision/nutrition-summary", h.handleGetDailyNutritionSummary)
	mux.HandleFunc("POST /api/v1/vision/meals", h.handleLogMeal)
	mux.HandleFunc("GET /api/v1/vision/meals", h.handleGetMealHistory)
}

func (h *RestHandler) handleAnalyzeFood(w http.ResponseWriter, r *http.Request) {
	if h.vision == nil {
		writeError(w, http.StatusServiceUnavailable, "vision service unavailable")
		return
	}
	var body struct {
		UserID   string `json:"user_id"`
		ImageUrl string `json:"image_url"`
		MealType string `json:"meal_type"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.vision.AnalyzeFood(r.Context(), &v1.AnalyzeFoodRequest{
		UserId: body.UserID, ImageUrl: body.ImageUrl, MealType: body.MealType,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetFoodAnalysis(w http.ResponseWriter, r *http.Request) {
	if h.vision == nil {
		writeError(w, http.StatusServiceUnavailable, "vision service unavailable")
		return
	}
	analysisID := r.PathValue("analysisId")
	resp, err := h.vision.GetFoodAnalysis(r.Context(), &v1.GetFoodAnalysisRequest{AnalysisId: analysisID})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleListFoodAnalyses(w http.ResponseWriter, r *http.Request) {
	if h.vision == nil {
		writeError(w, http.StatusServiceUnavailable, "vision service unavailable")
		return
	}
	userID := r.URL.Query().Get("user_id")
	resp, err := h.vision.ListFoodAnalyses(r.Context(), &v1.ListFoodAnalysesRequest{
		UserId: userID, Limit: queryInt(r, "limit", 20),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetDailyNutritionSummary(w http.ResponseWriter, r *http.Request) {
	if h.vision == nil {
		writeError(w, http.StatusServiceUnavailable, "vision service unavailable")
		return
	}
	userID := r.URL.Query().Get("user_id")
	date := r.URL.Query().Get("date")
	resp, err := h.vision.GetDailyNutritionSummary(r.Context(), &v1.GetDailyNutritionSummaryRequest{
		UserId: userID, Date: date,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleLogMeal(w http.ResponseWriter, r *http.Request) {
	if h.vision == nil {
		writeError(w, http.StatusServiceUnavailable, "vision service unavailable")
		return
	}
	var body struct {
		UserID      string `json:"user_id"`
		MealType    string `json:"meal_type"`
		Description string `json:"description"`
		AnalysisID  string `json:"analysis_id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.vision.LogMeal(r.Context(), &v1.LogMealRequest{
		UserId: body.UserID, MealType: body.MealType, Description: body.Description, AnalysisId: body.AnalysisID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetMealHistory(w http.ResponseWriter, r *http.Request) {
	if h.vision == nil {
		writeError(w, http.StatusServiceUnavailable, "vision service unavailable")
		return
	}
	userID := r.URL.Query().Get("user_id")
	resp, err := h.vision.GetMealHistory(r.Context(), &v1.GetMealHistoryRequest{
		UserId: userID, Limit: queryInt(r, "limit", 20),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}
