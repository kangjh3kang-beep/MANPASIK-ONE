import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/notification/data/notification_repository_rest.dart';
import 'package:manpasik/features/notification/domain/notification_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';

void main() {
  group('NotificationRepositoryRest', () {
    test('NotificationRepositoryRest는 NotificationRepository를 구현한다', () {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = NotificationRepositoryRest(client, 'user-1');
      expect(repo, isA<NotificationRepository>());
    });

    test('getNotifications는 DioException 시 빈 리스트를 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = NotificationRepositoryRest(client, 'user-1');
      final notifications = await repo.getNotifications();
      expect(notifications, isEmpty);
    });

    test('getUnreadCount는 DioException 시 0을 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = NotificationRepositoryRest(client, 'user-1');
      final count = await repo.getUnreadCount();
      expect(count, 0);
    });

    test('getNotifications 페이지네이션 파라미터', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = NotificationRepositoryRest(client, 'user-1');
      final notifications = await repo.getNotifications(page: 2, size: 10);
      expect(notifications, isEmpty);
    });
  });

  group('Notification 도메인 모델', () {
    test('AppNotification 생성 확인', () {
      final notification = AppNotification(
        id: 'n-1',
        type: NotificationType.health,
        title: '혈당 측정 알림',
        body: '오후 3시 혈당 측정 시간입니다',
        isRead: false,
        createdAt: DateTime(2026, 2, 19, 15, 0),
        deepLink: '/measurement',
      );
      expect(notification.id, 'n-1');
      expect(notification.isRead, isFalse);
      expect(notification.deepLink, '/measurement');
    });

    test('NotificationType enum 값 확인', () {
      expect(NotificationType.values, isNotEmpty);
    });
  });
}
