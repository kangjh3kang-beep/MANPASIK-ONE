# 초고해상도 메디컬 대시보드 및 홀로그램 완벽 구현 계획 (Ultimate Medical Dashboard)

## 1. 문제 해결 (Crash Fix)
*   **오류**: `RenderBox was not laid out` (RenderSliverList 관련).
*   **원인**: `DraggableScrollableSheet` 내부의 `ListView`나 `MonitoringDashboardScreen`의 `Stack` 레이아웃에서 제약조건(Constraints) 전달 불량으로 추정.
*   **해결**:
    *   `DraggableScrollableSheet`의 빌더 내부 구조를 `Container` -> `SingleChildScrollView` 대신 명확한 `ListView` 구조로 단순화하고, 부모 크기 제약을 명확히 함.
    *   불필요한 `RepaintBoundary` 제거 또는 위치 조정.

## 2. 홀로그램 품질 "완벽" 개선 (Reference Matching)
사용자가 제공한 레퍼런스(푸른색 와이어프레임 + 그라데이션 + 데이터 카드)와 100% 일치하는 **Cyber Mesh Engine**을 구현합니다.

### 2.1. HoloBody V5: Cyber Mesh & Volumetric Glow
*   **Texture**: 가로 스캔 라인뿐만 아니라, **세로 격자(Grid)**를 추가하여 면(Mesh)의 느낌을 살림.
*   **Rendering**:
    *   **Fresnel Glow**: 인체 외각선이 가장 밝게 빛나고, 내부는 투명하게 처리.
    *   **Depth Fade**: 뒤에 있는 메쉬는 어둡게, 앞에 있는 메쉬는 밝게 처리하여 공간감 극대화.
    *   **Body Scan Effect**: 위아래로 스캔 바가 지나갈 때 해당 단면이 하이라이트되는 효과 강화.

### 2.2. Medical HUD Interaction
*   **Floating Cards (데이터 카드)**:
    *   기존의 단순 아이콘 노드를 **"실시간 의료 데이터 카드(Pulse, SpO2, Brain, Gluc)"**로 교체.
    *   레퍼런스 이미지와 동일한 **Cyberpunk/Medical UI 스타일** (투명한 청록색 패널, 네온 텍스트).
*   **Anchoring Smart Lines**:
    *   데이터 카드가 단순히 주변을 도는 게 아니라, **신체의 정확한 부위(심장, 머리, 손목 등)**에서 선이 뻗어 나와 카드와 연결됨.
    *   점선(Dotted Line) + 끝점의 원형 앵커(Circle Anchor).

## 3. 구현 단계
1.  **Crash Fix**: `MonitoringDashboardScreen` 재구조화 (안정성 최우선).
2.  **HoloBody V5**: Cyber Mesh 엔진 탑재.
3.  **Medical Widgets**: `MedicalStatCard` 위젯 제작 및 `OrbitalNode` 대체.
4.  **Integration**: 신체 부위 좌표 - 위젯 연결 로직 구현.

---
**작성일**: 2026-02-19
**목표**: "조악함" 탈피, 레퍼런스급 "High-End Medical Dashboard" 구현.
