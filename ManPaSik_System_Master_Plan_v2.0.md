# ManPaSik (萬波息) 통합 시스템 구축 마스터플랜 v2.0

**문서번호**: MPS-SYS-MASTER-v2.0  
**작성일**: 2026-04-18  
**보안등급**: RESTRICTED  
**SSOT 기준**: CLAUDE.md v2.1 베이스라인  
**선행문서**: MPS-APP-ANALYSIS-v1.0 (종합분석보고서), ManPaSik_AI_Ecosystem_Plan_v1.0_FINAL

---

## 목차

```
1.  시스템 비전 및 하네스 엔지니어링 원칙
2.  통합 아키텍처 (6-Layer + Harness Abstraction)
3.  하드웨어-펌웨어 연동 계층
4.  Rust 핵심 엔진 (신호처리/차동측정/AI)
5.  모바일 애플리케이션 (Flutter + Riverpod)
6.  백엔드 마이크로서비스
7.  AI/ML 파이프라인 및 디지털 트윈
8.  SDK 생태계 및 카트리지 마켓플레이스
9.  보안 아키텍처 및 규제 준수
10. 펌웨어 OTA 및 배포 전략
11. 글로벌 확장 및 현지화
12. 미래 확장 혁신 아이디어
13. 개발 로드맵 (30개월 6-Phase)
14. 비용 분석 및 인력 계획
15. SSOT 정합성 매트릭스
```

---

## 1. 시스템 비전 및 하네스 엔지니어링 원칙

### 1.1 비전

만파식은 차동측정 기반 범용분석 기술을 핵심으로, 건강·환경·안전 데이터를 단일 플랫폼에서 측정·분석·관리하는 글로벌 POCT 생태계이다. 88→896→1792차원 핑거프린트 벡터와 비표적 역추론 AI를 결합하여, 기존 POCT가 접근하지 못하는 미지 물질 동정까지 수행한다.

### 1.2 하네스 엔지니어링 (Harness Engineering) 원칙

하네스 엔지니어링은 복잡한 시스템을 **사전인증된 모듈 단위로 분리**하고, 새로운 카트리지/센서 추가 시 **전체 재인증 없이 모듈 단위 인증**만으로 확장하는 설계 철학이다.

**6대 원칙:**

| 원칙 | 정의 | 만파식 적용 |
|---|---|---|
| **H1. 모듈 독립성** | 각 모듈은 명확한 인터페이스로 분리, 내부 변경이 외부에 전파되지 않음 | AFE 블록별 HAL 트레이트, 카트리지별 독립 규제 경로 |
| **H2. 후방 호환** | 상위 버전이 하위 버전을 반드시 지원 | CSI v1.0 16핀 → v2.0 24핀 전환 시 어댑터 1세대 의무 제공 |
| **H3. 사전인증 플랫폼** | 리더기 플랫폼 1회 인증 후 카트리지는 변경사항만 제출 | 510(k) predicate 기반 카트리지별 간소화 인증 |
| **H4. 점진적 확장** | Stage별로 AFE 블록 추가, 한 번에 전체를 구현하지 않음 | Stage-1(전기화학) → Stage-2(광학) → Stage-3(NAAT) 순차 |
| **H5. 인터페이스 계약** | 모듈 간 통신은 버전화된 프로토콜/스키마로 계약 | BLE GATT 서비스 UUID 고정, NFC 매니페스트 버전화 |
| **H6. 실패 격리** | 하나의 카트리지/센서 오류가 전체 시스템을 중단시키지 않음 | Rust Result<T> + 에러 격리 + 폴백 경로 |

### 1.3 SSOT 베이스라인 정합 선언

본 문서의 모든 기술 사양은 다음 SSOT 베이스라인을 기준으로 한다. 불일치 발견 시 SSOT가 우선한다.

```
환율:           ₩1,480/USD
커넥터:         CSI v1.0 = Samtec MECF-08-01-L-DV, 16핀, 1.27mm
                (구 E12-IF 12핀은 부재로 폐기. 앱 문서 갱신 필요)
PPM:            49.70×30×4.30mm (모듈 크기)
차분식:         Sdiff_n = S_n − α_n × R_n
Universal AFE:  9블록 (전기화학/SiPM-ECL/색측정/TEC/LAMP/eNose/EC어레이/방사선/분자진단)
특허:           Family A (APP2026-0022KR) / B (APP2025-0967KR) / C (APP2025-0968KR)
정확도:         3축 결합 92~98%
```

---

## 2. 통합 아키텍처 (6-Layer + Harness Abstraction)

### 2.1 전체 계층 구조

```
┌─────────────────────────────────────────────────────────────────┐
│  Layer 6: 인프라 (Infrastructure)                                │
│  AWS/GCP Multi-Cloud, Kubernetes, CI/CD, Monitoring              │
├─────────────────────────────────────────────────────────────────┤
│  Layer 5: 백엔드 서비스 (Backend Services)                       │
│  Go + gRPC MSA (21 서비스), PostgreSQL + TimescaleDB + Milvus    │
├─────────────────────────────────────────────────────────────────┤
│  Layer 4: 모바일 앱 (Mobile Application)                         │
│  Flutter 3.x + Riverpod + flutter_rust_bridge (FFI)              │
├─────────────────────────────────────────────────────────────────┤
│  Layer 3: Rust 핵심 엔진 (Core Engine)                           │
│  신호처리 + 차동측정 + AI 추론 + 보안 + 디지털 트윈              │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  ★ Harness Abstraction Layer (HAL)                        │  │
│  │  SensorTrait → 9블록 AFE 통합 인터페이스                   │  │
│  │  CartridgeManifest v2.0 → Stage별 확장 가능               │  │
│  └───────────────────────────────────────────────────────────┘  │
├─────────────────────────────────────────────────────────────────┤
│  Layer 2: 하드웨어 제어 (Hardware Control)                       │
│  embedded-hal, GPIO/SPI/I2C/ADC, RAFE 스위치, EHD 제어           │
├─────────────────────────────────────────────────────────────────┤
│  Layer 1: 하드웨어 (Hardware)                                    │
│  STM32F405 + CSI v1.0 (16핀) + BLE nRF52832 + NFC PN7150        │
└─────────────────────────────────────────────────────────────────┘
```

### 2.2 Harness Abstraction Layer (HAL) 상세

HAL은 Layer 2(하드웨어)와 Layer 3(Rust 엔진) 사이에 위치하며, 모든 센서/AFE 블록을 통합 인터페이스로 추상화한다.

```rust
// harness/sensor_trait.rs — 모든 AFE 블록이 구현하는 통합 트레이트
pub trait SensorTrait {
    type Config;
    type RawData;
    type ProcessedData;
    type Error;

    /// 센서 초기화 (카트리지 매니페스트 기반)
    fn init(&mut self, config: &Self::Config) -> Result<(), Self::Error>;
    
    /// 원시 데이터 읽기 (DMA 기반 고속 수집)
    fn read_raw(&self) -> Result<Self::RawData, Self::Error>;
    
    /// 차동측정 적용: S_det - alpha * S_ref
    fn apply_differential(&self, raw: &Self::RawData, alpha: f64) -> Result<Self::ProcessedData, Self::Error>;
    
    /// 보정 데이터 로드 (NFC 태그 → 폴백 체인)
    fn load_calibration(&mut self, manifest: &CartridgeManifest) -> Result<CalibrationData, Self::Error>;
    
    /// 센서 상태 진단 (디지털 트윈 잔차 기반)
    fn self_diagnose(&self) -> SensorHealth;
    
    /// 안전 종료
    fn shutdown(&mut self) -> Result<(), Self::Error>;
}

// 9블록 AFE별 구현체
pub struct ElectrochemAfe;      // Stage-1: LMP91000×4 + ADS1256
pub struct SipmEclAfe;          // Stage-2: SiPM-ECL 광학 검출
pub struct ColorimetricAfe;     // Stage-2: LED/PD 색측정
pub struct TecAfe;              // Stage-2: 열전소자 온도 제어
pub struct LampAfe;             // Stage-3: 등온 NAAT (LAMP/RPA)
pub struct EnoseAfe;            // Stage-1: 전자코 8채널 어레이
pub struct EcArrayAfe;          // Stage-1: 전기화학 어레이 확장
pub struct RadiationAfe;        // Stage-4: 방사선 검출
pub struct MolecularAfe;        // Stage-5: 분자진단 통합

impl SensorTrait for ElectrochemAfe {
    type Config = ElectrochemConfig;
    type RawData = AdcSamples;
    type ProcessedData = DifferentialSignal;
    type Error = AfeError;
    // ... 구현
}
```

### 2.3 데이터 흐름 (End-to-End)

```
카트리지 삽입
    │
    ▼
[NFC PN7150] ─── SUN/CMAC 인증 ─── CartridgeManifest v2.0 파싱
    │
    ▼
[HAL] ─── SensorTrait::init() ─── AFE 블록 자동 선택 (매니페스트 기반)
    │
    ▼
[RAFE] ─── 측정 모드 자동 설정 (CA/CV/EIS/LSV/SWV/DPV/임피던스/OCP)
    │
    ▼
[ADC ADS1256] ─── 24-bit 디지털 변환 ─── DMA 고속 수집
    │
    ▼
[SensorTrait::apply_differential()] ─── S_det - α × S_ref
    │
    ▼
[DSP Pipeline] ─── FFT + 피크 검출 + 칼만 필터 ─── 88→896차원 핑거프린트
    │
    ▼
[TFLite XGBoost] ─── 분류 + 정량 + 비표적 이상탐지
    │
    ▼
[Digital Twin] ─── 잔차 모니터링 (y - ŷ) ─── 드리프트 보정 / 교정 시점 예측
    │
    ▼
[BLE nRF52832] ─── GATT Notify ─── Flutter UI 실시간 시각화
    │
    ▼
[CRDT Sync] ─── 오프라인 우선 저장 ─── 클라우드 동기화 (온라인 시)
```

---

## 3. 하드웨어-펌웨어 연동 계층

