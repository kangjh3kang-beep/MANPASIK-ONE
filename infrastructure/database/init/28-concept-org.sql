-- ============================================================================
-- 목적별 컨셉 + 조직 (Concept/Organization Service) — Phase 3
-- ============================================================================

-- 컨셉
CREATE TABLE IF NOT EXISTS concepts (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT DEFAULT '',
    category VARCHAR(30) DEFAULT 'personal',
    icon_url VARCHAR(500) DEFAULT '',
    owner_id VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_concepts_owner ON concepts(owner_id);
CREATE INDEX IF NOT EXISTS idx_concepts_category ON concepts(category);

-- 컨셉-기기 할당
CREATE TABLE IF NOT EXISTS concept_device_assignments (
    id VARCHAR(64) PRIMARY KEY,
    concept_id VARCHAR(64) NOT NULL REFERENCES concepts(id) ON DELETE CASCADE,
    device_id VARCHAR(64) NOT NULL,
    assigned_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(concept_id, device_id)
);

-- 조직
CREATE TABLE IF NOT EXISTS organizations (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    type VARCHAR(30) DEFAULT 'enterprise',
    description TEXT DEFAULT '',
    owner_id VARCHAR(64) NOT NULL,
    member_count INT DEFAULT 1,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_organizations_owner ON organizations(owner_id);

-- 조직 멤버
CREATE TABLE IF NOT EXISTS organization_members (
    id VARCHAR(64) PRIMARY KEY,
    organization_id VARCHAR(64) NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id VARCHAR(64) NOT NULL,
    role VARCHAR(20) DEFAULT 'member',
    joined_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(organization_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_org_members_org ON organization_members(organization_id);
CREATE INDEX IF NOT EXISTS idx_org_members_user ON organization_members(user_id);
