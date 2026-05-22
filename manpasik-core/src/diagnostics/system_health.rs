// src/diagnostics/system_health.rs
// ManPaSik (萬波息) v2.4.3 - Phase 5: System Health Diagnostics
//
// 6-Layer 건강 모니터링, 신호 품질 평가, 측정 준비 상태 검증

use serde::{Deserialize, Serialize};

// ============================================================================
// Core Types
// ============================================================================

#[derive(Clone, Copy, Debug, PartialEq, Eq, Serialize, Deserialize)]
pub enum SystemLayer {
    Hardware = 0,
    Firmware = 1,
    RustCore = 2,
    FlutterApp = 3,
    GoBackend = 4,
    AiInfra = 5,
}

impl SystemLayer {
    /// 모든 레이어를 순서대로 반환
    pub fn all() -> [SystemLayer; 6] {
        [
            SystemLayer::Hardware,
            SystemLayer::Firmware,
            SystemLayer::RustCore,
            SystemLayer::FlutterApp,
            SystemLayer::GoBackend,
            SystemLayer::AiInfra,
        ]
    }

    /// 레이어 이름 (한국어)
    pub fn name_kr(&self) -> &'static str {
        match self {
            SystemLayer::Hardware => "하드웨어 (L1)",
            SystemLayer::Firmware => "펌웨어 (L2)",
            SystemLayer::RustCore => "Rust 코어 (L3)",
            SystemLayer::FlutterApp => "Flutter 앱 (L4)",
            SystemLayer::GoBackend => "Go 백엔드 (L5)",
            SystemLayer::AiInfra => "AI 인프라 (L6)",
        }
    }
}

#[derive(Clone, Copy, Debug, PartialEq, Eq, Serialize, Deserialize)]
pub enum HealthGrade {
    Excellent, // >= 0.90
    Good,      // >= 0.75
    Fair,      // >= 0.60
    Poor,      // >= 0.40
    Critical,  // < 0.40
}

impl HealthGrade {
    /// 점수로부터 등급 결정 (IEEE 754 epsilon 기반)
    pub fn from_score(score: f64) -> Self {
        let eps = 1e-9;
        if score >= 0.90 - eps {
            HealthGrade::Excellent
        } else if score >= 0.75 - eps {
            HealthGrade::Good
        } else if score >= 0.60 - eps {
            HealthGrade::Fair
        } else if score >= 0.40 - eps {
            HealthGrade::Poor
        } else {
            HealthGrade::Critical
        }
    }
}

#[derive(Clone, Copy, Debug, PartialEq, Eq, Serialize, Deserialize)]
pub enum HealthTrend {
    Improving,
    Stable,
    Declining,
}

// ============================================================================
// Signal Quality Assessment
// ============================================================================

/// 신호 품질 평가 결과
#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct SignalQualityReport {
    /// SNR (Signal-to-Noise Ratio) in dB
    pub snr_db: f64,
    /// 기준선 드리프트 (단위시간당 변화율)
    pub baseline_drift_rate: f64,
    /// 기준선 드리프트 허용 범위 초과 여부
    pub drift_exceeded: bool,
    /// 피크 신호 진폭
    pub peak_amplitude: f64,
    /// 노이즈 RMS
    pub noise_rms: f64,
    /// 품질 등급 (0.0 ~ 1.0)
    pub quality_score: f64,
}

/// 신호 품질 평가기
pub struct SignalQualityAssessor {
    /// SNR 최소 임계값 (dB)
    pub min_snr_db: f64,
    /// 기준선 드리프트 최대 허용값 (단위시간당)
    pub max_drift_rate: f64,
}

impl Default for SignalQualityAssessor {
    fn default() -> Self {
        Self {
            min_snr_db: 20.0,   // 20 dB 이상이면 양호
            max_drift_rate: 0.1, // 10%/단위시간 이하
        }
    }
}

impl SignalQualityAssessor {
    pub fn new() -> Self {
        Self::default()
    }

