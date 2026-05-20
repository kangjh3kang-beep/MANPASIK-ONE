# ManPaSik (萬波息) 기술 사양서 v2.4.3

**문서번호**: MPS-TECH-SPEC-v2.4.3
**작성일**: 2026-04-19
**갱신일**: 2026-04-19 Minor 개선 6건 — B04 metrics 튜플 타입 명시, B05 SelfDiagnostics/SelfHealing 역할 경계 서술, B06 Weibull 파라미터 상세, B07 SessionCheckpoint 직렬화 포맷
**용도**: AI IDE (Cursor/Antigravity/Windsurf) 기반 시스템 구현의 단일 참조 문서
**선행문서**: MPS-SYS-MASTER-v2.0, MPS-IDE-PROMPT-v1.0
**읽는 법**: IDE AI가 본 문서를 컨텍스트로 참조하여 모듈을 구현한다. 각 모듈의 [인터페이스 계약]은 구현 시 반드시 준수해야 하는 계약이며, [테스트 기준]은 구현 완료의 판단 기준이다.

---

## 목차

```
1.  프로젝트 개요 및 SSOT 상수
2.  프로젝트 구조 (디렉토리 트리)
3.  아키텍처 (6-Layer + HAL)
4.  핵심 타입 & 인터페이스 정의
5.  Layer 1-2: 하드웨어 + 펌웨어
6.  Layer 3: Rust Core Engine (모듈별 상세 사양)
7.  Layer 4: Flutter 모바일 앱
8.  Layer 5: Go gRPC 백엔드
9.  Layer 6: AI/ML + 인프라
10. SDK & 카트리지 마켓플레이스
11. 보안 아키텍처
12. 테스트 전략 (전체)
13. 구현 로드맵 + Phase별 체크리스트
Appendix A: IDE CLAUDE.md (프로젝트 루트 배치용)
Appendix B: 참조 프롬프트 모음
```

**v2.1 추가사항**: §6.8 CheckupSession, §6.9 DiseaseRiskEngine, §6.10 AutoClassifier, §7.5 종합검진 UI 상태관리, §8.4 종합검진 gRPC 서비스

**v2.2 추가사항**: §6.11 ContextEngine, §6.12 NutritionAdvisor, §6.13 ShoppingBridge, §6.14 HabitTracker, §7.6 유기적 연동 UI 상태관리, §8.5 유기적 연동 gRPC 서비스

**v2.3 추가사항**: §6.15 SelfDiagnostics(HwHealthMonitor, FwWatchdogBridge, DataIntegrityChecker, PredictiveMaintenanceEngine, ErrorReporter, SystemHealthOrchestrator), §7.7 자가검증 UI 상태관리, §8.6 자가검증 백엔드 서비스, §9.5 AI 모델 드리프트 감지, BLE GATT 0xFF08 신규

**v2.4 추가사항**: §6.16 SelfHealingOrchestrator(HwAutoRecovery, FwSelfRepair, DataPipelineReprocessor, MeasurementSessionRecovery), §7.8 Flutter 앱 복원 UI, §8.7 Go 백엔드 회복탄력성, §9.6 AI 자가교정 엔진, §11.1 IEC 62304/ISO 14971 규제 프로세스 매핑, §11.2 FHIR R4 의료 데이터 표준, §11.3 WCAG 2.1 AA 접근성, §5.4 MCU 리소스 예산, BLE 0xFF09 Healing Event, BLE f32→f64 타입 수정

**v2.4.1 수정사항**: §2 디렉토리 트리 구조 오류 수정(context/ 들여쓰기, healing/ 마지막 자식 표기), §12.2 통합 테스트 시나리오 9~11 추가(자가검증·자가치유·에스컬레이션 E2E), §13 Phase 3-5 체크리스트에 Self-Healing 6개 항목 추가

**v2.4.2 수정사항**: 무결점 감사 결과 반영 — §6.3 핑거프린트 차원 확장 공식 수정(Critical: "88×8=448" 산술 오류 → 교차반응 특성 선택 파이프라인으로 명확화, Stage-1 56차원 유래 상세화), §6.15 SystemHealthOrchestrator IEEE 754 부동소수점 정규화 주석 추가(등급 경계 비교 시 ε=1e-9 허용 오차)

---

## 1. 프로젝트 개요 및 SSOT 상수

### 1.1 시스템 정의

만파식은 차동측정 기반 범용분석 POCT 시스템이다. 하드웨어 리더기(STM32F405)에 일회용 카트리지를 삽입하면, 9블록 Universal AFE가 측정 모드를 자동 전환하고, 88→896→1792차원 핑거프린트 벡터와 비표적 역추론 AI로 분류·정량·미지물질 동정을 수행한다.

### 1.2 SSOT 상수 (코드에 반드시 반영)

```rust
// ── 이 블록의 값은 CLAUDE.md v2.1 베이스라인이며 코드 내 하드코딩 ──
pub const CONNECTOR_MODEL: &str     = "Samtec MECF-08-01-L-DV";
pub const CONNECTOR_PINS: u8        = 16;           // CSI v1.0
pub const CONNECTOR_PITCH_MM: f64   = 1.27;
pub const PPM_WIDTH_MM: f64         = 49.70;
pub const PPM_DEPTH_MM: f64         = 30.0;
pub const PPM_HEIGHT_MM: f64        = 4.30;
pub const ALPHA_DEFAULT: f64        = 0.98;
pub const ALPHA_MIN: f64            = 0.90;
pub const ALPHA_MAX: f64            = 1.10;
pub const AFE_BLOCK_COUNT: u8       = 9;
pub const ACCURACY_MIN_PERCENT: f64 = 92.0;
pub const ACCURACY_MAX_PERCENT: f64 = 98.0;
pub const FINGERPRINT_DIM_STAGE1: usize = 88;
pub const FINGERPRINT_DIM_ENOSE: usize  = 448;
pub const FINGERPRINT_DIM_FULL: usize   = 896;
pub const FINGERPRINT_DIM_MAX: usize    = 1792;
pub const EXCHANGE_RATE_KRW_USD: f64    = 1480.0;

// NFC 매니페스트
pub const MANIFEST_V1_SIZE: usize = 172;
pub const MANIFEST_V2_SIZE: usize = 256;
pub const MANIFEST_MAGIC_V1: [u8; 4] = *b"MPK1";
pub const MANIFEST_MAGIC_V2: [u8; 4] = *b"MPK2";
```

### 1.3 금지 사항

```
FORBIDDEN:
- E12-IF 12핀 참조 (구 사양, CSI v1.0 16핀으로 변경 확정)
- unwrap() 호출 (테스트 코드 제외)
- AI 생성 가상 검증 결과 삽입
- setState() 직접 호출 (Flutter, Riverpod 사용)
- 하드코딩 키/시크릿
- WidthType.PERCENTAGE (docx 테이블)
```

### 1.4 기술 스택 (확정)

| 레이어 | 기술 | 버전 |
|---|---|---|
| Rust Core | Rust + no_std + embedded-hal | stable 2024 |
| Flutter App | Flutter + Riverpod + freezed + flutter_rust_bridge | 3.x / 2.x |
| Backend | Go + gRPC + Fiber(REST gateway) | 1.22+ |
| DB | PostgreSQL + TimescaleDB + Redis + Milvus + ES | 16 / 2.x |
| AI Edge | TFLite + XGBoost (INT8 양자화) | — |
| AI Cloud | Triton Inference Server + Flower | — |

### 1.5 특허 매핑

| 특허 | 번호 | 핵심 청구항 → 모듈 매핑 |
|---|---|---|
| Family A | APP2026-0022KR | cl.1→디지털트윈, cl.4-5→차동측정, cl.11→NFC, cl.19→다층AI, cl.24→연합학습, cl.30-31→SiPM-ECL |
| Family B | APP2025-0967KR | cl.1→범용분석, cl.3→다중모드, cl.8→RAFE, cl.11-13→AI동적가중치, cl.17-20→eNose |
| Family C | APP2025-0968KR | 모듈형 플랫폼, OTA, 보안 |

---

## 2. 프로젝트 구조

```
manpasik/
├── CLAUDE.md                       # → Appendix A 참조
├── rust-core/                      # Rust 핵심 엔진 (Cargo workspace)
│   ├── Cargo.toml                  # workspace members
│   ├── manpasik-core/              # 메인 라이브러리
│   │   └── src/
│   │       ├── lib.rs              # FFI 엔트리포인트
│   │       ├── harness/            # ★ HAL 추상화 계층
│   │       │   ├── sensor_trait.rs     # §4.1
│   │       │   ├── cartridge_manifest.rs # §4.2
│   │       │   ├── afe_registry.rs     # §4.3
│   │       │   └── compatibility_checker.rs
│   │       ├── signal/             # 신호처리
│   │       │   ├── differential.rs     # §6.1
│   │       │   ├── dsp.rs              # §6.2
│   │       │   ├── fft.rs
│   │       │   ├── peak_detector.rs
│   │       │   ├── baseline_correction.rs
│   │       │   └── noise_reduction.rs
│   │       ├── calibration/        # 보정
│   │       │   ├── hybrid_correction.rs
│   │       │   ├── kdm_drift.rs
│   │       │   ├── temperature_comp.rs
│   │       │   ├── humidity_comp.rs
│   │       │   ├── matrix_correction.rs
│   │       │   └── calibration_store.rs
│   │       ├── quantification/     # 정량화
│   │       │   ├── concentration_engine.rs
│   │       │   ├── kalman_filter.rs
│   │       │   ├── multi_output.rs
│   │       │   └── uncertainty.rs
│   │       ├── fingerprint/        # 핑거프린트
│   │       │   ├── feature_extractor.rs    # §6.3
│   │       │   ├── multi_mode_expander.rs
│   │       │   ├── enose_fusion.rs
│   │       │   ├── etongue_fusion.rs
│   │       │   └── vector_normalizer.rs
│   │       ├── ai/                 # AI 추론
│   │       │   ├── tflite_runtime.rs   # §6.5
│   │       │   ├── xgboost_inference.rs
│   │       │   ├── classification.rs
│   │       │   ├── non_target_detector.rs
│   │       │   ├── model_registry.rs
│   │       │   ├── ab_test_engine.rs
│   │       │   └── xai_explainer.rs
│   │       ├── digital_twin/       # 디지털 트윈
│   │       │   ├── twin_engine.rs      # §6.4
│   │       │   ├── drift_detector.rs
│   │       │   ├── calibration_predictor.rs
│   │       │   ├── cartridge_life.rs
│   │       │   └── optimizer.rs
│   │       ├── communication/      # 통신
│   │       │   ├── ble_protocol.rs
│   │       │   ├── nfc_handler.rs
│   │       │   ├── wifi_direct.rs
│   │       │   └── usb_cdc.rs
│   │       ├── security/           # 보안
│   │       │   ├── hash_chain.rs       # §6.7
│   │       │   ├── rolling_hash.rs
│   │       │   ├── aes_gcm.rs
│   │       │   ├── tpm_interface.rs
│   │       │   └── pqc_hybrid.rs
│   │       ├── data/               # 데이터
│   │       │   ├── local_store.rs
│   │       │   ├── crdt_sync.rs
│   │       │   ├── transform_log.rs
│   │       │   └── fhir_export.rs
│   │       ├── cartridge/          # 카트리지
│   │       │   ├── registry.rs
│   │       │   ├── validator.rs
│   │       │   ├── nfc_data_parser.rs
│   │       │   └── fallback_chain.rs
│   │       ├── sipm_ecl/           # Stage-2 준비
│   │       │   ├── saturation_detector.rs  # §6.6
│   │       │   ├── nonlinear_correction.rs
│   │       │   ├── self_optimization.rs
│   │       │   └── ecl_signal_processor.rs
│   │       ├── checkup/            # 종합검진 모듈 (v2.1 신규)
│   │       │   ├── session.rs           # §6.8
│   │       │   ├── disease_risk.rs      # §6.9
│   │       │   └── auto_classifier.rs   # §6.10
│   │       ├── context/            # 유기적 연동 모듈 (v2.2 신규)
│   │       │   ├── engine.rs            # §6.11 ContextEngine
│   │       │   ├── nutrition.rs         # §6.12 NutritionAdvisor
│   │       │   ├── shopping.rs          # §6.13 ShoppingBridge
│   │       │   └── habit.rs             # §6.14 HabitTracker
│   │       ├── diagnostics/        # 자가검증 모듈 (v2.3 신규)
│   │       │   ├── system_health.rs     # §6.15 SystemHealthOrchestrator
│   │       │   ├── hw_monitor.rs        # §6.15.1 HwHealthMonitor
│   │       │   ├── fw_watchdog.rs       # §6.15.2 FwWatchdogBridge
│   │       │   ├── data_integrity.rs    # §6.15.3 DataIntegrityChecker
│   │       │   ├── predictive.rs        # §6.15.4 PredictiveMaintenanceEngine
│   │       │   └── error_reporter.rs    # §6.15.5 ErrorReporter
│   │       └── healing/            # 자가치유 모듈 (v2.4 신규)
│   │           ├── orchestrator.rs      # §6.16 SelfHealingOrchestrator
│   │           ├── hw_recovery.rs       # §6.16.1 HwAutoRecovery
│   │           ├── fw_repair.rs         # §6.16.2 FwSelfRepair
│   │           ├── pipeline_reprocess.rs # §6.16.3 DataPipelineReprocessor
│   │           └── session_recovery.rs  # §6.16.4 MeasurementSessionRecovery
│   ├── manpasik-ffi/               # flutter_rust_bridge FFI
│   └── manpasik-fw/                # 펌웨어 (no_std, STM32 HAL)
├── flutter-app/                    # Flutter 모바일 앱
│   ├── lib/
│   │   ├── providers/              # Riverpod providers
│   │   ├── screens/                # 화면
│   │   ├── services/               # BLE/NFC/CRDT 서비스
│   │   ├── models/                 # freezed 데이터 모델
│   │   └── widgets/                # 재사용 위젯
│   └── test/
├── go-backend/                     # Go gRPC 백엔드
│   ├── proto/                      # Protobuf 정의
│   ├── cmd/                        # 서비스 엔트리포인트
│   ├── internal/                   # 도메인 로직
│   └── deploy/                     # Kubernetes manifests
├── ai-ml/                          # Python AI/ML 파이프라인
│   ├── training/
│   ├── serving/
│   └── federated/
└── docs/                           # 본 사양서 등
```

---

## 3. 아키텍처

### 3.1 6-Layer 계층 구조

```
┌─────────────────────────────────────────────────────────────────┐
│  Layer 6: 인프라                                                 │
│  AWS/GCP Multi-Cloud, Kubernetes, CI/CD, Monitoring              │
├─────────────────────────────────────────────────────────────────┤
│  Layer 5: 백엔드 서비스                                          │
│  Go + gRPC MSA (21 서비스), PostgreSQL + TimescaleDB + Milvus    │
├─────────────────────────────────────────────────────────────────┤
│  Layer 4: 모바일 앱                                              │
│  Flutter 3.x + Riverpod + flutter_rust_bridge (FFI)              │
├─────────────────────────────────────────────────────────────────┤
│  Layer 3: Rust 핵심 엔진                                         │
│  신호처리 + 차동측정 + AI 추론 + 보안 + 디지털 트윈              │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  ★ Harness Abstraction Layer (HAL)                        │  │
│  │  SensorTrait → 9블록 AFE 통합 인터페이스                   │  │
│  │  CartridgeManifest v2.0 → Stage별 확장 가능               │  │
│  └───────────────────────────────────────────────────────────┘  │
├─────────────────────────────────────────────────────────────────┤
│  Layer 2: 하드웨어 제어                                          │
│  embedded-hal, GPIO/SPI/I2C/ADC, RAFE 스위치, EHD 제어           │
├─────────────────────────────────────────────────────────────────┤
│  Layer 1: 하드웨어                                               │
│  STM32F405 + CSI v1.0 (16핀) + BLE nRF52832 + NFC PN7150        │
└─────────────────────────────────────────────────────────────────┘
```

### 3.2 End-to-End 데이터 흐름

```
카트리지 삽입 → [NFC PN7150] SUN/CMAC 인증 → CartridgeManifest v2.0 파싱
    → [HAL] SensorTrait::init() → AFE 블록 자동 선택
    → [RAFE] 측정 모드 자동 설정 (CA/CV/EIS/LSV/SWV/DPV/IMP/OCP)
    → [ADC ADS1256] 24-bit 디지털 변환 → DMA 고속 수집
    → [DifferentialEngine] S_det - α × S_ref
    → [DSP Pipeline] FFT + 피크 검출 + 칼만 필터 → 88→896차원 핑거프린트
    → [TFLite XGBoost] 분류 + 정량 + 비표적 이상탐지
    → [Digital Twin] 잔차 모니터링 → 드리프트 보정 / 교정 시점 예측
    → [BLE nRF52832] GATT Notify → Flutter UI 실시간 시각화
    → [CRDT Sync] 오프라인 우선 저장 → 클라우드 동기화
```

### 3.3 모듈 의존성 그래프

```
                    ┌──────────────────┐
                    │ CartridgeManifest│
                    │     (§4.2)       │
                    └────────┬─────────┘
                             │ 파싱 결과
                    ┌────────▼─────────┐
                    │   AfeRegistry    │
                    │     (§4.3)       │
                    └────────┬─────────┘
                             │ AFE 블록 선택
              ┌──────────────┼──────────────┐
              │              │              │
     ┌────────▼───┐  ┌──────▼──────┐ ┌────▼────────┐
     │ SensorTrait│  │ SensorTrait │ │ SensorTrait │
     │ Electrochem│  │   eNose     │ │  SiPM-ECL   │
     │  (Stage-1) │  │  (Stage-1)  │ │  (Stage-2)  │
     └────────┬───┘  └──────┬──────┘ └────┬────────┘
              │              │              │
              └──────────────┼──────────────┘
                             │ RawData
                    ┌────────▼─────────┐
                    │ DifferentialEngine│
                    │     (§6.1)        │
                    └────────┬─────────┘
                             │ DifferentialSignal
              ┌──────────────┼──────────────┐
              │              │              │
     ┌────────▼───┐  ┌──────▼──────┐ ┌────▼────────┐
     │ DSP Pipeline│  │FeatureExtract│ │ Digital Twin│
     │   (§6.2)    │  │   (§6.3)    │ │   (§6.4)    │
     └────────┬───┘  └──────┬──────┘ └────┬────────┘
              │              │              │
              └──────────────┼──────────────┘
                             │ 896-dim fingerprint + twin status
                    ┌────────▼─────────┐
                    │  TFLite AI       │
                    │  Inference (§6.5)│
                    └────────┬─────────┘
                             │ classification + quantification
                    ┌────────▼─────────┐
        ┌───────────→│  Flutter App     │←────────────┐
        │           │  (§7, BLE 경유)  │             │
        │           └──────────────────┘             │
        │                                             │
        └──────────────────────────────────────────┘
           SystemHealthOrchestrator (§6.15 v2.3)
           (HwHealthMonitor + FwWatchdog +
            DataIntegrity + Predictive +
            ErrorReporter - 6계층 전체 모니터링)
```

---

## 4. 핵심 타입 & 인터페이스 정의

> **이 섹션의 모든 trait/struct는 프로젝트 전체의 계약(Contract)이다.
> 모듈 구현 시 이 정의를 변경하지 않고 준수해야 한다.**

### 4.1 SensorTrait — 통합 HAL 인터페이스

**파일**: `rust-core/manpasik-core/src/harness/sensor_trait.rs`
**역할**: 9블록 AFE를 단일 인터페이스로 추상화. 새 AFE 추가 시 이 trait만 구현하면 전체 파이프라인에 자동 통합.
**특허**: Family B cl.8 (RAFE 재구성)

```rust
pub trait SensorTrait {
    type Config;
    type RawData;
    type ProcessedData;
    type Error: core::fmt::Debug;

    /// 센서 초기화 (카트리지 매니페스트 기반)
    fn init(&mut self, config: &Self::Config) -> Result<(), Self::Error>;

    /// 원시 데이터 읽기 (DMA 기반 고속 수집)
    fn read_raw(&self) -> Result<Self::RawData, Self::Error>;

    /// 차동측정 적용: S_det - alpha * S_ref
    fn apply_differential(
        &self,
        raw: &Self::RawData,
        alpha: f64,
    ) -> Result<Self::ProcessedData, Self::Error>;

    /// 보정 데이터 로드 (NFC → QR → Cloud → 범용 폴백 체인)
    fn load_calibration(
        &mut self,
        manifest: &CartridgeManifest,
    ) -> Result<CalibrationData, Self::Error>;

    /// 센서 상태 진단 (디지털 트윈 잔차 기반)
    fn self_diagnose(&self) -> SensorHealth;

    /// 안전 종료
    fn shutdown(&mut self) -> Result<(), Self::Error>;
}

// ── 9블록 AFE 구현체 선언 ──
pub struct ElectrochemAfe;      // Stage-1: LMP91000×4 + ADS1256
pub struct SipmEclAfe;          // Stage-2: SiPM-ECL 광학
pub struct ColorimetricAfe;     // Stage-2: LED/PD 색측정
pub struct TecAfe;              // Stage-2: 열전소자 온도 제어
pub struct LampAfe;             // Stage-3: LAMP/RPA 등온 NAAT
pub struct EnoseAfe;            // Stage-1: 전자코 8채널
pub struct EcArrayAfe;          // Stage-1: 전기화학 어레이 확장
pub struct RadiationAfe;        // Stage-4: 방사선 검출
pub struct MolecularAfe;        // Stage-5: 분자진단 통합

// ── ElectrochemAfe 타입 계약 ──
impl SensorTrait for ElectrochemAfe {
    type Config = ElectrochemConfig;
    type RawData = AdcSamples;
    type ProcessedData = DifferentialSignal;
    type Error = AfeError;
    // ... 구현
}
```

**테스트 기준 (5개 이상)**:
1. `ElectrochemAfe::init()` — 유효한 Config → `Ok(())`
2. `ElectrochemAfe::init()` — 잘못된 채널 → `Err(AfeError::InvalidChannel)`
3. `read_raw()` — 4채널 ADC 값이 24-bit 범위 내 (±8,388,607)
4. `apply_differential()` — compute(100.0, 50.0) with alpha=0.98 → 51.0
5. `self_diagnose()` — 정상 센서 → `SensorHealth::Good`
6. `shutdown()` — 이중 호출 시에도 에러 없음

---

### 4.2 CartridgeManifest v2.0 — NFC 데이터 구조

**파일**: `rust-core/manpasik-core/src/harness/cartridge_manifest.rs`
**역할**: NFC 태그에서 읽은 256바이트를 파싱하여 카트리지 메타데이터 제공
**특허**: Family A cl.11 (NFC)

```rust
/// NFC 태그 데이터 구조 v2.0 (256바이트)
/// v1.0(172바이트) 카트리지는 from_v1()으로 업캐스트
#[repr(C, packed)]
pub struct CartridgeManifest {
    // === v1.0 호환 영역 (0-171 bytes) ===
    pub magic_number: [u8; 4],          // 'MPK1' 또는 'MPK2'
    pub manifest_version: u8,           // 1 또는 2
    pub cartridge_type: u16,            // 유형 코드
    pub serial_number: [u8; 12],        // 고유 일련번호
    pub manufacturing_date: u32,        // 제조일 (UNIX timestamp)
    pub expiration_date: u32,           // 만료일 (UNIX timestamp)
    pub max_uses: u16,                  // 최대 사용 횟수
    pub current_uses: u16,              // 현재 사용 횟수
    pub calibration_data: [u8; 128],    // 보정 계수
    pub firmware_min_ver: u32,          // 최소 펌웨어 버전
    pub regulatory_code: u32,           // 규제 승인 코드
    pub crc32: u32,                     // CRC32 (v1.0 영역)

    // === v2.0 확장 영역 (172-255 bytes) ===
    pub csi_version: u8,                // 1=v1.0(16핀), 2=v2.0(24핀)
    pub stage: u8,                      // 1~5
    pub afe_blocks: u16,                // AFE 비트마스크 (9비트)
    pub measurement_modes: u16,         // 측정모드 비트마스크 (8비트)
    pub pin_config_hash: [u8; 8],       // 핀 설정 해시
    pub optical_wavelength_nm: u16,     // LED 파장 (0=미사용)
    pub optical_excitation_mv: u16,     // 여기 전압 (0=미사용)
    pub sipm_bias_v: u16,               // SiPM 바이어스 (0.1V 단위)
    pub lamp_temp_target: u16,          // LAMP 목표 온도 (0.1°C)
    pub lamp_duration_sec: u16,         // LAMP 반응 시간
    pub recommended_model_id: [u8; 16], // 권장 AI 모델 ID
    pub developer_id: [u8; 8],          // SDK 개발자 ID
    pub reserved: [u8; 20],             // 확장 예약
    pub crc32_v2: u32,                  // CRC32 (전체 256바이트)
}

impl CartridgeManifest {
    pub fn is_v1(&self) -> bool { self.magic_number == *b"MPK1" }
    pub fn is_v2(&self) -> bool { self.magic_number == *b"MPK2" }

    /// v1.0 → v2.0 업캐스트 (후방호환, H2 원칙)
    pub fn from_v1(v1_data: &[u8; 172]) -> Self { /* 확장 필드 기본값 */ }

    /// 바이트 배열에서 파싱
    pub fn parse(data: &[u8]) -> Result<Self, ManifestError> { /* CRC 검증 포함 */ }

    /// 만료 여부 확인
    pub fn is_expired(&self, now_unix: u32) -> bool { /* expiration_date < now */ }

    /// 사용 가능 여부 (만료 + 횟수)
    pub fn is_usable(&self, now_unix: u32) -> bool { /* !expired && uses < max */ }

    /// AFE 비트마스크 → 활성 블록 목록
    pub fn active_afe_blocks(&self) -> Vec<AfeBlockType> { /* 비트 해석 */ }
}
```

**AFE 비트마스크 정의**:
```rust
pub enum AfeBlockType {
    Electrochem  = 0b0_0000_0001,  // bit 0
    SipmEcl      = 0b0_0000_0010,  // bit 1
    Colorimetric = 0b0_0000_0100,  // bit 2
    Tec          = 0b0_0000_1000,  // bit 3
    Lamp         = 0b0_0001_0000,  // bit 4
    Enose        = 0b0_0010_0000,  // bit 5
    EcArray      = 0b0_0100_0000,  // bit 6
    Radiation    = 0b0_1000_0000,  // bit 7
    Molecular    = 0b1_0000_0000,  // bit 8
}
```

**테스트 기준**:
1. v2.0 파싱: 256바이트 정상 → 모든 필드 정확
2. v1.0 후방호환: 172바이트 + `from_v1()` → csi_version=1, stage=1
3. CRC 오류: 변조된 데이터 → `Err(ManifestError::CrcMismatch)`
4. 만료 검증: 과거 날짜 → `is_expired() == true`
5. 사용 횟수: current_uses >= max_uses → `is_usable() == false`
6. AFE 비트마스크: `0b00100001` → `[Electrochem, Enose]`
7. 매직넘버 불일치 → `Err(ManifestError::InvalidMagic)`

---

### 4.3 AfeRegistry — AFE 블록 동적 등록

**파일**: `rust-core/manpasik-core/src/harness/afe_registry.rs`
**역할**: 런타임에 AFE 구현체를 등록하고, 카트리지 매니페스트에 따라 자동 선택
**원칙**: H1 모듈 독립성 — 새 AFE 추가 시 레지스트리에 등록만 하면 됨

```rust
pub struct AfeRegistry {
    handlers: HashMap<AfeBlockType, Box<dyn SensorTrait<
        Config = dyn Any,
        RawData = dyn Any,
        ProcessedData = dyn Any,
        Error = AfeError,
    >>>,
}

impl AfeRegistry {
    pub fn new() -> Self { /* 빈 레지스트리 */ }

    pub fn register(&mut self, block: AfeBlockType, handler: impl SensorTrait + 'static) { }

    /// 카트리지 매니페스트 기반 자동 선택
    pub fn select_for_cartridge(
        &self,
        manifest: &CartridgeManifest,
    ) -> Result<Vec<&dyn SensorTrait>, RegistryError> { }

    /// 등록된 블록 목록
    pub fn available_blocks(&self) -> Vec<AfeBlockType> { }
}
```

**테스트 기준**:
1. 빈 레지스트리에서 조회 → `Err(RegistryError::NoHandler)`
2. Electrochem 등록 후 Stage-1 매니페스트 → 해당 핸들러 반환
3. 미등록 Stage-2 매니페스트 → SipmEcl 부분 에러 (부분 성공 허용, H6)

---

### 4.4 공통 타입 정의

```rust
// ── 환경 데이터 ──
pub struct EnvironmentData {
    pub temperature: f64,       // °C
    pub humidity: f64,          // %
    pub sensor_age_hours: u32,  // 센서 누적 사용 시간
}

// ── 차동측정 결과 ──
pub struct DifferentialResult {
    pub timestamp_ms: u64,
    pub channel: u8,
    pub s_detection: f64,
    pub s_reference: f64,
    pub s_differential: f64,
    pub alpha_used: f64,
}

// ── 센서 건강 상태 ──
pub enum SensorHealth {
    Good,           // 정상
    Degraded(f64),  // 성능 저하 (잔차 비율 0.0~1.0)
    Replace,        // 교체 필요
}

// ── 보정 데이터 ──
pub struct CalibrationData {
    pub slope: f64,
    pub intercept: f64,
    pub temp_coefficient: f64,
    pub humid_coefficient: f64,
    pub valid_until: u32,       // UNIX timestamp
}

// ── 드리프트 수준 ──
pub enum DriftLevel {
    Normal,
    Warning,            // |잔차| > 3σ
    LongTermDrift,      // Allan 분산 τ > threshold
    SystematicBias,     // 연속 5회 단방향 잔차
}

// ── AI 추론 결과 ──
pub struct InferenceResult {
    pub classification: String,           // 물질 분류
    pub confidence: f64,                  // 분류 확신도
    pub concentrations: Vec<(String, f64)>, // (물질명, 농도) 최대 10종
    pub anomaly_score: f64,               // 비표적 이상탐지 점수
    pub model_id: String,                 // PCCP 추적용
    pub model_hash: String,               // SHA-256
}

// ── 에러 타입 ──
#[derive(Debug)]
pub enum AfeError {
    InitFailed(String),
    ReadTimeout,
    InvalidChannel(u8),
    CalibrationExpired,
    CommunicationError(String),
    SaturationDetected,         // SiPM-ECL 포화
}

pub enum ManifestError {
    InvalidMagic,
    CrcMismatch,
    UnsupportedVersion(u8),
    DataTooShort(usize),
}
```

