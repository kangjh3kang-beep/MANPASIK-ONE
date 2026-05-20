# ManPaSik UI/스토리라인 재점검 보고서 (2026-02-23)

## 1) 점검 목적
- 모바일/데스크톱에서 네비게이션 일관성 저하 문제를 제거한다.
- 상단/하단 프레임 노출 불안정(흰 박스, 중복 네비)을 제거한다.
- 의료 기능을 단순 통화 진입형에서 실제 진료 여정형(기관 선택-예약-진료-약국-결과)으로 복원한다.

## 2) 실제 원인 진단
- 전역 크롬(상단/하단)과 화면별 로컬 크롬이 동시에 렌더되어 중복/겹침 발생.
- 화면별 뒤로가기 처리 방식이 통일되지 않아 스택 없는 상태에서 이탈 실패.
- 모바일에서 좌측 네비가 사라지는 뷰포트에서 대체 진입점 부족.
- 의료 허브가 전화/통화 중심으로 축소되어 핵심 선행 작업(기관/약국 선택, 데이터 공유)이 분리됨.

## 3) 반영 완료 사항
- 공통 내비 fallback 추가: pop 불가 시 홈으로 이동.
- 경로 정책 기반 전역 크롬 제어로 중복 렌더 방지.
- 모바일 Drawer 전체 메뉴 추가.
- 홈/측정 화면에 전체 메뉴 바텀시트 추가.
- 의료 허브를 진료 여정형 액션 카드 중심으로 재구성.
- 의료 하위 화면 뒤로가기 로직 통일.
- /medical/video-call 경로의 글로벌 크롬 숨김 정책/테스트 추가.

## 4) 사용자 여정(권장 IA)
- 홈
- 데이터 허브
- 측정
- 의료
- 마켓
- 커뮤니티
- 가족
- 설정

의료 상세 여정
1. 의료기관 선택
2. 진료 예약
3. 화상진료 입장
4. 처방전/약국 전송
5. 결과/기록 확인

## 5) 연구/가이드 근거
- HHS Telehealth: 운영 워크플로우(코디네이터 지정, 예약-방문-사후 단계 분리) 권고
  - https://telehealth.hhs.gov/providers/telehealth-implementation-and-workflow-planning/creating-workflows
- HHS Telehealth: 서비스 준비/운영 체크리스트
  - https://telehealth.hhs.gov/providers/best-practice-guides/telehealth-implementation-playbook
- WHO: 디지털 헬스케어의 접근성/대기시간 개선 잠재 효과 및 안전한 도입 필요성
  - https://www.who.int/publications/i/item/9789241550505
- JAMA Network Open(2022): 응급/퇴원 영역에서 원격진료 사용과 운영 연계 분석
  - https://jamanetwork.com/journals/jamanetworkopen/fullarticle/2789923
- PubMed(2021, BMJ Open): 영상 기반 진료가 전화 대비 동등/우수한 임상결과 보고(체계적 문헌고찰)
  - https://pubmed.ncbi.nlm.nih.gov/34711580/
- HIRA(심평원): 비대면진료 시범사업 기관/약국 참여 구조 및 안내
  - https://www.hira.or.kr/dummy.do?pgmid=HIRAA020041000000
- ONC(미국): 환자 데이터 접근성(API) 확대는 환자 중심 디지털 경험 핵심 기반
  - https://www.healthit.gov/data/quickstats/quickstats/individuals-whose-health-care-providers-maintained-electronic-health-information
- WCAG 2.2: 고정 헤더/푸터가 포커스 콘텐츠를 가리지 않도록 요구(접근성)
  - https://www.w3.org/WAI/WCAG22/Understanding/focus-not-obscured-minimum
- Android 적응형 내비게이션: 화면 크기별 Bottom/NavigationRail/Drawer 권장
  - https://developer.android.com/develop/ui/views/layout/responsive-adaptive-design-with-views

## 6) 남은 고도화 백로그
- 의료기관-의사-예약 슬롯-결제-접속 토큰 발급까지 단일 트랜잭션 플로우 확장.
- 약국 선택 시 재고/조제 가능 시간/거리 기반 추천 추가.
- 진료 전 공유 데이터 묶음(측정값, 추세, 복약, 알레르기) 템플릿화.
- 진료 중 실시간 데이터 패널(필수 지표 3~5개) 고정.
- 진료 종료 후 결과/처방/재방문 예약을 하나의 타임라인으로 통합.

## 7) 검증 기록
- flutter analyze (변경 파일 대상): 통과
- flutter test test/core/router/bottom_nav_visibility_test.dart: 통과
- flutter build linux --debug: 통과