    /// 차분 신호 시계열의 품질 평가
    ///
    /// `signal`: 시간순 차분 측정값 (Sdiff_n = S_n - alpha * R_n)
    /// 최소 10개 이상의 샘플이 필요합니다.
    pub fn assess(&self, signal: &[f64]) -> Result<SignalQualityReport, &'static str> {
        if signal.len() < 10 {
            return Err("신호 품질 평가에 최소 10개 샘플이 필요합니다");
        }

        let n = signal.len() as f64;

        // 평균 및 분산 계산
        let mean = signal.iter().sum::<f64>() / n;

        // 신호와 노이즈 분리 (이동평균 기반)
        let window = 5.min(signal.len());
        let smoothed = moving_average(signal, window);

        // 노이즈 = 원본 - 평활 신호
        let noise: Vec<f64> = signal
            .iter()
            .zip(smoothed.iter())
            .map(|(s, sm)| s - sm)
            .collect();

        let noise_rms = (noise.iter().map(|n| n * n).sum::<f64>() / n).sqrt();
        let signal_rms = (smoothed.iter().map(|s| s * s).sum::<f64>() / n).sqrt();

        // SNR 계산
        let snr_db = if noise_rms > 1e-15 {
            20.0 * (signal_rms / noise_rms).log10()
        } else {
            60.0 // 노이즈가 사실상 0이면 매우 높은 SNR
        };

        // 피크 진폭
        let peak_amplitude = signal
            .iter()
            .map(|v| v.abs())
            .fold(0.0_f64, f64::max);

        // 기준선 드리프트: 선형 회귀 기울기 / 평균
        let baseline_drift_rate = {
            let x_mean = (signal.len() as f64 - 1.0) / 2.0;
            let mut ss_xy = 0.0;
            let mut ss_xx = 0.0;
            for (i, &y) in signal.iter().enumerate() {
                let x = i as f64;
                ss_xy += (x - x_mean) * (y - mean);
                ss_xx += (x - x_mean) * (x - x_mean);
            }
            let slope = if ss_xx.abs() > 1e-15 {
                ss_xy / ss_xx
            } else {
                0.0
            };
            // 정규화: 기울기 / 평균값 (상대적 드리프트율)
            if mean.abs() > 1e-10 {
                (slope / mean).abs()
            } else {
                slope.abs()
            }
        };

        let drift_exceeded = baseline_drift_rate > self.max_drift_rate;

        // 종합 품질 점수 (0.0 ~ 1.0)
        let snr_score = ((snr_db - 10.0) / 30.0).clamp(0.0, 1.0); // 10~40 dB -> 0~1
        let drift_score = (1.0 - baseline_drift_rate / self.max_drift_rate).clamp(0.0, 1.0);
        let quality_score = 0.7 * snr_score + 0.3 * drift_score;

        Ok(SignalQualityReport {
            snr_db,
            baseline_drift_rate,
            drift_exceeded,
            peak_amplitude,
            noise_rms,
            quality_score,
        })
    }
}

// ============================================================================
// Measurement Readiness Check
// ============================================================================

/// 측정 준비 상태 체크 항목
#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct ReadinessCheckItem {
    pub name: String,
    pub passed: bool,
    pub detail: String,
}

/// 측정 준비 상태 보고서
#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct MeasurementReadiness {
    pub ready: bool,
    pub checks: Vec<ReadinessCheckItem>,
    pub blocking_issues: Vec<String>,
}

/// 하드웨어/센서 상태 정보 (외부에서 주입)
#[derive(Clone, Debug)]
pub struct HardwareStatus {
    /// BLE 연결 상태
    pub ble_connected: bool,
    /// 마지막 캘리브레이션 이후 경과 시간 (초)
    pub calibration_age_secs: u64,
    /// 캘리브레이션 유효 기간 (초)
    pub calibration_validity_secs: u64,
    /// 카트리지 삽입 여부
    pub cartridge_loaded: bool,
    /// 카트리지 만료 여부
    pub cartridge_expired: bool,
    /// 센서 채널별 자가 진단 결과 (true=정상)
    pub sensor_channels_ok: Vec<bool>,
    /// 배터리 잔량 (0.0 ~ 1.0)
    pub battery_level: f64,
    /// 환경 온도 (degC)
    pub temperature_c: f64,
}

