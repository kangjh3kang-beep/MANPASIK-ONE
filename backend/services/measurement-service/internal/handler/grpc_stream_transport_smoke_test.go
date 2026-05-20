package handler

import (
	"bytes"
	"context"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/manpasik/backend/services/measurement-service/internal/repository/memory"
	"github.com/manpasik/backend/services/measurement-service/internal/service"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestStreamMeasurementOverGRPCTransportStoresAndResponds(t *testing.T) {
	handler, measureRepo, vectorRepo := newStreamMeasurementHandler(&service.MeasurementSession{
		ID:          "session-transport-1",
		DeviceID:    "device-transport-1",
		CartridgeID: "glucose",
		UserID:      "user-transport-1",
		StartedAt:   time.Now().UTC(),
		Status:      "active",
	})

	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	v1.RegisterMeasurementServiceServer(server, handler)
	go func() {
		if err := server.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			t.Errorf("bufconn gRPC server failed: %v", err)
		}
	}()
	t.Cleanup(func() {
		server.Stop()
		_ = listener.Close()
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return listener.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("dial bufconn: %v", err)
	}
	t.Cleanup(func() {
		_ = conn.Close()
	})

	stream, err := v1.NewMeasurementServiceClient(conn).StreamMeasurement(ctx)
	if err != nil {
		t.Fatalf("open StreamMeasurement: %v", err)
	}
	if err := stream.Send(&v1.MeasurementData{
		SessionId:   "session-transport-1",
		RawChannels: []float64{4, 5, 6},
		Differential: &v1.DifferentialCorrection{
			SDet:       120,
			SRef:       20,
			Alpha:      0.95,
			SCorrected: 101,
		},
		EnvMeta: &v1.EnvironmentMeta{
			TempC:       25.5,
			HumidityPct: 46,
			PressureKpa: 101.4,
		},
	}); err != nil {
		t.Fatalf("send measurement frame: %v", err)
	}
	if err := stream.CloseSend(); err != nil {
		t.Fatalf("close send: %v", err)
	}

	response, err := stream.Recv()
	if err != nil {
		t.Fatalf("recv measurement response: %v", err)
	}
	if response.SessionId != "session-transport-1" ||
		response.PrimaryValue != 101 ||
		response.Unit != "mg/dL" {
		t.Fatalf("response mismatch: %+v", response)
	}
	assertAssayDerivedConfidence(t, response.Confidence)
	assertFloat32Slice(t, response.FingerprintVector, []float32{4, 5, 6})
	if response.ProcessedAt == nil {
		t.Fatal("response ProcessedAt is nil")
	}
	if _, err := stream.Recv(); err != io.EOF {
		t.Fatalf("second recv err = %v, want EOF", err)
	}

	if len(measureRepo.stored) != 1 {
		t.Fatalf("stored measurement count = %d, want 1", len(measureRepo.stored))
	}
	stored := measureRepo.stored[0]
	if stored.SessionID != "session-transport-1" ||
		stored.DeviceID != "device-transport-1" ||
		stored.UserID != "user-transport-1" ||
		stored.CartridgeType != "glucose" {
		t.Fatalf("stored metadata mismatch: %+v", stored)
	}
	if stored.PrimaryValue != 101 || stored.SCorrected != 101 {
		t.Fatalf("stored corrected values mismatch: primary=%f corrected=%f", stored.PrimaryValue, stored.SCorrected)
	}
	assertFloat64Slice(t, stored.RawChannels, []float64{4, 5, 6})
	assertFloat32Slice(t, stored.FingerprintVector, []float32{4, 5, 6})
	assertFloat32Slice(t, vectorRepo.vectors["session-transport-1"], []float32{4, 5, 6})
}

func TestStreamMeasurementOverTCPLoopbackHandlesMultipleFrames(t *testing.T) {
	handler, measureRepo, vectorRepo := newStreamMeasurementHandler(&service.MeasurementSession{
		ID:          "session-tcp-1",
		DeviceID:    "device-tcp-1",
		CartridgeID: "glucose",
		UserID:      "user-tcp-1",
		StartedAt:   time.Now().UTC(),
		Status:      "active",
	})

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen tcp loopback: %v", err)
	}
	server := grpc.NewServer()
	v1.RegisterMeasurementServiceServer(server, handler)
	go func() {
		if err := server.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			t.Errorf("tcp gRPC server failed: %v", err)
		}
	}()
	t.Cleanup(func() {
		server.Stop()
		_ = listener.Close()
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(
		ctx,
		listener.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Fatalf("dial tcp loopback: %v", err)
	}
	t.Cleanup(func() {
		_ = conn.Close()
	})

	stream, err := v1.NewMeasurementServiceClient(conn).StreamMeasurement(ctx)
	if err != nil {
		t.Fatalf("open StreamMeasurement over tcp: %v", err)
	}
	frames := []*v1.MeasurementData{
		{
			SessionId:   "session-tcp-1",
			RawChannels: []float64{7, 8},
			Differential: &v1.DifferentialCorrection{
				SDet:       90,
				SRef:       10,
				Alpha:      0.95,
				SCorrected: 80.5,
			},
			EnvMeta: &v1.EnvironmentMeta{
				TempC:       26,
				HumidityPct: 47,
			},
		},
		{
			SessionId:   "session-tcp-1",
			RawChannels: []float64{9, 10, 11},
			Differential: &v1.DifferentialCorrection{
				SDet:       110,
				SRef:       12,
				Alpha:      0.95,
				SCorrected: 98.6,
			},
			EnvMeta: &v1.EnvironmentMeta{
				TempC:       26.5,
				HumidityPct: 48,
			},
		},
	}
	for _, frame := range frames {
		if err := stream.Send(frame); err != nil {
			t.Fatalf("send tcp measurement frame: %v", err)
		}
	}
	if err := stream.CloseSend(); err != nil {
		t.Fatalf("close tcp send: %v", err)
	}

	for i, want := range []struct {
		primary     float64
		fingerprint []float32
	}{
		{primary: 80.5, fingerprint: []float32{7, 8}},
		{primary: 98.6, fingerprint: []float32{9, 10, 11}},
	} {
		response, err := stream.Recv()
		if err != nil {
			t.Fatalf("recv tcp measurement response %d: %v", i, err)
		}
		if response.SessionId != "session-tcp-1" ||
			response.PrimaryValue != want.primary ||
			response.Unit != "mg/dL" {
			t.Fatalf("tcp response %d mismatch: %+v", i, response)
		}
		assertAssayDerivedConfidence(t, response.Confidence)
		assertFloat32Slice(t, response.FingerprintVector, want.fingerprint)
		if response.ProcessedAt == nil {
			t.Fatalf("tcp response %d ProcessedAt is nil", i)
		}
	}
	if _, err := stream.Recv(); err != io.EOF {
		t.Fatalf("tcp final recv err = %v, want EOF", err)
	}

	if len(measureRepo.stored) != 2 {
		t.Fatalf("stored tcp measurement count = %d, want 2", len(measureRepo.stored))
	}
	first := measureRepo.stored[0]
	second := measureRepo.stored[1]
	if first.PrimaryValue != 80.5 || first.SCorrected != 80.5 {
		t.Fatalf("first stored values mismatch: primary=%f corrected=%f", first.PrimaryValue, first.SCorrected)
	}
	if second.PrimaryValue != 98.6 || second.SCorrected != 98.6 {
		t.Fatalf("second stored values mismatch: primary=%f corrected=%f", second.PrimaryValue, second.SCorrected)
	}
	if first.DeviceID != "device-tcp-1" || first.UserID != "user-tcp-1" || first.CartridgeType != "glucose" {
		t.Fatalf("first stored metadata mismatch: %+v", first)
	}
	if second.DeviceID != "device-tcp-1" || second.UserID != "user-tcp-1" || second.CartridgeType != "glucose" {
		t.Fatalf("second stored metadata mismatch: %+v", second)
	}
	assertFloat64Slice(t, first.RawChannels, []float64{7, 8})
	assertFloat64Slice(t, second.RawChannels, []float64{9, 10, 11})
	assertFloat32Slice(t, vectorRepo.vectors["session-tcp-1"], []float32{9, 10, 11})
}

func TestStreamMeasurementAgainstHelperProcess(t *testing.T) {
	addr := startMeasurementServiceHelperProcess(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Fatalf("dial helper process measurement-service: %v", err)
	}
	t.Cleanup(func() {
		_ = conn.Close()
	})

	client := v1.NewMeasurementServiceClient(conn)
	session, err := client.StartSession(ctx, &v1.StartSessionRequest{
		DeviceId:    "device-process-1",
		CartridgeId: "glucose",
		UserId:      "user-process-1",
	})
	if err != nil {
		t.Fatalf("StartSession via helper process: %v", err)
	}
	if session.SessionId == "" || session.StartedAt == nil {
		t.Fatalf("invalid StartSession response: %+v", session)
	}

	stream, err := client.StreamMeasurement(ctx)
	if err != nil {
		t.Fatalf("open StreamMeasurement via helper process: %v", err)
	}
	if err := stream.Send(&v1.MeasurementData{
		SessionId:   session.SessionId,
		RawChannels: []float64{12, 13, 14},
		Differential: &v1.DifferentialCorrection{
			SDet:       140,
			SRef:       16,
			Alpha:      0.95,
			SCorrected: 124.8,
		},
		EnvMeta: &v1.EnvironmentMeta{
			TempC:       27.5,
			HumidityPct: 49,
		},
	}); err != nil {
		t.Fatalf("send helper process measurement frame: %v", err)
	}
	if err := stream.CloseSend(); err != nil {
		t.Fatalf("close helper process stream send: %v", err)
	}
	response, err := stream.Recv()
	if err != nil {
		t.Fatalf("recv helper process measurement response: %v", err)
	}
	if response.SessionId != session.SessionId ||
		response.PrimaryValue != 124.8 ||
		response.Unit != "mg/dL" {
		t.Fatalf("helper process response mismatch: %+v", response)
	}
	assertAssayDerivedConfidence(t, response.Confidence)
	assertFloat32Slice(t, response.FingerprintVector, []float32{12, 13, 14})
	if response.ProcessedAt == nil {
		t.Fatal("helper process response ProcessedAt is nil")
	}
	if _, err := stream.Recv(); err != io.EOF {
		t.Fatalf("helper process final recv err = %v, want EOF", err)
	}

	history, err := client.GetMeasurementHistory(ctx, &v1.GetHistoryRequest{
		UserId: "user-process-1",
		Limit:  10,
	})
	if err != nil {
		t.Fatalf("GetMeasurementHistory via helper process: %v", err)
	}
	if history.TotalCount != 1 || len(history.Measurements) != 1 {
		t.Fatalf("history count mismatch: total=%d len=%d", history.TotalCount, len(history.Measurements))
	}
	summary := history.Measurements[0]
	if summary.SessionId != session.SessionId ||
		summary.CartridgeType != "glucose" ||
		summary.PrimaryValue != 124.8 ||
		summary.Unit != "mg/dL" ||
		summary.MeasuredAt == nil {
		t.Fatalf("history summary mismatch: %+v", summary)
	}

	ended, err := client.EndSession(ctx, &v1.EndSessionRequest{
		SessionId: session.SessionId,
	})
	if err != nil {
		t.Fatalf("EndSession via helper process: %v", err)
	}
	if ended.SessionId != session.SessionId || ended.EndedAt == nil {
		t.Fatalf("invalid EndSession response: %+v", ended)
	}
}

func TestMeasurementServiceProcessHelper(t *testing.T) {
	if os.Getenv("MANPASIK_MEASUREMENT_PROCESS_HELPER") != "1" {
		return
	}

	addrFile := os.Getenv("MANPASIK_MEASUREMENT_PROCESS_ADDR_FILE")
	if addrFile == "" {
		t.Fatal("MANPASIK_MEASUREMENT_PROCESS_ADDR_FILE is required")
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("helper listen tcp loopback: %v", err)
	}

	measureSvc := service.NewMeasurementService(
		zap.NewNop(),
		memory.NewSessionRepository(),
		memory.NewMeasurementRepository(),
		memory.NewVectorRepository(),
		memory.NewEventPublisher(),
	)
	server := grpc.NewServer()
	v1.RegisterMeasurementServiceServer(server, NewMeasurementHandler(measureSvc, zap.NewNop()))
	go func() {
		if err := server.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			t.Errorf("helper grpc server failed: %v", err)
		}
	}()
	if err := os.WriteFile(addrFile, []byte(listener.Addr().String()), 0o600); err != nil {
		server.Stop()
		t.Fatalf("write helper address file: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	server.GracefulStop()
}

func startMeasurementServiceHelperProcess(t *testing.T) string {
	t.Helper()

	addrFile := filepath.Join(t.TempDir(), "measurement-helper.addr")
	cmd := exec.Command(os.Args[0], "-test.run=TestMeasurementServiceProcessHelper", "-test.v")
	cmd.Env = append(
		os.Environ(),
		"MANPASIK_MEASUREMENT_PROCESS_HELPER=1",
		"MANPASIK_MEASUREMENT_PROCESS_ADDR_FILE="+addrFile,
	)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Start(); err != nil {
		t.Fatalf("start helper process: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	t.Cleanup(func() {
		if cmd.Process == nil {
			return
		}
		_ = cmd.Process.Signal(os.Interrupt)
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			_ = cmd.Process.Kill()
			<-done
		}
	})

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		select {
		case err := <-done:
			t.Fatalf("helper process exited before ready: %v\n%s", err, output.String())
		default:
		}

		raw, err := os.ReadFile(addrFile)
		addr := strings.TrimSpace(string(raw))
		if err == nil && addr != "" {
			return addr
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("helper process did not publish address\n%s", output.String())
	return ""
}

func assertAssayDerivedConfidence(t *testing.T, confidence float64) {
	t.Helper()
	if confidence <= 0.90 || confidence >= 0.95 {
		t.Fatalf("confidence = %f, want assay-derived partial-channel confidence", confidence)
	}
}
