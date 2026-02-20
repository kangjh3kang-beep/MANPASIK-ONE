# 스프린트 14 — HoloBody v4.0 + 리더기 인터랙션 고도화 구현 보고서

> **문서 ID**: MPK-S14-HOLO-v4.0
> **작성일**: 2026-02-19
> **작성자**: Claude Opus 4.6 (AI 에이전트)
> **상태**: 완료 — flutter analyze 0 에러 / flutter build web 성공(72.8초)

---

## 목차

1. [Executive Summary](#1-executive-summary)
2. [변경 파일 총괄](#2-변경-파일-총괄)
3. [Part A: parentDataDirty 수정](#3-part-a-parentdatadirty-수정)
4. [Part B: HoloBody v4.0 구현](#4-part-b-holobody-v40-구현)
5. [Part C: 리더기 인터랙션 고도화](#5-part-c-리더기-인터랙션-고도화)
6. [Part D: Provider 추가](#6-part-d-provider-추가)
7. [Part E: 대시보드 UI 변경](#7-part-e-대시보드-ui-변경)
8. [발생 에러 및 해결 과정](#8-발생-에러-및-해결-과정)
9. [검증 결과](#9-검증-결과)
10. [기술 상세: 렌더링 파이프라인](#10-기술-상세-렌더링-파이프라인)
11. [기술 상세: 카트리지별 위험도 판정](#11-기술-상세-카트리지별-위험도-판정)
12. [성능 고려사항](#12-성능-고려사항)
13. [다음 단계](#13-다음-단계)

---

## 1. Executive Summary

### 정량 요약

| 지표 | 수치 |
|------|------|
| **변경 파일** | 4개 |
| **총 코드 줄** | 3,576줄 (4파일 합계) |
| **신규 렌더링 레이어** | +5개 (6→10단계 파이프라인) |
| **HoloBody 포인트 수** | ~3,500개 (성별 분화) |
| **삼각 메시 상한** | 2,000개 삼각형 |
| **카트리지 상세 타입** | 3종 (가스/환경/바이오) + 1 제네릭 |
| **flutter analyze** | 에러 0건 (정보/경고 736건) |
| **flutter build web** | 성공 (72.8초) |

### 작업 순서 (총 9단계)

| 단계 | 작업 | 상태 |
|------|------|------|
| 1 | parentDataDirty 수정 (Part A) | 완료 |
| 2 | Provider 추가 (Part D) | 완료 |
| 3 | HoloBody v4.0 Phase 1: 성별 체형 | 완료 |
| 4 | HoloBody v4.0 Phase 2: 신규 렌더링 5개 | 완료 |
| 5 | HoloBody v4.0 Phase 3: 격자 개선 | 완료 |
| 6 | 대시보드 UI: 성별 토글 + 파라미터 확장 (Part E) | 완료 |
| 7 | 호버 툴팁 (Part C-1) | 완료 |
| 8 | 카트리지별 상세 페이지 (Part C-2) | 완료 |
| 9 | 검증: flutter analyze + flutter build web | 완료 |

---

## 2. 변경 파일 총괄

| # | 파일 경로 | 변경 전 | 변경 후 | 변경 내용 |
|---|----------|--------|--------|----------|
| 1 | `shared/widgets/holo_body.dart` | ~630줄 | **1,081줄** | v4.0 완전 재작성: 성별 체형, 10단계 렌더링, 삼각 메시, ECG |
| 2 | `features/data_hub/presentation/monitoring_dashboard_screen.dart` | ~984줄 | **1,222줄** | parentDataDirty 수정, 성별 토글, 호버 툴팁 |
| 3 | `features/data_hub/presentation/providers/monitoring_providers.dart` | ~116줄 | **140줄** | holoGenderProvider, selectedBioDataProvider |
| 4 | `features/data_hub/presentation/widgets/device_detail_bottom_sheet.dart` | ~479줄 | **1,133줄** | 카트리지 타입별 상세 레이아웃 완전 재작성 |

---

## 3. Part A: parentDataDirty 수정

### 근본 원인

```
Failed assertion: line 4NN: '!semantics.parentDataDirty'
```

1. **AnimatedSwitcher** (line 257): 내부적으로 `Stack + FadeTransition`을 생성하여 old/new 위젯을 동시 배치. 외부 Stack의 `Positioned.fill` 안에서 AnimatedSwitcher가 자식의 parentData를 변경하면 semantics pass 시 assertion 발생.
2. **조건부 Positioned.fill** (line 302-314): `if (devices.isEmpty) Positioned.fill(...)` 이 Stack 자식 목록을 동적으로 변경하여 레이아웃 패스 불일치 유발.

### 수정 방안 (3단계)

#### A-1. AnimatedSwitcher → 항상-존재 + AnimatedOpacity

```dart
// 수정 전 (parentDataDirty 유발)
AnimatedSwitcher(
  duration: const Duration(milliseconds: 500),
  child: isBody ? HoloBody(...) : HoloGlobe(...),
)

// 수정 후 (Stack 자식 목록 고정)
Positioned.fill(
  child: AnimatedOpacity(
    opacity: isBody ? 0.0 : 1.0,
    duration: const Duration(milliseconds: 500),
    child: IgnorePointer(
      ignoring: isBody,
      child: HoloGlobe(...),
    ),
  ),
),
Positioned.fill(
  child: AnimatedOpacity(
    opacity: isBody ? 1.0 : 0.0,
    duration: const Duration(milliseconds: 500),
    child: IgnorePointer(
      ignoring: !isBody,
      child: HoloBody(...),
    ),
  ),
),
```

**핵심**: AnimatedSwitcher 제거 → Stack 자식 목록 완전 고정 → parentData 절대 변경 안 됨.

#### A-2. 조건부 빈 상태 → 항상 존재 + Opacity

```dart
// 수정 전
if (devices.isEmpty) Positioned.fill(child: ...)

// 수정 후
Positioned.fill(
  child: IgnorePointer(
    ignoring: devices.isNotEmpty,
    child: AnimatedOpacity(
      opacity: devices.isEmpty ? 1.0 : 0.0,
      duration: const Duration(milliseconds: 300),
      child: Center(child: Text('연결된 기기가 없습니다.')),
    ),
  ),
),
```

#### A-3. ExcludeSemantics 유지

기존 `_buildVisualization`의 내부 Stack 전체를 `ExcludeSemantics`로 감싸는 방식 유지.

---

## 4. Part B: HoloBody v4.0 구현

### Phase 1: 성별 체형 분화 시스템

#### HoloGender enum

```dart
enum HoloGender { male, female }
```

`holo_body.dart` 최상단에 선언 (circular import 방지를 위한 단일 출처).

#### 위젯 API 확장

```dart
class HoloBody extends StatefulWidget {
  final double width, height;
  final Color color;
  final Color? accentColor;
  final HoloGender gender;                  // 신규
  final Map<String, dynamic> bioData;       // 신규
  final bool showDataLabels;                // 신규
  final bool showEcg;                       // 신규
}
```

#### 성별 비율표

| 부위 | 남성 | 여성 |
|------|------|------|
| 어깨 삼각근 X | ±w×0.45 | ±w×0.38 |
| 흉곽 상단 폭 | w×0.48 | w×0.42 |
| 흉곽 하단 폭 | w×0.38 | w×0.36 |
| 복부 폭 | w×0.38→0.35 | w×0.35→0.28 (잘록) |
| 골반 상단 폭 | w×0.35 | w×0.42 (넓음) |
| 골반 하단 폭 | w×0.38 | w×0.46 (넓음) |
| 상완 시작 X | ±w×0.48 | ±w×0.42 |
| 대퇴 시작 X | ±w×0.18 | ±w×0.20 |
| 여성 전용: 가슴 | — | 각 35pt (±w×0.13, -h×0.23) |

#### didUpdateWidget 성별 변경 감지

```dart
@override
void didUpdateWidget(covariant HoloBody old) {
  super.didUpdateWidget(old);
  if (old.gender != widget.gender || old.width != widget.width || old.height != widget.height) {
    _points.clear();
    _triangles.clear();
    _generateAnatomicalPoints();
    _precomputeTriangleMesh();
  }
}
```

### Phase 2: 신규 렌더링 레이어 5개

기존 6단계 → **10단계 파이프라인**:

```
 1. _drawWireframeGrid()       (기존, 50줄로 증가)
 2. _drawPlatformRings()       ★ 신규 — 발 아래 3중 동심 타원
 3. _drawSkeletonWires()       (기존, 성별 비율 적용)
 4. _drawTriangleMesh()        ★ 신규 — 공간 해시 사전계산, 상한 2000삼각형
 5. _drawParticles()           (기존, 호흡 애니메이션)
 6. _drawBodyContourGlow()     ★ 신규 — 24 Y슬라이스 외곽 Path
 7. _drawScanningHoops()       ★ 신규 — 5개 후프 Y진동
 8. _drawScanLaser()           (기존)
 9. _drawHeartbeat()           (기존)
    _drawEcgWaveform()         ★ 신규 — PQRST 세그먼트 (showEcg일 때)
10. _drawEnergyWaves()         (기존)
    _drawOrganDataLabels()     ★ 신규 — 장기 데이터 라벨 (showDataLabels일 때)
```

#### 2-A. 플랫폼 링 (_drawPlatformRings)

- 발 아래(y=h×0.52)에 3중 동심 타원
- MaskFilter blur 글로우
- breathValue 기반 미세 펄스

#### 2-B. 삼각 메시 (_drawTriangleMesh)

- `_Triangle(i0, i1, i2)` 클래스
- `_precomputeTriangleMesh()`: 공간 해시(bucket=30px) 사전계산, 상한 2000삼각형
- initState/gender 변경 시 1회 계산
- 스캔 레이저 근처 alpha 강조

#### 2-C. 바디 컨투어 글로우 (_drawBodyContourGlow)

- body 포인트 24 Y슬라이스
- min/max X로 외곽 Path
- MaskFilter blur 8
- accentColor alpha 0.12 채움

#### 2-D. 스캐닝 후프 (_drawScanningHoops)

- `_hoopController` (6초, reverse) 추가
- 5개 후프 Y진동, 다른 위상 오프셋
- 타원 stroke + blur 4 글로우

#### 2-E. ECG 파형 (_drawEcgWaveform)

- 심장 우측에 PQRST 세그먼트 렌더링
- pulseValue 동기화 (pulse > 0.3일 때 R파 활성)
- clipRect로 영역 제한

#### 2-F. 장기 데이터 라벨 (_drawOrganDataLabels)

- bioData 맵 키 매칭:
  - `Pulse` / `pulse` / `HR` → 심장 (좌측)
  - `O2` / `SpO2` / `spo2` → 폐 (우측)
  - `Stress` / `stress` → 뇌 (좌측)
  - `Glucose` / `glucose` → 간 (우측)
- L자 연결선 + 반투명 RRect 배경 + 값 텍스트

### Phase 3: 와이어프레임 격자 개선

- 수평 라인 수: 30 → 50
- 골격 좌표 참조 정밀 폭 변조
- 양끝 alpha 페이드아웃
- 경계선 외 ±5% 마진

---

## 5. Part C: 리더기 인터랙션 고도화

### C-1. 호버 툴팁

#### 구현 구조

`_InteractiveNodesLayer`를 **StatefulWidget**으로 변환:

```dart
class _InteractiveNodesLayer extends StatefulWidget { ... }

class _InteractiveNodesLayerState extends State<_InteractiveNodesLayer> {
  int? _hoveredIndex;

  int? _findNodeAt(Offset pos) {
    for (int i = 0; i < positions.length; i++) {
      if ((pos - positions[i].nodePos).distance < 28) return i;
    }
    return null;
  }

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      onExit: (_) { if (_hoveredIndex != null) setState(() => _hoveredIndex = null); },
      child: Listener(
        onPointerHover: (event) {
          final found = _findNodeAt(event.localPosition);
          if (found != _hoveredIndex) setState(() => _hoveredIndex = found);
        },
        onPointerUp: (event) {
          final found = _findNodeAt(event.localPosition);
          if (found != null) widget.onNodeTap(widget.devices[found]);
        },
        child: Stack(children: [
          CustomPaint(painter: _NodesPainter(..., hoveredIndex: _hoveredIndex)),
          if (_hoveredIndex != null) _buildHoverTooltip(_hoveredIndex!),
        ]),
      ),
    );
  }
}
```

#### 호버 툴팁 디자인 (180×72px)

```
┌─────────────────────────────────────┐
│ ● LIVE  거실 공기질 측정기              │
│ CO2: 450 ppm · VOC: 0.05            │
│ 🔋 72%                               │
└─────────────────────────────────────┘
```

- BackdropFilter(blur 12) + 반투명 배경
- sanggamGold 테두리 (0.5px)
- 노드 상단(nodePos.dy - 80)에 배치, 화면 클램프
- 상태 배지 + 이름 + currentValues 상위 2개 + 배터리

#### _NodesPainter 호버 하이라이트

- hoveredIndex 노드: 반지름 +2, waveCyan alpha 0.25 글로우, 테두리 1.5→2.0
- selected가 우선 (선택된 노드는 기존 하이라이트 유지)

### C-2. 카트리지 타입별 상세 페이지

#### 공통 레이아웃

```
┌─────────────────────────────────────┐
│ ── 드래그 핸들 ──                      │
│ [아이콘]  기기명                       │
│ ● LIVE   가스 카트리지                │
├─────────────────────────────────────┤
│     [타입별 전용 콘텐츠 영역]           │
├─────────────────────────────────────┤
│ 배터리 ████████░░ 72%                │
│ 신호   ██████████ 95%                │
├─────────────────────────────────────┤
│ 측정 추이 (Sparkline 차트)            │
├─────────────────────────────────────┤
│ 기기 정보 (ID, 상태)                   │
│ [기기 관리] 버튼                       │
└─────────────────────────────────────┘
```

#### 가스 카트리지 (_GasDetailSection)

| 센서 | 안전 기준 | 주의 | 위험 |
|------|----------|------|------|
| CO | < 30ppm | 30~70ppm | > 70ppm |
| LNG | < 0.5% | 0.5~2% | > 2% |
| Smoke | 'None' | — | != 'None' |
| CO2 | < 1000ppm | 1000~2000ppm | > 2000ppm |
| VOC | < 0.5mg/m³ | 0.5~1.0 | > 1.0 |

- 위험도별 색상: 안전(`#00E676`), 주의(`#FFC107`), 위험(`#FF4D4D`)
- 1개라도 '주의' 이상 → "환기 권장" 표시

#### 환경 카트리지 (_EnvDetailSection)

- **쾌적도 점수 알고리즘**:
  - Temp 최적 22°C, 편차 1°C당 -5점 (100점 만점)
  - Humidity 최적 50%, 편차 1%당 -2점
  - Light 최적 300lux, 편차 1%당 -1점
  - 종합 = (temp + humi + light) / 3, clamp(0, 100)
- 쾌적도 인디케이터: 그라디언트 바 + 위치 표시 점
- 파라미터 범위 가이드: Temp(18~26°C), Humidity(40~60%), Light(100~500lux)

#### 바이오 카트리지 (_BioDetailSection)

| 지표 | 정상 | 주의 | 위험 |
|------|------|------|------|
| Pulse | 60~100bpm | <60 서맥 | >100 빈맥 |
| O2 | >95% | 90~95% | <90% |
| Stress | Low → 양호 | Medium → 보통 | High → 주의 |

- 미니 ECG 파형: `_MiniEcgPainter` CustomPainter (PQRST 세그먼트)
- 상태 라벨: 정상/주의/위험에 따른 색상 코드

---

## 6. Part D: Provider 추가

`monitoring_providers.dart`에 2개 Provider 추가:

```dart
// 1. 성별 토글 (바이오 탭 전용)
final holoGenderProvider = StateProvider<HoloGender>((ref) => HoloGender.male);

// 2. 선택된 바이오 기기의 bioData 파생
final selectedBioDataProvider = Provider<Map<String, dynamic>>((ref) {
  final device = ref.watch(selectedDeviceProvider);
  if (device != null && device.type == DeviceType.bioCartridge) {
    return device.currentValues;
  }
  final devicesAsync = ref.watch(pollingConnectedDevicesProvider);
  return devicesAsync.when(
    data: (devices) {
      final bio = devices.where((d) =>
        d.type == DeviceType.bioCartridge &&
        d.status == DeviceConnectionStatus.connected).toList();
      return bio.isNotEmpty ? bio.first.currentValues : <String, dynamic>{};
    },
    loading: () => <String, dynamic>{},
    error: (_, __) => <String, dynamic>{},
  );
});
```

`HoloGender` enum은 `holo_body.dart`에서 `show` import:
```dart
import 'package:manpasik/shared/widgets/holo_body.dart' show HoloGender;
```

---

## 7. Part E: 대시보드 UI 변경

### 성별 토글 (AppBar actions)

- filterTab == 3 (바이오) 일 때만 AppBar에 토글 표시
- 스타일: sanggamGold 테두리 캡슐, 아이콘(♂/♀) + 라벨(남/여)
- onTap: `holoGenderProvider` 토글

### HoloBody 호출 파라미터 확장

```dart
HoloBody(
  key: ValueKey('body_${gender.name}'),
  width: globeSize,
  height: bodyH,
  color: isDark ? AppTheme.waveCyan : const Color(0xFF00ACC1),
  accentColor: isDark ? AppTheme.sanggamGold : const Color(0xFFFF4D4D),
  gender: ref.watch(holoGenderProvider),
  bioData: ref.watch(selectedBioDataProvider),
  showDataLabels: true,
  showEcg: true,
)
```

---

## 8. 발생 에러 및 해결 과정

### 에러 1~3: HoloGender ambiguous import

**증상**: `HoloGender` enum이 `monitoring_providers.dart`와 `holo_body.dart` 양쪽에 정의되어 대시보드에서 import 시 모호성 발생.

**해결**: `monitoring_providers.dart`에서 enum 제거, `holo_body.dart`의 것을 `show HoloGender`로 import.

### 에러 4: StateProvider 타입 불일치

**증상**: ambiguous HoloGender로 인한 타입 해석 실패.

**해결**: 위 에러 1~3 해결로 자동 해소.

### 에러 5: Listener.onPointerExit 미정의

**증상**: `Listener` 위젯에 `onPointerExit` 속성이 존재하지 않음.

**해결**: `Listener`를 `MouseRegion`으로 감싸서 `onExit` 처리. `Listener`는 `onPointerHover`와 `onPointerUp`만 담당.

```dart
return MouseRegion(
  onExit: (_) {
    if (_hoveredIndex != null) setState(() => _hoveredIndex = null);
  },
  child: Listener(
    onPointerHover: (event) { ... },
    onPointerUp: (event) { ... },
    child: Stack(...),
  ),
);
```

---

## 9. 검증 결과

| 검증 항목 | 결과 |
|----------|------|
| **flutter analyze** | 에러 0건 (정보/경고 736건) |
| **flutter build web** | 성공 — `build/web` 생성 (72.8초) |
| **Wasm 경고** | flutter_secure_storage_web 관련 (기능 무관, 정보성) |
| **폰트 트리셰이킹** | CupertinoIcons 99.4%, MaterialIcons 97.2% 감소 |

---

## 10. 기술 상세: 렌더링 파이프라인

```
┌─────────────────────────────────────────────────────────────┐
│                    HoloBody v4.0 렌더링 파이프라인              │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Layer 1: 와이어프레임 격자 (50줄, alpha 페이드)                │
│    ↓                                                        │
│  Layer 2: 플랫폼 링 (3중 동심 타원, blur 글로우)               │
│    ↓                                                        │
│  Layer 3: 골격 와이어 (성별 비율 적용)                         │
│    ↓                                                        │
│  Layer 4: 삼각 메시 (공간 해시, ≤2000삼각형)                   │
│    ↓                                                        │
│  Layer 5: 파티클 호흡 (~3500개 포인트)                         │
│    ↓                                                        │
│  Layer 6: 바디 컨투어 글로우 (24슬라이스 외곽)                  │
│    ↓                                                        │
│  Layer 7: 스캐닝 후프 (5개, 6초 주기)                         │
│    ↓                                                        │
│  Layer 8: CT/MRI 스캔 레이저                                  │
│    ↓                                                        │
│  Layer 9: 심박 펄스 + ECG 파형 (PQRST)                        │
│    ↓                                                        │
│  Layer 10: 에너지 파동 + 장기 데이터 라벨                      │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 애니메이션 컨트롤러

| 컨트롤러 | 주기 | 용도 |
|----------|------|------|
| `_controller` | 8초 반복 | 기본 호흡/파티클 |
| `_scanController` | 4초 반복 | 스캔 레이저 |
| `_pulseController` | 1.2초 반복 | 심박 펄스 |
| `_hoopController` | 6초 왕복 | 스캐닝 후프 |

### 삼각 메시 알고리즘

```
1. 공간 해시 생성 (bucket = 30px)
2. 각 포인트 p에 대해:
   a. 같은/인접 버킷 내 포인트 수집
   b. 거리 < maxDist (w*0.12) 인 이웃 필터링
   c. 이웃 쌍 (n1, n2) 중:
      - n1~n2 거리 < maxDist
      - 삼각형 면적 > 최소 면적 (50.0)
      - 중복 검사 (정렬 키)
   d. _triangles에 추가
3. 상한 2000개 도달 시 중단
```

---

## 11. 기술 상세: 카트리지별 위험도 판정

### 가스 카트리지 — 위험도 분류 로직

```dart
enum _DangerLevel { safe, caution, danger }

_DangerLevel _assessGasDanger(String key, dynamic value) {
  final v = (value is num) ? value.toDouble() : 0.0;
  switch (key.toLowerCase()) {
    case 'co':     return v > 70 ? danger : v > 30 ? caution : safe;
    case 'lng':    return v > 2 ? danger : v > 0.5 ? caution : safe;
    case 'smoke':  return (value.toString() != 'None') ? danger : safe;
    case 'co2':    return v > 2000 ? danger : v > 1000 ? caution : safe;
    case 'voc':    return v > 1.0 ? danger : v > 0.5 ? caution : safe;
    default:       return safe;
  }
}
```

### 환경 카트리지 — 쾌적도 점수 알고리즘

```dart
double _computeComfortScore(Map<String, dynamic> values) {
  double tempScore = 100, humiScore = 100, lightScore = 100;

  if (values.containsKey('Temp')) {
    tempScore = max(0, 100 - (values['Temp'] - 22.0).abs() * 5);
  }
  if (values.containsKey('Humidity')) {
    humiScore = max(0, 100 - (values['Humidity'] - 50.0).abs() * 2);
  }
  if (values.containsKey('Light')) {
    lightScore = max(0, 100 - (values['Light'] - 300.0).abs() * 0.1);
  }

  return ((tempScore + humiScore + lightScore) / 3).clamp(0, 100);
}
```

### 바이오 카트리지 — 생체 판정

| 지표 | 키 매칭 | 정상 범위 | 판정 로직 |
|------|---------|----------|----------|
| Pulse | `Pulse`, `HR` | 60~100 bpm | <60→서맥, >100→빈맥 |
| SpO2 | `O2`, `SpO2` | >95% | 90~95→주의, <90→위험 |
| Stress | `Stress` | Low | Medium→보통, High→주의 |

---

## 12. 성능 고려사항

| 항목 | 전략 |
|------|------|
| **삼각 메시** | initState/gender변경 시 1회 사전계산, 렌더링은 인덱스 참조만 |
| **TextPainter** | 장기 라벨에만 사용, bioData 변경 드물어 매 프레임 허용 |
| **MaskFilter** | 플랫폼 링 + 컨투어 + 후프에만, 삼각 메시엔 미사용 |
| **RepaintBoundary** | 대시보드에서 HoloBody 감싸고 있음 |
| **호버 툴팁** | setState 최소화 — hoveredIndex 변경 시만 |
| **shouldRepaint** | 구체적 값 비교 (hoveredIndex, selectedId, 애니메이션 값) |
| **IgnorePointer** | 비활성 홀로그램(Globe/Body) 이벤트 차단 |

---

## 13. 다음 단계

| 우선순위 | 작업 | 관련 |
|----------|------|------|
| P1 | 실기기 연동 테스트 (실제 카트리지 데이터로 검증) | 스프린트 15 |
| P2 | 대시보드 성능 프로파일링 (60fps 유지 확인) | 성능 |
| P3 | 접근성 (Semantics) 추가 — ExcludeSemantics 내부 항목 | UX |
| P4 | 다크/라이트 모드 전환 시 색상 보정 | 디자인 |
| P5 | 가로 모드 대응 (현재 Portrait 전용) | UX |

---

> **문서 끝** — MPK-S14-HOLO-v4.0
> 다음 검증: 스프린트 15 전구간 통합 테스트
