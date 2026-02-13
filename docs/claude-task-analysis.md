# Claude 업무 분석 및 보완 계획

> **작성일**: 2026-02-09
> **목적**: Claude에게 할당된 전체 업무를 식별하고, 완료/보완/미착수 상태를 체계적으로 정리
> **근거 문서**: `AI_COLLABORATION.md`, `medical-compliance-spec.md`, `security-architecture-spec.md`

---

## 1. Claude 역할 정의 (AI_COLLABORATION.md 기준)

```
Claude - 심층 분석 & 문서화

담당 영역 6가지:
  ① 의료 규정 분석 (HIPAA, GDPR, MFDS, FDA)
  ② 보안 아키텍처 설계 및 검토
  ③ API 문서 작성 (OpenAPI, 사용자 가이드)
  ④ 코드 리뷰 및 취약점 분석
  ⑤ AI/ML 모델 설계 (연합학습, 엣지 AI)
  ⑥ 임상시험 설계 문서

명시적 산출물 3종:
  1. 의료기기 규제 준수 체크리스트 작성
  2. 데이터 보호 정책 문서
  3. AI 모델 학습 전략 설계
```

### Phase별 Claude 분담 (Week 단위)

| 주차 | 역할 | 상태 |
|------|------|------|
| Week 1-4 | 규정 분석 시작 | ✅ 완료 (체크리스트 + 안전등급 + STRIDE + 데이터보호 정책) |
| Week 5-8 | 보안 설계 | ✅ 부분 완료 (STRIDE + 데이터보호, 추가 보완 필요) |
| Week 9-12 | API 문서 | ❌ 미시작 |
| Week 13-16 | 코드 리뷰 | ❌ 미시작 |

---

## 2. 산출물 완료/미완료 전수 대조

### 2.1 medical-compliance-spec.md 요청 산출물 (5종)

| # | 요청 산출물 | 상태 | 산출물 경로 | 보완 필요 |
|---|-----------|------|-----------|---------|
| MC-1 | **규제 준수 체크리스트 (국가별)** | ✅ 완료 | `docs/compliance/regulatory-compliance-checklist.md` | 🔄 완료 항목 상태 갱신, 기술문서 목차 누락 |
| MC-2 | **데이터 보호 정책 초안** | ✅ 완료 | `docs/compliance/data-protection-policy.md` | 🔄 DPIA 템플릿, BAA 템플릿, ROPA 추가 |
| MC-3 | **소프트웨어 안전 등급 판정** | ✅ 완료 | `docs/compliance/software-safety-classification.md` | 🔴 서브시스템별 등급 할당 누락 (IEC 62304 5.3.3) |
| MC-4 | **위험관리 파일 목차** | ✅ 완료 | `docs/compliance/software-safety-classification.md` 내 섹션 4 | 🟡 개별 파일 실체 미생성 (목차만 존재) |
| MC-5 | **사이버보안 요구사항 문서** | ✅ 완료 | `docs/security/stride-threat-model.md` (STRIDE + SBOM + IRP) | 🔄 DFD 다이어그램, RTO/RPO 수치 미확정 |

### 2.2 security-architecture-spec.md 요청 산출물 (5종)

| # | 요청 산출물 | 상태 | 산출물 위치 | 보완 필요 |
|---|-----------|------|-----------|---------|
| SA-1 | **보안 아키텍처 다이어그램** | 🟡 부분 | `stride-threat-model.md` 섹션 1.1 (간이 다이어그램) | 🔴 전체 보안 계층 다이어그램 미작성 |
| SA-2 | **위협 모델링 결과 (STRIDE)** | ✅ 완료 | `docs/security/stride-threat-model.md` | 🔄 DFD 기반 분석 보완 |
| SA-3 | **암호화 정책 문서** | ✅ 완료 | `stride-threat-model.md` 섹션 4.3 + `data-protection-policy.md` 섹션 4.1 | 🔄 양자내성 암호화, 키 복구 프로세스 보완 |
| SA-4 | **접근제어 매트릭스** | ✅ 완료 | `stride-threat-model.md` 섹션 4.2 | 🟡 ABAC 정책 상세 미정의 |
| SA-5 | **침해사고 대응 계획** | ✅ 완료 | `stride-threat-model.md` 섹션 5 | 🔄 RTO/RPO 수치, 모의훈련 계획 보완 |

