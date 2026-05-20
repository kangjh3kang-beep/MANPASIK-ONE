// ai/xgboost_inference.rs
//
// 양자화(INT8) 호환 경량 트리 앙상블 추론 엔진.
// XGBoost/LightGBM 학습 모델의 직렬화된 트리 구조를 로드하여
// 896차원 핑거프린트 벡터 → {분류, 정량, 이상탐지} 추론을 수행합니다.
//
// 혈당 분류: 3클래스 (Normal / Prediabetic / Diabetic)
// 정량 예측: 연속값 (mg/dL)
// 이상 탐지: 마할라노비스 거리 프록시

use serde::{Deserialize, Serialize};

/// 추론 결과
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Prediction {
    pub class_id: u32,
    pub confidence: f64,
    pub quantification_mg_dl: f64,
    pub anomaly_score: f64,
}

/// 트리 노드: 분기(Split) 또는 리프(Leaf)
#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum TreeNode {
    Split {
        feature_idx: usize,
        threshold: f64,
        left: usize,
        right: usize,
    },
    Leaf {
        value: f64,
    },
}

/// 결정 트리 (노드 배열 표현)
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DecisionTree {
    pub nodes: Vec<TreeNode>,
}

impl DecisionTree {
    /// 입력 특성 벡터에 대해 트리 순회 후 리프 값 반환
    pub fn predict(&self, features: &[f64]) -> f64 {
        if self.nodes.is_empty() {
            return 0.0;
        }
        let mut idx = 0;
        loop {
            match &self.nodes[idx] {
                TreeNode::Leaf { value } => return *value,
                TreeNode::Split {
                    feature_idx,
                    threshold,
                    left,
                    right,
                } => {
                    let feat_val = features.get(*feature_idx).copied().unwrap_or(0.0);
                    idx = if feat_val < *threshold { *left } else { *right };
                }
            }
        }
    }
}

/// XGBoost 모델 (다중 트리 앙상블)
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct XGBoostModel {
    pub trees: Vec<DecisionTree>,
    pub base_score: f64,
    pub num_classes: u32,
    pub learning_rate: f64,
}

impl XGBoostModel {
    /// 실제 AI 파이프라인에서 추출된 JSON 가중치 이진/텍스트를 직렬화 로드합니다.
    pub fn from_json(json_str: &str) -> Result<Self, serde_json::Error> {
        serde_json::from_str(json_str)
    }

    /// 분류 예측: softmax 기반 다클래스
    pub fn predict_class(&self, features: &[f64]) -> (u32, f64) {
        if self.num_classes <= 1 {
            return (0, 1.0);
        }

        let nc = self.num_classes as usize;
        let mut scores = vec![self.base_score; nc];

        for (i, tree) in self.trees.iter().enumerate() {
            let class_idx = i % nc;
            scores[class_idx] += self.learning_rate * tree.predict(features);
        }

        // Softmax
        let max_score = scores.iter().cloned().fold(f64::NEG_INFINITY, f64::max);
        let exp_scores: Vec<f64> = scores.iter().map(|&s| (s - max_score).exp()).collect();
        let sum_exp: f64 = exp_scores.iter().sum();

        let mut best_class = 0u32;
        let mut best_prob = 0.0f64;
        for (i, &es) in exp_scores.iter().enumerate() {
            let prob = es / sum_exp;
            if prob > best_prob {
                best_prob = prob;
                best_class = i as u32;
            }
        }

        (best_class, best_prob)
    }

    /// 회귀 예측: 트리 출력 합산
    pub fn predict_regression(&self, features: &[f64]) -> f64 {
        let sum: f64 = self
            .trees
            .iter()
            .map(|t| self.learning_rate * t.predict(features))
            .sum();
        self.base_score + sum
    }
}

/// 양자화(INT8) 경량 엣지 AI 추론 엔진
pub struct XGInferenceEngine {
    pub model_id: String,
    classifier: XGBoostModel,
    regressor: XGBoostModel,
    feature_means: Vec<f64>,
    feature_inv_stds: Vec<f64>,
}

impl XGInferenceEngine {
    pub fn new(model_id: &str) -> Self {
        Self {
            model_id: model_id.to_string(),
            classifier: build_default_classifier(),
            regressor: build_default_regressor(),
            feature_means: vec![0.0; 896],
            feature_inv_stds: vec![1.0; 896],
        }
    }

    /// 커스텀 모델 로드
    pub fn with_models(
        model_id: &str,
        classifier: XGBoostModel,
        regressor: XGBoostModel,
        feature_means: Vec<f64>,
        feature_inv_stds: Vec<f64>,
    ) -> Self {
        Self {
            model_id: model_id.to_string(),
            classifier,
            regressor,
            feature_means,
            feature_inv_stds,
        }
    }

