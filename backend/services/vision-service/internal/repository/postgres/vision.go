package postgres

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/vision-service/internal/service"
)

type FoodAnalysisRepository struct {
	pool *pgxpool.Pool
}

func NewFoodAnalysisRepository(pool *pgxpool.Pool) *FoodAnalysisRepository {
	return &FoodAnalysisRepository{pool: pool}
}

func (r *FoodAnalysisRepository) Save(ctx context.Context, analysis *service.FoodAnalysis) error {
	itemsJSON, _ := json.Marshal(analysis.FoodItems)
	const q = `INSERT INTO food_analyses (id, user_id, image_url, status, total_calorie_kcal, food_items, meal_type, analyzed_at, created_at, error_message)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.pool.Exec(ctx, q,
		analysis.ID, analysis.UserID, analysis.ImageURL, int32(analysis.Status),
		analysis.TotalCalorieKcal, itemsJSON, analysis.MealType,
		analysis.AnalyzedAt, analysis.CreatedAt, analysis.ErrorMessage,
	)
	return err
}

func (r *FoodAnalysisRepository) FindByID(ctx context.Context, id string) (*service.FoodAnalysis, error) {
	const q = `SELECT id, user_id, image_url, status, total_calorie_kcal, food_items, meal_type, analyzed_at, created_at, error_message
		FROM food_analyses WHERE id = $1`
	var a service.FoodAnalysis
	var statusInt int32
	var itemsJSON []byte
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&a.ID, &a.UserID, &a.ImageURL, &statusInt,
		&a.TotalCalorieKcal, &itemsJSON, &a.MealType,
		&a.AnalyzedAt, &a.CreatedAt, &a.ErrorMessage,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	a.Status = service.FoodAnalysisStatus(statusInt)
	if len(itemsJSON) > 0 {
		json.Unmarshal(itemsJSON, &a.FoodItems)
	}
	return &a, nil
}

func (r *FoodAnalysisRepository) FindByUserID(ctx context.Context, userID string, limit, offset int32) ([]*service.FoodAnalysis, int32, error) {
	var total int32
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM food_analyses WHERE user_id = $1`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	const q = `SELECT id, user_id, image_url, status, total_calorie_kcal, food_items, meal_type, analyzed_at, created_at, error_message
		FROM food_analyses WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []*service.FoodAnalysis
	for rows.Next() {
		var a service.FoodAnalysis
		var statusInt int32
		var itemsJSON []byte
		if err := rows.Scan(
			&a.ID, &a.UserID, &a.ImageURL, &statusInt,
			&a.TotalCalorieKcal, &itemsJSON, &a.MealType,
			&a.AnalyzedAt, &a.CreatedAt, &a.ErrorMessage,
		); err != nil {
			return nil, 0, err
		}
		a.Status = service.FoodAnalysisStatus(statusInt)
		if len(itemsJSON) > 0 {
			json.Unmarshal(itemsJSON, &a.FoodItems)
		}
		list = append(list, &a)
	}
	return list, total, rows.Err()
}

func (r *FoodAnalysisRepository) Update(ctx context.Context, analysis *service.FoodAnalysis) error {
	itemsJSON, _ := json.Marshal(analysis.FoodItems)
	const q = `UPDATE food_analyses SET status=$1, total_calorie_kcal=$2, food_items=$3, analyzed_at=$4, error_message=$5
		WHERE id = $6`
	_, err := r.pool.Exec(ctx, q,
		int32(analysis.Status), analysis.TotalCalorieKcal, itemsJSON,
		analysis.AnalyzedAt, analysis.ErrorMessage, analysis.ID,
	)
	return err
}

func (r *FoodAnalysisRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM food_analyses WHERE id = $1`, id)
	return err
}
