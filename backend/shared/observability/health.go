package observability

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"
)

// HealthCheck provides standardized health check info
type HealthCheck struct {
	ServiceName string
	Version     string
	StartTime   time.Time
}

// NewHealthCheck creates a new health check handler
func NewHealthCheck(serviceName, version string) *HealthCheck {
	return &HealthCheck{ServiceName: serviceName, Version: version, StartTime: time.Now()}
}

// Handler returns an HTTP handler that serves health check info
func (h *HealthCheck) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":     "healthy",
			"service":    h.ServiceName,
			"version":    h.Version,
			"uptime":     time.Since(h.StartTime).String(),
			"goroutines": runtime.NumGoroutine(),
			"memory_mb":  memStats.Alloc / 1024 / 1024,
		})
	}
}
