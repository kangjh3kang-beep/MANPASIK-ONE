import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/admin/data/admin_repository_rest.dart';
import 'package:manpasik/features/admin/domain/admin_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';

void main() {
  group('AdminRepositoryRest', () {
    test('AdminRepositoryRest는 AdminRepository를 구현한다', () {
      final client =
          ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = AdminRepositoryRest(client);
      expect(repo, isA<AdminRepository>());
    });

    test('getDashboardStats는 DioException 시 기본 AdminStats를 반환한다', () async {
      final client =
          ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = AdminRepositoryRest(client);
      final stats = await repo.getDashboardStats();
      expect(stats, isA<AdminStats>());
      expect(stats.totalUsers, 0);
      expect(stats.activeUsers, 0);
      expect(stats.totalDevices, 0);
      expect(stats.activeDevices, 0);
      expect(stats.monthlyRevenue, 0.0);
      expect(stats.pendingOrders, 0);
    });

    test('getAuditLogs는 DioException 시 빈 리스트를 반환한다', () async {
      final client =
          ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = AdminRepositoryRest(client);
      final logs = await repo.getAuditLogs();
      expect(logs, isEmpty);
    });

    test('getAuditLogs page/size 파라미터 확인', () async {
      final client =
          ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = AdminRepositoryRest(client);
      final logs = await repo.getAuditLogs(page: 1, size: 25);
      expect(logs, isEmpty);
    });

    test('getRevenueStats는 DioException 시 빈 맵을 반환한다', () async {
      final client =
          ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = AdminRepositoryRest(client);
      final stats = await repo.getRevenueStats(period: 'monthly');
      expect(stats, isEmpty);
    });

    test('getInventoryStats는 DioException 시 빈 맵을 반환한다', () async {
      final client =
          ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = AdminRepositoryRest(client);
      final stats = await repo.getInventoryStats();
      expect(stats, isEmpty);
    });
  });

  group('Admin 도메인 모델', () {
    test('AdminStats 생성 확인', () {
      const stats = AdminStats(
        totalUsers: 1000,
        activeUsers: 500,
        totalDevices: 200,
        activeDevices: 150,
        monthlyRevenue: 5000000.0,
        pendingOrders: 25,
      );
      expect(stats.totalUsers, 1000);
      expect(stats.activeUsers, 500);
      expect(stats.totalDevices, 200);
      expect(stats.activeDevices, 150);
      expect(stats.monthlyRevenue, 5000000.0);
      expect(stats.pendingOrders, 25);
    });

    test('AuditLogEntry 생성 확인', () {
      final now = DateTime.now();
      final entry = AuditLogEntry(
        id: 'log-1',
        actorId: 'admin-1',
        actorName: '관리자',
        action: 'user.suspend',
        targetType: 'user',
        targetId: 'user-100',
        timestamp: now,
      );
      expect(entry.id, 'log-1');
      expect(entry.actorId, 'admin-1');
      expect(entry.actorName, '관리자');
      expect(entry.action, 'user.suspend');
      expect(entry.targetType, 'user');
      expect(entry.targetId, 'user-100');
      expect(entry.timestamp, now);
    });
  });
}
