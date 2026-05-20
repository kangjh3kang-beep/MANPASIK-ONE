// fingerprint/multi_mode_expander.rs
//
// 88차원 기본 핑거프린트를 896차원으로 확장하는 융합 모듈.
// SSOT: 88→448→896→1792차원 경로 중 88→896 확장 담당.
//
// 896차원 구성:
//   [0..88]     원본 특성 복사
//   [88..176]   제곱항 (polynomial degree 2)
//   [176..264]  세제곱 루트항 (비선형 변환)
//   [264..352]  로그 변환항
//   [352..616]  상위 교호작용 (feature_i × feature_j, 상위 264쌍)
//   [616..744]  eNose 8채널 × 16 파생 (교차 반응 패턴)
//   [744..808]  윈도우 통계 (88특성의 이동 통계량)
//   [808..896]  정규화 특성 (z-score, min-max, rank 기반)

/// 88차원 → 896차원 특성 공간 확장기
pub struct Expander;

impl Expander {
    /// 88차원 기본 특성 + eNose/eTongue 데이터를 융합하여 896차원 결합 벡터를 생성합니다.
    ///
    /// # Arguments
    /// * `base_features` - FeatureExtractor가 추출한 88차원 특성 벡터
    /// * `enose_data` - 전자코 8채널 교차 반응 데이터 (비어있으면 0으로 처리)
    ///
    /// # Returns
    /// 896차원 확장 벡터
    pub fn expand_to_896(base_features: &[f64; 88], enose_data: &[f64]) -> Vec<f64> {
        let mut expanded = Vec::with_capacity(896);

        // === Block 1: 원본 특성 복사 [0..88] ===
        expanded.extend_from_slice(base_features);

        // === Block 2: 제곱항 [88..176] ===
        for &v in base_features.iter() {
            expanded.push(v * v);
        }

        // === Block 3: 세제곱 루트 비선형 변환 [176..264] ===
        for &v in base_features.iter() {
            expanded.push(v.signum() * v.abs().cbrt());
        }

        // === Block 4: 로그 변환 [264..352] ===
        for &v in base_features.iter() {
            expanded.push(safe_log(v));
        }

        // === Block 5: 상위 교호작용 쌍 [352..616] = 264개 ===
        // 88차원에서 C(88,2) = 3828 쌍 중, 인접 그룹 교호작용 상위 264쌍 선택
        // 전략: 그룹 간 교차 (Group0×Group1, Group1×Group2, ...) + 그룹 내 인접
        let mut interaction_count = 0;
        const MAX_INTERACTIONS: usize = 264;

        // 그룹 간 교호작용: 인접 그룹의 특성을 교차 곱
        'outer: for group_offset in (0..80).step_by(8) {
            let next_group = group_offset + 8;
            if next_group >= 88 {
                break;
            }
            for i in group_offset..group_offset + 8 {
                for j in next_group..(next_group + 8).min(88) {
                    if interaction_count >= MAX_INTERACTIONS {
                        break 'outer;
                    }
                    expanded.push(base_features[i] * base_features[j]);
                    interaction_count += 1;
                }
            }
        }
        // 부족분 채움: 비인접 그룹 교호작용
        if interaction_count < MAX_INTERACTIONS {
            'fill: for i in 0..88 {
                for j in (i + 16)..88 {
                    if interaction_count >= MAX_INTERACTIONS {
                        break 'fill;
                    }
                    expanded.push(base_features[i] * base_features[j]);
                    interaction_count += 1;
                }
            }
        }

        // === Block 6: eNose 교차 반응 패턴 [616..744] = 128개 ===
        // 8채널 × 16 파생 특성
        let enose_8: [f64; 8] = {
            let mut arr = [0.0f64; 8];
            for (i, slot) in arr.iter_mut().enumerate() {
                *slot = enose_data.get(i).copied().unwrap_or(0.0);
            }
            arr
        };
        // 16 파생 per channel: {raw, squared, log, cbrt, ×base[0..3], ×base[4..7], ratio_to_mean, diff_to_prev, normalized}
        let enose_mean = enose_8.iter().sum::<f64>() / 8.0;
        for (ch, &e_val) in enose_8.iter().enumerate() {
            expanded.push(e_val); // raw
            expanded.push(e_val * e_val); // squared
            expanded.push(safe_log(e_val)); // log
            expanded.push(e_val.signum() * e_val.abs().cbrt()); // cbrt
                                                                // base feature 교차
            for k in 0..4 {
                let base_idx = ch.min(87);
                expanded.push(e_val * base_features[(base_idx + k * 22) % 88]);
            }
            // ratio to mean
            expanded.push(if enose_mean.abs() > 1e-12 {
                e_val / enose_mean
            } else {
                0.0
            });
            // diff to prev channel
            let prev = if ch > 0 { enose_8[ch - 1] } else { enose_8[7] };
            expanded.push(e_val - prev);
            // z-score normalization
            let enose_std = {
                let var: f64 = enose_8
                    .iter()
                    .map(|&x| (x - enose_mean).powi(2))
                    .sum::<f64>()
                    / 8.0;
                var.sqrt()
            };
            expanded.push(if enose_std > 1e-12 {
                (e_val - enose_mean) / enose_std
            } else {
                0.0
            });
            // energy contribution
            let total_energy: f64 = enose_8.iter().map(|x| x * x).sum();
            expanded.push(if total_energy > 1e-12 {
                (e_val * e_val) / total_energy
            } else {
                0.125
            });

            // padding to 16 per channel
            let emitted = 4 + 4 + 4; // raw/sq/log/cbrt + 4 cross + ratio/diff/zscore/energy
            let remaining = 16 - emitted;
            for r in 0..remaining {
                // 추가 파생: 지수 가중 평균
                let decay = 0.5f64.powi(r as i32 + 1);
                expanded.push(e_val * decay);
            }
        }

        // === Block 7: 윈도우 통계량 [744..808] = 64개 ===
        // 88차원을 8개 세그먼트로 나누어 각 세그먼트의 통계 (mean, std, max, min, range, energy, skew, kurtosis)
        for seg in 0..8 {
            let start = seg * 11;
            let end = ((seg + 1) * 11).min(88);
            let segment = &base_features[start..end];
            let (mean, var) = segment_mean_var(segment);
            let std = var.sqrt();
            let mx = segment.iter().cloned().fold(f64::NEG_INFINITY, f64::max);
            let mn = segment.iter().cloned().fold(f64::INFINITY, f64::min);
            expanded.push(mean);
            expanded.push(std);
            expanded.push(mx);
            expanded.push(mn);
            expanded.push(mx - mn); // range
            expanded.push(segment.iter().map(|x| x * x).sum::<f64>()); // energy
                                                                       // skewness
            if std > 1e-12 && segment.len() >= 3 {
                let m3 = segment
                    .iter()
                    .map(|x| ((x - mean) / std).powi(3))
                    .sum::<f64>()
                    / segment.len() as f64;
                expanded.push(m3);
            } else {
                expanded.push(0.0);
            }
            // kurtosis
            if std > 1e-12 && segment.len() >= 4 {
                let m4 = segment
                    .iter()
                    .map(|x| ((x - mean) / std).powi(4))
                    .sum::<f64>()
                    / segment.len() as f64
                    - 3.0;
                expanded.push(m4);
            } else {
                expanded.push(0.0);
            }
        }

        // === Block 8: 정규화 특성 [808..896] = 88개 ===
        // z-score 정규화
        let (global_mean, global_var) = segment_mean_var(base_features);
        let global_std = global_var.sqrt();
        // min-max 범위
        let global_min = base_features.iter().cloned().fold(f64::INFINITY, f64::min);
        let global_max = base_features
            .iter()
            .cloned()
            .fold(f64::NEG_INFINITY, f64::max);
        let global_range = global_max - global_min;

        // 88차원 중 교차 정규화: z-score(44) + min-max(44)
        for i in 0..44 {
            expanded.push(if global_std > 1e-12 {
                (base_features[i] - global_mean) / global_std
            } else {
                0.0
            });
        }
        for i in 44..88 {
            expanded.push(if global_range > 1e-12 {
                (base_features[i] - global_min) / global_range
            } else {
                0.0
            });
        }

        // 최종 크기 보정: 정확히 896차원
        expanded.truncate(896);
        while expanded.len() < 896 {
            expanded.push(0.0);
        }

        // NaN/Inf 보호
        for v in expanded.iter_mut() {
            if v.is_nan() || v.is_infinite() {
                *v = 0.0;
            }
        }

        expanded
    }
}

