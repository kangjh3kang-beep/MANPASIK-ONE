import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/ai_coach/domain/ai_coach_repository.dart';

void main() {
  group('HealthInsight', () {
    test('생성 확인', () {
      final now = DateTime.now();
      final insight = HealthInsight(
        summary: '혈당이 안정적입니다',
        detail: '지난 7일간 평균 혈당 95mg/dL로 정상 범위를 유지하고 있습니다.',
        confidence: 0.92,
        generatedAt: now,
      );
      expect(insight.summary, '혈당이 안정적입니다');
      expect(insight.confidence, 0.92);
      expect(insight.generatedAt, now);
    });

    test('confidence 범위 확인', () {
      final insight = HealthInsight(
        summary: '낮은 신뢰도',
        detail: '데이터 부족',
        confidence: 0.3,
        generatedAt: DateTime.now(),
      );
      expect(insight.confidence, lessThan(1.0));
      expect(insight.confidence, greaterThan(0.0));
    });
  });

  group('Recommendation', () {
    test('생성 확인', () {
      const rec = Recommendation(
        category: 'exercise',
        title: '아침 산책 추천',
        description: '식후 30분 걷기가 혈당 조절에 효과적입니다.',
        priority: 1,
      );
      expect(rec.category, 'exercise');
      expect(rec.title, '아침 산책 추천');
      expect(rec.priority, 1);
    });

    test('다양한 카테고리', () {
      for (final cat in ['exercise', 'diet', 'sleep', 'medication']) {
        final rec = Recommendation(
          category: cat,
          title: '테스트',
          description: '설명',
          priority: 1,
        );
        expect(rec.category, cat);
      }
    });

    test('우선순위 비교', () {
      const high = Recommendation(
        category: 'diet',
        title: '고위험',
        description: '즉시 조치 필요',
        priority: 1,
      );
      const low = Recommendation(
        category: 'exercise',
        title: '저위험',
        description: '여유있게 시작',
        priority: 5,
      );
      expect(high.priority, lessThan(low.priority));
    });
  });

  group('TrendAnalysis', () {
    test('생성 확인', () {
      const analysis = TrendAnalysis(
        biomarker: 'blood_glucose',
        trend: 'improving',
        values: [110.0, 105.0, 98.0, 95.0],
        dates: [],
      );
      expect(analysis.biomarker, 'blood_glucose');
      expect(analysis.trend, 'improving');
      expect(analysis.values, hasLength(4));
    });

    test('다양한 trend 값', () {
      for (final trend in ['improving', 'stable', 'declining']) {
        final analysis = TrendAnalysis(
          biomarker: 'cholesterol',
          trend: trend,
          values: [200.0, 210.0],
          dates: [DateTime(2026, 4, 1), DateTime(2026, 4, 15)],
        );
        expect(analysis.trend, trend);
      }
    });

    test('values와 dates 길이 일치', () {
      final dates = [DateTime(2026, 4, 1), DateTime(2026, 4, 8)];
      const values = [95.0, 92.0];
      final analysis = TrendAnalysis(
        biomarker: 'blood_glucose',
        trend: 'improving',
        values: values,
        dates: dates,
      );
      expect(analysis.values.length, analysis.dates.length);
    });
  });
}
