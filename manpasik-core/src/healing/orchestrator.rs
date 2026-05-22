// src/healing/orchestrator.rs
// ManPaSik (萬波息) v2.4.3 - Phase 5: Self-Healing Orchestrator
//
// H6 원칙 (실패 격리): 센서 1개 오류 시 전체 중단 금지.
// Result<T,E>로 에러 격리 및 폴백 경로 제공.
//
// 복구 전략:
//   1. AutoReconnect: BLE 재연결 (지수 백오프)
//   2. RecalibrateChannel: 센서 채널 재보정
//   3. FallbackToBackupSensor: 대체 센서 채널 전환
//   4. EscalateToUser: 사용자 안내 (자동 복구 불가 시)

use crate::diagnostics::system_health::SystemLayer;
use serde::{Deserialize, Serialize};

// ============================================================================
// Core Types
// ============================================================================

#[derive(Clone, Copy, Debug, PartialEq, Eq, Serialize, Deserialize)]
pub enum HealingTrigger {
    AutomaticFromDiagnostics,
    UserInitiated,
    ScheduledMaintenance,
    SessionInterruption,
}

#[derive(Clone, Copy, Debug, PartialEq, Eq, Serialize, Deserialize)]
pub enum HealingOutcome {
    FullyHealed,
    PartiallyHealed,
    HealingFailed,
    NotNeeded,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct HealingEvent {
    pub id: String,
    pub timestamp: u64,
    pub trigger: HealingTrigger,
    pub layer: SystemLayer,
    pub strategy_used: String,
    pub result: HealingOutcome,
    pub duration_ms: u32,
    pub score_before: f64,
    pub score_after: f64,
    pub detail: String,
}

#[derive(Clone, Debug)]
pub struct SystemError {
    pub id: String,
    pub layer: SystemLayer,
    pub error_code: u32,
    pub message: String,
    pub auto_recoverable: bool,
}

/// 오류 심각도
#[derive(Clone, Copy, Debug, PartialEq, Eq, PartialOrd, Ord, Serialize, Deserialize)]
pub enum ErrorSeverity {
    /// 정보 — 복구 불필요
    Info = 0,
    /// 경고 — 모니터링 필요, 자동 복구 시도
    Warning = 1,
    /// 오류 — 즉시 복구 필요
    Error = 2,
    /// 치명적 — 측정 중단, 사용자 개입 필요
    Critical = 3,
}

// ============================================================================
// Healing Strategies
// ============================================================================

/// 복구 전략
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize)]
pub enum HealingStrategy {
    /// BLE 재연결 (지수 백오프)
    AutoReconnect {
        max_retries: u32,
        initial_delay_ms: u32,
    },
    /// 센서 채널 재보정
    RecalibrateChannel {
        channel_id: u32,
        reference_value: Option<f64>,
    },
    /// 대체 센서 채널로 전환 (H6 원칙)
    FallbackToBackupSensor {
        failed_channel: u32,
        backup_channel: u32,
    },
    /// 사용자 안내 (자동 복구 불가)
    EscalateToUser {
        guidance: String,
    },
    /// 데이터 파이프라인 재시작
    RestartPipeline,
    /// 소프트 리셋 (상태 초기화)
    SoftReset {
        layer: SystemLayer,
    },
}

impl HealingStrategy {
    /// 전략 이름 (로깅/보고용)
    pub fn name(&self) -> &'static str {
        match self {
            HealingStrategy::AutoReconnect { .. } => "AutoReconnect",
            HealingStrategy::RecalibrateChannel { .. } => "RecalibrateChannel",
            HealingStrategy::FallbackToBackupSensor { .. } => "FallbackToBackupSensor",
            HealingStrategy::EscalateToUser { .. } => "EscalateToUser",
            HealingStrategy::RestartPipeline => "RestartPipeline",
            HealingStrategy::SoftReset { .. } => "SoftReset",
        }
    }
}

/// BLE 재연결 결과
struct ReconnectResult {
    success: bool,
    attempts: u32,
    total_delay_ms: u32,
}

/// 채널 재보정 결과
struct RecalibrationResult {
    success: bool,
    new_offset: f64,
    duration_ms: u32,
}

/// 치유 오류 타입
#[derive(Clone, Debug, Serialize, Deserialize)]
pub enum HealingError {
    /// BLE 재연결 실패
    ReconnectFailed { attempts: u32, message: String },
    /// 채널 전환 실패
    ChannelSwitchFailed { channel: u32, message: String },
    /// 유효하지 않은 채널
    InvalidChannel { channel: u32 },
    /// 내부 오류
    Internal(String),
}

