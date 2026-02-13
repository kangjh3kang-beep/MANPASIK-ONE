# 만파식 생태계 기획·계획 문서 마스터 인덱스

**문서번호**: MPK-INDEX-v1.0  
**최종 갱신**: 2026-02-14  
**목적**: 그동안 수립한 모든 기획·설계·명세·검증 문서를 계층별로 정리하고, 상호 참조 관계를 명시하여 어떤 문서에서든 전체 체계를 파악할 수 있게 한다.

---

## 1. 문서 계층 구조

```text
[L0] 원본 기획
 └── MPK-ECO-PLAN-v1.1-COMPLETE.md ················· 원본 기획안 완성본 (베이스라인)
     ├── original-detail-annex.md ··················· 원본 표·수치·시나리오 보조 문서
     └── plan-verification-report.md ··············· 원본 대비 검증 보고서

[L1] 종합 기획서
 └── COMPREHENSIVE-SYSTEM-BLUEPRINT-v3.0.md ········ 종합 시스템 상세 기획서 (Part 1~11)
     ├── COMPREHENSIVE-IMPLEMENTATION-MASTERPLAN-v2.0.md · 통합 세부 구현 마스터플랜
     └── terminology-and-tier-mapping.md ··········· 용어·티어·카트리지 접근 통일

[L2] 통합 구현·구축 계획
 └── FINAL-MASTER-IMPLEMENTATION-PLAN.md ··········· 최종 통합 구현·구축 계획 (I~V)
     └── ★ DETAILED-IMPLEMENTATION-PLAN.md ·········· 주차별 상세 구현 계획 (본 문서와 동시 생성)

[L3] 기능별 상세 명세서
 ├── AI-ASSISTANT-MASTER-SPEC.md ··················· AI 비서·주치의 통합 기획
 ├── MEASUREMENT-ANALYSIS-AI-SPEC.md ··············· 측정·분석·AI 확장 (88~1792차원)
 ├── CARTRIDGE-STORE-SDK-SPEC.md ··················· 카트리지 스토어 & 개발자 SDK
 ├── admin-settings-llm-assistant-spec.md ·········· 관리자 설정·LLM 어시스턴트
 ├── cartridge-system-spec.md ······················ 카트리지 무한확장·접근 제어
 ├── data-packet-family-c.md ······················· 표준 데이터 패킷 (패밀리C)
 ├── proto-extension-vision-service.md ············· Vision Service Proto 확장
 └── B3-toss-pg-integration.md ····················· Toss PG 결제 연동

[L4] 인프라·품질·규제 명세
 ├── event-schema-specification.md ················· Kafka 이벤트 스키마
 ├── non-functional-requirements.md ················ NFR (성능·보안·확장성)
 ├── offline-capability-matrix.md ·················· 오프라인 기능 매트릭스
 ├── test-strategy.md ······························ 테스트 전략
 ├── deployment-strategy.md ························ 배포 전략
 ├── service-communication-patterns.md ············· 서비스 간 통신 패턴 (gRPC/Kafka/Saga)
 ├── plan-traceability-matrix.md ··················· 기획–구현 추적성 (80 REQ)
 └── api-security-architecture.md ·················· API 보안 아키텍처 (인증/KMS/PHI)

[L5] Sprint 실행 계획
 ├── SPRINT2-EXECUTION-PLAN.md ····················· Sprint 2 실행 계획
 ├── SPRINT2-AS-BACKEND-DETAIL.md ·················· 관리자 설정 백엔드 세부
 ├── SPRINT2-FLUTTER-DETAIL.md ····················· Flutter 프론트엔드 세부
 ├── SPRINT2-INFRA-E2E-DETAIL.md ··················· 인프라·E2E 세부
 └── SPRINT2-COMPLIANCE-DETAIL.md ·················· IEC 62304 규정 문서 세부

[L6] Phase·로드맵 문서
 ├── phase_1c_stage_s5b.md ························· Phase 1C-S5b: Rust FFI·BLE·차트
 ├── phase_1c_stage_s6.md ·························· Phase 1C-S6: 전체 통합
 ├── phase_1d_integration_mvp.md ··················· Phase 1D: 통합 MVP
 ├── phase_2_commerce_ai.md ························ Phase 2: 커머스+AI
 ├── PHASE11-17_MASTER_IMPLEMENTATION_PLAN.md ····· Phase 11~17 구현 계획
 ├── msa-expansion-roadmap.md ······················ MSA 확장 로드맵
 ├── ai-agent-phase-mapping.md ····················· AI 에이전트 Phase 매핑
 └── unimplemented-features-and-implementation-plan.md · 미구현 식별·계획

[L7] 에이전트 운영
 ├── AGENT-PROMPTS.md ······························ 에이전트 프롬프트 모음
 ├── AGENT-WORK-DISTRIBUTION-2026-02-12.md ········ 에이전트 업무분장
 └── agent-work-orders-v1.0.md ····················· 에이전트 작업 지시서

[L8] 검증 보고서
 ├── SYSTEM-COMPLETENESS-AUDIT-v1.md ··············· 기획 완전성 검증
 ├── implementation-gap-analysis-2026-02-12.md ····· Gap 분석
 ├── implementation-status-2026-02-11.md ··········· 구현 완성도 분석
 └── system-verification-and-implementation-plan.md · 기획서 검증·구현 현황 (100/100)

[L9] 규제 준수 보강
 └── compliance-gap-resolution.md ·················· FMEA·SOUP·SMP·SCM 통합
```

---

## 2. 핵심 참조 경로

| 작업 상황 | 참조 순서 |
| --- | --- |
| 전체 구조 파악 | L1 Blueprint → L2 FINAL-MASTER → 본 인덱스 |
| 특정 기능 구현 | L1 Blueprint 해당 Part → L2 FINAL-MASTER II절 → L3 해당 명세서 → L5 Sprint 세부 |
| 측정 파이프라인 | L3 MEASUREMENT-ANALYSIS-AI-SPEC → L1 Blueprint 10.15 → L4 data-packet-family-c |
| 카트리지 스토어 | L3 CARTRIDGE-STORE-SDK-SPEC → L1 Blueprint 10.16 → L3 cartridge-system-spec |
| AI 비서 | L3 AI-ASSISTANT-MASTER-SPEC → L1 Blueprint 10.14 |
| Phase 일정·태스크 | L2 FINAL-MASTER IV절 → L2 DETAILED-IMPLEMENTATION-PLAN |
| 갭 확인 | L2 FINAL-MASTER III절 → L8 SYSTEM-COMPLETENESS-AUDIT |
| 규제·추적성 | L4 plan-traceability-matrix → L5 SPRINT2-COMPLIANCE-DETAIL |

---

## 3. 문서 총계

| 계층 | 문서 수 |
| --- | --- |
| L0 원본 기획 | 3 |
| L1 종합 기획서 | 3 |
| L2 통합 구현 계획 | 2 |
| L3 기능별 명세서 | 8 |
| L4 인프라·품질·규제 | 8 |
| L5 Sprint 실행 | 5 |
| L6 Phase·로드맵 | 8 |
| L7 에이전트 운영 | 3 |
| L8 검증 보고서 | 4 |
| L9 규제 준수 보강 | 1 |
| **합계** | **45** |
