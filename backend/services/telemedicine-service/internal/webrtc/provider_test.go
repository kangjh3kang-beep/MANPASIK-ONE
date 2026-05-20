package webrtc_test

import (
	"context"
	"testing"

	"github.com/manpasik/backend/services/telemedicine-service/internal/webrtc"
)

func TestNoopProvider_CreateRoom(t *testing.T) {
	provider := webrtc.NewNoopProvider()
	ctx := context.Background()

	room, err := provider.CreateRoom(ctx, "consult-123")
	if err != nil {
		t.Fatalf("CreateRoom 실패: %v", err)
	}
	if room == nil {
		t.Fatal("room이 nil입니다")
	}
	if room.RoomID != "noop-consult-123" {
		t.Errorf("RoomID = %q, want %q", room.RoomID, "noop-consult-123")
	}
	if room.RoomURL != "https://meet.manpasik.com/noop-consult-123" {
		t.Errorf("RoomURL = %q, want %q", room.RoomURL, "https://meet.manpasik.com/noop-consult-123")
	}
}

func TestNoopProvider_GenerateToken(t *testing.T) {
	provider := webrtc.NewNoopProvider()
	ctx := context.Background()

	token, err := provider.GenerateToken(ctx, "room-1", "user-1", webrtc.RolePatient)
	if err != nil {
		t.Fatalf("GenerateToken 실패: %v", err)
	}
	if token == nil {
		t.Fatal("token이 nil입니다")
	}
	if token.Token != "noop-token-user-1" {
		t.Errorf("Token = %q, want %q", token.Token, "noop-token-user-1")
	}
	if token.ChannelID != "room-1" {
		t.Errorf("ChannelID = %q, want %q", token.ChannelID, "room-1")
	}
	if token.UID != "user-1" {
		t.Errorf("UID = %q, want %q", token.UID, "user-1")
	}
}

func TestNoopProvider_GenerateToken_DoctorRole(t *testing.T) {
	provider := webrtc.NewNoopProvider()
	ctx := context.Background()

	token, err := provider.GenerateToken(ctx, "room-1", "doctor-1", webrtc.RoleDoctor)
	if err != nil {
		t.Fatalf("GenerateToken 실패: %v", err)
	}
	if token.Token != "noop-token-doctor-1" {
		t.Errorf("Token = %q, want %q", token.Token, "noop-token-doctor-1")
	}
}

func TestNoopProvider_EndRoom(t *testing.T) {
	provider := webrtc.NewNoopProvider()
	ctx := context.Background()

	err := provider.EndRoom(ctx, "room-1")
	if err != nil {
		t.Fatalf("EndRoom은 에러를 반환하지 않아야 합니다: %v", err)
	}
}

func TestNoopProvider_ProviderName(t *testing.T) {
	provider := webrtc.NewNoopProvider()
	if provider.ProviderName() != "noop" {
		t.Errorf("ProviderName = %q, want %q", provider.ProviderName(), "noop")
	}
}
