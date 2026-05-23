//! MMUP-OS 태스크 스케줄러 — RTOS 기반 실시간 태스크 관리
//!
//! IEC 62304 Class B · STM32F405 (ARM Cortex-M4F, 168MHz)
//! FreeRTOS 프로토타입 → SafeRTOS/ThreadX 양산 전환 대상

#![no_std]

/// 태스크 우선순위 (높을수록 우선)
#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord)]
#[repr(u8)]
pub enum TaskPriority {
    Idle = 0,
    Background = 1,      // 로깅, 텔레메트리
    Normal = 2,           // BLE 통신, NFC 읽기
    High = 3,             // 측정 시퀀서
    RealTime = 4,         // ADC 샘플링, 안전 워치독
    Critical = 5,         // 긴급 셧다운, 보안부팅 검증
}

/// 태스크 상태
#[derive(Debug, Clone, Copy, PartialEq)]
pub enum TaskState {
    Ready,
    Running,
    Blocked,
    Suspended,
    Terminated,
}

/// 태스크 정의
pub struct TaskConfig {
    pub name: &'static str,
    pub priority: TaskPriority,
    pub stack_size: usize,
    pub state: TaskState,
}

/// MMUP-OS 태스크 목록 (정적 할당)
pub const TASK_TABLE: &[TaskConfig] = &[
    TaskConfig { name: "watchdog",     priority: TaskPriority::Critical,  stack_size: 256,  state: TaskState::Ready },
    TaskConfig { name: "adc_sampler",  priority: TaskPriority::RealTime,  stack_size: 1024, state: TaskState::Ready },
    TaskConfig { name: "measure_seq",  priority: TaskPriority::High,      stack_size: 2048, state: TaskState::Ready },
    TaskConfig { name: "ble_manager",  priority: TaskPriority::Normal,    stack_size: 2048, state: TaskState::Ready },
    TaskConfig { name: "nfc_reader",   priority: TaskPriority::Normal,    stack_size: 1024, state: TaskState::Ready },
    TaskConfig { name: "power_mgmt",   priority: TaskPriority::Normal,    stack_size: 512,  state: TaskState::Ready },
    TaskConfig { name: "afe_control",  priority: TaskPriority::High,      stack_size: 1024, state: TaskState::Ready },
    TaskConfig { name: "data_logger",  priority: TaskPriority::Background,stack_size: 512,  state: TaskState::Ready },
    TaskConfig { name: "ota_handler",  priority: TaskPriority::Background,stack_size: 1024, state: TaskState::Ready },
    TaskConfig { name: "self_diag",    priority: TaskPriority::Background,stack_size: 512,  state: TaskState::Ready },
];

/// 전원 관리 모드
#[derive(Debug, Clone, Copy)]
pub enum PowerMode {
    Active,       // 전체 동작 (168MHz)
    LowPower,     // 측정 대기 (42MHz, 일부 주변 off)
    Sleep,        // 깊은 절전 (RTC + BLE 웨이크)
    Shutdown,     // 완전 종료 (RTC 알람으로만 복귀)
}

/// 전원 모드 전환 (H6: 실패 시 Active로 폴백)
pub fn transition_power_mode(current: PowerMode, target: PowerMode) -> Result<PowerMode, &'static str> {
    match (current, target) {
        (PowerMode::Active, PowerMode::LowPower) => Ok(PowerMode::LowPower),
        (PowerMode::Active, PowerMode::Sleep) => Ok(PowerMode::Sleep),
        (PowerMode::LowPower, PowerMode::Active) => Ok(PowerMode::Active),
        (PowerMode::Sleep, PowerMode::Active) => Ok(PowerMode::Active),
        (_, PowerMode::Shutdown) => Ok(PowerMode::Shutdown),
        _ => Err("잘못된 전원 모드 전환"),
    }
}