fn safe_log(v: f64) -> f64 {
    if v.abs() < 1e-12 {
        0.0
    } else {
        v.abs().ln().copysign(v)
    }
}

fn segment_mean_var(data: &[f64]) -> (f64, f64) {
    if data.is_empty() {
        return (0.0, 0.0);
    }
    let n = data.len() as f64;
    let mean = data.iter().sum::<f64>() / n;
    let var = data.iter().map(|x| (x - mean).powi(2)).sum::<f64>() / n;
    (mean, var)
}

#[cfg(test)]
mod tests {
    use super::*;

    fn make_base_features() -> [f64; 88] {
        let mut f = [0.0f64; 88];
        for i in 0..88 {
            f[i] = (i as f64 + 1.0) * 0.1;
        }
        f
    }

    #[test]
    fn test_output_length() {
        let base = make_base_features();
        let enose = vec![1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0];
        let result = Expander::expand_to_896(&base, &enose);
        assert_eq!(result.len(), 896, "출력 벡터는 정확히 896차원이어야 합니다");
    }

    #[test]
    fn test_first_88_preserved() {
        let base = make_base_features();
        let result = Expander::expand_to_896(&base, &[]);
        for i in 0..88 {
            assert!(
                (result[i] - base[i]).abs() < 1e-10,
                "원본 특성 [{}] 보존 실패: {} vs {}",
                i,
                result[i],
                base[i]
            );
        }
    }

