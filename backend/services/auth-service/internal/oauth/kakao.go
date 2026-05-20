package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// KakaoVerifier는 카카오 Access Token을 검증합니다.
type KakaoVerifier struct {
	httpClient *http.Client
	apiURL     string // 테스트에서 교체 가능
}

// NewKakaoVerifier는 KakaoVerifier를 생성합니다.
func NewKakaoVerifier() *KakaoVerifier {
	return &KakaoVerifier{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		apiURL:     "https://kapi.kakao.com/v2/user/me",
	}
}

// SetAPIURL은 카카오 API URL을 설정합니다 (테스트용).
func (k *KakaoVerifier) SetAPIURL(url string) {
	k.apiURL = url
}

// VerifyToken은 카카오 Access Token을 검증합니다.
func (k *KakaoVerifier) VerifyToken(ctx context.Context, _, accessToken string) (*OAuthUserInfo, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("카카오 Access Token이 비어있습니다")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, k.apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("카카오 API 요청 생성 실패: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	resp, err := k.httpClient.Do(req)
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
				Nickname     string `json:"nickname"`
				ProfileImage string `json:"profile_image_url"`
			} `json:"profile"`
		} `json:"kakao_account"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&kakaoResp); err != nil {
		return nil, fmt.Errorf("카카오 응답 파싱 실패: %w", err)
	}

	if kakaoResp.ID == 0 {
		return nil, fmt.Errorf("카카오 사용자 ID가 0입니다 — 유효하지 않은 토큰")
	}

	email := kakaoResp.KakaoAccount.Email
	if email == "" {
		email = fmt.Sprintf("kakao_%d@social.manpasik.com", kakaoResp.ID)
	}

	return &OAuthUserInfo{
		ProviderID:   strconv.FormatInt(kakaoResp.ID, 10),
		Email:        email,
		DisplayName:  kakaoResp.KakaoAccount.Profile.Nickname,
		ProfileImage: kakaoResp.KakaoAccount.Profile.ProfileImage,
		Provider:     "kakao",
	}, nil
}

// Provider는 "kakao"를 반환합니다.
func (k *KakaoVerifier) Provider() string { return "kakao" }