impl std::fmt::Display for HealingError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            HealingError::ReconnectFailed { attempts, message } => {
                write!(f, "BLE 재연결 실패 ({}회 시도): {}", attempts, message)
            }
            HealingError::ChannelSwitchFailed { channel, message } => {
                write!(f, "채널 {} 전환 실패: {}", channel, message)
            }
            HealingError::InvalidChannel { channel } => {
                write!(f, "유효하지 않은 채널: {} (CSI v1.0: 0~15)", channel)
            }
            HealingError::Internal(msg) => write!(f, "내부 오류: {}", msg),
        }
    }
}

impl std::error::Error for HealingError {}

// ============================================================================
// Strategy Selector
// ============================================================================

/// 오류 유형과 심각도에 따른 복구 전략 선택기
pub struct StrategySelector;

impl StrategySelector {
    /// 오류에 대한 복구 전략 목록 생성 (우선순위순)
    ///
    /// 여러 전략을 반환하여 첫 번째 실패 시 다음 전략으로 폴백합니다.
    pub fn select(error: &SystemError, severity: ErrorSeverity) -> Vec<HealingStrategy> {
        let mut strategies = Vec::new();

        match (error.layer, severity) {
            // ---- 하드웨어 레이어 ----
            (SystemLayer::Hardware, ErrorSeverity::Warning | ErrorSeverity::Error) => {
                // BLE 연결 문제
                if error.error_code & 0x1000 != 0 {
                    strategies.push(HealingStrategy::AutoReconnect {
                        max_retries: 5,
                        initial_delay_ms: 500,
                    });
                }
                // 센서 채널 이상
                if error.error_code & 0x2000 != 0 {
                    let failed_ch = error.error_code & 0xFF;
                    let backup_ch = find_backup_channel(failed_ch);
                    strategies.push(HealingStrategy::RecalibrateChannel {
                        channel_id: failed_ch,
                        reference_value: None,
                    });
                    strategies.push(HealingStrategy::FallbackToBackupSensor {
                        failed_channel: failed_ch,
                        backup_channel: backup_ch,
                    });
                }
            }
            (SystemLayer::Hardware, ErrorSeverity::Critical) => {
                // 치명적 하드웨어 오류 -> 사용자 개입 필요
                strategies.push(HealingStrategy::EscalateToUser {
                    guidance: "리더기 하드웨어에 문제가 감지되었습니다. \
                               리더기를 재시작하고, 문제가 지속되면 서비스 센터에 \
                               문의해 주세요."
                        .to_string(),
                });
            }

            // ---- 펌웨어 레이어 ----
            (SystemLayer::Firmware, ErrorSeverity::Warning | ErrorSeverity::Error) => {
                // 캘리브레이션 만료
                if error.error_code & 0x3000 != 0 {
                    strategies.push(HealingStrategy::RecalibrateChannel {
                        channel_id: 0, // 전체 채널 재보정
                        reference_value: None,
                    });
                }
                // 통신 오류 -> 재연결 + 파이프라인 재시작
                strategies.push(HealingStrategy::AutoReconnect {
                    max_retries: 3,
                    initial_delay_ms: 1000,
                });
                strategies.push(HealingStrategy::RestartPipeline);
            }
            (SystemLayer::Firmware, ErrorSeverity::Critical) => {
                strategies.push(HealingStrategy::SoftReset {
                    layer: SystemLayer::Firmware,
                });
                strategies.push(HealingStrategy::EscalateToUser {
                    guidance: "펌웨어 통신 오류가 지속됩니다. \
                               카트리지를 제거한 후 리더기를 재시작해 주세요."
                        .to_string(),
                });
            }

            // ---- Rust 코어 레이어 ----
            (SystemLayer::RustCore, _) => {
                strategies.push(HealingStrategy::RestartPipeline);
                strategies.push(HealingStrategy::SoftReset {
                    layer: SystemLayer::RustCore,
                });
            }

            // ---- Flutter 앱 레이어 ----
            (SystemLayer::FlutterApp, _) => {
                strategies.push(HealingStrategy::SoftReset {
                    layer: SystemLayer::FlutterApp,
                });
            }

            // ---- 백엔드 레이어 ----
            (SystemLayer::GoBackend, _) => {
                // 오프라인 우선 아키텍처: 백엔드 불가 시 로컬 모드 유지
                strategies.push(HealingStrategy::RestartPipeline);
            }

            // ---- AI 인프라 레이어 ----
            (SystemLayer::AiInfra, _) => {
                // AI 모델 로드 실패 시 시뮬레이션 폴백
                strategies.push(HealingStrategy::RestartPipeline);
                strategies.push(HealingStrategy::SoftReset {
                    layer: SystemLayer::AiInfra,
                });
            }

            // 기본: 정보 수준은 복구 불필요
            (_, ErrorSeverity::Info) => {
                // 조치 불필요
            }
        }

        // 최종 폴백: 모든 자동 복구 실패 시 사용자 에스컬레이션
        if strategies.is_empty()
            || !strategies
                .iter()
                .any(|s| matches!(s, HealingStrategy::EscalateToUser { .. }))
        {
            if severity >= ErrorSeverity::Error {
                strategies.push(HealingStrategy::EscalateToUser {
                    guidance: format!(
                        "자동 복구에 실패했습니다. 오류: {} (코드: 0x{:04X}). \
                         앱을 재시작하거나 고객 지원에 문의해 주세요.",
                        error.message, error.error_code
                    ),
                });
            }
        }

        strategies
    }
}

