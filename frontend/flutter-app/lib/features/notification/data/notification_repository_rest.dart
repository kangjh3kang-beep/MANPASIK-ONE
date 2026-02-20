import 'package:dio/dio.dart';
import 'package:manpasik/core/services/rest_client.dart';
import 'package:manpasik/features/notification/domain/notification_repository.dart';

/// REST Gateway를 사용하는 NotificationRepository 구현체
class NotificationRepositoryRest implements NotificationRepository {
  NotificationRepositoryRest(this._client, this._userId);

  final ManPaSikRestClient _client;
  final String _userId;

  @override
  Future<List<AppNotification>> getNotifications({
    int page = 0,
    int size = 20,
  }) async {
    try {
      final res = await _client.listNotifications(
        _userId,
        limit: size,
        offset: page * size,
      );
      final items = res['notifications'] as List<dynamic>? ?? [];
      return items.map((n) {
        final m = n as Map<String, dynamic>;
        return AppNotification(
          id: m['notification_id'] as String? ?? '',
          type: _parseType(m['type'] as String?),
          title: m['title'] as String? ?? '',
          body: m['body'] as String? ?? '',
          isRead: m['is_read'] as bool? ?? false,
          createdAt: DateTime.tryParse(m['created_at'] as String? ?? '') ??
              DateTime.now(),
          deepLink: m['deep_link'] as String?,
        );
      }).toList();
    } on DioException {
      return [];
    }
  }

  @override
  Future<int> getUnreadCount() async {
    try {
      final res = await _client.getUnreadCount(_userId);
      return res['count'] as int? ?? 0;
    } on DioException {
      return 0;
    }
  }

  @override
  Future<void> markAsRead(String notificationId) async {
    await _client.markNotificationAsRead(notificationId);
  }

  @override
  Future<void> markAllAsRead() async {
    await _client.markAllNotificationsAsRead(_userId);
  }

  @override
  Future<void> deleteNotification(String notificationId) async {
    // Mark as read as a soft delete (no dedicated delete endpoint)
    await _client.markNotificationAsRead(notificationId);
  }

  @override
  Future<Map<NotificationType, bool>> getNotificationSettings() async {
    try {
      final res = await _client.getNotificationPreferences(_userId);
      final prefs = res['preferences'] as Map<String, dynamic>? ?? {};
      return {
        for (final type in NotificationType.values)
          type: prefs[type.name] as bool? ?? true,
      };
    } on DioException {
      return {for (final t in NotificationType.values) t: true};
    }
  }

  @override
  Future<void> updateNotificationSettings(
      Map<NotificationType, bool> settings) async {
    final prefs = <String, dynamic>{};
    for (final entry in settings.entries) {
      prefs[entry.key.name] = entry.value;
    }
    await _client.updateNotificationPreferences(_userId, prefs);
  }

  NotificationType _parseType(String? type) {
    switch (type) {
      case 'measurement':
        return NotificationType.measurement;
      case 'health':
        return NotificationType.health;
      case 'community':
        return NotificationType.community;
      case 'family':
        return NotificationType.family;
      case 'market':
        return NotificationType.market;
      default:
        return NotificationType.system;
    }
  }
}
