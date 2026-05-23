//! 측정 시퀀서 — 급전→인큐→읽기→Sdiff 1차→결과 패킷
//!
//! IEC 62304 Class B · H6 실패 격리: 채널별 독립 처리
//! SSOT: Sdiff_n = S_n - α_n × R_n (α 기본 0.98, 범위 0.90~1.10)

#![no_std]

/// SSOT 차동 보정 계수
pub const DEFAULT_ALPHA: f64 = 0.98;
pub const ALPHA_MIN: f64 = 0.90;
pub const ALPHA_MAX: f64 = 1.10;

/// 측정 시퀀스 단계
#[derive(Debug, Clone, Copy, PartialEq)]
pub enum SequencePhase {
    Idle,              // 대기
    PowerOn,           // AFE 급전
    CartridgeDetect,   // NFC 카트리지 인식
    Incubation,        // 인큐베이션 (시료 반응 대기)
    QuietWindow,       // 정숙 구간 (노이즈 최소화)
    Sampling,          // ADC 샘플링
    DiffCorrection,    // 1차 차동 보정
    PacketBuild,       // 결과 패킷 생성
    Complete,          // 완료
    Error,             // 오류 (H6: 전체 중단 금지)
}

/// 채널별 측정 결과
#[derive(Debug, Clone)]
pub struct ChannelResult {
    pub channel_id: u8,
    pub s_det: f64,          // Detection 신호
    pub s_ref: f64,          // Reference 신호
    pub s_diff: f64,         // 차동 보정값
    pub alpha: f64,          // 사용된 α 계수
    pub valid: bool,         // 유효성 (H6: false면 이 채널만 스킵)
    pub error: Option<&'static str>,
}

/// 측정 시퀀서 상태
pub struct MeasureSequencer {
    pub phase: SequencePhase,
    pub alpha: f64,
    pub channels: u8,
    pub results: [Option<ChannelResult>; 16],
    pub incubation_ms: u32,
    pub quiet_window_ms: u32,
}

impl MeasureSequencer {
    /// 새 시퀀서 생성
    pub fn new(channels: u8, alpha: f64) -> Result<Self, &'static str> {
        let clamped = if alpha < ALPHA_MIN {
            ALPHA_MIN
        } else if alpha > ALPHA_MAX {
            ALPHA_MAX
        } else {
            alpha
        };

        if channels == 0 || channels > 16 {
            return Err("채널 수는 1~16 범위");
        }

        Ok(Self {
            phase: SequencePhase::Idle,
            alpha: clamped,
            channels,
            results: Default::default(),
            incubation_ms: 5000,
            quiet_window_ms: 500,
        })
    }

    /// 1차 차동 보정 수행 (H6: 채널별 독립, 실패 시 해당 채널만 스킵)
    pub fn compute_differential(
        &mut self,
        s_det: &[f64],
        s_ref: &[f64],
    ) -> Result<Vec<ChannelResult>, &'static str> {
        if s_det.len() != s_ref.len() {
            return Err("Detection/Reference 채널 수 불일치");
        }

        let mut results = Vec::new();

        for (i, (det, ref_val)) in s_det.iter().zip(s_ref.iter()).enumerate() {
            let channel_id = i as u8;

            // H6: 채널별 독립 처리 — NaN/Inf 감지
            if det.is_nan() || det.is_infinite() || ref_val.is_nan() || ref_val.is_infinite() {
                results.push(ChannelResult {
                    channel_id,
                    s_det: *det,
                    s_ref: *ref_val,
                    s_diff: 0.0,
                    alpha: self.alpha,
                    valid: false,
                    error: Some("비정상 신호값 (NaN/Inf)"),
                });
                continue;
            }

            // Sdiff = S_det - α × S_ref (SSOT 공식)
            let s_diff = det - self.alpha * ref_val;

            results.push(ChannelResult {
                channel_id,
                s_det: *det,
                s_ref: *ref_val,
                s_diff,
                alpha: self.alpha,
                valid: true,
                error: None,
            });
        }

        self.phase = SequencePhase::DiffCorrection;
        Ok(results)
    }

    /// 단계 전진 (상태 머신)
    pub fn advance(&mut self) -> Result<SequencePhase, &'static str> {
        let next = match self.phase {
            SequencePhase::Idle => SequencePhase::PowerOn,
            SequencePhase::PowerOn => SequencePhase::CartridgeDetect,
            SequencePhase::CartridgeDetect => SequencePhase::Incubation,
            SequencePhase::Incubation => SequencePhase::QuietWindow,
            SequencePhase::QuietWindow => SequencePhase::Sampling,
            SequencePhase::Sampling => SequencePhase::DiffCorrection,
            SequencePhase::DiffCorrection => SequencePhase::PacketBuild,
            SequencePhase::PacketBuild => SequencePhase::Complete,
            SequencePhase::Complete => return Err("시퀀스 이미 완료"),
            SequencePhase::Error => return Err("오류 상태 — 리셋 필요"),
        };
        self.phase = next;
        Ok(next)
    }

    /// 유효 채널 수 (H6: 실패 채널 제외)
    pub fn valid_channel_count(&self) -> usize {
        self.results.iter().filter(|r| r.as_ref().map_or(false, |c| c.valid)).count()
    }
}

impl Default for Option<ChannelResult> {
    fn default() -> Self {
        None
    }
}
