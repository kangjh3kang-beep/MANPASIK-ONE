/// 가족 건강 관리 도메인 모델 및 리포지토리
///
/// 가족 그룹 CRUD, 구성원 관리, 건강 데이터 공유

/// 가족 그룹
class FamilyGroup {
  final String id;
  final String name;
  final String ownerId;
  final List<FamilyMember> members;
  final DateTime createdAt;

  const FamilyGroup({
    required this.id,
    required this.name,
    required this.ownerId,
    required this.members,
    required this.createdAt,
  });
}

/// 가족 구성원
class FamilyMember {
  final String userId;
  final String displayName;
  final FamilyRole role;
  final SharingPermission permission;
  final DateTime? lastMeasurementAt;
  final String? latestHealthStatus; // 'normal', 'caution', 'alert'

  const FamilyMember({
    required this.userId,
    required this.displayName,
    required this.role,
    required this.permission,
    this.lastMeasurementAt,
    this.latestHealthStatus,
  });
}

/// 가족 역할
enum FamilyRole { owner, admin, member }

/// 데이터 공유 권한
class SharingPermission {
  final bool canViewResults;
  final bool canViewTrends;
  final bool canReceiveAlerts;
  final bool canSendReminders;

  const SharingPermission({
    this.canViewResults = false,
    this.canViewTrends = false,
    this.canReceiveAlerts = true,
    this.canSendReminders = false,
  });
}

/// 가족 초대
class FamilyInvitation {
  final String id;
  final String groupId;
  final String inviterName;
  final String inviteCode;
  final DateTime expiresAt;

  const FamilyInvitation({
    required this.id,
    required this.groupId,
    required this.inviterName,
    required this.inviteCode,
    required this.expiresAt,
  });
}

/// 가족 리포지토리 인터페이스
abstract class FamilyRepository {
  /// 가족 그룹 생성
  Future<FamilyGroup> createGroup(String name);

  /// 내 가족 그룹 목록
  Future<List<FamilyGroup>> getMyGroups();

  /// 가족 그룹 상세
  Future<FamilyGroup> getGroup(String groupId);

  /// 구성원 초대 링크 생성
  Future<FamilyInvitation> createInvitation(String groupId);

  /// 초대 수락
  Future<void> acceptInvitation(String inviteCode);

  /// 구성원 권한 변경
  Future<void> updateMemberPermission(String groupId, String userId, SharingPermission permission);

  /// 구성원 제거
  Future<void> removeMember(String groupId, String userId);

  /// 측정 리마인더 전송
  Future<void> sendMeasurementReminder(String groupId, String targetUserId);
}
