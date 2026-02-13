// Package postgres는 ai-inference-service의 PostgreSQL 저장소 구현입니다.
package postgres

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/ai-inference-service/internal/service"
)

// ============================================================================
// AnalysisRepository — PostgreSQL 기반
// ============================================================================

// AnalysisRepository는 PostgreSQL 기반 AnalysisRepository 구현입니다.
type AnalysisRepository struct {
	pool *pgxpool.Pool
}

// NewAnalysisRepository는 PostgreSQL AnalysisRepository를 생성합니다.
func NewAnalysisRepository(pool *pgxpool.Pool) *AnalysisRepository {
	return &AnalysisRepository{pool: pool}
}

// Save는 분석 결과를 저장합니다 (트랜잭션: analysis_results + biomarker_results + anomaly_flags).
func (r *AnalysisRepository) Save(ctx context.Context, result *service.AnalysisResult) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	const qAnalysis = `INSERT INTO analysis_results (id, user_id, measurement_id, overall_health_score, summary, analyzed_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	if _, err := tx.Exec(ctx, qAnalysis,
		result.AnalysisID, result.UserID, result.MeasurementID,
		result.OverallHealthScore, result.Summary, result.AnalyzedAt,
	); err != nil {
		return err
	}

	const qBiomarker = `INSERT INTO biomarker_results (analysis_id, biomarker_name, value, unit, classification, confidence, risk_level, reference_range)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	for _, bm := range result.Biomarkers {
		if _, err := tx.Exec(ctx, qBiomarker,
			result.AnalysisID, bm.BiomarkerName, bm.Value, bm.Unit,
			bm.Classification, bm.Confidence, riskLevelToString(bm.RiskLevel), bm.ReferenceRange,
		); err != nil {
			return err
		}
	}

	const qAnomaly = `INSERT INTO anomaly_flags (analysis_id, metric_name, value, expected_min, expected_max, anomaly_score, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	for _, a := range result.Anomalies {
		if _, err := tx.Exec(ctx, qAnomaly,
			result.AnalysisID, a.MetricName, a.Value,
			a.ExpectedMin, a.ExpectedMax, a.AnomalyScore, a.Description,
		); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// FindByID는 분석 결과를 ID로 조회합니다.
func (r *AnalysisRepository) FindByID(ctx context.Context, id string) (*service.AnalysisResult, error) {
	const q = `SELECT id, user_id, measurement_id, overall_health_score, COALESCE(summary, ''), analyzed_at
		FROM analysis_results WHERE id = $1`

	var ar service.AnalysisResult
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&ar.AnalysisID, &ar.UserID, &ar.MeasurementID,
		&ar.OverallHealthScore, &ar.Summary, &ar.AnalyzedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	biomarkers, err := r.loadBiomarkers(ctx, id)
	if err != nil {
		return nil, err
	}
	ar.Biomarkers = biomarkers

	anomalies, err := r.loadAnomalies(ctx, id)
	if err != nil {
		return nil, err
	}
	ar.Anomalies = anomalies

	return &ar, nil
}

// FindByUserID는 사용자의 분석 결과 목록을 조회합니다.
func (r *AnalysisRepository) FindByUserID(ctx context.Context, userID string, limit int) ([]*service.AnalysisResult, error) {
	const q = `SELECT id, user_id, measurement_id, overall_health_score, COALESCE(summary, ''), analyzed_at
		FROM analysis_results WHERE user_id = $1 ORDER BY analyzed_at DESC LIMIT $2`

	if limit <= 0 {
		limit = 100
	}

	rows, err := r.pool.Query(ctx, q, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*service.AnalysisResult
	for rows.Next() {
		var ar service.AnalysisResult
		if err := rows.Scan(
			&ar.AnalysisID, &ar.UserID, &ar.MeasurementID,
			&ar.OverallHealthScore, &ar.Summary, &ar.AnalyzedAt,
		); err != nil {
			return nil, err
		}
		results = append(results, &ar)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, ar := range results {
		bms, err := r.loadBiomarkers(ctx, ar.AnalysisID)
		if err != nil {
			return nil, err
		}
		ar.Biomarkers = bms

		ans, err := r.loadAnomalies(ctx, ar.AnalysisID)
		if err != nil {
			return nil, err
		}
		ar.Anomalies = ans
	}

	return results, nil
}

func (r *AnalysisRepository) loadBiomarkers(ctx context.Context, analysisID string) ([]service.BiomarkerResult, error) {
	const q = `SELECT biomarker_name, value, unit, classification, confidence, risk_level, COALESCE(reference_range, '')
		FROM biomarker_results WHERE analysis_id = $1`

	rows, err := r.pool.Query(ctx, q, analysisID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []service.BiomarkerResult
	for rows.Next() {
		var bm service.BiomarkerResult
		var rl string
		if err := rows.Scan(
			&bm.BiomarkerName, &bm.Value, &bm.Unit,
			&bm.Classification, &bm.Confidence, &rl, &bm.ReferenceRange,
		); err != nil {
			return nil, err
		}
		bm.RiskLevel = riskLevelFromString(rl)
		results = append(results, bm)
	}
	return results, rows.Err()
}

func (r *AnalysisRepository) loadAnomalies(ctx context.Context, analysisID string) ([]service.AnomalyFlag, error) {
	const q = `SELECT metric_name, value, expected_min, expected_max, anomaly_score, COALESCE(description, '')
		FROM anomaly_flags WHERE analysis_id = $1`

	rows, err := r.pool.Query(ctx, q, analysisID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []service.AnomalyFlag
	for rows.Next() {
		var a service.AnomalyFlag
		if err := rows.Scan(
			&a.MetricName, &a.Value, &a.ExpectedMin,
			&a.ExpectedMax, &a.AnomalyScore, &a.Description,
		); err != nil {
			return nil, err
		}
		results = append(results, a)
	}
	return results, rows.Err()
}

// ============================================================================
// HealthScoreRepository — PostgreSQL 기반
// ============================================================================

// HealthScoreRepository는 PostgreSQL 기반 HealthScoreRepository 구현입니다.
type HealthScoreRepository struct {
	pool *pgxpool.Pool
}

// NewHealthScoreRepository는 PostgreSQL HealthScoreRepository를 생성합니다.
func NewHealthScoreRepository(pool *pgxpool.Pool) *HealthScoreRepository {
	return &HealthScoreRepository{pool: pool}
}

// Save는 건강 점수를 저장합니다.
func (r *HealthScoreRepository) Save(ctx context.Context, score *service.HealthScore) error {
	catJSON, err := json.Marshal(score.CategoryScores)
	if err != nil {
		catJSON = []byte("{}")
	}

	const q = `INSERT INTO health_scores (user_id, overall_score, category_scores, trend, recommendation, calculated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = r.pool.Exec(ctx, q,
		score.UserID, score.OverallScore, catJSON,
		score.Trend, score.Recommendation, score.CalculatedAt,
	)
	return err
}

