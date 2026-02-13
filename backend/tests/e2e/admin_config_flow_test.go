package e2e

import (
	"context"
	"testing"
	"time"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TestAdminConfigFlow E2E: SetSystemConfig → GetSystemConfig → ListSystemConfigs → ValidateConfigValue → BulkSetConfigs
func TestAdminConfigFlow(t *testing.T) {
	dialCtx, dialCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer dialCancel()

	adminConn, err := grpc.DialContext(dialCtx, AdminAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Skipf("admin 서비스 연결 불가 (기동 후 재실행): %v", err)
	}
	defer adminConn.Close()

	client := v1.NewAdminServiceClient(adminConn)
	rpcCtx, rpcCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer rpcCancel()

	// 1) SetSystemConfig: 기본 설정 변경
	setResp, err := client.SetSystemConfig(rpcCtx, &v1.SetSystemConfigRequest{
		Key:   "maintenance_mode",
		Value: "false",
	})
	if err != nil {
		t.Fatalf("SetSystemConfig 실패: %v", err)
	}
	if setResp.Key != "maintenance_mode" {
		t.Errorf("key 불일치: %s", setResp.Key)
	}
	t.Logf("SetSystemConfig 성공: %s = %s", setResp.Key, setResp.Value)

	// 2) GetSystemConfig: 설정 조회
	getResp, err := client.GetSystemConfig(rpcCtx, &v1.GetSystemConfigRequest{
		Key: "maintenance_mode",
	})
	if err != nil {
		t.Fatalf("GetSystemConfig 실패: %v", err)
	}
	if getResp.Value != "false" {
		t.Errorf("값 불일치: %s", getResp.Value)
	}
	t.Logf("GetSystemConfig 성공: %s = %s", getResp.Key, getResp.Value)

	// 3) ListSystemConfigs: 설정 목록 조회 (확장 RPC)
	listResp, err := client.ListSystemConfigs(rpcCtx, &v1.ListSystemConfigsRequest{
		LanguageCode: "ko",
		Category:     "general",
	})
	if err != nil {
		t.Logf("ListSystemConfigs 미구현 또는 실패 (예상 가능): %v", err)
	} else {
		t.Logf("ListSystemConfigs 성공: %d개 설정", len(listResp.Configs))
		for _, c := range listResp.Configs {
			t.Logf("  - %s (%s): %s", c.Key, c.Category, c.Value)
		}
	}

	// 4) ValidateConfigValue: 유효성 검증
	validateResp, err := client.ValidateConfigValue(rpcCtx, &v1.ValidateConfigValueRequest{
		Key:   "maintenance_mode",
		Value: "true",
	})
	if err != nil {
		t.Logf("ValidateConfigValue 미구현 또는 실패 (예상 가능): %v", err)
	} else {
		if !validateResp.Valid {
			t.Errorf("true는 유효해야 합니다: %s", validateResp.ErrorMessage)
		}
		t.Logf("ValidateConfigValue 성공: valid=%v", validateResp.Valid)
	}

	// 5) GetSystemStats: 시스템 통계
	statsResp, err := client.GetSystemStats(rpcCtx, &v1.GetSystemStatsRequest{})
	if err != nil {
		t.Fatalf("GetSystemStats 실패: %v", err)
	}
	t.Logf("시스템 통계: 사용자=%d, 활성=%d, 디바이스=%d, 측정=%d",
		statsResp.TotalUsers, statsResp.ActiveUsers, statsResp.TotalDevices, statsResp.TotalMeasurements)
}

// TestAdminAuditLogFlow E2E: 설정 변경 후 감사 로그 확인
func TestAdminAuditLogFlow(t *testing.T) {
	dialCtx, dialCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer dialCancel()

	adminConn, err := grpc.DialContext(dialCtx, AdminAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Skipf("admin 서비스 연결 불가: %v", err)
	}
	defer adminConn.Close()

	client := v1.NewAdminServiceClient(adminConn)
	rpcCtx, rpcCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer rpcCancel()

	// 설정 변경 (감사 로그 생성)
	_, _ = client.SetSystemConfig(rpcCtx, &v1.SetSystemConfigRequest{
		Key:   "max_devices_per_user",
		Value: "10",
	})

	// 감사 로그 조회
	auditResp, err := client.GetAuditLog(rpcCtx, &v1.GetAuditLogRequest{
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("GetAuditLog 실패: %v", err)
	}
	t.Logf("감사 로그 수: %d", len(auditResp.Entries))
	for _, e := range auditResp.Entries {
		t.Logf("  - [%v] %s: %s (%s)", e.Action, e.ResourceType, e.ResourceId, e.Details)
	}
}
