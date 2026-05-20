-- nlp-service 초기화 스크립트
-- Phase A: NLP 건강 쿼리 영속화

-- 건강 쿼리 테이블
CREATE TABLE IF NOT EXISTS nlp_health_queries (
    id              VARCHAR(64) PRIMARY KEY,
    user_id         VARCHAR(64) NOT NULL,
    raw_text        TEXT NOT NULL,
    intent          VARCHAR(100),
    entities        JSONB DEFAULT '[]',
    confidence      DOUBLE PRECISION DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_nlp_queries_user ON nlp_health_queries(user_id);
CREATE INDEX IF NOT EXISTS idx_nlp_queries_intent ON nlp_health_queries(intent);

-- 증상 추출 테이블
CREATE TABLE IF NOT EXISTS nlp_symptom_extractions (
    id              VARCHAR(64) PRIMARY KEY,
    text            TEXT NOT NULL,
    symptoms        JSONB DEFAULT '[]',
    processed_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- NLP 제안 테이블
CREATE TABLE IF NOT EXISTS nlp_suggestions (
    id              VARCHAR(64) PRIMARY KEY,
    query_id        VARCHAR(64) NOT NULL REFERENCES nlp_health_queries(id),
    text            TEXT NOT NULL,
    category        VARCHAR(50),
    priority        INT NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_nlp_suggestions_query ON nlp_suggestions(query_id);