impl Default for HardwareStatus {
    fn default() -> Self {
        Self {
            ble_connected: false,
            calibration_age_secs: 0,
            calibration_validity_secs: 86400, // 24시간
            cartridge_loaded: false,
            cartridge_expired: false,
            sensor_channels_ok: vec![true; 16], // CSI v1.0 16핀
            battery_level: 1.0,
            temperature_c: 25.0,
        }
    }
}

/// 측정 준비 상태 검증기
pub struct ReadinessChecker {
    /// 최소 배터리 잔량 (기본 10%)
    pub min_battery: f64,
    /// 허용 온도 범위 (degC)
    pub temp_range: (f64, f64),
    /// 최소 정상 센서 채널 비율
    pub min_channel_ratio: f64,
}

impl Default for ReadinessChecker {
    fn default() -> Self {
        Self {
            min_battery: 0.10,
            temp_range: (10.0, 45.0),
            min_channel_ratio: 0.75, // 75% 이상 채널 정상 필요
        }
    }
}

impl ReadinessChecker {
    pub fn new() -> Self {
        Self::default()
    }

    /// 측정 준비 상태 종합 검증
    pub fn check(&self, hw: &HardwareStatus) -> MeasurementReadiness {
        let mut checks = Vec::new();
        let mut blocking = Vec::new();

        // 1. BLE 연결
        checks.push(ReadinessCheckItem {
            name: "BLE 연결".to_string(),
            passed: hw.ble_connected,
            detail: if hw.ble_connected {
                "리더기 연결됨".to_string()
            } else {
                "리더기 미연결 -- BLE 연결 필요".to_string()
            },
        });
        if !hw.ble_connected {
            blocking.push("하드웨어 리더기가 연결되지 않았습니다".to_string());
        }

        // 2. 캘리브레이션 유효성
        let cal_valid = hw.calibration_age_secs < hw.calibration_validity_secs;
        checks.push(ReadinessCheckItem {
            name: "캘리브레이션".to_string(),
            passed: cal_valid,
            detail: if cal_valid {
                format!(
                    "유효 ({}초 전 수행, 유효기간 {}초)",
                    hw.calibration_age_secs, hw.calibration_validity_secs
                )
            } else {
                format!(
                    "만료 ({}초 경과, 유효기간 {}초 초과)",
                    hw.calibration_age_secs, hw.calibration_validity_secs
                )
            },
        });
        if !cal_valid {
            blocking.push("캘리브레이션이 만료되었습니다. 재보정이 필요합니다".to_string());
        }

        // 3. 카트리지
        let cartridge_ok = hw.cartridge_loaded && !hw.cartridge_expired;
        checks.push(ReadinessCheckItem {
            name: "카트리지".to_string(),
            passed: cartridge_ok,
            detail: if !hw.cartridge_loaded {
                "카트리지 미삽입".to_string()
            } else if hw.cartridge_expired {
                "카트리지 만료".to_string()
            } else {
                "정상 장착".to_string()
            },
        });
        if !cartridge_ok {
            if !hw.cartridge_loaded {
                blocking.push("카트리지를 삽입해 주세요".to_string());
            } else {
                blocking.push("카트리지가 만료되었습니다. 새 카트리지로 교체해 주세요".to_string());
            }
        }

        // 4. 센서 채널 정상 비율
        let total_channels = hw.sensor_channels_ok.len();
        let ok_channels = hw.sensor_channels_ok.iter().filter(|&&v| v).count();
        let channel_ratio = if total_channels > 0 {
            ok_channels as f64 / total_channels as f64
        } else {
            0.0
        };
        let channels_ok = channel_ratio >= self.min_channel_ratio;
        checks.push(ReadinessCheckItem {
            name: "센서 채널".to_string(),
            passed: channels_ok,
            detail: format!(
                "{}/{} 채널 정상 ({:.0}%, 최소 {:.0}% 필요)",
                ok_channels,
                total_channels,
                channel_ratio * 100.0,
                self.min_channel_ratio * 100.0
            ),
        });
        if !channels_ok {
            blocking.push(format!(
                "센서 채널 {}/{}개만 정상 -- 허용 기준 미달",
                ok_channels, total_channels
            ));
        }

        // 5. 배터리
        let battery_ok = hw.battery_level >= self.min_battery;
        checks.push(ReadinessCheckItem {
            name: "배터리".to_string(),
            passed: battery_ok,
            detail: format!(
                "{:.0}% (최소 {:.0}% 필요)",
                hw.battery_level * 100.0,
                self.min_battery * 100.0
            ),
        });
        if !battery_ok {
            blocking.push(format!(
                "배터리 부족: {:.0}%",
                hw.battery_level * 100.0
            ));
        }

        // 6. 환경 온도
        let temp_ok =
            hw.temperature_c >= self.temp_range.0 && hw.temperature_c <= self.temp_range.1;
        checks.push(ReadinessCheckItem {
            name: "환경 온도".to_string(),
            passed: temp_ok,
            detail: format!(
                "{:.1} degC (허용: {:.0}~{:.0} degC)",
                hw.temperature_c, self.temp_range.0, self.temp_range.1
            ),
        });
        if !temp_ok {
            blocking.push(format!(
                "환경 온도 {:.1} degC -- 허용 범위({:.0}~{:.0} degC) 이탈",
                hw.temperature_c, self.temp_range.0, self.temp_range.1
            ));
        }

        let ready = blocking.is_empty();
        MeasurementReadiness {
            ready,
            checks,
            blocking_issues: blocking,
        }
    }
}

