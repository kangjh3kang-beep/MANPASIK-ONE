-- =============================================================================
-- Prescription fulfillment (약국 전달·조제 상태)
-- 24: 처방전 이행 컬럼 및 이력 (19-prescription.sql 보강)
-- =============================================================================

-- prescriptions 테이블에 이행 관련 컬럼 추가
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'prescriptions') THEN
        IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'prescriptions' AND column_name = 'fulfillment_type') THEN
            ALTER TABLE prescriptions ADD COLUMN fulfillment_type VARCHAR(20) DEFAULT '';
        END IF;
        IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'prescriptions' AND column_name = 'shipping_address_id') THEN
            ALTER TABLE prescriptions ADD COLUMN shipping_address_id UUID;
        END IF;
        IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'prescriptions' AND column_name = 'fulfillment_token') THEN
            ALTER TABLE prescriptions ADD COLUMN fulfillment_token VARCHAR(64) DEFAULT '';
        END IF;
        IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'prescriptions' AND column_name = 'dispensary_status') THEN
            ALTER TABLE prescriptions ADD COLUMN dispensary_status VARCHAR(32) DEFAULT '';
        END IF;
        IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'prescriptions' AND column_name = 'sent_to_pharmacy_at') THEN
            ALTER TABLE prescriptions ADD COLUMN sent_to_pharmacy_at TIMESTAMPTZ;
        END IF;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_prescriptions_fulfillment_token ON prescriptions(fulfillment_token) WHERE fulfillment_token != '';
CREATE INDEX IF NOT EXISTS idx_prescriptions_pharmacy ON prescriptions(pharmacy_id) WHERE pharmacy_id IS NOT NULL;

-- 조제 이행 이력 (감사/상태 추적)
CREATE TABLE IF NOT EXISTS prescription_fulfillment_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    prescription_id UUID NOT NULL REFERENCES prescriptions(id) ON DELETE CASCADE,
    event_type      VARCHAR(32) NOT NULL,  -- 'sent_to_pharmacy', 'dispensary_accepted', 'dispensed', 'picked_up', 'shipped'
    pharmacy_id     UUID,
    status_before   VARCHAR(32),
    status_after    VARCHAR(32),
    actor_type      VARCHAR(20) DEFAULT 'system',  -- 'system', 'pharmacy', 'patient'
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_fulfillment_logs_prescription ON prescription_fulfillment_logs(prescription_id);
CREATE INDEX IF NOT EXISTS idx_fulfillment_logs_created ON prescription_fulfillment_logs(created_at);