// FindLatestByUserID는 사용자의 최신 건강 점수를 조회합니다.
func (r *HealthScoreRepository) FindLatestByUserID(ctx context.Context, userID string) (*service.HealthScore, error) {
	const q = `SELECT user_id, overall_score, category_scores, trend, COALESCE(recommendation, ''), calculated_at
		FROM health_scores WHERE user_id = $1 ORDER BY calculated_at DESC LIMIT 1`

	var hs service.HealthScore
	var catJSON []byte
	err := r.pool.QueryRow(ctx, q, userID).Scan(
		&hs.UserID, &hs.OverallScore, &catJSON,
		&hs.Trend, &hs.Recommendation, &hs.CalculatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	hs.CategoryScores = make(map[string]float64)
	_ = json.Unmarshal(catJSON, &hs.CategoryScores)
	return &hs, nil
}

// ============================================================================
// ENUM 변환 헬퍼
// ============================================================================

func riskLevelToString(rl service.RiskLevel) string {
	switch rl {
	case service.RiskLow:
		return "LOW"
	case service.RiskModerate:
		return "MODERATE"
	case service.RiskHigh:
		return "HIGH"
	case service.RiskCritical:
		return "CRITICAL"
	default:
		return "LOW"
	}
}

func riskLevelFromString(s string) service.RiskLevel {
	switch s {
	case "LOW":
		return service.RiskLow
	case "MODERATE":
		return service.RiskModerate
	case "HIGH":
		return service.RiskHigh
	case "CRITICAL":
		return service.RiskCritical
	default:
		return service.RiskLow
	}
}
