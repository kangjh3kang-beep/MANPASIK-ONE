import 'dart:math';

/// 생체 핑거프린트 분석기 (C2)
///
/// 896차원 스펙트럼 데이터를 12개 바이오마커 클러스터로 축소하고,
/// 이상치 탐지 및 히트맵 데이터 변환을 수행합니다.
class FingerprintAnalyzer {
  /// 896차원 → 12클러스터 차원 축소 (평균 풀링)
  static List<ClusterData> reduceTo12Clusters(List<double> raw896) {
    const clusterNames = [
      '포도당',
      '콜레스테롤',
      '중성지방',
      '요산',
      '크레아티닌',
      'HbA1c',
      'CRP',
      'ALT',
      'AST',
      '빌리루빈',
      '알부민',
      '총단백',
    ];

    final clusters = <ClusterData>[];
    final chunkSize = raw896.length ~/ 12;

    for (var i = 0; i < 12; i++) {
      final start = i * chunkSize;
      final end = (i == 11) ? raw896.length : start + chunkSize;
      final chunk = raw896.sublist(start, end);

      final mean = chunk.reduce((a, b) => a + b) / chunk.length;
      final variance = chunk.map((v) => (v - mean) * (v - mean)).reduce((a, b) => a + b) / chunk.length;
      final stdDev = sqrt(variance);

      // z-score 기반 이상치 비율
      final anomalyRatio = chunk.where((v) => (v - mean).abs() > 2 * stdDev).length / chunk.length;

      clusters.add(ClusterData(
        name: clusterNames[i],
        value: _normalize(mean, 0, 1),
        anomalyScore: anomalyRatio,
        stdDev: stdDev,
      ));
    }

    return clusters;
  }

  /// 896차원 → 32x28 히트맵 그리드 변환
  static List<List<double>> toHeatmapGrid(List<double> raw896) {
    final grid = <List<double>>[];
    var idx = 0;
    for (var row = 0; row < 28; row++) {
      final rowData = <double>[];
      for (var col = 0; col < 32; col++) {
        if (idx < raw896.length) {
          rowData.add(_normalize(raw896[idx], 0, 1));
          idx++;
        } else {
          rowData.add(0.0);
        }
      }
      grid.add(rowData);
    }
    return grid;
  }

  /// 비표적 분석: 주요 이상 바이오마커 탐지 (C3)
  static List<AnomalyResult> detectAnomalies(List<double> raw896) {
    final clusters = reduceTo12Clusters(raw896);
    final anomalies = <AnomalyResult>[];

    for (final cluster in clusters) {
      if (cluster.anomalyScore > 0.15) {
        anomalies.add(AnomalyResult(
          biomarkerName: cluster.name,
          score: cluster.anomalyScore,
          severity: cluster.anomalyScore > 0.4
              ? AnomalySeverity.high
              : cluster.anomalyScore > 0.25
                  ? AnomalySeverity.medium
                  : AnomalySeverity.low,
          description: '${cluster.name} 영역에서 비정상 패턴이 감지되었습니다.',
        ));
      }
    }

    anomalies.sort((a, b) => b.score.compareTo(a.score));
    return anomalies;
  }

  /// 차동 계측 결과로부터 896차원 핑거프린트 벡터를 생성합니다.
  ///
  /// 88차원 기본 특성 추출 + 896차원 확장을 Dart-side에서 수행.
  /// Rust 네이티브 엔진과 동일한 알고리즘 구현.
  static List<double> fromDifferentialData(List<double> diffWave) {
    if (diffWave.isEmpty) return List.filled(896, 0.0);

    // 88차원 기본 특성 추출
    final base88 = _extractFeatures88(diffWave);
    // 896차원 확장
    return _expandTo896(base88);
  }

  /// 시뮬레이션용 896차원 데이터 생성 (테스트/데모 전용)
  ///
  /// 프로덕션에서는 [fromDifferentialData] 사용.
  static List<double> generateSimulatedData({int seed = 42}) {
    final rng = Random(seed);
    return List.generate(896, (i) {
      final base = 0.3 + 0.4 * sin(i * 0.02);
      final noise = (rng.nextDouble() - 0.5) * 0.3;
      final anomaly = (i > 300 && i < 380) ? 0.4 : 0.0;
      return (base + noise + anomaly).clamp(0.0, 1.0);
    });
  }

  static double _normalize(double value, double min, double max) {
    if (max == min) return 0.5;
    return ((value - min) / (max - min)).clamp(0.0, 1.0);
  }

  // ─── 내부: 88차원 특성 추출 (Rust FeatureExtractor Dart 대응) ───

