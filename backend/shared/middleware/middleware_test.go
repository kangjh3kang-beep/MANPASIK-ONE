package middleware

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// --- RBAC Tests ---

func TestRBACInterceptor_NilConfig(t *testing.T) {
	interceptor := RBACInterceptor(nil)
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	resp, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/test/Method"}, handler)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp != "ok" {
		t.Fatalf("expected 'ok', got %v", resp)
	}
}

func TestRBACInterceptor_UnrestrictedMethod(t *testing.T) {
	config := &RBACConfig{
		MethodRoles: map[string][]string{
			"/restricted/Method": {RoleMedicalStaff},
		},
	}
	interceptor := RBACInterceptor(config)
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	// Call an unrestricted method
	resp, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/unrestricted/Method"}, handler)
	if err != nil {
		t.Fatalf("expected no error for unrestricted method, got %v", err)
	}
	if resp != "ok" {
		t.Fatalf("expected 'ok', got %v", resp)
	}
}

func TestRBACInterceptor_AdminAlwaysAllowed(t *testing.T) {
	config := &RBACConfig{
		MethodRoles: map[string][]string{
			"/admin/OnlyMethod": {RoleMedicalStaff},
		},
	}
	interceptor := RBACInterceptor(config)
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	// Set admin role in context
	ctx := context.WithValue(context.Background(), UserRoleKey, RoleAdmin)

	resp, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/admin/OnlyMethod"}, handler)
	if err != nil {
		t.Fatalf("admin should always be allowed, got error: %v", err)
	}
	if resp != "ok" {
		t.Fatalf("expected 'ok', got %v", resp)
	}
}

func TestRBACInterceptor_AllowedRole(t *testing.T) {
	config := &RBACConfig{
		MethodRoles: map[string][]string{
			"/test/Method": {RoleMedicalStaff, RoleUser},
		},
	}
	interceptor := RBACInterceptor(config)
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	ctx := context.WithValue(context.Background(), UserRoleKey, RoleUser)

	resp, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test/Method"}, handler)
	if err != nil {
		t.Fatalf("user role should be allowed, got error: %v", err)
	}
	if resp != "ok" {
		t.Fatalf("expected 'ok', got %v", resp)
	}
}

func TestRBACInterceptor_UnauthorizedRole(t *testing.T) {
	config := &RBACConfig{
		MethodRoles: map[string][]string{
			"/test/Method": {RoleMedicalStaff},
		},
	}
	interceptor := RBACInterceptor(config)
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	ctx := context.WithValue(context.Background(), UserRoleKey, RoleUser)

	_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test/Method"}, handler)
	if err == nil {
		t.Fatal("expected permission denied error")
	}
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected gRPC status error, got %v", err)
	}
	if st.Code() != codes.PermissionDenied {
		t.Fatalf("expected PermissionDenied, got %v", st.Code())
	}
}

func TestRBACInterceptor_NoRoleInContext(t *testing.T) {
	config := &RBACConfig{
		MethodRoles: map[string][]string{
			"/test/Method": {RoleMedicalStaff},
		},
	}
	interceptor := RBACInterceptor(config)
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	// No role set in context
	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/test/Method"}, handler)
	if err == nil {
		t.Fatal("expected permission denied error when no role in context")
	}
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected gRPC status error, got %v", err)
	}
	if st.Code() != codes.PermissionDenied {
		t.Fatalf("expected PermissionDenied, got %v", st.Code())
	}
}

// --- Rate Limiter Tests ---

func TestRateLimiter_Allow(t *testing.T) {
	limiter := NewRateLimiter(5, time.Second, 5)

	// First 5 should be allowed
	for i := 0; i < 5; i++ {
		if !limiter.Allow("user1") {
			t.Errorf("request %d should be allowed", i+1)
		}
	}

	// 6th should be denied
	if limiter.Allow("user1") {
		t.Error("6th request should be denied")
	}

	// Different user should still be allowed
	if !limiter.Allow("user2") {
		t.Error("different user should be allowed")
	}
}

