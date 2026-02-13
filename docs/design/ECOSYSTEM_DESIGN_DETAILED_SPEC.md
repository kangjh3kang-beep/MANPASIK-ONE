# 만파식 월드(MANPASIK World) 상세 디자인 명세서

**문서번호**: MPK-DESIGN-DETAILED-v1.0
**작성일**: 2026-02-13
**상태**: Draft
**참조**: MPK-ECO-PLAN-v1.1, MPK-DESIGN-MASTER-v1.1, BRAND_GUIDELINE

---

## 1. 개요 (Overview)

본 문서는 **[만파식 월드 통합 디자인 마스터플랜](ECOSYSTEM_DESIGN_MASTER_PLAN.md)**의 하위 실행 문서로, 시스템 기획서(`sitemap`, `cartridge-spec`, `storyboard`)에 정의된 기능적 요구사항을 **구체적인 UI/UX 디자인 언어**로 번역한 명세서이다.
특히 **'오타 없는 완벽한 한국어'** 적용을 위해, 모든 디자인 컴포넌트의 텍스트는 본 문서의 **[5. UI 텍스트 표준 가이드](#5-ui-텍스트-표준-가이드-korean-text-standards)**를 원본으로 한다.

---

## 2. 모듈별 상세 디자인 전략 (Module Specifics)

### 2.1 App Module: MANPASIK Measurement System (Hands)

**핵심 목표**: 현장의 "Action & Control"에 최적화된 직관적이고 반응성 높은 인터페이스.

#### A. 내비게이션 구조 (Navigation) - `sitemap.md` 기반
- **구조**: Bottom Navigation Bar (5 Tabs)
  1.  **Home** (`/`): 대시보드 요약 (Smart Cards)
  2.  **Measurements** (`/measure`): 측정 및 결과 (Floating Action Button 강조)
  3.  **Data Hub** (`/data`): 건강 데이터 시각화 (Charts)
  4.  **AI Coach** (`/coach`): 대화형 인터페이스 (Chat UI)
  5.  **Market** (`/market`): 카트리지 구매 (Grid Layout)
- **특이사항**: `Measurement` 탭은 중앙에 돌출된 **'Hexagon Button'** 형태로 배치하여 접근성 및 상징성 강화. (MANPASIK의 6각형 화학 구조 은유)

#### B. 측정 프로세스 UI (Measurement Flow) - `storyboard-first-measurement.md` 기반
- **Step 1: 연결 & 인식 (Connection)**
  - **Interaction**: 리더기 근접 시 화면 하단에서 `Glassmorphism Sheet`가 부드럽게 올라오는 애니메이션.
  - **Feedback**: 연결 성공 시 기분 좋은 `Haptic Feedback (Light Impact)` + 리더기 아이콘의 가장자리가 `Cyan` 컬러로 발광.
- **Step 2: 카트리지 인식 (Cartridge Recognition)**
  - **Dynamic UI**: NFC 태그 정보를 읽어 (`cartridge-system-spec.md` 참조) 카트리지 종류에 따라 테마 컬러 즉시 변경.
    - *HealthBiomarker*: **Deep Sea Navy + Sanggam Gold**
    - *Environmental*: **Forest Green + Silver**
    - *FoodSafety*: **Warm Orange + Bronze**
  - **Infinite Scroll**: 2-Byte 코드로 무한 확장되는 카트리지 리스트를 `Virtual Scroll`로 처리하되, 자주 쓰는 항목은 상단 `Magnetic Area`에 고정.
- **Step 3: 측정 중 (Processing)**
  - **Visual**: 화면 중앙에 **'MANPASIK Wave'** 애니메이션 재생. 노이즈가 점차 줄어들며 직선(Flat Line)으로 수렴하는 과정을 시각화하여 "세상의 파동을 잠재운다"는 철학 전달.
  - **Progress**: 단순 `%` 표시 대신, `Stabilizing...` -> `Analyzing...` -> `Complete` 단계별 텍스트 모핑.

#### C. 결과 카드 (Result Card)
- **Design**: **SanggamContainer** 적용 (Deep Sea 배경 + 금박 인레이).
- **Data Vis**:
  - **수치**: `Oswald` 또는 `JetBrains Mono` 폰트로 크게 강조.
  - **Range**: 가로 바(Bar) 차트 위에 '나의 위치'를 다이아몬드(◆) 심볼로 표시.
  - **Ai Insight**: 하단에 `Typewriter Effect`로 AI 코멘트 한 줄 요약 출력.

### 2.2 Web Module: MANPASIK Health Management Lab (Brain)

**핵심 목표**: 방대한 데이터의 "Broad & Deep" 심층 분석 및 관제.

#### A. 대시보드 레이아웃 (Dashboard Layout)
- **Grid System**: 12-Column Fluid Grid. 해상도(Ultra-Wide)에 따라 패널 자동 재배치.
- **Glass Panel**: 모든 모듈은 `BackdropFilter(blur: 20px)`가 적용된 반투명 패널 사용. 배경의 심해(Deep Sea) 텍스처가 은은하게 비침.

#### B. 데이터 시각화 (Data Visualization)
- **896차원 핑거프린트**:
  - **WebGL**: `Three.js` 또는 `R3F`를 활용한 3D 구체(Sphere) 형태의 데이터 군집 시각화.
  - **Interaction**: 마우스 드래그로 회전, 줌인/아웃. 특정 포인트 클릭 시 상세 성분 정보 팝업.
- **Heatmap**: 시간(X축) x 차원(Y축) 히트맵으로 건강 트렌드 변화 추적. `Viridis` 또는 `Magma` 컬러맵 대신 `Manpasik Custom Gradient` (Navy -> Cyan -> Gold) 사용.

### 2.3 OS Module: MANPASIK Core OS (Heart)

**핵심 목표**: "Minimal & Stable" 시스템 상태의 명확한 전달.

#### A. 임베디드 디스플레이 (160x80 OLED)
- **Typography**: 가독성을 극대화한 `Pixel Font` (Monospaced).
- **Status Line**: 상단 1px 라인으로 상태 표시 (초록: 정상, 노랑: 경고, 빨강: 오류).
- **Animation**:
  - **Boot**: 로고가 왼쪽에서 오른쪽으로 그려지는 라인 드로잉.
  - **Active**: 심장 박동(Heartbeat) 리듬의 미세한 점멸.

---

## 3. 디자인 시스템 상세 (Design System Specs)

### 3.1 Sanggam Container (상감 컨테이너)
- **Border**:
  - Outer: 1px Solid `#D4AF37` (Sanggam Gold)
  - Inner: 2px Inset Gradient (Transparent -> Gold -> Transparent)
- **Background**:
  - Linear Gradient: TopLeft(`#0A192F`) -> BottomRight(`#112240`)
- **Shadow**:
  - Outer: `BoxShadow(color: black.withOpacity(0.5), blur: 10, offset: (0, 4))`
  - Inner: `BoxShadow(color: gold.withOpacity(0.1), blur: 5, spread: 1)`

### 3.2 Typography Scale (App 기준)
- **Display Large**: `Gowun Batang`, 32sp, Bold (헤드라인)
- **Title Medium**: `Noto Sans KR`, 18sp, SemiBold (카드 제목)
- **Body Large**: `Noto Sans KR`, 16sp, Regular (본문)
- **Label Small**: `JetBrains Mono`, 12sp, Medium (데이터, 캡션)

### 3.3 Iconography
- **Style**: **Thin Lined** (1.5px stroke), **Geometric**.
- **Metaphor**:
  - 측정: 육각형, 파동, 스코프
  - AI: 별(Star), 뉴런 노드
  - 데이터: 레이어, 큐브

---

## 4. 구현 가이드 (Implementation Guide)

### 4.1 Flutter (App)
- `lib/core/theme/sanggam_theme.dart`: `ThemeData` 확장 정의.
- `lib/shared/widgets/sanggam_container.dart`: `CustomPainter`로 상감 효과 구현.
- `lib/features/measure/animations/wave_painter.dart`: 파동 애니메이션 구현.

### 4.2 Next.js (Web)
- `tailwind.config.js`: 커스텀 컬러(`deep-sea`, `sanggam-gold`) 및 유틸리티 확장.
- `components/ui/GlassPanel.tsx`: `backdrop-filter` 및 Border 스타일링 컴포넌트화.

### 4.3 Rust (OS)
- `slint` 프레임워크 활용 시 `.slint` 파일에 스타일 변수 정의.
- 상태 LED 제어를 위한 `hal` 모듈 인터페이스 정의.

---

## 5. UI 텍스트 표준 가이드 (Korean Text Standards)

사용자 인터페이스에 노출되는 모든 한글 텍스트는 아래 표준을 따른다. 개발 및 디자인 시 복사/붙여넣기하여 오타를 방지한다.

### 5.1 공통 (Common)
- **Brand**: MANPASIK (영문), 만파식 (국문)
- **Values**: 초정밀 (High-Precision), 차동 계측 (Differential Measurement), 파동 (Wave)

### 5.2 App (Measurement System)
| 화면 | UI 요소 | 정확한 표기 (Korean) | 비고 |
| :--- | :--- | :--- | :--- |
| **Splash** | Title | **MANPASIK** | 영문 로고 사용 |
| | Subtitle | **초정밀 차동 계측 시스템** | 띄어쓰기 준수 |
| **Home** | Greeting | **안녕하세요, 연구원님.** | '님' 호칭 통일 |
| | Status | **현재 시스템 상태: 정상** | '안전', '양호' 대신 '정상' |
| **Measure** | Button | **측정 시작** | 'Start (X)' |
| | Progress | **파동 안정화 중...** | 말줄임표 3개 |
| | Result | **측정 완료** | |

### 5.3 Web (Health Management Lab)
| 화면 | UI 요소 | 정확한 표기 (Korean) | 비고 |
| :--- | :--- | :--- | :--- |
| **Login** | Title | **만파식 건강관리연구소** | |
| **Dash** | Header | **통합 관제 대시보드** | |
| | Label | **실시간 데이터 분석** | |

### 5.4 OS (Core)
| 화면 | UI 요소 | 정확한 표기 (Korean) | 비고 |
| :--- | :--- | :--- | :--- |
| **Boot** | Text | **시스템 초기화 중...** | 폰트 제약으로 영문 권장되나 한글 지원 시 사용 |
| **Ready** | Status | **준비됨** | |
| **Error** | Alert | **센서 연결 확인 필요** | 명확한 행동 유도 |
