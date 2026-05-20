// fingerprint/feature_extractor.rs
//
// 차동측정 파형으로부터 88차원 특성 벡터(핑거프린트)를 추출하는 모듈.
// SSOT: 88→448→896→1792차원 확장 경로의 첫 단계.
//
// 88차원 구성:
//   [0..8]   기본 통계 (peak, valley, mean, std, rms, skewness, kurtosis, median)
//   [8..16]  시간 도메인 (rise_time, fall_time, zero_crossing_rate, peak_count, ...)
//   [16..24] 적분/미분 (AUC, abs_AUC, max_slope, min_slope, ...)
//   [24..32] 형태 분석 (hjorth_activity/mobility/complexity, crest_factor, ...)
//   [32..48] 4-세그먼트 통계 (segment × {mean, std, max, min})
//   [48..56] 백분위수 (p5, p10, p25, p50, p75, p90, p95, IQR)
//   [56..64] 자기상관 (lag 1~8)
//   [64..72] 주파수 프록시 (mean_crossing_rate, dominant_period, spectral_entropy, ...)
//   [72..80] 시간 변화 (half_means, quarter_energies, trend, ...)
//   [80..88] 신호 품질 (baseline, noise_floor, dynamic_range, SNR, ...)

/// 88차원 핑거프린트 특성 벡터 추출기
pub struct FeatureExtractor {
    pub window_size: usize,
}

impl FeatureExtractor {
    pub fn new(window_size: usize) -> Self {
        Self { window_size }
    }

