package handler

import (
	"context"

	"github.com/manpasik/backend/services/ai-inference-service/internal/service"
	apperr "github.com/manpasik/backend/shared/errors"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type InferenceHandler struct {
	v1.UnimplementedAiInferenceServiceServer
	svc *service.InferenceService
}

func NewInferenceHandler(svc *service.InferenceService) *InferenceHandler {
	return &InferenceHandler{svc: svc}
}

// AnalyzeMeasurement implements v1.AiInferenceServiceServer.
func (h *InferenceHandler) AnalyzeMeasurement(ctx context.Context, req *v1.AnalyzeMeasurementRequest) (*v1.AnalysisResult, error) {
	models := protoModelsToService(req.Models)
	result, err := h.svc.AnalyzeMeasurement(ctx, req.UserId, req.MeasurementId, models)
	if err != nil {
		return nil, toGRPC(err)
	}
	return analysisResultToProto(result), nil
}

// GetHealthScore implements v1.AiInferenceServiceServer.
func (h *InferenceHandler) GetHealthScore(ctx context.Context, req *v1.GetHealthScoreRequest) (*v1.HealthScoreResponse, error) {
	score, err := h.svc.GetHealthScore(ctx, req.UserId, int(req.Days))
	if err != nil {
		return nil, toGRPC(err)
	}
	return healthScoreToProto(score), nil
}

// PredictTrend implements v1.AiInferenceServiceServer.
func (h *InferenceHandler) PredictTrend(ctx context.Context, req *v1.PredictTrendRequest) (*v1.TrendPrediction, error) {
	pred, err := h.svc.PredictTrend(ctx, req.UserId, req.MetricName, int(req.HistoryDays), int(req.PredictionDays))
	if err != nil {
		return nil, toGRPC(err)
	}
	return trendPredictionToProto(pred), nil
}

// GetModelInfo implements v1.AiInferenceServiceServer.
func (h *InferenceHandler) GetModelInfo(ctx context.Context, req *v1.GetModelInfoRequest) (*v1.ModelInfo, error) {
	mt := protoModelTypeToService(req.ModelType)
	info, err := h.svc.GetModelInfo(ctx, mt)
	if err != nil {
		return nil, toGRPC(err)
	}
	return modelInfoToProto(info), nil
}

// ListModels implements v1.AiInferenceServiceServer.
func (h *InferenceHandler) ListModels(ctx context.Context, _ *v1.ListModelsRequest) (*v1.ListModelsResponse, error) {
	models, err := h.svc.ListModels(ctx)
	if err != nil {
		return nil, toGRPC(err)
	}
	protoModels := make([]*v1.ModelInfo, len(models))
	for i, m := range models {
		protoModels[i] = modelInfoToProto(m)
	}
	return &v1.ListModelsResponse{Models: protoModels}, nil
}

// GenerateHealthInsight는 LLM 기반 건강 인사이트를 생성합니다.
// Proto에 별도 RPC가 정의되지 않았으므로, AnalyzeMeasurement의 Summary 필드에서
// LLM 향상이 자동으로 적용됩니다. 이 메서드는 내부/테스트용으로 직접 호출할 수 있습니다.
func (h *InferenceHandler) GenerateHealthInsight(ctx context.Context, userID string, measurements []service.MeasurementData) (string, error) {
	insight, err := h.svc.GenerateHealthInsight(ctx, userID, measurements)
	if err != nil {
		return "", toGRPC(err)
	}
	return insight, nil
}

// LLMEnabled는 LLM 클라이언트가 설정되어 있는지 반환합니다.
func (h *InferenceHandler) LLMEnabled() bool {
	return h.svc.LLMEnabled()
}

// ============================================================================
// Converters
// ============================================================================

func protoModelsToService(models []v1.AiModelType) []service.AiModelType {
	if len(models) == 0 {
		return nil
	}
	result := make([]service.AiModelType, len(models))
	for i, m := range models {
		result[i] = protoModelTypeToService(m)
	}
	return result
}

func protoModelTypeToService(mt v1.AiModelType) service.AiModelType {
	switch mt {
	case v1.AiModelType_AI_MODEL_TYPE_BIOMARKER_CLASSIFIER:
		return service.ModelBiomarkerClassifier
	case v1.AiModelType_AI_MODEL_TYPE_ANOMALY_DETECTOR:
		return service.ModelAnomalyDetector
	case v1.AiModelType_AI_MODEL_TYPE_TREND_PREDICTOR:
		return service.ModelTrendPredictor
	case v1.AiModelType_AI_MODEL_TYPE_HEALTH_SCORER:
		return service.ModelHealthScorer
	case v1.AiModelType_AI_MODEL_TYPE_FOOD_CALORIE_ESTIMATOR:
		return service.ModelFoodCalorieEstimator
	default:
		return service.ModelBiomarkerClassifier
	}
}

func serviceModelTypeToProto(mt service.AiModelType) v1.AiModelType {
	switch mt {
	case service.ModelBiomarkerClassifier:
		return v1.AiModelType_AI_MODEL_TYPE_BIOMARKER_CLASSIFIER
	case service.ModelAnomalyDetector:
		return v1.AiModelType_AI_MODEL_TYPE_ANOMALY_DETECTOR
	case service.ModelTrendPredictor:
		return v1.AiModelType_AI_MODEL_TYPE_TREND_PREDICTOR
	case service.ModelHealthScorer:
		return v1.AiModelType_AI_MODEL_TYPE_HEALTH_SCORER
	case service.ModelFoodCalorieEstimator:
		return v1.AiModelType_AI_MODEL_TYPE_FOOD_CALORIE_ESTIMATOR
	default:
		return v1.AiModelType_AI_MODEL_TYPE_UNSPECIFIED
	}
}

func serviceRiskLevelToProto(rl service.RiskLevel) v1.RiskLevel {
	switch rl {
	case service.RiskLow:
		return v1.RiskLevel_RISK_LEVEL_LOW
	case service.RiskModerate:
		return v1.RiskLevel_RISK_LEVEL_MODERATE
	case service.RiskHigh:
		return v1.RiskLevel_RISK_LEVEL_HIGH
	case service.RiskCritical:
		return v1.RiskLevel_RISK_LEVEL_CRITICAL
	default:
		return v1.RiskLevel_RISK_LEVEL_UNSPECIFIED
	}
}

func analysisResultToProto(r *service.AnalysisResult) *v1.AnalysisResult {
	biomarkers := make([]*v1.BiomarkerResult, len(r.Biomarkers))
	for i, b := range r.Biomarkers {
		biomarkers[i] = &v1.BiomarkerResult{
			BiomarkerName:  b.BiomarkerName,
			Value:          b.Value,
			Unit:           b.Unit,
			Classification: b.Classification,
			Confidence:     b.Confidence,
			RiskLevel:      serviceRiskLevelToProto(b.RiskLevel),
			ReferenceRange: b.ReferenceRange,
		}
	}
	anomalies := make([]*v1.AnomalyFlag, len(r.Anomalies))
	for i, a := range r.Anomalies {
		anomalies[i] = &v1.AnomalyFlag{
			MetricName:   a.MetricName,
			Value:        a.Value,
			ExpectedMin:  a.ExpectedMin,
			ExpectedMax:  a.ExpectedMax,
			AnomalyScore: a.AnomalyScore,
			Description:  a.Description,
		}
	}
	return &v1.AnalysisResult{
		AnalysisId:         r.AnalysisID,
		UserId:             r.UserID,
		MeasurementId:      r.MeasurementID,
		Biomarkers:         biomarkers,
		Anomalies:          anomalies,
		OverallHealthScore: r.OverallHealthScore,
		Summary:            r.Summary,
		AnalyzedAt:         timestamppb.New(r.AnalyzedAt),
	}
}

func healthScoreToProto(s *service.HealthScore) *v1.HealthScoreResponse {
	return &v1.HealthScoreResponse{
		UserId:         s.UserID,
		OverallScore:   s.OverallScore,
		CategoryScores: s.CategoryScores,
		Trend:          s.Trend,
		Recommendation: s.Recommendation,
		CalculatedAt:   timestamppb.New(s.CalculatedAt),
	}
}

func trendPredictionToProto(p *service.TrendPrediction) *v1.TrendPrediction {
	historical := make([]*v1.TrendDataPoint, len(p.Historical))
	for i, dp := range p.Historical {
		historical[i] = &v1.TrendDataPoint{
			Timestamp:  timestamppb.New(dp.Timestamp),
			Value:      dp.Value,
			LowerBound: dp.LowerBound,
			UpperBound: dp.UpperBound,
		}
	}
	predicted := make([]*v1.TrendDataPoint, len(p.Predicted))
	for i, dp := range p.Predicted {
		predicted[i] = &v1.TrendDataPoint{
			Timestamp:  timestamppb.New(dp.Timestamp),
			Value:      dp.Value,
			LowerBound: dp.LowerBound,
			UpperBound: dp.UpperBound,
		}
	}
	return &v1.TrendPrediction{
		UserId:     p.UserID,
		MetricName: p.MetricName,
		Historical: historical,
		Predicted:  predicted,
		Confidence: p.Confidence,
		Direction:  p.Direction,
		Insight:    p.Insight,
	}
}

func modelInfoToProto(m *service.ModelInfo) *v1.ModelInfo {
	return &v1.ModelInfo{
		ModelType:   serviceModelTypeToProto(m.ModelType),
		Name:        m.Name,
		Version:     m.Version,
		Description: m.Description,
		Accuracy:    m.Accuracy,
		LastTrained: timestamppb.New(m.LastTrained),
		Status:      string(m.Status),
	}
}

func toGRPC(err error) error {
	if err == nil {
		return nil
	}
	if ae, ok := err.(*apperr.AppError); ok {
		return ae.ToGRPC()
	}
	if s, ok := status.FromError(err); ok {
		return s.Err()
	}
	return status.Error(codes.Internal, "내부 오류가 발생했습니다")
}
