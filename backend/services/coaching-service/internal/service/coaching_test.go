package service

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"
)

// ============================================================================
// 테스트용 Fake 저장소
// ============================================================================

type fakeGoalRepo struct {
	goals map[string]*HealthGoal // key: goalID
}

func newFakeGoalRepo() *fakeGoalRepo {
	return &fakeGoalRepo{goals: make(map[string]*HealthGoal)}
}

func (r *fakeGoalRepo) Create(_ context.Context, goal *HealthGoal) error {
	r.goals[goal.GoalID] = goal
	return nil
}

func (r *fakeGoalRepo) GetByUserID(_ context.Context, userID string, statusFilter GoalStatus) ([]*HealthGoal, error) {
	var result []*HealthGoal
	for _, g := range r.goals {
		if g.UserID != userID {
			continue
		}
		if statusFilter != GoalStatusUnknown && g.Status != statusFilter {
			continue
		}
		result = append(result, g)
	}
	return result, nil
}

func (r *fakeGoalRepo) GetByID(_ context.Context, id string) (*HealthGoal, error) {
	g, ok := r.goals[id]
	if !ok {
		return nil, nil
	}
	return g, nil
}

func (r *fakeGoalRepo) Update(_ context.Context, goal *HealthGoal) error {
	r.goals[goal.GoalID] = goal
	return nil
}

type fakeMsgRepo struct {
	messages []*CoachingMessage
}

func newFakeMsgRepo() *fakeMsgRepo {
	return &fakeMsgRepo{messages: make([]*CoachingMessage, 0)}
}

func (r *fakeMsgRepo) Save(_ context.Context, msg *CoachingMessage) error {
	r.messages = append(r.messages, msg)
	return nil
}

func (r *fakeMsgRepo) ListByUserID(_ context.Context, userID string, typeFilter CoachingType, limit, offset int32) ([]*CoachingMessage, int32, error) {
	var filtered []*CoachingMessage
	for _, m := range r.messages {
		if m.UserID != userID {
			continue
		}
		if typeFilter != CoachingTypeUnknown && m.CoachingType != typeFilter {
			continue
		}
		filtered = append(filtered, m)
	}

	total := int32(len(filtered))
	start := int(offset)
	if start >= len(filtered) {
		return nil, total, nil
	}
	end := start + int(limit)
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[start:end], total, nil
}

type fakeReportRepo struct {
	reports []*DailyHealthReport
}

func newFakeReportRepo() *fakeReportRepo {
	return &fakeReportRepo{reports: make([]*DailyHealthReport, 0)}
}

func (r *fakeReportRepo) Save(_ context.Context, report *DailyHealthReport) error {
	r.reports = append(r.reports, report)
	return nil
}

func (r *fakeReportRepo) GetByUserAndDate(_ context.Context, userID string, date time.Time) (*DailyHealthReport, error) {
	dateStr := date.Format("2006-01-02")
	for _, rpt := range r.reports {
		if rpt.UserID == userID && rpt.ReportDate.Format("2006-01-02") == dateStr {
			return rpt, nil
		}
	}
	return nil, nil
}

func (r *fakeReportRepo) ListByUserAndRange(_ context.Context, userID string, start, end time.Time) ([]*DailyHealthReport, error) {
	var result []*DailyHealthReport
	for _, rpt := range r.reports {
		if rpt.UserID != userID {
			continue
		}
		rd := rpt.ReportDate
		if (rd.Equal(start) || rd.After(start)) && (rd.Equal(end) || rd.Before(end)) {
			result = append(result, rpt)
		}
	}
	return result, nil
}

// ============================================================================
// 테스트 헬퍼
// ============================================================================

func newTestCoachingService() (*CoachingService, *fakeGoalRepo, *fakeMsgRepo, *fakeReportRepo) {
	goalRepo := newFakeGoalRepo()
	msgRepo := newFakeMsgRepo()
	reportRepo := newFakeReportRepo()
	svc := NewCoachingService(zap.NewNop(), goalRepo, msgRepo, reportRepo)
	return svc, goalRepo, msgRepo, reportRepo
}

// ============================================================================
// TestSetHealthGoal — 건강 목표 생성 + 유효성 검증
// ============================================================================

