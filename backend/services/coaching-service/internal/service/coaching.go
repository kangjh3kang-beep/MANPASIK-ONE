// Package service는 coaching-service의 비즈니스 로직을 구현합니다.
package service

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/manpasik/backend/shared/errors"
	"go.uber.org/zap"
)

// ============================================================================
// 도메인 타입 (Domain Types)
// ============================================================================

// GoalCategory는 건강 목표 카테고리입니다.
type GoalCategory int32

const (
	GoalCategoryUnknown       GoalCategory = 0
	GoalCategoryBloodGlucose  GoalCategory = 1
	GoalCategoryBloodPressure GoalCategory = 2
	GoalCategoryCholesterol   GoalCategory = 3
	GoalCategoryWeight        GoalCategory = 4
	GoalCategoryExercise      GoalCategory = 5
	GoalCategoryNutrition     GoalCategory = 6
	GoalCategorySleep         GoalCategory = 7
	GoalCategoryStress        GoalCategory = 8
	GoalCategoryCustom        GoalCategory = 9
)

// GoalStatus는 건강 목표 상태입니다.
type GoalStatus int32

const (
	GoalStatusUnknown   GoalStatus = 0
	GoalStatusActive    GoalStatus = 1
	GoalStatusAchieved  GoalStatus = 2
	GoalStatusPaused    GoalStatus = 3
	GoalStatusCancelled GoalStatus = 4
)

// CoachingType은 코칭 메시지 유형입니다.
type CoachingType int32

const (
	CoachingTypeUnknown             CoachingType = 0
	CoachingTypeMeasurementFeedback CoachingType = 1
	CoachingTypeDailyTip            CoachingType = 2
	CoachingTypeGoalProgress        CoachingType = 3
	CoachingTypeAlert               CoachingType = 4
	CoachingTypeMotivation          CoachingType = 5
	CoachingTypeRecommendation      CoachingType = 6
)

// RiskLevel은 위험 수준입니다.
type RiskLevel int32

const (
	RiskLevelUnspecified RiskLevel = 0
	RiskLevelLow        RiskLevel = 1
	RiskLevelModerate   RiskLevel = 2
	RiskLevelHigh       RiskLevel = 3
	RiskLevelCritical   RiskLevel = 4
)

// RecommendationType은 추천 유형입니다.
type RecommendationType int32

const (
	RecommendationTypeUnknown    RecommendationType = 0
	RecommendationTypeFood       RecommendationType = 1
	RecommendationTypeExercise   RecommendationType = 2
	RecommendationTypeSupplement RecommendationType = 3
	RecommendationTypeLifestyle  RecommendationType = 4
	RecommendationTypeCheckup    RecommendationType = 5
)

// HealthGoal은 건강 목표 엔티티입니다.
type HealthGoal struct {
	GoalID      string
	UserID      string
	Category    GoalCategory
	MetricName  string
	TargetValue float64
	CurrentValue float64
	Unit        string
	ProgressPct float64
	Status      GoalStatus
	Description string
	CreatedAt   time.Time
	TargetDate  time.Time
	AchievedAt  *time.Time
}

// CoachingMessage는 코칭 메시지 엔티티입니다.
type CoachingMessage struct {
	MessageID     string
	UserID        string
	CoachingType  CoachingType
	Title         string
	Body          string
	RiskLevel     RiskLevel
	ActionItems   []string
	RelatedMetric string
	RelatedValue  float64
	CreatedAt     time.Time
}

// DailyHealthReport는 일일 건강 리포트입니다.
type DailyHealthReport struct {
	ReportID          string
	UserID            string
	ReportDate        time.Time
	OverallScore      float64
	MeasurementsCount int32
	Highlights        []*CoachingMessage
	Summary           string
	Recommendations   []string
}

// WeeklyHealthReport는 주간 건강 리포트입니다.
type WeeklyHealthReport struct {
	ReportID          string
	UserID            string
	WeekStart         time.Time
	WeekEnd           time.Time
	AverageScore      float64
	ScoreTrend        string // "improving", "stable", "declining"
	TotalMeasurements int32
	GoalsAchieved     int32
	GoalsActive       int32
	DailyReports      []*DailyHealthReport
	WeeklySummary     string
	KeyInsights       []string
}

// Recommendation은 개인화 추천 엔티티입니다.
type Recommendation struct {
	RecommendationID string
	Type             RecommendationType
	Title            string
	Description      string
	Reason           string
	Priority         RiskLevel
	ActionSteps      []string
	RelatedMetric    string
	CreatedAt        time.Time
}

// ============================================================================
// 저장소 인터페이스 (Repository Interfaces)
// ============================================================================

// HealthGoalRepository는 건강 목표 저장소 인터페이스입니다.
type HealthGoalRepository interface {
	Create(ctx context.Context, goal *HealthGoal) error
	GetByUserID(ctx context.Context, userID string, statusFilter GoalStatus) ([]*HealthGoal, error)
	GetByID(ctx context.Context, id string) (*HealthGoal, error)
	Update(ctx context.Context, goal *HealthGoal) error
}

// CoachingMessageRepository는 코칭 메시지 저장소 인터페이스입니다.
type CoachingMessageRepository interface {
	Save(ctx context.Context, msg *CoachingMessage) error
	ListByUserID(ctx context.Context, userID string, typeFilter CoachingType, limit, offset int32) ([]*CoachingMessage, int32, error)
}

// DailyReportRepository는 일일 리포트 저장소 인터페이스입니다.
type DailyReportRepository interface {
	Save(ctx context.Context, report *DailyHealthReport) error
	GetByUserAndDate(ctx context.Context, userID string, date time.Time) (*DailyHealthReport, error)
	ListByUserAndRange(ctx context.Context, userID string, start, end time.Time) ([]*DailyHealthReport, error)
}

// HealthDataProvider는 외부 건강 데이터를 제공하는 인터페이스입니다.
// ai-inference-service 등의 건강 점수/카테고리 데이터를 코칭에 활용합니다.
type HealthDataProvider interface {
	// GetUserHealthScore는 사용자의 건강 점수와 카테고리별 점수를 반환합니다.
	GetUserHealthScore(ctx context.Context, userID string) (*UserHealthData, error)
}

