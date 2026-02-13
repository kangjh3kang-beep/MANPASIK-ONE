-- Performance indexes for common queries
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_sessions_user_started ON measurement_sessions(user_id, started_at DESC);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_devices_user_status ON devices(user_id, status);
