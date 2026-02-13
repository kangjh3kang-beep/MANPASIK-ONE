# 만파식 시스템 구축 기획 완전성 검증 보고서

**문서번호**: MPK-AUDIT-v1.0  
**검증일**: 2026-02-12  
**목적**: 전체 시스템 구축 기획, 사이트맵, 스토리보드에 부합하는 상세 페이지 기획·구현 계획이 빠짐없이 수립되었는지 분석·검증한다.

---

## 1. 검증 범위 및 방법

### 1.1 검증 대상 문서

| 문서 | 역할 |
| --- | --- |
| COMPREHENSIVE-SYSTEM-BLUEPRINT-v3.0.md | 마스터 기획서 (Part 1~11) |
| FINAL-MASTER-IMPLEMENTATION-PLAN.md | 통합 구현·구축 계획 (I~V) |
| AI-ASSISTANT-MASTER-SPEC.md | AI 비서 상세 기획 |
| MEASUREMENT-ANALYSIS-AI-SPEC.md | 측정·분석·AI 확장 상세 기획 |
| CARTRIDGE-STORE-SDK-SPEC.md | 카트리지 스토어·SDK 상세 기획 |
| frontend/flutter-app/lib/core/router/app_router.dart | 실제 구현 라우트 |

### 1.2 검증 방법

1. **사이트맵 교차 대조**: Blueprint Part 3·6 사이트맵 ↔ FINAL-MASTER I.2·II절 ↔ 상세 기획서 ↔ 실제 Flutter 라우트
2. **기능-페이지-API-DB 4중 매핑**: Part 10 전 기능에 대해 (페이지, API, DB, 이벤트) 4요소 완전성 확인
3. **Phase 로드맵 정합성**: IV절 Phase에 모든 기능이 할당되었는지 확인
4. **갭 분석 유효성**: III절 기존 갭이 해결되었는지 확인

---

## 2. 사이트맵·스토리보드 대조 결과

### 2.1 전체 페이지·라우트 통합 현황 (50+ 화면)

| 구분 | 페이지/라우트 수 | 상세 기획 존재 | 실제 구현 |
| --- | --- | --- | --- |
| 인증·온보딩 (Splash, Login, Register) | 3 | O (II.1) | O |
| 홈 (HomeScreen) | 1 | O (6.5, II.2) | O |
| 측정 (Measurement, Result) | 2 | O (6.5, II.2, MEASURE-AI) | O |
| 데이터허브 (DataHub 4탭) | 4 | O (6.5) | X |
| AI 코치 (AiCoach) | 1 | O (6.5, 10.2) | X |
| 마켓 (Shop, Cart, Order, Payment) | 4 | O (II.4) | X |
| 커뮤니티 (Community 3탭) | 3 | O (II.5) | X |
| 의료 (Medical 4탭) | 4 | O (II.6) | X |
| 전문가 (ExpertDashboard, PatientDetail) | 2 | O (6.6, 10.11) | X |
| 기기 (Devices) | 1 | O (II.3) | O |
| 가족 (Family, 가족 건강 대시) | 2 | O (II.7) | X |
| 설정 (Settings 6그룹) | 1 | O (II.8) | O |
| 관리자 (Admin 3탭+LLM) | 3 | O (II.9) | X |
| 알림 (알림 센터 드로어) | 1 | O (II.10) | X |
| AI 비서 (전역 플로팅·대화) | 2 | O (II.11, AI-ASSISTANT) | X |
| 카트리지 스토어 (6라우트) | 6 | O (II.13, CART-STORE) | X |
| 목적별 대시 (dashboard, family/health, location, water, air, org) | 7 | O (10.13) | X |
| **합계** | **~47** | **47/47 = 100%** | **8/47 = 17%** |

**결론**: 모든 페이지·화면에 대한 상세 기획이 존재합니다. 구현은 Phase 1 범위(8개)가 완료되어 있습니다.

### 2.2 사이트맵 계층 구조 검증

