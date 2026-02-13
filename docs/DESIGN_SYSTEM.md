# ManPaSik Design System: "The Depths of Sanggam"

> **Code Name**: NadoBanana Pro Style
> **Concept**: 한국적 정서(K-Sentiment)와 첨단 기술(High-Tech)의 조화
> **Version**: v2.0 (2026-02-14 통합 팔레트 적용)

---

## 🎨 Color Palette: "The Depths of Sanggam"

모든 플랫폼(App, Web, OS)은 아래 컬러 시스템을 공유한다.

### Background: "Deep Sea Navy" (심해)
- **Color**: `#0A192F`
- **Meaning**: 심해의 깊이감, 몰입, 신중함
- **Usage**: 전체 배경, 다크 모드 기본 배경

### Primary: "Sanggam Gold" (상감 금색)
- **Color**: `#D4AF37`
- **Meaning**: 금속 상감 기법의 정교함, 프리미엄 가치, 하이라이트
- **Usage**: 헤드라인, 버튼, 테두리, 강조 포인트

### Secondary: "Wave Cyan" (파동 청록)
- **Color**: `#64FFDA` (Web) / `#00E5FF` (App)
- **Meaning**: 데이터의 흐름, 파동, 에너지, 생명력
- **Usage**: 데이터 라벨, 활성 상태, 진행률 표시

### Surface: "Glass Navy" (유리 네이비)
- **Color**: `#112240` (Solid) / `rgba(26, 35, 126, 0.6)` (Glass)
- **Meaning**: 글래스모피즘을 통한 공간감, 현대적 세련미
- **Usage**: 카드, 패널, 모달 배경 (Blur 20px)

### Alert: "Dancheong Red" (단청 적색)
- **Color**: `#D32F2F` (Web) / `#FF4D4D` (App)
- **Meaning**: 생명력, 경고, 역동성
- **Usage**: 알림, 위험 수치, 에러 상태

### Text: "Hanji White" (한지 백색)
- **Color**: `#FAFAFA`
- **Meaning**: 순수, 여백의 미
- **Usage**: 본문 텍스트, 라이트 모드 배경

### Extra Dark: "Ink Black" (먹색)
- **Color**: `#020617`
- **Meaning**: 깊이 있는 지식, 신뢰
- **Usage**: 최상위 배경 그라데이션 하단

---

## ✍️ Typography

### Korean
- **Headings**: `Gowun Batang` (고운바탕) — 정갈함, 전통적, 감성적
- **Body**: `Noto Sans KR` (본고딕) — 현대적, 가독성, 과학적

### English
- **Display**: `Outfit` — 현대적, 테크니컬
- **Brand**: `Playfair Display` — 품격, 프리미엄

### Code / Data
- **Mono**: `JetBrains Mono` — 데이터 신뢰성

---

## 🌊 Dynamic Interaction

### 1. Wave Ripple Effect (물결 파동)
- Legend of Manpasikjeok (파도를 잠재우는 피리) 설화에서 차용
- **Flutter**: `WaveRipplePainter` — 동심원이 중심에서 바깥으로 확산, Wave Cyan → Sanggam Gold 그라데이션
- **Web**: `@keyframes wave-ripple` — box-shadow 기반 동심원 확산
- **적용 화면**: Splash 배경, 디바이스 연결 중 상태

### 2. Breathing Animation (호흡 효과)
- 측정 진행 시 화면 전체가 호흡하듯 미세하게 움직임
- **Flutter**: `BreathingOverlay` — scale(0.98↔1.02) + opacity(0.85↔1.0) 반복
- **Web**: `@keyframes breathing` — scale + opacity 루프
- **적용 화면**: 측정 화면 (measuring 상태), AI 인사이트 패널

### 3. Wave Painter (파동 안정화)
- 측정 진행률에 따라 사인파 진폭이 감소하여 직선으로 수렴
- "세상의 파동을 잠재운다" 철학 시각화
- **Flutter**: `WavePainter` — 네온 글로우 효과의 사인파 라인
- **단계 텍스트**: `파동 안정화 중...` → `분석 중...` → `측정 완료`

### 4. Sanggam Glow Line (상감 금선)
- 헤더/섹션 구분선에 금색 그라데이션 라인
- **Web**: `@keyframes glow-line` — 좌→우 이동하는 금색 광택

### 5. Data Flow (데이터 흐름)
- 실시간 측정 데이터가 흐르는 듯한 라인 그래프
- 네온 글로우 효과로 첨단 과학 느낌 강조

---

## 📱 UI Components

### Sanggam Container (상감 컨테이너)
- **Border**: 1px Solid `#D4AF37` + 2px Inset Gradient (Transparent → Gold → Transparent)
- **Background**: Linear Gradient `#0A192F` → `#112240`
- **Shadow**: Outer `black.withOpacity(0.5), blur: 10` + Inner `gold.withOpacity(0.1), blur: 5`

### Sanggam Panel (Web 글래스 패널)
- **Background**: `rgba(26, 35, 126, 0.4)` + `backdrop-filter: blur(20px)`
- **Border**: `1px solid rgba(212, 175, 55, 0.3)`
- **Hover**: Border opacity 0.6, Shadow 강화

### Cards
- 반투명한 유리 질감 + 한지 텍스처 (Blur + Noise)

### Buttons
- 둥근 모서리 (기와 곡선 형상화)
- Sanggam Gold 테두리, 호버 시 반전

### Shadows
- 은은하고 깊이 있는 그림자 (Soft Ambient)
