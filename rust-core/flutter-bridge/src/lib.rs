//! Flutter-Rust Bridge API
//!
//! Flutter에서 호출 가능한 Rust 함수들

mod frb_generated;

use manpasik_engine::{
    ble::{BleManager, ConnectionState, DeviceInfo},
    differential::{
        CorrectionParams, DifferentialCorrection, DifferentialEngine, MeasurementResult,
    },
    fingerprint::{FingerprintVector, MeasurementType},
    nfc::{CartridgeInfo, CartridgeType, NfcReader},
};

// flutter_rust_bridge 매크로
pub use flutter_rust_bridge::frb;

/// ManPaSik 엔진 상태
pub struct ManpasikEngine {
    differential: DifferentialEngine,
    ble: BleManager,
    nfc: NfcReader,
}

impl ManpasikEngine {
    /// 새 엔진 인스턴스 생성
    #[frb(sync)]
    pub fn new(num_channels: usize) -> Self {
        Self {
            differential: DifferentialEngine::with_defaults(num_channels),
            ble: BleManager::new(),
            nfc: NfcReader::new(),
        }
    }

    /// 88차원 기본 엔진 생성
    #[frb(sync)]
    pub fn with_88_channels() -> Self {
        Self::new(88)
    }

    /// 896차원 완전 엔진 생성
    #[frb(sync)]
    pub fn with_896_channels() -> Self {
        Self::new(896)
    }

    /// 1792차원 궁극 엔진 생성 (Phase 5: E12-IF + 시간축 확장)
    #[frb(sync)]
    pub fn with_1792_channels() -> Self {
        Self::new(1792)
    }
}

// ============================================================================
// 차동측정 API
// ============================================================================

/// 차동측정 수행
#[frb(sync)]
pub fn differential_measure(
    s_det: Vec<f64>,
    s_ref: Vec<f64>,
    alpha: f64,
) -> Result<Vec<f64>, String> {
    let params = CorrectionParams {
        alpha,
        channel_offsets: Vec::new(),
        channel_gains: Vec::new(),
        temp_coefficient: 0.0,
    };

    let engine = DifferentialEngine::new(params, s_det.len());
    engine.measure(&s_det, &s_ref).map_err(|e| e.to_string())
}

/// 단일 채널 차동측정
#[frb(sync)]
pub fn differential_measure_single(
    s_det: f64,
    s_ref: f64,
    alpha: f64,
) -> DifferentialCorrectionDto {
    let params = CorrectionParams {
        alpha,
        ..Default::default()
    };
    let engine = DifferentialEngine::new(params, 1);
    let result = engine.measure_single(s_det, s_ref);

    DifferentialCorrectionDto {
        s_det: result.s_det,
        s_ref: result.s_ref,
        alpha: result.alpha,
        s_corrected: result.s_corrected,
    }
}

/// 차동측정 결과 DTO (Flutter 전달용)
#[frb(dart_metadata=("freezed"))]
pub struct DifferentialCorrectionDto {
    pub s_det: f64,
    pub s_ref: f64,
    pub alpha: f64,
    pub s_corrected: f64,
}

// ============================================================================
// 핑거프린트 API
// ============================================================================

/// 핑거프린트 벡터 생성 (88차원)
#[frb(sync)]
pub fn create_fingerprint_88(data: Vec<f32>) -> Result<FingerprintDto, String> {
    let fp = FingerprintVector::basic(data).map_err(|e| e.to_string())?;
    Ok(FingerprintDto {
        data: fp.data().to_vec(),
        dimension: fp.dimension(),
        measurement_type: format!("{:?}", fp.measurement_type()),
        normalized: fp.is_normalized(),
    })
}

/// 핑거프린트 벡터 생성 (896차원)
#[frb(sync)]
pub fn create_fingerprint_896(data: Vec<f32>) -> Result<FingerprintDto, String> {
    let fp = FingerprintVector::full(data).map_err(|e| e.to_string())?;
    Ok(FingerprintDto {
        data: fp.data().to_vec(),
        dimension: fp.dimension(),
        measurement_type: format!("{:?}", fp.measurement_type()),
        normalized: fp.is_normalized(),
    })
}

/// 핑거프린트 벡터 생성 (1792차원, Phase 5 궁극 확장)
#[frb(sync)]
pub fn create_fingerprint_1792(data: Vec<f32>) -> Result<FingerprintDto, String> {
    let fp = FingerprintVector::ultimate(data).map_err(|e| e.to_string())?;
    Ok(FingerprintDto {
        data: fp.data().to_vec(),
        dimension: fp.dimension(),
        measurement_type: format!("{:?}", fp.measurement_type()),
        normalized: fp.is_normalized(),
    })
}