    /// 차동 필터링된 측정 시계열 파형을 입력받아 88차원 특성 벡터를 추출합니다.
    ///
    /// # Arguments
    /// * `differential_wave` - 차동측정 결과 파형 (S_det - alpha * S_ref)
    ///
    /// # Returns
    /// 88차원 특성 벡터. 입력이 비어있으면 모두 0.0.
    pub fn extract(&self, differential_wave: &[f64]) -> [f64; 88] {
        let mut features = [0.0f64; 88];

        if differential_wave.is_empty() {
            return features;
        }

        let n = differential_wave.len();

        // 사전 계산: 1차/2차 미분, 정렬된 사본
        let first_deriv = first_derivative(differential_wave);
        let second_deriv = first_derivative(&first_deriv);
        let mut sorted = differential_wave.to_vec();
        sorted.sort_by(|a, b| a.partial_cmp(b).unwrap_or(std::cmp::Ordering::Equal));

        // === Group 1: 기본 통계 [0..8] ===
        let (mean, variance) = mean_and_variance(differential_wave);
        let std_dev = variance.sqrt();
        let rms = rms_value(differential_wave);

        features[0] = max_val(differential_wave); // peak
        features[1] = min_val(differential_wave); // valley
        features[2] = mean; // mean
        features[3] = std_dev; // std
        features[4] = rms; // RMS
        features[5] = skewness(differential_wave, mean, std_dev); // skewness
        features[6] = kurtosis(differential_wave, mean, std_dev); // kurtosis
        features[7] = percentile_sorted(&sorted, 0.50); // median

        // === Group 2: 시간 도메인 [8..16] ===
        features[8] = rise_time(differential_wave);
        features[9] = fall_time(differential_wave);
        features[10] = zero_crossing_rate(differential_wave);
        features[11] = count_peaks(differential_wave) as f64;
        features[12] = features[0] - features[1]; // peak-to-peak
        features[13] = waveform_length(differential_wave);
        features[14] = energy(differential_wave);
        features[15] = max_position_ratio(differential_wave);

        // === Group 3: 적분/미분 [16..24] ===
        features[16] = trapezoidal_auc(differential_wave);
        features[17] = absolute_auc(differential_wave);
        features[18] = max_val(&first_deriv); // max_slope
        features[19] = min_val(&first_deriv); // min_slope
        features[20] = mean_abs(&first_deriv); // avg_abs_slope
        features[21] = count_sign_changes(&second_deriv) as f64; // inflection_count
        features[22] = max_abs(&second_deriv); // max_second_derivative
        features[23] = {
            let (_, v) = mean_and_variance(&first_deriv);
            v.sqrt() // slope_std
        };

        // === Group 4: 형태 분석 [24..32] ===
        let (h_act, h_mob, h_comp) =
            hjorth_parameters(differential_wave, &first_deriv, &second_deriv);
        let mean_abs_val = mean_abs(differential_wave);
        features[24] = h_act; // hjorth_activity
        features[25] = h_mob; // hjorth_mobility
        features[26] = h_comp; // hjorth_complexity
        features[27] = if rms > 0.0 { features[0] / rms } else { 0.0 }; // crest_factor
        features[28] = if mean_abs_val > 0.0 {
            rms / mean_abs_val
        } else {
            0.0
        }; // form_factor
        features[29] = if mean_abs_val > 0.0 {
            features[0] / mean_abs_val
        } else {
            0.0
        }; // impulse_factor
        features[30] = clearance_factor(differential_wave, features[0]);
        features[31] = symmetry_index(differential_wave);

        // === Group 5: 4-세그먼트 통계 [32..48] ===
        let seg_size = n / 4;
        for seg in 0..4 {
            let start = seg * seg_size;
            let end = if seg == 3 { n } else { start + seg_size };
            let segment = &differential_wave[start..end];
            let (seg_mean, seg_var) = mean_and_variance(segment);
            let base = 32 + seg * 4;
            features[base] = seg_mean;
            features[base + 1] = seg_var.sqrt();
            features[base + 2] = max_val(segment);
            features[base + 3] = min_val(segment);
        }

        // === Group 6: 백분위수 [48..56] ===
        features[48] = percentile_sorted(&sorted, 0.05);
        features[49] = percentile_sorted(&sorted, 0.10);
        features[50] = percentile_sorted(&sorted, 0.25);
        features[51] = percentile_sorted(&sorted, 0.50);
        features[52] = percentile_sorted(&sorted, 0.75);
        features[53] = percentile_sorted(&sorted, 0.90);
        features[54] = percentile_sorted(&sorted, 0.95);
        features[55] = features[52] - features[50]; // IQR (Q3 - Q1)

        // === Group 7: 자기상관 [56..64] ===
        for lag in 1..=8 {
            features[55 + lag] = autocorrelation(differential_wave, mean, variance, lag);
        }

        // === Group 8: 주파수 프록시 [64..72] ===
        features[64] = mean_crossing_rate(differential_wave, mean);
        features[65] = dominant_period_estimate(differential_wave);
        features[66] = amplitude_histogram_entropy(differential_wave, 16);
        features[67] = spectral_centroid_proxy(differential_wave);
        features[68] = high_freq_energy_ratio(&first_deriv, differential_wave);
        features[69] = 1.0 - features[68]; // low_freq_energy_ratio
        features[70] = bandwidth_proxy(differential_wave, mean);
        features[71] = energy(&first_deriv); // first_diff_energy

        // === Group 9: 시간 변화 [72..80] ===
        let half = n / 2;
        features[72] = mean_of_slice(&differential_wave[..half]); // first_half_mean
        features[73] = mean_of_slice(&differential_wave[half..]); // second_half_mean
        let quarter = n / 4;
        features[74] = energy(&differential_wave[..quarter]); // first_quarter_energy
        features[75] = energy(&differential_wave[n - quarter..]); // last_quarter_energy
        let (slope, intercept, r_sq) = linear_regression(differential_wave);
        features[76] = slope; // trend_slope
        features[77] = intercept; // trend_intercept
        features[78] = r_sq; // trend_r_squared
        features[79] = max_consecutive_increase(differential_wave) as f64;

        // === Group 10: 신호 품질 [80..88] ===
        let baseline_len = (n as f64 * 0.1).max(1.0) as usize;
        let tail_len = (n as f64 * 0.2).max(1.0) as usize;
        features[80] = mean_of_slice(&differential_wave[..baseline_len]); // baseline_level
        features[81] = {
            let (_, v) = mean_and_variance(&differential_wave[..baseline_len]);
            v.sqrt() // noise_floor
        };
        let baseline = features[80];
        features[82] = if baseline.abs() > 1e-12 {
            (features[0] - features[1]) / baseline.abs() // dynamic_range
        } else {
            features[0] - features[1]
        };
        features[83] = {
            let tail_start = n.saturating_sub(baseline_len);
            let (_, v) = mean_and_variance(&differential_wave[tail_start..]);
            v.sqrt() // signal_stability
        };
        features[84] = response_time_index(differential_wave, features[0]);
        features[85] = rising_vs_falling_ratio(differential_wave);
        features[86] = {
            let tail_start = n.saturating_sub(tail_len);
            let (s, _, _) = linear_regression(&differential_wave[tail_start..]);
            s // tail_decay_rate
        };
        features[87] = if features[81] > 1e-12 {
            (features[0] - baseline).abs() / features[81] // overall_snr
        } else if (features[0] - baseline).abs() > 1e-12 {
            100.0 // 노이즈 플로어 ≈ 0이면 완벽한 SNR → 100.0 상당
        } else {
            0.0
        };

        // NaN 보호: 모든 차원에서 NaN을 0.0으로 대체
        for f in features.iter_mut() {
            if f.is_nan() || f.is_infinite() {
                *f = 0.0;
            }
        }

        features
    }
}