---

## 5. Layer 1-2: 하드웨어 + 펌웨어

### 5.1 CSI v1.0 커넥터 핀맵 (확정)

```
커넥터: Samtec MECF-08-01-L-DV, 16핀, 1.27mm 피치
```

| 핀 | 신호명 | 방향 | 설명 | Stage |
|---|---|---|---|---|
| 1 | WE1 | Analog In | 감지전극 1 | 1 |
| 2 | RE1 | Analog In | 참조전극 1 | 1 |
| 3 | CE1 | Analog Out | 대향전극 1 | 1 |
| 4 | WE2 | Analog In | 감지전극 2 (차동 쌍) | 1 |
| 5 | RE2 | Analog In | 참조전극 2 (차동 쌍) | 1 |
| 6 | CE2 | Analog Out | 대향전극 2 | 1 |
| 7 | TEMP | Analog In | NTC 온도 센서 | 1 |
| 8 | HUMID | Analog In | 습도 센서 | 1 |
| 9 | EHD_HV | Power | EHD 고전압 출력 | 1 |
| 10 | GND | Power | 접지 | 1 |
| 11 | VCC_3V3 | Power | 3.3V 전원 | 1 |
| 12 | SPI_CLK | Digital | SPI 클록 | 1 |
| 13 | SPI_MOSI | Digital | SPI 데이터 출력 | 1 |
| 14 | SPI_MISO | Digital | SPI 데이터 입력 | 1 |
| 15 | RESERVED_1 | — | Stage-2 광학 신호 예약 | 2 |
| 16 | RESERVED_2 | — | Stage-2 광학 전원 예약 | 2 |

**진화 로드맵**: CSI v1.0(16핀) → v2.0(24핀, Stage-2) → v3.0(32핀, Stage-3+). 핀 15-16은 Stage-1에서 NC(미연결), Stage-2 카트리지 삽입 시만 활성화.

### 5.2 MCU 펌웨어 태스크

```
STM32F405RGT6 (ARM Cortex-M4, 168MHz, 1MB Flash, 192KB RAM)

FreeRTOS Tasks:
├── ADC_Sampler     (Priority: Highest, 1kHz)      — ADS1256 SPI+DMA
├── BLE_Comm        (Priority: High)                — nRF52832 UART/SPI
├── NFC_Handler     (Priority: Medium)              — PN7150 I2C, SUN/CMAC
├── EHD_Controller  (Priority: Medium)              — PWM 고전압 펄스
├── Power_Manager   (Priority: Low)                 — BQ24195 충전/전원
└── OTA_Updater     (Priority: Lowest, 백그라운드)  — A/B 파티션 업데이트
```

### 5.3 BLE GATT v2.0 서비스 구조

**Service UUID**: `0000FF00-0000-1000-8000-00805F9B34FB`

| Characteristic | UUID | 권한 | 페이로드 |
|---|---|---|---|
| Configuration Control | 0xFF01 | Write | `measurement_mode: u8` + `sampling_rate: u16` + `afe_channel: u8` + `ehd_enable: bool` |
| Waveform Stream | 0xFF02 | Read/Notify | `timestamp: u32` + `channel: u8` + `raw_adc: i32` + `differential: f64` + `sequence_num: u16` |
| System Status | 0xFF03 | Read/Notify | `battery: u8` + `temp: i16` + `humidity: u16` + `afe_status: u8` + `error: u16` + `fw_ver: u32` |
| EHD Control | 0xFF04 | Write | `voltage_kv: u16` + `pulse_width_ms: u16` + `flow_direction: u8` |
| Security Hash | 0xFF05 | Read | `rolling_hash: [u8;32]` + `hash_index: u32` |
| OTA Control (v2.0) | 0xFF06 | Write/Notify | `command: u8` + `chunk_data: [u8;244]` + `chunk_index: u32` + `total_chunks: u32` |
| Digital Twin (v2.0) | 0xFF07 | Read/Notify | `predicted: f32` + `residual: f32` + `drift_score: f32` + `calibration_due: u32` + `health: u8` |
| System Health (v2.3) | 0xFF08 | Read/Notify | `overall_score: f32` + `grade: u8` + `hw_score: f32` + `fw_score: f32` + `rust_score: f32` + `measurement_ready: bool` + `active_alerts: u8` |
| Healing Event (v2.4) | 0xFF09 | Read/Notify | `event_type: u8` + `layer: u8` + `strategy: u8` + `outcome: u8` + `duration_ms: u16` + `timestamp: u32` |

**UUID 할당 현황** (v2.4 기준): 0xFF01~0xFF09 사용 중, 0xFF0A~0xFF0F 예약(향후 확장).

### 5.4 MCU 리소스 예산 (v2.4 보강)

**STM32F405 리소스 제약**: 168MHz Cortex-M4, 192KB SRAM, 1MB Flash

| 모듈 | SRAM 예산 | CPU 예산 (측정 중) | 비고 |
|---|---|---|---|
| ADC 샘플링 (1kHz, 9ch) | 36KB (DMA 더블 버퍼) | 15% | DMA 기반, CPU 최소 |
| DSP 파이프라인 (§6.2) | 24KB (FFT 버퍼) | 25% | Butterworth + Kalman |
| 차분식 엔진 (§6.1) | 4KB | 5% | f64 연산 |
| BLE 스택 + GATT | 32KB | 10% | SoftDevice 고정 |
| SelfDiagnostics (§6.15) | 8KB (점수 캐시) | 5% (백그라운드) | 측정 중 경량 모니터링만 |
| SelfHealing (§6.16) | 12KB (정책 + 로그) | 5% (이벤트 구동) | **측정 중 L0~L1만 허용** |
| FW 관리 + 기타 | 16KB | 10% | RTOS, HAL |
| **합계** | **132KB / 192KB (69%)** | **75%** | 여유: 60KB SRAM, 25% CPU |
| **Flash 사용** | 측정 코드 512KB + OTA A/B 각 256KB = **1024KB / 1024KB** | - | 여유 없음, 코드 최적화 필수 |

**제약 사항**: Flash가 거의 포화 상태이므로, SelfDiagnostics와 SelfHealing의 펌웨어 구현은 Rust core(모바일 측)에서 주로 실행하고, MCU는 최소 데이터 수집 + BLE 전송만 담당하는 분산 아키텍처를 채택한다. AI 추론(TFLite/XGBoost)은 모바일 디바이스에서 실행한다.

---

## 6. Layer 3: Rust Core Engine — 모듈별 상세 사양

### 6.1 차동측정 엔진 (DifferentialEngine)

**파일**: `signal/differential.rs`
**의존**: §4.4 공통 타입 (EnvironmentData, DifferentialResult)
**특허**: Family A cl.4-5 (차동측정), Family B cl.11-13 (AI 동적 가중치)

#### 인터페이스 계약

```rust
pub struct DifferentialEngine {
    alpha: f64,                          // 기본 0.98
    alpha_range: (f64, f64),             // (0.90, 1.10) 안전 범위
    history: VecDeque<DifferentialResult>, // 최근 N개 결과
}

impl DifferentialEngine {
    pub fn new() -> Self;

    /// 단일 채널 차동 연산
    /// 수식: s_detection - alpha * s_reference
    pub fn compute(&self, s_detection: f64, s_reference: f64) -> f64;

    /// 4채널 동시 차동 (CSI v1.0 16핀 WE1/WE2 + 2쌍 추가)
    pub fn compute_array(
        &self,
        detections: &[f64; 4],
        references: &[f64; 4],
    ) -> [f64; 4];

    /// 동적 alpha 조정
    /// AI 가중치 있으면 우선 사용 (Family B cl.11-13)
    /// 없으면 환경 보정:
    ///   temp_factor  = 1.0 + (temperature - 25.0) * 0.001
    ///   humid_factor = 1.0 + (humidity - 50.0) * 0.0005
    ///   aging_factor = 1.0 - sensor_age_hours * 0.00001
    ///   new_alpha = alpha * temp_factor * humid_factor * aging_factor
    ///   clamp to [0.90, 1.10]
    pub fn update_alpha(
        &mut self,
        env: &EnvironmentData,
        ai_weight: Option<f64>,
    ) -> f64;

    /// SNR 향상도 계산 (dB)
    /// 목표: ≥13.2 dB
    pub fn compute_snr_improvement(
        &self,
        s_det: f64,
        s_ref: f64,
        noise_rms: f64,
    ) -> f64;

    /// 매트릭스 간섭 제거율
    /// 수식: (before - after) / before * 100
    /// 목표: 92~96%
    pub fn matrix_removal_rate(&self, before: f64, after: f64) -> f64;
}
```

#### 테스트 기준 (10개)

| # | 테스트 | 입력 | 기대 출력 |
|---|---|---|---|
| 1 | 기본 차동 | compute(100.0, 50.0), alpha=0.98 | 51.0 |
| 2 | 완벽 공통모드 제거 | compute(50.0, 50.0), alpha=1.0 | 0.0 |
| 3 | 4채널 배열 | 4쌍 입력 | 4개 결과 각각 검증 |
| 4 | alpha 하한 클램핑 | 환경보정 결과 0.85 | 0.90으로 클램핑 |
| 5 | alpha 상한 클램핑 | 환경보정 결과 1.15 | 1.10으로 클램핑 |
| 6 | 환경보정 중립 | 25°C, 50%, age=0 | alpha 변화 없음(0.98) |
| 7 | 고온 보정 | 40°C | alpha > 0.98 |
| 8 | AI 가중치 우선 | ai_weight=Some(0.95) | alpha=0.95 (환경 무시) |
| 9 | SNR 향상 | 적절한 입력 | ≥13.2 dB |
| 10 | 매트릭스 제거율 | before=100, after=6 | 94% (92~96% 범위) |

---

### 6.2 DSP 파이프라인

**파일**: `signal/dsp.rs`, `signal/fft.rs`, `signal/peak_detector.rs`
**의존**: §6.1 DifferentialEngine 출력

#### 인터페이스 계약

```rust
pub struct DspPipeline {
    butterworth: ButterworthFilter,
    savgol: SavitzkyGolayFilter,
    kalman: AdaptiveKalmanFilter,
}

impl DspPipeline {
    /// 전체 파이프라인 실행
    /// DC 제거 → Butterworth BPF → Savitzky-Golay 평활 → 칼만 필터
    pub fn process(&mut self, raw: &[f64]) -> Result<ProcessedSignal, DspError>;
}

pub struct FftAnalyzer;
impl FftAnalyzer {
    /// 실수 FFT (Hanning 윈도우 적용)
    pub fn compute_rfft(&self, signal: &[f64]) -> Vec<Complex64>;

    /// 주파수 영역 피크 검출
    pub fn find_peaks(
        &self,
        spectrum: &[Complex64],
        threshold_db: f64,
    ) -> Vec<FrequencyPeak>;
}

pub struct PeakDetector;
impl PeakDetector {
    /// SWV/DPV 피크 검출 (2차 도함수 + 적응형 임계값)
    pub fn detect_peaks(
        &self,
        voltammogram: &[(f64, f64)],  // (전위, 전류)
        mode: MeasurementMode,
    ) -> Vec<ElectrochemPeak>;
}
```

**테스트 기준**:
1. Butterworth 50Hz 노이즈 → 30dB 이상 감쇠
2. 정현파 입력 → FFT 피크가 입력 주파수에 정확히 위치
3. SWV 피크 검출 → 알려진 농도에서 ±5% 이내 피크 전위
4. 칼만 필터 수렴 → 20 샘플 이내 안정화

---

### 6.3 핑거프린트 특성 추출기

**파일**: `fingerprint/feature_extractor.rs`, `fingerprint/multi_mode_expander.rs`
**의존**: §6.1 차동결과, §6.2 DSP 결과
**특허**: Family B cl.3 (다중모드)

#### 차원 확장 파이프라인

```
[Stage-1 기본] — 88차원
  기본 신호:  4 센서페어 × 8 모드(CA,CV,EIS,LSV,SWV,DPV,IMP,OCP) = 32
  시간영역:   모드별 가변 특성 추출 (아래 상세) = 56
              ├── CV:  4특성(피크전위,기울기,면적,시정수) × 4페어 = 16
              ├── CA:  2특성(피크전류,시정수)             × 4페어 =  8
              ├── EIS: 2특성(임피던스크기,위상각)         × 4페어 =  8
              ├── LSV: 2특성(피크전위,기울기)             × 4페어 =  8
              ├── SWV: 1특성(피크전류)                    × 4페어 =  4
              ├── DPV: 1특성(피크전류)                    × 4페어 =  4
              ├── IMP: 1특성(임피던스크기)                × 4페어 =  4
              └── OCP: 1특성(안정전위)                    × 4페어 =  4
                                              소계 = 16+8+8+8+4+4+4+4 = 56
  합계: 32 + 56 = 88차원 ✓

[전자코/전자혀 융합] — 448 + 448 = 896차원
  eNose 8채널 교차반응 특성 융합:
    ├── 8채널 개별 응답 × 88 기본 = 8개 채널 프로파일
    ├── C(8,2) = 28 채널쌍 교차반응 패턴
    └── 특성 선택(상호정보량 기반) → 448차원 (SSOT: FINGERPRINT_DIM_ENOSE)
  ※ 단순 곱(88×8=704)이 아닌 교차반응 패턴 추출 + 특성 선택 파이프라인
  eTongue도 동일 구조 → 448차원
  eNose(448) + eTongue(448) = 896차원 (SSOT: FINGERPRINT_DIM_FULL) ✓

[Stage-2+ 확장] — 1,792차원
  896 + 광학(SiPM-ECL) 256 + LAMP 형광 128 + 임피던스 512 = 1,792차원
  (SSOT: FINGERPRINT_DIM_MAX) ✓
```

#### 인터페이스 계약

```rust
pub struct FeatureExtractor {
    stage: u8,
}

impl FeatureExtractor {
    /// Stage-1: 88차원 기본 특성 추출
    pub fn extract_base_features(
        &self,
        differential_signals: &[DifferentialSignal],
        modes: &[MeasurementMode],
    ) -> Result<[f64; 88], FeatureError>;

    /// 전자코 융합으로 896차원 확장
    pub fn expand_with_enose(
        &self,
        base: &[f64; 88],
        enose_data: &EnoseData,
    ) -> Result<Vec<f64>, FeatureError>;  // 896차원

    /// L2 정규화
    pub fn normalize_l2(vector: &mut [f64]);
}
```

**테스트 기준**:
1. 88차원 추출 → 벡터 길이 정확히 88
2. L2 정규화 후 벡터 크기 = 1.0 (±1e-10)
3. 전자코 융합 → 896차원 벡터 반환
4. 올제로 입력 → 적절한 에러 처리 (0으로 나누기 방지)

---

### 6.4 디지털 트윈 엔진

**파일**: `digital_twin/twin_engine.rs`, `drift_detector.rs`, `calibration_predictor.rs`, `optimizer.rs`
**의존**: §6.1 DifferentialResult
**특허**: Family A cl.1 (자기최적화)

#### 아키텍처

```
실측(y) ─┐
          ├─→ 잔차(y-ŷ) ─┬─→ 드리프트 감지
예측(ŷ) ─┘               ├─→ 교정 시점 예측 (남은 측정 횟수)
                          ├─→ 카트리지 수명 추정
                          └─→ 자기최적화 피드백 (α 조정, 보정 갱신)
```

#### 인터페이스 계약

```rust
// ── twin_engine.rs ──
pub struct DigitalTwinEngine {
    physical_model: PhysicalModel,
    residual_history: VecDeque<Residual>,
    drift_state: DriftState,
}

impl DigitalTwinEngine {
    pub fn update(&mut self, measured: f64) -> TwinUpdate;
}

pub struct TwinUpdate {
    pub predicted: f64,
    pub residual: f64,
    pub drift_score: f64,         // 0.0~1.0
    pub calibration_due: u32,     // 남은 측정 횟수
    pub health: SensorHealth,
}

// ── drift_detector.rs ──
pub struct DriftDetector;
impl DriftDetector {
    /// 드리프트 판정
    /// |잔차| > 3σ → Warning
    /// Allan 분산 τ > threshold → LongTermDrift
    /// 연속 5회 단방향 → SystematicBias
    pub fn detect(&mut self, residual: f64) -> DriftLevel;
}

// ── calibration_predictor.rs ──
pub struct CalibrationPredictor;
impl CalibrationPredictor {
    /// 잔차 추세 선형 회귀 → 허용 오차 초과 시점 추정
    pub fn predict_remaining(
        &self,
        residuals: &[Residual],
    ) -> CalibrationPrediction;
}

pub struct CalibrationPrediction {
    pub remaining_measurements: u32,
    pub remaining_days: f64,
    pub confidence: f64,        // 0.0~1.0
}

// ── optimizer.rs ──
pub struct SelfOptimizer;
impl SelfOptimizer {
    /// 잔차 기반 자기최적화 (Family A cl.1)
    pub fn optimize(&mut self, twin: &TwinUpdate) -> OptimizationAction;
}

pub enum OptimizationAction {
    AdjustAlpha(f64),
    UpdateCalibration(CalibrationData),
    NoAction,
}
```

**테스트 기준**:
1. 잔차=0 (완벽 예측) → `DriftLevel::Normal`
2. 큰 잔차 (>3σ) → `DriftLevel::Warning`
3. 연속 5회 상승 잔차 → `SystematicBias`
4. 교정 시점 예측 — 잔차 증가 추세 → 남은 측정 횟수 양수
5. 자기최적화 — 편향 감지 → `AdjustAlpha` 반환
6. BLE GATT 0xFF07 직렬화 → 정확한 바이트 변환

---

### 6.5 TFLite AI 추론 엔진

**파일**: `ai/tflite_runtime.rs`, `ai/classification.rs`, `ai/non_target_detector.rs`
**의존**: §6.3 핑거프린트 벡터 (896차원)
**특허**: Family A cl.19 (다층 AI)

#### 인터페이스 계약

```rust
pub struct TfLiteEngine {
    model: TfLiteModel,       // INT8 양자화, <2MB
    model_id: String,         // PCCP 추적용
    model_hash: String,       // SHA-256
}

impl TfLiteEngine {
    /// 모델 로드 (model_id + hash 검증)
    pub fn load(model_bytes: &[u8], expected_hash: &str) -> Result<Self, ModelError>;

    /// 분류 + 정량 추론 (<1ms 목표)
    pub fn predict(&self, fingerprint: &[f64]) -> Result<InferenceResult, ModelError>;
}

pub struct NonTargetDetector;
impl NonTargetDetector {
    /// Mahalanobis 거리 기반 비표적 이상탐지
    /// score > threshold → 미지 물질 가능성
    pub fn detect_anomaly(
        &self,
        fingerprint: &[f64],
        reference_distribution: &Distribution,
    ) -> f64;  // anomaly score
}

pub struct AbTestEngine;
impl AbTestEngine {
    /// A/B 테스트: 기존 모델 vs 새 모델 병렬 추론
    /// 두 결과를 비교하여 새 모델 성능 검증
    pub fn run_parallel(
        &self,
        control: &TfLiteEngine,
        treatment: &TfLiteEngine,
        input: &[f64],
    ) -> AbTestResult;
}
```

**테스트 기준**:
1. 모델 로드 — hash 일치 → `Ok`, 불일치 → `Err(ModelError::HashMismatch)`
2. 896차원 입력 → InferenceResult 반환 (classification + concentrations)
3. 추론 시간 < 5ms (엣지 기기 기준)
4. 비표적 탐지 — 훈련 데이터 내 샘플 → 낮은 anomaly score
5. 비표적 탐지 — 미지 물질 → 높은 anomaly score

---

### 6.6 SiPM-ECL 포화보정 (Stage-2 준비)

**파일**: `sipm_ecl/saturation_detector.rs`, `sipm_ecl/nonlinear_correction.rs`
**특허**: Family A cl.30-31 (SiPM-ECL)

#### 인터페이스 계약

```rust
pub struct SaturationDetector {
    microcell_count: u32,       // SiPM 마이크로셀 수 (예: 14,410)
    dead_time_ns: f64,          // 회복 시간 (예: 21ns)
    saturation_threshold: f64,  // 포화 판정 임계값 (예: 0.7)
}

pub enum SaturationLevel {
    Linear,      // ratio < 0.3
    SubLinear,   // 0.3 ≤ ratio < 0.7
    Saturated,   // ratio ≥ 0.7
}

impl SaturationDetector {
    pub fn detect(&self, measured_photons: u32) -> SaturationLevel;
}

/// 4PL 시그모이드 비선형 보정
/// N_true = -N_max × ln(1 - N_measured/N_max) × correction_factor
pub fn correct_saturation(
    measured: f64,
    n_max: f64,
    correction_factor: f64,
) -> Result<f64, SipmError>;

pub enum SipmError {
    CompleteSaturation,  // ratio ≥ 1.0
    InvalidInput,
}
```

**테스트 기준**:
1. 낮은 광자수 → `SaturationLevel::Linear`
2. 중간 → `SubLinear`
3. 포화 → `Saturated`
4. ratio=0.5 보정 → 수학적으로 정확한 N_true
5. ratio≥1.0 → `Err(CompleteSaturation)`

---

### 6.15 자가검증 시스템 (SelfDiagnostics Engine) — v2.3 신규

**파일**: `diagnostics/system_health.rs`, `hw_monitor.rs`, `fw_watchdog.rs`, `data_integrity.rs`, `predictive.rs`, `error_reporter.rs`
**의존**: §4.1 SensorTrait, §6.4 DigitalTwinEngine, §6.7 HashChain
**특허**: Family A cl.1 (자기최적화), Family C (모듈형 플랫폼)

#### 6-Layer 자가검증 아키텍처

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    SystemHealthOrchestrator                             │
│    ┌──────────────────────────────────────────────────────────────┐    │
│    │  종합 건강 점수 = Σ(layer_score × weight) / Σ(weight)       │    │
│    │  L1:HW(0.30) + L2:FW(0.15) + L3:Rust(0.25)                 │    │
│    │  + L4:App(0.10) + L5:Backend(0.10) + L6:AI(0.10)           │    │
│    └──────────────────────────────────────────────────────────────┘    │
│                                                                         │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌──────┐│
│  │L1: HW   │ │L2: FW   │ │L3: Rust │ │L4: App  │ │L5: Back │ │L6: AI││
│  │Monitor  │ │Watchdog │ │Integrity│ │Health   │ │Health   │ │Drift ││
│  │(§6.15.1)│ │(§6.15.2)│ │(§6.15.3)│ │(§7.7)   │ │(§8.6)   │ │(§9.5)││
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘ └──────┘│
│                                                                         │
│  ┌──────────────────────────────────┐  ┌───────────────────────────┐   │
│  │ PredictiveMaintenanceEngine     │  │ ErrorReporter              │   │
│  │ (§6.15.4)                       │  │ (§6.15.5)                  │   │
│  │ 예방정비 스케줄링 + 고장 예측    │  │ 근본원인 분석 + 자동 보고  │   │
│  └──────────────────────────────────┘  └───────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────┘
```

#### 6.15.1 하드웨어 건강 모니터 (HwHealthMonitor)

**역할**: Layer 1-2 전체 HW 컴포넌트 상태를 측정 전/중/후 3단계로 점검

```rust
pub struct HwHealthMonitor {
    component_states: HashMap<HwComponent, ComponentHealth>,
    history: VecDeque<HwCheckResult>,
    thresholds: HwThresholds,
}

#[derive(Clone, Copy, Hash, Eq, PartialEq)]
pub enum HwComponent {
    Battery,           // BQ24195 — 전압, 충전율, 사이클수
    BleRadio,          // nRF52832 — RSSI, 패킷손실률, 연결안정성
    NfcReader,         // PN7150 — 응답시간, 인증성공률
    AdcMain,           // ADS1256 — 노이즈플로어, INL/DNL, 샘플링레이트
    TecController,     // 열전소자 — PID 안정도, 오버슈트
    AfeBlock(u8),      // 9블록 AFE 개별 — 임피던스, 누설전류
    EhdModule,         // EHD 고전압 — 출력안정도
    TemperatureSensor, // NTC — 응답시간, 정확도
    HumiditySensor,    // 습도센서 — 응답시간, 드리프트
}

pub struct ComponentHealth {
    pub component: HwComponent,
    pub status: ComponentStatus,
    pub score: f64,              // 0.0~1.0
    pub metrics: Vec<(String, f64, f64)>,  // (metric_name, current_value, pass_threshold) — 예: ("adc_noise_uV", 12.3, 50.0)
    pub last_checked: u64,       // UNIX timestamp
    pub degradation_rate: f64,   // 성능 저하율 (%/day)
    pub predicted_failure: Option<u32>,  // 예상 고장까지 남은 일수
}

#[derive(Clone, Copy)]
pub enum ComponentStatus {
    Healthy,       // 모든 지표 정상 범위
    Warning,       // 1개 이상 지표 경계값 (threshold 80~100%)
    Degraded,      // 1개 이상 지표 임계 초과, 측정 가능하나 정확도 저하 가능
    Failed,        // 측정 불가
    Unknown,       // 점검 미수행 또는 통신 불가
}

pub struct HwCheckResult {
    pub timestamp: u64,
    pub phase: CheckPhase,
    pub components: Vec<ComponentHealth>,
    pub overall_score: f64,
    pub pass: bool,              // 측정 진행 가능 여부
    pub blockers: Vec<HwComponent>,  // 측정 차단 사유
}

pub enum CheckPhase {
    PreMeasurement,   // 측정 전 — 전체 컴포넌트 점검 (2초 이내)
    DuringMeasurement, // 측정 중 — ADC/AFE/BLE 실시간 감시
    PostMeasurement,  // 측정 후 — 결과 신뢰도 평가 + 다음 측정 준비상태
}

impl HwHealthMonitor {
    /// 측정 전 전체 HW 점검 (2초 이내 완료)
    /// Battery ≥ 10% + BLE RSSI ≥ -80dBm + ADC 노이즈 ≤ 2LSB
    /// + 선택된 AFE 블록 임피던스 정상 → pass=true
    pub fn pre_check(&mut self) -> HwCheckResult;

    /// 측정 중 실시간 감시 (100ms 주기)
    /// ADC 오버플로우, BLE 패킷 손실, 온도 급변 감지
    pub fn monitor_during(&mut self, adc_sample: i32, ble_rssi: i8) -> Option<HwAlert>;

    /// 측정 후 결과 신뢰도 평가
    /// 측정 중 발생한 alert 수 + 환경 변동폭 → 신뢰도 점수
    pub fn post_evaluate(&mut self, alerts_during: &[HwAlert]) -> ReliabilityScore;

    /// 개별 컴포넌트 상세 진단
    pub fn diagnose_component(&mut self, component: HwComponent) -> ComponentHealth;

    /// 전체 HW 건강 점수 (L1 계층 점수)
    pub fn layer_score(&self) -> f64;
}

pub struct HwAlert {
    pub component: HwComponent,
    pub severity: AlertSeverity,
    pub message: String,
    pub timestamp: u64,
    pub auto_recoverable: bool,
}

pub enum AlertSeverity {
    Info,       // 로깅만
    Warning,    // 사용자 알림, 측정 계속
    Critical,   // 측정 중단 권고
    Fatal,      // 즉시 중단
}

pub struct ReliabilityScore {
    pub score: f64,       // 0.0~1.0
    pub factors: Vec<(String, f64)>,  // ("ADC 안정도", 0.95) 등
    pub recommendation: MeasurementRecommendation,
}

pub enum MeasurementRecommendation {
    Reliable,           // 결과 신뢰 가능
    AcceptWithCaveat(String),  // 조건부 수용 (사유 명시)
    RecommendRetry,     // 재측정 권고
    Discard,            // 결과 폐기
}
```

**테스트 기준 (8개)**:
1. `pre_check()` — 배터리 50%, BLE -60dBm, ADC 정상 → `pass=true`, `overall_score > 0.8`
2. `pre_check()` — 배터리 5% → `pass=false`, `blockers=[Battery]`
3. `pre_check()` — BLE RSSI -90dBm → `ComponentStatus::Warning`, `pass=true` (측정 가능하나 경고)
4. `monitor_during()` — ADC 값 ±8,388,607 (포화) → `HwAlert { severity: Critical }`
5. `monitor_during()` — BLE 패킷 3연속 손실 → `HwAlert { severity: Warning, auto_recoverable: true }`
6. `post_evaluate()` — alert 0건 → `ReliabilityScore { score: 1.0, recommendation: Reliable }`
7. `post_evaluate()` — Critical alert 2건 → `recommendation: RecommendRetry`
8. `diagnose_component(AdcMain)` — INL > 2LSB → `ComponentStatus::Degraded`, `degradation_rate > 0`

---

#### 6.15.2 펌웨어 워치독 브릿지 (FwWatchdogBridge)

**역할**: STM32 IWDG/WWDG + FreeRTOS 태스크 감시. Rust Core에서 FW 상태를 수신하고 이상 시 복구 명령 전송.

```rust
pub struct FwWatchdogBridge {
    task_states: HashMap<FwTask, TaskHealth>,
    watchdog_config: WatchdogConfig,
    reset_history: VecDeque<ResetEvent>,
}

#[derive(Clone, Copy, Hash, Eq, PartialEq)]
pub enum FwTask {
    AdcSampler,      // 최고 우선순위, 1kHz
    BleComm,         // 고 우선순위
    NfcHandler,      // 중 우선순위
    EhdController,   // 중 우선순위
    PowerManager,    // 저 우선순위
    OtaUpdater,      // 최저 우선순위
}

