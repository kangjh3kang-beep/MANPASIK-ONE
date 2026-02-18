import 'package:dio/dio.dart';
import 'package:manpasik/features/ai_coach/domain/ai_coach_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';

/// REST Gateway를 사용하는 AiCoachRepository 구현체
class AiCoachRepositoryRest implements AiCoachRepository {
  AiCoachRepositoryRest(this._client, {required this.userId});

  final ManPaSikRestClient _client;
  final String userId;

  @override
  Future<HealthInsight> getTodayInsight() async {
    try {
      final res = await _client.getHealthScore(userId);
      return HealthInsight(
        summary: res['summary'] as String? ?? '오늘의 건강 상태를 분석 중입니다.',
        detail: res['detail'] as String? ?? '',
        confidence: (res['confidence'] as num?)?.toDouble() ??
            (res['score'] as num?)?.toDouble() ??
            0.75,
        generatedAt: res['generated_at'] != null
            ? DateTime.tryParse(res['generated_at'] as String) ?? DateTime.now()
            : DateTime.now(),
      );
    } on DioException {
      return HealthInsight(
        summary: '건강 데이터를 불러올 수 없습니다.',
        detail: '서버 연결을 확인해주세요.',
        confidence: 0.0,
        generatedAt: DateTime.now(),
      );
    }
  }

  @override
  Future<List<Recommendation>> getRecommendations(String category) async {
    try {
      final res = await _client.getRecommendations(userId);
      final list = res['recommendations'] as List<dynamic>? ?? [];
      return list.map((r) {
        final m = r as Map<String, dynamic>;
        return Recommendation(
          category: m['category'] as String? ?? category,
          title: m['title'] as String? ?? '',
          description: m['description'] as String? ?? m['message'] as String? ?? '',
          priority: m['priority'] as int? ?? 0,
        );
      }).toList();
    } on DioException {
      return [];
    }
  }

  @override
  Future<TrendAnalysis> getTrendAnalysis(String biomarker, int days) async {
    try {
      final res = await _client.predictTrend(
        userId: userId,
        metricName: biomarker,
        historyDays: days,
        predictionDays: 7,
      );
      final points = res['data_points'] as List<dynamic>? ?? [];
      return TrendAnalysis(
        biomarker: biomarker,
        trend: res['trend'] as String? ?? 'stable',
        values: points
            .map((p) => (p is Map ? (p['value'] as num?)?.toDouble() : (p as num?)?.toDouble()) ?? 0.0)
            .toList(),
        dates: points.map((p) {
          if (p is Map && p['date'] != null) {
            return DateTime.tryParse(p['date'] as String) ?? DateTime.now();
          }
          return DateTime.now();
        }).toList(),
      );
    } on DioException {
      return TrendAnalysis(
        biomarker: biomarker,
        trend: 'insufficient',
        values: [],
        dates: [],
      );
    }
  }
}
