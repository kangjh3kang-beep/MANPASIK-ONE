// Package service는 analytics-service의 비즈니스 로직을 구현합니다.
package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// 이벤트 타입 화이트리스트
// ============================================================================

// validEventTypes는 허용된 이벤트 타입 목록입니다. "custom." 접두사는 항상 허용됩니다.
var validEventTypes = map[string]bool{
	"screen_view":         true,
	"button_click":        true,
	"measurement_start":   true,
	"measurement_complete": true,
	"login":               true,
	"logout":              true,
	"purchase":            true,
	"subscription_change": true,
	"device_connect":      true,
	"device_disconnect":   true,
	"coaching_view":       true,
	"reservation_create":  true,
	"food_analysis":       true,
	"emergency_report":    true,
	"page_view":           true,
	"app_open":            true,
	"app_close":           true,
}

// maxProperties는 이벤트 속성의 최대 개수입니다.
const maxProperties = 50

// maxPropertyValueLen는 속성 값의 최대 길이입니다.
const maxPropertyValueLen = 500

// ============================================================================
// 도메인 모델
// ============================================================================

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

// RetentionData는 사용자 유지율 정보입니다.
type RetentionData struct {
	UserID        string
	Days          int
	ActiveDays    int
	RetentionRate float64
}

// FunnelStep은 이벤트 퍼널의 한 단계입니다.
type FunnelStep struct {
	EventType string
	Count     int
	Completed bool
}

// ============================================================================
// Repository 인터페이스
// ============================================================================

// AnalyticsRepository는 분석 데이터 저장소 인터페이스입니다.
type AnalyticsRepository interface {
	TrackEvent(ctx context.Context, event *AnalyticsEvent) error
	GetUserAnalytics(ctx context.Context, userID string) (*UserAnalytics, error)
	GetDailyStats(ctx context.Context, date string) (*DailyStats, error)
	ListRecentEvents(ctx context.Context, limit int) ([]*AnalyticsEvent, error)
	ListEventsByUser(ctx context.Context, userID string) ([]*AnalyticsEvent, error)
}

// ============================================================================
// 에러 정의
// ============================================================================

var (
	ErrInvalidInput    = errors.New("invalid input")
	ErrInvalidEventType = errors.New("invalid event type")
)

// ============================================================================
// AnalyticsService
// ============================================================================

// AnalyticsService는 분석 서비스 핵심 로직입니다.
type AnalyticsService struct {
	repo AnalyticsRepository
}

// NewAnalyticsService는 AnalyticsService를 생성합니다.
func NewAnalyticsService(repo AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{repo: repo}
}

// isValidEventType은 이벤트 타입이 유효한지 확인합니다.
func isValidEventType(eventType string) bool {
	if validEventTypes[eventType] {
		return true
	}
	if strings.HasPrefix(eventType, "custom.") {
		return true
	}
	return false
}

// validateProperties는 속성의 유효성을 검증합니다.
func validateProperties(properties map[string]string) error {
	if len(properties) > maxProperties {
		return fmt.Errorf("%w: too many properties (max %d)", ErrInvalidInput, maxProperties)
	}
	for key, value := range properties {
		if strings.TrimSpace(key) == "" {
			return fmt.Errorf("%w: property key must not be empty", ErrInvalidInput)
		}
		if len(value) > maxPropertyValueLen {
			return fmt.Errorf("%w: property value too long for key %q (max %d)", ErrInvalidInput, key, maxPropertyValueLen)
		}
	}
	return nil
}

// TrackEvent는 새 분석 이벤트를 기록합니다.
func (s *AnalyticsService) TrackEvent(ctx context.Context, userID, eventType string, properties map[string]string) (string, error) {
	if userID == "" {
		return "", errors.New("user_id는 필수입니다")
	}
	if eventType == "" {
		return "", errors.New("event_type은 필수입니다")
	}

	if !isValidEventType(eventType) {
		return "", fmt.Errorf("%w: %q is not allowed (use 'custom.' prefix for custom events)", ErrInvalidEventType, eventType)
	}

	if properties != nil {
		if err := validateProperties(properties); err != nil {
			return "", err
		}
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

// GetUserRetention은 사용자의 최근 N일간 유지율을 계산합니다.
func (s *AnalyticsService) GetUserRetention(ctx context.Context, userID string, days int) (*RetentionData, error) {
	if userID == "" {
		return nil, errors.New("user_id는 필수입니다")
	}
	if days <= 0 {
		days = 7
	}

	events, err := s.repo.ListEventsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	cutoff := now.AddDate(0, 0, -days)
	activeDaySet := make(map[string]bool)
	for _, ev := range events {
		if ev.CreatedAt.After(cutoff) {
			day := ev.CreatedAt.Format("2006-01-02")
			activeDaySet[day] = true
		}
	}

	activeDays := len(activeDaySet)
	rate := 0.0
	if days > 0 {
		rate = float64(activeDays) / float64(days)
		if rate > 1.0 {
			rate = 1.0
		}
	}

	return &RetentionData{
		UserID:        userID,
		Days:          days,
		ActiveDays:    activeDays,
		RetentionRate: rate,
	}, nil
}

// GetEventFunnel은 이벤트 시퀀스의 달성률을 분석합니다.
func (s *AnalyticsService) GetEventFunnel(ctx context.Context, userID string, eventTypes []string) ([]FunnelStep, error) {
	if userID == "" {
		return nil, errors.New("user_id는 필수입니다")
	}
	if len(eventTypes) == 0 {
		return nil, errors.New("event_types는 필수입니다")
	}

	events, err := s.repo.ListEventsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 이벤트 타입별 발생 횟수
	counts := make(map[string]int)
	for _, ev := range events {
		counts[ev.EventType]++
	}

	steps := make([]FunnelStep, len(eventTypes))
	for i, et := range eventTypes {
		count := counts[et]
		steps[i] = FunnelStep{
			EventType: et,
			Count:     count,
			Completed: count > 0,
		}
	}

	return steps, nil
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