### 3.1 CSI v1.0 커넥터 사양 (확정)

```
커넥터:     Samtec MECF-08-01-L-DV
핀 수:      16핀 (기존 E12-IF 12핀에서 변경 확정)
피치:       1.27mm
형태:       메자닌 커넥터 (보드-투-보드)
이유:       12핀 에지 커넥터 적합 부품 부재
```

**16핀 할당 (Stage-1 기준):**

| 핀 | 신호명 | 방향 | 설명 |
|---|---|---|---|
| 1 | WE1 | Analog In | 감지전극 1 (Working Electrode) |
| 2 | RE1 | Analog In | 참조전극 1 (Reference Electrode) |
| 3 | CE1 | Analog Out | 대향전극 1 (Counter Electrode) |
| 4 | WE2 | Analog In | 감지전극 2 (차동 쌍) |
| 5 | RE2 | Analog In | 참조전극 2 (차동 쌍) |
| 6 | CE2 | Analog Out | 대향전극 2 |
| 7 | TEMP | Analog In | 온도 센서 (NTC) |
| 8 | HUMID | Analog In | 습도 센서 |
| 9 | EHD_HV | Power | EHD 고전압 출력 (기체 제어) |
| 10 | GND | Power | 접지 |
| 11 | VCC_3V3 | Power | 3.3V 전원 공급 |
| 12 | SPI_CLK | Digital | SPI 클록 (NFC 태그 직접 통신용) |
| 13 | SPI_MOSI | Digital | SPI 데이터 출력 |
| 14 | SPI_MISO | Digital | SPI 데이터 입력 |
| 15 | RESERVED_1 | - | Stage-2 광학 신호 예약 |
| 16 | RESERVED_2 | - | Stage-2 광학 전원 예약 |

> **후방호환 설계**: 핀 15-16은 Stage-1에서 NC(미연결)이며, Stage-2 광학 카트리지가 삽입될 때만 활성화. Stage-1 카트리지는 Stage-2+ 리더기에서 핀 15-16 무시로 정상 동작.

### 3.2 Stage별 커넥터 진화 로드맵

```
Stage-1 (확정): CSI v1.0 — Samtec 16핀 1.27mm
                전기화학 4채널 + EHD + 온/습도 + SPI + 예약 2핀
                
Stage-2 (계획): CSI v2.0 — 24핀 (미확정, Stage-2 진입 시 결정)
                + SiPM-ECL 광학 4핀 (HV_bias, SiPM_out, LED_drive, LED_GND)
                + LAMP 히터 2핀 (heater_ctrl, temp_feedback)
                + 확장 예약 2핀
                
Stage-3+ (비전): CSI v3.0 — 32핀 (장기)
                + 혈액학 임피던스 4핀
                + 고속 디지털 4핀 (향후 CMOS 이미징)
```

### 3.3 MCU 펌웨어 아키텍처

```
STM32F405RGT6 (ARM Cortex-M4, 168MHz, 1MB Flash, 192KB RAM)
│
├── RTOS Layer (FreeRTOS)
│   ├── Task: ADC_Sampler (Priority: Highest, 1kHz)
│   ├── Task: BLE_Comm (Priority: High)
│   ├── Task: NFC_Handler (Priority: Medium)
│   ├── Task: EHD_Controller (Priority: Medium)
│   ├── Task: Power_Manager (Priority: Low)
│   └── Task: OTA_Updater (Priority: Lowest, 백그라운드)
│
├── HAL Driver Layer
│   ├── adc_driver.c — ADS1256 SPI 통신, DMA 전송
│   ├── rafe_driver.c — LMP91000×4 I2C 설정
│   ├── ble_driver.c — nRF52832 UART/SPI 브릿지
│   ├── nfc_driver.c — PN7150 I2C, SUN/CMAC 인증
│   ├── ehd_driver.c — 고전압 펄스 제어 (PWM)
│   └── pmic_driver.c — BQ24195 충전/전원 관리
│
├── Measurement Engine (C → Rust FFI 브릿지 예정)
│   ├── differential_calc.c — S_det - α × S_ref
│   ├── signal_filter.c — 이동평균 + 칼만 필터
│   └── peak_detector.c — SWV/DPV 피크 검출
│
└── Security Layer
    ├── secure_boot.c — 부트 무결성 검증 (SHA-256)
    ├── firmware_sign.c — ECDSA P-256 서명 검증
    └── key_storage.c — STM32 OTP 영역 키 저장
```

### 3.4 BLE GATT 서비스 구조 (v2.0 확장)

```
ManPaSik Service (UUID: 0000FF00-0000-1000-8000-00805F9B34FB)
│
├── 0xFF01: Configuration Control [Write]
│   ├── measurement_mode: u8 (CA=0, CV=1, EIS=2, LSV=3, SWV=4, DPV=5, IMP=6, OCP=7)
│   ├── sampling_rate: u16 (Hz)
│   ├── afe_channel: u8 (0-3, 또는 0xFF=자동)
│   └── ehd_enable: bool
│
├── 0xFF02: Waveform Stream [Read/Notify]
│   ├── timestamp: u32 (ms)
│   ├── channel: u8
│   ├── raw_adc: i32 (24-bit signed)
│   ├── differential: f32 (S_det - α*S_ref)
│   └── sequence_num: u16 (패킷 재조립용)
│
├── 0xFF03: System Status [Read/Notify]
│   ├── battery_percent: u8
│   ├── temperature: i16 (0.1°C 단위)
│   ├── humidity: u16 (0.1% 단위)
│   ├── afe_status: u8 (비트마스크)
│   ├── error_code: u16
│   └── firmware_version: u32
│
├── 0xFF04: EHD Control [Write]
│   ├── voltage_kv: u16 (0.1kV 단위)
│   ├── pulse_width_ms: u16
│   └── flow_direction: u8 (0=흡입, 1=배출)
│
├── 0xFF05: Security Hash [Read]
│   ├── rolling_hash: [u8; 32] (SHA-256)
│   └── hash_index: u32
│
├── 0xFF06: OTA Control [Write/Notify] ← v2.0 신규
│   ├── ota_command: u8 (START=0, CHUNK=1, VERIFY=2, APPLY=3, ROLLBACK=4)
│   ├── chunk_data: [u8; 244] (BLE MTU 기반 최대 페이로드)
│   ├── chunk_index: u32
│   └── total_chunks: u32
│
└── 0xFF07: Digital Twin Sync [Read/Notify] ← v2.0 신규
    ├── predicted_value: f32 (ŷ)
    ├── residual: f32 (y - ŷ)
    ├── drift_score: f32 (0.0~1.0)
    ├── calibration_due: u32 (남은 측정 횟수)
    └── sensor_health: u8 (GOOD=0, DEGRADED=1, REPLACE=2)
```

---

## 4. Rust 핵심 엔진 (신호처리/차동측정/AI)

### 4.1 모듈 구조 (60개 모듈, Stage-1 기준 52개 활성)

