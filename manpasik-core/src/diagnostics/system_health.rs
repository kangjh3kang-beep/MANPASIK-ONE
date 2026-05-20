// src/diagnostics/system_health.rs
// ManPaSik (萬波息) v2.4.3 - §6.15.6 SystemHealthOrchestrator

#[derive(Clone, Copy, Debug, PartialEq)]
pub enum SystemLayer {
    Hardware = 0,
    Firmware = 1,
    RustCore = 2,
    FlutterApp = 3,
    GoBackend = 4,
    AiInfra = 5,
}

#[derive(Clone, Copy, Debug, PartialEq)]
pub enum HealthGrade {
    Excellent, // >= 0.90
    Good,      // >= 0.75
    Fair,      // >= 0.60
    Poor,      // >= 0.40
    Critical,  // < 0.40
}

#[derive(Clone, Copy, Debug, PartialEq)]
pub enum HealthTrend {
    Improving,
    Stable,
    Declining,
}

pub struct SystemHealthScore {
    pub overall: f64,
    pub grade: HealthGrade,
    pub layers: [(SystemLayer, f64); 6],
    pub timestamp: u64,
    pub trend: HealthTrend,
    pub measurement_ready: bool,
}

pub struct SystemHealthOrchestrator {
    // [HW, FW, Rust, App, Backend, AI] — sum must be 1.0
    layer_weights: [f64; 6],
    last_score: Option<f64>,
}

impl Default for SystemHealthOrchestrator {
    fn default() -> Self {
        Self {
            // [0.30, 0.15, 0.25, 0.10, 0.10, 0.10]
            layer_weights: [0.30, 0.15, 0.25, 0.10, 0.10, 0.10],
            last_score: None,
        }
    }
}

impl SystemHealthOrchestrator {
    pub fn new() -> Self {
        Self::default()
    }

    /// 전체 시스템 건강 점수 계산
    /// IEEE 754 부동소수점 오차 방지를 위해 epsilon 기반 매핑
    pub fn compute_health(&mut self, layer_scores: [f64; 6], timestamp: u64) -> SystemHealthScore {
        let mut overall = 0.0;
        let mut weight_sum = 0.0;
        let mut measurement_ready = true;

        for (i, score) in layer_scores.iter().enumerate() {
            overall += score * self.layer_weights[i];
            weight_sum += self.layer_weights[i];

            // L1 HW 상태가 치명적(Failed)이면 측정 불가
            if i == (SystemLayer::Hardware as usize) && *score < 0.40 {
                measurement_ready = false;
            }
        }

        let normalized_score = overall / weight_sum;
        let epsilon = 1e-9;

        // 등급 매핑
        let grade = if normalized_score >= 0.90 - epsilon {
            HealthGrade::Excellent
        } else if normalized_score >= 0.75 - epsilon {
            HealthGrade::Good
        } else if normalized_score >= 0.60 - epsilon {
            HealthGrade::Fair
        } else if normalized_score >= 0.40 - epsilon {
            HealthGrade::Poor
        } else {
            HealthGrade::Critical
        };

        // 트렌드 분석
        let trend = match self.last_score {
            Some(prev) => {
                let diff = normalized_score - prev;
                if diff >= 0.05 {
                    HealthTrend::Improving
                } else if diff <= -0.05 {
                    HealthTrend::Declining
                } else {
                    HealthTrend::Stable
                }
            }
            None => HealthTrend::Stable,
        };

        self.last_score = Some(normalized_score);

        SystemHealthScore {
            overall: normalized_score,
            grade,
            layers: [
                (SystemLayer::Hardware, layer_scores[0]),
                (SystemLayer::Firmware, layer_scores[1]),
                (SystemLayer::RustCore, layer_scores[2]),
                (SystemLayer::FlutterApp, layer_scores[3]),
                (SystemLayer::GoBackend, layer_scores[4]),
                (SystemLayer::AiInfra, layer_scores[5]),
            ],
            timestamp,
            trend,
            measurement_ready,
        }
    }
}
