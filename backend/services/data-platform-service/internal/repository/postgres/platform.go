package postgres

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/data-platform-service/internal/service"
)

type PlatformRepository struct {
	pool *pgxpool.Pool
}

func NewPlatformRepository(pool *pgxpool.Pool) *PlatformRepository {
	return &PlatformRepository{pool: pool}
}

func (r *PlatformRepository) GetRegionStats(ctx context.Context, code, period, biomarker string) (*service.RegionStats, error) {
	const q = `SELECT region_code, region_name, period, avg_value, min_value, max_value, sample_count, risk_level
		FROM dp_region_stats WHERE region_code = $1 AND period = $2`
	var s service.RegionStats
	err := r.pool.QueryRow(ctx, q, code, period).Scan(&s.RegionCode, &s.RegionName, &s.Period, &s.AvgValue, &s.MinValue, &s.MaxValue, &s.SampleCount, &s.RiskLevel)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &s, nil
}

func (r *PlatformRepository) ListHotspots(ctx context.Context, biomarker, risk string, limit int32) ([]*service.RegionStats, error) {
	const q = `SELECT region_code, region_name, period, avg_value, min_value, max_value, sample_count, risk_level
		FROM dp_region_stats WHERE risk_level = $1 ORDER BY avg_value DESC LIMIT $2`
	rows, err := r.pool.Query(ctx, q, risk, limit)
	if err != nil { return nil, err }
	defer rows.Close()
	var list []*service.RegionStats
	for rows.Next() {
		var s service.RegionStats
		if err := rows.Scan(&s.RegionCode, &s.RegionName, &s.Period, &s.AvgValue, &s.MinValue, &s.MaxValue, &s.SampleCount, &s.RiskLevel); err != nil { return nil, err }
		list = append(list, &s)
	}
	return list, rows.Err()
}

func (r *PlatformRepository) GetTrend(ctx context.Context, code, biomarker string) (*service.RegionTrend, error) {
	const q = `SELECT region_code, biomarker, trend_direction, points FROM dp_region_trends WHERE region_code = $1 AND biomarker = $2`
	var t service.RegionTrend
	var pointsJSON []byte
	err := r.pool.QueryRow(ctx, q, code, biomarker).Scan(&t.RegionCode, &t.Biomarker, &t.TrendDirection, &pointsJSON)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	if len(pointsJSON) > 0 { json.Unmarshal(pointsJSON, &t.Points) }
	return &t, nil
}

func (r *PlatformRepository) ListDatasets(ctx context.Context, dsType string, limit, offset int32) ([]*service.AnonymizedDataset, int32, error) {
	var total int32
	if dsType != "" {
		r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM dp_datasets WHERE dataset_type=$1`, dsType).Scan(&total)
	} else {
		r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM dp_datasets`).Scan(&total)
	}
	q := `SELECT id, dataset_type, description, record_count, date_range, created_at FROM dp_datasets`
	var args []interface{}
	if dsType != "" {
		q += ` WHERE dataset_type=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
		args = append(args, dsType, limit, offset)
	} else {
		q += ` ORDER BY created_at DESC LIMIT $1 OFFSET $2`
		args = append(args, limit, offset)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil { return nil, 0, err }
	defer rows.Close()
	var list []*service.AnonymizedDataset
	for rows.Next() {
		var d service.AnonymizedDataset
		if err := rows.Scan(&d.ID, &d.DatasetType, &d.Description, &d.RecordCount, &d.DateRange, &d.CreatedAt); err != nil { return nil, 0, err }
		list = append(list, &d)
	}
	return list, total, rows.Err()
}

func (r *PlatformRepository) GetDataset(ctx context.Context, id string) (*service.AnonymizedDataset, error) {
	const q = `SELECT id, dataset_type, description, record_count, date_range, created_at FROM dp_datasets WHERE id = $1`
	var d service.AnonymizedDataset
	err := r.pool.QueryRow(ctx, q, id).Scan(&d.ID, &d.DatasetType, &d.Description, &d.RecordCount, &d.DateRange, &d.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &d, nil
}

func (r *PlatformRepository) SaveDataset(ctx context.Context, ds *service.AnonymizedDataset) error {
	const q = `INSERT INTO dp_datasets (id, dataset_type, description, record_count, date_range, created_at)
		VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT (id) DO NOTHING`
	_, err := r.pool.Exec(ctx, q, ds.ID, ds.DatasetType, ds.Description, ds.RecordCount, ds.DateRange, ds.CreatedAt)
	return err
}

func (r *PlatformRepository) CreateAccessRequest(ctx context.Context, requesterID, datasetID, purpose string) (*service.DataAccessResp, error) {
	return &service.DataAccessResp{RequestID: "req-" + requesterID[:8], Status: "pending"}, nil
}

func (r *PlatformRepository) GetConsent(ctx context.Context, userID string) (*service.ConsentInfo, error) {
	const q = `SELECT user_id, anonymous_research, region_stats_consent, enterprise_data, updated_at
		FROM dp_consents WHERE user_id = $1`
	var c service.ConsentInfo
	err := r.pool.QueryRow(ctx, q, userID).Scan(&c.UserID, &c.AnonymousResearch, &c.RegionStatsConsent, &c.EnterpriseData, &c.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &c, nil
}

func (r *PlatformRepository) UpdateConsent(ctx context.Context, c *service.ConsentInfo) error {
	const q = `INSERT INTO dp_consents (user_id, anonymous_research, region_stats_consent, enterprise_data, updated_at)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (user_id) DO UPDATE SET anonymous_research=EXCLUDED.anonymous_research, region_stats_consent=EXCLUDED.region_stats_consent, enterprise_data=EXCLUDED.enterprise_data, updated_at=EXCLUDED.updated_at`
	_, err := r.pool.Exec(ctx, q, c.UserID, c.AnonymousResearch, c.RegionStatsConsent, c.EnterpriseData, c.UpdatedAt)
	return err
}

func (r *PlatformRepository) ListInsights(ctx context.Context, biomarker string, limit int32) ([]*service.AggregatedInsight, error) {
	const q = `SELECT id, biomarker, insight_type, summary, created_at FROM dp_insights WHERE biomarker = $1 ORDER BY created_at DESC LIMIT $2`
	rows, err := r.pool.Query(ctx, q, biomarker, limit)
	if err != nil { return nil, err }
	defer rows.Close()
	var list []*service.AggregatedInsight
	for rows.Next() {
		var i service.AggregatedInsight
		if err := rows.Scan(&i.ID, &i.Biomarker, &i.InsightType, &i.Summary, &i.CreatedAt); err != nil { return nil, err }
		list = append(list, &i)
	}
	return list, rows.Err()
}

func (r *PlatformRepository) SaveInsight(ctx context.Context, i *service.AggregatedInsight) error {
	const q = `INSERT INTO dp_insights (id, biomarker, insight_type, summary, created_at) VALUES ($1,$2,$3,$4,$5)`
	_, err := r.pool.Exec(ctx, q, i.ID, i.Biomarker, i.InsightType, i.Summary, i.CreatedAt)
	return err
}