```
manpasik-core/
├── Cargo.toml
├── src/
│   ├── lib.rs                          # FFI 엔트리포인트 (flutter_rust_bridge)
│   │
│   ├── harness/                        # ★ Harness Abstraction Layer
│   │   ├── mod.rs
│   │   ├── sensor_trait.rs             # SensorTrait 통합 인터페이스
│   │   ├── cartridge_manifest.rs       # CartridgeManifest v2.0 파서
│   │   ├── afe_registry.rs             # 9블록 AFE 동적 등록/조회
│   │   └── compatibility_checker.rs    # 카트리지-리더기 호환성 검증
│   │
│   ├── signal/                         # 신호처리
│   │   ├── mod.rs
│   │   ├── differential.rs             # 차동측정 S_det - α × S_ref
│   │   ├── dsp.rs                      # 디지털 필터 (Butterworth, Savitzky-Golay)
│   │   ├── fft.rs                      # FFT 분석 (실수 FFT, 윈도우 함수)
│   │   ├── peak_detector.rs            # SWV/DPV 피크 검출 알고리즘
│   │   ├── baseline_correction.rs      # 베이스라인 드리프트 보정
│   │   └── noise_reduction.rs          # 적응형 노이즈 제거
│   │
│   ├── calibration/                    # 보정 엔진
│   │   ├── mod.rs
│   │   ├── hybrid_correction.rs        # 4단계 하이브리드 보정
│   │   ├── kdm_drift.rs               # KDM 드리프트 보정 알고리즘
│   │   ├── temperature_comp.rs         # 온도 보상 (NTC 기반)
│   │   ├── humidity_comp.rs            # 습도 보상
│   │   ├── matrix_correction.rs        # 매트릭스 간섭 보정 (α 동적 조정)
│   │   └── calibration_store.rs        # NFC → QR → Cloud → 범용 폴백 체인
│   │
│   ├── quantification/                 # 정량화
│   │   ├── mod.rs
│   │   ├── concentration_engine.rs     # 농도 산출 (보정 곡선 기반)
│   │   ├── kalman_filter.rs            # 적응형 칼만 필터 (시계열 추적)
│   │   ├── multi_output.rs             # 다중 출력 (10종 동시 정량)
│   │   └── uncertainty.rs              # 측정 불확도 산출 (Clopper-Pearson)
│   │
│   ├── fingerprint/                    # 핑거프린트 벡터
│   │   ├── mod.rs
│   │   ├── feature_extractor.rs        # 88차원 기본 특성 추출
│   │   ├── multi_mode_expander.rs      # 88→448→896→1792차원 확장
│   │   ├── enose_fusion.rs             # 전자코 8채널 융합
│   │   ├── etongue_fusion.rs           # 전자혀 8채널 융합
│   │   └── vector_normalizer.rs        # L2 정규화 + 차원 축소 (PCA)
│   │
│   ├── ai/                             # AI/ML 추론
│   │   ├── mod.rs
│   │   ├── tflite_runtime.rs           # TFLite 엣지 추론 (양자화 INT8)
│   │   ├── xgboost_inference.rs        # XGBoost 로컬 추론
│   │   ├── classification.rs           # 물질 분류 (896차원 → 클래스)
│   │   ├── non_target_detector.rs      # 비표적 이상탐지 (역추론)
│   │   ├── model_registry.rs           # 모델 ID + SHA-256 해시 + 버전 관리
│   │   ├── ab_test_engine.rs           # A/B 테스트 (새 모델 vs 기존 병렬 추론)
│   │   └── xai_explainer.rs            # SHAP 값 기반 추론 근거 제공
│   │
│   ├── digital_twin/                   # 디지털 트윈
│   │   ├── mod.rs
│   │   ├── twin_engine.rs              # 센서 실측(y) vs 예측(ŷ) 실시간 비교
│   │   ├── drift_detector.rs           # 잔차(y-ŷ) 기반 드리프트 감지
│   │   ├── calibration_predictor.rs    # 교정 시점 예측 (남은 측정 횟수)
│   │   ├── cartridge_life.rs           # 카트리지 수명 추정
│   │   └── optimizer.rs                # 자기최적화 (Family A cl.1)
│   │
│   ├── communication/                  # 통신
│   │   ├── mod.rs
│   │   ├── ble_protocol.rs             # BLE GATT 서비스 (v2.0)
│   │   ├── nfc_handler.rs              # NFC SUN/CMAC 인증 + 매니페스트 파싱
│   │   ├── wifi_direct.rs              # Wi-Fi Direct 대용량 백업
│   │   └── usb_cdc.rs                  # USB-C 유선 통신
│   │
│   ├── security/                       # 보안
│   │   ├── mod.rs
│   │   ├── hash_chain.rs               # 해시 체인 (측정 데이터 무결성)
│   │   ├── rolling_hash.rs             # 롤링 해시 (실시간 변조 탐지)
│   │   ├── aes_gcm.rs                  # AES-256-GCM 암호화
│   │   ├── tpm_interface.rs            # TPM 2.0 인터페이스
│   │   └── pqc_hybrid.rs              # Post-Quantum 하이브리드 (향후)
│   │
│   ├── data/                           # 데이터 관리
│   │   ├── mod.rs
│   │   ├── local_store.rs              # SQLite 암호화 로컬 저장
│   │   ├── crdt_sync.rs               # CRDT 기반 오프라인 우선 동기화
│   │   ├── transform_log.rs            # 데이터 변환 감사 로그
│   │   └── fhir_export.rs             # FHIR R4 형식 내보내기
│   │
│   ├── cartridge/                      # 카트리지 관리
│   │   ├── mod.rs
│   │   ├── registry.rs                 # 카트리지 유형 동적 등록
│   │   ├── validator.rs                # 호환성 + 만료 + 사용횟수 검증
│   │   ├── nfc_data_parser.rs          # NFC 태그 데이터 구조 파싱
│   │   └── fallback_chain.rs           # NFC→QR→Cloud→범용보정 폴백
│   │
│   └── sipm_ecl/                       # ★ Stage-2 준비 (SiPM-ECL)
│       ├── mod.rs
│       ├── saturation_detector.rs      # 광자 파일업 포화 감지
│       ├── nonlinear_correction.rs     # 4PL 시그모이드 포화 보정
│       ├── self_optimization.rs        # 포화보정 → 디지털트윈 피드백
│       └── ecl_signal_processor.rs     # ECL 발광 신호 처리
```

### 4.2 차동측정 엔진 상세

```rust
// signal/differential.rs

/// SSOT 차분식: Sdiff_n = S_n − α_n × R_n
/// α는 매트릭스 보정 계수 (기본값 0.98, 범위 0.90~1.10)
pub struct DifferentialEngine {
    alpha: f64,
    alpha_range: (f64, f64),  // (0.90, 1.10)
    matrix_removal_rate: f64,  // 목표: 92~96%
}

impl DifferentialEngine {
    pub fn new() -> Self {
        Self {
            alpha: 0.98,
            alpha_range: (0.90, 1.10),
            matrix_removal_rate: 0.0,
        }
    }

    /// 단일 차동 연산
    pub fn compute(&self, s_detection: f64, s_reference: f64) -> f64 {
        s_detection - self.alpha * s_reference
    }

    /// 동적 α 조정 (Family B cl.11-13: AI 동적 가중치)
    /// 환경(온도/습도/노화) 기반 α 실시간 보정
    pub fn update_alpha(
        &mut self,
        temperature: f64,
        humidity: f64,
        sensor_age_hours: u32,
        ai_weight: Option<f64>,
    ) -> f64 {
        let temp_factor = 1.0 + (temperature - 25.0) * 0.001;  // 온도 계수
        let humid_factor = 1.0 + (humidity - 50.0) * 0.0005;   // 습도 계수
        let aging_factor = 1.0 - (sensor_age_hours as f64) * 0.00001;  // 노화 계수
        
        let new_alpha = match ai_weight {
            Some(w) => w,  // AI 모델이 제공하면 직접 사용
            None => self.alpha * temp_factor * humid_factor * aging_factor,
        };
        
        // 범위 클램핑 (안전 제한)
        self.alpha = new_alpha.clamp(self.alpha_range.0, self.alpha_range.1);
        self.alpha
    }

    /// 멀티채널 차동 배열 연산 (4채널 동시)
    pub fn compute_array(
        &self,
        detections: &[f64; 4],
        references: &[f64; 4],
    ) -> [f64; 4] {
        let mut results = [0.0f64; 4];
        for i in 0..4 {
            results[i] = self.compute(detections[i], references[i]);
        }
        results
    }
}
```

### 4.3 핑거프린트 확장 파이프라인

```
[Stage-1 기본]
4 센서페어 × 8 측정모드(CA,CV,EIS,LSV,SWV,DPV,IMP,OCP) = 32 기본 특성
+ 시간영역 특성 (피크, 기울기, 면적, 시정수 등) × 32 = 56 추가 특성
= 88차원 기본 벡터

[전자코 융합]
88차원 × 8채널 교차반응 패턴 = 448차원 (전자코/전자혀 각각)
eNose(448) + eTongue(448) = 896차원 통합 핑거프린트

[Stage-2+ 확장]
896차원 + 광학(SiPM-ECL) 256차원 + LAMP 형광 128차원 + 임피던스 512차원
= 1,792차원 풀스펙 핑거프린트

[AI 다중 출력]
896차원 → XGBoost → 분류 확률벡터 + 10종 동시 정량값
+ 비표적 이상탐지 점수 (Mahalanobis 거리)
```

### 4.4 SiPM-ECL 포화보정 (Stage-2 준비, Family A cl.30-31 대응)

```rust
// sipm_ecl/saturation_detector.rs

/// SiPM 광자 파일업 포화 검출
/// 고농도 시 광자 수가 SiPM 마이크로셀 수를 초과하면 포화 발생
/// → CV 급증 + 응답 비선형화
pub struct SaturationDetector {
    microcell_count: u32,       // SiPM 마이크로셀 수 (예: 14,410)
    dead_time_ns: f64,          // 마이크로셀 회복 시간 (예: 21ns)
    saturation_threshold: f64,  // 포화 판정 임계값 (예: 0.7)
}

impl SaturationDetector {
    /// 포화 여부 판정
    /// N_measured / N_max > threshold → 포화 영역
    pub fn detect(&self, measured_photons: u32) -> SaturationLevel {
        let ratio = measured_photons as f64 / self.microcell_count as f64;
        match ratio {
            r if r < 0.3 => SaturationLevel::Linear,
            r if r < 0.7 => SaturationLevel::SubLinear,
            _ => SaturationLevel::Saturated,
        }
    }
}

// sipm_ecl/nonlinear_correction.rs

/// 4PL 시그모이드 비선형 보정
/// N_true = -N_max × ln(1 - N_measured/N_max) × correction_factor
pub fn correct_saturation(
    measured: f64,
    n_max: f64,
    correction_factor: f64,
) -> Result<f64, SipmError> {
    let ratio = measured / n_max;
    if ratio >= 1.0 {
        return Err(SipmError::CompleteSaturation);
    }
    let n_true = -n_max * (1.0 - ratio).ln() * correction_factor;
    Ok(n_true)
}
```

---

## 5. 모바일 애플리케이션 (Flutter + Riverpod)

### 5.1 기술 스택 (확정)

```yaml
프레임워크:     Flutter 3.x (크로스 플랫폼 iOS/Android/Web PWA)
상태관리:       Riverpod 2.x + freezed (불변 상태)
Rust FFI:      flutter_rust_bridge 2.x
로컬 DB:       Hive (경량 구조화) + SQLite (암호화, sqlcipher)
BLE:           flutter_reactive_ble (최우선) + flutter_blue_plus (폴백)
NFC:           nfc_manager
차트:          fl_chart (경량) + syncfusion_flutter_charts (고급)
국제화:        flutter_localizations + intl (8→50개 언어)
테스트:        flutter_test + integration_test + mockito
```

### 5.2 화면 구조 (Information Architecture)

```
만파식 앱
├── 🏠 홈 (Home)
│   ├── 통합 건강/환경 스코어 (AI 산출)
│   ├── 최근 측정 요약 카드
│   ├── 긴급 알림 배너
│   └── 빠른 측정 버튼
│
├── 📊 측정 (Measure)
│   ├── 카트리지 스캔 (NFC 자동/QR 수동)
│   ├── 실시간 파형 표시 (BLE Stream)
│   ├── 측정 진행 상태 (단계별 가이드)
│   ├── 결과 표시 (정상/주의/경고/위험 4단계)
│   ├── AI 해석 코멘트
│   └── 결과 저장/공유
│
├── 📈 데이터 허브 (Data Hub)
│   ├── 시계열 추이 차트 (일/주/월/년)
│   ├── 건강-환경 상관관계 분석
│   ├── 리포트 생성 (PDF/FHIR)
│   ├── 데이터 내보내기 (CSV/JSON/FHIR)
│   └── 가족 구성원 비교 뷰
│
├── 🤖 AI 코치 (AI Coach)
│   ├── 대화형 건강 상담 (LLM 기반)
│   ├── 맞춤 운동/식단 추천
│   ├── 예측 경보 (트렌드 기반)
│   ├── 약물 상호작용 체크
│   └── 음성 AI 상담 (STT/TTS)
│
├── 🛍️ 마켓플레이스 (Marketplace)
│   ├── 카트리지 스토어 (건강/환경/안전)
│   ├── 구독 관리 (Basic/Bio-Opt/Clinical)
│   ├── 리더기 액세서리
│   └── SDK 개발자 포털 링크
│
└── ⚙️ 더보기 (More)
    ├── 리더기 관리 (최대 10대)
    ├── 가족 관리 / 보호자 모니터링
    ├── 의료 연동 (화상진료/병원예약)
    ├── 커뮤니티 (글로벌 + 실시간 번역)
    ├── 스마트홈 연동 (Matter/HomeKit)
    ├── 접근성 설정
    └── 계정/보안/언어 설정
```

