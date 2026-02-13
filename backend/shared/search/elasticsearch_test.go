package search

import (
	"encoding/json"
	"testing"
)

func TestNewESClient_InvalidURL(t *testing.T) {
	_, err := NewESClient("http://invalid-host:9999", "", "")
	if err == nil {
		t.Error("expected error for unreachable ES")
	}
}

func TestSearchResponse_Parse(t *testing.T) {
	// Test that SearchResponse struct can be unmarshaled
	data := `{"hits":{"total":{"value":1},"hits":[{"_id":"1","_score":1.0,"_source":{"name":"test"}}]}}`
	// Test JSON unmarshaling
	var resp SearchResponse
	if err := json.Unmarshal([]byte(data), &resp); err != nil {
		t.Errorf("unmarshal failed: %v", err)
	}
	if resp.Hits.Total.Value != 1 {
		t.Errorf("expected 1 hit, got %d", resp.Hits.Total.Value)
	}
	if resp.Hits.Hits[0].ID != "1" {
		t.Errorf("expected id '1', got '%s'", resp.Hits.Hits[0].ID)
	}
}
