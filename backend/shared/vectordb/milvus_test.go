package vectordb

import (
	"testing"
)

func TestSearchResult_Fields(t *testing.T) {
	r := SearchResult{
		ID:        "test-id-123",
		SessionID: "session-456",
		Score:     0.95,
	}

	if r.ID != "test-id-123" {
		t.Errorf("expected ID 'test-id-123', got '%s'", r.ID)
	}
	if r.SessionID != "session-456" {
		t.Errorf("expected SessionID 'session-456', got '%s'", r.SessionID)
	}
	if r.Score != 0.95 {
		t.Errorf("expected Score 0.95, got %f", r.Score)
	}
}

func TestNewMilvusClient_InvalidAddress(t *testing.T) {
	// Connecting to an invalid address should fail
	_, err := NewMilvusClient("invalid-host:99999", "test_collection", 128)
	if err == nil {
		t.Error("expected error when connecting to invalid address, got nil")
	}
}
