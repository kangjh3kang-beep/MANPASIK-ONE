package handler

import (
	"net/http"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
)

func (h *RestHandler) registerVoiceProfileRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/voice-profiles", h.handleCreateVoiceProfile)
	mux.HandleFunc("GET /api/v1/voice-profiles/{profileId}", h.handleGetVoiceProfile)
	mux.HandleFunc("GET /api/v1/voice-profiles", h.handleListVoiceProfiles)
	mux.HandleFunc("POST /api/v1/voice-profiles/synthesize", h.handleSynthesizeTranslation)
	mux.HandleFunc("DELETE /api/v1/voice-profiles/{profileId}", h.handleDeleteVoiceProfile)
}

func (h *RestHandler) handleCreateVoiceProfile(w http.ResponseWriter, r *http.Request) {
	if h.voiceProfile == nil {
		writeError(w, http.StatusServiceUnavailable, "voice profile service unavailable")
		return
	}
	var body struct {
		UserID      string `json:"user_id"`
		ProfileName string `json:"profile_name"`
		Language    string `json:"language"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.voiceProfile.CreateVoiceProfile(r.Context(), &v1.CreateVoiceProfileRequest{
		UserId: body.UserID, ProfileName: body.ProfileName, Language: body.Language,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleGetVoiceProfile(w http.ResponseWriter, r *http.Request) {
	if h.voiceProfile == nil {
		writeError(w, http.StatusServiceUnavailable, "voice profile service unavailable")
		return
	}
	resp, err := h.voiceProfile.GetVoiceProfile(r.Context(), &v1.GetVoiceProfileRequest{
		ProfileId: r.PathValue("profileId"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleListVoiceProfiles(w http.ResponseWriter, r *http.Request) {
	if h.voiceProfile == nil {
		writeError(w, http.StatusServiceUnavailable, "voice profile service unavailable")
		return
	}
	userID := r.URL.Query().Get("user_id")
	resp, err := h.voiceProfile.ListVoiceProfiles(r.Context(), &v1.ListVoiceProfilesRequest{UserId: userID})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleSynthesizeTranslation(w http.ResponseWriter, r *http.Request) {
	if h.voiceProfile == nil {
		writeError(w, http.StatusServiceUnavailable, "voice profile service unavailable")
		return
	}
	var body struct {
		ProfileID      string `json:"profile_id"`
		Text           string `json:"text"`
		TargetLanguage string `json:"target_language"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.voiceProfile.SynthesizeTranslation(r.Context(), &v1.SynthesizeTranslationRequest{
		ProfileId: body.ProfileID, Text: body.Text, TargetLanguage: body.TargetLanguage,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleDeleteVoiceProfile(w http.ResponseWriter, r *http.Request) {
	if h.voiceProfile == nil {
		writeError(w, http.StatusServiceUnavailable, "voice profile service unavailable")
		return
	}
	resp, err := h.voiceProfile.DeleteVoiceProfile(r.Context(), &v1.DeleteVoiceProfileRequest{
		ProfileId: r.PathValue("profileId"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}
