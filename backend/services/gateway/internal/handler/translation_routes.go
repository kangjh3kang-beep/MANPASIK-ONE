package handler

import (
	"net/http"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
)

// registerTranslationRoutes는 번역 관련 REST 엔드포인트를 등록합니다.
func (h *RestHandler) registerTranslationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/translations/text", h.handleTranslateText)
	mux.HandleFunc("POST /api/v1/translations/detect", h.handleDetectLanguage)
	mux.HandleFunc("GET /api/v1/translations/languages", h.handleListSupportedLanguages)
	mux.HandleFunc("POST /api/v1/translations/batch", h.handleTranslateBatch)
	mux.HandleFunc("GET /api/v1/translations/history", h.handleGetTranslationHistory)
	mux.HandleFunc("GET /api/v1/translations/usage", h.handleGetTranslationUsage)
	mux.HandleFunc("POST /api/v1/translations/realtime", h.handleTranslateRealtime)
}

func (h *RestHandler) handleTranslateText(w http.ResponseWriter, r *http.Request) {
	if h.translation == nil {
		writeError(w, http.StatusServiceUnavailable, "translation service unavailable")
		return
	}
	var body struct {
		Text           string `json:"text"`
		SourceLanguage string `json:"source_language"`
		TargetLanguage string `json:"target_language"`
		IsMedical      bool   `json:"is_medical"`
		Context        string `json:"context"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.translation.TranslateText(r.Context(), &v1.TranslateTextRequest{
		Text:           body.Text,
		SourceLanguage: body.SourceLanguage,
		TargetLanguage: body.TargetLanguage,
		IsMedical:      body.IsMedical,
		Context:        body.Context,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleDetectLanguage(w http.ResponseWriter, r *http.Request) {
	if h.translation == nil {
		writeError(w, http.StatusServiceUnavailable, "translation service unavailable")
		return
	}
	var body struct {
		Text string `json:"text"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.translation.DetectLanguage(r.Context(), &v1.DetectLanguageRequest{
		Text: body.Text,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleListSupportedLanguages(w http.ResponseWriter, r *http.Request) {
	if h.translation == nil {
		writeError(w, http.StatusServiceUnavailable, "translation service unavailable")
		return
	}
	resp, err := h.translation.ListSupportedLanguages(r.Context(), &v1.ListSupportedLanguagesRequest{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleTranslateBatch(w http.ResponseWriter, r *http.Request) {
	if h.translation == nil {
		writeError(w, http.StatusServiceUnavailable, "translation service unavailable")
		return
	}
	var body struct {
		Texts          []string `json:"texts"`
		SourceLanguage string   `json:"source_language"`
		TargetLanguage string   `json:"target_language"`
		IsMedical      bool     `json:"is_medical"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.translation.TranslateBatch(r.Context(), &v1.TranslateBatchRequest{
		Texts:          body.Texts,
		SourceLanguage: body.SourceLanguage,
		TargetLanguage: body.TargetLanguage,
		IsMedical:      body.IsMedical,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetTranslationHistory(w http.ResponseWriter, r *http.Request) {
	if h.translation == nil {
		writeError(w, http.StatusServiceUnavailable, "translation service unavailable")
		return
	}
	resp, err := h.translation.GetTranslationHistory(r.Context(), &v1.GetTranslationHistoryRequest{
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

func (h *RestHandler) handleGetTranslationUsage(w http.ResponseWriter, r *http.Request) {
	if h.translation == nil {
		writeError(w, http.StatusServiceUnavailable, "translation service unavailable")
		return
	}
	resp, err := h.translation.GetTranslationUsage(r.Context(), &v1.GetTranslationUsageRequest{
		UserId: r.URL.Query().Get("user_id"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleTranslateRealtime(w http.ResponseWriter, r *http.Request) {
	if h.translation == nil {
		writeError(w, http.StatusServiceUnavailable, "translation service unavailable")
		return
	}
	var body struct {
		Text           string `json:"text"`
		SourceLanguage string `json:"source_language"`
		TargetLanguage string `json:"target_language"`
		IsMedical      bool   `json:"is_medical"`
		Context        string `json:"context"`
		SessionID      string `json:"session_id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.translation.TranslateRealtime(r.Context(), &v1.TranslateRealtimeRequest{
		Text:           body.Text,
		SourceLanguage: body.SourceLanguage,
		TargetLanguage: body.TargetLanguage,
		IsMedical:      body.IsMedical,
		Context:        body.Context,
		SessionId:      body.SessionID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}