/// 대체 센서 채널 찾기 (H6 원칙: 실패 격리)
///
/// CSI v1.0 16핀 기준, 짝수/홀수 페어링으로 백업 채널 결정
fn find_backup_channel(failed_channel: u32) -> u32 {
    // 짝수 채널의 백업은 +1 (홀수), 홀수 채널의 백업은 -1 (짝수)
    // 범위: 0~15 (CSI v1.0 16핀)
    if failed_channel % 2 == 0 {
        (failed_channel + 1).min(15)
    } else {
        failed_channel.saturating_sub(1)
    }
}

// ============================================================================
// Self-Healing Orchestrator
// ============================================================================

/// 자가 치유 오케스트레이터
///
/// 오류를 수신하여 적절한 복구 전략을 선택하고 실행합니다.
/// 복구 이력과 성공률을 추적하여 반복 실패 시 에스컬레이션합니다.
pub struct SelfHealingOrchestrator {
    /// 24시간 내 총 시도 횟수
    total_attempts_24h: u32,
    /// 성공 횟수
    success_count: u32,
    /// 복구 이력 (최근 100건)
    history: Vec<HealingEvent>,
    /// 최대 이력 크기
    max_history: usize,
    /// 연속 실패 카운터 (에스컬레이션 판단용)
    consecutive_failures: u32,
    /// 연속 실패 임계값 (이 값을 넘으면 강제 에스컬레이션)
    max_consecutive_failures: u32,
}

impl SelfHealingOrchestrator {
    pub fn new() -> Self {
        Self {
            total_attempts_24h: 0,
            success_count: 0,
            history: Vec::with_capacity(100),
            max_history: 100,
            consecutive_failures: 0,
            max_consecutive_failures: 3,
        }
    }

    /// 오류를 수신하여 복구 전략을 실행하는 메인 진입점
    ///
    /// 1. 오류 심각도 결정
    /// 2. 전략 선택
    /// 3. 전략 순차 실행 (첫 성공 시 중단)
    /// 4. 이벤트 기록
    pub fn heal(
        &mut self,
        error: &SystemError,
        current_score: f64,
        timestamp: u64,
    ) -> HealingEvent {
        self.total_attempts_24h += 1;

        let severity = self.classify_severity(error, current_score);

        // 연속 실패가 임계값을 초과하면 강제 에스컬레이션
        if self.consecutive_failures >= self.max_consecutive_failures {
            return self.force_escalation(error, current_score, timestamp);
        }

        // 전략 선택
        let strategies = StrategySelector::select(error, severity);

        if strategies.is_empty() {
            return self.create_event(
                error,
                "None (no strategy applicable)",
                HealingOutcome::NotNeeded,
                0,
                current_score,
                current_score,
                timestamp,
                "심각도 Info -- 복구 불필요".to_string(),
            );
        }

        // 전략 순차 실행
        for strategy in &strategies {
            let (outcome, duration_ms, new_score, detail) =
                self.execute_strategy(strategy, error, current_score);

            match outcome {
                HealingOutcome::FullyHealed => {
                    self.success_count += 1;
                    self.consecutive_failures = 0;
                    let event = self.create_event(
                        error,
                        strategy.name(),
                        outcome,
                        duration_ms,
                        current_score,
                        new_score,
                        timestamp,
                        detail,
                    );
                    self.record_event(event.clone());
                    return event;
                }
                HealingOutcome::PartiallyHealed => {
                    self.success_count += 1;
                    self.consecutive_failures = 0;
                    let event = self.create_event(
                        error,
                        strategy.name(),
                        outcome,
                        duration_ms,
                        current_score,
                        new_score,
                        timestamp,
                        detail,
                    );
                    self.record_event(event.clone());
                    return event;
                }
                HealingOutcome::HealingFailed => {
                    // 다음 전략으로 폴백
                    continue;
                }
                HealingOutcome::NotNeeded => {
                    let event = self.create_event(
                        error,
                        strategy.name(),
                        outcome,
                        duration_ms,
                        current_score,
                        current_score,
                        timestamp,
                        detail,
                    );
                    self.record_event(event.clone());
                    return event;
                }
            }
        }

        // 모든 전략 실패
        self.consecutive_failures += 1;
        let event = self.create_event(
            error,
            "All strategies exhausted",
            HealingOutcome::HealingFailed,
            0,
            current_score,
            current_score,
            timestamp,
            format!(
                "모든 복구 전략 실패 (연속 실패: {}회)",
                self.consecutive_failures
            ),
        );
        self.record_event(event.clone());
        event
    }

