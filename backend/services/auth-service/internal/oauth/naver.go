package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// NaverVerifier는 네이버 Access Token을 검증합니다.
type NaverVerifier struct {
	httpClient *http.Client
	apiURL     string // 테스트에서 교체 가능
}

// NewNaverVerifier는 NaverVerifier를 생성합니다.
func NewNaverVerifier() *NaverVerifier {
	return &NaverVerifier{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		apiURL:     "https://openapi.naver.com/v1/nid/me",
	}
}

// SetAPIURL은 네이버 API URL을 설정합니다 (테스트용).
func (n *NaverVerifier) SetAPIURL(url string) {
	n.apiURL = url
}

type naverResponse struct {
	ResultCode string `json:"resultcode"`
	Message    string `json:"message"`
	Response   struct {
		ID           string `json:"id"`
		Email        string `json:"email"`
		Nickname     string `json:"nickname"`
		ProfileImage string `json:"profile_image"`
	} `json:"response"`
}

// VerifyToken은 네이버 Access Token을 검증합니다.
func (n *NaverVerifier) VerifyToken(ctx context.Context, _, accessToken string) (*OAuthUserInfo, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("네이버 Access Token이 비어있습니다")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, n.apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("네이버 API 요청 생성 실패: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("네이버 API 호출 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("네이버 API 오류 (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var naverResp naverResponse
	if err := json.NewDecoder(resp.Body).Decode(&naverResp); err != nil {
		return nil, fmt.Errorf("네이버 응답 파싱 실패: %w", err)
	}

	if naverResp.ResultCode != "00" {
		return nil, fmt.Errorf("네이버 API 결과 오류: %s (%s)", naverResp.ResultCode, naverResp.Message)
	}

	if naverResp.Response.ID == "" {
		return nil, fmt.Errorf("네이버 사용자 ID가 비어있습니다")
	}

	email := naverResp.Response.Email
	if email == "" {
		email = fmt.Sprintf("naver_%s@social.manpasik.com", naverResp.Response.ID)
	}

	return &OAuthUserInfo{
		ProviderID:   naverResp.Response.ID,
		Email:        email,
		DisplayName:  naverResp.Response.Nickname,
		ProfileImage: naverResp.Response.ProfileImage,
		Provider:     "naver",
	}, nil
}

// Provider는 "naver"를 반환합니다.
func (n *NaverVerifier) Provider() string { return "naver" }
