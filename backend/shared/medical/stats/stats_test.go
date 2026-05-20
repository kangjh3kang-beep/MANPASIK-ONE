package stats_test

import (
	"sort"
	"testing"

	"github.com/manpasik/backend/shared/medical/stats"
)

func TestMean(t *testing.T) {
	cases := []struct {
		values []float64
		want   float64
	}{
		{[]float64{}, 0},
		{[]float64{5}, 5},
		{[]float64{1, 2, 3, 4, 5}, 3},
		{[]float64{-10, 0, 10}, 0},
	}
	for _, c := range cases {
		got := stats.Mean(c.values)
		if got != c.want {
			t.Errorf("Mean(%v) = %f, want %f", c.values, got, c.want)
		}
	}
}

func TestStdDev_Sample(t *testing.T) {
	values := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	mean := stats.Mean(values)
	std := stats.StdDev(values, mean)
	// 표본 표준편차 (n-1 분모) = 2.138...
	if std < 2.13 || std > 2.14 {
		t.Errorf("StdDev = %f, want ~2.138", std)
	}
}

func TestStdDev_BelowTwo(t *testing.T) {
	if stats.StdDev([]float64{}, 0) != 0 {
		t.Error("empty StdDev != 0")
	}
	if stats.StdDev([]float64{5}, 5) != 0 {
		t.Error("single StdDev != 0")
	}
}

func TestPercentile_Boundaries(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5}
	if stats.Percentile(values, 0) != 1 {
		t.Error("P0 != min")
	}
	if stats.Percentile(values, 100) != 5 {
		t.Error("P100 != max")
	}
}

func TestPercentile_Median(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5}
	if stats.Percentile(values, 50) != 3 {
		t.Errorf("P50 = %f, want 3", stats.Percentile(values, 50))
	}
}

func TestPercentile_Interpolation(t *testing.T) {
	values := []float64{1, 100}
	// P50 = 50.5 (선형 보간)
	got := stats.Percentile(values, 50)
	if got < 50 || got > 51 {
		t.Errorf("P50 = %f, want ~50.5", got)
	}
}

func TestPercentile_Empty(t *testing.T) {
	if stats.Percentile([]float64{}, 50) != 0 {
		t.Error("빈 슬라이스 != 0")
	}
}

func TestPercentile_OutOfRange(t *testing.T) {
	values := []float64{1, 2, 3}
	if stats.Percentile(values, -10) != 1 {
		t.Error("음수 백분위 != min")
	}
	if stats.Percentile(values, 200) != 3 {
		t.Error("100+ 백분위 != max")
	}
}

func TestMedianSorted(t *testing.T) {
	values := []float64{5, 1, 3, 2, 4}
	sort.Float64s(values)
	if stats.MedianSorted(values) != 3 {
		t.Error("Median != 3")
	}
}

func TestVariance(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5}
	mean := stats.Mean(values)
	v := stats.Variance(values, mean)
	// 표본 분산 = 2.5
	if v < 2.49 || v > 2.51 {
		t.Errorf("Variance = %f, want ~2.5", v)
	}
}

func TestMinMax(t *testing.T) {
	min, max := stats.MinMax([]float64{3, 1, 4, 1, 5, 9, 2, 6})
	if min != 1 || max != 9 {
		t.Errorf("MinMax = %f, %f, want 1, 9", min, max)
	}
}

func TestMinMax_Empty(t *testing.T) {
	min, max := stats.MinMax([]float64{})
	if min != 0 || max != 0 {
		t.Errorf("empty MinMax = %f, %f", min, max)
	}
}

func TestPercentile_P95P99(t *testing.T) {
	// 1~100 분포
	values := make([]float64, 100)
	for i := 0; i < 100; i++ {
		values[i] = float64(i + 1)
	}

	if p95 := stats.Percentile(values, 95); p95 < 95 || p95 > 96 {
		t.Errorf("P95 = %f, want 95-96", p95)
	}
	if p99 := stats.Percentile(values, 99); p99 < 99 || p99 > 100 {
		t.Errorf("P99 = %f, want 99-100", p99)
	}
}
