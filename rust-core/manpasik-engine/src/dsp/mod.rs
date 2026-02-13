//! DSP (Digital Signal Processing) 모듈
//!
//! 실시간 센서 신호 처리, FFT, 디지털 필터링, 이동평균, 피크 검출
//!
//! # 주요 기능
//! - FFT/역FFT: 주파수 도메인 분석 (rustfft)
//! - 디지털 필터: 저역, 고역, 대역, 노치 필터
//! - 이동평균: 노이즈 제거를 위한 평활화
//! - 피크 검출: 신호 특성 추출
//! - 윈도우 함수: Hamming, Hann, Blackman

use rustfft::{num_complex::Complex, FftPlanner};
use thiserror::Error;

/// DSP 관련 에러
#[derive(Debug, Error)]
pub enum DspError {
    #[error("빈 신호 데이터")]
    EmptySignal,

    #[error("유효하지 않은 신호 길이: {0} (2의 거듭제곱 권장)")]
    InvalidSignalLength(usize),

    #[error("유효하지 않은 주파수: {frequency}Hz (샘플링 레이트 {sample_rate}Hz)")]
    InvalidFrequency { frequency: f64, sample_rate: f64 },

    #[error("유효하지 않은 필터 파라미터: {0}")]
    InvalidFilterParams(String),

    #[error("유효하지 않은 윈도우 크기: {0}")]
    InvalidWindowSize(usize),
}

/// 필터 타입
#[derive(Debug, Clone, Copy, PartialEq)]
pub enum FilterType {
    /// 저역 통과 필터 (cutoff 이하 통과)
    LowPass,
    /// 고역 통과 필터 (cutoff 이상 통과)
    HighPass,
    /// 대역 통과 필터 (low~high 범위 통과)
    BandPass,
    /// 노치 필터 (특정 주파수 제거)
    Notch,
}

/// 윈도우 함수 타입
#[derive(Debug, Clone, Copy)]
pub enum WindowType {
    /// 해밍 윈도우
    Hamming,
    /// 한 윈도우
    Hann,
    /// 블랙만 윈도우
    Blackman,
    /// 직사각형 윈도우 (윈도우 없음)
    Rectangular,
}

/// 필터 파라미터
#[derive(Debug, Clone)]
pub struct FilterParams {
    /// 필터 타입
    pub filter_type: FilterType,
    /// 차단 주파수 (Hz) - LowPass, HighPass에 사용
    pub cutoff_freq: f64,
    /// 상한 주파수 (Hz) - BandPass에 사용
    pub high_freq: Option<f64>,
    /// 대역폭 (Hz) - Notch에 사용
    pub bandwidth: Option<f64>,
    /// 샘플링 레이트 (Hz)
    pub sample_rate: f64,
}

/// 피크 정보
#[derive(Debug, Clone)]
pub struct PeakInfo {
    /// 피크 인덱스
    pub index: usize,
    /// 피크 값 (크기)
    pub amplitude: f64,
    /// 주파수 (Hz, FFT 결과에서 사용)
    pub frequency: Option<f64>,
}

/// FFT 결과
#[derive(Debug, Clone)]
pub struct FftResult {
    /// 주파수 빈 (Hz)
    pub frequencies: Vec<f64>,
    /// 크기 스펙트럼
    pub magnitudes: Vec<f64>,
    /// 위상 스펙트럼 (라디안)
    pub phases: Vec<f64>,
    /// 복소수 결과
    pub complex_spectrum: Vec<Complex<f64>>,
}

/// DSP 프로세서
///
/// 센서 신호의 전처리, 주파수 분석, 노이즈 필터링을 수행합니다.
pub struct DspProcessor {
    /// 샘플링 레이트 (Hz)
    sample_rate: f64,
}

impl DspProcessor {
    /// 새 DSP 프로세서 생성
    ///
    /// # Arguments
    /// * `sample_rate` - 샘플링 레이트 (Hz)
    pub fn new(sample_rate: f64) -> Self {
        Self { sample_rate }
    }