### 5.3 Riverpod 상태 관리 패턴

```dart
// providers/measurement_provider.dart

@riverpod
class MeasurementNotifier extends _$MeasurementNotifier {
  @override
  MeasurementState build() => const MeasurementState.idle();

  Future<void> startMeasurement(CartridgeManifest manifest) async {
    state = const MeasurementState.preparing();
    
    // 1. HAL 초기화 (Rust FFI 호출)
    final rustBridge = ref.read(rustBridgeProvider);
    await rustBridge.initSensor(manifest: manifest);
    
    // 2. BLE 스트리밍 시작
    state = const MeasurementState.measuring(progress: 0.0);
    final bleStream = ref.read(bleStreamProvider);
    
    await for (final packet in bleStream.waveformStream()) {
      // 3. 실시간 차동측정 (Rust 엔진)
      final differential = await rustBridge.applyDifferential(
        sDetection: packet.rawAdc.toDouble(),
        sReference: packet.referenceAdc.toDouble(),
      );
      
      // 4. 상태 업데이트 (UI 자동 반영)
      state = MeasurementState.measuring(
        progress: packet.sequenceNum / manifest.totalSequences,
        currentValue: differential,
        waveformData: [...state.waveformData, differential],
      );
    }
    
    // 5. AI 추론 (Rust TFLite)
    final result = await rustBridge.runInference(
      fingerprint: state.waveformData,
      modelId: manifest.recommendedModelId,
    );
    
    // 6. 디지털 트윈 업데이트
    final twinStatus = await rustBridge.updateDigitalTwin(
      measured: result.concentrations,
    );
    
    state = MeasurementState.completed(
      result: result,
      twinHealth: twinStatus,
    );
  }
}

// 상태 정의 (freezed 불변 객체)
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
```

### 5.4 오프라인 우선 아키텍처 (CRDT)

```dart
// data/crdt_sync.dart

/// CRDT (Conflict-free Replicated Data Type) 기반 동기화
/// 오프라인에서 100% 완전 구동 후, 온라인 복귀 시 자동 병합
class CrdtSyncManager {
  final LocalStore _local;
  final CloudApi _cloud;
  
  /// LWW-Register: 마지막 쓰기 승리 (측정 결과용)
  /// G-Counter: 단조증가 카운터 (사용 횟수용)
  /// OR-Set: 관찰된 제거 집합 (카트리지 목록용)
  
  Future<void> syncWhenOnline() async {
    final pendingOps = await _local.getPendingOperations();
    
    for (final op in pendingOps) {
      switch (op.type) {
        case OpType.measurement:
          // LWW-Register: 타임스탬프 기반 충돌 해결
          await _cloud.mergeMeasurement(op.data, op.hlcTimestamp);
          break;
        case OpType.cartridgeUse:
          // G-Counter: 양쪽 최대값 채택
          await _cloud.mergeCounter(op.cartridgeId, op.count);
          break;
        case OpType.settingsChange:
          // LWW-Register: 가장 최근 설정 유지
          await _cloud.mergeSettings(op.data, op.hlcTimestamp);
          break;
      }
      await _local.markSynced(op.id);
    }
  }
}
```

---

## 6. 백엔드 마이크로서비스

### 6.1 기술 스택 (SSOT 통일)

```yaml
# 백엔드 기술 스택 확정 (생태계 기획안 vs 기술문서 불일치 해소)
API 서버:       Go 1.22+ (Fiber v2 → gRPC 전환 예정)
API 게이트웨이: Kong / Envoy (gRPC-Web 프록시)
메시지 큐:      Apache Kafka (이벤트 소싱)
데이터베이스:
  - PostgreSQL 16 + JSONB (관계형 + 반정형)
  - TimescaleDB (시계열 측정 데이터)
  - Redis Cluster (세션 + 보정 데이터 캐시)
  - Milvus (896차원 벡터 검색, IVF_SQ8)
  - Elasticsearch (전문 검색 + 로그)
인증:           Keycloak (자체 호스팅 OAuth 2.0 + JWT)
AI 서빙:        Triton Inference Server (GPU) + TFServing (CPU)
객체 저장:      AWS S3 / MinIO (온프레미스)
```

### 6.2 마이크로서비스 구성 (21개, 7개 도메인)

```
┌─────────────────────────────────────────────────────────────┐
│                    API Gateway (Kong/Envoy)                   │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  [사용자 도메인]           [측정 도메인]          [AI 도메인]  │
│  ├── user-service         ├── measurement-svc    ├── ai-inference│
│  ├── auth-service         ├── calibration-svc    ├── ai-training │
│  └── family-service       ├── cartridge-svc      ├── digital-twin│
│                           └── data-pipeline      └── xai-service │
│                                                               │
│  [의료 도메인]            [커머스 도메인]     [관리자 도메인]  │
│  ├── health-record-svc    ├── subscription-svc  ├── admin-svc    │
│  ├── telemedicine-svc     ├── marketplace-svc   ├── analytics-svc│
│  └── emergency-svc        └── payment-svc       └── audit-svc   │
│                                                               │
│  [IoT/커뮤니티 도메인]                                        │
│  ├── device-management-svc                                    │
│  ├── ota-service                                              │
│  └── community-svc                                            │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

### 6.3 핵심 API 엔드포인트

```protobuf
// proto/measurement.proto (gRPC)

service MeasurementService {
  // 측정 데이터 업로드 (스트리밍)
  rpc StreamMeasurement(stream MeasurementChunk) returns (MeasurementResult);
  
  // 보정 데이터 조회 (NFC 폴백용)
  rpc GetCalibration(CalibrationRequest) returns (CalibrationData);
  
  // 측정 이력 조회
  rpc GetMeasurementHistory(HistoryRequest) returns (stream MeasurementRecord);
  
  // 디지털 트윈 상태 동기화
  rpc SyncDigitalTwin(TwinState) returns (TwinPrediction);
}

service CartridgeService {
  // SDK 개발자: 카트리지 등록
  rpc RegisterCartridge(CartridgeDefinition) returns (RegistrationResult);
  
  // 카트리지 호환성 검증
  rpc ValidateCompatibility(CompatibilityRequest) returns (CompatibilityReport);
  
  // 마켓플레이스: 카트리지 목록 조회
  rpc ListCartridges(ListRequest) returns (CartridgeCatalog);
}

service AIService {
  // 클라우드 AI 추론 (엣지 보완)
  rpc Predict(PredictionRequest) returns (PredictionResult);
  
  // 연합학습 모델 가중치 집계
  rpc AggregateFederatedWeights(stream ModelWeights) returns (AggregatedModel);
  
  // PCCP: 모델 버전 관리
  rpc GetApprovedModel(ModelRequest) returns (ModelArtifact);
}
```

### 6.4 데이터베이스 스키마 (핵심 테이블)

```sql
-- TimescaleDB 하이퍼테이블 (측정 데이터)
CREATE TABLE measurements (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    device_id       UUID NOT NULL REFERENCES devices(id),
    cartridge_id    UUID NOT NULL REFERENCES cartridges(id),
    measured_at     TIMESTAMPTZ NOT NULL,
    
    -- 차동측정 결과
    raw_detection   DOUBLE PRECISION[] NOT NULL,  -- S_det 배열
    raw_reference   DOUBLE PRECISION[] NOT NULL,  -- S_ref 배열
    alpha           DOUBLE PRECISION NOT NULL,     -- 적용된 α값
    differential    DOUBLE PRECISION[] NOT NULL,   -- S_det - α*S_ref
    
    -- AI 추론 결과
    fingerprint_896 DOUBLE PRECISION[896],
    classification  JSONB,  -- {class: "glucose", confidence: 0.97}
    concentrations  JSONB,  -- {glucose: 120, cholesterol: 180, ...}
    anomaly_score   DOUBLE PRECISION,  -- 비표적 이상탐지 점수
    model_id        VARCHAR(64) NOT NULL,  -- PCCP 추적용
    model_hash      VARCHAR(64) NOT NULL,  -- SHA-256
    
    -- 디지털 트윈
    twin_predicted  DOUBLE PRECISION[],
    twin_residual   DOUBLE PRECISION[],
    drift_score     DOUBLE PRECISION,
    
    -- 무결성
    hash_chain      VARCHAR(64) NOT NULL,  -- 이전 해시 연결
    data_integrity  VARCHAR(64) NOT NULL,  -- 본 레코드 해시
    
    -- FHIR 호환
    fhir_observation_id VARCHAR(64)
);

SELECT create_hypertable('measurements', 'measured_at');

-- Milvus 벡터 컬렉션 (핑거프린트 유사도 검색)
-- collection: fingerprints
-- fields:
--   id: INT64 (primary key)
--   measurement_id: VARCHAR(36)
--   vector: FLOAT_VECTOR(896)
--   substance_class: VARCHAR(64)
-- index: IVF_SQ8, nlist=1024, nprobe=16

