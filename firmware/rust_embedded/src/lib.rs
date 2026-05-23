//! MMUP-OS 임베디드 라이브러리
//!
//! STM32F405 타깃 빌드: `cargo build --target thumbv7em-none-eabihf`
//! 호스트 시뮬레이션: `cargo build` (기본 타깃, 테스트용)
//!
//! IEC 62304 Class B · H6 실패 격리 · Result<T,E> 필수

#![cfg_attr(not(test), no_std)]

/// RTOS 태스크 스케줄러
pub mod scheduler {
    include!("../../mmup_os/scheduler.rs");
}

/// AFE 9블록 HAL 드라이버
pub mod afe_hal {
    include!("../../drivers/afe_hal.rs");
}

/// 측정 시퀀서
pub mod sequencer {
    include!("../../measure_seq/sequencer.rs");
}

/// 보안 부팅 + A/B OTA
pub mod secure_boot {
    include!("../../secure/boot.rs");
}
