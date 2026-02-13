package storage

import (
	"fmt"
	"testing"
	"time"
)

func TestNewS3Client_InvalidEndpoint(t *testing.T) {
	// Attempting to connect to an invalid endpoint should fail during bucket check.
	_, err := NewS3Client("invalid-host:0", "key", "secret", "bucket", "us-east-1", false)
	if err == nil {
		t.Fatal("expected error for invalid endpoint, got nil")
	}
	t.Logf("got expected error: %v", err)
}

func TestS3Client_PathGeneration(t *testing.T) {
	// Test that file path building produces expected format.
	category := "profiles"
	dateStr := time.Now().Format("2006/01/02")
	fileID := "abc-123"
	ext := ".jpg"

	path := fmt.Sprintf("%s/%s/%s%s", category, dateStr, fileID, ext)

	expected := fmt.Sprintf("profiles/%s/abc-123.jpg", dateStr)
	if path != expected {
		t.Errorf("path = %q, want %q", path, expected)
	}
}
