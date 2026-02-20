import 'package:dio/dio.dart';
import 'package:manpasik/core/services/rest_client.dart';
import 'package:manpasik/features/admin/domain/admin_repository.dart';

/// REST Gateway를 사용하는 AdminRepository 구현체
class AdminRepositoryRest implements AdminRepository {
  AdminRepositoryRest(this._client);

  final ManPaSikRestClient _client;

  @override
  Future<AdminStats> getDashboardStats() async {
    try {
      final res = await _client.getSystemStats();
      return AdminStats(
        totalUsers: res['total_users'] as int? ?? 0,
        activeUsers: res['active_users'] as int? ?? 0,
        totalDevices: res['total_devices'] as int? ?? 0,
        activeDevices: res['active_devices'] as int? ?? 0,
        monthlyRevenue:
            (res['monthly_revenue'] as num?)?.toDouble() ?? 0.0,
        pendingOrders: res['pending_orders'] as int? ?? 0,
      );
    } on DioException {
      return const AdminStats(
        totalUsers: 0,
        activeUsers: 0,
        totalDevices: 0,
        activeDevices: 0,
        monthlyRevenue: 0,
        pendingOrders: 0,
      );
    }
  }

  @override
  Future<List<AuditLogEntry>> getAuditLogs({
    int page = 0,
    int size = 50,
  }) async {
    try {
      final res = await _client.getAuditLog(
        limit: size,
        offset: page * size,
      );
      final entries = res['entries'] as List<dynamic>? ?? [];
      return entries.map((e) {
        final m = e as Map<String, dynamic>;
        return AuditLogEntry(
          id: m['id'] as String? ?? '',
          actorId: m['actor_id'] as String? ?? '',
          actorName: m['actor_name'] as String? ?? '',
          action: m['action'] as String? ?? '',
          targetType: m['target_type'] as String? ?? '',
          targetId: m['target_id'] as String? ?? '',
          timestamp:
              DateTime.tryParse(m['timestamp'] as String? ?? '') ??
                  DateTime.now(),
        );
      }).toList();
    } on DioException {
      return [];
    }
  }

  @override
  Future<Map<String, dynamic>> getRevenueStats({
    required String period,
  }) async {
    try {
      return await _client.getRevenueStats(period: period);
    } on DioException {
      return {};
    }
  }

  @override
  Future<Map<String, dynamic>> getInventoryStats() async {
    try {
      return await _client.getInventoryStats();
    } on DioException {
      return {};
    }
  }

  @override
  Future<void> updateUserRole(String userId, String role) async {
    await _client.adminChangeRole(userId, role);
  }

  @override
  Future<void> suspendUser(String userId, {required String reason}) async {
    await _client.adminBulkAction(
      userIds: [userId],
      action: 'suspend',
    );
  }
}
