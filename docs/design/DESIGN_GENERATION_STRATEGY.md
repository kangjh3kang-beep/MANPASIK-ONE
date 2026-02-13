# 만파식 월드(MANPASIK World) 디자인 생성 및 기획 전략서

**문서번호**: MPK-DESIGN-GEN-v1.0
**작성일**: 2026-02-13
**작성자**: Manpasik Professional Designer

---

## 1. 개요 (Overview)

본 문서는 기 수립된 `ECOSYSTEM_DESIGN_MASTER_PLAN` 및 `DETAILED_SPEC`을 바탕으로, **실제 눈에 보이는 디자인(Visual Design)**을 생성하고 구체화하기 위한 실행 계획이다. 

**목표**: 단순한 기획서를 넘어, 개발자가 보고 코드로 옮길 수 있는 수준의 **시각적 산출물(Mockups)**과 **상세 명세(Specs)**를 병렬적으로 생산한다.

---

## 2. 플랫폼별 디자인 생성 대상 (Target Scope)

### 2.1 Web Module: MANPASIK Health Management Lab (Brain)
**컨셉**: 광활한 정보의 바다 (Broad & Deep)
- **[Priority 1] 통합 건강 관제 대시보드 (Health Command Center)**
  - 12-Grid 레이아웃 내 896차원 데이터 시각화 패널 배치.
  - 전 세계 리더기 연결 현황을 보여주는 3D Globe 위젯.
- **[Priority 2] 분석 상세 페이지 (Analysis Detail)**
  - 3D Sphere 핑거프린트 회전 뷰.
  - 시계열 히트맵 (Deep Sea Gradient 적용).

### 2.2 App Module: MANPASIK Measurement System (Hands)
**컨셉**: 즉각적인 제어와 반응 (Action & Control)
- **[Priority 1] 측정 컨트롤 패널 (Measurement Controller)**
  - 중앙 Hexagon Button과 파동 애니메이션(Wave) 결합.
  - SanggamContainer 기반의 실시간 상태 카드.
- **[Priority 2] 결과 리포트 카드 (Result Card)**
  - 모바일 화면에 최적화된 세로형 상감 프레임.
  - 직관적인 등급 표시 (Inlay Gold Icon).

### 2.3 OS Module: MANPASIK Core OS (Heart)
**컨셉**: 절제된 안정성 (Minimal & Stable)
- **[Priority 1] 부팅 및 대기 화면 (Boot & Idle)**
  - 160x80px 해상도 제한을 고려한 픽셀 퍼펙트 로고.
  - 시스템 상태(배터리, 네트워크)를 나타내는 1px 라인 인디케이터.

---

## 3. 병렬 진행 프로세스 (Parallel Process)

효율적인 진행을 위해 다음 3단계 사이클을 반복한다.

1.  **Define (정의)**: 각 화면에 필요한 데이터와 기능 명세 확인.
2.  **Visualize (생성)**: `generate_image` 도구를 활용하여 고품질 컨셉 이미지 생성.
3.  **Specify (명세)**: 생성된 이미지를 바탕으로 컬러 코드, 여백, 폰트 사이즈 등 퍼블리싱 스펙 확정.

---

## 4. 산출물 관리
- 모든 디자인 이미지는 `docs/design/assets/` 경로에 저장 및 관리.
- 확정된 디자인은 `ECOSYSTEM_DESIGN_DETAILED_SPEC.md`에 링크로 삽입하여 문서화.
