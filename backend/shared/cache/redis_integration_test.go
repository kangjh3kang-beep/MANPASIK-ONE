//go:build integration

package cache

import (
	"context"
	"os"
	"testing"
	"time"
)

func redisAddr() string {
	addr := os.Getenv("TEST_REDIS_ADDR")
	if addr == "" {
		return "localhost:6379"
	}
	return addr
}

func setupRedisClient(t *testing.T) *RedisClient {
	t.Helper()
	rc, err := NewRedisClient(redisAddr(), "", 15) // use DB 15 for tests
	if err != nil {
		t.Skipf("Redis not available, skipping: %v", err)
	}
	t.Cleanup(func() { rc.Close() })
	return rc
}

func TestNewRedisClient_Integration(t *testing.T) {
	rc := setupRedisClient(t)
	if rc == nil {
		t.Fatal("expected non-nil RedisClient")
	}
}

func TestSetGetDel_Integration(t *testing.T) {
	rc := setupRedisClient(t)
	ctx := context.Background()

	key := "test:setgetdel"
	value := "hello-redis"

	// Set
	if err := rc.Set(ctx, key, value, 10*time.Second); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Get
	got, err := rc.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got != value {
		t.Errorf("Get = %q, want %q", got, value)
	}

	// Del
	if err := rc.Del(ctx, key); err != nil {
		t.Fatalf("Del failed: %v", err)
	}

	// Get after Del
	_, err = rc.Get(ctx, key)
	if err == nil {
		t.Error("expected error after Del, got nil")
	}
}

func TestExists_Integration(t *testing.T) {
	rc := setupRedisClient(t)
	ctx := context.Background()

	key := "test:exists"

	// Should not exist
	exists, err := rc.Exists(ctx, key)
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if exists {
		t.Error("expected key to not exist")
	}

	// Set and check
	if err := rc.Set(ctx, key, "v", 10*time.Second); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	exists, err = rc.Exists(ctx, key)
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if !exists {
		t.Error("expected key to exist")
	}

	// Cleanup
	rc.Del(ctx, key)
}

func TestTTLExpiration_Integration(t *testing.T) {
	rc := setupRedisClient(t)
	ctx := context.Background()

	key := "test:ttl"

	// Set with 1-second TTL
	if err := rc.Set(ctx, key, "expires", 1*time.Second); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Should exist immediately
	exists, err := rc.Exists(ctx, key)
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if !exists {
		t.Error("expected key to exist immediately after Set")
	}

	// Wait for expiration
	time.Sleep(1500 * time.Millisecond)

	exists, err = rc.Exists(ctx, key)
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if exists {
		t.Error("expected key to have expired")
	}
}

func TestHealth_Integration(t *testing.T) {
	rc := setupRedisClient(t)
	ctx := context.Background()

	if err := rc.Health(ctx); err != nil {
		t.Fatalf("Health failed: %v", err)
	}
}
