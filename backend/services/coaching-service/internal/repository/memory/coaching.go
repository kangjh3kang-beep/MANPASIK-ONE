// Package memory는 인메모리 코칭 저장소입니다 (개발/테스트용).
package memory

import (
	"context"
	"sync"
	"time"

	"github.com/manpasik/backend/services/coaching-service/internal/service"
)

// ============================================================================
// HealthGoalRepository — 인메모리 건강 목표 저장소
// ============================================================================

// HealthGoalRepository는 인메모리 건강 목표 저장소입니다.
type HealthGoalRepository struct {
	mu    sync.RWMutex
	byID  map[string]*service.HealthGoal
	byUID map[string][]*service.HealthGoal // key: userID
}

// NewHealthGoalRepository는 인메모리 HealthGoalRepository를 생성합니다.
func NewHealthGoalRepository() *HealthGoalRepository {
	return &HealthGoalRepository{
		byID:  make(map[string]*service.HealthGoal),
		byUID: make(map[string][]*service.HealthGoal),
	}
}

// Create는 건강 목표를 생성합니다.
func (r *HealthGoalRepository) Create(_ context.Context, goal *service.HealthGoal) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *goal
	r.byID[goal.GoalID] = &cp
	r.byUID[goal.UserID] = append(r.byUID[goal.UserID], &cp)
	return nil
}

// GetByUserID는 사용자의 건강 목표를 조회합니다. statusFilter가 0이면 전체를 반환합니다.
func (r *HealthGoalRepository) GetByUserID(_ context.Context, userID string, statusFilter service.GoalStatus) ([]*service.HealthGoal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	goals, ok := r.byUID[userID]
	if !ok {
		return nil, nil
	}

	result := make([]*service.HealthGoal, 0, len(goals))
	for _, g := range goals {
		if statusFilter != service.GoalStatusUnknown && g.Status != statusFilter {
			continue
		}
		cp := *g
		result = append(result, &cp)
	}
	return result, nil
}

// GetByID는 목표 ID로 건강 목표를 조회합니다.
func (r *HealthGoalRepository) GetByID(_ context.Context, id string) (*service.HealthGoal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, ok := r.byID[id]
	if !ok {
		return nil, nil
	}
	cp := *g
	return &cp, nil
}

// Update는 건강 목표를 업데이트합니다.
func (r *HealthGoalRepository) Update(_ context.Context, goal *service.HealthGoal) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *goal
	r.byID[goal.GoalID] = &cp

	// byUID 슬라이스 내의 해당 목표도 업데이트
	goals := r.byUID[goal.UserID]
	for i, g := range goals {
		if g.GoalID == goal.GoalID {
			goals[i] = &cp
			break
		}
	}
	return nil
}

// ============================================================================
// CoachingMessageRepository — 인메모리 코칭 메시지 저장소
// ============================================================================

// CoachingMessageRepository는 인메모리 코칭 메시지 저장소입니다.
type CoachingMessageRepository struct {
	mu       sync.RWMutex
	messages []*service.CoachingMessage
}

// NewCoachingMessageRepository는 인메모리 CoachingMessageRepository를 생성합니다.
func NewCoachingMessageRepository() *CoachingMessageRepository {
	return &CoachingMessageRepository{
		messages: make([]*service.CoachingMessage, 0),
	}
}

// Save는 코칭 메시지를 저장합니다.
func (r *CoachingMessageRepository) Save(_ context.Context, msg *service.CoachingMessage) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *msg
	if msg.ActionItems != nil {
		cp.ActionItems = make([]string, len(msg.ActionItems))
		copy(cp.ActionItems, msg.ActionItems)
	}
	r.messages = append(r.messages, &cp)
	return nil
}

// ListByUserID는 사용자의 코칭 메시지를 조회합니다.
func (r *CoachingMessageRepository) ListByUserID(_ context.Context, userID string, typeFilter service.CoachingType, limit, offset int32) ([]*service.CoachingMessage, int32, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*service.CoachingMessage
	for _, m := range r.messages {
		if m.UserID != userID {
			continue
		}
		if typeFilter != service.CoachingTypeUnknown && m.CoachingType != typeFilter {
			continue
		}
		filtered = append(filtered, m)
	}

	total := int32(len(filtered))

	// 페이지네이션
	start := int(offset)
	if start >= len(filtered) {
		return nil, total, nil
	}
	end := start + int(limit)
	if end > len(filtered) {
		end = len(filtered)
	}

	result := make([]*service.CoachingMessage, 0, end-start)
	for _, m := range filtered[start:end] {
		cp := *m
		if m.ActionItems != nil {
			cp.ActionItems = make([]string, len(m.ActionItems))
			copy(cp.ActionItems, m.ActionItems)
		}
		result = append(result, &cp)
	}

	return result, total, nil
}

// ============================================================================
// DailyReportRepository — 인메모리 일일 리포트 저장소
// ============================================================================

// DailyReportRepository는 인메모리 일일 리포트 저장소입니다.
type DailyReportRepository struct {
	mu      sync.RWMutex
	reports []*service.DailyHealthReport
}

// NewDailyReportRepository는 인메모리 DailyReportRepository를 생성합니다.
func NewDailyReportRepository() *DailyReportRepository {
	return &DailyReportRepository{
		reports: make([]*service.DailyHealthReport, 0),
	}
}

// Save는 일일 리포트를 저장합니다.
func (r *DailyReportRepository) Save(_ context.Context, report *service.DailyHealthReport) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := copyDailyReport(report)
	r.reports = append(r.reports, cp)
	return nil
}

// GetByUserAndDate는 사용자의 특정 날짜 리포트를 조회합니다.
func (r *DailyReportRepository) GetByUserAndDate(_ context.Context, userID string, date time.Time) (*service.DailyHealthReport, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	dateStr := date.Format("2006-01-02")
	for _, rpt := range r.reports {
		if rpt.UserID == userID && rpt.ReportDate.Format("2006-01-02") == dateStr {
			cp := copyDailyReport(rpt)
			return cp, nil
		}
	}
	return nil, nil
}

// ListByUserAndRange는 사용자의 기간 내 리포트를 조회합니다.
func (r *DailyReportRepository) ListByUserAndRange(_ context.Context, userID string, start, end time.Time) ([]*service.DailyHealthReport, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*service.DailyHealthReport
	for _, rpt := range r.reports {
		if rpt.UserID != userID {
			continue
		}
		rd := rpt.ReportDate
		if (rd.Equal(start) || rd.After(start)) && (rd.Equal(end) || rd.Before(end)) {
			cp := copyDailyReport(rpt)
			result = append(result, cp)
		}
	}
	return result, nil
}

// copyDailyReport는 DailyHealthReport의 깊은 복사를 수행합니다.
func copyDailyReport(src *service.DailyHealthReport) *service.DailyHealthReport {
	cp := *src
	if src.Highlights != nil {
		cp.Highlights = make([]*service.CoachingMessage, len(src.Highlights))
		for i, h := range src.Highlights {
			hcp := *h
			if h.ActionItems != nil {
				hcp.ActionItems = make([]string, len(h.ActionItems))
				copy(hcp.ActionItems, h.ActionItems)
			}
			cp.Highlights[i] = &hcp
		}
	}
	if src.Recommendations != nil {
		cp.Recommendations = make([]string, len(src.Recommendations))
		copy(cp.Recommendations, src.Recommendations)
	}
	return &cp
}
