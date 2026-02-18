import 'package:dio/dio.dart';
import 'package:manpasik/features/data_hub/domain/data_hub_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';

/// REST Gateway를 사용하는 DataHubRepository 구현체
class DataHubRepositoryRest implements DataHubRepository {
  DataHubRepositoryRest(this._client, {required this.userId});

  final ManPaSikRestClient _client;
  final String userId;

  @override
  Future<List<TrendDataPoint>> getTrendData({
    required String biomarkerType,
    required DateTime from,
    required DateTime to,
  }) async {
    try {
      final res = await _client.getMeasurementHistory(userId, limit: 100);
      final measurements = res['measurements'] as List<dynamic>? ?? [];
      return measurements
          .map((m) {
            final map = m as Map<String, dynamic>;
            final dt = map['measured_at'] != null
                ? DateTime.tryParse(map['measured_at'] as String)
                : null;
            if (dt == null || dt.isBefore(from) || dt.isAfter(to)) return null;
            final type = map['cartridge_type'] as String? ?? '';
            if (biomarkerType.isNotEmpty && type != biomarkerType) return null;
            return TrendDataPoint(
              timestamp: dt,
              value: (map['primary_value'] as num?)?.toDouble() ?? 0.0,
              unit: map['unit'] as String? ?? '',
              biomarkerType: type,
              isWithinRange: map['is_within_range'] as bool? ?? true,
            );
          })
          .whereType<TrendDataPoint>()
          .toList();
    } on DioException {
      return [];
    }
  }

  @override
  Future<BiomarkerSummary> getBiomarkerSummary(String biomarkerType) async {
    try {
      final points = await getTrendData(
        biomarkerType: biomarkerType,
        from: DateTime.now().subtract(const Duration(days: 90)),
        to: DateTime.now(),
      );
      if (points.isEmpty) {
        return BiomarkerSummary(
          biomarkerType: biomarkerType,
          displayName: biomarkerType,
          unit: '',
          referenceMin: 0,
          referenceMax: 100,
          totalMeasurements: 0,
          trend: 'insufficient',
        );
      }
      final values = points.map((p) => p.value).toList();
      values.sort();
      final avg = values.reduce((a, b) => a + b) / values.length;
      return BiomarkerSummary(
        biomarkerType: biomarkerType,
        displayName: biomarkerType,
        unit: points.first.unit,
        latestValue: points.last.value,
        averageValue: avg,
        minValue: values.first,
        maxValue: values.last,
        referenceMin: 0,
        referenceMax: 200,
        totalMeasurements: points.length,
        trend: _computeTrend(values),
      );
    } on DioException {
      return BiomarkerSummary(
        biomarkerType: biomarkerType,
        displayName: biomarkerType,
        unit: '',
        referenceMin: 0,
        referenceMax: 100,
        totalMeasurements: 0,
        trend: 'insufficient',
      );
    }
  }

  @override
  Future<List<BiomarkerSummary>> getAllBiomarkerSummaries() async {
    try {
      final res = await _client.getMeasurementHistory(userId, limit: 200);
      final measurements = res['measurements'] as List<dynamic>? ?? [];
      final byType = <String, List<double>>{};
      final unitByType = <String, String>{};
      for (final m in measurements) {
        final map = m as Map<String, dynamic>;
        final type = map['cartridge_type'] as String? ?? '';
        if (type.isEmpty) continue;
        byType.putIfAbsent(type, () => []);
        byType[type]!.add((map['primary_value'] as num?)?.toDouble() ?? 0.0);
        unitByType.putIfAbsent(type, () => map['unit'] as String? ?? '');
      }
      return byType.entries.map((e) {
        final values = e.value..sort();
        final avg = values.reduce((a, b) => a + b) / values.length;
        return BiomarkerSummary(
          biomarkerType: e.key,
          displayName: e.key,
          unit: unitByType[e.key] ?? '',
          latestValue: values.last,
          averageValue: avg,
          minValue: values.first,
          maxValue: values.last,
          referenceMin: 0,
          referenceMax: 200,
          totalMeasurements: values.length,
          trend: _computeTrend(values),
        );
      }).toList();
    } on DioException {
      return [];
    }
  }

  @override
  Future<ExportResult> exportData({
    required ExportFormat format,
    DateTime? from,
    DateTime? to,
    List<String>? biomarkerTypes,
  }) async {
    try {
      final res = await _client.exportToFHIR(userId: userId);
      return ExportResult(
        filePath: res['file_path'] as String? ?? res['fhir_json'] as String? ?? '',
        format: format,
        recordCount: res['record_count'] as int? ?? 0,
        exportedAt: DateTime.now(),
      );
    } on DioException {
      return ExportResult(
        filePath: '',
        format: format,
        recordCount: 0,
        exportedAt: DateTime.now(),
      );
    }
  }

  @override
  Future<int> getTotalMeasurementCount() async {
    try {
      final res = await _client.getMeasurementHistory(userId, limit: 1);
      return res['total_count'] as int? ?? 0;
    } on DioException {
      return 0;
    }
  }

  String _computeTrend(List<double> sortedValues) {
    if (sortedValues.length < 3) return 'insufficient';
    final halfIdx = sortedValues.length ~/ 2;
    final firstHalf = sortedValues.sublist(0, halfIdx);
    final secondHalf = sortedValues.sublist(halfIdx);
    final avgFirst = firstHalf.reduce((a, b) => a + b) / firstHalf.length;
    final avgSecond = secondHalf.reduce((a, b) => a + b) / secondHalf.length;
    final diff = avgSecond - avgFirst;
    if (diff.abs() < avgFirst * 0.05) return 'stable';
    return diff > 0 ? 'rising' : 'falling';
  }
}
