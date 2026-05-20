package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/cartridge-store-service/internal/service"
)

type StoreRepository struct {
	pool *pgxpool.Pool
}

func NewStoreRepository(pool *pgxpool.Pool) *StoreRepository {
	return &StoreRepository{pool: pool}
}

func (r *StoreRepository) CreateDeveloper(ctx context.Context, d *service.Developer) error {
	const q = `INSERT INTO cs_developers (id, user_id, company_name, tier, status, created_at) VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.pool.Exec(ctx, q, d.ID, d.UserID, d.CompanyName, d.Tier, d.Status, d.CreatedAt)
	return err
}

func (r *StoreRepository) GetDeveloper(ctx context.Context, userID string) (*service.Developer, error) {
	const q = `SELECT id, user_id, company_name, tier, status, created_at FROM cs_developers WHERE user_id = $1`
	var d service.Developer
	err := r.pool.QueryRow(ctx, q, userID).Scan(&d.ID, &d.UserID, &d.CompanyName, &d.Tier, &d.Status, &d.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &d, nil
}

func (r *StoreRepository) CreateApiKey(ctx context.Context, k *service.ApiKey) error {
	const q = `INSERT INTO cs_api_keys (key_id, developer_id, key, status, created_at) VALUES ($1,$2,$3,$4,$5)`
	_, err := r.pool.Exec(ctx, q, k.KeyID, k.DeveloperID, k.Key, k.Status, k.CreatedAt)
	return err
}

func (r *StoreRepository) CreateItem(ctx context.Context, item *service.StoreItem) error {
	const q = `INSERT INTO cs_store_items (id, developer_id, name, description, category, version, price_krw, status, rating, downloads, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.pool.Exec(ctx, q, item.ID, item.DeveloperID, item.Name, item.Description, item.Category, item.Version, item.PriceKrw, item.Status, item.Rating, item.Downloads, item.CreatedAt)
	return err
}

func (r *StoreRepository) ListItems(ctx context.Context, category string, limit, offset int32) ([]*service.StoreItem, int32, error) {
	var total int32
	if category != "" {
		r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM cs_store_items WHERE category=$1`, category).Scan(&total)
	} else {
		r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM cs_store_items`).Scan(&total)
	}
	q := `SELECT id, developer_id, name, description, category, version, price_krw, status, rating, downloads, created_at FROM cs_store_items`
	var args []interface{}
	if category != "" {
		q += ` WHERE category=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
		args = append(args, category, limit, offset)
	} else {
		q += ` ORDER BY created_at DESC LIMIT $1 OFFSET $2`
		args = append(args, limit, offset)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil { return nil, 0, err }
	defer rows.Close()
	var list []*service.StoreItem
	for rows.Next() {
		var i service.StoreItem
		if err := rows.Scan(&i.ID, &i.DeveloperID, &i.Name, &i.Description, &i.Category, &i.Version, &i.PriceKrw, &i.Status, &i.Rating, &i.Downloads, &i.CreatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, &i)
	}
	return list, total, rows.Err()
}

func (r *StoreRepository) SearchItems(ctx context.Context, query string, limit int32) ([]*service.StoreItem, error) {
	const q = `SELECT id, developer_id, name, description, category, version, price_krw, status, rating, downloads, created_at
		FROM cs_store_items WHERE name ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%' LIMIT $2`
	rows, err := r.pool.Query(ctx, q, query, limit)
	if err != nil { return nil, err }
	defer rows.Close()
	var list []*service.StoreItem
	for rows.Next() {
		var i service.StoreItem
		if err := rows.Scan(&i.ID, &i.DeveloperID, &i.Name, &i.Description, &i.Category, &i.Version, &i.PriceKrw, &i.Status, &i.Rating, &i.Downloads, &i.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, &i)
	}
	return list, rows.Err()
}

func (r *StoreRepository) GetItem(ctx context.Context, id string) (*service.StoreItem, error) {
	const q = `SELECT id, developer_id, name, description, category, version, price_krw, status, rating, downloads, created_at
		FROM cs_store_items WHERE id = $1`
	var i service.StoreItem
	err := r.pool.QueryRow(ctx, q, id).Scan(&i.ID, &i.DeveloperID, &i.Name, &i.Description, &i.Category, &i.Version, &i.PriceKrw, &i.Status, &i.Rating, &i.Downloads, &i.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &i, nil
}

func (r *StoreRepository) CreatePurchase(ctx context.Context, p *service.Purchase) error {
	const q = `INSERT INTO cs_purchases (id, user_id, item_id, price_krw, status, purchased_at) VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.pool.Exec(ctx, q, p.ID, p.UserID, p.ItemID, p.PriceKrw, p.Status, p.PurchasedAt)
	return err
}

func (r *StoreRepository) ListPurchases(ctx context.Context, userID string, limit int32) ([]*service.Purchase, error) {
	const q = `SELECT id, user_id, item_id, price_krw, status, purchased_at FROM cs_purchases WHERE user_id = $1 ORDER BY purchased_at DESC LIMIT $2`
	rows, err := r.pool.Query(ctx, q, userID, limit)
	if err != nil { return nil, err }
	defer rows.Close()
	var list []*service.Purchase
	for rows.Next() {
		var p service.Purchase
		if err := rows.Scan(&p.ID, &p.UserID, &p.ItemID, &p.PriceKrw, &p.Status, &p.PurchasedAt); err != nil { return nil, err }
		list = append(list, &p)
	}
	return list, rows.Err()
}

func (r *StoreRepository) CreateReview(ctx context.Context, rv *service.ReviewStatus) error {
	const q = `INSERT INTO cs_reviews (submission_id, cartridge_id, status, reviewer_comment, submitted_at) VALUES ($1,$2,$3,$4,$5)`
	_, err := r.pool.Exec(ctx, q, rv.SubmissionID, rv.CartridgeID, rv.Status, rv.ReviewerComment, rv.SubmittedAt)
	return err
}

func (r *StoreRepository) GetReview(ctx context.Context, submissionID string) (*service.ReviewStatus, error) {
	const q = `SELECT submission_id, cartridge_id, status, COALESCE(reviewer_comment,''), submitted_at FROM cs_reviews WHERE submission_id = $1`
	var rv service.ReviewStatus
	err := r.pool.QueryRow(ctx, q, submissionID).Scan(&rv.SubmissionID, &rv.CartridgeID, &rv.Status, &rv.ReviewerComment, &rv.SubmittedAt)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &rv, nil
}

func (r *StoreRepository) UpdateReview(ctx context.Context, rv *service.ReviewStatus) error {
	const q = `UPDATE cs_reviews SET status=$1, reviewer_comment=$2 WHERE submission_id = $3`
	_, err := r.pool.Exec(ctx, q, rv.Status, rv.ReviewerComment, rv.SubmissionID)
	return err
}

func (r *StoreRepository) ListSubmissions(ctx context.Context, devID string, limit int32) ([]*service.ReviewStatus, error) {
	const q = `SELECT submission_id, cartridge_id, status, COALESCE(reviewer_comment,''), submitted_at FROM cs_reviews
		WHERE cartridge_id IN (SELECT id FROM cs_store_items WHERE developer_id = $1) ORDER BY submitted_at DESC LIMIT $2`
	rows, err := r.pool.Query(ctx, q, devID, limit)
	if err != nil { return nil, err }
	defer rows.Close()
	var list []*service.ReviewStatus
	for rows.Next() {
		var rv service.ReviewStatus
		if err := rows.Scan(&rv.SubmissionID, &rv.CartridgeID, &rv.Status, &rv.ReviewerComment, &rv.SubmittedAt); err != nil { return nil, err }
		list = append(list, &rv)
	}
	return list, rows.Err()
}
