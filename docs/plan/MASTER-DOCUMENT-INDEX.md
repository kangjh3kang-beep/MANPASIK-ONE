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
 ├── compliance-gap-resolution.md ·················· FMEA·SOUP·SMP·SCM 통합
 ├── iso14971-risk-management-plan.md ·············· ISO 14971 위험관리계획
 ├── iso14971-fmea.md ····························· FMEA (38개 고장 모드, 8개 서브시스템)
 ├── iso14971-pha.md ······························ 예비 위해 분석 (34개 위해, 7개 카테고리)
 ├── regulatory-compliance-checklist.md ············ 5개국 규제 준수 체크리스트
 ├── data-protection-policy.md ···················· 데이터 보호 정책 (HIPAA/GDPR/PIPA/PIPL/APPI)
 └── vnv-master-plan.md ·························· V&V 마스터 플랜

[L10] UX/디자인 산출물
 ├── sitemap.md ··································· 앱 사이트맵 (라우트·Phase 매핑)
 ├── BRAND_GUIDELINE.md ··························· 상감 디자인 시스템 브랜드 가이드라인
 ├── STYLE_GUIDE_CODE.md ·························· 코드 스타일 가이드
 ├── DESIGN_GENERATION_STRATEGY.md ················ 디자인 생성 전략
 ├── storyboard-auth-onboarding.md ················ UX 스토리보드: 인증/온보딩
 ├── storyboard-home-dashboard.md ················· UX 스토리보드: 홈 대시보드
 ├── storyboard-first-measurement.md ·············· UX 스토리보드: 첫 측정 여정
 ├── storyboard-device-management.md ·············· UX 스토리보드: 기기 관리
 ├── storyboard-settings.md ······················· UX 스토리보드: 설정
 ├── storyboard-offline-sync.md ··················· UX 스토리보드: 오프라인 동기화
 ├── storyboard-food-calorie.md ··················· UX 스토리보드: 음식 칼로리
 ├── storyboard-ai-assistant.md ··················· UX 스토리보드: AI 비서
 ├── storyboard-data-hub.md ······················· UX 스토리보드: 데이터 허브
 ├── storyboard-market-purchase.md ················ UX 스토리보드: 카트리지 구매
 ├── storyboard-subscription-upgrade.md ··········· UX 스토리보드: 구독 전환
 ├── storyboard-telemedicine.md ··················· UX 스토리보드: 비대면 진료
 ├── storyboard-family-management.md ·············· UX 스토리보드: 가족 관리
 ├── storyboard-community.md ····················· UX 스토리보드: 커뮤니티
 ├── storyboard-emergency-response.md ············· UX 스토리보드: 긴급 대응
 └── storyboard-admin-portal.md ··················· UX 스토리보드: 관리자 포탈

[L11] 테스트 산출물
 ├── tests/e2e/ (8개 파일) ························ E2E 테스트 (29개 시나리오)
 │   ├── helpers_test.go ·························· 공용 테스트 유틸리티
 │   ├── auth_flow_test.go ························ 인증 플로우 (5개)
 │   ├── measurement_flow_test.go ················· 측정 파이프라인 (5개)
 │   ├── device_management_test.go ················ 디바이스 관리 (4개)
 │   ├── medical_service_test.go ·················· 의료 서비스 (5개)
 │   ├── community_family_test.go ················· 커뮤니티/가족 (5개)
 │   └── admin_test.go ···························· 관리자 (5개)
 ├── tests/security/ (4개 파일) ··················· 보안 테스트 (OWASP Top 10)
 │   ├── auth_security_test.go ···················· 인증/권한 보안 (8개)
 │   ├── api_security_test.go ····················· API 보안 (9개)
 │   ├── data_security_test.go ···················· 데이터 보안 (11개)
 │   └── dependency_scan.sh ······················· 의존성 취약점 스캔
 └── tests/load/ (7개 파일) ······················· 부하 테스트 (k6)
     ├── config.js ································ 공용 설정
     ├── auth_load_test.js ························ 인증 부하
     ├── measurement_load_test.js ················· 측정 부하
     ├── stress_test.js ··························· 스트레스 테스트
     ├── spike_test.js ···························· 스파이크 테스트
     ├── concurrent_users_test.js ················· 동시 사용자
     └── api_gateway_load_test.js ················· API 게이트웨이 부하
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
| L9 규제 준수 보강 | 7 |
| L10 UX/디자인 산출물 | 24 |
| L11 테스트 산출물 | 19 |
| **합계** | **94** |

---

## 4. 구현 산출물 (코드)

| 영역 | 산출물 | 상태 |
| --- | --- | --- |
| Flutter 앱 (12개 Feature) | auth, home, measurement, devices, chat, settings, admin, ai_coach, data_hub, community, family, market, medical | ✅ screen + domain |
| Rust 코어 엔진 (10개 모듈) | ble, nfc, dsp, ai, crypto, sync, storage, ffi, config, error | ✅ 10모듈 + 보강 (BLE 상태머신, AI 안전검증) |
| Go 백엔드 (22개 서비스) | auth, user, measurement, device, cartridge, notification, family, health-record, telemedicine, reservation, community, admin, prescription, translation, video + gateway | ✅ gRPC 핸들러 + 서비스 + 리포지토리 |
| 인프라 | Docker Compose, PostgreSQL 21개 스키마, Kafka, Redis | ✅ 구성 완료 |
