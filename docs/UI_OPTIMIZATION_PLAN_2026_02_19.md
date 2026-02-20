# UI 최적화 및 홀로그램 고도화 계획 (2026-02-19)

## 1. 현황 분석 (Current Status & Issues)

### 1.1. 화면 가독성 및 겹침 문제 (Layout Overlaps)
*   **현상**: 상단 요약 바(Summary Bar)와 경고 배너(Alert Banner)가 화면 중앙의 홀로그램(Globe/Body) 및 리더기 노드와 겹쳐서 정보를 가리고 있음.
*   **원인**: 현재 `Stack` + `Positioned` 기반의 절대 좌표 방식을 사용하고 있어, 화면 높이가 낮거나 UI 요소(경고창 등)가 추가될 때 중앙 콘텐츠의 영역을 침범함.
*   **영향**: 핵심 데이터인 리더기 연결 상태와 홀로그램 시각화를 확인하기 어려움.

### 1.2. 인체 홀로그램 디테일 부족 (Low Fidelity HoloBody)
*   **현상**: 현재 `HoloBody`가 단순한 원기둥/구의 조합으로 구현되어 있어, 사람의 신체 곡선이나 관절 등의 디테일이 부족함. "파란 점들의 뭉치" 정도로만 보임.
*   **원인**: 파티클 생성 알고리즘이 단순 기하학 도형 기반임.
*   **목표**: "인간의 신체와 같은 모양"으로 고도화 (어깨, 허리 라인, 관절 등 해부학적 특징 반영).

### 1.3. 하단 네비게이션/패널 오류
*   **현상**: 하단 `DraggableScrollableSheet`가 제스처를 방해하거나, 위치가 애매하여 조작성이 떨어짐.

---

## 2. 상세 구현 계획 (Detailed Implementation Plan)

### 2.1. 레이아웃 엔진 재설계 (Responsive Layout Engine)
기존의 절대 좌표 방식을 버리고, **"영역 확보형(Safe Area)"** 레이아웃으로 전환합니다.

*   **Adaptive Column Layout**:
    1.  **Header Area**: AppBar + TabBar (높이 고정)
    2.  **Status Area**: Summary Bar + Alert Banner (내용에 따라 가변 높이)
    3.  **Visualization Area (Expanded)**: 남은 공간을 모두 차지하며, 이 영역의 **정중앙(Center)**에 홀로그램을 배치.
    4.  **Bottom Panel Area**: 하단 패널이 올라올 최소 공간 확보.
*   **Dynamic Scaling**: `Visualization Area`의 높이를 계산하여, 홀로그램의 크기(`Size`)와 궤도 반지름(`Radius`)을 동적으로 축소/확대. (예: UI에 가려지지 않도록 자동으로 작아짐)

### 2.2. HoloBody V2: 해부학적 파티클 시스템 (Anatomical Particle System)
더욱 정교한 인체 모델링을 위해 파티클 생성 로직을 전면 개편합니다.

*   **Body Part Segmentation**:
    *   **Torso**: 흉곽(넓음) -> 허리(잘록함) -> 골반(넓음)으로 이어지는 곡선 적용.
    *   **Limbs**: 허벅지/종아리, 팔뚝/전완근의 두께 차이 및 관절 부위(무릎, 팔꿈치) 파티클 밀도 강화.
    *   **Head**: 턱선과 두상 라인 디테일 추가.
*   **Volume Rendering**: 내부를 꽉 채우는 방식 대신, 피부층(Surface)과 신경계(Nervous System)를 구분하여 깊이감 부여.

### 2.3. 하단 패널 및 네비게이션 최적화
*   **Panel Interaction**: 드래그 패널의 초기 높이(`initialChildSize`)를 조정하여 중앙 홀로그램을 가리지 않도록 설정.
*   **Layering**: 패널이 확장될 때 반투명 배경(Blur)을 강화하여 뒤쪽 콘텐츠와 시각적으로 분리.

---

## 3. 실행 단계 (Execution Steps)

1.  **Step 1: Layout Refactoring**
    *   `Stack` 기반 코드를 `Column` + `Expanded` 구조로 변경.
    *   `LayoutBuilder`를 이용해 가용 높이 계산 로직 정밀화.
2.  **Step 2: HoloBody Upgrade**
    *   `holo_human.dart` 알고리즘 개선 (곡선 함수 적용).
3.  **Step 3: Panel Tuning**
    *   `DraggableScrollableSheet` 파라미터 미세 조정.
4.  **Step 4: Verification**
    *   다양한 해상도 및 경고 배너 유무 시나리오 테스트.

---
**목표 완료 시간**: 금일 내 완료 예정.