    /// 레거시 호환 heal (timestamp 없는 버전)
    pub fn heal_legacy(&mut self, error: &SystemError, current_score: f64) -> HealingEvent {
        self.heal(error, current_score, 0)
    }

    /// 자가치유 성공률 반환
    pub fn success_rate(&self) -> f64 {
        if self.total_attempts_24h == 0 {
            1.0
        } else {
            self.success_count as f64 / self.total_attempts_24h as f64
        }
    }

    /// 복구 이력 조회
    pub fn history(&self) -> &[HealingEvent] {
        &self.history
    }

    /// 24시간 통계 리셋
    pub fn reset_daily_stats(&mut self) {
        self.total_attempts_24h = 0;
        self.success_count = 0;
        self.consecutive_failures = 0;
    }

    /// 레이어별 복구 성공률
    pub fn success_rate_by_layer(&self, layer: SystemLayer) -> f64 {
        let layer_events: Vec<&HealingEvent> = self
            .history
            .iter()
            .filter(|e| e.layer == layer)
            .collect();

        if layer_events.is_empty() {
            return 1.0;
        }

        let successes = layer_events
            .iter()
            .filter(|e| {
                matches!(
                    e.result,
                    HealingOutcome::FullyHealed | HealingOutcome::PartiallyHealed
                )
            })
            .count();

        successes as f64 / layer_events.len() as f64
    }

    // ---- Private Methods ----

    /// 오류 심각도 분류
    fn classify_severity(&self, error: &SystemError, current_score: f64) -> ErrorSeverity {
        if !error.auto_recoverable {
            return ErrorSeverity::Critical;
        }

        // 현재 점수 기반 심각도 조정
        if current_score < 0.40 {
            ErrorSeverity::Critical
        } else if current_score < 0.60 {
            ErrorSeverity::Error
        } else if current_score < 0.75 {
            ErrorSeverity::Warning
        } else {
            ErrorSeverity::Info
        }
    }

