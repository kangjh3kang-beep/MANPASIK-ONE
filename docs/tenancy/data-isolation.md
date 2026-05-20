# 데이터 격리 기술 명세 (Phase AG-AH)

만파식 멀티테넌트 데이터 격리의 **계층별 구현 명세**.

## 격리 계층

```
┌─────────────────────────────────────────────────────────┐
│ L1: Network    Flutter TenantInterceptor (Dio)          │
│                X-Tenant-ID 헤더 자동 부착                  │
├─────────────────────────────────────────────────────────┤
│ L2: Edge       Gateway TenantPropagation (HTTP→ctx)     │
│                ctx → outgoing gRPC metadata             │
├─────────────────────────────────────────────────────────┤
│ L3: Service    grpc Interceptor + PolicyEngine          │
│                membership 검증 + Action 권한 검사            │
├─────────────────────────────────────────────────────────┤
│ L4: Repository BuildTenantClause                         │
│                SELECT WHERE tenant_id = $N AND ...      │
├─────────────────────────────────────────────────────────┤
│ L5: Database   인덱스 + CHECK 제약                          │
│                idx_*_user_tenant, role CHECK            │
└─────────────────────────────────────────────────────────┘
```

각 계층이 독립적으로 검증하므로 어느 한 계층 우회해도 다음 계층에서 차단.

## 계층별 명세

### L1: Flutter TenantInterceptor

**위치**: `frontend/flutter-app/lib/core/network/tenant_interceptor.dart`

**동작**:
1. SharedPreferences 의 `active_tenant_id` 키 읽음
2. `X-Tenant-ID: <value>` 헤더 자동 부착 (Dio interceptor)
3. 사용자 변경 시 `setActiveTenant(null)` → 헤더 제거 (개인 모드)

**검증**: `test/core/network/tenant_interceptor_test.dart` 5건

### L2: Gateway TenantPropagation

**위치**: `backend/services/gateway/internal/middleware/tenant.go`

**동작**:
1. HTTP 요청에서 `X-Tenant-ID` 추출
2. `tenancy.WithTenant(ctx, tid)` 로 ctx 주입
3. gRPC 호출 시 `WithTenantMetadata(ctx)` → outgoing metadata
4. 백엔드 서비스가 `metadata.FromIncomingContext` 로 복원

**검증**: `services/gateway/internal/middleware/tenant_test.go` 5건

### L3: gRPC Interceptor

**위치**: `backend/shared/tenancy/interceptor.go` + `helpers.go`

**동작 (TENANCY_ENFORCED=true 시)**:
1. ctx 에서 user_id, tenant_id 추출
2. `MembershipStore.Get(user, tenant)` 조회
3. 멤버 아니면 PermissionDenied
4. 비활성 멤버면 Unauthenticated

**적용 서비스 (9개)**:
auth, telemedicine, health-record, family,
notification, payment, measurement, coaching, community

### L4: Repository (BuildTenantClause)

**위치**: `backend/shared/tenancy/isolation.go`

**SQL 적용 패턴**:
```go
args := []interface{}{userID}
clause, tArgs := tenancy.BuildTenantClause(ctx, "md", len(args)+1)
args = append(args, tArgs...)
q := "SELECT * FROM measurement_data md WHERE md.user_id = $1" + clause
```

**적용 테이블**:
- `measurement_data` (Phase AG-1)
- `health_records` (Phase AH-1)
- `prescriptions` (Phase AG-1, schema 만 — repository 적용 예정)

**격리 동작**:
- ctx tenant 있음: `AND md.tenant_id = $2`
- ctx tenant 없음: `AND md.tenant_id IS NULL` (legacy/personal 만)

### L5: Database 스키마

**마이그레이션 SQL**:
- `41-tenancy.sql` — tenant_memberships, tenant_invitations 테이블
- `42-tenant-isolation.sql` — measurement_data/health_records/prescriptions tenant_id 컬럼 + 인덱스

**제약**:
```sql
ALTER TABLE tenant_memberships ADD CONSTRAINT
  tenant_memberships_role_check CHECK
  (role IN ('owner', 'admin', 'medical_staff', 'member', 'viewer'));
```

**인덱스**:
- `idx_measurement_data_user_tenant(user_id, tenant_id, time DESC)`
- `idx_health_records_user_tenant(user_id, tenant_id)`
- `idx_tenant_memberships_user(user_id) WHERE active = true`

## LLM 격리 (Phase AH-2)

LLM 호출은 **tenant 별 감사 + 비용 분리**를 위해 별도 어댑터 사용:

```go
inner := llm.NewOpenAIAdapter(apiKey, baseURL)
audit := llm.NewMemoryAuditLog() // 또는 영속 구현
adapter := llm.NewTenancyAdapter(inner, audit)
adapter.SetQuota(llm.NewMemoryQuota(100000)) // tenant 당 100k 토큰

// MedicalLLMResponder 에 주입
responder := assistant.NewMedicalLLMResponder(assistant.MedicalLLMConfig{
    Adapter: adapter,
})
```

**감사 항목**:
- TenantID, UserID
- Provider, Model
- PromptTokens, CompletionTokens, TotalTokens
- LatencyMs, Success, ErrorMessage

## 격리 위반 검증 (자동)

`backend/integration_test/tenancy_isolation_test.go` (9 시나리오):

| Test | 시나리오 | 기대 동작 |
|------|---------|----------|
| AdminListsOwnTenant | A 조직 admin → A 멤버 조회 | 200 OK |
| AdminListsOtherTenant | A admin → B 멤버 조회 | 403 Forbidden |
| NonAdminListsTenant | A 일반 멤버 → A 멤버 조회 | 403 (admin 권한 부족) |
| NonMemberListsTenant | 외부인 → A 멤버 조회 | 403 |
| CrossTenantRoleChange | A admin → B 멤버 역할 변경 | 403 |
| CrossTenantMemberRemoval | A admin → B 멤버 제거 | 403 |
| Unauthenticated | 인증 없음 | 401 |
| OwnerOverride | 본인 데이터 write | OK |
| PolicyEngine_CrossTenant | A admin → B 리소스 read | Decision.Allowed=false |

## 운영 시나리오 보장

위 5계층의 결과로 다음 시나리오는 **자동으로 차단**:

1. ✅ A 병원 의사가 B 병원 환자 데이터 조회 → L4 (SQL) 또는 L3 (인터셉터) 에서 차단
2. ✅ 가족 그룹 외부 사용자가 가족 측정 조회 → L1+L4
3. ✅ Admin 권한 없는 멤버가 조직 멤버 목록 조회 → L3 (PolicyEngine ActionAdmin)
4. ✅ 만료된 초대로 가입 시도 → InvitationService.Accept 거부
5. ✅ LLM 비용이 tenant 간 섞이지 않음 → TenancyAdapter audit