-- 카트리지 레지스트리 (SDK 개발자 등록)
CREATE TABLE cartridge_definitions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    developer_id    UUID NOT NULL REFERENCES developers(id),
    name            VARCHAR(128) NOT NULL,
    version         VARCHAR(16) NOT NULL,
    stage           SMALLINT NOT NULL,  -- 1, 2, 3, 4, 5
    
    -- 하네스 호환성
    csi_version     VARCHAR(8) NOT NULL,   -- "v1.0", "v2.0"
    pin_config      JSONB NOT NULL,        -- 핀 할당 맵
    afe_blocks      TEXT[] NOT NULL,        -- ["electrochem", "enose"]
    measurement_modes TEXT[] NOT NULL,      -- ["CA", "CV", "EIS"]
    
    -- 보정 파라미터
    calibration_schema JSONB NOT NULL,     -- NFC 태그 데이터 스키마
    
    -- 규제
    regulatory_class VARCHAR(16),          -- "Class I", "Class II"
    regulatory_path  VARCHAR(32),          -- "510k", "de_novo"
    approval_status  VARCHAR(16) DEFAULT 'pending',
    
    -- 마켓플레이스
    price_krw       INTEGER,
    revenue_share   DECIMAL(3,2) DEFAULT 0.70,  -- 개발자 70%
    
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    approved_at     TIMESTAMPTZ
);
```

---

## 7. AI/ML 파이프라인 및 디지털 트윈

### 7.1 3단계 AI 아키텍처

```
┌──────────────────────────────────────────────────────────┐
│  Level 1: 엣지 AI (리더기 + 앱 내)                        │
│  ├── TFLite INT8 양자화 XGBoost                           │
│  ├── 추론 시간: <1ms (Rust 네이티브)                      │
│  ├── 모델 크기: <2MB                                      │
│  ├── 오프라인 100% 동작                                   │
│  └── 분류 + 정량 + 기본 이상탐지                          │
├──────────────────────────────────────────────────────────┤
│  Level 2: 클라우드 AI (서버)                               │
│  ├── Triton Inference Server (GPU)                        │
│  ├── 대규모 앙상블 모델 (XGBoost + DNN)                   │
│  ├── 비표적 역추론 (Milvus 벡터 검색)                     │
│  ├── 환경-건강 상관관계 분석                               │
│  └── PCCP 관리 (모델 레지스트리 + A/B 테스트)             │
├──────────────────────────────────────────────────────────┤
│  Level 3: 연합학습 (Federated Learning)                    │
│  ├── Flower 프레임워크                                    │
│  ├── 기기별 로컬 학습 → 가중치만 서버 전송                │
│  ├── 개인정보 기기 외부 미유출 (GDPR/HIPAA 준수)          │
│  ├── 글로벌 질병 감시 네트워크 (FPCN)                     │
│  └── 모델 성능 지속 향상 (사용자 증가 → 정확도 향상)      │
└──────────────────────────────────────────────────────────┘
```

### 7.2 PCCP 대응 MLOps 파이프라인

FDA PCCP(Predetermined Change Control Plan) 최종 가이던스(2025)에 따라, AI 모델 변경을 사전 정의된 프로토콜 내에서 신규 마케팅 신청 없이 구현할 수 있는 체계를 구축한다.

```
[모델 개발] → [사전 검증] → [PCCP 범위 확인] → [A/B 테스트] → [단계적 배포]
                                    │
                    ┌───────────────┼───────────────┐
                    │               │               │
              범위 내 변경     범위 밖 변경     긴급 변경
              (자동 승인)     (FDA 재제출)    (72시간 사후보고)
```

**모델 레지스트리 구조:**

```yaml
model_registry:
  model_id: "MPS-XGB-GLUCOSE-v2.3.1"
  model_hash: "sha256:a1b2c3d4..."
  
  # PCCP 메타데이터
  pccp:
    change_description: "훈련 데이터 10,000건 추가, 하이퍼파라미터 튜닝"
    change_protocol: "MPS-PCCP-001 §3.2 (데이터 확장 허용 범위)"
    impact_assessment:
      accuracy_change: "+0.3% (95%CI: +0.1~+0.5%)"
      no_regression: true  # 어떤 서브그룹도 성능 저하 없음
      
  # 성능 스냅샷
  performance:
    accuracy: 0.967
    sensitivity: 0.952
    specificity: 0.978
    cv_percent: 1.8
    test_dataset: "MPS-VALIDATION-SET-v3"
    test_size: 5000
    
  # A/B 테스트 결과
  ab_test:
    duration_days: 14
    control_model: "MPS-XGB-GLUCOSE-v2.2.0"
    treatment_users: 500
    p_value: 0.003
    effect_size: 0.15
    
  # 배포 이력
  deployment:
    staged_rollout: [1%, 5%, 25%, 100%]
    rollback_trigger: "accuracy < 0.95 OR cv > 3.0%"
    current_stage: "100%"
    deployed_at: "2026-04-15T09:00:00Z"
```

### 7.3 디지털 트윈 엔진 (Family A cl.1)

```
┌─────────────────────────────────────────────────┐
│            Digital Twin Engine                    │
│                                                   │
│  ┌───────────┐    ┌───────────┐                  │
│  │ 실측값 (y) │    │ 예측값 (ŷ)│                  │
│  │ 센서 데이터│    │ 물리 모델 │                  │
│  └─────┬─────┘    └─────┬─────┘                  │
│        │                │                         │
│        └──────┬─────────┘                         │
│               │                                   │
│        ┌──────▼──────┐                            │
│        │ 잔차 (y-ŷ)  │                            │
│        │ 분석 엔진   │                            │
│        └──────┬──────┘                            │
│               │                                   │
│    ┌──────────┼──────────┐                        │
│    │          │          │                        │
│  ┌─▼──┐  ┌───▼───┐  ┌──▼───┐                    │
│  │드리프│  │교정    │  │수명  │                    │
│  │트감지│  │시점예측│  │추정  │                    │
│  └──┬──┘  └───┬───┘  └──┬───┘                    │
│     │         │         │                         │
│  ┌──▼─────────▼─────────▼──┐                     │
│  │ 자기최적화 피드백 루프    │                     │
│  │ α 동적 조정 + 보정 갱신  │                     │
│  └─────────────────────────┘                     │
└─────────────────────────────────────────────────┘

[드리프트 감지 기준]
- |잔차| > 3σ (3-시그마 룰) → 경고
- Allan 분산 τ > 임계값 → 장기 드리프트 판정
- 연속 5회 단방향 잔차 → 체계적 편향 감지

[교정 시점 예측]
- 잔차 추세 선형 회귀 → 허용 오차 초과 시점 추정
- 남은 측정 횟수 = (허용 오차 - 현재 잔차) / 드리프트 기울기
- BLE GATT 0xFF07 으로 앱에 실시간 전달

[카트리지 수명 추정]
- 사용 횟수 + 환경 노출(온도/습도 적산) + 드리프트 속도
- 잔여 수명 = f(사용횟수, 환경스트레스, 드리프트속도)
- 교체 권고 알림 → 마켓플레이스 재구매 유도
```

### 7.4 환경-건강 상관관계 AI (Cross-Domain Intelligence)

```python
# ai/cross_domain_analyzer.py
# 만파식 고유 강점: 건강+환경 동시 측정 데이터 상관관계 분석

class CrossDomainAnalyzer:
    """
    환경 지표(VOC, PM2.5, 라돈, CO)와 건강 지표(CRP, 혈당, 간기능)의
    시간지연 상관관계를 분석하여 인사이트를 생성한다.
    
    예시 인사이트:
    - "실내 VOC 상승 후 48시간 내 CRP 증가 패턴 감지"
    - "PM2.5 고농도 노출 3일 후 폐기능 지표 변화"
    - "라돈 만성 노출 지역의 간기능 지표 통계적 편향"
    """
    
    def analyze_lagged_correlation(
        self,
        env_timeseries: pd.DataFrame,  # 환경 측정 시계열
        health_timeseries: pd.DataFrame,  # 건강 측정 시계열
        max_lag_hours: int = 168,  # 최대 7일 시차
    ) -> List[CrossDomainInsight]:
        
        insights = []
        for env_col in env_timeseries.columns:
            for health_col in health_timeseries.columns:
                for lag in range(0, max_lag_hours, 6):
                    corr = self._compute_lagged_correlation(
                        env_timeseries[env_col],
                        health_timeseries[health_col],
                        lag_hours=lag,
                    )
                    if abs(corr.coefficient) > 0.3 and corr.p_value < 0.05:
                        insights.append(CrossDomainInsight(
                            env_factor=env_col,
                            health_factor=health_col,
                            lag_hours=lag,
                            correlation=corr.coefficient,
                            p_value=corr.p_value,
                            clinical_significance=self._assess_significance(corr),
                        ))
        
        return sorted(insights, key=lambda x: abs(x.correlation), reverse=True)