pub struct TaskHealth {
    pub task: FwTask,
    pub status: TaskStatus,
    pub last_heartbeat: u64,
    pub deadline_miss_count: u32,
    pub stack_usage_percent: f64,
    pub cpu_usage_percent: f64,
}

pub enum TaskStatus {
    Running,
    Blocked,      // 리소스 대기 중 (정상)
    Starved,      // 하트비트 지연 > 2× 주기
    Crashed,      // 하트비트 완전 정지
    Restarted,    // 크래시 후 자동 재시작됨
}

pub struct WatchdogConfig {
    pub iwdg_timeout_ms: u32,     // 독립 워치독 (기본 4000ms)
    pub wwdg_window_ms: u32,      // 윈도우 워치독 (기본 50ms)
    pub heartbeat_interval_ms: u32, // 태스크 하트비트 주기 (기본 1000ms)
    pub max_restart_count: u8,    // 자동 재시작 최대 횟수 (기본 3)
}

impl FwWatchdogBridge {
    /// 펌웨어로부터 태스크 하트비트 수신 (BLE 경유)
    pub fn receive_heartbeat(&mut self, task: FwTask, metrics: TaskMetrics);

    /// 전체 FW 태스크 상태 평가
    pub fn evaluate(&self) -> FwHealthReport;

    /// 크래시된 태스크 복구 명령 전송
    pub fn request_task_restart(&mut self, task: FwTask) -> Result<(), FwError>;

    /// 전체 MCU 소프트 리셋 (최후 수단)
    pub fn request_mcu_reset(&mut self) -> Result<(), FwError>;

    /// L2 계층 점수
    pub fn layer_score(&self) -> f64;
}

pub struct FwHealthReport {
    pub overall_status: FwOverallStatus,
    pub tasks: Vec<TaskHealth>,
    pub uptime_seconds: u64,
    pub reset_count_24h: u32,
    pub score: f64,        // 0.0~1.0
}

pub enum FwOverallStatus {
    Nominal,          // 모든 태스크 정상
    Degraded,         // 비핵심 태스크 1개 이상 이상
    CriticalFailure,  // 핵심 태스크(ADC/BLE) 이상
    RecoveryMode,     // 리셋 후 복구 중
}
```

**테스트 기준 (6개)**:
1. `receive_heartbeat()` — 정상 하트비트 → `TaskStatus::Running`
2. `evaluate()` — 모든 태스크 정상 → `FwOverallStatus::Nominal`, `score > 0.9`
3. `evaluate()` — AdcSampler 하트비트 3초 지연 → `TaskStatus::Starved`, `FwOverallStatus::CriticalFailure`
4. `evaluate()` — OtaUpdater 크래시 → `FwOverallStatus::Degraded` (비핵심이므로)
5. `request_task_restart()` — 크래시 태스크 재시작 → `TaskStatus::Restarted`, `Ok(())`
6. `request_task_restart()` — 재시작 3회 초과 → `Err(FwError::MaxRestartExceeded)`

---

#### 6.15.3 데이터 무결성 검증기 (DataIntegrityChecker)

**역할**: Layer 3 Rust Core의 데이터 파이프라인 무결성 검증. HashChain(§6.7) 확장.

```rust
pub struct DataIntegrityChecker {
    hash_chain: HashChain,           // §6.7 기존
    pipeline_validators: Vec<Box<dyn PipelineValidator>>,
    anomaly_history: VecDeque<IntegrityAnomaly>,
}

pub trait PipelineValidator {
    fn name(&self) -> &str;
    fn validate(&self, input: &[u8], output: &[u8]) -> ValidationResult;
}

pub struct ValidationResult {
    pub valid: bool,
    pub stage: PipelineStage,
    pub checksum_match: bool,
    pub range_check: bool,       // 출력값이 물리적으로 가능한 범위인지
    pub temporal_check: bool,    // 시간 순서 정합성
    pub details: String,
}

pub enum PipelineStage {
    AdcCapture,           // ADC → RawData
    DifferentialCompute,  // RawData → DifferentialSignal
    DspProcessing,        // DifferentialSignal → ProcessedSignal
    FeatureExtraction,    // ProcessedSignal → Fingerprint
    AiInference,          // Fingerprint → InferenceResult
    Storage,              // InferenceResult → LocalStore
    Transmission,         // LocalStore → Cloud Sync
}

impl DataIntegrityChecker {
    /// 전체 파이프라인 단계별 체크섬 검증
    /// ADC → Differential → DSP → Feature → AI → Storage → Transmission
    pub fn verify_pipeline(&mut self, measurement_id: &str) -> PipelineIntegrityReport;

    /// 특정 단계의 입출력 정합성 확인
    pub fn verify_stage(&self, stage: PipelineStage, input: &[u8], output: &[u8]) -> ValidationResult;

    /// 물리적 범위 검증 (예: 혈당 0~600 mg/dL, pH 0~14)
    pub fn range_check(&self, analyte: &str, value: f64) -> bool;

    /// 시간 순서 이상 탐지 (타임스탐프 역전, 중복 등)
    pub fn temporal_check(&self, timestamps: &[u64]) -> bool;

    /// HashChain 연속성 검증 (§6.7 확장)
    pub fn verify_chain_continuity(&self, from_index: u32, to_index: u32) -> bool;

    /// L3 계층 점수
    pub fn layer_score(&self) -> f64;
}

pub struct PipelineIntegrityReport {
    pub measurement_id: String,
    pub stages: Vec<ValidationResult>,
    pub overall_valid: bool,
    pub tamper_detected: bool,    // 데이터 위변조 의심
    pub score: f64,               // 0.0~1.0
}
```

**테스트 기준 (7개)**:
1. `verify_pipeline()` — 정상 측정 데이터 → `overall_valid=true`, `score=1.0`
2. `verify_stage(DifferentialCompute)` — compute(100,50)=51.0 → `valid=true`
3. `verify_stage(AiInference)` — 모델 출력 confidence < 0 → `valid=false` (범위 위반)
4. `range_check("glucose", 500.0)` → `true` (정상 범위 내)
5. `range_check("glucose", -10.0)` → `false` (물리적 불가)
6. `temporal_check([100, 200, 150])` → `false` (타임스탐프 역전)
7. `verify_chain_continuity()` — 중간 해시 변조 → `false`, `tamper_detected=true`

---

#### 6.15.4 예방정비 엔진 (PredictiveMaintenanceEngine)

**역할**: HW 열화 추세 + 디지털 트윈 잔차 + 사용 패턴 → 고장 시점 예측 및 선제 조치

```rust
pub struct PredictiveMaintenanceEngine {
    component_trends: HashMap<HwComponent, Vec<TrendPoint>>,
    twin_trends: Vec<TwinTrendPoint>,
    usage_pattern: UsagePattern,
    maintenance_schedule: Vec<MaintenanceTask>,
}

pub struct TrendPoint {
    pub timestamp: u64,
    pub score: f64,
    pub degradation_rate: f64,
}

pub struct UsagePattern {
    pub daily_measurements: f64,      // 일 평균 측정 횟수
    pub avg_session_duration_sec: f64,
    pub cartridge_change_frequency: f64,  // 일 평균 교체 횟수
    pub environment_stress: f64,      // 환경 스트레스 지수 (온습도 변동성)
}

pub struct MaintenanceTask {
    pub component: HwComponent,
    pub task_type: MaintenanceType,
    pub priority: MaintenancePriority,
    pub due_date: u64,               // 예상 필요 시점
    pub confidence: f64,
    pub description: String,
    pub user_action: String,          // 사용자에게 안내할 조치
}

pub enum MaintenanceType {
    CalibrationRenewal,    // 보정 갱신
    CartridgeReplacement,  // 카트리지 교체
    BatteryService,        // 배터리 점검/교체
    FirmwareUpdate,        // 펌웨어 업데이트
    CleaningSuggestion,    // 센서 청소 권고
    ProfessionalService,   // 전문 서비스 센터 방문
}

pub enum MaintenancePriority {
    Routine,     // 정기 점검 (30일+ 여유)
    Advisory,    // 권고 사항 (7~30일 이내)
    Urgent,      // 긴급 (7일 이내)
    Immediate,   // 즉시 조치 필요
}

impl PredictiveMaintenanceEngine {
    /// 전체 컴포넌트 열화 추세 분석 → 정비 스케줄 생성
    /// 선형 회귀 + Weibull 생존 분석 기반
    /// Weibull 파라미터: 형상(β) — β<1 초기고장, β≈1 랜덤고장, β>1 마모고장
    ///                  척도(η) — 63.2% 누적고장 시점(측정 횟수 기준)
    /// 컴포넌트별 초기값: ADC β=2.5/η=5000, BLE β=1.2/η=8000, 센서 β=3.0/η=3000
    /// 실측 데이터 축적 시 MLE(최대우도추정)로 파라미터 갱신
    pub fn analyze_trends(&mut self, hw_history: &[HwCheckResult], twin_history: &[TwinUpdate]) -> Vec<MaintenanceTask>;

    /// 다음 정비 필요 시점 예측 (일수)
    pub fn predict_next_maintenance(&self) -> Option<MaintenanceTask>;

    /// 카트리지 잔여 수명 예측 (디지털 트윈 + 사용 패턴)
    pub fn predict_cartridge_life(&self, twin: &TwinUpdate, usage: &UsagePattern) -> CartridgeLifePrediction;

    /// 배터리 수명 예측 (사이클수 + 충방전 패턴)
    pub fn predict_battery_life(&self, cycles: u32, avg_discharge_rate: f64) -> BatteryLifePrediction;
}

pub struct CartridgeLifePrediction {
    pub remaining_measurements: u32,
    pub remaining_days: f64,
    pub replacement_urgency: MaintenancePriority,
    pub confidence: f64,
}

pub struct BatteryLifePrediction {
    pub remaining_cycles: u32,
    pub health_percent: f64,    // SoH (State of Health)
    pub replacement_month: Option<u32>,
    pub confidence: f64,
}
```

**테스트 기준 (6개)**:
1. `analyze_trends()` — 안정적 score(30일 0.95→0.93) → `MaintenancePriority::Routine`
2. `analyze_trends()` — 급락 score(7일 0.90→0.70) → `MaintenancePriority::Urgent`
3. `predict_cartridge_life()` — drift_score=0.8, 일 3회 측정 → `remaining_days < 7`
4. `predict_cartridge_life()` — drift_score=0.1, 일 1회 → `remaining_days > 30`
5. `predict_battery_life()` — 500 사이클, SoH 80% → `remaining_cycles > 100`
6. `predict_next_maintenance()` — 정비 항목 없음 → `None`

---

#### 6.15.5 오류 리포터 (ErrorReporter)

**역할**: 6계층 전체 오류/이상 수집, 근본원인 분석(RCA), 자동 보고서 생성

```rust
pub struct ErrorReporter {
    error_log: VecDeque<SystemError>,
    rca_engine: RootCauseAnalyzer,
    report_queue: Vec<DiagnosticReport>,
}

pub struct SystemError {
    pub id: String,                 // UUID
    pub timestamp: u64,
    pub layer: SystemLayer,
    pub component: String,
    pub severity: AlertSeverity,
    pub error_code: u32,
    pub message: String,
    pub context: HashMap<String, String>,  // 발생 시점 환경 정보
    pub auto_resolved: bool,
    pub resolution: Option<String>,
}

pub enum SystemLayer {
    Hardware,       // L1
    Firmware,       // L2
    RustCore,       // L3
    FlutterApp,     // L4
    GoBackend,      // L5
    AiInfra,        // L6
}

pub struct RootCauseAnalyzer;
impl RootCauseAnalyzer {
    /// 연관 오류 클러스터링 → 근본원인 추론
    /// 시간 근접성(5초 이내) + 계층 인접성 + 컴포넌트 의존관계 기반
    pub fn analyze(&self, errors: &[SystemError]) -> RootCauseReport;
}

pub struct RootCauseReport {
    pub probable_cause: String,
    pub confidence: f64,
    pub affected_layers: Vec<SystemLayer>,
    pub error_chain: Vec<String>,     // 오류 전파 경로
    pub suggested_fix: String,
    pub user_action: String,          // 사용자 안내 문구
    pub requires_service: bool,       // 서비스센터 방문 필요 여부
}

pub struct DiagnosticReport {
    pub id: String,
    pub generated_at: u64,
    pub system_health_score: f64,
    pub layer_scores: HashMap<SystemLayer, f64>,
    pub active_issues: Vec<SystemError>,
    pub root_causes: Vec<RootCauseReport>,
    pub maintenance_tasks: Vec<MaintenanceTask>,
    pub recommendations: Vec<String>,
}

impl ErrorReporter {
    /// 오류 수집 (각 계층에서 호출)
    pub fn report_error(&mut self, error: SystemError);

    /// 근본원인 분석 실행
    pub fn run_rca(&mut self) -> Vec<RootCauseReport>;

    /// 종합 진단 보고서 생성
    pub fn generate_report(&self, health: &SystemHealthScore) -> DiagnosticReport;

    /// 사용자용 간소화 보고서 (기술 용어 제거)
    pub fn generate_user_report(&self, report: &DiagnosticReport) -> UserDiagnosticSummary;
}

pub struct UserDiagnosticSummary {
    pub health_emoji: char,           // 😊/😐/😟/🚨
    pub health_text: String,          // "리더기 상태 양호"
    pub issues: Vec<UserIssue>,
    pub actions: Vec<UserAction>,
}

pub struct UserIssue {
    pub icon: String,
    pub title: String,                // "배터리 수명 감소"
    pub description: String,          // "30일 이내 충전 성능 저하 예상"
}

pub struct UserAction {
    pub priority: MaintenancePriority,
    pub action: String,               // "가까운 서비스센터에서 배터리 점검"
    pub deep_link: String,            // "manpasik://settings/service-center"
}
```

**테스트 기준 (6개)**:
1. `report_error()` — HW 오류 등록 → error_log에 추가, 타임스탐프 자동 기록
2. `run_rca()` — BLE 끊김 + ADC 읽기 실패 (3초 이내) → `probable_cause: "BLE 통신 장애로 인한 데이터 수신 실패"`
3. `run_rca()` — 단일 오류 → 근본원인 = 자기 자신
4. `generate_report()` — score 0.85 → recommendations에 "정기 점검 권고" 포함
5. `generate_user_report()` — Critical 이슈 → `health_emoji: '🚨'`
6. `generate_user_report()` — 이슈 없음 → `health_emoji: '😊'`, actions 비어있음

---

#### 6.15.6 SystemHealthOrchestrator — 통합 오케스트레이터

**역할**: 6계층 점수 가중 합산 → 종합 건강 점수 + 의사결정

```rust
pub struct SystemHealthOrchestrator {
    hw_monitor: HwHealthMonitor,          // §6.15.1
    fw_watchdog: FwWatchdogBridge,        // §6.15.2
    data_checker: DataIntegrityChecker,   // §6.15.3
    predictive: PredictiveMaintenanceEngine, // §6.15.4
    error_reporter: ErrorReporter,        // §6.15.5
    // L4(§7.7), L5(§8.6), L6(§9.5)는 FFI/gRPC로 수신
    layer_weights: [f64; 6],  // [0.30, 0.15, 0.25, 0.10, 0.10, 0.10] 합=1.0
    // ⚠️ IEEE 754 부동소수점 합산 오차 방지: overall = Σ(score×weight) / Σ(weight)
    //    등급 경계(0.90, 0.75 등) 비교 시 ε=1e-9 허용 오차 적용
}

pub struct SystemHealthScore {
    pub overall: f64,                     // 0.0~1.0 정규화 가중 합산
    pub grade: HealthGrade,
    pub layers: [(SystemLayer, f64); 6],  // 각 계층 점수
    pub timestamp: u64,
    pub trend: HealthTrend,               // 전일 대비
    pub measurement_ready: bool,          // 측정 진행 가능 여부
}

pub enum HealthGrade {
    Excellent,   // ≥ 0.90
    Good,        // ≥ 0.75
    Fair,        // ≥ 0.60
    Poor,        // ≥ 0.40
    Critical,    // < 0.40
}

pub enum HealthTrend {
    Improving,    // 전일 대비 +5% 이상
    Stable,       // ±5% 이내
    Declining,    // 전일 대비 -5% 이상
}

impl SystemHealthOrchestrator {
    /// 전체 시스템 건강 점수 계산
    /// score = Σ(layer_score × weight) / Σ(weight)
    pub fn compute_health(&mut self) -> SystemHealthScore;

    /// 측정 전 자가검증 (pre-flight check)
    /// HW pre_check + FW evaluate + 데이터 무결성 → 측정 가능 판단
    pub fn pre_flight_check(&mut self) -> PreFlightResult;

    /// 측정 후 결과 신뢰도 종합 평가
    pub fn post_flight_evaluate(&mut self, measurement_id: &str) -> PostFlightResult;

    /// 주기적 백그라운드 건강 점검 (15분 간격)
    pub fn background_check(&mut self) -> SystemHealthScore;

    /// BLE GATT 0xFF08 직렬화 (새 Characteristic)
    pub fn serialize_for_ble(&self, score: &SystemHealthScore) -> [u8; 20];
}

pub struct PreFlightResult {
    pub ready: bool,
    pub hw_pass: bool,
    pub fw_pass: bool,
    pub data_pass: bool,
    pub blockers: Vec<String>,
    pub warnings: Vec<String>,
}

pub struct PostFlightResult {
    pub measurement_id: String,
    pub reliability: ReliabilityScore,
    pub data_integrity: bool,
    pub recommendation: MeasurementRecommendation,
}
```

**테스트 기준 (8개)**:
1. `compute_health()` — 모든 계층 1.0 → `overall=1.0`, `grade=Excellent`
2. `compute_health()` — HW=0.3, 나머지 1.0 → `overall < 0.75` (HW 가중치 30%가 반영)
3. `pre_flight_check()` — 전체 정상 → `ready=true`, `blockers` 비어있음
4. `pre_flight_check()` — 배터리 5% → `ready=false`, `blockers=["배터리 부족 (5%)"]`
5. `post_flight_evaluate()` — 정상 측정 → `reliability.score > 0.9`, `recommendation=Reliable`
6. `background_check()` — 15분 간격 호출 → 이전 대비 ±5% 이내 → `trend=Stable`
7. `serialize_for_ble()` — SystemHealthScore → 20바이트 직렬화 + 역직렬화 일치
8. `compute_health()` — L1 Failed → `measurement_ready=false`

---

### 6.16 자가치유 시스템 (Self-Healing Engine) — v2.4 신규

**파일**: `healing/orchestrator.rs`, `hw_recovery.rs`, `fw_repair.rs`, `pipeline_reprocess.rs`, `session_recovery.rs`
**의존**: §6.15 SelfDiagnostics, §6.4 DigitalTwinEngine, §6.1 DifferentialEngine, §6.7 HashChain
**특허**: Family A cl.1 (자기최적화), Family C (모듈형 플랫폼)

#### 자가점검 vs 자가치유 아키텍처

**역할 경계 원칙**: §6.15 SelfDiagnostics는 **읽기 전용** 관찰자로 HW/FW/데이터 상태를 감시·점수화하며 시스템을 변경하지 않는다. §6.16 SelfHealing은 SelfDiagnostics의 이상 이벤트를 수신한 후에만 **쓰기** 동작(재보정, 태스크 재시작, 파이프라인 재처리, 세션 복원)을 수행한다. 두 모듈은 이벤트 버스(`DiagnosticEvent`)를 통해 단방향 연결되며, SelfHealing이 SelfDiagnostics를 우회하여 독자적으로 복구를 개시하는 경로는 없다. FwWatchdogBridge(§6.15.2)의 `request_task_restart`는 예외적으로 FW 태스크 크래시 시 즉시 재시작 명령을 전송하는데, 이는 SelfHealing의 FwSelfRepair(§6.16.2)와 달리 단순 RTOS 레벨 재시작만 수행하며 힙 정리·OTA 롤백 등 고급 복구는 FwSelfRepair가 담당한다.

```
┌─────────────────────────────────────────────────────────────────────┐
│  §6.15 SelfDiagnostics (점검)          §6.16 SelfHealing (치유)     │
│  ┌────────────────────────────┐       ┌──────────────────────────┐ │
│  │ 이상 감지 (Detection)      │──────▶│ 자동 복구 (Recovery)     │ │
│  │                            │       │                          │ │
│  │ HwHealthMonitor            │       │ HwAutoRecovery           │ │
│  │ FwWatchdogBridge           │       │ FwSelfRepair             │ │
│  │ DataIntegrityChecker       │       │ DataPipelineReprocessor  │ │
│  │ PredictiveMaintenance      │       │ MeasurementSessionRecov. │ │
│  │ ErrorReporter              │       │ SelfHealingOrchestrator  │ │
│  └────────────────────────────┘       └──────────────────────────┘ │
│                                                                     │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │ 치유 정책: Detect → Classify → Attempt → Verify → Log         │ │
│  │                                                                │ │
│  │ 1. 이상 감지 (§6.15에서 수신)                                  │ │
│  │ 2. 자동복구 가능 여부 분류 (auto_recoverable?)                 │ │
│  │ 3. 복구 시도 (최대 retry_limit회)                              │ │
│  │ 4. 복구 성공 검증 (re-diagnose)                                │ │
│  │ 5. 실패 시 에스컬레이션 (사용자 안내 또는 서비스 티켓)        │ │
│  └────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
```

#### 6.16.1 하드웨어 자동 복구 (HwAutoRecovery)

**역할**: HwHealthMonitor(§6.15.1)가 감지한 HW 이상에 대해 소프트웨어로 복구 가능한 항목을 자동 처리

```rust
pub struct HwAutoRecovery {
    recovery_policies: HashMap<HwComponent, RecoveryPolicy>,
    recovery_log: VecDeque<RecoveryAttempt>,
    ble_reconnector: BleReconnector,
    adc_recalibrator: AdcRecalibrator,
}

pub struct RecoveryPolicy {
    pub component: HwComponent,
    pub max_retries: u8,
    pub backoff_ms: u32,       // 재시도 간격 (지수 백오프)
    pub auto_enabled: bool,    // 사용자 승인 없이 자동 실행 가능 여부
    pub strategies: Vec<RecoveryStrategy>,
}

pub enum RecoveryStrategy {
    // ── BLE 복구 ──
    BleReconnect {
        timeout_ms: u32,        // 기본 3000ms
        channel_hop: bool,      // 채널 변경 시도
    },
    BleSessionRestore {
        session_data: Vec<u8>,  // 세션 상태 백업에서 복원
    },

    // ── ADC 복구 ──
    AdcResample {
        discard_count: u8,      // 이상 샘플 폐기 후 재샘플링
        averaging_window: u8,   // 추가 평균화 윈도우
    },
    AdcGainRecalibrate,         // 자동 게인 재조정

    // ── NFC 복구 ──
    NfcRetryAuth {
        retry_count: u8,        // SUN/CMAC 재인증 시도
    },
    NfcFallbackQr,              // NFC 실패 시 QR 폴백 전환

    // ── AFE 복구 ──
    AfeBlockSwitch(u8),         // 고장 AFE → 대체 블록 자동 전환
    AfeImpedanceReset,          // 임피던스 초기화

    // ── 전원 복구 ──
    PowerCycleComponent(HwComponent),  // 개별 컴포넌트 전원 리셋

    // ── 환경 보상 ──
    TemperatureCompensation(f64),   // TEC 자동 보정
    HumidityCompensation(f64),      // 습도 보정값 재계산

    // ── 보정 복구 ──
    CalibrationCloudSync,       // 만료 보정 → 클라우드에서 최신 보정 자동 다운로드
    CalibrationFallbackDefault, // 클라우드 불가 시 범용 보정 폴백
}

pub struct RecoveryAttempt {
    pub component: HwComponent,
    pub strategy: String,
    pub timestamp: u64,
    pub success: bool,
    pub duration_ms: u32,
    pub before_score: f64,
    pub after_score: f64,
}

impl HwAutoRecovery {
    /// HW 이상 감지 시 자동 복구 시도
    /// auto_recoverable=true인 경우에만 사용자 승인 없이 실행
    pub fn attempt_recovery(
        &mut self,
        alert: &HwAlert,
    ) -> RecoveryResult;

    /// BLE 자동 재연결 + 세션 상태 복원
    /// 3초 타임아웃 → 채널 호핑 → 실패 시 사용자 안내
    pub fn recover_ble(&mut self) -> RecoveryResult;

    /// ADC 노이즈 이상 시 자동 리샘플링
    /// 이상 샘플 폐기 + 추가 평균화로 노이즈 저감
    pub fn recover_adc(&mut self, bad_samples: &[i32]) -> RecoveryResult;

    /// 보정 만료 시 클라우드 보정 자동 다운로드
    /// Cloud → QR 폴백 → 범용 보정 폴백 체인
    pub fn recover_calibration(&mut self, manifest: &CartridgeManifest) -> RecoveryResult;

    /// AFE 블록 고장 시 대체 블록 자동 전환
    /// 9블록 중 동일 모달리티의 대체 블록이 있으면 전환
    pub fn switch_afe_block(&mut self, failed_block: u8) -> RecoveryResult;
}

pub enum RecoveryResult {
    Success {
        strategy: String,
        duration_ms: u32,
        new_score: f64,
    },
    PartialSuccess {
        strategy: String,
        remaining_issues: Vec<String>,
        new_score: f64,
    },
    Failed {
        strategy: String,
        reason: String,
        escalation: EscalationType,
    },
}

pub enum EscalationType {
    UserNotification(String),   // 사용자에게 안내
    ServiceTicket,              // 서비스 티켓 자동 생성
    MeasurementBlock,           // 측정 차단
}
```

**테스트 기준 (8개)**:
1. `recover_ble()` — RSSI -85dBm(약함) → BleReconnect 시도 → 재연결 성공 → `RecoveryResult::Success`
2. `recover_ble()` — 완전 단절 + 3초 초과 → `RecoveryResult::Failed { escalation: UserNotification }`
3. `recover_adc()` — 포화 샘플 3개 → 폐기 + 리샘플링 → 정상 데이터 복원
4. `recover_adc()` — 연속 10개 포화 → `RecoveryResult::Failed { escalation: MeasurementBlock }`
5. `recover_calibration()` — 만료 보정 → Cloud 다운로드 성공 → `Success`
6. `recover_calibration()` — Cloud 불가 + QR 불가 → 범용 폴백 → `PartialSuccess { remaining: ["정확도 저하 가능"] }`
7. `switch_afe_block()` — Block 0 고장, Block 6(동일 모달리티) 가용 → 전환 성공
8. `attempt_recovery()` — `auto_recoverable=false` → 즉시 사용자 에스컬레이션 (자동 시도 안 함)

---

#### 6.16.2 펌웨어 자가수리 (FwSelfRepair)

**역할**: FwWatchdogBridge(§6.15.2)가 감지한 태스크 이상을 단계적으로 복구. 기존 `request_task_restart()`를 확장.

```rust
pub struct FwSelfRepair {
    repair_history: VecDeque<RepairAction>,
    escalation_level: EscalationLevel,
}

pub enum EscalationLevel {
    Level0_TaskRestart,    // 개별 태스크 재시작
    Level1_StackReset,     // 태스크 스택 초기화 + 재시작
    Level2_PeripheralReset,// 관련 주변장치 리셋 (SPI/I2C/UART)
    Level3_SoftReset,      // MCU 소프트 리셋
    Level4_FactoryFallback,// 팩토리 펌웨어 폴백 (A/B 파티션)
}

pub struct RepairAction {
    pub task: FwTask,
    pub level: EscalationLevel,
    pub timestamp: u64,
    pub success: bool,
    pub downtime_ms: u32,   // 복구 소요 시간
}

impl FwSelfRepair {
    /// 단계적 에스컬레이션 복구
    /// Level0 실패 → Level1 → Level2 → Level3 → Level4
    /// 각 단계에서 성공하면 즉시 중단
    pub fn escalated_repair(&mut self, task: FwTask) -> RepairResult;

    /// OTA 실패 시 이전 파티션 자동 롤백
    pub fn ota_rollback(&mut self) -> RepairResult;

    /// 메모리 누수 감지 시 자동 GC (FreeRTOS heap 정리)
    pub fn heap_cleanup(&mut self, task: FwTask) -> RepairResult;

    /// 주변장치 교착 해소 (SPI/I2C 버스 리셋)
    pub fn peripheral_bus_reset(&mut self, bus: PeripheralBus) -> RepairResult;
}

pub enum PeripheralBus {
    Spi1,    // ADS1256 ADC
    I2c1,    // PN7150 NFC
    Uart2,   // nRF52832 BLE
}

pub enum RepairResult {
    Repaired { level: EscalationLevel, downtime_ms: u32 },
    Degraded { level: EscalationLevel, remaining_issues: Vec<String> },
    Unrecoverable { final_level: EscalationLevel, recommendation: String },
}
```

**테스트 기준 (6개)**:
1. `escalated_repair(OtaUpdater)` — Level0 재시작 성공 → `Repaired { level: Level0, downtime_ms < 500 }`
2. `escalated_repair(AdcSampler)` — Level0~1 실패, Level2 성공 → `Repaired { level: Level2 }`
3. `escalated_repair()` — Level0~4 전부 실패 → `Unrecoverable { recommendation: "서비스 센터 방문" }`
4. `ota_rollback()` — 현재 파티션 부트 실패 → 이전 파티션 성공 → `Repaired`
5. `heap_cleanup()` — FreeRTOS heap 80%→45% → `Repaired`
6. `peripheral_bus_reset(Spi1)` — ADS1256 응답 없음 → SPI 리셋 → 응답 복원

---

#### 6.16.3 데이터 파이프라인 재처리기 (DataPipelineReprocessor)

**역할**: DataIntegrityChecker(§6.15.3)가 발견한 비위변조 오류를 자동 재처리. 위변조(tamper)는 재처리 불가 — 즉시 폐기.

```rust
pub struct DataPipelineReprocessor {
    retry_policies: HashMap<PipelineStage, RetryPolicy>,
    reprocess_log: VecDeque<ReprocessAttempt>,
}

