import 'dart:convert';
import 'dart:io';

import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/data_hub/data/data_hub_repository_rest.dart';
import 'package:manpasik/features/data_hub/domain/data_hub_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';
import 'package:shared_preferences/shared_preferences.dart';

void main() {
  group('DataHubRepositoryRest', () {
    setUp(() {
      SharedPreferences.setMockInitialValues({});
    });

    test('DataHubRepositoryRest는 DataHubRepository를 구현한다', () {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = DataHubRepositoryRest(client, userId: 'user-1');
      expect(repo, isA<DataHubRepository>());
    });

    test('getTrendData는 DioException 시 빈 리스트를 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = DataHubRepositoryRest(client, userId: 'user-1');
      final data = await repo.getTrendData(
        biomarkerType: 'glucose',
        from: DateTime(2026, 1, 1),
        to: DateTime(2026, 2, 19),
      );
      expect(data, isEmpty);
    });

    test('getTrendData는 snake_case evidence fields를 보존한다', () async {
      final server = await _startHistoryServer({
        'measurements': [
          _snakeMeasurement(),
        ],
        'total_count': 1,
      });
      addTearDown(() => server.close(force: true));

      final repo = DataHubRepositoryRest(
        ManPaSikRestClient(baseUrl: _baseUrl(server)),
        userId: 'user-1',
      );

      final points = await repo.getTrendData(
        biomarkerType: 'glucose',
        from: DateTime.now().subtract(const Duration(days: 1)),
        to: DateTime.now().add(const Duration(days: 1)),
      );

      expect(points, hasLength(1));
      expect(points.single.evidenceStatus, 'research_only');
      expect(points.single.diagnosticReady, isFalse);
      expect(points.single.evidenceGaps, contains('clinical_lock_required'));
    });

    test('getTrendData는 legacy camelCase evidence fields도 보존한다', () async {
      final server = await _startHistoryServer({
        'measurements': [
          _camelMeasurement(),
        ],
        'totalCount': 1,
      });
      addTearDown(() => server.close(force: true));

      final repo = DataHubRepositoryRest(
        ManPaSikRestClient(baseUrl: _baseUrl(server)),
        userId: 'user-1',
      );

      final points = await repo.getTrendData(
        biomarkerType: 'glucose',
        from: DateTime.now().subtract(const Duration(days: 1)),
        to: DateTime.now().add(const Duration(days: 1)),
      );

      expect(points, hasLength(1));
      expect(points.single.evidenceStatus, 'research_only');
      expect(points.single.diagnosticReady, isFalse);
      expect(points.single.evidenceGaps, contains('clinical_lock_required'));
    });

    test('getBiomarkerSummary는 latest evidence metadata를 보존한다', () async {
      final server = await _startHistoryServer({
        'measurements': [
          _snakeMeasurement(),
        ],
        'total_count': 1,
      });
      addTearDown(() => server.close(force: true));

      final repo = DataHubRepositoryRest(
        ManPaSikRestClient(baseUrl: _baseUrl(server)),
        userId: 'user-1',
      );

      final summary = await repo.getBiomarkerSummary('glucose');

      expect(summary.totalMeasurements, 1);
      expect(summary.latestEvidenceStatus, 'research_only');
      expect(summary.latestDiagnosticReady, isFalse);
      expect(
        summary.latestEvidenceGaps,
        contains('clinical_lock_required'),
      );
    });

    test('getTrendData는 measured_at 기준 오름차순으로 정렬한다', () async {
      final now = DateTime.now().toUtc();
      final newer = now.subtract(const Duration(days: 1));
      final older = now.subtract(const Duration(days: 3));
      final server = await _startHistoryServer({
        'measurements': [
          _snakeMeasurement(measuredAt: newer, value: 120.0),
          _snakeMeasurement(measuredAt: older, value: 80.0),
        ],
        'total_count': 2,
      });
      addTearDown(() => server.close(force: true));

      final repo = DataHubRepositoryRest(
        ManPaSikRestClient(baseUrl: _baseUrl(server)),
        userId: 'user-1',
      );

      final points = await repo.getTrendData(
        biomarkerType: 'glucose',
        from: now.subtract(const Duration(days: 7)),
        to: now.add(const Duration(days: 1)),
      );

      expect(points, hasLength(2));
      expect(points.first.timestamp, older);
      expect(points.last.timestamp, newer);
    });

    test('getBiomarkerSummary는 measured_at 최신 evidence를 선택한다', () async {
      final now = DateTime.now().toUtc();
      final newer = now.subtract(const Duration(days: 1));
      final older = now.subtract(const Duration(days: 3));
      final server = await _startHistoryServer({
        'measurements': [
          _snakeMeasurement(
            measuredAt: newer,
            value: 120.0,
            evidenceStatus: 'clinical_locked',
            diagnosticReady: true,
            evidenceGaps: const [],
          ),
          _snakeMeasurement(
            measuredAt: older,
            value: 80.0,
            evidenceStatus: 'research_only',
            diagnosticReady: false,
            evidenceGaps: const ['clinical_lock_required'],
          ),
        ],
        'total_count': 2,
      });
      addTearDown(() => server.close(force: true));

      final repo = DataHubRepositoryRest(
        ManPaSikRestClient(baseUrl: _baseUrl(server)),
        userId: 'user-1',
      );

      final summary = await repo.getBiomarkerSummary('glucose');

      expect(summary.latestValue, 120.0);
      expect(summary.latestEvidenceStatus, 'clinical_locked');
      expect(summary.latestDiagnosticReady, isTrue);
      expect(summary.latestEvidenceGaps, isEmpty);
    });

    test('getAllBiomarkerSummaries는 measured_at 최신 evidence를 선택한다', () async {
      final now = DateTime.now().toUtc();
      final newer = now.subtract(const Duration(hours: 4));
      final older = now.subtract(const Duration(days: 2));
      final server = await _startHistoryServer({
        'measurements': [
          _snakeMeasurement(
            measuredAt: older,
            value: 130.0,
            evidenceStatus: 'research_only',
            diagnosticReady: false,
            evidenceGaps: const ['clinical_lock_required'],
          ),
          _snakeMeasurement(
            measuredAt: newer,
            value: 95.0,
            evidenceStatus: 'clinical_locked',
            diagnosticReady: true,
            evidenceGaps: const [],
          ),
        ],
        'total_count': 2,
      });
      addTearDown(() => server.close(force: true));

      final repo = DataHubRepositoryRest(
        ManPaSikRestClient(baseUrl: _baseUrl(server)),
        userId: 'user-1',
      );

      final summaries = await repo.getAllBiomarkerSummaries();

      expect(summaries, hasLength(1));
      expect(summaries.single.latestValue, 95.0);
      expect(summaries.single.latestEvidenceStatus, 'clinical_locked');
      expect(summaries.single.latestDiagnosticReady, isTrue);
      expect(summaries.single.latestEvidenceGaps, isEmpty);
    });

    test('getBiomarkerSummary는 timestamp 순서 기준 하강 추세를 계산한다', () async {
      final now = DateTime.now().toUtc();
      final older = now.subtract(const Duration(days: 3));
      final middle = now.subtract(const Duration(days: 2));
      final newer = now.subtract(const Duration(days: 1));
      final server = await _startHistoryServer({
        'measurements': [
          _snakeMeasurement(measuredAt: newer, value: 70.0),
          _snakeMeasurement(measuredAt: older, value: 130.0),
          _snakeMeasurement(measuredAt: middle, value: 100.0),
        ],
        'total_count': 3,
      });
      addTearDown(() => server.close(force: true));

      final repo = DataHubRepositoryRest(
        ManPaSikRestClient(baseUrl: _baseUrl(server)),
        userId: 'user-1',
      );

      final summary = await repo.getBiomarkerSummary('glucose');
      final allSummaries = await repo.getAllBiomarkerSummaries();

      expect(summary.trend, 'falling');
      expect(allSummaries.single.trend, 'falling');
    });

    test('getTotalMeasurementCount는 legacy camelCase totalCount를 보존한다',
        () async {
      final server = await _startHistoryServer({
        'measurements': [],
        'totalCount': 7,
      });
      addTearDown(() => server.close(force: true));

      final repo = DataHubRepositoryRest(
        ManPaSikRestClient(baseUrl: _baseUrl(server)),
        userId: 'user-1',
      );

      final count = await repo.getTotalMeasurementCount();

      expect(count, 7);
    });

    test('exportData는 camelCase filePath와 recordCount를 보존한다', () async {
      final server = await _startExportServer({
        'filePath': '/tmp/export.fhir.json',
        'recordCount': 3,
      });
      addTearDown(() => server.close(force: true));

      final repo = DataHubRepositoryRest(
        ManPaSikRestClient(baseUrl: _baseUrl(server)),
        userId: 'user-1',
      );

      final result = await repo.exportData(format: ExportFormat.json);

      expect(result.filePath, '/tmp/export.fhir.json');
      expect(result.recordCount, 3);
      expect(result.format, ExportFormat.json);
    });
  });

  group('DataHub 도메인 모델', () {
    test('TrendDataPoint 생성 확인', () {
      final point = TrendDataPoint(
        timestamp: DateTime(2026, 2, 19, 10, 0),
        value: 95.0,
        unit: 'mg/dL',
        biomarkerType: 'glucose',
        isWithinRange: true,
      );
      expect(point.value, 95.0);
      expect(point.isWithinRange, isTrue);
    });
  });
}

