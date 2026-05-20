# ManPaSik (萬波息) — 프로젝트 CLAUDE.md (v2.1 베이스라인)

## 프로젝트 정의
차동측정 기반 범용분석 POCT 시스템의 소프트웨어 스택.
하드웨어 리더기(STM32F405) + 일회용 카트리지를 연동하는 모바일 앱(Flutter) + Rust 핵심 엔진 + Go 백엔드 + AI/ML 파이프라인으로 구성된 6-Layer 통합 생태계.

## SSOT 베이스라인 (절대 변경 금지, 코드에 반드시 반영)
- 환율: ₩1,480/USD
- 커넥터: CSI v1.0 = Samtec MECF-08-01-L-DV, 16핀, 1.27mm (E12-IF 절대 사용 금지)
- 차분식: Sdiff_n = S_n - α_n * R_n (alpha 기본값 0.98, 동적 보정 범위 0.90~1.10)
- PPM 크기: 49.70×30×4.30mm
- Universal AFE: 9블록 (Stage별 자동 활성화)
- 정확도: 3축 결합 기준 92~98%
- 핑거프린트: 88차원(기본) → 448차원 → 896차원(통합) → 1792차원(풀스펙)

## 통합 아키텍처 기술 스택 (확정)
- Layer 6 (인프라): AWS/GCP Multi-Cloud, Kubernetes
- Layer 5 (백엔드): Go 1.22+ (gRPC MSA), PostgreSQL 16 + TimescaleDB + Milvus
- Layer 4 (앱): Flutter 3.x (Riverpod 2.x + freezed) + flutter_rust_bridge 2.x
- Layer 3 (코어): Rust 핵심 엔진 (no_std 호환, Harness Abstraction Layer)
- Layer 2 (하드웨어 제어): embedded-hal, RAFE 스위치
- Layer 1 (하드웨어): STM32F405, nRF52832(BLE), PN7150(NFC)

## 하네스 엔지니어링 (Harness Engineering) 6대 원칙
1. H1 모듈 독립성: AFE 블록은 HAL(SensorTrait)로 분리, 구체 타입 직접 참조 금지.
2. H2 후방 호환: CartridgeManifest v1.0 카트리지는 v2.0 리더기/앱에서 반드시 동작해야 함.
3. H3 사전인증 플랫폼: 510(k) predicate 기반 모듈 단위 간소화 인증 고려 설계.
4. H4 점진적 확장: Stage-1(전기화학)부터 Stage-3(NAAT)까지 순차 확장 고려 구현.
5. H5 인터페이스 계약: BLE GATT 서비스 UUID 및 NFC 매니페스트 구조는 버전화하여 계약 명시.
6. H6 실패 격리: 센서 1개 오류시 전체 중단 금지. `Result<T,E>`로 에러 격리 및 폴백 경로 제공.

## 코딩/설계 규칙
1. Rust: Result<T, E> 패턴 필수. unwrap() 절대 금지 (테스트 제외). clippy 경고 0건 유지.
2. Flutter: 오프라인 우선(CRDT). 상태관리는 Riverpod+freezed 불변 상태.
3. 백엔드: Go gRPC, FHIR R4 호환 직렬화, PCCP 대응 MLOps (모델 레지스트리) 파이프라인 대비.
4. AI/ML: TFLite 엣지 양자화(INT8) 우선, XAI(SHAP) 적용 및 클라우드-연합학습(FL) 하이브리드.
5. 테스트/보안: 커버리지 80% 이상, TPM 2.0 / 해시체인 / AES-256-GCM / PII 토큰화.

## 절대 금지 사항
- E12-IF 12핀 참조 금지 (기존 문서들에서 발견 시 CSI v1.0 16핀으로 일괄 수정 대상).
- AI가 생성한 가상의 검증/테스트 데이터를 마치 실제인 것처럼 삽입 금지.
- 검증되지 않은 알고리즘, 파라미터를 코드 주석에 "검증됨"으로 표기 금지.

## AI 어시스턴트(Agent) 기본 동작 및 보고 지침 (Core Workflow)
1. **작업 기록 및 실시간 공유**: 모든 코드 변경사항 및 개발 결과물은 즉시 시스템 작업 일지나 Artifact로 기록하고 사용자에게 투명하게 진행 상황을 공유한다.
2. **반복적 검증 (코드리뷰 > 린트 > 빌드 > 테스트)**: 코드 작성 완료 즉시 다음 단계로 넘어가기 전에, 자체 코드리뷰를 거치고 각 언어별 린트(Rust clippy, Flutter analyzer, Go vet 등) 수행 및 단위 컴파일/테스트를 반드시 통과시켜 무결성을 확보한다.