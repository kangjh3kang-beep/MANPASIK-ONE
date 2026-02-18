/// AI 코치 도메인 레이어 인터페이스
abstract class AiCoachRepository {
  /// 오늘의 건강 인사이트 조회
  Future<HealthInsight> getTodayInsight();

  /// 카테고리별 추천 조회
  Future<List<Recommendation>> getRecommendations(String category);

  /// 건강 트렌드 분석 조회
  Future<TrendAnalysis> getTrendAnalysis(String biomarker, int days);
}

/// 건강 인사이트 모델
class HealthInsight {
  final String summary;
  final String detail;
  final double confidence;
  final DateTime generatedAt;

  const HealthInsight({
    required this.summary,
    required this.detail,
    required this.confidence,
    required this.generatedAt,
  });
}

/// 추천 모델
class Recommendation {
  final String category;
  final String title;
  final String description;
  final int priority;

  const Recommendation({
    required this.category,
    required this.title,
    required this.description,
    required this.priority,
  });
}

/// 트렌드 분석 모델
class TrendAnalysis {
  final String biomarker;
  final String trend; // 'improving', 'stable', 'declining'
  final List<double> values;
  final List<DateTime> dates;

  const TrendAnalysis({
    required this.biomarker,
    required this.trend,
    required this.values,
    required this.dates,
  });
}
