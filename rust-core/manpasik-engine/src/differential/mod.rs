//! 차동측정 엔진 (Differential Measurement Engine)
//!
//! 핵심 공식: S_corrected = S_det - α × S_ref
//!
//! - S_det: 검출 전극 신호
//! - S_ref: 기준 전극 신호
//! - α: 보정 계수 (SSOT 기본값 0.98, 범위 0.90~1.10)
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
    /// 알파 계수 (SSOT 기본값 0.98, 범위 0.90~1.10)
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

    /// 기본 설정으로 엔진 생성 (alpha = SSOT 0.98)
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

    /// 알파 값 설정 (SSOT 범위 0.90~1.10 클램프)
    pub fn set_alpha(&mut self, alpha: f64) {
        self.params.alpha = alpha.clamp(crate::ALPHA_MIN, crate::ALPHA_MAX);
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_default_alpha_is_ssot() {
        let params = CorrectionParams::default();
        assert!((params.alpha - 0.98).abs() < f64::EPSILON);
    }

    #[test]
    fn test_differential_measurement() {
        let engine = DifferentialEngine::with_defaults(4);

        let s_det = vec![1.0, 2.0, 3.0, 4.0];
        let s_ref = vec![0.1, 0.2, 0.3, 0.4];

        let result = engine.measure(&s_det, &s_ref).unwrap();

        // S_corrected = S_det - 0.98 * S_ref
        assert!((result[0] - 0.902).abs() < 0.001);
        assert!((result[1] - 1.804).abs() < 0.001);
        assert!((result[2] - 2.706).abs() < 0.001);
        assert!((result[3] - 3.608).abs() < 0.001);
    }

    #[test]
    fn test_single_measurement() {
        let engine = DifferentialEngine::with_defaults(1);

        let result = engine.measure_single(1.234, 0.012);

        // 1.234 - 0.98 * 0.012 = 1.234 - 0.01176 = 1.22224
        assert!((result.s_corrected - 1.22224).abs() < 0.0001);
        assert!((result.alpha - 0.98).abs() < f64::EPSILON);
    }

    #[test]
    fn test_channel_mismatch() {
        let engine = DifferentialEngine::with_defaults(4);

        let s_det = vec![1.0, 2.0]; // 2개 (4개 필요)
        let s_ref = vec![0.1, 0.2, 0.3, 0.4];

        let result = engine.measure(&s_det, &s_ref);
        assert!(result.is_err());
    }

    #[test]
    fn test_alpha_range_clamp() {
        let mut engine = DifferentialEngine::with_defaults(1);

        engine.set_alpha(0.80); // 아래 범위
        assert!((engine.alpha() - 0.90).abs() < f64::EPSILON);

        engine.set_alpha(1.20); // 위 범위
        assert!((engine.alpha() - 1.10).abs() < f64::EPSILON);

        engine.set_alpha(1.05); // 범위 내
        assert!((engine.alpha() - 1.05).abs() < f64::EPSILON);
    }

    // ========================================================================
    // Phase M 엣지 케이스 테스트 (+15건)
    // ========================================================================

    #[test]
    fn test_zero_signal_input() {
        // 0 입력 → 0 출력
        let engine = DifferentialEngine::with_defaults(4);
        let s_det = vec![0.0, 0.0, 0.0, 0.0];
        let s_ref = vec![0.0, 0.0, 0.0, 0.0];

        let result = engine.measure(&s_det, &s_ref).unwrap();
        for v in result {
            assert!(v.abs() < f64::EPSILON);
        }
    }

    #[test]
    fn test_negative_signal_input() {
        // 음수 입력도 처리되어야 함
        let engine = DifferentialEngine::with_defaults(2);
        let s_det = vec![-1.0, -2.0];
        let s_ref = vec![0.1, 0.2];

        let result = engine.measure(&s_det, &s_ref).unwrap();
        // -1.0 - 0.98 * 0.1 = -1.098
        assert!((result[0] - (-1.098)).abs() < 0.001);
    }

    #[test]
    fn test_alpha_min_boundary() {
        let mut engine = DifferentialEngine::with_defaults(1);
        engine.set_alpha(crate::ALPHA_MIN);
        assert!((engine.alpha() - crate::ALPHA_MIN).abs() < f64::EPSILON);
    }

    #[test]
    fn test_alpha_max_boundary() {
        let mut engine = DifferentialEngine::with_defaults(1);
        engine.set_alpha(crate::ALPHA_MAX);
        assert!((engine.alpha() - crate::ALPHA_MAX).abs() < f64::EPSILON);
    }

    #[test]
    fn test_alpha_just_below_min() {
        let mut engine = DifferentialEngine::with_defaults(1);
        engine.set_alpha(crate::ALPHA_MIN - 0.001);
        // 클램프 → ALPHA_MIN
        assert!((engine.alpha() - crate::ALPHA_MIN).abs() < f64::EPSILON);
    }

    #[test]
    fn test_alpha_just_above_max() {
        let mut engine = DifferentialEngine::with_defaults(1);
        engine.set_alpha(crate::ALPHA_MAX + 0.001);
        assert!((engine.alpha() - crate::ALPHA_MAX).abs() < f64::EPSILON);
    }

    #[test]
    fn test_extreme_negative_alpha() {
        let mut engine = DifferentialEngine::with_defaults(1);
        engine.set_alpha(-100.0);
        assert!((engine.alpha() - crate::ALPHA_MIN).abs() < f64::EPSILON);
    }

    #[test]
    fn test_extreme_positive_alpha() {
        let mut engine = DifferentialEngine::with_defaults(1);
        engine.set_alpha(100.0);
        assert!((engine.alpha() - crate::ALPHA_MAX).abs() < f64::EPSILON);
    }

    #[test]
    fn test_nan_in_signal() {
        // NaN 입력은 결과에도 NaN이 나올 수 있음 (구현에 따라)
        let engine = DifferentialEngine::with_defaults(2);
        let s_det = vec![1.0, f64::NAN];
        let s_ref = vec![0.1, 0.2];

        let result = engine.measure(&s_det, &s_ref).unwrap();
        assert!(result[1].is_nan(), "NaN 입력은 NaN 출력");
    }

    #[test]
    fn test_infinity_signal() {
        let engine = DifferentialEngine::with_defaults(2);
        let s_det = vec![f64::INFINITY, 1.0];
        let s_ref = vec![0.1, 0.2];

        let result = engine.measure(&s_det, &s_ref).unwrap();
        assert!(result[0].is_infinite(), "Infinity 입력은 Infinity 출력");
    }

    #[test]
    fn test_large_channel_count() {
        let n = 1000;
        let engine = DifferentialEngine::with_defaults(n);
        let s_det = vec![1.0; n];
        let s_ref = vec![0.5; n];

        let result = engine.measure(&s_det, &s_ref).unwrap();
        assert_eq!(result.len(), n);
        // 1.0 - 0.98 * 0.5 = 0.51
        for v in result {
            assert!((v - 0.51).abs() < 0.001);
        }
    }

    #[test]
    fn test_single_channel_measurement() {
        let engine = DifferentialEngine::with_defaults(1);
        let s_det = vec![5.0];
        let s_ref = vec![1.0];

        let result = engine.measure(&s_det, &s_ref).unwrap();
        assert_eq!(result.len(), 1);
        // 5.0 - 0.98 * 1.0 = 4.02
        assert!((result[0] - 4.02).abs() < 0.001);
    }

    #[test]
    fn test_empty_signal_arrays() {
        let engine = DifferentialEngine::with_defaults(0);
        let s_det: Vec<f64> = vec![];
        let s_ref: Vec<f64> = vec![];

        let result = engine.measure(&s_det, &s_ref).unwrap();
        assert!(result.is_empty());
    }

    #[test]
    fn test_alpha_default_is_098() {
        let engine = DifferentialEngine::with_defaults(1);
        assert!((engine.alpha() - 0.98).abs() < f64::EPSILON);
    }

    #[test]
    fn test_correction_params_clone() {
        let params = CorrectionParams::default();
        let cloned = params.clone();
        assert!((params.alpha - cloned.alpha).abs() < f64::EPSILON);
    }
}
