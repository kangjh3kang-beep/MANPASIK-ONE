package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ESClient is a lightweight Elasticsearch client
type ESClient struct {
	baseURL  string
	username string
	password string
	http     *http.Client
}

// NewESClient creates a new Elasticsearch client
func NewESClient(url, username, password string) (*ESClient, error) {
	client := &ESClient{
		baseURL:  url,
		username: username,
		password: password,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	// Health check
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Health(ctx); err != nil {
		return nil, fmt.Errorf("elasticsearch health check failed: %w", err)
	}

	return client, nil
}

// Health checks ES cluster health
func (c *ESClient) Health(ctx context.Context) error {
	req, _ := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/_cluster/health", nil)
	c.setAuth(req)
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("ES health status: %d", resp.StatusCode)
	}
	return nil
}

func (c *ESClient) setAuth(req *http.Request) {
	if c.username != "" {
		req.SetBasicAuth(c.username, c.password)
	}
}

// IndexDocument indexes a document
func (c *ESClient) IndexDocument(ctx context.Context, index, id string, doc interface{}) error {
	data, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/%s/_doc/%s", c.baseURL, index, id)
	req, _ := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	c.setAuth(req)
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("index failed (%d): %s", resp.StatusCode, string(body))
	}
	return nil
}

// Search performs a search query
func (c *ESClient) Search(ctx context.Context, index string, query map[string]interface{}) (*SearchResponse, error) {
	data, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/%s/_search", c.baseURL, index)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	c.setAuth(req)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search failed (%d): %s", resp.StatusCode, string(body))
	}
	var result SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteDocument deletes a document
func (c *ESClient) DeleteDocument(ctx context.Context, index, id string) error {
	url := fmt.Sprintf("%s/%s/_doc/%s", c.baseURL, index, id)
	req, _ := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	c.setAuth(req)
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// CreateIndex creates an ES index with mappings
func (c *ESClient) CreateIndex(ctx context.Context, index string, mappings map[string]interface{}) error {
	data, err := json.Marshal(mappings)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/%s", c.baseURL, index)
	req, _ := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	c.setAuth(req)
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 && resp.StatusCode != 400 { // 400 may mean already exists
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("create index failed (%d): %s", resp.StatusCode, string(body))
	}
	return nil
}

// Close closes the client (no-op for HTTP)
func (c *ESClient) Close() error {
	return nil
}

// SearchResponse represents ES search results
type SearchResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []Hit `json:"hits"`
	} `json:"hits"`
}

// Hit represents a single search hit
type Hit struct {
	ID     string          `json:"_id"`
	Score  float64         `json:"_score"`
	Source json.RawMessage `json:"_source"`
}
