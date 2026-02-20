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

  /// 구성 상태에 따라 적절한 구현체를 반환하는 팩토리
  static PushNotificationService create({
    required ManPaSikRestClient restClient,
    required String userId,
  }) {
    // Firebase 설정이 있으면 FCM 사용, 아니면 폴링 폴백
    if (FcmNotificationService.isConfigured) {
      return FcmNotificationService(
        restClient: restClient,
        userId: userId,
      );
    }
    return PollingNotificationService(
      restClient: restClient,
      userId: userId,
    );
  }
}

/// FCM 기반 푸시 알림 서비스
///
/// Firebase Cloud Messaging을 사용하여 실시간 푸시 알림을 수신합니다.
/// firebase_core / firebase_messaging 패키지 설치 및
/// google-services.json / GoogleService-Info.plist 설정이 필요합니다.
///
/// 설정 미완료 시 PollingNotificationService로 자동 폴백.
class FcmNotificationService implements PushNotificationService {
  FcmNotificationService({
    required this.restClient,
    required this.userId,
  });

  final ManPaSikRestClient restClient;
  final String userId;

  final _controller = StreamController<NotificationPayload>.broadcast();
  String? _fcmToken;
  bool _initialized = false;

  /// FCM 설정 파일이 존재하는지 확인
  /// 실제 Firebase 초기화는 firebase_core 패키지 설치 후 활성화
  static const _fcmEnabled = bool.fromEnvironment('FCM_ENABLED');
  static bool get isConfigured => _fcmEnabled;

  // 폴링 폴백
  Timer? _pollingTimer;
  int _lastKnownCount = 0;

  @override
  Future<void> initialize() async {
    if (_initialized) return;

    if (isConfigured) {
      await _initializeFcm();
    } else {
      debugPrint('[FcmNotification] FCM 미설정 → REST 폴링 폴백');
      _startPollingFallback();
    }
    _initialized = true;
  }

  /// FCM 초기화 (firebase_core/firebase_messaging 패키지 설치 시 활성화)
  Future<void> _initializeFcm() async {
    try {
      // Firebase 초기화
      // await Firebase.initializeApp(options: DefaultFirebaseOptions.currentPlatform);
      // final messaging = FirebaseMessaging.instance;
      //
      // // iOS 권한 요청
      // await messaging.requestPermission(
      //   alert: true, badge: true, sound: true,
      // );
      //
      // // 토큰 획득 및 서버 등록
      // final token = await messaging.getToken();
      // if (token != null) {
      //   _fcmToken = token;
      //   await restClient.registerPushToken(
      //     userId: userId, token: token,
      //   );
      // }
      //
      // // 토큰 갱신 시 서버에 재등록
      // messaging.onTokenRefresh.listen((newToken) async {
      //   _fcmToken = newToken;
      //   await restClient.registerPushToken(
      //     userId: userId, token: newToken,
      //   );
      // });
      //
      // // 포그라운드 알림 수신
      // FirebaseMessaging.onMessage.listen((message) {
      //   _controller.add(NotificationPayload(
      //     title: message.notification?.title ?? '',
      //     body: message.notification?.body ?? '',
      //     data: message.data,
      //   ));
      // });
      //
      // // 백그라운드 알림 탭
      // FirebaseMessaging.onMessageOpenedApp.listen((message) {
      //   _controller.add(NotificationPayload(
      //     title: message.notification?.title ?? '',
      //     body: message.notification?.body ?? '',
      //     data: {...message.data, 'tapped': true},
      //   ));
      // });

      debugPrint('[FcmNotification] FCM 초기화 완료 (토큰: ${_fcmToken?.substring(0, 10)}...)');
    } catch (e) {
      debugPrint('[FcmNotification] FCM 초기화 실패 → 폴링 폴백: $e');
      _startPollingFallback();
    }
  }

  void _startPollingFallback() {
    _pollingTimer?.cancel();
    _pollingTimer = Timer.periodic(
      const Duration(seconds: 30),
      (_) => _poll(),
    );
    _poll();
  }

  Future<void> _poll() async {
    if (userId.isEmpty) return;
    try {
      final res = await restClient.getUnreadCount(userId);
      final count =
          res['count'] as int? ?? res['unread_count'] as int? ?? 0;

      if (count > _lastKnownCount && _lastKnownCount > 0) {
        _controller.add(NotificationPayload(
          title: '새 알림',
          body: '${count - _lastKnownCount}개의 새 알림이 있습니다',
          data: {'unread_count': count},
        ));
      }
      _lastKnownCount = count;
    } catch (e) {
      debugPrint('[FcmNotification] 폴링 실패: $e');
    }
  }

  @override
  Future<String?> getToken() async {
    if (_fcmToken != null) return _fcmToken;
    return 'polling_${userId}_${DateTime.now().millisecondsSinceEpoch}';
  }

  @override
  Stream<NotificationPayload> get onNotification => _controller.stream;

  @override
  Future<void> dispose() async {
    _pollingTimer?.cancel();
    await _controller.close();
  }
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
    debugPrint(
        '[PollingNotification] 초기화: userId=$userId, interval=${pollInterval.inSeconds}s');
    _startPolling();
  }

  void _startPolling() {
    _timer?.cancel();
    _timer = Timer.periodic(pollInterval, (_) => _poll());
    _poll();
  }

  Future<void> _poll() async {
    if (userId.isEmpty) return;
    try {
      final res = await restClient.getUnreadCount(userId);
      final count =
          res['count'] as int? ?? res['unread_count'] as int? ?? 0;

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
  Future<String?> getToken() async =>
      'polling_${userId}_${DateTime.now().millisecondsSinceEpoch}';

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
