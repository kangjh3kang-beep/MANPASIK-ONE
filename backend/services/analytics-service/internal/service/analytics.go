// Package service는 analytics-service의 비즈니스 로직을 구현합니다.
package service

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/google/uuid"
)

// AnalyticsEvent는 분석 이벤트 도메인 객체입니다.
type AnalyticsEvent struct {
	ID         string
	UserID     string
	EventType  string
	Properties map[string]string
	CreatedAt  time.Time
}

// DailyStats는 일별 분석 통계입니다.
type DailyStats struct {
	Date         string
	TotalEvents  int
	UniqueUsers  int
	EventsByType map[string]int
}

// UserAnalytics는 사용자별 분석 요약입니다.
type UserAnalytics struct {
	UserID        string
	TotalEvents   int
	LastEventAt   time.Time
	TopEventTypes []string
}

// AnalyticsRepository는 분석 데이터 저장소 인터페이스입니다.
type AnalyticsRepository interface {
	TrackEvent(ctx context.Context, event *AnalyticsEvent) error
	GetUserAnalytics(ctx context.Context, userID string) (*UserAnalytics, error)
	GetDailyStats(ctx context.Context, date string) (*DailyStats, error)
	ListRecentEvents(ctx context.Context, limit int) ([]*AnalyticsEvent, error)
}

// AnalyticsService는 분석 서비스 핵심 로직입니다.
type AnalyticsService struct {
	repo AnalyticsRepository
}

// NewAnalyticsService는 AnalyticsService를 생성합니다.
func NewAnalyticsService(repo AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{repo: repo}
}

// TrackEvent는 새 분석 이벤트를 기록합니다.
func (s *AnalyticsService) TrackEvent(ctx context.Context, userID, eventType string, properties map[string]string) (string, error) {
	if userID == "" {
		return "", errors.New("user_id는 필수입니다")
	}
	if eventType == "" {
		return "", errors.New("event_type은 필수입니다")
	}

	event := &AnalyticsEvent{
		ID:         uuid.New().String(),
		UserID:     userID,
		EventType:  eventType,
		Properties: properties,
		CreatedAt:  time.Now(),
	}

	if err := s.repo.TrackEvent(ctx, event); err != nil {
		return "", err
	}

	return event.ID, nil
}

// GetUserAnalytics는 사용자 분석 요약을 반환합니다.
func (s *AnalyticsService) GetUserAnalytics(ctx context.Context, userID string) (*UserAnalytics, error) {
	if userID == "" {
		return nil, errors.New("user_id는 필수입니다")
	}
	return s.repo.GetUserAnalytics(ctx, userID)
}

// GetDailyStats는 일별 분석 통계를 반환합니다.
func (s *AnalyticsService) GetDailyStats(ctx context.Context, date string) (*DailyStats, error) {
	if date == "" {
		return nil, errors.New("date는 필수입니다")
	}
	return s.repo.GetDailyStats(ctx, date)
}

// ListRecentEvents는 최근 이벤트 목록을 반환합니다.
func (s *AnalyticsService) ListRecentEvents(ctx context.Context, limit int) ([]*AnalyticsEvent, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.ListRecentEvents(ctx, limit)
}

// TopEventTypes는 이벤트 타입별 횟수를 내림차순으로 정렬하여 상위 N개를 반환합니다.
func TopEventTypes(eventsByType map[string]int, n int) []string {
	type kv struct {
		Key   string
		Value int
	}
	var sorted []kv
	for k, v := range eventsByType {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})
	var result []string
	for i := 0; i < len(sorted) && i < n; i++ {
		result = append(result, sorted[i].Key)
	}
	return result
}