// ─── 내부 유틸리티 함수 ───────────────────────────────────────────

fn max_val(data: &[f64]) -> f64 {
    data.iter().cloned().fold(f64::NEG_INFINITY, f64::max)
}

fn min_val(data: &[f64]) -> f64 {
    data.iter().cloned().fold(f64::INFINITY, f64::min)
}

fn max_abs(data: &[f64]) -> f64 {
    data.iter().map(|x| x.abs()).fold(0.0f64, f64::max)
}

fn mean_of_slice(data: &[f64]) -> f64 {
    if data.is_empty() {
        return 0.0;
    }
    data.iter().sum::<f64>() / data.len() as f64
}

fn mean_abs(data: &[f64]) -> f64 {
    if data.is_empty() {
        return 0.0;
    }
    data.iter().map(|x| x.abs()).sum::<f64>() / data.len() as f64
}

fn mean_and_variance(data: &[f64]) -> (f64, f64) {
    if data.is_empty() {
        return (0.0, 0.0);
    }
    let n = data.len() as f64;
    let mean = data.iter().sum::<f64>() / n;
    let var = data.iter().map(|x| (x - mean).powi(2)).sum::<f64>() / n;
    (mean, var)
}

fn rms_value(data: &[f64]) -> f64 {
    if data.is_empty() {
        return 0.0;
    }
    (data.iter().map(|x| x * x).sum::<f64>() / data.len() as f64).sqrt()
}

fn energy(data: &[f64]) -> f64 {
    data.iter().map(|x| x * x).sum::<f64>()
}

fn skewness(data: &[f64], mean: f64, std: f64) -> f64 {
    if std < 1e-12 || data.len() < 3 {
        return 0.0;
    }
    let n = data.len() as f64;
    let m3 = data.iter().map(|x| ((x - mean) / std).powi(3)).sum::<f64>() / n;
    m3
}

fn kurtosis(data: &[f64], mean: f64, std: f64) -> f64 {
    if std < 1e-12 || data.len() < 4 {
        return 0.0;
    }
    let n = data.len() as f64;
    let m4 = data.iter().map(|x| ((x - mean) / std).powi(4)).sum::<f64>() / n;
    m4 - 3.0 // excess kurtosis
}

fn first_derivative(data: &[f64]) -> Vec<f64> {
    if data.len() < 2 {
        return vec![0.0];
    }
    data.windows(2).map(|w| w[1] - w[0]).collect()
}

fn percentile_sorted(sorted: &[f64], p: f64) -> f64 {
    if sorted.is_empty() {
        return 0.0;
    }
    let idx = (p * (sorted.len() - 1) as f64).round() as usize;
    sorted[idx.min(sorted.len() - 1)]
}

