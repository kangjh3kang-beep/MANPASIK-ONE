package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type auditRecordEvent struct {
	AdminID      string            `json:"admin_id"`
	Action       string            `json:"action"`
	ResourceType string            `json:"resource_type"`
	ResourceID   string            `json:"resource_id"`
	Description  string            `json:"description"`
	IPAddress    string            `json:"ip_address"`
	UserAgent    string            `json:"user_agent"`
	Metadata     map[string]string `json:"metadata"`
	OccurredAt   string            `json:"occurred_at"`
}

type auditRecordResult struct {
	EntryID string
	Status  string
}

type auditEventRecorder interface {
	RecordAuditEvent(ctx context.Context, event auditRecordEvent) (*auditRecordResult, error)
}

type HTTPAuditEventRecorder struct {
	endpoint string
	client   *http.Client
}

func NewHTTPAuditEventRecorder(endpoint string, client *http.Client) *HTTPAuditEventRecorder {
	if client == nil {
		client = &http.Client{Timeout: 2 * time.Second}
	}
	return &HTTPAuditEventRecorder{endpoint: endpoint, client: client}
}

func (r *HTTPAuditEventRecorder) RecordAuditEvent(ctx context.Context, event auditRecordEvent) (*auditRecordResult, error) {
	if r == nil || r.endpoint == "" {
		return nil, fmt.Errorf("audit intake endpoint is not configured")
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("audit intake returned %d", resp.StatusCode)
	}

	var body struct {
		EntryID string `json:"entry_id"`
		Status  string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	if body.Status == "" {
		body.Status = "persisted"
	}
	return &auditRecordResult{EntryID: body.EntryID, Status: body.Status}, nil
}
