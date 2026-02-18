/// 알림 도메인 모델 및 리포지토리
///
/// 푸시/인앱 알림, 읽음 처리, 알림 설정

/// 알림 타입
enum NotificationType { measurement, health, community, system, family, market }

/// 알림
class AppNotification {
  final String id;
  final NotificationType type;
  final String title;
  final String body;
  final bool isRead;
  final DateTime createdAt;
  final String? deepLink;

  const AppNotification({
    required this.id,
    required this.type,
    required this.title,
    required this.body,
    required this.isRead,
    required this.createdAt,
    this.deepLink,
  });
}

/// 알림 리포지토리 인터페이스
abstract class NotificationRepository {
  Future<List<AppNotification>> getNotifications({int page = 0, int size = 20});
  Future<int> getUnreadCount();
  Future<void> markAsRead(String notificationId);
  Future<void> markAllAsRead();
  Future<void> deleteNotification(String notificationId);
  Future<Map<NotificationType, bool>> getNotificationSettings();
  Future<void> updateNotificationSettings(Map<NotificationType, bool> settings);
}