fn zero_crossing_rate(data: &[f64]) -> f64 {
    if data.len() < 2 {
        return 0.0;
    }
    let crossings = data
        .windows(2)
        .filter(|w| (w[0] >= 0.0 && w[1] < 0.0) || (w[0] < 0.0 && w[1] >= 0.0))
        .count();
    crossings as f64 / (data.len() - 1) as f64
}

fn mean_crossing_rate(data: &[f64], mean: f64) -> f64 {
    if data.len() < 2 {
        return 0.0;
    }
    let crossings = data
        .windows(2)
        .filter(|w| (w[0] >= mean && w[1] < mean) || (w[0] < mean && w[1] >= mean))
        .count();
    crossings as f64 / (data.len() - 1) as f64
}

fn count_peaks(data: &[f64]) -> usize {
    if data.len() < 3 {
        return 0;
    }
    data.windows(3)
        .filter(|w| w[1] > w[0] && w[1] > w[2])
        .count()
}

fn count_sign_changes(data: &[f64]) -> usize {
    if data.len() < 2 {
        return 0;
    }
    data.windows(2)
        .filter(|w| (w[0] > 0.0 && w[1] < 0.0) || (w[0] < 0.0 && w[1] > 0.0))
        .count()
}

fn waveform_length(data: &[f64]) -> f64 {
    data.windows(2).map(|w| (w[1] - w[0]).abs()).sum::<f64>()
}

fn max_position_ratio(data: &[f64]) -> f64 {
    if data.is_empty() {
        return 0.0;
    }
    let max_idx = data
        .iter()
        .enumerate()
        .max_by(|(_, a), (_, b)| a.partial_cmp(b).unwrap_or(std::cmp::Ordering::Equal))
        .map(|(i, _)| i)
        .unwrap_or(0);
    max_idx as f64 / data.len() as f64
}

fn trapezoidal_auc(data: &[f64]) -> f64 {
    if data.len() < 2 {
        return 0.0;
    }
    data.windows(2).map(|w| (w[0] + w[1]) * 0.5).sum::<f64>()
}

fn absolute_auc(data: &[f64]) -> f64 {
    if data.len() < 2 {
        return 0.0;
    }
    data.windows(2)
        .map(|w| (w[0].abs() + w[1].abs()) * 0.5)
        .sum::<f64>()
}

/// 10%→90% 도달 시간 (상승 구간)
fn rise_time(data: &[f64]) -> f64 {
    if data.is_empty() {
        return 0.0;
    }
    let peak = max_val(data);
    let valley = min_val(data);
    let range = peak - valley;
    if range < 1e-12 {
        return 0.0;
    }
    let t10 = valley + 0.1 * range;
    let t90 = valley + 0.9 * range;

    let mut i10 = None;
    let mut i90 = None;
    for (i, &v) in data.iter().enumerate() {
        if i10.is_none() && v >= t10 {
            i10 = Some(i);
        }
        if i10.is_some() && i90.is_none() && v >= t90 {
            i90 = Some(i);
            break;
        }
    }
    match (i10, i90) {
        (Some(a), Some(b)) => (b - a) as f64 / data.len() as f64,
        _ => 0.0,
    }
}

/// 90%→10% 하강 시간 (하강 구간, 피크 이후)
fn fall_time(data: &[f64]) -> f64 {
    if data.is_empty() {
        return 0.0;
    }
    let peak_idx = data
        .iter()
        .enumerate()
        .max_by(|(_, a), (_, b)| a.partial_cmp(b).unwrap_or(std::cmp::Ordering::Equal))
        .map(|(i, _)| i)
        .unwrap_or(0);

    let tail = &data[peak_idx..];
    if tail.is_empty() {
        return 0.0;
    }

    let peak = max_val(data);
    let valley = min_val(data);
    let range = peak - valley;
    if range < 1e-12 {
        return 0.0;
    }
    let t90 = valley + 0.9 * range;
    let t10 = valley + 0.1 * range;

    let mut i90 = None;
    let mut i10 = None;
    for (i, &v) in tail.iter().enumerate() {
        if i90.is_none() && v <= t90 {
            i90 = Some(i);
        }
        if i90.is_some() && i10.is_none() && v <= t10 {
            i10 = Some(i);
            break;
        }
    }
    match (i90, i10) {
        (Some(a), Some(b)) => (b - a) as f64 / data.len() as f64,
        _ => 0.0,
    }
}