  static List<double> _extractFeatures88(List<double> data) {
    final f = List.filled(88, 0.0);
    if (data.isEmpty) return f;
    final n = data.length;

    double sum = 0, sumSq = 0;
    double peak = double.negativeInfinity, valley = double.infinity;
    for (final v in data) {
      sum += v; sumSq += v * v;
      if (v > peak) peak = v;
      if (v < valley) valley = v;
    }
    final mean = sum / n;
    final variance = (sumSq / n - mean * mean).clamp(0.0, double.infinity);
    final std = sqrt(variance);
    f[0] = peak; f[1] = valley; f[2] = mean; f[3] = std;
    f[4] = sqrt(sumSq / n); // RMS

    final sorted = List<double>.from(data)..sort();
    f[7] = sorted[n ~/ 2]; // median
    f[12] = peak - valley; // peak-to-peak
    f[14] = sumSq; // energy

    // AUC
    double auc = 0;
    for (int i = 1; i < n; i++) { auc += (data[i - 1] + data[i]) * 0.5; }
    f[16] = auc;

    // 세그먼트 통계
    final segSize = n ~/ 4;
    for (int seg = 0; seg < 4; seg++) {
      final start = seg * segSize;
      final end = seg == 3 ? n : start + segSize;
      final chunk = data.sublist(start, end);
      double cS = 0, cMx = double.negativeInfinity, cMn = double.infinity;
      for (final v in chunk) { cS += v; if (v > cMx) cMx = v; if (v < cMn) cMn = v; }
      final cM = cS / chunk.length;
      double cV = 0;
      for (final v in chunk) { cV += (v - cM) * (v - cM); }
      f[32 + seg * 4] = cM;
      f[33 + seg * 4] = sqrt(cV / chunk.length);
      f[34 + seg * 4] = cMx;
      f[35 + seg * 4] = cMn;
    }

    // 백분위수
    f[48] = sorted[(0.05 * (n - 1)).round().clamp(0, n - 1)];
    f[50] = sorted[(0.25 * (n - 1)).round().clamp(0, n - 1)];
    f[52] = sorted[(0.75 * (n - 1)).round().clamp(0, n - 1)];
    f[54] = sorted[(0.95 * (n - 1)).round().clamp(0, n - 1)];

    // 트렌드
    final half = n ~/ 2;
    if (half > 0) {
      f[72] = data.sublist(0, half).reduce((a, b) => a + b) / half;
      f[73] = data.sublist(half).reduce((a, b) => a + b) / (n - half);
    }

    return f;
  }

  // ─── 내부: 88→896 확장 (Rust MultiModeExpander Dart 대응) ───

  static List<double> _expandTo896(List<double> base88) {
    final expanded = <double>[];
    expanded.addAll(base88);                              // [0..88]
    expanded.addAll(base88.map((v) => v * v));             // [88..176]
    expanded.addAll(base88.map((v) =>                      // [176..264]
        v.sign * pow(v.abs(), 1.0 / 3.0).toDouble()));
    expanded.addAll(base88.map((v) =>                      // [264..352]
        v.abs() > 1e-12 ? log(v.abs()) * v.sign : 0.0));

    // [352..616] 교호작용
    int count = 0;
    for (int g = 0; g < 80 && count < 264; g += 8) {
      final ng = g + 8;
      if (ng >= 88) break;
      for (int i = g; i < g + 8 && count < 264; i++) {
        for (int j = ng; j < min(ng + 8, 88) && count < 264; j++) {
          expanded.add(base88[i] * base88[j]);
          count++;
        }
      }
    }
    while (count < 264) { expanded.add(0.0); count++; }

    // [616..744] eNose (네이티브 없을 시 0)
    expanded.addAll(List.filled(128, 0.0));

    // [744..808] 윈도우 통계
    for (int seg = 0; seg < 8; seg++) {
      final start = seg * 11;
      final end = min((seg + 1) * 11, 88);
      final chunk = base88.sublist(start, end);
      double s = 0, mx = double.negativeInfinity, mn = double.infinity;
      for (final v in chunk) { s += v; if (v > mx) mx = v; if (v < mn) mn = v; }
      final m = s / chunk.length;
      double vr = 0;
      for (final v in chunk) { vr += (v - m) * (v - m); }
      expanded.addAll([m, sqrt(vr / chunk.length), mx, mn, mx - mn,
        chunk.map((v) => v * v).reduce((a, b) => a + b), 0.0, 0.0]);
    }

    // [808..896] 정규화
    double gMean = 0, gMin = double.infinity, gMax = double.negativeInfinity;
    for (final v in base88) { gMean += v; if (v < gMin) gMin = v; if (v > gMax) gMax = v; }
    gMean /= 88;
    double gVar = 0;
    for (final v in base88) { gVar += (v - gMean) * (v - gMean); }
    final gStd = sqrt(gVar / 88);
    final gRange = gMax - gMin;
    for (int i = 0; i < 44; i++) {
      expanded.add(gStd > 1e-12 ? (base88[i] - gMean) / gStd : 0.0);
    }
    for (int i = 44; i < 88; i++) {
      expanded.add(gRange > 1e-12 ? (base88[i] - gMin) / gRange : 0.0);
    }

    while (expanded.length < 896) { expanded.add(0.0); }
    if (expanded.length > 896) return expanded.sublist(0, 896);
    return expanded;
  }
}

class ClusterData {
  final String name;
  final double value;
  final double anomalyScore;
  final double stdDev;

  const ClusterData({
    required this.name,
    required this.value,
    required this.anomalyScore,
    required this.stdDev,
  });
}

enum AnomalySeverity { low, medium, high }

class AnomalyResult {
  final String biomarkerName;
  final double score;
  final AnomalySeverity severity;
  final String description;

  const AnomalyResult({
    required this.biomarkerName,
    required this.score,
    required this.severity,
    required this.description,
  });
}