func TestRateLimiter_Refill(t *testing.T) {
	limiter := NewRateLimiter(1, 10*time.Millisecond, 1)

	if !limiter.Allow("user1") {
		t.Error("first request should be allowed")
	}
	if limiter.Allow("user1") {
		t.Error("second request should be denied")
	}

	time.Sleep(15 * time.Millisecond)

	if !limiter.Allow("user1") {
		t.Error("request after refill should be allowed")
	}
}

func TestRateLimitInterceptor_Allowed(t *testing.T) {
	limiter := NewRateLimiter(10, time.Second, 10)
	interceptor := RateLimitInterceptor(limiter)
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	resp, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/test/Method"}, handler)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp != "ok" {
		t.Fatalf("expected 'ok', got %v", resp)
	}
}

func TestRateLimitInterceptor_Exhausted(t *testing.T) {
	limiter := NewRateLimiter(1, time.Minute, 1)
	interceptor := RateLimitInterceptor(limiter)
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	// First request should pass
	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/test/Method"}, handler)
	if err != nil {
		t.Fatalf("first request should succeed, got %v", err)
	}

	// Second request should be rate limited
	_, err = interceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/test/Method"}, handler)
	if err == nil {
		t.Fatal("expected rate limit error")
	}
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected gRPC status error, got %v", err)
	}
	if st.Code() != codes.ResourceExhausted {
		t.Fatalf("expected ResourceExhausted, got %v", st.Code())
	}
}

func TestRateLimitInterceptor_WithUserID(t *testing.T) {
	limiter := NewRateLimiter(1, time.Minute, 1)
	interceptor := RateLimitInterceptor(limiter)
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	// User1 first request
	ctx1 := context.WithValue(context.Background(), UserIDKey, "user-1")
	_, err := interceptor(ctx1, nil, &grpc.UnaryServerInfo{FullMethod: "/test/Method"}, handler)
	if err != nil {
		t.Fatalf("user1 first request should succeed, got %v", err)
	}

	// User2 first request (different user should succeed)
	ctx2 := context.WithValue(context.Background(), UserIDKey, "user-2")
	_, err = interceptor(ctx2, nil, &grpc.UnaryServerInfo{FullMethod: "/test/Method"}, handler)
	if err != nil {
		t.Fatalf("user2 first request should succeed, got %v", err)
	}

	// User1 second request (should be rate limited)
	_, err = interceptor(ctx1, nil, &grpc.UnaryServerInfo{FullMethod: "/test/Method"}, handler)
	if err == nil {
		t.Fatal("user1 second request should be rate limited")
	}
}

// --- Request ID Tests ---

func TestRequestIDFromContext_Empty(t *testing.T) {
	id := RequestIDFromContext(context.Background())
	if id != "" {
		t.Fatalf("expected empty string, got %q", id)
	}
}

func TestRequestIDFromContext_WithValue(t *testing.T) {
	ctx := context.WithValue(context.Background(), RequestIDKey, "test-id-123")
	id := RequestIDFromContext(ctx)
	if id != "test-id-123" {
		t.Fatalf("expected 'test-id-123', got %q", id)
	}
}

func TestRequestIDInterceptor_GeneratesID(t *testing.T) {
	interceptor := RequestIDInterceptor()
	var capturedCtx context.Context
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		capturedCtx = ctx
		return "ok", nil
	}

	// Create context with incoming metadata (no request ID)
	md := metadata.New(map[string]string{})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test/Method"}, handler)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	id := RequestIDFromContext(capturedCtx)
	if id == "" {
		t.Fatal("expected generated request ID, got empty string")
	}
}

func TestRequestIDInterceptor_PropagatesExisting(t *testing.T) {
	interceptor := RequestIDInterceptor()
	var capturedCtx context.Context
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		capturedCtx = ctx
		return "ok", nil
	}

	// Create context with incoming metadata containing a request ID
	md := metadata.New(map[string]string{RequestIDHeader: "existing-id-456"})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test/Method"}, handler)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	id := RequestIDFromContext(capturedCtx)
	if id != "existing-id-456" {
		t.Fatalf("expected 'existing-id-456', got %q", id)
	}
}
