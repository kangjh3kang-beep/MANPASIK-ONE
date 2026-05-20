import 'dart:convert';
import 'dart:io';

import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/measurement/data/measurement_repository_rest.dart';
import 'package:manpasik/features/measurement/domain/measurement_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';
import 'package:shared_preferences/shared_preferences.dart';

void main() {
  group('MeasurementRepositoryRest', () {
    setUp(() {
      SharedPreferences.setMockInitialValues({});
    });

    test('MeasurementRepositoryRest는 MeasurementRepository를 구현한다', () {
      final client =
          ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = MeasurementRepositoryRest(client);
      expect(repo, isA<MeasurementRepository>());
    });

    test('getHistory는 DioException 시 stale/error 결과를 반환한다', () async {
      final client =
          ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = MeasurementRepositoryRest(client);
      final result = await repo.getHistory(userId: 'user-1');
      expect(result.items, isEmpty);
      expect(result.totalCount, 0);
      expect(result.isStale, isTrue);
      expect(result.hasError, isTrue);
      expect(result.errorMessage, contains('측정 기록'));
    });

    test('getHistory 커스텀 limit/offset', () async {
      final client =
          ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = MeasurementRepositoryRest(client);
      final result =
          await repo.getHistory(userId: 'user-1', limit: 5, offset: 10);
      expect(result.items, isEmpty);
      expect(result.isStale, isTrue);
    });

    test('getHistory는 Gateway snake_case evidence fields를 매핑한다', () async {
      final server = await _startMeasurementHistoryServer({
        'measurements': [
          {
            'session_id': 'session-rest-1',
            'cartridge_type': 'glucose',
            'primary_value': 99.5,
            'unit': 'mg/dL',
            'evidence_status': 'research_only',
            'diagnostic_ready': false,
            'evidence_gaps': ['clinical_lock_required'],
          },
        ],
        'total_count': 1,
      });
      addTearDown(() => server.close(force: true));

      final repo = MeasurementRepositoryRest(
        ManPaSikRestClient(baseUrl: _baseUrl(server)),
      );

      final result = await repo.getHistory(userId: 'user-1');

      expect(result.totalCount, 1);
      expect(result.items.single.evidenceStatus, 'research_only');
      expect(result.items.single.diagnosticReady, isFalse);
      expect(
        result.items.single.evidenceGaps,
        contains('clinical_lock_required'),
      );
    });

    test('getHistory는 legacy camelCase evidence fields도 매핑한다', () async {
      final server = await _startMeasurementHistoryServer({
        'measurements': [
          {
            'sessionId': 'session-rest-2',
            'cartridgeType': 'glucose',
            'primaryValue': 98.4,
            'unit': 'mg/dL',
            'evidenceStatus': 'research_only',
            'diagnosticReady': false,
            'evidenceGaps': ['clinical_lock_required'],
          },
        ],
        'totalCount': 1,
      });
      addTearDown(() => server.close(force: true));

      final repo = MeasurementRepositoryRest(
        ManPaSikRestClient(baseUrl: _baseUrl(server)),
      );

      final result = await repo.getHistory(userId: 'user-1');

      expect(result.totalCount, 1);
      expect(result.items.single.sessionId, 'session-rest-2');
      expect(result.items.single.primaryValue, 98.4);
      expect(result.items.single.evidenceStatus, 'research_only');
      expect(result.items.single.diagnosticReady, isFalse);
      expect(
        result.items.single.evidenceGaps,
        contains('clinical_lock_required'),
      );
    });
  });

  group('Measurement 도메인 모델', () {
    test('StartSessionResult 생성 확인', () {
      const result = StartSessionResult(sessionId: 'sess-1');
      expect(result.sessionId, 'sess-1');
      expect(result.startedAt, isNull);
    });

    test('EndSessionResult 생성 확인', () {
      const result = EndSessionResult(
        sessionId: 'sess-1',
        totalMeasurements: 5,
      );
      expect(result.totalMeasurements, 5);
      expect(result.endedAt, isNull);
    });

    test('MeasurementHistoryResult 빈 결과 생성', () {
      const result = MeasurementHistoryResult(items: [], totalCount: 0);
      expect(result.items, isEmpty);
      expect(result.totalCount, 0);
      expect(result.isStale, isFalse);
      expect(result.hasError, isFalse);
    });

    test('MeasurementHistoryResult stale/error 생성', () {
      const result = MeasurementHistoryResult(
        items: [],
        totalCount: 0,
        isStale: true,
        errorMessage: 'sync failed',
      );
      expect(result.isStale, isTrue);
      expect(result.hasError, isTrue);
      expect(result.errorMessage, 'sync failed');
    });

    test('MeasurementHistoryItem 생성', () {
      final item = MeasurementHistoryItem(
        sessionId: 'sess-1',
        cartridgeType: 'glucose',
        primaryValue: 95.5,
        unit: 'mg/dL',
        measuredAt: DateTime(2026, 2, 19, 10, 0),
      );
      expect(item.primaryValue, 95.5);
      expect(item.cartridgeType, 'glucose');
      expect(item.evidenceStatus, 'unknown');
      expect(item.diagnosticReady, isFalse);
      expect(item.evidenceGaps, isEmpty);
    });
  });
}

String _baseUrl(HttpServer server) {
  return 'http://${server.address.address}:${server.port}/api/v1';
}

Future<HttpServer> _startMeasurementHistoryServer(
  Map<String, dynamic> responseBody,
) async {
  final server = await HttpServer.bind(InternetAddress.loopbackIPv4, 0);

  server.listen((request) async {
    request.response.headers.contentType = ContentType.json;
    request.response.statusCode =
        request.uri.path == '/api/v1/measurements/history'
            ? HttpStatus.ok
            : HttpStatus.notFound;
    request.response.write(jsonEncode(responseBody));
    await request.response.close();
  });

  return server;
}
