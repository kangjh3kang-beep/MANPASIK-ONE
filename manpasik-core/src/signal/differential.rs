// signal/differential.rs

/// SSOT 차분식: Sdiff_n = S_n - alpha_n * R_n
/// alpha 는 매트릭스 보정 계수 (기본값 0.98, 범위 0.90~1.10)
pub struct DifferentialEngine {
    alpha: f64,
    alpha_range: (f64, f64),
    // matrix_removal_rate: f64, // 목표치 평가용
}

impl DifferentialEngine {
    pub fn new() -> Self {
        Self {
            alpha: 0.98,
            alpha_range: (0.90, 1.10),
        }
    }

    /// 단일 차동 연산 수행 (S_det - alpha * S_ref)
    pub fn compute(&self, s_detection: f64, s_reference: f64) -> f64 {
        s_detection - self.alpha * s_reference
    }

    /// 동적 alpha 계수 조정 (환경 온도, 습도, 노화도 및 AI 모델 피드백 반영)
    pub fn update_alpha(
        &mut self,
        temperature: f64,
        humidity: f64,
        sensor_age_hours: u32,
        ai_weight: Option<f64>,
    ) -> f64 {
        let temp_factor = 1.0 + (temperature - 25.0) * 0.001;
        let humid_factor = 1.0 + (humidity - 50.0) * 0.0005;
        let aging_factor = 1.0 - (sensor_age_hours as f64) * 0.00001;

        let mut new_alpha = match ai_weight {
            Some(w) => w,
            None => 0.98 * temp_factor * humid_factor * aging_factor,
        };

        // 범위를 SSOT 기준 (0.90 ~ 1.10) 로 제한하여 오도작 방지
        new_alpha = new_alpha.clamp(self.alpha_range.0, self.alpha_range.1);
        self.alpha = new_alpha;
        self.alpha
    }

    /// 4채널 동시 차동 배열 연산
    pub fn compute_array(&self, detections: &[f64; 4], references: &[f64; 4]) -> [f64; 4] {
        let mut results = [0.0; 4];
        for i in 0..4 {
            results[i] = self.compute(detections[i], references[i]);
        }
        results
    }
}
