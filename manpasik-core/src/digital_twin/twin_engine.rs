// digital_twin/twin_engine.rs
//
// 센서 디지털 트윈: 잔차 기반 자가진단 및 교정 시기 예측 엔진.
// 실측값(y)과 예측값(ŷ)의 시계열 잔차를 모니터링하여
// 센서 노화/오염/드리프트를 감지하고 교정 필요 시점을 예측합니다.
//
// 통계적 프로세스 모니터링 기법:
//   - EWMA (지수가중이동평균): 점진적 드리프트 감지
//   - CUSUM (누적합): 평균 이동 조기 탐지
//   - 자기상관 분석: 잔차 패턴 검출
//   - 잔여 수명 추정: 교정까지 남은 측정 횟수

use serde::{Deserialize, Serialize};

/// 센서 건강 상태 (4단계)
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
pub enum SensorHealth {
    Good,     // 정상 범위 내
    Watch,    // 주의 관찰 필요
    Degraded, // 성능 저하 (교정 권장)
    Replace,  // 교체 필요
}

/// 교정 예측 결과
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CalibrationForecast {
    pub health: SensorHealth,
    pub drift_magnitude: f64,
    pub drift_direction: f64,
    pub remaining_measurements: u32,
    pub confidence: f64,
    pub ewma_value: f64,
    pub cusum_pos: f64,
    pub cusum_neg: f64,
}

/// EWMA 제어 차트 상태
#[derive(Debug, Clone)]
struct EwmaState {
    value: f64,
    lambda: f64,
    sigma: f64,
    control_limit_factor: f64,
}

impl EwmaState {
    fn new(lambda: f64, sigma: f64) -> Self {
        Self {
            value: 0.0,
            lambda,
            sigma,
            control_limit_factor: 3.0,
        }
    }

    fn update(&mut self, residual: f64) {
        self.value = self.lambda * residual + (1.0 - self.lambda) * self.value;
    }

    fn control_limit(&self, n: usize) -> f64 {
        if n == 0 {
            return self.sigma * self.control_limit_factor;
        }
        let factor = self.lambda / (2.0 - self.lambda);
        let correction = 1.0 - (1.0 - self.lambda).powi(2 * n as i32);
        self.control_limit_factor * self.sigma * (factor * correction).sqrt()
    }

    fn is_out_of_control(&self, n: usize) -> bool {
        self.value.abs() > self.control_limit(n)
    }
}

/// CUSUM 제어 차트 상태
#[derive(Debug, Clone)]
struct CusumState {
    pos: f64,
    neg: f64,
    slack: f64,
    threshold: f64,
}

impl CusumState {
    fn new(slack: f64, threshold: f64) -> Self {
        Self {
            pos: 0.0,
            neg: 0.0,
            slack,
            threshold,
        }
    }

    fn update(&mut self, residual: f64) {
        self.pos = (self.pos + residual - self.slack).max(0.0);
        self.neg = (self.neg - residual - self.slack).max(0.0);
    }

    fn is_out_of_control(&self) -> bool {
        self.pos > self.threshold || self.neg > self.threshold
    }
}

/// 디지털 트윈 엔진: 센서 자가진단 및 교정 예측
pub struct TwinEngine {
    pub drift_tolerance: f64,
    ewma: EwmaState,
    cusum: CusumState,
    residual_history: Vec<f64>,
    max_history: usize,
    total_measurements: u32,
    max_measurements_before_cal: u32,
}

impl TwinEngine {
    pub fn new(drift_tolerance: f64) -> Self {
        let sigma = drift_tolerance * 0.3;
        Self {
            drift_tolerance,
            ewma: EwmaState::new(0.2, sigma),
            cusum: CusumState::new(sigma * 0.5, drift_tolerance * 4.0),
            residual_history: Vec::with_capacity(256),
            max_history: 256,
            total_measurements: 0,
            max_measurements_before_cal: 500,
        }
    }

    /// 최대 교정 간격 설정
    pub fn set_max_measurements(&mut self, max: u32) {
        self.max_measurements_before_cal = max;
    }

    /// 단일 잔차 입력 → 상태 업데이트 + 건강 진단
    pub fn evaluate_health(&mut self, actual_y: f64, predicted_y_hat: f64) -> SensorHealth {
        let forecast = self.update_and_forecast(actual_y, predicted_y_hat);
        forecast.health
    }

