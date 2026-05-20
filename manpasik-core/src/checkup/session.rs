// src/checkup/session.rs
// ManPaSik (萬波息) v2.4.3 - §6.8 CheckupSession

#[derive(Debug, Clone, PartialEq)]
pub enum CheckupPackage {
    Quick,            // 혈당, 전해질, CRP (1종)
    Standard,         // Quick + 추가 (3종)
    Premium,          // Standard + 다수 (5종)
    Custom(Vec<u16>), // 사용자 선택
}

#[derive(Debug, Clone, PartialEq)]
pub enum CheckupSessionState {
    Preparing,
    InProgress { current_step: u8, total_steps: u8 },
    Analyzing,
    Completed,
}

pub struct CheckupSession {
    pub session_id: String,
    pub package: CheckupPackage,
    pub state: CheckupSessionState,
}

impl CheckupSession {
    pub fn new(session_id: String, package: CheckupPackage) -> Self {
        Self {
            session_id,
            package,
            state: CheckupSessionState::Preparing,
        }
    }

    /// 테스트 패키지별 필요한 카트리지 개수 반환
    pub fn required_cartridge_types(&self) -> u8 {
        match &self.package {
            CheckupPackage::Quick => 1,
            CheckupPackage::Standard => 3,
            CheckupPackage::Premium => 5,
            CheckupPackage::Custom(list) => list.len() as u8,
        }
    }

    /// 단일 카트리지 측정 완료 처리
    pub fn advance_step(&mut self) {
        let total = self.required_cartridge_types();

        match &self.state {
            CheckupSessionState::Preparing => {
                if total > 0 {
                    self.state = CheckupSessionState::InProgress {
                        current_step: 1,
                        total_steps: total,
                    };
                } else {
                    self.state = CheckupSessionState::Completed;
                }
            }
            CheckupSessionState::InProgress {
                current_step,
                total_steps,
            } => {
                if current_step + 1 > *total_steps {
                    self.state = CheckupSessionState::Analyzing; // 분석 단계로 진입
                } else {
                    self.state = CheckupSessionState::InProgress {
                        current_step: current_step + 1,
                        total_steps: *total_steps,
                    };
                }
            }
            CheckupSessionState::Analyzing => {
                self.state = CheckupSessionState::Completed;
            }
            CheckupSessionState::Completed => {}
        }
    }
}
