# ManPaSik 배포 전략 명세서 (Deployment Strategy Specification)

**문서번호**: MPK-DEPLOY-STRATEGY-v1.0  
**갱신일**: 2026-02-12  
**목적**: Phase별 배포 환경, 방식, 파이프라인, 롤백 전략을 정의  
**참조**: .github/workflows/ci.yml, .github/workflows/cd.yml, infrastructure/kubernetes/

---

## 1. 환경 구성

### 1.1 환경 목록

| 환경 | 목적 | 인프라 | 배포 빈도 | 접근 권한 |
|------|------|--------|---------|----------|
| **Local** | 개발자 로컬 | Docker Compose | 수시 | 개발자 |
| **Dev** | 개발 통합 | Docker Compose (서버) | PR 머지 시 | 개발팀 |
| **Staging** | QA 검증 | Kubernetes (단일 노드) | 일 1회 또는 수동 | 개발팀 + QA |
| **Production** | 실서비스 | Kubernetes (다중 노드) | 주 1~2회 (Phase 3+) | 운영팀 (승인 필요) |

### 1.2 환경별 리소스 스펙

| 환경 | 노드 | CPU/노드 | 메모리/노드 | 스토리지 |
|------|------|---------|-----------|---------|
| Local | 1 | 4 vCPU | 8GB | 50GB |
| Dev | 1 | 8 vCPU | 16GB | 100GB |
| Staging | 3 | 4 vCPU | 8GB | 200GB |
| Production (Phase 2) | 5 | 8 vCPU | 16GB | 1TB |
| Production (Phase 3) | 10 | 8 vCPU | 32GB | 5TB |
| Production (Phase 4) | 20+ | 16 vCPU | 64GB | 20TB+ |

---

## 2. Phase별 배포 전략

### 2.1 Phase 1 (MVP) — 현재

| 항목 | 전략 |
|------|------|
| **환경** | Local (Docker Compose) + Dev |
| **배포 방식** | `docker compose up -d` (수동) |
| **이미지 빌드** | 멀티스테이지 Dockerfile (Go 빌드 → alpine 실행) |
| **DB 마이그레이션** | SQL 초기화 스크립트 (infrastructure/database/init/) |
| **롤백** | `docker compose down && docker compose up -d` |
| **CI** | GitHub Actions (빌드 + 테스트 + 린트) |
| **CD** | 수동 (docker-compose) |

```
Developer → git push → GitHub Actions CI
                           ├── Rust: clippy + test + build
                           ├── Go: lint + test + build (18서비스)
                           ├── Docker: 이미지 빌드
                           └── E2E: 기본 플로우 테스트
```

### 2.2 Phase 2 (Core) — 목표

| 항목 | 전략 |
|------|------|
| **환경** | Dev + Staging (K8s) |
| **배포 방식** | Rolling Update (K8s Deployment) |
| **이미지 레지스트리** | GitHub Container Registry (ghcr.io) |
| **DB 마이그레이션** | golang-migrate (자동, 배포 전 Job) |
| **롤백** | `kubectl rollout undo deployment/<서비스명>` (< 2분) |
| **시크릿 관리** | Kubernetes Secrets + Sealed Secrets |
| **CI** | GitHub Actions (빌드 + 테스트 + 보안 스캔 + 이미지 푸시) |
| **CD** | GitHub Actions → Staging 자동 배포 |

```
Developer → git push → GitHub Actions CI
                           ├── 빌드/테스트/린트/보안스캔
                           ├── Docker 이미지 빌드 + ghcr.io 푸시
                           └── tag: v*.*.* 시
                                ├── Staging 자동 배포 (K8s Rolling)
                                ├── DB 마이그레이션 (Job)
                                ├── E2E 테스트 (Staging)
                                └── ✅ 통과 시 Production 승인 대기
```

**Rolling Update 설정:**
```yaml
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 0      # 무중단
      maxSurge: 1            # 1개씩 교체
  minReadySeconds: 10        # 10초 안정화 대기
```

### 2.3 Phase 3 (Advanced) — 프로덕션 출시