String _baseUrl(HttpServer server) {
  return 'http://${server.address.address}:${server.port}/api/v1';
}

Future<HttpServer> _startHistoryServer(Map<String, dynamic> body) async {
  final server = await HttpServer.bind(InternetAddress.loopbackIPv4, 0);

  server.listen((request) async {
    request.response.headers.contentType = ContentType.json;
    request.response.statusCode =
        request.uri.path == '/api/v1/measurements/history'
            ? HttpStatus.ok
            : HttpStatus.notFound;
    request.response.write(jsonEncode(body));
    await request.response.close();
  });

  return server;
}

Future<HttpServer> _startExportServer(Map<String, dynamic> body) async {
  final server = await HttpServer.bind(InternetAddress.loopbackIPv4, 0);

  server.listen((request) async {
    request.response.headers.contentType = ContentType.json;
    request.response.statusCode =
        request.uri.path == '/api/v1/health-records/export/fhir'
            ? HttpStatus.ok
            : HttpStatus.notFound;
    request.response.write(jsonEncode(body));
    await request.response.close();
  });

  return server;
}

Map<String, dynamic> _snakeMeasurement({
  DateTime? measuredAt,
  double value = 99.5,
  String evidenceStatus = 'research_only',
  bool diagnosticReady = false,
  List<String> evidenceGaps = const ['clinical_lock_required'],
}) {
  return {
    'session_id': 'session-datahub-1',
    'cartridge_type': 'glucose',
    'primary_value': value,
    'unit': 'mg/dL',
    'measured_at': (measuredAt ?? DateTime.now().toUtc()).toIso8601String(),
    'evidence_status': evidenceStatus,
    'diagnostic_ready': diagnosticReady,
    'evidence_gaps': evidenceGaps,
  };
}

Map<String, dynamic> _camelMeasurement({
  DateTime? measuredAt,
  double value = 98.4,
  String evidenceStatus = 'research_only',
  bool diagnosticReady = false,
  List<String> evidenceGaps = const ['clinical_lock_required'],
}) {
  return {
    'sessionId': 'session-datahub-2',
    'cartridgeType': 'glucose',
    'primaryValue': value,
    'unit': 'mg/dL',
    'measuredAt': (measuredAt ?? DateTime.now().toUtc()).toIso8601String(),
    'evidenceStatus': evidenceStatus,
    'diagnosticReady': diagnosticReady,
    'evidenceGaps': evidenceGaps,
  };
}