/// Hjorth 파라미터 (Activity, Mobility, Complexity)
fn hjorth_parameters(data: &[f64], d1: &[f64], d2: &[f64]) -> (f64, f64, f64) {
    let (_, var0) = mean_and_variance(data);
    let (_, var1) = mean_and_variance(d1);
    let (_, var2) = mean_and_variance(d2);

    let activity = var0;
    let mobility = if var0 > 1e-12 {
        (var1 / var0).sqrt()
    } else {
        0.0
    };
    let mob_d1 = if var1 > 1e-12 {
        (var2 / var1).sqrt()
    } else {
        0.0
    };
    let complexity = if mobility > 1e-12 {
        mob_d1 / mobility
    } else {
        0.0
    };

    (activity, mobility, complexity)
}

fn clearance_factor(data: &[f64], peak: f64) -> f64 {
    if data.is_empty() {
        return 0.0;
    }
    let mean_sqrt = data.iter().map(|x| x.abs().sqrt()).sum::<f64>() / data.len() as f64;
    let denom = mean_sqrt * mean_sqrt;
    if denom > 1e-12 {
        peak / denom
    } else {
        0.0
    }
}

fn symmetry_index(data: &[f64]) -> f64 {
    if data.len() < 2 {
        return 0.0;
    }
    let n = data.len();
    let half = n / 2;
    let first_energy: f64 = data[..half].iter().map(|x| x * x).sum();
    let second_energy: f64 = data[n - half..].iter().map(|x| x * x).sum();
    let total = first_energy + second_energy;
    if total > 1e-12 {
        first_energy / total
    } else {
        0.5
    }
}

fn autocorrelation(data: &[f64], mean: f64, variance: f64, lag: usize) -> f64 {
    if lag >= data.len() || variance < 1e-12 {
        return 0.0;
    }
    let n = data.len();
    let sum: f64 = (0..n - lag)
        .map(|i| (data[i] - mean) * (data[i + lag] - mean))
        .sum();
    sum / (n as f64 * variance)
}

fn dominant_period_estimate(data: &[f64]) -> f64 {
    if data.len() < 3 {
        return 0.0;
    }
    let mut peak_positions = Vec::new();
    for i in 1..data.len() - 1 {
        if data[i] > data[i - 1] && data[i] > data[i + 1] {
            peak_positions.push(i);
        }
    }
    if peak_positions.len() < 2 {
        return data.len() as f64;
    }
    let total_gaps: usize = peak_positions.windows(2).map(|w| w[1] - w[0]).sum();
    total_gaps as f64 / (peak_positions.len() - 1) as f64
}

/// 진폭 히스토그램 기반 엔트로피 (주파수 영역 프록시)
fn amplitude_histogram_entropy(data: &[f64], bins: usize) -> f64 {
    if data.is_empty() || bins == 0 {
        return 0.0;
    }
    let mn = min_val(data);
    let mx = max_val(data);
    let range = mx - mn;
    if range < 1e-12 {
        return 0.0;
    }

    let mut hist = vec![0usize; bins];
    for &v in data {
        let idx = ((v - mn) / range * (bins - 1) as f64).round() as usize;
        hist[idx.min(bins - 1)] += 1;
    }

    let n = data.len() as f64;
    let mut entropy = 0.0;
    for &count in &hist {
        if count > 0 {
            let p = count as f64 / n;
            entropy -= p * p.ln();
        }
    }
    entropy / (bins as f64).ln() // 정규화 [0, 1]
}

