import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/data_hub/domain/data_hub_repository.dart';

void main() {
  group('ExportFormat enum', () {
    test('값 목록 확인 (3개)', () {
      expect(ExportFormat.values, hasLength(3));
      expect(ExportFormat.values, contains(ExportFormat.csv));
      expect(ExportFormat.values, contains(ExportFormat.pdf));
      expect(ExportFormat.values, contains(ExportFormat.json));
    });
  });

  group('TrendDataPoint', () {
    test('생성 확인', () {
      final point = TrendDataPoint(
        timestamp: DateTime(2026, 4, 19, 10, 0),
        value: 95.0,
        unit: 'mg/dL',
        biomarkerType: 'blood_glucose',
        isWithinRange: true,
      );
      expect(point.value, 95.0);
      expect(point.unit, 'mg/dL');
      expect(point.biomarkerType, 'blood_glucose');
      expect(point.isWithinRange, isTrue);
      expect(point.evidenceStatus, 'unknown');
      expect(point.diagnosticReady, isFalse);
      expect(point.evidenceGaps, isEmpty);
    });

    test('evidence metadata 포함 생성', () {
      final point = TrendDataPoint(
        timestamp: DateTime(2026, 4, 19, 10, 0),
        value: 95.0,
        unit: 'mg/dL',
        biomarkerType: 'blood_glucose',
        isWithinRange: true,
        evidenceStatus: 'research_only',
        diagnosticReady: false,
        evidenceGaps: const ['clinical_lock_required'],
      );
      expect(point.evidenceStatus, 'research_only');
      expect(point.diagnosticReady, isFalse);
      expect(point.evidenceGaps, contains('clinical_lock_required'));
    });

    test('범위 초과 데이터 포인트', () {
      final point = TrendDataPoint(
        timestamp: DateTime(2026, 4, 19),
        value: 250.0,
        unit: 'mg/dL',
        biomarkerType: 'blood_glucose',
        isWithinRange: false,
      );
      expect(point.isWithinRange, isFalse);
      expect(point.value, 250.0);
    });
  });

  group('BiomarkerSummary', () {
    test('기본 생성 확인', () {
      const summary = BiomarkerSummary(
        biomarkerType: 'blood_glucose',
        displayName: '혈당',
        unit: 'mg/dL',
        referenceMin: 70.0,
        referenceMax: 110.0,
        totalMeasurements: 30,
        trend: 'stable',
      );
      expect(summary.biomarkerType, 'blood_glucose');
      expect(summary.displayName, '혈당');
      expect(summary.latestValue, isNull);
      expect(summary.averageValue, isNull);
      expect(summary.totalMeasurements, 30);
      expect(summary.trend, 'stable');
      expect(summary.latestEvidenceStatus, 'unknown');
      expect(summary.latestDiagnosticReady, isFalse);
      expect(summary.latestEvidenceGaps, isEmpty);
    });

    test('모든 통계 값 포함 생성', () {
      const summary = BiomarkerSummary(
        biomarkerType: 'cholesterol',
        displayName: '콜레스테롤',
        unit: 'mg/dL',
        latestValue: 195.0,
        averageValue: 200.0,
        minValue: 180.0,
        maxValue: 220.0,
        referenceMin: 0.0,
        referenceMax: 200.0,
        totalMeasurements: 10,
        trend: 'falling',
        latestEvidenceStatus: 'research_only',
        latestDiagnosticReady: false,
        latestEvidenceGaps: ['clinical_lock_required'],
      );
      expect(summary.latestValue, 195.0);
      expect(summary.averageValue, 200.0);
      expect(summary.minValue, 180.0);
      expect(summary.maxValue, 220.0);
      expect(summary.latestEvidenceStatus, 'research_only');
      expect(summary.latestDiagnosticReady, isFalse);
      expect(summary.latestEvidenceGaps, contains('clinical_lock_required'));
    });

    test('다양한 trend 값', () {
      for (final trend in ['rising', 'falling', 'stable', 'insufficient']) {
        final summary = BiomarkerSummary(
          biomarkerType: 'test',
          displayName: '테스트',
          unit: 'unit',
          referenceMin: 0,
          referenceMax: 100,
          totalMeasurements: 0,
          trend: trend,
        );
        expect(summary.trend, trend);
      }
    });
  });

  group('ExportResult', () {
    test('생성 확인', () {
      final now = DateTime.now();
      final result = ExportResult(
        filePath: '/data/export_2026-04-19.csv',
        format: ExportFormat.csv,
        recordCount: 100,
        exportedAt: now,
      );
      expect(result.filePath, contains('.csv'));
      expect(result.format, ExportFormat.csv);
      expect(result.recordCount, 100);
      expect(result.exportedAt, now);
    });

    test('각 형식별 ExportResult 생성', () {
      for (final format in ExportFormat.values) {
        final result = ExportResult(
          filePath: '/data/export.${format.name}',
          format: format,
          recordCount: 50,
          exportedAt: DateTime.now(),
        );
        expect(result.format, format);
      }
    });
  });
}