### 2.3 AI_COLLABORATION.md 명시적 산출물 (3종)

| # | 요청 산출물 | 상태 | 비고 |
|---|-----------|------|------|
| AC-1 | 의료기기 규제 준수 체크리스트 | ✅ 완료 | MC-1과 동일 |
| AC-2 | 데이터 보호 정책 문서 | ✅ 완료 | MC-2와 동일 |
| AC-3 | **AI 모델 학습 전략 설계** | 🔴 **미착수** | `ml-model-design-spec.md`도 부재, 산출물도 미작성 |

### 2.4 AI_COLLABORATION.md 역할 기반 미착수 업무 (6종)

| # | 업무 | 상태 | Phase | 의존성 |
|---|------|------|-------|--------|
| R-1 | **API 문서 작성 (OpenAPI)** | ❌ 미착수 | Week 9-12 | Go 서비스 구현 후 (ChatGPT 의존) |
| R-2 | **코드 리뷰 및 취약점 분석** | ❌ 미착수 | Week 13-16 | 코드 구현 진행 후 |
| R-3 | **AI/ML 모델 설계** | 🔴 미착수 | Week 5-8 | 독립 (즉시 가능) |
| R-4 | **임상시험 설계 문서** | ❌ 미착수 | Phase 2+ | 규정 분석 완료 후 |
| R-5 | **사용자 가이드** | ❌ 미착수 | Phase 2+ | Flutter 앱 구현 후 |
| R-6 | **설계 검토 (Design Review)** | ❌ 미착수 | 지속 | 코드 변경 시 마다 |

---

## 3. 기존 산출물 보완 필요 상세

### 3.1 🔴 Critical: 서브시스템별 안전 등급 할당 (IEC 62304 5.3.3)

**현재 상태**: 전체 시스템을 Class B로 판정했으나, IEC 62304 5.3.3은 소프트웨어 아키텍처의 각 소프트웨어 아이템(서브시스템)에 안전 등급을 할당할 것을 요구합니다.

**필요한 보완**: 각 모듈/서비스에 개별 안전 등급 부여

### 3.2 🔴 Critical: 기술문서(Technical File) 목차 누락

**현재 상태**: `medical-compliance-spec.md`에서 "기술문서(Technical File) 목차"를 산출물로 명시적 요청했으나 미작성.

**필요한 보완**: FDA 510(k) 및 CE-IVDR Annex II/III 양식의 기술문서 목차

### 3.3 🔴 Critical: AI/ML 모델 설계 전략 미착수

**현재 상태**: `AI_COLLABORATION.md`에서 Claude 산출물 #3으로 명시. `ml-model-design-spec.md` 파일 자체도 부재.

**필요한 보완**: FDA의 AI/ML SaMD 가이던스를 반영한 모델 설계 전략

### 3.4 🟡 High: V&V(검증/확인) 마스터 플랜 미작성

**현재 상태**: `medical-compliance-spec.md`에서 "검증/확인 활동 정의"를 요청. IEC 62304 Class B에서 통합 시험(5.6), 시스템 시험(5.7)은 필수.

### 3.5 🟡 High: 검증 보고서 오래된 상태 반영

**현재 상태**: `system-plan-verification.md`가 산출물 작성 전에 작성되어 "의료규정 분석 0%, 보안 아키텍처 0%"로 표시됨. 현재 실제로는 ~40% 진행.

---

## 4. Claude 전체 업무 로드맵 (우선순위 정렬)

### Phase 1A: 즉시 보완 (현재 세션)

| 순서 | 업무 | 유형 | 근거 |
|------|------|------|------|
| 1 | 서브시스템별 안전 등급 할당 | 기존 보완 | IEC 62304 5.3.3 필수, 개발 프로세스 결정 |
| 2 | 기술문서(Technical File) 목차 | 신규 | 스펙 요청 미이행, 인허가 구조 필수 |
| 3 | AI/ML 모델 설계 전략 | 신규 | AI_COLLABORATION 지정 산출물 #3 누락 |
| 4 | 검증 보고서 업데이트 | 기존 갱신 | 정확한 현황 반영 |