/// 코사인 유사도 계산
#[frb(sync)]
pub fn fingerprint_cosine_similarity(
    fp1_data: Vec<f32>,
    fp2_data: Vec<f32>,
    dimension: usize,
) -> Result<f32, String> {
    let measurement_type = match dimension {
        88 => MeasurementType::Basic,
        448 => MeasurementType::Enhanced,
        896 => MeasurementType::Full,
        1792 => MeasurementType::Ultimate,
        _ => return Err("Unsupported dimension".to_string()),
    };

    let fp1 = FingerprintVector::new(fp1_data, measurement_type).map_err(|e| e.to_string())?;
    let fp2 = FingerprintVector::new(fp2_data, measurement_type).map_err(|e| e.to_string())?;

    fp1.cosine_similarity(&fp2).map_err(|e| e.to_string())
}

/// 트렌드 분석 결과 DTO (선형회귀 + 예측)
#[frb(dart_metadata=("freezed"))]
pub struct TrendAnalysisDto {
    /// 트렌드 방향: improving, worsening, stable, insufficient
    pub direction: String,
    /// 선형회귀 기울기 (단위시간당 변화량)
    pub slope: f64,
    /// 선형회귀 절편
    pub intercept: f64,
    /// 결정계수 R² (0~1)
    pub r_squared: f64,
    /// 기울기 95% 신뢰구간 [하한, 상한]
    pub slope_ci_95: Vec<f64>,
    /// 미래 예측값
    pub predictions: Vec<f64>,
    /// 이동평균 (최근 3개)
    pub moving_average: f64,
    /// 변동계수 CV%
    pub coefficient_of_variation: f64,
    /// 한국어 요약
    pub summary: String,
}

/// 핑거프린트 DTO
#[frb(dart_metadata=("freezed"))]
pub struct FingerprintDto {
    pub data: Vec<f32>,
    pub dimension: usize,
    pub measurement_type: String,
    pub normalized: bool,
}

// ============================================================================
// BLE API
// ============================================================================

/// BLE 디바이스 스캔
#[frb]
pub async fn ble_scan() -> Vec<DeviceInfoDto> {
    let ble = BleManager::new();
    let devices = ble.scan().await;
    devices
        .into_iter()
        .map(|d| DeviceInfoDto {
            device_id: d.device_id,
            name: d.name,
            rssi: d.rssi,
            state: format!("{:?}", d.state),
        })
        .collect()
}

/// BLE 디바이스 연결
#[frb]
pub async fn ble_connect(device_id: String) -> Result<bool, String> {
    let mut ble = BleManager::new();
    ble.connect(&device_id).await.map(|_| true).map_err(|e| e.to_string())
}

/// 디바이스 정보 DTO
#[frb(dart_metadata=("freezed"))]
pub struct DeviceInfoDto {
    pub device_id: String,
    pub name: String,
    pub rssi: i8,
    pub state: String,
}

// ============================================================================
// NFC API
// ============================================================================

/// 카트리지 읽기
#[frb]
pub async fn nfc_read_cartridge() -> Result<CartridgeInfoDto, String> {
    let mut nfc = NfcReader::new();
    let info = nfc.read_cartridge().await.map_err(|e| e.to_string())?;

    Ok(CartridgeInfoDto {
        cartridge_id: info.cartridge_id,
        cartridge_type: format!("{:?}", info.cartridge_type),
        lot_id: info.lot_id,
        expiry_date: info.expiry_date,
        remaining_uses: info.remaining_uses,
    })
}

/// NFC 태그 원시 데이터 파싱 (플랫폼 독립적)
///
/// Flutter에서 nfc_manager 등으로 읽은 원시 바이트를 전달받아
/// Rust의 v1.0/v2.0 자동 감지 파서로 CartridgeInfo를 추출합니다.
/// 하드웨어 없이도 NFC 태그 데이터를 파싱할 수 있습니다.
#[frb(sync)]
pub fn nfc_parse_tag(data: Vec<u8>) -> Result<CartridgeInfoDto, String> {
    let nfc = NfcReader::new();
    let info = nfc.parse_tag_data(&data).map_err(|e| e.to_string())?;

    Ok(CartridgeInfoDto {
        cartridge_id: info.cartridge_id,
        cartridge_type: format!("{:?}", info.cartridge_type),
        lot_id: info.lot_id,
        expiry_date: info.expiry_date,
        remaining_uses: info.remaining_uses,
    })
}

