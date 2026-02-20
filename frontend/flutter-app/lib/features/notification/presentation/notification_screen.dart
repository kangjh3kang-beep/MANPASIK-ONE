import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';

/// 알림 목록 데이터
class NotificationItem {
  final String id;
  final String title;
  final String body;
  final String type;
  final bool isRead;
  final DateTime createdAt;

  const NotificationItem({
    required this.id,
    required this.title,
    required this.body,
    required this.type,
    required this.isRead,
    required this.createdAt,
  });
}

/// 알림 목록 Provider
final notificationsListProvider = FutureProvider<List<NotificationItem>>((ref) async {
  final userId = ref.watch(authProvider).userId;
  if (userId == null || userId.isEmpty) return [];
  try {
    final client = ref.read(restClientProvider);
    final res = await client.listNotifications(userId);
    final list = res['notifications'] as List<dynamic>? ?? [];
    return list.map((n) {
      final m = n as Map<String, dynamic>;
      return NotificationItem(
        id: m['id'] as String? ?? m['notification_id'] as String? ?? '',
        title: m['title'] as String? ?? '',
        body: m['body'] as String? ?? m['message'] as String? ?? '',
        type: m['type'] as String? ?? '',
        isRead: m['is_read'] as bool? ?? false,
        createdAt: m['created_at'] != null
            ? DateTime.tryParse(m['created_at'] as String) ?? DateTime.now()
            : DateTime.now(),
      );
    }).toList();
  } catch (_) {
    return [];
  }
});

/// 탭별 필터 타입
const _tabFilters = <String?>[
  null, // 전체
  'health_alert,measurement', // 건강
  'order,community', // 시스템
  'family', // 가족
];

/// 알림 센터 화면
class NotificationScreen extends ConsumerWidget {
  const NotificationScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final notiAsync = ref.watch(notificationsListProvider);

    return DefaultTabController(
      length: 4,
      child: Scaffold(
        appBar: AppBar(
          title: const Text('알림'),
          centerTitle: true,
          bottom: TabBar(
            isScrollable: false,
            labelColor: theme.colorScheme.primary,
            unselectedLabelColor: theme.colorScheme.onSurfaceVariant,
            indicatorColor: theme.colorScheme.primary,
            tabs: const [
              Tab(text: '전체'),
              Tab(text: '건강'),
              Tab(text: '시스템'),
              Tab(text: '가족'),
            ],
          ),
        ),
        body: notiAsync.when(
          data: (notifications) {
            return TabBarView(
              children: List.generate(4, (tabIndex) {
                final filter = _tabFilters[tabIndex];
                final filtered = filter == null
                    ? notifications
                    : notifications.where((n) {
                        final types = filter.split(',');
                        return types.contains(n.type);
                      }).toList();

                if (filtered.isEmpty) {
                  return Center(
                    child: Column(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Icon(Icons.notifications_off_outlined, size: 48, color: theme.colorScheme.outline),
                        const SizedBox(height: 12),
                        Text(
                          '알림이 없습니다',
                          style: theme.textTheme.bodyLarge?.copyWith(color: theme.colorScheme.onSurfaceVariant),
                        ),
                      ],
                    ),
                  );
                }

                return RefreshIndicator(
                  onRefresh: () async {
                    ref.invalidate(notificationsListProvider);
                    ref.invalidate(unreadNotificationCountProvider);
                  },
                  child: ListView.builder(
                    itemCount: filtered.length,
                    itemBuilder: (context, index) {
                      return _NotificationTile(noti: filtered[index]);
                    },
                  ),
                );
              }),
            );
          },
          loading: () => const Center(child: CircularProgressIndicator()),
          error: (_, __) => Center(
            child: Text('알림을 불러올 수 없습니다', style: theme.textTheme.bodyLarge),
          ),
        ),
      ),
    );
  }
}

class _NotificationTile extends ConsumerWidget {
  const _NotificationTile({required this.noti});
  final NotificationItem noti;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final icon = switch (noti.type) {
      'measurement' => Icons.query_stats,
      'health_alert' => Icons.warning_amber,
      'family' => Icons.family_restroom,
      'community' => Icons.forum,
      'order' => Icons.local_shipping,
      _ => Icons.notifications,
    };
    final iconColor = switch (noti.type) {
      'health_alert' => Colors.red,
      'measurement' => theme.colorScheme.primary,
      _ => theme.colorScheme.secondary,
    };

    return ListTile(
      leading: CircleAvatar(
        backgroundColor: iconColor.withOpacity(0.1),
        child: Icon(icon, color: iconColor, size: 20),
      ),
      title: Text(
        noti.title,
        style: theme.textTheme.bodyMedium?.copyWith(
          fontWeight: noti.isRead ? FontWeight.normal : FontWeight.bold,
        ),
      ),
      subtitle: Text(noti.body, maxLines: 2, overflow: TextOverflow.ellipsis),
      trailing: Text(
        _formatTime(noti.createdAt),
        style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.outline),
      ),
      tileColor: noti.isRead ? null : theme.colorScheme.primaryContainer.withOpacity(0.1),
      onTap: () async {
        if (!noti.isRead) {
          try {
            await ref.read(restClientProvider).markNotificationAsRead(noti.id);
            ref.invalidate(notificationsListProvider);
            ref.invalidate(unreadNotificationCountProvider);
          } catch (_) {}
        }
        if (noti.type == 'health_alert' && context.mounted) {
          context.push('/family/alert/${noti.id}');
        }
      },
    );
  }

  String _formatTime(DateTime dt) {
    final diff = DateTime.now().difference(dt);
    if (diff.inMinutes < 60) return '${diff.inMinutes}분 전';
    if (diff.inHours < 24) return '${diff.inHours}시간 전';
    if (diff.inDays < 7) return '${diff.inDays}일 전';
    return '${dt.month}/${dt.day}';
  }
}