fn spectral_centroid_proxy(data: &[f64]) -> f64 {
    if data.len() < 3 {
        return 0.0;
    }
    let mut peak_positions = Vec::new();
    for i in 1..data.len() - 1 {
        if data[i] > data[i - 1] && data[i] > data[i + 1] {
            peak_positions.push(i);
        }
    }
    if peak_positions.is_empty() {
        return 0.0;
    }
    let weighted_sum: f64 = peak_positions
        .iter()
        .map(|&i| i as f64 * data[i].abs())
        .sum();
    let weight_total: f64 = peak_positions.iter().map(|&i| data[i].abs()).sum();
    if weight_total > 1e-12 {
        weighted_sum / weight_total / data.len() as f64
    } else {
        0.0
    }
}

fn high_freq_energy_ratio(deriv: &[f64], original: &[f64]) -> f64 {
    let deriv_energy = energy(deriv);
    let orig_energy = energy(original);
    if orig_energy > 1e-12 {
        (deriv_energy / orig_energy).min(1.0)
    } else {
        0.0
    }
}

fn bandwidth_proxy(data: &[f64], mean: f64) -> f64 {
    if data.len() < 2 {
        return 0.0;
    }
    let mut intervals = Vec::new();
    let mut above = data[0] >= mean;
    let mut last_cross = 0;
    for i in 1..data.len() {
        let now_above = data[i] >= mean;
        if now_above != above {
            if last_cross > 0 {
                intervals.push(i - last_cross);
            }
            last_cross = i;
            above = now_above;
        }
    }
    if intervals.is_empty() {
        return 0.0;
    }
    let mean_interval = intervals.iter().sum::<usize>() as f64 / intervals.len() as f64;
    let var_interval: f64 = intervals
        .iter()
        .map(|&x| (x as f64 - mean_interval).powi(2))
        .sum::<f64>()
        / intervals.len() as f64;
    var_interval.sqrt() / mean_interval.max(1e-12)
}

fn linear_regression(data: &[f64]) -> (f64, f64, f64) {
    let n = data.len() as f64;
    if n < 2.0 {
        return (0.0, 0.0, 0.0);
    }
    let x_mean = (n - 1.0) / 2.0;
    let y_mean = data.iter().sum::<f64>() / n;

    let mut ss_xy = 0.0;
    let mut ss_xx = 0.0;
    let mut ss_yy = 0.0;
    for (i, &y) in data.iter().enumerate() {
        let dx = i as f64 - x_mean;
        let dy = y - y_mean;
        ss_xy += dx * dy;
        ss_xx += dx * dx;
        ss_yy += dy * dy;
    }

    let slope = if ss_xx > 1e-12 { ss_xy / ss_xx } else { 0.0 };
    let intercept = y_mean - slope * x_mean;
    let r_squared = if ss_yy > 1e-12 {
        (ss_xy * ss_xy) / (ss_xx * ss_yy)
    } else {
        0.0
    };

    (slope, intercept, r_squared)
}

fn max_consecutive_increase(data: &[f64]) -> usize {
    if data.len() < 2 {
        return 0;
    }
    let mut max_run = 0;
    let mut current_run = 0;
    for w in data.windows(2) {
        if w[1] > w[0] {
            current_run += 1;
            max_run = max_run.max(current_run);
        } else {
            current_run = 0;
        }
    }
    max_run
}

fn response_time_index(data: &[f64], peak: f64) -> f64 {
    if data.is_empty() || peak.abs() < 1e-12 {
        return 0.0;
    }
    let threshold = peak * 0.5;
    for (i, &v) in data.iter().enumerate() {
        if v >= threshold {
            return i as f64 / data.len() as f64;
        }
    }
    1.0
}

fn rising_vs_falling_ratio(data: &[f64]) -> f64 {
    if data.len() < 2 {
        return 0.5;
    }
    let peak_idx = data
        .iter()
        .enumerate()
        .max_by(|(_, a), (_, b)| a.partial_cmp(b).unwrap_or(std::cmp::Ordering::Equal))
        .map(|(i, _)| i)
        .unwrap_or(0);
    let rising = if peak_idx > 0 { peak_idx as f64 } else { 1.0 };
    let falling = if peak_idx < data.len() - 1 {
        (data.len() - 1 - peak_idx) as f64
    } else {
        1.0
    };
    rising / (rising + falling)
}