### Phase 1B: 2주 내 (다음 세션)

| 순서 | 업무 | 유형 | 근거 |
|------|------|------|------|
| 5 | ISO 14971 위험관리 계획서 | 신규 | 인허가 핵심 문서, 목차만 존재 |
| 6 | V&V 마스터 플랜 | 신규 | IEC 62304 Class B 필수 |
| 7 | DPIA 템플릿 | 보완 | GDPR Art.35 필수 |
| 8 | Predicate Device 조사 (FDA 510(k)) | 신규 | 인허가 경로 결정 |
| 9 | FHIR R4 통합 설계 | 신규 | 데이터 호환성 표준 |

### Phase 1C: 1개월 내

| 순서 | 업무 | 유형 | 근거 |
|------|------|------|------|
| 10 | IEC 62366-1 사용적합성 엔지니어링 계획 | 신규 | 5개국 공통 필수 표준 |
| 11 | Rust 코어 코드 리뷰 (보안 관점) | R-2 | 구현 코드 존재 |
| 12 | SBOM 첫 생성 가이드 | 보완 | FDA Cybersecurity 필수 |
| 13 | 시판 후 감시(PMS) 계획 | 신규 | IVDR Art.78, MFDS 필수 |
| 14 | 임상시험 프로토콜 초안 | R-4 | MFDS/CE-IVDR 필수 |

### Phase 1D: Phase 1 완료까지

| 순서 | 업무 | 유형 | 근거 |
|------|------|------|------|
| 15 | OpenAPI 문서 작성 | R-1 | Go 서비스 구현 후 |
| 16 | QMS 매뉴얼 초안 | 신규 | ISO 13485 기반 |
| 17 | 형상관리(SCM) 절차서 | 신규 | IEC 62304 Clause 8 |
| 18 | 문제 해결 프로세스 절차서 | 신규 | IEC 62304 Clause 9 |
| 19 | CAPA 절차서 | 신규 | ISO 13485 8.5 |

---

## 5. 업무 통계 요약

| 구분 | 건수 | 비율 |
|------|------|------|
| ✅ 완료 (품질 양호) | 5 | 26% |
| 🔄 완료 but 보완 필요 | 5 | 26% |
| 🔴 미착수 (Critical) | 3 | 16% |
| ❌ 미착수 (Phase 1) | 6 | 32% |
| **전체** | **19** | **100%** |

### Claude 업무 부하 추정

| Phase | 업무 수 | 예상 소요 |
|-------|--------|---------|
| 1A (즉시) | 4건 | 현재 세션 |
| 1B (2주) | 5건 | 2-3 세션 |
| 1C (1개월) | 5건 | 3-4 세션 |
| 1D (Phase 1 완료) | 5건 | 2-3 세션 |
| **합계** | **19건** | **~10-12 세션** |

---

## 6. 핵심 의존성 맵

```
Claude 업무 의존성 구조:

[독립 - 즉시 가능]
  ├── 서브시스템 안전 등급 ──→ V&V 플랜 ──→ 테스트 전략
  ├── 기술문서 목차 ──→ 각국 인허가 서류 패키지
  ├── AI/ML 모델 전략 ──→ FDA AI/ML SaMD 제출
  ├── 위험관리 계획서 ──→ FMEA ──→ 위험관리 보고서
  └── Predicate Device 조사 ──→ 510(k) Substantial Equivalence

[ChatGPT 의존]
  ├── Go 서비스 구현 후 ──→ OpenAPI 문서
  ├── Flutter 앱 구현 후 ──→ 사용적합성 테스트
  └── DB 스키마 확정 후 ──→ 감사 추적 설계 검토

[Antigravity 의존]
  ├── crypto 모듈 AES-256-GCM 구현 후 ──→ 암호화 구현 검증
  ├── tenant_id 결정 후 ──→ 멀티테넌시 보안 리뷰
  └── BLE Secure Connections 구현 후 ──→ BLE 보안 검증
```

---

**Document Version**: 1.0.0
**작성일**: 2026-02-09
**작성자**: Claude (Security & Architecture Agent)