pub struct RetryPolicy {
    pub max_retries: u8,
    pub recompute_from: PipelineStage,  // 어디서부터 재처리
    pub timeout_ms: u32,
}

impl DataPipelineReprocessor {
    /// 체크섬 불일치 시 해당 스테이지부터 재처리
    /// tamper_detected=true → 즉시 Discard (재처리 불가)
    /// tamper_detected=false → 해당 스테이지 재연산
    pub fn reprocess(
        &mut self,
        report: &PipelineIntegrityReport,
        raw_data: &[u8],
    ) -> ReprocessResult;

    /// 차동측정 결과 재연산
    /// DifferentialEngine::compute() 재실행 + 결과 비교
    pub fn recompute_differential(
        &self,
        s_det: f64,
        s_ref: f64,
        alpha: f64,
    ) -> f64;

    /// AI 추론 재실행 (다른 모델 경로로 교차 검증)
    /// 1차: 동일 모델 재실행
    /// 2차: 백업 모델(XGBoost ↔ TFLite 교차) 실행
    /// 두 결과 불일치 시 confidence 낮은 쪽 채택 + 주석 기록
    pub fn cross_validate_inference(
        &self,
        fingerprint: &[f64],
    ) -> CrossValidationResult;

    /// CRDT 동기화 충돌 자동 해소
    /// 타임스탬프 최신 우선 + 측정 데이터는 항상 보존 (삭제 금지)
    pub fn resolve_crdt_conflict(
        &mut self,
        local: &CrdtState,
        remote: &CrdtState,
    ) -> MergedState;
}

pub struct CrossValidationResult {
    pub primary_result: InferenceResult,
    pub backup_result: Option<InferenceResult>,
    pub agreement: bool,          // 두 모델 결과 일치 여부
    pub chosen: InferenceResult,  // 최종 채택 결과
    pub annotation: Option<String>,  // 불일치 시 주석
}

pub enum ReprocessResult {
    Recovered {
        stage: PipelineStage,
        retries: u8,
        new_checksum_valid: bool,
    },
    Discarded {
        reason: String,  // "tamper_detected" 또는 "max_retries_exceeded"
    },
}
```

**테스트 기준 (7개)**:
1. `reprocess()` — DifferentialCompute 체크섬 불일치, tamper=false → 재연산 → `Recovered`
2. `reprocess()` — tamper_detected=true → `Discarded { reason: "tamper_detected" }`
3. `reprocess()` — 3회 재시도 후에도 불일치 → `Discarded { reason: "max_retries_exceeded" }`
4. `cross_validate_inference()` — TFLite·XGBoost 일치 → `agreement=true`
5. `cross_validate_inference()` — TFLite glucose=115, XGBoost glucose=118 → `agreement=true` (허용 오차 내)
6. `cross_validate_inference()` — TFLite glucose=115, XGBoost glucose=180 → `agreement=false`, confidence 낮은 쪽 채택, annotation 기록
7. `resolve_crdt_conflict()` — local(12:00) vs remote(12:05) 동일 필드 → remote 채택 (최신)

---

#### 6.16.4 측정 세션 복구 (MeasurementSessionRecovery)

**역할**: 측정 중 비정상 중단(BLE 끊김, 앱 크래시, 배터리 방전) 시 세션 상태를 자동 저장/복원

```rust
pub struct MeasurementSessionRecovery {
    checkpoint_store: CheckpointStore,
    recovery_timeout_ms: u32,  // 복원 가능 최대 시간 (기본 300,000ms = 5분)
}

/// 직렬화 포맷: bincode v2 (Little-Endian, VarInt 길이 접두사)
/// 저장 위치: SQLite `checkpoints` 테이블 BLOB 컬럼 + 메모리 LRU 캐시(최대 3개)
/// 체크섬: SHA-256(session_id ‖ bincode_payload)
/// 최대 크기: ~64KB (twin_state 압축 후 기준, LZ4 적용)
pub struct SessionCheckpoint {
    pub session_id: String,
    pub measurement_type: String,
    pub cartridge_manifest: CartridgeManifest,
    pub progress_percent: f64,
    pub collected_samples: u32,
    pub total_expected_samples: u32,
    pub partial_results: Option<Vec<u8>>,  // 중간 결과 — bincode 직렬화된 PartialMeasurement
    pub alpha_state: f64,
    pub twin_state: Vec<u8>,              // DigitalTwin 스냅샷 — LZ4 압축 bincode
    pub timestamp: u64,
    pub checksum: [u8; 32],               // SHA-256 무결성 해시
}

impl MeasurementSessionRecovery {
    /// 측정 진행 중 100ms마다 체크포인트 자동 저장
    /// 로컬 스토리지(SQLite) + 메모리 캐시 이중 저장
    pub fn save_checkpoint(&mut self, checkpoint: &SessionCheckpoint) -> Result<(), StorageError>;

    /// 앱 재시작 시 미완료 세션 탐색
    /// 5분 이내 중단 + 카트리지 동일 + 체크섬 유효 → 복원 가능
    pub fn find_recoverable_session(&self) -> Option<SessionCheckpoint>;

    /// 세션 복원 — 중단 지점부터 측정 재개
    /// BLE 재연결 → 카트리지 재인증 → 수집 재개
    pub fn restore_session(
        &mut self,
        checkpoint: &SessionCheckpoint,
    ) -> SessionRestoreResult;

    /// 종합검진(CheckupSession) 부분 복구
    /// 다중 카트리지 순차 측정 중 실패 시, 완료된 카트리지 결과는 보존
    pub fn restore_checkup_session(
        &mut self,
        checkup_checkpoint: &CheckupCheckpoint,
    ) -> CheckupRestoreResult;

    /// 체크포인트 만료 정리 (5분 초과 또는 성공 완료)
    pub fn cleanup_expired(&mut self) -> u32;
}

pub enum SessionRestoreResult {
    Resumed {
        from_progress: f64,   // 재개 시점 (%)
        estimated_remaining_sec: u32,
    },
    PartialResult {
        collected_percent: f64,
        result: Option<InferenceResult>,  // 수집된 데이터로 부분 추론
        caveat: String,  // "75% 데이터 기반 추정치"
    },
    Expired {
        reason: String,  // "5분 초과" 또는 "카트리지 변경됨"
    },
}

pub struct CheckupRestoreResult {
    pub completed_tests: Vec<String>,  // 완료된 카트리지 결과 보존
    pub failed_test: String,           // 실패한 카트리지
    pub remaining_tests: Vec<String>,  // 아직 안 한 카트리지
    pub can_continue: bool,            // 실패 카트리지 건너뛰고 계속 가능 여부
}
```

**테스트 기준 (7개)**:
1. `save_checkpoint()` — 측정 50% 진행 → 체크포인트 저장 성공
2. `find_recoverable_session()` — 2분 전 중단 + 동일 카트리지 → `Some(checkpoint)`
3. `find_recoverable_session()` — 10분 전 중단 → `None` (5분 초과)
4. `restore_session()` — BLE 재연결 + 카트리지 일치 → `Resumed { from_progress: 50% }`
5. `restore_session()` — 카트리지 변경됨 → `Expired`
6. `restore_checkup_session()` — 3/5 카트리지 완료 후 중단 → `completed_tests: 3개`, `can_continue: true`
7. `cleanup_expired()` — 만료 체크포인트 5개 → 5개 삭제, 반환값 5

---

#### 6.16.5 SelfHealingOrchestrator — 자가치유 통합 오케스트레이터

**역할**: §6.15 SelfDiagnostics의 이상 감지 → §6.16 자가치유 자동 연결. "감지 → 분류 → 시도 → 검증 → 기록" 5단계 자가치유 루프.

```rust
pub struct SelfHealingOrchestrator {
    hw_recovery: HwAutoRecovery,             // §6.16.1
    fw_repair: FwSelfRepair,                 // §6.16.2
    pipeline_reprocessor: DataPipelineReprocessor, // §6.16.3
    session_recovery: MeasurementSessionRecovery,  // §6.16.4
    diagnostics: SystemHealthOrchestrator,   // §6.15.6 (점검과 연동)
    healing_log: VecDeque<HealingEvent>,
    healing_stats: HealingStatistics,
}

pub struct HealingEvent {
    pub id: String,
    pub timestamp: u64,
    pub trigger: HealingTrigger,
    pub layer: SystemLayer,
    pub strategy_used: String,
    pub result: HealingOutcome,
    pub duration_ms: u32,
    pub score_before: f64,
    pub score_after: f64,
}

pub enum HealingTrigger {
    AutomaticFromDiagnostics,  // §6.15 자동 감지 → 자동 치유
    UserInitiated,             // H-036 보고서에서 [자동 수정] 버튼
    ScheduledMaintenance,      // §6.15.4 예방정비 스케줄
    SessionInterruption,       // 측정 세션 비정상 중단
}

pub enum HealingOutcome {
    FullyHealed,       // 완전 복구 — 점수 정상 복원
    PartiallyHealed,   // 부분 복구 — 측정 가능하나 제약
    HealingFailed,     // 복구 실패 — 에스컬레이션
    NotNeeded,         // 진단 결과 치유 불필요
}

pub struct HealingStatistics {
    pub total_attempts_24h: u32,
    pub success_rate: f64,          // 0.0~1.0
    pub avg_recovery_time_ms: u32,
    pub most_common_issue: String,
    pub auto_heal_count: u32,       // 사용자 무개입 자동 치유 횟수
}

impl SelfHealingOrchestrator {
    /// 자가검증 이상 → 자가치유 자동 연결 (메인 루프)
    /// 1. §6.15에서 이상 수신
    /// 2. auto_recoverable 분류
    /// 3. 해당 계층 recovery 모듈 호출
    /// 4. 치유 후 re-diagnose (§6.15 재실행)
    /// 5. 성공/실패 기록
    pub fn heal(&mut self, error: &SystemError) -> HealingEvent;

    /// 측정 중 실시간 자가치유
    /// monitor_during()에서 Warning → 즉시 recover 시도
    /// 사용자는 모르는 사이에 자동 복구 (seamless)
    pub fn heal_during_measurement(&mut self, alert: &HwAlert) -> HealingOutcome;

    /// 앱 시작 시 자동 복구 점검
    /// 미완료 세션 복원 + 펌웨어 상태 확인 + 보정 유효성
    pub fn startup_recovery(&mut self) -> StartupRecoveryReport;

    /// 24시간 자가치유 통계
    pub fn get_statistics(&self) -> HealingStatistics;

    /// 치유 이벤트 로그 조회 (진단 보고서 포함용)
    pub fn get_healing_log(&self, hours: u32) -> Vec<HealingEvent>;
}

pub struct StartupRecoveryReport {
    pub session_restored: bool,
    pub fw_repaired: bool,
    pub calibration_updated: bool,
    pub issues_auto_resolved: u32,
    pub remaining_issues: Vec<String>,
}
```

**테스트 기준 (8개)**:
1. `heal()` — BLE Warning, auto_recoverable=true → HwAutoRecovery 호출 → `FullyHealed`
2. `heal()` — ADC Fatal, auto_recoverable=false → 즉시 `HealingFailed` + UserNotification
3. `heal_during_measurement()` — BLE 패킷 손실 → 자동 재연결 → 사용자 무개입 복구 → `FullyHealed`
4. `heal_during_measurement()` — ADC 포화 연속 → 리샘플링 실패 → `HealingFailed` → 측정 중단 안내
5. `startup_recovery()` — 미완료 세션 발견 → `session_restored: true`
6. `startup_recovery()` — 보정 만료 → 클라우드 다운로드 → `calibration_updated: true`
7. `get_statistics()` — 24시간 10회 시도, 8회 성공 → `success_rate: 0.8`
8. `heal()` → re-diagnose 후 score_after > score_before 검증

---

### 6.7 보안 모듈

**파일**: `security/hash_chain.rs`, `security/rolling_hash.rs`, `security/aes_gcm.rs`

#### 인터페이스 계약

```rust
pub struct HashChain {
    previous_hash: [u8; 32],  // SHA-256
}

impl HashChain {
    /// 측정 데이터에 해시 체인 연결 (무결성 보장)
    pub fn append(&mut self, data: &[u8]) -> [u8; 32];

    /// 체인 검증
    pub fn verify(chain: &[([u8; 32], &[u8])]) -> bool;
}

pub struct RollingHash;
impl RollingHash {
    /// 실시간 변조 탐지
    pub fn update(&mut self, new_data: &[u8]) -> [u8; 32];
}

/// AES-256-GCM 암호화/복호화
pub fn encrypt_aes_gcm(
    plaintext: &[u8],
    key: &[u8; 32],
    nonce: &[u8; 12],
) -> Result<Vec<u8>, CryptoError>;

pub fn decrypt_aes_gcm(
    ciphertext: &[u8],
    key: &[u8; 32],
    nonce: &[u8; 12],
) -> Result<Vec<u8>, CryptoError>;
```

**테스트 기준**:
1. 해시 체인 3개 연결 → 전체 검증 성공
2. 중간 데이터 변조 → 검증 실패
3. AES-GCM 왕복 → 평문 동일
4. AES-GCM 잘못된 키 → 복호화 실패

---

### 6.8 종합검진 세션 관리 (CheckupSession)

**파일**: `checkup/session.rs`
**역할**: 다중 카트리지 순차 측정을 하나의 검진 세션으로 묶고, 진행 상태 관리
**배경**: 단건 측정 중심에서 다중 검진 패키지로 확장 (혈액 3종 + 타액 1종 등)

#### 검진 패키지 정의

| 패키지 | 시료 | 카트리지 수 | 측정 항목 | 소요시간 |
|---|---|---|---|---|
| Quick | 혈액 1종 | 1 | 혈당, 전해질(Na/K/Cl), CRP | ~5분 |
| Standard | 혈액 2종 + 타액 1종 | 3 | Quick + HbA1c, 지질(TC/HDL/LDL/TG), 코르티솔, pH | ~15분 |
| Premium | 혈액 3종 + 소변 1종 + 타액 1종 | 5 | Standard + 간기능(ALT/AST/GGT), 신장(Creatinine/BUN), TSH, 요단백, 요당 | ~25분 |
| Custom | 사용자 선택 | 1~8 | 마켓플레이스 카트리지 자유 조합 | 가변 |

#### 인터페이스 계약

```rust
#[derive(Debug, Clone)]
pub enum CheckupPackage {
    Quick,
    Standard,
    Premium,
    Custom(Vec<u16>),  // cartridge_type 목록
}

#[derive(Debug, Clone)]
pub enum CheckupSessionState {
    Preparing,
    InProgress {
        current_step: u8,
        total_steps: u8,
        completed_results: Vec<MeasurementResult>,
    },
    Analyzing,
    Completed(CheckupReport),
    PartiallyCompleted {
        completed_results: Vec<MeasurementResult>,
        remaining_types: Vec<u16>,
    },
}

pub struct CheckupSession {
    pub session_id: Uuid,
    pub package: CheckupPackage,
    pub state: CheckupSessionState,
    pub started_at: u64,

    /// 패키지에 필요한 카트리지 유형 목록 반환 (측정 최적 순서)
    pub fn required_cartridge_types(&self) -> Vec<CartridgeTypeInfo>;

    /// 현재 단계의 카트리지를 NFC 인식 결과와 매칭
    pub fn validate_cartridge(&self, manifest: &CartridgeManifest) -> Result<StepMatch, CheckupError>;

    /// 개별 측정 결과를 세션에 추가
    pub fn add_measurement_result(&mut self, result: MeasurementResult) -> Result<CheckupProgress, CheckupError>;

    /// 세션 중간 저장 (이탈 대비)
    pub fn save_partial(&self) -> Result<PartialCheckup, CheckupError>;

    /// 부분 완료 세션 복원
    pub fn resume(partial: PartialCheckup) -> Result<Self, CheckupError>;
}
```

#### 테스트 기준 (8개)

| # | 테스트 | 기대 결과 |
|---|---|---|
| 1 | Quick 패키지 → required_cartridge_types() | [혈당/전해질/CRP 콤보] 1종 반환 |
| 2 | Standard 패키지 → required_cartridge_types() | 3종 반환, 혈액 먼저 정렬 |
| 3 | validate_cartridge() — 올바른 카트리지 | Ok(StepMatch) |
| 4 | validate_cartridge() — 틀린 카트리지 | Err(WrongCartridgeType) |
| 5 | validate_cartridge() — 만료 카트리지 | Err(CartridgeExpired) |
| 6 | add_measurement_result() 3회 호출 | state가 Analyzing으로 전환 |
| 7 | save_partial() → resume() | 이전 결과 유지 + 남은 단계 정확 |
| 8 | Custom([type1, type2, type3]) → 3종 반환 | 3개 요소 |

---

### 6.9 질병 리스크 엔진 (DiseaseRiskEngine)

**파일**: `checkup/disease_risk.rs`
**역할**: 검진 결과로부터 8대 질병 영역 리스크 점수 산출
**특허**: 결과 해석 및 예측 알고리즘 기반

#### 8대 질병 리스크 영역

```rust
#[derive(Debug, Clone, Copy)]
pub enum DiseaseRiskDomain {
    MetabolicSyndrome,    // 대사증후군
    DiabetesGlycemic,     // 당뇨/혈당관리
    Cardiovascular,       // 심혈관
    Liver,                // 간 건강
    Kidney,               // 신장 건강
    Thyroid,              // 갑상선
    StressHormone,        // 스트레스/호르몬
    InflammationImmune,   // 염증/면역
}

pub struct DiseaseRiskScore {
    pub domain: DiseaseRiskDomain,
    pub score: f64,             // 0.0~100.0
    pub level: RiskLevel,       // Normal/Caution/Warning/Danger
    pub contributing_factors: Vec<ContributingFactor>,
    pub trend: Option<RiskTrend>, // 이전 검진 대비 개선/유지/악화
}

pub enum RiskLevel {
    Normal,    // 0~25
    Caution,   // 26~50
    Warning,   // 51~75
    Danger,    // 76~100
}

pub struct DiseaseRiskEngine {
    /// 검진 결과로 8대 영역 리스크 점수 산출
    pub fn calculate_risk(
        &self,
        results: &[MeasurementResult],
        patient_profile: &PatientProfile,  // 나이, 성별, BMI, 가족력
    ) -> Result<Vec<DiseaseRiskScore>, RiskError>;

    /// 이전 검진과 비교하여 추세 분석
    pub fn compare_with_previous(
        &self,
        current: &[DiseaseRiskScore],
        previous: &[DiseaseRiskScore],
    ) -> Vec<RiskTrendAnalysis>;

    /// 관리 권고사항 생성
    pub fn generate_recommendations(
        &self,
        risks: &[DiseaseRiskScore],
        profile: &PatientProfile,
    ) -> Vec<HealthRecommendation>;
}
```

#### 테스트 기준 (7개)

| # | 테스트 | 기대 결과 |
|---|---|---|
| 1 | 모든 수치 정상 | 8개 영역 전부 RiskLevel::Normal |
| 2 | 혈당 250mg/dL | DiabetesGlycemic = Danger (score > 75) |
| 3 | LDL 190 + CRP 5 | Cardiovascular = Warning |
| 4 | Quick 검진 (3항목만) | 3개 영역만 점수, 나머지 None |
| 5 | compare_with_previous 혈당 개선 | trend = Improving |
| 6 | generate_recommendations Danger 영역 | "즉시 의료기관 방문" 포함 |
| 7 | contributing_factors | 혈당 250이 기여도 1위 |

---

### 6.10 AI 자동분류 엔진 (AutoClassifier)

**파일**: `checkup/auto_classifier.rs`
**역할**: 축적된 측정 데이터를 도메인/카테고리별로 자동 분류 + 통계 요약

#### 측정 데이터 도메인

```rust
pub enum MeasurementDomain {
    Health(HealthCategory),
    Environment(EnvironmentCategory),
    Safety(SafetyCategory),
    Industrial(IndustrialCategory),
}

pub enum HealthCategory {
    BloodGlucose,
    Lipids,
    Liver,
    Kidney,
    Thyroid,
    Electrolytes,
    Inflammation,
    Hormone,
}

pub struct AutoClassifier {
    /// CartridgeManifest의 cartridge_type으로 도메인 자동 분류
    pub fn classify_by_cartridge(
        &self,
        manifest: &CartridgeManifest,
    ) -> MeasurementDomain;

    /// 측정 결과에 AI 스마트 태그 부여
    pub fn assign_smart_tags(
        &self,
        result: &MeasurementResult,
        history: &[MeasurementResult],
    ) -> Vec<SmartTag>;

    /// 카테고리별 통계 요약 생성
    pub fn generate_category_stats(
        &self,
        results: &[MeasurementResult],
        period: TimePeriod,
    ) -> Vec<CategoryStats>;
}

pub struct SmartTag {
    pub tag: String,        // "급격한상승", "안정적", "주기적패턴", "이상값"
    pub confidence: f64,    // 0.0~1.0
    pub explanation: String,
}

pub struct CategoryStats {
    pub domain: MeasurementDomain,
    pub measurement_count: u32,
    pub mean: f64,
    pub std_dev: f64,
    pub trend_direction: TrendDirection,  // Rising/Stable/Falling
    pub anomaly_count: u32,
    pub last_measured: u64,
}
```

#### 테스트 기준 (5개)

| # | 테스트 | 기대 결과 |
|---|---|---|
| 1 | 혈당 카트리지 type=0x0001 | Health(BloodGlucose) |
| 2 | VOC 카트리지 type=0x0101 | Environment(AirQuality) |
| 3 | assign_smart_tags 3회 연속 상승 | "급격한상승" 태그 생성 |
| 4 | generate_category_stats 30일 혈당 10회 | mean/std_dev 정확 |
| 5 | 미등록 카트리지 type=0xFFFF | Industrial(Unknown) 기본값 |

---

### 6.11 ContextEngine — 유기적 연동 엔진 (v2.2 신규)

**파일**: `context/engine.rs`
**역할**: 사용자의 현재 상태(최근 측정, 습관, 환경, 구매이력, 시간)를 종합하여 상황별 최적 행동을 추천하는 핵심 엔진. 모든 측정 데이터가 인사이트→영양추천→제품추천→복용추적→재측정의 순환 루프를 형성.

#### 핵심 구조

```rust
/// 사용자 컨텍스트 스냅샷 — 모든 데이터 소스의 현재 상태
pub struct UserContext {
    pub recent_measurements: Vec<MeasurementResult>,  // 최근 7일
    pub active_risks: Vec<DiseaseRiskScore>,           // 현재 리스크
    pub nutrition_status: NutritionStatus,              // 영양 현황
    pub habit_streaks: Vec<HabitStreak>,                // 습관 스트릭
    pub supplement_log: Vec<SupplementEntry>,           // 복용 기록
    pub environment: Option<EnvironmentSnapshot>,       // 환경 데이터
    pub purchase_history: Vec<PurchaseRecord>,          // 구매 이력
    pub user_profile: PatientProfile,                   // 프로필
    pub timestamp: u64,
}

/// 컨텍스트 카드 — 화면에 표시되는 지능형 추천 단위
pub struct ContextCard {
    pub card_id: Uuid,
    pub card_type: ContextCardType,
    pub priority: u8,              // 1~100 (높을수록 먼저)
    pub title: String,
    pub body: String,
    pub action: CardAction,        // 탭 시 이동할 화면 + 파라미터
    pub evidence: Vec<Evidence>,   // 추천 근거
    pub expires_at: Option<u64>,
    pub dismissed: bool,
}

pub enum ContextCardType {
    NutritionRecommendation,     // 영양 추천
    ProductSuggestion,           // 제품 제안
    HabitNudge,                  // 습관 넛지
    CheckupReminder,             // 검진 알림
    EnvironmentAlert,            // 환경 연계
    DataInsight,                 // 데이터 인사이트
    GoalProgress,                // 목표 진척
    RepurchaseReminder,          // 재구매 알림
    MedicationReminder,          // 복용 알림
    AchievementUnlocked,         // 달성 보상
}

pub struct ContextEngine {
    /// 전체 사용자 컨텍스트를 수집하여 스냅샷 생성
    pub fn build_context(
        &self,
        user_id: &Uuid,
        data_sources: &DataSources,
    ) -> Result<UserContext, ContextError>;

    /// 컨텍스트 기반 카드 생성 (우선순위 정렬, 최대 10장)
    pub fn generate_cards(
        &self,
        context: &UserContext,
        rules: &[ContextRule],
    ) -> Vec<ContextCard>;

    /// 넛지/푸시알림 스케줄 생성
    pub fn schedule_nudges(
        &self,
        context: &UserContext,
    ) -> Vec<ScheduledNudge>;

    /// 사용자 반응 학습 (카드 탭/무시/닫기 → 향후 우선순위 조정)
    pub fn learn_from_interaction(
        &mut self,
        card_id: &Uuid,
        interaction: CardInteraction,  // Tapped / Dismissed / Ignored
    );
}
```

#### 인터페이스 계약

- build_context()는 7개 데이터 소스를 30초 내 병합, 누락 소스는 Option::None으로 graceful degradation
- generate_cards()는 최대 10장, priority 내림차순, ContextCardType 중복 최대 2장
- 위험(Danger) 레벨 카드는 항상 priority 90+ (안전 > 편의)
- learn_from_interaction()은 Dismissed 3회 이상 → 해당 card_type 우선순위 자동 하향
- evidence 필드는 최소 1개 이상 (빈 근거 카드 금지)

#### 테스트 기준 (8개)

| # | 테스트 | 기대 결과 |
|---|---|---|
| 1 | 혈당 115mg/dL 측정 | NutritionRecommendation 카드 생성 (크롬/마그네슘 권고) |
| 2 | 보충제 잔여 5일분 | RepurchaseReminder 카드 생성 |
| 3 | 습관 3일 미달 | HabitNudge 카드 priority > 50 |
| 4 | Danger 리스크 감지 | 해당 카드 priority ≥ 90 |
| 5 | 동일 card_type 3장 생성 시도 | 2장만 반환 |
| 6 | 환경 데이터 없음 (Option::None) | 에러 없이 나머지 카드 정상 생성 |
| 7 | Dismissed 3회 → 해당 type 다음 생성 | priority 20% 감소 |
| 8 | generate_cards() 반환 → 전체 evidence 비어있는 카드 | 0건 |

---

### 6.12 NutritionAdvisor — 영양 추천 엔진 (v2.2 신규)

**파일**: `context/nutrition.rs`
**역할**: 바이오마커 데이터 + 식단기록 + 보충제 기록으로 개인화 영양 추천. 영양 정보 제공 목적이며 의료 처방이 아님을 명시.

#### 핵심 구조

```rust
/// 영양소 상태 분석 결과
pub struct NutrientAnalysis {
    pub nutrient: Nutrient,
    pub current_intake: f64,           // mg 또는 μg
    pub recommended_intake: f64,        // 개인화 RDA
    pub fulfillment_percent: f64,       // 0~200+
    pub biomarker_correlation: Option<BiomarkerCorrelation>, // 관련 바이오마커
    pub status: NutrientStatus,         // Deficient / Low / Adequate / Excess
    pub recommendation: String,         // "크롬 200μg 보충 권장 (정보 제공용)"
}

pub enum NutrientStatus {
    Deficient,    // <50% RDA 또는 바이오마커 위험
    Low,          // 50~80% RDA
    Adequate,     // 80~120% RDA
    Excess,       // >150% RDA — 과잉 경고
}

pub struct BiomarkerCorrelation {
    pub biomarker_name: String,      // "혈당", "페리틴" 등
    pub current_value: f64,
    pub reference_range: (f64, f64),
    pub correlation_strength: f64,    // -1.0 ~ 1.0
    pub evidence_pmid: Option<String>, // PubMed ID (검증됨만)
}

pub struct NutritionAdvisor {
    /// 바이오마커 기반 영양 분석 (측정 결과 → 영양소 부족/과잉 판정)
    pub fn analyze_from_biomarkers(
        &self,
        measurements: &[MeasurementResult],
        profile: &PatientProfile,
    ) -> Vec<NutrientAnalysis>;

    /// 식단 기록 분석 (일일 섭취 → 영양소별 충족률)
    pub fn analyze_meal_log(
        &self,
        meals: &[MealEntry],
        profile: &PatientProfile,
    ) -> DailyNutritionSummary;

    /// 보충제 효과 분석 (복용 전후 바이오마커 비교)
    pub fn evaluate_supplement_effect(
        &self,
        supplement: &SupplementEntry,
        before_measurements: &[MeasurementResult],
        after_measurements: &[MeasurementResult],
    ) -> SupplementEffectReport;

    /// 영양소 → 추천 식품/보충제 매핑
    pub fn recommend_sources(
        &self,
        deficiencies: &[NutrientAnalysis],
        allergies: &[String],
        current_medications: &[Medication],  // 약물 상호작용 체크
    ) -> Vec<NutritionRecommendation>;
}
```

#### 인터페이스 계약

- analyze_from_biomarkers()의 바이오마커→영양소 매핑: 혈당→크롬/마그네슘/식이섬유, 페리틴→철분, 25(OH)D→비타민D, CRP→오메가3/커큐민, 전해질→Na/K/Mg
- recommend_sources()는 약물 상호작용 체크 후 안전한 것만 반환. 상호작용 있으면 contraindicated 필드에 사유 기록
- SupplementEffectReport는 최소 2주(측정 2회 이상) 데이터 필요, 미달 시 InsufficientData 반환
- evidence_pmid는 [미검증]이면 None, 검증된 근거만 기입

#### 테스트 기준 (7개)

| # | 테스트 | 기대 결과 |
|---|---|---|
| 1 | 혈당 115mg/dL | 크롬 Deficient, 마그네슘 Low 반환 |
| 2 | 페리틴 10ng/mL | 철분 Deficient (status=Deficient, recommend="철분 18mg 보충 권장") |
| 3 | 모든 바이오마커 정상 | 모든 영양소 Adequate |
| 4 | recommend_sources() + 약물 상호작용 있음 | 해당 보충제 제외 + contraindicated 기록 |
| 5 | 알레르기 "유제품" | 유제품 기반 칼슘 추천 제외 |
| 6 | evaluate_supplement_effect() 데이터 1주일만 | InsufficientData |
| 7 | 비타민D 3000IU 4주 복용 후 25(OH)D 15→32 | "효과 있음" + improvement_percent 계산 |

---

### 6.13 ShoppingBridge — 쇼핑 연동 브릿지 (v2.2 신규)

**파일**: `context/shopping.rs`
**역할**: 영양추천→제품매칭→재구매알림→효과추적의 쇼핑 순환 루프

#### 핵심 구조

```rust
pub struct ProductRecommendation {
    pub product_id: String,
    pub product_name: String,
    pub category: ProductCategory,
    pub match_reason: Vec<MatchReason>,    // 왜 이 제품인지
    pub match_score: f64,                   // 0.0~1.0
    pub price: u64,                         // ₩
    pub contraindicated: bool,              // 약물 상호작용 위험
    pub certifications: Vec<String>,        // GMP, KFDA 등
}

