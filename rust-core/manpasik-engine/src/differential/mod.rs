//! 차동측정 엔진 (Differential Measurement Engine)
//!
//! 핵심 공식: S_corrected = S_det - α × S_ref
//!
//! - S_det: 검출 전극 신호
//! - S_ref: 기준 전극 신호
//! - α: 보정 계수 (기본값 0.95)
//!
//! 이 방식으로 99%의 매트릭스 노이즈를 제거합니다.

use serde::{Deserialize, Serialize};
use thiserror::Error;

/// 차동측정 엔진
pub struct DifferentialEngine {
    /// 보정 파라미터
    params: CorrectionParams,
    /// 채널 수
    num_channels: usize,
}

/// 보정 파라미터
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CorrectionParams {
    /// 알파 계수 (기본값 0.95)
    pub alpha: f64,
    /// 채널별 오프셋
    pub channel_offsets: Vec<f64>,
    /// 채널별 게인
    pub channel_gains: Vec<f64>,
    /// 온도 보정 계수
    pub temp_coefficient: f64,
}

impl Default for CorrectionParams {
    fn default() -> Self {
        Self {
            alpha: crate::DEFAULT_ALPHA,
            channel_offsets: Vec::new(),
            channel_gains: Vec::new(),
            temp_coefficient: 0.0,
        }
    }
}

/// 측정 결과
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MeasurementResult {
    /// 주요 측정값
    pub primary_value: f64,
    /// 단위
    pub unit: String,
    /// 신뢰도 (0.0 ~ 1.0)
    pub confidence: f64,
    /// 차동 보정 상세
    pub differential_correction: DifferentialCorrection,
}

/// 차동 보정 상세
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DifferentialCorrection {
    /// 검출 전극 신호
    pub s_det: f64,
    /// 기준 전극 신호
    pub s_ref: f64,
    /// 알파 계수
    pub alpha: f64,
    /// 보정된 신호
    pub s_corrected: f64,
}

#[derive(Debug, Error)]
pub enum DifferentialError {
    #[error("채널 수 불일치: expected {expected}, got {got}")]
    ChannelCountMismatch { expected: usize, got: usize },

    #[error("유효하지 않은 측정값: {0}")]
    InvalidMeasurement(String),

    #[error("보정 파라미터 오류: {0}")]
    CalibrationError(String),
}

impl DifferentialEngine {
    /// 새 엔진 생성
    pub fn new(params: CorrectionParams, num_channels: usize) -> Self {
        Self {
            params,
            num_channels,
        }
    }

    /// 기본 설정으로 엔진 생성
    pub fn with_defaults(num_channels: usize) -> Self {
        Self::new(CorrectionParams::default(), num_channels)
    }

    /// 차동측정 수행
    ///
    /// # Arguments
    ///
    /// * `s_det` - 검출 전극 신호 배열
    /// * `s_ref` - 기준 전극 신호 배열
    ///
    /// # Returns
    ///
    /// 보정된 신호 배열
    pub fn measure(&self, s_det: &[f64], s_ref: &[f64]) -> Result<Vec<f64>, DifferentialError> {
        // 채널 수 검증
        if s_det.len() != self.num_channels {
            return Err(DifferentialError::ChannelCountMismatch {
                expected: self.num_channels,
                got: s_det.len(),
            });
        }
        if s_ref.len() != self.num_channels {
            return Err(DifferentialError::ChannelCountMismatch {
                expected: self.num_channels,
                got: s_ref.len(),
            });
        }

        // 차동측정: S_corrected = S_det - α × S_ref
        let corrected: Vec<f64> = s_det
            .iter()
            .zip(s_ref.iter())
            .enumerate()
            .map(|(i, (&det, &ref_val))| {
                let offset = self.params.channel_offsets.get(i).unwrap_or(&0.0);
                let gain = self.params.channel_gains.get(i).unwrap_or(&1.0);

                let raw = det - self.params.alpha * ref_val;
                (raw - offset) * gain
            })
            .collect();

        Ok(corrected)
    }

    /// 단일 채널 측정
    pub fn measure_single(&self, s_det: f64, s_ref: f64) -> DifferentialCorrection {
        let s_corrected = s_det - self.params.alpha * s_ref;

        DifferentialCorrection {
            s_det,
            s_ref,
            alpha: self.params.alpha,
            s_corrected,
        }
    }

    /// 보정 파라미터 업데이트
    pub fn update_params(&mut self, params: CorrectionParams) {
        self.params = params;
    }

    /// 현재 알파 값 조회
    pub fn alpha(&self) -> f64 {
        self.params.alpha
    }

    /// 알파 값 설정
    pub fn set_alpha(&mut self, alpha: f64) {
        self.params.alpha = alpha;
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_differential_measurement() {
        let engine = DifferentialEngine::with_defaults(4);

        let s_det = vec![1.0, 2.0, 3.0, 4.0];
        let s_ref = vec![0.1, 0.2, 0.3, 0.4];

        let result = engine.measure(&s_det, &s_ref).unwrap();

        // S_corrected = S_det - 0.95 * S_ref
        assert!((result[0] - 0.905).abs() < 0.001);
        assert!((result[1] - 1.810).abs() < 0.001);
        assert!((result[2] - 2.715).abs() < 0.001);
        assert!((result[3] - 3.620).abs() < 0.001);
    }

    #[test]
    fn test_single_measurement() {
        let engine = DifferentialEngine::with_defaults(1);

        let result = engine.measure_single(1.234, 0.012);

        assert!((result.s_corrected - 1.2226).abs() < 0.0001);
        assert_eq!(result.alpha, 0.95);
    }

    #[test]
    fn test_channel_mismatch() {
        let engine = DifferentialEngine::with_defaults(4);

        let s_det = vec![1.0, 2.0]; // 2개 (4개 필요)
        let s_ref = vec![0.1, 0.2, 0.3, 0.4];

        let result = engine.measure(&s_det, &s_ref);
        assert!(result.is_err());
    }
}