func TestSetHealthGoal(t *testing.T) {
	svc, _, _, _ := newTestCoachingService()
	ctx := context.Background()

	goal, err := svc.SetHealthGoal(ctx, "user-1", GoalCategoryBloodGlucose, "fasting_glucose", 100.0, "mg/dL", "공복혈당 100 이하 유지", time.Now().AddDate(0, 3, 0))
	if err != nil {
		t.Fatalf("SetHealthGoal 실패: %v", err)
	}
	if goal.GoalID == "" {
		t.Error("GoalID가 비어 있습니다")
	}
	if goal.UserID != "user-1" {
		t.Errorf("UserID: got %s, want user-1", goal.UserID)
	}
	if goal.Category != GoalCategoryBloodGlucose {
		t.Errorf("Category: got %d, want %d", goal.Category, GoalCategoryBloodGlucose)
	}
	if goal.MetricName != "fasting_glucose" {
		t.Errorf("MetricName: got %s, want fasting_glucose", goal.MetricName)
	}
	if goal.TargetValue != 100.0 {
		t.Errorf("TargetValue: got %f, want 100.0", goal.TargetValue)
	}
	if goal.ProgressPct != 0 {
		t.Errorf("ProgressPct: got %f, want 0", goal.ProgressPct)
	}
	if goal.Status != GoalStatusActive {
		t.Errorf("Status: got %d, want %d", goal.Status, GoalStatusActive)
	}
	if goal.Unit != "mg/dL" {
		t.Errorf("Unit: got %s, want mg/dL", goal.Unit)
	}

	// 유효성 검증: user_id 누락
	_, err = svc.SetHealthGoal(ctx, "", GoalCategoryBloodGlucose, "fasting_glucose", 100.0, "mg/dL", "", time.Now())
	if err == nil {
		t.Error("user_id 누락 시 에러가 발생해야 합니다")
	}

	// 유효성 검증: 카테고리 누락
	_, err = svc.SetHealthGoal(ctx, "user-1", GoalCategoryUnknown, "fasting_glucose", 100.0, "mg/dL", "", time.Now())
	if err == nil {
		t.Error("카테고리 누락 시 에러가 발생해야 합니다")
	}

	// 유효성 검증: 지표 이름 누락
	_, err = svc.SetHealthGoal(ctx, "user-1", GoalCategoryBloodGlucose, "", 100.0, "mg/dL", "", time.Now())
	if err == nil {
		t.Error("metric_name 누락 시 에러가 발생해야 합니다")
	}

	// 유효성 검증: target_value <= 0
	_, err = svc.SetHealthGoal(ctx, "user-1", GoalCategoryBloodGlucose, "fasting_glucose", 0, "mg/dL", "", time.Now())
	if err == nil {
		t.Error("target_value가 0 이하일 때 에러가 발생해야 합니다")
	}
}

// ============================================================================
// TestGetHealthGoals — 목표 목록 조회
// ============================================================================

func TestGetHealthGoals(t *testing.T) {
	svc, _, _, _ := newTestCoachingService()
	ctx := context.Background()

	// 목표 3개 생성
	_, _ = svc.SetHealthGoal(ctx, "user-goals", GoalCategoryBloodGlucose, "fasting_glucose", 100.0, "mg/dL", "혈당 관리", time.Now().AddDate(0, 3, 0))
	_, _ = svc.SetHealthGoal(ctx, "user-goals", GoalCategoryWeight, "weight", 70.0, "kg", "체중 관리", time.Now().AddDate(0, 6, 0))
	_, _ = svc.SetHealthGoal(ctx, "user-goals", GoalCategoryExercise, "steps", 10000, "steps", "만보 걷기", time.Now().AddDate(0, 1, 0))

	goals, err := svc.GetHealthGoals(ctx, "user-goals", GoalStatusUnknown)
	if err != nil {
		t.Fatalf("GetHealthGoals 실패: %v", err)
	}
	if len(goals) != 3 {
		t.Errorf("목표 수: got %d, want 3", len(goals))
	}

	// user_id 누락 시 에러
	_, err = svc.GetHealthGoals(ctx, "", GoalStatusUnknown)
	if err == nil {
		t.Error("user_id 누락 시 에러가 발생해야 합니다")
	}
}