pub enum MatchReason {
    BiomarkerDeficiency { nutrient: String, current: f64, target: f64 },
    DiseaseRiskReduction { domain: String, expected_reduction: f64 },
    RepurchaseDue { last_order_date: u64, estimated_remaining_days: u8 },
    CoachRecommendation { session_id: Uuid },
}

pub struct ShoppingBridge {
    /// 영양 분석 → 제품 매칭 (카탈로그 검색 + 스코어링)
    pub fn match_products(
        &self,
        deficiencies: &[NutrientAnalysis],
        catalog: &ProductCatalog,
        user_prefs: &UserPreferences,  // 알레르기, 비건, 예산
    ) -> Vec<ProductRecommendation>;

    /// 재구매 알림 계산 (복용 기록 기반 잔여일 추정)
    pub fn calculate_repurchase_timing(
        &self,
        supplement_log: &[SupplementEntry],
        purchase: &PurchaseRecord,
    ) -> Option<RepurchaseAlert>;

    /// 제품 효과 추적 (구매 후 바이오마커 변화)
    pub fn track_product_effect(
        &self,
        product_id: &str,
        purchase_date: u64,
        measurements_before: &[MeasurementResult],
        measurements_after: &[MeasurementResult],
    ) -> ProductEffectReport;
}
```

#### 인터페이스 계약

- match_products()는 contraindicated=true인 제품은 반환하되 구매 불가 표시
- match_score = 0.4×영양매칭 + 0.2×가격대비가치 + 0.2×인증점수 + 0.1×리뷰평점 + 0.1×개인이력
- calculate_repurchase_timing() — 일일 복용량 × 구매 수량 → 남은 일수, 5일 전 알림
- 모든 추천에 match_reason 최소 1개 (근거 투명성)

#### 테스트 기준 (6개)

| # | 테스트 | 기대 결과 |
|---|---|---|
| 1 | 크롬 Deficient | catalog에서 크롬 보충제 match_score > 0.7 |
| 2 | 약물 상호작용 위험 제품 | contraindicated=true, 구매 차단 |
| 3 | 비건 프리퍼런스 | 동물성 제품 필터링 |
| 4 | 60정 제품 + 일 1정 복용 50일 경과 | 잔여 10일 → RepurchaseAlert(remaining_days: 10) |
| 5 | 제품 구매 후 4주 바이오마커 개선 | ProductEffectReport improvement > 0 |
| 6 | match_reason 비어있는 ProductRecommendation | 생성 불가 (에러) |

---

### 6.14 HabitTracker — 습관 추적 엔진 (v2.2 신규)

**파일**: `context/habit.rs`
**역할**: 5대 습관(수면/운동/수분/식단/스트레스) 추적 + 바이오마커 상관분석 + 게이미피케이션

#### 핵심 구조

```rust
pub enum HabitCategory {
    Sleep,      // 수면 시간 + 품질
    Exercise,   // 운동 유형 + 시간
    Hydration,  // 수분 섭취량
    Diet,       // 식단 품질 점수
    Stress,     // 스트레스 레벨 (자가평가)
}

pub struct HabitEntry {
    pub category: HabitCategory,
    pub date: u32,           // YYYYMMDD
    pub value: f64,          // 시간/리터/점수 등
    pub target: f64,         // 목표값
    pub achieved: bool,      // value >= target
}

pub struct HabitStreak {
    pub category: HabitCategory,
    pub current_streak: u16,     // 연속 달성일
    pub longest_streak: u16,
    pub total_achieved: u32,
    pub achievement_rate_30d: f64, // 최근 30일 달성률
}

pub struct HabitBiomarkerCorrelation {
    pub habit: HabitCategory,
    pub biomarker: String,
    pub correlation: f64,         // -1.0~1.0
    pub sample_size: u32,         // 관측 수
    pub insight: String,          // "운동 주 3회 이상 시 CRP 평균 15% 감소"
    pub confidence: ConfidenceLevel, // Low/Medium/High
}

pub struct HabitTracker {
    /// 일일 습관 기록
    pub fn log_habit(&mut self, entry: HabitEntry) -> Result<HabitStreak, HabitError>;

    /// 습관-바이오마커 상관분석 (최소 30일 데이터)
    pub fn analyze_correlation(
        &self,
        habits: &[HabitEntry],
        measurements: &[MeasurementResult],
    ) -> Vec<HabitBiomarkerCorrelation>;

    /// 게이미피케이션 보상 판정
    pub fn check_achievements(
        &self,
        streaks: &[HabitStreak],
        goals: &[HealthGoal],
    ) -> Vec<Achievement>;

    /// AI 넛지 메시지 생성 (시간+컨텍스트 기반)
    pub fn generate_nudge(
        &self,
        context: &UserContext,
        time_of_day: TimeOfDay,  // Morning/Afternoon/Evening/Night
    ) -> Option<NudgeMessage>;
}
```

#### 인터페이스 계약

- log_habit() 호출 시 스트릭 자동 계산. 하루 중복 입력은 마지막 값으로 덮어쓰기
- analyze_correlation()은 최소 30일, 측정 최소 4회 필요. 미달 시 confidence=Low로 표시
- generate_nudge()는 하루 최대 5회 (과도한 알림 방지), Morning은 습관/복용, Evening은 수면/회고
- check_achievements()의 배지: 7일스트릭(브론즈), 30일(실버), 90일(골드), 365일(플래티넘)

#### 테스트 기준 (6개)

| # | 테스트 | 기대 결과 |
|---|---|---|
| 1 | 7일 연속 운동 달성 | current_streak=7, 브론즈 배지 |
| 2 | 동일 날짜 중복 입력 | 마지막 값으로 덮어쓰기 |
| 3 | 30일 운동 + CRP 4회 측정 | correlation 계산 (수치 반환) |
| 4 | 15일 데이터만 | confidence=Low |
| 5 | generate_nudge(Morning) | 습관/복용 관련 메시지 (수면 관련 아님) |
| 6 | 넛지 6회째 호출 | None 반환 (하루 5회 제한) |

---

## 7. Layer 4: Flutter 모바일 앱

### 7.1 상태 관리 패턴 (Riverpod + freezed)

```dart
// ── 측정 상태 (freezed 불변 객체) ──
@freezed
class MeasurementState with _$MeasurementState {
  const factory MeasurementState.idle() = _Idle;
  const factory MeasurementState.preparing() = _Preparing;
  const factory MeasurementState.measuring({
    required double progress,
    double? currentValue,
    @Default([]) List<double> waveformData,
  }) = _Measuring;
  const factory MeasurementState.completed({
    required InferenceResult result,
    required TwinHealth twinHealth,
  }) = _Completed;
  const factory MeasurementState.error(String message) = _Error;
}

// ── Riverpod Notifier ──
@riverpod
class MeasurementNotifier extends _$MeasurementNotifier {
  @override
  MeasurementState build() => const MeasurementState.idle();

  Future<void> startMeasurement(CartridgeManifest manifest) async {
    state = const MeasurementState.preparing();
    final rustBridge = ref.read(rustBridgeProvider);
    await rustBridge.initSensor(manifest: manifest);

    state = const MeasurementState.measuring(progress: 0.0);
    final bleStream = ref.read(bleStreamProvider);

    await for (final packet in bleStream.waveformStream()) {
      final differential = await rustBridge.applyDifferential(
        sDetection: packet.rawAdc.toDouble(),
        sReference: packet.referenceAdc.toDouble(),
      );
      state = MeasurementState.measuring(
        progress: packet.sequenceNum / manifest.totalSequences,
        currentValue: differential,
        waveformData: [...state.waveformData, differential],
      );
    }

    final result = await rustBridge.runInference(
      fingerprint: state.waveformData,
      modelId: manifest.recommendedModelId,
    );
    final twinStatus = await rustBridge.updateDigitalTwin(measured: result.concentrations);
    state = MeasurementState.completed(result: result, twinHealth: twinStatus);
  }
}
```

### 7.2 BLE 통신 서비스

**패키지**: `flutter_reactive_ble` (1순위) + `flutter_blue_plus` (폴백)

```dart
class BleService {
  /// 리더기 스캔 (ManPaSik Service UUID 필터)
  Stream<DiscoveredDevice> scanForReaders();

  /// 연결 (AES-128 Secure Connections)
  Future<BleConnection> connect(String deviceId);

  /// GATT 0xFF02 Waveform 스트림 구독
  Stream<WaveformPacket> waveformStream();

  /// GATT 0xFF01 설정 쓰기
  Future<void> writeConfig(MeasurementConfig config);

  /// 연결 상태 모니터링 + 자동 재연결
  Stream<ConnectionState> connectionState();
}
```

**테스트 기준**:
1. Mock BLE 기기 → 스캔 → 연결 → Config 쓰기 성공
2. Waveform 스트림 100 패킷 → 순서 보장 (sequence_num)
3. 연결 끊김 → 3초 내 자동 재연결 시도

### 7.3 NFC 카트리지 인식

```dart
class NfcService {
  /// NFC 태그 읽기 → CartridgeManifest 파싱
  Future<CartridgeManifest> readCartridge();

  /// SUN/CMAC 인증 검증
  Future<bool> verifyAuthenticity(CartridgeManifest manifest);

  /// 사용 횟수 업데이트 (NFC 태그 쓰기)
  Future<void> incrementUseCount(String serialNumber);
}
```

**테스트 기준**:
1. v2.0 태그 → 정상 파싱
2. v1.0 태그 → 후방호환 파싱 (from_v1)
3. 만료된 카트리지 → 사용 불가 경고
4. 위조 태그 → SUN/CMAC 인증 실패

### 7.4 CRDT 오프라인 동기화

```dart
class CrdtSyncManager {
  /// LWW-Register: 측정 결과 (타임스탬프 기반 충돌 해결)
  /// G-Counter: 카트리지 사용 횟수 (단조증가)
  /// OR-Set: 카트리지 목록 (관찰된 제거)

  Future<void> saveOffline(MeasurementRecord record);
  Future<void> syncWhenOnline();
  Future<int> pendingCount();
}
```

**테스트 기준**:
1. 오프라인 저장 10건 → pendingCount == 10
2. 온라인 동기화 → pendingCount == 0
3. 동일 레코드 양쪽 수정 → LWW 최신 타임스탬프 승리
4. 100% 오프라인에서 측정 완료 가능

### 7.5 종합검진 UI 상태 관리 (v2.1 신규)

**파일**: `providers/checkup_provider.dart`
**역할**: 다중 카트리지 순차 측정 UI 상태 관리 + 리포트 표시

#### Riverpod 상태 정의

```dart
// ── 종합검진 세션 상태 ──
@freezed
class CheckupSessionState with _$CheckupSessionState {
  const factory CheckupSessionState.idle() = _CheckupIdle;
  const factory CheckupSessionState.packageSelected({
    required CheckupPackage package,
    required List<CartridgeRequirement> requirements,
  }) = _PackageSelected;
  const factory CheckupSessionState.inProgress({
    required int currentStep,
    required int totalSteps,
    required String currentCartridgeType,
    required List<CheckupStepResult> completedSteps,
  }) = _CheckupInProgress;
  const factory CheckupSessionState.analyzing() = _CheckupAnalyzing;
  const factory CheckupSessionState.completed({
    required CheckupReport report,
    required List<DiseaseRiskScore> risks,
    required List<HealthRecommendation> recommendations,
  }) = _CheckupCompleted;
  const factory CheckupSessionState.partial({
    required List<CheckupStepResult> completedSteps,
    required List<CartridgeRequirement> remaining,
  }) = _CheckupPartial;
}

@riverpod
class CheckupNotifier extends _$CheckupNotifier {
  @override
  CheckupSessionState build() => const CheckupSessionState.idle();

  Future<void> selectPackage(CheckupPackage package) async {
    state = const CheckupSessionState.packageSelected(package: package, requirements: []);
    final rustBridge = ref.read(rustBridgeProvider);
    final requirements = await rustBridge.getCheckupRequirements(package);
    state = CheckupSessionState.packageSelected(package: package, requirements: requirements);
  }

  Future<void> startStep(CartridgeManifest manifest) async {
    if (state is! _CheckupInProgress) return;
    final current = state as _CheckupInProgress;
    state = CheckupSessionState.inProgress(
      currentStep: current.currentStep,
      totalSteps: current.totalSteps,
      currentCartridgeType: manifest.cartridge_type.toString(),
      completedSteps: current.completedSteps,
    );
  }

  Future<void> completeStep(MeasurementResult result) async {
    if (state is! _CheckupInProgress) return;
    final current = state as _CheckupInProgress;
    final updated = [...current.completedSteps, CheckupStepResult(result: result)];

    if (current.currentStep >= current.totalSteps) {
      state = const CheckupSessionState.analyzing();
      final rustBridge = ref.read(rustBridgeProvider);
      final risks = await rustBridge.calculateDiseaseRisk(updated);
      final recommendations = await rustBridge.generateRecommendations(risks);
      final report = CheckupReport(steps: updated, risks: risks);
      state = CheckupSessionState.completed(
        report: report,
        risks: risks,
        recommendations: recommendations,
      );
    } else {
      state = CheckupSessionState.inProgress(
        currentStep: current.currentStep + 1,
        totalSteps: current.totalSteps,
        currentCartridgeType: "",
        completedSteps: updated,
      );
    }
  }

  Future<void> savePartial() async {
    if (state is! _CheckupInProgress) return;
    final current = state as _CheckupInProgress;
    final rustBridge = ref.read(rustBridgeProvider);
    await rustBridge.savePartialCheckup(current.completedSteps);
  }

  Future<void> resumeSession(String sessionId) async {
    state = const CheckupSessionState.idle();
    final rustBridge = ref.read(rustBridgeProvider);
    final partial = await rustBridge.resumeCheckup(sessionId);
    state = CheckupSessionState.partial(
      completedSteps: partial.completedSteps,
      remaining: partial.remainingRequirements,
    );
  }
}

// ── 질병 리스크 Provider ──
@riverpod
Future<List<DiseaseRiskScore>> diseaseRisk(DiseaseRiskRef ref, String sessionId) async {
  final api = ref.read(apiServiceProvider);
  return api.getDiseaseRisk(sessionId);
}

// ── AI 자동분류 Provider ──
@riverpod
Future<List<CategoryStats>> autoClassification(AutoClassificationRef ref, TimePeriod period) async {
  final api = ref.read(apiServiceProvider);
  return api.getCategoryStats(period);
}
```

#### 화면 구조 (v2.1 추가)

```
만파식 앱
├── 홈           — 통합 스코어 + 최근 측정 + 긴급 알림
├── 종합검진      — ★v2.1 신규 ★ 패키지 선택 → 다중 카트리지 순차 측정 → 리스크 리포트
├── 측정         — NFC 스캔 → 실시간 파형 → 결과 (4단계: 정상/주의/경고/위험)
├── 데이터 허브   — 시계열 차트 + 리포트(PDF/FHIR) + 내보내기
├── AI 코치      — 대화형 상담 + 약물 상호작용 + 예측 경보
├── 마켓플레이스  — 카트리지 스토어 + 구독 + SDK 링크
└── 더보기       — 리더기 관리 + 가족 + 의료 연동 + 스마트홈 + 접근성
```

---

### 7.6 유기적 연동 UI 상태관리 (v2.2 신규)

**파일**: `lib/providers/` 내 Riverpod providers
**역할**: 컨텍스트 카드, 영양 분석, 습관 추적, 쇼핑 추천을 한 번에 관리하는 통합 상태 레이어

#### Riverpod Provider 정의

```dart
// ── 컨텍스트 카드 피드 ──
@riverpod
Future<List<ContextCard>> contextCardFeed(ContextCardFeedRef ref) async {
  final api = ref.read(apiServiceProvider);
  final context = await api.buildUserContext();
  return api.generateContextCards(context);
}

@riverpod
class CardInteractionNotifier extends _$CardInteractionNotifier {
  Future<void> recordInteraction(String cardId, String interaction) async {
    final api = ref.read(apiServiceProvider);
    await api.recordCardInteraction(cardId, interaction);
    // 로컬 상태 업데이트 후 리드
    state = await asyncValue;
  }
}

// ── 영양 어드바이저 ──
@freezed
class NutritionState with _$NutritionState {
  const factory NutritionState.loading() = _NutritionLoading;
  const factory NutritionState.loaded({
    required DailyNutritionSummary daily,
    required List<NutrientAnalysis> biomarkerBased,
    required List<NutritionRecommendation> recommendations,
  }) = _NutritionLoaded;
  const factory NutritionState.error(String message) = _NutritionError;
}

@riverpod
class NutritionNotifier extends _$NutritionNotifier {
  @override
  Future<NutritionState> build() async {
    final api = ref.read(apiServiceProvider);
    try {
      final report = await api.analyzeNutrition();
      return NutritionState.loaded(
        daily: report.dailySummary,
        biomarkerBased: report.biomarkerAnalysis,
        recommendations: report.recommendations,
      );
    } catch (e) {
      return NutritionState.error(e.toString());
    }
  }

  Future<void> logMeal(MealEntry meal) async {
    final api = ref.read(apiServiceProvider);
    await api.logMeal(meal);
    state = await AsyncValue.guard(() => build());
  }

  Future<void> logSupplement(SupplementEntry supplement) async {
    final api = ref.read(apiServiceProvider);
    await api.logSupplement(supplement);
    state = await AsyncValue.guard(() => build());
  }
}

// ── 습관 트래커 ──
@freezed
class HabitState with _$HabitState {
  const factory HabitState({
    required Map<String, HabitEntry> todayEntries,
    required List<HabitStreak> streaks,
    required List<Achievement> recentAchievements,
    required List<HabitBiomarkerCorrelation> correlations,
  }) = _HabitState;
}

@riverpod
class HabitNotifier extends _$HabitNotifier {
  @override
  Future<HabitState> build() async {
    final api = ref.read(apiServiceProvider);
    final habits = await api.getHabits();
    final achievements = await api.getAchievements();
    final correlations = await api.getHabitCorrelations();

    return HabitState(
      todayEntries: Map.fromIterable(habits, key: (h) => h.category),
      streaks: habits.map((h) => h.streak).toList(),
      recentAchievements: achievements,
      correlations: correlations,
    );
  }

  Future<void> logHabit(HabitEntry entry) async {
    final api = ref.read(apiServiceProvider);
    await api.logHabit(entry);
    state = await AsyncValue.guard(() => build());
  }
}

// ── 추천 제품 ──
@riverpod
Future<List<ProductRecommendation>> recommendedProducts(
  RecommendedProductsRef ref, {
  List<String>? nutrientFilter
}) async {
  final api = ref.read(apiServiceProvider);
  return api.getRecommendedProducts(nutrientFilter: nutrientFilter);
}

@riverpod
Future<List<RepurchaseAlert>> repurchaseAlerts(RepurchaseAlertsRef ref) async {
  final api = ref.read(apiServiceProvider);
  return api.getRepurchaseAlerts();
}

// ── 정기배송 관리 ──
@freezed
class SubscriptionState with _$SubscriptionState {
  const factory SubscriptionState({
    required List<SubscriptionDelivery> active,
    required List<SubscriptionDelivery> paused,
  }) = _SubscriptionState;
}

@riverpod
class SubscriptionDeliveryNotifier extends _$SubscriptionDeliveryNotifier {
  @override
  Future<SubscriptionState> build() async {
    final api = ref.read(apiServiceProvider);
    final all = await api.getSubscriptions();
    return SubscriptionState(
      active: all.where((s) => s.status == 'active').toList(),
      paused: all.where((s) => s.status == 'paused').toList(),
    );
  }

  Future<void> pauseSubscription(String subscriptionId) async {
    final api = ref.read(apiServiceProvider);
    await api.updateSubscription(subscriptionId, status: 'paused');
    state = await AsyncValue.guard(() => build());
  }

  Future<void> resumeSubscription(String subscriptionId) async {
    final api = ref.read(apiServiceProvider);
    await api.updateSubscription(subscriptionId, status: 'active');
    state = await AsyncValue.guard(() => build());
  }
}

// ── 포인트/리워드 ──
@riverpod
Future<RewardStatus> rewardStatus(RewardStatusRef ref) async {
  final api = ref.read(apiServiceProvider);
  return api.getRewardStatus();
}

@riverpod
class RewardNotifier extends _$RewardNotifier {
  Future<void> redeemPoints(int points) async {
    final api = ref.read(apiServiceProvider);
    await api.redeemPoints(points);
    ref.invalidate(rewardStatus);
  }
}

// ── 넛지 알림 스케줄 ──
@riverpod
Future<List<NudgeMessage>> scheduledNudges(ScheduledNudgesRef ref) async {
  final api = ref.read(apiServiceProvider);
  return api.getScheduledNudges();
}
```

#### 화면 상태 계층

- **ContextCardFeed**: 매 30초 자동 새로고침, Dismissed 인터랙션 즉시 반영
- **NutritionNotifier**: 식단/보충제 로그 입력 시 자동 리로드
- **HabitNotifier**: 하루 1회 새로고침(자정), 당일 입력 즉시 반영
- **RecommendedProducts**: 영양 분석 변경 시 연쇄 리로드
- **SubscriptionDeliveryNotifier**: 배송 상태 변경 시 푸시 알림 + 로컬 업데이트

---

### 7.7 자가검증 UI 상태관리 (v2.3 신규)

**파일**: `flutter-app/lib/providers/diagnostics_providers.dart`
**의존**: §6.15 SelfDiagnostics FFI, §7.1 BLE Provider

#### Riverpod Provider 구조

```dart
// ── 시스템 건강 점수 (15분 자동 갱신) ──
@riverpod
class SystemHealthNotifier extends _$SystemHealthNotifier {
  // Rust FFI: SystemHealthOrchestrator::compute_health()
  // 15분 Timer + 측정 전후 강제 리프레시
  // 출력: SystemHealthScore (overall, grade, layers, trend, measurement_ready)
}

// ── 측정 전 자가검증 ──
@riverpod
Future<PreFlightResult> preFlightCheck(ref) async {
  // Rust FFI: SystemHealthOrchestrator::pre_flight_check()
  // M-010 측정 대기 진입 시 자동 호출
  // ready=false → M-085 자가검증 실패 화면으로 라우팅
}

// ── HW 컴포넌트 상세 ──
@riverpod
class HwComponentDetailNotifier extends _$HwComponentDetailNotifier {
  // 특정 HwComponent 선택 → diagnose_component() 호출
  // H-035 HW 진단 상세 화면에서 사용
}

// ── 예방정비 스케줄 ──
@riverpod
class MaintenanceScheduleNotifier extends _$MaintenanceScheduleNotifier {
  // Rust FFI: PredictiveMaintenanceEngine::analyze_trends()
  // 일 1회 갱신, 푸시 알림 트리거
}

// ── 진단 보고서 ──
@riverpod
Future<UserDiagnosticSummary> diagnosticReport(ref) async {
  // Rust FFI: ErrorReporter::generate_user_report()
  // H-036 진단 보고서 화면에서 사용
}
```

#### 상태 전이

```
앱 시작 → systemHealthNotifier 초기화 (Rust FFI background_check)
    ├── grade=Excellent/Good → 정상 동작
    ├── grade=Fair → H-010 홈에 주의 배너 표시
    ├── grade=Poor → OV-080 "시스템 점검 필요" 오버레이
    └── grade=Critical → 측정 기능 비활성화 + OV-081 "서비스 필요" 모달

측정 시작 → preFlightCheck 호출
    ├── ready=true → M-020 카트리지 스캔 진행
    ├── ready=false + HW → M-085 HW 점검 안내
    └── ready=false + FW → M-086 FW 업데이트 안내
```

---

### 7.8 Self-Healing 복원 UI (v2.4 신규)

**역할**: 자가치유 이벤트의 실시간 시각화, 세션 복원 다이얼로그, 힐링 로그 열람

```dart
/// AppCrashRecoveryProvider — 앱 비정상 종료 후 복원 관리
@riverpod
class AppCrashRecoveryNotifier extends _$AppCrashRecoveryNotifier {
  @override
  AsyncValue<CrashRecoveryState> build() => const AsyncLoading();

  /// 앱 시작 시 이전 세션 크래시 여부 확인
  Future<void> checkPreviousCrash() async {
    state = const AsyncLoading();
    final checkpoint = await _rustBridge.findRecoverableSession();
    if (checkpoint != null) {
      state = AsyncData(CrashRecoveryState(
        hasRecoverableSession: true,
        sessionId: checkpoint.sessionId,
        progressPercent: checkpoint.progressPercent,
        measurementType: checkpoint.measurementType,
        crashTimestamp: DateTime.fromMillisecondsSinceEpoch(checkpoint.timestamp),
      ));
    } else {
      state = const AsyncData(CrashRecoveryState.noRecovery());
    }
  }

  /// 사용자 승인 후 세션 복원 실행
  Future<SessionRestoreResult> restoreSession(String sessionId) async {
    final result = await _rustBridge.restoreSession(sessionId);
    if (result.isSuccess) {
      // 복원된 측정 화면으로 네비게이션
      ref.read(navigationProvider.notifier).navigateToRestoredMeasurement(result);
    }
    return result;
  }
}

/// HealingLogNotifier — 자가치유 이벤트 로그 조회
@riverpod
class HealingLogNotifier extends _$HealingLogNotifier {
  @override
  AsyncValue<List<HealingEventDisplay>> build() => const AsyncLoading();

  Future<void> loadHealingLog({int hours = 24}) async {
    state = const AsyncLoading();
    final events = await _rustBridge.getHealingLog(hours);
    state = AsyncData(events.map((e) => HealingEventDisplay(
      timestamp: e.timestamp,
      layer: e.layer,
      errorType: e.errorType,
      strategy: e.strategy,
      outcome: e.outcome,        // FullyHealed / PartiallyHealed / Failed
      durationMs: e.durationMs,
      icon: _outcomeIcon(e.outcome),
      color: _outcomeColor(e.outcome),
    )).toList());
  }

  String _outcomeIcon(HealingOutcome o) => switch (o) {
    HealingOutcome.fullyHealed     => '✅',
    HealingOutcome.partiallyHealed => '⚠️',
    HealingOutcome.healingFailed   => '❌',
    HealingOutcome.notNeeded       => 'ℹ️',
  };
}

/// HealingStatsNotifier — 치유 통계 대시보드
@riverpod
class HealingStatsNotifier extends _$HealingStatsNotifier {
  @override
  AsyncValue<HealingStatsSummary> build() => const AsyncLoading();

  Future<void> loadStats() async {
    final stats = await _rustBridge.getHealingStatistics();
    state = AsyncData(HealingStatsSummary(
      totalAttempts: stats.totalAttempts,
      successRate: stats.successRate,            // 0.0 ~ 1.0
      avgRecoveryMs: stats.avgRecoveryMs,
      byLayer: {
        'HW': stats.hwRecoveries,
        'FW': stats.fwRepairs,
        'Data': stats.dataReprocesses,
        'Session': stats.sessionRestores,
      },
      topStrategies: stats.topStrategies,        // 가장 많이 사용된 복구 전략 Top 5
      lastHealingEvent: stats.lastEventTimestamp,
    ));
  }
}
```

**위젯 (3개)**:

| 위젯 | 역할 | 위치 |
|---|---|---|
| `MpsSessionRestoreDialog` | 복원 가능 세션 발견 시 모달 — 진행률·측정유형·시간 표시, 복원/폐기 선택 | 앱 시작 시 오버레이 |
| `MpsHealingLogTimeline` | 24시간 치유 이벤트 타임라인 — 레이어별 색상, 성공/실패 아이콘 | H-038 힐링 로그 화면 |
| `MpsHealingStatsCard` | 치유 성공률·평균 복구시간·레이어별 통계 미니 카드 | H-035 진단 센터 하단 |

**화면 네비게이션 흐름**:

```
앱 시작
  ├── RecoverableSession 발견 → MpsSessionRestoreDialog 표시
  │     ├── "복원" → Rust restoreSession → M-020 측정 화면 (중단 지점)
  │     └── "폐기" → cleanupExpired → 정상 홈 화면
  │
  └── 자가치유 발생 중 (측정 도중)
        ├── 백그라운드 치유 성공 → 토스트 "센서 자동 보정 완료"
        ├── 부분 성공 → OV-083 치유 상태 오버레이 (경고)
        └── 치유 실패 → M-087 측정 중단 안내 + ServiceTicket
```

**테스트 기준 (5개)**:
1. `checkPreviousCrash()` — 복원 가능 세션 존재 → 다이얼로그 표시
2. `checkPreviousCrash()` — 복원 가능 세션 없음 → 정상 진행
3. `restoreSession()` — 유효한 checkpoint → 측정 화면 복원 + 진행률 유지
4. `loadHealingLog(24)` — 이벤트 3건 → 타임라인에 3건 표시, 시간순 정렬
5. `loadStats()` — successRate 0.85 → "85%" 표시, 레이어별 차트 렌더링

---

## 8. Layer 5: Go gRPC 백엔드

### 8.1 기술 스택

```yaml
API:         Go 1.22+ gRPC (REST는 gRPC-Gateway로 변환)
Gateway:     Kong / Envoy
메시지큐:    Apache Kafka (이벤트 소싱)
DB:          PostgreSQL 16 + TimescaleDB + Redis + Milvus + ES
인증:        Keycloak (OAuth 2.0 + JWT)
AI 서빙:     Triton Inference Server (GPU)
```

### 8.2 서비스 구성 (21개, 7개 도메인)

```
[사용자]     user-svc, auth-svc, family-svc
[측정]       measurement-svc, calibration-svc, cartridge-svc, data-pipeline
[AI]         ai-inference, ai-training, digital-twin, xai-service
[의료]       health-record-svc, telemedicine-svc, emergency-svc
[커머스]     subscription-svc, marketplace-svc, payment-svc
[관리자]     admin-svc, analytics-svc, audit-svc
[IoT]        device-management-svc, ota-service, community-svc
```

### 8.3 핵심 Protobuf 정의

```protobuf
service MeasurementService {
  rpc StreamMeasurement(stream MeasurementChunk) returns (MeasurementResult);
  rpc GetCalibration(CalibrationRequest) returns (CalibrationData);
  rpc GetMeasurementHistory(HistoryRequest) returns (stream MeasurementRecord);
  rpc SyncDigitalTwin(TwinState) returns (TwinPrediction);
}

