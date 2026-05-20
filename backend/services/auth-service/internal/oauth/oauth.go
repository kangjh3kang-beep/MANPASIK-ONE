// Package oauth는 소셜 로그인 제공자(Google, Apple, Kakao, Naver) 토큰 검증을 구현합니다.
package oauth

import (
	"context"
	"fmt"
)

// OAuthUserInfo는 소셜 로그인 제공자에서 가져온 사용자 정보입니다.
type OAuthUserInfo struct {
	ProviderID   string // 제공자별 고유 ID
	Email        string
	DisplayName  string
	ProfileImage string
	Provider     string // "google", "apple", "kakao", "naver"
}

// OAuthVerifier는 소셜 로그인 토큰을 검증하는 인터페이스입니다.
type OAuthVerifier interface {
	VerifyToken(ctx context.Context, idToken, accessToken string) (*OAuthUserInfo, error)
	Provider() string
}

// NoopVerifier는 검증 없이 시뮬레이션 사용자 정보를 반환합니다.
// 개발/테스트 환경 또는 API 키 미설정 시 사용합니다.
type NoopVerifier struct {
	providerName string
}

// NewNoopVerifier는 NoopVerifier를 생성합니다.
func NewNoopVerifier(provider string) *NoopVerifier {
	return &NoopVerifier{providerName: provider}
}

// VerifyToken은 시뮬레이션 사용자 정보를 반환합니다.
func (n *NoopVerifier) VerifyToken(_ context.Context, idToken, accessToken string) (*OAuthUserInfo, error) {
	token := idToken
	if token == "" {
		token = accessToken
	}
	if token == "" {
		return nil, fmt.Errorf("토큰이 비어있습니다")
	}

	// 토큰 앞 8자를 이메일 생성에 사용 (기존 SocialLogin 패턴)
	prefix := token
	if len(prefix) > 8 {
		prefix = prefix[:8]
	}

	return &OAuthUserInfo{
		ProviderID:  fmt.Sprintf("noop_%s_%s", n.providerName, prefix),
		Email:       fmt.Sprintf("%s_%s@social.manpasik.com", n.providerName, prefix),
		DisplayName: n.providerName + " 사용자",
		Provider:    n.providerName,
	}, nil
}

// Provider는 제공자 이름을 반환합니다.
func (n *NoopVerifier) Provider() string {
	return n.providerName
}