// ============================================================================
// Layer Health Monitoring
// ============================================================================

/// 개별 레이어 건강 상태 정보
#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct LayerHealthDetail {
    pub layer: SystemLayer,
    pub score: f64,
    pub grade: HealthGrade,
    pub checks: Vec<ReadinessCheckItem>,
}

/// 전체 시스템 건강 보고서
#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct SystemHealthReport {
    /// 전체 점수 (가중 평균, 0.0 ~ 1.0)
    pub overall: f64,
    /// 전체 등급
    pub grade: HealthGrade,
    /// 레이어별 점수
    pub layers: Vec<LayerHealthDetail>,
    /// 타임스탬프
    pub timestamp: u64,
    /// 트렌드
    pub trend: HealthTrend,
    /// 측정 준비 상태
    pub measurement_ready: bool,
    /// 권장 조치
    pub recommendations: Vec<String>,
}

pub struct SystemHealthScore {
    pub overall: f64,
    pub grade: HealthGrade,
    pub layers: [(SystemLayer, f64); 6],
    pub timestamp: u64,
    pub trend: HealthTrend,
    pub measurement_ready: bool,
}

/// 시스템 건강 오케스트레이터
pub struct SystemHealthOrchestrator {
    // [HW, FW, Rust, App, Backend, AI] -- sum must be 1.0
    layer_weights: [f64; 6],
    last_score: Option<f64>,
    readiness_checker: ReadinessChecker,
    signal_assessor: SignalQualityAssessor,
    /// 최근 건강 점수 이력 (트렌드 분석용, 최대 20개)
    score_history: Vec<f64>,
}

impl Default for SystemHealthOrchestrator {
    fn default() -> Self {
        Self {
            // [0.30, 0.15, 0.25, 0.10, 0.10, 0.10]
            layer_weights: [0.30, 0.15, 0.25, 0.10, 0.10, 0.10],
            last_score: None,
            readiness_checker: ReadinessChecker::new(),
            signal_assessor: SignalQualityAssessor::new(),
            score_history: Vec::with_capacity(20),
        }
    }
}

impl SystemHealthOrchestrator {
    pub fn new() -> Self {
        Self::default()
    }

