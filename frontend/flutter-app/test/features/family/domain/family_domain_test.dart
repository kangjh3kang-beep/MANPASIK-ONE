import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/family/domain/family_repository.dart';

void main() {
  group('FamilyRole enum', () {
    test('값 목록 확인 (3개)', () {
      expect(FamilyRole.values, hasLength(3));
      expect(FamilyRole.values, contains(FamilyRole.owner));
      expect(FamilyRole.values, contains(FamilyRole.admin));
      expect(FamilyRole.values, contains(FamilyRole.member));
    });
  });

  group('SharingPermission', () {
    test('기본값 확인', () {
      const perm = SharingPermission();
      expect(perm.canViewResults, isFalse);
      expect(perm.canViewTrends, isFalse);
      expect(perm.canReceiveAlerts, isTrue); // 기본값 true
      expect(perm.canSendReminders, isFalse);
    });

    test('모든 권한 활성화', () {
      const perm = SharingPermission(
        canViewResults: true,
        canViewTrends: true,
        canReceiveAlerts: true,
        canSendReminders: true,
      );
      expect(perm.canViewResults, isTrue);
      expect(perm.canViewTrends, isTrue);
      expect(perm.canReceiveAlerts, isTrue);
      expect(perm.canSendReminders, isTrue);
    });
  });

  group('FamilyMember', () {
    test('기본 생성 확인', () {
      const member = FamilyMember(
        userId: 'user-1',
        displayName: '홍길동',
        role: FamilyRole.owner,
        permission: SharingPermission(
          canViewResults: true,
          canViewTrends: true,
          canReceiveAlerts: true,
          canSendReminders: true,
        ),
      );
      expect(member.userId, 'user-1');
      expect(member.displayName, '홍길동');
      expect(member.role, FamilyRole.owner);
      expect(member.lastMeasurementAt, isNull);
      expect(member.latestHealthStatus, isNull);
    });

    test('건강 상태 포함 생성', () {
      final now = DateTime.now();
      final member = FamilyMember(
        userId: 'user-2',
        displayName: '김건강',
        role: FamilyRole.member,
        permission: const SharingPermission(),
        lastMeasurementAt: now,
        latestHealthStatus: 'caution',
      );
      expect(member.lastMeasurementAt, now);
      expect(member.latestHealthStatus, 'caution');
    });
  });

  group('FamilyGroup', () {
    test('기본 생성 확인', () {
      final now = DateTime.now();
      final group = FamilyGroup(
        id: 'grp-1',
        name: '우리 가족',
        ownerId: 'user-1',
        members: const [],
        createdAt: now,
      );
      expect(group.id, 'grp-1');
      expect(group.name, '우리 가족');
      expect(group.ownerId, 'user-1');
      expect(group.members, isEmpty);
      expect(group.createdAt, now);
    });

    test('구성원이 있는 그룹 생성', () {
      final group = FamilyGroup(
        id: 'grp-2',
        name: '부모님',
        ownerId: 'user-1',
        members: const [
          FamilyMember(
            userId: 'user-1',
            displayName: '홍길동',
            role: FamilyRole.owner,
            permission: SharingPermission(
              canViewResults: true,
              canViewTrends: true,
            ),
          ),
          FamilyMember(
            userId: 'user-2',
            displayName: '어머니',
            role: FamilyRole.member,
            permission: SharingPermission(
              canReceiveAlerts: true,
            ),
          ),
        ],
        createdAt: DateTime(2026, 1, 1),
      );
      expect(group.members, hasLength(2));
      expect(group.members[0].role, FamilyRole.owner);
      expect(group.members[1].role, FamilyRole.member);
    });
  });

  group('FamilyInvitation', () {
    test('생성 확인', () {
      final invitation = FamilyInvitation(
        id: 'inv-1',
        groupId: 'grp-1',
        inviterName: '홍길동',
        inviteCode: 'ABC123',
        expiresAt: DateTime(2026, 4, 26),
      );
      expect(invitation.id, 'inv-1');
      expect(invitation.groupId, 'grp-1');
      expect(invitation.inviterName, '홍길동');
      expect(invitation.inviteCode, 'ABC123');
      expect(invitation.expiresAt.day, 26);
    });
  });
}