    /// 기본 샘플링 레이트(1000Hz)로 생성
    pub fn with_default_rate() -> Self {
        Self::new(1000.0)
    }

    // =========================================================================
    // FFT (Fast Fourier Transform)
    // =========================================================================

    /// FFT 수행
    ///
    /// 시간 도메인 신호를 주파수 도메인으로 변환합니다.
    ///
    /// # Arguments
    /// * `signal` - 입력 신호 (시간 도메인)
    ///
    /// # Returns
    /// FFT 결과 (주파수, 크기, 위상)
    pub fn fft(&self, signal: &[f64]) -> Result<FftResult, DspError> {
        if signal.is_empty() {
            return Err(DspError::EmptySignal);
        }

        let n = signal.len();
        let mut planner = FftPlanner::new();
        let fft = planner.plan_fft_forward(n);

        // 실수 → 복소수 변환
        let mut buffer: Vec<Complex<f64>> = signal.iter().map(|&x| Complex::new(x, 0.0)).collect();

        fft.process(&mut buffer);

        // 주파수 빈 계산
        let freq_resolution = self.sample_rate / n as f64;
        let frequencies: Vec<f64> = (0..n / 2 + 1).map(|i| i as f64 * freq_resolution).collect();

        // 크기, 위상 계산 (양의 주파수만)
        let half_n = n / 2 + 1;
        let magnitudes: Vec<f64> = buffer[..half_n]
            .iter()
            .map(|c| c.norm() / n as f64 * 2.0)
            .collect();

        let phases: Vec<f64> = buffer[..half_n].iter().map(|c| c.arg()).collect();

        Ok(FftResult {
            frequencies,
            magnitudes,
            phases,
            complex_spectrum: buffer,
        })
    }

    /// 역 FFT 수행
    ///
    /// 주파수 도메인 데이터를 시간 도메인으로 변환합니다.
    pub fn ifft(&self, spectrum: &[Complex<f64>]) -> Result<Vec<f64>, DspError> {
        if spectrum.is_empty() {
            return Err(DspError::EmptySignal);
        }

        let n = spectrum.len();
        let mut planner = FftPlanner::new();
        let ifft = planner.plan_fft_inverse(n);

        let mut buffer = spectrum.to_vec();
        ifft.process(&mut buffer);

        // 정규화 후 실수부 추출
        Ok(buffer.iter().map(|c| c.re / n as f64).collect())
    }

    // =========================================================================
    // 주파수 도메인 필터링
    // =========================================================================

    /// 주파수 도메인 필터 적용
    ///
    /// FFT → 필터 마스크 적용 → 역FFT 방식의 필터링
    pub fn filter(&self, signal: &[f64], params: &FilterParams) -> Result<Vec<f64>, DspError> {
        if signal.is_empty() {
            return Err(DspError::EmptySignal);
        }

        let n = signal.len();
        let mut planner = FftPlanner::new();
        let fft = planner.plan_fft_forward(n);

        // FFT
        let mut buffer: Vec<Complex<f64>> = signal.iter().map(|&x| Complex::new(x, 0.0)).collect();
        fft.process(&mut buffer);

        // 주파수 마스크 생성 및 적용
        let freq_resolution = params.sample_rate / n as f64;

        for (i, sample) in buffer.iter_mut().enumerate() {
            let freq = if i <= n / 2 {
                i as f64 * freq_resolution
            } else {
                (n - i) as f64 * freq_resolution
            };

            let pass = match params.filter_type {
                FilterType::LowPass => freq <= params.cutoff_freq,
                FilterType::HighPass => freq >= params.cutoff_freq,
                FilterType::BandPass => {
                    let high = params.high_freq.unwrap_or(params.cutoff_freq * 2.0);
                    freq >= params.cutoff_freq && freq <= high
                }
                FilterType::Notch => {
                    let bw = params.bandwidth.unwrap_or(10.0);
                    let low = params.cutoff_freq - bw / 2.0;
                    let high = params.cutoff_freq + bw / 2.0;
                    !(freq >= low && freq <= high)
                }
            };

            if !pass {
                *sample = Complex::new(0.0, 0.0);
            }
        }

        // 역FFT
        let ifft = planner.plan_fft_inverse(n);
        ifft.process(&mut buffer);

        // 정규화
        Ok(buffer.iter().map(|c| c.re / n as f64).collect())
    }

