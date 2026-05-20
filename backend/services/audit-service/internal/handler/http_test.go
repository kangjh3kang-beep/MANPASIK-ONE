package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/manpasik/backend/services/audit-service/internal/service"
	"go.uber.org/zap"
)

type httpTestAuditRepo struct {
	entries []*service.AuditEntry
}

func (r *httpTestAuditRepo) Store(_ context.Context, entry *service.AuditEntry) error {
	r.entries = append(r.entries, entry)
	return nil
}

func (r *httpTestAuditRepo) List(_ context.Context, _ service.AuditFilter) ([]*service.AuditEntry, int, error) {
	return r.entries, len(r.entries), nil
}

func (r *httpTestAuditRepo) Get(_ context.Context, entryID string) (*service.AuditEntry, error) {
	for _, entry := range r.entries {
		if entry.ID == entryID {
			return entry, nil
		}
	}
	return nil, nil
}

func TestHTTPRecordAuditEventPersistsEntry(t *testing.T) {
	repo := &httpTestAuditRepo{}
	svc := service.NewAuditService(zap.NewNop(), repo)
	mux := http.NewServeMux()
	NewHTTPHandler(svc, zap.NewNop()).RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodPost, "/audit/events", strings.NewReader(`{
		"admin_id":"system:gateway",
		"action":"measure.trace.serverProcessed",
		"resource_type":"measurement_trace",
		"resource_id":"session-1",
		"description":"Measure trace phase serverProcessed",
		"ip_address":"10.0.0.1",
		"user_agent":"gateway-test",
		"occurred_at":"2026-05-01T12:00:00Z",
		"metadata":{"phase":"serverProcessed","schema_version":"measure_trace.v1"}
	}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusCreated, rec.Body.String())
	}
	if len(repo.entries) != 1 {
		t.Fatalf("stored entries = %d, want 1", len(repo.entries))
	}
	entry := repo.entries[0]
	if entry.Action != "measure.trace.serverProcessed" {
		t.Fatalf("action = %q", entry.Action)
	}
	if entry.ResourceType != "measurement_trace" || entry.ResourceID != "session-1" {
		t.Fatalf("resource mismatch: type=%q id=%q", entry.ResourceType, entry.ResourceID)
	}
	if entry.Metadata["phase"] != "serverProcessed" {
		t.Fatalf("metadata phase = %q", entry.Metadata["phase"])
	}
	if !entry.Timestamp.Equal(time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)) {
		t.Fatalf("timestamp = %s", entry.Timestamp)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("response json: %v", err)
	}
	if body["status"] != "persisted" || body["entry_id"] == "" {
		t.Fatalf("response mismatch: %v", body)
	}
}

func TestHTTPRecordAuditEventRequiresActionAndResourceType(t *testing.T) {
	repo := &httpTestAuditRepo{}
	svc := service.NewAuditService(zap.NewNop(), repo)
	mux := http.NewServeMux()
	NewHTTPHandler(svc, zap.NewNop()).RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodPost, "/audit/events", strings.NewReader(`{"action":""}`))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if len(repo.entries) != 0 {
		t.Fatalf("stored entries = %d, want 0", len(repo.entries))
	}
}