// UserHealthData는 외부 건강 데이터 제공자로부터 받은 건강 정보입니다.
type UserHealthData struct {
	OverallScore   float64
	CategoryScores map[string]float64 // cardiovascular, metabolic, nutritional, fitness
	Trend          string             // improving, stable, declining
}

// ============================================================================
// 코칭 서비스 (Coaching Service)
// ============================================================================

// CoachingService는 AI 건강 코칭 비즈니스 로직입니다.
type CoachingService struct {
	logger       *zap.Logger
	goalRepo     HealthGoalRepository
	msgRepo      CoachingMessageRepository
	reportRepo   DailyReportRepository
	healthData   HealthDataProvider // nil이면 시뮬레이션 사용
	rng          *rand.Rand
}

// NewCoachingService는 새 CoachingService를 생성합니다.
func NewCoachingService(
	logger *zap.Logger,
	goalRepo HealthGoalRepository,
	msgRepo CoachingMessageRepository,
	reportRepo DailyReportRepository,
) *CoachingService {
	return &CoachingService{
		logger:     logger,
		goalRepo:   goalRepo,
		msgRepo:    msgRepo,
		reportRepo: reportRepo,
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// SetHealthDataProvider는 건강 데이터 제공자를 주입합니다.
func (s *CoachingService) SetHealthDataProvider(provider HealthDataProvider) {
	s.healthData = provider
}

// ============================================================================
// SetHealthGoal — 건강 목표 설정
// ============================================================================

// SetHealthGoal은 새 건강 목표를 생성합니다.
func (s *CoachingService) SetHealthGoal(ctx context.Context, userID string, category GoalCategory, metricName string, targetValue float64, unit, description string, targetDate time.Time) (*HealthGoal, error) {
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}
	if category == GoalCategoryUnknown {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "목표 카테고리는 필수입니다")
	}
	if metricName == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "지표 이름은 필수입니다")
	}
	if targetValue <= 0 {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "목표 값은 0보다 커야 합니다")
	}

	now := time.Now().UTC()
	goal := &HealthGoal{
		GoalID:       uuid.New().String(),
		UserID:       userID,
		Category:     category,
		MetricName:   metricName,
		TargetValue:  targetValue,
		CurrentValue: 0,
		Unit:         unit,
		ProgressPct:  0,
		Status:       GoalStatusActive,
		Description:  description,
		CreatedAt:    now,
		TargetDate:   targetDate,
	}

	if err := s.goalRepo.Create(ctx, goal); err != nil {
		s.logger.Error("건강 목표 생성 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "건강 목표 생성에 실패했습니다")
	}

	s.logger.Info("건강 목표 생성 완료",
		zap.String("user_id", userID),
		zap.String("goal_id", goal.GoalID),
		zap.Int32("category", int32(category)),
	)
	return goal, nil
}

// ============================================================================
// GetHealthGoals — 건강 목표 조회
// ============================================================================

// GetHealthGoals는 사용자의 건강 목표 목록을 반환합니다.
func (s *CoachingService) GetHealthGoals(ctx context.Context, userID string, statusFilter GoalStatus) ([]*HealthGoal, error) {
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}

	goals, err := s.goalRepo.GetByUserID(ctx, userID, statusFilter)
	if err != nil {
		s.logger.Error("건강 목표 조회 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "건강 목표 조회에 실패했습니다")
	}

	return goals, nil
}

// ============================================================================
// GenerateCoaching — AI 코칭 메시지 생성
// ============================================================================

// GenerateCoaching은 AI 코칭 메시지를 생성합니다.
func (s *CoachingService) GenerateCoaching(ctx context.Context, userID, measurementID string, coachingType CoachingType) (*CoachingMessage, error) {
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}

	// 코칭 타입이 지정되지 않으면 자동 선택
	if coachingType == CoachingTypeUnknown {
		if measurementID != "" {
			coachingType = CoachingTypeMeasurementFeedback
		} else {
			coachingType = CoachingTypeDailyTip
		}
	}

	var msg *CoachingMessage
	switch coachingType {
	case CoachingTypeMeasurementFeedback:
		msg = s.generateMeasurementFeedback(userID, measurementID)
	case CoachingTypeDailyTip:
		msg = s.generateDailyTip(userID)
	case CoachingTypeGoalProgress:
		msg = s.generateGoalProgress(ctx, userID)
	case CoachingTypeAlert:
		msg = s.generateAlert(userID)
	case CoachingTypeMotivation:
		msg = s.generateMotivation(userID)
	case CoachingTypeRecommendation:
		msg = s.generateRecommendationMessage(userID)
	default:
		msg = s.generateDailyTip(userID)
	}

	if err := s.msgRepo.Save(ctx, msg); err != nil {
		s.logger.Error("코칭 메시지 저장 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "코칭 메시지 저장에 실패했습니다")
	}

	s.logger.Info("코칭 메시지 생성 완료",
		zap.String("user_id", userID),
		zap.String("message_id", msg.MessageID),
		zap.Int32("coaching_type", int32(coachingType)),
	)
	return msg, nil
}

// ============================================================================
// ListCoachingMessages — 코칭 메시지 이력 조회
// ============================================================================

// ListCoachingMessages는 코칭 메시지 이력을 반환합니다.
func (s *CoachingService) ListCoachingMessages(ctx context.Context, userID string, typeFilter CoachingType, limit, offset int32) ([]*CoachingMessage, int32, error) {
	if userID == "" {
		return nil, 0, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}
	if limit <= 0 {
		limit = 20
	}

	messages, total, err := s.msgRepo.ListByUserID(ctx, userID, typeFilter, limit, offset)
	if err != nil {
		s.logger.Error("코칭 메시지 조회 실패", zap.Error(err))
		return nil, 0, apperrors.New(apperrors.ErrInternal, "코칭 메시지 조회에 실패했습니다")
	}

	return messages, total, nil
}

// ============================================================================
// GenerateDailyReport — 일일 건강 리포트 생성
// ============================================================================

