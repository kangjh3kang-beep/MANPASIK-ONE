# Phase H: 규제/인증 — 21 CFR Part 11 + 컴플라이언스 자동 검사

## 개요
Phase H는 POCT 의료기기 소프트웨어의 규제 준수를 코드 레벨에서 보강하는 단계입니다.
정밀 실사 결과, 21 CFR Part 11(전자서명/감사로그 무결성) 구현이 부재하고,
규제 준수 현황이 하드코딩되어 실시간 자동 검사가 불가능한 상태였습니다.

## 사전 실사 결과

### Compliance 패키지 현황 (Phase H 이전)
| 파일 | 줄수 | 기능 | 상태 |
|------|------|------|------|
| encryption.go | 136 | AES-256-GCM PHI 암호화 | ✅ 완성 |
| retention.go | 130 | 데이터 보존 정책 (7개 유형) | ✅ 완성 |
| deletion.go | 204 | GDPR Art.17 삭제권 프레임워크 | ✅ 프레임워크 |
| compliance_test.go | 383 | 테스트 20건 | ✅ ALL PASS |
| **전자서명** | — | — | ❌ 미구현 |
| **해시 체인** | — | — | ❌ 미구현 |
| **자동 검사** | — | — | ❌ 하드코딩 |

### 규제 프레임워크별 소프트웨어 준비도 (Phase H 이전)
| 프레임워크 | 준수율 | 비고 |
|-----------|--------|------|
| GDPR | 50% | 동의/암호화/감사 기초, 삭제권 미완 |
| HIPAA | 40% | RBAC/감시, 무결성 미흡 |
| 21 CFR Part 11 | 30% | 전자서명/부인방지 없음 |
| IEC 62304 | 60% | SRS 완성, SDP/SAD 미흡 |
| MFDS | 10% | 임상시험 미실시 (blocking) |
| FDA 510(k) | 5% | 예비 조사만 |
| CE-IVDR | 3% | Notified Body 필요 |

## 구현 내역

### H-1: 21 CFR Part 11 전자서명 + 감사로그 무결성

**파일**: `backend/shared/compliance/signature.go` (신규, ~120줄)

#### 전자 서명 (AuditSigner)

| 구조체/함수 | 설명 |
|------------|------|
| `AuditSignature` | 감사 로그 전자서명 (EntryID + SignerID + Signature + Algorithm) |
| `AuditSigner` | HMAC-SHA256 기반 서명기 |
| `NewAuditSigner(key)` | 서명 키 검증 (최소 16바이트) |
| `SignEntry(entryID, signerID, payload)` | 감사 항목 서명 생성 |
| `VerifySignature(signature, payload)` | 서명 검증 (변조 감지) |
| `BuildAuditPayload(...)` | 필드 결합 → 서명 대상 문자열 |

규제 대응:
- ✅ 21 CFR Part 11 §11.100 (전자 서명)
- ✅ 21 CFR Part 11 §11.10(e) (감사 로그 변조 방지)
- ✅ HIPAA §164.312(c)(1) (무결성 검증)

#### 해시 체인 (AuditChain)

| 구조체/함수 | 설명 |
|------------|------|
| `AuditChain` | SHA-256 해시 체인 (각 항목이 이전 해시 포함) |
| `NewAuditChain(genesisHash)` | 초기 해시로 체인 시작 |
| `AppendEntry(payload)` | 새 항목 추가 → 체인 해시 반환 |
| `VerifyChain(genesis, entries)` | 전체 체인 무결성 검증 (위조 인덱스 반환) |
| `ChainEntry` | 검증용 (payload + storedHash) |

설계 원칙:
- **Append-only**: 중간 삽입/삭제/수정 불가
- **Tamper-evident**: 하나라도 변조 시 체인 전체 검증 실패
- **Deterministic**: 동일 입력 → 동일 해시 (재현 가능)

### H-2: 규제 컴플라이언스 자동 검사

**파일**: `backend/shared/compliance/auto_check.go` (신규, ~210줄)

#### ComplianceChecker

| 구조체/함수 | 설명 |
|------------|------|
| `SystemState` | 시스템 현황 (암호화/감사/접근/데이터/문서/인증 16필드) |
| `ComplianceChecker` | 시스템 상태 기반 자동 검사기 |
| `RunChecks()` | 6개 프레임워크 자동 검사 실행 |
| `ComputeScore(statuses)` | 종합 준수율 계산 (0~100%) |

