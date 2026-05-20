package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GoogleVerifier는 Google OAuth ID 토큰을 검증합니다.
type GoogleVerifier struct {
	clientID   string
	httpClient *http.Client
	tokenURL   string // 테스트에서 교체 가능
}

// NewGoogleVerifier는 GoogleVerifier를 생성합니다.
func NewGoogleVerifier(clientID string) *GoogleVerifier {
	return &GoogleVerifier{
		clientID:   clientID,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		tokenURL:   "https://oauth2.googleapis.com/tokeninfo",
	}
}

// SetTokenURL은 토큰 검증 URL을 설정합니다 (테스트용).
func (g *GoogleVerifier) SetTokenURL(url string) {
	g.tokenURL = url
}

type googleTokenInfo struct {
	Iss           string `json:"iss"`
	Sub           string `json:"sub"`
	Aud           string `json:"aud"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

// VerifyToken은 Google ID 토큰을 검증합니다.
func (g *GoogleVerifier) VerifyToken(ctx context.Context, idToken, _ string) (*OAuthUserInfo, error) {
	if idToken == "" {
		return nil, fmt.Errorf("Google ID 토큰이 비어있습니다")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s?id_token=%s", g.tokenURL, idToken), nil)
	if err != nil {
		return nil, fmt.Errorf("Google 요청 생성 실패: %w", err)
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Google API 호출 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Google 토큰 검증 실패 (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var info googleTokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("Google 응답 파싱 실패: %w", err)
	}

	// issuer 검증
	if info.Iss != "accounts.google.com" && info.Iss != "https://accounts.google.com" {
		return nil, fmt.Errorf("Google issuer 불일치: %s", info.Iss)
	}

	// audience 검증
	if info.Aud != g.clientID {
		return nil, fmt.Errorf("Google audience 불일치: got %s, want %s", info.Aud, g.clientID)
	}

	if info.Sub == "" {
		return nil, fmt.Errorf("Google 사용자 ID(sub)가 비어있습니다")
	}

	email := info.Email
	if email == "" {
		email = fmt.Sprintf("google_%s@social.manpasik.com", info.Sub)
	}

	return &OAuthUserInfo{
		ProviderID:   info.Sub,
		Email:        email,
		DisplayName:  info.Name,
		ProfileImage: info.Picture,
		Provider:     "google",
	}, nil
}

// Provider는 "google"을 반환합니다.
func (g *GoogleVerifier) Provider() string { return "google" }