// GenerateDailyReport는 일일 건강 리포트를 생성합니다.
func (s *CoachingService) GenerateDailyReport(ctx context.Context, userID string, date time.Time) (*DailyHealthReport, error) {
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}

	if date.IsZero() {
		date = time.Now().UTC()
	}
	// 날짜 정규화 (시간 제거)
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	// 건강 데이터 제공자에서 실제 점수를 가져옴 (없으면 시뮬레이션)
	var overallScore float64
	measurementsCount := int32(s.rng.Intn(8) + 1) // 1~8 측정
	if s.healthData != nil {
		hd, err := s.healthData.GetUserHealthScore(ctx, userID)
		if err == nil && hd != nil {
			overallScore = math.Round(hd.OverallScore*10) / 10
		} else {
			overallScore = math.Round((60+s.rng.Float64()*40)*10) / 10
		}
	} else {
		overallScore = math.Round((60+s.rng.Float64()*40)*10) / 10
	}

	highlights := make([]*CoachingMessage, 0)
	// 주요 하이라이트 2~3개 생성
	highlightCount := s.rng.Intn(2) + 2
	for i := 0; i < highlightCount; i++ {
		hl := s.generateMeasurementFeedback(userID, fmt.Sprintf("msmt_%d_%d", date.Day(), i))
		highlights = append(highlights, hl)
	}

	recommendations := s.generateDailyRecommendations()

	summary := s.generateDailySummary(overallScore, measurementsCount)

	report := &DailyHealthReport{
		ReportID:          uuid.New().String(),
		UserID:            userID,
		ReportDate:        date,
		OverallScore:      overallScore,
		MeasurementsCount: measurementsCount,
		Highlights:        highlights,
		Summary:           summary,
		Recommendations:   recommendations,
	}

	if err := s.reportRepo.Save(ctx, report); err != nil {
		s.logger.Error("일일 리포트 저장 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "일일 리포트 저장에 실패했습니다")
	}

	s.logger.Info("일일 건강 리포트 생성 완료",
		zap.String("user_id", userID),
		zap.String("report_id", report.ReportID),
	)
	return report, nil
}

// ============================================================================
// GetWeeklyReport — 주간 건강 리포트 조회
// ============================================================================

// GetWeeklyReport는 주간 건강 리포트를 생성합니다.
func (s *CoachingService) GetWeeklyReport(ctx context.Context, userID string, weekStart time.Time) (*WeeklyHealthReport, error) {
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}

	if weekStart.IsZero() {
		// 이번 주 월요일로 설정
		now := time.Now().UTC()
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		weekStart = time.Date(now.Year(), now.Month(), now.Day()-(weekday-1), 0, 0, 0, 0, time.UTC)
	}
	weekEnd := weekStart.AddDate(0, 0, 6)

	// 기존 일일 리포트 조회 또는 시뮬레이션 생성
	dailyReports, err := s.reportRepo.ListByUserAndRange(ctx, userID, weekStart, weekEnd)
	if err != nil {
		s.logger.Error("일일 리포트 범위 조회 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "일일 리포트 조회에 실패했습니다")
	}

	// 기존 리포트가 없으면 7일치 시뮬레이션 생성
	if len(dailyReports) == 0 {
		for i := 0; i < 7; i++ {
			d := weekStart.AddDate(0, 0, i)
			report, genErr := s.GenerateDailyReport(ctx, userID, d)
			if genErr != nil {
				continue
			}
			dailyReports = append(dailyReports, report)
		}
	}

	// 집계
	var totalScore float64
	var totalMeasurements int32
	for _, dr := range dailyReports {
		totalScore += dr.OverallScore
		totalMeasurements += dr.MeasurementsCount
	}

	avgScore := float64(0)
	if len(dailyReports) > 0 {
		avgScore = math.Round(totalScore/float64(len(dailyReports))*10) / 10
	}

	// 목표 집계
	goals, _ := s.goalRepo.GetByUserID(ctx, userID, GoalStatusUnknown)
	var goalsAchieved, goalsActive int32
	for _, g := range goals {
		switch g.Status {
		case GoalStatusAchieved:
			goalsAchieved++
		case GoalStatusActive:
			goalsActive++
		}
	}

	scoreTrend := s.evaluateScoreTrend(dailyReports)

	weeklyReport := &WeeklyHealthReport{
		ReportID:          uuid.New().String(),
		UserID:            userID,
		WeekStart:         weekStart,
		WeekEnd:           weekEnd,
		AverageScore:      avgScore,
		ScoreTrend:        scoreTrend,
		TotalMeasurements: totalMeasurements,
		GoalsAchieved:     goalsAchieved,
		GoalsActive:       goalsActive,
		DailyReports:      dailyReports,
		WeeklySummary:     s.generateWeeklySummary(avgScore, totalMeasurements, scoreTrend),
		KeyInsights:       s.generateKeyInsights(avgScore, totalMeasurements, goalsActive),
	}

	s.logger.Info("주간 건강 리포트 생성 완료",
		zap.String("user_id", userID),
		zap.String("report_id", weeklyReport.ReportID),
	)
	return weeklyReport, nil
}

// ============================================================================
// GetRecommendations — 개인화 추천 조회
// ============================================================================

// GetRecommendations는 개인화 추천 목록을 반환합니다.
func (s *CoachingService) GetRecommendations(_ context.Context, userID string, typeFilter RecommendationType, limit int32) ([]*Recommendation, error) {
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}
	if limit <= 0 {
		limit = 10
	}

	allRecs := s.generateAllRecommendations(userID)

	// 타입 필터 적용
	if typeFilter != RecommendationTypeUnknown {
		filtered := make([]*Recommendation, 0)
		for _, r := range allRecs {
			if r.Type == typeFilter {
				filtered = append(filtered, r)
			}
		}
		allRecs = filtered
	}

	// 제한 적용
	if int32(len(allRecs)) > limit {
		allRecs = allRecs[:limit]
	}

	return allRecs, nil
}

// ============================================================================
// 내부 헬퍼 — 시뮬레이션 코칭 엔진
// ============================================================================

