// Package memory는 analytics-service의 인메모리 저장소입니다.
package memory

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/manpasik/backend/services/analytics-service/internal/service"
)

// AnalyticsRepository는 분석 이벤트 인메모리 저장소입니다.
type AnalyticsRepository struct {
	mu    sync.RWMutex
	store map[string]*service.AnalyticsEvent // ID -> Event
}

// NewAnalyticsRepository는 새 인메모리 분석 저장소를 생성합니다.
func NewAnalyticsRepository() *AnalyticsRepository {
	return &AnalyticsRepository{
		store: make(map[string]*service.AnalyticsEvent),
	}
}

// TrackEvent는 분석 이벤트를 저장합니다.
func (r *AnalyticsRepository) TrackEvent(_ context.Context, event *service.AnalyticsEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[event.ID] = event
	return nil
}

// GetUserAnalytics는 사용자별 분석 요약을 반환합니다.
func (r *AnalyticsRepository) GetUserAnalytics(_ context.Context, userID string) (*service.UserAnalytics, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var lastEventAt time.Time
	eventsByType := make(map[string]int)
	totalEvents := 0

	for _, event := range r.store {
		if event.UserID != userID {
			continue
		}
		totalEvents++
		eventsByType[event.EventType]++
		if event.CreatedAt.After(lastEventAt) {
			lastEventAt = event.CreatedAt
		}
	}

	topTypes := service.TopEventTypes(eventsByType, 5)

	return &service.UserAnalytics{
		UserID:        userID,
		TotalEvents:   totalEvents,
		LastEventAt:   lastEventAt,
		TopEventTypes: topTypes,
	}, nil
}

// GetDailyStats는 특정 날짜의 분석 통계를 반환합니다.
func (r *AnalyticsRepository) GetDailyStats(_ context.Context, date string) (*service.DailyStats, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	uniqueUsers := make(map[string]struct{})
	eventsByType := make(map[string]int)
	totalEvents := 0

	for _, event := range r.store {
		eventDate := event.CreatedAt.Format("2006-01-02")
		if eventDate != date {
			continue
		}
		totalEvents++
		uniqueUsers[event.UserID] = struct{}{}
		eventsByType[event.EventType]++
	}

	return &service.DailyStats{
		Date:         date,
		TotalEvents:  totalEvents,
		UniqueUsers:  len(uniqueUsers),
		EventsByType: eventsByType,
	}, nil
}

// ListRecentEvents는 최근 이벤트를 시간 역순으로 반환합니다.
func (r *AnalyticsRepository) ListRecentEvents(_ context.Context, limit int) ([]*service.AnalyticsEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	events := make([]*service.AnalyticsEvent, 0, len(r.store))
	for _, event := range r.store {
		events = append(events, event)
	}

	// 시간 역순 정렬 (최신 먼저)
	sort.Slice(events, func(i, j int) bool {
		return events[i].CreatedAt.After(events[j].CreatedAt)
	})

	if limit > len(events) {
		limit = len(events)
	}

	return events[:limit], nil
}
