// Package stats는 medical 도메인 공용 통계 헬퍼입니다.
//
// poise/sla/observability 등 여러 모듈에서 중복 정의되던 mean/percentile/stddev를
// 단일 검증 지점으로 통합합니다 (DRY + 단일 회귀 테스트).
package stats

import "math"

// Mean은 슬라이스의 산술 평균을 반환합니다. 빈 슬라이스는 0.
func Mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// StdDev는 표본 표준편차(분모 n-1)를 반환합니다.
//
// 길이가 2 미만이면 0.
func StdDev(values []float64, mean float64) float64 {
	if len(values) < 2 {
		return 0
	}
	sumSq := 0.0
	for _, v := range values {
		sumSq += (v - mean) * (v - mean)
	}
	return math.Sqrt(sumSq / float64(len(values)-1))
}

// Percentile은 정렬된 값 슬라이스의 백분위수를 선형 보간으로 계산합니다.
//
// 호출자가 사전에 sort.Float64s로 정렬해야 합니다.
// p는 0~100 범위, 범위 밖은 자동 클램프.
func Percentile(sortedValues []float64, p float64) float64 {
	if len(sortedValues) == 0 {
		return 0
	}
	if p <= 0 {
		return sortedValues[0]
	}
	if p >= 100 {
		return sortedValues[len(sortedValues)-1]
	}
	idx := (p / 100.0) * float64(len(sortedValues)-1)
	lower := int(idx)
	upper := lower + 1
	if upper >= len(sortedValues) {
		return sortedValues[lower]
	}
	weight := idx - float64(lower)
	return sortedValues[lower]*(1-weight) + sortedValues[upper]*weight
}

// MedianSorted는 정렬된 값의 중앙값(P50)을 반환합니다.
func MedianSorted(sortedValues []float64) float64 {
	return Percentile(sortedValues, 50)
}

// Variance는 표본 분산(분모 n-1)을 반환합니다.
func Variance(values []float64, mean float64) float64 {
	if len(values) < 2 {
		return 0
	}
	sumSq := 0.0
	for _, v := range values {
		sumSq += (v - mean) * (v - mean)
	}
	return sumSq / float64(len(values)-1)
}

// MinMax는 슬라이스의 최솟값과 최댓값을 반환합니다.
func MinMax(values []float64) (min, max float64) {
	if len(values) == 0 {
		return 0, 0
	}
	min = values[0]
	max = values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max
}
