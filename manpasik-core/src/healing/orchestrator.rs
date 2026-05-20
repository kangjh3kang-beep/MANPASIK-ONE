// src/healing/orchestrator.rs
// ManPaSik (萬波息) v2.4.3 - §6.16.5 SelfHealingOrchestrator

use crate::diagnostics::system_health::SystemLayer;

pub enum HealingTrigger {
    AutomaticFromDiagnostics,
    UserInitiated,
    ScheduledMaintenance,
    SessionInterruption,
}

pub enum HealingOutcome {
    FullyHealed,
    PartiallyHealed,
    HealingFailed,
    NotNeeded,
}

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
}

pub struct SystemError {
    pub id: String,
    pub layer: SystemLayer,
    pub error_code: u32,
    pub message: String,
    pub auto_recoverable: bool,
}

pub struct SelfHealingOrchestrator {
    // 통계 캐시
    total_attempts_24h: u32,
    success_count: u32,
}

impl SelfHealingOrchestrator {
    pub fn new() -> Self {
        Self {
            total_attempts_24h: 0,
            success_count: 0,
        }
    }

    /// 오류를 수신하여 복구 전략을 지시하는 메인 루프 (단방향 이벤트 처리)
    /// 1. auto_recoverable 확인
    /// 2. 해당 모듈의 Recovery 전략 호출 (모의)
    /// 3. HwAutoRecovery, FwSelfRepair 등의 구체적인 동작 개시
    pub fn heal(&mut self, error: &SystemError, current_score: f64) -> HealingEvent {
        self.total_attempts_24h += 1;

        if !error.auto_recoverable {
            // 자가치유 불가 오류 -> 에스컬레이션(서비스 티켓 등)
            return HealingEvent {
                id: format!("hev-{}", self.total_attempts_24h),
                timestamp: 0, // Mock timestamp
                trigger: HealingTrigger::AutomaticFromDiagnostics,
                layer: error.layer,
                strategy_used: "Escalation to User".into(),
                result: HealingOutcome::HealingFailed,
                duration_ms: 10,
                score_before: current_score,
                score_after: current_score,
            };
        }

        // 복구 시뮬레이션
        self.success_count += 1;
        let new_score = current_score + 0.1; // 점수 회복 시뮬레이션

        HealingEvent {
            id: format!("hev-{}", self.total_attempts_24h),
            timestamp: 0, // Mock
            trigger: HealingTrigger::AutomaticFromDiagnostics,
            layer: error.layer,
            strategy_used: "Auto-Reconnect / Pipeline Reprocess".into(),
            result: HealingOutcome::FullyHealed,
            duration_ms: 320,
            score_before: current_score,
            score_after: if new_score > 1.0 { 1.0 } else { new_score },
        }
    }

    /// 자가치유 성공률 반환
    pub fn success_rate(&self) -> f64 {
        if self.total_attempts_24h == 0 {
            1.0
        } else {
            self.success_count as f64 / self.total_attempts_24h as f64
        }
    }
}
