package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
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

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestMeasurementServiceBinarySmoke(t *testing.T) {
	binary := buildMeasurementServiceBinary(t)
	grpcAddr := reserveLoopbackAddress(t)
	httpAddr := reserveLoopbackAddress(t)

	cmd := exec.Command(binary)
	cmd.Env = measurementServiceSmokeEnv(grpcAddr, httpAddr)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Start(); err != nil {
		t.Fatalf("start measurement-service binary: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	t.Cleanup(func() {
		stopMeasurementServiceProcess(t, cmd, done)
	})

	waitForHTTPHealth(t, done, &output, httpAddr)
	conn := waitForGRPCServing(t, done, &output, grpcAddr)
	t.Cleanup(func() {
		_ = conn.Close()
	})

	runMeasurementServiceLifecycleSmoke(t, conn, "binary", "device-binary-1", "user-binary-1", []float64{21, 22, 23}, 161)
}

func TestMeasurementServiceExternalEndpointSmoke(t *testing.T) {
	grpcAddr := os.Getenv("MANPASIK_MEASUREMENT_SERVICE_SMOKE_ADDR")
	if grpcAddr == "" {
		t.Skip("MANPASIK_MEASUREMENT_SERVICE_SMOKE_ADDR is not set")
	}
	httpAddr := os.Getenv("MANPASIK_MEASUREMENT_SERVICE_SMOKE_HTTP_ADDR")
	if httpAddr != "" {
		waitForExternalHTTPHealth(t, httpAddr, os.Getenv("MANPASIK_MEASUREMENT_SERVICE_SMOKE_VERSION"))
	}

	conn := waitForExternalGRPCServing(t, grpcAddr)
	t.Cleanup(func() {
		_ = conn.Close()
	})

	runMeasurementServiceLifecycleSmoke(t, conn, "external", "device-compose-1", "user-compose-1", []float64{31, 32, 33}, 251)
}

func runMeasurementServiceLifecycleSmoke(
	t *testing.T,
	conn *grpc.ClientConn,
	label string,
	deviceID string,
	userID string,
	rawChannels []float64,
	corrected float64,
) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client := v1.NewMeasurementServiceClient(conn)

	session, err := client.StartSession(ctx, &v1.StartSessionRequest{
		DeviceId:    deviceID,
		CartridgeId: "glucose",
		UserId:      userID,
	})
	if err != nil {
		t.Fatalf("StartSession against %s endpoint: %v", label, err)
	}
	if session.SessionId == "" || session.StartedAt == nil {
		t.Fatalf("invalid %s StartSession response: %+v", label, session)
	}

	stream, err := client.StreamMeasurement(ctx)
	if err != nil {
		t.Fatalf("open StreamMeasurement against %s endpoint: %v", label, err)
	}
	if err := stream.Send(&v1.MeasurementData{
		SessionId:   session.SessionId,
		RawChannels: rawChannels,
		Differential: &v1.DifferentialCorrection{
			SDet:       180,
			SRef:       20,
			Alpha:      0.95,
			SCorrected: corrected,
		},
		EnvMeta: &v1.EnvironmentMeta{
			TempC:       28.2,
			HumidityPct: 50.5,
		},
	}); err != nil {
		t.Fatalf("send %s measurement frame: %v", label, err)
	}
	if err := stream.CloseSend(); err != nil {
		t.Fatalf("close %s stream send: %v", label, err)
	}

	response, err := stream.Recv()
	if err != nil {
		t.Fatalf("recv %s measurement response: %v", label, err)
	}
	if response.SessionId != session.SessionId ||
		response.PrimaryValue != corrected ||
		response.Unit != "mg/dL" ||
		response.Confidence != 0.95 {
		t.Fatalf("%s measurement response mismatch: %+v", label, response)
	}
	assertFloat32Slice(t, response.FingerprintVector, float64ToFloat32(rawChannels))
	if response.ProcessedAt == nil {
		t.Fatalf("%s measurement response ProcessedAt is nil", label)
	}
	if _, err := stream.Recv(); err != io.EOF {
		t.Fatalf("%s final recv err = %v, want EOF", label, err)
	}

	history, err := client.GetMeasurementHistory(ctx, &v1.GetHistoryRequest{
		UserId: userID,
		Limit:  10,
	})
	if err != nil {
		t.Fatalf("GetMeasurementHistory against %s endpoint: %v", label, err)
	}
	if history.TotalCount != 1 || len(history.Measurements) != 1 {
		t.Fatalf("%s history count mismatch: total=%d len=%d", label, history.TotalCount, len(history.Measurements))
	}
	summary := history.Measurements[0]
	if summary.SessionId != session.SessionId ||
		summary.CartridgeType != "glucose" ||
		summary.PrimaryValue != corrected ||
		summary.Unit != "mg/dL" ||
		summary.MeasuredAt == nil {
		t.Fatalf("%s history summary mismatch: %+v", label, summary)
	}

	ended, err := client.EndSession(ctx, &v1.EndSessionRequest{
		SessionId: session.SessionId,
	})
	if err != nil {
		t.Fatalf("EndSession against %s endpoint: %v", label, err)
	}
	if ended.SessionId != session.SessionId || ended.EndedAt == nil {
		t.Fatalf("invalid %s EndSession response: %+v", label, ended)
	}
}

