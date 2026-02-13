-- =============================================================================
-- Community Service Schema (Phase 3B)
-- =============================================================================

-- 게시글 카테고리
CREATE TYPE post_category AS ENUM (
    'general', 'question', 'tip', 'review', 'experience', 'recipe', 'exercise'
);

-- 챌린지 상태
CREATE TYPE challenge_status AS ENUM (
    'upcoming', 'active', 'completed', 'cancelled'
);

-- 챌린지 유형
CREATE TYPE challenge_type AS ENUM (
    'measurement', 'exercise', 'nutrition', 'sleep', 'weight', 'custom'
);

-- 게시글
CREATE TABLE IF NOT EXISTS posts (
    post_id             VARCHAR(36) PRIMARY KEY,
    author_user_id      VARCHAR(36) NOT NULL,
    author_display_name VARCHAR(100),
    category            post_category NOT NULL DEFAULT 'general',
    title               VARCHAR(300) NOT NULL,
    content             TEXT NOT NULL,
    tags                TEXT[],
    like_count          INTEGER DEFAULT 0,
    comment_count       INTEGER DEFAULT 0,
    view_count          INTEGER DEFAULT 0,
    is_anonymous        BOOLEAN DEFAULT false,
    is_pinned           BOOLEAN DEFAULT false,
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);

-- 댓글
CREATE TABLE IF NOT EXISTS comments (
    comment_id          VARCHAR(36) PRIMARY KEY,
    post_id             VARCHAR(36) NOT NULL REFERENCES posts(post_id) ON DELETE CASCADE,
    author_user_id      VARCHAR(36) NOT NULL,
    author_display_name VARCHAR(100),
    content             TEXT NOT NULL,
    parent_comment_id   VARCHAR(36) REFERENCES comments(comment_id),
    like_count          INTEGER DEFAULT 0,
    created_at          TIMESTAMPTZ DEFAULT NOW()
);

-- 좋아요
CREATE TABLE IF NOT EXISTS post_likes (
    post_id     VARCHAR(36) NOT NULL REFERENCES posts(post_id) ON DELETE CASCADE,
    user_id     VARCHAR(36) NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (post_id, user_id)
);

-- 챌린지
CREATE TABLE IF NOT EXISTS challenges (
    challenge_id        VARCHAR(36) PRIMARY KEY,
    creator_user_id     VARCHAR(36) NOT NULL,
    type                challenge_type NOT NULL DEFAULT 'custom',
    title               VARCHAR(300) NOT NULL,
    description         TEXT,
    goal_description    TEXT,
    target_value        DECIMAL(12,2),
    target_unit         VARCHAR(50),
    status              challenge_status NOT NULL DEFAULT 'upcoming',
    participant_count   INTEGER DEFAULT 0,
    max_participants    INTEGER DEFAULT 100,
    duration_days       INTEGER DEFAULT 30,
    start_date          TIMESTAMPTZ NOT NULL,
    end_date            TIMESTAMPTZ NOT NULL,
    created_at          TIMESTAMPTZ DEFAULT NOW()
);

-- 챌린지 참가자
CREATE TABLE IF NOT EXISTS challenge_participants (
    challenge_id    VARCHAR(36) NOT NULL REFERENCES challenges(challenge_id) ON DELETE CASCADE,
    user_id         VARCHAR(36) NOT NULL,
    joined_at       TIMESTAMPTZ DEFAULT NOW(),
    progress        DECIMAL(5,2) DEFAULT 0.00,
    PRIMARY KEY (challenge_id, user_id)
);

-- 인덱스
CREATE INDEX IF NOT EXISTS idx_posts_author ON posts(author_user_id);
CREATE INDEX IF NOT EXISTS idx_posts_category ON posts(category);
CREATE INDEX IF NOT EXISTS idx_posts_created ON posts(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_posts_tags ON posts USING GIN(tags);
CREATE INDEX IF NOT EXISTS idx_comments_post ON comments(post_id);
CREATE INDEX IF NOT EXISTS idx_comments_parent ON comments(parent_comment_id);
CREATE INDEX IF NOT EXISTS idx_challenges_status ON challenges(status);
CREATE INDEX IF NOT EXISTS idx_challenges_type ON challenges(type);
CREATE INDEX IF NOT EXISTS idx_challenges_start ON challenges(start_date);
