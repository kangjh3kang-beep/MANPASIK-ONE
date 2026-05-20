package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPAuditEventRecorderPostsEvent(t *testing.T) {
	var received auditRecordEvent
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"status":"persisted","entry_id":"audit-1"}`))
	}))
	defer server.Close()

	recorder := NewHTTPAuditEventRecorder(server.URL, server.Client())
	result, err := recorder.RecordAuditEvent(context.Background(), auditRecordEvent{
		AdminID:      "system:gateway",
		Action:       "measure.trace.serverProcessed",
		ResourceType: "measurement_trace",
		ResourceID:   "session-1",
		Metadata:     map[string]string{"phase": "serverProcessed"},
	})

	if err != nil {
		t.Fatalf("RecordAuditEvent failed: %v", err)
	}
	if result.EntryID != "audit-1" || result.Status != "persisted" {
		t.Fatalf("result mismatch: %#v", result)
	}
	if received.Action != "measure.trace.serverProcessed" {
		t.Fatalf("received action = %q", received.Action)
	}
	if received.Metadata["phase"] != "serverProcessed" {
		t.Fatalf("received metadata phase = %q", received.Metadata["phase"])
	}
}
