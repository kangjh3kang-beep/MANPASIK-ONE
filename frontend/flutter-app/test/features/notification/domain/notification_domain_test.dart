import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/notification/domain/notification_repository.dart';

void main() {
  group('NotificationType enum', () {
    test('값 목록 확인 (6개)', () {
      expect(NotificationType.values, hasLength(6));
      expect(NotificationType.values, contains(NotificationType.measurement));
      expect(NotificationType.values, contains(NotificationType.health));
      expect(NotificationType.values, contains(NotificationType.community));
      expect(NotificationType.values, contains(NotificationType.system));
      expect(NotificationType.values, contains(NotificationType.family));
      expect(NotificationType.values, contains(NotificationType.market));
    });
  });

  group('AppNotification', () {
    test('기본 생성 확인', () {
      final now = DateTime.now();
      final notification = AppNotification(
        id: 'notif-1',
        type: NotificationType.health,
        title: '혈당 측정 알림',
        body: '오전 혈당 측정 시간입니다.',
        isRead: false,
        createdAt: now,
      );
      expect(notification.id, 'notif-1');
      expect(notification.type, NotificationType.health);
      expect(notification.title, '혈당 측정 알림');
      expect(notification.body, '오전 혈당 측정 시간입니다.');
      expect(notification.isRead, isFalse);
      expect(notification.createdAt, now);
      expect(notification.deepLink, isNull);
    });

    test('deepLink 포함 생성', () {
      final notification = AppNotification(
        id: 'notif-2',
        type: NotificationType.community,
        title: '새 댓글',
        body: '게시물에 새 댓글이 달렸습니다.',
        isRead: true,
        createdAt: DateTime(2026, 4, 19),
        deepLink: '/community/post/123',
      );
      expect(notification.deepLink, '/community/post/123');
      expect(notification.isRead, isTrue);
    });

    test('읽음/안읽음 상태 확인', () {
      final unread = AppNotification(
        id: 'n-1',
        type: NotificationType.system,
        title: '시스템 알림',
        body: '업데이트가 있습니다.',
        isRead: false,
        createdAt: DateTime.now(),
      );
      final read = AppNotification(
        id: 'n-2',
        type: NotificationType.system,
        title: '시스템 알림',
        body: '업데이트 완료.',
        isRead: true,
        createdAt: DateTime.now(),
      );
      expect(unread.isRead, isFalse);
      expect(read.isRead, isTrue);
    });

    test('모든 NotificationType으로 생성 가능', () {
      for (final type in NotificationType.values) {
        final notification = AppNotification(
          id: 'test-${type.name}',
          type: type,
          title: '테스트',
          body: '테스트 알림',
          isRead: false,
          createdAt: DateTime.now(),
        );
        expect(notification.type, type);
      }
    });
  });
}
