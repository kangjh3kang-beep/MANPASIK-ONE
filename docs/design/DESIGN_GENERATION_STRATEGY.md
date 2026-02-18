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

---

## 5. 디자인 토큰 및 컴포넌트 라이브러리

### 5.1 디자인 토큰 (Design Tokens)

상감 디자인 시스템의 모든 시각적 속성을 토큰으로 관리한다.

| 카테고리 | 토큰 | 값 | 용도 |
|---------|------|-----|------|
| Color/Primary | `--sanggam-deep-sea-navy` | #0A192F | 배경, 텍스트 |
| Color/Secondary | `--sanggam-gold` | #D4AF37 | 강조, CTA |
| Color/Accent | `--sanggam-wave-cyan-web` | #64FFDA | Web 인터랙션 |
| Color/Accent | `--sanggam-wave-cyan-app` | #00E5FF | App 인터랙션 |
| Spacing | `--space-xs` ~ `--space-3xl` | 4px ~ 64px | 여백 체계 |
| Radius | `--radius-sm` ~ `--radius-full` | 8px ~ 9999px | 모서리 곡률 |
| Shadow | `--shadow-sanggam` | 0 4px 32px rgba(212,175,55,0.15) | 상감 컨테이너 |
| Motion | `--duration-fast` ~ `--duration-slow` | 100ms ~ 3000ms | 애니메이션 |
| Typography | `--font-display` | Gowun Batang | 제목용 서체 |
| Typography | `--font-body` | Noto Sans KR | 본문용 서체 |

### 5.2 핵심 컴포넌트 라이브러리

| 컴포넌트 | 플랫폼 | 상태 | 비고 |
|---------|--------|------|------|
| SanggamContainer | Flutter/Web | 구현 완료 | 금선 보더 + 그라데이션 |
| PrimaryButton | Flutter/Web | 구현 완료 | Wave Cyan CTA |
| MeasurementCard | Flutter | 구현 완료 | 측정 결과 카드 |
| BiomarkerChart | Flutter/Web | 설계 완료 | 트렌드 차트 |
| HexagonButton | Flutter | 설계 완료 | 측정 시작 |
| WaveAnimation | Flutter | 설계 완료 | 파동 리플 효과 |
| StatusIndicator | All | 구현 완료 | Normal/Caution/Alert |
| NavBottomBar | Flutter | 구현 완료 | 5-탭 네비게이션 |

---

## 6. 플랫폼별 디자인 생성 우선순위

### 6.1 Phase 1 (MVP) — Month 1-4

| 순위 | 화면 | 플랫폼 | 산출물 |
|------|------|--------|--------|
| P1 | 온보딩 + 회원가입 | App | 모바일 스크린 5장 |
| P1 | 측정 플로우 (BLE 연결→측정→결과) | App | 모바일 스크린 6장 |
| P1 | 홈 대시보드 | App | 모바일 스크린 1장 |
| P2 | 디바이스 페어링 | App | 모바일 스크린 3장 |
| P2 | 설정/프로필 | App | 모바일 스크린 2장 |

### 6.2 Phase 2 (Core) — Month 5-8

| 순위 | 화면 | 플랫폼 | 산출물 |
|------|------|--------|--------|
| P1 | 건강 대시보드 (Health Command Center) | Web | 데스크톱 스크린 2장 |
| P1 | AI 코치 대화 | App | 모바일 스크린 3장 |
| P2 | 카트리지 마켓 | App + Web | 크로스플랫폼 4장 |
| P2 | 가족 건강 관리 | App | 모바일 스크린 4장 |
| P3 | 커뮤니티 | App | 모바일 스크린 3장 |

### 6.3 Phase 3 (Advanced) — Month 9-12

| 순위 | 화면 | 플랫폼 | 산출물 |
|------|------|--------|--------|
| P1 | 비대면 진료 (WebRTC) | App + Web | 크로스플랫폼 5장 |
| P2 | 리더기 OS UI (Boot/Idle/Measuring) | OS (OLED) | 픽셀 스크린 4장 |
| P3 | 관리자 대시보드 | Web | 데스크톱 스크린 3장 |
| P3 | 3D 핑거프린트 시각화 | Web | 인터랙티브 컴포넌트 |

---

## 7. 디자인 QA 프로세스

### 7.1 품질 검증 체크리스트

| # | 항목 | 검증 기준 | 도구 |
|---|------|---------|------|
| DQ01 | 색상 일관성 | 상감 팔레트 100% 준수 | Design Lint |
| DQ02 | 접근성 (WCAG AA) | 명도비 4.5:1 이상 | Stark / axe |
| DQ03 | 타이포그래피 | 스케일 8단계 준수 | Style Dictionary |
| DQ04 | 간격/여백 | 4px 그리드 기반 | Figma Auto Layout |
| DQ05 | 다크 모드 | Deep Sea Navy 기반 반전 | 수동 검증 |
| DQ06 | RTL 지원 | 아랍어/히브리어 미러링 | Flutter RTL 테스트 |
| DQ07 | 반응형 | 5개 브레이크포인트 | 브라우저 DevTools |
| DQ08 | 모션 | 100ms-3000ms 범위 | 프로토타입 검증 |

### 7.2 디자인-코드 핸드오프

```
[Figma 디자인] → [디자인 토큰 Export]
                        │
                   JSON/YAML
                        │
            ┌───────────┼───────────┐
            ▼           ▼           ▼
      [Flutter]    [Web CSS]    [OS OLED]
     ThemeData    Tailwind v4   Slint .60
```

- **자동 동기화**: Figma Tokens → Style Dictionary → 플랫폼별 코드 생성
- **검증**: CI에서 디자인 토큰 diff 감지 → 변경 시 리뷰 필수