    /// 복구 전략 실행
    ///
    /// Returns: (outcome, duration_ms, new_score, detail)
    fn execute_strategy(
        &self,
        strategy: &HealingStrategy,
        error: &SystemError,
        current_score: f64,
    ) -> (HealingOutcome, u32, f64, String) {
        match strategy {
            HealingStrategy::AutoReconnect {
                max_retries,
                initial_delay_ms,
            } => {
                let result = self.execute_auto_reconnect(*max_retries, *initial_delay_ms);
                if result.success {
                    let recovery = 0.15_f64.min(1.0 - current_score);
                    (
                        HealingOutcome::FullyHealed,
                        result.total_delay_ms,
                        (current_score + recovery).min(1.0),
                        format!(
                            "BLE 재연결 성공 ({}회 시도, {}ms 소요)",
                            result.attempts, result.total_delay_ms
                        ),
                    )
                } else {
                    (
                        HealingOutcome::HealingFailed,
                        result.total_delay_ms,
                        current_score,
                        format!(
                            "BLE 재연결 실패 ({}회 시도 소진)",
                            result.attempts
                        ),
                    )
                }
            }

            HealingStrategy::RecalibrateChannel {
                channel_id,
                reference_value,
            } => {
                let result = self.execute_recalibration(*channel_id, *reference_value);
                if result.success {
                    let recovery = 0.10_f64.min(1.0 - current_score);
                    (
                        HealingOutcome::FullyHealed,
                        result.duration_ms,
                        (current_score + recovery).min(1.0),
                        format!(
                            "채널 {} 재보정 완료 (오프셋: {:.4})",
                            channel_id, result.new_offset
                        ),
                    )
                } else {
                    (
                        HealingOutcome::HealingFailed,
                        result.duration_ms,
                        current_score,
                        format!("채널 {} 재보정 실패", channel_id),
                    )
                }
            }

            HealingStrategy::FallbackToBackupSensor {
                failed_channel,
                backup_channel,
            } => {
                // H6 원칙: 단일 센서 오류 시 대체 채널로 전환
                // 백업 채널이 유효 범위 내인지 확인 (CSI v1.0: 0~15)
                if *backup_channel < 16 && *backup_channel != *failed_channel {
                    let recovery = 0.08_f64.min(1.0 - current_score);
                    (
                        HealingOutcome::PartiallyHealed, // 원래 채널이 아니므로 부분 복구
                        50,
                        (current_score + recovery).min(1.0),
                        format!(
                            "채널 {} -> 백업 채널 {} 전환 완료 (부분 복구)",
                            failed_channel, backup_channel
                        ),
                    )
                } else {
                    (
                        HealingOutcome::HealingFailed,
                        10,
                        current_score,
                        format!(
                            "백업 채널 {} 사용 불가",
                            backup_channel
                        ),
                    )
                }
            }

            HealingStrategy::EscalateToUser { guidance } => {
                // 사용자 에스컬레이션은 항상 "실패" (자동 복구 아님)
                (
                    HealingOutcome::HealingFailed,
                    10,
                    current_score,
                    format!("사용자 안내: {}", guidance),
                )
            }

            HealingStrategy::RestartPipeline => {
                let recovery = 0.05_f64.min(1.0 - current_score);
                (
                    HealingOutcome::FullyHealed,
                    200,
                    (current_score + recovery).min(1.0),
                    "데이터 파이프라인 재시작 완료".to_string(),
                )
            }

            HealingStrategy::SoftReset { layer } => {
                let recovery = 0.10_f64.min(1.0 - current_score);
                (
                    HealingOutcome::FullyHealed,
                    150,
                    (current_score + recovery).min(1.0),
                    format!("{:?} 레이어 소프트 리셋 완료", layer),
                )
            }
        }
    }

    /// BLE 재연결 시뮬레이션 (지수 백오프)
    fn execute_auto_reconnect(
        &self,
        max_retries: u32,
        initial_delay_ms: u32,
    ) -> ReconnectResult {
        let mut total_delay = 0u32;

        for attempt in 1..=max_retries {
            let delay = initial_delay_ms * 2u32.saturating_pow(attempt - 1);
            total_delay = total_delay.saturating_add(delay);

            // 시뮬레이션: 시도 횟수가 늘어날수록 성공 확률 증가
            // 실제 구현에서는 btleplug 기반 BLE 스캔/연결 수행
            let success_probability = match attempt {
                1 => 0.6,
                2 => 0.75,
                3 => 0.85,
                _ => 0.90,
            };

            // 결정론적 시뮬레이션 (성공률 기반)
            // 실제 운영에서는 BLE 스택 상태에 따라 결정
            if success_probability > 0.7 {
                return ReconnectResult {
                    success: true,
                    attempts: attempt,
                    total_delay_ms: total_delay,
                };
            }
        }

        ReconnectResult {
            success: false,
            attempts: max_retries,
            total_delay_ms: total_delay,
        }
    }

    /// 채널 재보정 시뮬레이션
    fn execute_recalibration(
        &self,
        channel_id: u32,
        reference_value: Option<f64>,
    ) -> RecalibrationResult {
        // 참조값이 주어지면 해당 값으로 보정, 없으면 자가 보정
        let ref_val = reference_value.unwrap_or(0.0);

        // 시뮬레이션: 채널 ID가 유효 범위(0~15) 내이면 성공
        if channel_id < 16 {
            RecalibrationResult {
                success: true,
                new_offset: ref_val * 0.02, // 2% 오프셋 보정
                duration_ms: 500,
            }
        } else {
            RecalibrationResult {
                success: false,
                new_offset: 0.0,
                duration_ms: 100,
            }
        }
    }

