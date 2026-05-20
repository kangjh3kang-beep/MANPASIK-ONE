//! ManPaSik AI 코어 엔진
//!
//! 카트리지 비전 분석 및 바이오마커 추론을 위한 네이티브 Rust 라이브러리.
//! flutter_rust_bridge를 통해 Flutter에서 FFI 호출.

// ============================================================================
// SSOT 상수 (CLAUDE.md 베이스라인 — 절대 변경 금지)
// ============================================================================

/// 엔진 버전 (Cargo.toml 기준)
pub const VERSION: &str = env!("CARGO_PKG_VERSION");

/// 차동 보정 기본 알파 계수 (S_diff = S_det - alpha * S_ref)
/// SSOT 범위: 0.90 ~ 1.10
pub const DEFAULT_ALPHA: f64 = 0.98;

/// 알파 허용 최솟값
pub const ALPHA_MIN: f64 = 0.90;

/// 알파 허용 최댓값
pub const ALPHA_MAX: f64 = 1.10;

/// 최대 채널 수 (1792차원 = 896 x 2 시간축 확장)
pub const MAX_CHANNELS: usize = 1792;

/// CSI v1.0 커넥터 핀 수 (Samtec MECF-08-01-L-DV, 1.27mm)
pub const CSI_PIN_COUNT: usize = 16;

/// PPM 물리 치수 (mm)
pub const PPM_DIMENSIONS: [f64; 3] = [49.70, 30.0, 4.30];

// ============================================================================
// 모듈
// ============================================================================

pub mod ai;
pub mod api;
pub mod ble;
pub mod crypto;
pub mod differential;
pub mod dsp;
pub mod fingerprint;
pub mod nfc;
pub mod sync;
