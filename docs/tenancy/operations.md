# 만파식 멀티테넌시 운영 가이드 v2 (Phase AH-3, AL-3 통합)

## 개요

만파식 시스템의 멀티테넌트 격리는 **A'~AL 38+ 단계** 에 걸쳐 구축됨.
이 문서는 **운영자/보안 담당자/온콜 엔지니어** 가 일상 운영에서 참조하는 매뉴얼.

## v2 주요 변경 사항 (Phase AI~AL)

- **LLM Audit 영속화**: 모든 LLM 호출이 `llm_audit` 테이블에 자동 기록 (BatchAuditLog 5초 flush)
- **동적 Quota**: tenant 별 일/월 토큰 한도 + 60초 캐시 (PostgresQuotaStore 영속)
- **Webhook 자동 발송**: Invitation/Membership 변경 → Slack/Generic 외부 알림
- **Webhook 통계**: sent/failed/retried/dropped Prometheus 메트릭
- **REST API**: `/ops/tenancy/quota/{tid}` GET/PUT/DELETE 동적 운영

## 운영 환경변수 (Phase AL-1, AL-2)

| 변수 | 기본값 | 설명 |
|------|--------|------|
| `WEBHOOK_URL` | (없음) | Slack incoming webhook 또는 일반 webhook URL. 미설정 시 webhook 비활성 |
| `WEBHOOK_MODE` | `generic` | `slack` 또는 `generic` — 페이로드 형식 결정 |
| `DB_HOST` | (없음) | 설정 시 PostgreSQL 영속 저장소 자동 활성 (Tenancy/Quota/Audit 모두) |
| `TENANCY_ENFORCED` | `false` | `true` 시 9개 보호 서비스에서 인터셉터 활성 |

## 1. 아키텍처 요약

```
[Flutter 앱]
  └─ TenantInterceptor: SharedPreferences.active_tenant_id → X-Tenant-ID 헤더
       │
[Gateway :8080]
  ├─ TenantPropagation: HTTP 헤더 → ctx → outgoing gRPC metadata
  └─ /api/v1/tenancy/*: Invitation/Membership REST API (단독 마운트)
       │
[gRPC 백엔드 9개 서비스]
  ├─ tenancy.UnaryInterceptor: ctx tenant_id 검증 → 비멤버 거부
  ├─ tenancy.PolicyEngine: ActionRead/Write/Admin 권한 검사
  └─ Repository: BuildTenantClause → SELECT/UPDATE 자동 격리
       │
[PostgreSQL]
  ├─ measurement_data.tenant_id (Phase AG-1)
  ├─ health_records.tenant_id    (Phase AH-1)
  ├─ tenant_memberships          (Phase Y-2)
  └─ tenant_invitations          (Phase AA-1)
```

**핵심 원칙**: tenant 미설정 (`X-Tenant-ID` 헤더 없음) 시 `tenant_id IS NULL`
인 데이터만 반환 → legacy/personal 데이터 격리 + 다른 조직 데이터 차단.

## 2. 일반 운영 시나리오

### 2.1 새 조직 (병원/가족) 추가

운영자가 직접 SQL 또는 Admin UI 사용:
```sql
INSERT INTO tenant_memberships (user_id, tenant_id, role, active, joined_at)
VALUES ('admin-user-id', 'hospitalA', 'owner', true, EXTRACT(EPOCH FROM NOW())::BIGINT);
```

또는 family-service 가 자동 생성 (TenancySync):
- 사용자가 가족 그룹 생성 → owner 자동 등록
- 멤버 가입 → member 자동 등록 (Phase AE-2)

### 2.2 멤버 초대 흐름

```
1. admin → POST /api/v1/tenancy/invitations
2. (자동) InvitationService → InvitationNotifier → 이메일/카카오 발송
   ├─ 본문 4개 언어 자동 (Phase AF-3)
   └─ LocaleResolver 등록 시 invitee.lang_code 기반 (Phase AG-3)
3. invitee → 딥링크 또는 토큰 입력 → POST /accept
4. (자동) MembershipStore.Add + Invitation status=accepted
```

### 2.3 만료 초대 자동 정리

`InvitationCleaner` (Phase AC-1) 가 1시간 주기로 실행:
- pending + ExpiresAt 지난 초대 → status=expired
- 운영 중단 없이 백그라운드 동작

수동 강제 정리:
```bash
# admin-service 컨테이너에서
curl -X POST http://localhost:9100/ops/tenancy/cleanup
```
(현재 미구현 — InvitationCleaner.MarkExpiredAll 직접 호출 필요)

## 3. 트러블슈팅

### 3.1 "조직 데이터가 안 보여요"

**증상**: 사용자가 자신의 조직 데이터를 못 봄.

**진단**:
1. `X-Tenant-ID` 헤더 확인 (gateway 로그)
2. 멤버십 확인: `SELECT * FROM tenant_memberships WHERE user_id = ? AND active = true;`
3. 데이터 tenant_id 확인: `SELECT DISTINCT tenant_id FROM measurement_data WHERE user_id = ?;`

