import 'package:dio/dio.dart';
import 'package:manpasik/features/family/domain/family_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';

/// REST Gateway를 사용하는 FamilyRepository 구현체
///
/// 가족 서비스는 현재 REST Gateway에 전용 엔드포인트가 없어
/// 일부 메서드는 placeholder 구현. 백엔드 확장 시 연동 가능.
class FamilyRepositoryRest implements FamilyRepository {
  FamilyRepositoryRest(this._client, {required this.userId});

  final ManPaSikRestClient _client;
  final String userId;

  @override
  Future<FamilyGroup> createGroup(String name) async {
    // Placeholder: family REST endpoints not yet in gateway
    return FamilyGroup(
      id: 'group-${DateTime.now().millisecondsSinceEpoch}',
      name: name,
      ownerId: userId,
      members: [
        FamilyMember(
          userId: userId,
          displayName: '나',
          role: FamilyRole.owner,
          permission: const SharingPermission(
            canViewResults: true,
            canViewTrends: true,
            canReceiveAlerts: true,
            canSendReminders: true,
          ),
        ),
      ],
      createdAt: DateTime.now(),
    );
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
    final res = await _client.getFamilyGroupReport(groupId);
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
    return FamilyInvitation(
      id: 'inv-${DateTime.now().millisecondsSinceEpoch}',
      groupId: groupId,
      inviterName: '나',
      inviteCode: 'MPK-${DateTime.now().millisecondsSinceEpoch.toRadixString(36).toUpperCase()}',
      expiresAt: DateTime.now().add(const Duration(days: 7)),
    );
  }

  @override
  Future<void> acceptInvitation(String inviteCode) async {
    // Placeholder
  }

  @override
  Future<void> updateMemberPermission(
    String groupId,
    String userId,
    SharingPermission permission,
  ) async {
    // Placeholder
  }

  @override
  Future<void> removeMember(String groupId, String userId) async {
    // Placeholder
  }

  @override
  Future<void> sendMeasurementReminder(String groupId, String targetUserId) async {
    // Placeholder
  }
}