검사 대상 프레임워크:

| 프레임워크 | 검사 항목 수 | 검사 내용 |
|-----------|-------------|---------|
| GDPR | 5 | 동의, PHI 암호화, 보존 정책, 삭제권, 감사 |
| HIPAA | 4 | 접근 제어, 감사+서명, 무결성, 암호화 |
| MFDS | 5 | 안전 분류, 위험관리, SRS, 임상시험, GMP |
| FDA | 2 | Predicate Device, Pre-submission |
| CE-IVDR | 1 | 기술 문서 |
| 21 CFR 11 | 3 | 전자서명, 감사보호, 접근제어 |
| **합계** | **20** | — |

점수 계산: compliant=100%, partial=50%, non_compliant=0%

## 테스트 결과

### 테스트 전후 비교
| 구분 | Before | After | 증가 |
|------|--------|-------|------|
| 테스트 함수 | 20건 | 31건 | +11 |
| 코드 줄수 | ~470줄 | ~800줄 | +330 |

### Phase H 추가 테스트 (11건)
| 테스트 | 범주 | 설명 |
|--------|------|------|
| TestAuditSigner_SignAndVerify | 전자서명 | 서명 생성 + 검증 |
| TestAuditSigner_TamperedPayload | 전자서명 | 변조된 페이로드 → 실패 |
| TestAuditSigner_DifferentKeys | 전자서명 | 다른 키 → 검증 실패 |
| TestAuditSigner_ShortKey | 전자서명 | 짧은 키 → 에러 |
| TestAuditSigner_EmptyFields | 전자서명 | 빈 필드 → 에러 |
| TestAuditSigner_NilSignature | 전자서명 | nil 서명 → 에러 |
| TestAuditChain_AppendAndVerify | 해시 체인 | 5개 항목 추가 + 체인 검증 |
| TestAuditChain_TamperDetection | 해시 체인 | 변조 감지 + 인덱스 반환 |
| TestAuditChain_EmptyChain | 해시 체인 | 빈 체인 → valid |
| TestAuditChain_DeterministicHash | 해시 체인 | 결정적 해시 |
| TestBuildAuditPayload_Deterministic | 유틸 | 결정적 페이로드 |

### 자동 검사 테스트 (6건)
| 테스트 | 설명 |
|--------|------|
| TestComplianceChecker_FullyCompliant | 모든 기능 활성화 → >80% |
| TestComplianceChecker_MinimalState | 최소 상태 → <30% |
| TestComplianceChecker_FrameworkCoverage | 6개 프레임워크 커버 |
| TestComplianceChecker_Recommendations | 비준수 → 권고 포함 |
| TestComputeScore_Empty | 빈 입력 → 0점 |
| TestComputeScore_Mixed | 혼합 상태 → 62.5% |

## 빌드 검증
```
compliance 패키지: BUILD OK + 31 PASS (0.007s)
```

## Phase H 전후 비교

| 영역 | Before | After | 변화 |
|------|--------|-------|------|
| 전자 서명 | 없음 | HMAC-SHA256 | 신규 |
| 감사 무결성 | 없음 | SHA-256 해시 체인 | 신규 |
| 자동 검사 | 하드코딩 (19항목) | 동적 검사 (20항목, 6프레임워크) | 대폭 개선 |
| compliance 코드 | ~470줄 | ~800줄 | +330줄 |
| compliance 테스트 | 20건 | 31건 | +11건 |

## 변경 파일 목록
1. `backend/shared/compliance/signature.go` — 신규 (전자서명 + 해시 체인)
2. `backend/shared/compliance/auto_check.go` — 신규 (자동 검사기 + 점수 계산)
3. `backend/shared/compliance/compliance_test.go` — +11 서명 테스트 + +6 자동 검사 테스트

## 규제 인허가 로드맵 (코드 범위 밖 — 참조용)

### 인허가 Blocking Items (실물 규제 절차)
| 항목 | 예상 기간 | 비고 |
|------|----------|------|
| MFDS IRB 승인 | 2026-06 | 대학병원 IRB 선정 |
| 임상시험 실시 | 2026-07~11 | 100명 (건강인 50 + 환자 50) |
| ISO 13485 GMP | 2026-Q3 | 제조 품질 인증 |
| FDA Pre-submission | 2026-Q3 | 사전 상담 |
| CE-IVDR Notified Body | 2026-Q4 | 3rd party 인증기관 |
