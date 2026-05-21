# shared/fhir-adapter

이 디렉토리는 placeholder 입니다. **실제 FHIR 어댑터 구현은 [`backend/shared/medical/fhir/`](../medical/fhir/) 로 이전되었습니다.**

## 이전 경로

| 구버전 (이 디렉토리) | 신버전 |
|--------------------|--------|
| `backend/shared/fhir-adapter/` | [`backend/shared/medical/fhir/`](../medical/fhir/) |

## 구현 내용 (medical/fhir/)

- **R4 ↔ R5 Bridge** (`bridge.go`) — FHIR R4 와 R5 리소스 양방향 변환
- **Observation / DiagnosticReport / DocumentReference** — 측정 결과 → FHIR 리소스 매핑
- **Bundle 트랜잭션** — 다중 리소스 일괄 처리

상세 사용 예는 [`bridge_test.go`](../medical/fhir/bridge_test.go) 참고.

## 디렉토리 유지 사유

git history 보존 및 외부 빌드 도구의 캐시 경로 안정성을 위해 디렉토리만 유지하고, 코드는 medical 도메인 모듈로 통합되었습니다.
