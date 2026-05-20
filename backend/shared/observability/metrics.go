package observability

import (
	"net/http"
	"sort"
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

// MetricsSnapshot 은 Metrics 의 현재 상태를 외부 모듈(예: canary 어댑터) 이
// 사용할 수 있는 형태로 노출.
//
// 누적값 기준이며, 슬라이딩 윈도우가 필요하면 외부에서 주기적으로
// 새 인스턴스를 생성하여 사용 (현재 Metrics 는 reset 미지원).
type MetricsSnapshot struct {
	TotalRequests int64
	TotalErrors   int64
	ErrorRate     float64       // 0~1
	LatencyP50    time.Duration
	LatencyP99    time.Duration
	UptimeSeconds float64
}

// Snapshot 은 모든 메서드의 통계를 한 번에 집계하여 반환.
//
// LatencyP99 는 정렬 기반 백분위 (선형 보간 없음). 데이터 0 건이면 0.
func (m *Metrics) Snapshot() MetricsSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var totalReqs, totalErrs int64
	var allDurations []time.Duration
	for _, count := range m.requestCount {
		totalReqs += count
	}
	for _, count := range m.errorCount {
		totalErrs += count
	}
	for _, durs := range m.requestDuration {
		allDurations = append(allDurations, durs...)
	}

	snap := MetricsSnapshot{
		TotalRequests: totalReqs,
		TotalErrors:   totalErrs,
		UptimeSeconds: time.Since(m.startTime).Seconds(),
	}
	if totalReqs > 0 {
		snap.ErrorRate = float64(totalErrs) / float64(totalReqs)
	}
	if len(allDurations) > 0 {
		sort.Slice(allDurations, func(i, j int) bool {
			return allDurations[i] < allDurations[j]
		})
		snap.LatencyP50 = percentileDuration(allDurations, 50)
		snap.LatencyP99 = percentileDuration(allDurations, 99)
	}
	return snap
}

// percentileDuration 은 정렬된 durations 에서 p 번째 백분위 반환 (p: 0~100).
func percentileDuration(sortedDurations []time.Duration, p int) time.Duration {
	n := len(sortedDurations)
	if n == 0 {
		return 0
	}
	if p <= 0 {
		return sortedDurations[0]
	}
	if p >= 100 {
		return sortedDurations[n-1]
	}
	idx := (p * (n - 1)) / 100
	return sortedDurations[idx]
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
