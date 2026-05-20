import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/measurement/domain/measurement_repository.dart';

void main() {
  group('StartSessionResult', () {
    test('기본 생성 확인', () {
      const result = StartSessionResult(sessionId: 'sess-1');
      expect(result.sessionId, 'sess-1');
      expect(result.startedAt, isNull);
    });

    test('startedAt 포함 생성', () {
      final now = DateTime.now();
      final result = StartSessionResult(
        sessionId: 'sess-2',
        startedAt: now,
      );
      expect(result.sessionId, 'sess-2');
      expect(result.startedAt, now);
    });
  });

  group('EndSessionResult', () {
    test('기본 생성 확인', () {
      const result = EndSessionResult(
        sessionId: 'sess-1',
        totalMeasurements: 10,
      );
      expect(result.sessionId, 'sess-1');
      expect(result.totalMeasurements, 10);
      expect(result.endedAt, isNull);
    });

    test('endedAt 포함 생성', () {
      final now = DateTime.now();
      final result = EndSessionResult(
        sessionId: 'sess-3',
        totalMeasurements: 25,
        endedAt: now,
      );
      expect(result.endedAt, now);
      expect(result.totalMeasurements, 25);
    });
  });

  group('MeasurementHistoryResult', () {
    test('빈 목록 생성', () {
      const result = MeasurementHistoryResult(
        items: [],
        totalCount: 0,
      );
      expect(result.items, isEmpty);
      expect(result.totalCount, 0);
    });

    test('항목이 있는 결과 생성', () {
      final items = [
        MeasurementHistoryItem(
          sessionId: 'sess-1',
          cartridgeType: 'BIO-GLU',
          primaryValue: 95.0,
          unit: 'mg/dL',
          measuredAt: DateTime(2026, 4, 19),
        ),
        const MeasurementHistoryItem(
          sessionId: 'sess-2',
          cartridgeType: 'BIO-CHO',
          primaryValue: 180.0,
          unit: 'mg/dL',
        ),
      ];
      final result = MeasurementHistoryResult(
        items: items,
        totalCount: 50,
      );
      expect(result.items, hasLength(2));
      expect(result.totalCount, 50);
    });
  });

  group('MeasurementHistoryItem', () {
    test('기본 생성 확인', () {
      const item = MeasurementHistoryItem(
        sessionId: 'sess-1',
        cartridgeType: 'BIO-GLU',
        primaryValue: 95.5,
        unit: 'mg/dL',
      );
      expect(item.sessionId, 'sess-1');
      expect(item.cartridgeType, 'BIO-GLU');
      expect(item.primaryValue, 95.5);
      expect(item.unit, 'mg/dL');
      expect(item.measuredAt, isNull);
    });

    test('measuredAt 포함 생성', () {
      final now = DateTime.now();
      final item = MeasurementHistoryItem(
        sessionId: 'sess-5',
        cartridgeType: 'ENV-CO2',
        primaryValue: 450.0,
        unit: 'ppm',
        measuredAt: now,
      );
      expect(item.measuredAt, now);
    });

    test('다양한 cartridgeType 값', () {
      for (final type in ['BIO-GLU', 'BIO-CHO', 'ENV-CO2', 'GAS-VOC']) {
        final item = MeasurementHistoryItem(
          sessionId: 'test',
          cartridgeType: type,
          primaryValue: 100.0,
          unit: 'unit',
        );
        expect(item.cartridgeType, type);
      }
    });
  });
}
