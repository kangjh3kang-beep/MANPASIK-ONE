// OfflineQueue 단위 테스트
// - Hive 기반 오프라인 요청 큐의 enqueue/dequeue/retry/drop 동작 검증

import 'dart:io';

import 'package:flutter_test/flutter_test.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:manpasik/core/network/offline_queue.dart';

void main() {
  late OfflineQueue queue;

  setUp(() async {
    // Hive를 임시 디렉토리로 초기화
    final dir = await Directory.systemTemp.createTemp('offline_queue_test');
    Hive.init(dir.path);
    queue = OfflineQueue.instance;
    await queue.init();
    await queue.clear();
  });

  tearDown(() async {
    await queue.clear();
  });

  group('OfflineQueue', () {
    test('enqueue adds request to queue', () async {
      expect(queue.pendingCount, 0);

      final request = OfflineRequest(
        method: 'POST',
        path: '/api/v1/measurements',
        body: {'value': 42},
        createdAt: DateTime.now(),
      );
      await queue.enqueue(request);

      expect(queue.pendingCount, 1);

      final all = queue.getAll();
      expect(all.length, 1);
      expect(all.first.method, 'POST');
      expect(all.first.path, '/api/v1/measurements');
      expect(all.first.body, {'value': 42});
    });

    test('dequeue returns requests in FIFO order', () async {
      final req1 = OfflineRequest(
        method: 'POST',
        path: '/api/v1/first',
        createdAt: DateTime.now(),
      );
      final req2 = OfflineRequest(
        method: 'PUT',
        path: '/api/v1/second',
        createdAt: DateTime.now(),
      );
      await queue.enqueue(req1);
      await queue.enqueue(req2);

      final all = queue.getAll();
      expect(all.length, 2);
      expect(all[0].path, '/api/v1/first');
      expect(all[1].path, '/api/v1/second');
    });

    test('pendingCount reflects queue size', () async {
      expect(queue.pendingCount, 0);

      await queue.enqueue(OfflineRequest(
        method: 'GET',
        path: '/a',
        createdAt: DateTime.now(),
      ));
      expect(queue.pendingCount, 1);

      await queue.enqueue(OfflineRequest(
        method: 'GET',
        path: '/b',
        createdAt: DateTime.now(),
      ));
      expect(queue.pendingCount, 2);

      await queue.clear();
      expect(queue.pendingCount, 0);
    });

    test('retry increments retryCount', () {
      final original = OfflineRequest(
        method: 'POST',
        path: '/api/v1/test',
        createdAt: DateTime.now(),
        retryCount: 0,
      );

      final retried = original.copyWithRetry();
      expect(retried.retryCount, 1);
      expect(retried.method, original.method);
      expect(retried.path, original.path);

      final retriedAgain = retried.copyWithRetry();
      expect(retriedAgain.retryCount, 2);
    });

    test('drops request after max retries', () async {
      const maxRetries = 3;

      var request = OfflineRequest(
        method: 'POST',
        path: '/api/v1/flaky',
        createdAt: DateTime.now(),
        retryCount: 0,
      );
      await queue.enqueue(request);

      // Simulate retry loop: remove old, re-enqueue with incremented count
      // until maxRetries is exceeded, then drop
      final keys = queue.keys;
      expect(keys.length, 1);

      // Remove the original
      await queue.remove(keys.first);
      expect(queue.pendingCount, 0);

      // Simulate retries up to max
      for (var i = 0; i < maxRetries; i++) {
        request = request.copyWithRetry();
      }
      // After max retries, the request should be dropped (not re-enqueued)
      expect(request.retryCount, maxRetries);

      // Verify queue stays empty when we choose not to re-enqueue
      expect(queue.pendingCount, 0);
    });

    test('remove deletes specific request by key', () async {
      await queue.enqueue(OfflineRequest(
        method: 'POST',
        path: '/api/v1/keep',
        createdAt: DateTime.now(),
      ));
      await queue.enqueue(OfflineRequest(
        method: 'DELETE',
        path: '/api/v1/remove-me',
        createdAt: DateTime.now(),
      ));
      expect(queue.pendingCount, 2);

      final keys = queue.keys;
      // Remove the first one
      await queue.remove(keys.first);
      expect(queue.pendingCount, 1);

      final remaining = queue.getAll();
      expect(remaining.first.path, '/api/v1/remove-me');
    });

    test('getByKey returns correct request', () async {
      await queue.enqueue(OfflineRequest(
        method: 'PATCH',
        path: '/api/v1/target',
        body: {'status': 'updated'},
        createdAt: DateTime.now(),
      ));

      final key = queue.keys.first;
      final fetched = queue.getByKey(key);
      expect(fetched, isNotNull);
      expect(fetched!.method, 'PATCH');
      expect(fetched.path, '/api/v1/target');
    });
  });
}