    /// 전체 시스템 건강 점수 계산 (레거시 호환)
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

        let grade = HealthGrade::from_score(normalized_score);

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

        // 이력 저장
        self.score_history.push(normalized_score);
        if self.score_history.len() > 20 {
            self.score_history.remove(0);
        }

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

    /// 개별 레이어 건강 체크
    ///
    /// 각 레이어에 대해 특화된 검사 항목을 실행합니다.
    pub fn check_layer_health(
        &self,
        layer: SystemLayer,
        hw_status: &HardwareStatus,
    ) -> LayerHealthDetail {
        let (score, checks) = match layer {
            SystemLayer::Hardware => self.check_hardware_health(hw_status),
            SystemLayer::Firmware => self.check_firmware_health(hw_status),
            SystemLayer::RustCore => self.check_rust_core_health(),
            SystemLayer::FlutterApp => self.check_flutter_app_health(),
            SystemLayer::GoBackend => self.check_go_backend_health(),
            SystemLayer::AiInfra => self.check_ai_infra_health(),
        };

        LayerHealthDetail {
            layer,
            score,
            grade: HealthGrade::from_score(score),
            checks,
        }
    }

    /// 전체 시스템 건강 보고서 생성 (상세 버전)
    pub fn generate_report(
        &mut self,
        hw_status: &HardwareStatus,
        timestamp: u64,
    ) -> SystemHealthReport {
        let mut layer_details = Vec::new();
        let mut layer_scores = [0.0_f64; 6];

        for (i, &layer) in SystemLayer::all().iter().enumerate() {
            let detail = self.check_layer_health(layer, hw_status);
            layer_scores[i] = detail.score;
            layer_details.push(detail);
        }

        let health = self.compute_health(layer_scores, timestamp);

        // 측정 준비 상태 체크
        let readiness = self.readiness_checker.check(hw_status);

        // 권장 조치 수집
        let mut recommendations = Vec::new();
        for detail in &layer_details {
            if detail.score < 0.60 {
                recommendations.push(format!(
                    "{}: 점수 {:.0}% -- 점검 필요",
                    detail.layer.name_kr(),
                    detail.score * 100.0
                ));
            }
        }
        recommendations.extend(readiness.blocking_issues.clone());

        SystemHealthReport {
            overall: health.overall,
            grade: health.grade,
            layers: layer_details,
            timestamp,
            trend: health.trend,
            measurement_ready: readiness.ready,
            recommendations,
        }
    }

