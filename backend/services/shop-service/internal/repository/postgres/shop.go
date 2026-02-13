// Package postgres는 shop-service의 PostgreSQL 저장소 구현입니다.
package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/shop-service/internal/service"
)

// ============================================================================
// ProductRepository
// ============================================================================

// ProductRepository는 PostgreSQL 기반 ProductRepository 구현입니다.
type ProductRepository struct {
	pool *pgxpool.Pool
}

// NewProductRepository는 ProductRepository를 생성합니다.
func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

// List는 상품 목록을 조회합니다. category가 0이면 전체 조회.
func (r *ProductRepository) List(ctx context.Context, category service.ProductCategory, limit, offset int32) ([]*service.Product, int32, error) {
	var total int32

	// 전체 수 조회
	if category == service.CategoryUnknown {
		err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM products WHERE is_active = TRUE`).Scan(&total)
		if err != nil {
			return nil, 0, err
		}
	} else {
		err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM products WHERE is_active = TRUE AND category = $1`, int32(category)).Scan(&total)
		if err != nil {
			return nil, 0, err
		}
	}

	// 목록 조회
	var q string
	var rows pgx.Rows
	var err error

	if category == service.CategoryUnknown {
		q = `SELECT id, name, description, category, price_krw, stock, COALESCE(image_url, ''), is_active, created_at
			FROM products WHERE is_active = TRUE ORDER BY created_at DESC LIMIT $1 OFFSET $2`
		rows, err = r.pool.Query(ctx, q, limit, offset)
	} else {
		q = `SELECT id, name, description, category, price_krw, stock, COALESCE(image_url, ''), is_active, created_at
			FROM products WHERE is_active = TRUE AND category = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
		rows, err = r.pool.Query(ctx, q, int32(category), limit, offset)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []*service.Product
	for rows.Next() {
		var p service.Product
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.Category,
			&p.PriceKRW, &p.Stock, &p.ImageURL, &p.IsActive, &p.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		products = append(products, &p)
	}
	return products, total, rows.Err()
}

// GetByID는 상품 ID로 조회합니다.
func (r *ProductRepository) GetByID(ctx context.Context, id string) (*service.Product, error) {
	const q = `SELECT id, name, description, category, price_krw, stock, COALESCE(image_url, ''), is_active, created_at
		FROM products WHERE id = $1`
	var p service.Product
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.Category,
		&p.PriceKRW, &p.Stock, &p.ImageURL, &p.IsActive, &p.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// Create는 새 상품을 생성합니다.
func (r *ProductRepository) Create(ctx context.Context, product *service.Product) error {
	const q = `INSERT INTO products (id, name, description, category, price_krw, stock, image_url, is_active, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := r.pool.Exec(ctx, q,
		product.ID, product.Name, product.Description, int32(product.Category),
		product.PriceKRW, product.Stock, product.ImageURL, product.IsActive, product.CreatedAt,
	)
	return err
}

// UpdateStock는 상품 재고를 변경합니다 (delta: 양수=증가, 음수=감소).
func (r *ProductRepository) UpdateStock(ctx context.Context, id string, delta int32) error {
	const q = `UPDATE products SET stock = stock + $1, updated_at = NOW() WHERE id = $2`
	_, err := r.pool.Exec(ctx, q, delta, id)
	return err
}

// ============================================================================
// CartRepository
// ============================================================================

// CartRepository는 PostgreSQL 기반 CartRepository 구현입니다.
type CartRepository struct {
	pool *pgxpool.Pool
}

// NewCartRepository는 CartRepository를 생성합니다.
func NewCartRepository(pool *pgxpool.Pool) *CartRepository {
	return &CartRepository{pool: pool}
}

// GetByUserID는 사용자의 장바구니를 조회합니다.
func (r *CartRepository) GetByUserID(ctx context.Context, userID string) ([]*service.CartItem, error) {
	const q = `SELECT ci.id, ci.user_id, ci.product_id, p.name, ci.quantity, p.price_krw,
		(ci.quantity * p.price_krw) AS total_price_krw
		FROM cart_items ci JOIN products p ON ci.product_id = p.id
		WHERE ci.user_id = $1 ORDER BY ci.created_at DESC`
	rows, err := r.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*service.CartItem
	for rows.Next() {
		var item service.CartItem
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.ProductID, &item.ProductName,
			&item.Quantity, &item.UnitPriceKRW, &item.TotalPriceKRW,
		); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, rows.Err()
}

// AddItem는 장바구니에 항목을 추가합니다.
func (r *CartRepository) AddItem(ctx context.Context, item *service.CartItem) error {
	const q = `INSERT INTO cart_items (id, user_id, product_id, quantity, created_at)
		VALUES ($1,$2,$3,$4,NOW())`
	_, err := r.pool.Exec(ctx, q,
		item.ID, item.UserID, item.ProductID, item.Quantity,
	)
	return err
}

// RemoveItem는 장바구니에서 항목을 제거합니다.
func (r *CartRepository) RemoveItem(ctx context.Context, userID, cartItemID string) error {
	const q = `DELETE FROM cart_items WHERE id = $1 AND user_id = $2`
	_, err := r.pool.Exec(ctx, q, cartItemID, userID)
	return err
}

// Clear는 사용자의 장바구니를 비웁니다.
func (r *CartRepository) Clear(ctx context.Context, userID string) error {
	const q = `DELETE FROM cart_items WHERE user_id = $1`
	_, err := r.pool.Exec(ctx, q, userID)
	return err
}

// ============================================================================
// OrderRepository
// ============================================================================

// OrderRepository는 PostgreSQL 기반 OrderRepository 구현입니다.
type OrderRepository struct {
	pool *pgxpool.Pool
}

// NewOrderRepository는 OrderRepository를 생성합니다.
func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

// Create는 주문 및 주문 항목을 생성합니다 (트랜잭션).
func (r *OrderRepository) Create(ctx context.Context, order *service.Order) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	const orderQ = `INSERT INTO orders (id, user_id, total_price_krw, status, shipping_address, payment_id, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err = tx.Exec(ctx, orderQ,
		order.ID, order.UserID, order.TotalPriceKRW, int32(order.Status),
		order.ShippingAddress, nilIfEmpty(order.PaymentID),
		order.CreatedAt, order.UpdatedAt,
	)
	if err != nil {
		return err
	}

	const itemQ = `INSERT INTO order_items (id, order_id, product_id, product_name, quantity, unit_price_krw, total_price_krw)
		VALUES (gen_random_uuid(), $1,$2,$3,$4,$5,$6)`
	for _, item := range order.Items {
		_, err = tx.Exec(ctx, itemQ,
			order.ID, item.ProductID, item.ProductName,
			item.Quantity, item.UnitPriceKRW, item.TotalPriceKRW,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// GetByID는 주문 ID로 조회합니다 (주문 항목 포함).
func (r *OrderRepository) GetByID(ctx context.Context, id string) (*service.Order, error) {
	const orderQ = `SELECT id, user_id, total_price_krw, status, COALESCE(shipping_address, ''),
		COALESCE(payment_id::text, ''), created_at, updated_at
		FROM orders WHERE id = $1`

	var o service.Order
	err := r.pool.QueryRow(ctx, orderQ, id).Scan(
		&o.ID, &o.UserID, &o.TotalPriceKRW, &o.Status,
		&o.ShippingAddress, &o.PaymentID, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	items, err := r.getOrderItems(ctx, id)
	if err != nil {
		return nil, err
	}
	o.Items = items
	return &o, nil
}

// ListByUserID는 사용자의 주문 목록을 조회합니다.
func (r *OrderRepository) ListByUserID(ctx context.Context, userID string, limit, offset int32) ([]*service.Order, int32, error) {
	var total int32
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM orders WHERE user_id = $1`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	const q = `SELECT id, user_id, total_price_krw, status, COALESCE(shipping_address, ''),
		COALESCE(payment_id::text, ''), created_at, updated_at
		FROM orders WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []*service.Order
	for rows.Next() {
		var o service.Order
		if err := rows.Scan(
			&o.ID, &o.UserID, &o.TotalPriceKRW, &o.Status,
			&o.ShippingAddress, &o.PaymentID, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		orders = append(orders, &o)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	// 각 주문의 항목을 로드
	for _, o := range orders {
		items, err := r.getOrderItems(ctx, o.ID)
		if err != nil {
			return nil, 0, err
		}
		o.Items = items
	}

	return orders, total, nil
}

// UpdateStatus는 주문 상태를 업데이트합니다.
func (r *OrderRepository) UpdateStatus(ctx context.Context, id string, status service.OrderStatus) error {
	const q = `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.pool.Exec(ctx, q, int32(status), id)
	return err
}

// getOrderItems는 주문의 항목 목록을 조회합니다.
func (r *OrderRepository) getOrderItems(ctx context.Context, orderID string) ([]*service.OrderItem, error) {
	const q = `SELECT product_id, product_name, quantity, unit_price_krw, total_price_krw
		FROM order_items WHERE order_id = $1`
	rows, err := r.pool.Query(ctx, q, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*service.OrderItem
	for rows.Next() {
		var item service.OrderItem
		if err := rows.Scan(
			&item.ProductID, &item.ProductName,
			&item.Quantity, &item.UnitPriceKRW, &item.TotalPriceKRW,
		); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, rows.Err()
}

// nilIfEmpty는 빈 문자열을 nil로 변환합니다 (UUID 컬럼용).
func nilIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
