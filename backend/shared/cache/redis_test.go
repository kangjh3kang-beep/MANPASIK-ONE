package cache

import "testing"

func TestRedisClientConfig(t *testing.T) {
	// Test that NewRedisClient returns error on invalid address
	_, err := NewRedisClient("invalid:0", "", 0)
	if err == nil {
		t.Error("expected error for invalid Redis address")
	}
}