```text
만파식 앱
├── / (Splash) → 토큰 검증 → /home 또는 /login     ✅ 기획+구현
├── /login                                          ✅ 기획+구현
├── /register                                       ✅ 기획+구현
├── /home (HomeScreen)                              ✅ 기획+구현
│   ├── 건강 요약 카드                               ✅ 기획 (6.5)
│   ├── 환경 요약                                    ✅ 기획 (10.12)
│   ├── AI 코칭 요약                                 ✅ 기획 (10.2)
│   ├── 빠른 측정                                    ✅ 기획 (II.2)
│   ├── AI 비서 진입                                 ✅ 기획 (10.14)
│   └── 알림 아이콘                                  ✅ 기획 (II.10)
├── /measurement                                    ✅ 기획+구현
│   ├── 카트리지 인식                                ✅ 기획 (MEASURE-AI §3)
│   ├── 세션 관리                                    ✅ 기획 (MEASURE-AI §3)
│   └── /measurement/result                         ✅ 기획+구현
├── /data (DataHub)                                 ✅ 기획 (6.5)
│   ├── 요약 탭                                     ✅ 기획
│   ├── 타임라인 탭                                  ✅ 기획
│   ├── 트렌드 탭                                    ✅ 기획
│   └── 내 기준선 탭                                 ✅ 기획
├── /coach (AiCoach)                                ✅ 기획 (6.5, 10.2)
│   ├── 대화형 상담                                  ✅ 기획
│   ├── 일/주간 리포트                               ✅ 기획
│   ├── 목표 설정                                    ✅ 기획
│   └── 음식 칼로리 분석                              ✅ 기획 (10.6)
├── /market (Shop)                                  ✅ 기획 (II.4)
│   ├── 상품 목록·필터·검색                           ✅ 기획
│   ├── 장바구니                                     ✅ 기획
│   ├── 주문·결제                                    ✅ 기획
│   └── 구독 관리                                    ✅ 기획
├── /store (카트리지 스토어)                          ✅ 기획 (CART-STORE §4)
│   ├── /store/search                               ✅ 기획
│   ├── /store/category/:id                         ✅ 기획
│   ├── /store/item/:id                             ✅ 기획
│   ├── /store/my-cartridges                        ✅ 기획
│   └── /store/purchase-history                     ✅ 기획
├── /community                                      ✅ 기획 (II.5)
│   ├── 피드                                        ✅ 기획
│   ├── 챌린지                                       ✅ 기획
│   └── 내 글                                       ✅ 기획
├── /medical                                        ✅ 기획 (II.6)
│   ├── 예약                                        ✅ 기획
│   ├── 화상진료                                     ✅ 기획 (10.4)
│   ├── 처방                                        ✅ 기획
│   ├── 건강데이터 공유                               ✅ 기획
│   ├── /medical/expert                             ✅ 기획 (6.6)
│   └── /medical/expert/:id                         ✅ 기획 (6.6, 10.11)
├── /devices                                        ✅ 기획+구현
├── /family                                         ✅ 기획 (II.7)
│   └── /family/health                              ✅ 기획 (10.13)
├── /settings                                       ✅ 기획+구현
│   ├── 계정                                        ✅ 기획
│   ├── 구독                                        ✅ 기획
│   ├── 알림                                        ✅ 기획
│   ├── 접근성                                       ✅ 기획
│   ├── 긴급 대응                                    ✅ 기획
│   └── 약관·개인정보                                 ✅ 기획
├── /admin/*                                        ✅ 기획 (II.9)
│   ├── 설정 관리                                    ✅ 기획
│   ├── LLM 어시스턴트                               ✅ 기획
│   └── 감사·통계                                    ✅ 기획
├── /dashboard (목적별 본인)                          ✅ 기획 (10.13)
├── /location/:id/env (위치 환경)                     ✅ 기획 (10.13)
├── /water (수질)                                    ✅ 기획 (10.13)
├── /air (공기질)                                    ✅ 기획 (10.13)
├── /org (기업 포털)                                 ✅ 기획 (10.13)
│   └── /org/unit/:id (부서/공장/지역)               ✅ 기획 (10.13)
└── [전역] AI 비서 플로팅                             ✅ 기획 (10.14, AI-ASSISTANT)
```

---

## 3. 기능-페이지-API-DB 4중 매핑 검증

### 3.1 완전 매핑된 기능 (문제 없음)