// --- 측정 피드백 ---

func (s *CoachingService) generateMeasurementFeedback(userID, measurementID string) *CoachingMessage {
	metrics := []struct {
		name     string
		unit     string
		refMin   float64
		refMax   float64
		refRange string
	}{
		{"공복혈당", "mg/dL", 70, 100, "70-100 mg/dL"},
		{"총콜레스테롤", "mg/dL", 0, 200, "< 200 mg/dL"},
		{"혈압(수축기)", "mmHg", 90, 120, "90-120 mmHg"},
		{"HbA1c", "%", 4.0, 5.6, "4.0-5.6%"},
		{"중성지방", "mg/dL", 0, 150, "< 150 mg/dL"},
	}

	m := metrics[s.rng.Intn(len(metrics))]
	value := m.refMin + s.rng.Float64()*(m.refMax*1.4-m.refMin)
	value = math.Round(value*10) / 10

	risk := RiskLevelLow
	title := fmt.Sprintf("%s 측정 결과 분석", m.name)
	var body string
	var actionItems []string

	switch {
	case value <= m.refMax*0.9:
		risk = RiskLevelLow
		body = fmt.Sprintf("%s 수치가 %.1f%s으로 정상 범위 내에 있습니다. 현재 건강 관리를 잘 유지하고 계십니다.", m.name, value, m.unit)
		actionItems = []string{
			"현재 식단과 운동 패턴을 유지하세요",
			"정기적인 측정을 계속해 주세요",
		}
	case value <= m.refMax:
		risk = RiskLevelLow
		body = fmt.Sprintf("%s 수치가 %.1f%s으로 정상 범위입니다. 기준 범위는 %s입니다.", m.name, value, m.unit, m.refRange)
		actionItems = []string{
			"현재 생활 습관을 유지하세요",
			"균형 잡힌 식단을 유지하세요",
		}
	case value <= m.refMax*1.15:
		risk = RiskLevelModerate
		body = fmt.Sprintf("%s 수치가 %.1f%s으로 기준 범위(%s)보다 약간 높습니다. 생활 습관 개선을 권장합니다.", m.name, value, m.unit, m.refRange)
		actionItems = []string{
			"식이 조절을 시작해 보세요",
			"규칙적인 유산소 운동을 30분 이상 하세요",
			"1~2주 후 재측정을 권장합니다",
		}
	default:
		risk = RiskLevelHigh
		body = fmt.Sprintf("⚠️ %s 수치가 %.1f%s으로 기준 범위(%s)를 초과했습니다. 전문의 상담을 권장합니다.", m.name, value, m.unit, m.refRange)
		actionItems = []string{
			"가까운 시일 내에 전문의 상담을 받으세요",
			"과식과 고지방 음식을 피하세요",
			"충분한 수면과 스트레스 관리를 하세요",
			"3일 이내 재측정을 해 주세요",
		}
	}

	return &CoachingMessage{
		MessageID:     uuid.New().String(),
		UserID:        userID,
		CoachingType:  CoachingTypeMeasurementFeedback,
		Title:         title,
		Body:          body,
		RiskLevel:     risk,
		ActionItems:   actionItems,
		RelatedMetric: m.name,
		RelatedValue:  value,
		CreatedAt:     time.Now().UTC(),
	}
}

// --- 일일 건강 팁 ---

var dailyTips = []struct {
	title       string
	body        string
	actionItems []string
}{
	{
		"혈당 관리 팁",
		"식후 30분 이내에 가벼운 산책을 하면 혈당 스파이크를 줄일 수 있습니다. 하루 15분 걷기만으로도 혈당 관리에 큰 도움이 됩니다.",
		[]string{"식후 15-30분 산책하기", "계단 이용하기", "식사 속도 늦추기"},
	},
	{
		"수분 섭취의 중요성",
		"충분한 수분 섭취는 혈액 순환을 돕고, 노폐물 배출을 촉진합니다. 하루 8잔(약 2L)의 물을 마시는 것을 목표로 하세요.",
		[]string{"아침 기상 직후 물 한 잔 마시기", "식사 30분 전에 물 마시기", "카페인 음료 대신 물 선택하기"},
	},
	{
		"수면과 건강",
		"수면은 면역력 강화와 호르몬 균형에 필수적입니다. 매일 같은 시간에 취침하고 기상하면 체내 리듬이 안정됩니다.",
		[]string{"취침 1시간 전 스크린 사용 줄이기", "침실 온도를 18-20°C로 유지하기", "카페인은 오후 2시 이전에만 섭취하기"},
	},
	{
		"스트레스 관리법",
		"만성 스트레스는 혈당, 혈압, 콜레스테롤 수치를 악화시킵니다. 하루 10분의 명상이나 심호흡으로 스트레스를 관리하세요.",
		[]string{"4-7-8 호흡법 실천하기", "하루 10분 명상 또는 요가하기", "좋아하는 취미 활동 시간 만들기"},
	},
	{
		"식이섬유 섭취",
		"식이섬유는 혈당 상승을 완만하게 하고, 콜레스테롤 수치를 낮추는 데 도움을 줍니다. 채소, 과일, 통곡물을 매 끼니에 포함하세요.",
		[]string{"매 끼니 채소를 먼저 먹기", "흰 빵 대신 통밀빵 선택하기", "간식으로 과일이나 견과류 먹기"},
	},
	{
		"규칙적인 운동의 효과",
		"주 5회 30분 이상의 중강도 유산소 운동은 심혈관 질환 위험을 50% 줄일 수 있습니다. 걷기, 수영, 자전거 등이 좋습니다.",
		[]string{"매일 같은 시간에 운동하기", "운동 전후로 스트레칭하기", "점진적으로 운동 강도 높이기"},
	},
	{
		"나트륨 섭취 줄이기",
		"과도한 나트륨 섭취는 혈압 상승의 주요 원인입니다. 하루 나트륨 섭취량을 2,000mg 이하로 유지하세요.",
		[]string{"음식 조리 시 소금 줄이기", "가공식품 섭취 줄이기", "외식 시 싱겁게 주문하기"},
	},
	{
		"건강한 장내 환경",
		"장내 미생물 균형은 면역력, 정신 건강, 대사 기능에 영향을 줍니다. 발효 식품과 프로바이오틱스를 규칙적으로 섭취하세요.",
		[]string{"김치, 요거트 등 발효식품 섭취하기", "프리바이오틱스가 풍부한 식품 먹기", "과도한 항생제 사용 피하기"},
	},
}