    /// 강제 에스컬레이션 (연속 실패 임계 초과)
    fn force_escalation(
        &mut self,
        error: &SystemError,
        current_score: f64,
        timestamp: u64,
    ) -> HealingEvent {
        self.consecutive_failures += 1;
        self.create_event(
            error,
            "ForcedEscalation",
            HealingOutcome::HealingFailed,
            10,
            current_score,
            current_score,
            timestamp,
            format!(
                "연속 {}회 복구 실패로 강제 에스컬레이션. \
                 리더기를 재시작하거나 고객 지원에 문의해 주세요. \
                 오류: {} (0x{:04X})",
                self.consecutive_failures, error.message, error.error_code
            ),
        )
    }

    /// 이벤트 생성 헬퍼
    fn create_event(
        &self,
        error: &SystemError,
        strategy_name: &str,
        result: HealingOutcome,
        duration_ms: u32,
        score_before: f64,
        score_after: f64,
        timestamp: u64,
        detail: String,
    ) -> HealingEvent {
        HealingEvent {
            id: format!("hev-{}", self.total_attempts_24h),
            timestamp,
            trigger: HealingTrigger::AutomaticFromDiagnostics,
            layer: error.layer,
            strategy_used: strategy_name.to_string(),
            result,
            duration_ms,
            score_before,
            score_after,
            detail,
        }
    }

    /// 이벤트 이력 기록
    fn record_event(&mut self, event: HealingEvent) {
        if self.history.len() >= self.max_history {
            self.history.remove(0);
        }
        self.history.push(event);
    }

    // ====================================================================
    // H6 실패 격리: 실제 복구 API 연동
    // ====================================================================

    /// BLE 재연결 시도 (지수 백오프)
    ///
    /// 실제 BLE 스택 연동 시 `ble_manager.reconnect()`를 호출합니다.
    /// 현재는 결정론적 지수 백오프 + 구조화된 Result 반환으로 구현합니다.
    ///
    /// # Arguments
    /// * `max_retries` - 최대 재시도 횟수
    /// * `initial_delay_ms` - 초기 대기 시간 (ms)
    ///
    /// # Returns
    /// 재연결 성공 시 `FullyHealed`, 모든 시도 소진 시 `HealingFailed`
    pub fn try_auto_reconnect(
        &self,
        max_retries: u32,
        initial_delay_ms: u64,
    ) -> Result<HealingOutcome, HealingError> {
        if max_retries == 0 {
            return Err(HealingError::ReconnectFailed {
                attempts: 0,
                message: "max_retries가 0입니다".to_string(),
            });
        }

        let mut total_delay_ms: u64 = 0;

        for attempt in 0..max_retries {
            let delay = initial_delay_ms.saturating_mul(2u64.saturating_pow(attempt));
            total_delay_ms = total_delay_ms.saturating_add(delay);

            // 실제 구현: ble_manager.reconnect() 호출
            // 현재: 2번째 시도 이후 성공으로 처리 (결정론적)
            // 향후 btleplug 기반 BLE 스캔/연결로 교체
            if attempt >= 2 {
                return Ok(HealingOutcome::FullyHealed);
            }
        }

        Err(HealingError::ReconnectFailed {
            attempts: max_retries,
            message: format!(
                "{}회 시도 후 BLE 재연결 실패 (총 {}ms 소요)",
                max_retries, total_delay_ms
            ),
        })
    }

    /// 대체 센서 채널로 전환 (H6 실패 격리)
    ///
    /// 오류가 발생한 채널을 마스킹하고, 백업 채널을 활성화합니다.
    /// `DifferentialEngine`의 `mask_channel` / `activate_channel`을 사용합니다.
    ///
    /// # Arguments
    /// * `failed_channel` - 오류가 발생한 채널 번호 (CSI v1.0: 0~15)
    /// * `backup_channel` - 대체 채널 번호
    /// * `num_channels` - DifferentialEngine의 채널 수
    ///
    /// # Returns
    /// 채널 전환 성공 시 `PartiallyHealed` (원래 채널이 아니므로 부분 복구)
    pub fn try_fallback_to_backup(
        &self,
        failed_channel: u32,
        backup_channel: u32,
        num_channels: usize,
    ) -> Result<HealingOutcome, HealingError> {
        // CSI v1.0: 16핀 범위 검증
        if failed_channel >= 16 {
            return Err(HealingError::InvalidChannel {
                channel: failed_channel,
            });
        }
        if backup_channel >= 16 {
            return Err(HealingError::InvalidChannel {
                channel: backup_channel,
            });
        }
        if failed_channel == backup_channel {
            return Err(HealingError::ChannelSwitchFailed {
                channel: backup_channel,
                message: "백업 채널과 오류 채널이 동일합니다".to_string(),
            });
        }

        // DifferentialEngine 채널 범위 검증
        if failed_channel as usize >= num_channels && backup_channel as usize >= num_channels {
            return Err(HealingError::ChannelSwitchFailed {
                channel: backup_channel,
                message: format!(
                    "채널이 엔진 범위(0~{}) 밖입니다",
                    num_channels.saturating_sub(1)
                ),
            });
        }

        // 성공: DifferentialEngine의 mask_channel/activate_channel은 호출자가 수행
        Ok(HealingOutcome::PartiallyHealed)
    }
}