| 기능 | 페이지 | API | DB | 이벤트 | 판정 |
| --- | --- | --- | --- | --- | --- |
| 인증 | O | O | O | O | ✅ 완전 |
| 본인 건강 요약 | O | O | O | O | ✅ 완전 |
| 측정 | O | O | O | O | ✅ 완전 |
| AI 주치의·코칭 (10.2) | O | O | O | O | ✅ 완전 |
| 전문가 진료 입력 (10.11) | O | O | O | O | ✅ 완전 |
| 마켓·결제 | O | O | O | O | ✅ 완전 |
| 예약·화상·처방 | O | O | O | O | ✅ 완전 |
| 커뮤니티 | O | O | O | O | ✅ 완전 |
| 기기·카트리지 | O | O | O | O | ✅ 완전 |
| 관리자 | O | O | O | O | ✅ 완전 |
| AI 비서 (10.14) | O | O | O | O | ✅ 완전 |
| 측정·분석·AI (10.15) | O | O | O | O | ✅ 완전 |
| 카트리지 스토어 (10.16) | O | O | O | O | ✅ 완전 |

### 3.2 부분 누락이 있는 기능 (보강 필요)

| 기능 | 페이지 | API | DB | 이벤트 | 누락 사항 |
| --- | --- | --- | --- | --- | --- |
| 식단관리 (10.5) | △ (AiCoach 내) | △ (제안만) | ✗ | ✗ | API 명세 미확정, DB 테이블·이벤트 미정의 |
| 칼로리관리 (10.6) | △ (AiCoach 내) | △ (제안만) | ✗ | ✗ | AnalyzeFoodImage 제안만, DB·이벤트 미정의 |
| 헬스코칭 보상 (10.7) | △ (AiCoach 내) | △ (제안만) | △ (제안만) | △ (제안만) | RewardService 미확정, 보상 정책 테이블 미정의 |
| 음성 인식·번역 (10.8) | △ (설정 내) | △ (제안만) | △ (제안만) | ✗ | VoiceProfileService 미확정 |
| 익명 학습 AI (10.9) | ✗ (백엔드) | △ (제안만) | △ (제안만) | △ (제안만) | 익명화 파이프라인·모델 저장소 미정의 |
| 데이터 제공 (10.10) | △ (기업 포탈) | △ (제안만) | △ (제안만) | ✗ | DataProvisionService·B2B 계약 테이블 미정의 |
| 지역 통계 (10.12) | △ (Home/Data) | △ (제안만) | △ (제안만) | △ (제안만) | LocationStatsService·location_statistics 미정의 |
| 목적별·레고형 (10.13) | O | △ (제안만) | △ (제안만) | ✗ | ConceptService·OrganizationService 미확정 |
| 가족 | O | O | O | ✗ | 이벤트 미정의 |
| 설정·알림 | O | O | O | ✗ | 이벤트 미정의 |
| 검증된 전문가 매칭 (10.3) | △ | △ (제안만) | △ | ✗ | 매칭 알고리즘·DB 미정의 |
| 화상·번역 (10.4) | O | △ (TranslationService 제안) | △ | △ | TranslationService 미확정 |

---

## 4. FINAL-MASTER 정합성 검증

### 4.1 I.2 매트릭스 이벤트 열 빈 값

| 기능 영역 | 현재 이벤트 열 | 보강 필요 이벤트 |
| --- | --- | --- |
| 식단·칼로리 | `-` | `meal.logged`, `nutrition.analyzed` |
| 가족 | `-` | `family.member.invited`, `family.sharing.updated` |
| 설정·알림 | `-` | `notification.preferences.updated`, `profile.updated` |

### 4.2 IV절 Phase 로드맵 누락

| 기능 | II절 위치 | Phase 포함 여부 | 권장 Phase |
| --- | --- | --- | --- |
| AI 비서 (II.11) | II.11 | ✗ 미포함 | Phase 2~3 (기본) + Phase 5 (음성 고도화) |
| 측정·분석·AI 확장 (II.12) | II.12 | ✗ 미포함 | Phase 1(88-dim) + Phase 2(448) + Phase 3(896) + Phase 5(1792) |
| 카트리지 스토어·SDK (II.13) | II.13 | ✗ 미포함 | Phase 3(1st-party) + Phase 4(SDK 공개) + Phase 5(전면 오픈) |

### 4.3 III절 갭 항목 유효성

모든 10개 갭 항목이 **여전히 유효**합니다 (미해결 상태).