| 항목 | 전략 |
|------|------|
| **환경** | Dev + Staging + Production |
| **배포 방식** | Canary 배포 (5% → 25% → 50% → 100%) |
| **트래픽 관리** | Istio Service Mesh |
| **DB 마이그레이션** | 무중단 마이그레이션 (Schema 호환성 유지) |
| **롤백** | 자동 (에러율 > 1% 시 즉시 롤백) |
| **모니터링** | Canary 메트릭 비교 (베이스라인 vs 카나리) |
| **CI** | + 카오스 테스트, 성능 테스트 |
| **CD** | GitOps (ArgoCD) + 자동 Canary |

```
Release Tag → GitHub Actions CI
                 ├── 전체 테스트 스위트
                 └── Docker 이미지 빌드 + 푸시
                      ↓
ArgoCD 감지 → Canary 배포 시작
                 ├── Step 1: 5% 트래픽 → 카나리 Pod (10분 관찰)
                 │   └── 에러율 < 1%, 지연 < 200ms → 다음 단계
                 ├── Step 2: 25% 트래픽 (10분 관찰)
                 ├── Step 3: 50% 트래픽 (10분 관찰)
                 └── Step 4: 100% 트래픽 → 배포 완료
                      └── 실패 시: 자동 롤백 (< 30초)
```

**Canary 자동 롤백 조건:**
- 5xx 에러율 > 1%
- P99 응답시간 > 1초
- Pod 크래시 > 0

### 2.4 Phase 4 (Ecosystem) — 글로벌 확장

| 항목 | 전략 |
|------|------|
| **환경** | Dev + Staging + Production (Multi-Region) |
| **배포 방식** | Blue-Green 배포 |
| **트래픽 관리** | Istio + Global Load Balancer |
| **리전** | KR (서울), US (버지니아), EU (프랑크푸르트), JP (도쿄) |
| **DB** | CockroachDB 또는 PostgreSQL + Citus (멀티리전 샤딩) |
| **CDN** | CloudFront / Cloud CDN (정적 자산, Flutter 웹) |
| **롤백** | Blue-Green 전환 (< 10초) |

```
릴리스 → Blue 환경에 배포 (현재 Green이 서비스 중)
           ├── Blue에서 전체 테스트 실행
           ├── 스모크 테스트 PASS
           └── 트래픽 전환: Green → Blue (DNS/LB 전환)
                ├── 10분 모니터링
                └── 이상 시: Blue → Green 즉시 전환
```

---

## 3. DB 마이그레이션 전략

### 3.1 마이그레이션 도구

| Phase | 도구 | 방식 |
|-------|------|------|
| Phase 1 | SQL 초기화 스크립트 | Docker entrypoint |
| Phase 2+ | golang-migrate | K8s Job (배포 전 실행) |

### 3.2 무중단 마이그레이션 규칙 (Phase 3+)

1. **컬럼 추가**: 항상 `DEFAULT NULL` 또는 `DEFAULT 값`으로 추가
2. **컬럼 삭제**: 3단계 (미사용 코드 배포 → 컬럼 nullable → 컬럼 삭제)
3. **테이블 이름 변경**: 2단계 (새 테이블 생성 + 트리거 → 구 테이블 삭제)
4. **인덱스 생성**: `CREATE INDEX CONCURRENTLY` 사용
5. **대용량 데이터 변환**: 배치 처리 (1,000행씩)

### 3.3 마이그레이션 실행 순서

```
1. DB 마이그레이션 Job 실행 (K8s Job)
2. 마이그레이션 성공 확인
3. 새 버전 서비스 배포 (Rolling/Canary)
4. 구 버전 호환성 확인 (Backward Compatible)
5. 다음 릴리스에서 구 스키마 정리
```

---

## 4. Docker 이미지 관리

### 4.1 이미지 태깅 전략

| 태그 패턴 | 용도 | 예시 |
|----------|------|------|
| `latest` | Dev 환경 | `ghcr.io/manpasik/auth-service:latest` |
| `v{major}.{minor}.{patch}` | 릴리스 | `ghcr.io/manpasik/auth-service:v1.2.3` |
| `sha-{commit}` | CI 추적 | `ghcr.io/manpasik/auth-service:sha-abc1234` |
| `staging` | Staging 환경 | `ghcr.io/manpasik/auth-service:staging` |

### 4.2 이미지 최적화

| 항목 | 방법 |
|------|------|
| 베이스 이미지 | `alpine:3.19` (최소 크기) |
| 멀티스테이지 빌드 | Go 빌드 → alpine 복사 |
| 레이어 캐시 | `go.mod/go.sum` 먼저 COPY |
| 보안 스캔 | Trivy (매 빌드) |
| 목표 크기 | Go 서비스: < 30MB |

