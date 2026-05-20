// Package webrtc는 화상진료를 위한 WebRTC 제공자 인터페이스를 정의합니다.
package webrtc

import "context"

// RoomInfo는 WebRTC 방 생성 결과입니다.
type RoomInfo struct {
	RoomID  string // 제공자 측 방 ID
	RoomURL string // 클라이언트가 접속할 URL
}

// TokenInfo는 WebRTC 접속 토큰입니다.
type TokenInfo struct {
	Token     string // 클라이언트 접속 토큰
	ChannelID string // 채널/방 식별자
	UID       string // 사용자 식별자
}

// Provider는 WebRTC 화상 통화 제공자 인터페이스입니다.
type Provider interface {
	// CreateRoom은 새 화상 통화 방을 생성합니다.
	CreateRoom(ctx context.Context, consultationID string) (*RoomInfo, error)

	// GenerateToken은 방에 접속하기 위한 토큰을 발급합니다.
	GenerateToken(ctx context.Context, roomID, userID string, role TokenRole) (*TokenInfo, error)

	// EndRoom은 화상 통화 방을 종료합니다.
	EndRoom(ctx context.Context, roomID string) error

	// ProviderName은 제공자 이름을 반환합니다 ("agora", "twilio", "noop").
	ProviderName() string
}

// TokenRole은 토큰 발급 시 역할입니다.
type TokenRole int

const (
	RolePatient TokenRole = iota
	RoleDoctor
)

// NoopProvider는 외부 WebRTC API 없이 목 URL/토큰을 반환합니다.
type NoopProvider struct{}

// NewNoopProvider는 NoopProvider를 생성합니다.
func NewNoopProvider() *NoopProvider {
	return &NoopProvider{}
}

// CreateRoom은 목 방 정보를 반환합니다.
func (n *NoopProvider) CreateRoom(_ context.Context, consultationID string) (*RoomInfo, error) {
	return &RoomInfo{
		RoomID:  "noop-" + consultationID,
		RoomURL: "https://meet.manpasik.com/noop-" + consultationID,
	}, nil
}

// GenerateToken은 목 토큰을 반환합니다.
func (n *NoopProvider) GenerateToken(_ context.Context, roomID, userID string, _ TokenRole) (*TokenInfo, error) {
	return &TokenInfo{
		Token:     "noop-token-" + userID,
		ChannelID: roomID,
		UID:       userID,
	}, nil
}

// EndRoom은 아무 작업도 하지 않습니다.
func (n *NoopProvider) EndRoom(_ context.Context, _ string) error {
	return nil
}

// ProviderName은 "noop"을 반환합니다.
func (n *NoopProvider) ProviderName() string {
	return "noop"
}
