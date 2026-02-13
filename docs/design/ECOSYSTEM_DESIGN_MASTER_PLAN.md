# 만파식 월드(MANPASIK World) 통합 디자인 마스터플랜

**문서번호**: MPK-DESIGN-MASTER-v1.1
**작성일**: 2026-02-13
**상태**: Revised

---

## 1. 개요 (Overview)

본 문서는 **만파식 월드(MANPASIK World)**라는 거대한 통합 생태계를 구축하기 위한 디자인 가이드라인이다. `~/Manpasik` 프로젝트는 이 세계관을 구현하는 메인 실체로서, 하드웨어(OS), 모바일(App), 클라우드(Web)가 유기적으로 연결된 **초연결 AI 계측 시스템**을 지향한다.

### 1.1 핵심 철학: "Korean Futuristic" (한국적 미래주의)
- **전통(Heritage)**: Manpasikjeok(萬波息笛) 설화의 의미 계승. 거친 파도(노이즈)를 잠재우고 평온(데이터)을 찾음.
- **기술(Tech)**: 인공지능과 초정밀 계측 기술의 융합.
- **시각화(Visualization)**: 심해(Deep Sea)의 고요함 속에 빛나는 상감(Sanggam Gold)의 정교함.

---

## 2. 생태계 구조 및 네이밍 전략 (Ecosystem Structure)

'만파식 월드'라는 최상위 브랜드 아래, 각 모듈은 고유한 역할과 네이밍을 갖는다.

| 플랫폼 모듈 | 코드명 | 정식 명칭 | 역할 (Role) | 디자인 키워드 |
| :--- | :--- | :--- | :--- | :--- |
| **Ecosystem** | `~/Manpasik` | **MANPASIK World** | **세계관 (Universe)**<br>- Web, App, OS를 포괄하는 통합 시스템<br>- 만파식 생태계의 본체 | **Infinity & Harmony**<br>- 무한한 확장성<br>- 전체의 조화 |
| **Web Module** | `frontend/web-app` | **MANPASIK Health Management Lab** | **두뇌 (Brain)**<br>- 데이터 심층 분석<br>- AI 모델 학습 및 관리<br>- 중앙 관제 대시보드 | **Broad & Deep**<br>- 광활한 데이터 시각화<br>- 복잡한 정보를 정갈하게 정리 |
| **App Module** | `frontend/flutter-app` | **MANPASIK Measurement System** | **손발 (Hands)**<br>- 현장 디바이스 제어<br>- 실시간 측정 및 수집<br>- 즉각적인 피드백 | **Action & Control**<br>- 직관적인 컨트롤 인터페이스<br>- 햅틱 및 동적 반응 |
| **OS Module** | `rust-core` | **MANPASIK Core OS** | **심장 (Heart)**<br>- 임베디드 리더기 구동<br>- Rust 기반 초고속 연산<br>- 하드웨어 추상화 계층 | **Minimal & Stable**<br>- 극도로 절제된 상태 표시<br>- 시스템 안정성 시각화 |

---

## 3. 통합 디자인 언어 (Unified Design Language)

### 3.1 컬러 팔레트: "The Depths of Sanggam"
모든 플랫폼은 아래의 컬러 시스템을 공유하여 통일감을 부여한다.

- **Background**: `Deep Sea Navy` (#0A192F)
  - 심해의 깊이감, 몰입, 신중함. (기존 NadoBanana Black에서 진화)
- **Primary**: `Sanggam Gold` (#D4AF37)
  - 금속 상감 기법의 정교함, 프리미엄 가치, 하이라이트.
- **Secondary**: `Wave Cyan` (#64FFDA)
  - 데이터의 흐름, 파동, 에너지, 생명력.
- **Surface**: `Glass Navy` (Blur 20px, Opacity 60%)
  - 글래스모피즘을 통한 공간감, 현대적 세련미.

### 3.2 그래픽 모티프: "Sanggam Line (상감 라인)"
- **정의**: 모든 컨테이너와 주요 구획은 1px~2px 두께의 금색/청록색 그라데이션 라인으로 마감한다.
- **의미**: 칠기나 도자기에 문양을 파 넣듯, 데이터와 기능을 정교하게 새겨 넣음.
- **구현**:
  - **Flutter**: `CustomPainter`를 이용한 `SanggamContainer`.
  - **Web**: Tailwind `border-image` 및 `box-shadow` 활용.
  - **OS**: OLED 디스플레이의 픽셀 단위 발광 라인.

### 3.3 타이포그래피
- **Korean**: `Noto Sans KR` (본문) + `Gowun Batang` (헤드라인, 감성적 강조)
- **English**: `Outfit` (현대적, 테크니컬) + `Playfair Display` (브랜드, 품격)
- **Code/Data**: `JetBrains Mono` (데이터 신뢰성)

---

## 4. 플랫폼별 상세 가이드 (Platform Specifics)

### 4.1 MANPASIK R&D Lab (Web)
- **레이아웃**: 와이드 스크린을 활용한 다단 그리드.
- **데이터 시각화**: 3D 구체(Sphere), 복잡한 파동 그래프, 히트맵 등 고밀도 정보 표현.
- **인터랙션**: 마우스 호버 시 글래스 표면이 빛나는 `Luster Effect`.

### 4.2 MANPASIK Measurement System (App)
- **레이아웃**: 한 손 조작에 최적화된 하단 내비게이션 및 카드 리스트.
- **컨트롤**: 물리 버튼을 누르는 듯한 `Haptic Touch` 및 `Neumorphism` 일부 차용.
- **피드백**: 측정 진행 시 화면 전체가 호흡하듯 움직이는 `Breathing Animation`.

### 4.3 MANPASIK Core OS (Embedded)
- **레이아웃**: 해상도 제약(예: 160x80)을 고려한 픽셀 퍼펙트 아이콘.
- **상태 표시**: LED 인디케이터와 연동되는 심플한 라인 애니메이션.
- **부팅**: 만파식적 피리 소리를 시각화한 사운드 웨이브 로고.

---

## 5. 구현 로드맵 (Implementation Roadmap)

1.  **Phase 1: Identity Establish (완료)**
    - App/Web 공통 컬러 및 로고 브랜딩 적용.
    - `SanggamContainer` 프로토타입 구현 및 검증.
2.  **Phase 2: System Optimization (진행 중)**
    - App: Measurement System 전용 UI/UX 고도화 (워크쓰루 반영).
    - Web: R&D Lab 전용 데이터 시각화 컴포넌트 개발.
3.  **Phase 3: OS Integration (예정)**
    - Rust Core 기반 임베디드 GUI (Slint/Flutter Embedder) 설계.
    - 하드웨어-소프트웨어 인터랙션 프로토콜 정의.

---

**결론**: 만파식 생태계는 단순한 기능의 집합이 아니라, **'한국적 미학이 깃든 첨단 기술의 결정체'**로서 사용자에게 일관되고 깊이 있는 경험을 제공해야 한다.