---

## 5. 종합 판정

### 5.1 완전성 점수

| 검증 항목 | 점수 | 판정 |
| --- | --- | --- |
| 사이트맵·페이지 커버리지 | 47/47 (100%) | ✅ 완전 |
| 상세 기획 존재 여부 | 47/47 (100%) | ✅ 완전 |
| 기능-API 매핑 | 13/25 완전, 12/25 부분 | ⚠️ 보강 필요 |
| 기능-DB 매핑 | 13/25 완전, 12/25 부분 | ⚠️ 보강 필요 |
| 기능-이벤트 매핑 | 13/25 완전, 12/25 부분 | ⚠️ 보강 필요 |
| Phase 로드맵 정합성 | 10/13 포함, 3/13 누락 | ⚠️ 보강 필요 |
| 구현 코드 | 8/47 (17%) | Phase 1 수준 정상 |

### 5.2 발견된 갭 요약 (우선순위별)

**P0 — 즉시 보강 (기획 문서 갭)**:
1. I.2 매트릭스 이벤트 열 3건 빈 값
2. IV절 Phase 로드맵에 II.11·II.12·II.13 미포함
3. 10.5 식단·10.6 칼로리: API·DB·이벤트 미정의

**P1 — 단기 보강 (상세 명세 부족)**:
4. 10.7 보상시스템: RewardService API·DB 미확정
5. 10.8 음성인식·번역: VoiceProfileService 미확정
6. 10.3 전문가 매칭: 매칭 알고리즘·DB 미정의
7. 10.4 실시간 번역: TranslationService 미확정
8. 10.12 지역 통계: LocationStatsService·테이블 미정의
9. 10.13 레고형: ConceptService·OrganizationService 미확정

**P2 — 중기 보강 (인프라·백엔드)**:
10. 10.9 익명 학습 AI: 익명화 파이프라인·모델 저장소
11. 10.10 데이터 제공: DataProvisionService·B2B 계약

---

## 6. 보강 조치 (본 검증에서 즉시 반영)

아래 P0 항목을 본 보고서 작성과 동시에 반영합니다.

### 6.1 I.2 매트릭스 이벤트 열 보강
- 식단·칼로리: `meal.logged`, `nutrition.analyzed`
- 가족: `family.member.invited`, `family.sharing.updated`
- 설정·알림: `notification.preferences.updated`, `profile.updated`

### 6.2 IV절 Phase 로드맵 보강
- Phase 1: II.12 기본(88-dim 파이프라인)
- Phase 2: II.11 기본(텍스트 AI 비서), II.12 확장(448-dim)
- Phase 3: II.12 확장(896-dim), II.13(1st-party 카트리지 스토어)
- Phase 4: II.13(SDK 공개·서드파티)
- Phase 5: II.11 고도화(음성), II.12 궁극(1792-dim), II.13(전면 오픈·AI 모델 마켓)

### 6.3 10.5 식단·10.6 칼로리 API·DB·이벤트 정의
- DietService RPC: LogMeal, ListMeals, GetDailyNutritionSummary, DeleteMeal
- AiInferenceService RPC: AnalyzeFoodImage (이미지→영양소 벡터)
- DB: diet_logs(id, user_id, meal_type, items_json, calories, nutrients_json, image_url, logged_at), daily_nutrition_summary(user_id, date, total_calories, macros_json)
- 이벤트: `meal.logged`, `nutrition.daily.summarized`

---

## 7. 결론

만파식 생태계의 전체 사이트맵(47+ 페이지·화면)에 대한 **상세 페이지 기획은 100% 수립**되어 있습니다. 핵심 13개 기능 영역은 기능-페이지-API-DB-이벤트 4중 매핑이 완전합니다.

12개 부가 기능 영역(10.3~10.13 일부)은 API·DB·이벤트 명세가 "제안" 수준이므로, Phase별 구현 시점에 맞춰 상세 명세를 확정해야 합니다. 이 중 P0 3건은 본 보고서에서 즉시 보강 반영했습니다.

**기획 완성도**: 전체 시스템 구축에 필요한 기획·구현 계획이 **체계적으로 수립**되어 있으며, 사이트맵·스토리보드 대비 **누락된 페이지는 없습니다**.
