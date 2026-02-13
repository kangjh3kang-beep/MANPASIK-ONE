//! AI 추론 모듈
//!
//! TFLite + ONNX 기반 엣지 AI 추론 엔진

use serde::{Deserialize, Serialize};
use std::path::Path;
use thiserror::Error;

/// AI 모델 타입
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
pub enum ModelType {
    /// 보정 모델 (Calibration)
    Calibration,
    /// 핑거프린트 분류 (Classification)
    FingerprintClassifier,
    /// 이상 탐지 (Anomaly Detection)
    AnomalyDetection,
    /// 값 예측 (Regression)
    ValuePredictor,
    /// 품질 평가 (Quality Assessment)
    QualityAssessment,
}

/// 모델 포맷
#[derive(Debug, Clone, Copy)]
pub enum ModelFormat {
    TFLite,
    Onnx,
}

#[derive(Debug, Error)]
pub enum InferenceError {
    #[error("모델을 찾을 수 없습니다: {0}")]
    ModelNotFound(String),

    #[error("모델 로드 실패: {0}")]
    LoadError(String),

    #[error("입력 형식 오류: expected {expected}, got {got}")]
    InputShapeMismatch { expected: usize, got: usize },

    #[error("추론 실패: {0}")]
    InferenceError(String),

    #[error("출력 형식 오류")]
    OutputError,
}

/// 추론 결과
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct InferenceResult {
    /// 출력 값
    pub values: Vec<f32>,
    /// 신뢰도 (0.0 ~ 1.0)
    pub confidence: f32,
    /// 추론 시간 (밀리초)
    pub inference_time_ms: f32,
    /// 모델 타입
    pub model_type: ModelType,
}

/// AI 추론 엔진
pub struct InferenceEngine {
    model_type: ModelType,
    input_size: usize,
    output_size: usize,
    model_loaded: bool,
}

impl InferenceEngine {
    /// 새 추론 엔진 생성
    pub fn new(model_type: ModelType) -> Self {
        let (input_size, output_size) = match model_type {
            ModelType::Calibration => (88, 88), // 채널별 보정값 출력
            ModelType::FingerprintClassifier => (1792, 30), // 1792차원 궁극 입력, 30종 분류 (29종 + NonTarget1792)
            ModelType::AnomalyDetection => (88, 1),         // 이상 스코어
            ModelType::ValuePredictor => (88, 1),           // 단일 예측값
            ModelType::QualityAssessment => (88, 3),        // 품질 등급 (좋음/보통/나쁨)
        };

        Self {
            model_type,
            input_size,
            output_size,
            model_loaded: false,
        }
    }

    /// 보정 모델 엔진 생성
    pub fn calibration() -> Self {
        Self::new(ModelType::Calibration)
    }

    /// 핑거프린트 분류 엔진 생성
    pub fn classifier() -> Self {
        Self::new(ModelType::FingerprintClassifier)
    }

    /// 이상 탐지 엔진 생성
    pub fn anomaly_detector() -> Self {
        Self::new(ModelType::AnomalyDetection)
    }

    /// 모델 로드
    pub fn load_model(&mut self, model_path: &Path) -> Result<(), InferenceError> {
        if !model_path.exists() {
            return Err(InferenceError::ModelNotFound(
                model_path.to_string_lossy().to_string(),
            ));
        }

        #[cfg(feature = "ai")]
        {
            // TFLite 모델 로드 (tflitec 사용)
            // TODO: 실제 모델 로드 구현
        }

        self.model_loaded = true;
        Ok(())
    }

    /// 추론 수행
    pub fn predict(&self, input: &[f32]) -> Result<InferenceResult, InferenceError> {
        // 입력 크기 검증
        if input.len() != self.input_size {
            return Err(InferenceError::InputShapeMismatch {
                expected: self.input_size,
                got: input.len(),
            });
        }

        let start = std::time::Instant::now();

        // 모델이 로드되지 않은 경우 시뮬레이션 추론
        let (values, confidence) = if self.model_loaded {
            #[cfg(feature = "ai")]
            {
                // 실제 TFLite 추론
                // TODO: tflitec 구현
                self.simulate_inference(input)
            }
            #[cfg(not(feature = "ai"))]
            {
                self.simulate_inference(input)
            }
        } else {
            self.simulate_inference(input)
        };

        let inference_time_ms = start.elapsed().as_secs_f32() * 1000.0;

        Ok(InferenceResult {
            values,
            confidence,
            inference_time_ms,
            model_type: self.model_type,
        })
    }

    /// 시뮬레이션 추론 (모델 없이 테스트용)
    fn simulate_inference(&self, input: &[f32]) -> (Vec<f32>, f32) {
        match self.model_type {
            ModelType::Calibration => {
                // 입력에 약간의 보정 적용
                let values: Vec<f32> = input.iter().map(|&v| v * 0.98 + 0.01).collect();
                (values, 0.95)
            }
            ModelType::FingerprintClassifier => {
                // 30개 클래스에 대한 확률 분포 (29종 레거시 + NonTarget1792)
                let num_classes = self.output_size;
                let mut values = vec![0.0f32; num_classes];
                // 간단한 해시 기반 클래스 선택 (시뮬레이션)
                let class_idx = (input.iter().sum::<f32>().abs() as usize) % num_classes;
                values[class_idx] = 0.85;
                (values, 0.85)
            }
            ModelType::AnomalyDetection => {
                // 입력 범위 기반 이상 스코어 계산
                let mean: f32 = input.iter().sum::<f32>() / input.len() as f32;
                let variance: f32 =
                    input.iter().map(|x| (x - mean).powi(2)).sum::<f32>() / input.len() as f32;

                // 분산이 높으면 이상
                let anomaly_score = (variance / 10.0).min(1.0);
                (vec![anomaly_score], 0.90)
            }
            ModelType::ValuePredictor => {
                // 입력 평균 기반 값 예측
                let mean = input.iter().sum::<f32>() / input.len() as f32;
                (vec![mean * 100.0], 0.92)
            }
            ModelType::QualityAssessment => {
                // 품질 등급 (좋음, 보통, 나쁨) 확률
                let variance: f32 = {
                    let mean = input.iter().sum::<f32>() / input.len() as f32;
                    input.iter().map(|x| (x - mean).powi(2)).sum::<f32>() / input.len() as f32
                };

                if variance < 0.1 {
                    (vec![0.9, 0.08, 0.02], 0.90) // 좋음
                } else if variance < 0.5 {
                    (vec![0.2, 0.7, 0.1], 0.70) // 보통
                } else {
                    (vec![0.1, 0.2, 0.7], 0.70) // 나쁨
                }
            }
        }
    }

