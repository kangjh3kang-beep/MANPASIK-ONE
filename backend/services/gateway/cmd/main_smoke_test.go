package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestGatewayAuthLiveBinarySmoke(t *testing.T) {
	root := repositoryRoot(t)
	authBinary := buildGoBinary(t, filepath.Join(root, "backend/services/auth-service/cmd"), "auth-service-live-smoke")
	gatewayBinary := buildGoBinary(t, filepath.Join(root, "backend/services/gateway/cmd"), "gateway-live-smoke")

	authGRPCAddr := reserveLoopbackAddress(t)
	authMetricsAddr := reserveLoopbackAddress(t)
	gatewayHTTPAddr := reserveLoopbackAddress(t)

	authCmd, authDone, authOutput := startSmokeProcess(t, authBinary, authServiceSmokeEnv(authGRPCAddr, authMetricsAddr))
	t.Cleanup(func() {
		stopSmokeProcess(t, authCmd, authDone)
	})
	authConn := waitForSmokeGRPCServing(t, "auth-service", authDone, authOutput, authGRPCAddr)
	t.Cleanup(func() {
		_ = authConn.Close()
	})

	gatewayCmd, gatewayDone, gatewayOutput := startSmokeProcess(t, gatewayBinary, gatewaySmokeEnv(gatewayHTTPAddr, authGRPCAddr))
	t.Cleanup(func() {
		stopSmokeProcess(t, gatewayCmd, gatewayDone)
	})
	waitForGatewayHTTPHealth(t, gatewayDone, gatewayOutput, gatewayHTTPAddr)

	runGatewayAuthRESTLifecycleSmoke(t, gatewayHTTPAddr)
}

func runGatewayAuthRESTLifecycleSmoke(t *testing.T, gatewayHTTPAddr string) {
	t.Helper()

	client := &http.Client{Timeout: 2 * time.Second}
	baseURL := "http://" + gatewayHTTPAddr + "/api/v1/auth"
	email := fmt.Sprintf("gateway-live-%d@manpasik.com", time.Now().UnixNano())
	password := "Pass123!Live"

	registerPayload := map[string]any{
		"email":        email,
		"password":     password,
		"display_name": "Gateway Live",
	}
	registerResp := postJSON(t, client, baseURL+"/register", registerPayload, http.StatusCreated)
	userID := readStringField(t, registerResp, "user_id", "userId")
	if userID == "" {
		t.Fatalf("register response user_id is empty: %+v", registerResp)
	}

	loginResp := postJSON(t, client, baseURL+"/login", map[string]any{
		"email":    email,
		"password": password,
	}, http.StatusOK)
	if got := readStringField(t, loginResp, "user_id", "userId"); got != userID {
		t.Fatalf("login user_id = %q, want registered user_id %q (payload=%+v)", got, userID, loginResp)
	}
	accessToken := readStringField(t, loginResp, "access_token", "accessToken")
	refreshToken := readStringField(t, loginResp, "refresh_token", "refreshToken")
	if accessToken == "" || refreshToken == "" {
		t.Fatalf("login response missing tokens: %+v", loginResp)
	}

	refreshResp := postJSON(t, client, baseURL+"/refresh", map[string]any{
		"refresh_token": refreshToken,
	}, http.StatusOK)
	if got := readStringField(t, refreshResp, "user_id", "userId"); got != userID {
		t.Fatalf("refresh user_id = %q, want registered user_id %q (payload=%+v)", got, userID, refreshResp)
	}
	if readStringField(t, refreshResp, "access_token", "accessToken") == "" ||
		readStringField(t, refreshResp, "refresh_token", "refreshToken") == "" {
		t.Fatalf("refresh response missing tokens: %+v", refreshResp)
	}
}

func postJSON(
	t *testing.T,
	client *http.Client,
	url string,
	payload map[string]any,
	wantStatus int,
) map[string]any {
	t.Helper()

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal request payload: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("create POST %s: %v", url, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("POST %s: %v", url, err)
	}
	defer resp.Body.Close()

	var decoded map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatalf("decode POST %s response: %v", url, err)
	}
	if resp.StatusCode != wantStatus {
		t.Fatalf("POST %s status = %d, want %d, body=%+v", url, resp.StatusCode, wantStatus, decoded)
	}
	return decoded
}

func readStringField(t *testing.T, source map[string]any, keys ...string) string {
	t.Helper()

	for _, key := range keys {
		if value, ok := source[key].(string); ok {
			return value
		}
	}
	return ""
}

func buildGoBinary(t *testing.T, packageDir string, binaryName string) string {
	t.Helper()

	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	binary := filepath.Join(t.TempDir(), binaryName)
	build := exec.Command(goTool(t), "build", "-o", binary, ".")
	build.Dir = packageDir
	var output bytes.Buffer
	build.Stdout = &output
	build.Stderr = &output
	if err := build.Run(); err != nil {
		t.Fatalf("build %s: %v\n%s", packageDir, err, output.String())
	}
	return binary
}

func repositoryRoot(t *testing.T) string {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve gateway cmd package path")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(filename), "../../../.."))
}

