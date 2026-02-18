# ManPaSik AI 생태계 구축 기획안 v1.1 완성본

**문서번호**: MPK-ECO-PLAN-v1.1-20260210-COMPLETE  
**기반**: 원본 v1.0 FINAL (MPK-ECO-PLAN-v1.0-20260208-FINAL) + 원본 대비 검증·발전 수립안 반영  
**목적**: 원본에 완벽 부합하고, 수정·보완·추가가 반영된 단일 기획 기준 문서. 개발 과정에서 고도화·유기적 대응의 베이스라인.

---

## 목차

1. [총괄 개요 (I)](#i-총괄-개요)
2. [생태계 철학 및 비전 (II)](#ii-생태계-철학-및-비전)
3. [기술 스택 (III)](#iii-기술-스택)
4. [전체 시스템 아키텍처 (IV)](#iv-전체-시스템-아키텍처)
5. [핵심 기능 상세 (V)](#v-핵심-기능-상세)
6. [사이트맵 (VI)](#vi-사이트맵)
7. [스토리보드 (VII)](#vii-스토리보드)
8. [AI 에이전트 자동화 (VIII)](#viii-ai-에이전트-자동화)
9. [데이터 아키텍처 (IX)](#ix-데이터-아키텍처)
10. [보안 및 규제 (X)](#x-보안-및-규제)
11. [개발 로드맵 (XI)](#xi-개발-로드맵)
12. [비용 분석 (XII)](#xii-비용-분석)
13. [가상 시뮬레이션 (XIII)](#xiii-가상-시뮬레이션)
14. [전문가 패널 검증 (XIV)](#xiv-전문가-패널-검증)
15. [결론 (XV)](#xv-결론)
16. [현행 시스템 매핑표 (XVI) — v1.1 신설](#xvi-현행-시스템-매핑표-v11-신설)
17. [단계별 MSA 확장 로드맵 (XVII) — v1.1 신설](#xvii-단계별-msa-확장-로드맵-v11-신설)
18. [기획-구현 추적성 (XVIII) — v1.1 신설](#xviii-기획-구현-추적성-v11-신설)
19. [세부 명세 참조 (XIX) — v1.1 신설](#xix-세부-명세-참조-v11-신설)

---

## I. 총괄 개요

- **프로젝트 정의**: 차동측정 기반 범용 분석(리더기+카트리지), 88→448→896차원 핑거프린트, 전자코/전자혀 융합, 세계 유일 범용 분석 생태계.
- **핵심 수치**: 측정 88→448→896차원, 정확도 92–96%, 측정 90초, 리더기 160×80×18mm, 카트리지 단가 1,135원(양산), 미지물질 탐지 97%, 8개 언어(초기)·50+ 확장, 동시 리더기 구독별 제한(최대 10대/Clinical), 오프라인 100%, SaaS 4티어(Free/Basic/Pro/Clinical), 카트리지 **무한확장 체계** (기본 29종 + 레지스트리 기반 무제한 추가, 최대 65,536종/2byte, 약 43억종/4byte), 구독 등급별 카트리지 접근 제어, TPM+해시체인+AES-256.
- **구독 티어 공식 대응**: [docs/plan/terminology-and-tier-mapping.md](terminology-and-tier-mapping.md) 참조.

**현행 반영**: CONTEXT.md, AGENTS.md §1–2.

---

## II. 생태계 철학 및 비전

- **홍익인간**: "세상의 모든 파동을 분석하여 인간을 이롭게 한다."
- **비전 계층**: 최상위 홍익인간 → 가치(치료→예방) → 플랫폼(만파식 AI 생태계) → 영역(건강/환경/식품/산업/안전) → 기술(차동측정+AI+전자코/전자혀) → 하드웨어(리더기+카트리지+E12-IF).
- **유기적 성장**: 88→448→896→1792차원, 생물체 단계 비유.

**현행 반영**: AGENTS.md 프로젝트 아이덴티티.

---

## III. 기술 스택

- **프론트**: Flutter 3.x, Next.js 14+, Material Design 3, Riverpod 2.x, go_router, SQLite+Hive(Flutter), Sled(Rust).
- **코어**: Rust(no_std), TFLite+ONNX, btleplug, nfc, rustfft+dasp, ring+aes-gcm, tokio, flutter_rust_bridge.
- **백엔드**: Kong, Go+gRPC, Keycloak OIDC, PostgreSQL 16, TimescaleDB, Redis, Elasticsearch, MinIO, Kafka, K8s, Prometheus+Grafana, ELK.
- **AI/ML**: PyTorch, TFLite, Milvus, 연합학습(Flower), MLOps(MLflow 등). 상세 [docs/ai-specs/claude/ml-model-design-spec.md](../ai-specs/claude/ml-model-design-spec.md).

**현행 반영**: .cursor/rules/manpasik-project.mdc, CONTEXT.md 기술 스택 표.

---

## IV. 전체 시스템 아키텍처

- **MSA**: 원본 8도메인 30+ 서비스. **현행 Phase 1**: auth, user, device, measurement 4개. Phase 2–4 확장은 [docs/plan/msa-expansion-roadmap.md](msa-expansion-roadmap.md) 참조.
- **데이터 흐름**: 카트리지→RAFE→STM32→BLE→Rust(차동측정·DSP·AI·핑거프린트)→로컬 저장→API Gateway→measurement/ai/coaching 등→TimescaleDB+Milvus→UI.
- **오프라인-온라인**: CRDT 동기화, Delta Sync, LWW+벡터 클록. [docs/specs/offline-capability-matrix.md](../specs/offline-capability-matrix.md) 참조.

**현행 반영**: AGENTS.md §3, backend/services/, rust-core/.

---

## V. 핵심 기능 상세

- **5.1** SaaS 구독·개인 주치의: 4티어(Free/Basic/Pro/Clinical). [terminology-and-tier-mapping.md](terminology-and-tier-mapping.md).
- **5.2** 다수 리더기: BLE 최대 구독별 N대, Wi-Fi Direct, 위치 대시보드(Phase 2).
- **5.3** 계층형 관리자: 총괄→국가→지역→지점→판매점. Phase 3, admin-service.
- **5.4** 화상진료/병원·약국 예약: Phase 3.
- **5.5** 통합 쇼핑몰: Phase 2.
- **5.6** SDK·카트리지 마켓: Phase 4. manpasik-sdk 기반 서드파티 카트리지 개발→검증→등록→마켓플레이스 게시 워크플로우. 카테고리 0xF0~0xFD 동적 할당, 수익분배(만파식:벤더 = 30:70). 서드파티 카트리지도 등급별 접근 제어 적용 (기본: ADD_ON).
- **5.7** 건강관리 코칭: Phase 2, [ml-model-design-spec.md](../ai-specs/claude/ml-model-design-spec.md).
- **5.8** 오프라인 완전 구동: [offline-capability-matrix.md](../specs/offline-capability-matrix.md).
- **5.9** 카트리지 무한확장 체계 및 자동인식:
  - **무한확장 레지스트리**: 2-Byte 계층형 코드(카테고리 u8 × 타입 u8 = 65,536종), 4-Byte 확장(Phase 4+, 약 43억종). 코드 배포 없이 서버 레지스트리(DB) 등록만으로 신규 카트리지 추가.
  - **자동인식**: NFC ISO 14443A + Rust nfc 모듈. v2.0 태그 포맷(카테고리+타입+레거시 호환 코드+보정데이터). v1.0 레거시(29종) 자동 변환.
  - **등급별 접근 제어**: 구독 티어(Free/Basic/Pro/Clinical)에 따라 사용 가능 카트리지 범위 결정. 카테고리/타입별 세분화 정책(INCLUDED/LIMITED/ADD_ON/RESTRICTED/BETA). 관리자가 코드 배포 없이 DB 정책 변경 가능. 상세 [docs/specs/cartridge-system-spec.md](../specs/cartridge-system-spec.md).
  - **카테고리 체계**: HealthBiomarker(0x01), Environmental(0x02), FoodSafety(0x03), ElectronicSensor(0x04), AdvancedAnalysis(0x05), Industrial(0x06), Veterinary(0x07), Pharmaceutical(0x08), Agricultural(0x09), Cosmetic(0x0A), Forensic(0x0B), Marine(0x0C), Reserved(0x0D~0xEF), ThirdParty(0xF0~0xFD), Beta(0xFE), CustomResearch(0xFF).
  - **서드파티 확장**: Phase 4에서 manpasik-sdk로 외부 개발자가 카트리지 설계·제조·마켓플레이스 등록. 카테고리 0xF0~0xFD 동적 할당.
- **5.10** 글로벌 규제: 5개국 인증·데이터보호·품질·SW수명주기. [docs/compliance/regulatory-compliance-checklist.md](../compliance/regulatory-compliance-checklist.md).
- **5.11** 글로벌 커뮤니티·실시간 번역: Phase 3.
- **5.12** 자기학습 성장 AI: 연합학습·지속학습, [ml-model-design-spec.md](../ai-specs/claude/ml-model-design-spec.md).
- **5.13** 음성 명령/접근성: Phase 4–5.
- **5.14** 유기적 확장: SDK, OTA, E12-IF, 1792차원 경로.

**현행 반영**: CONTEXT, AGENTS, rust-core 모듈, backend 서비스.

---

## VI. 사이트맵

- **공식 문서**: [docs/ux/sitemap.md](../ux/sitemap.md). 원본 VI 절 트리 구조, 라우트 ID·Phase·인증 매핑.
- **요약**: 홈, 측정, 데이터허브, AI코치, 마켓, 커뮤니티, 의료서비스, 기기관리, 가족, 설정, 관리자 포탈.

---

## VII. 스토리보드

### Phase 1 (핵심 경험)
- **인증/온보딩**: [storyboard-auth-onboarding.md](../ux/storyboard-auth-onboarding.md) — 인트로, 소셜로그인, 약관동의, 온보딩, 리더기 페어링
- **홈 대시보드**: [storyboard-home-dashboard.md](../ux/storyboard-home-dashboard.md) — 건강 요약, AI 코칭 요약, 빠른 측정
- **첫 측정**: [storyboard-first-measurement.md](../ux/storyboard-first-measurement.md) — 카트리지 인식→측정→결과→AI 분석
- **기기 관리**: [storyboard-device-management.md](../ux/storyboard-device-management.md) — BLE 페어링, 펌웨어, 구독별 대수 제한
- **설정**: [storyboard-settings.md](../ux/storyboard-settings.md) — 계정, 구독, 알림, 접근성, 긴급대응 설정
- **오프라인 동기화**: [storyboard-offline-sync.md](../ux/storyboard-offline-sync.md) — 오프라인 측정, 데이터 큐, 자동 동기화

### Phase 2 (AI/커머스)
- **음식 촬영→칼로리**: [storyboard-food-calorie.md](../ux/storyboard-food-calorie.md) — 사진 분석, 영양소, 식단 기록
- **AI 비서**: [storyboard-ai-assistant.md](../ux/storyboard-ai-assistant.md) — 대화형 AI 상담, 건강 코칭, 운동/식단 추천
- **데이터 허브**: [storyboard-data-hub.md](../ux/storyboard-data-hub.md) — 타임라인, 트렌드 차트, FHIR 내보내기
- **마켓 구매**: [storyboard-market-purchase.md](../ux/storyboard-market-purchase.md) — 카트리지 스토어, 장바구니, Toss PG 결제
- **구독 전환**: [storyboard-subscription-upgrade.md](../ux/storyboard-subscription-upgrade.md) — 티어 비교, 업/다운그레이드, 해지

### Phase 3 (의료/커뮤니티/가족)
- **화상진료**: [storyboard-telemedicine.md](../ux/storyboard-telemedicine.md) — 예약, 대기실, WebRTC 진료, 처방전
- **가족 관리**: [storyboard-family-management.md](../ux/storyboard-family-management.md) — 그룹 생성, 보호자 대시보드, 119 연동
- **커뮤니티**: [storyboard-community.md](../ux/storyboard-community.md) — 건강 포럼, 챌린지, 전문가 Q&A, 번역 채팅
- **긴급 대응**: [storyboard-emergency-response.md](../ux/storyboard-emergency-response.md) — 이상 감지, 에스컬레이션, 119 자동 신고
- **관리자 포탈**: [storyboard-admin-portal.md](../ux/storyboard-admin-portal.md) — KPI 대시보드, 회원/재고/매출 관리

---

## VIII. AI 에이전트 자동화

- **Phase·서비스 매핑**: [docs/plan/ai-agent-phase-mapping.md](ai-agent-phase-mapping.md). 사용자 인터페이스, 측정 자동화, 건강 관리, 시스템 관리, 긴급 대응 블록별 담당 Phase·서비스·구현 담당.

**현행 반영**: Rust 측정 자동화(ble, nfc, differential, ai), device-service(OTA).

---

## IX. 데이터 아키텍처

- **데이터 패킷 표준(패밀리C)**: [docs/specs/data-packet-family-c.md](../specs/data-packet-family-c.md). header/payload/footer, transform_log, state_meta, Proto 매핑.
- **DB 스키마**: users, families, devices, cartridges, measurements, measurement_results, fingerprint_vectors(Milvus), health_records, audit_logs 등. tenant_id 필수. [AGENTS.md §5](../../AGENTS.md), [plan-original-vs-current-and-development-proposal.md](../plan-original-vs-current-and-development-proposal.md) 3.2 구현·스키마 보완 참조.

---

## X. 보안 및 규제

- **보안 계층**: TPM, 펌웨어 보안 부팅, BLE AES-CCM, HTTPS TLS 1.3, AES-256-GCM, 해시체인, 감사추적, OIDC, MFA, RBAC.
- **규제**: MFDS/FDA/CE-IVDR/NMPA/PMDA, PIPA/HIPAA/GDPR/PIPL/APPI, ISO 13485, IEC 62304, ISO 14971.
- **문서**: [STRIDE](../security/stride-threat-model.md), [데이터 보호 정책](../compliance/data-protection-policy.md), [위험관리 계획서](../compliance/iso14971-risk-management-plan.md), [규제 체크리스트](../compliance/regulatory-compliance-checklist.md).

---

## XI. 개발 로드맵

- **Phase 1 (1–4개월)**: 기본 UI, BLE/NFC, 차동측정, 88차원, 시각화, 클라우드 연동, 오프라인 기본. 리더기 최대 1대(Free). **현행**: Rust 코어·FFI 완료, Go 4서비스 구현 중, Flutter 대기.
- **Phase 2 (5–8)**: AI 코칭, 음식 칼로리, 마켓·결제·구독, HealthKit 연동, 데이터 허브. 리더기 최대: Free 1대, Basic 2대, Pro 5대.
- **Phase 3 (9–12)**: 화상진료, 병원/약국 예약, 가족, 커뮤니티, 실시간 번역, 계층형 관리자, 긴급대응(119). 리더기 최대: Clinical 10대.
- **Phase 4 (13–18)**: SDK, 서드파티 카트리지 마켓, 896차원·전자코/전자혀, B2B API, 글로벌 인증(FDA/CE-IVDR), 연합학습, FHIR R4 의료 데이터 교환.
- **Phase 5 (19–24)**: AI 에이전트 완전 자동화, 1792차원, 음성 명령, 웨어러블 연동, 스마트홈 IoT.

---

## XII. 비용 분석

- **24개월 개발 비용**: 약 67억원 (모바일 8, 웹 4, 백엔드 12, Rust 6, AI/ML 10, 실시간번역 3, 화상진료 2.5, 관리자 3, UI/UX 2.5, 인프라 6, QA 3, 규제 5, PM 2).
- **피크 인력**: 약 32명.

---

## XIII. 가상 시뮬레이션

- **부하**: 동시 100,000 MAU 시나리오. Phase 2/3 인프라·성능 검증 시 재수행.
- **오프라인**: 72시간 완전 오프라인. Phase 1 통합 테스트에서 검증.
- **896차원 검색**: 100만 벡터 Milvus. Phase 2 벡터 검색 수용 기준으로 활용.

---

## XIV. 전문가 패널 검증

- 원본 v1.0: 25인 전문가 패널 만장일치 승인. 12개 항목 93–99점, 종합 96.6/100. **APPROVED FOR IMPLEMENTATION**.

---

## XV. 결론

- **5대 가치**: 범용성, 정밀성, 접근성, 확장성, 연결성.
- **7대 차별화**: 2단계 매트릭스 제거, 전자코+전자혀 융합, 896차원, 3단계 역추론, TPM+해시체인, 디지털 트윈 자기최적화, EHD 무소음 기체 제어.

---

## XVI. 현행 시스템 매핑표 (v1.1 신설)

| 원본 절 | 현재 반영 위치 | Phase/비고 |
|---------|----------------|------------|
| I 총괄 | CONTEXT.md, AGENTS.md §1–2 | — |
| II 철학·비전 | AGENTS.md 프로젝트 아이덴티티 | — |
| III 기술 스택 | manpasik-project.mdc, CONTEXT 기술 스택 표 | — |
| IV 아키텍처 | AGENTS.md §3, backend/services/, rust-core/ | Phase 1: 4서비스 |
| V 핵심 기능 | CONTEXT, AGENTS, rust-core, backend, docs/specs·ux | 5.8·5.9·5.10 구현/문서 반영 |
| VI 사이트맵 | docs/ux/sitemap.md | — |
| VII 스토리보드 | docs/ux/storyboard-*.md | — |
| VIII 에이전트 | docs/plan/ai-agent-phase-mapping.md | Phase별 |
| IX 데이터 | docs/specs/data-packet-family-c.md, AGENTS §5, Proto | — |
| X 보안·규제 | docs/security, docs/compliance | — |
| XI 로드맵 | CONTEXT 현재 진행 단계, 본 문서 XI | — |
| XII 비용 | AGENTS.md | — |
| XIII 시뮬레이션 | Phase 2/3 검증 시 재사용 | — |
| XIV·XV | 원본 유지 | — |

---

## XVII. 단계별 MSA 확장 로드맵 (v1.1 신설)

- **상세**: [docs/plan/msa-expansion-roadmap.md](msa-expansion-roadmap.md).
- **요약**: Phase 1 auth/user/device/measurement → Phase 2 subscription, shop, payment, ai-inference, coaching, cartridge, calibration → Phase 3 family, health-record, telemedicine, reservation, prescription, community, translation, video, notification, admin → Phase 4 marketplace, ai-training, vision, nlp, inventory, logistics, analytics, iot-gateway, location, emergency.

---

## XVIII. 기획-구현 추적성 (v1.1 신설)

- **상세**: [docs/plan/plan-traceability-matrix.md](plan-traceability-matrix.md). 요구 ID(REQ-*) ↔ 설계(DES) ↔ 구현(IMP) ↔ V&V.
- **유지**: 신규 요구·설계·구현 시 매트릭스 갱신. 의료기기 인허가·감사 대비.

---

## XIX. 세부 명세 참조 (v1.1 신설)

| 명세 | 경로 |
|------|------|
| **카트리지 무한확장 체계 명세** | [docs/specs/cartridge-system-spec.md](../specs/cartridge-system-spec.md) |
| 데이터 패킷 (패밀리C) | [docs/specs/data-packet-family-c.md](../specs/data-packet-family-c.md) |
| 사이트맵 | [docs/ux/sitemap.md](../ux/sitemap.md) |
| 스토리보드 (16건) | [docs/ux/storyboard-*.md](../ux/) — 인증, 홈, 측정, 기기, 설정, 오프라인, 음식칼로리, AI비서, 데이터허브, 마켓, 구독, 화상진료, 가족, 커뮤니티, 긴급대응, 관리자 |
| 오프라인 기능 매트릭스 | [docs/specs/offline-capability-matrix.md](../specs/offline-capability-matrix.md) |
| MSA 확장 로드맵 | [docs/plan/msa-expansion-roadmap.md](msa-expansion-roadmap.md) |
| 기획-구현 추적성 | [docs/plan/plan-traceability-matrix.md](plan-traceability-matrix.md) |
| AI 에이전트-Phase 매핑 | [docs/plan/ai-agent-phase-mapping.md](ai-agent-phase-mapping.md) |
| 용어·티어 통일 | [docs/plan/terminology-and-tier-mapping.md](terminology-and-tier-mapping.md) |
| 개발 철학 | [docs/development-philosophy.md](../development-philosophy.md) |
| 원본 대비 검증·발전 수립안 | [docs/plan-original-vs-current-and-development-proposal.md](../plan-original-vs-current-and-development-proposal.md) |
| **최종 검증 보고서** | [docs/plan/plan-verification-report.md](plan-verification-report.md) |
| **원본 상세 반영 보조(Annex)** | [docs/plan/original-detail-annex.md](original-detail-annex.md) |

---

**문서 종료**

*본 기획안 v1.1 완성본은 원본 MPK-ECO-PLAN-v1.0-20260208-FINAL에 완벽 부합하도록 수정·보완·추가하였으며, 개발 과정에서 [docs/development-philosophy.md](../development-philosophy.md)에 따라 고도화 및 유기적 대응으로 갱신한다.*