// ============================================================================
// TestGetHealthGoals_StatusFilter — 상태 필터 테스트
// ============================================================================

func TestGetHealthGoals_StatusFilter(t *testing.T) {
	svc, goalRepo, _, _ := newTestCoachingService()
	ctx := context.Background()

	// 활성 목표 2개, 달성 목표 1개
	g1, _ := svc.SetHealthGoal(ctx, "user-filter", GoalCategoryBloodGlucose, "fasting_glucose", 100.0, "mg/dL", "혈당 관리", time.Now().AddDate(0, 3, 0))
	_, _ = svc.SetHealthGoal(ctx, "user-filter", GoalCategoryWeight, "weight", 70.0, "kg", "체중 관리", time.Now().AddDate(0, 6, 0))

	// 첫 번째 목표를 달성 상태로 변경
	g1.Status = GoalStatusAchieved
	now := time.Now().UTC()
	g1.AchievedAt = &now
	_ = goalRepo.Update(ctx, g1)

	// Active 필터
	activeGoals, err := svc.GetHealthGoals(ctx, "user-filter", GoalStatusActive)
	if err != nil {
		t.Fatalf("GetHealthGoals(Active) 실패: %v", err)
	}
	if len(activeGoals) != 1 {
		t.Errorf("Active 목표 수: got %d, want 1", len(activeGoals))
	}

	// Achieved 필터
	achievedGoals, err := svc.GetHealthGoals(ctx, "user-filter", GoalStatusAchieved)
	if err != nil {
		t.Fatalf("GetHealthGoals(Achieved) 실패: %v", err)
	}
	if len(achievedGoals) != 1 {
		t.Errorf("Achieved 목표 수: got %d, want 1", len(achievedGoals))
	}

	// 전체 (필터 없음)
	allGoals, err := svc.GetHealthGoals(ctx, "user-filter", GoalStatusUnknown)
	if err != nil {
		t.Fatalf("GetHealthGoals(전체) 실패: %v", err)
	}
	if len(allGoals) != 2 {
		t.Errorf("전체 목표 수: got %d, want 2", len(allGoals))
	}
}

// ============================================================================
// TestGenerateCoaching_MeasurementFeedback — 측정 피드백 코칭
// ============================================================================

func TestGenerateCoaching_MeasurementFeedback(t *testing.T) {
	svc, _, _, _ := newTestCoachingService()
	ctx := context.Background()

	msg, err := svc.GenerateCoaching(ctx, "user-feedback", "msmt-001", CoachingTypeMeasurementFeedback)
	if err != nil {
		t.Fatalf("GenerateCoaching(MeasurementFeedback) 실패: %v", err)
	}
	if msg.MessageID == "" {
		t.Error("MessageID가 비어 있습니다")
	}
	if msg.UserID != "user-feedback" {
		t.Errorf("UserID: got %s, want user-feedback", msg.UserID)
	}
	if msg.CoachingType != CoachingTypeMeasurementFeedback {
		t.Errorf("CoachingType: got %d, want %d", msg.CoachingType, CoachingTypeMeasurementFeedback)
	}
	if msg.Title == "" {
		t.Error("Title이 비어 있습니다")
	}
	if msg.Body == "" {
		t.Error("Body가 비어 있습니다")
	}
	if msg.RiskLevel < RiskLevelLow || msg.RiskLevel > RiskLevelCritical {
		t.Errorf("RiskLevel 범위 오류: %d", msg.RiskLevel)
	}
	if len(msg.ActionItems) == 0 {
		t.Error("ActionItems가 비어 있습니다")
	}
	if msg.RelatedMetric == "" {
		t.Error("RelatedMetric이 비어 있습니다")
	}
	if msg.RelatedValue <= 0 {
		t.Errorf("RelatedValue는 양수여야 합니다: got %f", msg.RelatedValue)
	}
}

// ============================================================================
// TestGenerateCoaching_DailyTip — 일일 건강 팁
// ============================================================================