/// 카트리지 정보 DTO
#[frb(dart_metadata=("freezed"))]
pub struct CartridgeInfoDto {
    pub cartridge_id: String,
    pub cartridge_type: String,
    pub lot_id: String,
    pub expiry_date: String,
    pub remaining_uses: u32,
}

// ============================================================================
// 측정 파이프라인 API (BLE → DSP → AI 통합)
// ============================================================================

/// 측정 파이프라인 결과 DTO
#[frb(dart_metadata=("freezed"))]
pub struct MeasurementPipelineDto {
    pub primary_value: f64,
    pub reference_value: f64,
    pub differential_value: f64,
    pub snr: f64,
    pub confidence: f64,
    pub biomarker: String,
    pub unit: String,
    pub risk_level: String,
    pub health_score: f64,
    pub recommendations: Vec<String>,
    pub pipeline_duration_ms: f64,
}

/// AI 분석 결과 DTO
#[frb(dart_metadata=("freezed"))]
pub struct AiAnalysisResultDto {
    pub risk_level: String,
    pub health_score: f64,
    pub summary: String,
    pub recommendations: Vec<String>,
    pub trend: String,
}

/// 차동 계측 수행 + AI 안전 검증 파이프라인
#[frb]
pub fn run_measurement_pipeline(
    s_det: Vec<f64>,
    s_ref: Vec<f64>,
    alpha: f64,
    biomarker: String,
    unit: String,
) -> Result<MeasurementPipelineDto, String> {
    use manpasik_engine::ai::{InferenceEngine, ModelType, SafetyValidator};

    let start = std::time::Instant::now();

    // 1. 차동 계측
    let params = CorrectionParams {
        alpha,
        channel_offsets: Vec::new(),
        channel_gains: Vec::new(),
        temp_coefficient: 0.0,
    };
    let engine = DifferentialEngine::new(params, s_det.len());
    let corrected = engine.measure(&s_det, &s_ref).map_err(|e| e.to_string())?;

    // 2. 주요 값 추출 (평균)
    let primary_value = if corrected.is_empty() {
        0.0
    } else {
        corrected.iter().sum::<f64>() / corrected.len() as f64
    };
    let reference_value = if s_ref.is_empty() {
        0.0
    } else {
        s_ref.iter().sum::<f64>() / s_ref.len() as f64
    };

    // 3. AI 분석 (시뮬레이션 모드)
    let ai_engine = InferenceEngine::new(ModelType::ValuePredictor);
    let input: Vec<f32> = corrected.iter().take(88).map(|&v| v as f32).collect();
    let padded_input = if input.len() < 88 {
        let mut padded = input;
        padded.resize(88, 0.0);
        padded
    } else {
        input
    };
    let ai_result = ai_engine.predict(&padded_input).map_err(|e| e.to_string())?;

    // 4. 안전 검증
    let validator = SafetyValidator::new();
    let safety = validator.validate(&ai_result, primary_value, &biomarker);
    let risk_level = match safety.final_verdict {
        manpasik_engine::ai::SafetyVerdict::Normal => "normal",
        manpasik_engine::ai::SafetyVerdict::Caution => "caution",
        manpasik_engine::ai::SafetyVerdict::Alert => "warning",
        manpasik_engine::ai::SafetyVerdict::Uncertain => "uncertain",
    };

    let pipeline_duration_ms = start.elapsed().as_secs_f64() * 1000.0;

    Ok(MeasurementPipelineDto {
        primary_value,
        reference_value,
        differential_value: primary_value - reference_value,
        snr: ai_result.confidence as f64 * 50.0,
        confidence: ai_result.confidence as f64,
        biomarker,
        unit,
        risk_level: risk_level.to_string(),
        health_score: match risk_level {
            "normal" => 90.0,
            "caution" => 65.0,
            "warning" => 40.0,
            _ => 50.0,
        },
        recommendations: safety.warnings,
        pipeline_duration_ms,
    })
}

/// AI 분석만 수행 (측정값 → 위험도 판정)
#[frb]
pub fn analyze_measurement(
    value: f64,
    biomarker: String,
) -> AiAnalysisResultDto {
    use manpasik_engine::ai::{InferenceEngine, ModelType, SafetyValidator};

    let ai_engine = InferenceEngine::new(ModelType::ValuePredictor);
    let input = vec![value as f32; 88];
    let ai_result = ai_engine.predict(&input).unwrap_or(manpasik_engine::ai::InferenceResult {
        values: vec![0.0],
        confidence: 0.5,
        inference_time_ms: 0.0,
        model_type: ModelType::ValuePredictor,
    });

    let validator = SafetyValidator::new();
    let safety = validator.validate(&ai_result, value, &biomarker);

    let risk_level = match safety.final_verdict {
        manpasik_engine::ai::SafetyVerdict::Normal => "normal",
        manpasik_engine::ai::SafetyVerdict::Caution => "caution",
        manpasik_engine::ai::SafetyVerdict::Alert => "warning",
        manpasik_engine::ai::SafetyVerdict::Uncertain => "uncertain",
    };

    AiAnalysisResultDto {
        risk_level: risk_level.to_string(),
        health_score: match risk_level {
            "normal" => 90.0,
            "caution" => 65.0,
            "warning" => 40.0,
            _ => 50.0,
        },
        summary: format!("{} {:.1} — {}", biomarker, value, risk_level),
        recommendations: safety.warnings,
        trend: "stable".to_string(),
    }
}