func (s *CoachingService) generateDailyTip(userID string) *CoachingMessage {
	tip := dailyTips[s.rng.Intn(len(dailyTips))]
	return &CoachingMessage{
		MessageID:    uuid.New().String(),
		UserID:       userID,
		CoachingType: CoachingTypeDailyTip,
		Title:        tip.title,
		Body:         tip.body,
		RiskLevel:    RiskLevelLow,
		ActionItems:  tip.actionItems,
		CreatedAt:    time.Now().UTC(),
	}
}

// --- 목표 진행 알림 ---

func (s *CoachingService) generateGoalProgress(ctx context.Context, userID string) *CoachingMessage {
	goals, _ := s.goalRepo.GetByUserID(ctx, userID, GoalStatusActive)

	if len(goals) == 0 {
		return &CoachingMessage{
			MessageID:    uuid.New().String(),
			UserID:       userID,
			CoachingType: CoachingTypeGoalProgress,
			Title:        "건강 목표를 설정해 보세요",
			Body:         "아직 설정된 건강 목표가 없습니다. 혈당, 혈압, 체중 등 관리하고 싶은 목표를 설정하면 맞춤형 코칭을 받을 수 있습니다.",
			RiskLevel:    RiskLevelLow,
			ActionItems: []string{
				"건강 목표 설정 메뉴에서 새 목표를 추가하세요",
				"작은 목표부터 시작하는 것이 좋습니다",
			},
			CreatedAt: time.Now().UTC(),
		}
	}

	goal := goals[s.rng.Intn(len(goals))]
	// 시뮬레이션: 목표 진행률 업데이트
	progress := math.Round(s.rng.Float64() * 100 * 10) / 10

	var body string
	risk := RiskLevelLow
	var actionItems []string

	switch {
	case progress >= 90:
		body = fmt.Sprintf("🎉 %s 목표 달성이 거의 완료되었습니다! 진행률 %.1f%%입니다. 마지막까지 화이팅!", goal.Description, progress)
		actionItems = []string{"현재 페이스를 유지하세요", "달성 후 새로운 목표를 설정하세요"}
	case progress >= 50:
		body = fmt.Sprintf("👍 %s 목표를 향해 순조롭게 진행 중입니다. 현재 진행률 %.1f%%입니다.", goal.Description, progress)
		actionItems = []string{"꾸준함이 핵심입니다", "주간 리포트에서 트렌드를 확인하세요"}
	case progress >= 20:
		body = fmt.Sprintf("💪 %s 목표 진행률이 %.1f%%입니다. 조금 더 노력해 볼까요?", goal.Description, progress)
		risk = RiskLevelModerate
		actionItems = []string{"목표 달성 전략을 재점검해 보세요", "작은 단계별 목표를 설정하세요"}
	default:
		body = fmt.Sprintf("📋 %s 목표 진행률이 %.1f%%입니다. 시작이 반입니다! 오늘부터 실천해 보세요.", goal.Description, progress)
		risk = RiskLevelModerate
		actionItems = []string{"오늘 할 수 있는 작은 행동 하나를 시작하세요", "알림 설정으로 규칙적인 관리를 해 보세요"}
	}

	return &CoachingMessage{
		MessageID:     uuid.New().String(),
		UserID:        userID,
		CoachingType:  CoachingTypeGoalProgress,
		Title:         fmt.Sprintf("%s 목표 진행 상황", goal.MetricName),
		Body:          body,
		RiskLevel:     risk,
		ActionItems:   actionItems,
		RelatedMetric: goal.MetricName,
		RelatedValue:  progress,
		CreatedAt:     time.Now().UTC(),
	}
}

// --- 건강 경고 ---

func (s *CoachingService) generateAlert(userID string) *CoachingMessage {
	alerts := []struct {
		title       string
		body        string
		risk        RiskLevel
		metric      string
		value       float64
		actionItems []string
	}{
		{
			"혈당 수치 경고",
			"⚠️ 최근 공복혈당 수치가 126mg/dL 이상으로 당뇨병 기준을 초과했습니다. 전문의 상담이 필요합니다.",
			RiskLevelHigh,
			"공복혈당",
			135.0,
			[]string{"가능한 빨리 전문의를 방문하세요", "고탄수화물 식품 섭취를 제한하세요", "규칙적인 혈당 모니터링을 하세요"},
		},
		{
			"혈압 주의보",
			"⚠️ 수축기 혈압이 140mmHg를 초과하여 고혈압 1기 수준입니다. 생활 습관 개선과 의료 상담을 권장합니다.",
			RiskLevelHigh,
			"수축기혈압",
			145.0,
			[]string{"소금 섭취를 줄이세요", "매일 30분 이상 유산소 운동을 하세요", "이번 주 내에 의료 상담을 받으세요"},
		},
		{
			"콜레스테롤 경고",
			"총콜레스테롤이 240mg/dL 이상으로 높은 수준입니다. 심혈관 질환 위험이 증가할 수 있으므로 주의가 필요합니다.",
			RiskLevelHigh,
			"총콜레스테롤",
			255.0,
			[]string{"포화지방 섭취를 줄이세요", "식이섬유가 풍부한 음식을 섭취하세요", "전문의 상담 후 약물 치료 여부를 결정하세요"},
		},
		{
			"체중 변화 감지",
			"최근 1주일간 체중이 2kg 이상 급격히 변화했습니다. 갑작스러운 체중 변화는 건강 이상 신호일 수 있습니다.",
			RiskLevelModerate,
			"체중",
			72.5,
			[]string{"식사 일지를 기록해 보세요", "수분 섭취량을 확인하세요", "지속되면 의료 상담을 받으세요"},
		},
	}

	alert := alerts[s.rng.Intn(len(alerts))]
	return &CoachingMessage{
		MessageID:     uuid.New().String(),
		UserID:        userID,
		CoachingType:  CoachingTypeAlert,
		Title:         alert.title,
		Body:          alert.body,
		RiskLevel:     alert.risk,
		ActionItems:   alert.actionItems,
		RelatedMetric: alert.metric,
		RelatedValue:  alert.value,
		CreatedAt:     time.Now().UTC(),
	}
}

