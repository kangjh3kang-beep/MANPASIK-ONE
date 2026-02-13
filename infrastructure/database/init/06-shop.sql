-- Shop Service Database Initialization
-- 테이블은 POSTGRES_DB(기본 manpasik)에 생성됩니다.

-- 상품 테이블
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    category INTEGER NOT NULL DEFAULT 0,  -- 0: Unknown, 1: Cartridge, 2: Reader, 3: Accessory, 4: Bundle
    price_krw INTEGER NOT NULL DEFAULT 0,
    stock INTEGER NOT NULL DEFAULT 0,
    image_url TEXT DEFAULT '',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_products_category ON products (category);
CREATE INDEX IF NOT EXISTS idx_products_is_active ON products (is_active);

-- 기본 상품 데이터
INSERT INTO products (id, name, description, category, price_krw, stock, is_active) VALUES
    ('00000000-0000-0000-0000-000000000001', '혈당 카트리지 10팩', '혈당 측정용 카트리지 10개입', 1, 29000, 1000, TRUE),
    ('00000000-0000-0000-0000-000000000002', '콜레스테롤 카트리지 10팩', '콜레스테롤 측정용 카트리지 10개입', 1, 35000, 800, TRUE),
    ('00000000-0000-0000-0000-000000000003', '헤모글로빈 카트리지 10팩', '헤모글로빈 측정용 카트리지 10개입', 1, 32000, 600, TRUE),
    ('00000000-0000-0000-0000-000000000004', 'ManPaSik 리더기 V2', '차동측정 기반 범용 분석 리더기 2세대', 2, 199000, 200, TRUE),
    ('00000000-0000-0000-0000-000000000005', '리더기 V2 + 카트리지 번들', '리더기 V2 + 혈당·콜레스테롤 카트리지 각 10팩', 4, 249000, 100, TRUE),
    ('00000000-0000-0000-0000-000000000006', '리더기 보호 케이스', '프리미엄 실리콘 보호 케이스', 3, 25000, 500, TRUE)
ON CONFLICT (id) DO NOTHING;

-- 장바구니 테이블
CREATE TABLE IF NOT EXISTS cart_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    product_id UUID NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_cart_items_user_id ON cart_items (user_id);

-- 주문 테이블
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    total_price_krw INTEGER NOT NULL DEFAULT 0,
    status INTEGER NOT NULL DEFAULT 1,  -- 1: Pending, 2: Paid, 3: Shipped, 4: Delivered, 5: Cancelled, 6: Refunded
    shipping_address TEXT DEFAULT '',
    payment_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders (user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders (status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders (created_at);

-- 주문 항목 테이블
CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price_krw INTEGER NOT NULL DEFAULT 0,
    total_price_krw INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items (order_id);
