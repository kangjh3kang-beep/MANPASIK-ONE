# Phase G: 운영/배포 — CI/CD + K8s 보안 + 모니터링

## 개요
Phase G는 운영/배포 인프라의 GAP을 보강하는 단계입니다.
정밀 실사 결과, CI 파이프라인이 Node.js만 검증하고 Go/Flutter/Rust를 누락하고 있었으며,
K8s에 NetworkPolicy가 없고, Prometheus에 인프라 모니터링이 부재한 상태였습니다.

## 사전 실사 결과

### CI/CD 현황 (Phase G 이전)
| 구분 | 상태 | 비고 |
|------|------|------|
| ci.yml | Node.js만 검증 (41줄) | Go/Flutter/Rust 누락 |
| cd.yml | 완성 (534줄) | 23 서비스 Docker + 3단계 배포 + 롤백 |
| Dockerfile | 완성 | 멀티스테이지 Go 빌드 (golang:1.24-alpine → scratch) |

### K8s 현황 (Phase G 이전)
| 구분 | 상태 | 비고 |
|------|------|------|
| base/ | 32 서비스 YAML | Deployment + Service |
| overlays/ | dev/staging/production | HPA(9개) + PDB |
| Ingress | 기본 설정 | CORS + SSL + 프록시 타임아웃 |
| NetworkPolicy | **없음** | 서비스 간 격리 미적용 |
| Rate Limiting | **없음** | DDoS/남용 방어 없음 |

### 모니터링 현황 (Phase G 이전)
| 구분 | 상태 | 비고 |
|------|------|------|
| Prometheus | 35 서비스 스크래핑 | 15s 인터벌 |
| Alert Rules | 12개 규칙 | 서비스/레이턴시/리소스/DB |
| kube-state-metrics | **없음** | 클러스터 상태 미수집 |
| node-exporter | **없음** | 노드 시스템 메트릭 미수집 |
| 인프라 exporter | **없음** | PostgreSQL/Redis/Kafka 미수집 |

## 구현 내역

### G-1: CI 파이프라인 확장

**파일**: `.github/workflows/ci.yml` (41줄 → 130줄)

기존 Node.js 웹앱 Job에 3개 Job을 추가하여 전체 스택을 검증:

| Job | 스택 | 검증 항목 |
|-----|------|-----------|
| `web-app` | Node.js/pnpm | type check + lint + test + build |
| `go-backend` | Go 1.22 | go vet + build all 35 services + go test |
| `flutter-app` | Flutter 3.32 | flutter analyze + flutter test |
| `rust-core` | Rust stable | cargo build + cargo clippy + cargo test |

주요 설계 결정:
- Go: `backend/` 디렉토리 기준, `go.sum` 캐시
- Flutter: `subosiatech/flutter-action@v2` + stable 채널
- Rust: `dtolnay/rust-toolchain@stable` + `actions/cache@v4` (Cargo 캐시)
- 4개 Job은 **병렬 실행** (의존 관계 없음)

### G-2: K8s 보안 강화

#### NetworkPolicy (신규)
**파일**: `infrastructure/kubernetes/base/network-policy.yaml`

| 정책 | 설명 |
|------|------|
| `default-deny-ingress` | 네임스페이스 기본 인바운드 거부 |
| `allow-ingress-to-gateway` | Ingress Controller → Gateway (8080) |
| `allow-gateway-to-grpc` | Gateway → gRPC 서비스 (50051) |
| `allow-prometheus-scrape` | monitoring NS → 모든 Pod (9100) |
| `allow-inter-service-grpc` | 서비스 간 gRPC 통신 (50051) |
| `allow-service-to-postgres` | 서비스 → PostgreSQL (5432) |
| `allow-service-to-kafka` | 서비스 → Kafka (9092) |

#### Ingress Rate Limiting (추가)
**파일**: `infrastructure/kubernetes/base/ingress.yaml`

| 설정 | 값 | 설명 |
|------|-----|------|
| `limit-rps` | 30 | IP당 초당 최대 요청 |
| `limit-burst-multiplier` | 3 | 버스트 허용 배수 (90 rps 순간 허용) |
| `limit-connections` | 10 | IP당 동시 연결 제한 |
| `enable-modsecurity` | true | WAF 기본 보호 (XSS/SQL Injection) |

#### kustomization.yaml 갱신
- `network-policy.yaml` 리소스 추가

### G-3: 모니터링 개선

#### Prometheus 스크래핑 추가 (6개)
**파일**: `infrastructure/monitoring/prometheus/prometheus.yml`

| Job | 타겟 | 설명 |
|-----|------|------|
| `kube-state-metrics` | monitoring:8080 | K8s 클러스터 상태 |
| `node-exporter` | kubernetes_sd node | 노드 시스템 메트릭 |
| `postgresql` | postgres-exporter:9187 | DB 메트릭 |
| `redis` | redis-exporter:9121 | 캐시 메트릭 |
| `kafka` | kafka-exporter:9308 | 메시지 큐 메트릭 |
| `nginx-ingress` | ingress-nginx:10254 | Ingress 메트릭 |

- `external_labels` 추가: cluster=manpasik, environment=production

#### 알림 규칙 추가 (6개)
**파일**: `infrastructure/monitoring/prometheus/alert_rules.yml`

| 알림 | 조건 | 심각도 |
|------|------|--------|
| `NodeHighCPU` | 노드 CPU > 85% (10분) | warning |
| `NodeHighMemory` | 노드 메모리 > 90% (5분) | critical |
| `NodeDiskAlmostFull` | 디스크 사용 > 85% (10분) | warning |
| `KafkaConsumerLag` | 소비자 지연 > 10,000건 (5분) | warning |
| `IngressHighErrorRate` | Ingress 5xx > 5% (5분) | warning |
| `IngressHighLatency` | Ingress p99 > 5초 (5분) | warning |

## Phase G 전후 비교

| 영역 | Before | After | 변화 |
|------|--------|-------|------|
| CI Jobs | 1 (Node.js) | 4 (Node/Go/Flutter/Rust) | +3 |
| CI 줄 수 | 41줄 | 130줄 | +89 |
| NetworkPolicy | 0개 | 7개 | +7 |
| Ingress 보안 | CORS만 | CORS + Rate Limit + WAF | 대폭 강화 |
| Prometheus Jobs | 35개 | 41개 | +6 |
| Alert Rules | 12개 | 18개 | +6 |
| 모니터링 범위 | 앱 서비스만 | 앱 + 인프라 (노드/DB/Redis/Kafka/Ingress) | 전체 |

## 변경 파일 목록
1. `.github/workflows/ci.yml` — Go/Flutter/Rust CI Job 추가
2. `infrastructure/kubernetes/base/network-policy.yaml` — 신규 7개 NetworkPolicy
3. `infrastructure/kubernetes/base/ingress.yaml` — Rate Limit + WAF 어노테이션 추가
4. `infrastructure/kubernetes/base/kustomization.yaml` — network-policy.yaml 리소스 등록
5. `infrastructure/monitoring/prometheus/prometheus.yml` — 6개 인프라 스크래핑 + external_labels
6. `infrastructure/monitoring/prometheus/alert_rules.yml` — 6개 인프라 알림 규칙 추가
