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

  /// 시뮬레이션용 896차원 데이터 생성
  static List<double> generateSimulatedData({int seed = 42}) {
    final rng = Random(seed);
    return List.generate(896, (i) {
      final base = 0.3 + 0.4 * sin(i * 0.02);
      final noise = (rng.nextDouble() - 0.5) * 0.3;
      // 일부 클러스터에 이상치 삽입
      final anomaly = (i > 300 && i < 380) ? 0.4 : 0.0;
      return (base + noise + anomaly).clamp(0.0, 1.0);
    });
  }

  static double _normalize(double value, double min, double max) {
    if (max == min) return 0.5;
    return ((value - min) / (max - min)).clamp(0.0, 1.0);
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
