-- =============================================================================
-- ManPaSik Phase 2: ai-inference-service 테이블 초기화
-- =============================================================================

-- AI 모델 타입 ENUM
DO $$ BEGIN
  CREATE TYPE ai_model_type AS ENUM (
    'BIOMARKER_CLASSIFIER',
    'ANOMALY_DETECTOR',
    'TREND_PREDICTOR',
    'HEALTH_SCORER',
    'FOOD_CALORIE_ESTIMATOR'
  );
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- 리스크 레벨 ENUM
DO $$ BEGIN
  CREATE TYPE risk_level AS ENUM ('LOW', 'MODERATE', 'HIGH', 'CRITICAL');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- 분석 결과 테이블
CREATE TABLE IF NOT EXISTS analysis_results (
  id               VARCHAR(64)   PRIMARY KEY,
  user_id          VARCHAR(64)   NOT NULL,
  measurement_id   VARCHAR(64)   NOT NULL,
  overall_health_score DOUBLE PRECISION NOT NULL DEFAULT 0.0,
  summary          TEXT,
  analyzed_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_analysis_results_user_id        ON analysis_results (user_id);
CREATE INDEX IF NOT EXISTS idx_analysis_results_measurement_id ON analysis_results (measurement_id);

-- 바이오마커 결과 (분석 당 N개)
CREATE TABLE IF NOT EXISTS biomarker_results (
  id               BIGSERIAL     PRIMARY KEY,
  analysis_id      VARCHAR(64)   NOT NULL REFERENCES analysis_results(id),
  biomarker_name   VARCHAR(128)  NOT NULL,
  value            DOUBLE PRECISION NOT NULL,
  unit             VARCHAR(32)   NOT NULL,
  classification   VARCHAR(32)   NOT NULL,  -- normal, borderline, abnormal
  confidence       DOUBLE PRECISION NOT NULL,
  risk_level       risk_level    NOT NULL DEFAULT 'LOW',
  reference_range  VARCHAR(128)
);

CREATE INDEX IF NOT EXISTS idx_biomarker_results_analysis_id ON biomarker_results (analysis_id);

-- 이상치 플래그
CREATE TABLE IF NOT EXISTS anomaly_flags (
  id               BIGSERIAL     PRIMARY KEY,
  analysis_id      VARCHAR(64)   NOT NULL REFERENCES analysis_results(id),
  metric_name      VARCHAR(128)  NOT NULL,
  value            DOUBLE PRECISION NOT NULL,
  expected_min     DOUBLE PRECISION NOT NULL,
  expected_max     DOUBLE PRECISION NOT NULL,
  anomaly_score    DOUBLE PRECISION NOT NULL,
  description      TEXT
);

CREATE INDEX IF NOT EXISTS idx_anomaly_flags_analysis_id ON anomaly_flags (analysis_id);

-- 건강 점수 히스토리
CREATE TABLE IF NOT EXISTS health_scores (
  id               BIGSERIAL     PRIMARY KEY,
  user_id          VARCHAR(64)   NOT NULL,
  overall_score    DOUBLE PRECISION NOT NULL,
  category_scores  JSONB,
  trend            VARCHAR(32)   NOT NULL,  -- improving, stable, declining
  recommendation   TEXT,
  calculated_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_health_scores_user_id ON health_scores (user_id);

-- AI 모델 레지스트리
CREATE TABLE IF NOT EXISTS ai_models (
  model_type       ai_model_type PRIMARY KEY,
  name             VARCHAR(128)  NOT NULL,
  version          VARCHAR(32)   NOT NULL,
  description      TEXT,
  accuracy         DOUBLE PRECISION NOT NULL DEFAULT 0.0,
  last_trained     TIMESTAMPTZ,
  status           VARCHAR(32)   NOT NULL DEFAULT 'active'  -- active, training, deprecated
);

-- 초기 모델 데이터
INSERT INTO ai_models (model_type, name, version, description, accuracy, last_trained, status) VALUES
  ('BIOMARKER_CLASSIFIER', 'BiomarkerClassifier', '1.0.0', '바이오마커 분류 모델', 0.942, NOW() - INTERVAL '1 day', 'active'),
  ('ANOMALY_DETECTOR', 'AnomalyDetector', '1.0.0', '이상치 탐지 모델', 0.918, NOW() - INTERVAL '1 day', 'active'),
  ('TREND_PREDICTOR', 'TrendPredictor', '1.0.0', '트렌드 예측 모델', 0.876, NOW() - INTERVAL '1 day', 'active'),
  ('HEALTH_SCORER', 'HealthScorer', '1.0.0', '건강 점수 산출 모델', 0.905, NOW() - INTERVAL '1 day', 'active'),
  ('FOOD_CALORIE_ESTIMATOR', 'FoodCalorieEstimator', '0.9.0-beta', '음식 칼로리 추정 모델 (Phase 2 후반)', 0.823, NOW() - INTERVAL '1 day', 'training')
ON CONFLICT (model_type) DO NOTHING;