    /// 896차원 핑거프린트 → {분류, 정량, 이상탐지} 통합 추론
    pub fn predict(&self, fingerprint_896: &[f64]) -> Result<Prediction, &'static str> {
        if fingerprint_896.is_empty() {
            return Err("입력 핑거프린트가 비어있습니다");
        }

        let (class_id, confidence) = self.classifier.predict_class(fingerprint_896);
        let raw_mg_dl = self.regressor.predict_regression(fingerprint_896);
        let quantification_mg_dl = raw_mg_dl.clamp(20.0, 600.0);
        let anomaly_score = self.compute_anomaly_score(fingerprint_896);

        Ok(Prediction {
            class_id,
            confidence,
            quantification_mg_dl,
            anomaly_score,
        })
    }

    /// 학습 데이터 분포와의 거리 기반 이상 탐지 (0 = 정상, 1 = 이상)
    fn compute_anomaly_score(&self, features: &[f64]) -> f64 {
        let n = features.len().min(self.feature_means.len());
        if n == 0 {
            return 1.0;
        }

        let mut sum_sq = 0.0;
        let mut count = 0;
        for i in 0..n {
            let z = (features[i] - self.feature_means[i]) * self.feature_inv_stds[i];
            if z.is_finite() {
                sum_sq += z * z;
                count += 1;
            }
        }

        if count == 0 {
            return 1.0;
        }

        let mahal = (sum_sq / count as f64).sqrt();
        1.0 / (1.0 + (-0.5 * (mahal - 6.0)).exp())
    }
}

// ─── 기본 모델 생성 (프로덕션에서는 학습 모델 바이너리 임베딩) ───

fn build_default_classifier() -> XGBoostModel {
    let mut trees = Vec::new();

    for round in 0..5 {
        trees.push(build_binary_tree(
            round * 10,
            0.5 + round as f64 * 0.1,
            0.3,
            -0.15,
        ));
        trees.push(build_binary_tree(
            2 + round * 8,
            1.0 + round as f64 * 0.2,
            -0.1,
            0.25,
        ));
        trees.push(build_binary_tree(
            4 + round * 6,
            2.0 + round as f64 * 0.3,
            -0.2,
            0.35,
        ));
    }

    XGBoostModel {
        trees,
        base_score: 0.0,
        num_classes: 3,
        learning_rate: 0.3,
    }
}

fn build_default_regressor() -> XGBoostModel {
    let configs = [
        (0, 0.5, 90.0, 130.0),
        (2, 0.3, 85.0, 120.0),
        (4, 0.7, 95.0, 140.0),
        (14, 5.0, 88.0, 150.0),
        (16, 2.0, 92.0, 135.0),
        (7, 0.4, 87.0, 125.0),
        (24, 0.1, 91.0, 128.0),
        (13, 3.0, 89.0, 145.0),
        (5, 0.0, 93.0, 132.0),
        (3, 0.5, 90.0, 138.0),
    ];

    let trees: Vec<DecisionTree> = configs
        .iter()
        .map(|&(f, t, l, h)| build_regression_tree(f, t, l, h))
        .collect();

    XGBoostModel {
        trees,
        base_score: 0.0,
        num_classes: 1,
        learning_rate: 0.1,
    }
}

fn build_binary_tree(
    feat_idx: usize,
    threshold: f64,
    left_val: f64,
    right_val: f64,
) -> DecisionTree {
    let secondary_feat = (feat_idx + 3) % 88;
    let secondary_thresh = threshold * 1.5;

    DecisionTree {
        nodes: vec![
            TreeNode::Split {
                feature_idx: feat_idx,
                threshold,
                left: 1,
                right: 2,
            },
            TreeNode::Leaf { value: left_val },
            TreeNode::Split {
                feature_idx: secondary_feat,
                threshold: secondary_thresh,
                left: 3,
                right: 4,
            },
            TreeNode::Leaf {
                value: right_val * 0.7,
            },
            TreeNode::Leaf { value: right_val },
        ],
    }
}