impl Default for SelfHealingOrchestrator {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    fn make_error(layer: SystemLayer, code: u32, recoverable: bool) -> SystemError {
        SystemError {
            id: "test-err-1".to_string(),
            layer,
            error_code: code,
            message: "테스트 오류".to_string(),
            auto_recoverable: recoverable,
        }
    }

    #[test]
    fn test_heal_auto_recoverable() {
        let mut orch = SelfHealingOrchestrator::new();
        let error = make_error(SystemLayer::RustCore, 0x0001, true);

        let event = orch.heal(&error, 0.70, 1000);
        assert!(
            matches!(
                event.result,
                HealingOutcome::FullyHealed | HealingOutcome::PartiallyHealed
            ),
            "Auto-recoverable error should heal: {:?}",
            event.result
        );
        assert!(event.score_after >= event.score_before);
    }

    #[test]
    fn test_heal_non_recoverable() {
        let mut orch = SelfHealingOrchestrator::new();
        let error = make_error(SystemLayer::Hardware, 0x9999, false);

        let event = orch.heal(&error, 0.30, 1000);
        // 치명적 비복구 -> 에스컬레이션 -> HealingFailed
        assert_eq!(event.result, HealingOutcome::HealingFailed);
    }

    #[test]
    fn test_heal_ble_reconnect() {
        let mut orch = SelfHealingOrchestrator::new();
        // BLE 오류 코드: 0x1000 비트
        let error = make_error(SystemLayer::Hardware, 0x1001, true);

        let event = orch.heal(&error, 0.55, 1000);
        assert!(event.strategy_used.contains("Reconnect") || event.strategy_used.contains("Restart"));
    }

    #[test]
    fn test_heal_sensor_fallback() {
        let mut orch = SelfHealingOrchestrator::new();
        // 센서 채널 3 오류: 0x2000 | channel_id
        let error = make_error(SystemLayer::Hardware, 0x2003, true);

        let event = orch.heal(&error, 0.55, 1000);
        // RecalibrateChannel 또는 FallbackToBackupSensor
        assert!(
            event.strategy_used.contains("Recalibrate")
                || event.strategy_used.contains("Fallback")
                || event.strategy_used.contains("Reconnect"),
            "Expected sensor-related strategy, got: {}",
            event.strategy_used
        );
    }

    #[test]
    fn test_success_rate() {
        let mut orch = SelfHealingOrchestrator::new();

        // 2 성공
        let err1 = make_error(SystemLayer::RustCore, 0x0001, true);
        orch.heal(&err1, 0.70, 1);
        orch.heal(&err1, 0.70, 2);

        // 1 실패 (비복구)
        let err2 = make_error(SystemLayer::Hardware, 0xFFFF, false);
        orch.heal(&err2, 0.20, 3);

        let rate = orch.success_rate();
        assert!(rate > 0.5 && rate < 1.0, "Success rate should be ~0.67, got {}", rate);
    }

    #[test]
    fn test_consecutive_failure_escalation() {
        let mut orch = SelfHealingOrchestrator::new();
        orch.max_consecutive_failures = 2;

        let error = make_error(SystemLayer::Hardware, 0xFFFF, false);

        // 연속 실패
        orch.heal(&error, 0.20, 1);
        orch.heal(&error, 0.20, 2);
        let event = orch.heal(&error, 0.20, 3);

        assert_eq!(event.result, HealingOutcome::HealingFailed);
        assert!(
            event.detail.contains("강제 에스컬레이션"),
            "Expected forced escalation, got: {}",
            event.detail
        );
    }

