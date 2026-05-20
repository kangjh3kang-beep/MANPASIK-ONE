package postgres

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/nlp-service/internal/service"
)

type NLPRepository struct {
	pool *pgxpool.Pool
}

func NewNLPRepository(pool *pgxpool.Pool) *NLPRepository {
	return &NLPRepository{pool: pool}
}

func (r *NLPRepository) SaveQuery(ctx context.Context, query *service.HealthQuery) error {
	entJSON, _ := json.Marshal(query.Entities)
	const q = `INSERT INTO nlp_health_queries (id, user_id, raw_text, intent, entities, confidence, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(ctx, q, query.ID, query.UserID, query.RawText, query.Intent, entJSON, query.Confidence, query.CreatedAt)
	return err
}

func (r *NLPRepository) GetQuery(ctx context.Context, queryID string) (*service.HealthQuery, error) {
	const q = `SELECT id, user_id, raw_text, intent, entities, confidence, created_at
		FROM nlp_health_queries WHERE id = $1`
	var hq service.HealthQuery
	var entJSON []byte
	err := r.pool.QueryRow(ctx, q, queryID).Scan(&hq.ID, &hq.UserID, &hq.RawText, &hq.Intent, &entJSON, &hq.Confidence, &hq.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	if len(entJSON) > 0 { json.Unmarshal(entJSON, &hq.Entities) }
	return &hq, nil
}

func (r *NLPRepository) SaveExtraction(ctx context.Context, extraction *service.SymptomExtraction) error {
	sympJSON, _ := json.Marshal(extraction.Symptoms)
	const q = `INSERT INTO nlp_symptom_extractions (id, text, symptoms, processed_at)
		VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, q, extraction.ID, extraction.Text, sympJSON, extraction.ProcessedAt)
	return err
}

func (r *NLPRepository) GetSuggestions(ctx context.Context, queryID string) ([]service.Suggestion, error) {
	const q = `SELECT id, query_id, text, category, priority FROM nlp_suggestions WHERE query_id = $1 ORDER BY priority DESC`
	rows, err := r.pool.Query(ctx, q, queryID)
	if err != nil { return nil, err }
	defer rows.Close()
	var list []service.Suggestion
	for rows.Next() {
		var s service.Suggestion
		if err := rows.Scan(&s.ID, &s.QueryID, &s.Text, &s.Category, &s.Priority); err != nil { return nil, err }
		list = append(list, s)
	}
	return list, rows.Err()
}

func (r *NLPRepository) SaveSuggestions(ctx context.Context, queryID string, suggestions []service.Suggestion) error {
	for _, s := range suggestions {
		const q = `INSERT INTO nlp_suggestions (id, query_id, text, category, priority) VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (id) DO UPDATE SET text = EXCLUDED.text, category = EXCLUDED.category, priority = EXCLUDED.priority`
		if _, err := r.pool.Exec(ctx, q, s.ID, queryID, s.Text, s.Category, s.Priority); err != nil {
			return err
		}
	}
	return nil
}