func TestGenerateCoaching_DailyTip(t *testing.T) {
	svc, _, _, _ := newTestCoachingService()
	ctx := context.Background()

	msg, err := svc.GenerateCoaching(ctx, "user-tip", "", CoachingTypeDailyTip)
	if err != nil {
		t.Fatalf("GenerateCoaching(DailyTip) 실패: %v", err)
	}
	if msg.CoachingType != CoachingTypeDailyTip {
		t.Errorf("CoachingType: got %d, want %d", msg.CoachingType, CoachingTypeDailyTip)
	}
	if msg.Title == "" {
		t.Error("Title이 비어 있습니다")
	}
	if msg.Body == "" {
		t.Error("Body가 비어 있습니다")
	}
	if msg.RiskLevel != RiskLevelLow {
		t.Errorf("DailyTip의 RiskLevel은 Low여야 합니다: got %d", msg.RiskLevel)
	}
	if len(msg.ActionItems) == 0 {
		t.Error("ActionItems가 비어 있습니다")
	}
}

// ============================================================================
// TestGenerateCoaching_Motivation — 동기부여 메시지
// ============================================================================

func TestGenerateCoaching_Motivation(t *testing.T) {
	svc, _, _, _ := newTestCoachingService()
	ctx := context.Background()

	msg, err := svc.GenerateCoaching(ctx, "user-motivation", "", CoachingTypeMotivation)
	if err != nil {
		t.Fatalf("GenerateCoaching(Motivation) 실패: %v", err)
	}
	if msg.CoachingType != CoachingTypeMotivation {
		t.Errorf("CoachingType: got %d, want %d", msg.CoachingType, CoachingTypeMotivation)
	}
	if msg.Title == "" {
		t.Error("Title이 비어 있습니다")
	}
	if msg.Body == "" {
		t.Error("Body가 비어 있습니다")
	}
	if msg.RiskLevel != RiskLevelLow {
		t.Errorf("Motivation의 RiskLevel은 Low여야 합니다: got %d", msg.RiskLevel)
	}

	// user_id 누락 시 에러
	_, err = svc.GenerateCoaching(ctx, "", "", CoachingTypeMotivation)
	if err == nil {
		t.Error("user_id 누락 시 에러가 발생해야 합니다")
	}
}

// ============================================================================
// TestListCoachingMessages — 코칭 메시지 이력 조회
// ============================================================================

func TestListCoachingMessages(t *testing.T) {
	svc, _, _, _ := newTestCoachingService()
	ctx := context.Background()

	// 여러 타입의 코칭 메시지 생성
	_, _ = svc.GenerateCoaching(ctx, "user-list", "msmt-1", CoachingTypeMeasurementFeedback)
	_, _ = svc.GenerateCoaching(ctx, "user-list", "", CoachingTypeDailyTip)
	_, _ = svc.GenerateCoaching(ctx, "user-list", "", CoachingTypeMotivation)
	_, _ = svc.GenerateCoaching(ctx, "user-list", "", CoachingTypeDailyTip)

	// 전체 조회
	messages, total, err := svc.ListCoachingMessages(ctx, "user-list", CoachingTypeUnknown, 10, 0)
	if err != nil {
		t.Fatalf("ListCoachingMessages 실패: %v", err)
	}
	if total != 4 {
		t.Errorf("전체 메시지 수: got %d, want 4", total)
	}
	if len(messages) != 4 {
		t.Errorf("반환 메시지 수: got %d, want 4", len(messages))
	}

	// 타입 필터 (DailyTip만)
	tipMessages, tipTotal, err := svc.ListCoachingMessages(ctx, "user-list", CoachingTypeDailyTip, 10, 0)
	if err != nil {
		t.Fatalf("ListCoachingMessages(DailyTip 필터) 실패: %v", err)
	}
	if tipTotal != 2 {
		t.Errorf("DailyTip 메시지 수: got %d, want 2", tipTotal)
	}
	if len(tipMessages) != 2 {
		t.Errorf("DailyTip 반환 수: got %d, want 2", len(tipMessages))
	}

	// 페이지네이션
	page1, _, err := svc.ListCoachingMessages(ctx, "user-list", CoachingTypeUnknown, 2, 0)
	if err != nil {
		t.Fatalf("ListCoachingMessages(page1) 실패: %v", err)
	}
	if len(page1) != 2 {
		t.Errorf("Page1 메시지 수: got %d, want 2", len(page1))
	}

	page2, _, err := svc.ListCoachingMessages(ctx, "user-list", CoachingTypeUnknown, 2, 2)
	if err != nil {
		t.Fatalf("ListCoachingMessages(page2) 실패: %v", err)
	}
	if len(page2) != 2 {
		t.Errorf("Page2 메시지 수: got %d, want 2", len(page2))
	}
}