// --- 동기부여 메시지 ---

func (s *CoachingService) generateMotivation(userID string) *CoachingMessage {
	motivations := []struct {
		title string
		body  string
	}{
		{
			"꾸준함이 건강의 비결입니다",
			"매일 조금씩 건강을 관리하는 습관이 큰 변화를 만듭니다. 오늘도 건강 측정을 완료하셨네요! 이 작은 실천이 쌓여 큰 건강을 지켜줍니다. 💪",
		},
		{
			"건강은 최고의 투자입니다",
			"건강에 투자하는 시간은 절대 낭비가 아닙니다. 규칙적인 건강 체크는 미래의 나를 위한 최고의 선물입니다. 오늘도 건강한 하루 되세요! 🌟",
		},
		{
			"작은 변화가 큰 차이를 만듭니다",
			"어제보다 한 잔 더 마신 물, 10분 더 걸은 거리, 조금 일찍 잠든 시간... 이 작은 변화들이 모여 당신의 건강을 지켜줍니다. 오늘도 하나의 좋은 습관을 실천해 보세요! 🎯",
		},
		{
			"당신의 건강 관리에 박수를 보냅니다",
			"정기적으로 건강을 모니터링하는 것만으로도 대단한 일입니다. 많은 사람들이 시작조차 하지 못하는데, 당신은 이미 행동으로 옮기고 있습니다. 자신에게 칭찬 한마디 해 주세요! 👏",
		},
		{
			"오늘의 건강이 내일의 행복입니다",
			"건강한 몸은 행복한 삶의 기초입니다. 오늘 측정한 데이터는 더 건강한 내일을 만들기 위한 소중한 정보입니다. 앞으로도 함께 건강을 지켜나가요! 🏃‍♂️",
		},
	}

	m := motivations[s.rng.Intn(len(motivations))]
	return &CoachingMessage{
		MessageID:    uuid.New().String(),
		UserID:       userID,
		CoachingType: CoachingTypeMotivation,
		Title:        m.title,
		Body:         m.body,
		RiskLevel:    RiskLevelLow,
		ActionItems: []string{
			"오늘의 건강 측정을 완료하세요",
			"가족이나 친구와 건강 목표를 공유해 보세요",
		},
		CreatedAt: time.Now().UTC(),
	}
}

// --- 추천 메시지 ---

func (s *CoachingService) generateRecommendationMessage(userID string) *CoachingMessage {
	recs := s.generateAllRecommendations(userID)
	if len(recs) == 0 {
		return s.generateDailyTip(userID)
	}
	rec := recs[s.rng.Intn(len(recs))]
	return &CoachingMessage{
		MessageID:     uuid.New().String(),
		UserID:        userID,
		CoachingType:  CoachingTypeRecommendation,
		Title:         rec.Title,
		Body:          rec.Description,
		RiskLevel:     rec.Priority,
		ActionItems:   rec.ActionSteps,
		RelatedMetric: rec.RelatedMetric,
		CreatedAt:     time.Now().UTC(),
	}
}

// --- 개인화 추천 생성 ---

func (s *CoachingService) generateAllRecommendations(userID string) []*Recommendation {
	now := time.Now().UTC()

	// 건강 데이터가 있으면 카테고리 점수 기반 동적 우선순위 적용
	var healthData *UserHealthData
	if s.healthData != nil {
		hd, err := s.healthData.GetUserHealthScore(context.Background(), userID)
		if err == nil && hd != nil {
			healthData = hd
		}
	}

	recs := s.buildRecommendationPool(now, healthData)

	// 건강 데이터 기반 우선순위 재정렬
	if healthData != nil {
		s.prioritizeRecommendations(recs, healthData)
	}

	return recs
}

