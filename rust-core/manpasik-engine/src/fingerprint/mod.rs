//! 1792차원 핑거프린트 벡터 생성 모듈
//!
//! 88차원 기본 측정 → 448차원 전자코/전자혀 → 896차원 완전 융합 → 1792차원 궁극 확장
//!
//! 원본 기획안: "88→448→896→1792차원, 생물체 단계 비유"
//! - 88차원: 단세포(기본 바이오마커)
//! - 448차원: 다세포(전자코/전자혀 융합)
//! - 896차원: 유기체(완전 센서 융합)
//! - 1792차원: 생태계(E12-IF 다중 리더기 융합 + 시간축 확장) — Phase 5

use serde::{Deserialize, Serialize};
use thiserror::Error;

/// 핑거프린트 차원 상수
pub const DIM_88: usize = 88;
pub const DIM_448: usize = 448;
pub const DIM_896: usize = 896;
/// 궁극 확장 차원 (Phase 5: E12-IF 다중 리더기 융합 + 시간축 확장)
/// 896차원 × 2 시간 윈도우 = 1792차원
pub const DIM_1792: usize = 1792;

/// 핑거프린트 벡터
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct FingerprintVector {
    /// 벡터 데이터
    data: Vec<f32>,
    /// 벡터 차원
    dimension: usize,
    /// 측정 타입
    measurement_type: MeasurementType,
    /// 정규화 여부
    normalized: bool,
}

/// 측정 타입 (기획안 4단계 성장 경로)
#[derive(Debug, Clone, Copy, Serialize, Deserialize, PartialEq, Eq)]
pub enum MeasurementType {
    /// 기본 측정 (88차원) — 단세포 단계
    Basic,
    /// 전자코/전자혀 포함 (448차원) — 다세포 단계
    Enhanced,
    /// 완전 융합 (896차원) — 유기체 단계
    Full,
    /// 궁극 확장 (1792차원) — 생태계 단계 (Phase 5: E12-IF + 시간축)
    Ultimate,
}

#[derive(Debug, Error)]
pub enum FingerprintError {
    #[error("차원 불일치: expected {expected}, got {got}")]
    DimensionMismatch { expected: usize, got: usize },

    #[error("빈 벡터")]
    EmptyVector,

    #[error("정규화 오류: {0}")]
    NormalizationError(String),
}

impl FingerprintVector {
    /// 새 핑거프린트 벡터 생성
    pub fn new(
        data: Vec<f32>,
        measurement_type: MeasurementType,
    ) -> Result<Self, FingerprintError> {
        if data.is_empty() {
            return Err(FingerprintError::EmptyVector);
        }

        let expected_dim = match measurement_type {
            MeasurementType::Basic => DIM_88,
            MeasurementType::Enhanced => DIM_448,
            MeasurementType::Full => DIM_896,
            MeasurementType::Ultimate => DIM_1792,
        };

        if data.len() != expected_dim {
            return Err(FingerprintError::DimensionMismatch {
                expected: expected_dim,
                got: data.len(),
            });
        }

        Ok(Self {
            data,
            dimension: expected_dim,
            measurement_type,
            normalized: false,
        })
    }

    /// 88차원 기본 벡터 생성
    pub fn basic(data: Vec<f32>) -> Result<Self, FingerprintError> {
        Self::new(data, MeasurementType::Basic)
    }

    /// 448차원 확장 벡터 생성
    pub fn enhanced(data: Vec<f32>) -> Result<Self, FingerprintError> {
        Self::new(data, MeasurementType::Enhanced)
    }

    /// 896차원 완전 벡터 생성
    pub fn full(data: Vec<f32>) -> Result<Self, FingerprintError> {
        Self::new(data, MeasurementType::Full)
    }

    /// 1792차원 궁극 벡터 생성 (Phase 5: E12-IF + 시간축 확장)
    pub fn ultimate(data: Vec<f32>) -> Result<Self, FingerprintError> {
        Self::new(data, MeasurementType::Ultimate)
    }

