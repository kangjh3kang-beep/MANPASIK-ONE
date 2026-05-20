# Security Release Gates

## 목적

dev compose의 기본 계정과 비활성 보안 설정이 production profile로 번지는 것을 CI에서 차단한다. 이 게이트는 FDA 사이버보안/QMS 제출 준비의 최소 release hygiene으로, secret management와 production hardening의 초기 방어선이다.

## 변경 범위

- `scripts/security_release_gate.sh`
  - production 후보 파일에서 `dev-secret-change-in-production`, `manpasik_dev`, `minioadmin`, 기본 Keycloak admin password, Elasticsearch security disabled, `DB_SSLMODE=disable` 패턴을 거부한다.
  - Kubernetes base config의 `MAX_FINGERPRINT_DIMENSION=1792`, `DEFAULT_ALPHA=0.98`을 검증한다.
  - production overlay의 `DB_SSLMODE=require`, `TENANCY_ENFORCED=true`를 검증한다.
- `.github/workflows/ci.yml`
  - `ssot-governance` job에서 SSOT constants check 이후 security release gate를 실행한다.
- `infrastructure/kubernetes/base/config/configmap.yaml`
  - runtime config 상수를 최신 SSOT 기준인 `1792`, `0.98`로 정렬했다.

## 품질 게이트

- `bash scripts/security_release_gate.sh`: PASS

## 잔여 리스크

- 이 스크립트는 production profile로 새는 가장 위험한 문자열 패턴만 막는다. 실제 운영 배포 전에는 Kubernetes Secret/External Secrets/Vault 연동 검증, image signing, SBOM, vulnerability SLA, mTLS policy 검증을 추가해야 한다.
- root `docker-compose.yml`과 `infrastructure/docker/docker-compose.dev.yml`의 dev 기본값은 현재 개발 전용으로 허용한다. production compose 파일이 추가되면 반드시 이 게이트의 검사 대상에 포함해야 한다.
