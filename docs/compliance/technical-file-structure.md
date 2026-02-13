# 만파식 기술문서(Technical File) 목차

> **문서 ID**: MPS-QMS-TFS-001
> **작성일**: 2026-02-09
> **작성자**: Claude (Regulatory & Security Analysis)
> **적용 대상**: FDA 510(k), CE-IVDR Annex II/III, MFDS 체외진단기기 허가, NMPA, PMDA

---

## 1. 공통 기술문서 마스터 구조

```
docs/technical-file/
│
├── 01-device-description/                    # 기기 설명
│   ├── 01-product-overview.md                # 제품 개요 및 의도된 사용
│   ├── 02-system-architecture.md             # 시스템 아키텍처 (HW + SW + Cloud)
│   ├── 03-functional-description.md          # 기능 설명 (차동측정, 핑거프린트, AI)
│   ├── 04-hardware-description.md            # 리더기 하드웨어 설명
│   ├── 05-software-description.md            # 소프트웨어 설명 (아키텍처, 모듈, 버전)
│   ├── 06-cartridge-system.md                # 29종 카트리지 시스템 상세
│   ├── 07-accessories-and-components.md      # 부속품 및 구성품 목록
│   └── 08-comparison-with-predicate.md       # 유사 기기 비교 (510(k)용)
│
├── 02-design-and-manufacturing/              # 설계 및 제조
│   ├── 01-design-input.md                    # 설계 입력 (요구사항)
│   ├── 02-design-output.md                   # 설계 출력 (아키텍처, 코드, 문서)
│   ├── 03-design-review-records.md           # 설계 검토 기록
│   ├── 04-design-transfer.md                 # 설계 이관 (개발→제조)
│   ├── 05-design-changes.md                  # 설계 변경 이력
│   └── 06-manufacturing-process.md           # 제조 프로세스 (카트리지, 리더기)
│
├── 03-software-documentation/                # 소프트웨어 문서 (IEC 62304)
│   ├── 01-software-development-plan.md       # 소프트웨어 개발 계획
│   ├── 02-software-requirements-spec.md      # SRS (소프트웨어 요구사항 명세)
│   ├── 03-software-architecture-doc.md       # SAD (소프트웨어 아키텍처 설계)
│   ├── 04-software-safety-classification.md  # → 기존 완료 문서 참조
│   ├── 05-subsystem-safety-allocation.md     # → 기존 완료 문서 참조 (섹션 5)
│   ├── 06-unit-test-results.md               # 단위 테스트 결과
│   ├── 07-integration-test-plan.md           # 통합 테스트 계획 및 결과
│   ├── 08-system-test-plan.md                # 시스템 테스트 계획 및 결과
│   ├── 09-traceability-matrix.md             # 추적성 매트릭스 (요구사항↔설계↔테스트↔코드)
│   ├── 10-soup-list-and-assessment.md        # SOUP 목록 및 위험 평가
│   ├── 11-anomaly-list.md                    # 알려진 이상/결함 목록
│   └── 12-release-notes.md                   # 릴리스 노트
│
├── 04-risk-management/                       # 위험관리 (ISO 14971)
│   └── → docs/risk-management/ 참조          # 별도 위험관리 파일 구조 (10+4 문서)
│
├── 05-verification-and-validation/           # 검증 및 확인
│   ├── 01-vv-master-plan.md                  # V&V 마스터 플랜
│   ├── 02-analytical-performance.md          # 분석 성능 (민감도, 특이도, 정확도, 정밀도)
│   │   ├── 바이오마커별 분석 성능
│   │   ├── 교차반응성 시험
│   │   ├── 간섭 시험
│   │   └── 측정 범위 (Reportable Range)
│   ├── 03-clinical-performance.md            # 임상 성능 평가
│   │   ├── 임상시험 프로토콜
│   │   ├── 대조 방법 비교 (기존 IVD 기기)
│   │   └── 임상 민감도/특이도
│   ├── 04-software-verification.md           # 소프트웨어 검증 결과
│   ├── 05-electrical-safety.md               # 전기 안전 (IEC 60601-1, 리더기)
│   ├── 06-emc-testing.md                     # EMC 시험 (IEC 60601-1-2, 리더기)
│   ├── 07-biocompatibility.md                # 생체적합성 (ISO 10993, 리더기)
│   ├── 08-stability-testing.md               # 안정성 시험 (카트리지 유효기간)
│   ├── 09-usability-testing.md               # 사용적합성 시험 (IEC 62366-1)
│   └── 10-cybersecurity-testing.md           # 사이버보안 시험 (침투 테스트 등)
│
├── 06-cybersecurity/                         # 사이버보안
│   ├── 01-threat-model.md                    # → 기존 STRIDE 문서 참조
│   ├── 02-sbom.md                            # 소프트웨어 BOM (CycloneDX)
│   ├── 03-vulnerability-management.md        # 취약점 관리 프로세스
│   ├── 04-security-update-plan.md            # 보안 업데이트 계획
│   └── 05-incident-response-plan.md          # → 기존 STRIDE 문서 참조
│
├── 07-data-protection/                       # 데이터 보호
│   ├── 01-data-protection-policy.md          # → 기존 완료 문서 참조
│   ├── 02-dpia.md                            # 데이터 보호 영향평가 (GDPR Art.35)
│   ├── 03-consent-management-spec.md         # 동의 관리 상세 사양
│   └── 04-data-localization-plan.md          # 국가별 데이터 현지화 계획
│
├── 08-labeling/                              # 라벨링 및 IFU
│   ├── 01-labeling-artwork.md                # 라벨 디자인
│   ├── 02-ifu-instructions.md                # 사용설명서 (Instructions for Use)
│   ├── 03-quick-start-guide.md               # 빠른 시작 가이드
│   └── translations/                         # 다국어 번역 (ko, en, ja, zh, de, fr)
│
├── 09-clinical-evidence/                     # 임상 근거 (CE-IVDR)
│   ├── 01-literature-review.md               # 문헌 조사
│   ├── 02-clinical-investigation-plan.md     # 임상시험 계획서
│   ├── 03-clinical-investigation-report.md   # 임상시험 보고서
│   └── 04-clinical-performance-summary.md    # 임상 성능 요약 (SSCP for IVDR)
│
├── 10-post-market/                           # 시판 후 관리
│   ├── 01-pms-plan.md                        # 시판 후 감시 계획
│   ├── 02-pmpf-plan.md                       # 시판 후 성능 추적 계획 (IVDR)
│   ├── 03-psur-template.md                   # 주기적 안전 업데이트 보고 양식
│   ├── 04-vigilance-reporting.md             # 의료기기 이상사례 보고 절차
│   └── 05-field-safety-corrective-action.md  # 현장 안전 시정 조치 절차
│
└── 11-regulatory-submissions/                # 인허가 제출 문서
    ├── kr-mfds/                              # 한국 MFDS 허가 서류
    │   ├── application-form.md
    │   └── country-specific-requirements.md
    ├── us-fda/                               # 미국 FDA 510(k)
    │   ├── 510k-cover-letter.md
    │   ├── substantial-equivalence.md
    │   ├── predicate-device-comparison.md
    │   └── level-of-concern.md
    ├── eu-ivdr/                              # EU CE-IVDR
    │   ├── declaration-of-conformity.md
    │   ├── gspr-checklist.md                 # General Safety & Performance Requirements
    │   └── udi-registration.md
    ├── cn-nmpa/                              # 중국 NMPA
    │   ├── registration-application.md
    │   └── cybersecurity-report.md
    └── jp-pmda/                              # 일본 PMDA
        ├── certification-application.md
        └── essential-principles.md
```