    /// L2 정규화
    pub fn normalize(&mut self) -> Result<(), FingerprintError> {
        let norm: f32 = self.data.iter().map(|x| x * x).sum::<f32>().sqrt();

        if norm == 0.0 {
            return Err(FingerprintError::NormalizationError(
                "Zero norm vector".to_string(),
            ));
        }

        for x in &mut self.data {
            *x /= norm;
        }

        self.normalized = true;
        Ok(())
    }

    /// 코사인 유사도 계산
    pub fn cosine_similarity(&self, other: &FingerprintVector) -> Result<f32, FingerprintError> {
        if self.dimension != other.dimension {
            return Err(FingerprintError::DimensionMismatch {
                expected: self.dimension,
                got: other.dimension,
            });
        }

        let dot: f32 = self
            .data
            .iter()
            .zip(other.data.iter())
            .map(|(a, b)| a * b)
            .sum();

        let norm_a: f32 = self.data.iter().map(|x| x * x).sum::<f32>().sqrt();
        let norm_b: f32 = other.data.iter().map(|x| x * x).sum::<f32>().sqrt();

        if norm_a == 0.0 || norm_b == 0.0 {
            return Ok(0.0);
        }

        Ok(dot / (norm_a * norm_b))
    }

    /// 유클리드 거리 계산
    pub fn euclidean_distance(&self, other: &FingerprintVector) -> Result<f32, FingerprintError> {
        if self.dimension != other.dimension {
            return Err(FingerprintError::DimensionMismatch {
                expected: self.dimension,
                got: other.dimension,
            });
        }

        let sum: f32 = self
            .data
            .iter()
            .zip(other.data.iter())
            .map(|(a, b)| (a - b).powi(2))
            .sum();

        Ok(sum.sqrt())
    }

    /// 벡터 데이터 조회
    pub fn data(&self) -> &[f32] {
        &self.data
    }

    /// 차원 조회
    pub fn dimension(&self) -> usize {
        self.dimension
    }

    /// 측정 타입 조회
    pub fn measurement_type(&self) -> MeasurementType {
        self.measurement_type
    }

    /// 정규화 여부 조회
    pub fn is_normalized(&self) -> bool {
        self.normalized
    }

    /// Milvus 저장용 벡터 변환
    pub fn to_milvus_vector(&self) -> Vec<f32> {
        self.data.clone()
    }
}

/// 핑거프린트 빌더 (88 → 448 → 896 → 1792 확장)
pub struct FingerprintBuilder {
    base_channels: Vec<f32>,
    e_nose_channels: Option<Vec<f32>>,
    e_tongue_channels: Option<Vec<f32>>,
    /// Phase 5: 시간축 확장 데이터 (이전 시점의 896차원 벡터)
    temporal_channels: Option<Vec<f32>>,
}

impl FingerprintBuilder {
    /// 88차원 기본 채널로 시작
    pub fn new(base_channels: Vec<f32>) -> Self {
        Self {
            base_channels,
            e_nose_channels: None,
            e_tongue_channels: None,
            temporal_channels: None,
        }
    }

    /// 전자코 채널 추가 (8채널 × 복수 시점)
    pub fn with_e_nose(mut self, channels: Vec<f32>) -> Self {
        self.e_nose_channels = Some(channels);
        self
    }

    /// 전자혀 채널 추가 (8채널 × 복수 시점)
    pub fn with_e_tongue(mut self, channels: Vec<f32>) -> Self {
        self.e_tongue_channels = Some(channels);
        self
    }

    /// 시간축 확장 데이터 추가 (Phase 5: 이전 시점 896차원 벡터)
    /// 현재 896차원 + 이전 시점 896차원 = 1792차원 궁극 핑거프린트
    pub fn with_temporal(mut self, previous_full_vector: Vec<f32>) -> Self {
        self.temporal_channels = Some(previous_full_vector);
        self
    }

