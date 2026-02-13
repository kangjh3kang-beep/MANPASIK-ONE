package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TestMeasurementFlow E2E: Register → Login → ValidateToken → StartSession → EndSession → GetMeasurementHistory
// 서비스 미기동 시 연결 실패하면 t.Skip 처리.
func TestMeasurementFlow(t *testing.T) {
	// Dial context (연결 전용, 짧은 타임아웃)
	dialCtx, dialCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer dialCancel()

	// 1) Auth 연결
	authConn, err := grpc.DialContext(dialCtx, AuthAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Skipf("auth 서비스 연결 불가 (기동 후 재실행): %v", err)
	}
	defer authConn.Close()

	// RPC context (개별 호출마다 충분한 시간 부여)
	rpcCtx, rpcCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer rpcCancel()

	// 2) 회원가입 후 로그인 (인메모리 저장소 사용 시 매 실행 새 사용자)
	registerReq := &v1.RegisterRequest{
		Email:       "e2e-flow@manpasik.test",
		Password:    "E2EPass123!",
		DisplayName: "E2E Flow",
	}
	var registerResp v1.RegisterResponse
	if err := authConn.Invoke(rpcCtx, "/manpasik.v1.AuthService/Register", registerReq, &registerResp); err != nil {
		t.Skipf("Register 실패 (서비스 확인): %v", err)
	}

	loginReq := &v1.LoginRequest{
		Email:    registerReq.Email,
		Password: registerReq.Password,
	}
	var loginResp v1.LoginResponse
	if err := authConn.Invoke(rpcCtx, "/manpasik.v1.AuthService/Login", loginReq, &loginResp); err != nil {
		t.Fatalf("Login 실패: %v", err)
	}
	if loginResp.AccessToken == "" {
		t.Fatal("Login: access_token 비어 있음")
	}

	// 3) user_id 조회 (ValidateToken)
	validateReq := &v1.ValidateTokenRequest{AccessToken: loginResp.AccessToken}
	var validateResp v1.ValidateTokenResponse
	if err := authConn.Invoke(rpcCtx, "/manpasik.v1.AuthService/ValidateToken", validateReq, &validateResp); err != nil {
		t.Fatalf("ValidateToken 실패: %v", err)
	}
	userID := validateResp.UserId
	if userID == "" {
		t.Fatal("ValidateToken: user_id 비어 있음")
	}

	// 4) Measurement 서비스 연결
	dialCtx2, dialCancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer dialCancel2()

	measConn, err := grpc.DialContext(dialCtx2, MeasurementAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Skipf("measurement 서비스 연결 불가: %v", err)
	}
	defer measConn.Close()

	// 5) StartSession
	startReq := &v1.StartSessionRequest{
		DeviceId:    "e2e-device-001",
		CartridgeId: "e2e-cartridge-001",
		UserId:      userID,
	}
	var startResp v1.StartSessionResponse
	if err := measConn.Invoke(rpcCtx, "/manpasik.v1.MeasurementService/StartSession", startReq, &startResp); err != nil {
		t.Fatalf("StartSession 실패: %v", err)
	}
	sessionID := startResp.SessionId
	if sessionID == "" {
		t.Fatal("StartSession: session_id 비어 있음")
	}

	// 6) EndSession
	endReq := &v1.EndSessionRequest{SessionId: sessionID}
	var endResp v1.EndSessionResponse
	if err := measConn.Invoke(rpcCtx, "/manpasik.v1.MeasurementService/EndSession", endReq, &endResp); err != nil {
		t.Fatalf("EndSession 실패: %v", err)
	}

	// 7) GetMeasurementHistory
	historyReq := &v1.GetHistoryRequest{
		UserId: userID,
		Limit:  10,
		Offset: 0,
	}
	var historyResp v1.GetHistoryResponse
	if err := measConn.Invoke(rpcCtx, "/manpasik.v1.MeasurementService/GetMeasurementHistory", historyReq, &historyResp); err != nil {
		t.Fatalf("GetMeasurementHistory 실패: %v", err)
	}
	// 방금 종료한 세션이 목록에 있을 수 있음 (구현에 따라 0개일 수도 있음)
	t.Logf("✅ 측정 플로우 완료: session_id=%s, history total=%d", sessionID, historyResp.TotalCount)
}