---

## 2. FDA 510(k) 제출 문서 구조 (상세)

| 섹션 | 문서 | 근거 | 소스 |
|------|------|------|------|
| I | Cover Letter | 21 CFR 807.87 | 신규 작성 |
| II | Indications for Use | 21 CFR 807.87(e) | `01-device-description/01-product-overview.md` |
| III | 510(k) Summary or Statement | 21 CFR 807.92 | 신규 작성 |
| IV | Truthful and Accuracy Statement | 21 CFR 807.87 | 서명 필요 |
| V | Device Description | 21 CFR 807.87(e) | `01-device-description/*` |
| VI | Substantial Equivalence | 21 CFR 807.87(f) | `11-regulatory-submissions/us-fda/substantial-equivalence.md` |
| VII | Performance Data | 21 CFR 807.87(g) | `05-verification-and-validation/02-analytical-performance.md` |
| VIII | Software Documentation | FDA SW Guidance | `03-software-documentation/*` |
| IX | Cybersecurity Documentation | FDA Cyber 2023 | `06-cybersecurity/*` |
| X | Labeling | 21 CFR 809.10 | `08-labeling/*` |
| XI | Biocompatibility | ISO 10993 | `05-verification-and-validation/07-biocompatibility.md` |
| XII | EMC/Electrical Safety | IEC 60601 | `05-verification-and-validation/05,06` |

---

## 3. CE-IVDR Technical Documentation (Annex II/III) 매핑

| IVDR Annex | 내용 | 대응 섹션 |
|------------|------|---------|
| Annex II, 1 | 기기 설명 및 사양 | `01-device-description/` |
| Annex II, 2 | 제조 정보 | `02-design-and-manufacturing/06` |
| Annex II, 3 | 설계/제조 정보 | `02-design-and-manufacturing/` + `03-software-documentation/` |
| Annex II, 4 | GSPR (Annex I) | `11-regulatory-submissions/eu-ivdr/gspr-checklist.md` |
| Annex II, 5 | 위험-편익 분석 | `04-risk-management/` |
| Annex II, 6 | 제품 검증/확인 | `05-verification-and-validation/` |
| Annex III | 기술문서 (성능 평가) | `09-clinical-evidence/` |
| Annex XIII | 성능 연구 | `09-clinical-evidence/02-clinical-investigation-plan.md` |

---

## 4. 현재 완료 문서 매핑

| 기술문서 섹션 | 기존 완료 문서 | 추가 작업 필요 |
|-------------|-------------|-------------|
| 03-04 (안전 등급) | `docs/compliance/software-safety-classification.md` | ❌ 완료 |
| 03-05 (서브시스템 등급) | 위 문서 섹션 5 | ❌ 완료 |
| 06-01 (STRIDE) | `docs/security/stride-threat-model.md` | 🔄 DFD 보완 |
| 06-05 (IRP) | 위 문서 섹션 5 | ❌ 완료 |
| 07-01 (데이터 보호) | `docs/compliance/data-protection-policy.md` | 🔄 DPIA 추가 |
| 04 (위험관리) | 목차만 존재 | 🔴 실체 문서 작성 필요 |
| 기타 전체 | - | 🔴 대부분 미착수 |

---

**Document Version**: 1.0.0
**작성일**: 2026-02-09
**작성자**: Claude (Regulatory & Security Analysis)