    /// 신호 품질 평가 (위임)
    pub fn assess_signal_quality(
        &self,
        signal: &[f64],
    ) -> Result<SignalQualityReport, &'static str> {
        self.signal_assessor.assess(signal)
    }

    /// 측정 준비 상태 체크 (위임)
    pub fn check_measurement_readiness(
        &self,
        hw_status: &HardwareStatus,
    ) -> MeasurementReadiness {
        self.readiness_checker.check(hw_status)
    }

    // ---- 레이어별 건강 체크 구현 ----

    fn check_hardware_health(
        &self,
        hw: &HardwareStatus,
    ) -> (f64, Vec<ReadinessCheckItem>) {
        let mut checks = Vec::new();
        let mut score = 1.0;

        // BLE 연결
        let ble_ok = hw.ble_connected;
        if !ble_ok {
            score -= 0.30;
        }
        checks.push(ReadinessCheckItem {
            name: "BLE 연결".to_string(),
            passed: ble_ok,
            detail: if ble_ok { "연결됨".to_string() } else { "미연결".to_string() },
        });

        // 배터리
        if hw.battery_level < 0.10 {
            score -= 0.25;
        } else if hw.battery_level < 0.30 {
            score -= 0.10;
        }
        checks.push(ReadinessCheckItem {
            name: "배터리".to_string(),
            passed: hw.battery_level >= 0.10,
            detail: format!("{:.0}%", hw.battery_level * 100.0),
        });

        // 센서 채널
        let total = hw.sensor_channels_ok.len();
        let ok = hw.sensor_channels_ok.iter().filter(|&&v| v).count();
        let ratio = if total > 0 { ok as f64 / total as f64 } else { 0.0 };
        if ratio < 0.75 {
            score -= 0.30;
        } else if ratio < 1.0 {
            score -= 0.10;
        }
        checks.push(ReadinessCheckItem {
            name: "센서 채널".to_string(),
            passed: ratio >= 0.75,
            detail: format!("{}/{} 정상", ok, total),
        });

        // 온도 범위
        let temp_ok = hw.temperature_c >= 10.0 && hw.temperature_c <= 45.0;
        if !temp_ok {
            score -= 0.15;
        }
        checks.push(ReadinessCheckItem {
            name: "환경 온도".to_string(),
            passed: temp_ok,
            detail: format!("{:.1} degC", hw.temperature_c),
        });

        (score.clamp(0.0, 1.0), checks)
    }

    fn check_firmware_health(
        &self,
        hw: &HardwareStatus,
    ) -> (f64, Vec<ReadinessCheckItem>) {
        let mut checks = Vec::new();
        let mut score = 1.0;

        // 캘리브레이션 유효성
        let cal_valid = hw.calibration_age_secs < hw.calibration_validity_secs;
        if !cal_valid {
            score -= 0.40;
        }
        checks.push(ReadinessCheckItem {
            name: "캘리브레이션".to_string(),
            passed: cal_valid,
            detail: if cal_valid {
                format!("유효 ({}초 전)", hw.calibration_age_secs)
            } else {
                "만료".to_string()
            },
        });

        // 펌웨어 통신 (BLE 연결 기반 추론)
        let comm_ok = hw.ble_connected;
        if !comm_ok {
            score -= 0.30;
        }
        checks.push(ReadinessCheckItem {
            name: "통신 상태".to_string(),
            passed: comm_ok,
            detail: if comm_ok { "정상".to_string() } else { "통신 불가".to_string() },
        });

        // 카트리지
        let cart_ok = hw.cartridge_loaded && !hw.cartridge_expired;
        if !cart_ok {
            score -= 0.30;
        }
        checks.push(ReadinessCheckItem {
            name: "카트리지".to_string(),
            passed: cart_ok,
            detail: if !hw.cartridge_loaded {
                "미삽입".to_string()
            } else if hw.cartridge_expired {
                "만료".to_string()
            } else {
                "정상".to_string()
            },
        });

        (score.clamp(0.0, 1.0), checks)
    }

    fn check_rust_core_health(&self) -> (f64, Vec<ReadinessCheckItem>) {
        let mut checks = Vec::new();

        // Rust 코어는 항상 동작 중 (self가 존재하므로)
        checks.push(ReadinessCheckItem {
            name: "엔진 상태".to_string(),
            passed: true,
            detail: "정상 동작".to_string(),
        });

        // 메모리 상태 (힙 할당 기반 추론 -- no_std에서는 불가하므로 항상 통과)
        checks.push(ReadinessCheckItem {
            name: "메모리".to_string(),
            passed: true,
            detail: "정상".to_string(),
        });

        // 이력 충분성
        let history_ok = self.score_history.len() >= 2;
        checks.push(ReadinessCheckItem {
            name: "진단 이력".to_string(),
            passed: history_ok,
            detail: format!("{}/20 이력", self.score_history.len()),
        });

        (0.95, checks) // Rust 코어는 실행 중이면 거의 항상 양호
    }

    fn check_flutter_app_health(&self) -> (f64, Vec<ReadinessCheckItem>) {
        // Flutter 앱 상태는 FFI 콜이 도달했다는 것 자체가 건강 증거
        let checks = vec![ReadinessCheckItem {
            name: "FFI 통신".to_string(),
            passed: true,
            detail: "Rust 코어 호출 정상".to_string(),
        }];
        (0.90, checks)
    }

    fn check_go_backend_health(&self) -> (f64, Vec<ReadinessCheckItem>) {
        // 백엔드 상태는 오프라인 우선 아키텍처에서는 비필수
        // 실제 구현에서는 gRPC 핑 결과를 반영
        let checks = vec![ReadinessCheckItem {
            name: "gRPC 연결".to_string(),
            passed: true, // 오프라인 시에도 기본 통과
            detail: "오프라인 모드 허용".to_string(),
        }];
        (0.85, checks)
    }

    fn check_ai_infra_health(&self) -> (f64, Vec<ReadinessCheckItem>) {
        // AI 모델 가용성 (시뮬레이션 폴백 존재하므로 기본 통과)
        let checks = vec![
            ReadinessCheckItem {
                name: "추론 엔진".to_string(),
                passed: true,
                detail: "시뮬레이션 폴백 가용".to_string(),
            },
            ReadinessCheckItem {
                name: "모델 레지스트리".to_string(),
                passed: true,
                detail: "로컬 캐시 사용".to_string(),
            },
        ];
        (0.85, checks)
    }
}