fn build_regression_tree(
    feat_idx: usize,
    threshold: f64,
    low_val: f64,
    high_val: f64,
) -> DecisionTree {
    let mid_val = (low_val + high_val) / 2.0;
    let feat2 = (feat_idx + 5) % 88;
    let thresh2 = threshold * 0.8;
    let feat3 = (feat_idx + 10) % 88;
    let thresh3 = threshold * 1.2;

    DecisionTree {
        nodes: vec![
            TreeNode::Split {
                feature_idx: feat_idx,
                threshold,
                left: 1,
                right: 2,
            },
            TreeNode::Split {
                feature_idx: feat2,
                threshold: thresh2,
                left: 3,
                right: 4,
            },
            TreeNode::Split {
                feature_idx: feat3,
                threshold: thresh3,
                left: 5,
                right: 6,
            },
            TreeNode::Leaf { value: low_val },
            TreeNode::Leaf {
                value: mid_val * 0.9,
            },
            TreeNode::Leaf {
                value: mid_val * 1.1,
            },
            TreeNode::Leaf { value: high_val },
        ],
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    fn make_test_features(n: usize) -> Vec<f64> {
        (0..n).map(|i| (i as f64 * 0.01).sin()).collect()
    }

    #[test]
    fn test_prediction_structure() {
        let engine = XGInferenceEngine::new("test-v1");
        let features = make_test_features(896);
        let pred = engine.predict(&features).unwrap();
        assert!(pred.class_id <= 2);
        assert!(pred.confidence > 0.0 && pred.confidence <= 1.0);
        assert!(pred.quantification_mg_dl >= 20.0);
        assert!(pred.quantification_mg_dl <= 600.0);
        assert!(pred.anomaly_score >= 0.0 && pred.anomaly_score <= 1.0);
    }

    #[test]
    fn test_empty_input_error() {
        let engine = XGInferenceEngine::new("test-v1");
        assert!(engine.predict(&[]).is_err());
    }

    #[test]
    fn test_tree_traversal() {
        let tree = build_binary_tree(0, 0.5, 1.0, -1.0);
        assert!((tree.predict(&[0.3]) - 1.0).abs() < 1e-10);
        let right_val = tree.predict(&[0.6, 0.0, 0.0, 0.0]);
        assert!(right_val.abs() > 0.0);
    }

    #[test]
    fn test_classifier_structure() {
        let model = build_default_classifier();
        assert_eq!(model.num_classes, 3);
        assert_eq!(model.trees.len(), 15);
    }

    #[test]
    fn test_regressor_output() {
        let model = build_default_regressor();
        let features = make_test_features(896);
        let mg_dl = model.predict_regression(&features);
        assert!(mg_dl.is_finite());
    }

    #[test]
    fn test_anomaly_score_normal() {
        let engine = XGInferenceEngine::new("test-v1");
        let features = vec![0.0; 896];
        let pred = engine.predict(&features).unwrap();
        assert!(
            pred.anomaly_score < 0.5,
            "영 벡터 이상점수 < 0.5 (실제: {})",
            pred.anomaly_score
        );
    }

    #[test]
    fn test_anomaly_score_outlier() {
        let engine = XGInferenceEngine::new("test-v1");
        let features = vec![1000.0; 896];
        let pred = engine.predict(&features).unwrap();
        assert!(
            pred.anomaly_score > 0.5,
            "이상치 점수 > 0.5 (실제: {})",
            pred.anomaly_score
        );
    }

    #[test]
    fn test_different_inputs_different_outputs() {
        let engine = XGInferenceEngine::new("test-v1");
        let p1 = engine.predict(&vec![0.1; 896]).unwrap();
        let p2 = engine.predict(&vec![5.0; 896]).unwrap();
        assert!((p1.quantification_mg_dl - p2.quantification_mg_dl).abs() > 0.01);
    }

    #[test]
    fn test_reproducibility() {
        let engine = XGInferenceEngine::new("test-v1");
        let features = make_test_features(896);
        let p1 = engine.predict(&features).unwrap();
        let p2 = engine.predict(&features).unwrap();
        assert_eq!(p1.class_id, p2.class_id);
        assert!((p1.confidence - p2.confidence).abs() < 1e-15);
        assert!((p1.quantification_mg_dl - p2.quantification_mg_dl).abs() < 1e-15);
    }

    #[test]
    fn test_short_input_handled() {
        let engine = XGInferenceEngine::new("test-v1");
        assert!(engine.predict(&vec![1.0; 10]).is_ok());
    }

    #[test]
    fn test_softmax_sums_to_one() {
        let model = build_default_classifier();
        let features = make_test_features(896);
        let nc = model.num_classes as usize;
        let mut scores = vec![model.base_score; nc];
        for (i, tree) in model.trees.iter().enumerate() {
            scores[i % nc] += model.learning_rate * tree.predict(&features);
        }
        let max_s = scores.iter().cloned().fold(f64::NEG_INFINITY, f64::max);
        let exp_s: Vec<f64> = scores.iter().map(|&s| (s - max_s).exp()).collect();
        let sum_exp: f64 = exp_s.iter().sum();
        let total: f64 = exp_s.iter().map(|&e| e / sum_exp).sum();
        assert!((total - 1.0).abs() < 1e-10);
    }

    #[test]
    fn test_custom_model() {
        let tree = DecisionTree {
            nodes: vec![
                TreeNode::Split {
                    feature_idx: 0,
                    threshold: 100.0,
                    left: 1,
                    right: 2,
                },
                TreeNode::Leaf { value: 80.0 },
                TreeNode::Leaf { value: 150.0 },
            ],
        };
        let model = XGBoostModel {
            trees: vec![tree],
            base_score: 10.0,
            num_classes: 1,
            learning_rate: 1.0,
        };
        assert!((model.predict_regression(&[50.0]) - 90.0).abs() < 1e-10);
        assert!((model.predict_regression(&[200.0]) - 160.0).abs() < 1e-10);
    }
}