#[cfg(test)]
mod tests {
    use super::*;

    fn make_sine_wave(n: usize, freq: f64) -> Vec<f64> {
        (0..n)
            .map(|i| {
                let t = i as f64 / n as f64;
                (2.0 * std::f64::consts::PI * freq * t).sin()
            })
            .collect()
    }

    fn make_pulse(n: usize, center: usize, width: usize, height: f64) -> Vec<f64> {
        (0..n)
            .map(|i| {
                if i >= center.saturating_sub(width / 2) && i <= center + width / 2 {
                    height
                } else {
                    0.0
                }
            })
            .collect()
    }

    #[test]
    fn test_empty_input() {
        let fe = FeatureExtractor::new(256);
        let result = fe.extract(&[]);
        assert!(result.iter().all(|&v| v == 0.0));
    }

    #[test]
    fn test_constant_signal() {
        let fe = FeatureExtractor::new(256);
        let data = vec![5.0; 100];
        let f = fe.extract(&data);
        assert!((f[0] - 5.0).abs() < 1e-10); // peak = 5.0
        assert!((f[1] - 5.0).abs() < 1e-10); // valley = 5.0
        assert!((f[2] - 5.0).abs() < 1e-10); // mean = 5.0
        assert!(f[3] < 1e-10); // std ≈ 0
        assert!(f[12] < 1e-10); // peak-to-peak ≈ 0
    }

    #[test]
    fn test_all_88_dimensions_populated() {
        let fe = FeatureExtractor::new(256);
        let wave = make_sine_wave(512, 5.0);
        let f = fe.extract(&wave);

        // 사인파에서 모든 차원에 의미있는 값이 있어야 함
        let nonzero_count = f.iter().filter(|&&v| v.abs() > 1e-12).count();
        assert!(
            nonzero_count >= 60,
            "사인파에서 최소 60개 차원은 비영 값이어야 합니다 (실제: {})",
            nonzero_count
        );
    }

    #[test]
    fn test_no_nan_or_inf() {
        let fe = FeatureExtractor::new(256);
        let wave = make_sine_wave(256, 3.0);
        let f = fe.extract(&wave);
        for (i, &v) in f.iter().enumerate() {
            assert!(!v.is_nan(), "차원 {}에 NaN 값", i);
            assert!(!v.is_infinite(), "차원 {}에 Inf 값", i);
        }
    }

    #[test]
    fn test_peak_and_valley() {
        let fe = FeatureExtractor::new(256);
        let wave = make_sine_wave(256, 2.0);
        let f = fe.extract(&wave);
        assert!((f[0] - 1.0).abs() < 0.1); // sin peak ≈ 1.0
        assert!((f[1] + 1.0).abs() < 0.1); // sin valley ≈ -1.0
    }

    #[test]
    fn test_sine_mean_near_zero() {
        let fe = FeatureExtractor::new(256);
        let wave = make_sine_wave(1024, 4.0);
        let f = fe.extract(&wave);
        assert!(f[2].abs() < 0.05, "사인파 평균 ≈ 0 (실제: {})", f[2]);
    }

    #[test]
    fn test_pulse_features() {
        let fe = FeatureExtractor::new(256);
        let pulse = make_pulse(200, 100, 20, 10.0);
        let f = fe.extract(&pulse);
        assert!((f[0] - 10.0).abs() < 1e-10); // peak = 10.0
        assert!((f[1]).abs() < 1e-10); // valley = 0.0
        assert!(f[14] > 0.0); // energy > 0
        assert!(f[16] > 0.0); // AUC > 0
    }

    #[test]
    fn test_hjorth_sine() {
        let fe = FeatureExtractor::new(256);
        let wave = make_sine_wave(512, 5.0);
        let f = fe.extract(&wave);
        assert!(f[24] > 0.0, "Hjorth Activity > 0");
        assert!(f[25] > 0.0, "Hjorth Mobility > 0");
        assert!(f[26] > 0.0, "Hjorth Complexity > 0");
    }