    /// 저역 통과 필터 (편의 메서드)
    pub fn low_pass_filter(&self, signal: &[f64], cutoff_hz: f64) -> Result<Vec<f64>, DspError> {
        self.filter(
            signal,
            &FilterParams {
                filter_type: FilterType::LowPass,
                cutoff_freq: cutoff_hz,
                high_freq: None,
                bandwidth: None,
                sample_rate: self.sample_rate,
            },
        )
    }

    /// 고역 통과 필터 (편의 메서드)
    pub fn high_pass_filter(&self, signal: &[f64], cutoff_hz: f64) -> Result<Vec<f64>, DspError> {
        self.filter(
            signal,
            &FilterParams {
                filter_type: FilterType::HighPass,
                cutoff_freq: cutoff_hz,
                high_freq: None,
                bandwidth: None,
                sample_rate: self.sample_rate,
            },
        )
    }

    // =========================================================================
    // 이동평균 (Moving Average)
    // =========================================================================

    /// 단순 이동평균 (SMA)
    ///
    /// 노이즈 제거를 위한 평활화 필터
    ///
    /// # Arguments
    /// * `signal` - 입력 신호
    /// * `window_size` - 윈도우 크기
    pub fn moving_average(&self, signal: &[f64], window_size: usize) -> Result<Vec<f64>, DspError> {
        if signal.is_empty() {
            return Err(DspError::EmptySignal);
        }
        if window_size == 0 || window_size > signal.len() {
            return Err(DspError::InvalidWindowSize(window_size));
        }

        let n = signal.len();
        let mut result = Vec::with_capacity(n);

        // 중심 이동평균 (centered moving average)
        for i in 0..n {
            let start = i.saturating_sub(window_size / 2);
            let end = (i + window_size / 2 + 1).min(n);
            let window_sum: f64 = signal[start..end].iter().sum();
            let count = end - start;
            result.push(window_sum / count as f64);
        }

        Ok(result)
    }

    /// 지수 이동평균 (EMA)
    ///
    /// # Arguments
    /// * `signal` - 입력 신호
    /// * `alpha` - 평활화 계수 (0.0~1.0, 클수록 최근 값에 가중치)
    pub fn exponential_moving_average(
        &self,
        signal: &[f64],
        alpha: f64,
    ) -> Result<Vec<f64>, DspError> {
        if signal.is_empty() {
            return Err(DspError::EmptySignal);
        }
        if !(0.0..=1.0).contains(&alpha) {
            return Err(DspError::InvalidFilterParams(format!(
                "alpha={}, 0.0~1.0 범위여야 합니다",
                alpha
            )));
        }

        let mut result = Vec::with_capacity(signal.len());
        result.push(signal[0]);

        for i in 1..signal.len() {
            let ema = alpha * signal[i] + (1.0 - alpha) * result[i - 1];
            result.push(ema);
        }

        Ok(result)
    }

    // =========================================================================
    // 윈도우 함수
    // =========================================================================

    /// 윈도우 함수 적용
    pub fn apply_window(&self, signal: &[f64], window_type: WindowType) -> Vec<f64> {
        let n = signal.len();
        signal
            .iter()
            .enumerate()
            .map(|(i, &x)| x * self.window_coefficient(i, n, window_type))
            .collect()
    }