service CartridgeService {
  rpc RegisterCartridge(CartridgeDefinition) returns (RegistrationResult);
  rpc ValidateCompatibility(CompatibilityRequest) returns (CompatibilityReport);
  rpc ListCartridges(ListRequest) returns (CartridgeCatalog);
}

service AIService {
  rpc Predict(PredictionRequest) returns (PredictionResult);
  rpc AggregateFederatedWeights(stream ModelWeights) returns (AggregatedModel);
  rpc GetApprovedModel(ModelRequest) returns (ModelArtifact);
}
```

### 8.4 종합검진 gRPC 서비스 (v2.1 신규)

**파일**: `proto/checkup.proto` + `internal/checkup-svc/`
**역할**: 다중 카트리지 검진 세션 관리 + 질병 리스크 산출 + 이력 비교

#### Protobuf 정의

```protobuf
service CheckupService {
  rpc CreateSession(CreateSessionRequest) returns (CheckupSession);
  rpc AddStepResult(AddStepResultRequest) returns (CheckupProgress);
  rpc CompleteSession(CompleteSessionRequest) returns (CheckupReport);
  rpc GetDiseaseRisk(GetDiseaseRiskRequest) returns (DiseaseRiskResponse);
  rpc GetCheckupHistory(GetHistoryRequest) returns (CheckupHistoryResponse);
  rpc CompareCheckups(CompareRequest) returns (CompareResponse);
}

// REST Gateway 엔드포인트
POST   /checkup/sessions              → CreateSession
POST   /checkup/sessions/{id}/steps   → AddStepResult
POST   /checkup/sessions/{id}/complete → CompleteSession
GET    /checkup/sessions/{id}/risk    → GetDiseaseRisk
GET    /checkup/history               → GetCheckupHistory
GET    /checkup/compare?ids=a,b       → CompareCheckups
```

#### 메시지 정의

```protobuf
message CreateSessionRequest {
  string user_id = 1;
  string package_type = 2;  // 'quick'/'standard'/'premium'/'custom'
  repeated uint32 custom_cartridge_types = 3;  // custom 패키지용
}

message CheckupSession {
  string session_id = 1;
  string user_id = 2;
  string package_type = 3;
  string status = 4;  // 'in_progress'/'completed'/'partial'
  int64 started_at_unix = 5;
  int64 completed_at_unix = 6;
  repeated CheckupStep completed_steps = 7;
}

message CheckupStep {
  int32 step_number = 1;
  uint32 cartridge_type = 2;
  string measurement_id = 3;
  int64 completed_at_unix = 4;
}

message AddStepResultRequest {
  string session_id = 1;
  string measurement_id = 2;
  bytes measurement_data = 3;  // Rust-serialized MeasurementResult
}

message CheckupProgress {
  string session_id = 1;
  int32 current_step = 2;
  int32 total_steps = 3;
  string status = 4;
}

message CompleteSessionRequest {
  string session_id = 1;
  string patient_profile_json = 2;  // age, sex, BMI, family_history
}

message CheckupReport {
  string session_id = 1;
  repeated CheckupStep steps = 2;
  repeated DiseaseRiskScore risks = 3;
  repeated HealthRecommendation recommendations = 4;
  int64 generated_at_unix = 5;
}

message DiseaseRiskScore {
  string domain = 1;  // 'metabolic_syndrome', 'diabetes_glycemic', etc
  double score = 2;   // 0.0~100.0
  string level = 3;   // 'normal'/'caution'/'warning'/'danger'
  repeated ContributingFactor factors = 4;
  string trend = 5;   // 'improving'/'stable'/'worsening'
}

message ContributingFactor {
  string measurement_name = 1;
  double value = 2;
  double contribution_weight = 3;  // 0.0~1.0
}

message HealthRecommendation {
  string category = 1;  // 'immediate'/'urgent'/'routine'
  string text = 2;
  int32 priority = 3;   // 1~5
}

message GetDiseaseRiskRequest {
  string session_id = 1;
}

message DiseaseRiskResponse {
  repeated DiseaseRiskScore risks = 1;
  string summary = 2;
}

message GetHistoryRequest {
  string user_id = 1;
  int32 limit = 2;
  string offset_token = 3;
}

message CheckupHistoryResponse {
  repeated CheckupSessionSummary sessions = 1;
  string next_offset_token = 2;
}

message CheckupSessionSummary {
  string session_id = 1;
  string package_type = 2;
  int64 completed_at_unix = 3;
  string overall_risk_level = 4;  // 'normal'/'caution'/'warning'/'danger'
}

message CompareRequest {
  string user_id = 1;
  repeated string session_ids = 2;  // 2개 이상
}

message CompareResponse {
  repeated DomainComparison comparisons = 1;
  string trend_summary = 2;
}

message DomainComparison {
  string domain = 1;
  repeated ScoreHistory history = 2;
}

message ScoreHistory {
  string session_id = 1;
  int64 measured_at_unix = 2;
  double risk_score = 3;
  string risk_level = 4;
}
```

#### DB 마이그레이션

```sql
CREATE TABLE checkup_sessions (
    session_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    package_type    VARCHAR(20) NOT NULL,
    status          VARCHAR(20) NOT NULL,
    started_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at    TIMESTAMPTZ,
    report_json     JSONB,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    INDEX idx_user_created (user_id, created_at DESC)
);

CREATE TABLE checkup_steps (
    step_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id      UUID NOT NULL REFERENCES checkup_sessions(session_id),
    step_number     SMALLINT NOT NULL,
    cartridge_type  INTEGER NOT NULL,
    measurement_id  UUID NOT NULL REFERENCES measurements(id),
    completed_at    TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(session_id, step_number)
);

CREATE TABLE disease_risk_scores (
    risk_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id      UUID NOT NULL REFERENCES checkup_sessions(session_id),
    domain          VARCHAR(30) NOT NULL,
    score           DECIMAL(5,2) NOT NULL,
    level           VARCHAR(10) NOT NULL,
    factors_json    JSONB,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    INDEX idx_session_domain (session_id, domain)
);

CREATE TABLE auto_classifications (
    classification_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    measurement_id    UUID NOT NULL REFERENCES measurements(id),
    domain            VARCHAR(30) NOT NULL,
    category          VARCHAR(30) NOT NULL,
    smart_tags        JSONB,
    created_at        TIMESTAMPTZ DEFAULT NOW(),
    INDEX idx_measurement (measurement_id)
);

-- 검진 이력 조회용 VIEW
CREATE VIEW checkup_summary AS
SELECT
    cs.session_id,
    cs.user_id,
    cs.package_type,
    cs.completed_at,
    MAX(CASE WHEN drs.level = 'danger' THEN 'danger'
         WHEN drs.level = 'warning' THEN 'warning'
         WHEN drs.level = 'caution' THEN 'caution'
         ELSE 'normal' END) as overall_risk_level
FROM checkup_sessions cs
LEFT JOIN disease_risk_scores drs ON cs.session_id = drs.session_id
WHERE cs.status = 'completed'
GROUP BY cs.session_id, cs.user_id, cs.package_type, cs.completed_at;
```

#### 테스트 기준

| 시나리오 | 검증 항목 |
|---|---|
| Standard 검진 E2E | 패키지 선택 → 3카트리지 순차 NFC → 3회 측정 → AI 리스크 산출 → 리포트 |
| 검진 중간이탈+복원 | 2/3 완료 → 앱 종료 → 재시작 → resume → 3/3 완료 |
| 리스크 비교 | 이전 검진과 비교 → 개선/악화 추세 표시 |
| 이력 조회 | 최근 10개 검진 세션 조회 + 페이지네이션 |

---

### 8.5 유기적 연동 gRPC 서비스 (v2.2 신규)

**역할**: ContextEngine, NutritionAdvisor, ShoppingBridge, HabitTracker의 Rust Core 기능을 gRPC로 노출. REST Gateway는 Fiber로 구현하여 모바일 클라이언트와 웹 대시보드 모두 지원.

#### Protobuf 정의 (핵심)

```protobuf
syntax = "proto3";
package manpasik.context.v1;

service ContextService {
  rpc BuildUserContext(UserContextRequest) returns (UserContext);
  rpc GenerateContextCards(GenerateCardsRequest) returns (ContextCardList);
  rpc RecordCardInteraction(CardInteractionRequest) returns (Empty);
}

service NutritionService {
  rpc AnalyzeNutrition(NutritionAnalysisRequest) returns (NutritionReport);
  rpc LogMeal(LogMealRequest) returns (MealEntry);
  rpc LogSupplement(LogSupplementRequest) returns (SupplementEntry);
  rpc GetSupplementEffect(SupplementEffectRequest) returns (SupplementEffectReport);
  rpc GetRecommendedSources(RecommendedSourcesRequest) returns (RecommendationList);
}

service ShoppingService {
  rpc GetRecommendedProducts(ProductRecommendationRequest) returns (ProductList);
  rpc GetRepurchaseAlerts(RepurchaseAlertRequest) returns (RepurchaseAlertList);
  rpc TrackProductEffect(ProductEffectRequest) returns (ProductEffectReport);
  rpc CreateSubscription(CreateSubscriptionRequest) returns (SubscriptionDelivery);
  rpc UpdateSubscription(UpdateSubscriptionRequest) returns (SubscriptionDelivery);
}

service HabitService {
  rpc LogHabit(LogHabitRequest) returns (HabitStreak);
  rpc GetHabitCorrelations(HabitCorrelationRequest) returns (CorrelationList);
  rpc GetAchievements(AchievementRequest) returns (AchievementList);
  rpc GetNudge(NudgeRequest) returns (NudgeMessage);
}

service RewardService {
  rpc GetRewardStatus(RewardStatusRequest) returns (RewardStatus);
  rpc RedeemPoints(RedeemRequest) returns (RedeemResult);
}

// 메시지 정의 (요약)
message UserContext {
  repeated MeasurementResult recent_measurements = 1;
  repeated DiseaseRiskScore active_risks = 2;
  NutritionStatus nutrition_status = 3;
  repeated HabitStreak habit_streaks = 4;
  repeated SupplementEntry supplement_log = 5;
  EnvironmentSnapshot environment = 6;
  repeated PurchaseRecord purchase_history = 7;
  PatientProfile user_profile = 8;
  int64 timestamp = 9;
}

message ContextCard {
  string card_id = 1;
  string card_type = 2;  // enum: NutritionRecommendation/ProductSuggestion/HabitNudge/...
  int32 priority = 3;
  string title = 4;
  string body = 5;
  CardAction action = 6;
  repeated Evidence evidence = 7;
  int64 expires_at = 8;
  bool dismissed = 9;
}

message NutritionReport {
  DailyNutritionSummary daily_summary = 1;
  repeated NutrientAnalysis biomarker_analysis = 2;
  repeated NutritionRecommendation recommendations = 3;
}

message ProductRecommendation {
  string product_id = 1;
  string product_name = 2;
  string category = 3;
  repeated string match_reason = 4;
  double match_score = 5;
  int64 price_krw = 6;
  bool contraindicated = 7;
  repeated string certifications = 8;
}

message SubscriptionDelivery {
  string subscription_id = 1;
  string user_id = 2;
  string product_id = 3;
  int32 interval_days = 4;
  string next_delivery = 5;  // YYYY-MM-DD
  string status = 6;  // active/paused/cancelled
}

message RewardStatus {
  int32 total_points = 1;
  repeated PointSource point_sources = 2;
  repeated RedeemablePrize prizes = 3;
}
```

#### REST Gateway 엔드포인트 (Fiber)

```
POST   /context/v1/build                         → BuildUserContext
POST   /context/v1/cards                         → GenerateContextCards
POST   /context/v1/cards/{id}/interaction        → RecordCardInteraction

POST   /nutrition/v1/analyze                     → AnalyzeNutrition
POST   /nutrition/v1/meals/log                   → LogMeal
POST   /nutrition/v1/supplements/log              → LogSupplement
GET    /nutrition/v1/supplements/{id}/effect     → GetSupplementEffect
GET    /nutrition/v1/recommended-sources         → GetRecommendedSources

GET    /shopping/v1/recommended-products         → GetRecommendedProducts
GET    /shopping/v1/repurchase-alerts            → GetRepurchaseAlerts
GET    /shopping/v1/products/{id}/effect         → TrackProductEffect
POST   /shopping/v1/subscriptions                → CreateSubscription
PUT    /shopping/v1/subscriptions/{id}           → UpdateSubscription

POST   /habits/v1/log                            → LogHabit
GET    /habits/v1/correlations                   → GetHabitCorrelations
GET    /habits/v1/achievements                   → GetAchievements
GET    /habits/v1/nudge                          → GetNudge