// ============================================================================
// BLE 배터리/품질 API
// ============================================================================

/// BLE 배터리 레벨 읽기
#[frb]
pub async fn ble_read_battery(device_id: String) -> Result<u8, String> {
    let ble = BleManager::new();
    ble.read_battery_level(&device_id)
        .await
        .map_err(|e| e.to_string())
}

/// BLE 연결 품질 조회
#[frb(sync)]
pub fn ble_connection_quality(rssi: i8) -> String {
    use manpasik_engine::ble::{ConnectionQuality, RssiMonitor};
    let mut monitor = RssiMonitor::new(5);
    monitor.add_sample(rssi);
    match monitor.quality() {
        ConnectionQuality::Excellent => "excellent".to_string(),
        ConnectionQuality::Good => "good".to_string(),
        ConnectionQuality::Fair => "fair".to_string(),
        ConnectionQuality::Poor => "poor".to_string(),
    }
}

// ============================================================================
// 유틸리티 API
// ============================================================================

/// 엔진 버전 조회
#[frb(sync)]
pub fn get_engine_version() -> String {
    manpasik_engine::VERSION.to_string()
}

/// 최대 채널 수 조회
#[frb(sync)]
pub fn get_max_channels() -> usize {
    manpasik_engine::MAX_CHANNELS
}

/// 시계열 트렌드 분석 (선형회귀 + 신뢰구간 + 예측)
#[frb(sync)]
pub fn analyze_trend(
    values: Vec<f64>,
    biomarker: String,
    predict_ahead: usize,
) -> TrendAnalysisDto {
    use manpasik_engine::ai::{TrendAnalyzer, TrendDirection};
    let analyzer = TrendAnalyzer::new();
    let result = analyzer.analyze(&values, &biomarker, predict_ahead);

    let direction = match result.direction {
        TrendDirection::Improving => "improving",
        TrendDirection::Worsening => "worsening",
        TrendDirection::Stable => "stable",
        TrendDirection::Insufficient => "insufficient",
    };

    let (slope, intercept, r_squared, slope_ci_low, slope_ci_high) =
        if let Some(ref reg) = result.regression {
            (reg.slope, reg.intercept, reg.r_squared, reg.slope_ci_95[0], reg.slope_ci_95[1])
        } else {
            (0.0, 0.0, 0.0, 0.0, 0.0)
        };

    TrendAnalysisDto {
        direction: direction.to_string(),
        slope,
        intercept,
        r_squared,
        slope_ci_95: vec![slope_ci_low, slope_ci_high],
        predictions: result.predictions,
        moving_average: result.moving_average,
        coefficient_of_variation: result.coefficient_of_variation,
        summary: result.summary,
    }
}

