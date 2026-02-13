// Package postgres는 coaching-service의 PostgreSQL 저장소 구현입니다.
package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/coaching-service/internal/service"
)

// ============================================================================
// HealthGoalRepository — PostgreSQL 기반
// ============================================================================

// HealthGoalRepository는 PostgreSQL 기반 건강 목표 저장소입니다.
type HealthGoalRepository struct {
	pool *pgxpool.Pool
}

// NewHealthGoalRepository는 PostgreSQL HealthGoalRepository를 생성합니다.
func NewHealthGoalRepository(pool *pgxpool.Pool) *HealthGoalRepository {
	return &HealthGoalRepository{pool: pool}
}

// Create는 건강 목표를 생성합니다.
func (r *HealthGoalRepository) Create(ctx context.Context, goal *service.HealthGoal) error {
	const q = `INSERT INTO health_goals
		(id, user_id, category, metric_name, target_value, current_value, unit,
		 progress_pct, status, description, target_date, achieved_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err := r.pool.Exec(ctx, q,
		goal.GoalID,
		goal.UserID,
		goalCategoryToString(goal.Category),
		goal.MetricName,
		goal.TargetValue,
		goal.CurrentValue,
		goal.Unit,
		goal.ProgressPct,
		goalStatusToString(goal.Status),
		goal.Description,
		goal.TargetDate,
		goal.AchievedAt,
		goal.CreatedAt,
	)
	return err
}

// GetByUserID는 사용자의 건강 목표를 조회합니다. statusFilter가 0이면 전체를 반환합니다.
func (r *HealthGoalRepository) GetByUserID(ctx context.Context, userID string, statusFilter service.GoalStatus) ([]*service.HealthGoal, error) {
	var rows pgx.Rows
	var err error

	if statusFilter == service.GoalStatusUnknown {
		const q = `SELECT id, user_id, category, metric_name, target_value, current_value, unit,
			progress_pct, status, COALESCE(description, ''), target_date, achieved_at, created_at
			FROM health_goals WHERE user_id = $1 ORDER BY created_at DESC`
		rows, err = r.pool.Query(ctx, q, userID)
	} else {
		const q = `SELECT id, user_id, category, metric_name, target_value, current_value, unit,
			progress_pct, status, COALESCE(description, ''), target_date, achieved_at, created_at
			FROM health_goals WHERE user_id = $1 AND status = $2 ORDER BY created_at DESC`
		rows, err = r.pool.Query(ctx, q, userID, goalStatusToString(statusFilter))
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []*service.HealthGoal
	for rows.Next() {
		g, err := scanGoal(rows)
		if err != nil {
			return nil, err
		}
		goals = append(goals, g)
	}
	return goals, rows.Err()
}

// GetByID는 목표 ID로 건강 목표를 조회합니다.
func (r *HealthGoalRepository) GetByID(ctx context.Context, id string) (*service.HealthGoal, error) {
	const q = `SELECT id, user_id, category, metric_name, target_value, current_value, unit,
		progress_pct, status, COALESCE(description, ''), target_date, achieved_at, created_at
		FROM health_goals WHERE id = $1`

	var g service.HealthGoal
	var cat, status string
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&g.GoalID, &g.UserID, &cat, &g.MetricName,
		&g.TargetValue, &g.CurrentValue, &g.Unit,
		&g.ProgressPct, &status, &g.Description,
		&g.TargetDate, &g.AchievedAt, &g.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	g.Category = goalCategoryFromString(cat)
	g.Status = goalStatusFromString(status)
	return &g, nil
}

// Update는 건강 목표를 업데이트합니다.
func (r *HealthGoalRepository) Update(ctx context.Context, goal *service.HealthGoal) error {
	const q = `UPDATE health_goals SET
		current_value = $1, progress_pct = $2, status = $3, description = $4,
		target_date = $5, achieved_at = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7`
	_, err := r.pool.Exec(ctx, q,
		goal.CurrentValue,
		goal.ProgressPct,
		goalStatusToString(goal.Status),
		goal.Description,
		goal.TargetDate,
		goal.AchievedAt,
		goal.GoalID,
	)
	return err
}

func scanGoal(rows pgx.Rows) (*service.HealthGoal, error) {
	var g service.HealthGoal
	var cat, status string
	if err := rows.Scan(
		&g.GoalID, &g.UserID, &cat, &g.MetricName,
		&g.TargetValue, &g.CurrentValue, &g.Unit,
		&g.ProgressPct, &status, &g.Description,
		&g.TargetDate, &g.AchievedAt, &g.CreatedAt,
	); err != nil {
		return nil, err
	}
	g.Category = goalCategoryFromString(cat)
	g.Status = goalStatusFromString(status)
	return &g, nil
}

// ============================================================================
// CoachingMessageRepository — PostgreSQL 기반
// ============================================================================

// CoachingMessageRepository는 PostgreSQL 기반 코칭 메시지 저장소입니다.
type CoachingMessageRepository struct {
	pool *pgxpool.Pool
}

// NewCoachingMessageRepository는 PostgreSQL CoachingMessageRepository를 생성합니다.
func NewCoachingMessageRepository(pool *pgxpool.Pool) *CoachingMessageRepository {
	return &CoachingMessageRepository{pool: pool}
}

// Save는 코칭 메시지를 저장합니다.
func (r *CoachingMessageRepository) Save(ctx context.Context, msg *service.CoachingMessage) error {
	actionJSON, err := json.Marshal(msg.ActionItems)
	if err != nil {
		actionJSON = []byte("[]")
	}

	const q = `INSERT INTO coaching_messages
		(id, user_id, coaching_type, title, body, risk_level, action_items, related_metric, related_value, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err = r.pool.Exec(ctx, q,
		msg.MessageID,
		msg.UserID,
		coachingTypeToString(msg.CoachingType),
		msg.Title,
		msg.Body,
		int32(msg.RiskLevel),
		actionJSON,
		msg.RelatedMetric,
		msg.RelatedValue,
		msg.CreatedAt,
	)
	return err
}

// ListByUserID는 사용자의 코칭 메시지를 조회합니다.
func (r *CoachingMessageRepository) ListByUserID(ctx context.Context, userID string, typeFilter service.CoachingType, limit, offset int32) ([]*service.CoachingMessage, int32, error) {
	// 총 개수
	var totalCount int32
	if typeFilter == service.CoachingTypeUnknown {
		const countQ = `SELECT COUNT(*) FROM coaching_messages WHERE user_id = $1`
		if err := r.pool.QueryRow(ctx, countQ, userID).Scan(&totalCount); err != nil {
			return nil, 0, err
		}
	} else {
		const countQ = `SELECT COUNT(*) FROM coaching_messages WHERE user_id = $1 AND coaching_type = $2`
		if err := r.pool.QueryRow(ctx, countQ, userID, coachingTypeToString(typeFilter)).Scan(&totalCount); err != nil {
			return nil, 0, err
		}
	}

	var rows pgx.Rows
	var err error
	if typeFilter == service.CoachingTypeUnknown {
		const q = `SELECT id, user_id, coaching_type, title, body, risk_level, action_items,
			COALESCE(related_metric, ''), related_value, created_at
			FROM coaching_messages WHERE user_id = $1
			ORDER BY created_at DESC LIMIT $2 OFFSET $3`
		rows, err = r.pool.Query(ctx, q, userID, limit, offset)
	} else {
		const q = `SELECT id, user_id, coaching_type, title, body, risk_level, action_items,
			COALESCE(related_metric, ''), related_value, created_at
			FROM coaching_messages WHERE user_id = $1 AND coaching_type = $2
			ORDER BY created_at DESC LIMIT $3 OFFSET $4`
		rows, err = r.pool.Query(ctx, q, userID, coachingTypeToString(typeFilter), limit, offset)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var messages []*service.CoachingMessage
	for rows.Next() {
		var msg service.CoachingMessage
		var ct string
		var rl int32
		var actionJSON []byte
		if err := rows.Scan(
			&msg.MessageID, &msg.UserID, &ct, &msg.Title, &msg.Body,
			&rl, &actionJSON, &msg.RelatedMetric, &msg.RelatedValue, &msg.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		msg.CoachingType = coachingTypeFromString(ct)
		msg.RiskLevel = service.RiskLevel(rl)
		_ = json.Unmarshal(actionJSON, &msg.ActionItems)
		messages = append(messages, &msg)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return messages, totalCount, nil
}

// ============================================================================
// DailyReportRepository — PostgreSQL 기반
// ============================================================================

// DailyReportRepository는 PostgreSQL 기반 일일 리포트 저장소입니다.
type DailyReportRepository struct {
	pool *pgxpool.Pool
}

// NewDailyReportRepository는 PostgreSQL DailyReportRepository를 생성합니다.
func NewDailyReportRepository(pool *pgxpool.Pool) *DailyReportRepository {
	return &DailyReportRepository{pool: pool}
}

// Save는 일일 리포트를 저장합니다.
func (r *DailyReportRepository) Save(ctx context.Context, report *service.DailyHealthReport) error {
	recsJSON, err := json.Marshal(report.Recommendations)
	if err != nil {
		recsJSON = []byte("[]")
	}

	const q = `INSERT INTO daily_health_reports
		(id, user_id, report_date, overall_score, measurements_count, summary, recommendations)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id, report_date) DO UPDATE SET
			overall_score = EXCLUDED.overall_score,
			measurements_count = EXCLUDED.measurements_count,
			summary = EXCLUDED.summary,
			recommendations = EXCLUDED.recommendations`
	_, err = r.pool.Exec(ctx, q,
		report.ReportID,
		report.UserID,
		report.ReportDate,
		report.OverallScore,
		report.MeasurementsCount,
		report.Summary,
		recsJSON,
	)
	return err
}

// GetByUserAndDate는 사용자의 특정 날짜 리포트를 조회합니다.
func (r *DailyReportRepository) GetByUserAndDate(ctx context.Context, userID string, date time.Time) (*service.DailyHealthReport, error) {
	const q = `SELECT id, user_id, report_date, overall_score, measurements_count,
		COALESCE(summary, ''), recommendations
		FROM daily_health_reports WHERE user_id = $1 AND report_date = $2`

	rpt, err := r.scanReport(ctx, q, userID, date)
	if err != nil {
		return nil, err
	}
	return rpt, nil
}

// ListByUserAndRange는 사용자의 기간 내 리포트를 조회합니다.
func (r *DailyReportRepository) ListByUserAndRange(ctx context.Context, userID string, start, end time.Time) ([]*service.DailyHealthReport, error) {
	const q = `SELECT id, user_id, report_date, overall_score, measurements_count,
		COALESCE(summary, ''), recommendations
		FROM daily_health_reports
		WHERE user_id = $1 AND report_date >= $2 AND report_date <= $3
		ORDER BY report_date ASC`

	rows, err := r.pool.Query(ctx, q, userID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []*service.DailyHealthReport
	for rows.Next() {
		var rpt service.DailyHealthReport
		var recsJSON []byte
		if err := rows.Scan(
			&rpt.ReportID, &rpt.UserID, &rpt.ReportDate,
			&rpt.OverallScore, &rpt.MeasurementsCount,
			&rpt.Summary, &recsJSON,
		); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(recsJSON, &rpt.Recommendations)
		reports = append(reports, &rpt)
	}
	return reports, rows.Err()
}

func (r *DailyReportRepository) scanReport(ctx context.Context, query string, args ...any) (*service.DailyHealthReport, error) {
	var rpt service.DailyHealthReport
	var recsJSON []byte
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&rpt.ReportID, &rpt.UserID, &rpt.ReportDate,
		&rpt.OverallScore, &rpt.MeasurementsCount,
		&rpt.Summary, &recsJSON,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	_ = json.Unmarshal(recsJSON, &rpt.Recommendations)
	return &rpt, nil
}

// ============================================================================
// ENUM 변환 헬퍼
// ============================================================================

func goalCategoryToString(c service.GoalCategory) string {
	switch c {
	case service.GoalCategoryBloodGlucose:
		return "BLOOD_GLUCOSE"
	case service.GoalCategoryBloodPressure:
		return "BLOOD_PRESSURE"
	case service.GoalCategoryCholesterol:
		return "CHOLESTEROL"
	case service.GoalCategoryWeight:
		return "WEIGHT"
	case service.GoalCategoryExercise:
		return "EXERCISE"
	case service.GoalCategoryNutrition:
		return "NUTRITION"
	case service.GoalCategorySleep:
		return "SLEEP"
	case service.GoalCategoryStress:
		return "STRESS"
	case service.GoalCategoryCustom:
		return "CUSTOM"
	default:
		return "CUSTOM"
	}
}

func goalCategoryFromString(s string) service.GoalCategory {
	switch s {
	case "BLOOD_GLUCOSE":
		return service.GoalCategoryBloodGlucose
	case "BLOOD_PRESSURE":
		return service.GoalCategoryBloodPressure
	case "CHOLESTEROL":
		return service.GoalCategoryCholesterol
	case "WEIGHT":
		return service.GoalCategoryWeight
	case "EXERCISE":
		return service.GoalCategoryExercise
	case "NUTRITION":
		return service.GoalCategoryNutrition
	case "SLEEP":
		return service.GoalCategorySleep
	case "STRESS":
		return service.GoalCategoryStress
	case "CUSTOM":
		return service.GoalCategoryCustom
	default:
		return service.GoalCategoryUnknown
	}
}

func goalStatusToString(s service.GoalStatus) string {
	switch s {
	case service.GoalStatusActive:
		return "ACTIVE"
	case service.GoalStatusAchieved:
		return "ACHIEVED"
	case service.GoalStatusPaused:
		return "PAUSED"
	case service.GoalStatusCancelled:
		return "CANCELLED"
	default:
		return "ACTIVE"
	}
}

func goalStatusFromString(s string) service.GoalStatus {
	switch s {
	case "ACTIVE":
		return service.GoalStatusActive
	case "ACHIEVED":
		return service.GoalStatusAchieved
	case "PAUSED":
		return service.GoalStatusPaused
	case "CANCELLED":
		return service.GoalStatusCancelled
	default:
		return service.GoalStatusUnknown
	}
}

func coachingTypeToString(ct service.CoachingType) string {
	switch ct {
	case service.CoachingTypeMeasurementFeedback:
		return "MEASUREMENT_FEEDBACK"
	case service.CoachingTypeDailyTip:
		return "DAILY_TIP"
	case service.CoachingTypeGoalProgress:
		return "GOAL_PROGRESS"
	case service.CoachingTypeAlert:
		return "ALERT"
	case service.CoachingTypeMotivation:
		return "MOTIVATION"
	case service.CoachingTypeRecommendation:
		return "RECOMMENDATION"
	default:
		return "DAILY_TIP"
	}
}

func coachingTypeFromString(s string) service.CoachingType {
	switch s {
	case "MEASUREMENT_FEEDBACK":
		return service.CoachingTypeMeasurementFeedback
	case "DAILY_TIP":
		return service.CoachingTypeDailyTip
	case "GOAL_PROGRESS":
		return service.CoachingTypeGoalProgress
	case "ALERT":
		return service.CoachingTypeAlert
	case "MOTIVATION":
		return service.CoachingTypeMotivation
	case "RECOMMENDATION":
		return service.CoachingTypeRecommendation
	default:
		return service.CoachingTypeUnknown
	}
}
