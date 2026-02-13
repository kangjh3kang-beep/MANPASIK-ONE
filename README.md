# ManPaSik AI Ecosystem (만파식 AI 생태계)

> **철학**: 홍익인간(弘益人間) - 세상의 모든 파동을 분석해 인간을 이롭게 한다

[![License](https://img.shields.io/badge/license-Proprietary-red)]()
[![Flutter](https://img.shields.io/badge/Flutter-3.x-blue)]()
[![Rust](https://img.shields.io/badge/Rust-1.75+-orange)]()
[![Go](https://img.shields.io/badge/Go-1.22+-cyan)]()

## 개요

만파식(萬波息)은 차동측정 기반 범용 분석 시스템으로, 88~896차원 핑거프린트 벡터를 생성하여 
혈액, 타액, 체액, 소변, 수질, 식품, 공기질, 방사능, 라돈 등 측정 가능한 모든 항목을 
체계적으로 측정, 분석, 관리하는 통합 헬스케어 플랫폼입니다.

### 공통 개발 규칙 (모든 IDE)

**어떤 IDE에서 작업하든** 적용하는 공통 규칙: **[docs/COMMON_RULES.md](docs/COMMON_RULES.md)**  
- 단계 완료 시: **코드 리뷰 → 린트 → 테스트·빌드** 필수  
- 상세: [QUALITY_GATES.md](QUALITY_GATES.md)

## 핵심 기술

- **차동측정 엔진**: `S_det - α × S_ref` (99% 매트릭스 노이즈 제거)
- **896차원 핑거프린트**: 전자코(8ch) + 전자혀(8ch) 융합
- **엣지 AI**: TFLite 기반 오프라인 100% 동작
- **무제한 리더기 확장**: BLE/Wi-Fi Hub/Cloud Gateway 계층형 연결

## 프로젝트 구조

```
Manpasik/
├── docs/             # 문서
├── infrastructure/   # 인프라 (K8s, Terraform, Docker)
├── backend/          # Go 마이크로서비스 (30+)
├── rust-core/        # Rust 코어 엔진
├── ai-ml/            # AI/ML 파이프라인
├── frontend/         # Flutter 앱 + Next.js 웹
├── sdk/              # 개발자 SDK
└── tests/            # 테스트
```

## 빠른 시작

### ⚠️ 실행 위치

**다음 명령은 반드시 프로젝트 루트(`Manpasik/`)에서 실행하세요.** `backend/` 또는 다른 하위 폴더에서 실행하면 Makefile·경로를 찾지 못합니다.

| 목적 | 루트에서 실행 |
|------|----------------|
| `make proto` | ✅ `cd ~/Manpasik && make proto` |
| `make build-go` | ✅ `cd ~/Manpasik && make build-go` |
| Docker Compose | ✅ `cd ~/Manpasik && docker compose -f infrastructure/docker/docker-compose.dev.yml up -d` |
| E2E 테스트 | `cd ~/Manpasik/backend && go test -v ./tests/e2e/...` (backend가 모듈 루트) |

### 개발 환경 실행

**프로젝트 루트(Manpasik)에서** 실행하세요. `docker-compose.dev.yml`은 `infrastructure/docker/` 안에 있습니다.

```bash
cd ~/Manpasik
cd infrastructure/docker
docker compose -f docker-compose.dev.yml up -d
```

> WSL 2에서 `docker`/`docker-compose`를 찾을 수 없다면 [Docker Desktop WSL 통합](https://docs.docker.com/go/wsl2/)을 활성화하세요.

### 서비스 접속

| 서비스 | URL | 기본 계정 |
|--------|-----|-----------|
| Kong API Gateway | http://localhost:8000 | - |
| Kong Admin | http://localhost:8001 | - |
| Keycloak | http://localhost:8080 | admin / admin |
| Grafana | http://localhost:3000 | admin / admin |
| MinIO Console | http://localhost:9010 | manpasik / manpasik_dev_2026 |
| Milvus | localhost:19530 | - |
| PostgreSQL | localhost:5432 | manpasik / manpasik_dev_2026 |
| TimescaleDB | localhost:5433 | manpasik / manpasik_dev_2026 |
| Redis | localhost:6379 | - |
| Elasticsearch | localhost:9200 | - |
| MQTT | localhost:1883 | - |
| Kafka | localhost:19092 | - |

**Go gRPC 서비스 (Phase 1C)**  
| 서비스 | 주소 | 비고 |
|--------|------|------|
| auth-service | localhost:50051 | gRPC + health |
| user-service | localhost:50052 | gRPC + health |
| device-service | localhost:50053 | gRPC + health |
| measurement-service | localhost:50054 | gRPC + health |

## 기술 스택

### 프론트엔드
- Flutter 3.x (모바일/데스크톱)
- Next.js 14+ (웹 PWA)
- Riverpod (상태관리)

### 코어 엔진
- Rust no_std (차동측정)
- TFLite (엣지 AI)
- flutter_rust_bridge (FFI)

### 백엔드
- Go + gRPC (마이크로서비스)
- Kong (API Gateway)
- Keycloak (인증/인가)

### 데이터
- PostgreSQL 16 (메인 DB)
- TimescaleDB (시계열)
- Milvus (벡터 DB)
- Redis (캐시)
- Apache Kafka (이벤트 스트리밍)

### 인프라
- Kubernetes (오케스트레이션)
- Docker (컨테이너)
- GitHub Actions (CI/CD)

## 라이선스

Proprietary - All Rights Reserved

---

**Document Version**: 2.0.0  
**Created**: 2026-02-09  
**Based on**: MPK-ECO-PLAN-v1.0-20260208-FINAL
