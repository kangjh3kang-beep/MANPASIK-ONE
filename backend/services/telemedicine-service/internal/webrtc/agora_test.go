package webrtc_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/manpasik/backend/services/telemedicine-service/internal/webrtc"
)

func TestAgoraProvider_CreateRoom_Success(t *testing.T) {
	provider := webrtc.NewAgoraProvider(webrtc.AgoraConfig{
		AppID:   "test-app-id",
		AppCert: "test-cert",
	})
	ctx := context.Background()

	room, err := provider.CreateRoom(ctx, "consult-456")
	if err != nil {
		t.Fatalf("CreateRoom 실패: %v", err)
	}
	if room.RoomID != "manpasik-consult-456" {
		t.Errorf("RoomID = %q, want %q", room.RoomID, "manpasik-consult-456")
	}
	if room.RoomURL != "https://meet.manpasik.com/agora/manpasik-consult-456" {
		t.Errorf("RoomURL = %q, want %q", room.RoomURL, "https://meet.manpasik.com/agora/manpasik-consult-456")
	}
}

func TestAgoraProvider_CreateRoom_EmptyAppID(t *testing.T) {
	provider := webrtc.NewAgoraProvider(webrtc.AgoraConfig{
		AppID: "",
	})
	ctx := context.Background()

	_, err := provider.CreateRoom(ctx, "consult-1")
	if err == nil {
		t.Fatal("빈 AppID에 에러가 반환되어야 합니다")
	}
}

func TestAgoraProvider_GenerateToken_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Basic Auth 확인
		user, pass, ok := r.BasicAuth()
		if !ok || user != "test-customer" || pass != "test-secret" {
			t.Error("Basic Auth 인증 실패")
		}

		resp := map[string]string{"resourceId": "res-123"}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	provider := webrtc.NewAgoraProvider(webrtc.AgoraConfig{
		AppID:      "test-app-id",
		AppCert:    "test-cert",
		CustomerID: "test-customer",
		Secret:     "test-secret",
		BaseURL:    server.URL,
	})
	ctx := context.Background()

	token, err := provider.GenerateToken(ctx, "room-1", "user-1", webrtc.RolePatient)
	if err != nil {
		t.Fatalf("GenerateToken 실패: %v", err)
	}
	if token.Token == "" {
		t.Fatal("Token이 비어 있습니다")
	}
	if token.ChannelID != "room-1" {
		t.Errorf("ChannelID = %q, want %q", token.ChannelID, "room-1")
	}
	if token.UID != "user-1" {
		t.Errorf("UID = %q, want %q", token.UID, "user-1")
	}
}

func TestAgoraProvider_GenerateToken_EmptyAppID(t *testing.T) {
	provider := webrtc.NewAgoraProvider(webrtc.AgoraConfig{
		AppID: "",
	})
	ctx := context.Background()

	_, err := provider.GenerateToken(ctx, "room-1", "user-1", webrtc.RolePatient)
	if err == nil {
		t.Fatal("빈 AppID에 에러가 반환되어야 합니다")
	}
}

func TestAgoraProvider_GenerateToken_EmptyAppCert(t *testing.T) {
	provider := webrtc.NewAgoraProvider(webrtc.AgoraConfig{
		AppID:   "test-app-id",
		AppCert: "",
	})
	ctx := context.Background()

	_, err := provider.GenerateToken(ctx, "room-1", "user-1", webrtc.RolePatient)
	if err == nil {
		t.Fatal("빈 AppCert에 에러가 반환되어야 합니다")
	}
}

func TestAgoraProvider_GenerateToken_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	provider := webrtc.NewAgoraProvider(webrtc.AgoraConfig{
		AppID:      "test-app-id",
		AppCert:    "test-cert",
		CustomerID: "cust",
		Secret:     "sec",
		BaseURL:    server.URL,
	})
	ctx := context.Background()

	_, err := provider.GenerateToken(ctx, "room-1", "user-1", webrtc.RolePatient)
	if err == nil {
		t.Fatal("HTTP 503에 에러가 반환되어야 합니다")
	}
}

func TestAgoraProvider_EndRoom(t *testing.T) {
	provider := webrtc.NewAgoraProvider(webrtc.AgoraConfig{
		AppID: "test-app-id",
	})
	ctx := context.Background()

	err := provider.EndRoom(ctx, "room-1")
	if err != nil {
		t.Fatalf("EndRoom 에러: %v", err)
	}
}

func TestAgoraProvider_ProviderName(t *testing.T) {
	provider := webrtc.NewAgoraProvider(webrtc.AgoraConfig{
		AppID: "test-app-id",
	})
	if provider.ProviderName() != "agora" {
		t.Errorf("ProviderName = %q, want %q", provider.ProviderName(), "agora")
	}
}