```

---

## 8. SDK 생태계 및 카트리지 마켓플레이스

### 8.1 전략적 포지셔닝

전 세계 POCT 업계에서 SDK 기반 카트리지 마켓플레이스를 운영하는 사업자는 없다. Abbott i-STAT, Roche cobas, Cue Health 모두 폐쇄형 카트리지 생태계를 유지한다. 만파식의 개방형 SDK 전략은 "POCT의 앱스토어"로서 글로벌 최초 포지션을 확보한다.

### 8.2 SDK 아키텍처

```
manpasik-cartridge-sdk/
├── Cargo.toml                     # Rust 메인 (crates.io 배포)
├── bindings/
│   ├── python/                    # PyPI: manpasik-sdk
│   ├── typescript/                # npm: @manpasik/sdk
│   └── c/                        # vcpkg: manpasik-sdk
│
├── src/
│   ├── builder.rs                 # CartridgeBuilder (빌더 패턴)
│   │   ├── set_name()
│   │   ├── set_stage()           # Stage-1~5
│   │   ├── set_csi_version()     # "v1.0" (16핀)
│   │   ├── add_measurement_mode() # CA, CV, EIS, ...
│   │   ├── set_calibration_schema()
│   │   ├── set_nfc_manifest()
│   │   └── build() → CartridgeDefinition
│   │
│   ├── validator.rs               # 자동 검증
│   │   ├── validate_pin_config()  # CSI 핀 맵 호환성
│   │   ├── validate_afe_support() # 리더기 AFE 블록 지원 여부
│   │   ├── validate_power()       # 전력 소비 예산
│   │   ├── validate_nfc_size()    # NFC 태그 용량 (172→256 바이트)
│   │   └── generate_report() → CompatibilityReport
│   │
│   ├── simulator.rs               # 시뮬레이션 검증
│   │   ├── simulate_measurement() # 가상 측정 시뮬레이션
│   │   ├── simulate_noise()       # 노이즈 주입 테스트
│   │   ├── simulate_drift()       # 드리프트 시뮬레이션
│   │   └── performance_report() → SimulationResult
│   │
│   ├── regulatory.rs              # 규제 분류 자동화
│   │   ├── classify_fda()         # FDA Class I/II/III 자동 분류
│   │   ├── classify_ce_ivdr()     # EU IVDR Class A/B/C/D
│   │   ├── classify_mfds()        # MFDS KGMP 등급
│   │   └── suggest_pathway() → RegulatoryPathway
│   │
│   └── publisher.rs               # 마켓플레이스 등록
│       ├── package()              # 카트리지 패키지 생성
│       ├── submit_review()        # 리뷰 제출
│       └── publish() → MarketplaceEntry
│
├── templates/                     # 카트리지 템플릿
│   ├── electrochemistry_basic.toml
│   ├── optical_immunoassay.toml
│   ├── lamp_naat.toml
│   └── gas_analysis.toml
│
└── docs/
    ├── getting_started.md
    ├── api_reference.md
    ├── hardware_guide.md          # CSI v1.0 핀맵 + 전기 사양
    └── regulatory_guide.md        # FDA/CE/MFDS 인증 가이드
```

### 8.3 카트리지 매니페스트 v2.0

```rust
// NFC 태그 데이터 구조 v2.0 (256바이트, v1.0 172바이트에서 확장)
#[repr(C, packed)]
pub struct CartridgeManifest {
    // === v1.0 호환 영역 (0-171 바이트, 기존과 동일) ===
    magic_number: [u8; 4],          // 0x4D504B32 ('MPK2', v1.0은 'MPK1')
    manifest_version: u8,           // 2 (v2.0)
    cartridge_type: u16,            // 카트리지 유형 코드
    serial_number: [u8; 12],        // 고유 일련번호
    manufacturing_date: u32,        // 제조일 (UNIX timestamp)
    expiration_date: u32,           // 만료일
    max_uses: u16,                  // 최대 사용 횟수
    current_uses: u16,              // 현재 사용 횟수
    calibration_data: [u8; 128],    // 보정 계수 (기울기/절편/온도/습도)
    firmware_min_ver: u32,          // 최소 펌웨어 버전
    regulatory_code: u32,           // 규제 승인 코드
    crc32: u32,                     // CRC32 (v1.0 영역)
    
    // === v2.0 확장 영역 (172-255 바이트) ===
    csi_version: u8,                // CSI 버전 (1=v1.0 16핀, 2=v2.0 24핀)
    stage: u8,                      // Stage (1-5)
    afe_blocks: u16,                // AFE 블록 비트마스크 (9비트)
    measurement_modes: u16,         // 측정 모드 비트마스크 (8비트)
    pin_config_hash: [u8; 8],       // 핀 설정 해시 (빠른 호환성 검증)
    
    // Stage-2+ 광학 파라미터 (v1.0 카트리지는 0으로 채움)
    optical_wavelength_nm: u16,     // LED 파장 (nm), 0=미사용
    optical_excitation_mv: u16,     // 여기 전압 (mV), 0=미사용
    sipm_bias_v: u16,               // SiPM 바이어스 전압 (0.1V 단위)
    
    // LAMP 파라미터 (Stage-3)
    lamp_temp_target: u16,          // LAMP 목표 온도 (0.1°C 단위)
    lamp_duration_sec: u16,         // LAMP 반응 시간 (초)
    
    // AI 모델 권장
    recommended_model_id: [u8; 16], // 권장 AI 모델 ID
    
    // 개발자 정보
    developer_id: [u8; 8],          // SDK 개발자 ID
    
    // 확장 예약
    reserved: [u8; 20],             // 향후 확장용
    
    crc32_v2: u32,                  // CRC32 (전체 256바이트)
}

// v1.0 카트리지 후방호환
impl CartridgeManifest {
    pub fn is_v1(&self) -> bool {
        self.magic_number == *b"MPK1"
    }
    
    pub fn is_v2(&self) -> bool {
        self.magic_number == *b"MPK2"
    }
    
    /// v1.0 카트리지를 v2.0 구조로 업캐스트 (확장 필드는 기본값)
    pub fn from_v1(v1_data: &[u8; 172]) -> Self {
        let mut manifest = Self::default();
        manifest.calibration_data.copy_from_slice(&v1_data[..172]);
        manifest.csi_version = 1;  // v1.0 = 16핀
        manifest.stage = 1;         // Stage-1 기본
        manifest
    }
}
```

### 8.4 수익분배 모델 (확정)

```
┌─────────────────────────────────────────────────────┐
│  카트리지 판매 수익 흐름                              │
│                                                       │
│  소비자 결제 (₩1,135~₩15,000/개)                     │
│       │                                               │
│       ├── 70% → 카트리지 개발자                       │
│       │   (MPS 자체 카트리지는 100% MPS)               │
│       │                                               │
│       ├── 15% → MPS 플랫폼 수수료                     │
│       │                                               │
│       ├── 10% → 제조/물류 비용                        │
│       │                                               │
│       └──  5% → 규제/인증 기금                        │
│           (카트리지별 인증 비용 적립)                   │
│                                                       │
│  ※ 기존 불일치 해소:                                  │
│    생태계 기획안 "70:30" + 기술문서 "15%"              │
│    → 개발자 70% / 플랫폼 15% / 제조 10% / 규제 5%    │
│      로 통합 확정                                     │
└─────────────────────────────────────────────────────┘
```

---

## 9. 보안 아키텍처 및 규제 준수

### 9.1 다층 보안 아키텍처

```
Layer 7: 애플리케이션 보안
├── RBAC + ABAC 하이브리드 접근 제어
├── 생체 인증 (Face ID / Touch ID)
├── JWT + Refresh Token 회전
└── 5-Tier 관리자 권한 분리 (총괄→국가→시군구→지점→판매점)

Layer 6: API 보안
├── Kong API Gateway + Rate Limiting
├── OAuth 2.0 + PKCE (모바일)
├── API Key + HMAC (SDK 개발자)
└── mTLS (서비스 간 통신)

Layer 5: 데이터 보안
├── AES-256-GCM (저장 암호화)
├── TLS 1.3 + Certificate Pinning (전송 암호화)
├── 해시 체인 + 롤링 해시 (측정 데이터 무결성)
└── PHI/PII 분리 저장 + 토큰화

Layer 4: 기기 보안
├── TPM 2.0 (리더기 하드웨어 신뢰 루트)
├── Secure Boot (부트 체인 무결성)
├── NFC SUN/CMAC (카트리지 정품 인증)
└── 펌웨어 ECDSA P-256 서명 검증

Layer 3: 네트워크 보안
├── BLE AES-128 페어링 + Secure Connections
├── Wi-Fi WPA3 + ECDHE 키교환
└── VPN (관리자/의료기관 전용)

Layer 2: 인프라 보안
├── AWS WAF + CloudFlare DDoS 방어
├── Kubernetes RBAC + NetworkPolicy
├── 컨테이너 이미지 서명 (cosign)
└── 시크릿 관리 (AWS Secrets Manager / Vault)

Layer 1: Post-Quantum 준비
├── PQC 하이브리드 알고리즘 대비 (ML-KEM + X25519)
├── 모듈형 암호 라이브러리 (알고리즘 교체 가능)
└── NIST PQC 표준 확정 시 마이그레이션 계획
```

### 9.2 규제 준수 매트릭스

| 규제 | 지역 | 대상 | 현재 상태 | Stage |
|---|---|---|---|---|
| ISO 13485:2016 | 글로벌 | QMS | 설계 착수 시 구축 시작 | 전 Stage |
| IEC 62304:2015 | 글로벌 | SW 수명주기 | Class B 설계 적용 | 전 Stage |
| ISO 14971:2019 | 글로벌 | 위험관리 | FMEA 수행 예정 | 전 Stage |
| FDA 21 CFR 820 | 미국 | QSR → QMSR | QMSR 전환 대응 설계 | Stage-1~ |
| FDA 510(k) | 미국 | 시판전 승인 | Predicate 조사 중 | Stage-1 |
| FDA PCCP | 미국 | AI 변경관리 | MLOps 파이프라인 설계 완료 | Stage-1~ |
| CE-IVDR 2017/746 | EU | IVD | Class C 이하 설계 | Stage-1+ |
| MFDS KGMP | 한국 | 체외진단 | 설계관리 문서 준비 | Stage-1 |
| GDPR | EU | 개인정보 | DPO 지정 + DPIA 설계 | 전 Stage |
| HIPAA | 미국 | PHI | BAA 체결 프로세스 설계 | 전 Stage |
| KC 인증 | 한국 | EMC/전기안전 | EMI 차폐 설계 적용 | MVP |
| FCC Part 15 | 미국 | 전자파 | BLE/Wi-Fi 인증 계획 | MVP |

---

## 10. 펌웨어 OTA 및 배포 전략

### 10.1 OTA 아키텍처

```
┌─────────────────────────────────────────────┐
│  OTA 업데이트 파이프라인                       │
│                                               │
│  [빌드] → [서명] → [스테이징] → [단계적 배포]  │
│                                               │
│  빌드: GitHub Actions CI/CD                   │
│  서명: ECDSA P-256 (HSM 기반 키 관리)         │
│  스테이징: 내부 테스트 그룹 (100대)            │
│  배포: 1% → 5% → 25% → 100% (4단계)          │
│                                               │
│  전송 경로 (우선순위):                         │
│  1. Wi-Fi Direct (최대 250Mbps, 대용량)       │
│  2. BLE (최대 2Mbps, 델타 업데이트)           │
│  3. USB-C (최대 480Mbps, 유선 폴백)           │
│                                               │
│  안전 메커니즘:                                │
│  ├── A/B 파티션 (현재/백업 이중화)             │
│  ├── 자동 롤백 (부팅 실패 3회 → 이전 버전)     │
│  ├── 배터리 잔량 확인 (>30% 시만 업데이트)      │
│  ├── 측정 중 업데이트 차단                     │
│  └── 무결성 검증 실패 시 업데이트 거부          │
└─────────────────────────────────────────────┘
```

### 10.2 델타 업데이트

```
전체 펌웨어:  ~512KB (STM32F405 1MB Flash 중 A/B 각 512KB)
델타 패치:    ~20-50KB (bsdiff 알고리즘)
전송 시간:    BLE 기준 10-25초 (델타), 4-5분 (전체)