func goTool(t *testing.T) string {
	t.Helper()

	if configured := os.Getenv("MANPASIK_GO_BINARY"); configured != "" {
		return configured
	}
	if goroot := runtime.GOROOT(); goroot != "" {
		candidate := filepath.Join(goroot, "bin", "go")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	if found, err := exec.LookPath("go"); err == nil {
		return found
	}
	t.Fatal("go tool not found; set MANPASIK_GO_BINARY")
	return ""
}

func reserveLoopbackAddress(t *testing.T) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve loopback address: %v", err)
	}
	addr := listener.Addr().String()
	if err := listener.Close(); err != nil {
		t.Fatalf("release loopback address: %v", err)
	}
	return addr
}

func startSmokeProcess(t *testing.T, binary string, env []string) (*exec.Cmd, chan error, *bytes.Buffer) {
	t.Helper()

	cmd := exec.Command(binary)
	cmd.Env = env
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Start(); err != nil {
		t.Fatalf("start %s: %v", binary, err)
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	return cmd, done, &output
}

func authServiceSmokeEnv(grpcAddr, metricsAddr string) []string {
	return baseSmokeEnv([]string{
		"VERSION=auth-live-smoke",
		"GRPC_PORT=" + grpcAddr,
		"METRICS_PORT=" + metricsAddr,
		"SHUTDOWN_TIMEOUT_SECONDS=1",
		"TENANCY_ENFORCED=false",
		"JWT_SECRET=gateway-live-smoke-secret-32-bytes",
		"JWT_ISSUER=manpasik-live-smoke",
	})
}

func gatewaySmokeEnv(httpAddr, authGRPCAddr string) []string {
	return baseSmokeEnv([]string{
		"VERSION=gateway-live-smoke",
		"HTTP_PORT=" + httpAddr,
		"AUTH_SERVICE_ADDR=" + authGRPCAddr,
		"SHUTDOWN_TIMEOUT_SECONDS=1",
		"TENANCY_ENFORCED=false",
		"JWT_SECRET=gateway-live-smoke-secret-32-bytes",
		"JWT_ISSUER=manpasik-live-smoke",
	})
}

func baseSmokeEnv(values []string) []string {
	env := make([]string, 0, len(values)+4)
	env = append(env, values...)
	for _, key := range []string{"PATH", "HOME", "TMPDIR", "USER"} {
		if value := os.Getenv(key); value != "" {
			env = append(env, key+"="+value)
		}
	}
	return env
}

func waitForSmokeGRPCServing(
	t *testing.T,
	service string,
	done chan error,
	output *bytes.Buffer,
	addr string,
) *grpc.ClientConn {
	t.Helper()

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		select {
		case err := <-done:
			done <- err
			t.Fatalf("%s binary exited before gRPC health ready: %v\n%s", service, err, output.String())
		default:
		}

		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		conn, err := grpc.DialContext(
			ctx,
			addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		cancel()
		if err == nil {
			checkCtx, checkCancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
			resp, checkErr := healthpb.NewHealthClient(conn).Check(checkCtx, &healthpb.HealthCheckRequest{
				Service: service,
			})
			checkCancel()
			if checkErr == nil && resp.GetStatus() == healthpb.HealthCheckResponse_SERVING {
				return conn
			}
			_ = conn.Close()
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("%s gRPC health not ready at %s\n%s", service, addr, output.String())
	return nil
}

func waitForGatewayHTTPHealth(t *testing.T, done chan error, output *bytes.Buffer, addr string) {
	t.Helper()

	client := &http.Client{Timeout: 300 * time.Millisecond}
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		select {
		case err := <-done:
			done <- err
			t.Fatalf("gateway binary exited before HTTP health ready: %v\n%s", err, output.String())
		default:
		}

		resp, err := client.Get("http://" + addr + "/health")
		if err == nil {
			var payload struct {
				Status string `json:"status"`
			}
			decodeErr := json.NewDecoder(resp.Body).Decode(&payload)
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK && decodeErr == nil && payload.Status == "ok" {
				return
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("gateway HTTP health not ready at %s\n%s", addr, output.String())
}

func stopSmokeProcess(t *testing.T, cmd *exec.Cmd, done chan error) {
	t.Helper()

	if cmd.Process == nil {
		return
	}
	select {
	case <-done:
		return
	default:
	}
	_ = cmd.Process.Signal(syscall.SIGTERM)
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		_ = cmd.Process.Kill()
		select {
		case <-done:
		case <-time.After(time.Second):
		}
	}
}

func TestGatewayAuthExternalEndpointSmoke(t *testing.T) {
	gatewayHTTPAddr := os.Getenv("MANPASIK_GATEWAY_AUTH_SMOKE_HTTP_ADDR")
	if gatewayHTTPAddr == "" {
		t.Skip("MANPASIK_GATEWAY_AUTH_SMOKE_HTTP_ADDR is not set")
	}
	if strings.HasPrefix(gatewayHTTPAddr, "http://") || strings.HasPrefix(gatewayHTTPAddr, "https://") {
		gatewayHTTPAddr = strings.TrimPrefix(strings.TrimPrefix(gatewayHTTPAddr, "http://"), "https://")
	}
	runGatewayAuthRESTLifecycleSmoke(t, gatewayHTTPAddr)
}