---

## 5. Kubernetes 매니페스트 관리

### 5.1 디렉토리 구조

```
infrastructure/kubernetes/
├── base/                    # 공통 매니페스트 (Kustomize base)
│   ├── config/
│   │   ├── configmap.yaml   # 환경 설정
│   │   └── secrets.yaml     # 시크릿 (Sealed Secrets)
│   ├── services/            # 서비스별 Deployment + Service
│   │   ├── auth-service.yaml
│   │   ├── ... (20개)
│   │   └── gateway.yaml
│   ├── ingress.yaml
│   └── kustomization.yaml
├── overlays/
│   ├── dev/                 # Dev 오버레이 (리소스 최소)
│   ├── staging/             # Staging 오버레이 (중간)
│   └── production/          # Production 오버레이 (HA, 리소스 확대)
```

### 5.2 환경별 오버레이

| 환경 | 레플리카 | CPU 요청 | 메모리 요청 | HPA |
|------|---------|---------|-----------|-----|
| Dev | 1 | 100m | 128Mi | 없음 |
| Staging | 2 | 250m | 256Mi | min:2 max:5 |
| Production | 3 | 500m | 512Mi | min:3 max:20 |

---

## 6. 모니터링 및 알림

### 6.1 배포 모니터링 대시보드 (Grafana)

| 대시보드 | 메트릭 |
|---------|--------|
| **배포 현황** | 현재 버전, 배포 시각, 롤백 횟수 |
| **서비스 상태** | Pod 수, CPU/메모리, 재시작 횟수 |
| **에러율** | 5xx/4xx 비율, gRPC 에러 코드 |
| **응답시간** | P50/P95/P99 히스토그램 |
| **처리량** | RPS, 활성 연결 수 |

### 6.2 알림 규칙

| 조건 | 심각도 | 알림 채널 | 응답 시간 |
|------|--------|---------|----------|
| Pod 크래시 반복 (3회/5분) | Critical | Slack + PagerDuty | 5분 |
| 에러율 > 5% | Critical | Slack + PagerDuty | 5분 |
| 에러율 > 1% | Warning | Slack | 15분 |
| 응답시간 P95 > 500ms | Warning | Slack | 30분 |
| CPU > 80% (5분 지속) | Warning | Slack | 30분 |
| 디스크 > 85% | Warning | Slack + Email | 1시간 |
| 인증서 만료 7일 전 | Info | Email | 24시간 |

---

## 7. 롤백 절차

### 7.1 자동 롤백 (Phase 3+ Canary)

```
1. Canary 메트릭 이상 감지 (에러율 > 1%)
2. ArgoCD/Istio 자동 롤백 트리거
3. 카나리 Pod 제거, 이전 버전 트래픽 100% 복구
4. Slack 알림: "자동 롤백 실행 - 서비스명, 사유, 이전 버전"
5. 개발팀 원인 분석 → 수정 후 재배포
```

### 7.2 수동 롤백

```bash
# Phase 2 (K8s Rolling)
kubectl rollout undo deployment/auth-service -n manpasik

# Phase 3 (ArgoCD)
argocd app rollback manpasik-auth-service

# 확인
kubectl rollout status deployment/auth-service -n manpasik
```

### 7.3 DB 롤백

```bash
# golang-migrate 다운 마이그레이션
migrate -source file://backend/migrations \
  -database "postgres://..." \
  down 1
```

---

## 8. Phase별 배포 체크리스트

### Phase 2 배포 체크리스트

- [ ] CI 전체 통과 (빌드, 테스트, 린트, 보안)
- [ ] Docker 이미지 빌드 성공 (18서비스)
- [ ] DB 마이그레이션 Staging 검증
- [ ] E2E 테스트 Staging 통과
- [ ] 성능 테스트 통과 (5,000 RPS)
- [ ] Staging 24시간 안정성 확인
- [ ] 운영팀 배포 승인
- [ ] Production 배포 실행
- [ ] 배포 후 30분 모니터링
- [ ] 배포 완료 확인 및 CHANGELOG 갱신

---

**참조**: .github/workflows/ci.yml, .github/workflows/cd.yml, infrastructure/kubernetes/, QUALITY_GATES.md
