# UI 가독성 개선 및 인체 홀로그램 실루엣 고도화 계획

## 1. 현황 및 문제점 분석
*   **시스템 안정성**: `semantics.parentDataDirty` 크래시가 지속적으로 발생하여 앱 사용이 불가능한 수준.
    *   *원인*: 복잡한 Stack/Positioned 애니메이션이 접근성 트리(Semantics Tree) 갱신과 충돌.
*   **UI 가독성 저하**:
    *   배경(파티클/링)과 전경(텍스트/데이터)이 겹쳐서 글자가 잘 안 보임.
    *   화면 구성이 산만하여 정보의 위계가 불명확함.
*   **HoloBody 품질**:
    *   현재의 파티클 방식은 "점들의 집합"으로만 보여, 사용자가 "사람의 형체"로 인식하기에 직관성이 떨어짐("조악하다"는 평가).

## 2. 개선 목표
1.  **Crash Free**: 시각화 영역의 Semantics를 제외(`ExcludeSemantics`)하여 렌더링 충돌 원천 차단.
2.  **High-Fidelity Contour (실루엣 고도화)**:
    *   단순 점이 아닌, **"외각 라인(Outline)"**을 디테일하게 그려 인간의 신체 윤곽을 명확히 함.
    *   '사이버네틱'한 느낌의 와이어프레임(Wireframe) 스타일 추가.
3.  **Readable Interface (인터페이스 개편)**:
    *   **Zone Separation**: [시각화 영역]과 [정보 영역]을 명확히 분리.
    *   **High Contrast**: 텍스트 가독성을 위해 정보 카드에 확실한 배경(Backdrop Blur)과 테두리 적용.

## 3. 상세 구현 계획

### 3.1. HoloBody V3: Silhouette & Wireframe
*   **기존**: 무작위 파티클 생성 (Random Points).
*   **변경**:
    *   **Contour Paths**: 머리, 어깨, 팔, 다리의 외각선을 따라 흐르는 **곡선(Path)** 데이터 생성.
    *   **Connectors**: 주요 관절(어깨-팔꿈치-손목)을 잇는 뼈대 라인(Bone Lines) 추가.
    *   **Effect**: 외각선은 밝게(Glow), 내부는 은은한 파티클로 채워 "투명한 사이보그" 느낌 구현.

### 3.2. MonitoringDashboard 리팩토링 (가독성 중심)
*   **Layout**:
    *   기존: 전체 화면 Stack (요소가 둥둥 떠다님).
    *   변경:
        *   **Background**: 어두운 우주/클라우드 배경 (Darkened).
        *   **Center**: HoloBody/HoloGlobe (크기 최적화).
        *   **Foreground UI**: 상단/하단에 **고정된 패널(Glassmorphism Panel)** 배치. 텍스트가 절대 홀로그램 위에 겹치지 않게 처리.
*   **Semantics Fix**: `ExcludeSemantics`로 시각화 위젯 감싸기 (가장 시급).

### 3.3. 실행 순서
1.  **Urgent Fix**: `MonitoringDashboardScreen.dart`에 `ExcludeSemantics` 적용 및 레이아웃 정리.
2.  **HoloBody Upgrade**: `holo_body.dart`를 외각선(Outline) 드로잉 로직으로 전면 교체.
3.  **UI Polish**: 배경 밝기 조절 및 텍스트 명도 대비 강화.

---
**작성일**: 2026-02-19
**작성자**: Antigravity AI
