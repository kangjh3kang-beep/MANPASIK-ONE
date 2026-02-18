package handler

import (
	"net/http"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
)

// registerAuthRoutes는 인증 관련 REST 엔드포인트를 등록합니다.
func (h *RestHandler) registerAuthRoutes(mux *http.ServeMux) {
	// POST /api/v1/auth/register
	mux.HandleFunc("POST /api/v1/auth/register", h.handleRegister)
	// POST /api/v1/auth/login
	mux.HandleFunc("POST /api/v1/auth/login", h.handleLogin)
	// POST /api/v1/auth/social-login
	mux.HandleFunc("POST /api/v1/auth/social-login", h.handleSocialLogin)
	// POST /api/v1/auth/refresh
	mux.HandleFunc("POST /api/v1/auth/refresh", h.handleRefreshToken)
	// POST /api/v1/auth/logout
	mux.HandleFunc("POST /api/v1/auth/logout", h.handleLogout)
	// POST /api/v1/auth/reset-password
	mux.HandleFunc("POST /api/v1/auth/reset-password", h.handleResetPassword)
}

func (h *RestHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	if h.auth == nil {
		writeError(w, http.StatusServiceUnavailable, "auth service unavailable")
		return
	}

	var body struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		DisplayName string `json:"display_name"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.auth.Register(r.Context(), &v1.RegisterRequest{
		Email:       body.Email,
		Password:    body.Password,
		DisplayName: body.DisplayName,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	if h.auth == nil {
		writeError(w, http.StatusServiceUnavailable, "auth service unavailable")
		return
	}

	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.auth.Login(r.Context(), &v1.LoginRequest{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleSocialLogin(w http.ResponseWriter, r *http.Request) {
	if h.auth == nil {
		writeError(w, http.StatusServiceUnavailable, "auth service unavailable")
		return
	}

	var body struct {
		Provider string `json:"provider"`
		IDToken  string `json:"id_token"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.auth.SocialLogin(r.Context(), &v1.SocialLoginRequest{
		Provider: v1.SocialProvider(v1.SocialProvider_value[body.Provider]),
		IdToken:  body.IDToken,
	})
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	if h.auth == nil {
		writeError(w, http.StatusServiceUnavailable, "auth service unavailable")
		return
	}

	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.auth.RefreshToken(r.Context(), &v1.RefreshTokenRequest{
		RefreshToken: body.RefreshToken,
	})
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleLogout(w http.ResponseWriter, r *http.Request) {
	if h.auth == nil {
		writeError(w, http.StatusServiceUnavailable, "auth service unavailable")
		return
	}

	var body struct {
		UserID string `json:"user_id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.auth.Logout(r.Context(), &v1.LogoutRequest{
		UserId: body.UserID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleResetPassword(w http.ResponseWriter, r *http.Request) {
	if h.auth == nil {
		writeError(w, http.StatusServiceUnavailable, "auth service unavailable")
		return
	}

	var body struct {
		Email string `json:"email"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.auth.ResetPassword(r.Context(), &v1.ResetPasswordRequest{
		Email: body.Email,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeProtoJSON(w, http.StatusOK, resp)
}
