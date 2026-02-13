# ManPaSik 공통 개발 규칙 (모든 IDE·환경)

> **적용 범위**: Cursor, VS Code, JetBrains, Vim, 터미널 등 **어떤 IDE나 환경에서 작업하든** 동일하게 적용하는 프로젝트 공통 규칙입니다.  
> **기준 문서**: 이 파일이 단일 기준( single source of truth )이며, 상세 절차는 [QUALITY_GATES.md](../QUALITY_GATES.md)를 따릅니다.

---

## 1. 단계 완료 시 필수 3단계

기능 추가, 서비스 추가, Stage/Phase 완료 시 **완료 선언 전에** 아래 순서로 반드시 수행합니다.

| 순서 | 항목 | 내용 |
|------|------|------|
| **1** | **코드 리뷰** | 변경/추가된 코드 자기 점검 — 보안(입력 검증·인젝션·비밀 하드코딩), 프로젝트 패턴(에러 처리·API 일관성·nil/panic), 의존성 정합성(import·타입·스텁 일치) |
| **2** | **린트** | 해당 언어 린터 실행, **에러 0개** (Go: `golangci-lint run`, Rust: `cargo clippy`, Dart: `dart analyze` 등) |
| **3** | **테스트·빌드** | 해당 범위 테스트 및 빌드 실행 (Go: `go mod tidy && go test ./...` 및 `go build ./...` 등) |

- 세 단계 모두 통과해야 작업을 **완료**로 간주합니다.
- 상세 명령어·언어별 절차: [QUALITY_GATES.md §2](../QUALITY_GATES.md#2-level-1-매-작업-즉시-검증-every-change)

---

## 2. 핵심 원칙 요약

- **보안 우선**: 입력 검증·ORM 사용·JWT/RBAC·비밀/스택 노출 금지. ([QUALITY_GATES](../QUALITY_GATES.md), 프로젝트 규칙 참조)
- **테스트 필수**: "No Code Without Tests" — 기능 구현 전/완료 시 테스트 작성·실행, 커버리지 80%+ 목표.
- **우회·미루기 금지**: 구현을 "나중에" 또는 "우회"로 미루지 않음. 부득이한 경우 KNOWN_ISSUES 등록 및 해결 조건·시한 명시.

---

## 3. 참고 문서 위치

| 문서 | 용도 |
|------|------|
| [QUALITY_GATES.md](../QUALITY_GATES.md) | 단계별 품질 검증·Level 1/2/3·코드 리뷰→린트→테스트 상세 |
| [CONTEXT.md](../CONTEXT.md) | 프로젝트 현황·Phase·다음 단계 |
| [CHANGELOG.md](../CHANGELOG.md) | 작업 기록·검증 결과 기록 |
| [KNOWN_ISSUES.md](../KNOWN_ISSUES.md) | 알려진 이슈·우회 사항 |
| `.cursor/rules/` | **Cursor IDE 전용** — Cursor AI가 자동으로 참조하는 규칙 (다른 IDE에서는 적용되지 않음) |

---

**문서 버전**: 1.0  
**최종 업데이트**: 2026-02-10
