/// 관리자 도메인 모델 및 리포지토리
///
/// 사용자 관리, 매출/재고 통계, 감사 로그, 컴플라이언스

/// 관리자 대시보드 통계
class AdminStats {
  final int totalUsers;
  final int activeUsers;
  final int totalDevices;
  final int activeDevices;
  final double monthlyRevenue;
  final int pendingOrders;

  const AdminStats({
    required this.totalUsers,
    required this.activeUsers,
    required this.totalDevices,
    required this.activeDevices,
    required this.monthlyRevenue,
    required this.pendingOrders,
  });
}

/// 감사 로그 엔트리
class AuditLogEntry {
  final String id;
  final String actorId;
  final String actorName;
  final String action;
  final String targetType;
  final String targetId;
  final DateTime timestamp;

  const AuditLogEntry({
    required this.id,
    required this.actorId,
    required this.actorName,
    required this.action,
    required this.targetType,
    required this.targetId,
    required this.timestamp,
  });
}

/// 관리자 리포지토리 인터페이스
abstract class AdminRepository {
  Future<AdminStats> getDashboardStats();
  Future<List<AuditLogEntry>> getAuditLogs({int page = 0, int size = 50});
  Future<Map<String, dynamic>> getRevenueStats({required String period});
  Future<Map<String, dynamic>> getInventoryStats();
  Future<void> updateUserRole(String userId, String role);
  Future<void> suspendUser(String userId, {required String reason});
}
