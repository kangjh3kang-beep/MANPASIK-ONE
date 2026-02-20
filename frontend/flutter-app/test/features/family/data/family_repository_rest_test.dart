import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/family/data/family_repository_rest.dart';
import 'package:manpasik/features/family/domain/family_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';

void main() {
  group('FamilyRepositoryRest', () {
    test('FamilyRepositoryRest는 FamilyRepository를 구현한다', () {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = FamilyRepositoryRest(client, userId: 'user-1');
      expect(repo, isA<FamilyRepository>());
    });

    test('getMyGroups는 DioException 시 빈 리스트를 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = FamilyRepositoryRest(client, userId: 'user-1');
      final groups = await repo.getMyGroups();
      expect(groups, isEmpty);
    });
  });

  group('Family 도메인 모델', () {
    test('FamilyRole enum 값 확인', () {
      expect(FamilyRole.values, isNotEmpty);
      expect(FamilyRole.values, contains(FamilyRole.owner));
      expect(FamilyRole.values, contains(FamilyRole.member));
    });

    test('SharingPermission 기본값 확인', () {
      const permission = SharingPermission();
      expect(permission.canViewResults, isFalse);
      expect(permission.canViewTrends, isFalse);
      expect(permission.canReceiveAlerts, isTrue);
      expect(permission.canSendReminders, isFalse);
    });

    test('SharingPermission 커스텀 값', () {
      const permission = SharingPermission(
        canViewResults: true,
        canViewTrends: true,
        canReceiveAlerts: false,
      );
      expect(permission.canViewResults, isTrue);
      expect(permission.canReceiveAlerts, isFalse);
    });

    test('FamilyMember 생성 확인', () {
      const member = FamilyMember(
        userId: 'user-1',
        displayName: '홍길동',
        role: FamilyRole.owner,
        permission: SharingPermission(
          canViewResults: true,
          canViewTrends: true,
          canReceiveAlerts: true,
        ),
      );
      expect(member.displayName, '홍길동');
      expect(member.role, FamilyRole.owner);
      expect(member.lastMeasurementAt, isNull);
    });

    test('FamilyInvitation 생성 확인', () {
      final invitation = FamilyInvitation(
        id: 'inv-1',
        groupId: 'grp-1',
        inviterName: '홍길동',
        inviteCode: 'ABC123',
        expiresAt: DateTime(2026, 2, 26),
      );
      expect(invitation.inviteCode, 'ABC123');
    });
  });
}
