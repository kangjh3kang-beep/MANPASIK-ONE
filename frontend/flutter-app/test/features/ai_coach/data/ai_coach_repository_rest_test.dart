import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/ai_coach/data/ai_coach_repository_rest.dart';
import 'package:manpasik/features/ai_coach/domain/ai_coach_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';

void main() {
  group('AiCoachRepositoryRest', () {
    test('AiCoachRepositoryRest는 AiCoachRepository를 구현한다', () {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = AiCoachRepositoryRest(client, userId: 'user-1');
      expect(repo, isA<AiCoachRepository>());
    });

    test('getTodayInsight는 DioException 시 기본 인사이트를 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = AiCoachRepositoryRest(client, userId: 'user-1');
      final insight = await repo.getTodayInsight();
      expect(insight, isA<HealthInsight>());
      expect(insight.summary, isNotEmpty);
      expect(insight.confidence, 0.0);
    });

    test('getRecommendations는 DioException 시 빈 리스트를 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = AiCoachRepositoryRest(client, userId: 'user-1');
      final recs = await repo.getRecommendations('exercise');
      expect(recs, isEmpty);
    });
  });

  group('AI Coach 도메인 모델', () {
    test('HealthInsight 생성 확인', () {
      final insight = HealthInsight(
        summary: '오늘 건강 상태 양호',
        detail: '혈당 정상, 혈압 정상',
        confidence: 0.85,
        generatedAt: DateTime(2026, 2, 19),
      );
      expect(insight.summary, contains('양호'));
      expect(insight.confidence, 0.85);
    });

    test('Recommendation 생성 확인', () {
      const rec = Recommendation(
        category: 'exercise',
        title: '걷기 운동 추천',
        description: '하루 30분 걷기',
        priority: 1,
      );
      expect(rec.category, 'exercise');
      expect(rec.priority, 1);
    });
  });
}