    /// 잔차 입력 → 상태 업데이트 + 상세 교정 예측
    pub fn update_and_forecast(
        &mut self,
        actual_y: f64,
        predicted_y_hat: f64,
    ) -> CalibrationForecast {
        let residual = actual_y - predicted_y_hat;
        self.total_measurements += 1;

        // 잔차 기록
        if self.residual_history.len() >= self.max_history {
            self.residual_history.remove(0);
        }
        self.residual_history.push(residual);

        // EWMA 업데이트
        self.ewma.update(residual);

        // CUSUM 업데이트
        self.cusum.update(residual);

        // 드리프트 크기 및 방향
        let drift_magnitude = self.ewma.value.abs();
        let drift_direction = self.ewma.value.signum();

        // 건강 상태 판정 (다중 기준 통합)
        let health = self.determine_health(drift_magnitude);

        // 잔여 측정 횟수 추정
        let remaining = self.estimate_remaining(drift_magnitude);

        // 판정 신뢰도
        let confidence = self.compute_confidence();

        CalibrationForecast {
            health,
            drift_magnitude,
            drift_direction,
            remaining_measurements: remaining,
            confidence,
            ewma_value: self.ewma.value,
            cusum_pos: self.cusum.pos,
            cusum_neg: self.cusum.neg,
        }
    }

    /// 히스토리 리셋 (교정 후 호출)
    pub fn reset_after_calibration(&mut self) {
        self.ewma.value = 0.0;
        self.cusum.pos = 0.0;
        self.cusum.neg = 0.0;
        self.residual_history.clear();
        self.total_measurements = 0;
    }

    /// 잔차 자기상관 계산 (패턴 탐지용)
    pub fn residual_autocorrelation(&self, lag: usize) -> f64 {
        let data = &self.residual_history;
        if data.len() < lag + 2 {
            return 0.0;
        }
        let n = data.len();
        let mean = data.iter().sum::<f64>() / n as f64;
        let var: f64 = data.iter().map(|x| (x - mean).powi(2)).sum::<f64>() / n as f64;
        if var < 1e-12 {
            return 0.0;
        }
        let cov: f64 = (0..n - lag)
            .map(|i| (data[i] - mean) * (data[i + lag] - mean))
            .sum::<f64>()
            / n as f64;
        cov / var
    }

    // ─── 내부 메서드 ─────────────────────────────────────

    fn determine_health(&self, drift_magnitude: f64) -> SensorHealth {
        let n = self.residual_history.len();

        // 측정 횟수 초과
        if self.total_measurements >= self.max_measurements_before_cal {
            return SensorHealth::Replace;
        }

        // CUSUM 이탈
        if self.cusum.is_out_of_control() {
            return SensorHealth::Replace;
        }

        // EWMA 이탈
        if self.ewma.is_out_of_control(n) {
            return SensorHealth::Degraded;
        }

        // 드리프트 크기 기반
        if drift_magnitude >= self.drift_tolerance {
            return SensorHealth::Replace;
        }
        if drift_magnitude >= self.drift_tolerance * 0.7 {
            return SensorHealth::Degraded;
        }
        if drift_magnitude >= self.drift_tolerance * 0.3 {
            return SensorHealth::Watch;
        }

        // 잔차 자기상관 패턴 검출 (높은 자기상관 = 체계적 드리프트)
        if n >= 10 {
            let ac1 = self.residual_autocorrelation(1);
            if ac1.abs() > 0.7 {
                return SensorHealth::Watch;
            }
        }

        SensorHealth::Good
    }

    fn estimate_remaining(&self, drift_magnitude: f64) -> u32 {
        // 측정 횟수 기반 잔여
        let count_remaining = self
            .max_measurements_before_cal
            .saturating_sub(self.total_measurements);

        // 드리프트 속도 기반 잔여
        let drift_remaining = if drift_magnitude > 1e-12 && self.total_measurements > 5 {
            let rate = drift_magnitude / self.total_measurements as f64;
            if rate > 1e-12 {
                let headroom = self.drift_tolerance - drift_magnitude;
                (headroom / rate).max(0.0) as u32
            } else {
                count_remaining
            }
        } else {
            count_remaining
        };

        // 보수적 추정 (둘 중 작은 값)
        count_remaining.min(drift_remaining)
    }