// ============================================================================
// TestGenerateDailyReport — 일일 건강 리포트 생성
// ============================================================================

func TestGenerateDailyReport(t *testing.T) {
	svc, _, _, _ := newTestCoachingService()
	ctx := context.Background()

	today := time.Now().UTC()
	report, err := svc.GenerateDailyReport(ctx, "user-daily", today)
	if err != nil {
		t.Fatalf("GenerateDailyReport 실패: %v", err)
	}
	if report.ReportID == "" {
		t.Error("ReportID가 비어 있습니다")
	}
	if report.UserID != "user-daily" {
		t.Errorf("UserID: got %s, want user-daily", report.UserID)
	}
	if report.OverallScore < 60 || report.OverallScore > 100 {
		t.Errorf("OverallScore 범위 오류: %f (expected 60-100)", report.OverallScore)
	}
	if report.MeasurementsCount < 1 || report.MeasurementsCount > 8 {
		t.Errorf("MeasurementsCount 범위 오류: %d (expected 1-8)", report.MeasurementsCount)
	}
	if len(report.Highlights) < 2 {
		t.Errorf("Highlights 수: got %d, want >= 2", len(report.Highlights))
	}
	if report.Summary == "" {
		t.Error("Summary가 비어 있습니다")
	}
	if len(report.Recommendations) < 2 {
		t.Errorf("Recommendations 수: got %d, want >= 2", len(report.Recommendations))
	}

	// user_id 누락 시 에러
	_, err = svc.GenerateDailyReport(ctx, "", today)
	if err == nil {
		t.Error("user_id 누락 시 에러가 발생해야 합니다")
	}
}

// ============================================================================
// TestGetWeeklyReport — 주간 건강 리포트 조회
// ============================================================================

func TestGetWeeklyReport(t *testing.T) {
	svc, _, _, _ := newTestCoachingService()
	ctx := context.Background()

	// 이번 주 월요일로 설정
	now := time.Now().UTC()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	weekStart := time.Date(now.Year(), now.Month(), now.Day()-(weekday-1), 0, 0, 0, 0, time.UTC)

	report, err := svc.GetWeeklyReport(ctx, "user-weekly", weekStart)
	if err != nil {
		t.Fatalf("GetWeeklyReport 실패: %v", err)
	}
	if report.ReportID == "" {
		t.Error("ReportID가 비어 있습니다")
	}
	if report.UserID != "user-weekly" {
		t.Errorf("UserID: got %s, want user-weekly", report.UserID)
	}
	if report.WeekStart.IsZero() {
		t.Error("WeekStart가 비어 있습니다")
	}
	if report.WeekEnd.IsZero() {
		t.Error("WeekEnd가 비어 있습니다")
	}
	if report.AverageScore < 0 || report.AverageScore > 100 {
		t.Errorf("AverageScore 범위 오류: %f", report.AverageScore)
	}
	validTrends := map[string]bool{"improving": true, "stable": true, "declining": true}
	if !validTrends[report.ScoreTrend] {
		t.Errorf("ScoreTrend 값 오류: %s", report.ScoreTrend)
	}
	if report.TotalMeasurements <= 0 {
		t.Errorf("TotalMeasurements는 양수여야 합니다: got %d", report.TotalMeasurements)
	}
	if len(report.DailyReports) == 0 {
		t.Error("DailyReports가 비어 있습니다")
	}
	if report.WeeklySummary == "" {
		t.Error("WeeklySummary가 비어 있습니다")
	}
	if len(report.KeyInsights) == 0 {
		t.Error("KeyInsights가 비어 있습니다")
	}

	// user_id 누락 시 에러
	_, err = svc.GetWeeklyReport(ctx, "", weekStart)
	if err == nil {
		t.Error("user_id 누락 시 에러가 발생해야 합니다")
	}
}

