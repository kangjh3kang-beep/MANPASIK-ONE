// api.rs (flutter_rust_bridge 연동 진입점)

use crate::ai::xgboost_inference::XGInferenceEngine;
use crate::fingerprint::feature_extractor::FeatureExtractor;
use crate::fingerprint::multi_mode_expander::Expander;
use crate::signal::differential::DifferentialEngine;

/// Flutter Riverpod에서 직접 호출할 FFI 브릿지 인터페이스
pub struct ManpasikFfiBridge {
    diff_engine: DifferentialEngine,
    feature_extractor: FeatureExtractor,
    ai_engine: XGInferenceEngine,
}

impl Default for ManpasikFfiBridge {
    fn default() -> Self {
        Self {
            diff_engine: DifferentialEngine::new(),
            feature_extractor: FeatureExtractor::new(256),
            ai_engine: XGInferenceEngine::new("MPS-XGB-INT8-v2.1"),
        }
    }
}

impl ManpasikFfiBridge {
    pub fn new() -> Self {
        Self::default()
    }

    /// BLE 스트리밍 수신 직후 Flutter에서 동기적으로 차동연산 수행
    pub fn apply_differential(&self, s_detection: f64, s_reference: f64) -> f64 {
        self.diff_engine.compute(s_detection, s_reference)
    }

    /// 핑거프린트 추출 및 인공지능 로컬 추론을 한 큐에 실행 (CRDT 기록용 결과 반환)
    pub fn run_inference(&self, differential_wave: Vec<f64>, _temperature: f64) -> String {
        // 1. 88차원 특성 추출
        let base_features = self.feature_extractor.extract(&differential_wave);

        // 2. 896차원으로 공간 확장 (eNose, eTongue 등 융합 모의)
        let dummy_enose = vec![0.0; 8];
        let expanded_features = Expander::expand_to_896(&base_features, &dummy_enose);

        // 3. TFLite 로컬 추론 동작 및 모델 평가
        match self.ai_engine.predict(&expanded_features) {
            Ok(prediction) => {
                format!(
                    "정량 예측치: {} mg/dL (Anomoly Score: {})",
                    prediction.quantification_mg_dl, prediction.anomaly_score
                )
            }
            Err(e) => format!("로컬 추론 에러: {}", e),
        }
    }
}