/// SHA-256 해시 계산
#[frb(sync)]
pub fn calculate_sha256(data: Vec<u8>) -> String {
    use manpasik_engine::crypto::CryptoEngine;
    let crypto = CryptoEngine::new();
    crypto.hash_sha256(&data)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_differential_single() {
        let result = differential_measure_single(1.234, 0.012, 0.95);
        assert!((result.s_corrected - 1.2226).abs() < 0.0001);
    }

    #[test]
    fn test_fingerprint_88() {
        let data = vec![0.5f32; 88];
        let fp = create_fingerprint_88(data).unwrap();
        assert_eq!(fp.dimension, 88);
    }

    #[test]
    fn test_fingerprint_1792() {
        let data = vec![0.5f32; 1792];
        let fp = create_fingerprint_1792(data).unwrap();
        assert_eq!(fp.dimension, 1792);
        assert_eq!(fp.measurement_type, "Ultimate");
    }

    #[test]
    fn test_max_channels_1792() {
        assert_eq!(get_max_channels(), 1792);
    }

    #[test]
    fn test_engine_version() {
        let version = get_engine_version();
        assert!(!version.is_empty());
    }

    #[test]
    fn test_measurement_pipeline() {
        let s_det = vec![1.0f64; 88];
        let s_ref = vec![0.01f64; 88];
        let result = run_measurement_pipeline(
            s_det,
            s_ref,
            0.95,
            "glucose".to_string(),
            "mg/dL".to_string(),
        );
        assert!(result.is_ok());
        let r = result.unwrap();
        assert!(r.confidence > 0.0);
        assert!(r.pipeline_duration_ms >= 0.0);
        assert!(!r.risk_level.is_empty());
    }

    #[test]
    fn test_analyze_measurement() {
        let result = analyze_measurement(85.0, "glucose".to_string());
        assert_eq!(result.risk_level, "normal");
        assert!(result.health_score > 0.0);
    }

    #[test]
    fn test_ble_connection_quality() {
        assert_eq!(ble_connection_quality(-45), "excellent");
        assert_eq!(ble_connection_quality(-65), "good");
        assert_eq!(ble_connection_quality(-80), "fair");
        assert_eq!(ble_connection_quality(-95), "poor");
    }

    #[test]
    fn test_nfc_parse_tag_v1() {
        // v1.0 태그 데이터 생성 (64 바이트)
        let mut data = vec![0u8; 64];
        // [0-7] 카트리��� ID
        data[0..8].copy_from_slice(&[0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08]);
        // [8] 카트리지 타입 (Glucose = 0x01)
        data[8] = 0x01;
        // [9-16] 로트 ID
        data[9..17].copy_from_slice(b"LOT2026A");
        // [17-24] 유효 기간
        data[17..25].copy_from_slice(b"20270630");
        // [25-26] 잔여 사용 횟수 (u16 LE = 10)
        data[25..27].copy_from_slice(&10u16.to_le_bytes());
        // [27-28] 최대 사용 횟수 (u16 LE = 50)
        data[27..29].copy_from_slice(&50u16.to_le_bytes());
        // [29-36] alpha 계수 (f64 LE = 0.98)
        data[29..37].copy_from_slice(&0.98f64.to_le_bytes());
        // [37-44] 온도 보정 계수
        data[37..45].copy_from_slice(&0.01f64.to_le_bytes());
        // [45-52] 습도 보정 계수
        data[45..53].copy_from_slice(&0.005f64.to_le_bytes());

        let result = nfc_parse_tag(data);
        assert!(result.is_ok(), "v1.0 태그 파싱 성공: {:?}", result.err());
        let info = result.unwrap();
        assert_eq!(info.remaining_uses, 10);
    }

    #[test]
    fn test_nfc_parse_tag_too_short() {
        // 너무 짧은 데이터 → 에러
        let data = vec![0u8; 10];
        let result = nfc_parse_tag(data);
        assert!(result.is_err());
    }

    #[test]
    fn test_sha256_hash() {
        let hash = calculate_sha256(b"ManPaSik".to_vec());
        assert_eq!(hash.len(), 64); // SHA-256 hex = 64자
        assert!(!hash.is_empty());
    }

    #[test]
    fn test_differential_measure_vec() {
        let s_det = vec![1.0f64; 16];
        let s_ref = vec![0.01f64; 16];
        let result = differential_measure(s_det, s_ref, 0.98);
        assert!(result.is_ok());
        let corrected = result.unwrap();
        assert_eq!(corrected.len(), 16);
        // S_corrected = 1.0 - 0.98 * 0.01 = 0.9902
        assert!((corrected[0] - 0.9902).abs() < 0.001);
    }

    #[test]
    fn test_analyze_trend_improving() {
        let values = vec![120.0, 115.0, 110.0, 105.0, 100.0];
        let result = analyze_trend(values, "glucose".to_string(), 2);
        assert_eq!(result.direction, "improving");
        assert!(result.r_squared > 0.95);
        assert_eq!(result.predictions.len(), 2);
        assert!(result.predictions[0] < 100.0);
    }

    #[test]
    fn test_analyze_trend_insufficient() {
        let values = vec![100.0];
        let result = analyze_trend(values, "glucose".to_string(), 0);
        assert_eq!(result.direction, "insufficient");
    }

    #[test]
    fn test_cosine_similarity_identical() {
        let fp = vec![1.0f32, 0.0, 0.5];
        // 동일 벡터의 코사인 유사도 = 1.0
        // 88차원으로 맞춤
        let data = vec![0.5f32; 88];
        let similarity =
            fingerprint_cosine_similarity(data.clone(), data.clone(), 88);
        assert!(similarity.is_ok());
        let sim = similarity.unwrap();
        assert!((sim - 1.0).abs() < 0.001, "동일 벡터 유사도 = 1.0");
    }
}