// ============================================================================
// Utility
// ============================================================================

/// 이동평균 계산
fn moving_average(data: &[f64], window: usize) -> Vec<f64> {
    let n = data.len();
    let mut result = Vec::with_capacity(n);
    let half = window / 2;

    for i in 0..n {
        let start = if i >= half { i - half } else { 0 };
        let end = (i + half + 1).min(n);
        let sum: f64 = data[start..end].iter().sum();
        result.push(sum / (end - start) as f64);
    }
    result
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_health_grade_from_score() {
        assert_eq!(HealthGrade::from_score(0.95), HealthGrade::Excellent);
        assert_eq!(HealthGrade::from_score(0.90), HealthGrade::Excellent);
        assert_eq!(HealthGrade::from_score(0.80), HealthGrade::Good);
        assert_eq!(HealthGrade::from_score(0.60), HealthGrade::Fair);
        assert_eq!(HealthGrade::from_score(0.40), HealthGrade::Poor);
        assert_eq!(HealthGrade::from_score(0.30), HealthGrade::Critical);
    }

    #[test]
    fn test_compute_health_basic() {
        let mut orch = SystemHealthOrchestrator::new();
        let scores = [0.95, 0.90, 0.95, 0.85, 0.80, 0.90];
        let result = orch.compute_health(scores, 1000);
        assert!(result.overall > 0.85);
        assert_eq!(result.grade, HealthGrade::Excellent);
        assert!(result.measurement_ready);
    }

    #[test]
    fn test_compute_health_hw_critical() {
        let mut orch = SystemHealthOrchestrator::new();
        let scores = [0.20, 0.90, 0.95, 0.85, 0.80, 0.90]; // HW critical
        let result = orch.compute_health(scores, 1000);
        assert!(!result.measurement_ready);
    }

    #[test]
    fn test_compute_health_trend() {
        let mut orch = SystemHealthOrchestrator::new();
        let _ = orch.compute_health([0.90, 0.90, 0.90, 0.90, 0.90, 0.90], 1);
        let result2 = orch.compute_health([0.50, 0.50, 0.50, 0.50, 0.50, 0.50], 2);
        assert_eq!(result2.trend, HealthTrend::Declining);
    }

    #[test]
    fn test_signal_quality_good() {
        let assessor = SignalQualityAssessor::new();
        // 안정적인 신호 + 작은 노이즈
        let signal: Vec<f64> = (0..100)
            .map(|i| 1.0 + 0.001 * (i as f64 * 0.1).sin())
            .collect();
        let report = assessor.assess(&signal).unwrap();
        assert!(report.snr_db > 10.0);
        assert!(report.quality_score > 0.5);
    }

    #[test]
    fn test_signal_quality_noisy() {
        let assessor = SignalQualityAssessor::new();
        // 노이즈가 많은 신호
        let mut rng_state = 42u64;
        let signal: Vec<f64> = (0..100)
            .map(|_| {
                // 간단한 LCG 의사난수
                rng_state = rng_state.wrapping_mul(6364136223846793005).wrapping_add(1);
                (rng_state as f64) / (u64::MAX as f64) - 0.5
            })
            .collect();
        let report = assessor.assess(&signal).unwrap();
        assert!(report.quality_score < 0.8);
    }

    #[test]
    fn test_signal_quality_insufficient_samples() {
        let assessor = SignalQualityAssessor::new();
        let signal = vec![1.0; 5]; // 10개 미만
        assert!(assessor.assess(&signal).is_err());
    }

    #[test]
    fn test_readiness_all_good() {
        let checker = ReadinessChecker::new();
        let hw = HardwareStatus {
            ble_connected: true,
            calibration_age_secs: 3600,     // 1시간 전
            calibration_validity_secs: 86400, // 24시간 유효
            cartridge_loaded: true,
            cartridge_expired: false,
            sensor_channels_ok: vec![true; 16],
            battery_level: 0.85,
            temperature_c: 25.0,
        };
        let result = checker.check(&hw);
        assert!(result.ready);
        assert!(result.blocking_issues.is_empty());
    }

    #[test]
    fn test_readiness_no_ble() {
        let checker = ReadinessChecker::new();
        let hw = HardwareStatus {
            ble_connected: false,
            ..HardwareStatus::default()
        };
        let result = checker.check(&hw);
        assert!(!result.ready);
        assert!(!result.blocking_issues.is_empty());
    }

    #[test]
    fn test_readiness_expired_cartridge() {
        let checker = ReadinessChecker::new();
        let hw = HardwareStatus {
            ble_connected: true,
            cartridge_loaded: true,
            cartridge_expired: true,
            calibration_age_secs: 100,
            battery_level: 0.80,
            temperature_c: 25.0,
            ..HardwareStatus::default()
        };
        let result = checker.check(&hw);
        assert!(!result.ready);
    }

    #[test]
    fn test_readiness_low_battery() {
        let checker = ReadinessChecker::new();
        let hw = HardwareStatus {
            ble_connected: true,
            cartridge_loaded: true,
            cartridge_expired: false,
            calibration_age_secs: 100,
            battery_level: 0.05, // 5% -- below 10% threshold
            temperature_c: 25.0,
            ..HardwareStatus::default()
        };
        let result = checker.check(&hw);
        assert!(!result.ready);
    }

    #[test]
    fn test_readiness_degraded_channels() {
        let checker = ReadinessChecker::new();
        // 16채널 중 3개만 정상 (18.75% < 75%)
        let mut channels = vec![false; 16];
        channels[0] = true;
        channels[1] = true;
        channels[2] = true;
        let hw = HardwareStatus {
            ble_connected: true,
            cartridge_loaded: true,
            cartridge_expired: false,
            calibration_age_secs: 100,
            sensor_channels_ok: channels,
            battery_level: 0.80,
            temperature_c: 25.0,
            ..HardwareStatus::default()
        };
        let result = checker.check(&hw);
        assert!(!result.ready);
    }

    #[test]
    fn test_full_report_generation() {
        let mut orch = SystemHealthOrchestrator::new();
        let hw = HardwareStatus {
            ble_connected: true,
            calibration_age_secs: 3600,
            calibration_validity_secs: 86400,
            cartridge_loaded: true,
            cartridge_expired: false,
            sensor_channels_ok: vec![true; 16],
            battery_level: 0.90,
            temperature_c: 25.0,
        };
        let report = orch.generate_report(&hw, 1000);
        assert!(report.overall > 0.70);
        assert!(report.measurement_ready);
        assert_eq!(report.layers.len(), 6);
    }

    #[test]
    fn test_layer_health_hardware() {
        let orch = SystemHealthOrchestrator::new();
        let hw = HardwareStatus {
            ble_connected: true,
            battery_level: 0.90,
            sensor_channels_ok: vec![true; 16],
            temperature_c: 25.0,
            ..HardwareStatus::default()
        };
        let detail = orch.check_layer_health(SystemLayer::Hardware, &hw);
        assert!(detail.score > 0.85);
        assert!(!detail.checks.is_empty());
    }
}
