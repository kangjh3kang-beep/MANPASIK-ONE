//! ManPaSik Core Engine
//!
//! 만파식(萬波息) - 차동측정 기반 범용 분석 엔진
//!
//! # 핵심 모듈
//!
//! - `differential` - 차동측정 엔진 (S_det - α × S_ref)
//! - `ai` - TFLite 기반 엣지 AI 추론
//! - `ble` - BLE 5.0 GATT 통신
//! - `nfc` - NFC 카트리지 인식
//! - `dsp` - 실시간 신호 처리
//! - `crypto` - 암호화 (AES-256, SHA-256)
//! - `fingerprint` - 88→448→896→1792차원 핑거프린트 생성
//! - `sync` - CRDT 기반 오프라인 동기화

#![warn(clippy::all)]
#![warn(rust_2018_idioms)]

pub mod ai;
pub mod ble;
pub mod crypto;
pub mod differential;
pub mod dsp;
pub mod fingerprint;
pub mod nfc;
pub mod sync;

// Re-exports
pub use differential::{CorrectionParams, DifferentialEngine, MeasurementResult};
pub use fingerprint::FingerprintVector;

/// 엔진 버전
pub const VERSION: &str = env!("CARGO_PKG_VERSION");

/// 최대 센서 채널 수 (1792차원 = 896 × 2 시간 윈도우, Phase 5 궁극 확장)
/// Phase 1-4: 88→448→896, Phase 5: 1792
pub const MAX_CHANNELS: usize = 1792;

/// 차동측정 기본 알파 값
pub const DEFAULT_ALPHA: f64 = 0.95;

/// 측정 데이터 패킷
#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct MeasurementPacket {
    /// 헤더 정보
    pub header: PacketHeader,
    /// 페이로드
    pub payload: PacketPayload,
    /// 푸터 (체크섬 등)
    pub footer: PacketFooter,
}

#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct PacketHeader {
    /// 리더기 고유 ID
    pub device_id: String,
    /// 제조 로트 ID
    pub lot_id: String,
    /// 펌웨어 버전
    pub fw_ver: String,
    /// 카트리지 ID
    pub cartridge_id: String,
    /// 측정 세션 ID
    pub session_id: uuid::Uuid,
    /// 타임스탬프
    pub timestamp: chrono::DateTime<chrono::Utc>,
}

#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct PacketPayload {
    /// 원시 채널 데이터
    pub raw_channels: Vec<f64>,
    /// 측정 결과
    pub result: Option<MeasurementResult>,
    /// 환경 메타데이터
    pub env_meta: EnvironmentMeta,
    /// 상태 메타데이터
    pub state_meta: StateMeta,
}

#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct EnvironmentMeta {
    /// 온도 (섭씨)
    pub temp_c: f32,
    /// 습도 (%)
    pub humidity_pct: f32,
    /// 기압 (kPa)
    pub pressure_kpa: f32,
}

#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct StateMeta {
    /// 배터리 잔량 (%)
    pub battery_pct: u8,
    /// 신호 품질
    pub signal_quality: SignalQuality,
    /// 자가 진단 결과
    pub self_diagnostic: DiagnosticResult,
    /// 불확도
    pub uncertainty: f64,
}

#[derive(Debug, Clone, Copy, serde::Serialize, serde::Deserialize)]
pub enum SignalQuality {
    High,
    Medium,
    Low,
    Poor,
}

#[derive(Debug, Clone, Copy, serde::Serialize, serde::Deserialize)]
pub enum DiagnosticResult {
    Pass,
    Warning,
    Fail,
}

#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct PacketFooter {
    /// SHA-256 체크섬
    pub checksum: String,
    /// 스키마 버전
    pub schema_ver: String,
    /// 변환 로그 (해시체인)
    pub transform_log: Vec<TransformStep>,
}

#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct TransformStep {
    pub step: String,
    pub hash: String,
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_version() {
        assert!(!VERSION.is_empty());
    }

    #[test]
    fn test_max_channels() {
        assert_eq!(MAX_CHANNELS, 1792);
    }
}