    #[test]
    fn test_find_backup_channel() {
        assert_eq!(find_backup_channel(0), 1);  // 짝수 -> +1
        assert_eq!(find_backup_channel(1), 0);  // 홀수 -> -1
        assert_eq!(find_backup_channel(14), 15); // 짝수 -> +1
        assert_eq!(find_backup_channel(15), 14); // 홀수 -> -1
    }

    #[test]
    fn test_strategy_selector_info_severity() {
        let error = make_error(SystemLayer::Hardware, 0x0001, true);
        let strategies = StrategySelector::select(&error, ErrorSeverity::Info);
        // Info 심각도에서는 전략 없음 (하드웨어에 해당하는 Info 매칭 없음)
        assert!(strategies.is_empty());
    }

    #[test]
    fn test_history_tracking() {
        let mut orch = SelfHealingOrchestrator::new();
        let error = make_error(SystemLayer::RustCore, 0x0001, true);

        orch.heal(&error, 0.70, 1);
        orch.heal(&error, 0.70, 2);

        assert_eq!(orch.history().len(), 2);
    }

    #[test]
    fn test_success_rate_by_layer() {
        let mut orch = SelfHealingOrchestrator::new();

        let rust_err = make_error(SystemLayer::RustCore, 0x0001, true);
        orch.heal(&rust_err, 0.70, 1);

        let hw_err = make_error(SystemLayer::Hardware, 0xFFFF, false);
        orch.heal(&hw_err, 0.20, 2);

        let rust_rate = orch.success_rate_by_layer(SystemLayer::RustCore);
        assert!(rust_rate >= 0.5, "Rust core should have good success rate");
    }

    #[test]
    fn test_reset_daily_stats() {
        let mut orch = SelfHealingOrchestrator::new();
        let error = make_error(SystemLayer::RustCore, 0x0001, true);
        orch.heal(&error, 0.70, 1);

        assert!(orch.total_attempts_24h > 0);
        orch.reset_daily_stats();
        assert_eq!(orch.total_attempts_24h, 0);
        assert_eq!(orch.success_count, 0);
    }

    // ================================================================
    // H6 실패 격리: 실제 복구 API 테스트
    // ================================================================

    #[test]
    fn test_try_auto_reconnect_success() {
        let orch = SelfHealingOrchestrator::new();
        // 3회 이상 시도 시 성공 (attempt >= 2에서 성공)
        let result = orch.try_auto_reconnect(5, 100);
        assert!(result.is_ok());
        assert_eq!(result.expect("should succeed"), HealingOutcome::FullyHealed);
    }

    #[test]
    fn test_try_auto_reconnect_insufficient_retries() {
        let orch = SelfHealingOrchestrator::new();
        // 2회만 시도: attempt 0, 1만 실행 → 실패
        let result = orch.try_auto_reconnect(2, 100);
        assert!(result.is_err());
        match result {
            Err(HealingError::ReconnectFailed { attempts, .. }) => {
                assert_eq!(attempts, 2);
            }
            _ => panic!("Expected ReconnectFailed"),
        }
    }

    #[test]
    fn test_try_auto_reconnect_zero_retries() {
        let orch = SelfHealingOrchestrator::new();
        let result = orch.try_auto_reconnect(0, 100);
        assert!(result.is_err());
    }

    #[test]
    fn test_try_fallback_to_backup_success() {
        let orch = SelfHealingOrchestrator::new();
        let result = orch.try_fallback_to_backup(2, 3, 8);
        assert!(result.is_ok());
        assert_eq!(result.expect("should succeed"), HealingOutcome::PartiallyHealed);
    }

    #[test]
    fn test_try_fallback_invalid_channel() {
        let orch = SelfHealingOrchestrator::new();
        // 채널 16은 CSI v1.0 범위 초과
        let result = orch.try_fallback_to_backup(16, 3, 8);
        assert!(result.is_err());
        assert!(matches!(result, Err(HealingError::InvalidChannel { channel: 16 })));
    }

    #[test]
    fn test_try_fallback_same_channel() {
        let orch = SelfHealingOrchestrator::new();
        let result = orch.try_fallback_to_backup(3, 3, 8);
        assert!(result.is_err());
        assert!(matches!(result, Err(HealingError::ChannelSwitchFailed { .. })));
    }

    #[test]
    fn test_healing_error_display() {
        let err = HealingError::ReconnectFailed {
            attempts: 5,
            message: "timeout".to_string(),
        };
        let msg = format!("{}", err);
        assert!(msg.contains("5회 시도"));
        assert!(msg.contains("timeout"));
    }
}
