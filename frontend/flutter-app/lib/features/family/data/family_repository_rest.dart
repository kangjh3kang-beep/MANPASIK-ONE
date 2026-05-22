import 'package:dio/dio.dart';
import 'package:manpasik/features/family/domain/family_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';
import 'package:manpasik/core/network/offline_queue.dart';

/// REST Gateway를 사용하는 FamilyRepository 구현체
///
/// 가족 서비스는 현재 REST Gateway에 전용 엔드포인트가 없어
/// 일부 메서드는 placeholder 구현. 백엔드 확장 시 연동 가능.
class FamilyRepositoryRest implements FamilyRepository {
  FamilyRepositoryRest(this._client, {required this.userId, OfflineQueue? offlineQueue})
      : _offlineQueue = offlineQueue ?? OfflineQueue.instance;

  final ManPaSikRestClient _client;
  final String userId;
  final OfflineQueue _offlineQueue;

  @override
  Future<FamilyGroup> createGroup(String name) async {
    try {
      final res = await _client.createFamilyGroup(
        userId: userId,
        name: name,
      );
      return _mapGroup(res);
    } on DioException catch (e) {
      if (e.type == DioExceptionType.connectionError ||
          e.type == DioExceptionType.connectionTimeout) {
        await _offlineQueue.enqueue(OfflineRequest(
          method: 'POST',
          path: '/family/groups',
          body: {'user_id': userId, 'name': name},
          createdAt: DateTime.now(),
        ));
        // Optimistic result
        return FamilyGroup(
          id: 'pending-${DateTime.now().millisecondsSinceEpoch}',
          name: name,
          ownerId: userId,
          members: [],
          createdAt: DateTime.now(),
        );
      }
      rethrow;
    }
  }

  @override
  Future<List<FamilyGroup>> getMyGroups() async {
    try {
      final res = await _client.listFamilyGroups(userId);
      final list = res['groups'] as List<dynamic>? ?? [];
      return list.map((g) => _mapGroup(g as Map<String, dynamic>)).toList();
    } on DioException {
      return [];
    }
  }

  @override
  Future<FamilyGroup> getGroup(String groupId) async {
    final res = await _client.getFamilyGroup(groupId);
    return _mapGroup(res);
  }

  FamilyGroup _mapGroup(Map<String, dynamic> m) {
    final members = (m['members'] as List<dynamic>? ?? []).map((mem) {
      final mm = mem as Map<String, dynamic>;
      return FamilyMember(
        userId: mm['user_id'] as String? ?? '',
        displayName: mm['display_name'] as String? ?? mm['name'] as String? ?? '',
        role: _parseRole(mm['role']),
        permission: const SharingPermission(
          canViewResults: true,
          canViewTrends: true,
          canReceiveAlerts: true,
        ),
        lastMeasurementAt: mm['last_measurement_at'] != null
            ? DateTime.tryParse(mm['last_measurement_at'] as String)
            : null,
        latestHealthStatus: mm['health_status'] as String?,
      );
    }).toList();
    return FamilyGroup(
      id: m['group_id'] as String? ?? m['id'] as String? ?? '',
      name: m['name'] as String? ?? '',
      ownerId: m['owner_id'] as String? ?? '',
      members: members,
      createdAt: m['created_at'] != null
          ? DateTime.tryParse(m['created_at'] as String) ?? DateTime.now()
          : DateTime.now(),
    );
  }

  FamilyRole _parseRole(dynamic v) {
    if (v is String) {
      switch (v.toLowerCase()) {
        case 'owner': return FamilyRole.owner;
        case 'admin': return FamilyRole.admin;
      }
    }
    return FamilyRole.member;
  }

  @override
  Future<FamilyInvitation> createInvitation(String groupId) async {
    final res = await _client.inviteFamilyMember(
      groupId: groupId,
      inviteePhone: '',
    );
    return FamilyInvitation(
      id: res['invitation_id'] as String? ?? 'inv-${DateTime.now().millisecondsSinceEpoch}',
      groupId: groupId,
      inviterName: res['inviter_name'] as String? ?? '나',
      inviteCode: res['invite_code'] as String? ?? '',
      expiresAt: res['expires_at'] != null
          ? DateTime.tryParse(res['expires_at'] as String) ??
              DateTime.now().add(const Duration(days: 7))
          : DateTime.now().add(const Duration(days: 7)),
    );
  }

  @override
  Future<void> acceptInvitation(String inviteCode) async {
    await _client.respondToInvitation(
      invitationId: inviteCode,
      accept: true,
    );
  }

  @override
  Future<void> updateMemberPermission(
    String groupId,
    String targetUserId,
    SharingPermission permission,
  ) async {
    final body = {
      'user_id': targetUserId,
      'can_view_results': permission.canViewResults,
      'can_view_trends': permission.canViewTrends,
      'can_receive_alerts': permission.canReceiveAlerts,
      'can_send_reminders': permission.canSendReminders,
    };
    try {
      await _client.setSharingPreferences(groupId, body);
    } on DioException catch (e) {
      if (e.type == DioExceptionType.connectionError ||
          e.type == DioExceptionType.connectionTimeout) {
        await _offlineQueue.enqueue(OfflineRequest(
          method: 'PUT',
          path: '/family/groups/$groupId/sharing',
          body: body,
          createdAt: DateTime.now(),
        ));
        return; // Optimistic success
      }
      rethrow;
    }
  }

  @override
  Future<void> removeMember(String groupId, String targetUserId) async {
    try {
      await _client.removeFamilyMember(groupId, targetUserId);
    } on DioException catch (e) {
      if (e.type == DioExceptionType.connectionError ||
          e.type == DioExceptionType.connectionTimeout) {
        await _offlineQueue.enqueue(OfflineRequest(
          method: 'DELETE',
          path: '/family/groups/$groupId/members/$targetUserId',
          createdAt: DateTime.now(),
        ));
        return; // Optimistic success
      }
      rethrow;
    }
  }

  @override
  Future<void> sendMeasurementReminder(String groupId, String targetUserId) async {
    try {
      await _client.sendNotification(
        userId: targetUserId,
        title: '측정 리마인더',
        body: '가족 그룹에서 건강 측정을 요청했습니다.',
        type: 'family',
      );
    } on DioException catch (e) {
      if (e.type == DioExceptionType.connectionError ||
          e.type == DioExceptionType.connectionTimeout) {
        await _offlineQueue.enqueue(OfflineRequest(
          method: 'POST',
          path: '/notifications',
          body: {
            'user_id': targetUserId,
            'title': '측정 리마인더',
            'body': '가족 그룹에서 건강 측정을 요청했습니다.',
            'type': 'family',
          },
          createdAt: DateTime.now(),
        ));
        return; // Enqueued for later delivery
      }
      rethrow;
    }
  }
}
