package router

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (r *Router) dialAuth() (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return grpc.DialContext(ctx, r.authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

// POST /api/v1/auth/register
func (r *Router) handleRegister(w http.ResponseWriter, req *http.Request) {
	var body struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		DisplayName string `json:"display_name"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := r.dialAuth()
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "인증 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Register(ctx, &v1.RegisterRequest{
		Email:       body.Email,
		Password:    body.Password,
		DisplayName: body.DisplayName,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"user_id":      resp.GetUserId(),
		"email":        resp.GetEmail(),
		"display_name": resp.GetDisplayName(),
		"role":         resp.GetRole(),
	})
}

// POST /api/v1/auth/login
func (r *Router) handleLogin(w http.ResponseWriter, req *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := r.dialAuth()
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "인증 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Login(ctx, &v1.LoginRequest{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"access_token":  resp.GetAccessToken(),
		"refresh_token": resp.GetRefreshToken(),
		"expires_in":    resp.GetExpiresIn(),
		"token_type":    resp.GetTokenType(),
	})
}

// POST /api/v1/auth/refresh
func (r *Router) handleRefreshToken(w http.ResponseWriter, req *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := r.dialAuth()
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "인증 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.RefreshToken(ctx, &v1.RefreshTokenRequest{
		RefreshToken: body.RefreshToken,
	})
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"access_token":  resp.GetAccessToken(),
		"refresh_token": resp.GetRefreshToken(),
		"expires_in":    resp.GetExpiresIn(),
		"token_type":    resp.GetTokenType(),
	})
}

// POST /api/v1/auth/logout
func (r *Router) handleLogout(w http.ResponseWriter, req *http.Request) {
	var body struct {
		UserId string `json:"user_id"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := r.dialAuth()
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "인증 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Logout(ctx, &v1.LogoutRequest{
		UserId: body.UserId,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "로그아웃 완료"})
}