    /// 윈도우 계수 계산
    fn window_coefficient(&self, i: usize, n: usize, window_type: WindowType) -> f64 {
        let pi2 = 2.0 * std::f64::consts::PI;
        let ratio = i as f64 / (n - 1) as f64;

        match window_type {
            WindowType::Hamming => 0.54 - 0.46 * (pi2 * ratio).cos(),
            WindowType::Hann => 0.5 * (1.0 - (pi2 * ratio).cos()),
            WindowType::Blackman => {
                0.42 - 0.5 * (pi2 * ratio).cos() + 0.08 * (2.0 * pi2 * ratio).cos()
            }
            WindowType::Rectangular => 1.0,
        }
    }

    // =========================================================================
    // 피크 검출
    // =========================================================================

    /// 피크 검출
    ///
    /// 로컬 최댓값을 찾아 피크 목록을 반환합니다.
    ///
    /// # Arguments
    /// * `signal` - 입력 신호
    /// * `min_amplitude` - 최소 진폭 임계값
    /// * `min_distance` - 피크 간 최소 거리 (샘플 수)
    pub fn find_peaks(
        &self,
        signal: &[f64],
        min_amplitude: f64,
        min_distance: usize,
    ) -> Result<Vec<PeakInfo>, DspError> {
        if signal.is_empty() {
            return Err(DspError::EmptySignal);
        }

        let mut peaks = Vec::new();
        let n = signal.len();

        for i in 1..n - 1 {
            if signal[i] > signal[i - 1] && signal[i] > signal[i + 1] && signal[i] >= min_amplitude
            {
                // 최소 거리 조건 확인
                let too_close = peaks
                    .iter()
                    .any(|p: &PeakInfo| i.abs_diff(p.index) < min_distance);

                if !too_close {
                    peaks.push(PeakInfo {
                        index: i,
                        amplitude: signal[i],
                        frequency: Some(i as f64 * self.sample_rate / n as f64),
                    });
                }
            }
        }

        // 진폭 내림차순 정렬
        peaks.sort_by(|a, b| b.amplitude.partial_cmp(&a.amplitude).unwrap());

        Ok(peaks)
    }

    // =========================================================================
    // 유틸리티
    // =========================================================================

    /// 신호의 RMS (Root Mean Square) 계산
    pub fn rms(&self, signal: &[f64]) -> Result<f64, DspError> {
        if signal.is_empty() {
            return Err(DspError::EmptySignal);
        }
        let sum_sq: f64 = signal.iter().map(|x| x * x).sum();
        Ok((sum_sq / signal.len() as f64).sqrt())
    }

    /// 신호의 SNR (Signal-to-Noise Ratio) 계산 (dB)
    ///
    /// # Arguments
    /// * `signal` - 원본 신호
    /// * `noise` - 노이즈 신호
    pub fn snr_db(&self, signal: &[f64], noise: &[f64]) -> Result<f64, DspError> {
        let signal_rms = self.rms(signal)?;
        let noise_rms = self.rms(noise)?;

        if noise_rms == 0.0 {
            return Ok(f64::INFINITY);
        }

        Ok(20.0 * (signal_rms / noise_rms).log10())
    }

    /// 신호 정규화 (0~1 범위)
    pub fn normalize(&self, signal: &[f64]) -> Result<Vec<f64>, DspError> {
        if signal.is_empty() {
            return Err(DspError::EmptySignal);
        }

        let min = signal.iter().cloned().fold(f64::INFINITY, f64::min);
        let max = signal.iter().cloned().fold(f64::NEG_INFINITY, f64::max);
        let range = max - min;

        if range == 0.0 {
            return Ok(vec![0.5; signal.len()]);
        }

        Ok(signal.iter().map(|&x| (x - min) / range).collect())
    }

    /// 샘플링 레이트 조회
    pub fn sample_rate(&self) -> f64 {
        self.sample_rate
    }
}

impl Default for DspProcessor {
    fn default() -> Self {
        Self::with_default_rate()
    }
}

// =============================================================================
// 테스트
// =============================================================================

#[cfg(test)]
mod tests {
    use super::*;
    use std::f64::consts::PI;