델타 생성:    bsdiff(old_firmware, new_firmware) → patch
델타 적용:    bspatch(current_firmware, patch) → new_firmware
검증:         SHA-256(new_firmware) == expected_hash
```

---

## 11. 글로벌 확장 및 현지화

### 11.1 지역별 확장 로드맵

```
Phase 1 (0-18개월): 한국 시장 집중
├── MFDS KGMP 인증
├── KC 인증 (EMC/전기안전)
├── 한국어 + 영어 지원
└── 국내 임상 시험

Phase 2 (12-24개월): 미국 시장 진입
├── FDA 510(k) 제출 (Stage-1)
├── FCC Part 15 인증
├── 영어 + 스페인어 추가
└── CLIA waiver 취득

Phase 3 (18-30개월): 유럽/일본 확장
├── CE-IVDR (Class C 이하)
├── PMDA 인증 (일본)
├── 독일어, 프랑스어, 일본어 추가
└── 현지 임상 데이터 보강

Phase 4 (24-36개월): 아시아/신흥국
├── 중국 NMPA, 인도 CDSCO
├── 중국어(간체/번체) 추가
├── 현지 제조 파트너십
└── 저가 카트리지 라인업
```

### 11.2 다국어 실시간 번역 시스템

```yaml
기술 스택:
  기본: flutter_localizations + ARB 파일 (정적 번역)
  실시간: SeamlessM4T (Meta) 또는 NLLB-200
  음성: Whisper STT + VITS TTS (다국어)
  
초기 8개 언어: 한국어, 영어, 일본어, 중국어(간체), 중국어(번체), 스페인어, 독일어, 프랑스어
확장 50+ 언어: 실시간 번역 엔진으로 커뮤니티/AI 코치 대화 지원

커뮤니티 번역:
  - 전문가 포럼 글 → 실시간 번역 → 전 세계 사용자 접근
  - 의료 용어는 SNOMED-CT/LOINC 코드 기반 정확도 보장
  - 번역 품질 피드백 루프 (사용자 교정 → 모델 학습)
```

---

## 12. 미래 확장 혁신 아이디어

### 12.1 Federated POCT Network (FPCN) — 글로벌 질병 감시

```
[개념]
전 세계 만파식 리더기가 로컬에서 학습한 모델 가중치만 중앙 서버로 전송하여,
글로벌 질병 감시(Disease Surveillance) 네트워크를 구축한다.

[메커니즘]
1. 각 리더기: 측정 데이터로 로컬 모델 미세 조정
2. Flower 서버: 가중치 집계 (FedAvg / FedProx)
3. 집계 모델: 지역별 이상 패턴 감지
4. 경보 시스템: 특정 지역 CRP 급증 → 감염병 유행 조기 경보

[프라이버시]
- 개인 측정 데이터는 기기를 떠나지 않음
- 모델 가중치만 전송 (차분 프라이버시 적용)
- GDPR/HIPAA 완전 준수

[비즈니스]
- B2G: 질병관리청, CDC, WHO에 감시 데이터 판매
- B2B: 보험사에 지역별 위험도 데이터 제공
- B2R: 제약사에 임상 시험 대상자 발굴 데이터 제공

[기술 요구사항]
- Flower 프레임워크 서버 구축
- 차분 프라이버시 (ε=1.0) 적용
- 지역 이상탐지 알고리즘 (Isolation Forest)
- 실시간 대시보드 (지역별 히트맵)
```

### 12.2 환경-건강 교차분석 AI (Cross-Domain Intelligence)

```
[세계 유일의 데이터]
만파식은 건강(혈액/타액) + 환경(공기/수질) 데이터를 동일 사용자로부터 수집하는
세계 유일의 POCT 플랫폼이다. 이 데이터 조합은 기존 어떤 의료기기도 보유하지 못한다.

[인사이트 예시]
- "실내 VOC 150ppb 초과 → 48시간 후 CRP 0.3mg/dL 상승 (p<0.01)"
- "수돗물 잔류염소 > 1.5ppm 지역 → 피부 바이오마커 이상 빈도 2.3배"
- "PM2.5 > 75μg/m³ 노출 72시간 → 폐기능 지표 3% 감소"
- "라돈 > 4 pCi/L 만성 노출 → 간기능 지표 경계선 이상 OR=1.8"

[수익화]
- B2B: 역학 연구기관, 보험사, 건설사에 익명화 통계 데이터 판매
- B2C: 프리미엄 구독(Clinical Guard)에 교차분석 리포트 포함
- B2G: 환경부/질병관리청에 실시간 환경-건강 상관 대시보드 제공
```

### 12.3 AR/XR 측정 가이드

```
[Phase 1: 스마트폰 AR]
- ARKit/ARCore 기반 카메라 오버레이
- 카트리지 삽입 위치 AR 가이드
- 시료 주입량 시각적 안내 (20μL 정확한 주입)
- 측정 중 실시간 AR 상태 표시

[Phase 2: 스마트글래스]
- Apple Vision Pro / Meta Quest 연동
- 핸즈프리 측정 (의료/실험실 환경)
- 결과 공간 고정 표시 (Spatial Anchor)

[Phase 3: 원격 AR 지원]
- WebRTC + AR 오버레이
- 전문가가 사용자 카메라를 보며 실시간 AR 마킹으로 가이드
- 개발도상국 현장 지원에 활용
```

### 12.4 음성 AI "제로터치" 측정

```
[개념]
시각장애인, 손이 불편한 사용자가 음성만으로 전체 측정 프로세스를 완료하는 모드

[워크플로우]
1. "만파식, 혈당 측정 시작" → LLM이 의도 파악
2. "카트리지를 삽입해주세요" → TTS 음성 안내
3. [NFC 자동 인식] → "혈당 카트리지가 인식되었습니다"
4. "시료를 주입해주세요" → 진동 + 음성 가이드
5. [자동 측정 시작] → "측정 중입니다. 약 90초 소요됩니다"
6. "혈당 120 밀리그램입니다. 정상 범위입니다" → 결과 음성 안내
7. "추가 조언을 원하시면 말씀해주세요" → AI 코치 연계

[기술 스택]
- STT: Whisper (로컬, 오프라인 가능)
- NLU: LLM 기반 의도 파악 (온라인) / 규칙 기반 (오프라인)
- TTS: VITS (다국어, 자연스러운 음성)
- 대화관리: FSM + LLM 하이브리드
```

### 12.5 카트리지 개발자 커뮤니티 플랫폼

```
[개념]
GitHub + npm + 앱스토어를 결합한 카트리지 개발자 전용 플랫폼

[기능]
- 카트리지 설계 공유 (오픈소스 / 상용)
- 시뮬레이션 클라우드 (가상 측정 테스트)
- 규제 가이드 위자드 (FDA/CE/MFDS 경로 자동 추천)
- 개발자 인증 프로그램 (Certified ManPaSik Developer)
- 해커톤 / 챌린지 (신규 측정 항목 개발 대회)
- 수익 대시보드 (판매량, 리뷰, 수수료 정산)

[생태계 성장 전략]
Year 1: 내부 5종 카트리지 + SDK 베타 (개발자 50명)
Year 2: 외부 10종 추가 + SDK 정식 출시 (개발자 200명)
Year 3: 50종 카트리지 + 글로벌 개발자 1,000명
Year 5: 200종+ + 자생적 생태계 (개발자 5,000명+)
```

### 12.6 예측 보전 + 구독 연계 (Predictive Maintenance as a Service)

```
[개념]
디지털 트윈 엔진의 드리프트 감지 + 카트리지 수명 예측을 구독 서비스와 결합

[사용자 경험]
1. 디지털 트윈: "리더기 센서 드리프트 감지. 20회 측정 후 교정 필요"
2. 앱 알림: "교정 카트리지를 자동 주문합니다" (Clinical Guard 구독자)
3. 마켓플레이스: 교정 카트리지 자동 배송
4. 사용자: 교정 카트리지 삽입 → 자동 교정 수행

[비즈니스 가치]
- 구독 해지율(Churn) 감소: 예측 보전이 주는 가치 체감
- 카트리지 반복 구매 자동화: LTV(고객생애가치) 극대화
- 측정 정확도 유지: 임상 신뢰도 + 규제 준수
```

---

## 13. 개발 로드맵 (30개월 6-Phase)

### 13.1 전체 타임라인

```
Phase 0: 기반 구축 (Month 1-3)
├── Rust 코어 엔진 아키텍처 확정 + HAL 트레이트 정의
├── Flutter 프로젝트 스캐폴딩 + 디자인 시스템 구축
├── Go gRPC 백엔드 기본 구조 + DB 스키마
├── CSI v1.0 16핀 HW 프로토타입 브링업
├── CI/CD 파이프라인 (GitHub Actions + ArgoCD)
└── ISO 13485 QMS 초기 문서화