    #[test]
    fn test_percentiles_monotonic() {
        let fe = FeatureExtractor::new(256);
        let wave: Vec<f64> = (0..200).map(|i| i as f64).collect();
        let f = fe.extract(&wave);
        // p5 < p10 < p25 < p50 < p75 < p90 < p95
        assert!(f[48] <= f[49]);
        assert!(f[49] <= f[50]);
        assert!(f[50] <= f[51]);
        assert!(f[51] <= f[52]);
        assert!(f[52] <= f[53]);
        assert!(f[53] <= f[54]);
    }

    #[test]
    fn test_autocorrelation_sine() {
        let fe = FeatureExtractor::new(256);
        let wave = make_sine_wave(512, 2.0);
        let f = fe.extract(&wave);
        // 사인파의 lag-1 자기상관은 양수 (높은 상관)
        assert!(f[56] > 0.5, "사인파 lag-1 자기상관 > 0.5 (실제: {})", f[56]);
    }

    #[test]
    fn test_linear_trend() {
        let fe = FeatureExtractor::new(256);
        let wave: Vec<f64> = (0..100).map(|i| 2.0 * i as f64 + 1.0).collect();
        let f = fe.extract(&wave);
        assert!(
            (f[76] - 2.0).abs() < 0.1,
            "선형 추세 기울기 ≈ 2.0 (실제: {})",
            f[76]
        );
        assert!(f[78] > 0.99, "선형 추세 R² ≈ 1.0 (실제: {})", f[78]);
    }

    #[test]
    fn test_snr_noisy_vs_clean() {
        let fe = FeatureExtractor::new(256);
        let clean = make_pulse(200, 100, 20, 10.0);
        let f_clean = fe.extract(&clean);
        // 깨끗한 펄스: baseline이 조용하므로 SNR이 높아야 함
        assert!(
            f_clean[87] > 1.0,
            "깨끗한 펄스 SNR > 1 (실제: {})",
            f_clean[87]
        );
    }

    #[test]
    fn test_single_sample() {
        let fe = FeatureExtractor::new(256);
        let f = fe.extract(&[42.0]);
        assert!((f[0] - 42.0).abs() < 1e-10);
        assert!((f[2] - 42.0).abs() < 1e-10);
        // 나머지는 0.0 또는 합리적 값 (NaN 없음)
        for &v in f.iter() {
            assert!(!v.is_nan());
        }
    }

    #[test]
    fn test_reproducibility() {
        let fe = FeatureExtractor::new(256);
        let wave = make_sine_wave(256, 3.0);
        let f1 = fe.extract(&wave);
        let f2 = fe.extract(&wave);
        for i in 0..88 {
            assert!((f1[i] - f2[i]).abs() < 1e-15, "차원 {} 재현성 실패", i);
        }
    }

    #[test]
    fn test_symmetry_of_symmetric_signal() {
        let fe = FeatureExtractor::new(256);
        let wave: Vec<f64> = (0..200)
            .map(|i| {
                let center = 100.0;
                let x = i as f64 - center;
                (-x * x / 500.0).exp() * 10.0
            })
            .collect();
        let f = fe.extract(&wave);
        // 대칭 신호: symmetry ≈ 0.5
        assert!(
            (f[31] - 0.5).abs() < 0.15,
            "대칭 신호 대칭 지수 ≈ 0.5 (실제: {})",
            f[31]
        );
    }

    #[test]
    fn test_segment_features_exist() {
        let fe = FeatureExtractor::new(256);
        let wave = make_sine_wave(400, 3.0);
        let f = fe.extract(&wave);
        // 4개 세그먼트 모두 mean, std, max, min 값이 있어야 함
        for seg in 0..4 {
            let base = 32 + seg * 4;
            // std > 0 (사인파 세그먼트의 표준편차)
            assert!(f[base + 1] > 0.0, "세그먼트 {} std > 0", seg);
            // max > min
            assert!(f[base + 2] >= f[base + 3], "세그먼트 {} max >= min", seg);
        }
    }
}
