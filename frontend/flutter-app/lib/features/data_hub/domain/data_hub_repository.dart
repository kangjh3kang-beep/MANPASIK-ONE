/// 데이터 허브 도메인 모델 및 리포지토리
///
/// 측정 데이터 통합 관리, 트렌드 분석, 데이터 내보내기

/// 바이오마커 트렌드 데이터 포인트
class TrendDataPoint {
  final DateTime timestamp;
  final double value;
  final String unit;
  final String biomarkerType;
  final bool isWithinRange;

  const TrendDataPoint({
    required this.timestamp,
    required this.value,
    required this.unit,
    required this.biomarkerType,
    required this.isWithinRange,
  });
}

/// 바이오마커 요약 통계
class BiomarkerSummary {
  final String biomarkerType;
  final String displayName;
  final String unit;
  final double? latestValue;
  final double? averageValue;
  final double? minValue;
  final double? maxValue;
  final double referenceMin;
  final double referenceMax;
  final int totalMeasurements;
  final String trend; // 'rising', 'falling', 'stable', 'insufficient'

  const BiomarkerSummary({
    required this.biomarkerType,
    required this.displayName,
    required this.unit,
    this.latestValue,
    this.averageValue,
    this.minValue,
    this.maxValue,
    required this.referenceMin,
    required this.referenceMax,
    required this.totalMeasurements,
    required this.trend,
  });
}

/// 내보내기 형식
enum ExportFormat { csv, pdf, json }

/// 내보내기 결과
class ExportResult {
  final String filePath;
  final ExportFormat format;
  final int recordCount;
  final DateTime exportedAt;

  const ExportResult({
    required this.filePath,
    required this.format,
    required this.recordCount,
    required this.exportedAt,
  });
}

/// 데이터 허브 리포지토리 인터페이스
abstract class DataHubRepository {
  /// 기간별 트렌드 데이터 조회
  Future<List<TrendDataPoint>> getTrendData({
    required String biomarkerType,
    required DateTime from,
    required DateTime to,
  });

  /// 바이오마커별 요약 통계
  Future<BiomarkerSummary> getBiomarkerSummary(String biomarkerType);

  /// 전체 바이오마커 요약 목록
  Future<List<BiomarkerSummary>> getAllBiomarkerSummaries();

  /// 데이터 내보내기
  Future<ExportResult> exportData({
    required ExportFormat format,
    DateTime? from,
    DateTime? to,
    List<String>? biomarkerTypes,
  });

  /// 총 측정 횟수
  Future<int> getTotalMeasurementCount();
}
