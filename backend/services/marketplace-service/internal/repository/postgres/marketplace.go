package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/marketplace-service/internal/service"
)

// ProductRepository is the PostgreSQL implementation.
type ProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

func (r *ProductRepository) FindByPartner(ctx context.Context, partnerID string) ([]*service.PartnerProduct, error) {
	return r.queryProducts(ctx, `SELECT id, partner_id, name, description, price, category, image_url, is_active, created_at, updated_at
		FROM marketplace_products WHERE partner_id = $1 AND is_active = TRUE ORDER BY created_at DESC`, partnerID)
}

func (r *ProductRepository) FindByCategory(ctx context.Context, category string) ([]*service.PartnerProduct, error) {
	return r.queryProducts(ctx, `SELECT id, partner_id, name, description, price, category, image_url, is_active, created_at, updated_at
		FROM marketplace_products WHERE category = $1 AND is_active = TRUE ORDER BY created_at DESC`, category)
}

func (r *ProductRepository) FindByPartnerAndCategory(ctx context.Context, partnerID, category string) ([]*service.PartnerProduct, error) {
	return r.queryProducts(ctx, `SELECT id, partner_id, name, description, price, category, image_url, is_active, created_at, updated_at
		FROM marketplace_products WHERE partner_id = $1 AND category = $2 AND is_active = TRUE ORDER BY created_at DESC`, partnerID, category)
}

func (r *ProductRepository) FindAll(ctx context.Context) ([]*service.PartnerProduct, error) {
	return r.queryProducts(ctx, `SELECT id, partner_id, name, description, price, category, image_url, is_active, created_at, updated_at
		FROM marketplace_products WHERE is_active = TRUE ORDER BY created_at DESC LIMIT 100`)
}

func (r *ProductRepository) FindByID(ctx context.Context, id string) (*service.PartnerProduct, error) {
	const q = `SELECT id, partner_id, name, description, price, category, image_url, is_active, created_at, updated_at
		FROM marketplace_products WHERE id = $1`
	var p service.PartnerProduct
	err := r.pool.QueryRow(ctx, q, id).Scan(&p.ID, &p.PartnerID, &p.Name, &p.Description, &p.Price, &p.Category, &p.ImageURL, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Save(ctx context.Context, product *service.PartnerProduct) error {
	const q = `INSERT INTO marketplace_products (id, partner_id, name, description, price, category, image_url, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.pool.Exec(ctx, q, product.ID, product.PartnerID, product.Name, product.Description, product.Price, product.Category, product.ImageURL, product.IsActive, product.CreatedAt, product.UpdatedAt)
	return err
}

func (r *ProductRepository) Update(ctx context.Context, product *service.PartnerProduct) error {
	const q = `UPDATE marketplace_products SET name=$1, description=$2, price=$3, category=$4, image_url=$5, is_active=$6, updated_at=$7 WHERE id = $8`
	_, err := r.pool.Exec(ctx, q, product.Name, product.Description, product.Price, product.Category, product.ImageURL, product.IsActive, product.UpdatedAt, product.ID)
	return err
}

func (r *ProductRepository) queryProducts(ctx context.Context, q string, args ...interface{}) ([]*service.PartnerProduct, error) {
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil { return nil, err }
	defer rows.Close()
	var list []*service.PartnerProduct
	for rows.Next() {
		var p service.PartnerProduct
		if err := rows.Scan(&p.ID, &p.PartnerID, &p.Name, &p.Description, &p.Price, &p.Category, &p.ImageURL, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &p)
	}
	return list, rows.Err()
}

// PartnerRepository is the PostgreSQL implementation.
type PartnerRepository struct {
	pool *pgxpool.Pool
}

func NewPartnerRepository(pool *pgxpool.Pool) *PartnerRepository {
	return &PartnerRepository{pool: pool}
}

func (r *PartnerRepository) Save(ctx context.Context, partner *service.Partner) error {
	const q = `INSERT INTO marketplace_partners (id, name, description, contact_email, status, created_at)
		VALUES ($1, $2, $3, $4, $5::partner_status, $6)`
	_, err := r.pool.Exec(ctx, q, partner.ID, partner.Name, partner.Description, partner.ContactEmail, partner.Status, partner.CreatedAt)
	return err
}

func (r *PartnerRepository) FindByID(ctx context.Context, id string) (*service.Partner, error) {
	const q = `SELECT id, name, description, contact_email, status::text, created_at FROM marketplace_partners WHERE id = $1`
	var p service.Partner
	err := r.pool.QueryRow(ctx, q, id).Scan(&p.ID, &p.Name, &p.Description, &p.ContactEmail, &p.Status, &p.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &p, nil
}

func (r *PartnerRepository) FindAll(ctx context.Context) ([]*service.Partner, error) {
	const q = `SELECT id, name, description, contact_email, status::text, created_at FROM marketplace_partners ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, q)
	if err != nil { return nil, err }
	defer rows.Close()
	var list []*service.Partner
	for rows.Next() {
		var p service.Partner
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.ContactEmail, &p.Status, &p.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, &p)
	}
	return list, rows.Err()
}

// StatsRepository is the PostgreSQL implementation.
type StatsRepository struct {
	pool *pgxpool.Pool
}

func NewStatsRepository(pool *pgxpool.Pool) *StatsRepository {
	return &StatsRepository{pool: pool}
}

func (r *StatsRepository) FindByPartnerID(ctx context.Context, partnerID string) (*service.PartnerStats, error) {
	const q = `SELECT partner_id, total_products, total_orders, revenue FROM marketplace_stats WHERE partner_id = $1`
	var s service.PartnerStats
	err := r.pool.QueryRow(ctx, q, partnerID).Scan(&s.PartnerID, &s.TotalProducts, &s.TotalOrders, &s.Revenue)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &s, nil
}

func (r *StatsRepository) Save(ctx context.Context, stats *service.PartnerStats) error {
	const q = `INSERT INTO marketplace_stats (partner_id, total_products, total_orders, revenue)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (partner_id) DO UPDATE SET total_products=EXCLUDED.total_products, total_orders=EXCLUDED.total_orders, revenue=EXCLUDED.revenue`
	_, err := r.pool.Exec(ctx, q, stats.PartnerID, stats.TotalProducts, stats.TotalOrders, stats.Revenue)
	return err
}