GET    /rewards/v1/status                        → GetRewardStatus
POST   /rewards/v1/redeem                        → RedeemPoints
```

#### 데이터베이스 스키마 추가 (v2.2)

```sql
-- ContextEngine 테이블
CREATE TABLE context_cards (
    card_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID NOT NULL REFERENCES users(id),
    card_type     VARCHAR(30) NOT NULL,
    priority      SMALLINT NOT NULL CHECK (priority >= 1 AND priority <= 100),
    title         TEXT NOT NULL,
    body          TEXT,
    action_json   JSONB NOT NULL,
    evidence_json JSONB NOT NULL,
    expires_at    TIMESTAMPTZ,
    dismissed     BOOLEAN DEFAULT FALSE,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_context_cards_user_priority
ON context_cards(user_id, priority DESC, created_at DESC);

CREATE TABLE card_interactions (
    interaction_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id        UUID NOT NULL REFERENCES context_cards(card_id),
    interaction    VARCHAR(10) NOT NULL,  -- 'tapped'/'dismissed'/'ignored'
    created_at     TIMESTAMPTZ DEFAULT NOW()
);

-- NutritionAdvisor 테이블
CREATE TABLE meal_logs (
    meal_id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID NOT NULL REFERENCES users(id),
    meal_type      VARCHAR(10) NOT NULL,    -- 'breakfast'/'lunch'/'dinner'/'snack'
    meal_json      JSONB NOT NULL,
    nutrients_json JSONB,
    logged_at      TIMESTAMPTZ NOT NULL,
    created_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE supplement_logs (
    log_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID NOT NULL REFERENCES users(id),
    product_id     TEXT,
    supplement_name TEXT NOT NULL,
    dosage         TEXT NOT NULL,
    taken_at       TIMESTAMPTZ NOT NULL,
    created_at     TIMESTAMPTZ DEFAULT NOW()
);

-- ShoppingBridge 테이블
CREATE TABLE product_catalog (
    product_id     TEXT PRIMARY KEY,
    product_name   TEXT NOT NULL,
    category       VARCHAR(30) NOT NULL,
    nutrients_json JSONB,
    price_krw      INTEGER NOT NULL,
    certifications TEXT[],
    contraindications TEXT[],
    rating         DECIMAL(2,1),
    review_count   INTEGER DEFAULT 0,
    active         BOOLEAN DEFAULT TRUE,
    created_at     TIMESTAMPTZ DEFAULT NOW(),
    updated_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE subscriptions_delivery (
    subscription_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    product_id      TEXT NOT NULL REFERENCES product_catalog(product_id),
    interval_days   SMALLINT NOT NULL,
    next_delivery   DATE NOT NULL,
    status          VARCHAR(15) DEFAULT 'active',
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_subscriptions_user_next
ON subscriptions_delivery(user_id, next_delivery);

-- HabitTracker 테이블
CREATE TABLE habit_logs (
    log_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id),
    category   VARCHAR(15) NOT NULL,
    value      DECIMAL(8,2) NOT NULL,
    target     DECIMAL(8,2) NOT NULL,
    achieved   BOOLEAN NOT NULL,
    log_date   DATE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, category, log_date)
);

CREATE TABLE achievements (
    achievement_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID NOT NULL REFERENCES users(id),
    badge_type     VARCHAR(30) NOT NULL,
    badge_tier     VARCHAR(10) NOT NULL,  -- 'bronze'/'silver'/'gold'/'platinum'
    unlocked_at    TIMESTAMPTZ NOT NULL,
    created_at     TIMESTAMPTZ DEFAULT NOW()
);

-- Rewards 테이블
CREATE TABLE reward_points (
    point_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id),
    points      INTEGER NOT NULL,
    source      VARCHAR(30) NOT NULL,     -- 'measurement'/'habit'/'purchase'/'referral'
    description TEXT,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_reward_points_user_date
ON reward_points(user_id, created_at DESC);
```

#### 인터페이스 계약

- 모든 gRPC 엔드포인트는 JWT 토큰 기반 인증 필수 (Authorization: Bearer token)
- 응답은 항상 일관된 에러 구조 사용 (code, message, details 필드)
- 타임아웃: 쿼리는 5초, 분석은 15초
- 캐싱: ContextCard는 30초, ProductCatalog는 1시간 레디스 캐시
- 비동기 작업(SupplementEffectReport 계산 등)은 Job ID 반환 후 polling 또는 WebSocket으로 업데이트 전달

#### 테스트 기준

| # | 테스트 | 기대 결과 |
|---|---|---|
| 1 | GenerateContextCards 호출 | 최대 10장 카드 반환, priority 내림차순 |
| 2 | LogMeal → AnalyzeNutrition | 식단 기록 후 영양 분석 즉시 갱신 |
| 3 | LogSupplement → RepurchaseAlert | 보충제 로그 후 5일 이내 재구매 알림 생성 |
| 4 | GetRecommendedProducts 필터 | nutrientFilter=["chrome"] → 크롬 제품만 |
| 5 | CreateSubscription + 5일 경과 | next_delivery 자동 연장 |
| 6 | LogHabit 7일 연속 | check_achievements → Bronze 배지 반환 |
| 7 | GetRewardStatus | total_points 누적, point_sources 분류 |
| 8 | REST Gateway /context/v1/cards | gRPC 응답을 JSON으로 변환 반환 |

---

### 8.6 자가검증 백엔드 서비스 (v2.3 신규)

#### Protobuf 정의

```protobuf
service DiagnosticsService {
  rpc UploadDiagnosticReport(DiagnosticReportRequest) returns (DiagnosticReportResponse);
  rpc GetDeviceHealthHistory(DeviceHealthRequest) returns (DeviceHealthHistoryResponse);
  rpc GetFleetHealthDashboard(FleetHealthRequest) returns (FleetHealthResponse);
  rpc CreateServiceTicket(ServiceTicketRequest) returns (ServiceTicketResponse);
  rpc GetMaintenanceSchedule(MaintenanceRequest) returns (MaintenanceScheduleResponse);
}

message DiagnosticReportRequest {
  string device_id = 1;
  string user_id = 2;
  double overall_score = 3;
  repeated LayerScore layer_scores = 4;
  repeated SystemError active_errors = 5;
  repeated MaintenanceTask pending_tasks = 6;
  int64 timestamp = 7;
}

message LayerScore {
  string layer = 1;     // "hardware", "firmware", "rust_core", "flutter", "backend", "ai"
  double score = 2;
  string status = 3;
}

message FleetHealthResponse {
  int32 total_devices = 1;
  int32 healthy_devices = 2;
  int32 warning_devices = 3;
  int32 critical_devices = 4;
  repeated DeviceHealthSummary top_issues = 5;
  double fleet_average_score = 6;
}
```

#### REST Gateway 추가 엔드포인트

```
POST   /diagnostics/v1/report          → 진단 보고서 업로드
GET    /diagnostics/v1/health/{device}  → 디바이스 건강 이력
GET    /diagnostics/v1/fleet            → 전체 디바이스 대시보드 (관리자)
POST   /diagnostics/v1/ticket           → 서비스 티켓 자동 생성
GET    /diagnostics/v1/maintenance/{device} → 예방정비 스케줄
```

#### DB 스키마 (TimescaleDB)

```sql
-- 디바이스 건강 이력 (시계열)
CREATE TABLE device_health_log (
    id              BIGSERIAL,
    device_id       UUID NOT NULL,
    timestamp       TIMESTAMPTZ NOT NULL,
    overall_score   DOUBLE PRECISION,
    hw_score        DOUBLE PRECISION,
    fw_score        DOUBLE PRECISION,
    rust_score      DOUBLE PRECISION,
    app_score       DOUBLE PRECISION,
    backend_score   DOUBLE PRECISION,
    ai_score        DOUBLE PRECISION,
    grade           VARCHAR(20),
    measurement_ready BOOLEAN,
    PRIMARY KEY (id, timestamp)
);
SELECT create_hypertable('device_health_log', 'timestamp');

-- 시스템 오류 로그
CREATE TABLE system_error_log (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id       UUID NOT NULL,
    timestamp       TIMESTAMPTZ NOT NULL,
    layer           VARCHAR(20) NOT NULL,
    component       VARCHAR(50),
    severity        VARCHAR(20),
    error_code      INTEGER,
    message         TEXT,
    context         JSONB,
    auto_resolved   BOOLEAN DEFAULT FALSE,
    resolution      TEXT,
    rca_cause       TEXT
);

-- 정비 이력
CREATE TABLE maintenance_log (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id       UUID NOT NULL,
    task_type       VARCHAR(30),
    priority        VARCHAR(20),
    scheduled_date  TIMESTAMPTZ,
    completed_date  TIMESTAMPTZ,
    status          VARCHAR(20) DEFAULT 'pending',  -- pending/completed/overdue/cancelled
    description     TEXT,
    user_action     TEXT
);
```

**테스트 기준 (6개)**:
1. `UploadDiagnosticReport` — 유효한 보고서 → 200 OK, device_health_log에 INSERT
2. `GetDeviceHealthHistory` — 30일 조회 → 시계열 데이터 반환
3. `GetFleetHealthDashboard` — 디바이스 100대 → 통계 집계 정확
4. `CreateServiceTicket` — Critical 이슈 → 티켓 생성 + 관리자 알림
5. `GetMaintenanceSchedule` — pending 정비 3건 → 우선순위순 정렬
6. `UploadDiagnosticReport` — device_id 누락 → 400 Bad Request

---

### 8.7 백엔드 회복탄력성 엔진 (v2.4 신규)

**역할**: Go 백엔드 서비스의 장애 전파 차단, 자동 재시도, 격리 및 우아한 성능 저하

```go
// CircuitBreaker — 외부 서비스(Lab API, Cloud ML, Push Gateway) 장애 전파 차단
type CircuitBreaker struct {
    State          CircuitState  // Closed → Open → HalfOpen
    FailureCount   int
    SuccessCount   int
    Threshold      int           // Open 전환 임계값 (기본 5회)
    HalfOpenMax    int           // HalfOpen 시 허용 시도 횟수 (기본 3회)
    Timeout        time.Duration // Open 유지 시간 (기본 30초)
    LastFailure    time.Time
    OnStateChange  func(from, to CircuitState)
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
    // Closed: 정상 실행, 연속 실패 Threshold 도달 시 → Open
    // Open: 즉시 실패 반환, Timeout 경과 후 → HalfOpen
    // HalfOpen: 제한적 시도, 성공 시 → Closed, 실패 시 → Open
}

// RetryPolicy — 지수 백오프 + 지터 기반 재시도
type RetryPolicy struct {
    MaxAttempts    int           // 최대 재시도 (기본 3회)
    InitialDelay   time.Duration // 초기 대기 (기본 100ms)
    MaxDelay       time.Duration // 최대 대기 (기본 10초)
    Multiplier     float64       // 백오프 배율 (기본 2.0)
    JitterFraction float64       // 지터 비율 (기본 0.1)
    RetryableErrors []error      // 재시도 대상 에러 유형
}

func (rp *RetryPolicy) ExecuteWithRetry(ctx context.Context, fn func() error) error {
    // delay = min(InitialDelay * Multiplier^attempt + jitter, MaxDelay)
    // 비멱등(non-idempotent) 요청은 재시도 금지
}

// BulkheadIsolation — 서비스별 고루틴 풀 격리
type BulkheadIsolation struct {
    Pools map[ServiceDomain]*WorkerPool
    // ServiceDomain: Measurement, Sync, Analytics, Diagnostics, Healing
}

type WorkerPool struct {
    MaxConcurrent  int           // 최대 동시 처리 (도메인별 차등)
    QueueSize      int           // 대기열 크기
    Timeout        time.Duration // 작업 타임아웃
    ActiveWorkers  atomic.Int64
    DroppedCount   atomic.Int64
}

// 도메인별 격리 설정
var DefaultBulkheadConfig = map[ServiceDomain]PoolConfig{
    Measurement:  {MaxConcurrent: 50,  QueueSize: 100, Timeout: 30 * time.Second},
    Sync:         {MaxConcurrent: 20,  QueueSize: 50,  Timeout: 60 * time.Second},
    Analytics:    {MaxConcurrent: 10,  QueueSize: 30,  Timeout: 120 * time.Second},
    Diagnostics:  {MaxConcurrent: 15,  QueueSize: 30,  Timeout: 45 * time.Second},
    Healing:      {MaxConcurrent: 30,  QueueSize: 60,  Timeout: 30 * time.Second},
}

// GracefulDegradation — 부하/장애 시 단계적 기능 축소
type GracefulDegradation struct {
    CurrentLevel   DegradationLevel  // Normal → Reduced → Minimal → Emergency
    Thresholds     DegradationThresholds
}

type DegradationLevel int
const (
    Normal    DegradationLevel = iota  // 전체 기능
    Reduced                             // Analytics/추천 비활성화
    Minimal                             // 측정 저장 + 동기화만
    Emergency                           // 측정 저장만 (로컬 큐)
)

func (gd *GracefulDegradation) Evaluate(metrics SystemMetrics) DegradationLevel {
    // CPU > 90% || Memory > 85% || ErrorRate > 5% → Reduced
    // CPU > 95% || ErrorRate > 15% → Minimal
    // DB unreachable || Critical service down → Emergency
}

// HealthCheckEndpoints — 쿠버네티스 호환 헬스체크
type HealthCheckEndpoints struct{}

func (h *HealthCheckEndpoints) Liveness() HealthResponse {
    // 프로세스 생존 확인 (DB 연결 불필요)
}

func (h *HealthCheckEndpoints) Readiness() HealthResponse {
    // DB + Redis + 필수 서비스 연결 확인
}

func (h *HealthCheckEndpoints) Startup() HealthResponse {
    // 초기화 완료 확인 (마이그레이션, 캐시 워밍)
}
```

**gRPC 서비스 정의**:

```protobuf
service ResilienceService {
    rpc GetCircuitBreakerStatus(Empty) returns (CircuitBreakerStatusList);
    rpc GetBulkheadStatus(Empty) returns (BulkheadStatusList);
    rpc GetDegradationLevel(Empty) returns (DegradationLevelResponse);
    rpc ForceCircuitBreakerState(ForceStateRequest) returns (ForceStateResponse);
    rpc GetResilienceMetrics(TimeRangeRequest) returns (ResilienceMetricsResponse);
}
```

**테스트 기준 (6개)**:
1. `CircuitBreaker.Execute()` — 연속 5회 실패 → State=Open, 즉시 실패 반환
2. `CircuitBreaker` — Open 30초 후 → HalfOpen, 성공 1회 → Closed
3. `RetryPolicy.ExecuteWithRetry()` — 2회 실패 후 3회차 성공 → 최종 성공, 지수 백오프 대기시간 검증
4. `BulkheadIsolation` — Measurement 풀 50 포화 → Analytics 요청 정상 처리 (격리 확인)
5. `GracefulDegradation` — ErrorRate 20% → Minimal 레벨, 측정 저장만 허용
6. `HealthCheckEndpoints.Readiness()` — DB 연결 끊김 → 503 + 상세 사유

---

### 8.8 핵심 DB 스키마

```sql
-- TimescaleDB 하이퍼테이블
CREATE TABLE measurements (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    device_id       UUID NOT NULL REFERENCES devices(id),
    cartridge_id    UUID NOT NULL REFERENCES cartridges(id),
    measured_at     TIMESTAMPTZ NOT NULL,
    raw_detection   DOUBLE PRECISION[] NOT NULL,
    raw_reference   DOUBLE PRECISION[] NOT NULL,
    alpha           DOUBLE PRECISION NOT NULL,
    differential    DOUBLE PRECISION[] NOT NULL,
    fingerprint_896 DOUBLE PRECISION[896],
    classification  JSONB,
    concentrations  JSONB,
    anomaly_score   DOUBLE PRECISION,
    model_id        VARCHAR(64) NOT NULL,
    model_hash      VARCHAR(64) NOT NULL,
    twin_predicted  DOUBLE PRECISION[],
    twin_residual   DOUBLE PRECISION[],
    drift_score     DOUBLE PRECISION,
    hash_chain      VARCHAR(64) NOT NULL,
    data_integrity  VARCHAR(64) NOT NULL,
    fhir_observation_id VARCHAR(64)
);
SELECT create_hypertable('measurements', 'measured_at');

-- Milvus 벡터 컬렉션
-- collection: fingerprints
-- fields: id(INT64), measurement_id(VARCHAR), vector(FLOAT_VECTOR(896)), substance_class(VARCHAR)
-- index: IVF_SQ8, nlist=1024, nprobe=16

-- 카트리지 레지스트리
CREATE TABLE cartridge_definitions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    developer_id    UUID NOT NULL REFERENCES developers(id),
    name            VARCHAR(128) NOT NULL,
    version         VARCHAR(16) NOT NULL,
    stage           SMALLINT NOT NULL,
    csi_version     VARCHAR(8) NOT NULL,
    pin_config      JSONB NOT NULL,
    afe_blocks      TEXT[] NOT NULL,
    measurement_modes TEXT[] NOT NULL,
    calibration_schema JSONB NOT NULL,
    regulatory_class VARCHAR(16),
    price_krw       INTEGER,
    revenue_share   DECIMAL(3,2) DEFAULT 0.70,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
```

---

## 9. Layer 6: AI/ML + 인프라

### 9.1 3단계 AI 아키텍처

| Level | 위치 | 기술 | 추론 시간 | 모델 크기 |
|---|---|---|---|---|
| 1 (Edge) | 리더기+앱 | TFLite INT8 XGBoost | <1ms | <2MB |
| 2 (Cloud) | 서버 | Triton + DNN 앙상블 | <100ms | ~500MB |
| 3 (Federated) | 분산 | Flower FedAvg | — | 가중치만 전송 |

### 9.2 PCCP MLOps 모델 레지스트리

```yaml
model_registry:
  model_id: "MPS-XGB-GLUCOSE-v2.3.1"
  model_hash: "sha256:a1b2c3d4..."
  pccp:
    change_description: "훈련 데이터 10,000건 추가"
    change_protocol: "MPS-PCCP-001 §3.2"
    impact_assessment:
      accuracy_change: "+0.3% (95%CI: +0.1~+0.5%)"
      no_regression: true
  performance:
    accuracy: 0.967
    sensitivity: 0.952
    specificity: 0.978
    cv_percent: 1.8
  deployment:
    staged_rollout: [1%, 5%, 25%, 100%]
    rollback_trigger: "accuracy < 0.95 OR cv > 3.0%"
```

### 9.3 환경-건강 교차분석 AI

```python
class CrossDomainAnalyzer:
    """
    건강 지표(CRP, 혈당) + 환경 지표(VOC, PM2.5)의
    시간지연 상관관계 분석. 만파식 고유 강점.
    """
    def analyze_lagged_correlation(
        self,
        env_timeseries: pd.DataFrame,
        health_timeseries: pd.DataFrame,
        max_lag_hours: int = 168,
    ) -> List[CrossDomainInsight]:
        # lag 0~168h, 6h 간격으로 상관계수 계산
        # |corr| > 0.3 and p < 0.05 → 유의미한 인사이트
```

### 9.4 연합학습 (FPCN)

```
메커니즘:
1. 각 리더기 → 로컬 모델 미세 조정
2. Flower 서버 → 가중치 집계 (FedAvg/FedProx)
3. 차분 프라이버시 (ε=1.0) 적용
4. 지역 이상탐지 → 감염병 조기 경보

프라이버시: 개인 측정 데이터는 기기를 떠나지 않음 (GDPR/HIPAA 준수)
```

---

### 9.5 AI 모델 드리프트 감지 (v2.3 신규)

**역할**: 배포된 TFLite/XGBoost 모델의 입력 분포 변화(Data Drift) + 출력 정확도 저하(Concept Drift) 감지

```python
class ModelDriftDetector:
    """
    L6 계층 자가검증 — AI 모델 건강 모니터

    - Data Drift: KL-divergence / PSI (Population Stability Index)
      - 학습 데이터 분포 vs 최근 7일 추론 입력 분포
      - PSI > 0.25 → Warning, > 0.50 → Critical

    - Concept Drift:
      - 예측 confidence 이동 평균 모니터링
      - 사용자 피드백(결과 수정) 비율 모니터링
      - ADWIN (Adaptive Windowing) 알고리즘

    - 자동 대응:
      - Warning → 모니터링 강화 + 관리자 알림
      - Critical → 이전 안정 모델로 자동 롤백 + 재학습 트리거
    """

    def compute_psi(self, reference: np.ndarray, current: np.ndarray) -> float: ...
    def detect_concept_drift(self, predictions: List[Prediction], feedback: List[Feedback]) -> DriftStatus: ...
    def recommend_action(self, drift_status: DriftStatus) -> ModelAction: ...
    def layer_score(self) -> float: ...
```

**테스트 기준 (4개)**:
1. `compute_psi()` — 동일 분포 → PSI < 0.1
2. `compute_psi()` — 분포 이동 → PSI > 0.25, Warning
3. `detect_concept_drift()` — confidence 평균 0.6→0.4 (7일) → DriftStatus::Detected
4. `recommend_action(Critical)` → `ModelAction::Rollback`

---

### 9.6 AI 자가교정 엔진 (v2.4 신규)

**역할**: §9.5 드리프트 감지가 탐지한 이상에 대해 **자동 교정**을 수행하는 능동적 자가치유 계층. 온라인 점진 학습, 피처 드리프트 보상, 모델 앙상블 폴백, 신뢰도 재보정을 통해 AI 파이프라인의 연속적 정확도를 유지한다.

```python
class AiSelfCorrectionEngine:
    """
    L6 계층 자가치유 — AI 모델 자동 교정기

    §9.5 ModelDriftDetector가 Warning/Critical 판정 시 호출.
    4단계 교정 전략을 순차 적용:

    Phase 1: ConfidenceRecalibrator
      - Platt Scaling / Temperature Scaling 재적합
      - 최근 7일 검증 데이터로 보정 곡선 갱신
      - 보정 전후 ECE(Expected Calibration Error) 비교

    Phase 2: FeatureDriftCompensator
      - 입력 피처 분포 이동 감지 → 정규화 파라미터 갱신
      - Z-score / Min-Max 스케일러 온라인 업데이트
      - 차분식 α 계수 적응형 조정: α_new = α_old × (σ_ref / σ_current)

    Phase 3: OnlineIncrementalLearner
      - XGBoost: 기존 모델에 새 데이터 부스팅 라운드 추가
      - TFLite: 마지막 Dense 레이어만 재학습 (Transfer Learning)
      - 학습률: η = 0.001 (보수적), 최대 50 iteration
      - 검증 성능 하락 시 즉시 중단 + 롤백

    Phase 4: ModelEnsembleFallback
      - Primary 모델 실패 시 Secondary 모델로 자동 전환
      - TFLite ↔ XGBoost 교차 검증 불일치 > 15% → 앙상블 합의
      - 가중 투표: w_tflite × P_tflite + w_xgb × P_xgb (가중치는 최근 정확도 기반)
    """

    def __init__(self, drift_detector: ModelDriftDetector):
        self.drift_detector = drift_detector
        self.calibrator = ConfidenceRecalibrator()
        self.feature_compensator = FeatureDriftCompensator()
        self.incremental_learner = OnlineIncrementalLearner()
        self.ensemble_manager = ModelEnsembleFallback()
        self.correction_log: List[CorrectionEvent] = []

    def auto_correct(self, drift_status: DriftStatus) -> CorrectionResult:
        """
        드리프트 심각도에 따라 단계적 교정 실행.

        Returns:
            CorrectionResult with:
              - phase_applied: 적용된 최종 교정 단계 (1~4)
              - metrics_before: 교정 전 정확도/ECE
              - metrics_after: 교정 후 정확도/ECE
              - rollback_needed: 교정 실패 시 True
              - next_review: 다음 재검증 시점 (시간)
        """
        ...

    def recalibrate_confidence(self, predictions: List[Prediction],
                                actuals: List[bool]) -> CalibrationResult:
        """Temperature Scaling 기반 신뢰도 재보정"""
        ...

    def compensate_feature_drift(self, reference_stats: FeatureStats,
                                  current_stats: FeatureStats) -> CompensationResult:
        """피처 정규화 파라미터 온라인 갱신 + α 계수 적응"""
        ...

    def incremental_update(self, new_samples: List[Sample],
                           model_type: ModelType) -> UpdateResult:
        """보수적 온라인 학습 (검증 성능 모니터링 포함)"""
        ...

    def ensemble_fallback(self, primary_pred: Prediction,
                          secondary_pred: Prediction) -> EnsembleResult:
        """교차 검증 불일치 시 가중 투표 기반 합의"""
        ...

    def layer_score(self) -> float:
        """AI 자가교정 건강 점수 (0.0 ~ 1.0)"""
        ...

class CorrectionResult:
    phase_applied: int          # 1~4
    metrics_before: ModelMetrics
    metrics_after: ModelMetrics
    rollback_needed: bool
    correction_strategy: str
    duration_ms: int
    next_review_hours: int      # 교정 후 재검증까지 시간
```

**교정 트리거 매트릭스**:

| 드리프트 상태 | PSI 범위 | 교정 단계 | 최대 허용 시간 |
|---|---|---|---|
| Warning (경미) | 0.25 ~ 0.35 | Phase 1 (재보정) | 즉시 (자동) |
| Warning (중등) | 0.35 ~ 0.50 | Phase 1 + 2 (재보정 + 피처 보상) | 1시간 이내 |
| Critical (경중) | 0.50 ~ 0.70 | Phase 1 + 2 + 3 (+ 점진 학습) | 4시간 이내 |
| Critical (심각) | > 0.70 | Phase 4 (앙상블 폴백) + 관리자 알림 | 즉시 (자동) |
| Concept Drift | confidence 급감 | Phase 3 + 4 + 클라우드 재학습 요청 | 24시간 이내 |

**안전장치**:
- 모든 교정은 검증 데이터셋 대비 성능 비교 후 적용/롤백 결정
- 연속 3회 교정 실패 → 이전 안정 모델로 강제 롤백 + `ServiceTicket` 생성
- 온라인 학습 시 환자 안전 관련 바이오마커(예: Troponin-I 심근경색 판정)는 학습 대상에서 제외 → 클라우드 전문 팀 검증 후에만 모델 갱신
- 교정 이력은 `HashChain`(§6.7)에 기록하여 감사 추적성 보장

**테스트 기준 (6개)**:
1. `recalibrate_confidence()` — ECE 0.15 → 재보정 후 ECE < 0.05
2. `compensate_feature_drift()` — σ 2배 변화 → α 계수 적응 후 차분식 오차 < 3%
3. `incremental_update(XGBoost)` — 신규 50샘플 → 정확도 유지 또는 개선, 기존 검증셋 회귀 없음
4. `incremental_update()` — 학습 중 검증 정확도 하락 → 즉시 중단 + 롤백 확인
5. `ensemble_fallback()` — TFLite 예측 0.8, XGBoost 예측 0.3 (불일치 50%) → 앙상블 합의 + 알림
6. `auto_correct(Critical, PSI=0.65)` — Phase 1~3 순차 적용 → 최종 metrics_after.accuracy ≥ metrics_before.accuracy

---

## 10. SDK & 카트리지 마켓플레이스

### 10.1 SDK 구조

```
manpasik-cartridge-sdk/
├── src/
│   ├── builder.rs        — CartridgeBuilder (빌더 패턴)
│   ├── validator.rs      — CSI 핀 호환성, AFE 지원, 전력 예산, NFC 용량 검증
│   ├── simulator.rs      — 가상 측정/노이즈/드리프트 시뮬레이션
│   ├── regulatory.rs     — FDA/CE-IVDR/MFDS 자동 분류
│   └── publisher.rs      — 마켓플레이스 등록
├── bindings/             — Python(PyPI) + TypeScript(npm) + C(vcpkg)
└── templates/            — 전기화학/면역분석/LAMP/가스 템플릿
```

### 10.2 수익분배 (확정)

```
소비자 결제 → 개발자 70% / 플랫폼 15% / 제조 10% / 규제 5%
(MPS 자체 카트리지는 100% MPS)
```

---

## 11. 보안 아키텍처 (7-Layer)

```
Layer 7: 앱 보안       — RBAC+ABAC, 생체 인증, JWT 회전
Layer 6: API 보안      — Kong Gateway, OAuth 2.0+PKCE, mTLS
Layer 5: 데이터 보안   — AES-256-GCM, TLS 1.3, 해시 체인, PHI/PII 분리
Layer 4: 기기 보안     — TPM 2.0, Secure Boot, NFC SUN/CMAC, ECDSA P-256
Layer 3: 네트워크 보안 — BLE AES-128, WPA3, VPN
Layer 2: 인프라 보안   — WAF, K8s RBAC, 컨테이너 서명
Layer 1: PQC 준비     — ML-KEM + X25519 하이브리드 대비
```

### 11.1 의료기기 규제 프로세스 매핑 (v2.4 보강)

**IEC 62304 소프트웨어 안전 분류**:

| 소프트웨어 시스템 | 안전 분류 | 근거 |
|---|---|---|
| Rust Core (측정·차분식·AI) | **Class C** | 잘못된 결과 → 환자 사망/심각한 부상 가능 (심근경색·당뇨 오판) |
| Flutter 앱 (결과 표시) | **Class B** | 표시 오류 → 부적절한 치료 결정 가능 (경미한 부상) |
| Go 백엔드 (동기화·분석) | **Class B** | 데이터 손실 → 추적성 상실, 간접적 환자 위해 |
| Self-Healing (§6.16) | **Class C** | 잘못된 자가치유 → 측정 정확도 직접 훼손 |
| SDK/마켓플레이스 | **Class A** | 개발 도구, 환자 직접 영향 없음 |

**IEC 62304 프로세스 매핑**:

```
SW 요구사항 분석 (§4 인터페이스 정의 + §6 모듈 사양)
  → SW 아키텍처 설계 (§3 6-Layer + §2 디렉토리 구조)
    → SW 상세 설계 (§6.1~§6.16 모듈별 인터페이스 계약)
      → SW 구현 (Appendix B 프롬프트 기반 코드 생성)
        → SW 단위 검증 (§12.1 계층별 테스트)
          → SW 통합 검증 (§12.2 시나리오 테스트)
            → SW 시스템 검증 (§12.3 SSOT 자동 검증)
```

**ISO 14971 위험관리 핵심 항목**:

| 위험 시나리오 | 심각도 | 발생 확률 | 위험 수준 | 통제 조치 | 잔여 위험 |
|---|---|---|---|---|---|
| 혈당 위험 수준 오탐(위음성) | 치명적 | 낮음(AI 3축 92%+) | High | 차분식+디지털트윈+AI 교차검증, 3σ 이상치 거부 | 수용 가능 |
| BLE 데이터 변조 | 심각 | 매우 낮음 | Medium | AES-128 암호화 + 해시 체인(§6.7) + Sequence 번호 검증 | 수용 가능 |
| 자가치유 중 잘못된 α 교정 | 심각 | 낮음 | High | α 클램핑 [0.90,1.10] + 교차검증 + 연속 3회 실패 시 롤백 | 수용 가능 |
| 카트리지 만료 미감지 | 중등 | 낮음 | Medium | NFC 매니페스트 만료 필드 + 서버 교차 검증 | 수용 가능 |
| AI 온라인 학습 데이터 오염 | 심각 | 중간 | High | 품질 게이트(§9.6), 환자 안전 바이오마커 학습 제외, 롤백 | 수용 가능 |
| 펌웨어 업데이트 중 벽돌화 | 심각 | 낮음 | Medium | A/B 파티션 + 자동 롤백(§6.16.2) | 수용 가능 |

**FDA Cybersecurity Guidance 요구사항 매핑**:

| FDA 요구사항 | 구현 위치 | 상태 |
|---|---|---|
| Asset inventory | §5.3 BLE GATT + §2 디렉토리 | 구현됨 |
| Threat modeling (STRIDE) | §11 7-Layer + ISO 14971 | [보강 필요] — STRIDE 모델 별도 문서 |
| Authentication | Layer 7 RBAC + JWT + 생체 | 구현됨 |
| Encryption (at rest/transit) | Layer 5 AES-256-GCM + TLS 1.3 | 구현됨 |
| Software Bill of Materials (SBOM) | [미구현] | Phase 2 SBOM 자동 생성(cargo-sbom + syft) |
| Vulnerability management | [미구현] | Phase 2 cargo-audit + npm-audit CI 통합 |
| Incident response | §6.16 SelfHealing + §6.15.5 ErrorReporter | 부분 구현 |
| OTA update integrity | §5.2 A/B 파티션 + ECDSA 서명 | 구현됨 |

**제약**: 자가치유 중 측정 안전성 — 측정 진행 중(M-040 활성 상태) Level 2 이상 에스컬레이션(PeripheralReset, SoftReset, FactoryFallback)은 **금지**한다. 측정 완료 후에만 실행 가능하며, 측정 중에는 Level 0~1(TaskRestart, StackReset)만 허용한다.

### 11.2 의료 데이터 표준 연동 (v2.4 보강)

**FHIR R4 리소스 매핑**:

| ManPaSik 개체 | FHIR 리소스 | 주요 필드 매핑 |
|---|---|---|
| MeasurementResult | `Observation` | `code`(LOINC), `valueQuantity`, `effectiveDateTime`, `device` |
| CheckupSession | `DiagnosticReport` | `result[]`(Observation 참조), `conclusionCode`, `status` |
| DiseaseRisk | `RiskAssessment` | `prediction[].probabilityDecimal`, `condition`, `basis[]` |
| CartridgeManifest | `Device` | `type`, `version`, `expirationDate`, `lotNumber` |
| UserProfile | `Patient` | `identifier`, `birthDate`, `gender` (PHI 분리 원칙 준수) |

**HL7 CDA 지원**: Phase 4에서 `DiagnosticReport` → CDA XML 자동 변환 서비스(Go gRPC) 추가 예정. 병원 EHR 연동 시 필수.

### 11.3 접근성 표준 (WCAG 2.1 AA) (v2.4 보강)

| 기준 | 요구사항 | 구현 전략 |
|---|---|---|
| 1.1.1 비텍스트 콘텐츠 | 모든 차트·그래프에 대체 텍스트 | `Semantics(label:)` 위젯 감싸기 |
| 1.4.3 대비 (최소) | 텍스트 대비 비율 4.5:1 이상 | MPS 디자인 토큰에 대비 검증 |
| 1.4.11 비텍스트 대비 | UI 컴포넌트 대비 3:1 이상 | 아이콘·버튼 색상 검증 |
| 2.1.1 키보드 | 모든 기능 키보드 접근 가능 | Flutter FocusNode + tab order |
| 2.5.5 타겟 크기 | 최소 44×44 dp | 모든 버튼·탭 타겟 규격화 |
| 3.3.1 오류 식별 | 입력 오류 시 텍스트로 안내 | FormField validator + 음성 안내 |
| 4.1.2 이름·역할·값 | 보조 기술이 UI 상태 파악 | Semantics + ExcludeSemantics 적용 |

**고령자 모드**: 설정에서 "큰 글씨 모드" 활성화 시 기본 폰트 14pt→18pt, 버튼 높이 48dp→56dp, 터치 타겟 48dp→56dp. 측정 결과 화면에 음성 안내(TTS) 자동 재생 옵션.

---

## 12. 테스트 전략

### 12.1 계층별 테스트

| 레이어 | 테스트 유형 | 도구 | 커버리지 목표 |
|---|---|---|---|
| Rust Core | 단위 테스트 | `cargo test` | 80% |
| Rust Core | 속성 기반 | `proptest` | 핵심 알고리즘 |
| Flutter | 단위 + 위젯 | `flutter_test` | 80% |
| Flutter | 통합 | `integration_test` | 핵심 플로우 |
| Go Backend | 단위 + API | `go test` + `grpcurl` | 80% |
| E2E | 시나리오 | 자동화 스크립트 | 4개 시나리오 |

### 12.2 통합 테스트 시나리오

| # | 시나리오 | 검증 항목 |
|---|---|---|
| 1 | 혈당 측정 E2E | NFC→BLE→차동측정→AI→결과표시→클라우드 저장 |
| 2 | 오프라인 측정+동기화 | 오프라인 측정→로컬 저장→온라인 복귀→CRDT 병합 |
| 3 | 카트리지 만료 거부 | 만료 카트리지 삽입→경고 표시→측정 차단 |
| 4 | OTA 업데이트+롤백 | 업데이트 시작→검증 실패→자동 롤백→이전 버전 정상 동작 |
| 5 | Standard 검진 E2E (v2.1) | 패키지 선택 → 3카트리지 순차 NFC → 3회 측정 → AI 리스크 산출 → 리포트 표시 |
| 6 | 검진 중간이탈+복원 (v2.1) | 2/3 완료 → 앱 종료 → 재시작 → resume → 3/3 완료 |
| 7 | 유기적 순환 E2E (v2.2) | 혈당 측정 → 영양 부족 판정 → 제품 추천 → 구매 → 복용 기록 → 2주 후 재측정 → 효과 리포트 |
| 8 | 컨텍스트 카드 생성 (v2.2) | 7개 데이터소스 → ContextEngine → 10장 카드 → 우선순위 정렬 → Danger 최우선 |
| 9 | 자가검증→측정 차단/허용 (v2.3) | preFlightCheck→HW 이상 감지→M-085 차단 화면→사용자 재검증→통과→정상 측정 진행→결과에 ReliabilityScore 배지 |
| 10 | 자가치유+세션 복원 E2E (v2.4) | 측정 중 BLE 단절→SelfHealingOrchestrator 자동 감지→BleReconnect 시도→OV-083 토스트→복원 성공→측정 계속. 이후 앱 크래시→재시작→OV-084 세션 복원 다이얼로그→체크포인트 복원→측정 완료 |
| 11 | 자가치유 실패→에스컬레이션 (v2.4) | ADC 포화 지속→HwAutoRecovery 3회 시도 실패→M-087 측정 중단 화면→서비스 티켓 자동 생성→힐링 로그(H-038) 자동 첨부 확인 |

### 12.3 SSOT 정합성 자동 검증

```bash
#!/bin/bash
# ssot_check.sh — CI에서 자동 실행
# Rust 상수가 SSOT와 일치하는지 검증

grep -r "ALPHA_DEFAULT" rust-core/ | grep -v "0.98" && echo "FAIL: ALPHA_DEFAULT" && exit 1
grep -r "CONNECTOR_PINS" rust-core/ | grep -v "16" && echo "FAIL: CONNECTOR_PINS" && exit 1
grep -r "FINGERPRINT_DIM_FULL" rust-core/ | grep -v "896" && echo "FAIL: FINGERPRINT_DIM" && exit 1
grep -rn "E12-IF\|12핀\|12-pin" rust-core/ flutter-app/ go-backend/ && echo "FAIL: E12-IF 참조" && exit 1
echo "SSOT CHECK PASSED"
```

---

## 13. 구현 로드맵 + Phase별 체크리스트

### Phase 0: 기반 구축 (Month 1-3)

```
□ Cargo workspace 초기화 (manpasik-core, manpasik-ffi, manpasik-fw)
□ SensorTrait HAL 정의 (§4.1)
□ CartridgeManifest v2.0 파서 (§4.2)
□ AfeRegistry (§4.3)
□ DifferentialEngine 구현 + 10개 테스트 통과 (§6.1) ★핵심
□ DSP 파이프라인 (§6.2)
□ FeatureExtractor 88차원 (§6.3)
□ Flutter 프로젝트 초기화 (Riverpod + freezed)
□ Go gRPC 백엔드 초기화 (proto + 기본 서비스)
□ SSOT 자동 검증 스크립트 (§12.3)
□ CI/CD 파이프라인 (GitHub Actions)
```

**Gate G0 통과 기준**: DifferentialEngine 10개 테스트 100% 통과, SensorTrait 컴파일, CartridgeManifest v1.0/v2.0 파싱 성공

### Phase 1: MVP (Month 4-8)

```
□ BLE 통신 서비스 (§7.2) — 리더기 연동
□ NFC 카트리지 인식 (§7.3)
□ 측정 화면 (§7.1 상태 관리)
□ 로컬 저장 + CRDT 기본 (§7.4)
□ 해시 체인 보안 (§6.7)
□ Stage-1 전기화학 카트리지 3종 (혈당/전해질/CRP)
□ 테스트 커버리지 80%
□ 통합 테스트 시나리오 1-2 (§12.2)
```

**Gate G1**: BLE 측정 성공률>95%, 차동측정 CV<3%, 내부 알파 NPS>40

### Phase 2: 핵심 기능 (Month 9-14)

```
□ TFLite AI 추론 엔진 (§6.5) ★특허 핵심
□ 디지털 트윈 v1.0 (§6.4) ★특허 핵심
□ 종합검진 Quick 기본 (§6.8 CheckupSession + §7.5 UI) ★v2.1 신규
□ AI 자동분류 기본 (§6.10 AutoClassifier)
□ ContextEngine 기본 (§6.11) ★v2.2 신규
□ 컨텍스트 카드 피드 + 사용자 반응 학습 (§7.6) ★v2.2 신규
□ SDK CartridgeBuilder (§10.1)
□ CRDT 완전 오프라인 (§7.4)
□ AI 코치 v1.0
□ 데이터 허브 + 리포트
□ SaaS 구독 3티어 + 결제
□ 통합 테스트 시나리오 3-6 (§12.2) + 시나리오 7-8 (v2.2)
```

**Gate G2**: AI 정확도>92%, 오프라인 100%, SDK 빌드 가능, Quick 검진 E2E 통과, ContextCard 우선순위 정렬 검증

### Phase 3-5: 확장 (Month 15-30)

```
□ Standard/Premium 검진 확장 (§6.8, §6.9, §7.5) ★v2.1 신규
□ 질병 리스크 엔진 완전 구현 (§6.9 DiseaseRiskEngine)
□ 검진 이력 비교 + 추세 분석 (§8.4)
□ Custom 맞춤검진 + 검진 패키지 스토어 (§10)
□ NutritionAdvisor 완전 구현 (§6.12) ★v2.2 신규
□ ShoppingBridge + 정기배송 (§6.13) ★v2.2 신규
□ HabitTracker + 게이미피케이션 (§6.14) ★v2.2 신규
□ 유기적 연동 gRPC 서비스 (§8.5) ★v2.2 신규
□ 제품 효과 추적 + 영양-바이오마커 Deep 인사이트
□ 포인트/리워드 시스템 완성
□ HwAutoRecovery + SessionRecovery 구현 (§6.16.1, §6.16.4) ★v2.4 신규 Phase 3 우선
□ FwSelfRepair + DataPipelineReprocessor 구현 (§6.16.2, §6.16.3) ★v2.4 신규 Phase 4
□ Flutter 복원 UI (§7.8) — H-038, M-087, OV-083~084 ★v2.4 신규
□ Go 백엔드 회복탄력성 (§8.7) — CircuitBreaker, BulkheadIsolation ★v2.4 신규
□ AI 자가교정 엔진 (§9.6) — Phase 5, 온라인 학습 포함 ★v2.4 신규
□ 통합 테스트 시나리오 9~11 (§12.2) — 자가검증·자가치유 E2E ★v2.4 신규
□ SiPM-ECL 포화보정 (§6.6) — Stage-2 준비
□ 연합학습 Flower (§9.4)
□ 환경-건강 교차분석 (§9.3)
□ OTA 서비스 (A/B 파티션, 델타 업데이트)
□ Matter 스마트홈 연동
□ AR 측정 가이드, 음성 제로터치
□ 규제 인증 (MFDS, FDA 510(k), CE-IVDR)
□ 글로벌 출시
```

### 비용 요약

```
총 예산: ₩85억 (30개월)
├── 인건비:        ₩54억 (63%, 피크 28명)
├── 인프라:        ₩12억 (14%)
├── 규제:          ₩10억 (12%)
├── 외주/라이선스: ₩5억 (6%)
└── 예비비:        ₩4억 (5%)
```

---

## Appendix A: IDE CLAUDE.md (프로젝트 루트 배치용)

> 아래 내용을 `manpasik/CLAUDE.md`로 저장하면 IDE AI가 자동 참조한다.

```markdown
# ManPaSik (萬波息) — 프로젝트 CLAUDE.md

## 프로젝트 정의
차동측정 기반 범용분석 POCT 시스템의 소프트웨어 스택.
하드웨어 리더기(STM32F405) + 일회용 카트리지를 연동하는
모바일 앱(Flutter) + Rust 핵심 엔진 + Go 백엔드 + AI/ML 파이프라인.

## 참조 문서
- 기술 사양서: docs/ManPaSik_Technical_Specification_v2.0.md ← 이 문서
- 마스터플랜: docs/ManPaSik_System_Master_Plan_v2.0.md

## SSOT 베이스라인 (절대 변경 금지)
- 커넥터: CSI v1.0 = Samtec MECF-08-01-L-DV, 16핀, 1.27mm
- 차분식: S_diff = S_detection - alpha * S_reference (alpha 기본 0.98, 범위 0.90~1.10)
- PPM: 49.70×30×4.30mm
- Universal AFE: 9블록 (Stage별 활성화)
- 정확도: 3축 결합 92~98%
- 핑거프린트: 88→448→896→1792차원 (Stage별)

## 기술 스택
- Rust: 핵심 엔진 (no_std, embedded-hal)
- Flutter 3.x: 모바일 앱 (Riverpod + freezed + flutter_rust_bridge)
- Go 1.22+: 백엔드 (gRPC)
- PostgreSQL 16 + TimescaleDB + Redis + Milvus
- TFLite: 엣지 AI (XGBoost INT8)

## 코딩 규칙
1. Rust: Result<T,E> 필수. unwrap() 금지(테스트 제외). clippy 0건.
2. Flutter: Riverpod 2.x + freezed. setState() 금지.
3. Go: protobuf gRPC. REST는 gRPC-Gateway.
4. 테스트: 모듈당 최소 5개. 커버리지 80%.
5. 문서: pub 함수에 rustdoc/dartdoc 필수.
6. 보안: 하드코딩 키 금지. AES-256-GCM + SHA-256.

## 하네스 원칙
- H1 모듈 독립: trait/interface로만 의존. 구체 타입 직접 참조 금지.
- H2 후방 호환: CartridgeManifest v1.0은 v2.0 코드에서 반드시 동작.
- H5 인터페이스 계약: BLE GATT UUID, NFC 매니페스트 버전화.
- H6 실패 격리: 센서 1개 오류가 전체 측정을 중단시키지 않음.

## 금지 사항
- E12-IF 12핀 참조 금지 (CSI v1.0 16핀으로 변경 확정)
- AI 생성 가상 검증 결과 삽입 금지
- 미검증 수치를 "검증됨"으로 표기 금지
```

---

## Appendix B: 참조 프롬프트 모음

> 각 프롬프트는 사양서의 해당 섹션을 참조하여 구현을 트리거한다.
> IDE에서 "§N.N 사양서 참조" 부분을 AI가 자동으로 사양서에서 찾아 읽는다.

### Phase 0 기반구축

```
B-00: Cargo workspace 초기화
"§2 프로젝트 구조대로 manpasik/ Cargo workspace를 생성해.
 members: manpasik-core, manpasik-ffi, manpasik-fw.
 manpasik-core/src/lib.rs에 harness, signal, calibration, quantification,
 fingerprint, ai, digital_twin, communication, security, data, cartridge,
 sipm_ecl 모듈을 mod 선언해."
```

```
B-01: SensorTrait HAL 정의
"§4.1 SensorTrait 인터페이스 계약을 그대로 구현해.
 sensor_trait.rs에 trait + 9블록 구현체 선언 + ElectrochemAfe 타입 정의.
 테스트 6개 작성."
```

```
B-02: CartridgeManifest v2.0
"§4.2 CartridgeManifest 구조체를 구현해.
 256바이트 파싱 + v1.0 후방호환(from_v1) + CRC32 검증 + 만료/사용횟수 체크.
 AFE 비트마스크 해석. 테스트 7개."
```

```
B-03: AfeRegistry
"§4.3 AfeRegistry를 구현해. HashMap 기반 동적 등록, 매니페스트 기반 자동 선택.
 테스트 3개."
```

```
B-04: 차동측정 엔진 ★핵심
"§6.1 DifferentialEngine을 정확히 구현해.
 compute, compute_array, update_alpha(환경보정+AI가중치), snr, matrix_removal_rate.
 alpha 클램핑 [0.90, 1.10]. 테스트 10개 전부 통과시켜."
```

```
B-05: DSP 파이프라인
"§6.2 DspPipeline을 구현해. Butterworth BPF + Savitzky-Golay + 칼만 필터.
 FftAnalyzer(Hanning 윈도우 실수 FFT). PeakDetector(SWV/DPV).
 테스트 4개."
```

```
B-06: 핑거프린트 추출기
"§6.3 FeatureExtractor를 구현해.
 88차원 기본 추출 + 전자코 융합 896차원 확장 + L2 정규화.
 테스트 4개."
```

```
B-07: Flutter 프로젝트 초기화
"§7 사양대로 Flutter 프로젝트를 생성해.
 pubspec.yaml: riverpod, freezed, flutter_rust_bridge, flutter_reactive_ble, nfc_manager, fl_chart.
 §7.1 MeasurementState freezed 클래스 + MeasurementNotifier 작성.
 디렉토리: providers/, screens/, services/, models/, widgets/"
```

```
B-08: Go 백엔드 초기화
"§8 사양대로 Go gRPC 프로젝트를 생성해.
 §8.3 Protobuf 정의 (measurement.proto, cartridge.proto, ai.proto).
 §8.4 measurements 테이블 + cartridge_definitions 테이블 DDL.
 기본 서비스 stub 3개 (measurement-svc, cartridge-svc, ai-svc)."
```

```
B-09: Rust 종합검진 모듈 초기화 (v2.1)
"§6.8~6.10 종합검진 모듈을 manpasik-core에 생성해.
 src/checkup/ 디렉토리 생성.
 session.rs, disease_risk.rs, auto_classifier.rs 기본 구조 + 모듈 등록.
 Cargo.toml에 필요한 의존성 추가."
```

```
B-10: Flutter 종합검진 화면 초기화 (v2.1)
"§7.5 종합검진 화면을 Flutter 프로젝트에 생성해.
 screens/checkup/ 디렉토리 생성.
 checkup_screen.dart, checkup_notifier.dart 기본 구조.
 Riverpod providers 등록 (checkupNotifierProvider, diseaseRiskProvider 등)."
```

### Phase 1 MVP

```
C-01: BLE 통신 서비스
"§7.2 BleService를 구현해. flutter_reactive_ble 사용.
 scanForReaders + connect + waveformStream(0xFF02) + writeConfig(0xFF01).
 자동 재연결. 테스트 3개."
```

```
C-02: NFC 카트리지 인식
"§7.3 NfcService를 구현해. nfc_manager 사용.
 readCartridge → CartridgeManifest 파싱 (Rust FFI 호출).
 v1.0/v2.0 후방호환. 만료/횟수 검증. 테스트 4개."
```

```
C-03: 측정 화면
"§7.1 + §7.5 화면 구조의 '측정' 탭을 구현해.
 NFC 스캔 → BLE 연결 → 실시간 파형(fl_chart) → 결과 4단계 표시.
 MeasurementNotifier 연동. Riverpod Consumer 패턴."
```

```
C-04: CRDT 동기화
"§7.4 CrdtSyncManager를 구현해.
 LWW-Register(측정), G-Counter(사용횟수), OR-Set(카트리지목록).
 오프라인 저장 → 온라인 자동 병합. 테스트 4개."
```

```
C-05: 종합검진 UI (v2.1 신규)
"§7.5 CheckupNotifier + CheckupScreen을 구현해.
 패키지 선택(Quick/Standard/Premium) → 다중 카트리지 순차 스캔 → 결과 통합 표시.
 검진 이력 목록(compare) + 리스크 스코어 4단계 표시. 테스트 4개."
```

### Phase 2-3 핵심

```
D-01: TFLite AI 추론
"§6.5 TfLiteEngine을 구현해.
 모델 로드(hash 검증) + predict(896차원→InferenceResult) + <5ms.
 NonTargetDetector(Mahalanobis). AbTestEngine(병렬추론). 테스트 5개."
```

```
D-02: 디지털 트윈 ★특허 핵심
"§6.4 전체를 구현해.
 DigitalTwinEngine(update→TwinUpdate), DriftDetector(3σ/Allan/연속5회),
 CalibrationPredictor(잔차 회귀), SelfOptimizer(AdjustAlpha).
 테스트 6개."
```

```
D-03: SDK 카트리지 빌더
"§10.1 SDK 구조의 builder.rs + validator.rs를 구현해.
 CartridgeBuilder(빌더 패턴) → CartridgeDefinition.
 validate_pin_config, validate_afe_support, validate_nfc_size."
```

```
D-04: SiPM-ECL 포화보정
"§6.6 SaturationDetector + correct_saturation을 구현해.
 3단계 포화 레벨 판정. 4PL 시그모이드 보정.
 테스트 5개."
```

```
D-05: 종합검진 세션 관리 (v2.1 신규)
"§6.8 CheckupSession을 구현해.
 4개 패키지 정의(Quick/Standard/Premium/Custom).
 순차 측정 관리 + 부분 저장/복원. 테스트 8개."
```

```
D-06: 질병 리스크 엔진 (v2.1 신규)
"§6.9 DiseaseRiskEngine을 구현해.
 8대 영역 리스크 점수(0~100) + 이전 비교 + 관리 권고.
 contribution_factors SHAP 유사 기여도. 테스트 7개."
```

```
D-07: AI 자동분류 엔진 (v2.1 신규)
"§6.10 AutoClassifier를 구현해.
 cartridge_type → 도메인 분류 + 스마트 태그 + 카테고리 통계.
 테스트 5개."
```

```
D-08: 종합검진 UI (v2.1 신규)
"§7.5 CheckupNotifier + CheckupScreen을 구현해.
 패키지 선택 → 다중 카트리지 순차 → 결과 종합.
 검진 이력 목록 + 리스크 스코어 표시."
```

```
D-09: 종합검진 백엔드 (v2.1 신규)
"§8.4 CheckupService gRPC + REST.
 세션 CRUD + 리스크 산출 + 이력 비교.
 measurements/disease_risk_scores/auto_classifications 테이블 DDL."
```

```
B-11: ContextEngine 구현 (v2.2 신규)
"§6.11 ContextEngine을 rust-core/manpasik-core/src/context/engine.rs에 구현해.
 build_context(7개 데이터 소스 병합, 30초 내, graceful degradation).
 generate_cards(최대 10장, priority 내림차순, ContextCardType 중복 최대 2장, Danger ≥90).
 learn_from_interaction(Dismissed 3회 → 우선순위 20% 하향).
 테스트 8개 전부 통과시켜."
```

```
B-12: NutritionAdvisor 구현 (v2.2 신규)
"§6.12 NutritionAdvisor를 rust-core/manpasik-core/src/context/nutrition.rs에 구현해.
 analyze_from_biomarkers(혈당→크롬/마그네슘/식이섬유, 페리틴→철분, 25(OH)D→비타민D, CRP→오메가3/커큐민).
 recommend_sources(약물 상호작용 체크, contraindicated 필드 기록).
 evaluate_supplement_effect(최소 2주 데이터 필요, 미달 시 InsufficientData).
 테스트 7개."
```

```
B-13: ShoppingBridge 구현 (v2.2 신규)
"§6.13 ShoppingBridge를 rust-core/manpasik-core/src/context/shopping.rs에 구현해.
 match_products(match_score = 0.4×영양매칭 + 0.2×가격대비가치 + 0.2×인증점수 + 0.1×리뷰평점 + 0.1×개인이력).
 calculate_repurchase_timing(일일 복용량 × 구매 수량 → 남은 일수, 5일 전 알림).
 match_reason 최소 1개 (근거 투명성).
 테스트 6개."
```

```
B-14: HabitTracker 구현 (v2.2 신규)
"§6.14 HabitTracker를 rust-core/manpasik-core/src/context/habit.rs에 구현해.
 log_habit(스트릭 자동 계산, 중복 입력은 마지막 값으로 덮어쓰기).
 analyze_correlation(최소 30일, 최소 4회 측정, 미달 시 confidence=Low).
 generate_nudge(하루 최대 5회, Morning은 습관/복용, Evening은 수면/회고).
 check_achievements(7일/30일/90일/365일 → 브론즈/실버/골드/플래티넘 배지).
 테스트 6개."
```

```
C-06: 유기적 연동 UI 상태관리 (v2.2 신규)
"§7.6 컨텍스트 카드 피드 + 영양대시보드 + 습관트래커 + 추천스토어 UI를 Flutter에 구현해.
 @riverpod Provider: contextCardFeed, NutritionNotifier, HabitNotifier.
 recommendedProducts, repurchaseAlerts, SubscriptionDeliveryNotifier, RewardNotifier.
 ContextCardFeed 30초 자동 새로고침, NutritionNotifier 식단/보충제 입력 시 자동 리로드.
 Riverpod freezed 패턴 준수."
```

```
D-10: 유기적 연동 백엔드 (v2.2 신규)
"§8.5 ContextService + NutritionService + ShoppingService + HabitService + RewardService gRPC 구현.
 Protobuf: UserContext, ContextCard, NutritionReport, ProductRecommendation, SubscriptionDelivery, RewardStatus.
 REST Gateway: Fiber로 gRPC→JSON 변환 (POST /context/v1/build, GET /nutrition/v1/analyze 등).
 DB 마이그레이션: context_cards, meal_logs, supplement_logs, product_catalog, subscriptions_delivery, habit_logs, achievements, reward_points 테이블 생성.
 테스트 8개."
```

```
B-15: SelfDiagnostics 6계층 자가검증 (v2.3 신규)
"§6.15 SystemHealthOrchestrator를 구현해.
 HwHealthMonitor(L1, §6.15.1) + FwWatchdogBridge(L2) + DataIntegrityChecker(L3) +
 PredictiveMaintenanceEngine(L4) + ErrorReporter(L5) + SystemHealthOrchestrator(L6).
 가중합산: HW 30% FW 15% Rust 25% App 10% Backend 10% AI 10%.
 pre_flight_check 2초 이내. 테스트 26개."
```

```
B-16: SelfHealingOrchestrator (v2.4 신규)
"§6.16 SelfHealingOrchestrator를 구현해.
 HwAutoRecovery(BLE 재연결, ADC 리샘플, AFE 블록 스위치, 보정 클라우드 싱크).
 FwSelfRepair(5단계 에스컬레이션: TaskRestart→StackReset→PeripheralReset→SoftReset→FactoryFallback).
 DataPipelineReprocessor(체크섬 재처리, AI 교차검증, CRDT 충돌 해소).
 MeasurementSessionRecovery(100ms 체크포인트, 5분 복원 윈도우).
 SelfHealingOrchestrator(Detect→Classify→Attempt→Verify→Log 5-phase 루프).
 테스트 30개."
```

```
C-08: Flutter Self-Healing 복원 UI (v2.4 신규)
"§7.8 AppCrashRecoveryNotifier + HealingLogNotifier + HealingStatsNotifier 구현.
 MpsSessionRestoreDialog(복원/폐기 선택), MpsHealingLogTimeline(24h 이벤트),
 MpsHealingStatsCard(성공률, 평균 복구시간, 레이어별 통계).
 앱 시작 시 복원 가능 세션 자동 탐지 → 다이얼로그 표시.
 Riverpod freezed 패턴 준수. 테스트 5개."
```

```
D-12: AI 자가교정 엔진 (v2.4 신규)
"§9.6 AiSelfCorrectionEngine을 구현해.
 Phase 1: ConfidenceRecalibrator(Temperature Scaling, ECE < 0.05).
 Phase 2: FeatureDriftCompensator(α 계수 적응, σ_ref/σ_current 비율 기반).
 Phase 3: OnlineIncrementalLearner(XGBoost 부스팅 라운드 추가, TFLite Dense 재학습, η=0.001).
 Phase 4: ModelEnsembleFallback(TFLite↔XGBoost 가중 투표).
 안전장치: 연속 3회 실패 → 강제 롤백, 환자 안전 바이오마커 학습 제외.
 테스트 6개."
```

### Phase 3-5 확장

```
E-01: 연합학습 — "§9.4 Flower FedAvg 클라이언트 구현. 차분 프라이버시 ε=1.0."
E-02: 교차분석 AI — "§9.3 CrossDomainAnalyzer. 시간지연 상관관계, |corr|>0.3, p<0.05."
E-03: OTA — "§5.2 OTA_Updater. A/B 파티션 + 델타(bsdiff) + 자동 롤백."
E-04: Matter — "스마트홈 연동. Matter over Thread. 환경 데이터 → HomeKit/Google Home."
```

### 검증

```
F-01: 커버리지 검증 — "cargo tarpaulin + flutter test --coverage. 80% 이상 확인."
F-02: SSOT 자동 검증 — "§12.3 ssot_check.sh 실행. 모든 상수 일치 + E12-IF 참조 0건."
F-03: 통합 테스트 — "§12.2 시나리오 1-4 전체 실행."
G-01: CI/CD — "GitHub Actions: lint → test → coverage → SSOT check → build → deploy(ArgoCD)."
```

---

## 14. 향후 로드맵 앵커 (v2.4 신설)

> 본 섹션은 정식 구현 사양이 아닌 **"반드시 다뤄야 할 체크리스트"**이다.
> 각 앵커는 구현 시점에 정식 섹션(§N.N)으로 승격되며, 승격 전까지 `[앵커]` 상태를 유지한다.
> CLAUDE.md §4 진화 원칙에 따라, 입증된 필요성이 확인될 때만 승격한다.

### 14.1 [필수] STRIDE 위협 모델링 별도 문서

| 항목 | 내용 |
|---|---|
| **필요성** | FDA Cybersecurity Guidance 필수 요구사항. §11.1에 매핑만 기재, 상세 위협 분석 미실시 |
| **Phase** | Phase 2 (Gate G3 이전 완료) |
| **선행 조건** | §11 보안 아키텍처 확정, BLE GATT v2.0 확정, FHIR 연동 범위 확정 |
| **완료 기준** | STRIDE 6개 범주(Spoofing, Tampering, Repudiation, Information Disclosure, DoS, Elevation) × 주요 자산(BLE 채널, gRPC API, NFC 매니페스트, OTA 채널, PHI 저장소) 매트릭스 완성. 각 위협에 대한 통제 조치 + 잔여 위험 수용 판정. 별도 문서 `MPS-SEC-STRIDE-v1.0` |
| **산출물** | 위협 모델 문서 (독립), 기술사양서 §11에 참조 링크 추가 |

### 14.2 [필수] SBOM 자동 생성 CI 통합

| 항목 | 내용 |
|---|---|
| **필요성** | FDA SBOM 의무화(2023~), EU Cyber Resilience Act 요구. 현재 미구현 |
| **Phase** | Phase 2 (CI/CD 파이프라인 구축 시 동시) |
| **선행 조건** | Cargo workspace 확정, Flutter pubspec 확정, Go mod 확정 |
| **완료 기준** | `cargo-sbom`(Rust) + `syft`(컨테이너) + `cyclonedx-npm`(Flutter) 자동 실행. SPDX 또는 CycloneDX 형식 출력. CI 빌드마다 갱신. 알려진 CVE 0건 (cargo-audit + npm-audit 통과) |
| **산출물** | CI 스크립트 (§13 로드맵 G-01에 통합), SBOM 파일 (빌드 아티팩트) |

### 14.3 [필수] 다국어 i18n/l10n 아키텍처

| 항목 | 내용 |
|---|---|
| **필요성** | EU IVDR 다국어 라벨링 요구, 글로벌 시장 진입(미국/EU/한국) 필수 |
| **Phase** | Phase 2 시작 시 기반 구축, Phase 3 영어/중국어 추가 |
| **선행 조건** | Flutter 프로젝트 초기화 완료 (B-07), UI 텍스트 확정 |
| **완료 기준** | Flutter `intl` + ARB 파일 기반 번역 시스템. 모든 사용자 표시 문자열 하드코딩 0건. 한국어(기본) + 영어 + 중국어 간체 3개 로캘. 의료 용어 번역 검증(의료 전문가 리뷰). 측정 단위(mg/dL ↔ mmol/L) 자동 변환 |
| **산출물** | §7에 i18n 아키텍처 서브섹션 승격, ARB 파일 구조 정의 |

### 14.4 [강력 권고] Federated Learning 구현 상세

| 항목 | 내용 |
|---|---|
| **필요성** | §9.4에 1문단 수준 기술. 개인 맞춤 AI 차별화 + PII 미전송 규제 이점 |
| **Phase** | Phase 5 |
| **선행 조건** | §9.6 AI 자가교정 엔진 구현 완료, 최소 1,000대 기기 배포, 10만+ 측정 데이터 축적 |
| **완료 기준** | Flower FedAvg 클라이언트 구현. 차분 프라이버시 ε=1.0. 그래디언트 압축(대역폭 90% 절감). 로컬 미세조정(XGBoost 1,000행). 중앙 집계 서버(Go gRPC). 정확도 92%→96% 향상 실증 |
| **산출물** | §9.4를 정식 상세 섹션으로 확장 |

### 14.5 [강력 권고] B2B/B2G 배포 모델 + 관리자 대시보드

| 항목 | 내용 |
|---|---|
| **필요성** | 매출 3배 성장 경로. 병원/보건소 대량 배포 시 다중 테넌트·관리자 기능 필수 |
| **Phase** | Phase 3-4 |
| **선행 조건** | §8 Go 백엔드 안정화, FHIR 연동(§11.2) 완료, 인증(FDA/MFDS) 1개 이상 획득 |
| **완료 기준** | 다중 테넌트 DB 격리(PostgreSQL RLS). 관리자 웹 대시보드(실시간 기기 모니터링, 환자 리스트, 감사 로그). On-Premise 배포 옵션(Docker Compose). B2G API 게이트웨이(공공데이터 연동) |
| **산출물** | 별도 문서 `MPS-B2B-ARCH-v1.0`, 기술사양서 §8에 다중 테넌트 서브섹션 승격 |

### 14.6 [강력 권고] 시계열 DB 벤치마크 (TimescaleDB vs QuestDB)

| 항목 | 내용 |
|---|---|
| **필요성** | 10만대 기기 × 일 10회 측정 = 100만 행/일. 현재 TimescaleDB 선택 근거 미제시 |
| **Phase** | Phase 3 초입 (기술 검증) |
| **선행 조건** | §8 백엔드 기본 구현 완료, 벤치마크용 시뮬레이션 데이터 생성기 |
| **완료 기준** | 3개 DB(TimescaleDB, QuestDB, InfluxDB)에 대해 100만 행 삽입/조회 벤치마크. p99 쿼리 응답시간, 저장 비용(GB/월), 운영 복잡도 비교. CLAUDE.md §4 진화 원칙 4조건(≥10% 우위, 회귀 없음, 독립 검증, 영향도 통과) 충족 시만 교체 |
| **산출물** | 벤치마크 보고서 (독립), §8.8 DB 스키마에 결과 반영 |

### 14.7 [권고] SiPM-ECL 실측 검증 보고서

| 항목 | 내용 |
|---|---|
| **필요성** | §6.6 포화보정 알고리즘이 실측 데이터 없이 설계됨. Stage-2 진입 필수 선행 |
| **Phase** | Phase 3 (Stage-2 진입 Gate) |
| **선행 조건** | SiPM 센서 모듈 프로토타입 확보, ECL 시약 키트 확보, 광학 차폐 테스트 환경 |
| **완료 기준** | SiPM 펄스 데드타임(21ns) 실측 + 4PL 보정 곡선 실측 정확도 95% 이상. 포화 임계값(Count Rate) 3단계 검증. 전기화학 센서와의 시약 호환성 확인. 결합 후 3축 정확도 92% 이상 유지 확인 |
| **산출물** | `MPS-SIPM-VALIDATION-v1.0`, §6.6에 실측 데이터 반영 + [미검증]→[검증됨] 전환 |

### 14.8 [권고] Edge AI 성능 벤치마크 (TFLite <5ms 검증)

| 항목 | 내용 |
|---|---|
| **필요성** | §6.5 TFLite <5ms 목표가 모바일 디바이스 기준 실측 미입증 |
| **Phase** | Phase 2 (AI 모듈 구현 시) |
| **선행 조건** | TFLite 모델 학습 완료, 테스트 디바이스(중급 Android + iPhone SE 3) 확보 |
| **완료 기준** | 896차원 입력 기준: 모델 크기 <2MB, 추론 레이턴시 <10ms(p99), 에너지 소비 <5mJ/추론. 5ms 불가 시 목표 완화(10ms) + 사용자 체감 영향 분석. XGBoost 교차검증 포함 시 총 <20ms |
| **산출물** | 벤치마크 보고서, §6.5 성능 목표 갱신 |

---

### 앵커 상태 요약

| ID | 제목 | 우선순위 | Phase | 상태 |
|---|---|---|---|---|
| 14.1 | STRIDE 위협 모델링 | 필수 | 2 | [앵커] |
| 14.2 | SBOM 자동 생성 | 필수 | 2 | [앵커] |
| 14.3 | 다국어 i18n | 필수 | 2-3 | [앵커] |
| 14.4 | Federated Learning 상세 | 강력 권고 | 5 | [앵커] |
| 14.5 | B2B/B2G + 관리자 대시보드 | 강력 권고 | 3-4 | [앵커] |
| 14.6 | 시계열 DB 벤치마크 | 강력 권고 | 3 | [앵커] |
| 14.7 | SiPM-ECL 실측 검증 | 권고 | 3 | [앵커] |
| 14.8 | Edge AI 벤치마크 | 권고 | 2 | [앵커] |

> **승격 절차**: 앵커 항목의 구현이 완료되면, (1) 완료 기준 충족 확인, (2) 정식 섹션 번호 부여, (3) 본 표의 상태를 `[승격→§N.N]`으로 변경, (4) §8 변경이력에 기록.

---

## §Z 의견

**관찰**: 본 기술 사양서 v2.4는 v2.3의 자가검증(SelfDiagnostics, §6.15)에 **자가치유(SelfHealing, §6.16)**를 추가함으로써, 만파식 플랫폼이 "문제를 감지하는 것"을 넘어 "문제를 스스로 해결하는" 자율 복원 시스템으로 진화하였다. 이로써 4대 축이 완성되었다: ①측정(Measurement) ②분석(Analysis) ③자가진단(Self-Diagnosis) ④자가치유(Self-Healing). 자가치유는 6계층 전체를 관통한다. L1-HW(§6.16.1 HwAutoRecovery: BLE 재연결, ADC 리샘플, AFE 블록 스위치, 보정 클라우드 싱크), L2-FW(§6.16.2 FwSelfRepair: 5단계 에스컬레이션 Task→Stack→Peripheral→Soft→Factory), L3-Data(§6.16.3 DataPipelineReprocessor: 체크섬 재처리, AI 교차검증, CRDT 충돌 해소), L4-Session(§6.16.4 MeasurementSessionRecovery: 100ms 체크포인트, 5분 복원 윈도우), L4-App(§7.8 Flutter 복원 UI: 세션 복원 다이얼로그, 힐링 로그 타임라인), L5-Backend(§8.7 회복탄력성: CircuitBreaker, 지수 백오프 재시도, Bulkhead 격리, GracefulDegradation 4단계), L6-AI(§9.6 자가교정 엔진: Temperature Scaling 재보정, 피처 드리프트 α 적응, 온라인 점진 학습, 앙상블 폴백). 핵심 설계 원칙은 Detect→Classify→Attempt→Verify→Log의 5-phase 루프이며, 모든 치유 시도는 검증 후 적용/롤백이 결정되고 HashChain(§6.7)에 감사 기록이 남는다.

**리스크**: (1) 자가치유의 가장 큰 위험은 "잘못된 치유"다. 특히 AI 자가교정(§9.6) Phase 3(온라인 학습)에서 노이즈 데이터가 유입되면 모델이 악화될 수 있다. 현재 설계의 "연속 3회 실패 시 강제 롤백"과 "환자 안전 바이오마커 학습 제외"는 적절한 안전장치이나, 학습 데이터의 품질 검증 기준이 명시되지 않았다. (2) MeasurementSessionRecovery(§6.16.4)의 5분 복원 윈도우는 카트리지 시약의 화학적 안정성에 의존한다. 일부 면역센서(SiPM-ECL)는 반응 중 시약이 변질될 수 있으므로, 카트리지 유형별 복원 가능 시간을 차등 적용해야 한다. (3) Go 백엔드의 GracefulDegradation(§8.7)에서 Emergency 모드(측정 저장만)가 장시간 지속되면, 동기화 큐가 폭발적으로 증가하여 복구 시 부하 급증이 발생할 수 있다. 큐 크기 상한 및 우선순위 기반 처리 전략이 필요하다. (4) 자가치유와 자가검증이 동시에 동작할 때 리소스 경합(CPU, 메모리)이 발생할 수 있다. STM32F405의 168MHz/192KB SRAM 제약 하에서 치유 로직과 진단 로직의 동시 실행이 가능한지 실측이 필요하다.

**권고**: (1) 자가치유 구현은 §6.16.1(HwAutoRecovery)과 §6.16.4(SessionRecovery)를 Phase 3에서 우선 구현하고, FwSelfRepair(§6.16.2)와 DataPipelineReprocessor(§6.16.3)는 Phase 4에서, AI 자가교정(§9.6)은 Phase 5에서 구현할 것을 권장한다. (2) AI 온라인 학습(§9.6 Phase 3)의 학습 데이터 품질 게이트를 추가한다 — 최소 조건: 샘플 수 ≥ 50, 결측치 비율 < 5%, 이상치(3σ 이상) 비율 < 2%. (3) SessionRecovery의 복원 윈도우를 카트리지 매니페스트의 `stability_window_sec` 필드로 동적 결정하도록 확장한다. (4) Go 백엔드 Emergency 모드의 큐 상한을 10,000건으로 설정하고, 초과 시 oldest-first 드롭 + 알림 정책을 추가한다.

**대안**: (1) 현재 5-phase(Detect→Classify→Attempt→Verify→Log) 루프 대신, 3-phase(Detect→Heal→Log)로 단순화하는 것도 고려할 수 있다. Classify와 Verify 단계가 별도로 존재하면 치유 지연이 늘어나므로, 긴급 치유(예: 측정 중 BLE 단절)에서는 fast-path를 제공하는 것이 현실적이다. 현재 설계에서 `heal_during_measurement()`가 이 역할을 하지만, 일반 `heal()`과의 경로 분리가 더 명확해야 한다. (2) AI 자가교정에서 온라인 학습 대신 "모델 버전 풀(Model Pool)" 방식도 대안이다. 클라우드에서 사전 학습된 여러 버전의 모델을 보유하고, 드리프트 발생 시 가장 적합한 버전을 선택하는 것이 온라인 학습보다 안전할 수 있다. 이 경우 §9.6 Phase 3의 위험(데이터 품질 문제)을 근본적으로 제거할 수 있다.

**반대의견**: (1) 자가치유 시스템 자체가 복잡성의 원천이 될 수 있다. POCT 기기의 본질은 "정확한 측정"이며, 자가치유 로직이 측정 경로와 강하게 결합되면 디버깅 난이도가 급격히 상승한다. 최소한 v1.0 출시 시점에서는 HwAutoRecovery(BLE 재연결, ADC 리샘플)와 SessionRecovery만 포함하고, 나머지 치유 모듈은 OTA로 점진 배포하는 것이 리스크 관리 측면에서 현명하다. (2) AI 자가교정의 α 계수 적응(§9.6 Phase 2)은 차분식(§6.1)의 핵심 파라미터를 자동으로 변경하는 것이므로, 잘못된 적응이 측정 정확도를 직접 훼손할 수 있다. α 변경 범위를 현재 SSOT 기준 [0.90, 1.10]으로 제한하더라도, 자동 변경과 수동 변경의 경계를 명확히 구분해야 한다.

---

*본 문서는 ManPaSik SSOT 베이스라인 v1.1 기준이며, CLAUDE.md v2.1 §4 진화 원칙에 따라 운영한다. v2.4.2는 무결점 감사(4-way 병렬 에이전트)를 통해 Critical 1건(§6.3 핑거프린트 차원 공식), Major 4건(부동소수점 정규화, 사이트맵 화면 수/Phase 배분), Minor 4건(인터페이스 명확성)을 수정하였다. 다음 갱신: SiPM-ECL 실측, STRIDE 위협 모델링 별도 문서, SBOM 자동 생성 CI 통합, 또는 연합학습 파일럿 시점.*