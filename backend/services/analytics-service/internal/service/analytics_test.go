package service_test

import (
	"context"
	"testing"

	"github.com/manpasik/backend/services/analytics-service/internal/repository/memory"
	"github.com/manpasik/backend/services/analytics-service/internal/service"
)

func setupTestService() *service.AnalyticsService {
	repo := memory.NewAnalyticsRepository()
	return service.NewAnalyticsService(repo)
}

func TestTrackEvent_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	props := map[string]string{"page": "/home", "action": "click"}
	id, err := svc.TrackEvent(ctx, "user-1", "page_view", props)
	if err != nil {
		t.Fatalf("이벤트 기록 실패: %v", err)
	}
	if id == "" {
		t.Fatal("이벤트 ID가 비어 있음")
	}
}

func TestTrackEvent_MissingUserID(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, err := svc.TrackEvent(ctx, "", "page_view", nil)
	if err == nil {
		t.Fatal("빈 user_id에 에러가 반환되어야 함")
	}
}

func TestTrackEvent_MissingEventType(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, err := svc.TrackEvent(ctx, "user-1", "", nil)
	if err == nil {
		t.Fatal("빈 event_type에 에러가 반환되어야 함")
	}
}

func TestGetUserAnalytics_Empty(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	analytics, err := svc.GetUserAnalytics(ctx, "user-no-events")
	if err != nil {
		t.Fatalf("사용자 분석 조회 실패: %v", err)
	}
	if analytics == nil {
		t.Fatal("UserAnalytics가 nil이면 안 됨")
	}
	if analytics.TotalEvents != 0 {
		t.Fatalf("이벤트 없는 사용자의 TotalEvents가 0이어야 함: got %d", analytics.TotalEvents)
	}
	if analytics.UserID != "user-no-events" {
		t.Fatalf("UserID 불일치: got %s, want user-no-events", analytics.UserID)
	}
	if len(analytics.TopEventTypes) != 0 {
		t.Fatalf("이벤트 없는 사용자의 TopEventTypes가 비어야 함: got %v", analytics.TopEventTypes)
	}
}

func TestGetDailyStats_Empty(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	stats, err := svc.GetDailyStats(ctx, "2026-01-01")
	if err != nil {
		t.Fatalf("일별 통계 조회 실패: %v", err)
	}
	if stats == nil {
		t.Fatal("DailyStats가 nil이면 안 됨")
	}
	if stats.TotalEvents != 0 {
		t.Fatalf("이벤트 없는 날의 TotalEvents가 0이어야 함: got %d", stats.TotalEvents)
	}
	if stats.UniqueUsers != 0 {
		t.Fatalf("이벤트 없는 날의 UniqueUsers가 0이어야 함: got %d", stats.UniqueUsers)
	}
	if stats.Date != "2026-01-01" {
		t.Fatalf("Date 불일치: got %s, want 2026-01-01", stats.Date)
	}
}

func TestTrackEvent_MultipleEvents(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	ids := make(map[string]bool)
	for i := 0; i < 10; i++ {
		id, err := svc.TrackEvent(ctx, "user-multi", "page_view", nil)
		if err != nil {
			t.Fatalf("이벤트 %d 기록 실패: %v", i, err)
		}
		if ids[id] {
			t.Fatalf("중복 이벤트 ID: %s", id)
		}
		ids[id] = true
	}

	analytics, err := svc.GetUserAnalytics(ctx, "user-multi")
	if err != nil {
		t.Fatalf("사용자 분석 조회 실패: %v", err)
	}
	if analytics.TotalEvents != 10 {
		t.Fatalf("TotalEvents 불일치: got %d, want 10", analytics.TotalEvents)
	}
}

func TestGetUserAnalytics_WithEvents(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	svc.TrackEvent(ctx, "user-stats", "page_view", nil)
	svc.TrackEvent(ctx, "user-stats", "page_view", nil)
	svc.TrackEvent(ctx, "user-stats", "button_click", nil)
	svc.TrackEvent(ctx, "user-stats", "purchase", nil)

	analytics, err := svc.GetUserAnalytics(ctx, "user-stats")
	if err != nil {
		t.Fatalf("사용자 분석 조회 실패: %v", err)
	}
	if analytics.TotalEvents != 4 {
		t.Fatalf("TotalEvents 불일치: got %d, want 4", analytics.TotalEvents)
	}
	if len(analytics.TopEventTypes) == 0 {
		t.Fatal("TopEventTypes가 비어 있으면 안 됨")
	}
	// page_view가 2회로 가장 많으므로 첫 번째여야 함
	if analytics.TopEventTypes[0] != "page_view" {
		t.Fatalf("TopEventTypes[0] 불일치: got %s, want page_view", analytics.TopEventTypes[0])
	}
}

func TestListRecentEvents(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		_, err := svc.TrackEvent(ctx, "user-recent", "event_type", nil)
		if err != nil {
			t.Fatalf("이벤트 기록 실패: %v", err)
		}
	}

	events, err := svc.ListRecentEvents(ctx, 3)
	if err != nil {
		t.Fatalf("최근 이벤트 조회 실패: %v", err)
	}
	if len(events) != 3 {
		t.Fatalf("반환 이벤트 수 불일치: got %d, want 3", len(events))
	}
}

func TestListRecentEvents_DefaultLimit(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		svc.TrackEvent(ctx, "user-default", "test_event", nil)
	}

	// limit <= 0 이면 기본값 50 사용
	events, err := svc.ListRecentEvents(ctx, 0)
	if err != nil {
		t.Fatalf("최근 이벤트 조회 실패: %v", err)
	}
	if len(events) != 3 {
		t.Fatalf("반환 이벤트 수 불일치: got %d, want 3", len(events))
	}
}

func TestGetUserAnalytics_MissingUserID(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, err := svc.GetUserAnalytics(ctx, "")
	if err == nil {
		t.Fatal("빈 user_id에 에러가 반환되어야 함")
	}
}

func TestGetDailyStats_MissingDate(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, err := svc.GetDailyStats(ctx, "")
	if err == nil {
		t.Fatal("빈 date에 에러가 반환되어야 함")
	}
}

func TestTrackEvent_WithProperties(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	props := map[string]string{
		"page":     "/products/123",
		"referrer": "/home",
		"device":   "mobile",
	}
	id, err := svc.TrackEvent(ctx, "user-props", "page_view", props)
	if err != nil {
		t.Fatalf("이벤트 기록 실패: %v", err)
	}
	if id == "" {
		t.Fatal("이벤트 ID가 비어 있음")
	}

	events, err := svc.ListRecentEvents(ctx, 10)
	if err != nil {
		t.Fatalf("최근 이벤트 조회 실패: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("이벤트 수 불일치: got %d, want 1", len(events))
	}
	if events[0].Properties["page"] != "/products/123" {
		t.Fatalf("Properties[page] 불일치: got %s", events[0].Properties["page"])
	}
	if events[0].Properties["device"] != "mobile" {
		t.Fatalf("Properties[device] 불일치: got %s", events[0].Properties["device"])
	}
}
