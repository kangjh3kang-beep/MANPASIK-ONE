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
      final points = measurements
          .map((m) {
            final map = m as Map<String, dynamic>;
            final point = _trendPointFromMeasurement(map);
            if (point == null ||
                point.timestamp.isBefore(from) ||
                point.timestamp.isAfter(to)) {
              return null;
            }
            if (biomarkerType.isNotEmpty &&
                point.biomarkerType != biomarkerType) {
              return null;
            }
            return point;
          })
          .whereType<TrendDataPoint>()
          .toList();
      points.sort((a, b) => a.timestamp.compareTo(b.timestamp));
      return points;
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
      final orderedValues = points.map((p) => p.value).toList();
      final sortedValues = [...orderedValues]..sort();
      final latestPoint = points.last;
      final avg =
          orderedValues.reduce((a, b) => a + b) / orderedValues.length;
      return BiomarkerSummary(
        biomarkerType: biomarkerType,
        displayName: biomarkerType,
        unit: latestPoint.unit,
        latestValue: latestPoint.value,
        averageValue: avg,
        minValue: sortedValues.first,
        maxValue: sortedValues.last,
        referenceMin: 0,
        referenceMax: 200,
        totalMeasurements: points.length,
        trend: _computeTrend(orderedValues),
        latestEvidenceStatus: latestPoint.evidenceStatus,
        latestDiagnosticReady: latestPoint.diagnosticReady,
        latestEvidenceGaps: latestPoint.evidenceGaps,
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
      final pointsByType = <String, List<TrendDataPoint>>{};
      final unitByType = <String, String>{};
      for (final m in measurements) {
        final map = m as Map<String, dynamic>;
        final point = _trendPointFromMeasurement(map);
        if (point == null) continue;
        final type = point.biomarkerType;
        if (type.isEmpty) continue;
        pointsByType.putIfAbsent(type, () => []).add(point);
        unitByType.putIfAbsent(type, () => point.unit);
      }
      return pointsByType.entries.map((e) {
        final points = e.value
          ..sort((a, b) => a.timestamp.compareTo(b.timestamp));
        final orderedValues = points.map((p) => p.value).toList();
        final sortedValues = [...orderedValues]..sort();
        final avg =
            orderedValues.reduce((a, b) => a + b) / orderedValues.length;
        final latestEvidence = points.last;
        return BiomarkerSummary(
          biomarkerType: e.key,
          displayName: e.key,
          unit: unitByType[e.key] ?? '',
          latestValue: latestEvidence.value,
          averageValue: avg,
          minValue: sortedValues.first,
          maxValue: sortedValues.last,
          referenceMin: 0,
          referenceMax: 200,
          totalMeasurements: orderedValues.length,
          trend: _computeTrend(orderedValues),
          latestEvidenceStatus: latestEvidence.evidenceStatus,
          latestDiagnosticReady: latestEvidence.diagnosticReady,
          latestEvidenceGaps: latestEvidence.evidenceGaps,
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
      final filePath = _stringField(
        res,
        'file_path',
        'filePath',
        fallback: _stringField(res, 'fhir_json', 'fhirJson'),
      );
      return ExportResult(
        filePath: filePath,
        format: format,
        recordCount: _intField(res, 'record_count', 'recordCount') ?? 0,
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
      return _intField(res, 'total_count', 'totalCount') ?? 0;
    } on DioException {
      return 0;
    }
  }

  String _computeTrend(List<double> orderedValues) {
    if (orderedValues.length < 3) return 'insufficient';
    final halfIdx = orderedValues.length ~/ 2;
    final firstHalf = orderedValues.sublist(0, halfIdx);
    final secondHalf = orderedValues.sublist(halfIdx);
    final avgFirst = firstHalf.reduce((a, b) => a + b) / firstHalf.length;
    final avgSecond = secondHalf.reduce((a, b) => a + b) / secondHalf.length;
    final diff = avgSecond - avgFirst;
    if (diff.abs() < avgFirst * 0.05) return 'stable';
    return diff > 0 ? 'rising' : 'falling';
  }
}

TrendDataPoint? _trendPointFromMeasurement(Map<String, dynamic> map) {
  final measuredAt = _stringField(map, 'measured_at', 'measuredAt');
  final timestamp =
      measuredAt.isNotEmpty ? DateTime.tryParse(measuredAt) : null;
  if (timestamp == null) {
    return null;
  }
  return TrendDataPoint(
    timestamp: timestamp,
    value: _numberField(map, 'primary_value', 'primaryValue') ?? 0.0,
    unit: _stringField(map, 'unit', 'unit'),
    biomarkerType: _stringField(map, 'cartridge_type', 'cartridgeType'),
    isWithinRange: _boolField(
      map,
      'is_within_range',
      'isWithinRange',
      fallback: true,
    ),
    evidenceStatus: _stringField(
      map,
      'evidence_status',
      'evidenceStatus',
      fallback: 'unknown',
    ),
    diagnosticReady: _boolField(map, 'diagnostic_ready', 'diagnosticReady'),
    evidenceGaps: _stringListField(map, 'evidence_gaps', 'evidenceGaps'),
  );
}

String _stringField(
  Map<String, dynamic> map,
  String snakeKey,
  String camelKey, {
  String fallback = '',
}) {
  final value = map[snakeKey] ?? map[camelKey];
  return value is String ? value : fallback;
}

double? _numberField(
  Map<String, dynamic> map,
  String snakeKey,
  String camelKey,
) {
  final value = map[snakeKey] ?? map[camelKey];
  return value is num ? value.toDouble() : null;
}

int? _intField(
  Map<String, dynamic> map,
  String snakeKey,
  String camelKey,
) {
  final value = map[snakeKey] ?? map[camelKey];
  return value is num ? value.toInt() : null;
}

bool _boolField(
  Map<String, dynamic> map,
  String snakeKey,
  String camelKey, {
  bool fallback = false,
}) {
  final value = map[snakeKey] ?? map[camelKey];
  return value is bool ? value : fallback;
}

List<String> _stringListField(
  Map<String, dynamic> map,
  String snakeKey,
  String camelKey,
) {
  final value = map[snakeKey] ?? map[camelKey];
  if (value is! List) {
    return const [];
  }
  return List.unmodifiable(value.whereType<String>());
}