Phase 1: MVP (Month 4-8)
├── BLE/NFC 리더기-앱 연동 완성
├── 차동측정 엔진 Rust 구현 + 검증
├── Stage-1 전기화학 카트리지 3종 (혈당/전해질/CRP)
├── 기본 UI (홈/측정/결과/설정)
├── 오프라인 로컬 저장 + 기본 클라우드 동기화
├── 기본 보안 (AES-256 + 해시 체인)
└── 내부 알파 테스트 (50명)

Phase 2: 핵심 기능 (Month 9-14)
├── AI 파이프라인 (TFLite XGBoost + 분류 + 정량)
├── 디지털 트윈 v1.0 (드리프트 감지 + 교정 예측)
├── CRDT 오프라인 완전 구동
├── AI 코치 v1.0 (규칙 기반 + 기본 LLM)
├── 데이터 허브 (시계열 차트 + 리포트)
├── 카트리지 마켓플레이스 MVP (자체 5종)
├── SaaS 구독 3티어 + 결제 시스템
├── SDK v0.9 베타 (내부 개발자 10명)
├── HealthKit/Health Connect 연동
└── 클로즈드 베타 테스트 (500명)

Phase 3: 규제 인증 + 확장 (Month 15-20)
├── MFDS KGMP 인증 신청
├── KC/FCC 인증 획득
├── FDA 510(k) Pre-Submission (Stage-1)
├── PCCP MLOps 파이프라인 운영 개시
├── SDK v1.0 정식 출시 (외부 개발자 50명)
├── 전자코 융합 → 896차원 핑거프린트
├── 비표적 역추론 v1.0
├── 화상진료 연동 (WebRTC)
├── 가족 관리 + 보호자 모니터링
├── 다국어 8개 언어 지원
└── 오픈 베타 (5,000명)

Phase 4: Stage-2 준비 + 글로벌 (Month 21-26)
├── SiPM-ECL 광학 모듈 HW 개발
├── CSI v2.0 (24핀) 설계 검토 (§4 진화 원칙)
├── sipm_ecl Rust 모듈 구현 + 포화보정
├── Stage-2 카트리지 프로토타입 (면역분석)
├── FDA 510(k) 정식 제출 (Stage-1)
├── CE-IVDR 신청 준비
├── 연합학습 v1.0 (Flower, 1,000+ 기기)
├── 환경-건강 교차분석 AI v1.0
├── Matter 스마트홈 연동
├── 음성 AI 제로터치 v1.0
└── 글로벌 정식 출시 (한국 + 미국 동시)

Phase 5: 생태계 성숙 (Month 27-30)
├── SDK 마켓플레이스 본격 운영 (외부 카트리지 10종+)
├── FPCN 글로벌 질병 감시 파일럿
├── AR 측정 가이드 v1.0
├── Stage-3 LAMP/NAAT 모듈 설계 착수
├── B2B 대시보드 (기업/의료기관)
├── 개발자 커뮤니티 플랫폼 런칭
├── 예측 보전 + 구독 자동화
└── 카트리지 200종 달성 목표 로드맵 수립
```

### 13.2 Phase별 마일스톤 게이트

| Gate | 시점 | 통과 기준 |
|---|---|---|
| G1: MVP Ready | Month 8 | BLE 측정 성공률 >95%, 차동측정 CV <3%, 내부 알파 NPS >40 |
| G2: Beta Ready | Month 14 | AI 정확도 >92%, 오프라인 100% 구동, SDK 빌드 가능, 구독 결제 정상 |
| G3: Regulatory Submit | Month 20 | MFDS 서류 완비, 510(k) Pre-Sub 피드백 반영, 임상 데이터 100건+ |
| G4: Global Launch | Month 26 | FDA/MFDS 승인, SDK 마켓플레이스 5종+, MAU 10,000+ |
| G5: Ecosystem Maturity | Month 30 | 외부 카트리지 10종, 개발자 200명, ARR ₩10억+ |

---

## 14. 비용 분석 및 인력 계획

### 14.1 총 예산 (30개월)

```
총 예산: ₩85억 (기존 ₩67억 대비 +₩18억, 확장 범위 반영)

[인건비] ₩54억 (63%)
├── Phase 0-1: 18명 × 8개월 = ₩14.4억
├── Phase 2-3: 28명 × 12개월 = ₩28.0억
└── Phase 4-5: 24명 × 10개월 = ₩11.6억

[인프라/클라우드] ₩12억 (14%)
├── AWS/GCP: ₩5억
├── CI/CD + 모니터링: ₩2억
├── AI 학습 GPU: ₩3억
└── 테스트 장비/리더기: ₩2억

[규제/인증] ₩10억 (12%)
├── MFDS KGMP: ₩2억
├── FDA 510(k): ₩4억
├── CE-IVDR: ₩2억
├── KC/FCC: ₩1억
└── 임상 시험: ₩1억

[외주/라이선스] ₩5억 (6%)
├── UI/UX 디자인: ₩2억
├── 보안 감사: ₩1억
├── 법률/특허: ₩1억
└── 소프트웨어 라이선스: ₩1억

[예비비] ₩4억 (5%)
```

### 14.2 핵심 인력 구성 (피크 28명)

| 역할 | 인원 | Phase | 핵심 스킬 |
|---|---|---|---|
| PM/PO | 2 | 전 Phase | 의료기기 PM 경험, 규제 이해 |
| Rust 엔진 개발자 | 4 | 전 Phase | embedded Rust, DSP, no_std |
| Flutter 개발자 | 4 | Phase 1~ | Riverpod, BLE, 크로스플랫폼 |
| Go 백엔드 개발자 | 3 | Phase 1~ | gRPC, Kubernetes, 분산시스템 |
| AI/ML 엔지니어 | 3 | Phase 2~ | XGBoost, TFLite, 연합학습 |
| 펌웨어 엔지니어 | 2 | 전 Phase | STM32, FreeRTOS, BLE/NFC |
| HW 엔지니어 | 2 | Phase 0-4 | PCB, 아날로그, CSI 커넥터 |
| QA/RA | 3 | Phase 2~ | IEC 62304, FDA 510(k), ISO 13485 |
| DevOps/SRE | 2 | Phase 1~ | K8s, CI/CD, 모니터링 |
| UI/UX 디자이너 | 2 | Phase 0-3 | 의료기기 UI, 접근성 |
| 보안 전문가 | 1 | Phase 2~ | TPM, 암호화, HIPAA |

---

## 15. SSOT 정합성 매트릭스

본 문서의 모든 기술 사양이 SSOT 베이스라인과 일치하는지 최종 교차검증한 결과이다.

| 항목 | SSOT 베이스라인 | 본 문서 기재 | 일치 |
|---|---|---|---|
| 환율 | ₩1,480/USD | ₩1,480/USD | ✅ |
| 커넥터 | CSI v1.0 Samtec 16핀 | CSI v1.0 Samtec 16핀 | ✅ |
| PPM | 49.70×30×4.30mm | 49.70×30×4.30mm (모듈) | ✅ |
| 차분식 | Sdiff_n = S_n − α_n × R_n | S_det - α × S_ref | ✅ (표기 변환) |
| Universal AFE | 9블록 | 9블록 (Stage별 활성화) | ✅ |
| 특허 A | APP2026-0022KR | cl.1,4-5,11,19,24 반영 | ✅ |
| 특허 B | APP2025-0967KR | cl.1,3,8,11-13,17-20 반영 | ✅ |
| 특허 C | APP2025-0968KR | 모듈형 플랫폼, OTA, 보안 반영 | ✅ |
| SiPM-ECL | AFE 9블록 포함 + cl.30-31 예정 | Stage-2 sipm_ecl 모듈 설계 | ✅ |
| 정확도 | 3축 결합 92~98% | 92~98% (Stage별 검증 예정) | ✅ |
| 실현가능성 | Stage 1: 93% | 구현 계획 반영 [추정] | ✅ |

---

## §Z 의견

**관찰**: 본 마스터플랜은 기존 생태계 기획안(v1.0)과 기술문서(v1.0)의 7건 불일치를 해소하고, 하네스 엔지니어링 6대 원칙을 전 계층에 적용하였다. 60개 Rust 모듈, 21개 백엔드 서비스, 30개월 6-Phase 로드맵으로 구성되며, 기존 문서 대비 AFE HAL 추상화, PCCP MLOps, 카트리지 매니페스트 v2.0, SiPM-ECL 준비, 5대 혁신 아이디어가 신규 추가되었다.

**리스크**: 30개월 ₩85억은 의욕적 목표이며, Phase 3 규제 인증 단계에서 일정 지연이 가장 큰 리스크이다. FDA 510(k) Pre-Submission 피드백에 따라 추가 임상 데이터가 요구될 수 있으며, 이 경우 6-12개월 추가 소요가 가능하다. 두 번째 리스크는 SDK 생태계 초기 개발자 확보이다. 카트리지 하드웨어 개발 진입장벽이 소프트웨어 대비 높으므로, 초기에는 대학/연구소 파트너십으로 생태계를 점화해야 한다.

**권고**: Phase 0-1에서 CSI v1.0 16핀 기반 HW-SW 정합성을 완벽히 검증하는 것이 전체 프로젝트의 성패를 좌우한다. "작동하는 MVP를 빠르게" 전략으로 Month 8까지 혈당 1종 카트리지의 엔드투엔드 측정을 완성하고, 이를 기반으로 투자 유치와 규제 대응을 병행할 것을 권고한다.

**반대의견**: 30개월 로드맵이 공격적이라는 우려가 있을 수 있다. 특히 규제 인증과 글로벌 확장을 동시에 진행하는 Phase 4는 인력 부하가 집중된다. 이 경우 글로벌 확장을 Phase 5로 이동하고 Phase 4를 규제 전담으로 재편하는 보수적 대안도 검토할 수 있다.

---

*문서 끝. 본 문서는 ManPaSik SSOT 베이스라인 v1.1 기준으로 작성되었으며, CLAUDE.md v2.1의 §4 진화 원칙에 따라 입증된 우위가 확인될 때만 갱신한다.*


