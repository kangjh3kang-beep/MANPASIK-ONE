// Package handler는 coaching-service의 gRPC 핸들러입니다.
package handler

import (
	"context"

	"github.com/manpasik/backend/services/coaching-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CoachingHandler는 CoachingService gRPC 서버를 구현합니다.
type CoachingHandler struct {
	v1.UnimplementedCoachingServiceServer
	svc *service.CoachingService
	log *zap.Logger
}

// NewCoachingHandler는 CoachingHandler를 생성합니다.
func NewCoachingHandler(svc *service.CoachingService, log *zap.Logger) *CoachingHandler {
	return &CoachingHandler{svc: svc, log: log}
}

// ============================================================================
// SetHealthGoal — 건강 목표 설정
// ============================================================================

// SetHealthGoal은 건강 목표 설정 RPC입니다.
func (h *CoachingHandler) SetHealthGoal(ctx context.Context, req *v1.SetHealthGoalRequest) (*v1.HealthGoal, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}
	if req.MetricName == "" {
		return nil, status.Error(codes.InvalidArgument, "metric_name은 필수입니다")
	}
	if req.TargetValue <= 0 {
		return nil, status.Error(codes.InvalidArgument, "target_value는 0보다 커야 합니다")
	}

	var targetDate = req.TargetDate.AsTime()

	goal, err := h.svc.SetHealthGoal(
		ctx,
		req.UserId,
		protoGoalCategoryToService(req.Category),
		req.MetricName,
		req.TargetValue,
		req.Unit,
		req.Description,
		targetDate,
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	return goalToProto(goal), nil
}

// ============================================================================
// GetHealthGoals — 건강 목표 조회
// ============================================================================

// GetHealthGoals는 건강 목표 조회 RPC입니다.
func (h *CoachingHandler) GetHealthGoals(ctx context.Context, req *v1.GetHealthGoalsRequest) (*v1.GetHealthGoalsResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	goals, err := h.svc.GetHealthGoals(ctx, req.UserId, protoGoalStatusToService(req.StatusFilter))
	if err != nil {
		return nil, toGRPC(err)
	}

	protoGoals := make([]*v1.HealthGoal, 0, len(goals))
	for _, g := range goals {
		protoGoals = append(protoGoals, goalToProto(g))
	}

	return &v1.GetHealthGoalsResponse{Goals: protoGoals}, nil
}

// ============================================================================
// GenerateCoaching — AI 코칭 메시지 생성
// ============================================================================

// GenerateCoaching은 AI 코칭 메시지 생성 RPC입니다.
func (h *CoachingHandler) GenerateCoaching(ctx context.Context, req *v1.GenerateCoachingRequest) (*v1.CoachingMessage, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	msg, err := h.svc.GenerateCoaching(
		ctx,
		req.UserId,
		req.MeasurementId,
		protoCoachingTypeToService(req.CoachingType),
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	return coachingMessageToProto(msg), nil
}

// ============================================================================
// ListCoachingMessages — 코칭 메시지 이력 조회
// ============================================================================

// ListCoachingMessages는 코칭 메시지 이력 조회 RPC입니다.
func (h *CoachingHandler) ListCoachingMessages(ctx context.Context, req *v1.ListCoachingMessagesRequest) (*v1.ListCoachingMessagesResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	messages, total, err := h.svc.ListCoachingMessages(
		ctx,
		req.UserId,
		protoCoachingTypeToService(req.TypeFilter),
		req.Limit,
		req.Offset,
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	protoMsgs := make([]*v1.CoachingMessage, 0, len(messages))
	for _, m := range messages {
		protoMsgs = append(protoMsgs, coachingMessageToProto(m))
	}

	return &v1.ListCoachingMessagesResponse{
		Messages:   protoMsgs,
		TotalCount: total,
	}, nil
}

// ============================================================================
// GenerateDailyReport — 일일 건강 리포트 생성
// ============================================================================

// GenerateDailyReport는 일일 건강 리포트 생성 RPC입니다.
func (h *CoachingHandler) GenerateDailyReport(ctx context.Context, req *v1.GenerateDailyReportRequest) (*v1.DailyHealthReport, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	var date = req.Date.AsTime()

	report, err := h.svc.GenerateDailyReport(ctx, req.UserId, date)
	if err != nil {
		return nil, toGRPC(err)
	}

	return dailyReportToProto(report), nil
}

// ============================================================================
// GetWeeklyReport — 주간 건강 리포트 조회
// ============================================================================

// GetWeeklyReport는 주간 건강 리포트 조회 RPC입니다.
func (h *CoachingHandler) GetWeeklyReport(ctx context.Context, req *v1.GetWeeklyReportRequest) (*v1.WeeklyHealthReport, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	var weekStart = req.WeekStart.AsTime()

	report, err := h.svc.GetWeeklyReport(ctx, req.UserId, weekStart)
	if err != nil {
		return nil, toGRPC(err)
	}

	return weeklyReportToProto(report), nil
}

// ============================================================================
// GetRecommendations — 개인화 추천 조회
// ============================================================================

// GetRecommendations는 개인화 추천 조회 RPC입니다.
func (h *CoachingHandler) GetRecommendations(ctx context.Context, req *v1.GetRecommendationsRequest) (*v1.GetRecommendationsResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	recs, err := h.svc.GetRecommendations(
		ctx,
		req.UserId,
		protoRecommendationTypeToService(req.TypeFilter),
		req.Limit,
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	protoRecs := make([]*v1.Recommendation, 0, len(recs))
	for _, r := range recs {
		protoRecs = append(protoRecs, recommendationToProto(r))
	}

	return &v1.GetRecommendationsResponse{Recommendations: protoRecs}, nil
}

// ============================================================================
// 헬퍼 함수 — Proto ↔ Service 변환
// ============================================================================

func goalToProto(g *service.HealthGoal) *v1.HealthGoal {
	pg := &v1.HealthGoal{
		GoalId:       g.GoalID,
		UserId:       g.UserID,
		Category:     serviceGoalCategoryToProto(g.Category),
		MetricName:   g.MetricName,
		TargetValue:  g.TargetValue,
		CurrentValue: g.CurrentValue,
		Unit:         g.Unit,
		ProgressPct:  g.ProgressPct,
		Status:       serviceGoalStatusToProto(g.Status),
		Description:  g.Description,
		CreatedAt:    timestamppb.New(g.CreatedAt),
		TargetDate:   timestamppb.New(g.TargetDate),
	}
	if g.AchievedAt != nil {
		pg.AchievedAt = timestamppb.New(*g.AchievedAt)
	}
	return pg
}

func coachingMessageToProto(m *service.CoachingMessage) *v1.CoachingMessage {
	return &v1.CoachingMessage{
		MessageId:     m.MessageID,
		UserId:        m.UserID,
		CoachingType:  serviceCoachingTypeToProto(m.CoachingType),
		Title:         m.Title,
		Body:          m.Body,
		RiskLevel:     serviceRiskLevelToProto(m.RiskLevel),
		ActionItems:   m.ActionItems,
		RelatedMetric: m.RelatedMetric,
		RelatedValue:  m.RelatedValue,
		CreatedAt:     timestamppb.New(m.CreatedAt),
	}
}

func dailyReportToProto(r *service.DailyHealthReport) *v1.DailyHealthReport {
	highlights := make([]*v1.CoachingMessage, 0, len(r.Highlights))
	for _, h := range r.Highlights {
		highlights = append(highlights, coachingMessageToProto(h))
	}
	return &v1.DailyHealthReport{
		ReportId:          r.ReportID,
		UserId:            r.UserID,
		ReportDate:        timestamppb.New(r.ReportDate),
		OverallScore:      r.OverallScore,
		MeasurementsCount: r.MeasurementsCount,
		Highlights:        highlights,
		Summary:           r.Summary,
		Recommendations:   r.Recommendations,
	}
}

func weeklyReportToProto(r *service.WeeklyHealthReport) *v1.WeeklyHealthReport {
	dailyReports := make([]*v1.DailyHealthReport, 0, len(r.DailyReports))
	for _, dr := range r.DailyReports {
		dailyReports = append(dailyReports, dailyReportToProto(dr))
	}
	return &v1.WeeklyHealthReport{
		ReportId:          r.ReportID,
		UserId:            r.UserID,
		WeekStart:         timestamppb.New(r.WeekStart),
		WeekEnd:           timestamppb.New(r.WeekEnd),
		AverageScore:      r.AverageScore,
		ScoreTrend:        r.ScoreTrend,
		TotalMeasurements: r.TotalMeasurements,
		GoalsAchieved:     r.GoalsAchieved,
		GoalsActive:       r.GoalsActive,
		DailyReports:      dailyReports,
		WeeklySummary:     r.WeeklySummary,
		KeyInsights:       r.KeyInsights,
	}
}

func recommendationToProto(r *service.Recommendation) *v1.Recommendation {
	return &v1.Recommendation{
		RecommendationId: r.RecommendationID,
		Type:             serviceRecommendationTypeToProto(r.Type),
		Title:            r.Title,
		Description:      r.Description,
		Reason:           r.Reason,
		Priority:         serviceRiskLevelToProto(r.Priority),
		ActionSteps:      r.ActionSteps,
		RelatedMetric:    r.RelatedMetric,
		CreatedAt:        timestamppb.New(r.CreatedAt),
	}
}

// --- GoalCategory 변환 ---

func protoGoalCategoryToService(c v1.GoalCategory) service.GoalCategory {
	switch c {
	case v1.GoalCategory_GOAL_CATEGORY_BLOOD_GLUCOSE:
		return service.GoalCategoryBloodGlucose
	case v1.GoalCategory_GOAL_CATEGORY_BLOOD_PRESSURE:
		return service.GoalCategoryBloodPressure
	case v1.GoalCategory_GOAL_CATEGORY_CHOLESTEROL:
		return service.GoalCategoryCholesterol
	case v1.GoalCategory_GOAL_CATEGORY_WEIGHT:
		return service.GoalCategoryWeight
	case v1.GoalCategory_GOAL_CATEGORY_EXERCISE:
		return service.GoalCategoryExercise
	case v1.GoalCategory_GOAL_CATEGORY_NUTRITION:
		return service.GoalCategoryNutrition
	case v1.GoalCategory_GOAL_CATEGORY_SLEEP:
		return service.GoalCategorySleep
	case v1.GoalCategory_GOAL_CATEGORY_STRESS:
		return service.GoalCategoryStress
	case v1.GoalCategory_GOAL_CATEGORY_CUSTOM:
		return service.GoalCategoryCustom
	default:
		return service.GoalCategoryUnknown
	}
}

func serviceGoalCategoryToProto(c service.GoalCategory) v1.GoalCategory {
	switch c {
	case service.GoalCategoryBloodGlucose:
		return v1.GoalCategory_GOAL_CATEGORY_BLOOD_GLUCOSE
	case service.GoalCategoryBloodPressure:
		return v1.GoalCategory_GOAL_CATEGORY_BLOOD_PRESSURE
	case service.GoalCategoryCholesterol:
		return v1.GoalCategory_GOAL_CATEGORY_CHOLESTEROL
	case service.GoalCategoryWeight:
		return v1.GoalCategory_GOAL_CATEGORY_WEIGHT
	case service.GoalCategoryExercise:
		return v1.GoalCategory_GOAL_CATEGORY_EXERCISE
	case service.GoalCategoryNutrition:
		return v1.GoalCategory_GOAL_CATEGORY_NUTRITION
	case service.GoalCategorySleep:
		return v1.GoalCategory_GOAL_CATEGORY_SLEEP
	case service.GoalCategoryStress:
		return v1.GoalCategory_GOAL_CATEGORY_STRESS
	case service.GoalCategoryCustom:
		return v1.GoalCategory_GOAL_CATEGORY_CUSTOM
	default:
		return v1.GoalCategory_GOAL_CATEGORY_UNKNOWN
	}
}

// --- GoalStatus 변환 ---

func protoGoalStatusToService(s v1.GoalStatus) service.GoalStatus {
	switch s {
	case v1.GoalStatus_GOAL_STATUS_ACTIVE:
		return service.GoalStatusActive
	case v1.GoalStatus_GOAL_STATUS_ACHIEVED:
		return service.GoalStatusAchieved
	case v1.GoalStatus_GOAL_STATUS_PAUSED:
		return service.GoalStatusPaused
	case v1.GoalStatus_GOAL_STATUS_CANCELLED:
		return service.GoalStatusCancelled
	default:
		return service.GoalStatusUnknown
	}
}

func serviceGoalStatusToProto(s service.GoalStatus) v1.GoalStatus {
	switch s {
	case service.GoalStatusActive:
		return v1.GoalStatus_GOAL_STATUS_ACTIVE
	case service.GoalStatusAchieved:
		return v1.GoalStatus_GOAL_STATUS_ACHIEVED
	case service.GoalStatusPaused:
		return v1.GoalStatus_GOAL_STATUS_PAUSED
	case service.GoalStatusCancelled:
		return v1.GoalStatus_GOAL_STATUS_CANCELLED
	default:
		return v1.GoalStatus_GOAL_STATUS_UNKNOWN
	}
}

// --- CoachingType 변환 ---

func protoCoachingTypeToService(t v1.CoachingType) service.CoachingType {
	switch t {
	case v1.CoachingType_COACHING_TYPE_MEASUREMENT_FEEDBACK:
		return service.CoachingTypeMeasurementFeedback
	case v1.CoachingType_COACHING_TYPE_DAILY_TIP:
		return service.CoachingTypeDailyTip
	case v1.CoachingType_COACHING_TYPE_GOAL_PROGRESS:
		return service.CoachingTypeGoalProgress
	case v1.CoachingType_COACHING_TYPE_ALERT:
		return service.CoachingTypeAlert
	case v1.CoachingType_COACHING_TYPE_MOTIVATION:
		return service.CoachingTypeMotivation
	case v1.CoachingType_COACHING_TYPE_RECOMMENDATION:
		return service.CoachingTypeRecommendation
	default:
		return service.CoachingTypeUnknown
	}
}

func serviceCoachingTypeToProto(t service.CoachingType) v1.CoachingType {
	switch t {
	case service.CoachingTypeMeasurementFeedback:
		return v1.CoachingType_COACHING_TYPE_MEASUREMENT_FEEDBACK
	case service.CoachingTypeDailyTip:
		return v1.CoachingType_COACHING_TYPE_DAILY_TIP
	case service.CoachingTypeGoalProgress:
		return v1.CoachingType_COACHING_TYPE_GOAL_PROGRESS
	case service.CoachingTypeAlert:
		return v1.CoachingType_COACHING_TYPE_ALERT
	case service.CoachingTypeMotivation:
		return v1.CoachingType_COACHING_TYPE_MOTIVATION
	case service.CoachingTypeRecommendation:
		return v1.CoachingType_COACHING_TYPE_RECOMMENDATION
	default:
		return v1.CoachingType_COACHING_TYPE_UNKNOWN
	}
}

// --- RiskLevel 변환 ---

func serviceRiskLevelToProto(r service.RiskLevel) v1.RiskLevel {
	switch r {
	case service.RiskLevelLow:
		return v1.RiskLevel_RISK_LEVEL_LOW
	case service.RiskLevelModerate:
		return v1.RiskLevel_RISK_LEVEL_MODERATE
	case service.RiskLevelHigh:
		return v1.RiskLevel_RISK_LEVEL_HIGH
	case service.RiskLevelCritical:
		return v1.RiskLevel_RISK_LEVEL_CRITICAL
	default:
		return v1.RiskLevel_RISK_LEVEL_UNSPECIFIED
	}
}

// --- RecommendationType 변환 ---

func protoRecommendationTypeToService(t v1.RecommendationType) service.RecommendationType {
	switch t {
	case v1.RecommendationType_RECOMMENDATION_TYPE_FOOD:
		return service.RecommendationTypeFood
	case v1.RecommendationType_RECOMMENDATION_TYPE_EXERCISE:
		return service.RecommendationTypeExercise
	case v1.RecommendationType_RECOMMENDATION_TYPE_SUPPLEMENT:
		return service.RecommendationTypeSupplement
	case v1.RecommendationType_RECOMMENDATION_TYPE_LIFESTYLE:
		return service.RecommendationTypeLifestyle
	case v1.RecommendationType_RECOMMENDATION_TYPE_CHECKUP:
		return service.RecommendationTypeCheckup
	default:
		return service.RecommendationTypeUnknown
	}
}

func serviceRecommendationTypeToProto(t service.RecommendationType) v1.RecommendationType {
	switch t {
	case service.RecommendationTypeFood:
		return v1.RecommendationType_RECOMMENDATION_TYPE_FOOD
	case service.RecommendationTypeExercise:
		return v1.RecommendationType_RECOMMENDATION_TYPE_EXERCISE
	case service.RecommendationTypeSupplement:
		return v1.RecommendationType_RECOMMENDATION_TYPE_SUPPLEMENT
	case service.RecommendationTypeLifestyle:
		return v1.RecommendationType_RECOMMENDATION_TYPE_LIFESTYLE
	case service.RecommendationTypeCheckup:
		return v1.RecommendationType_RECOMMENDATION_TYPE_CHECKUP
	default:
		return v1.RecommendationType_RECOMMENDATION_TYPE_UNKNOWN
	}
}

// --- 에러 변환 ---

// toGRPC는 AppError를 gRPC status로 변환합니다.
func toGRPC(err error) error {
	if err == nil {
		return nil
	}
	if ae, ok := err.(*apperrors.AppError); ok {
		return ae.ToGRPC()
	}
	if s, ok := status.FromError(err); ok {
		return s.Err()
	}
	return status.Error(codes.Internal, "내부 오류가 발생했습니다")
}