func buildMeasurementServiceBinary(t *testing.T) string {
	t.Helper()

	binaryName := "measurement-service-smoke"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	binary := filepath.Join(t.TempDir(), binaryName)
	build := exec.Command(goTool(t), "build", "-o", binary, ".")
	build.Dir = currentPackageDir(t)
	var output bytes.Buffer
	build.Stdout = &output
	build.Stderr = &output
	if err := build.Run(); err != nil {
		t.Fatalf("build measurement-service binary: %v\n%s", err, output.String())
	}
	return binary
}

func currentPackageDir(t *testing.T) string {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve current package dir")
	}
	return filepath.Dir(filename)
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

func measurementServiceSmokeEnv(grpcAddr, httpAddr string) []string {
	env := []string{
		"VERSION=smoke-test",
		"GRPC_PORT=" + grpcAddr,
		"HTTP_PORT=" + httpAddr,
		"SHUTDOWN_TIMEOUT_SECONDS=1",
		"TENANCY_ENFORCED=false",
		"LOG_LEVEL=debug",
	}
	for _, key := range []string{"PATH", "HOME", "TMPDIR", "USER"} {
		if value := os.Getenv(key); value != "" {
			env = append(env, key+"="+value)
		}
	}
	return env
}

func waitForHTTPHealth(t *testing.T, done chan error, output *bytes.Buffer, addr string) {
	t.Helper()

	client := &http.Client{Timeout: 300 * time.Millisecond}
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		select {
		case err := <-done:
			done <- err
			t.Fatalf("measurement-service binary exited before HTTP health ready: %v\n%s", err, output.String())
		default:
		}

		resp, err := client.Get("http://" + addr + "/health")
		if err == nil {
			var payload struct {
				Status  string `json:"status"`
				Service string `json:"service"`
				Version string `json:"version"`
			}
			decodeErr := json.NewDecoder(resp.Body).Decode(&payload)
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK &&
				decodeErr == nil &&
				payload.Status == "healthy" &&
				payload.Service == serviceName &&
				payload.Version == "smoke-test" {
				return
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("measurement-service binary HTTP health not ready at %s\n%s", addr, output.String())
}

func waitForGRPCServing(t *testing.T, done chan error, output *bytes.Buffer, addr string) *grpc.ClientConn {
	t.Helper()

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		select {
		case err := <-done:
			done <- err
			t.Fatalf("measurement-service binary exited before gRPC health ready: %v\n%s", err, output.String())
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
				Service: serviceName,
			})
			checkCancel()
			if checkErr == nil && resp.GetStatus() == healthpb.HealthCheckResponse_SERVING {
				return conn
			}
			_ = conn.Close()
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("measurement-service binary gRPC health not ready at %s\n%s", addr, output.String())
	return nil
}

func waitForExternalHTTPHealth(t *testing.T, addr string, expectedVersion string) {
	t.Helper()

	healthURL := addr
	if !strings.HasPrefix(healthURL, "http://") && !strings.HasPrefix(healthURL, "https://") {
		healthURL = "http://" + healthURL + "/health"
	}

	client := &http.Client{Timeout: 500 * time.Millisecond}
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := client.Get(healthURL)
		if err == nil {
			var payload struct {
				Status  string `json:"status"`
				Service string `json:"service"`
				Version string `json:"version"`
			}
			decodeErr := json.NewDecoder(resp.Body).Decode(&payload)
			_ = resp.Body.Close()
			versionMatches := expectedVersion == "" || payload.Version == expectedVersion
			if resp.StatusCode == http.StatusOK &&
				decodeErr == nil &&
				payload.Status == "healthy" &&
				payload.Service == serviceName &&
				versionMatches {
				return
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatalf("external measurement-service HTTP health not ready at %s", healthURL)
}

func waitForExternalGRPCServing(t *testing.T, addr string) *grpc.ClientConn {
	t.Helper()

	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		conn, err := grpc.DialContext(
			ctx,
			addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		cancel()
		if err == nil {
			checkCtx, checkCancel := context.WithTimeout(context.Background(), time.Second)
			resp, checkErr := healthpb.NewHealthClient(conn).Check(checkCtx, &healthpb.HealthCheckRequest{
				Service: serviceName,
			})
			checkCancel()
			if checkErr == nil && resp.GetStatus() == healthpb.HealthCheckResponse_SERVING {
				return conn
			}
			_ = conn.Close()
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatalf("external measurement-service gRPC health not ready at %s", addr)
	return nil
}

func stopMeasurementServiceProcess(t *testing.T, cmd *exec.Cmd, done chan error) {
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

func float64ToFloat32(values []float64) []float32 {
	converted := make([]float32, 0, len(values))
	for _, value := range values {
		converted = append(converted, float32(value))
	}
	return converted
}

func assertFloat32Slice(t *testing.T, got, want []float32) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("float32 slice length = %d, want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("float32 slice[%d] = %f, want %f (all=%v)", i, got[i], want[i], got)
		}
	}
}