// buildRecommendationPool은 건강 데이터에 따라 맞춤형 추천 풀을 생성합니다.
func (s *CoachingService) buildRecommendationPool(now time.Time, healthData *UserHealthData) []*Recommendation {
	recs := make([]*Recommendation, 0, 5)

	// 식단 추천 (대사 점수 기반)
	foodRec := &Recommendation{
		RecommendationID: uuid.New().String(),
		Type:             RecommendationTypeFood,
		RelatedMetric:    "공복혈당",
		CreatedAt:        now,
	}
	if healthData != nil && healthData.CategoryScores["metabolic"] < 60 {
		foodRec.Title = "혈당 관리 긴급 식단 조정"
		foodRec.Description = "대사 지표가 낮습니다. 즉시 저GI 식품 위주로 식단을 전환하고, 정제 탄수화물을 최소화하세요."
		foodRec.Reason = fmt.Sprintf("대사 건강 점수 %.0f점으로 적극적인 식이 관리가 필요합니다", healthData.CategoryScores["metabolic"])
		foodRec.Priority = RiskLevelHigh
		foodRec.ActionSteps = []string{"정제 탄수화물(흰빵, 흰쌀) → 통곡물로 교체", "매 끼니 식이섬유 풍부한 채소 반 접시 이상", "당분 음료 중단, 물 또는 무가당 차로 대체", "식후 15분 이내 가벼운 산책"}
	} else {
		foodRec.Title = "혈당 관리를 위한 식단 추천"
		foodRec.Description = "혈당 수치 안정화를 위해 저GI 식품 위주의 식단을 권장합니다."
		foodRec.Reason = "최근 혈당 측정 결과를 기반으로 식이 관리가 필요합니다"
		foodRec.Priority = RiskLevelModerate
		foodRec.ActionSteps = []string{"아침: 귀리죽 + 계란 + 채소 샐러드", "점심: 현미밥 + 두부구이 + 나물", "저녁: 잡곡밥 + 생선구이 + 미역국", "간식: 견과류 한 줌 또는 사과 반 개"}
	}
	recs = append(recs, foodRec)

	// 운동 추천 (심혈관+체력 점수 기반)
	exerciseRec := &Recommendation{
		RecommendationID: uuid.New().String(),
		Type:             RecommendationTypeExercise,
		RelatedMetric:    "심혈관",
		CreatedAt:        now,
	}
	if healthData != nil && healthData.CategoryScores["cardiovascular"] < 60 {
		exerciseRec.Title = "심혈관 건강 개선 집중 프로그램"
		exerciseRec.Description = "심혈관 지표가 낮습니다. 유산소 운동을 매일 실시하고, 나트륨 섭취를 하루 2g 이하로 제한하세요."
		exerciseRec.Reason = fmt.Sprintf("심혈관 건강 점수 %.0f점으로 적극적 운동이 필요합니다", healthData.CategoryScores["cardiovascular"])
		exerciseRec.Priority = RiskLevelHigh
		exerciseRec.ActionSteps = []string{"매일: 빠르게 걷기 30분 (심박수 120-140)", "주 3회: 수영 또는 자전거 40분", "매일: 혈압 모니터링 (아침/저녁)", "나트륨 섭취 제한: 국물, 찌개 줄이기"}
	} else {
		exerciseRec.Title = "심혈관 건강을 위한 운동 프로그램"
		exerciseRec.Description = "유산소 운동은 심혈관 건강과 혈당 조절에 효과적입니다."
		exerciseRec.Reason = "정기적인 운동은 건강 점수를 개선하는 가장 효과적인 방법입니다"
		exerciseRec.Priority = RiskLevelModerate
		exerciseRec.ActionSteps = []string{"월수금: 빠르게 걷기 30분", "화목: 가벼운 근력 운동 20분", "주말: 자전거 타기 또는 수영 40분", "매일: 스트레칭 10분"}
	}
	recs = append(recs, exerciseRec)

	// 영양 보충 추천 (영양 점수 기반)
	suppRec := &Recommendation{
		RecommendationID: uuid.New().String(),
		Type:             RecommendationTypeSupplement,
		RelatedMetric:    "영양",
		CreatedAt:        now,
	}
	if healthData != nil && healthData.CategoryScores["nutritional"] < 60 {
		suppRec.Title = "영양 결핍 보충 프로그램"
		suppRec.Description = "영양 지표가 낮습니다. 전문의 상담 후 부족 영양소를 적극 보충하세요."
		suppRec.Reason = fmt.Sprintf("영양 점수 %.0f점으로 영양 보충이 시급합니다", healthData.CategoryScores["nutritional"])
		suppRec.Priority = RiskLevelHigh
		suppRec.ActionSteps = []string{"철분 + 비타민 C 복합 보충제 (빈혈 예방)", "비타민 D: 하루 2,000IU (결핍 위험)", "단백질: 체중 kg당 1.2g 이상 섭취", "혈액 검사로 정확한 결핍 확인 필수"}
	} else {
		suppRec.Title = "영양 보충제 추천"
		suppRec.Description = "균형 잡힌 영양 섭취를 위해 부족할 수 있는 영양소의 보충을 권장합니다."
		suppRec.Reason = "현대인의 식습관으로 부족하기 쉬운 영양소를 보충합니다"
		suppRec.Priority = RiskLevelLow
		suppRec.ActionSteps = []string{"비타민 D: 하루 1,000IU (실내 생활이 많은 경우)", "오메가-3: 하루 1,000mg (생선 섭취 부족 시)", "마그네슘: 하루 300-400mg (스트레스가 많은 경우)", "의사와 상담 후 복용을 시작하세요"}
	}
	recs = append(recs, suppRec)

	// 생활습관 추천
	lifestyleRec := &Recommendation{
		RecommendationID: uuid.New().String(),
		Type:             RecommendationTypeLifestyle,
		Title:            "스트레스 관리와 수면 개선",
		Description:      "만성 스트레스는 혈당, 혈압, 면역력에 부정적 영향을 미칩니다.",
		Reason:           "스트레스와 수면은 전반적 건강 지표에 영향을 미칩니다",
		Priority:         RiskLevelModerate,
		ActionSteps:      []string{"매일 같은 시간에 취침하기 (23시 이전)", "취침 전 1시간은 디지털 기기 사용 중단", "하루 10분 명상 또는 심호흡 연습", "주 1회 이상 자연에서 산책하기"},
		RelatedMetric:    "스트레스",
		CreatedAt:        now,
	}
	recs = append(recs, lifestyleRec)

	// 건강검진 추천 (전체 점수 기반)
	checkupRec := &Recommendation{
		RecommendationID: uuid.New().String(),
		Type:             RecommendationTypeCheckup,
		RelatedMetric:    "종합",
		CreatedAt:        now,
	}
	if healthData != nil && healthData.OverallScore < 60 {
		checkupRec.Title = "긴급 건강검진 권장"
		checkupRec.Description = "종합 건강 점수가 낮습니다. 가능한 빨리 종합 건강검진을 받으시기 바랍니다."
		checkupRec.Reason = fmt.Sprintf("종합 건강 점수 %.0f점으로 전문 검진이 필요합니다", healthData.OverallScore)
		checkupRec.Priority = RiskLevelHigh
		checkupRec.ActionSteps = []string{"2주 이내 종합 건강검진 예약", "혈액 검사 (공복혈당, HbA1c, 지질 패널, CBC)", "심전도 + 흉부 X-ray", "전문의 상담 및 추적 관리 계획 수립"}
	} else {
		checkupRec.Title = "정기 건강검진 안내"
		checkupRec.Description = "정기적인 건강검진은 질병을 조기에 발견하고 예방하는 가장 좋은 방법입니다."
		checkupRec.Reason = "측정 데이터를 기반으로 전문적인 검진이 권장됩니다"
		checkupRec.Priority = RiskLevelModerate
		checkupRec.ActionSteps = []string{"연 1회 종합 건강검진 받기", "혈액 검사 (공복혈당, HbA1c, 지질 패널)", "혈압 및 심전도 검사", "눈, 치아 정기 검진"}
	}
	recs = append(recs, checkupRec)

	return recs
}