// ============================================================================
// TestGetRecommendations — 개인화 추천 조회
// ============================================================================

func TestGetRecommendations(t *testing.T) {
	svc, _, _, _ := newTestCoachingService()
	ctx := context.Background()

	recs, err := svc.GetRecommendations(ctx, "user-recs", RecommendationTypeUnknown, 10)
	if err != nil {
		t.Fatalf("GetRecommendations 실패: %v", err)
	}
	if len(recs) != 5 {
		t.Errorf("추천 수: got %d, want 5", len(recs))
	}

	// 각 추천 항목 검증
	for _, r := range recs {
		if r.RecommendationID == "" {
			t.Error("RecommendationID가 비어 있습니다")
		}
		if r.Title == "" {
			t.Error("Title이 비어 있습니다")
		}
		if r.Description == "" {
			t.Error("Description이 비어 있습니다")
		}
		if r.Reason == "" {
			t.Error("Reason이 비어 있습니다")
		}
		if len(r.ActionSteps) == 0 {
			t.Error("ActionSteps가 비어 있습니다")
		}
		if r.Type == RecommendationTypeUnknown {
			t.Error("Type이 Unknown입니다")
		}
	}

	// 유형별 확인
	typeSet := make(map[RecommendationType]bool)
	for _, r := range recs {
		typeSet[r.Type] = true
	}
	expectedTypes := []RecommendationType{
		RecommendationTypeFood,
		RecommendationTypeExercise,
		RecommendationTypeSupplement,
		RecommendationTypeLifestyle,
		RecommendationTypeCheckup,
	}
	for _, et := range expectedTypes {
		if !typeSet[et] {
			t.Errorf("추천 유형 %d가 누락되었습니다", et)
		}
	}

	// user_id 누락 시 에러
	_, err = svc.GetRecommendations(ctx, "", RecommendationTypeUnknown, 10)
	if err == nil {
		t.Error("user_id 누락 시 에러가 발생해야 합니다")
	}
}

// ============================================================================
// TestGetRecommendations_TypeFilter — 추천 유형 필터
// ============================================================================

func TestGetRecommendations_TypeFilter(t *testing.T) {
	svc, _, _, _ := newTestCoachingService()
	ctx := context.Background()

	// FOOD 필터
	foodRecs, err := svc.GetRecommendations(ctx, "user-filter-recs", RecommendationTypeFood, 10)
	if err != nil {
		t.Fatalf("GetRecommendations(Food) 실패: %v", err)
	}
	if len(foodRecs) != 1 {
		t.Errorf("Food 추천 수: got %d, want 1", len(foodRecs))
	}
	if len(foodRecs) > 0 && foodRecs[0].Type != RecommendationTypeFood {
		t.Errorf("Food 추천 타입: got %d, want %d", foodRecs[0].Type, RecommendationTypeFood)
	}

	// EXERCISE 필터
	exerciseRecs, err := svc.GetRecommendations(ctx, "user-filter-recs", RecommendationTypeExercise, 10)
	if err != nil {
		t.Fatalf("GetRecommendations(Exercise) 실패: %v", err)
	}
	if len(exerciseRecs) != 1 {
		t.Errorf("Exercise 추천 수: got %d, want 1", len(exerciseRecs))
	}

	// CHECKUP 필터
	checkupRecs, err := svc.GetRecommendations(ctx, "user-filter-recs", RecommendationTypeCheckup, 10)
	if err != nil {
		t.Fatalf("GetRecommendations(Checkup) 실패: %v", err)
	}
	if len(checkupRecs) != 1 {
		t.Errorf("Checkup 추천 수: got %d, want 1", len(checkupRecs))
	}

	// limit 적용 테스트
	limitedRecs, err := svc.GetRecommendations(ctx, "user-filter-recs", RecommendationTypeUnknown, 2)
	if err != nil {
		t.Fatalf("GetRecommendations(limit=2) 실패: %v", err)
	}
	if len(limitedRecs) != 2 {
		t.Errorf("제한 추천 수: got %d, want 2", len(limitedRecs))
	}
}
