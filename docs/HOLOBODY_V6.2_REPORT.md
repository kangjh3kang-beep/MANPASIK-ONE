# HoloBody v6.2 작업 보고서

**날짜**: 2026-02-19
**범위**: 홀로그램 품질 고도화 4종 + 대시보드 크래시 해결 2건

---

## 1. 작업 배경

v6.0/v6.1에서 13단계 렌더링 파이프라인과 인체 비율을 수정했으나:
- 인체 홀로그램이 여전히 조잡한 파티클 노이즈 수준
- `child.hasSize` 어설션 에러로 앱 크래시 발생 (Lost connection to device)
- 참조 영상(의료급 3D 스캔 홀로그램)과 비교 시 큰 품질 격차

---

## 2. 구현 내역

### 2.1 Catmull-Rom 스플라인 실루엣
**파일**: `holo_body.dart` — `_drawBodyContourGlow`

직선으로 연결하던 인체 외곽선을 Catmull-Rom 스플라인으로 변환하여 매끄러운 곡선 실루엣 생성.

- Catmull-Rom → Cubic Bezier 변환 (tension 1/6)
- 이중 레이어 렌더링:
  - 외곽 글로우: 6px stroke, MaskFilter.blur(12)
  - 내부 실루엣: 1.5px stroke, MaskFilter.blur(2)

### 2.2 HSL 깊이 색상 그라디언트
**파일**: `holo_body.dart` — `_drawParticles`

단색 cyan 파티클을 Z축 깊이에 따라 HSL 색공간에서 자동 변환.

| 속성 | 앞쪽 (가까운) | 뒤쪽 (먼) |
|------|:----------:|:--------:|
| 밝기 | 0.8 | 0.3 |
| 채도 | 0.9 | 0.4 |
| blur | 1.2 | 1.2 |

스캔 영역 파티클: blur 3.0, 크기 x1.3

### 2.3 메시 투명화 + 볼류메트릭 글로우
**파일**: `holo_body.dart` — `_drawTriangleMesh`

메시를 미세화하여 파티클이 더 두드러지도록 조정.

| 속성 | v6.1 (스캔/비스캔) | v6.2 (스캔/비스캔) |
|------|:----------------:|:----------------:|
| fill alpha | 0.12 / 0.04 | 0.08 / 0.02 |
| stroke alpha | 0.60 / 0.20 | 0.35 / 0.10 |
| stroke width | 0.8 / 0.5 | 0.6 / 0.3 |
| 꼭짓점 blur | 2 | 4 |
| 꼭짓점 radius | 1.5 | 2.0 |
| 꼭짓점 alpha | 0.15 고정 | 0.1~0.2 깊이비례 |

### 2.4 파티클 밀도 증가 (~3500 → ~5000+)
**파일**: `holo_body.dart` — `_generateAnatomicalPoints`

모든 신체 부위의 파티클 수를 약 1.6~1.8배 증가. surfaceBias 전체 0.95 통일.

| 부위 | v6.1 | v6.2 | 증가율 |
|------|-----:|-----:|------:|
| 머리 | 365 | 580 | 1.59x |
| 목 | 50 | 80 | 1.60x |
| 몸통 | 790 | 1300 | 1.65x |
| 팔 (x2) | 330 | 540 | 1.64x |
| 다리 (x2) | 440 | 730 | 1.66x |
| **총 body** | **~2500** | **~4500+** | **~1.8x** |
| 장기+관절 | ~340 | ~340 | 유지 |
| **총합** | **~2840** | **~4840+** | **~1.7x** |

### 2.5 DraggableScrollableSheet 크래시 해결
**파일**: `monitoring_dashboard_screen.dart`

| 변경 전 | 변경 후 |
|---------|---------|
| `ListView(children: [...devices.map(...)])` | `CustomScrollView(slivers: [...])` |
| 모든 children 동시 빌드/교체 | `SliverList.builder` 가시영역만 빌드 |
| 폴링 리빌드 시 child.hasSize 실패 | `ValueKey('dev_${device.id}')` 안정적 재사용 |

구조:
- `SliverPadding` + `SliverChildListDelegate`: 드래그 핸들, 알림배너, 제목 (정적 헤더)
- `SliverPadding` + `SliverList.builder`: 기기 카드 목록 (동적, 가시영역만)
- `SliverPadding` + `SliverToBoxAdapter`: 빈 목록 안내 문구

---

## 3. 해결된 크래시 (누적 2건)

| # | 에러 | 원인 | 수정 | 버전 |
|---|------|------|------|------|
| 1 | `hitTest` Null check viewport | DraggableSheet 폴링 리빌드 타이밍 | `_SafeHitTestWrapper` (RenderProxyBox try-catch) | v6.1 |
| 2 | `child.hasSize` assertion | ListView children spread 전체 교체 | `CustomScrollView` + `SliverList.builder` | v6.2 |

---

## 4. 빌드 검증

```
flutter analyze  → 에러 0건 (737 info/warning)
flutter build web → 성공 (68.1초)
```

---

## 5. Before vs After (v6.1 → v6.2)

| 항목 | v6.1 | v6.2 |
|------|------|------|
| 파티클 수 | ~2500 body | ~4500+ body |
| surfaceBias | 0.88~0.92 혼재 | 0.95 통일 |
| 파티클 색상 | 단색 cyan | HSL 깊이 그라디언트 |
| 실루엣 외곽 | 직선 연결 | Catmull-Rom 스플라인 |
| 메시 가시성 | stroke 0.60 (눈에 띔) | stroke 0.35 (미세) |
| 꼭짓점 글로우 | blur 2, 고정 alpha | blur 4, 깊이비례 alpha |
| 기기목록 | ListView children | CustomScrollView+SliverList.builder |
| 크래시 | child.hasSize 발생 | 해결 |

---

## 6. 변경 파일 목록

| 파일 | 변경 유형 | 주요 변경 |
|------|----------|----------|
| `shared/widgets/holo_body.dart` | 수정 | 4종 품질 개선 |
| `data_hub/presentation/monitoring_dashboard_screen.dart` | 수정 | CustomScrollView 전환 |