    /// 배치 추론
    pub fn predict_batch(
        &self,
        inputs: &[Vec<f32>],
    ) -> Result<Vec<InferenceResult>, InferenceError> {
        inputs.iter().map(|input| self.predict(input)).collect()
    }

    /// 모델 타입 조회
    pub fn model_type(&self) -> ModelType {
        self.model_type
    }

    /// 입력 크기 조회
    pub fn input_size(&self) -> usize {
        self.input_size
    }

    /// 출력 크기 조회
    pub fn output_size(&self) -> usize {
        self.output_size
    }

    /// 모델 로드 상태
    pub fn is_loaded(&self) -> bool {
        self.model_loaded
    }
}

/// 모델 관리자 (여러 모델 관리)
pub struct ModelManager {
    calibration: Option<InferenceEngine>,
    classifier: Option<InferenceEngine>,
    anomaly_detector: Option<InferenceEngine>,
}

impl ModelManager {
    pub fn new() -> Self {
        Self {
            calibration: None,
            classifier: None,
            anomaly_detector: None,
        }
    }

    /// 기본 모델들 로드
    pub fn load_defaults(&mut self, models_dir: &Path) -> Result<(), InferenceError> {
        // 보정 모델
        let mut calibration = InferenceEngine::calibration();
        let cal_path = models_dir.join("calibration.tflite");
        if cal_path.exists() {
            calibration.load_model(&cal_path)?;
        }
        self.calibration = Some(calibration);

        // 분류 모델
        let mut classifier = InferenceEngine::classifier();
        let cls_path = models_dir.join("classifier.tflite");
        if cls_path.exists() {
            classifier.load_model(&cls_path)?;
        }
        self.classifier = Some(classifier);

        // 이상 탐지 모델
        let mut anomaly = InferenceEngine::anomaly_detector();
        let ano_path = models_dir.join("anomaly.tflite");
        if ano_path.exists() {
            anomaly.load_model(&ano_path)?;
        }
        self.anomaly_detector = Some(anomaly);

        Ok(())
    }

    /// 보정 추론
    pub fn calibrate(&self, input: &[f32]) -> Result<InferenceResult, InferenceError> {
        self.calibration
            .as_ref()
            .ok_or(InferenceError::ModelNotFound("calibration".to_string()))?
            .predict(input)
    }

    /// 분류 추론
    pub fn classify(&self, input: &[f32]) -> Result<InferenceResult, InferenceError> {
        self.classifier
            .as_ref()
            .ok_or(InferenceError::ModelNotFound("classifier".to_string()))?
            .predict(input)
    }

    /// 이상 탐지 추론
    pub fn detect_anomaly(&self, input: &[f32]) -> Result<InferenceResult, InferenceError> {
        self.anomaly_detector
            .as_ref()
            .ok_or(InferenceError::ModelNotFound(
                "anomaly_detector".to_string(),
            ))?
            .predict(input)
    }
}

impl Default for ModelManager {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_calibration_engine() {
        let engine = InferenceEngine::calibration();
        let input = vec![1.0f32; 88];

        let result = engine.predict(&input).unwrap();

        assert_eq!(result.values.len(), 88);
        assert!(result.confidence > 0.0);
    }

    #[test]
    fn test_classifier_engine() {
        let engine = InferenceEngine::classifier();
        let input = vec![0.5f32; 1792];

        let result = engine.predict(&input).unwrap();

        assert_eq!(result.values.len(), 30);
        // 최대 확률이 있어야 함
        let max_prob: f32 = result
            .values
            .iter()
            .cloned()
            .fold(f32::NEG_INFINITY, f32::max);
        assert!(max_prob > 0.0);
    }

    #[test]
    fn test_anomaly_detector() {
        let engine = InferenceEngine::anomaly_detector();

        // 정상 데이터
        let normal_input = vec![1.0f32; 88];
        let normal_result = engine.predict(&normal_input).unwrap();

        // 이상 데이터 (높은 분산)
        let anomaly_input: Vec<f32> = (0..88).map(|i| (i as f32) * 0.5).collect();
        let anomaly_result = engine.predict(&anomaly_input).unwrap();

        // 이상 데이터가 더 높은 이상 스코어를 가져야 함
        assert!(anomaly_result.values[0] > normal_result.values[0]);
    }

    #[test]
    fn test_input_shape_mismatch() {
        let engine = InferenceEngine::calibration();
        let wrong_input = vec![1.0f32; 50]; // 88이 아닌 50

        let result = engine.predict(&wrong_input);
        assert!(matches!(
            result,
            Err(InferenceError::InputShapeMismatch { .. })
        ));
    }
}