    #[test]
    fn test_squared_terms() {
        let base = make_base_features();
        let result = Expander::expand_to_896(&base, &[]);
        for i in 0..88 {
            let expected = base[i] * base[i];
            assert!(
                (result[88 + i] - expected).abs() < 1e-10,
                "제곱항 [{}] 오류: {} vs {}",
                i,
                result[88 + i],
                expected
            );
        }
    }

    #[test]
    fn test_no_nan_or_inf() {
        let base = make_base_features();
        let enose = vec![0.0, 1.0, -1.0, 0.5, -0.5, 100.0, -100.0, 0.001];
        let result = Expander::expand_to_896(&base, &enose);
        for (i, &v) in result.iter().enumerate() {
            assert!(!v.is_nan(), "차원 {}에 NaN", i);
            assert!(!v.is_infinite(), "차원 {}에 Inf", i);
        }
    }

    #[test]
    fn test_empty_enose() {
        let base = make_base_features();
        let result = Expander::expand_to_896(&base, &[]);
        assert_eq!(result.len(), 896);
        // eNose 블록이 0으로 채워져야 함
        // Block 6 starts at 616
        let enose_block = &result[616..744];
        // 빈 eNose → raw 값들이 0.0
        for (i, &v) in enose_block.iter().enumerate().step_by(16) {
            assert!((v - 0.0).abs() < 1e-10, "빈 eNose raw 값 [{}]: {}", i, v);
        }
    }

    #[test]
    fn test_zero_features() {
        let base = [0.0f64; 88];
        let result = Expander::expand_to_896(&base, &[]);
        assert_eq!(result.len(), 896);
        // 모든 값이 0.0 또는 유효한 숫자여야 함
        for (i, &v) in result.iter().enumerate() {
            assert!(!v.is_nan(), "제로 입력에서 NaN at {}", i);
        }
    }

    #[test]
    fn test_interactions_nonzero() {
        let base = make_base_features();
        let result = Expander::expand_to_896(&base, &[]);
        // 교호작용 블록 [352..616] 에 비영 값이 있어야 함
        let interaction_block = &result[352..616];
        let nonzero = interaction_block
            .iter()
            .filter(|&&v| v.abs() > 1e-12)
            .count();
        assert!(
            nonzero >= 200,
            "교호작용 블록에 최소 200개 비영 값 (실제: {})",
            nonzero
        );
    }

    #[test]
    fn test_reproducibility() {
        let base = make_base_features();
        let enose = vec![1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0];
        let r1 = Expander::expand_to_896(&base, &enose);
        let r2 = Expander::expand_to_896(&base, &enose);
        for i in 0..896 {
            assert!((r1[i] - r2[i]).abs() < 1e-15, "재현성 실패 at dim {}", i);
        }
    }

    #[test]
    fn test_normalized_block_range() {
        let base = make_base_features();
        let result = Expander::expand_to_896(&base, &[]);
        // 정규화 블록 [808..896] - min-max 부분은 [0,1] 범위
        for i in 852..896 {
            assert!(
                result[i] >= -0.01 && result[i] <= 1.01,
                "min-max 정규화 범위 초과 at {}: {}",
                i,
                result[i]
            );
        }
    }
}
