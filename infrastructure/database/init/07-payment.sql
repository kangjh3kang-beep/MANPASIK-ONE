-- =============================================================================
-- ManPaSik Phase 2: payment-service 테이블 초기화
-- =============================================================================

-- 결제 유형 ENUM
DO $$ BEGIN
  CREATE TYPE payment_type AS ENUM ('CARD', 'BANK_TRANSFER', 'VIRTUAL_ACCOUNT', 'MOBILE');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- 결제 상태 ENUM
DO $$ BEGIN
  CREATE TYPE payment_status AS ENUM ('PENDING', 'CONFIRMED', 'FAILED', 'REFUNDED', 'PARTIAL_REFUND');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- 결제 테이블
CREATE TABLE IF NOT EXISTS payments (
  id            VARCHAR(64)   PRIMARY KEY,
  user_id       VARCHAR(64)   NOT NULL,
  order_id      VARCHAR(64)   NOT NULL,
  amount        BIGINT        NOT NULL CHECK (amount >= 0),     -- 원 단위
  currency      VARCHAR(8)    NOT NULL DEFAULT 'KRW',
  type          payment_type  NOT NULL DEFAULT 'CARD',
  status        payment_status NOT NULL DEFAULT 'PENDING',
  pg_tx_id      VARCHAR(128),                                    -- PG사 거래 ID
  description   TEXT,
  metadata      JSONB,
  created_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
  updated_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payments_user_id   ON payments (user_id);
CREATE INDEX IF NOT EXISTS idx_payments_order_id  ON payments (order_id);
CREATE INDEX IF NOT EXISTS idx_payments_status    ON payments (status);

-- 환불 테이블
CREATE TABLE IF NOT EXISTS refunds (
  id            VARCHAR(64)   PRIMARY KEY,
  payment_id    VARCHAR(64)   NOT NULL REFERENCES payments(id),
  amount        BIGINT        NOT NULL CHECK (amount > 0),
  reason        TEXT,
  status        VARCHAR(32)   NOT NULL DEFAULT 'COMPLETED',
  created_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_refunds_payment_id ON refunds (payment_id);
