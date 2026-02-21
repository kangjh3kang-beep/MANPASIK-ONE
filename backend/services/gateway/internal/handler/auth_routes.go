package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

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
	// POST /api/v1/auth/kakao — 카카오 소셜 로그인
	mux.HandleFunc("POST /api/v1/auth/kakao", h.handleKakaoLogin)
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

// handleKakaoLogin은 카카오 소셜 로그인을 처리합니다.
//
// POST /api/v1/auth/kakao
// Body: { "kakao_access_token": "...", "email": "...", "nickname": "..." }
//
// 흐름:
// 1. 카카오 Access Token으로 카카오 사용자 정보 API 호출 → 사용자 ID 검증
// 2. 검증된 카카오 사용자 정보로 auth-service SocialLogin gRPC 호출
// 3. 만파식 JWT(Access/Refresh Token) 발급 → JSON 응답
func (h *RestHandler) handleKakaoLogin(w http.ResponseWriter, r *http.Request) {
	if h.auth == nil {
		writeError(w, http.StatusServiceUnavailable, "auth service unavailable")
		return
	}

	var body struct {
		KakaoAccessToken string `json:"kakao_access_token"`
		Email            string `json:"email"`
		Nickname         string `json:"nickname"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.KakaoAccessToken == "" {
		writeError(w, http.StatusBadRequest, "kakao_access_token is required")
		return
	}

	// 1단계: 카카오 사용자 정보 API로 토큰 검증
	kakaoUser, err := verifyKakaoAccessToken(r.Context(), body.KakaoAccessToken)
	if err != nil {
		writeError(w, http.StatusUnauthorized, fmt.Sprintf("카카오 인증 실패: %v", err))
		return
	}

	// 2단계: 검증된 카카오 ID를 idToken으로 사용하여 auth-service에 전달
	kakaoID := fmt.Sprintf("kakao_%d", kakaoUser.ID)
	resp, err := h.auth.SocialLogin(r.Context(), &v1.SocialLoginRequest{
		Provider:    v1.SocialProvider(v1.SocialProvider_value["kakao"]),
		IdToken:     kakaoID,
		AccessToken: body.KakaoAccessToken,
	})
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// 3단계: 추가 사용자 정보를 포함한 JSON 응답
	displayName := kakaoUser.Nickname
	if displayName == "" {
		displayName = body.Nickname
	}
	email := kakaoUser.Email
	if email == "" {
		email = body.Email
	}

	result := map[string]interface{}{
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"expires_in":    resp.ExpiresIn,
		"token_type":    resp.TokenType,
		"user_id":       kakaoID,
		"email":         email,
		"display_name":  displayName,
		"kakao_id":      kakaoUser.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// kakaoUserResponse는 카카오 사용자 정보 API 응답 구조체입니다.
type kakaoUserResponse struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

// verifyKakaoAccessToken은 카카오 Access Token으로 사용자 정보 API를 호출하여
// 토큰을 검증하고 사용자 고유 ID를 반환합니다.
func verifyKakaoAccessToken(ctx context.Context, accessToken string) (*kakaoUserResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://kapi.kakao.com/v2/user/me", nil)
	if err != nil {
		return nil, fmt.Errorf("요청 생성 실패: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("카카오 API 호출 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("카카오 API 오류 (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var kakaoResp struct {
		ID           int64 `json:"id"`
		KakaoAccount struct {
			Email   string `json:"email"`
			Profile struct {
				Nickname string `json:"nickname"`
			} `json:"profile"`
		} `json:"kakao_account"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&kakaoResp); err != nil {
		return nil, fmt.Errorf("응답 파싱 실패: %w", err)
	}

	if kakaoResp.ID == 0 {
		return nil, fmt.Errorf("유효하지 않은 카카오 토큰 — 사용자 ID 없음")
	}

	return &kakaoUserResponse{
		ID:       kakaoResp.ID,
		Email:    kakaoResp.KakaoAccount.Email,
		Nickname: kakaoResp.KakaoAccount.Profile.Nickname,
	}, nil
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
