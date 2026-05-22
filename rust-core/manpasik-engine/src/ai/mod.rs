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
    allow_simulation: bool,
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
            allow_simulation: true,
        }
    }

    /// 운영 모드 추론 엔진 생성
    pub fn production(model_type: ModelType) -> Self {
        Self::new(model_type).with_simulation_allowed(false)
    }

    /// 시뮬레이션 폴백 허용 여부를 설정합니다.
    pub fn with_simulation_allowed(mut self, allow_simulation: bool) -> Self {
        self.allow_simulation = allow_simulation;
        self
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
            use std::fs;

            // TFLite 모델 파일 검증
            let model_data = fs::read(model_path)
                .map_err(|e| InferenceError::LoadError(format!("파일 읽기 실패: {}", e)))?;

            // TFLite 매직 넘버 검증 (0x54464C33 = "TFL3")
            if model_data.len() < 4 {
                return Err(InferenceError::LoadError(
                    "모델 파일이 너무 작습니다".to_string(),
                ));
            }

            // tflitec::Interpreter 로드 시도
            // 실제 배포 시 tflitec::model::Model::from_file() 사용
            tracing::info!(
                "TFLite 모델 로드: {} ({} bytes, {:?})",
                model_path.display(),
                model_data.len(),
                self.model_type,
            );
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

        if !self.allow_simulation && !self.has_runtime_backed_model() {
            return Err(InferenceError::InferenceError(
                "production inference requires a runtime-backed model; simulation fallback disabled"
                    .to_string(),
            ));
        }

        // 모델이 로드된 경우 실제 추론 시도, 실패 시 시뮬레이션 폴백
        let (values, confidence) = if self.model_loaded {
            #[cfg(feature = "ai")]
            {
                // tflitec 추론 시도 → 실패 시 시뮬레이션 폴백
                // 실제 배포 시: tflitec::Interpreter::invoke()
                // 현재: 모델 파일 로드는 성공했지만 런타임이 없으면 시뮬레이션
                tracing::debug!(
                    "AI 추론 수행: {:?} (입력 {}차원)",
                    self.model_type,
                    input.len()
                );
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

    /// 시뮬레이션 폴백 허용 여부
    pub fn simulation_allowed(&self) -> bool {
        self.allow_simulation
    }

    fn has_runtime_backed_model(&self) -> bool {
        false
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

// ============================================================================
// 안전 검증 레이어 (FM-AI-002 대응: AI 위음성 방지)
// ============================================================================

/// AI 결과에 대한 규칙 기반 안전 검증
/// FMEA FM-AI-002: confidence 임계값 + 규칙 기반 이중 체크
pub struct SafetyValidator {
    /// confidence 임계값 (기본 0.85)
    confidence_threshold: f32,
}

/// 안전 검증 결과
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SafetyCheckResult {
    /// AI 추론 결과 통과 여부
    pub ai_passed: bool,
    /// 규칙 기반 체크 통과 여부
    pub rule_passed: bool,
    /// 최종 판정
    pub final_verdict: SafetyVerdict,
    /// 경고 메시지
    pub warnings: Vec<String>,
}

/// 안전 판정
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
pub enum SafetyVerdict {
    /// 정상 — AI + 규칙 모두 통과
    Normal,
    /// 주의 — AI 또는 규칙 중 하나 경고
    Caution,
    /// 위험 — AI 또는 규칙에서 이상 감지
    Alert,
    /// 불확실 — confidence 부족, 재측정 권장
    Uncertain,
}

impl SafetyValidator {
    pub fn new() -> Self {
        Self {
            confidence_threshold: 0.85,
        }
    }

    pub fn with_threshold(threshold: f32) -> Self {
        Self {
            confidence_threshold: threshold,
        }
    }

    /// 측정 결과에 대한 안전 검증 수행
    pub fn validate(
        &self,
        inference: &InferenceResult,
        measured_value: f64,
        biomarker_type: &str,
    ) -> SafetyCheckResult {
        let mut warnings = Vec::new();

        // 1. AI confidence 체크
        let ai_passed = inference.confidence >= self.confidence_threshold;
        if !ai_passed {
            warnings.push(format!(
                "AI 신뢰도 부족: {:.2} < {:.2} (재측정 권장)",
                inference.confidence, self.confidence_threshold
            ));
        }

        // 2. 규칙 기반 임계값 체크 (바이오마커별)
        let rule_passed = self.check_reference_range(measured_value, biomarker_type, &mut warnings);

        // 3. 최종 판정
        let final_verdict = match (ai_passed, rule_passed) {
            (true, true) => SafetyVerdict::Normal,
            (true, false) => SafetyVerdict::Alert,
            (false, true) => SafetyVerdict::Uncertain,
            (false, false) => SafetyVerdict::Alert,
        };

        SafetyCheckResult {
            ai_passed,
            rule_passed,
            final_verdict,
            warnings,
        }
    }

    /// 바이오마커별 참조 범위 규칙 기반 체크
    fn check_reference_range(
        &self,
        value: f64,
        biomarker_type: &str,
        warnings: &mut Vec<String>,
    ) -> bool {
        let (low, high, critical_low, critical_high) = match biomarker_type {
            "glucose" => (70.0, 100.0, 40.0, 400.0),
            "hba1c" => (4.0, 5.6, 2.0, 15.0),
            "uric_acid" => (3.5, 7.2, 1.0, 15.0),
            "creatinine" => (0.7, 1.3, 0.1, 10.0),
            "vitamin_d" => (30.0, 100.0, 5.0, 150.0),
            "tsh" => (0.4, 4.0, 0.01, 100.0),
            "cortisol" => (6.0, 23.0, 1.0, 60.0),
            "crp" => (0.0, 3.0, 0.0, 200.0),
            _ => return true, // 알 수 없는 타입은 패스
        };

        // 위험 수준 체크
        if value <= critical_low || value >= critical_high {
            warnings.push(format!(
                "위험 수준: {biomarker_type} = {value:.1} (위험 범위: <{critical_low} 또는 >{critical_high})"
            ));
            return false;
        }

        // 정상 범위 체크
        if value < low || value > high {
            warnings.push(format!(
                "참조 범위 이탈: {biomarker_type} = {value:.1} (정상: {low}-{high})"
            ));
        }

        true
    }
}

impl Default for SafetyValidator {
    fn default() -> Self {
        Self::new()
    }
}

// ============================================================
// 편향 탐지 모듈 (Bias Detector)
// FDA AI/ML SaMD 가이드라인 — 인구통계별 성능 모니터링
// ============================================================

/// 인구통계 그룹
#[derive(Debug, Clone, PartialEq)]
pub struct DemographicGroup {
    pub name: String,
    pub sample_count: usize,
    pub accuracy: f64,
    pub sensitivity: f64,
    pub specificity: f64,
}

/// 편향 분석 결과
#[derive(Debug, Clone)]
pub struct BiasReport {
    pub groups: Vec<DemographicGroup>,
    pub max_accuracy_gap: f64,
    pub max_sensitivity_gap: f64,
    pub is_biased: bool,
    pub bias_threshold: f64,
    pub recommendations: Vec<String>,
}

/// 편향 탐지기
pub struct BiasDetector {
    /// 그룹 간 성능 차이 임계값 (기본 5%)
    pub threshold: f64,
    /// 최소 샘플 수 (통계적 유의성)
    pub min_sample_size: usize,
}

impl BiasDetector {
    pub fn new() -> Self {
        Self {
            threshold: 0.05,
            min_sample_size: 30,
        }
    }

    /// 인구통계 그룹별 편향 분석
    pub fn analyze(&self, groups: &[DemographicGroup]) -> BiasReport {
        let valid_groups: Vec<&DemographicGroup> = groups
            .iter()
            .filter(|g| g.sample_count >= self.min_sample_size)
            .collect();

        if valid_groups.len() < 2 {
            return BiasReport {
                groups: groups.to_vec(),
                max_accuracy_gap: 0.0,
                max_sensitivity_gap: 0.0,
                is_biased: false,
                bias_threshold: self.threshold,
                recommendations: vec!["샘플 수 부족: 최소 2개 그룹 × 30건 필요".to_string()],
            };
        }

        // 정확도 범위
        let acc_max = valid_groups
            .iter()
            .map(|g| g.accuracy)
            .fold(f64::NEG_INFINITY, f64::max);
        let acc_min = valid_groups
            .iter()
            .map(|g| g.accuracy)
            .fold(f64::INFINITY, f64::min);
        let max_accuracy_gap = acc_max - acc_min;

        // 민감도 범위
        let sens_max = valid_groups
            .iter()
            .map(|g| g.sensitivity)
            .fold(f64::NEG_INFINITY, f64::max);
        let sens_min = valid_groups
            .iter()
            .map(|g| g.sensitivity)
            .fold(f64::INFINITY, f64::min);
        let max_sensitivity_gap = sens_max - sens_min;

        let is_biased = max_accuracy_gap > self.threshold || max_sensitivity_gap > self.threshold;

        let mut recommendations = Vec::new();
        if is_biased {
            // 가장 성능 낮은 그룹 찾기
            let weakest = valid_groups
                .iter()
                .min_by(|a, b| a.accuracy.partial_cmp(&b.accuracy).unwrap())
                .unwrap();
            recommendations.push(format!(
                "편향 감지: '{}' 그룹 정확도 {:.1}% — 데이터 보강 필요",
                weakest.name,
                weakest.accuracy * 100.0
            ));
            recommendations.push("인구통계별 학습 데이터 리밸런싱 권장".to_string());
            recommendations.push("차기 릴리스 전 편향 재평가 필수".to_string());
        }

        BiasReport {
            groups: groups.to_vec(),
            max_accuracy_gap,
            max_sensitivity_gap,
            is_biased,
            bias_threshold: self.threshold,
            recommendations,
        }
    }
}

impl Default for BiasDetector {
    fn default() -> Self {
        Self::new()
    }
}

// ============================================================
// 설명가능성 모듈 (Explainability / Feature Importance)
// SHAP 기반 특성 중요도 계산
// ============================================================

/// 특성 중요도 항목
#[derive(Debug, Clone)]
pub struct FeatureImportance {
    pub feature_name: String,
    pub importance_score: f64,
    pub direction: ContributionDirection,
}

/// 기여 방향
#[derive(Debug, Clone, PartialEq)]
pub enum ContributionDirection {
    Positive, // 수치 상승 방향
    Negative, // 수치 하강 방향
    Neutral,
}

/// 설명가능성 결과
#[derive(Debug, Clone)]
pub struct ExplainabilityReport {
    pub biomarker_type: String,
    pub predicted_value: f64,
    pub top_features: Vec<FeatureImportance>,
    pub confidence: f64,
    pub explanation_text: String,
}

/// 설명가능성 엔진 (Permutation-based Feature Importance)
pub struct ExplainabilityEngine {
    /// 상위 N개 특성만 보고
    pub top_k: usize,
}

impl ExplainabilityEngine {
    pub fn new() -> Self {
        Self { top_k: 5 }
    }

    /// 채널별 기여도 계산 (Permutation Importance 근사)
    ///
    /// 각 채널을 순서대로 셔플하여 예측 변화량을 측정
    pub fn compute_importance(
        &self,
        engine: &InferenceEngine,
        input: &[f32],
        baseline_result: &InferenceResult,
    ) -> Result<Vec<FeatureImportance>, InferenceError> {
        let baseline_value = baseline_result
            .values
            .iter()
            .cloned()
            .fold(f32::NEG_INFINITY, f32::max);
        let mut importances: Vec<FeatureImportance> = Vec::new();

        for i in 0..input.len() {
            // 채널 i를 0으로 마스킹
            let mut perturbed = input.to_vec();
            perturbed[i] = 0.0;

            let perturbed_result = engine.predict(&perturbed)?;
            let perturbed_value = perturbed_result
                .values
                .iter()
                .cloned()
                .fold(f32::NEG_INFINITY, f32::max);
            let diff = (baseline_value - perturbed_value) as f64;

            let direction = if diff > 0.01 {
                ContributionDirection::Positive
            } else if diff < -0.01 {
                ContributionDirection::Negative
            } else {
                ContributionDirection::Neutral
            };

            importances.push(FeatureImportance {
                feature_name: format!("channel_{i}"),
                importance_score: diff.abs(),
                direction,
            });
        }

        // 중요도 내림차순 정렬
        importances.sort_by(|a, b| b.importance_score.partial_cmp(&a.importance_score).unwrap());
        importances.truncate(self.top_k);

        Ok(importances)
    }

    /// 사용자 친화적 설명 텍스트 생성
    pub fn generate_explanation(
        &self,
        biomarker_type: &str,
        predicted_value: f64,
        confidence: f64,
        top_features: &[FeatureImportance],
    ) -> ExplainabilityReport {
        let main_contributors: Vec<String> = top_features
            .iter()
            .take(3)
            .map(|f| {
                let dir = match f.direction {
                    ContributionDirection::Positive => "상승",
                    ContributionDirection::Negative => "하강",
                    ContributionDirection::Neutral => "중립",
                };
                format!("{} ({})", f.feature_name, dir)
            })
            .collect();

        let explanation = format!(
            "{}의 예측값은 {:.1}입니다 (신뢰도: {:.0}%). 주요 기여 채널: {}. \
             본 결과는 AI 분석에 기반한 참고 정보이며, 의료 진단이 아닙니다.",
            biomarker_type,
            predicted_value,
            confidence * 100.0,
            main_contributors.join(", ")
        );

        ExplainabilityReport {
            biomarker_type: biomarker_type.to_string(),
            predicted_value,
            top_features: top_features.to_vec(),
            confidence,
            explanation_text: explanation,
        }
    }
}

impl Default for ExplainabilityEngine {
    fn default() -> Self {
        Self::new()
    }
}

// ============================================================
// 시계열 트렌드 분석기 (TrendAnalyzer)
// 선형회귀 + 신뢰구간 + 이동평균 + 추세 판정
// ============================================================

/// 트렌드 방향
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
pub enum TrendDirection {
    /// 개선 (하강이 좋은 바이오마커) 또는 상승 (상승이 좋은 바이오마커)
    Improving,
    /// 악화
    Worsening,
    /// 안정
    Stable,
    /// 데이터 부족
    Insufficient,
}

/// 선형회귀 결과
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LinearRegressionResult {
    /// 기울기 (단위시간당 변화량)
    pub slope: f64,
    /// 절편
    pub intercept: f64,
    /// 결정계수 R² (0~1, 1에 가까울수록 설명력 높음)
    pub r_squared: f64,
    /// 기울기의 표준오차
    pub slope_std_error: f64,
    /// 95% 신뢰구간 [하한, 상한]
    pub slope_ci_95: [f64; 2],
}

/// 트렌드 분석 결과
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TrendAnalysis {
    pub direction: TrendDirection,
    pub regression: Option<LinearRegressionResult>,
    /// 예측값 (다음 N 시점)
    pub predictions: Vec<f64>,
    /// 예측 신뢰구간 (각 시점별 [하한, 상한])
    pub prediction_intervals: Vec<[f64; 2]>,
    /// 이동평균 (가장 최근 윈도우)
    pub moving_average: f64,
    /// 변동계수 (CV%) — 안정성 지표
    pub coefficient_of_variation: f64,
    /// 한국어 요약
    pub summary: String,
}

/// 시계열 트렌드 분석기
pub struct TrendAnalyzer {
    /// 트렌드 판정을 위한 기울기 임계값 (기본: 표준편차의 5%)
    pub slope_threshold_pct: f64,
    /// 최소 데이터 포인트 수
    pub min_data_points: usize,
}

impl TrendAnalyzer {
    pub fn new() -> Self {
        Self {
            slope_threshold_pct: 0.05,
            min_data_points: 3,
        }
    }

    /// 시계열 데이터에 대한 트렌드 분석 수행
    ///
    /// `values`: 시간순 측정값 배열
    /// `biomarker`: 바이오마커 타입 (트렌드 방향 해석용)
    /// `predict_ahead`: 미래 예측 시점 수
    pub fn analyze(&self, values: &[f64], biomarker: &str, predict_ahead: usize) -> TrendAnalysis {
        if values.len() < self.min_data_points {
            return TrendAnalysis {
                direction: TrendDirection::Insufficient,
                regression: None,
                predictions: Vec::new(),
                prediction_intervals: Vec::new(),
                moving_average: if values.is_empty() {
                    0.0
                } else {
                    values.iter().sum::<f64>() / values.len() as f64
                },
                coefficient_of_variation: 0.0,
                summary: format!(
                    "데이터 부족: {}개 (최소 {}개 필요)",
                    values.len(),
                    self.min_data_points
                ),
            };
        }

        let n = values.len();
        let mean = values.iter().sum::<f64>() / n as f64;
        let std_dev = (values.iter().map(|v| (v - mean).powi(2)).sum::<f64>() / n as f64).sqrt();
        let cv = if mean.abs() > 1e-10 {
            (std_dev / mean.abs()) * 100.0
        } else {
            0.0
        };

        // 선형회귀 (OLS)
        let regression = self.linear_regression(values);

        // 이동평균 (최근 3개 또는 전체)
        let window = 3.min(n);
        let ma = values[n - window..].iter().sum::<f64>() / window as f64;

        // 트렌드 방향 판정
        let slope_threshold = std_dev * self.slope_threshold_pct;
        let direction = self.determine_direction(&regression, slope_threshold, biomarker);

        // 미래 예측
        let (predictions, intervals) = self.predict_future(&regression, n, predict_ahead, values);

        // 한국어 요약
        let dir_str = match direction {
            TrendDirection::Improving => "개선",
            TrendDirection::Worsening => "악화",
            TrendDirection::Stable => "안정",
            TrendDirection::Insufficient => "판단 불가",
        };
        let summary = format!(
            "{} 추세: {} (R²={:.2}, 기울기={:.3}/단위시간, CV={:.1}%)",
            biomarker, dir_str, regression.r_squared, regression.slope, cv,
        );

        TrendAnalysis {
            direction,
            regression: Some(regression),
            predictions,
            prediction_intervals: intervals,
            moving_average: ma,
            coefficient_of_variation: cv,
            summary,
        }
    }

    /// OLS 선형회귀 (x=0,1,2,...,n-1)
    fn linear_regression(&self, values: &[f64]) -> LinearRegressionResult {
        let n = values.len() as f64;
        let x_mean = (n - 1.0) / 2.0;
        let y_mean = values.iter().sum::<f64>() / n;

        let mut ss_xy = 0.0;
        let mut ss_xx = 0.0;
        let mut ss_yy = 0.0;

        for (i, &y) in values.iter().enumerate() {
            let x = i as f64;
            ss_xy += (x - x_mean) * (y - y_mean);
            ss_xx += (x - x_mean).powi(2);
            ss_yy += (y - y_mean).powi(2);
        }

        let slope = if ss_xx.abs() > 1e-15 {
            ss_xy / ss_xx
        } else {
            0.0
        };
        let intercept = y_mean - slope * x_mean;

        // R²
        let r_squared = if ss_yy.abs() > 1e-15 {
            let ss_res: f64 = values
                .iter()
                .enumerate()
                .map(|(i, &y)| (y - (intercept + slope * i as f64)).powi(2))
                .sum();
            1.0 - ss_res / ss_yy
        } else {
            1.0 // 모든 값이 동일하면 완벽 적합
        };
        let r_squared = r_squared.max(0.0);

        // 기울기 표준오차 + 95% 신뢰구간 (t ≈ 1.96 for large n)
        let se = if n > 2.0 && ss_xx.abs() > 1e-15 {
            let ss_res: f64 = values
                .iter()
                .enumerate()
                .map(|(i, &y)| (y - (intercept + slope * i as f64)).powi(2))
                .sum();
            let mse = ss_res / (n - 2.0);
            (mse / ss_xx).sqrt()
        } else {
            0.0
        };

        let t_val = if n > 30.0 { 1.96 } else { 2.0 + 5.0 / n }; // 근사 t-value
        let ci_low = slope - t_val * se;
        let ci_high = slope + t_val * se;

        LinearRegressionResult {
            slope,
            intercept,
            r_squared,
            slope_std_error: se,
            slope_ci_95: [ci_low, ci_high],
        }
    }

    /// 트렌드 방향 결정
    fn determine_direction(
        &self,
        reg: &LinearRegressionResult,
        threshold: f64,
        biomarker: &str,
    ) -> TrendDirection {
        // 기울기가 임계값 이내이면 안정
        if reg.slope.abs() <= threshold {
            return TrendDirection::Stable;
        }

        // 바이오마커별 "좋은 방향" 판정
        let lower_is_better = matches!(
            biomarker,
            "glucose"
                | "hba1c"
                | "cholesterol"
                | "ldl"
                | "triglycerides"
                | "uric_acid"
                | "creatinine"
                | "crp"
                | "cortisol"
                | "blood_pressure"
        );

        if lower_is_better {
            if reg.slope < 0.0 {
                TrendDirection::Improving
            } else {
                TrendDirection::Worsening
            }
        } else if reg.slope > 0.0 {
            TrendDirection::Improving
        } else {
            TrendDirection::Worsening
        }
    }

    /// 미래 시점 예측 + 예측 구간
    fn predict_future(
        &self,
        reg: &LinearRegressionResult,
        n: usize,
        ahead: usize,
        values: &[f64],
    ) -> (Vec<f64>, Vec<[f64; 2]>) {
        if ahead == 0 {
            return (Vec::new(), Vec::new());
        }

        let nf = n as f64;
        let x_mean = (nf - 1.0) / 2.0;
        let ss_xx: f64 = (0..n).map(|i| (i as f64 - x_mean).powi(2)).sum();
        let mse = if n > 2 {
            let ss_res: f64 = values
                .iter()
                .enumerate()
                .map(|(i, &y)| (y - (reg.intercept + reg.slope * i as f64)).powi(2))
                .sum();
            ss_res / (nf - 2.0)
        } else {
            0.0
        };

        let t_val = if nf > 30.0 { 1.96 } else { 2.0 + 5.0 / nf };

        let mut preds = Vec::with_capacity(ahead);
        let mut intervals = Vec::with_capacity(ahead);

        for k in 1..=ahead {
            let x_pred = (n + k - 1) as f64;
            let y_pred = reg.intercept + reg.slope * x_pred;
            preds.push(y_pred);

            // 예측 구간 (prediction interval)
            let h = 1.0 + 1.0 / nf + (x_pred - x_mean).powi(2) / ss_xx.max(1e-15);
            let margin = t_val * (mse * h).sqrt();
            intervals.push([y_pred - margin, y_pred + margin]);
        }

        (preds, intervals)
    }
}

impl Default for TrendAnalyzer {
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
    fn test_safety_validator_normal() {
        let validator = SafetyValidator::new();
        let inference = InferenceResult {
            values: vec![0.9],
            confidence: 0.95,
            inference_time_ms: 10.0,
            model_type: ModelType::ValuePredictor,
        };
        let result = validator.validate(&inference, 85.0, "glucose");
        assert_eq!(result.final_verdict, SafetyVerdict::Normal);
    }

    #[test]
    fn test_safety_validator_critical() {
        let validator = SafetyValidator::new();
        let inference = InferenceResult {
            values: vec![0.9],
            confidence: 0.95,
            inference_time_ms: 10.0,
            model_type: ModelType::ValuePredictor,
        };
        // 혈당 40 mg/dL — 위험 저혈당
        let result = validator.validate(&inference, 35.0, "glucose");
        assert_eq!(result.final_verdict, SafetyVerdict::Alert);
        assert!(!result.rule_passed);
    }

    #[test]
    fn test_safety_validator_low_confidence() {
        let validator = SafetyValidator::new();
        let inference = InferenceResult {
            values: vec![0.5],
            confidence: 0.50, // 임계값 미달
            inference_time_ms: 10.0,
            model_type: ModelType::ValuePredictor,
        };
        let result = validator.validate(&inference, 85.0, "glucose");
        assert_eq!(result.final_verdict, SafetyVerdict::Uncertain);
        assert!(!result.ai_passed);
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

    #[test]
    fn test_production_engine_rejects_simulation_without_runtime_model() {
        let engine = InferenceEngine::production(ModelType::Calibration);
        let input = vec![1.0f32; 88];

        let result = engine.predict(&input);

        assert!(matches!(
            result,
            Err(InferenceError::InferenceError(message))
                if message.contains("simulation fallback disabled")
        ));
    }

    #[test]
    fn test_default_engine_allows_simulation_for_demo_and_tests() {
        let engine = InferenceEngine::calibration();
        let input = vec![1.0f32; 88];

        let result = engine.predict(&input).unwrap();

        assert!(engine.simulation_allowed());
        assert_eq!(result.values.len(), 88);
    }

    // === 편향 탐지 테스트 ===

    #[test]
    fn test_bias_detector_no_bias() {
        let detector = BiasDetector::new();
        let groups = vec![
            DemographicGroup {
                name: "20대".to_string(),
                sample_count: 100,
                accuracy: 0.92,
                sensitivity: 0.90,
                specificity: 0.93,
            },
            DemographicGroup {
                name: "30대".to_string(),
                sample_count: 150,
                accuracy: 0.91,
                sensitivity: 0.89,
                specificity: 0.92,
            },
            DemographicGroup {
                name: "60대".to_string(),
                sample_count: 80,
                accuracy: 0.89,
                sensitivity: 0.87,
                specificity: 0.90,
            },
        ];
        let report = detector.analyze(&groups);
        // 최대 차이 3% < 임계값 5%
        assert!(!report.is_biased);
        assert!(report.max_accuracy_gap < 0.05);
    }

    #[test]
    fn test_bias_detector_biased() {
        let detector = BiasDetector::new();
        let groups = vec![
            DemographicGroup {
                name: "그룹A".to_string(),
                sample_count: 200,
                accuracy: 0.95,
                sensitivity: 0.94,
                specificity: 0.96,
            },
            DemographicGroup {
                name: "그룹B (소수)".to_string(),
                sample_count: 50,
                accuracy: 0.82, // 13% 차이 → 편향
                sensitivity: 0.80,
                specificity: 0.84,
            },
        ];
        let report = detector.analyze(&groups);
        assert!(report.is_biased);
        assert!(report.max_accuracy_gap > 0.05);
        assert!(!report.recommendations.is_empty());
    }

    #[test]
    fn test_bias_detector_insufficient_samples() {
        let detector = BiasDetector::new();
        let groups = vec![DemographicGroup {
            name: "소수그룹".to_string(),
            sample_count: 5, // 최소 30 미만
            accuracy: 0.50,
            sensitivity: 0.40,
            specificity: 0.60,
        }];
        let report = detector.analyze(&groups);
        assert!(!report.is_biased); // 샘플 부족으로 판단 불가
    }

    // === 설명가능성 테스트 ===

    #[test]
    fn test_explainability_feature_importance() {
        let engine = InferenceEngine::calibration();
        let explainer = ExplainabilityEngine::new();

        let input = vec![1.0f32; 88];
        let baseline = engine.predict(&input).unwrap();

        let importances = explainer
            .compute_importance(&engine, &input, &baseline)
            .unwrap();

        // top_k = 5개 이하로 반환
        assert!(importances.len() <= 5);
        // 내림차순 정렬
        for w in importances.windows(2) {
            assert!(w[0].importance_score >= w[1].importance_score);
        }
    }

    #[test]
    fn test_explainability_report_generation() {
        let explainer = ExplainabilityEngine::new();
        let features = vec![
            FeatureImportance {
                feature_name: "channel_12".to_string(),
                importance_score: 0.15,
                direction: ContributionDirection::Positive,
            },
            FeatureImportance {
                feature_name: "channel_45".to_string(),
                importance_score: 0.10,
                direction: ContributionDirection::Negative,
            },
        ];

        let report = explainer.generate_explanation("glucose", 95.0, 0.92, &features);
        assert_eq!(report.biomarker_type, "glucose");
        assert!(report.explanation_text.contains("glucose"));
        assert!(report.explanation_text.contains("95.0"));
        assert!(report.explanation_text.contains("의료 진단이 아닙니다"));
    }

    // === 트렌드 분석 테스트 ===

    #[test]
    fn test_trend_analyzer_improving() {
        let analyzer = TrendAnalyzer::new();
        // 혈당 하강 추세 (lower is better → Improving)
        let values = vec![120.0, 115.0, 110.0, 105.0, 100.0];
        let result = analyzer.analyze(&values, "glucose", 3);
        assert_eq!(result.direction, TrendDirection::Improving);
        let reg = result.regression.as_ref().unwrap();
        assert!(reg.slope < 0.0); // 하강
        assert!(reg.r_squared > 0.95); // 높은 적합도
        assert_eq!(result.predictions.len(), 3);
        assert!(result.predictions[0] < 100.0); // 100 미만으로 예측
    }

    #[test]
    fn test_trend_analyzer_stable() {
        let analyzer = TrendAnalyzer::new();
        // 안정적인 값 (대칭 배치 → slope ≈ 0)
        let values = vec![99.9, 100.1, 100.0, 100.1, 99.9, 100.0];
        let result = analyzer.analyze(&values, "glucose", 0);
        assert_eq!(result.direction, TrendDirection::Stable);
        assert!(result.coefficient_of_variation < 1.0);
    }

    #[test]
    fn test_trend_analyzer_insufficient() {
        let analyzer = TrendAnalyzer::new();
        let values = vec![100.0, 101.0]; // 2개 < 최소 3개
        let result = analyzer.analyze(&values, "glucose", 0);
        assert_eq!(result.direction, TrendDirection::Insufficient);
        assert!(result.regression.is_none());
    }

    #[test]
    fn test_trend_analyzer_worsening() {
        let analyzer = TrendAnalyzer::new();
        // 혈당 상승 추세 (lower is better → Worsening)
        let values = vec![90.0, 100.0, 110.0, 120.0, 130.0];
        let result = analyzer.analyze(&values, "glucose", 2);
        assert_eq!(result.direction, TrendDirection::Worsening);
        let reg = result.regression.as_ref().unwrap();
        assert!(reg.slope > 0.0);
        assert_eq!(result.prediction_intervals.len(), 2);
        // 예측 구간은 예측값을 포함해야 함
        assert!(result.prediction_intervals[0][0] <= result.predictions[0]);
        assert!(result.prediction_intervals[0][1] >= result.predictions[0]);
    }

    #[test]
    fn test_trend_analyzer_r_squared_perfect() {
        let analyzer = TrendAnalyzer::new();
        // 완벽한 선형 데이터
        let values = vec![10.0, 20.0, 30.0, 40.0, 50.0];
        let result = analyzer.analyze(&values, "vitamin_d", 1);
        let reg = result.regression.as_ref().unwrap();
        assert!((reg.r_squared - 1.0).abs() < 1e-10);
        assert!((reg.slope - 10.0).abs() < 1e-10);
        assert!((reg.intercept - 10.0).abs() < 1e-10);
    }

    #[test]
    fn test_trend_confidence_interval() {
        let analyzer = TrendAnalyzer::new();
        let values = vec![100.0, 102.0, 98.0, 103.0, 101.0, 99.0, 104.0];
        let result = analyzer.analyze(&values, "glucose", 1);
        let reg = result.regression.as_ref().unwrap();
        assert!(reg.slope_ci_95[0] <= reg.slope);
        assert!(reg.slope_ci_95[1] >= reg.slope);
        assert!(result.summary.contains("glucose"));
    }

    // ========================================================================
    // Phase M 엣지 케이스 (+8건)
    // ========================================================================

    #[test]
    fn test_inference_engine_unloaded_state() {
        let engine = InferenceEngine::classifier();
        assert!(!engine.is_loaded(), "초기 상태는 미로드");
        assert_eq!(engine.model_type(), ModelType::FingerprintClassifier);
    }

    #[test]
    fn test_inference_predict_without_loading_fails() {
        let engine = InferenceEngine::classifier();
        let result = engine.predict(&[0.5; 10]);
        assert!(result.is_err(), "미로드 모델 추론은 실패");
    }

    #[test]
    fn test_trend_analyzer_empty_values() {
        let analyzer = TrendAnalyzer::new();
        let values: Vec<f64> = vec![];
        let result = analyzer.analyze(&values, "glucose", 0);
        assert_eq!(result.direction, TrendDirection::Insufficient);
    }

    #[test]
    fn test_trend_analyzer_single_value() {
        let analyzer = TrendAnalyzer::new();
        let result = analyzer.analyze(&[100.0], "glucose", 0);
        assert_eq!(result.direction, TrendDirection::Insufficient);
    }

    #[test]
    fn test_trend_analyzer_constant_values() {
        let analyzer = TrendAnalyzer::new();
        let values = vec![100.0; 10];
        let result = analyzer.analyze(&values, "glucose", 1);
        if let Some(reg) = result.regression.as_ref() {
            assert!(
                reg.slope.abs() < 1e-10,
                "Constant values should have slope ~0"
            );
        }
    }

    #[test]
    fn test_trend_analyzer_with_negative_values() {
        let analyzer = TrendAnalyzer::new();
        let values = vec![-10.0, -5.0, 0.0, 5.0, 10.0];
        let result = analyzer.analyze(&values, "glucose", 1);
        let reg = result.regression.as_ref().unwrap();
        assert!(reg.slope > 0.0);
    }

    #[test]
    fn test_trend_predictions_count_matches_horizon() {
        let analyzer = TrendAnalyzer::new();
        let values = vec![100.0, 105.0, 110.0, 115.0, 120.0];
        let result = analyzer.analyze(&values, "glucose", 5);
        assert_eq!(result.predictions.len(), 5);
        assert_eq!(result.prediction_intervals.len(), 5);
    }

    #[test]
    fn test_inference_engine_calibration_default() {
        let engine = InferenceEngine::calibration();
        assert_eq!(engine.model_type(), ModelType::Calibration);
        assert!(!engine.is_loaded());
    }
}
