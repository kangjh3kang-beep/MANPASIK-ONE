package observability

import (
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Metrics collects basic application metrics
type Metrics struct {
	mu              sync.RWMutex
	requestCount    map[string]int64
	requestDuration map[string][]time.Duration
	errorCount      map[string]int64
	startTime       time.Time
}

// NewMetrics creates a new metrics collector
func NewMetrics() *Metrics {
	return &Metrics{
		requestCount:    make(map[string]int64),
		requestDuration: make(map[string][]time.Duration),
		errorCount:      make(map[string]int64),
		startTime:       time.Now(),
	}
}

// RecordRequest records a request
func (m *Metrics) RecordRequest(method string, duration time.Duration, statusCode int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := method
	m.requestCount[key]++
	m.requestDuration[key] = append(m.requestDuration[key], duration)
	if statusCode >= 400 {
		m.errorCount[key]++
	}
}

// GetStats returns current metrics
func (m *Metrics) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	totalReqs := int64(0)
	totalErrs := int64(0)
	for _, v := range m.requestCount {
		totalReqs += v
	}
	for _, v := range m.errorCount {
		totalErrs += v
	}
	return map[string]interface{}{
		"uptime_seconds": time.Since(m.startTime).Seconds(),
		"total_requests": totalReqs,
		"total_errors":   totalErrs,
		"methods":        m.requestCount,
	}
}

// PrometheusHandler returns metrics in Prometheus format
func (m *Metrics) PrometheusHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.mu.RLock()
		defer m.mu.RUnlock()

		w.Header().Set("Content-Type", "text/plain; version=0.0.4")

		// Uptime
		uptime := time.Since(m.startTime).Seconds()
		w.Write([]byte("# HELP manpasik_uptime_seconds Service uptime\n"))
		w.Write([]byte("# TYPE manpasik_uptime_seconds gauge\n"))
		w.Write([]byte("manpasik_uptime_seconds " + strconv.FormatFloat(uptime, 'f', 2, 64) + "\n\n"))

		// Request count
		w.Write([]byte("# HELP manpasik_requests_total Total requests\n"))
		w.Write([]byte("# TYPE manpasik_requests_total counter\n"))
		for method, count := range m.requestCount {
			w.Write([]byte("manpasik_requests_total{method=\"" + method + "\"} " + strconv.FormatInt(count, 10) + "\n"))
		}
		w.Write([]byte("\n"))

		// Error count
		w.Write([]byte("# HELP manpasik_errors_total Total errors\n"))
		w.Write([]byte("# TYPE manpasik_errors_total counter\n"))
		for method, count := range m.errorCount {
			w.Write([]byte("manpasik_errors_total{method=\"" + method + "\"} " + strconv.FormatInt(count, 10) + "\n"))
		}
		w.Write([]byte("\n"))

		// Average duration per method
		w.Write([]byte("# HELP manpasik_request_duration_seconds Average request duration\n"))
		w.Write([]byte("# TYPE manpasik_request_duration_seconds gauge\n"))
		for method, durations := range m.requestDuration {
			if len(durations) == 0 {
				continue
			}
			total := time.Duration(0)
			for _, d := range durations {
				total += d
			}
			avg := total / time.Duration(len(durations))
			w.Write([]byte("manpasik_request_duration_seconds{method=\"" + method + "\"} " + strconv.FormatFloat(avg.Seconds(), 'f', 6, 64) + "\n"))
		}
	}
}