// prioritizeRecommendations는 건강 데이터에 따라 추천 우선순위를 재정렬합니다.
// 가장 낮은 카테고리 점수 관련 추천이 최상위로 올라갑니다.
func (s *CoachingService) prioritizeRecommendations(recs []*Recommendation, healthData *UserHealthData) {
	// 카테고리→추천 타입 매핑
	categoryToType := map[string]RecommendationType{
		"metabolic":      RecommendationTypeFood,
		"cardiovascular": RecommendationTypeExercise,
		"nutritional":    RecommendationTypeSupplement,
		"fitness":        RecommendationTypeExercise,
	}

	// 가장 낮은 카테고리 찾기
	worstCategory := ""
	worstScore := 100.0
	for cat, sc := range healthData.CategoryScores {
		if sc < worstScore {
			worstScore = sc
			worstCategory = cat
		}
	}

	if worstCategory == "" {
		return
	}

	targetType := categoryToType[worstCategory]

	// 해당 타입의 추천을 맨 앞으로 이동
	for i, rec := range recs {
		if rec.Type == targetType && i > 0 {
			// 위치 교환
			recs[0], recs[i] = recs[i], recs[0]
			break
		}
	}
}

// --- 일일 리포트 헬퍼 ---

func (s *CoachingService) generateDailyRecommendations() []string {
	allRecs := []string{
		"식후 30분 이내에 가벼운 산책을 해 보세요",
		"오늘 물 8잔(2L) 마시기를 목표로 하세요",
		"채소와 과일을 매 끼니에 포함하세요",
		"취침 전 스마트폰 사용을 줄여 보세요",
		"계단을 이용하는 작은 운동부터 시작하세요",
		"하루 10분 명상으로 마음을 안정시키세요",
		"짠 음식을 줄이고 싱겁게 드세요",
		"아침 식사를 거르지 마세요",
	}

	count := s.rng.Intn(3) + 2 // 2~4개
	perm := s.rng.Perm(len(allRecs))
	result := make([]string, 0, count)
	for i := 0; i < count && i < len(perm); i++ {
		result = append(result, allRecs[perm[i]])
	}
	return result
}

func (s *CoachingService) generateDailySummary(score float64, measurements int32) string {
	switch {
	case score >= 85:
		return fmt.Sprintf("오늘 %d건의 측정을 완료했습니다. 건강 점수 %.1f점으로 매우 양호한 상태입니다. 꾸준한 건강 관리가 빛을 발하고 있습니다!", measurements, score)
	case score >= 70:
		return fmt.Sprintf("오늘 %d건의 측정을 완료했습니다. 건강 점수 %.1f점으로 전반적으로 양호합니다. 몇 가지 개선 포인트를 확인해 보세요.", measurements, score)
	case score >= 55:
		return fmt.Sprintf("오늘 %d건의 측정을 완료했습니다. 건강 점수 %.1f점으로 보통 수준입니다. 생활 습관 개선으로 점수를 높일 수 있습니다.", measurements, score)
	default:
		return fmt.Sprintf("오늘 %d건의 측정을 완료했습니다. 건강 점수 %.1f점으로 주의가 필요합니다. 추천 사항을 확인하고 실천해 보세요.", measurements, score)
	}
}

// --- 주간 리포트 헬퍼 ---

func (s *CoachingService) evaluateScoreTrend(reports []*DailyHealthReport) string {
	if len(reports) < 2 {
		return "stable"
	}
	firstHalf := reports[:len(reports)/2]
	secondHalf := reports[len(reports)/2:]

	var avgFirst, avgSecond float64
	for _, r := range firstHalf {
		avgFirst += r.OverallScore
	}
	for _, r := range secondHalf {
		avgSecond += r.OverallScore
	}
	avgFirst /= float64(len(firstHalf))
	avgSecond /= float64(len(secondHalf))

	diff := avgSecond - avgFirst
	if diff > 3 {
		return "improving"
	}
	if diff < -3 {
		return "declining"
	}
	return "stable"
}

func (s *CoachingService) generateWeeklySummary(avgScore float64, totalMeasurements int32, trend string) string {
	trendKo := "안정적"
	switch trend {
	case "improving":
		trendKo = "개선 추세"
	case "declining":
		trendKo = "하락 추세"
	}

	return fmt.Sprintf("이번 주 총 %d건의 측정을 수행했으며, 평균 건강 점수는 %.1f점입니다. 전반적인 건강 추세는 '%s'입니다.",
		totalMeasurements, avgScore, trendKo)
}

func (s *CoachingService) generateKeyInsights(avgScore float64, totalMeasurements int32, goalsActive int32) []string {
	insights := make([]string, 0, 4)

	if avgScore >= 80 {
		insights = append(insights, "이번 주 건강 점수가 우수합니다. 현재 생활 습관을 유지하세요.")
	} else if avgScore >= 60 {
		insights = append(insights, "건강 점수를 높이려면 규칙적인 운동과 식단 관리를 병행하세요.")
	} else {
		insights = append(insights, "건강 점수가 낮습니다. 전문의 상담과 함께 생활 습관을 재점검하세요.")
	}

	if totalMeasurements >= 14 {
		insights = append(insights, fmt.Sprintf("이번 주 %d건의 측정을 완료했습니다. 꾸준한 모니터링에 감사드립니다!", totalMeasurements))
	} else {
		insights = append(insights, "측정 빈도를 높이면 더 정확한 건강 트렌드를 파악할 수 있습니다.")
	}

	if goalsActive > 0 {
		insights = append(insights, fmt.Sprintf("현재 %d개의 활성 목표가 있습니다. 목표를 향해 꾸준히 노력하세요!", goalsActive))
	} else {
		insights = append(insights, "새로운 건강 목표를 설정하여 동기부여를 받아 보세요.")
	}

	insights = append(insights, "다음 주에도 건강한 습관을 유지하세요. 만파식이 함께합니다! 💪")

	return insights
}