    /// 테스트용 사인파 생성
    fn generate_sine(freq: f64, sample_rate: f64, duration: f64) -> Vec<f64> {
        let num_samples = (sample_rate * duration) as usize;
        (0..num_samples)
            .map(|i| (2.0 * PI * freq * i as f64 / sample_rate).sin())
            .collect()
    }

    /// 복합 신호 생성 (여러 주파수 합성)
    fn generate_composite(
        freqs: &[(f64, f64)], // (주파수, 진폭)
        sample_rate: f64,
        duration: f64,
    ) -> Vec<f64> {
        let num_samples = (sample_rate * duration) as usize;
        (0..num_samples)
            .map(|i| {
                freqs
                    .iter()
                    .map(|(f, a)| a * (2.0 * PI * f * i as f64 / sample_rate).sin())
                    .sum()
            })
            .collect()
    }

    #[test]
    fn test_fft_단일_주파수() {
        let dsp = DspProcessor::new(1000.0);
        let signal = generate_sine(50.0, 1000.0, 1.0); // 50Hz 사인파

        let result = dsp.fft(&signal).unwrap();

        // 50Hz 근처에서 피크 확인
        let peak_idx = result
            .magnitudes
            .iter()
            .enumerate()
            .max_by(|(_, a), (_, b)| a.partial_cmp(b).unwrap())
            .unwrap()
            .0;

        let peak_freq = result.frequencies[peak_idx];
        assert!(
            (peak_freq - 50.0).abs() < 2.0,
            "피크 주파수: {}Hz (예상: 50Hz)",
            peak_freq
        );
    }

    #[test]
    fn test_fft_ifft_왕복() {
        let dsp = DspProcessor::new(1000.0);
        let original = generate_sine(100.0, 1000.0, 0.1);

        // FFT → IFFT
        let fft_result = dsp.fft(&original).unwrap();
        let reconstructed = dsp.ifft(&fft_result.complex_spectrum).unwrap();

        // 원본과 복원된 신호 비교 (오차 허용)
        for (orig, recon) in original.iter().zip(reconstructed.iter()) {
            assert!(
                (orig - recon).abs() < 1e-10,
                "FFT-IFFT 왕복 오차: {} vs {}",
                orig,
                recon
            );
        }
    }

    #[test]
    fn test_저역통과_필터() {
        let dsp = DspProcessor::new(1000.0);

        // 10Hz (신호) + 200Hz (노이즈) 합성
        let signal = generate_composite(&[(10.0, 1.0), (200.0, 0.5)], 1000.0, 1.0);

        // 50Hz 저역통과 필터
        let filtered = dsp.low_pass_filter(&signal, 50.0).unwrap();

        // 필터 후 200Hz 성분이 감쇠되었는지 확인
        let fft_before = dsp.fft(&signal).unwrap();
        let fft_after = dsp.fft(&filtered).unwrap();

        // 200Hz 빈 인덱스
        let idx_200 = (200.0 / (1000.0 / signal.len() as f64)).round() as usize;
        let idx_200 = idx_200.min(fft_after.magnitudes.len() - 1);

        assert!(
            fft_after.magnitudes[idx_200] < fft_before.magnitudes[idx_200] * 0.1,
            "200Hz 성분이 충분히 감쇠되지 않음"
        );
    }

    #[test]
    fn test_이동평균() {
        let dsp = DspProcessor::new(1000.0);
        let signal = vec![1.0, 3.0, 5.0, 7.0, 5.0, 3.0, 1.0];

        let smoothed = dsp.moving_average(&signal, 3).unwrap();

        assert_eq!(smoothed.len(), signal.len());
        // 중앙 값은 윈도우 내 평균에 가까워야 함
        assert!((smoothed[3] - 5.666).abs() < 0.01);
    }