    /// 핑거프린트 벡터 빌드
    pub fn build(self) -> Result<FingerprintVector, FingerprintError> {
        let mut data = self.base_channels;

        // 전자코/전자혀 융합
        if let Some(e_nose) = self.e_nose_channels {
            data.extend(e_nose);
        }
        if let Some(e_tongue) = self.e_tongue_channels {
            data.extend(e_tongue);
        }

        // 시간축 확장 (Phase 5: 896 + 896 → 1792)
        if let Some(temporal) = self.temporal_channels {
            data.extend(temporal);
        }

        // 차원에 따라 타입 결정 (4단계 성장 경로)
        let measurement_type = match data.len() {
            DIM_88 => MeasurementType::Basic,
            DIM_448 => MeasurementType::Enhanced,
            DIM_896 => MeasurementType::Full,
            DIM_1792 => MeasurementType::Ultimate,
            _ => {
                return Err(FingerprintError::DimensionMismatch {
                    expected: DIM_1792,
                    got: data.len(),
                })
            }
        };

        FingerprintVector::new(data, measurement_type)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_basic_fingerprint() {
        let data = vec![0.1f32; DIM_88];
        let fp = FingerprintVector::basic(data).unwrap();

        assert_eq!(fp.dimension(), DIM_88);
        assert_eq!(fp.measurement_type(), MeasurementType::Basic);
    }

    #[test]
    fn test_cosine_similarity() {
        let data1 = vec![1.0f32; DIM_88];
        let data2 = vec![1.0f32; DIM_88];

        let fp1 = FingerprintVector::basic(data1).unwrap();
        let fp2 = FingerprintVector::basic(data2).unwrap();

        let similarity = fp1.cosine_similarity(&fp2).unwrap();
        assert!((similarity - 1.0).abs() < 0.001);
    }

    #[test]
    fn test_ultimate_fingerprint() {
        let data = vec![0.1f32; DIM_1792];
        let fp = FingerprintVector::ultimate(data).unwrap();

        assert_eq!(fp.dimension(), DIM_1792);
        assert_eq!(fp.measurement_type(), MeasurementType::Ultimate);
    }

    #[test]
    fn test_full_growth_path_88_448_896_1792() {
        // 기획안: "88→448→896→1792차원, 생물체 단계 비유"
        let basic = FingerprintVector::basic(vec![0.1f32; DIM_88]).unwrap();
        assert_eq!(basic.dimension(), 88);
        assert_eq!(basic.measurement_type(), MeasurementType::Basic);

        let enhanced = FingerprintVector::enhanced(vec![0.1f32; DIM_448]).unwrap();
        assert_eq!(enhanced.dimension(), 448);
        assert_eq!(enhanced.measurement_type(), MeasurementType::Enhanced);

        let full = FingerprintVector::full(vec![0.1f32; DIM_896]).unwrap();
        assert_eq!(full.dimension(), 896);
        assert_eq!(full.measurement_type(), MeasurementType::Full);

        let ultimate = FingerprintVector::ultimate(vec![0.1f32; DIM_1792]).unwrap();
        assert_eq!(ultimate.dimension(), 1792);
        assert_eq!(ultimate.measurement_type(), MeasurementType::Ultimate);
    }

    #[test]
    fn test_builder_1792_temporal_expansion() {
        // Phase 5: 896차원 현재 + 896차원 이전 시점 → 1792차원
        // 빌더 경로: 896(현재 완전 융합) + 896(이전 시점) = 1792(궁극)
        let current_896 = vec![0.1f32; 896];
        let previous_896 = vec![0.2f32; 896];

        let fp = FingerprintBuilder::new(current_896)
            .with_temporal(previous_896)
            .build()
            .unwrap();

        assert_eq!(fp.dimension(), DIM_1792);
        assert_eq!(fp.measurement_type(), MeasurementType::Ultimate);
    }

    #[test]
    fn test_normalization() {
        let data = vec![3.0f32, 4.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0]
            .into_iter()
            .chain(std::iter::repeat(0.0).take(80))
            .collect();

        let mut fp = FingerprintVector::basic(data).unwrap();
        fp.normalize().unwrap();

        // L2 norm should be 1.0
        let norm: f32 = fp.data().iter().map(|x| x * x).sum::<f32>().sqrt();
        assert!((norm - 1.0).abs() < 0.001);
    }
}
