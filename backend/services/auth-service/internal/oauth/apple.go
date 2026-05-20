package oauth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AppleVerifier는 Apple Sign-In JWT를 검증합니다.
type AppleVerifier struct {
	clientID   string
	httpClient *http.Client
	keysURL    string // 테스트에서 교체 가능
	cachedKeys []appleJWK
	cachedAt   time.Time
	keysTTL    time.Duration
	mu         sync.RWMutex
}

// NewAppleVerifier는 AppleVerifier를 생성합니다.
func NewAppleVerifier(clientID string) *AppleVerifier {
	return &AppleVerifier{
		clientID:   clientID,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		keysURL:    "https://appleid.apple.com/auth/keys",
		keysTTL:    1 * time.Hour,
	}
}

// SetKeysURL은 Apple 공개키 URL을 설정합니다 (테스트용).
func (a *AppleVerifier) SetKeysURL(url string) {
	a.keysURL = url
}

type appleJWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type appleKeysResponse struct {
	Keys []appleJWK `json:"keys"`
}

// VerifyToken은 Apple ID 토큰(JWT)을 검증합니다.
func (a *AppleVerifier) VerifyToken(ctx context.Context, idToken, _ string) (*OAuthUserInfo, error) {
	if idToken == "" {
		return nil, fmt.Errorf("Apple ID 토큰이 비어있습니다")
	}

	// JWT 헤더에서 kid 추출 (검증 전 파싱)
	parser := jwt.NewParser()
	token, parts, err := parser.ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil || len(parts) < 2 {
		return nil, fmt.Errorf("Apple JWT 파싱 실패: %w", err)
	}

	kid, _ := token.Header["kid"].(string)
	if kid == "" {
		return nil, fmt.Errorf("Apple JWT에 kid 헤더가 없습니다")
	}

	// Apple 공개키 가져오기 (캐시 활용)
	keys, err := a.fetchKeys(ctx)
	if err != nil {
		return nil, fmt.Errorf("Apple 공개키 조회 실패: %w", err)
	}

	// kid에 해당하는 키 찾기
	var matchedKey *appleJWK
	for i := range keys {
		if keys[i].Kid == kid {
			matchedKey = &keys[i]
			break
		}
	}
	if matchedKey == nil {
		return nil, fmt.Errorf("Apple 공개키에서 kid=%s를 찾을 수 없습니다", kid)
	}

	// RSA 공개키 구성
	pubKey, err := jwkToRSAPublicKey(matchedKey)
	if err != nil {
		return nil, fmt.Errorf("Apple RSA 키 변환 실패: %w", err)
	}

	// JWT 검증
	claims := jwt.MapClaims{}
	verifiedToken, err := jwt.ParseWithClaims(idToken, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("예상하지 못한 서명 알고리즘: %v", t.Header["alg"])
		}
		return pubKey, nil
	})
	if err != nil || !verifiedToken.Valid {
		return nil, fmt.Errorf("Apple JWT 검증 실패: %w", err)
	}

	// issuer, audience 검증
	iss, _ := claims["iss"].(string)
	if iss != "https://appleid.apple.com" {
		return nil, fmt.Errorf("Apple issuer 불일치: %s", iss)
	}
	aud, _ := claims["aud"].(string)
	if aud != a.clientID {
		return nil, fmt.Errorf("Apple audience 불일치: got %s, want %s", aud, a.clientID)
	}

	sub, _ := claims["sub"].(string)
	email, _ := claims["email"].(string)
	if email == "" {
		email = fmt.Sprintf("apple_%s@social.manpasik.com", sub)
	}

	return &OAuthUserInfo{
		ProviderID:  sub,
		Email:       email,
		DisplayName: "Apple 사용자",
		Provider:    "apple",
	}, nil
}

// Provider는 "apple"을 반환합니다.
func (a *AppleVerifier) Provider() string { return "apple" }

// fetchKeys는 Apple 공개키를 가져옵니다 (1시간 캐시).
func (a *AppleVerifier) fetchKeys(ctx context.Context) ([]appleJWK, error) {
	a.mu.RLock()
	if len(a.cachedKeys) > 0 && time.Since(a.cachedAt) < a.keysTTL {
		keys := a.cachedKeys
		a.mu.RUnlock()
		return keys, nil
	}
	a.mu.RUnlock()

	a.mu.Lock()
	defer a.mu.Unlock()

	// double-check
	if len(a.cachedKeys) > 0 && time.Since(a.cachedAt) < a.keysTTL {
		return a.cachedKeys, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.keysURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Apple 키 조회 실패 (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var keysResp appleKeysResponse
	if err := json.NewDecoder(resp.Body).Decode(&keysResp); err != nil {
		return nil, fmt.Errorf("Apple 키 파싱 실패: %w", err)
	}

	a.cachedKeys = keysResp.Keys
	a.cachedAt = time.Now()
	return a.cachedKeys, nil
}

// jwkToRSAPublicKey는 Apple JWK를 RSA 공개키로 변환합니다.
func jwkToRSAPublicKey(key *appleJWK) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		return nil, fmt.Errorf("N 디코딩 실패: %w", err)
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		return nil, fmt.Errorf("E 디코딩 실패: %w", err)
	}

	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)

	return &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}, nil
}