    #[test]
    fn test_지수_이동평균() {
        let dsp = DspProcessor::new(1000.0);
        let signal = vec![1.0, 2.0, 3.0, 4.0, 5.0];

        let ema = dsp.exponential_moving_average(&signal, 0.5).unwrap();

        assert_eq!(ema.len(), signal.len());
        assert_eq!(ema[0], 1.0); // 첫 값은 그대로
        assert!((ema[1] - 1.5).abs() < 1e-10); // 0.5 * 2.0 + 0.5 * 1.0
        assert!((ema[2] - 2.25).abs() < 1e-10); // 0.5 * 3.0 + 0.5 * 1.5
    }

    #[test]
    fn test_지수_이동평균_유효하지_않은_알파() {
        let dsp = DspProcessor::new(1000.0);
        let signal = vec![1.0, 2.0, 3.0];

        let result = dsp.exponential_moving_average(&signal, 1.5);
        assert!(result.is_err());
    }

    #[test]
    fn test_윈도우_함수() {
        let dsp = DspProcessor::new(1000.0);
        let signal = vec![1.0; 100];

        // 해밍 윈도우: 양 끝이 감쇠
        let windowed = dsp.apply_window(&signal, WindowType::Hamming);
        assert!(windowed[0] < 1.0); // 시작부 감쇠
        assert!(windowed[50] > windowed[0]); // 중앙부 높음

        // 직사각형 윈도우: 변화 없음
        let rectangular = dsp.apply_window(&signal, WindowType::Rectangular);
        assert!((rectangular[0] - 1.0).abs() < 1e-10);
    }

    #[test]
    fn test_피크_검출() {
        let dsp = DspProcessor::new(1000.0);

        // 두 개의 명확한 피크가 있는 신호
        let mut signal = vec![0.0; 100];
        signal[20] = 5.0; // 첫 번째 피크
        signal[19] = 2.0;
        signal[21] = 2.0;
        signal[60] = 8.0; // 두 번째 피크
        signal[59] = 3.0;
        signal[61] = 3.0;

        let peaks = dsp.find_peaks(&signal, 1.0, 5).unwrap();

        assert_eq!(peaks.len(), 2);
        assert_eq!(peaks[0].amplitude, 8.0); // 가장 큰 피크
        assert_eq!(peaks[0].index, 60);
        assert_eq!(peaks[1].amplitude, 5.0);
        assert_eq!(peaks[1].index, 20);
    }

    #[test]
    fn test_rms_계산() {
        let dsp = DspProcessor::new(1000.0);

        // DC 신호
        let dc = vec![3.0; 100];
        let rms = dsp.rms(&dc).unwrap();
        assert!((rms - 3.0).abs() < 1e-10);

        // 사인파 RMS = amplitude / sqrt(2)
        let sine = generate_sine(50.0, 1000.0, 1.0);
        let sine_rms = dsp.rms(&sine).unwrap();
        assert!((sine_rms - 1.0 / 2.0_f64.sqrt()).abs() < 0.01);
    }

    #[test]
    fn test_snr_계산() {
        let dsp = DspProcessor::new(1000.0);

        let signal = vec![1.0; 100];
        let noise = vec![0.1; 100];

        let snr = dsp.snr_db(&signal, &noise).unwrap();
        assert!((snr - 20.0).abs() < 0.01); // 20 * log10(1.0/0.1) = 20dB
    }

    #[test]
    fn test_신호_정규화() {
        let dsp = DspProcessor::new(1000.0);
        let signal = vec![-5.0, 0.0, 5.0, 10.0];

        let normalized = dsp.normalize(&signal).unwrap();

        assert!((normalized[0] - 0.0).abs() < 1e-10); // 최솟값 → 0
        assert!((normalized[3] - 1.0).abs() < 1e-10); // 최댓값 → 1
        assert!((normalized[1] - 1.0 / 3.0).abs() < 1e-10);
    }

    #[test]
    fn test_빈_신호_에러() {
        let dsp = DspProcessor::new(1000.0);
        let empty: Vec<f64> = vec![];

        assert!(dsp.fft(&empty).is_err());
        assert!(dsp.moving_average(&empty, 3).is_err());
        assert!(dsp.rms(&empty).is_err());
        assert!(dsp.normalize(&empty).is_err());
    }
}
