import 'dart:async';
import 'package:flutter/foundation.dart';

import 'package:manpasik/core/services/rest_client.dart';

/// 푸시 알림 서비스 인터페이스 (B3)
///
/// FCM / APNs 실제 SDK 연동 시 이 인터페이스를 구현합니다.
/// SDK 미설치 시 PollingNotificationService로 REST 폴링 기반 알림 대체.
abstract class PushNotificationService {
  Future<void> initialize();
  Future<String?> getToken();
  Stream<NotificationPayload> get onNotification;
  Future<void> dispose();
}

/// REST 폴링 기반 알림 서비스 (FCM 미사용 시 대체)
class PollingNotificationService implements PushNotificationService {
  PollingNotificationService({
    required this.restClient,
    required this.userId,
    this.pollInterval = const Duration(seconds: 30),
  });

  final ManPaSikRestClient restClient;
  final String userId;
  final Duration pollInterval;

  Timer? _timer;
  int _lastKnownCount = 0;
  final _controller = StreamController<NotificationPayload>.broadcast();

  @override
  Future<void> initialize() async {
    debugPrint('[PollingNotification] 초기화: userId=$userId, interval=${pollInterval.inSeconds}s');
    _startPolling();
  }

  void _startPolling() {
    _timer?.cancel();
    _timer = Timer.periodic(pollInterval, (_) => _poll());
    _poll(); // 즉시 첫 폴링
  }

  Future<void> _poll() async {
    if (userId.isEmpty) return;
    try {
      final res = await restClient.getUnreadCount(userId);
      final count = res['count'] as int? ?? res['unread_count'] as int? ?? 0;

      if (count > _lastKnownCount && _lastKnownCount > 0) {
        _controller.add(NotificationPayload(
          title: '새 알림',
          body: '${count - _lastKnownCount}개의 새 알림이 있습니다',
          data: {'unread_count': count},
        ));
      }
      _lastKnownCount = count;
    } catch (e) {
      debugPrint('[PollingNotification] 폴링 실패: $e');
    }
  }

  @override
  Future<String?> getToken() async => 'polling_${userId}_${DateTime.now().millisecondsSinceEpoch}';

  @override
  Stream<NotificationPayload> get onNotification => _controller.stream;

  @override
  Future<void> dispose() async {
    _timer?.cancel();
    await _controller.close();
  }
}

class NotificationPayload {
  final String title;
  final String body;
  final Map<String, dynamic>? data;

  const NotificationPayload({
    required this.title,
    required this.body,
    this.data,
  });
}