    fn compute_confidence(&self) -> f64 {
        let n = self.residual_history.len();
        if n < 5 {
            return 0.3;
        }
        if n < 20 {
            return 0.5;
        }
        if n < 50 {
            return 0.7;
        }
        if n < 100 {
            return 0.85;
        }
        0.95
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_good_sensor() {
        let mut engine = TwinEngine::new(1.0);
        for _ in 0..20 {
            let health = engine.evaluate_health(100.0, 100.05);
            assert_eq!(health, SensorHealth::Good);
        }
    }

    #[test]
    fn test_degraded_sensor() {
        let mut engine = TwinEngine::new(1.0);
        // 점진적 드리프트 주입
        for i in 0..50 {
            let drift = i as f64 * 0.02;
            engine.evaluate_health(100.0 + drift, 100.0);
        }
        let forecast = engine.update_and_forecast(101.0, 100.0);
        assert!(
            forecast.health == SensorHealth::Watch
                || forecast.health == SensorHealth::Degraded
                || forecast.health == SensorHealth::Replace,
            "드리프트 누적 후 Watch 이상 상태여야 함: {:?}",
            forecast.health
        );
    }

    #[test]
    fn test_replace_sensor() {
        let mut engine = TwinEngine::new(0.5);
        // 큰 잔차 반복
        for _ in 0..30 {
            engine.evaluate_health(110.0, 100.0);
        }
        let health = engine.evaluate_health(110.0, 100.0);
        assert_eq!(health, SensorHealth::Replace);
    }

    #[test]
    fn test_ewma_tracks_drift() {
        let mut engine = TwinEngine::new(5.0);
        // 일정한 양의 바이어스
        for _ in 0..20 {
            engine.update_and_forecast(101.0, 100.0);
        }
        assert!(
            engine.ewma.value > 0.0,
            "EWMA가 양의 드리프트를 추적해야 합니다"
        );
    }

    #[test]
    fn test_cusum_detects_shift() {
        let mut engine = TwinEngine::new(1.0);
        for _ in 0..50 {
            engine.update_and_forecast(102.0, 100.0);
        }
        assert!(engine.cusum.pos > 0.0, "CUSUM 양수 누적 > 0");
    }

    #[test]
    fn test_reset_after_calibration() {
        let mut engine = TwinEngine::new(1.0);
        for _ in 0..20 {
            engine.update_and_forecast(105.0, 100.0);
        }
        engine.reset_after_calibration();
        assert!((engine.ewma.value).abs() < 1e-12);
        assert!((engine.cusum.pos).abs() < 1e-12);
        assert!(engine.residual_history.is_empty());
        assert_eq!(engine.total_measurements, 0);
    }

    #[test]
    fn test_remaining_measurements_decrease() {
        let mut engine = TwinEngine::new(5.0);
        engine.set_max_measurements(100);

        let f1 = engine.update_and_forecast(100.5, 100.0);
        for _ in 0..50 {
            engine.update_and_forecast(100.5, 100.0);
        }
        let f2 = engine.update_and_forecast(100.5, 100.0);

        assert!(
            f2.remaining_measurements < f1.remaining_measurements,
            "측정 후 잔여 횟수 감소: {} -> {}",
            f1.remaining_measurements,
            f2.remaining_measurements
        );
    }

    #[test]
    fn test_max_measurements_triggers_replace() {
        let mut engine = TwinEngine::new(10.0);
        engine.set_max_measurements(10);
        for _ in 0..10 {
            engine.evaluate_health(100.0, 100.0);
        }
        let health = engine.evaluate_health(100.0, 100.0);
        assert_eq!(
            health,
            SensorHealth::Replace,
            "최대 측정 횟수 초과 → Replace"
        );
    }

    #[test]
    fn test_autocorrelation_random() {
        let mut engine = TwinEngine::new(5.0);
        // 교대 부호 노이즈 (낮은 자기상관)
        for i in 0..100 {
            let sign = if i % 2 == 0 { 1.0 } else { -1.0 };
            let magnitude = ((i * 37 + 17) % 51) as f64 * 0.01;
            engine.update_and_forecast(100.0 + sign * magnitude, 100.0);
        }
        let ac1 = engine.residual_autocorrelation(1);
        assert!(
            ac1 < 0.3,
            "교대 노이즈의 lag-1 자기상관 < 0.3 (실제: {})",
            ac1
        );
    }

    #[test]
    fn test_autocorrelation_systematic() {
        let mut engine = TwinEngine::new(5.0);
        // 체계적 증가 잔차 (높은 자기상관)
        for i in 0..100 {
            let drift = i as f64 * 0.01;
            engine.update_and_forecast(100.0 + drift, 100.0);
        }
        let ac1 = engine.residual_autocorrelation(1);
        assert!(
            ac1 > 0.5,
            "체계적 드리프트의 lag-1 자기상관 > 0.5 (실제: {})",
            ac1
        );
    }

    #[test]
    fn test_confidence_increases_with_data() {
        let mut engine = TwinEngine::new(5.0);
        let f1 = engine.update_and_forecast(100.0, 100.0);
        for _ in 0..100 {
            engine.update_and_forecast(100.0, 100.0);
        }
        let f2 = engine.update_and_forecast(100.0, 100.0);
        assert!(
            f2.confidence > f1.confidence,
            "데이터 누적 후 신뢰도 증가: {} -> {}",
            f1.confidence,
            f2.confidence
        );
    }

    #[test]
    fn test_forecast_fields() {
        let mut engine = TwinEngine::new(2.0);
        let forecast = engine.update_and_forecast(100.5, 100.0);
        assert!(forecast.drift_magnitude >= 0.0);
        assert!(forecast.drift_direction.abs() <= 1.0 || forecast.drift_direction.abs() < 1e-12);
        assert!(forecast.remaining_measurements > 0);
        assert!(forecast.confidence > 0.0 && forecast.confidence <= 1.0);
    }
}