**원인 후보**:
- 클라이언트가 `TenantInterceptor.setActiveTenant` 미호출 → 헤더 없음 → IS NULL 필터
- 멤버십이 비활성 (active=false)
- 데이터가 다른 tenant_id 로 저장됨 (이력 데이터)

### 3.2 "다른 조직 데이터가 보여요" (격리 위반)

**즉시 대응**:
1. 영향 범위 측정: 어떤 사용자가 어떤 데이터에 접근?
2. 임시 차단: 해당 사용자 멤버십 비활성화
   ```sql
   UPDATE tenant_memberships SET active = false WHERE user_id = ?;
   ```
3. 감사 로그 확인: gateway 액세스 로그에서 헤더/응답 검증

**근본 원인 분석**:
- Repository 가 `BuildTenantClause` 미사용 (개발자 누락)
- ctx 가 잘못된 tenant 로 설정됨
- TenancyAdapter (LLM) 의 audit 로그 확인

### 3.3 "초대 알림이 안 와요"

1. invitation 상태: `SELECT * FROM tenant_invitations WHERE token = ?;`
2. notification-service 로그: `Dispatch failed` 또는 `contact lookup failed`
3. invitee_hint 가 이메일/전화 형식인지
4. SMS/Email 어댑터 healthcheck

## 4. Prometheus 알림 대응 (Phase AF-2)

| Alert | 의미 | 대응 |
|-------|------|------|
| `TenancyExpiredInvitationsPiling` | 만료 초대 >20 누적 | InvitationCleaner 동작 확인. `kubectl logs gateway` 에서 cycle count 확인 |
| `TenancyPendingInvitationsStuck` | pending >50 (1시간 지속) | 알림 발송 실패 의심. notification-service 헬스체크 |
| `TenancyZeroMemberTenants` | 모든 조직 멤버 0명 | **CRITICAL**. DB 데이터 손실 가능성. 즉시 백업 복구 검토 |
| `TenancyHighInactiveMemberRatio` | 비활성 >30% | 정상 패턴인지 확인 (대량 이벤트 vs 시스템 오류) |
| `TenancyAdminCountTooLow` | admin/owner 부족 | 일부 조직에 권한 위임 필요 |

## 5. 성능 튜닝

### 5.1 인덱스 활용 확인

```sql
EXPLAIN ANALYZE
SELECT * FROM measurement_data
WHERE user_id = 'u1' AND tenant_id = 'hospA' AND time > NOW() - INTERVAL '7 days';
```

복합 인덱스 `idx_measurement_data_user_tenant` 사용 확인 (Index Scan).
Sequential Scan 발생 시 ANALYZE/REINDEX 검토.

### 5.2 LLM Quota 활용 (Phase AH-2)

```go
quota := llm.NewMemoryQuota(100000) // tenant 당 100k 토큰/일
adapter := llm.NewTenancyAdapter(openaiAdapter, auditLog)
adapter.SetQuota(quota)
```
- 대용량 사용 tenant 식별: `MemoryAuditLog.TotalTokensByTenant`
- Grafana 대시보드: `manpasik_llm_tokens_total{tenant=...}` (별도 wire-up 필요)

## 6. 보안 체크리스트

운영 배포 전 확인:

- [ ] `TENANCY_ENFORCED=true` (production overlay)
- [ ] PolicyEngine 등록 (HTTPHandler + InvitationService 모두)
- [ ] TenantInterceptor 미들웨어가 모든 보호 서비스에 연결됨
- [ ] Repository 의 SELECT/UPDATE 가 `BuildTenantClause` 호출
- [ ] `Sanitize TenantQuery` 가 사용자 입력 SQL 부분에 적용됨
- [ ] LLM 어댑터가 `TenancyAdapter` 로 wrapping 됨
- [ ] Prometheus alert rules 활성 (alert_rules.yml tenancy 그룹)
- [ ] Vault Bootstrap 자동 회전 활성 (운영 환경)
- [ ] Audit log 영속화 (DB 또는 외부 저장소)

## 7. 정기 점검 (월간)

1. 만료 초대 정리 효율: ExpiredInvitationCount / TotalInvitations
2. 멤버십 분포: SELECT role, COUNT(*) FROM tenant_memberships GROUP BY role
3. tenant 별 데이터 크기: `SELECT tenant_id, pg_size_pretty(...) FROM measurement_data;`
4. 격리 위반 시도 로그 (gateway access logs 의 403)
5. LLM 비용 분배: tenant 별 토큰 사용량 → 청구 분리 가능

## 참고 자료

- [data-isolation.md](data-isolation.md) — 격리 기술 명세
- [/docs/deeplinks/](../deeplinks/) — Universal Link 운영
- [Phase AG-AH 변경 이력](../../CONTEXT.md) — 최신 보안 강화 내역
